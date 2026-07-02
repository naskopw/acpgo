package acp_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	acp "github.com/naskopw/acpgo"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/require"
)

func compileSchemaDef(t *testing.T, defPath string) *jsonschema.Schema {
	t.Helper()
	c := jsonschema.NewCompiler()

	data, err := os.ReadFile("testdata/schema/v1/schema.json")
	require.NoError(t, err)

	var doc any
	require.NoError(t, json.Unmarshal(data, &doc))
	require.NoError(t, c.AddResource("schema.json", doc))

	sch, err := c.Compile("schema.json#/$defs/" + defPath)
	require.NoError(t, err)
	return sch
}

type goCase struct {
	name    string
	val     any
	defPath string
}

func allRoundTripCases() []goCase {
	return []goCase{
		// ── Client → Agent Requests ──
		{name: "InitializeRequest", val: acp.InitializeRequest{ProtocolVersion: 1, ClientInfo: &acp.Implementation{Name: "test", Version: "1.0"}}, defPath: "InitializeRequest"},
		{name: "AuthenticateRequest", val: acp.AuthenticateRequest{MethodID: "password"}, defPath: "AuthenticateRequest"},
		{name: "LogoutRequest", val: acp.LogoutRequest{}, defPath: "LogoutRequest"},
		{name: "NewSessionRequest", val: acp.NewSessionRequest{CWD: "/home", MCPServers: []acp.MCPServer{{Name: "fs", Command: "/bin/ls", Args: []string{"-la"}, Env: []acp.EnvVariable{{Name: "PATH", Value: "/usr/bin"}}}}}, defPath: "NewSessionRequest"},
		{name: "LoadSessionRequest", val: acp.LoadSessionRequest{SessionID: "s1", CWD: "/home", MCPServers: []acp.MCPServer{{Name: "fs", Command: "/bin/ls", Args: []string{"-la"}, Env: []acp.EnvVariable{{Name: "PATH", Value: "/usr/bin"}}}}}, defPath: "LoadSessionRequest"},
		{name: "ListSessionsRequest", val: acp.ListSessionsRequest{CWD: "/home"}, defPath: "ListSessionsRequest"},
		{name: "DeleteSessionRequest", val: acp.DeleteSessionRequest{SessionID: "s1"}, defPath: "DeleteSessionRequest"},
		{name: "ResumeSessionRequest", val: acp.ResumeSessionRequest{SessionID: "s1", CWD: "/home"}, defPath: "ResumeSessionRequest"},
		{name: "CloseSessionRequest", val: acp.CloseSessionRequest{SessionID: "s1"}, defPath: "CloseSessionRequest"},
		{name: "ContentPromptRequest", val: acp.ContentPromptRequest{SessionID: "s1", Prompt: []acp.ContentBlock{{Type: "text", Text: "hello"}}}, defPath: "PromptRequest"},
		{name: "SetModeRequest", val: acp.SetModeRequest{SessionID: "s1", ModeID: "code"}, defPath: "SetSessionModeRequest"},
		{name: "SetConfigOptionRequest", val: acp.SetConfigOptionRequest{SessionID: "s1", ConfigID: "model", Value: "gpt-4"}, defPath: "SetSessionConfigOptionRequest"},

		// ── Agent → Client Requests ──
		{name: "ReadTextFileRequest", val: acp.ReadTextFileRequest{SessionID: "s1", Path: "/foo.txt"}, defPath: "ReadTextFileRequest"},
		{name: "WriteTextFileRequest", val: acp.WriteTextFileRequest{SessionID: "s1", Path: "/foo.txt", Content: "data"}, defPath: "WriteTextFileRequest"},
		{name: "RequestPermissionRequest", val: acp.RequestPermissionRequest{SessionID: "s1", ToolCall: &acp.ToolCallUpdate{ToolCallID: "tc1"}, Options: []acp.PermissionOption{{OptionID: "opt1", Name: "Allow", Kind: "allow_once"}}}, defPath: "RequestPermissionRequest"},
		{name: "CreateTerminalRequest", val: acp.CreateTerminalRequest{SessionID: "s1", Command: "/bin/bash"}, defPath: "CreateTerminalRequest"},
		{name: "TerminalOutputRequest", val: acp.TerminalOutputRequest{SessionID: "s1", TerminalID: "t1"}, defPath: "TerminalOutputRequest"},
		{name: "ReleaseTerminalRequest", val: acp.ReleaseTerminalRequest{SessionID: "s1", TerminalID: "t1"}, defPath: "ReleaseTerminalRequest"},
		{name: "WaitForTerminalExitRequest", val: acp.WaitForTerminalExitRequest{SessionID: "s1", TerminalID: "t1"}, defPath: "WaitForTerminalExitRequest"},
		{name: "KillTerminalRequest", val: acp.KillTerminalRequest{SessionID: "s1", TerminalID: "t1"}, defPath: "KillTerminalRequest"},

		// ── Agent → Client Responses ──
		{name: "InitializeResponse", val: acp.InitializeResponse{ProtocolVersion: 1, AgentInfo: &acp.Implementation{Name: "agent", Version: "1.0"}}, defPath: "InitializeResponse"},
		{name: "AuthenticateResponse", val: acp.AuthenticateResponse{}, defPath: "AuthenticateResponse"},
		{name: "LogoutResponse", val: acp.LogoutResponse{}, defPath: "LogoutResponse"},
		{name: "NewSessionResponse", val: acp.NewSessionResponse{SessionID: "s1"}, defPath: "NewSessionResponse"},
		{name: "LoadSessionResponse", val: acp.LoadSessionResponse{}, defPath: "LoadSessionResponse"},
		{name: "ListSessionsResponse", val: acp.ListSessionsResponse{Sessions: []acp.SessionInfo{}}, defPath: "ListSessionsResponse"},
		{name: "DeleteSessionResponse", val: acp.DeleteSessionResponse{}, defPath: "DeleteSessionResponse"},
		{name: "ResumeSessionResponse", val: acp.ResumeSessionResponse{}, defPath: "ResumeSessionResponse"},
		{name: "CloseSessionResponse", val: acp.CloseSessionResponse{}, defPath: "CloseSessionResponse"},
		{name: "SetSessionModeResponse", val: acp.SetSessionModeResponse{}, defPath: "SetSessionModeResponse"},
		{name: "SetConfigOptionResponse", val: acp.SetConfigOptionResponse{ConfigOptions: []acp.ConfigOption{}}, defPath: "SetSessionConfigOptionResponse"},
		{name: "PromptResponse", val: acp.PromptResponse{StopReason: "end_turn"}, defPath: "PromptResponse"},

		// ── Client → Agent Responses ──
		{name: "ReadTextFileResponse", val: acp.ReadTextFileResponse{Content: "file content"}, defPath: "ReadTextFileResponse"},
		{name: "WriteTextFileResponse", val: acp.WriteTextFileResponse{}, defPath: "WriteTextFileResponse"},
		{name: "RequestPermissionResponse", val: acp.RequestPermissionResponse{Outcome: &acp.PermissionOutcome{Outcome: "selected", OptionID: "opt1"}}, defPath: "RequestPermissionResponse"},
		{name: "CreateTerminalResponse", val: acp.CreateTerminalResponse{TerminalID: "t1"}, defPath: "CreateTerminalResponse"},
		{name: "TerminalOutputResponse", val: acp.TerminalOutputResponse{Output: "hello", Truncated: false}, defPath: "TerminalOutputResponse"},
		{name: "ReleaseTerminalResponse", val: acp.ReleaseTerminalResponse{}, defPath: "ReleaseTerminalResponse"},
		{name: "WaitForTerminalExitResponse", val: acp.WaitForTerminalExitResponse{}, defPath: "WaitForTerminalExitResponse"},
		{name: "KillTerminalResponse", val: acp.KillTerminalResponse{}, defPath: "KillTerminalResponse"},

		// ── Notifications ──
		{name: "CancelNotification", val: acp.CancelNotification{SessionID: "s1"}, defPath: "CancelNotification"},
		{name: "CancelRequestNotification", val: acp.CancelRequestNotification{RequestID: "req1"}, defPath: "CancelRequestNotification"},
		{name: "SessionNotification", val: acp.SessionNotification{SessionID: "s1", Update: acp.SessionUpdate{SessionUpdateVariant: "tool_call", ToolCallID: "tc1", Title: "list"}}, defPath: "SessionNotification"},

		// ── Shared / Nested Types ──
		{name: "ContentBlock/text", val: acp.ContentBlock{Type: "text", Text: "hello"}, defPath: "ContentBlock"},
		{name: "ContentBlock/image", val: acp.ContentBlock{Type: "image", Data: "base64data", MimeType: "image/png"}, defPath: "ContentBlock"},
		{name: "ContentBlock/audio", val: acp.ContentBlock{Type: "audio", Data: "base64data", MimeType: "audio/mp3"}, defPath: "ContentBlock"},
		{name: "ContentBlock/resource_link", val: acp.ContentBlock{Type: "resource_link", Name: "doc", URI: "file:///doc.md"}, defPath: "ContentBlock"},
		{name: "ToolCall", val: acp.ToolCall{ToolCallID: "tc1", Title: "List files", Kind: "read", Status: "completed"}, defPath: "ToolCall"},
		{name: "ToolCallUpdate", val: acp.ToolCallUpdate{ToolCallID: "tc1", Kind: "read", Status: "completed"}, defPath: "ToolCallUpdate"},
		{name: "ToolCallContent/content", val: acp.ToolCallContent{Type: "content", Content: &acp.ContentBlock{Type: "text", Text: "result"}}, defPath: "ToolCallContent"},
		{name: "ToolCallContent/diff", val: acp.ToolCallContent{Type: "diff", DiffPath: "/foo.txt", DiffNewText: "new content"}, defPath: "ToolCallContent"},
		{name: "ToolCallContent/terminal", val: acp.ToolCallContent{Type: "terminal", TerminalID: "t1"}, defPath: "ToolCallContent"},
		{name: "ToolCallLocation", val: acp.ToolCallLocation{Path: "/foo.txt", Line: 10}, defPath: "ToolCallLocation"},
		{name: "PermissionOption", val: acp.PermissionOption{OptionID: "opt1", Name: "Allow", Kind: "allow_once"}, defPath: "PermissionOption"},
		{name: "PermissionOutcome/selected", val: acp.PermissionOutcome{Outcome: "selected", OptionID: "opt1"}, defPath: "RequestPermissionOutcome"},
		{name: "PermissionOutcome/cancelled", val: acp.PermissionOutcome{Outcome: "cancelled"}, defPath: "RequestPermissionOutcome"},
		{name: "PlanEntry", val: acp.PlanEntry{Content: "Do the thing", Priority: "high", Status: "pending"}, defPath: "PlanEntry"},
		{name: "Cost", val: acp.Cost{Amount: 0.01, Currency: "USD"}, defPath: "Cost"},
		{name: "Annotations", val: acp.Annotations{Audience: []string{"assistant"}}, defPath: "Annotations"},
		{name: "EnvVariable", val: acp.EnvVariable{Name: "PATH", Value: "/usr/bin"}, defPath: "EnvVariable"},
		{name: "HTTPHeader", val: acp.HTTPHeader{Name: "Authorization", Value: "Bearer token"}, defPath: "HttpHeader"},
		{name: "MCPServer/stdio", val: acp.MCPServer{Name: "fs", Command: "/bin/ls", Args: []string{"-la"}, Env: []acp.EnvVariable{{Name: "PATH", Value: "/usr/bin"}}}, defPath: "McpServer"},
		{name: "MCPServer/http", val: acp.MCPServer{Type: "http", Name: "remote", URL: "https://example.com", Headers: []acp.HTTPHeader{{Name: "Authorization", Value: "Bearer token"}}}, defPath: "McpServer"},
		{name: "ConfigOption", val: acp.ConfigOption{ID: "model", Name: "Model", Type: "select", CurrentValue: "gpt-4", Options: []acp.ConfigOptionValue{{Value: "gpt-4", Name: "GPT-4"}}, Description: "The model"}, defPath: "SessionConfigOption"},
		{name: "ConfigOptionValue", val: acp.ConfigOptionValue{Value: "gpt-4", Name: "GPT-4", Description: "GPT-4 model"}, defPath: "SessionConfigSelectOption"},
		{name: "ClientCapabilities", val: acp.ClientCapabilities{FS: &acp.FileSystemCapabilities{ReadTextFile: true}, Terminal: true}, defPath: "ClientCapabilities"},
		{name: "AgentCapabilities", val: acp.AgentCapabilities{LoadSession: true}, defPath: "AgentCapabilities"},
		{name: "FileSystemCapabilities", val: acp.FileSystemCapabilities{ReadTextFile: true}, defPath: "FileSystemCapabilities"},
		{name: "PromptCapabilities", val: acp.PromptCapabilities{Image: true}, defPath: "PromptCapabilities"},
		{name: "MCPCapabilities", val: acp.MCPCapabilities{HTTP: true}, defPath: "McpCapabilities"},
		{name: "SessionCapabilities", val: acp.SessionCapabilities{}, defPath: "SessionCapabilities"},
		{name: "AgentAuthCapabilities", val: acp.AgentAuthCapabilities{}, defPath: "AgentAuthCapabilities"},
		{name: "LogoutCapabilities", val: acp.LogoutCapabilities{}, defPath: "LogoutCapabilities"},
		{name: "Implementation", val: acp.Implementation{Name: "test", Version: "1.0"}, defPath: "Implementation"},
		{name: "AuthMethod", val: acp.AuthMethod{Type: "agent", ID: "password", Name: "Password"}, defPath: "AuthMethod"},
		{name: "SessionInfo", val: acp.SessionInfo{SessionID: "s1", CWD: "/home"}, defPath: "SessionInfo"},
		{name: "SessionModeState", val: acp.SessionModeState{CurrentModeID: "code", AvailableModes: []acp.SessionMode{{ID: "code", Name: "Code"}}}, defPath: "SessionModeState"},
		{name: "SessionMode", val: acp.SessionMode{ID: "code", Name: "Code", Description: "Code mode"}, defPath: "SessionMode"},
		{name: "AvailableCommand", val: acp.AvailableCommand{Name: "think", Description: "Think about a problem"}, defPath: "AvailableCommand"},
		{name: "UnstructuredCommandInput", val: acp.UnstructuredCommandInput{Hint: "Enter text"}, defPath: "UnstructuredCommandInput"},
		{name: "ContentChunk", val: acp.ContentChunk{Content: &acp.ContentBlock{Type: "text", Text: "chunk"}}, defPath: "ContentChunk"},
		{name: "UsageUpdate", val: acp.UsageUpdate{Used: 100, Size: 1000}, defPath: "UsageUpdate"},
		{name: "TerminalExitStatus", val: acp.TerminalExitStatus{ExitCode: 1, Signal: "SIGTERM"}, defPath: "TerminalExitStatus"},
	}
}

