# caniscrape — false-positive audit report

## Run metadata

```yaml
timestamp: 2026-08-02T00:00:00Z
repository: caniscrape
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape
branch: main
commit: 3624e55
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape
chunk_path: scripts/caniscrape/chunks
function_context_path: scripts/caniscrape/findings/functions
```

## Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/caniscrape/chunks -context-dir scripts/caniscrape/findings/functions real-repos/caniscrape`
- Findings: `338`
- Chunks reviewed: `scripts/caniscrape/chunks/Chunk_1_25.txt`, `Chunk_26_50.txt`, `Chunk_51_75.txt`, `Chunk_76_100.txt`, `Chunk_101_125.txt`, `Chunk_126_150.txt`, `Chunk_151_175.txt`, `Chunk_176_200.txt`, `Chunk_201_225.txt`, `Chunk_226_250.txt`, `Chunk_251_275.txt`, `Chunk_276_300.txt`, `Chunk_301_325.txt`, `Chunk_326_338.txt`
- Function contexts reviewed: `scripts/caniscrape/findings/functions/<id>.txt` for every proposed false positive, plus the enclosing source files (`caniscrape/cli.py`, `caniscrape/telemetry.py`, `caniscrape/config.py`, `caniscrape/upload_handler.py`, `caniscrape/diff.py`, `caniscrape/api_client.py`, `caniscrape/analyzers/*.py`, `caniscrape/commands/*.py`, `caniscrape/utils/*.py`)

## Audit checklist

- [x] Read every assigned chunk under `scripts/caniscrape/chunks`.
- [x] Read `scripts/caniscrape/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

Rule conditions verified with `./bin/goslop -explain <RULE-ID> --config templates/goslop-python.toml` (BP-PY-1/2/3/46, CWE-215/390/396/397/1071/1121/1124/1341, PERF-PY-26) and, for the narrow-pattern rules, against the detector implementations in `internal/lang/python/detectors/` (BP-PY-46 skip logic in `bad_practices/rules_observability.go`, CWE-1341 `pyTierBDoubleCloseRE`, CWE-1124 control-frame counter, CWE-1121 `minimumRouteBranches` counter, CWE-215 `sensitiveValueRE`/`isPureStringLiteral`, CWE-396 `pyGenericExceptRE`, PERF-PY-26 `decodeHotRE`/`windowHas`).

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 259 | 14, 17, 31, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 164, 165, 166, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 180, 185, 186, 187, 188, 189, 190, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 207, 208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239, 240, 243, 248, 249, 250, 252, 253, 254, 255, 256, 257, 258, 259, 260, 261, 262, 263, 264, 265, 266, 267, 268, 269, 270, 271, 272, 273, 274, 275, 276, 277, 278, 279, 280, 281, 282, 283, 284, 293, 294, 295, 296, 297, 298, 299, 300, 301, 302, 303, 304, 305, 306, 307, 308, 309, 310, 311, 312, 313, 314, 315, 316, 317, 321 |
| True positive | 79 | 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 15, 16, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 32, 33, 54, 55, 133, 162, 163, 167, 181, 182, 183, 184, 191, 192, 205, 206, 241, 242, 244, 245, 246, 247, 251, 285, 286, 287, 288, 289, 290, 291, 292, 318, 319, 320, 322, 323, 324, 325, 326, 327, 328, 329, 330, 331, 332, 333, 334, 335, 336, 337, 338 |
| Uncertain | 0 | — |

## False positives

### [ ] Finding `14` — `CWE-1341`

- Function context: `scripts/caniscrape/findings/functions/14.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/analyzers/integrity_analyzer.py:102:13`
- Checklist pattern: the two `.close()` calls pair up distinct resource handles, each released once

Source excerpt:

```
            clean_context.close()

            target_context_options = {}
            ...
            target_context.close()
```

Why this is a false positive: `pyTierBDoubleCloseRE` (`\w+\.close\s*\(\s*\)[\s\S]{0,180}\w+\.close\s*\(\s*\)`) matches any two `.close()` calls within 180 characters; here `clean_context.close()` and `browser.close()` are the matched pair around `target_context.close()` — three different handles, each closed exactly once.

Checklist evidence: the "same resource handle is released twice" condition is unmet — `clean_context`, `target_context` and `browser` are distinct Playwright objects and no handle is released more than once.

### [ ] Finding `17` — `PERF-PY-26`

- Function context: `scripts/caniscrape/findings/functions/17.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/analyzers/js_detector.py:50:1`
- Checklist pattern: one-shot parse call on a CLI scan path; the hot-path predicate is matched only by a name substring

Source excerpt:

```
            if proxies:
                proxy = random.choice(proxies)
                proxy_config = parse_proxy_for_playwright(proxy)
```

Why this is a false positive: `parse_proxy_for_playwright(proxy)` runs at most once per scan (inside `if proxies:`, no loop). `windowHas` finds the "render" needle only because the enclosing function is named `analyze_js_rendering`, so the "expensive decode/parse runs on a hot path" condition is not satisfied.

Checklist evidence: the parse is not inside a loop and the function window contains no `handle_job`/`handle_request`/`build_`/`process` call — the only "hot path" signal is the substring `render` inside the function name.

### [ ] Finding `31` — `PERF-PY-26`

- Function context: `scripts/caniscrape/findings/functions/31.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/analyzers/waf_detector.py:31:1`
- Checklist pattern: one-shot parse call after a single subprocess run; hot-path predicate matched only by a name substring

Source excerpt:

```
        result = subprocess.run(
            command,
            capture_output=True,
            text=True,
            timeout=60,
            check=False
        )

        wafs_found = parse_wafw00f_output(result.stdout, result.stderr)
```

Why this is a false positive: `parse_wafw00f_output` runs exactly once per scan and the module caches nothing because there is nothing to cache. `windowHas` finds the "process" needle only as a substring of `subprocess.run`, so the "runs on a hot path without a visible cache" condition is not satisfied.

Checklist evidence: the function window has no loop and no `handle_job`/`handle_request`/`render`/`build_`/`process` call — the only match is the substring `process` inside the identifier `subprocess`.

### [ ] Finding `321` — `CWE-215`

- Function context: `scripts/caniscrape/findings/functions/321.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/upload_handler.py:65:13`
- Checklist pattern: the "sensitive" argument is a static message containing the word `token`, not a token value

Source excerpt:

```
        if 'Authentication failed' in error_msg or 'expired' in error_msg.lower():
            print(f'\n[yellow]⚠️  Upload failed: API token expired[/yellow]')
