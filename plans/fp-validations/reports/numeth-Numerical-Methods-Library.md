# False-positive audit — numeth-Numerical-Methods-Library

## Run metadata

```yaml
timestamp: 2026-08-02T07:13:32Z
repository: numeth-Numerical-Methods-Library
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/numeth-Numerical-Methods-Library
branch: main
commit: 14f2ae3df5201465d319c2874cfac9462e30c20b
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/numeth-Numerical-Methods-Library
chunk_path: scripts/numeth-Numerical-Methods-Library/chunks
function_context_path: scripts/numeth-Numerical-Methods-Library/findings/functions
```

## Scan evidence

- Build command: `make build` (bin/goslop prebuilt)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/numeth-Numerical-Methods-Library/chunks -context-dir scripts/numeth-Numerical-Methods-Library/findings/functions real-repos/numeth-Numerical-Methods-Library`
- Findings: `5`
- Chunks reviewed: `scripts/numeth-Numerical-Methods-Library/chunks/Chunk_1_5.txt`
- Function contexts reviewed: `scripts/numeth-Numerical-Methods-Library/findings/functions/1.txt`, `2.txt`, `3.txt`, `4.txt`, `5.txt`

## Audit checklist

- [x] Read every assigned chunk under `scripts/numeth-Numerical-Methods-Library/chunks`.
- [x] Read `scripts/numeth-Numerical-Methods-Library/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 1 | 1 |
| True positive | 4 | 2, 3, 4, 5 |
| Uncertain | 0 | — |

## False positives

### [ ] Finding `1` — `PERF-PY-25`

- Function context: `scripts/numeth-Numerical-Methods-Library/findings/functions/1.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/numeth-Numerical-Methods-Library/src/numeth/interpolation.py:37:1`
- Checklist pattern: lambda is constructed at most once per function invocation, not once per loop element — the enclosing loop returns unconditionally immediately after the construction.

Source excerpt (from `src/numeth/interpolation.py`):

```
    for i in range(len(points) - 1):
        if xs[i] <= x <= xs[i + 1]:
            x1, y1 = points[i]
            x2, y2 = points[i + 1]
            result = y1 + (y2 - y1) * (x - x1) / (x2 - x1)
            f_interp = lambda xi: linear_interpolation(points, xi)
            return NumericalResult(result, method_info={'type': 'interpolation', 'method': 'linear_interpolation', 'points': points, 'f_interp': f_interp})
    raise ValueError("Interpolation failed")  # Should not reach here
```

Why this is a false positive: the rule flags a lambda constructed per homogeneous loop element, but here the lambda at line 37 is immediately followed by an unconditional `return` at line 38 — the loop is a linear search over intervals and executes the construct-and-return branch at most once per function call, so no lambda is built per element.

Checklist evidence: the line flagged is inside a loop body but the loop always returns on the first matching iteration; the lambda is constructed zero or one times per call, never `len(points)` times, so the "per homogeneous loop element" condition of the rule is not satisfied.

## True positives

### BP-PY-46 — `print` Debugging In Library Code

`results.py` is a library module (`src/numeth/results.py` of the `numeth` package) with no `if __name__ == "__main__"` guard, no argparse CLI, and no CLI decorators; every flagged `print(` is a real call in a method body, which satisfies the rule condition (print in non-script module, outside main guard, non-test).

| Finding ID | Source | Reason |
| --- | --- | --- |
| 2 | `src/numeth/results.py:17` | `print("Matplotlib is required for visualization.")` in `NumericalResult.graph()` — user-facing print in library method. |
| 3 | `src/numeth/results.py:28` | `print(f"Visualization not implemented for type: {m_type}")` in `NumericalResult.graph()` else-branch. |
| 4 | `src/numeth/results.py:102` | `print("Matplotlib is required for visualization.")` in `IterativeResult.graph()` — second class, distinct line. |
| 5 | `src/numeth/results.py:111` | `print(f"Visualization not implemented for type: {m_type}")` in `IterativeResult.graph()` else-branch. |

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/numeth-Numerical-Methods-Library/chunks/Chunk_1_5.txt`
- Function evidence: `scripts/numeth-Numerical-Methods-Library/findings/functions/1.txt`–`5.txt`
- Validation: `git diff --check` — pass
