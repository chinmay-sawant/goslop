// Package python is a minimal LanguagePlugin stub for multi-language WIP.
//
// It is intentionally not registered in engine.DefaultRegistry. Production
// remains Go-only (core.DefaultEnabledLanguages). Tests may compose a registry
// with python.NewPlugin() until Phase 3 wires languages config filtering.
//
// Parse path is source-only (no AST / no CGO). Zero detectors.
package python

import "github.com/chinmay-sawant/goslop/internal/core"

// Plugin is a source-only Python language stub (no detectors).
type Plugin struct {
	core.BasePlugin
}

// NewPlugin returns a LanguagePlugin for Python source files (.py).
func NewPlugin() core.LanguagePlugin {
	return &Plugin{}
}

// Register adds the Python stub to an engine registry (tests / experimental).
// Prefer NewRegistry([]core.LanguagePlugin{..., python.NewPlugin()}) so
// DefaultRegistry stays Go-only.
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

// Detectors implements core.LanguagePlugin (empty catalogue).
func (p *Plugin) Detectors() []core.Detector { return nil }

// NewDetectors implements core.LanguagePlugin (empty session set).
func (p *Plugin) NewDetectors() []core.Detector { return nil }

// ParseSource returns a source-only ParsedUnit with LanguagePython.
// Overrides BasePlugin so unit.Language is not left as Go.
func (p *Plugin) ParseSource(path, source string) (*core.ParseResult, error) {
	return &core.ParseResult{
		Unit:    core.NewParsedUnit(core.LanguagePython, path, source),
		Quality: core.ParseQualitySourceOnly,
	}, nil
}
