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

func detectCWE502(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	src := unit.Source

	// pickle sinks — any call is treated as unsafe deserialization of potentially
	// untrusted data (conservative; matches priority-batch museum style).
	for _, call := range findCalls(src, "pickle.loads", "pickle.load", "pickle.Unpickler") {
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
	for _, call := range findCalls(src, "yaml.unsafe_load") {
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
	for _, call := range findCalls(src, "yaml.load") {
		// Do not treat yaml.safe_load as yaml.load (boundary check already prevents prefix match of safe_load when searching "yaml.load" — verify)
		// "yaml.load" is not a prefix of "yaml.safe_load" in reverse; but "yaml.load" could match inside nothing else.
		if yamlLoadLooksSafe(call.ArgsText) {
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

func detectCWE78(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	src := unit.Source

	// os.system / os.popen with dynamic command
	for _, name := range []string{"os.system", "os.popen"} {
		for _, call := range findCalls(src, name) {
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
		for _, call := range findCalls(src, name) {
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

func detectCWE89(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	src := unit.Source

	// Prefer method-style .execute( / .executemany( and bare execute(
	// Scan for ".execute(" and "execute(" carefully.
	for _, call := range findExecuteCalls(src) {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) == 0 {
			continue
		}
		sqlArg := args[0]
		// Parameterized: first arg is pure string with placeholders AND second arg is present (bound params)
		if len(args) >= 2 && isPureStringLiteral(sqlArg) && !looksSQLFormatted(sqlArg) {
			// bound args tuple/list → safe
			second := strings.TrimSpace(args[1])
			if strings.HasPrefix(second, "(") || strings.HasPrefix(second, "[") || isIdentOnly(second) {
				continue
			}
		}
		// Dynamic SQL construction in first arg
		if looksSQLFormatted(sqlArg) || (isDynamicExpr(sqlArg) && !isPureStringLiteral(sqlArg)) {
			// Avoid flagging completely static pure string without format
			if isPureStringLiteral(sqlArg) {
				continue
			}
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

func findExecuteCalls(src string) []callSite {
	// Match .execute( / .executemany( and also bare execute( / executemany(
	out := findCalls(src, ".execute", ".executemany")
	// bare names — findCalls with "execute" would match .execute too if boundary allows '.'
	// Our boundary treats '.' on left as failure for "execute", so bare works; but ".execute" already matched method form.
	// findCalls("execute") won't match ".execute" because left boundary sees '.'.
	out = append(out, findCalls(src, "execute", "executemany")...)
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

func detectCWE22(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	src := unit.Source

	// Documented safe suppressions
	if hasSafePathConfinement(src) {
		return
	}

	// open(os.path.join(...)) style — join of non-all-literals into open
	for _, call := range findCalls(src, "open") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) == 0 {
			continue
		}
		pathArg := args[0]
		if isPureStringLiteral(pathArg) {
			continue
		}
		// open(os.path.join(root, user)) or open(Path(...) / user)
		if strings.Contains(pathArg, "os.path.join(") ||
			strings.Contains(pathArg, "pathlib.Path(") ||
			strings.Contains(pathArg, "Path(") ||
			strings.Contains(pathArg, "/") {
			// dynamic path composition
			if isDynamicExpr(pathArg) || strings.Contains(pathArg, "os.path.join(") {
				// if join args are all pure literals, skip
				if joinAllLiterals(pathArg) {
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
		}
		// open(variable) with join elsewhere and open of that var is harder; flag open(os.path.join) primarily
	}

	// Path(root) / user then used — flag join-style Path division with dynamic right-hand side near open/read
	// Simple pattern: ` / ` with request-like or bare identifiers after Path(
	if !strings.Contains(src, "Path(") || !strings.Contains(src, " / ") || !hasPathSink(src) || !hasDynamicPathDiv(src) {
		return
	}
	i := strings.Index(src, "Path(")
	if i < 0 {
		i = strings.Index(src, " / ")
	}
	line, col := unit.LineCol(i)
	rules.PushFindingWithConfidence(&MetaCWE22, unitFile(unit), line, col,
		"pathlib path joined with dynamic segment without resolve+prefix confinement", confidence65, out)
}

func hasPathSink(source string) bool {
	return strings.Contains(source, "open(") || strings.Contains(source, "read_text(") || strings.Contains(source, "write_text(") ||
		strings.Contains(source, "unlink(") || strings.Contains(source, "os.remove(")
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

func hasDynamicPathDiv(src string) bool {
	// Find " / " at top-ish level that is not pure string on both sides — lightweight
	start := 0
	for {
		idx := strings.Index(src[start:], " / ")
		if idx < 0 {
			return false
		}
		abs := start + idx
		// take a short token/expression on the right-hand side
		end := abs + 3
		for end < len(src) && src[end] != '\n' && src[end] != ')' && src[end] != ';' && src[end] != '#' {
			end++
			if end-abs > pathDivisionContextWindow {
				break
			}
		}
		rhsExpr := strings.TrimSpace(src[abs+3 : end])
		// strip trailing call punctuation
		rhsExpr = strings.TrimRight(rhsExpr, ",)")
		if rhsExpr != "" && !isPureStringLiteral(rhsExpr) {
			return true
		}
		start = abs + 3
	}
}

// --- CWE-79 XSS ---

func detectCWE79(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	src := unit.Source

	// mark_safe(...) with dynamic arg
	for _, call := range findCalls(src, "mark_safe") {
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
	for _, call := range findCalls(src, "Markup") {
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
	for _, call := range findCalls(src, "render_template_string") {
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
