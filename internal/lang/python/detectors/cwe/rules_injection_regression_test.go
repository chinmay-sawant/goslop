package cwe

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/fixture"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestDetectCWE93PreservesMultilineSourceAlignment(t *testing.T) {
	fixturePath := filepath.Join("..", "..", "..", "..", "..", "tests", "fixtures", "python", "cwe", "CWE-93-multiline-docstring-vulnerable.txt")
	contents, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("read %s: %v", fixturePath, err)
	}
	fx, err := fixture.ParseFixture(string(contents), filepath.Base(fixturePath))
	if err != nil {
		t.Fatalf("parse %s: %v", fixturePath, err)
	}
	path := fx.Filename
	if path == "" {
		path = "app.py"
	}
	unit := core.NewParsedUnit(core.LanguagePython, path, fx.Source)
	var findings []rules.Finding

	detectCWE93(unit, nil, &findings)

	if len(findings) != 1 {
		t.Fatalf("findings = %d, want 1: %+v", len(findings), findings)
	}
	if findings[0].Line != 4 {
		t.Fatalf("finding line = %d, want 4", findings[0].Line)
	}
}
