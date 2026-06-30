package acp

// ConfigOption is a session configuration option (model, mode, etc.).
type ConfigOption struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description,omitempty"`
	Category     string                 `json:"category,omitempty"`
	Type         string                 `json:"type"`
	CurrentValue string                 `json:"currentValue"`
	Options      []ConfigOptionValue    `json:"options"`
	Meta         map[string]any         `json:"_meta,omitempty"`
}

// ConfigOptionValue is a possible value for a config option.
type ConfigOptionValue struct {
	Value       string         `json:"value"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Meta        map[string]any `json:"_meta,omitempty"`
}
