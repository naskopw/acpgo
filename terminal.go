package acp

// CreateTerminalRequest creates a new terminal to execute a command.
type CreateTerminalRequest struct {
	SessionID      string         `json:"sessionId"`
	Command        string         `json:"command"`
	Args           []string       `json:"args,omitempty"`
	Env            []EnvVariable  `json:"env,omitempty"`
	CWD            string         `json:"cwd,omitempty"`
	OutputByteLimit int64         `json:"outputByteLimit,omitempty"`
	Meta           map[string]any `json:"_meta,omitempty"`
}

// CreateTerminalResponse contains the ID of the created terminal.
type CreateTerminalResponse struct {
	TerminalID string         `json:"terminalId"`
	Meta       map[string]any `json:"_meta,omitempty"`
}

// TerminalOutputRequest gets terminal output.
type TerminalOutputRequest struct {
	SessionID  string         `json:"sessionId"`
	TerminalID string         `json:"terminalId"`
	Meta       map[string]any `json:"_meta,omitempty"`
}

// TerminalOutputResponse contains terminal output and exit status.
type TerminalOutputResponse struct {
	Output     string             `json:"output"`
	Truncated  bool               `json:"truncated"`
	ExitStatus *TerminalExitStatus `json:"exitStatus,omitempty"`
	Meta       map[string]any     `json:"_meta,omitempty"`
}

// TerminalExitStatus describes how a terminal command exited.
type TerminalExitStatus struct {
	ExitCode int            `json:"exitCode,omitempty"`
	Signal   string         `json:"signal,omitempty"`
	Meta     map[string]any `json:"_meta,omitempty"`
}

// ReleaseTerminalRequest releases a terminal.
type ReleaseTerminalRequest struct {
	SessionID  string         `json:"sessionId"`
	TerminalID string         `json:"terminalId"`
	Meta       map[string]any `json:"_meta,omitempty"`
}

type ReleaseTerminalResponse struct {
	Meta map[string]any `json:"_meta,omitempty"`
}

// WaitForTerminalExitRequest waits for a terminal command to exit.
type WaitForTerminalExitRequest struct {
	SessionID  string         `json:"sessionId"`
	TerminalID string         `json:"terminalId"`
	Meta       map[string]any `json:"_meta,omitempty"`
}

// WaitForTerminalExitResponse contains the exit status.
type WaitForTerminalExitResponse struct {
	ExitCode int            `json:"exitCode,omitempty"`
	Signal   string         `json:"signal,omitempty"`
	Meta     map[string]any `json:"_meta,omitempty"`
}

// KillTerminalRequest kills a terminal command without releasing it.
type KillTerminalRequest struct {
	SessionID  string         `json:"sessionId"`
	TerminalID string         `json:"terminalId"`
	Meta       map[string]any `json:"_meta,omitempty"`
}

type KillTerminalResponse struct {
	Meta map[string]any `json:"_meta,omitempty"`
}
