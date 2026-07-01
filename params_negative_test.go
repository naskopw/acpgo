package acp_test

import (
	"encoding/json"
	"testing"

	"github.com/naskopw/acpgo"
	"github.com/stretchr/testify/require"
)

func TestInitializeRequestMissingProtocolVersion(t *testing.T) {
	var req acp.InitializeRequest
	err := json.Unmarshal([]byte(`{}`), &req)
	require.NoError(t, err)
	require.Equal(t, 0, req.ProtocolVersion)
}

func TestNewSessionRequestMissingCWD(t *testing.T) {
	var req acp.NewSessionRequest
	err := json.Unmarshal([]byte(`{"mcpServers":[]}`), &req)
	require.NoError(t, err)
	require.Empty(t, req.CWD)
	require.Empty(t, req.MCPServers)
}

func TestPromptResponseMissingStopReason(t *testing.T) {
	var resp acp.PromptResponse
	err := json.Unmarshal([]byte(`{}`), &resp)
	require.NoError(t, err)
	require.Empty(t, resp.StopReason)
}

func TestContentBlockTextMissingText(t *testing.T) {
	var cb acp.ContentBlock
	err := json.Unmarshal([]byte(`{"type":"text"}`), &cb)
	require.NoError(t, err)
	require.Empty(t, cb.Text)
	require.Equal(t, "text", cb.Type)
}

func TestContentBlockResourceLinkMissingURI(t *testing.T) {
	var cb acp.ContentBlock
	err := json.Unmarshal([]byte(`{"type":"resource_link","name":"foo"}`), &cb)
	require.NoError(t, err)
	require.Empty(t, cb.URI)
}

func TestToolCallMissingToolCallID(t *testing.T) {
	var su acp.SessionUpdate
	err := json.Unmarshal([]byte(`{"sessionUpdate":"tool_call","title":"test"}`), &su)
	require.NoError(t, err)
	require.Empty(t, su.ToolCallID)
}

func TestPlanEntryMissingFields(t *testing.T) {
	var entry acp.PlanEntry
	err := json.Unmarshal([]byte(`{}`), &entry)
	require.NoError(t, err)
	require.Empty(t, entry.Content)
	require.Empty(t, entry.Priority)
	require.Empty(t, entry.Status)
}

func TestSessionUpdateNullContent(t *testing.T) {
	var su acp.SessionUpdate
	err := json.Unmarshal([]byte(`{"sessionUpdate":"agent_message_chunk","content":null}`), &su)
	require.NoError(t, err)
	require.NotNil(t, su.Content)
	require.Equal(t, json.RawMessage("null"), su.Content)
	cb, err := su.ContentBlock()
	require.NoError(t, err)
	require.NotNil(t, cb)
	require.Empty(t, cb.Type)
}

func TestSessionUpdateUnknownVariant(t *testing.T) {
	var su acp.SessionUpdate
	err := json.Unmarshal([]byte(`{"sessionUpdate":"nonexistent_variant"}`), &su)
	require.NoError(t, err)
	require.Equal(t, "nonexistent_variant", su.SessionUpdateVariant)
}

func TestPermissionOutcomeMissingOptionID(t *testing.T) {
	var po acp.PermissionOutcome
	err := json.Unmarshal([]byte(`{"outcome":"selected"}`), &po)
	require.NoError(t, err)
	require.Empty(t, po.OptionID)
	require.Equal(t, "selected", po.Outcome)
}

func TestRequestPermissionResponseNullOutcome(t *testing.T) {
	var rpr acp.RequestPermissionResponse
	err := json.Unmarshal([]byte(`{"outcome":null}`), &rpr)
	require.NoError(t, err)
	require.Nil(t, rpr.Outcome)
}

func TestContentPromptRequestEmptyPrompt(t *testing.T) {
	var req acp.ContentPromptRequest
	err := json.Unmarshal([]byte(`{"sessionId":"s1","prompt":[]}`), &req)
	require.NoError(t, err)
	require.Empty(t, req.Prompt)
	require.Equal(t, "s1", req.SessionID)
}

func TestSessionNotificationMissingUpdate(t *testing.T) {
	var sn acp.SessionNotification
	err := json.Unmarshal([]byte(`{"sessionId":"s1"}`), &sn)
	require.NoError(t, err)
	require.Empty(t, sn.Update.SessionUpdateVariant)
}
