package perf

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/fixture"
)

func TestBuildCodeLinesBlanksTripleQuotedBodies(t *testing.T) {
	t.Parallel()
	src := loadPerfFactsFixture(t, "PERF-PY-26-code-lines-docstring.txt")
	lines := buildCodeLines(src)
	joined := ""
	for _, line := range lines {
		joined += line.text
	}
	if containsFold(joined, ".objects.filter") {
		t.Fatalf("docstring ORM call should be blanked, got:\n%q", joined)
	}
}

func TestBuildCodeLinesKeepsOrdinaryStringKeywords(t *testing.T) {
	t.Parallel()
	src := loadPerfFactsFixture(t, "PERF-PY-26-code-lines-string-keyword.txt")
	lines := buildCodeLines(src)
	if len(lines) == 0 || !containsFold(lines[0].text, "pending") {
		t.Fatalf("ordinary string keyword should remain for heuristics: %#v", lines)
	}
}

func TestComputeInLoopMatchesScan(t *testing.T) {
	t.Parallel()
	assertInLoopConsistent(t, loadPerfFactsFixture(t, "PERF-PY-26-inloop-nested.txt"))
}

func TestComputeInLoopStopsAtFunctionBoundary(t *testing.T) {
	t.Parallel()
	src := loadPerfFactsFixture(t, "PERF-PY-26-cli-parse-safe.txt")
	lines := buildCodeLines(src)
	got := computeInLoop(lines)
	assertInLoopConsistent(t, src)
	for i, line := range lines {
		if strings.Contains(line.text, "parse_report(xml_data)") && got[i] {
			t.Fatalf("main() parse_report call must not inherit parse_report's loop: line=%q", line.raw)
		}
	}
}

func assertInLoopConsistent(t *testing.T, src string) {
	t.Helper()
	lines := buildCodeLines(src)
	got := computeInLoop(lines)
	if len(got) != len(lines) {
		t.Fatalf("len mismatch: got %d want %d", len(got), len(lines))
	}
	for i := range lines {
		if got[i] != inLoop(lines, i) {
			t.Fatalf("inLoop[%d]=%v want %v\nline=%q", i, got[i], inLoop(lines, i), lines[i].raw)
		}
	}
}

func loadPerfFactsFixture(t *testing.T, name string) string {
	t.Helper()
	txtPath := filepath.Join(perfFixturesRoot(t), name)
	data, err := os.ReadFile(txtPath)
	if err != nil {
		t.Fatalf("read fixture %s: %v", txtPath, err)
	}
	fx, err := fixture.ParseFixture(string(data), txtPath)
	if err != nil {
		t.Fatalf("parse fixture %s: %v", txtPath, err)
	}
	return fx.Source
}

func perfFixturesRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "..", ".."))
	fx := filepath.Join(root, "tests", "fixtures", "python", "perf")
	if _, err := os.Stat(fx); err != nil {
		t.Fatalf("fixtures root %s: %v", fx, err)
	}
	return fx
}

func containsFold(s, needle string) bool {
	return len(s) >= len(needle) && (func() bool {
		for i := 0; i+len(needle) <= len(s); i++ {
			if s[i:i+len(needle)] == needle {
				return true
			}
		}
		return false
	})()
}
