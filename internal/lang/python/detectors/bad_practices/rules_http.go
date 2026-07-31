package badpractices

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	metaByID["BP-PY-14"] = &rules.RuleMetadata{
		ID: "BP-PY-14", Title: "requests Without Timeout",
		Description: "`requests` HTTP calls omit `timeout=`, risking hung workers.",
		Severity:    rules.SeverityMedium, Pack: rules.PackBadPractice,
		Fix: "Pass timeout= seconds or timeout=(connect, read) on every requests call.",
	}
	metaByID["BP-PY-15"] = &rules.RuleMetadata{
		ID: "BP-PY-15", Title: "httpx Async Client Not Closed",
		Description: "An `httpx.AsyncClient` is created without `async with` or explicit aclose.",
		Severity:    rules.SeverityMedium, Pack: rules.PackBadPractice,
		Fix: "Prefer `async with httpx.AsyncClient() as client:` or await client.aclose().",
	}
	RegisterRule("BP-PY-14", detectBPPY14)
	RegisterRule("BP-PY-15", detectBPPY15)
}

// requests.<method>(
var requestsCallRe = regexp.MustCompile(`\brequests\.(get|post|put|patch|delete|request|head|options)\s*\(`)

// session-like method calls: session.get(, s.post(, client.get( when Session visible is noisy;
// v0 also flags bare .get/.post etc only when the receiver is clearly a requests Session
// via common names session / sess / req_session.
var sessionHTTPCallRe = regexp.MustCompile(`\b(session|sess|req_session)\.(get|post|put|patch|delete|request|head|options)\s*\(`)

// BP-PY-14: requests.get/post/... without timeout keyword.
func detectBPPY14(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-14")
	if isPythonTestFile(unit) {
		return
	}
	if !facts.has("requests.") && !facts.hasAny("session.", "sess.") {
		// Still allow session.* when source mentions Session
		if !strings.Contains(unit.Source, "requests") && !strings.Contains(unit.Source, "Session") {
			return
		}
	}
	src := unit.Source
	scanHTTPCallsNoTimeout(unit, src, requestsCallRe, meta, out)
	scanHTTPCallsNoTimeout(unit, src, sessionHTTPCallRe, meta, out)
}

func scanHTTPCallsNoTimeout(unit *core.ParsedUnit, src string, re *regexp.Regexp, meta *rules.RuleMetadata, out *[]rules.Finding) {
	for _, m := range re.FindAllStringIndex(src, -1) {
		open := strings.IndexByte(src[m[0]:m[1]], '(')
		if open < 0 {
			continue
		}
		openAbs := m[0] + open
		inner, ok := callArgsRegion(src, openAbs)
		if !ok {
			// Fallback window for unmatched paren
			windowEnd := m[1] + httpTimeoutFallbackBytes
			if windowEnd > len(src) {
				windowEnd = len(src)
			}
			inner = src[m[0]:windowEnd]
		}
		if callArgsHasTimeout(inner) {
			continue
		}
		pushAt(unit, meta, m[0], "requests call missing timeout=; hung network can stall workers", out)
	}
}

const httpTimeoutFallbackBytes = 240

// callArgsHasTimeout reports timeout= as a keyword (not a variable named timeout alone).
func callArgsHasTimeout(args string) bool {
	// Match timeout= outside strings roughly via substring; good enough for v0.
	// Avoid matching timeout in URL query strings inside string literals when possible:
	// look for timeout\s*=
	return timeoutKwRe.MatchString(args)
}

var timeoutKwRe = regexp.MustCompile(`\btimeout\s*=`)

// AsyncClient construction
var asyncClientCallRe = regexp.MustCompile(`(?:httpx\.)?AsyncClient\s*\(`)

// assignment of AsyncClient: name = ...AsyncClient(
var asyncClientAssignRe = regexp.MustCompile(`(?m)^[ \t]*([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(?:httpx\.)?AsyncClient\s*\(`)

// BP-PY-15: httpx.AsyncClient assigned without async with / with and without await <name>.aclose().
// v0: flag bare assignments; miss context-manager forms and same-file aclose of the bound name.
func detectBPPY15(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-15")
	if isPythonTestFile(unit) {
		return
	}
	if !facts.has("AsyncClient") && !strings.Contains(unit.Source, "AsyncClient") {
		return
	}
	src := unit.Source
	lines := codeLinesFacts(facts, unit.Source)

	// Collect bound names that are closed via aclose.
	closed := map[string]struct{}{}
	for _, line := range lines {
		t := strings.TrimSpace(line.text)
		// await client.aclose() / client.aclose()
		if i := strings.Index(t, ".aclose("); i > 0 {
			// extract receiver
			recv := extractReceiverBefore(t, i)
			if recv != "" {
				closed[recv] = struct{}{}
			}
		}
	}

	// Flag assignment forms not in a with/async with line.
	for _, m := range asyncClientAssignRe.FindAllStringSubmatchIndex(src, -1) {
		// m[0]:m[1] full; m[2]:m[3] name
		if len(m) < 4 {
			continue
		}
		fullStart := m[0]
		name := src[m[2]:m[3]]
		// Determine line text containing the assignment.
		lineStart := fullStart
		for lineStart > 0 && src[lineStart-1] != '\n' {
			lineStart--
		}
		lineEnd := fullStart
		for lineEnd < len(src) && src[lineEnd] != '\n' {
			lineEnd++
		}
		lineText := src[lineStart:lineEnd]
		trimmed := strings.TrimSpace(lineText)
		// Context manager: async with httpx.AsyncClient() as client:  (not an assignment form)
		// Also: with AsyncClient() as c: — rare but miss
		if strings.HasPrefix(trimmed, "async with ") || strings.HasPrefix(trimmed, "with ") {
			continue
		}
		// "async with" may appear mid-line after await? skip if line contains " with " before call.
		if containsWithBefore(lineText, "AsyncClient") {
			continue
		}
		if _, ok := closed[name]; ok {
			continue
		}
		pushAt(unit, meta, fullStart, "httpx.AsyncClient not closed; use async with or await client.aclose()", out)
	}

	// Also flag bare AsyncClient() used without assignment/with if it looks like a call statement?
	// Skip — too noisy. Context-manager-only miss covered above; assignment is the main hit path.
	_ = asyncClientCallRe
}

func extractReceiverBefore(line string, dotIdx int) string {
	// Walk back from dotIdx-1 over identifier.
	i := dotIdx - 1
	for i >= 0 && isIdentByte(line[i]) {
		i--
	}
	recv := strings.TrimSpace(line[i+1 : dotIdx])
	// Strip optional "await "
	recv = strings.TrimPrefix(recv, "await ")
	recv = strings.TrimSpace(recv)
	// Only simple names (not chained attr)
	if recv == "" || strings.Contains(recv, ".") {
		// if await client — after strip, client; if obj.client, skip chained
		parts := strings.Fields(recv)
		if len(parts) == 0 {
			return ""
		}
		recv = parts[len(parts)-1]
		if strings.Contains(recv, ".") {
			return ""
		}
	}
	return recv
}

func containsWithBefore(line, needle string) bool {
	idx := strings.Index(line, needle)
	if idx < 0 {
		return false
	}
	prefix := line[:idx]
	return strings.Contains(prefix, "with ") || strings.Contains(prefix, "async with ")
}
