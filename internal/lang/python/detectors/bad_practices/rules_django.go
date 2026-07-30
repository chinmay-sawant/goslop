package badpractices

import (
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	// Metadata first so RegisterRule resolves catalogue titles/severities.
	metaByID["BP-PY-22"] = &rules.RuleMetadata{
		ID: "BP-PY-22", Title: "Django SECRET_KEY Hardcoded",
		Description: "Django `SECRET_KEY` is a literal in settings.",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Load SECRET_KEY from the environment or a secrets backend.",
	}
	metaByID["BP-PY-23"] = &rules.RuleMetadata{
		ID: "BP-PY-23", Title: "Django ALLOWED_HOSTS Empty Or Star",
		Description: "`ALLOWED_HOSTS` is empty with DEBUG off risk, or uses `*` in production-like settings.",
		Severity:    rules.SeverityMedium, Pack: rules.PackBadPractice,
		Fix: "Set explicit host names; never use '*' or leave empty when DEBUG is False.",
	}
	metaByID["BP-PY-24"] = &rules.RuleMetadata{
		ID: "BP-PY-24", Title: "Django raw SQL With Format",
		Description: "Django `raw()` / cursor.execute builds SQL via f-string or %-format with variables.",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Use parameterized SQL (params= / %s placeholders with a params list), not string formatting.",
	}
	metaByID["BP-PY-25"] = &rules.RuleMetadata{
		ID: "BP-PY-25", Title: "Django mark_safe On Dynamic Data",
		Description: "`mark_safe` / `SafeString` wraps non-literal values that may be untrusted.",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Prefer auto-escaping templates; only mark trusted, static HTML as safe.",
	}
	metaByID["BP-PY-26"] = &rules.RuleMetadata{
		ID: "BP-PY-26", Title: "Django csrf_exempt On State-Changing View",
		Description: "`@csrf_exempt` disables CSRF on a view that mutates state.",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Remove csrf_exempt for cookie-session views; for webhooks use explicit auth and document the exemption.",
	}
	metaByID["BP-PY-27"] = &rules.RuleMetadata{
		ID: "BP-PY-27", Title: "Django Mass Assignment From request.POST",
		Description: "Model is constructed/updated from full `request.POST` / `request.data` without field allowlists.",
		Severity:    rules.SeverityMedium, Pack: rules.PackBadPractice,
		Fix: "Bind explicit fields or use ModelForm / serializer Meta.fields allowlists.",
	}
	metaByID["BP-PY-28"] = &rules.RuleMetadata{
		ID: "BP-PY-28", Title: "Django N+1 Query In Loop",
		Description: "Related objects are accessed in a loop without select_related/prefetch_related.",
		Severity:    rules.SeverityMedium, Pack: rules.PackBadPractice,
		Fix: "Add select_related/prefetch_related on the queryset (review-only heuristic).",
	}

	RegisterRule("BP-PY-22", detectBPPY22)
	RegisterRule("BP-PY-23", detectBPPY23)
	RegisterRule("BP-PY-24", detectBPPY24)
	RegisterRule("BP-PY-25", detectBPPY25)
	RegisterRule("BP-PY-26", detectBPPY26)
	RegisterRule("BP-PY-27", detectBPPY27)
	RegisterRule("BP-PY-28", detectBPPY28)
}

var (
	djangoSecretKeyRe = regexp.MustCompile(`(?i)\bSECRET_KEY\s*=\s*`)
	allowedHostsRe    = regexp.MustCompile(`(?i)\bALLOWED_HOSTS\s*=\s*`)
	debugFalseAssign  = regexp.MustCompile(`(?i)\bDEBUG\s*=\s*False\b`)
	objectsRawRe      = regexp.MustCompile(`\.objects\.raw\s*\(`)
	markSafeCallRe    = regexp.MustCompile(`\b(?:mark_safe|SafeString)\s*\(`)
	massAssignRe      = regexp.MustCompile(`\(\s*\*\*\s*request\.(?:POST|data|POST\.dict\s*\(\s*\))\s*\)|\.objects\.create\s*\(\s*\*\*\s*request\.|\.update\s*\(\s*\*\*\s*request\.`)
	forObjectsRe      = regexp.MustCompile(`(?i)^\s*for\s+([A-Za-z_][\w]*)\s+in\s+(.+):\s*$`)
)

