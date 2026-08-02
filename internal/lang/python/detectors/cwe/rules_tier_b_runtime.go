package cwe

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	RegisterRule("CWE-908", detectCWE908, &MetaCWE908,
		"= None", "=None")
	RegisterRule("CWE-909", detectCWE909, &MetaCWE909,
		"db.execute(")
	RegisterRule("CWE-910", detectCWE910, &MetaCWE910,
		".close()")
	RegisterRule("CWE-911", detectCWE911, &MetaCWE911,
		"ctypes.pythonapi", "Py_IncRef", "Py_DecRef")
	RegisterRule("CWE-920", detectCWE920, &MetaCWE920,
		"while True", "hashlib", "sha256", "calculate", "compute")
	RegisterRule("CWE-939", detectCWE939, &MetaCWE939,
		"webbrowser.open")
	RegisterRule("CWE-1007", detectCWE1007, &MetaCWE1007,
		"render_template", "request.args", "request.form")
	RegisterRule("CWE-1021", detectCWE1021, &MetaCWE1021,
		"make_response", "render_template")
	RegisterRule("CWE-1046", detectCWE1046, &MetaCWE1046,
		"+=")
	RegisterRule("CWE-1050", detectCWE1050, &MetaCWE1050,
		"open(", "for ", "while ")
	RegisterRule("CWE-1060", detectCWE1060, &MetaCWE1060,
		".objects.all", ".objects.filter")
	RegisterRule("CWE-1067", detectCWE1067, &MetaCWE1067,
		"__contains=", "__icontains=", ".filter")
	RegisterRule("CWE-1071", detectCWE1071, &MetaCWE1071,
		"except")
	RegisterRule("CWE-1072", detectCWE1072, &MetaCWE1072,
		"psycopg2.connect", ".route")
	RegisterRule("CWE-1084", detectCWE1084, &MetaCWE1084,
		"open(", ".execute(")
}

