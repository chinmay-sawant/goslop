package core

import (
	"testing"

	"github.com/chinmay/codehound/internal/rules"
)

func TestAllowsMatrix(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(*ScanContext)
		ruleID string
		want   bool
	}{
		{"default allows CWE", func(*ScanContext) {}, "CWE-89", true},
		{"default allows BP", func(*ScanContext) {}, "BP-1", true},
		{"BP denied when disabled", func(c *ScanContext) { c.BadPracticesEnabled = false }, "BP-1", false},
		{"BP disabled allows CWE", func(c *ScanContext) { c.BadPracticesEnabled = false }, "CWE-22", true},
		{"exact skip denies", func(c *ScanContext) { c.Skip = []string{"CWE-89"} }, "CWE-89", false},
		{"exact skip leaves others", func(c *ScanContext) { c.Skip = []string{"CWE-89"} }, "CWE-22", true},
		{"prefix skip denies", func(c *ScanContext) { c.Skip = []string{"CWE-*"} }, "CWE-22", false},
		{"prefix skip no overmatch", func(c *ScanContext) { c.Skip = []string{"CWE-8*"} }, "CWE-22", true},
		{"only exact allows", func(c *ScanContext) { c.Only = []string{"PERF-101"} }, "PERF-101", true},
		{"only exact denies others", func(c *ScanContext) { c.Only = []string{"PERF-101"} }, "PERF-102", false},
		{"only prefix allows family", func(c *ScanContext) { c.Only = []string{"PERF-*"} }, "PERF-50", true},
		{"skip wins over only", func(c *ScanContext) { c.Only = []string{"CWE-*"}; c.Skip = []string{"CWE-89"} }, "CWE-89", false},
		{"BP disabled wins before only", func(c *ScanContext) { c.BadPracticesEnabled = false; c.Only = []string{"BP-*"} }, "BP-7", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := DefaultScanContext()
			tc.setup(ctx)
			if got := ctx.Allows(tc.ruleID); got != tc.want {
				t.Fatalf("Allows(%q)=%v want %v", tc.ruleID, got, tc.want)
			}
		})
	}
	var nilCtx *ScanContext
	if !nilCtx.Allows("CWE-1") {
		t.Fatal("nil allows")
	}
}

func TestFailPolicyShouldFail(t *testing.T) {
	if FailNone.ShouldFail(rules.SeverityCritical) || FailNever.ShouldFail(rules.SeverityHigh) {
		t.Fatal("none should not fail")
	}
	if !FailHigh.ShouldFail(rules.SeverityHigh) || FailHigh.ShouldFail(rules.SeverityMedium) {
		t.Fatal("high policy")
	}
	if !FailMedium.ShouldFail(rules.SeverityMedium) {
		t.Fatal("medium policy")
	}
}

func TestScanProfile(t *testing.T) {
	p, err := ParseScanProfile("ci")
	if err != nil || p != ProfileRecommended {
		t.Fatalf("ci: %v %v", p, err)
	}
	p2, ok := ParseProfile("bp")
	if !ok || p2 != ProfileStyle {
		t.Fatalf("bp: %v %v", p2, ok)
	}
	if ProfileRecommended.DefaultFailPolicy() != FailHigh {
		t.Fatal("fail")
	}
	if !ProfileSecurity.EnablesTaint() || ProfileRecommended.EnablesBadPractices() {
		t.Fatal("taint/bp")
	}
	patterns := ProfileRecommended.OnlyPatterns()
	hasCWE22, hasPERF101 := false, false
	for _, id := range patterns {
		if id == "CWE-22" {
			hasCWE22 = true
		}
		if id == "PERF-101" {
			hasPERF101 = true
		}
		if rules.PackFromRuleID(id).IsBadPractice() {
			t.Fatalf("BP in recommended: %s", id)
		}
	}
	if !hasCWE22 || !hasPERF101 {
		t.Fatal("missing core ids")
	}
	ctx := NewScanContext(ProfileRecommended, nil, nil)
	if ctx.FailPolicy != FailHigh || ctx.BadPracticesEnabled || !ctx.Allows("CWE-89") || ctx.Allows("BP-1") {
		t.Fatalf("applied ctx: %+v", ctx)
	}
}

func TestParsedUnitLineCol(t *testing.T) {
	src := "package main\n\nfunc f() {}\n"
	u := NewParsedUnit(LanguageGo, "main.go", src)
	if line, col := u.LineCol(0); line != 1 || col != 1 {
		t.Fatalf("0 → %d:%d", line, col)
	}
	if line, col := u.LineCol(13); line != 2 || col != 1 {
		t.Fatalf("13 → %d:%d", line, col)
	}
	if line, col := u.LineCol(14); line != 3 || col != 1 {
		t.Fatalf("14 → %d:%d", line, col)
	}
}
