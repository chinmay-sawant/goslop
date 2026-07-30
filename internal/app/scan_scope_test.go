package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveScanScopeUsesDirectoryAsProjectRoot(t *testing.T) {
	dir := t.TempDir()
	scope := resolveScanScope([]string{dir})
	if scope.projectRoot != dir {
		t.Fatalf("project root=%q want %q", scope.projectRoot, dir)
	}
	if got := scope.cacheKey(filepath.Join(dir, "nested", "file.go")); got != "nested/file.go" {
		t.Fatalf("cache key=%q", got)
	}
}

func TestResolveScanScopeUsesFileParentAsProjectRoot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.go")
	if err := os.WriteFile(path, []byte("package p\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	scope := resolveScanScope([]string{path})
	if scope.projectRoot != dir {
		t.Fatalf("project root=%q want %q", scope.projectRoot, dir)
	}
	if got := scope.cacheKey(path); got != "file.go" {
		t.Fatalf("cache key=%q", got)
	}
}