var (
	pyTierBNoneAssignRE  = regexp.MustCompile(`(?m)\b([A-Za-z_][A-Za-z0-9_]*)\s*=\s*None\b`)
	pyTierBCloseCallRE   = regexp.MustCompile(`\b([A-Za-z_][A-Za-z0-9_]*)\.close\s*\(\s*\)`)
	pyTierBBusyLoopRE    = regexp.MustCompile(`(?is)while\s+True\s*:[\s\S]{0,260}(?:hashlib|sha256|calculate|compute)`)
	pyTierBHomoglyphRE   = regexp.MustCompile(`(?is)render_template\s*\([^\n]*request\.(?:args|form)\.get\s*\([^\n]*username`)
	pyTierBFrameRE       = regexp.MustCompile(`(?is)response\s*=\s*make_response\s*\([^\n]*render_template[\s\S]{0,240}return\s+response`)
	pyTierBLoopHeaderRE  = regexp.MustCompile(`^(?:for|while)\b.*:\s*$`)
	pyTierBAugAssignRE   = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)\s*\+=\s*(.+)$`)
	pyTierBOpenLoopRE    = regexp.MustCompile(`(?is)(?:for|while)\s+[^\n]+:\s*\n(?:\s+[^\n]*\n){0,4}\s*open\s*\(`)
	pyTierBNPlusOneRE    = regexp.MustCompile(`(?is)for\s+\w+\s+in\s+\w+\.objects\.(?:all|filter)\s*\([^\n]*\)\s*:[\s\S]{0,180}\w+\.[a-z_]+\.(?:all|filter)\s*\(`)
	pyTierBEmptyExceptRE = regexp.MustCompile(`(?is)except(?:\s+[A-Za-z_][A-Za-z0-9_.]*)?\s*:\s*\n\s*pass\b`)
	pyTierBPoolRouteRE   = regexp.MustCompile(`(?is)@\w+\.route\s*\([^\n]*\)[\s\S]{0,500}psycopg2\.connect\s*\(`)
)

func detectCWE908(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	code := facts.Masked
	if !strings.Contains(code, "= None") && !strings.Contains(code, "=None") {
		return
	}
	search := 0
	for search <= len(code) {
		match := pyTierBNoneAssignRE.FindStringSubmatchIndex(code[search:])
		if match == nil {
			return
		}
		abs0 := search + match[0]
		abs1 := search + match[1]
		name := code[search+match[2] : search+match[3]]
		if resourceUseStart(unit.Source, code, name, abs1, uninitializedResourceWindow, "read", "write", "execute", "connect") >= 0 {
			emitTierBFinding(unit, &MetaCWE908, abs0, "resource initialized to None is used without initialization", confidence78, out)
			return
		}
		if abs1 <= search {
			search++
			continue
		}
		search = abs1
	}
}

func detectCWE909(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if !strings.Contains(facts.Masked, "db.execute(") {
		return
	}
	for _, fn := range facts.Functions() {
		code := facts.codeMask(fn.body, fn.bodyStart)
		if strings.Contains(code, "db.execute(") && !strings.Contains(code, "db =") && !strings.Contains(code, "get_db(") {
			emitTierBFinding(unit, &MetaCWE909, fn.bodyStart, "database resource is used without local initialization evidence", confidence68, out)
			return
		}
	}
}

func detectCWE910(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	code := facts.Masked
	if !strings.Contains(code, ".close(") {
		return
	}
	search := 0
	for search <= len(code) {
		match := pyTierBCloseCallRE.FindStringSubmatchIndex(code[search:])
		if match == nil {
			return
		}
		abs0 := search + match[0]
		abs1 := search + match[1]
		name := code[search+match[2] : search+match[3]]
		if resourceUseStart(unit.Source, code, name, abs1, closedResourceWindow, "read", "write", "flush") >= 0 {
			emitTierBFinding(unit, &MetaCWE910, abs0, "closed file descriptor is used again", confidence86, out)
			return
		}
		if abs1 <= search {
			search++
			continue
		}
		search = abs1
	}
}

func resourceUseStart(source, code, name string, from, span int, methods ...string) int {
	end := min(from+span, len(code))
	assignment := regexp.MustCompile(`(?m)\b` + regexp.QuoteMeta(name) + `\s*=`)
	if next := assignment.FindStringIndex(code[from:end]); next != nil {
		end = from + next[0]
	}
	qualified := make([]string, 0, len(methods))
	for _, method := range methods {
		qualified = append(qualified, name+"."+method)
	}
	calls := findCallsMasked(source[from:end], code[from:end], qualified...)
	if len(calls) == 0 {
		return -1
	}
	return from + calls[0].Start
}

func detectCWE911(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	for _, call := range findCalls(facts, unit.Source, "ctypes.pythonapi.Py_IncRef", "ctypes.pythonapi.Py_DecRef") {
		emitTierBFinding(unit, &MetaCWE911, call.Start, "manual CPython reference-count API is invoked", confidence80, out)
		return
	}
}

func detectCWE920(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(facts, unit.Source, pyTierBBusyLoopRE); start >= 0 && !strings.Contains(unit.Source[start:], "time.sleep(") {
		emitTierBFinding(unit, &MetaCWE920, start, "unbounded busy loop repeatedly performs expensive computation", confidence76, out)
	}
}

func detectCWE939(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	for _, call := range findCalls(facts, unit.Source, "webbrowser.open") {
		if strings.Contains(call.ArgsText, "request.") {
			emitTierBFinding(unit, &MetaCWE939, call.Start, "custom URL handler opens a request-controlled target", confidence78, out)
			return
		}
	}
}

func detectCWE1007(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstLiteralMatchStartIfContains(facts, unit, pyTierBHomoglyphRE,
		"username", "render_template"); start >= 0 && !strings.Contains(unit.Source, "normalize(") {
		emitTierBFinding(unit, &MetaCWE1007, start, "request username is rendered without Unicode normalization", confidence70, out)
	}
}

func detectCWE1021(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(facts, unit.Source, pyTierBFrameRE); start >= 0 && !strings.Contains(unit.Source, "X-Frame-Options") && !strings.Contains(unit.Source, "frame-ancestors") {
		emitTierBFinding(unit, &MetaCWE1021, start, "HTML response is returned without an observable frame-embedding restriction", confidence68, out)
	}
}

func detectCWE1046(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if !strings.Contains(facts.Masked, "+=") {
		return
	}
	maskedLines := strings.Split(facts.Masked, "\n")
	originalLines := strings.Split(unit.Source, "\n")
	loopIndents := make([]int, 0, 4)
	offset := 0
	for i, masked := range maskedLines {
		trimmed := strings.TrimSpace(masked)
		if trimmed == "" {
			offset += len(masked) + 1
			continue
		}
		indent := len(masked) - len(strings.TrimLeft(masked, " \t"))
		for len(loopIndents) > 0 && indent <= loopIndents[len(loopIndents)-1] {
			loopIndents = loopIndents[:len(loopIndents)-1]
		}
		if pyTierBLoopHeaderRE.MatchString(trimmed) {
			loopIndents = append(loopIndents, indent)
			offset += len(masked) + 1
			continue
		}
		if len(loopIndents) == 0 || !strings.Contains(trimmed, "+=") {
			offset += len(masked) + 1
			continue
		}
		match := pyTierBAugAssignRE.FindStringSubmatch(trimmed)
		if len(match) == 3 && textAccumulatorEvidence(originalLines[:i], match[1], match[2]) {
			emitTierBFinding(unit, &MetaCWE1046, offset, "immutable text is repeatedly concatenated inside a loop", confidence76, out)
			return
		}
		offset += len(masked) + 1
	}
}

func textAccumulatorEvidence(previous []string, name, rhs string) bool {
	for i := len(previous) - 1; i >= 0; i-- {
		line := strings.TrimSpace(previous[i])
		if !strings.HasPrefix(line, name+" = ") {
			continue
		}
		value := strings.TrimSpace(strings.TrimPrefix(line, name+" = "))
		// bytearray += is in-place mutable buffer extend, not immutable text concat.
		if strings.HasPrefix(value, "bytearray(") {
			return false
		}
		if strings.HasPrefix(value, "\"") || strings.HasPrefix(value, "'") || strings.HasPrefix(value, "str(") {
			return true
		}
		break
	}
	lowerName := strings.ToLower(name)
	if strings.Contains(lowerName, "text") || strings.Contains(lowerName, "string") || strings.Contains(lowerName, "message") || strings.Contains(lowerName, "output") {
		return true
	}
	trimmedRHS := strings.TrimSpace(rhs)
	return strings.HasPrefix(trimmedRHS, "f\"") || strings.HasPrefix(trimmedRHS, "f'") ||
		strings.HasPrefix(trimmedRHS, "\"") || strings.HasPrefix(trimmedRHS, "'")
}

func detectCWE1050(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(facts, unit.Source, pyTierBOpenLoopRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE1050, start, "platform resource is opened for every loop iteration", confidence76, out)
	}
}

func detectCWE1060(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	code := facts.Masked
	if !strings.Contains(code, ".objects.") {
		return
	}
	if start := firstCodeMatchStart(facts, unit.Source, pyTierBNPlusOneRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE1060, start, "ORM relation is loaded once per query result", confidence76, out)
	}
}

func detectCWE1067(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	code := facts.Masked
	if !strings.Contains(code, ".filter") {
		return
	}
	if !strings.Contains(code, "__contains=") && !strings.Contains(code, "__icontains=") {
		return
	}
	for _, call := range findCalls(facts, unit.Source, ".filter") {
		if strings.Contains(call.ArgsText, "__contains=") || strings.Contains(call.ArgsText, "__icontains=") {
			emitTierBFinding(unit, &MetaCWE1067, call.Start, "data-resource search uses an unanchored contains lookup", confidence72, out)
			return
		}
	}
}

func detectCWE1071(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	// Expected-exception FPs use exceptPassIsSafe; do not path-skip modules.

	lines := facts.MaskedLines()
	rawLines := buildMaskedPythonLines(unit.Source)
	for i, line := range lines {
		if !pyExceptStartRE.MatchString(line.text) {
			continue
		}
		if !exceptPassOnly(lines, i) {
			continue
		}
		if exceptPassIsSafe(unit, lines, rawLines, i) {
			continue
		}
		emitTierBFinding(unit, &MetaCWE1071, line.start, "exception handler silently contains only pass", confidence78, out)
		return
	}
}

func detectCWE1072(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(facts, unit.Source, pyTierBPoolRouteRE); start >= 0 && !strings.Contains(unit.Source, "ThreadedConnectionPool") {
		emitTierBFinding(unit, &MetaCWE1072, start, "route opens a direct database connection without pool evidence", confidence72, out)
	}
}

func detectCWE1084(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	for _, fn := range facts.Functions() {
		masked := facts.codeMask(fn.body, fn.bodyStart)
		if len(findCallsMasked(fn.body, masked, "open", ".execute")) >= 3 {
			emitTierBFinding(unit, &MetaCWE1084, fn.bodyStart, "single function performs many file or data-access operations", confidence70, out)
			return
		}
	}
}
