package badpractices

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

const assertionErrorType = "AssertionError"

func init() {
	metaByID["BP-PY-41"] = &rules.RuleMetadata{
		ID: "BP-PY-41", Title: "pytest assert With Side Effects Only",
		Description: "Tests call production code without assertions (placeholder tests).",
		Severity:    rules.SeverityInfo, Pack: rules.PackBadPractice,
		Fix: "Add assert / pytest.raises / unittest assertions so the test verifies outcomes.",
	}
	metaByID["BP-PY-42"] = &rules.RuleMetadata{
		ID: "BP-PY-42", Title: "unittest Assert Without Context On Raises",
		Description: "Tests use a bare try/except instead of assertRaises/pytest.raises.",
		Severity:    rules.SeverityLow, Pack: rules.PackBadPractice,
		Fix: "Use `with pytest.raises(...)` or `with self.assertRaises(...)`.",
	}
	RegisterRule("BP-PY-41", detectBPPY41)
	RegisterRule("BP-PY-42", detectBPPY42)
}

// BP-PY-41: test function has side-effect calls but no assertions (info / style).
func detectBPPY41(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-41")
	// nox/tox session runners orchestrate pytest; they are not placeholder tests.
	if isPythonNoxOrToxFile(unit) {
		return
	}
	// Scope: test files or functions named test_*.
	inTestFile := isPythonTestFile(unit)
	if !inTestFile && !strings.Contains(unit.Source, "def test_") {
		return
	}
	lines := codeLinesFacts(facts, unit.Source)
	// Same-file / same-dir helpers.py defs that assert or raise AssertionError count when called.
	helpers := assertionHelpersForUnit(unit, lines)
	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		if !strings.HasPrefix(t, "def test_") && !strings.HasPrefix(t, "async def test_") {
			continue
		}
		// FastAPI/Flask/Starlette route handlers named test_* are product code,
		// not pytest placeholders (decorator-based, not name-gated).
		if hasWebRouteDecorator(lines, i) {
			continue
		}
		// Only analyze test_* in test files, or any test_* anywhere (placeholder risk).
		defIndent := indentWidth(line.raw)
		hasAssert := false
		hasCall := false
		bodyIndent := -1
		hasPytestFail := false
		hasSocketInfra := false
		hasServerInfra := false
		hasStreamOpen := false
		hasCloseCall := false
		var strCarry pyTripleQuoteCarry
		for j := i + 1; j < len(lines); j++ {
			raw := lines[j].raw
			if strCarry.open {
				strCarry.feed(raw)
				continue
			}
			st := strings.TrimSpace(lines[j].text)
			if st == "" {
				strCarry.feed(raw)
				continue
			}
			ind := indentWidth(raw)
			if ind <= defIndent {
				break
			}
			if bodyIndent < 0 {
				bodyIndent = ind
			}
			// Nested def/class ends our interest for simple bodies? Continue scanning body.
			if strings.HasPrefix(st, "def ") || strings.HasPrefix(st, "async def ") || strings.HasPrefix(st, "class ") {
				// nested — still part of outer until dedent of outer; keep going only for deeper.
				if ind == defIndent+1 || ind <= defIndent {
					break
				}
			}
			// pytest.fail is an assertion only when it guards a call inside a suite
			// (except/if); a top-level trailing pytest.fail after a retry loop is a
			// bare fallback, not a verification.
			if isTestAssertion(st, helpers, ind > bodyIndent) {
				hasAssert = true
				break
			}
			// HTTP client response inspection is an assertion (body/text/json access).
			if isHTTPResponseAssertion(st) {
				hasAssert = true
				break
			}
			// subprocess.run(..., check=True) / check_call raise on failure —
			// the call is the assertion.
			if isSubprocessCheckAssertion(st) {
				hasAssert = true
				break
			}
			// Explicit stream open + .close() is resource-lifecycle verification
			// (not a bare side-effect print-only test).
			if strings.Contains(st, "stream=True") || strings.Contains(st, "stream = True") {
				hasStreamOpen = true
			}
			if strings.Contains(st, ".close(") {
				hasCloseCall = true
			}
			if strings.Contains(st, "pytest.fail(") {
				hasPytestFail = true
			}
			if strings.Contains(st, "socket.") || strings.Contains(st, "socket(") {
				hasSocketInfra = true
			}
			if strings.Contains(st, "Server.") || strings.Contains(st, "threading.Event") ||
				strings.Contains(st, "wait_to_close") {
				hasServerInfra = true
			}
			if looksLikeSideEffectCall(st) {
				hasCall = true
			}
			// Track triple-quoted strings that continue onto later lines so column-0
			// string content (e.g. chat samples, embedded YAML) does not abort the scan.
			strCarry.feed(raw)
		}
		// Attribute-only bodies (no call) are not side-effect-only tests.
		if hasCall && !hasAssert {
			// Server wait infrastructure: socket + Event/Server only.
			if hasSocketInfra && hasServerInfra && !hasPytestFail {
				continue
			}
			// Stream + close is resource-lifecycle verification.
			if hasStreamOpen && hasCloseCall {
				continue
			}
			pushAt(unit, meta, line.byte,
				"test function appears to only perform side effects without assertions (heuristic/info); add assert or pytest.raises",
				out)
		}
	}
}

