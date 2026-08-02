# Batch 1 — BP-PY-46 false-positive reduction (2026-08-02)

## Scope

- Rule: **BP-PY-46** only (`detectBPPY46` in `internal/lang/python/detectors/bad_practices/rules_observability.go`)
- Helpers: `isPythonScriptModule` / `isPythonShebangScript` in `common.go`
- Checklist reports sampled (BP-PY-46 FPs): pdf_oxide, httpmorph, caniscrape, movielite, voicetag, sync-with-uv, FlashySurf, whatsapp-wrapped, WHEN-Language, astroz, pytogether, tenso, wse, rendercv, Project_Parva, logxide
- Cross-repo root cause: MASTER.md §1

## FP mechanism (from function context)

Detector flagged `print(` in non-library contexts because it lacked:

1. **Script-path exemption** — `examples/`, `scripts/`, `tools/`, `demos/`, `commands/`, `setup.py`, `cli.py`, `__main__.py` (pdf_oxide/httpmorph/logxide/movielite/tenso/wse/rendercv/astroz setup.py).
2. **CLI decorator coverage** — only Flask `.cli.command(`; missed Click `@cli.command`, Typer `@app.command()`, Cyclopts `@app.default()` (caniscrape, voicetag, sync-with-uv, rendercv).
3. **Docstring / string masking** — line scanner treated doctest/docstring/template `print(` as code (tenso, astroz, Project_Parva, pytogether); argparse epilog column-0 lines (`Examples:`) reset `cliIndent` so `main()` prints fired (whatsapp-wrapped).
4. **Shebang entry scripts** — whole-file CLI like WHEN-Language `when.py`.
5. **Rich CLI UX** — `from rich import print` presentation modules (caniscrape telemetry/upload/diff).

Library TPs (e.g. httpmorph `src/httpmorph/_async_client.py`, pytogether middleware/settings prints, whatsapp-wrapped helper prints outside `main`) must still fire.

## Detector conditions changed

| Guardrail | Behavior |
| --- | --- |
| `isPythonScriptModule` | Early-return for path dirs `examples|example|scripts|script|tools|tool|demos|demo|commands` and basenames `setup.py`, `__main__.py`, `cli.py` |
| `isPythonShebangScript` | Early-return when source starts with `#!` |
| `pythonUsesRichPrint` | Early-return when module does `from rich import print` |
| `isPythonCLIDecorator` | Recognize Click/Typer/Cyclopts/Flask: `.cli.command(`, `click.command/group(`, `cli.command/group(`, `.command(`, `.default(`, `.callback(` |
| Masked twin lines (`pytext.Mask`) | Skip wholly blanked lines for indent tracking; only flag `print(` that survives masking |
| `pythonCLIPrintSkipFunc` | Also skip `show_*` when `__main__` invokes `main`; `run_*` with argparse; `*_command` when Clickish CLI imports present |

## Findings / repos sampled

| Repo | Sample finding IDs | Pattern | Outcome after fix |
| --- | --- | ---: | --- |
| voicetag | 1–29 | Typer `@app.command` | **0** BP-PY-46 (was 29) |
| sync-with-uv | 1–2, 5–15 | cyclopts `@app.default` + `scripts/` | **0** |
| caniscrape | 34+ (CLI cluster) | `@cli.command` + `commands/` + rich print | **12** residual (was ~255); `cli.py`/telemetry/commands clean |
| whatsapp-wrapped | 22–35, 37 | `main()` + epilog indent bug | main prints gone; **16** helper TPs retained |
| httpmorph | 225+ FP / **564 TP** | examples vs `src/` | examples **0**; src TP **1** retained |
| pdf_oxide | 6–9+ | `examples/` | examples **0** |
| pytogether | 12 FP | string template | template silent; **29** library TPs retained |
| FlashySurf | 1,2,4,5 | root hyphen scripts, no shebang/path | **4** residual (no safe path/shebang signal) |
| WHEN-Language | 86+ | shebang `when.py` | `when.py` clean; residual in `interpreter.py` |
| logxide / tenso / wse / rendercv / movielite / astroz / Project_Parva | various | examples/scripts/docstrings/setup | covered by path + mask + setup.py |

## Fixtures (safe + vulnerable pairs)

All under `tests/fixtures/python/bp/`:

- `BP-PY-46-script-path-{safe,vulnerable}.txt`
- `BP-PY-46-click-cli-{safe,vulnerable}.txt`
- `BP-PY-46-typer-cli-{safe,vulnerable}.txt`
- `BP-PY-46-cyclopts-cli-{safe,vulnerable}.txt`
- `BP-PY-46-docstring-print-{safe,vulnerable}.txt`
- `BP-PY-46-main-epilog-{safe,vulnerable}.txt`
- `BP-PY-46-shebang-script-{safe,vulnerable}.txt`
- `BP-PY-46-setup-script-{safe,vulnerable}.txt`
- `BP-PY-46-string-template-{safe,vulnerable}.txt`
- `BP-PY-46-cli-module-{safe,vulnerable}.txt`
- `BP-PY-46-commands-path-{safe,vulnerable}.txt`
- `BP-PY-46-rich-print-{safe,vulnerable}.txt`

Wired in:

- `rules_observability_test.go` (`TestBPPY46PrintInLibrary`)
- `audit_variants_test.go` (append-only BP-PY-46 cases)
- Integration matrix auto-discovers pairs (`TestPythonBPFixturesMatrix`)

## Test results

```
go test ./internal/lang/python/detectors/bad_practices/ -count=1 \
  -run 'TestBPPY46PrintInLibrary|TestBPFalsePositiveAuditFixtureVariants/BP-PY-46'
→ PASS
```

Full package / full python integration currently also surface **unrelated** parallel-agent failures (`BP-PY-49-fingerprint-pin` safe fixture). No BP-PY-46 matrix failures.

## Rescan (post `go build -o bin/goslop`)

| Target | BP-PY-46 count |
| --- | ---: |
| voicetag | 0 |
| sync-with-uv | 0 |
| caniscrape | 12 |
| whatsapp-wrapped | 16 (helper TPs) |
| httpmorph/examples | 0 |
| httpmorph/src | **1** (TP retained) |
| pdf_oxide/examples | 0 |
| pytogether | 29 (library TPs; string-template FP gone) |
| FlashySurf | 4 (residual) |

## Remaining uncertainty

1. **FlashySurf** root scripts (`data-process.py`, `semantic-classification.py`) — no shebang, not under script dirs; hyphenated basenames were not exempted (would break fixture `file:` paths).
2. **caniscrape** 12 leftovers in `analyzers/*`, `utils/*`, `config.py` — audit called them CLI presentation, but they lack decorators/rich-print/commands-path; skipped as too broad without stronger evidence.
3. **WHEN-Language `interpreter.py`** prints — not in the when.py FP cluster; left alone (possible TPs).
4. **`from rich import print` whole-module skip** — correct for audited CLIs; a library that rebinds rich.print for debug would be silenced (rare; note for follow-up).
5. Parallel agents may still be editing shared `audit_variants_test.go` / BP matrix inventory thresholds.
