# False-positive audit — FlashySurf

## Run metadata

```yaml
timestamp: 2026-08-02T07:14:56Z
repository: FlashySurf
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FlashySurf
branch: main
commit: 3d14f8ca328e7f11d3fc6f2fe86bc15a2463184e
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FlashySurf
chunk_path: scripts/FlashySurf/chunks
function_context_path: scripts/FlashySurf/findings/functions
```

## Scan evidence

- Build command: `make build` (bin/goslop prebuilt)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/FlashySurf/chunks -context-dir scripts/FlashySurf/findings/functions real-repos/FlashySurf`
- Findings: `5`
- Chunks reviewed: `scripts/FlashySurf/chunks/Chunk_1_5.txt`
- Function contexts reviewed: `scripts/FlashySurf/findings/functions/1.txt`, `2.txt`, `3.txt`, `4.txt`, `5.txt`

## Audit checklist

- [x] Read every assigned chunk under `scripts/FlashySurf/chunks`.
- [x] Read `scripts/FlashySurf/findings/functions/<finding-id>.txt` for every proposed false positive.
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

- Function context: `scripts/FlashySurf/findings/functions/1.txt`
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

- Function context: `scripts/FlashySurf/findings/functions/2.txt`
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

- Function context: `scripts/FlashySurf/findings/functions/4.txt`
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

- Function context: `scripts/FlashySurf/findings/functions/5.txt`
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
- Chunk evidence: `scripts/FlashySurf/chunks/Chunk_1_5.txt`
- Function evidence: `scripts/FlashySurf/findings/functions/1.txt`–`5.txt`
- Validation: `git diff --check` — pass

## Post-fix remaining-FP audit (2026-08-02)

### Run metadata

```yaml
timestamp: 2026-08-02T16:38:00Z
repository: FlashySurf
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FlashySurf
branch: main
commit: 3d14f8ca328e7f11d3fc6f2fe86bc15a2463184e
scanner_revision: b5b8fde (post-FP-fix binary, rebuilt 2026-08-02 16:29)
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FlashySurf
chunk_path: scripts/FlashySurf/chunks
function_context_path: scripts/FlashySurf/findings/functions
```

### Scan evidence

- Build command: `make build` (bin/goslop prebuilt, fix `b5b8fde`)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/FlashySurf/chunks -context-dir scripts/FlashySurf/findings/functions real-repos/FlashySurf`
- Findings: `5`
- Chunks reviewed: `scripts/FlashySurf/chunks/Chunk_1_5.txt`
- Function contexts reviewed: `scripts/FlashySurf/findings/functions/1.txt`, `2.txt`, `3.txt`, `4.txt`, `5.txt`

### Audit checklist

- [x] Read every assigned chunk under `scripts/FlashySurf/chunks`.
- [x] Read `scripts/FlashySurf/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

### Classification summary

Fresh findings matched to the audited run by `Source:` path (file:line:col); the fix did not change any finding for this repo.

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 4 | 1, 2, 4, 5 |
| True positive | 1 | 3 |
| Uncertain | 0 | — |

Remaining false positives after the fix: `4` (unchanged from the pre-fix audit — the fix did not suppress any finding in FlashySurf).

### False positives

#### [ ] Findings `1`, `2` — `BP-PY-46`

- Function contexts: `scripts/FlashySurf/findings/functions/1.txt`, `2.txt`
- Sources: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FlashySurf/data-process.py:167:1`, `data-process.py:168:1`
- Checklist pattern: adjacent module-level (column 0) `print(` calls in a standalone top-to-bottom script with no `if __name__ == "__main__"` guard; identical source construct, so grouped (both IDs listed).

Source excerpt (from `data-process.py`):

```
print(f"Math questions: {len(output['math'])}/{mathCount}")
print(f"English questions: {len(output['english'])}/{englishCount}")



with open("questions.json", "w+") as f:
    json.dump(output, f, indent=4)
```

Why this is a false positive: lines 167–168 are zero-indent module-level prints in the standalone `data-process.py` script, whose entire body (file I/O at line 22, the transformation loop, output writes at lines 172–176) executes at import time with no `__main__` guard and no importable module API — the prints are the script's intended summary output, so the rule condition ("print is used for operational logging in non-script modules") is not satisfied.

Checklist evidence: the flagged lines are module-level `print(` calls at column 0 in a script with no `__main__` guard and no library API, so the "non-script modules" condition of the rule is not satisfied.

