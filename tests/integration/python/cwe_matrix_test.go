package python_test

import (
	"fmt"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/tests/integration"
)

// TestPythonCWEFixturesMatrix validates priority CWE fixtures under
// tests/fixtures/python/cwe with a Python-only registry.
// Kept separate from Go CWE matrix tests in tests/integration/cwe_matrix_test.go.
func TestPythonCWEFixturesMatrix(t *testing.T) {
	cases, err := integration.DiscoverPythonCWECases()
	if err != nil {
		t.Fatal(err)
	}
	if len(cases) == 0 {
		t.Fatal("expected Python CWE fixture pairs under tests/fixtures/python/cwe")
	}
	t.Logf("Python CWE fixture cases: %d (×2 files)", len(cases))

	var failures []string
	for _, c := range cases {
		// Stem CWE-78 → CWE-78
		rule := c
		if len(c) > 4 && c[:4] == "CWE-" {
			// already CWE-N
			rule = c
		}
		opts := integration.DefaultMatrixOptions()
		opts.Only = []string{rule}
		opts.Languages = []core.LanguageID{core.LanguagePython}
		// CWE-502 is fixture-only maturity under recommended packs; pin only.
		opts.Profile = core.ProfileAll

		vulnRel := integration.PythonCWEFixtureRel(c, true)
		vres, err := integration.MaterializeAndScanOpts(vulnRel, opts)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s vuln: %v", c, err))
		} else if aerr := integration.AssertVulnerable(vres.Findings, rule, vulnRel); aerr != nil {
			failures = append(failures, aerr.Error())
		}

		safeRel := integration.PythonCWEFixtureRel(c, false)
		sres, err := integration.MaterializeAndScanOpts(safeRel, opts)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s safe: %v", c, err))
		} else if integration.HasRule(sres.Findings, rule) {
			failures = append(failures, fmt.Sprintf("%s: expected no %s on safe, got %s",
				safeRel, rule, integration.SummarizeFindings(sres.Findings)))
		}
	}

	if len(failures) > 0 {
		t.Fatalf("Python CWE matrix failures %d of %d:\n%s",
			len(failures), len(cases), formatFailures(failures, 40))
	}
}
