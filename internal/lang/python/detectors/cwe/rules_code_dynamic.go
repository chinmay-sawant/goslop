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
		"__import__(", "importlib.import_module", "spec_from_file_location", "runpy.run_path")
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
)

// CWE-749 reports dynamic execution only when it is in a route-decorated
// handler. This same-file boundary deliberately avoids treating internal admin
// scripts or sandboxed helpers as an exposed API.
func detectCWE749(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, handler := range routeHandlerBodies(unit.Source) {
		for _, name := range []string{"eval", "exec", "compile", "__import__", "os.system"} {
			for _, call := range findCalls(handler.body, name) {
				if len(splitTopLevelArgs(call.ArgsText)) == 0 || !isDynamicExpr(call.ArgsText) {
					continue
				}
				emitCodeDynamic(unit, &MetaCWE749, handler.start+call.Start,
					"route handler exposes dynamic code execution to external callers", 0.84, out)
				return
			}
		}
		for _, call := range findCalls(handler.body, "importlib.import_module") {
			args := splitTopLevelArgs(call.ArgsText)
			if len(args) > 0 && isDynamicExpr(args[0]) {
				emitCodeDynamic(unit, &MetaCWE749, handler.start+call.Start,
					"route handler exposes a dynamically selected module import", 0.8, out)
				return
			}
		}
	}
}

// CWE-829 reports only imports and source-file execution selected by a dynamic
// expression; package-controlled literal module names and paths are suppressed.
func detectCWE829(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"__import__", "importlib.import_module", "runpy.run_path"} {
		for _, call := range findCalls(unit.Source, name) {
			args := splitTopLevelArgs(call.ArgsText)
			if len(args) > 0 && isDynamicExpr(args[0]) {
				emitCodeDynamic(unit, &MetaCWE829, call.Start,
					"dynamically selected module or code path is loaded for execution", 0.82, out)
				return
			}
		}
	}
	for _, call := range findCalls(unit.Source, "importlib.util.spec_from_file_location") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) >= 2 && isDynamicExpr(args[1]) {
			emitCodeDynamic(unit, &MetaCWE829, call.Start,
				"dynamically selected file path is prepared for module execution", 0.8, out)
			return
		}
	}
}

// CWE-695 is intentionally narrow: importing a module alone is common, while
// these calls cross into native code or raw virtual-memory functionality.
func detectCWE695(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"ctypes.CDLL", "ctypes.PyDLL", "ctypes.cast", "cffi.FFI", "mmap.mmap"} {
		if calls := findCalls(unit.Source, name); len(calls) > 0 {
			emitCodeDynamic(unit, &MetaCWE695, calls[0].Start,
				"low-level native or memory interface bypasses higher-level safety controls", 0.76, out)
			return
		}
	}
}

// CWE-214 recognizes dynamic secret-bearing process arguments or subprocess
// environment maps. Literal examples and password-file options remain safe.
func detectCWE214(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"subprocess.run", "subprocess.call", "subprocess.check_call", "subprocess.check_output", "subprocess.Popen"} {
		for _, call := range findCalls(unit.Source, name) {
			if !sensitiveCommandRE.MatchString(call.ArgsText) && !strings.Contains(strings.ToLower(call.ArgsText), "env=") {
				continue
			}
			if sensitiveValueRE.MatchString(call.ArgsText) && isDynamicExpr(call.ArgsText) {
				emitCodeDynamic(unit, &MetaCWE214, call.Start,
					"subprocess receives a sensitive value through visible arguments or environment", 0.82, out)
				return
			}
		}
	}
}

// CWE-215 reports debug sinks only when their arguments include a sensitive
// identifier. Generic debug logs and intentionally redacted literals are safe.
func detectCWE215(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"print", "logging.debug", ".debug"} {
		for _, call := range findCalls(unit.Source, name) {
			if sensitiveValueRE.MatchString(call.ArgsText) && !isPureStringLiteral(call.ArgsText) {
				emitCodeDynamic(unit, &MetaCWE215, call.Start,
					"debug output includes a sensitive value; redact it before logging", 0.86, out)
				return
			}
		}
	}
}

type routeHandlerBody struct {
	start int
	body  string
}

func routeHandlerBodies(src string) []routeHandlerBody {
	locs := pythonRouteDecoratorRE.FindAllStringIndex(src, -1)
	bodies := make([]routeHandlerBody, 0, len(locs))
	for _, loc := range locs {
		defAt := strings.Index(src[loc[1]:], "\ndef ")
		if defAt < 0 {
			continue
		}
		start := loc[1] + defAt + 1
		end := len(src)
		if next := strings.Index(src[start+4:], "\ndef "); next >= 0 {
			end = start + 4 + next
		}
		bodies = append(bodies, routeHandlerBody{start: start, body: src[start:end]})
	}
	return bodies
}

func emitCodeDynamic(unit *core.ParsedUnit, meta *rules.RuleMetadata, offset int, message string, confidence float32, out *[]rules.Finding) {
	line, col := unit.LineCol(offset)
	rules.PushFindingWithConfidence(meta, unitFile(unit), line, col, message, confidence, out)
}
