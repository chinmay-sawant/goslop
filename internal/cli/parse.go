package cli

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

// Parse parses args (typically os.Args[1:]) into Options.
// Supports both -flag and --flag forms via the stdlib flag package.
func Parse(args []string) (*Options, error) {
	return ParseWithOutput(args, nil)
}

// ParseWithOutput is like Parse but directs flag usage/errors to w when non-nil.
func ParseWithOutput(args []string, w io.Writer) (*Options, error) {
	opts := &Options{
		Profile: "recommended",
		Format:  FormatText,
	}

	// Subcommand: init (must be the first token).
	if len(args) > 0 && args[0] == "init" {
		opts.Command = "init"
		return opts, nil
	}

	fs := flag.NewFlagSet("codehound", flag.ContinueOnError)
	if w != nil {
		fs.SetOutput(w)
	} else {
		fs.SetOutput(io.Discard)
	}

	var (
		only   string
		skip   string
		format string
	)

	fs.StringVar(&opts.Profile, "profile", "recommended", "product pack: recommended|perf|security|style|all")
	fs.StringVar(&only, "only", "", "only run these rule IDs (comma-separated)")
	fs.StringVar(&skip, "skip", "", "skip these rule IDs (comma-separated)")
	fs.StringVar(&format, "format", "text", "output format: text|json|sarif")
	fs.BoolVar(&opts.ListRules, "list-rules", false, "list registered rules and exit")
	fs.BoolVar(&opts.IncludeTests, "include-tests", false, "include test files (*_test.*) in analysis")
	fs.BoolVar(&opts.NoCache, "no-cache", false, "disable the incremental analysis cache")
	fs.BoolVar(&opts.Taint, "taint", false, "enable experimental taint tracking (CWE-22/78/79/89)")
	fs.BoolVar(&opts.NoTaint, "no-taint", false, "disable taint tracking even under security profile")
	fs.IntVar(&opts.TaintDepth, "taint-depth", 0, "inter-procedural taint depth 1–4 (0 = profile default)")
	fs.BoolVar(&opts.TaintShowPaths, "taint-show-paths", false, "attach taint hop evidence to findings")
	fs.BoolVar(&opts.Version, "version", false, "print version and exit")

	fs.Usage = func() {
		out := fs.Output()
		fmt.Fprintf(out, "Usage: codehound [flags] [PATH...]\n")
		fmt.Fprintf(out, "       codehound init\n\n")
		fmt.Fprintf(out, "Flags:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	opts.Only = splitCSV(only)
	opts.Skip = splitCSV(skip)

	fmtNorm := strings.ToLower(strings.TrimSpace(format))
	switch OutputFormat(fmtNorm) {
	case FormatText, FormatJSON, FormatSARIF:
		opts.Format = OutputFormat(fmtNorm)
	default:
		return nil, fmt.Errorf("invalid -format %q (want text|json|sarif)", format)
	}

	opts.Paths = append([]string(nil), fs.Args()...)
	if len(opts.Paths) == 0 {
		opts.Paths = []string{"."}
	}

	return opts, nil
}

func splitCSV(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
