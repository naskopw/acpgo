# acp Library Production Readiness — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix wire-format bugs and fill spec gaps in the acp Go library.

**Architecture:** Single Go package (`github.com/naskopw/acp`). All fixes are in-place changes to protocol.go, tool.go, params.go, capabilities.go, README.md, plus tests.

**Tech Stack:** Go 1.22, testify, no external dependencies beyond standard library.

## Global Constraints

- All constant/method value changes MUST match the ACP v1 schema exactly (no aliasing old values)
- Exported types keep their names — only values, fields, and tags change
- All existing tests MUST pass after each task
- No new dependencies beyond `github.com/stretchr/testify`

---

### Task 1: Fix ToolCallStatus constants

**Files:**
- Modify: `protocol.go:82-87`
- Test: `types_test.go:10-102`

**Interfaces:**
- Consumes: no prior tasks
- Produces: `ToolStatusCompleted = "completed"`, `ToolStatusFailed = "failed"`

- [ ] **Step 1: Change constants**

In `protocol.go`, change line 84 from `ToolStatusCompleted = "tool_call_completed"` to `ToolStatusCompleted = "completed"`.
Change line 85 from `ToolStatusFailed = "tool_call_failed"` to `ToolStatusFailed = "failed"`.

- [ ] **Step 2: Run tests to verify**

