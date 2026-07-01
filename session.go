package acp

import (
	"encoding/json"
	"fmt"
)

// SessionNotification is the params of the session/update notification.
// It wraps a SessionUpdate with the session ID.
type SessionNotification struct {
	SessionID string         `json:"sessionId"`
	Update    SessionUpdate  `json:"update"`
	Meta      map[string]any `json:"_meta,omitempty"`
}

// SessionUpdate is a discriminated union sent via session/update notifications.
// The SessionUpdateVariant field (JSON: "sessionUpdate") determines which variant.
// All variant fields are flat — the spec does NOT nest them under sub-keys.
type SessionUpdate struct {
	SessionUpdateVariant string             `json:"sessionUpdate"`
	MessageID            string             `json:"messageId,omitempty"`
	Content              json.RawMessage    `json:"content,omitempty"`
	ToolCallID           string             `json:"toolCallId,omitempty"`
	Title                string             `json:"title,omitempty"`
	Kind                 string             `json:"kind,omitempty"`
	Status               string             `json:"status,omitempty"`
	Locations            []ToolCallLocation `json:"locations,omitempty"`
	RawInput             any                `json:"rawInput,omitempty"`
	RawOutput            any                `json:"rawOutput,omitempty"`
	Entries              []PlanEntry        `json:"entries,omitempty"`
	AvailableCommands    []AvailableCommand `json:"availableCommands,omitempty"`
	CurrentModeID        string             `json:"currentModeId,omitempty"`
	ConfigOptions        []ConfigOption     `json:"configOptions,omitempty"`
	UpdatedAt            string             `json:"updatedAt,omitempty"`
	Used                 uint64             `json:"used,omitempty"`
	Size                 uint64             `json:"size,omitempty"`
	Cost                 *Cost              `json:"cost,omitempty"`
	Meta                 map[string]any     `json:"_meta,omitempty"`
}

// ContentBlock decodes the polymorphic content field as a single ContentBlock
// (used by user_message_chunk, agent_message_chunk, agent_thought_chunk).
func (su *SessionUpdate) ContentBlock() (*ContentBlock, error) {
	if len(su.Content) == 0 {
		return nil, nil
	}
	var cb ContentBlock
	if err := json.Unmarshal(su.Content, &cb); err != nil {
		return nil, fmt.Errorf("decode content block: %w", err)
	}
	return &cb, nil
}

// SetContentBlock encodes a ContentBlock into the content field.
func (su *SessionUpdate) SetContentBlock(cb *ContentBlock) {
	if cb == nil {
		su.Content = nil
		return
	}
	su.Content, _ = json.Marshal(cb)
}

// ToolCallContent decodes the polymorphic content field as a []ToolCallContent
// (used by tool_call, tool_call_update).
func (su *SessionUpdate) ToolCallContent() ([]ToolCallContent, error) {
	if len(su.Content) == 0 {
		return nil, nil
	}
	var tcc []ToolCallContent
	if err := json.Unmarshal(su.Content, &tcc); err != nil {
		return nil, fmt.Errorf("decode tool call content: %w", err)
	}
	return tcc, nil
}

// SetToolCallContent encodes a []ToolCallContent into the content field.
func (su *SessionUpdate) SetToolCallContent(tcc []ToolCallContent) {
	if tcc == nil {
		su.Content = nil
		return
	}
	su.Content, _ = json.Marshal(tcc)
}

// UsageUpdate contains context window and cumulative cost information.
type UsageUpdate struct {
	Used uint64         `json:"used"`
	Size uint64         `json:"size"`
	Cost *Cost          `json:"cost,omitempty"`
	Meta map[string]any `json:"_meta,omitempty"`
}

// Cost represents monetary cost of a turn.
type Cost struct {
	Amount   float64        `json:"amount"`
	Currency string         `json:"currency"`
	Meta     map[string]any `json:"_meta,omitempty"`
}

// PlanEntry is a single step in an execution plan.
type PlanEntry struct {
	Content  string         `json:"content"`
	Priority string         `json:"priority"`
	Status   string         `json:"status"`
	Meta     map[string]any `json:"_meta,omitempty"`
}
