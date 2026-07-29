package perf_test

import (
	"testing"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/lang/go/detectors/perf"
	"github.com/chinmay/codehound/internal/rules"
)

func TestPERF116Vulnerable(t *testing.T) {
	src := `package sample

import "strings"

func Has(s, sub string) bool {
	return strings.Index(s, sub) != -1
}
`
	findings := runPERF116(t, src, "PERF-116-vulnerable.go")
	if !hasRule(findings, "PERF-116") {
		t.Fatalf("expected PERF-116 finding, got %#v", findings)
	}
}

func TestPERF116Safe(t *testing.T) {
	src := `package sample

import "strings"

func Has(s, sub string) bool {
	return strings.Contains(s, sub)
}
`
	findings := runPERF116(t, src, "PERF-116-safe.go")
	if hasRule(findings, "PERF-116") {
		t.Fatalf("safe fixture should not emit PERF-116, got %#v", findings)
	}
}

func TestPERF116IndexWithoutCompare(t *testing.T) {
	src := `package sample
import "strings"
func Pos(s, sub string) int { return strings.Index(s, sub) }
`
	findings := runPERF116(t, src, "index-only.go")
	if hasRule(findings, "PERF-116") {
		t.Fatalf("Index without -1 compare should be silent, got %#v", findings)
	}
}

func runPERF116(t *testing.T, src, path string) []rules.Finding {
	t.Helper()
	d := perf.NewPERF116()
	unit := core.NewParsedUnit(core.LangGo, path, src)
	ctx := core.DefaultScanContext()
	var out []rules.Finding
	d.Run(ctx, unit, &out)
	return out
}

func hasRule(findings []rules.Finding, id string) bool {
	for _, f := range findings {
		if f.RuleID == id {
			return true
		}
	}
	return false
}
