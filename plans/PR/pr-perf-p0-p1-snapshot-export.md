## Summary

Speed up product `make run` on the **gopdfsuit reference corpus** by ~**40%** (scan wall ~**190 ms → ~112 ms** mean over 20 runs) without changing §12.4 parity (**915** findings, severity histogram, top rules, **915+37** exports). Implements the pprof checklist through **P0–P2**, **P3.1**, **CWE/PERF needle gates**, **shared AST facts (PERF+BP)**, and formal **bench/re-profile** evidence.

---

## Motivation / context

- Plans: `plans/v0.0.1/perf-pprof-optimization-checklist.md`
- Measurement: `plans/v0.0.1/perf-p0p3-measurement.md`
- Product timing: 20× `make run` sample (`1.txt`)
- Prior pprof: `/tmp/goslop-pprof/*`
- Issues: see **Related issues**

---

## Changes

### P0 / P1 (major wall win)

- BP `ProjectSnapshot`: `sync.Once` per absolute root (no thundering-herd rebuild)
- Export: function spans + line starts once per file; dual-export formats each finding once
- BP `codeLines`: build once in `buildFacts`, reuse via `codeLinesFacts`

### P2 (smaller incremental)

- Snapshot walk: skip tests/examples/dot/build dirs; stop content scan when flags complete; preallocated `bytes` needles
- `goparse.Tree.Slice`: zero-copy view over immutable `Source`
- Export early-out for package-level rules / sources without `func`

### P3.1 (dispatch)

- Immutable PERF/CWE/BP rule catalogues + cached sorted `RuleIDs()` (no per-file lock/copy/sort)
- Filter with `Allows` per unit (race-safe with shared detectors)

### Residual follow-ups (also in this PR)

#### Needle / domain gates

- **CWE:** ~169 table rules skip when `SourceIndex` has none of the rule’s group needles (same condition as `runNeedleRule` → FN-safe)
- **PERF:** optional gates on `RegisterRule`; high-hit rules wired in `zz_gates.go` (e.g. PERF-1/5/6/32, 16/21–23/28/35/46/47/101/122/186) against expanded `perfNeedles`

#### Cross-pack parse + AST reuse

- `goparse.TreeForUnit`: one parse pinned on `unit.Tree` for PERF / BP / taint
- `astfacts.Shared` on `unit.FactCache`: **one shared two-pass AST walk** for PERF + BP (no second `Inspect`)
- CWE structural still uses its own needle index (different table)

#### GOGC (P3.2)

- Re-profiled after opts: export allocs −56% but `gcDrain` still ~24% of export samples
- **No product GOGC change** (relative GC share is not a wall-time regression)

#### Bench harness

- Landed `internal/bench/` + Makefile `make bench` / `BENCHTIME` / `GOSLOP_BENCH_SCAN_PATH`

### Makefile / docs naming

- Optional timing-loop comment:  
  `for i in {1..20}; do make run || break; done > 1.txt 2>&1`
- Informal metrics-gate naming cleaned up (`make reference-metrics`, reference corpus wording)

---

## Impact

| Area | Impact |
|------|--------|
| **Performance** | `make run` scan ~**190 → ~112 ms** mean (~**41%** / ~**1.7×**). Formal 20x benches: Scan **−31%**, ScanAndExport **−50%**, ExportOnly **−67%** (see table below) |
| **Memory** | Large drop in B/op and allocs/op (scan B/op **−69%**, export allocs **−56%** on 20x benches) |
| **Behavior / correctness** | Unchanged product parity: **915** findings; **10h / 197i / 312l / 396m**; top BP-1×181…; **915+37** exports |
| **API / CLI** | `make bench` added; metrics gate remains `make reference-metrics` |
| **Dependencies** | None |
| **Binary size / build time** | Negligible |

### Formal benches (`benchtime=20x`, gopdfsuit)

