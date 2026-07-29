package integration_test

import (
	"path/filepath"
	"testing"

	"github.com/chinmay/goslop/tests/integration"
)

// TestSeedFixtureOracle exercises the Phase 12 integration harness on a small
// seed set of materialized fixtures (CWE-78/89 + PERF-6). This is scaffolding
// only — full fixture matrix and §12.4 product oracle remain deferred.
func TestSeedFixtureOracle(t *testing.T) {
	if _, err := integration.FixturesRoot(); err != nil {
		t.Fatalf("fixtures not available: %v", err)
	}

	for _, tc := range integration.SeedCases() {
		tc := tc
		t.Run(filepath.ToSlash(tc.RelPath), func(t *testing.T) {
			t.Parallel()
			outRoot := t.TempDir()
			res, err := integration.MaterializeAndScan(tc, outRoot)
			if err != nil {
				t.Fatalf("MaterializeAndScan: %v", err)
			}
			if res.Fired != tc.ExpectFire {
				t.Fatalf("rule %s fire=%v want %v; findings: %s",
					tc.RuleID, res.Fired, tc.ExpectFire,
					integration.SummarizeFindings(res.Findings))
			}
		})
	}
}

func TestRepoRootResolves(t *testing.T) {
	root, err := integration.RepoRoot()
	if err != nil {
		t.Fatal(err)
	}
	if root == "" {
		t.Fatal("empty root")
	}
	fx, err := integration.FixturesRoot()
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(fx) != "fixtures" {
		t.Fatalf("fixtures base: %s", fx)
	}
}
