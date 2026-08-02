package perf_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/fixture"
	perf "github.com/chinmay-sawant/goslop/internal/lang/python/detectors/perf"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestPERFFalsePositiveAuditFixtureVariants(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		rule     string
		caseName string
	}{
		{rule: "PERF-PY-27", caseName: "PERF-PY-27-unique-path"},
		{rule: "PERF-PY-27", caseName: "PERF-PY-27-analyze-once"},
		{rule: "PERF-PY-27", caseName: "PERF-PY-27-derived-path"},
		{rule: "PERF-PY-23", caseName: "PERF-PY-23-seed"},
		{rule: "PERF-PY-26", caseName: "PERF-PY-26-cli-parse"},
		{rule: "PERF-PY-26", caseName: "PERF-PY-26-parser-descent"},
		{rule: "PERF-PY-26", caseName: "PERF-PY-26-field-decode"},
		{rule: "PERF-PY-25", caseName: "PERF-PY-25-key-lambda"},
		{rule: "PERF-PY-25", caseName: "PERF-PY-25-early-return"},
	} {
		tc := tc
		t.Run(tc.caseName, func(t *testing.T) {
			t.Parallel()
			assertPERFFixtureCase(t, tc.rule, tc.caseName)
		})
	}
}

func assertPERFFixtureCase(t *testing.T, rule, fixtureCase string) {
	t.Helper()
	vuln := loadPERFFixture(t, fixtureCase, true)
	assertPERFRuleFire(t, rule, vuln.path, vuln.body, true)

	safe := loadPERFFixture(t, fixtureCase, false)
	assertPERFRuleFire(t, rule, safe.path, safe.body, false)
}

type perfFixtureSource struct {
	path string
	body string
}

func loadPERFFixture(t *testing.T, caseName string, vulnerable bool) perfFixtureSource {
	t.Helper()
	suf := "safe"
	if vulnerable {
		suf = "vulnerable"
	}
	txtPath := filepath.Join(repoPERFFixturesRoot(t), caseName+"-"+suf+".txt")
	data, err := os.ReadFile(txtPath)
	if err != nil {
		t.Fatalf("read fixture %s: %v", txtPath, err)
	}
	fx, err := fixture.ParseFixture(string(data), txtPath)
	if err != nil {
		t.Fatalf("parse fixture %s: %v", txtPath, err)
	}
	path := fx.Filename
	if path == "" {
		path = caseName + "-" + suf + ".py"
	}
	return perfFixtureSource{path: path, body: fx.Source}
}

func assertPERFRuleFire(t *testing.T, rule, path, source string, want bool) {
	t.Helper()
	ctx := core.DefaultScanContext()
	ctx.Only = []string{rule}
	unit := core.NewParsedUnit(core.LanguagePython, path, source)
	var findings []rules.Finding
	perf.NewPythonPerfScan().Run(ctx, unit, &findings)
	got := false
	for _, finding := range findings {
		if finding.RuleID == rule {
			got = true
			break
		}
	}
	if got != want {
		t.Fatalf("%s on %s: got fire=%v, want %v; findings=%#v\nsource:\n%s",
			rule, path, got, want, findings, source)
	}
}

func repoPERFFixturesRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "..", ".."))
	fx := filepath.Join(root, "tests", "fixtures", "python", "perf")
	if _, err := os.Stat(fx); err != nil {
		t.Fatalf("fixtures root %s: %v", fx, err)
	}
	return fx
}
