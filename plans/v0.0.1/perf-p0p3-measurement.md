# P0–P3.1 measurement (post-change re-profile)

> **Date:** 2026-07-30  
> **Branch:** `perf/p0-p1-snapshot-export` (P0 + P1 + P2 + P3.1 code landed)  
> **Corpus:** gopdfsuit (`GOSLOP_BENCH_SCAN_PATH=/home/chinmay/ChinmayPersonalProjects/gopdfsuit`)  
> **Toolchain:** Go 1.25.0 · `CGO_ENABLED=0`  
> **Harness:** `./internal/bench/` · formal `benchtime=20x` + CPU `benchtime=5s`  
> **Artifacts:** `/tmp/goslop-pprof/bench-after-p0p3-20x.txt`, `scan-cpu-after.prof`, `export-cpu-after.prof`, `*-top*.txt`

---

## Bench table (before → after)

Pre-change medians from `/tmp/goslop-pprof/Benchmark*-bench.txt` (short runs, same machine/corpus).  
After: `go test -run='^$' -bench=. -benchmem -benchtime=20x ./internal/bench/`.

| Benchmark | ns/op (ms) | B/op | allocs/op | Δ time | Δ B/op | Δ allocs |
|-----------|------------|------|-----------|--------|--------|----------|
| **ScanProfileAll** before | 164 177 113 (**164.2 ms**) | 199 386 155 | 1 190 864 | | | |
| **ScanProfileAll** after | 113 655 862 (**113.7 ms**) | 61 296 874 | 434 657 | **−30.8%** | **−69.3%** | **−63.5%** |
| **ScanAndExport** before | 366 688 659 (**366.7 ms**) | 345 187 011 | 2 644 075 | | | |
| **ScanAndExport** after | 181 693 366 (**181.7 ms**) | 145 374 862 | 1 080 914 | **−50.5%** | **−57.9%** | **−59.1%** |
| **ExportOnly** before | 200 627 296 (**200.6 ms**) | 148 525 689 | 1 467 566 | | | |
| **ExportOnly** after | 65 689 184 (**65.7 ms**) | 84 764 598 | 646 516 | **−67.3%** | **−42.9%** | **−55.9%** |

5s profile re-runs (confirming stability, not formal 20x):

- ScanProfileAll: **109.0 ms**/op · 60.6 MB · 434k allocs  
- ExportOnly: **64.5 ms**/op · 84.8 MB · 647k allocs  

**Product `make run` wall** (from checklist / PR notes): pre ~**190 ms** → post ~**112–122 ms** scan (915 findings / 915+37 exports unchanged).

---

## Hotspot status (target residual items)

CPU compared via `go tool pprof -top -cum -focus=…` on pre (`Benchmark*-cpu.prof`) vs post (`scan-cpu-after.prof` / `export-cpu-after.prof`).

| Hotspot | Package / symbol | Pre (cum % of samples) | Post | Status |
|---------|------------------|------------------------|------|--------|
| **buildProjectSnapshot** / project walk | `bad_practices.buildProjectSnapshot` | **~31%** scan CPU | **~1.5%** | **gone / drastically reduced** (Once + walk filters) |
| **enclosingFunctionLines** / per-finding `ast.Inspect` | `export.enclosingFunctionLines` | **~18%** export CPU (`Walk` ~19%) | **0%** symbol; residual `ast.Inspect` ~**2.6%** (once-per-file span build) | **gone** as per-finding walk |
| **stripLineComment / codeLines** | `bad_practices.stripLineComment` / `codeLines` | flat **1.6%** / cum **~3.3%** scan | strip **0.25%** via `buildCodeLines`; `codeLines` not a multi-call hotspot | **reduced** (once per file) |
| **numberedLines / fmt.Sprintf** | `export.numberedLines` | numbered **~15%** · `fmt.Sprintf` **~8%** export | `numberedLinesCached` modest; `fmt.Sprintf` **~0.7%** | **reduced** |
| **catalogue lock / RuleIDs sort** | `RuleIDs` + `sort.Strings` / register path | small but visible (~0.5% with sort under RuleIDs) | `RuleIDs` **not in profile**; `Allows` ~**0.04%** | **gone** as measurable hotspot |

### Scan after — remaining cost shape

Dominated by real detector work, not the old infrastructure:

- `strings.Contains` / `Index` / `indexbytebody` (needle search across PERF/CWE/BP)
- Domain `Run` paths (`perf`, `bad_practices`, `cwe`) and fact builders
- Parse / `ast.Build` / residual `go/ast.Walk` inside detectors (~6% cum)
- GC: `gcDrain` ~**6.8%** cum (secondary on scan)

### Export after — remaining cost shape

- `ExportFindings` / `formatFindingBlock` / `functionWindow` (content assembly + whole-function windows)
- Disk: `os.WriteFile` / `RemoveAll` / syscalls (I/O bound portion)
- Parse for files still needing AST (`getParsed` / `go/parser`)
- **GC still large share:** `gcDrain` ~**24%** cum (see below)

---

## P3.2 GC note

| Metric | Export pre | Export post |
|--------|------------|-------------|
| `runtime.gcDrain` cum | ~**22.6%** | ~**24.2%** |
| `runtime.scanObjectsSmall` cum | ~**14%** | **not visible** (runtime/GC path samples shifted; `scanobject` ~21% cum) |
| Export wall (20x) | ~201 ms | ~**66 ms** |

**Conclusion for P3.2:** Absolute export time and allocs dropped a lot (allocs/op −56%, B/op −43%), but **GC mark work still dominates relative export CPU share** (`gcDrain` ~¼ of samples). That is expected when application work shrinks: GC percentage can hold or rise while absolute GC time falls.  

**Do not** change product GOGC defaults yet. Optional follow-on only if product wall is still GC-bound after residual needle/fact work: document CLI `GOGC` experiment notes, measure with same 20x benches, keep default unchanged unless win is clear and correctness-stable.

---

## Gates checklist mapping

| Gate | Result |
|------|--------|
| Formal `benchtime=20x` vs pre median | **Pass** — all three benches much faster + fewer allocs |
| pprof hotspots drop | **Pass** — snapshot, span Inspect, codeLines, numbered/sprintf, RuleIDs all reduced/gone |
| §12.4 findings/exports | Already green on branch (**915** / **915+37**) |
| P3.2 GC | **Still hot on export relative %** — open for optional tuning notes only |

---

## Commands used

```sh
export CGO_ENABLED=0
export GOSLOP_BENCH_SCAN_PATH=/home/chinmay/ChinmayPersonalProjects/gopdfsuit
# Formal
go test -run='^$' -bench=. -benchmem -benchtime=20x ./internal/bench/ \
  | tee /tmp/goslop-pprof/bench-after-p0p3-20x.txt
# CPU profiles
go test -run='^$' -bench=BenchmarkScanProfileAll -benchtime=5s -benchmem \
  -cpuprofile=/tmp/goslop-pprof/scan-cpu-after.prof ./internal/bench/
go test -run='^$' -bench=BenchmarkExportOnly -benchtime=5s -benchmem \
  -cpuprofile=/tmp/goslop-pprof/export-cpu-after.prof ./internal/bench/
go tool pprof -top -cum /tmp/goslop-pprof/scan-cpu-after.prof
go tool pprof -top -cum /tmp/goslop-pprof/export-cpu-after.prof
```

**Note:** `internal/bench/` was restored from `origin/perf/benchmarks-go126-pprof` for this measurement only (not committed on this branch).
