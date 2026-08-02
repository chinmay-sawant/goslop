# False-positive audit — FlashySurf

## Run metadata

```yaml
timestamp: 2026-08-02T07:14:56Z
repository: FlashySurf
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FlashySurf
branch: main
commit: 3d14f8ca328e7f11d3fc6f2fe86bc15a2463184e
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FlashySurf
chunk_path: ./scripts/chunks
function_context_path: ./scripts/findings/functions
```

## Scan evidence

- Build command: `make build` (bin/goslop prebuilt)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir real-repos/FlashySurf/scripts/chunks -context-dir real-repos/FlashySurf/scripts/findings/functions real-repos/FlashySurf`
- Findings: `5`
- Chunks reviewed: `./scripts/chunks/Chunk_1_5.txt`
- Function contexts reviewed: `./scripts/findings/functions/1.txt`, `2.txt`, `3.txt`, `4.txt`, `5.txt`

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
| False positive | 4 | 1, 2, 4, 5 |
| True positive | 1 | 3 |
| Uncertain | 0 | — |

## False positives

### [ ] Finding `1` — `BP-PY-46`

- Function context: `./scripts/findings/functions/1.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FlashySurf/data-process.py:167:1`
- Checklist pattern: print sits at module scope (zero indentation) of a standalone script; the file is not a non-script module, so the rule's stated condition ("print is used for operational logging in non-script modules") is not satisfied.

Source excerpt (from `data-process.py`):

```
        output["english"].append(cleanUp(flashcard))

print(f"Math questions: {len(output['math'])}/{mathCount}")
print(f"English questions: {len(output['english'])}/{englishCount}")
```

Why this is a false positive: the print at line 167 is at module level (column 0) in a standalone script that runs top-to-bottom — file I/O (`with open("cb-digital-questions.json", "r")` at line 22), the whole transformation loop, and the output writes at lines 172–176 all execute at import time with no `if __name__ == "__main__"` guard. The print is the script's intended summary output, and the file is a script, not a "non-script module" the rule targets; the detector has no script/library distinction and flags every unguarded module-level `print(`.

Checklist evidence: the flagged line is a zero-indent, module-level `print(` in a file whose entire body is top-level executable code (no guard, no importable module API), so the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `2` — `BP-PY-46`

- Function context: `./scripts/findings/functions/2.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FlashySurf/data-process.py:168:1`
- Checklist pattern: same construct as finding 1 — module-scope print in the same standalone script, adjacent to line 167; distinct line, so kept as its own finding.

Source excerpt (from `data-process.py`):

```
print(f"Math questions: {len(output['math'])}/{mathCount}")
print(f"English questions: {len(output['english'])}/{englishCount}")
```

Why this is a false positive: same reasoning as finding 1 — line 168 is a zero-indent module-level print in the standalone `data-process.py` script (no `__main__` guard, no library API), producing the script's intended output, so the rule condition ("print in non-script modules") is not met.

Checklist evidence: the flagged line is a module-level `print(` (column 0) in a top-to-bottom script with no `__main__` guard, so the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `4` — `BP-PY-46`

- Function context: `./scripts/findings/functions/4.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FlashySurf/semantic-classification.py:42:1`
- Checklist pattern: print sits at module scope (zero indentation) of a standalone script; the file is not a non-script module, so the rule's stated condition is not satisfied.

Source excerpt (from `semantic-classification.py`):

```
n_clusters = min(87, n_samples - 1)  # 87 clusters as requested

print(f"Dataset size: {n_samples} questions")
print(f"Using {n_neighbors} neighbors, {n_components} components, {n_clusters} clusters")
```

Why this is a false positive: `semantic-classification.py` is a standalone analysis script — the model is loaded at module level (line 9), embeddings/clustering run top-to-bottom, and results are written at module level (lines 89–90); there is no `if __name__ == "__main__"` guard and no importable module API. The print at line 42 is the script's intended progress/output, not operational logging in a non-script module, so the rule condition is not satisfied.

Checklist evidence: the flagged line is a zero-indent, module-level `print(` in a file whose entire body is top-level executable code (no guard, no importable module API), so the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `5` — `BP-PY-46`

- Function context: `./scripts/findings/functions/5.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FlashySurf/semantic-classification.py:43:1`
- Checklist pattern: same construct as finding 4 — module-scope print in the same standalone script, adjacent to line 42; distinct line, so kept as its own finding.

Source excerpt (from `semantic-classification.py`):

```
print(f"Dataset size: {n_samples} questions")
print(f"Using {n_neighbors} neighbors, {n_components} components, {n_clusters} clusters")
```

Why this is a false positive: same reasoning as finding 4 — line 43 is a zero-indent module-level print in the standalone `semantic-classification.py` script (no `__main__` guard, no library API), reporting the script's run parameters, so the rule condition ("print in non-script modules") is not met.

Checklist evidence: the flagged line is a module-level `print(` (column 0) in a top-to-bottom script with no `__main__` guard, so the "non-script modules" condition of the rule is not satisfied.

## True positives

### BP-PY-7 — `open` Without Context Manager

`semantic-classification.py:10` calls `open("questions.json")` without a `with` statement and passes the handle to `json.load`, which does not close it — the file object is never closed by any code path (only CPython refcounting frees it), which is exactly the resource-leak condition the rule flags. No detector exemption (import/comment/`with open`) applies.

| Finding ID | Source | Reason |
| --- | --- | --- |
| 3 | `semantic-classification.py:10:23` | `questions = json.load(open("questions.json"))` — `open` outside a `with` statement, handle never explicitly closed, satisfying the rule condition. |

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `./scripts/chunks/Chunk_1_5.txt`
- Function evidence: `./scripts/findings/functions/1.txt`–`5.txt`
- Validation: `git diff --check` — pass
