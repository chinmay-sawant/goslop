package perf

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	RegisterRule("PERF-PY-5", detectPERFPY5SequentialDelivery)
	RegisterRule("PERF-PY-7", detectPERFPY7BaseHTTPMiddleware)
	RegisterRule("PERF-PY-17", detectPERFPY17ConnectionControls)
	RegisterRule("PERF-PY-18", detectPERFPY18RepeatedRegexRewrite)
	RegisterRule("PERF-PY-22", detectPERFPY22SQLiteConcurrency)
}

func detectPERFPY5SequentialDelivery(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	for i, line := range facts.lines {
		if !runtimeDeliveryCall(line.text) || !facts.lineInLoop(i) {
			continue
		}
		loop, ok := enclosingLoopHeader(facts.lines, i)
		if !ok {
			continue
		}
		_, iterable, bindOK := perfLoopBinding(loop.text)
		if !bindOK {
			continue
		}
		start, end := functionWindow(facts.lines, i)
		if !runtimeHasBatchClaim(facts.lines, start, end) || runtimeHasBoundedFanout(facts.lines, start, end) {
			continue
		}
		// Require the loop iterable to be the claimed batch (or a name assigned from a claim).
		if !runtimeIterableFromClaim(facts.lines, start, end, iterable) {
			continue
		}
		pushLine(unit, "PERF-PY-5", line, strings.TrimSpace(line.text), "claimed batch is delivered synchronously inside a loop", out)
		return
	}
}

func detectPERFPY7BaseHTTPMiddleware(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	for _, line := range facts.lines {
		className := runtimeBaseMiddlewareClass(line.text)
		if className == "" {
			continue
		}
		for _, registration := range facts.lines {
			if strings.Contains(registration.text, ".add_middleware(") && strings.Contains(registration.text, className) {
				pushLine(unit, "PERF-PY-7", registration, className, "BaseHTTPMiddleware subclass is registered on the request path", out)
				return
			}
		}
	}
}

func detectPERFPY17ConnectionControls(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if runtimeConfigSuppressed(unit) || !runtimeProductionEvidence(facts.Source) {
		return
	}
	lower := strings.ToLower(facts.Source)
	if strings.Contains(lower, "sqlite") || !runtimeDatabaseConfiguration(lower) || runtimeHasConnectionControls(lower) {
		return
	}
	for _, line := range facts.lines {
		if strings.Contains(line.text, "create_engine(") || strings.Contains(line.text, "create_async_engine(") || strings.Contains(line.text, "DATABASES") {
			pushLine(unit, "PERF-PY-17", line, strings.TrimSpace(line.text), "production database configuration has no pool lifecycle, health, or timeout control", out)
			return
		}
	}
}

func detectPERFPY18RepeatedRegexRewrite(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	for i, line := range facts.lines {
		name, ok := runtimeRegexSelfAssignment(line.text)
		if !ok || runtimeDeliberateRegexStage(line.raw) {
			continue
		}
		start, _ := functionWindow(facts.lines, i)
		if !runtimeHotPathFunction(facts.lines, start) {
			continue
		}
		for previous := i - 1; previous >= start && previous >= i-3; previous-- {
			previousName, priorOK := runtimeRegexSelfAssignment(facts.lines[previous].text)
			if priorOK && previousName == name && !runtimeDeliberateRegexStage(facts.lines[previous].raw) {
				pushLine(unit, "PERF-PY-18", line, name, "same input is repeatedly rewritten with regex substitutions on a hot path", out)
				return
			}
		}
	}
}

func detectPERFPY22SQLiteConcurrency(unit *core.ParsedUnit, facts *pyPerfFacts, out *[]rules.Finding) {
	if runtimeConfigSuppressed(unit) || !runtimeProductionEvidence(facts.Source) {
		return
	}
	lower := strings.ToLower(facts.Source)
	if !strings.Contains(lower, "sqlite") || !runtimeConcurrencyEvidence(lower) {
		return
	}
	for _, line := range facts.lines {
		if strings.Contains(strings.ToLower(line.text), "sqlite") {
			pushLine(unit, "PERF-PY-22", line, "sqlite", "SQLite is configured for a concurrent production service write path", out)
			return
		}
	}
}

func runtimeDeliveryCall(line string) bool {
	lower := strings.ToLower(line)
	if strings.Contains(lower, "asyncio.gather") || strings.Contains(lower, "executor") || strings.Contains(lower, "semaphore") || strings.Contains(lower, "limiter") {
		return false
	}
	return strings.Contains(lower, "deliver") || strings.Contains(lower, "send_") ||
		strings.Contains(lower, "publish") || strings.Contains(lower, "webhook") ||
		strings.Contains(lower, "requests.") || strings.Contains(lower, "httpx.")
}

func runtimeHasBatchClaim(lines []codeLine, start, end int) bool {
	start, end = safeLineRange(lines, start, end)
	for _, line := range lines[start:end] {
		lower := strings.ToLower(line.text)
		if strings.Contains(lower, "claim_") || strings.Contains(lower, ".claim(") ||
			strings.Contains(lower, "reserve_") || strings.Contains(lower, "lease_") ||
			strings.Contains(lower, "fetch_pending") || strings.Contains(lower, "select_for_update") {
			return true
		}
	}
	return false
}