func TestSchemaRoundTrip(t *testing.T) {
	for _, tc := range allRoundTripCases() {
		t.Run(tc.name, func(t *testing.T) {
			sch := compileSchemaDef(t, tc.defPath)

			data, err := json.Marshal(tc.val)
			require.NoError(t, err)

			var v any
			require.NoError(t, json.Unmarshal(data, &v))
			err = sch.Validate(v)
			require.NoError(t, err, "schema violation for %s\nJSON: %s", tc.name, string(data))

			typ := reflect.TypeOf(tc.val)
			got := reflect.New(typ).Interface()
			require.NoError(t, json.Unmarshal(data, got))

			require.Equal(t, tc.val, reflect.ValueOf(got).Elem().Interface())
		})
	}
}

// fixture maps a fixture filename (without .json) to its schema $defs path
// and the expected Go type to unmarshal into.
type fixtureCase struct {
	defPath  string
	emptyVal any // zero-value instance of the expected Go type
}

func fixtureMap() map[string]fixtureCase {
	return map[string]fixtureCase{
		// Client → Agent Requests
		"request_initialize":         {defPath: "InitializeRequest", emptyVal: acp.InitializeRequest{}},
		"request_authenticate":       {defPath: "AuthenticateRequest", emptyVal: acp.AuthenticateRequest{}},
		"request_logout":             {defPath: "LogoutRequest", emptyVal: acp.LogoutRequest{}},
		"request_new_session":        {defPath: "NewSessionRequest", emptyVal: acp.NewSessionRequest{}},
		"request_load_session":       {defPath: "LoadSessionRequest", emptyVal: acp.LoadSessionRequest{}},
		"request_list_sessions":      {defPath: "ListSessionsRequest", emptyVal: acp.ListSessionsRequest{}},
		"request_delete_session":     {defPath: "DeleteSessionRequest", emptyVal: acp.DeleteSessionRequest{}},
		"request_resume_session":     {defPath: "ResumeSessionRequest", emptyVal: acp.ResumeSessionRequest{}},
		"request_close_session":      {defPath: "CloseSessionRequest", emptyVal: acp.CloseSessionRequest{}},
		"request_prompt":             {defPath: "PromptRequest", emptyVal: acp.ContentPromptRequest{}},
		"request_set_mode":           {defPath: "SetSessionModeRequest", emptyVal: acp.SetModeRequest{}},
		"request_set_config_option":  {defPath: "SetSessionConfigOptionRequest", emptyVal: acp.SetConfigOptionRequest{}},
		"request_read_text_file":     {defPath: "ReadTextFileRequest", emptyVal: acp.ReadTextFileRequest{}},
		"request_write_text_file":    {defPath: "WriteTextFileRequest", emptyVal: acp.WriteTextFileRequest{}},
		"request_permission":         {defPath: "RequestPermissionRequest", emptyVal: acp.RequestPermissionRequest{}},
		"request_create_terminal":    {defPath: "CreateTerminalRequest", emptyVal: acp.CreateTerminalRequest{}},
		"request_terminal_output":    {defPath: "TerminalOutputRequest", emptyVal: acp.TerminalOutputRequest{}},
		"request_release_terminal":   {defPath: "ReleaseTerminalRequest", emptyVal: acp.ReleaseTerminalRequest{}},
		"request_wait_terminal_exit": {defPath: "WaitForTerminalExitRequest", emptyVal: acp.WaitForTerminalExitRequest{}},
		"request_kill_terminal":      {defPath: "KillTerminalRequest", emptyVal: acp.KillTerminalRequest{}},

		// Agent → Client Responses
		"response_initialize":        {defPath: "InitializeResponse", emptyVal: acp.InitializeResponse{}},
		"response_authenticate":      {defPath: "AuthenticateResponse", emptyVal: acp.AuthenticateResponse{}},
		"response_logout":            {defPath: "LogoutResponse", emptyVal: acp.LogoutResponse{}},
		"response_new_session":       {defPath: "NewSessionResponse", emptyVal: acp.NewSessionResponse{}},
		"response_load_session":      {defPath: "LoadSessionResponse", emptyVal: acp.LoadSessionResponse{}},
		"response_list_sessions":     {defPath: "ListSessionsResponse", emptyVal: acp.ListSessionsResponse{}},
		"response_delete_session":    {defPath: "DeleteSessionResponse", emptyVal: acp.DeleteSessionResponse{}},
		"response_resume_session":    {defPath: "ResumeSessionResponse", emptyVal: acp.ResumeSessionResponse{}},
		"response_close_session":     {defPath: "CloseSessionResponse", emptyVal: acp.CloseSessionResponse{}},
		"response_set_mode":          {defPath: "SetSessionModeResponse", emptyVal: acp.SetSessionModeResponse{}},
		"response_set_config_option": {defPath: "SetSessionConfigOptionResponse", emptyVal: acp.SetConfigOptionResponse{}},
		"response_prompt":            {defPath: "PromptResponse", emptyVal: acp.PromptResponse{}},

		// Client → Agent Responses
		"response_read_text_file":     {defPath: "ReadTextFileResponse", emptyVal: acp.ReadTextFileResponse{}},
		"response_write_text_file":    {defPath: "WriteTextFileResponse", emptyVal: acp.WriteTextFileResponse{}},
		"response_request_permission": {defPath: "RequestPermissionResponse", emptyVal: acp.RequestPermissionResponse{}},
		"response_create_terminal":    {defPath: "CreateTerminalResponse", emptyVal: acp.CreateTerminalResponse{}},
		"response_terminal_output":    {defPath: "TerminalOutputResponse", emptyVal: acp.TerminalOutputResponse{}},
		"response_release_terminal":   {defPath: "ReleaseTerminalResponse", emptyVal: acp.ReleaseTerminalResponse{}},
		"response_wait_terminal_exit": {defPath: "WaitForTerminalExitResponse", emptyVal: acp.WaitForTerminalExitResponse{}},
		"response_kill_terminal":      {defPath: "KillTerminalResponse", emptyVal: acp.KillTerminalResponse{}},

		// Notifications
		"notification_cancel":         {defPath: "CancelNotification", emptyVal: acp.CancelNotification{}},
		"notification_cancel_request": {defPath: "CancelRequestNotification", emptyVal: acp.CancelRequestNotification{}},
		"notification_session_update": {defPath: "SessionNotification", emptyVal: acp.SessionNotification{}},

		// Shared types
		"content_block_text":           {defPath: "ContentBlock", emptyVal: acp.ContentBlock{}},
		"content_block_image":          {defPath: "ContentBlock", emptyVal: acp.ContentBlock{}},
		"tool_call":                    {defPath: "ToolCall", emptyVal: acp.ToolCall{}},
		"tool_call_update":             {defPath: "ToolCallUpdate", emptyVal: acp.ToolCallUpdate{}},
		"tool_call_content_content":    {defPath: "ToolCallContent", emptyVal: acp.ToolCallContent{}},
		"tool_call_content_diff":       {defPath: "ToolCallContent", emptyVal: acp.ToolCallContent{}},
		"tool_call_content_terminal":   {defPath: "ToolCallContent", emptyVal: acp.ToolCallContent{}},
		"permission_option":            {defPath: "PermissionOption", emptyVal: acp.PermissionOption{}},
		"permission_outcome_selected":  {defPath: "RequestPermissionOutcome", emptyVal: acp.PermissionOutcome{}},
		"permission_outcome_cancelled": {defPath: "RequestPermissionOutcome", emptyVal: acp.PermissionOutcome{}},
		"plan_entry":                   {defPath: "PlanEntry", emptyVal: acp.PlanEntry{}},
		"session_update_tool_call":     {defPath: "SessionUpdate", emptyVal: acp.SessionUpdate{}},
		"session_update_agent_message": {defPath: "SessionUpdate", emptyVal: acp.SessionUpdate{}},
	}
}

