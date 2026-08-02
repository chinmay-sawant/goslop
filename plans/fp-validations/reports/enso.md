# False-positive audit: enso

## Run metadata

```yaml
timestamp: 2026-08-02T07:42:22Z
repository: enso
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso
branch: master
commit: 516c06caa712e2e454e8673ef5f365616362a9a9
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso
chunk_path: scripts/enso/chunks
function_context_path: scripts/enso/findings/functions
```

## Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/enso/chunks -context-dir scripts/enso/findings/functions real-repos/enso`
- Findings: `138`
- Chunks reviewed: `scripts/enso/chunks/Chunk_1_25.txt`, `scripts/enso/chunks/Chunk_26_50.txt`, `scripts/enso/chunks/Chunk_51_75.txt`, `scripts/enso/chunks/Chunk_76_100.txt`, `scripts/enso/chunks/Chunk_101_125.txt`, `scripts/enso/chunks/Chunk_126_138.txt`
- Function contexts reviewed: `scripts/enso/findings/functions/11.txt`, `14.txt`, `25.txt`, `85.txt`, `95.txt`, `96.txt`, `97.txt`, `98.txt`, `99.txt`, `100.txt`, `101.txt`

## Audit checklist

- [x] Read every assigned chunk under `scripts/enso/chunks`.
- [x] Read `scripts/enso/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient. (No delegated reviews.)
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 11 | 11, 14, 25, 85, 95, 96, 97, 98, 99, 100, 101 |
| True positive | 127 | 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 13, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 86, 87, 88, 89, 90, 91, 92, 93, 94, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138 |
| Uncertain | 0 | — |

## False positives

### [ ] Finding `11` — `CWE-89`

- Function context: `scripts/enso/findings/functions/11.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso/Base.py:210:5`
- Checklist pattern: flagged construct is a `def` statement, not a call

Source excerpt:

```
def execute(source_code, global_vars=None, local_vars=None):
    ''' custom exec for running files, does AST xforms after parsing and imports '''