```bash
go test -v -race -count=1 ./...
```
Expected: existing tests should still pass (the constants aren't tested for specific values in the existing test suite beyond non-empty checks).

- [ ] **Step 3: Commit**

```bash
git add protocol.go
git commit -m "fix: align ToolCallStatus constants with ACP v1 schema"
```

---

### Task 2: Fix error code constants

**Files:**
- Modify: `protocol.go:127-135`

**Interfaces:**
- Consumes: no prior tasks
- Produces: `ErrCodeAuthRequired = -32000`, `ErrCodeResourceNotFound = -32002`, `ErrCodeServer` removed

- [ ] **Step 1: Fix error code constants**

In `protocol.go`, make these changes:
- `ErrCodeAuthRequired = -32001` → `ErrCodeAuthRequired = -32000`
- Remove `ErrCodeServer = -32000`
- Add `ErrCodeResourceNotFound = -32002`

- [ ] **Step 2: Update the test that validates error code negativity**

In `envelope_test.go:95-101`, the test `TestErrorCodeConstants` checks `ErrCodeServer` which is being removed. Either remove that check or replace it with `ErrCodeResourceNotFound`.

Change:
```go
func TestErrorCodeConstants(t *testing.T) {
    require.True(t, acp.ErrCodeParse < 0)
    require.True(t, acp.ErrCodeInvalidRequest < 0)
    require.True(t, acp.ErrCodeMethodNotFound < 0)
    require.True(t, acp.ErrCodeInvalidParams < 0)
    require.True(t, acp.ErrCodeInternal < 0)
    require.True(t, acp.ErrCodeResourceNotFound < 0)
    require.True(t, acp.ErrCodeAuthRequired < 0)
}
```

- [ ] **Step 3: Run tests**

```bash
go test -v -race -count=1 ./...
```
Expected: all pass.

- [ ] **Step 4: Commit**

```bash
git add protocol.go envelope_test.go
git commit -m "fix: align error code constants with ACP v1 schema"
```

---

### Task 3: Fix ToolCallContent diff serialization

**Files:**
- Modify: `tool.go:29-36`

**Interfaces:**
- Consumes: no prior tasks
- Produces: `ToolCallContent` with inline `Diff` fields instead of nested `*Diff`

- [ ] **Step 1: Update ToolCallContent struct**

Change `tool.go` from:
```go
type ToolCallContent struct {
    Type       string         `json:"type"`
    Content    *ContentBlock  `json:"content,omitempty"`
    Diff       *Diff          `json:"diff,omitempty"`
    TerminalID string         `json:"terminalId,omitempty"`
    Meta       map[string]any `json:"_meta,omitempty"`
}
```

To:
```go
type ToolCallContent struct {
    Type       string         `json:"type"`
    Content    *ContentBlock  `json:"content,omitempty"`
    DiffPath   string         `json:"path,omitempty"`
    DiffOldText string        `json:"oldText,omitempty"`
    DiffNewText string        `json:"newText,omitempty"`
    TerminalID string         `json:"terminalId,omitempty"`
    Meta       map[string]any `json:"_meta,omitempty"`
}
```

Note: `DiffOldText` is `omitempty` because schema allows `null` (for new files). We use `omitempty` which means an empty string (new file) serializes as absent — but that matches how `Diff` was used before since Go zero-value string is empty.

- [ ] **Step 2: Run tests**

```bash
go test -v -race -count=1 ./...
```
Expected: all pass.

- [ ] **Step 3: Commit**

```bash
git add tool.go
git commit -m "fix: flatten diffs in ToolCallContent per ACP v1 schema"
```

---

### Task 4: Fix SessionModeState JSON tags

**Files:**
- Modify: `params.go:276-281`

**Interfaces:**
- Consumes: no prior tasks
- Produces: `SessionModeState` with corrected tags

- [ ] **Step 1: Fix JSON tags**

In `params.go`, change:
```go
type SessionModeState struct {
    Available []SessionMode `json:"available"`
    CurrentID string        `json:"currentId,omitempty"`
    Meta      map[string]any `json:"_meta,omitempty"`
}
```
To:
```go
type SessionModeState struct {
    AvailableModes []SessionMode `json:"availableModes"`
    CurrentModeID  string        `json:"currentModeId,omitempty"`
    Meta           map[string]any `json:"_meta,omitempty"`
}
```

- [ ] **Step 2: Run tests**

```bash
go test -v -race -count=1 ./...
```
Expected: all pass.

- [ ] **Step 3: Commit**

```bash
git add params.go
git commit -m "fix: align SessionModeState JSON tags with ACP v1 schema"
```

---

### Task 5: Add missing domain constants

**Files:**
- Modify: `protocol.go`

**Interfaces:**
- Consumes: no prior tasks
- Produces: `PermOptionAllowOnce`, `PermOptionAllowAlways`, `PermOptionRejectOnce`, `PermOptionRejectAlways`, `PlanPriorityHigh`, `PlanPriorityMedium`, `PlanPriorityLow`, `PlanStatusPending`, `PlanStatusInProgress`, `PlanStatusCompleted`, `RoleAssistant`, `RoleUser`

- [ ] **Step 1: Add PermissionOptionKind constants**

At the end of `protocol.go`, add:
```go
// PermissionOptionKind values for PermissionOption.kind.
const (
    PermOptionAllowOnce     = "allow_once"
    PermOptionAllowAlways   = "allow_always"
    PermOptionRejectOnce    = "reject_once"
    PermOptionRejectAlways  = "reject_always"
)
```

- [ ] **Step 2: Add PlanEntry priority and status constants**

Add:
```go
// PlanEntryPriority values.
const (
    PlanPriorityHigh   = "high"
    PlanPriorityMedium = "medium"
    PlanPriorityLow    = "low"
)

// PlanEntryStatus values.
const (
    PlanStatusPending    = "pending"
    PlanStatusInProgress = "in_progress"
    PlanStatusCompleted  = "completed"
)
```

- [ ] **Step 3: Add Role constants**

Add:
```go
// Role values for Annotations.audience.
const (
    RoleAssistant = "assistant"
    RoleUser      = "user"
)
```

- [ ] **Step 4: Write tests for new constants**

Add to `types_test.go`:
```go
func TestPermissionOptionKindConstants(t *testing.T) {
    vals := []string{acp.PermOptionAllowOnce, acp.PermOptionAllowAlways, acp.PermOptionRejectOnce, acp.PermOptionRejectAlways}
    for _, v := range vals {
        require.NotEqual(t, "", v)
    }
}

func TestPlanConstants(t *testing.T) {
    require.Equal(t, "high", acp.PlanPriorityHigh)
    require.Equal(t, "medium", acp.PlanPriorityMedium)
    require.Equal(t, "low", acp.PlanPriorityLow)
    require.Equal(t, "pending", acp.PlanStatusPending)
    require.Equal(t, "in_progress", acp.PlanStatusInProgress)
    require.Equal(t, "completed", acp.PlanStatusCompleted)
}

func TestRoleConstants(t *testing.T) {
    require.Equal(t, "assistant", acp.RoleAssistant)
    require.Equal(t, "user", acp.RoleUser)
}
```

- [ ] **Step 5: Run tests**

```bash
go test -v -race -count=1 ./...
```
Expected: all pass.

- [ ] **Step 6: Commit**

```bash
git add protocol.go types_test.go
git commit -m "feat: add missing domain constants for ACP v1 schema"
```

---

### Task 6: Add empty response types, SlashCommand input, extension types, AuthMethod fields

**Files:**
- Modify: `params.go`, `capabilities.go`
- Test: `params_test.go`, `capabilities_test.go` (new if needed)

**Interfaces:**
- Produces: empty response types, `AvailableCommandInput`, `UnstructuredCommandInput`, `ExtRequest`, `ExtResponse`, `ExtNotification`, updated `AuthMethod` and `SlashCommand`

- [ ] **Step 1: Add empty response types**

In `params.go`, find the `AuthenticateResponse` placeholder comment (line 28: `// AuthenticateResponse is the result of authentication.`) and replace it with:
```go
type AuthenticateResponse struct {
    Meta map[string]any `json:"_meta,omitempty"`
}
```

Near `LogoutRequest`, add:
```go
type LogoutResponse struct {
    Meta map[string]any `json:"_meta,omitempty"`
}
```

Near `DeleteSessionRequest`, add:
```go
type DeleteSessionResponse struct {
    Meta map[string]any `json:"_meta,omitempty"`
}
```

Near `CloseSessionRequest`, add:
```go
type CloseSessionResponse struct {
    Meta map[string]any `json:"_meta,omitempty"`
}
```

In `terminal.go`, add `ReleaseTerminalResponse` and `KillTerminalResponse` after their request types:
```go
type ReleaseTerminalResponse struct {
    Meta map[string]any `json:"_meta,omitempty"`
}

type KillTerminalResponse struct {
    Meta map[string]any `json:"_meta,omitempty"`
}
```

Also add missing `WriteTextFileResponse` (currently missing from params.go). Add near `WriteTextFileRequest`:
```go
type WriteTextFileResponse struct {
    Meta map[string]any `json:"_meta,omitempty"`
}
```

In `params.go`, find `SetConfigOptionResponse` and add after `SetModeRequest`:
```go
type SetSessionModeResponse struct {
    Meta map[string]any `json:"_meta,omitempty"`
}
```

- [ ] **Step 2: Add SlashCommand input types**

In `params.go`, replace the `SlashCommand` struct with:
```go
type SlashCommand struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description,omitempty"`
    Input       *AvailableCommandInput `json:"input,omitempty"`
    Meta        map[string]any         `json:"_meta,omitempty"`
}