func TestValidation_Negative(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want error
	}{
		{name: "ContentBlock/invalid_type", val: acp.ContentBlock{Type: "invalid"}, want: acp.ErrInvalidContentType},
		{name: "ToolCall/invalid_kind", val: acp.ToolCall{ToolCallID: "t1", Title: "t", Kind: "invalid"}, want: acp.ErrInvalidToolKind},
		{name: "ToolCall/invalid_status", val: acp.ToolCall{ToolCallID: "t1", Title: "t", Status: "invalid"}, want: acp.ErrInvalidToolStatus},
		{name: "ToolCallUpdate/invalid_kind", val: acp.ToolCallUpdate{ToolCallID: "t1", Kind: "invalid"}, want: acp.ErrInvalidToolKind},
		{name: "ToolCallUpdate/invalid_status", val: acp.ToolCallUpdate{ToolCallID: "t1", Status: "invalid"}, want: acp.ErrInvalidToolStatus},
		{name: "ToolCallContent/invalid_type", val: acp.ToolCallContent{Type: "invalid"}, want: acp.ErrInvalidToolContentType},
		{name: "PlanEntry/invalid_priority", val: acp.PlanEntry{Content: "c", Priority: "invalid", Status: "pending"}, want: acp.ErrInvalidPlanPriority},
		{name: "PlanEntry/invalid_status", val: acp.PlanEntry{Content: "c", Priority: "high", Status: "invalid"}, want: acp.ErrInvalidPlanStatus},
		{name: "PermissionOption/invalid_kind", val: acp.PermissionOption{OptionID: "o1", Name: "Allow", Kind: "invalid"}, want: acp.ErrInvalidPermissionKind},
		{name: "PermissionOutcome/invalid_outcome", val: acp.PermissionOutcome{Outcome: "invalid"}, want: acp.ErrInvalidPermissionOutcome},
		{name: "SessionUpdate/invalid_variant", val: acp.SessionUpdate{SessionUpdateVariant: "invalid"}, want: acp.ErrInvalidSessionUpdateVariant},
		{name: "PromptResponse/invalid_stop_reason", val: acp.PromptResponse{StopReason: "invalid"}, want: acp.ErrInvalidStopReason},
		{name: "ConfigOption/invalid_category", val: acp.ConfigOption{ID: "id", Name: "n", Type: "select", CurrentValue: "v", Options: []acp.ConfigOptionValue{}, Category: "invalid"}, want: acp.ErrInvalidConfigCategory},
		{name: "Annotations/invalid_audience", val: acp.Annotations{Audience: []string{"invalid"}}, want: acp.ErrInvalidAudience},
		{name: "MCPServer/invalid_type", val: acp.MCPServer{Type: "invalid", Name: "m"}, want: acp.ErrInvalidMCPServerType},
		{name: "NewSessionRequest/relative_cwd", val: acp.NewSessionRequest{CWD: "relative", MCPServers: []acp.MCPServer{}}, want: acp.ErrInvalidCWD},
		{name: "NewSessionRequest/relative_additional_dir", val: acp.NewSessionRequest{CWD: "/abs", AdditionalDirectories: []string{"relative"}, MCPServers: []acp.MCPServer{}}, want: acp.ErrInvalidAdditionalDirectory},
		{name: "NewSessionRequest/relative_mcp_command", val: acp.NewSessionRequest{CWD: "/abs", MCPServers: []acp.MCPServer{{Type: "stdio", Command: "relative"}}}, want: acp.ErrInvalidMCPCommand},
		{name: "LoadSessionRequest/relative_cwd", val: acp.LoadSessionRequest{SessionID: "s1", CWD: "relative", MCPServers: []acp.MCPServer{}}, want: acp.ErrInvalidCWD},
		{name: "ResumeSessionRequest/relative_cwd", val: acp.ResumeSessionRequest{SessionID: "s1", CWD: "relative", MCPServers: []acp.MCPServer{}}, want: acp.ErrInvalidCWD},
		{name: "ListSessionsRequest/relative_cwd", val: acp.ListSessionsRequest{CWD: "relative"}, want: acp.ErrInvalidCWD},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v, ok := tc.val.(interface{ Validate() error })
			if !ok {
				t.Fatal("type does not implement Validate()")
			}
			err := v.Validate()
			require.ErrorIs(t, err, tc.want, "expected %v, got %v", tc.want, err)
		})
	}
}

