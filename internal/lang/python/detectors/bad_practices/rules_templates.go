package badpractices

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	RegisterRule("BP-PY-33", detectBPPY33)
	RegisterRule("BP-PY-34", detectBPPY34)
}

var (
	jinjaEnvCallRe  = regexp.MustCompile(`(?:jinja2\.)?Environment\s*\(`)
	autoescapeFalse = regexp.MustCompile(`(?i)\bautoescape\s*=\s*False\b`)
	jinjaSafeFilter = regexp.MustCompile(`\|\s*safe\b`)
)

// BP-PY-33: Jinja2 Environment(autoescape=False).
func detectBPPY33(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-33")
	if isPythonTestFile(unit) {
		return
	}
	src := unit.Source
	if !strings.Contains(src, "Environment") && !strings.Contains(src, "jinja2") {
		return
	}
	if !strings.Contains(src, "autoescape") && !strings.Contains(src, "Environment") {
		return
	}

	// Scan for Environment( calls (possibly multi-line).
	start := 0
	for {
		loc := jinjaEnvCallRe.FindStringIndex(src[start:])
		if loc == nil {
			return
		}
		abs := start + loc[0]
		// Find '(' after Environment
		open := strings.Index(src[abs:], "(")
		if open < 0 {
			start = abs + 1
			continue
		}
		openAbs := abs + open
		// Avoid false positives: not FooEnvironment if we matched Environment at end —
		// regex already anchors Environment as whole via (?:jinja2\.)?Environment
		// Skip if preceded by identifier char (e.g. MyEnvironment)
		if abs > 0 && isIdentByte(src[abs-1]) {
			start = abs + 1
			continue
		}
		inner, ok := callArgsRegion(src, openAbs)
		if !ok {
			// Window fallback
			end := openAbs + templateCallFallbackBytes
			if end > len(src) {
				end = len(src)
			}
			inner = src[openAbs:end]
		}
		if autoescapeFalse.MatchString(inner) {
			pushAt(unit, meta, abs, "Jinja2 Environment has autoescape=False; enable autoescape or select_autoescape for HTML", out)
		}
		start = openAbs + 1
	}
}

// BP-PY-34: Markup(non-literal) or |safe on dynamic template fragments in .py.
func detectBPPY34(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-34")
	if isPythonTestFile(unit) {
		return
	}
	src := unit.Source
	hasMarkup := strings.Contains(src, "Markup") || strings.Contains(src, "markupsafe")
	hasSafe := strings.Contains(src, "|safe") || strings.Contains(src, "| safe")
	if !hasMarkup && !hasSafe {
		return
	}

	// Markup( non-literal )
	if hasMarkup {
		scanDynamicMarkup(unit, meta, src, out)
	}

	if hasSafe {
		scanJinjaSafeFilter(unit, meta, facts, src, out)
	}
}

const templateCallFallbackBytes = 200

func scanDynamicMarkup(unit *core.ParsedUnit, meta *rules.RuleMetadata, src string, out *[]rules.Finding) {
	start := 0
	for {
		// Prefer markupsafe.Markup and bare Markup(
		idx := indexOfMarkupCall(src, start)
		if idx < 0 {
			break
		}
		open := strings.Index(src[idx:], "(")
		if open < 0 {
			break
		}
		openAbs := idx + open
		arg, ok := firstCallArg(src, openAbs)
		if !ok {
			start = idx + 6
			continue
		}
		arg = strings.TrimSpace(arg)
		if arg == "" {
			start = idx + 6
			continue
		}
		// Miss: plain string/bytes literal (non-f)
		if isPlainStringLiteral(arg) {
			start = idx + 6
			continue
		}
		// Hit: f-string, name, call, concat, etc.
		pushAt(unit, meta, idx, "Markup marks dynamic HTML as safe; only use on trusted/sanitized literals", out)
		start = idx + 6
	}
}

func scanJinjaSafeFilter(unit *core.ParsedUnit, meta *rules.RuleMetadata, facts *bpFacts, src string, out *[]rules.Finding) {
	// Python source embedding Jinja |safe (string templates in .py)
	lines := codeLinesFacts(facts, src)
	for _, line := range lines {
		t := line.text
		if !jinjaSafeFilter.MatchString(t) {
			continue
		}
		// Skip pure comments already stripped; skip if only in docs without template context.
		// Flag when appears inside a string-looking template fragment.
		if strings.Contains(t, "{{") || strings.Contains(t, "|safe") || strings.Contains(t, "| safe") {
			// Avoid flagging documentation that says "use |safe carefully" without template braces —
			// still flag |safe as high-signal in .py per plan.
			loc := jinjaSafeFilter.FindStringIndex(t)
			off := line.byte
			if loc != nil {
				off += loc[0]
			}
			pushAt(unit, meta, off, "Jinja |safe disables escaping on a value; ensure the value is trusted", out)
		}
	}
}

func indexOfMarkupCall(src string, start int) int {
	// markupsafe.Markup( or Markup(
	needles := []string{"markupsafe.Markup", "Markup"}
	best := -1
	for _, n := range needles {
		idx := indexOfIdent(src[start:], n)
		if idx < 0 {
			continue
		}
		abs := start + idx
		// Must be followed by optional space and (
		rest := src[abs+len(n):]
		restTrim := strings.TrimLeft(rest, " \t")
		if !strings.HasPrefix(restTrim, "(") {
			continue
		}
		// Prefer longer match (markupsafe.Markup over Markup)
		if best < 0 || abs < best || (abs == best && len(n) > 0) {
			// For Markup, ensure not part of markupsafe.Markup already counted
			if n == "Markup" && abs >= len("markupsafe.") {
				prefix := src[abs-len("markupsafe.") : abs]
				if prefix == "markupsafe." {
					continue
				}
			}
			if best < 0 || abs < best {
				best = abs
			}
		}
	}
	// Also search further if first Markup skipped
	if best >= 0 {
		return best
	}
	// Manual scan for any Markup(
	from := start
	for {
		idx := strings.Index(src[from:], "Markup")
		if idx < 0 {
			return -1
		}
		abs := from + idx
		if abs > 0 && isIdentByte(src[abs-1]) {
			from = abs + 6
			continue
		}
		end := abs + len("Markup")
		if end < len(src) && isIdentByte(src[end]) {
			from = end
			continue
		}
		rest := strings.TrimLeft(src[end:], " \t")
		if strings.HasPrefix(rest, "(") {
			return abs
		}
		from = end
	}
}
