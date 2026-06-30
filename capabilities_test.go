package acp_test

import (
	"encoding/json"
	"testing"

	"github.com/naskopw/acpgo"
	"github.com/stretchr/testify/require"
)

func TestAuthMethodJSON(t *testing.T) {
	m := acp.AuthMethod{
		Type:        "oauth2",
		ID:          "github",
		Name:        "GitHub OAuth",
		Description: "Sign in with GitHub",
	}
	data, err := json.Marshal(m)
	require.NoError(t, err)
	var got acp.AuthMethod
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, "oauth2", got.Type)
	require.Equal(t, "github", got.ID)
	require.Equal(t, "GitHub OAuth", got.Name)
}

func TestSetSessionModeResponseJSON(t *testing.T) {
	p := acp.SetSessionModeResponse{}
	data, err := json.Marshal(p)
	require.NoError(t, err)
	var got acp.SetSessionModeResponse
	require.NoError(t, json.Unmarshal(data, &got))
}