// looksDjangoish reports settings-path modules or sources with Django markers.
func looksDjangoish(unit *core.ParsedUnit) bool {
	if looksDjangoSettings(unit) {
		return true
	}
	if unit == nil {
		return false
	}
	src := unit.Source
	if strings.Contains(src, "django") || strings.Contains(src, "INSTALLED_APPS") ||
		strings.Contains(src, "MIDDLEWARE") || strings.Contains(src, "get_user_model") ||
		strings.Contains(src, "django.db") || strings.Contains(src, "csrf_exempt") ||
		strings.Contains(src, "mark_safe") || strings.Contains(src, "objects.raw") ||
		strings.Contains(src, "request.POST") || strings.Contains(src, "django.utils") {
		return true
	}
	p := filepath.ToSlash(strings.ToLower(fileDisplayPath(unit)))
	return strings.Contains(p, "/django") || strings.Contains(p, "views.py") ||
		strings.Contains(p, "models.py") || strings.Contains(p, "urls.py")
}

// BP-PY-22: Django SECRET_KEY = '...' in settings modules.
func detectBPPY22(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-22")
	if isPythonTestFile(unit) {
		return
	}
	// Own Django settings; leave Flask-ish modules to BP-PY-17.
	if !looksDjangoSettings(unit) {
		if !looksDjangoish(unit) || looksFlaskish(unit, unit.Source) {
			return
		}
		// Content markers alone without settings path: only fire when SECRET_KEY is bare
		// module-level and path looks settings-like or has INSTALLED_APPS/MIDDLEWARE.
		src := unit.Source
		if !strings.Contains(src, "INSTALLED_APPS") && !strings.Contains(src, "MIDDLEWARE") &&
			!strings.Contains(src, "django") {
			return
		}
	}
	if !facts.hasAny("SECRET_KEY") && !strings.Contains(unit.Source, "SECRET_KEY") {
		return
	}
	lines := codeLinesFacts(facts, unit.Source)
	for _, line := range lines {
		t := strings.TrimSpace(line.text)
		loc := djangoSecretKeyRe.FindStringIndex(t)
		if loc == nil {
			continue
		}
		if strings.Contains(t, "os.environ") || strings.Contains(t, "getenv") ||
			strings.Contains(t, "environ.get") || strings.Contains(t, "config(") ||
			strings.Contains(t, "env(") || strings.Contains(t, "env[") {
			continue
		}
		eq := strings.Index(t[loc[0]:], "=")
		if eq < 0 {
			continue
		}
		rhs := strings.TrimSpace(t[loc[0]+eq+1:])
		if !isStringLiteral(rhs) {
			continue
		}
		if looksLikePlaceholderSecret(rhs) {
			continue
		}
		val := unwrapStringLiteral(rhs)
		if len(val) < 4 {
			continue
		}
		pushAt(unit, meta, line.byte+loc[0], "hardcoded Django SECRET_KEY; load from environment or a secrets backend", out)
	}
}

// BP-PY-23: ALLOWED_HOSTS = ['*'] always; empty list when DEBUG is False nearby.
func detectBPPY23(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-23")
	if isPythonTestFile(unit) {
		return
	}
	if !looksDjangoSettings(unit) && !looksDjangoish(unit) {
		return
	}
	if !facts.has("ALLOWED_HOSTS") && !strings.Contains(unit.Source, "ALLOWED_HOSTS") {
		return
	}
	debugFalse := fileHasDebugFalse(unit.Source)
	lines := codeLinesFacts(facts, unit.Source)
	for _, line := range lines {
		t := strings.TrimSpace(line.text)
		loc := allowedHostsRe.FindStringIndex(t)
		if loc == nil {
			continue
		}
		eq := strings.Index(t[loc[0]:], "=")
		if eq < 0 {
			continue
		}
		rhs := strings.TrimSpace(t[loc[0]+eq+1:])
		if rhs == "" {
			continue
		}
		if allowedHostsHasStar(rhs) {
			pushAt(unit, meta, line.byte+loc[0], "ALLOWED_HOSTS uses '*'; set explicit host names", out)
			continue
		}
		if allowedHostsEmpty(rhs) && debugFalse {
			pushAt(unit, meta, line.byte+loc[0], "ALLOWED_HOSTS is empty while DEBUG is False; set explicit hosts", out)
		}
	}
}

func fileHasDebugFalse(src string) bool {
	for _, raw := range strings.Split(src, "\n") {
		t := strings.TrimSpace(stripPyComment(raw))
		if debugFalseAssign.MatchString(t) {
			return true
		}
	}
	return false
}

func allowedHostsHasStar(rhs string) bool {
	// Star as list/tuple element: '*', "*", or bare * inside brackets.
	if strings.Contains(rhs, `'*'`) || strings.Contains(rhs, `"*"`) {
		return true
	}
	return false
}

