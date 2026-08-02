# False-positive audit — pyhash-complete

## Run metadata

```yaml
timestamp: 2026-08-02T07:15:05Z
repository: pyhash-complete
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pyhash-complete
branch: main
commit: 16dfd3d31673ff8399da51bfdb9844036f8785fb
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pyhash-complete
chunk_path: scripts/pyhash-complete/chunks
function_context_path: scripts/pyhash-complete/findings/functions
```

## Scan evidence

- Build command: `(pre-built) ./bin/goslop`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/pyhash-complete/chunks -context-dir scripts/pyhash-complete/findings/functions real-repos/pyhash-complete`
- Findings: `11`
- Chunks reviewed: `scripts/pyhash-complete/chunks/Chunk_1_11.txt`
- Function contexts reviewed: `scripts/pyhash-complete/findings/functions/1.txt … 11.txt`

## Audit checklist

- [x] Read every assigned chunk under `scripts/pyhash-complete/chunks`.
- [x] Read `scripts/pyhash-complete/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 1 | 11 |
| True positive | 10 | 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 |
| Uncertain | 0 | — |

## False positives

### [ ] Finding 11 — BP-PY-12

- Function context: `scripts/pyhash-complete/findings/functions/11.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pyhash-complete/simpleui.py:145:18`
- Checklist pattern: `exec` is a Qt (PySide6 `QApplication`) event-loop method, not Python's `exec` builtin; the identifier match fires on attribute access.

Source excerpt:

```
    app = QApplication(sys.argv)
    window = Editor(args.dataset)
    window.show()
    sys.exit(app.exec())
```

Why this is a false positive: `app.exec()` is a fixed zero-argument method of the Qt `QApplication` object that runs the GUI event loop; it is not Python's `exec` builtin and evaluates no dynamic code.

Checklist evidence: `indexOfIdent` only rejects matches preceded by an ident byte (`isIdentByte`); the `.` before `exec` in `app.exec(` is not an ident byte, so the boundary check passes and the method name is misreported as the builtin. The rule condition "eval/exec on dynamic input" requires a builtin call taking code; here the call is an attribute method with no arguments.

## True positives

### PERF-PY-23

| Finding | Source | Reason |
| --- | --- | --- |
| 1 | `get_builtins.py:13` | `json.dumps(` matches `encodeInLoopRE` and the line is inside the `for i in range(...)` loop at line 11; no skip pattern (`encode_record/cell/leaf`), not a test file. |

### BP-PY-46

| Finding | Source | Reason |
| --- | --- | --- |
| 2 | `get_builtins.py:26` | Module-level `print(jsonify(build_map()))`; no `__main__` guard, no argparse, so no skip rule applies — the print runs on import. |
| 3 | `lookup.py:24` | `print` inside `load()` — a reusable library function (also imported by `simpleui.py`); not in `main`, so `pythonCLIPrintSkipFunc` does not skip it. |
| 4 | `lookup.py:28` | Same construct as 3, distinct line: `print(f"bad magic number: ...")` inside `load()`. |
| 5 | `lookup.py:33` | Same construct as 3, distinct line: `print("corrupt or truncated slots section")` inside `load()`. |
| 8 | `simpleui.py:71` | Debug `print([c.name for c in comps][:10])` in `Editor.on_text_changed`; `__main__` guard never invokes `main(` and no argparse `print_*`/`cmd_*` name matches, so no skip applies. |
| 9 | `simpleui.py:74` | Same handler, distinct line: `print(completions[:10])` — leftover debug output (neighbouring lines are commented-out prints). |

### CWE-829 / CWE-94 (same construct, lookup.py:74)

| Finding | Source | Reason |
| --- | --- | --- |
| 6 | `lookup.py:74` | `__import__(name)` where `name = sys.argv[1].split('.')[0]` — a non-literal, user-supplied expression reaches a dynamic-import sink; `isDynamicExpr` returns true. |
| 7 | `lookup.py:74` | Same call site: dynamic module name reaches the code-generation/dynamic-import sink per `detectCWE94` (non-literal first arg). |

### CWE-396

| Finding | Source | Reason |
| --- | --- | --- |
| 10 | `simpleui.py:75` | `except Exception as e:` matches `pyGenericExceptRE` exactly; rule spec reports the declaration of a generic catch and has no re-raise exemption. |

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/pyhash-complete/chunks/Chunk_1_11.txt`
- Function evidence: `scripts/pyhash-complete/findings/functions/1.txt` … `11.txt`
- Validation: `git diff --check` — `pass`
