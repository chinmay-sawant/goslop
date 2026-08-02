# False-positive audit — whatsapp-wrapped

## Run metadata

```yaml
timestamp: 2026-08-02T12:41:00Z
repository: whatsapp-wrapped
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/whatsapp-wrapped
branch: main
commit: 815689e
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/whatsapp-wrapped
chunk_path: scripts/whatsapp-wrapped/chunks
function_context_path: scripts/whatsapp-wrapped/findings/functions
```

## Scan evidence

- Build command: `not recorded (pre-existing ./bin/goslop, built 2026-08-02 02:40)`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/<REPO_NAME>/chunks -context-dir scripts/<REPO_NAME>/findings/functions real-repos/<REPO_NAME>`
- Findings: `40`
- Chunks reviewed: `scripts/whatsapp-wrapped/chunks/Chunk_1_25.txt`, `scripts/whatsapp-wrapped/chunks/Chunk_26_40.txt` (re-run of the scan reproduced identical chunks byte-for-byte)
- Function contexts reviewed: `scripts/whatsapp-wrapped/findings/functions/1.txt`–`40.txt`

## Audit checklist

- [x] Read every assigned chunk under `scripts/whatsapp-wrapped/chunks`.
- [x] Read `scripts/whatsapp-wrapped/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 17 | 1, 2, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 37 |
| True positive | 23 | 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 36, 38, 39, 40 |
| Uncertain | 0 | — |

## False positives

### [ ] Finding `1` — `BP-PY-41`

- Function context: `scripts/whatsapp-wrapped/findings/functions/1.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/whatsapp-wrapped/tests/test_parser.py:358:1`
- Checklist pattern: rule condition is `hasCall && !hasAssert` — the test body must contain no assertion for the rule to fire.

Source excerpt:

```
    df = parse_chat(chat_spanish)
    assert len(df) == 5
    assert df["message_type"][0] == "video"
    assert df["message_type"][1] == "sticker"
```

Why this is a false positive: the function body contains multiple `assert` statements (lines 372–377, 387–392, 402–406), so the `hasCall && !hasAssert` condition is not satisfied. The detector's body scan stops early at the column-0 content lines of the triple-quoted chat sample strings (e.g. line 366 `21-06-23, 11:20:04 - Cristhian: sticker omitido`) and never reaches the asserts.

Checklist evidence: `assert` statements exist in the same function body at body indent; the rule's "placeholder test without assertions" condition fails.

### [ ] Finding `2` — `BP-PY-41`

- Function context: `scripts/whatsapp-wrapped/findings/functions/2.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/whatsapp-wrapped/tests/test_parser.py:409:1`
- Checklist pattern: rule condition is `hasCall && !hasAssert` — the test body must contain no assertion for the rule to fire.

Source excerpt:

```
    df = parse_chat(chat_normal)
    assert len(df) == 4
    assert df["message_type"][0] == "text"
```

Why this is a false positive: the function body contains four `assert` statements (lines 420, 422–425), so the `hasCall && !hasAssert` condition is not satisfied. As with finding 1, the column-0 content lines of the triple-quoted sample strings (lines 414–417) abort the assertion scan early.

Checklist evidence: `assert` statements exist in the same function body at body indent; the rule's "placeholder test without assertions" condition fails.

### [ ] Finding `22` — `BP-PY-46`

- Function context: `scripts/whatsapp-wrapped/findings/functions/22.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/whatsapp-wrapped/whatsapp_wrapped/generator.py:524:9`
- Checklist pattern: rule flags `print(` in non-script modules but exempts prints inside `main()` when `if __name__ == "__main__":` invokes `main()` (`pythonCLIPrintSkipFunc`).

Source excerpt:

```
def main():
    """Command-line interface for the report generator."""
    ...
    if not chat_file.exists():
        print(f"[!] Error: Chat file not found: {chat_file}")
        sys.exit(1)

...
if __name__ == "__main__":
    main()
```

Why this is a false positive: the print is CLI user output inside `main()`, which is invoked from the `if __name__ == "__main__":` guard at lines 562–563. The rule condition exempts this exact pattern (rule fix text: "keep print under `if __name__ == "__main__"` for CLIs"). The finding fired because the content line `Examples:` (column 0) inside the argparse `epilog` triple-quoted string (lines 451–458) resets the detector's CLI-indent tracking, so the prints inside `main()` were scanned as library code.

