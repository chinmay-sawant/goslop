package perf

import (
	"strings"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/lang/go/detectors/sourceutil"
	"github.com/chinmay/codehound/internal/rules"
)

const rulePERF116 = "PERF-116"

// PERF116Detector flags strings.Index(...) compared to -1 (use strings.Contains).
type PERF116Detector struct {
	core.BaseDetector
}

// NewPERF116 returns a PERF-116 detector.
func NewPERF116() *PERF116Detector {
	return &PERF116Detector{}
}

// Language implements core.Detector.
func (d *PERF116Detector) Language() core.LanguageID { return core.LangGo }

// RuleIDs implements core.Detector.
func (d *PERF116Detector) RuleIDs() []string { return []string{rulePERF116} }

// MetadataFor implements core.Detector.
func (d *PERF116Detector) MetadataFor(ruleID string) *rules.RuleMetadata {
	if ruleID == rulePERF116 {
		return &MetaPERF116
	}
	return nil
}

// Run implements core.Detector.
func (d *PERF116Detector) Run(ctx *core.ScanContext, unit *core.ParsedUnit, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if ctx != nil && !ctx.Allows(rulePERF116) {
		return
	}
	src := unit.Source
	if !strings.Contains(src, "strings.Index") {
		return
	}
	file := unit.DisplayPath
	if file == "" {
		file = unit.Path
	}

	for _, call := range sourceutil.FindCalls(src, "strings.Index") {
		// Surrounding statement window: include a bit before and after the call.
		windowStart := call.Start
		if windowStart > 8 {
			// allow parenthesized forms; mostly we care about after the call
			windowStart = call.Start
		}
		windowEnd := call.End + 32
		if windowEnd > len(src) {
			windowEnd = len(src)
		}
		// Also peek slightly after for " != -1"
		window := src[windowStart:windowEnd]
		// Expand to end of line for typical `return strings.Index(s, sub) != -1`
		if nl := strings.IndexByte(src[call.Start:], '\n'); nl >= 0 {
			lineEnd := call.Start + nl
			if lineEnd > windowEnd {
				window = src[windowStart:lineEnd]
			}
		}
		if !strings.Contains(window, "!= -1") && !strings.Contains(window, "== -1") &&
			!strings.Contains(window, "!=-1") && !strings.Contains(window, "==-1") {
			continue
		}
		line, col := unit.LineCol(call.Start)
		rules.PushFinding(
			&MetaPERF116,
			file,
			line,
			col,
			"strings.Index(s, sub) compared to -1 should be strings.Contains(s, sub)",
			out,
		)
	}
}
