package cwe

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	// Gate on LDAP library evidence only. Bare ".search(" matches re.Pattern.search
	// and must not schedule this rule on non-LDAP files.
	RegisterRule("CWE-90", detectCWE90, &MetaCWE90,
		"ldap3", "ldap.initialize", "import ldap", "from ldap ", ".search_s(")
	RegisterRule("CWE-91", detectCWE91, &MetaCWE91,
		".xpath(", "XPath(", ".fromstring(")
	RegisterRule("CWE-93", detectCWE93, &MetaCWE93,
		"response.headers[", ".headers[", ".set_header(", ".add_header(", "HttpResponseRedirect(")
	RegisterRule("CWE-94", detectCWE94, &MetaCWE94,
		"eval(", "exec(", "compile(", "__import__(", "importlib.import_module")
	RegisterRule("CWE-88", detectCWE88, &MetaCWE88,
		"subprocess.")
	RegisterRule("CWE-117", detectCWE117, &MetaCWE117,
		"logging.", "logger.", "log.")
}

// CWE-90: only dynamic LDAP filter expressions reach LDAP search APIs. Filter
// values escaped through the standard ldap3 helper are intentionally suppressed.
// Regex .search / Pattern.search sinks are excluded via LDAP library context and
// receiver shape (re / *_RE), not treated as LDAP injection.
func detectCWE90(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if !pythonFileLooksLikeLDAP(unit.Source) {
		return
	}
	for _, call := range findCalls(facts, unit.Source, ".search", ".search_s") {
		if call.Name == ".search" && looksLikeRegexSearchReceiver(methodReceiver(unit.Source, call.Start)) {
			continue
		}
		args := splitTopLevelArgs(call.ArgsText)
		filterArg := ldapFilterArg(args, call.Name)
		if filterArg == "" || !isDynamicExpr(filterArg) || ldapFilterLooksEscaped(filterArg) {
			continue
		}
		pushInjectionFinding(unit, call.Start, &MetaCWE90,
			"dynamic value reaches an LDAP search filter without LDAP filter escaping", confidence78, out)
		return
	}
}

func pythonFileLooksLikeLDAP(source string) bool {
	return strings.Contains(source, "ldap3") ||
		strings.Contains(source, "ldap.initialize") ||
		strings.Contains(source, "import ldap") ||
		strings.Contains(source, "from ldap ") ||
		strings.Contains(source, ".search_s(")
}

func ldapFilterLooksEscaped(expr string) bool {
	return strings.Contains(expr, "escape_filter_chars(") || strings.Contains(expr, "escape_filter_value(")
}

// ldapFilterArg picks the LDAP filter operand. python-ldap search_s is
// (base, scope, filter, ...); ldap3 Connection.search is often (base, filter, ...)
// or a single filter expression in simplified call sites.
func ldapFilterArg(args []string, callName string) string {
	switch {
	case callName == ".search_s" && len(args) >= 3:
		return args[2]
	case len(args) >= 2:
		return args[1]
	case len(args) == 1:
		return args[0]
	default:
		return ""
	}
}

// methodReceiver returns the identifier immediately left of a ".method" match
// start (call.Start points at the '.').
func methodReceiver(source string, dotAt int) string {
	if dotAt <= 0 || dotAt > len(source) {
		return ""
	}
	i := dotAt - 1
	for i >= 0 {
		c := source[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
			i--
			continue
		}
		break
	}
	return source[i+1 : dotAt]
}

func looksLikeRegexSearchReceiver(recv string) bool {
	r := strings.TrimSpace(recv)
	if r == "" {
		return false
	}
	if r == "re" || r == "regex" {
		return true
	}
	return strings.HasSuffix(r, "_RE") || strings.HasSuffix(r, "_re") ||
		strings.HasSuffix(r, "Pattern") || strings.HasSuffix(r, "pattern")
}

// CWE-91: flag dynamic XML/XPath expression construction at XPath APIs and at
// XML parse APIs only when the document argument is itself constructed
// (f-string / format / concat). Bare ET.fromstring(xml_data) is document
// parsing, not injection. Literal XPath with bound variables remains safe.
func detectCWE91(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(facts, unit.Source, ".xpath", "etree.XPath", "ElementTree.XPath", "ET.XPath", "etree.fromstring", "ElementTree.fromstring", "ET.fromstring") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) == 0 || !isDynamicExpr(args[0]) {
			continue
		}
		if isXMLParseCall(call.Name) && !looksXMLConstructed(args[0]) {
			continue
		}
		pushInjectionFinding(unit, call.Start, &MetaCWE91,
			"dynamic value is formatted into an XML or XPath expression", confidence75, out)
		return
	}
}

