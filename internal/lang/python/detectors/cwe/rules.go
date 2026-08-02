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
		// F-string SQL whose interpolations are only bound-placeholder calls
		// (store.param()) or only allowlist-guarded identifiers carries no
		// dynamic user content; values travel through the bind argument.
		if looksSQLFormatted(sqlArg) {
			if sqlFStringIsParamsOnly(facts, src, call.Start, sqlArg) || sqlFStringIdentsAllowlistGuarded(facts, src, call.Start, sqlArg) {
				continue
			}
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
	// sqlalchemy.text(<expr>): a bound static text() is safe, but a dynamic
	// inner expression is exactly the injection case (e.g.
	// db.execute(text(sql)) with sql built from f-strings — WeThePeople TPs).
	if inner, ok := sqlalchemyTextInner(t); ok {
		inner = strings.TrimSpace(inner)
		if isStaticSQLArg(inner) && !looksSQLFormatted(inner) {
			return false
		}
		if isIdentOnly(inner) {
			// Schema DDL built purely from SQLAlchemy model metadata (loop
			// variables over .columns/.tables) is developer-defined, not data.
			if sqlTextInnerMetadataOnly(facts, source, callStart, inner) {
				return false
			}
			return localAssignLooksDynamicSQL(facts, source, callStart, inner)
		}
		if looksSQLFormatted(inner) {
			return true
		}
		return isDynamicExpr(inner) && (containsSQLKeyword(inner) || strings.Contains(inner, "\"") || strings.Contains(inner, "'"))
	}
	// Dotted module/class constants (SQL.CREATE_TEST_TABLE) are not dynamic:
	// the constant name itself is not SQL evidence; only its assigned value is.
	if isDottedConstantRef(t) {
		return dottedConstantResolvesDynamic(facts, source, callStart, t)
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
		if rhs == "(" || rhs == `"""` || rhs == `'''` {
			var b strings.Builder
			b.WriteString(rhs)
			if rhs == "(" {
				depth := 1
				for j := i + 1; j < len(lines) && depth > 0; j++ {
					continued := strings.TrimSpace(lines[j])
					if continued == "" {
						continue
					}
					depth += pyOpenBracketBalance(continued)
					b.WriteByte(' ')
					b.WriteString(continued)
				}
			} else {
				for j := i + 1; j < len(lines); j++ {
					continued := strings.TrimSpace(lines[j])
					b.WriteString(continued)
					if strings.Contains(continued, rhs) {
						break
					}
				}
			}
			rhs = b.String()
		}
		if isQueryBuilderConstructor(rhs) || (isStaticSQLArg(rhs) && !looksSQLFormatted(rhs)) {
			return false
		}
		if looksSQLFormatted(rhs) {
			return true
		}
		// Ternary / subscripts over module-level SQL constants
		// (SQLITE_SCHEMA[0] if ... else POSTGRES_SCHEMA[0]) are static.
		if sqlExprIsStaticConstantRef(source, callStart, rhs) {
			return false
		}
		if isDynamicExpr(rhs) && (containsSQLKeyword(rhs) || strings.Contains(rhs, "\"") || strings.Contains(rhs, "'")) {
			return true
		}
		return false
	}
	return false
}

// sqlExprIsStaticConstantRef reports an expression that selects from
// module-level SQL constants without interpolation: a ternary over constant
// references, a subscript of a literal collection, or a plain literal.
func sqlExprIsStaticConstantRef(source string, callStart int, expr string) bool {
	t := strings.TrimSpace(expr)
	if t == "" {
		return false
	}
	if isPureStringLiteral(t) || isAdjacentStringLiterals(t) {
		return true
	}
	if a, _, b, ok := sqlTernaryParts(t); ok {
		return sqlExprIsStaticConstantRef(source, callStart, a) && sqlExprIsStaticConstantRef(source, callStart, b)
	}
	if base, key, isGet, ok := splitIndexOrGet(t); ok && !isGet {
		if !isNumericLiteral(key) && !isPureStringLiteral(key) {
			return false
		}
		return sqlIdentResolvesStaticCollection(source, callStart, base)
	}
	return false
}

