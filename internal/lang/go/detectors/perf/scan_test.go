package perf_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/fixture"
	"github.com/chinmay/goslop/internal/lang/go/detectors/perf"
	"github.com/chinmay/goslop/internal/rules"
)

func TestPERF116Vulnerable(t *testing.T) {
	src := `package sample

import "strings"

func Has(s, sub string) bool {
	return strings.Index(s, sub) != -1
}
`
	findings := runPerf(t, src, "PERF-116-vulnerable.go")
	if !hasRule(findings, "PERF-116") {
		t.Fatalf("expected PERF-116 finding, got %#v", findings)
	}
}

func TestPERF116Safe(t *testing.T) {
	src := `package sample

import "strings"

func Has(s, sub string) bool {
	return strings.Contains(s, sub)
}
`
	findings := runPerf(t, src, "PERF-116-safe.go")
	if hasRule(findings, "PERF-116") {
		t.Fatalf("safe fixture should not emit PERF-116, got %#v", findings)
	}
}

func TestPERF116IndexWithoutCompare(t *testing.T) {
	src := `package sample
import "strings"
func Pos(s, sub string) int { return strings.Index(s, sub) }
`
	findings := runPerf(t, src, "index-only.go")
	if hasRule(findings, "PERF-116") {
		t.Fatalf("Index without -1 compare should be silent, got %#v", findings)
	}
}

func TestPERF6FmtInLoop(t *testing.T) {
	src := `package sample
import "fmt"
func FormatAll(values []int) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		out = append(out, fmt.Sprintf("value=%d", v))
	}
	return out
}
`
	findings := runPerf(t, src, "PERF-006-vulnerable.go")
	if !hasRule(findings, "PERF-6") {
		t.Fatalf("expected PERF-6, got %#v", findings)
	}
}

func TestPERF6SafeNoLoop(t *testing.T) {
	src := `package sample
import (
	"bytes"
	"strconv"
)
func FormatOne(buf *bytes.Buffer, v int) string {
	buf.Reset()
	buf.WriteString("value=")
	buf.WriteString(strconv.Itoa(v))
	return buf.String()
}
`
	findings := runPerf(t, src, "PERF-006-safe.go")
	if hasRule(findings, "PERF-6") {
		t.Fatalf("safe should not emit PERF-6, got %#v", findings)
	}
}

func TestPERF7DeferInLoop(t *testing.T) {
	src := `package sample
import "os"
func CloseAll(paths []string) error {
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()
	}
	return nil
}
`
	findings := runPerf(t, src, "PERF-007-vulnerable.go")
	if !hasRule(findings, "PERF-7") {
		t.Fatalf("expected PERF-7, got %#v", findings)
	}
}

func TestPERF8TimeParseInLoop(t *testing.T) {
	src := `package sample
import (
	"strings"
	"time"
)
func ParseTimestamps(input string) []time.Time {
	out := make([]time.Time, 0)
	for _, line := range strings.Split(input, "\n") {
		t, _ := time.Parse("2006-01-02", line)
		out = append(out, t)
	}
	return out
}
`
	findings := runPerf(t, src, "PERF-008-vulnerable.go")
	if !hasRule(findings, "PERF-8") {
		t.Fatalf("expected PERF-8, got %#v", findings)
	}
}

func TestPERF32ByteConversionInLoop(t *testing.T) {
	src := `package sample
func EncodeParts(parts []string) [][]byte {
	out := make([][]byte, 0, len(parts))
	for _, v := range parts {
		out = append(out, []byte(v))
	}
	return out
}
`
	findings := runPerf(t, src, "PERF-032-vulnerable.go")
	if !hasRule(findings, "PERF-32") {
		t.Fatalf("expected PERF-32, got %#v", findings)
	}
}

func TestPERF32LiteralSafe(t *testing.T) {
	src := `package sample
func EncodeOK() []byte {
	return []byte("ok")
}
`
	findings := runPerf(t, src, "PERF-032-safe.go")
	if hasRule(findings, "PERF-32") {
		t.Fatalf("literal conversion should be silent, got %#v", findings)
	}
}

func TestPERF50RegexpMatchInLoop(t *testing.T) {
	src := `package sample
import "regexp"
func FindHits(rows []string) []bool {
	hits := make([]bool, 0, len(rows))
	for _, r := range rows {
		matched, _ := regexp.MatchString(` + "`^[A-Z]+$`" + `, r)
		hits = append(hits, matched)
	}
	return hits
}
`
	findings := runPerf(t, src, "PERF-050-vulnerable.go")
	if !hasRule(findings, "PERF-50") {
		t.Fatalf("expected PERF-50, got %#v", findings)
	}
}

