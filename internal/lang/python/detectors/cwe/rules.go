package cwe

import (
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	RegisterRule("CWE-502", detectCWE502, &MetaCWE502,
		"pickle.loads", "pickle.load", "pickle.Unpickler", "yaml.load", "yaml.unsafe_load")
	RegisterRule("CWE-78", detectCWE78, &MetaCWE78,
		"os.system", "os.popen", "subprocess.", "shell=True", "commands.")
	RegisterRule("CWE-89", detectCWE89, &MetaCWE89,
		"execute(", "executemany(", ".execute(", ".executemany(", "raw(")
	RegisterRule("CWE-22", detectCWE22, &MetaCWE22,
		"open(", "pathlib", "os.path.join", "os.remove", "os.unlink", "Path(")
	RegisterRule("CWE-79", detectCWE79, &MetaCWE79,
		"mark_safe", "Markup(", "render_template_string", "|safe", "HttpResponse(")
}

// --- CWE-502 Deserialization of Untrusted Data ---
//
// Choice: flag any pickle.loads/load/Unpickler call site (museum-style, high signal
// for Python platform CWE). For yaml.load, flag unless Loader is SafeLoader /
// CSafeLoader / FullLoader-safe pattern or call is yaml.safe_load.

func detectCWE502(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	src := unit.Source

	// pickle sinks — any call is treated as unsafe deserialization of potentially
	// untrusted data (conservative; matches priority-batch museum style).
	for _, call := range findCalls(facts, src, "pickle.loads", "pickle.load", "pickle.Unpickler") {
		line, col := unit.LineCol(call.Start)
		rules.PushFindingWithConfidence(
			&MetaCWE502,
			unitFile(unit),
			line, col,
			"unsafe pickle deserialization sink (untrusted data must not be unpickled)",
			confidence80,
			out,
		)
		return
	}

	// yaml.unsafe_load always unsafe
	for _, call := range findCalls(facts, src, "yaml.unsafe_load") {
		line, col := unit.LineCol(call.Start)
		rules.PushFindingWithConfidence(
			&MetaCWE502,
			unitFile(unit),
			line, col,
			"yaml.unsafe_load deserializes untrusted data without a safe loader",
			confidence85,
			out,
		)
		return
	}

	// yaml.load without SafeLoader / CSafeLoader
	for _, call := range findCalls(facts, src, "yaml.load") {
		// Do not treat yaml.safe_load as yaml.load (boundary check already prevents prefix match of safe_load when searching "yaml.load" — verify)
		// "yaml.load" is not a prefix of "yaml.safe_load" in reverse; but "yaml.load" could match inside nothing else.
		if yamlLoadLooksSafe(call.ArgsText) {
			continue
		}
		// ruamel.yaml YAML().load is a safe constructor — not PyYAML's unsafe load.
		if yamlLoadLooksLikeRuamel(src, call.ArgsText) {
			continue
		}
		line, col := unit.LineCol(call.Start)
		rules.PushFindingWithConfidence(
			&MetaCWE502,
			unitFile(unit),
			line, col,
			"yaml.load without SafeLoader (use yaml.safe_load or Loader=yaml.SafeLoader)",
			confidence80,
			out,
		)
		return
	}
}

func yamlLoadLooksSafe(args string) bool {
	compact := compactWhitespace(args)
	// Loader=yaml.SafeLoader / CSafeLoader / SafeDumper not relevant
	if strings.Contains(compact, "Loader=yaml.SafeLoader") ||
		strings.Contains(compact, "Loader=yaml.CSafeLoader") ||
		strings.Contains(compact, "Loader=SafeLoader") ||
		strings.Contains(compact, "Loader=CSafeLoader") {
		return true
	}
	// FullLoader is safer than unsafe_load but still allows arbitrary Python objects in older PyYAML — do not treat as fully safe for this museum.
	return false
}

// --- CWE-78 OS command injection ---

