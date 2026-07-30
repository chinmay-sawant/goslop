package cwe

import (
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/lang/go/detectors/sourceutil"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// CWE78Detector is a single-rule adapter retained for seed unit tests.
// Production registration uses GoCweScan via RegisterRule("CWE-78", ...).
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
func (d *CWE78Detector) RuleIDs() []string { return []string{"CWE-78"} }

// MetadataFor implements core.Detector.
func (d *CWE78Detector) MetadataFor(ruleID string) *rules.RuleMetadata {
	if ruleID == "CWE-78" {
		return &MetaCWE78
	}
	return nil
}

// Run implements core.Detector.
// When taint is enabled, the taint detector owns CWE-78 (seed is the fallback).
func (d *CWE78Detector) Run(ctx *core.ScanContext, unit *core.ParsedUnit, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if ctx != nil && ctx.TaintEnabled {
		return
	}
	if ctx != nil && !ctx.Allows("CWE-78") {
		return
	}
	detectCWE78(unit, BuildFacts(unit), out)
}

func shouldFlagCommandInjection(call sourceutil.CallSite, tainted map[string]struct{}, hasReq bool) bool {
	args := sourceutil.SplitTopLevelArgs(call.ArgsText)
	if len(args) == 0 {
		return false
	}

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

	cmdArgs := args
	if call.Name == "exec.CommandContext" && len(cmdArgs) > 0 {
		cmdArgs = cmdArgs[1:]
	}
	if !hasShellArgv(cmdArgs) {
		return false
	}
	if len(cmdArgs) < 3 {
		return false
	}
	cmdBody := cmdArgs[2]
	if sourceutil.IsPureStringLiteral(cmdBody) {
		return false
	}
	return hasReq
}

func hasShellArgv(args []string) bool {
	if len(args) < 2 {
		return false
	}
	compact := sourceutil.CompactWhitespace(strings.Join(args[:2], ","))
	return strings.Contains(compact, `"sh","-c"`) ||
		strings.Contains(compact, `"bash","-c"`) ||
		strings.Contains(compact, "`sh`,`-c`") ||
		strings.Contains(compact, "`bash`,`-c`")
}