// sqlIdentResolvesStaticCollection reports that ident is assigned a literal
// collection (tuple/list/dict of string literals) or a generator over one.
func sqlIdentResolvesStaticCollection(source string, callStart int, ident string) bool {
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
			if strings.HasPrefix(trimmed, "class ") || defParamsContain(trimmed, ident) {
				return false
			}
			continue
		}
		lhs, rhs, ok := strings.Cut(trimmed, "=")
		if !ok {
			continue
		}
		lhs = strings.TrimSpace(lhs)
		if lhs != ident {
			// tolerate `name: Type = value` annotations
			if colon := strings.IndexByte(lhs, ':'); colon < 0 || strings.TrimSpace(lhs[:colon]) != ident {
				continue
			}
		}
		rhs = stripPyComment(strings.TrimSpace(rhs))
		if rhs == "" {
			return false
		}
		if rhs == "(" || rhs == "[" || rhs == "{" || strings.HasPrefix(rhs, "tuple(") {
			var b strings.Builder
			depth := pyOpenBracketBalance(rhs)
			b.WriteString(rhs)
			for j := i + 1; j < len(lines) && depth > 0; j++ {
				next := strings.TrimSpace(lines[j])
				if next == "" {
					continue
				}
				depth += pyOpenBracketBalance(next)
				b.WriteByte(' ')
				b.WriteString(next)
			}
			rhs = b.String()
		}
		if looksStaticStringList(rhs) || looksStaticDictLiteral(rhs) {
			return true
		}
		// generator expression: tuple(x for x in <static base>)
		if forIdx := strings.Index(rhs, " for "); forIdx >= 0 {
			if inIdx := strings.Index(rhs[forIdx+5:], " in "); inIdx >= 0 {
				iterable := strings.TrimSpace(rhs[forIdx+5+inIdx+4:])
				if close := strings.LastIndexByte(iterable, ')'); close >= 0 {
					iterable = strings.TrimSpace(iterable[:close])
				}
				if isIdentOnly(iterable) {
					return sqlIdentResolvesStaticCollection(source, callStart, iterable)
				}
			}
		}
		return false
	}
	return false
}

// sqlTernaryParts splits a `A if C else B` expression into its three parts.
func sqlTernaryParts(expr string) (a, cond, b string, ok bool) {
	t := strings.TrimSpace(expr)
	ifIdx := strings.Index(t, " if ")
	if ifIdx < 0 {
		return "", "", "", false
	}
	elseIdx := strings.Index(t[ifIdx+4:], " else ")
	if elseIdx < 0 {
		return "", "", "", false
	}
	elseIdx += ifIdx + 4
	return strings.TrimSpace(t[:ifIdx]), strings.TrimSpace(t[ifIdx+4 : elseIdx]), strings.TrimSpace(t[elseIdx+6:]), true
}

// sqlalchemyTextInner returns the first-argument expression of a
// sqlalchemy.text(...) call when expr is exactly that call.
func sqlalchemyTextInner(expr string) (string, bool) {
	t := strings.TrimSpace(expr)
	idx := strings.LastIndex(t, "text(")
	if idx < 0 {
		return "", false
	}
	if idx > 0 {
		if isIdentByteCWE(t[idx-1]) {
			return "", false
		}
	}
	open := idx + len("text")
	if open >= len(t) || t[open] != '(' {
		return "", false
	}
	closeAt, argsText := scanCallArgs(t, open)
	if closeAt < 0 || strings.TrimSpace(t[closeAt+1:]) != "" {
		return "", false
	}
	args := splitTopLevelArgs(argsText)
	if len(args) == 0 {
		return "", false
	}
	return args[0], true
}

