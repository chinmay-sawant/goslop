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

## Follow-up reduction: CWE-90 regex and CWE-88 fixed-argv (2026-08-02)

- Audit source: `./plans/fp-review-reduce/python/false-positive-audit-2026-08-02.md`
- Scope: Assignment A (CWE-90 IDs 23, 31, 41, 51, 76, 79, 82–84, 86, 90, 92) and Assignment B (CWE-88 IDs 53, 87, 91; retain TP finding 7).
- Do not change uncertain finding 40; audit checklist boxes left unchanged.

### Completion record

| Findings | Rule | Detector condition changed | Safe fixture | Vulnerable fixture | Tests |
| --- | --- | --- | --- | --- | --- |
| 23, 31, 41, 51, 76, 79, 82, 83, 84, 86, 90, 92 | CWE-90 | Drop bare `.search(` gate; require LDAP library evidence (`ldap3` / `import ldap` / `.search_s(`); skip `re` / `*_RE` receivers; pick LDAP filter arg for `search`/`search_s`. | `tests/fixtures/python/cwe/CWE-90-regex-search-safe.txt` | `tests/fixtures/python/cwe/CWE-90-regex-search-vulnerable.txt` | `audit_variants_test.go` (`CWE-90-regex-search`); `rules_injection_test.go` pair; `tests/integration/python` |
| 53, 87, 91 | CWE-88 | Skip conventional Python test modules; treat pure string/bytes literals and static literal concatenations as non-dynamic argv segments. Library `run_verapdf_check`-style dynamic `flavour`/`pdf` still reports. | `tests/fixtures/python/cwe/CWE-88-fixed-argv-safe.txt` | `tests/fixtures/python/cwe/CWE-88-fixed-argv-vulnerable.txt` | `audit_variants_test.go` (`CWE-88-fixed-argv`); `rules_injection_test.go` pair; `tests/integration/python` |

### Validation

- Focused: `go test ./internal/lang/python/detectors/cwe/... ./tests/integration/python -count=1` — passed.
- Diff hygiene: `git diff --check` — passed.
- Full `make lint-all` / `make test` / corpus rescan: not run in this assignment (focused coverage only).
- Remaining uncertainty: finding 40 (CWE-328) left untouched per assignment.

## Follow-up reduction: path + secrets FPs (2026-08-02)

- Audit source: `./plans/fp-review-reduce/python/false-positive-audit-2026-08-02.md`
- Scope: CWE-22 IDs 19, 22, 29, 48, 94, 99; CWE-73 ID 44; CWE-256/798 IDs 42, 43.
- Do not change uncertain finding 40; audit checklist boxes left unchanged; no commit.

### Completion record

| Findings | Rule | Detector condition changed | Safe fixture | Vulnerable fixture | Tests |
| --- | --- | --- | --- | --- | --- |
| 19, 22, 29, 48, 94, 99 | CWE-22 | Require confined-root join of a dynamic segment (`os.path.join` / same-line `Path(...) / dynamic`); skip `Path(__file__)`, intentional CLI `Path(argv)`, and `Path(path)+open` without join-escape; ignore comment/arithmetic ` / `. | `tests/fixtures/python/cwe/CWE-22-no-restricted-join-safe.txt` | `tests/fixtures/python/cwe/CWE-22-no-restricted-join-vulnerable.txt` | `audit_variants_test.go`; `rules_test.go` pair; `tests/integration/python` |
| 44 | CWE-73 | Suppress `sys.argv`-only path sinks under `if __name__ == "__main__"` (CLI/fixture out-dir); keep request-controlled sinks. | `tests/fixtures/python/cwe/CWE-73-main-cli-outdir-safe.txt` | `tests/fixtures/python/cwe/CWE-73-main-cli-outdir-vulnerable.txt` | `audit_variants_test.go`; `rules_path_fs_test.go` pair; `tests/integration/python` |
| 42 | CWE-256 | Suppress `fixture-password` literals and `fixtures.py` builder functions (`*fixture*` / `generate_*` / `*encrypted*`); production password literals still report. | `tests/fixtures/python/cwe/CWE-256-fixture-password-safe.txt` | `tests/fixtures/python/cwe/CWE-256-fixture-password-vulnerable.txt` | `audit_variants_test.go`; `rules_secrets_test.go` pair; `tests/integration/python` |
| 43 | CWE-798 | Same fixture-credential guard as CWE-256 (extends existing test/bench suppressions). | `tests/fixtures/python/cwe/CWE-798-fixture-password-safe.txt` | `tests/fixtures/python/cwe/CWE-798-fixture-password-vulnerable.txt` | `audit_variants_test.go`; `rules_secrets_test.go` pair; `tests/integration/python` |

