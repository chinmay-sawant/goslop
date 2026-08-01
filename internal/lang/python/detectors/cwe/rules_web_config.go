package cwe

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	// CWE-756/489/1051/1125 stay ungated (varied assignment/call shapes).
	// Inventory firstCodeMatchStart rules use broad any-of FN-safe gates.
	RegisterRule("CWE-756", detectCWE756, &MetaCWE756)
	RegisterRule("CWE-489", detectCWE489, &MetaCWE489)
	RegisterRule("CWE-15", detectCWE15, &MetaCWE15,
		"request.args", "request.form", "request.values", "request.json", "request.data",
		"request.GET", "request.POST", "request.query_params", "os.putenv", "os.environ[")
	RegisterRule("CWE-1051", detectCWE1051, &MetaCWE1051)
	RegisterRule("CWE-1052", detectCWE1052, &MetaCWE1052,
		"SECRET_KEY", "DATABASE_URL", "SQLALCHEMY_DATABASE_URI",
		"create_engine", "psycopg2.connect")
	RegisterRule("CWE-1125", detectCWE1125, &MetaCWE1125)
	RegisterRule("CWE-1188", detectCWE1188, &MetaCWE1188,
		"ALLOWED_HOSTS", "verify=False", "requests.")
	RegisterRule("CWE-921", detectCWE921, &MetaCWE921,
		"open(", "/tmp/", "/var/tmp/", "/dev/shm/")
}

var (
	pyDebugAssignmentRE = regexp.MustCompile(`(?im)^\s*(?:(?:[A-Za-z_][A-Za-z0-9_]*\s*\.\s*)?debug|DEBUG)\s*=\s*True\b|^\s*[A-Za-z_][A-Za-z0-9_]*\s*\.\s*config\s*\[\s*['"]DEBUG['"]\s*\]\s*=\s*True\b`)
	pyExternalConfigRE  = regexp.MustCompile(`(?im)^\s*(?:(?:[A-Za-z_][A-Za-z0-9_]*\s*\.\s*)?config\s*\[\s*['"][^'"]+['"]\s*\]|(?:settings|django\.conf\.settings)\s*\.\s*[A-Za-z_][A-Za-z0-9_]*)\s*=\s*request\.(?:args|form|values|json|data|GET|POST|query_params)\b|^\s*os\s*\.\s*environ\s*\[\s*['"][^'"]+['"]\s*\]\s*=\s*request\.(?:args|form|values|json|data|GET|POST|query_params)\b`)
	pyHardcodedInitRE   = regexp.MustCompile(`(?im)^\s*(?:SECRET_KEY|DATABASE_URL|SQLALCHEMY_DATABASE_URI)\s*=\s*(?:[bruBRU]*)['"][^'"\r\n]{8,}['"]\s*$`)
	pyWildcardHostsRE   = regexp.MustCompile(`(?im)^\s*ALLOWED_HOSTS\s*=\s*\[\s*['"]\*['"]\s*\]`)
	pyTempSensitiveRE   = regexp.MustCompile(`(?i)^/(?:tmp|var/tmp|dev/shm)/[^'"\r\n]*(?:secret|token|password|passwd|credential|private|key|pem)[^'"\r\n]*$`)
	pyNetworkResourceRE = regexp.MustCompile(`(?i)^[bru]*['"]https?://(?:(?:10|127)(?:\.\d{1,3}){3}|192\.168(?:\.\d{1,3}){2}|172\.(?:1[6-9]|2\d|3[01])(?:\.\d{1,3}){2}|localhost):\d{1,5}(?:[/?#][^'"\r\n]*)?['"]$`)
	pyCredentialDSNRE   = regexp.MustCompile(`(?i)^[bru]*['"][a-z][a-z0-9+.-]*://[^'"\s/:]+:[^'"\s@]+@[^'"\s/]+`)
)

