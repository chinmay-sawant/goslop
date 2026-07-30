package integration

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/engine"
)

// TestDogfoodDetectorSourcesNoHang scans the Go detector implementation tree.
//
// Why this is not covered by fixture matrices:
//   - BP/CWE/PERF fixtures are small, balanced Go snippets (real interfaces,
//     real select blocks). They never contain the self-referential string
//     needles that live only in detector source, e.g. strings.Index(..., "interface {").
//   - The BP-28/BP-29 hang required scanning rules_api.go against itself: the
//     text scanner matched its own needle inside a string literal, brace
//     matching never returned to depth 0, and start=end looped forever.
//
// A dogfood scan of internal/lang/go/detectors is the cheapest way to keep
// "analyzer hung on its own source" from shipping again.
func TestDogfoodDetectorSourcesNoHang(t *testing.T) {
	root, err := RepoRoot()
	if err != nil {
		t.Fatal(err)
	}
	detectors := filepath.Join(root, "internal", "lang", "go", "detectors")

	// Profile all: BP text scanners + PERF + CWE. Must finish well under this wall.
	const wall = 10 * time.Second

	ctx := core.NewScanContext(core.ProfileAll, nil, nil)
	ctx.NoCache = true
	ctx.NoBaseline = true
	a := engine.NewAnalyzerBuilder().
		Registry(engine.DefaultRegistry()).
		ScanContext(ctx).
		WalkOptions(engine.DefaultWalkOptions()).
		ProjectRoot(root).
		Build()

	done := make(chan error, 1)
	var files, findings int
	go func() {
		res, err := a.AnalyzePaths([]string{detectors})
		if err != nil {
			done <- err
			return
		}
		if res != nil && res.Stats != nil {
			files = res.Stats.FilesScanned
		}
		if res != nil {
			findings = len(res.Findings)
		}
		done <- nil
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("dogfood AnalyzePaths: %v", err)
		}
		if files == 0 {
			t.Fatal("dogfood scanned 0 files; path wrong?")
		}
		t.Logf("dogfood detectors: files=%d findings=%d (wall limit %s)", files, findings, wall)
	case <-time.After(wall):
		t.Fatalf("dogfood scan of %s did not finish within %s (likely infinite loop in a text scanner)", detectors, wall)
	}
}
