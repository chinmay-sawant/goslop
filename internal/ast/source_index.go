package ast

import "strings"

// SourceIndex records which needles from a table appear as substrings in a
// source file. Lookup via Has is O(1) average after Build.
//
// Current implementation: one strings.Contains per needle at build time.
// Upgrade path: replace Build with Aho-Corasick (or similar multi-pattern
// matcher) when needle tables grow large enough that N×M scanning matters —
// see Rust codehound::lang::source_index which uses aho-corasick with a
// process-lifetime matcher cache keyed by needle table identity.
type SourceIndex struct {
	flags  []bool
	byName map[string]int
}

// Build constructs a SourceIndex for source against needles.
// Duplicate needles: first occurrence wins for Has lookup.
// Empty needles are never marked present.
func Build(source string, needles []string) SourceIndex {
	byName := make(map[string]int, len(needles))
	flags := make([]bool, len(needles))
	for i, n := range needles {
		if _, exists := byName[n]; !exists {
			byName[n] = i
		}
		if n != "" && strings.Contains(source, n) {
			flags[i] = true
		}
	}
	return SourceIndex{flags: flags, byName: byName}
}

// BuildBytes is Build for a byte slice (converted once).
func BuildBytes(source []byte, needles []string) SourceIndex {
	return Build(string(source), needles)
}

// Has reports whether needle was present in the indexed source.
// Unknown needles (not in the build table) return false.
func (s SourceIndex) Has(needle string) bool {
	idx, ok := s.byName[needle]
	if !ok {
		return false
	}
	if idx < 0 || idx >= len(s.flags) {
		return false
	}
	return s.flags[idx]
}

// HasAny is true if any of needles is present.
func (s SourceIndex) HasAny(needles []string) bool {
	for _, n := range needles {
		if s.Has(n) {
			return true
		}
	}
	return false
}

// Len returns the number of needles in the backing table.
func (s SourceIndex) Len() int { return len(s.flags) }

// Empty reports whether the needle table is empty.
func (s SourceIndex) Empty() bool { return len(s.flags) == 0 }
