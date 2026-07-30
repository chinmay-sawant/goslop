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
