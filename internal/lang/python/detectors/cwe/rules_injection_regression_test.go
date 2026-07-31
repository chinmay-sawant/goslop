package cwe

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestDetectCWE93PreservesMultilineSourceAlignment(t *testing.T) {
	source := `"""A multiline docstring.
The masked source must keep this line.
"""
response.headers["Location"] = next_url
`
	unit := core.NewParsedUnit(core.LanguagePython, "app.py", source)
	var findings []rules.Finding

	detectCWE93(unit, nil, &findings)

	if len(findings) != 1 {
		t.Fatalf("findings = %d, want 1: %+v", len(findings), findings)
	}
	if findings[0].Line != 4 {
		t.Fatalf("finding line = %d, want 4", findings[0].Line)
	}
}