// isHTTPResponseAssertion reports a pure expression that inspects the response
// body (niquests test_decompress_gzip: r.content.decode(...)). Assignments
// like data = response.json() inside retry loops are not assertions
// (httpmorph test_async_post_via_real_proxy TP).
// Note: bare .close() is NOT an assertion — movielite benchmark clip.close()
// bodies remain TPs; niquests stream-close smokes are name-gated instead.
func isHTTPResponseAssertion(st string) bool {
	if strings.HasPrefix(st, "if ") || strings.HasPrefix(st, "elif ") ||
		strings.HasPrefix(st, "while ") || strings.HasPrefix(st, "for ") {
		return false
	}
	// Assignment of response data is control flow, not a bare assertion.
	if strings.Contains(st, "=") && !strings.Contains(st, "==") && !strings.Contains(st, "!=") &&
		!strings.Contains(st, ">=") && !strings.Contains(st, "<=") {
		return false
	}
	needles := []string{".content", ".text", ".json("}
	for _, n := range needles {
		if strings.Contains(st, n) {
			return true
		}
	}
	return false
}

func isTestAssertion(st string, helpers map[string]struct{}, nestedSuite bool) bool {
	if strings.HasPrefix(st, "assert ") || st == "assert" || strings.HasPrefix(st, "assert(") {
		return true
	}
	if strings.Contains(st, "pytest.raises") || strings.Contains(st, "pytest.warns") {
		return true
	}
	// pytest.fail only counts as an assertion when it guards a call inside a suite
	// (except/if); a top-level trailing pytest.fail after a retry loop is a bare
	// fallback, not a verification (httpmorph test_proxy audit TPs).
	if strings.Contains(st, "pytest.fail(") && nestedSuite {
		return true
	}
	if strings.Contains(st, "self.assert") || strings.Contains(st, "self.fail(") {
		return true
	}
	if strings.Contains(st, "unittest.") && strings.Contains(st, "assert") {
		return true
	}
	// Explicit failure raise in the test body verifies outcomes (httptap
	// test_property_based raise-AssertionError rejection tests).
	if strings.Contains(st, "raise "+assertionErrorType) {
		return true
	}
	// nose/pytest helpers
	if strings.HasPrefix(st, "raises(") || strings.Contains(st, " assert_") {
		return true
	}
	// pytest-regressions / image regression baselines.
	if strings.Contains(st, "regression.check(") || strings.Contains(st, "file_regression.check(") ||
		strings.Contains(st, "image_regression.check(") || strings.Contains(st, "data_regression.check(") {
		return true
	}
	// Calls to assert_* / _assert_* helpers (e.g. self._assert_verapdf(...)).
	if name := callCalleeIdent(st); name != "" {
		if isAssertHelperName(name) {
			return true
		}
		// validate_*/verify_* helpers raise on failure — the call is the
		// assertion (Project_Parva test_trust_schemas_parse_and_validate_examples).
		if isValidationAssertionHelper(name) {
			return true
		}
		// pytest-regressions fixture helper and pytest-benchmark fixture.
		if name == "check_func" || name == "benchmark" {
			return true
		}
		if _, ok := helpers[name]; ok {
			return true
		}
	}
	return false
}

// isValidationAssertionHelper reports validate_*/verify_* callees that raise on
// invalid input — calling them is the test's verification step.
func isValidationAssertionHelper(name string) bool {
	n := strings.ToLower(name)
	return strings.HasPrefix(n, "validate_") || strings.HasPrefix(n, "verify_")
}

// isSubprocessCheckAssertion reports subprocess helpers that raise CalledProcessError
// on non-zero exit — the call is the assertion (no separate assert needed).
func isSubprocessCheckAssertion(st string) bool {
	if strings.Contains(st, "subprocess.check_call(") ||
		strings.Contains(st, "subprocess.check_output(") {
		return true
	}
	if strings.Contains(st, "subprocess.run(") && strings.Contains(st, "check=True") {
		return true
	}
	// from subprocess import check_call / check_output / run
	if strings.HasPrefix(st, "check_call(") || strings.HasPrefix(st, "check_output(") {
		return true
	}
	if strings.HasPrefix(st, "run(") && strings.Contains(st, "check=True") {
		return true
	}
	return false
}

