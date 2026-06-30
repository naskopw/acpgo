package acp_test

import (
	"encoding/json"
	"testing"

	"github.com/naskopw/acp"
	"github.com/stretchr/testify/require"
)

func TestRequestJSON(t *testing.T) {
	req := acp.Request{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`1`),
		Method:  "initialize",
	}
	data, err := json.Marshal(req)
	require.NoError(t, err)
	var got acp.Request
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, "2.0", got.JSONRPC)
	require.Equal(t, "initialize", got.Method)
}

func TestRequestWithParams(t *testing.T) {
	req := acp.Request{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`1`),
		Method:  "setModel",
		Params:  json.RawMessage(`{"model_id":"gpt-4"}`),
	}
	data, err := json.Marshal(req)
	require.NoError(t, err)
	var got acp.Request
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, `{"model_id":"gpt-4"}`, string(got.Params))
}

func TestResponseJSON(t *testing.T) {
	resp := acp.Response{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`1`),
		Result:  json.RawMessage(`"ok"`),
	}
	data, err := json.Marshal(resp)
	require.NoError(t, err)
	var got acp.Response
	require.NoError(t, json.Unmarshal(data, &got))
	require.Nil(t, got.Error)
}

func TestResponseError(t *testing.T) {
	resp := acp.Response{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`1`),
		Error: &acp.RPCError{
			Code:    acp.ErrCodeMethodNotFound,
			Message: "method not found",
		},
	}
	data, err := json.Marshal(resp)
	require.NoError(t, err)
	var got acp.Response
	require.NoError(t, json.Unmarshal(data, &got))
	require.NotNil(t, got.Error)
	require.Equal(t, acp.ErrCodeMethodNotFound, got.Error.Code)
}

func TestNotificationJSON(t *testing.T) {
	n := acp.Notification{
		JSONRPC: "2.0",
		Method:  "session/update",
		Params:  json.RawMessage(`{"type":"text"}`),
	}
	data, err := json.Marshal(n)
	require.NoError(t, err)
	var got acp.Notification
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, "session/update", got.Method)
}

func TestRPCErrorJSON(t *testing.T) {
	e := acp.RPCError{
		Code:    acp.ErrCodeParse,
		Message: "parse error",
	}
	data, err := json.Marshal(e)
	require.NoError(t, err)
	var got acp.RPCError
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, acp.ErrCodeParse, got.Code)
	require.Equal(t, "parse error", got.Message)
}

func TestErrorCodeConstants(t *testing.T) {
	require.True(t, acp.ErrCodeParse < 0, "ErrCodeParse should be negative")
	require.True(t, acp.ErrCodeInvalidRequest < 0, "ErrCodeInvalidRequest should be negative")
	require.True(t, acp.ErrCodeMethodNotFound < 0, "ErrCodeMethodNotFound should be negative")
	require.True(t, acp.ErrCodeInvalidParams < 0, "ErrCodeInvalidParams should be negative")
	require.True(t, acp.ErrCodeInternal < 0, "ErrCodeInternal should be negative")
	require.True(t, acp.ErrCodeServer < 0, "ErrCodeServer should be negative")
}
