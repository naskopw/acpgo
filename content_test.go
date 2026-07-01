package acp_test

import (
	"encoding/json"
	"testing"

	"github.com/naskopw/acpgo"
	"github.com/stretchr/testify/require"
)

func TestContentBlockResourceLinkTitle(t *testing.T) {
	cb := acp.ContentBlock{
		Type:  "resource_link",
		URI:   "file:///home/user/doc.pdf",
		Name:  "doc.pdf",
		Title: "My Document",
	}
	data, err := json.Marshal(cb)
	require.NoError(t, err)
	require.Contains(t, string(data), `"title":"My Document"`)

	var got acp.ContentBlock
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, "My Document", got.Title)
}
