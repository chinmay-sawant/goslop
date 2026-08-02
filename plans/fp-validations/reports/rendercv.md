# False-positive audit — rendercv

## Run metadata

```yaml
timestamp: 2026-08-02T07:33:44Z
repository: rendercv
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv
branch: main
commit: 1d4b87bc427e4cf61c0ef49623c971b0e2224708
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv
chunk_path: scripts/rendercv/chunks
function_context_path: scripts/rendercv/findings/functions
```

## Scan evidence

- Build command: `(pre-built) ./bin/goslop`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/rendercv/chunks -context-dir scripts/rendercv/findings/functions real-repos/rendercv`
- Findings: `73`
- Chunks reviewed: `scripts/rendercv/chunks/Chunk_1_25.txt`, `scripts/rendercv/chunks/Chunk_26_50.txt`, `scripts/rendercv/chunks/Chunk_51_73.txt`
- Function contexts reviewed: `scripts/rendercv/findings/functions/1.txt` .. `scripts/rendercv/findings/functions/73.txt`

## Audit checklist

- [x] Read every assigned chunk under `scripts/rendercv/chunks`.
- [x] Read `scripts/rendercv/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 46 | 3, 5, 6, 7, 9, 10, 11, 12, 13, 14, 15, 16, 17, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 35, 36, 38, 39, 40, 41, 42, 43, 44, 45, 46, 49, 50, 55, 56, 57, 64, 65, 71, 72, 73 |
| True positive | 27 | 1, 2, 4, 8, 18, 31, 32, 33, 34, 37, 47, 48, 51, 52, 53, 54, 58, 59, 60, 61, 62, 63, 66, 67, 68, 69, 70 |
| Uncertain | 0 | — |

## False positives

### [ ] Finding `3` — `BP-PY-7`

- Function context: `scripts/rendercv/findings/functions/3.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/ats_proof/analyze_pdfs.py:46:15`
- Checklist pattern: the flagged call is a third-party library method (`fitz.open`), not the builtin `open()` file-handle API the rule targets

Source excerpt:

```
def extract_pymupdf(pdf_path: Path) -> str:
    """Extract text using PyMuPDF."""
    doc = fitz.open(str(pdf_path))
    text = "".join(page.get_text() for page in doc)
    doc.close()
    return text
```

Why this is a false positive: `fitz.open` is PyMuPDF's document opener, not Python's builtin `open`; the rule's fix (`with open(...) as f:`) is not applicable to a third-party document object, and the resource is explicitly released by `doc.close()` on the next line.

Checklist evidence: the "open without `with` risks resource leaks" condition targets builtin file handles; here the file-backed resource is closed explicitly, and no builtin `open` is involved.

### [ ] Finding `5` — `BP-PY-11`

- Function context: `scripts/rendercv/findings/functions/5.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/ats_proof/analyze_pdfs.py:106:20`
- Checklist pattern: the receiver is a `ruamel.yaml.YAML()` instance (import at line 25), whose `load()` is safe by default and cannot construct arbitrary Python objects

Source excerpt:

```
from ruamel.yaml import YAML
...
        yaml = YAML()
        with yaml_path.open(encoding="utf-8") as f:
            data = yaml.load(f)
```

Why this is a false positive: the flagged `yaml` is a `ruamel.yaml.YAML()` instance — ruamel's constructor does not support arbitrary object instantiation, so the rule's premise "yaml.load without SafeLoader can execute code" does not apply; `yaml.safe_load` is not even available on this API.

Checklist evidence: `detectBPPY11` fires on the text `yaml.load(` without a `SafeLoader` argument, but the receiver is ruamel's safe YAML class, not PyYAML's `yaml` module.

### [ ] Finding `6` — `CWE-502`

- Function context: `scripts/rendercv/findings/functions/6.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/ats_proof/analyze_pdfs.py:106:20`
- Checklist pattern: same construct as finding 5 — ruamel `YAML()` instance, no untrusted-deserialization surface

Source excerpt:

```
from ruamel.yaml import YAML
...
        yaml = YAML()
        with yaml_path.open(encoding="utf-8") as f:
            data = yaml.load(f)
```

Why this is a false positive: CWE-502's condition ("yaml.load without a SafeLoader on attacker-controlled bytes") is unmet — ruamel `YAML().load` is a safe constructor that cannot instantiate arbitrary objects.

Checklist evidence: `detectCWE502` exempts only explicit SafeLoader arguments; the shown source's receiver is ruamel's safe YAML class, so the "without a safe loader" premise fails.

### [ ] Finding `7` — `BP-PY-1`

- Function context: `scripts/rendercv/findings/functions/7.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/ats_proof/analyze_pdfs.py:143:1`
- Checklist pattern: the exception is recorded into the structured result, so it is handled, not swallowed

Source excerpt:

```
        except Exception as e:
            result["extractors"][ext_name] = {"success": False, "error": str(e)}
```

Why this is a false positive: the handler reports the failure — it stores the exception message in the per-extractor result record — so the "without handling ... swallows failures and hides bugs" condition is not satisfied.

Checklist evidence: rule condition is a broad except "without handling or re-raise"; the shown handler records the exception type/message into the output envelope.

### [ ] Finding `9` — `BP-PY-11`

- Function context: `scripts/rendercv/findings/functions/9.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/ats_proof/common.py:102:16`
- Checklist pattern: ruamel `YAML()` instance (import at line 6), safe `load()`

Source excerpt:

```
from ruamel.yaml import YAML
...
    yaml = YAML()
    with yaml_path.open(encoding="utf-8") as f:
        data = yaml.load(f)
```

Why this is a false positive: same as finding 5 — ruamel's safe YAML constructor cannot execute code; the rule's PyYAML `yaml.load` premise does not apply.

Checklist evidence: `detectBPPY11`'s `yaml.load(` needle matches the ruamel instance method; the "can execute code" condition is false for ruamel.

### [ ] Finding `10` — `CWE-502`

- Function context: `scripts/rendercv/findings/functions/10.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/ats_proof/common.py:102:16`
- Checklist pattern: same construct as finding 9 — ruamel `YAML()` instance

Source excerpt:

```
from ruamel.yaml import YAML
...
    yaml = YAML()
    with yaml_path.open(encoding="utf-8") as f:
        data = yaml.load(f)
```

Why this is a false positive: the load is performed by ruamel's safe constructor, so no unsafe deserialization of untrusted data can occur.

Checklist evidence: CWE-502's "yaml.load without a SafeLoader" condition presupposes PyYAML semantics; the shown receiver is ruamel's safe YAML class.

### [ ] Finding `11` — `BP-PY-11`

- Function context: `scripts/rendercv/findings/functions/11.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/ats_proof/common.py:207:16`
- Checklist pattern: same construct as findings 5/9 — ruamel `YAML()` instance, distinct line

Source excerpt:

```
    yaml = YAML()
    with yaml_path.open(encoding="utf-8") as f:
        data = yaml.load(f)
```

Why this is a false positive: identical to finding 9 — ruamel's `YAML().load` cannot construct arbitrary Python objects.

Checklist evidence: the flagged `yaml` is a ruamel `YAML()` instance, not PyYAML's module-level `yaml.load`.

### [ ] Finding `12` — `BP-PY-46`

- Function context: `scripts/rendercv/findings/functions/12.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/ats_proof/run_all.py:33:5`
- Checklist pattern: print in a standalone CLI script (argparse + `if __name__ == "__main__": main()`), user-facing step output, not library logging

Source excerpt:

```
def run_step(script: str, description: str) -> None:
    """Run a script, exit on failure."""
    print(f"\n{'=' * 60}")  # noqa: T201
    print(f"  {description}")  # noqa: T201
```

