# v0.0.2 / BP-PY Batch 05 — Templates + Database heuristics

> **Parent:** `plans/v0.0.2/heuristics/bp-plans/README.md` — batch inventory + PR order  
> **Canonical ledger:** `plans/v0.0.2/heuristics/python-heuristics-bp.md` (#53)  
> **Catalogue:** `ruleset/python/bad-practices.json`  
> **Status:** planned — not started  
> **Estimated effort:** one PR (medium; 5 rules across two domains)  
> **Branch suggestion:** `feat/python-bp-batch-05-templates-db`  
> **Ledger rule:** mark `[x]` only with detector + unit hit/miss proof; record `make lint` / `make test` outcomes on the implement branch.

---

## Constraints (this PR)

- **Package:** `internal/lang/python/detectors/bad_practices/`
- **Pure-Go source heuristics** over `ParsedUnit.Source` (no Python AST / CGO)
- **`RegisterRule`** from `init()`; IDs **`BP-PY-*` only** (never bare Go `BP-<n>`)
- **Files ≤1500 lines preferred, hard max 2000 lines** — split as needed:
  - **`rules_templates.go`** — `BP-PY-33`, `BP-PY-34`
  - **`rules_db.go`** — `BP-PY-35`, `BP-PY-36`, `BP-PY-37`
- **Batchwise PR** for this file only (templates + DB)
- **Zero missing rules left unplanned** in the assigned set: `BP-PY-33`, `34`, `35`, `36`, `37`

---

## Overview

| Item | Evidence / target |
|------|-------------------|
| Templates | `BP-PY-33` Jinja2 autoescape off; `BP-PY-34` Markup / `\|safe` |
| Database | `BP-PY-35` SQLAlchemy `text` f-string; `BP-PY-36` Session not closed; `BP-PY-37` DB-API `%` / f-string execute |
| Target files | **`rules_templates.go`**, **`rules_db.go`** (both new) |
| Tests | `scan_test.go` and/or `rules_templates_test.go` + `rules_db_test.go` |
| Fixtures | optional under `tests/fixtures/python/bp/` |
| Non-`.py` templates | v0 scans Python units only; `|safe` in `.html`/`.j2` is **out of band** unless a non-py path appears later — implement Python-side `Markup` / string templates in `.py` first |
| Inventory | all five currently `missing` in `_inventory.json` |
| Related | Django `BP-PY-24` (raw SQL) is batch-03; keep SQLAlchemy/DB-API rules here |

**Detection strategy:** pure-Go call-site needles. SQLi-adjacent rules (`35`,`37`) favor high precision on f-string / format / `%` of SQL. Resource rule `36` mirrors `open`/`with` style used by `BP-PY-7` / httpx patterns.

---

## Executive Summary

### Why

- Completes parent Batch D (`33`,`35`) plus deferred template/DB siblings (`34`,`36`,`37`).
- XSS-adjacent Jinja2 misconfig and SQL injection construction are high product signal.
- Domain split (`rules_templates.go` / `rules_db.go`) keeps review and file-size gates clean.

### Non-goals

- Full Jinja template file scanning (non-Python units)
- Proving all session close paths (CFG / finally on every branch)
- Emitting CWE-79 / CWE-89 rule IDs (BP stream only)
- Django ORM raw SQL (`BP-PY-24` belongs to batch-03)

### Dependency graph

```text
Phase 1 scaffold (done)
  └─ this batch
       Phase 1: rules_templates.go (33, 34)
       Phase 2: rules_db.go (35, 37 high-sev SQL construction)
       Phase 3: rules_db.go (36 session lifecycle)
       Phase 4: file-size + lint/test + inventory
```

---

## Phase 1: Templates — Jinja2 (`BP-PY-33`, `BP-PY-34`)

### 1.1 `BP-PY-33` Jinja2 autoescape Disabled

| Field | Catalogue |
|-------|-----------|
| Severity | high |
| Category | Templates |
| detection_notes | Match `jinja2.Environment(autoescape=False)` or `select_autoescape` not used for `.html` templates. Complements CWE-79. |

- [ ] **Detection approach:**
  - Match `Environment(` / `jinja2.Environment(` with `autoescape=False` in args (single-line regex + multi-line call window via `callArgsRegion` if available)
  - Optional: `Environment(` without any `autoescape=` when constructing HTML environments — noisier; v0 prioritize explicit `autoescape=False`
  - Miss: `Environment(autoescape=True)`, `Environment(autoescape=select_autoescape(...))`, `select_autoescape(['html', 'xml'])`
- [ ] Path: `internal/lang/python/detectors/bad_practices/rules_templates.go` — `detectBPPY33`
- [ ] Register: `RegisterRule("BP-PY-33", detectBPPY33)` in `init()`
- [ ] Unit hit: `jinja2.Environment(autoescape=False)`
- [ ] Unit miss: `jinja2.Environment(autoescape=True)` or `autoescape=select_autoescape(['html','xml'])`
- [ ] Optional fixture hit/miss
- [ ] Severity **high**

### 1.2 `BP-PY-34` Jinja2 Markup Or safe Filter On Variables

| Field | Catalogue |
|-------|-----------|
| Severity | high |
| Category | Templates |
| detection_notes | Match `Markup(...)`, `\|safe`, `\|forceescape` misuse in templates or Python. Flag non-literal arguments to `Markup`. |

- [ ] **Detection approach:**
  - Python: `Markup(` with first arg not a string literal (reuse `firstCallArg` / `isStringLiteral`)
  - Python string templates / Jinja embedded in `.py`: needle `|safe` / `| safe` on non-static context if present in source
  - Optional: `markupsafe.Markup`
  - Miss: `Markup("<b>ok</b>")` literal; pure `|e` escape use
  - **Limit:** `.html`/`.j2` files not in Python unit stream — document as non-goal for v0; Python-side only
- [ ] Path: `rules_templates.go` — `detectBPPY34`
- [ ] Register `BP-PY-34`
- [ ] Unit hit: `Markup(user_html)` / `Markup(request.args['x'])`
- [ ] Unit miss: `Markup("<br/>")`
- [ ] Optional: hit on Python source containing `{{ x|safe }}` if easy; else `[~]` note for template files
- [ ] Optional fixture
- [ ] Severity **high**

### 1.3 Phase 1 validation

- [ ] Both template rules registered with hit/miss tests
- [ ] `gofmt -w` on `rules_templates.go` (+ tests)

---

## Phase 2: Database — SQL construction (`BP-PY-35`, `BP-PY-37`)

### 2.1 `BP-PY-35` SQLAlchemy text With F-String

| Field | Catalogue |
|-------|-----------|
| Severity | high |
| Category | Database |
| detection_notes | Match `text(f"...")`, `text("...".format(...))`, or execute with concatenated SQL. Prefer `text("... WHERE id = :id")` with bindparams. Complements CWE-89. |

- [ ] **Detection approach:**
  - Match `text(` call (sqlalchemy `text`) when first arg is f-string, or chained `.format(`, or string `%` format, or `+` concat of SQL fragments
  - Prefer identifier context: `from sqlalchemy` / `sqlalchemy.text` / bare `text(` when file has sqlalchemy import needles (reduce false positives on other `text(` helpers — if ambiguous, require `sqlalchemy` import or `sa.text`)
  - Miss: `text("SELECT * FROM t WHERE id = :id")` + bindparams
- [ ] Path: `internal/lang/python/detectors/bad_practices/rules_db.go` — `detectBPPY35`
- [ ] Register `BP-PY-35`
- [ ] Unit hit: `from sqlalchemy import text` + `text(f"SELECT * FROM users WHERE id = {uid}")`
- [ ] Unit hit: `text("SELECT * FROM users WHERE id = {}".format(uid))`
- [ ] Unit miss: `text("SELECT * FROM users WHERE id = :id")`
- [ ] Optional fixture
- [ ] Severity **high**
- [ ] Do not double-count as Django `BP-PY-24` (different APIs)

### 2.2 `BP-PY-37` DB-API Cursor Execute With Percent Format

| Field | Catalogue |
|-------|-----------|
| Severity | high |
| Category | Database |
| detection_notes | Match `execute("...%s..." % (...))` or `execute(f"...")` patterns. Prefer `execute(sql, params)` binding. Complements CWE-89. |

- [ ] **Detection approach:**
  - Match `.execute(` where first arg is f-string, or string `%` formatting expression, or `.format(` on SQL string
  - Scope: any `.execute(` is OK for v0 (DB-API + many drivers); optional require cursor/connection naming (`cur.execute`, `cursor.execute`, `conn.execute`)
  - Miss: `cursor.execute("SELECT * FROM t WHERE id = %s", (uid,))` — static SQL + params tuple/list as second arg
  - Coordinate with SQLAlchemy: `session.execute(text(f"..."))` may fire both `35` and `37` if nested — prefer fire `35` on `text(` and `37` on bare execute f-string; avoid duplicate messages on same byte offset if both match same call (document policy: both OK or suppress one)
- [ ] Path: `rules_db.go` — `detectBPPY37`
- [ ] Register `BP-PY-37`
- [ ] Unit hit: `cursor.execute(f"SELECT * FROM t WHERE id = {uid}")`
- [ ] Unit hit: `cursor.execute("SELECT * FROM t WHERE id = %s" % (uid,))`
- [ ] Unit miss: `cursor.execute("SELECT * FROM t WHERE id = %s", (uid,))`
- [ ] Optional fixture
- [ ] Severity **high**

### 2.3 Phase 2 validation

- [ ] `35` and `37` hit/miss green
- [ ] Overlap policy with each other and with Django `24` documented in code comments if needed

---

## Phase 3: Database — session lifecycle (`BP-PY-36`)

### 3.1 `BP-PY-36` SQLAlchemy Session Not Closed

| Field | Catalogue |
|-------|-----------|
| Severity | medium |
| Category | Database |
| detection_notes | Match `Session()` / `SessionLocal()` without `with`/`contextlib.closing` and without `.close()` on exit paths. Prefer session scope context managers. |

- [ ] **Detection approach:**
  - Match assignment `session = Session(` / `SessionLocal(` / `sessionmaker(...)(` patterns
  - Fire when same function (indent window) has no `with Session` / `with SessionLocal` / `with sessionmaker` and no `.close()` call on that name
  - Miss: `with Session() as session:` / `with SessionLocal() as session:`; or `session.close()` present in function
  - False positives: sessions closed in callers — accept v0 noise; skip tests via `isPythonTestFile` if desired
- [ ] Path: `rules_db.go` — `detectBPPY36`
- [ ] Register `BP-PY-36`
- [ ] Unit hit: `session = SessionLocal()` + use without close/with
- [ ] Unit miss: `with SessionLocal() as session:` body
- [ ] Unit miss: `session = SessionLocal()` + `session.close()` before return
- [ ] Optional fixture
- [ ] Severity **medium**

### 3.2 Phase 3 validation

- [ ] `36` registered + hit/miss
- [ ] Comment that full CFG close proof is out of scope

---

## Phase 4: Wire-up, file-size gate, closure

### 4.1 Registration & catalogue surface

- [ ] Registered set includes `BP-PY-33`…`BP-PY-37`
- [ ] Metadata non-nil for each; pack `PackBadPractice`; `BP-PY-` prefix only
- [ ] Collision guard unchanged

### 4.2 File-size gate (required after implementation)

- [ ] `wc -l internal/lang/python/detectors/bad_practices/*.go`
- [ ] Each production `.go` file **≤1500 preferred**, **hard max 2000**
- [ ] Confirm split:
  - `rules_templates.go` holds `33`–`34` only (plus small helpers if private)
  - `rules_db.go` holds `35`–`37`
- [ ] If either file nears limit mid-implementation, extract shared SQL-arg helpers into `common.go` only if still under limit there (`common.go` baseline ~446 lines)
- [ ] Split tests if `scan_test.go` exceeds comfort
- [ ] Record line counts when closing the PR

### 4.3 Quality gates (required for non-docs)

- [ ] `gofmt -w` on all touched Go files
- [ ] `make lint` — leave unchecked until green; record outcome
- [ ] `make test` — leave unchecked until green; record outcome
- [ ] Optional: `go test ./internal/lang/python/detectors/bad_practices/ -count=1`

### 4.4 PR hygiene

- [ ] One PR for this batch (`batch-05-templates-db`)
- [ ] `Relates to #53` / `#51`
- [ ] Update `python-heuristics-bp.md` Batch D rows + inventory for `33`–`37` when proven
- [ ] Update `_inventory.json` implemented/missing when batch lands

---

## Rule tracker (this batch)

| ID | Name | Sev | Domain | Target file | Status |
|----|------|-----|--------|-------------|--------|
| BP-PY-33 | Jinja2 autoescape Disabled | high | Templates | `rules_templates.go` | [ ] |
| BP-PY-34 | Jinja2 Markup Or safe Filter On Variables | high | Templates | `rules_templates.go` | [ ] |
| BP-PY-35 | SQLAlchemy text With F-String | high | Database | `rules_db.go` | [ ] |
| BP-PY-36 | SQLAlchemy Session Not Closed | medium | Database | `rules_db.go` | [ ] |
| BP-PY-37 | DB-API Cursor Execute With Percent Format | high | Database | `rules_db.go` | [ ] |

**Coverage check:** assigned set fully planned — zero missing unplanned.

---

## Dependencies

| Depends on | Note |
|------------|------|
| Package scaffold | `RegisterRule`, scan loop, `pushAt`, facts, string/call helpers |
| `common.go` helpers | `isStringLiteral`, `firstCallArg` / `callArgsRegion`, `codeLinesFacts`, `isPythonTestFile` |
| Catalogue | `ruleset/python/bad-practices.json` keys `BP-PY-33`…`37` |
| Parent README | `plans/v0.0.2/heuristics/bp-plans/README.md` |
| Sibling | Django raw SQL `BP-PY-24` in `batch-03-django.md` — do not reimplement here |
| Out of scope | FastAPI (`batch-04`); non-Python template file types |

---

## References

- Catalogue: `ruleset/python/bad-practices.json`
- Related open-without-with pattern: `BP-PY-7` in `rules_core.go` (session close analogy)
- Tests pattern: `internal/lang/python/detectors/bad_practices/scan_test.go`
- Inventory: `plans/v0.0.2/heuristics/bp-plans/_inventory.json`
- Parent ledger: `plans/v0.0.2/heuristics/python-heuristics-bp.md`
- Skill: `plans/skills/phase-wise-checklist/SKILLS.md`