Checklist evidence: the finding line sits in the `main()` body; `pythonMainGuardInvokesMain` returns true for this file, so `pythonCLIPrintSkipFunc("main", ...)` exempts the block — the rule's own condition is not met.

### [ ] Finding `23` — `BP-PY-46`

- Function context: `scripts/whatsapp-wrapped/findings/functions/23.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/whatsapp-wrapped/whatsapp_wrapped/generator.py:528:9`
- Checklist pattern: same as finding 22 — prints in `main()` invoked by the `__main__` guard are exempt.

Source excerpt:

```
    if not args.quiet:
        print()
        print("=" * 60)
```

Why this is a false positive: CLI banner output inside `main()`, invoked by the `__main__` guard; exempt per the rule's CLI-output exception. Same `epilog`-string indent-tracking artifact as finding 22.

Checklist evidence: the finding line is inside the `main()` body, which `pythonCLIPrintSkipFunc` exempts when `mainInvoked` is true.

### [ ] Finding `24` — `BP-PY-46`

- Function context: `scripts/whatsapp-wrapped/findings/functions/24.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/whatsapp-wrapped/whatsapp_wrapped/generator.py:529:9`
- Checklist pattern: same as finding 22.

Source excerpt:

```
    if not args.quiet:
        print()
        print("=" * 60)
```

Why this is a false positive: CLI banner output inside `main()`; exempt per rule condition, fired only due to the `epilog`-string indent-tracking artifact.

Checklist evidence: the finding line is inside the `main()` body; `mainInvoked` is true so the rule exempts it.

### [ ] Finding `25` — `BP-PY-46`

- Function context: `scripts/whatsapp-wrapped/findings/functions/25.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/whatsapp-wrapped/whatsapp_wrapped/generator.py:530:9`
- Checklist pattern: same as finding 22.

Source excerpt:

```
        print("=" * 60)
        print("  WHATSAPP WRAPPED - REPORT GENERATOR")
        print("=" * 60)
```

Why this is a false positive: CLI banner output inside `main()`; exempt per rule condition, fired only due to the `epilog`-string indent-tracking artifact.

Checklist evidence: the finding line is inside the `main()` body; `mainInvoked` is true so the rule exempts it.

### [ ] Finding `26` — `BP-PY-46`

- Function context: `scripts/whatsapp-wrapped/findings/functions/26.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/whatsapp-wrapped/whatsapp_wrapped/generator.py:531:9`
- Checklist pattern: same as finding 22.

Source excerpt:

```
        print("=" * 60)
        print("  WHATSAPP WRAPPED - REPORT GENERATOR")
        print("=" * 60)
```

Why this is a false positive: CLI banner output inside `main()`; exempt per rule condition, fired only due to the `epilog`-string indent-tracking artifact.

Checklist evidence: the finding line is inside the `main()` body; `mainInvoked` is true so the rule exempts it.

### [ ] Finding `27` — `BP-PY-46`

- Function context: `scripts/whatsapp-wrapped/findings/functions/27.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/whatsapp-wrapped/whatsapp_wrapped/generator.py:532:9`
- Checklist pattern: same as finding 22.

Source excerpt:

```
        print("  WHATSAPP WRAPPED - REPORT GENERATOR")
        print("=" * 60)
        print()
```

Why this is a false positive: CLI banner output inside `main()`; exempt per rule condition, fired only due to the `epilog`-string indent-tracking artifact.

Checklist evidence: the finding line is inside the `main()` body; `mainInvoked` is true so the rule exempts it.

### [ ] Finding `28` — `BP-PY-46`

- Function context: `scripts/whatsapp-wrapped/findings/functions/28.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/whatsapp-wrapped/whatsapp_wrapped/generator.py:547:13`
- Checklist pattern: same as finding 22.

Source excerpt:

```
        if not args.quiet:
            print()
            print("=" * 60)
```

Why this is a false positive: CLI completion banner inside `main()`; exempt per rule condition, fired only due to the `epilog`-string indent-tracking artifact.

Checklist evidence: the finding line is inside the `main()` body; `mainInvoked` is true so the rule exempts it.

### [ ] Finding `29` — `BP-PY-46`

- Function context: `scripts/whatsapp-wrapped/findings/functions/29.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/whatsapp-wrapped/whatsapp_wrapped/generator.py:548:13`
- Checklist pattern: same as finding 22.

