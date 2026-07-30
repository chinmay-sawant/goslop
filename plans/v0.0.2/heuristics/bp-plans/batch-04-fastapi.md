# v0.0.2 / BP-PY Batch 04 — FastAPI / Starlette heuristics

> **Parent:** `plans/v0.0.2/heuristics/bp-plans/README.md` — batch inventory + PR order  
> **Canonical ledger:** `plans/v0.0.2/heuristics/python-heuristics-bp.md` (#53)  
> **Catalogue:** `ruleset/python/bad-practices.json`  
> **Status:** planned — not started  
> **Estimated effort:** one PR (medium; 4 rules)  
> **Branch suggestion:** `feat/python-bp-batch-04-fastapi`  
> **Ledger rule:** mark `[x]` only with detector + unit hit/miss proof; record `make lint` / `make test` outcomes on the implement branch.

---

## Constraints (this PR)

- **Package:** `internal/lang/python/detectors/bad_practices/`
- **Pure-Go source heuristics** over `ParsedUnit.Source` (no Python AST / CGO)
- **`RegisterRule`** from `init()`; IDs **`BP-PY-*` only** (never bare Go `BP-<n>`)
- **Files ≤1500 lines preferred, hard max 2000 lines** — put FastAPI/Starlette rules in **`rules_fastapi.go`** (new)
- **Batchwise PR** for this file only (FastAPI batch)
- **Zero missing rules left unplanned** in the assigned set: `BP-PY-29`, `30`, `31`, `32`

---

## Overview

| Item | Evidence / target |
|------|-------------------|
| Assigned IDs | `BP-PY-29`, `BP-PY-30`, `BP-PY-31`, `BP-PY-32` (all missing) |
| Target file | **`rules_fastapi.go`** (new) |
| Tests | `scan_test.go` and/or `rules_fastapi_test.go` via `assertRule` hit/miss |
| Fixtures | optional `tests/fixtures/python/bp/BP-PY-2{9,0,1,2}-*.txt` |
| Related overlap | `BP-PY-39` (time.sleep in any async def) is **not** this batch; keep IDs distinct — `30` is FastAPI-route-scoped blocking I/O set |
| Inventory | `_inventory.json` lists all four under `missing` |

**Detection strategy:** source needles for FastAPI/Starlette imports and decorators (`@app.`, `@router.`, `APIRouter`, `FastAPI(`, `FileResponse`, `Depends(`). Async body windows for blocking calls. No full type inference for ORM return values — document heuristic limits for `31`.

---

## Executive Summary

### Why

- High-signal framework rules for async Python APIs: event-loop blocking, path traversal via `FileResponse`, mutable global deps, ORM leakage without `response_model`.
- Parent ledger Batch C deferred `30`/`32`; Batch E deferred `29`/`31` — this PR lands the full FastAPI catalogue slice in one place.

### Non-goals

- Proving all blocking I/O (socket, sync redis, etc.) — v0 needle list is finite
- Cross-file Depends graph
- Emitting CWE-22 / CWE-79 instead of `BP-PY-*`
- Implementing `BP-PY-39` here (async sleep global rule is a later batch)

### Dependency graph

```text
Phase 1 scaffold (done)
  └─ this batch: rules_fastapi.go
       Phase 1: BP-PY-30 blocking I/O + BP-PY-32 FileResponse (high sev)
       Phase 2: BP-PY-29 mutable global + BP-PY-31 response_model (medium)
       Phase 3: register tests + file-size + lint/test gates
```

---

## Phase 1: High-severity FastAPI / Starlette (`BP-PY-30`, `BP-PY-32`)

### 1.1 `BP-PY-30` FastAPI Blocking I/O In Async Route

| Field | Catalogue |
|-------|-----------|
| Severity | high |
| Category | FastAPI |
| detection_notes | Match `async def` route/dependency bodies calling `time.sleep`, `requests.*`, sync sqlalchemy session, `open().read` large, or `subprocess` without asyncio. Prefer `asyncio.sleep`, `httpx.AsyncClient`, `run_in_executor`, or `def` routes. |

- [ ] **Detection approach:**
  - Identify async route/dependency defs: `async def` under FastAPI-ish module (import/`FastAPI`/`APIRouter`/`@app.`/`@router.` needles) **or** any `async def` decorated with `@app.` / `@router.` / that takes `Depends`
  - In function body window (indent-based until next top-level/sibling def): flag calls matching blocking needles:
    - `time.sleep(`
    - `requests.get(` / `requests.post(` / `requests.put(` / `requests.patch(` / `requests.delete(` / `requests.request(` / `requests.Session`
    - `subprocess.` (`run`/`call`/`Popen`/…)
    - optional: `open(` + read patterns; sync `Session(` / `session.query` / `session.execute` when sqlalchemy markers present
  - Miss: `await asyncio.sleep(`; `httpx.AsyncClient`; sync `def` routes (not async)
  - Distinct from future `BP-PY-39`: this rule requires FastAPI/route context **or** documents broader async-def scope if v0 chooses module-wide async def + blocking needle (prefer route-scoped for lower noise)
- [ ] Path: `internal/lang/python/detectors/bad_practices/rules_fastapi.go` — `detectBPPY30`
- [ ] Register: `RegisterRule("BP-PY-30", detectBPPY30)` in `init()`
- [ ] Unit hit: FastAPI app + `@app.get` + `async def` body with `time.sleep(1)` or `requests.get(...)`
- [ ] Unit miss: same with `await asyncio.sleep(1)`
- [ ] Unit miss: sync `def` route with `time.sleep` (not this rule)
- [ ] Optional fixture hit/miss
- [ ] Severity **high** from metadata

### 1.2 `BP-PY-32` Starlette FileResponse User Path

| Field | Catalogue |
|-------|-----------|
| Severity | high |
| Category | FastAPI |
| detection_notes | Match `starlette.responses.FileResponse` or `StaticFiles` usage with path from path params/query without resolve+prefix check. Complements CWE-22. |

- [ ] **Detection approach:**
  - Match `FileResponse(` (and optional `StaticFiles`)
  - Fire when first/path argument expression text contains user-input needles: `request.`, path param name reuse, f-string / concatenation with variables, or bare name that is a function parameter (heuristic: param names used in `FileResponse(param)` for route with `{param}` or typed path params)
  - v0 practical rule: `FileResponse(` with non-literal path arg **and** (file has FastAPI/Starlette markers) **and** (arg involves f-string / `request` / `+` concat / `.format` / path param identifier from signature)
  - Miss: `FileResponse("/var/www/static/report.pdf")` constant; or path passed through `os.path.realpath` + prefix check needles (optional miss if too hard — document)
- [ ] Path: `rules_fastapi.go` — `detectBPPY32`
- [ ] Register `BP-PY-32`
- [ ] Unit hit: `@app.get("/files/{name}")` + `return FileResponse(name)` or `FileResponse(f"/data/{name}")`
- [ ] Unit miss: `FileResponse("/safe/fixed.pdf")`
- [ ] Optional fixture
- [ ] Severity **high**

### 1.3 Phase 1 validation

- [ ] Both high-sev rules registered with hit/miss tests
- [ ] `gofmt -w` on touched files

---

## Phase 2: Medium-severity FastAPI design (`BP-PY-29`, `BP-PY-31`)

### 2.1 `BP-PY-29` FastAPI Depends On Mutable Global

| Field | Catalogue |
|-------|-----------|
| Severity | medium |
| Category | FastAPI |
| detection_notes | Match `global`/`nonlocal` writes or module-level dict/list mutation inside route handlers/dependencies. Prefer request-scoped deps and proper stores. |

- [ ] **Detection approach:**
  - In FastAPI-ish modules, inside route/dependency functions: flag `global ` statements; `nonlocal ` writes; or mutations of module-level names (hard) — v0 minimum: `global foo` in route/dep body; optional `cache[key] =` when `cache` assigned at module level as `{}`/`[]`
  - Miss: pure function deps without global; request-state patterns
- [ ] Path: `rules_fastapi.go` — `detectBPPY29`
- [ ] Register `BP-PY-29`
- [ ] Unit hit: module-level `STORE = {}` + route/dep does `global STORE` or `STORE[k] = v`
- [ ] Unit miss: route uses local dict only
- [ ] Optional fixture
- [ ] Severity **medium**

### 2.2 `BP-PY-31` FastAPI response_model Disabled Unsafely

| Field | Catalogue |
|-------|-----------|
| Severity | medium |
| Category | FastAPI |
| detection_notes | Heuristic: FastAPI route returns SQLAlchemy/Django model instances without `response_model=` or `response_model_exclude`. Prefer Pydantic response models. |

- [ ] **Detection approach (heuristic, document limits):**
  - Find route decorators `@app.get/post/...` / `@router.get/...` **without** `response_model=` in decorator args
  - And function body `return` of names that look like ORM instances: `return db_user`, `return User.query...`, `return session.get(...)`, `return models.X` patterns — v0: fire when decorator lacks `response_model` **and** return expr matches ORM-ish needles (`session.query`, `.query.get`, `db.query`, SQLAlchemy model class call) **or** simpler v0: any FastAPI route decorator missing `response_model=` that returns a non-dict non-literal name when sqlalchemy/django models imported
  - Prefer precision: require both missing `response_model` and ORM return needle
  - Miss: `@app.get(..., response_model=UserOut)` or returns `UserOut(...)` / plain `dict`
- [ ] Path: `rules_fastapi.go` — `detectBPPY31`
- [ ] Register `BP-PY-31`
- [ ] Unit hit: `@app.get("/u/{id}")` without response_model + `return session.query(User).get(id)` (or similar)
- [ ] Unit miss: same with `response_model=UserOut`
- [ ] Optional fixture
- [ ] Severity **medium**; message notes heuristic

### 2.3 Phase 2 validation

- [ ] `29` and `31` registered + hit/miss
- [ ] No ID collision with Django/Flask rules

---

## Phase 3: Wire-up, file-size gate, closure

### 3.1 Registration & catalogue surface

- [ ] `TestBPRulesRegistered` (or fastapi-specific) includes `BP-PY-29`…`BP-PY-32`
- [ ] Each id: non-nil metadata, `PackBadPractice`, `BP-PY-` prefix only
- [ ] Collision guard unchanged

### 3.2 File-size gate (required after implementation)

- [ ] `wc -l internal/lang/python/detectors/bad_practices/*.go` — each file **≤1500 preferred**, **hard max 2000**
- [ ] Default: all four detectors live in **`rules_fastapi.go`**
- [ ] If tests bloat `scan_test.go` past comfort, use `rules_fastapi_test.go`
- [ ] Record line counts when closing the PR

### 3.3 Quality gates (required for non-docs)

- [ ] `gofmt -w` on all touched Go files
- [ ] `make lint` — leave unchecked until green; record outcome
- [ ] `make test` — leave unchecked until green; record outcome
- [ ] Optional: `go test ./internal/lang/python/detectors/bad_practices/ -count=1`

### 3.4 PR hygiene

- [ ] One PR for this batch (`batch-04-fastapi`)
- [ ] `Relates to #53` / `#51`
- [ ] Update `python-heuristics-bp.md` inventory rows for `29`–`32` when proven
- [ ] Update `_inventory.json` implemented/missing when batch lands

---

## Rule tracker (this batch)

| ID | Name | Sev | Status |
|----|------|-----|--------|
| BP-PY-29 | FastAPI Depends On Mutable Global | medium | [ ] |
| BP-PY-30 | FastAPI Blocking I/O In Async Route | high | [ ] |
| BP-PY-31 | FastAPI response_model Disabled Unsafely | medium | [ ] |
| BP-PY-32 | Starlette FileResponse User Path | high | [ ] |

**Coverage check:** assigned set fully planned — zero missing unplanned.

---

## Dependencies

| Depends on | Note |
|------------|------|
| Package scaffold | `RegisterRule`, scan loop, `pushAt`, facts bag |
| Common helpers | `isPythonTestFile`, `codeLinesFacts`, call-arg helpers in `common.go` |
| Catalogue | `ruleset/python/bad-practices.json` keys `BP-PY-29`…`32` |
| Parent README | `plans/v0.0.2/heuristics/bp-plans/README.md` |
| Sibling note | `BP-PY-39` (generic async sleep) may later overlap needles — keep messages/rule IDs distinct |
| Out of scope | Django (`batch-03`), templates/DB (`batch-05`) |

---

## References

- Catalogue: `ruleset/python/bad-practices.json`
- Tests pattern: `internal/lang/python/detectors/bad_practices/scan_test.go`
- Inventory: `plans/v0.0.2/heuristics/bp-plans/_inventory.json`
- Parent ledger: `plans/v0.0.2/heuristics/python-heuristics-bp.md`
- Skill: `plans/skills/phase-wise-checklist/SKILLS.md`