Why this is a false positive: `run_all.py` is a script (argparse CLI, `__main__` guard invoking `main()`); the prints are the CLI's progress output, so the rule's "operational logging in non-script modules" condition is not satisfied.

Checklist evidence: the module is a script whose only execution path is `python run_all.py`; `pythonCLIPrintSkipFunc` covers `main` but not `run_step`, yet the semantic condition ("non-script modules") is unmet.

### [ ] Finding `13` — `BP-PY-46`

- Function context: `scripts/rendercv/findings/functions/13.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/ats_proof/run_all.py:34:5`
- Checklist pattern: same construct as finding 12, distinct line

Source excerpt:

```
    print(f"\n{'=' * 60}")  # noqa: T201
    print(f"  {description}")  # noqa: T201
    print(f"{'=' * 60}")  # noqa: T201
```

Why this is a false positive: CLI progress output in a standalone script, same reasoning as finding 12.

Checklist evidence: script module, CLI presentation output, not library debug logging.

### [ ] Finding `14` — `BP-PY-46`

- Function context: `scripts/rendercv/findings/functions/14.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/ats_proof/run_all.py:35:5`
- Checklist pattern: same construct as finding 12, distinct line

Source excerpt:

```
    print(f"\n{'=' * 60}")  # noqa: T201
    print(f"  {description}")  # noqa: T201
    print(f"{'=' * 60}")  # noqa: T201
```

Why this is a false positive: CLI progress output in a standalone script, same reasoning as finding 12.

Checklist evidence: script module, CLI presentation output, not library debug logging.

### [ ] Finding `15` — `CWE-88`

- Function context: `scripts/rendercv/findings/functions/15.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/ats_proof/run_all.py:36:14`
- Checklist pattern: the dynamic argv segment comes from fixed developer-maintained constant tuples, not externally influenced input

Source excerpt:

```
STEPS_LOCAL: list[tuple[str, str]] = [
    ("render_pdfs.py", "Render PDFs across all themes"),
    ("analyze_pdfs.py", "Run structural + extraction analysis"),
]
...
    result = subprocess.run(
        [sys.executable, str(SCRIPT_DIR / script)],
        check=False,
    )
```

Why this is a false positive: `script` is always one of the literal filenames in the module's own constant tuples; no externally influenced value can reach the argv, so the "can become an unintended option" condition is unmet.

Checklist evidence: rule condition requires "externally influenced input"; the argv operand is drawn from the source's own constant list, not from external input.

### [ ] Finding `16` — `BP-PY-46`

- Function context: `scripts/rendercv/findings/functions/16.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/ats_proof/run_all.py:41:9`
- Checklist pattern: same construct as finding 12 — failure message output of a standalone script

Source excerpt:

```
    if result.returncode != 0:
        print(f"\nFAILED: {description}")  # noqa: T201
        sys.exit(result.returncode)
```

Why this is a false positive: the print is the CLI's error report in a standalone script, not library logging.

Checklist evidence: script module, CLI presentation output, not library debug logging.

### [ ] Finding `17` — `BP-PY-1`

- Function context: `scripts/rendercv/findings/functions/17.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/ats_proof/submit_commercial.py:86:1`
- Checklist pattern: the exception is reported (printed with its message) and counted, so it is handled, not swallowed

Source excerpt:

```
            except Exception as e:
                failed += 1
                print(f"  [{i + 1}/{len(pdfs)}] ERROR: {rel} - {e}")  # noqa: T201
```

Why this is a false positive: the handler records the failure (increments the counter and prints the error message), so the "without handling ... swallows failures" condition is not satisfied.

Checklist evidence: rule condition is a broad except "without handling or re-raise"; the shown handler reports the exception type/message and counts the failure.

### [ ] Finding `19` — `CWE-88`

- Function context: `scripts/rendercv/findings/functions/19.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/create_executable.py:35:5`
- Checklist pattern: all argv segments are developer/environment constants (sys.executable, in-process temp path), no external input

Source excerpt:

```
    # Run PyInstaller
    subprocess.run(
        [
            sys.executable,
            "-m",
            "PyInstaller",
            ...
            "bin",
            str(rendercv_file),
        ],
        check=True,
    )
```

Why this is a false positive: `rendercv_file` is a path created in-process inside `tempfile.TemporaryDirectory()` and `sys.executable` is the interpreter; no externally influenced value reaches the argv, so the rule's "externally influenced input" condition is unmet.

Checklist evidence: `argvSegmentLooksDynamic` treats any non-literal as dynamic, but the rule's semantic condition requires externally influenced input; the shown values are environment/in-process constants.

### [ ] Finding `20` — `BP-PY-11`

- Function context: `scripts/rendercv/findings/functions/20.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/rendercv_skill/generate.py:245:12`
- Checklist pattern: explicit `ruamel.yaml.YAML()` instance — safe constructor

Source excerpt:

```
    yaml = ruamel.yaml.YAML()
    yaml.preserve_quotes = True
    data = yaml.load(cv_yaml)
```

Why this is a false positive: the receiver is an explicitly constructed `ruamel.yaml.YAML()` instance whose `load()` cannot execute code; the rule's PyYAML premise is false here.

Checklist evidence: `detectBPPY11`'s `yaml.load(` needle matches the ruamel instance method; "can execute code" is false for ruamel.

### [ ] Finding `21` — `CWE-502`

- Function context: `scripts/rendercv/findings/functions/21.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/rendercv_skill/generate.py:245:12`
- Checklist pattern: same construct as finding 20 — ruamel `YAML()` instance

Source excerpt:

```
    yaml = ruamel.yaml.YAML()
    yaml.preserve_quotes = True
    data = yaml.load(cv_yaml)
```

Why this is a false positive: ruamel's safe constructor cannot deserialize into arbitrary objects, so no untrusted deserialization occurs.

Checklist evidence: CWE-502's "yaml.load without a SafeLoader" premise does not hold for ruamel `YAML()` instances.

### [ ] Finding `22` — `BP-PY-11`

- Function context: `scripts/rendercv/findings/functions/22.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/rendercv_skill/generate.py:283:20`
- Checklist pattern: same construct as finding 20 — ruamel `YAML()` instance, distinct line

Source excerpt:

```
        if yaml_path.exists():
            yaml = ruamel.yaml.YAML()
            data = yaml.load(yaml_path.read_text(encoding="utf-8"))
```

Why this is a false positive: ruamel's safe constructor is used; no arbitrary object construction is possible.

Checklist evidence: `detectBPPY11`'s needle matches the ruamel instance method; "can execute code" is false for ruamel.

### [ ] Finding `23` — `PERF-PY-27`

- Function context: `scripts/rendercv/findings/functions/23.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/rendercv_skill/generate.py:283:25`
- Checklist pattern: the load is inside a `for theme in SKILL_THEMES` loop, each iteration loading a distinct per-theme path — the rule's own "once-per-distinct-path" exemption

Source excerpt:

```
    for theme in SKILL_THEMES:
        if theme == "classic":
            continue
        yaml_path = other_themes_dir / f"{theme}.yaml"
        if yaml_path.exists():
            yaml = ruamel.yaml.YAML()
            data = yaml.load(yaml_path.read_text(encoding="utf-8"))
```

Why this is a false positive: `yaml_path` is recomputed per iteration as a different theme file, so no path is loaded repeatedly; the rule's stated exemption ("once-per-distinct-path loads are not 'repeated same path'") applies, and the detector's loop-variable check misses the f-string-derived path.

