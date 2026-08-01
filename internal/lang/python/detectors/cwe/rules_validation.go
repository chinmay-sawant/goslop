package cwe

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	// Most validation rules stay ungated (varied framework spellings).
	// CWE-1230 only needs Content-Disposition plus a request token.
	RegisterRule("CWE-1173", detectCWE1173, &MetaCWE1173)
	RegisterRule("CWE-1230", detectCWE1230, &MetaCWE1230,
		"Content-Disposition")
	RegisterRule("CWE-1236", detectCWE1236, &MetaCWE1236)
	RegisterRule("CWE-1286", detectCWE1286, &MetaCWE1286)
	RegisterRule("CWE-1289", detectCWE1289, &MetaCWE1289)
	RegisterRule("CWE-1333", detectCWE1333, &MetaCWE1333)
	RegisterRule("CWE-1389", detectCWE1389, &MetaCWE1389)
	RegisterRule("CWE-140", detectCWE140, &MetaCWE140)
}

var (
	pyContentDispositionRequestRE = regexp.MustCompile(`(?is)(?:headers\s*\[\s*['"]Content-Disposition['"]\s*\]|headers\s*\.\s*(?:add|set)\s*\(\s*['"]Content-Disposition['"])[^\n]*request\s*\.(?:args|form|values|GET|POST|query_params)\b`)
	pyPathDenyListRE              = regexp.MustCompile(`(?im)^\s*if\s+([A-Za-z_][A-Za-z0-9_]*)\s*==\s*['"][^'"]*(?:private|secret|admin)[^'"]*['"]\s*:`)
	pyNestedRegexQuantifierRE     = regexp.MustCompile(`\([^\n()]*[+*][^\n()]*\)[+*]`)
)

// CWE-1173 reports a deliberately narrow unsafe construction: a route reads
// JSON and persists it through a model save without a visible schema,
// serializer, or validation call in that handler. Simple JSON readers and
// handlers that validate elsewhere are intentionally outside this heuristic.
func detectCWE1173(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, handler := range routeHandlerBodies(facts, unit.Source) {
		code := strings.ToLower(facts.codeMask(handler.body, handler.start))
		if !hasRequestBodyRead(code) || !strings.Contains(code, ".save(") || hasValidationFrameworkEvidence(code) {
			continue
		}
		if start := strings.Index(code, "request.get_json("); start >= 0 {
			emitValidationFinding(unit, &MetaCWE1173, handler.start+start, "route persists JSON request data without an observable schema or serializer validation step", confidence76, out)
			return
		}
		if start := strings.Index(code, "request.json"); start >= 0 {
			emitValidationFinding(unit, &MetaCWE1173, handler.start+start, "route persists JSON request data without an observable schema or serializer validation step", confidence76, out)
			return
		}
	}
}

func hasValidationFrameworkEvidence(code string) bool {
	for _, marker := range []string{".is_valid(", ".validated_data", ".model_validate(", ".parse_obj(", ".load(", "schema.validate(", "serializer("} {
		if strings.Contains(code, marker) {
			return true
		}
	}
	return false
}

// CWE-1230 only reports executable Content-Disposition construction that uses
// request metadata. Generated, fixed download names are deliberately safe.
func detectCWE1230(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if start := firstCodeMatchStart(facts, unit.Source, pyContentDispositionRequestRE); start >= 0 {
		emitValidationFinding(unit, &MetaCWE1230, start, "response exposes a request-controlled filename through Content-Disposition metadata", confidence80, out)
	}
}

// CWE-1236 recognizes only a direct request expression at a CSV writer row.
// Sanitizers that visibly strip or quote formula prefixes remain out of scope.
func detectCWE1236(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(facts, unit.Source, ".writerow") {
		if isDirectRequestExpr(call.ArgsText) && !strings.Contains(call.ArgsText, "lstrip") && !strings.Contains(call.ArgsText, "safe_csv") {
			emitValidationFinding(unit, &MetaCWE1236, call.Start, "CSV row writes a request-controlled field without formula neutralization", confidence82, out)
			return
		}
	}
}

