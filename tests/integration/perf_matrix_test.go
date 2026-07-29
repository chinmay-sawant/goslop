package integration_test

import (
	"fmt"
	"testing"

	"github.com/chinmay/goslop/tests/integration"
)

// TestPERFFixturesMatrix ports Rust go_perf_fixtures_fire_vulnerable_and_silence_safe
// with --only PERF-N pinning (primary-rule oracle).
func TestPERFFixturesMatrix(t *testing.T) {
	cases, err := integration.DiscoverPERFCases()
	if err != nil {
		t.Fatal(err)
	}
	if len(cases) == 0 {
		t.Fatal("expected at least one PERF fixture case")
	}
	t.Logf("PERF fixture cases: %d (×2 files)", len(cases))

	var failures []string
	for _, c := range cases {
		rule := integration.PERFRuleID(c)
		opts := integration.DefaultMatrixOptions()
		opts.Only = []string{rule}

		vulnRel := integration.PERFFixtureRel(c, true)
		vres, err := integration.MaterializeAndScanOpts(vulnRel, opts)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s vuln: %v", c, err))
		} else if aerr := integration.AssertVulnerable(vres.Findings, rule, vulnRel); aerr != nil {
			failures = append(failures, aerr.Error())
		}

		safeRel := integration.PERFFixtureRel(c, false)
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
		t.Logf("PERF quarantined gaps: %d (see fixture_quarantine.txt)", len(known))
	}
	if len(unexpected) > 0 {
		t.Fatalf("PERF unexpected failures %d (quarantined %d) of %d cases:\n%s",
			len(unexpected), len(known), len(cases),
			integration.FormatMatrixReport(unexpected, known, 40))
	}
}
