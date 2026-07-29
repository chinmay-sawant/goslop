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
	// NoCache disables the incremental analysis cache.
	NoCache bool
	// CacheDir overrides the default .codehound-cache directory.
	CacheDir string
	// RebuildCache purges the cache directory before the scan.
	RebuildCache bool
	// PruneCache prunes stale cache entries for the given paths and exits.
	PruneCache bool
	// NoBaseline disables baseline loading/filtering.
	NoBaseline bool
	// BaselineFile is an explicit baseline path (default: discover).
	BaselineFile string
	// ShowIgnored keeps findings suppressed by codehound-ignore directives.
	ShowIgnored bool
	// ShowBaselined keeps findings present in the baseline.
	ShowBaselined bool
	// Version requests version printing.
	Version bool
	// Command is a top-level subcommand when set ("init").
	Command string
}
