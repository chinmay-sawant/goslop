## Summary

Speed up product `make run` on the gopdfsuit reference corpus by ~**40%** (scan wall ~**190 ms → ~112 ms** mean) without changing §12.4 parity (915 findings, severity histogram, top rules, 915+37 exports). Implements pprof checklist P0–P2 and P3.1 dispatch cleanup.

---

## Motivation / context

- Plans: `plans/v0.0.1/perf-pprof-optimization-checklist.md`
- Evidence: prior `/tmp/goslop-pprof` profiles; product runs in `1.txt` (20× `make run`)
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

### Makefile

- Comment documenting optional timing loop:  
  `for i in {1..20}; do make run || break; done > 1.txt 2>&1`

---

## Impact

| Area | Impact |
|------|--------|
| **Performance** | `make run` scan wall ~**190 ms → ~112 ms** mean (~**41%** faster); min ~106 ms, max ~125 ms over 20 runs |
| **Memory** | Lower export/scan allocs (span cache, zero-copy Slice, shared codeLines) |
| **Behavior / correctness** | Unchanged: 915 findings; 10h/197i/312l/396m; top BP-1×181…; 915+37 exports |
| **API / CLI** | None (Makefile comment only for ops) |
| **Dependencies** | None |
| **Binary size / build time** | Negligible |

---

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| None | - |

---

## Test plan

- [x] `go test ./internal/export/ ./internal/lang/go/detectors/... ./tests/integration/ -count=1`
- [x] `make run` on gopdfsuit — 915 findings / 915+37 exports
- [x] 20× `make run` timing (`1.txt`) — mean ~112 ms scan
- [ ] `make bench BENCHTIME=20x` vs prior medians (optional follow-up)
- [ ] Re-profile pprof hotspots after merge (optional)

### Commands

```sh
make run
# optional timing:
# for i in {1..20}; do make run || break; done > 1.txt 2>&1
```

---

## Screenshots / sample output

**Before (approx.):** `scanned 78 files … in ~190ms` · 915 findings  

**After (from 20-run sample):**

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

**Summary:** mean ~**112 ms** scan (~**1.7×** vs ~190 ms baseline), parity intact.

---

## Related issues

- Relates to `plans/v0.0.1/perf-pprof-optimization-checklist.md` (P0–P2, P3.1)
- No open GitHub issue ID

---

## PR metadata checklist (author)

- [x] Self-assigned (`--assignee @me`)
- [x] Labels applied
- [x] Related issues filled
- [x] Filled body under `plans/PR/pr-perf-p0-p1-snapshot-export.md`

---

## Follow-ups (out of scope)

- P3 residual: needle/domain gates when index has no hits
- Cross-pack fact reuse (one AST walk for PERF/CWE/BP)
- P3.2 GC pressure notes after re-profile
- Formal `make bench` / pprof before-after attachment

---

## Reviewer checklist

- [ ] Behavior matches summary and test plan
- [ ] No unrelated changes in diff
- [ ] Public API / CLI changes documented
- [ ] New rules have fixture coverage when applicable
- [ ] PR has assignee and labels
- [ ] Related issues use correct Closes/Relates keywords
- [ ] No secrets or generated artifacts committed