func TestPERF230StableHelperInLoop(t *testing.T) {
	src := `package sample
func parseProps(props string) int { return len(props) }
func EstimateTextWidth(font string, text string, size float64) float64 {
	return float64(len(text)) * size * 0.5
}
func GetTextWidth(font string, text string) float64 { return float64(len(text)) }
func resolveFontName(props string) string { return props }
func DrawCells(cells []string, props string, font string, size float64) float64 {
	total := 0.0
	for _, cell := range cells {
		w := parseProps(props)
		_ = resolveFontName(props)
		tw := EstimateTextWidth(font, cell, size)
		gw := GetTextWidth(font, cell)
		total += float64(w) + tw + gw
	}
	return total
}
`
	findings := runPerf(t, src, "PERF-230-vulnerable.go")
	if !hasRule(findings, "PERF-230") {
		t.Fatalf("expected PERF-230, got %#v", findings)
	}
}

func TestPERF1MustCompileInLoop(t *testing.T) {
	src := `package sample
import "regexp"
func MatchAll(patterns, rows []string) int {
	n := 0
	for _, p := range patterns {
		re := regexp.MustCompile(p)
		for _, r := range rows {
			if re.MatchString(r) {
				n++
			}
		}
	}
	return n
}
`
	findings := runPerf(t, src, "PERF-001-vulnerable.go")
	if !hasRule(findings, "PERF-1") {
		t.Fatalf("expected PERF-1, got %#v", findings)
	}
}

func TestGoPerfScanRuleIDs(t *testing.T) {
	d := perf.NewGoPerfScan()
	ids := d.RuleIDs()
	// Full registry catalogue (all domains) must be registered.
	if len(ids) < 239 {
		t.Fatalf("expected at least 239 PERF rules, got %d", len(ids))
	}
	for _, want := range []string{"PERF-1", "PERF-6", "PERF-32", "PERF-116", "PERF-230", "PERF-61", "PERF-242"} {
		found := false
		for _, id := range ids {
			if id == want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing rule id %s", want)
		}
	}
}

func TestFixturePERF006(t *testing.T) {
	runFixtureRule(t, "PERF-006-vulnerable.txt", "PERF-6", true)
	runFixtureRule(t, "PERF-006-safe.txt", "PERF-6", false)
}

func TestFixturePERF032(t *testing.T) {
	runFixtureRule(t, "PERF-032-vulnerable.txt", "PERF-32", true)
	runFixtureRule(t, "PERF-032-safe.txt", "PERF-32", false)
}

func TestFixturePERF230(t *testing.T) {
	runFixtureRule(t, "PERF-230-vulnerable.txt", "PERF-230", true)
}

func runFixtureRule(t *testing.T, name, rule string, wantFire bool) {
	t.Helper()
	cwd, _ := os.Getwd()
	var (
		data []byte
		err  error
		path string
	)
	// go test runs with package dir as cwd (…/internal/lang/go/detectors/perf).
	candidates := []string{
		filepath.Join(cwd, "..", "..", "..", "..", "..", "tests", "fixtures", "go", "perf", name),
		filepath.Join(cwd, "tests", "fixtures", "go", "perf", name),
	}
	for d := cwd; d != string(filepath.Separator) && d != "."; d = filepath.Dir(d) {
		candidates = append(candidates, filepath.Join(d, "tests", "fixtures", "go", "perf", name))
		if filepath.Dir(d) == d {
			break
		}
	}
	for _, c := range candidates {
		data, err = os.ReadFile(c)
		if err == nil {
			path = c
			break
		}
	}
	if err != nil {
		t.Fatalf("read fixture %s: %v (cwd=%s)", name, err, cwd)
	}
	fx, err := fixture.ParseFixture(string(data), path)
	if err != nil {
		t.Fatal(err)
	}
	findings := runPerf(t, fx.Source, fx.Filename)
	got := hasRule(findings, rule)
	if got != wantFire {
		t.Fatalf("fixture %s rule %s want fire=%v got=%v findings=%v", name, rule, wantFire, got, findings)
	}
}

func runPerf(t *testing.T, src, path string) []rules.Finding {
	t.Helper()
	d := perf.NewGoPerfScan()
	unit := core.NewParsedUnit(core.LangGo, path, src)
	ctx := core.DefaultScanContext()
	var out []rules.Finding
	d.Run(ctx, unit, &out)
	return out
}

func hasRule(findings []rules.Finding, id string) bool {
	for _, f := range findings {
		if f.RuleID == id {
			return true
		}
	}
	return false
}
