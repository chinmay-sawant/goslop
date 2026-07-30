package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/engine"
	"github.com/chinmay-sawant/goslop/internal/fixture"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// ScanOptions configures a fixture matrix scan (Rust Analyzer::builder defaults).
type ScanOptions struct {
	// IncludeTests matches Rust BP path_filters.exclude_tests = false.
	IncludeTests bool
	// TaintEnabled enables Phase 9 taint (needed for CWE-22/78/79/89 matrix).
	TaintEnabled bool
	// Profile defaults to all.
	Profile core.ScanProfile
	// Only restricts detectors (empty = full catalogue, Rust matrix style).
	// When set to the target rule, safe fixtures only need that rule silent
	// (avoids cross-rule BP-41/PERF noise that Rust suppresses via flat-fixture gates).
	Only []string
	// ClassPrefix for safe-fixture silence checks: "BP-", "PERF-", "CWE-".
	// When empty and Only is set, safe checks only the Only rules.
	ClassPrefix string
	// Languages selects built-in language plugins for the scan registry.
	// Empty defaults to Go-only (DefaultRegistry). Python fixture matrices
	// must set Languages to include core.LanguagePython.
	Languages []core.LanguageID
}

// DefaultMatrixOptions returns product-style options for heuristic matrices.
func DefaultMatrixOptions() ScanOptions {
	return ScanOptions{
		IncludeTests: true,
		Profile:      core.ProfileAll,
	}
}

// MaterializeAndScanOpts materializes one .txt fixture and analyzes it.
// Unlike Seed MaterializeAndScan, does not pin --only so sibling rules may fire
// (Rust assert_fixture_rules semantics: require target rule on vulnerable;
// silence only the rule class on safe).
func MaterializeAndScanOpts(relPath string, opts ScanOptions) (Result, error) {
	return materializeAndScanOpts(relPath, opts, scanMaterializedFixture)
}

type matrixFixtureScanner func(ctx *core.ScanContext, srcPath string) (*engine.AnalysisResult, error)

func materializeAndScanOpts(relPath string, opts ScanOptions, scan matrixFixtureScanner) (Result, error) {
	fx, err := FixturesRoot()
	if err != nil {
		return Result{}, err
	}
	txt := filepath.Join(fx, filepath.FromSlash(relPath))
	// Use a per-call temp dir under the system temp so paths stay absolute
	// but isRealProjectScan treats /tmp materializations as fixtures.
	outRoot, err := osMkdirTempFixture()
	if err != nil {
		return Result{}, err
	}
	defer os.RemoveAll(outRoot)
	srcPath, err := fixture.MaterializeFixtureFile(txt, outRoot)
	if err != nil {
		return Result{}, fmt.Errorf("materialize %s: %w", relPath, err)
	}

	profile := opts.Profile
	if profile == 0 {
		profile = core.ProfileAll
	}
	ctx := core.NewScanContext(profile, opts.Only, nil)
	ctx.IncludeTests = opts.IncludeTests
	if opts.TaintEnabled {
		ctx.TaintEnabled = true
	}

	// Prefer language-aware scanner so Python fixtures use the Python plugin.
	var res *engine.AnalysisResult
	if scan == nil || len(opts.Languages) > 0 {
		res, err = scanMaterializedFixtureWithLangs(ctx, srcPath, opts.Languages)
	} else {
		res, err = scan(ctx, srcPath)
	}
	if err != nil {
		return Result{}, fmt.Errorf("analyze %s: %w", srcPath, err)
	}
	findings := res.Findings
	if findings == nil {
		findings = []rules.Finding{}
	}
	ruleID := ruleIDFromFixtureRel(relPath)
	return Result{
		Case: Case{
			RelPath:    relPath,
			RuleID:     ruleID,
			ExpectFire: strings.Contains(relPath, "vulnerable"),
		},
		SourcePath: srcPath,
		Findings:   findings,
		Fired:      HasRule(findings, ruleID),
	}, nil
}

func scanMaterializedFixture(ctx *core.ScanContext, srcPath string) (*engine.AnalysisResult, error) {
	// Back-compat for call sites that do not pass languages through opts.
	return scanMaterializedFixtureWithLangs(ctx, srcPath, nil)
}

func scanMaterializedFixtureWithLangs(ctx *core.ScanContext, srcPath string, langs []core.LanguageID) (*engine.AnalysisResult, error) {
	reg := engine.DefaultRegistry()
	if len(langs) > 0 {
		var err error
		reg, err = engine.NewRegistryWithLanguages(langs...)
		if err != nil {
			return nil, err
		}
	}
	analyzer := engine.NewAnalyzerBuilder().
		Registry(reg).
		ScanContext(ctx).
		Workers(1).
		Build()

	// Analyze the file path when possible so single-file fixtures work.
	res, err := analyzer.AnalyzePaths([]string{srcPath})
	if err == nil {
		return res, nil
	}
	// Fall back to a parent-directory walk.
	return analyzer.AnalyzePaths([]string{filepath.Dir(srcPath)})
}

// AssertVulnerable requires ruleID among findings.
func AssertVulnerable(findings []rules.Finding, ruleID, label string) error {
	if HasRule(findings, ruleID) {
		return nil
	}
	return fmt.Errorf("%s: expected %s, got %s", label, ruleID, SummarizeFindings(findings))
}

// AssertSafeClass requires no finding with the given class prefix (BP-/PERF-/CWE-).
func AssertSafeClass(findings []rules.Finding, classPrefix, label string) error {
	var bad []string
	for _, f := range findings {
		if strings.HasPrefix(f.RuleID, classPrefix) {
			bad = append(bad, f.RuleID)
		}
	}
	if len(bad) == 0 {
		return nil
	}
	return fmt.Errorf("%s: expected no %s findings, got %v", label, classPrefix, bad)
}

func ruleIDFromFixtureRel(rel string) string {
	base := filepath.Base(rel)
	base = strings.TrimSuffix(base, ".txt")
	base = strings.TrimSuffix(base, "-vulnerable")
	base = strings.TrimSuffix(base, "-safe")
	switch {
	case strings.HasPrefix(base, "BP-PY-"):
		// BP-PY-1 / BP-PY-12 → full id (Python catalogue).
		return PythonBPRuleID(base)
	case strings.HasPrefix(base, "BP-"):
		return BPRuleID(base)
	case strings.HasPrefix(base, "PERF-"):
		return PERFRuleID(base)
	case strings.HasPrefix(base, "CWE-"):
		// CWE-89-prepare-same-var → CWE-89
		parts := strings.Split(base, "-")
		if len(parts) >= 2 {
			return parts[0] + "-" + parts[1]
		}
		return base
	default:
		return base
	}
}

// osMkdirTempFixture is a tiny seam for tests (uses os.MkdirTemp).
var osMkdirTempFixture = func() (string, error) {
	return os.MkdirTemp("", "goslop-fixture-*")
}
