package badpractices

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
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
		if t == "except:" || strings.HasPrefix(t, "except:") {
			// "except:" only (no type)
			rest := strings.TrimSpace(strings.TrimPrefix(t, "except"))
			if rest == ":" || strings.HasPrefix(rest, ":") {
				pushAt(unit, meta, line.byte, "bare except: swallows all exceptions; catch specific types", out)
				continue
			}
		}
		// except Exception: / except BaseException: (optional "as x")
		if isBroadExcept(t) {
			// Flag broad Exception/BaseException unless suite clearly re-raises.
			if suiteReraises(lines, i) {
				continue
			}
			pushAt(unit, meta, line.byte, "broad except Exception/BaseException hides failures; catch specific types or re-raise", out)
		}
	}
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
	return rest == "Exception" || rest == "BaseException"
}

func suiteReraises(lines []codeLine, exceptIdx int) bool {
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
		if t == "raise" || strings.HasPrefix(t, "raise ") {
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
		if t != "except:" && !strings.HasPrefix(t, "except ") && !strings.HasPrefix(t, "except:") {
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
			if j-i > 30 {
				break
			}
		}
		// Only look inside parentheses of the def.
		open := strings.Index(sig, "(")
		close := strings.LastIndex(sig, ")")
		if open < 0 || close <= open {
			continue
		}
		params := sig[open+1 : close]
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

// BP-PY-6: assert in non-test modules.
func detectBPPY6(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-6")
	if isPythonTestFile(unit) {
		return
	}
	if !facts.has("assert ") && !strings.Contains(unit.Source, "assert ") {
		// also bare "assert(" unusual
		if !strings.Contains(unit.Source, "assert") {
			return
		}
	}
	lines := codeLinesFacts(facts, unit.Source)
	for _, line := range lines {
		t := strings.TrimSpace(line.text)
		// assert expr / assert expr, msg
		if strings.HasPrefix(t, "assert ") || t == "assert" || strings.HasPrefix(t, "assert(") {
			pushAt(unit, meta, line.byte, "assert is stripped with python -O; use if + raise for runtime validation", out)
		}
	}
}

// BP-PY-7: open( / .open( without with.
func detectBPPY7(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-7")
	if !facts.has("open(") && !facts.has(".open(") {
		return
	}
	lines := codeLinesFacts(facts, unit.Source)
	for _, line := range lines {
		t := line.text
		if !strings.Contains(t, "open(") {
			continue
		}
		if lineContainsWithOpen(t) {
			continue
		}
		// Skip imports / comments already stripped.
		trimmed := strings.TrimSpace(t)
		if strings.HasPrefix(trimmed, "from ") || strings.HasPrefix(trimmed, "import ") {
			continue
		}
		// Find open( occurrences not as part of another ident (e.g. reopen is rare; open is builtin).
		// Flag assignment or bare call: f = open(...), open(...).read()
		// Skip if the open is only inside a string (comment-stripped already; strings still present).
		if openCallOutsideString(t) {
			// Locate byte offset of open(
			off := line.byte + strings.Index(t, "open(")
			// Prefer .open( or open(
			if i := strings.Index(t, ".open("); i >= 0 {
				off = line.byte + i
			} else if i := indexOfIdent(t, "open"); i >= 0 {
				off = line.byte + i
			}
			pushAt(unit, meta, off, "open without with risks resource leaks; use a context manager", out)
		}
	}
}

func openCallOutsideString(line string) bool {
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
		// match open(
		if c == 'o' && i+5 <= len(line) && line[i:i+5] == "open(" {
			if i == 0 || !isIdentByte(line[i-1]) {
				return true
			}
		}
		// .open(
		if c == '.' && i+6 <= len(line) && line[i:i+6] == ".open(" {
			return true
		}
	}
	return false
}
