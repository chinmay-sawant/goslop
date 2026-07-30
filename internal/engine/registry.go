package engine

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/chinmay-sawant/goslop/internal/core"
	golang "github.com/chinmay-sawant/goslop/internal/lang/go"
)

// RegistryError is a plugin composition error.
type RegistryError struct {
	Msg string
}

func (e *RegistryError) Error() string { return e.Msg }

// Registry holds language plugins and detectors indexed by language.
type Registry struct {
	plugins     []core.LanguagePlugin
	byExtension map[string]int
	byID        map[core.LanguageID]int
	detectors   []core.Detector
	byLanguage  map[core.LanguageID][]int
	allIndices  []int
}

// scanSession owns the detector instances and lifecycle state for exactly one
// Analyzer.AnalyzePaths call. Registry keeps the immutable catalogue used by
// rule listing and file discovery; a session gets fresh detector instances.
type scanSession struct {
	detectors  []core.Detector
	byLanguage map[core.LanguageID][]int
}

// DetectorCount returns the number of detectors available to this scan.
func (s *scanSession) DetectorCount() int {
	if s == nil {
		return 0
	}
	return len(s.detectors)
}

// Detectors returns session-local detectors in registration order.
func (s *scanSession) Detectors() []core.Detector {
	if s == nil {
		return nil
	}
	return s.detectors
}

// DetectorIndices returns session-local detector indices for a language.
func (s *scanSession) DetectorIndices(language core.LanguageID) []int {
	if s == nil {
		return nil
	}
	return s.byLanguage[language]
}

// Detector returns a session-local detector by index.
func (s *scanSession) Detector(index int) core.Detector {
	if s == nil || index < 0 || index >= len(s.detectors) {
		return nil
	}
	return s.detectors[index]
}

// NewRegistry builds a registry from plugins. Each plugin's Detectors() is
// invoked once; detectors are validated against the plugin language.
func NewRegistry(plugins []core.LanguagePlugin) (*Registry, error) {
	r := &Registry{
		byExtension: make(map[string]int),
		byID:        make(map[core.LanguageID]int),
		byLanguage:  make(map[core.LanguageID][]int),
	}
	for _, plugin := range plugins {
		if err := r.RegisterPlugin(plugin); err != nil {
			return nil, err
		}
	}
	return r, nil
}

// RegisterPlugin appends a plugin (and its detectors) to an existing registry.
func (r *Registry) RegisterPlugin(plugin core.LanguagePlugin) error {
	if r == nil {
		return &RegistryError{Msg: "nil registry"}
	}
	if plugin == nil {
		return &RegistryError{Msg: "nil language plugin"}
	}
	if r.byExtension == nil {
		r.byExtension = make(map[string]int)
	}
	if r.byID == nil {
		r.byID = make(map[core.LanguageID]int)
	}
	if r.byLanguage == nil {
		r.byLanguage = make(map[core.LanguageID][]int)
	}

	id := plugin.ID()
	if _, exists := r.byID[id]; exists {
		return &RegistryError{Msg: fmt.Sprintf("duplicate language plugin: %s", id)}
	}

	idx := len(r.plugins)
	r.byID[id] = idx
	for _, ext := range plugin.Extensions() {
		ext = strings.TrimPrefix(ext, ".")
		if _, exists := r.byExtension[ext]; exists {
			return &RegistryError{Msg: fmt.Sprintf("extension %q already registered", ext)}
		}
		r.byExtension[ext] = idx
	}
	for _, det := range plugin.Detectors() {
		if det == nil {
			return &RegistryError{Msg: "nil detector"}
		}
		if det.Language() != id {
			return &RegistryError{
				Msg: fmt.Sprintf("detector language %s does not match plugin language %s", det.Language(), id),
			}
		}
		dIdx := len(r.detectors)
		r.detectors = append(r.detectors, det)
		r.byLanguage[id] = append(r.byLanguage[id], dIdx)
		r.allIndices = append(r.allIndices, dIdx)
	}
	r.plugins = append(r.plugins, plugin)
	return nil
}

// newScanSession creates fresh detector instances for one scan. NewDetectors
// is an explicit factory seam so a plugin cannot accidentally reuse lifecycle
// state that the registry holds for catalogue listing.
func (r *Registry) newScanSession() (*scanSession, error) {
	if r == nil {
		return nil, &RegistryError{Msg: "nil registry"}
	}
	session := &scanSession{byLanguage: make(map[core.LanguageID][]int, len(r.plugins))}
	for _, plugin := range r.plugins {
		id := plugin.ID()
		for _, det := range plugin.NewDetectors() {
			if det == nil {
				return nil, &RegistryError{Msg: "nil scan detector"}
			}
			if det.Language() != id {
				return nil, &RegistryError{
					Msg: fmt.Sprintf("scan detector language %s does not match plugin language %s", det.Language(), id),
				}
			}
			idx := len(session.detectors)
			session.detectors = append(session.detectors, det)
			session.byLanguage[id] = append(session.byLanguage[id], idx)
		}
	}
	return session, nil
}

// DetectorCount returns the number of registered detectors.
func (r *Registry) DetectorCount() int {
	if r == nil {
		return 0
	}
	return len(r.detectors)
}

