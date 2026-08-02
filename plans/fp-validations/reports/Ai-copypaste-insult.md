# False-positive audit — Ai-copypaste-insult

## Run metadata

```yaml
timestamp: 2026-08-02T07:11:55Z
repository: Ai-copypaste-insult
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Ai-copypaste-insult
branch: main
commit: 77f3d97d832566ce37688acedf7c4096f168beb4
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Ai-copypaste-insult
chunk_path: scripts/Ai-copypaste-insult/chunks
function_context_path: scripts/Ai-copypaste-insult/findings/functions
```

## Scan evidence

- Build command: `not provided in audit prompt`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/Ai-copypaste-insult/chunks -context-dir scripts/Ai-copypaste-insult/findings/functions real-repos/Ai-copypaste-insult`
- Findings: `2`
- Chunks reviewed: `scripts/Ai-copypaste-insult/chunks/Chunk_1_2.txt`
- Function contexts reviewed: `scripts/Ai-copypaste-insult/findings/functions/1.txt`, `scripts/Ai-copypaste-insult/findings/functions/2.txt`

## Audit checklist

- [x] Read every assigned chunk under `scripts/Ai-copypaste-insult/chunks`.
- [x] Read `scripts/Ai-copypaste-insult/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 0 | — |
| True positive | 2 | 1, 2 |
| Uncertain | 0 | — |

## False positives

No findings were classified as false positives.

## True positives

### BP-PY-1 — Bare Except Clause

| Finding ID | Source | Reason |
| --- | --- | --- |
| 1 | `src/detector/test.py:30:1` | `except Exception as e:` suite only prints a fixed fallback and discards the exception (`e` never used); it neither re-raises nor is it a test-module evidence collector, so the broad-except condition fires. |

### CWE-396 — Declaration of Catch for Generic Exception

| Finding ID | Source | Reason |
| --- | --- | --- |
| 2 | `src/detector/test.py:30:1` | `except Exception as e:` matches the `pyGenericExceptRE` trigger (`except\s+Exception\s+as\s+e\s*:`) in a module not recognized as a test module, and the generic handler swallows the failure conditions. |

## Uncertain findings

No findings were classified as uncertain.

## Final evidence

- Delegated reviewers: `none`
- Chunk evidence: `scripts/Ai-copypaste-insult/chunks`
- Function evidence: `scripts/Ai-copypaste-insult/findings/functions`
- Validation: `git diff --check` — `pass`
