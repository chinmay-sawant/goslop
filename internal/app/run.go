// Package app orchestrates the CodeHound CLI.
package app

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/chinmay/codehound/internal/cli"
	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/engine"
	"github.com/chinmay/codehound/internal/reporting"
	"github.com/chinmay/codehound/internal/rules"
)

// Run is the CLI entry used by cmd/codehound.
func Run(args []string) error {
	return run(args, os.Stdout, os.Stderr)
}

func run(args []string, stdout, stderr io.Writer) error {
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}

	opts, err := cli.ParseWithOutput(args, stderr)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return &ExitCodeError{Code: ExitConfig, Err: err}
	}

	if opts.Command == "init" {
		return runInit()
	}
	if opts.Version {
		fmt.Fprintln(stdout, Version)
		return nil
	}
	if opts.ListRules {
		return listRules(stdout)
	}

	profile, ok := core.ParseProfile(opts.Profile)
	if !ok {
		return &ExitCodeError{
			Code: ExitConfig,
			Err:  fmt.Errorf("unknown profile %q", opts.Profile),
		}
	}

	ctx := core.NewScanContext(profile, opts.Only, opts.Skip)
	ctx.IncludeTests = opts.IncludeTests
	ctx.NoCache = opts.NoCache

	reg := engine.DefaultRegistry()
	analyzer := engine.NewAnalyzer(ctx, reg)
	res, err := analyzer.AnalyzePaths(opts.Paths)
	if err != nil {
		return &ExitCodeError{Code: ExitInternal, Err: err}
	}

	findings := res.Findings
	if findings == nil {
		findings = []rules.Finding{}
	}

	rep, err := reporting.New(string(opts.Format))
	if err != nil {
		return &ExitCodeError{Code: ExitConfig, Err: err}
	}
	switch r := rep.(type) {
	case reporting.JSONReporter:
		r.Version = Version
		rep = r
	case reporting.SARIFReporter:
		r.Version = Version
		rep = r
	}
	if err := rep.Write(findings, stdout); err != nil {
		return &ExitCodeError{Code: ExitInternal, Err: err}
	}

	if res.ShouldFail(ctx.FailPolicy) {
		return &ExitCodeError{Code: ExitFailing}
	}
	return nil
}

func listRules(w io.Writer) error {
	reg := engine.DefaultRegistry()
	ids := reg.AllRuleIDs()
	if len(ids) == 0 {
		fmt.Fprintln(w, "no rules registered")
		return nil
	}
	for _, d := range reg.Detectors() {
		for _, id := range d.RuleIDs() {
			title := ""
			if meta := d.MetadataFor(id); meta != nil {
				title = meta.Title
			}
			if title != "" {
				fmt.Fprintf(w, "%s\t%s\n", id, title)
			} else {
				fmt.Fprintln(w, id)
			}
		}
	}
	return nil
}
