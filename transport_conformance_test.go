package acp_test

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"

	acp "github.com/naskopw/acpgo"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/require"
)

type captureSubprocess struct {
	stdin    *os.File
	stdout   *os.File
	mu       sync.Mutex
	captured [][]byte
	received chan struct{}
}

func newCaptureSubprocess(t *testing.T) (*acp.StdioTransport, *captureSubprocess) {
	t.Helper()

	stdinR, stdinW, err := os.Pipe()
	require.NoError(t, err)

	stdoutR, stdoutW, err := os.Pipe()
	require.NoError(t, err)

	transport := acp.NewStdioTransportWithIO(
		context.Background(),
		slog.Default(),
		stdinW,
		stdoutR,
	)

	sub := &captureSubprocess{
		stdin:    stdinR,
		stdout:   stdoutW,
		received: make(chan struct{}, 100),
	}

	go sub.readLoop(t)

	t.Cleanup(func() {
		_ = transport.Close()
		_ = stdinR.Close()
		_ = stdoutW.Close()
	})

	return transport, sub
}

func (s *captureSubprocess) readLoop(t *testing.T) {
	t.Helper()
	scanner := bufio.NewScanner(s.stdin)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := make([]byte, len(scanner.Bytes()))
		copy(line, scanner.Bytes())

		var raw struct {
			ID     json.RawMessage `json:"id"`
			Method string          `json:"method"`
		}
		if err := json.Unmarshal(line, &raw); err != nil {
			continue
		}

		s.mu.Lock()
		s.captured = append(s.captured, line)
		s.mu.Unlock()

		s.received <- struct{}{}

		var resp string
		if len(raw.ID) > 0 && string(raw.ID) != "null" {
			resp = fmt.Sprintf(`{"jsonrpc":"2.0","id":%s,"result":{}}`, string(raw.ID))
		} else if raw.Method != "" {
			resp = `{"jsonrpc":"2.0"}`
		} else {
			continue
		}

		_, _ = fmt.Fprintln(s.stdout, resp)
	}
}

func (s *captureSubprocess) waitForCapture(t *testing.T) []byte {
	t.Helper()
	select {
	case <-s.received:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for capture")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.captured) == 0 {
		return nil
	}
	return s.captured[len(s.captured)-1]
}

func compileRootSchema(t *testing.T) *jsonschema.Schema {
	t.Helper()
	c := jsonschema.NewCompiler()

	data, err := os.ReadFile("testdata/schema/v1/schema.json")
	require.NoError(t, err)

	var doc any
	require.NoError(t, json.Unmarshal(data, &doc))
	require.NoError(t, c.AddResource("schema.json", doc))

	sch, err := c.Compile("schema.json")
	require.NoError(t, err)
	return sch
}

func unmarshalAny(t *testing.T, data []byte) any {
	t.Helper()
	var v any
	require.NoError(t, json.Unmarshal(data, &v))
	return v
}

type transportTestCase struct {
	name      string
	method    string
	params    any
	side      string // "client" (client->agent) or "agent" (agent->client)
	msgType   string // "request" or "notification"
	defPath   string
}

