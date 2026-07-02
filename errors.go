package acp

import "errors"

var (
	// ErrHarnessNotFound is returned when the harness binary cannot be found.
	ErrHarnessNotFound = errors.New("harness not found")
	// ErrProtocolMismatch is returned when the protocol versions are incompatible.
	ErrProtocolMismatch = errors.New("protocol version mismatch")
	// ErrRequestCancelled is returned when a request was cancelled.
	ErrRequestCancelled = errors.New("request cancelled")
	// ErrInvalidCWD is returned when a CWD path is set but is not absolute.
	ErrInvalidCWD = errors.New("cwd must be an absolute path")
	// ErrInvalidAdditionalDirectory is returned when an additionalDirectory entry is not absolute.
	ErrInvalidAdditionalDirectory = errors.New("additionalDirectories entries must be absolute paths")
	// ErrInvalidMCPCommand is returned when an MCP server command is not an absolute path.
	ErrInvalidMCPCommand = errors.New("mcpServers command must be an absolute path")

	// ErrInvalidContentType is returned when a ContentBlock.Type is not a valid value.
	ErrInvalidContentType        = errors.New("invalid content block type")
	// ErrInvalidToolKind is returned when a ToolCall.Kind is not a valid value.
	ErrInvalidToolKind           = errors.New("invalid tool kind")
	// ErrInvalidToolStatus is returned when a ToolCall.Status is not a valid value.
	ErrInvalidToolStatus         = errors.New("invalid tool status")
	// ErrInvalidToolContentType is returned when a ToolCallContent.Type is not a valid value.
	ErrInvalidToolContentType    = errors.New("invalid tool content type")
	// ErrInvalidSessionUpdateVariant is returned when a SessionUpdate.SessionUpdateVariant is not a valid value.
	ErrInvalidSessionUpdateVariant = errors.New("invalid session update variant")
	// ErrInvalidPlanPriority is returned when a PlanEntry.Priority is not a valid value.
	ErrInvalidPlanPriority       = errors.New("invalid plan entry priority")
	// ErrInvalidPlanStatus is returned when a PlanEntry.Status is not a valid value.
	ErrInvalidPlanStatus         = errors.New("invalid plan entry status")
	// ErrInvalidPermissionKind is returned when a PermissionOption.Kind is not a valid value.
	ErrInvalidPermissionKind     = errors.New("invalid permission option kind")
	// ErrInvalidPermissionOutcome is returned when a PermissionOutcome.Outcome is not a valid value.
	ErrInvalidPermissionOutcome  = errors.New("invalid permission outcome")
	// ErrInvalidMCPServerType is returned when an MCPServer.Type is not a valid value.
	ErrInvalidMCPServerType      = errors.New("invalid MCP server type")
	// ErrInvalidStopReason is returned when a PromptResponse.StopReason is not a valid value.
	ErrInvalidStopReason         = errors.New("invalid stop reason")
	// ErrInvalidConfigCategory is returned when a ConfigOption.Category is not a valid value.
	ErrInvalidConfigCategory     = errors.New("invalid config category")
	// ErrInvalidAudience is returned when an Annotations.Audience entry is not a valid value.
	ErrInvalidAudience           = errors.New("invalid audience value")
	// ErrInvalidRequiredField is returned when a required field is empty.
	ErrInvalidRequiredField      = errors.New("required field is empty")
)
