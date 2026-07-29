package core

// ProjectContext is language-neutral project context for dependency extraction.
type ProjectContext struct {
	Root string
}

// LanguagePlugin is a per-language backend: parse sources and supply detectors.
//
// # Adding a second language (no CGO)
//
// The engine is language-agnostic. To add e.g. Python without CGO:
//
//  1. Reserve / reuse a LanguageID in language.go (LanguagePython already exists).
//  2. Implement LanguagePlugin in internal/lang/<lang>/ (mirror lang/go):
//     - ID() / Extensions() — file selection via engine.Registry
//     - Detectors() — rule catalogue for that language
//     - ParseSource(path, source) — pure-Go parser preferred (stdlib or pure-Go
//       library). Attach any AST as ParsedUnit.Tree (opaque any). Do not pull
//       tree-sitter or other CGO deps into the default build.
//     - PrepareProject / ExtractDeps — optional project prep and import edges
//  3. Register the plugin next to golang.Register in engine.DefaultRegistry
//     (or behind a build tag / feature flag if experimental).
//  4. Detectors type-assert unit.Tree only inside the language package
//     (e.g. *mypython.Tree); engine code never imports language AST types.
//
// Go reference implementation: internal/lang/go (goparse = go/parser + go/ast).
// BasePlugin supplies no-op defaults for PrepareProject / ExtractDeps and a
// source-only ParseSource fallback.
type LanguagePlugin interface {
	ID() LanguageID
	Extensions() []string
	Detectors() []Detector
	ParseSource(path, source string) (*ParsedUnit, error)
	PrepareProject(ctx *ScanContext, projectRoots []string)
	ExtractDeps(unit *ParsedUnit, project ProjectContext) []string
}

// BasePlugin provides default implementations for optional LanguagePlugin methods.
type BasePlugin struct{}

// ParseSource returns a source-only ParsedUnit (no AST). Language defaults to Go;
// plugins should override ParseSource when language is not Go.
func (BasePlugin) ParseSource(path, source string) (*ParsedUnit, error) {
	return NewParsedUnit(LanguageGo, path, source), nil
}

func (BasePlugin) PrepareProject(*ScanContext, []string) {}

func (BasePlugin) ExtractDeps(*ParsedUnit, ProjectContext) []string { return nil }

// ParserConfigurer is implemented by plugins that configure a reusable parser.
// Optional; pure-Go plugins typically parse per-file in ParseSource with no pool.
type ParserConfigurer interface {
	ConfigureParser(parser any) error
}

// ProjectPreparer is an optional capability for project-level prep.
// LanguagePlugin already includes PrepareProject; this alias supports type asserts.
type ProjectPreparer interface {
	PrepareProject(ctx *ScanContext, projectRoots []string)
}

// UnitParser is an optional capability for custom source parsing.
// LanguagePlugin already includes ParseSource; this alias supports type asserts.
type UnitParser interface {
	ParseSource(path, source string) (*ParsedUnit, error)
}