func TestValidation_Positive(t *testing.T) {
	tests := []struct {
		name string
		val  any
	}{
		{name: "ContentBlock/text", val: acp.ContentBlock{Type: "text"}},
		{name: "ContentBlock/image", val: acp.ContentBlock{Type: "image"}},
		{name: "ContentBlock/audio", val: acp.ContentBlock{Type: "audio"}},
		{name: "ContentBlock/resource_link", val: acp.ContentBlock{Type: "resource_link"}},
		{name: "ContentBlock/resource", val: acp.ContentBlock{Type: "resource"}},
		{name: "ToolCall/valid", val: acp.ToolCall{ToolCallID: "t1", Title: "t", Kind: "read", Status: "completed"}},
		{name: "ToolCall/empty_kind", val: acp.ToolCall{ToolCallID: "t1", Title: "t"}},
		{name: "ToolCall/empty_status", val: acp.ToolCall{ToolCallID: "t1", Title: "t", Kind: "read"}},
		{name: "ToolCallUpdate/valid", val: acp.ToolCallUpdate{ToolCallID: "t1", Kind: "read", Status: "completed"}},
		{name: "ToolCallUpdate/empty_optional", val: acp.ToolCallUpdate{ToolCallID: "t1"}},
		{name: "ToolCallContent/content", val: acp.ToolCallContent{Type: "content"}},
		{name: "ToolCallContent/diff", val: acp.ToolCallContent{Type: "diff"}},
		{name: "ToolCallContent/terminal", val: acp.ToolCallContent{Type: "terminal"}},
		{name: "PlanEntry/valid", val: acp.PlanEntry{Content: "c", Priority: "high", Status: "pending"}},
		{name: "PermissionOption/allow_once", val: acp.PermissionOption{OptionID: "o1", Name: "Allow", Kind: "allow_once"}},
		{name: "PermissionOption/allow_always", val: acp.PermissionOption{OptionID: "o1", Name: "Allow", Kind: "allow_always"}},
		{name: "PermissionOption/reject_once", val: acp.PermissionOption{OptionID: "o1", Name: "Reject", Kind: "reject_once"}},
		{name: "PermissionOption/reject_always", val: acp.PermissionOption{OptionID: "o1", Name: "Reject", Kind: "reject_always"}},
		{name: "PermissionOutcome/cancelled", val: acp.PermissionOutcome{Outcome: "cancelled"}},
		{name: "PermissionOutcome/selected", val: acp.PermissionOutcome{Outcome: "selected"}},
		{name: "SessionUpdate/user_message_chunk", val: acp.SessionUpdate{SessionUpdateVariant: "user_message_chunk"}},
		{name: "SessionUpdate/tool_call", val: acp.SessionUpdate{SessionUpdateVariant: "tool_call"}},
		{name: "PromptResponse/end_turn", val: acp.PromptResponse{StopReason: "end_turn"}},
		{name: "PromptResponse/max_tokens", val: acp.PromptResponse{StopReason: "max_tokens"}},
		{name: "ConfigOption/no_category", val: acp.ConfigOption{ID: "id", Name: "n", Type: "select", CurrentValue: "v", Options: []acp.ConfigOptionValue{}}},
		{name: "ConfigOption/valid_category", val: acp.ConfigOption{ID: "id", Name: "n", Type: "select", CurrentValue: "v", Options: []acp.ConfigOptionValue{}, Category: "mode"}},
		{name: "Annotations/empty", val: acp.Annotations{}},
		{name: "Annotations/valid", val: acp.Annotations{Audience: []string{"assistant", "user"}}},
		{name: "MCPServer/stdio", val: acp.MCPServer{Type: "stdio", Name: "m"}},
		{name: "MCPServer/http", val: acp.MCPServer{Type: "http", Name: "m", URL: "http://localhost"}},
		{name: "MCPServer/sse", val: acp.MCPServer{Type: "sse", Name: "m", URL: "http://localhost"}},
		{name: "MCPServer/empty_type", val: acp.MCPServer{Name: "m"}},
		{name: "NewSessionRequest/valid", val: acp.NewSessionRequest{CWD: "/abs", MCPServers: []acp.MCPServer{}}},
		{name: "LoadSessionRequest/valid", val: acp.LoadSessionRequest{SessionID: "s1", CWD: "/abs", MCPServers: []acp.MCPServer{}}},
		{name: "ResumeSessionRequest/valid", val: acp.ResumeSessionRequest{SessionID: "s1", CWD: "/abs", MCPServers: []acp.MCPServer{}}},
		{name: "ListSessionsRequest/valid", val: acp.ListSessionsRequest{CWD: "/abs"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v, ok := tc.val.(interface{ Validate() error })
			if !ok {
				t.Fatal("type does not implement Validate()")
			}
			require.NoError(t, v.Validate())
		})
	}
}

