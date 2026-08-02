# False-positive audit — PyDepends

## Run metadata

```yaml
timestamp: 2026-08-02T07:14:55Z
repository: PyDepends
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/PyDepends
branch: main
commit: bc1e400990bb6b50a76da423fbfad6df8074c048
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/PyDepends
chunk_path: scripts/PyDepends/chunks
function_context_path: scripts/PyDepends/findings/functions
```

## Scan evidence

- Build command: `bin/goslop` prebuilt (`make build`)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/PyDepends/chunks -context-dir scripts/PyDepends/findings/functions real-repos/PyDepends`
- Findings: `4`
- Chunks reviewed: `scripts/PyDepends/chunks/Chunk_1_4.txt`
- Function contexts reviewed: `scripts/PyDepends/findings/functions/1.txt`, `2.txt`, `3.txt`, `4.txt`

## Audit checklist

- [x] Read every assigned chunk under `scripts/PyDepends/chunks`.
- [x] Read `scripts/PyDepends/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 0 | — |
| True positive | 4 | 1, 2, 3, 4 |
| Uncertain | 0 | — |

## False positives

No findings were classified as false positives.

## True positives

All four findings sit on the two `except …: pass` handlers used to close generator dependencies in `pydepends/depends.py`. Each rule's condition is purely structural — it matches any `except` suite whose sole body statement is `pass` — and none of the three rules exempts `StopIteration`/`StopAsyncIteration` or generator-teardown `finally` blocks. The shown source satisfies each condition literally, and the audit instructions forbid overriding the rule condition with context knowledge, so all four are true positives.

### BP-PY-2 — Except Pass

Rule condition (from `detectBPPY2` in `internal/lang/python/detectors/bad_practices/rules_core.go`): scan lines for `except …:` clauses; flag when the immediate suite is exactly one statement and that statement is `pass`. Detection notes explicitly direct flagging library code ("Flag library and request-path code more strictly than top-level CLI mains") — `pydepends/depends.py` is library code. Both findings satisfy `len(suite) == 1 && suite[0] == "pass"`.

| Finding id | Source | One-line reason |
| --- | --- | --- |
| 1 | `pydepends/depends.py:108:1` | `except StopIteration: pass` in the `finally` of `_managed_dependency` — suite is solely `pass` in a library module; the rule condition has no exception-type or teardown exemption. |
| 4 | `pydepends/depends.py:120:1` | `except StopAsyncIteration: pass` in the `finally` of `_async_managed_dependency` — same condition satisfied on a distinct line. |

### CWE-390 — Detection of Error Condition Without Action

Rule condition (from `detectCWE390` in `internal/lang/python/detectors/cwe/rules_platform.go`): "recognizes only an except clause whose direct body is pass. It does not infer whether logging, recovery, re-raising, or a caller's behaviour is sufficient error handling." The source matches `exceptPassStart` exactly.

| Finding id | Source | One-line reason |
| --- | --- | --- |
| 2 | `pydepends/depends.py:108:1` | `except StopIteration:` whose only body line is `pass` — direct-body-pass matches the structural trigger; the rule deliberately does not judge whether the handling is adequate. |

### CWE-1071 — Empty Code Block

Rule condition (from `detectCWE1071` in `internal/lang/python/detectors/cwe/rules_tier_b_runtime.go`): regex `except(?:\s+[A-Za-z_][A-Za-z0-9_.]*)?\s*:\s*\n\s*pass\b` matches an exception handler containing only `pass`. The source matches.

| Finding id | Source | One-line reason |
| --- | --- | --- |
| 3 | `pydepends/depends.py:108:9` | `except StopIteration:` immediately followed by a `pass` body — empty-code-block trigger matches verbatim. |

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/PyDepends/chunks/Chunk_1_4.txt`
- Function evidence: `scripts/PyDepends/findings/functions/1.txt`–`4.txt`
- Validation: `git diff --check` — pass