```

Why this is a false positive: `sensitiveValueRE` matches the keyword `token` inside the fixed message text, and the f-string (which embeds no expressions) is treated as dynamic by `isPureStringLiteral`. No password, token, secret, or API key value is ever printed — the output is a constant string.

Checklist evidence: the "debug output logs a password, token, secret... that could be exposed" condition is unmet — the printed argument contains no identifier value, only the literal words "API token expired".

## False positives — BP-PY-46 (255)

All 255 findings below share one pattern: caniscrape is a click-based CLI application, and each flagged `print(` is the user-facing presentation output of a CLI command implementation (or a display/prompt/status function it calls). BP-PY-46's condition is "`print` is used for operational logging in non-script modules"; these prints are the CLI's intended output, and the detector's own skip logic (`bad_practices/rules_observability.go`: main-guard blocks, argparse-registered commands, click `.cli.command(` decorators, `print_*`/`cmd_*` functions) exists precisely to exempt this presentation layer. The findings fire only because the detector's CLI patterns miss this code shape: `@cli.command(...)` is not matched by the `.cli.command(` substring (the group is named `cli` without a leading dot), and the `commands/*.py`/`telemetry.py` command bodies carry no decorator of their own. The same reasoning as the `pprint` exemption (calgebra audit) applies: the print is the function's documented output, not operational logging.

### [ ] Finding `34` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/34.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:228:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f"[yellow]⚠️  URL scheme missing. Assuming 'http://'. Analyzing: [bold blue]{url}[/bold blue]...[/yellow]")
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `35` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/35.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:230:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'🔍 Analyzing: [bold blue]{url}[/bold blue]...')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `36` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/36.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:233:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [yellow]⚠️  Running with --find-all is aggressive and may trigger rate limits or temporary IP bans.[/yellow]\n')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `37` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/37.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:235:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [yellow]⚠️  --thorough scan selected. Behavioral analysis may take several minutes on large sites.[/yellow]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `38` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/38.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:237:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [yellow]⚠️  --deep scan selected. Behavioral analysis may take 10+ minutes on large sites.[/yellow]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `39` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/39.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:239:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('    [yellow]You have 5 seconds after the above message(s) to cancel. (Ctrl + C to cancel)[/yellow]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `40` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/40.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:246:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('Checking robots.txt...')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `41` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/41.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:250:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('Analyzing TLS fingerprint...')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `42` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/42.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:253:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('Analyzing for advanced fingerprinting...')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `43` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/43.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:256:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('Performing function integrity analysis...')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `44` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/44.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:259:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('Analyzing JavaScript rendering...')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `45` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/45.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:263:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('Analyzing for behavioral traps (default scan)...')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `46` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/46.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:265:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'Analyzing for behavioral traps ({scan_depth} scan)...')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `47` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/47.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:268:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('Detecting CAPTCHA...')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `48` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/48.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:272:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('Profiling rate limits with browser-like client...')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `49` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/49.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:274:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('Profiling rate limits with Python client...')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `50` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/50.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:277:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('Running WAF detection...')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `51` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/51.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:320:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[yellow]⚠️  Upload failed. Results cached locally.[/yellow]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `52` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/52.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:321:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[dim]    Run [cyan]caniscrape push[/cyan] to retry.[/dim]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `53` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/53.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:324:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[dim]💾 Results saved locally. Run [cyan]caniscrape push[/cyan] to upload.[/dim]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `56` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/56.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:348:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('\n')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `57` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/57.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:349:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(Rule(f"[bold white on blue] DIFFICULTY SCORE: {score_card['score']}/10 ({score_card['label']}) [/]", style="blue"))
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `58` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/58.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:350:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print()
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `59` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/59.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:352:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print()
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `60` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/60.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:353:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("[bold yellow]🛡️  ACTIVE PROTECTIONS[/bold yellow]")
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `61` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/61.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:354:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print()
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `62` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/62.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:360:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('    [red]❌ robots.txt: Explicitly disallows scraping for all bots (\'Disallow: /\')[/red]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `63` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/63.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:366:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [green]✅ robots.txt: {message}[/green]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `64` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/64.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:368:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('    [green]✅ robots.txt: Website does not have a robots.txt file (no explicit restrictions).[/green]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `65` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/65.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:370:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [yellow]⚠️  robots.txt: Could not be analyzed. Reason: {robots_result["message"]}[/yellow]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `66` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/66.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:375:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [red]❌ TLS Fingerprinting: {tls_result["details"]}[/red]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `67` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/67.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:377:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [green]✅ TLS Fingerprinting: {tls_result["details"]}[/green]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `68` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/68.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:379:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [yellow]⚠️  TLS Fingerprinting: {tls_result["details"]}[/yellow]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `69` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/69.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:399:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('    [red]❌ Advanced Bot Detection:[/red]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `70` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/70.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:401:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [red]- {flag}[/red]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `71` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/71.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:404:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('    [yellow]⚠️  Suspicious Signals:[/yellow]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `72` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/72.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:406:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [yellow]- {flag}[/yellow]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `73` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/73.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:409:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('    [green]✅ Advanced Bot Detection: No obvious fingerprinting services or signals detected.[/green]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `74` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/74.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:411:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [yellow]⚠️  Fingerprinting: Analysis failed. Reason: {fingerprint_result["message"]}[/yellow]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `75` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/75.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:418:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('    [red]❌ Browser Integrity Compromised:[/red]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `76` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/76.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:420:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [red]- Function "{func}" was modified.[/red]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `77` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/77.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:421:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'      [red]Reason: {reason}[/red]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `78` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/78.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:423:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('    [green]✅ Browser Integrity: No modifications detected.[/green]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `79` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/79.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:425:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [yellow]⚠️  Integrity Analysis: Test failed. Reason: {integrity_result["message"]}[/yellow]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `80` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/80.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:431:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [red]❌ JavaScript: Required (React/Vue/Angular SPA). {js_result["content_difference_%"]}% of content is missing without JS.[/red]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `81` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/81.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:433:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [yellow]⚠️  JavaScript: Required for some content. {js_result["content_difference_%"]}% of content is missing without JS.[/yellow]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `82` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/82.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:435:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [green]✅ JavaScript: Not required for main content.[/green]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `83` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/83.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:437:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [yellow]⚠️  JavaScript: Analysis failed. Reason: {js_result["message"]}[/yellow]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `84` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/84.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:445:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [red]❌ Behavioral Analysis: Found {count} invisible "honeypot" links (out of {checked} checked). There are many bot traps.[/red]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `85` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/85.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:447:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [green]✅ Behavioral Analysis: No obvious honeypot traps detected.[/green]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `86` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/86.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:449:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [yellow]⚠️  Behavioral Analysis: Test failed. Reason: {behavioral_result["message"]}[/yellow]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `87` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/87.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:457:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [red]❌ CAPTCHA: {captcha_type} detected ({trigger}).[/red]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `88` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/88.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:459:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'        [green]✅ {captcha_result["details"]}[/green]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `89` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/89.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:461:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'        [red]❌ {captcha_result["details"]}[/red]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `90` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/90.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:463:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'        [blue]ℹ️  {captcha_result["details"]}[/blue]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `91` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/91.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:465:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [green]✅ CAPTCHA: No CAPTCHA detected during initial analysis.[/green]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `92` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/92.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:467:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [yellow]⚠️  CAPTCHA: Analysis failed. Reason: {captcha_result["message"]}[/yellow]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `93` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/93.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:474:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [red]❌ Rate Limiting: Blocked Immediately ({results["details"]})[/red]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `94` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/94.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:475:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [yellow]💡 [bold]Advice:[/bold] This is likely due to client fingerprinting (TLS fingerprinting, User-Agent, etc.), not a classic rate limit.[/yellow]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `95` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/95.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:476:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'       [yellow]Run the analysis again. A different browser identity will be used, which may not be blocked.[/yellow]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `96` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/96.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:477:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [yellow]   Otherwise, try the --impersonate flag, it will take longer but is likely to succeed.[/yellow]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `97` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/97.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:479:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [green]✅ Rate Limiting: {results["details"]}[/green]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `98` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/98.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:482:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [yellow]⚠️  Rate Limiting: Test failed. Reason: {error_message}[/yellow]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `99` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/99.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:491:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[bold red]Error: "wafw00f" command not found.[/bold red]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `100` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/100.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:492:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[yellow]To fix this, please follow these steps in your terminal:')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `101` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/101.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:493:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[yellow]1. Install pipx: [bold]python -m pip install --user pipx[/bold][/yellow]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `102` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/102.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:494:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[yellow]2. Install wafw00f: [bold]pipx install wafw00f[/bold][/yellow]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `103` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/103.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:495:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[yellow](You may need to restart your terminal or restart your IDE after step 1 if step 2 doesn\'t work.)')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `104` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/104.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:497:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[yellow]WAF detection timed out.[/yellow]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `105` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/105.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:499:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'[yellow]WAF detection failed. Wafw00f stderr: {message}[/yellow]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `106` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/106.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:504:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('    [green]✅ No WAF detected.[/green]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `107` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/107.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:507:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [blue]ℹ️  WAF: A generic firewall or server security rule might be present (low confidence).[/blue]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `108` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/108.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:517:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [red]❌ WAF: {display_lines[0]}[/red]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `109` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/109.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:519:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'    [red]❌ WAFs Detected:[/red]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `110` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/110.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:521:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'        [red]- {line}[/red]')
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `111` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/111.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:523:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print()
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `112` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/112.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:524:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(Rule("[bold]💡 RECOMMENDATIONS[/bold]", style="cyan"))
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `113` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/113.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:526:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("\n[bold]Required Tools:[/bold]")
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `114` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/114.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:529:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f"  • {tool}")
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `115` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/115.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:531:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("  • No special tools required.")
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `116` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/116.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:533:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("\n[bold]Scraping Strategy:[/bold]")
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `117` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/117.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:535:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f"  • {tip}")
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `118` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/118.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:537:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print()
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `119` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/119.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:538:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(Rule("[bold]Analysis Complete[/bold]", style="green"))
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `120` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/120.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/cli.py:539:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print()
```

Why this is a false positive: print inside `scan_command` — the print statement is the user-facing output of the `scan` click command implementation; the caniscrape package is a CLI application (console script), and the detector's own CLI carve-outs (main-guard, argparse commands, click `.cli.command(` decorators) are intended to exempt exactly this presentation layer — `@cli.command(...)` is only missed because the decorator is registered on a group named `cli` without a leading dot. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `121` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/121.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/config_cmd.py:11:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[red]❌ Not linked to a project. Run [cyan]caniscrape init[/cyan] first.[/red]')
```

Why this is a false positive: print inside `set_config_command / show_config_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `122` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/122.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/config_cmd.py:19:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[green]✅ Auto-upload enabled[/green]')
```

Why this is a false positive: print inside `set_config_command / show_config_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `123` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/123.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/config_cmd.py:20:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[dim]Scans will now automatically sync to your cloud project.[/dim]')
```

Why this is a false positive: print inside `set_config_command / show_config_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `124` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/124.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/config_cmd.py:22:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[blue]ℹ️  Auto-upload disabled[/blue]')
```

Why this is a false positive: print inside `set_config_command / show_config_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `125` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/125.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/config_cmd.py:23:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[dim]Run [cyan]caniscrape push[/cyan] to upload scans manually.[/dim]')
```

Why this is a false positive: print inside `set_config_command / show_config_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `126` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/126.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/config_cmd.py:32:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[red]❌ Not linked to a project[/red]')
```

Why this is a false positive: print inside `set_config_command / show_config_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `127` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/127.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/config_cmd.py:33:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[dim]Run [cyan]caniscrape init[/cyan] or [cyan]caniscrape link[/cyan] to get started.[/dim]')
```

Why this is a false positive: print inside `set_config_command / show_config_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `128` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/128.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/config_cmd.py:36:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('\n[bold cyan]📋 Project Configuration[/bold cyan]\n')
```

Why this is a false positive: print inside `set_config_command / show_config_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `129` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/129.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/config_cmd.py:37:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'Project: [bold]{config.get("project_name", "Unknown")}[/bold]')
```

Why this is a false positive: print inside `set_config_command / show_config_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `130` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/130.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/config_cmd.py:38:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'Project ID: [dim]{config.get_project_id()}[/dim]')
```

Why this is a false positive: print inside `set_config_command / show_config_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `131` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/131.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/config_cmd.py:39:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'Auto-upload: {"[green]✅ Enabled[/green]" if config.get("auto_upload") else "[yellow]❌ Disabled[/yellow]"}')
```

Why this is a false positive: print inside `set_config_command / show_config_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `132` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/132.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/config_cmd.py:40:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print()
```

Why this is a false positive: print inside `set_config_command / show_config_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `134` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/134.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:11:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[bold blue]🚀 caniscrape Cloud Setup[/bold blue]\n')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `135` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/135.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:16:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'[yellow]⚠️  This directory is already linked to project: [bold]{project_name}[/bold][/yellow]')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `136` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/136.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:19:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[dim]Cancelled.[/dim]')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `137` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/137.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:23:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[green]✅ Unlinked from previous project.[/green]\n')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `138` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/138.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:25:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[bold]Step 1: Authentication[/bold]')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `139` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/139.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:26:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('You need to authenticate with caniscrape Cloud.')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `140` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/140.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:27:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('Visit [link]https://caniscrape.org/dashboard/settings [/link]to get your API token.\n')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `141` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/141.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:32:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[red]❌ Token cannot be empty.[/red]')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `142` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/142.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:37:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[dim]Authenticating...[/dim]')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `143` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/143.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:42:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'[green]✅ Authenticated successfully![/green]')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `144` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/144.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:45:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'\n[dim]You have {len(projects)} existing project(s):[/dim]')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `145` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/145.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:47:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'  - {p["name"]}')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `146` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/146.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:49:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'  ... and {len(projects) - 3} more')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `147` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/147.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:51:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'\n[dim]You don\'t have any existing projects. Create one now![/dim]')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `148` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/148.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:53:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'[red]❌ Authentication failed: {str(e)}[/red]')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `149` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/149.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:56:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('\n[bold]Step 2: Create Project[/bold]')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `150` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/150.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:57:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('Give this project a name (e.g., "E-commerce Scraper", "News Aggregator")\n')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `151` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/151.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:62:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[red]❌ Project name cannot be empty.[/red]')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `152` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/152.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:67:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('\n[dim]Creating project...[/dim]')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `153` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/153.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:70:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'[green]✅ Project created: {project["name"]}[/green]')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `154` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/154.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:72:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'[red]❌ Failed to create project: {str(e)}[/red]')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `155` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/155.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:75:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('\n[bold]Step 3: Auto-Upload Settings[/bold]')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `156` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/156.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:76:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('Automatically push scan results to your cloud project?\n')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `157` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/157.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:77:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[dim]  • Yes: Every scan uploads instantly (recommended)[/dim]')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `158` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/158.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:78:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[dim]  • No: Results stay local until you run "caniscrape push"[/dim]\n')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `159` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/159.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:86:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[green]✅ Auto-upload enabled[/green]')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `160` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/160.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:88:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[blue]ℹ️  Auto-upload disabled. Use [cyan]caniscrape push[/cyan] to upload manually.[/blue]')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `161` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/161.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:98:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'[green]✅ Configuration saved to .caniscrape/config[/green]')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `164` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/164.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:100:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'[red]❌ Failed to save config: {str(e)}[/red]')
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `165` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/165.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:103:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('\n' + '='*60)
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `166` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/166.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/init.py:104:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(Panel.fit(
```

Why this is a false positive: print inside `init_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `168` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/168.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/link.py:16:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[bold blue]🔗 Link to Existing Project[/bold blue]\n')
```

Why this is a false positive: print inside `link_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `169` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/169.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/link.py:21:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'[yellow]⚠️  Already linked to: [bold]{project_name}[/bold][/yellow]')
```

Why this is a false positive: print inside `link_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `170` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/170.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/link.py:24:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[dim]Cancelled.[/dim]')
```

Why this is a false positive: print inside `link_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `171` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/171.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/link.py:28:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[green]✅ Unlinked.[/green]\n')
```

Why this is a false positive: print inside `link_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `172` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/172.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/link.py:30:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[bold]Step 1: Authentication[/bold]')
```

Why this is a false positive: print inside `link_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `173` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/173.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/link.py:31:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('Visit [link]https://caniscrape.org/settings[/link] to get your API token.\n')
```

Why this is a false positive: print inside `link_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `174` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/174.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/link.py:36:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[red]❌ Token cannot be empty.[/red]')
```

Why this is a false positive: print inside `link_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `175` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/175.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/link.py:42:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('\n[dim]Fetching your projects...[/dim]')
```

Why this is a false positive: print inside `link_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `176` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/176.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/link.py:46:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'[red]❌ Authentication failed: {str(e)}[/red]')
```

Why this is a false positive: print inside `link_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `177` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/177.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/link.py:50:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[yellow]⚠️  You have no projects yet.[/yellow]')
```

Why this is a false positive: print inside `link_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `178` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/178.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/link.py:51:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[dim]Create one at [link]https://caniscrape.org/projects[/link] or run [cyan]caniscrape init[/cyan][/dim]')
```

Why this is a false positive: print inside `link_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `179` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/179.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/link.py:54:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'\n[bold]Step 2: Select Project[/bold]')
```

Why this is a false positive: print inside `link_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `180` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/180.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/link.py:55:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'[dim]Found {len(projects)} project(s)[/dim]\n')
```

Why this is a false positive: print inside `link_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `185` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/185.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/link.py:79:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(table)
```

Why this is a false positive: print inside `link_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `186` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/186.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/link.py:80:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print()
```

Why this is a false positive: print inside `link_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `187` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/187.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/link.py:88:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[red]Please enter a number.[/red]')
```

Why this is a false positive: print inside `link_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `188` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/188.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/link.py:99:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[red]Invalid choice. Try again.[/red]')
```

Why this is a false positive: print inside `link_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `189` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/189.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/link.py:101:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'\n[bold]Step 3: Auto-Upload Settings[/bold]')
```

Why this is a false positive: print inside `link_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `190` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/190.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/link.py:102:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('Automatically push scan results to this project?\n')
```

Why this is a false positive: print inside `link_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `193` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/193.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/link.py:115:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'[red]❌ Failed to save config: {str(e)}[/red]')
```

Why this is a false positive: print inside `link_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `194` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/194.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/link.py:118:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('\n' + '='*60)
```

Why this is a false positive: print inside `link_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `195` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/195.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/link.py:119:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(Panel.fit(
```

Why this is a false positive: print inside `link_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `196` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/196.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/push.py:13:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[bold blue]📤 Pushing scan results to cloud...[/bold blue]\n')
```

Why this is a false positive: print inside `push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `197` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/197.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/push.py:18:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[red]❌ Not linked to a cloud project.[/red]')
```

Why this is a false positive: print inside `push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `198` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/198.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/push.py:19:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[yellow]Run [cyan]caniscrape init[/cyan] to link this directory to a project.[/yellow]')
```

Why this is a false positive: print inside `push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `199` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/199.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/push.py:30:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[yellow]⚠️  No local scan results found.[/yellow]')
```

Why this is a false positive: print inside `push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `200` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/200.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/push.py:31:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[dim]Run some scans first, then push them to the cloud.[/dim]')
```

Why this is a false positive: print inside `push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `201` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/201.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/push.py:37:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[yellow]⚠️  No scan results to push.[/yellow]')
```

Why this is a false positive: print inside `push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `202` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/202.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/push.py:40:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'[dim]Found {len(scan_files)} scan result(s) to push...[/dim]\n')
```

Why this is a false positive: print inside `push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `203` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/203.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/push.py:62:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'[green]✅ Pushed: {url}[/green]')
```

Why this is a false positive: print inside `push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `204` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/204.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/push.py:66:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'[red]❌ Failed to push {file.name}: {str(e)}[/red]')
```

Why this is a false positive: print inside `push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `207` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/207.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/push.py:69:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'[yellow]⚠️  Error reading {file.name}: {str(e)}[/yellow]')
```

Why this is a false positive: print inside `push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `208` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/208.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/push.py:72:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print()
```

Why this is a false positive: print inside `push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `209` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/209.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/push.py:74:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'[green]🚀 Successfully pushed {success_count} scan(s) to \'{project_name}\'[/green]')
```

Why this is a false positive: print inside `push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `210` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/210.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/push.py:76:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'[yellow]⚠️  Failed to push {failed_count} scan(s)[/yellow]')
```

Why this is a false positive: print inside `push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `211` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/211.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:12:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[bold cyan]📊 Contribute to Public Telemetry[/bold cyan]\n')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `212` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/212.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:17:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[red]❌ Not linked to a cloud project.[/red]')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `213` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/213.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:18:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[yellow]Run [cyan]caniscrape init[/cyan] first to set up your account.[/yellow]')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `214` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/214.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:24:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[bold]What is public telemetry?[/bold]')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `215` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/215.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:25:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('Public telemetry aggregates scan data from all users to create a')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `216` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/216.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:26:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('searchable database of website protections (like Shodan for anti-bot defenses).\n')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `217` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/217.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:28:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[bold green]Benefits:[/bold green]')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `218` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/218.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:29:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('  • See how site protections change over time')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `219` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/219.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:30:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('  • Search for any URL to see its defense history')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `220` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/220.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:31:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('  • Compare different sites\' protection strategies')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `221` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/221.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:32:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('  • Free access while we build the database\n')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `222` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/222.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:34:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[bold yellow]What we collect:[/bold yellow]')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `223` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/223.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:35:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('  • URLs you scan')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `224` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/224.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:36:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('  • Scan results (WAF, CAPTCHA, rate limits, etc.)')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `225` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/225.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:37:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('  • Timestamps of scans\n')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `226` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/226.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:39:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[bold red]What we DON\'T collect:[/bold red]')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `227` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/227.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:40:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('  • Your IP address or personal info')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `228` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/228.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:41:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('  • Any authentication tokens or credentials\n')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `229` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/229.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:43:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[dim]You can disable this anytime with: caniscrape telemetry disable[/dim]\n')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `230` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/230.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:46:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[blue]Telemetry contributions disabled.[/blue]')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `231` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/231.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:51:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[green]✅ Telemetry contributions enabled![/green]\n')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `232` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/232.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:57:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[dim]Checking for scans to contribute...[/dim]')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `233` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/233.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:66:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[yellow]⚠️  No scans found in your project.[/yellow]')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `234` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/234.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:67:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[dim]Run some scans first, then contribute them to telemetry.[/dim]')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `235` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/235.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:82:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'[green]✅ Contributed: {url}[/green]')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `236` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/236.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:85:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'[yellow]⚠️  Failed to contribute {url}: {str(e)}[/yellow]')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `237` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/237.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:88:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'\n[green]✨ Successfully contributed {contributed} scan(s) to public telemetry![/green]')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `238` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/238.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:89:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('[dim]Thank you for helping build the database![/dim]')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `239` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/239.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:91:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print('\n[dim]All scans have already been contributed to telemetry.[/dim]')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `240` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/240.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:93:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'[red]❌ Failed to contribute to telemetry: {str(e)}[/red]')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `243` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/243.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/commands/telemetry_push.py:95:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'[red]❌ Unexpected error: {str(e)}[/red]')
```

Why this is a false positive: print inside `telemetry_push_command` — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `252` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/252.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:87:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("\n[bold cyan]📊 Help us improve caniscrape![/bold cyan]\n")
```

Why this is a false positive: print inside `prompt_usage_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `253` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/253.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:88:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("We'd like to collect [bold]anonymous usage data[/bold] to improve the tool.")
```

Why this is a false positive: print inside `prompt_usage_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `254` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/254.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:89:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("\n[dim]What we collect:[/dim]")
```

Why this is a false positive: print inside `prompt_usage_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `255` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/255.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:90:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("  • CLI version and Python version")
```

Why this is a false positive: print inside `prompt_usage_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `256` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/256.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:91:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("  • Operating system type")
```

Why this is a false positive: print inside `prompt_usage_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `257` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/257.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:92:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("  • Commands used (scan/init)")
```

Why this is a false positive: print inside `prompt_usage_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `258` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/258.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:93:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("  • Success/failure rates")
```

Why this is a false positive: print inside `prompt_usage_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `259` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/259.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:94:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("  • Error types (no URLs or personal data)")
```

Why this is a false positive: print inside `prompt_usage_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `260` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/260.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:96:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("\n[dim]What we DON'T collect:[/dim]")
```

Why this is a false positive: print inside `prompt_usage_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `261` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/261.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:97:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("  • Your name, email, or IP address")
```

Why this is a false positive: print inside `prompt_usage_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `262` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/262.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:98:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("  • URLs you scan")
```

Why this is a false positive: print inside `prompt_usage_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `263` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/263.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:99:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("  • Any personally identifiable information")
```

Why this is a false positive: print inside `prompt_usage_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `264` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/264.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:101:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("\n[dim]You can opt-out anytime with:[/dim] [bold]caniscrape telemetry usage off[/bold]\n")
```

Why this is a false positive: print inside `prompt_usage_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `265` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/265.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:110:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("\n[green]✅ Usage telemetry enabled.[/green]")
```

Why this is a false positive: print inside `prompt_usage_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `266` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/266.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:112:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("\n[blue]Usage telemetry disabled.[/blue]")
```

Why this is a false positive: print inside `prompt_usage_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `267` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/267.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:123:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("\n[bold cyan]🌍 Contribute to Public Scan Database[/bold cyan]\n")
```

Why this is a false positive: print inside `prompt_scan_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `268` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/268.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:124:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("Help build a [bold]searchable database of website protections[/bold].")
```

Why this is a false positive: print inside `prompt_scan_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `269` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/269.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:126:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("\n[dim]What we collect:[/dim]")
```

Why this is a false positive: print inside `prompt_scan_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `270` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/270.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:127:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("  • URLs you scan")
```

Why this is a false positive: print inside `prompt_scan_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `271` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/271.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:128:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("  • Scan results (WAF, CAPTCHA, difficulty scores)")
```

Why this is a false positive: print inside `prompt_scan_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `272` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/272.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:129:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("  • Protection details")
```

Why this is a false positive: print inside `prompt_scan_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `273` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/273.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:131:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("\n[dim]What we DON'T collect:[/dim]")
```

Why this is a false positive: print inside `prompt_scan_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `274` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/274.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:132:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("  • Your name, email, or IP address")
```

Why this is a false positive: print inside `prompt_scan_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `275` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/275.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:133:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("  • Any authentication tokens")
```

Why this is a false positive: print inside `prompt_scan_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `276` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/276.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:135:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("\n[bold green]Benefits:[/bold green]")
```

Why this is a false positive: print inside `prompt_scan_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `277` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/277.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:136:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("  • Search any URL to see its protection history")
```

Why this is a false positive: print inside `prompt_scan_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `278` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/278.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:137:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("  • Compare sites' defense strategies")
```

Why this is a false positive: print inside `prompt_scan_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `279` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/279.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:138:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("  • Track changes over time")
```

Why this is a false positive: print inside `prompt_scan_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `280` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/280.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:139:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("  • The more everyone contributes the better it is!")
```

Why this is a false positive: print inside `prompt_scan_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `281` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/281.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:141:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("\n[dim]You can opt-out anytime with:[/dim] [bold]caniscrape telemetry scans off[/bold]")
```

Why this is a false positive: print inside `prompt_scan_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `282` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/282.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:142:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("[dim]This is separate from usage telemetry.[/dim]\n")
```

Why this is a false positive: print inside `prompt_scan_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `283` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/283.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:151:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("\n[green]✅ Scan contributions enabled. Thank you![/green]")
```

Why this is a false positive: print inside `prompt_scan_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `284` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/284.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:153:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("\n[blue]Scan contributions disabled.[/blue]")
```

Why this is a false positive: print inside `prompt_scan_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `293` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/293.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:261:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("[green]✅ Usage telemetry enabled.[/green]")
```

Why this is a false positive: print inside `enable_usage_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `294` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/294.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:271:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("[blue]Usage telemetry disabled.[/blue]")
```

Why this is a false positive: print inside `disable_usage_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `295` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/295.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:281:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("[green]✅ Scan contributions enabled.[/green]")
```

Why this is a false positive: print inside `enable_scan_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `296` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/296.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:291:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("[blue]Scan contributions disabled.[/blue]")
```

Why this is a false positive: print inside `disable_scan_telemetry`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `297` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/297.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:298:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("[yellow]No telemetry data found locally.[/yellow]")
```

Why this is a false positive: print inside `request_data_deletion`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `298` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/298.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:305:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("[yellow]No device ID found. No data to delete.[/yellow]")
```

Why this is a false positive: print inside `request_data_deletion`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `299` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/299.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:308:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f"\n[bold yellow]⚠️  Delete all usage telemetry data?[/bold yellow]")
```

Why this is a false positive: print inside `request_data_deletion`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `300` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/300.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:309:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f"[dim]Device ID: {device_id}[/dim]")
```

Why this is a false positive: print inside `request_data_deletion`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `301` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/301.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:310:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f"[dim]Note: This only deletes usage stats, not public scan contributions.[/dim]\n")
```

Why this is a false positive: print inside `request_data_deletion`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `302` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/302.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:315:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("[blue]Deletion cancelled.[/blue]")
```

Why this is a false positive: print inside `request_data_deletion`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `303` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/303.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:326:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("\n[green]✅ All usage telemetry data deleted.[/green]")
```

Why this is a false positive: print inside `request_data_deletion`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `304` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/304.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:330:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f"\n[red]❌ Failed to delete data: {error}[/red]")
```

Why this is a false positive: print inside `request_data_deletion`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `305` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/305.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:334:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f"\n[red]❌ Network error: {str(e)}[/red]")
```

Why this is a false positive: print inside `request_data_deletion`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `306` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/306.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:343:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("\n[bold cyan]📊 Telemetry Status[/bold cyan]\n")
```

Why this is a false positive: print inside `show_status`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `307` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/307.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:346:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f"[bold]Usage Telemetry:[/bold] {'[green]✅ Enabled[/green]' if usage_enabled else '[red]❌ Disabled[/red]'}")
```

Why this is a false positive: print inside `show_status`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `308` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/308.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:348:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f"  [dim]Decided: {config['usage_telemetry_decided_at'][:10]}[/dim]")
```

Why this is a false positive: print inside `show_status`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `309` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/309.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:351:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f"\n[bold]Public Scan Contributions:[/bold] {'[green]✅ Enabled[/green]' if scan_enabled else '[red]❌ Disabled[/red]'}")
```

Why this is a false positive: print inside `show_status`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `310` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/310.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:353:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f"  [dim]Decided: {config['scan_telemetry_decided_at'][:10]}[/dim]")
```

Why this is a false positive: print inside `show_status`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `311` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/311.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:356:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f"\n[dim]Device ID: {config['device_id']}[/dim]")
```

Why this is a false positive: print inside `show_status`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `312` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/312.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:358:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("\n[bold]Commands:[/bold]")
```

Why this is a false positive: print inside `show_status`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `313` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/313.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:359:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("  • [cyan]caniscrape telemetry usage on/off[/cyan]  - Toggle usage telemetry")
```

Why this is a false positive: print inside `show_status`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `314` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/314.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:360:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("  • [cyan]caniscrape telemetry scans on/off[/cyan]  - Toggle scan contributions")
```

Why this is a false positive: print inside `show_status`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `315` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/315.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:361:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("  • [cyan]caniscrape telemetry delete[/cyan]       - Delete usage data (GDPR)")
```

Why this is a false positive: print inside `show_status`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `316` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/316.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:362:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print("  • [cyan]caniscrape telemetry status[/cyan]       - Show this status")
```

Why this is a false positive: print inside `show_status`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `317` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/317.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/telemetry.py:363:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print()
```

Why this is a false positive: print inside `show_status`, the implementation of the `caniscrape telemetry ...` commands (wired via click wrappers `usage_on`/`usage_off`/`scans_on`/`scans_off`/`delete`/`telemetry_status` in `cli.py`) — the print statement is the user-facing output of a CLI command implementation (invoked from the thin click wrappers in `cli.py`), not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `248` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/248.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/diff.py:103:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(f'\n[dim]ℹ️  No changes detected since last scan ({previous_scan_date})[/dim]')
```

Why this is a false positive: print inside `display_diff`, a presentation function whose documented purpose is to render the scan comparison to stdout (same shape as the `pprint` exemption) — the print is the function's intended output, not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `249` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/249.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/diff.py:131:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print()
```

Why this is a false positive: print inside `display_diff`, a presentation function whose documented purpose is to render the scan comparison to stdout (same shape as the `pprint` exemption) — the print is the function's intended output, not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


### [ ] Finding `250` — `BP-PY-46`

- Function context: `scripts/caniscrape/findings/functions/250.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/caniscrape/caniscrape/diff.py:132:1`
- Checklist pattern: `print` is user-facing CLI output, not operational logging in a non-script module

Source excerpt:

```
print(Panel(
```

Why this is a false positive: print inside `display_diff`, a presentation function whose documented purpose is to render the scan comparison to stdout (same shape as the `pprint` exemption) — the print is the function's intended output, not operational logging. The rule's condition (`print` used for operational logging in non-script modules) is not satisfied.

Checklist evidence: the flagged `print(` implements the CLI's presentation output (report, prompt, or status) inside a command implementation, so the "operational logging" predicate of the rule condition is unmet.


## True positives

### BP-PY-1 — Bare Except Clause (22)

Rule condition (`rules_core.go` `detectBPPY1`): a bare `except:` is always flagged; a broad `except Exception`/`except BaseException` is flagged unless the suite clearly re-raises (`suiteReraises`). None of these suites re-raises.

| Finding | Source | Reason |
| --- | --- | --- |
| 1 | `caniscrape/analyzers/behavioral_detector.py:63` | `except Exception as e:` returns an error dict; no re-raise |
| 3 | `caniscrape/analyzers/captcha_detector.py:62` | bare `except:` swallows page-load errors |
| 9 | `caniscrape/analyzers/captcha_detector.py:112` | `except Exception as e:` returns error dict; no re-raise |
| 12 | `caniscrape/analyzers/fingerprint_analyzer.py:165` | `except Exception as e:` sets message and returns; no re-raise |
| 15 | `caniscrape/analyzers/integrity_analyzer.py:120` | `except Exception as e:` sets message and returns; no re-raise |
| 18 | `caniscrape/analyzers/js_detector.py:62` | bare `except:` swallows page-load errors |
| 22 | `caniscrape/analyzers/js_detector.py:83` | `except Exception as e:` returns error dict; no re-raise |
| 24 | `caniscrape/analyzers/rate_limit_profiler.py:44` | `except Exception:` returns 999; no re-raise |
| 26 | `caniscrape/analyzers/rate_limit_profiler.py:99` | `except Exception as e:` returns error dict; no re-raise |
| 29 | `caniscrape/analyzers/tls_analyzer.py:41` | `except Exception as e:` sets flag; no re-raise |
| 32 | `caniscrape/api_client.py:52` | bare `except:` inside HTTPError handling |
| 55 | `caniscrape/cli.py:344` | bare `except:` swallows date-parse failure |
| 162 | `caniscrape/commands/init.py:99` | `except Exception as e:` prints and returns; no re-raise |
| 181 | `caniscrape/commands/link.py:69` | bare `except:` with pass |
| 191 | `caniscrape/commands/link.py:114` | `except Exception as e:` prints and returns; no re-raise |
| 205 | `caniscrape/commands/push.py:68` | `except Exception as e:` prints and continues; no re-raise |
| 241 | `caniscrape/commands/telemetry_push.py:94` | `except Exception as e:` prints; no re-raise |
| 251 | `caniscrape/telemetry.py:40` | bare `except:` returns {} |
| 285 | `caniscrape/telemetry.py:197` | `except Exception:` with pass |
| 292 | `caniscrape/telemetry.py:250` | `except Exception:` returns False |
| 326 | `caniscrape/upload_handler.py:74` | `except Exception:` returns False |
| 328 | `caniscrape/upload_handler.py:95` | bare `except:` returns None |

### BP-PY-2 — Except Pass (4)

Rule condition (`detectBPPY2`): the except suite consists solely of `pass`.

| Finding | Source | Reason |
| --- | --- | --- |
| 4 | `caniscrape/analyzers/captcha_detector.py:62` | bare `except:` suite is only `pass` |
| 19 | `caniscrape/analyzers/js_detector.py:62` | bare `except:` suite is only `pass` |
| 182 | `caniscrape/commands/link.py:69` | bare `except:` suite is only `pass` |
| 286 | `caniscrape/telemetry.py:197` | `except Exception:` suite is only `pass` |

### BP-PY-3 — Raise Generic Exception (1)

Rule condition (`detectBPPY3`): code raises bare `Exception`/`BaseException`.

| Finding | Source | Reason |
| --- | --- | --- |
| 246 | `caniscrape/config.py:41` | `raise Exception(f'Failed to save config.')` |

### CWE-390 — Detection of Error Condition Without Action (4)

Rule condition (`detectCWE390`): an except clause whose direct body is `pass`.

| Finding | Source | Reason |
| --- | --- | --- |
| 5 | `caniscrape/analyzers/captcha_detector.py:62` | handler takes no action |
| 20 | `caniscrape/analyzers/js_detector.py:62` | handler takes no action |
| 183 | `caniscrape/commands/link.py:69` | handler takes no action |
| 287 | `caniscrape/telemetry.py:197` | handler takes no action |

### CWE-396 — Declaration of Catch for Generic Exception (16)

Rule condition (`detectCWE396`): the first `except Exception`/`except BaseException` in a non-test module (`pyGenericExceptRE`); no re-raise exemption.

| Finding | Source | Reason |
| --- | --- | --- |
| 2 | `caniscrape/analyzers/behavioral_detector.py:63` | generic handler hides distinct failures |
| 10 | `caniscrape/analyzers/captcha_detector.py:112` | generic handler hides distinct failures |
| 13 | `caniscrape/analyzers/fingerprint_analyzer.py:165` | generic handler hides distinct failures |
| 16 | `caniscrape/analyzers/integrity_analyzer.py:120` | generic handler hides distinct failures |
| 23 | `caniscrape/analyzers/js_detector.py:83` | generic handler hides distinct failures |
| 25 | `caniscrape/analyzers/rate_limit_profiler.py:44` | generic handler hides distinct failures |
| 30 | `caniscrape/analyzers/tls_analyzer.py:41` | generic handler hides distinct failures |
| 54 | `caniscrape/cli.py:326` | generic handler in scan command |
| 163 | `caniscrape/commands/init.py:99` | generic handler hides distinct failures |
| 192 | `caniscrape/commands/link.py:114` | generic handler hides distinct failures |
| 206 | `caniscrape/commands/push.py:68` | generic handler hides distinct failures |
| 242 | `caniscrape/commands/telemetry_push.py:94` | generic handler hides distinct failures |
| 245 | `caniscrape/config.py:38` | generic handler hides distinct failures |
| 288 | `caniscrape/telemetry.py:197` | generic handler hides distinct failures |
| 327 | `caniscrape/upload_handler.py:74` | generic handler hides distinct failures |
| 329 | `caniscrape/utils/captcha_solvers.py:39` | generic handler wraps initialization errors |

### CWE-397 — Declaration of Throws for Generic Exception (1)

Rule condition (`detectCWE397`): `raise Exception`/`raise BaseException` (`pyGenericRaiseRE`).

| Finding | Source | Reason |
| --- | --- | --- |
| 247 | `caniscrape/config.py:41` | `raise Exception(f'Failed to save config.')` |

### CWE-1071 — Empty Code Block (4)

Rule condition (`detectCWE1071`): an exception handler containing only `pass`.

| Finding | Source | Reason |
| --- | --- | --- |
| 6 | `caniscrape/analyzers/captcha_detector.py:62` | handler body is only `pass` |
| 21 | `caniscrape/analyzers/js_detector.py:62` | handler body is only `pass` |
| 184 | `caniscrape/commands/link.py:69` | handler body is only `pass` |
| 289 | `caniscrape/telemetry.py:197` | handler body is only `pass` |

### CWE-1121 — Excessive McCabe Cyclomatic Complexity (3)

Rule condition (`detectCWE1121`): a function body whose masked code contains at least 12 occurrences of `if `/`elif `/`for `/`while `/`except `.

| Finding | Source | Reason |
| --- | --- | --- |
| 33 | `caniscrape/cli.py:220` | `scan_command` counts 30+ control-flow keywords |
| 133 | `caniscrape/commands/init.py:10` | `init_command` counts 13 control-flow keywords |
| 167 | `caniscrape/commands/link.py:12` | `link_command` counts 12+ control-flow keywords |

### CWE-1124 — Excessively Deep Nesting (3)

Rule condition (`detectCWE1124`): an executable statement nested under at least six control frames (`try`/`if`/`for`/`while`/`with`).

| Finding | Source | Reason |
| --- | --- | --- |
| 8 | `caniscrape/analyzers/captcha_detector.py:79` | `raise` sits under try/with/try/if/if/if (6 frames) |
| 11 | `caniscrape/analyzers/fingerprint_analyzer.py:145` | `append` sits under try/with/for/for/if/if (6 frames) |
| 27 | `caniscrape/analyzers/robots_checker.py:63` | assignment sits under try/if/for/if/if/if (6 frames) |

### BP-PY-46 — print Debugging In Library Code (21)

These prints are operational/informational output inside worker and utility functions (analyzers, solvers, parsers, config loader, upload/telemetry workers) — not CLI presentation output, so the rule condition ("`print` is used for operational logging in non-script modules") is satisfied.

| Finding | Source | Reason |
| --- | --- | --- |
| 7 | `caniscrape/analyzers/captcha_detector.py:72` | `print(f"INFO: {captcha_on_load} detected. Attempting to solve...")` in analyzer worker |
| 28 | `caniscrape/analyzers/robots_checker.py:75` | `print("Entering else block")` — debug leftover in analyzer worker |
| 244 | `caniscrape/config.py:27` | `print(f"[yellow]Warning: Config file corrupted, resetting: ...")` in `Config._load` worker |
| 290 | `caniscrape/telemetry.py:243` | `print('[dim]🌍 Scan contributed to public database (new entry)[/dim]')` in `contribute_scan` worker |
| 291 | `caniscrape/telemetry.py:245` | `print('[dim]🌍 Scan contributed to public database (updated)[/dim]')` in `contribute_scan` worker |
| 318 | `caniscrape/upload_handler.py:58` | `print(f'\n[green]✨ Results synced to \'{project_name}\'[/green]')` in `try_upload_scan` worker |
| 319 | `caniscrape/upload_handler.py:59` | `print(f'[dim]   View: https://caniscrape.org/projects/{project_id}[/dim]')` in worker |
| 320 | `caniscrape/upload_handler.py:65` | `print(f'\n[yellow]⚠️  Upload failed: API token expired[/yellow]')` in worker error path |
| 322 | `caniscrape/upload_handler.py:66` | `print(f'[dim]   Run [cyan]caniscrape init[/cyan] to re-authenticate[/dim]')` in worker error path |
| 323 | `caniscrape/upload_handler.py:68` | `print(f'\n[yellow]⚠️  Upload failed: Rate limit exceeded[/yellow]')` in worker error path |
| 324 | `caniscrape/upload_handler.py:69` | `print(f'[dim]   Results cached. Try again later.[/dim]')` in worker error path |
| 325 | `caniscrape/upload_handler.py:71` | `print(f'\n[yellow]⚠️  Upload failed: {error_msg}[/yellow]')` in worker error path |
| 330 | `caniscrape/utils/captcha_solvers.py:43` | `print('INFO: Solving reCAPTCHA v2 with CapSolver...')` in solver worker |
| 331 | `caniscrape/utils/captcha_solvers.py:81` | `print(f'INFO: Using proxy {proxy_address}:{proxy_port} for CAPTCHA solving')` |
| 332 | `caniscrape/utils/captcha_solvers.py:98` | `print('INFO: Solving hCaptcha with CapSolver...')` |
| 333 | `caniscrape/utils/captcha_solvers.py:136` | `print(f'INFO: Using proxy {proxy_address}:{proxy_port} for CAPTCHA solving')` |
| 334 | `caniscrape/utils/captcha_solvers.py:168` | `print('INFO: Solving reCAPTCHA v2 with 2Captcha...')` |
| 335 | `caniscrape/utils/captcha_solvers.py:178` | `print(f'INFO: Using proxy for CAPTCHA solving')` |
| 336 | `caniscrape/utils/captcha_solvers.py:189` | `print('INFO: Solving hCAPTCHA using 2Captcha...')` |
| 337 | `caniscrape/utils/captcha_solvers.py:199` | `print(f'INFO: Using proxy for CAPTCHA solving')` |
| 338 | `caniscrape/utils/playwright_proxy_parser.py:24` | `print(f"⚠️  Failed to parse proxy: {proxy_url[:50]}...")` in parser worker |

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/caniscrape/chunks`
- Function evidence: `scripts/caniscrape/findings/functions`
- Validation: `git diff --check` — pass
