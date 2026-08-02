# False-positive audit report — graphzero

## Run metadata

```yaml
timestamp: 2026-08-02T07:14:37Z
repository: graphzero
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/graphzero
branch: master
commit: 0133f90898dc8b6de9b61994c3de622e591b2012
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/graphzero
chunk_path: ./scripts/chunks
function_context_path: ./scripts/findings/functions
```

## Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir real-repos/graphzero/scripts/chunks -context-dir real-repos/graphzero/scripts/findings/functions real-repos/graphzero`
- Findings: `7`
- Chunks reviewed: `./scripts/chunks/Chunk_1_7.txt`
- Function contexts reviewed: `./scripts/findings/functions/1.txt`, `2.txt`, `3.txt`, `4.txt`, `5.txt`, `6.txt`, `7.txt`

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
| True positive | 7 | 1, 2, 3, 4, 5, 6, 7 |
| Uncertain | 0 | — |

## False positives

None. Every finding matches its rule's detection condition exactly (see `## True positives`).

## Uncertain findings

None.

## True positives

### BP-PY-1 (Bare Except Clause)

| Finding | Source | Reason |
| --- | --- | --- |
| 1 | benchmark/benchmark_papers100M_Pyg.py:33 | `except Exception as e:` is a broad handler in a non-test file whose suite only prints and calls `exit(1)` — no re-raise; neither BP-PY-1 exemption (`isPythonTestFile` + evidence collection, `suiteReraises`) applies. |
| 3 | benchmark/benchmark_papers100M_Pyg.py:56 | `except Exception as e:` broad handler in a non-test file; suite only prints (`print(f"❌ Failed during access: {e}")`), no re-raise. |
| 4 | benchmark/benchmark_papers100M_gz.py:36 | `except Exception as e:` broad handler in a non-test file; suite prints and exits, no re-raise. |

Source excerpt (finding 1, benchmark/benchmark_papers100M_Pyg.py:33):

```python
except Exception as e:
    print(f"\n❌ CRASHED as expected: {e}")
    print(f"💾 RAM at Crash: {get_ram_usage():.4f} GB")
    exit(1)
```

Checklist evidence: the handler is `except Exception` (broad), the file is a benchmark script (path `benchmark/…`, no `test_*.py`/`*_test.py`/`tests/` match for `isPythonTestFile`), and the suite contains no `raise` statement, so `suiteReraises` returns false. Findings 3 and 4 are the same construct in the same files at distinct lines (33/56, 36), so each fires independently on its own source construct.

### CWE-396 (Declaration of Catch for Generic Exception)

| Finding | Source | Reason |
| --- | --- | --- |
| 2 | benchmark/benchmark_papers100M_Pyg.py:33 | `except Exception as e:` matches `pyGenericExceptRE` (`^[\t ]*except\s+(?:Exception|BaseException)(?:\s+as\s+\w+)?\s*:`), and the file is not a Python test module, so the `isPythonTestModule` gate does not suppress it. |
| 5 | benchmark/benchmark_papers100M_gz.py:36 | Same construct at a separate location; matches the regex and is not a test module. |

Checklist evidence: the exact exception token gate (`"except Exception"`) is present and the regex matches the handler line; `isPythonTestModule` returns false for `benchmark/*.py` paths (no `test_`/`_test`/`tests/` component).

### CWE-1084 (Invokable Control Element with Excessive File or Data Access Operations)

| Finding | Source | Reason |
| --- | --- | --- |
| 6 | benchmark/train_graphsage.py:18 | `generate_synthetic_data()` contains exactly 3 `open(...)` calls (lines 27, 31, 37), meeting the rule condition of ≥ 3 `open`/`.execute` calls in a single function body. |

Source excerpt (benchmark/train_graphsage.py:18-39):

```python
def generate_synthetic_data():
    """Generates synthetic CSVs if they don't exist yet."""
    if os.path.exists("dataset/edges.csv"): return
    os.makedirs("dataset", exist_ok=True)
    ...
    with open("dataset/edges.csv", "w") as f:
        for s, d in zip(src, dst): f.write(f"{s},{d}\n")
    with open("dataset/features.csv", "w") as f:
        ...
    with open("dataset/labels.csv", "w") as f:
        ...
```

Checklist evidence: `detectCWE1084` counts `open`/`.execute` calls inside the function body via `findCallsMasked`; three `open()` calls are present (`>= 3`), so the finding fires on the rule condition.

### CWE-367 (Time-of-check Time-of-use (TOCTOU) Race Condition)

| Finding | Source | Reason |
| --- | --- | --- |
| 7 | tests/test.py:93 | `os.path.exists(gd_path)` is immediately followed by `os.remove(gd_path)` on the same path — the exact exists/lexists-then-use pattern matched by `pyTierBTOCTOURE`; the state between check and use can change (e.g. file vanishes → `FileNotFoundError`). |

Source excerpt (tests/test.py:88-97):

```python
@pytest.fixture
def feature_store():
    csv_path = "tests/features.csv"
    gd_path = "tests/features.gd"
    # Compile to binary using FLOAT32
    if os.path.exists(gd_path):
        os.remove(gd_path)
    gz.convert_csv_to_gd(csv_path, gd_path, gz.DataType.FLOAT32)
```

Checklist evidence: `pyTierBTOCTOURE` requires `os.path.(exists|lexists)(<ident>)` followed within 300 chars by `open(`/`os.remove(`/`os.unlink(`; the source satisfies it verbatim, and CWE-367 has no test-file exemption. (Severity is low-impact: only a possible `FileNotFoundError` in a test fixture.)

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `./scripts/chunks/Chunk_1_7.txt` (all 7 findings)
- Function evidence: `./scripts/findings/functions/1.txt` … `7.txt`
- Validation: `git diff --check` — pass
