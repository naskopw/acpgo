package acp

// Standard ACP protocol method names (client → agent).
const (
	MethodInitialize         = "initialize"
	MethodAuthenticate       = "authenticate"
	MethodLogout             = "logout"
	MethodNewSession         = "session/new"
	MethodLoadSession        = "session/load"
	MethodListSessions       = "session/list"
	MethodDeleteSession      = "session/delete"
	MethodResumeSession      = "session/resume"
	MethodCloseSession       = "session/close"
	MethodPrompt             = "session/prompt"
	MethodSetMode            = "session/set_mode"
	MethodSetConfigOption    = "session/set_config_option"
	MethodCancel             = "session/cancel"
)

// ACP protocol method names (agent → client).
const (
	MethodReadTextFile         = "fs/read_text_file"
	MethodWriteTextFile        = "fs/write_text_file"
	MethodRequestPermission    = "session/request_permission"
	MethodCreateTerminal       = "terminal/create"
	MethodTerminalOutput       = "terminal/output"
	MethodReleaseTerminal      = "terminal/release"
	MethodWaitForTerminalExit  = "terminal/wait_for_exit"
	MethodKillTerminal         = "terminal/kill"
)

// Standard ACP protocol notification method names.
const (
	NotificationSessionUpdate  = "session/update"
	NotificationCancel         = "session/cancel"
	NotificationCancelRequest  = "$/cancel_request"
)

// SessionUpdate variants (discriminated by sessionUpdate field).
const (
	SessionUpdateUserMessageChunk     = "user_message_chunk"
	SessionUpdateAgentMessageChunk    = "agent_message_chunk"
	SessionUpdateAgentThoughtChunk    = "agent_thought_chunk"
	SessionUpdateToolCall             = "tool_call"
	SessionUpdateToolCallUpdate       = "tool_call_update"
	SessionUpdatePlan                 = "plan"
	SessionUpdateAvailableCommands    = "available_commands_update"
	SessionUpdateCurrentMode          = "current_mode_update"
	SessionUpdateConfigOption         = "config_option_update"
	SessionUpdateSessionInfo          = "session_info_update"
	SessionUpdateUsage                = "usage_update"
	SessionUpdateEndTurn              = "end_turn"
	SessionUpdateError                = "error"
)

// ContentBlock type discriminators.
const (
	ContentTypeText         = "text"
	ContentTypeImage        = "image"
	ContentTypeAudio        = "audio"
	ContentTypeResourceLink = "resource_link"
	ContentTypeResource     = "resource"
)

// ToolCallContent type discriminators.
const (
	ToolContentTypeContent  = "content"
	ToolContentTypeDiff     = "diff"
	ToolContentTypeTerminal = "terminal"
)

// StopReason values.
const (
	StopReasonEndTurn       = "end_turn"
	StopReasonMaxTokens     = "max_tokens"
	StopReasonMaxTurnReqs   = "max_turn_requests"
	StopReasonRefusal       = "refusal"
	StopReasonCancelled     = "cancelled"
)

// ToolCallStatus values.
const (
	ToolStatusPending     = "pending"
	ToolStatusInProgress  = "in_progress"
	ToolStatusCompleted   = "completed"
	ToolStatusFailed      = "failed"
)

// ToolKind values.
const (
	ToolKindRead       = "read"
	ToolKindEdit       = "edit"
	ToolKindDelete     = "delete"
	ToolKindMove       = "move"
	ToolKindSearch     = "search"
	ToolKindExecute    = "execute"
	ToolKindThink      = "think"
	ToolKindFetch      = "fetch"
	ToolKindSwitchMode = "switch_mode"
	ToolKindOther      = "other"
)

// RequestPermissionOutcome values.
const (
	PermOutcomeCancelled = "cancelled"
	PermOutcomeSelected  = "selected"
)

// SelectedPermissionOutcomeOption values.
const (
	PermSelectedOnce        = "once"
	PermSelectedAlways      = "always"
	PermSelectedReject      = "reject"
	PermSelectedRejectAlways = "reject_always"
)

// ConfigOptionCategory values.
const (
	ConfigCategoryMode         = "mode"
	ConfigCategoryModel        = "model"
	ConfigCategoryModelConfig  = "model_config"
	ConfigCategoryThoughtLevel = "thought_level"
)

// Error codes (ACP standard).
const (
	ErrCodeParse           = -32700
	ErrCodeInvalidRequest  = -32600
	ErrCodeMethodNotFound  = -32601
	ErrCodeInvalidParams   = -32602
	ErrCodeInternal        = -32603
	ErrCodeCancelled       = -32800
	ErrCodeAuthRequired       = -32000
	ErrCodeResourceNotFound   = -32002
)
