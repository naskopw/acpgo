package acp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"strconv"
	"sync"
	"sync/atomic"
)

// StdioError represents a JSON-RPC error returned by a stdio subprocess.
type StdioError struct {
	Code    int
	Message string
	Data    json.RawMessage
}

func (e *StdioError) Error() string {
	return fmt.Sprintf("JSON-RPC error %d: %s", e.Code, e.Message)
}

// StdioTransport manages JSON-RPC communication with a subprocess over stdin/stdout.
// This implements the ACP stdio transport specification.
type StdioTransport struct {
	logger *slog.Logger
	cancel context.CancelFunc
	cmd    *exec.Cmd

	stdin  io.WriteCloser
	stdout io.Reader

	writeMu   sync.Mutex
	pendingMu sync.Mutex
	pending   map[string]chan Response
	nextID    atomic.Int64

	notifHandler func(method string, params json.RawMessage)
	notifMu      sync.RWMutex
}

// NewStdioTransport starts the given binary as a subprocess and returns a transport
// that communicates with it over JSON-RPC using the ACP stdio transport.
func NewStdioTransport(ctx context.Context, logger *slog.Logger, binary string, args ...string) (*StdioTransport, error) {
	ctx, cancel := context.WithCancel(ctx)

	cmd := exec.CommandContext(ctx, binary, args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		cancel()
		return nil, fmt.Errorf("start: %w", err)
	}

	t := &StdioTransport{
		logger:  logger.With("component", "acpgo"),
		cancel:  cancel,
		cmd:     cmd,
		stdin:   stdin,
		stdout:  stdout,
		pending: make(map[string]chan Response),
	}

	go t.readLoop()

	return t, nil
}

// NewStdioTransportWithIO creates a transport from existing io.ReadWriter pairs.
func NewStdioTransportWithIO(ctx context.Context, logger *slog.Logger, stdin io.WriteCloser, stdout io.Reader) *StdioTransport {
	_, cancel := context.WithCancel(ctx)
	t := &StdioTransport{
		logger:  logger.With("component", "acpgo"),
		cancel:  cancel,
		stdin:   stdin,
		stdout:  stdout,
		pending: make(map[string]chan Response),
	}
	go t.readLoop()
	return t
}

// Call sends a JSON-RPC request and waits for the response.
func (t *StdioTransport) Call(ctx context.Context, method string, params interface{}) (json.RawMessage, error) {
	id := strconv.FormatInt(t.nextID.Add(1), 10)

	paramsRaw, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal params: %w", err)
	}

	req := Request{
		JSONRPC: "2.0",
		ID:      json.RawMessage(id),
		Method:  method,
		Params:  paramsRaw,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	ch := make(chan Response, 1)

	t.pendingMu.Lock()
	t.pending[id] = ch
	t.pendingMu.Unlock()

	defer func() {
		t.pendingMu.Lock()
		delete(t.pending, id)
		t.pendingMu.Unlock()
	}()

	t.logger.Debug("stdio call", "method", method, "id", id)

	t.writeMu.Lock()
	_, err = fmt.Fprintln(t.stdin, string(data))
	t.writeMu.Unlock()
	if err != nil {
		return nil, fmt.Errorf("write request: %w", err)
	}

	select {
	case resp := <-ch:
		t.logger.Debug("stdio response", "method", method, "id", id, "hasError", resp.Error != nil)
		if resp.Error != nil {
			return nil, &StdioError{
				Code:    resp.Error.Code,
				Message: resp.Error.Message,
				Data:    resp.Error.Data,
			}
		}
		return resp.Result, nil
	case <-ctx.Done():
		t.logger.Debug("stdio call cancelled", "method", method, "id", id)
		_ = t.Notify(NotificationCancelRequest, CancelRequestNotification{RequestID: id})
		return nil, ctx.Err()
	}
}

// Notify sends a JSON-RPC notification (no response expected).
func (t *StdioTransport) Notify(method string, params interface{}) error {
	paramsRaw, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshal params: %w", err)
	}

	notif := Notification{
		JSONRPC: "2.0",
		Method:  method,
		Params:  paramsRaw,
	}

	data, err := json.Marshal(notif)
	if err != nil {
		return fmt.Errorf("marshal notification: %w", err)
	}

	t.logger.Debug("stdio notify", "method", method)

	t.writeMu.Lock()
	_, err = fmt.Fprintln(t.stdin, string(data))
	t.writeMu.Unlock()
	if err != nil {
		return fmt.Errorf("write notification: %w", err)
	}

	return nil
}

// Close shuts down the transport, killing the subprocess and unblocking readLoop.
func (t *StdioTransport) Close() error {
	t.cancel()

	var errs []error

	if err := t.stdin.Close(); err != nil {
		errs = append(errs, fmt.Errorf("close stdin: %w", err))
	}

	if t.cmd != nil {
		if waitErr := t.cmd.Wait(); waitErr != nil {
			var exitErr *exec.ExitError
			if errors.As(waitErr, &exitErr) && !exitErr.Exited() {
				// Process was killed by signal (likely our context cancellation) — expected
			} else {
				errs = append(errs, fmt.Errorf("subprocess wait: %w", waitErr))
			}
		}
	} else {
		if rc, ok := t.stdout.(io.Closer); ok {
			if err := rc.Close(); err != nil {
				t.logger.Debug("close stdout", "error", err)
			}
		}
	}

	return errors.Join(errs...)
}

// SetNotificationHandler registers a handler for incoming notifications.
func (t *StdioTransport) SetNotificationHandler(fn func(method string, params json.RawMessage)) {
	t.notifMu.Lock()
	defer t.notifMu.Unlock()
	t.notifHandler = fn
}

func (t *StdioTransport) readLoop() {
	scanner := bufio.NewScanner(t.stdout)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		lineCopy := make([]byte, len(line))
		copy(lineCopy, line)

		var raw struct {
			ID     json.RawMessage `json:"id"`
			Method string          `json:"method"`
		}
		if err := json.Unmarshal(lineCopy, &raw); err != nil {
			t.logger.Error("failed to parse JSON-RPC message", "error", err)
			continue
		}

		switch {
		case len(raw.ID) > 0 && string(raw.ID) != "null":
			var resp Response
			if err := json.Unmarshal(lineCopy, &resp); err != nil {
				t.logger.Error("failed to parse JSON-RPC response", "error", err)
				continue
			}

			key := string(bytes.TrimSpace(resp.ID))
			t.pendingMu.Lock()
			ch, ok := t.pending[key]
			t.pendingMu.Unlock()

			if ok {
				t.logger.Debug("stdio response received", "id", key, "hasError", resp.Error != nil)
				ch <- resp
			} else {
				t.logger.Warn("response for unknown request ID", "id", key)
			}

		case raw.Method != "":
			t.logger.Debug("stdio notification received", "method", raw.Method)

			var notif struct {
				Params json.RawMessage `json:"params"`
			}
			if err := json.Unmarshal(lineCopy, &notif); err != nil {
				t.logger.Error("failed to parse notification", "error", err)
				continue
			}

			t.notifMu.RLock()
			fn := t.notifHandler
			t.notifMu.RUnlock()

			if fn != nil {
				fn(raw.Method, notif.Params)
			}

		default:
			t.logger.Warn("unrecognized JSON-RPC message", "line", string(lineCopy))
		}
	}

	if err := scanner.Err(); err != nil {
		t.logger.Error("read loop scanner error", "error", err)
	}
}
