package cwe

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	RegisterRule("CWE-749", detectCWE749, &MetaCWE749,
		"eval(", "exec(", "compile(", "__import__(", "importlib.import_module", "os.system")
	RegisterRule("CWE-829", detectCWE829, &MetaCWE829,
		"__import__(", "importlib.import_module", "runpy.run_path", "spec_from_file_location")
	RegisterRule("CWE-695", detectCWE695, &MetaCWE695,
		"ctypes.", "cffi.FFI", "mmap.mmap")
	RegisterRule("CWE-214", detectCWE214, &MetaCWE214,
		"subprocess.")
	RegisterRule("CWE-215", detectCWE215, &MetaCWE215,
		"print(", "logging.debug", ".debug(")
}

var (
	pythonRouteDecoratorRE = regexp.MustCompile(`(?m)^\s*@(?:[A-Za-z_][\w]*\.)?(?:route|api_route|get|post|put|delete|patch)\s*\(`)
	sensitiveValueRE       = regexp.MustCompile(`(?i)\b(?:password|passwd|token|secret|api[_-]?key|credential)s?\b`)
	sensitiveCommandRE     = regexp.MustCompile(`(?i)(?:--(?:password|passwd|token|secret|api[_-]?key|credential)(?:=|\b)|\b(?:password|passwd|token|secret|api[_-]?key|credential)\s*=)`)
	sensitiveFileOptionRE  = regexp.MustCompile(`(?i)--(?:password|passwd|token|secret|api[_-]?key|credential)s?-file\b`)
)

// CWE-749 reports dynamic execution only when it is in a route-decorated
// handler. This same-file boundary deliberately avoids treating internal admin
// scripts or sandboxed helpers as an exposed API.
func detectCWE749(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, handler := range routeHandlerBodies(facts, unit.Source) {
		masked := facts.codeMask(handler.body, handler.start)
		for _, name := range []string{"eval", "exec", "compile", "__import__", "os.system"} {
			for _, call := range findCallsMasked(handler.body, masked, name) {
				if len(splitTopLevelArgs(call.ArgsText)) == 0 || !isDynamicExpr(call.ArgsText) {
					continue
				}
				emitCodeDynamic(unit, &MetaCWE749, handler.start+call.Start,
					"route handler exposes dynamic code execution to external callers", confidence84, out)
				return
			}
		}
		for _, call := range findCallsMasked(handler.body, masked, "importlib.import_module") {
			args := splitTopLevelArgs(call.ArgsText)
			if len(args) > 0 && isDynamicExpr(args[0]) {
				emitCodeDynamic(unit, &MetaCWE749, handler.start+call.Start,
					"route handler exposes a dynamically selected module import", confidence80, out)
				return
			}
		}
	}
}

// CWE-829 reports only imports and source-file execution selected by a dynamic
// expression; package-controlled literal module names and paths are suppressed.
func detectCWE829(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"__import__", "importlib.import_module", "runpy.run_path"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			args := splitTopLevelArgs(call.ArgsText)
			if len(args) > 0 && isDynamicExpr(args[0]) {
				emitCodeDynamic(unit, &MetaCWE829, call.Start,
					"dynamically selected module or code path is loaded for execution", confidence82, out)
				return
			}
		}
	}
	for _, call := range findCalls(facts, unit.Source, "importlib.util.spec_from_file_location") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) >= 2 && isDynamicExpr(args[1]) {
			emitCodeDynamic(unit, &MetaCWE829, call.Start,
				"dynamically selected file path is prepared for module execution", confidence80, out)
			return
		}
	}
}

// CWE-695 is intentionally narrow: importing a module alone is common, while
// these calls cross into native code or raw virtual-memory functionality.
func detectCWE695(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"ctypes.CDLL", "ctypes.PyDLL", "ctypes.cast", "cffi.FFI", "mmap.mmap"} {
		if calls := findCalls(facts, unit.Source, name); len(calls) > 0 {
			emitCodeDynamic(unit, &MetaCWE695, calls[0].Start,
				"low-level native or memory interface bypasses higher-level safety controls", confidence76, out)
			return
		}
	}
}

// CWE-214 recognizes dynamic secret-bearing process arguments or subprocess
// environment maps. Literal examples and password-file options remain safe.
func detectCWE214(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"subprocess.run", "subprocess.call", "subprocess.check_call", "subprocess.check_output", "subprocess.Popen"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			visibleArgs := sensitiveFileOptionRE.ReplaceAllString(call.ArgsText, "")
			if !sensitiveCommandRE.MatchString(visibleArgs) && !strings.Contains(strings.ToLower(visibleArgs), "env=") {
				continue
			}
			if sensitiveValueRE.MatchString(visibleArgs) && isDynamicExpr(visibleArgs) {
				emitCodeDynamic(unit, &MetaCWE214, call.Start,
					"subprocess receives a sensitive value through visible arguments or environment", confidence82, out)
				return
			}
		}
	}
}

// CWE-215 reports debug sinks only when their arguments include a sensitive
// identifier. Generic debug logs and intentionally redacted literals are safe.
func detectCWE215(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"print", "logging.debug", ".debug"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			if sensitiveValueRE.MatchString(call.ArgsText) && !isPureStringLiteral(call.ArgsText) {
				emitCodeDynamic(unit, &MetaCWE215, call.Start,
					"debug output includes a sensitive value; redact it before logging", confidence86, out)
				return
			}
		}
	}
}

func routeHandlerBodies(facts *PyCweFacts, src string) []routeHandlerBody {
	if facts != nil && src == facts.Source {
		return facts.RouteHandlers()
	}
	var code string
	if facts != nil {
		code = facts.codeMask(src, fragStartHint(facts, src))
	} else {
		code = pythonCodeMask(src)
	}
	return buildRouteHandlerBodies(src, code)
}

// buildRouteHandlerBodies finds @route/@get/… handlers via iterative FindStringIndex
// on masked text (never FindAllStringIndex).
func buildRouteHandlerBodies(src, code string) []routeHandlerBody {
	if code == "" {
		code = pythonCodeMask(src)
	}
	var bodies []routeHandlerBody
	search := 0
	for search <= len(code) {
		loc := pythonRouteDecoratorRE.FindStringIndex(code[search:])
		if loc == nil {
			break
		}
		decStart := search + loc[0]
		decEnd := search + loc[1]
		defAt := strings.Index(code[decEnd:], "\ndef ")
		if defAt < 0 {
			search = decEnd
			continue
		}
		start := decEnd + defAt + 1
		end := len(src)
		if next := strings.Index(code[start+4:], "\ndef "); next >= 0 {
			end = start + 4 + next
		}
		if end > len(src) {
			end = len(src)
		}
		bodies = append(bodies, routeHandlerBody{
			decoratorStart: decStart,
			start:          start,
			body:           src[start:end],
		})
		search = decEnd
	}
	return bodies
}

func emitCodeDynamic(unit *core.ParsedUnit, meta *rules.RuleMetadata, offset int, message string, confidence float32, out *[]rules.Finding) {
	line, col := unit.LineCol(offset)
	rules.PushFindingWithConfidence(meta, unitFile(unit), line, col, message, confidence, out)
}
