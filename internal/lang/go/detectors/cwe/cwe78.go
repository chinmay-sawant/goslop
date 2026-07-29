package cwe

import (
	"strings"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/lang/go/detectors/sourceutil"
	"github.com/chinmay/codehound/internal/rules"
)

const ruleCWE78 = "CWE-78"

// CWE78Detector flags same-file command injection heuristics around exec.Command.
//
// This is a taint-lite seed: it is not a full inter-procedural taint graph.
// It fires when exec.Command / exec.CommandContext consumes clearly request-derived
// values (or known shell + dynamic command shapes co-located with request sources).
type CWE78Detector struct {
	core.BaseDetector
}

// NewCWE78 returns a CWE-78 detector.
func NewCWE78() *CWE78Detector {
	return &CWE78Detector{}
}

// Language implements core.Detector.
func (d *CWE78Detector) Language() core.LanguageID { return core.LangGo }

// RuleIDs implements core.Detector.
func (d *CWE78Detector) RuleIDs() []string { return []string{ruleCWE78} }

// MetadataFor implements core.Detector.
func (d *CWE78Detector) MetadataFor(ruleID string) *rules.RuleMetadata {
	if ruleID == ruleCWE78 {
		return &MetaCWE78
	}
	return nil
}

// Run implements core.Detector.
func (d *CWE78Detector) Run(ctx *core.ScanContext, unit *core.ParsedUnit, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if ctx != nil && !ctx.Allows(ruleCWE78) {
		return
	}
	src := unit.Source
	if !strings.Contains(src, "exec.Command") {
		return
	}

	tainted := sourceutil.FindTaintedIdents(src)
	hasReq := sourceutil.HasRequestSource(src)
	calls := sourceutil.FindCalls(src, "exec.Command", "exec.CommandContext")
	file := unit.DisplayPath
	if file == "" {
		file = unit.Path
	}

	for _, call := range calls {
		if shouldFlagCommandInjection(call, tainted, hasReq) {
			line, col := unit.LineCol(call.Start)
			rules.PushFindingWithConfidence(
				&MetaCWE78,
				file,
				line,
				col,
				"user-controlled input reaches a shell command execution sink",
				0.75,
				out,
			)
		}
	}
}

func shouldFlagCommandInjection(call sourceutil.CallSite, tainted map[string]struct{}, hasReq bool) bool {
	args := sourceutil.SplitTopLevelArgs(call.ArgsText)
	if len(args) == 0 {
		return false
	}

	// Direct request pattern in any argument.
	for _, a := range args {
		if sourceutil.HasRequestSource(a) {
			return true
		}
		for name := range tainted {
			if sourceutil.ContainsIdent(a, name) {
				return true
			}
		}
	}

	// Shell form: exec.Command("sh"|"bash", "-c", <dynamic>) with request sources present.
	// CommandContext prepends context as arg0.
	cmdArgs := args
	if call.Name == "exec.CommandContext" && len(cmdArgs) > 0 {
		cmdArgs = cmdArgs[1:]
	}
	if !hasShellArgv(cmdArgs) {
		// Without shell argv, only flag when a tainted/request value is already in args (above).
		return false
	}
	// Third argv after shell+flag (index 2) is the command string.
	if len(cmdArgs) < 3 {
		return false
	}
	cmdBody := cmdArgs[2]
	if sourceutil.IsPureStringLiteral(cmdBody) {
		return false
	}
	// Dynamic shell body + any request source in file (fixture-shaped / same-file heuristic).
	if hasReq {
		return true
	}
	// Or dynamic body referencing any identifier (still risky but lower signal without request).
	return false
}

func hasShellArgv(args []string) bool {
	if len(args) < 2 {
		return false
	}
	// compact check across first two args for "sh"/"bash" and "-c"
	compact := sourceutil.CompactWhitespace(strings.Join(args[:2], ","))
	return strings.Contains(compact, `"sh","-c"`) ||
		strings.Contains(compact, `"bash","-c"`) ||
		strings.Contains(compact, "`sh`,`-c`") ||
		strings.Contains(compact, "`bash`,`-c`")
}
