package cwe

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	// CWE-396/397 are gated on the exact exception tokens they match.
	// Remaining platform rules use FN-safe structural tokens.
	RegisterRule("CWE-396", detectCWE396, &MetaCWE396,
		"except Exception", "except BaseException")
	RegisterRule("CWE-397", detectCWE397, &MetaCWE397,
		"raise Exception", "raise BaseException")
	RegisterRule("CWE-478", detectCWE478, &MetaCWE478,
		"match ")
	RegisterRule("CWE-252", detectCWE252, &MetaCWE252,
		"subprocess.run", "subprocess.call", "os.system")
	RegisterRule("CWE-390", detectCWE390, &MetaCWE390,
		"except")
	RegisterRule("CWE-584", detectCWE584, &MetaCWE584,
		"finally")
}

var (
	pyGenericExceptLineRE = regexp.MustCompile(`^[\t ]*except\s+(?:Exception|BaseException)(?:\s+as\s+[A-Za-z_][A-Za-z0-9_]*)?\s*:`)
	pyGenericRaiseRE      = regexp.MustCompile(`(?m)^[\t ]*raise\s+(?:Exception|BaseException)(?:\s*\(|\b)`)
	pyMatchStartRE        = regexp.MustCompile(`^[\t ]*match\s+[^:\n]+:\s*$`)
	pyCaseStartRE         = regexp.MustCompile(`^[\t ]*case\s+[^:\n]+:\s*$`)
	pyDefaultCaseRE       = regexp.MustCompile(`^[\t ]*case\s+_\s*(?:if\s+[^:\n]+)?\s*:\s*$`)
	pyExceptStartRE       = regexp.MustCompile(`^[\t ]*except(?:\s+[^:\n]+)?\s*:\s*$`)
	pyFinallyStartRE      = regexp.MustCompile(`^[\t ]*finally\s*:\s*$`)
	pyReturnLineRE        = regexp.MustCompile(`^[\t ]*return\b`)
)

// suiteSurfacesFailureMasked reports whether an except suite re-raises,
// logs with traceback, forwards via set_exception, or records into an
// error/result field — i.e. does not hide the failure.
func suiteSurfacesFailureMasked(lines []pyMaskedLine, exceptIdx int) bool {
	if exceptIdx < 0 || exceptIdx >= len(lines) {
		return false
	}
	exceptIndent := lines[exceptIdx].indent
	for _, body := range lines[exceptIdx+1:] {
		trimmed := strings.TrimSpace(body.text)
		if trimmed == "" {
			continue
		}
		if body.indent <= exceptIndent {
			break
		}
		if suiteLineSurfacesFailure(trimmed) {
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
	if i := strings.Index(t, ".error"); i >= 0 {
		rest := strings.TrimSpace(t[i+len(".error"):])
		if strings.HasPrefix(rest, "=") {
			return true
		}
	}
	return false
}

// CWE-396 reports only the two Python root exception classes. Specific
// exception handlers intentionally remain outside this narrow heuristic.
// Handlers whose suite surfaces the failure (re-raise, exc_info /
// logger.exception, set_exception, or error-field recording) are skipped —
// they do not hide distinct failure conditions.
func detectCWE396(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if isPythonTestModule(unit) {
		return
	}
	if unit == nil || out == nil {
		return
	}
	if !containsAnyNeedle(unit.Source, "except Exception", "except BaseException") {
		return
	}
	lines := facts.MaskedLines()
	for i, line := range lines {
		if !pyGenericExceptLineRE.MatchString(line.text) {
			continue
		}
		if suiteSurfacesFailureMasked(lines, i) {
			continue
		}
		emitPlatformFinding(unit, &MetaCWE396, line.start, "generic Exception or BaseException handler can hide distinct failure conditions", confidence84, out)
		return
	}
}

// CWE-397 recognizes direct construction or re-raising of Python's generic
// root exception classes. Raising an application-specific exception is safe.
func detectCWE397(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstMatchStartIfContains(facts, unit, pyGenericRaiseRE,
		"raise Exception", "raise BaseException"); start >= 0 {
		emitPlatformFinding(unit, &MetaCWE397, start, "generic Exception or BaseException is raised directly", confidence82, out)
	}
}

// CWE-478 reports a Python match statement only when it has two or more
// immediate case branches and lacks a wildcard case. The indentation-aware
// walk avoids confusing nested match cases with branches of the outer match.
func detectCWE478(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	lines := facts.MaskedLines()
	if start := matchWithoutDefaultStart(lines); start >= 0 {
		emitPlatformFinding(unit, &MetaCWE478, start, "multiple-case match expression has no wildcard default branch", confidence76, out)
	}
}

// CWE-252 is deliberately limited to process calls used as standalone
// statements. Assigned results and subprocess.run(..., check=True) have
// explicit success handling paths and are not reported.
func detectCWE252(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"subprocess.run", "subprocess.call", "os.system"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			if !standaloneCall(facts, unit.Source, call) || (name == "subprocess.run" && hasKwargTrue(call.ArgsText, "check")) {
				continue
			}
			emitPlatformFinding(unit, &MetaCWE252, call.Start, "process call return status is discarded without checking success", confidence82, out)
			return
		}
	}
}

