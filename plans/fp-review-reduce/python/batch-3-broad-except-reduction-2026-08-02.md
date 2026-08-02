# Batch 3 — broad-except family FP reduction (2026-08-02)

## Scope

MASTER.md root cause #3: handlers that re-raise / record / log / `set_exception`,
but detectors treated them as swallowing failures.

Rules in scope: **BP-PY-1**, **CWE-396**, **BP-PY-42** (primary).  
**BP-PY-2 / CWE-390 / CWE-1071** left unchanged (see below).

Checklist sources: `niquests`, `logxide`, `calgebra`, `Project_Parva`,
`onlymaps`, `html2pic`, `FuncToWeb`, `CourtScrapper`, `requestSpeedTest`.

## Missing precondition (stated before edit)

Rule conditions require the handler to *hide* / *swallow* failures. Detectors
matched the `except Exception` / bare-try shape without inspecting the suite for:

| Surface action | Example FP |
| --- | --- |
| `raise` / `raise … from e` | html2pic CWE-396 `#6,#7,#10,#32`; Project_Parva `#37` |
| `exc_info=True` / `logger.exception(` | calgebra BP-PY-1 `#13`; CourtScrapper CWE-396 `#280`; logxide `#230` |
| `future.set_exception(e)` | niquests BP-PY-1 `#134+` |
| `result.error =` / `_error_result(e)` | logxide `#4`; calgebra `#12,#14` |
| BP-PY-42 whole-file fallback on `*_test.py` without `def test_*` | requestSpeedTest `#1,#2` |

Bare / `except: pass` / warn-and-return-default handlers must remain TPs
(WeThePeople / httpmorph / html2pic BP-PY-1 TPs).

## Detector changes

### BP-PY-1 — `internal/lang/python/detectors/bad_practices/rules_core.go`

Expanded `suiteReraises` into `suiteSurfacesFailure` / `suiteLineSurfacesFailure`.
Skip broad `except Exception|BaseException` when the suite contains any of:

- `raise` / `raise …`
- `exc_info`
- `.exception(`
- `set_exception(`
- `_error_result(`
- `.error =` (result/error field recording)

Bare `except:` still always fires.

### CWE-396 — `internal/lang/python/detectors/cwe/rules_platform.go`

Previously: first generic except in the file, no suite inspection.  
Now: walk masked lines, skip suites that `suiteSurfacesFailureMasked`, emit the
first remaining unhandled generic except (preserves single-finding-per-file).

### BP-PY-42 — `internal/lang/python/detectors/bad_practices/rules_testing.go`

1. Removed whole-file fallback for `*_test.py` modules with no `def test_*`
   (fixes helper/benchmark modules that only look test-named).
2. Skip expect-failure excepts whose suite surfaces the failure (re-raise etc.).

### BP-PY-2 / CWE-390 / CWE-1071 — no detector change

Audited “FPs” for these IDs are almost entirely `except …: pass` bodies
(optional import, teardown, best-effort probe). That is a different judgment
from root cause #3 suite-handling; narrowing them would suppress real TPs.
They remain structural except-pass detectors.

## Fixtures (unique `batch3` names)

### BP-PY-1

| Case | Safe pattern | Vulnerable pattern |
| --- | --- | --- |
| `BP-PY-1-batch3-exc-info` | `exc_info=True` warning | warning without traceback |
| `BP-PY-1-batch3-error-result` | `return _error_result(e)` | `return []` |
| `BP-PY-1-batch3-set-exception` | `future.set_exception(e)` | `pass` |
| `BP-PY-1-batch3-result-error` | `result.error = …` | `result.ok = False` only |

### BP-PY-42

| Case | Safe pattern | Vulnerable pattern |
| --- | --- | --- |
| `BP-PY-42-batch3-reraise` | `def test_*` + `raise e` | `def test_*` + `pass` |
| `BP-PY-42-batch3-non-test-helper` | `httpx_test.py` helper, no `test_*` | real `test_*` + `pass` |

### CWE-396

| Case | Safe pattern | Vulnerable pattern |
| --- | --- | --- |
| `CWE-396-batch3-reraise-from` | `raise RenderError(...) from e` | warn + `return None` |
| `CWE-396-batch3-exc-info` | `logger.error(..., exc_info=True)` | silent default return |
| `CWE-396-batch3-set-exception` | `set_exception(e)` | `set_result(None)` |
| `CWE-396-batch3-result-error` | `result.error = …` | `result.ok = False` only |

## Tests updated

- `internal/lang/python/detectors/bad_practices/audit_variants_test.go` — batch3 BP cases
- `internal/lang/python/detectors/cwe/audit_variants_test.go` — batch3 CWE-396 cases
- `tests/integration/python/bp_matrix_test.go` — inventory floor updated for parallel fixture growth (`>= 178`)

## Validation results

Focused command:

```text
go test ./internal/lang/python/detectors/cwe/... \
        ./internal/lang/python/detectors/bad_practices/... \
        ./tests/integration/python -count=1
```

| Surface | Result |
| --- | --- |
| `cwe` package | PASS |
| `bad_practices` package | PASS |
| `tests/integration/python` | PASS |
| Batch3 BP-PY-1 / BP-PY-42 audit variants | PASS (all 6) |
| Batch3 CWE-396 audit variants | PASS (all 4) |

Incidental during parallel work: RE2 compile fix in `cwe/fp_guards.go`
(`\1` → `\w+`) so the package could load while batch 5 was mid-flight. Not a
broad-except change.

## Sample finding IDs addressed by guardrail

| Repo | IDs | Rule | Guardrail hit |
| --- | --- | --- | --- |
| html2pic | 6, 7, 10, 32 | CWE-396 | `raise … from e` |
| requestSpeedTest | 1, 2 | BP-PY-42 | no `def test_*` / re-raise |
| requestSpeedTest | 3 | BP-PY-1 | (partial — error counters; not all accounting forms) |
| calgebra | 12, 14–16 | BP-PY-1 | `_error_result` |
| calgebra | 13 | BP-PY-1 | `exc_info=True` |
| CourtScrapper | 280 | CWE-396 | `exc_info=True` |
| logxide | 4, 5 | BP-PY-1 / CWE-396 | `result.error =` |
| logxide | 230–232 | BP-PY-1 / CWE-396 | `.exception(` / `exc_info` |
| niquests | 134+ | BP-PY-1 | `set_exception(` |
| Project_Parva | 37 | CWE-396 | `raise` (+ `.exception(`) |

## Remaining uncertainty

- Retry-then-re-raise via library helpers (`retries.increment(...)` in niquests)
  without an in-suite `raise` / `set_exception` is **not** covered.
- Outer-flow re-raise *after* the except suite (calgebra JSON-message builders
  `#9–11`) is **not** covered — raise is a sibling, not suite content.
- Pure `logger.warning("…", exc)` without `exc_info` (Project_Parva `#38`) is
  **not** exempted (would over-suppress warn-and-continue TPs).
- BP-PY-2 / CWE-390 / CWE-1071 best-effort/`ImportError: pass` audit FPs remain;
  fixing them needs a separate optional-dependency / teardown policy.

## Checklist artifacts

Audit report checkboxes were **not** modified (per batch instructions).
