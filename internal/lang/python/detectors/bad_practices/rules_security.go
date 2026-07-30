package badpractices

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
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
		inner, _, ok := callArgsRegion(src, openAbs)
		if !ok {
			// Fallback: window search for shell=True near the call.
			windowEnd := m[1] + 200
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

// BP-PY-10: pickle.load / pickle.loads / _pickle / cloudpickle
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
			pushAt(unit, meta, abs, "pickle load can execute arbitrary code; avoid on untrusted data", out)
			start = abs + len(n)
		}
	}
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
		inner, _, ok := callArgsRegion(src, abs+len("yaml.load")-1+1) // point at '('
		// abs is start of "yaml.load(" so open paren at abs+len("yaml.load")
		open := abs + len("yaml.load")
		if open < len(src) && src[open] == '(' {
			inner, _, ok = callArgsRegion(src, open)
		}
		if ok {
			if strings.Contains(inner, "SafeLoader") || strings.Contains(inner, "CSafeLoader") ||
				strings.Contains(inner, "safe_load") {
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
func detectBPPY12(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-12")
	if !facts.hasAny("eval(", "exec(", "compile(") {
		return
	}
	src := unit.Source
	for _, name := range []string{"eval", "exec"} {
		for _, off := range findAllIdent(src, name) {
			end := off + len(name)
			if end >= len(src) || src[end] != '(' {
				continue
			}
			arg, _, ok := firstCallArg(src, end)
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
	for _, off := range findAllIdent(src, "compile") {
		end := off + len("compile")
		if end >= len(src) || src[end] != '(' {
			continue
		}
		inner, _, ok := callArgsRegion(src, end)
		if !ok {
			continue
		}
		if strings.Contains(inner, "'exec'") || strings.Contains(inner, `"exec"`) ||
			strings.Contains(inner, "'eval'") || strings.Contains(inner, `"eval"`) {
			// First arg literal skip
			arg, _, aok := firstCallArg(src, end)
			if aok && isStringLiteral(arg) {
				continue
			}
			pushAt(unit, meta, off, "compile with exec/eval mode on dynamic input is unsafe", out)
		}
	}
}

var secretNameRe = regexp.MustCompile(`(?i)\b([A-Za-z_][A-Za-z0-9_]*(?:password|passwd|secret|api_key|apikey|token|private_key|access_key)[A-Za-z0-9_]*)\s*=\s*`)

// Also match SECRET_KEY, API_KEY, etc. as whole names.
var secretExactAssignRe = regexp.MustCompile(`(?i)\b(password|passwd|secret|secret_key|api_key|apikey|token|private_key|access_key|auth_token)\s*=\s*`)

// BP-PY-13: hardcoded secret heuristic (conservative).
func detectBPPY13(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-13")
	if !facts.hasAny("password", "secret", "api_key", "token", "private_key", "SECRET_KEY") {
		// still check SECRET_KEY case via source
		if !strings.Contains(strings.ToLower(unit.Source), "password") &&
			!strings.Contains(strings.ToLower(unit.Source), "secret") &&
			!strings.Contains(strings.ToLower(unit.Source), "token") &&
			!strings.Contains(strings.ToLower(unit.Source), "api_key") {
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
			// SECRET_KEY = '...'
			if m := regexp.MustCompile(`(?i)\bSECRET_KEY\s*=\s*`).FindStringIndex(t); m != nil {
				loc = m
			}
		}
		if loc == nil {
			continue
		}
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
		// Only string literals.
		if !isStringLiteral(rhs) {
			// RHS might be "foo" + something — still skip non-pure for conservative.
			continue
		}
		if looksLikePlaceholderSecret(rhs) {
			continue
		}
		// Very short literals often placeholders.
		val := unwrapStringLiteral(rhs)
		if len(val) < 6 {
			continue
		}
		pushAt(unit, meta, line.byte+loc[0], "hardcoded secret-like string in source; load from environment or a secret manager", out)
	}
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