### Validation

- Focused: `go test ./internal/lang/python/detectors/cwe/ -count=1 -run 'TestCoreCWEFixturePairs|TestCWEPathFSFixturePairs|TestCWESecretsFixturePairs|TestCWEFalsePositiveAuditFixtureVariants'` — passed.
- Integration: `go test ./tests/integration/python -count=1 -run TestPythonCWEFixturesMatrix` — passed.
- Package: `go test ./internal/lang/python/detectors/cwe/ -count=1` — passed.
- Diff hygiene: `git diff --check` on touched detector/fixture paths — passed.
- Synthetic rescan (`templates/goslop-python.toml`, only CWE-22/73/256/798): FP samples silent; `Path(root)/user`, `request.args` open, and production `password`/`API_KEY` still emit.
- Full `make lint-all` / `make test` / corpus `make run-python`: not required for this assignment (focused coverage only).
- Remaining uncertainty: finding 40 (CWE-328) left untouched per assignment.

## Follow-up reduction: pythoncoreengine audit 2026-08-02 (consolidated)

- Audit source: `./plans/fp-review-reduce/python/false-positive-audit-2026-08-02.md`
- Target corpus: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/pythoncoreengine`
- Baseline: **102** findings (96 FP, 5 TP, 1 Uncertain)
- Validation timestamp: `2026-08-01T20:25:00Z` (local 2026-08-02)
- Delegated reducers: BP-PY-46, BP-PY-41, CWE-90/88, CWE-328, CWE-22/73/256/798, PERF-PY-27/CWE-1046

### Completion record

| Findings | Rule(s) | Detector condition changed | Safe fixture(s) | Vulnerable fixture(s) |
| --- | --- | --- | --- | --- |
| 2–6, 8–17, 20–21, 25–28, 32–35, 93, 96–98, 101–102 | BP-PY-46 | Skip argparse/`main()`/`print_*`/`cmd_*` CLI presentation; expand `bench*.py` stem matching. Keep library debug prints (TP 39). | `BP-PY-46-argparse-cli-safe.txt` | `BP-PY-46-argparse-cli-vulnerable.txt` |
| 54–75, 88–89 | BP-PY-41 | Credit `_assert*`/`assert*` callees and same-file / sibling `helpers.py` helpers that `self.assert*` or raise `AssertionError`. | `BP-PY-41-assert-helper-safe.txt` | `BP-PY-41-assert-helper-vulnerable.txt` |
| 23, 31, 41, 51, 76, 79, 82–84, 86, 90, 92 | CWE-90 | LDAP-only gates/receivers; skip `re` / `*_RE.search`. | `CWE-90-regex-search-safe.txt` | `CWE-90-regex-search-vulnerable.txt` |
| 53, 87, 91 | CWE-88 | Skip test modules; treat pure literal argv segments as non-dynamic. Keep TP 7. | `CWE-88-fixed-argv-safe.txt` | `CWE-88-fixed-argv-vulnerable.txt` |
| 18, 24, 30, 37, 38, 45, 47, 52, 77, 78, 80, 81, 85 | CWE-328 | Require `securityHashContext` (password/credential/token/auth/key-derivation). Fingerprint MD5 ignored. | `CWE-328-fingerprint-safe.txt` | `CWE-328-fingerprint-vulnerable.txt` |
| 19, 22, 29, 48, 94, 99 | CWE-22 | Require confined-root dynamic join; skip `Path(__file__)`, CLI argv roots, plain `Path(path)+open`. | `CWE-22-no-restricted-join-safe.txt` | `CWE-22-no-restricted-join-vulnerable.txt` |
| 44 | CWE-73 | Suppress `sys.argv` sinks under `__main__`. | `CWE-73-main-cli-outdir-safe.txt` | `CWE-73-main-cli-outdir-vulnerable.txt` |
| 42, 43 | CWE-256, CWE-798 | Suppress fixture-password / `fixtures.py` builders. | `CWE-256-fixture-password-safe.txt`, `CWE-798-fixture-password-safe.txt` | matching `*-vulnerable.txt` |
| 1, 95, 100 | PERF-PY-27 | Require invariant path (not loop var / callee param). | `PERF-PY-27-unique-path-safe.txt`, `PERF-PY-27-analyze-once-safe.txt` | matching `*-vulnerable.txt` |
| 36 | CWE-1046 | `bytearray()` accumulators are not immutable text concat. | `CWE-1046-bytearray-loop-safe.txt` | `CWE-1046-bytearray-loop-vulnerable.txt` |

### Retained findings (rescan)

| Rule | Source | Maps to audit |
| --- | --- | --- |
| CWE-88 | `compliance/verapdf_report.py` | TP 7 |
| BP-PY-46 | `engine/doc.py` (`ENGINE_DEBUG_BUFFERS`) | TP 39 |
| CWE-328 | `engine/encrypt.py:236` password MD5 | Uncertain 40 family (security-context hit retained; `_md5_iterate` helper silent) |
| PERF-PY-24 | `engine/layout.py` | TP 46 |
| BP-PY-1 | `engine/render.py` | TP 49 |
| CWE-396 | `engine/render.py` | TP 50 |

### Validation

- Focused detectors + Python integration: `go test ./internal/lang/python/detectors/cwe/... ./internal/lang/python/detectors/bad_practices/... ./internal/lang/python/detectors/perf/... ./tests/integration/python -count=1` — **passed**.
- Full tests: `make test` — **passed**.
- Diff hygiene: `git diff --check` — **passed**.
- Exact no-cache rescan: `make run-python` — **6 findings** from 63 files / 20566 lines (was 102). Removes all **96** audited false positives; retains the 5 true positives plus one security-context CWE-328.
- `make lint-all` — **fails** on pre-existing branch issues outside this reduction (unused helpers in `rules_auth.go` / `perf/facts.go`, revive `unexported-return` on `PyCweFacts`, etc.). Introduced G304/wastedassign/`close` issues from this work were fixed before rescan.
- Audit checklist `[ ]` boxes left unchanged (per skill).
- Remaining uncertainty: finding 40 dual-use MD5 — not specially suppressed; password MD5 at `encrypt.py:236` still reports.

## Follow-up reduction: full-corpus perf-targets audit (2026-08-02)

- Audit source: `./plans/fp-review-reduce/python/false-positive-audit-2026-08-02-perf-targets.md`
- Target corpus: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets` (all projects)
- Baseline: **45** findings (2 FP, 43 TP, 0 Uncertain)
- Validation timestamp: `2026-08-01T21:15:00Z` (local 2026-08-02)

