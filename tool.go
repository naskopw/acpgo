package acp

// ToolCall represents a tool call requested by the language model.
type ToolCall struct {
	ToolCallID string            `json:"toolCallId"`
	Title      string            `json:"title"`
	Kind       string            `json:"kind,omitempty"`
	Status     string            `json:"status,omitempty"`
	Content    []ToolCallContent `json:"content,omitempty"`
	Locations  []ToolCallLocation `json:"locations,omitempty"`
	RawInput   any               `json:"rawInput,omitempty"`
	RawOutput  any               `json:"rawOutput,omitempty"`
	Meta       map[string]any    `json:"_meta,omitempty"`
}

// ToolCallUpdate is an update to an existing tool call.
type ToolCallUpdate struct {
	ToolCallID string            `json:"toolCallId"`
	Kind       string            `json:"kind,omitempty"`
	Status     string            `json:"status,omitempty"`
	Title      string            `json:"title,omitempty"`
	Content    []ToolCallContent `json:"content,omitempty"`
	Locations  []ToolCallLocation `json:"locations,omitempty"`
	RawInput   any               `json:"rawInput,omitempty"`
	RawOutput  any               `json:"rawOutput,omitempty"`
	Meta       map[string]any    `json:"_meta,omitempty"`
}

// ToolCallContent is content produced by a tool call (content, diff, or terminal).
type ToolCallContent struct {
	Type        string         `json:"type"`
	Content     *ContentBlock  `json:"content,omitempty"`
	DiffPath    string         `json:"path,omitempty"`
	DiffOldText string         `json:"oldText,omitempty"`
	DiffNewText string         `json:"newText,omitempty"`
	TerminalID  string         `json:"terminalId,omitempty"`
	Meta        map[string]any `json:"_meta,omitempty"`
}

// ToolCallLocation is a file location accessed by a tool.
type ToolCallLocation struct {
	Path string         `json:"path"`
	Line int            `json:"line,omitempty"`
	Meta map[string]any `json:"_meta,omitempty"`
}
