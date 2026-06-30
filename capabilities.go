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
	Meta        map[string]any `json:"_meta,omitempty"`
}

// AuthMethodID is an auth method identifier.
type AuthMethodID string

// Capability is a legacy boolean-capability struct for backward compatibility.
type Capability struct {
	Models        bool `json:"models"`
	Sessions      bool `json:"sessions"`
	History       bool `json:"history"`
	SlashCommands bool `json:"slashCommands"`
	Cancel        bool `json:"cancel"`
}

// ParseCapabilities parses a capability string slice into a Capability struct.
func ParseCapabilities(caps []string) Capability {
	var c Capability
	for _, s := range caps {
		switch s {
		case "models": c.Models = true
		case "sessions": c.Sessions = true
		case "history": c.History = true
		case "slashCommands": c.SlashCommands = true
		case "cancel": c.Cancel = true
		}
	}
	return c
}

// List returns the capability strings for enabled capabilities.
func (c Capability) List() []string {
	var caps []string
	if c.Models { caps = append(caps, "models") }
	if c.Sessions { caps = append(caps, "sessions") }
	if c.History { caps = append(caps, "history") }
	if c.SlashCommands { caps = append(caps, "slashCommands") }
	if c.Cancel { caps = append(caps, "cancel") }
	return caps
}