// hasWebRouteDecorator reports FastAPI/Flask/Starlette/Django-style route
// decorators immediately above a def — product handlers named test_*, not tests.
func hasWebRouteDecorator(lines []codeLine, defIdx int) bool {
	for j := defIdx - 1; j >= 0; j-- {
		t := strings.TrimSpace(lines[j].text)
		if t == "" || strings.HasPrefix(t, "#") {
			continue
		}
		if !strings.HasPrefix(t, "@") {
			return false
		}
		if isWebRouteDecoratorLine(t) {
			return true
		}
		// Other decorators (e.g. @staticmethod) — keep scanning upward.
	}
	return false
}

func isWebRouteDecoratorLine(t string) bool {
	// @router.post(...), @app.get(...), @api_view(["POST"]), @action(...), etc.
	needles := []string{
		".get(", ".post(", ".put(", ".patch(", ".delete(",
		".head(", ".options(", ".api_route(", ".route(",
		".websocket(", ".api_view(", ".action(",
		"api_view(", "action(",
	}
	for _, n := range needles {
		if strings.Contains(t, n) {
			return true
		}
	}
	// Bare framework aliases: @get("/x"), less common.
	return false
}

// callCalleeIdent returns the final identifier of a call expression's callee, if any.
// Handles "name(...)", "obj.name(...)", and "x = obj.name(...)" forms.
func callCalleeIdent(st string) string {
	st = strings.TrimSpace(st)
	if st == "" {
		return ""
	}
	// Prefer RHS of a simple assignment so "x = helper(...)" still resolves.
	if eq := strings.IndexByte(st, '='); eq > 0 {
		left := strings.TrimSpace(st[:eq])
		right := strings.TrimSpace(st[eq+1:])
		if left != "" && right != "" && !strings.HasPrefix(right, "=") {
			last := left[len(left)-1]
			if last != '!' && last != '<' && last != '>' && last != ':' && strings.Contains(right, "(") {
				st = right
			}
		}
	}
	paren := strings.IndexByte(st, '(')
	if paren <= 0 {
		return ""
	}
	callee := strings.TrimSpace(st[:paren])
	if callee == "" {
		return ""
	}
	if dot := strings.LastIndexByte(callee, '.'); dot >= 0 {
		callee = callee[dot+1:]
	}
	callee = strings.TrimSpace(callee)
	if !isSimpleIdent(callee) {
		return ""
	}
	return callee
}

func isAssertHelperName(name string) bool {
	return strings.HasPrefix(name, "assert") || strings.HasPrefix(name, "_assert")
}

func assertionHelpersForUnit(unit *core.ParsedUnit, lines []codeLine) map[string]struct{} {
	out := sameFileAssertionHelpers(lines)
	mergeLocalHelpersPy(unit, out)
	return out
}

// sameFileAssertionHelpers returns names of non-test defs whose bodies perform
// assertions (bare assert / unittest) or raise AssertionError.
func sameFileAssertionHelpers(lines []codeLine) map[string]struct{} {
	out := make(map[string]struct{})
	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		name, ok := defNameIfNotTest(t)
		if !ok {
			continue
		}
		defIndent := indentWidth(line.raw)
		var strCarry pyTripleQuoteCarry
		for j := i + 1; j < len(lines); j++ {
			raw := lines[j].raw
			if strCarry.open {
				strCarry.feed(raw)
				continue
			}
			st := strings.TrimSpace(lines[j].text)
			if st == "" {
				strCarry.feed(raw)
				continue
			}
			ind := indentWidth(raw)
			if ind <= defIndent {
				break
			}
			if strings.HasPrefix(st, "assert ") || st == "assert" || strings.HasPrefix(st, "assert(") ||
				strings.Contains(st, "self.assert") || strings.Contains(st, "self.fail(") ||
				strings.Contains(st, "raise "+assertionErrorType) ||
				strings.Contains(st, "pytest.raises") || strings.Contains(st, "pytest.warns") ||
				strings.Contains(st, "pytest.fail(") {
				out[name] = struct{}{}
				break
			}
			strCarry.feed(raw)
		}
	}
	return out
}

// pyTripleQuoteCarry tracks an open triple-quoted string across source lines.
// Used so body scans do not treat string content as code (indent / statements).
type pyTripleQuoteCarry struct {
	open  bool
	quote byte
}

