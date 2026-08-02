package cwe

import (
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	// These rules deliberately have no SourceIndex gates: a gate that omits a
	// valid request spelling would be a false-negative source of truth.
	RegisterRule("CWE-918", detectCWE918, &MetaCWE918)
	RegisterRule("CWE-601", detectCWE601, &MetaCWE601)
	RegisterRule("CWE-605", detectCWE605, &MetaCWE605)
	RegisterRule("CWE-924", detectCWE924, &MetaCWE924)
	RegisterRule("CWE-940", detectCWE940, &MetaCWE940)
	RegisterRule("CWE-941", detectCWE941, &MetaCWE941)
}

// CWE-918 reports only direct request-derived URLs at known outbound HTTP
// sinks. Variables passed between functions are deliberately out of scope.
func detectCWE918(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(facts, unit.Source,
		"requests.get", "requests.post", "requests.put", "requests.patch", "requests.delete", "requests.request",
		"httpx.get", "httpx.post", "httpx.put", "httpx.patch", "httpx.delete", "httpx.request",
		"urllib.request.urlopen") {
		if !callHasDirectRequestURL(call.ArgsText) {
			continue
		}
		pushSSRFfinding(unit, call.Start, &MetaCWE918, "request-controlled URL reaches an outbound HTTP client without destination validation", confidence84, out)
		return
	}
}

func callHasDirectRequestURL(argsText string) bool {
	for _, arg := range splitTopLevelArgs(argsText) {
		candidate := strings.TrimSpace(arg)
		if strings.HasPrefix(candidate, "url=") {
			candidate = strings.TrimSpace(strings.TrimPrefix(candidate, "url="))
		} else if strings.Contains(candidate, "=") {
			continue
		}
		return isDirectRequestExpr(candidate)
	}
	return false
}

// CWE-601 requires direct request evidence at a framework redirect sink. A
// fixed local path and an already-validated local variable are not reported.
func detectCWE601(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(facts, unit.Source, "redirect", "flask.redirect", "django.shortcuts.redirect", "HttpResponseRedirect", "RedirectResponse") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) == 0 || !isDirectRequestExpr(args[0]) {
			continue
		}
		pushSSRFfinding(unit, call.Start, &MetaCWE601, "request-controlled URL is passed directly to a redirect response", confidence82, out)
		return
	}
}

// CWE-605 requires both a reuse socket option and a wildcard bind in the same
// file. Either operation alone is common and insufficient evidence.
func detectCWE605(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, fn := range facts.Functions() {
		masked := facts.codeMask(fn.body, fn.bodyStart)
		if start := reuseThenWildcardBind(fn.body, masked); start >= 0 {
			pushSSRFfinding(unit, fn.bodyStart+start, &MetaCWE605, "socket reuse is enabled before binding the same socket to a wildcard interface", confidence76, out)
			return
		}
	}
}

func reuseThenWildcardBind(source, masked string) int {
	for _, call := range findCallsMasked(source, masked, ".setsockopt") {
		if !strings.Contains(call.ArgsText, "SO_REUSEADDR") && !strings.Contains(call.ArgsText, "SO_REUSEPORT") {
			continue
		}
		receiver := callReceiver(source, call.Start)
		if receiver == "" {
			continue
		}
		for _, bind := range findCallsMasked(source, masked, ".bind") {
			if bind.Start <= call.Start || receiver != callReceiver(source, bind.Start) ||
				(!strings.Contains(bind.ArgsText, "0.0.0.0") && !strings.Contains(bind.ArgsText, "::")) {
				continue
			}
			return call.Start
		}
	}
	return -1
}

func callReceiver(source string, callStart int) string {
	if callStart <= 0 || callStart > len(source) {
		return ""
	}
	prefix := source[:callStart]
	lineStart := strings.LastIndex(prefix, "\n") + 1
	receiver := strings.TrimSpace(prefix[lineStart:])
	if receiver == "" || strings.ContainsAny(receiver, " =()[]{}.:+") {
		return ""
	}
	return receiver
}

