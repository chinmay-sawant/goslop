# v0.0.1 — Performance optimization checklist (from pprof)

> **Status:** **done for P0–P3** (gates, shared AST, benches, docs) — optional further PERF gates only  
> **Branch:** `perf/p0-p1-snapshot-export`  
> **Evidence:** `/tmp/goslop-pprof/*` · `plans/v0.0.1/perf-p0p3-measurement.md`  
> **Benches (pre → after `benchtime=5s`, gopdfsuit, 2026-07-30):**  
> - `BenchmarkScanProfileAll` ~**164 ms** → **112.7 ms**/op · 199 MB → 62 MB · 1.19M → 444k allocs  
> - `BenchmarkScanAndExport` ~**367 ms** → **186.8 ms**/op · 345 MB → 146 MB · 2.64M → 1.09M allocs  
> - `BenchmarkExportOnly` ~**201 ms** → **71.7 ms**/op · 149 MB → 85 MB · 1.47M → 647k allocs  
> - Product **`make run` scan**: pre ~**190 ms** → post ~**112–117 ms** mean · best **99 ms** (915 findings, 915+37 exports)  
> **Corpus:** gopdfsuit · `internal/bench`

Legend: `[ ]` not started · `[~]` partial · `[x]` done with evidence

---

## Gates (prove after each batch)

- [x] Duration benches `benchtime=5s` — `/tmp/goslop-pprof/bench-after-5s.txt` (Scan −31%, ScanAndExport −49%, Export −64%)
- [x] Fixed-iter 20x confirmation earlier same day — `/tmp/goslop-pprof/bench-after-p0p3-20x.txt`
- [x] `go tool pprof -top -cum` on scan + export CPU profiles — see `perf-p0p3-measurement.md`
- [x] §12.4 **915** findings / **915+37** exports; `make run` scan **~99–117 ms** (was ~190 ms)
- [x] `go test ./internal/export/ ./internal/lang/go/detectors/... ./tests/integration/ -count=1`

**Measured:** product scan ~**1.6–1.9×** faster; benches confirm; parity OK.

---

## P0 — highest leverage

### P0.1 BP project snapshot memoization (~31% scan CPU)

Package: `internal/lang/go/detectors/bad_practices` (`project.go`)

- [x] Cache `ProjectSnapshot` **once per scan root** for the whole `AnalyzePaths` (not per file / per BP rule)
- [x] Memoize `projectSnapshotForRoot` (`sync.Once` + map keyed by absolute root)
- [x] Ensure concurrent workers share one snapshot safely (no races / no thundering-herd rebuild)
- [x] Invalidate only across scans (`clear()` / new session caches)
- [x] Re-profile: `buildProjectSnapshot` / `WalkDir` cum % drops substantially — **~31% → ~1.5%** scan cum (`perf-p0p3-measurement.md`)
- [x] Document expected win in PR (tens of % scan CPU on multi-file trees)

### P0.2 Export function span index (~18% export CPU)

Package: `internal/export` (`export.go` — `enclosingFunctionLines` / `functionWindow`)

- [x] Build `FuncDecl` / `FuncLit` **line spans once per file** (at parse, stored on `parsedFile`)
- [x] Look up hit line via span table (no full `go/ast.Inspect` per finding)
- [x] Prefer outermost `FuncDecl` over inner `FuncLit` (current product behavior)
- [x] Fall back to line window when no enclosing function
- [x] Context + chunks still share the same span builder / `astCache`
- [x] Re-profile: `enclosingFunctionLines` / `go/ast.Inspect` cum % drops — symbol **gone** (~18% → 0%; residual Inspect ~2.6% once/file)
- [x] Unit tests: multi-finding same file; defer-closure prefers outer decl (`export_test.go`)

---

## P1 — high value

### P1.1 `codeLines` once per file (~3%+ CPU, ~16% alloc_space)

Package: `internal/lang/go/detectors/bad_practices` (`common.go` + callers)

- [x] Compute `codeLines` **once per file** in `buildFacts`
- [x] Pass shared line view into BP detectors (`codeLinesFacts`)
- [x] Named `codeLine` type + pre-sized slice
- [x] Line walk without repeated `strings.Split` growslice
- [x] Re-profile: `codeLines` / `stripLineComment` alloc_space flat % drops — CPU cum **~3.3% → ~0.25%** (`buildCodeLines` once/file)

### P1.2 `numberedLines` without per-line `fmt.Sprintf` (~15% export CPU, ~40% alloc_objects)

Package: `internal/export` (`numberedLines`)

