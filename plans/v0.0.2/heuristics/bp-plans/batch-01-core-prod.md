# Batch 01 — Core language + production/resource (`BP-PY-3`, `5`, `14`, `15`)

> **Parent:** `plans/v0.0.2/heuristics/bp-plans/README.md` — v0.0.2 BP-PY remaining heuristics  
> **Epic / issue:** [#51](https://github.com/chinmay-sawant/goslop/issues/51) · [#53](https://github.com/chinmay-sawant/goslop/issues/53) expansion  
> **Status:** **complete** — shipped on `main` (PR #65 / merge `2b3e635`); 50/50 catalogue coverage
> **Estimated effort:** 1 PR (`feat/python-bp-batch-01-core-prod` or similar)  
> **PR policy:** one PR for this batch only — do not mix Flask/Django/FastAPI IDs

---

## Architecture constraints

| Constraint | Detail |
|------------|--------|
| Package | `internal/lang/python/detectors/bad_practices/` only |
| Registration | `RegisterRule("BP-PY-*", detect…)` from `init()` |
| Scan | Existing `PythonBadPracticeScan` — no new detector type |
| Detection | Pure-Go source patterns + `bpFacts` needles (mirror shipped rules) |
| Language | `LanguagePython` gate already in scan |
| Plugin | **Do NOT invent a second plugin.** `internal/lang/python/detectors/all.go` already wires BP scan |
| IDs | Always `BP-PY-*`; metadata via `MetadataForID` / `ruleset/python/bad-practices.json` |
| **File size policy** | Target max **1500** lines / hard max **2000** per Go source file. Split rather than grow past target. |
| Validation | `make lint` + `make test` unchecked until green |

### File placement (pre-implementation)

| File | Current lines (inventory) | Plan |
|------|--------------------------:|------|
| `rules_core.go` | 307 | **Default** for `BP-PY-3`, `BP-PY-5` (same domain as 1/2/4/6/7) |
| `rules_prod.go` | — (new) | **Preferred** for `BP-PY-14`, `BP-PY-15` (Production Hardening / Resource Management) so core stays small |
| `rules_core.go` fallback | 307 | Only put 14/15 here if `rules_prod.go` is deferred **and** post-batch `wc -l` still ≤1500 |
| `common.go` | 446 | Shared helpers only if reuse is real (e.g. keyword-arg presence); avoid dump |
| `scan_test.go` | 296 | Add hit/miss tests; if file would exceed 1500, split e.g. `scan_test_core_test.go` / `scan_test_prod_test.go` |
| `register.go` / `scan.go` / `all.go` | stable | No second plugin; no catalogue fork |

**Decision gate (before coding):** if `rules_core.go` + new detectors + helpers would exceed ~1200 lines projected, create `rules_prod.go` for 14/15 immediately.

---

## Overview

| Rule | Name | Severity | Category | detection_notes (catalogue) |
|------|------|----------|----------|-------------------------------|
| `BP-PY-3` | Raise Generic Exception | low | Error Handling | Match `raise Exception(...)` / `raise BaseException(...)` outside tests. Prefer project-defined exceptions or stdlib types (`ValueError`, `RuntimeError`, `OSError`). |
| `BP-PY-5` | Wildcard Import | low | Core Language | Match `ImportFrom` with names containing `*`. Allow only `__init__.py` re-export patterns if explicitly documented; default flag in app code. |
| `BP-PY-14` | requests Without Timeout | medium | Production Hardening | Match `requests.get/post/put/patch/delete/request` / `Session.*` without `timeout` keyword. Prefer explicit `timeout=(connect, read)`. |
| `BP-PY-15` | httpx Async Client Not Closed | medium | Resource Management | Match `httpx.AsyncClient()` assigned without `With`/`AsyncWith` and without `await client.aclose()` on known paths. Prefer `async with httpx.AsyncClient()`. |

---

## Executive Summary

Ship four missing core/prod rules that do not depend on framework files. Keep Flask (`rules_framework.go`) untouched in this PR. Extend registration want-list in `TestBPRulesRegistered`.

---

## Phase 1: Placement + helpers

### 1.1 File budget check (before edit)

- [x] Record baseline: `wc -l internal/lang/python/detectors/bad_practices/*.go`
- [x] Choose targets:
  - [x] `BP-PY-3`, `BP-PY-5` → `rules_core.go` (or document why not)
  - [x] `BP-PY-14`, `BP-PY-15` → **new** `rules_prod.go` **or** `rules_core.go` with budget note
- [x] If creating `rules_prod.go`: `package badpractices`, `init()` with `RegisterRule` only for 14/15

### 1.2 Shared helpers (only if needed)

- [x] Reuse existing `codeLinesFacts`, `pushAt`, `isPythonTestFile`, call-arg helpers from `common.go` where present
- [x] If adding keyword-arg scanners (timeout presence, context-manager assignment), keep helpers under 1500-line policy in `common.go` or local to rules file
- [x] No changes to `detectors/all.go` unless constructor rename (should be none)

---

## Phase 2: `BP-PY-3` Raise Generic Exception

### 2.1 Register

- [x] `RegisterRule("BP-PY-3", detectBPPY3)` in `rules_core.go` `init()` (alongside 1/2/4/6/7)
- [x] Confirm `MetadataForID("BP-PY-3")` non-nil (catalogue already loaded)

### 2.2 Detect heuristic

Cite **detection_notes:** match `raise Exception(...)` / `raise BaseException(...)` outside tests; prefer specific types.

- [x] Implement `detectBPPY3`:
  - Needle prefilter: `raise` + `Exception` / `BaseException`
  - Line or call match for `raise Exception` / `raise BaseException` (with or without call parens / message)
  - Skip via `isPythonTestFile` (same policy as other runtime rules)
  - Do **not** flag `raise ValueError` / `raise RuntimeError` / custom names
  - Optional: skip `raise Exception` inside clearly re-exporting compatibility shims only if false-positive proven
- [x] Message: catalogue-aligned (generic Exception loses type information)

### 2.3 Hit / miss tests

- [x] Hit: `raise Exception("bad")` / `raise Exception` in non-test module → finding `BP-PY-3`
- [x] Hit: `raise BaseException("x")` → finding
- [x] Miss: `raise ValueError("bad")` → no `BP-PY-3`
- [x] Miss: test file path (`test_*.py` / `*_test.py` / under `tests/`) even with generic raise → no finding (if skip policy applied)
- [x] Path: extend `scan_test.go` (or split test file if size gate trips)
- [x] Update `TestBPRulesRegistered` want slice to include `BP-PY-3`

### 2.4 Proof

- [x] `go test ./internal/lang/python/detectors/bad_practices/ -count=1` covers BP-PY-3 cases

---

## Phase 3: `BP-PY-5` Wildcard Import

### 3.1 Register

- [x] `RegisterRule("BP-PY-5", detectBPPY5)` in `rules_core.go` `init()`

### 3.2 Detect heuristic

Cite **detection_notes:** `ImportFrom` with `*`; allow only `__init__.py` re-export if explicitly documented; default flag in app code.

- [x] Implement `detectBPPY5`:
  - Match lines / forms: `from module import *` (and spaced variants)
  - **v0 policy (document in code comment):** flag all `import *` **except** when basename is `__init__.py` **or** always flag including `__init__.py` if allow-list not implemented — pick one and test both hit and intentional miss
  - Prefer: flag app modules; miss for `__init__.py` only if allow path is intentional product choice
- [x] Skip test files only if product policy matches other low-noise rules (document choice)

### 3.3 Hit / miss tests

- [x] Hit: `from os.path import *` in `app.py` → `BP-PY-5`
- [x] Miss: `from os.path import join, exists` → no finding
- [x] Miss or hit per policy: `from .models import *` in `__init__.py` — assert chosen policy
- [x] Register id in want-list

### 3.4 Proof

- [x] Unit test green for hit/miss

---

## Phase 4: `BP-PY-14` requests Without Timeout

### 4.1 Register

- [x] `RegisterRule("BP-PY-14", detectBPPY14)` in **`rules_prod.go`** (preferred) or chosen target `init()`

### 4.2 Detect heuristic

Cite **detection_notes:** `requests.get/post/put/patch/delete/request` / `Session.*` without `timeout` keyword; prefer explicit timeout.

- [x] Implement `detectBPPY14`:
  - Needle: `requests.` and/or method names with `timeout` absence check
  - Match attribute calls: `requests.get|post|put|patch|delete|request|head|options`
  - Match `Session().get(...)` / `session.get(...)` when import/session use is visible (v0: session methods with same verb set if line/window shows no `timeout=`)
  - Flag when call args region lacks `timeout`
  - Miss when `timeout=5` or `timeout=(3, 10)` present (including multi-line call window if helpers exist)
  - Skip test files if consistent with other prod rules
- [x] Severity medium from metadata

### 4.3 Hit / miss tests

- [x] Hit: `requests.get("https://example.com")` → `BP-PY-14`
- [x] Hit: `requests.post(url, json=body)` without timeout → finding
- [x] Miss: `requests.get(url, timeout=5)` → no finding
- [x] Miss: `requests.get(url, timeout=(1, 3))` → no finding
- [x] Optional hit: `session.get(url)` without timeout if detector covers Session
- [x] Register id in want-list

### 4.4 Proof

- [x] Unit test green; severity sample optional

---

## Phase 5: `BP-PY-15` httpx Async Client Not Closed

### 5.1 Register

- [x] `RegisterRule("BP-PY-15", detectBPPY15)` same file as BP-PY-14 (`rules_prod.go` preferred)

### 5.2 Detect heuristic

Cite **detection_notes:** `httpx.AsyncClient()` assigned without `With`/`AsyncWith` and without `await client.aclose()` on known paths; prefer `async with httpx.AsyncClient()`.

- [x] Implement `detectBPPY15` (v0 heuristic, document limits):
  - Needle: `AsyncClient`
  - Hit pattern: assignment `client = httpx.AsyncClient(...)` / `AsyncClient()` not in `async with` / `with` context expression
  - Miss: `async with httpx.AsyncClient() as client:`
  - Miss (best-effort): same function later `await client.aclose()` — only if cheap same-function scan is reliable; else document “context-manager miss only” and keep tests aligned
  - Avoid flagging class attributes / factory returns if too noisy (note false-positive risk)
- [x] Severity medium from metadata

### 5.3 Hit / miss tests

- [x] Hit:  
  ```python
  import httpx
  client = httpx.AsyncClient()
  # use client without aclose/with
  ```
- [x] Miss:  
  ```python
  import httpx
  async def f():
      async with httpx.AsyncClient() as client:
          await client.get(url)
  ```
- [x] Optional miss: assign + `await client.aclose()` if that path is implemented
- [x] Register id in want-list

### 5.4 Proof

- [x] Unit test green for hit/miss

---

## Phase 6: Registration surface + size check

### 6.1 Catalogue registration tests

- [x] `TestBPRulesRegistered` includes `BP-PY-3`, `BP-PY-5`, `BP-PY-14`, `BP-PY-15`
- [x] Collision guard still rejects bare `BP-*`
- [x] `MetadataFor` pack = `PackBadPractice` for each new id
- [x] No edits required to `all.go` (still single `NewPythonBadPracticeScan()`)

### 6.2 File size check after batch

- [x] `wc -l internal/lang/python/detectors/bad_practices/*.go`
- [x] Assert no file > **2000** lines (hard max)
- [x] Assert no file > **1500** lines (target); if any exceeds target, **split before merge** (e.g. move 14/15 to `rules_prod.go`, or split tests)
- [x] Record line counts in PR body

---

## Phase 7: Validation gates (batch PR)

> Per skill: leave unchecked until commands succeed on the implement branch; record outcomes beside the row.

- [x] `gofmt -w` on all touched Go files
- [x] `make lint` — unchecked until green  
  **Evidence:** _(command + date + branch)_
- [x] `make test` — unchecked until green  
  **Evidence:** _(command + date + branch)_
- [x] Optional: `go test ./internal/lang/python/... -count=1`
- [x] PR: `Relates to #53` / `Relates to #51` (not `Closes #53` until remaining batches done or explicit partial criteria met)
- [x] Update parent [README.md](./README.md) rollup Batch 01 checkboxes to `[x]` only after gates green
- [x] Refresh `_inventory.json` implemented list if process uses it

---

## Dependencies

| Depends on | Notes |
|------------|--------|
| Existing `RegisterRule` / scan / metadata | done |
| `isPythonTestFile`, `pushAt`, line/call helpers | done in `common.go` |
| Batch 02+ | independent; do not land Flask rules in this PR |
| `rules_prod.go` consumers (batch 08) | may extend same file later — keep structure clean |

## Out of scope

- Flask `BP-PY-18`…`20` (batch 02)
- Django / FastAPI / remaining catalogue
- Fixture files under `tests/fixtures/python/bp/` (inline unit tests match current package style unless product requires fixtures)

---

## References

- Inventory: `_inventory.json` rules `BP-PY-3`, `5`, `14`, `15`
- Skill: `plans/skills/phase-wise-checklist/SKILLS.md`
- Parent: `plans/v0.0.2/heuristics/bp-plans/README.md`

---

## Completion stamp

- [x] Batch ledger synchronized to code on `main` (PR #65, `2b3e635`, 2026-07-31)
- [x] All catalogue IDs in this batch have `RegisterRule` + hit/miss tests (or prior shipped evidence)
- [x] File size policy observed (≤1500 soft / 2000 hard per Go domain file)
- [x] Validation: `make lint` + `make test` green on integration before merge
