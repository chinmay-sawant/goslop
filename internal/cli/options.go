// Package cli parses CodeHound command-line flags into Options.
package cli

// OutputFormat is the reporter selection.
type OutputFormat string

const (
	FormatText  OutputFormat = "text"
	FormatJSON  OutputFormat = "json"
	FormatSARIF OutputFormat = "sarif"
)

// Options holds parsed CLI state for the app layer.
type Options struct {
	// Paths are scan roots (default ["."] when empty after parse for scan).
	Paths []string
	// Profile is recommended|perf|security|style|all (default recommended).
	Profile string
	// Only / Skip are rule ID filters (comma-separated on the CLI).
	Only []string
	Skip []string
	// Format is text|json|sarif.
	Format OutputFormat
	// ListRules prints registered rules and exits.
	ListRules bool
	// IncludeTests includes *_test.* files (excluded by default).
	IncludeTests bool
	// NoCache disables the incremental cache (stored for later wiring).
	NoCache bool
	// Taint enables experimental taint tracking (overrides profile default when set).
	Taint bool
	// NoTaint disables taint even for security profile.
	NoTaint bool
	// TaintDepth is inter-procedural hop budget (1–4; 0 = use profile default).
	TaintDepth int
	// TaintShowPaths attaches hop evidence to findings.
	TaintShowPaths bool
	// Version requests version printing.
	Version bool
	// Command is a top-level subcommand when set ("init").
	Command string
}