func isXMLParseCall(name string) bool {
	switch name {
	case "etree.fromstring", "ElementTree.fromstring", "ET.fromstring", ".fromstring":
		return true
	default:
		return strings.HasSuffix(name, ".fromstring")
	}
}

func looksXMLConstructed(expr string) bool {
	t := strings.TrimSpace(expr)
	return isFStringExpr(t) || strings.Contains(t, ".format(") || indexTopLevelPercent(t) > 0 ||
		(strings.Contains(t, "+") && !isPureStringLiteral(t))
}

// CWE-93: header sinks are limited to explicit header mutation APIs. Values
// that visibly remove both CR and LF are not reported by this source heuristic.
// Header reads and == comparisons are not sinks (assignment '=' only), and
// writes on outgoing client request objects (client.headers[...]) are not
// response-header writes.
func detectCWE93(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if isPythonTestModule(unit) {
		return
	}
	for _, call := range findCalls(facts, unit.Source, ".set_header", ".add_header", "HttpResponseRedirect") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) < 2 || !isDynamicExpr(args[len(args)-1]) || headerValueLooksSanitized(args[len(args)-1]) {
			continue
		}
		pushInjectionFinding(unit, call.Start, &MetaCWE93,
			"dynamic value is written to an HTTP response header without CRLF neutralization", confidence72, out)
		return
	}
	masked := ""
	if facts != nil {
		masked = facts.Masked
	}
	if masked == "" {
		masked = pythonCodeMask(unit.Source)
	}
	originalLines := strings.Split(unit.Source, "\n")
	maskedLines := strings.Split(masked, "\n")
	lineOffset := 0
	for i, line := range originalLines {
		maskedLine := ""
		if i < len(maskedLines) {
			maskedLine = maskedLines[i]
		}
		lhs, _, ok := splitAssignmentEq(maskedLine)
		_, rhs, rhsOK := splitAssignmentEq(line)
		if !ok || !rhsOK || !strings.Contains(lhs, ".headers[") || !headerWriteReceiverLooksResponse(lhs) || !isDynamicExpr(rhs) || headerValueLooksSanitized(rhs) || headerValueIsInternalNumeric(rhs) || headerValueLooksConstant(rhs, facts, unit.Source) {
			lineOffset += len(line) + 1
			continue
		}
		start := lineOffset + strings.Index(line, ".headers[")
		pushInjectionFinding(unit, start, &MetaCWE93,
			"dynamic value is written to an HTTP response header without CRLF neutralization", confidence72, out)
		return
	}
}

func headerValueLooksSanitized(expr string) bool {
	compact := compactWhitespace(expr)
	return (strings.Contains(compact, `replace("\r","")`) || strings.Contains(compact, `replace('\r','')`)) &&
		(strings.Contains(compact, `replace("\n","")`) || strings.Contains(compact, `replace('\n','')`))
}

func headerValueIsInternalNumeric(expr string) bool {
	compact := compactWhitespace(expr)
	if strings.Contains(compact, "str(int(") || strings.Contains(compact, "str(round(") {
		return true
	}
	if strings.Contains(compact, ".shape") {
		return true
	}
	if strings.Contains(compact, "os.path.getsize(") || strings.Contains(compact, ".getsize(") {
		return true
	}
	return false
}

