package acp_test

import (
	"encoding/json"
	"testing"

	"github.com/naskopw/acpgo"
	"github.com/stretchr/testify/require"
)

func mustMarshalJSON(t *testing.T, v any) json.RawMessage {
	t.Helper()
	data, err := json.Marshal(v)
	require.NoError(t, err)
	return data
}

func TestSessionNotificationJSON(t *testing.T) {
	sn := acp.SessionNotification{
		SessionID: "sess_abc123",
		Update: acp.SessionUpdate{
			SessionUpdateVariant: acp.SessionUpdateAgentMessageChunk,
			MessageID:            "msg_001",
			Content:              mustMarshalJSON(t, acp.ContentBlock{Type: "text", Text: "Hello!"}),
		},
	}
	data, err := json.Marshal(sn)
	require.NoError(t, err)

	var raw map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(data, &raw))
	require.Contains(t, raw, "sessionId")
	require.Contains(t, raw, "update")
	require.NotContains(t, raw, "sessionUpdate")

	var updateRaw map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(raw["update"], &updateRaw))
	require.Contains(t, updateRaw, "sessionUpdate")
	require.Contains(t, updateRaw, "content")
}

func TestSessionNotificationUnmarshal(t *testing.T) {
	input := `{
		"sessionId": "sess_abc123",
		"update": {
			"sessionUpdate": "agent_message_chunk",
			"messageId": "msg_001",
			"content": {"type": "text", "text": "Hello!"}
		}
	}`
	var sn acp.SessionNotification
	require.NoError(t, json.Unmarshal([]byte(input), &sn))
	require.Equal(t, "sess_abc123", sn.SessionID)
	require.Equal(t, acp.SessionUpdateAgentMessageChunk, sn.Update.SessionUpdateVariant)
	require.Equal(t, "msg_001", sn.Update.MessageID)

	cb, err := sn.Update.ContentBlock()
	require.NoError(t, err)
	require.NotNil(t, cb)
	require.Equal(t, "text", cb.Type)
	require.Equal(t, "Hello!", cb.Text)
}

func TestSessionUpdateJSON(t *testing.T) {
	su := acp.SessionUpdate{
		SessionUpdateVariant: acp.SessionUpdateAgentMessageChunk,
	}
	su.SetContentBlock(&acp.ContentBlock{Type: "text", Text: "Hello!"})
	data, err := json.Marshal(su)
	require.NoError(t, err)
	var got acp.SessionUpdate
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, acp.SessionUpdateAgentMessageChunk, got.SessionUpdateVariant)

	cb, err := got.ContentBlock()
	require.NoError(t, err)
	require.NotNil(t, cb)
	require.Equal(t, "Hello!", cb.Text)
}

func TestSessionUpdateToolCallFlat(t *testing.T) {
	input := `{
		"sessionUpdate": "tool_call",
		"toolCallId": "call_001",
		"title": "Reading configuration file",
		"kind": "read",
		"status": "pending"
	}`
	var su acp.SessionUpdate
	require.NoError(t, json.Unmarshal([]byte(input), &su))
	require.Equal(t, acp.SessionUpdateToolCall, su.SessionUpdateVariant)
	require.Equal(t, "call_001", su.ToolCallID)
	require.Equal(t, "Reading configuration file", su.Title)
	require.Equal(t, "read", su.Kind)
	require.Equal(t, "pending", su.Status)
}

func TestSessionUpdatePlanFlat(t *testing.T) {
	input := `{
		"sessionUpdate": "plan",
		"entries": [
			{"content": "Step 1", "priority": "high", "status": "pending"},
			{"content": "Step 2", "priority": "low", "status": "completed"}
		]
	}`
	var su acp.SessionUpdate
	require.NoError(t, json.Unmarshal([]byte(input), &su))
	require.Equal(t, acp.SessionUpdatePlan, su.SessionUpdateVariant)
	require.Len(t, su.Entries, 2)
	require.Equal(t, "Step 1", su.Entries[0].Content)
}

func TestSessionUpdateUsageFlat(t *testing.T) {
	input := `{
		"sessionUpdate": "usage_update",
		"used": 53000,
		"size": 200000
	}`
	var su acp.SessionUpdate
	require.NoError(t, json.Unmarshal([]byte(input), &su))
	require.Equal(t, acp.SessionUpdateUsage, su.SessionUpdateVariant)
	require.Equal(t, uint64(53000), su.Used)
	require.Equal(t, uint64(200000), su.Size)
}