func TestFixtureValidation(t *testing.T) {
	fm := fixtureMap()
	dir := "testdata/fixtures/v1"

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		name := entry.Name()
		base := name[:len(name)-len(".json")]
		fc, ok := fm[base]
		if !ok {
			t.Logf("SKIP: no fixture mapping for %s", base)
			continue
		}

		t.Run(base, func(t *testing.T) {
			sch := compileSchemaDef(t, fc.defPath)

			data, err := os.ReadFile(filepath.Join(dir, name))
			require.NoError(t, err)

			// Skip empty fixtures (only `{}`)
			if len(data) <= 2 {
				t.Skip("empty fixture:", base)
			}

			// Validate against schema
			var v any
			require.NoError(t, json.Unmarshal(data, &v))
			err = sch.Validate(v)
			require.NoError(t, err, "schema validation failed for fixture %s", base)

			// Unmarshal to correct Go type
			typ := reflect.TypeOf(fc.emptyVal)
			got := reflect.New(typ).Interface()
			require.NoError(t, json.Unmarshal(data, got))

			// Re-marshal and compare for round-trip
			gotData, err := json.Marshal(got)
			require.NoError(t, err)

			var gotV any
			require.NoError(t, json.Unmarshal(gotData, &gotV))
			err = sch.Validate(gotV)
			require.NoError(t, err, "schema validation failed after round-trip for fixture %s", base)
		})
	}
}
