# False-positive audit — httptap

## Run metadata

```yaml
timestamp: 2026-08-02T07:43:18Z
repository: httptap
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap
branch: main
commit: 58dc6816ec71a7685b87823a0793b4d0d7cff933
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap
chunk_path: scripts/httptap/chunks
function_context_path: scripts/httptap/findings/functions
```

## Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop` (goslop binary used: `./bin/goslop`)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/httptap/chunks -context-dir scripts/httptap/findings/functions real-repos/httptap`
- Findings: `103`
- Chunks reviewed: `scripts/httptap/chunks/Chunk_1_25.txt`, `scripts/httptap/chunks/Chunk_26_50.txt`, `scripts/httptap/chunks/Chunk_51_75.txt`, `scripts/httptap/chunks/Chunk_76_100.txt`, `scripts/httptap/chunks/Chunk_101_103.txt`
- Function contexts reviewed: `scripts/httptap/findings/functions/1.txt` .. `scripts/httptap/findings/functions/103.txt` (all 103; plus the enclosing source file for every distinct construct)

## Audit checklist

- [x] Read every assigned chunk under `scripts/httptap/chunks`.
- [x] Read `scripts/httptap/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [ ] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 103 | 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103 |
| True positive | 0 | — |
| Uncertain | 0 | — |

## False positives

### [ ] Finding `1` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/1.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/_pkgmeta.py:54:13`
- Checklist pattern: `print(` token inside a docstring doctest example (`>>> print(...)`)

Source excerpt:

```
    Examples:
        >>> info = get_package_info()
        >>> print(f"{info.version} by {info.author}")
        0.1.0 by Sergei Ozeranskii
```

Why this is a false positive: the `print(` token sits inside the `get_package_info` docstring's `Examples:` block (`_pkgmeta.py:43-57`), i.e. documentation text, not an executable statement — no `print` call exists in the code.

Checklist evidence: BP-PY-46's condition is "`print` is used for operational logging in non-script modules"; `printCallOutsideString` matches per line and cannot see the enclosing triple-quoted docstring, so it reports documentation text that is never executed.

### [ ] Finding `2` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/2.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/analyzer.py:144:21`
- Checklist pattern: `print(` token inside a docstring doctest example

Source excerpt:

```
        Examples:
            Basic usage without redirects:
                >>> analyzer = HTTPTapAnalyzer()
                >>> steps = analyzer.analyze_url("https://example.com")
                >>> print(f"Total time: {steps[0].timing.total_ms}ms")
                Total time: 234.5ms
```

Why this is a false positive: the line is inside the `analyze_url` docstring's `Examples:` block (`analyzer.py:140-145`); the `print(` text is a doctest example, not executable code.

Checklist evidence: no executable `print` call exists at the flagged location; the rule's per-line string check cannot see the enclosing docstring.

### [ ] Finding `3` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/3.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/analyzer.py:151:25`
- Checklist pattern: `print(` token inside a docstring doctest example

Source excerpt:

```
            Following redirect chain:
                >>> analyzer = HTTPTapAnalyzer(follow_redirects=True)
                >>> steps = analyzer.analyze_url("http://example.com")
                >>> for i, step in enumerate(steps, 1):
                ...     print(f"Step {i}: {step.response.status}")
                Step 1: 301
                Step 2: 200
```

Why this is a false positive: same construct as finding 2 — a doctest line inside the `analyze_url` docstring, not an executed `print`.

Checklist evidence: the flagged line is documentation inside a triple-quoted docstring, so the rule's "print used for logging" condition is not met.

### [ ] Finding `4` — `CWE-89`

- Function context: `scripts/httptap/findings/functions/4.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/analyzer.py:249:52`
- Checklist pattern: `.execute(...)` is an HTTP request executor method, not a DB-API cursor call

Source excerpt:

```
            outcome: RequestOutcome = self._request.execute(options)
```

Why this is a false positive: `self._request` is a `RequestExecutor` protocol object (defined in `request_executor.py:63`) and `execute(options)` performs an HTTP request with a `RequestOptions` object; no SQL command string is constructed or executed.

Checklist evidence: CWE-89's condition is that `execute/executemany` build SQL via f-strings/`%`/`.format`; the shown source passes a `RequestOptions` object, so there is no SQL command at all.

### [ ] Finding `5` — `BP-PY-1`

- Function context: `scripts/httptap/findings/functions/5.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/analyzer.py:262:1`
- Checklist pattern: broad except whose suite records the failure into the per-step result instead of swallowing it

Source excerpt:

```
        except Exception as exc:  # noqa: BLE001
            # Unexpected error
            step.error = str(exc)
            step.note = f"Step {step_number}: Unexpected error"
```

Why this is a false positive: the exception is not swallowed — it is captured into `step.error`/`step.note`, returned in the `StepMetrics` and rendered to the user (`render.py:150` prints "Step N: ERROR - …"); the analyzer deliberately continues per-step analysis with the failure surfaced.

Checklist evidence: BP-PY-1's condition is "without handling or re-raise swallows failures and hides bugs"; the shown source handles the failure by recording and displaying it, so the condition is not satisfied.

### [ ] Finding `6` — `CWE-396`

- Function context: `scripts/httptap/findings/functions/6.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/analyzer.py:262:1`
- Checklist pattern: same handler as finding 5 — the generic catch records the failure into the step result

Source excerpt:

```
        except Exception as exc:  # noqa: BLE001
            # Unexpected error
            step.error = str(exc)
            step.note = f"Step {step_number}: Unexpected error"
```

Why this is a false positive: the generic handler's only purpose is to convert per-step request failures into visible step errors; distinct failure conditions are intentionally folded into the step result, which the renderer displays.

Checklist evidence: CWE-396's condition is that a generic handler "can hide failures that require distinct handling"; here each step's failure is recorded and shown, so no failure is hidden.

### [ ] Finding `7` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/7.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/cli.py:85:17`
- Checklist pattern: Rich `Console.print` user-facing output in the CLI entry-point module

Source excerpt:

```
    def error(self, message: str) -> NoReturn:
        """Override error to provide Rich formatted error messages."""
        console.print(
            Panel(
```

Why this is a false positive: `cli.py` is the CLI script module (ends with `if __name__ == "__main__": sys.exit(main())`), and this call renders the user-facing argparse error panel — intentional CLI output, not operational logging.

Checklist evidence: BP-PY-46's condition is "`print` is used for operational logging in non-script modules"; `cli.py` is the script module and the call is a Rich `Console.print` producing the CLI's error UI.

### [ ] Finding `8` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/8.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/cli.py:305:17`
- Checklist pattern: Rich `Console.print` user-facing output in the CLI entry-point module

Source excerpt:

```
        console.print("\n[yellow]⚠ Interrupted by user[/yellow]")
        sys.exit(UNIX_SIGNAL_EXIT_OFFSET + signum)  # Standard Unix convention
```

Why this is a false positive: the signal handler prints the CLI's interrupt notice to the user before exiting — intended CLI output in the script module.

Checklist evidence: same as finding 7 — script module, user-facing `Console.print`, not library debug logging.

### [ ] Finding `9` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/9.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/cli.py:349:17`
- Checklist pattern: Rich `Console.print` user-facing output in the CLI entry-point module

Source excerpt:

```
    except OSError as export_error:
        console.print(
            f"[yellow]⚠ Warning:[/yellow] Failed to export JSON: {export_error}",
```

Why this is a false positive: `_export_results` (cli.py:335) warns the user when the `--json` export fails; it is CLI presentation output.

Checklist evidence: the call is a Rich `Console.print` warning inside the CLI script module, not operational logging in library code.

### [ ] Finding `10` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/10.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/cli.py:405:17`
- Checklist pattern: Rich `Console.print` user-facing output in the CLI entry-point module

Source excerpt:

```
        console.print(
            Panel(
                error_text,
                title="[bold red]❌ Validation Error[/bold red]",
```

Why this is a false positive: URL-validation error panel rendered to the user of the CLI; same module as finding 7.

Checklist evidence: user-facing `Console.print` panel in the CLI script module.

### [ ] Finding `11` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/11.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/cli.py:422:17`
- Checklist pattern: Rich `Console.print` user-facing output in the CLI entry-point module

Source excerpt:

```
        console.print(
            Panel(
                error_text,
                title="[bold red]❌ Validation Error[/bold red]",
```

Why this is a false positive: timeout-validation error panel rendered to the CLI user.

Checklist evidence: same as finding 7.

### [ ] Finding `12` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/12.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/cli.py:435:17`
- Checklist pattern: Rich `Console.print` user-facing output in the CLI entry-point module

Source excerpt:

```
    except ValueError as exc:
        console.print(
            Panel(
                str(exc),
                title="[bold red]❌ Header Error[/bold red]",
```

Why this is a false positive: header-parsing error panel rendered to the CLI user.

Checklist evidence: same as finding 7.

### [ ] Finding `13` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/13.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/cli.py:447:21`
- Checklist pattern: Rich `Console.print` user-facing output in the CLI entry-point module

Source excerpt:

```
        if not ca_bundle_str:
            console.print(
                Panel(
```

Why this is a false positive: empty CA-bundle validation error panel rendered to the CLI user.

Checklist evidence: same as finding 7.

### [ ] Finding `14` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/14.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/cli.py:467:21`
- Checklist pattern: Rich `Console.print` user-facing output in the CLI entry-point module

