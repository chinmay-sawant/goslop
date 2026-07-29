package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// DiscoverCasesWithSuffix returns case stems for files ending in suffix
// (e.g. "-vulnerable.txt" → "BP-1", "PERF-038-done").
func DiscoverCasesWithSuffix(dir, suffix string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read_dir %s: %w", dir, err)
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if stem, ok := strings.CutSuffix(name, suffix); ok {
			out = append(out, stem)
		}
	}
	sort.Strings(out)
	return out, nil
}

// DiscoverPairedCases requires matching -vulnerable/-safe stems under dir.
func DiscoverPairedCases(dir string) ([]string, error) {
	vuln, err := DiscoverCasesWithSuffix(dir, "-vulnerable.txt")
	if err != nil {
		return nil, err
	}
	safe, err := DiscoverCasesWithSuffix(dir, "-safe.txt")
	if err != nil {
		return nil, err
	}
	if len(vuln) != len(safe) {
		return nil, fmt.Errorf("%s: vulnerable/safe set sizes differ: %d vs %d", dir, len(vuln), len(safe))
	}
	for i := range vuln {
		if vuln[i] != safe[i] {
			return nil, fmt.Errorf("%s: vulnerable/safe sets drifted at %q vs %q", dir, vuln[i], safe[i])
		}
	}
	return vuln, nil
}

// DiscoverBPCases lists BP fixture cases under tests/fixtures/go/bad_practices.
func DiscoverBPCases() ([]string, error) {
	fx, err := FixturesRoot()
	if err != nil {
		return nil, err
	}
	cases, err := DiscoverPairedCases(filepath.Join(fx, "go", "bad_practices"))
	if err != nil {
		return nil, err
	}
	sort.Slice(cases, func(i, j int) bool {
		return compareBPCase(cases[i], cases[j]) < 0
	})
	return cases, nil
}

// BPRuleID maps case stem BP-1 / BP-1-variant → BP-1.
func BPRuleID(caseName string) string {
	// BP-<num>[-variant...]
	parts := strings.Split(caseName, "-")
	if len(parts) < 2 || parts[0] != "BP" {
		return caseName
	}
	return "BP-" + parts[1]
}

// BPFixtureRel returns path relative to tests/fixtures/.
func BPFixtureRel(caseName string, vulnerable bool) string {
	suf := "safe"
	if vulnerable {
		suf = "vulnerable"
	}
	return filepath.ToSlash(filepath.Join("go", "bad_practices", caseName+"-"+suf+".txt"))
}

func compareBPCase(a, b string) int {
	an, at := parseBPCase(a)
	bn, bt := parseBPCase(b)
	if an != bn {
		return an - bn
	}
	return strings.Compare(at, bt)
}

func parseBPCase(s string) (int, string) {
	rest, ok := strings.CutPrefix(s, "BP-")
	if !ok {
		return 0, s
	}
	numStr, tail, found := strings.Cut(rest, "-")
	if !found {
		numStr, tail = rest, ""
	}
	n, _ := strconv.Atoi(numStr)
	return n, tail
}

// DiscoverPERFCases lists PERF fixture cases under tests/fixtures/go/perf.
func DiscoverPERFCases() ([]string, error) {
	fx, err := FixturesRoot()
	if err != nil {
		return nil, err
	}
	cases, err := DiscoverPairedCases(filepath.Join(fx, "go", "perf"))
	if err != nil {
		return nil, err
	}
	sort.Slice(cases, func(i, j int) bool {
		return comparePERFCase(cases[i], cases[j]) < 0
	})
	return cases, nil
}

// PERFRuleID maps PERF-038 / PERF-038-done → PERF-38.
func PERFRuleID(caseName string) string {
	rest, ok := strings.CutPrefix(caseName, "PERF-")
	if !ok {
		return caseName
	}
	numStr, _, _ := strings.Cut(rest, "-")
	n, err := strconv.Atoi(numStr)
	if err != nil {
		return caseName
	}
	return fmt.Sprintf("PERF-%d", n)
}

// PERFFixtureRel returns path relative to tests/fixtures/.
func PERFFixtureRel(caseName string, vulnerable bool) string {
	suf := "safe"
	if vulnerable {
		suf = "vulnerable"
	}
	return filepath.ToSlash(filepath.Join("go", "perf", caseName+"-"+suf+".txt"))
}

func comparePERFCase(a, b string) int {
	an, at := parsePERFCase(a)
	bn, bt := parsePERFCase(b)
	if an != bn {
		return an - bn
	}
	return strings.Compare(at, bt)
}

func parsePERFCase(s string) (int, string) {
	rest, ok := strings.CutPrefix(s, "PERF-")
	if !ok {
		return 0, s
	}
	numStr, tail, found := strings.Cut(rest, "-")
	if !found {
		numStr, tail = rest, ""
	}
	n, _ := strconv.Atoi(numStr)
	return n, tail
}

// DiscoverCWECases lists CWE ids present in both frameworks/ and stdlib/ suites.
func DiscoverCWECases() ([]string, error) {
	fx, err := FixturesRoot()
	if err != nil {
		return nil, err
	}
	fw, err := DiscoverPairedCases(filepath.Join(fx, "go", "frameworks"))
	if err != nil {
		return nil, err
	}
	st, err := DiscoverPairedCases(filepath.Join(fx, "go", "stdlib"))
	if err != nil {
		return nil, err
	}
	if len(fw) != len(st) {
		return nil, fmt.Errorf("frameworks/stdlib CWE inventory sizes differ: %d vs %d", len(fw), len(st))
	}
	for i := range fw {
		if fw[i] != st[i] {
			return nil, fmt.Errorf("frameworks/stdlib CWE inventory drifted: %q vs %q", fw[i], st[i])
		}
	}
	sort.Slice(fw, func(i, j int) bool {
		return parseCWENumber(fw[i]) < parseCWENumber(fw[j])
	})
	return fw, nil
}

// CWEFixtureRel returns path relative to tests/fixtures/ for suite in {frameworks,stdlib}.
func CWEFixtureRel(suite, cwe string, vulnerable bool) string {
	suf := "safe"
	if vulnerable {
		suf = "vulnerable"
	}
	return filepath.ToSlash(filepath.Join("go", suite, cwe+"-"+suf+".txt"))
}

func parseCWENumber(cwe string) int {
	rest, ok := strings.CutPrefix(cwe, "CWE-")
	if !ok {
		return 0
	}
	n, _ := strconv.Atoi(rest)
	return n
}

// DiscoverTaintCases lists CWE-* stems under tests/fixtures/go/taint.
func DiscoverTaintCases() ([]string, error) {
	fx, err := FixturesRoot()
	if err != nil {
		return nil, err
	}
	return DiscoverPairedCases(filepath.Join(fx, "go", "taint"))
}

// TaintFixtureRel returns path relative to tests/fixtures/.
func TaintFixtureRel(caseName string, vulnerable bool) string {
	suf := "safe"
	if vulnerable {
		suf = "vulnerable"
	}
	return filepath.ToSlash(filepath.Join("go", "taint", caseName+"-"+suf+".txt"))
}
