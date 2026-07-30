package engine_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/engine"
	"github.com/chinmay-sawant/goslop/internal/engine/cache"
	golang "github.com/chinmay-sawant/goslop/internal/lang/go"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestAnalyzerCacheHitAccumulatesStateWithoutRerunningDetection(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "one.go", "package sample\nfunc One() {}\n")
	writeFile(t, dir, "two.go", "package sample\nfunc Two() {}\n")

	det := &statefulCacheDetector{}
	reg, err := engine.NewRegistry([]core.LanguagePlugin{&goTestPlugin{dets: []core.Detector{det}}})
	if err != nil {
		t.Fatal(err)
	}
	ctx := core.DefaultScanContext()
	ctx.Only = []string{"TEST-CACHE-STATE"}
	a := engine.NewAnalyzerBuilder().
		Registry(reg).
		ScanContext(ctx).
		Cache(cache.InMemory("test-stateful-cache")).
		ProjectRoot(dir).
		Workers(1).
		Build()

	cold, err := a.AnalyzePaths([]string{dir})
	if err != nil {
		t.Fatal(err)
	}
	if det.runCount != 2 || det.finalizedUnits != 2 {
		t.Fatalf("cold scan counts run=%d finalized=%d, want 2/2", det.runCount, det.finalizedUnits)
	}

	warm, err := a.AnalyzePaths([]string{dir})
	if err != nil {
		t.Fatal(err)
	}
	if warm.Stats.CacheHits != 2 || warm.Stats.CacheMisses != 0 {
		t.Fatalf("warm cache stats=%+v, want two hits and no misses", warm.Stats)
	}
	if det.runCount != 0 {
		t.Fatalf("ordinary detection ran %d times on cache hits", det.runCount)
	}
	if det.finalizedUnits != 2 {
		t.Fatalf("cache hits accumulated %d state units, want 2", det.finalizedUnits)
	}
	assertFindingFingerprintParity(t, cold.Findings, warm.Findings)
}

func TestAnalyzerCacheTaintParityAcrossColdWarmAndMixedScans(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "handler.go", `package app

import "net/http"

func handler(r *http.Request) {
	input := r.URL.Query().Get("input")
	openPath(input)
}
`)
	writeFile(t, dir, "helper.go", `package app

import "os"

func openPath(input string) {
	os.Open(input)
}
`)

	reg, err := engine.NewRegistry([]core.LanguagePlugin{golang.NewPlugin()})
	if err != nil {
		t.Fatal(err)
	}
	ctx := core.DefaultScanContext()
	ctx.TaintEnabled = true
	ctx.TaintMaxDepth = 4
	ctx.Only = []string{"CWE-22"}
	a := engine.NewAnalyzerBuilder().
		Registry(reg).
		ScanContext(ctx).
		Cache(cache.InMemory("test-taint-cache-parity")).
		ProjectRoot(dir).
		Build()

	cold, err := a.AnalyzePaths([]string{dir})
	if err != nil {
		t.Fatal(err)
	}
	if !hasFindingForRule(cold.Findings, "CWE-22") {
		t.Fatalf("cold scan did not report the inter-file taint flow: %#v", cold.Findings)
	}

	warm, err := a.AnalyzePaths([]string{dir})
	if err != nil {
		t.Fatal(err)
	}
	if warm.Stats.CacheHits != 2 || warm.Stats.CacheMisses != 0 {
		t.Fatalf("warm cache stats=%+v, want two hits and no misses", warm.Stats)
	}
	assertFindingFingerprintParity(t, cold.Findings, warm.Findings)

	// Only handler.go becomes stale; helper.go must still contribute its
	// cached state to the inter-file finalization result.
	writeFile(t, dir, "handler.go", `package app

import "net/http"

func handler(r *http.Request) {
	// Deliberately edit only this file to produce a mixed cache scan.
	input := r.URL.Query().Get("input")
	openPath(input)
}
`)
	mixed, err := a.AnalyzePaths([]string{dir})
	if err != nil {
		t.Fatal(err)
	}
	if mixed.Stats.CacheHits != 1 || mixed.Stats.CacheMisses != 1 {
		t.Fatalf("mixed cache stats=%+v, want one hit and one miss", mixed.Stats)
	}
	assertFindingFingerprintParity(t, cold.Findings, mixed.Findings)
}

var statefulCacheMeta = &rules.RuleMetadata{
	ID:       "TEST-CACHE-STATE",
	Title:    "stateful cache detector",
	Severity: rules.SeverityInfo,
}

// statefulCacheDetector models a detector with Finalize-time state. It makes
// the cache contract explicit: a hit may reuse ordinary findings but must not
// skip state accumulation.
type statefulCacheDetector struct {
	core.BaseDetector
	runCount       int
	stateUnits     int
	finalizedUnits int
}

func (d *statefulCacheDetector) Language() core.LanguageID { return core.LangGo }

func (d *statefulCacheDetector) RuleIDs() []string { return []string{"TEST-CACHE-STATE"} }

func (d *statefulCacheDetector) MetadataFor(ruleID string) *rules.RuleMetadata {
	if ruleID == statefulCacheMeta.ID {
		return statefulCacheMeta
	}
	return nil
}

func (d *statefulCacheDetector) BeginScan(*core.ScanContext) {
	d.runCount = 0
	d.stateUnits = 0
	d.finalizedUnits = 0
}

func (d *statefulCacheDetector) RequiresCacheState(*core.ScanContext) bool { return true }

func (d *statefulCacheDetector) Run(_ *core.ScanContext, unit *core.ParsedUnit, out *[]rules.Finding) {
	d.runCount++
	rules.PushFinding(statefulCacheMeta, unit.DisplayPath, 1, 1, "ordinary detection", out)
}

func (d *statefulCacheDetector) AccumulateState(*core.ScanContext, *core.ParsedUnit) {
	d.stateUnits++
}

func (d *statefulCacheDetector) Finalize(_ *core.ScanContext, out *[]rules.Finding) {
	d.finalizedUnits = d.stateUnits
	rules.PushFinding(statefulCacheMeta, "project", 1, 1, "finalized state", out)
}

func assertFindingFingerprintParity(t *testing.T, want, got []rules.Finding) {
	t.Helper()
	wantFingerprints := findingFingerprints(want)
	gotFingerprints := findingFingerprints(got)
	if strings.Join(wantFingerprints, "\n") != strings.Join(gotFingerprints, "\n") {
		t.Fatalf("finding fingerprints differ\nwant: %v\n got: %v", wantFingerprints, gotFingerprints)
	}
}

func findingFingerprints(findings []rules.Finding) []string {
	fingerprints := make([]string, len(findings))
	for i, finding := range findings {
		fingerprints[i] = finding.FingerprintString()
	}
	sort.Strings(fingerprints)
	return fingerprints
}

func hasFindingForRule(findings []rules.Finding, ruleID string) bool {
	for _, finding := range findings {
		if finding.RuleID == ruleID {
			return true
		}
	}
	return false
}
