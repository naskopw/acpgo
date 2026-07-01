package acp

import "path/filepath"

// Validate checks that cwd (if set) is an absolute path.
func (r NewSessionRequest) Validate() error {
	if r.CWD != "" && !filepath.IsAbs(r.CWD) {
		return ErrInvalidCWD
	}
	return nil
}

// Validate checks that cwd (if set) is an absolute path.
func (r LoadSessionRequest) Validate() error {
	if r.CWD != "" && !filepath.IsAbs(r.CWD) {
		return ErrInvalidCWD
	}
	return nil
}

// Validate checks that cwd (if set) is an absolute path.
func (r ResumeSessionRequest) Validate() error {
	if r.CWD != "" && !filepath.IsAbs(r.CWD) {
		return ErrInvalidCWD
	}
	return nil
}