func TestSessionUpdateAvailableCommandsFlat(t *testing.T) {
	input := `{
		"sessionUpdate": "available_commands_update",
		"availableCommands": [
			{"name": "create_plan", "description": "Create a plan"}
		]
	}`
	var su acp.SessionUpdate
	require.NoError(t, json.Unmarshal([]byte(input), &su))
	require.Equal(t, acp.SessionUpdateAvailableCommands, su.SessionUpdateVariant)
	require.Len(t, su.AvailableCommands, 1)
	require.Equal(t, "create_plan", su.AvailableCommands[0].Name)
}

func TestSessionUpdateToolCallContentArray(t *testing.T) {
	input := `{
		"sessionUpdate": "tool_call_update",
		"toolCallId": "call_001",
		"content": [
			{"type": "content", "content": {"type": "text", "text": "Done"}}
		]
	}`
	var su acp.SessionUpdate
	require.NoError(t, json.Unmarshal([]byte(input), &su))
	require.Equal(t, acp.SessionUpdateToolCallUpdate, su.SessionUpdateVariant)

	tcc, err := su.ToolCallContent()
	require.NoError(t, err)
	require.Len(t, tcc, 1)
	require.Equal(t, "content", tcc[0].Type)
}

func TestSessionUpdateContentBlockHelper(t *testing.T) {
	su := acp.SessionUpdate{
		SessionUpdateVariant: acp.SessionUpdateAgentMessageChunk,
		Content:              mustMarshalJSON(t, acp.ContentBlock{Type: "text", Text: "Hi"}),
	}
	cb, err := su.ContentBlock()
	require.NoError(t, err)
	require.Equal(t, "Hi", cb.Text)
}

func TestSessionUpdateToolCallContentHelper(t *testing.T) {
	su := acp.SessionUpdate{
		SessionUpdateVariant: acp.SessionUpdateToolCall,
		Content: mustMarshalJSON(t, []acp.ToolCallContent{
			{Type: "content", Content: &acp.ContentBlock{Type: "text", Text: "result"}},
		}),
	}
	tcc, err := su.ToolCallContent()
	require.NoError(t, err)
	require.Len(t, tcc, 1)
	require.Equal(t, "result", tcc[0].Content.Text)
}

func TestSessionUpdateSetContentBlock(t *testing.T) {
	su := acp.SessionUpdate{
		SessionUpdateVariant: acp.SessionUpdateAgentMessageChunk,
	}
	su.SetContentBlock(&acp.ContentBlock{Type: "text", Text: "Hello"})
	cb, err := su.ContentBlock()
	require.NoError(t, err)
	require.Equal(t, "Hello", cb.Text)
}

func TestSessionUpdateSetToolCallContent(t *testing.T) {
	su := acp.SessionUpdate{
		SessionUpdateVariant: acp.SessionUpdateToolCall,
	}
	su.SetToolCallContent([]acp.ToolCallContent{
		{Type: "diff", DiffPath: "/foo.go", DiffNewText: "package main"},
	})
	tcc, err := su.ToolCallContent()
	require.NoError(t, err)
	require.Len(t, tcc, 1)
	require.Equal(t, "/foo.go", tcc[0].DiffPath)
}

func TestPlanEntryJSON(t *testing.T) {
	entry := acp.PlanEntry{
		Content:  "Analyze the existing codebase structure",
		Priority: acp.PlanPriorityHigh,
		Status:   acp.PlanStatusPending,
	}
	data, err := json.Marshal(entry)
	require.NoError(t, err)

	expected := `{"content":"Analyze the existing codebase structure","priority":"high","status":"pending"}`
	require.JSONEq(t, expected, string(data))

	var got acp.PlanEntry
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, "Analyze the existing codebase structure", got.Content)
	require.Equal(t, acp.PlanPriorityHigh, got.Priority)
	require.Equal(t, acp.PlanStatusPending, got.Status)
}

func TestPlanEntryOmitsNonSpecFields(t *testing.T) {
	entry := acp.PlanEntry{Content: "test", Priority: "high", Status: "pending"}
	data, err := json.Marshal(entry)
	require.NoError(t, err)
	require.NotContains(t, string(data), `"title"`)
	require.NotContains(t, string(data), `"description"`)
	require.NotContains(t, string(data), `"id"`)
}

