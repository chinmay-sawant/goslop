package badpractices

import (
	"strings"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/rules"
)

func init() {
	RegisterRule("BP-16", detectBP16)
	RegisterRule("BP-17", detectBP17)
	RegisterRule("BP-18", detectBP18)
	RegisterRule("BP-19", detectBP19)
	RegisterRule("BP-20", detectBP20)
	RegisterRule("BP-21", detectBP21)
	RegisterRule("BP-22", detectBP22)
	RegisterRule("BP-23", detectBP23)
	RegisterRule("BP-24", detectBP24)
	RegisterRule("BP-25", detectBP25)
	RegisterRule("BP-161", detectBP161)
	RegisterRule("BP-162", detectBP162)
	RegisterRule("BP-163", detectBP163)
}

func detectBP16(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-16")
	if !isTestFile(unit) || !strings.Contains(unit.Source, "time.Sleep(") {
		return
	}
	for _, line := range codeLines(unit.Source) {
		if strings.Contains(line.text, "time.Sleep(") {
			pushAt(unit, meta, line.byte, "time.Sleep in a test is brittle; prefer deterministic synchronization", out)
		}
	}
}

func detectBP17(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-17")
	if !isTestFile(unit) {
		return
	}
	lines := codeLines(unit.Source)
	for i := 0; i < len(lines)-1; i++ {
		cur := strings.TrimSpace(lines[i].text)
		next := strings.TrimSpace(lines[i+1].text)
		if (strings.HasPrefix(cur, "t.Error(") || strings.HasPrefix(cur, "t.Errorf(")) &&
			(strings.HasPrefix(next, "t.Fatal(") || strings.HasPrefix(next, "t.Fatalf(")) {
			pushAt(unit, meta, lines[i].byte, "t.Error immediately followed by t.Fatal is redundant; use Fatal alone", out)
		}
	}
}

func detectBP18(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-18")
	if !isTestFile(unit) {
		return
	}
	// t.Error without return/t.FailNow in same if block — heuristic: t.Error then continues with use of failed value
	lines := codeLines(unit.Source)
	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		if !(strings.HasPrefix(t, "t.Error(") || strings.HasPrefix(t, "t.Errorf(")) {
			continue
		}
		// if next non-empty is not return/Fatal and not closing brace alone of empty block
		hasExit := false
		for j := i + 1; j < len(lines) && j < i+5; j++ {
			n := strings.TrimSpace(lines[j].text)
			if n == "}" {
				break
			}
			if strings.HasPrefix(n, "return") || strings.HasPrefix(n, "t.FailNow") || strings.HasPrefix(n, "t.Fatal") {
				hasExit = true
				break
			}
		}
		if !hasExit {
			// only flag when inside if err
			if i > 0 && strings.Contains(strings.TrimSpace(lines[i-1].text), "err") {
				pushAt(unit, meta, line.byte, "t.Error without early exit continues the test after a failure", out)
			}
		}
	}
}

func detectBP19(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-19")
	if !isTestFile(unit) {
		return
	}
	// helper funcs taking *testing.T without t.Helper()
	src := unit.Source
	if !strings.Contains(src, "*testing.T") {
		return
	}
	for _, line := range codeLines(src) {
		t := strings.TrimSpace(line.text)
		if strings.HasPrefix(t, "func ") && strings.Contains(t, "*testing.T") && !strings.HasPrefix(t, "func Test") {
			// check body for t.Helper
			// find function name
			if !strings.Contains(src, "t.Helper()") && !strings.Contains(src, "tb.Helper()") {
				pushAt(unit, meta, line.byte, "test helper should call t.Helper()", out)
				return
			}
		}
	}
}

func detectBP20(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-20")
	if !isTestFile(unit) {
		return
	}
	// table tests: for _, tc := range tests { without t.Run
	if strings.Contains(unit.Source, "for _, ") && strings.Contains(unit.Source, " range ") &&
		(strings.Contains(unit.Source, "tests") || strings.Contains(unit.Source, "cases") || strings.Contains(unit.Source, "tt ")) {
		if !strings.Contains(unit.Source, "t.Run(") {
			if pos := strings.Index(unit.Source, "for _, "); pos >= 0 {
				pushAt(unit, meta, pos, "table-driven test without t.Run; name subtests for isolation", out)
			}
		}
	}
}

func detectBP21(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-21")
	if !isTestFile(unit) {
		return
	}
	if strings.Contains(unit.Source, "t.Run(") && !strings.Contains(unit.Source, "t.Parallel()") {
		if pos := strings.Index(unit.Source, "t.Run("); pos >= 0 {
			pushAt(unit, meta, pos, "subtest could call t.Parallel() for independent cases", out)
		}
	}
}

