# False-positive audit: cylinder

## Run metadata

```yaml
timestamp: 2026-08-02T07:39:32Z
repository: cylinder
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder
branch: main
commit: 7592dfac8c3c8141770f449ab3e65973b84268fc
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder
chunk_path: scripts/cylinder/chunks
function_context_path: scripts/cylinder/findings/functions
```

## Scan evidence

- Build command: `make build` (`go build -o bin/goslop ./cmd/goslop`)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/cylinder/chunks -context-dir scripts/cylinder/findings/functions real-repos/cylinder`
- Findings: `17`
- Chunks reviewed: `scripts/cylinder/chunks/Chunk_1_17.txt`
- Function contexts reviewed: `scripts/cylinder/findings/functions/1.txt` … `scripts/cylinder/findings/functions/17.txt`

## Audit checklist

- [x] Read every assigned chunk under `scripts/cylinder/chunks`.
- [x] Read `scripts/cylinder/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 11 | 2, 4, 7, 8, 11, 12, 13, 14, 15, 16, 17 |
| True positive | 6 | 1, 3, 5, 6, 9, 10 |
| Uncertain | 0 | — |

## False positives

### [ ] Finding `2` — CWE-396

- Function context: `scripts/cylinder/findings/functions/2.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder/src/cylinder.py:268:1`

Source excerpt:

```python
    try:
        response = validate_response(response, run_func_with_dict(params, module.main))
    except Exception as e:
        if isinstance(e, werkzeug.exceptions.HTTPException):
            raise e
        else:
            raise werkzeug.exceptions.InternalServerError(original_exception=e)
```

Why this is a false positive: the handler unconditionally re-raises — `HTTPException` is passed through unchanged and every other exception is wrapped into `InternalServerError` — so no failure condition is hidden and the CWE-396 claim ("can hide distinct failure conditions") does not apply to this source. The sibling rule BP-PY-1 deliberately suppresses this exact construct via its `suiteReraises` check, confirming a re-raising suite is not the intended target.

Checklist evidence: the except suite unconditionally re-raises (`raise e` / `raise werkzeug.exceptions.InternalServerError`); failures are propagated, not hidden.

### [ ] Finding `4` — CWE-93

- Function context: `scripts/cylinder/findings/functions/4.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder/src/cylinder.py:369:17`

Source excerpt:

```python
        response.content_encoding = content_encoding or "identity"
        response.mimetype = mimetype or "application/octet-stream"
        response.headers["Content-Length"] = os.path.getsize(direct_path)
```

Why this is a false positive: the header value is the integer return of `os.path.getsize(...)` — a file size in bytes. An integer cannot contain CR or LF, so CRLF injection into the `Content-Length` header is impossible. The rule itself suppresses numeric-derived values (`headerValueIsInternalNumeric`: `str(int(...))` / `str(round(...))`), confirming numeric values are out of scope; `os.path.getsize` has the same non-injectable shape and is only missed because it is not an `int(...)`/`round(...)` literal.

Checklist evidence: the value written to the header is an integer (file size) with no possible CR/LF bytes; not text that can carry an injected sequence.

### [ ] Finding `7` — BP-PY-6

- Function context: `scripts/cylinder/findings/functions/7.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder/test_sites/foo_site.500.py:5:1`

Source excerpt:

```python
def main(e, response, request):
    # e.original_exception has the original exception
    # e.get_response() can give you the default response for this error
    response.data = str(e.original_exception)
    assert request.shallow == False
    return response
```

Why this is a false positive: `request.shallow` is a framework-internal boolean — the framework itself sets it (`src/cylinder.py:263`: `request.shallow = shallow_request`) from an internal parameter, not from request data. The assert checks an internal invariant in a test-site fixture, not untrusted request input, so the BP-PY-6 rationale (validation stripped by `python -O` weakens a control) does not apply; the `request.` needle fires the rule but the checked value is not request input.

Checklist evidence: the asserted value is an internal framework flag set by `process_module`, not attacker-controlled request data; no security control depends on it.

### [ ] Finding `8` — BP-PY-6

- Function context: `scripts/cylinder/findings/functions/8.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder/test_sites/foo_site.eh.default.py:13:1`

Source excerpt:

```python
    if request.base_url.endswith("/") and request.path != "/":
        new_base_url = request.base_url.rstrip("/")
        if request.query_string:
            abort(308, f"{new_base_url}?{request.query_string.decode('utf-8')}")
        else:
            abort(308, new_base_url)

    assert request.shallow == True
```

