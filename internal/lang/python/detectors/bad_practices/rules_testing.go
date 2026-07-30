package badpractices

import (
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

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
			if isTestAssertion(st) {
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

func isTestAssertion(st string) bool {
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
	return false
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
	var tests []region
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
	case "Exception", "BaseException", "AssertionError", "Exception as e":
		return true
	}
	// bare multi? keep conservative: only broad / AssertionError
	if rest == "AssertionError" || rest == "Exception" || rest == "BaseException" {
		return true
	}
	return rest == "AssertionError" || strings.HasPrefix(rest, "AssertionError")
}
