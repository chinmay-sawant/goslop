package badpractices

import (
	goast "go/ast"
	"strings"

	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/rules"
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

func detectBP16(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-16")
	if !isTestFile(unit) || !strings.Contains(unit.Source, "time.Sleep(") {
		return
	}
	for _, line := range codeLinesFacts(facts, unit.Source) {
		if strings.Contains(line.text, "time.Sleep(") {
			pushAt(unit, meta, line.byte, "time.Sleep in a test is brittle; prefer deterministic synchronization", out)
		}
	}
}

func detectBP17(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-17")
	if !isTestFile(unit) {
		return
	}
	lines := codeLinesFacts(facts, unit.Source)
	for i := 0; i < len(lines)-1; i++ {
		cur := strings.TrimSpace(lines[i].text)
		next := strings.TrimSpace(lines[i+1].text)
		if (strings.HasPrefix(cur, "t.Error(") || strings.HasPrefix(cur, "t.Errorf(")) &&
			(strings.HasPrefix(next, "t.Fatal(") || strings.HasPrefix(next, "t.Fatalf(")) {
			pushAt(unit, meta, lines[i].byte, "t.Error immediately followed by t.Fatal is redundant; use Fatal alone", out)
		}
	}
}

func detectBP18(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-18")
	if !isTestFile(unit) {
		return
	}
	// Rust parity: t.Error/f without an immediate terminating next statement.
	lines := codeLinesFacts(facts, unit.Source)
	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		if !(strings.HasPrefix(t, "t.Error(") || strings.HasPrefix(t, "t.Errorf(")) {
			continue
		}
		nextIdx := i + 1
		for nextIdx < len(lines) && strings.TrimSpace(lines[nextIdx].text) == "" {
			nextIdx++
		}
		if nextIdx >= len(lines) {
			continue
		}
		n := strings.TrimSpace(lines[nextIdx].text)
		terminates := n == "return" ||
			strings.HasPrefix(n, "return ") ||
			strings.HasPrefix(n, "t.FailNow(") ||
			strings.HasPrefix(n, "t.Fatal(") ||
			strings.HasPrefix(n, "t.Fatalf(") ||
			strings.HasPrefix(n, "t.Skip(") ||
			strings.HasPrefix(n, "t.Skipf(") ||
			strings.HasPrefix(n, "t.SkipNow(")
		if !terminates {
			pushAt(unit, meta, line.byte, "t.Error continues the test path; return or fail immediately after the error", out)
		}
	}
}