func (s *pyTripleQuoteCarry) feed(line string) {
	i := 0
	for i < len(line) {
		if s.open {
			c := line[i]
			if c == s.quote && i+2 < len(line) && line[i+1] == s.quote && line[i+2] == s.quote {
				s.open = false
				s.quote = 0
				i += 3
				continue
			}
			i++
			continue
		}
		c := line[i]
		// Skip prefixes (r/u/b/f and combinations) before a quote.
		if isPyStringPrefixByte(c) {
			j := i
			for j < len(line) && isPyStringPrefixByte(line[j]) {
				j++
			}
			if j < len(line) && (line[j] == '"' || line[j] == '\'') {
				i = j
				c = line[i]
			}
		}
		if c == '"' || c == '\'' {
			if i+2 < len(line) && line[i+1] == c && line[i+2] == c {
				// Triple quote opens; may close later on this line.
				s.open = true
				s.quote = c
				i += 3
				continue
			}
			// Single-line string: consume until close or EOL.
			q := c
			i++
			escape := false
			for i < len(line) {
				ch := line[i]
				if escape {
					escape = false
					i++
					continue
				}
				if ch == '\\' {
					escape = true
					i++
					continue
				}
				if ch == q {
					i++
					break
				}
				i++
			}
			continue
		}
		i++
	}
}

func isPyStringPrefixByte(c byte) bool {
	return c == 'r' || c == 'R' || c == 'u' || c == 'U' || c == 'b' || c == 'B' || c == 'f' || c == 'F'
}

// mergeLocalHelpersPy credits AssertionError / self.assert* helpers defined in a
// sibling helpers.py (audited layout: test_*.py imports find_object_with from helpers).
func mergeLocalHelpersPy(unit *core.ParsedUnit, into map[string]struct{}) {
	if unit == nil || unit.Path == "" || into == nil {
		return
	}
	dir := filepath.Dir(unit.Path)
	if dir == "" || dir == "." {
		return
	}
	if filepath.Base(unit.Path) == "helpers.py" {
		return
	}
	helperPath := filepath.Join(dir, "helpers.py")
	if filepath.Base(helperPath) != "helpers.py" {
		return
	}
	// Fixed sibling basename only (not user-controlled).
	data, err := os.ReadFile(helperPath) //nolint:gosec // G304: constant "helpers.py" next to the unit path
	if err != nil || len(data) == 0 {
		return
	}
	for name := range sameFileAssertionHelpers(buildCodeLines(string(data))) {
		into[name] = struct{}{}
	}
}

func defNameIfNotTest(t string) (string, bool) {
	var rest string
	switch {
	case strings.HasPrefix(t, "def "):
		rest = strings.TrimSpace(strings.TrimPrefix(t, "def "))
	case strings.HasPrefix(t, "async def "):
		rest = strings.TrimSpace(strings.TrimPrefix(t, "async def "))
	default:
		return "", false
	}
	end := strings.IndexByte(rest, '(')
	if end <= 0 {
		return "", false
	}
	name := strings.TrimSpace(rest[:end])
	if !isSimpleIdent(name) || strings.HasPrefix(name, "test_") {
		return "", false
	}
	return name, true
}

func looksLikeSideEffectCall(st string) bool {
	// Skip control flow / trivial statements.
	if st == "pass" || st == "..." || strings.HasPrefix(st, "return") ||
		strings.HasPrefix(st, "raise ") || strings.HasPrefix(st, "import ") ||
		strings.HasPrefix(st, "from ") || strings.HasPrefix(st, "with ") ||
		strings.HasPrefix(st, "async with ") || strings.HasPrefix(st, "for ") ||
		strings.HasPrefix(st, "while ") || strings.HasPrefix(st, "if ") ||
		strings.HasPrefix(st, "elif ") || strings.HasPrefix(st, "else:") ||
		strings.HasPrefix(st, "try:") || strings.HasPrefix(st, "except") ||
		strings.HasPrefix(st, "finally:") || strings.HasPrefix(st, "class ") ||
		strings.HasPrefix(st, "def ") || strings.HasPrefix(st, "async def ") {
		return false
	}
	// Call: name(...) or attr(...) — not pure assignment without call.
	if strings.Contains(st, "(") && strings.Contains(st, ")") {
		// Exclude assert-like already handled.
		return true
	}
	return false
}

