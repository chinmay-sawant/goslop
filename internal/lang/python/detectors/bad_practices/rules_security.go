package badpractices

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/lang/python/pytext"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	RegisterRule("BP-PY-8", detectBPPY8)
	RegisterRule("BP-PY-9", detectBPPY9)
	RegisterRule("BP-PY-10", detectBPPY10)
	RegisterRule("BP-PY-11", detectBPPY11)
	RegisterRule("BP-PY-12", detectBPPY12)
	RegisterRule("BP-PY-13", detectBPPY13)
}

var subprocessCallRe = regexp.MustCompile(`subprocess\.(run|Popen|call|check_output|check_call)\s*\(`)

const subprocessCallFallbackBytes = 200

// BP-PY-8: subprocess.*(… shell=True …)
func detectBPPY8(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-8")
	if !facts.has("subprocess.") || !facts.has("shell=True") {
		return
	}
	src := unit.Source
	for _, m := range subprocessCallRe.FindAllStringIndex(src, -1) {
		open := strings.IndexByte(src[m[0]:m[1]], '(')
		if open < 0 {
			continue
		}
		openAbs := m[0] + open
		inner, ok := callArgsRegion(src, openAbs)
		if !ok {
			// Fallback: window search for shell=True near the call.
			windowEnd := m[1] + subprocessCallFallbackBytes
			if windowEnd > len(src) {
				windowEnd = len(src)
			}
			window := src[m[0]:windowEnd]
			if strings.Contains(window, "shell=True") {
				pushAt(unit, meta, m[0], "subprocess with shell=True enables shell injection; use argv list and shell=False", out)
			}
			continue
		}
		if strings.Contains(inner, "shell=True") {
			pushAt(unit, meta, m[0], "subprocess with shell=True enables shell injection; use argv list and shell=False", out)
		}
	}
}

// BP-PY-9: os.system( / os.popen(
func detectBPPY9(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-9")
	if !facts.has("os.system") && !facts.has("os.popen") {
		return
	}
	src := unit.Source
	for _, needle := range []string{"os.system(", "os.popen("} {
		for _, off := range findAllIdent(src, strings.TrimSuffix(needle, "(")) {
			// ensure followed by (
			end := off + len(strings.TrimSuffix(needle, "("))
			if end < len(src) && src[end] == '(' {
				pushAt(unit, meta, off, "os.system/os.popen runs a shell; prefer subprocess with list argv", out)
			}
		}
	}
}

// BP-PY-10: pickle.load / pickle.loads / _pickle / cloudpickle on non-constant /
// untrusted-looking sources (request body, user path, generic payload names).
// Trusted/local cache-style constant loads and same-process pickle.dumps
// round-trips (especially in tests) are missed.
func detectBPPY10(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-10")
	if !facts.hasAny("pickle.", "cloudpickle.", "_pickle.") {
		return
	}
	src := unit.Source
	needles := []string{
		"pickle.loads(", "pickle.load(",
		"_pickle.loads(", "_pickle.load(",
		"cloudpickle.loads(", "cloudpickle.load(",
	}
	for _, n := range needles {
		start := 0
		for {
			idx := strings.Index(src[start:], n)
			if idx < 0 {
				break
			}
			abs := start + idx
			open := abs + len(n) - 1 // points at '('
			inner, ok := callArgsRegion(src, open)
			if !ok {
				windowEnd := abs + len(n) + 120
				if windowEnd > len(src) {
					windowEnd = len(src)
				}
				inner = src[abs:windowEnd]
			}
			// Same-process round-trip: pickle.loads(pickle.dumps(x)) is trusted
			// (niquests pickleability tests).
			if pickleArgIsRoundTrip(inner) {
				start = abs + len(n)
				continue
			}
			if pickleArgLooksUntrusted(inner) {
				pushAt(unit, meta, abs, "pickle load can execute arbitrary code; avoid on untrusted data", out)
			}
			start = abs + len(n)
		}
	}
}

// pickleArgIsRoundTrip reports pickle.loads(pickle.dumps(...)) style args.
func pickleArgIsRoundTrip(arg string) bool {
	lower := strings.ToLower(arg)
	return strings.Contains(lower, "pickle.dumps(") || strings.Contains(lower, "cloudpickle.dumps(") ||
		strings.Contains(lower, "_pickle.dumps(")
}

