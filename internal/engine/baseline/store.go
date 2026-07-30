// Package baseline stores and filters baselined finding fingerprints.
package baseline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/chinmay-sawant/goslop/internal/rules"
)

// FileName is the default baseline filename discovered by walking up from cwd.
const FileName = ".goslop-baseline.json"

// Version is the wire format version.
const Version = "1"

// Entry is one baselined finding location/fingerprint.
type Entry struct {
	File        string  `json:"file"`
	Line        int     `json:"line"`
	Column      int     `json:"column"`
	Fingerprint string  `json:"fingerprint"`
	Reason      *string `json:"reason,omitempty"`
	Expires     *string `json:"expires,omitempty"`
}

// Baseline is the on-disk document of suppressed finding fingerprints.
type Baseline struct {
	Version     string             `json:"version"`
	GeneratedAt string             `json:"generated_at"`
	ToolVersion string             `json:"tool_version"`
	Entries     map[string][]Entry `json:"entries"` // rule_id → entries
}

// FromFindings builds a baseline from live findings.
func FromFindings(findings []rules.Finding, toolVersion string) *Baseline {
	b := &Baseline{
		Version:     Version,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		ToolVersion: toolVersion,
		Entries:     make(map[string][]Entry),
	}
	for i := range findings {
		f := &findings[i]
		f.EnsureFingerprint()
		b.Entries[f.RuleID] = append(b.Entries[f.RuleID], Entry{
			File:        f.File,
			Line:        f.Line,
			Column:      f.Column,
			Fingerprint: f.FingerprintString(),
		})
	}
	return b
}

// Load reads a baseline JSON file.
func Load(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, err
	}
	if b.Entries == nil {
		b.Entries = make(map[string][]Entry)
	}
	return &b, nil
}

// Save writes the baseline atomically.
func (b *Baseline) Save(path string) error {
	if b == nil {
		return nil
	}
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// Discover walks up from start looking for .goslop-baseline.json (stops at .git).
func Discover(start string) string {
	dir := start
	if dir == "" {
		dir = "."
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		abs = dir
	}
	// If start is a file, begin at parent.
	if fi, err := os.Stat(abs); err == nil && !fi.IsDir() {
		abs = filepath.Dir(abs)
	}
	for {
		candidate := filepath.Join(abs, FileName)
		if fi, err := os.Stat(candidate); err == nil && !fi.IsDir() {
			return candidate
		}
		if _, err := os.Stat(filepath.Join(abs, ".git")); err == nil {
			return ""
		}
		parent := filepath.Dir(abs)
		if parent == abs {
			return ""
		}
		abs = parent
	}
}

// Contains reports whether finding is baselined (fingerprint first, then location).
func (b *Baseline) Contains(f *rules.Finding) bool {
	if b == nil || f == nil {
		return false
	}
	entries, ok := b.Entries[f.RuleID]
	if !ok {
		return false
	}
	now := time.Now().UTC().Format(time.RFC3339)
	fp := f.FingerprintString()
	locationMatch := false
	for i := range entries {
		e := &entries[i]
		if entryExpired(e, now) {
			continue
		}
		if e.Fingerprint == fp {
			return true
		}
		if e.File == f.File && e.Line == f.Line && e.Column == f.Column {
			locationMatch = true
		}
	}
	return locationMatch
}

// Filter removes baselined findings. When showBaselined is true, baselined
// findings are kept with Suppressed=true and severity Info.
// Returns (filtered, baselinedCount). Does not mutate the input slice header.
func (b *Baseline) Filter(findings []rules.Finding, showBaselined bool) ([]rules.Finding, int) {
	if b == nil || len(findings) == 0 {
		return findings, 0
	}
	out := make([]rules.Finding, 0, len(findings))
	n := 0
	for i := range findings {
		f := findings[i]
		if !b.Contains(&f) {
			out = append(out, f)
			continue
		}
		n++
		if showBaselined {
			f.Severity = rules.SeverityInfo
			f.Suppressed = true
			out = append(out, f)
		}
	}
	return out, n
}

// EntryCount returns total baselined entries.
func (b *Baseline) EntryCount() int {
	if b == nil {
		return 0
	}
	n := 0
	for _, ents := range b.Entries {
		n += len(ents)
	}
	return n
}

func entryExpired(e *Entry, nowISO string) bool {
	if e.Expires == nil || *e.Expires == "" {
		return false
	}
	return *e.Expires < nowISO
}