func detectCWE78(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	src := unit.Source

	// os.system / os.popen with dynamic command
	for _, name := range []string{"os.system", "os.popen"} {
		for _, call := range findCalls(facts, src, name) {
			args := splitTopLevelArgs(call.ArgsText)
			if len(args) == 0 {
				continue
			}
			if !isDynamicExpr(args[0]) {
				continue
			}
			line, col := unit.LineCol(call.Start)
			rules.PushFindingWithConfidence(
				&MetaCWE78,
				unitFile(unit),
				line, col,
				"dynamic input reaches an OS command sink (os.system/os.popen)",
				confidence80,
				out,
			)
			return
		}
	}

	// subprocess.* with shell=True and dynamic command
	for _, name := range []string{
		"subprocess.run", "subprocess.call", "subprocess.Popen",
		"subprocess.check_call", "subprocess.check_output",
	} {
		for _, call := range findCalls(facts, src, name) {
			if !hasKwargTrue(call.ArgsText, "shell") {
				continue
			}
			args := splitTopLevelArgs(call.ArgsText)
			if len(args) == 0 {
				continue
			}
			cmd := args[0]
			// strip keyword args that might be first if unusual — typically cmd is first positional
			if strings.HasPrefix(cmd, "shell=") {
				continue
			}
			if !isDynamicExpr(cmd) {
				continue
			}
			line, col := unit.LineCol(call.Start)
			rules.PushFindingWithConfidence(
				&MetaCWE78,
				unitFile(unit),
				line, col,
				"dynamic command with subprocess shell=True (prefer list argv and shell=False)",
				confidence80,
				out,
			)
			return
		}
	}
}

// --- CWE-89 SQL injection ---

func detectCWE89(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	src := unit.Source

	// Prefer method-style .execute( / .executemany( and bare execute(
	// Scan for ".execute(" and "execute(" carefully.
	for _, call := range findExecuteCalls(facts, src) {
		if isDefHeaderCall(src, call.Start) {
			continue
		}
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) == 0 {
			continue
		}
		if isExecuteCallbackPassthrough(facts, src, call.Start, args) {
			continue
		}
		sqlArg := args[0]
		if isBoundQueryBuilderExpression(facts, src, call.Start, sqlArg) {
			continue
		}
		// Static SQL (including adjacent literals and sqlalchemy.text("…")) is not injection.
		if isStaticSQLArg(sqlArg) && !looksSQLFormatted(sqlArg) {
			continue
		}
		// Dynamic SQL construction in first arg — require SQL evidence so HTTP /
		// custom execute wrappers with non-SQL first args stay silent.
		if looksSQLFormatted(sqlArg) || sqlArgLooksInjected(facts, src, call.Start, sqlArg) {
			line, col := unit.LineCol(call.Start)
			rules.PushFindingWithConfidence(
				&MetaCWE89,
				unitFile(unit),
				line, col,
				"dynamic SQL string reaches execute/executemany (use bound parameters)",
				confidence75,
				out,
			)
			return
		}
	}
}

// isBoundQueryBuilderExpression recognizes a local SQLAlchemy-style statement
// variable. These objects carry bound values into execute and are not dynamic
// SQL strings merely because their variable name is non-literal.
func isBoundQueryBuilderExpression(facts *PyCweFacts, source string, callStart int, arg string) bool {
	name := strings.TrimSpace(arg)
	if isQueryBuilderConstructor(name) {
		return true
	}
	if !isIdentOnly(name) || callStart < 0 || callStart > len(source) {
		return false
	}
	prefix := source[:callStart]
	lines := strings.Split(facts.codeMask(prefix, 0), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		lhs, rhs, ok := strings.Cut(line, "=")
		if !ok || strings.TrimSpace(lhs) != name {
			continue
		}
		rhs = strings.TrimSpace(rhs)
		if isQueryBuilderConstructor(rhs) {
			return true
		}
		if rhs != "(" {
			return false
		}
		for j := i + 1; j < len(lines); j++ {
			continued := strings.TrimSpace(lines[j])
			if continued == "" {
				continue
			}
			return isQueryBuilderConstructor(continued)
		}
		return false
	}
	return false
}

func isQueryBuilderConstructor(expression string) bool {
	return strings.HasPrefix(expression, "select(") || strings.HasPrefix(expression, "delete(") ||
		strings.HasPrefix(expression, "update(") || strings.HasPrefix(expression, "insert(")
}

