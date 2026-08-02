package badpractices

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Production hardening batch (BP-PY-48..50).
// BP-PY-14 (requests without timeout) is owned by batch-01 — do not register here.
func init() {
	metaByID["BP-PY-48"] = &rules.RuleMetadata{
		ID: "BP-PY-48", Title: "CORS Allow Origins Star With Credentials",
		Description: "CORS is configured with `*` origins while allowing credentials.",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Name concrete origins instead of `*`, or disable credentials when reflecting any origin.",
	}
	metaByID["BP-PY-49"] = &rules.RuleMetadata{
		ID: "BP-PY-49", Title: "TLS Verification Disabled",
		Description: "HTTP clients disable TLS verification (`verify=False`).",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Keep TLS verification enabled; use a custom CA bundle path when needed instead of verify=False.",
	}
	metaByID["BP-PY-50"] = &rules.RuleMetadata{
		ID: "BP-PY-50", Title: "Django/Flask CSRF Or Session Cookie Insecure Flags",
		Description: "Session/CSRF cookies lack Secure/HttpOnly/SameSite in production settings.",
		Severity:    rules.SeverityMedium, Pack: rules.PackBadPractice,
		Fix: "Set SESSION_COOKIE_SECURE, CSRF_COOKIE_SECURE, and SESSION_COOKIE_HTTPONLY to True in production.",
	}

	RegisterRule("BP-PY-48", detectBPPY48)
	RegisterRule("BP-PY-49", detectBPPY49)
	RegisterRule("BP-PY-50", detectBPPY50)
}

var (
	// FastAPI / Starlette CORSMiddleware kwargs.
	corsAllowOriginsStarRe = regexp.MustCompile(
		`(?i)allow_origins\s*=\s*(?:\[\s*['"]\*['"]\s*\]|['"]\*['"])`,
	)
	corsAllowCredentialsTrueRe = regexp.MustCompile(`(?i)allow_credentials\s*=\s*True\b`)

	// flask-cors
	flaskCORSCallRe            = regexp.MustCompile(`\bCORS\s*\(`)
	flaskSupportsCredentialsRe = regexp.MustCompile(`(?i)supports_credentials\s*=\s*True\b`)
	flaskOriginsStarRe         = regexp.MustCompile(`(?i)(?:origins|resources)\s*=\s*(?:\[\s*['"]\*['"]\s*\]|['"]\*['"]|\{[^}]*['"]\*['"])`)

	// django-cors-headers
	djangoCORSAllOriginsTrueRe = regexp.MustCompile(
		`(?i)\b(?:CORS_ALLOW_ALL_ORIGINS|CORS_ORIGIN_ALLOW_ALL)\s*=\s*True\b`,
	)
	djangoCORSCredentialsTrueRe = regexp.MustCompile(`(?i)\bCORS_ALLOW_CREDENTIALS\s*=\s*True\b`)

	// TLS verification disabled.
	verifyFalseRe = regexp.MustCompile(`(?i)\bverify\s*=\s*False\b`)
	// ssl._create_unverified_context / CERT_NONE (stdlib ssl / urllib3)
	unverifiedContextRe = regexp.MustCompile(`(?i)(?:ssl\.)?_create_unverified_context\s*\(`)
	certNoneRe          = regexp.MustCompile(`(?i)(?:ssl\.)?CERT_NONE\b`)
	assertHostnameFalse = regexp.MustCompile(`(?i)\bassert_hostname\s*=\s*False\b`)

	// Cookie flags set to False.
	cookieFlagFalseRe = regexp.MustCompile(
		`(?i)\b(SESSION_COOKIE_SECURE|CSRF_COOKIE_SECURE|SESSION_COOKIE_HTTPONLY)\s*=\s*False\b`,
	)
	flaskConfigCookieFalseRe = regexp.MustCompile(
		`(?i)(?:app\.config|\.config)\s*\[\s*['"](?:SESSION_COOKIE_SECURE|CSRF_COOKIE_SECURE|SESSION_COOKIE_HTTPONLY)['"]\s*\]\s*=\s*False\b`,
	)
)