func transportRoundTripCases() []transportTestCase {
	return []transportTestCase{
		// Agent → Client requests
		{name: "ReadTextFileRequest", method: acp.MethodReadTextFile, params: acp.ReadTextFileRequest{SessionID: "s1", Path: "/foo.txt"}, side: "agent", msgType: "request", defPath: "ReadTextFileRequest"},
		{name: "WriteTextFileRequest", method: acp.MethodWriteTextFile, params: acp.WriteTextFileRequest{SessionID: "s1", Path: "/foo.txt", Content: "data"}, side: "agent", msgType: "request", defPath: "WriteTextFileRequest"},
		{name: "RequestPermissionRequest", method: acp.MethodRequestPermission, params: acp.RequestPermissionRequest{SessionID: "s1", ToolCall: &acp.ToolCallUpdate{ToolCallID: "tc1"}, Options: []acp.PermissionOption{{OptionID: "opt1", Name: "Allow", Kind: "allow_once"}}}, side: "agent", msgType: "request", defPath: "RequestPermissionRequest"},
		{name: "CreateTerminalRequest", method: acp.MethodCreateTerminal, params: acp.CreateTerminalRequest{SessionID: "s1", Command: "/bin/bash"}, side: "agent", msgType: "request", defPath: "CreateTerminalRequest"},
		{name: "TerminalOutputRequest", method: acp.MethodTerminalOutput, params: acp.TerminalOutputRequest{SessionID: "s1", TerminalID: "t1"}, side: "agent", msgType: "request", defPath: "TerminalOutputRequest"},
		{name: "ReleaseTerminalRequest", method: acp.MethodReleaseTerminal, params: acp.ReleaseTerminalRequest{SessionID: "s1", TerminalID: "t1"}, side: "agent", msgType: "request", defPath: "ReleaseTerminalRequest"},
		{name: "WaitForTerminalExitRequest", method: acp.MethodWaitForTerminalExit, params: acp.WaitForTerminalExitRequest{SessionID: "s1", TerminalID: "t1"}, side: "agent", msgType: "request", defPath: "WaitForTerminalExitRequest"},
		{name: "KillTerminalRequest", method: acp.MethodKillTerminal, params: acp.KillTerminalRequest{SessionID: "s1", TerminalID: "t1"}, side: "agent", msgType: "request", defPath: "KillTerminalRequest"},

		// Client → Agent notifications
		{name: "CancelNotification", method: acp.NotificationCancel, params: acp.CancelNotification{SessionID: "s1"}, side: "client", msgType: "notification", defPath: "CancelNotification"},
		{name: "CancelRequestNotification", method: acp.NotificationCancelRequest, params: acp.CancelRequestNotification{RequestID: "req1"}, side: "client", msgType: "notification", defPath: "CancelRequestNotification"},
	}
}

func TestTransportConformance_RequestSchema(t *testing.T) {
	rootSch := compileRootSchema(t)

	for _, tc := range transportRoundTripCases() {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.msgType {
			case "request":
				testTransportRequest(t, rootSch, tc)
			case "notification":
				testTransportNotification(t, rootSch, tc)
			}
		})
	}
}

func testTransportRequest(t *testing.T, rootSch *jsonschema.Schema, tc transportTestCase) {
	t.Helper()
	transport, sub := newCaptureSubprocess(t)

	_, err := transport.Call(context.Background(), tc.method, tc.params)
	require.NoError(t, err)

	line := sub.waitForCapture(t)
	require.NotNil(t, line, "no request was sent")

	err = rootSch.Validate(unmarshalAny(t, line))
	require.NoError(t, err, "request JSON violates ACP schema\nJSON: %s", string(line))

	var parsed struct {
		ID     json.RawMessage `json:"id"`
		Method string          `json:"method"`
	}
	require.NoError(t, json.Unmarshal(line, &parsed))
	require.NotEmpty(t, parsed.ID, "request must have an id")
	require.Equal(t, tc.method, parsed.Method, "method must match")
}

func testTransportNotification(t *testing.T, rootSch *jsonschema.Schema, tc transportTestCase) {
	t.Helper()
	transport, sub := newCaptureSubprocess(t)

	err := transport.Notify(tc.method, tc.params)
	require.NoError(t, err)

	line := sub.waitForCapture(t)
	require.NotNil(t, line, "no notification was sent")

	err = rootSch.Validate(unmarshalAny(t, line))
	require.NoError(t, err, "notification JSON violates ACP schema\nJSON: %s", string(line))

	var parsed struct {
		ID     json.RawMessage `json:"id"`
		Method string          `json:"method"`
	}
	require.NoError(t, json.Unmarshal(line, &parsed))
	require.Empty(t, parsed.ID, "notification must not have an id")
	require.Equal(t, tc.method, parsed.Method, "method must match")
}

func TestTransportConformance_ResponseSchema(t *testing.T) {
	rootSch := compileRootSchema(t)

	transport, sub := newCaptureSubprocess(t)

	result, err := transport.Call(context.Background(), acp.MethodReadTextFile, acp.ReadTextFileRequest{SessionID: "s1", Path: "/foo.txt"})
	require.NoError(t, err)

	err = rootSch.Validate(unmarshalAny(t, []byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":1,"result":%s}`, string(result)))))
	require.NoError(t, err)

	line := sub.waitForCapture(t)
	require.NotNil(t, line, "no request was sent")
}

func TestTransportConformance_StdioSubprocess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping subprocess test in short mode")
	}

	rootSch := compileRootSchema(t)

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

	result, err := transport.Call(context.Background(), acp.MethodReadTextFile, acp.ReadTextFileRequest{SessionID: "s1", Path: "/foo.txt"})
	require.NoError(t, err)

	err = rootSch.Validate(unmarshalAny(t, []byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":1,"result":%s}`, string(result)))))
	require.NoError(t, err)
}
