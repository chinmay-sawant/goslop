package cwe_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/tests/integration"
)

// assertCWEFixturePair makes the text fixtures the unit-test corpus as well
// as the integration-test corpus. Keep detector test files to rule IDs and
// use fixture edits for all Python source examples and regressions.
func assertCWEFixturePair(t *testing.T, rule string) {
	t.Helper()
	opts := integration.DefaultMatrixOptions()
	opts.Only = []string{rule}
	opts.Languages = []core.LanguageID{core.LanguagePython}
	opts.Profile = core.ProfileAll

	vulnerable := integration.PythonCWEFixtureRel(rule, true)
	vulnResult, err := integration.MaterializeAndScanOpts(vulnerable, opts)
	if err != nil {
		t.Fatalf("scan vulnerable fixture %s: %v", vulnerable, err)
	}
	if err := integration.AssertVulnerable(vulnResult.Findings, rule, vulnerable); err != nil {
		t.Fatal(err)
	}

	safe := integration.PythonCWEFixtureRel(rule, false)
	safeResult, err := integration.MaterializeAndScanOpts(safe, opts)
	if err != nil {
		t.Fatalf("scan safe fixture %s: %v", safe, err)
	}
	if integration.HasRule(safeResult.Findings, rule) {
		t.Fatalf("safe fixture %s unexpectedly emitted %s: %s", safe, rule, integration.SummarizeFindings(safeResult.Findings))
	}
}
