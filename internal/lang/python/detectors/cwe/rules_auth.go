package cwe

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	// These rules intentionally have no SourceIndex gates: individual routes
	// and framework calls are detected from combinations of source evidence,
	// and a smaller gate could incorrectly skip a valid combination.
	RegisterRule("CWE-306", detectCWE306, &MetaCWE306)
	RegisterRule("CWE-307", detectCWE307, &MetaCWE307)
	RegisterRule("CWE-346", detectCWE346, &MetaCWE346)
	RegisterRule("CWE-359", detectCWE359, &MetaCWE359)
	RegisterRule("CWE-613", detectCWE613, &MetaCWE613)
	RegisterRule("CWE-565", detectCWE565, &MetaCWE565)
	RegisterRule("CWE-807", detectCWE807, &MetaCWE807)
	RegisterRule("CWE-698", detectCWE698, &MetaCWE698)
}

var (
	pyCriticalRouteRE          = regexp.MustCompile(`(?m)^\s*@(?:[A-Za-z_][A-Za-z0-9_]*\.)?(?:route|get|post|put|patch|delete)\s*\(\s*[rRuUfF]*["']/(?:admin|manage|internal)(?:/|["'])`)
	pyLoginRouteRE             = regexp.MustCompile(`(?m)^\s*@(?:[A-Za-z_][A-Za-z0-9_]*\.)?(?:route|get|post)\s*\(\s*[rRuUfF]*["'][^"'\r\n]*(?:login|sign-in|signin)[^"'\r\n]*["']`)
	pySessionNeverExpiresRE    = regexp.MustCompile(`(?m)^\s*(?:SESSION_COOKIE_AGE|PERMANENT_SESSION_LIFETIME)\s*=\s*0\b`)
	pyDirectSecurityDecisionRE = regexp.MustCompile(`(?im)^\s*if\s+request\.(?:headers|cookies|args)\.get\(\s*["'](?:x-admin|x-role|is_admin|admin|role|user_role)["']\s*\)`)
	pyStandaloneRedirectRE     = regexp.MustCompile(`(?m)^\s*(?:redirect|HttpResponseRedirect)\s*\(`)
	pyPersonalFieldRE          = regexp.MustCompile(`(?i)\b(?:user|profile|customer|account)\.(?:password|passwd|ssn|social_security_number|email|phone|date_of_birth|dob)\b`)
)

// CWE-306 limits itself to framework route decorators on explicitly critical
// paths. Public routes and routes that use a recognizable protection decorator
// are intentionally outside this same-file heuristic.
func detectCWE306(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, match := range pyCriticalRouteRE.FindAllStringIndex(unit.Source, -1) {
		if !codeRangeHasContent(unit.Source, match[0], match[1]) {
			continue
		}
		fn, ok := firstPythonFunctionAfter(unit.Source, match[1])
		if !ok || protectedRouteDecorators(unit.Source[match[0]:fn.start]) {
			continue
		}
		emitAuthFinding(unit, &MetaCWE306, match[0], "critical route lacks a same-route authentication or authorization decorator", 0.76, out)
		return
	}
}

// CWE-307 reports only a login-shaped route that performs an explicit
// password-authentication operation and has no recognized throttle decorator.
func detectCWE307(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, match := range pyLoginRouteRE.FindAllStringIndex(unit.Source, -1) {
		if !codeRangeHasContent(unit.Source, match[0], match[1]) {
			continue
		}
		fn, ok := firstPythonFunctionAfter(unit.Source, match[1])
		if !ok {
			continue
		}
		decorators := strings.ToLower(unit.Source[match[0]:fn.start])
		body := strings.ToLower(pythonCodeMask(fn.body))
		if rateLimitedDecorators(decorators) || !strings.Contains(body, "check_password_hash(") && !strings.Contains(body, "authenticate(") {
			continue
		}
		emitAuthFinding(unit, &MetaCWE307, match[0], "password-authentication route lacks a same-route rate-limit or throttle decorator", 0.74, out)
		return
	}
}

func detectCWE346(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"CORS", "app.add_middleware"} {
		for _, call := range findCalls(unit.Source, name) {
			args := strings.ToLower(compactWhitespace(call.ArgsText))
			if allowsWildcardOrigin(args) && (strings.Contains(args, "supports_credentials=true") || strings.Contains(args, "allow_credentials=true")) {
				emitAuthFinding(unit, &MetaCWE346, call.Start, "CORS allows every origin while credentialed requests are enabled", 0.86, out)
				return
			}
		}
	}
}

func detectCWE359(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"print", "logging.debug", "logging.info", "logging.warning", "logging.error", "logger.debug", "logger.info", "logger.warning", "logger.error"} {
		for _, call := range findCalls(unit.Source, name) {
			if pyPersonalFieldRE.MatchString(call.ArgsText) {
				emitAuthFinding(unit, &MetaCWE359, call.Start, "personal data field is written directly to a log or console sink", 0.8, out)
				return
			}
		}
	}
}

func detectCWE613(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if start := firstMatchStart(unit, pySessionNeverExpiresRE); start >= 0 {
		emitAuthFinding(unit, &MetaCWE613, start, "session lifetime is explicitly configured to never expire", 0.8, out)
	}
}

