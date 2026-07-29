package cli

import (
	"flag"
	"strings"
	"testing"
)

func TestParseDefaults(t *testing.T) {
	opts, err := Parse(nil)
	if err != nil {
		t.Fatal(err)
	}
	if opts.Profile != "recommended" {
		t.Fatalf("profile: got %q", opts.Profile)
	}
	if opts.Format != FormatText {
		t.Fatalf("format: got %q", opts.Format)
	}
	if len(opts.Paths) != 1 || opts.Paths[0] != "." {
		t.Fatalf("paths: got %#v", opts.Paths)
	}
	if opts.ListRules || opts.IncludeTests || opts.NoCache || opts.Version {
		t.Fatalf("unexpected bools: %+v", opts)
	}
}

func TestParseFlagsDoubleAndSingleDash(t *testing.T) {
	opts, err := Parse([]string{
		"--profile", "security",
		"-only", "CWE-22,CWE-89",
		"--skip", "PERF-1, BP-2",
		"-format", "json",
		"--list-rules",
		"--include-tests",
		"--no-cache",
		"./cmd", "pkg/",
	})
	if err != nil {
		t.Fatal(err)
	}
	if opts.Profile != "security" {
		t.Fatalf("profile: %q", opts.Profile)
	}
	if opts.Format != FormatJSON {
		t.Fatalf("format: %q", opts.Format)
	}
	if !opts.ListRules || !opts.IncludeTests || !opts.NoCache {
		t.Fatalf("bools: list=%v tests=%v nocache=%v", opts.ListRules, opts.IncludeTests, opts.NoCache)
	}
	wantOnly := []string{"CWE-22", "CWE-89"}
	if !stringSlicesEqual(opts.Only, wantOnly) {
		t.Fatalf("only: %#v", opts.Only)
	}
	wantSkip := []string{"PERF-1", "BP-2"}
	if !stringSlicesEqual(opts.Skip, wantSkip) {
		t.Fatalf("skip: %#v", opts.Skip)
	}
	if !stringSlicesEqual(opts.Paths, []string{"./cmd", "pkg/"}) {
		t.Fatalf("paths: %#v", opts.Paths)
	}
}

func TestParseVersion(t *testing.T) {
	opts, err := Parse([]string{"-version"})
	if err != nil {
		t.Fatal(err)
	}
	if !opts.Version {
		t.Fatal("expected Version")
	}
}

func TestParseTaintFlags(t *testing.T) {
	opts, err := Parse([]string{
		"--taint",
		"--taint-depth", "3",
		"--taint-show-paths",
		".",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !opts.Taint {
		t.Fatal("expected Taint")
	}
	if opts.TaintDepth != 3 {
		t.Fatalf("depth: %d", opts.TaintDepth)
	}
	if !opts.TaintShowPaths {
		t.Fatal("expected TaintShowPaths")
	}

	opts2, err := Parse([]string{"--no-taint", "."})
	if err != nil {
		t.Fatal(err)
	}
	if !opts2.NoTaint {
		t.Fatal("expected NoTaint")
	}
}

func TestParseInit(t *testing.T) {
	opts, err := Parse([]string{"init"})
	if err != nil {
		t.Fatal(err)
	}
	if opts.Command != "init" {
		t.Fatalf("command: %q", opts.Command)
	}
}

func TestParseInvalidFormat(t *testing.T) {
	_, err := Parse([]string{"-format", "xml"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "format") {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestParseHelp(t *testing.T) {
	var buf strings.Builder
	_, err := ParseWithOutput([]string{"-h"}, &buf)
	if err != flag.ErrHelp {
		t.Fatalf("want flag.ErrHelp, got %v", err)
	}
	if !strings.Contains(buf.String(), "Usage:") {
		t.Fatalf("help text missing: %q", buf.String())
	}
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