```

Why this is a false positive: the flagged line is a function *definition* named `execute` (the repo's custom interpreter exec wrapper), not a call site — the detector's `findExecuteCalls` matched the `execute(` token inside `def execute(...)`. No SQL command reaches any sink; the rule's condition ("dynamic SQL string reaches execute/executemany") is not satisfied by a definition.

Checklist evidence: the flagged token appears in a `def` header; there is no call and no SQL on the line or in the enclosing function.

### [ ] Finding `14` — `BP-PY-46`

- Function context: `scripts/enso/findings/functions/14.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso/Builtin.py:537:5`
- Checklist pattern: flagged construct is a `def` statement, not a print call

Source excerpt:

```
def print(*x:T, end="\n") -> Void:
```

Why this is a false positive: the flagged line is a function *definition* that declares the language's `print` builtin; the `print(` match is the parameter list of the `def`. The rule condition flags `print(` *usage* in library code — a definition is not usage, so the condition is not met.

Checklist evidence: the token `print(` matched by `printCallOutsideString` is inside a `def print(...)` header; no print call occurs on the line.

### [ ] Finding `25` — `CWE-89`

- Function context: `scripts/enso/findings/functions/25.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso/Callables.py:895:12`
- Checklist pattern: call target is the repository's own interpreter exec function, not a DB-API execute sink

Source excerpt:

```
def ex(s:Str) -> T:
    from base import execute
    return execute(s, globals(), {})
```

Why this is a false positive: the shown source imports the call target `execute` from the repo's own `base` module, where `Base.py:210` defines it as "custom exec for running files" (an AST-transforming interpreter entry point). The first argument is enso-dialect source text, not an SQL string, and the repo contains no database driver or cursor. The rule condition (dynamic *SQL* reaching a DB execute sink) is not satisfied.

Checklist evidence: the callee is a local import (`from base import execute`) whose definition is a code-execution wrapper; the argument is program source, not an SQL string.

### [ ] Finding `85` — `BP-PY-46`

- Function context: `scripts/enso/findings/functions/85.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso/Test.py:4:1`
- Checklist pattern: print in a test file — the rule's own condition excludes test files

Source excerpt:

```
from Shell import Command as cmd

print("hello i am a test file")
```

Why this is a false positive: the rule's condition is "print used in library code (not under `__main__` guard, not tests)"; this is a standalone test script whose module-level print is test scaffolding. The detector's `isPythonTestFile` exemption failed only because the file is named `Test.py` (capital `T`, not `test_*.py`) and is not under a `tests/` path — the heuristic misses it while the rule intends to skip it.

Checklist evidence: the file is a test script (its own output declares "i am a test file"); the rule condition explicitly excludes tests, so the print is outside the rule's scope.

### [ ] Finding `95` — `BP-PY-46`

- Function context: `scripts/enso/findings/functions/95.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso/Tests/Testing.py:12:5`
- Checklist pattern: print in the test harness (Tests/ directory) — the rule's own condition excludes tests

Source excerpt:

```
def section(section_name:Str) -> Void:
    print()
    header $ f"#### {section_name} ####"
```

Why this is a false positive: the print is test-harness output in `Tests/Testing.py` (the repo's test framework). The rule condition excludes test files; the detector's `isPythonTestFile` misses the path only because its check is case-sensitive (`/tests/` vs the repo's `Tests/`), so the exemption does not apply while the rule intends it to.

Checklist evidence: the file lives in a `Tests/` directory and implements test sections/reports; the print is the harness's user-facing output, which the rule's test exemption covers.

### [ ] Finding `96` — `BP-PY-46`

- Function context: `scripts/enso/findings/functions/96.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso/Tests/Testing.py:14:5`
- Checklist pattern: print in the test harness (same construct family as finding 95)

Source excerpt:

```
    print()
    header $ f"#### {section_name} ####"
    print()
```

Why this is a false positive: same test-harness case as finding 95 — the print is test output in the `Tests/` directory file that the rule's test-file exemption covers; the exemption is defeated only by the case-sensitive `tests/` path check.

Checklist evidence: file is the repo's test framework under `Tests/`; the print is harness formatting output, not library debug logging.

### [ ] Finding `97` — `BP-PY-46`

- Function context: `scripts/enso/findings/functions/97.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso/Tests/Testing.py:25:9`
- Checklist pattern: print in the test harness

Source excerpt:

```
    if failed:
        print(f"{rfill(60, '.', name)} {Color.red('FAILURE')}")
```

Why this is a false positive: FAILURE/PASSED lines are the test framework's report output in `Tests/Testing.py`; the rule's test exemption applies and is defeated only by the case-sensitive `tests/` path check in `isPythonTestFile`.

Checklist evidence: the print renders a test result report inside the repo's test framework file; the rule condition excludes tests.

### [ ] Finding `98` — `BP-PY-46`

- Function context: `scripts/enso/findings/functions/98.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso/Tests/Testing.py:28:9`
- Checklist pattern: print in the test harness

Source excerpt:

```
    else:
        testlist.append(True)
        print(f"{rfill(60, '.', name)} {Color.green('PASSED')}")
```

Why this is a false positive: same as finding 97 — PASSED report output in the test framework; the rule's test-file exemption covers it, and the detector only misses it because of the case-sensitive path check.

Checklist evidence: the print emits a test-result line inside the repo's test framework; the rule condition excludes tests.

### [ ] Finding `99` — `BP-PY-46`

- Function context: `scripts/enso/findings/functions/99.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso/Tests/Testing.py:33:9`
- Checklist pattern: print in the test harness

Source excerpt:

```
    failures = testlist | failing?
    if failures:
        print()
        warning $ "warning, failing tests"
```

Why this is a false positive: final-report formatting output of the test framework in `Tests/Testing.py`; covered by the rule's test exemption, missed only by the case-sensitive `tests/` path check.

Checklist evidence: the print is part of `final_report` in the repo's test framework; the rule condition excludes tests.

### [ ] Finding `100` — `BP-PY-46`

- Function context: `scripts/enso/findings/functions/100.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso/Tests/Testing.py:36:9`
- Checklist pattern: print in the test harness

Source excerpt:

```
        warning $ f"failed {len(failures)}/{len(testlist)} tests"
        print()
```

Why this is a false positive: final-report formatting output of the test framework in `Tests/Testing.py`; covered by the rule's test exemption, missed only by the case-sensitive `tests/` path check.

Checklist evidence: the print is part of `final_report` in the repo's test framework; the rule condition excludes tests.

### [ ] Finding `101` — `BP-PY-46`

- Function context: `scripts/enso/findings/functions/101.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso/Tests/Testing.py:38:13`
- Checklist pattern: print in the test harness

Source excerpt:

```
        for failure in failures:
            print(failure)
        exit(1)
```

Why this is a false positive: failure listing output of the test framework in `Tests/Testing.py`; covered by the rule's test exemption, missed only by the case-sensitive `tests/` path check.

Checklist evidence: the print is part of `final_report` in the repo's test framework; the rule condition excludes tests.

## True positives

All remaining findings satisfy their rule's condition on the shown source. Grouped per rule:

| Rule | Finding IDs | Reason |
| --- | --- | --- |
| BP-PY-5 | 1, 2, 3, 16, 17, 26, 31, 32, 33, 53, 59, 68, 77, 86, 87, 88, 89, 90, 91, 92, 93, 94, 102, 103, 109, 116–138 | `from X import *` in a non-`__init__.py` file (rule condition only exempts files literally named `__init__.py`; `so.py` is an aggregator but not `__init__.py`). |
| BP-PY-12 | 4, 7, 8, 9, 12, 13 | `compile(..., 'exec'/'eval')` / `exec` / `eval` with non-literal arguments (interpreter core — still satisfies the rule condition). |
| BP-PY-12 | 5 | `eval(TYPE_IMPORT)` with a non-literal loop-variable argument. |
| BP-PY-1 | 19, 20, 23, 30, 34, 35, 40, 45, 46, 48, 49, 50, 58, 62, 70, 71, 73, 104, 110, 111, 112, 113, 115 | Bare `except:` or broad `except Exception` whose suite does not re-raise. |
| BP-PY-2 | 36 | `except: pass` — handler body is only `pass`. |
| BP-PY-3 | 21, 28, 54, 56, 57, 60, 65, 66, 75, 78, 80, 81, 82, 83 | `raise Exception(...)` loses type information. |
| BP-PY-6 | 108 | `assert` used for runtime validation in a non-test module. |
| BP-PY-7 | 41, 42, 43, 44, 51 | `open(...)` without a `with` statement (context-manager condition met; e.g. `File.py:52` closes only after the loop). |
| BP-PY-46 | 15, 27, 64, 67, 69, 84, 105 | Real `print(` calls in library modules (builtin wrapper, color/console helpers, REPL loop, etc.). |
| CWE-94 | 6, 39 | `eval` reached by a non-literal/dynamic expression (`TYPE_IMPORT`; `f"o.{p}"`). |
| CWE-1333 | 107 | `re.compile(r"[^'\\]*(?:\\.[^'\\]*)*'")` contains a nested unbounded quantifier shape `(?:…*)*`. |
| CWE-396 | 10, 24, 47, 63, 72 | `except Exception` generic handlers. |
| CWE-397 | 22, 29, 55, 61, 76, 79 | `raise Exception` generic raise. |
| CWE-390 | 37 | Exception detected but handler is `pass`. |
| CWE-1071 | 38 | Exception handler contains only `pass`. |
| CWE-1046 | 74 | `x += c` string accumulation inside a loop with prior `x = ""` evidence. |
| CWE-1121 | 18, 106, 114 | Function bodies with ≥ 12 counted control-flow branches (14 in `Callables.__call__`, tokenizer, `FunctionType.__eq__`). |
| PERF-PY-26 | 52 | `parse_macro(` called inside a loop without cache evidence. |

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/enso/chunks`
- Function evidence: `scripts/enso/findings/functions`
- Validation: `git diff --check` — pass

## Post-fix remaining-FP audit (2026-08-02)

### Run metadata

```yaml
timestamp: 2026-08-02T16:30:00Z (fresh scan; binary rebuilt 2026-08-02 16:29, fix b5b8fde)
repository: enso
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso
branch: master
commit: 516c06caa712e2e454e8673ef5f365616362a9a9
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso
chunk_path: scripts/enso/chunks
function_context_path: scripts/enso/findings/functions
```

### Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/enso/chunks -context-dir scripts/enso/findings/functions real-repos/enso`
- Findings: `134`
- Chunks reviewed: `scripts/enso/chunks/Chunk_1_25.txt`, `Chunk_26_50.txt`, `Chunk_51_75.txt`, `Chunk_76_100.txt`, `Chunk_101_125.txt`, `Chunk_126_134.txt`
- Function contexts reviewed: `scripts/enso/findings/functions/12.txt`, `81.txt`, `91.txt`, `92.txt`, `93.txt`, `94.txt`, `95.txt`, `96.txt`, `97.txt`

### Audit checklist

- [x] Read every assigned chunk under `scripts/enso/chunks`.
- [x] Read `scripts/enso/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient. (No delegated reviews.)
- [x] Ran `git diff --check` after updating this report.

### Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 9 | 12, 81, 91, 92, 93, 94, 95, 96, 97 |
| True positive | 125 | 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 82, 83, 84, 85, 86, 87, 88, 89, 90, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134 |
| Uncertain | 0 | — |

Fresh findings were matched to the audited findings by `Source:` (file:line:col). All 9 false positives below are re-appearing audited FPs (old IDs 14, 85, 95, 96, 97, 98, 99, 100, 101). The remaining 125 findings match audited TP sources one-to-one (per-rule counts agree with the audited TP grouping: BP-PY-5 48, BP-PY-12 7, BP-PY-1 23, BP-PY-2 1, BP-PY-3 14, BP-PY-6 1, BP-PY-7 5, BP-PY-46 7, CWE-94 2, CWE-1333 1, CWE-397 6, CWE-390 1, CWE-1071 1, CWE-1046 1, CWE-1121 3) except CWE-396, which is now 4 of the audited 5 (see review item below).

Two audited FPs were fixed by `b5b8fde` and no longer fire: CWE-89 `Base.py:210:5` (old FP 11, `def execute` header) and CWE-89 `Callables.py:895:12` (old FP 25, non-SQL `execute` call) — the new non-SQL-`execute` guard suppresses both. No new findings were introduced.

Review item (over-suppression, Mode B concern): two audited TPs are absent from the fresh scan while their source constructs remain present —
- CWE-396 `Base.py:206:9` (old TP 10): handler body is `raise`; the new `suiteSurfacesFailureMasked` guard skips re-raising handlers, so this audited TP is suppressed-but-present. The suppression matches the rule's intent (a re-raise does not hide the failure) and is flagged for reviewer confirmation.
- PERF-PY-26 `Macro.py:180` (old TP 52): `macros.append(parse_macro(macro))` inside the `while ...:` loop is now skipped because `parse_*` calls only fire inside explicit `handle_job`/`handle_request` windows (`perfLineIsParseCall`); the construct is still in the source, suppressed-but-present.

### False positives

### [ ] Finding `12` — `BP-PY-46`

- Function context: `scripts/enso/findings/functions/12.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso/Builtin.py:537:5`
- Checklist pattern: flagged construct is a `def` statement, not a print call (re-appearing audited FP 14)

Source excerpt:

```
def print(*x:T, end="\n") -> Void:
```

Why this is a false positive: the flagged line is the declaration of the language's `print` builtin; the `print(` token is the `def` parameter list, not a call site. The fix's def-header guards were applied to BP-PY-12/7 and CWE-89/94, not to BP-PY-46, so the rule still fires on the definition.

Checklist evidence: the flagged token is inside a `def print(...)` header; no `print(` call occurs on the line, so the rule condition ("print used in library code") is not met by usage.

### [ ] Finding `81` — `BP-PY-46`

- Function context: `scripts/enso/findings/functions/81.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso/Test.py:4:1`
- Checklist pattern: print in a test file — the rule's own test exemption is defeated by the case-sensitive path check (re-appearing audited FP 85)

Source excerpt:

```
from Shell import Command as cmd

print("hello i am a test file")
```

Why this is a false positive: this is a standalone test script (its own output declares "i am a test file"). `isPythonTestFile` still does not match it: the filename `Test.py` matches neither `test_*.py` nor `*_test.py`, and the `Tests/`/`Test.py` paths do not match the case-sensitive `tests/`/`test/` path components, so the rule's test exemption does not apply while the rule intends it to.

Checklist evidence: `isPythonTestFile` (common.go:24) only matches `test_*.py`, `*_test.py`, `tests.py`, `conftest.py`, and `tests/`/`test/` path components; `Test.py` satisfies none of them.

### [ ] Finding `91` — `BP-PY-46`

- Function context: `scripts/enso/findings/functions/91.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso/Tests/Testing.py:12:5`
- Checklist pattern: print in the test harness (`Tests/` directory) — case-sensitive path exemption miss (re-appearing audited FP 95)

Source excerpt:

```
def section(section_name:Str) -> Void:
    print()
    header $ f"#### {section_name} ####"
```

Why this is a false positive: the print is test-harness output in `Tests/Testing.py`, the repo's test framework; the rule's test exemption covers it but is defeated only by the case-sensitive `tests/` path check in `isPythonTestFile` (`Tests/` ≠ `tests/`).

Checklist evidence: the file lives in a `Tests/` directory and implements test sections; the rule condition excludes tests, so the print is outside the rule's scope.

### [ ] Finding `92` — `BP-PY-46`

- Function context: `scripts/enso/findings/functions/92.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso/Tests/Testing.py:14:5`
- Checklist pattern: print in the test harness (same exemption-miss as finding 91) (re-appearing audited FP 96)

Source excerpt:

```
    print()
    header $ f"#### {section_name} ####"
    print()
```

Why this is a false positive: same test-harness case as finding 91 — the print is section formatting output in the `Tests/` test framework file that the rule's test exemption covers.

Checklist evidence: file is the repo's test framework under `Tests/`; the exemption applies and is missed only by the case-sensitive path check.

### [ ] Finding `93` — `BP-PY-46`

- Function context: `scripts/enso/findings/functions/93.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso/Tests/Testing.py:25:9`
- Checklist pattern: print in the test harness (re-appearing audited FP 97)

Source excerpt:

```
    if failed:
        print(f"{rfill(60, '.', name)} {Color.red('FAILURE')}")
```

Why this is a false positive: FAILURE/PASSED lines are the test framework's report output in `Tests/Testing.py`; covered by the rule's test exemption, missed only by the case-sensitive `tests/` path check.

Checklist evidence: the print renders a test-result report inside the repo's test framework file; the rule condition excludes tests.

### [ ] Finding `94` — `BP-PY-46`

- Function context: `scripts/enso/findings/functions/94.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso/Tests/Testing.py:28:9`
- Checklist pattern: print in the test harness (re-appearing audited FP 98)

Source excerpt:

```
    else:
        testlist.append(True)
        print(f"{rfill(60, '.', name)} {Color.green('PASSED')}")
```

Why this is a false positive: same as finding 93 — PASSED report output in the test framework; the rule's test-file exemption covers it, and the detector only misses it because of the case-sensitive path check.

Checklist evidence: the print emits a test-result line inside the repo's test framework; the rule condition excludes tests.

### [ ] Finding `95` — `BP-PY-46`

- Function context: `scripts/enso/findings/functions/95.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso/Tests/Testing.py:33:9`
- Checklist pattern: print in the test harness (re-appearing audited FP 99)

Source excerpt:

```
    failures = testlist | failing?
    if failures:
        print()
        warning $ "warning, failing tests"
```

Why this is a false positive: final-report formatting output of the test framework in `Tests/Testing.py`; covered by the rule's test exemption, missed only by the case-sensitive `tests/` path check.

Checklist evidence: the print is part of `final_report` in the repo's test framework; the rule condition excludes tests.

### [ ] Finding `96` — `BP-PY-46`

- Function context: `scripts/enso/findings/functions/96.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso/Tests/Testing.py:36:9`
- Checklist pattern: print in the test harness (re-appearing audited FP 100)

Source excerpt:

```
        warning $ f"failed {len(failures)}/{len(testlist)} tests"
        print()
```

Why this is a false positive: final-report formatting output of the test framework in `Tests/Testing.py`; covered by the rule's test exemption, missed only by the case-sensitive `tests/` path check.

Checklist evidence: the print is part of `final_report` in the repo's test framework; the rule condition excludes tests.

### [ ] Finding `97` — `BP-PY-46`

- Function context: `scripts/enso/findings/functions/97.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso/Tests/Testing.py:38:13`
- Checklist pattern: print in the test harness (re-appearing audited FP 101)

Source excerpt:

```
        for failure in failures:
            print(failure)
        exit(1)
```

Why this is a false positive: failure-listing output of the test framework in `Tests/Testing.py`; covered by the rule's test exemption, missed only by the case-sensitive `tests/` path check.

Checklist evidence: the print is part of `final_report` in the repo's test framework; the rule condition excludes tests.

### Uncertain findings

None.

### Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/enso/chunks`
- Function evidence: `scripts/enso/findings/functions`
- Validation: `git diff --check` — pass

## Post-fix v2 audit (latest binary)

### Run metadata

```yaml
timestamp: 2026-08-02T18:05:00Z (fresh scan; binary rebuilt 2026-08-02 17:56)
repository: enso
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso
branch: master
commit: 516c06caa712e2e454e8673ef5f365616362a9a9
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso
chunk_path: scripts/enso/chunks
function_context_path: scripts/enso/findings/functions
```

### Scan evidence

- Build command: `make build` (~17:56)
- Findings: `135`
- Chunks reviewed: `scripts/enso/chunks/Chunk_1_25.txt`, `Chunk_26_50.txt`, `Chunk_51_75.txt`, `Chunk_76_100.txt`, `Chunk_101_125.txt`, `Chunk_126_135.txt` (all fresh, read in full)
- Function contexts reviewed: every fresh finding's exported context was in its chunk; enclosing source read for `Base.py:206`, `Structure.py:115`, `Macro.py:173-180`

### Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 9 | 13, 82, 92, 93, 94, 95, 96, 97, 98 |
| True positive | 126 | all other fresh findings (1–12, 14–81, 83–91, 99–135) |
| Uncertain | 0 | — |

All 135 fresh findings were matched by `Source:` to prior audit entries; the 9 FPs are the same re-appearing audited FPs (old IDs 14, 85, 95–101; Mode B IDs 12, 81, 91–97). Per-rule counts are identical to Mode B (BP-PY-5 48, BP-PY-12 7, BP-PY-1 23, BP-PY-2 1, BP-PY-3 14, BP-PY-6 1, BP-PY-7 5, BP-PY-46 7TP/9FP, CWE-94 2, CWE-1333 1, CWE-397 6, CWE-390 1, CWE-1071 1, CWE-1046 1, CWE-1121 3) except CWE-396, now 5 (Mode B: 4). The delta is fresh finding `10` — CWE-396 `Base.py:206:1` (`except Exception as e: raise`) re-appearing as an audited TP (old TP 10); the Mode B "suppressed-but-present" over-suppression noted there is resolved by the latest binary (re-raise handlers are flagged again, matching the rule's intent). PERF-PY-26 `Macro.py:180` is still suppressed-but-present (construct in source at `parse_macro_file`, no fresh finding). CWE-89 (old FPs 11, 25) stays fixed — zero CWE-89 findings.

## Fix checklist (FP patterns)

| Pattern # | Rule | Trigger shape | Count | Example sources |
| --- | --- | --- | ---: | --- |
| 1 | BP-PY-46 | `print(` token inside a `def print(...)` header (declaration) matched as a call site | 1 | `Builtin.py:537:5` |
| 2 | BP-PY-46 | `print(...)` in test files whose basename/path is a case-variant of `test` (`Test.py`, `Tests/`) — `isPythonTestFile` path/name check is case-sensitive and misses them | 8 | `Test.py:4:1`; `Tests/Testing.py:12:5`, `14:5`, `25:9`, `28:9`, `33:9`, `36:9`, `38:13` |

Pattern 1 fix condition: skip tokens that are part of a `def`/`class` header (parameter list), i.e. only flag actual `print(` call expressions. Pattern 2 fix condition: the test exemption should match `test` case-insensitively in the basename (`Test.py`, `Testing.py`) and in path components (`Tests/`), or skip any file under a directory/basename whose lowercased form contains `test`.

## New findings

None — every fresh finding has a prior classification (9 re-appearing audited FPs, 126 audited TPs). Only delta vs Mode B: CWE-396 `Base.py:206:1` (old TP 10) recovered from the b5b8fde over-suppression; PERF-PY-26 `Macro.py:180` remains suppressed-but-present.