Checklist evidence: the rule requires "repeated load of the same filesystem path"; each iteration reads `other_themes_dir/<theme>.yaml` with a distinct `<theme>`.

### [ ] Finding `24` — `PERF-PY-27`

- Function context: `scripts/rendercv/findings/functions/24.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/rendercv_skill/generate.py:293:1`
- Checklist pattern: the load line uses the loop binding `path` (a distinct filesystem path per iteration)

Source excerpt:

```
    model_sources: dict[str, str] = {}
    for name, path in MODEL_SOURCE_FILES.items():
        raw_source = path.read_text(encoding="utf-8")
        model_sources[name] = strip_to_schema(raw_source)
```

Why this is a false positive: each iteration reads a different path from the `MODEL_SOURCE_FILES` mapping, so the rule's own "once-per-distinct-path" exemption applies; the detector only checks the first loop binding (`name`) and misses that the load line uses `path`.

Checklist evidence: `enclosingForLoopVar` returns only the first binding of `for name, path in ...`; the load site uses the second binding, which is exactly the loop-variable usage the rule exempts.

### [ ] Finding `25` — `BP-PY-7`

- Function context: `scripts/rendercv/findings/functions/25.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/update_entry_figures.py:30:15`
- Checklist pattern: third-party library method (`fitz.open`), not the builtin `open()`

Source excerpt:

```
    png_file_name = pdf_file_path.stem
    png_files = []
    pdf = fitz.open(pdf_file_path)  # open the PDF file
    for page in pdf:  # iterate the pages
```

Why this is a false positive: `fitz.open` is PyMuPDF's document API — the rule's "use `with open(...) as f:`" fix is not applicable, and the document's lifecycle is managed by the library (auto-closed on GC).

Checklist evidence: BP-PY-7's condition targets builtin `open()` file handles; the flagged call is a third-party method on a path argument.

### [ ] Finding `26` — `BP-PY-46`

- Function context: `scripts/rendercv/findings/functions/26.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/update_entry_figures.py:115:1`
- Checklist pattern: completion print at the end of a standalone maintenance script

Source excerpt:

```
print("Entry figures generated successfully.")  # NOQA: T201
```

Why this is a false positive: the module is a standalone maintenance script (top-level executable code, run via `python scripts/update_entry_figures.py`); the print is the script's completion message, not library logging.

Checklist evidence: rule condition is "print used for operational logging in non-script modules"; this is a script with no importable module API.

### [ ] Finding `27` — `BP-PY-46`

- Function context: `scripts/rendercv/findings/functions/27.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/update_examples.py:59:1`
- Checklist pattern: same construct as finding 26 — standalone maintenance script completion print

Source excerpt:

```
print("Examples generated successfully.")  # NOQA: T201
```

Why this is a false positive: standalone maintenance script; the print is the script's completion message, not library logging.

Checklist evidence: script module, not a non-script module.

### [ ] Finding `28` — `BP-PY-46`

- Function context: `scripts/rendercv/findings/functions/28.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/update_schema.py:7:1`
- Checklist pattern: same construct as finding 26 — standalone maintenance script completion print

Source excerpt:

```
json_schema_file_path = pathlib.Path(__file__).parent.parent / "schema.json"
generate_json_schema_file(json_schema_file_path)
print("Schema generated successfully.")  # NOQA: T201
```

Why this is a false positive: the module is a 7-line standalone script whose whole body is top-level executable code; the print is the script's completion message.

Checklist evidence: rule condition is "non-script modules"; the shown module is a script.

### [ ] Finding `29` — `BP-PY-46`

- Function context: `scripts/rendercv/findings/functions/29.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/src/rendercv/cli/app.py:40:9`
- Checklist pattern: user-facing version output of the typer CLI callback, not library debug logging

Source excerpt:

```
@app.callback()
def cli_command_no_args(
    ctx: typer.Context,
    version_requested: Annotated[
        bool | None, typer.Option("--version", "-v", help="Show the version")
    ] = None,
):
    warn_if_new_version_is_available()

    if version_requested:
        print(f"RenderCV v{__version__}")
```

Why this is a false positive: the print is the CLI's user-facing `--version` output in the typer command layer; the rule's "operational logging in a non-script module" condition targets library debug prints, not CLI presentation.

Checklist evidence: CLI presentation output in a typer CLI callback (the module is the CLI command layer, not a reusable library).

### [ ] Finding `30` — `BP-PY-46`

- Function context: `scripts/rendercv/findings/functions/30.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/src/rendercv/cli/app.py:43:9`
- Checklist pattern: same construct as finding 29 — CLI help output

Source excerpt:

```
    elif ctx.invoked_subcommand is None:
        # No command was provided, show help
        print(ctx.get_help())
        raise typer.Exit()
```

Why this is a false positive: printing the CLI help text is the command's user-facing behavior, not library debug logging.

Checklist evidence: CLI presentation output, not operational logging in library code.

### [ ] Finding `35` — `BP-PY-40`

- Function context: `scripts/rendercv/findings/functions/35.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/src/rendercv/cli/app.py:123:15`
- Checklist pattern: the thread is constructed with `daemon=True` — the shutdown protocol the rule's fix demands

Source excerpt:

```
    if not cache or (time.time() - cache["last_check"]) >= VERSION_CHECK_TTL_SECONDS:
        thread = threading.Thread(target=fetch_and_cache_latest_version, daemon=True)
        thread.start()
```

Why this is a false positive: the thread is explicitly daemon — the rule's fix says to "avoid fire-and-forget non-daemon threads", and a daemon thread cannot block process exit, so the "clear shutdown protocol" condition is satisfied; the detector only skips `daemon=True` when it appears on the same line as `.start(`.

Checklist evidence: rule condition is a thread started without join or a clear shutdown protocol; `daemon=True` on the construction line is that shutdown protocol.

### [ ] Finding `36` — `BP-PY-46`

- Function context: `scripts/rendercv/findings/functions/36.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/src/rendercv/cli/app.py:130:17`
- Checklist pattern: same construct as finding 29 — user-facing update notice in the CLI

Source excerpt:

```
            if current < latest:
                print(
                    "\n[bold yellow]A new version of RenderCV is available!"
```

Why this is a false positive: the print is user-facing CLI notification output, not library debug logging.

Checklist evidence: CLI presentation output, not operational logging in a non-script module.

### [ ] Finding `38` — `CWE-829`

- Function context: `scripts/rendercv/findings/functions/38.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/src/rendercv/cli/app.py:149:14`
- Checklist pattern: the dynamically imported modules are discovered by globbing the package's own `cli/` directory — package-controlled selection, not an untrusted control sphere

Source excerpt:

```
# Auto import all commands so that they are registered with the app:
cli_folder_path = pathlib.Path(__file__).parent
for file in cli_folder_path.rglob("*_command.py"):
    # Enforce folder structure: ./name_command/name_command.py
    folder_name = file.parent.name  # e.g. "foo_command"
    py_file_name = file.stem  # e.g. "foo_command"

    # Build full module path: <parent_pkg>.foo_command.foo_command
    full_module = f"{__package__}.{folder_name}.{py_file_name}"

    module = importlib.import_module(full_module)
```

Why this is a false positive: the module names derive from the package's own source tree (`rglob` over the installed `rendercv/cli/` directory), so no untrusted control sphere selects what executes; the rule's "request-derived module names" condition is unmet.

Checklist evidence: rule condition requires "a dynamically selected module or file path, allowing an untrusted control sphere to select what executes"; the selection is a deterministic scan of package-controlled files.

### [ ] Finding `39` — `CWE-94`