#### [ ] Findings `4`, `5` — `BP-PY-46`

- Function contexts: `scripts/FlashySurf/findings/functions/4.txt`, `5.txt`
- Sources: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FlashySurf/semantic-classification.py:42:1`, `semantic-classification.py:43:1`
- Checklist pattern: adjacent module-level (column 0) `print(` calls in a standalone top-to-bottom analysis script with no `__main__` guard; identical source construct, so grouped (both IDs listed).

Source excerpt (from `semantic-classification.py`):

```
n_clusters = min(87, n_samples - 1)  # 87 clusters as requested

print(f"Dataset size: {n_samples} questions")
print(f"Using {n_neighbors} neighbors, {n_components} components, {n_clusters} clusters")
```

Why this is a false positive: lines 42–43 are zero-indent module-level prints in the standalone `semantic-classification.py` script — the model loads at module level (line 9), embeddings/clustering run top-to-bottom, and results are written at module level; there is no `__main__` guard and no importable module API. The prints report the script's run parameters, not operational logging in a non-script module, so the rule condition is not satisfied.

Checklist evidence: the flagged lines are module-level `print(` calls at column 0 in a script with no `__main__` guard and no library API, so the "non-script modules" condition of the rule is not satisfied.

### True positives

#### `BP-PY-7` — `open` Without Context Manager

| Finding ID | Source | Reason |
| --- | --- | --- |
| 3 | `semantic-classification.py:10:23` | `questions = json.load(open("questions.json"))` — source matches audited TP; `open` outside a `with` statement, handle never explicitly closed, satisfying the rule condition. |

### Uncertain findings

None.

### Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/FlashySurf/chunks/Chunk_1_5.txt`
- Function evidence: `scripts/FlashySurf/findings/functions/1.txt`–`5.txt`
- Rule conditions: `./bin/goslop -explain BP-PY-46 --config templates/goslop-python.toml` (non-script modules condition), `./bin/goslop -explain BP-PY-7 --config templates/goslop-python.toml` (open without `with`)
- Validation: `git diff --check` — pass

## Post-fix v2 audit (latest binary)

### Run metadata

```yaml
timestamp: 2026-08-02T18:05:00Z
repository: FlashySurf
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FlashySurf
branch: main
commit: 3d14f8ca328e7f11d3fc6f2fe86bc15a2463184e
scanner_revision: latest binary (make build ~17:56)
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FlashySurf
chunk_path: scripts/FlashySurf/chunks
function_context_path: scripts/FlashySurf/findings/functions
```

### Scan evidence

- Build command: `make build` (bin/goslop rebuilt ~2026-08-02 17:56)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/FlashySurf/chunks -context-dir scripts/FlashySurf/findings/functions real-repos/FlashySurf`
- Findings: `5`
- Chunks reviewed: `scripts/FlashySurf/chunks/Chunk_1_5.txt`
- Function contexts reviewed: `scripts/FlashySurf/findings/functions/1.txt`–`5.txt`

### Classification summary (fresh counts)

All 5 fresh findings match prior audited `Source:` paths exactly (verified against both source files — no line drift). Classifications reused from the original audit + Mode A/B appends.

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| True positive | 1 | 3 |
| False positive | 4 | 1, 2, 4, 5 |
| Uncertain | 0 | — |

- Fresh findings: `5`; fresh TP `1` / FP `4` / U `0`.
- New findings (no prior classification): none — the latest binary changed no finding for this repo (same 5 sources as the b5b8fde audit; fingerprints now carry `goslop:2:` but sources/rules are identical).

## Fix checklist (FP patterns)

| Pattern # | Rule | Trigger shape | Count | Example sources |
| --- | --- | --- | ---: | --- |
| 1 | BP-PY-46 | module-level (column 0) `print(` in a standalone top-to-bottom script with no `if __name__ == "__main__"` guard and no importable module API (whole body executes at import); rule targets "non-script modules" but detector lacks a script/library distinction — safe condition: print at module scope in a script file (no guard needed), vulnerable: print in an importable module with library API | 4 | `data-process.py:167:1`, `data-process.py:168:1`, `semantic-classification.py:42:1`, `semantic-classification.py:43:1` |

## New findings

None — every fresh finding has a prior classification; the latest binary produced no new findings for FlashySurf.