func TestUsageUpdateJSON(t *testing.T) {
	uu := acp.UsageUpdate{
		Used: 53000,
		Size: 200000,
		Cost: &acp.Cost{Amount: 0.045, Currency: "USD"},
	}
	data, err := json.Marshal(uu)
	require.NoError(t, err)

	expected := `{"used":53000,"size":200000,"cost":{"amount":0.045,"currency":"USD"}}`
	require.JSONEq(t, expected, string(data))

	var got acp.UsageUpdate
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, uint64(53000), got.Used)
	require.Equal(t, uint64(200000), got.Size)
	require.NotNil(t, got.Cost)
	require.Equal(t, 0.045, got.Cost.Amount)
}

func TestUsageUpdateOmitsNonSpecFields(t *testing.T) {
	uu := acp.UsageUpdate{Used: 100, Size: 200}
	data, err := json.Marshal(uu)
	require.NoError(t, err)
	require.NotContains(t, string(data), `"tokensIn"`)
	require.NotContains(t, string(data), `"tokensOut"`)
}

func TestConfigOptionJSON(t *testing.T) {
	co := acp.ConfigOption{
		ID:           "model",
		Name:         "Model",
		Type:         "select",
		CurrentValue: "gpt-4",
		Options: []acp.ConfigOptionValue{
			{Value: "gpt-4", Name: "GPT-4"},
			{Value: "claude-3", Name: "Claude 3"},
		},
	}
	data, err := json.Marshal(co)
	require.NoError(t, err)
	var got acp.ConfigOption
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, "model", got.ID)
	require.Equal(t, "gpt-4", got.CurrentValue)
	require.Len(t, got.Options, 2)
}

func TestContentBlockNilContent(t *testing.T) {
	su := acp.SessionUpdate{
		SessionUpdateVariant: acp.SessionUpdateAgentMessageChunk,
	}
	cb, err := su.ContentBlock()
	require.NoError(t, err)
	require.Nil(t, cb)
}

func TestContentBlockEmptyContent(t *testing.T) {
	su := acp.SessionUpdate{
		SessionUpdateVariant: acp.SessionUpdateAgentMessageChunk,
		Content:              json.RawMessage{},
	}
	cb, err := su.ContentBlock()
	require.NoError(t, err)
	require.Nil(t, cb)
}

func TestContentBlockInvalidJSON(t *testing.T) {
	su := acp.SessionUpdate{
		SessionUpdateVariant: acp.SessionUpdateAgentMessageChunk,
		Content:              json.RawMessage(`{invalid`),
	}
	_, err := su.ContentBlock()
	require.Error(t, err)
}

func TestToolCallContentNilContent(t *testing.T) {
	su := acp.SessionUpdate{
		SessionUpdateVariant: acp.SessionUpdateToolCall,
	}
	tcc, err := su.ToolCallContent()
	require.NoError(t, err)
	require.Nil(t, tcc)
}

func TestToolCallContentEmptyContent(t *testing.T) {
	su := acp.SessionUpdate{
		SessionUpdateVariant: acp.SessionUpdateToolCall,
		Content:              json.RawMessage{},
	}
	tcc, err := su.ToolCallContent()
	require.NoError(t, err)
	require.Nil(t, tcc)
}

func TestToolCallContentInvalidJSON(t *testing.T) {
	su := acp.SessionUpdate{
		SessionUpdateVariant: acp.SessionUpdateToolCall,
		Content:              json.RawMessage(`[invalid`),
	}
	_, err := su.ToolCallContent()
	require.Error(t, err)
}

func TestSetContentBlockNil(t *testing.T) {
	su := acp.SessionUpdate{
		SessionUpdateVariant: acp.SessionUpdateAgentMessageChunk,
		Content:              mustMarshalJSON(t, acp.ContentBlock{Type: "text", Text: "old"}),
	}
	su.SetContentBlock(nil)
	require.Nil(t, su.Content)
}

func TestSetToolCallContentNil(t *testing.T) {
	su := acp.SessionUpdate{
		SessionUpdateVariant: acp.SessionUpdateToolCall,
		Content:              mustMarshalJSON(t, []acp.ToolCallContent{{Type: "content"}}),
	}
	su.SetToolCallContent(nil)
	require.Nil(t, su.Content)
}