func allowedHostsEmpty(rhs string) bool {
	r := strings.TrimSpace(rhs)
	// Strip trailing comma / comments already handled.
	if i := strings.IndexAny(r, "#"); i >= 0 {
		r = strings.TrimSpace(r[:i])
	}
	return r == "[]" || r == "()" || r == "list()" || r == "tuple()"
}

// BP-PY-24: objects.raw / cursor.execute with f-string / .format / % formatting.
func detectBPPY24(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-24")
	src := unit.Source
	if !strings.Contains(src, ".raw(") && !strings.Contains(src, "execute(") &&
		!strings.Contains(src, "objects.raw") {
		return
	}
	// Prefer django-ish modules; still fire on objects.raw anywhere (Django ORM API).
	djangoHint := looksDjangoish(unit) || strings.Contains(src, "objects.raw") ||
		strings.Contains(src, "connection.cursor") || strings.Contains(src, "connections[") ||
		strings.Contains(src, "django.db")

	// objects.raw(
	for _, m := range objectsRawRe.FindAllStringIndex(src, -1) {
		open := m[1] - 1 // points at '('
		if open < 0 || open >= len(src) || src[open] != '(' {
			// find '(' after match
			open = strings.IndexByte(src[m[0]:m[1]], '(')
			if open < 0 {
				continue
			}
			open = m[0] + open
		}
		arg, _, ok := firstCallArg(src, open)
		if !ok {
			continue
		}
		if sqlArgUsesFormat(arg) {
			pushAt(unit, meta, m[0], "Django raw() builds SQL with string formatting; use params binding", out)
		}
	}

	// .execute( — gate when djangoish or connection/cursor nearby.
	if !djangoHint && !strings.Contains(src, "cursor") {
		return
	}
	start := 0
	for {
		idx := strings.Index(src[start:], ".execute(")
		if idx < 0 {
			break
		}
		abs := start + idx
		// Skip non-cursor noise conservatively only when clearly django-related context nearby.
		windowStart := abs - 80
		if windowStart < 0 {
			windowStart = 0
		}
		window := src[windowStart:abs]
		if !djangoHint && !strings.Contains(window, "cursor") && !strings.Contains(window, "connection") {
			start = abs + 9
			continue
		}
		open := abs + len(".execute")
		if open >= len(src) || src[open] != '(' {
			start = abs + 9
			continue
		}
		arg, _, ok := firstCallArg(src, open)
		if ok && sqlArgUsesFormat(arg) {
			pushAt(unit, meta, abs, "cursor.execute builds SQL with string formatting; use params binding", out)
		}
		start = abs + 9
	}
}

// sqlArgUsesFormat reports f-string / .format / % interpolation on a call arg.
func sqlArgUsesFormat(arg string) bool {
	a := strings.TrimSpace(arg)
	if a == "" {
		return false
	}
	// f-string / rf / fr prefixes.
	if isFStringLiteral(a) {
		return true
	}
	if strings.Contains(a, ".format(") {
		return true
	}
	// "..." % expr  or '...' % (
	if sqlPercentFormat(a) {
		return true
	}
	// Implicit concat of f-parts is rare; skip.
	return false
}

func isFStringLiteral(s string) bool {
	s = strings.TrimSpace(s)
	// Strip leading prefixes until we hit quote, tracking f.
	hasF := false
	i := 0
	for i < len(s) {
		c := s[i]
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
	if !hasF || i >= len(s) {
		return false
	}
	return s[i] == '"' || s[i] == '\''
}

func sqlPercentFormat(s string) bool {
	// Detect string-literal % non-string at top level.
	depth := 0
	inStr := byte(0)
	escape := false
	triple := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if inStr != 0 {
			if escape {
				escape = false
				continue
			}
			if c == '\\' && !triple {
				escape = true
				continue
			}
			if triple {
				if c == inStr && i+2 < len(s) && s[i+1] == inStr && s[i+2] == inStr {
					inStr = 0
					triple = false
					i += 2
				}
				continue
			}
			if c == inStr {
				inStr = 0
			}
			continue
		}
		switch c {
		case '"', '\'':
			if i+2 < len(s) && s[i+1] == c && s[i+2] == c {
				inStr = c
				triple = true
				i += 2
			} else {
				inStr = c
			}
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			if depth > 0 {
				depth--
			}
		case '%':
			if depth == 0 {
				// RHS after % should not be another string-only join edge case.
				rest := strings.TrimSpace(s[i+1:])
				if rest == "" {
					return false
				}
				// "%s" % x  or "..." % (
				if rest[0] != '"' && rest[0] != '\'' {
					return true
				}
			}
		}
	}
	return false
}

