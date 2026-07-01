package acp_test

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/naskopw/acpgo"
	"github.com/stretchr/testify/require"
)

type mockSubprocess struct {
	stdin  io.Reader
	stdout io.Writer
	done   chan struct{}
}

func closeQuietly(c io.Closer) {
	_ = c.Close()
}

func newMockTransport(t *testing.T) (*acp.StdioTransport, *mockSubprocess) {
	t.Helper()

	subprocStdinR, subprocStdinW, err := os.Pipe()
	require.NoError(t, err)

	subprocStdoutR, subprocStdoutW, err := os.Pipe()
	require.NoError(t, err)

	transport := acp.NewStdioTransportWithIO(
		context.Background(),
		slog.Default(),
		subprocStdinW,
		subprocStdoutR,
	)

	mock := &mockSubprocess{
		stdin:  subprocStdinR,
		stdout: subprocStdoutW,
		done:   make(chan struct{}),
	}

	t.Cleanup(func() {
		closeQuietly(subprocStdinR)
		closeQuietly(subprocStdoutW)
	})

	return transport, mock
}

func (m *mockSubprocess) runEchoHandler(t *testing.T) {
	t.Helper()
	go func() {
		defer close(m.done)
		scanner := bufio.NewScanner(m.stdin)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for scanner.Scan() {
			line := scanner.Bytes()
			var raw struct {
				ID     json.RawMessage `json:"id"`
				Method string          `json:"method"`
			}
			if err := json.Unmarshal(line, &raw); err != nil {
				continue
			}
			if len(raw.ID) > 0 && string(raw.ID) != "null" {
				resp := fmt.Sprintf(`{"jsonrpc":"2.0","id":%s,"result":{"echo":true}}`, string(raw.ID))
				_, _ = fmt.Fprintln(m.stdout, resp)
			}
		}
	}()
}

func (m *mockSubprocess) runEchoHandlerWithDelay(t *testing.T, delay time.Duration) {
	t.Helper()
	go func() {
		defer close(m.done)
		scanner := bufio.NewScanner(m.stdin)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for scanner.Scan() {
			line := scanner.Bytes()
			var raw struct {
				ID     json.RawMessage `json:"id"`
				Method string          `json:"method"`
			}
			if err := json.Unmarshal(line, &raw); err != nil {
				continue
			}
			if len(raw.ID) > 0 && string(raw.ID) != "null" {
				time.Sleep(delay)
				resp := fmt.Sprintf(`{"jsonrpc":"2.0","id":%s,"result":{"echo":true}}`, string(raw.ID))
				_, _ = fmt.Fprintln(m.stdout, resp)
			}
		}
	}()
}

func (m *mockSubprocess) runErrorHandler(t *testing.T) {
	t.Helper()
	go func() {
		defer close(m.done)
		scanner := bufio.NewScanner(m.stdin)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for scanner.Scan() {
			line := scanner.Bytes()
			var raw struct {
				ID     json.RawMessage `json:"id"`
				Method string          `json:"method"`
			}
			if err := json.Unmarshal(line, &raw); err != nil {
				continue
			}
			if len(raw.ID) > 0 && string(raw.ID) != "null" {
				resp := fmt.Sprintf(`{"jsonrpc":"2.0","id":%s,"error":{"code":-32603,"message":"internal error"}}`, string(raw.ID))
				_, _ = fmt.Fprintln(m.stdout, resp)
			}
		}
	}()
}

func TestStdioCallBasic(t *testing.T) {
	transport, mock := newMockTransport(t)
	mock.runEchoHandler(t)
	defer func() { _ = transport.Close() }()

	result, err := transport.Call(context.Background(), "test/method", map[string]string{"key": "value"})
	require.NoError(t, err)
	require.NotNil(t, result)

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal(result, &parsed))
	require.Equal(t, true, parsed["echo"])
}

func TestStdioCallErrorResponse(t *testing.T) {
	transport, mock := newMockTransport(t)
	mock.runErrorHandler(t)
	defer func() { _ = transport.Close() }()

	result, err := transport.Call(context.Background(), "test/method", nil)
	require.Error(t, err)
	require.Nil(t, result)

	var stdioErr *acp.StdioError
	require.True(t, errors.As(err, &stdioErr))
	require.Equal(t, -32603, stdioErr.Code)
	require.Equal(t, "internal error", stdioErr.Message)
}

func TestStdioNotify(t *testing.T) {
	transport, mock := newMockTransport(t)
	mock.runEchoHandler(t)
	defer func() { _ = transport.Close() }()

	err := transport.Notify("test/notification", map[string]string{"msg": "hello"})
	require.NoError(t, err)
}

