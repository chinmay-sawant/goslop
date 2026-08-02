# False-positive audit report — FuncToWeb

## Run metadata

```yaml
timestamp: 2026-08-02T07:20:00Z
repository: FuncToWeb
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb
branch: main
commit: 96a90a2b14d667c9fa5957bff4405912e31b0462
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb
chunk_path: scripts/FuncToWeb/chunks
function_context_path: scripts/FuncToWeb/findings/functions
```

## Scan evidence

- Build command: `n/a` (pure Python package, no build step)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/FuncToWeb/chunks -context-dir scripts/FuncToWeb/findings/functions real-repos/FuncToWeb`
- Findings: `57`
- Chunks reviewed: `scripts/FuncToWeb/chunks/Chunk_1_25.txt`, `scripts/FuncToWeb/chunks/Chunk_26_50.txt`, `scripts/FuncToWeb/chunks/Chunk_51_57.txt`
- Function contexts reviewed: `scripts/FuncToWeb/findings/functions/{6,7,18,20,21,23,24,31,32,33,35,36,42,44,45,47,48,49,57}.txt`

## Audit checklist

- [x] Read every assigned chunk under `scripts/FuncToWeb/chunks`.
- [x] Read `scripts/FuncToWeb/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient. (No delegated reviews.)
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 19 | 6, 7, 18, 20, 21, 23, 24, 31, 32, 33, 35, 36, 42, 44, 45, 47, 48, 49, 57 |
| True positive | 38 | 1, 2, 3, 4, 5, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 19, 22, 25, 26, 27, 28, 29, 30, 34, 37, 38, 39, 40, 41, 43, 46, 50, 51, 52, 53, 54, 55, 56 |
| Uncertain | 0 | — |

## False positives

### [ ] Finding 7 — BP-PY-46

- Function context: `scripts/FuncToWeb/findings/functions/7.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/examples/project/gallery.py:101:11`
- Checklist pattern: `print(...)` call token inside a string literal

Source excerpt:

```python
    """Write one thumbnail per photo, reporting each one as it lands.

    The lines appear in the modal while the function runs, because the page it
    embeds is already talking to /invoke-stream. Streaming asks nothing of the
    host: print() is the whole API.
    """
```

Why this is a false positive: Line 101 is the last line of the `bulk_thumbnails` docstring; the text `print() is the whole API.` is documentation, not a `print(...)` call.

Checklist evidence: BP-PY-46 matches `print(...)` calls in non-`__main__` modules; there is no `print(...)` call node on line 101.

### [ ] Finding 6 — CWE-22

- Function context: `scripts/FuncToWeb/findings/functions/6.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/examples/project/gallery.py:62:27`
- Checklist pattern: pathlib join with dynamic segment; segment is typed integer

Source excerpt:

```python
def resize_photo(
    photo: Image.Image,
    width: Annotated[int, Min(16), Max(2000), Label("Width")] = 320,
    height: Annotated[int, Min(16), Max(2000), Label("Height")] = 320,
) -> tuple[Image.Image, Annotated[Path, Download()]]:
    ...
    path = Path(mkdtemp()) / f"resized-{width}x{height}.png"
    resized.save(path)
```

