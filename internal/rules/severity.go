// Package rules defines findings, severity, metadata, packs, and emit helpers.
package rules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Severity is a finding severity used by failure policy and reporters.
// JSON wire form is lowercase ("info", "low", …).
type Severity int

const (
	SeverityInfo Severity = iota
	SeverityLow
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

// String returns the stable lowercase wire representation.
func (s Severity) String() string {
	switch s {
	case SeverityInfo:
		return "info"
	case SeverityLow:
		return "low"
	case SeverityMedium:
		return "medium"
	case SeverityHigh:
		return "high"
	case SeverityCritical:
		return "critical"
	default:
		return fmt.Sprintf("Severity(%d)", int(s))
	}
}

// IsFailure reports whether this severity participates in the default failure gate.
func (s Severity) IsFailure() bool {
	return s >= SeverityMedium
}

// ParseSeverity parses a case-insensitive severity name.
func ParseSeverity(s string) (Severity, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "info":
		return SeverityInfo, nil
	case "low":
		return SeverityLow, nil
	case "medium":
		return SeverityMedium, nil
	case "high":
		return SeverityHigh, nil
	case "critical":
		return SeverityCritical, nil
	default:
		return 0, fmt.Errorf("unknown severity %q", s)
	}
}

// MarshalJSON encodes Severity as a lowercase JSON string.
func (s Severity) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON decodes a lowercase (or mixed-case) severity string.
func (s *Severity) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	parsed, err := ParseSeverity(raw)
	if err != nil {
		return err
	}
	*s = parsed
	return nil
}
