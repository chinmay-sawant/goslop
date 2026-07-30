// Package bench hosts product-level Go benchmarks for goslop.
//
// These are standard library tests (testing.B). Run with:
//
//	go test -bench=. -benchmem -benchtime=100x ./internal/bench/
//
// Reports ns/op, B/op, and allocs/op (with -benchmem). No external harness.
//
// Compare two toolchains (write go.mod go directive to match each binary first):
//
//	go1.25.0 test -bench=. -benchmem -benchtime=100x ./internal/bench/ | tee /tmp/goslop-bench-go1.25.txt
//	go1.26.4 test -bench=. -benchmem -benchtime=100x ./internal/bench/ | tee /tmp/goslop-bench-go1.26.4.txt
//	# optional: go install golang.org/x/perf/cmd/benchstat@latest && benchstat a.txt b.txt
package bench_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/engine"
	"github.com/chinmay/goslop/internal/export"
)

// defaultScanPath matches Makefile SCAN_PATH (gopdfsuit §12.4 reference corpus).
const defaultScanPath = "/home/chinmay/ChinmayPersonalProjects/gopdfsuit"

func scanPath(b *testing.B) string {
	b.Helper()
	p := os.Getenv("GOSLOP_BENCH_SCAN_PATH")
	if p == "" {
		p = defaultScanPath
	}
	if st, err := os.Stat(p); err != nil || !st.IsDir() {
		b.Skipf("scan path unavailable %q (set GOSLOP_BENCH_SCAN_PATH): %v", p, err)
	}
	return p
}

func newAnalyzer(profile core.ScanProfile, retainSources bool) *engine.Analyzer {
	ctx := core.NewScanContext(profile, nil, nil)
	ctx.NoCache = true
	ctx.NoBaseline = true
	ctx.RetainSources = retainSources
	walk := engine.DefaultWalkOptions()
	return engine.NewAnalyzerBuilder().
		Registry(engine.DefaultRegistry()).
		ScanContext(ctx).
		WalkOptions(walk).
		ProjectRoot(".").
		Build()
}

// BenchmarkScanProfileAll is the engine half of `make run` (profile all, no cache).
// Best signal for runtime/GC changes: no disk export noise.
func BenchmarkScanProfileAll(b *testing.B) {
	path := scanPath(b)
	a := newAnalyzer(core.ProfileAll, false)

	// Sanity: one warm run outside the timed loop.
	res, err := a.AnalyzePaths([]string{path})
	if err != nil {
		b.Fatalf("warmup AnalyzePaths: %v", err)
	}
	if n := len(res.Findings); n == 0 {
		b.Fatalf("warmup: expected findings, got 0")
	}
	wantFindings := len(res.Findings)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := a.AnalyzePaths([]string{path})
		if err != nil {
			b.Fatalf("AnalyzePaths: %v", err)
		}
		if len(res.Findings) != wantFindings {
			b.Fatalf("findings=%d want %d", len(res.Findings), wantFindings)
		}
		runtime.KeepAlive(res)
	}
}

// BenchmarkScanAndExport matches product `make run` surface more closely:
// profile all + export context/chunks (whole-function Context) + no cache.
// Disk I/O dominates; still useful for end-to-end wall regressions.
func BenchmarkScanAndExport(b *testing.B) {
	path := scanPath(b)
	a := newAnalyzer(core.ProfileAll, true)
	outRoot := b.TempDir()

	// Warmup
	res, err := a.AnalyzePaths([]string{path})
	if err != nil {
		b.Fatalf("warmup: %v", err)
	}
	if len(res.Findings) == 0 {
		b.Fatal("warmup: no findings")
	}
	wantFindings := len(res.Findings)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ctxDir := filepath.Join(outRoot, "ctx")
		chunkDir := filepath.Join(outRoot, "chunks")
		_ = os.RemoveAll(ctxDir)
		_ = os.RemoveAll(chunkDir)
		b.StartTimer()

		res, err := a.AnalyzePaths([]string{path})
		if err != nil {
			b.Fatalf("AnalyzePaths: %v", err)
		}
		if len(res.Findings) != wantFindings {
			b.Fatalf("findings=%d want %d", len(res.Findings), wantFindings)
		}
		_, err = export.ExportFindings(res.Findings, export.Options{
			ExportContext:    true,
			ExportChunks:     true,
			ContextOutputDir: ctxDir,
			ChunksOutputDir:  chunkDir,
			// WholeFunction nil → default true
		}, res.SourceCache)
		if err != nil {
			b.Fatalf("ExportFindings: %v", err)
		}
		runtime.KeepAlive(res)
	}
}

// BenchmarkExportOnly times on-disk export alone (fixed finding set from one scan).
// Isolates whole-function context expansion cost from the analyzer.
func BenchmarkExportOnly(b *testing.B) {
	path := scanPath(b)
	a := newAnalyzer(core.ProfileAll, true)
	res, err := a.AnalyzePaths([]string{path})
	if err != nil {
		b.Fatalf("scan: %v", err)
	}
	if len(res.Findings) == 0 {
		b.Fatal("no findings")
	}
	outRoot := b.TempDir()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ctxDir := filepath.Join(outRoot, "ctx")
		chunkDir := filepath.Join(outRoot, "chunks")
		_ = os.RemoveAll(ctxDir)
		_ = os.RemoveAll(chunkDir)
		b.StartTimer()

		_, err := export.ExportFindings(res.Findings, export.Options{
			ExportContext:    true,
			ExportChunks:     true,
			ContextOutputDir: ctxDir,
			ChunksOutputDir:  chunkDir,
		}, res.SourceCache)
		if err != nil {
			b.Fatalf("ExportFindings: %v", err)
		}
	}
}