// Detectors returns all detectors in registration order.
func (r *Registry) Detectors() []core.Detector {
	if r == nil {
		return nil
	}
	return r.detectors
}

// Plugins returns all language plugins in registration order.
func (r *Registry) Plugins() []core.LanguagePlugin {
	if r == nil {
		return nil
	}
	return r.plugins
}

// Plugin returns the language plugin for id.
func (r *Registry) Plugin(id core.LanguageID) (core.LanguagePlugin, bool) {
	if r == nil {
		return nil, false
	}
	idx, ok := r.byID[id]
	if !ok {
		return nil, false
	}
	return r.plugins[idx], true
}

// DetectorIndices returns detector indices for language.
func (r *Registry) DetectorIndices(language core.LanguageID) []int {
	if r == nil {
		return nil
	}
	return r.byLanguage[language]
}

// Detector returns the detector at index, or nil.
func (r *Registry) Detector(index int) core.Detector {
	if r == nil || index < 0 || index >= len(r.detectors) {
		return nil
	}
	return r.detectors[index]
}

// PluginForPath resolves a language plugin by file extension.
func (r *Registry) PluginForPath(path string) core.LanguagePlugin {
	if r == nil {
		return nil
	}
	ext := strings.TrimPrefix(filepath.Ext(path), ".")
	idx, ok := r.byExtension[ext]
	if !ok {
		return nil
	}
	return r.plugins[idx]
}

// PluginForID resolves a language plugin by language id.
func (r *Registry) PluginForID(id core.LanguageID) core.LanguagePlugin {
	p, _ := r.Plugin(id)
	return p
}

// ExtensionMap returns extension → language for walk collection.
func (r *Registry) ExtensionMap() map[string]core.LanguageID {
	out := make(map[string]core.LanguageID, len(r.byExtension))
	for ext, idx := range r.byExtension {
		out[ext] = r.plugins[idx].ID()
	}
	return out
}

// AllRuleIDs returns every rule id from every detector (deduplicated, registration order).
func (r *Registry) AllRuleIDs() []string {
	if r == nil {
		return nil
	}
	seen := make(map[string]struct{})
	var out []string
	for _, det := range r.detectors {
		for _, id := range det.RuleIDs() {
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			out = append(out, id)
		}
	}
	return out
}

// RuleIDs is an alias for AllRuleIDs.
func (r *Registry) RuleIDs() []string { return r.AllRuleIDs() }

// DetectorsForLanguage returns detector instances registered for language.
func (r *Registry) DetectorsForLanguage(language core.LanguageID) []core.Detector {
	idxs := r.DetectorIndices(language)
	out := make([]core.Detector, 0, len(idxs))
	for _, i := range idxs {
		if d := r.Detector(i); d != nil {
			out = append(out, d)
		}
	}
	return out
}

// RuleInfo is a listed rule for --list-rules.
type RuleInfo struct {
	ID    string
	Title string
}

// --- default process-wide registry ---

var (
	defaultRegMu sync.RWMutex
	defaultReg   *Registry
)

// SetDefaultRegistry installs the process-wide registry.
func SetDefaultRegistry(r *Registry) {
	defaultRegMu.Lock()
	defer defaultRegMu.Unlock()
	defaultReg = r
}

// DefaultRegistry returns the process-wide registry with the Go language plugin
// and seed detectors registered. Language plugins must not import this package
// (avoids an import cycle); the engine owns default composition.
func DefaultRegistry() *Registry {
	defaultRegMu.RLock()
	r := defaultReg
	defaultRegMu.RUnlock()
	if r != nil {
		return r
	}
	defaultRegMu.Lock()
	defer defaultRegMu.Unlock()
	if defaultReg == nil {
		reg, err := NewRegistry([]core.LanguagePlugin{golang.NewPlugin()})
		if err != nil {
			reg, _ = NewRegistry(nil)
		}
		defaultReg = reg
	}
	return defaultReg
}

// RegisterDefaultPlugin registers a plugin onto the process-wide registry.
// Duplicate language plugins are ignored (idempotent for init re-entry in tests).
func RegisterDefaultPlugin(plugin core.LanguagePlugin) error {
	if plugin == nil {
		return &RegistryError{Msg: "nil language plugin"}
	}
	reg := DefaultRegistry()
	if _, ok := reg.Plugin(plugin.ID()); ok {
		return nil
	}
	return reg.RegisterPlugin(plugin)
}

// ListRules returns rule IDs from the default registry.
func ListRules() []RuleInfo {
	r := DefaultRegistry()
	ids := r.AllRuleIDs()
	out := make([]RuleInfo, len(ids))
	for i, id := range ids {
		title := ""
		for _, det := range r.Detectors() {
			if meta := det.MetadataFor(id); meta != nil {
				title = meta.Title
				break
			}
		}
		out[i] = RuleInfo{ID: id, Title: title}
	}
	return out
}

// AnalyzePaths is a package-level entry that uses the default registry.
func AnalyzePaths(paths []string, ctx *core.ScanContext) (*AnalysisResult, error) {
	if ctx == nil {
		ctx = core.DefaultScanContext()
	}
	return NewAnalyzer(ctx, DefaultRegistry()).AnalyzePaths(paths)
}