// CWE-756 reports production-style debug settings that expose framework error
// detail rather than a custom error page. A generic error handler cannot be
// proved from a single file, so only explicit debug enablement is reported.
func detectCWE756(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := debugEnabledStart(unit, facts); start >= 0 {
		emitWebConfigFinding(unit, &MetaCWE756, start, "debug error output is enabled instead of a custom error page", confidence76, out)
	}
}

// CWE-489 recognizes active framework debug mode and executable debugger
// breakpoints. findCalls masks comments and docstrings before matching calls.
func detectCWE489(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := debugEnabledStart(unit, facts); start >= 0 {
		emitWebConfigFinding(unit, &MetaCWE489, start, "debug mode is enabled in application source", confidence84, out)
		return
	}
	for _, name := range []string{"pdb.set_trace", "breakpoint"} {
		if calls := findCalls(facts, unit.Source, name); len(calls) > 0 {
			emitWebConfigFinding(unit, &MetaCWE489, calls[0].Start, "executable debugger breakpoint remains in application source", confidence84, out)
			return
		}
	}
}

// CWE-15 requires a direct request value to be assigned to an application or
// process setting. Reading an environment variable is intentionally excluded.
func detectCWE15(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstMatchStart(facts, unit, pyExternalConfigRE); start >= 0 {
		emitWebConfigFinding(unit, &MetaCWE15, start, "request-controlled data directly changes an application or process configuration setting", confidence82, out)
		return
	}
	for _, call := range findCalls(facts, unit.Source, "os.putenv") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) >= 2 && isDirectRequestExpr(args[1]) {
			emitWebConfigFinding(unit, &MetaCWE15, call.Start, "request-controlled data is passed to os.putenv", confidence82, out)
			return
		}
	}
}

// CWE-1051 recognizes literal private-network or localhost URLs with an
// explicit port only at outbound HTTP client calls. Static public service URLs
// and runtime-provided endpoints are outside this narrow heuristic.
func detectCWE1051(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"requests.get", "requests.post", "requests.put", "requests.patch", "requests.delete", "requests.request", "httpx.get", "httpx.post", "httpx.put", "httpx.request", "urllib.request.urlopen"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			if pyNetworkResourceRE.MatchString(strings.TrimSpace(httpCallURLArgument(call.Name, call.ArgsText))) {
				emitWebConfigFinding(unit, &MetaCWE1051, call.Start, "outbound client is initialized with a hard-coded private network resource URL", confidence78, out)
				return
			}
		}
	}
}

// CWE-1052 is limited to literal security or database initialization values,
// including credential-bearing DSNs. Ordinary constants are deliberately not
// treated as excessive hard-coded initialization.
func detectCWE1052(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstMatchStart(facts, unit, pyHardcodedInitRE); start >= 0 {
		emitWebConfigFinding(unit, &MetaCWE1052, start, "security or database initialization uses a hard-coded literal", confidence78, out)
		return
	}
	for _, name := range []string{"create_engine", "sqlalchemy.create_engine", "psycopg2.connect"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			if credentialBearingDSN(call.ArgsText) {
				emitWebConfigFinding(unit, &MetaCWE1052, call.Start, "database initialization uses a credential-bearing hard-coded DSN", confidence80, out)
				return
			}
		}
	}
}

// CWE-1125 reports explicit debug-only HTTP routes. Generic administrative
// routes are common and not sufficient evidence of an excessive attack surface.
func detectCWE1125(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{".route", "add_url_rule"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			if isDebugRoute(callURLArgument(call.ArgsText)) {
				emitWebConfigFinding(unit, &MetaCWE1125, call.Start, "debug-only HTTP route is exposed by the application", confidence76, out)
				return
			}
		}
	}
}

// CWE-1188 reports a wildcard Django host policy or an explicit TLS
// verification bypass, both of which select an insecure resource default.
func detectCWE1188(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstMatchStart(facts, unit, pyWildcardHostsRE); start >= 0 {
		emitWebConfigFinding(unit, &MetaCWE1188, start, "host validation is initialized with a wildcard default", confidence80, out)
		return
	}
	for _, name := range []string{"requests.get", "requests.post", "requests.put", "requests.request", "httpx.get", "httpx.post", "httpx.request"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			if hasKwargFalse(call.ArgsText, "verify") {
				emitWebConfigFinding(unit, &MetaCWE1188, call.Start, "network resource is initialized with TLS verification disabled", confidence82, out)
				return
			}
		}
	}
}