### FP mechanisms

| Finding | Rule | Mechanism |
| --- | --- | --- |
| 38 | CWE-91 | `ET.fromstring(xml_data)` parses a whole document; rule condition is formatting into XML/XPath. Detector treated any dynamic parse arg as injection. |
| 39 | PERF-PY-26 | One-shot CLI `main()` calling `parse_report`; `computeInLoop`/`inLoop` treated later lines as in-loop after earlier `for` loops in a sibling function, so hot-path fired incorrectly. |

### Completion record

| Findings | Rule | Detector condition changed | Safe fixture | Vulnerable fixture | Tests |
| --- | --- | --- | --- | --- | --- |
| 38 | CWE-91 | For `*.fromstring` only, require constructed markup (`f-string` / `.format` / `%` / concat). Bare dynamic parse args skipped. XPath sinks unchanged. | `tests/fixtures/python/cwe/CWE-91-fromstring-parse-safe.txt` | `tests/fixtures/python/cwe/CWE-91-fromstring-parse-vulnerable.txt` | `audit_variants_test.go` (`CWE-91-fromstring-parse`); `rules_injection_test.go` pair; `tests/integration/python` |
| 39 | PERF-PY-26 | Function-scoped loop membership: `def`/`class` boundaries clear prior loops in `computeInLoop` / `inLoop`. Cross-function loop poison removed for all PERF rules using `lineInLoop`. | `tests/fixtures/python/perf/PERF-PY-26-cli-parse-safe.txt` | `tests/fixtures/python/perf/PERF-PY-26-cli-parse-vulnerable.txt` | `audit_variants_test.go` (`PERF-PY-26-cli-parse`); `facts_test.go` boundary case; existing `PERF-PY-26` pair |

### Validation

- Focused: `go test ./internal/lang/python/detectors/cwe/... ./internal/lang/python/detectors/perf/... ./tests/integration/python -count=1` — passed.
- Full tests: `go test ./... -count=1` — passed.
- Diff hygiene: `git diff --check` — passed.
- Exact no-cache rescan: `make run-python` — **43 findings** from 116 files / 23033 lines (was 45). Removes both audited FPs; CWE-91 and PERF-PY-26 absent from chunks.
- Retained TPs still present: CWE-88 (`verapdf_report.py`), BP-PY-46 (`doc.py`), CWE-328 (`encrypt.py:236`), PERF-PY-24 (`layout.py`), plus Django/FastAPI/Flask TPs.
- `make lint-all` — fails on pre-existing branch issues outside this reduction.
- Audit checklist `[ ]` boxes left unchanged (per skill).
- Remaining uncertainty: none.
