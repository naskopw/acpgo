package acp_test

import (
	"encoding/json"
	"testing"

	"github.com/naskopw/acpgo"
	"github.com/stretchr/testify/require"
)

func TestInitializeRequestJSON(t *testing.T) {
	p := acp.InitializeRequest{
		ProtocolVersion: 1,
		ClientInfo: &acp.Implementation{
			Name:    "test-client",
			Version: "1.0",
		},
	}
	data, err := json.Marshal(p)
	require.NoError(t, err)
	var got acp.InitializeRequest
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, 1, got.ProtocolVersion)
	require.NotNil(t, got.ClientInfo)
	require.Equal(t, "test-client", got.ClientInfo.Name)
}

func TestInitializeResponseJSON(t *testing.T) {
	p := acp.InitializeResponse{
		ProtocolVersion: 1,
		AgentInfo: &acp.Implementation{
			Name:    "agent",
			Version: "1.0",
		},
		AgentCapabilities: &acp.AgentCapabilities{
			LoadSession: true,
		},
	}
	data, err := json.Marshal(p)
	require.NoError(t, err)
	var got acp.InitializeResponse
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, 1, got.ProtocolVersion)
	require.True(t, got.AgentCapabilities.LoadSession)
}

func TestSetModelParamsJSON(t *testing.T) {
	p := acp.SetModelParams{ModelID: "gpt-4"}
	data, err := json.Marshal(p)
	require.NoError(t, err)
	var got acp.SetModelParams
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, "gpt-4", got.ModelID)
}

func TestNewSessionParamsJSON(t *testing.T) {
	p := acp.NewSessionParams{
		CWD: "/home",
		MCPServers: []acp.MCPServer{
			{Name: "fs", Command: "npx", Args: []string{"-y", "@modelcontextprotocol/server-filesystem"}, Env: []acp.EnvVariable{{Name: "PATH", Value: "/usr/bin"}}},
		},
	}
	data, err := json.Marshal(p)
	require.NoError(t, err)
	var got acp.NewSessionParams
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, "/home", got.CWD)
	require.Len(t, got.MCPServers, 1)
	require.Equal(t, "fs", got.MCPServers[0].Name)
}

func TestSessionIDParamsJSON(t *testing.T) {
	p := acp.SessionIDParams{SessionID: "ses-1"}
	data, err := json.Marshal(p)
	require.NoError(t, err)
	var got acp.SessionIDParams
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, "ses-1", got.SessionID)
}

func TestCancelParamsJSON(t *testing.T) {
	p := acp.CancelParams{SessionID: "ses-1", MessageID: "msg-1"}
	data, err := json.Marshal(p)
	require.NoError(t, err)
	var got acp.CancelParams
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, "ses-1", got.SessionID)
	require.Equal(t, "msg-1", got.MessageID)
}

func TestNewSessionRequestJSON(t *testing.T) {
	p := acp.NewSessionRequest{
		CWD: "/home",
		MCPServers: []acp.MCPServer{
			{Name: "filesystem", Command: "npx", Args: []string{"-y", "@modelcontextprotocol/server-filesystem"}},
		},
	}
	data, err := json.Marshal(p)
	require.NoError(t, err)
	var got acp.NewSessionRequest
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, "/home", got.CWD)
}

func TestSetConfigOptionRequestJSON(t *testing.T) {
	p := acp.SetConfigOptionRequest{
		SessionID: "ses-1",
		ConfigID:  "model",
		Value:     "gpt-4",
	}
	data, err := json.Marshal(p)
	require.NoError(t, err)
	var got acp.SetConfigOptionRequest
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, "model", got.ConfigID)
}

func TestContentPromptRequestJSON(t *testing.T) {
	p := acp.ContentPromptRequest{
		SessionID: "ses-1",
		Prompt: []acp.ContentBlock{
			{Type: "text", Text: "hello"},
		},
	}
	data, err := json.Marshal(p)
	require.NoError(t, err)
	var got acp.ContentPromptRequest
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, "ses-1", got.SessionID)
	require.Len(t, got.Prompt, 1)
}

func TestWriteTextFileResponseJSON(t *testing.T) {
	p := acp.WriteTextFileResponse{}
	data, err := json.Marshal(p)
	require.NoError(t, err)
	var got acp.WriteTextFileResponse
	require.NoError(t, json.Unmarshal(data, &got))
}

func TestLogoutResponseJSON(t *testing.T) {
	p := acp.LogoutResponse{}
	data, err := json.Marshal(p)
	require.NoError(t, err)
	var got acp.LogoutResponse
	require.NoError(t, json.Unmarshal(data, &got))
}

func TestAvailableCommandJSON(t *testing.T) {
	cmd := acp.AvailableCommand{
		Name:        "create_plan",
		Description: "Create an execution plan",
	}
	data, err := json.Marshal(cmd)
	require.NoError(t, err)

	var got acp.AvailableCommand
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, "create_plan", got.Name)
	require.Equal(t, "Create an execution plan", got.Description)
}

func TestAvailableCommandWithInputJSON(t *testing.T) {
	p := acp.AvailableCommand{
		Name:        "think",
		Description: "Think about a problem",
		Input: &acp.AvailableCommandInput{
			Unstructured: &acp.UnstructuredCommandInput{Hint: "What should I think about?"},
		},
	}
	data, err := json.Marshal(p)
	require.NoError(t, err)
	var got acp.AvailableCommand
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, "think", got.Name)
	require.NotNil(t, got.Input)
	require.NotNil(t, got.Input.Unstructured)
	require.Equal(t, "What should I think about?", got.Input.Unstructured.Hint)
}

func TestExtRequestJSON(t *testing.T) {
	p := acp.ExtRequest{Method: "custom/do_thing", Params: json.RawMessage(`{"key":"val"}`)}
	data, err := json.Marshal(p)
	require.NoError(t, err)
	var got acp.ExtRequest
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, "custom/do_thing", got.Method)
}