// CWE-921 reports writes of secret-shaped files to common globally accessible
// temporary locations. It does not speculate about permission modes or data
// flow, and only reports an explicit sensitive filename plus a write mode.
func detectCWE921(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(facts, unit.Source, "open") {
		if insecureSensitiveTempWrite(call.ArgsText) {
			emitWebConfigFinding(unit, &MetaCWE921, call.Start, "sensitive data file is opened for writing in a shared temporary location", confidence82, out)
			return
		}
	}
}

func debugEnabledStart(unit *core.ParsedUnit, facts *PyCweFacts) int {
	if start := firstMatchStart(facts, unit, pyDebugAssignmentRE); start >= 0 {
		return start
	}
	if unit == nil {
		return -1
	}
	for _, call := range findCalls(facts, unit.Source, ".run") {
		if hasKwargTrue(call.ArgsText, "debug") {
			return call.Start
		}
	}
	return -1
}

func callURLArgument(argsText string) string {
	for _, arg := range splitTopLevelArgs(argsText) {
		key, value, isKeyword := strings.Cut(arg, "=")
		if isKeyword {
			if strings.TrimSpace(key) == "url" {
				return strings.TrimSpace(value)
			}
			continue
		}
		return strings.TrimSpace(arg)
	}
	return ""
}

func httpCallURLArgument(name, argsText string) string {
	args := splitTopLevelArgs(argsText)
	for _, arg := range args {
		key, value, isKeyword := strings.Cut(arg, "=")
		if isKeyword && strings.TrimSpace(key) == "url" {
			return strings.TrimSpace(value)
		}
	}
	if (name == "requests.request" || name == "httpx.request") && len(args) >= 2 {
		return strings.TrimSpace(args[1])
	}
	return callURLArgument(argsText)
}

func credentialBearingDSN(argsText string) bool {
	for _, arg := range splitTopLevelArgs(argsText) {
		_, value, isKeyword := strings.Cut(arg, "=")
		candidate := arg
		if isKeyword {
			candidate = value
		}
		if pyCredentialDSNRE.MatchString(strings.TrimSpace(candidate)) {
			return true
		}
	}
	return false
}

func isDebugRoute(expr string) bool {
	t := strings.TrimSpace(expr)
	if !isPureStringLiteral(t) {
		return false
	}
	t = strings.Trim(t, "'\"")
	t = strings.ToLower(t)
	return t == "/debug" || t == "/_debug" || t == "/__debug__" || strings.HasPrefix(t, "/debug/")
}

func insecureSensitiveTempWrite(argsText string) bool {
	args := splitTopLevelArgs(argsText)
	if len(args) < 2 || !isPureStringLiteral(args[0]) {
		return false
	}
	path := strings.Trim(strings.TrimSpace(args[0]), "'\"")
	if !pyTempSensitiveRE.MatchString(path) {
		return false
	}
	for _, arg := range args[1:] {
		key, value, isKeyword := strings.Cut(arg, "=")
		if isKeyword && strings.TrimSpace(key) != "mode" {
			continue
		}
		mode := strings.Trim(strings.TrimSpace(value), "'\"")
		if !isKeyword {
			mode = strings.Trim(strings.TrimSpace(arg), "'\"")
		}
		if strings.ContainsAny(mode, "wax") {
			return true
		}
	}
	return false
}

func emitWebConfigFinding(unit *core.ParsedUnit, meta *rules.RuleMetadata, offset int, message string, confidence float32, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	line, col := unit.LineCol(offset)
	rules.PushFindingWithConfidence(meta, unitFile(unit), line, col, message, confidence, out)
}