Why this is a false positive: `request.shallow` is a framework-internal boolean set by `process_module` (`src/cylinder.py:263`) from an internal parameter; the assert verifies the framework passed `shallow=True` to the early hook in this test-site fixture. It is an internal invariant check, not validation of untrusted input, so `python -O` stripping it cannot weaken any control.

Checklist evidence: the asserted value is an internal framework flag, not request input; the assert has no runtime-validation or security purpose.

### [ ] Finding `11` — BP-PY-6

- Function context: `scripts/cylinder/findings/functions/11.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder/test_sites/foo_site.lh.default.py:3:1`

Source excerpt:

```python
def main(request, response, init, g, log):
    response.headers["late_hook"] = "good"
    assert request.shallow == False
    return response
```

Why this is a false positive: same construct as findings 7 and 8 — `request.shallow` is the framework-internal flag set by `process_module` (`src/cylinder.py:263`), asserted to verify the late hook receives `shallow=False`. This is an internal invariant check in a test-site fixture, not request-input validation; no security control is lost if the assert is stripped by `python -O`.

Checklist evidence: the asserted value is an internal framework flag, not request input; the assert has no runtime-validation or security purpose.

### [ ] Finding `12` — CWE-93

- Function context: `scripts/cylinder/findings/functions/12.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder/tests/cylinder_test.py:15:20`

Source excerpt:

```python
    response = foo_site_client.get("/")
    assert b"hello, Chris" in response.data
    assert response.headers["Early_hook"] == "good"
    assert response.headers["Late_hook"] == "good"
```

Why this is a false positive: this line is an equality comparison (`==`), not a header assignment. The detector splits the line on the first `=` (`strings.Cut(maskedLine, "=")`), so the `==` of the comparison is misread as an assignment `lhs = rhs`; the line only reads the header set by the site hook to the literal `"good"` and asserts equality. No header is written and no CRLF-injectable value exists here.

Checklist evidence: the line is a read/comparison (`assert ... == "good"`), not a write to `response.headers[...]`; no dynamic value reaches a header sink.

### [ ] Finding `13` — BP-PY-7

- Function context: `scripts/cylinder/findings/functions/13.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder/tests/cylinder_test.py:47:31`

Source excerpt:

```python
def test_reserved_method(foo_site_client):
    response = foo_site_client.open("/", method="DEFAULT")
    assert response.status_code == 500
```

Why this is a false positive: `foo_site_client.open(...)` is the `werkzeug.test.Client.open` HTTP-request method, not the builtin file-open function. The detector flags any `.open(` call outside a `with`, but no file handle is opened here, so the resource-leak rationale of the rule does not apply.

Checklist evidence: the trigger is a method call on a test client object performing an HTTP request; no file is opened.

### [ ] Finding `14` — BP-PY-7

- Function context: `scripts/cylinder/findings/functions/14.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder/tests/cylinder_test.py:53:31`

Source excerpt:

```python
def test_bad_method_case(foo_site_client):
    response = foo_site_client.open("/", method="f12")
    assert response.status_code == 500
```

Why this is a false positive: `foo_site_client.open(...)` is the werkzeug test client's HTTP-request method, not the builtin `open`; no file is opened, so the "open without with risks resource leaks" condition is not satisfied.

Checklist evidence: the trigger is a method call on a test client object performing an HTTP request; no file is opened.

### [ ] Finding `15` — BP-PY-7

- Function context: `scripts/cylinder/findings/functions/15.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder/tests/cylinder_test.py:260:31`

Source excerpt:

```python
def test_no_method_bypass(foo_site_client):
    response = foo_site_client.open("/", method="early")
    assert response.status_code == 501
```

Why this is a false positive: `foo_site_client.open(...)` is the werkzeug test client's HTTP-request method, not the builtin `open`; no file is opened, so the "open without with risks resource leaks" condition is not satisfied.

Checklist evidence: the trigger is a method call on a test client object performing an HTTP request; no file is opened.

### [ ] Finding `16` — BP-PY-7

- Function context: `scripts/cylinder/findings/functions/16.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder/tests/cylinder_test.py:266:31`

Source excerpt:

```python
def test_custom_method(foo_site_client):
    response = foo_site_client.open("/custom_method", method="BLARG")
    assert response.status_code == 200
```

Why this is a false positive: `foo_site_client.open(...)` is the werkzeug test client's HTTP-request method, not the builtin `open`; no file is opened, so the "open without with risks resource leaks" condition is not satisfied.

