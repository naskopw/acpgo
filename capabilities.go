package acp

// ClientCapabilities describes what a client supports.
type ClientCapabilities struct {
	FS        *FileSystemCapabilities `json:"fs,omitempty"`
	Terminal  bool                    `json:"terminal,omitempty"`
	Meta      map[string]any          `json:"_meta,omitempty"`
}

// FileSystemCapabilities describes file system support.
type FileSystemCapabilities struct {
	ReadTextFile  bool `json:"readTextFile,omitempty"`
	WriteTextFile bool `json:"writeTextFile,omitempty"`
	Meta          map[string]any `json:"_meta,omitempty"`
}

// AgentCapabilities describes what an agent supports.
type AgentCapabilities struct {
	LoadSession      bool                    `json:"loadSession,omitempty"`
	PromptCapabilities *PromptCapabilities   `json:"promptCapabilities,omitempty"`
	MCPCapabilities  *MCPCapabilities        `json:"mcpCapabilities,omitempty"`
	SessionCapabilities *SessionCapabilities `json:"sessionCapabilities,omitempty"`
	Auth             *AgentAuthCapabilities  `json:"auth,omitempty"`
	Meta             map[string]any          `json:"_meta,omitempty"`
}

// PromptCapabilities describes prompt content support.
type PromptCapabilities struct {
	Image           bool `json:"image,omitempty"`
	Audio           bool `json:"audio,omitempty"`
	EmbeddedContext bool `json:"embeddedContext,omitempty"`
	Meta            map[string]any `json:"_meta,omitempty"`
}

// MCPCapabilities describes MCP transport support.
type MCPCapabilities struct {
	HTTP bool `json:"http,omitempty"`
	SSE  bool `json:"sse,omitempty"`
	Meta map[string]any `json:"_meta,omitempty"`
}

// SessionCapabilities describes session-related capabilities.
type SessionCapabilities struct {
	List                  interface{} `json:"list,omitempty"`
	Delete                interface{} `json:"delete,omitempty"`
	AdditionalDirectories interface{} `json:"additionalDirectories,omitempty"`
	Resume                interface{} `json:"resume,omitempty"`
	Close                 interface{} `json:"close,omitempty"`
	Meta                  map[string]any `json:"_meta,omitempty"`
}

// AgentAuthCapabilities describes authentication capabilities.
type AgentAuthCapabilities struct {
	Logout *LogoutCapabilities `json:"logout,omitempty"`
	Meta   map[string]any      `json:"_meta,omitempty"`
}

// LogoutCapabilities describes logout support.
type LogoutCapabilities struct {
	Meta map[string]any `json:"_meta,omitempty"`
}

// Implementation describes a client or agent implementation.
type Implementation struct {
	Name    string         `json:"name"`
	Title   string         `json:"title,omitempty"`
	Version string         `json:"version"`
	Meta    map[string]any `json:"_meta,omitempty"`
}

// AuthMethod describes an authentication method.
type AuthMethod struct {
	Type        string         `json:"type"`
	ID          string         `json:"id"`
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Meta        map[string]any `json:"_meta,omitempty"`
}

