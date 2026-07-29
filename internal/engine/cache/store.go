package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/chinmay/codehound/internal/rules"
)

// Store is an on-disk incremental analysis cache.
//
// Layout:
//
//	<codehound-cache>/
//	  manifest.json
//	  files/<sha256(path)>.json
//
// Safe for concurrent Lookup; Put/Flush/Prune require exclusive use or the
// store's internal mutex (methods take the lock).
type Store struct {
	mu               sync.Mutex
	cacheDir         string
	filesDir         string
	manifest         Manifest
	dirty            bool
	toolVersion      string
	ruleConfigHash   string
	maxSizeBytes     int64
	evictTarget      float64 // fraction of max to retain after eviction
	maxFileSizeBytes int64
	ephemeral        bool
	memoryEntries    map[string]*Entry // ephemeral only
}

// OpenOptions configures Store.Open.
type OpenOptions struct {
	// MaxSizeMB is the max size of files/ in MiB. 0 disables eviction.
	MaxSizeMB uint64
	// EvictTargetRatio is the fraction of MaxSizeMB to retain after eviction (default 0.9).
	EvictTargetRatio float64
	// MaxFileSizeMB skips caching files larger than this. 0 disables. Default 4.
	MaxFileSizeMB uint64
	// ToolVersion is written into the manifest (required for version invalidation).
	ToolVersion string
}

// Open opens or creates a disk-backed cache at cacheDir.
func Open(cacheDir string, opts OpenOptions) (*Store, error) {
	if cacheDir == "" {
		cacheDir = DEFAULT_CACHE_DIR
	}
	if opts.ToolVersion == "" {
		opts.ToolVersion = "unknown"
	}
	evict := opts.EvictTargetRatio
	if evict <= 0 || evict >= 1 || !(evict >= 0.1 && evict <= 0.99) {
		evict = 0.9
	}
	// Default max file size 4 MiB when not specified.
	maxFileMB := opts.MaxFileSizeMB
	if maxFileMB == 0 {
		maxFileMB = 4
	}

	filesDir := filepath.Join(cacheDir, filesSubdir)
	if err := os.MkdirAll(filesDir, 0o755); err != nil {
		return nil, &Error{Op: "mkdir", Path: filesDir, Err: err}
	}

	s := &Store{
		cacheDir:         cacheDir,
		filesDir:         filesDir,
		toolVersion:      opts.ToolVersion,
		maxSizeBytes:     int64(opts.MaxSizeMB) * 1024 * 1024,
		evictTarget:      evict,
		maxFileSizeBytes: int64(maxFileMB) * 1024 * 1024,
	}

	manifestPath := filepath.Join(cacheDir, manifestName)
	if data, err := os.ReadFile(manifestPath); err == nil {
		var m Manifest
		if jerr := json.Unmarshal(data, &m); jerr != nil {
			// Corrupt → empty
			s.manifest = emptyManifest(opts.ToolVersion)
			s.dirty = true
		} else if m.SchemaVersion != CACHE_VERSION {
			// Incompatible schema → empty
			s.manifest = emptyManifest(opts.ToolVersion)
			s.dirty = true
		} else if m.ToolVersion != opts.ToolVersion {
			// Tool version change → mass-stale (empty manifest, drop entries)
			s.manifest = emptyManifest(opts.ToolVersion)
			s.dirty = true
			_ = clearFilesDir(filesDir)
		} else {
			if m.Files == nil {
				m.Files = make(map[string]FileCacheMeta)
			}
			s.manifest = m
			s.ruleConfigHash = m.RuleConfigHash
		}
	} else if !os.IsNotExist(err) {
		return nil, &Error{Op: "read", Path: manifestPath, Err: err}
	} else {
		s.manifest = emptyManifest(opts.ToolVersion)
		s.dirty = true
	}

	return s, nil
}

// InMemory returns an ephemeral store that never touches disk (tests).
func InMemory(toolVersion string) *Store {
	if toolVersion == "" {
		toolVersion = "test"
	}
	return &Store{
		manifest:         emptyManifest(toolVersion),
		toolVersion:      toolVersion,
		evictTarget:      0.9,
		ephemeral:        true,
		maxFileSizeBytes: 0,
	}
}

func emptyManifest(toolVersion string) Manifest {
	return Manifest{
		SchemaVersion: CACHE_VERSION,
		ToolVersion:   toolVersion,
		Files:         make(map[string]FileCacheMeta),
	}
}

// Dir returns the cache root directory.
func (s *Store) Dir() string {
	if s == nil {
		return ""
	}
	return s.cacheDir
}

