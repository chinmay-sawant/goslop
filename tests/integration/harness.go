// Package integration provides a small fixture-materialization harness for
// Phase 12 parity scaffolding. Full §12.4 parity baseline comparison is intentionally
// out of scope until Phases 7–11 land.
package integration

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/engine"
	"github.com/chinmay/goslop/internal/fixture"
	"github.com/chinmay/goslop/internal/rules"
)

// RepoRoot returns the repository root (directory that contains go.mod),
// resolved from this source file's location under tests/integration/.
func RepoRoot() (string, error) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("runtime.Caller failed")
	}
	// tests/integration/harness.go → repo root is ../..
	root := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		return "", fmt.Errorf("repo root %s: %w", root, err)
	}
	return root, nil
}

// FixturesRoot is tests/fixtures under the repo root.
func FixturesRoot() (string, error) {
	root, err := RepoRoot()
	if err != nil {
		return "", err
	}
	p := filepath.Join(root, "tests", "fixtures")
	if _, err := os.Stat(p); err != nil {
		return "", fmt.Errorf("fixtures root %s: %w", p, err)
	}
	return p, nil
}

// Case describes one fixture file and the expected rule-fire behavior.
type Case struct {
	// RelPath is relative to tests/fixtures/ (e.g. "go/taint/CWE-78-vulnerable.txt").
	RelPath string
	// RuleID is the rule expected to fire (or stay silent when ExpectFire is false).
	RuleID string
	// ExpectFire is true for vulnerable fixtures, false for safe ones.
	ExpectFire bool
	// Only restricts analysis to these rule IDs (defaults to [RuleID] when empty).
	Only []string
}

// Result is the analyzer output for one materialized fixture.
type Result struct {
	Case       Case
	SourcePath string
	Findings   []rules.Finding
	Fired      bool
}

// HasRule reports whether any finding carries ruleID.
func HasRule(findings []rules.Finding, ruleID string) bool {
	for _, f := range findings {
		if f.RuleID == ruleID {
			return true
		}
	}
	return false
}

// MaterializeAndScan writes a single .txt fixture into outRoot and runs the
// default registry analyzer on the resulting source tree.
func MaterializeAndScan(c Case, outRoot string) (Result, error) {
	fix, err := FixturesRoot()
	if err != nil {
		return Result{}, err
	}
	txt := filepath.Join(fix, filepath.FromSlash(c.RelPath))
	srcPath, err := fixture.MaterializeFixtureFile(txt, outRoot)
	if err != nil {
		return Result{}, fmt.Errorf("materialize %s: %w", c.RelPath, err)
	}

	only := c.Only
	if len(only) == 0 && c.RuleID != "" {
		only = []string{c.RuleID}
	}
	ctx := core.NewScanContext(core.ProfileAll, only, nil)
	// Single-file trees are small; one worker is enough.
	analyzer := engine.NewAnalyzerBuilder().
		Registry(engine.DefaultRegistry()).
		ScanContext(ctx).
		Workers(1).
		Build()

	// Scan the language directory that contains the materialized file so the
	// walk still finds *.go even if nested under go/.
	scanRoot := filepath.Dir(srcPath)
	res, err := analyzer.AnalyzePaths([]string{scanRoot})
	if err != nil {
		return Result{}, fmt.Errorf("analyze %s: %w", srcPath, err)
	}
	findings := res.Findings
	if findings == nil {
		findings = []rules.Finding{}
	}
	return Result{
		Case:       c,
		SourcePath: srcPath,
		Findings:   findings,
		Fired:      HasRule(findings, c.RuleID),
	}, nil
}

// SeedCases is a small, always-on subset used by the scaffolding suite.
// Expand as detectors land; do not treat this as the full fixture matrix.
func SeedCases() []Case {
	return []Case{
		{
			RelPath:    "go/taint/CWE-78-vulnerable.txt",
			RuleID:     "CWE-78",
			ExpectFire: true,
		},
		{
			RelPath:    "go/taint/CWE-78-safe.txt",
			RuleID:     "CWE-78",
			ExpectFire: false,
		},
		{
			RelPath:    "go/taint/CWE-89-vulnerable.txt",
			RuleID:     "CWE-89",
			ExpectFire: true,
		},
		{
			RelPath:    "go/taint/CWE-89-safe.txt",
			RuleID:     "CWE-89",
			ExpectFire: false,
		},
		// PERF-006 fixture maps to rule PERF-6 (zero-padded fixture stem).
		{
			RelPath:    "go/perf/PERF-006-vulnerable.txt",
			RuleID:     "PERF-6",
			ExpectFire: true,
		},
		{
			RelPath:    "go/perf/PERF-006-safe.txt",
			RuleID:     "PERF-6",
			ExpectFire: false,
		},
	}
}

// RuleIDsFromFindings returns a sorted, de-duplicated list of rule IDs (stable for tests).
func RuleIDsFromFindings(findings []rules.Finding) []string {
	seen := make(map[string]struct{}, len(findings))
	ids := make([]string, 0, len(findings))
	for _, f := range findings {
		if _, ok := seen[f.RuleID]; ok {
			continue
		}
		seen[f.RuleID] = struct{}{}
		ids = append(ids, f.RuleID)
	}
	// Insertion order of first occurrence is enough for small seed cases.
	return ids
}

// SummarizeFindings is a short debug string for failed assertions.
func SummarizeFindings(findings []rules.Finding) string {
	if len(findings) == 0 {
		return "(no findings)"
	}
	parts := make([]string, 0, len(findings))
	for _, f := range findings {
		parts = append(parts, fmt.Sprintf("%s@%s:%d", f.RuleID, f.File, f.Line))
	}
	return strings.Join(parts, ", ")
}
