// Package python_test holds Python-language fixture matrix integration tests.
//
// These are intentionally separate from Go matrix tests under tests/integration/
// (bp_matrix_test.go, cwe_matrix_test.go, perf_matrix_test.go) so:
//   - `go test ./tests/integration` remains the Go-only parity surface
//   - `go test ./tests/integration/python` validates Python detectors
//   - registries differ (DefaultRegistry is Go-only; Python needs LanguagePython)
package python_test

import (
	"fmt"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/tests/integration"
)

// TestPythonBPFixturesMatrix scans every tests/fixtures/python/bp/*-{vulnerable,safe}.txt
// pair with languages=[python] and --only BP-PY-N.
func TestPythonBPFixturesMatrix(t *testing.T) {
	cases, err := integration.DiscoverPythonBPCases()
	if err != nil {
		t.Fatal(err)
	}
	if len(cases) == 0 {
		t.Fatal("expected at least one Python BP-PY fixture case under tests/fixtures/python/bp")
	}
	if len(cases) < 151 {
		t.Fatalf("expected at least 151 BP-PY cases, got %d", len(cases))
	}
	t.Logf("Python BP-PY fixture cases: %d (×2 files)", len(cases))

	var failures []string
	for _, c := range cases {
		rule := integration.PythonBPRuleID(c)
		opts := integration.DefaultMatrixOptions()
		opts.Only = []string{rule}
		opts.IncludeTests = true // BP-PY-41/42 are test-file oriented
		opts.Languages = []core.LanguageID{core.LanguagePython}

		vulnRel := integration.PythonBPFixtureRel(c, true)
		vres, err := integration.MaterializeAndScanOpts(vulnRel, opts)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s vuln: %v", c, err))
		} else if aerr := integration.AssertVulnerable(vres.Findings, rule, vulnRel); aerr != nil {
			failures = append(failures, aerr.Error())
		}

		safeRel := integration.PythonBPFixtureRel(c, false)
		sres, err := integration.MaterializeAndScanOpts(safeRel, opts)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s safe: %v", c, err))
		} else if integration.HasRule(sres.Findings, rule) {
			failures = append(failures, fmt.Sprintf("%s: expected no %s on safe, got %s",
				safeRel, rule, integration.SummarizeFindings(sres.Findings)))
		}
	}

	if len(failures) > 0 {
		t.Fatalf("Python BP matrix failures %d of %d cases:\n%s",
			len(failures), len(cases), formatFailures(failures, 60))
	}
}

// TestPythonBPFixtureInventorySorted ensures paired discovery is stable.
func TestPythonBPFixtureInventorySorted(t *testing.T) {
	cases, err := integration.DiscoverPythonBPCases()
	if err != nil {
		t.Fatal(err)
	}
	if len(cases) != 151 {
		t.Fatalf("inventory size = %d, want 151", len(cases))
	}
	seen := map[string]struct{}{}
	for _, c := range cases {
		if _, ok := seen[c]; ok {
			t.Fatalf("duplicate case %q", c)
		}
		seen[c] = struct{}{}
		if integration.PythonBPRuleID(c) == "" {
			t.Fatalf("empty rule id for %q", c)
		}
	}
}

func formatFailures(failures []string, limit int) string {
	if len(failures) <= limit {
		out := ""
		for _, f := range failures {
			out += "  - " + f + "\n"
		}
		return out
	}
	out := ""
	for _, f := range failures[:limit] {
		out += "  - " + f + "\n"
	}
	out += fmt.Sprintf("  ... +%d more\n", len(failures)-limit)
	return out
}
