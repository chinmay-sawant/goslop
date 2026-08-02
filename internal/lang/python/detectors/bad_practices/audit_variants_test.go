package badpractices_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/fixture"
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
		{rule: "BP-PY-41", caseName: "BP-PY-41-assert-helper"},
		{rule: "BP-PY-41", caseName: "BP-PY-41-assertionerror-helper"},
		{rule: "BP-PY-41", caseName: "BP-PY-41-sibling-helpers"},
		{rule: "BP-PY-42", caseName: "BP-PY-42-thread-collection"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-cli-output"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-benchmark-script"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-argparse-cli"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-main-guard"},
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
	vuln := loadBPFixture(t, fixtureCase, true)
	assertRule(t, rule, vuln.path, vuln.body, true)

	safe := loadBPFixture(t, fixtureCase, false)
	assertRule(t, rule, safe.path, safe.body, false)
}

type bpFixtureSource struct {
	path string
	body string
}

func loadBPFixture(t *testing.T, caseName string, vulnerable bool) bpFixtureSource {
	t.Helper()
	suf := "safe"
	if vulnerable {
		suf = "vulnerable"
	}
	return loadBPFixtureFile(t, caseName+"-"+suf+".txt")
}

// loadBPFixtureFile loads any .txt under tests/fixtures/python/bp/ (including
// non-paired helpers used by TempDir sibling tests).
func loadBPFixtureFile(t *testing.T, fileName string) bpFixtureSource {
	t.Helper()
	txtPath := filepath.Join(repoFixturesRoot(t), "python", "bp", fileName)
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
		path = strings.TrimSuffix(fileName, ".txt") + ".py"
	}
	return bpFixtureSource{path: path, body: fx.Source}
}

func repoFixturesRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	// .../detectors/bad_practices/audit_variants_test.go → repo root is ../../../../..
	root := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "..", ".."))
	fx := filepath.Join(root, "tests", "fixtures")
	if _, err := os.Stat(fx); err != nil {
		t.Fatalf("fixtures root %s: %v", fx, err)
	}
	return fx
}