func TestStdioNotificationHandler(t *testing.T) {
	transport, mock := newMockTransport(t)
	defer func() { _ = transport.Close() }()

	received := make(chan struct {
		method string
		params json.RawMessage
	}, 1)

	transport.SetNotificationHandler(func(method string, params json.RawMessage) {
		received <- struct {
			method string
			params json.RawMessage
		}{method, params}
	})

	notif := `{"jsonrpc":"2.0","method":"test/notification","params":{"msg":"hi"}}`
	_, err := fmt.Fprintln(mock.stdout, notif)
	require.NoError(t, err)

	select {
	case got := <-received:
		require.Equal(t, "test/notification", got.method)
		var params map[string]string
		require.NoError(t, json.Unmarshal(got.params, &params))
		require.Equal(t, "hi", params["msg"])
	case <-time.After(2 * time.Second):
		t.Fatal("notification handler not called")
	}
}

func TestStdioCallContextCancel(t *testing.T) {
	transport, mock := newMockTransport(t)
	mock.runEchoHandlerWithDelay(t, 500*time.Millisecond)
	defer func() { _ = transport.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	result, err := transport.Call(ctx, "test/method", nil)
	require.Error(t, err)
	require.Nil(t, result)
	require.True(t, errors.Is(err, context.DeadlineExceeded))
}

func TestStdioConcurrentCalls(t *testing.T) {
	transport, mock := newMockTransport(t)
	mock.runEchoHandler(t)
	defer func() { _ = transport.Close() }()

	const n = 10
	var wg sync.WaitGroup
	wg.Add(n)
	errs := make([]error, n)

	for i := 0; i < n; i++ {
		go func(idx int) {
			defer wg.Done()
			_, errs[idx] = transport.Call(context.Background(), "test/method", map[string]int{"idx": idx})
		}(i)
	}

	wg.Wait()
	for i, err := range errs {
		require.NoError(t, err, "call %d failed", i)
	}
}

func TestStdioCloseUnblocksReadLoop(t *testing.T) {
	transport, mock := newMockTransport(t)
	mock.runEchoHandler(t)

	_, err := transport.Call(context.Background(), "test/method", nil)
	require.NoError(t, err)

	err = transport.Close()
	require.NoError(t, err)

	select {
	case <-mock.done:
	case <-time.After(2 * time.Second):
		t.Fatal("mock subprocess goroutine did not exit after Close — readLoop goroutine leak")
	}
}

type closeTrackingReader struct {
	r       io.Reader
	closed  bool
	mu      sync.Mutex
}

func (c *closeTrackingReader) Read(p []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return 0, os.ErrClosed
	}
	return c.r.Read(p)
}

func (c *closeTrackingReader) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
	return nil
}

func TestStdioCloseClosesStdout(t *testing.T) {
	stdinR, stdinW, err := os.Pipe()
	require.NoError(t, err)
	defer func() { _ = stdinR.Close() }()

	stdoutR, stdoutW, err := os.Pipe()
	require.NoError(t, err)
	defer func() { _ = stdoutW.Close() }()

	trackingReader := &closeTrackingReader{r: stdoutR}

	transport := acp.NewStdioTransportWithIO(
		context.Background(),
		slog.Default(),
		stdinW,
		trackingReader,
	)

	err = transport.Close()
	require.NoError(t, err)

	require.True(t, trackingReader.closed, "Close() should close stdout to unblock readLoop goroutine")
}

func TestStdioMalformedJSON(t *testing.T) {
	transport, mock := newMockTransport(t)
	defer func() { _ = transport.Close() }()

	handlerCalled := make(chan struct{}, 1)
	transport.SetNotificationHandler(func(method string, params json.RawMessage) {
		handlerCalled <- struct{}{}
	})

	_, err := fmt.Fprintln(mock.stdout, `{invalid json}`)
	require.NoError(t, err)

	_, err = fmt.Fprintln(mock.stdout, `{"jsonrpc":"2.0","method":"test/valid","params":{}}`)
	require.NoError(t, err)

	select {
	case <-handlerCalled:
	case <-time.After(2 * time.Second):
		t.Fatal("valid notification after malformed JSON was not processed — readLoop stopped on bad input")
	}
}

func TestStdioResponseForUnknownID(t *testing.T) {
	transport, mock := newMockTransport(t)
	mock.runEchoHandler(t)
	defer func() { _ = transport.Close() }()

	orphanResp := `{"jsonrpc":"2.0","id":"99999","result":{}}`
	_, err := fmt.Fprintln(mock.stdout, orphanResp)
	require.NoError(t, err)

	result, err := transport.Call(context.Background(), "test/method", nil)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestStdioCallAfterClose(t *testing.T) {
	transport, mock := newMockTransport(t)
	mock.runEchoHandler(t)

	require.NoError(t, transport.Close())

	_, err := transport.Call(context.Background(), "test/method", nil)
	require.Error(t, err)
}

func TestStdioCloseKilledSubprocessNoError(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess test in short mode")
	}

	transport, err := acp.NewStdioTransport(
		context.Background(),
		slog.Default(),
		"sleep", "10",
	)
	require.NoError(t, err)

	closeErr := transport.Close()
	require.NoError(t, closeErr, "Close should not return error when subprocess is killed by context cancellation")
}

