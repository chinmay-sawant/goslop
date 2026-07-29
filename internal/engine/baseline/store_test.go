package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chinmay/codehound/internal/engine/baseline"
	"github.com/chinmay/codehound/internal/rules"
)

func TestBaselineRoundTripAndFilter(t *testing.T) {
	findings := []rules.Finding{
		{RuleID: "CWE-78", File: "a.go", Line: 10, Column: 1, Message: "cmd", Severity: rules.SeverityHigh},
		{RuleID: "PERF-1", File: "b.go", Line: 2, Column: 1, Message: "slow", Severity: rules.SeverityLow},
	}
	for i := range findings {
		findings[i].EnsureFingerprint()
	}
	b := baseline.FromFindings(findings[:1], "0.1.0-dev")
	if b.EntryCount() != 1 {
		t.Fatalf("count: %d", b.EntryCount())
	}

	dir := t.TempDir()
	path := filepath.Join(dir, baseline.FileName)
	if err := b.Save(path); err != nil {
		t.Fatal(err)
	}
	loaded, err := baseline.Load(path)
	if err != nil {
		t.Fatal(err)
	}

	out, n := loaded.Filter(findings, false)
	if n != 1 || len(out) != 1 || out[0].RuleID != "PERF-1" {
		t.Fatalf("filter: n=%d out=%+v", n, out)
	}

	out, n = loaded.Filter(findings, true)
	if n != 1 || len(out) != 2 {
		t.Fatalf("show baselined: n=%d len=%d", n, len(out))
	}
}

func TestDiscover(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "pkg")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, baseline.FileName)
	b := baseline.FromFindings(nil, "t")
	if err := b.Save(path); err != nil {
		t.Fatal(err)
	}
	found := baseline.Discover(sub)
	if found != path {
		// Compare cleaned paths
		if filepath.Clean(found) != filepath.Clean(path) {
			t.Fatalf("discover: got %q want %q", found, path)
		}
	}
}


