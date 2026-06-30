package acp

// ModelInfo describes an AI model available through a harness.
type ModelInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ChunkType identifies the kind of content in a streaming chunk.
type ChunkType int

// Chunk types for streaming session output.
const (
	ChunkText    ChunkType = iota
	ChunkThought
	ChunkToolUse
	ChunkDone
	ChunkError
)

// Message represents a single turn in a conversation.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// StructuredError carries detailed error information from the harness.
type StructuredError struct {
	Source   string `json:"source"`
	Message  string `json:"message"`
	Code     *int32 `json:"code,omitempty"`
	DataJSON string `json:"data_json,omitempty"`
}
