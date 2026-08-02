package badpractices

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/lang/python/pytext"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

const (
	bareExceptClause  = "except:"
	exceptionName     = "Exception"
	baseExceptionName = "BaseException"
	maxSignatureLines = 30
)

func init() {
	RegisterRule("BP-PY-1", detectBPPY1)
	RegisterRule("BP-PY-2", detectBPPY2)
	RegisterRule("BP-PY-4", detectBPPY4)
	RegisterRule("BP-PY-6", detectBPPY6)
	RegisterRule("BP-PY-7", detectBPPY7)
}

// BP-PY-1: bare `except:` or broad `except Exception` / `except BaseException`
// with weak handling (pass / bare continue).
func detectBPPY1(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-1")
	if !facts.has("except") {
		return
	}
	lines := codeLinesFacts(facts, unit.Source)
	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		if t == "" {
			continue
		}
		// Bare except:
		if t == bareExceptClause || strings.HasPrefix(t, bareExceptClause) {
			// "except:" only (no type)
			rest := strings.TrimSpace(strings.TrimPrefix(t, "except"))
			if rest == ":" || strings.HasPrefix(rest, ":") {
				pushAt(unit, meta, line.byte, "bare except: swallows all exceptions; catch specific types", out)
				continue
			}
		}
		// except Exception: / except BaseException: (optional "as x")
		if isBroadExcept(t) {
			if isPythonTestFile(unit) && broadExceptCollectsTestEvidence(lines, i) {
				continue
			}
			// Flag broad Exception/BaseException unless the suite surfaces the
			// failure (re-raise, exc_info / logger.exception, set_exception,
			// or records into an error/result field).
			if suiteSurfacesFailure(lines, i) {
				continue
			}
			pushAt(unit, meta, line.byte, "broad except Exception/BaseException hides failures; catch specific types or re-raise", out)
		}
	}
}

func broadExceptCollectsTestEvidence(lines []codeLine, exceptIdx int) bool {
	if exceptIdx < 0 || exceptIdx >= len(lines) {
		return false
	}
	exceptIndent := indentWidth(lines[exceptIdx].raw)
	collects := false
	for j := exceptIdx + 1; j < len(lines); j++ {
		t := strings.TrimSpace(lines[j].text)
		if t == "" {
			continue
		}
		if indentWidth(lines[j].raw) <= exceptIndent {
			break
		}
		collects = collects || strings.Contains(t, ".append(")
	}
	if !collects {
		return false
	}
	for _, line := range lines[exceptIdx+1:] {
		t := strings.TrimSpace(line.text)
		if strings.HasPrefix(t, "assert ") || strings.Contains(t, ".assert") {
			return true
		}
	}
	return false
}

func isBroadExcept(t string) bool {
	// except Exception: / except Exception as e: / except BaseException...
	t = strings.TrimSpace(t)
	if !strings.HasPrefix(t, "except ") {
		return false
	}
	rest := strings.TrimSpace(strings.TrimPrefix(t, "except "))
	// Strip trailing :
	if i := strings.Index(rest, ":"); i >= 0 {
		rest = strings.TrimSpace(rest[:i])
	}
	// as name
	if i := strings.Index(rest, " as "); i >= 0 {
		rest = strings.TrimSpace(rest[:i])
	}
	return rest == exceptionName || rest == baseExceptionName
}

func suiteReraises(lines []codeLine, exceptIdx int) bool {
	return suiteSurfacesFailure(lines, exceptIdx)
}

// suiteSurfacesFailure reports whether an except suite propagates or records
// the failure instead of swallowing it. Bare/swallowing handlers (pass,
// continue, warn-and-return-default) remain reportable.
func suiteSurfacesFailure(lines []codeLine, exceptIdx int) bool {
	if exceptIdx < 0 || exceptIdx >= len(lines) {
		return false
	}
	exceptIndent := indentWidth(lines[exceptIdx].raw)
	for j := exceptIdx + 1; j < len(lines); j++ {
		raw := lines[j].raw
		t := strings.TrimSpace(lines[j].text)
		if t == "" {
			continue
		}
		ind := indentWidth(raw)
		if ind <= exceptIndent {
			break
		}
		if suiteLineSurfacesFailure(t) {
			return true
		}
	}
	return false
}