func TestStdioCloseReturnsStdinError(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess test in short mode")
	}

	transport, err := acp.NewStdioTransport(
		context.Background(),
		slog.Default(),
		"sleep", "10",
	)
	require.NoError(t, err)

	require.NoError(t, transport.Close())
	err = transport.Close()
	require.Error(t, err, "second Close should return error from already-closed stdin")
}

func TestStdioNewStdioTransportBasic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess test in short mode")
	}

	echoScript := `#!/bin/sh
while IFS= read -r line; do
  id=$(echo "$line" | sed -n 's/.*"id":\([0-9]*\).*/\1/p')
  if [ -n "$id" ]; then
    printf '{"jsonrpc":"2.0","id":%s,"result":{"ok":true}}\n' "$id"
  fi
done
`
	tmpDir := t.TempDir()
	scriptPath := tmpDir + "/echo.sh"
	require.NoError(t, os.WriteFile(scriptPath, []byte(echoScript), 0755))

	transport, err := acp.NewStdioTransport(
		context.Background(),
		slog.Default(),
		"sh", scriptPath,
	)
	require.NoError(t, err)
	defer func() { _ = transport.Close() }()

	result, err := transport.Call(context.Background(), "test/method", map[string]string{"hello": "world"})
	require.NoError(t, err)
	require.NotNil(t, result)

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal(result, &parsed))
	require.Equal(t, true, parsed["ok"])
}

func TestStdioNewStdioTransportNotFound(t *testing.T) {
	_, err := acp.NewStdioTransport(
		context.Background(),
		slog.Default(),
		"/nonexistent/binary/that/does/not/exist",
	)
	require.Error(t, err)
}

func TestStdioStdioErrorString(t *testing.T) {
	e := &acp.StdioError{Code: -32603, Message: "internal error"}
	require.Contains(t, e.Error(), "-32603")
	require.Contains(t, e.Error(), "internal error")
}

func TestStdioCallCancelSendsCancelNotification(t *testing.T) {
	transport, mock := newMockTransport(t)
	defer func() { _ = transport.Close() }()

	var receivedCancel bool
	var mu sync.Mutex

	go func() {
		scanner := bufio.NewScanner(mock.stdin)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for scanner.Scan() {
			line := scanner.Bytes()
			var raw struct {
				Method string          `json:"method"`
				Params json.RawMessage `json:"params"`
			}
			if err := json.Unmarshal(line, &raw); err != nil {
				continue
			}
			if raw.Method == acp.NotificationCancelRequest {
				mu.Lock()
				receivedCancel = true
				mu.Unlock()
			}
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	_, _ = transport.Call(ctx, "test/method", nil)

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	require.True(t, receivedCancel, "transport should send $/cancel_request notification when Call context is cancelled")
}

func TestStdioEmptyLine(t *testing.T) {
	transport, mock := newMockTransport(t)
	defer func() { _ = transport.Close() }()

	handlerCalled := make(chan struct{}, 1)
	transport.SetNotificationHandler(func(method string, params json.RawMessage) {
		handlerCalled <- struct{}{}
	})

	_, err := fmt.Fprintln(mock.stdout, "")
	require.NoError(t, err)

	_, err = fmt.Fprintln(mock.stdout, `{"jsonrpc":"2.0","method":"test/valid","params":{}}`)
	require.NoError(t, err)

	select {
	case <-handlerCalled:
	case <-time.After(2 * time.Second):
		t.Fatal("readLoop should skip empty lines and continue processing")
	}
}

func TestStdioUnrecognizedMessage(t *testing.T) {
	transport, mock := newMockTransport(t)
	defer func() { _ = transport.Close() }()

	handlerCalled := make(chan struct{}, 1)
	transport.SetNotificationHandler(func(method string, params json.RawMessage) {
		handlerCalled <- struct{}{}
	})

	_, err := fmt.Fprintln(mock.stdout, `{"jsonrpc":"2.0"}`)
	require.NoError(t, err)

	_, err = fmt.Fprintln(mock.stdout, `{"jsonrpc":"2.0","method":"test/valid","params":{}}`)
	require.NoError(t, err)

	select {
	case <-handlerCalled:
	case <-time.After(2 * time.Second):
		t.Fatal("readLoop should skip unrecognized messages and continue processing")
	}
}

func TestStdioLargePayload(t *testing.T) {
	transport, mock := newMockTransport(t)
	mock.runEchoHandler(t)
	defer func() { _ = transport.Close() }()

	largeText := strings.Repeat("A", 100*1024)
	result, err := transport.Call(context.Background(), "test/method", map[string]string{"data": largeText})
	require.NoError(t, err)
	require.NotNil(t, result)
}
