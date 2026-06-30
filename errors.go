package acp

import "errors"

var (
	// ErrHarnessNotFound is returned when the harness binary cannot be found.
	ErrHarnessNotFound = errors.New("harness not found")
	// ErrProtocolMismatch is returned when the protocol versions are incompatible.
	ErrProtocolMismatch = errors.New("protocol version mismatch")
	// ErrRequestCancelled is returned when a request was cancelled.
	ErrRequestCancelled = errors.New("request cancelled")
)