// CWE-94: code-generation and dynamic-import APIs are only findings when the
// code/module argument is non-literal. Literal developer-owned expressions are
// out of scope for this same-file heuristic. Attribute methods named exec/eval
// (SQLModel session.exec, Qt app.exec) are rejected by findCalls' left-boundary
// check on '.' unless the receiver is the explicit builtins module; def
// signatures named exec/eval are declarations, not call sites; dynamic imports
// over developer-controlled constants (own-package names, literal allowlists,
// own-tree enumeration) are not untrusted control spheres; test harnesses
// exercising repo-committed code are not deployment targets.
func detectCWE94(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if isPythonTestModule(unit) {
		return
	}
	for _, call := range findCalls(facts, unit.Source, "eval", "exec", "compile", "builtins.exec", "builtins.eval", "__import__", "importlib.import_module") {
		if call.Start > 0 && unit.Source[call.Start-1] == '.' {
			if !cweReceiverIsBuiltins(unit.Source, call.Start) {
				continue
			}
		}
		switch call.Name {
		case "eval", "exec", "compile":
			if isDefHeaderCall(unit.Source, call.Start) {
				continue
			}
		case "__import__", "importlib.import_module":
			args := splitTopLevelArgs(call.ArgsText)
			if len(args) > 0 && importTargetIsDevControlled(unit.Source, call.Start, args[0]) {
				continue
			}
		}
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) == 0 || !isDynamicExpr(args[0]) {
			continue
		}
		pushInjectionFinding(unit, call.Start, &MetaCWE94,
			"dynamic value reaches a Python code-generation or dynamic-import sink", confidence82, out)
		return
	}
}

// cweReceiverIsBuiltins reports whether the attribute receiver immediately
// left of the identifier at identAt is the explicit builtins module
// (builtins.exec / builtins.eval are the builtin, unlike session.exec).
func cweReceiverIsBuiltins(source string, identAt int) bool {
	dot := identAt - 1
	if dot <= 0 || source[dot] != '.' {
		return false
	}
	end := dot
	start := end
	for start > 0 {
		c := source[start-1]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
			start--
			continue
		}
		break
	}
	return source[start:end] == "builtins"
}

// CWE-88: shell=False is not sufficient when an untrusted argument can become
// a tool option. Detect only explicit argv literals with a dynamic segment; a
// pre-built argv variable has insufficient same-file evidence to report.
// Conventional Python test/benchmark modules are skipped. Module constants
// (UPPER_SNAKE), sys.executable, internally generated values (tempfile/socket
// paths, ports), literal+constant concatenations, and self re-exec argv
// spreads are not treated as attacker-controlled option text.
func detectCWE88(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if isPythonTestModule(unit) || isPythonBenchmarkFile(unit) {
		return
	}
	ctx := buildPyFileCtx(facts, unit.Source)
	for _, call := range findCalls(facts, unit.Source,
		"subprocess.run", "subprocess.call", "subprocess.Popen", "subprocess.check_call", "subprocess.check_output") {
		if hasKwargTrue(call.ArgsText, "shell") {
			continue
		}
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) == 0 || !looksDynamicArgv(args[0], ctx) {
			continue
		}
		pushInjectionFinding(unit, call.Start, &MetaCWE88,
			"dynamic value is embedded in a subprocess argument vector and can become an unintended option", confidence70, out)
		return
	}
}

func looksDynamicArgv(expr string, ctx *pyFileCtx) bool {
	t := strings.TrimSpace(expr)
	if len(t) < 2 || (t[0] != '[' && t[0] != '(') {
		return false
	}
	closeDelimiter := byte(']')
	if t[0] == '(' {
		closeDelimiter = ')'
	}
	if t[len(t)-1] != closeDelimiter {
		return false
	}
	segments := splitTopLevelArgs(t[1 : len(t)-1])
	// Self re-exec of the operator's own argv under a guard is not injection.
	for _, seg := range segments {
		compact := compactWhitespace(seg)
		if compact == "*sys.argv" || strings.HasPrefix(compact, "*sys.argv[") {
			return false
		}
	}
	for i, seg := range segments {
		if argvSegmentIsOptionValue(segments, i, ctx) {
			continue
		}
		if argvSegmentLooksDynamic(seg, ctx) {
			return true
		}
	}
	return false
}

