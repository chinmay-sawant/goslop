package engine

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/chinmay-sawant/goslop/internal/core"
)

// ScanEntry is a source file queued for analysis.
type ScanEntry struct {
	Path     string
	Language core.LanguageID
}

// WalkOptions control filesystem discovery.
//
// Skip counting matches Rust `ignore::WalkBuilder` + standard_filters:
// gitignored and hidden paths are not visited (and never increment skipped);
// only yielded files that fail language/extension/test filters count as skipped.
type WalkOptions struct {
	// IncludeTests includes *_test.go files (default false).
	IncludeTests bool
	// SkipVendor skips directories named "vendor" (default true).
	SkipVendor bool
	// Extensions restricts extensions (without dot). Empty means use registry plugins.
	Extensions []string
	// HonorIgnoreFiles applies .gitignore / .goslopignore / .ignore (default true).
	HonorIgnoreFiles bool
	// SkipCacheDir skips .goslop-cache directories (default true).
	SkipCacheDir bool
	// SkipHidden skips dotfiles and dot-directories (Rust standard_filters; default true).
	SkipHidden bool
	// Include is optional gitignore-style allow-list (empty = all paths allowed).
	Include []string
	// Exclude is optional gitignore-style deny-list.
	Exclude []string
}

// DefaultWalkOptions returns production walk defaults.
func DefaultWalkOptions() WalkOptions {
	return WalkOptions{
		IncludeTests:     false,
		SkipVendor:       true,
		Extensions:       []string{"go"},
		HonorIgnoreFiles: true,
		SkipCacheDir:     true,
		SkipHidden:       true,
	}
}

// CollectGoFiles walks roots and returns .go files suitable for analysis.
//
// Skips *_test.go unless opts.IncludeTests; skips vendor/ when opts.SkipVendor.
// Honors .gitignore / .goslopignore / .ignore when opts.HonorIgnoreFiles.
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
// plus a count of regular files skipped by test / extension / path filters.
//
// Parity with Rust engine::walk::FilesystemWalker:
//   - standard_filters: gitignore + hidden paths are omitted from the walk
//     entirely (not counted as skipped)
//   - skipped = files the walker yields that fail plugin/extension, language,
//     test, or include/exclude filters
//
// pluginExt maps extension (no dot) → language. When nil, opts.Extensions are
// treated as Go (LangGo).
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
				if pathAllowed(root, filepath.Dir(root), opts) {
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
		// Per-root include/exclude (gitignore-style), relative to this root.
		pathMatcher := newRootPathMatcher(rootAbs, opts.Include, opts.Exclude, !opts.IncludeTests)

		err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			rel, relErr := filepath.Rel(rootAbs, path)
			if relErr != nil {
				rel = path
			}
			rel = filepath.ToSlash(rel)
			name := d.Name()

			if d.IsDir() {
				if name == ".git" || name == ".hg" || name == ".svn" {
					return filepath.SkipDir
				}
				if skipVendor && name == "vendor" {
					return filepath.SkipDir
				}
				if opts.SkipCacheDir && name == ".goslop-cache" {
					return filepath.SkipDir
				}
				// Rust standard_filters: skip hidden directories entirely.
				if opts.SkipHidden && rel != "." && strings.HasPrefix(name, ".") {
					return filepath.SkipDir
				}
				if rel != "." && matcher != nil && matcher.ignored(rel, true) {
					return filepath.SkipDir
				}
				return nil
			}

			// Hidden files: omitted like Rust (not counted as skipped).
			if opts.SkipHidden && strings.HasPrefix(name, ".") {
				return nil
			}
			// Gitignored / goslop-ignored: omitted (not counted as skipped).
			if matcher != nil && rel != "." && matcher.ignored(rel, false) {
				return nil
			}
			// Path include/exclude + test-name filter (Rust RootPathMatcher).
			if pathMatcher != nil && !pathMatcher.allows(path) {
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
	// When pathMatcher already applies exclude_tests, this is a second guard for
	// CollectGoFiles callers that do not use pathMatcher include/exclude.
	if !opts.IncludeTests && strings.HasSuffix(base, "_test.go") {
		return ScanEntry{}, false
	}
	return ScanEntry{Path: path, Language: lang}, true
}

func pathAllowed(path, root string, opts WalkOptions) bool {
	m := newRootPathMatcher(root, opts.Include, opts.Exclude, !opts.IncludeTests)
	if m == nil {
		return true
	}
	return m.allows(path)
}

// rootPathMatcher mirrors Rust RootPathMatcher (include / exclude / tests).
type rootPathMatcher struct {
	include      *pathIgnoreMatcher
	exclude      *pathIgnoreMatcher
	excludeTests bool
	root         string
}

func newRootPathMatcher(root string, include, exclude []string, excludeTests bool) *rootPathMatcher {
	// Test exclusion is also enforced in acceptFile; path matcher is only
	// needed when include/exclude globs are configured.
	if len(include) == 0 && len(exclude) == 0 {
		return nil
	}
	m := &rootPathMatcher{root: root, excludeTests: excludeTests}
	if len(include) > 0 {
		m.include = matcherFromPatterns(include)
	}
	if len(exclude) > 0 {
		m.exclude = matcherFromPatterns(exclude)
	}
	return m
}

func matcherFromPatterns(patterns []string) *pathIgnoreMatcher {
	m := &pathIgnoreMatcher{}
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if p == "" || strings.HasPrefix(p, "#") {
			continue
		}
		m.patterns = append(m.patterns, parseIgnorePattern(p))
	}
	if len(m.patterns) == 0 {
		return nil
	}
	return m
}

func (m *rootPathMatcher) allows(path string) bool {
	if m == nil {
		return true
	}
	base := filepath.Base(path)
	if m.excludeTests && strings.Contains(base, "_test") {
		return false
	}
	rel := path
	if m.root != "" {
		if r, err := filepath.Rel(m.root, path); err == nil {
			rel = filepath.ToSlash(r)
		}
	}
	rel = filepath.ToSlash(rel)
	// include: must match at least one pattern when set (gitignore "whitelist" via matched-as-ignore).
	if m.include != nil && !m.include.ignored(rel, false) {
		// Our matcher.ignored means "matches a positive pattern". For include
		// lists we want "matches any pattern" — reuse ignored() which is true
		// when a non-negated pattern matches.
		return false
	}
	if m.exclude != nil && m.exclude.ignored(rel, false) {
		return false
	}
	return true
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