func pickleArgLooksUntrusted(arg string) bool {
	lower := strings.ToLower(arg)
	if strings.Contains(lower, "request.") || strings.Contains(lower, "sys.stdin") {
		return true
	}
	untrusted := []string{
		"body", "payload", "data", "raw", "user", "upload", "content",
		"message", "socket", "recv", "input",
	}
	for _, n := range untrusted {
		if strings.Contains(lower, n) {
			return true
		}
	}
	// Literal constant / ALL_CAPS cache names → miss.
	trimmed := strings.TrimSpace(arg)
	if isStringLiteral(trimmed) {
		return false
	}
	if isSimpleIdent(trimmed) && trimmed == strings.ToUpper(trimmed) && len(trimmed) > 1 {
		return false
	}
	if strings.Contains(lower, "cache") && !strings.Contains(lower, "request") {
		return false
	}
	// Bare non-literal name with no trust signal → still flag (conservative).
	return strings.TrimSpace(arg) != ""
}

// BP-PY-11: yaml.load without SafeLoader
func detectBPPY11(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-11")
	if !facts.has("yaml.load") {
		return
	}
	src := unit.Source
	// Do not flag yaml.safe_load
	start := 0
	for {
		idx := strings.Index(src[start:], "yaml.load(")
		if idx < 0 {
			break
		}
		abs := start + idx
		// Skip yaml.safe_load — different needle.
		// Confirm not safe_load: "yaml.load(" is exact.
		inner, ok := callArgsRegion(src, abs+len("yaml.load")-1+1) // point at '('
		// abs is start of "yaml.load(" so open paren at abs+len("yaml.load")
		open := abs + len("yaml.load")
		if open < len(src) && src[open] == '(' {
			inner, ok = callArgsRegion(src, open)
		}
		if ok {
			if strings.Contains(inner, "SafeLoader") || strings.Contains(inner, "CSafeLoader") ||
				strings.Contains(inner, "safe_load") {
				start = abs + len("yaml.load(")
				continue
			}
			// ruamel.YAML().load is not PyYAML's unsafe yaml.load.
			if bpYAMLLoadLooksLikeRuamel(src, inner) {
				start = abs + len("yaml.load(")
				continue
			}
			// Loader= present but not safe: still flag (FullLoader/UnsafeLoader).
			// Bare yaml.load without Loader also flags.
			pushAt(unit, meta, abs, "yaml.load without SafeLoader can execute code; use yaml.safe_load", out)
		} else {
			// Unbalanced — flag conservatively.
			pushAt(unit, meta, abs, "yaml.load without SafeLoader can execute code; use yaml.safe_load", out)
		}
		start = abs + len("yaml.load(")
	}
}

// BP-PY-12: eval( / exec( / compile(..., 'exec')
// Only the builtins are in scope. Attribute methods such as session.exec /
// app.exec / db.exec are SQL/Qt/container APIs, not the exec builtin, unless
// the receiver is the explicit builtins module. def/async def signatures that
// merely declare a method named exec/eval are not call sites. Matches inside
// comments and string literals are ignored via pytext.Mask.
func detectBPPY12(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-12")
	if !facts.hasAny("eval(", "exec(", "compile(") {
		return
	}
	src := unit.Source
	masked := pytext.Mask(src)
	for _, name := range []string{"eval", "exec"} {
		for _, off := range findAllIdent(masked, name) {
			if off > 0 && masked[off-1] == '.' {
				if !bpReceiverIsBuiltins(masked, off) {
					continue
				}
			}
			end := off + len(name)
			if end >= len(src) || src[end] != '(' {
				continue
			}
			if bpIsDefHeaderIdent(src, off) {
				continue
			}
			arg, ok := firstCallArg(src, end)
			if !ok {
				pushAt(unit, meta, off, "eval/exec enables arbitrary code execution", out)
				continue
			}
			// Skip pure string-literal args (static snippets) at v0 — still often unsafe;
			// plan: flag non-literal. Also flag literal for high signal? Plan says non-literal.
			if isStringLiteral(arg) {
				continue
			}
			pushAt(unit, meta, off, "eval/exec on dynamic input enables arbitrary code execution", out)
		}
	}
	// compile(..., 'exec') / "exec"
	for _, off := range findAllIdent(masked, "compile") {
		if off > 0 && masked[off-1] == '.' {
			continue
		}
		end := off + len("compile")
		if end >= len(src) || src[end] != '(' {
			continue
		}
		if bpIsDefHeaderIdent(src, off) {
			continue
		}
		inner, ok := callArgsRegion(src, end)
		if !ok {
			continue
		}
		if strings.Contains(inner, "'exec'") || strings.Contains(inner, `"exec"`) ||
			strings.Contains(inner, "'eval'") || strings.Contains(inner, `"eval"`) {
			// First arg literal skip
			arg, aok := firstCallArg(src, end)
			if aok && isStringLiteral(arg) {
				continue
			}
			pushAt(unit, meta, off, "compile with exec/eval mode on dynamic input is unsafe", out)
		}
	}
}

