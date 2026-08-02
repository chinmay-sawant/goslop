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
		{rule: "BP-PY-1", caseName: "BP-PY-1-batch3-exc-info"},
		{rule: "BP-PY-1", caseName: "BP-PY-1-batch3-error-result"},
		{rule: "BP-PY-1", caseName: "BP-PY-1-batch3-set-exception"},
		{rule: "BP-PY-1", caseName: "BP-PY-1-batch3-result-error"},
		{rule: "BP-PY-26", caseName: "BP-PY-26-read-only"},
		{rule: "BP-PY-38", caseName: "BP-PY-38-task-list"},
		{rule: "BP-PY-41", caseName: "BP-PY-41-assert-helper"},
		{rule: "BP-PY-41", caseName: "BP-PY-41-assertionerror-helper"},
		{rule: "BP-PY-41", caseName: "BP-PY-41-sibling-helpers"},
		{rule: "BP-PY-41", caseName: "BP-PY-41-check-func"},
		{rule: "BP-PY-41", caseName: "BP-PY-41-benchmark"},
		{rule: "BP-PY-41", caseName: "BP-PY-41-bare-assert-helper"},
		{rule: "BP-PY-41", caseName: "BP-PY-41-string-body-scan"},
		{rule: "BP-PY-42", caseName: "BP-PY-42-thread-collection"},
		{rule: "BP-PY-42", caseName: "BP-PY-42-batch3-reraise"},
		{rule: "BP-PY-42", caseName: "BP-PY-42-batch3-non-test-helper"},
		{rule: "BP-PY-45", caseName: "BP-PY-45-docs-conf"},
		{rule: "BP-PY-45", caseName: "BP-PY-45-file-bootstrap"},
		{rule: "BP-PY-45", caseName: "BP-PY-45-guarded-bootstrap"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-cli-output"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-benchmark-script"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-argparse-cli"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-main-guard"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-script-path"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-click-cli"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-typer-cli"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-cyclopts-cli"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-docstring-print"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-main-epilog"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-shebang-script"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-setup-script"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-string-template"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-cli-module"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-commands-path"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-rich-print"},
		{rule: "BP-PY-7", caseName: "BP-PY-7-attr-method-open"},
		{rule: "BP-PY-7", caseName: "BP-PY-7-def-open"},
		{rule: "BP-PY-12", caseName: "BP-PY-12-attr-method-exec"},
		{rule: "BP-PY-12", caseName: "BP-PY-12-string-literal-exec"},
		{rule: "BP-PY-11", caseName: "BP-PY-11-ruamel"},
		{rule: "BP-PY-13", caseName: "BP-PY-13-bench-secret"},
		{rule: "BP-PY-49", caseName: "BP-PY-49-fingerprint-pin"},
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