func detectCWE565(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, fn := range pythonFunctions(unit.Source) {
		code := strings.ToLower(pythonCodeMask(fn.body))
		if hasCookieValidation(code) || !strings.Contains(code, "request.cookies") {
			continue
		}
		for _, call := range findCalls(fn.body, "request.cookies.get", "request.COOKIES.get") {
			if !cookieLooksSecuritySensitive(call.ArgsText) {
				continue
			}
			if !strings.Contains(code, "load_user(") && !strings.Contains(code, "get_user(") && !strings.Contains(code, "authorize(") {
				continue
			}
			emitAuthFinding(unit, &MetaCWE565, fn.bodyStart+call.Start, "security-sensitive cookie value is used without same-function validation evidence", 0.76, out)
			return
		}
	}
}

func detectCWE807(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, match := range pyDirectSecurityDecisionRE.FindAllStringIndex(unit.Source, -1) {
		if !codeRangeHasContent(unit.Source, match[0], match[1]) {
			continue
		}
		fn, ok := containingPythonFunction(unit.Source, match[0])
		if ok && hasCookieValidation(strings.ToLower(pythonCodeMask(fn.body))) {
			continue
		}
		emitAuthFinding(unit, &MetaCWE807, match[0], "client-controlled request value directly controls an authorization decision", 0.82, out)
		return
	}
}

func detectCWE698(unit *core.ParsedUnit, _ *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, fn := range pythonFunctions(unit.Source) {
		for _, match := range pyStandaloneRedirectRE.FindAllStringIndex(fn.body, -1) {
			if !codeRangeHasContent(fn.body, match[0], match[1]) {
				continue
			}
			if next := nextExecutableLine(pythonCodeMask(fn.body), match[1]); next >= 0 {
				emitAuthFinding(unit, &MetaCWE698, fn.bodyStart+match[0], "redirect response is not returned before later code executes", 0.82, out)
				return
			}
		}
	}
}

func firstPythonFunctionAfter(source string, offset int) (pythonFunction, bool) {
	for _, fn := range pythonFunctions(source) {
		if fn.start >= offset {
			return fn, true
		}
	}
	return pythonFunction{}, false
}

func containingPythonFunction(source string, offset int) (pythonFunction, bool) {
	for _, fn := range pythonFunctions(source) {
		if offset >= fn.bodyStart && offset < fn.bodyStart+len(fn.body) {
			return fn, true
		}
	}
	return pythonFunction{}, false
}

func codeRangeHasContent(source string, start, end int) bool {
	if start < 0 || end > len(source) || start >= end {
		return false
	}
	return strings.TrimSpace(pythonCodeMask(source)[start:end]) != ""
}

func protectedRouteDecorators(decorators string) bool {
	lower := strings.ToLower(decorators)
	return strings.Contains(lower, "login_required") || strings.Contains(lower, "auth_required") ||
		strings.Contains(lower, "permission_required") || strings.Contains(lower, "requires_auth") || strings.Contains(lower, "depends(")
}

func rateLimitedDecorators(decorators string) bool {
	return strings.Contains(decorators, "limiter.limit") || strings.Contains(decorators, "ratelimit") ||
		strings.Contains(decorators, "throttle") || strings.Contains(decorators, "axes")
}

func allowsWildcardOrigin(args string) bool {
	return strings.Contains(args, "origins=\"*\"") || strings.Contains(args, "origins='*'") ||
		strings.Contains(args, "allow_origins=[\"*\"]") || strings.Contains(args, "allow_origins=['*']")
}

func cookieLooksSecuritySensitive(args string) bool {
	lower := strings.ToLower(args)
	for _, marker := range []string{"user", "uid", "session", "role", "admin", "auth"} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

func hasCookieValidation(code string) bool {
	return strings.Contains(code, "hmac.compare_digest") || strings.Contains(code, "verify_") ||
		strings.Contains(code, "validate_") || strings.Contains(code, "itsdangerous")
}

func nextExecutableLine(masked string, start int) int {
	for lineStart := start; lineStart < len(masked); {
		if next := strings.IndexByte(masked[lineStart:], '\n'); next >= 0 {
			lineStart += next + 1
		} else {
			return -1
		}
		lineEnd := len(masked)
		if next := strings.IndexByte(masked[lineStart:], '\n'); next >= 0 {
			lineEnd = lineStart + next
		}
		trimmed := strings.TrimSpace(masked[lineStart:lineEnd])
		if trimmed == "" {
			if lineEnd == len(masked) {
				return -1
			}
			continue
		}
		if strings.HasPrefix(trimmed, "return") || strings.HasPrefix(trimmed, "raise") {
			return -1
		}
		return lineStart
	}
	return -1
}

func emitAuthFinding(unit *core.ParsedUnit, meta *rules.RuleMetadata, start int, message string, confidence float32, out *[]rules.Finding) {
	line, col := unit.LineCol(start)
	rules.PushFindingWithConfidence(meta, unitFile(unit), line, col, message, confidence, out)
}
