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
	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/lang/go/detectors"
	"github.com/chinmay-sawant/goslop/internal/lang/go/goparse"
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

// NewDetectors creates the session-local Go detector set. The registry keeps
// catalogue metadata separately, so lifecycle state never crosses scans.
func (p *Plugin) NewDetectors() []core.Detector { return detectors.All() }

// ParseSource parses Go source with go/parser and attaches *goparse.Tree.
// Syntax recovery keeps text-level analysis available and returns a diagnostic
// so callers can distinguish complete and partial parse quality.
func (p *Plugin) ParseSource(path, source string) (*core.ParseResult, error) {
	tree, parseErr := goparse.Parse([]byte(source))
	if tree == nil || tree.File == nil {
		result := &core.ParseResult{
			Unit:    core.NewParsedUnit(core.LangGo, path, source),
			Quality: core.ParseQualitySourceOnly,
		}
		if parseErr != nil {
			result.Diagnostic = parseErr.Error()
		}
		return result, nil
	}
	result := &core.ParseResult{
		Unit:    core.NewParsedUnitWithTree(core.LangGo, path, source, tree),
		Quality: core.ParseQualityComplete,
	}
	if parseErr != nil {
		result.Quality = core.ParseQualityPartial
		result.Diagnostic = parseErr.Error()
	}
	return result, nil
}

// FunctionNodeKinds documents function-like constructs (for tooling; go/ast uses FuncDecl).
func (p *Plugin) FunctionNodeKinds() []string {
	return []string{"FuncDecl", "FuncLit"}
}

// LoopNodeKinds documents loop constructs (for tooling; go/ast uses ForStmt/RangeStmt).
func (p *Plugin) LoopNodeKinds() []string {
	return []string{"ForStmt", "RangeStmt"}
}