// CWE-390 recognizes only an except clause whose direct body is pass. It does
// not infer whether logging, recovery, re-raising, or a caller's behaviour is
// sufficient error handling.
func detectCWE390(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if start := exceptPassStart(facts.MaskedLines()); start >= 0 {
		emitPlatformFinding(unit, &MetaCWE390, start, "exception is detected but the handler takes no action", confidence82, out)
	}
}

// CWE-584 limits reporting to direct returns in a finally suite. A return in
// an unrelated nested definition is excluded by requiring the suite's direct
// indentation level.
func detectCWE584(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if start := finallyReturnStart(facts.MaskedLines()); start >= 0 {
		emitPlatformFinding(unit, &MetaCWE584, start, "return inside finally can suppress an exception from the protected block", confidence90, out)
	}
}

type pyMaskedLine struct {
	start  int
	text   string
	indent int
}

// buildMaskedPythonLines splits a pre-masked file into line spans.
func buildMaskedPythonLines(masked string) []pyMaskedLine {
	lines := make([]pyMaskedLine, 0, strings.Count(masked, "\n")+1)
	for start := 0; start <= len(masked); {
		end := len(masked)
		if next := strings.IndexByte(masked[start:], '\n'); next >= 0 {
			end = start + next
		}
		line := masked[start:end]
		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		lines = append(lines, pyMaskedLine{start: start, text: line, indent: indent})
		if end == len(masked) {
			break
		}
		start = end + 1
	}
	return lines
}

func matchWithoutDefaultStart(lines []pyMaskedLine) int {
	for i, match := range lines {
		if !pyMatchStartRE.MatchString(match.text) {
			continue
		}
		caseIndent, cases, hasDefault := -1, 0, false
		for _, line := range lines[i+1:] {
			if strings.TrimSpace(line.text) == "" {
				continue
			}
			if line.indent <= match.indent {
				break
			}
			if !pyCaseStartRE.MatchString(line.text) {
				continue
			}
			if caseIndent < 0 {
				caseIndent = line.indent
			}
			if line.indent != caseIndent {
				continue
			}
			cases++
			hasDefault = hasDefault || pyDefaultCaseRE.MatchString(line.text)
		}
		if cases >= 2 && !hasDefault {
			return match.start
		}
	}
	return -1
}

func standaloneCall(facts *PyCweFacts, source string, call callSite) bool {
	lineStart := strings.LastIndex(source[:call.Start], "\n") + 1
	prefix := source[lineStart:call.Start]
	var masked string
	if facts != nil {
		masked = facts.codeMask(prefix, lineStart)
	} else {
		masked = pythonCodeMask(prefix)
	}
	return strings.TrimSpace(masked) == ""
}

func exceptPassStart(lines []pyMaskedLine) int {
	for i, except := range lines {
		if !pyExceptStartRE.MatchString(except.text) {
			continue
		}
		for _, body := range lines[i+1:] {
			trimmed := strings.TrimSpace(body.text)
			if trimmed == "" {
				continue
			}
			if body.indent <= except.indent {
				break
			}
			if body.indent > except.indent && trimmed == "pass" {
				return except.start
			}
			break
		}
	}
	return -1
}

func finallyReturnStart(lines []pyMaskedLine) int {
	for i, finally := range lines {
		if !pyFinallyStartRE.MatchString(finally.text) {
			continue
		}
		bodyIndent := -1
		for _, body := range lines[i+1:] {
			trimmed := strings.TrimSpace(body.text)
			if trimmed == "" {
				continue
			}
			if body.indent <= finally.indent {
				break
			}
			if bodyIndent < 0 {
				bodyIndent = body.indent
			}
			if body.indent == bodyIndent && pyReturnLineRE.MatchString(body.text) {
				return body.start
			}
		}
	}
	return -1
}

func emitPlatformFinding(unit *core.ParsedUnit, meta *rules.RuleMetadata, offset int, message string, confidence float32, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	line, col := unit.LineCol(offset)
	rules.PushFindingWithConfidence(meta, unitFile(unit), line, col, message, confidence, out)
}
