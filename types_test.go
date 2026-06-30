package acp_test

import (
	"encoding/json"
	"testing"

	"github.com/naskopw/acpgo"
	"github.com/stretchr/testify/require"
)

func TestModelInfoJSON(t *testing.T) {
	m := acp.ModelInfo{ID: "gpt-4", Name: "GPT-4"}
	data, err := json.Marshal(m)
	require.NoError(t, err)
	var got acp.ModelInfo
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, "gpt-4", got.ID)
	require.Equal(t, "GPT-4", got.Name)
}

func TestPromptRequestJSON(t *testing.T) {
	modelID := "gpt-4"
	req := acp.PromptRequest{
		RequestID: "req-1",
		SessionID: "ses-1",
		Prompt:    "hello",
		ModelID:   &modelID,
	}
	data, err := json.Marshal(req)
	require.NoError(t, err)
	var got acp.PromptRequest
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, "req-1", got.RequestID)
	require.Equal(t, "hello", got.Prompt)
	require.NotNil(t, got.ModelID)
	require.Equal(t, "gpt-4", *got.ModelID)
}

func TestInitializeResultJSON(t *testing.T) {
	init := acp.InitializeResult{
		ProtocolVersion: 1,
		AgentCapabilities: &acp.AgentCapabilities{
			LoadSession: true,
			SessionCapabilities: &acp.SessionCapabilities{
				Delete: struct{}{},
				List:   struct{}{},
			},
		},
		AgentName:    "test-agent",
		AgentVersion: "1.0.0",
	}
	data, err := json.Marshal(init)
	require.NoError(t, err)
	var got acp.InitializeResult
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, 1, got.ProtocolVersion)
	require.NotNil(t, got.AgentCapabilities)
	require.True(t, got.AgentCapabilities.LoadSession)
	require.NotNil(t, got.AgentCapabilities.SessionCapabilities)
	require.NotNil(t, got.AgentCapabilities.SessionCapabilities.Delete)
	require.Equal(t, "test-agent", got.AgentName)
}

func TestChunkTypeValues(t *testing.T) {
	require.True(t, acp.ChunkText == 0 && acp.ChunkThought == 1 && acp.ChunkToolUse == 2 && acp.ChunkDone == 3 && acp.ChunkError == 4, "chunk type values changed")
}

func TestMethodConstants(t *testing.T) {
	expected := map[string]string{
		"Initialize":        acp.MethodInitialize,
		"NewSession":        acp.MethodNewSession,
		"ListSessions":      acp.MethodListSessions,
		"ResumeSession":     acp.MethodResumeSession,
		"DeleteSession":     acp.MethodDeleteSession,
		"Prompt":            acp.MethodPrompt,
		"LoadSession":       acp.MethodLoadSession,
		"CloseSession":      acp.MethodCloseSession,
		"Authenticate":      acp.MethodAuthenticate,
		"Logout":            acp.MethodLogout,
		"SetConfigOption":   acp.MethodSetConfigOption,
		"SetMode":           acp.MethodSetMode,
		"ReadTextFile":      acp.MethodReadTextFile,
		"WriteTextFile":     acp.MethodWriteTextFile,
		"RequestPermission": acp.MethodRequestPermission,
		"CreateTerminal":    acp.MethodCreateTerminal,
	}
	for name, val := range expected {
		require.NotEqual(t, "", val, "constant %s is empty", name)
	}
}

func TestNotificationConstants(t *testing.T) {
	require.NotEqual(t, "", acp.NotificationSessionUpdate, "NotificationSessionUpdate is empty")
	require.NotEqual(t, "", acp.NotificationCancel, "NotificationCancel is empty")
	require.NotEqual(t, "", acp.NotificationCancelRequest, "NotificationCancelRequest is empty")
}

func TestErrorSentinels(t *testing.T) {
	require.NotNil(t, acp.ErrHarnessNotFound, "ErrHarnessNotFound is nil")
	require.NotNil(t, acp.ErrProtocolMismatch, "ErrProtocolMismatch is nil")
	require.NotNil(t, acp.ErrRequestCancelled, "ErrRequestCancelled is nil")
}

func TestPermissionOptionKindConstants(t *testing.T) {
	vals := []string{acp.PermOptionAllowOnce, acp.PermOptionAllowAlways, acp.PermOptionRejectOnce, acp.PermOptionRejectAlways}
	for _, v := range vals {
		require.NotEqual(t, "", v)
	}
}

func TestPlanConstants(t *testing.T) {
	require.Equal(t, "high", acp.PlanPriorityHigh)
	require.Equal(t, "medium", acp.PlanPriorityMedium)
	require.Equal(t, "low", acp.PlanPriorityLow)
	require.Equal(t, "pending", acp.PlanStatusPending)
	require.Equal(t, "in_progress", acp.PlanStatusInProgress)
	require.Equal(t, "completed", acp.PlanStatusCompleted)
}

func TestRoleConstants(t *testing.T) {
	require.Equal(t, "assistant", acp.RoleAssistant)
	require.Equal(t, "user", acp.RoleUser)
}

func TestToolCallContentDiffJSON(t *testing.T) {
	c := acp.ToolCallContent{
		Type:        "diff",
		DiffPath:    "/foo/bar.go",
		DiffNewText: "package main",
	}
	data, err := json.Marshal(c)
	require.NoError(t, err)
	var got acp.ToolCallContent
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, "diff", got.Type)
	require.Equal(t, "/foo/bar.go", got.DiffPath)
	require.Equal(t, "package main", got.DiffNewText)
	require.Contains(t, string(data), `"path"`)
	require.NotContains(t, string(data), `"diff":`)
}