// BP-PY-25: mark_safe / SafeString on non-literal first arg.
func detectBPPY25(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-25")
	src := unit.Source
	if !strings.Contains(src, "mark_safe") && !strings.Contains(src, "SafeString") {
		return
	}
	for _, m := range markSafeCallRe.FindAllStringIndex(src, -1) {
		// Find '('
		open := strings.IndexByte(src[m[0]:m[1]], '(')
		if open < 0 {
			// match ends at '(' inclusive via \s*\(
			open = m[1] - 1
		} else {
			open = m[0] + open
		}
		if open < 0 || open >= len(src) || src[open] != '(' {
			// scan forward
			j := m[0]
			for j < len(src) && src[j] != '(' {
				j++
			}
			if j >= len(src) {
				continue
			}
			open = j
		}
		arg, _, ok := firstCallArg(src, open)
		if !ok {
			continue
		}
		arg = strings.TrimSpace(arg)
		if arg == "" {
			continue
		}
		// Static string/bytes literal only → miss.
		if isStringLiteral(arg) && !isFStringLiteral(arg) {
			continue
		}
		// f-string is dynamic.
		pushAt(unit, meta, m[0], "mark_safe/SafeString on dynamic data may enable XSS; prefer auto-escaping", out)
	}
}

// BP-PY-26: @csrf_exempt on views that touch POST/body or non-GET methods.
func detectBPPY26(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-26")
	src := unit.Source
	if !strings.Contains(src, "csrf_exempt") {
		return
	}
	// Find @csrf_exempt then the next def/async def; inspect body for state-changing signals.
	lines := codeLinesFacts(facts, src)
	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		if !strings.HasPrefix(t, "@csrf_exempt") {
			// Also allow csrf_exempt(view) wrapper assignment forms lightly.
			if strings.Contains(t, "csrf_exempt(") && !strings.HasPrefix(t, "def ") &&
				!strings.HasPrefix(t, "from ") && !strings.HasPrefix(t, "import ") {
				// csrf_exempt(my_view) — flag when obvious.
				if strings.Contains(t, "csrf_exempt(") && !strings.Contains(t, "import") {
					// Only decorator-style or assignment wrap; skip import lines.
					if strings.Contains(t, "=") || strings.HasPrefix(t, "csrf_exempt(") {
						pushAt(unit, meta, line.byte, "csrf_exempt wraps a view; ensure CSRF is not required for state changes without alternate auth", out)
					}
				}
			}
			continue
		}
		// Look ahead for def within a few decorator lines.
		defIdx := -1
		for j := i + 1; j < len(lines) && j <= i+6; j++ {
			tj := strings.TrimSpace(lines[j].text)
			if strings.HasPrefix(tj, "@") {
				continue
			}
			if strings.HasPrefix(tj, "def ") || strings.HasPrefix(tj, "async def ") {
				defIdx = j
				break
			}
			break
		}
		if defIdx < 0 {
			// Decorator without immediate def still suspicious.
			pushAt(unit, meta, line.byte, "csrf_exempt on view disables CSRF protection", out)
			continue
		}
		body := functionBodyText(lines, defIdx)
		if viewLooksStateChanging(body) {
			pushAt(unit, meta, line.byte, "csrf_exempt on state-changing view; CSRF bypass risk", out)
		} else {
			// v0: still flag any csrf_exempt on a def (high signal per plan); message notes review.
			pushAt(unit, meta, line.byte, "csrf_exempt disables CSRF on this view; verify it is not state-changing without alternate auth", out)
		}
	}
}

func functionBodyText(lines []codeLine, defIdx int) string {
	if defIdx < 0 || defIdx >= len(lines) {
		return ""
	}
	defIndent := indentWidth(lines[defIdx].raw)
	var b strings.Builder
	for j := defIdx + 1; j < len(lines); j++ {
		raw := lines[j].raw
		trim := strings.TrimSpace(raw)
		if trim == "" {
			continue
		}
		ind := indentWidth(raw)
		if ind <= defIndent && trim != "" {
			break
		}
		b.WriteString(lines[j].text)
		b.WriteByte('\n')
	}
	return b.String()
}

func viewLooksStateChanging(body string) bool {
	if body == "" {
		return false
	}
	needles := []string{
		"request.POST", "request.body", "request.FILES",
		"request.method", "POST", "PUT", "PATCH", "DELETE",
		".save(", ".delete(", ".create(", ".update(",
	}
	for _, n := range needles {
		if strings.Contains(body, n) {
			// Avoid false "POST" in comments already stripped; still may match strings.
			return true
		}
	}
	return false
}

