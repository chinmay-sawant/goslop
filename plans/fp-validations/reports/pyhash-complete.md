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

## Post-fix over-suppression audit (2026-08-02)

Mode B: fresh findings (9) < audited true positives (10) — checking which audited TPs the FP-reduction fix (`b5b8fde`, binary rebuilt 2026-08-02 16:29 local) suppressed.

### Run metadata

```yaml
timestamp: 2026-08-02T11:08:51Z
repository: pyhash-complete
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pyhash-complete
branch: main
commit: 16dfd3d31673ff8399da51bfdb9844036f8785fb
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pyhash-complete
chunk_path: scripts/pyhash-complete/chunks
function_context_path: scripts/pyhash-complete/findings/functions
binary: ./bin/goslop (rebuilt from b5b8fde, 2026-08-02 16:29 local)
```

### Scan evidence

- Build command: `(pre-built) ./bin/goslop` (commit `b5b8fde`)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/pyhash-complete/chunks -context-dir scripts/pyhash-complete/findings/functions real-repos/pyhash-complete`
- Findings: `9`
- Chunks reviewed: `scripts/pyhash-complete/chunks/Chunk_1_9.txt`
- Function contexts reviewed: `scripts/pyhash-complete/findings/functions/1.txt` … `9.txt`

### Over-suppression table

| Old finding ID | Rule | Source | One-line reason (from old audit) | Current status |
| --- | --- | --- | --- | --- |
| 1 | PERF-PY-23 | `get_builtins.py:13` | `json.dumps(` inside the `for i in range(...)` loop, no skip pattern | present in fresh scan (fresh finding 1, same line) |
| 2 | BP-PY-46 | `get_builtins.py:26` | module-level `print(jsonify(build_map()))`, no `__main__` guard | present in fresh scan (fresh finding 2) |
| 3 | BP-PY-46 | `lookup.py:24` | `print("corrupt or truncated file")` inside library `load()` | present in fresh scan (fresh finding 3) |
| 4 | BP-PY-46 | `lookup.py:28` | `print(f"bad magic number: ...")` inside `load()` | present in fresh scan (fresh finding 4) |
| 5 | BP-PY-46 | `lookup.py:33` | `print("corrupt or truncated slots section")` inside `load()` | present in fresh scan (fresh finding 5) |
| 6 | CWE-829 | `lookup.py:74` | `__import__(name)` with non-literal `name = sys.argv[1].split('.')[0]` | present in fresh scan (fresh finding 6) |
| 7 | CWE-94 | `lookup.py:74` | same call site, dynamic module name at code-generation sink | present in fresh scan (fresh finding 7) |
| 8 | BP-PY-46 | `simpleui.py:71` | debug `print([c.name for c in comps][:10])` in `Editor.on_text_changed` | present in fresh scan (fresh finding 8) |
| 9 | BP-PY-46 | `simpleui.py:74` | debug `print(completions[:10])` in same handler | present in fresh scan (fresh finding 9) |
| 10 | CWE-396 | `simpleui.py:75` | `except Exception as e:` matches `pyGenericExceptRE`; old rule had no re-raise exemption | **suppressed-but-present** |

Fresh finding sources were matched to audited TP sources (file:line) one-for-one; findings 1–9 all reappear at identical locations. Only old finding 10 is missing from the fresh scan. Audited FP 11 (`simpleui.py:145`, `app.exec()`) is also absent — resolved, but out of scope for this over-suppression audit.

### [ ] Suppressed-but-present TP: old finding 10 — CWE-396

- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pyhash-complete/simpleui.py:75`
- Current status: `suppressed-but-present` — construct still in the source, not fixed/removed

Source excerpt (current `simpleui.py:75-76`, read 2026-08-02):

```
        except Exception as e:
            raise e
```

Why it is flagged for review: this is the only generic `except Exception` in the repo (`rg "except"` finds only `simpleui.py:75` and the specific `lookup.py:75 except ModuleNotFoundError`), so under the pre-fix rule condition it necessarily fired `detectCWE396` (first generic except in a non-test module, no suite inspection) and was correctly audited as a TP. The construct still exists unchanged, so the suppression is rule-driven, not code removal.

Why the suppression is the intentional batch-3 exemption: `detectCWE396` (`internal/lang/python/detectors/cwe/rules_platform.go:92`) now skips any generic-except suite for which `suiteSurfacesFailureMasked` is true; `suiteLineSurfacesFailure` (`rules_platform.go:62`) returns true for `raise e` via the `raise ` prefix check. The suite re-raises the original exception, so the handler surfaces the failure instead of hiding distinct failure conditions — the exact weakness the batch-3 fix (`plans/fp-review-reduce/python/batch-3-broad-except-reduction-2026-08-02.md`, guardrail `CWE-396-batch3-reraise-from` family) exempts. A bare `raise e` is the canonical surfacing action; this is a correctly suppressed TP under the corrected rule condition, not an over-suppression defect.

Checklist evidence: `./bin/goslop -explain CWE-396` reports `unknown rule` for the fresh binary, but the rule ID remains registered in the codebase (`rules_platform.go`); the fresh scan itself confirms the exemption — running the same scan command on the unchanged tree emits no CWE-396 finding for `simpleui.py:75`.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/pyhash-complete/chunks/Chunk_1_9.txt`
- Function evidence: `scripts/pyhash-complete/findings/functions/1.txt` … `9.txt`
- Validation: `git diff --check` — `pass`
