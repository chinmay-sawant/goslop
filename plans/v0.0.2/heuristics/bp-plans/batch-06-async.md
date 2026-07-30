# Batch 06 — Async / threading BP-PY heuristics

> **Parent:** `plans/v0.0.2/heuristics/bp-plans/README.md` — batch index  
> **Canonical #53 ledger:** `plans/v0.0.2/heuristics/python-heuristics-bp.md`  
> **Issue:** [#53](https://github.com/chinmay-sawant/goslop/issues/53) expansion under epic [#51](https://github.com/chinmay-sawant/goslop/issues/51)  
> **Status:** not started — implementation checklist  
> **Estimated effort:** 1 PR (small–medium)  
> **PR policy:** single batchwise PR; title e.g. `python(bp): batch-06 async rules (BP-PY-38..40)`

---

## Overview

| Item | Value |
|------|-------|
| IDs (all must ship in this PR) | **BP-PY-38**, **BP-PY-39**, **BP-PY-40** |
| Category | Async (catalogue) |
| Target rule file | **`internal/lang/python/detectors/bad_practices/rules_async.go`** (new) |
| Do not grow | `rules_core.go` / `rules_security.go` / `rules_framework.go` for these IDs |
| Line budget | keep `rules_async.go` well under **1500 soft / 2000 hard**; package already ~2k lines total — new file required |
| Detection | source-pattern + `bpFacts` needles; pure-Go; no Python AST |
| Tests | hit/miss in `scan_test.go` or `rules_async_test.go` if test file size warrants split |
| Overlap note | BP-PY-39 overlaps FastAPI BP-PY-30 (`time.sleep` in async); **keep distinct rule IDs** — 39 is general async-def sleep, 30 is FastAPI route-scoped (other batch) |

### Catalogue contract

| ID | Name | Sev | Detection notes (v0) |
|----|------|-----|----------------------|
| BP-PY-38 | asyncio create_task Without Reference | medium | `create_task` / `ensure_future` as bare expression (return unused) |
| BP-PY-39 | time.sleep In Async Function | high | `time.sleep` inside `async def` body; prefer `await asyncio.sleep` |
| BP-PY-40 | threading Without Join Or Shutdown | low | `threading.Thread(...).start()` without later `.join` / documented daemon — review-only OK if message says so |

---

## Executive Summary

Add three async-domain detectors via `RegisterRule` in a new `rules_async.go`, wire needles in `facts.go` if useful, prove with unit hit/miss, then `make lint` + `make test`.

```text
rules_async.go init()
  ├─ RegisterRule("BP-PY-38", detectBPPY38)
  ├─ RegisterRule("BP-PY-39", detectBPPY39)
  └─ RegisterRule("BP-PY-40", detectBPPY40)
```

---

## Phase 1: Scaffold `rules_async.go`

### 1.1 File + registration

- [ ] Create `internal/lang/python/detectors/bad_practices/rules_async.go` (`package badpractices`)
- [ ] `init()` registers **BP-PY-38**, **BP-PY-39**, **BP-PY-40** only (no other IDs)
- [ ] Confirm `wc -l rules_async.go` stays under 1500 soft max after implementation
- [ ] Do **not** touch shipped detectors in `rules_core.go` / `rules_security.go` / `rules_framework.go` except shared helpers if truly necessary

### 1.2 Facts / needles

- [ ] Extend `bpNeedles` in `facts.go` with cheap prefilters useful to 38–40, e.g. `create_task`, `ensure_future`, `time.sleep`, `threading.`, `.start(`, `.join(` (only if used by detectors)
- [ ] Keep needle table ordered / documented; avoid duplicating entire source scans when `facts.has` can early-out

### 1.3 Metadata

- [ ] `MetadataForID("BP-PY-38")` / `39` / `40` non-nil after catalogue load path (existing `metadata.go`)
- [ ] Severities match catalogue: 38 **medium**, 39 **high**, 40 **low**
- [ ] Pack remains `PackBadPractice` via metadata / ID prefix

---

## Phase 2: `BP-PY-38` — create_task without reference

### 2.1 Detector behavior

- [ ] Implement `detectBPPY38` in `rules_async.go`
- [ ] Hit: bare call expressions such as `asyncio.create_task(...)`, `create_task(...)`, `asyncio.ensure_future(...)`, `ensure_future(...)` where the call is a **statement** (not assigned, not awaited into a name, not passed to `gather`/`wait` as stored task list in simple cases)
- [ ] Miss: `task = asyncio.create_task(...)`, `tasks.append(asyncio.create_task(...))`, `await asyncio.gather(asyncio.create_task(...))` if return is consumed (v0: assignment / attribute store / list append is enough)
- [ ] Message: store tasks and await/gather with exception handling (align catalogue description)
- [ ] Early-out when source lacks `create_task` and `ensure_future`

### 2.2 Proof

- [ ] Unit test hit: path e.g. `async_fire.py` with bare `asyncio.create_task(coro())` inside `async def`
- [ ] Unit test miss: `t = asyncio.create_task(coro())` then `await t` (or at least assignment)
- [ ] Assert finding `RuleID == "BP-PY-38"` only on hit
- [ ] Update `TestBPRulesRegistered` want-list to include `BP-PY-38`

---

## Phase 3: `BP-PY-39` — time.sleep in async def

### 3.1 Detector behavior

- [ ] Implement `detectBPPY39` in `rules_async.go`
- [ ] Hit: inside an `async def` body (indent/block heuristic from line scan), call `time.sleep(...)` or `sleep(...)` only when clearly bound to `time` import if cheap; v0 may flag `time.sleep(` anywhere under an enclosing `async def` indent
- [ ] Miss: `time.sleep` only under plain `def` (sync); miss `await asyncio.sleep(...)`
- [ ] Miss / skip: optional skip of pure test files if noise is high — document choice; prefer still flagging production async
- [ ] Message: blocks event loop; use `await asyncio.sleep`
- [ ] Do **not** emit `BP-PY-30` from this detector (separate FastAPI rule / other batch)

### 3.2 Proof

- [ ] Unit test hit: `async def handler():\n    time.sleep(1)\n`
- [ ] Unit test miss: sync `def f(): time.sleep(1)` does **not** fire BP-PY-39
- [ ] Unit test miss: `async def handler():\n    await asyncio.sleep(1)\n`
- [ ] Severity from metadata is **high**
- [ ] `TestBPRulesRegistered` includes `BP-PY-39`

---

## Phase 4: `BP-PY-40` — threading without join / shutdown

### 4.1 Detector behavior

- [ ] Implement `detectBPPY40` in `rules_async.go`
- [ ] Hit: `threading.Thread(...).start()` or `t = threading.Thread(...); t.start()` when **no** `.join(` appears in the same unit (file-level heuristic is OK for v0)
- [ ] Prefer flagging non-daemon threads; if `daemon=True` present, miss or softer message (document in message / comment)
- [ ] Miss: same unit has `t.join(` / `.join()` on the thread name when assignment is simple enough; or `ThreadPoolExecutor` shutdown patterns out of scope
- [ ] Severity **low**; message marks **review-only** heuristic if confidence is file-level
- [ ] Early-out without `threading.` / `.start(`

### 4.2 Proof

- [ ] Unit test hit: start thread, no join in file
- [ ] Unit test miss: `t.start()` then `t.join()` in same source
- [ ] Optional miss: `daemon=True` policy if implemented
- [ ] `TestBPRulesRegistered` includes `BP-PY-40`

---

## Phase 5: Integration within package

### 5.1 Catalogue + plugin surface

- [ ] `NewPythonBadPracticeScan().RuleIDs()` contains 38, 39, 40
- [ ] `MetadataFor` non-nil for each; no bare `BP-*` IDs introduced
- [ ] No changes required to `plugin.go` if detectors already wired via `detectors.All()` → BP scan (confirm still true)

### 5.2 Shared helpers

- [ ] Reuse `pushAt`, `codeLinesFacts` / line helpers from `common.go` rather than forking new finding APIs
- [ ] If async indent walking is shared between 38/39, keep helpers **local** to `rules_async.go` unless second domain needs them

### 5.3 Line-limit gate

- [ ] `wc -l` on every touched Go file under `bad_practices/` — none over **2000** hard; prefer split if any approaches **1500** soft
- [ ] If `scan_test.go` exceeds soft cap, move batch-06 tests to `rules_async_test.go`

---

## Phase 6: Validation gates (required for code)

### 6.1 Format + tests

- [ ] `gofmt -w` on all touched Go files
- [ ] `go test ./internal/lang/python/detectors/bad_practices/ -count=1` — all batch-06 tests green
- [ ] `make lint` — green; record command + outcome here on implement branch: ________
- [ ] `make test` — green; record: ________
- [ ] Optional: `CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop`

### 6.2 Parent / inventory sync

- [ ] Mark BP-PY-38/39/40 `[x]` in parent `python-heuristics-bp.md` inventory only with evidence
- [ ] Update `_inventory.json` `implemented` / `missing` when batch lands
- [ ] PR body under `plans/PR/v0.0.2/` (or link): `Relates to #53`, `Relates to #51` — do not `Closes #53` from this batch alone

---

## Dependencies

| Depends on | Note |
|------------|------|
| Batch 00 scaffold | `RegisterRule`, `PythonBadPracticeScan`, metadata, facts, common helpers **shipped** |
| Catalogue keys | `BP-PY-38`–`40` already in `ruleset/python/bad-practices.json` |
| Not blocked on | FastAPI BP-PY-30 batch (keep IDs distinct) |
| Out of scope | tree-sitter async CFG; cross-file task lifetime proof |

---

## Non-goals

- Full asyncio task graph / cancellation analysis  
- Detecting all forms of fire-and-forget (`loop.call_soon` etc.) beyond create_task/ensure_future  
- Replacing or implementing BP-PY-30 in this PR  
- Shipping BP-PY-14/48–50 here  

---

## Complete ID checklist (none deferred)

- [ ] **BP-PY-38** — detector + register + hit/miss + registered in `RuleIDs`
- [ ] **BP-PY-39** — detector + register + hit/miss + registered in `RuleIDs`
- [ ] **BP-PY-40** — detector + register + hit/miss + registered in `RuleIDs`

**Batch-06 ID list:** BP-PY-38, BP-PY-39, BP-PY-40