Source excerpt:

```
        except SLOSpecError as exc:
            console.print(
                Panel(
```

Why this is a false positive: SLO-spec parse error panel rendered to the CLI user.

Checklist evidence: same as finding 7.

### [ ] Finding `15` — `BP-PY-1`

- Function context: `scripts/httptap/findings/functions/15.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/cli.py:589:1`
- Checklist pattern: top-level CLI `main()` catch-all that logs the traceback, displays the error and returns a fatal exit code

Source excerpt:

```
    except Exception as e:
        logger.exception("Unexpected error")

        error_panel = Panel(
            f"[red]{e}[/red]",
            title="[bold red]❌ Internal Error[/bold red]",
            border_style="red",
            padding=(1, 2),
        )
        console.print(error_panel)
        return EXIT_FATAL_ERROR
```

Why this is a false positive: the process-boundary handler fully surfaces the failure — `logger.exception` writes the traceback, the error is displayed, and `main` returns `EXIT_FATAL_ERROR` (non-zero) — so nothing is hidden from the user.

Checklist evidence: BP-PY-1's condition is that the handler "swallows failures and hides bugs"; the shown source logs, displays and propagates a fatal exit code instead.

### [ ] Finding `16` — `CWE-396`

- Function context: `scripts/httptap/findings/functions/16.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/cli.py:589:1`
- Checklist pattern: same handler as finding 15 — generic catch at the CLI process boundary that logs and reports

Source excerpt:

```
    except Exception as e:
        logger.exception("Unexpected error")
```

Why this is a false positive: this is the top-level guard of the CLI entry point; the failure is logged with a full traceback, rendered to the user and converted to a fatal exit code, so distinct failure conditions are not hidden.

Checklist evidence: CWE-396's condition ("can hide failures that require distinct handling") is not met because the handler explicitly logs and reports every failure before terminating.

### [ ] Finding `17` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/17.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/exporter.py:171:22`
- Checklist pattern: Rich `Console.print` user-facing success message of the CLI output layer

Source excerpt:

```
    def _print_success(self, output_path: str) -> None:
        """Print success message.
        ...
        """
        self.console.print(f"\n[green]✓ Exported analysis to {output_path}[/green]")
```

Why this is a false positive: `_print_success` is the CLI's confirmation message when `--json` export succeeds — the module is the CLI output layer, and the call is a Rich `Console.print` of intended user output.

Checklist evidence: the flagged call is user-facing product output, not debug/operational logging, so BP-PY-46's condition is unmet.

### [ ] Finding `18` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/18.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/http_client.py:41:13`
- Checklist pattern: `print(` token inside a docstring doctest example

Source excerpt:

```
Examples:
    Basic usage:
        >>> timing, network, response = make_request(
        ...     "https://example.com",
        ...     timeout=10.0,
        ...     http2=True,
        ... )
        >>> print(f"Total: {timing.total_ms:.1f}ms")
        Total: 234.5ms
```

Why this is a false positive: module docstring (`http_client.py:34-44`) doctest example; the `print(` text is documentation, not code.

Checklist evidence: same as finding 1 — no executable `print` call exists at the flagged location.

### [ ] Finding `19` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/19.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/http_client.py:122:17`
- Checklist pattern: `print(` token inside a docstring doctest example

Source excerpt:

```
    Examples:
        >>> try:
        ...     make_request("https://invalid.example", timeout=5.0)
        ... except HTTPClientError as e:
        ...     print(f"Request failed: {e}")
        Request failed: DNS resolution failed: invalid.example
```

Why this is a false positive: doctest block inside the `HTTPClientError` class docstring (`http_client.py:118-123`); not executable code.

Checklist evidence: same as finding 1.

### [ ] Finding `20` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/20.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/http_client.py:461:17`
- Checklist pattern: `print(` token inside a docstring doctest example

Source excerpt:

```
    Examples:
        Basic request:
            >>> timing, network, response = make_request("https://example.com")
            >>> print(f"Status: {response.status}")
            Status: 200
```

Why this is a false positive: doctest example inside the `make_request` docstring (`http_client.py:458-466`); not executable code.

Checklist evidence: same as finding 1.

### [ ] Finding `21` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/21.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/http_client.py:463:17`
- Checklist pattern: `print(` token inside a docstring doctest example

Source excerpt:

```
            >>> print(f"Status: {response.status}")
            Status: 200
            >>> print(f"IP: {network.ip}")
            IP: 93.184.216.34
```

Why this is a false positive: doctest example inside the `make_request` docstring; adjacent to finding 20, distinct line.

Checklist evidence: same as finding 1.

### [ ] Finding `22` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/22.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/http_client.py:465:17`
- Checklist pattern: `print(` token inside a docstring doctest example

Source excerpt:

```
            >>> print(f"IP: {network.ip}")
            IP: 93.184.216.34
            >>> print(f"Total time: {timing.total_ms:.1f}ms")
            Total time: 234.5ms
```

Why this is a false positive: doctest example inside the `make_request` docstring; distinct line from findings 20-21.

Checklist evidence: same as finding 1.

### [ ] Finding `23` — `CWE-93`

- Function context: `scripts/httptap/findings/functions/23.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/http_client.py:576:19`
- Checklist pattern: assignment to an outgoing client request header with a constant value

Source excerpt:

```
            client.headers["User-Agent"] = USER_AGENT
```

Why this is a false positive: `client` is an `httpx.Client` (http_client.py:568) and this sets an *outgoing request* header to the module constant `USER_AGENT` — neither a response header nor an externally influenced value, so no CRLF injection surface exists.

Checklist evidence: CWE-93's condition is writing an externally influenced value into an HTTP *response* header; the shown source sets a request header from a constant.

### [ ] Finding `24` — `BP-PY-2`

- Function context: `scripts/httptap/findings/functions/24.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/http_client.py:616:1`
- Checklist pattern: exception handler deliberately documented as non-fatal best-effort enrichment

Source excerpt:

```
            except TLSInspectionError:
                # TLS inspection is non-fatal, continue without it
                pass
```

Why this is a false positive: the `pass` implements the explicitly documented design — optional TLS metadata enrichment (http_client.py:604-618); when the inspection fails, the request result is still complete and the failure is declared non-fatal in the source.

Checklist evidence: BP-PY-2's condition targets silently discarding failures; the shown source documents and intends the discard because TLS inspection is optional enrichment.

### [ ] Finding `25` — `CWE-390`

- Function context: `scripts/httptap/findings/functions/25.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/http_client.py:616:1`
- Checklist pattern: same handler as finding 24 — documented non-fatal best-effort enrichment

Source excerpt:

```
            except TLSInspectionError:
                # TLS inspection is non-fatal, continue without it
                pass
```

Why this is a false positive: the "error condition without action" is the intended control flow for optional TLS metadata; the program continues with a well-defined state (TLS fields simply remain unset).

Checklist evidence: CWE-390's condition assumes the program "continues without any response" to a real error; here the exception is an expected outcome of an optional probe, documented as non-fatal in the source.

### [ ] Finding `26` — `CWE-1071`

- Function context: `scripts/httptap/findings/functions/26.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/http_client.py:616:13`
- Checklist pattern: same handler as findings 24-25 — empty block is documented non-fatal control flow

Source excerpt:

```
            except TLSInspectionError:
                # TLS inspection is non-fatal, continue without it
                pass
```

Why this is a false positive: the "empty code block" is the deliberate, documented outcome of the non-fatal inspection probe (same construct as findings 24-25, distinct rule).

Checklist evidence: CWE-1071's empty-block heuristic does not apply because the block's emptiness is intentional and documented.

### [ ] Finding `27` — `CWE-396`

- Function context: `scripts/httptap/findings/functions/27.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/http_client.py:628:1`
- Checklist pattern: generic catch that wraps and re-raises as the module's own exception type

Source excerpt:

```
    except HTTPClientError:
        raise
    except Exception as exc:
        msg = f"Unexpected error: {exc}"
        raise HTTPClientError(msg) from exc
```

Why this is a false positive: the suite explicitly re-raises a wrapped `HTTPClientError` with `from exc` (http_client.py:628-630), so no failure is hidden; expected failure modes are handled by the preceding typed handlers.

Checklist evidence: CWE-396's condition ("can hide failures") is not met because the exception is re-raised to the caller.

### [ ] Finding `28` — `BP-PY-1`

- Function context: `scripts/httptap/findings/functions/28.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/http_client.py:655:1`
- Checklist pattern: documented defensive fallback for optional certificate metadata

Source excerpt:

```
        try:
            cert_info = extract_certificate_info(ssl_object)
        except Exception:  # pragma: no cover - defensive  # noqa: BLE001
            cert_info = None
```

Why this is a false positive: the handler is a deliberate defensive fallback — when certificate extraction fails, optional metadata is set to `None` and the request result remains valid; the source marks it `defensive` and `# noqa: BLE001`.

Checklist evidence: BP-PY-1's condition ("swallows failures and hides bugs") is not met: the failure is handled by a defined fallback state and is expected on non-certificate sockets.

### [ ] Finding `29` — `BP-PY-1`

- Function context: `scripts/httptap/findings/functions/29.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/http_client.py:697:1`
- Checklist pattern: documented defensive fallback returning `None`

Source excerpt:

```
    try:
        ssl_candidate = getter("ssl_object")
    except Exception:  # pragma: no cover - defensive  # noqa: BLE001
        return None
```

