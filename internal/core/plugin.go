package core

// ProjectContext is language-neutral project context for dependency extraction.
type ProjectContext struct {
	Root string
}

// LanguagePlugin is a per-language backend: parse sources and supply detectors.
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