const corsCallFallbackBytes = 400

// BP-PY-48: CORS `*` origins combined with credentials True.
//
// Hits when the same unit enables wildcard origins and credentials together:
//   - FastAPI/Starlette CORSMiddleware: allow_origins=["*"] + allow_credentials=True
//   - flask-cors: CORS(..., origins="*", supports_credentials=True)
//   - django-cors-headers: CORS_ALLOW_ALL_ORIGINS=True + CORS_ALLOW_CREDENTIALS=True
//
// Miss: star without credentials; explicit origin list with credentials.
// Test modules are skipped to limit noise from CORS unit tests.
func detectBPPY48(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-48")
	if meta == nil || unit == nil {
		return
	}
	if isPythonTestFile(unit) {
		return
	}
	if !facts.hasAny(
		"CORSMiddleware", "allow_origins", "allow_credentials",
		"supports_credentials", "CORS_ALLOW_ALL_ORIGINS", "CORS_ORIGIN_ALLOW_ALL",
		"CORS_ALLOW_CREDENTIALS", "CORS(",
	) {
		// Fallback when needles not indexed yet for partial matches.
		srcQuick := unit.Source
		if !strings.Contains(srcQuick, "allow_origins") &&
			!strings.Contains(srcQuick, "CORSMiddleware") &&
			!strings.Contains(srcQuick, "CORS_ALLOW") &&
			!strings.Contains(srcQuick, "supports_credentials") &&
			!strings.Contains(srcQuick, "CORS(") {
			return
		}
	}

	src := unit.Source
	msg := "CORS allows credentials with wildcard origins; name concrete origins instead of *"

	// FastAPI / Starlette: unit-level co-occurrence is enough for v0.
	if corsAllowOriginsStarRe.MatchString(src) && corsAllowCredentialsTrueRe.MatchString(src) {
		loc := corsAllowOriginsStarRe.FindStringIndex(src)
		if loc != nil {
			pushAt(unit, meta, loc[0], msg, out)
			return
		}
	}

	// django-cors-headers settings pair.
	if djangoCORSAllOriginsTrueRe.MatchString(src) && djangoCORSCredentialsTrueRe.MatchString(src) {
		loc := djangoCORSAllOriginsTrueRe.FindStringIndex(src)
		if loc != nil {
			pushAt(unit, meta, loc[0], msg, out)
			return
		}
	}

	// flask-cors: CORS(...) call with both supports_credentials=True and star origins.
	for _, m := range flaskCORSCallRe.FindAllStringIndex(src, -1) {
		open := m[1] - 1 // points at '(' if match ends with (
		if open < 0 || open >= len(src) || src[open] != '(' {
			// locate '(' after CORS
			idx := strings.IndexByte(src[m[0]:m[1]], '(')
			if idx < 0 {
				continue
			}
			open = m[0] + idx
		}
		inner, ok := callArgsRegion(src, open)
		if !ok {
			// Window fallback for multi-line or unmatched paren.
			end := m[0] + corsCallFallbackBytes
			if end > len(src) {
				end = len(src)
			}
			inner = src[m[0]:end]
		}
		if flaskSupportsCredentialsRe.MatchString(inner) &&
			(flaskOriginsStarRe.MatchString(inner) || strings.Contains(inner, `["*"]`) ||
				strings.Contains(inner, `['*']`) || strings.Contains(inner, `"*"`) ||
				strings.Contains(inner, `'*'`)) {
			// Prefer origins-related offset when present.
			if loc := flaskOriginsStarRe.FindStringIndex(inner); loc != nil {
				pushAt(unit, meta, open+1+loc[0], msg, out)
			} else {
				pushAt(unit, meta, m[0], msg, out)
			}
			return
		}
	}
}

