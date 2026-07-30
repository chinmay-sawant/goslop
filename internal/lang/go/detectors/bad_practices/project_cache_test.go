package badpractices

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
)

func TestPackageTypeFactsAreScopedAndMemoizedPerScan(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "types.go")
	src := "package sample\n\ntype Reader interface { Read() }\n\ntype impl struct{}\nfunc (impl) Read() {}\n"
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	unit := core.NewParsedUnit(core.LangGo, path, src)

	firstScan := newProjectCaches()
	first := packageTypeFactsForUnit(unit, firstScan)
	if again := packageTypeFactsForUnit(unit, firstScan); first != again {
		t.Fatal("package facts should be reused within one scan")
	}
	secondScan := newProjectCaches()
	if next := packageTypeFactsForUnit(unit, secondScan); next == first {
		t.Fatal("package facts must not cross scan sessions")
	}
}

func TestPackageDocSnapshotBuildsOnceForConcurrentScanWorkers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "doc.go")
	source := "// Package sample documents the package.\npackage sample\n"
	if err := os.WriteFile(path, []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}

	caches := newProjectCaches()
	started := make(chan struct{}, 16)
	release := make(chan struct{})
	var builds atomic.Int32
	caches.packageDocBuilder = func(root string) *PackageDocSnapshot {
		builds.Add(1)
		started <- struct{}{}
		<-release
		return buildPackageDocSnapshot(root)
	}

	results := make(chan *PackageDocSnapshot, 16)
	var workers sync.WaitGroup
	for range 16 {
		workers.Add(1)
		go func() {
			defer workers.Done()
			results <- packageDocSnapshotForDir(dir, caches)
		}()
	}
	<-started
	close(release)
	workers.Wait()
	close(results)

	if got := builds.Load(); got != 1 {
		t.Fatalf("package doc builds=%d want 1", got)
	}
	var first *PackageDocSnapshot
	for snapshot := range results {
		if first == nil {
			first = snapshot
			continue
		}
		if snapshot != first {
			t.Fatal("concurrent workers received distinct package doc snapshots")
		}
	}
	if _, documented := first.DocumentedPackages["sample"]; !documented {
		t.Fatalf("snapshot did not retain package documentation: %#v", first)
	}
}

func TestConcurrentDistinctRootScansDoNotShareProjectFacts(t *testing.T) {
	left := newProjectFactsTestProject(t, "left", "LeftType")
	right := newProjectFactsTestProject(t, "right", "RightType")

	type result struct {
		snapshot *ProjectSnapshot
		types    *packageTypeFacts
		caches   *bpProjectCaches
	}
	results := make(chan result, 2)
	start := make(chan struct{})
	var workers sync.WaitGroup
	for _, project := range []struct {
		unit *core.ParsedUnit
	}{left, right} {
		workers.Add(1)
		go func(unit *core.ParsedUnit) {
			defer workers.Done()
			detector := NewGoBadPracticeScan()
			detector.BeginScan(core.DefaultScanContext())
			facts := buildFacts(unit, detector.caches)
			defer facts.close()
			<-start
			results <- result{
				snapshot: facts.projectSnapshot(unit),
				types:    facts.packageTypeFacts(unit),
				caches:   detector.caches,
			}
			detector.EndScan()
		}(project.unit)
	}
	close(start)
	workers.Wait()
	close(results)

	var scans []result
	for scan := range results {
		scans = append(scans, scan)
	}
	if len(scans) != 2 {
		t.Fatalf("scan results=%d want 2", len(scans))
	}
	if scans[0].caches == scans[1].caches || scans[0].snapshot == scans[1].snapshot || scans[0].types == scans[1].types {
		t.Fatal("distinct scans shared project or package facts")
	}

	wantModules := map[string]string{
		"module example.com/left":  "LeftType",
		"module example.com/right": "RightType",
	}
	for _, scan := range scans {
		var typeName string
		for module, expectedType := range wantModules {
			if strings.Contains(scan.snapshot.GoModText, module) {
				typeName = expectedType
				break
			}
		}
		if typeName == "" {
			t.Fatalf("unexpected module snapshot: %#v", scan.snapshot)
		}
		if _, ok := scan.types.methods[typeName]; !ok {
			t.Fatalf("package facts for %s did not contain %s: %#v", scan.snapshot.Root, typeName, scan.types.methods)
		}
	}
}

func newProjectFactsTestProject(t *testing.T, name, typeName string) struct{ unit *core.ParsedUnit } {
	t.Helper()
	dir := t.TempDir()
	goMod := "module example.com/" + name + "\n\ngo 1.26\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}
	source := "package sample\n\ntype " + typeName + " struct{}\n\nfunc (" + typeName + ") Run() {}\n"
	path := filepath.Join(dir, "types.go")
	if err := os.WriteFile(path, []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	return struct{ unit *core.ParsedUnit }{unit: core.NewParsedUnit(core.LangGo, path, source)}
}
