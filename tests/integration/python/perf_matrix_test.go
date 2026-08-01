package python_test

import (
	"fmt"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/tests/integration"
)

// TestPythonPERFFixturesMatrix validates each PERF-PY vulnerable/safe pair with
// the Python plugin explicitly enabled. It never exercises Go PERF discovery.
func TestPythonPERFFixturesMatrix(t *testing.T) {
	cases, err := integration.DiscoverPythonPERFCases()
	if err != nil {
		t.Fatal(err)
	}
	if len(cases) != 30 {
		t.Fatalf("PERF-PY fixture cases = %d, want 30", len(cases))
	}
	var failures []string
	for _, c := range cases {
		rule := integration.PythonPERFRuleID(c)
		opts := integration.DefaultMatrixOptions()
		opts.Only = []string{rule}
		opts.Languages = []core.LanguageID{core.LanguagePython}
		opts.Profile = core.ProfileAll
		vulnRel := integration.PythonPERFFixtureRel(c, true)
		vres, err := integration.MaterializeAndScanOpts(vulnRel, opts)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s vulnerable: %v", c, err))
		} else if assertErr := integration.AssertVulnerable(vres.Findings, rule, vulnRel); assertErr != nil {
			failures = append(failures, assertErr.Error())
		}
		safeRel := integration.PythonPERFFixtureRel(c, false)
		sres, err := integration.MaterializeAndScanOpts(safeRel, opts)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s safe: %v", c, err))
		} else if integration.HasRule(sres.Findings, rule) {
			failures = append(failures, fmt.Sprintf("%s: expected no %s, got %s", safeRel, rule, integration.SummarizeFindings(sres.Findings)))
		}
	}
	if len(failures) > 0 {
		t.Fatalf("Python PERF matrix failures %d of %d:\n%s", len(failures), len(cases), formatFailures(failures, 40))
	}
}

func TestPythonPERFFixtureInventory(t *testing.T) {
	cases, err := integration.DiscoverPythonPERFCases()
	if err != nil {
		t.Fatal(err)
	}
	if len(cases) != 30 {
		t.Fatalf("inventory = %d, want 30", len(cases))
	}
	for i, c := range cases {
		want := fmt.Sprintf("PERF-PY-%d", i+1)
		if got := integration.PythonPERFRuleID(c); got != want {
			t.Fatalf("case %q maps to %q, want %q", c, got, want)
		}
	}
}