// argvSegmentIsOptionValue reports a segment bound as the value of the named
// option that precedes it: a bare identifier or str() of a numeric-typed
// parameter in that position is consumed by the tool as the option's value and
// cannot be parsed as an unintended switch in a no-shell argv.
func argvSegmentIsOptionValue(segments []string, i int, ctx *pyFileCtx) bool {
	if i == 0 {
		return false
	}
	prev := strings.TrimSpace(segments[i-1])
	if len(prev) >= 2 && (prev[0] == '"' || prev[0] == '\'') && prev[len(prev)-1] == prev[0] {
		prev = prev[1 : len(prev)-1]
	}
	prev = compactWhitespace(prev)
	if len(prev) < 2 || prev[0] != '-' || prev == "--" {
		return false
	}
	t := strings.TrimSpace(segments[i])
	if bareIdentRE.MatchString(t) {
		return true
	}
	if strings.HasPrefix(t, "str(") && strings.HasSuffix(t, ")") {
		inner := strings.TrimSpace(t[4 : len(t)-1])
		if bareIdentRE.MatchString(inner) && ctx != nil && ctx.numericAnnotated[inner] {
			return true
		}
	}
	return false
}

// argvSegmentLooksDynamic is true when a list/tuple argv element can carry
// attacker-controlled option text. Pure string/bytes literals, static
// concatenations of those literals, and trusted module constants are not dynamic.
func argvSegmentLooksDynamic(expr string, ctx *pyFileCtx) bool {
	t := strings.TrimSpace(expr)
	if t == "" {
		return false
	}
	if isPureStringLiteral(t) || isPureBytesLiteral(t) {
		return false
	}
	if looksStaticStringList(t) {
		return false
	}
	if looksStaticLiteralConcat(t) {
		return false
	}
	if argvSegmentLooksTrusted(t, ctx) {
		return false
	}
	return true
}

// headerWriteReceiverLooksResponse restricts .headers[...] write sinks to
// response objects (response/resp/res/self), excluding outgoing client or
// request-object header mutation.
func headerWriteReceiverLooksResponse(lhs string) bool {
	idx := strings.Index(lhs, ".headers[")
	if idx < 0 {
		return false
	}
	recv := strings.TrimSpace(lhs[:idx])
	if i := strings.LastIndex(recv, "."); i >= 0 {
		recv = strings.TrimSpace(recv[i+1:])
	}
	lower := strings.ToLower(recv)
	if lower == "self" || lower == "res" {
		return true
	}
	return strings.Contains(lower, "response") || strings.HasPrefix(lower, "resp")
}

// headerValueLooksConstant resolves a header value that is a same-file call
// whose callee returns only literals (constant helpers), or a bare module
// constant assigned a literal.
func headerValueLooksConstant(expr string, facts *PyCweFacts, source string) bool {
	t := strings.TrimSpace(expr)
	if t == "" {
		return false
	}
	if upperSnakeRE.MatchString(t) && source != "" {
		return moduleConstAssignedLiteral(t, source)
	}
	name := calleeNameOf(t)
	if name == "" {
		return false
	}
	return functionReturnsOnlyLiteralsForName(name, source)
}

// calleeNameOf returns the plain callee identifier of a call expression.
func calleeNameOf(expr string) string {
	t := strings.TrimSpace(expr)
	open := strings.IndexByte(t, '(')
	if open <= 0 || !strings.HasSuffix(t, ")") {
		return ""
	}
	callee := strings.TrimSpace(t[:open])
	if i := strings.LastIndex(callee, "."); i >= 0 {
		callee = strings.TrimSpace(callee[i+1:])
	}
	if !bareIdentRE.MatchString(callee) {
		return ""
	}
	return callee
}

// moduleConstAssignedLiteral reports UPPER_SNAKE = <literal> evidence.
func moduleConstAssignedLiteral(name, source string) bool {
	masked := pythonCodeMask(source)
	for _, line := range strings.Split(masked, "\n") {
		lhs, rhs, ok := splitAssignmentEq(line)
		if !ok || strings.TrimSpace(lhs) != name {
			continue
		}
		return isPureStringLiteral(rhs) || isNumericLiteral(rhs)
	}
	return false
}

// functionReturnsOnlyLiterals reports a def body whose every return statement
// returns a literal (constant helper).
func functionReturnsOnlyLiterals(body string) bool {
	if body == "" {
		return false
	}
	lines := buildMaskedPythonLines(pythonCodeMask(body))
	seenReturn := false
	for _, line := range lines {
		masked := strings.TrimSpace(line.text)
		if masked == "" {
			continue
		}
		if masked != "return" && !strings.HasPrefix(masked, "return ") && !strings.HasPrefix(masked, "return(") {
			continue
		}
		seenReturn = true
		orig := body[line.start:]
		if len(orig) > len(line.text) {
			orig = orig[:len(line.text)]
		}
		expr := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(orig), "return"))
		if strings.HasPrefix(expr, "(") && strings.HasSuffix(expr, ")") {
			expr = strings.TrimSpace(expr[1 : len(expr)-1])
		}
		if expr == "" || isPureStringLiteral(expr) || isNumericLiteral(expr) {
			continue
		}
		return false
	}
	return seenReturn
}

