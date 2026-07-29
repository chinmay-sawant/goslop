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
}

// DefaultWalkOptions returns production walk defaults.
func DefaultWalkOptions() WalkOptions {
	return WalkOptions{
		IncludeTests: false,
		SkipVendor:   true,
		Extensions:   []string{"go"},
	}
}

// CollectGoFiles walks roots and returns .go files suitable for analysis.
//
// Skips *_test.go unless opts.IncludeTests; skips vendor/ when opts.SkipVendor.
// Basic .gitignore is not applied yet (optional stub).
func CollectGoFiles(roots []string, opts WalkOptions) ([]string, error) {
	opts.Extensions = []string{"go"}
	entries, err := CollectFiles(roots, opts, nil)
	if err != nil {
		return nil, err
	}
	out := make([]string, len(entries))
	for i, e := range entries {
		out[i] = e.Path
	}
	return out, nil
}

// CollectFiles walks roots and returns ScanEntry values for supported extensions.
//
// pluginExt maps extension (no dot) → language. When nil, opts.Extensions are
// treated as Go (LangGo) for MVP simplicity.
//
// .gitignore / .codehoundignore respect is a future hook (Phase 10 / walk polish).
func CollectFiles(roots []string, opts WalkOptions, pluginExt map[string]core.LanguageID) ([]ScanEntry, error) {
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
	// Default SkipVendor true when using DefaultWalkOptions; zero-value leaves it false.
	// Callers that want vendor skipped should set SkipVendor or use DefaultWalkOptions.

	var entries []ScanEntry
	seen := make(map[string]struct{})

	for _, root := range roots {
		if root == "" {
			root = "."
		}
		info, err := os.Stat(root)
		if err != nil {
			return nil, err
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
			}
			continue
		}

		err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() {
				name := d.Name()
				if name == ".git" || name == ".hg" || name == ".svn" {
					return filepath.SkipDir
				}
				if skipVendor && name == "vendor" {
					return filepath.SkipDir
				}
				return nil
			}
			e, ok := acceptFile(path, opts, pluginExt)
			if !ok {
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
			return nil, err
		}
	}
	return entries, nil
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