func detectBP19(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-19")
	if !isTestFile(unit) {
		return
	}
	// helper funcs taking *testing.T without t.Helper()
	src := unit.Source
	if !strings.Contains(src, "*testing.T") {
		return
	}
	for _, line := range codeLinesFacts(facts, src) {
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

func detectBP20(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
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

func detectBP21(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
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

func detectBP22(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
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

func detectBP23(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-23")
	if !isTestFile(unit) {
		return
	}
	// Rust parity: Test* body spanning ≥20 lines without testing.Short().
	msg := "long-running test should gate itself with testing.Short()"
	if facts != nil && facts.tree != nil && facts.tree.File != nil {
		tree := facts.tree
		for _, decl := range tree.File.Decls {
			fd, ok := decl.(*goast.FuncDecl)
			if !ok || fd.Name == nil || fd.Body == nil {
				continue
			}
			if !strings.HasPrefix(fd.Name.Name, "Test") {
				continue
			}
			startLine := tree.Fset.Position(fd.Pos()).Line
			endLine := tree.Fset.Position(fd.End()).Line
			if endLine-startLine < 20 {
				continue
			}
			body := tree.NodeText(fd.Body)
			if strings.Contains(body, "testing.Short") {
				continue
			}
			pushAt(unit, meta, tree.Offset(fd.Pos()), msg, out)
		}
		return
	}
	// Text fallback: any Test* whose source span is long without Short.
	src := unit.Source
	if strings.Contains(src, "testing.Short") {
		return
	}
	lines := strings.Split(src, "\n")
	if len(lines) < 20 {
		return
	}
	if pos := strings.Index(src, "func Test"); pos >= 0 {
		pushAt(unit, meta, pos, msg, out)
	}
}

func detectBP24(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-24")
	if !isTestFile(unit) {
		return
	}
	if !strings.Contains(unit.Source, "func Test") && !strings.Contains(unit.Source, "func Benchmark") && !strings.Contains(unit.Source, "func Fuzz") {
		pushAt(unit, meta, 0, "test file contains no Test/Benchmark/Fuzz functions", out)
	}
}

func detectBP25(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-25")
	if !isTestFile(unit) {
		return
	}
	for _, line := range codeLinesFacts(facts, unit.Source) {
		t := strings.TrimSpace(line.text)
		if strings.HasPrefix(t, "func ") && strings.Contains(t, "*testing.T") && strings.Contains(t, " error") &&
			!strings.HasPrefix(t, "func Test") {
			pushAt(unit, meta, line.byte, "test helper should fail the test rather than returning an error", out)
		}
	}
}

func detectBP161(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-161")
	if !isTestFile(unit) {
		return
	}
	src := unit.Source
	// sql.Open or gorm.Open with a literal production marker in the call text.
	for _, line := range codeLinesFacts(facts, src) {
		t := line.text
		if !(strings.Contains(t, "sql.Open") || strings.Contains(t, "gorm.Open") ||
			strings.Contains(t, "postgres.Open") || strings.Contains(t, "mysql.")) {
			// gorm may span lines: collect multi-line call roughly via full source later
			continue
		}
		if containsLiteralProductionMarker(t) {
			pushAt(unit, meta, line.byte, "test uses a production-looking DSN; point tests at local/ephemeral databases", out)
			return
		}
	}
	// Multi-line Open calls: check a window around Open.
	if strings.Contains(src, "sql.Open") || strings.Contains(src, "gorm.Open") || strings.Contains(src, "postgres.Open") {
		if containsLiteralProductionMarker(src) &&
			(strings.Contains(src, "sql.Open") || strings.Contains(src, "gorm.Open") || strings.Contains(src, "postgres.Open")) {
			// Avoid flagging local targets only when no prod marker — already gated.
			// But exclude pure local: if only localhost/127.0.0.1 and no prod, skip.
			if containsLiteralProductionMarker(src) {
				pos := strings.Index(src, "Open(")
				if pos < 0 {
					pos = 0
				}
				pushAt(unit, meta, pos, "test uses a production-looking DSN; point tests at local/ephemeral databases", out)
			}
		}
	}
}

func containsLiteralProductionMarker(s string) bool {
	lower := strings.ToLower(s)
	// Standalone prod / production markers in a string literal-ish context.
	markers := []string{
		"prod-", "-prod", ".prod.", "production", "orders-prod",
		"prod-db", "prod_db", "prod.", "live.",
		"rds.amazonaws", "azure.com",
	}
	for _, m := range markers {
		if strings.Contains(lower, m) {
			return true
		}
	}
	// host=...prod...
	if strings.Contains(lower, "host=") && strings.Contains(lower, "prod") {
		return true
	}
	// postgres://...@prod...
	if strings.Contains(lower, "postgres://") || strings.Contains(lower, "mysql://") {
		if strings.Contains(lower, "prod") && !strings.Contains(lower, "localhost") {
			return true
		}
	}
	return false
}

func detectBP162(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-162")
	if !isTestFile(unit) || !strings.Contains(unit.Source, "t.Parallel()") {
		return
	}
	src := unit.Source
	globals := packageLevelVarNames(src)
	if len(globals) == 0 {
		return
	}
	// Only fire when a parallel test body assigns to a package-level name.
	lines := codeLinesFacts(facts, src)
	inParallelTest := false
	depth := 0
	testDepth := 0
	for _, line := range lines {
		t := strings.TrimSpace(line.text)
		if strings.HasPrefix(t, "func Test") || strings.HasPrefix(t, "func Benchmark") || strings.HasPrefix(t, "func Fuzz") {
			inParallelTest = false
			testDepth = 0
			depth = 0
		}
		for _, ch := range line.text {
			if ch == '{' {
				depth++
				if testDepth == 0 && (strings.HasPrefix(t, "func Test") || strings.Contains(line.text, "func Test")) {
					testDepth = depth
				}
			} else if ch == '}' {
				if depth > 0 {
					depth--
				}
				if testDepth > 0 && depth < testDepth {
					inParallelTest = false
					testDepth = 0
				}
			}
		}
		if strings.Contains(t, "t.Parallel()") {
			inParallelTest = true
		}
		if !inParallelTest {
			continue
		}
		// assignment or index assign to global
		if !strings.Contains(t, "=") && !strings.HasSuffix(t, "++") && !strings.HasSuffix(t, "--") {
			continue
		}
		if strings.Contains(t, ":=") || strings.Contains(t, "==") || strings.Contains(t, "!=") {
			continue
		}
		if strings.HasPrefix(t, "t.") || strings.HasPrefix(t, "var ") || strings.HasPrefix(t, "const ") {
			continue
		}
		for g := range globals {
			if strings.HasPrefix(t, g+" ") || strings.HasPrefix(t, g+"=") ||
				strings.HasPrefix(t, g+".") || strings.HasPrefix(t, g+"[") ||
				strings.HasPrefix(t, g+"++") || strings.HasPrefix(t, g+"--") {
				pushAt(unit, meta, line.byte, "parallel test mutates package-level state; use an isolated per-test fixture", out)
				return
			}
		}
	}
}

func detectBP163(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-163")
	if !isTestFile(unit) {
		return
	}
	src := unit.Source
	if strings.Contains(src, "testing.Short") {
		return
	}
	// Require an update flag declaration.
	hasUpdateFlag := strings.Contains(src, `flag.Bool("update"`) ||
		strings.Contains(src, `flag.Bool("update-golden"`) ||
		(strings.Contains(src, "flag.BoolVar(") && strings.Contains(src, `"update"`))
	if !hasUpdateFlag {
		return
	}
	// Update branch with WriteFile or Create.
	lines := codeLinesFacts(facts, src)
	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		// if *update / if *updateGolden
		if !strings.HasPrefix(t, "if ") {
			continue
		}
		cond := strings.TrimPrefix(t, "if ")
		condNorm := strings.Map(func(r rune) rune {
			if r == '*' || r == ' ' || r == '\t' {
				return -1
			}
			return r
		}, cond)
		// strip trailing {
		condNorm = strings.TrimSuffix(condNorm, "{")
		if !(condNorm == "update" || condNorm == "updateGolden" ||
			strings.HasPrefix(condNorm, "update&&") || strings.HasPrefix(condNorm, "updateGolden&&")) {
			// also allow contains update as sole condition pieces
			if !strings.Contains(cond, "update") && !strings.Contains(cond, "Update") {
				continue
			}
			// require *update-style
			if !strings.Contains(cond, "*update") && !strings.Contains(cond, "*updateGolden") {
				continue
			}
		}
		// Scan consequence for WriteFile/Create
		depth := 0
		started := false
		for j := i; j < len(lines); j++ {
			lt := lines[j].text
			for _, ch := range lt {
				if ch == '{' {
					depth++
					started = true
				} else if ch == '}' {
					if depth > 0 {
						depth--
					}
				}
			}
			bt := strings.TrimSpace(lt)
			if strings.Contains(bt, "os.WriteFile") || strings.Contains(bt, "ioutil.WriteFile") ||
				strings.Contains(bt, "os.Create(") {
				pushAt(unit, meta, lines[j].byte, "golden-file update path writes without a testing.Short() guard", out)
				return
			}
			if started && depth == 0 && j > i {
				break
			}
		}
	}
}