- Function context: `scripts/rendercv/findings/functions/39.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/src/rendercv/cli/app.py:149:14`
- Checklist pattern: same construct as finding 38 — the imported module names are not externally influenced

Source excerpt:

```
cli_folder_path = pathlib.Path(__file__).parent
for file in cli_folder_path.rglob("*_command.py"):
    folder_name = file.parent.name  # e.g. "foo_command"
    py_file_name = file.stem  # e.g. "foo_command"
    full_module = f"{__package__}.{folder_name}.{py_file_name}"
    module = importlib.import_module(full_module)
```

Why this is a false positive: CWE-94's condition is "externally influenced text" reaching a dynamic-import sink; the module name is composed from the package's own file names, not from external input.

Checklist evidence: rule condition requires externally influenced text; the sink argument is derived from a package-internal directory scan.

### [ ] Finding `40` — `BP-PY-46`

- Function context: `scripts/rendercv/findings/functions/40.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/src/rendercv/cli/create_theme_command/create_theme_command.py:59:5`
- Checklist pattern: CLI command presentation (rich panel), not library logging

Source excerpt:

```
    print(
        rich.panel.Panel(
```

Why this is a false positive: the print renders a user-facing panel for the `create-theme` CLI command, not library debug logging.

Checklist evidence: CLI presentation output in the CLI command layer.

### [ ] Finding `41` — `BP-PY-46`

- Function context: `scripts/rendercv/findings/functions/41.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/src/rendercv/cli/error_handler.py:41:17`
- Checklist pattern: CLI error display to the user, not library logging

Source excerpt:

```
        except RenderCVUserError as e:
            if e.message:
                print(
                    rich.panel.Panel(
```

Why this is a false positive: the print is the CLI's user-facing error panel, the module's whole purpose is presenting errors to the CLI user.

Checklist evidence: CLI presentation output, not operational logging in a non-script module.

### [ ] Finding `42` — `BP-PY-46`

- Function context: `scripts/rendercv/findings/functions/42.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/src/rendercv/cli/new_command/new_command.py:121:5`
- Checklist pattern: CLI command presentation panel, not library logging

Source excerpt:

```
    print(build_creation_panel(input_file_path, created_items, existing_items))
```

Why this is a false positive: the print outputs the creation summary panel for the `new` CLI command, not library debug logging.

Checklist evidence: CLI presentation output in the CLI command layer.

### [ ] Finding `43` — `BP-PY-46`

- Function context: `scripts/rendercv/findings/functions/43.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/src/rendercv/cli/new_command/print_welcome.py:14:5`
- Checklist pattern: welcome banner — the module is a CLI presentation helper

Source excerpt:

```
    print(f"\nWelcome to [dodger_blue3]RenderCV v{__version__}[/dodger_blue3]!\n")
```

Why this is a false positive: the module exists solely to print the CLI welcome screen; the print is the user-facing banner, not library logging.

Checklist evidence: CLI presentation output in a presentation-only module.

### [ ] Finding `44` — `BP-PY-46`

- Function context: `scripts/rendercv/findings/functions/44.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/src/rendercv/cli/new_command/print_welcome.py:32:5`
- Checklist pattern: same construct as finding 43 — welcome links panel

Source excerpt:

```
    print(link_panel)
```

Why this is a false positive: same reasoning as finding 43 — CLI welcome presentation output.

Checklist evidence: CLI presentation output in a presentation-only module.

### [ ] Finding `45` — `PERF-PY-26`

- Function context: `scripts/rendercv/findings/functions/45.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/src/rendercv/cli/render_command/render_command.py:227:1`
- Checklist pattern: the "hot path" trigger is the substring `render` inside identifiers (`run_rendercv`), and the call is a one-time CLI argument parse

Source excerpt:

```
    arguments: BuildRendercvModelArguments = {
        ...
        "dont_generate_pdf": dont_generate_pdf,
        "dont_generate_png": dont_generate_png,
        "overrides": parse_override_arguments(extra_data_model_override_arguments),
    }
```

Why this is a false positive: `parse_override_arguments` runs once per CLI invocation over typer option strings — there is no hot path; the rule's `windowHas(..., "render", ...)` marker is satisfied only by the substring `render` inside `run_rendercv`/`rendercv_model` identifiers.

Checklist evidence: the "expensive decode/parse on a hot path" condition requires a service/render path; the shown call is a single parse of CLI override arguments.

### [ ] Finding `46` — `PERF-PY-26`

- Function context: `scripts/rendercv/findings/functions/46.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/src/rendercv/renderer/templater/connections.py:62:1`
- Checklist pattern: the hot-path marker is the substring `render` inside the parameter name `rendercv_model`; the function performs no decode/parse and is not in a loop

Source excerpt:

```
def parse_connections(rendercv_model: RenderCVModel) -> list[Connection]:
    """Extract contact information from CV model into normalized connection format.
    ...
    connections: list[Connection] = []
    for key in rendercv_model.cv._key_order:
        match key:
```

Why this is a false positive: the flagged def is a plain extraction over an already-validated pydantic model (no decode/parse, no loop); the `parse_` name matches `decodeHotRE` and the window "hot path" check is satisfied only by the substring `render` in every `rendercv_model` reference.

Checklist evidence: rule condition requires "expensive decode/parse on a hot path"; there is no serialized-data decode and the only hot-path marker is a substring of the parameter name.

### [ ] Finding `49` — `PERF-PY-26`

- Function context: `scripts/rendercv/findings/functions/49.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/src/rendercv/renderer/templater/connections.py:200:1`
- Checklist pattern: same substring artifact as finding 46 — call site of `parse_connections`

Source excerpt:

```
def compute_connections_for_typst(rendercv_model: RenderCVModel) -> list[str]:
    ...
    connections = parse_connections(rendercv_model)
```

Why this is a false positive: the call runs once per typst connection formatting pass over already-parsed model data; the window hot-path marker is the `render` substring in `rendercv_model`.

Checklist evidence: no decode/parse of serialized data and no hot path; the marker match is a variable-name substring.

### [ ] Finding `50` — `PERF-PY-26`

- Function context: `scripts/rendercv/findings/functions/50.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/src/rendercv/renderer/templater/connections.py:235:1`
- Checklist pattern: same substring artifact as finding 49 — markdown call site

Source excerpt:

```
def compute_connections_for_markdown(rendercv_model: RenderCVModel) -> list[str]:
    ...
    connections = parse_connections(rendercv_model)
```

Why this is a false positive: same reasoning as finding 49 — one-time formatting call, hot-path marker is the `render` substring in `rendercv_model`.

Checklist evidence: no decode/parse of serialized data and no hot path; the marker match is a variable-name substring.

### [ ] Finding `55` — `PERF-PY-26`

- Function context: `scripts/rendercv/findings/functions/55.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/src/rendercv/schema/rendercv_model_builder.py:186:1`
- Checklist pattern: the hot-path marker is the substring `render` inside `RenderCVModel`; the call runs once in the validation-error path

Source excerpt:

```
        model = RenderCVModel.model_validate(commented_map, context=validation_context)
    except pydantic.ValidationError as e:
        validation_errors = parse_validation_errors(e, commented_map, overlay_sources)
        raise RenderCVUserValidationError(validation_errors) from e
```

Why this is a false positive: the parse runs once, only when pydantic validation fails; the window hot-path marker is satisfied by the `render` substring inside the type name `RenderCVModel`, not by a render/service path.

Checklist evidence: no loop and no service/render handler; the marker match is a type-name substring.

### [ ] Finding `56` — `BP-PY-11`

