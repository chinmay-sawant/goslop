package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/config"
	"github.com/chinmay-sawant/goslop/internal/core"
)

func TestParseAndMergeAdditiveOnlySkip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "goslop.toml")
	body := `
[goslop]
only = ["PERF-1"]
skip = ["BP-1"]
fail_on = "high"
exclude_tests = false
include = ["**/*.go"]
exclude = ["vendor/**"]

[goslop.cache]
enabled = false

[goslop.taint]
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
	path := filepath.Join(dir, "goslop.toml")
	if err := os.WriteFile(path, []byte(`[goslop]
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
	_, err := config.Parse([]byte(`[goslop]
not_a_real_field = true
`))
	if err == nil {
		t.Fatal("expected unknown field error")
	}
}

func TestUnsupportedLanguageAndTypedConfigurationRejected(t *testing.T) {
	for name, body := range map[string]string{
		"languages": `[goslop]
languages = ["go"]
`,
		"typed": `[goslop.typed]
enabled = true
`,
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := config.Parse([]byte(body)); err == nil {
				t.Fatalf("expected unsupported %s configuration to be rejected", name)
			}
		})
	}
}

func TestEvictTargetRatioValidatedAndMerged(t *testing.T) {
	const body = `[goslop.cache]
evict_target_ratio = 0.5
`
	dir := t.TempDir()
	path := filepath.Join(dir, "goslop.toml")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	merged, err := config.LoadAndMerge(config.MergeInput{ConfigPath: path})
	if err != nil {
		t.Fatal(err)
	}
	if merged.CacheEvictTargetRatio == nil || *merged.CacheEvictTargetRatio != 0.5 {
		t.Fatalf("ratio=%v", merged.CacheEvictTargetRatio)
	}
	if _, err := config.Parse([]byte(`[goslop.cache]
evict_target_ratio = 0.05
`)); err == nil {
		t.Fatal("expected out-of-range eviction target ratio to be rejected")
	}
}

func TestExportWholeFunctionConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "goslop.toml")
	if err := os.WriteFile(path, []byte(`[goslop.export]
whole_function = false
`), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := config.LoadAndMerge(config.MergeInput{ConfigPath: path})
	if err != nil {
		t.Fatal(err)
	}
	if m.ExportWholeFunction == nil || *m.ExportWholeFunction {
		t.Fatalf("expected whole_function=false, got %v", m.ExportWholeFunction)
	}

	// Unset → nil (export package defaults to true).
	path2 := filepath.Join(dir, "empty.toml")
	if writeErr := os.WriteFile(path2, []byte(`[goslop]
fail_on = "none"
`), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}
	m2, err := config.LoadAndMerge(config.MergeInput{ConfigPath: path2})
	if err != nil {
		t.Fatal(err)
	}
	if m2.ExportWholeFunction != nil {
		t.Fatalf("expected nil ExportWholeFunction when unset, got %v", *m2.ExportWholeFunction)
	}
}

func TestDiscover(t *testing.T) {
	root := t.TempDir()
	sub := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := filepath.Join(root, "goslop.toml")
	if err := os.WriteFile(cfg, []byte("[goslop]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got := config.Discover(sub)
	if got != cfg {
		t.Fatalf("Discover=%q want %q", got, cfg)
	}
}
