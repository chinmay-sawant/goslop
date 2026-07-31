package cwe

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	RegisterRule("CWE-908", detectCWE908, &MetaCWE908)
	RegisterRule("CWE-909", detectCWE909, &MetaCWE909)
	RegisterRule("CWE-910", detectCWE910, &MetaCWE910)
	RegisterRule("CWE-911", detectCWE911, &MetaCWE911)
	RegisterRule("CWE-920", detectCWE920, &MetaCWE920)
	RegisterRule("CWE-939", detectCWE939, &MetaCWE939)
	RegisterRule("CWE-1007", detectCWE1007, &MetaCWE1007)
	RegisterRule("CWE-1021", detectCWE1021, &MetaCWE1021)
	RegisterRule("CWE-1046", detectCWE1046, &MetaCWE1046)
	RegisterRule("CWE-1050", detectCWE1050, &MetaCWE1050)
	RegisterRule("CWE-1060", detectCWE1060, &MetaCWE1060)
	RegisterRule("CWE-1067", detectCWE1067, &MetaCWE1067)
	RegisterRule("CWE-1071", detectCWE1071, &MetaCWE1071)
	RegisterRule("CWE-1072", detectCWE1072, &MetaCWE1072)
	RegisterRule("CWE-1084", detectCWE1084, &MetaCWE1084)
}

var (
	pyTierBNoneUseRE     = regexp.MustCompile(`(?is)\w+\s*=\s*None[\s\S]{0,220}\w+\.(?:read|write|execute|connect)\s*\(`)
	pyTierBClosedUseRE   = regexp.MustCompile(`(?is)\w+\.close\s*\(\s*\)[\s\S]{0,180}\w+\.(?:read|write|flush)\s*\(`)
	pyTierBBusyLoopRE    = regexp.MustCompile(`(?is)while\s+True\s*:[\s\S]{0,260}(?:hashlib|sha256|calculate|compute)`)
	pyTierBHomoglyphRE   = regexp.MustCompile(`(?is)render_template\s*\([^\n]*request\.(?:args|form)\.get\s*\([^\n]*username`)
	pyTierBFrameRE       = regexp.MustCompile(`(?is)response\s*=\s*make_response\s*\([^\n]*render_template[\s\S]{0,240}return\s+response`)
	pyTierBConcatLoopRE  = regexp.MustCompile(`(?is)(?:for|while)\s+[^\n]+:\s*\n(?:\s+[^\n]*\n){0,4}\s*\w+\s*\+=\s*`)
	pyTierBOpenLoopRE    = regexp.MustCompile(`(?is)(?:for|while)\s+[^\n]+:\s*\n(?:\s+[^\n]*\n){0,4}\s*open\s*\(`)
	pyTierBNPlusOneRE    = regexp.MustCompile(`(?is)for\s+\w+\s+in\s+\w+\.objects\.(?:all|filter)\s*\([^\n]*\)\s*:[\s\S]{0,180}\w+\.[a-z_]+\.(?:all|filter)\s*\(`)
	pyTierBEmptyExceptRE = regexp.MustCompile(`(?is)except(?:\s+[A-Za-z_][A-Za-z0-9_.]*)?\s*:\s*\n\s*pass\b`)
	pyTierBPoolRouteRE   = regexp.MustCompile(`(?is)@\w+\.route\s*\([^\n]*\)[\s\S]{0,500}psycopg2\.connect\s*\(`)
)

func detectCWE908(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBNoneUseRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE908, start, "resource initialized to None is used without initialization", 0.78, out)
	}
}

func detectCWE909(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	for _, fn := range pythonFunctions(unit.Source) {
		code := pythonCodeMask(fn.body)
		if strings.Contains(code, "db.execute(") && !strings.Contains(code, "db =") && !strings.Contains(code, "get_db(") {
			emitTierBFinding(unit, &MetaCWE909, fn.bodyStart, "database resource is used without local initialization evidence", 0.68, out)
			return
		}
	}
}

func detectCWE910(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBClosedUseRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE910, start, "closed file descriptor is used again", 0.86, out)
	}
}

func detectCWE911(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	for _, call := range findCalls(unit.Source, "ctypes.pythonapi.Py_IncRef", "ctypes.pythonapi.Py_DecRef") {
		emitTierBFinding(unit, &MetaCWE911, call.Start, "manual CPython reference-count API is invoked", 0.8, out)
		return
	}
}

func detectCWE920(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBBusyLoopRE); start >= 0 && !strings.Contains(unit.Source[start:], "time.sleep(") {
		emitTierBFinding(unit, &MetaCWE920, start, "unbounded busy loop repeatedly performs expensive computation", 0.76, out)
	}
}

func detectCWE939(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	for _, call := range findCalls(unit.Source, "webbrowser.open") {
		if strings.Contains(call.ArgsText, "request.") {
			emitTierBFinding(unit, &MetaCWE939, call.Start, "custom URL handler opens a request-controlled target", 0.78, out)
			return
		}
	}
}

func detectCWE1007(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBHomoglyphRE); start >= 0 && !strings.Contains(unit.Source, "normalize(") {
		emitTierBFinding(unit, &MetaCWE1007, start, "request username is rendered without Unicode normalization", 0.7, out)
	}
}

func detectCWE1021(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBFrameRE); start >= 0 && !strings.Contains(unit.Source, "X-Frame-Options") && !strings.Contains(unit.Source, "frame-ancestors") {
		emitTierBFinding(unit, &MetaCWE1021, start, "HTML response is returned without an observable frame-embedding restriction", 0.68, out)
	}
}

func detectCWE1046(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBConcatLoopRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE1046, start, "immutable text is repeatedly concatenated inside a loop", 0.76, out)
	}
}

func detectCWE1050(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBOpenLoopRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE1050, start, "platform resource is opened for every loop iteration", 0.76, out)
	}
}

func detectCWE1060(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBNPlusOneRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE1060, start, "ORM relation is loaded once per query result", 0.76, out)
	}
}

func detectCWE1067(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	for _, call := range findCalls(unit.Source, ".filter") {
		if strings.Contains(call.ArgsText, "__contains=") || strings.Contains(call.ArgsText, "__icontains=") {
			emitTierBFinding(unit, &MetaCWE1067, call.Start, "data-resource search uses an unanchored contains lookup", 0.72, out)
			return
		}
	}
}

func detectCWE1071(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBEmptyExceptRE); start >= 0 {
		emitTierBFinding(unit, &MetaCWE1071, start, "exception handler silently contains only pass", 0.78, out)
	}
}

func detectCWE1072(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstCodeMatchStart(unit.Source, pyTierBPoolRouteRE); start >= 0 && !strings.Contains(unit.Source, "ThreadedConnectionPool") {
		emitTierBFinding(unit, &MetaCWE1072, start, "route opens a direct database connection without pool evidence", 0.72, out)
	}
}

func detectCWE1084(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	for _, fn := range pythonFunctions(unit.Source) {
		if len(findCalls(fn.body, "open", ".execute")) >= 3 {
			emitTierBFinding(unit, &MetaCWE1084, fn.bodyStart, "single function performs many file or data-access operations", 0.7, out)
			return
		}
	}
}
