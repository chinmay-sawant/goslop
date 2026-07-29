// Package tsparse parses Go source with tree-sitter.
//
// CGO is required: both github.com/tree-sitter/go-tree-sitter and
// github.com/tree-sitter/tree-sitter-go/bindings/go compile C sources.
// Build with CGO_ENABLED=1 and a working C toolchain (gcc/clang).
//
// Usage for other agents:
//
//	tree, err := tsparse.Parse(source)
//	if err != nil { ... }
//	defer tree.Close()
//	root := tree.RootNode()
//	line, col := tree.LineCol(int(node.StartByte()))
//	ast.WalkCalls(root, func(n *sitter.Node) { ... })
//
// For worker pools, reuse a Parser via NewParser / ParseWith instead of Parse.
package tsparse

import (
	"errors"
	"fmt"
	"sync"

	"github.com/chinmay/codehound/internal/ast"
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

// Re-export common tree-sitter types so callers can stay on this package.
type (
	Node       = sitter.Node
	TreeCursor = sitter.TreeCursor
)

// Tree wraps a tree-sitter syntax tree plus the source bytes and line-start table.
type Tree struct {
	inner      *sitter.Tree
	Source     []byte
	LineStarts []int
}

// RootNode returns the root of the syntax tree.
func (t *Tree) RootNode() *sitter.Node {
	if t == nil || t.inner == nil {
		return nil
	}
	return t.inner.RootNode()
}

// Close frees the underlying tree-sitter tree. Safe on nil receiver.
func (t *Tree) Close() {
	if t == nil || t.inner == nil {
		return
	}
	t.inner.Close()
	t.inner = nil
}

// LineCol returns 1-indexed (line, column) for a byte offset into Source.
func (t *Tree) LineCol(byteOffset int) (line, col int) {
	return ast.LineColWithStarts(t.LineStarts, byteOffset)
}

// LineColAt returns a LineCol for the byte offset.
func (t *Tree) LineColAt(byteOffset int) ast.LineCol {
	return ast.LineColAt(t.LineStarts, byteOffset)
}

// Inner exposes the raw tree-sitter tree (for advanced callers).
func (t *Tree) Inner() *sitter.Tree { return t.inner }

// Parser is a reusable Go-language tree-sitter parser (one per worker).
type Parser struct {
	p *sitter.Parser
}

var languageOnce sync.Once
var goLanguage *sitter.Language
var errLanguage error

func loadLanguage() (*sitter.Language, error) {
	languageOnce.Do(func() {
		goLanguage = sitter.NewLanguage(tree_sitter_go.Language())
		if goLanguage == nil {
			errLanguage = errors.New("tree-sitter-go: Language() returned nil")
		}
	})
	return goLanguage, errLanguage
}

// NewParser creates a parser configured for the Go grammar.
// Callers must Close the parser when done.
func NewParser() (*Parser, error) {
	lang, err := loadLanguage()
	if err != nil {
		return nil, err
	}
	p := sitter.NewParser()
	if err := p.SetLanguage(lang); err != nil {
		p.Close()
		return nil, fmt.Errorf("set go language: %w", err)
	}
	return &Parser{p: p}, nil
}

// Close frees the parser.
func (p *Parser) Close() {
	if p == nil || p.p == nil {
		return
	}
	p.p.Close()
	p.p = nil
}

// Parse parses source with this parser. The returned Tree owns the CST;
// call Tree.Close when finished. The Parser may be reused after Parse.
func (p *Parser) Parse(source []byte) (*Tree, error) {
	if p == nil || p.p == nil {
		return nil, errors.New("tsparse: nil parser")
	}
	tree := p.p.Parse(source, nil)
	if tree == nil {
		return nil, errors.New("tree-sitter returned nil tree")
	}
	return &Tree{
		inner:      tree,
		Source:     source,
		LineStarts: ast.ComputeLineStartsBytes(source),
	}, nil
}

// Parse is a convenience that creates a short-lived parser.
// Prefer NewParser + Parse for bulk work.
func Parse(source []byte) (*Tree, error) {
	p, err := NewParser()
	if err != nil {
		return nil, err
	}
	defer p.Close()
	return p.Parse(source)
}

// Walk visits every node under root (pre-order). Prefer ast.Walk for shared helpers.
func Walk(root *sitter.Node, visit func(n *sitter.Node) bool) {
	ast.Walk(root, visit)
}