// bpIsDefHeaderIdent reports whether the identifier at off is the declared
// name of a def/async def signature (e.g. `def exec(self, sql): ...`), not a
// call site.
func bpIsDefHeaderIdent(src string, off int) bool {
	if off <= 0 || off > len(src) {
		return false
	}
	lineStart := off
	for lineStart > 0 && src[lineStart-1] != '\n' {
		lineStart--
	}
	prefix := strings.TrimSpace(src[lineStart:off])
	return prefix == "def" || prefix == "async def"
}

// bpReceiverIsBuiltins reports whether the attribute receiver immediately left
// of the identifier at identAt is the explicit builtins module (builtins.exec
// / builtins.eval are the builtin, unlike session.exec / app.exec).
func bpReceiverIsBuiltins(masked string, identAt int) bool {
	dot := identAt - 1
	if dot <= 0 || masked[dot] != '.' {
		return false
	}
	end := dot
	start := end
	for start > 0 {
		c := masked[start-1]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
			start--
			continue
		}
		break
	}
	return masked[start:end] == "builtins"
}

var secretNameRe = regexp.MustCompile(`(?i)\b([A-Za-z_][A-Za-z0-9_]*(?:password|passwd|secret|api_key|apikey|token|private_key|access_key)[A-Za-z0-9_]*)\s*=\s*`)

// Also match SECRET_KEY, API_KEY, etc. as whole names (case-insensitive).
var secretExactAssignRe = regexp.MustCompile(`(?i)\b(password|passwd|secret|secret_key|api_key|apikey|token|private_key|access_key|auth_token)\s*=\s*`)

// BP-PY-13: hardcoded secret heuristic (conservative).
func detectBPPY13(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-13")
	if isPythonTestFile(unit) || isPythonBenchmarkFile(unit) {
		return
	}
	if !facts.hasAny("password", "secret", "api_key", "token", "private_key", "SECRET_KEY") {
		// Case-fold once for the gate-miss path (avoids up to 4× ToLower on full source).
		lower := strings.ToLower(unit.Source)
		if !strings.Contains(lower, "password") &&
			!strings.Contains(lower, "secret") &&
			!strings.Contains(lower, "token") &&
			!strings.Contains(lower, "api_key") {
			return
		}
	}
	lines := codeLinesFacts(facts, unit.Source)
	for _, line := range lines {
		t := strings.TrimSpace(line.text)
		if t == "" || strings.HasPrefix(t, "def ") || strings.HasPrefix(t, "class ") {
			continue
		}
		// Skip function kwargs that are not simple assignments at module/class level? Still flag.
		// Match assignment.
		if !strings.Contains(t, "=") {
			continue
		}
		// env patterns
		if strings.Contains(t, "os.environ") || strings.Contains(t, "getenv(") ||
			strings.Contains(t, "os.getenv") || strings.Contains(t, "environ.get") {
			continue
		}
		loc := secretExactAssignRe.FindStringIndex(t)
		if loc == nil {
			loc = secretNameRe.FindStringIndex(t)
		}
		if loc == nil {
			continue
		}
		// LHS name (for *_NAME / env-key holders).
		lhs := strings.TrimSpace(t[loc[0] : loc[0]+strings.Index(t[loc[0]:], "=")])
		// RHS after =
		eq := strings.Index(t[loc[0]:], "=")
		if eq < 0 {
			continue
		}
		rhs := strings.TrimSpace(t[loc[0]+eq+1:])
		// Strip trailing comments already handled; strip inline call junk.
		if rhs == "" {
			continue
		}
		// Only pure string literals — reject concatenations like
		// token = "{" + key + "}" (Project_Parva route binder) and f-strings
		// that interpolate runtime values (secret = f"…{secrets.token_…}").
		val, ok := pureHardcodedSecretLiteral(rhs)
		if !ok {
			continue
		}
		if looksLikePlaceholderSecret(rhs) {
			continue
		}
		// Env-var *name* holders (OPENAI_API_KEY_NAME = "PYCAPS_OPENAI_API_KEY")
		// are configuration keys, not committed credentials (pycaps).
		if secretAssignLooksLikeEnvKeyName(lhs, val) {
			continue
		}
		// Very short literals often placeholders.
		if len(val) < 6 {
			continue
		}
		pushAt(unit, meta, line.byte+loc[0], "hardcoded secret-like string in source; load from environment or a secret manager", out)
	}
}

