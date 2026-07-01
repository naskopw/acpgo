package acp

import "encoding/json"

// ---- Client → Agent Requests ----

// InitializeRequest starts the handshake.
type InitializeRequest struct {
	ProtocolVersion   int                `json:"protocolVersion"`
	ClientCapabilities *ClientCapabilities `json:"clientCapabilities,omitempty"`
	ClientInfo        *Implementation    `json:"clientInfo,omitempty"`
	Meta              map[string]any     `json:"_meta,omitempty"`
}

// InitializeResponse is the result of a successful initialize.
type InitializeResponse struct {
	ProtocolVersion  int                 `json:"protocolVersion"`
	AgentCapabilities *AgentCapabilities `json:"agentCapabilities,omitempty"`
	AuthMethods      []AuthMethod        `json:"authMethods,omitempty"`
	AgentInfo        *Implementation     `json:"agentInfo,omitempty"`
	Meta             map[string]any      `json:"_meta,omitempty"`
}

// AuthenticateRequest authenticates with the agent.
type AuthenticateRequest struct {
	MethodID string         `json:"methodId"`
	Meta     map[string]any `json:"_meta,omitempty"`
}

// AuthenticateResponse is the result of authentication.
type AuthenticateResponse struct {
	Meta map[string]any `json:"_meta,omitempty"`
}

// LogoutRequest ends the authenticated session.
type LogoutRequest struct {
	Meta map[string]any `json:"_meta,omitempty"`
}

// LogoutResponse is the result of a successful logout.
type LogoutResponse struct {
	Meta map[string]any `json:"_meta,omitempty"`
}

// NewSessionRequest creates a new session.
type NewSessionRequest struct {
	CWD                  string         `json:"cwd"`
	MCPServers           []MCPServer    `json:"mcpServers"`
	AdditionalDirectories []string       `json:"additionalDirectories,omitempty"`
	Meta                 map[string]any `json:"_meta,omitempty"`
}

// NewSessionResponse is the result of creating a session.
type NewSessionResponse struct {
	SessionID    string         `json:"sessionId"`
	ConfigOptions []ConfigOption `json:"configOptions,omitempty"`
	Modes        *SessionModeState `json:"modes,omitempty"`
	Meta         map[string]any `json:"_meta,omitempty"`
}

// LoadSessionRequest loads an existing session.
type LoadSessionRequest struct {
	SessionID            string         `json:"sessionId"`
	CWD                  string         `json:"cwd"`
	MCPServers           []MCPServer    `json:"mcpServers"`
	AdditionalDirectories []string       `json:"additionalDirectories,omitempty"`
	Meta                 map[string]any `json:"_meta,omitempty"`
}

// LoadSessionResponse is the result of loading a session.
type LoadSessionResponse struct {
	ConfigOptions []ConfigOption   `json:"configOptions,omitempty"`
	Modes        *SessionModeState `json:"modes,omitempty"`
	Meta         map[string]any   `json:"_meta,omitempty"`
}

// ListSessionsRequest lists existing sessions.
type ListSessionsRequest struct {
	CWD    string         `json:"cwd,omitempty"`
	Cursor string         `json:"cursor,omitempty"`
	Meta   map[string]any `json:"_meta,omitempty"`
}

// ListSessionsResponse is the result of listing sessions.
type ListSessionsResponse struct {
	Sessions   []SessionInfo  `json:"sessions"`
	NextCursor string         `json:"nextCursor,omitempty"`
	Meta       map[string]any `json:"_meta,omitempty"`
}

// DeleteSessionRequest deletes a session.
type DeleteSessionRequest struct {
	SessionID string         `json:"sessionId"`
	Meta      map[string]any `json:"_meta,omitempty"`
}

// DeleteSessionResponse is the result of deleting a session.
type DeleteSessionResponse struct {
	Meta map[string]any `json:"_meta,omitempty"`
}

// ResumeSessionRequest resumes an existing session.
type ResumeSessionRequest struct {
	SessionID            string         `json:"sessionId"`
	CWD                  string         `json:"cwd"`
	MCPServers           []MCPServer    `json:"mcpServers,omitempty"`
	AdditionalDirectories []string       `json:"additionalDirectories,omitempty"`
	Meta                 map[string]any `json:"_meta,omitempty"`
}

