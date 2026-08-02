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
		{rule: "BP-PY-1", caseName: "BP-PY-1-print-surfaced"},
		{rule: "BP-PY-1", caseName: "BP-PY-1-error-record-return"},
		{rule: "BP-PY-1", caseName: "BP-PY-1-fallthrough-raise"},
		{rule: "BP-PY-1", caseName: "BP-PY-1-raw-capture"},
		{rule: "BP-PY-1", caseName: "BP-PY-1-docstring"},
		{rule: "BP-PY-1", caseName: "BP-PY-1-documented-fallback"},
		{rule: "BP-PY-1", caseName: "BP-PY-1-type-counter"},
		{rule: "BP-PY-1", caseName: "BP-PY-1-ladder-print"},
		{rule: "BP-PY-1", caseName: "BP-PY-1-log-reraise"},
		{rule: "BP-PY-1", caseName: "BP-PY-1-js-bridge-probe"},
		{rule: "BP-PY-1", caseName: "BP-PY-1-return-constant-fallback"},
		{rule: "BP-PY-1", caseName: "BP-PY-1-soft-warning-continue"},
		{rule: "BP-PY-1", caseName: "BP-PY-1-error-dict-payload"},
		{rule: "BP-PY-1", caseName: "BP-PY-1-import-exception-fallback"},
		{rule: "BP-PY-1", caseName: "BP-PY-1-exception-print-exit"},
		{rule: "BP-PY-1", caseName: "BP-PY-1-defensive-return-none"},
		{rule: "BP-PY-2", caseName: "BP-PY-2-import-fallback"},
		{rule: "BP-PY-2", caseName: "BP-PY-2-parsing-fallback"},
		{rule: "BP-PY-2", caseName: "BP-PY-2-expected-exception"},
		{rule: "BP-PY-2", caseName: "BP-PY-2-js-bridge-probe"},
		{rule: "BP-PY-2", caseName: "BP-PY-2-extension-load"},
		{rule: "BP-PY-2", caseName: "BP-PY-2-fall-back-marker"},
		{rule: "BP-PY-2", caseName: "BP-PY-2-conversion-guard"},
		{rule: "BP-PY-2", caseName: "BP-PY-2-except-else"},
		{rule: "BP-PY-2", caseName: "BP-PY-2-index-lookup-pass"},
		{rule: "BP-PY-2", caseName: "BP-PY-2-async-cancel-wait"},
		{rule: "BP-PY-2", caseName: "BP-PY-2-await-expected-exception"},
		{rule: "BP-PY-32", caseName: "BP-PY-32-confined-basename"},
		{rule: "BP-PY-32", caseName: "BP-PY-32-constant-dir"},
		{rule: "BP-PY-2", caseName: "BP-PY-2-attr-probe"},
		{rule: "BP-PY-2", caseName: "BP-PY-2-termination"},
		{rule: "BP-PY-2", caseName: "BP-PY-2-pytest-fail"},
		{rule: "BP-PY-26", caseName: "BP-PY-26-read-only"},
		{rule: "BP-PY-38", caseName: "BP-PY-38-task-list"},
		{rule: "BP-PY-41", caseName: "BP-PY-41-assert-helper"},
		{rule: "BP-PY-41", caseName: "BP-PY-41-assertionerror-helper"},
		{rule: "BP-PY-41", caseName: "BP-PY-41-sibling-helpers"},
		{rule: "BP-PY-41", caseName: "BP-PY-41-check-func"},
		{rule: "BP-PY-41", caseName: "BP-PY-41-benchmark"},
		{rule: "BP-PY-41", caseName: "BP-PY-41-bare-assert-helper"},
		{rule: "BP-PY-41", caseName: "BP-PY-41-string-body-scan"},
		{rule: "BP-PY-41", caseName: "BP-PY-41-raise-assertion"},
		{rule: "BP-PY-41", caseName: "BP-PY-41-retry-fallback"},
		{rule: "BP-PY-41", caseName: "BP-PY-41-noxfile"},
		{rule: "BP-PY-41", caseName: "BP-PY-41-pytest-raises-helper"},
		{rule: "BP-PY-42", caseName: "BP-PY-42-thread-collection"},
		{rule: "BP-PY-42", caseName: "BP-PY-42-batch3-reraise"},
		{rule: "BP-PY-42", caseName: "BP-PY-42-batch3-non-test-helper"},
		{rule: "BP-PY-45", caseName: "BP-PY-45-docs-conf"},
		{rule: "BP-PY-45", caseName: "BP-PY-45-file-bootstrap"},
		{rule: "BP-PY-45", caseName: "BP-PY-45-guarded-bootstrap"},
		{rule: "BP-PY-45", caseName: "BP-PY-45-bootstrap-script"},
		{rule: "BP-PY-45", caseName: "BP-PY-45-library-module"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-cli-output"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-benchmark-script"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-argparse-cli"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-main-guard"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-script-path"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-examples-library-print"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-click-cli"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-typer-cli"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-cyclopts-cli"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-docstring-print"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-main-epilog"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-shebang-script"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-shebang-vestigial-library"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-setup-script"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-string-template"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-cli-module"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-commands-path"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-rich-print"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-console-method"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-demo-script-examples"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-demo-guard-local"},
		{rule: "BP-PY-46", caseName: "BP-PY-46-script-completion"},
		{rule: "BP-PY-7", caseName: "BP-PY-7-attr-method-open"},
		{rule: "BP-PY-7", caseName: "BP-PY-7-def-open"},
		{rule: "BP-PY-12", caseName: "BP-PY-12-attr-method-exec"},
		{rule: "BP-PY-12", caseName: "BP-PY-12-string-literal-exec"},
		{rule: "BP-PY-12", caseName: "BP-PY-12-def-exec"},
		{rule: "BP-PY-12", caseName: "BP-PY-12-builtins-exec"},
		{rule: "BP-PY-11", caseName: "BP-PY-11-ruamel"},
		{rule: "BP-PY-13", caseName: "BP-PY-13-bench-secret"},
		{rule: "BP-PY-13", caseName: "BP-PY-13-env-key-name"},
		{rule: "BP-PY-13", caseName: "BP-PY-13-concat-token"},
		{rule: "BP-PY-13", caseName: "BP-PY-13-fstring-secret"},
		{rule: "BP-PY-49", caseName: "BP-PY-49-fingerprint-pin"},
		{rule: "BP-PY-49", caseName: "BP-PY-49-error-message"},
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
