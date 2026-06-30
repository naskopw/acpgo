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
