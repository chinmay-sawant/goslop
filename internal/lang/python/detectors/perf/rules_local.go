package perf

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	RegisterRule("PERF-PY-1", detectPYPERF1)
	RegisterRule("PERF-PY-9", detectPYPERF9)
	RegisterRule("PERF-PY-12", detectPYPERF12)
	RegisterRule("PERF-PY-14", detectPYPERF14)
}

var (
	materializedAssignRE = regexp.MustCompile(`\b([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(?:await\s+)?[^\n]*\.(?:scalars\s*\(\s*\)\s*\.)?all\s*\(\s*\)`)
	jsonParseAssignRE    = regexp.MustCompile(`\b([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(?:await\s+)?(?:request\.(?:get_json|json)\s*\(|json\.loads\s*\(\s*request\.(?:body|data))`)
)

// PERF-PY-1: materializing an ORM result and sorting it in Python turns an
// aggregate-style read into an unbounded transfer plus an O(N log N) sort.
func detectPYPERF1(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if facts == nil || isPythonTestFile(unit) {
		return
	}
	for i, line := range facts.lines {
		m := materializedAssignRE.FindStringSubmatch(line.text)
		if len(m) < 2 || strings.Contains(line.text, ".limit(") || strings.Contains(line.text, "[:") {
			continue
		}
		start, end := functionWindow(facts.lines, i)
		if windowHas(facts.lines, start, end, "percentile_cont", "func.count", "func.avg", "over(") {
			continue
		}
		name := regexp.QuoteMeta(m[1])
		sortRE := regexp.MustCompile(`\bsorted\s*\(\s*` + name + `\b|\b` + name + `\s*\.sort\s*\(`)
		laterStart, laterEnd := safeLineRange(facts.lines, i+1, end)
		for _, later := range facts.lines[laterStart:laterEnd] {
			if sortRE.MatchString(later.text) {
				pushLine(unit, "PERF-PY-1", line, ".all", "ORM result is fully materialized then sorted in Python; prefer a database aggregate or ordered query", out)
				break
			}
		}
	}
}

// PERF-PY-9: parsing a request payload only to dump it unchanged before
// persistence adds CPU and allocation churn; keep raw bytes when possible.
func detectPYPERF9(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if facts == nil || isPythonTestFile(unit) {
		return
	}
	for i, line := range facts.lines {
		m := jsonParseAssignRE.FindStringSubmatch(line.text)
		if len(m) < 2 {
			continue
		}
		start, end := functionWindow(facts.lines, i)
		if windowHas(facts.lines, start, end, "redact", "normalize", "sanitiz", "transform") {
			continue
		}
		dumpRE := regexp.MustCompile(`\bjson\.dumps\s*\(\s*` + regexp.QuoteMeta(m[1]) + `\b`)
		dumped := false
		laterStart, laterEnd := safeLineRange(facts.lines, i+1, end)
		for _, later := range facts.lines[laterStart:laterEnd] {
			if dumpRE.MatchString(later.text) {
				dumped = true
			}
			if dumped && (strings.Contains(later.text, ".add(") || strings.Contains(later.text, ".objects.create(") || strings.Contains(later.text, ".create(")) {
				pushLine(unit, "PERF-PY-9", line, m[1], "request JSON is parsed and re-serialized before persistence; persist the raw payload when no transformation is needed", out)
				break
			}
		}
	}
}

// PERF-PY-12: route-level JSON decoding should have a visible body-size bound.
func detectPYPERF12(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if facts == nil || isPythonTestFile(unit) {
		return
	}
	for i, line := range facts.lines {
		if !strings.Contains(line.text, "request.get_json(") && !strings.Contains(line.text, "request.json(") && !strings.Contains(line.text, "json.loads(request.body") {
			continue
		}
		start, end := functionWindow(facts.lines, i)
		if !windowHas(facts.lines, start, end, "@app.", "@router.", "@api_view", "APIView", "request") {
			continue
		}
		if windowHas(facts.lines, start, end, "content_length", "Content-Length", "MAX_BODY", "max_body", "max_content", "body_limit", "stream") {
			continue
		}
		pushLine(unit, "PERF-PY-12", line, "request", "request JSON is parsed without a visible body-size limit; bound the payload before decoding", out)
	}
}

// PERF-PY-14: checking an idempotency key and later creating the same model
// is a select-then-insert race unless the database operation is an upsert.
func detectPYPERF14(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if facts == nil || isPythonTestFile(unit) {
		return
	}
	for i, line := range facts.lines {
		if !isIdempotencyLookup(line.text) {
			continue
		}
		start, end := functionWindow(facts.lines, i)
		if windowHas(facts.lines, start, end, "get_or_create", "update_or_create", "on_conflict", "ON CONFLICT", "IntegrityError", "insert_or_ignore") {
			continue
		}
		laterStart, laterEnd := safeLineRange(facts.lines, i+1, end)
		for _, later := range facts.lines[laterStart:laterEnd] {
			if strings.Contains(later.text, ".objects.create(") || strings.Contains(later.text, ".add(") || strings.Contains(later.text, "insert(") {
				pushLine(unit, "PERF-PY-14", line, "idempotency", "idempotency lookup is followed by an insert; use a database upsert or conflict-safe create", out)
				break
			}
		}
	}
}

func isIdempotencyLookup(line string) bool {
	lower := strings.ToLower(line)
	keyed := strings.Contains(lower, "idempotency") || strings.Contains(lower, "event_key") || strings.Contains(lower, "request_key")
	lookup := strings.Contains(lower, ".filter(") || strings.Contains(lower, ".get(") || strings.Contains(lower, "find_") || strings.Contains(lower, "lookup_") || strings.Contains(lower, ".query(")
	return keyed && lookup
}