// isDottedConstantRef reports an identifier attribute chain (SQL.CREATE_TEST
// _TABLE) with no operators, quotes, or interpolation.
func isDottedConstantRef(expr string) bool {
	t := strings.TrimSpace(expr)
	if len(t) < 3 || !strings.Contains(t, ".") {
		return false
	}
	expectIdent := true
	first := true
	for i := 0; i < len(t); i++ {
		c := t[i]
		if c == '.' {
			if expectIdent || first {
				return false
			}
			expectIdent = true
			continue
		}
		if !isIdentByteCWE(c) {
			return false
		}
		if first && c >= '0' && c <= '9' {
			return false
		}
		first = false
		expectIdent = false
	}
	return !expectIdent
}

// dottedConstantResolvesDynamic resolves the root of a dotted constant
// reference (SQL.CREATE_TEST_TABLE → SQL) to the assignment preceding the
// call. Only a formatted/dynamic assignment makes the reference dynamic.
func dottedConstantResolvesDynamic(facts *PyCweFacts, source string, callStart int, expr string) bool {
	root := expr
	if i := strings.IndexByte(expr, '.'); i > 0 {
		root = expr[:i]
	}
	if callStart <= 0 || callStart > len(source) {
		return false
	}
	lines := strings.Split(source[:callStart], "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if strings.HasPrefix(trimmed, "class "+root+"(") || strings.HasPrefix(trimmed, "class "+root+" ") ||
			trimmed == "class "+root {
			return false
		}
		lhs, rhs, ok := strings.Cut(trimmed, "=")
		if !ok || strings.TrimSpace(lhs) != root {
			continue
		}
		rhs = strings.TrimSpace(rhs)
		if rhs == "" {
			return false
		}
		if isStaticSQLArg(rhs) && !looksSQLFormatted(rhs) {
			return false
		}
		return looksSQLFormatted(rhs) ||
			(isDynamicExpr(rhs) && (containsSQLKeyword(rhs) || strings.Contains(rhs, "\"") || strings.Contains(rhs, "'")))
	}
	return false
}

// sqlFStringIsParamsOnly reports an f-string SQL literal whose every
// interpolation is a bound-placeholder token: a .param() call directly or a
// local alias assigned from one (placeholder = self.store.param()). All data
// travels through the second execute argument, so the SQL text is static.
func sqlFStringIsParamsOnly(facts *PyCweFacts, source string, callStart int, expr string) bool {
	parts := fstringInterpolations(expr)
	if len(parts) == 0 {
		return false
	}
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t == "" {
			continue
		}
		if isParamPlaceholderCall(t) {
			continue
		}
		if !isIdentOnly(t) || !identAssignedFromParamCall(source, callStart, t, map[string]bool{}) {
			return false
		}
	}
	return true
}

// isParamPlaceholderCall reports a `<obj>.param(...)` call expression.
func isParamPlaceholderCall(t string) bool {
	idx := strings.LastIndex(t, ".param")
	if idx < 0 {
		return false
	}
	open := idx + len(".param")
	if open >= len(t) || t[open] != '(' {
		return false
	}
	closeAt, _ := scanCallArgs(t, open)
	return closeAt >= 0 && strings.TrimSpace(t[closeAt+1:]) == ""
}

// identAssignedFromParamCall resolves ident through local assignments to a
// .param() placeholder call (placeholder = self.store.param()).
func identAssignedFromParamCall(source string, callStart int, ident string, seen map[string]bool) bool {
	if seen[ident] {
		return false
	}
	seen[ident] = true
	if callStart <= 0 || callStart > len(source) {
		return false
	}
	lines := strings.Split(source[:callStart], "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if strings.HasPrefix(trimmed, "def ") || strings.HasPrefix(trimmed, "async def ") {
			// A function parameter is runtime-controlled, not a placeholder.
			return false
		}
		lhs, rhs, ok := strings.Cut(trimmed, "=")
		if !ok {
			continue
		}
		if strings.TrimSpace(lhs) != ident {
			continue
		}
		rhs = stripPyComment(strings.TrimSpace(rhs))
		if isParamPlaceholderCall(rhs) {
			return true
		}
		if isIdentOnly(rhs) {
			return identAssignedFromParamCall(source, callStart, rhs, seen)
		}
		return false
	}
	return false
}

