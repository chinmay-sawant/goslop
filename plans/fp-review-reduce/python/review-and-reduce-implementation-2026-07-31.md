# Python false-positive reduction — implementation evidence

- Date: 2026-07-31
- Audit source: `./plans/fp-review-reduce/python/false-positive-audit-2026-07-31.md`
- Target corpus: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets`
- Baseline: 49 findings; 27 confirmed false positives, 20 true positives, 2 uncertain findings.

## Implemented guardrails

| Audit pattern | Rules refined | Guardrail |
| --- | --- | --- |
| Indentation/nesting | CWE-1124 | Count structural blocks rather than collection layout; do not count exception-handler clauses as an additional nesting level. |
| ORM expression | CWE-89 | Recognize local SQLAlchemy `select`/`delete`/`update`/`insert` statements, including parenthesized multiline assignments. |
| Internal numeric header | CWE-93 | Do not treat `str(int(...))` or `str(round(...))` as externally controlled header content. |
| Test-only behavior | CWE-312, CWE-798, CWE-396, BP-PY-1, BP-PY-42 | Recognize conventional Python test modules and collected thread-error assertions. |
| Counter accumulation | CWE-1046 | Require text-accumulator evidence before treating `+=` in a loop as string concatenation. |
| Read-only endpoint | BP-PY-26 | Require a state-changing request or write operation before reporting a CSRF-exempt view. |
| Managed task lifecycle | BP-PY-38 | Recognize `create_task` calls retained in a collection. |
| Intentional command output | BP-PY-46 | Recognize functions registered by a CLI command decorator. |

## Regression coverage

Each guardrail has a safe and vulnerable text fixture under `./tests/fixtures/python/{cwe,bp}/`. Dedicated per-heuristic tests and the Python CWE/BP integration matrices consume those fixtures; no Python snippets are embedded in the Go tests.

## Validation

- `go test ./internal/lang/python/detectors/cwe/... ./internal/lang/python/detectors/bad_practices/... -count=1` — passed.
- `go test ./tests/integration/python -count=1` — passed.
- `make lint-all` — passed.
- `make test` — passed.
- Full no-cache corpus scan — 22 findings, exactly `49 - 27` confirmed false positives.

The remaining scan output preserves the two audit entries intentionally left uncertain: CWE-1124 at `app/tasks.py:38` and CWE-924 at `app/routes.py:19`.
