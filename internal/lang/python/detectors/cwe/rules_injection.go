package cwe

import (
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	RegisterRule("CWE-90", detectCWE90, &MetaCWE90, "ldap3", "ldap.initialize", ".search(", ".search_s(")
	RegisterRule("CWE-91", detectCWE91, &MetaCWE91, ".xpath(", "XPath(", ".fromstring(")
	RegisterRule("CWE-93", detectCWE93, &MetaCWE93, "response.headers[", ".set_header(", ".add_header(", "HttpResponseRedirect(")
	RegisterRule("CWE-94", detectCWE94, &MetaCWE94, "eval(", "exec(", "compile(", "__import__(", "importlib.import_module")
	RegisterRule("CWE-88", detectCWE88, &MetaCWE88, "subprocess.")
	RegisterRule("CWE-117", detectCWE117, &MetaCWE117, "logging.", "logger.", "log.")
}

// CWE-90: only dynamic LDAP filter expressions reach LDAP search APIs. Filter
// values escaped through the standard ldap3 helper are intentionally suppressed.
func detectCWE90(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(unit.Source, ".search", ".search_s") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) == 0 || !isDynamicExpr(args[0]) || ldapFilterLooksEscaped(args[0]) {
			continue
		}
		pushInjectionFinding(unit, call.Start, &MetaCWE90,
			"dynamic value reaches an LDAP search filter without LDAP filter escaping", 0.78, out)
		return
	}
}

func ldapFilterLooksEscaped(expr string) bool {
	return strings.Contains(expr, "escape_filter_chars(") || strings.Contains(expr, "escape_filter_value(")
}

// CWE-91: flag dynamic XML/XPath expression construction only at parser and
// XPath APIs. Literal XPath with bound variables remains a supported safe form.
func detectCWE91(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(unit.Source, ".xpath", "etree.XPath", "ElementTree.XPath", "ET.XPath", "etree.fromstring", "ElementTree.fromstring", "ET.fromstring") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) == 0 || !isDynamicExpr(args[0]) {
			continue
		}
		pushInjectionFinding(unit, call.Start, &MetaCWE91,
			"dynamic value is formatted into an XML or XPath expression", 0.75, out)
		return
	}
}

// CWE-93: header sinks are limited to explicit header mutation APIs. Values
// that visibly remove both CR and LF are not reported by this source heuristic.
func detectCWE93(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(unit.Source, ".set_header", ".add_header", "HttpResponseRedirect") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) < 2 || !isDynamicExpr(args[len(args)-1]) || headerValueLooksSanitized(args[len(args)-1]) {
			continue
		}
		pushInjectionFinding(unit, call.Start, &MetaCWE93,
			"dynamic value is written to an HTTP response header without CRLF neutralization", 0.72, out)
		return
	}
	for _, line := range strings.Split(unit.Source, "\n") {
		lhs, rhs, ok := strings.Cut(line, "=")
		if !ok || !strings.Contains(lhs, ".headers[") || !isDynamicExpr(rhs) || headerValueLooksSanitized(rhs) {
			continue
		}
		start := strings.Index(unit.Source, line) + strings.Index(line, ".headers[")
		pushInjectionFinding(unit, start, &MetaCWE93,
			"dynamic value is written to an HTTP response header without CRLF neutralization", 0.72, out)
		return
	}
}

func headerValueLooksSanitized(expr string) bool {
	compact := compactWhitespace(expr)
	return (strings.Contains(compact, `replace("\r","")`) || strings.Contains(compact, `replace('\r','')`)) &&
		(strings.Contains(compact, `replace("\n","")`) || strings.Contains(compact, `replace('\n','')`))
}

// CWE-94: code-generation and dynamic-import APIs are only findings when the
// code/module argument is non-literal. Literal developer-owned expressions are
// out of scope for this same-file heuristic.
func detectCWE94(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(unit.Source, "eval", "exec", "compile", "__import__", "importlib.import_module") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) == 0 || !isDynamicExpr(args[0]) {
			continue
		}
		pushInjectionFinding(unit, call.Start, &MetaCWE94,
			"dynamic value reaches a Python code-generation or dynamic-import sink", 0.82, out)
		return
	}
}

// CWE-88: shell=False is not sufficient when an untrusted argument can become
// a tool option. Detect only explicit argv literals with a dynamic segment; a
// pre-built argv variable has insufficient same-file evidence to report.
func detectCWE88(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(unit.Source,
		"subprocess.run", "subprocess.call", "subprocess.Popen", "subprocess.check_call", "subprocess.check_output") {
		if hasKwargTrue(call.ArgsText, "shell") {
			continue
		}
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) == 0 || !looksDynamicArgv(args[0]) {
			continue
		}
		pushInjectionFinding(unit, call.Start, &MetaCWE88,
			"dynamic value is embedded in a subprocess argument vector and can become an unintended option", 0.7, out)
		return
	}
}

func looksDynamicArgv(expr string) bool {
	t := strings.TrimSpace(expr)
	if len(t) < 2 || (t[0] != '[' && t[0] != '(') {
		return false
	}
	close := byte(']')
	if t[0] == '(' {
		close = ')'
	}
	if t[len(t)-1] != close {
		return false
	}
	for _, arg := range splitTopLevelArgs(t[1 : len(t)-1]) {
		if isDynamicExpr(arg) {
			return true
		}
	}
	return false
}

// CWE-117: keep the signal to dynamically formatted messages passed to known
// Python logging APIs. Structured literal messages with separate arguments are
// intentionally excluded.
func detectCWE117(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(unit.Source,
		"logging.debug", "logging.info", "logging.warning", "logging.error", "logging.critical",
		"logger.debug", "logger.info", "logger.warning", "logger.error", "logger.critical",
		"log.debug", "log.info", "log.warning", "log.error", "log.critical") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) == 0 || !looksLogFormatted(args[0]) || headerValueLooksSanitized(args[0]) {
			continue
		}
		pushInjectionFinding(unit, call.Start, &MetaCWE117,
			"dynamic value is formatted directly into a log message without CRLF neutralization", 0.68, out)
		return
	}
}

func looksLogFormatted(expr string) bool {
	t := strings.TrimSpace(expr)
	return isFStringExpr(t) || strings.Contains(t, ".format(") || indexTopLevelPercent(t) > 0 ||
		(strings.Contains(t, "+") && !isPureStringLiteral(t))
}

func pushInjectionFinding(unit *core.ParsedUnit, offset int, meta *rules.RuleMetadata, message string, confidence float32, out *[]rules.Finding) {
	line, col := unit.LineCol(offset)
	rules.PushFindingWithConfidence(meta, unitFile(unit), line, col, message, confidence, out)
}
