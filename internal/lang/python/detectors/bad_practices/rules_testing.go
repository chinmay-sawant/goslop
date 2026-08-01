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
		// Only analyze test_* in test files, or any test_* anywhere (placeholder risk).
		defIndent := indentWidth(line.raw)
		hasAssert := false
		hasCall := false
		for j := i + 1; j < len(lines); j++ {
			st := strings.TrimSpace(lines[j].text)
			if st == "" {
				continue
			}
			ind := indentWidth(lines[j].raw)
			if ind <= defIndent {
				break
			}
			// Nested def/class ends our interest for simple bodies? Continue scanning body.
			if strings.HasPrefix(st, "def ") || strings.HasPrefix(st, "async def ") || strings.HasPrefix(st, "class ") {
				// nested — still part of outer until dedent of outer; keep going only for deeper.
				if ind == defIndent+1 || ind <= defIndent {
					break
				}
			}
			if isTestAssertion(st, helpers) {
				hasAssert = true
				break
			}
			if looksLikeSideEffectCall(st) {
				hasCall = true
			}
		}
		if hasCall && !hasAssert {
			pushAt(unit, meta, line.byte,
				"test function appears to only perform side effects without assertions (heuristic/info); add assert or pytest.raises",
				out)
		}
	}
}

func isTestAssertion(st string, helpers map[string]struct{}) bool {
	if strings.HasPrefix(st, "assert ") || st == "assert" || strings.HasPrefix(st, "assert(") {
		return true
	}
	if strings.Contains(st, "pytest.raises") || strings.Contains(st, "pytest.warns") {
		return true
	}
	if strings.Contains(st, "self.assert") || strings.Contains(st, "self.fail(") {
		return true
	}
	if strings.Contains(st, "unittest.") && strings.Contains(st, "assert") {
		return true
	}
	// nose/pytest helpers
	if strings.HasPrefix(st, "raises(") || strings.Contains(st, " assert_") {
		return true
	}
	// Calls to assert_* / _assert_* helpers (e.g. self._assert_verapdf(...)).
	if name := callCalleeIdent(st); name != "" {
		if isAssertHelperName(name) {
			return true
		}
		if _, ok := helpers[name]; ok {
			return true
		}
	}
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
// unittest assertions or raise AssertionError (presence-check helpers).
func sameFileAssertionHelpers(lines []codeLine) map[string]struct{} {
	out := make(map[string]struct{})
	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		name, ok := defNameIfNotTest(t)
		if !ok {
			continue
		}
		defIndent := indentWidth(line.raw)
		for j := i + 1; j < len(lines); j++ {
			st := strings.TrimSpace(lines[j].text)
			if st == "" {
				continue
			}
			ind := indentWidth(lines[j].raw)
			if ind <= defIndent {
				break
			}
			if strings.Contains(st, "self.assert") || strings.Contains(st, "self.fail(") ||
				strings.Contains(st, "raise "+assertionErrorType) {
				out[name] = struct{}{}
				break
			}
		}
	}
	return out
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
	if len(tests) == 0 && isPythonTestFile(unit) {
		// Whole file is test module — scan all.
		tests = append(tests, region{start: 0, end: len(lines), indent: -1})
	}
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
					if isExpectFailureExcept(st) {
						// Miss if same region uses assertRaises/pytest.raises nearby (already preferred).
						// Flag this try/except pattern.
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
	switch rest {
	case "Exception", "BaseException", assertionErrorType, "Exception as e":
		return true
	}
	// bare multi? keep conservative: only broad / AssertionError
	if rest == assertionErrorType || rest == "Exception" || rest == "BaseException" {
		return true
	}
	return rest == assertionErrorType || strings.HasPrefix(rest, assertionErrorType)
}
