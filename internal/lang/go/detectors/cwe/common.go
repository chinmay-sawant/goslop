package cwe

import "github.com/chinmay-sawant/goslop/internal/core"

// unitFile returns the display path for findings.
func unitFile(unit *core.ParsedUnit) string {
	if unit == nil {
		return ""
	}
	if unit.DisplayPath != "" {
		return unit.DisplayPath
	}
	return unit.Path
}