type AvailableCommandInput struct {
    Unstructured *UnstructuredCommandInput `json:"unstructured,omitempty"`
    Meta         map[string]any            `json:"_meta,omitempty"`
}

type UnstructuredCommandInput struct {
    Hint string         `json:"hint,omitempty"`
    Meta map[string]any `json:"_meta,omitempty"`
}
```

- [ ] **Step 3: Add extension method types**

Add to `params.go`:
```go
type ExtRequest struct {
    Method string          `json:"method"`
    Params json.RawMessage `json:"params,omitempty"`
    Meta   map[string]any  `json:"_meta,omitempty"`
}

type ExtResponse struct {
    Result json.RawMessage `json:"result,omitempty"`
    Error  *RPCError       `json:"error,omitempty"`
    Meta   map[string]any  `json:"_meta,omitempty"`
}

type ExtNotification struct {
    Method string          `json:"method"`
    Params json.RawMessage `json:"params,omitempty"`
    Meta   map[string]any  `json:"_meta,omitempty"`
}
```

- [ ] **Step 4: Add AuthMethod fields**

In `capabilities.go`, update `AuthMethod`:
```go
type AuthMethod struct {
    Type        string         `json:"type"`
    ID          string         `json:"id"`
    Name        string         `json:"name,omitempty"`
    Description string         `json:"description,omitempty"`
    Meta        map[string]any `json:"_meta,omitempty"`
}
```

- [ ] **Step 5: Write JSON roundtrip tests**

Add to `params_test.go`:
```go
func TestWriteTextFileResponseJSON(t *testing.T) {
    p := acp.WriteTextFileResponse{}
    data, err := json.Marshal(p)
    require.NoError(t, err)
    var got acp.WriteTextFileResponse
    require.NoError(t, json.Unmarshal(data, &got))
}

