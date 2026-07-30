# v0.0.1 — Performance optimization checklist (from pprof)

> **Status:** in progress — **P0 + P1 implemented**; measurement gates pending  
> **Branch:** `perf/p0-p1-snapshot-export`  
> **Evidence:** `/tmp/goslop-pprof/*` (pre-change) + unit/integration tests  
> **Benches (pre-change median, gopdfsuit):**  
> - `BenchmarkScanProfileAll` ~**164 ms**/op · 199 MB · 1.19M allocs  
> - `BenchmarkScanAndExport` ~**367 ms**/op · 345 MB · 2.64M allocs  
> - `BenchmarkExportOnly` ~**201 ms**/op · 149 MB · 1.47M allocs  
> - Product **`make run` wall**: pre ~**190 ms** → post ~**122 ms** scan (915 findings, 915+37 exports unchanged; 2026-07-30)  
> **Corpus:** gopdfsuit · `internal/bench`

Legend: `[ ]` not started · `[~]` partial · `[x]` done with evidence

---

## Gates (prove after each batch)

- [ ] `make bench BENCHTIME=20x` (or `100x`) vs pre-change median
- [ ] `go tool pprof -top -cum` on scan + export CPU profiles (spot-check hotspots drop)
- [x] §12.4 still **915** findings / **915+37** exports on gopdfsuit (`make run` 121.7ms scan)
- [x] `go test ./internal/export/ ./internal/lang/go/detectors/... ./tests/integration/ -count=1` (2026-07-30)

**Measured:** `make run` scan **121.7ms** (was ~190ms); findings/exports parity OK.

---

## P0 — highest leverage

### P0.1 BP project snapshot memoization (~31% scan CPU)

Package: `internal/lang/go/detectors/bad_practices` (`project.go`)

- [x] Cache `ProjectSnapshot` **once per scan root** for the whole `AnalyzePaths` (not per file / per BP rule)
- [x] Memoize `projectSnapshotForRoot` (`sync.Once` + map keyed by absolute root)
- [x] Ensure concurrent workers share one snapshot safely (no races / no thundering-herd rebuild)
- [x] Invalidate only across scans (`clear()` / new session caches)
- [ ] Re-profile: `buildProjectSnapshot` / `WalkDir` cum % drops substantially
- [ ] Document expected win in PR (tens of % scan CPU on multi-file trees)

### P0.2 Export function span index (~18% export CPU)

Package: `internal/export` (`export.go` — `enclosingFunctionLines` / `functionWindow`)

- [x] Build `FuncDecl` / `FuncLit` **line spans once per file** (at parse, stored on `parsedFile`)
- [x] Look up hit line via span table (no full `go/ast.Inspect` per finding)
- [x] Prefer outermost `FuncDecl` over inner `FuncLit` (current product behavior)
- [x] Fall back to line window when no enclosing function
- [x] Context + chunks still share the same span builder / `astCache`
- [ ] Re-profile: `enclosingFunctionLines` / `go/ast.Inspect` cum % drops
- [x] Unit tests: multi-finding same file; defer-closure prefers outer decl (`export_test.go`)

---

## P1 — high value

### P1.1 `codeLines` once per file (~3%+ CPU, ~16% alloc_space)

Package: `internal/lang/go/detectors/bad_practices` (`common.go` + callers)

- [x] Compute `codeLines` **once per file** in `buildFacts`
- [x] Pass shared line view into BP detectors (`codeLinesFacts`)
- [x] Named `codeLine` type + pre-sized slice
- [x] Line walk without repeated `strings.Split` growslice
- [ ] Re-profile: `codeLines` / `stripLineComment` alloc_space flat % drops

### P1.2 `numberedLines` without per-line `fmt.Sprintf` (~15% export CPU, ~40% alloc_objects)

Package: `internal/export` (`numberedLines`)

- [x] Precompute line starts once per file (`computeLineStarts` on `parsedFile`)
- [x] Avoid re-`strings.Split` of full content per finding
- [x] Format lines with `strings.Builder` + `strconv` (not `fmt.Sprintf` per line)
- [x] Keep `>` marker and line numbers product-compatible
- [ ] Re-profile: `numberedLines` / `fmt.Sprintf` alloc_objects drop
- [x] Export unit tests still pass (`export_test.go`)

