package engine

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/chinmay/codehound/internal/core"
)

// ScanEntry is a source file queued for analysis.
type ScanEntry struct {
	Path     string
	Language core.LanguageID
}

// WalkOptions control filesystem discovery.
type WalkOptions struct {
	// IncludeTests includes *_test.go files (default false).
	IncludeTests bool
	// SkipVendor skips directories named "vendor" (default true).
	SkipVendor bool
	// Extensions restricts extensions (without dot). Empty means use registry plugins.
	Extensions []string
	// HonorIgnoreFiles applies .gitignore / .codehoundignore / .ignore (default true).
	HonorIgnoreFiles bool
	// SkipCacheDir skips .codehound-cache directories (default true).
	SkipCacheDir bool
}

// DefaultWalkOptions returns production walk defaults.
func DefaultWalkOptions() WalkOptions {
	return WalkOptions{
		IncludeTests:     false,
		SkipVendor:       true,
		Extensions:       []string{"go"},
		HonorIgnoreFiles: true,
		SkipCacheDir:     true,
	}
}

// CollectGoFiles walks roots and returns .go files suitable for analysis.
//
// Skips *_test.go unless opts.IncludeTests; skips vendor/ when opts.SkipVendor.
// Honors .gitignore / .codehoundignore / .ignore when opts.HonorIgnoreFiles.
func CollectGoFiles(roots []string, opts WalkOptions) ([]string, error) {
	opts.Extensions = []string{"go"}
	entries, _, err := CollectFiles(roots, opts, nil)
	if err != nil {
		return nil, err
	}
	out := make([]string, len(entries))
	for i, e := range entries {
		out[i] = e.Path
	}
	return out, nil
}

// CollectFiles walks roots and returns ScanEntry values for supported extensions
// plus a count of regular files skipped by ignore / extension / test filters
// (Rust walk parity: every non-accepted file increments skipped).
//
// pluginExt maps extension (no dot) → language. When nil, opts.Extensions are
// treated as Go (LangGo) for MVP simplicity.
//
// Honors .gitignore / .codehoundignore / .ignore when opts.HonorIgnoreFiles
// (default true via DefaultWalkOptions). Vendor and VCS dirs are always skipped
// when SkipVendor is set.
func CollectFiles(roots []string, opts WalkOptions, pluginExt map[string]core.LanguageID) ([]ScanEntry, int, error) {
	if pluginExt == nil {
		pluginExt = map[string]core.LanguageID{}
		exts := opts.Extensions
		if len(exts) == 0 {
			exts = []string{"go"}
		}
		for _, ext := range exts {
			pluginExt[ext] = core.LangGo
		}
	}

	skipVendor := opts.SkipVendor

	var entries []ScanEntry
	skipped := 0
	seen := make(map[string]struct{})

	for _, root := range roots {
		if root == "" {
			root = "."
		}
		info, err := os.Stat(root)
		if err != nil {
			return nil, 0, err
		}
		if !info.IsDir() {
			if e, ok := acceptFile(root, opts, pluginExt); ok {
				abs, aerr := filepath.Abs(root)
				if aerr != nil {
					abs = root
				}
				if _, dup := seen[abs]; !dup {
					seen[abs] = struct{}{}
					e.Path = abs
					entries = append(entries, e)
				}
			} else {
				skipped++
			}
			continue
		}

		var matcher *pathIgnoreMatcher
		if opts.HonorIgnoreFiles {
			matcher = loadPathIgnores(root)
		}
		rootAbs, rerr := filepath.Abs(root)
		if rerr != nil {
			rootAbs = root
		}

		err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			rel, relErr := filepath.Rel(rootAbs, path)
			if relErr != nil {
				rel = path
			}
			rel = filepath.ToSlash(rel)

			if d.IsDir() {
				name := d.Name()
				if name == ".git" || name == ".hg" || name == ".svn" {
					return filepath.SkipDir
				}
				if skipVendor && name == "vendor" {
					return filepath.SkipDir
				}
				if opts.SkipCacheDir && name == ".codehound-cache" {
					return filepath.SkipDir
				}
				// Do not apply ignore to the walk root itself.
				if rel != "." && matcher != nil && matcher.ignored(rel, true) {
					return filepath.SkipDir
				}
				return nil
			}
			// Regular file: count skips like Rust (ignore / unsupported / tests).
			if matcher != nil && rel != "." && matcher.ignored(rel, false) {
				skipped++
				return nil
			}
			e, ok := acceptFile(path, opts, pluginExt)
			if !ok {
				skipped++
				return nil
			}
			abs, aerr := filepath.Abs(path)
			if aerr != nil {
				abs = path
			}
			if _, dup := seen[abs]; dup {
				return nil
			}
			seen[abs] = struct{}{}
			e.Path = abs
			entries = append(entries, e)
			return nil
		})
		if err != nil {
			return nil, 0, err
		}
	}
	return entries, skipped, nil
}

func acceptFile(path string, opts WalkOptions, pluginExt map[string]core.LanguageID) (ScanEntry, bool) {
	ext := strings.TrimPrefix(filepath.Ext(path), ".")
	lang, ok := pluginExt[ext]
	if !ok {
		return ScanEntry{}, false
	}
	base := filepath.Base(path)
	if !opts.IncludeTests && strings.HasSuffix(base, "_test.go") {
		return ScanEntry{}, false
	}
	return ScanEntry{Path: path, Language: lang}, true
}

// ReadUTF8 reads path and verifies valid UTF-8.
func ReadUTF8(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	if !utf8.Valid(b) {
		return "", &ScanError{
			Path:    path,
			Kind:    ScanErrorEncoding,
			Message: "source is not valid UTF-8",
		}
	}
	return string(b), nil
}