Checklist evidence: the trigger is a method call on a test client object performing an HTTP request; no file is opened.

### [ ] Finding `17` — CWE-117

- Function context: `scripts/cylinder/findings/functions/17.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder/tests/cylinder_test.py:396:9`

Source excerpt:

```python
    log = tiny_queue_app.logger
    for num in range(10):
        log.error(f"test_logger_full {num}")
```

Why this is a false positive: the interpolated value is the local integer loop counter `num` from `range(10)` inside a test. It is not externally controlled and cannot contain CR/LF, so the log-injection concern (attacker-controlled CRLF reaching the log sink) cannot materialize; the f-string form is the only reason the rule's `looksLogFormatted` fires.

Checklist evidence: the formatted value is an internal integer loop counter in test code; no externally controlled input reaches the log message.

## True positives

### CWE-1121 — Excessive McCabe Cyclomatic Complexity

Rule condition (from `detectCWE1121` in `internal/lang/python/detectors/cwe/rules_tier_b_quality.go`): count `if ` / `elif ` / `for ` / `while ` / `except ` tokens in the masked body; flag when ≥ 12.

| Finding id | Source | One-line reason |
| --- | --- | --- |
| 1 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder/src/cylinder.py:75` | `app(environ, start_response)` body genuinely contains ≥ 12 visible branch tokens (13 counted: 8×`if`, 1×`for`, 4×`except`); the complexity finding is accurate (author even added `# ruff: disable[C901, PLR0915]`). |

### CWE-829 — Inclusion of Functionality from Untrusted Control Sphere

Rule condition (from `detectCWE829` in `internal/lang/python/detectors/cwe/rules_code_dynamic.go`): `importlib.util.spec_from_file_location` with a dynamic second argument.

| Finding id | Source | One-line reason |
| --- | --- | --- |
| 3 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder/src/cylinder.py:287` | `spec_from_file_location(path, path)` executes a runtime-resolved `.py` file; `path` comes from `find_processor_path`, which searches paths derived from the request URL (`get_search_paths(site_root, request.path)`) — a genuine dynamically selected code-load sink. |

### BP-PY-1 — Bare Except Clause

Rule condition (from `detectBPPY1` in `internal/lang/python/detectors/bad_practices/rules_core.go`): `except Exception`/`except BaseException` whose suite does not re-raise (matches the rule's own `BP-PY-1-broad-vulnerable` fixture shape).

| Finding id | Source | One-line reason |
| --- | --- | --- |
| 5 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder/src/cylinder.py:390` | `except Exception:` silently swallows the failure and substitutes `record.request_id = ""` — no re-raise, no log; a real defect in the assignment would be hidden. |
| 9 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder/test_sites/foo_site.eh.default.py:24` | `except Exception:` around `json.loads` without re-raise collapses every failure (including handler bugs) into one `abort(400)`; matches the rule's broad-except vulnerable fixture. |

### CWE-396 — Declaration of Catch for Generic Exception

Rule condition (from `detectCWE396` in `internal/lang/python/detectors/cwe/rules_platform.go`): `except Exception` / `except BaseException` in a non-test module.

| Finding id | Source | One-line reason |
| --- | --- | --- |
| 10 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder/test_sites/foo_site.eh.default.py:24` | Generic `except Exception` without re-raise maps every failure to a single 400 response, hiding distinct failure conditions; finding 2 on the same rule is a false positive only because its suite re-raises unconditionally. |

### BP-PY-45 — sys.path Mutation At Runtime

Rule condition (from `detectBPPY45` in `internal/lang/python/detectors/bad_practices/rules_deps.go`): `sys.path.insert/append/extend` in a non-test, non-bootstrap module.

| Finding id | Source | One-line reason |
| --- | --- | --- |
| 6 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder/test_sites/_run_local.py:11` | `sys.path.insert(0, ...)` mutates the import path at runtime to make `src.cylinder` importable — exactly the path hack the rule targets; `_run_local.py` is not recognized as a test file by the rule's name/path checks. |

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/cylinder/chunks`
- Function evidence: `scripts/cylinder/findings/functions`
- Validation: `git diff --check` — pass

## Post-fix remaining-FP audit (2026-08-02)

### Run metadata

```yaml
timestamp: 2026-08-02T16:38:00Z
repository: cylinder
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder
branch: main
commit: 7592dfac8c3c8141770f449ab3e65973b84268fc
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder
chunk_path: scripts/cylinder/chunks
function_context_path: scripts/cylinder/findings/functions
```