Why this is a false positive: `_extract_ssl_object` is a best-effort helper (http_client.py:690-699); if the probe attribute is missing the helper returns `None` and the caller continues — the failure is handled by the documented fallback.

Checklist evidence: same as finding 28 — defined fallback handling, not silent swallowing.

### [ ] Finding `30` — `BP-PY-1`

- Function context: `scripts/httptap/findings/functions/30.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/implementations/dns.py:83:1`
- Checklist pattern: worker-thread handler that collects the exception and re-raises it after `join`

Source excerpt:

```
            except Exception as exc:  # pragma: no cover - handled below  # noqa: BLE001
                worker_error = exc
```

Why this is a false positive: the exception is stored in `worker_error` and, after `thread.join(timeout)`, is re-raised to the caller as `DNSResolutionError` (dns.py:94-99); the failure is propagated, not swallowed.

Checklist evidence: BP-PY-1's condition is that failures are swallowed; the shown source collects the error for later re-raise, so the failure surfaces to the caller.

### [ ] Finding `31` — `CWE-396`

- Function context: `scripts/httptap/findings/functions/31.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/implementations/dns.py:83:1`
- Checklist pattern: same handler as finding 30 — exception collected and re-raised after thread join

Source excerpt:

```
            except Exception as exc:  # pragma: no cover - handled below  # noqa: BLE001
                worker_error = exc
```

Why this is a false positive: the generic catch exists only to transport the exception out of the worker thread; the caller re-raises it as `DNSResolutionError`, so failures are not hidden.

Checklist evidence: CWE-396's condition is unmet because the collected exception is propagated after `join` (dns.py:94-99).

### [ ] Finding `32` — `CWE-396`

- Function context: `scripts/httptap/findings/functions/32.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/implementations/tls.py:59:1`
- Checklist pattern: generic catch that wraps and re-raises as a typed exception

Source excerpt:

```
        except Exception as exc:
            msg = f"TLS inspection failed for {host}:{port}: {exc}"
            raise TLSInspectionError(
                msg,
            ) from exc
```

Why this is a false positive: the suite re-raises a wrapped `TLSInspectionError` with the original as cause (tls.py:59-63); the failure is propagated to the caller.

Checklist evidence: CWE-396's condition is unmet because the exception is re-raised, not hidden.

### [ ] Finding `33` — `BP-PY-2`

- Function context: `scripts/httptap/findings/functions/33.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/implementations/tls.py:79:1`
- Checklist pattern: documented best-effort socket-info probe

Source excerpt:

```
        except OSError:  # pragma: no cover - best effort
            pass
```

Why this is a false positive: `_populate_network_info` (tls.py:72-80) treats `getpeername()` failure as expected and marks the probe `# pragma: no cover - best effort`; the IP/family fields simply remain unset and the inspection continues.

Checklist evidence: BP-PY-2's condition targets silently discarding real failures; the shown source documents the `pass` as intentional best-effort enrichment.

### [ ] Finding `34` — `CWE-390`

- Function context: `scripts/httptap/findings/functions/34.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/implementations/tls.py:79:1`
- Checklist pattern: same handler as finding 33 — documented best-effort probe

Source excerpt:

```
        except OSError:  # pragma: no cover - best effort
            pass
```

Why this is a false positive: the `OSError` is an expected outcome of the optional socket probe; continuing without IP metadata is the intended, documented behavior.

Checklist evidence: CWE-390's condition is unmet because the "error condition" is an expected, non-fatal probe result explicitly marked best-effort.

### [ ] Finding `35` — `CWE-1071`

- Function context: `scripts/httptap/findings/functions/35.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/implementations/tls.py:79:9`
- Checklist pattern: same handler as findings 33-34 — empty block is documented best-effort control flow

Source excerpt:

```
        except OSError:  # pragma: no cover - best effort
            pass
```

Why this is a false positive: the empty block is the documented, intentional response to an optional probe failure (same construct as findings 33-34, distinct rule).

Checklist evidence: CWE-1071's empty-block heuristic does not apply to a deliberately empty best-effort handler.

### [ ] Finding `36` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/36.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/interfaces.py:40:21`
- Checklist pattern: `print(` token inside a docstring doctest example

Source excerpt:

```
    Examples:
        >>> class CustomVisualizer:
        ...     def render(self, step: StepMetrics) -> None:
        ...         print(f"Step {step.step_number}: {step.timing.total_ms}ms")
```

Why this is a false positive: doctest example inside the `Visualizer` protocol docstring (interfaces.py:37-41); the `print(` text is documentation, not code.

Checklist evidence: same as finding 1.

### [ ] Finding `37` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/37.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/models.py:24:13`
- Checklist pattern: `print(` token inside a docstring doctest example

Source excerpt:

```
Examples:
    Creating and using metrics:
        >>> timing = TimingMetrics(dns_ms=12.5, connect_ms=45.0)
        >>> timing.calculate_derived()
        >>> print(timing.to_dict())
        {'dns_ms': 12.5, 'connect_ms': 45.0, ...}
```

Why this is a false positive: module docstring doctest example (models.py:20-25); not executable code.

Checklist evidence: same as finding 1.

### [ ] Finding `38` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/38.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/models.py:93:17`
- Checklist pattern: `print(` token inside a docstring doctest example

Source excerpt:

```
            >>> timing.calculate_derived()
            >>> print(f"Wait: {timing.wait_ms}ms, Transfer: {timing.xfer_ms}ms")
            Wait: 20.0ms, Transfer: 50.0ms
```

Why this is a false positive: doctest example inside the `TimingMetrics.calculate_derived` docstring; distinct line from finding 37.

Checklist evidence: same as finding 1.

### [ ] Finding `39` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/39.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/render.py:123:30`
- Checklist pattern: Rich `Console.print` report output of the CLI rendering layer

Source excerpt:

```
        for index, step in enumerate(steps):
            self._render_step(step)
            if index < len(steps) - 1:
                self.console.print(Rule(style="dim"))
                self.console.print()
```

Why this is a false positive: `render.py` is the CLI's output layer — `self.console` is a Rich `Console` and these calls emit the analysis report (separators, panels, metric lines) to the user; they are the product's output, not debug logging.

Checklist evidence: BP-PY-46's condition is "`print` used for operational logging"; the shown calls render the CLI report via the injected Rich `Console`.

### [ ] Finding `40` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/40.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/render.py:124:30`
- Checklist pattern: Rich `Console.print` report output of the CLI rendering layer

Source excerpt:

```
            if index < len(steps) - 1:
                self.console.print(Rule(style="dim"))
                self.console.print()
```

Why this is a false positive: blank-line separator in the report output; same module and purpose as finding 39, distinct line.

Checklist evidence: same as finding 39.

### [ ] Finding `41` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/41.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/render.py:135:22`
- Checklist pattern: Rich `Console.print` report output of the CLI rendering layer

Source excerpt:

```
        if slo_result is None:
            return
        self.console.print()
        self.console.print(format_slo_panel(slo_result))
```

Why this is a false positive: SLO panel spacing/output in `_print_slo_panel` (render.py:131-136); user-facing report rendering.

Checklist evidence: same as finding 39.

### [ ] Finding `42` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/42.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/render.py:136:22`
- Checklist pattern: Rich `Console.print` report output of the CLI rendering layer

Source excerpt:

```
        self.console.print()
        self.console.print(format_slo_panel(slo_result))
```

Why this is a false positive: SLO panel output line; adjacent to finding 41, distinct line.

Checklist evidence: same as finding 39.

### [ ] Finding `43` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/43.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/render.py:150:30`
- Checklist pattern: Rich `Console.print` report output of the CLI rendering layer

Source excerpt:

```
        for step in steps:
            if step.has_error:
                self.console.print(f"Step {step.step_number}: ERROR - {step.error}")
            else:
                self.console.print(format_compact_line(step))
```

Why this is a false positive: compact-mode per-step summary line in `_render_compact` (render.py:138-152); user-facing report output.

Checklist evidence: same as finding 39.

### [ ] Finding `44` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/44.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/render.py:152:30`
- Checklist pattern: Rich `Console.print` report output of the CLI rendering layer

Source excerpt:

```
            else:
                self.console.print(format_compact_line(step))
```

Why this is a false positive: compact-mode step line; adjacent to finding 43, distinct line.

Checklist evidence: same as finding 39.

### [ ] Finding `45` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/45.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/render.py:195:22`
- Checklist pattern: Rich `Console.print` report output of the CLI rendering layer

Source excerpt:

```
            padding=(0, 1),
        )
        self.console.print()
        self.console.print(panel)
        self.console.print()
```

Why this is a false positive: header-panel spacing/output in `_print_header` (render.py:177+); user-facing report rendering.

Checklist evidence: same as finding 39.

### [ ] Finding `46` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/46.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/render.py:196:22`
- Checklist pattern: Rich `Console.print` report output of the CLI rendering layer

Source excerpt:

```
        self.console.print()
        self.console.print(panel)
        self.console.print()
```

Why this is a false positive: header-panel output line; adjacent to findings 45/47, distinct line.

Checklist evidence: same as finding 39.

### [ ] Finding `47` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/47.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/render.py:197:22`
- Checklist pattern: Rich `Console.print` report output of the CLI rendering layer

Source excerpt:

```
        self.console.print(panel)
        self.console.print()
```

