package app

import (
	"os"
	"path/filepath"

	"github.com/chinmay/goslop/internal/engine/cache"
)

// scanScope defines the CLI path identity used by cache operations. The first
// path is the cache root; individual files use their parent directory.
type scanScope struct {
	projectRoot string
}

func resolveScanScope(paths []string) scanScope {
	root := "."
	if len(paths) == 0 {
		return scanScope{projectRoot: root}
	}

	path := paths[0]
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		if abs, absErr := filepath.Abs(path); absErr == nil {
			root = abs
		} else {
			root = path
		}
		return scanScope{projectRoot: root}
	}

	if abs, err := filepath.Abs(filepath.Dir(path)); err == nil {
		root = abs
	}
	return scanScope{projectRoot: root}
}

func (s scanScope) cacheKey(path string) string {
	rel, err := filepath.Rel(s.projectRoot, path)
	if err != nil {
		rel = path
	}
	return cache.NormalizePath(rel)
}
