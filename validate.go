package acp

import "path/filepath"

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
