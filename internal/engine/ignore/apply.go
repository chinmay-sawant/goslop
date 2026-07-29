package ignore

import (
	"strings"

	"github.com/chinmay/goslop/internal/rules"
)

const suppressedTag = " (suppressed)"

// ApplyOptions controls how suppressions are reported.
type ApplyOptions struct {
	// ShowIgnored keeps suppressed findings (severity forced to Info, Suppressed=true).
	ShowIgnored bool
}

// Apply applies file-level and inline ignore directives to findings.
// Returns the number of findings suppressed (removed or tagged).
func Apply(source string, findings []rules.Finding, opts ApplyOptions) ([]rules.Finding, int) {
	if len(findings) == 0 {
		return findings, 0
	}
	fileIgnore, hasFile := ParseFileIgnore(source)
	if !opts.ShowIgnored && hasFile && fileIgnore.IsAll() {
		n := len(findings)
		return findings[:0], n
	}

	suppressed := 0
	if hasFile {
		var n int
		findings, n = applyPredicate(findings, func(f *rules.Finding) bool {
			return fileIgnore.Matches(f.RuleID)
		}, opts.ShowIgnored)
		suppressed += n
	}

	inline := ParseInlineIgnores(source)
	if len(inline) == 0 {
		return findings, suppressed
	}
	var n int
	findings, n = applyPredicate(findings, func(f *rules.Finding) bool {
		d, ok := inline[f.Line]
		if !ok {
			return false
		}
		return d.Matches(f.RuleID)
	}, opts.ShowIgnored)
	return findings, suppressed + n
}

func applyPredicate(findings []rules.Finding, shouldSuppress func(*rules.Finding) bool, showIgnored bool) ([]rules.Finding, int) {
	if len(findings) == 0 {
		return findings, 0
	}
	suppressed := 0
	out := make([]rules.Finding, 0, len(findings))
	for i := range findings {
		f := findings[i]
		if !shouldSuppress(&f) {
			out = append(out, f)
			continue
		}
		suppressed++
		if showIgnored {
			f.Severity = rules.SeverityInfo
			f.Suppressed = true
			if !strings.HasSuffix(f.Message, suppressedTag) {
				f.Message += suppressedTag
			}
			out = append(out, f)
		}
	}
	return out, suppressed
}