### Scan evidence

- Build command: `make build` (`go build -o bin/goslop ./cmd/goslop`)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/cylinder/chunks -context-dir scripts/cylinder/findings/functions real-repos/cylinder`
- Findings: `12`
- Chunks reviewed: `scripts/cylinder/chunks/Chunk_1_12.txt`
- Function contexts reviewed: `scripts/cylinder/findings/functions/1.txt` … `scripts/cylinder/findings/functions/12.txt`

### Audit checklist

- [x] Read every assigned chunk under `scripts/cylinder/chunks`.
- [x] Read `scripts/cylinder/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

### Classification summary (fresh run)

Fresh finding IDs were renumbered (1–12) and matched to the old audit by `Source:` path.

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 5 | 3, 7, 8, 11, 12 |
| True positive | 7 | 1, 2, 4, 5, 6, 9, 10 |
| Uncertain | 0 | — |

All five remaining false positives are re-appearing audited FPs: fresh 3 = old 4, fresh 7 = old 7, fresh 8 = old 8, fresh 11 = old 11, fresh 12 = old 17. Old audited FPs 2, 12, 13, 14, 15, 16 no longer fire in the fresh scan. Fresh finding 5 (CWE-396 at `src/cylinder.py:390`) is a new finding, not present in the old audit; it is a true positive because the suite does not surface the failure (see below).

## False positives

### [ ] Finding `3` — CWE-93

- Function context: `scripts/cylinder/findings/functions/3.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder/src/cylinder.py:369:17`

Source excerpt (from `scripts/cylinder/findings/functions/3.txt`, verified at `src/cylinder.py:369`):

```python
        response.content_encoding = content_encoding or "identity"
        response.mimetype = mimetype or "application/octet-stream"
        response.headers["Content-Length"] = os.path.getsize(direct_path)
```

Why this is a false positive: the header value is the integer return of `os.path.getsize(...)` — a file size in bytes, which cannot contain CR or LF bytes, so CRLF injection is impossible. The rule's own suppression `headerValueIsInternalNumeric` treats numeric-derived values (`str(int(...))` / `str(round(...))`) as out of scope; `os.path.getsize` has the same non-injectable shape and is only missed because it is not an `int(...)`/`round(...)` literal. (Same finding as audited FP 4.)

Checklist evidence: the value written to the header is an integer (file size) with no possible CR/LF bytes; not text that can carry an injected sequence.

### [ ] Finding `7` — BP-PY-6

- Function context: `scripts/cylinder/findings/functions/7.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder/test_sites/foo_site.500.py:5:1`

Source excerpt (from `scripts/cylinder/findings/functions/7.txt`, verified at `test_sites/foo_site.500.py:5`):

```python
def main(e, response, request):
    # e.original_exception has the original exception
    # e.get_response() can give you the default response for this error
    response.data = str(e.original_exception)
    assert request.shallow == False
    return response
```

Why this is a false positive: `request.shallow` is a framework-internal boolean set by the framework itself (`src/cylinder.py:263`), not request data; the assert verifies an internal invariant in a test-site fixture, so stripping it with `python -O` cannot weaken any security control — the BP-PY-6 rationale does not apply. (Same finding as audited FP 7.)

Checklist evidence: the asserted value is an internal framework flag, not request input; the assert has no runtime-validation or security purpose.

### [ ] Finding `8` — BP-PY-6

- Function context: `scripts/cylinder/findings/functions/8.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder/test_sites/foo_site.eh.default.py:13:1`

Source excerpt (from `scripts/cylinder/findings/functions/8.txt`, verified at `test_sites/foo_site.eh.default.py:13`):

```python
    if request.base_url.endswith("/") and request.path != "/":
        new_base_url = request.base_url.rstrip("/")
        if request.query_string:
            abort(308, f"{new_base_url}?{request.query_string.decode('utf-8')}")
        else:
            abort(308, new_base_url)

    assert request.shallow == True
```

Why this is a false positive: same construct as finding 7 — `request.shallow` is the framework-internal flag set by `process_module` (`src/cylinder.py:263`), asserted to verify the early hook receives `shallow=True` in this test-site fixture. It is an internal invariant check, not validation of untrusted input, so `python -O` stripping it cannot weaken any control. (Same finding as audited FP 8.)

Checklist evidence: the asserted value is an internal framework flag, not request input; the assert has no runtime-validation or security purpose.

### [ ] Finding `11` — BP-PY-6

