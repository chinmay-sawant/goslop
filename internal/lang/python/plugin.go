// Package python implements the Python LanguagePlugin for goslop.
//
// Status: source-only ParseSource (no Python AST tree) plus BP-PY bad-practice
// heuristics registered via detectors.All(). CWE / PERF families land in sibling
// workstreams. Production DefaultRegistry stays Go-only; include Python via
// engine.NewRegistryWithLanguages or python.Register when intentionally enabling
// the language (tests / config languages).
//
// Callers use:
//
//	python.NewPlugin() / python.Register(reg)
package python

import (
	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/lang/python/detectors"
)

// Plugin is the Python language plugin (source-only parse + BP detectors).
type Plugin struct {
	core.BasePlugin
}

// NewPlugin returns a LanguagePlugin for Python source files (.py).
func NewPlugin() core.LanguagePlugin {
	return &Plugin{}
}

// Register adds the Python plugin to an engine registry.
// Accepts any value that implements RegisterPlugin(core.LanguagePlugin) with or
// without an error return (same shape as golang.Register).
// Prefer NewRegistry / NewRegistryWithLanguages so DefaultRegistry stays Go-only.
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
func (p *Plugin) ID() core.LanguageID { return core.LanguagePython }

// Extensions implements core.LanguagePlugin.
func (p *Plugin) Extensions() []string { return []string{"py"} }

// Detectors implements core.LanguagePlugin.
func (p *Plugin) Detectors() []core.Detector { return detectors.All() }

// NewDetectors creates the session-local Python detector set.
func (p *Plugin) NewDetectors() []core.Detector { return detectors.All() }

// ParseSource returns a source-only ParsedUnit tagged LanguagePython.
// Overrides BasePlugin.ParseSource which defaults language to Go.
func (p *Plugin) ParseSource(path, source string) (*core.ParseResult, error) {
	return &core.ParseResult{
		Unit:    core.NewParsedUnit(core.LanguagePython, path, source),
		Quality: core.ParseQualitySourceOnly,
	}, nil
}