Why this is a false positive: trailing spacing after the header panel; distinct line from findings 45-46.

Checklist evidence: same as finding 39.

### [ ] Finding `48` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/48.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/render.py:212:22`
- Checklist pattern: Rich `Console.print` report output of the CLI rendering layer

Source excerpt:

```
            padding=(0, 1),
        )
        self.console.print(header_panel)
```

Why this is a false positive: per-step header panel output in `_render_step`; user-facing report rendering.

Checklist evidence: same as finding 39.

### [ ] Finding `49` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/49.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/render.py:217:26`
- Checklist pattern: Rich `Console.print` report output of the CLI rendering layer

Source excerpt:

```
        if step.has_error:
            error_panel = format_error(step)
            self.console.print(error_panel)
            return
```

Why this is a false positive: per-step error panel output in `_render_step`; user-facing report rendering.

Checklist evidence: same as finding 39.

### [ ] Finding `50` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/50.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/render.py:223:26`
- Checklist pattern: Rich `Console.print` report output of the CLI rendering layer

Source excerpt:

```
        network_line = format_network_info(step)
        if network_line:
            self.console.print(network_line)
```

Why this is a false positive: per-step network info line in `_render_step`; user-facing report rendering.

Checklist evidence: same as finding 39.

### [ ] Finding `51` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/51.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/render.py:228:26`
- Checklist pattern: Rich `Console.print` report output of the CLI rendering layer

Source excerpt:

```
        response_line = format_response_info(step)
        if response_line:
            self.console.print(response_line)
```

Why this is a false positive: per-step response info line; adjacent to finding 50, distinct line.

Checklist evidence: same as finding 39.

### [ ] Finding `52` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/52.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/render.py:254:30`
- Checklist pattern: Rich `Console.print` report output of the CLI rendering layer

Source excerpt:

```
        for step in steps:
            if step.has_error:
                self.console.print(f"Step {step.step_number}: ERROR - {step.error}")
                continue
```

Why this is a false positive: metrics-only mode error line in `_render_metrics_only`; user-facing report output.

Checklist evidence: same as finding 39.

### [ ] Finding `53` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/53.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/render.py:258:26`
- Checklist pattern: Rich `Console.print` report output of the CLI rendering layer

Source excerpt:

```
            step_slo = slo_result if step is slo_target else None
            self.console.print(format_metrics_line(step, slo_result=step_slo))
```

Why this is a false positive: metrics-only mode metrics line; user-facing report output.

Checklist evidence: same as finding 39.

### [ ] Finding `54` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/54.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/render.py:268:22`
- Checklist pattern: Rich `Console.print` report output of the CLI rendering layer

Source excerpt:

```
        table = self._build_redirect_table(steps)
        self.console.print()
        self.console.print(table)
```

Why this is a false positive: redirect-summary spacing in `_render_redirect_summary`; user-facing report rendering.

Checklist evidence: same as finding 39.

### [ ] Finding `55` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/55.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/render.py:269:22`
- Checklist pattern: Rich `Console.print` report output of the CLI rendering layer

Source excerpt:

```
        self.console.print()
        self.console.print(table)
```

Why this is a false positive: redirect table output line; adjacent to finding 54, distinct line.

Checklist evidence: same as finding 39.

### [ ] Finding `56` — `CWE-89`

- Function context: `scripts/httptap/findings/functions/56.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/request_executor.py:63:9`
- Checklist pattern: `execute` is a protocol method performing HTTP requests, not SQL

Source excerpt:

```
    def execute(self, options: RequestOptions) -> RequestOutcome:
        """Perform an HTTP request based on provided options."""
```

Why this is a false positive: `execute` is the `RequestExecutor` protocol's method that performs an HTTP request (http_client.py's `HTTPClientRequestExecutor` implements it); the first argument is a `RequestOptions` object, not SQL.

Checklist evidence: CWE-89's condition requires SQL construction reaching a DB-API `execute`; the shown source is an HTTP request-executor interface.

### [ ] Finding `57` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/57.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/slo.py:18:17`
- Checklist pattern: `print(` token inside a docstring doctest example

Source excerpt:

```
Typical usage:
    >>> thresholds = parse_slo_spec("total=500,ttfb=200")
    >>> result = evaluate_slo(step, thresholds)
    >>> if not result.passed:
    ...     for v in result.violations:
    ...         print(f"{v.key}: {v.actual_ms}ms > {v.threshold_ms}ms")
```

Why this is a false positive: module docstring doctest example (slo.py:13-19); the `print(` text is documentation, not code.

Checklist evidence: same as finding 1.

### [ ] Finding `58` — `CWE-396`

- Function context: `scripts/httptap/findings/functions/58.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/tls_inspector.py:159:1`
- Checklist pattern: generic catch that wraps and re-raises as a typed exception

Source excerpt:

```
    except Exception as e:
        msg = f"Failed to extract certificate info: {e}"
        raise TLSInspectionError(msg) from e
```

Why this is a false positive: the suite re-raises a wrapped `TLSInspectionError` with the original as cause (tls_inspector.py:159-161); the failure is propagated to the caller.

Checklist evidence: CWE-396's condition is unmet because the exception is re-raised, not hidden.

### [ ] Finding `59` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/59.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/tls_inspector.py:178:17`
- Checklist pattern: `print(` token inside a docstring doctest example

Source excerpt:

```
        >>> with context.wrap_socket(sock, server_hostname=host) as tls_sock:
        ...     version, cipher, cert_info = extract_tls_info(tls_sock)
        ...     print(f"TLS: {version}, Cipher: {cipher}")
        TLS: TLSv1.3, Cipher: TLS_AES_256_GCM_SHA384
```

Why this is a false positive: doctest example inside the `extract_tls_info` docstring; not executable code.

Checklist evidence: same as finding 1.

### [ ] Finding `60` — `BP-PY-49`

- Function context: `scripts/httptap/findings/functions/60.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/utils.py:185:23`
- Checklist pattern: TLS verification disabled only in the explicit, operator-chosen legacy mode

Source excerpt:

```
    if verify_ssl:
        context = ssl.create_default_context()
        ...
        return context

    # For legacy mode create a mutable context allowing older protocols.
    context = ssl.SSLContext(ssl.PROTOCOL_TLS)

    context.check_hostname = False
    context.verify_mode = ssl.CERT_NONE
```

Why this is a false positive: `create_ssl_context` (utils.py:157) disables verification only in the `verify_ssl=False` branch, which is reached solely when the operator passes the CLI's `--ignore-ssl` flag (`verify_ssl=not args.ignore_ssl`, cli.py:567); the default path uses `ssl.create_default_context()`. This is the documented, opt-in behavior of a TLS diagnostics tool, not accidental production disabling.

Checklist evidence: BP-PY-49's condition ("HTTP clients disable TLS verification") fires on the raw `CERT_NONE` marker without context; the shown source gates the disablement behind an explicit parameter whose default path keeps full verification.

### [ ] Finding `61` — `CWE-295`

- Function context: `scripts/httptap/findings/functions/61.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/utils.py:185:27`
- Checklist pattern: same construct as finding 60 — verification disabled only in the explicit legacy mode

Source excerpt:

```
    context.check_hostname = False
    context.verify_mode = ssl.CERT_NONE
```

