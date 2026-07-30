package cwe

import (
	"github.com/chinmay-sawant/goslop/internal/ast"
	"github.com/chinmay-sawant/goslop/internal/core"
)

// PyCweFacts is the per-file fact bag for Python CWE source-pattern detectors.
// Priority batch uses SourceIndex needle presence; call heuristics scan unit.Source.
type PyCweFacts struct {
	Index ast.SourceIndex
}

// BuildFacts constructs PyCweFacts for a unit.
//
// Rules only need a SourceIndex over pyCweNeedles (no Python AST / CGO).
// Pack-local needle table — do not reuse Go cweNeedles.
func BuildFacts(unit *core.ParsedUnit) *PyCweFacts {
	facts := &PyCweFacts{}
	if unit == nil || unit.Source == "" {
		facts.Index = ast.Build("", pyCweNeedles)
		return facts
	}
	facts.Index = ast.Build(unit.Source, pyCweNeedles)
	return facts
}