// CWE-1286 is constrained to one expression so it does not claim arbitrary
// data-flow: json.loads reads request bytes directly into an HTTP client URL.
func detectCWE1286(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"requests.get", "requests.post", "requests.request", "httpx.get", "httpx.post", "httpx.request", "urllib.request.urlopen"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			compact := compactWhitespace(call.ArgsText)
			if strings.Contains(compact, "json.loads(request.") && !strings.Contains(compact, "urlparse(") {
				emitValidationFinding(unit, &MetaCWE1286, call.Start, "request JSON is used as an outbound URL without syntactic validation", confidence80, out)
				return
			}
		}
	}
}

// CWE-1289 requires all three local facts: request-derived path assignment, a
// deny-list equality check, and filesystem use without canonicalization.
func detectCWE1289(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, fn := range facts.Functions() {
		code := facts.codeMask(fn.body, fn.bodyStart)
		match := pyPathDenyListRE.FindStringSubmatchIndex(fn.body)
		if match == nil || strings.Contains(code, "realpath(") || strings.Contains(code, "resolve(") {
			continue
		}
		// The deny-list value is a string literal and therefore absent from the
		// code mask. Still require executable syntax around it so examples in a
		// comment or docstring cannot become a finding.
		if strings.TrimSpace(code[match[0]:match[1]]) == "" {
			continue
		}
		pathName := fn.body[match[2]:match[3]]
		if !strings.Contains(code, pathName+" = request.") || !hasFilesystemUse(code, pathName) {
			continue
		}
		emitValidationFinding(unit, &MetaCWE1289, fn.start+match[0], "request path is protected by exact deny-list equality before filesystem use", confidence80, out)
		return
	}
}

func hasFilesystemUse(code, name string) bool {
	// code is already masked; search it directly without remasking.
	for _, call := range findCallsMasked(code, code, "open", "os.remove", "os.unlink", "send_file", "FileResponse") {
		if containsIdent(call.ArgsText, name) {
			return true
		}
	}
	return false
}

// CWE-1333 restricts the pattern to re.compile so a regular parenthesized
// expression in application code cannot be mistaken for a regular expression.
func detectCWE1333(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(facts, unit.Source, "re.compile", "regex.compile") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) > 0 && pyNestedRegexQuantifierRE.MatchString(args[0]) {
			emitValidationFinding(unit, &MetaCWE1333, call.Start, "compiled regular expression contains nested unbounded quantifiers", confidence84, out)
			return
		}
	}
}

// CWE-1389 reports base-zero parsing only when the parsed argument is visibly
// request-controlled. Application-controlled protocol parsing is not covered.
func detectCWE1389(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(facts, unit.Source, "int") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) >= 2 && strings.TrimSpace(args[1]) == "0" && isDirectRequestExpr(args[0]) {
			emitValidationFinding(unit, &MetaCWE1389, call.Start, "request-controlled numeric input is parsed with base 0", confidence82, out)
			return
		}
	}
}

// CWE-140 is limited to a direct request value manually joined for a response.
// csv.writer and other structured encoders are intentionally not considered.
func detectCWE140(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(facts, unit.Source, ".join") {
		if isDirectRequestExpr(call.ArgsText) && strings.Contains(call.ArgsText, "request.") {
			emitValidationFinding(unit, &MetaCWE140, call.Start, "response manually joins request-controlled fields with a delimiter", confidence76, out)
			return
		}
	}
}

func emitValidationFinding(unit *core.ParsedUnit, meta *rules.RuleMetadata, offset int, message string, confidence float32, out *[]rules.Finding) {
	line, col := unit.LineCol(offset)
	rules.PushFindingWithConfidence(meta, unitFile(unit), line, col, message, confidence, out)
}
