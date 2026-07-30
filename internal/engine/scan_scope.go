package engine

import (
	"os"
	"path/filepath"

	"github.com/chinmay-sawant/goslop/internal/engine/cache"
)

// ScanScope owns the path identity shared by CLI cache maintenance and the
// analyzer: the first requested root establishes cache-relative keys, while
// display paths remain relative to the current working directory when safe.
type ScanScope struct {
	projectRoot string
}

// ResolveScanScope resolves the first requested directory, or a requested
// file's parent, into the cache root used for this scan.
func ResolveScanScope(paths []string) ScanScope {
	if len(paths) == 0 {
		return ScanScope{projectRoot: "."}
	}
	path := paths[0]
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return NewScanScope(path)
	}
	return NewScanScope(filepath.Dir(path))
}

// NewScanScope creates a scope for an already-resolved cache root.
func NewScanScope(root string) ScanScope {
	if root == "" {
		root = "."
	}
	if abs, err := filepath.Abs(root); err == nil {
		root = abs
	}
	return ScanScope{projectRoot: root}
}

// ProjectRoot is the root used for cache-relative keys.
func (s ScanScope) ProjectRoot() string { return s.projectRoot }

// CacheKey returns the normalized cache key for path.
func (s ScanScope) CacheKey(path string) string {
	rel, err := filepath.Rel(s.projectRoot, path)
	if err != nil {
		rel = path
	}
	return cache.NormalizePath(rel)
}

// DisplayPath returns a path relative to the current directory when it stays
// within that directory; otherwise it preserves the supplied path.
func (s ScanScope) DisplayPath(path string) string {
	if rel, err := filepath.Rel(".", path); err == nil && !filepath.IsAbs(rel) {
		return rel
	}
	return path
}
