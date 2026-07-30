package app

import (
	"bytes"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/engine"
)

func TestRunVersion(t *testing.T) {
	var out, errBuf bytes.Buffer
	if err := run([]string{"-version"}, &out, &errBuf); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(out.String()) != Version {
		t.Fatalf("version output: %q", out.String())
	}
}

func TestRunListRulesNoPanic(t *testing.T) {
	var out, errBuf bytes.Buffer
	if err := run([]string{"--list-rules"}, &out, &errBuf); err != nil {
		t.Fatal(err)
	}
	s := out.String()
	if s == "" {
		t.Fatal("expected list-rules output")
	}
	// Either seed rules or the empty-registry stub message.
	if !strings.Contains(s, "no rules registered") &&
		!strings.Contains(s, "CWE-") &&
		!strings.Contains(s, "PERF-") {
		t.Fatalf("list-rules unexpected: %q", s)
	}
}

func TestRunUnknownProfile(t *testing.T) {
	var out, errBuf bytes.Buffer
	err := run([]string{"-profile", "nope", "."}, &out, &errBuf)
	if err == nil {
		t.Fatal("expected error")
	}
	if ExitCode(err) != ExitConfig {
		t.Fatalf("code: %d (%v)", ExitCode(err), err)
	}
}

func TestRunHelp(t *testing.T) {
	var out, errBuf bytes.Buffer
	if err := run([]string{"-h"}, &out, &errBuf); err != nil {
		t.Fatalf("help should be clean exit: %v", err)
	}
}

