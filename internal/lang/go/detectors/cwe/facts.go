package cwe

import (
	"github.com/chinmay-sawant/goslop/internal/ast"
	"github.com/chinmay-sawant/goslop/internal/core"
)

// GoCweFacts is the per-file fact bag for structural CWE detectors.
// Phase 7 uses SourceIndex needle presence; call/assignment facts can be
// layered later without changing the rule registration API.
type GoCweFacts struct {
	Index ast.SourceIndex
}

// BuildFacts constructs GoCweFacts for a unit.
//
// Structural CWE rules only need a SourceIndex over cweNeedles (no AST walk).
// PERF/BP needle tables differ, so indexes stay pack-local; parse reuse for
// AST-based packs is handled by goparse.TreeForUnit / unit.Tree, not here.
func BuildFacts(unit *core.ParsedUnit) *GoCweFacts {
	facts := &GoCweFacts{}
	if unit == nil || unit.Source == "" {
		facts.Index = ast.Build("", cweNeedles)
		return facts
	}
	facts.Index = ast.Build(unit.Source, cweNeedles)
	return facts
}
