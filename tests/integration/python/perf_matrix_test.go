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
	if len(cases) < 50 {
		t.Fatalf("PERF-PY fixture cases = %d, want at least 50", len(cases))
	}
	t.Logf("Python PERF-PY fixture cases: %d (×2 files)", len(cases))
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
	if len(cases) < 50 {
		t.Fatalf("inventory = %d, want at least 50", len(cases))
	}
	seen := map[string]struct{}{}
	rules := map[string]struct{}{}
	for _, c := range cases {
		if _, ok := seen[c]; ok {
			t.Fatalf("duplicate case %q", c)
		}
		seen[c] = struct{}{}
		rule := integration.PythonPERFRuleID(c)
		if rule == "" {
			t.Fatalf("empty rule id for %q", c)
		}
		rules[rule] = struct{}{}
	}
	for i := 1; i <= 30; i++ {
		want := fmt.Sprintf("PERF-PY-%d", i)
		if _, ok := rules[want]; !ok {
			t.Fatalf("missing base fixture case for %s", want)
		}
	}
}
