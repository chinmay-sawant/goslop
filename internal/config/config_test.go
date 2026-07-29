package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chinmay/codehound/internal/config"
	"github.com/chinmay/codehound/internal/core"
)

func TestParseAndMergeAdditiveOnlySkip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "codehound.toml")
	body := `
[codehound]
only = ["PERF-1"]
skip = ["BP-1"]
fail_on = "high"
exclude_tests = false
include = ["**/*.go"]
exclude = ["vendor/**"]

[codehound.cache]
enabled = false

[codehound.taint]
enabled = true
show_paths = true
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	m, err := config.LoadAndMerge(config.MergeInput{
		Only:       []string{"PERF-2"},
		Skip:       []string{"BP-2"},
		ConfigPath: path,
		Paths:      []string{dir},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Only) != 2 || m.Only[0] != "PERF-1" || m.Only[1] != "PERF-2" {
		t.Fatalf("only=%v", m.Only)
	}
	if len(m.Skip) != 2 {
		t.Fatalf("skip=%v", m.Skip)
	}
	if !m.IncludeTests {
		t.Fatal("exclude_tests=false should include tests")
	}
	if !m.NoCache {
		t.Fatal("cache.enabled=false → NoCache")
	}
	if !m.Taint || !m.TaintShowPaths {
		t.Fatalf("taint flags: taint=%v show=%v", m.Taint, m.TaintShowPaths)
	}
	if m.FailPolicy == nil || *m.FailPolicy != core.FailHigh {
		t.Fatalf("fail policy=%v", m.FailPolicy)
	}
	if len(m.Include) != 1 || m.Include[0] != "**/*.go" {
		t.Fatalf("include=%v", m.Include)
	}
}

func TestCLINoFailOverridesConfigFailOn(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "codehound.toml")
	if err := os.WriteFile(path, []byte(`[codehound]
fail_on = "high"
`), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := config.LoadAndMerge(config.MergeInput{
		ConfigPath: path,
		NoFail:     true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if m.FailPolicy != nil {
		t.Fatalf("expected nil fail policy when --no-fail, got %v", *m.FailPolicy)
	}
}

func TestUnknownFieldRejected(t *testing.T) {
	_, err := config.Parse([]byte(`[codehound]
not_a_real_field = true
`))
	if err == nil {
		t.Fatal("expected unknown field error")
	}
}

func TestDiscover(t *testing.T) {
	root := t.TempDir()
	sub := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := filepath.Join(root, "codehound.toml")
	if err := os.WriteFile(cfg, []byte("[codehound]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got := config.Discover(sub)
	if got != cfg {
		t.Fatalf("Discover=%q want %q", got, cfg)
	}
}