Why this is a false positive: The dynamic segment is built only from `width` and `height`, which the signature types as `Annotated[int, ...]`; an integer cannot carry `/`, `\`, or `..`, so the joined path can never leave the fresh `mkdtemp()` directory.

Checklist evidence: CWE-22 requires a user-influenced pathname that can escape a restricted directory; the only variable parts are validated integers that cannot contain traversal sequences.

### [ ] Finding 18 — CWE-1333

- Function context: `scripts/FuncToWeb/findings/functions/18.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/src/func_to_web/models/function.py:17:16`
- Checklist pattern: nested quantifier `(?:-[…]+)*`; repetition disambiguated by literal separator

Source excerpt:

```python
SLUG_PATTERN = re.compile(r"^[A-Za-z0-9_]+(?:-[A-Za-z0-9_]+)*$")
```

Why this is a false positive: The outer `*` repeats a group that must begin with a literal `-`, which is not a member of the inner class `[A-Za-z0-9_]`; every match position is unambiguous, so the pattern matches in linear time and cannot exhibit catastrophic backtracking.

Checklist evidence: CWE-1333 requires worst-case inefficient/exponential matching; for this anchored pattern no two repetition splits can produce the same string, so complexity is linear.

### [ ] Finding 20 — CWE-89

- Function context: `scripts/FuncToWeb/findings/functions/20.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/src/func_to_web/web/execution.py:33:11`
- Checklist pattern: identifier `execute` is a user-defined async function, not a DB cursor

Source excerpt:

```python
async def execute(
    web_function: WebFunction,
    data: dict[str, Any],
    *,
    capture: PrintCapture | None = None,
    form: FormAction | None = None,
) -> tuple[dict[str, Any], int]:
```

Why this is a false positive: The flag sits on the definition of a custom async wrapper that decodes and calls user functions; the codebase contains no SQL driver, cursor, or query string at all, so no SQL command text reaches any `execute` sink.

Checklist evidence: CWE-89 requires a dynamic SQL string reaching `execute`/`executemany`; the flagged name is a plain function, and no SQL string exists in the module.

### [ ] Finding 31 — CWE-89

- Function context: `scripts/FuncToWeb/findings/functions/31.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/src/func_to_web/web/router.py:105:34`
- Checklist pattern: call to the custom `execute` wrapper, not a DB cursor call

Source excerpt:

```python
    @router.post(f"/{web_function.slug}/invoke")
    async def invoke(data: dict[str, Any]) -> JSONResponse:
        envelope, status = await execute(web_function, data, form=form)
```

Why this is a false positive: `execute` is `execution.execute`, an async wrapper defined in this project; no SQL string is involved, so the CWE-89 sink condition is never reached.

Checklist evidence: CWE-89 requires a dynamic SQL string; the call passes a `WebFunction` and a request dict to a user-defined function.

### [ ] Finding 35 — CWE-89

- Function context: `scripts/FuncToWeb/findings/functions/35.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/src/func_to_web/web/stream.py:56:40`
- Checklist pattern: call to the custom `execute` wrapper, not a DB cursor call

Source excerpt:

```python
    async def run() -> None:
        try:
            holder["response"] = await execute(web_function, data,
                                               capture=capture, form=form)
```

Why this is a false positive: Same custom `execution.execute` wrapper as findings 20 and 31; there is no SQL command text anywhere in the call path.

Checklist evidence: CWE-89 requires a dynamic SQL string reaching a DB sink; the argument is a `WebFunction` object, not SQL text.

### [ ] Finding 32 — BP-PY-32

- Function context: `scripts/FuncToWeb/findings/functions/32.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/src/func_to_web/web/router.py:188:20`
- Checklist pattern: FileResponse path already confined by `segment_of` + `stored_of` resolve/prefix checks

Source excerpt:

```python
            try:
                path, filename = stored_return(reference)
            except (ValueError, FileNotFoundError) as error:
                raise HTTPException(404) from error

            return FileResponse(path, filename=filename,
                                media_type="application/octet-stream")