### P1.3 Format finding block once when dual-exporting

Package: `internal/export`

- [x] When both `--export-context` and `--export-chunks` are on, format each finding **once** and reuse for chunk body
- [x] Confirm no behavior change in chunk separators / headers (unit tests)

---

## P2 — medium

### P2.1 Snapshot walk filters

Package: `internal/lang/go/detectors/bad_practices` (`buildProjectSnapshot`)

- [x] Skip `vendor/`, `.git/`, `node_modules/`, and other non-source dirs in snapshot walk (already in `skipProjectDirs`)
- [ ] Prefer path heuristics before full-file `ReadFile` + `Contains`
- [ ] Only open likely anchors (`go.mod`, `*.go`, known entrypoints) when possible

### P2.2 Reduce `goparse.(*Tree).Slice` copies (~5% allocs)

Package: `internal/lang/go/goparse` + PERF hot paths

- [ ] Audit hot callers of `Tree.Slice` (`recordCallAST`, assign paths, etc.)
- [ ] Prefer shared source + offset range over new string when safe
- [ ] Re-profile: `Tree.Slice` alloc_space flat % drops

### P2.3 Package-level / non-function export early-out

Package: `internal/export`

- [ ] Skip whole-function parse for package-level / non-Go findings when span lookup is impossible
- [ ] Keep `<context unavailable>` / line-window fallbacks

---

## P3 — later / catalogue scale

### P3.1 PERF / CWE rule scheduling

Packages: `internal/lang/go/detectors/perf`, `.../cwe`

- [ ] Skip rule domains when needles / facts are absent (callee / token index)
- [ ] Reuse facts across PERF / CWE / BP where the same AST walk would repeat
- [ ] Optional: domain-level disable in profiles without full catalogue cost

### P3.2 GC pressure (follow-on)

- [ ] After P0–P1 alloc cuts, re-check `runtime.scanObjectsSmall` / `gcDrain` share on export profile
- [ ] Only then consider GOGC tuning notes for CLI (product default unchanged unless measured)

---

## Measurement checklist (each PR)

- [ ] Before numbers: prior medians above (or `make bench BENCHTIME=20x`)
- [ ] After numbers: same command, same `GOSLOP_BENCH_SCAN_PATH`
- [ ] CPU: `go tool pprof -top -cum <cpu.prof>` for changed path
- [ ] Mem (if alloc work): `go tool pprof -top -sample_index=alloc_space <mem.prof>`
- [ ] Correctness: gopdfsuit findings count / export file counts unchanged (**`make run` ~190 ms baseline**)
- [ ] Link PR to this checklist items closed (`[x]`)

### Reproduce profiles

```sh
export CGO_ENABLED=0 GOTOOLCHAIN=local
export GOSLOP_BENCH_SCAN_PATH=/home/chinmay/ChinmayPersonalProjects/gopdfsuit

go test -run='^$' -bench=BenchmarkScanProfileAll -benchtime=5s -benchmem \
  -cpuprofile=/tmp/goslop-pprof/scan-cpu.prof \
  -memprofile=/tmp/goslop-pprof/scan-mem.prof \
  ./internal/bench/

go tool pprof -top -cum /tmp/goslop-pprof/scan-cpu.prof
```

---

## Suggested PR order

| PR | Checklist items | Theme | Status |
|----|-----------------|--------|--------|
| 1 | P0.1 + P0.2 + P1.1 + P1.2 + P1.3 | Snapshot Once + export span/format + codeLines | **code done** |
| 2 | P2.* | Walk filters + Slice + early-out | open |
| 3 | P3.* | Rule scheduling / residual GC | open |

---

## Out of scope (this checklist)

- [ ] Changing detector *semantics* / finding counts for speed
- [ ] Dropping whole-function export default
- [ ] Committing raw `.prof` blobs to git
- [ ] Micro-benchmarking `strings.Index` in isolation without fixing callers

---

## References

- pprof data: `/tmp/goslop-pprof/`
- Bench harness: `internal/bench/make_run_bench_test.go`
- Snapshot code: `internal/lang/go/detectors/bad_practices/project.go`
- Export code: `internal/export/export.go`
