package cwe

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

var pyFunctionDefRE = regexp.MustCompile(`(?m)^([ \t]*)def\s+([A-Za-z_][A-Za-z0-9_]*)\s*\([^\n]*\)\s*:`)

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
func detectCWE918(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(unit.Source,
		"requests.get", "requests.post", "requests.put", "requests.patch", "requests.delete", "requests.request",
		"httpx.get", "httpx.post", "httpx.put", "httpx.patch", "httpx.delete", "httpx.request",
		"urllib.request.urlopen") {
		if !callHasDirectRequestURL(call.ArgsText) {
			continue
		}
		pushSSRFfinding(unit, call.Start, &MetaCWE918, "request-controlled URL reaches an outbound HTTP client without destination validation", 0.84, out)
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
func detectCWE601(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(unit.Source, "redirect", "flask.redirect", "django.shortcuts.redirect", "HttpResponseRedirect", "RedirectResponse") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) == 0 || !isDirectRequestExpr(args[0]) {
			continue
		}
		pushSSRFfinding(unit, call.Start, &MetaCWE601, "request-controlled URL is passed directly to a redirect response", 0.82, out)
		return
	}
}

// CWE-605 requires both a reuse socket option and a wildcard bind in the same
// file. Either operation alone is common and insufficient evidence.
func detectCWE605(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, fn := range pythonFunctions(unit.Source) {
		if start := reuseThenWildcardBind(fn.body); start >= 0 {
			pushSSRFfinding(unit, fn.bodyStart+start, &MetaCWE605, "socket reuse is enabled before binding the same socket to a wildcard interface", 0.76, out)
			return
		}
	}
}

func reuseThenWildcardBind(source string) int {
	for _, call := range findCalls(source, ".setsockopt") {
		if !strings.Contains(call.ArgsText, "SO_REUSEADDR") && !strings.Contains(call.ArgsText, "SO_REUSEPORT") {
			continue
		}
		receiver := callReceiver(source, call.Start)
		if receiver == "" {
			continue
		}
		for _, bind := range findCalls(source, ".bind") {
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
func detectCWE924(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, fn := range pythonFunctions(unit.Source) {
		name := strings.ToLower(fn.name)
		if !strings.Contains(name, "webhook") && !strings.Contains(name, "_hook") {
			continue
		}
		code := pythonCodeMask(fn.body)
		if !hasRequestBodyRead(code) || hasMessageIntegrityVerification(code) {
			continue
		}
		pushSSRFfinding(unit, fn.start, &MetaCWE924, "webhook handler consumes a request body without same-handler signature verification", 0.76, out)
		return
	}
}

// CWE-940 targets a particularly dangerous callback form: a callback handler
// turns a direct request identity into a login without an apparent state/nonce
// or source verification. General login handlers are intentionally excluded.
func detectCWE940(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, fn := range pythonFunctions(unit.Source) {
		name := strings.ToLower(fn.name)
		if !strings.Contains(name, "callback") && !strings.Contains(name, "oauth") {
			continue
		}
		code := pythonCodeMask(fn.body)
		if !strings.Contains(code, "request.args") && !strings.Contains(code, "request.GET") && !strings.Contains(code, "request.query_params") {
			continue
		}
		if !strings.Contains(code, "login_user(") && !strings.Contains(code, "session[") {
			continue
		}
		if hasCallbackSourceVerification(code) {
			continue
		}
		pushSSRFfinding(unit, fn.start, &MetaCWE940, "authentication callback trusts a request-controlled identity without source or state verification", 0.78, out)
		return
	}
}

// CWE-941 is intentionally distinct from SSRF: it looks only at direct
// request-controlled recipients passed to SMTP or framework mail APIs.
func detectCWE941(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(unit.Source, "send_mail", ".sendmail", "sendmail") {
		args := splitTopLevelArgs(call.ArgsText)
		if directRequestMailDestination(call.Name, args) {
			pushSSRFfinding(unit, call.Start, &MetaCWE941, "request-controlled recipient is used as an outbound mail destination", 0.8, out)
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

type pythonFunction struct {
	name      string
	start     int
	bodyStart int
	body      string
}

func pythonFunctions(source string) []pythonFunction {
	code := pythonCodeMask(source)
	matches := pyFunctionDefRE.FindAllStringSubmatchIndex(code, -1)
	out := make([]pythonFunction, 0, len(matches))
	for _, match := range matches {
		indent := len(code[match[2]:match[3]])
		bodyStart := match[1]
		bodyEnd := pythonFunctionBodyEnd(code, bodyStart, indent)
		out = append(out, pythonFunction{name: source[match[4]:match[5]], start: match[0], bodyStart: bodyStart, body: source[bodyStart:bodyEnd]})
	}
	return out
}

func pythonFunctionBodyEnd(code string, bodyStart, indent int) int {
	for lineStart := bodyStart; lineStart < len(code); {
		if code[lineStart] == '\n' {
			lineStart++
		}
		lineEnd := len(code)
		if next := strings.IndexByte(code[lineStart:], '\n'); next >= 0 {
			lineEnd = lineStart + next
		}
		line := code[lineStart:lineEnd]
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			lineIndent := len(line) - len(strings.TrimLeft(line, " \t"))
			if lineIndent <= indent {
				return lineStart
			}
		}
		if lineEnd == len(code) {
			break
		}
		lineStart = lineEnd + 1
	}
	return len(code)
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