// isDefHeaderCall is true when callStart names a def/async def parameter list
// (e.g. `def execute(...):`), not an execute/executemany call site.
func isDefHeaderCall(source string, callStart int) bool {
	if callStart <= 0 || callStart > len(source) {
		return false
	}
	lineStart := callStart
	for lineStart > 0 && source[lineStart-1] != '\n' {
		lineStart--
	}
	prefix := strings.TrimSpace(source[lineStart:callStart])
	return prefix == "def" || prefix == "async def"
}

// isStaticSQLArg reports a non-interpolated SQL text expression: a string
// literal, adjacent/implicitly concatenated string literals, or
// sqlalchemy.text("literal") / text("literal").
func isStaticSQLArg(expr string) bool {
	t := strings.TrimSpace(expr)
	if t == "" {
		return false
	}
	if isPureStringLiteral(t) || isAdjacentStringLiterals(t) {
		return true
	}
	return isStaticSQLAlchemyText(t)
}

// isAdjacentStringLiterals is true for Python's implicit string concatenation
// such as "SELECT … " "WHERE …" (no f-strings / formats).
func isAdjacentStringLiterals(expr string) bool {
	t := strings.TrimSpace(expr)
	if t == "" {
		return false
	}
	parts := 0
	for t != "" {
		rest, ok := consumeLeadingStringLiteral(t)
		if !ok {
			return false
		}
		parts++
		t = strings.TrimSpace(rest)
	}
	return parts >= 2
}

// consumeLeadingStringLiteral strips one non-f string literal from the front of
// expr and returns the remainder. f-strings are rejected (dynamic).
func consumeLeadingStringLiteral(expr string) (string, bool) {
	t := strings.TrimSpace(expr)
	if t == "" {
		return "", false
	}
	i := 0
	for i < len(t) {
		c := t[i]
		if c == 'f' || c == 'F' {
			return "", false
		}
		if c == 'b' || c == 'B' || c == 'r' || c == 'R' || c == 'u' || c == 'U' {
			i++
			continue
		}
		break
	}
	if i >= len(t) {
		return "", false
	}
	quote := t[i]
	if quote != '"' && quote != '\'' {
		return "", false
	}
	if i+2 < len(t) && t[i+1] == quote && t[i+2] == quote {
		end := strings.Index(t[i+3:], string([]byte{quote, quote, quote}))
		if end < 0 {
			return "", false
		}
		return t[i+3+end+3:], true
	}
	escape := false
	for j := i + 1; j < len(t); j++ {
		if escape {
			escape = false
			continue
		}
		if t[j] == '\\' {
			escape = true
			continue
		}
		if t[j] == quote {
			return t[j+1:], true
		}
	}
	return "", false
}

// isStaticSQLAlchemyText recognizes text("…") / sqlalchemy.text("…") whose
// first argument is a static (possibly adjacent) string literal.
func isStaticSQLAlchemyText(expr string) bool {
	t := strings.TrimSpace(expr)
	idx := strings.LastIndex(t, "text(")
	if idx < 0 {
		return false
	}
	if idx > 0 {
		prev := t[idx-1]
		if isIdentByteCWE(prev) {
			return false
		}
	}
	open := idx + len("text")
	if open >= len(t) || t[open] != '(' {
		return false
	}
	closeAt, argsText := scanCallArgs(t, open)
	if closeAt < 0 || strings.TrimSpace(t[closeAt+1:]) != "" {
		return false
	}
	args := splitTopLevelArgs(argsText)
	if len(args) == 0 {
		return false
	}
	inner := args[0]
	return (isPureStringLiteral(inner) || isAdjacentStringLiterals(inner)) && !looksSQLFormatted(inner)
}

