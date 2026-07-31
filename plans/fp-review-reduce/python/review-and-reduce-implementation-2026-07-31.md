# Python false-positive reduction — implementation evidence

- Date: 2026-07-31
- Audit source: `./plans/fp-review-reduce/python/false-positive-audit-2026-07-31.md`
- Target corpus: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets`
- Baseline: 49 findings; 29 confirmed false positives, 20 true positives, and no remaining uncertain findings after the CWE-1124 and CWE-924 scope clarifications.

## Implemented guardrails

| Audit pattern | Rules refined | Guardrail |
| --- | --- | --- |
| Indentation/nesting | CWE-1124 | Count executable control-flow headers rather than collection layout or class/function declarations; do not count exception-handler clauses as an additional nesting level. |
| Authenticated webhook route | CWE-924 | Respect a module-level header-authentication gate; do not require an HMAC signature when the route’s message-integrity boundary is established outside the handler. |
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
- Full no-cache corpus scan — 20 findings, exactly `49 - 29` confirmed false positives.

All 29 confirmed false positives are now suppressed; the remaining 20 findings satisfy their respective rule conditions.

## Follow-up reduction: benchmark and SQL false positives (2026-08-01)

- Audit source: `./plans/fp-review-reduce/python/false-positive-audit-2026-07-31.md`
- Baseline: 56 findings from the no-cache `make run-python` scan, including 16 confirmed false positives and no uncertain findings.
- Validation timestamp: `2026-08-01T00:10:39+05:30`.

### Completion record

| Findings | Rule(s) | Detector condition and evidence | Safe fixture | Vulnerable fixture | Regression coverage |
| --- | --- | --- | --- | --- | --- |
| 1, 45 | CWE-89 | Suppress the exact four-argument `__call__` DB-API execute passthrough unless the wrapper rewrites `sql`; recognize direct SQLAlchemy Core query-builder expressions as bound ORM/query-builder use. | `tests/fixtures/python/cwe/CWE-89-execute-wrapper-safe.txt`; `tests/fixtures/python/cwe/CWE-89-orm-direct-safe.txt` | Matching `...-vulnerable.txt` fixtures | `internal/lang/python/detectors/cwe/audit_variants_test.go`; `tests/integration/python` |
| 2, 3, 41–44, 46, 50–53, 56 | BP-PY-46 | Suppress intentional output in files whose path contains the exact `bench` or `benchmarks` component; ordinary application/library output remains reportable. | `tests/fixtures/python/bp/BP-PY-46-benchmark-script-safe.txt` | `tests/fixtures/python/bp/BP-PY-46-benchmark-script-vulnerable.txt` | `internal/lang/python/detectors/bad_practices/audit_variants_test.go`; `tests/integration/python/bp_matrix_test.go` |
| 54, 55 | CWE-312, CWE-798 | Suppress synthetic benchmark secrets only when the file is under a benchmark path; deployed application configuration remains reportable. | `tests/fixtures/python/cwe/CWE-312-benchmark-secret-safe.txt`; `tests/fixtures/python/cwe/CWE-798-benchmark-secret-safe.txt` | Matching `...-vulnerable.txt` fixtures | `internal/lang/python/detectors/cwe/audit_variants_test.go`; `tests/integration/python` |

### Validation

- Focused detector and Python integration tests: `go test ./internal/lang/python/detectors/cwe/... ./internal/lang/python/detectors/bad_practices/... ./tests/integration/python -count=1` — passed.
- Full lint: `make lint-all` — passed.
- Full tests: `make test` — passed.
- Diff hygiene: `git diff --check` — passed.
- Exact no-cache rescan: `make run-python` — 40 findings from 58 files and 2,716 lines; 7 high, 0 info, 0 low, 33 medium; 40 context files and 2 chunk files exported. This removes exactly the 16 audited false-positive findings from the 56-finding baseline.
- Retained `CWE-312` and `CWE-798` findings are outside benchmark paths in Django settings and remain true positives. No audited false-positive `CWE-89` or `BP-PY-46` context remains in the current chunks.
- Remaining uncertainty: none.
