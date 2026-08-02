# False-positive audit — voicetag

## Run metadata

```yaml
timestamp: 2026-08-02T07:22:58Z
repository: voicetag
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag
branch: main
commit: d5ddf73a2ceb644674f7091a7efa7e9c29e6d621
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag
chunk_path: scripts/voicetag/chunks
function_context_path: scripts/voicetag/findings/functions
```

## Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/voicetag/chunks -context-dir scripts/voicetag/findings/functions real-repos/voicetag`
- Findings: `42`
- Chunks reviewed: `scripts/voicetag/chunks/Chunk_1_25.txt`, `scripts/voicetag/chunks/Chunk_26_42.txt`
- Function contexts reviewed: `scripts/voicetag/findings/functions/1.txt` … `42.txt` (all)

## Audit checklist

- [x] Read every assigned chunk under `scripts/voicetag/chunks`.
- [x] Read `scripts/voicetag/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 29 | 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29 |
| True positive | 13 | 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42 |
| Uncertain | 0 | — |

## False positives

All false positives share one checklist pattern: the flagged call is an intentional, user-facing output of a Typer CLI command function in a script module (`voicetag/cli.py`) that is guarded by `if __name__ == "__main__": main()`. BP-PY-46's condition is "print used for operational logging in non-script modules", and its own fix permits print for CLIs ("keep print under `if __name__ == "__main__"` for CLIs"). The detector's CLI exemption machinery (`pythonCLIPrintSkipFunc`, `.cli.command(` decorator skip) recognizes argparse/`__main__` and Click-style CLI entrypoints but not Typer's `@app.command()` shape, so every command body is flagged even though none is a library debug print.

### [ ] Finding 1 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/1.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:105:21`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
@app.command()
def enroll(
        ...
        console.print(f"[dim]Loaded existing profiles from {profiles}[/dim]")
```

Why this is a false positive: The print is CLI status output inside the Typer `enroll` command of a CLI script module, not operational logging in a non-script module.

Checklist evidence: the module is a script (Typer app + `if __name__ == "__main__": main()` at cli.py:512-513); the call renders a Rich `[dim]…[/dim]` status line for the user.

### [ ] Finding 2 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/2.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:107:25`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
        except VoiceTagError as exc:
            err_console.print(
                Panel(
```

Why this is a false positive: The call renders an error `Panel` to stderr for the CLI user inside the `enroll` command; it is not a library debug print.

Checklist evidence: `err_console` is `Console(stderr=True)` (cli.py:41) and the output is a styled `Panel` error message for the command's user.

### [ ] Finding 3 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/3.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:134:17`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
@app.command()
def enroll(
        ...
        console.print(
            Panel(
```

Why this is a false positive: It is the "Enrollment complete" success `Panel` the CLI prints for the user at the end of the `enroll` command.

Checklist evidence: the `Panel` title is `[green]Enrollment complete[/green]` (cli.py:139) — CLI presentation output, not debug logging.

### [ ] Finding 4 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/4.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:144:21`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
    except VoiceTagError as exc:
        err_console.print(Panel(str(exc), title="[red]Enrollment Error[/red]", border_style="red"))
        raise typer.Exit(code=1)
```

Why this is a false positive: The error `Panel` is the CLI's user-facing error report for the `enroll` command, followed by a `typer.Exit`.

Checklist evidence: `err_console` (stderr console) renders a styled error panel; no debug intent.

### [ ] Finding 5 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/5.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:187:21`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
@app.command()
def identify(
        ...
            console.print(f"[dim]Loaded profiles from {profiles}[/dim]")
```

Why this is a false positive: Status line of the Typer `identify` command telling the user profiles were loaded.

Checklist evidence: Rich `[dim]` status formatting inside a command function of the guarded CLI script.

### [ ] Finding 6 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/6.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:234:17`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
        console.print()
        console.print(table)
```

Why this is a false positive: Blank line separating the CLI's speaker-timeline table output in `identify`.

Checklist evidence: part of the result-table rendering sequence (cli.py:234-236) in a command function.

### [ ] Finding 7 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/7.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:235:17`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
        console.print(table)
```

Why this is a false positive: Prints the speaker-timeline `Table` — the command's primary user output.

Checklist evidence: `table` is a Rich `Table` built for display (cli.py:199-232) inside `identify`.

### [ ] Finding 8 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/8.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:236:17`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
        console.print(table)
        console.print()
```

Why this is a false positive: Trailing blank line closing the timeline output of `identify`.

Checklist evidence: same table-rendering sequence as findings 6-7; presentation only.

### [ ] Finding 9 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/9.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:250:17`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
        console.print(
            Panel(
                "\n".join(summary_lines),
                title="[bold]Summary[/bold]",
```

Why this is a false positive: The `Summary` result `Panel` of the `identify` command.

Checklist evidence: `summary_lines` are user-facing stats (duration, speakers, segments) rendered in a styled panel.

### [ ] Finding 10 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/10.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:263:21`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
            console.print(f"[dim]Results saved to {output}[/dim]")
```

Why this is a false positive: Confirmation line telling the user where JSON results were saved.

Checklist evidence: status output adjacent to the `json.dump` write in `identify`; CLI presentation.

### [ ] Finding 11 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/11.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:266:21`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
    except VoiceTagError as exc:
        err_console.print(Panel(str(exc), title="[red]Error[/red]", border_style="red"))
        raise typer.Exit(code=1)
```

Why this is a false positive: User-facing error `Panel` followed by `typer.Exit` in `identify`.

Checklist evidence: `err_console` stderr rendering of the caught `VoiceTagError`; error presentation.

### [ ] Finding 12 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/12.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:283:17`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
@profiles_app.command("list")
def profiles_list(
        ...
        console.print(
            Panel(
```

Why this is a false positive: The "No Profiles" guidance `Panel` of the `profiles list` command.

Checklist evidence: panel explains how to enroll speakers first (cli.py:283-291); user guidance, not debug output.

### [ ] Finding 13 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/13.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:297:21`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
    except VoiceTagError as exc:
        err_console.print(Panel(str(exc), title="[red]Error[/red]", border_style="red"))
        raise typer.Exit(code=1)
```

Why this is a false positive: User-facing error `Panel` for profile-load failure in `profiles list`.

Checklist evidence: `err_console` stderr error rendering followed by `typer.Exit`; CLI error presentation.

### [ ] Finding 14 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/14.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:302:17`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
        console.print("[yellow]No speakers enrolled.[/yellow]")
        raise typer.Exit(code=0)
```

Why this is a false positive: Direct user message when the profile store is empty in `profiles list`.

Checklist evidence: styled `[yellow]` informational line; presentation output.

### [ ] Finding 15 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/15.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:318:13`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
    console.print(table)
```

Why this is a false positive: Prints the "Enrolled Speakers" `Table` — the `profiles list` command's main output.

Checklist evidence: `table` is a Rich `Table` titled "Enrolled Speakers" (cli.py:305); presentation.

### [ ] Finding 16 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/16.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:319:13`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
    console.print(f"\n[dim]{len(speakers)} speaker(s) enrolled.[/dim]")
```

Why this is a false positive: Footer line of the `profiles list` table output.

Checklist evidence: Rich `[dim]` summary under the table; user-facing.

### [ ] Finding 17 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/17.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:335:21`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
@profiles_app.command("remove")
def profiles_remove(
        ...
        err_console.print(
            Panel(
```

Why this is a false positive: Error `Panel` telling the user the profiles file was not found in `profiles remove`.

Checklist evidence: `err_console` stderr rendering of a styled error panel; error presentation.

### [ ] Finding 18 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/18.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:350:17`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
        console.print(
            Panel(
                f"Removed speaker [cyan]{name}[/cyan] from profiles.\n"
```

Why this is a false positive: Success `Panel` confirming the speaker removal in `profiles remove`.

Checklist evidence: `Panel` title `[green]Speaker Removed[/green]` (cli.py:354); user confirmation.

### [ ] Finding 19 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/19.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:359:21`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
    except VoiceTagError as exc:
        err_console.print(Panel(str(exc), title="[red]Error[/red]", border_style="red"))
        raise typer.Exit(code=1)
```

Why this is a false positive: User-facing error `Panel` for failures in `profiles remove`.

Checklist evidence: `err_console` stderr error rendering followed by `typer.Exit`; error presentation.

### [ ] Finding 20 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/20.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:411:21`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
@app.command()
def transcribe(
        ...
            console.print(f"[dim]Loaded profiles from {profiles}[/dim]")
```

Why this is a false positive: Status line of the `transcribe` command telling the user profiles were loaded.

Checklist evidence: Rich `[dim]` status formatting inside a command function of the guarded CLI script.

### [ ] Finding 21 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/21.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:447:17`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
        console.print()
        console.print(table)
```

Why this is a false positive: Blank line separating the CLI's transcript table output in `transcribe`.

Checklist evidence: part of the result-table rendering sequence (cli.py:447-449) in a command function.

### [ ] Finding 22 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/22.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:448:17`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
        console.print(table)
```

Why this is a false positive: Prints the transcript `Table` — the `transcribe` command's primary user output.

Checklist evidence: `table` is a Rich `Table` titled `Transcript — …` (cli.py:429); presentation.

### [ ] Finding 23 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/23.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:449:17`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
        console.print(table)
        console.print()
```

Why this is a false positive: Trailing blank line closing the transcript output of `transcribe`.

Checklist evidence: same table-rendering sequence as findings 21-22; presentation only.

### [ ] Finding 24 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/24.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:461:17`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
        console.print(
            Panel(
                "\n".join(summary_lines),
                title="[bold]Summary[/bold]",
```

Why this is a false positive: The `Summary` result `Panel` of the `transcribe` command.

Checklist evidence: `summary_lines` are user-facing stats (duration, speakers, provider) rendered in a styled panel.

### [ ] Finding 25 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/25.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:474:21`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
            console.print(f"[dim]Results saved to {output}[/dim]")
```

Why this is a false positive: Confirmation line telling the user where JSON results were saved.

Checklist evidence: status output adjacent to the `json.dump` write in `transcribe`; CLI presentation.

### [ ] Finding 26 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/26.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:477:21`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
    except VoiceTagError as exc:
        err_console.print(Panel(str(exc), title="[red]Error[/red]", border_style="red"))
        raise typer.Exit(code=1)
```

Why this is a false positive: User-facing error `Panel` followed by `typer.Exit` in `transcribe`.

Checklist evidence: `err_console` stderr rendering of the caught `VoiceTagError`; error presentation.

### [ ] Finding 27 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/27.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:493:13`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
@app.command()
def providers() -> None:
    ...
    console.print(table)
```

Why this is a false positive: Prints the "Available STT Providers" `Table` — the `providers` command's entire output.

Checklist evidence: `table` is a Rich `Table` titled "Available STT Providers" (cli.py:487); presentation.

### [ ] Finding 28 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/28.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:494:13`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
    console.print(f"\n[dim]{len(provider_list)} provider(s) available.[/dim]")
```

Why this is a false positive: Footer line of the `providers` command's table output.

Checklist evidence: Rich `[dim]` summary under the table; user-facing.

### [ ] Finding 29 — BP-PY-46

- Function context: `scripts/voicetag/findings/functions/29.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/voicetag/voicetag/cli.py:500:13`
- Checklist pattern: print call is intentional user-facing output in a Typer CLI command function of a script module with a `__main__` guard

Source excerpt:

```
@app.command()
def version() -> None:
    """Show the voicetag version."""
    console.print(f"voicetag [bold cyan]{__version__}[/bold cyan]")
```

Why this is a false positive: The `version` command's sole purpose is printing the version banner to the user.

Checklist evidence: styled `[bold cyan]` version banner inside a Typer command of the guarded CLI script.

## True positives

### CWE-396 — Declaration of Catch for Generic Exception (findings 30, 32, 34, 35, 36, 37, 38, 39, 42)

| Finding id | Source | Reason |
| --- | --- | --- |
| 30 | `voicetag/diarizer.py:70:1` | `except Exception as exc:` matches `pyGenericExceptRE` in a non-test module; the handler inspects the message and folds distinct conditions (auth vs pipeline-load failure) into a single `DiarizationError`. |
| 32 | `voicetag/encoder.py:79:1` | `except Exception as exc:` matches the regex; the handler only logs (`logger.warning`) and continues the loop, so per-file failures are swallowed and hidden. |
| 34 | `voicetag/pipeline.py:157:1` | `except Exception as exc:` matches the regex; the handler logs and returns `None`, hiding the embedding-failure condition from the caller. |
| 35 | `voicetag/providers/deepgram_stt.py:59:1` | `except Exception as exc:` matches the regex; all API/SDK failure conditions are collapsed into `TranscriptionError`. |
| 36 | `voicetag/providers/fireworks_stt.py:61:1` | `except Exception as exc:` matches the regex; HTTP/parse failures are collapsed into `TranscriptionError`. |
| 37 | `voicetag/providers/groq_stt.py:60:1` | `except Exception as exc:` matches the regex; SDK/API failures are collapsed into `TranscriptionError`. |
| 38 | `voicetag/providers/openai_stt.py:59:1` | `except Exception as exc:` matches the regex; SDK/API failures are collapsed into `TranscriptionError`. |
| 39 | `voicetag/providers/whisper_local.py:61:1` | `except Exception as exc:` matches the regex; decode/model failures are collapsed into `TranscriptionError`. |
| 42 | `voicetag/utils.py:72:1` | `except Exception as exc:` matches the regex; distinct decode/read conditions are collapsed into `AudioLoadError`. |

### BP-PY-1 — Bare Except Clause (findings 31, 33)

| Finding id | Source | Reason |
| --- | --- | --- |
| 31 | `voicetag/encoder.py:79:1` | Broad `except Exception as exc:` whose suite only calls `logger.warning(...)` — no `raise`, so `suiteReraises` is false and the rule condition fires (non-test file). |
| 33 | `voicetag/pipeline.py:157:1` | Broad `except Exception as exc:` whose suite only calls `logger.warning(...)` and `return None` — no re-raise; matches the rule condition. |

### CWE-829 — Inclusion of Functionality from Untrusted Control Sphere (finding 40)

| Finding id | Source | Reason |
| --- | --- | --- |
| 40 | `voicetag/transcriber.py:90:18` | `importlib.import_module(module_path)` with a non-literal first argument (`module_path` is a variable, not a literal), so `isDynamicExpr` is true and the dynamic-import condition fires. |

### CWE-94 — Improper Control of Generation of Code (finding 41)

| Finding id | Source | Reason |
| --- | --- | --- |
| 41 | `voicetag/transcriber.py:90:18` | Same call site: a dynamic (non-literal) expression reaches the `importlib.import_module` code-generation/dynamic-import sink, satisfying `detectCWE94`'s non-literal-first-arg condition. |

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/voicetag/chunks`
- Function evidence: `scripts/voicetag/findings/functions`
- Validation: `git diff --check` — `pass`
