# acp Go Library — Production Readiness

## Context

The `acp` Go library provides shared protocol types and the stdio transport for the
Agent Client Protocol. This spec covers making it production-ready: fixing wire-format
bugs, filling spec gaps, and correcting the README.

## Scope

This is **not** adding high-level agent or client runtime APIs. Those belong in
`../harnessd`. This library stays focused on low-level protocol types, JSON-RPC
envelopes, and the standard stdio transport.

## Tier 1 — Wire Compatibility Bugs (MUST FIX)

These break communication with spec-compliant peers.

### 1.1 ToolCallStatus values

| Constant | Current (wrong) | Schema value (correct) |
|---|---|---|
| `ToolStatusCompleted` | `"tool_call_completed"` | `"completed"` |
| `ToolStatusFailed` | `"tool_call_failed"` | `"failed"` |

**Fix:** Change the constant values. This is a breaking change for any existing
serialized data but matches the spec.

**Files:** `protocol.go`

**Migration:** Old values are not aliased. If backward compatibility with saved
data is needed, a migration utility or unmarshal shim can be added in a follow-up.

### 1.2 ToolCallContent diff serialization

Current Go struct nests `Diff` fields under a `"diff"` key:

```json
{"type": "diff", "diff": {"path": "...", "newText": "..."}}
```

Schema requires flattened:

```json
{"type": "diff", "path": "...", "newText": "..."}
```

**Fix:** Replace the single `Diff *Diff` field with inline `Diff` fields in
`ToolCallContent`. The `content` and `terminal` variants already serialize
correctly and are unchanged.

**Files:** `tool.go`, plus update any code constructing diff-type ToolCallContent.

### 1.3 SessionModeState JSON tags

| Field | Current tag | Correct tag |
|---|---|---|
| `Available` | `json:"available"` | `json:"availableModes"` |
| `CurrentID` | `json:"currentId"` | `json:"currentModeId"` |

**Fix:** Update the JSON struct tags.

**Files:** `params.go`

### 1.4 Error code AuthRequired

| Constant | Current | Schema |
|---|---|---|
| `ErrCodeAuthRequired` | -32001 | -32000 |

The schema defines error codes:
- -32700 Parse error
- -32600 Invalid Request
- -32601 Method not found
- -32602 Invalid params
- -32603 Internal error
- -32000 Auth required
- -32002 Resource not found
- -32800 Cancelled

No `-32001` in the schema. Go's `ErrCodeServer = -32000` also conflicts with
`ErrCodeAuthRequired = -32000` after correction — `ErrCodeServer` should
become `-32002` (renamed to `ErrCodeResourceNotFound`).

**Fix:**

- `ErrCodeAuthRequired`: -32001 → -32000
- `ErrCodeServer`: remove or rename to `ErrCodeResourceNotFound = -32002`

**Files:** `protocol.go`

### 1.5 Missing error code

Add `ErrCodeResourceNotFound = -32002`.

## Tier 2 — Spec Completeness (SHOULD ADD)

### 2.1 PermissionOptionKind constants

Add constants matching schema values for `PermissionOption.kind`:

| Constant | Value |
|---|---|
| `PermOptionAllowOnce` | `"allow_once"` |
| `PermOptionAllowAlways` | `"allow_always"` |
| `PermOptionRejectOnce` | `"reject_once"` |
| `PermOptionRejectAlways` | `"reject_always"` |

The existing `PermSelectedOnce = "once"` etc. constants are for the
`SelectedPermission.reply` field and are kept as-is.

**Files:** `protocol.go`

### 2.2 Empty response types

The following methods return empty response bodies (`{_meta: ...}`) but have no
Go types. Add empty structs for generic handler use:

- `WriteTextFileResponse`
- `ReleaseTerminalResponse`
- `KillTerminalResponse`
- `AuthenticateResponse`
- `LogoutResponse`
- `DeleteSessionResponse`
- `CloseSessionResponse`
- `SetSessionModeResponse` (corresponds to `SetModeRequest`)

**Files:** `params.go`

### 2.3 Domain constants

Add constants for schema enums used in session updates and plans:

| Constant group | Values |
|---|---|
| `PlanEntryPriority` | `PlanPriorityHigh = "high"`, `PlanPriorityMedium = "medium"`, `PlanPriorityLow = "low"` |
| `PlanEntryStatus` | `PlanStatusPending = "pending"`, `PlanStatusInProgress = "in_progress"`, `PlanStatusCompleted = "completed"` |
| `Role` | `RoleAssistant = "assistant"`, `RoleUser = "user"` |

**Files:** `protocol.go`

### 2.4 AuthMethod fields

Add `Name` and `Description` fields to `AuthMethod`:

```go
type AuthMethod struct {
    Type        string         `json:"type"`
    ID          string         `json:"id"`
    Name        string         `json:"name,omitempty"`
    Description string         `json:"description,omitempty"`
    Meta        map[string]any `json:"_meta,omitempty"`
}
```

**Files:** `capabilities.go`

### 2.5 SlashCommand / AvailableCommand

Add `Input` field to `SlashCommand` to match schema's `AvailableCommand.input`:

```go
type SlashCommand struct {
    Name        string                  `json:"name"`
    Description string                  `json:"description,omitempty"`
    Input       *AvailableCommandInput  `json:"input,omitempty"`
    Meta        map[string]any          `json:"_meta,omitempty"`
}

type AvailableCommandInput struct {
    Unstructured *UnstructuredCommandInput `json:"unstructured,omitempty"`
    Meta         map[string]any             `json:"_meta,omitempty"`
}

type UnstructuredCommandInput struct {
    Hint string         `json:"hint,omitempty"`
    Meta map[string]any `json:"_meta,omitempty"`
}
```

**Files:** `params.go`

### 2.6 Extension method types

Add opaque types for extension methods (non-standard methods outside the ACP spec):

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

**Files:** `params.go`

## Tier 3 — Intentional Skips

These schema constructs exist but the current flat Go structs work correctly at
the JSON level and provide a better Go API:

- `ContentBlock` discriminated union → flat struct (simpler for consumers)
- `SessionNotification` wrapper → `SessionUpdate` flat (fewer allocations)
- MCPServer discriminated union → flat struct (simpler construction)
- PlanEntry fields (`ID`, `Title`, `Description`) → extra fields beyond schema (more useful)
- `SessionConfigOption` discriminated union → flat `ConfigOption` (simpler)
- `EmbeddedResource` nesting → flat fields (simpler)

## README

Replace "building both **ACP clients** (editor/consumer side) and **ACP agents** (harness/server side)" with:

> "building both **ACP clients** (editor/consumer side) and **ACP servers** (the harness/agent side of the protocol). This library provides the shared protocol types — actual agent runtime logic is implemented by the harness layer."

## Testing

No behavioral changes to existing logic — the fixes are:
- Constant value changes (tests verify values)
- JSON tag changes (tests verify roundtrip)
- New types added (test JSON roundtrip)

All existing tests continue to pass with updated values.

## Non-Goals

- No runtime API for agent or client implementation
- No harness logic
- No breaking of existing working code beyond the intentional constant fixes
