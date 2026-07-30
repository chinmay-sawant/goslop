package core

import (
	"fmt"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/rules"
)

// FailPolicy controls when the process should exit non-zero based on finding severity.
type FailPolicy int

const (
	// FailNone always exits 0 for findings (Rust: NoFail).
	FailNone FailPolicy = iota
	// FailHigh fails only on high or critical (Rust: Strict).
	FailHigh
	// FailMedium fails on medium and above (Rust: MediumAsErrors).
	FailMedium
)

// FailNever is an alias for FailNone.
const FailNever = FailNone

// String returns a stable lowercase policy name.
func (p FailPolicy) String() string {
	switch p {
	case FailNone:
		return "none"
	case FailHigh:
		return "high"
	case FailMedium:
		return "medium"
	default:
		return fmt.Sprintf("FailPolicy(%d)", int(p))
	}
}

// ParseFailPolicy parses a fail policy name (aliases accepted).
func ParseFailPolicy(s string) (FailPolicy, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "none", "no_fail", "nofail", "never":
		return FailNone, nil
	case "high", "strict":
		return FailHigh, nil
	case "medium", "medium_as_errors", "default":
		return FailMedium, nil
	default:
		return 0, fmt.Errorf("unknown fail policy %q", s)
	}
}

// ShouldFail reports whether a finding of severity should fail the run.
func (p FailPolicy) ShouldFail(severity rules.Severity) bool {
	switch p {
	case FailNone:
		return false
	case FailHigh:
		return severity >= rules.SeverityHigh
	case FailMedium:
		return severity.IsFailure()
	default:
		return false
	}
}