Source excerpt:

```
        if not args.quiet:
            print()
            print("=" * 60)
            print("  GENERATION COMPLETE")
```

Why this is a false positive: CLI completion banner inside `main()`; exempt per rule condition, fired only due to the `epilog`-string indent-tracking artifact.

Checklist evidence: the finding line is inside the `main()` body; `mainInvoked` is true so the rule exempts it.

### [ ] Finding `30` — `BP-PY-46`

- Function context: `scripts/whatsapp-wrapped/findings/functions/30.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/whatsapp-wrapped/whatsapp_wrapped/generator.py:549:13`
- Checklist pattern: same as finding 22.

Source excerpt:

```
            print("=" * 60)
            print("  GENERATION COMPLETE")
            print("=" * 60)
```

Why this is a false positive: CLI completion banner inside `main()`; exempt per rule condition, fired only due to the `epilog`-string indent-tracking artifact.

Checklist evidence: the finding line is inside the `main()` body; `mainInvoked` is true so the rule exempts it.

### [ ] Finding `31` — `BP-PY-46`

- Function context: `scripts/whatsapp-wrapped/findings/functions/31.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/whatsapp-wrapped/whatsapp_wrapped/generator.py:550:13`
- Checklist pattern: same as finding 22.

Source excerpt:

```
            print("  GENERATION COMPLETE")
            print("=" * 60)
            print()
```

Why this is a false positive: CLI completion banner inside `main()`; exempt per rule condition, fired only due to the `epilog`-string indent-tracking artifact.

Checklist evidence: the finding line is inside the `main()` body; `mainInvoked` is true so the rule exempts it.

### [ ] Finding `32` — `BP-PY-46`

- Function context: `scripts/whatsapp-wrapped/findings/functions/32.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/whatsapp-wrapped/whatsapp_wrapped/generator.py:551:13`
- Checklist pattern: same as finding 22.

Source excerpt:

```
            print("=" * 60)
            print()
            print(f"  HTML: {html_path}")
```

Why this is a false positive: CLI completion output inside `main()`; exempt per rule condition, fired only due to the `epilog`-string indent-tracking artifact.

Checklist evidence: the finding line is inside the `main()` body; `mainInvoked` is true so the rule exempts it.

### [ ] Finding `33` — `BP-PY-46`

- Function context: `scripts/whatsapp-wrapped/findings/functions/33.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/whatsapp-wrapped/whatsapp_wrapped/generator.py:552:13`
- Checklist pattern: same as finding 22.

Source excerpt:

```
            print()
            print(f"  HTML: {html_path}")
            if static_path:
```

Why this is a false positive: CLI completion output inside `main()`; exempt per rule condition, fired only due to the `epilog`-string indent-tracking artifact.

Checklist evidence: the finding line is inside the `main()` body; `mainInvoked` is true so the rule exempts it.

### [ ] Finding `34` — `BP-PY-46`

- Function context: `scripts/whatsapp-wrapped/findings/functions/34.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/whatsapp-wrapped/whatsapp_wrapped/generator.py:554:17`
- Checklist pattern: same as finding 22.

Source excerpt:

```
            if static_path:
                print(f"  Static HTML: {static_path}")
            print()
```

Why this is a false positive: CLI completion output inside `main()`; exempt per rule condition, fired only due to the `epilog`-string indent-tracking artifact.

Checklist evidence: the finding line is inside the `main()` body; `mainInvoked` is true so the rule exempts it.

### [ ] Finding `35` — `BP-PY-46`

- Function context: `scripts/whatsapp-wrapped/findings/functions/35.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/whatsapp-wrapped/whatsapp_wrapped/generator.py:555:13`
- Checklist pattern: same as finding 22.

Source excerpt:

```
                print(f"  Static HTML: {static_path}")
            print()
```

Why this is a false positive: CLI completion output inside `main()`; exempt per rule condition, fired only due to the `epilog`-string indent-tracking artifact.

Checklist evidence: the finding line is inside the `main()` body; `mainInvoked` is true so the rule exempts it.

### [ ] Finding `37` — `BP-PY-46`

- Function context: `scripts/whatsapp-wrapped/findings/functions/37.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/whatsapp-wrapped/whatsapp_wrapped/generator.py:558:9`
- Checklist pattern: same as finding 22.

