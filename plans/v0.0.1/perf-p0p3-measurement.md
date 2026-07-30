# P0–P3.1 measurement (post-change re-profile)

> **Date:** 2026-07-30  
> **Branch:** `perf/p0-p1-snapshot-export` (P0 + P1 + P2 + P3.1 code landed)  
> **Corpus:** gopdfsuit (`GOSLOP_BENCH_SCAN_PATH=/home/chinmay/ChinmayPersonalProjects/gopdfsuit`)  
> **Toolchain:** Go 1.25.0 · `CGO_ENABLED=0`  
> **Harness:** `./internal/bench/` · formal `benchtime=20x` + duration `benchtime=5s`  
> **Artifacts:** `/tmp/goslop-pprof/bench-after-p0p3-20x.txt`, `bench-after-5s.txt`, `scan-cpu-after.prof`, `export-cpu-after.prof`

---

## Bench table (before → after)

Pre-change medians from `/tmp/goslop-pprof/Benchmark*-bench.txt` (same machine/corpus).  
**Primary after (duration-based):** `go test -run='^$' -bench=. -benchmem -benchtime=5s ./internal/bench/`  
(Confirming fixed-iter 20x in parentheses where useful.)

| Benchmark | Before | After **5s** | Δ time | Δ B/op | Δ allocs |
|-----------|--------|--------------|--------|--------|----------|
| **ScanProfileAll** | **164.2 ms** · 199 MB · 1.19M | **112.7 ms** · 62.0 MB · 444k | **−31%** | **−69%** | **−63%** |
| **ScanAndExport** | **366.7 ms** · 345 MB · 2.64M | **186.8 ms** · 145.9 MB · 1.09M | **−49%** | **−58%** | **−59%** |
| **ExportOnly** | **200.6 ms** · 149 MB · 1.47M | **71.7 ms** · 84.8 MB · 647k | **−64%** | **−43%** | **−56%** |

20x fixed-iter (earlier same day) for comparison: Scan **113.7 ms** · ScanAndExport **181.7 ms** · Export **65.7 ms** — agrees with 5s within noise.

**Product `make run` scan wall** (gopdfsuit, profile all + exports, no-cache):

| | Scan wall | Findings / exports |
|--|----------:|--------------------|
| Pre-opt | ~**190 ms** | 915 / 915+37 |
| Post-opt (20-run means) | ~**112–117 ms** | 915 / 915+37 |
| Best observed | **99.0 ms** | 915 / 915+37 |

Parity unchanged: severity 10/197/312/396; top BP-1×181, PERF-6×94, PERF-32×59, BP-5×50, PERF-230×44.
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
