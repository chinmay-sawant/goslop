## Summary

Implements **BP-PY batch 01 + 02** remainder heuristics under `internal/lang/python/detectors/bad_practices/`:

| Batch | IDs | Domain |
|-------|-----|--------|
| 01 | `BP-PY-3`, `BP-PY-5`, `BP-PY-14`, `BP-PY-15` | Core language + production/resource (requests/httpx) |
| 02 | `BP-PY-18`, `BP-PY-19`, `BP-PY-20` | Flask remainder (routes, errorhandler leak, send_file) |

Pure-Go source-pattern heuristics via existing `PythonBadPracticeScan` + `RegisterRule` in `init()`. Metadata for new IDs is set in each rules file’s `init()` (no `metadata.go` edits) to reduce merge conflicts.

---

## Motivation / context

- Plans: `plans/v0.0.2/heuristics/bp-plans/batch-01-core-prod.md`, `batch-02-flask.md`
- Catalogue: `ruleset/python/bad-practices.json` detection_notes
- Issues: **Relates to #53**, **Relates to #51** (do not Closes #53 yet — remaining batches still open)

---

## Changes

### New files (package `badpractices` only)

| File | Rules |
|------|--------|
| `rules_core_extra.go` | BP-PY-3 Raise Generic Exception; BP-PY-5 Wildcard Import |
| `rules_http.go` | BP-PY-14 requests Without Timeout; BP-PY-15 httpx Async Client Not Closed |
| `rules_flask.go` | BP-PY-18 Route Missing Methods; BP-PY-19 jsonify Error Leak; BP-PY-20 send_file User Path |
| `rules_batch0102_test.go` | Hit/miss + registration + severity tests for 01+02 |

### Small shared updates

- `facts.go` — additive needles for batch 01/02 fast-paths only

### Intentionally untouched

- `rules_framework.go` (Flask 16/17 + Django 21 stay put)
- `metadata.go` (new metas via per-file `init()`)
- `scan_test.go` (new coverage in `rules_batch0102_test.go`)
- Plugin / `all.go` (still single BP scan)

### Heuristics (short)

| ID | Detection |
|----|-----------|
| BP-PY-3 | `raise Exception` / `raise BaseException` outside tests |
| BP-PY-5 | `from … import *` in non-`__init__.py` modules |
| BP-PY-14 | `requests.{get,post,…}` / `session.*` call args missing `timeout=` |
| BP-PY-15 | `name = httpx.AsyncClient(...)` without `async with` / `await name.aclose()` |
| BP-PY-18 | `@*.route` without `methods=` whose body uses `request.form` / `get_json` / `json` / `files` |
| BP-PY-19 | `@*.errorhandler` body returns `str(e)` / `repr(e)` / `traceback.format_exc()` |
| BP-PY-20 | `send_file` / `send_from_directory` path from `request.*` without `safe_join` |

---

## Impact

| Area | Impact |
|------|--------|
| **Performance** | Light line/regex scans; needle prefilters |
| **Memory** | Negligible |
| **Behavior** | New BP-PY findings when Python + bad-practices enabled |
| **API / CLI** | No CLI change; more IDs under `--list-rules` for Python |
| **Dependencies** | None |

---

## Line counts (post-change)

```
  446 common.go
   69 facts.go
  225 metadata.go
   56 register.go
  161 rules_batch0102_test.go
  307 rules_core.go
   84 rules_core_extra.go
  352 rules_flask.go
  200 rules_framework.go
  190 rules_http.go
  283 rules_security.go
   94 scan.go
  296 scan_test.go
 2763 total
```

Hard max 2000 / soft max 1500: **all files under budget**.

---

## Validation

```sh
make lint   # green
make test   # green
go test ./internal/lang/python/detectors/bad_practices/ -count=1  # green
wc -l internal/lang/python/detectors/bad_practices/*.go | awk '$1>2000{print; exit 1}'  # ok
```

---

## Related issues

- Relates to #53
- Relates to #51

---

## Test plan

- [x] Unit: BP-PY-3/5/14/15/18/19/20 hit + miss cases in `rules_batch0102_test.go`
- [x] Registration + catalogue severity alignment for new IDs
- [x] Existing Flask 16/17 / Django 21 / core / security tests still pass
- [x] `make lint` + `make test`
