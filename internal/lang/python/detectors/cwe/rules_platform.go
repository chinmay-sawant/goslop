package cwe

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	// These are source-only checks with no SourceIndex gates. A gate would need
	// to cover every valid exception, match, and process-call spelling, so it
	// would be prone to false negatives before the rule can inspect source.
	RegisterRule("CWE-396", detectCWE396, &MetaCWE396)
	RegisterRule("CWE-397", detectCWE397, &MetaCWE397)
	RegisterRule("CWE-478", detectCWE478, &MetaCWE478)
	RegisterRule("CWE-252", detectCWE252, &MetaCWE252)
	RegisterRule("CWE-390", detectCWE390, &MetaCWE390)
	RegisterRule("CWE-584", detectCWE584, &MetaCWE584)
}

var (
	pyGenericExceptRE = regexp.MustCompile(`(?m)^[\t ]*except\s+(?:Exception|BaseException)(?:\s+as\s+[A-Za-z_][A-Za-z0-9_]*)?\s*:`)
	pyGenericRaiseRE  = regexp.MustCompile(`(?m)^[\t ]*raise\s+(?:Exception|BaseException)(?:\s*\(|\b)`)
	pyMatchStartRE    = regexp.MustCompile(`^[\t ]*match\s+[^:\n]+:\s*$`)
	pyCaseStartRE     = regexp.MustCompile(`^[\t ]*case\s+[^:\n]+:\s*$`)
	pyDefaultCaseRE   = regexp.MustCompile(`^[\t ]*case\s+_\s*(?:if\s+[^:\n]+)?\s*:\s*$`)
	pyExceptStartRE   = regexp.MustCompile(`^[\t ]*except(?:\s+[^:\n]+)?\s*:\s*$`)
	pyFinallyStartRE  = regexp.MustCompile(`^[\t ]*finally\s*:\s*$`)
	pyReturnLineRE    = regexp.MustCompile(`^[\t ]*return\b`)
)

// CWE-396 reports only the two Python root exception classes. Specific
// exception handlers intentionally remain outside this narrow heuristic.
func detectCWE396(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstMatchStart(unit, pyGenericExceptRE); start >= 0 {
		emitPlatformFinding(unit, &MetaCWE396, start, "generic Exception or BaseException handler can hide distinct failure conditions", confidence84, out)
	}
}

// CWE-397 recognizes direct construction or re-raising of Python's generic
// root exception classes. Raising an application-specific exception is safe.
func detectCWE397(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstMatchStart(unit, pyGenericRaiseRE); start >= 0 {
		emitPlatformFinding(unit, &MetaCWE397, start, "generic Exception or BaseException is raised directly", confidence82, out)
	}
}

// CWE-478 reports a Python match statement only when it has two or more
// immediate case branches and lacks a wildcard case. The indentation-aware
// walk avoids confusing nested match cases with branches of the outer match.
func detectCWE478(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if start := matchWithoutDefaultStart(unit.Source); start >= 0 {
		emitPlatformFinding(unit, &MetaCWE478, start, "multiple-case match expression has no wildcard default branch", confidence76, out)
	}
}

// CWE-252 is deliberately limited to process calls used as standalone
// statements. Assigned results and subprocess.run(..., check=True) have
// explicit success handling paths and are not reported.
func detectCWE252(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"subprocess.run", "subprocess.call", "os.system"} {
		for _, call := range findCalls(unit.Source, name) {
			if !standaloneCall(unit.Source, call) || (name == "subprocess.run" && hasKwargTrue(call.ArgsText, "check")) {
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
func detectCWE390(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if start := exceptPassStart(unit.Source); start >= 0 {
		emitPlatformFinding(unit, &MetaCWE390, start, "exception is detected but the handler takes no action", confidence82, out)
	}
}

// CWE-584 limits reporting to direct returns in a finally suite. A return in
// an unrelated nested definition is excluded by requiring the suite's direct
// indentation level.
func detectCWE584(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if start := finallyReturnStart(unit.Source); start >= 0 {
		emitPlatformFinding(unit, &MetaCWE584, start, "return inside finally can suppress an exception from the protected block", confidence90, out)
	}
}

type pyMaskedLine struct {
	start  int
	text   string
	indent int
}

func maskedPythonLines(source string) []pyMaskedLine {
	masked := pythonCodeMask(source)
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

func matchWithoutDefaultStart(source string) int {
	for i, match := range maskedPythonLines(source) {
		if !pyMatchStartRE.MatchString(match.text) {
			continue
		}
		caseIndent, cases, hasDefault := -1, 0, false
		for _, line := range maskedPythonLines(source)[i+1:] {
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

func standaloneCall(source string, call callSite) bool {
	lineStart := strings.LastIndex(source[:call.Start], "\n") + 1
	return strings.TrimSpace(pythonCodeMask(source[lineStart:call.Start])) == ""
}

func exceptPassStart(source string) int {
	lines := maskedPythonLines(source)
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

func finallyReturnStart(source string) int {
	lines := maskedPythonLines(source)
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