// sqlFStringIdentsAllowlistGuarded reports an f-string SQL literal whose
// interpolated identifiers are all allowlist-constrained earlier in the same
// function: `if <ident> not in <set>: raise ...` (violit audited FP). Bound
// values travel through the parameter list, so the interpolated token cannot
// be attacker-controlled.
func sqlFStringIdentsAllowlistGuarded(facts *PyCweFacts, source string, callStart int, expr string) bool {
	if callStart <= 0 || callStart > len(source) {
		return false
	}
	var idents []string
	for _, p := range fstringInterpolations(expr) {
		t := strings.TrimSpace(p)
		if t == "" {
			continue
		}
		if !isIdentOnly(t) {
			return false
		}
		idents = append(idents, t)
	}
	if len(idents) == 0 {
		return false
	}
	lines := strings.Split(source[:callStart], "\n")
	callLine := len(lines) - 1
	funcIdx := -1
	for i := len(lines) - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "def ") || strings.HasPrefix(trimmed, "async def ") {
			funcIdx = i
			break
		}
	}
	if funcIdx < 0 {
		return false
	}
	for _, id := range idents {
		if !identAllowlistGuarded(lines, funcIdx, callLine, id) {
			return false
		}
	}
	return true
}

func identAllowlistGuarded(lines []string, from, to int, ident string) bool {
	for i := from; i < to; i++ {
		t := strings.TrimSpace(lines[i])
		if t == "" || !strings.HasPrefix(t, "if ") {
			continue
		}
		notIdx := strings.Index(t, " not in ")
		if notIdx < 0 {
			continue
		}
		lhs := strings.TrimSpace(t[len("if "):notIdx])
		if lhs != ident {
			continue
		}
		rest := strings.TrimSpace(t[notIdx+len(" not in "):])
		if strings.HasPrefix(rest, ": raise") || strings.HasPrefix(rest, ":raise") {
			return true
		}
		if !strings.HasSuffix(rest, ":") {
			continue
		}
		ifIndent := indentOfCWE(lines[i])
		for j := i + 1; j < to; j++ {
			if strings.TrimSpace(lines[j]) == "" {
				continue
			}
			if indentOfCWE(lines[j]) <= ifIndent {
				break
			}
			body := strings.TrimSpace(lines[j])
			if body == "raise" || strings.HasPrefix(body, "raise ") {
				return true
			}
		}
	}
	return false
}

func indentOfCWE(line string) int {
	n := 0
	for n < len(line) {
		if line[n] != ' ' && line[n] != '\t' {
			break
		}
		n++
	}
	return n
}

// sqlTextInnerMetadataOnly reports that the identifier was assigned a
// formatted SQL string whose interpolations all derive from SQLAlchemy model
// metadata: loop variables over <x>.columns / <x>.tables (and values assigned
// from attribute or method expressions on them). Such DDL is developer-defined
// schema, not injected data (violit db.py:151 audited FP).
func sqlTextInnerMetadataOnly(facts *PyCweFacts, source string, callStart int, ident string) bool {
	if callStart <= 0 || callStart > len(source) {
		return false
	}
	lines := strings.Split(source[:callStart], "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if strings.HasPrefix(trimmed, "def ") || strings.HasPrefix(trimmed, "async def ") {
			return false
		}
		lhs, rhs, ok := strings.Cut(trimmed, "=")
		if !ok || strings.TrimSpace(lhs) != ident {
			continue
		}
		rhs = strings.TrimSpace(rhs)
		if rhs == "" {
			return false
		}
		if rhs == "(" {
			// parenthesized adjacent f-string concatenation
			var b strings.Builder
			depth := 1
			for j := i + 1; j < len(lines) && depth > 0; j++ {
				next := strings.TrimSpace(lines[j])
				if next == "" {
					continue
				}
				depth += pyOpenBracketBalance(next)
				b.WriteByte(' ')
				b.WriteString(next)
			}
			rhs = b.String()
		}
		if !looksSQLFormatted(rhs) {
			return false
		}
		parts := allFstringInterpolations(rhs)
		if len(parts) == 0 {
			return false
		}
		for _, part := range parts {
			if strings.TrimSpace(part) != "" && !sqlInterpMetadataDerived(source, callStart, part) {
				return false
			}
		}
		return true
	}
	return false
}

