package ast

// This package is language-agnostic: source index + line/col helpers.
// Language-specific tree walks live in each language plugin
// (e.g. go/ast for Go, future pure-Go parsers for other languages).
//
// Former tree-sitter walk helpers were removed with the CGO parse path.