// functionReturnsOnlyLiteralsForName locates a same-file def (arrow-annotation
// tolerant) and analyzes its body returns. Facts.Funcs cannot be reused here
// because its def regex requires `):` and misses `-> str:` annotations.
func functionReturnsOnlyLiteralsForName(name, source string) bool {
	if name == "" || source == "" {
		return false
	}
	re := regexp.MustCompile(`(?m)^[ \t]*def\s+` + regexp.QuoteMeta(name) + `\s*\([^\n]*\)\s*(->[^\n]*)?:`)
	masked := pythonCodeMask(source)
	loc := re.FindStringIndex(masked)
	if loc == nil {
		return false
	}
	colonAt := loc[1] - 1 // match ends at ':'
	indent := 0
	for loc[0]+indent < colonAt && (masked[loc[0]+indent] == ' ' || masked[loc[0]+indent] == '\t') {
		indent++
	}
	bodyStart := colonAt + 1
	if bodyStart < len(masked) && masked[bodyStart] == '\n' {
		bodyStart++
	}
	if bodyStart > len(source) {
		return false
	}
	bodyEnd := pythonFunctionBodyEnd(masked, bodyStart, indent)
	if bodyStart > bodyEnd || bodyEnd > len(source) {
		return false
	}
	return functionReturnsOnlyLiterals(source[bodyStart:bodyEnd])
}

// isPureBytesLiteral reports a b"..." / B'...' literal.
func isPureBytesLiteral(expr string) bool {
	t := strings.TrimSpace(expr)
	if len(t) < 3 {
		return false
	}
	if (t[0] == 'b' || t[0] == 'B') && (t[1] == '"' || t[1] == '\'') {
		return isPureStringLiteral(t[1:])
	}
	return false
}

func looksStaticLiteralConcat(expr string) bool {
	t := strings.TrimSpace(expr)
	if !strings.Contains(t, "+") {
		return false
	}
	parts := splitTopLevelConcat(t)
	if len(parts) < 2 {
		return false
	}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if !isPureStringLiteral(p) && !isPureBytesLiteral(p) {
			return false
		}
	}
	return true
}

func splitTopLevelConcat(expr string) []string {
	var parts []string
	depth := 0
	start := 0
	inStr := byte(0)
	esc := false
	for i := 0; i < len(expr); i++ {
		c := expr[i]
		if inStr != 0 {
			if esc {
				esc = false
				continue
			}
			if c == '\\' {
				esc = true
				continue
			}
			if c == inStr {
				inStr = 0
			}
			continue
		}
		switch c {
		case '\'', '"':
			inStr = c
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			if depth > 0 {
				depth--
			}
		case '+':
			if depth == 0 {
				parts = append(parts, strings.TrimSpace(expr[start:i]))
				start = i + 1
			}
		}
	}
	parts = append(parts, strings.TrimSpace(expr[start:]))
	return parts
}

// CWE-117: keep the signal to dynamically formatted messages passed to known
// Python logging APIs. Structured literal messages with separate arguments are
// intentionally excluded. Constant / numeric-only interpolations, loop
// counters, and internally generated names cannot carry CRLF.
func detectCWE117(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	ctx := buildPyFileCtx(facts, unit.Source)
	for _, call := range findCalls(facts, unit.Source,
		"logging.debug", "logging.info", "logging.warning", "logging.error", "logging.critical",
		"logger.debug", "logger.info", "logger.warning", "logger.error", "logger.critical",
		"log.debug", "log.info", "log.warning", "log.error", "log.critical") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) == 0 || !looksLogFormatted(args[0]) || headerValueLooksSanitized(args[0]) {
			continue
		}
		if !logMessageHasCRLFCapableValue(args[0], ctx) {
			continue
		}
		pushInjectionFinding(unit, call.Start, &MetaCWE117,
			"dynamic value is formatted directly into a log message without CRLF neutralization", confidence68, out)
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
