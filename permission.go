package acp

// PermissionOption is an option presented to the user for a permission request.
type PermissionOption struct {
	OptionID string         `json:"optionId"`
	Name     string         `json:"name"`
	Kind     string         `json:"kind"`
	Meta     map[string]any `json:"_meta,omitempty"`
}

// RequestPermissionResponse is the result of a permission request.
type RequestPermissionResponse struct {
	Outcome *PermissionOutcome `json:"outcome"`
	Meta    map[string]any     `json:"_meta,omitempty"`
}

// PermissionOutcome is the user's decision on a permission request.
type PermissionOutcome struct {
	Outcome  string         `json:"outcome"`
	OptionID string         `json:"optionId,omitempty"`
	Meta     map[string]any `json:"_meta,omitempty"`
}
