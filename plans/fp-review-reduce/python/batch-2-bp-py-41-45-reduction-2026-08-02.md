# Batch 2 — BP-PY-41 / BP-PY-45 false-positive reduction

- Date: 2026-08-02
- Agent: Batch 2 of 5
- Checklist: `plans/fp-validations/reports/` (pictex, httptap, niquests, Project_Parva, safer, rendercv, whatsapp-wrapped) + MASTER root cause #2
- Scope: `detectBPPY41` (`rules_testing.go`), `detectBPPY45` (`rules_deps.go`); fixtures `BP-PY-41-*` / `BP-PY-45-*` only
- Audit checklist `[ ]` boxes left unchanged

## Mechanisms (from function contexts)

| Pattern | Repos / evidence | Why the old detector fired |
| --- | --- | --- |
| `check_func(file_regression, …)` | pictex (160) | Only credited `assert` / `pytest.raises` / `assert*` callees |
| `benchmark(...)` | httptap (25+) | pytest-benchmark fixture is the verification mechanism |
| `_test` / `_inner_test_*` helpers with bare `assert` | safer 28–30; niquests 209+ | Helper map only credited `self.assert*` / `raise AssertionError` |
| Dedented triple-quoted samples before asserts | whatsapp-wrapped 1–2; rendercv 64, 72–73 | Body scan treated column-0 string lines as function exit (`ind <= defIndent`) |
| Sphinx `docs/conf.py` `sys.path.insert` | niquests finding 1 | Treated build-time docs config as runtime library mutation |
| Module-top `__file__` / `not in sys.path` bootstrap | Project_Parva (~133) | Import-time script bootstrap flagged as runtime path hacks |

Skipped uncertain: nox `tests_impl(...)` session orchestration (niquests finding 2) — not a same-file assert-bearing helper.

## Detector changes

### BP-PY-41 (`rules_testing.go`)

1. **Assertion idioms:** credit `check_func`, `benchmark`, `pytest.fail`, and `*regression.check(` / `file_regression.check(` / `image_regression.check(` / `data_regression.check(`.
2. **Helper bodies:** `sameFileAssertionHelpers` also credits bare `assert` (covers safer `_test` and niquests `_inner_test_*`).
3. **Body scan:** `pyTripleQuoteCarry` tracks open triple-quoted strings across lines so string content cannot abort the scan before real asserts.

### BP-PY-45 (`rules_deps.go`)

1. Skip Sphinx `docs/**/conf.py`.
2. Skip mutations whose line contains `__file__` (classic in-tree bootstrap).
3. Skip module-level inserts immediately under `if … not in sys.path:` (guarded bootstrap); **keep** in-function runtime inserts and hard-coded vendor inserts (e.g. `./vendor`).

## Completion record

| Finding IDs (representative) | Rule | Detector condition changed | Safe fixture | Vulnerable fixture |
| --- | --- | --- | --- | --- |
| pictex 41+ | BP-PY-41 | Credit `check_func` / regression `.check(` | `tests/fixtures/python/bp/BP-PY-41-check-func-safe.txt` | `…/BP-PY-41-check-func-vulnerable.txt` |
| httptap 67–91 | BP-PY-41 | Credit `benchmark(...)` | `…/BP-PY-41-benchmark-safe.txt` | `…/BP-PY-41-benchmark-vulnerable.txt` |
| safer 28–30; niquests 209+ | BP-PY-41 | Bare `assert` in same-file helpers | `…/BP-PY-41-bare-assert-helper-safe.txt` | `…/BP-PY-41-bare-assert-helper-vulnerable.txt` |
| whatsapp-wrapped 1–2; rendercv 64,72–73 | BP-PY-41 | Triple-quote carry in body scan | `…/BP-PY-41-string-body-scan-safe.txt` | `…/BP-PY-41-string-body-scan-vulnerable.txt` |
| niquests 1 | BP-PY-45 | Skip `docs/conf.py` | `…/BP-PY-45-docs-conf-safe.txt` | `…/BP-PY-45-docs-conf-vulnerable.txt` |
| Project_Parva 76+ (`Path(__file__)`) | BP-PY-45 | Skip `__file__` bootstrap lines | `…/BP-PY-45-file-bootstrap-safe.txt` | `…/BP-PY-45-file-bootstrap-vulnerable.txt` |
| Project_Parva guarded `not in sys.path` | BP-PY-45 | Skip module-level guarded bootstrap; keep in-function inserts | `…/BP-PY-45-guarded-bootstrap-safe.txt` | `…/BP-PY-45-guarded-bootstrap-vulnerable.txt` |

## Tests

- Appended cases to `audit_variants_test.go` and `rules_deps_test.go` / `rules_testing_test.go`.
- `go test ./internal/lang/python/detectors/bad_practices/ -count=1 -run 'TestBPPY41|TestBPPY45|TestBPFalsePositiveAuditFixtureVariants/BP-PY-41|TestBPFalsePositiveAuditFixtureVariants/BP-PY-45'` — **PASS** (all new + existing BP-PY-41/45 cases).
- Full `go test ./internal/lang/python/detectors/bad_practices/... ./tests/integration/python -count=1` — **blocked by parallel agents**: BP-PY-46 safe fixtures regress under concurrent edits; `tests/integration/python` panics on import from `cwe/fp_guards.go` invalid `\1` regexp (Batch 5). Not in this batch’s exclusive scope.
- Inventory floor in `bp_matrix_test.go` raised to ≥178 (paired fixture count at write time; exact equality softened for parallel appends).

## Remaining uncertainty

- Nox `tests_impl` orchestration left as-is (uncertain / out of helper-delegation pattern).
- Corpus rescan of `real-repos/*` not run (focused fixture tests only, per batch instructions).
