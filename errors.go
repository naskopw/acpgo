package acp

import "errors"

var (
	// ErrHarnessNotFound is returned when the harness binary cannot be found.
	ErrHarnessNotFound = errors.New("harness not found")
	// ErrProtocolMismatch is returned when the protocol versions are incompatible.
	ErrProtocolMismatch = errors.New("protocol version mismatch")
	// ErrRequestCancelled is returned when a request was cancelled.
	ErrRequestCancelled = errors.New("request cancelled")
	// ErrInvalidCWD is returned when a CWD path is set but is not absolute.
	ErrInvalidCWD = errors.New("cwd must be an absolute path")
	// ErrInvalidAdditionalDirectory is returned when an additionalDirectory entry is not absolute.
	ErrInvalidAdditionalDirectory = errors.New("additionalDirectories entries must be absolute paths")
	// ErrInvalidMCPCommand is returned when an MCP server command is not an absolute path.
	ErrInvalidMCPCommand = errors.New("mcpServers command must be an absolute path")
)
