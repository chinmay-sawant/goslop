package core

// ParsedUnit is a single source file prepared for analysis.
type ParsedUnit struct {
	Language    LanguageID
	Path        string
	DisplayPath string
	Source      string
	LineStarts  []int
	// Tree is an opaque language-specific AST handle (e.g. *goparse.Tree for Go);
	// nil when source-only.
	Tree any
}

// NewParsedUnit builds a source unit. Optional tree is an opaque AST handle.
// DisplayPath defaults to path.
func NewParsedUnit(language LanguageID, path, source string, tree ...any) *ParsedUnit {
	var t any
	if len(tree) > 0 {
		t = tree[0]
	}
	return &ParsedUnit{
		Language:    language,
		Path:        path,
		DisplayPath: path,
		Source:      source,
		LineStarts:  ComputeLineStarts(source),
		Tree:        t,
	}
}

// NewParsedUnitWithTree builds a unit that also carries a syntax tree handle.
func NewParsedUnitWithTree(language LanguageID, path, source string, tree any) *ParsedUnit {
	return NewParsedUnit(language, path, source, tree)
}

// ComputeLineStarts builds per-line start-offset table from source text.
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

// LineCol converts a byte offset into a one-indexed line and column.
func (u *ParsedUnit) LineCol(byteOffset int) (line, col int) {
	if u == nil {
		return 1, 1
	}
	starts := u.LineStarts
	if len(starts) == 0 {
		starts = ComputeLineStarts(u.Source)
	}
	return lineColWithStarts(starts, byteOffset)
}

func lineColWithStarts(lineStarts []int, byteOffset int) (line, col int) {
	if len(lineStarts) == 0 {
		return 1, 1
	}
	if byteOffset < 0 {
		byteOffset = 0
	}
	lo, hi := 0, len(lineStarts)
	for lo < hi {
		mid := (lo + hi) / 2
		if lineStarts[mid] <= byteOffset {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	idx := lo - 1
	if idx < 0 {
		return 1, 1
	}
	line = idx + 1
	col = byteOffset - lineStarts[idx] + 1
	if col < 1 {
		col = 1
	}
	return line, col
}