- Function context: `scripts/cylinder/findings/functions/11.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder/test_sites/foo_site.lh.default.py:3:1`

Source excerpt (from `scripts/cylinder/findings/functions/11.txt`, verified at `test_sites/foo_site.lh.default.py:3`):

```python
def main(request, response, init, g, log):
    response.headers["late_hook"] = "good"
    assert request.shallow == False
    return response
```

Why this is a false positive: same construct as findings 7 and 8 — `request.shallow` is the framework-internal flag set by `process_module` (`src/cylinder.py:263`), asserted to verify the late hook receives `shallow=False` in this test-site fixture. An internal invariant check in test-site code; no security control is lost if the assert is stripped by `python -O`. (Same finding as audited FP 11.)

Checklist evidence: the asserted value is an internal framework flag, not request input; the assert has no runtime-validation or security purpose.

### [ ] Finding `12` — CWE-117

- Function context: `scripts/cylinder/findings/functions/12.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/cylinder/tests/cylinder_test.py:396:9`

Source excerpt (from `scripts/cylinder/findings/functions/12.txt`, verified at `tests/cylinder_test.py:395-396`):

```python
    log = tiny_queue_app.logger
    for num in range(10):
        log.error(f"test_logger_full {num}")
```

Why this is a false positive: the interpolated value is the local integer loop counter `num` from `range(10)` inside a test — not externally controlled and cannot contain CR/LF — so the log-injection concern cannot materialize; the f-string form is the only reason the rule's `looksLogFormatted` fires. (Same finding as audited FP 17.)

Checklist evidence: the formatted value is an internal integer loop counter in test code; no externally controlled input reaches the log message.

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/cylinder/chunks/Chunk_1_12.txt`
- Function evidence: `scripts/cylinder/findings/functions/1.txt` … `12.txt`
- Validation: `git diff --check` — pass

## Post-fix v2 audit (latest binary)

Run metadata (build: `make build` ~17:56), scan evidence: `scripts/cylinder/chunks/Chunk_1_12.txt` (12 findings, ids 1–12, same scan command as prior runs), contexts `scripts/cylinder/findings/functions/1.txt`…`12.txt`.

Classification summary (fresh counts): **12 fresh findings — TP 6 / FP 6 / U 0**. Every fresh finding matched a prior classification by `Source:`; no fresh Uncertain, no new findings (all 12 appear in the prior audit: 6 TPs from the original audit, 6 FPs from the Mode A/B appends). Fresh IDs → old IDs: 1→old TP 1, 2→old FP 2, 3→old TP 3, 4→old FP 4, 5→old TP 5, 6→old TP 6, 7→old FP 7, 8→old FP 8, 9→old TP 9, 10→old TP 10, 11→old FP 11, 12→old FP 17.

## Fix checklist (FP patterns)

| Pattern # | Rule | Trigger shape | Count | Example sources |
| --- | --- | --- | ---: | --- |
| 1 | CWE-396 | `except Exception as e:` whose suite **unconditionally re-raises** (`raise e` / `raise ...InternalServerError(original_exception=e)`); the CWE-396 "hides failure conditions" claim does not apply — condition to add: suppress generic-except when every path of the except suite raises (same as BP-PY-1's `suiteReraises`) | 1 | src/cylinder.py:268 |
| 2 | CWE-93 | header value derived from a **numeric function call** (`os.path.getsize(...)` → integer file size, cannot contain CR/LF); condition to add: extend `headerValueIsInternalNumeric` beyond `int(...)`/`round(...)` literals to numeric-returning calls | 1 | src/cylinder.py:369 |
| 3 | BP-PY-6 | `assert <obj>.attr == <literal>` on a **framework-internal attribute** (`request.shallow`) inside a hook/test-site fixture — internal invariant check, not validation of untrusted input; condition to add: skip assert when the asserted target is an attribute of a framework object whose value is set internally (or restrict to request-input-derived targets) | 3 | test_sites/foo_site.500.py:5, test_sites/foo_site.eh.default.py:13, test_sites/foo_site.lh.default.py:3 |
| 4 | CWE-117 | f-string log message interpolating a **local integer loop counter** (`for num in range(10): log.error(f"...{num}")`) in a test file; not externally controlled, no CR/LF possible; condition to add: `looksLogFormatted` should ignore interpolations that resolve to local integer literals/counters (or suppress in test files) | 1 | tests/cylinder_test.py:396 |

## New findings

None — every fresh finding has a prior classification by `Source:` (6 TP / 6 FP as above); no unmatched finding needed fresh classification.