// CWE-924 keeps webhook integrity checking local to a clearly named handler.
// A signature verification in an unrelated function is not a safe suppression.
func detectCWE924(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	// A module-level authentication gate establishes the caller boundary for
	// internal webhook relays. CWE-924 should not require a second payload
	// signature when the route is protected by header authentication over TLS.
	if hasApplicationAuthGuard(facts) {
		return
	}
	for _, fn := range facts.Functions() {
		name := strings.ToLower(fn.name)
		if !strings.Contains(name, "webhook") && !strings.Contains(name, "_hook") {
			continue
		}
		code := facts.codeMask(fn.body, fn.bodyStart)
		if !hasRequestBodyRead(code) || hasMessageIntegrityVerification(code) {
			continue
		}
		pushSSRFfinding(unit, fn.start, &MetaCWE924, "webhook handler consumes a request body without same-handler signature verification", confidence76, out)
		return
	}
}

func hasApplicationAuthGuard(facts *PyCweFacts) bool {
	if facts == nil {
		return false
	}
	code := facts.Masked
	if !strings.Contains(code, "before_request") || !strings.Contains(code, "request.headers") {
		return false
	}
	return strings.Contains(code, "current_app.config") || strings.Contains(code, "config[") ||
		strings.Contains(code, "api_key") || strings.Contains(code, "authorization")
}

// CWE-940 targets a particularly dangerous callback form: a callback handler
// turns a direct request identity into a login without an apparent state/nonce
// or source verification. General login handlers are intentionally excluded.
func detectCWE940(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, fn := range facts.Functions() {
		name := strings.ToLower(fn.name)
		if !strings.Contains(name, "callback") && !strings.Contains(name, "oauth") {
			continue
		}
		code := facts.codeMask(fn.body, fn.bodyStart)
		if !strings.Contains(code, "request.args") && !strings.Contains(code, "request.GET") && !strings.Contains(code, "request.query_params") {
			continue
		}
		if !strings.Contains(code, "login_user(") && !strings.Contains(code, "session[") {
			continue
		}
		if hasCallbackSourceVerification(code) {
			continue
		}
		pushSSRFfinding(unit, fn.start, &MetaCWE940, "authentication callback trusts a request-controlled identity without source or state verification", confidence78, out)
		return
	}
}

// CWE-941 is intentionally distinct from SSRF: it looks only at direct
// request-controlled recipients passed to SMTP or framework mail APIs.
func detectCWE941(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(facts, unit.Source, "send_mail", ".sendmail", "sendmail") {
		args := splitTopLevelArgs(call.ArgsText)
		if directRequestMailDestination(call.Name, args) {
			pushSSRFfinding(unit, call.Start, &MetaCWE941, "request-controlled recipient is used as an outbound mail destination", confidence80, out)
			return
		}
	}
}

func isDirectRequestExpr(expr string) bool {
	if !isDynamicExpr(expr) {
		return false
	}
	compact := compactWhitespace(expr)
	return strings.Contains(compact, "request.args") || strings.Contains(compact, "request.form") ||
		strings.Contains(compact, "request.data") || strings.Contains(compact, "request.json") ||
		strings.Contains(compact, "request.get_json(") || strings.Contains(compact, "request.GET") ||
		strings.Contains(compact, "request.POST") || strings.Contains(compact, "request.query_params")
}

func hasRequestBodyRead(code string) bool {
	return strings.Contains(code, "request.data") || strings.Contains(code, "request.json") ||
		strings.Contains(code, "request.get_json(") || strings.Contains(code, "request.get_data(") || strings.Contains(code, "request.body")
}

func hasMessageIntegrityVerification(code string) bool {
	return strings.Contains(code, "hmac.compare_digest") || strings.Contains(code, "verify_signature(") ||
		strings.Contains(code, "verify_webhook(") || strings.Contains(code, ".verify(")
}

func hasCallbackSourceVerification(code string) bool {
	return strings.Contains(code, "verify_state(") || strings.Contains(code, "validate_state(") ||
		strings.Contains(code, "verify_nonce(") || strings.Contains(code, "validate_nonce(") ||
		strings.Contains(code, "hmac.compare_digest")
}

func directRequestMailDestination(name string, args []string) bool {
	if name == "send_mail" {
		return len(args) >= 4 && isDirectRequestExpr(args[3])
	}
	// smtplib.SMTP.sendmail(from_addr, to_addrs, msg) puts the recipient
	// address in its second positional argument.
	return len(args) >= 2 && isDirectRequestExpr(args[1])
}

func pushSSRFfinding(unit *core.ParsedUnit, start int, meta *rules.RuleMetadata, message string, confidence float32, out *[]rules.Finding) {
	line, col := unit.LineCol(start)
	rules.PushFindingWithConfidence(meta, unitFile(unit), line, col, message, confidence, out)
}