- Function context: `scripts/rendercv/findings/functions/56.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/src/rendercv/schema/yaml_reader.py:52:40`
- Checklist pattern: module-level `yaml = ruamel.yaml.YAML()` instance — safe constructor

Source excerpt:

```
yaml = ruamel.yaml.YAML()
yaml.Scanner = ScannerNoAlias
...
    yaml_as_dictionary: CommentedMap = yaml.load(file_content)
```

Why this is a false positive: the receiver is a `ruamel.yaml.YAML()` instance (only a scanner/constructor tweak for timestamps is added); ruamel's `load()` cannot construct arbitrary Python objects, so the rule's "can execute code" premise is false.

Checklist evidence: `detectBPPY11`'s `yaml.load(` needle matches the ruamel instance method; "can execute code" is false for ruamel.

### [ ] Finding `57` — `CWE-502`

- Function context: `scripts/rendercv/findings/functions/57.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/src/rendercv/schema/yaml_reader.py:52:40`
- Checklist pattern: same construct as finding 56 — ruamel `YAML()` instance

Source excerpt:

```
yaml = ruamel.yaml.YAML()
yaml.Scanner = ScannerNoAlias
...
    yaml_as_dictionary: CommentedMap = yaml.load(file_content)
```

Why this is a false positive: ruamel's safe constructor cannot deserialize into arbitrary objects; the CWE-502 "yaml.load without a SafeLoader" premise does not hold.

Checklist evidence: the load is performed by ruamel's safe YAML class, not PyYAML's `yaml.load`.

### [ ] Finding `64` — `BP-PY-41`

- Function context: `scripts/rendercv/findings/functions/64.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/tests/cli/render_command/test_run_rendercv.py:110:1`
- Checklist pattern: the test body contains both `pytest.raises` and an `assert`; the detector's line-based body scan terminates early at a dedented line inside a multi-line triple-quoted string

Source excerpt:

```
    def test_template_syntax_error(self, tmp_path):
        os.chdir(tmp_path)
        ...
        yaml_file.write_text(
            """cv:
    name: John Doe
design:
    theme: badtheme
""",
            encoding="utf-8",
        )
        ...
        with pytest.raises(typer.Exit) as exc_info, progress:
            run_rendercv(yaml_file, progress)

        assert exc_info.value.exit_code == 1
```

Why this is a false positive: the test does contain `pytest.raises(typer.Exit)` and `assert exc_info.value.exit_code == 1`; the continuation lines of the triple-quoted string indented at the `def` level (`    name: John Doe`) terminate the body scan (`ind <= defIndent`), so the assertions are never seen (reproduced with a minimal fixture).

Checklist evidence: the rule's condition "side effects without assertions" is false — the body has both `pytest.raises` and `assert`; the miss is caused by multi-line string continuation lines at def indentation.

### [ ] Finding `65` — `BP-PY-40`

- Function context: `scripts/rendercv/findings/functions/65.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/tests/cli/render_command/test_watcher.py:41:23`
- Checklist pattern: the thread is constructed with `daemon=True` (line 39) — clear shutdown protocol

Source excerpt:

```
        watcher_thread = threading.Thread(
            target=watcher.run_function_if_files_change,
            args=([watched_file], tracked_function),
            daemon=True,
        )
        watcher_thread.start()
```

Why this is a false positive: the thread is explicitly daemon, so it cannot prevent process exit; the rule's "avoid fire-and-forget non-daemon threads" condition is not violated — the detector only skips `daemon=True` when it appears on the same line as `.start(`.

Checklist evidence: rule condition is a thread started without join or a clear shutdown protocol; `daemon=True` is that shutdown protocol.

### [ ] Finding `71` — `CWE-94`

- Function context: `scripts/rendercv/findings/functions/71.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/tests/schema/models/cv/test_section.py:59:17`
- Checklist pattern: the eval argument is built from `@pytest.mark.parametrize` literal values in a test — not externally influenced text

Source excerpt:

```
@pytest.mark.parametrize(
    "entry, expected_entry_type, expected_section_type",
    [
        ("education_entry", "EducationEntry", "SectionWithEducationEntries"),
        ...
    ],
)
def test_get_entry_type_name_and_section_model(...):
    entry = request.getfixturevalue(entry)
    entry_type, SectionModel = get_entry_type_name_and_section_model(entry)
    ...
    if entry_type != "TextEntry":
        entry = eval(f"{entry_type}(**entry)")
```

Why this is a false positive: CWE-94's condition is "externally influenced text" reaching a code-generation sink; `entry_type` is one of the literal class names declared in the test's own parametrize list, so no external input reaches the eval.

Checklist evidence: rule condition requires externally influenced text; the eval argument derives from test-controlled literals.

### [ ] Finding `72` — `BP-PY-41`

- Function context: `scripts/rendercv/findings/functions/72.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/tests/schema/models/design/test_design.py:85:1`
- Checklist pattern: the test body ends with two `assert` statements; the body scan terminates at a dedented line inside the triple-quoted `__init__.py` content

Source excerpt:

```
    def test_accepts_custom_theme_with_valid_init_file(self, design_adapter, tmp_path):
        custom_theme_path = tmp_path / "mytheme"
        ...
        init_file.write_text(
            """
from pydantic import BaseModel

class MythemeTheme(BaseModel):
    theme: str
    custom_option: str = "default_value"
""",
            encoding="utf-8",
        )
        ...
        assert design.theme == "mytheme"
        assert design.custom_option == "test_value"
```

Why this is a false positive: the test contains two final `assert` statements; the multi-line string's content lines at indentation 0 (`class MythemeTheme(BaseModel):`) terminate the detector's body scan before the asserts are reached (reproduced with a minimal fixture).

Checklist evidence: the rule's condition "side effects without assertions" is false — the body has two asserts; the miss is the same multi-line-string scan bug as finding 64.

### [ ] Finding `73` — `BP-PY-41`

- Function context: `scripts/rendercv/findings/functions/73.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/tests/test_pyodide.py:72:1`
- Checklist pattern: the test ends with `assert result.returncode == 0`; the body scan terminates at dedented lines inside the triple-quoted JS script

Source excerpt:

```
def test_rendercv_installs_in_pyodide(tmp_path: pathlib.Path, js_runtime: str) -> None:
    ...
    script.write_text(f"""\
import fs from "node:fs";
...
`); 
}})

