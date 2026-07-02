package acp

import "path/filepath"

func validateEnum[T comparable](val T, enumErr error, allowed ...T) error {
	for _, a := range allowed {
		if val == a {
			return nil
		}
	}
	return enumErr
}

// Validate checks that all path fields are absolute per ACP spec.
func (r NewSessionRequest) Validate() error {
	if r.CWD != "" && !filepath.IsAbs(r.CWD) {
		return ErrInvalidCWD
	}
	for _, d := range r.AdditionalDirectories {
		if !filepath.IsAbs(d) {
			return ErrInvalidAdditionalDirectory
		}
	}
	for _, s := range r.MCPServers {
		if err := s.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate checks that all path fields are absolute per ACP spec.
func (r LoadSessionRequest) Validate() error {
	if r.CWD != "" && !filepath.IsAbs(r.CWD) {
		return ErrInvalidCWD
	}
	for _, d := range r.AdditionalDirectories {
		if !filepath.IsAbs(d) {
			return ErrInvalidAdditionalDirectory
		}
	}
	for _, s := range r.MCPServers {
		if err := s.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate checks that all path fields are absolute per ACP spec.
func (r ResumeSessionRequest) Validate() error {
	if r.CWD != "" && !filepath.IsAbs(r.CWD) {
		return ErrInvalidCWD
	}
	for _, d := range r.AdditionalDirectories {
		if !filepath.IsAbs(d) {
			return ErrInvalidAdditionalDirectory
		}
	}
	for _, s := range r.MCPServers {
		if err := s.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate checks that the MCP server command (if set) is an absolute path.
func (s MCPServer) Validate() error {
	if s.Type != "" {
		if err := validateEnum(s.Type, ErrInvalidMCPServerType, "stdio", "http", "sse"); err != nil {
			return err
		}
	}
	if s.Type == "stdio" && s.Command != "" && !filepath.IsAbs(s.Command) {
		return ErrInvalidMCPCommand
	}
	return nil
}

// Validate checks that cwd (if set) is an absolute path.
func (r ListSessionsRequest) Validate() error {
	if r.CWD != "" && !filepath.IsAbs(r.CWD) {
		return ErrInvalidCWD
	}
	return nil
}

// Validate checks that the ContentBlock type is valid.
func (b ContentBlock) Validate() error {
	return validateEnum(b.Type, ErrInvalidContentType, ContentTypeText, ContentTypeImage, ContentTypeAudio, ContentTypeResourceLink, ContentTypeResource)
}

// Validate checks that the ToolCall kind and status are valid (if set).
func (t ToolCall) Validate() error {
	if t.Kind != "" {
		if err := validateEnum(t.Kind, ErrInvalidToolKind, ToolKindRead, ToolKindEdit, ToolKindDelete, ToolKindMove, ToolKindSearch, ToolKindExecute, ToolKindThink, ToolKindFetch, ToolKindSwitchMode, ToolKindOther); err != nil {
			return err
		}
	}
	if t.Status != "" {
		if err := validateEnum(t.Status, ErrInvalidToolStatus, ToolStatusPending, ToolStatusInProgress, ToolStatusCompleted, ToolStatusFailed); err != nil {
			return err
		}
	}
	return nil
}

// Validate checks that the ToolCallUpdate kind and status are valid (if set).
func (t ToolCallUpdate) Validate() error {
	if t.Kind != "" {
		if err := validateEnum(t.Kind, ErrInvalidToolKind, ToolKindRead, ToolKindEdit, ToolKindDelete, ToolKindMove, ToolKindSearch, ToolKindExecute, ToolKindThink, ToolKindFetch, ToolKindSwitchMode, ToolKindOther); err != nil {
			return err
		}
	}
	if t.Status != "" {
		if err := validateEnum(t.Status, ErrInvalidToolStatus, ToolStatusPending, ToolStatusInProgress, ToolStatusCompleted, ToolStatusFailed); err != nil {
			return err
		}
	}
	return nil
}

// Validate checks that the ToolCallContent type is valid.
func (t ToolCallContent) Validate() error {
	return validateEnum(t.Type, ErrInvalidToolContentType, ToolContentTypeContent, ToolContentTypeDiff, ToolContentTypeTerminal)
}

// Validate checks that the PlanEntry priority and status are valid.
func (e PlanEntry) Validate() error {
	if err := validateEnum(e.Priority, ErrInvalidPlanPriority, PlanPriorityHigh, PlanPriorityMedium, PlanPriorityLow); err != nil {
		return err
	}
	if err := validateEnum(e.Status, ErrInvalidPlanStatus, PlanStatusPending, PlanStatusInProgress, PlanStatusCompleted); err != nil {
		return err
	}
	return nil
}

// Validate checks that the PermissionOption kind is valid.
func (o PermissionOption) Validate() error {
	return validateEnum(o.Kind, ErrInvalidPermissionKind, PermOptionAllowOnce, PermOptionAllowAlways, PermOptionRejectOnce, PermOptionRejectAlways)
}

// Validate checks that the PermissionOutcome outcome is valid.
func (o PermissionOutcome) Validate() error {
	return validateEnum(o.Outcome, ErrInvalidPermissionOutcome, PermOutcomeCancelled, PermOutcomeSelected)
}

// Validate checks that the SessionUpdate variant is valid.
func (s SessionUpdate) Validate() error {
	return validateEnum(s.SessionUpdateVariant, ErrInvalidSessionUpdateVariant,
		SessionUpdateUserMessageChunk, SessionUpdateAgentMessageChunk, SessionUpdateAgentThoughtChunk,
		SessionUpdateToolCall, SessionUpdateToolCallUpdate, SessionUpdatePlan,
		SessionUpdateAvailableCommands, SessionUpdateCurrentMode, SessionUpdateConfigOption,
		SessionUpdateSessionInfo, SessionUpdateUsage)
}

// Validate checks that the PromptResponse stop reason is valid.
func (r PromptResponse) Validate() error {
	return validateEnum(r.StopReason, ErrInvalidStopReason, StopReasonEndTurn, StopReasonMaxTokens, StopReasonMaxTurnReqs, StopReasonRefusal, StopReasonCancelled)
}

// Validate checks that the ConfigOption category is valid (if set).
func (c ConfigOption) Validate() error {
	if c.Category != "" {
		return validateEnum(c.Category, ErrInvalidConfigCategory, ConfigCategoryMode, ConfigCategoryModel, ConfigCategoryModelConfig, ConfigCategoryThoughtLevel)
	}
	return nil
}

// Validate checks that each Annotations audience entry is valid.
func (a Annotations) Validate() error {
	for _, role := range a.Audience {
		if err := validateEnum(role, ErrInvalidAudience, RoleAssistant, RoleUser); err != nil {
			return err
		}
	}
	return nil
}
