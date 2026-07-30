## Summary

Implements **BP-PY batch 03 (Django)** heuristics: `BP-PY-22`…`BP-PY-28` as pure-Go source-pattern detectors under `internal/lang/python/detectors/bad_practices/rules_django.go`. `BP-PY-21` remains shipped prior evidence (not re-implemented). Metadata is registered via `init()` into `metaByID`; unit hit/miss tests live in `rules_django_test.go`.

---

## Motivation / context

- Plan: `plans/v0.0.2/heuristics/bp-plans/batch-03-django.md`
- Catalogue: `ruleset/python/bad-practices.json`
- Completes the Django family after `BP-PY-21` (DEBUG) from the Phase 1 BP scaffold

---

## Changes

### Detectors (`rules_django.go`)

| ID | Heuristic | Sev |
|----|-----------|-----|
| BP-PY-22 | `SECRET_KEY = '...'` in Django settings; skip env/`getenv`/`config(`/`env(` | high |
| BP-PY-23 | `ALLOWED_HOSTS` contains `*` always; empty `[]`/`()` when `DEBUG = False` nearby | medium |
| BP-PY-24 | `Model.objects.raw` / `.execute` first arg is f-string / `.format` / `%` format | high |
| BP-PY-25 | `mark_safe` / `SafeString` on non-literal first arg | high |
| BP-PY-26 | `@csrf_exempt` on views (state-changing signal preferred; any def still flagged) | high |
| BP-PY-27 | `Model(**request.POST)` / `objects.create(**request.data)` mass assignment | medium |
| BP-PY-28 | for-loop over queryset with multi-hop attr access and no `select_related`/`prefetch_related` (**review-only heuristic**) | medium |

### Wire-up

- `init()` writes catalogue-aligned metadata into `metaByID` then `RegisterRule` for `22`…`28`
- Does **not** grow `rules_framework.go` with new Django rules
- Does **not** re-implement `BP-PY-21`
- Coordinates with Flask `BP-PY-17` (Django settings path ownership for bare `SECRET_KEY`)

### Tests

- `rules_django_test.go` — hit/miss per rule + registration + severity spot-checks
- `scan_test.go` — `TestBPRulesRegistered` extended with `BP-PY-22`…`28`

### Plans / inventory

- `_inventory.json` implemented/missing lists updated
- `python-heuristics-bp.md` Django rows marked shipped

---

## Impact

| Area | Impact |
|------|--------|
| **Performance** | Light pure-Go line/regex scans on Python units when BP pack is enabled |
| **Memory** | Negligible (shared fact bag + line table) |
| **Behavior** | New `BP-PY-22`…`28` findings when `languages` includes Python |
| **API / CLI** | `--list-rules` surfaces the new IDs |
| **Dependencies** | None |

---

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| None | - |

---

## Test plan

- [x] `gofmt` on touched Go files
- [x] `go test ./internal/lang/python/detectors/bad_practices/ -count=1`
- [x] `make lint`
- [x] `make test`
- [x] `wc -l` each new/changed package file ≤2000 (prefer ≤1500)

### File-size snapshot (post-change)

| File | Lines |
|------|------:|
| `rules_django.go` | ~658 |
| `rules_django_test.go` | ~147 |
| `scan_test.go` | ~297 |
| `rules_framework.go` | ~200 (unchanged rules) |

---

## Related issues

- Relates to #53 (Python BP heuristics)
- Relates to #51 (Python heuristics epic)

Does **not** close #53 (catalogue remains incomplete outside this batch).