// EnsureRuleConfigHash mass-stales entries when the finding-affecting config changes.
func (s *Store) EnsureRuleConfigHash(hash string) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.manifest.RuleConfigHash == hash {
		s.ruleConfigHash = hash
		return
	}
	if s.manifest.RuleConfigHash != "" || len(s.manifest.Files) > 0 {
		for file := range s.manifest.Files {
			s.removeLocked(file)
		}
	}
	s.manifest.RuleConfigHash = hash
	s.ruleConfigHash = hash
	s.dirty = true
}

// ShouldCacheBytes reports whether a file of sizeBytes is eligible for cache writes.
func (s *Store) ShouldCacheBytes(sizeBytes int64) bool {
	if s == nil {
		return false
	}
	if s.maxFileSizeBytes <= 0 {
		return true
	}
	return sizeBytes <= s.maxFileSizeBytes
}

// Lookup returns hit/stale/miss for file+contentHash.
func (s *Store) Lookup(file, contentHash string) (Lookup, *Entry) {
	if s == nil {
		return LookupMiss, nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	file = NormalizePath(file)
	meta, ok := s.manifest.Files[file]
	if !ok {
		return LookupMiss, nil
	}
	if meta.ContentHash != contentHash {
		return LookupStale, nil
	}
	entry, err := s.loadEntryLocked(file)
	if err != nil || entry == nil {
		return LookupStale, nil
	}
	if NormalizePath(entry.File) != file ||
		entry.ContentHash != meta.ContentHash ||
		entry.RuleConfigHash != s.manifest.RuleConfigHash {
		return LookupStale, nil
	}
	return LookupHit, entry
}

// Put stores findings for a file (overwrites existing).
func (s *Store) Put(file, contentHash string, findings []rules.Finding, suppressed int) error {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	file = NormalizePath(file)
	now := time.Now().UTC().Format(time.RFC3339)
	entry := Entry{
		SchemaVersion:   CACHE_VERSION,
		File:            file,
		ContentHash:     contentHash,
		RuleConfigHash:  s.manifest.RuleConfigHash,
		Findings:        cloneFindings(findings),
		SuppressedCount: suppressed,
		CachedAt:        now,
	}
	if entry.Findings == nil {
		entry.Findings = []rules.Finding{}
	}

	if !s.ephemeral {
		key := CacheKeyForPath(file)
		path := filepath.Join(s.filesDir, key+".json")
		data, err := json.MarshalIndent(entry, "", "  ")
		if err != nil {
			return &Error{Op: "marshal", Path: path, Err: err}
		}
		if err := os.WriteFile(path, data, 0o644); err != nil {
			return &Error{Op: "write", Path: path, Err: err}
		}
	} else {
		// Ephemeral: keep entry only in an in-memory side map via manifest meta + filesDir unused.
		// We store JSON in a private map by abusing filesDir memory — use dedicated field.
		if s.memoryEntries == nil {
			s.memoryEntries = make(map[string]*Entry)
		}
		cp := entry
		s.memoryEntries[file] = &cp
	}

	s.manifest.Files[file] = FileCacheMeta{
		ContentHash: contentHash,
		CachedAt:    now,
	}
	s.dirty = true
	return nil
}

// Prune removes entries for files not in scanned (normalized paths).
func (s *Store) Prune(scanned map[string]struct{}) (int, error) {
	if s == nil {
		return 0, nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	normScanned := make(map[string]struct{}, len(scanned))
	for f := range scanned {
		normScanned[NormalizePath(f)] = struct{}{}
	}
	var toRemove []string
	for f := range s.manifest.Files {
		if _, ok := normScanned[NormalizePath(f)]; !ok {
			toRemove = append(toRemove, f)
		}
	}
	for _, f := range toRemove {
		s.removeLocked(f)
	}
	return len(toRemove), nil
}

// CleanOrphans removes files/*.json not referenced by the manifest.
func (s *Store) CleanOrphans() (int, error) {
	if s == nil || s.ephemeral {
		return 0, nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	active := make(map[string]struct{}, len(s.manifest.Files))
	for f := range s.manifest.Files {
		active[CacheKeyForPath(f)+".json"] = struct{}{}
	}
	entries, err := os.ReadDir(s.filesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, &Error{Op: "readdir", Path: s.filesDir, Err: err}
	}
	removed := 0
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if _, ok := active[name]; ok {
			continue
		}
		path := filepath.Join(s.filesDir, name)
		if err := os.Remove(path); err != nil {
			return removed, &Error{Op: "remove", Path: path, Err: err}
		}
		removed++
	}
	return removed, nil
}

// Flush persists the manifest and runs size-based eviction when configured.
func (s *Store) Flush() error {
	if s == nil || s.ephemeral {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.maxSizeBytes > 0 {
		if err := s.evictLocked(); err != nil {
			return err
		}
	}
	if !s.dirty {
		return nil
	}
	return s.writeManifestLocked()
}

// Rebuild wipes the cache directory (after safety checks) and reopens empty.
func Rebuild(cacheDir string, opts OpenOptions) (*Store, error) {
	if cacheDir == "" {
		cacheDir = DEFAULT_CACHE_DIR
	}
	if err := validatePurgePath(cacheDir); err != nil {
		return nil, err
	}
	if err := os.RemoveAll(cacheDir); err != nil && !os.IsNotExist(err) {
		return nil, &Error{Op: "purge", Path: cacheDir, Err: err}
	}
	return Open(cacheDir, opts)
}

// EntryCount returns the number of tracked files.
func (s *Store) EntryCount() int {
	if s == nil {
		return 0
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.manifest.Files)
}

// --- internal ---

func (s *Store) loadEntryLocked(file string) (*Entry, error) {
	if s.ephemeral {
		if s.memoryEntries == nil {
			return nil, nil
		}
		e := s.memoryEntries[file]
		if e == nil {
			return nil, nil
		}
		cp := *e
		cp.Findings = cloneFindings(e.Findings)
		return &cp, nil
	}
	key := CacheKeyForPath(file)
	path := filepath.Join(s.filesDir, key+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

func (s *Store) removeLocked(file string) {
	file = NormalizePath(file)
	if _, ok := s.manifest.Files[file]; !ok {
		return
	}
	delete(s.manifest.Files, file)
	if s.ephemeral {
		if s.memoryEntries != nil {
			delete(s.memoryEntries, file)
		}
	} else {
		path := filepath.Join(s.filesDir, CacheKeyForPath(file)+".json")
		_ = os.Remove(path)
	}
	s.dirty = true
}

func (s *Store) writeManifestLocked() error {
	s.manifest.SchemaVersion = CACHE_VERSION
	s.manifest.ToolVersion = s.toolVersion
	data, err := json.MarshalIndent(s.manifest, "", "  ")
	if err != nil {
		return &Error{Op: "marshal", Path: manifestName, Err: err}
	}
	path := filepath.Join(s.cacheDir, manifestName)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return &Error{Op: "write", Path: tmp, Err: err}
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return &Error{Op: "rename", Path: path, Err: err}
	}
	s.dirty = false
	return nil
}

// evictLocked drops oldest entries until total files/ size is under target.
func (s *Store) evictLocked() error {
	type item struct {
		file string
		at   string
		size int64
	}
	var items []item
	var total int64
	for file, meta := range s.manifest.Files {
		path := filepath.Join(s.filesDir, CacheKeyForPath(file)+".json")
		fi, err := os.Stat(path)
		sz := int64(0)
		if err == nil {
			sz = fi.Size()
		}
		total += sz
		items = append(items, item{file: file, at: meta.CachedAt, size: sz})
	}
	if total <= s.maxSizeBytes {
		return nil
	}
	target := int64(float64(s.maxSizeBytes) * s.evictTarget)
	sort.Slice(items, func(i, j int) bool {
		return items[i].at < items[j].at // oldest first
	})
	for _, it := range items {
		if total <= target {
			break
		}
		s.removeLocked(it.file)
		total -= it.size
	}
	return nil
}

func clearFilesDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		_ = os.RemoveAll(filepath.Join(dir, e.Name()))
	}
	return nil
}

func cloneFindings(in []rules.Finding) []rules.Finding {
	if in == nil {
		return nil
	}
	out := make([]rules.Finding, len(in))
	copy(out, in)
	return out
}

// validatePurgePath refuses to delete non-cache or sensitive paths.
func validatePurgePath(dir string) error {
	fi, err := os.Lstat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return &Error{Op: "stat", Path: dir, Err: err}
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		return &Error{Op: "purge", Path: dir, Err: errSymlink}
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return &Error{Op: "abs", Path: dir, Err: err}
	}
	// Must look like a CodeHound cache (conventional name or manifest+files).
	base := filepath.Base(abs)
	manifest := filepath.Join(abs, manifestName)
	files := filepath.Join(abs, filesSubdir)
	conventional := base == DEFAULT_CACHE_DIR
	hasLayout := fileExists(manifest) && dirExists(files)
	if !conventional && !hasLayout {
		return &Error{Op: "purge", Path: dir, Err: errNotCache}
	}
	return nil
}

func fileExists(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && !fi.IsDir()
}

func dirExists(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && fi.IsDir()
}

type simpleError string

func (e simpleError) Error() string { return string(e) }

const (
	errSymlink  simpleError = "refusing to purge a symlinked cache path"
	errNotCache simpleError = "path is not a CodeHound cache directory"
)
