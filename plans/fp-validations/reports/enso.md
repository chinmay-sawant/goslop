# False-positive audit: enso

## Run metadata

```yaml
timestamp: 2026-08-02T07:42:22Z
repository: enso
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso
branch: master
commit: 516c06caa712e2e454e8673ef5f365616362a9a9
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso
chunk_path: ./scripts/chunks
function_context_path: ./scripts/findings/functions
```

## Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir real-repos/enso/scripts/chunks -context-dir real-repos/enso/scripts/findings/functions real-repos/enso`
- Findings: `138`
- Chunks reviewed: `./scripts/chunks/Chunk_1_25.txt`, `./scripts/chunks/Chunk_26_50.txt`, `./scripts/chunks/Chunk_51_75.txt`, `./scripts/chunks/Chunk_76_100.txt`, `./scripts/chunks/Chunk_101_125.txt`, `./scripts/chunks/Chunk_126_138.txt`
- Function contexts reviewed: `./scripts/findings/functions/11.txt`, `14.txt`, `25.txt`, `85.txt`, `95.txt`, `96.txt`, `97.txt`, `98.txt`, `99.txt`, `100.txt`, `101.txt`

## Audit checklist

- [x] Read every assigned chunk under `./scripts/chunks`.
- [x] Read `./scripts/findings/functions/<finding-id>.txt` for every proposed false positive.
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

- Function context: `./scripts/findings/functions/11.txt`
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

- Function context: `./scripts/findings/functions/14.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/enso/Builtin.py:537:5`
- Checklist pattern: flagged construct is a `def` statement, not a print call

Source excerpt:

```
def print(*x:T, end="\n") -> Void:
```

Why this is a false positive: the flagged line is a function *definition* that declares the language's `print` builtin; the `print(` match is the parameter list of the `def`. The rule condition flags `print(` *usage* in library code — a definition is not usage, so the condition is not met.

Checklist evidence: the token `print(` matched by `printCallOutsideString` is inside a `def print(...)` header; no print call occurs on the line.

### [ ] Finding `25` — `CWE-89`

- Function context: `./scripts/findings/functions/25.txt`
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

- Function context: `./scripts/findings/functions/85.txt`
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

- Function context: `./scripts/findings/functions/95.txt`
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

- Function context: `./scripts/findings/functions/96.txt`
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

- Function context: `./scripts/findings/functions/97.txt`
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

- Function context: `./scripts/findings/functions/98.txt`
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

- Function context: `./scripts/findings/functions/99.txt`
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

- Function context: `./scripts/findings/functions/100.txt`
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

- Function context: `./scripts/findings/functions/101.txt`
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
- Chunk evidence: `./scripts/chunks`
- Function evidence: `./scripts/findings/functions`
- Validation: `git diff --check` — pass