// BP-PY-49: TLS verification disabled in HTTP clients / SSL context.
//
// Hits: verify=False, ssl._create_unverified_context, CERT_NONE, optional assert_hostname=False.
// Miss: verify=True, omitted verify, verify="/path/to/ca.pem".
// Does not implement BP-PY-14 (timeout) — batch-01 ownership.
// Test modules are skipped.
func detectBPPY49(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-49")
	if meta == nil || unit == nil {
		return
	}
	if isPythonTestFile(unit) {
		return
	}
	// Dual-mode SSL context factories (verify_ssl switch + create_default_context
	// vs CERT_NONE) are intentional (httptap utils) — audited FPs.
	if strings.Contains(unit.Source, "create_default_context") &&
		(strings.Contains(unit.Source, "verify_ssl") || strings.Contains(unit.Source, "legacy")) &&
		(strings.Contains(unit.Source, "CERT_NONE") || strings.Contains(unit.Source, "check_hostname")) {
		return
	}
	if !facts.hasAny("verify=False", "_create_unverified_context", "CERT_NONE", "assert_hostname") {
		srcQuick := unit.Source
		if !strings.Contains(srcQuick, "verify=False") &&
			!strings.Contains(srcQuick, "verify = False") &&
			!strings.Contains(srcQuick, "_create_unverified_context") &&
			!strings.Contains(srcQuick, "CERT_NONE") &&
			!strings.Contains(srcQuick, "assert_hostname") {
			return
		}
	}

	src := unit.Source
	msg := "TLS verification disabled; do not use verify=False or CERT_NONE in production"

	// Line-oriented reporting for clear locations; also covers multi-kwarg calls on one line.
	lines := codeLinesFacts(facts, src)
	seen := map[int]struct{}{}
	report := func(byteOff int) {
		line, _ := unit.LineCol(byteOff)
		if _, ok := seen[line]; ok {
			return
		}
		seen[line] = struct{}{}
		pushAt(unit, meta, byteOff, msg, out)
	}

	for _, line := range lines {
		t := strings.TrimSpace(line.text)
		if t == "" {
			continue
		}
		if verifyFalseRe.MatchString(t) {
			// Fingerprint pinning assigns verify=False next to assert_fingerprint.
			if tlsFingerprintPinNearby(lines, line) || tlsVerifyFalseIsLocalAssign(t) {
				continue
			}
			if loc := verifyFalseRe.FindStringIndex(t); loc != nil {
				report(line.byte + loc[0])
			}
			continue
		}
		if unverifiedContextRe.MatchString(t) {
			if loc := unverifiedContextRe.FindStringIndex(t); loc != nil {
				report(line.byte + loc[0])
			}
			continue
		}
		if certNoneRe.MatchString(t) {
			// State comparisons (cert_reqs != "CERT_NONE") are not disables.
			if strings.Contains(t, "!=") || strings.Contains(t, "==") {
				continue
			}
			if tlsFingerprintPinNearby(lines, line) {
				continue
			}
			if loc := certNoneRe.FindStringIndex(t); loc != nil {
				report(line.byte + loc[0])
			}
			continue
		}
		if assertHostnameFalse.MatchString(t) {
			if loc := assertHostnameFalse.FindStringIndex(t); loc != nil {
				report(line.byte + loc[0])
			}
		}
	}

	// Multi-line call: verify=False may sit alone on a continued line already covered above.
	// Also catch source-wide when only multi-line patterns exist without line trim issues.
	// Do not resurrect fingerprint-pinning or local verify=False assignments skipped above.
	if len(seen) == 0 {
		reportTLSFallbackHit(src, report)
	}
}

