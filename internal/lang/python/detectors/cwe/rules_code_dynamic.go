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
// expression from a genuinely untrusted control sphere. Package-controlled
// literal module names, own-package references (__name__/__package__), and
// own-tree enumeration (pkgutil.walk_packages, glob/rglob over the package's
// own directories) are suppressed, as are test harnesses that import
// repo-committed content.
func detectCWE829(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if isPythonTestModule(unit) {
		return
	}
	for _, name := range []string{"__import__", "importlib.import_module", "runpy.run_path"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			args := splitTopLevelArgs(call.ArgsText)
			if len(args) > 0 && isDynamicExpr(args[0]) && !importTargetIsDevControlled(unit.Source, call.Start, args[0]) {
				emitCodeDynamic(unit, &MetaCWE829, call.Start,
					"dynamically selected module or code path is loaded for execution", confidence82, out)
				return
			}
		}
	}
	for _, call := range findCalls(facts, unit.Source, "importlib.util.spec_from_file_location") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) >= 2 && isDynamicExpr(args[1]) && !importTargetIsDevControlled(unit.Source, call.Start, args[1]) {
			emitCodeDynamic(unit, &MetaCWE829, call.Start,
				"dynamically selected file path is prepared for module execution", confidence80, out)
			return
		}
	}
}

// --- Developer-controlled import targets (CWE-829 / CWE-94) ---

// importTargetIsDevControlled reports whether a dynamic-import name or path
// derives only from developer-controlled source rather than runtime input:
// string literals, the current package's own names (__name__/__package__),
// enumeration of the package's own tree (pkgutil.walk_packages, glob/rglob
// over __file__-relative directories), local assignment chains, or loops over
// literal collections. Audited FP shapes: niquests packages.py (own fixed
// package list), pdf_oxide benchmark (literal dict), rendercv cli rglob,
// Project_Parva lazy __getattr__ shim, httpmorph own-directory glob. The
// audited TP (voicetag _PROVIDERS[provider] with a runtime key) stays dynamic.
func importTargetIsDevControlled(source string, callStart int, expr string) bool {
	return importTargetDevControlled(source, callStart, expr, map[string]bool{})
}

func importTargetDevControlled(source string, callStart int, expr string, seen map[string]bool) bool {
	t := strings.TrimSpace(expr)
	if t == "" {
		return false
	}
	if isPureStringLiteral(t) || looksStaticStringList(t) || looksStaticDictLiteral(t) || isNumericLiteral(t) {
		return true
	}
	// f-string: every interpolation must itself be developer-controlled.
	// (Checked before the __name__/__package__ and own-tree tests so a single
	// runtime interpolation cannot hide behind a static prefix.)
	if isFStringExpr(t) {
		for _, part := range fstringInterpolations(t) {
			if strings.TrimSpace(part) != "" && !importTargetDevControlled(source, callStart, part, seen) {
				return false
			}
		}
		return true
	}
	// Current-package self references.
	if strings.Contains(t, "__name__") || strings.Contains(t, "__package__") {
		return true
	}
	// Own-directory / own-tree enumeration inside the expression.
	if strings.Contains(t, "walk_packages(") || strings.Contains(t, "rglob(") ||
		strings.Contains(t, "glob.glob(") || strings.Contains(t, "listdir(") {
		return true
	}
	// Dict subscript / allowlist .get over a developer-controlled base.
	if base, key, isGet, ok := splitIndexOrGet(t); ok {
		if !importTargetDevControlled(source, callStart, base, seen) {
			return false
		}
		if isGet {
			// .get(<any key>) on a literal allowlist dict yields one of the
			// dict's literal values regardless of the key.
			return true
		}
		return importTargetDevControlled(source, callStart, key, seen)
	}
	// Attribute access on a developer-controlled base (pkgutil walk results).
	if root, ok := dottedAttrParts(t); ok {
		return importTargetDevControlled(source, callStart, root, seen)
	}
	if isIdentOnly(t) {
		if seen[t] {
			return false
		}
		if loopTargetFromStaticCollection(source, callStart, t) {
			return true
		}
		seen[t] = true
		return localChainStatic(source, callStart, t, seen)
	}
	return false
}

