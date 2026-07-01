# acpgo — Agent Client Protocol for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/naskopw/acpgo.svg)](https://pkg.go.dev/github.com/naskopw/acpgo)
[![Go Version](https://img.shields.io/github/go-mod/go-version/naskopw/acpgo)](go.mod)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue)](LICENSE)

**acpgo** is a Go implementation of the [Agent Client Protocol (ACP)](https://agentclientprotocol.com/) — an open standard for communication between code editors and AI coding agents.

## What is ACP?

ACP standardizes how editors and agents exchange messages over JSON-RPC 2.0. It covers session lifecycle management, prompt/response turns with streaming, tool calling, file operations, terminal management, capability negotiation, and authentication.

- [Official Specification](https://github.com/agentclientprotocol/agent-client-protocol)
- [ACP Website](https://agentclientprotocol.com/)

## What This Library Provides

A complete set of Go types and utilities for building **ACP clients** (editor/consumer side) and implementing **ACP servers** (the harness/agent side of the protocol). This library provides the shared protocol types — actual agent runtime logic lives in the harness layer.

This library is harness-agnostic by design. For the harness server implementation, see the harnessd package.

| Area | Description |
|---|---|
| **JSON-RPC Envelopes** | `Request`, `Response`, `Notification`, `RPCError` |
| **Protocol Methods** | All client→agent and agent→client methods as typed constants |
| **Session Management** | Create, load, list, resume, close, delete sessions |
| **Streaming** | `SessionUpdate` discriminated union with message chunks, tool calls, plans, usage, mode/config updates |
| **Tool Calling** | `ToolCall`, `ToolCallUpdate`, tool content types, locations |
| **Content Types** | Text, image, audio, resource links, embedded resources, diffs |
| **File Operations** | Read/write text file requests and responses |
| **Terminal Management** | Create, output, wait, kill, release terminal |
| **Permission Requests** | Permission option types, outcomes, and responses |
| **Capability Negotiation** | Client and agent capability types for `initialize` handshake |
| **Configuration** | Config option types for mode, model, and other settings |
| **Authentication** | Auth method types and logout support |
| **Stdio Transport** | JSON-RPC over subprocess stdin/stdout with concurrent request tracking and context-based cancellation |

## Installation

```bash
go get github.com/naskopw/acpgo
```

## Usage

### Client-Side (Connecting to an Agent)

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "log/slog"

    "github.com/naskopw/acpgo"
)

func main() {
    ctx := context.Background()
    logger := slog.Default()

    // Start an agent subprocess over stdio
    transport, err := acp.NewStdioTransport(ctx, logger, "/path/to/agent", "--flag")
    if err != nil {
        log.Fatal(err)
    }
    defer func() { _ = transport.Close() }()

    // Register handler for streaming updates
    transport.SetNotificationHandler(func(method string, params json.RawMessage) {
        if method == acp.NotificationSessionUpdate {
            var sn acp.SessionNotification
            json.Unmarshal(params, &sn)
            // handle update based on sn.Update.SessionUpdateVariant
            switch sn.Update.SessionUpdateVariant {
            case acp.SessionUpdateAgentMessageChunk:
                cb, _ := sn.Update.ContentBlock()
                if cb != nil {
                    log.Printf("agent: %s", cb.Text)
                }
            }
        }
    })

    // Initialize handshake
    initResp, err := transport.Call(ctx, acp.MethodInitialize, acp.InitializeRequest{
        ProtocolVersion: 1,
        ClientInfo:      &acp.Implementation{Name: "my-client", Version: "1.0.0"},
    })
    if err != nil {
        log.Fatal(err)
    }
    _ = initResp

    // Create a session
    sessionResp, err := transport.Call(ctx, acp.MethodNewSession, acp.NewSessionRequest{
        CWD:        "/path/to/project",
        MCPServers: []acp.MCPServer{},
    })
    if err != nil {
        log.Fatal(err)
    }

    var session acp.NewSessionResponse
    json.Unmarshal(sessionResp, &session)

    // Send a prompt
    _, err = transport.Call(ctx, acp.MethodPrompt, acp.ContentPromptRequest{
        SessionID: session.SessionID,
        Prompt:    []acp.ContentBlock{{Type: "text", Text: "Add error handling to the HTTP handler"}},
    })
    if err != nil {
        log.Fatal(err)
    }
}
```

### Sending Notifications

```go
transport.Notify(acp.NotificationCancel, acp.CancelNotification{
    SessionID: sessionID,
})
```

## Development

```bash
make test    # run tests with race detector
make lint    # run golangci-lint
make build   # build the binary
```

## License

Apache 2.0
