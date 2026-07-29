package integration_test

import (
	"fmt"
	"testing"

	"github.com/chinmay/goslop/tests/integration"
)

// TestBPFixturesMatrix ports Rust go_bad_practice_fixtures_fire_vulnerable_and_silence_safe.
// Each tests/fixtures/go/bad_practices/*-{vulnerable,safe}.txt pair is scanned with
// --only BP-N (pinned) so the suite measures that rule’s heuristic, not cross-rule noise.
func TestBPFixturesMatrix(t *testing.T) {
	cases, err := integration.DiscoverBPCases()
	if err != nil {
		t.Fatal(err)
	}
	if len(cases) == 0 {
		t.Fatal("expected at least one BP fixture case")
	}
	t.Logf("BP fixture cases: %d (×2 files)", len(cases))

	var failures []string
	for _, c := range cases {
		rule := integration.BPRuleID(c)
		opts := integration.DefaultMatrixOptions()
		opts.Only = []string{rule}

		vulnRel := integration.BPFixtureRel(c, true)
		vres, err := integration.MaterializeAndScanOpts(vulnRel, opts)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s vuln: %v", c, err))
		} else if aerr := integration.AssertVulnerable(vres.Findings, rule, vulnRel); aerr != nil {
			failures = append(failures, aerr.Error())
		}

		safeRel := integration.BPFixtureRel(c, false)
		sres, err := integration.MaterializeAndScanOpts(safeRel, opts)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s safe: %v", c, err))
		} else if integration.HasRule(sres.Findings, rule) {
			failures = append(failures, fmt.Sprintf("%s: expected no %s on safe, got %s",
				safeRel, rule, integration.SummarizeFindings(sres.Findings)))
		}
	}
	unexpected, known := integration.PartitionFailures(failures)
	if len(known) > 0 {
		t.Logf("BP quarantined gaps: %d (see fixture_quarantine.txt)", len(known))
	}
	if len(unexpected) > 0 {
		t.Fatalf("BP unexpected failures %d (quarantined %d) of %d cases:\n%s",
			len(unexpected), len(known), len(cases),
			integration.FormatMatrixReport(unexpected, known, 40))
	}
}

// TestBPFixtureInventorySorted mirrors Rust inventory shape checks.
func TestBPFixtureInventorySorted(t *testing.T) {
	cases, err := integration.DiscoverBPCases()
	if err != nil {
		t.Fatal(err)
	}
	if len(cases) == 0 {
		t.Fatal("empty BP inventory")
	}
	seen := map[string]struct{}{}
	for _, c := range cases {
		if _, ok := seen[c]; ok {
			t.Fatalf("duplicate BP case %q", c)
		}
		seen[c] = struct{}{}
	}
}