// splitIndexOrGet recognizes M[k] and M.get(k) top-level expressions.
func splitIndexOrGet(expr string) (base, key string, isGet, ok bool) {
	t := strings.TrimSpace(expr)
	if t == "" {
		return "", "", false, false
	}
	if t[len(t)-1] == ']' {
		open := strings.LastIndexByte(t, '[')
		if open <= 0 {
			return "", "", false, false
		}
		return strings.TrimSpace(t[:open]), strings.TrimSpace(t[open+1 : len(t)-1]), false, true
	}
	idx := strings.LastIndex(t, ".get(")
	if idx <= 0 || !strings.HasSuffix(t, ")") {
		return "", "", false, false
	}
	open := idx + len(".get")
	if open >= len(t) || t[open] != '(' {
		return "", "", false, false
	}
	closeAt, inner := scanCallArgs(t, open)
	if closeAt < 0 || strings.TrimSpace(t[closeAt+1:]) != "" {
		return "", "", false, false
	}
	args := splitTopLevelArgs(inner)
	if len(args) < 1 {
		return "", "", false, false
	}
	return strings.TrimSpace(t[:idx]), args[0], true, true
}

// dottedAttrParts splits a pure identifier attribute chain at its first dot.
func dottedAttrParts(expr string) (string, bool) {
	t := strings.TrimSpace(expr)
	i := strings.IndexByte(t, '.')
	if i <= 0 {
		return "", false
	}
	root := t[:i]
	if !isIdentOnly(root) {
		return "", false
	}
	for _, seg := range strings.Split(t[i+1:], ".") {
		if !isIdentOnly(seg) {
			return "", false
		}
	}
	return root, true
}

// looksStaticDictLiteral reports a dict literal whose keys and values are all
// string literals — a developer-controlled allowlist mapping.
func looksStaticDictLiteral(expr string) bool {
	t := strings.TrimSpace(expr)
	if len(t) < 2 || t[0] != '{' || t[len(t)-1] != '}' {
		return false
	}
	inner := strings.TrimSpace(t[1 : len(t)-1])
	if inner == "" {
		return true
	}
	for _, item := range splitTopLevelArgs(inner) {
		if strings.TrimSpace(item) == "" {
			continue
		}
		key, val, ok := strings.Cut(item, ":")
		if !ok || !isPureStringLiteral(key) || !isPureStringLiteral(val) {
			return false
		}
	}
	return true
}

// localChainStatic walks the lines before callStart for the nearest assignment
// to name and checks whether its right-hand side is developer-controlled.
// def/class/for boundaries stop the walk only when they bind the name
// (parameter, loop target); module-level constants above a function are fair
// game, while a shadowing parameter keeps the name runtime-controlled.
func localChainStatic(source string, callStart int, name string, seen map[string]bool) bool {
	if callStart <= 0 || callStart > len(source) {
		return false
	}
	lines := strings.Split(source[:callStart], "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if strings.HasPrefix(trimmed, "def ") || strings.HasPrefix(trimmed, "async def ") ||
			strings.HasPrefix(trimmed, "class ") {
			if strings.HasPrefix(trimmed, "class ") || defParamsContain(trimmed, name) {
				return false
			}
			continue
		}
		if strings.HasPrefix(trimmed, "for ") {
			if loopLineTargetsName(trimmed, name) {
				return false
			}
			continue
		}
		lhs, rhs, ok := strings.Cut(trimmed, "=")
		if !ok {
			continue
		}
		lhs = strings.TrimSpace(lhs)
		if lhs != name {
			// tolerate `name: Type = value` annotations
			if colon := strings.IndexByte(lhs, ':'); colon < 0 || strings.TrimSpace(lhs[:colon]) != name {
				continue
			}
		}
		rhs = stripPyComment(strings.TrimSpace(rhs))
		if rhs == "" {
			return false
		}
		// Multi-line right-hand sides (dict / string literals spanning lines).
		if rhs == "{" || rhs == "[" || rhs == "(" {
			rhs = joinRHSSource(lines, i, rhs)
		}
		return importTargetDevControlled(source, callStart, rhs, seen)
	}
	return false
}

