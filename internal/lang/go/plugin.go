// Package golang implements the Go LanguagePlugin for goslop.
//
// Import path ends in /go which is a reserved keyword as a package name, so the
// package is named "golang". Callers use:
//
//	golang.NewPlugin() / golang.Register(reg)
//
// Parse path is pure Go (go/parser + go/ast via goparse) — no CGO.
//
// This package is the reference LanguagePlugin. To add another language without
// CGO, implement core.LanguagePlugin the same way: pure-Go ParseSource, opaque
// unit.Tree, detectors that type-assert only inside the language package.
// See the package doc on core.LanguagePlugin.
package golang

import (
	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/lang/go/detectors"
	"github.com/chinmay/codehound/internal/lang/go/goparse"
)

// Plugin is the Go language plugin.
type Plugin struct {
	core.BasePlugin
}

// NewPlugin returns a LanguagePlugin for Go source files.
func NewPlugin() core.LanguagePlugin {
	return &Plugin{}
}

// Register adds the Go plugin to an engine registry.
func Register(reg any) {
	if reg == nil {
		return
	}
	p := NewPlugin()
	switch r := reg.(type) {
	case interface {
		RegisterPlugin(core.LanguagePlugin) error
	}:
		_ = r.RegisterPlugin(p)
	case interface{ RegisterPlugin(core.LanguagePlugin) }:
		r.RegisterPlugin(p)
	}
}

// ID implements core.LanguagePlugin.
func (p *Plugin) ID() core.LanguageID { return core.LangGo }

// Extensions implements core.LanguagePlugin.
func (p *Plugin) Extensions() []string { return []string{"go"} }

// Detectors implements core.LanguagePlugin.
func (p *Plugin) Detectors() []core.Detector { return detectors.All() }

// ParseSource parses Go source with go/parser and attaches *goparse.Tree.
// On parse failure, falls back to a source-only unit so text-level detectors still run.
func (p *Plugin) ParseSource(path, source string) (*core.ParsedUnit, error) {
	tree, parseErr := goparse.Parse([]byte(source))
	// Intentionally ignore parseErr for the fallback path: text-level detectors
	// still need a unit when the AST is missing. Partial ASTs are still attached.
	if tree == nil || tree.File == nil {
		return core.NewParsedUnit(core.LangGo, path, source), nil
	}
	_ = parseErr
	return core.NewParsedUnitWithTree(core.LangGo, path, source, tree), nil
}

// FunctionNodeKinds documents function-like constructs (for tooling; go/ast uses FuncDecl).
func (p *Plugin) FunctionNodeKinds() []string {
	return []string{"FuncDecl", "FuncLit"}
}

// LoopNodeKinds documents loop constructs (for tooling; go/ast uses ForStmt/RangeStmt).
func (p *Plugin) LoopNodeKinds() []string {
	return []string{"ForStmt", "RangeStmt"}
}
