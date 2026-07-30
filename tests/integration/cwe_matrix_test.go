package integration_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/chinmay-sawant/goslop/tests/integration"
)

// TestCWEFixturesMatrix ports Rust go_cwe_fixtures_fire_vulnerable_and_silence_safe
// for frameworks/ + stdlib/ inventories (taint on, --only CWE-N).
func TestCWEFixturesMatrix(t *testing.T) {
	cases, err := integration.DiscoverCWECases()
	if err != nil {
		t.Fatal(err)
	}
	if len(cases) == 0 {
		t.Fatal("expected at least one CWE fixture case")
	}
	t.Logf("CWE fixture cases: %d ×2 suites ×2 files", len(cases))

	var failures []string
	for _, cwe := range cases {
		opts := integration.DefaultMatrixOptions()
		opts.TaintEnabled = true
		opts.Only = []string{cwe}

		for _, suite := range []string{"frameworks", "stdlib"} {
			vulnRel := integration.CWEFixtureRel(suite, cwe, true)
			vres, err := integration.MaterializeAndScanOpts(vulnRel, opts)
			if err != nil {
				failures = append(failures, fmt.Sprintf("%s/%s vuln: %v", suite, cwe, err))
			} else if aerr := integration.AssertVulnerable(vres.Findings, cwe, vulnRel); aerr != nil {
				failures = append(failures, aerr.Error())
			}
			safeRel := integration.CWEFixtureRel(suite, cwe, false)
			sres, err := integration.MaterializeAndScanOpts(safeRel, opts)
			if err != nil {
				failures = append(failures, fmt.Sprintf("%s/%s safe: %v", suite, cwe, err))
			} else if integration.HasRule(sres.Findings, cwe) {
				failures = append(failures, fmt.Sprintf("%s: expected no %s on safe, got %s",
					safeRel, cwe, integration.SummarizeFindings(sres.Findings)))
			}
		}
	}
	unexpected, known := integration.PartitionFailures(failures)
	if len(known) > 0 {
		t.Logf("CWE quarantined gaps: %d (see fixture_quarantine.txt)", len(known))
	}
	if len(unexpected) > 0 {
		t.Fatalf("CWE unexpected failures %d (quarantined %d):\n%s",
			len(unexpected), len(known),
			integration.FormatMatrixReport(unexpected, known, 40))
	}
}

// TestTaintFixturesMatrix ports Rust taint_cwe_fixtures for go/taint CWE-* pairs
// (skips IP-* project variants which use different rule ids).
func TestTaintFixturesMatrix(t *testing.T) {
	cases, err := integration.DiscoverTaintCases()
	if err != nil {
		t.Fatal(err)
	}
	var cweCases []string
	for _, c := range cases {
		if strings.HasPrefix(c, "CWE-") {
			cweCases = append(cweCases, c)
		}
	}
	if len(cweCases) == 0 {
		t.Fatal("expected CWE taint fixtures")
	}
	t.Logf("taint CWE cases: %d", len(cweCases))

	var failures []string
	for _, c := range cweCases {
		ruleID := splitCWECase(c)
		opts := integration.DefaultMatrixOptions()
		opts.TaintEnabled = true
		opts.Only = []string{ruleID}

		vulnRel := integration.TaintFixtureRel(c, true)
		vres, err := integration.MaterializeAndScanOpts(vulnRel, opts)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s vuln: %v", c, err))
		} else if aerr := integration.AssertVulnerable(vres.Findings, ruleID, vulnRel); aerr != nil {
			failures = append(failures, aerr.Error())
		}
		safeRel := integration.TaintFixtureRel(c, false)
		sres, err := integration.MaterializeAndScanOpts(safeRel, opts)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s safe: %v", c, err))
		} else if integration.HasRule(sres.Findings, ruleID) {
			failures = append(failures, fmt.Sprintf("%s: expected no %s on safe, got %s",
				safeRel, ruleID, integration.SummarizeFindings(sres.Findings)))
		}
	}
	unexpected, known := integration.PartitionFailures(failures)
	if len(known) > 0 {
		t.Logf("taint quarantined gaps: %d", len(known))
	}
	if len(unexpected) > 0 {
		t.Fatalf("taint unexpected failures %d (quarantined %d):\n%s",
			len(unexpected), len(known),
			integration.FormatMatrixReport(unexpected, known, 40))
	}
}

func splitCWECase(c string) string {
	if len(c) < 5 || c[:4] != "CWE-" {
		return c
	}
	rest := c[4:]
	for i, ch := range rest {
		if ch == '-' {
			return "CWE-" + rest[:i]
		}
	}
	return c
}