// defParamsContain reports whether a def/async def header declares name as a
// parameter (including *args/**kwargs-style bindings).
func defParamsContain(header string, name string) bool {
	open := strings.IndexByte(header, '(')
	if open < 0 {
		return false
	}
	inner := header[open+1:]
	if i := strings.LastIndexByte(inner, ')'); i >= 0 {
		inner = inner[:i]
	}
	for _, p := range strings.Split(inner, ",") {
		p = strings.TrimSpace(p)
		if i := strings.IndexByte(p, ':'); i >= 0 {
			p = p[:i]
		}
		if i := strings.IndexByte(p, '='); i >= 0 {
			p = p[:i]
		}
		p = strings.TrimSpace(p)
		if strings.HasPrefix(p, "*") {
			p = strings.TrimPrefix(p, "*")
			p = strings.TrimPrefix(p, "*")
		}
		if p == name {
			return true
		}
	}
	return false
}

// loopLineTargetsName reports whether a for line binds name as one of its loop
// targets (`for name in ...` or `for a, name in ...`).
func loopLineTargetsName(line string, name string) bool {
	body := strings.TrimSpace(strings.TrimPrefix(line, "for "))
	if inIdx := strings.Index(body, " in "); inIdx >= 0 {
		body = body[:inIdx]
	}
	if body == name {
		return true
	}
	for _, part := range strings.Split(body, ",") {
		if strings.TrimSpace(part) == name {
			return true
		}
	}
	return false
}

// stripPyComment removes a trailing Python comment, honoring string literals.
func stripPyComment(line string) string {
	inStr := byte(0)
	escaped := false
	for i := 0; i < len(line); i++ {
		c := line[i]
		if inStr != 0 {
			if escaped {
				escaped = false
				continue
			}
			if c == '\\' {
				escaped = true
				continue
			}
			if c == inStr {
				inStr = 0
			}
			continue
		}
		if c == '\'' || c == '"' {
			inStr = c
			continue
		}
		if c == '#' {
			return strings.TrimSpace(line[:i])
		}
	}
	return line
}

// joinRHSSource folds the following lines into a multi-line right-hand side
// until its brackets balance (string-aware).
func joinRHSSource(lines []string, start int, rhs string) string {
	depth := pyOpenBracketBalance(rhs)
	var b strings.Builder
	b.WriteString(rhs)
	for j := start + 1; j < len(lines) && depth > 0; j++ {
		next := strings.TrimSpace(lines[j])
		if next == "" {
			continue
		}
		depth += pyOpenBracketBalance(next)
		b.WriteByte(' ')
		b.WriteString(next)
	}
	return b.String()
}

// pyOpenBracketBalance returns the ( [ { minus ) ] } depth change of s,
// ignoring string-literal and comment content.
func pyOpenBracketBalance(s string) int {
	depth := 0
	inStr := byte(0)
	escaped := false
	inComment := false
	triple := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if inComment {
			if c == '\n' {
				inComment = false
			}
			continue
		}
		if inStr != 0 {
			if escaped {
				escaped = false
				continue
			}
			if c == '\\' {
				escaped = true
				continue
			}
			if triple {
				if c == inStr && i+2 < len(s) && s[i+1] == inStr && s[i+2] == inStr {
					inStr = 0
					triple = false
					i += 2
				}
				continue
			}
			if c == inStr {
				inStr = 0
			}
			continue
		}
		switch c {
		case '#':
			inComment = true
		case '\'', '"':
			inStr = c
			if i+2 < len(s) && s[i+1] == c && s[i+2] == c {
				triple = true
				i += 2
			}
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			depth--
		}
	}
	return depth
}

// loopTargetFromStaticCollection reports whether the nearest enclosing for
// loop over target iterates a developer-controlled collection (literal tuple,
// own-tree enumeration, or .items()/.values() of a literal dict).
func loopTargetFromStaticCollection(source string, callStart int, target string) bool {
	if callStart <= 0 || callStart > len(source) {
		return false
	}
	lines := strings.Split(source[:callStart], "\n")
	offsets := make([]int, len(lines))
	off := 0
	for i, l := range lines {
		offsets[i] = off
		off += len(l) + 1
	}
	for i := len(lines) - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "def ") || strings.HasPrefix(trimmed, "async def ") ||
			strings.HasPrefix(trimmed, "class ") {
			return false
		}
		if !strings.HasPrefix(trimmed, "for ") {
			continue
		}
		body := trimmed[len("for "):]
		inIdx := strings.Index(body, " in ")
		if inIdx < 0 {
			continue
		}
		targets := strings.TrimSpace(body[:inIdx])
		collection := strings.TrimSpace(body[inIdx+len(" in "):])
		collection = strings.TrimSpace(strings.TrimSuffix(collection, ":"))
		if !loopTargetMatches(targets, target) {
			continue
		}
		return loopCollectionStatic(source, offsets[i], collection)
	}
	return false
}

