package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// ContentHash returns "sha256:<hex>" of source text.
func ContentHash(source string) string {
	sum := sha256.Sum256([]byte(source))
	return "sha256:" + hex.EncodeToString(sum[:])
}

// CacheKeyForPath returns the lowercase hex filename for a relative path entry.
func CacheKeyForPath(relPath string) string {
	norm := NormalizePath(relPath)
	sum := sha256.Sum256([]byte(norm))
	return hex.EncodeToString(sum[:])
}

// NormalizePath normalizes project-relative paths for cache keys (forward slashes, no leading ./).
func NormalizePath(p string) string {
	p = strings.ReplaceAll(p, "\\", "/")
	p = strings.TrimPrefix(p, "./")
	for strings.HasPrefix(p, "/") {
		p = strings.TrimPrefix(p, "/")
	}
	return p
}
