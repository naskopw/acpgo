package acp_test

import (
	"encoding/json"
	"testing"

	acp "github.com/naskopw/acpgo"
	"github.com/stretchr/testify/require"
)

func TestPermissionOutcomeSelectedJSON(t *testing.T) {
	outcome := acp.PermissionOutcome{
		Outcome:  acp.PermOutcomeSelected,
		OptionID: "allow-once",
	}
	data, err := json.Marshal(outcome)
	require.NoError(t, err)

	expected := `{"outcome":"selected","optionId":"allow-once"}`
	require.JSONEq(t, expected, string(data))
}

func TestPermissionOutcomeCancelledJSON(t *testing.T) {
	outcome := acp.PermissionOutcome{
		Outcome: acp.PermOutcomeCancelled,
	}
	data, err := json.Marshal(outcome)
	require.NoError(t, err)

	expected := `{"outcome":"cancelled"}`
	require.JSONEq(t, expected, string(data))
}

func TestRequestPermissionResponseJSON(t *testing.T) {
	resp := acp.RequestPermissionResponse{
		Outcome: &acp.PermissionOutcome{
			Outcome:  acp.PermOutcomeSelected,
			OptionID: "allow-once",
		},
	}
	data, err := json.Marshal(resp)
	require.NoError(t, err)
	require.Contains(t, string(data), `"outcome"`)
	require.Contains(t, string(data), `"optionId"`)
}
