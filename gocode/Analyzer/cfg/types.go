package cfg

import (
	"fmt"
)

type (
	// Config is the main configuration structure.
	SemanticVersion struct {
		Major    int
		Minor    int
		Revision int
	}

	Config struct {
		Verbosity      bool
		AccessPatterns string
		ToolName       string
		Version        SemanticVersion
	}
)

func (s *SemanticVersion) init(major, minor, revision int) {
	s.Major = major
	s.Minor = minor
	s.Revision = revision
}

func (s SemanticVersion) String() string {
	return fmt.Sprintf("v%v.%v.%v", s.Major, s.Minor, s.Revision)
}