| Benchmark | Before | After | Δ time | Δ B/op | Δ allocs |
|-----------|--------|-------|--------|--------|----------|
| **ScanProfileAll** | 164.2 ms · 199 MB · 1.19M | **112.7 ms** · 62 MB · 444k | **−31%** | **−69%** | **−63%** |
| **ScanAndExport** | 366.7 ms · 345 MB · 2.64M | **186.8 ms** · 146 MB · 1.09M | **−49%** | **−58%** | **−59%** |
| **ExportOnly** | 200.6 ms · 149 MB · 1.47M | **71.7 ms** · 85 MB · 647k | **−64%** | **−43%** | **−56%** |

### Hotspot status (pprof cum)

| Hotspot | Status |
|---------|--------|
| `buildProjectSnapshot` | ~31% → ~**1.5%** (gone / drastically reduced) |
| per-finding `enclosingFunctionLines` / `ast.Inspect` | **gone** (once-per-file spans) |
| `codeLines` / `stripLineComment` | **reduced** (once per file) |
| catalogue `RuleIDs` sort / register lock | **gone** as measurable cost |
| export `numberedLines` / `fmt.Sprintf` | **reduced** |

---

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| None for CLI detectors | - |
| Legacy Makefile metrics target name | Use `make reference-metrics` |

---

## Test plan

- [x] `go test ./internal/export/ ./internal/lang/go/detectors/... ./internal/engine/ ./internal/lang/go/goparse/ ./tests/integration/ -count=1`
- [x] `make run` on gopdfsuit — 915 findings / 915+37 exports
- [x] 20× `make run` timing — mean ~112 ms scan
- [x] Duration benches `benchtime=5s` — Scan 112.7 / ScanAndExport 186.8 / Export 71.7 ms
- [x] Fixed-iter 20x confirmation earlier same day
- [x] pprof re-profile — hotspots in `perf-p0p3-measurement.md`
- [x] GOGC decision: leave product default unchanged
- [x] Docs: `documents/make-run.md` performance section

### Commands

```sh
make run
# optional timing:
# for i in {1..20}; do make run || break; done > 1.txt 2>&1
make bench BENCHTIME=5s
make reference-metrics
```

---

## Screenshots / sample output

**Before (approx.):** `scanned 78 files … in ~190ms` · 915 findings  

**After (20-run product sample):**

```
scanned 78 files (28042 lines) in 107.4ms   # min-ish
scanned 78 files (28042 lines) in 112.3ms
…
scanned 78 files (28042 lines) in 124.5ms   # max
915 findings
  severity: 10 high, 197 info, 312 low, 396 medium
  top rules: BP-1 ×181, PERF-6 ×94, PERF-32 ×59, BP-5 ×50, PERF-230 ×44
exported 915 context file(s) …; exported 37 chunk file(s) …
```

**Summary:** mean ~**112 ms** scan (~**1.7×** vs ~190 ms), parity intact.

---

## Related issues

- Relates to `plans/v0.0.1/perf-pprof-optimization-checklist.md`
- Measurement: `plans/v0.0.1/perf-p0p3-measurement.md`
- No open GitHub issue ID

---

## PR metadata checklist (author)

- [x] Self-assigned (`--assignee @me`)
- [x] Labels applied
- [x] Related issues filled
- [x] Filled body under `plans/PR/pr-perf-p0-p1-snapshot-export.md`

---

## Follow-ups (still open / optional)

- More PERF gates only where FN-safe (many rules still ungated by design)
- Shared CWE fact bag only if needle tables can be unified without parity risk
- GOGC CLI experiment notes **only if** wall is still GC-bound after further residual work
- Optional: attach raw `.prof` artifacts outside git

---

## Reviewer checklist

- [ ] Behavior matches summary and test plan
- [ ] No unrelated changes in diff
- [ ] Public API / CLI changes documented
- [ ] New rules have fixture coverage when applicable
- [ ] PR has assignee and labels
- [ ] Related issues use correct Closes/Relates keywords
- [ ] No secrets or generated artifacts committed