func TestLogoutResponseJSON(t *testing.T) {
    p := acp.LogoutResponse{}
    data, err := json.Marshal(p)
    require.NoError(t, err)
    var got acp.LogoutResponse
    require.NoError(t, json.Unmarshal(data, &got))
}

func TestSlashCommandWithInputJSON(t *testing.T) {
    p := acp.SlashCommand{
        Name:        "think",
        Description: "Think about a problem",
        Input: &acp.AvailableCommandInput{
            Unstructured: &acp.UnstructuredCommandInput{Hint: "What should I think about?"},
        },
    }
    data, err := json.Marshal(p)
    require.NoError(t, err)
    var got acp.SlashCommand
    require.NoError(t, json.Unmarshal(data, &got))
    require.Equal(t, "think", got.Name)
    require.NotNil(t, got.Input)
    require.NotNil(t, got.Input.Unstructured)
    require.Equal(t, "What should I think about?", got.Input.Unstructured.Hint)
}

func TestExtRequestJSON(t *testing.T) {
    p := acp.ExtRequest{Method: "custom/do_thing", Params: json.RawMessage(`{"key":"val"}`)}
    data, err := json.Marshal(p)
    require.NoError(t, err)
    var got acp.ExtRequest
    require.NoError(t, json.Unmarshal(data, &got))
    require.Equal(t, "custom/do_thing", got.Method)
}
```

Add to `capabilities_test.go` (new file):
```go
package acp_test

import (
    "encoding/json"
    "testing"

    "github.com/naskopw/acp"
    "github.com/stretchr/testify/require"
)

func TestAuthMethodJSON(t *testing.T) {
    m := acp.AuthMethod{
        Type:        "oauth2",
        ID:          "github",
        Name:        "GitHub OAuth",
        Description: "Sign in with GitHub",
    }
    data, err := json.Marshal(m)
    require.NoError(t, err)
    var got acp.AuthMethod
    require.NoError(t, json.Unmarshal(data, &got))
    require.Equal(t, "oauth2", got.Type)
    require.Equal(t, "github", got.ID)
    require.Equal(t, "GitHub OAuth", got.Name)
}

func TestSetSessionModeResponseJSON(t *testing.T) {
    p := acp.SetSessionModeResponse{}
    data, err := json.Marshal(p)
    require.NoError(t, err)
    var got acp.SetSessionModeResponse
    require.NoError(t, json.Unmarshal(data, &got))
}
```

- [ ] **Step 6: Run tests**

```bash
go test -v -race -count=1 ./...
```
Expected: all pass.

- [ ] **Step 7: Commit**

```bash
git add params.go terminal.go capabilities.go params_test.go capabilities_test.go
git commit -m "feat: add missing types and fields per ACP v1 schema"
```

---

### Task 7: Fix README wording

**Files:**
- Modify: `README.md:18`

- [ ] **Step 1: Fix the misleading description**

In `README.md`, change:
```
building both **ACP clients** (editor/consumer side) and **ACP agents** (harness/server side):
```
To:
```
building **ACP clients** (editor/consumer side) and implementing **ACP servers** (the harness/agent side of the protocol). This library provides the shared protocol types — actual agent runtime logic lives in the harness layer.

This library is harness-agnostic by design. For the harness server implementation, see the harnessd package.
```

Also update the `doc.go` if it has similar language (line 4):
```
// It provides types and a stdio transport for building both ACP agents (harnesses)
// and ACP clients (consumers), including session management, tool definitions,
```
To:
```
// It provides shared protocol types and a stdio transport for implementing
// ACP clients and ACP servers, including session management, tool definitions,
```

- [ ] **Step 2: Run lint and test**

```bash
go test -v -race -count=1 ./...
golangci-lint run
```
Expected: clean.

- [ ] **Step 3: Commit**

```bash
git add README.md doc.go
git commit -m "docs: clarify library scope — protocol types, not agent runtime"
```
