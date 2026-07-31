package badpractices_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/tests/integration"
)

func TestBPFalsePositiveAuditFixtureVariants(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		rule     string
		caseName string
	}{
		{rule: "BP-PY-1", caseName: "BP-PY-1-thread-collection"},
		{rule: "BP-PY-26", caseName: "BP-PY-26-read-only"},
		{rule: "BP-PY-38", caseName: "BP-PY-38-task-list"},
		{rule: "BP-PY-42", caseName: "BP-PY-42-thread-collection"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-cli-output"},
	} {
		tc := tc
		t.Run(tc.caseName, func(t *testing.T) {
			t.Parallel()
			assertBPFixtureCase(t, tc.rule, tc.caseName)
		})
	}
}

func assertBPFixtureCase(t *testing.T, rule, fixtureCase string) {
	t.Helper()
	opts := integration.DefaultMatrixOptions()
	opts.Only = []string{rule}
	opts.IncludeTests = true
	opts.Languages = []core.LanguageID{core.LanguagePython}

	vulnerable := integration.PythonBPFixtureRel(fixtureCase, true)
	vulnResult, err := integration.MaterializeAndScanOpts(vulnerable, opts)
	if err != nil {
		t.Fatalf("scan vulnerable fixture %s: %v", vulnerable, err)
	}
	if assertErr := integration.AssertVulnerable(vulnResult.Findings, rule, vulnerable); assertErr != nil {
		t.Fatal(assertErr)
	}

	safe := integration.PythonBPFixtureRel(fixtureCase, false)
	safeResult, err := integration.MaterializeAndScanOpts(safe, opts)
	if err != nil {
		t.Fatalf("scan safe fixture %s: %v", safe, err)
	}
	if integration.HasRule(safeResult.Findings, rule) {
		t.Fatalf("safe fixture %s unexpectedly emitted %s: %s", safe, rule, integration.SummarizeFindings(safeResult.Findings))
	}
}
