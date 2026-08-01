package cwe

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	RegisterRule("CWE-434", detectCWE434, &MetaCWE434,
		"request.files")
	RegisterRule("CWE-427", detectCWE427, &MetaCWE427,
		"os.putenv", "os.environ", "LD_LIBRARY_PATH", "PYTHONPATH")
	RegisterRule("CWE-379", detectCWE379, &MetaCWE379,
		"open(", "/tmp/", "/var/tmp/")
	RegisterRule("CWE-459", detectCWE459, &MetaCWE459,
		"tempfile.mkstemp", "NamedTemporaryFile")
	RegisterRule("CWE-772", detectCWE772, &MetaCWE772,
		"open(", "socket.socket", "urlopen")
	RegisterRule("CWE-770", detectCWE770, &MetaCWE770,
		"request.get_data")
	RegisterRule("CWE-708", detectCWE708, &MetaCWE708,
		"os.chown")
	RegisterRule("CWE-477", detectCWE477, &MetaCWE477,
		"tempfile.mktemp", "cgi.escape", "asyncore.loop", "imp.load_module", "imp.load_source")
}

var (
	pyResourceAssignmentRE = regexp.MustCompile(`(?m)^\s*([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(?:open|socket\.socket|urllib\.request\.urlopen)\s*\(`)
	pyUploadAssignmentRE   = regexp.MustCompile(`(?m)^\s*([A-Za-z_][A-Za-z0-9_]*)\s*=\s*request\.files(?:\s*\[[^\n]+\]|\s*\.get\s*\()`)
)

// CWE-434 requires a Flask-shaped request.files assignment followed by saving
// that uploaded object without a same-function allowlist decision. Filename
// normalization alone is not an allowlist. Interprocedural upload helpers and
// content inspection deliberately remain outside this source-only rule.
func detectCWE434(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, fn := range facts.Functions() {
		masked := facts.codeMask(fn.body, fn.bodyStart)
		if !strings.Contains(masked, "request.files") {
			continue
		}
		search := 0
		uploads := map[string]struct{}{}
		for search <= len(masked) {
			match := pyUploadAssignmentRE.FindStringSubmatchIndex(masked[search:])
			if match == nil {
				break
			}
			uploads[masked[search+match[2]:search+match[3]]] = struct{}{}
			next := search + match[1]
			if next <= search {
				search++
			} else {
				search = next
			}
		}
		if len(uploads) == 0 || hasUploadTypeAllowlist(masked) {
			continue
		}
		for _, call := range findCallsMasked(fn.body, masked, ".save") {
			if _, ok := uploads[methodReceiverBefore(masked, call.Start)]; ok {
				emitResourceFinding(unit, &MetaCWE434, fn.bodyStart+call.Start, "uploaded file is saved without a same-function dangerous-type allowlist", confidence80, out)
				return
			}
		}
	}
}

func hasUploadTypeAllowlist(code string) bool {
	return strings.Contains(code, "if") && (strings.Contains(code, "allowed_file(") || strings.Contains(code, "allowed_extension("))
}

func methodReceiverBefore(source string, dot int) string {
	if dot <= 0 || dot > len(source) {
		return ""
	}
	end := dot
	start := end
	for start > 0 {
		c := source[start-1]
		if c != '_' && (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && (c < '0' || c > '9') {
			break
		}
		start--
	}
	return source[start:end]
}

// CWE-427 reports only direct process-library and Python-module path mutation
// from a non-literal expression. Fixed paths remain outside the rule.
func detectCWE427(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(facts, unit.Source, "os.putenv") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) >= 2 && searchPathName(args[0]) && isDynamicExpr(args[1]) {
			emitResourceFinding(unit, &MetaCWE427, call.Start, "externally controlled value is assigned to a process search-path environment variable", confidence84, out)
			return
		}
	}
	for _, line := range strings.Split(unit.Source, "\n") {
		compact := compactWhitespace(line)
		if strings.HasPrefix(strings.TrimSpace(line), "#") || !strings.Contains(compact, "os.environ[") || !strings.Contains(compact, "]=") {
			continue
		}
		if (strings.Contains(line, "LD_LIBRARY_PATH") || strings.Contains(line, "PYTHONPATH")) && !assignmentHasLiteralValue(line) {
			offset := strings.Index(unit.Source, line)
			emitResourceFinding(unit, &MetaCWE427, offset, "externally controlled value is assigned to a process search-path environment variable", confidence84, out)
			return
		}
	}
}

func searchPathName(expr string) bool {
	lower := strings.ToLower(expr)
	return strings.Contains(lower, "ld_library_path") || strings.Contains(lower, "pythonpath")
}

