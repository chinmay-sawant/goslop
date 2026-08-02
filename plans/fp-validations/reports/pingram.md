# False-positive audit — pingram

## Run metadata

```yaml
timestamp: 2026-08-02T07:14:11Z
repository: pingram
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pingram
branch: main
commit: e06eed8e008c79afa56970b9a576757df20630ed
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pingram
chunk_path: ./scripts/chunks
function_context_path: ./scripts/findings/functions
```

## Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir real-repos/pingram/scripts/chunks -context-dir real-repos/pingram/scripts/findings/functions real-repos/pingram`
- Findings: `6`
- Chunks reviewed: `./scripts/chunks/Chunk_1_6.txt`
- Function contexts reviewed: `./scripts/findings/functions/1.txt`, `2.txt`, `3.txt`, `4.txt`, `5.txt`, `6.txt`

## Audit checklist

- [x] Read every assigned chunk under `./scripts/chunks`.
- [x] Read `./scripts/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 0 | — |
| True positive | 6 | 1, 2, 3, 4, 5, 6 |
| Uncertain | 0 | — |

## False positives

No false positives found. Every finding satisfies the literal rule condition of the reported rule (verified against `internal/lang/python/detectors/bad_practices/rules_core.go` for BP-PY-1/BP-PY-2 and `internal/lang/python/detectors/cwe/rules_platform.go` for CWE-396/CWE-390), and none of the built-in exemptions (re-raise, test-module, test-evidence collection) applies to the shown source.

## True positives

### BP-PY-1 — Bare Except Clause

| Finding | Source | Reason |
| --- | --- | --- |
| 1 | `src/pingram/_errors.py:73` | `except Exception:` in non-test module whose suite is only `return None` — neither re-raise (`suiteReraises`) nor the test-evidence exemption applies. |
| 3 | `src/pingram/_retry.py:140` | Same construct: `except Exception:` with `return None`, no re-raise, non-test file. |

### CWE-396 — Declaration of Catch for Generic Exception

| Finding | Source | Reason |
| --- | --- | --- |
| 2 | `src/pingram/_errors.py:73` | Matches `pyGenericExceptRE` (`except Exception:`); `isPythonTestModule` is false for `src/pingram/_errors.py`. |
| 4 | `src/pingram/_retry.py:140` | Matches `pyGenericExceptRE` (`except Exception:`); `isPythonTestModule` is false for `src/pingram/_retry.py`. |

### BP-PY-2 — Except Pass

| Finding | Source | Reason |
| --- | --- | --- |
| 5 | `tests/unit/conftest.py:66` | Except suite contains only `pass` after comment masking (`suite[0] == "pass"`); the rule has no test-file exemption. |

### CWE-390 — Detection of Error Condition Without Action

| Finding | Source | Reason |
| --- | --- | --- |
| 6 | `tests/unit/conftest.py:66` | `exceptPassStart` finds an except clause whose first direct body statement is `pass`; the rule has no test-file exemption. |

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `./scripts/chunks/Chunk_1_6.txt`
- Function evidence: `./scripts/findings/functions/1.txt` … `6.txt`
- Validation: `git diff --check` — pass
