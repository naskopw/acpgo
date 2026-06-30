package acp_test

import (
	"encoding/json"
	"testing"

	"github.com/naskopw/acp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChunkTypeString(t *testing.T) {
	tests := []struct {
		ct   acp.ChunkType
		want string
	}{
		{acp.ChunkText, "text"},
		{acp.ChunkThought, "thought"},
		{acp.ChunkToolUse, "tool_use"},
		{acp.ChunkDone, "done"},
		{acp.ChunkError, "error"},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.ct.String(), "ChunkType(%d).String()", int(tt.ct))
	}
}

func TestChunkTypeMarshalJSON(t *testing.T) {
	tests := []struct {
		ct   acp.ChunkType
		want string
	}{
		{acp.ChunkText, `"text"`},
		{acp.ChunkThought, `"thought"`},
		{acp.ChunkToolUse, `"tool_use"`},
		{acp.ChunkDone, `"done"`},
		{acp.ChunkError, `"error"`},
	}
	for _, tt := range tests {
		data, err := json.Marshal(tt.ct)
		require.NoError(t, err)
		assert.Equal(t, tt.want, string(data), "json.Marshal(%d)", int(tt.ct))
	}
}

func TestChunkTypeUnmarshalJSON(t *testing.T) {
	tests := []struct {
		data string
		want acp.ChunkType
	}{
		{`"text"`, acp.ChunkText},
		{`"thought"`, acp.ChunkThought},
		{`"tool_use"`, acp.ChunkToolUse},
		{`"done"`, acp.ChunkDone},
		{`"error"`, acp.ChunkError},
	}
	for _, tt := range tests {
		var ct acp.ChunkType
		require.NoError(t, json.Unmarshal([]byte(tt.data), &ct))
		assert.Equal(t, tt.want, ct, "Unmarshal(%s)", tt.data)
	}
}

func TestChunkTypeRoundtrip(t *testing.T) {
	for _, ct := range []acp.ChunkType{acp.ChunkText, acp.ChunkThought, acp.ChunkToolUse, acp.ChunkDone, acp.ChunkError} {
		data, err := json.Marshal(ct)
		require.NoError(t, err)
		var got acp.ChunkType
		require.NoError(t, json.Unmarshal(data, &got))
		require.Equal(t, ct, got)
	}
}

func TestChunkTypeUnmarshalInvalid(t *testing.T) {
	var ct acp.ChunkType
	require.Error(t, json.Unmarshal([]byte(`"invalid"`), &ct))
}

func TestParseChunkType(t *testing.T) {
	ct, err := acp.ParseChunkType("thought")
	require.NoError(t, err)
	require.Equal(t, acp.ChunkThought, ct)
}

func TestParseChunkTypeInvalid(t *testing.T) {
	_, err := acp.ParseChunkType("invalid")
	require.Error(t, err)
}

func TestChunkTypeMarshalUnknown(t *testing.T) {
	unknown := acp.ChunkType(99)
	_, err := json.Marshal(unknown)
	require.Error(t, err)
}
