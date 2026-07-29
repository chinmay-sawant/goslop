package engine

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// pathIgnoreMatcher is a minimal gitignore-style matcher for walk filtering.
// Supports: comments, blank lines, negation (!), directory trailing /, simple
// globs with * (path.Match), and path prefixes. Not a full gitignore engine.
type pathIgnoreMatcher struct {
	// patterns are evaluated in order; last match wins for negation.
	patterns []ignorePattern
}

type ignorePattern struct {
	raw      string
	negated  bool
	dirOnly  bool
	matchAll bool // "**" or empty-ish
	segments []string
	hasSlash bool
}

// loadPathIgnores loads .gitignore, .codehoundignore, and .ignore from root.
func loadPathIgnores(root string) *pathIgnoreMatcher {
	m := &pathIgnoreMatcher{}
	for _, name := range []string{".gitignore", ".codehoundignore", ".ignore"} {
		path := filepath.Join(root, name)
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			m.patterns = append(m.patterns, parseIgnorePattern(line))
		}
		_ = f.Close()
	}
	if len(m.patterns) == 0 {
		return nil
	}
	return m
}

func parseIgnorePattern(line string) ignorePattern {
	p := ignorePattern{raw: line}
	if strings.HasPrefix(line, "!") {
		p.negated = true
		line = line[1:]
	}
	// Strip leading slash (repo-root relative); treat as relative to walk root.
	line = strings.TrimPrefix(line, "/")
	if strings.HasSuffix(line, "/") {
		p.dirOnly = true
		line = strings.TrimSuffix(line, "/")
	}
	line = filepath.ToSlash(line)
	p.hasSlash = strings.Contains(line, "/")
	if line == "**" || line == "*" {
		p.matchAll = true
		return p
	}
	p.segments = strings.Split(line, "/")
	return p
}

// ignored reports whether relPath (slash-separated, relative to root) should be skipped.
// isDir is true when path is a directory.
func (m *pathIgnoreMatcher) ignored(relPath string, isDir bool) bool {
	if m == nil || len(m.patterns) == 0 {
		return false
	}
	relPath = filepath.ToSlash(relPath)
	relPath = strings.TrimPrefix(relPath, "./")
	matched := false
	for _, p := range m.patterns {
		if p.dirOnly && !isDir {
			// Directory patterns only skip when the path itself is a dir,
			// but also when any parent matches... handled by walk SkipDir.
			// For files, still match if a path component equals the dir name.
			if !matchPattern(p, relPath, false) {
				continue
			}
			// File under ignored directory name: treat as match when pattern is a prefix.
		}
		if matchPattern(p, relPath, isDir) {
			if p.negated {
				matched = false
			} else {
				matched = true
			}
		}
	}
	return matched
}

func matchPattern(p ignorePattern, relPath string, isDir bool) bool {
	if p.matchAll {
		return true
	}
	if p.dirOnly && !isDir {
		// Match files whose path is under the directory prefix.
		prefix := strings.Join(p.segments, "/")
		if prefix == "" {
			return false
		}
		return relPath == prefix || strings.HasPrefix(relPath, prefix+"/")
	}

	// Pattern without slash: match against any path segment (basename) or full path.
	if !p.hasSlash {
		pat := p.segments[0]
		base := filepath.Base(relPath)
		if matchGlob(pat, base) {
			return true
		}
		// Also match full relative path for convenience.
		return matchGlob(pat, relPath)
	}

	// Pattern with slash: match full relative path (and ** style loosely).
	pat := strings.Join(p.segments, "/")
	if matchGlob(pat, relPath) {
		return true
	}
	// Prefix directory match: "foo/bar" matches "foo/bar/baz.go"
	if strings.HasPrefix(relPath, pat+"/") {
		return true
	}
	return false
}

// matchGlob supports * wildcards via path.Match-like rules on each segment.
func matchGlob(pattern, name string) bool {
	// Fast paths
	if pattern == name {
		return true
	}
	if pattern == "*" {
		return true
	}
	// Use filepath.Match which treats * as within a single segment (no /).
	// For multi-segment patterns, match segment-wise when both have same count,
	// else fall back to full-string Match with * only (no ** expansion beyond prefix).
	ok, err := filepath.Match(pattern, name)
	if err == nil && ok {
		return true
	}
	// Simple trailing "/*"
	if strings.HasSuffix(pattern, "/*") {
		prefix := strings.TrimSuffix(pattern, "/*")
		if name == prefix {
			return false
		}
		if strings.HasPrefix(name, prefix+"/") {
			rest := name[len(prefix)+1:]
			return !strings.Contains(rest, "/")
		}
	}
	// Simple "**/" prefix
	if strings.HasPrefix(pattern, "**/") {
		suffix := pattern[3:]
		if matchGlob(suffix, name) {
			return true
		}
		if idx := strings.Index(name, "/"); idx >= 0 {
			return matchGlob(pattern, name[idx+1:])
		}
	}
	return false
}