func loopTargetMatches(targets, target string) bool {
	if targets == target {
		return true
	}
	if !strings.Contains(targets, ",") {
		return false
	}
	for _, part := range strings.Split(targets, ",") {
		if strings.TrimSpace(part) == target {
			return true
		}
	}
	return false
}

func loopCollectionStatic(source string, forLineStart int, collection string) bool {
	t := strings.TrimSpace(collection)
	if t == "" {
		return false
	}
	if looksStaticStringList(t) {
		return true
	}
	if strings.Contains(t, "walk_packages(") || strings.Contains(t, "rglob(") ||
		strings.Contains(t, "glob.glob(") || strings.Contains(t, "listdir(") ||
		strings.Contains(t, "__name__") || strings.Contains(t, "__package__") {
		return true
	}
	for _, method := range []string{".items()", ".keys()", ".values()", ".iteritems()"} {
		if strings.HasSuffix(t, method) {
			return localChainStatic(source, forLineStart, strings.TrimSpace(t[:len(t)-len(method)]), map[string]bool{})
		}
	}
	if isIdentOnly(t) {
		return localChainStatic(source, forLineStart, t, map[string]bool{})
	}
	return false
}

// fstringInterpolations returns the {…} expression bodies of an f-string
// literal expression, skipping escaped braces ({{, }}).
func fstringInterpolations(expr string) []string {
	t := strings.TrimSpace(expr)
	i := 0
	for i < len(t) {
		c := t[i]
		if c == 'f' || c == 'F' || c == 'r' || c == 'R' || c == 'b' || c == 'B' || c == 'u' || c == 'U' {
			i++
			continue
		}
		break
	}
	if i >= len(t) {
		return nil
	}
	quote := t[i]
	if quote != '"' && quote != '\'' {
		return nil
	}
	var body string
	if i+2 < len(t) && t[i+1] == quote && t[i+2] == quote {
		after := t[i+3:]
		end := strings.Index(after, string([]byte{quote, quote, quote}))
		if end < 0 {
			return nil
		}
		body = after[:end]
	} else {
		escaped := false
		for j := i + 1; j < len(t); j++ {
			if escaped {
				escaped = false
				continue
			}
			if t[j] == '\\' {
				escaped = true
				continue
			}
			if t[j] == quote {
				body = t[i+1 : j]
				break
			}
		}
		if body == "" {
			return nil
		}
	}
	return collectFstringParts(body)
}

func collectFstringParts(body string) []string {
	var out []string
	depth := 0
	var cur strings.Builder
	for k := 0; k < len(body); k++ {
		c := body[k]
		if c == '{' {
			if k+1 < len(body) && body[k+1] == '{' {
				k++
				continue
			}
			if depth == 0 {
				cur.Reset()
			}
			depth++
			continue
		}
		if c == '}' {
			if k+1 < len(body) && body[k+1] == '}' {
				k++
				continue
			}
			depth--
			if depth == 0 {
				out = append(out, cur.String())
				cur.Reset()
			}
			continue
		}
		if depth > 0 {
			cur.WriteByte(c)
		}
	}
	return out
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
// The English word "password" inside a message string is not a sensitive value.
//
// Offline release/CLI tooling (scripts/release key generators, tools/) that
// deliberately print a freshly minted secret to stdout for the operator is
// not deployment debug code.
func detectCWE215(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if isPythonBenchmarkFile(unit) {
		return
	}
	for _, name := range []string{"print", "logging.debug", ".debug"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			if sensitiveIdentOutsideLiterals(call.ArgsText) && !isPureStringLiteral(call.ArgsText) {
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