main().catch(err => {{
    console.error(err.message);
    process.exit(1);
}});
""")
    ...
    result = subprocess.run(
        [js_runtime, str(script)],
        capture_output=True,
        text=True,
        timeout=300,
        check=False,
    )
    assert result.returncode == 0, f"Pyodide install failed:\n{result.stderr}"
```

Why this is a false positive: the test contains `assert result.returncode == 0`; the multi-line f-string's content lines at indentation 0 (`import fs from "node:fs";`, `});`) terminate the detector's body scan before the assertion is reached — the same scan bug as findings 64 and 72.

Checklist evidence: the rule's condition "side effects without assertions" is false — the body ends with an assert; the miss is the multi-line-string scan bug.

## True positives

### CWE-1121 — Excessive McCabe Cyclomatic Complexity

| Finding | Source | Reason |
| --- | --- | --- |
| 1 | `docs/docs_templating.py:72` | `define_env` contains ≥ 12 control-flow tokens (`for entry_name...` loop plus 10+ comprehension `for ... in ...` clauses); the detector's branch count is met. |

### CWE-88 — Argument Injection

| Finding | Source | Reason |
| --- | --- | --- |
| 2 | `scripts/ats_proof/analyze_pdfs.py:34` | `["pdftotext", "-layout", str(pdf_path), "-"]` — a runtime-derived path (`find_rendered_pdfs()` directory scan) is embedded in an argv literal without `--`; the dynamic-segment condition is met. |

### CWE-252 — Unchecked Return Value

| Finding | Source | Reason |
| --- | --- | --- |
| 4 | `scripts/ats_proof/analyze_pdfs.py:66` | `subprocess.run(["pdftotext", "-v"], ..., check=False)` is a standalone statement; the return status is discarded and success is inferred only from the absence of `FileNotFoundError`. |

### CWE-396 — Declaration of Catch for Generic Exception

| Finding | Source | Reason |
| --- | --- | --- |
| 8 | `scripts/ats_proof/analyze_pdfs.py:143` | `except Exception as e:` generic catch that consumes the failure (records it and continues); no re-raise, so the failure condition is not propagated. |
| 18 | `scripts/ats_proof/submit_commercial.py:86` | `except Exception as e:` generic catch that consumes the failure (counts and continues); no re-raise. |

### BP-PY-2 — Except Pass

| Finding | Source | Reason |
| --- | --- | --- |
| 31 | `src/rendercv/cli/app.py:73` | `except (OSError, json.JSONDecodeError, KeyError): pass` — suite is solely `pass`, exactly the rule condition. |
| 33 | `src/rendercv/cli/app.py:87` | `except OSError: pass` — suite is solely `pass`. |
| 37 | `src/rendercv/cli/app.py:135` | `except packaging.version.InvalidVersion: pass` — suite is solely `pass`; the rule has no specific-type exemption. |

### CWE-390 — Detection of Error Condition Without Action

| Finding | Source | Reason |
| --- | --- | --- |
| 32 | `src/rendercv/cli/app.py:73` | except clause whose direct body is `pass` — error detected, no action taken. |

### CWE-1071 — Empty Code Block

| Finding | Source | Reason |
| --- | --- | --- |
| 34 | `src/rendercv/cli/app.py:87` | exception handler silently contains only `pass`; matches the empty-except pattern. |

### CWE-478 — Missing Default Case in Multiple Condition Expression

| Finding | Source | Reason |
| --- | --- | --- |
| 47 | `src/rendercv/renderer/templater/connections.py:78` | `match key:` over `cv._key_order` has 6 cases and no wildcard, and the domain is *not* exhaustively covered (keys such as `name`, `photo`, `sections` can appear in `_key_order` and fall through) — unlike the pycaps exhaustive-coverage false positives. |
| 51 | `src/rendercv/schema/models/cv/social_network.py:86` | `match network:` covers only 8 of the 16 `SocialNetworkName` literal members with no wildcard — the domain is not exhaustively covered. |

### CWE-1124 — Excessively Deep Nesting

| Finding | Source | Reason |
| --- | --- | --- |
| 48 | `src/rendercv/renderer/templater/connections.py:152` | `body = "Google Scholar"` sits under for → match → case → for → if → else → match → case (8 control frames ≥ 6); the nesting threshold is genuinely met. |

### CWE-829 — Inclusion of Functionality from Untrusted Control Sphere

| Finding | Source | Reason |
| --- | --- | --- |
| 52 | `src/rendercv/schema/models/design/design.py:96` | `spec_from_file_location("theme", custom_theme_folder / "__init__.py")` + `exec_module` executes a module selected from the user-supplied `design.theme` value — externally influenced path reaches an execution sink. |

### PERF-PY-26 — Expensive Decode Or Parse On Hot Path

| Finding | Source | Reason |
| --- | --- | --- |
| 53 | `src/rendercv/schema/pydantic_error_handling.py:150` | `parse_plain_pydantic_error(...)` is called directly inside the `for plain_error in all_plain_errors:` loop — the in-loop condition is met. |
| 54 | `src/rendercv/schema/pydantic_error_handling.py:162` | `parse_plain_pydantic_error(...)` inside the nested `for plain_cause_error in ...` loop — in-loop condition met. |

### BP-PY-41 — pytest assert With Side Effects Only

| Finding | Source | Reason |
| --- | --- | --- |
| 58 | `tests/cli/render_command/test_progress_panel.py:117` | `test_respects_quiet_mode` calls `panel.print_progress_panel("Test Title")` with no assertion and no exception-capture construct. |
| 59 | `tests/cli/render_command/test_progress_panel.py:123` | `test_displays_step_without_paths` — same construct, no assertion. |
| 60 | `tests/cli/render_command/test_progress_panel.py:130` | `test_displays_step_with_single_path` — no assertion. |
| 61 | `tests/cli/render_command/test_progress_panel.py:137` | `test_displays_step_with_multiple_paths` — no assertion. |
| 62 | `tests/cli/render_command/test_progress_panel.py:147` | `test_displays_multiple_steps` — no assertion. |
| 63 | `tests/cli/render_command/test_progress_panel.py:159` | `test_handles_empty_steps` — no assertion. |
| 66 | `tests/cli/test_app.py:146` | `test_silently_ignores_os_error` only calls `write_version_cache("2.0.0")` after monkeypatching — no assertion. |
| 67 | `tests/cli/test_error_handler.py:9` | `test_returns_function_result_on_success` only calls the decorated function — no assertion. |
| 68 | `tests/renderer/templater/test_entry_templates_from_input.py:761` | `test_never_crashes` calls `clean_trailing_parts(text)` with no assertion — assertion-free by design. |
| 69 | `tests/renderer/templater/test_markdown_parser.py:63` | `test_never_crashes_on_arbitrary_input` calls `escape_typst_characters(text)` with no assertion. |
| 70 | `tests/renderer/templater/test_markdown_parser.py:248` | `test_never_crashes_on_arbitrary_input` calls `markdown_to_typst(text)` with no assertion. |

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/rendercv/chunks/Chunk_1_25.txt`, `scripts/rendercv/chunks/Chunk_26_50.txt`, `scripts/rendercv/chunks/Chunk_51_73.txt`
- Function evidence: `scripts/rendercv/findings/functions/1.txt` .. `scripts/rendercv/findings/functions/73.txt`
- Validation: `git diff --check` — `pass`

## Post-fix remaining-FP audit (2026-08-02)

## Run metadata (fresh scan)

```yaml
timestamp: 2026-08-02T11:08:51Z
repository: rendercv
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv
branch: main
commit: 1d4b87bc427e4cf61c0ef49623c971b0e2224708
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv
scanner_fix: b5b8fde6e5850be5a4bf6957e90110b146a2e584 (binary rebuilt 2026-08-02 16:29)
chunk_path: scripts/rendercv/chunks
function_context_path: scripts/rendercv/findings/functions
```

## Scan evidence

- Build command: `(pre-built) ./bin/goslop` from fix commit `b5b8fde` (rebuilt 2026-08-02 16:29)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/rendercv/chunks -context-dir scripts/rendercv/findings/functions real-repos/rendercv`
- Findings: `36`
- Chunks reviewed: `scripts/rendercv/chunks/Chunk_1_25.txt`, `scripts/rendercv/chunks/Chunk_26_36.txt`
- Function contexts reviewed: `scripts/rendercv/findings/functions/4.txt`, `6.txt`, `7.txt`, `9.txt`, `10.txt`, `11.txt`, `16.txt`, `18.txt`, `19.txt`, `30.txt`, `36.txt`

## Audit checklist

- [x] Read every assigned chunk under `scripts/rendercv/chunks`.
- [x] Read `scripts/rendercv/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary (fresh run)

Fresh finding IDs do not correspond to old IDs; matching was done by `Source:` path. Every fresh finding maps 1:1 to an audited old finding — there are no new findings.

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 11 | 4, 6, 7, 9, 10, 11, 16, 18, 19, 30, 36 |
| True positive | 25 | 1, 2, 3, 5, 8, 12, 13, 14, 15, 17, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 31, 32, 33, 34, 35 |
| Uncertain | 0 | — |

True positive mapping (fresh → audited TP, by source): 1→1 (CWE-1121 `docs_templating.py:72`), 2→2 (CWE-88 `analyze_pdfs.py:34`), 3→4 (CWE-252 `analyze_pdfs.py:66`), 5→8 (CWE-396 `analyze_pdfs.py:143`), 8→18 (CWE-396 `submit_commercial.py:86`), 12→31 (BP-PY-2 `app.py:73`), 13→32 (CWE-390 `app.py:73`), 14→33 (BP-PY-2 `app.py:87`), 15→34 (CWE-1071 `app.py:87`), 17→37 (BP-PY-2 `app.py:135`), 20→47 (CWE-478 `connections.py:78`), 21→48 (CWE-1124 `connections.py:152`), 22→51 (CWE-478 `social_network.py:86`), 23→52 (CWE-829 `design.py:96`), 24→58 / 25→59 / 26→60 / 27→61 / 28→62 / 29→63 (BP-PY-41 `test_progress_panel.py`), 31→66 (BP-PY-41 `test_app.py:146`), 32→67 (BP-PY-41 `test_error_handler.py:9`), 33→68 (BP-PY-41 `test_entry_templates_from_input.py:761`), 34→69 / 35→70 (BP-PY-41 `test_markdown_parser.py`). All audited TPs still present and re-detected; none over-suppressed.

False positive mapping (fresh → audited FP): 4→7, 6→15, 7→17, 9→19, 10→23, 11→24, 16→35, 18→38, 19→39, 30→65, 36→71. These 11 audited FPs re-appear after the fix; the fix removed the other 35 audited FPs (ruamel `yaml.load` family, BP-PY-46 CLI prints, `fitz.open` BP-PY-7, BP-PY-41 multi-line-string scan bug cases, PERF-PY-26 render-substring artifacts, CWE-88 `submit_commercial` siblings).

## False positives

### [ ] Findings `18`, `19` — `CWE-829` / `CWE-94`

- Function context: `scripts/rendercv/findings/functions/18.txt`, `scripts/rendercv/findings/functions/19.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/src/rendercv/cli/app.py:149:14`
- Checklist pattern: the dynamically imported module names are composed from the package's own installed source tree — package-controlled selection, not an untrusted control sphere

Source excerpt:

```
cli_folder_path = pathlib.Path(__file__).parent
for file in cli_folder_path.rglob("*_command.py"):
    folder_name = file.parent.name  # e.g. "foo_command"
    py_file_name = file.stem  # e.g. "foo_command"

    # Build full module path: <parent_pkg>.foo_command.foo_command
    full_module = f"{__package__}.{folder_name}.{py_file_name}"

    module = importlib.import_module(full_module)
```

Why this is a false positive: the module names derive from a deterministic `rglob` over the package's own `rendercv/cli/` directory, so no untrusted control sphere selects what executes (CWE-829's "request-derived module names" condition unmet) and no externally influenced text reaches the dynamic-import sink (CWE-94's "externally influenced text" condition unmet). Same construct and reasoning as audited FPs 38 and 39.

Checklist evidence: CWE-829 requires "a dynamically selected module or file path, allowing an untrusted control sphere to select what executes"; the selection is a deterministic scan of package-controlled files. CWE-94 requires externally influenced text; the sink argument derives from the package-internal directory scan.

### [ ] Finding `4` — `BP-PY-1`

- Function context: `scripts/rendercv/findings/functions/4.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/ats_proof/analyze_pdfs.py:143:1`
- Checklist pattern: the exception is recorded into the structured result, so it is handled, not swallowed

Source excerpt:

```
        except Exception as e:
            result["extractors"][ext_name] = {"success": False, "error": str(e)}
```

Why this is a false positive: the handler stores the exception message in the per-extractor result record, propagating the failure to the caller — BP-PY-1's "weak handling (pass / bare continue / log-only without re-raise)" condition is unmet. Same construct as audited FP 7; note the same `except` line is correctly flagged as CWE-396 (finding 5), which has no handling exemption.

Checklist evidence: BP-PY-1 condition is a broad `except` with weak handling; the shown handler reports the exception type/message into the output envelope.

### [ ] Finding `6` — `CWE-88`

- Function context: `scripts/rendercv/findings/functions/6.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/ats_proof/run_all.py:36:14`
- Checklist pattern: the dynamic argv segment comes from fixed developer-maintained constant tuples, not externally influenced input

Source excerpt:

```
STEPS_LOCAL: list[tuple[str, str]] = [
    ("render_pdfs.py", "Render PDFs across all themes"),
    ("analyze_pdfs.py", "Run structural + extraction analysis"),
]
...
    result = subprocess.run(
        [sys.executable, str(SCRIPT_DIR / script)],
        check=False,
    )
```

Why this is a false positive: `script` is always one of the literal filenames in the module's own constant tuples (`STEPS_LOCAL`/`STEPS_COMMERCIAL`/`STEPS_REPORT`), so no externally influenced value can reach the argv and become an unintended option. Same construct as audited FP 15.

Checklist evidence: rule condition requires externally influenced input; the argv operand is drawn from the source's own constant list, not from external input.

### [ ] Finding `7` — `BP-PY-1`

- Function context: `scripts/rendercv/findings/functions/7.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/ats_proof/submit_commercial.py:86:1`
- Checklist pattern: the exception is reported (printed with its message) and counted, so it is handled, not swallowed

Source excerpt:

```
            except Exception as e:
                failed += 1
                print(f"  [{i + 1}/{len(pdfs)}] ERROR: {rel} - {e}")  # noqa: T201
```

Why this is a false positive: the handler increments the failure counter and prints the error message — handled, not swallowed; BP-PY-1's weak-handling condition is unmet. Same construct as audited FP 17; the same `except` line is correctly flagged as CWE-396 (finding 8).

Checklist evidence: BP-PY-1 condition is a broad except with weak handling; the shown handler reports the exception and counts the failure.

### [ ] Finding `9` — `CWE-88`

- Function context: `scripts/rendercv/findings/functions/9.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/create_executable.py:35:5`
- Checklist pattern: all argv segments are developer/environment constants (sys.executable, fixed flags, in-process temp path), no external input

Source excerpt:

```
    # Run PyInstaller
    subprocess.run(
        [
            sys.executable,
            "-m",
            "PyInstaller",
            "--onefile",
            ...
            "bin",
            str(rendercv_file),
        ],
        check=True,
    )
```

Why this is a false positive: `rendercv_file` is a path created in-process inside `tempfile.TemporaryDirectory()` and the remaining segments are the interpreter and fixed PyInstaller flags; no externally influenced value reaches the argv. Same construct as audited FP 19.

Checklist evidence: the detector treats any non-literal as dynamic, but the rule's semantic condition requires externally influenced input; the shown values are environment/in-process constants.

### [ ] Finding `10` — `PERF-PY-27`

- Function context: `scripts/rendercv/findings/functions/10.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/rendercv_skill/generate.py:283:25`
- Checklist pattern: the load is inside a `for theme in SKILL_THEMES` loop, each iteration loading a distinct per-theme path — the rule's own "once-per-distinct-path" exemption

Source excerpt:

```
    for theme in SKILL_THEMES:
        if theme == "classic":
            continue
        yaml_path = other_themes_dir / f"{theme}.yaml"
        if yaml_path.exists():
            yaml = ruamel.yaml.YAML()
            data = yaml.load(yaml_path.read_text(encoding="utf-8"))
```

Why this is a false positive: `yaml_path` is recomputed per iteration as a different theme file, so no path is loaded repeatedly; the rule's stated exemption (once-per-distinct-path loads are not "repeated same path") applies, and the detector's loop-variable check misses the f-string-derived path. Same construct as audited FP 23.

Checklist evidence: the rule requires repeated load of the same filesystem path; each iteration reads `other_themes_dir/<theme>.yaml` with a distinct `<theme>`.

### [ ] Finding `11` — `PERF-PY-27`

- Function context: `scripts/rendercv/findings/functions/11.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/scripts/rendercv_skill/generate.py:293:1`
- Checklist pattern: the load line uses the loop binding `path` (a distinct filesystem path per iteration)

Source excerpt:

```
    model_sources: dict[str, str] = {}
    for name, path in MODEL_SOURCE_FILES.items():
        raw_source = path.read_text(encoding="utf-8")
        model_sources[name] = strip_to_schema(raw_source)
```

Why this is a false positive: each iteration reads a different path from the `MODEL_SOURCE_FILES` mapping, so the rule's own once-per-distinct-path exemption applies; the detector only checks the first loop binding (`name`) and misses that the load line uses `path`. Same construct as audited FP 24.

Checklist evidence: `enclosingForLoopVar` returns only the first binding of `for name, path in ...`; the load site uses the second binding, which is exactly the loop-variable usage the rule exempts.

### [ ] Finding `16` — `BP-PY-40`

- Function context: `scripts/rendercv/findings/functions/16.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/src/rendercv/cli/app.py:123:15`
- Checklist pattern: the thread is constructed with `daemon=True` — the shutdown protocol the rule's fix demands

Source excerpt:

```
    if not cache or (time.time() - cache["last_check"]) >= VERSION_CHECK_TTL_SECONDS:
        thread = threading.Thread(target=fetch_and_cache_latest_version, daemon=True)
        thread.start()
```

Why this is a false positive: the thread is explicitly daemon, so it cannot block process exit; the rule's "avoid fire-and-forget non-daemon threads" condition is not violated — the detector only skips `daemon=True` when it appears on the same line as `.start(`. Same construct as audited FP 35.

Checklist evidence: rule condition is a thread started without join or a clear shutdown protocol; `daemon=True` on the construction line is that shutdown protocol.

### [ ] Finding `30` — `BP-PY-40`

- Function context: `scripts/rendercv/findings/functions/30.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/tests/cli/render_command/test_watcher.py:41:23`
- Checklist pattern: same construct as finding 16 — thread constructed with `daemon=True` (line 39)

Source excerpt:

```
        watcher_thread = threading.Thread(
            target=watcher.run_function_if_files_change,
            args=([watched_file], tracked_function),
            daemon=True,
        )
        watcher_thread.start()
```

Why this is a false positive: the thread is explicitly daemon, so it cannot prevent process exit; the detector only skips `daemon=True` when it appears on the same line as `.start(`. Same construct as audited FP 65.

Checklist evidence: rule condition is a thread started without join or a clear shutdown protocol; `daemon=True` is that shutdown protocol.

### [ ] Finding `36` — `CWE-94`

- Function context: `scripts/rendercv/findings/functions/36.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/rendercv/tests/schema/models/cv/test_section.py:59:17`
- Checklist pattern: the eval argument is built from `@pytest.mark.parametrize` literal values in a test — not externally influenced text

Source excerpt:

```
    if entry_type != "TextEntry":
        entry = eval(f"{entry_type}(**entry)")
```

Why this is a false positive: CWE-94's condition is externally influenced text reaching a code-generation sink; `entry_type` is one of the literal class names declared in the test's own parametrize list, so no external input reaches the eval. Same construct as audited FP 71.

Checklist evidence: rule condition requires externally influenced text; the eval argument derives from test-controlled literals.

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/rendercv/chunks/Chunk_1_25.txt`, `scripts/rendercv/chunks/Chunk_26_36.txt`
- Function evidence: `scripts/rendercv/findings/functions/4.txt`, `6.txt`, `7.txt`, `9.txt`, `10.txt`, `11.txt`, `16.txt`, `18.txt`, `19.txt`, `30.txt`, `36.txt`
- Validation: `git diff --check` — `pass`

## Post-fix v2 audit (latest binary)

Run metadata: binary rebuilt ~2026-08-02 17:56 (`make build`); scan: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/rendercv/chunks -context-dir scripts/rendercv/findings/functions real-repos/rendercv`; findings `47`; chunks `Chunk_1_25.txt` + `Chunk_26_47.txt`; repo unchanged (commit `1d4b87b`).

Classification summary (fresh): TP **25** (1, 2, 3, 5, 8, 17, 18, 19, 20, 23, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 42, 43, 44, 45, 46) / FP **22** (4, 6, 7, 9, 10, 11, 12, 13, 14, 15, 16, 21, 22, 24, 25, 26, 27, 28, 29, 30, 41, 47) / U **0**. Every fresh finding matched a prior audit entry 1:1 by `Source:` — no new findings. All 25 audited TPs re-detected. The 11 Mode-A FPs all persist; the latest binary re-flags the 11 BP-PY-46 CLI/script presentation prints that the b5b8fde fix had suppressed (fresh 12–16, 22, 26–30).

## Fix checklist (FP patterns)

| Pattern # | Rule | Trigger shape | Count | Example sources |
| --- | --- | --- | --- | --- |
| 1 | BP-PY-1 | `except Exception as e:` whose body records the error (into result dict, or counter+print) — handled, not swallowed | 2 | `analyze_pdfs.py:143`, `submit_commercial.py:86` |
| 2 | BP-PY-46 | print in standalone script (top-level completion message) or CLI command layer (typer `--version`/help, update notice, rich panels, welcome banner) — user-facing presentation, not library logging | 11 | `update_*.py:7/59/115`, `cli/app.py:40/43/130`, `cli/create_theme_command.py:59`, `cli/error_handler.py:41`, `cli/new_command/new_command.py:121`, `cli/new_command/print_welcome.py:14/32` |
| 3 | CWE-88 | subprocess argv segment that looks dynamic but is a developer constant (module-level tuple member, `sys.executable`, in-process temp path) | 2 | `scripts/ats_proof/run_all.py:36`, `scripts/create_executable.py:35` |
| 4 | PERF-PY-27 | `read_text`/`load` inside a loop where the path is derived from the loop binding (`f"{theme}.yaml"` or 2nd binding `path`) — distinct path per iteration | 2 | `scripts/rendercv_skill/generate.py:283/293` |
| 5 | BP-PY-40 | `threading.Thread(..., daemon=True)` with `.start()` on a later line — daemon flag is the shutdown protocol; skip daemon across the multi-line construction, not only same line as `.start(` | 2 | `cli/app.py:123`, `tests/cli/render_command/test_watcher.py:41` |
| 6 | CWE-829 / CWE-94 | `importlib.import_module` of names built from a deterministic `rglob` over the package's own `cli/` tree — package-controlled selection, no untrusted control sphere / no externally influenced text | 2 | `cli/app.py:149` |
| 7 | CWE-94 | `eval(f"{entry_type}(**entry)")` where `entry_type` is a literal class name from the test's own `@pytest.mark.parametrize` list | 1 | `tests/schema/models/cv/test_section.py:59` |

## New findings

None — all 47 fresh findings match a prior classification by `Source:` (TP → prior TP, FP → prior FP with reused reason); 22 FP / 25 TP / 0 uncertain, no fresh-only constructs.
