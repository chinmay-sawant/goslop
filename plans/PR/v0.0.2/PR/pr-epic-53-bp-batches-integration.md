## Summary

Integration of all remaining BP-PY heuristic batches (01–08) so the full catalogue of **50** `BP-PY-*` rules is registered under `internal/lang/python/detectors/bad_practices/`. Prefer this PR over individual batch PRs.

## Child PRs (superseded)

| Batch | PR | Branch | Rules |
|------:|----|--------|-------|
| 01–02 | #60 | `feat/bp-batch-01-02` | 3,5,14,15,18–20 |
| 03 | #62 | `feat/bp-batch-03-django` | 22–28 |
| 04–05 | #64 | `feat/bp-batch-04-05` | 29–37 |
| 06–07 | #63 | `feat/bp-batch-06-07` | 38–47 |
| 08 | #61 | `feat/bp-batch-08-prod` | 48–50 |

Prior shipped (PR #58): 1,2,4,6,7,8–13,16,17,21

## Architecture

- Package: `internal/lang/python/detectors/bad_practices/`
- Domain-split rule files (`rules_django.go`, `rules_fastapi.go`, …)
- **Soft max 1500 / hard max 2000 lines** per Go file (largest ~658)
- Pure-Go source patterns; `RegisterRule`; `BP-PY-*` IDs only

## Test plan

- [x] `make lint`
- [x] `make test`
- [x] All 50 IDs in `RuleIDs()`

## Related issues

- Closes #53
- Relates to #51
- Relates to #54 (PERF still deferred)

## Follow-ups

- CWE expansion beyond priority 5
- PERF catalogue + heuristics (#54)