func suiteLineSurfacesFailure(t string) bool {
	if t == "raise" || strings.HasPrefix(t, "raise ") {
		return true
	}
	if strings.Contains(t, "exc_info") {
		return true
	}
	if strings.Contains(t, ".exception(") {
		return true
	}
	if strings.Contains(t, "set_exception(") {
		return true
	}
	if strings.Contains(t, "_error_result(") {
		return true
	}
	// result.error = ... / health.error = ... — records failure into a field.
	if i := strings.Index(t, ".error"); i >= 0 {
		rest := strings.TrimSpace(t[i+len(".error"):])
		if strings.HasPrefix(rest, "=") {
			return true
		}
	}
	return false
}

// BP-PY-2: except suite is solely pass (optional comment already stripped).
func detectBPPY2(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-2")
	if !facts.has("except") || !facts.has("pass") {
		return
	}
	lines := codeLinesFacts(facts, unit.Source)
	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		if !strings.HasPrefix(t, "except") || !strings.Contains(t, ":") {
			continue
		}
		// Must look like an except clause (not a variable named except_foo).
		if t != bareExceptClause && !strings.HasPrefix(t, "except ") && !strings.HasPrefix(t, bareExceptClause) {
			continue
		}
		exceptIndent := indentWidth(line.raw)
		// Collect suite statements at greater indent until dedent.
		var suite []string
		for j := i + 1; j < len(lines); j++ {
			st := strings.TrimSpace(lines[j].text)
			if st == "" {
				continue
			}
			ind := indentWidth(lines[j].raw)
			if ind <= exceptIndent {
				break
			}
			// Nested block headers still count as suite content.
			suite = append(suite, st)
			// Only consider immediate suite lines at the first indent level for "solely pass".
			// If we already have more than one statement, not solely pass.
			if len(suite) > 1 {
				break
			}
		}
		if len(suite) == 1 && suite[0] == "pass" {
			pushAt(unit, meta, line.byte, "except handler body is only pass; failures are discarded silently", out)
		}
	}
}

var mutableDefaultRe = regexp.MustCompile(`=\s*(\[\s*\]|\{\s*\}|set\s*\(\s*\)|list\s*\(\s*\)|dict\s*\(\s*\))`)

// BP-PY-4: mutable default arguments [] {} set() list() dict()
func detectBPPY4(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-4")
	if !facts.has("def ") && !facts.has("async def ") {
		return
	}
	src := unit.Source
	// Scan for def / async def signatures (possibly multi-line until ')' before ':').
	lines := codeLinesFacts(facts, unit.Source)
	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		if !strings.HasPrefix(t, "def ") && !strings.HasPrefix(t, "async def ") {
			continue
		}
		// Accumulate signature text until we see ')' that closes params and optional '->' then ':'.
		sig := t
		j := i
		// If already complete on one line
		for !signatureComplete(sig) && j+1 < len(lines) {
			j++
			sig += " " + strings.TrimSpace(lines[j].text)
			// Safety: stop after many lines
			if j-i > maxSignatureLines {
				break
			}
		}
		// Only look inside parentheses of the def.
		open := strings.Index(sig, "(")
		closeParen := strings.LastIndex(sig, ")")
		if open < 0 || closeParen <= open {
			continue
		}
		params := sig[open+1 : closeParen]
		if mutableDefaultRe.MatchString(params) {
			pushAt(unit, meta, line.byte, "mutable default argument is shared across calls; use None and assign inside the body", out)
		}
	}
	_ = src
}

