// Package ast provides shared tree-sitter walk helpers, line/column mapping,
// and multi-needle SourceIndex utilities for detector hot paths.
package ast

import "sort"

// LineCol is a 1-indexed source location (matches Rust rules::LineCol).
// Standalone until internal/core.ParsedUnit / internal/rules land.
type LineCol struct {
	Line   int
	Column int
}

// Valid reports whether line and column are at least 1.
func (lc LineCol) Valid() bool {
	return lc.Line >= 1 && lc.Column >= 1
}

// ComputeLineStarts builds a per-line start-offset table from source text.
// The returned slice contains, in order, the byte offset of the first byte of
// each line (always starting with 0). Used so LineColWithStarts is O(log N).
func ComputeLineStarts(source string) []int {
	starts := make([]int, 0, 64)
	starts = append(starts, 0)
	for i := 0; i < len(source); i++ {
		if source[i] == '\n' {
			starts = append(starts, i+1)
		}
	}
	return starts
}

// ComputeLineStartsBytes is ComputeLineStarts for a byte slice (no copy).
func ComputeLineStartsBytes(source []byte) []int {
	starts := make([]int, 0, 64)
	starts = append(starts, 0)
	for i := 0; i < len(source); i++ {
		if source[i] == '\n' {
			starts = append(starts, i+1)
		}
	}
	return starts
}

// LineColWithStarts returns 1-indexed (line, column) for a byte offset using a
// precomputed per-line start-offset table. O(log N).
func LineColWithStarts(lineStarts []int, byteOffset int) (line, col int) {
	if len(lineStarts) == 0 {
		return 1, 1
	}
	if byteOffset < 0 {
		return 1, 1
	}
	// Find largest index with lineStarts[i] <= byteOffset.
	i := sort.Search(len(lineStarts), func(i int) bool {
		return lineStarts[i] > byteOffset
	})
	if i == 0 {
		return 1, 1
	}
	idx := i - 1
	return idx + 1, byteOffset - lineStarts[idx] + 1
}

// LineColAt is a convenience that returns a LineCol for the offset.
func LineColAt(lineStarts []int, byteOffset int) LineCol {
	line, col := LineColWithStarts(lineStarts, byteOffset)
	return LineCol{Line: line, Column: col}
}
