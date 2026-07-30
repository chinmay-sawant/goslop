package cache

import (
	"fmt"

	"github.com/chinmay-sawant/goslop/internal/rules"
)

// CACHE_VERSION is the on-disk schema version. Bump on breaking JSON shape changes.
const CACHE_VERSION uint32 = 2

// DEFAULT_CACHE_DIR is the conventional cache directory name.
const DEFAULT_CACHE_DIR = ".goslop-cache"

const (
	manifestName = "manifest.json"
	filesSubdir  = "files"
)

// Manifest is the cheap O(1) lookup index for per-file cache state.
type Manifest struct {
	SchemaVersion  uint32                   `json:"schema_version"`
	ToolVersion    string                   `json:"tool_version"`
	RuleConfigHash string                   `json:"rule_config_hash"`
	Files          map[string]FileCacheMeta `json:"files"`
}

// FileCacheMeta is per-file metadata stored in the manifest (no findings).
type FileCacheMeta struct {
	ContentHash  string   `json:"content_hash"`
	Dependencies []string `json:"dependencies,omitempty"`
	CachedAt     string   `json:"cached_at"`
}

// Entry is the full per-file cache payload at files/<key>.json.
type Entry struct {
	SchemaVersion   uint32          `json:"schema_version"`
	File            string          `json:"file"`
	ContentHash     string          `json:"content_hash"`
	RuleConfigHash  string          `json:"rule_config_hash"`
	Findings        []rules.Finding `json:"findings"`
	SuppressedCount int             `json:"suppressed_count"`
	CachedAt        string          `json:"cached_at"`
}

// Lookup is the outcome of a cache lookup.
type Lookup int

const (
	// LookupMiss means the file has no manifest entry.
	LookupMiss Lookup = iota
	// LookupStale means the entry exists but is invalid.
	LookupStale
	// LookupHit means a fresh entry was returned.
	LookupHit
)

// String returns a short name for diagnostics.
func (l Lookup) String() string {
	switch l {
	case LookupHit:
		return "hit"
	case LookupStale:
		return "stale"
	case LookupMiss:
		return "miss"
	default:
		return "miss"
	}
}

// Error is a cache-layer failure (never panics on corrupt entries).
type Error struct {
	Op   string
	Path string
	Err  error
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Path != "" {
		return fmt.Sprintf("cache %s %s: %v", e.Op, e.Path, e.Err)
	}
	return fmt.Sprintf("cache %s: %v", e.Op, e.Err)
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}
