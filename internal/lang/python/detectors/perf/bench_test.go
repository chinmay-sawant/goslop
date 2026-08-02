package perf

import (
	"strings"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Microbenchmark for facts build + full PERF-PY catalogue on a synthetic service file.
func BenchmarkPythonPerfScanSynthetic(b *testing.B) {
	src := strings.Repeat(`def worker(batch):
    for job in batch:
        sku = Sku.objects.get(code=job["code"])
        deliver_webhook(job)
`, 40)
	unit := core.NewParsedUnit(core.LanguagePython, "bench.py", src)
	ctx := core.DefaultScanContext()
	ctx.Only = nil
	detector := NewPythonPerfScan()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var out []rules.Finding
		detector.Run(ctx, unit, &out)
	}
}