Why this is a false positive: the disabling is conditional on the explicit `verify_ssl` parameter (opt-in via the CLI's `--ignore-ssl` flag); with the default setting the function returns `ssl.create_default_context()` with full verification.

Checklist evidence: CWE-295's condition is met only in the explicitly requested legacy mode; the shown source demonstrates the disablement is parameter-gated and opt-in, so no accidental bypass exists.

### [ ] Finding `62` — `CWE-523`

- Function context: `scripts/httptap/findings/functions/62.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/utils.py:185:27`
- Checklist pattern: same construct as findings 60-61 — opt-in legacy mode; no tool-owned credential flow

Source excerpt:

```
    context.check_hostname = False
    context.verify_mode = ssl.CERT_NONE
```

Why this is a false positive: the tool performs no login or credential-bearing operation of its own — it only relays the operator's explicitly configured request; disabling verification is the operator's opt-in choice for the diagnostic CLI.

Checklist evidence: CWE-523's condition requires a credential-bearing transport; the shown source disables verification only in the operator-chosen legacy mode and the tool itself establishes no credential exchange.

### [ ] Finding `63` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/63.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/visualizer.py:39:22`
- Checklist pattern: Rich `Console.print` timeline output of the CLI visualizer

Source excerpt:

```
        scale = step.timing.total_ms / used_width

        self.console.print("\n  [bold]Request Timeline:[/bold]")
```

Why this is a false positive: `visualizer.py` renders the waterfall timeline for the CLI user via an injected Rich `Console` (visualizer.py:23-25); the call is product output, not debug logging.

Checklist evidence: BP-PY-46's condition is "`print` used for operational logging"; the shown call renders the user-facing timeline.

### [ ] Finding `64` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/64.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/visualizer.py:88:22`
- Checklist pattern: Rich `Console.print` timeline output of the CLI visualizer

Source excerpt:

```
        line = f"  {offset_str}[{color}]{bar}[/{color}] [{color}]{label_field}{timing_str}[/{color}]"
        self.console.print(line)
```

Why this is a false positive: waterfall phase bar line; user-facing timeline rendering.

Checklist evidence: same as finding 63.

### [ ] Finding `65` — `BP-PY-46`

- Function context: `scripts/httptap/findings/functions/65.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/visualizer.py:152:22`
- Checklist pattern: Rich `Console.print` timeline output of the CLI visualizer

Source excerpt:

```
        total_line = f"\n  [bold]Total:[/bold] [bold cyan]{total_str}[/bold cyan] [dim]{scale_info}[/dim]"
        self.console.print(total_line)
```

Why this is a false positive: waterfall total line; user-facing timeline rendering.

Checklist evidence: same as finding 63.

### [ ] Finding `66` — `CWE-89`

- Function context: `scripts/httptap/findings/functions/66.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_analyzer.py:20:9`
- Checklist pattern: `execute` stub mimicking the HTTP request-executor protocol

Source excerpt:

```
class StubExecutor:
    def __init__(self, results: list[tuple[int, str | None]]) -> None:
        ...
    def execute(self, options: RequestOptions) -> RequestOutcome:
        if not self.results:
            msg = "no more results"
            raise HTTPClientError(msg)
```

Why this is a false positive: `StubExecutor.execute` (test_analyzer.py:15-31) fakes the HTTP request executor; the argument is a `RequestOptions` object, and no SQL statement is involved.

Checklist evidence: CWE-89's condition requires SQL construction in `execute`; the shown source is a test stub for an HTTP executor.

### [ ] Finding `67` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/67.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:112:1`
- Checklist pattern: pytest-benchmark micro-benchmark whose `benchmark(...)` fixture call is the test's purpose

Source excerpt:

```
@pytest.mark.benchmark(group="models")
def test_bench_timing_calculate_derived(benchmark: BenchmarkFixture, sample_timing: TimingMetrics) -> None:
    benchmark(sample_timing.calculate_derived)
```

Why this is a false positive: this is a pytest-benchmark performance test — the `benchmark(...)` fixture executes and times the callable and fails the test if it raises; the rule's "placeholder test" condition (a test that passes even if the code breaks) is not satisfied.

Checklist evidence: BP-PY-41's condition targets side-effect-only tests without verification; the benchmark fixture executes the target and errors on failure, so the test does exercise and validate the call.

### [ ] Finding `68` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/68.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:117:1`
- Checklist pattern: same as finding 67 — benchmark test calling `benchmark(sample_timing.to_dict)`

Source excerpt:

```
@pytest.mark.benchmark(group="models")
def test_bench_timing_to_dict(benchmark: BenchmarkFixture, sample_timing: TimingMetrics) -> None:
    benchmark(sample_timing.to_dict)
```

Why this is a false positive: benchmark-only test; the `benchmark(...)` fixture is the designated verification mechanism for pytest-benchmark tests.

Checklist evidence: same as finding 67.

### [ ] Finding `69` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/69.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:122:1`
- Checklist pattern: same as finding 67 — benchmark test calling `benchmark(sample_network.to_dict)`

Source excerpt:

```
@pytest.mark.benchmark(group="models")
def test_bench_network_to_dict(benchmark: BenchmarkFixture, sample_network: NetworkInfo) -> None:
    benchmark(sample_network.to_dict)
```

Why this is a false positive: benchmark-only test, same pattern as finding 67.

Checklist evidence: same as finding 67.

### [ ] Finding `70` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/70.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:127:1`
- Checklist pattern: same as finding 67 — benchmark test calling `benchmark(sample_response.to_dict)`

Source excerpt:

```
@pytest.mark.benchmark(group="models")
def test_bench_response_to_dict(benchmark: BenchmarkFixture, sample_response: ResponseInfo) -> None:
    benchmark(sample_response.to_dict)
```

Why this is a false positive: benchmark-only test, same pattern as finding 67.

Checklist evidence: same as finding 67.

### [ ] Finding `71` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/71.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:132:1`
- Checklist pattern: same as finding 67 — benchmark test calling `benchmark(sample_step.to_dict)`

Source excerpt:

```
@pytest.mark.benchmark(group="models")
def test_bench_step_to_dict(benchmark: BenchmarkFixture, sample_step: StepMetrics) -> None:
    benchmark(sample_step.to_dict)
```

Why this is a false positive: benchmark-only test, same pattern as finding 67.

Checklist evidence: same as finding 67.

### [ ] Finding `72` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/72.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:137:1`
- Checklist pattern: same as finding 67 — benchmark test calling `benchmark(lambda: sample_step.is_redirect)`

Source excerpt:

```
@pytest.mark.benchmark(group="models")
def test_bench_step_is_redirect(benchmark: BenchmarkFixture, sample_step: StepMetrics) -> None:
    benchmark(lambda: sample_step.is_redirect)
```

Why this is a false positive: benchmark-only test, same pattern as finding 67.

Checklist evidence: same as finding 67.

### [ ] Finding `73` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/73.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:142:1`
- Checklist pattern: same as finding 67 — benchmark test calling `benchmark(lambda: sample_step.has_error)`

Source excerpt:

```
@pytest.mark.benchmark(group="models")
def test_bench_step_has_error(benchmark: BenchmarkFixture, sample_step: StepMetrics) -> None:
    benchmark(lambda: sample_step.has_error)
```

Why this is a false positive: benchmark-only test, same pattern as finding 67.

Checklist evidence: same as finding 67.

### [ ] Finding `74` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/74.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:147:1`
- Checklist pattern: same as finding 67 — benchmark test calling `benchmark(format_bytes_human, 512)`

Source excerpt:

```
@pytest.mark.benchmark(group="formatters")
def test_bench_format_bytes_human_small(benchmark: BenchmarkFixture) -> None:
    benchmark(format_bytes_human, 512)
```

Why this is a false positive: benchmark-only test, same pattern as finding 67.

Checklist evidence: same as finding 67.

### [ ] Finding `75` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/75.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:152:1`
- Checklist pattern: same as finding 67 — benchmark test calling `benchmark(format_bytes_human, 1_048_576)`

Source excerpt:

```
@pytest.mark.benchmark(group="formatters")
def test_bench_format_bytes_human_large(benchmark: BenchmarkFixture) -> None:
    benchmark(format_bytes_human, 1_048_576)
```

Why this is a false positive: benchmark-only test, same pattern as finding 67.

Checklist evidence: same as finding 67.

### [ ] Finding `76` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/76.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:157:1`
- Checklist pattern: same as finding 67 — benchmark test calling `benchmark(format_metrics_line, sample_step)`

Source excerpt:

```
@pytest.mark.benchmark(group="formatters")
def test_bench_format_metrics_line(benchmark: BenchmarkFixture, sample_step: StepMetrics) -> None:
    benchmark(format_metrics_line, sample_step)
```

Why this is a false positive: benchmark-only test, same pattern as finding 67.

Checklist evidence: same as finding 67.

### [ ] Finding `77` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/77.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:162:1`
- Checklist pattern: same as finding 67 — benchmark test calling `benchmark(format_network_info, sample_step)`

Source excerpt:

```
@pytest.mark.benchmark(group="formatters")
def test_bench_format_network_info(benchmark: BenchmarkFixture, sample_step: StepMetrics) -> None:
    benchmark(format_network_info, sample_step)
```

Why this is a false positive: benchmark-only test, same pattern as finding 67.

Checklist evidence: same as finding 67.

### [ ] Finding `78` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/78.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:167:1`
- Checklist pattern: same as finding 67 — benchmark test calling `benchmark(format_response_info, sample_step)`

Source excerpt:

```
@pytest.mark.benchmark(group="formatters")
def test_bench_format_response_info(benchmark: BenchmarkFixture, sample_step: StepMetrics) -> None:
    benchmark(format_response_info, sample_step)
```

Why this is a false positive: benchmark-only test, same pattern as finding 67.

Checklist evidence: same as finding 67.

### [ ] Finding `79` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/79.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:172:1`
- Checklist pattern: same as finding 67 — benchmark test calling `benchmark(format_step_header, sample_step)`

Source excerpt:

```
@pytest.mark.benchmark(group="formatters")
def test_bench_format_step_header(benchmark: BenchmarkFixture, sample_step: StepMetrics) -> None:
    benchmark(format_step_header, sample_step)
```

Why this is a false positive: benchmark-only test, same pattern as finding 67.

Checklist evidence: same as finding 67.

### [ ] Finding `80` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/80.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:177:1`
- Checklist pattern: same as finding 67 — benchmark test calling `benchmark(format_error, sample_error_step)`

Source excerpt:

```
@pytest.mark.benchmark(group="formatters")
def test_bench_format_error(benchmark: BenchmarkFixture, sample_error_step: StepMetrics) -> None:
    benchmark(format_error, sample_error_step)
```

Why this is a false positive: benchmark-only test, same pattern as finding 67.

Checklist evidence: same as finding 67.

### [ ] Finding `81` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/81.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:182:1`
- Checklist pattern: same as finding 67 — benchmark test calling `benchmark(mask_sensitive_value, "Bearer …")`

Source excerpt:

```
@pytest.mark.benchmark(group="utils")
def test_bench_mask_sensitive_value(benchmark: BenchmarkFixture) -> None:
    benchmark(mask_sensitive_value, "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.secret")
```

Why this is a false positive: benchmark-only test, same pattern as finding 67.

Checklist evidence: same as finding 67.

### [ ] Finding `82` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/82.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:187:1`
- Checklist pattern: same as finding 67 — benchmark test calling `benchmark(sanitize_headers, headers)`

Source excerpt:

```
@pytest.mark.benchmark(group="utils")
def test_bench_sanitize_headers(benchmark: BenchmarkFixture) -> None:
    headers = {
```

Why this is a false positive: benchmark-only test, same pattern as finding 67 (the body builds sample headers then calls `benchmark(sanitize_headers, headers)`).

Checklist evidence: same as finding 67.

### [ ] Finding `83` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/83.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:199:1`
- Checklist pattern: same as finding 67 — benchmark test calling `benchmark(parse_http_date, "Mon, …")`

Source excerpt:

```
@pytest.mark.benchmark(group="utils")
def test_bench_parse_http_date(benchmark: BenchmarkFixture) -> None:
    benchmark(parse_http_date, "Mon, 22 Oct 2025 12:00:00 GMT")
```

Why this is a false positive: benchmark-only test, same pattern as finding 67.

Checklist evidence: same as finding 67.

### [ ] Finding `84` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/84.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:204:1`
- Checklist pattern: same as finding 67 — benchmark test calling `benchmark(parse_certificate_date, "Oct 22 …")`

Source excerpt:

```
@pytest.mark.benchmark(group="utils")
def test_bench_parse_certificate_date(benchmark: BenchmarkFixture) -> None:
    benchmark(parse_certificate_date, "Oct 22 12:00:00 2025 GMT")
```

Why this is a false positive: benchmark-only test, same pattern as finding 67.

Checklist evidence: same as finding 67.

### [ ] Finding `85` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/85.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:209:1`
- Checklist pattern: same as finding 67 — benchmark test calling `benchmark(validate_url, "https://…")`

Source excerpt:

```
@pytest.mark.benchmark(group="utils")
def test_bench_validate_url_valid(benchmark: BenchmarkFixture) -> None:
    benchmark(validate_url, "https://example.com/api/v1/data?key=value")
```

Why this is a false positive: benchmark-only test, same pattern as finding 67.

Checklist evidence: same as finding 67.

### [ ] Finding `86` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/86.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:214:1`
- Checklist pattern: same as finding 67 — benchmark test calling `benchmark(validate_url, "ftp://…")`

Source excerpt:

```
@pytest.mark.benchmark(group="utils")
def test_bench_validate_url_invalid(benchmark: BenchmarkFixture) -> None:
    benchmark(validate_url, "ftp://example.com/file")
```

Why this is a false positive: benchmark-only test, same pattern as finding 67.

Checklist evidence: same as finding 67.

### [ ] Finding `87` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/87.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:219:1`
- Checklist pattern: same as finding 67 — benchmark test calling `benchmark(calculate_days_until, target)`

Source excerpt:

```
@pytest.mark.benchmark(group="utils")
def test_bench_calculate_days_until(benchmark: BenchmarkFixture) -> None:
    target = datetime.now(timezone.utc) + timedelta(days=120)
```

Why this is a false positive: benchmark-only test, same pattern as finding 67 (body calls `benchmark(calculate_days_until, target)`).

Checklist evidence: same as finding 67.

### [ ] Finding `88` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/88.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:225:1`
- Checklist pattern: same as finding 67 — benchmark test calling `benchmark(build_summary, steps)`

Source excerpt:

```
@pytest.mark.benchmark(group="exporter")
def test_bench_build_summary(benchmark: BenchmarkFixture, sample_step: StepMetrics) -> None:
    steps = [sample_step] * 5
```

Why this is a false positive: benchmark-only test, same pattern as finding 67 (body calls `benchmark(build_summary, steps)`).

Checklist evidence: same as finding 67.

### [ ] Finding `89` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/89.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:231:1`
- Checklist pattern: same as finding 67 — benchmark test calling `benchmark(lambda: json.dumps(sample_step.to_dict()))`

Source excerpt:

```
@pytest.mark.benchmark(group="exporter")
def test_bench_step_to_json(benchmark: BenchmarkFixture, sample_step: StepMetrics) -> None:
    benchmark(lambda: json.dumps(sample_step.to_dict()))
```

Why this is a false positive: benchmark-only test, same pattern as finding 67.

Checklist evidence: same as finding 67.

### [ ] Finding `90` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/90.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:236:1`
- Checklist pattern: same as finding 67 — benchmark test calling `benchmark(WaterfallVisualizer._get_phases, sample_step)`

Source excerpt:

```
@pytest.mark.benchmark(group="visualizer")
def test_bench_get_phases(benchmark: BenchmarkFixture, sample_step: StepMetrics) -> None:
    benchmark(WaterfallVisualizer._get_phases, sample_step)
```

Why this is a false positive: benchmark-only test, same pattern as finding 67.

Checklist evidence: same as finding 67.

### [ ] Finding `91` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/91.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_benchmarks.py:241:1`
- Checklist pattern: same as finding 67 — benchmark test calling `benchmark(visualizer._compute_phase_widths, …)`

Source excerpt:

```
@pytest.mark.benchmark(group="visualizer")
def test_bench_compute_phase_widths(benchmark: BenchmarkFixture) -> None:
    visualizer = WaterfallVisualizer(console=Console(quiet=True))
```

Why this is a false positive: benchmark-only test, same pattern as finding 67 (body calls `benchmark(visualizer._compute_phase_widths, visualizer)`).

Checklist evidence: same as finding 67.

### [ ] Finding `92` — `BP-PY-13`

- Function context: `scripts/httptap/findings/functions/92.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_http_client.py:90:1`
- Checklist pattern: token value generated by `faker` at runtime, not a literal secret

Source excerpt:

```
    url = "https://example.test/api"
    body = b'{"ok": true}'
    token = f"Bearer {faker.hexify('^' * 16)}"
```

Why this is a false positive: the RHS is a runtime f-string calling `faker.hexify(...)` — a generated fake token in a test — not a hardcoded secret literal, and the file is a test module.

Checklist evidence: BP-PY-13's condition is "a secret-like name is assigned a non-empty string literal in source"; the shown source assigns a dynamically generated faker value, so no literal secret exists.

### [ ] Finding `93` — `CWE-617`

- Function context: `scripts/httptap/findings/functions/93.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_http_client.py:93:1`
- Checklist pattern: assertion inside a test mock handler that verifies the test's own expectation

Source excerpt:

```
    def handler(request: httpx.Request) -> httpx.Response:
        assert request.headers["Authorization"] == token
        assert request.headers["Accept"] == "application/json"
```

Why this is a false positive: the `assert` is the test's verification mechanism — the mock transport handler asserts the expected request headers, and the values are controlled by the test itself (test_http_client.py:83-100); this is a test assertion, not a reachable production assertion.

Checklist evidence: CWE-617's condition is a reachable assertion on request-controlled state in production; the shown source is test code where the assert is the intended check.

### [ ] Finding `94` — `CWE-93`

- Function context: `scripts/httptap/findings/functions/94.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_http_client.py:93:23`
- Checklist pattern: assertion *reading* a request header, not writing a response header

Source excerpt:

```
    def handler(request: httpx.Request) -> httpx.Response:
        assert request.headers["Authorization"] == token
```

Why this is a false positive: the line reads the incoming request header inside a test assertion; no header value is written, and the response headers returned are constants.

Checklist evidence: CWE-93's condition requires writing a value into an HTTP response header; the shown source only compares a request header in a test.

### [ ] Finding `95` — `CWE-295`

- Function context: `scripts/httptap/findings/functions/95.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_http_client.py:1073:42`
- Checklist pattern: test assertion verifying the documented `verify_ssl=False` behavior

Source excerpt:

```
        verify_arg = created_clients[0].kwargs["verify"]
        assert isinstance(verify_arg, ssl.SSLContext)
        assert verify_arg.verify_mode == ssl.CERT_NONE
        assert verify_arg.check_hostname is False
```

Why this is a false positive: the line asserts that the opt-in legacy mode produces `CERT_NONE` (test_http_client.py:1056-1075); the assert is test verification of the deliberate parameter-gated behavior, not a production disablement.

Checklist evidence: CWE-295's condition is an explicit production bypass; the shown source is a test assertion verifying documented opt-in behavior.

### [ ] Finding `96` — `CWE-523`

- Function context: `scripts/httptap/findings/functions/96.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_http_client.py:1073:42`
- Checklist pattern: same line as finding 95 — test assertion, no credential-bearing transport established by the test

Source excerpt:

```
        assert verify_arg.verify_mode == ssl.CERT_NONE
```

Why this is a false positive: the test asserts the opt-in context's mode; it does not transmit credentials over an unverified transport.

Checklist evidence: CWE-523's condition requires disabling validation over a credential-bearing transport; the shown source is a behavior assertion in a test.

### [ ] Finding `97` — `CWE-208`

- Function context: `scripts/httptap/findings/functions/97.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_pkgmeta.py:60:12`
- Checklist pattern: equality comparison of package metadata (author string), not a secret

Source excerpt:

```
    assert info.version == version
    assert info.author == primary_author
    assert info.homepage == primary_homepage
```

Why this is a false positive: the compared values are package metadata strings (`version`, `author`, `homepage`, `license`) supplied by the test itself (test_pkgmeta.py:57-62); they are not secrets or authentication values, so timing side-channels are irrelevant.

Checklist evidence: CWE-208's condition is comparing security-sensitive values with `==`; the shown source compares package metadata in a test assertion.

### [ ] Finding `98` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/98.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_property_based.py:117:1`
- Checklist pattern: test that raises `AssertionError` explicitly on failure

Source excerpt:

```
    def test_non_positive_values_rejected(self, value: float) -> None:
        try:
            parse_slo_spec(f"total={value}")
        except SLOSpecError:
            return
        msg = f"expected SLOSpecError for non-positive value {value!r}"
        raise AssertionError(msg)
```

Why this is a false positive: the test verifies rejection by raising `AssertionError` when `SLOSpecError` is not raised; BP-PY-41's assertion detector does not recognize `raise AssertionError`, so it misclassifies a verifying test as assertion-less.

Checklist evidence: BP-PY-41's condition is "tests call production code without assertions"; the shown source contains an explicit `raise AssertionError(msg)` on the failure path, so outcomes are verified.

### [ ] Finding `99` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/99.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_property_based.py:126:1`
- Checklist pattern: same as finding 98 — explicit `raise AssertionError` on failure

Source excerpt:

```
    def test_non_finite_values_rejected(self, literal: str) -> None:
        try:
            parse_slo_spec(f"total={literal}")
        except SLOSpecError:
            return
        msg = f"expected SLOSpecError for non-finite literal {literal!r}"
        raise AssertionError(msg)
```

Why this is a false positive: same construct as finding 98 — the test raises `AssertionError` when the value is not rejected.

Checklist evidence: same as finding 98.

### [ ] Finding `100` — `BP-PY-41`

- Function context: `scripts/httptap/findings/functions/100.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_property_based.py:142:1`
- Checklist pattern: same as finding 98 — explicit `raise AssertionError` on failure

Source excerpt:

```
    def test_unknown_key_rejected(self, key: str) -> None:
        assume("=" not in key and "," not in key and key.strip() != "")
        try:
            parse_slo_spec(f"{key}=100")
        except SLOSpecError:
            return
        msg = f"expected SLOSpecError for unknown key {key!r}"
        raise AssertionError(msg)
```

Why this is a false positive: same construct as finding 98 — the test raises `AssertionError` when an unknown SLO key is accepted.

Checklist evidence: same as finding 98.

### [ ] Finding `101` — `CWE-89`

- Function context: `scripts/httptap/findings/functions/101.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_request_executor.py:44:23`
- Checklist pattern: `.execute` is the HTTP request executor, not a DB cursor

Source excerpt:

```
    executor = HTTPClientRequestExecutor()
    outcome = executor.execute(options)
```

Why this is a false positive: `executor` is an `HTTPClientRequestExecutor` performing an HTTP request with a `RequestOptions` object (test_request_executor.py:43-44); no SQL is constructed or executed.

Checklist evidence: CWE-89's condition requires SQL command construction at a DB-API `execute`; the shown source calls the HTTP executor.

### [ ] Finding `102` — `CWE-295`

- Function context: `scripts/httptap/findings/functions/102.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_utils.py:341:35`
- Checklist pattern: test assertion verifying the documented `verify_ssl=False` behavior

Source excerpt:

```
    def test_create_ssl_context_verification_disabled(self) -> None:
        """Legacy mode disables verification and relaxes protocol bounds."""

        ctx = create_ssl_context(verify_ssl=False)

        assert ctx.verify_mode == ssl.CERT_NONE
        assert ctx.check_hostname is False
```

Why this is a false positive: the line asserts the deliberate opt-in legacy mode of `create_ssl_context`; the test itself verifies documented behavior and makes no production request.

Checklist evidence: CWE-295's condition is an explicit production bypass; the shown source is a unit-test assertion of the parameter-gated behavior.

### [ ] Finding `103` — `CWE-523`

- Function context: `scripts/httptap/findings/functions/103.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_utils.py:341:35`
- Checklist pattern: same line as finding 102 — test assertion, no credential transport

Source excerpt:

```
        assert ctx.verify_mode == ssl.CERT_NONE
        assert ctx.check_hostname is False
```

Why this is a false positive: the test asserts the opt-in context's mode without transporting any credentials.

Checklist evidence: CWE-523's condition requires disabling validation over a credential-bearing transport; the shown source is a behavior assertion in a test.

## True positives

None. All 103 findings were classified as false positives; no finding satisfies its rule condition on the shown source.

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/httptap/chunks`
- Function evidence: `scripts/httptap/findings/functions`
- Validation: `git diff --check` — `pass`

## Post-fix remaining-FP audit (2026-08-02)

### Run metadata

```yaml
timestamp: 2026-08-02T11:10:00Z
repository: httptap
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap
branch: main
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap
chunk_path: scripts/httptap/chunks
function_context_path: scripts/httptap/findings/functions
binary: ./bin/goslop (post-fix rebuild 2026-08-02 16:29)
```

### Scan evidence (fresh run)

- Build command: `go build -o bin/goslop ./cmd/goslop` (binary rebuilt 2026-08-02 16:29)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/httptap/chunks -context-dir scripts/httptap/findings/functions real-repos/httptap` (fresh run 2026-08-02 16:38)
- Findings: `45`
- Chunks reviewed: `scripts/httptap/chunks/Chunk_1_25.txt`, `scripts/httptap/chunks/Chunk_26_45.txt`
- Function contexts reviewed: `scripts/httptap/findings/functions/1.txt` .. `scripts/httptap/findings/functions/45.txt` (all 45)

### Audit checklist

- [x] Read every assigned chunk under `scripts/httptap/chunks`.
- [x] Read `scripts/httptap/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [ ] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

### Classification summary (fresh run)

Matching method: the old audit (above) contains **0 audited true positives**, so no fresh finding can match an audited TP. 44 fresh findings match audited-FP sources exactly (file:line); 1 fresh finding (7) is new but fails its rule condition on the shown source.

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 45 | 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45 |
| True positive | 0 | — |
| Uncertain | 0 | — |

Remaining false positives after fix: `45` (the fix removed 58 of the 103 old findings; every fresh finding is still a false positive).

### False positives

#### [ ] Finding `1` — `BP-PY-46` (re-appearing audited FP, old `17`)

- Function context: `scripts/httptap/findings/functions/1.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/exporter.py:171:22`
- Checklist pattern: Rich `Console.print` user-facing success message of the CLI output layer

```
        self.console.print(f"\n[green]✓ Exported analysis to {output_path}[/green]")
```

Why this is a false positive: `_print_success` (exporter.py:171) is the CLI's confirmation message rendered via the injected Rich `Console` — product output, not operational logging.

Checklist evidence: BP-PY-46's condition is "`print` is used for operational logging in non-script modules"; the shown call is user-facing output of the CLI export layer.

#### [ ] Finding `2` — `CWE-93` (re-appearing audited FP, old `23`)

- Function context: `scripts/httptap/findings/functions/2.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/http_client.py:576:19`
- Checklist pattern: assignment to an outgoing client request header with a constant value

```
            client.headers["User-Agent"] = USER_AGENT
```

Why this is a false positive: `client` is an `httpx.Client` and the line sets an *outgoing request* header to the module constant `USER_AGENT` — no response header, no externally influenced value, so no CRLF injection surface.

Checklist evidence: CWE-93's condition is writing an externally influenced value into an HTTP *response* header; the shown source sets a request header from a constant.

#### [ ] Findings `3`, `4`, `5` — `BP-PY-2`, `CWE-390`, `CWE-1071` (re-appearing audited FPs, old `24`, `25`, `26`)

- Function contexts: `scripts/httptap/findings/functions/3.txt`, `4.txt`, `5.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/http_client.py:616:1` (`:13` for `5`)
- Checklist pattern: documented non-fatal best-effort TLS enrichment handler

```
            except TLSInspectionError:
                # TLS inspection is non-fatal, continue without it
                pass
```

Why this is a false positive: the `pass` is the explicitly documented design for optional TLS metadata enrichment (http_client.py:604-618); the request result stays complete when the probe fails, so no failure is discarded.

Checklist evidence: the handlers are identical source constructs flagged by three rules — the "silently discarded failure" (BP-PY-2/CWE-390) and "empty block" (CWE-1071) conditions are not met because the discard is intentional and documented in the source.

#### [ ] Findings `6`, `7` — `BP-PY-1`, `CWE-396` (`6` re-appearing audited FP, old `28`; `7` is a NEW finding)

- Function contexts: `scripts/httptap/findings/functions/6.txt`, `7.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/http_client.py:655:1`
- Checklist pattern: documented defensive fallback for optional certificate metadata

```
        try:
            cert_info = extract_certificate_info(ssl_object)
        except Exception:  # pragma: no cover - defensive  # noqa: BLE001
            cert_info = None
```

Why this is a false positive: the handler is a deliberate defensive fallback — when certificate extraction fails, optional metadata is set to `None` and the request result remains valid; the source marks it `defensive`/`# noqa: BLE001` (http_client.py:653-656). Finding 7 is new in this run (the old CWE-396 at this file was at line 628) but fails the same condition: no distinct failure handling is required for optional metadata, and the failure is handled by a defined fallback state rather than hidden.

Checklist evidence: BP-PY-1's condition ("swallows failures and hides bugs") and CWE-396's condition ("can hide failures that require distinct handling") are unmet — the exception is expected on non-certificate sockets and is converted to a defined `None` fallback.

#### [ ] Finding `8` — `BP-PY-1` (re-appearing audited FP, old `29`)

- Function context: `scripts/httptap/findings/functions/8.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/http_client.py:697:1`
- Checklist pattern: documented defensive fallback returning `None`

```
    try:
        ssl_candidate = getter("ssl_object")
    except Exception:  # pragma: no cover - defensive  # noqa: BLE001
        return None
```

Why this is a false positive: `_extract_ssl_object` is a best-effort helper (http_client.py:690-699); a missing probe attribute yields `None` and the caller continues — failure handled by documented fallback.

Checklist evidence: BP-PY-1's condition is unmet: the failure is handled by a defined fallback return, not silently swallowed.

#### [ ] Findings `9`, `10` — `BP-PY-1`, `CWE-396` (re-appearing audited FPs, old `30`, `31`)

- Function contexts: `scripts/httptap/findings/functions/9.txt`, `10.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/implementations/dns.py:83:1`
- Checklist pattern: worker-thread handler that collects the exception and re-raises it after `join`

```
            except Exception as exc:  # pragma: no cover - handled below  # noqa: BLE001
                worker_error = exc
```

Why this is a false positive: the exception is stored in `worker_error` and re-raised to the caller as `DNSResolutionError` after `thread.join(timeout)` (dns.py:94-99); the failure is propagated, not swallowed.

Checklist evidence: BP-PY-1/CWE-396 conditions require hidden failures; the shown handler exists only to transport the exception out of the worker thread and the caller re-raises it.

#### [ ] Findings `11`, `12`, `13` — `BP-PY-2`, `CWE-390`, `CWE-1071` (re-appearing audited FPs, old `33`, `34`, `35`)

- Function contexts: `scripts/httptap/findings/functions/11.txt`, `12.txt`, `13.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/implementations/tls.py:79:1` (`:9` for `13`)
- Checklist pattern: documented best-effort socket-info probe

```
        except OSError:  # pragma: no cover - best effort
            pass
```

Why this is a false positive: `_populate_network_info` (tls.py:72-80) treats `getpeername()` failure as expected, marked `# pragma: no cover - best effort`; IP/family fields simply stay unset and inspection continues.

Checklist evidence: identical construct flagged by three rules — BP-PY-2/CWE-390 ("failure discarded silently") and CWE-1071 ("empty block") conditions are unmet because the pass is the documented response to an expected optional-probe failure.

#### [ ] Findings `14`..`30` — `BP-PY-46` (re-appearing audited FPs, old `39`..`55`)

- Function contexts: `scripts/httptap/findings/functions/14.txt` .. `30.txt`
- Sources: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/render.py:123:30, 124:30, 135:22, 136:22, 150:30, 152:30, 195:22, 196:22, 197:22, 212:22, 217:26, 223:26, 228:26, 254:30, 258:26, 268:22, 269:22`
- Checklist pattern: Rich `Console.print` report output of the CLI rendering layer

```
        for index, step in enumerate(steps):
            self._render_step(step)
            if index < len(steps) - 1:
                self.console.print(Rule(style="dim"))
                self.console.print()
```

Why this is a false positive: all 17 calls sit in `render.py`, the CLI's output layer — `self.console` is an injected Rich `Console` and every call emits the user-facing analysis report (separators, panels, step/metric lines); the module is imported and used exclusively by the CLI entry point, so these are product output, not operational logging.

Checklist evidence: BP-PY-46's condition ("print used for operational logging in non-script modules") is unmet — every flagged line renders CLI report output via the Rich `Console`.

#### [ ] Findings `31`, `32`, `33` — `BP-PY-49`, `CWE-295`, `CWE-523` (re-appearing audited FPs, old `60`, `61`, `62`)

- Function contexts: `scripts/httptap/findings/functions/31.txt`, `32.txt`, `33.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/utils.py:185:23` (`:27` for `32`, `33`)
- Checklist pattern: TLS verification disabled only in the explicit, operator-chosen legacy mode

```
    # For legacy mode create a mutable context allowing older protocols.
    context = ssl.SSLContext(ssl.PROTOCOL_TLS)

    context.check_hostname = False
    context.verify_mode = ssl.CERT_NONE
```

Why this is a false positive: `create_ssl_context` (utils.py:157-190) disables verification only in the `verify_ssl=False` branch, reached solely via the CLI's opt-in `--ignore-ssl` flag; the default path returns `ssl.create_default_context()` with full verification, and the tool relays the operator's own requests rather than performing credential exchanges.

Checklist evidence: BP-PY-49/CWE-295 fire on the raw `CERT_NONE` marker; the shown source gates the disablement behind an explicit parameter whose default path keeps verification. CWE-523 additionally requires a credential-bearing transport, which the tool itself never establishes.

#### [ ] Findings `34`, `35`, `36` — `BP-PY-46` (re-appearing audited FPs, old `63`, `64`, `65`)

- Function contexts: `scripts/httptap/findings/functions/34.txt`, `35.txt`, `36.txt`
- Sources: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/httptap/visualizer.py:39:22, 88:22, 152:22`
- Checklist pattern: Rich `Console.print` timeline output of the CLI visualizer

```
        scale = step.timing.total_ms / used_width

        self.console.print("\n  [bold]Request Timeline:[/bold]")
```

Why this is a false positive: `visualizer.py` renders the waterfall timeline for the CLI user via an injected Rich `Console` (visualizer.py:23-25); the calls are user-facing report output, not debug logging.

Checklist evidence: same as findings 14-30 — BP-PY-46's operational-logging condition is unmet.

#### [ ] Finding `37` — `CWE-617` (re-appearing audited FP, old `93`)

- Function context: `scripts/httptap/findings/functions/37.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_http_client.py:93:1`
- Checklist pattern: assertion inside a test mock handler that verifies the test's own expectation

```
    def handler(request: httpx.Request) -> httpx.Response:
        assert request.headers["Authorization"] == token
```

Why this is a false positive: the `assert` is the test's verification mechanism inside a mock transport handler; request values are controlled by the test itself (test_http_client.py:83-100), so this is test code, not a reachable production assertion.

Checklist evidence: CWE-617's condition is a reachable assertion on request-controlled state in production; the shown source is a test where the assert is the intended check.

#### [ ] Findings `38`, `39` — `CWE-295`, `CWE-523` (re-appearing audited FPs, old `95`, `96`)

- Function contexts: `scripts/httptap/findings/functions/38.txt`, `39.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_http_client.py:1073:42`
- Checklist pattern: test assertion verifying the documented `verify_ssl=False` behavior

```
        assert verify_arg.verify_mode == ssl.CERT_NONE
        assert verify_arg.check_hostname is False
```

Why this is a false positive: the lines assert the opt-in legacy context's mode (test_http_client.py:1056-1075); they verify documented behavior and transmit no credentials over an unverified transport.

Checklist evidence: CWE-295 requires an explicit production bypass, CWE-523 a credential-bearing transport — the shown source is a behavior assertion in a test.

#### [ ] Finding `40` — `CWE-208` (re-appearing audited FP, old `97`)

- Function context: `scripts/httptap/findings/functions/40.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_pkgmeta.py:60:12`
- Checklist pattern: equality comparison of package metadata (author string), not a secret

```
    assert info.author == primary_author
```

Why this is a false positive: the compared values are package metadata strings supplied by the test itself (test_pkgmeta.py:57-62); they are not secrets or authentication values, so timing side-channels are irrelevant.

Checklist evidence: CWE-208's condition is comparing security-sensitive values with `==`; the shown source compares package metadata in a test assertion.

#### [ ] Findings `41`, `42`, `43` — `BP-PY-41` (re-appearing audited FPs, old `98`, `99`, `100`)

- Function contexts: `scripts/httptap/findings/functions/41.txt`, `42.txt`, `43.txt`
- Sources: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_property_based.py:117:1, 126:1, 142:1`
- Checklist pattern: tests that raise `AssertionError` explicitly on failure

```
    def test_non_positive_values_rejected(self, value: float) -> None:
        try:
            parse_slo_spec(f"total={value}")
        except SLOSpecError:
            return
        msg = f"expected SLOSpecError for non-positive value {value!r}"
        raise AssertionError(msg)
```

Why this is a false positive: each test verifies rejection by raising `AssertionError` when `SLOSpecError` is not raised; BP-PY-41's assertion detector does not recognize `raise AssertionError`, so it misclassifies verifying tests as assertion-less.

Checklist evidence: BP-PY-41's condition is "tests call production code without assertions"; the shown sources contain explicit `raise AssertionError(msg)` on the failure path, so outcomes are verified.

#### [ ] Findings `44`, `45` — `CWE-295`, `CWE-523` (re-appearing audited FPs, old `102`, `103`)

- Function contexts: `scripts/httptap/findings/functions/44.txt`, `45.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httptap/tests/test_utils.py:341:35`
- Checklist pattern: test assertion verifying the documented `verify_ssl=False` behavior

```
        ctx = create_ssl_context(verify_ssl=False)

        assert ctx.verify_mode == ssl.CERT_NONE
```

Why this is a false positive: the lines assert the deliberate opt-in legacy mode of `create_ssl_context`; the test verifies documented behavior and makes no production request or credential transport.

Checklist evidence: same as findings 38-39 — test assertions, not production bypasses.

### Uncertain findings

None. All 45 fresh findings were classified as false positives with rule conditions verified against the shown source; the one new finding (7) fails the same documented-fallback condition as the adjacent audited constructs.

### Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/httptap/chunks`
- Function evidence: `scripts/httptap/findings/functions`
- Validation: `git diff --check` — `pass`
