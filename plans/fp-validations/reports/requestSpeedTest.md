# False-positive audit — requestSpeedTest

## Run metadata

```yaml
timestamp: 2026-08-02T07:15:28Z
repository: requestSpeedTest
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/requestSpeedTest
branch: main
commit: 60626f3d548258a76c8df962db1b26aa9baa0e87
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/requestSpeedTest
chunk_path: ./scripts/chunks
function_context_path: ./scripts/findings/functions
```

## Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop` (goslop binary used: `./bin/goslop`)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir real-repos/requestSpeedTest/scripts/chunks -context-dir real-repos/requestSpeedTest/scripts/findings/functions real-repos/requestSpeedTest`
- Findings: `13`
- Chunks reviewed: `./scripts/chunks/Chunk_1_13.txt`
- Function contexts reviewed: `./scripts/findings/functions/1.txt` .. `./scripts/findings/functions/13.txt`

## Audit checklist

- [x] Read every assigned chunk under `./scripts/chunks`.
- [x] Read `./scripts/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [ ] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 4 | 1, 2, 3, 10 |
| True positive | 9 | 4, 5, 6, 7, 8, 9, 11, 12, 13 |
| Uncertain | 0 | — |

## False positives

### [x] Finding `1` — BP-PY-42

- Function context: `./scripts/findings/functions/1.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/requestSpeedTest/send_request/httpx_test.py:11:1`
- Checklist pattern: rule condition not satisfied — the try/except is not used "instead of assertRaises/pytest.raises"

Source excerpt:

```python
async def send_request(client, url):
    try:
        resp = await client.get(url)
        return resp.status_code
    except Exception as e:
        log.exception(e)
        raise e
```

Why this is a false positive: The rule's condition is that a test uses a bare try/except to *expect failure* instead of `assertRaises`/`pytest.raises`. Here the except suite re-raises the exception (`raise e`) and logs it — the failure is treated as an error to propagate, not as an expected test outcome, and the file contains no assertions at all. The trigger is only the `*_test.py` file-name heuristic.

Checklist evidence: the except handler re-raises (`raise e`), so the construct is a logging/propagation wrapper, not a failure-expectation substitute; the "instead of assertRaises/pytest.raises" condition is unmet.

### [x] Finding `2` — BP-PY-42

- Function context: `./scripts/findings/functions/2.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/requestSpeedTest/send_request/rnet_test.py:32:1`
- Checklist pattern: rule condition not satisfied — the try/except is error accounting, not an expectation assertion

Source excerpt:

```python
        try:
            resp = await client.get(url)
            status = resp.status.as_int() if hasattr(resp, "status") else getattr(resp, "status_code", None)
            local_statuses[status] += 1
            local_success += 1
        except Exception as e:
            name = type(e).__name__
            local_errors[name] += 1
            local_fail += 1
            if local_errors[name] == 1:
                log.error(f"Worker {wid} error: {e}")
```

Why this is a false positive: The try/except is inside `worker()`, a load-test loop that deliberately counts request failures as benchmark metrics; it is not asserting that an exception occurs. The file is a benchmark CLI (argparse + `main()` + `if __name__ == "__main__"`) that merely carries a `_test.py` name, and it contains no assertions.

Checklist evidence: the except suite records error statistics and logs the first error — there is no assertion intent and no assertion anywhere in the file, so the construct is not a substitute for `assertRaises`/`pytest.raises`.

### [x] Finding `3` — BP-PY-1

- Function context: `./scripts/findings/functions/3.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/requestSpeedTest/send_request/rnet_test.py:37:1`
- Checklist pattern: rule condition not satisfied — the handler actively handles the exception instead of swallowing it

Source excerpt:

```python
        except Exception as e:
            name = type(e).__name__
            local_errors[name] += 1
            local_fail += 1
            if local_errors[name] == 1:
                log.error(f"Worker {wid} error: {e}")
    return local_success, local_fail, local_statuses, local_errors
