// Package python implements a minimal Python LanguagePlugin stub for goslop.
//
// Status: WIP foundation for multi-language enable/disable (epic #39). This
// package provides ID / Extensions / source-only ParseSource so the engine can
// resolve .py files when the plugin is registered. There are no detectors yet
// (Detectors / NewDetectors return empty). Full parse trees, detectors, config
// TOML wiring, and ruleset JSON catalogues are out of scope for this stub.
//
// Callers use:
//
//	python.NewPlugin() / python.Register(reg)
//
// DefaultRegistry stays Go-only; include Python via engine.NewRegistryWithLanguages
// or Register when intentionally enabling the language (tests / future config).
package python

import (
	"github.com/chinmay-sawant/goslop/internal/core"
)

// Plugin is the Python language plugin stub (source-only, zero detectors).
type Plugin struct {
	core.BasePlugin
}

// NewPlugin returns a LanguagePlugin for Python source files.
func NewPlugin() core.LanguagePlugin {
	return &Plugin{}
}

// Register adds the Python plugin to an engine registry.
// Accepts any value that implements RegisterPlugin(core.LanguagePlugin) with or
// without an error return (same shape as golang.Register).
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

// Detectors implements core.LanguagePlugin. Empty until Python detectors land.
func (p *Plugin) Detectors() []core.Detector { return nil }

// NewDetectors implements core.LanguagePlugin. Empty catalogue matches Detectors.
func (p *Plugin) NewDetectors() []core.Detector { return nil }

// ParseSource returns a source-only ParsedUnit tagged LanguagePython.
// Overrides BasePlugin.ParseSource which defaults language to Go.
func (p *Plugin) ParseSource(path, source string) (*core.ParseResult, error) {
	return &core.ParseResult{
		Unit:    core.NewParsedUnit(core.LanguagePython, path, source),
		Quality: core.ParseQualitySourceOnly,
	}, nil
}
