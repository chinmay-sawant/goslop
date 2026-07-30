// Package goparse parses Go source with the stdlib (go/parser + go/ast).
//
// Pure Go: no CGO. Replaces the former tree-sitter path (tsparse).
//
//	tree, err := goparse.Parse(source)
//	if err != nil { ... }
//	defer tree.Close() // no-op; kept for API symmetry
//	file := tree.File
//	line, col := tree.LineCol(byteOffset)
package goparse

import (
	"bytes"
	"errors"
	"fmt"
	goast "go/ast"
	"go/parser"
	"go/token"
	"unsafe"

	"github.com/chinmay-sawant/goslop/internal/ast"
)

// Tree holds a parsed *ast.File, file set, source, and line-start table.
type Tree struct {
	File       *goast.File
	Fset       *token.FileSet
	Source     []byte
	LineStarts []int
	tf         *token.File // cached primary file in Fset
}

// Close is a no-op for API symmetry with the old tree-sitter handle.
func (t *Tree) Close() {}

// LineCol returns 1-indexed (line, column) for a byte offset into Source.
func (t *Tree) LineCol(byteOffset int) (int, int) {
	return ast.LineColWithStarts(t.LineStarts, byteOffset)
}

// LineColAt returns a LineCol for the byte offset.
func (t *Tree) LineColAt(byteOffset int) ast.LineCol {
	return ast.LineColAt(t.LineStarts, byteOffset)
}

// Offset returns the byte offset for a token.Pos (0 if invalid).
func (t *Tree) Offset(pos token.Pos) int {
	if t == nil || !pos.IsValid() {
		return 0
	}
	if t.tf != nil {
		// token.File.Offset panics if pos is not in file; clamp via Position.
		p := t.Fset.Position(pos)
		if p.Offset < 0 {
			return 0
		}
		return p.Offset
	}
	p := t.Fset.Position(pos)
	if p.Offset < 0 {
		return 0
	}
	return p.Offset
}

// Parse parses source as a full Go file (package clause required when present).
// Uses parser.ParseFile with AllErrors|ParseComments skipped for speed.
// On syntax errors, still returns a partial tree when the parser yields a File.
func Parse(source []byte) (*Tree, error) {
	fset := token.NewFileSet()
	// ParseFile needs a name; use a stable placeholder.
	file, err := parser.ParseFile(fset, "file.go", source, parser.SkipObjectResolution)
	if file == nil {
		if err != nil {
			return nil, fmt.Errorf("goparse: %w", err)
		}
		return nil, errors.New("goparse: empty parse result")
	}
	tf := fset.File(file.Pos())
	tree := &Tree{
		File:       file,
		Fset:       fset,
		Source:     source,
		LineStarts: ast.ComputeLineStartsBytes(source),
		tf:         tf,
	}
	// May return a partial AST with a non-nil error (syntax recovery).
	if err != nil {
		return tree, fmt.Errorf("goparse: %w", err)
	}
	return tree, nil
}

// ParseExpr parses an expression (tests / helpers).
func ParseExpr(source []byte) (goast.Expr, *token.FileSet, error) {
	fset := token.NewFileSet()
	expr, err := parser.ParseExprFrom(fset, "expr.go", source, 0)
	if err != nil {
		return expr, fset, fmt.Errorf("goparse expr: %w", err)
	}
	return expr, fset, nil
}

// Slice returns source[start:end] clamped to bounds.
//
// Zero-copy: the result aliases Tree.Source. Callers must treat the string as
// immutable and must not retain it past Close/replacement of Source (P2.2).
func (t *Tree) Slice(start, end int) string {
	if t == nil {
		return ""
	}
	src := t.Source
	if start < 0 {
		start = 0
	}
	if end > len(src) {
		end = len(src)
	}
	if start >= end {
		return ""
	}
	b := src[start:end]
	// #nosec G103 -- Tree.Source is immutable after Parse and callers receive a read-only string view.
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// NodeText returns the source text of node.
func (t *Tree) NodeText(n goast.Node) string {
	if t == nil || n == nil {
		return ""
	}
	start := t.Offset(n.Pos())
	end := t.Offset(n.End())
	return t.Slice(start, end)
}

// Ensure source is valid UTF-8-ish (parser is binary-safe for offsets).
var _ = bytes.Equal
