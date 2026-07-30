package goparse

import "github.com/chinmay-sawant/goslop/internal/core"

// AsTree type-asserts v to a *Tree with a non-nil File suitable for inspection.
func AsTree(v any) *Tree {
	t, ok := v.(*Tree)
	if !ok || t == nil || t.File == nil {
		return nil
	}
	return t
}

// TreeForUnit returns the Go AST for unit, parsing at most once per unit.
//
// When unit.Tree already holds a usable *Tree, it is returned as-is.
// Otherwise source is parsed and the result is stored on unit.Tree so PERF,
// BP, CWE-taint, and any other pack share one parse per file (cross-pack reuse).
// Lifetime of an attached tree is owned by the unit / analyzer closeUnitTree.
func TreeForUnit(unit *core.ParsedUnit) *Tree {
	if unit == nil {
		return nil
	}
	if t := AsTree(unit.Tree); t != nil {
		return t
	}
	if unit.Source == "" {
		return nil
	}
	t, _ := Parse([]byte(unit.Source))
	if t == nil || t.File == nil {
		return nil
	}
	unit.Tree = t
	return t
}
