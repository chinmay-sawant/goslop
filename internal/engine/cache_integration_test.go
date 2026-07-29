package engine_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/engine"
	"github.com/chinmay/goslop/internal/engine/cache"
	"github.com/chinmay/goslop/internal/engine/ignore"
	"github.com/chinmay/goslop/internal/rules"
)

func TestAnalyzerCacheHitMiss(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "vulnerable.go", `package main

import "os/exec"

func run(cmd string) {
	exec.Command("sh", "-c", cmd)
}
`)

	reg, err := engine.NewRegistry([]core.LanguagePlugin{&goTestPlugin{}})
	if err != nil {
		t.Fatal(err)
	}
	ctx := core.DefaultScanContext()
	ctx.BadPracticesEnabled = true

	store := cache.InMemory("test-ver")
	a := engine.NewAnalyzerBuilder().
		Registry(reg).
		ScanContext(ctx).
		Cache(store).
		ProjectRoot(dir).
		Build()

	res1, err := a.AnalyzePaths([]string{dir})
	if err != nil {
		t.Fatal(err)
	}
	if res1.Stats.CacheMisses < 1 {
		t.Fatalf("expected cache miss on first scan, stats=%+v", res1.Stats)
	}
	if len(res1.Findings) == 0 {
		t.Fatal("expected findings on first scan")
	}

	// Second scan should hit.
	res2, err := a.AnalyzePaths([]string{dir})
	if err != nil {
		t.Fatal(err)
	}
	if res2.Stats.CacheHits < 1 {
		t.Fatalf("expected cache hit on second scan, stats=%+v", res2.Stats)
	}
	if len(res2.Findings) != len(res1.Findings) {
		t.Fatalf("findings mismatch hit=%d miss=%d", len(res2.Findings), len(res1.Findings))
	}

	// Content change → miss again
	writeFile(t, dir, "vulnerable.go", `package main

import "os/exec"

func run(cmd string) {
	// changed
	exec.Command("sh", "-c", cmd)
}
`)
	res3, err := a.AnalyzePaths([]string{dir})
	if err != nil {
		t.Fatal(err)
	}
	if res3.Stats.CacheMisses < 1 {
		t.Fatalf("expected miss after edit, stats=%+v", res3.Stats)
	}
}

func TestAnalyzerInlineIgnoreSuppresses(t *testing.T) {
	dir := t.TempDir()
	// Materialize fixture-style suppressed inline
	body := `package main

import (
	"net/http"
	"os/exec"
)

func handler(w http.ResponseWriter, r *http.Request) {
	cmd := r.URL.Query().Get("cmd")
	// goslop-ignore: TEST-EXEC-COMMAND
	exec.Command("sh", "-c", cmd).Run()
}
`
	writeFile(t, dir, "suppressed_inline.go", body)

	reg, err := engine.NewRegistry([]core.LanguagePlugin{&goTestPlugin{}})
	if err != nil {
		t.Fatal(err)
	}
	ctx := core.DefaultScanContext()
	// No cache to isolate ignore behaviour
	ctx.NoCache = true
	a := engine.NewAnalyzer(ctx, reg)
	res, err := a.AnalyzePaths([]string{dir})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 0 {
		t.Fatalf("expected ignore suppression, got %+v", res.Findings)
	}
	if res.SuppressedCount < 1 {
		t.Fatalf("expected suppressed count, got %d", res.SuppressedCount)
	}
}

func TestWalkHonorsGoslopIgnore(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "keep.go", "package main\n")
	writeFile(t, dir, "skip_me.go", "package main\n")
	if err := os.WriteFile(filepath.Join(dir, ".goslopignore"), []byte("skip_me.go\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	entries, err := engine.CollectGoFiles([]string{dir}, engine.DefaultWalkOptions())
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range entries {
		if strings.Contains(p, "skip_me.go") {
			t.Fatalf("skip_me.go should be ignored, got %v", entries)
		}
	}
	if len(entries) != 1 {
		t.Fatalf("entries: %v", entries)
	}
}

func TestIgnorePackageUnit(t *testing.T) {
	// Sanity: package is importable from engine tests
	if !ignore.All().IsAll() {
		t.Fatal()
	}
	_ = rules.SeverityHigh
}