Source excerpt:

```
    except Exception as e:
        print(f"\n[!] Error: {e}")
        sys.exit(1)
```

Why this is a false positive: error reporting for the CLI top-level handler inside `main()`; exempt per rule condition, fired only due to the `epilog`-string indent-tracking artifact.

Checklist evidence: the finding line is inside the `main()` body; `mainInvoked` is true so the rule exempts it.

## True positives

### Rule `BP-PY-33` — Jinja2 autoescape Disabled

| Finding | Source | Reason |
| --- | --- | --- |
| 3 | generator.py:66 | `Environment(loader=FileSystemLoader(...), autoescape=False)` explicitly sets `autoescape=False`; rule condition is met verbatim. |

### Rule `BP-PY-46` — print Debugging In Library Code

| Finding | Source | Reason |
| --- | --- | --- |
| 4 | generator.py:74 | `print("[*] Rendering HTML report...")` inside `_render_report_template`, a library function outside the `main()`/`__main__` exempt zones. |
| 5 | generator.py:106 | `print(f"[*] Loading chat file: {chat_file.name}")` inside `_generate_report_data`. |
| 6 | generator.py:117 | `print(f"[+] Parsed {len(df)} messages ...")` inside `_generate_report_data`. |
| 7 | generator.py:121 | `print("[*] Running analytics...")` inside `_generate_report_data`. |
| 8 | generator.py:126 | `print("[*] Generating charts...")` inside `_generate_report_data`. |
| 9 | generator.py:132 | `print("[*] Generating user sparklines...")` inside `_generate_report_data`. |
| 10 | generator.py:148 | `print("[*] Calculating achievement badges...")` inside `_generate_report_data`. |
| 11 | generator.py:210 | `print(f"[+] HTML report saved: {output_path}")` inside `generate_html_report`. |
| 12 | generator.py:240 | `print("[!] Playwright not installed...")` inside nested `_generate_static`. |
| 13 | generator.py:241 | `print("[!] Then run: uv run playwright install webkit")` inside nested `_generate_static`. |
| 14 | generator.py:260 | `print("[*] Converting to static HTML with Playwright...")` inside nested `_generate_static`. |
| 15 | generator.py:342 | `print(f"[+] Static HTML report saved: {output_path}")` inside nested `_generate_static`. |
| 16 | generator.py:425 | `print(f"[+] HTML report saved: {html_path}")` inside `generate_full_report`. |
| 17 | generator.py:433 | `print("[!] Skipping static HTML generation...")` inside `generate_full_report`. |
| 18 | generator.py:434 | `print("[!] Install with: uv pip install playwright ...")` inside `generate_full_report`. |
| 21 | generator.py:438 | `print(f"[!] Static HTML generation failed: {e}")` inside `generate_full_report`. |

### Rule `BP-PY-1` — Bare Except Clause

| Finding | Source | Reason |
| --- | --- | --- |
| 19 | generator.py:436 | `except Exception as e:` in `generate_full_report`; broad handler does not re-raise (`suiteReraises` false), condition met. |
| 36 | generator.py:557 | `except Exception as e:` in `main()`; handler prints and `sys.exit(1)` but does not re-raise, condition met. |

### Rule `CWE-396` — Declaration of Catch for Generic Exception

| Finding | Source | Reason |
| --- | --- | --- |
| 20 | generator.py:436 | Line matches `^[\t ]*except\s+(?:Exception|BaseException)...:` regex exactly. |

### Rule `BP-PY-2` — Except Pass

| Finding | Source | Reason |
| --- | --- | --- |
| 38 | parser.py:233 | `except ValueError:` handler body is solely `pass`; condition met verbatim. |

### Rule `CWE-390` — Detection of Error Condition Without Action

| Finding | Source | Reason |
| --- | --- | --- |
| 39 | parser.py:233 | Except clause whose direct body is only `pass`; `exceptPassStart` matches. |

### Rule `CWE-1071` — Empty Code Block

| Finding | Source | Reason |
| --- | --- | --- |
| 40 | parser.py:233 | Narrow source pattern of an exception handler containing only `pass`; condition met. |

## Uncertain findings

No findings classified as uncertain.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/whatsapp-wrapped/chunks`
- Function evidence: `scripts/whatsapp-wrapped/findings/functions`
- Validation: `git diff --check` — pass
