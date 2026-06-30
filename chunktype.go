package acp

import (
	"encoding/json"
	"fmt"
)

var chunkTypeStrings = map[ChunkType]string{
	ChunkText:    "text",
	ChunkThought: "thought",
	ChunkToolUse: "tool_use",
	ChunkDone:    "done",
	ChunkError:   "error",
}

var chunkTypeValues = map[string]ChunkType{
	"text":     ChunkText,
	"thought":  ChunkThought,
	"tool_use": ChunkToolUse,
	"done":     ChunkDone,
	"error":    ChunkError,
}

func (ct ChunkType) String() string {
	if s, ok := chunkTypeStrings[ct]; ok {
		return s
	}
	return fmt.Sprintf("ChunkType(%d)", int(ct))
}

// MarshalJSON implements json.Marshaler for ChunkType.
func (ct ChunkType) MarshalJSON() ([]byte, error) {
	s, ok := chunkTypeStrings[ct]
	if !ok {
		return nil, fmt.Errorf("unknown ChunkType: %d", int(ct))
	}
	return json.Marshal(s)
}

// UnmarshalJSON implements json.Unmarshaler for ChunkType.
func (ct *ChunkType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("unmarshal ChunkType: %w", err)
	}
	v, ok := chunkTypeValues[s]
	if !ok {
		return fmt.Errorf("unknown ChunkType: %q", s)
	}
	*ct = v
	return nil
}

// ParseChunkType parses a string into a ChunkType.
func ParseChunkType(s string) (ChunkType, error) {
	v, ok := chunkTypeValues[s]
	if !ok {
		return 0, fmt.Errorf("unknown ChunkType: %q", s)
	}
	return v, nil
}
