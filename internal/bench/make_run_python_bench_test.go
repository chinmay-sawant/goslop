// Python product-level benches for `make run-python` (languages=["python"]).
//
//	GOSLOP_BENCH_PYTHON_SCAN_PATH=... go test -run='^$' -bench=BenchmarkPython \
//	  -benchmem -benchtime=5s -cpuprofile=/tmp/goslop-python-pprof/scan-cpu.prof \
//	  -memprofile=/tmp/goslop-python-pprof/scan-mem.prof ./internal/bench/
package bench_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/engine"
	"github.com/chinmay-sawant/goslop/internal/export"
)

const defaultPythonScanPath = "/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/pythoncoreengine"

func pythonScanPath(b *testing.B) string {
	b.Helper()
	p := os.Getenv("GOSLOP_BENCH_PYTHON_SCAN_PATH")
	if p == "" {
		p = defaultPythonScanPath
	}
	if st, err := os.Stat(p); err != nil || !st.IsDir() {
		b.Skipf("python scan path unavailable %q (set GOSLOP_BENCH_PYTHON_SCAN_PATH): %v", p, err)
	}
	return p
}

func newPythonAnalyzer(retainSources bool) *engine.Analyzer {
	reg, err := engine.NewRegistryWithLanguages(core.LanguagePython)
	if err != nil {
		panic(err)
	}
	ctx := core.NewScanContext(core.ProfileAll, nil, nil)
	ctx.NoCache = true
	ctx.NoBaseline = true
	ctx.RetainSources = retainSources
	walk := engine.DefaultWalkOptions()
	return engine.NewAnalyzerBuilder().
		Registry(reg).
		ScanContext(ctx).
		WalkOptions(walk).
		ProjectRoot(".").
		Build()
}

// BenchmarkPythonScanProfileAll is the engine half of `make run-python`.
func BenchmarkPythonScanProfileAll(b *testing.B) {
	path := pythonScanPath(b)
	a := newPythonAnalyzer(false)

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

// BenchmarkPythonScanAndExport matches product `make run-python` export surface.
func BenchmarkPythonScanAndExport(b *testing.B) {
	path := pythonScanPath(b)
	a := newPythonAnalyzer(true)
	outRoot := b.TempDir()

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
		}, res.SourceCache)
		if err != nil {
			b.Fatalf("ExportFindings: %v", err)
		}
		runtime.KeepAlive(res)
	}
}
