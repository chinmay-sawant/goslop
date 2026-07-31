# goslop - Python false-positive reduction

## Summary

This change reduces confirmed false positives from the Python detector scan by tightening SQL, benchmark-output, and benchmark-secret classification. It adds source-backed safe/vulnerable fixture pairs and records the audit and rescan evidence in the Python review ledger.

## Motivation / context

- Plans: `plans/fp-review-reduce/python/false-positive-audit-2026-07-31.md`
- Implementation evidence: `plans/fp-review-reduce/python/review-and-reduce-implementation-2026-07-31.md`
- Baseline: 56 findings, including 16 confirmed false positives and no uncertain findings.
- Related issue: see **Related issues**.

## Changes

### Detector guardrails

- CWE-89 now recognizes the exact DB-API execute-wrapper passthrough shape and direct SQLAlchemy Core query-builder expressions while retaining detection for SQL rewritten by the wrapper and dynamically constructed SQL.
- BP-PY-46 ignores intentional console output in exact `bench`/`benchmarks` path components while continuing to report ordinary application/library output.
- CWE-312 and CWE-798 ignore synthetic benchmark secrets only under benchmark paths; deployed application settings remain reportable.

### Fixture and audit coverage

- Added five safe/vulnerable text-fixture pairs covering the two CWE-89 cases, BP-PY-46 benchmark output, and benchmark secrets for CWE-312 and CWE-798.
- Extended the Python detector variants and integration matrix to consume the fixtures.
- Updated the audit report and implementation ledger with per-finding evidence, fixture paths, and the final rescan result.

## Impact

| Area | Impact |
|------|--------|
| **Performance** | Small bounded path checks and a narrow callback-shape check; the exact scan remains within the existing sub-second runtime. |
| **Memory** | No material change. |
| **Behavior / correctness** | Removes the 16 audited false positives while retaining true positives for dynamic SQL, application output, and deployed credentials. |
| **API / CLI** | No changes. |
| **Dependencies** | No changes. |
| **Binary size / build time** | No material change. |

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| None | - |

## Test plan

- [x] `make test`
- [x] `make lint`
- [x] `make lint-all`
- [x] Focused Python detector and integration tests
- [x] `make run-python` exact no-cache rescan

### Commands

```sh
go test ./internal/lang/python/detectors/cwe/... ./internal/lang/python/detectors/bad_practices/... ./tests/integration/python -count=1
make lint
make lint-all
make test
make run-python
```

## Screenshots / sample output

```text
scanned 58 files (2716 lines) in 125.9ms
40 findings
  severity: 7 high, 0 info, 0 low, 33 medium
exported 40 context file(s) ...; exported 2 chunk file(s) ...
```

## Related issues

- Relates to #51 — Epic: Python heuristic detectors (CWE, BP, PERF)

## PR metadata checklist (author)

- [x] Self-assigned (`--assignee @me`)
- [x] Labels applied
- [x] Related issues filled with a real ticket ID
- [x] Filled body committed under `plans/PR/pr-python-fp-reduction.md`

## Follow-ups (out of scope)

- None.

## Reviewer checklist

- [ ] Behavior matches summary and test plan
- [ ] No unrelated changes in diff
- [ ] New rules have fixture coverage when applicable
- [ ] PR has assignee and labels
- [ ] Related issues use the correct `Relates to` keyword
- [ ] No secrets or generated artifacts committed
