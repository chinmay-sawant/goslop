# False-positive audit — sync-with-uv

## Run metadata

```yaml
timestamp: 2026-08-02T07:22:28Z
repository: sync-with-uv
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/sync-with-uv
branch: main
commit: 27f1d39
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/sync-with-uv
chunk_path: scripts/sync-with-uv/chunks
function_context_path: scripts/sync-with-uv/findings/functions
```

## Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop` (goslop binary used: `./bin/goslop`)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/sync-with-uv/chunks -context-dir scripts/sync-with-uv/findings/functions real-repos/sync-with-uv`
- Findings: `20`
- Chunks reviewed: `scripts/sync-with-uv/chunks/Chunk_1_20.txt`
- Function contexts reviewed: `scripts/sync-with-uv/findings/functions/1.txt` .. `scripts/sync-with-uv/findings/functions/20.txt`

## Audit checklist

- [x] Read every assigned chunk under `scripts/sync-with-uv/chunks`.
- [x] Read `scripts/sync-with-uv/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 15 | 1, 2, 3, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 20 |
| True positive | 5 | 4, 16, 17, 18, 19 |
| Uncertain | 0 | — |

## False positives

### [ ] Finding `1` — `BP-PY-46`

- Function context: `scripts/sync-with-uv/findings/functions/1.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/sync-with-uv/scripts/gen_ref_pages.py:30:9`
- Checklist pattern: the flagged `print(` sits at module scope of a standalone docs-generation script whose whole body is top-level executable code (no `__main__` guard, no importable module API), and it writes generated content to a file handle — the "print used in non-script modules" condition is not satisfied.

Source excerpt:

```
    with mkdocs_gen_files.open(full_doc_path, "w") as fd:
        identifier = ".".join(parts)
        print("::: " + identifier, file=fd)
```

Why this is a false positive: `gen_ref_pages.py` is a script, not a library module — the file loop, the `mkdocs_gen_files.open(...)` writes and the `nav.build_literate_nav()` all execute top-to-bottom at import time. The flagged print writes page content into the file object opened by the caller (`file=fd`), i.e. it is the script's file-output mechanism, not operational logging in a non-script module, so the rule's stated condition (`print` is used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged line is a module-level `print(` with `file=` targeting a caller-opened file object in a top-to-bottom script with no library API; the "non-script modules" condition of the rule is not met.

### [ ] Finding `2` — `BP-PY-46`

- Function context: `scripts/sync-with-uv/findings/functions/2.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/sync-with-uv/src/sync_with_uv/cli.py:132:9`
- Checklist pattern: the flagged `print(` is CLI presentation output (user-facing error text to stderr) in the cyclopts CLI command module, not operational logging in a non-script module.

Source excerpt:

```
    except ValueError as e:
        print("Error:", e, file=sys.stderr)
        return 1
```

Why this is a false positive: `cli.py` is the CLI entry module — `App(name="sync-with-uv")` (cli.py:15) with the command registered via `@app.default()` (cli.py:80) and invoked only from `__main__.py` (`from .cli import app; app()`); nothing imports it as a library, so the prints only fire in CLI context. This print reports a config-resolution error to the user, which is exactly the CLI presentation category the rule's own fix text allows ("keep print under `if __name__ == "__main__"` for CLIs") and that the detector already exempts for argparse `main()` / `print_*` / `cmd_*` CLI presentation. The "print in library code" condition is not satisfied.

Checklist evidence: the print is user-facing error output in the cyclopts-registered CLI command body (`@app.default()`), the module has no library consumer, so the "print used for operational logging in non-script modules" condition is unmet.

### [ ] Finding `3` — `BP-PY-1`

- Function context: `scripts/sync-with-uv/findings/functions/3.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/sync-with-uv/src/sync_with_uv/cli.py:161:1`
- Checklist pattern: the rule condition flags a broad `except Exception` that "hides failures"; this handler surfaces the failure instead of swallowing it.

Source excerpt:

```
    except Exception as e:  # noqa: BLE001
        print("Error:", e, file=sys.stderr)
        return 123
```

Why this is a false positive: the handler deliberately reports the failure — it prints the exception detail to stderr and returns the documented internal-error exit code `123` (docstring: "Return code 123 means there was an internal error"), so the failure is surfaced to the user and signaled via a non-zero status, not hidden. The "hides failures" clause of the rule condition is not satisfied; the only way to reach this handler is through the CLI entry path.

Checklist evidence: the except suite prints the exception (`print("Error:", e, file=sys.stderr)`) and returns a non-zero status, i.e. the exception is handled by explicit reporting — the "without handling ... hides failures" clause of the rule condition is not met.

### [ ] Finding `5` — `BP-PY-46`

- Function context: `scripts/sync-with-uv/findings/functions/5.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/sync-with-uv/src/sync_with_uv/cli.py:162:9`
- Checklist pattern: same construct as finding 2 — user-facing error output in the cyclopts CLI command body; distinct line, so kept as its own finding.

Source excerpt:

```
    except Exception as e:  # noqa: BLE001
        print("Error:", e, file=sys.stderr)
        return 123
```

Why this is a false positive: same reasoning as finding 2 — the print writes the error message to stderr in the CLI entry module registered with `@app.default()`, which no code imports as a library; it is CLI presentation, not operational logging in library code, so the rule condition is not satisfied.

Checklist evidence: user-facing stderr output inside the cyclopts CLI command body, module has no library consumer — "print in non-script modules / library code" condition unmet.

### [ ] Finding `6` — `BP-PY-46`

- Function context: `scripts/sync-with-uv/findings/functions/6.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/sync-with-uv/src/sync_with_uv/cli.py:169:13`
- Checklist pattern: `print(` in `_print_changes`, a CLI presentation helper producing the `--verbose` status report to stderr — not operational logging in a non-script module.

Source excerpt:

```
    for package, change in changes.repos.items():
        if isinstance(change, tuple):
            print(f"{package}: {change[0]} -> {change[1]}", file=sys.stderr)
        elif change:
```

Why this is a false positive: `_print_changes` (cli.py:166) is the CLI's verbose status report, invoked only from the cyclopts command body (`if verbose: _print_changes(changes)`, cli.py:148-149) and writing to stderr for the user. It is CLI presentation in the CLI entry module with no library consumer, so the rule condition ("print for operational logging in library code") is not satisfied.

Checklist evidence: the print is user-facing status output from a CLI presentation helper called only from the `@app.default()` command; "print in library code" condition unmet.

### [ ] Finding `7` — `BP-PY-46`

- Function context: `scripts/sync-with-uv/findings/functions/7.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/sync-with-uv/src/sync_with_uv/cli.py:171:13`
- Checklist pattern: same construct as finding 6 — CLI verbose status output in `_print_changes`; distinct line, so kept as its own finding.

Source excerpt:

```
            print(f"{package}: {change[0]} -> {change[1]}", file=sys.stderr)
        elif change:
            print(f"{package}: unchanged", file=sys.stderr)
        else:
```

Why this is a false positive: same reasoning as finding 6 — user-facing verbose status text in the CLI presentation helper of the cyclopts entry module; not operational logging in library code.

Checklist evidence: user-facing stderr status output from a CLI presentation helper called only from the CLI command; "print in library code" condition unmet.

### [ ] Finding `8` — `BP-PY-46`

- Function context: `scripts/sync-with-uv/findings/functions/8.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/sync-with-uv/src/sync_with_uv/cli.py:173:13`
- Checklist pattern: same construct as finding 6 — CLI verbose status output in `_print_changes`; distinct line, so kept as its own finding.

Source excerpt:

```
            print(f"{package}: unchanged", file=sys.stderr)
        else:
            print(f"{package}: not managed in uv", file=sys.stderr)
```

Why this is a false positive: same reasoning as finding 6 — user-facing verbose status text in the CLI presentation helper of the cyclopts entry module; not operational logging in library code.

Checklist evidence: user-facing stderr status output from a CLI presentation helper called only from the CLI command; "print in library code" condition unmet.

### [ ] Finding `9` — `BP-PY-46`

- Function context: `scripts/sync-with-uv/findings/functions/9.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/sync-with-uv/src/sync_with_uv/cli.py:177:13`
- Checklist pattern: same construct as finding 6 — CLI verbose status output in `_print_changes`; distinct line, so kept as its own finding.

Source excerpt:

```
        if dep.changed:
            old_spec = dep.old_spec or "(unpinned)"
            print(
                f"line {line_number}: {dep.package} {old_spec} -> {dep.new_spec}",
```

Why this is a false positive: same reasoning as finding 6 — the print reports per-dependency changes to the user under `--verbose`; CLI presentation output, not operational logging in library code.

Checklist evidence: user-facing stderr status output from a CLI presentation helper called only from the CLI command; "print in library code" condition unmet.

### [ ] Finding `10` — `BP-PY-46`

- Function context: `scripts/sync-with-uv/findings/functions/10.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/sync-with-uv/src/sync_with_uv/cli.py:182:13`
- Checklist pattern: same construct as finding 6 — CLI verbose status output in `_print_changes`; distinct line, so kept as its own finding.

Source excerpt:

```
        else:
            print(f"line {line_number}: {dep.package} unchanged", file=sys.stderr)
```

Why this is a false positive: same reasoning as finding 6 — user-facing verbose status text in the CLI presentation helper of the cyclopts entry module; not operational logging in library code.

Checklist evidence: user-facing stderr status output from a CLI presentation helper called only from the CLI command; "print in library code" condition unmet.

### [ ] Finding `11` — `BP-PY-46`

- Function context: `scripts/sync-with-uv/findings/functions/11.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/sync-with-uv/src/sync_with_uv/cli.py:183:5`
- Checklist pattern: same construct as finding 6 — CLI verbose status output in `_print_changes`; distinct line, so kept as its own finding.

Source excerpt:

```
            print(f"line {line_number}: {dep.package} unchanged", file=sys.stderr)
    print(file=sys.stderr)
```

Why this is a false positive: the trailing `print(file=sys.stderr)` emits the blank separator that closes the CLI's verbose report — presentation formatting, not operational logging in library code.

Checklist evidence: user-facing stderr formatting output in the CLI presentation helper of the cyclopts entry module; "print in library code" condition unmet.

### [ ] Finding `12` — `BP-PY-46`

- Function context: `scripts/sync-with-uv/findings/functions/12.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/sync-with-uv/src/sync_with_uv/cli.py:199:5`
- Checklist pattern: `print(` in `_print_diff`, the CLI's `--diff` presentation that writes the unified diff to stdout — the primary user output of the diff mode.

Source excerpt:

```
    if color:
        diff_lines = get_colored_diff(diff_lines)
    print("\n".join(diff_lines))
```

Why this is a false positive: `_print_diff` (cli.py:186) is invoked from the `@app.default()` command only when `--diff` is given (`if diff: _print_diff(...)`, cli.py:151-152) and emits the diff the user asked for to stdout. It is the CLI's intended output, not operational logging in library code.

Checklist evidence: the print emits the requested `--diff` output to stdout from a CLI presentation helper called only from the CLI command; "print in library code" condition unmet.

### [ ] Finding `13` — `BP-PY-46`

- Function context: `scripts/sync-with-uv/findings/functions/13.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/sync-with-uv/src/sync_with_uv/cli.py:208:5`
- Checklist pattern: `print(` in `_print_summary`, the CLI's completion summary to stderr — not operational logging in a non-script module.

Source excerpt:

```
def _print_summary(changes: Changes, *, dry_mode: bool) -> None:
    print("All done!", file=sys.stderr)
    would_be = "would be " if dry_mode else ""
```

Why this is a false positive: `_print_summary` (cli.py:207) is the CLI's user-facing completion summary, called from the command body (cli.py:157-158) and writing to stderr; it is presentation output of the cyclopts CLI entry, not operational logging in library code.

Checklist evidence: user-facing stderr summary output in the CLI entry module with no library consumer; "print in library code" condition unmet.

### [ ] Finding `14` — `BP-PY-46`

- Function context: `scripts/sync-with-uv/findings/functions/14.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/sync-with-uv/src/sync_with_uv/cli.py:212:5`
- Checklist pattern: same construct as finding 13 — CLI completion summary in `_print_summary`; distinct line, so kept as its own finding.

Source excerpt:

```
    n_pkg_changed = sum(isinstance(c, tuple) for c in changes.repos.values())
    n_pkg_unchanged = len(changes.repos) - n_pkg_changed
    print(
        f"{n_pkg_changed} {_plural(n_pkg_changed, 'package', 'packages')} "
```

Why this is a false positive: same reasoning as finding 13 — the print reports changed/unchanged package counts to the user; CLI presentation output, not operational logging in library code.

Checklist evidence: user-facing stderr summary output from the CLI command path; "print in library code" condition unmet.

### [ ] Finding `15` — `BP-PY-46`

- Function context: `scripts/sync-with-uv/findings/functions/15.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/sync-with-uv/src/sync_with_uv/cli.py:222:9`
- Checklist pattern: same construct as finding 13 — CLI completion summary in `_print_summary`; distinct line, so kept as its own finding.

Source excerpt:

```
        n_changed = sum(d.changed for d in changes.lines.values())
        n_unchanged = len(changes.lines) - n_changed
        print(
            f"{n_changed} {_plural(n_changed, 'dependency', 'dependencies')} "
```

Why this is a false positive: same reasoning as finding 13 — the print reports changed/unchanged dependency counts to the user; CLI presentation output, not operational logging in library code.

Checklist evidence: user-facing stderr summary output from the CLI command path; "print in library code" condition unmet.

### [ ] Finding `20` — `CWE-252`

- Function context: `scripts/sync-with-uv/findings/functions/20.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/sync-with-uv/tests/test_precommit.py:54:9`
- Checklist pattern: rule condition not satisfied — the call passes `check=True`, which is the rule's explicit exemption for `subprocess.run`.

Source excerpt:

```
    try:
        subprocess.run(
            ["pre-commit", "install"], cwd=repo_dir, check=True  # noqa: S607
        )
    except subprocess.CalledProcessError as e:
        # pre-commit install might crash on prerelease Python versions (SIGSEGV)
        pytest.skip(f"pre-commit install failed: {e}")
```

Why this is a false positive: the rule's condition (`detectCWE252` in `internal/lang/python/detectors/cwe/rules_platform.go`) states that "`subprocess.run(..., check=True)` ha[s] explicit success handling paths and [is] not reported" — and the shown call passes `check=True`, so the "process call return status is discarded without checking success" condition is not satisfied; a non-zero status raises `CalledProcessError`, which the surrounding `try/except` handles deliberately.

Checklist evidence: the call includes `check=True`, the rule's documented exemption, so the return status is not discarded — the "unchecked return value" condition is unmet.

## True positives

### CWE-396 — Declaration of Catch for Generic Exception

| Finding ID | Source | Reason |
| --- | --- | --- |
| 4 | `src/sync_with_uv/cli.py:161:1` | `except Exception as e:` matches `pyGenericExceptRE` verbatim; `detectCWE396` has no handler-inspection or re-raise exemption and `cli.py` is not a test module. |

### BP-PY-2 — Except Pass

| Finding ID | Source | Reason |
| --- | --- | --- |
| 16 | `src/sync_with_uv/repo_data.py:106:1` | `except KeyError:` whose immediate suite is solely `pass` — `detectBPPY2` requires exactly `len(suite) == 1 && suite[0] == "pass"`, satisfied on this line. |
| 19 | `src/sync_with_uv/repo_data.py:110:1` | Same construct on a distinct line: `except KeyError: pass` with a suite consisting only of `pass`. |

### CWE-390 — Detection of Error Condition Without Action

| Finding ID | Source | Reason |
| --- | --- | --- |
| 17 | `src/sync_with_uv/repo_data.py:106:1` | `exceptPassStart` matches the except clause whose first direct body statement is `pass`; the rule deliberately does not judge whether subsequent fallback logic is adequate handling. |

### CWE-1071 — Empty Code Block

| Finding ID | Source | Reason |
| --- | --- | --- |
| 18 | `src/sync_with_uv/repo_data.py:106:5` | The exception handler contains only `pass` — the empty-code-block trigger matches the shown source. |

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/sync-with-uv/chunks/Chunk_1_20.txt`
- Function evidence: `scripts/sync-with-uv/findings/functions/1.txt` .. `20.txt`
- Validation: `git diff --check` — pass

## Post-fix remaining-FP audit (2026-08-02)

## Run metadata (fresh scan)

```yaml
timestamp: 2026-08-02T16:38:00Z
repository: sync-with-uv
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/sync-with-uv
branch: main
commit: 27f1d39
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/sync-with-uv
chunk_path: scripts/sync-with-uv/chunks
function_context_path: scripts/sync-with-uv/findings/functions
```

## Scan evidence (fresh scan)

- Build command: `go build -o bin/goslop ./cmd/goslop` (goslop binary: `./bin/goslop`, rebuilt 2026-08-02 16:29)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/sync-with-uv/chunks -context-dir scripts/sync-with-uv/findings/functions real-repos/sync-with-uv`
- Findings: `7`
- Chunks reviewed: `scripts/sync-with-uv/chunks/Chunk_1_7.txt`
- Function contexts reviewed: `scripts/sync-with-uv/findings/functions/1.txt` .. `scripts/sync-with-uv/findings/functions/7.txt`

## Audit checklist

- [x] Read every assigned chunk under `scripts/sync-with-uv/chunks`.
- [x] Read `scripts/sync-with-uv/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary (fresh run)

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 2 | 1, 7 |
| True positive | 5 | 2, 3, 4, 5, 6 |
| Uncertain | 0 | — |

Matching: fresh 2 → audited TP 4 (CWE-396, `cli.py:161:1`); fresh 3 → audited TP 16 (BP-PY-2, `repo_data.py:106:1`); fresh 4 → audited TP 17 (CWE-390, `repo_data.py:106:1`); fresh 5 → audited TP 18 (CWE-1071, `repo_data.py:106:5`); fresh 6 → audited TP 19 (BP-PY-2, `repo_data.py:110:1`). Fresh 1 → audited FP 3 (BP-PY-1, `cli.py:161:1`); fresh 7 → audited FP 20 (CWE-252, `test_precommit.py:54:9`).

## False positives

### [ ] Finding `1` — `BP-PY-1`

- Function context: `scripts/sync-with-uv/findings/functions/1.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/sync-with-uv/src/sync_with_uv/cli.py:161:1`
- Checklist pattern: re-appearing audited FP (old finding 3, same source) — the broad `except Exception` handler surfaces the failure instead of hiding it, so the rule condition ("hides failures") is not satisfied.

Source excerpt:

```
        # return 1 if check and changed
        return int(check and fixed_text != config_text)
    except Exception as e:  # noqa: BLE001
        print("Error:", e, file=sys.stderr)
        return 123
```

Why this is a false positive: identical to the audited decision for old finding 3 — the handler reports the exception detail to stderr and returns the documented internal-error exit code `123` (docstring: "Return code 123 means there was an internal error"), so the failure is surfaced to the user and signaled via a non-zero status, not hidden. The "hides failures" clause of the rule condition is not satisfied; the only path to this handler is the cyclopts CLI entry (`@app.default()`), reachable only from `__main__.py`.

Checklist evidence: the except suite prints the exception (`print("Error:", e, file=sys.stderr)`) and returns a non-zero status, i.e. the exception is handled by explicit reporting — the "without handling ... hides failures" clause of the rule condition is not met.

### [ ] Finding `7` — `CWE-252`

- Function context: `scripts/sync-with-uv/findings/functions/7.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/sync-with-uv/tests/test_precommit.py:54:9`
- Checklist pattern: re-appearing audited FP (old finding 20, same source) — the call passes `check=True`, which is the rule's explicit exemption for `subprocess.run`.

Source excerpt:

```
    try:
        subprocess.run(
            ["pre-commit", "install"], cwd=repo_dir, check=True  # noqa: S607
        )
    except subprocess.CalledProcessError as e:
        # pre-commit install might crash on prerelease Python versions (SIGSEGV)
        pytest.skip(f"pre-commit install failed: {e}")
```

Why this is a false positive: identical to the audited decision for old finding 20 — the rule's condition ("process call return status is discarded without checking success") is not met because the call passes `check=True`; the rule's own fix text directs callers to "use check=True so a failed command cannot silently continue execution", and a non-zero status raises `CalledProcessError`, which the surrounding `try/except` handles deliberately by skipping the test.

Checklist evidence: the call includes `check=True`, the rule's documented exemption, so the return status is not discarded — the "unchecked return value" condition is unmet.

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/sync-with-uv/chunks/Chunk_1_7.txt`
- Function evidence: `scripts/sync-with-uv/findings/functions/1.txt` .. `7.txt`
- Validation: `git diff --check` — pass

## Post-fix v2 audit (latest binary)

Run metadata (build: `make build` ~17:56), scan evidence, classification summary (fresh counts):

- Chunk: `scripts/sync-with-uv/chunks/Chunk_1_8.txt`; contexts `functions/1.txt`..`8.txt`; findings: `8`
- Matching by `Source:`: fresh 1 → audited FP 1 (BP-PY-46, `gen_ref_pages.py:30:9`); fresh 2 → audited FP 3 / Mode-B FP 1 (BP-PY-1, `cli.py:161:1`); fresh 3 → audited TP 4 / Mode-B TP 2 (CWE-396, `cli.py:161:1`); fresh 4 → audited TP 16 / Mode-B TP 3 (BP-PY-2, `repo_data.py:106:1`); fresh 5 → audited TP 17 / Mode-B TP 4 (CWE-390, `repo_data.py:106:1`); fresh 6 → audited TP 18 / Mode-B TP 5 (CWE-1071, `repo_data.py:106:5`); fresh 7 → audited TP 19 / Mode-B TP 6 (BP-PY-2, `repo_data.py:110:1`); fresh 8 → audited FP 20 / Mode-B FP 7 (CWE-252, `test_precommit.py:54:9`).
- Fresh classification: **False positive 3** (1, 2, 8) / **True positive 5** (3, 4, 5, 6, 7) / **Uncertain 0**. New findings: **0** (fresh 1 is a re-appeared audited FP — absent from the Mode B run, back under the latest binary).

## Fix checklist (FP patterns)

| Pattern # | Rule | Trigger shape | Count | Example sources |
| --- | --- | --- | ---: | --- |
| 1 | BP-PY-46 | Module-scope `print("::: " + identifier, file=fd)` in a standalone top-to-bottom docs script (no `__main__` guard, no importable module API) writing generated content into a caller-opened file handle | 1 | `scripts/gen_ref_pages.py:30:9` |
| 2 | BP-PY-1 | `except Exception as e:` whose suite surfaces the failure — `print("Error:", e, file=sys.stderr)` then `return 123` (non-zero); "hides failures" clause unmet | 1 | `src/sync_with_uv/cli.py:161:1` |
| 3 | CWE-252 | `subprocess.run([...], check=True)` — the rule's own documented exemption ("use check=True"); status raises `CalledProcessError` handled by `try/except` | 1 | `tests/test_precommit.py:54:9` |

## New findings (any fresh finding with no prior classification; classify each)

None — every fresh finding matches an already-audited Source.

Validation: `git diff --check` — see run below.
