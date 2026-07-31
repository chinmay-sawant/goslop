package cwe

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	RegisterRule("CWE-1104", detectCWE1104, &MetaCWE1104)
	RegisterRule("CWE-1106", detectCWE1106, &MetaCWE1106)
	RegisterRule("CWE-1108", detectCWE1108, &MetaCWE1108)
	RegisterRule("CWE-1121", detectCWE1121, &MetaCWE1121)
	RegisterRule("CWE-1123", detectCWE1123, &MetaCWE1123)
	RegisterRule("CWE-1124", detectCWE1124, &MetaCWE1124)
	RegisterRule("CWE-1220", detectCWE1220, &MetaCWE1220)
	RegisterRule("CWE-1265", detectCWE1265, &MetaCWE1265)
	RegisterRule("CWE-1284", detectCWE1284, &MetaCWE1284)
	RegisterRule("CWE-1285", detectCWE1285, &MetaCWE1285)
	RegisterRule("CWE-1287", detectCWE1287, &MetaCWE1287)
	RegisterRule("CWE-1288", detectCWE1288, &MetaCWE1288)
	RegisterRule("CWE-1322", detectCWE1322, &MetaCWE1322)
	RegisterRule("CWE-1339", detectCWE1339, &MetaCWE1339)
	RegisterRule("CWE-1341", detectCWE1341, &MetaCWE1341)
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
	pyTierBDoubleCloseRE = regexp.MustCompile(`(?is)\w+\.close\s*\(\s*\)[\s\S]{0,180}\w+\.close\s*\(\s*\)`)
)

func detectCWE1104(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, regexp.MustCompile(`(?im)^\s*(?:import|from)\s+imp\b`)); start >= 0 {
		emitTierBFinding(unit, &MetaCWE1104, start, "unmaintained imp component is imported", 0.74, out)
	}
}

func detectCWE1106(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBMagicCookieRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE1106, start, "security-sensitive cookie lifetime uses an unexplained magic number", 0.68, out)
	}
}

func detectCWE1108(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	for _, fn := range pythonFunctions(unit.Source) {
		if strings.Count(pythonCodeMask(fn.body), "global ") >= 3 {
			emitTierBFinding(unit, &MetaCWE1108, fn.bodyStart, "function relies on multiple global variables", 0.7, out)
			return
		}
	}
}

func detectCWE1121(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	for _, fn := range pythonFunctions(unit.Source) {
		code := pythonCodeMask(fn.body)
		branches := strings.Count(code, "if ") + strings.Count(code, "elif ") + strings.Count(code, "for ") + strings.Count(code, "while ") + strings.Count(code, "except ")
		if branches >= 12 {
			emitTierBFinding(unit, &MetaCWE1121, fn.bodyStart, "function has at least twelve visible control-flow branches", 0.7, out)
			return
		}
	}
}

func detectCWE1123(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	for _, call := range findCalls(unit.Source, "types.FunctionType") {
		emitTierBFinding(unit, &MetaCWE1123, call.Start, "function object is constructed dynamically from code", 0.8, out)
		return
	}
}

func detectCWE1124(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	for offset := 0; offset < len(unit.Source); {
		next := strings.IndexByte(unit.Source[offset:], '\n')
		end := len(unit.Source)
		if next >= 0 {
			end = offset + next
		}
		text := pythonCodeMask(unit.Source[offset:end])
		if strings.TrimSpace(text) != "" && len(text)-len(strings.TrimLeft(text, " ")) >= 24 {
			emitTierBFinding(unit, &MetaCWE1124, offset, "executable statement is nested at least six indentation levels", 0.7, out)
			return
		}
		if end == len(unit.Source) {
			break
		}
		offset = end + 1
	}
}

func detectCWE1220(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBUserLookupRE); start >= 0 && !strings.Contains(unit.Source[start:start+minInt(400, len(unit.Source)-start)], "request.user") {
		emitTierBFinding(unit, &MetaCWE1220, start, "object is fetched by request identifier without an owner constraint", 0.76, out)
	}
}

func detectCWE1265(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBReentrantRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE1265, start, "non-reentrant lock is acquired again in a nested scope", 0.8, out)
	}
}

func detectCWE1284(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBLimitRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE1284, start, "request quantity controls a data-access limit without a visible bound", 0.76, out)
	}
}

func detectCWE1285(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBIndexRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE1285, start, "request value is converted directly into a collection index", 0.78, out)
	}
}

func detectCWE1287(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBTypeRE); start >= 0 && !strings.Contains(unit.Source, "isinstance(") {
		emitTierBFinding(unit, &MetaCWE1287, start, "JSON request object is used for execution without a visible type check", 0.7, out)
	}
}

func detectCWE1288(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBConsistencyRE); start >= 0 && !strings.Contains(unit.Source, "len(items)") {
		emitTierBFinding(unit, &MetaCWE1288, start, "request count and item collection are used without a consistency check", 0.7, out)
	}
}

func detectCWE1322(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBAsyncSleepRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE1322, start, "blocking sleep runs inside an async function", 0.84, out)
	}
}

func detectCWE1339(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBFloatMoneyRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE1339, start, "request monetary value is parsed as a binary floating-point number", 0.8, out)
	}
}

func detectCWE1341(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBDoubleCloseRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE1341, start, "same resource handle is released twice", 0.82, out)
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