func detectBP22(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-22")
	if !isTestFile(unit) || !strings.Contains(unit.Source, "func TestMain(") {
		return
	}
	if strings.Contains(unit.Source, "TestMain") && !strings.Contains(unit.Source, "os.Exit") {
		if pos := strings.Index(unit.Source, "func TestMain"); pos >= 0 {
			pushAt(unit, meta, pos, "TestMain should call os.Exit with m.Run()'s status", out)
		}
	}
}

func detectBP23(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-23")
	if !isTestFile(unit) {
		return
	}
	// long-looking tests: integration/slow markers without testing.Short
	src := unit.Source
	if (strings.Contains(src, "integration") || strings.Contains(src, "time.Sleep") || strings.Contains(src, "http.Get")) &&
		!strings.Contains(src, "testing.Short") {
		if pos := strings.Index(src, "func Test"); pos >= 0 {
			pushAt(unit, meta, pos, "long-running test lacks testing.Short guard", out)
		}
	}
}

func detectBP24(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-24")
	if !isTestFile(unit) {
		return
	}
	if !strings.Contains(unit.Source, "func Test") && !strings.Contains(unit.Source, "func Benchmark") && !strings.Contains(unit.Source, "func Fuzz") {
		pushAt(unit, meta, 0, "test file contains no Test/Benchmark/Fuzz functions", out)
	}
}

func detectBP25(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-25")
	if !isTestFile(unit) {
		return
	}
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if strings.HasPrefix(t, "func ") && strings.Contains(t, "*testing.T") && strings.Contains(t, " error") &&
			!strings.HasPrefix(t, "func Test") {
			pushAt(unit, meta, line.byte, "test helper should fail the test rather than returning an error", out)
		}
	}
}

func detectBP161(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-161")
	if !isTestFile(unit) {
		return
	}
	src := strings.ToLower(unit.Source)
	prodHints := []string{"prod", "production", "live.", "rds.amazonaws", "azure.com"}
	if strings.Contains(src, "sql.open") || strings.Contains(src, "postgres://") || strings.Contains(src, "mysql://") {
		for _, h := range prodHints {
			if strings.Contains(src, h) {
				pushAt(unit, meta, 0, "test uses a production-looking DSN; point tests at local/ephemeral databases", out)
				return
			}
		}
	}
}

func detectBP162(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-162")
	if !isTestFile(unit) || !strings.Contains(unit.Source, "t.Parallel()") {
		return
	}
	// package-level var mutated in parallel test
	if strings.Contains(unit.Source, "var ") && strings.Contains(unit.Source, "t.Parallel()") {
		// look for assignment to package level — weak
		for _, line := range codeLines(unit.Source) {
			t := strings.TrimSpace(line.text)
			if strings.Contains(t, "=") && !strings.Contains(t, ":=") && !strings.HasPrefix(t, "var ") &&
				!strings.HasPrefix(t, "const ") && !strings.HasPrefix(t, "type ") &&
				!strings.Contains(t, "==") && !strings.Contains(t, "!=") {
				// skip local-looking
				if strings.HasPrefix(t, "t.") || strings.HasPrefix(t, "err ") {
					continue
				}
			}
		}
		// if global var exists and Parallel
		hasGlobal := false
		for _, line := range codeLines(unit.Source) {
			t := strings.TrimSpace(line.text)
			if strings.HasPrefix(t, "var ") && !strings.Contains(t, "func") {
				// at package level roughly if previous was not inside func — accept any top var
				hasGlobal = true
			}
		}
		if hasGlobal && strings.Contains(unit.Source, "t.Parallel()") {
			if pos := strings.Index(unit.Source, "t.Parallel()"); pos >= 0 {
				pushAt(unit, meta, pos, "parallel test may mutate shared package state", out)
			}
		}
	}
}

func detectBP163(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-163")
	if !isTestFile(unit) {
		return
	}
	if (strings.Contains(unit.Source, "update") || strings.Contains(unit.Source, "Update")) &&
		(strings.Contains(unit.Source, "golden") || strings.Contains(unit.Source, "os.WriteFile") || strings.Contains(unit.Source, "ioutil.WriteFile")) {
		if !strings.Contains(unit.Source, "testing.Short") && !strings.Contains(unit.Source, "*update") {
			if pos := strings.Index(unit.Source, "WriteFile"); pos >= 0 {
				pushAt(unit, meta, pos, "golden update path lacks a short-test or flag guard", out)
			}
		}
	}
}
