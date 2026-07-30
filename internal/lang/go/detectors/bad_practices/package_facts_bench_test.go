package badpractices

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
)

func BenchmarkPackageTypeFactsMultiFile(b *testing.B) {
	units := benchmarkPackageFactUnits(b)

	b.Run("cold_scan", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			caches := newProjectCaches()
			for _, unit := range units {
				facts := packageTypeFactsForUnit(unit, caches)
				if len(facts.methods) == 0 {
					b.Fatal("expected package methods")
				}
				runtime.KeepAlive(facts)
			}
		}
	})

	b.Run("cached_scan", func(b *testing.B) {
		caches := newProjectCaches()
		for _, unit := range units {
			_ = packageTypeFactsForUnit(unit, caches)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, unit := range units {
				facts := packageTypeFactsForUnit(unit, caches)
				if len(facts.methods) == 0 {
					b.Fatal("expected cached package methods")
				}
				runtime.KeepAlive(facts)
			}
		}
	})
}

func benchmarkPackageFactUnits(b *testing.B) []*core.ParsedUnit {
	b.Helper()
	dir := b.TempDir()
	const files = 12
	units := make([]*core.ParsedUnit, 0, files)
	for i := 0; i < files; i++ {
		typeName := fmt.Sprintf("Type%d", i)
		source := fmt.Sprintf("package benchmarkpkg\n\ntype %s struct{}\n\nfunc (%s) Run%d() {}\n", typeName, typeName, i)
		path := filepath.Join(dir, fmt.Sprintf("type_%d.go", i))
		if err := os.WriteFile(path, []byte(source), 0o644); err != nil {
			b.Fatal(err)
		}
		units = append(units, core.NewParsedUnit(core.LangGo, path, source))
	}
	return units
}
