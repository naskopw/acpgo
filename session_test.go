package acp_test

import (
	"encoding/json"
	"testing"

	"github.com/naskopw/acpgo"
	"github.com/stretchr/testify/require"
)

func TestSessionUpdateJSON(t *testing.T) {
	su := acp.SessionUpdate{
		SessionUpdateVariant: acp.SessionUpdateAgentMessageChunk,
		SessionID:            "ses-1",
		Content: &acp.ContentBlock{
			Type: "text",
			Text: "Hello!",
		},
	}
	data, err := json.Marshal(su)
	require.NoError(t, err)
	var got acp.SessionUpdate
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, acp.SessionUpdateAgentMessageChunk, got.SessionUpdateVariant)
	require.NotNil(t, got.Content)
	require.Equal(t, "Hello!", got.Content.Text)
}

func TestSessionUpdateEndTurn(t *testing.T) {
	su := acp.SessionUpdate{
		SessionUpdateVariant: acp.SessionUpdateEndTurn,
		SessionID:            "ses-1",
		StopReason:           acp.StopReasonEndTurn,
	}
	data, err := json.Marshal(su)
	require.NoError(t, err)
	var got acp.SessionUpdate
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, acp.SessionUpdateEndTurn, got.SessionUpdateVariant)
	require.Equal(t, acp.StopReasonEndTurn, got.StopReason)
}

func TestSessionModelChangedJSON(t *testing.T) {
	smc := acp.SessionModelChanged{
		SessionID: "ses-1",
		ModelID:   "gpt-4",
	}
	data, err := json.Marshal(smc)
	require.NoError(t, err)
	var got acp.SessionModelChanged
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, "ses-1", got.SessionID)
	require.Equal(t, "gpt-4", got.ModelID)
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
