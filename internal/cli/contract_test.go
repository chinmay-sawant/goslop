package cli

import (
	"strings"
	"testing"
)

// Profile / CLI contract tests (Phase 12.1 scaffolding).
// These lock the flag surface and default profile contract without requiring a full scan.

func TestParseAllProfiles(t *testing.T) {
	// CLI accepts these profile names; app layer maps them via core.ParseProfile.
	for _, name := range []string{
		"recommended", "default", "ci",
		"perf", "performance",
		"security", "sec",
		"style", "bp",
		"all", "full",
	} {
		opts, err := Parse([]string{"-profile", name, "."})
		if err != nil {
			t.Fatalf("profile %q: %v", name, err)
		}
		if opts.Profile != name {
			t.Fatalf("profile stored: got %q want %q", opts.Profile, name)
		}
	}
}

func TestParseFormats(t *testing.T) {
	for _, format := range []string{"text", "json", "sarif"} {
		opts, err := Parse([]string{"-format", format})
		if err != nil {
			t.Fatalf("format %q: %v", format, err)
		}
		if string(opts.Format) != format {
			t.Fatalf("format: got %q want %q", opts.Format, format)
		}
	}
}

func TestParseOnlySkipCSV(t *testing.T) {
	opts, err := Parse([]string{
		"--only", "CWE-78, PERF-6 ,CWE-89",
		"--skip", "PERF-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(opts.Only) != 3 || opts.Only[0] != "CWE-78" || opts.Only[1] != "PERF-6" || opts.Only[2] != "CWE-89" {
		t.Fatalf("only: %#v", opts.Only)
	}
	if len(opts.Skip) != 1 || opts.Skip[0] != "PERF-1" {
		t.Fatalf("skip: %#v", opts.Skip)
	}
}

func TestParseDefaultPathWhenNone(t *testing.T) {
	opts, err := Parse([]string{"--profile", "all"})
	if err != nil {
		t.Fatal(err)
	}
	if len(opts.Paths) != 1 || opts.Paths[0] != "." {
		t.Fatalf("default paths: %#v", opts.Paths)
	}
}

func TestParseSARIFAlias(t *testing.T) {
	// Ensure sarif is accepted (contract for --format sarif used in CI docs).
	opts, err := Parse([]string{"--format", "sarif", "."})
	if err != nil {
		t.Fatal(err)
	}
	if opts.Format != FormatSARIF {
		t.Fatalf("format: %q", opts.Format)
	}
}

func TestParseRejectsUnknownFormatMessage(t *testing.T) {
	_, err := Parse([]string{"-format", "yaml"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "format") {
		t.Fatalf("error should mention format: %v", err)
	}
}
