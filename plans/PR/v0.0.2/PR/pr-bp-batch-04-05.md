## Summary

Implements **BP-PY batches 04 + 05** (FastAPI/Starlette, Jinja2 templates, SQLAlchemy/DB-API) as pure-Go source heuristics under `internal/lang/python/detectors/bad_practices/`. Adds nine catalogue rules `BP-PY-29`…`BP-PY-37` in dedicated files (`rules_fastapi.go`, `rules_templates.go`, `rules_db.go`) with hit/miss unit tests. Relates to #53 / #51.

---

## Motivation / context

- Plans: `plans/v0.0.2/heuristics/bp-plans/batch-04-fastapi.md`, `batch-05-templates-db.md`
- Ledger: `plans/v0.0.2/heuristics/python-heuristics-bp.md` (#53)
- Catalogue: `ruleset/python/bad-practices.json` (metadata already present; detectors were missing)

---

## Changes

### Batch 04 — FastAPI / Starlette (`rules_fastapi.go`)

| ID | Heuristic |
|----|-----------|
| BP-PY-29 | `global`/`nonlocal` or module-level `{}`/`[]` mutation in routes/deps |
| BP-PY-30 | Blocking I/O (`time.sleep`, `requests.*`, `subprocess.*`, sync ORM) in `async def` routes |
| BP-PY-31 | Route without `response_model=` returning ORM-ish values |
| BP-PY-32 | `FileResponse` with non-literal / user-influenced path |

### Batch 05 — Templates (`rules_templates.go`)

| ID | Heuristic |
|----|-----------|
| BP-PY-33 | `jinja2.Environment(..., autoescape=False)` |
| BP-PY-34 | `Markup(non-literal)` and Python-side `\|safe` |

### Batch 05 — Database (`rules_db.go`)

| ID | Heuristic |
|----|-----------|
| BP-PY-35 | `sqlalchemy.text` with f-string / `.format` / `%` / concat |
| BP-PY-36 | `Session`/`SessionLocal()` without `with` or `.close()` (no full CFG) |
| BP-PY-37 | `.execute(` with f-string / `%` format SQL (skips `text(...)` args; overlap policy with 35) |

### Shared

- Metadata entries in `metadata.go` (`metaByID`) matching catalogue titles/severities
- Needles extended in `facts.go`
- Tests: `rules_fastapi_test.go`, `rules_templates_test.go`, `rules_db_test.go`; registration list updated in `scan_test.go`
- File-size gate: production files well under 2000 lines (largest new file ~589)

---

## Impact

| Area | Impact |
|------|--------|
| **Performance** | Light extra source scans on Python units when BP pack enabled |
| **Memory** | Negligible |
| **Behavior** | New `BP-PY-29`…`37` findings on matching Python sources |
| **API / CLI** | `--list-rules` with Python lists the new IDs |
| **Dependencies** | None |

---

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| None | - |

---

## Test plan

- [x] `make lint`
- [x] `make test`
- [x] `go test ./internal/lang/python/detectors/bad_practices/ -count=1`
- [x] Hit/miss unit tests for all 9 rules
- [x] No production Go file > 2000 lines

### Commands

```sh
make lint
make test
go test ./internal/lang/python/detectors/bad_practices/ -count=1
wc -l internal/lang/python/detectors/bad_practices/*.go
```

---

## Related issues

- Relates to #53 (Python BP heuristics)
- Relates to #51 (Python heuristics epic)

---

## Checklist

- [x] Title: `feat(python/bp): batch 04-05 FastAPI templates and database`
- [x] Self-assigned
- [x] Labels: `enhancement`, `documentation`
- [x] Base: `main`
- [x] Pure-Go source patterns; `RegisterRule`; `BP-PY-*` only