```

Why this is a false positive: `stored_return` runs `segment_of` (rejects `/`, `\`, `..`, control characters) and `stored_of` (resolves and requires `candidate.parent == RETURNS_DIR.resolve()`), so the FileResponse path is confined to the storage directory before it reaches the sink.

Checklist evidence: BP-PY-32 requires a user-input path without resolve+prefix confinement; the source performs exactly the resolve+prefix confinement the rule prescribes.

### [ ] Finding 33 — BP-PY-32

- Function context: `scripts/FuncToWeb/findings/functions/33.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/src/func_to_web/web/router.py:204:20`
- Checklist pattern: FileResponse path confined by `static_asset` resolve+prefix checks

Source excerpt:

```python
    for root in STATIC_ROOTS:
        try:
            candidate = (root / name).resolve()

            if candidate.is_relative_to(root) and candidate.is_file():
                return candidate
        except (OSError, ValueError):
            continue
    ...
        response = FileResponse(
            path,
            stat_result=path.stat(),
```

Why this is a false positive: `static_asset` resolves the candidate and requires `candidate.is_relative_to(root)` before returning it, so the path handed to `FileResponse` cannot escape the static roots.

Checklist evidence: BP-PY-32 requires no resolve+prefix check; the shown source performs `Path.resolve()` plus `is_relative_to(root)` confinement on the user-supplied `name`.

### [ ] Finding 36 — BP-PY-40

- Function context: `scripts/FuncToWeb/findings/functions/36.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/src/func_to_web/web/upload.py:210:17`
- Checklist pattern: thread is explicitly `daemon=True`; rule targets non-daemon workers

Source excerpt:

```python
    if needed and _sweeper is None:
        _sweeper = Thread(target=_sweeping, name=SWEEPER_NAME, daemon=True)
        _sweeper.start()
```

Why this is a false positive: The rule's own detection note scopes it as a "review-only heuristic for non-daemon workers"; this sweeper thread is constructed with `daemon=True`, the documented policy that lets the process exit without a join, and its lifecycle is a named infinite loop (`_sweeping`).

Checklist evidence: BP-PY-40 matches `Thread(...).start()` without a join/daemon policy; the shown source sets the daemon policy explicitly at construction.

### [ ] Finding 42 — CWE-1341

- Function context: `scripts/FuncToWeb/findings/functions/42.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/tests/e2e/conftest.py:34:5`
- Checklist pattern: single `close()` in fixture teardown; no second release in the shown source

Source excerpt:

```python
    try:
        instance = playwright_engine.chromium.launch()
    except Error as error:
        pytest.skip(f"no playwright browser available: {error}")

    yield instance

    instance.close()
```

Why this is a false positive: The fixture opens the browser once and releases it exactly once in teardown; nothing in the shown source releases `instance` a second time.

Checklist evidence: CWE-1341 requires the same handle to be released twice; the shown source contains a single `close()` after the `yield`.

### [ ] Finding 44 — CWE-93

- Function context: `scripts/FuncToWeb/findings/functions/44.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/tests/integration/test_http_invoke.py:238:19`
- Checklist pattern: header read inside an `assert`, not a header write

Source excerpt:

```python
    assert [answer.status_code for answer in answers] == [200, 422, 500]
    assert {answer.headers["content-type"] for answer in answers} == {
        "application/json"
    }
```

Why this is a false positive: The flagged expression reads `response.headers["content-type"]` to assert on it; nothing is written to a response header, so no CRLF injection sink exists.

Checklist evidence: CWE-93 requires a dynamic value written to an HTTP response header; the shown source is a read-only assertion on the test client's response.

### [ ] Finding 45 — CWE-93

- Function context: `scripts/FuncToWeb/findings/functions/45.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/tests/integration/test_http_stream.py:307:20`
- Checklist pattern: header read inside an `assert`, not a header write

Source excerpt:

```python
    assert response.status_code == 200
    assert response.headers["content-type"].startswith("text/event-stream")
    assert response.headers["cache-control"] == "no-cache"
```

Why this is a false positive: `response.headers[...]` is a read of the received response used in assertions; there is no code path that sets a header value from dynamic input.

Checklist evidence: CWE-93 requires a header-write sink; the shown source only reads header values.

### [ ] Finding 47 — BP-PY-1

- Function context: `scripts/FuncToWeb/findings/functions/47.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/tests/integration/test_http_stream.py:479:1`
- Checklist pattern: handler records the exception and re-raises it after the thread joins

Source excerpt:

```python
    def worker():
        try:
            caught["value"] = asyncio.run(work())
        except BaseException as error:
            caught["error"] = error
    ...
    if "error" in caught:
        raise caught["error"]
```

Why this is a false positive: The `except BaseException` handler stores the exception in `caught["error"]` and the caller re-raises it afterwards; the failure is recorded and surfaced, so it is not swallowed.

Checklist evidence: BP-PY-1 matches handlers that pass, only log without re-raise, or continue without recording the exception type; here the exception type is recorded and re-raised.

### [ ] Finding 48 — CWE-93

- Function context: `scripts/FuncToWeb/findings/functions/48.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/tests/integration/test_prefix.py:120:20`
- Checklist pattern: header read inside an `assert`, not a header write

Source excerpt:

```python
    assert response.status_code == 200
    assert response.headers["content-type"] == "image/svg+xml"
```

Why this is a false positive: Reading `response.headers["content-type"]` in an assertion is a header read, not a write; no CRLF-injectable sink is present.

Checklist evidence: CWE-93 requires a dynamic value written to an HTTP response header; the shown source reads the header only.

### [ ] Finding 49 — BP-PY-42

- Function context: `scripts/FuncToWeb/findings/functions/49.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/tests/shared.py:23:1`
- Checklist pattern: helper `warm_uvicorn()`, not a `test_*` function expecting failure

Source excerpt:

```python
def warm_uvicorn() -> None:
    # uvicorn.Config.load() imports its protocol backends lazily ...
    # Importing the backends once here, with warnings ignored, leaves them in
    # sys.modules and lets live_server work as it was meant to.
    with warnings.catch_warnings():
        warnings.simplefilter("ignore")

        try:
            uvicorn.Config(FastAPI(), host="127.0.0.1", port=0).load()
        except Exception:
            pass
```

Why this is a false positive: `warm_uvicorn` is an import-time test helper, not a test function, and the try/except suppresses a known import-time `DeprecationWarning` turned into an error by `filterwarnings=["error"]`; it does not use try/except to "expect failure" in a test.

Checklist evidence: BP-PY-42 matches test functions using bare try/except instead of `assertRaises`/`pytest.raises`; the flagged code is a non-test helper whose try/except is a documented import warm-up.

### [ ] Finding 57 — CWE-93

- Function context: `scripts/FuncToWeb/findings/functions/57.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/tests/unit/test_outputs_downloads.py:272:20`
- Checklist pattern: header read inside an `assert`, not a header write

Source excerpt:

```python
    response = client.get(f"/returns/{result['value']}")

    assert response.headers["content-type"] == "application/octet-stream"
```

Why this is a false positive: The flagged expression reads a response header inside an assertion; there is no header-write sink, so no CRLF injection is possible.

Checklist evidence: CWE-93 requires a dynamic value written to an HTTP response header; the shown source only reads the header.

### [ ] Finding 21 — BP-PY-1

- Function context: `scripts/FuncToWeb/findings/functions/21.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/src/func_to_web/web/execution.py:45:1`
- Checklist pattern: handler records the exception type in the failure envelope returned to the caller

Source excerpt:

```python
    try:
        values = decode(web_function.schema, data,
                        file_resolver=stored_file)

        kwargs = web_function.schema.build(values)
    except Exception as error:
        return _failure(error, UNPROCESSABLE)
```

Why this is a false positive: The handler passes the exception to `_failure`, which records `f"{type(error).__name__}: {error}"` in the JSON failure envelope returned to the caller; the failure is neither passed over nor left unrecorded, so it does not match the rule's swallow patterns.

Checklist evidence: BP-PY-1 matches handlers that pass, only log without re-raise, or continue without recording the exception type; here the exception type and message are recorded into the response.

### [ ] Finding 23 — BP-PY-1

- Function context: `scripts/FuncToWeb/findings/functions/23.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/src/func_to_web/web/execution.py:50:1`
- Checklist pattern: handler records the exception type in the failure envelope returned to the caller

Source excerpt:

```python
    try:
        value = await _call(web_function.fn, kwargs, capture)
    except Exception as error:
        return _failure(error, FAILED)
```

Why this is a false positive: Same error-boundary pattern as finding 21 — `_failure` embeds the exception type and message in the returned envelope, so the error is recorded and surfaced.

Checklist evidence: BP-PY-1's swallow patterns (pass, log-without-reraise, continue-without-recording) do not apply; the exception is recorded in the failure response.

### [ ] Finding 24 — BP-PY-1

- Function context: `scripts/FuncToWeb/findings/functions/24.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/src/func_to_web/web/execution.py:61:1`
- Checklist pattern: handler records the exception type in the failure envelope returned to the caller

Source excerpt:

```python
        return {"result": output_of(resolved, download=_download_of)}, OK
    except Exception as error:
        return _failure(error, FAILED)
```

Why this is a false positive: Same deliberate boundary as findings 21 and 23 — the exception is recorded by `_failure` into the structured error response rather than hidden.

Checklist evidence: BP-PY-1's swallow patterns do not apply; the exception type is recorded in the failure envelope.

## Uncertain findings

None.

## True positives

### BP-PY-46 — print Debugging In Library Code (real `print(...)` calls, non-`__main__`, non-test)

| Finding | Source | Reason |
| --- | --- | --- |
| 1 | examples/fastapi/sdk_frontend.py:79 | Real `print(f"file {index} of {files}")` in a module function, not under `__main__` |
| 2 | examples/http/server.py:38 | Real `print(f"{index}...")` in `countdown`, not under `__main__` |
| 4 | examples/outputs/no_result.py:8 | Real `print(f"Archiving order {reference}")` in `archive_order` |
| 5 | examples/outputs/no_result.py:13 | Real `print(f"Clearing queue {name}")` in `clear_queue` |
| 8 | examples/project/gallery.py:110 | Real `print(f"[{index}/{len(photos)}] ...")` in `bulk_thumbnails` |
| 9 | examples/streaming/async_progress.py:13 | Real `print(f"fetching {pages} page(s)")` in `fetch_pages` |
| 10 | examples/streaming/async_progress.py:17 | Real `print(f"page {index} of {pages} ready")` in loop |
| 11 | examples/streaming/basic_prints.py:10 | Real `print(f"hello {name}")` in `announce` |
| 12 | examples/streaming/basic_prints.py:18 | Real `print(f"step {index} of {items} done")` in loop |
| 13 | examples/streaming/capture_policy.py:11 | Real `print(f"audited step {index}")` in loop |
| 14 | examples/streaming/capture_policy.py:19 | Real `print(f"silent step {index}")` in loop |
| 15 | examples/streaming/print_then_error.py:13 | Real `print(f"row {index} read")` in loop |
| 16 | examples/streaming/progress.py:13 | Real `print(f"converting {files} file(s)")` in `convert` |
| 17 | examples/streaming/progress.py:18 | Real `print(f"[{percent:3d}%] file {index} of {files}")` in loop |
| 34 | src/func_to_web/web/run.py:84 | Real `print(f"{NAME} {__version__}...")` banner inside library `run()`, not under `__main__` |

### CWE-22 — Path Traversal

| Finding | Source | Reason |
| --- | --- | --- |
| 3 | examples/outputs/download_optional.py:28 | Free-form `reference: str` from the web request joined without confinement into `path.open("w")` sink |

### CWE-396 — Declaration of Catch for Generic Exception

| Finding | Source | Reason |
| --- | --- | --- |
| 19 | src/func_to_web/models/return_parser.py:486 | `except Exception as error` re-wrapped into `ReturnContractError` |
| 22 | src/func_to_web/web/execution.py:45 | Generic `except Exception` catch present |
| 29 | src/func_to_web/web/pending.py:174 | Generic `except Exception` catch present |
| 38 | src/func_to_web/web/upload.py:233 | Generic `except Exception` catch present in sweeper thread |

### BP-PY-1 — Bare Except Clause (swallow variants)

| Finding | Source | Reason |
| --- | --- | --- |
| 28 | src/func_to_web/web/pending.py:174 | `except Exception: continue` — continues without recording the exception type |
| 30 | src/func_to_web/web/pending.py:199 | `except Exception: continue` — continues without recording the exception type |
| 37 | src/func_to_web/web/upload.py:233 | `except Exception: continue` in sweeper loop, error unrecorded |
| 50 | tests/shared.py:25 | `except Exception: pass` — handler is bare pass |

### BP-PY-2 — Except Pass

| Finding | Source | Reason |
| --- | --- | --- |
| 25 | src/func_to_web/web/pending.py:150 | `except OSError: pass` — handler suite is solely `pass` |
| 39 | src/func_to_web/web/upload.py:321 | `except OSError: pass` — handler suite is solely `pass` |
| 51 | tests/shared.py:25 | `except Exception: pass` — handler suite is solely `pass` |

### CWE-390 — Detection of Error Condition Without Action

| Finding | Source | Reason |
| --- | --- | --- |
| 26 | src/func_to_web/web/pending.py:150 | Exception detected, handler takes no action |
| 40 | src/func_to_web/web/upload.py:321 | Exception detected, handler takes no action |
| 52 | tests/shared.py:25 | Exception detected, handler takes no action |

### CWE-1071 — Empty Code Block

| Finding | Source | Reason |
| --- | --- | --- |
| 27 | src/func_to_web/web/pending.py:150 | Exception handler contains only `pass` |
| 41 | src/func_to_web/web/upload.py:321 | Exception handler contains only `pass` |
| 53 | tests/shared.py:25 | Exception handler contains only `pass` |

### BP-PY-41 — pytest assert With Side Effects Only

| Finding | Source | Reason |
| --- | --- | --- |
| 43 | tests/integration/test_docs_snippets.py:390 | `test_every_complete_python_block_parses` calls `ast.parse(block.code)` with no assert |
| 54 | tests/unit/test_examples_import.py:49 | `test_every_module_of_the_collections_compiles` calls `compile(...)` with no assert |

### CWE-1046 — Creation of Immutable Text Using String Concatenation

| Finding | Source | Reason |
| --- | --- | --- |
| 46 | tests/integration/test_http_stream.py:426 | `body += chunk` repeatedly concatenates inside the `for chunk in response.iter_text()` loop |

### BP-PY-12 — eval Or exec On Dynamic Input

| Finding | Source | Reason |
| --- | --- | --- |
| 55 | tests/unit/test_examples_import.py:50 | `compile(source.read_text(...), str(source), "exec")` — non-literal input in exec mode |

### CWE-94 — Code Injection

| Finding | Source | Reason |
| --- | --- | --- |
| 56 | tests/unit/test_examples_import.py:50 | Dynamic value (file content) reaches the `compile` code-generation sink |

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/FuncToWeb/chunks`
- Function evidence: `scripts/FuncToWeb/findings/functions`
- Validation: `git diff --check` — pass

## Post-fix over-suppression audit (2026-08-02)

Run metadata:

```yaml
timestamp: 2026-08-02T16:45:00Z
repository: FuncToWeb
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb
branch: main
commit: 96a90a2b14d667c9fa5957bff4405912e31b0462
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb
scanner_binary: ./bin/goslop rebuilt at b5b8fde (FP-reduction fix, 2026-08-02 16:29)
chunk_path: scripts/FuncToWeb/chunks
function_context_path: scripts/FuncToWeb/findings/functions
```

Fresh scan produced 32 findings; the audited report lists 38 true positives. 23 of those TPs re-appear in the fresh scan (matched by `Source:` file:line, e.g. old 3→fresh 1, 22→5, 25–30→8–13, 34→16, 37–41→18–22, 43→23, 46→24, 50–53→26–29, 54–56→30–32). The remaining 15 audited TPs are absent from the fresh scan. Each was checked against the current source at the audited location: all 15 constructs are still present, so all 15 are over-suppressed, not fixed.

| Old finding ID | Rule | Source | One-line reason (from old audit) | Current status |
| --- | --- | --- | --- | --- |
| 1 | BP-PY-46 | examples/fastapi/sdk_frontend.py:79 | Real `print(f"file {index} of {files}")` in a module function, not under `__main__` | suppressed-but-present |
| 2 | BP-PY-46 | examples/http/server.py:38 | Real `print(f"{index}...")` in `countdown`, not under `__main__` | suppressed-but-present |
| 4 | BP-PY-46 | examples/outputs/no_result.py:8 | Real `print(f"Archiving order {reference}")` in `archive_order` | suppressed-but-present |
| 5 | BP-PY-46 | examples/outputs/no_result.py:13 | Real `print(f"Clearing queue {name}")` in `clear_queue` | suppressed-but-present |
| 8 | BP-PY-46 | examples/project/gallery.py:110 | Real `print(f"[{index}/{len(photos)}] ...")` in `bulk_thumbnails` | suppressed-but-present |
| 9 | BP-PY-46 | examples/streaming/async_progress.py:13 | Real `print(f"fetching {pages} page(s)")` in `fetch_pages` | suppressed-but-present |
| 10 | BP-PY-46 | examples/streaming/async_progress.py:17 | Real `print(f"page {index} of {pages} ready")` in loop | suppressed-but-present |
| 11 | BP-PY-46 | examples/streaming/basic_prints.py:10 | Real `print(f"hello {name}")` in `announce` | suppressed-but-present |
| 12 | BP-PY-46 | examples/streaming/basic_prints.py:18 | Real `print(f"step {index} of {items} done")` in loop | suppressed-but-present |
| 13 | BP-PY-46 | examples/streaming/capture_policy.py:11 | Real `print(f"audited step {index}")` in loop | suppressed-but-present |
| 14 | BP-PY-46 | examples/streaming/capture_policy.py:19 | Real `print(f"silent step {index}")` in loop | suppressed-but-present |
| 15 | BP-PY-46 | examples/streaming/print_then_error.py:13 | Real `print(f"row {index} read")` in loop | suppressed-but-present |
| 16 | BP-PY-46 | examples/streaming/progress.py:13 | Real `print(f"converting {files} file(s)")` in `convert` | suppressed-but-present |
| 17 | BP-PY-46 | examples/streaming/progress.py:18 | Real `print(f"[{percent:3d}%] file {index} of {files}")` in loop | suppressed-but-present |
| 19 | CWE-396 | src/func_to_web/models/return_parser.py:486 | `except Exception as error` re-wrapped into `ReturnContractError` | suppressed-but-present |

Fixed-removed: 0. Root cause: the fix's `isPythonScriptModule` exemption (added in `internal/lang/python/detectors/bad_practices/common.go:72-98`, path component `examples/` → BP-PY-46 skipped entirely) suppresses all 14 example prints, and the broad-except "re-raise" exemption suppresses the CWE-396 at return_parser.py:486. None of the suppressed prints sit inside an `if __name__ == "__main__":` block (the guards appear at lines 67–213, after the flagged functions), so the main-guard exemption would not have applied; the path exemption is the sole cause.

### [ ] Finding 1 — BP-PY-46

- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/examples/fastapi/sdk_frontend.py:79`

Source excerpt:

```python
def convert(files: Annotated[int, Min(1), Max(8)] = 5) -> str:
    """Report the progress of a slow job while it advances."""
    for index in range(1, files + 1):
        time.sleep(STEP_SECONDS)
        print(f"file {index} of {files}")
```

Why it satisfies the rule condition: module-level `print(...)` call in an importable module (`examples/fastapi/sdk_frontend.py` is imported by the host server; the `__main__` guard is at line 97, after `convert`), matching BP-PY-46's non-script print-in-library pattern.

### [ ] Finding 2 — BP-PY-46

- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/examples/http/server.py:38`

Source excerpt:

```python
def countdown(steps: Annotated[int, Min(1), Max(10)] = 3) -> str:
    """Print one line per step so the stream carries print events."""
    for index in reversed(range(1, steps + 1)):
        print(f"{index}...")
```

Why it satisfies the rule condition: executable `print(f"{index}...")` inside `countdown`, a module-level function not under the `__main__` guard at line 67.

### [ ] Finding 4 — BP-PY-46

- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/examples/outputs/no_result.py:8`

Source excerpt:

```python
def archive_order(reference: str) -> None:
    """Do the work and return nothing, which the web shows as "Done"."""
    print(f"Archiving order {reference}")
```

Why it satisfies the rule condition: `print(...)` at module level in an importable function; the `__main__` guard at line 18 only wraps `run([...])`.

### [ ] Finding 5 — BP-PY-46

- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/examples/outputs/no_result.py:13`

Source excerpt:

```python
def clear_queue(name: str) -> list[str]:
    """Return an empty collection, which also shows "Done"."""
    print(f"Clearing queue {name}")
```

Why it satisfies the rule condition: same module as finding 4, second module-level `print(...)` in an importable function, not under a `__main__` guard.

### [ ] Finding 8 — BP-PY-46

- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/examples/project/gallery.py:110`

Source excerpt:

```python
    for index, path in enumerate(photos, start=1):
        with Image.open(path) as picture:
            picture.thumbnail((size, size))
            picture.convert("RGB").save(folder / f"thumb-{index}.png")

        print(f"[{index}/{len(photos)}] {Path(path).name} at {size}px")
```

Why it satisfies the rule condition: `print(...)` inside `bulk_thumbnails`, a module-level function (guard at line 213); note line 101's `print() is the whole API.` is a docstring, but line 110 is a real call node.

### [ ] Finding 9 — BP-PY-46

- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/examples/streaming/async_progress.py:13`

Source excerpt:

```python
async def fetch_pages(pages: Annotated[int, Min(1), Max(8)] = 4) -> str:
    """Report the progress of an awaited job."""
    print(f"fetching {pages} page(s)")
```

Why it satisfies the rule condition: `print(f"fetching {pages} page(s)")` is executable code in an async module-level function, before the `__main__` guard at line 22.

### [ ] Finding 10 — BP-PY-46

- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/examples/streaming/async_progress.py:17`

Source excerpt:

```python
    for index in range(1, pages + 1):
        await asyncio.sleep(STEP_SECONDS)
        print(f"page {index} of {pages} ready")
```

Why it satisfies the rule condition: second module-level `print(...)` in `fetch_pages`, not under the `__main__` guard.

### [ ] Finding 11 — BP-PY-46

- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/examples/streaming/basic_prints.py:10`

Source excerpt:

```python
def announce(name: str = "world") -> str:
    """Print a single line before returning."""
    print(f"hello {name}")
```

Why it satisfies the rule condition: module-level `print(f"hello {name}")` in `announce`; the guard at line 23 only runs `run(...)`.

### [ ] Finding 12 — BP-PY-46

- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/examples/streaming/basic_prints.py:18`

Source excerpt:

```python
def checklist(items: Annotated[int, Min(1), Max(20)] = 4) -> str:
    """Print one line per item and report the total."""
    for index in range(1, items + 1):
        print(f"step {index} of {items} done")
```

Why it satisfies the rule condition: second module-level `print(...)` in `basic_prints.py`, not under the `__main__` guard.

### [ ] Finding 13 — BP-PY-46

- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/examples/streaming/capture_policy.py:11`

Source excerpt:

```python
def audited(steps: Annotated[int, Min(1), Max(10)] = 3) -> str:
    """Inherit the policy of the space, so its prints are streamed."""
    for index in range(steps):
        print(f"audited step {index}")
```

Why it satisfies the rule condition: module-level `print(f"audited step {index}")` in `audited`, before the guard at line 24.

### [ ] Finding 14 — BP-PY-46

- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/examples/streaming/capture_policy.py:19`

Source excerpt:

```python
def silent(steps: Annotated[int, Min(1), Max(10)] = 3) -> str:
    """Print only to the server console: the browser sees no print event."""
    for index in range(steps):
        print(f"silent step {index}")
```

Why it satisfies the rule condition: second module-level `print(...)` in `capture_policy.py`, not under the `__main__` guard.

### [ ] Finding 15 — BP-PY-46

- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/examples/streaming/print_then_error.py:13`

Source excerpt:

```python
def import_rows(rows: Annotated[int, Min(1), Max(10)] = 5) -> str:
    """Read rows one by one and fail on the malformed one."""
    for index in range(1, rows + 1):
        print(f"row {index} read")
```

Why it satisfies the rule condition: module-level `print(f"row {index} read")` in `import_rows`; the guard at line 21 comes after the function.

### [ ] Finding 16 — BP-PY-46

- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/examples/streaming/progress.py:13`

Source excerpt:

```python
def convert(files: Annotated[int, Min(1), Max(8)] = 5) -> str:
    """Report the progress of a slow job while it advances."""
    print(f"converting {files} file(s)")
```

Why it satisfies the rule condition: module-level `print(f"converting {files} file(s)")` in `convert`, before the guard at line 23.

### [ ] Finding 17 — BP-PY-46

- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/examples/streaming/progress.py:18`

Source excerpt:

```python
    for index in range(1, files + 1):
        time.sleep(STEP_SECONDS)
        percent = index * 100 // files
        print(f"[{percent:3d}%] file {index} of {files}")
```

Why it satisfies the rule condition: second module-level `print(...)` in `progress.py`'s `convert`, not under the `__main__` guard.

### [ ] Finding 19 — CWE-396

- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/FuncToWeb/src/func_to_web/models/return_parser.py:486`

Source excerpt:

```python
    try:
        produced = filename(value, index)
    except Exception as error:
        raise ReturnContractError(
            f"Download filename callable failed: "
            f"{type(error).__name__}: {error}"
        ) from error
```

Why it satisfies the rule condition: `except Exception as error` is a generic catch declared in library code, which CWE-396 flags regardless of how the handler re-uses `error`; the fix's re-raise exemption suppressed it, but the audit classified the generic catch itself as the true positive.

## Final evidence (post-fix over-suppression audit)

- Delegated reviewers: none
- Chunk evidence: `scripts/FuncToWeb/chunks/Chunk_1_25.txt`, `scripts/FuncToWeb/chunks/Chunk_26_32.txt`
- Source evidence: current files at `real-repos/FuncToWeb` commit `96a90a2b14d667c9fa5957bff4405912e31b0462` (unchanged since the original audit)
- Validation: `git diff --check` — pass
