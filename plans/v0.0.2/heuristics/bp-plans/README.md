# v0.0.2 / #53 — BP-PY batch plans (execution index)

> **Parent:** `plans/v0.0.2/heuristics/python-heuristics-bp.md` — canonical #53 ledger  
> **Epic:** [#51](https://github.com/chinmay-sawant/goslop/issues/51)  
> **Issue:** [#53](https://github.com/chinmay-sawant/goslop/issues/53) expansion — batchwise PRs for remaining `BP-PY-*`  
> **Status:** all 50 BP-PY RegisterRule-implemented on integration branch (2026-07-31); batch ledgers evidence  
> **Inventory snapshot:** `plans/v0.0.2/heuristics/bp-plans/_inventory.json`

---

## Overview

Batch plans under this directory are **live execution ledgers** for #53 follow-up work. Each batch file follows `plans/skills/phase-wise-checklist/SKILLS.md`:

- Atomic `[ ]` / `[x]` / `[~]` rows with path, rule ID, expected behavior, proof.
- Mark `[x]` only after detector + hit/miss proof and (for code) `make lint` + `make test`.
- One **batchwise PR** per batch file (do not merge unrelated batches in one PR unless tiny).

### Architecture (all batches)

| Constraint | Detail |
|------------|--------|
| Package | `internal/lang/python/detectors/bad_practices/` |
| Registration | `RegisterRule(id, fn)` from `init()` in domain rule files |
| Detector surface | single `PythonBadPracticeScan` (`scan.go`); do not add a second `core.Detector` |
| Detection v0 | pure-Go **source-pattern** heuristics over `ParsedUnit.Source` (+ `bpFacts` needles) |
| Metadata | `MetadataForID` / catalogue `ruleset/python/bad-practices.json` — IDs stay `BP-PY-*` |
| File size | **1500 soft / 2000 hard max lines per Go file** — split domain files before growing existing ones |
| Split targets | `rules_async.go`, `rules_testing.go`, `rules_deps.go`, `rules_observability.go`, `rules_prod.go` (plus existing `rules_core.go`, `rules_security.go`, `rules_framework.go`) |
| Tests | extend `scan_test.go` (or domain `*_test.go` if `scan_test.go` approaches line cap) |
| Product default | Go-only registry; Python opt-in via `languages = ["python"]` |

**Do not** reopen shipped IDs as pending work — see [batch-00-shipped.md](./batch-00-shipped.md).

---

## Batch index

| Batch | File | IDs | Theme | Target files |
|------:|------|-----|-------|--------------|
| 00 | [batch-00-shipped.md](./batch-00-shipped.md) | 1,2,4,6,7,8–13,16,17,21 | **Ledger only** — already `RegisterRule` | `rules_core.go`, `rules_security.go`, `rules_framework.go` |
| 01 | [batch-01-core-prod.md](./batch-01-core-prod.md) | 3, 5, **14**, 15 | core leftovers + HTTP timeout / httpx close | `rules_core.go`, **`rules_prod.go`** |
| 02 | [batch-02-flask.md](./batch-02-flask.md) | 18, 19, 20 | Flask remainder | `rules_framework.go` or flask split |
| 03 | [batch-03-django.md](./batch-03-django.md) | **22–28** (21 already shipped) | Django | **`rules_django.go`** |
| 04 | [batch-04-fastapi.md](./batch-04-fastapi.md) | **29–32** | FastAPI / Starlette | **`rules_fastapi.go`** |
| 05 | [batch-05-templates-db.md](./batch-05-templates-db.md) | **33–37** | Jinja2 + SQLAlchemy/DB-API | **`rules_templates.go`**, **`rules_db.go`** |
| 06 | [batch-06-async.md](./batch-06-async.md) | **38, 39, 40** | Async / threading | **`rules_async.go`** |
| 07 | [batch-07-testing-deps-obs.md](./batch-07-testing-deps-obs.md) | **41–47** | Testing, deps, observability | **`rules_testing.go`**, **`rules_deps.go`**, **`rules_observability.go`** |
| 08 | [batch-08-prod-hardening.md](./batch-08-prod-hardening.md) | **48, 49, 50** | Production hardening (**not** 14) | **`rules_prod.go`** (append; may share file with batch-01) |

### ID ownership (no double-booking)

| ID | Owner batch |
|----|-------------|
| BP-PY-1,2,4,6,7,8–13,16,17,21 | **00 shipped** — closed |
| BP-PY-14 | **batch-01** (or sibling core/HTTP batch) — **not** batch-08 |
| BP-PY-38,39,40 | **06** |
| BP-PY-41–47 | **07** |
| BP-PY-48,49,50 | **08** |
| Other missing IDs | other sibling batch plans (03–05, etc.) |

---

## PR policy

1. One PR per batch ledger (title: `python(bp): batch-NN …` / `Relates to #53` / `Relates to #51`).
2. Before PR: `gofmt`, `make lint`, `make test`; record outcomes in the batch checklist Phase N validation.
3. Update `_inventory.json` `implemented` / `missing` when a batch lands (optional automation; keep honest).
4. Do **not** `Closes #53` until catalogue coverage + deferred inventory is honest in parent `python-heuristics-bp.md`.

---

## Current package size (pre-batch baseline)

From `_inventory.json` / `wc` at plan time (~2026-07-31):

| File | ~lines | Note |
|------|-------:|------|
| `common.go` | 446 | shared helpers — prefer not to dump rules here |
| `rules_core.go` | 307 | shipped core |
| `rules_security.go` | 283 | shipped security |
| `rules_framework.go` | 200 | shipped Flask/Django subset |
| `scan_test.go` | 296 | grow carefully; split tests if >1500 |
| **package total** | ~1968 | new domains **must** be new files |

---

## Dependencies

| Depends on | Note |
|------------|------|
| Parent ledger | `plans/v0.0.2/heuristics/python-heuristics-bp.md` |
| Scaffold | Phase 1 package already on tree (`register.go`, `scan.go`, `metadata.go`, `facts.go`, `common.go`) |
| Catalogue | `ruleset/python/bad-practices.json` (50 keys) |
| Skill | `plans/skills/phase-wise-checklist/SKILLS.md` |

---

## References

- Parent epic rollup: `plans/v0.0.2/heuristics/python-heuristics.md`
- Catalogue README: `ruleset/python/README.md`
- Go pattern reference: `internal/lang/go/detectors/bad_practices/`
