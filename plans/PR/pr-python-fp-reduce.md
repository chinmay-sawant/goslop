## Summary

Tightens Python CWE and bad-practice heuristics so they stop flagging common false-positive patterns found in the 2026-07-31 corpus audit (49 findings → 22 after reduction). Adds structural nesting, ORM expression, test-module, and lifecycle guardrails plus paired safe/vulnerable fixtures for each refined rule.

---

## Motivation / context

- Plans: `plans/fp-review-reduce/python/false-positive-audit-2026-07-31.md`, `plans/fp-review-reduce/python/review-and-reduce-implementation-2026-07-31.md`
- Issues: see **Related issues**
- Prior work: FP audit/review-reduce skills landed in #70; this PR implements the detector guardrails from that audit

Baseline audit: 27 confirmed false positives, 20 true positives, 2 uncertain. Post-change full no-cache corpus scan: 22 findings (`49 - 27`).

---

## Changes

### CWE detector guardrails (`internal/lang/python/detectors/cwe/`)

| Rule | Guardrail |
|------|-----------|
| **CWE-1124** | Count structural control-flow blocks instead of raw indent / data-layout; do not count `except`/`finally` as extra nesting |
| **CWE-89** | Recognize local SQLAlchemy-style `select`/`delete`/`update`/`insert` statements (including parenthesized multiline assigns) |
| **CWE-93** | Do not treat `str(int(...))` / `str(round(...))` header values as externally controlled |
| **CWE-1046** | Require text-accumulator evidence before flagging `+=` inside a loop as string concat |
| **CWE-312 / 798 / 396** | Skip conventional Python test modules (`tests.py`, `test_*.py`, `/tests/`, etc.) |

Shared helper: `isPythonTestModule` in `common.go`. Fixture helpers: `assertCWEFixtureCase` for variant stems.

### Bad-practice detector guardrails (`internal/lang/python/detectors/bad_practices/`)

| Rule | Guardrail |
|------|-----------|
| **BP-PY-1 / BP-PY-42** | Allow broad-except / try-except patterns that collect thread errors for later assert in tests |
| **BP-PY-26** | Flag `csrf_exempt` only when the view body looks state-changing (drops “flag all” v0 behavior) |
| **BP-PY-38** | Recognize `create_task` / `ensure_future` retained in a list/collection |
| **BP-PY-46** | Suppress `print` inside CLI command decorator bodies |
| **Test file detection** | Treat Django-style `tests.py` / `conftest.py` as test modules |

### Fixtures and integration matrix

- New safe/vulnerable pairs under `tests/fixtures/python/{cwe,bp}/` for each audit guardrail variant (e.g. `CWE-89-orm-expression`, `BP-PY-26-read-only`, thread-collection cases)
- `PythonCWERuleID` maps fixture stems like `CWE-89-orm-expression` → `CWE-89`
- BP inventory expectation raised to 55 cases; dedicated `audit_variants_test.go` for both packs

### Plans / evidence

- Audit report and implementation evidence under `plans/fp-review-reduce/python/`
- This PR body: `plans/PR/pr-python-fp-reduce.md`

---

## Impact

| Area | Impact |
|------|--------|
| **Performance** | Negligible; slightly more line-walk work for CWE-1124/1046 on Python-enabled scans only |
| **Memory** | Negligible |
| **Behavior / correctness** | Fewer false positives on the audited patterns; true-positive corpus cases retained. BP-PY-26 no longer reports non-state-changing `csrf_exempt` views |
| **API / CLI** | None |
| **Dependencies** | None |
| **Binary size / build time** | Negligible pure-Go changes |

---

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| Consumers expecting BP-PY-26 on every `csrf_exempt` | Findings only when the view looks state-changing |
| Integration tests hard-coding BP case count 50 | Updated to 55 |

---

## Test plan

- [x] `go test ./internal/lang/python/detectors/cwe/... ./internal/lang/python/detectors/bad_practices/... -count=1`
- [x] `go test ./tests/integration/python -count=1`
- [x] `make lint-all`
- [x] `make test`
- [ ] `make run` wall time vs baseline — N/A for default Go path (Python opt-in)
- [ ] `make reference-metrics` — N/A (Go detector surface unchanged)
- [x] Full no-cache corpus scan on `codehound-python-perf-targets` → 22 findings

### Commands

```sh
go test ./internal/lang/python/detectors/cwe/... ./internal/lang/python/detectors/bad_practices/... -count=1
go test ./tests/integration/python -count=1
make lint-all
make test
```

### Recorded outcome (2026-07-31)

- Unit + audit-variant fixtures green
- Python CWE/BP integration matrices green
- `make lint-all` / `make test` green
- Corpus: 22 findings remaining (20 TP + 2 uncertain from audit)

---

## Screenshots / sample output

```
Full no-cache corpus scan after guardrails: 22 findings
(49 baseline − 27 confirmed false positives)
Uncertain retained: CWE-1124 app/tasks.py:38; CWE-924 app/routes.py:19
```

---

## Related issues

- Relates to #51
- Relates to #70 (FP audit / review-reduce skills)

---

## PR metadata checklist (author)

- [x] Self-assigned (`--assignee @me`)
- [x] Labels applied
- [x] Related issues filled with real ticket IDs
- [x] Filled body under `plans/PR/pr-python-fp-reduce.md`

---

## Follow-ups (out of scope)

- Resolve the two uncertain audit findings (CWE-1124 tasks nesting; CWE-924 host header)
- Further FP passes on new corpora after PERF heuristics land (#54)

---

## Reviewer checklist

- [ ] Behavior matches summary and test plan
- [ ] No unrelated changes in diff
- [ ] Public API / CLI changes documented
- [ ] New rules have fixture coverage when applicable
- [ ] PR has assignee and labels
- [ ] Related issues use correct Closes/Relates keywords
- [ ] No secrets or generated artifacts committed