func runtimeIterableFromClaim(lines []codeLine, start, end int, iterable string) bool {
	iterable = strings.TrimSpace(iterable)
	if iterable == "" {
		return false
	}
	start, end = safeLineRange(lines, start, end)
	assignRE := regexp.MustCompile(`\b` + regexp.QuoteMeta(iterable) + `\s*=\s*(.+)$`)
	for _, line := range lines[start:end] {
		m := assignRE.FindStringSubmatch(strings.TrimSpace(line.text))
		if len(m) != 2 {
			continue
		}
		rhs := strings.ToLower(m[1])
		if strings.Contains(rhs, "claim_") || strings.Contains(rhs, ".claim(") ||
			strings.Contains(rhs, "reserve_") || strings.Contains(rhs, "lease_") ||
			strings.Contains(rhs, "fetch_pending") || strings.Contains(rhs, "select_for_update") {
			return true
		}
	}
	// Direct for x in claim_pending_batch():
	lower := strings.ToLower(iterable)
	return strings.Contains(lower, "claim_") || strings.Contains(lower, ".claim(") ||
		strings.Contains(lower, "reserve_") || strings.Contains(lower, "lease_") ||
		strings.Contains(lower, "fetch_pending")
}

func runtimeHasBoundedFanout(lines []codeLine, start, end int) bool {
	start, end = safeLineRange(lines, start, end)
	for _, line := range lines[start:end] {
		lower := strings.ToLower(line.text)
		if strings.Contains(lower, "asyncio.gather") || strings.Contains(lower, "semaphore") ||
			strings.Contains(lower, "threadpoolexecutor") || strings.Contains(lower, "processpoolexecutor") ||
			strings.Contains(lower, "limiter") || strings.Contains(lower, "bounded") {
			return true
		}
	}
	return false
}

func runtimeBaseMiddlewareClass(line string) string {
	t := strings.TrimSpace(line)
	if !strings.HasPrefix(t, "class ") || !strings.Contains(t, "BaseHTTPMiddleware") {
		return ""
	}
	rest := strings.TrimPrefix(t, "class ")
	if at := strings.IndexByte(rest, '('); at >= 0 {
		rest = rest[:at]
	}
	name := strings.TrimSpace(rest)
	if !isPyIdentifier(name) {
		return ""
	}
	return name
}

func runtimeConfigSuppressed(unit *core.ParsedUnit) bool {
	if isPythonTestFile(unit) {
		return true
	}
	path := strings.ToLower(fileDisplayPath(unit))
	return strings.Contains(path, "local_") || strings.Contains(path, "/local/") ||
		strings.Contains(path, "dev_") || strings.Contains(path, "/dev/") || strings.Contains(path, "/fixtures/")
}

func runtimeProductionEvidence(source string) bool {
	lower := strings.ToLower(source)
	return strings.Contains(lower, "production") || strings.Contains(lower, "environment=\"prod\"") ||
		strings.Contains(lower, "environment = \"prod\"") || strings.Contains(lower, "deployment_env")
}

func runtimeDatabaseConfiguration(source string) bool {
	return strings.Contains(source, "create_engine(") || strings.Contains(source, "create_async_engine(") || strings.Contains(source, "databases") ||
		strings.Contains(source, "postgresql") || strings.Contains(source, "mysql") || strings.Contains(source, "mariadb")
}

func runtimeHasConnectionControls(source string) bool {
	for _, control := range []string{
		"pool_pre_ping", "pool_recycle", "pool_timeout", "connect_timeout", "conn_max_age", "conn_health_checks", "statement_timeout",
	} {
		if strings.Contains(source, control) {
			return true
		}
	}
	return false
}

func runtimeConcurrencyEvidence(source string) bool {
	for _, signal := range []string{
		"celery", "worker_concurrency", "threadpoolexecutor", "processpoolexecutor",
		"threading.", "multiprocessing", "gunicorn", "uvicorn", "task_queue", "queue.",
	} {
		if strings.Contains(source, signal) {
			return true
		}
	}
	return false
}

func runtimeRegexSelfAssignment(line string) (string, bool) {
	assignment := strings.SplitN(line, "=", 2)
	if len(assignment) != 2 || strings.Contains(assignment[0], "!") || strings.Contains(assignment[0], ">") || strings.Contains(assignment[0], "<") {
		return "", false
	}
	name := strings.TrimSpace(assignment[0])
	rhs := assignment[1]
	if !isPyIdentifier(name) || (!strings.Contains(rhs, "re.sub(") && !strings.Contains(rhs, ".sub(")) || !containsPyIdentifier(rhs, name) {
		return "", false
	}
	return name, true
}

func runtimeDeliberateRegexStage(raw string) bool {
	lower := strings.ToLower(raw)
	return strings.Contains(lower, "staged") || strings.Contains(lower, "overlap") || strings.Contains(lower, "intentional")
}

func runtimeHotPathFunction(lines []codeLine, start int) bool {
	if start < 0 || start >= len(lines) {
		return false
	}
	lower := strings.ToLower(lines[start].text)
	for _, name := range []string{"normalize", "normalise", "clean", "handle", "process", "render", "format"} {
		if strings.Contains(lower, name) {
			return true
		}
	}
	return false
}

func isPyIdentifier(name string) bool {
	if name == "" {
		return false
	}
	for i := 0; i < len(name); i++ {
		c := name[i]
		if !(c == '_' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || i > 0 && c >= '0' && c <= '9') {
			return false
		}
	}
	return true
}

func containsPyIdentifier(text, want string) bool {
	for start := 0; start < len(text); {
		at := strings.Index(text[start:], want)
		if at < 0 {
			return false
		}
		at += start
		end := at + len(want)
		if (at == 0 || !isPyIdentByte(text[at-1])) && (end == len(text) || !isPyIdentByte(text[end])) {
			return true
		}
		start = end
	}
	return false
}

func isPyIdentByte(c byte) bool {
	return c == '_' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9'
}
