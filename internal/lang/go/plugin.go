// Package golang implements the Go LanguagePlugin for CodeHound.
//
// Import path ends in /go which is a reserved keyword as a package name, so the
// package is named "golang". Callers use:
//
//	golang.NewPlugin() / golang.Register(reg)
//
// The engine package's DefaultRegistry registers this plugin (do not import
// engine from this package — that creates an import cycle).
package golang

import (
	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/lang/go/detectors"
	"github.com/chinmay/codehound/internal/lang/go/tsparse"
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
// Accepts either RegisterPlugin(plugin) or RegisterPlugin(plugin) error.
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

// ParseSource parses Go source with tree-sitter and attaches the CST to the unit.
// On parse failure, falls back to a source-only unit so text-level detectors still run.
func (p *Plugin) ParseSource(path, source string) (*core.ParsedUnit, error) {
	tree, err := tsparse.Parse([]byte(source))
	if err != nil {
		return core.NewParsedUnit(core.LangGo, path, source), nil
	}
	return core.NewParsedUnitWithTree(core.LangGo, path, source, tree), nil
}

// FunctionNodeKinds returns tree-sitter function-like node types.
func (p *Plugin) FunctionNodeKinds() []string {
	return []string{"function_declaration", "method_declaration"}
}

// LoopNodeKinds returns tree-sitter loop node types.
func (p *Plugin) LoopNodeKinds() []string {
	return []string{"for_statement"}
}