// pureHardcodedSecretLiteral accepts only a single static string literal RHS
// (optional b/r/u/f prefixes). Concatenation, calls, and f-string
// interpolations (`{…}`) are rejected — those are not hardcoded secrets.
func pureHardcodedSecretLiteral(rhs string) (string, bool) {
	s := strings.TrimSpace(rhs)
	if s == "" {
		return "", false
	}
	isF := false
	for {
		if len(s) < 2 {
			return "", false
		}
		c0 := s[0]
		if c0 == 'f' || c0 == 'F' {
			isF = true
			s = s[1:]
			continue
		}
		if c0 == 'r' || c0 == 'R' || c0 == 'u' || c0 == 'U' || c0 == 'b' || c0 == 'B' {
			s = s[1:]
			continue
		}
		break
	}
	if len(s) < 2 {
		return "", false
	}
	q := s[0]
	if q != '"' && q != '\'' {
		return "", false
	}
	// Triple-quoted: require exact triple open/close and nothing after.
	if len(s) >= 6 && s[1] == q && s[2] == q {
		close := string([]byte{q, q, q})
		if !strings.HasSuffix(s, close) {
			return "", false
		}
		inner := s[3 : len(s)-3]
		if isF && strings.Contains(inner, "{") {
			return "", false
		}
		return inner, true
	}
	// Single-quoted: find matching end quote with escapes; refuse trailing junk.
	escaped := false
	end := -1
	for i := 1; i < len(s); i++ {
		c := s[i]
		if escaped {
			escaped = false
			continue
		}
		if c == '\\' {
			escaped = true
			continue
		}
		if c == q {
			end = i
			break
		}
	}
	if end < 0 {
		return "", false
	}
	if strings.TrimSpace(s[end+1:]) != "" {
		// e.g. "{" + key + "}" — not a pure literal.
		return "", false
	}
	inner := s[1:end]
	if isF && strings.Contains(inner, "{") {
		return "", false
	}
	return inner, true
}

// secretAssignLooksLikeEnvKeyName reports env-var name holders / ALL_CAPS key
// strings (configuration indirection), not committed secret values.
func secretAssignLooksLikeEnvKeyName(lhs, val string) bool {
	l := strings.ToUpper(strings.TrimSpace(lhs))
	if strings.HasSuffix(l, "_NAME") || strings.HasSuffix(l, "_ENV") ||
		strings.HasSuffix(l, "_ENV_VAR") || strings.HasSuffix(l, "_VAR") ||
		strings.Contains(l, "_KEY_NAME") || strings.Contains(l, "ENV_KEY") {
		return true
	}
	// Value is an env-style name: ALL_CAPS / digits / underscores only, with _.
	if val == "" {
		return false
	}
	hasUnderscore := false
	for i := 0; i < len(val); i++ {
		c := val[i]
		if c == '_' {
			hasUnderscore = true
			continue
		}
		if (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			continue
		}
		return false
	}
	// Require underscore + reasonable length so "SECRET" alone still flags.
	return hasUnderscore && len(val) >= 8
}

func unwrapStringLiteral(s string) string {
	s = strings.TrimSpace(s)
	for {
		if len(s) < 2 {
			break
		}
		c0 := s[0]
		if c0 == 'r' || c0 == 'R' || c0 == 'u' || c0 == 'U' || c0 == 'b' || c0 == 'B' || c0 == 'f' || c0 == 'F' {
			s = s[1:]
			continue
		}
		break
	}
	if len(s) < 2 {
		return s
	}
	q := s[0]
	if q != '"' && q != '\'' {
		return s
	}
	if len(s) >= 6 && s[1] == q && s[2] == q {
		return s[3 : len(s)-3]
	}
	return s[1 : len(s)-1]
}

func bpYAMLLoadLooksLikeRuamel(src, args string) bool {
	compact := strings.ReplaceAll(strings.ReplaceAll(args, " ", ""), "\t", "")
	if strings.Contains(compact, "Loader=") {
		return false
	}
	return strings.Contains(src, "ruamel.yaml") ||
		strings.Contains(src, "from ruamel") ||
		strings.Contains(src, "import ruamel") ||
		strings.Contains(src, "YAML()")
}