func reportTLSFallbackHit(src string, report func(int)) {
	if loc := verifyFalseRe.FindStringIndex(src); loc != nil {
		snippet := src[loc[0]:]
		if end := strings.IndexByte(snippet, '\n'); end >= 0 {
			snippet = snippet[:end]
		}
		if !tlsVerifyFalseIsLocalAssign(snippet) && !strings.Contains(src, "assert_fingerprint") {
			report(loc[0])
		}
		return
	}
	if loc := unverifiedContextRe.FindStringIndex(src); loc != nil {
		report(loc[0])
		return
	}
	loc := certNoneRe.FindStringIndex(src)
	if loc == nil {
		return
	}
	line := src[loc[0]:]
	if end := strings.IndexByte(line, '\n'); end >= 0 {
		line = line[:end]
	}
	if !strings.Contains(line, "!=") && !strings.Contains(line, "==") {
		report(loc[0])
	}
}

// BP-PY-50: Django/Flask session or CSRF cookie flags set to False.
//
// Hits: SESSION_COOKIE_SECURE / CSRF_COOKIE_SECURE / SESSION_COOKIE_HTTPONLY = False
// and Flask app.config['SESSION_COOKIE_SECURE'] = False (same keys).
// Miss: same flags = True.
// Explicit False only for v0 (no "missing when DEBUG=False" inference).
// Prefer settings-like paths but still flag explicit assignments anywhere outside tests.
func detectBPPY50(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-50")
	if meta == nil || unit == nil {
		return
	}
	if isPythonTestFile(unit) {
		return
	}
	if !facts.hasAny(
		"SESSION_COOKIE_SECURE", "CSRF_COOKIE_SECURE", "SESSION_COOKIE_HTTPONLY",
	) {
		srcQuick := unit.Source
		if !strings.Contains(srcQuick, "SESSION_COOKIE_SECURE") &&
			!strings.Contains(srcQuick, "CSRF_COOKIE_SECURE") &&
			!strings.Contains(srcQuick, "SESSION_COOKIE_HTTPONLY") {
			return
		}
	}

	lines := codeLinesFacts(facts, unit.Source)
	for _, line := range lines {
		t := strings.TrimSpace(line.text)
		if t == "" {
			continue
		}
		if loc := cookieFlagFalseRe.FindStringIndex(t); loc != nil {
			pushAt(unit, meta, line.byte+loc[0],
				"session/CSRF cookie Secure or HttpOnly set to False; enable in production", out)
			continue
		}
		if loc := flaskConfigCookieFalseRe.FindStringIndex(t); loc != nil {
			pushAt(unit, meta, line.byte+loc[0],
				"session/CSRF cookie Secure or HttpOnly set to False; enable in production", out)
		}
	}
}

func tlsVerifyFalseIsLocalAssign(t string) bool {
	trimmed := strings.TrimSpace(t)
	// Bare assignment "verify = False" (not a call kwarg like get(..., verify=False)).
	if strings.Contains(trimmed, "(") {
		return false
	}
	return regexp.MustCompile(`(?i)^verify\s*=\s*False\b`).MatchString(trimmed)
}

func tlsFingerprintPinNearby(lines []codeLine, line codeLine) bool {
	const window = 12
	idx := -1
	for i, l := range lines {
		if l.byte == line.byte {
			idx = i
			break
		}
	}
	if idx < 0 {
		// Fall back to scanning all lines when byte identity differs.
		for _, l := range lines {
			if strings.Contains(l.text, "assert_fingerprint") {
				// Still require proximity by line content match below.
				break
			}
		}
		for i, l := range lines {
			if strings.TrimSpace(l.text) == strings.TrimSpace(line.text) {
				idx = i
				break
			}
		}
	}
	if idx < 0 {
		return false
	}
	lo := idx - window
	if lo < 0 {
		lo = 0
	}
	hi := idx + window
	if hi >= len(lines) {
		hi = len(lines) - 1
	}
	for i := lo; i <= hi; i++ {
		if strings.Contains(lines[i].text, "assert_fingerprint") {
			return true
		}
	}
	return false
}