func TestRunScanEmptyDirJSON(t *testing.T) {
	dir := t.TempDir()
	var out, errBuf bytes.Buffer
	if err := run([]string{"-format", "json", "-profile", "all", dir}, &out, &errBuf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"findings"`) {
		t.Fatalf("json: %q", out.String())
	}
}

func TestRunPartialScanWritesErrorsToStderrAndFails(t *testing.T) {
	dir := t.TempDir()
	badSource := filepath.Join(dir, "invalid.go")
	if err := os.WriteFile(badSource, []byte("package invalid\n\xff"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errBuf bytes.Buffer
	err := run([]string{"--format", "json", "--no-cache", dir}, &out, &errBuf)
	if ExitCode(err) != ExitInternal {
		t.Fatalf("exit code=%d err=%v", ExitCode(err), err)
	}
	if out.Len() != 0 {
		t.Fatalf("partial scan must not emit a JSON payload: %q", out.String())
	}
	if !strings.Contains(errBuf.String(), "analysis error [encoding]") ||
		!strings.Contains(errBuf.String(), "incomplete: 1 file(s)") {
		t.Fatalf("stderr=%q", errBuf.String())
	}
}

func TestRunUnreadableFileWritesErrorsToStderrAndFails(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "unreadable.go")
	if err := os.Symlink(filepath.Join(dir, "missing.go"), path); err != nil {
		t.Skipf("cannot create dangling symlink: %v", err)
	}
	var out, errBuf bytes.Buffer
	err := run([]string{"--format", "sarif", "--no-cache", dir}, &out, &errBuf)
	if ExitCode(err) != ExitInternal {
		t.Fatalf("exit code=%d err=%v", ExitCode(err), err)
	}
	if out.Len() != 0 {
		t.Fatalf("partial scan must not emit a SARIF payload: %q", out.String())
	}
	if !strings.Contains(errBuf.String(), "analysis error [") ||
		!strings.Contains(errBuf.String(), "unreadable.go") ||
		!strings.Contains(errBuf.String(), "incomplete: 1 file(s)") {
		t.Fatalf("stderr=%q", errBuf.String())
	}
}

func TestReportScanErrorsWritesParseFailuresToStderr(t *testing.T) {
	var stderr bytes.Buffer
	reportScanErrors(&stderr, []engine.ScanError{{
		Path: "broken.go", Kind: engine.ScanErrorParse, Message: "expected declaration",
	}})
	if got := stderr.String(); !strings.Contains(got, "analysis error [parse]: broken.go: expected declaration") {
		t.Fatalf("stderr=%q", got)
	}
}

func TestRunParseFallbackWritesWarningAndKeepsMachineOutput(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "broken.go")
	if err := os.WriteFile(path, []byte("package sample\nfunc broken( {\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errBuf bytes.Buffer
	if err := run([]string{"--format", "json", "--no-cache", dir}, &out, &errBuf); err != nil {
		t.Fatalf("parse fallback should remain non-fatal: %v", err)
	}
	if !strings.Contains(out.String(), `"findings"`) {
		t.Fatalf("expected JSON output, got %q", out.String())
	}
	if got := strings.Count(errBuf.String(), "analysis warning [parse]"); got != 1 {
		t.Fatalf("parse diagnostic count=%d stderr=%q", got, errBuf.String())
	}
}

func TestRunLanguagesGoOnlyIgnoresPythonFiles(t *testing.T) {
	dir := t.TempDir()
	// Default (unset languages) and explicit languages=["go"] both ignore .py.
	if err := os.WriteFile(filepath.Join(dir, "goslop.toml"), []byte(`[goslop]
languages = ["go"]
fail_on = "none"
`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "ok.go"), []byte("package sample\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// A deliberately "bad" python file that would fail encoding/parse if scanned.
	if err := os.WriteFile(filepath.Join(dir, "bad.py"), []byte("print('hi')\n\xff"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errBuf bytes.Buffer
	if err := run([]string{"--format", "json", "--no-cache", dir}, &out, &errBuf); err != nil {
		t.Fatalf("go-only scan should ignore .py: %v\nstderr=%s", err, errBuf.String())
	}
	// scanned 1 file (the .go), not the .py
	if !strings.Contains(errBuf.String(), "scanned 1 files") {
		t.Fatalf("expected single go file scanned, stderr=%q", errBuf.String())
	}
}

func TestRunLanguagesPythonOnlyScansPyNotGo(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "goslop.toml"), []byte(`[goslop]
languages = ["python"]
fail_on = "none"
`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package sample\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "app.py"), []byte("x = 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errBuf bytes.Buffer
	// Built-in Python stub is registered via NewRegistryWithLanguages: scan .py only.
	if err := run([]string{"--format", "json", "--no-cache", dir}, &out, &errBuf); err != nil {
		t.Fatalf("python-only scan must not crash: %v\nstderr=%s", err, errBuf.String())
	}
	if !strings.Contains(errBuf.String(), "scanned 1 files") {
		t.Fatalf("expected only app.py scanned, stderr=%q", errBuf.String())
	}
	if !strings.Contains(out.String(), `"findings"`) {
		t.Fatalf("json: %q", out.String())
	}
}

func TestRunListRulesRespectsLanguagesConfig(t *testing.T) {
	dir := t.TempDir()
	// languages=["python"] uses the stub plugin → empty detector catalogue.
	cfg := filepath.Join(dir, "goslop.toml")
	if err := os.WriteFile(cfg, []byte(`[goslop]
languages = ["python"]
`), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errBuf bytes.Buffer
	if err := run([]string{"--list-rules", "--config", cfg}, &out, &errBuf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "no rules registered") {
		t.Fatalf("expected empty rule list for Python stub (zero detectors): %q", out.String())
	}

	// languages=["go"] still lists Go rules.
	if err := os.WriteFile(cfg, []byte(`[goslop]
languages = ["go"]
`), 0o644); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	errBuf.Reset()
	if err := run([]string{"--list-rules", "--config", cfg}, &out, &errBuf); err != nil {
		t.Fatal(err)
	}
	s := out.String()
	if !strings.Contains(s, "CWE-") && !strings.Contains(s, "PERF-") {
		t.Fatalf("expected go rules: %q", s)
	}
}

func TestRunInvalidLanguagesConfigRejected(t *testing.T) {
	dir := t.TempDir()
	cfg := filepath.Join(dir, "goslop.toml")
	if err := os.WriteFile(cfg, []byte(`[goslop]
languages = ["ruby"]
`), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errBuf bytes.Buffer
	err := run([]string{"--list-rules", "--config", cfg}, &out, &errBuf)
	if ExitCode(err) != ExitConfig {
		t.Fatalf("exit=%d err=%v", ExitCode(err), err)
	}
}

func TestRunConfigFailPolicyHasObservableEffect(t *testing.T) {
	makeProject := func(t *testing.T, configBody string) string {
		t.Helper()
		dir := t.TempDir()
		if configBody != "" {
			if err := os.WriteFile(filepath.Join(dir, "goslop.toml"), []byte(configBody), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		source := `package sample

import "regexp"

func slow() {
	for i := 0; i < 2; i++ {
		_ = regexp.MustCompile("x")
	}
}
`
		if err := os.WriteFile(filepath.Join(dir, "slow.go"), []byte(source), 0o644); err != nil {
			t.Fatal(err)
		}
		return dir
	}

	withoutConfig := makeProject(t, "")
	var out, errBuf bytes.Buffer
	if err := run([]string{"--profile", "all", "--no-cache", withoutConfig}, &out, &errBuf); ExitCode(err) != ExitFailing {
		t.Fatalf("default policy exit=%d err=%v stderr=%q", ExitCode(err), err, errBuf.String())
	}

	withConfig := makeProject(t, "[goslop]\nfail_on = \"none\"\n")
	out.Reset()
	errBuf.Reset()
	if err := run([]string{"--profile", "all", "--no-cache", withConfig}, &out, &errBuf); err != nil {
		t.Fatalf("config fail_on=none should be observable as a clean exit: %v\nstderr=%s", err, errBuf.String())
	}
}

func TestRunConfigExportWholeFunctionAffectsContext(t *testing.T) {
	dir := t.TempDir()
	contextDir := filepath.Join(dir, "context")
	if err := os.WriteFile(filepath.Join(dir, "goslop.toml"), []byte("[goslop.export]\nwhole_function = false\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	source := `package sample

import "fmt"

func slow() {
	first := 1
	second := first + 1
	third := second + 1
	fourth := third + 1
	for i := 0; i < fourth; i++ {
		rendered := fmt.Sprintf("%d", i)
		_ = rendered
	}
}
`
	if err := os.WriteFile(filepath.Join(dir, "slow.go"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}

	var out, errBuf bytes.Buffer
	if err := run([]string{
		"--profile", "all", "--only", "PERF-6", "--no-cache", "--no-fail",
		"--export-context", "--context-dir", contextDir, dir,
	}, &out, &errBuf); err != nil {
		t.Fatalf("run: %v\nstderr=%s", err, errBuf.String())
	}
	body, err := os.ReadFile(filepath.Join(contextDir, "1.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(body), "func slow()") {
		t.Fatalf("whole_function=false should use a nearby window:\n%s", body)
	}
	if !strings.Contains(string(body), "fmt.Sprintf") {
		t.Fatalf("export did not include the finding context:\n%s", body)
	}
}

func TestRunPruneCacheOpenFailureIsExplicit(t *testing.T) {
	dir := t.TempDir()
	cachePath := filepath.Join(dir, "not-a-cache-directory")
	if err := os.WriteFile(cachePath, []byte("not a directory"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errBuf bytes.Buffer
	err := run([]string{"--prune-cache", "--cache-dir", cachePath, dir}, &out, &errBuf)
	if ExitCode(err) != ExitInternal {
		t.Fatalf("exit code=%d err=%v", ExitCode(err), err)
	}
	if err == nil || !strings.Contains(err.Error(), "open cache for pruning") {
		t.Fatalf("error=%v", err)
	}
}

func TestRunInit(t *testing.T) {
	dir := t.TempDir()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if chdirErr := os.Chdir(dir); chdirErr != nil {
		t.Fatal(chdirErr)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })

	var out bytes.Buffer
	if initErr := run([]string{"init"}, &out, ioDiscard{}); initErr != nil {
		t.Fatal(initErr)
	}
	if !strings.Contains(out.String(), "wrote starter goslop.toml to") {
		t.Fatalf("init output was not written to the supplied stdout: %q", out.String())
	}
	if _, statErr := os.Stat(filepath.Join(dir, "goslop.toml")); statErr != nil {
		t.Fatal(statErr)
	}
	// Second init should fail (already exists).
	err = run([]string{"init"}, ioDiscard{}, ioDiscard{})
	if err == nil || ExitCode(err) != ExitConfig {
		t.Fatalf("second init: %v", err)
	}
}

func TestExitCodeMapping(t *testing.T) {
	if ExitCode(nil) != ExitClean {
		t.Fatal()
	}
	if ExitCode(&ExitCodeError{Code: ExitFailing}) != ExitFailing {
		t.Fatal()
	}
	if ExitCode(errors.New("x")) != ExitConfig {
		t.Fatal()
	}
	if !errors.Is(flag.ErrHelp, flag.ErrHelp) {
		t.Fatal()
	}
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) { return len(p), nil }