// sqlInterpMetadataDerived reports whether a single f-string interpolation is
// derived from SQLAlchemy model metadata: an attribute chain rooted at a
// metadata loop variable, a method call on such a root (or with one as an
// argument), or an identifier assigned from such an expression.
func sqlInterpMetadataDerived(source string, callStart int, part string) bool {
	t := strings.TrimSpace(part)
	if t == "" {
		return true
	}
	if strings.HasSuffix(t, ")") {
		if open := strings.IndexByte(t, '('); open > 0 {
			callee := strings.TrimSpace(t[:open])
			if isDottedConstantRef(callee) {
				root := callee[:strings.IndexByte(callee, '.')]
				if metadataLoopVariable(source, callStart, root) {
					return true
				}
			}
			args := splitTopLevelArgs(t[open+1 : len(t)-1])
			for _, a := range args {
				if sqlInterpMetadataDerived(source, callStart, a) {
					return true
				}
			}
			return false
		}
		return false
	}
	if isDottedConstantRef(t) {
		root := t[:strings.IndexByte(t, '.')]
		return metadataLoopVariable(source, callStart, root)
	}
	if isIdentOnly(t) {
		if metadataLoopVariable(source, callStart, t) {
			return true
		}
		return identAssignedFromMetadata(source, callStart, t, map[string]bool{})
	}
	return false
}

// metadataLoopVariable reports that root is a loop target of a for loop over a
// SQLAlchemy metadata collection (<x>.columns / <x>.tables.items()) or over a
// same-file dict derived from schema inspection (inspector.get_columns(...)).
func metadataLoopVariable(source string, callStart int, root string) bool {
	for _, line := range strings.Split(source, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "for ") {
			continue
		}
		body := strings.TrimPrefix(trimmed, "for ")
		inIdx := strings.Index(body, " in ")
		if inIdx < 0 {
			continue
		}
		targets := strings.TrimSpace(body[:inIdx])
		collection := strings.TrimSpace(body[inIdx+4:])
		collection = strings.TrimSpace(strings.TrimSuffix(collection, ":"))
		if !loopTargetMatches(targets, root) {
			continue
		}
		if strings.HasSuffix(collection, ".columns") ||
			strings.Contains(collection, ".tables.items()") ||
			strings.Contains(collection, ".metadata.tables.items()") {
			return true
		}
		// unwrap list(...) wrappers: for x in list(dict.keys())
		if strings.HasPrefix(collection, "list(") && strings.HasSuffix(collection, ")") {
			collection = strings.TrimSpace(collection[len("list(") : len(collection)-1])
		}
		for _, method := range []string{".items()", ".keys()", ".values()"} {
			if strings.HasSuffix(collection, method) {
				base := strings.TrimSpace(collection[:len(collection)-len(method)])
				if isIdentOnly(base) && identAssignedFromInspector(source, base) {
					return true
				}
			}
		}
	}
	return false
}

