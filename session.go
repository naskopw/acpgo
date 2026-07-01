package acp

// SessionUpdate is a discriminated union sent via session/update notifications.
// The SessionUpdateVariant field determines which variant it is.
type SessionUpdate struct {
	SessionUpdateVariant string            `json:"sessionUpdate"`
	SessionID            string            `json:"sessionId,omitempty"`
	MessageID            string            `json:"messageId,omitempty"`
	Content              *ContentBlock     `json:"content,omitempty"`
	ToolCall             *ToolCall         `json:"toolCall,omitempty"`
	ToolCallUpdate       *ToolCallUpdate   `json:"toolCallUpdate,omitempty"`
	Plan                 *Plan             `json:"plan,omitempty"`
	Commands             []SlashCommand    `json:"commands,omitempty"`
	CurrentModeID        string            `json:"currentModeId,omitempty"`
	ConfigOptions        []ConfigOption    `json:"configOptions,omitempty"`
	Title                string            `json:"title,omitempty"`
	UpdatedAt            string            `json:"updatedAt,omitempty"`
	Usage                *UsageUpdate      `json:"usage,omitempty"`
	Error                *StructuredError  `json:"error,omitempty"`
	StopReason           string            `json:"stopReason,omitempty"`
	Meta                 map[string]any    `json:"_meta,omitempty"`
}

// AvailableCommandsUpdate is sent when available commands change.
type AvailableCommandsUpdate struct {
	AvailableCommands []SlashCommand `json:"availableCommands"`
	Meta              map[string]any `json:"_meta,omitempty"`
}

// CurrentModeUpdate is sent when the session mode changes.
type CurrentModeUpdate struct {
	CurrentModeID string         `json:"currentModeId"`
	Meta          map[string]any `json:"_meta,omitempty"`
}

// ConfigOptionUpdate is sent when config options change.
type ConfigOptionUpdate struct {
	ConfigOptions []ConfigOption `json:"configOptions"`
	Meta          map[string]any `json:"_meta,omitempty"`
}

// SessionInfoUpdate is sent when session metadata (title) changes.
type SessionInfoUpdate struct {
	Title     string         `json:"title,omitempty"`
	UpdatedAt string         `json:"updatedAt,omitempty"`
	Meta      map[string]any `json:"_meta,omitempty"`
}

// UsageUpdate contains token usage information.
type UsageUpdate struct {
	TokensIn  int64          `json:"tokensIn,omitempty"`
	TokensOut int64          `json:"tokensOut,omitempty"`
	Cost      *Cost          `json:"cost,omitempty"`
	Meta      map[string]any `json:"_meta,omitempty"`
}

// Cost represents monetary cost of a turn.
type Cost struct {
	Amount   float64        `json:"amount"`
	Currency string         `json:"currency"`
	Meta     map[string]any `json:"_meta,omitempty"`
}

// Plan describes the agent's execution plan.
type Plan struct {
	Entries []PlanEntry    `json:"entries"`
	Meta    map[string]any `json:"_meta,omitempty"`
}

// PlanEntry is a single step in an execution plan.
type PlanEntry struct {
	Content  string         `json:"content"`
	Priority string         `json:"priority"`
	Status   string         `json:"status"`
	Meta     map[string]any `json:"_meta,omitempty"`
}

// SessionModelChanged is a notification that the model changed externally.
type SessionModelChanged struct {
	SessionID string `json:"sessionId"`
	ModelID   string `json:"modelId"`
}
