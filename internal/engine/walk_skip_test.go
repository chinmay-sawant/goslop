package engine_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/engine"
)

func TestCollectFiles_SkipParityHiddenAndIgnored(t *testing.T) {
	dir := t.TempDir()
	// accepted
	mustWrite(t, filepath.Join(dir, "main.go"), "package main\n")
	// test → skipped count
	mustWrite(t, filepath.Join(dir, "main_test.go"), "package main\n")
	// non-go → skipped count
	mustWrite(t, filepath.Join(dir, "README.md"), "# hi\n")
	// hidden file → not counted
	mustWrite(t, filepath.Join(dir, ".env"), "X=1\n")
	// hidden dir → not walked
	mustWrite(t, filepath.Join(dir, ".cache", "x.txt"), "x\n")
	// gitignored file → not counted
	mustWrite(t, filepath.Join(dir, ".gitignore"), "secret.bin\nignored/\n")
	mustWrite(t, filepath.Join(dir, "secret.bin"), "bin\n")
	mustWrite(t, filepath.Join(dir, "ignored", "a.go"), "package ignored\n")

	entries, skipped, err := engine.CollectFiles([]string{dir}, engine.DefaultWalkOptions(), map[string]core.LanguageID{"go": core.LangGo})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("entries=%d %#v", len(entries), entries)
	}
	// skipped: main_test.go + README.md only (not .env, not secret.bin, not ignored/)
	if skipped != 2 {
		t.Fatalf("skipped=%d want 2", skipped)
	}
}

func TestCollectFiles_IncludeExclude(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, "keep", "a.go"), "package keep\n")
	mustWrite(t, filepath.Join(dir, "drop", "b.go"), "package drop\n")
	opts := engine.DefaultWalkOptions()
	opts.Include = []string{"keep/*"}
	entries, _, err := engine.CollectFiles([]string{dir}, opts, map[string]core.LanguageID{"go": core.LangGo})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || !strings.Contains(entries[0].Path, "keep") {
		t.Fatalf("entries=%v", entries)
	}
}

func mustWrite(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}