// ResumeSessionResponse is the result of resuming a session.
type ResumeSessionResponse struct {
	ConfigOptions []ConfigOption   `json:"configOptions,omitempty"`
	Modes        *SessionModeState `json:"modes,omitempty"`
	Meta         map[string]any   `json:"_meta,omitempty"`
}

// CloseSessionRequest closes an active session.
type CloseSessionRequest struct {
	SessionID string         `json:"sessionId"`
	Meta      map[string]any `json:"_meta,omitempty"`
}

// CloseSessionResponse is the result of closing a session.
type CloseSessionResponse struct {
	Meta map[string]any `json:"_meta,omitempty"`
}

// PromptRequest contains parameters for a prompt request.
type PromptRequest struct {
	RequestID string  `json:"requestId"`
	SessionID string  `json:"sessionId"`
	Prompt    string  `json:"prompt"`
	ModelID   *string `json:"modelId,omitempty"`
}

// ContentPromptRequest is the ACP-standard session/prompt request with content blocks.
type ContentPromptRequest struct {
	SessionID string         `json:"sessionId"`
	Prompt    []ContentBlock `json:"prompt"`
	Meta      map[string]any `json:"_meta,omitempty"`
}

// PromptResponse is the result of processing a prompt.
type PromptResponse struct {
	StopReason string         `json:"stopReason"`
	Meta       map[string]any `json:"_meta,omitempty"`
}

// SetModeRequest sets the current mode for a session.
type SetModeRequest struct {
	SessionID string         `json:"sessionId"`
	ModeID    string         `json:"modeId"`
	Meta      map[string]any `json:"_meta,omitempty"`
}

// SetSessionModeResponse is the result of setting the session mode.
type SetSessionModeResponse struct {
	Meta map[string]any `json:"_meta,omitempty"`
}

// SetConfigOptionRequest sets a config option value.
type SetConfigOptionRequest struct {
	SessionID string         `json:"sessionId"`
	ConfigID  string         `json:"configId"`
	Value     string         `json:"value"`
	Meta      map[string]any `json:"_meta,omitempty"`
}

// SetConfigOptionResponse returns all config options after setting one.
type SetConfigOptionResponse struct {
	ConfigOptions []ConfigOption `json:"configOptions"`
	Meta          map[string]any `json:"_meta,omitempty"`
}

// ---- Agent → Client Requests ----

// ReadTextFileRequest reads a text file.
type ReadTextFileRequest struct {
	SessionID string         `json:"sessionId"`
	Path      string         `json:"path"`
	Line      int            `json:"line,omitempty"`
	Limit     int            `json:"limit,omitempty"`
	Meta      map[string]any `json:"_meta,omitempty"`
}

// ReadTextFileResponse contains file content.
type ReadTextFileResponse struct {
	Content string         `json:"content"`
	Meta    map[string]any `json:"_meta,omitempty"`
}

// WriteTextFileRequest writes content to a text file.
type WriteTextFileRequest struct {
	SessionID string         `json:"sessionId"`
	Path      string         `json:"path"`
	Content   string         `json:"content"`
	Meta      map[string]any `json:"_meta,omitempty"`
}

// WriteTextFileResponse is the result of writing a text file.
type WriteTextFileResponse struct {
	Meta map[string]any `json:"_meta,omitempty"`
}

// RequestPermissionRequest asks for user permission.
type RequestPermissionRequest struct {
	SessionID string            `json:"sessionId"`
	ToolCall  *ToolCallUpdate   `json:"toolCall"`
	Options   []PermissionOption `json:"options"`
	Meta      map[string]any    `json:"_meta,omitempty"`
}

// ---- Notifications ----

// CancelNotification cancels an ongoing prompt turn (client → agent).
type CancelNotification struct {
	SessionID string         `json:"sessionId"`
	Meta      map[string]any `json:"_meta,omitempty"`
}

