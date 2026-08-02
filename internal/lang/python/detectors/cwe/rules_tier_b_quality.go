package cwe

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	RegisterRule("CWE-1104", detectCWE1104, &MetaCWE1104,
		"import imp", "from imp")
	RegisterRule("CWE-1106", detectCWE1106, &MetaCWE1106,
		"set_cookie")
	// CWE-1108/1121/1124 stay ungated: structural/global-count heuristics fire
	// without a reliable single-token prefilter.
	RegisterRule("CWE-1108", detectCWE1108, &MetaCWE1108)
	RegisterRule("CWE-1121", detectCWE1121, &MetaCWE1121)
	RegisterRule("CWE-1123", detectCWE1123, &MetaCWE1123,
		"types.FunctionType")
	RegisterRule("CWE-1124", detectCWE1124, &MetaCWE1124)
	RegisterRule("CWE-1220", detectCWE1220, &MetaCWE1220,
		".objects.get", "request.args", "request.view_args")
	RegisterRule("CWE-1265", detectCWE1265, &MetaCWE1265,
		"with lock")
	RegisterRule("CWE-1284", detectCWE1284, &MetaCWE1284,
		"request.args", "request.form")
	RegisterRule("CWE-1285", detectCWE1285, &MetaCWE1285,
		"request.args", "request.form")
	RegisterRule("CWE-1287", detectCWE1287, &MetaCWE1287,
		"request.get_json", ".execute(")
	RegisterRule("CWE-1288", detectCWE1288, &MetaCWE1288,
		"request.get_json", "expected_count")
	RegisterRule("CWE-1322", detectCWE1322, &MetaCWE1322,
		"async def", "time.sleep")
	RegisterRule("CWE-1339", detectCWE1339, &MetaCWE1339,
		"float(", "request.args", "request.form")
	RegisterRule("CWE-1341", detectCWE1341, &MetaCWE1341,
		".close()")
}

var (
	pyTierBMagicCookieRE = regexp.MustCompile(`(?is)set_cookie\s*\([^\n]*max_age\s*=\s*(?:3600|86400)\b`)
	pyTierBUserLookupRE  = regexp.MustCompile(`(?is)\.objects\.get\s*\(\s*id\s*=\s*request\.(?:args|view_args)`)
	pyTierBReentrantRE   = regexp.MustCompile(`(?is)with\s+lock\s*:[\s\S]{0,220}with\s+lock\s*:`)
	pyTierBLimitRE       = regexp.MustCompile(`(?is)limit\s*=\s*int\s*\(\s*request\.(?:args|form)[\s\S]{0,240}limit\s*=\s*limit\b`)
	pyTierBIndexRE       = regexp.MustCompile(`(?is)\w+\s*\[\s*int\s*\(\s*request\.(?:args|form)`)
	pyTierBTypeRE        = regexp.MustCompile(`(?is)request\.get_json\s*\(\s*\)[\s\S]{0,220}\.execute\s*\(`)
	pyTierBConsistencyRE = regexp.MustCompile(`(?is)expected_count\s*=\s*request\.get_json\s*\([^\n]*\)[\s\S]{0,260}items\s*=\s*request\.get_json\s*\(`)
	pyTierBAsyncSleepRE  = regexp.MustCompile(`(?is)async\s+def\s+\w+\s*\([^)]*\)\s*:[\s\S]{0,260}time\.sleep\s*\(`)
	pyTierBFloatMoneyRE  = regexp.MustCompile(`(?is)float\s*\(\s*request\.(?:args|form)\.get\s*\([^\n]*(?:price|amount|balance)`)
	pyControlHeaderRE    = regexp.MustCompile(`^(?:async\s+)?(?:if|elif|else|for|while|try|except|finally|with|match|case)\b.*:\s*$`)
	pyTierBImpImportRE   = regexp.MustCompile(`(?im)^\s*(?:import|from)\s+imp\b`)
)

func detectCWE1104(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(facts, unit.Source, pyTierBImpImportRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE1104, start, "unmaintained imp component is imported", confidence74, out)
	}
}

func detectCWE1106(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(facts, unit.Source, pyTierBMagicCookieRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE1106, start, "security-sensitive cookie lifetime uses an unexplained magic number", confidence68, out)
	}
}

func detectCWE1108(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	for _, fn := range facts.Functions() {
		if strings.Count(facts.codeMask(fn.body, fn.bodyStart), "global ") >= 3 {
			emitTierBFinding(unit, &MetaCWE1108, fn.bodyStart, "function relies on multiple global variables", confidence70, out)
			return
		}
	}
}