- [x] Precompute line starts once per file (`computeLineStarts` on `parsedFile`)
- [x] Avoid re-`strings.Split` of full content per finding
- [x] Format lines with `strings.Builder` + `strconv` (not `fmt.Sprintf` per line)
- [x] Keep `>` marker and line numbers product-compatible
- [x] Re-profile: `numberedLines` / `fmt.Sprintf` alloc_objects drop — export `fmt.Sprintf` cum **~8% → ~0.7%**; numbered path via cache
- [x] Export unit tests still pass (`export_test.go`)

### P1.3 Format finding block once when dual-exporting

Package: `internal/export`

- [x] When both `--export-context` and `--export-chunks` are on, format each finding **once** and reuse for chunk body
- [x] Confirm no behavior change in chunk separators / headers (unit tests)

---

## P2 — medium

### P2.1 Snapshot walk filters

Package: `internal/lang/go/detectors/bad_practices` (`buildProjectSnapshot`)

- [x] Skip `vendor/`, `.git/`, `node_modules/`, and other non-source dirs in snapshot walk (`skipProjectDirs` + all dot-dirs + bin/dist/scripts/…)
- [x] Prefer path heuristics before full-file `ReadFile` + `Contains` (skip `_test.go` / examples; stop content scan when flags complete)
- [x] Root `go.mod`/`go.sum` only; non-test `.go` names for anchor; content uses preallocated `bytes.Contains` needles

### P2.2 Reduce `goparse.(*Tree).Slice` copies (~5% allocs)

Package: `internal/lang/go/goparse` + PERF hot paths

- [x] Audit hot callers of `Tree.Slice` / `NodeText` (facts, taint, BP)
- [x] Zero-copy `Slice` via `unsafe.String` over immutable `Tree.Source`
- [x] Re-profile: `Tree.Slice` alloc_space flat % drops — covered in post P0–P3.1 alloc cut (scan B/op −69%); Slice no longer a standout

### P2.3 Package-level / non-function export early-out

Package: `internal/export`

- [x] Skip whole-function parse for package-level BP rules (BP-41, BP-57–65) and files without `func`
- [x] Keep `<context unavailable>` / line-window fallbacks

---

## P3 — later / catalogue scale

### P3.1 PERF / CWE rule scheduling

Packages: `internal/lang/go/detectors/perf`, `.../cwe`, `.../bad_practices`

- [x] Immutable catalogue snapshot once after `init` (no per-file register lock + copy)
- [x] Cached sorted `RuleIDs()` (no alloc/sort per file / per `anyRuleAllowed`)
- [x] Keep `Allows` per file (safe for shared detector instances + parallel tests)
- [x] Skip rule domains when needles / facts are absent (callee / token index) — residual
  - CWE: table needle rules gated via `RegisterRule(..., gates...)` + `Index.HasAny` (group needles flattened)
  - PERF: optional `gates` on `RegisterRule` / `ruleEntry` (infrastructure; individual rules opt-in)
- [x] Reuse parse/`unit.Tree` across PERF / CWE-taint / BP (`goparse.TreeForUnit`; pack-local SourceIndex + AST inspect still separate)
- [ ] Optional: domain-level disable in profiles without full catalogue cost

### P3.2 GC pressure (follow-on)

- [x] After P0–P1 alloc cuts, re-check `runtime.scanObjectsSmall` / `gcDrain` share on export profile — **`gcDrain` still ~24% export cum** (was ~23%); absolute export much faster; `scanObjectsSmall` no longer a top named sample
- [ ] Only then consider GOGC tuning notes for CLI (product default unchanged unless measured) — **optional**; still GC-relative-hot on export after app-path wins

---

## Measurement checklist (each PR)

- [x] Before numbers: prior medians above (or `make bench BENCHTIME=20x`)
- [x] After numbers: same command, same `GOSLOP_BENCH_SCAN_PATH` — `bench-after-p0p3-20x.txt`
- [x] CPU: `go tool pprof -top -cum <cpu.prof>` for changed path — `scan-cpu-after.prof` / `export-cpu-after.prof`
- [ ] Mem (if alloc work): `go tool pprof -top -sample_index=alloc_space <mem.prof>` — optional; bench B/op/allocs already show large cuts
- [x] Correctness: gopdfsuit findings count / export file counts unchanged (**`make run` ~190 ms baseline** → ~112–122 ms)
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
| 1 | P0.1 + P0.2 + P1.1 + P1.2 + P1.3 | Snapshot Once + export span/format + codeLines | **committed** (`9105596`, ~122ms scan) |
| 2 | P2.* | Walk filters + Slice + early-out | **committed** |
| 3 | P3.1 catalogue snapshot + RuleIDs cache | Dispatch overhead | **code done** |
| 4 | P3 residual | Needle gates / cross-pack facts / GC | open |

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