func assignmentHasLiteralValue(line string) bool {
	_, value, ok := strings.Cut(line, "=")
	return ok && isPureStringLiteral(value)
}

// CWE-379 is restricted to creating a fixed temporary pathname through the
// ordinary open APIs. Secure tempfile constructors are intentionally silent.
func detectCWE379(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"open", "os.open"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			args := splitTopLevelArgs(call.ArgsText)
			if len(args) > 0 && temporaryLiteralPath(args[0]) {
				emitResourceFinding(unit, &MetaCWE379, call.Start, "temporary file is created at a predictable pathname in a shared temporary directory", confidence82, out)
				return
			}
		}
	}
}

func temporaryLiteralPath(expr string) bool {
	t := strings.TrimSpace(expr)
	return isPureStringLiteral(t) && (strings.Contains(t, "/tmp/") || strings.Contains(t, "/var/tmp/"))
}

// CWE-459 is a same-function check for explicit persistent temporary files.
// It reports only when no unlink cleanup appears in that function.
func detectCWE459(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, fn := range facts.Functions() {
		code := facts.codeMask(fn.body, fn.bodyStart)
		if hasTemporaryCleanup(code) {
			continue
		}
		for _, call := range findCallsMasked(fn.body, code, "tempfile.mkstemp", "tempfile.NamedTemporaryFile") {
			if call.Name == "tempfile.NamedTemporaryFile" && !hasBooleanKwarg(call.ArgsText, "delete", "False") {
				continue
			}
			emitResourceFinding(unit, &MetaCWE459, fn.bodyStart+call.Start, "persistent temporary file has no same-function unlink cleanup", confidence78, out)
			return
		}
	}
}

func hasTemporaryCleanup(code string) bool {
	return strings.Contains(code, "os.unlink(") || strings.Contains(code, "os.remove(") || strings.Contains(code, ".unlink(")
}

// CWE-772 reports an assigned file, socket, or urlopen response only when the
// variable is not closed in the same function. Context-manager forms do not
// match the assignment shape.
func detectCWE772(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, fn := range facts.Functions() {
		code := facts.codeMask(fn.body, fn.bodyStart)
		if !strings.Contains(code, "open(") && !strings.Contains(code, "socket.socket") && !strings.Contains(code, "urlopen") {
			continue
		}
		search := 0
		for search <= len(code) {
			match := pyResourceAssignmentRE.FindStringSubmatchIndex(code[search:])
			if match == nil {
				break
			}
			abs0 := search + match[0]
			abs1 := search + match[1]
			name := code[search+match[2] : search+match[3]]
			if !strings.Contains(code[abs1:], name+".close(") {
				emitResourceFinding(unit, &MetaCWE772, fn.bodyStart+abs0, "resource is assigned without a same-function close or context-manager release", confidence76, out)
				return
			}
			if abs1 <= search {
				search++
			} else {
				search = abs1
			}
		}
	}
}

// CWE-770 identifies Flask's direct unbounded request-body reader only when
// the module has no MAX_CONTENT_LENGTH configuration evidence.
func detectCWE770(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil || strings.Contains(unit.Source, "MAX_CONTENT_LENGTH") {
		return
	}
	if calls := findCalls(facts, unit.Source, "request.get_data"); len(calls) > 0 {
		emitResourceFinding(unit, &MetaCWE770, calls[0].Start, "request body is read without a module-level maximum content-length limit", confidence78, out)
	}
}

// CWE-708 restricts the ownership rule to an explicit root owner or group,
// which is a verifiable high-risk assignment without policy inference.
func detectCWE708(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(facts, unit.Source, "os.chown") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) >= 3 && (strings.TrimSpace(args[1]) == "0" || strings.TrimSpace(args[2]) == "0") {
			emitResourceFinding(unit, &MetaCWE708, call.Start, "resource ownership is explicitly assigned to the root user or group", confidence82, out)
			return
		}
	}
}

// CWE-477 reports calls to removed or deprecated standard-library APIs. It
// does not flag imports alone because compatibility wrappers may import them.
func detectCWE477(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"tempfile.mktemp", "cgi.escape", "asyncore.loop", "imp.load_module", "imp.load_source"} {
		if calls := findCalls(facts, unit.Source, name); len(calls) > 0 {
			emitResourceFinding(unit, &MetaCWE477, calls[0].Start, "deprecated or obsolete Python standard-library function is called", confidence82, out)
			return
		}
	}
}

func emitResourceFinding(unit *core.ParsedUnit, meta *rules.RuleMetadata, offset int, message string, confidence float32, out *[]rules.Finding) {
	line, col := unit.LineCol(offset)
	rules.PushFindingWithConfidence(meta, unitFile(unit), line, col, message, confidence, out)
}