func signatureComplete(sig string) bool {
	// Rough: has '(' and matching ')' and ends with ':' (possibly after return annotation).
	if !strings.Contains(sig, "(") {
		return false
	}
	depth := 0
	inStr := byte(0)
	escape := false
	for i := 0; i < len(sig); i++ {
		c := sig[i]
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
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				// rest should contain ':'
				return strings.Contains(sig[i:], ":")
			}
		}
	}
	return false
}

// BP-PY-6: assert used for request/CLI/authz/path validation in non-test modules.
// Internal invariant asserts (no security/input needle) are intentionally missed.
func detectBPPY6(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-6")
	if isPythonTestFile(unit) {
		return
	}
	if !facts.has("assert ") && !strings.Contains(unit.Source, "assert ") {
		if !strings.Contains(unit.Source, "assert") {
			return
		}
	}
	lines := codeLinesFacts(facts, unit.Source)
	for _, line := range lines {
		t := strings.TrimSpace(line.text)
		if !(strings.HasPrefix(t, "assert ") || t == "assert" || strings.HasPrefix(t, "assert(")) {
			continue
		}
		if !assertLooksLikeRuntimeValidation(t) {
			continue
		}
		pushAt(unit, meta, line.byte, "assert is stripped with python -O; use if + raise for runtime validation", out)
	}
}

func assertLooksLikeRuntimeValidation(line string) bool {
	lower := strings.ToLower(line)
	needles := []string{
		"request.", "request(", "request ",
		"args.", "form.", "files.",
		"is_authenticated", "is_anonymous", "has_perm", "permission", "authorize",
		"csrf", "safe_join",
		"filename", "filepath", "dirname",
		"sys.argv", "click.", "argparse",
		"user_input", "untrusted",
	}
	for _, n := range needles {
		if strings.Contains(lower, n) {
			return true
		}
	}
	// Word-ish tokens to avoid matching "author"/"apathy"/etc.
	for _, re := range []*regexp.Regexp{
		regexp.MustCompile(`\bauth\b`),
		regexp.MustCompile(`\brole\b`),
		regexp.MustCompile(`\btoken\b`),
		regexp.MustCompile(`\bpath\b`),
		regexp.MustCompile(`\bcli\b`),
		regexp.MustCompile(`\bargv\b`),
	} {
		if re.MatchString(lower) {
			return true
		}
	}
	return false
}

// BP-PY-7: builtin open( without with.
// Attribute methods (fitz.open, Image.open, self.open, Client.open, os.open)
// and function definitions (def open) are out of scope. Docstring/comment
// matches are blanked via pytext.Mask before the line scan.
func detectBPPY7(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-7")
	if !facts.has("open(") && !facts.has(".open(") {
		return
	}
	masked := pytext.Mask(unit.Source)
	lines := buildCodeLines(masked)
	for _, line := range lines {
		t := line.text
		if !strings.Contains(t, "open(") {
			continue
		}
		if lineContainsWithOpen(t) {
			continue
		}
		trimmed := strings.TrimSpace(t)
		if strings.HasPrefix(trimmed, "from ") || strings.HasPrefix(trimmed, "import ") {
			continue
		}
		if strings.HasPrefix(trimmed, "def open(") || strings.HasPrefix(trimmed, "async def open(") {
			continue
		}
		if i := indexOfBareOpenCall(t); i >= 0 {
			pushAt(unit, meta, line.byte+i, "open without with risks resource leaks; use a context manager", out)
		}
	}
}

// indexOfBareOpenCall finds builtin open( — not mid-ident and not attribute .open(.
func indexOfBareOpenCall(line string) int {
	start := 0
	for {
		i := strings.Index(line[start:], "open(")
		if i < 0 {
			return -1
		}
		abs := start + i
		if abs > 0 {
			prev := line[abs-1]
			if prev == '.' || isIdentByte(prev) {
				start = abs + 4
				continue
			}
		}
		return abs
	}
}