// CancelRequestNotification cancels any in-flight request (protocol level).
type CancelRequestNotification struct {
	RequestID RequestID      `json:"requestId"`
	Meta      map[string]any `json:"_meta,omitempty"`
}

// RequestID is a JSON-RPC request identifier.
type RequestID any

// ---- Shared Types ----

// MCPServer describes an MCP server connection.
type MCPServer struct {
	Type    string         `json:"type,omitempty"`
	Name    string         `json:"name"`
	Command string         `json:"command,omitempty"`
	Args    []string       `json:"args,omitempty"`
	Env     []EnvVariable  `json:"env,omitempty"`
	URL     string         `json:"url,omitempty"`
	Headers []HTTPHeader   `json:"headers,omitempty"`
	Meta    map[string]any `json:"_meta,omitempty"`
}

// HTTPHeader is an HTTP header for MCP server connections.
type HTTPHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// EnvVariable represents an environment variable.
type EnvVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Harness interface parameter types (simplified for the Go harness interface).

// NewSessionParams contains harness-level session creation parameters.
type NewSessionParams struct {
	CWD        string      `json:"cwd"`
	MCPServers []MCPServer `json:"mcpServers,omitempty"`
}

// CancelParams identifies a session (and optionally message) to cancel.
type CancelParams struct {
	SessionID string `json:"sessionId"`
	MessageID string `json:"messageId,omitempty"`
}

// SetModelParams contains the model ID to set.
type SetModelParams struct {
	ModelID string `json:"modelId"`
}

// InitializeResult is the result of an initialize handshake.
type InitializeResult struct {
	ProtocolVersion   int                `json:"protocolVersion"`
	AgentCapabilities *AgentCapabilities `json:"agentCapabilities,omitempty"`
	AgentName         string             `json:"agentName,omitempty"`
	AgentVersion      string             `json:"agentVersion,omitempty"`
}

// SessionIDParams identifies a session by ID.
type SessionIDParams struct {
	SessionID string `json:"sessionId"`
}

// SessionInfo describes a session in list results.
type SessionInfo struct {
	SessionID            string            `json:"sessionId"`
	CWD                  string            `json:"cwd"`
	AdditionalDirectories []string          `json:"additionalDirectories,omitempty"`
	Title                string            `json:"title,omitempty"`
	UpdatedAt            string            `json:"updatedAt,omitempty"`
	ConfigOptions        []ConfigOption    `json:"configOptions,omitempty"`
	Meta                 map[string]any    `json:"_meta,omitempty"`
}

// AvailableCommand describes a command the agent can execute.
type AvailableCommand struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Input       *AvailableCommandInput `json:"input,omitempty"`
	Meta        map[string]any         `json:"_meta,omitempty"`
}

// AvailableCommandInput describes the input format for a slash command.
type AvailableCommandInput struct {
	Unstructured *UnstructuredCommandInput `json:"unstructured,omitempty"`
	Meta         map[string]any            `json:"_meta,omitempty"`
}

// UnstructuredCommandInput describes an unstructured text input for a command.
type UnstructuredCommandInput struct {
	Hint string         `json:"hint,omitempty"`
	Meta map[string]any `json:"_meta,omitempty"`
}

// ExtRequest is a generic extension method request.
type ExtRequest struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
	Meta   map[string]any  `json:"_meta,omitempty"`
}

// ExtResponse is a generic extension method response.
type ExtResponse struct {
	Result json.RawMessage `json:"result,omitempty"`
	Error  *RPCError       `json:"error,omitempty"`
	Meta   map[string]any  `json:"_meta,omitempty"`
}

// ExtNotification is a generic extension method notification.
type ExtNotification struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
	Meta   map[string]any  `json:"_meta,omitempty"`
}

// SessionModeState describes the available and current session modes.
type SessionModeState struct {
	AvailableModes []SessionMode `json:"availableModes"`
	CurrentModeID  string        `json:"currentModeId,omitempty"`
	Meta           map[string]any `json:"_meta,omitempty"`
}

// SessionMode describes a single session mode.
type SessionMode struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Meta        map[string]any `json:"_meta,omitempty"`
}