func isIdentByteCWE(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

// sqlArgLooksInjected requires dynamic SQL evidence: formatted SQL text, a
// local assignment of dynamic SQL, or an expression that both looks dynamic
// and carries SQL/string content. Bare idents like HTTP `options` or custom
// `web_function` wrappers do not qualify.
func sqlArgLooksInjected(facts *PyCweFacts, source string, callStart int, sqlArg string) bool {
	t := strings.TrimSpace(sqlArg)
	if t == "" || isStaticSQLArg(t) {
		return false
	}
	if looksSQLFormatted(t) {
		return true
	}
	if isIdentOnly(t) {
		return localAssignLooksDynamicSQL(facts, source, callStart, t)
	}
	if !isDynamicExpr(t) {
		return false
	}
	return containsSQLKeyword(t) || strings.Contains(t, "\"") || strings.Contains(t, "'")
}

func localAssignLooksDynamicSQL(_ *PyCweFacts, source string, callStart int, name string) bool {
	if callStart < 0 || callStart > len(source) {
		return false
	}
	// Use unmasked source so f-string / formatted SQL text remains visible.
	lines := strings.Split(source[:callStart], "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		lhs, rhs, ok := strings.Cut(line, "=")
		if !ok || strings.TrimSpace(lhs) != name {
			continue
		}
		rhs = strings.TrimSpace(rhs)
		if rhs == "(" {
			for j := i + 1; j < len(lines); j++ {
				continued := strings.TrimSpace(lines[j])
				if continued == "" {
					continue
				}
				rhs = continued
				break
			}
		}
		if isQueryBuilderConstructor(rhs) || (isStaticSQLArg(rhs) && !looksSQLFormatted(rhs)) {
			return false
		}
		if looksSQLFormatted(rhs) {
			return true
		}
		if isDynamicExpr(rhs) && (containsSQLKeyword(rhs) || strings.Contains(rhs, "\"") || strings.Contains(rhs, "'")) {
			return true
		}
		return false
	}
	return false
}

func containsSQLKeyword(expr string) bool {
	upper := strings.ToUpper(expr)
	for _, kw := range []string{"SELECT", "INSERT", "UPDATE", "DELETE", "CREATE", "DROP", "ALTER", "WITH "} {
		if strings.Contains(upper, kw) {
			return true
		}
	}
	return false
}

// isExecuteCallbackPassthrough recognizes the DB-API execute-wrapper shape
// used by Django's execute_wrapper hook. The callback receives SQL and its
// bound arguments, then forwards them unchanged; this is not SQL construction.
// Keep the guard exact so a wrapper that rewrites sql before forwarding still
// reaches the dynamic-SQL check.
func isExecuteCallbackPassthrough(facts *PyCweFacts, source string, callStart int, args []string) bool {
	if len(args) != 4 || strings.TrimSpace(args[0]) != "sql" ||
		strings.TrimSpace(args[1]) != "params" || strings.TrimSpace(args[2]) != "many" ||
		strings.TrimSpace(args[3]) != "context" || callStart <= 0 || callStart > len(source) {
		return false
	}

	prefix := source[:callStart]
	lines := strings.Split(facts.codeMask(prefix, 0), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if !strings.HasPrefix(line, "def __call__(") && !strings.HasPrefix(line, "async def __call__(") {
			continue
		}
		for _, parameter := range []string{"execute", "sql", "params", "many", "context"} {
			if !strings.Contains(line, parameter) {
				return false
			}
		}
		for _, bodyLine := range lines[i+1:] {
			trimmed := strings.TrimSpace(bodyLine)
			if strings.HasPrefix(trimmed, "sql =") || strings.HasPrefix(trimmed, "sql=") ||
				strings.HasPrefix(trimmed, "sql +=") || strings.HasPrefix(trimmed, "sql-=") ||
				strings.HasPrefix(trimmed, "sql.format(") {
				return false
			}
		}
		return true
	}
	return false
}

func findExecuteCalls(facts *PyCweFacts, src string) []callSite {
	// Match .execute( / .executemany( and also bare execute( / executemany(
	out := findCalls(facts, src, ".execute", ".executemany")
	// bare names — findCalls with "execute" would match .execute too if boundary allows '.'
	// Our boundary treats '.' on left as failure for "execute", so bare works; but ".execute" already matched method form.
	// findCalls("execute") won't match ".execute" because left boundary sees '.'.
	out = append(out, findCalls(facts, src, "execute", "executemany")...)
	return out
}

func isIdentOnly(s string) bool {
	t := strings.TrimSpace(s)
	if t == "" {
		return false
	}
	for i, r := range t {
		if i == 0 {
			if r != '_' && (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
				return false
			}
			continue
		}
		if r != '_' && (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}

// --- CWE-22 Path traversal ---
//
// Classic traversal only: a dynamic segment joined into a restricted/base
// directory (os.path.join / Path(base) / segment) without confinement.
// Mere Path(__file__), Path(argv[0]) roots, or open(whole_caller_path) are
// intentionally outside this rule.

func detectCWE22(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	src := unit.Source

	// Documented safe suppressions
	if hasSafePathConfinement(src) {
		return
	}

	// open(os.path.join(root, user)) or open(Path(root) / user)
	for _, call := range findCalls(facts, src, "open") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) == 0 {
			continue
		}
		pathArg := args[0]
		if isPureStringLiteral(pathArg) || !pathArgLooksJoined(pathArg) {
			continue
		}
		if !openPathArgIsUnsafeJoin(pathArg) {
			continue
		}
		line, col := unit.LineCol(call.Start)
		rules.PushFindingWithConfidence(
			&MetaCWE22,
			unitFile(unit),
			line, col,
			"user-influenced path segment reaches open() without confinement (basename/resolve+prefix)",
			confidence70,
			out,
		)
		return
	}

	// Path(base) / dynamic_segment near a path sink (same-statement Path join)
	if !hasPathSink(src) {
		return
	}
	if start := pathlibUnsafeJoinStart(facts, src); start >= 0 {
		line, col := unit.LineCol(start)
		rules.PushFindingWithConfidence(&MetaCWE22, unitFile(unit), line, col,
			"pathlib path joined with dynamic segment without resolve+prefix confinement", confidence65, out)
	}
}

func hasPathSink(source string) bool {
	return strings.Contains(source, "open(") || strings.Contains(source, "read_text(") || strings.Contains(source, "write_text(") ||
		strings.Contains(source, "unlink(") || strings.Contains(source, "os.remove(")
}

// pathArgLooksJoined reports join/concat composition into a base path — not
// merely Path(whole) or a bare variable passed to open.
func pathArgLooksJoined(pathArg string) bool {
	if strings.Contains(pathArg, "os.path.join(") || strings.Contains(pathArg, " / ") {
		return true
	}
	if strings.Contains(pathArg, "+") &&
		(strings.Contains(pathArg, `"/"`) || strings.Contains(pathArg, `'/'`)) {
		return true
	}
	return false
}

func openPathArgIsUnsafeJoin(pathArg string) bool {
	if strings.Contains(pathArg, "os.path.join(") {
		return !joinAllLiterals(pathArg)
	}
	if strings.Contains(pathArg, " / ") {
		return pathExprHasDynamicPathDivision(pathArg)
	}
	return isDynamicExpr(pathArg)
}

func joinAllLiterals(pathArg string) bool {
	// Extract first os.path.join( interior and check args
	const needle = "os.path.join("
	i := strings.Index(pathArg, needle)
	if i < 0 {
		return false
	}
	open := i + len(needle) - 1
	closeAt, inner := scanCallArgs(pathArg, open)
	if closeAt < 0 {
		return false
	}
	parts := splitTopLevelArgs(inner)
	if len(parts) == 0 {
		return false
	}
	for _, p := range parts {
		if !isPureStringLiteral(p) {
			return false
		}
	}
	return true
}

// pathlibUnsafeJoinStart finds Path(...) / dynamic_segment on a code line.
// Arithmetic and comment ` / ` tokens are ignored (masked + Path-left requirement).
func pathlibUnsafeJoinStart(facts *PyCweFacts, src string) int {
	masked := src
	if facts != nil && facts.Masked != "" {
		masked = facts.Masked
	}
	start := 0
	for {
		idx := strings.Index(masked[start:], " / ")
		if idx < 0 {
			return -1
		}
		abs := start + idx
		if pathlibJoinLeftHasPathCtor(src, abs) {
			rhs := pathDivisionRHS(src, abs+3)
			if rhs != "" && !isPureStringLiteral(rhs) && !looksNumericOrArithmetic(rhs) {
				return abs
			}
		}
		start = abs + 3
	}
}

func pathlibJoinLeftHasPathCtor(src string, divAt int) bool {
	if divAt <= 0 || divAt > len(src) {
		return false
	}
	lineStart := strings.LastIndex(src[:divAt], "\n") + 1
	left := strings.TrimSpace(src[lineStart:divAt])
	return strings.Contains(left, "Path(") || strings.Contains(left, "pathlib.Path(")
}

func pathExprHasDynamicPathDivision(expr string) bool {
	start := 0
	for {
		idx := strings.Index(expr[start:], " / ")
		if idx < 0 {
			return false
		}
		abs := start + idx
		rhs := pathDivisionRHS(expr, abs+3)
		if rhs != "" && !isPureStringLiteral(rhs) && !looksNumericOrArithmetic(rhs) {
			return true
		}
		start = abs + 3
	}
}

func pathDivisionRHS(src string, from int) string {
	if from < 0 || from >= len(src) {
		return ""
	}
	end := from
	depth := 0
	for end < len(src) {
		c := src[end]
		if c == '\n' || c == '#' || c == ';' {
			break
		}
		if c == '(' || c == '[' {
			depth++
		}
		if c == ')' || c == ']' {
			if depth == 0 {
				break
			}
			depth--
		}
		if depth == 0 && (c == ',' || (c == ' ' && end+2 < len(src) && src[end:end+3] == " / ")) {
			break
		}
		end++
		if end-from > pathDivisionContextWindow {
			break
		}
	}
	return strings.TrimSpace(strings.TrimRight(src[from:end], ",)"))
}

func looksNumericOrArithmetic(expr string) bool {
	t := strings.TrimSpace(expr)
	if t == "" {
		return false
	}
	if t[0] >= '0' && t[0] <= '9' {
		return true
	}
	if strings.HasPrefix(t, "float(") || strings.HasPrefix(t, "int(") {
		return true
	}
	if strings.Contains(t, "<<") || strings.Contains(t, ">>") {
		return true
	}
	return false
}

// --- CWE-79 XSS ---

func detectCWE79(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	src := unit.Source

	// mark_safe(...) with dynamic arg
	for _, call := range findCalls(facts, src, "mark_safe") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) == 0 {
			continue
		}
		if isPureStringLiteral(args[0]) {
			continue
		}
		line, col := unit.LineCol(call.Start)
		rules.PushFindingWithConfidence(
			&MetaCWE79,
			unitFile(unit),
			line, col,
			"mark_safe on dynamic content disables autoescaping (XSS risk)",
			confidence80,
			out,
		)
		return
	}

	// Markup(...) dynamic
	for _, call := range findCalls(facts, src, "Markup") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) == 0 {
			continue
		}
		if isPureStringLiteral(args[0]) {
			continue
		}
		line, col := unit.LineCol(call.Start)
		rules.PushFindingWithConfidence(
			&MetaCWE79,
			unitFile(unit),
			line, col,
			"Markup() wraps dynamic HTML without escaping (XSS risk)",
			confidence80,
			out,
		)
		return
	}

	// render_template_string with dynamic template
	for _, call := range findCalls(facts, src, "render_template_string") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) == 0 {
			continue
		}
		if isPureStringLiteral(args[0]) && !isFStringExpr(args[0]) {
			// static template string with variables via kwargs is still XSS-prone if template uses |safe
			// but plan says: dynamic HTML; static template string is lower signal — skip pure literals
			continue
		}
		if isDynamicExpr(args[0]) || isFStringExpr(args[0]) {
			line, col := unit.LineCol(call.Start)
			rules.PushFindingWithConfidence(
				&MetaCWE79,
				unitFile(unit),
				line, col,
				"render_template_string with dynamic template/HTML (prefer render_template + autoescape)",
				confidence75,
				out,
			)
			return
		}
	}

	// |safe filter in template strings embedded in source
	if strings.Contains(src, "|safe") && (strings.Contains(src, "render_template") || strings.Contains(src, "{{")) {
		i := strings.Index(src, "|safe")
		line, col := unit.LineCol(i)
		rules.PushFindingWithConfidence(
			&MetaCWE79,
			unitFile(unit),
			line, col,
			"Jinja |safe filter disables escaping for dynamic content (XSS risk)",
			confidence65,
			out,
		)
		return
	}

	// Note: plain render_template("x.html", name=name) intentionally does NOT fire.
}