func detectCWE1121(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	for _, fn := range facts.Functions() {
		code := facts.codeMask(fn.body, fn.bodyStart)
		branches := strings.Count(code, "if ") + strings.Count(code, "elif ") + strings.Count(code, "for ") + strings.Count(code, "while ") + strings.Count(code, "except ")
		if branches >= minimumRouteBranches {
			emitTierBFinding(unit, &MetaCWE1121, fn.bodyStart, "function has at least twelve visible control-flow branches", confidence70, out)
			return
		}
	}
}

func detectCWE1123(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	for _, call := range findCalls(facts, unit.Source, "types.FunctionType") {
		emitTierBFinding(unit, &MetaCWE1123, call.Start, "function object is constructed dynamically from code", confidence80, out)
		return
	}
}

func detectCWE1124(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	type controlFrame struct{ indent int }
	var frames []controlFrame
	offset := 0
	for _, line := range strings.Split(facts.Masked, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			offset += len(line) + 1
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		for len(frames) > 0 && indent <= frames[len(frames)-1].indent {
			frames = frames[:len(frames)-1]
		}
		if pyControlHeaderRE.MatchString(trimmed) {
			// Exception-handler clauses belong to the preceding try control flow;
			// counting them as an additional nesting level overstates ordinary
			// recovery and retry paths.
			if !strings.HasPrefix(trimmed, "except") && !strings.HasPrefix(trimmed, "finally") {
				frames = append(frames, controlFrame{indent: indent})
			}
			offset += len(line) + 1
			continue
		}
		if len(frames) >= 6 && isExecutableNestedStatement(trimmed) {
			emitTierBFinding(unit, &MetaCWE1124, offset, "executable statement is nested at least six control-flow levels", confidence70, out)
			return
		}
		offset += len(line) + 1
	}
}

func isExecutableNestedStatement(line string) bool {
	if line == "" || strings.HasPrefix(line, "@") || strings.HasPrefix(line, "def ") || strings.HasPrefix(line, "class ") {
		return false
	}
	if strings.HasPrefix(line, "return ") || strings.HasPrefix(line, "raise ") || strings.HasPrefix(line, "break") || strings.HasPrefix(line, "continue") {
		return true
	}
	return strings.Contains(line, "=") || strings.Contains(line, "(")
}

func detectCWE1220(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(facts, unit.Source, pyTierBUserLookupRE); start >= 0 && !strings.Contains(unit.Source[start:start+minInt(userLookupContextWindow, len(unit.Source)-start)], "request.user") {
		emitTierBFinding(unit, &MetaCWE1220, start, "object is fetched by request identifier without an owner constraint", confidence76, out)
	}
}

func detectCWE1265(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(facts, unit.Source, pyTierBReentrantRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE1265, start, "non-reentrant lock is acquired again in a nested scope", confidence80, out)
	}
}

func detectCWE1284(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(facts, unit.Source, pyTierBLimitRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE1284, start, "request quantity controls a data-access limit without a visible bound", confidence76, out)
	}
}

func detectCWE1285(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(facts, unit.Source, pyTierBIndexRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE1285, start, "request value is converted directly into a collection index", confidence78, out)
	}
}

func detectCWE1287(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(facts, unit.Source, pyTierBTypeRE); start >= 0 && !strings.Contains(unit.Source, "isinstance(") {
		emitTierBFinding(unit, &MetaCWE1287, start, "JSON request object is used for execution without a visible type check", confidence70, out)
	}
}

func detectCWE1288(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(facts, unit.Source, pyTierBConsistencyRE); start >= 0 && !strings.Contains(unit.Source, "len(items)") {
		emitTierBFinding(unit, &MetaCWE1288, start, "request count and item collection are used without a consistency check", confidence70, out)
	}
}

func detectCWE1322(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(facts, unit.Source, pyTierBAsyncSleepRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE1322, start, "blocking sleep runs inside an async function", confidence84, out)
	}
}

func detectCWE1339(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstLiteralMatchStartIfContains(facts, unit, pyTierBFloatMoneyRE,
		"price", "amount", "balance"); start >= 0 {
		emitTierBFinding(unit, &MetaCWE1339, start, "request monetary value is parsed as a binary floating-point number", confidence80, out)
	}
}

func detectCWE1341(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	masked := ""
	if facts != nil {
		masked = facts.codeMask(unit.Source, fragStartHint(facts, unit.Source))
	} else if unit != nil {
		masked = pythonCodeMask(unit.Source)
	}
	if start := sameHandleDoubleCloseStart(masked); start >= 0 {
		emitTierBFinding(unit, &MetaCWE1341, start, "same resource handle is released twice", confidence82, out)
	}
}

func emitTierBFinding(unit *core.ParsedUnit, meta *rules.RuleMetadata, offset int, message string, confidence float32, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	line, col := unit.LineCol(offset)
	rules.PushFindingWithConfidence(meta, unitFile(unit), line, col, message, confidence, out)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
