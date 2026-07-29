package app

import (
	"bytes"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
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

	if initErr := run([]string{"init"}, ioDiscard{}, ioDiscard{}); initErr != nil {
		t.Fatal(initErr)
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