// BP-PY-27: Model(**request.POST) / objects.create(**request.data).
func detectBPPY27(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-27")
	src := unit.Source
	if !strings.Contains(src, "request.POST") && !strings.Contains(src, "request.data") {
		return
	}
	// Primary: **request.POST / **request.data splats.
	for _, m := range massAssignRe.FindAllStringIndex(src, -1) {
		pushAt(unit, meta, m[0], "mass assignment from request.POST/data; bind explicit fields or use a form/serializer allowlist", out)
	}
	// Broader: **request.POST without requiring create( form.
	start := 0
	for {
		idx := strings.Index(src[start:], "**request.")
		if idx < 0 {
			break
		}
		abs := start + idx
		rest := src[abs:]
		if strings.HasPrefix(rest, "**request.POST") || strings.HasPrefix(rest, "**request.data") {
			// Avoid double-report if massAssignRe already covered this offset.
			already := false
			for _, f := range *out {
				if f.RuleID == "BP-PY-27" {
					line, _ := unit.LineCol(abs)
					if f.Line == line {
						already = true
						break
					}
				}
			}
			if !already {
				pushAt(unit, meta, abs, "mass assignment from request.POST/data; bind explicit fields or use a form/serializer allowlist", out)
			}
		}
		start = abs + 10
	}
}

// BP-PY-28: N+1 heuristic — for-loop over queryset with multi-hop attr access and no select/prefetch.
// Review-only; medium confidence; false positives expected.
func detectBPPY28(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-28")
	src := unit.Source
	if !strings.Contains(src, "for ") || !strings.Contains(src, ".objects") {
		return
	}
	// File-level opt-out when every queryset uses select/prefetch is too broad;
	// we check per for-loop iterable expression instead.
	lines := codeLinesFacts(facts, src)
	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		m := forObjectsRe.FindStringSubmatch(t)
		if m == nil {
			continue
		}
		loopVar := m[1]
		iter := m[2]
		if !querysetLike(iter) {
			continue
		}
		// Miss when the iterable (or chained call) already has select/prefetch.
		if strings.Contains(iter, "select_related") || strings.Contains(iter, "prefetch_related") {
			continue
		}
		// Also miss if nearby assignment of the same name used select/prefetch (simple lookback).
		if priorQSHasSelectPrefetch(lines, i, iter) {
			continue
		}
		body := functionBodyText(lines, i) // reuse indent walker: body after for line
		// functionBodyText expects def indent; for for-loops, indent of for line works the same.
		if !loopBodyHasRelationAccess(body, loopVar) {
			continue
		}
		// If body itself calls select_related on items, still N+1 — still flag.
		pushAt(unit, meta, line.byte,
			"possible Django N+1: related attributes accessed in a loop without select_related/prefetch_related (heuristic/review-only)",
			out)
	}
}

func querysetLike(iter string) bool {
	it := strings.TrimSpace(iter)
	if strings.Contains(it, ".objects.") || strings.Contains(it, ".objects(") {
		return true
	}
	// Name that was clearly a queryset call chain ending in all/filter/exclude/iterator.
	for _, suf := range []string{".all()", ".filter(", ".exclude(", ".iterator(", ".all("} {
		if strings.Contains(it, suf) && strings.Contains(it, "objects") {
			return true
		}
	}
	// Variable name only: check is weak; require .objects somewhere in iter.
	return strings.Contains(it, "objects")
}

func priorQSHasSelectPrefetch(lines []codeLine, forIdx int, iter string) bool {
	// If iter is a bare name, look back for name = ...select_related...
	name := strings.TrimSpace(iter)
	if name == "" || strings.ContainsAny(name, ".([{") {
		return false
	}
	// validate identifier
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	assignPrefix := name + " ="
	assignPrefix2 := name + "="
	for j := forIdx - 1; j >= 0 && j >= forIdx-30; j-- {
		t := strings.TrimSpace(lines[j].text)
		if strings.HasPrefix(t, assignPrefix) || strings.HasPrefix(t, assignPrefix2) {
			if strings.Contains(t, "select_related") || strings.Contains(t, "prefetch_related") {
				return true
			}
			return false
		}
	}
	return false
}

func loopBodyHasRelationAccess(body, loopVar string) bool {
	if body == "" || loopVar == "" {
		return false
	}
	// Multi-hop: loopVar.foo.bar (relation-like). Single hop loopVar.pk / .id is OK.
	// Conservative v0: require multi-hop only to reduce noise.
	pat := regexp.MustCompile(`\b` + regexp.QuoteMeta(loopVar) + `\.[A-Za-z_][\w]*\.[A-Za-z_][\w]*`)
	return pat.MatchString(body)
}
