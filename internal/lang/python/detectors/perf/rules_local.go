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
	materializedAssignRE   = regexp.MustCompile(`\b([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(?:await\s+)?[^\n]*\.(?:scalars\s*\(\s*\)\s*\.)?all\s*\(\s*\)`)
	jsonParseAssignRE      = regexp.MustCompile(`\b([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(?:await\s+)?(?:request\.(?:get_json|json)\s*\(|json\.loads\s*\(\s*request\.(?:body|data))`)
	idempotencyKeyAssignRE = regexp.MustCompile(`(?i)(?:idempotency_key|event_key|request_key)\s*=\s*([A-Za-z_][A-Za-z0-9_\.]*)`)
	ormObjectsModelRE      = regexp.MustCompile(`\b([A-Z][A-Za-z0-9_]*)\.objects\.`)
	ormCtorModelRE         = regexp.MustCompile(`\b([A-Z][A-Za-z0-9_]*)\s*\(`)
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
// The dumps(...) result (or dumps call) must appear in the create/add argument list.
func detectPYPERF9(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if facts == nil || isPythonTestFile(unit) {
		return
	}
	for i, line := range facts.lines {
		m := jsonParseAssignRE.FindStringSubmatch(line.text)
		if len(m) < 2 {
			continue
		}
		name := m[1]
		start, end := functionWindow(facts.lines, i)
		if windowHas(facts.lines, start, end, "redact", "normalize", "sanitiz", "transform") {
			continue
		}
		inlineDumpRE := regexp.MustCompile(`\bjson\.dumps\s*\(\s*` + regexp.QuoteMeta(name) + `\b`)
		dumpAssignRE := regexp.MustCompile(`\b([A-Za-z_][A-Za-z0-9_]*)\s*=\s*json\.dumps\s*\(\s*` + regexp.QuoteMeta(name) + `\b`)
		laterStart, laterEnd := safeLineRange(facts.lines, i+1, end)
		dumpNames := map[string]struct{}{}
		for _, later := range facts.lines[laterStart:laterEnd] {
			if dm := dumpAssignRE.FindStringSubmatch(later.text); len(dm) == 2 {
				dumpNames[dm[1]] = struct{}{}
			}
			if !strings.Contains(later.text, ".add(") && !strings.Contains(later.text, ".objects.create(") && !strings.Contains(later.text, ".create(") {
				continue
			}
			usesDump := inlineDumpRE.MatchString(later.text)
			if !usesDump {
				for dumped := range dumpNames {
					if strings.Contains(later.text, dumped) {
						usesDump = true
						break
					}
				}
			}
			if usesDump {
				pushLine(unit, "PERF-PY-9", line, name, "request JSON is parsed and re-serialized before persistence; persist the raw payload when no transformation is needed", out)
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
// The later insert must share an idempotency/key token with the lookup.
func detectPYPERF14(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if facts == nil || isPythonTestFile(unit) {
		return
	}
	for i, line := range facts.lines {
		if !isIdempotencyLookup(line.text) {
			continue
		}
		keys := idempotencyKeyTokens(line.text)
		if len(keys) == 0 {
			continue
		}
		start, end := functionWindow(facts.lines, i)
		if windowHas(facts.lines, start, end, "get_or_create", "update_or_create", "on_conflict", "ON CONFLICT", "IntegrityError", "insert_or_ignore") {
			continue
		}
		model := ormModelToken(line.text)
		laterStart, laterEnd := safeLineRange(facts.lines, i+1, end)
		for _, later := range facts.lines[laterStart:laterEnd] {
			if !strings.Contains(later.text, ".objects.create(") && !strings.Contains(later.text, ".add(") && !strings.Contains(later.text, "insert(") {
				continue
			}
			if model != "" {
				if createModel := ormModelToken(later.text); createModel != "" && createModel != model && !strings.Contains(later.text, model+"(") {
					continue
				}
			}
			if !lineHasAnyToken(later.text, keys) {
				continue
			}
			pushLine(unit, "PERF-PY-14", line, "idempotency", "idempotency lookup is followed by an insert; use a database upsert or conflict-safe create", out)
			break
		}
	}
}

func isIdempotencyLookup(line string) bool {
	lower := strings.ToLower(line)
	keyed := strings.Contains(lower, "idempotency") || strings.Contains(lower, "event_key") || strings.Contains(lower, "request_key")
	lookup := strings.Contains(lower, ".filter(") || strings.Contains(lower, ".get(") || strings.Contains(lower, "find_") || strings.Contains(lower, "lookup_") || strings.Contains(lower, ".query(")
	return keyed && lookup
}

func idempotencyKeyTokens(line string) []string {
	lower := strings.ToLower(line)
	var keys []string
	for _, tok := range []string{"idempotency_key", "idempotency", "event_key", "request_key"} {
		if strings.Contains(lower, tok) {
			keys = append(keys, tok)
		}
	}
	// Also capture a simple RHS name: filter(idempotency_key=key) → key
	if m := idempotencyKeyAssignRE.FindStringSubmatch(line); len(m) == 2 {
		keys = append(keys, m[1])
	}
	return keys
}

func ormModelToken(line string) string {
	if m := ormObjectsModelRE.FindStringSubmatch(line); len(m) == 2 {
		return m[1]
	}
	if m := ormCtorModelRE.FindStringSubmatch(line); len(m) == 2 {
		return m[1]
	}
	return ""
}

func lineHasAnyToken(line string, tokens []string) bool {
	lower := strings.ToLower(line)
	for _, tok := range tokens {
		if tok != "" && strings.Contains(lower, strings.ToLower(tok)) {
			return true
		}
	}
	return false
}
