package badpractices

import (
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	metaByID["BP-PY-46"] = &rules.RuleMetadata{
		ID: "BP-PY-46", Title: "print Debugging In Library Code",
		Description: "`print` is used for operational logging in non-script modules.",
		Severity:    rules.SeverityInfo, Pack: rules.PackBadPractice,
		Fix: "Use the logging module; keep print under `if __name__ == \"__main__\"` for CLIs.",
	}
	metaByID["BP-PY-47"] = &rules.RuleMetadata{
		ID: "BP-PY-47", Title: "logging With String Format Before Logger",
		Description: "Log messages are eagerly formatted with f-strings instead of lazy `%s`/`{}` args.",
		Severity:    rules.SeverityInfo, Pack: rules.PackBadPractice,
		Fix: "Prefer lazy logging: logger.info(\"x=%s\", val) so formatting is skipped when disabled.",
	}
	RegisterRule("BP-PY-46", detectBPPY46)
	RegisterRule("BP-PY-47", detectBPPY47)
}

// BP-PY-46: print( in library code (not under __main__ guard, not tests).
func detectBPPY46(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-46")
	if isPythonTestFile(unit) {
		return
	}
	if isRequirementsPath(unit) {
		return
	}
	if !facts.has("print(") && !strings.Contains(unit.Source, "print(") {
		return
	}
	lines := codeLinesFacts(facts, unit.Source)
	// Track whether we are inside if __name__ == "__main__": block.
	mainIndent := -1
	for _, line := range lines {
		t := strings.TrimSpace(line.text)
		if t == "" {
			continue
		}
		ind := indentWidth(line.raw)
		if mainIndent >= 0 && ind <= mainIndent {
			mainIndent = -1
		}
		if isMainGuard(t) {
			mainIndent = ind
			continue
		}
		if mainIndent >= 0 && ind > mainIndent {
			// Under main guard — skip print.
			continue
		}
		// Flag print( calls that are not import-like.
		if !strings.Contains(t, "print(") {
			continue
		}
		// Skip comments already stripped; skip if print is only in a string (cheap check).
		if printCallOutsideString(t) {
			off := line.byte
			if i := indexOfIdent(t, "print"); i >= 0 {
				// Use offset in raw text when possible
				if j := strings.Index(line.text, "print"); j >= 0 {
					off = line.byte + j
				}
			}
			pushAt(unit, meta, off,
				"print used in library code; prefer logging (or keep print under if __name__ == \"__main__\")",
				out)
		}
	}
}

func isMainGuard(t string) bool {
	// if __name__ == "__main__":  and variants with single quotes / spaces
	if !strings.HasPrefix(t, "if ") {
		return false
	}
	if !strings.Contains(t, "__name__") {
		return false
	}
	if !strings.Contains(t, "__main__") {
		return false
	}
	return strings.HasSuffix(strings.TrimSpace(t), ":") || strings.Contains(t, ":")
}

func printCallOutsideString(line string) bool {
	inStr := byte(0)
	escape := false
	for i := 0; i < len(line); i++ {
		c := line[i]
		if inStr != 0 {
			if escape {
				escape = false
				continue
			}
			if c == '\\' {
				escape = true
				continue
			}
			if c == inStr {
				inStr = 0
			}
			continue
		}
		if c == '"' || c == '\'' {
			inStr = c
			continue
		}
		if c == 'p' && i+6 <= len(line) && line[i:i+6] == "print(" {
			if i == 0 || !isIdentByte(line[i-1]) {
				return true
			}
		}
	}
	return false
}

// BP-PY-47: logger.*/logging.* with eager f-string or .format( as first arg.
func detectBPPY47(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-47")
	if isRequirementsPath(unit) {
		return
	}
	if !facts.hasAny("logger.", "logging.") &&
		!strings.Contains(unit.Source, "logger.") && !strings.Contains(unit.Source, "logging.") {
		return
	}
	src := unit.Source
	methods := []string{
		"debug", "info", "warning", "error", "critical", "exception", "log",
	}
	// Patterns: logger.METHOD( or logging.METHOD(
	prefixes := []string{"logger.", "logging."}
	for _, pref := range prefixes {
		for _, m := range methods {
			needle := pref + m + "("
			start := 0
			for {
				idx := strings.Index(src[start:], needle)
				if idx < 0 {
					break
				}
				abs := start + idx
				open := abs + len(needle) - 1 // '('
				if open >= len(src) || src[open] != '(' {
					start = abs + len(needle)
					continue
				}
				arg, ok := firstCallArg(src, open)
				if !ok {
					start = abs + len(needle)
					continue
				}
				if isEagerLogFormat(arg) {
					pushAt(unit, meta, abs,
						"log message eagerly formatted (f-string/.format); prefer lazy %s args so work is skipped when disabled",
						out)
				}
				start = abs + len(needle)
			}
		}
	}
}

func isEagerLogFormat(arg string) bool {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		return false
	}
	// f-string / rf / fr prefixes
	lower := arg
	// Strip b/u not relevant; f/F and combinations
	a := arg
	// Common f-string prefixes
	for _, p := range []string{"f", "F", "rf", "Rf", "rF", "RF", "fr", "Fr", "fR", "FR"} {
		if strings.HasPrefix(a, p+"\"") || strings.HasPrefix(a, p+"'") {
			return true
		}
	}
	// .format( on the first arg expression: "...".format( or 'x{}'.format(
	if strings.Contains(arg, ".format(") {
		return true
	}
	// % formatting done before call: ("x %s" % val) — paren wrapped
	// Conservative: only flag if top-level % between quotes and rest.
	if isPercentFormattedArg(arg) {
		return true
	}
	_ = lower
	return false
}

func isPercentFormattedArg(arg string) bool {
	// ("msg %s" % x) or 'msg %s' % x
	// Not a plain "msg %s" literal (lazy-style first arg without %).
	depth := 0
	inStr := byte(0)
	escape := false
	for i := 0; i < len(arg); i++ {
		c := arg[i]
		if inStr != 0 {
			if escape {
				escape = false
				continue
			}
			if c == '\\' {
				escape = true
				continue
			}
			if c == inStr {
				inStr = 0
			}
			continue
		}
		switch c {
		case '"', '\'':
			inStr = c
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			if depth > 0 {
				depth--
			}
		case '%':
			if depth == 0 {
				// binary % operator outside string
				// ensure not %% and not start of f-string nonsense
				if i+1 < len(arg) && arg[i+1] == '%' {
					i++
					continue
				}
				// Must have something after that looks like RHS
				rest := strings.TrimSpace(arg[i+1:])
				if rest != "" {
					return true
				}
			}
		}
	}
	return false
}
