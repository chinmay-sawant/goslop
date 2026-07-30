package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveScanScopeUsesDirectoryAndFileParent(t *testing.T) {
	dir := t.TempDir()
	dirScope := ResolveScanScope([]string{dir})
	if dirScope.ProjectRoot() != dir {
		t.Fatalf("directory root=%q want %q", dirScope.ProjectRoot(), dir)
	}
	path := filepath.Join(dir, "nested", "file.go")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("package p\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := dirScope.CacheKey(path); got != "nested/file.go" {
		t.Fatalf("cache key=%q", got)
	}
	fileScope := ResolveScanScope([]string{path})
	if fileScope.ProjectRoot() != filepath.Dir(path) {
		t.Fatalf("file root=%q want %q", fileScope.ProjectRoot(), filepath.Dir(path))
	}
	if got := fileScope.CacheKey(path); got != "file.go" {
		t.Fatalf("file cache key=%q", got)
	}
}

func TestScanScopeKeepsDistinctMultiRootKeysAndDisplayPaths(t *testing.T) {
	first := t.TempDir()
	second := t.TempDir()
	scope := ResolveScanScope([]string{first, second})
	if got := scope.CacheKey(filepath.Join(second, "other.go")); got == "other.go" {
		t.Fatal("second root must not collide with a first-root file cache key")
	}
	if got := scope.DisplayPath(filepath.Join(second, "other.go")); got == "" {
		t.Fatal("display path must not be empty")
	}
}