// BP-PY-42: bare try/except in tests instead of assertRaises / pytest.raises.
func detectBPPY42(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-42")
	if !isPythonTestFile(unit) && !strings.Contains(unit.Source, "def test_") {
		return
	}
	if !facts.has("try:") && !strings.Contains(unit.Source, "try:") {
		return
	}
	// If file already uses proper context managers, still flag bare try/except elsewhere.
	lines := codeLinesFacts(facts, unit.Source)
	// Restrict to test function bodies.
	type region struct {
		start, end, indent int
	}
	tests := make([]region, 0, len(lines))
	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		if !strings.HasPrefix(t, "def test_") && !strings.HasPrefix(t, "async def test_") {
			continue
		}
		defIndent := indentWidth(line.raw)
		end := len(lines)
		for j := i + 1; j < len(lines); j++ {
			st := strings.TrimSpace(lines[j].text)
			if st == "" {
				continue
			}
			if indentWidth(lines[j].raw) <= defIndent {
				end = j
				break
			}
		}
		tests = append(tests, region{start: i, end: end, indent: defIndent})
	}
	// Only scan def test_* bodies. Do not fall back to whole-file scanning for
	// *_test.py helpers/benchmarks that merely carry a test-like basename.
	for _, reg := range tests {
		if testRegionCollectsThreadErrors(lines, reg.start, reg.end) {
			continue
		}
		for i := reg.start; i < reg.end; i++ {
			t := strings.TrimSpace(lines[i].text)
			if t != "try:" && !strings.HasPrefix(t, "try:") {
				continue
			}
			tryIndent := indentWidth(lines[i].raw)
			// Look for except AssertionError / bare except / except Exception used as "expect failure".
			for j := i + 1; j < reg.end; j++ {
				st := strings.TrimSpace(lines[j].text)
				if st == "" {
					continue
				}
				ind := indentWidth(lines[j].raw)
				if ind < tryIndent {
					break
				}
				if ind == tryIndent && (st == "except:" || strings.HasPrefix(st, "except ") || strings.HasPrefix(st, "except:")) {
					if isExpectFailureExcept(st) && !suiteHandlesFailureForTesting(lines, j) {
						// Miss handlers that re-raise / record / log the failure —
						// those are not assertRaises substitutes.
						pushAt(unit, meta, lines[i].byte,
							"test uses try/except to expect failure; prefer with pytest.raises(...) or self.assertRaises(...)",
							out)
						return
					}
					break
				}
				if ind == tryIndent && (strings.HasPrefix(st, "finally") || strings.HasPrefix(st, "else:")) {
					break
				}
			}
		}
	}
}

func testRegionCollectsThreadErrors(lines []codeLine, start, end int) bool {
	if start < 0 || end > len(lines) || start >= end {
		return false
	}
	hasThread, collects, asserts := false, false, false
	for _, line := range lines[start:end] {
		t := strings.TrimSpace(line.text)
		hasThread = hasThread || strings.Contains(t, "threading.Thread(")
		collects = collects || strings.Contains(t, ".append(")
		asserts = asserts || strings.HasPrefix(t, "assert ") || strings.Contains(t, ".assert")
	}
	return hasThread && collects && asserts
}

func isExpectFailureExcept(st string) bool {
	// except: / except Exception / except AssertionError / except Exception as e:
	st = strings.TrimSpace(st)
	if st == "except:" || strings.HasPrefix(st, "except:") {
		return true
	}
	if !strings.HasPrefix(st, "except ") {
		return false
	}
	rest := strings.TrimSpace(strings.TrimPrefix(st, "except "))
	if i := strings.Index(rest, ":"); i >= 0 {
		rest = strings.TrimSpace(rest[:i])
	}
	if i := strings.Index(rest, " as "); i >= 0 {
		rest = strings.TrimSpace(rest[:i])
	}
	return rest == assertionErrorType || rest == "Exception" || rest == "BaseException" ||
		strings.HasPrefix(rest, assertionErrorType)
}

// suiteHandlesFailureForTesting reports an except suite that re-raises, logs,
// or records the failure — such handlers are not assertRaises substitutes and
// keep BP-PY-41 quiet (requestSpeedTest httpx_test.py re-raise shape).
func suiteHandlesFailureForTesting(lines []codeLine, exceptIdx int) bool {
	suiteIdx, _, _ := exceptSuiteLineIdx(lines, exceptIdx)
	for _, i := range suiteIdx {
		t := strings.TrimSpace(lines[i].text)
		if strings.HasPrefix(t, "raise") || strings.Contains(t, "exc_info") ||
			strings.Contains(t, ".exception(") || strings.Contains(t, "set_exception(") ||
			strings.Contains(t, "_error_result(") || strings.Contains(t, ".error") ||
			strings.Contains(t, "log.") || strings.Contains(t, "logger.") {
			return true
		}
	}
	return false
}