// identAssignedFromInspector reports that ident is assigned from a schema
// inspection call (inspector.get_columns / get_table_names / get_foreign_keys).
func identAssignedFromInspector(source, ident string) bool {
	lines := strings.Split(source, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		lhs, rhs, ok := strings.Cut(trimmed, "=")
		if !ok {
			continue
		}
		lhs = strings.TrimSpace(lhs)
		if lhs != ident {
			// tolerate `name: Type = value` annotations
			if colon := strings.IndexByte(lhs, ':'); colon < 0 || strings.TrimSpace(lhs[:colon]) != ident {
				continue
			}
		}
		rhs = strings.TrimSpace(rhs)
		if rhs == "" {
			continue
		}
		// multi-line dict comprehensions
		if rhs == "{" || rhs == "[" || rhs == "(" {
			var b strings.Builder
			depth := pyOpenBracketBalance(rhs)
			b.WriteString(rhs)
			for j := i + 1; j < len(lines) && depth > 0; j++ {
				next := strings.TrimSpace(lines[j])
				if next == "" {
					continue
				}
				depth += pyOpenBracketBalance(next)
				b.WriteByte(' ')
				b.WriteString(next)
			}
			rhs = b.String()
		}
		return strings.Contains(rhs, "get_columns(") ||
			strings.Contains(rhs, "get_table_names(") ||
			strings.Contains(rhs, "get_foreign_keys(")
	}
	return false
}

// identAssignedFromMetadata resolves ident through local assignments to a
// metadata-derived expression.
func identAssignedFromMetadata(source string, callStart int, ident string, seen map[string]bool) bool {
	if seen[ident] {
		return false
	}
	seen[ident] = true
	if callStart <= 0 || callStart > len(source) {
		return false
	}
	lines := strings.Split(source[:callStart], "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if strings.HasPrefix(trimmed, "def ") || strings.HasPrefix(trimmed, "async def ") {
			return false
		}
		lhs, rhs, ok := strings.Cut(trimmed, "=")
		if !ok {
			continue
		}
		if strings.TrimSpace(lhs) != ident {
			continue
		}
		rhs = stripPyComment(strings.TrimSpace(rhs))
		if rhs == "" {
			return false
		}
		// Multi-line call/dict right-hand sides.
		depth := pyOpenBracketBalance(rhs)
		if depth > 0 {
			var b strings.Builder
			b.WriteString(rhs)
			for j := i + 1; j < len(lines) && depth > 0; j++ {
				next := strings.TrimSpace(lines[j])
				if next == "" {
					continue
				}
				depth += pyOpenBracketBalance(next)
				b.WriteByte(' ')
				b.WriteString(next)
			}
			rhs = b.String()
		}
		return sqlInterpMetadataDerived(source, callStart, rhs)
	}
	return false
}

// allFstringInterpolations collects {…} bodies from every f-string literal in
// expr, including parenthesized adjacent f-string concatenations.
func allFstringInterpolations(expr string) []string {
	t := strings.TrimSpace(expr)
	if strings.HasPrefix(t, "(") && strings.HasSuffix(t, ")") {
		t = strings.TrimSpace(t[1 : len(t)-1])
	}
	var out []string
	for {
		body, rest, ok := consumeFStringLiteral(t)
		if !ok {
			break
		}
		out = append(out, collectFstringParts(body)...)
		t = strings.TrimSpace(rest)
		if t == "" {
			break
		}
	}
	return out
}

// consumeFStringLiteral strips one leading f-string literal from expr and
// returns its raw interior and the remainder.
func consumeFStringLiteral(expr string) (body, rest string, ok bool) {
	t := strings.TrimSpace(expr)
	i := 0
	hasF := false
	for i < len(t) {
		c := t[i]
		if c == 'f' || c == 'F' {
			hasF = true
			i++
			continue
		}
		if c == 'r' || c == 'R' || c == 'b' || c == 'B' || c == 'u' || c == 'U' {
			i++
			continue
		}
		break
	}
	if !hasF || i >= len(t) {
		return "", "", false
	}
	quote := t[i]
	if quote != '"' && quote != '\'' {
		return "", "", false
	}
	if i+2 < len(t) && t[i+1] == quote && t[i+2] == quote {
		after := t[i+3:]
		end := strings.Index(after, string([]byte{quote, quote, quote}))
		if end < 0 {
			return "", "", false
		}
		return after[:end], t[i+3+end+3:], true
	}
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
			return t[i+1 : j], t[j+1:], true
		}
	}
	return "", "", false
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