```

Why this is a false positive: The rule condition is a broad `except Exception` "without handling or re-raise [that] swallows failures and hides bugs". Here the handler deliberately handles the failure: it records the exception type, increments the failure/error counters that are returned and printed by the caller, and logs the first occurrence. The failure is surfaced, not swallowed or hidden.

Checklist evidence: the except suite performs explicit error handling (accounting + logging) and the results propagate to the caller, so the "without handling ... swallows failures" clause of the rule condition is not satisfied.

### [x] Finding `10` — BP-PY-1

- Function context: `./scripts/findings/functions/10.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/requestSpeedTest/utils/increase_limits.py:38:1`
- Checklist pattern: rule condition not satisfied — the handler surfaces the error instead of swallowing it

Source excerpt:

```python
            except Exception as e:
                print(f"Failed to set a high limit: {e}")
```

Why this is a false positive: The rule condition requires a broad except "without handling or re-raise" that swallows failures. This handler deliberately reports the failure to the user with the exception detail (and is the final fallback of a two-step limit-raising strategy), so the error is surfaced rather than hidden.

Checklist evidence: the except suite prints the failure message, i.e., the exception is handled by explicit reporting, satisfying the "handling" clause that exempts the construct from the rule.

## True positives

### BP-PY-46 — print used in library code (7 findings)

`utils/increase_limits.py` is a library module imported by other modules in the repo (e.g. `send_request/rnet_test.py:12` `from utils.increase_limits import set_max_open_files`, called at module level at line 19). All seven `print()` calls sit inside `set_max_open_files()` — outside the `if __name__ == "__main__"` guard — and the file has no argparse/`main()` CLI exemption, so the printed output fires in library/import context and pollutes the host program's stdout. The rule condition ("`print` is used for operational logging in non-script modules", keep print under the `__main__` guard for CLIs) is satisfied.

| Finding | Source | Reason |
| --- | --- | --- |
| 5 | `utils/increase_limits.py:12:9` | `print(...)` in library function `set_max_open_files()`, not under `__main__` guard |
| 6 | `utils/increase_limits.py:21:9` | `print(f"Current limits: ...")` in library function, not under `__main__` guard |
| 7 | `utils/increase_limits.py:27:13` | `print(f"Successfully raised soft limit ...")` in library function, not under `__main__` guard |
| 8 | `utils/increase_limits.py:29:13` | `print(f"Could not raise soft limit ...")` in library function, not under `__main__` guard |
| 9 | `utils/increase_limits.py:37:17` | `print(f"Successfully raised soft limit ...")` in library function, not under `__main__` guard |
| 12 | `utils/increase_limits.py:39:17` | `print(f"Failed to set a high limit: {e}")` in library function, not under `__main__` guard |
| 13 | `utils/increase_limits.py:42:9` | `print(...)` in library function, not under `__main__` guard |

### CWE-117 — Improper Output Neutralization for Logs (1 finding)

| Finding | Source | Reason |
| --- | --- | --- |
| 4 | `send_request/rnet_test.py:42:17` | `log.error(f"Worker {wid} error: {e}")` formats a dynamic value (exception from a network client call, i.e. potentially externally influenced) directly into the log message without CRLF neutralization — the exact construct the rule condition targets |

### CWE-396 — Declaration of Catch for Generic Exception (1 finding)

| Finding | Source | Reason |
| --- | --- | --- |
| 11 | `utils/increase_limits.py:38:1` | Generic `except Exception as e:` around `resource.setrlimit` hides distinct failure conditions (e.g. `ValueError` vs `OSError`); the surrounding code itself distinguishes `ValueError` at line 28, so the broad catch genuinely merges conditions that merit distinct handling |

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `./scripts/chunks/Chunk_1_13.txt`
- Function evidence: `./scripts/findings/functions/1.txt` .. `13.txt`
- Validation: `git diff --check` — pass
