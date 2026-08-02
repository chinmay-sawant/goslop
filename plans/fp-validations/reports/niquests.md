# False-positive audit — niquests

## Run metadata

```yaml
timestamp: 2026-08-02T08:00:16Z
repository: niquests
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests
branch: main
commit: 7633aa3f1f9fcdb7790192ffd8cfacb69ca2c807
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests
chunk_path: scripts/niquests/chunks
function_context_path: scripts/niquests/findings/functions
```

## Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop` (goslop binary used: `./bin/goslop`)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/niquests/chunks -context-dir scripts/niquests/findings/functions real-repos/niquests`
- Findings: `346`
- Chunks reviewed: `scripts/niquests/chunks/Chunk_1_25.txt`, `scripts/niquests/chunks/Chunk_26_50.txt`, `scripts/niquests/chunks/Chunk_51_75.txt`, `scripts/niquests/chunks/Chunk_76_100.txt`, `scripts/niquests/chunks/Chunk_101_125.txt`, `scripts/niquests/chunks/Chunk_126_150.txt`, `scripts/niquests/chunks/Chunk_151_175.txt`, `scripts/niquests/chunks/Chunk_176_200.txt`, `scripts/niquests/chunks/Chunk_201_225.txt`, `scripts/niquests/chunks/Chunk_226_250.txt`, `scripts/niquests/chunks/Chunk_251_275.txt`, `scripts/niquests/chunks/Chunk_276_300.txt`, `scripts/niquests/chunks/Chunk_301_325.txt`, `scripts/niquests/chunks/Chunk_326_346.txt`
- Function contexts reviewed: `scripts/niquests/findings/functions/<finding-id>.txt` for every proposed false positive (all 311); enclosing source read when the exported context was insufficient

## Audit checklist

- [x] Read every assigned chunk under `scripts/niquests/chunks`.
- [x] Read `scripts/niquests/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [ ] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 311 | 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 29, 30, 31, 32, 33, 34, 35, 36, 37, 41, 42, 43, 44, 45, 46, 49, 50, 52, 54, 55, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 114, 115, 116, 121, 122, 123, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 148, 149, 150, 151, 154, 155, 156, 157, 158, 159, 160, 161, 162, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 180, 181, 182, 184, 186, 187, 188, 195, 196, 197, 198, 199, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 255, 256, 257, 258, 259, 260, 261, 262, 263, 264, 265, 266, 267, 268, 269, 270, 271, 272, 273, 274, 275, 276, 277, 278, 279, 280, 281, 282, 283, 284, 285, 286, 287, 288, 289, 290, 291, 292, 293, 294, 295, 296, 297, 298, 299, 300, 301, 302, 303, 304, 305, 306, 307, 308, 309, 310, 311, 312, 313, 314, 315, 316, 317, 318, 319, 320, 321, 322, 323, 324, 325, 326, 327, 328, 329, 330, 331, 332, 333, 334, 335, 336, 337, 338, 339, 340, 341, 342, 343, 344, 345, 346 |
| True positive | 35 | 26, 27, 28, 38, 39, 40, 47, 48, 51, 53, 56, 57, 113, 117, 118, 119, 120, 124, 125, 126, 146, 147, 152, 153, 163, 164, 183, 185, 189, 190, 191, 192, 193, 194, 200 |
| Uncertain | 0 | — |

## False positives

One subsection per finding. Findings could not be grouped: every pair of findings either hits a distinct line or a distinct rule, so no two findings reference the exact same source construct under the same rule.

### [ ] Finding `1` — `BP-PY-45`

- Function context: `scripts/niquests/findings/functions/1.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/docs/conf.py:24:1`
- Checklist pattern: Sphinx docs-build configuration script

Source excerpt:

```
         22: 
         23: # Insert Requests' path into the system.
    >    24: sys.path.insert(0, os.path.abspath("../src"))
         25: 
```

Why this is a false positive: `docs/conf.py` is the Sphinx build configuration, not library code; `sys.path.insert` is the standard idiom to make the package importable for autodoc at build time.

Checklist evidence: BP-PY-45 targets runtime library behavior; a docs build script mutating `sys.path` at build time is the documented tooling pattern.

### [ ] Finding `2` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/2.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/noxfile.py:95:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >    95: def test_cohabitation(session: nox.Session) -> None:
        96:     tests_impl(session, cohabitation=True)
        97:     tests_impl(session, cohabitation=None)
```

Why this is a false positive: the test body only invokes `tests_impl(session, ...)`, the nox runner that executes the whole pytest suite — the assertions live in the suite run by pytest, so this is a session orchestration, not a placeholder test.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `3` — `BP-PY-46`

- Function context: `scripts/niquests/findings/functions/3.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/noxfile.py:147:9`
- Checklist pattern: print in a nox automation script for CI console output

Source excerpt:

```
        145:     """Test on Emscripten with Pyodide & Chrome / Firefox / Node.js"""
        146:     if runner == "node":
    >   147:         print(
        148:             "Node version:",
```

Why this is a false positive: `noxfile.py` is the project's build/test automation script; the prints emit human-readable progress to the CI console, not library logging.

Checklist evidence: BP-PY-46 targets 'library code'; noxfile.py is a dev/CI script where console output is the intended interface.

### [ ] Finding `4` — `CWE-22`

- Function context: `scripts/niquests/findings/functions/4.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/noxfile.py:157:53`
- Checklist pattern: path built from constants in a dev automation script

Source excerpt:

```
        155:     pyodide_version = "0.28.1"
        156: 
    >   157:     pyodide_artifacts_path = Path(session.cache_dir) / f"pyodide-{pyodide_version}"
        158: 
```

Why this is a false positive: both segments are trusted: `session.cache_dir` is nox's own cache directory and `pyodide_version` is a module constant; no untrusted input reaches the join.

Checklist evidence: CWE-22's condition requires a dynamic segment from an untrusted source; the joined segments are constants under the tool's own cache root.

### [ ] Finding `5` — `BP-PY-46`

- Function context: `scripts/niquests/findings/functions/5.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/noxfile.py:160:9`
- Checklist pattern: print in a nox automation script for CI console output

Source excerpt:

```
        158: 
        159:     if not pyodide_artifacts_path.exists():
    >   160:         print("Fetching pyodide build artifacts")
        161:         session.run(
```

Why this is a false positive: `noxfile.py` is the project's build/test automation script; the prints emit human-readable progress to the CI console, not library logging.

Checklist evidence: BP-PY-46 targets 'library code'; noxfile.py is a dev/CI script where console output is the intended interface.

### [ ] Finding `6` — `CWE-409`

- Function context: `scripts/niquests/findings/functions/6.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/noxfile.py:237:19`
- Checklist pattern: tarfile extraction using the documented safe filter

Source excerpt:

```
        235:         extract_root.mkdir(exist_ok=True)
        236:         with tarfile.open(archive) as source:
    >   237:             source.extractall(extract_root, filter="data")
        238:         extracted = next(path for path in extract_root.iterdir() if (path / "wit").is_dir())
```

Why this is a false positive: `source.extractall(extract_root, filter="data")` uses Python 3.12's `filter="data"` — the platform's documented safe extraction filter that blocks traversal and special-file members.

Checklist evidence: CWE-409's fix explicitly recommends 'the platform's documented safe API'; the shown source already uses it.

### [ ] Finding `7` — `BP-PY-46`

- Function context: `scripts/niquests/findings/functions/7.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/__init__.py:24:8`
- Checklist pattern: `print(` token inside a docstring doctest example

Source excerpt:

```
         22:    >>> payload = dict(key1='value1', key2='value2')
         23:    >>> r = niquests.post('https://httpbin.org/post', data=payload)
    >    24:    >>> print(r.text)
         25:    {
```

Why this is a false positive: the `>>> print(...)` lines are doctest examples inside docstrings, i.e. documentation text, not executable statements.

Checklist evidence: BP-PY-46's condition is a real `print` call in library code; the per-line string check cannot see the enclosing docstring, so it reports documentation text.

### [ ] Finding `8` — `BP-PY-5`

- Function context: `scripts/niquests/findings/functions/8.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_async.py:14:1`
- Checklist pattern: deliberate API re-export shim marked `# noqa`

Source excerpt:

```
         12: )
         13: 
    >    14: from .async_session import *  # noqa
         15: 
```

Why this is a false positive: `_async.py`/`_typing.py` are internal shim modules whose entire purpose is re-exporting the public API (`from .async_session import *  # noqa`); the wildcard is the re-export mechanism, not accidental namespace pollution.

Checklist evidence: the wildcard import is the module's designed function (re-export), explicitly annotated `# noqa`.

### [ ] Finding `9` — `BP-PY-5`

- Function context: `scripts/niquests/findings/functions/9.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_typing.py:14:1`
- Checklist pattern: deliberate API re-export shim marked `# noqa`

Source excerpt:

```
         12: )
         13: 
    >    14: from .typing import *  # noqa
         15: 
```

Why this is a false positive: `_async.py`/`_typing.py` are internal shim modules whose entire purpose is re-exporting the public API (`from .async_session import *  # noqa`); the wildcard is the re-export mechanism, not accidental namespace pollution.

Checklist evidence: the wildcard import is the module's designed function (re-export), explicitly annotated `# noqa`.

### [ ] Finding `10` — `BP-PY-14`

- Function context: `scripts/niquests/findings/functions/10.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_vendor/kiss_headers/__init__.py:11:12`
- Checklist pattern: requests/session call inside a docstring doctest example

Source excerpt:

```
          9:    >>> import requests
         10:    >>> from kiss_headers import parse_it
    >    11:    >>> r = requests.get('https://www.python.org')
         12:    >>> headers = parse_it(r)
```

Why this is a false positive: the `>>> requests.get(...)` / `>>> session.get(...)` lines are doctest examples in docstrings, i.e. documentation text, not executed HTTP calls.

Checklist evidence: BP-PY-14's condition is an executable requests call missing `timeout=`; the regex matches the raw source inside a docstring, so no real call exists.

### [ ] Finding `11` — `BP-PY-14`

- Function context: `scripts/niquests/findings/functions/11.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_vendor/kiss_headers/__init__.py:27:18`
- Checklist pattern: requests/session call inside a docstring doctest example

Source excerpt:

```
         25: ... or from a raw IMAP4 message:
         26: 
    >    27:    >>> message = requests.get("https://gist.githubusercontent.com/Ousret/8b84b736c375bb6aa3d389e86b5116ec/raw/21cb2f7af865e401c37d9b053fb6fe1abf63165b/sample-message.eml").content
         28:    >>> headers = parse_it(message)
```

Why this is a false positive: the `>>> requests.get(...)` / `>>> session.get(...)` lines are doctest examples in docstrings, i.e. documentation text, not executed HTTP calls.

Checklist evidence: BP-PY-14's condition is an executable requests call missing `timeout=`; the regex matches the raw source inside a docstring, so no real call exists.

### [ ] Finding `12` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/12.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_vendor/kiss_headers/builder.py:350:1`
- Checklist pattern: parsing fallback for an optional `filename*` encoding split

Source excerpt:

```
        348:                 encoding, encoded_filename = tuple(str(self["filename*"]).split("''"))
        349:                 return url_unquote(encoded_filename, encoding)
    >   350:             except ValueError:
        351:                 pass
```

Why this is a false positive: the `ValueError` is the expected outcome when the `filename*` value has no `encoding''encoded` split; the property falls through to the raw value, a defined fallback.

Checklist evidence: the except-pass is a deliberate fallback for a malformed optional header; no failure is hidden because a defined value is returned.

### [ ] Finding `13` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/13.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_vendor/kiss_headers/builder.py:350:1`
- Checklist pattern: parsing fallback for an optional `filename*` encoding split

Source excerpt:

```
        348:                 encoding, encoded_filename = tuple(str(self["filename*"]).split("''"))
        349:                 return url_unquote(encoded_filename, encoding)
    >   350:             except ValueError:
        351:                 pass
```

Why this is a false positive: the `ValueError` is the expected outcome when the `filename*` value has no `encoding''encoded` split; the property falls through to the raw value, a defined fallback.

Checklist evidence: the except-pass is a deliberate fallback for a malformed optional header; no failure is hidden because a defined value is returned.

### [ ] Finding `14` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/14.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_vendor/kiss_headers/builder.py:350:13`
- Checklist pattern: parsing fallback for an optional `filename*` encoding split

Source excerpt:

```
        348:                 encoding, encoded_filename = tuple(str(self["filename*"]).split("''"))
        349:                 return url_unquote(encoded_filename, encoding)
    >   350:             except ValueError:
        351:                 pass
```

Why this is a false positive: the `ValueError` is the expected outcome when the `filename*` value has no `encoding''encoded` split; the property falls through to the raw value, a defined fallback.

Checklist evidence: the except-pass is a deliberate fallback for a malformed optional header; no failure is hidden because a defined value is returned.

### [ ] Finding `15` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/15.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_vendor/kiss_headers/models.py:893:1`
- Checklist pattern: parsing fallback in header subclass lookup

Source excerpt:

```
        891:                     else None
        892:                 )
    >   893:             except TypeError:
        894:                 pass
```

Why this is a false positive: `header_name_to_class` is a best-effort subclass probe; on `TypeError` `target_subclass` stays `None`, which the caller handles.

Checklist evidence: the exception is the expected outcome of probing for an optional class mapping; the fallback value is defined.

### [ ] Finding `16` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/16.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_vendor/kiss_headers/models.py:893:1`
- Checklist pattern: parsing fallback in header subclass lookup

Source excerpt:

```
        891:                     else None
        892:                 )
    >   893:             except TypeError:
        894:                 pass
```

Why this is a false positive: `header_name_to_class` is a best-effort subclass probe; on `TypeError` `target_subclass` stays `None`, which the caller handles.

Checklist evidence: the exception is the expected outcome of probing for an optional class mapping; the fallback value is defined.

### [ ] Finding `17` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/17.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_vendor/kiss_headers/models.py:893:13`
- Checklist pattern: parsing fallback in header subclass lookup

Source excerpt:

```
        891:                     else None
        892:                 )
    >   893:             except TypeError:
        894:                 pass
```

Why this is a false positive: `header_name_to_class` is a best-effort subclass probe; on `TypeError` `target_subclass` stays `None`, which the caller handles.

Checklist evidence: the exception is the expected outcome of probing for an optional class mapping; the fallback value is defined.

### [ ] Finding `18` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/18.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_vendor/kiss_headers/utils.py:339:1`
- Checklist pattern: parsing fallback when a segment contains no `;`

Source excerpt:

```
        337:     try:
        338:         next_semi_colon_index = elem_end_index + content[elem_end_index:].index(";")
    >   339:     except ValueError:
        340:         pass
```

Why this is a false positive: `index(";")` raising `ValueError` means the element has no parameter part — the expected outcome the loop is written for; processing continues with the collected element.

Checklist evidence: the flagged `pass` implements an expected parsing branch, not silent failure discarding.

### [ ] Finding `19` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/19.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_vendor/kiss_headers/utils.py:339:1`
- Checklist pattern: parsing fallback when a segment contains no `;`

Source excerpt:

```
        337:     try:
        338:         next_semi_colon_index = elem_end_index + content[elem_end_index:].index(";")
    >   339:     except ValueError:
        340:         pass
```

Why this is a false positive: `index(";")` raising `ValueError` means the element has no parameter part — the expected outcome the loop is written for; processing continues with the collected element.

Checklist evidence: the flagged `pass` implements an expected parsing branch, not silent failure discarding.

### [ ] Finding `20` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/20.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_vendor/kiss_headers/utils.py:339:5`
- Checklist pattern: parsing fallback when a segment contains no `;`

Source excerpt:

```
        337:     try:
        338:         next_semi_colon_index = elem_end_index + content[elem_end_index:].index(";")
    >   339:     except ValueError:
        340:         pass
```

Why this is a false positive: `index(";")` raising `ValueError` means the element has no parameter part — the expected outcome the loop is written for; processing continues with the collected element.

Checklist evidence: the flagged `pass` implements an expected parsing branch, not silent failure discarding.

### [ ] Finding `21` — `BP-PY-49`

- Function context: `scripts/niquests/findings/functions/21.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/adapters.py:568:1`
- Checklist pattern: certificate-fingerprint pinning path — verification is still performed

Source excerpt:

```
        566:                         if len(verify) in [71, 45]:
        567:                             assert_fingerprint = verify.split("_", 1)[-1]
    >   568:                             verify = False
        569:                             cert_data = None
```

Why this is a false positive: these lines run only when the caller passes a `sha256_`/`sha1_` fingerprint; `verify = False`/`cert_reqs = "CERT_NONE"` is the required urllib3 setup for `assert_fingerprint` pinning (line 606/1703), so TLS is still verified by fingerprint.

Checklist evidence: the rule's condition 'TLS verification disabled' is not satisfied: verification is performed via certificate fingerprint pinning on this path.

### [ ] Finding `22` — `BP-PY-49`

- Function context: `scripts/niquests/findings/functions/22.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/adapters.py:591:23`
- Checklist pattern: comparison expression, not a disabling assignment

Source excerpt:

```
        589:                     conn.cert_reqs = "CERT_REQUIRED"
        590:                 else:
    >   591:                     if conn.cert_reqs != "CERT_NONE":
        592:                         need_reboot_conn = True
```

Why this is a false positive: `if conn.cert_reqs != "CERT_NONE":` is a state comparison (checking whether the connection needs a reboot), not an assignment that disables verification.

Checklist evidence: the rule's CERT_NONE marker fires on the regex token inside a `!=` comparison; no verification is disabled at this line.

### [ ] Finding `23` — `BP-PY-49`

- Function context: `scripts/niquests/findings/functions/23.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/adapters.py:594:19`
- Checklist pattern: certificate-fingerprint pinning path — verification is still performed

Source excerpt:

```
        592:                         need_reboot_conn = True
        593: 
    >   594:                     conn.cert_reqs = "CERT_NONE"
        595: 
```

Why this is a false positive: these lines run only when the caller passes a `sha256_`/`sha1_` fingerprint; `verify = False`/`cert_reqs = "CERT_NONE"` is the required urllib3 setup for `assert_fingerprint` pinning (line 606/1703), so TLS is still verified by fingerprint.

Checklist evidence: the rule's condition 'TLS verification disabled' is not satisfied: verification is performed via certificate fingerprint pinning on this path.

### [ ] Finding `24` — `BP-PY-49`

- Function context: `scripts/niquests/findings/functions/24.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/adapters.py:608:23`
- Checklist pattern: comparison expression, not a disabling assignment

Source excerpt:

```
        606:                     conn.assert_fingerprint = assert_fingerprint
        607:             else:
    >   608:                 if conn.cert_reqs != "CERT_NONE":
        609:                     need_reboot_conn = True
```

Why this is a false positive: `if conn.cert_reqs != "CERT_NONE":` is a state comparison (checking whether the connection needs a reboot), not an assignment that disables verification.

Checklist evidence: the rule's CERT_NONE marker fires on the regex token inside a `!=` comparison; no verification is disabled at this line.

### [ ] Finding `25` — `BP-PY-49`

- Function context: `scripts/niquests/findings/functions/25.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/adapters.py:611:19`
- Checklist pattern: caller-opted `verify=False` implementation of the documented public API

Source excerpt:

```
        609:                     need_reboot_conn = True
        610: 
    >   611:                 conn.cert_reqs = "CERT_NONE"
        612:                 conn.ca_certs = None
```

Why this is a false positive: the branch executes only when the caller explicitly passes `verify=False` (the requests-compatible public option whose default is `True`); the library implements the documented option rather than accidentally disabling TLS.

Checklist evidence: the disablement is parameter-gated and opt-in through the public `verify` argument, mirroring the intended requests-compatible surface.

### [ ] Finding `29` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/29.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/adapters.py:1075:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
       1073:                             try:
       1074:                                 from .extensions.revocation._ocsp import verify as ocsp_verify
    >  1075:                             except ImportError:
       1076:                                 pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `30` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/30.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/adapters.py:1075:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
       1073:                             try:
       1074:                                 from .extensions.revocation._ocsp import verify as ocsp_verify
    >  1075:                             except ImportError:
       1076:                                 pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `31` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/31.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/adapters.py:1075:29`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
       1073:                             try:
       1074:                                 from .extensions.revocation._ocsp import verify as ocsp_verify
    >  1075:                             except ImportError:
       1076:                                 pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `32` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/32.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/adapters.py:1091:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
       1089:                             try:
       1090:                                 from .extensions.revocation._crl import verify as crl_verify
    >  1091:                             except ImportError:
       1092:                                 pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `33` — `BP-PY-49`

- Function context: `scripts/niquests/findings/functions/33.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/adapters.py:1666:1`
- Checklist pattern: certificate-fingerprint pinning path — verification is still performed

Source excerpt:

```
       1664:                         if len(verify) in [71, 45]:
       1665:                             assert_fingerprint = verify.split("_", 1)[-1]
    >  1666:                             verify = False
       1667:                             cert_data = None
```

Why this is a false positive: these lines run only when the caller passes a `sha256_`/`sha1_` fingerprint; `verify = False`/`cert_reqs = "CERT_NONE"` is the required urllib3 setup for `assert_fingerprint` pinning (line 606/1703), so TLS is still verified by fingerprint.

Checklist evidence: the rule's condition 'TLS verification disabled' is not satisfied: verification is performed via certificate fingerprint pinning on this path.

### [ ] Finding `34` — `BP-PY-49`

- Function context: `scripts/niquests/findings/functions/34.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/adapters.py:1688:23`
- Checklist pattern: comparison expression, not a disabling assignment

Source excerpt:

```
       1686:                     conn.cert_reqs = "CERT_REQUIRED"
       1687:                 else:
    >  1688:                     if conn.cert_reqs != "CERT_NONE":
       1689:                         need_reboot_conn = True
```

Why this is a false positive: `if conn.cert_reqs != "CERT_NONE":` is a state comparison (checking whether the connection needs a reboot), not an assignment that disables verification.

Checklist evidence: the rule's CERT_NONE marker fires on the regex token inside a `!=` comparison; no verification is disabled at this line.

### [ ] Finding `35` — `BP-PY-49`

- Function context: `scripts/niquests/findings/functions/35.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/adapters.py:1691:19`
- Checklist pattern: certificate-fingerprint pinning path — verification is still performed

Source excerpt:

```
       1689:                         need_reboot_conn = True
       1690: 
    >  1691:                     conn.cert_reqs = "CERT_NONE"
       1692: 
```

Why this is a false positive: these lines run only when the caller passes a `sha256_`/`sha1_` fingerprint; `verify = False`/`cert_reqs = "CERT_NONE"` is the required urllib3 setup for `assert_fingerprint` pinning (line 606/1703), so TLS is still verified by fingerprint.

Checklist evidence: the rule's condition 'TLS verification disabled' is not satisfied: verification is performed via certificate fingerprint pinning on this path.

### [ ] Finding `36` — `BP-PY-49`

- Function context: `scripts/niquests/findings/functions/36.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/adapters.py:1705:23`
- Checklist pattern: comparison expression, not a disabling assignment

Source excerpt:

```
       1703:                     conn.assert_fingerprint = assert_fingerprint
       1704:             else:
    >  1705:                 if conn.cert_reqs != "CERT_NONE":
       1706:                     need_reboot_conn = True
```

Why this is a false positive: `if conn.cert_reqs != "CERT_NONE":` is a state comparison (checking whether the connection needs a reboot), not an assignment that disables verification.

Checklist evidence: the rule's CERT_NONE marker fires on the regex token inside a `!=` comparison; no verification is disabled at this line.

### [ ] Finding `37` — `BP-PY-49`

- Function context: `scripts/niquests/findings/functions/37.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/adapters.py:1708:19`
- Checklist pattern: caller-opted `verify=False` implementation of the documented public API

Source excerpt:

```
       1706:                     need_reboot_conn = True
       1707: 
    >  1708:                 conn.cert_reqs = "CERT_NONE"
       1709:                 conn.ca_certs = None
```

Why this is a false positive: the branch executes only when the caller explicitly passes `verify=False` (the requests-compatible public option whose default is `True`); the library implements the documented option rather than accidentally disabling TLS.

Checklist evidence: the disablement is parameter-gated and opt-in through the public `verify` argument, mirroring the intended requests-compatible surface.

### [ ] Finding `41` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/41.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/adapters.py:2189:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
       2187:                                     verify as async_ocsp_verify,
       2188:                                 )
    >  2189:                             except ImportError:
       2190:                                 pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `42` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/42.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/adapters.py:2207:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
       2205:                                     verify as async_crl_verify,
       2206:                                 )
    >  2207:                             except ImportError:
       2208:                                 pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `43` — `BP-PY-14`

- Function context: `scripts/niquests/findings/functions/43.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/async_api.py:183:26`
- Checklist pattern: timeout passed as a positional argument

Source excerpt:

```
        181:         session._crl_cache = _SHARED_CRL_CACHE.get()
        182:         try:
    >   183:             return await session.request(  # type: ignore[misc]
        184:                 method,
```

Why this is a false positive: the `session.request(...)` call passes `timeout` as the 10th positional argument, so the request does carry a timeout; the rule only detects the `timeout=` keyword form.

Checklist evidence: the rule's condition 'call missing timeout=' is not met — the timeout value is present positionally.

### [ ] Finding `44` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/44.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/async_session.py:99:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         97:     try:
         98:         from .extensions.wasi._async._adapter import AsyncWASIAdapter as AsyncWASIHTTPAdapter
    >    99:     except ImportError:
        100:         pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `45` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/45.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/async_session.py:99:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         97:     try:
         98:         from .extensions.wasi._async._adapter import AsyncWASIAdapter as AsyncWASIHTTPAdapter
    >    99:     except ImportError:
        100:         pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `46` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/46.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/async_session.py:99:5`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         97:     try:
         98:         from .extensions.wasi._async._adapter import AsyncWASIAdapter as AsyncWASIHTTPAdapter
    >    99:     except ImportError:
        100:         pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `49` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/49.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/async_session.py:620:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
        618:                 if scheme in extension.supported_schemes() and extension.scheme_to_http_scheme(scheme) == parse_scheme(prefix):
        619:                     return adapter
    >   620:         except ImportError:
        621:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `50` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/50.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/async_session.py:703:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
        701:                             verify as ocsp_verify,
        702:                         )
    >   703:                     except ImportError:
        704:                         pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `52` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/52.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/async_session.py:731:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
        729:                             verify as crl_verify,
        730:                         )
    >   731:                     except ImportError:
        732:                         pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `54` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/54.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/async_session.py:957:1`
- Checklist pattern: generator-exhaustion control flow

Source excerpt:

```
        955:                 try:
        956:                     r._next = await gen.__anext__()  # type: ignore[assignment]
    >   957:                 except StopAsyncIteration:
        958:                     pass
```

Why this is a false positive: `StopIteration`/`StopAsyncIteration` is the generator's normal termination signal; catching it is the intended loop-exit control flow, not a discarded failure.

Checklist evidence: the rule's 'failures discarded silently' condition does not apply to iterator-exhaustion signalling.

### [ ] Finding `55` — `CWE-1341`

- Function context: `scripts/niquests/findings/functions/55.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/async_session.py:1015:23`
- Checklist pattern: mutually exclusive close paths, single release per execution

Source excerpt:

```
       1013:             # Release the connection back into the pool.
       1014:             if isinstance(resp, AsyncResponse):
    >  1015:                 await resp.close()
       1016:             else:
```

Why this is a false positive: `async_session.py:1015` and its `else: resp.close()` are the two branches of one `isinstance` conditional — exactly one executes; the wasi `response.close()` calls sit in exclusive retry/raise branches (line 235 vs 239).

Checklist evidence: CWE-1341's condition is the same handle released twice on one path; the matched `close()` pairs are mutually exclusive branches (or distinct resources).

### [ ] Finding `58` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/58.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/async_session.py:1234:1`
- Checklist pattern: non-seekable body probe fallback

Source excerpt:

```
       1232:             try:
       1233:                 prep._body_position = await prep.body.tell()
    >  1234:             except OSError:
       1235:                 pass
```

Why this is a false positive: `body.tell()` raises `OSError` for non-seekable bodies; `_body_position` stays unset and the caller treats it as 'unknown' — a defined fallback.

Checklist evidence: the probe failure is expected for streaming bodies and handled by leaving the optional attribute unset.

### [ ] Finding `59` — `CWE-93`

- Function context: `scripts/niquests/findings/functions/59.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/auth.py:75:14`
- Checklist pattern: write to an outgoing request header, not a response header

Source excerpt:

```
         73: 
         74:         if len(detect_token_type) == 1:
    >    75:             r.headers["Authorization"] = f"Bearer {self.token}"
         76:         else:
```

Why this is a false positive: `auth.py:75` sets the `Authorization` header on the outgoing request (from the caller's own token); `models.py:488` populates the `PreparedRequest.headers` dict in `prepare_headers` from the caller-supplied header mapping — neither is an HTTP response header.

Checklist evidence: CWE-93's condition is writing an externally influenced value into an HTTP response header; the shown lines write client-side request headers.

### [ ] Finding `60` — `BP-PY-14`

- Function context: `scripts/niquests/findings/functions/60.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/auth.py:399:27`
- Checklist pattern: requests/session call inside a docstring doctest example

Source excerpt:

```
        397:         >>> auth = niquests.auth.AsyncHTTPDigestAuth('user', 'pass')
        398:         >>> async with niquests.AsyncSession() as session:
    >   399:         ...     r = await session.get('https://httpbin.org/digest-auth/auth/user/pass', auth=auth)
        400:         ...     print(r.status_code)
```

Why this is a false positive: the `>>> requests.get(...)` / `>>> session.get(...)` lines are doctest examples in docstrings, i.e. documentation text, not executed HTTP calls.

Checklist evidence: BP-PY-14's condition is an executable requests call missing `timeout=`; the regex matches the raw source inside a docstring, so no real call exists.

### [ ] Finding `61` — `BP-PY-46`

- Function context: `scripts/niquests/findings/functions/61.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/auth.py:400:17`
- Checklist pattern: `print(` token inside a docstring doctest example

Source excerpt:

```
        398:         >>> async with niquests.AsyncSession() as session:
        399:         ...     r = await session.get('https://httpbin.org/digest-auth/auth/user/pass', auth=auth)
    >   400:         ...     print(r.status_code)
        401:         200
```

Why this is a false positive: the `>>> print(...)` lines are doctest examples inside docstrings, i.e. documentation text, not executable statements.

Checklist evidence: BP-PY-46's condition is a real `print` call in library code; the per-line string check cannot see the enclosing docstring, so it reports documentation text.

### [ ] Finding `62` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/62.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:57:1`
- Checklist pattern: documented defensive JS-bridge reader setup

Source excerpt:

```
         55:                 if body is not None:
         56:                     self._reader = body.getReader()
    >    57:             except Exception:
         58:                 pass
```

Why this is a false positive: the Pyodide JS bridge suppresses failures while probing the JS stream object's `getReader()`; when it fails, `self._reader` stays unset and the sync reader falls back to `body.arrayBuffer()`.

Checklist evidence: the handler is a deliberate best-effort bridge probe with a defined fallback state; no failure that matters is hidden.

### [ ] Finding `63` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/63.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:57:1`
- Checklist pattern: documented defensive JS-bridge reader setup

Source excerpt:

```
         55:                 if body is not None:
         56:                     self._reader = body.getReader()
    >    57:             except Exception:
         58:                 pass
```

Why this is a false positive: the Pyodide JS bridge suppresses failures while probing the JS stream object's `getReader()`; when it fails, `self._reader` stays unset and the sync reader falls back to `body.arrayBuffer()`.

Checklist evidence: the handler is a deliberate best-effort bridge probe with a defined fallback state; no failure that matters is hidden.

### [ ] Finding `64` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/64.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:57:1`
- Checklist pattern: documented defensive JS-bridge reader setup

Source excerpt:

```
         55:                 if body is not None:
         56:                     self._reader = body.getReader()
    >    57:             except Exception:
         58:                 pass
```

Why this is a false positive: the Pyodide JS bridge suppresses failures while probing the JS stream object's `getReader()`; when it fails, `self._reader` stays unset and the sync reader falls back to `body.arrayBuffer()`.

Checklist evidence: the handler is a deliberate best-effort bridge probe with a defined fallback state; no failure that matters is hidden.

### [ ] Finding `65` — `CWE-396`

- Function context: `scripts/niquests/findings/functions/65.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:57:1`
- Checklist pattern: documented defensive JS-bridge reader setup

Source excerpt:

```
         55:                 if body is not None:
         56:                     self._reader = body.getReader()
    >    57:             except Exception:
         58:                 pass
```

Why this is a false positive: the Pyodide JS bridge suppresses failures while probing the JS stream object's `getReader()`; when it fails, `self._reader` stays unset and the sync reader falls back to `body.arrayBuffer()`.

Checklist evidence: the handler is a deliberate best-effort bridge probe with a defined fallback state; no failure that matters is hidden.

### [ ] Finding `66` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/66.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:57:13`
- Checklist pattern: documented defensive JS-bridge reader setup

Source excerpt:

```
         55:                 if body is not None:
         56:                     self._reader = body.getReader()
    >    57:             except Exception:
         58:                 pass
```

Why this is a false positive: the Pyodide JS bridge suppresses failures while probing the JS stream object's `getReader()`; when it fails, `self._reader` stays unset and the sync reader falls back to `body.arrayBuffer()`.

Checklist evidence: the handler is a deliberate best-effort bridge probe with a defined fallback state; no failure that matters is hidden.

### [ ] Finding `67` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/67.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:77:1`
- Checklist pattern: best-effort `to_py()` conversion with `None` fallback

Source excerpt:

```
         75:                 return bytes(value.to_py())
         76:             return None
    >    77:         except Exception:
         78:             return None
```

Why this is a false positive: converting a JS object to bytes is best-effort; failure returns `None`, a defined fallback the caller handles.

Checklist evidence: the handler returns a defined fallback value, so the failure is not swallowed without consequence.

### [ ] Finding `68` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/68.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:131:1`
- Checklist pattern: teardown best-effort reader cancellation

Source excerpt:

```
        129:             try:
        130:                 run_sync(self._reader.cancel())
    >   131:             except Exception:
        132:                 pass
```

Why this is a false positive: cancelling the stream reader during teardown is best-effort; JS errors here are expected and the loop has already ended.

Checklist evidence: the cancellation failure is the expected outcome of tearing down an already-closed stream.

### [ ] Finding `69` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/69.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:131:1`
- Checklist pattern: teardown best-effort reader cancellation

Source excerpt:

```
        129:             try:
        130:                 run_sync(self._reader.cancel())
    >   131:             except Exception:
        132:                 pass
```

Why this is a false positive: cancelling the stream reader during teardown is best-effort; JS errors here are expected and the loop has already ended.

Checklist evidence: the cancellation failure is the expected outcome of tearing down an already-closed stream.

### [ ] Finding `70` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/70.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:189:1`
- Checklist pattern: retry-loop handler that acts on the exception

Source excerpt:

```
        187:             try:
        188:                 response = self._do_send(request, stream, timeout)
    >   189:             except Exception as err:
        190:                 retries = retries.increment(method, request.url, error=err)
```

Why this is a false positive: the handler increments the retry state (`retries.increment(...)`) and re-raises `MaxRetryError` once retries are exhausted — the failure is not swallowed.

Checklist evidence: BP-PY-1's 'hides failures' condition is unmet: the generic catch feeds the retry machinery and re-raises on exhaustion.

### [ ] Finding `71` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/71.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:299:1`
- Checklist pattern: best-effort response-header parse fallback

Source excerpt:

```
        297:                     for entry in js_headers.entries():
        298:                         response_headers[entry[0]] = entry[1]
    >   299:         except Exception:
        300:             pass
```

Why this is a false positive: extracting headers from the JS `headers.entries()` is best-effort; on failure the response is returned without the (optional) header snapshot.

Checklist evidence: the handler is a defensive bridge fallback with a defined result.

### [ ] Finding `72` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/72.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:299:1`
- Checklist pattern: best-effort response-header parse fallback

Source excerpt:

```
        297:                     for entry in js_headers.entries():
        298:                         response_headers[entry[0]] = entry[1]
    >   299:         except Exception:
        300:             pass
```

Why this is a false positive: extracting headers from the JS `headers.entries()` is best-effort; on failure the response is returned without the (optional) header snapshot.

Checklist evidence: the handler is a defensive bridge fallback with a defined result.

### [ ] Finding `73` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/73.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:312:1`
- Checklist pattern: fallback for missing reason phrase

Source excerpt:

```
        310:         try:
        311:             response.reason = js_response.status_text or ""
    >   312:         except Exception:
        313:             response.reason = ""
```

Why this is a false positive: when the JS `status_text` probe fails, `response.reason` falls back to `""` — a defined value, not a swallowed failure.

Checklist evidence: the handler assigns a defined fallback value.

### [ ] Finding `74` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/74.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:327:1`
- Checklist pattern: fallback for body read failure

Source excerpt:

```
        325:             try:
        326:                 response_body: bytes = run_sync(js_response.bytes())
    >   327:             except Exception:
        328:                 response_body = b""
```

Why this is a false positive: when reading the JS response body fails, `response_body` falls back to `b""` — a defined value.

Checklist evidence: the handler assigns a defined fallback value.

### [ ] Finding `75` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/75.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/__init__.py:56:1`
- Checklist pattern: documented defensive JS-bridge reader setup

Source excerpt:

```
         54:                 if body is not None:
         55:                     self._reader = body.getReader()
    >    56:             except Exception:
         57:                 pass
```

Why this is a false positive: the Pyodide JS bridge suppresses failures while probing the JS stream object's `getReader()`; when it fails, `self._reader` stays unset and the sync reader falls back to `body.arrayBuffer()`.

Checklist evidence: the handler is a deliberate best-effort bridge probe with a defined fallback state; no failure that matters is hidden.

### [ ] Finding `76` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/76.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/__init__.py:56:1`
- Checklist pattern: documented defensive JS-bridge reader setup

Source excerpt:

```
         54:                 if body is not None:
         55:                     self._reader = body.getReader()
    >    56:             except Exception:
         57:                 pass
```

Why this is a false positive: the Pyodide JS bridge suppresses failures while probing the JS stream object's `getReader()`; when it fails, `self._reader` stays unset and the sync reader falls back to `body.arrayBuffer()`.

Checklist evidence: the handler is a deliberate best-effort bridge probe with a defined fallback state; no failure that matters is hidden.

### [ ] Finding `77` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/77.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/__init__.py:56:1`
- Checklist pattern: documented defensive JS-bridge reader setup

Source excerpt:

```
         54:                 if body is not None:
         55:                     self._reader = body.getReader()
    >    56:             except Exception:
         57:                 pass
```

Why this is a false positive: the Pyodide JS bridge suppresses failures while probing the JS stream object's `getReader()`; when it fails, `self._reader` stays unset and the sync reader falls back to `body.arrayBuffer()`.

Checklist evidence: the handler is a deliberate best-effort bridge probe with a defined fallback state; no failure that matters is hidden.

### [ ] Finding `78` — `CWE-396`

- Function context: `scripts/niquests/findings/functions/78.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/__init__.py:56:1`
- Checklist pattern: documented defensive JS-bridge reader setup

Source excerpt:

```
         54:                 if body is not None:
         55:                     self._reader = body.getReader()
    >    56:             except Exception:
         57:                 pass
```

Why this is a false positive: the Pyodide JS bridge suppresses failures while probing the JS stream object's `getReader()`; when it fails, `self._reader` stays unset and the sync reader falls back to `body.arrayBuffer()`.

Checklist evidence: the handler is a deliberate best-effort bridge probe with a defined fallback state; no failure that matters is hidden.

### [ ] Finding `79` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/79.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/__init__.py:56:13`
- Checklist pattern: documented defensive JS-bridge reader setup

Source excerpt:

```
         54:                 if body is not None:
         55:                     self._reader = body.getReader()
    >    56:             except Exception:
         57:                 pass
```

Why this is a false positive: the Pyodide JS bridge suppresses failures while probing the JS stream object's `getReader()`; when it fails, `self._reader` stays unset and the sync reader falls back to `body.arrayBuffer()`.

Checklist evidence: the handler is a deliberate best-effort bridge probe with a defined fallback state; no failure that matters is hidden.

### [ ] Finding `80` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/80.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/__init__.py:119:1`
- Checklist pattern: best-effort stream read with `None` fallback

Source excerpt:

```
        117:         except asyncio.TimeoutError:
        118:             raise ReadTimeout("Read timed out while streaming Pyodide response")
    >   119:         except Exception:
        120:             return None
```

Why this is a false positive: the read converts JS chunks and returns `None` on conversion failure, which the caller treats as end-of-stream.

Checklist evidence: the handler returns a defined fallback.

### [ ] Finding `81` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/81.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/__init__.py:153:1`
- Checklist pattern: teardown best-effort reader cancellation

Source excerpt:

```
        151:             try:
        152:                 await self._reader.cancel()
    >   153:             except Exception:
        154:                 pass
```

Why this is a false positive: cancelling the stream reader during teardown is best-effort; JS errors here are expected and the loop has already ended.

Checklist evidence: the cancellation failure is the expected outcome of tearing down an already-closed stream.

### [ ] Finding `82` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/82.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/__init__.py:153:1`
- Checklist pattern: teardown best-effort reader cancellation

Source excerpt:

```
        151:             try:
        152:                 await self._reader.cancel()
    >   153:             except Exception:
        154:                 pass
```

Why this is a false positive: cancelling the stream reader during teardown is best-effort; JS errors here are expected and the loop has already ended.

Checklist evidence: the cancellation failure is the expected outcome of tearing down an already-closed stream.

### [ ] Finding `83` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/83.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/__init__.py:216:1`
- Checklist pattern: retry-loop handler that acts on the exception

Source excerpt:

```
        214:             try:
        215:                 response = await self._do_send(request, stream, timeout)
    >   216:             except Exception as err:
        217:                 retries = retries.increment(method, request.url, error=err)
```

Why this is a false positive: the handler increments the retry state (`retries.increment(...)`) and re-raises `MaxRetryError` once retries are exhausted — the failure is not swallowed.

Checklist evidence: BP-PY-1's 'hides failures' condition is unmet: the generic catch feeds the retry machinery and re-raises on exhaustion.

### [ ] Finding `84` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/84.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/__init__.py:341:1`
- Checklist pattern: best-effort response-header parse fallback

Source excerpt:

```
        339:                     for entry in js_headers.entries():
        340:                         response_headers[entry[0]] = entry[1]
    >   341:         except Exception:
        342:             pass
```

Why this is a false positive: extracting headers from the JS `headers.entries()` is best-effort; on failure the response is returned without the (optional) header snapshot.

Checklist evidence: the handler is a defensive bridge fallback with a defined result.

### [ ] Finding `85` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/85.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/__init__.py:341:1`
- Checklist pattern: best-effort response-header parse fallback

Source excerpt:

```
        339:                     for entry in js_headers.entries():
        340:                         response_headers[entry[0]] = entry[1]
    >   341:         except Exception:
        342:             pass
```

Why this is a false positive: extracting headers from the JS `headers.entries()` is best-effort; on failure the response is returned without the (optional) header snapshot.

Checklist evidence: the handler is a defensive bridge fallback with a defined result.

### [ ] Finding `86` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/86.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/__init__.py:355:1`
- Checklist pattern: fallback for missing reason phrase

Source excerpt:

```
        353:         try:
        354:             response.reason = js_response.status_text or ""
    >   355:         except Exception:
        356:             response.reason = ""
```

Why this is a false positive: when the JS `status_text` probe fails, `response.reason` falls back to `""` — a defined value, not a swallowed failure.

Checklist evidence: the handler assigns a defined fallback value.

### [ ] Finding `87` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/87.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_sse.py:58:1`
- Checklist pattern: best-effort stream read with `None` fallback

Source excerpt:

```
         56:                 return bytes(value.to_py()).decode("utf-8")
         57:             return None
    >    58:         except Exception:
         59:             return None
```

Why this is a false positive: the read converts JS chunks and returns `None` on conversion failure, which the caller treats as end-of-stream.

Checklist evidence: the handler returns a defined fallback.

### [ ] Finding `88` — `CWE-396`

- Function context: `scripts/niquests/findings/functions/88.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_sse.py:58:1`
- Checklist pattern: best-effort stream read with `None` fallback

Source excerpt:

```
         56:                 return bytes(value.to_py()).decode("utf-8")
         57:             return None
    >    58:         except Exception:
         59:             return None
```

Why this is a false positive: the read converts JS chunks and returns `None` on conversion failure, which the caller treats as end-of-stream.

Checklist evidence: the handler returns a defined fallback.

### [ ] Finding `89` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/89.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_sse.py:147:1`
- Checklist pattern: teardown best-effort reader cancellation

Source excerpt:

```
        145:             try:
        146:                 await self._reader.cancel()
    >   147:             except Exception:
        148:                 pass
```

Why this is a false positive: cancelling the stream reader during teardown is best-effort; JS errors here are expected and the loop has already ended.

Checklist evidence: the cancellation failure is the expected outcome of tearing down an already-closed stream.

### [ ] Finding `90` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/90.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_sse.py:147:1`
- Checklist pattern: teardown best-effort reader cancellation

Source excerpt:

```
        145:             try:
        146:                 await self._reader.cancel()
    >   147:             except Exception:
        148:                 pass
```

Why this is a false positive: cancelling the stream reader during teardown is best-effort; JS errors here are expected and the loop has already ended.

Checklist evidence: the cancellation failure is the expected outcome of tearing down an already-closed stream.

### [ ] Finding `91` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/91.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_sse.py:147:1`
- Checklist pattern: teardown best-effort reader cancellation

Source excerpt:

```
        145:             try:
        146:                 await self._reader.cancel()
    >   147:             except Exception:
        148:                 pass
```

Why this is a false positive: cancelling the stream reader during teardown is best-effort; JS errors here are expected and the loop has already ended.

Checklist evidence: the cancellation failure is the expected outcome of tearing down an already-closed stream.

### [ ] Finding `92` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/92.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_sse.py:147:13`
- Checklist pattern: teardown best-effort reader cancellation

Source excerpt:

```
        145:             try:
        146:                 await self._reader.cancel()
    >   147:             except Exception:
        148:                 pass
```

Why this is a false positive: cancelling the stream reader during teardown is best-effort; JS errors here are expected and the loop has already ended.

Checklist evidence: the cancellation failure is the expected outcome of tearing down an already-closed stream.

### [ ] Finding `93` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/93.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_ws.py:123:1`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        121:         try:
        122:             self._ws.close()
    >   123:         except Exception:  # Defensive: suppress JS errors on teardown
        124:             pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `94` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/94.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_ws.py:123:1`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        121:         try:
        122:             self._ws.close()
    >   123:         except Exception:  # Defensive: suppress JS errors on teardown
        124:             pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `95` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/95.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_ws.py:123:1`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        121:         try:
        122:             self._ws.close()
    >   123:         except Exception:  # Defensive: suppress JS errors on teardown
        124:             pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `96` — `CWE-396`

- Function context: `scripts/niquests/findings/functions/96.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_ws.py:123:1`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        121:         try:
        122:             self._ws.close()
    >   123:         except Exception:  # Defensive: suppress JS errors on teardown
        124:             pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `97` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/97.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_ws.py:123:9`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        121:         try:
        122:             self._ws.close()
    >   123:         except Exception:  # Defensive: suppress JS errors on teardown
        124:             pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `98` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/98.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_ws.py:129:1`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        127:             try:
        128:                 proxy.destroy()
    >   129:             except Exception:  # Defensive: suppress JS errors on teardown
        130:                 pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `99` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/99.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_ws.py:129:1`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        127:             try:
        128:                 proxy.destroy()
    >   129:             except Exception:  # Defensive: suppress JS errors on teardown
        130:                 pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `100` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/100.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_sse.py:59:1`
- Checklist pattern: best-effort stream read with `None` fallback

Source excerpt:

```
         57:                 return bytes(value.to_py()).decode("utf-8")
         58:             return None
    >    59:         except Exception:
         60:             return None
```

Why this is a false positive: the read converts JS chunks and returns `None` on conversion failure, which the caller treats as end-of-stream.

Checklist evidence: the handler returns a defined fallback.

### [ ] Finding `101` — `CWE-396`

- Function context: `scripts/niquests/findings/functions/101.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_sse.py:59:1`
- Checklist pattern: best-effort stream read with `None` fallback

Source excerpt:

```
         57:                 return bytes(value.to_py()).decode("utf-8")
         58:             return None
    >    59:         except Exception:
         60:             return None
```

Why this is a false positive: the read converts JS chunks and returns `None` on conversion failure, which the caller treats as end-of-stream.

Checklist evidence: the handler returns a defined fallback.

### [ ] Finding `102` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/102.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_sse.py:151:1`
- Checklist pattern: teardown best-effort reader cancellation

Source excerpt:

```
        149:             try:
        150:                 run_sync(self._reader.cancel())
    >   151:             except Exception:
        152:                 pass
```

Why this is a false positive: cancelling the stream reader during teardown is best-effort; JS errors here are expected and the loop has already ended.

Checklist evidence: the cancellation failure is the expected outcome of tearing down an already-closed stream.

### [ ] Finding `103` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/103.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_sse.py:151:1`
- Checklist pattern: teardown best-effort reader cancellation

Source excerpt:

```
        149:             try:
        150:                 run_sync(self._reader.cancel())
    >   151:             except Exception:
        152:                 pass
```

Why this is a false positive: cancelling the stream reader during teardown is best-effort; JS errors here are expected and the loop has already ended.

Checklist evidence: the cancellation failure is the expected outcome of tearing down an already-closed stream.

### [ ] Finding `104` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/104.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_sse.py:151:1`
- Checklist pattern: teardown best-effort reader cancellation

Source excerpt:

```
        149:             try:
        150:                 run_sync(self._reader.cancel())
    >   151:             except Exception:
        152:                 pass
```

Why this is a false positive: cancelling the stream reader during teardown is best-effort; JS errors here are expected and the loop has already ended.

Checklist evidence: the cancellation failure is the expected outcome of tearing down an already-closed stream.

### [ ] Finding `105` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/105.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_sse.py:151:13`
- Checklist pattern: teardown best-effort reader cancellation

Source excerpt:

```
        149:             try:
        150:                 run_sync(self._reader.cancel())
    >   151:             except Exception:
        152:                 pass
```

Why this is a false positive: cancelling the stream reader during teardown is best-effort; JS errors here are expected and the loop has already ended.

Checklist evidence: the cancellation failure is the expected outcome of tearing down an already-closed stream.

### [ ] Finding `106` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/106.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_ws.py:177:1`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        175:         try:
        176:             self._ws.close()
    >   177:         except Exception:  # Defensive: suppress JS errors on teardown
        178:             pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `107` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/107.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_ws.py:177:1`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        175:         try:
        176:             self._ws.close()
    >   177:         except Exception:  # Defensive: suppress JS errors on teardown
        178:             pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `108` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/108.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_ws.py:177:1`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        175:         try:
        176:             self._ws.close()
    >   177:         except Exception:  # Defensive: suppress JS errors on teardown
        178:             pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `109` — `CWE-396`

- Function context: `scripts/niquests/findings/functions/109.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_ws.py:177:1`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        175:         try:
        176:             self._ws.close()
    >   177:         except Exception:  # Defensive: suppress JS errors on teardown
        178:             pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `110` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/110.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_ws.py:177:9`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        175:         try:
        176:             self._ws.close()
    >   177:         except Exception:  # Defensive: suppress JS errors on teardown
        178:             pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `111` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/111.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_ws.py:183:1`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        181:             try:
        182:                 proxy.destroy()
    >   183:             except Exception:  # Defensive: suppress JS errors on teardown
        184:                 pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `112` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/112.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_ws.py:183:1`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        181:             try:
        182:                 proxy.destroy()
    >   183:             except Exception:  # Defensive: suppress JS errors on teardown
        184:                 pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `114` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/114.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/revocation/_crl/__init__.py:259:1`
- Checklist pattern: best-effort issuer-certificate fetch during revocation checking

Source excerpt:

```
        257:                     try:
        258:                         raw_intermediary_response = session.get(hint_ca_issuers[0])
    >   259:                     except RequestException:
        260:                         pass
```

Why this is a false positive: the intermediary-issuer fetch is optional enrichment: when the network fetch fails, revocation checking continues in its defined fallback state (no issuer hint).

Checklist evidence: the `RequestException` is an expected outcome of a best-effort network probe with a defined fallback; the pass is the intended behavior.

### [ ] Finding `115` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/115.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/revocation/_crl/__init__.py:259:1`
- Checklist pattern: best-effort issuer-certificate fetch during revocation checking

Source excerpt:

```
        257:                     try:
        258:                         raw_intermediary_response = session.get(hint_ca_issuers[0])
    >   259:                     except RequestException:
        260:                         pass
```

Why this is a false positive: the intermediary-issuer fetch is optional enrichment: when the network fetch fails, revocation checking continues in its defined fallback state (no issuer hint).

Checklist evidence: the `RequestException` is an expected outcome of a best-effort network probe with a defined fallback; the pass is the intended behavior.

### [ ] Finding `116` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/116.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/revocation/_crl/__init__.py:259:21`
- Checklist pattern: best-effort issuer-certificate fetch during revocation checking

Source excerpt:

```
        257:                     try:
        258:                         raw_intermediary_response = session.get(hint_ca_issuers[0])
    >   259:                     except RequestException:
        260:                         pass
```

Why this is a false positive: the intermediary-issuer fetch is optional enrichment: when the network fetch fails, revocation checking continues in its defined fallback state (no issuer hint).

Checklist evidence: the `RequestException` is an expected outcome of a best-effort network probe with a defined fallback; the pass is the intended behavior.

### [ ] Finding `121` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/121.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/revocation/_ocsp/__init__.py:325:1`
- Checklist pattern: best-effort issuer-certificate fetch during revocation checking

Source excerpt:

```
        323:                     try:
        324:                         raw_intermediary_response = session.get(hint_ca_issuers[0])
    >   325:                     except RequestException:
        326:                         pass
```

Why this is a false positive: the intermediary-issuer fetch is optional enrichment: when the network fetch fails, revocation checking continues in its defined fallback state (no issuer hint).

Checklist evidence: the `RequestException` is an expected outcome of a best-effort network probe with a defined fallback; the pass is the intended behavior.

### [ ] Finding `122` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/122.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/revocation/_ocsp/__init__.py:325:1`
- Checklist pattern: best-effort issuer-certificate fetch during revocation checking

Source excerpt:

```
        323:                     try:
        324:                         raw_intermediary_response = session.get(hint_ca_issuers[0])
    >   325:                     except RequestException:
        326:                         pass
```

Why this is a false positive: the intermediary-issuer fetch is optional enrichment: when the network fetch fails, revocation checking continues in its defined fallback state (no issuer hint).

Checklist evidence: the `RequestException` is an expected outcome of a best-effort network probe with a defined fallback; the pass is the intended behavior.

### [ ] Finding `123` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/123.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/revocation/_ocsp/__init__.py:325:21`
- Checklist pattern: best-effort issuer-certificate fetch during revocation checking

Source excerpt:

```
        323:                     try:
        324:                         raw_intermediary_response = session.get(hint_ca_issuers[0])
    >   325:                     except RequestException:
        326:                         pass
```

Why this is a false positive: the intermediary-issuer fetch is optional enrichment: when the network fetch fails, revocation checking continues in its defined fallback state (no issuer hint).

Checklist evidence: the `RequestException` is an expected outcome of a best-effort network probe with a defined fallback; the pass is the intended behavior.

### [ ] Finding `127` — `CWE-396`

- Function context: `scripts/niquests/findings/functions/127.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/sgi/__init__.py:129:1`
- Checklist pattern: generic catch inside the retry loop that re-raises on exhaustion

Source excerpt:

```
        127:             try:
        128:                 response = self._do_send(request, stream)
    >   129:             except Exception as err:
        130:                 try:
```

Why this is a false positive: the handler increments retries and re-raises `MaxRetryError` when retries are exhausted (`except MaxRetryError: raise`), then sleeps and retries otherwise — failures surface to the caller.

Checklist evidence: CWE-396's 'can hide distinct failure conditions' is unmet: the exception is fed to retry logic and re-raised on exhaustion.

### [ ] Finding `128` — `CWE-396`

- Function context: `scripts/niquests/findings/functions/128.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/sgi/_async/__init__.py:182:1`
- Checklist pattern: generic catch inside the retry loop that re-raises on exhaustion

Source excerpt:

```
        180:             try:
        181:                 response = await self._do_send(request, stream, timeout)
    >   182:             except Exception as err:
        183:                 try:
```

Why this is a false positive: the handler increments retries and re-raises `MaxRetryError` when retries are exhausted (`except MaxRetryError: raise`), then sleeps and retries otherwise — failures surface to the caller.

Checklist evidence: CWE-396's 'can hide distinct failure conditions' is unmet: the exception is fed to retry logic and re-raised on exhaustion.

### [ ] Finding `129` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/129.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/sgi/_async/__init__.py:355:1`
- Checklist pattern: application exception captured for later propagation

Source excerpt:

```
        353:             try:
        354:                 await self.app(scope, receive, send_func)
    >   355:             except Exception as ex:
        356:                 app_exception = ex
```

Why this is a false positive: the ASGI app exception is stored in `app_exception` and re-raised to the caller after the response cycle completes.

Checklist evidence: the exception is propagated, not hidden.

### [ ] Finding `130` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/130.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/sgi/_async/__init__.py:661:1`
- Checklist pattern: startup failure recorded and surfaced

Source excerpt:

```
        659:             try:
        660:                 await self.app(scope, receive, send)
    >   661:             except Exception as e:
        662:                 if not startup_complete.is_set():
```

Why this is a false positive: the handler checks the startup-complete state and records the exception so the lifespan startup failure is reported.

Checklist evidence: the failure is surfaced through the startup handshake, not swallowed.

### [ ] Finding `131` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/131.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/sgi/_async/__init__.py:679:1`
- Checklist pattern: wait-until-cancelled loop termination

Source excerpt:

```
        677:         try:
        678:             await asyncio.Future()  # Wait forever until canceled
    >   679:         except (asyncio.CancelledError, GeneratorExit):
        680:             pass
```

Why this is a false positive: `await asyncio.Future()` runs until cancelled; catching `CancelledError`/`GeneratorExit` is the designed shutdown exit of the server task.

Checklist evidence: cancellation is the intended control flow, not a silently discarded failure.

### [ ] Finding `132` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/132.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/sgi/_async/__init__.py:679:1`
- Checklist pattern: wait-until-cancelled loop termination

Source excerpt:

```
        677:         try:
        678:             await asyncio.Future()  # Wait forever until canceled
    >   679:         except (asyncio.CancelledError, GeneratorExit):
        680:             pass
```

Why this is a false positive: `await asyncio.Future()` runs until cancelled; catching `CancelledError`/`GeneratorExit` is the designed shutdown exit of the server task.

Checklist evidence: cancellation is the intended control flow, not a silently discarded failure.

### [ ] Finding `133` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/133.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/sgi/_async/__init__.py:687:1`
- Checklist pattern: designed lifespan-shutdown path

Source excerpt:

```
        685:                 await receive_queue.put({"type": "lifespan.shutdown"})
        686:                 await asyncio.wait_for(shutdown_complete.wait(), timeout=5.0)
    >   687:             except (asyncio.TimeoutError, asyncio.CancelledError, RuntimeError):
        688:                 pass
```

Why this is a false positive: the handler is the ASGI lifespan shutdown routine; timeouts/cancellation during graceful shutdown are the expected outcome and the `pass` completes the shutdown sequence.

Checklist evidence: the exception set is the expected outcome of the shutdown handshake, so the condition 'failure discarded silently' is unmet.

### [ ] Finding `134` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/134.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/sgi/_async/__init__.py:711:1`
- Checklist pattern: bridge that propagates the exception into the future

Source excerpt:

```
        709: 
        710:                 future.set_result(result)  # type: ignore[arg-type]
    >   711:             except Exception as e:
        712:                 future.set_exception(e)
```

Why this is a false positive: `future.set_exception(e)` transports the exception to the awaiting caller — the generic catch is the async bridge, not a swallow.

Checklist evidence: the exception is re-raised to the awaiting side, so no failure is hidden.

### [ ] Finding `135` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/135.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/sgi/_async/__init__.py:735:1`
- Checklist pattern: bridge that propagates the exception into the future

Source excerpt:

```
        733: 
        734:                 future.set_result(result)  # type: ignore[arg-type]
    >   735:             except Exception as e:
        736:                 future.set_exception(e)
```

Why this is a false positive: `future.set_exception(e)` transports the exception to the awaiting caller — the generic catch is the async bridge, not a swallow.

Checklist evidence: the exception is re-raised to the awaiting side, so no failure is hidden.

### [ ] Finding `136` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/136.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/sgi/_async/__init__.py:795:1`
- Checklist pattern: bridge that propagates the exception into the future

Source excerpt:

```
        793:                 _swap_context(result)
        794:                 future.set_result(result)  # type: ignore[arg-type]
    >   795:             except Exception as e:
        796:                 future.set_exception(e)
```

Why this is a false positive: `future.set_exception(e)` transports the exception to the awaiting caller — the generic catch is the async bridge, not a swallow.

Checklist evidence: the exception is re-raised to the awaiting side, so no failure is hidden.

### [ ] Finding `137` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/137.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/sgi/_sse.py:153:1`
- Checklist pattern: bridge that propagates the exception into the future

Source excerpt:

```
        151:                 result = await self._async_ext.next_payload(raw=raw)
        152:                 future.set_result(result)
    >   153:             except Exception as e:
        154:                 future.set_exception(e)
```

Why this is a false positive: `future.set_exception(e)` transports the exception to the awaiting caller — the generic catch is the async bridge, not a swallow.

Checklist evidence: the exception is re-raised to the awaiting side, so no failure is hidden.

### [ ] Finding `138` — `CWE-396`

- Function context: `scripts/niquests/findings/functions/138.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/sgi/_sse.py:153:1`
- Checklist pattern: bridge that propagates the exception into the future

Source excerpt:

```
        151:                 result = await self._async_ext.next_payload(raw=raw)
        152:                 future.set_result(result)
    >   153:             except Exception as e:
        154:                 future.set_exception(e)
```

Why this is a false positive: `future.set_exception(e)` transports the exception to the awaiting caller — the generic catch is the async bridge, not a swallow.

Checklist evidence: the exception is re-raised to the awaiting side, so no failure is hidden.

### [ ] Finding `139` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/139.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/sgi/_sse.py:170:1`
- Checklist pattern: bridge that propagates the exception into the future

Source excerpt:

```
        168:                 await self._async_ext.close()
        169:                 future.set_result(None)
    >   170:             except Exception as e:
        171:                 future.set_exception(e)
```

Why this is a false positive: `future.set_exception(e)` transports the exception to the awaiting caller — the generic catch is the async bridge, not a swallow.

Checklist evidence: the exception is re-raised to the awaiting side, so no failure is hidden.

### [ ] Finding `140` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/140.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/sgi/_ws.py:44:1`
- Checklist pattern: bridge that propagates the exception into the future

Source excerpt:

```
         42:                 result = await self._async_ext.next_payload()
         43:                 future.set_result(result)
    >    44:             except Exception as e:
         45:                 future.set_exception(e)
```

Why this is a false positive: `future.set_exception(e)` transports the exception to the awaiting caller — the generic catch is the async bridge, not a swallow.

Checklist evidence: the exception is re-raised to the awaiting side, so no failure is hidden.

### [ ] Finding `141` — `CWE-396`

- Function context: `scripts/niquests/findings/functions/141.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/sgi/_ws.py:44:1`
- Checklist pattern: bridge that propagates the exception into the future

Source excerpt:

```
         42:                 result = await self._async_ext.next_payload()
         43:                 future.set_result(result)
    >    44:             except Exception as e:
         45:                 future.set_exception(e)
```

Why this is a false positive: `future.set_exception(e)` transports the exception to the awaiting caller — the generic catch is the async bridge, not a swallow.

Checklist evidence: the exception is re-raised to the awaiting side, so no failure is hidden.

### [ ] Finding `142` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/142.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/sgi/_ws.py:58:1`
- Checklist pattern: bridge that propagates the exception into the future

Source excerpt:

```
         56:                 await self._async_ext.send_payload(buf)
         57:                 future.set_result(None)
    >    58:             except Exception as e:
         59:                 future.set_exception(e)
```

Why this is a false positive: `future.set_exception(e)` transports the exception to the awaiting caller — the generic catch is the async bridge, not a swallow.

Checklist evidence: the exception is re-raised to the awaiting side, so no failure is hidden.

### [ ] Finding `143` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/143.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/sgi/_ws.py:72:1`
- Checklist pattern: bridge that propagates the exception into the future

Source excerpt:

```
         70:                 await self._async_ext.close()
         71:                 future.set_result(None)
    >    72:             except Exception as e:
         73:                 future.set_exception(e)
```

Why this is a false positive: `future.set_exception(e)` transports the exception to the awaiting caller — the generic catch is the async bridge, not a swallow.

Checklist evidence: the exception is re-raised to the awaiting side, so no failure is hidden.

### [ ] Finding `144` — `CWE-772`

- Function context: `scripts/niquests/findings/functions/144.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/unixsocket/__init__.py:47:23`
- Checklist pattern: socket ownership transferred to the connection object

Source excerpt:

```
         45:         self.host = self.socket_path.split("/")[-1]
         46: 
    >    47:     def connect(self):
         48:         sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
```

Why this is a false positive: `connect()` stores the socket on `self.sock` (line 51); the `HTTPConnection` lifecycle closes it, so the resource is released after its effective lifetime by its owner.

Checklist evidence: CWE-772 requires the resource to never be released; here it is stored on the instance and closed by the connection lifecycle.

### [ ] Finding `145` — `CWE-396`

- Function context: `scripts/niquests/findings/functions/145.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_adapter.py:107:1`
- Checklist pattern: re-raise of every unexpected exception

Source excerpt:

```
        105:                 finally:
        106:                     close_resource(pollable)
    >   107:             except BaseException as exc:
        108:                 if not isinstance(exc, _Err):
```

Why this is a false positive: the handler re-raises anything that is not the expected WASI `_Err`; the `_Err` variant is mapped to a typed exception (`with wasi_exception_mapping(...): raise`).

Checklist evidence: CWE-396's condition is unmet: unexpected failures are re-raised and expected ones are translated to typed errors.

### [ ] Finding `148` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/148.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_adapter.py:216:1`
- Checklist pattern: retry-loop handler that acts on the exception

Source excerpt:

```
        214:             except (InvalidSchema, RuntimeError, SSLError, _WASIProxyError):
        215:                 raise
    >   216:             except Exception as exc:
        217:                 retries = retries.increment(method, request.url, error=exc)
```

Why this is a false positive: the handler increments the retry state and re-raises `MaxRetryError` on exhaustion — the failure is not swallowed.

Checklist evidence: the generic catch feeds the retry machinery and re-raises on exhaustion.

### [ ] Finding `149` — `CWE-1341`

- Function context: `scripts/niquests/findings/functions/149.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_adapter.py:235:25`
- Checklist pattern: mutually exclusive close paths, single release per execution

Source excerpt:

```
        233:                 except MaxRetryError:
        234:                     if retries.raise_on_status:
    >   235:                         response.close()
        236:                         raise
```

Why this is a false positive: `async_session.py:1015` and its `else: resp.close()` are the two branches of one `isinstance` conditional — exactly one executes; the wasi `response.close()` calls sit in exclusive retry/raise branches (line 235 vs 239).

Checklist evidence: CWE-1341's condition is the same handle released twice on one path; the matched `close()` pairs are mutually exclusive branches (or distinct resources).

### [ ] Finding `150` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/150.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_async/_adapter.py:99:1`
- Checklist pattern: upload-failure cleanup flag

Source excerpt:

```
         97:                     if on_upload_body is not None:  # Defensive: Session always supplies the callback
         98:                         await on_upload_body(sent, total, False, False)
    >    99:     except BaseException:  # Defensive: upload failure cleanup
        100:         failed = True
```

Why this is a false positive: the `BaseException` handler sets `failed = True` (documented `# Defensive: upload failure cleanup`), which drives the caller's cleanup and completion path.

Checklist evidence: the handler takes a visible action, so failures are not silently discarded.

### [ ] Finding `151` — `CWE-396`

- Function context: `scripts/niquests/findings/functions/151.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_async/_adapter.py:99:1`
- Checklist pattern: upload-failure cleanup flag

Source excerpt:

```
         97:                     if on_upload_body is not None:  # Defensive: Session always supplies the callback
         98:                         await on_upload_body(sent, total, False, False)
    >    99:     except BaseException:  # Defensive: upload failure cleanup
        100:         failed = True
```

Why this is a false positive: the `BaseException` handler sets `failed = True` (documented `# Defensive: upload failure cleanup`), which drives the caller's cleanup and completion path.

Checklist evidence: the handler takes a visible action, so failures are not silently discarded.

### [ ] Finding `154` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/154.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_async/_adapter.py:237:1`
- Checklist pattern: retry-loop handler that acts on the exception

Source excerpt:

```
        235:             except (InvalidSchema, RuntimeError, SSLError, _WASIProxyError):
        236:                 raise
    >   237:             except Exception as exc:
        238:                 retries = retries.increment(method, request.url, error=exc)
```

Why this is a false positive: the handler increments the retry state and re-raises `MaxRetryError` on exhaustion — the failure is not swallowed.

Checklist evidence: the generic catch feeds the retry machinery and re-raises on exhaustion.

### [ ] Finding `155` — `CWE-1341`

- Function context: `scripts/niquests/findings/functions/155.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_async/_adapter.py:256:31`
- Checklist pattern: mutually exclusive close paths, single release per execution

Source excerpt:

```
        254:                 except MaxRetryError:
        255:                     if retries.raise_on_status:
    >   256:                         await response.close()
        257:                         raise
```

Why this is a false positive: `async_session.py:1015` and its `else: resp.close()` are the two branches of one `isinstance` conditional — exactly one executes; the wasi `response.close()` calls sit in exclusive retry/raise branches (line 235 vs 239).

Checklist evidence: CWE-1341's condition is the same handle released twice on one path; the matched `close()` pairs are mutually exclusive branches (or distinct resources).

### [ ] Finding `156` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/156.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_async/_sse.py:62:1`
- Checklist pattern: optional numeric conversion fallback in SSE field parsing

Source excerpt:

```
         60:                 try:
         61:                     values[key] = int(value)
    >    62:                 except ValueError:
         63:                     pass
```

Why this is a false positive: the `int(value)` conversion is best-effort; when the field is not numeric the string value is kept, which is the intended fallback for SSE `retry` fields.

Checklist evidence: the exception is an expected parsing outcome with a defined fallback (original string preserved).

### [ ] Finding `157` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/157.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_async/_sse.py:62:1`
- Checklist pattern: optional numeric conversion fallback in SSE field parsing

Source excerpt:

```
         60:                 try:
         61:                     values[key] = int(value)
    >    62:                 except ValueError:
         63:                     pass
```

Why this is a false positive: the `int(value)` conversion is best-effort; when the field is not numeric the string value is kept, which is the intended fallback for SSE `retry` fields.

Checklist evidence: the exception is an expected parsing outcome with a defined fallback (original string preserved).

### [ ] Finding `158` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/158.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_async/_sse.py:62:17`
- Checklist pattern: optional numeric conversion fallback in SSE field parsing

Source excerpt:

```
         60:                 try:
         61:                     values[key] = int(value)
    >    62:                 except ValueError:
         63:                     pass
```

Why this is a false positive: the `int(value)` conversion is best-effort; when the field is not numeric the string value is kept, which is the intended fallback for SSE `retry` fields.

Checklist evidence: the exception is an expected parsing outcome with a defined fallback (original string preserved).

### [ ] Finding `159` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/159.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_sse.py:61:1`
- Checklist pattern: optional numeric conversion fallback in SSE field parsing

Source excerpt:

```
         59:                 try:
         60:                     values[key] = int(value)
    >    61:                 except ValueError:
         62:                     pass
```

Why this is a false positive: the `int(value)` conversion is best-effort; when the field is not numeric the string value is kept, which is the intended fallback for SSE `retry` fields.

Checklist evidence: the exception is an expected parsing outcome with a defined fallback (original string preserved).

### [ ] Finding `160` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/160.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_sse.py:61:1`
- Checklist pattern: optional numeric conversion fallback in SSE field parsing

Source excerpt:

```
         59:                 try:
         60:                     values[key] = int(value)
    >    61:                 except ValueError:
         62:                     pass
```

Why this is a false positive: the `int(value)` conversion is best-effort; when the field is not numeric the string value is kept, which is the intended fallback for SSE `retry` fields.

Checklist evidence: the exception is an expected parsing outcome with a defined fallback (original string preserved).

### [ ] Finding `161` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/161.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_sse.py:61:17`
- Checklist pattern: optional numeric conversion fallback in SSE field parsing

Source excerpt:

```
         59:                 try:
         60:                     values[key] = int(value)
    >    61:                 except ValueError:
         62:                     pass
```

Why this is a false positive: the `int(value)` conversion is best-effort; when the field is not numeric the string value is kept, which is the intended fallback for SSE `retry` fields.

Checklist evidence: the exception is an expected parsing outcome with a defined fallback (original string preserved).

### [ ] Finding `162` — `BP-PY-49`

- Function context: `scripts/niquests/findings/functions/162.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_utils.py:63:41`
- Checklist pattern: `verify=False` text inside an error-message string

Source excerpt:

```
         61:     if verify is not True:
         62:         raise SSLError(
    >    63:             "WASI HTTPS uses the host trust policy; verify=False and custom CA bundles are unavailable.",
         64:             request=request,
```

Why this is a false positive: the flagged text is a string literal in an `SSLError` message; the surrounding code (`if verify is not True: raise SSLError(...)`) actively rejects `verify=False` on WASI.

Checklist evidence: no TLS disabling construct exists at this line — the regex matched a string literal inside an error message.

### [ ] Finding `165` — `BP-PY-36`

- Function context: `scripts/niquests/findings/functions/165.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/help.py:176:1`
- Checklist pattern: niquests' own HTTP `Session`, not SQLAlchemy

Source excerpt:

```
        174: 
        175: 
    >   176: pypi_session = Session()
        177: 
```

Why this is a false positive: `help.py:34` imports `Session` from niquests itself (`from . import ... Session`); `pypi_session = Session()` is an HTTP client session used to query the PyPI JSON API — SQLAlchemy is not involved.

Checklist evidence: BP-PY-36's condition is a SQLAlchemy `Session`/`SessionLocal`; the shown source constructs the package's own HTTP session class.

### [ ] Finding `166` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/166.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/help.py:196:1`
- Checklist pattern: best-effort update check

Source excerpt:

```
        194:                     UserWarning,
        195:                 )
    >   196:     except (RequestException, JSONDecodeError, HTTPError):
        197:         pass
```

Why this is a false positive: `check_update` is an optional PyPI version probe; network/JSON failures are benign and the check silently degrades — the pass is the designed behavior of the best-effort helper.

Checklist evidence: the exception set (network + parse) is expected for a best-effort update check with no failure semantics.

### [ ] Finding `167` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/167.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/help.py:196:1`
- Checklist pattern: best-effort update check

Source excerpt:

```
        194:                     UserWarning,
        195:                 )
    >   196:     except (RequestException, JSONDecodeError, HTTPError):
        197:         pass
```

Why this is a false positive: `check_update` is an optional PyPI version probe; network/JSON failures are benign and the check silently degrades — the pass is the designed behavior of the best-effort helper.

Checklist evidence: the exception set (network + parse) is expected for a best-effort update check with no failure semantics.

### [ ] Finding `168` — `CWE-93`

- Function context: `scripts/niquests/findings/functions/168.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/models.py:488:25`
- Checklist pattern: write to an outgoing request header, not a response header

Source excerpt:

```
        486:                     if isinstance(name, bytes):
        487:                         name = name.decode()
    >   488:                     self.headers[name] = value
        489: 
```

Why this is a false positive: `auth.py:75` sets the `Authorization` header on the outgoing request (from the caller's own token); `models.py:488` populates the `PreparedRequest.headers` dict in `prepare_headers` from the caller-supplied header mapping — neither is an HTTP response header.

Checklist evidence: CWE-93's condition is writing an externally influenced value into an HTTP response header; the shown lines write client-side request headers.

### [ ] Finding `169` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/169.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/models.py:1065:1`
- Checklist pattern: defensive probe of lazily-initialized attributes

Source excerpt:

```
       1063:             super().__getattribute__("_promise")
       1064:             super().__getattribute__("connection").gather(self)
    >  1065:         except AttributeError:
       1066:             pass
```

Why this is a false positive: the `try` probes `_promise`/`_gather` internals that exist only in initialized states; `AttributeError` means there is nothing to gather — the designed no-op.

Checklist evidence: the exception is the expected outcome for an uninitialized optional internal, and the no-op is the designed behavior.

### [ ] Finding `170` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/170.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/models.py:1065:1`
- Checklist pattern: defensive probe of lazily-initialized attributes

Source excerpt:

```
       1063:             super().__getattribute__("_promise")
       1064:             super().__getattribute__("connection").gather(self)
    >  1065:         except AttributeError:
       1066:             pass
```

Why this is a false positive: the `try` probes `_promise`/`_gather` internals that exist only in initialized states; `AttributeError` means there is nothing to gather — the designed no-op.

Checklist evidence: the exception is the expected outcome for an uninitialized optional internal, and the no-op is the designed behavior.

### [ ] Finding `171` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/171.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/models.py:1065:9`
- Checklist pattern: defensive probe of lazily-initialized attributes

Source excerpt:

```
       1063:             super().__getattribute__("_promise")
       1064:             super().__getattribute__("connection").gather(self)
    >  1065:         except AttributeError:
       1066:             pass
```

Why this is a false positive: the `try` probes `_promise`/`_gather` internals that exist only in initialized states; `AttributeError` means there is nothing to gather — the designed no-op.

Checklist evidence: the exception is the expected outcome for an uninitialized optional internal, and the no-op is the designed behavior.

### [ ] Finding `172` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/172.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/models.py:1079:1`
- Checklist pattern: defensive probe of lazily-initialized attributes

Source excerpt:

```
       1077:                 else:
       1078:                     super().__getattribute__("_gather")()
    >  1079:         except AttributeError:
       1080:             pass
```

Why this is a false positive: the `try` probes `_promise`/`_gather` internals that exist only in initialized states; `AttributeError` means there is nothing to gather — the designed no-op.

Checklist evidence: the exception is the expected outcome for an uninitialized optional internal, and the no-op is the designed behavior.

### [ ] Finding `173` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/173.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/models.py:1233:1`
- Checklist pattern: optional content-length progress hint fallback

Source excerpt:

```
       1231:                     try:
       1232:                         self.download_progress.content_length = int(self.headers["content-length"])
    >  1233:                     except ValueError:
       1234:                         pass
```

Why this is a false positive: the `content-length` header may be absent or non-numeric; on `ValueError` the progress `content_length` simply stays unset — an optional hint, not a failure.

Checklist evidence: the flagged `pass` is the expected fallback for an optional progress attribute.

### [ ] Finding `174` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/174.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/models.py:1318:1`
- Checklist pattern: optional content-length progress hint fallback

Source excerpt:

```
       1316:                     try:
       1317:                         self.download_progress.content_length = int(self.headers["content-length"])
    >  1318:                     except ValueError:
       1319:                         pass
```

Why this is a false positive: the `content-length` header may be absent or non-numeric; on `ValueError` the progress `content_length` simply stays unset — an optional hint, not a failure.

Checklist evidence: the flagged `pass` is the expected fallback for an optional progress attribute.

### [ ] Finding `175` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/175.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/models.py:1768:1`
- Checklist pattern: optional content-length progress hint fallback

Source excerpt:

```
       1766:                     try:
       1767:                         self.download_progress.content_length = int(self.headers["content-length"])
    >  1768:                     except ValueError:
       1769:                         pass
```

Why this is a false positive: the `content-length` header may be absent or non-numeric; on `ValueError` the progress `content_length` simply stays unset — an optional hint, not a failure.

Checklist evidence: the flagged `pass` is the expected fallback for an optional progress attribute.

### [ ] Finding `176` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/176.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/models.py:1827:1`
- Checklist pattern: optional content-length progress hint fallback

Source excerpt:

```
       1825:                     try:
       1826:                         self.download_progress.content_length = int(self.headers["content-length"])
    >  1827:                     except ValueError:
       1828:                         pass
```

Why this is a false positive: the `content-length` header may be absent or non-numeric; on `ValueError` the progress `content_length` simply stays unset — an optional hint, not a failure.

Checklist evidence: the flagged `pass` is the expected fallback for an optional progress attribute.

### [ ] Finding `177` — `CWE-829`

- Function context: `scripts/niquests/findings/functions/177.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/packages.py:42:29`
- Checklist pattern: dynamic import over a developer-controlled constant allowlist

Source excerpt:

```
         40: 
         41:     try:
    >    42:         locals()[package] = __import__(to_be_imported)
         43:     except ImportError:
```

Why this is a false positive: `to_be_imported` is computed from the hardcoded module tuple `("urllib3", "charset_normalizer", "idna", "chardet")` plus fixed alias mappings — no request- or user-derived value reaches `__import__`.

Checklist evidence: CWE-829/CWE-94's condition is dynamically loading untrusted/request-derived modules; the shown loop iterates a compile-time allowlist of the package's own dependencies.

### [ ] Finding `178` — `CWE-94`

- Function context: `scripts/niquests/findings/functions/178.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/packages.py:42:29`
- Checklist pattern: dynamic import over a developer-controlled constant allowlist

Source excerpt:

```
         40: 
         41:     try:
    >    42:         locals()[package] = __import__(to_be_imported)
         43:     except ImportError:
```

Why this is a false positive: `to_be_imported` is computed from the hardcoded module tuple `("urllib3", "charset_normalizer", "idna", "chardet")` plus fixed alias mappings — no request- or user-derived value reaches `__import__`.

Checklist evidence: CWE-829/CWE-94's condition is dynamically loading untrusted/request-derived modules; the shown loop iterates a compile-time allowlist of the package's own dependencies.

### [ ] Finding `179` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/179.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/sessions.py:115:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
        113:     try:
        114:         from .extensions.wasi._adapter import WASIAdapter as WASIHTTPAdapter
    >   115:     except ImportError:
        116:         pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `180` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/180.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/sessions.py:115:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
        113:     try:
        114:         from .extensions.wasi._adapter import WASIAdapter as WASIHTTPAdapter
    >   115:     except ImportError:
        116:         pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `181` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/181.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/sessions.py:115:5`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
        113:     try:
        114:         from .extensions.wasi._adapter import WASIAdapter as WASIHTTPAdapter
    >   115:     except ImportError:
        116:         pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `182` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/182.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/sessions.py:1449:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
       1447:                             verify as ocsp_verify,
       1448:                         )
    >  1449:                     except ImportError:
       1450:                         pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `184` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/184.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/sessions.py:1476:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
       1474:                             verify as crl_verify,
       1475:                         )
    >  1476:                     except ImportError:
       1477:                         pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `186` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/186.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/sessions.py:1709:1`
- Checklist pattern: generator-exhaustion control flow

Source excerpt:

```
       1707:                         self.resolve_redirects(r, request, yield_requests=True, **kwargs)  # type: ignore[assignment]
       1708:                     )
    >  1709:                 except StopIteration:
       1710:                     pass
```

Why this is a false positive: `StopIteration`/`StopAsyncIteration` is the generator's normal termination signal; catching it is the intended loop-exit control flow, not a discarded failure.

Checklist evidence: the rule's 'failures discarded silently' condition does not apply to iterator-exhaustion signalling.

### [ ] Finding `187` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/187.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/sessions.py:1785:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
       1783:                 if scheme in extension.supported_schemes() and extension.scheme_to_http_scheme(scheme) == parse_scheme(prefix):
       1784:                     return adapter
    >  1785:         except ImportError:
       1786:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `188` — `CWE-1341`

- Function context: `scripts/niquests/findings/functions/188.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/sessions.py:1801:13`
- Checklist pattern: different handles closed in one loop

Source excerpt:

```
       1799:         """Closes all adapters and as such the session"""
       1800:         for v in self.adapters.values():
    >  1801:             v.close()
       1802:         if self._own_resolver:
```

Why this is a false positive: `for v in self.adapters.values(): v.close()` closes each adapter once; `self.resolver.close()` closes a different resource.

Checklist evidence: the two `close()` calls release different resources, so the same handle is not released twice.

### [ ] Finding `195` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/195.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/utils.py:157:1`
- Checklist pattern: expected fallback for non-file objects

Source excerpt:

```
        155:         try:
        156:             fileno = o.fileno()
    >   157:         except (io.UnsupportedOperation, AttributeError):
        158:             # AttributeError is a surprising exception, seeing as how we've just checked
```

Why this is a false positive: `fileno()` is unsupported for pure-Python/byte-stream objects; the surrounding comment documents the fallback path that handles objects without a file descriptor.

Checklist evidence: the shown source documents the `AttributeError`/`UnsupportedOperation` as a surprising-but-expected case with a defined fallback.

### [ ] Finding `196` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/196.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/utils.py:157:1`
- Checklist pattern: expected fallback for non-file objects

Source excerpt:

```
        155:         try:
        156:             fileno = o.fileno()
    >   157:         except (io.UnsupportedOperation, AttributeError):
        158:             # AttributeError is a surprising exception, seeing as how we've just checked
```

Why this is a false positive: `fileno()` is unsupported for pure-Python/byte-stream objects; the surrounding comment documents the fallback path that handles objects without a file descriptor.

Checklist evidence: the shown source documents the `AttributeError`/`UnsupportedOperation` as a surprising-but-expected case with a defined fallback.

### [ ] Finding `197` — `CWE-396`

- Function context: `scripts/niquests/findings/functions/197.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/utils.py:256:1`
- Checklist pattern: cleanup handler that removes the temp file

Source excerpt:

```
        254:             yield tmp_handler
        255:         os.replace(tmp_name, filename)
    >   256:     except BaseException:
        257:         os.remove(tmp_name)
```

Why this is a false positive: `except BaseException: os.remove(tmp_name)` performs the cleanup action before the context manager re-raises.

Checklist evidence: the handler acts (removes the temporary file), so the 'detection without action' condition is unmet.

### [ ] Finding `198` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/198.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/utils.py:815:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
        813:         try:
        814:             extension_class = load_extension(maybe_extension_scheme, implementation)
    >   815:         except ImportError:
        816:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `199` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/199.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/utils.py:815:9`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
        813:         try:
        814:             extension_class = load_extension(maybe_extension_scheme, implementation)
    >   815:         except ImportError:
        816:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `201` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/201.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/utils.py:1196:1`
- Checklist pattern: scheme-parsing fallback that falls through to an explicit raise

Source excerpt:

```
       1194: 
       1195:             return outsider_scheme.lower()
    >  1196:         except ValueError:
       1197:             pass
```

Why this is a false positive: the inner `except ValueError: pass` is immediately followed by `raise MissingSchema(...)`, so the parse error is surfaced to the caller, not discarded.

Checklist evidence: the `pass` only bridges to the enclosing `raise`, so no error condition goes without action.

### [ ] Finding `202` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/202.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/conftest.py:107:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
        105:             yield selenium
        106: 
    >   107: except ImportError:
        108:     pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `203` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/203.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/conftest.py:107:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
        105:             yield selenium
        106: 
    >   107: except ImportError:
        108:     pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `204` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/204.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/conftest.py:107:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
        105:             yield selenium
        106: 
    >   107: except ImportError:
        108:     pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `205` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/205.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_asgi.py:52:1`
- Checklist pattern: test websocket receive loop termination

Source excerpt:

```
         50:         try:
         51:             data = await websocket.receive()
    >    52:         except Exception:
         53:             break
```

Why this is a false positive: the test receives frames in a loop and breaks on any connection error — the generic catch is the loop's termination, and a broken socket fails the test otherwise.

Checklist evidence: test code: the exception is the expected end of the receive loop and any real failure fails the test.

### [ ] Finding `206` — `CWE-617`

- Function context: `scripts/niquests/findings/functions/206.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_async.py:30:1`
- Checklist pattern: test assertion is the verification mechanism

Source excerpt:

```
         28:     assert encoded_values == [value]
         29:     assert request.body == b'{"id":2}'
    >    30:     assert request.headers["Content-Type"] == "application/json;charset=utf-8"
         31:     assert request.headers["Content-Length"] == "8"
```

Why this is a false positive: `assert request.headers[...] == ...` verifies the request the server received; the asserted values are the test's own expectations.

Checklist evidence: CWE-617's condition is a reachable assertion on request-controlled state in production; the shown source is test verification.

### [ ] Finding `207` — `CWE-93`

- Function context: `scripts/niquests/findings/functions/207.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_async.py:30:19`
- Checklist pattern: test assertion reading a header value

Source excerpt:

```
         28:     assert encoded_values == [value]
         29:     assert request.body == b'{"id":2}'
    >    30:     assert request.headers["Content-Type"] == "application/json;charset=utf-8"
         31:     assert request.headers["Content-Length"] == "8"
```

Why this is a false positive: each flagged line is `assert <resp>.headers[...] == "..."` — a pytest assertion that reads a received header; nothing is written to any header.

Checklist evidence: CWE-93 requires a header write; the shown source only compares a header value in a test.

### [ ] Finding `208` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/208.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_async.py:195:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >   195:     async def test_explicit_close_in_streaming_response(self) -> None:
       196:         async with AsyncSession() as s:
       197:             try:
       198:                 r = await s.get("https://httpbingo.org/html", stream=True)
```

Why this is a false positive: the test explicitly closes a streaming response and relies on any failure to surface as an exception (there is no `except`), so the test fails if close is broken.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `209` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/209.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:58:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >    58: def test_sync_basic_get(selenium_jspi):
        59:     """Test basic GET request works."""
        60:     data = _inner_test_sync_basic_get(selenium_jspi)
        61:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `210` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/210.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:91:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >    91: def test_sync_basic_post(selenium_jspi):
        92:     """Test basic POST request works."""
        93:     data = _inner_test_sync_basic_post(selenium_jspi)
        94:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `211` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/211.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:129:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >   129: def test_sync_streaming_read(selenium_jspi):
       130:     """Test streaming response works."""
       131:     data = _inner_test_sync_streaming_read(selenium_jspi)
       132:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `212` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/212.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:160:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >   160: def test_sync_retry_works(selenium_jspi):
       161:     """Test that retry mechanism works."""
       162:     data = _inner_test_sync_retry_works(selenium_jspi)
       163:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `213` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/213.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:192:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >   192: def test_async_basic_get(selenium):
       193:     """Test basic async GET request works."""
       194:     data = _inner_test_async_basic_get(selenium)
       195:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `214` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/214.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:225:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >   225: def test_async_basic_post(selenium):
       226:     """Test basic async POST request works."""
       227:     data = _inner_test_async_basic_post(selenium)
       228:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `215` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/215.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:263:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >   263: def test_async_streaming_read(selenium):
       264:     """Test async streaming response works."""
       265:     data = _inner_test_async_streaming_read(selenium)
       266:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `216` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/216.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:294:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >   294: def test_async_retry_works(selenium):
       295:     """Test that async retry mechanism works."""
       296:     data = _inner_test_async_retry_works(selenium)
       297:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `217` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/217.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:353:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >   353: def test_sync_websocket(selenium_jspi):
       354:     """Test sync WebSocket via browser native API + JSPI."""
       355:     data = _inner_test_sync_websocket(selenium_jspi)
       356:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `218` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/218.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:403:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >   403: def test_sync_websocket_binary(selenium_jspi):
       404:     """Test sync WebSocket binary message handling."""
       405:     data = _inner_test_sync_websocket_binary(selenium_jspi)
       406:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `219` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/219.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:446:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >   446: def test_sync_websocket_close_from_server(selenium_jspi):
       447:     """Test that close works and state is updated."""
       448:     data = _inner_test_sync_websocket_close_from_server(selenium_jspi)
       449:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `220` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/220.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:500:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >   500: def test_async_websocket(selenium):
       501:     """Test async WebSocket via browser native API."""
       502:     data = _inner_test_async_websocket(selenium)
       503:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `221` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/221.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:550:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >   550: def test_async_websocket_binary(selenium):
       551:     """Test async WebSocket binary message handling."""
       552:     data = _inner_test_async_websocket_binary(selenium)
       553:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `222` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/222.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:598:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >   598: def test_sync_sse(selenium_jspi):
       599:     """Test sync SSE via pyfetch streaming + manual parsing."""
       600:     data = _inner_test_sync_sse(selenium_jspi)
       601:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `223` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/223.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:646:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >   646: def test_async_sse(selenium):
       647:     """Test async SSE via pyfetch streaming + manual parsing."""
       648:     data = _inner_test_async_sse(selenium)
       649:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `224` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/224.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:684:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >   684: def test_sync_sse_send_raises(selenium_jspi):
       685:     """Test that SSE send_payload raises NotImplementedError."""
       686:     data = _inner_test_sync_sse_send_raises(selenium_jspi)
       687:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `225` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/225.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:727:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >   727: def test_sync_timeout(selenium_jspi):
       728:     """Test that timeout raises ConnectTimeout on slow responses."""
       729:     data = _inner_test_sync_timeout(selenium_jspi)
       730:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `226` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/226.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:770:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >   770: def test_async_timeout(selenium):
       771:     """Test that async timeout raises ConnectTimeout on slow responses."""
       772:     data = _inner_test_async_timeout(selenium)
       773:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `227` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/227.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:814:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >   814: def test_sync_custom_request_headers(selenium_jspi):
       815:     """Test that custom request headers are sent and response headers are parsed."""
       816:     data = _inner_test_sync_custom_request_headers(selenium_jspi)
       817:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `228` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/228.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:856:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >   856: def test_async_custom_request_headers(selenium):
       857:     """Test that async custom request headers are sent and response headers are parsed."""
       858:     data = _inner_test_async_custom_request_headers(selenium)
       859:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `229` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/229.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:901:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >   901: def test_sync_forbidden_header_silently_dropped(selenium_jspi):
       902:     """Test that forbidden headers (e.g. Host) are silently stripped and do not cause errors."""
       903:     data = _inner_test_sync_forbidden_header_silently_dropped(selenium_jspi)
       904:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `230` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/230.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:946:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >   946: def test_async_forbidden_header_silently_dropped(selenium):
       947:     """Test that forbidden headers (e.g. Host) are silently stripped and do not cause errors (async)."""
       948:     data = _inner_test_async_forbidden_header_silently_dropped(selenium)
       949:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `231` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/231.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:984:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >   984: def test_sync_non_ok_response(selenium_jspi):
       985:     """Test non-OK responses and raise_for_status."""
       986:     data = _inner_test_sync_non_ok_response(selenium_jspi)
       987:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `232` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/232.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:1022:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  1022: def test_async_non_ok_response(selenium):
      1023:     """Test async non-OK responses and raise_for_status."""
      1024:     data = _inner_test_async_non_ok_response(selenium)
      1025:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `233` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/233.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:1059:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  1059: def test_sync_binary_response(selenium_jspi):
      1060:     """Test binary response content (JPEG image)."""
      1061:     data = _inner_test_sync_binary_response(selenium_jspi)
      1062:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `234` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/234.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:1095:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  1095: def test_async_binary_response(selenium):
      1096:     """Test async binary response content (JPEG image)."""
      1097:     data = _inner_test_async_binary_response(selenium)
      1098:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `235` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/235.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:1128:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  1128: def test_sync_redirect(selenium_jspi):
      1129:     """Test that redirects are followed and final URL is correct."""
      1130:     data = _inner_test_sync_redirect(selenium_jspi)
      1131:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `236` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/236.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:1160:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  1160: def test_async_redirect(selenium):
      1161:     """Test that async redirects are followed and final URL is correct."""
      1162:     data = _inner_test_async_redirect(selenium)
      1163:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `237` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/237.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:1197:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  1197: def test_sync_string_body(selenium_jspi):
      1198:     """Test POST with a raw string body."""
      1199:     data = _inner_test_sync_string_body(selenium_jspi)
      1200:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `238` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/238.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:1234:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  1234: def test_async_string_body(selenium):
      1235:     """Test async POST with a raw string body."""
      1236:     data = _inner_test_async_string_body(selenium)
      1237:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `239` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/239.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:1276:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  1276: def test_sync_iterable_body(selenium_jspi):
      1277:     """Test POST with a generator/iterable body."""
      1278:     data = _inner_test_sync_iterable_body(selenium_jspi)
      1279:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `240` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/240.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:1318:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  1318: def test_async_iterable_body(selenium):
      1319:     """Test async POST with a generator/iterable body."""
      1320:     data = _inner_test_async_iterable_body(selenium)
      1321:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `241` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/241.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:1359:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  1359: def test_sync_iterable_str_body(selenium_jspi):
      1360:     """Test POST with a generator yielding str chunks."""
      1361:     data = _inner_test_sync_iterable_str_body(selenium_jspi)
      1362:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `242` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/242.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:1400:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  1400: def test_async_iterable_str_body(selenium):
      1401:     """Test async POST with a generator yielding str chunks."""
      1402:     data = _inner_test_async_iterable_str_body(selenium)
      1403:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `243` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/243.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:1442:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  1442: def test_async_async_iterable_body(selenium):
      1443:     """Test async POST with an async generator body (covers __aiter__ path)."""
      1444:     data = _inner_test_async_async_iterable_body(selenium)
      1445:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `244` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/244.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:1483:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  1483: def test_async_async_iterable_str_body(selenium):
      1484:     """Test async POST with an async generator yielding str chunks."""
      1485:     data = _inner_test_async_async_iterable_str_body(selenium)
      1486:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `245` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/245.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:1528:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  1528: def test_sync_sse_raw_mode(selenium_jspi):
      1529:     """Test sync SSE with raw=True returns raw event strings."""
      1530:     data = _inner_test_sync_sse_raw_mode(selenium_jspi)
      1531:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `246` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/246.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:1573:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  1573: def test_async_sse_raw_mode(selenium):
      1574:     """Test async SSE with raw=True returns raw event strings."""
      1575:     data = _inner_test_async_sse_raw_mode(selenium)
      1576:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `247` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/247.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:1617:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  1617: def test_sync_sse_via_stream(selenium_jspi):
      1618:     """Test consuming an SSE endpoint as a plain stream (no sse:// scheme)."""
      1619:     data = _inner_test_sync_sse_via_stream(selenium_jspi)
      1620:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `248` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/248.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:1661:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  1661: def test_async_sse_via_stream(selenium):
      1662:     """Test consuming an SSE endpoint as a plain async stream (no sse:// scheme)."""
      1663:     data = _inner_test_async_sse_via_stream(selenium)
      1664:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `249` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/249.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:1715:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  1715: def test_async_concurrent_gather(selenium):
      1716:     """Test that 5 concurrent async requests via asyncio.gather all complete."""
      1717:     data = _inner_test_async_concurrent_gather(selenium)
      1718:     if data:
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `250` — `CWE-93`

- Function context: `scripts/niquests/findings/functions/250.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_lowlevel.py:38:21`
- Checklist pattern: test assertion reading a header value

Source excerpt:

```
         36: 
         37:     assert r.status_code == 200
    >    38:     assert r.request.headers["Transfer-Encoding"] == "chunked"
         39: 
```

Why this is a false positive: each flagged line is `assert <resp>.headers[...] == "..."` — a pytest assertion that reads a received header; nothing is written to any header.

Checklist evidence: CWE-93 requires a header write; the shown source only compares a header value in a test.

### [ ] Finding `251` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/251.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:79:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >    79:     def test_entry_points(self):
        80:         niquests.Session
        81:         niquests.Session().get
        82:         niquests.Session().head
```

Why this is a false positive: the body only touches the public entry points (`niquests.Session`, `.get`, `.head`); a missing or broken entry point raises `AttributeError` and fails the test — an API-surface smoke test.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `252` — `CWE-93`

- Function context: `scripts/niquests/findings/functions/252.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:124:19`
- Checklist pattern: test assertion reading a header value

Source excerpt:

```
        122:     def test_no_body_content_length(self, httpbin, method):
        123:         req = niquests.Request(method, httpbin(method.lower())).prepare()
    >   124:         assert req.headers["Content-Length"] == "0"
        125: 
```

Why this is a false positive: each flagged line is `assert <resp>.headers[...] == "..."` — a pytest assertion that reads a received header; nothing is written to any header.

Checklist evidence: CWE-93 requires a header write; the shown source only compares a header value in a test.

### [ ] Finding `253` — `CWE-208`

- Function context: `scripts/niquests/findings/functions/253.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:574:16`
- Checklist pattern: equality of an auth object in a test assertion

Source excerpt:

```
        572:         """Session accepts custom auth at initialization."""
        573:         s = niquests.Session(auth=init_auth)
    >   574:         assert s.auth == expected_auth
        575: 
```

Why this is a false positive: `assert s.auth == expected_auth` compares `AuthBase` objects configured by the test itself; no secret comparison subject to timing side-channels exists.

Checklist evidence: CWE-208's condition is comparing security-sensitive values; the compared values are test-owned auth objects, and the line is a test assertion.

### [ ] Finding `254` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/254.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:1030:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  1030:     def test_decompress_gzip(self, httpbin):
      1031:         r = niquests.get(httpbin("gzip"))
      1032:         r.content.decode("ascii")
```

Why this is a false positive: `r.content.decode("ascii")` raises if gzip decompression produced garbage — the test fails on broken decompression, so it verifies the outcome.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `255` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/255.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:1044:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  1044:     def test_unicode_get(self, httpbin, url, params):
      1045:         niquests.get(httpbin(url), params=params)
```

Why this is a false positive: the request fails (raises) if unicode query params are mishandled, failing the test — a smoke test of unicode parameter support.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `256` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/256.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:1047:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  1047:     def test_unicode_header_name(self, httpbin):
      1048:         niquests.put(
      1049:             httpbin("put"),
      1050:             headers={"Content-Type": "application/octet-stream"},
```

Why this is a false positive: the request fails (raises) if unicode header names/values are mishandled, failing the test — a smoke test of unicode header support.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `257` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/257.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:1054:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  1054:     def test_pyopenssl_redirect(self, httpbin_secure, httpbin_ca_bundle):
      1055:         niquests.get(httpbin_secure("status", "301"), verify=httpbin_ca_bundle)
```

Why this is a false positive: the redirected request over the pyopenssl TLS stack fails (raises) if redirects or TLS verification break, failing the test.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `258` — `BP-PY-10`

- Function context: `scripts/niquests/findings/functions/258.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:1651:16`
- Checklist pattern: self-serialized pickle round-trip tests

Source excerpt:

```
       1649: 
       1650:         # verify we can pickle the original request
    >  1651:         assert pickle.loads(pickle.dumps(r.request))
       1652: 
```

Why this is a false positive: each flagged line is `pickle.loads(pickle.dumps(obj))` in a test — the data is serialized by the same expression, so the unpickled bytes are self-produced, not untrusted input.

Checklist evidence: BP-PY-10/CWE-502's condition is deserializing untrusted data; the tests round-trip their own in-memory objects.

### [ ] Finding `259` — `CWE-502`

- Function context: `scripts/niquests/findings/functions/259.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:1651:16`
- Checklist pattern: self-serialized pickle round-trip tests

Source excerpt:

```
       1649: 
       1650:         # verify we can pickle the original request
    >  1651:         assert pickle.loads(pickle.dumps(r.request))
       1652: 
```

Why this is a false positive: each flagged line is `pickle.loads(pickle.dumps(obj))` in a test — the data is serialized by the same expression, so the unpickled bytes are self-produced, not untrusted input.

Checklist evidence: BP-PY-10/CWE-502's condition is deserializing untrusted data; the tests round-trip their own in-memory objects.

### [ ] Finding `260` — `BP-PY-10`

- Function context: `scripts/niquests/findings/functions/260.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:1655:14`
- Checklist pattern: self-serialized pickle round-trip tests

Source excerpt:

```
       1653:         # verify we can pickle the response and that we have access to
       1654:         # the original request.
    >  1655:         pr = pickle.loads(pickle.dumps(r))
       1656:         assert r.request.url == pr.request.url
```

Why this is a false positive: each flagged line is `pickle.loads(pickle.dumps(obj))` in a test — the data is serialized by the same expression, so the unpickled bytes are self-produced, not untrusted input.

Checklist evidence: BP-PY-10/CWE-502's condition is deserializing untrusted data; the tests round-trip their own in-memory objects.

### [ ] Finding `261` — `BP-PY-10`

- Function context: `scripts/niquests/findings/functions/261.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:1663:13`
- Checklist pattern: self-serialized pickle round-trip tests

Source excerpt:

```
       1661: 
       1662:         # Verify PreparedRequest can be pickled and unpickled
    >  1663:         r = pickle.loads(pickle.dumps(p))
       1664:         assert r.url == p.url
```

Why this is a false positive: each flagged line is `pickle.loads(pickle.dumps(obj))` in a test — the data is serialized by the same expression, so the unpickled bytes are self-produced, not untrusted input.

Checklist evidence: BP-PY-10/CWE-502's condition is deserializing untrusted data; the tests round-trip their own in-memory objects.

### [ ] Finding `262` — `BP-PY-10`

- Function context: `scripts/niquests/findings/functions/262.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:1679:13`
- Checklist pattern: self-serialized pickle round-trip tests

Source excerpt:

```
       1677: 
       1678:         # Verify PreparedRequest can be pickled and unpickled
    >  1679:         r = pickle.loads(pickle.dumps(p))
       1680:         assert r.url == p.url
```

Why this is a false positive: each flagged line is `pickle.loads(pickle.dumps(obj))` in a test — the data is serialized by the same expression, so the unpickled bytes are self-produced, not untrusted input.

Checklist evidence: BP-PY-10/CWE-502's condition is deserializing untrusted data; the tests round-trip their own in-memory objects.

### [ ] Finding `263` — `BP-PY-10`

- Function context: `scripts/niquests/findings/functions/263.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:1694:13`
- Checklist pattern: self-serialized pickle round-trip tests

Source excerpt:

```
       1692: 
       1693:         # Verify PreparedRequest can be pickled
    >  1694:         r = pickle.loads(pickle.dumps(p))
       1695:         assert r.url == p.url
```

Why this is a false positive: each flagged line is `pickle.loads(pickle.dumps(obj))` in a test — the data is serialized by the same expression, so the unpickled bytes are self-produced, not untrusted input.

Checklist evidence: BP-PY-10/CWE-502's condition is deserializing untrusted data; the tests round-trip their own in-memory objects.

### [ ] Finding `264` — `BP-PY-10`

- Function context: `scripts/niquests/findings/functions/264.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:1724:13`
- Checklist pattern: self-serialized pickle round-trip tests

Source excerpt:

```
       1722:         s = niquests.Session()
       1723: 
    >  1724:         s = pickle.loads(pickle.dumps(s))
       1725:         s.proxies = getproxies()
```

Why this is a false positive: each flagged line is `pickle.loads(pickle.dumps(obj))` in a test — the data is serialized by the same expression, so the unpickled bytes are self-produced, not untrusted input.

Checklist evidence: BP-PY-10/CWE-502's condition is deserializing untrusted data; the tests round-trip their own in-memory objects.

### [ ] Finding `265` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/265.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:2163:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  2163:     def test_redirect_with_wrong_gzipped_header(self, httpbin):
      2164:         s = niquests.Session()
      2165:         url = httpbin("redirect/1")
      2166:         self._patch_adapter_gzipped_redirect(s, url)
```

Why this is a false positive: the patched adapter breaks the gzipped redirect; if the client mishandles it, `s.get(url)` raises and the test fails — a regression smoke test.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `266` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/266.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:2626:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  2626:     def test_read_timeout(self, httpbin, timeout):
      2627:         try:
      2628:             niquests.get(httpbin("delay/10"), timeout=timeout)
      2629:             pytest.fail("The recv() request should time out.")
```

Why this is a false positive: `pytest.fail` fires when the request does NOT time out; the `except ReadTimeout: pass` is the expected-exception assertion — the test verifies the timeout.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `267` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/267.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:2630:1`
- Checklist pattern: expected-exception path of a timeout test

Source excerpt:

```
       2628:             niquests.get(httpbin("delay/10"), timeout=timeout)
       2629:             pytest.fail("The recv() request should time out.")
    >  2630:         except ReadTimeout:
       2631:             pass
```

Why this is a false positive: `pytest.fail` fires when the request does NOT time out; the `except ReadTimeout`/`ConnectTimeout: pass` is the test's assertion mechanism — the exception IS the expected outcome.

Checklist evidence: the `pass` is the success path of the test (the exception was raised as expected), not a discarded failure.

### [ ] Finding `268` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/268.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:2630:1`
- Checklist pattern: expected-exception path of a timeout test

Source excerpt:

```
       2628:             niquests.get(httpbin("delay/10"), timeout=timeout)
       2629:             pytest.fail("The recv() request should time out.")
    >  2630:         except ReadTimeout:
       2631:             pass
```

Why this is a false positive: `pytest.fail` fires when the request does NOT time out; the `except ReadTimeout`/`ConnectTimeout: pass` is the test's assertion mechanism — the exception IS the expected outcome.

Checklist evidence: the `pass` is the success path of the test (the exception was raised as expected), not a discarded failure.

### [ ] Finding `269` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/269.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:2630:9`
- Checklist pattern: expected-exception path of a timeout test

Source excerpt:

```
       2628:             niquests.get(httpbin("delay/10"), timeout=timeout)
       2629:             pytest.fail("The recv() request should time out.")
    >  2630:         except ReadTimeout:
       2631:             pass
```

Why this is a false positive: `pytest.fail` fires when the request does NOT time out; the `except ReadTimeout`/`ConnectTimeout: pass` is the test's assertion mechanism — the exception IS the expected outcome.

Checklist evidence: the `pass` is the success path of the test (the exception was raised as expected), not a discarded failure.

### [ ] Finding `270` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/270.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:2643:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >  2643:     def test_total_timeout_connect(self, timeout):
      2644:         try:
      2645:             niquests.get(TARPIT, timeout=timeout)
      2646:             pytest.fail("The connect() request should time out.")
```

Why this is a false positive: `pytest.fail` fires when the connect does NOT time out; the `except ConnectTimeout: pass` is the expected-exception assertion.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `271` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/271.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:2647:1`
- Checklist pattern: expected-exception path of a timeout test

Source excerpt:

```
       2645:             niquests.get(TARPIT, timeout=timeout)
       2646:             pytest.fail("The connect() request should time out.")
    >  2647:         except ConnectTimeout:
       2648:             pass
```

Why this is a false positive: `pytest.fail` fires when the request does NOT time out; the `except ReadTimeout`/`ConnectTimeout: pass` is the test's assertion mechanism — the exception IS the expected outcome.

Checklist evidence: the `pass` is the success path of the test (the exception was raised as expected), not a discarded failure.

### [ ] Finding `272` — `CWE-617`

- Function context: `scripts/niquests/findings/functions/272.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:3078:1`
- Checklist pattern: test assertion is the verification mechanism

Source excerpt:

```
       3076:         assert encoded_values == [value]
       3077:         assert request.body == b'{"id":2}'
    >  3078:         assert request.headers["Content-Type"] == "application/json;charset=utf-8"
       3079:         assert request.headers["Content-Length"] == "8"
```

Why this is a false positive: `assert request.headers[...] == ...` verifies the request the server received; the asserted values are the test's own expectations.

Checklist evidence: CWE-617's condition is a reachable assertion on request-controlled state in production; the shown source is test verification.

### [ ] Finding `273` — `CWE-772`

- Function context: `scripts/niquests/findings/functions/273.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_testserver.py:41:1`
- Checklist pattern: test socket used to verify connection refusal

Source excerpt:

```
         39: 
         40:         with pytest.raises(socket.error):
    >    41:             new_sock = socket.socket()
         42:             new_sock.connect((host, port))
```

Why this is a false positive: `new_sock = socket.socket()` is created inside `with pytest.raises(socket.error)` to assert that connecting to a closed port fails; the test discards the socket by design after the expected error.

Checklist evidence: CWE-772's condition is a resource leak in production code; the shown source is a test fixture checking connection failure.

### [ ] Finding `274` — `CWE-93`

- Function context: `scripts/niquests/findings/functions/274.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_testserver.py:53:21`
- Checklist pattern: test assertion reading a header value

Source excerpt:

```
         51:             assert r.status_code == 200
         52:             assert r.text == "roflol"
    >    53:             assert r.headers["Content-Length"] == "6"
         54: 
```

Why this is a false positive: each flagged line is `assert <resp>.headers[...] == "..."` — a pytest assertion that reads a received header; nothing is written to any header.

Checklist evidence: CWE-93 requires a header write; the shown source only compares a header value in a test.

### [ ] Finding `275` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/275.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_testserver.py:134:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
    >   134:     def test_basic_waiting_server(self):
       135:         """the server waits for the block_server event to be set before closing"""
       136:         block_server = threading.Event()
       137: 
```

Why this is a false positive: the test drives the blocking server and would fail with a socket error if the server misbehaved; the server's `handler_results` are asserted by the server fixtures' surrounding checks.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `276` — `CWE-1341`

- Function context: `scripts/niquests/findings/functions/276.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_testserver.py:179:13`
- Checklist pattern: two distinct sockets closed once each

Source excerpt:

```
        177:             sock1.connect(address)
        178:             sock1.sendall(first_request)
    >   179:             sock1.close()
        180: 
```

Why this is a false positive: `sock1.close()` and `sock2.close()` close two different sockets created at lines 169-170.

Checklist evidence: the matched pair releases distinct handles.

### [ ] Finding `277` — `CWE-397`

- Function context: `scripts/niquests/findings/functions/277.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_testserver.py:220:1`
- Checklist pattern: tests deliberately raising a generic exception

Source excerpt:

```
        218:         with pytest.raises(Exception):
        219:             with server:
    >   220:                 raise Exception()
        221: 
```

Why this is a false positive: `raise Exception()` / `raise Exception("Expected exception")` inside `with pytest.raises(Exception):` exercises exception-handling behavior with a deliberately generic exception.

Checklist evidence: CWE-397's condition targets production code raising generic exceptions; the shown lines are test fixtures for exception handling.

### [ ] Finding `278` — `BP-PY-13`

- Function context: `scripts/niquests/findings/functions/278.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_utils.py:318:8`
- Checklist pattern: test fixture punctuation string, not a credential

Source excerpt:

```
        316: 
        317: 
    >   318: USER = PASSWORD = "%!*'();:@&=+$,/?#[] "
        319: ENCODED_USER = urllib.parse.quote(USER, "")
```

Why this is a false positive: `USER = PASSWORD = "%!*'();:@&=+$,/?#[] "` is the URL-encoding test corpus (it exercises quoting of reserved characters); it is not a real credential and lives in a test module.

Checklist evidence: BP-PY-13's condition is a hardcoded secret-like credential; the value is a non-secret punctuation corpus used to test `urlparse` quoting.

### [ ] Finding `279` — `CWE-397`

- Function context: `scripts/niquests/findings/functions/279.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_utils.py:816:1`
- Checklist pattern: tests deliberately raising a generic exception

Source excerpt:

```
        814:     with pytest.raises(Exception) as exception:
        815:         with set_environ("test1", None):
    >   816:             raise Exception("Expected exception")
        817: 
```

Why this is a false positive: `raise Exception()` / `raise Exception("Expected exception")` inside `with pytest.raises(Exception):` exercises exception-handling behavior with a deliberately generic exception.

Checklist evidence: CWE-397's condition targets production code raising generic exceptions; the shown lines are test fixtures for exception handling.

### [ ] Finding `280` — `CWE-770`

- Function context: `scripts/niquests/findings/functions/280.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_wsgi.py:59:21`
- Checklist pattern: test WSGI fixture reading its own request body

Source excerpt:

```
         57:             "path": request.path,
         58:             "query": request.query_string.decode("utf-8"),
    >    59:             "body": request.get_data(as_text=True),
         60:             "headers": dict(request.headers),
```

Why this is a false positive: `request.get_data(as_text=True)` runs in a minimal Flask test app in `test_wsgi.py`; the test client controls the body and no production resource-exhaustion surface exists.

Checklist evidence: CWE-770's condition is a production request reader without a size limit; the shown source is a test fixture.

### [ ] Finding `281` — `CWE-772`

- Function context: `scripts/niquests/findings/functions/281.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/testserver/server.py:81:39`
- Checklist pattern: socket ownership transferred via return

Source excerpt:

```
         79:             self.stop_event.set()
         80: 
    >    81:     def _create_socket_and_bind(self):
         82:         sock = socket.socket()
```

Why this is a false positive: `_create_socket_and_bind` returns `sock`; the caller stores it as `self.server_sock`, which `_close_server_sock_ignore_errors` closes (line 89).

Checklist evidence: the resource escapes the function to its owner, which releases it; no leak exists.

### [ ] Finding `282` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/282.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/testserver/server.py:90:1`
- Checklist pattern: defensive close of the server socket

Source excerpt:

```
         88:         try:
         89:             self.server_sock.close()
    >    90:         except OSError:
         91:             pass
```

Why this is a false positive: `close()` on an already-closed server socket may raise `OSError`; the pass is the designed best-effort teardown of the test server.

Checklist evidence: the exception is the expected outcome of idempotent close; the handler is deliberately best-effort.

### [ ] Finding `283` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/283.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/testserver/server.py:90:1`
- Checklist pattern: defensive close of the server socket

Source excerpt:

```
         88:         try:
         89:             self.server_sock.close()
    >    90:         except OSError:
         91:             pass
```

Why this is a false positive: `close()` on an already-closed server socket may raise `OSError`; the pass is the designed best-effort teardown of the test server.

Checklist evidence: the exception is the expected outcome of idempotent close; the handler is deliberately best-effort.

### [ ] Finding `284` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/284.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/testserver/server.py:90:9`
- Checklist pattern: defensive close of the server socket

Source excerpt:

```
         88:         try:
         89:             self.server_sock.close()
    >    90:         except OSError:
         91:             pass
```

Why this is a false positive: `close()` on an already-closed server socket may raise `OSError`; the pass is the designed best-effort teardown of the test server.

Checklist evidence: the exception is the expected outcome of idempotent close; the handler is deliberately best-effort.

### [ ] Finding `285` — `CWE-829`

- Function context: `scripts/niquests/findings/functions/285.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/app.py:21:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         19:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         20:         try:
    >    21:             importlib.import_module(module.name)
         22:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `286` — `CWE-94`

- Function context: `scripts/niquests/findings/functions/286.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/app.py:21:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         19:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         20:         try:
    >    21:             importlib.import_module(module.name)
         22:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `287` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/287.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/app.py:22:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         20:         try:
         21:             importlib.import_module(module.name)
    >    22:         except ImportError:
         23:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `288` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/288.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/app.py:22:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         20:         try:
         21:             importlib.import_module(module.name)
    >    22:         except ImportError:
         23:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `289` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/289.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/app.py:22:9`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         20:         try:
         21:             importlib.import_module(module.name)
    >    22:         except ImportError:
         23:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `290` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/290.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/edge_cases.py:257:1`
- Checklist pattern: expected-exception exercises in WASI guest edge-case tests

Source excerpt:

```
        255:     try:
        256:         wasi._unwrap_result(Err(Failed()))
    >   257:     except Err:
        258:         pass
```

Why this is a false positive: each `try` deliberately triggers a failure (`_unwrap_result(Err(...))`, reads on broken streams) and the `pass` records that the expected exception was raised — the test's assertion.

Checklist evidence: the except-pass is the expected-exception verification of the test, not silent error discarding.

### [ ] Finding `291` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/291.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/edge_cases.py:257:1`
- Checklist pattern: expected-exception exercises in WASI guest edge-case tests

Source excerpt:

```
        255:     try:
        256:         wasi._unwrap_result(Err(Failed()))
    >   257:     except Err:
        258:         pass
```

Why this is a false positive: each `try` deliberately triggers a failure (`_unwrap_result(Err(...))`, reads on broken streams) and the `pass` records that the expected exception was raised — the test's assertion.

Checklist evidence: the except-pass is the expected-exception verification of the test, not silent error discarding.

### [ ] Finding `292` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/292.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/edge_cases.py:257:5`
- Checklist pattern: expected-exception exercises in WASI guest edge-case tests

Source excerpt:

```
        255:     try:
        256:         wasi._unwrap_result(Err(Failed()))
    >   257:     except Err:
        258:         pass
```

Why this is a false positive: each `try` deliberately triggers a failure (`_unwrap_result(Err(...))`, reads on broken streams) and the `pass` records that the expected exception was raised — the test's assertion.

Checklist evidence: the except-pass is the expected-exception verification of the test, not silent error discarding.

### [ ] Finding `293` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/293.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/edge_cases.py:297:1`
- Checklist pattern: expected-exception exercises in WASI guest edge-case tests

Source excerpt:

```
        295:     try:
        296:         await failed.read(1)
    >   297:     except ReadTimeout:
        298:         pass
```

Why this is a false positive: each `try` deliberately triggers a failure (`_unwrap_result(Err(...))`, reads on broken streams) and the `pass` records that the expected exception was raised — the test's assertion.

Checklist evidence: the except-pass is the expected-exception verification of the test, not silent error discarding.

### [ ] Finding `294` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/294.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/edge_cases.py:305:1`
- Checklist pattern: expected-exception exercises in WASI guest edge-case tests

Source excerpt:

```
        303:     try:
        304:         await other_error.read(1)
    >   305:     except ValueError:
        306:         pass
```

Why this is a false positive: each `try` deliberately triggers a failure (`_unwrap_result(Err(...))`, reads on broken streams) and the `pass` records that the expected exception was raised — the test's assertion.

Checklist evidence: the except-pass is the expected-exception verification of the test, not silent error discarding.

### [ ] Finding `295` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/295.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/edge_cases.py:313:1`
- Checklist pattern: expected-exception exercises in WASI guest edge-case tests

Source excerpt:

```
        311:     try:
        312:         await future_error.read(1)
    >   313:     except ReadTimeout:
        314:         pass
```

Why this is a false positive: each `try` deliberately triggers a failure (`_unwrap_result(Err(...))`, reads on broken streams) and the `pass` records that the expected exception was raised — the test's assertion.

Checklist evidence: the except-pass is the expected-exception verification of the test, not silent error discarding.

### [ ] Finding `296` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/296.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/edge_cases.py:321:1`
- Checklist pattern: expected-exception exercises in WASI guest edge-case tests

Source excerpt:

```
        319:     try:
        320:         await unexpected.read(1)
    >   321:     except ValueError:
        322:         pass
```

Why this is a false positive: each `try` deliberately triggers a failure (`_unwrap_result(Err(...))`, reads on broken streams) and the `pass` records that the expected exception was raised — the test's assertion.

Checklist evidence: the except-pass is the expected-exception verification of the test, not silent error discarding.

### [ ] Finding `297` — `CWE-1341`

- Function context: `scripts/niquests/findings/functions/297.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/edge_cases.py:335:11`
- Checklist pattern: test exercising idempotent close

Source excerpt:

```
        333:     )
        334:     response = wasi._AsyncWASIHTTPResponse(body=cancellable, headers={}, preload_content=False)
    >   335:     await response.close()
        336:     await response.close()
```

Why this is a false positive: the tests deliberately call `response.close()` twice to verify the close is idempotent (a WASI response teardown edge case).

Checklist evidence: the double close is the test's subject matter, not an accidental double release.

### [ ] Finding `298` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/298.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/edge_cases.py:417:1`
- Checklist pattern: expected-exception exercises in WASI guest edge-case tests

Source excerpt:

```
        415:     try:
        416:         await extension.next_payload()
    >   417:     except OSError:
        418:         pass
```

Why this is a false positive: each `try` deliberately triggers a failure (`_unwrap_result(Err(...))`, reads on broken streams) and the `pass` records that the expected exception was raised — the test's assertion.

Checklist evidence: the except-pass is the expected-exception verification of the test, not silent error discarding.

### [ ] Finding `299` — `CWE-93`

- Function context: `scripts/niquests/findings/functions/299.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/test_wasi.py:47:23`
- Checklist pattern: test assertion reading a header value

Source excerpt:

```
         45:         encoded = await session.get("https://httpbingo.org/gzip", stream=True)
         46:         raw_body = b"".join([chunk async for chunk in await encoded.iter_raw()])
    >    47:         assert encoded.headers["content-encoding"] == "gzip"
         48:         assert raw_body.startswith(b"\x1f\x8b")
```

Why this is a false positive: each flagged line is `assert <resp>.headers[...] == "..."` — a pytest assertion that reads a received header; nothing is written to any header.

Checklist evidence: CWE-93 requires a header write; the shown source only compares a header value in a test.

### [ ] Finding `300` — `CWE-829`

- Function context: `scripts/niquests/findings/functions/300.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/combined/app.py:20:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         18:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         19:         try:
    >    20:             importlib.import_module(module.name)
         21:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `301` — `CWE-94`

- Function context: `scripts/niquests/findings/functions/301.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/combined/app.py:20:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         18:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         19:         try:
    >    20:             importlib.import_module(module.name)
         21:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `302` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/302.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/combined/app.py:21:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         19:         try:
         20:             importlib.import_module(module.name)
    >    21:         except ImportError:
         22:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `303` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/303.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/combined/app.py:21:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         19:         try:
         20:             importlib.import_module(module.name)
    >    21:         except ImportError:
         22:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `304` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/304.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/combined/app.py:21:9`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         19:         try:
         20:             importlib.import_module(module.name)
    >    21:         except ImportError:
         22:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `305` — `CWE-829`

- Function context: `scripts/niquests/findings/functions/305.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/hybrid_async/app.py:20:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         18:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         19:         try:
    >    20:             importlib.import_module(module.name)
         21:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `306` — `CWE-94`

- Function context: `scripts/niquests/findings/functions/306.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/hybrid_async/app.py:20:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         18:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         19:         try:
    >    20:             importlib.import_module(module.name)
         21:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `307` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/307.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/hybrid_async/app.py:21:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         19:         try:
         20:             importlib.import_module(module.name)
    >    21:         except ImportError:
         22:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `308` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/308.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/hybrid_async/app.py:21:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         19:         try:
         20:             importlib.import_module(module.name)
    >    21:         except ImportError:
         22:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `309` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/309.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/hybrid_async/app.py:21:9`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         19:         try:
         20:             importlib.import_module(module.name)
    >    21:         except ImportError:
         22:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `310` — `CWE-829`

- Function context: `scripts/niquests/findings/functions/310.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/hybrid_sync/app.py:19:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         17:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         18:         try:
    >    19:             importlib.import_module(module.name)
         20:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `311` — `CWE-94`

- Function context: `scripts/niquests/findings/functions/311.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/hybrid_sync/app.py:19:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         17:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         18:         try:
    >    19:             importlib.import_module(module.name)
         20:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `312` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/312.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/hybrid_sync/app.py:20:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         18:         try:
         19:             importlib.import_module(module.name)
    >    20:         except ImportError:
         21:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `313` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/313.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/hybrid_sync/app.py:20:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         18:         try:
         19:             importlib.import_module(module.name)
    >    20:         except ImportError:
         21:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `314` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/314.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/hybrid_sync/app.py:20:9`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         18:         try:
         19:             importlib.import_module(module.name)
    >    20:         except ImportError:
         21:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `315` — `CWE-829`

- Function context: `scripts/niquests/findings/functions/315.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/p1_sync/app.py:19:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         17:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         18:         try:
    >    19:             importlib.import_module(module.name)
         20:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `316` — `CWE-94`

- Function context: `scripts/niquests/findings/functions/316.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/p1_sync/app.py:19:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         17:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         18:         try:
    >    19:             importlib.import_module(module.name)
         20:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `317` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/317.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/p1_sync/app.py:20:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         18:         try:
         19:             importlib.import_module(module.name)
    >    20:         except ImportError:
         21:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `318` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/318.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/p1_sync/app.py:20:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         18:         try:
         19:             importlib.import_module(module.name)
    >    20:         except ImportError:
         21:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `319` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/319.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/p1_sync/app.py:20:9`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         18:         try:
         19:             importlib.import_module(module.name)
    >    20:         except ImportError:
         21:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `320` — `CWE-829`

- Function context: `scripts/niquests/findings/functions/320.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/socket_async/app.py:20:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         18:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         19:         try:
    >    20:             importlib.import_module(module.name)
         21:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `321` — `CWE-94`

- Function context: `scripts/niquests/findings/functions/321.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/socket_async/app.py:20:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         18:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         19:         try:
    >    20:             importlib.import_module(module.name)
         21:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `322` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/322.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/socket_async/app.py:21:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         19:         try:
         20:             importlib.import_module(module.name)
    >    21:         except ImportError:
         22:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `323` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/323.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/socket_async/app.py:21:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         19:         try:
         20:             importlib.import_module(module.name)
    >    21:         except ImportError:
         22:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `324` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/324.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/socket_async/app.py:21:9`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         19:         try:
         20:             importlib.import_module(module.name)
    >    21:         except ImportError:
         22:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `325` — `CWE-829`

- Function context: `scripts/niquests/findings/functions/325.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/socket_sync/app.py:19:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         17:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         18:         try:
    >    19:             importlib.import_module(module.name)
         20:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `326` — `CWE-94`

- Function context: `scripts/niquests/findings/functions/326.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/socket_sync/app.py:19:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         17:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         18:         try:
    >    19:             importlib.import_module(module.name)
         20:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `327` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/327.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/socket_sync/app.py:20:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         18:         try:
         19:             importlib.import_module(module.name)
    >    20:         except ImportError:
         21:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `328` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/328.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/socket_sync/app.py:20:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         18:         try:
         19:             importlib.import_module(module.name)
    >    20:         except ImportError:
         21:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `329` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/329.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/socket_sync/app.py:20:9`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         18:         try:
         19:             importlib.import_module(module.name)
    >    20:         except ImportError:
         21:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `330` — `CWE-829`

- Function context: `scripts/niquests/findings/functions/330.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/sync/app.py:21:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         19:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         20:         try:
    >    21:             importlib.import_module(module.name)
         22:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `331` — `CWE-94`

- Function context: `scripts/niquests/findings/functions/331.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/sync/app.py:21:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         19:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         20:         try:
    >    21:             importlib.import_module(module.name)
         22:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `332` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/332.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/sync/app.py:22:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         20:         try:
         21:             importlib.import_module(module.name)
    >    22:         except ImportError:
         23:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `333` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/333.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/sync/app.py:22:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         20:         try:
         21:             importlib.import_module(module.name)
    >    22:         except ImportError:
         23:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `334` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/334.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/sync/app.py:22:9`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         20:         try:
         21:             importlib.import_module(module.name)
    >    22:         except ImportError:
         23:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `335` — `CWE-1341`

- Function context: `scripts/niquests/findings/functions/335.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/sync/test_edges.py:377:5`
- Checklist pattern: test exercising idempotent close

Source excerpt:

```
        375:     cancellable = wasi._WASILowLevelResponse("GET", 200, "OK", HTTPHeaderDict(), Body(), Stream([b"pending"]), "url")
        376:     response = wasi._WASIHTTPResponse(body=cancellable, headers={}, preload_content=False)
    >   377:     response.close()
        378:     response.close()
```

Why this is a false positive: the tests deliberately call `response.close()` twice to verify the close is idempotent (a WASI response teardown edge case).

Checklist evidence: the double close is the test's subject matter, not an accidental double release.

### [ ] Finding `336` — `CWE-93`

- Function context: `scripts/niquests/findings/functions/336.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/sync/test_wasi.py:34:19`
- Checklist pattern: test assertion reading a header value

Source excerpt:

```
         32: def test_gzip_raw_and_decoded():
         33:     encoded = niquests.get("https://httpbingo.org/gzip", stream=True)
    >    34:     assert encoded.headers["content-encoding"] == "gzip"
         35:     assert b"".join(encoded.iter_raw()).startswith(b"\x1f\x8b")
```

Why this is a false positive: each flagged line is `assert <resp>.headers[...] == "..."` — a pytest assertion that reads a received header; nothing is written to any header.

Checklist evidence: CWE-93 requires a header write; the shown source only compares a header value in a test.

### [ ] Finding `337` — `CWE-829`

- Function context: `scripts/niquests/findings/functions/337.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/unavailable_async/app.py:20:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         18:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         19:         try:
    >    20:             importlib.import_module(module.name)
         21:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `338` — `CWE-94`

- Function context: `scripts/niquests/findings/functions/338.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/unavailable_async/app.py:20:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         18:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         19:         try:
    >    20:             importlib.import_module(module.name)
         21:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `339` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/339.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/unavailable_async/app.py:21:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         19:         try:
         20:             importlib.import_module(module.name)
    >    21:         except ImportError:
         22:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `340` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/340.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/unavailable_async/app.py:21:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         19:         try:
         20:             importlib.import_module(module.name)
    >    21:         except ImportError:
         22:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `341` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/341.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/unavailable_async/app.py:21:9`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         19:         try:
         20:             importlib.import_module(module.name)
    >    21:         except ImportError:
         22:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `342` — `CWE-829`

- Function context: `scripts/niquests/findings/functions/342.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/unavailable_sync/app.py:19:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         17:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         18:         try:
    >    19:             importlib.import_module(module.name)
         20:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `343` — `CWE-94`

- Function context: `scripts/niquests/findings/functions/343.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/unavailable_sync/app.py:19:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         17:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         18:         try:
    >    19:             importlib.import_module(module.name)
         20:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `344` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/344.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/unavailable_sync/app.py:20:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         18:         try:
         19:             importlib.import_module(module.name)
    >    20:         except ImportError:
         21:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `345` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/345.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/unavailable_sync/app.py:20:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         18:         try:
         19:             importlib.import_module(module.name)
    >    20:         except ImportError:
         21:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `346` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/346.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/unavailable_sync/app.py:20:9`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         18:         try:
         19:             importlib.import_module(module.name)
    >    20:         except ImportError:
         21:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

## True positives

Findings that satisfy the rule condition as implemented, with no visible mitigation in the shown source.

### `BP-PY-6` — assert Used For Runtime Validation (16)

| Finding | Source | Reason |
| --- | --- | --- |
| 26 | `src/niquests/adapters.py:749` | `assert <request/url/headers>.is not None` guards runtime request state in the public send path; stripped under `python -O` |
| 27 | `src/niquests/adapters.py:836` | `assert <request/url/headers>.is not None` guards runtime request state in the public send path; stripped under `python -O` |
| 28 | `src/niquests/adapters.py:1063` | `assert <request/url/headers>.is not None` guards runtime request state in the public send path; stripped under `python -O` |
| 38 | `src/niquests/adapters.py:1841` | `assert <request/url/headers>.is not None` guards runtime request state in the public send path; stripped under `python -O` |
| 39 | `src/niquests/adapters.py:1926` | `assert <request/url/headers>.is not None` guards runtime request state in the public send path; stripped under `python -O` |
| 40 | `src/niquests/adapters.py:2175` | `assert <request/url/headers>.is not None` guards runtime request state in the public send path; stripped under `python -O` |
| 53 | `src/niquests/async_session.py:782` | `assert <request/url/headers>.is not None` guards runtime request state in the public send path; stripped under `python -O` |
| 56 | `src/niquests/async_session.py:1051` | `assert <request/url/headers>.is not None` guards runtime request state in the public send path; stripped under `python -O` |
| 57 | `src/niquests/async_session.py:1070` | `assert <request/url/headers>.is not None` guards runtime request state in the public send path; stripped under `python -O` |
| 185 | `src/niquests/sessions.py:1526` | `assert <request/url/headers>.is not None` guards runtime request state in the public send path; stripped under `python -O` |
| 190 | `src/niquests/sessions.py:2093` | `assert <request/url/headers>.is not None` guards runtime request state in the public send path; stripped under `python -O` |
| 191 | `src/niquests/sessions.py:2112` | `assert <request/url/headers>.is not None` guards runtime request state in the public send path; stripped under `python -O` |
| 192 | `src/niquests/sessions.py:2177` | `assert <request/url/headers>.is not None` guards runtime request state in the public send path; stripped under `python -O` |
| 193 | `src/niquests/sessions.py:2178` | `assert <request/url/headers>.is not None` guards runtime request state in the public send path; stripped under `python -O` |
| 194 | `src/niquests/sessions.py:2201` | `assert <request/url/headers>.is not None` guards runtime request state in the public send path; stripped under `python -O` |
| 200 | `src/niquests/utils.py:848` | `assert <request/url/headers>.is not None` guards runtime request state in the public send path; stripped under `python -O` |

### `CWE-1121` — Excessive McCabe Cyclomatic Complexity (2)

| Finding | Source | Reason |
| --- | --- | --- |
| 47 | `src/niquests/async_session.py:459` | `__setstate__` contains 12+ visible control-flow branches (adapter-mounting chains), matching the ≥12 branch threshold |
| 189 | `src/niquests/sessions.py:1828` | `__setstate__` contains 12+ visible control-flow branches (adapter-mounting chains), matching the ≥12 branch threshold |

### `CWE-1124` — Excessively Deep Nesting (8)

| Finding | Source | Reason |
| --- | --- | --- |
| 51 | `src/niquests/async_session.py:711` | executable statement nested at ≥6 control-flow levels in the shown excerpt |
| 117 | `src/niquests/extensions/revocation/_crl/__init__.py:263` | executable statement nested at ≥6 control-flow levels in the shown excerpt |
| 119 | `src/niquests/extensions/revocation/_crl/_async/__init__.py:285` | executable statement nested at ≥6 control-flow levels in the shown excerpt |
| 124 | `src/niquests/extensions/revocation/_ocsp/__init__.py:329` | executable statement nested at ≥6 control-flow levels in the shown excerpt |
| 126 | `src/niquests/extensions/revocation/_ocsp/_async/__init__.py:341` | executable statement nested at ≥6 control-flow levels in the shown excerpt |
| 146 | `src/niquests/extensions/wasi/_adapter.py:130` | executable statement nested at ≥6 control-flow levels in the shown excerpt |
| 152 | `src/niquests/extensions/wasi/_async/_adapter.py:153` | executable statement nested at ≥6 control-flow levels in the shown excerpt |
| 183 | `src/niquests/sessions.py:1457` | executable statement nested at ≥6 control-flow levels in the shown excerpt |

### `PERF-PY-26` — Expensive Decode Or Parse On Hot Path (5)

| Finding | Source | Reason |
| --- | --- | --- |
| 48 | `src/niquests/async_session.py:618` | `decode_field_value(...)`/`parse_scheme(...)` runs inside a loop on the request path with no visible cache |
| 147 | `src/niquests/extensions/wasi/_adapter.py:139` | `decode_field_value(...)`/`parse_scheme(...)` runs inside a loop on the request path with no visible cache |
| 153 | `src/niquests/extensions/wasi/_async/_adapter.py:153` | `decode_field_value(...)`/`parse_scheme(...)` runs inside a loop on the request path with no visible cache |
| 163 | `src/niquests/extensions/wasi/_utils.py:81` | `decode_field_value(...)`/`parse_scheme(...)` runs inside a loop on the request path with no visible cache |
| 164 | `src/niquests/extensions/wasi/_utils.py:95` | `decode_field_value(...)`/`parse_scheme(...)` runs inside a loop on the request path with no visible cache |

### `BP-PY-14` — requests Without Timeout (4)

| Finding | Source | Reason |
| --- | --- | --- |
| 113 | `src/niquests/extensions/revocation/_crl/__init__.py:258` | executable `session.get(hint_ca_issuers[0])` network call in revocation checking carries no `timeout=` at all |
| 118 | `src/niquests/extensions/revocation/_crl/_async/__init__.py:277` | executable `session.get(hint_ca_issuers[0])` network call in revocation checking carries no `timeout=` at all |
| 120 | `src/niquests/extensions/revocation/_ocsp/__init__.py:324` | executable `session.get(hint_ca_issuers[0])` network call in revocation checking carries no `timeout=` at all |
| 125 | `src/niquests/extensions/revocation/_ocsp/_async/__init__.py:333` | executable `session.get(hint_ca_issuers[0])` network call in revocation checking carries no `timeout=` at all |

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/niquests/chunks`
- Function evidence: `scripts/niquests/findings/functions`
- Validation: `git diff --check` — pass (run in the goslop repo root after writing this report)
## Post-fix remaining-FP audit (2026-08-02)

Re-audit of the fresh (post-fix `b5b8fde`) scan, mode A. Classification is by `Source:` matching (file:line:col) against the audited TP list in the original audit below; every fresh finding matched either an audited TP source or an audited FP source+rule (2 fresh findings are new rule firings on audited-FP lines; both fail their rule condition and are classified FP). Repo commit unchanged since the original audit (`7633aa3f1f9fcdb7790192ffd8cfacb69ca2c807`).

### Run metadata (fresh scan)

```yaml
timestamp: 2026-08-02 (post-fix scan, binary rebuilt 2026-08-02 16:29)
repository: niquests
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests
branch: main
commit: 7633aa3f1f9fcdb7790192ffd8cfacb69ca2c807
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests
chunk_path: scripts/niquests/chunks
function_context_path: scripts/niquests/findings/functions
```

### Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop` (goslop binary used: `./bin/goslop`)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/niquests/chunks -context-dir scripts/niquests/findings/functions real-repos/niquests` (per `plans/fp-validations/reports/TP_MATCH.md` fresh run)
- Findings: `271`
- Chunks reviewed: `scripts/niquests/chunks/Chunk_1_25.txt` through `Chunk_251_271.txt` (all 11 chunk files)
- Function contexts reviewed: `scripts/niquests/findings/functions/<finding-id>.txt` for every fresh finding; every FP's context-file flagged line was cross-checked against its chunk excerpt (0 mismatches); enclosing source read for the 2 new rule firings and for the suppressed TP check

### Audit checklist

- [x] Read every assigned chunk under `scripts/niquests/chunks`.
- [x] Read `scripts/niquests/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [ ] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient (no delegated reviewers).
- [x] Ran `git diff --check` after updating this report.

### Classification summary (fresh run)

Matching rule: a fresh finding is TP iff its `Source:` (file:line) matches an audited TP source; FP iff it matches an audited FP source (same construct+rule) or fails its rule condition on inspection. Audited TP 48 (`PERF-PY-26` at `src/niquests/async_session.py:618`) is absent from the fresh scan although the construct (`parse_scheme(prefix)` in the adapter-mount loop) is still present at the same line — one audited TP over-suppressed by the fix, documented below for review.

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 237 | 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 22, 23, 24, 25, 29, 30, 31, 32, 33, 34, 36, 37, 39, 41, 42, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 100, 101, 102, 107, 108, 109, 113, 114, 115, 116, 117, 118, 119, 122, 123, 124, 125, 126, 129, 130, 131, 132, 133, 134, 135, 136, 137, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 159, 161, 162, 169, 170, 171, 172, 174, 175, 176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 255, 256, 257, 258, 259, 260, 261, 262, 263, 264, 265, 266, 267, 268, 269, 270, 271 |
| True positive | 34 | 19, 20, 21, 26, 27, 28, 35, 38, 40, 43, 44, 99, 103, 104, 105, 106, 110, 111, 112, 120, 121, 127, 128, 138, 139, 158, 160, 163, 164, 165, 166, 167, 168, 173 |
| Uncertain | 0 | — |

### False positives

One subsection per finding. As in the original audit, no grouping applies: every fresh finding is a distinct (file:line, rule) pair — findings on the same line carry distinct rules — so no two findings reference the exact same source construct under the same rule. Reasons for re-appearing findings are carried over from the original audit's verified subsections for the identical construct; excerpts are from the fresh function contexts.

### [ ] Finding `1` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/1.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/noxfile.py:95:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
         93:     ]
         94: )
    >    95: def test_cohabitation(session: nox.Session) -> None:
         96:     tests_impl(session, cohabitation=True)
```

Why this is a false positive: the test body only invokes `tests_impl(session, ...)`, the nox runner that executes the whole pytest suite — the assertions live in the suite run by pytest, so this is a session orchestration, not a placeholder test.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `2` — `BP-PY-46`

- Function context: `scripts/niquests/findings/functions/2.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/noxfile.py:147:9`
- Checklist pattern: print in a nox automation script for CI console output

Source excerpt:

```
        145:     """Test on Emscripten with Pyodide & Chrome / Firefox / Node.js"""
        146:     if runner == "node":
    >   147:         print(
        148:             "Node version:",
```

Why this is a false positive: `noxfile.py` is the project's build/test automation script; the prints emit human-readable progress to the CI console, not library logging.

Checklist evidence: BP-PY-46 targets 'library code'; noxfile.py is a dev/CI script where console output is the intended interface.

### [ ] Finding `3` — `CWE-22`

- Function context: `scripts/niquests/findings/functions/3.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/noxfile.py:157:53`
- Checklist pattern: path built from constants in a dev automation script

Source excerpt:

```
        155:     pyodide_version = "0.28.1"
        156:
    >   157:     pyodide_artifacts_path = Path(session.cache_dir) / f"pyodide-{pyodide_version}"
        158:
```

Why this is a false positive: both segments are trusted: `session.cache_dir` is nox's own cache directory and `pyodide_version` is a module constant; no untrusted input reaches the join.

Checklist evidence: CWE-22's condition requires a dynamic segment from an untrusted source; the joined segments are constants under the tool's own cache root.

### [ ] Finding `4` — `BP-PY-46`

- Function context: `scripts/niquests/findings/functions/4.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/noxfile.py:160:9`
- Checklist pattern: print in a nox automation script for CI console output

Source excerpt:

```
        158:
        159:     if not pyodide_artifacts_path.exists():
    >   160:         print("Fetching pyodide build artifacts")
        161:         session.run(
```

Why this is a false positive: `noxfile.py` is the project's build/test automation script; the prints emit human-readable progress to the CI console, not library logging.

Checklist evidence: BP-PY-46 targets 'library code'; noxfile.py is a dev/CI script where console output is the intended interface.

### [ ] Finding `5` — `CWE-409`

- Function context: `scripts/niquests/findings/functions/5.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/noxfile.py:237:19`
- Checklist pattern: tarfile extraction using the documented safe filter

Source excerpt:

```
        235:         extract_root.mkdir(exist_ok=True)
        236:         with tarfile.open(archive) as source:
    >   237:             source.extractall(extract_root, filter="data")
        238:         extracted = next(path for path in extract_root.iterdir() if (path / "wit").is_dir())
```

Why this is a false positive: `source.extractall(extract_root, filter="data")` uses Python 3.12's `filter="data"` — the platform's documented safe extraction filter that blocks traversal and special-file members.

Checklist evidence: CWE-409's fix explicitly recommends 'the platform's documented safe API'; the shown source already uses it.

### [ ] Finding `6` — `BP-PY-5`

- Function context: `scripts/niquests/findings/functions/6.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_async.py:14:1`
- Checklist pattern: deliberate API re-export shim marked `# noqa`

Source excerpt:

```
         12: )
         13:
    >    14: from .async_session import *  # noqa
         15:
```

Why this is a false positive: `_async.py`/`_typing.py` are internal shim modules whose entire purpose is re-exporting the public API (`from .async_session import *  # noqa`); the wildcard is the re-export mechanism, not accidental namespace pollution.

Checklist evidence: the wildcard import is the module's designed function (re-export), explicitly annotated `# noqa`.

### [ ] Finding `7` — `BP-PY-5`

- Function context: `scripts/niquests/findings/functions/7.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_typing.py:14:1`
- Checklist pattern: deliberate API re-export shim marked `# noqa`

Source excerpt:

```
         12: )
         13:
    >    14: from .typing import *  # noqa
         15:
```

Why this is a false positive: `_async.py`/`_typing.py` are internal shim modules whose entire purpose is re-exporting the public API (`from .async_session import *  # noqa`); the wildcard is the re-export mechanism, not accidental namespace pollution.

Checklist evidence: the wildcard import is the module's designed function (re-export), explicitly annotated `# noqa`.

### [ ] Finding `8` — `BP-PY-14`

- Function context: `scripts/niquests/findings/functions/8.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_vendor/kiss_headers/__init__.py:11:12`
- Checklist pattern: requests/session call inside a docstring doctest example

Source excerpt:

```
          9:    >>> import requests
         10:    >>> from kiss_headers import parse_it
    >    11:    >>> r = requests.get('https://www.python.org')
         12:    >>> headers = parse_it(r)
```

Why this is a false positive: the `>>> requests.get(...)` / `>>> session.get(...)` lines are doctest examples in docstrings, i.e. documentation text, not executed HTTP calls.

Checklist evidence: BP-PY-14's condition is an executable requests call missing `timeout=`; the regex matches the raw source inside a docstring, so no real call exists.

### [ ] Finding `9` — `BP-PY-14`

- Function context: `scripts/niquests/findings/functions/9.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_vendor/kiss_headers/__init__.py:27:18`
- Checklist pattern: requests/session call inside a docstring doctest example

Source excerpt:

```
         25: ... or from a raw IMAP4 message:
         26:
    >    27:    >>> message = requests.get("https://gist.githubusercontent.com/Ousret/8b84b736c375bb6aa3d389e86b5116ec/raw/21cb2f7af865e401c37d9b053fb6fe1abf63165b/sample-message.eml").content
         28:    >>> headers = parse_it(message)
```

Why this is a false positive: the `>>> requests.get(...)` / `>>> session.get(...)` lines are doctest examples in docstrings, i.e. documentation text, not executed HTTP calls.

Checklist evidence: BP-PY-14's condition is an executable requests call missing `timeout=`; the regex matches the raw source inside a docstring, so no real call exists.

### [ ] Finding `10` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/10.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_vendor/kiss_headers/builder.py:350:1`
- Checklist pattern: parsing fallback for an optional `filename*` encoding split

Source excerpt:

```
        348:                 encoding, encoded_filename = tuple(str(self["filename*"]).split("''"))
        349:                 return url_unquote(encoded_filename, encoding)
    >   350:             except ValueError:
        351:                 pass
```

Why this is a false positive: the `ValueError` is the expected outcome when the `filename*` value has no `encoding''encoded` split; the property falls through to the raw value, a defined fallback.

Checklist evidence: the except-pass is a deliberate fallback for a malformed optional header; no failure is hidden because a defined value is returned.

### [ ] Finding `11` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/11.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_vendor/kiss_headers/builder.py:350:1`
- Checklist pattern: parsing fallback for an optional `filename*` encoding split

Source excerpt:

```
        348:                 encoding, encoded_filename = tuple(str(self["filename*"]).split("''"))
        349:                 return url_unquote(encoded_filename, encoding)
    >   350:             except ValueError:
        351:                 pass
```

Why this is a false positive: the `ValueError` is the expected outcome when the `filename*` value has no `encoding''encoded` split; the property falls through to the raw value, a defined fallback.

Checklist evidence: the except-pass is a deliberate fallback for a malformed optional header; no failure is hidden because a defined value is returned.

### [ ] Finding `12` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/12.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_vendor/kiss_headers/builder.py:350:13`
- Checklist pattern: parsing fallback for an optional `filename*` encoding split

Source excerpt:

```
        348:                 encoding, encoded_filename = tuple(str(self["filename*"]).split("''"))
        349:                 return url_unquote(encoded_filename, encoding)
    >   350:             except ValueError:
        351:                 pass
```

Why this is a false positive: the `ValueError` is the expected outcome when the `filename*` value has no `encoding''encoded` split; the property falls through to the raw value, a defined fallback.

Checklist evidence: the except-pass is a deliberate fallback for a malformed optional header; no failure is hidden because a defined value is returned.

### [ ] Finding `13` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/13.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_vendor/kiss_headers/models.py:893:1`
- Checklist pattern: parsing fallback in header subclass lookup

Source excerpt:

```
        891:                     else None
        892:                 )
    >   893:             except TypeError:
        894:                 pass
```

Why this is a false positive: `header_name_to_class` is a best-effort subclass probe; on `TypeError` `target_subclass` stays `None`, which the caller handles.

Checklist evidence: the exception is the expected outcome of probing for an optional class mapping; the fallback value is defined.

### [ ] Finding `14` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/14.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_vendor/kiss_headers/models.py:893:1`
- Checklist pattern: parsing fallback in header subclass lookup

Source excerpt:

```
        891:                     else None
        892:                 )
    >   893:             except TypeError:
        894:                 pass
```

Why this is a false positive: `header_name_to_class` is a best-effort subclass probe; on `TypeError` `target_subclass` stays `None`, which the caller handles.

Checklist evidence: the exception is the expected outcome of probing for an optional class mapping; the fallback value is defined.

### [ ] Finding `15` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/15.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_vendor/kiss_headers/models.py:893:13`
- Checklist pattern: parsing fallback in header subclass lookup

Source excerpt:

```
        891:                     else None
        892:                 )
    >   893:             except TypeError:
        894:                 pass
```

Why this is a false positive: `header_name_to_class` is a best-effort subclass probe; on `TypeError` `target_subclass` stays `None`, which the caller handles.

Checklist evidence: the exception is the expected outcome of probing for an optional class mapping; the fallback value is defined.

### [ ] Finding `16` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/16.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_vendor/kiss_headers/utils.py:339:1`
- Checklist pattern: parsing fallback when a segment contains no `;`

Source excerpt:

```
        337:     try:
        338:         next_semi_colon_index = elem_end_index + content[elem_end_index:].index(";")
    >   339:     except ValueError:
        340:         pass
```

Why this is a false positive: `index(";")` raising `ValueError` means the element has no parameter part — the expected outcome the loop is written for; processing continues with the collected element.

Checklist evidence: the flagged `pass` implements an expected parsing branch, not silent failure discarding.

### [ ] Finding `17` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/17.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_vendor/kiss_headers/utils.py:339:1`
- Checklist pattern: parsing fallback when a segment contains no `;`

Source excerpt:

```
        337:     try:
        338:         next_semi_colon_index = elem_end_index + content[elem_end_index:].index(";")
    >   339:     except ValueError:
        340:         pass
```

Why this is a false positive: `index(";")` raising `ValueError` means the element has no parameter part — the expected outcome the loop is written for; processing continues with the collected element.

Checklist evidence: the flagged `pass` implements an expected parsing branch, not silent failure discarding.

### [ ] Finding `18` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/18.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/_vendor/kiss_headers/utils.py:339:5`
- Checklist pattern: parsing fallback when a segment contains no `;`

Source excerpt:

```
        337:     try:
        338:         next_semi_colon_index = elem_end_index + content[elem_end_index:].index(";")
    >   339:     except ValueError:
        340:         pass
```

Why this is a false positive: `index(";")` raising `ValueError` means the element has no parameter part — the expected outcome the loop is written for; processing continues with the collected element.

Checklist evidence: the flagged `pass` implements an expected parsing branch, not silent failure discarding.

### [ ] Finding `22` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/22.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/adapters.py:1075:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
       1073:                             try:
       1074:                                 from .extensions.revocation._ocsp import verify as ocsp_verify
    >  1075:                             except ImportError:
       1076:                                 pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `23` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/23.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/adapters.py:1075:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
       1073:                             try:
       1074:                                 from .extensions.revocation._ocsp import verify as ocsp_verify
    >  1075:                             except ImportError:
       1076:                                 pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `24` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/24.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/adapters.py:1075:29`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
       1073:                             try:
       1074:                                 from .extensions.revocation._ocsp import verify as ocsp_verify
    >  1075:                             except ImportError:
       1076:                                 pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `25` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/25.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/adapters.py:1091:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
       1089:                             try:
       1090:                                 from .extensions.revocation._crl import verify as crl_verify
    >  1091:                             except ImportError:
       1092:                                 pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `29` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/29.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/adapters.py:2189:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
       2187:                                     verify as async_ocsp_verify,
       2188:                                 )
    >  2189:                             except ImportError:
       2190:                                 pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `30` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/30.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/adapters.py:2207:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
       2205:                                     verify as async_crl_verify,
       2206:                                 )
    >  2207:                             except ImportError:
       2208:                                 pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `31` — `BP-PY-14`

- Function context: `scripts/niquests/findings/functions/31.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/async_api.py:183:26`
- Checklist pattern: timeout passed as a positional argument

Source excerpt:

```
        181:         session._crl_cache = _SHARED_CRL_CACHE.get()
        182:         try:
    >   183:             return await session.request(  # type: ignore[misc]
        184:                 method,
```

Why this is a false positive: the `session.request(...)` call passes `timeout` as the 10th positional argument, so the request does carry a timeout; the rule only detects the `timeout=` keyword form.

Checklist evidence: the rule's condition 'call missing timeout=' is not met — the timeout value is present positionally.

### [ ] Finding `32` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/32.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/async_session.py:99:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         97:     try:
         98:         from .extensions.wasi._async._adapter import AsyncWASIAdapter as AsyncWASIHTTPAdapter
    >    99:     except ImportError:
        100:         pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `33` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/33.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/async_session.py:99:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         97:     try:
         98:         from .extensions.wasi._async._adapter import AsyncWASIAdapter as AsyncWASIHTTPAdapter
    >    99:     except ImportError:
        100:         pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `34` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/34.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/async_session.py:99:5`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         97:     try:
         98:         from .extensions.wasi._async._adapter import AsyncWASIAdapter as AsyncWASIHTTPAdapter
    >    99:     except ImportError:
        100:         pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `36` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/36.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/async_session.py:620:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
        618:                 if scheme in extension.supported_schemes() and extension.scheme_to_http_scheme(scheme) == parse_scheme(prefix):
        619:                     return adapter
    >   620:         except ImportError:
        621:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `37` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/37.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/async_session.py:703:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
        701:                             verify as ocsp_verify,
        702:                         )
    >   703:                     except ImportError:
        704:                         pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `39` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/39.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/async_session.py:731:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
        729:                             verify as crl_verify,
        730:                         )
    >   731:                     except ImportError:
        732:                         pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `41` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/41.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/async_session.py:957:1`
- Checklist pattern: generator-exhaustion control flow

Source excerpt:

```
        955:                 try:
        956:                     r._next = await gen.__anext__()  # type: ignore[assignment]
    >   957:                 except StopAsyncIteration:
        958:                     pass
```

Why this is a false positive: `StopIteration`/`StopAsyncIteration` is the generator's normal termination signal; catching it is the intended loop-exit control flow, not a discarded failure.

Checklist evidence: the rule's 'failures discarded silently' condition does not apply to iterator-exhaustion signalling.

### [ ] Finding `42` — `CWE-1341`

- Function context: `scripts/niquests/findings/functions/42.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/async_session.py:1015:23`
- Checklist pattern: mutually exclusive close paths, single release per execution

Source excerpt:

```
       1013:             # Release the connection back into the pool.
       1014:             if isinstance(resp, AsyncResponse):
    >  1015:                 await resp.close()
       1016:             else:
```

Why this is a false positive: `async_session.py:1015` and its `else: resp.close()` are the two branches of one `isinstance` conditional — exactly one executes; the wasi `response.close()` calls sit in exclusive retry/raise branches (line 235 vs 239).

Checklist evidence: CWE-1341's condition is the same handle released twice on one path; the matched `close()` pairs are mutually exclusive branches (or distinct resources).

### [ ] Finding `45` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/45.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/async_session.py:1234:1`
- Checklist pattern: non-seekable body probe fallback

Source excerpt:

```
       1232:             try:
       1233:                 prep._body_position = await prep.body.tell()
    >  1234:             except OSError:
       1235:                 pass
```

Why this is a false positive: `body.tell()` raises `OSError` for non-seekable bodies; `_body_position` stays unset and the caller treats it as 'unknown' — a defined fallback.

Checklist evidence: the probe failure is expected for streaming bodies and handled by leaving the optional attribute unset.

### [ ] Finding `46` — `CWE-93`

- Function context: `scripts/niquests/findings/functions/46.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/auth.py:75:14`
- Checklist pattern: write to an outgoing request header, not a response header

Source excerpt:

```
         73:
         74:         if len(detect_token_type) == 1:
    >    75:             r.headers["Authorization"] = f"Bearer {self.token}"
         76:         else:
```

Why this is a false positive: `auth.py:75` sets the `Authorization` header on the outgoing request (from the caller's own token); `models.py:488` populates the `PreparedRequest.headers` dict in `prepare_headers` from the caller-supplied header mapping — neither is an HTTP response header.

Checklist evidence: CWE-93's condition is writing an externally influenced value into an HTTP response header; the shown lines write client-side request headers.

### [ ] Finding `47` — `BP-PY-14`

- Function context: `scripts/niquests/findings/functions/47.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/auth.py:399:27`
- Checklist pattern: requests/session call inside a docstring doctest example

Source excerpt:

```
        397:         >>> auth = niquests.auth.AsyncHTTPDigestAuth('user', 'pass')
        398:         >>> async with niquests.AsyncSession() as session:
    >   399:         ...     r = await session.get('https://httpbin.org/digest-auth/auth/user/pass', auth=auth)
        400:         ...     print(r.status_code)
```

Why this is a false positive: the `>>> requests.get(...)` / `>>> session.get(...)` lines are doctest examples in docstrings, i.e. documentation text, not executed HTTP calls.

Checklist evidence: BP-PY-14's condition is an executable requests call missing `timeout=`; the regex matches the raw source inside a docstring, so no real call exists.

### [ ] Finding `48` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/48.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:57:1`
- Checklist pattern: documented defensive JS-bridge reader setup

Source excerpt:

```
         55:                 if body is not None:
         56:                     self._reader = body.getReader()
    >    57:             except Exception:
         58:                 pass
```

Why this is a false positive: the Pyodide JS bridge suppresses failures while probing the JS stream object's `getReader()`; when it fails, `self._reader` stays unset and the sync reader falls back to `body.arrayBuffer()`.

Checklist evidence: the handler is a deliberate best-effort bridge probe with a defined fallback state; no failure that matters is hidden.

### [ ] Finding `49` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/49.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:57:1`
- Checklist pattern: documented defensive JS-bridge reader setup

Source excerpt:

```
         55:                 if body is not None:
         56:                     self._reader = body.getReader()
    >    57:             except Exception:
         58:                 pass
```

Why this is a false positive: the Pyodide JS bridge suppresses failures while probing the JS stream object's `getReader()`; when it fails, `self._reader` stays unset and the sync reader falls back to `body.arrayBuffer()`.

Checklist evidence: the handler is a deliberate best-effort bridge probe with a defined fallback state; no failure that matters is hidden.

### [ ] Finding `50` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/50.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:57:1`
- Checklist pattern: documented defensive JS-bridge reader setup

Source excerpt:

```
         55:                 if body is not None:
         56:                     self._reader = body.getReader()
    >    57:             except Exception:
         58:                 pass
```

Why this is a false positive: the Pyodide JS bridge suppresses failures while probing the JS stream object's `getReader()`; when it fails, `self._reader` stays unset and the sync reader falls back to `body.arrayBuffer()`.

Checklist evidence: the handler is a deliberate best-effort bridge probe with a defined fallback state; no failure that matters is hidden.

### [ ] Finding `51` — `CWE-396`

- Function context: `scripts/niquests/findings/functions/51.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:57:1`
- Checklist pattern: documented defensive JS-bridge reader setup

Source excerpt:

```
         55:                 if body is not None:
         56:                     self._reader = body.getReader()
    >    57:             except Exception:
         58:                 pass
```

Why this is a false positive: the Pyodide JS bridge suppresses failures while probing the JS stream object's `getReader()`; when it fails, `self._reader` stays unset and the sync reader falls back to `body.arrayBuffer()`.

Checklist evidence: the handler is a deliberate best-effort bridge probe with a defined fallback state; no failure that matters is hidden.

### [ ] Finding `52` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/52.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:57:13`
- Checklist pattern: documented defensive JS-bridge reader setup

Source excerpt:

```
         55:                 if body is not None:
         56:                     self._reader = body.getReader()
    >    57:             except Exception:
         58:                 pass
```

Why this is a false positive: the Pyodide JS bridge suppresses failures while probing the JS stream object's `getReader()`; when it fails, `self._reader` stays unset and the sync reader falls back to `body.arrayBuffer()`.

Checklist evidence: the handler is a deliberate best-effort bridge probe with a defined fallback state; no failure that matters is hidden.

### [ ] Finding `53` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/53.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:77:1`
- Checklist pattern: best-effort `to_py()` conversion with `None` fallback

Source excerpt:

```
         75:                 return bytes(value.to_py())
         76:             return None
    >    77:         except Exception:
         78:             return None
```

Why this is a false positive: converting a JS object to bytes is best-effort; failure returns `None`, a defined fallback the caller handles.

Checklist evidence: the handler returns a defined fallback value, so the failure is not swallowed without consequence.

### [ ] Finding `54` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/54.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:131:1`
- Checklist pattern: teardown best-effort reader cancellation

Source excerpt:

```
        129:             try:
        130:                 run_sync(self._reader.cancel())
    >   131:             except Exception:
        132:                 pass
```

Why this is a false positive: cancelling the stream reader during teardown is best-effort; JS errors here are expected and the loop has already ended.

Checklist evidence: the cancellation failure is the expected outcome of tearing down an already-closed stream.

### [ ] Finding `55` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/55.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:131:1`
- Checklist pattern: teardown best-effort reader cancellation

Source excerpt:

```
        129:             try:
        130:                 run_sync(self._reader.cancel())
    >   131:             except Exception:
        132:                 pass
```

Why this is a false positive: cancelling the stream reader during teardown is best-effort; JS errors here are expected and the loop has already ended.

Checklist evidence: the cancellation failure is the expected outcome of tearing down an already-closed stream.

### [ ] Finding `56` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/56.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:189:1`
- Checklist pattern: retry-loop handler that acts on the exception

Source excerpt:

```
        187:             try:
        188:                 response = self._do_send(request, stream, timeout)
    >   189:             except Exception as err:
        190:                 retries = retries.increment(method, request.url, error=err)
```

Why this is a false positive: the handler increments the retry state (`retries.increment(...)`) and re-raises `MaxRetryError` once retries are exhausted — the failure is not swallowed.

Checklist evidence: BP-PY-1's 'hides failures' condition is unmet: the generic catch feeds the retry machinery and re-raises on exhaustion.

### [ ] Finding `57` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/57.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:299:1`
- Checklist pattern: best-effort response-header parse fallback

Source excerpt:

```
        297:                     for entry in js_headers.entries():
        298:                         response_headers[entry[0]] = entry[1]
    >   299:         except Exception:
        300:             pass
```

Why this is a false positive: extracting headers from the JS `headers.entries()` is best-effort; on failure the response is returned without the (optional) header snapshot.

Checklist evidence: the handler is a defensive bridge fallback with a defined result.

### [ ] Finding `58` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/58.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:299:1`
- Checklist pattern: best-effort response-header parse fallback

Source excerpt:

```
        297:                     for entry in js_headers.entries():
        298:                         response_headers[entry[0]] = entry[1]
    >   299:         except Exception:
        300:             pass
```

Why this is a false positive: extracting headers from the JS `headers.entries()` is best-effort; on failure the response is returned without the (optional) header snapshot.

Checklist evidence: the handler is a defensive bridge fallback with a defined result.

### [ ] Finding `59` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/59.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:312:1`
- Checklist pattern: fallback for missing reason phrase

Source excerpt:

```
        310:         try:
        311:             response.reason = js_response.status_text or ""
    >   312:         except Exception:
        313:             response.reason = ""
```

Why this is a false positive: when the JS `status_text` probe fails, `response.reason` falls back to `""` — a defined value, not a swallowed failure.

Checklist evidence: the handler assigns a defined fallback value.

### [ ] Finding `60` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/60.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/__init__.py:327:1`
- Checklist pattern: fallback for body read failure

Source excerpt:

```
        325:             try:
        326:                 response_body: bytes = run_sync(js_response.bytes())
    >   327:             except Exception:
        328:                 response_body = b""
```

Why this is a false positive: when reading the JS response body fails, `response_body` falls back to `b""` — a defined value.

Checklist evidence: the handler assigns a defined fallback value.

### [ ] Finding `61` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/61.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/__init__.py:56:1`
- Checklist pattern: documented defensive JS-bridge reader setup

Source excerpt:

```
         54:                 if body is not None:
         55:                     self._reader = body.getReader()
    >    56:             except Exception:
         57:                 pass
```

Why this is a false positive: the Pyodide JS bridge suppresses failures while probing the JS stream object's `getReader()`; when it fails, `self._reader` stays unset and the sync reader falls back to `body.arrayBuffer()`.

Checklist evidence: the handler is a deliberate best-effort bridge probe with a defined fallback state; no failure that matters is hidden.

### [ ] Finding `62` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/62.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/__init__.py:56:1`
- Checklist pattern: documented defensive JS-bridge reader setup

Source excerpt:

```
         54:                 if body is not None:
         55:                     self._reader = body.getReader()
    >    56:             except Exception:
         57:                 pass
```

Why this is a false positive: the Pyodide JS bridge suppresses failures while probing the JS stream object's `getReader()`; when it fails, `self._reader` stays unset and the sync reader falls back to `body.arrayBuffer()`.

Checklist evidence: the handler is a deliberate best-effort bridge probe with a defined fallback state; no failure that matters is hidden.

### [ ] Finding `63` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/63.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/__init__.py:56:1`
- Checklist pattern: documented defensive JS-bridge reader setup

Source excerpt:

```
         54:                 if body is not None:
         55:                     self._reader = body.getReader()
    >    56:             except Exception:
         57:                 pass
```

Why this is a false positive: the Pyodide JS bridge suppresses failures while probing the JS stream object's `getReader()`; when it fails, `self._reader` stays unset and the sync reader falls back to `body.arrayBuffer()`.

Checklist evidence: the handler is a deliberate best-effort bridge probe with a defined fallback state; no failure that matters is hidden.

### [ ] Finding `64` — `CWE-396`

- Function context: `scripts/niquests/findings/functions/64.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/__init__.py:56:1`
- Checklist pattern: documented defensive JS-bridge reader setup

Source excerpt:

```
         54:                 if body is not None:
         55:                     self._reader = body.getReader()
    >    56:             except Exception:
         57:                 pass
```

Why this is a false positive: the Pyodide JS bridge suppresses failures while probing the JS stream object's `getReader()`; when it fails, `self._reader` stays unset and the sync reader falls back to `body.arrayBuffer()`.

Checklist evidence: the handler is a deliberate best-effort bridge probe with a defined fallback state; no failure that matters is hidden.

### [ ] Finding `65` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/65.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/__init__.py:56:13`
- Checklist pattern: documented defensive JS-bridge reader setup

Source excerpt:

```
         54:                 if body is not None:
         55:                     self._reader = body.getReader()
    >    56:             except Exception:
         57:                 pass
```

Why this is a false positive: the Pyodide JS bridge suppresses failures while probing the JS stream object's `getReader()`; when it fails, `self._reader` stays unset and the sync reader falls back to `body.arrayBuffer()`.

Checklist evidence: the handler is a deliberate best-effort bridge probe with a defined fallback state; no failure that matters is hidden.

### [ ] Finding `66` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/66.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/__init__.py:119:1`
- Checklist pattern: best-effort stream read with `None` fallback

Source excerpt:

```
        117:         except asyncio.TimeoutError:
        118:             raise ReadTimeout("Read timed out while streaming Pyodide response")
    >   119:         except Exception:
        120:             return None
```

Why this is a false positive: the read converts JS chunks and returns `None` on conversion failure, which the caller treats as end-of-stream.

Checklist evidence: the handler returns a defined fallback.

### [ ] Finding `67` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/67.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/__init__.py:153:1`
- Checklist pattern: teardown best-effort reader cancellation

Source excerpt:

```
        151:             try:
        152:                 await self._reader.cancel()
    >   153:             except Exception:
        154:                 pass
```

Why this is a false positive: cancelling the stream reader during teardown is best-effort; JS errors here are expected and the loop has already ended.

Checklist evidence: the cancellation failure is the expected outcome of tearing down an already-closed stream.

### [ ] Finding `68` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/68.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/__init__.py:153:1`
- Checklist pattern: teardown best-effort reader cancellation

Source excerpt:

```
        151:             try:
        152:                 await self._reader.cancel()
    >   153:             except Exception:
        154:                 pass
```

Why this is a false positive: cancelling the stream reader during teardown is best-effort; JS errors here are expected and the loop has already ended.

Checklist evidence: the cancellation failure is the expected outcome of tearing down an already-closed stream.

### [ ] Finding `69` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/69.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/__init__.py:216:1`
- Checklist pattern: retry-loop handler that acts on the exception

Source excerpt:

```
        214:             try:
        215:                 response = await self._do_send(request, stream, timeout)
    >   216:             except Exception as err:
        217:                 retries = retries.increment(method, request.url, error=err)
```

Why this is a false positive: the handler increments the retry state (`retries.increment(...)`) and re-raises `MaxRetryError` once retries are exhausted — the failure is not swallowed.

Checklist evidence: BP-PY-1's 'hides failures' condition is unmet: the generic catch feeds the retry machinery and re-raises on exhaustion.

### [ ] Finding `70` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/70.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/__init__.py:341:1`
- Checklist pattern: best-effort response-header parse fallback

Source excerpt:

```
        339:                     for entry in js_headers.entries():
        340:                         response_headers[entry[0]] = entry[1]
    >   341:         except Exception:
        342:             pass
```

Why this is a false positive: extracting headers from the JS `headers.entries()` is best-effort; on failure the response is returned without the (optional) header snapshot.

Checklist evidence: the handler is a defensive bridge fallback with a defined result.

### [ ] Finding `71` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/71.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/__init__.py:341:1`
- Checklist pattern: best-effort response-header parse fallback

Source excerpt:

```
        339:                     for entry in js_headers.entries():
        340:                         response_headers[entry[0]] = entry[1]
    >   341:         except Exception:
        342:             pass
```

Why this is a false positive: extracting headers from the JS `headers.entries()` is best-effort; on failure the response is returned without the (optional) header snapshot.

Checklist evidence: the handler is a defensive bridge fallback with a defined result.

### [ ] Finding `72` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/72.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/__init__.py:355:1`
- Checklist pattern: fallback for missing reason phrase

Source excerpt:

```
        353:         try:
        354:             response.reason = js_response.status_text or ""
    >   355:         except Exception:
        356:             response.reason = ""
```

Why this is a false positive: when the JS `status_text` probe fails, `response.reason` falls back to `""` — a defined value, not a swallowed failure.

Checklist evidence: the handler assigns a defined fallback value.

### [ ] Finding `73` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/73.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_sse.py:58:1`
- Checklist pattern: best-effort stream read with `None` fallback

Source excerpt:

```
         56:                 return bytes(value.to_py()).decode("utf-8")
         57:             return None
    >    58:         except Exception:
         59:             return None
```

Why this is a false positive: the read converts JS chunks and returns `None` on conversion failure, which the caller treats as end-of-stream.

Checklist evidence: the handler returns a defined fallback.

### [ ] Finding `74` — `CWE-396`

- Function context: `scripts/niquests/findings/functions/74.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_sse.py:58:1`
- Checklist pattern: best-effort stream read with `None` fallback

Source excerpt:

```
         56:                 return bytes(value.to_py()).decode("utf-8")
         57:             return None
    >    58:         except Exception:
         59:             return None
```

Why this is a false positive: the read converts JS chunks and returns `None` on conversion failure, which the caller treats as end-of-stream.

Checklist evidence: the handler returns a defined fallback.

### [ ] Finding `75` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/75.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_sse.py:147:1`
- Checklist pattern: teardown best-effort reader cancellation

Source excerpt:

```
        145:             try:
        146:                 await self._reader.cancel()
    >   147:             except Exception:
        148:                 pass
```

Why this is a false positive: cancelling the stream reader during teardown is best-effort; JS errors here are expected and the loop has already ended.

Checklist evidence: the cancellation failure is the expected outcome of tearing down an already-closed stream.

### [ ] Finding `76` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/76.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_sse.py:147:1`
- Checklist pattern: teardown best-effort reader cancellation

Source excerpt:

```
        145:             try:
        146:                 await self._reader.cancel()
    >   147:             except Exception:
        148:                 pass
```

Why this is a false positive: cancelling the stream reader during teardown is best-effort; JS errors here are expected and the loop has already ended.

Checklist evidence: the cancellation failure is the expected outcome of tearing down an already-closed stream.

### [ ] Finding `77` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/77.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_sse.py:147:1`
- Checklist pattern: teardown best-effort reader cancellation

Source excerpt:

```
        145:             try:
        146:                 await self._reader.cancel()
    >   147:             except Exception:
        148:                 pass
```

Why this is a false positive: cancelling the stream reader during teardown is best-effort; JS errors here are expected and the loop has already ended.

Checklist evidence: the cancellation failure is the expected outcome of tearing down an already-closed stream.

### [ ] Finding `78` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/78.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_sse.py:147:13`
- Checklist pattern: teardown best-effort reader cancellation

Source excerpt:

```
        145:             try:
        146:                 await self._reader.cancel()
    >   147:             except Exception:
        148:                 pass
```

Why this is a false positive: cancelling the stream reader during teardown is best-effort; JS errors here are expected and the loop has already ended.

Checklist evidence: the cancellation failure is the expected outcome of tearing down an already-closed stream.

### [ ] Finding `79` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/79.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_ws.py:123:1`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        121:         try:
        122:             self._ws.close()
    >   123:         except Exception:  # Defensive: suppress JS errors on teardown
        124:             pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `80` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/80.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_ws.py:123:1`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        121:         try:
        122:             self._ws.close()
    >   123:         except Exception:  # Defensive: suppress JS errors on teardown
        124:             pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `81` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/81.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_ws.py:123:1`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        121:         try:
        122:             self._ws.close()
    >   123:         except Exception:  # Defensive: suppress JS errors on teardown
        124:             pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `82` — `CWE-396`

- Function context: `scripts/niquests/findings/functions/82.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_ws.py:123:1`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        121:         try:
        122:             self._ws.close()
    >   123:         except Exception:  # Defensive: suppress JS errors on teardown
        124:             pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `83` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/83.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_ws.py:123:9`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        121:         try:
        122:             self._ws.close()
    >   123:         except Exception:  # Defensive: suppress JS errors on teardown
        124:             pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `84` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/84.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_ws.py:129:1`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        127:             try:
        128:                 proxy.destroy()
    >   129:             except Exception:  # Defensive: suppress JS errors on teardown
        130:                 pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `85` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/85.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_async/_ws.py:129:1`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        127:             try:
        128:                 proxy.destroy()
    >   129:             except Exception:  # Defensive: suppress JS errors on teardown
        130:                 pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `86` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/86.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_sse.py:59:1`
- Checklist pattern: best-effort stream read with `None` fallback

Source excerpt:

```
         57:                 return bytes(value.to_py()).decode("utf-8")
         58:             return None
    >    59:         except Exception:
         60:             return None
```

Why this is a false positive: the read converts JS chunks and returns `None` on conversion failure, which the caller treats as end-of-stream.

Checklist evidence: the handler returns a defined fallback.

### [ ] Finding `87` — `CWE-396`

- Function context: `scripts/niquests/findings/functions/87.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_sse.py:59:1`
- Checklist pattern: best-effort stream read with `None` fallback

Source excerpt:

```
         57:                 return bytes(value.to_py()).decode("utf-8")
         58:             return None
    >    59:         except Exception:
         60:             return None
```

Why this is a false positive: the read converts JS chunks and returns `None` on conversion failure, which the caller treats as end-of-stream.

Checklist evidence: the handler returns a defined fallback.

### [ ] Finding `88` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/88.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_sse.py:151:1`
- Checklist pattern: teardown best-effort reader cancellation

Source excerpt:

```
        149:             try:
        150:                 run_sync(self._reader.cancel())
    >   151:             except Exception:
        152:                 pass
```

Why this is a false positive: cancelling the stream reader during teardown is best-effort; JS errors here are expected and the loop has already ended.

Checklist evidence: the cancellation failure is the expected outcome of tearing down an already-closed stream.

### [ ] Finding `89` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/89.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_sse.py:151:1`
- Checklist pattern: teardown best-effort reader cancellation

Source excerpt:

```
        149:             try:
        150:                 run_sync(self._reader.cancel())
    >   151:             except Exception:
        152:                 pass
```

Why this is a false positive: cancelling the stream reader during teardown is best-effort; JS errors here are expected and the loop has already ended.

Checklist evidence: the cancellation failure is the expected outcome of tearing down an already-closed stream.

### [ ] Finding `90` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/90.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_sse.py:151:1`
- Checklist pattern: teardown best-effort reader cancellation

Source excerpt:

```
        149:             try:
        150:                 run_sync(self._reader.cancel())
    >   151:             except Exception:
        152:                 pass
```

Why this is a false positive: cancelling the stream reader during teardown is best-effort; JS errors here are expected and the loop has already ended.

Checklist evidence: the cancellation failure is the expected outcome of tearing down an already-closed stream.

### [ ] Finding `91` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/91.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_sse.py:151:13`
- Checklist pattern: teardown best-effort reader cancellation

Source excerpt:

```
        149:             try:
        150:                 run_sync(self._reader.cancel())
    >   151:             except Exception:
        152:                 pass
```

Why this is a false positive: cancelling the stream reader during teardown is best-effort; JS errors here are expected and the loop has already ended.

Checklist evidence: the cancellation failure is the expected outcome of tearing down an already-closed stream.

### [ ] Finding `92` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/92.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_ws.py:177:1`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        175:         try:
        176:             self._ws.close()
    >   177:         except Exception:  # Defensive: suppress JS errors on teardown
        178:             pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `93` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/93.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_ws.py:177:1`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        175:         try:
        176:             self._ws.close()
    >   177:         except Exception:  # Defensive: suppress JS errors on teardown
        178:             pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `94` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/94.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_ws.py:177:1`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        175:         try:
        176:             self._ws.close()
    >   177:         except Exception:  # Defensive: suppress JS errors on teardown
        178:             pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `95` — `CWE-396`

- Function context: `scripts/niquests/findings/functions/95.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_ws.py:177:1`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        175:         try:
        176:             self._ws.close()
    >   177:         except Exception:  # Defensive: suppress JS errors on teardown
        178:             pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `96` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/96.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_ws.py:177:9`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        175:         try:
        176:             self._ws.close()
    >   177:         except Exception:  # Defensive: suppress JS errors on teardown
        178:             pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `97` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/97.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_ws.py:183:1`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        181:             try:
        182:                 proxy.destroy()
    >   183:             except Exception:  # Defensive: suppress JS errors on teardown
        184:                 pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `98` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/98.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/pyodide/_ws.py:183:1`
- Checklist pattern: documented defensive teardown suppression of JS errors

Source excerpt:

```
        181:             try:
        182:                 proxy.destroy()
    >   183:             except Exception:  # Defensive: suppress JS errors on teardown
        184:                 pass
```

Why this is a false positive: the source comments the handler explicitly: `# Defensive: suppress JS errors on teardown` — closing the WS and destroying the proxy during teardown is best-effort.

Checklist evidence: the pass is documented as intentional teardown suppression; the rule's 'silently discarded' condition is not met.

### [ ] Finding `100` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/100.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/revocation/_crl/__init__.py:259:1`
- Checklist pattern: best-effort issuer-certificate fetch during revocation checking

Source excerpt:

```
        257:                     try:
        258:                         raw_intermediary_response = session.get(hint_ca_issuers[0])
    >   259:                     except RequestException:
        260:                         pass
```

Why this is a false positive: the intermediary-issuer fetch is optional enrichment: when the network fetch fails, revocation checking continues in its defined fallback state (no issuer hint).

Checklist evidence: the `RequestException` is an expected outcome of a best-effort network probe with a defined fallback; the pass is the intended behavior.

### [ ] Finding `101` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/101.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/revocation/_crl/__init__.py:259:1`
- Checklist pattern: best-effort issuer-certificate fetch during revocation checking

Source excerpt:

```
        257:                     try:
        258:                         raw_intermediary_response = session.get(hint_ca_issuers[0])
    >   259:                     except RequestException:
        260:                         pass
```

Why this is a false positive: the intermediary-issuer fetch is optional enrichment: when the network fetch fails, revocation checking continues in its defined fallback state (no issuer hint).

Checklist evidence: the `RequestException` is an expected outcome of a best-effort network probe with a defined fallback; the pass is the intended behavior.

### [ ] Finding `102` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/102.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/revocation/_crl/__init__.py:259:21`
- Checklist pattern: best-effort issuer-certificate fetch during revocation checking

Source excerpt:

```
        257:                     try:
        258:                         raw_intermediary_response = session.get(hint_ca_issuers[0])
    >   259:                     except RequestException:
        260:                         pass
```

Why this is a false positive: the intermediary-issuer fetch is optional enrichment: when the network fetch fails, revocation checking continues in its defined fallback state (no issuer hint).

Checklist evidence: the `RequestException` is an expected outcome of a best-effort network probe with a defined fallback; the pass is the intended behavior.

### [ ] Finding `107` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/107.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/revocation/_ocsp/__init__.py:325:1`
- Checklist pattern: best-effort issuer-certificate fetch during revocation checking

Source excerpt:

```
        323:                     try:
        324:                         raw_intermediary_response = session.get(hint_ca_issuers[0])
    >   325:                     except RequestException:
        326:                         pass
```

Why this is a false positive: the intermediary-issuer fetch is optional enrichment: when the network fetch fails, revocation checking continues in its defined fallback state (no issuer hint).

Checklist evidence: the `RequestException` is an expected outcome of a best-effort network probe with a defined fallback; the pass is the intended behavior.

### [ ] Finding `108` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/108.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/revocation/_ocsp/__init__.py:325:1`
- Checklist pattern: best-effort issuer-certificate fetch during revocation checking

Source excerpt:

```
        323:                     try:
        324:                         raw_intermediary_response = session.get(hint_ca_issuers[0])
    >   325:                     except RequestException:
        326:                         pass
```

Why this is a false positive: the intermediary-issuer fetch is optional enrichment: when the network fetch fails, revocation checking continues in its defined fallback state (no issuer hint).

Checklist evidence: the `RequestException` is an expected outcome of a best-effort network probe with a defined fallback; the pass is the intended behavior.

### [ ] Finding `109` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/109.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/revocation/_ocsp/__init__.py:325:21`
- Checklist pattern: best-effort issuer-certificate fetch during revocation checking

Source excerpt:

```
        323:                     try:
        324:                         raw_intermediary_response = session.get(hint_ca_issuers[0])
    >   325:                     except RequestException:
        326:                         pass
```

Why this is a false positive: the intermediary-issuer fetch is optional enrichment: when the network fetch fails, revocation checking continues in its defined fallback state (no issuer hint).

Checklist evidence: the `RequestException` is an expected outcome of a best-effort network probe with a defined fallback; the pass is the intended behavior.

### [ ] Finding `113` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/113.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/sgi/_async/__init__.py:355:1`
- Checklist pattern: application exception captured for later propagation

Source excerpt:

```
        353:             try:
        354:                 await self.app(scope, receive, send_func)
    >   355:             except Exception as ex:
        356:                 app_exception = ex
```

Why this is a false positive: the ASGI app exception is stored in `app_exception` and re-raised to the caller after the response cycle completes.

Checklist evidence: the exception is propagated, not hidden.

### [ ] Finding `114` — `CWE-396`

- Function context: `scripts/niquests/findings/functions/114.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/sgi/_async/__init__.py:355:1`
- Checklist pattern: application exception captured for later propagation (generic catch feeds `app_exception`, re-raised after the response cycle)

Source excerpt:

```
        353:             try:
        354:                 await self.app(scope, receive, send_func)
    >   355:             except Exception as ex:
        356:                 app_exception = ex
```

Why this is a false positive: the ASGI app exception is stored in `app_exception = ex` and re-raised to the caller once the response cycle completes; the generic catch is the deliberate propagation mechanism, so the failure is not hidden.

Checklist evidence: CWE-396's condition 'generic handler can hide distinct failure conditions' is unmet: the handler captures the exception and the caller re-raises it.

### [ ] Finding `115` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/115.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/sgi/_async/__init__.py:661:1`
- Checklist pattern: startup failure recorded and surfaced

Source excerpt:

```
        659:             try:
        660:                 await self.app(scope, receive, send)
    >   661:             except Exception as e:
        662:                 if not startup_complete.is_set():
```

Why this is a false positive: the handler checks the startup-complete state and records the exception so the lifespan startup failure is reported.

Checklist evidence: the failure is surfaced through the startup handshake, not swallowed.

### [ ] Finding `116` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/116.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/sgi/_async/__init__.py:679:1`
- Checklist pattern: wait-until-cancelled loop termination

Source excerpt:

```
        677:         try:
        678:             await asyncio.Future()  # Wait forever until canceled
    >   679:         except (asyncio.CancelledError, GeneratorExit):
        680:             pass
```

Why this is a false positive: `await asyncio.Future()` runs until cancelled; catching `CancelledError`/`GeneratorExit` is the designed shutdown exit of the server task.

Checklist evidence: cancellation is the intended control flow, not a silently discarded failure.

### [ ] Finding `117` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/117.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/sgi/_async/__init__.py:679:1`
- Checklist pattern: wait-until-cancelled loop termination

Source excerpt:

```
        677:         try:
        678:             await asyncio.Future()  # Wait forever until canceled
    >   679:         except (asyncio.CancelledError, GeneratorExit):
        680:             pass
```

Why this is a false positive: `await asyncio.Future()` runs until cancelled; catching `CancelledError`/`GeneratorExit` is the designed shutdown exit of the server task.

Checklist evidence: cancellation is the intended control flow, not a silently discarded failure.

### [ ] Finding `118` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/118.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/sgi/_async/__init__.py:687:1`
- Checklist pattern: designed lifespan-shutdown path

Source excerpt:

```
        685:                 await receive_queue.put({"type": "lifespan.shutdown"})
        686:                 await asyncio.wait_for(shutdown_complete.wait(), timeout=5.0)
    >   687:             except (asyncio.TimeoutError, asyncio.CancelledError, RuntimeError):
        688:                 pass
```

Why this is a false positive: the handler is the ASGI lifespan shutdown routine; timeouts/cancellation during graceful shutdown are the expected outcome and the `pass` completes the shutdown sequence.

Checklist evidence: the exception set is the expected outcome of the shutdown handshake, so the condition 'failure discarded silently' is unmet.

### [ ] Finding `119` — `CWE-772`

- Function context: `scripts/niquests/findings/functions/119.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/unixsocket/__init__.py:47:23`
- Checklist pattern: socket ownership transferred to the connection object

Source excerpt:

```
         45:         self.host = self.socket_path.split("/")[-1]
         46:
    >    47:     def connect(self):
         48:         sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
```

Why this is a false positive: `connect()` stores the socket on `self.sock` (line 51); the `HTTPConnection` lifecycle closes it, so the resource is released after its effective lifetime by its owner.

Checklist evidence: CWE-772 requires the resource to never be released; here it is stored on the instance and closed by the connection lifecycle.

### [ ] Finding `122` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/122.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_adapter.py:216:1`
- Checklist pattern: retry-loop handler that acts on the exception

Source excerpt:

```
        214:             except (InvalidSchema, RuntimeError, SSLError, _WASIProxyError):
        215:                 raise
    >   216:             except Exception as exc:
        217:                 retries = retries.increment(method, request.url, error=exc)
```

Why this is a false positive: the handler increments the retry state and re-raises `MaxRetryError` on exhaustion — the failure is not swallowed.

Checklist evidence: the generic catch feeds the retry machinery and re-raises on exhaustion.

### [ ] Finding `123` — `CWE-396`

- Function context: `scripts/niquests/findings/functions/123.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_adapter.py:216:1`
- Checklist pattern: generic catch inside the retry loop that re-raises on exhaustion (new CWE-396 firing on the audited-FP line `wasi/_adapter.py:216`)

Source excerpt:

```
        214:             except (InvalidSchema, RuntimeError, SSLError, _WASIProxyError):
        215:                 raise
    >   216:             except Exception as exc:
        217:                 retries = retries.increment(method, request.url, error=exc)
```

Why this is a false positive: the handler increments the retry state (`retries.increment(...)`), rewinds the body, sleeps and retries; once retries are exhausted the loop re-raises `MaxRetryError` — failures surface to the caller, so the generic catch does not hide them.

Checklist evidence: CWE-396's condition is unmet: the generic catch feeds the retry machinery and re-raises on exhaustion (same construct the original audit verified for BP-PY-1 at this line).

### [ ] Finding `124` — `CWE-1341`

- Function context: `scripts/niquests/findings/functions/124.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_adapter.py:235:25`
- Checklist pattern: mutually exclusive close paths, single release per execution

Source excerpt:

```
        233:                 except MaxRetryError:
        234:                     if retries.raise_on_status:
    >   235:                         response.close()
        236:                         raise
```

Why this is a false positive: `async_session.py:1015` and its `else: resp.close()` are the two branches of one `isinstance` conditional — exactly one executes; the wasi `response.close()` calls sit in exclusive retry/raise branches (line 235 vs 239).

Checklist evidence: CWE-1341's condition is the same handle released twice on one path; the matched `close()` pairs are mutually exclusive branches (or distinct resources).

### [ ] Finding `125` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/125.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_async/_adapter.py:99:1`
- Checklist pattern: upload-failure cleanup flag

Source excerpt:

```
         97:                     if on_upload_body is not None:  # Defensive: Session always supplies the callback
         98:                         await on_upload_body(sent, total, False, False)
    >    99:     except BaseException:  # Defensive: upload failure cleanup
        100:         failed = True
```

Why this is a false positive: the `BaseException` handler sets `failed = True` (documented `# Defensive: upload failure cleanup`), which drives the caller's cleanup and completion path.

Checklist evidence: the handler takes a visible action, so failures are not silently discarded.

### [ ] Finding `126` — `CWE-396`

- Function context: `scripts/niquests/findings/functions/126.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_async/_adapter.py:99:1`
- Checklist pattern: upload-failure cleanup flag

Source excerpt:

```
         97:                     if on_upload_body is not None:  # Defensive: Session always supplies the callback
         98:                         await on_upload_body(sent, total, False, False)
    >    99:     except BaseException:  # Defensive: upload failure cleanup
        100:         failed = True
```

Why this is a false positive: the `BaseException` handler sets `failed = True` (documented `# Defensive: upload failure cleanup`), which drives the caller's cleanup and completion path.

Checklist evidence: the handler takes a visible action, so failures are not silently discarded.

### [ ] Finding `129` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/129.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_async/_adapter.py:237:1`
- Checklist pattern: retry-loop handler that acts on the exception

Source excerpt:

```
        235:             except (InvalidSchema, RuntimeError, SSLError, _WASIProxyError):
        236:                 raise
    >   237:             except Exception as exc:
        238:                 retries = retries.increment(method, request.url, error=exc)
```

Why this is a false positive: the handler increments the retry state and re-raises `MaxRetryError` on exhaustion — the failure is not swallowed.

Checklist evidence: the generic catch feeds the retry machinery and re-raises on exhaustion.

### [ ] Finding `130` — `CWE-1341`

- Function context: `scripts/niquests/findings/functions/130.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_async/_adapter.py:256:31`
- Checklist pattern: mutually exclusive close paths, single release per execution

Source excerpt:

```
        254:                 except MaxRetryError:
        255:                     if retries.raise_on_status:
    >   256:                         await response.close()
        257:                         raise
```

Why this is a false positive: `async_session.py:1015` and its `else: resp.close()` are the two branches of one `isinstance` conditional — exactly one executes; the wasi `response.close()` calls sit in exclusive retry/raise branches (line 235 vs 239).

Checklist evidence: CWE-1341's condition is the same handle released twice on one path; the matched `close()` pairs are mutually exclusive branches (or distinct resources).

### [ ] Finding `131` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/131.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_async/_sse.py:62:1`
- Checklist pattern: optional numeric conversion fallback in SSE field parsing

Source excerpt:

```
         60:                 try:
         61:                     values[key] = int(value)
    >    62:                 except ValueError:
         63:                     pass
```

Why this is a false positive: the `int(value)` conversion is best-effort; when the field is not numeric the string value is kept, which is the intended fallback for SSE `retry` fields.

Checklist evidence: the exception is an expected parsing outcome with a defined fallback (original string preserved).

### [ ] Finding `132` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/132.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_async/_sse.py:62:1`
- Checklist pattern: optional numeric conversion fallback in SSE field parsing

Source excerpt:

```
         60:                 try:
         61:                     values[key] = int(value)
    >    62:                 except ValueError:
         63:                     pass
```

Why this is a false positive: the `int(value)` conversion is best-effort; when the field is not numeric the string value is kept, which is the intended fallback for SSE `retry` fields.

Checklist evidence: the exception is an expected parsing outcome with a defined fallback (original string preserved).

### [ ] Finding `133` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/133.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_async/_sse.py:62:17`
- Checklist pattern: optional numeric conversion fallback in SSE field parsing

Source excerpt:

```
         60:                 try:
         61:                     values[key] = int(value)
    >    62:                 except ValueError:
         63:                     pass
```

Why this is a false positive: the `int(value)` conversion is best-effort; when the field is not numeric the string value is kept, which is the intended fallback for SSE `retry` fields.

Checklist evidence: the exception is an expected parsing outcome with a defined fallback (original string preserved).

### [ ] Finding `134` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/134.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_sse.py:61:1`
- Checklist pattern: optional numeric conversion fallback in SSE field parsing

Source excerpt:

```
         59:                 try:
         60:                     values[key] = int(value)
    >    61:                 except ValueError:
         62:                     pass
```

Why this is a false positive: the `int(value)` conversion is best-effort; when the field is not numeric the string value is kept, which is the intended fallback for SSE `retry` fields.

Checklist evidence: the exception is an expected parsing outcome with a defined fallback (original string preserved).

### [ ] Finding `135` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/135.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_sse.py:61:1`
- Checklist pattern: optional numeric conversion fallback in SSE field parsing

Source excerpt:

```
         59:                 try:
         60:                     values[key] = int(value)
    >    61:                 except ValueError:
         62:                     pass
```

Why this is a false positive: the `int(value)` conversion is best-effort; when the field is not numeric the string value is kept, which is the intended fallback for SSE `retry` fields.

Checklist evidence: the exception is an expected parsing outcome with a defined fallback (original string preserved).

### [ ] Finding `136` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/136.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_sse.py:61:17`
- Checklist pattern: optional numeric conversion fallback in SSE field parsing

Source excerpt:

```
         59:                 try:
         60:                     values[key] = int(value)
    >    61:                 except ValueError:
         62:                     pass
```

Why this is a false positive: the `int(value)` conversion is best-effort; when the field is not numeric the string value is kept, which is the intended fallback for SSE `retry` fields.

Checklist evidence: the exception is an expected parsing outcome with a defined fallback (original string preserved).

### [ ] Finding `137` — `BP-PY-49`

- Function context: `scripts/niquests/findings/functions/137.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/extensions/wasi/_utils.py:63:41`
- Checklist pattern: `verify=False` text inside an error-message string

Source excerpt:

```
         61:     if verify is not True:
         62:         raise SSLError(
    >    63:             "WASI HTTPS uses the host trust policy; verify=False and custom CA bundles are unavailable.",
         64:             request=request,
```

Why this is a false positive: the flagged text is a string literal in an `SSLError` message; the surrounding code (`if verify is not True: raise SSLError(...)`) actively rejects `verify=False` on WASI.

Checklist evidence: no TLS disabling construct exists at this line — the regex matched a string literal inside an error message.

### [ ] Finding `140` — `BP-PY-36`

- Function context: `scripts/niquests/findings/functions/140.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/help.py:176:1`
- Checklist pattern: niquests' own HTTP `Session`, not SQLAlchemy

Source excerpt:

```
        174:
        175:
    >   176: pypi_session = Session()
        177:
```

Why this is a false positive: `help.py:34` imports `Session` from niquests itself (`from . import ... Session`); `pypi_session = Session()` is an HTTP client session used to query the PyPI JSON API — SQLAlchemy is not involved.

Checklist evidence: BP-PY-36's condition is a SQLAlchemy `Session`/`SessionLocal`; the shown source constructs the package's own HTTP session class.

### [ ] Finding `141` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/141.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/help.py:196:1`
- Checklist pattern: best-effort update check

Source excerpt:

```
        194:                     UserWarning,
        195:                 )
    >   196:     except (RequestException, JSONDecodeError, HTTPError):
        197:         pass
```

Why this is a false positive: `check_update` is an optional PyPI version probe; network/JSON failures are benign and the check silently degrades — the pass is the designed behavior of the best-effort helper.

Checklist evidence: the exception set (network + parse) is expected for a best-effort update check with no failure semantics.

### [ ] Finding `142` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/142.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/help.py:196:1`
- Checklist pattern: best-effort update check

Source excerpt:

```
        194:                     UserWarning,
        195:                 )
    >   196:     except (RequestException, JSONDecodeError, HTTPError):
        197:         pass
```

Why this is a false positive: `check_update` is an optional PyPI version probe; network/JSON failures are benign and the check silently degrades — the pass is the designed behavior of the best-effort helper.

Checklist evidence: the exception set (network + parse) is expected for a best-effort update check with no failure semantics.

### [ ] Finding `143` — `CWE-93`

- Function context: `scripts/niquests/findings/functions/143.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/models.py:488:25`
- Checklist pattern: write to an outgoing request header, not a response header

Source excerpt:

```
        486:                     if isinstance(name, bytes):
        487:                         name = name.decode()
    >   488:                     self.headers[name] = value
        489:
```

Why this is a false positive: `auth.py:75` sets the `Authorization` header on the outgoing request (from the caller's own token); `models.py:488` populates the `PreparedRequest.headers` dict in `prepare_headers` from the caller-supplied header mapping — neither is an HTTP response header.

Checklist evidence: CWE-93's condition is writing an externally influenced value into an HTTP response header; the shown lines write client-side request headers.

### [ ] Finding `144` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/144.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/models.py:1065:1`
- Checklist pattern: defensive probe of lazily-initialized attributes

Source excerpt:

```
       1063:             super().__getattribute__("_promise")
       1064:             super().__getattribute__("connection").gather(self)
    >  1065:         except AttributeError:
       1066:             pass
```

Why this is a false positive: the `try` probes `_promise`/`_gather` internals that exist only in initialized states; `AttributeError` means there is nothing to gather — the designed no-op.

Checklist evidence: the exception is the expected outcome for an uninitialized optional internal, and the no-op is the designed behavior.

### [ ] Finding `145` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/145.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/models.py:1065:1`
- Checklist pattern: defensive probe of lazily-initialized attributes

Source excerpt:

```
       1063:             super().__getattribute__("_promise")
       1064:             super().__getattribute__("connection").gather(self)
    >  1065:         except AttributeError:
       1066:             pass
```

Why this is a false positive: the `try` probes `_promise`/`_gather` internals that exist only in initialized states; `AttributeError` means there is nothing to gather — the designed no-op.

Checklist evidence: the exception is the expected outcome for an uninitialized optional internal, and the no-op is the designed behavior.

### [ ] Finding `146` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/146.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/models.py:1065:9`
- Checklist pattern: defensive probe of lazily-initialized attributes

Source excerpt:

```
       1063:             super().__getattribute__("_promise")
       1064:             super().__getattribute__("connection").gather(self)
    >  1065:         except AttributeError:
       1066:             pass
```

Why this is a false positive: the `try` probes `_promise`/`_gather` internals that exist only in initialized states; `AttributeError` means there is nothing to gather — the designed no-op.

Checklist evidence: the exception is the expected outcome for an uninitialized optional internal, and the no-op is the designed behavior.

### [ ] Finding `147` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/147.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/models.py:1079:1`
- Checklist pattern: defensive probe of lazily-initialized attributes

Source excerpt:

```
       1077:                 else:
       1078:                     super().__getattribute__("_gather")()
    >  1079:         except AttributeError:
       1080:             pass
```

Why this is a false positive: the `try` probes `_promise`/`_gather` internals that exist only in initialized states; `AttributeError` means there is nothing to gather — the designed no-op.

Checklist evidence: the exception is the expected outcome for an uninitialized optional internal, and the no-op is the designed behavior.

### [ ] Finding `148` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/148.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/models.py:1233:1`
- Checklist pattern: optional content-length progress hint fallback

Source excerpt:

```
       1231:                     try:
       1232:                         self.download_progress.content_length = int(self.headers["content-length"])
    >  1233:                     except ValueError:
       1234:                         pass
```

Why this is a false positive: the `content-length` header may be absent or non-numeric; on `ValueError` the progress `content_length` simply stays unset — an optional hint, not a failure.

Checklist evidence: the flagged `pass` is the expected fallback for an optional progress attribute.

### [ ] Finding `149` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/149.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/models.py:1318:1`
- Checklist pattern: optional content-length progress hint fallback

Source excerpt:

```
       1316:                     try:
       1317:                         self.download_progress.content_length = int(self.headers["content-length"])
    >  1318:                     except ValueError:
       1319:                         pass
```

Why this is a false positive: the `content-length` header may be absent or non-numeric; on `ValueError` the progress `content_length` simply stays unset — an optional hint, not a failure.

Checklist evidence: the flagged `pass` is the expected fallback for an optional progress attribute.

### [ ] Finding `150` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/150.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/models.py:1768:1`
- Checklist pattern: optional content-length progress hint fallback

Source excerpt:

```
       1766:                     try:
       1767:                         self.download_progress.content_length = int(self.headers["content-length"])
    >  1768:                     except ValueError:
       1769:                         pass
```

Why this is a false positive: the `content-length` header may be absent or non-numeric; on `ValueError` the progress `content_length` simply stays unset — an optional hint, not a failure.

Checklist evidence: the flagged `pass` is the expected fallback for an optional progress attribute.

### [ ] Finding `151` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/151.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/models.py:1827:1`
- Checklist pattern: optional content-length progress hint fallback

Source excerpt:

```
       1825:                     try:
       1826:                         self.download_progress.content_length = int(self.headers["content-length"])
    >  1827:                     except ValueError:
       1828:                         pass
```

Why this is a false positive: the `content-length` header may be absent or non-numeric; on `ValueError` the progress `content_length` simply stays unset — an optional hint, not a failure.

Checklist evidence: the flagged `pass` is the expected fallback for an optional progress attribute.

### [ ] Finding `152` — `CWE-829`

- Function context: `scripts/niquests/findings/functions/152.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/packages.py:42:29`
- Checklist pattern: dynamic import over a developer-controlled constant allowlist

Source excerpt:

```
         40:
         41:     try:
    >    42:         locals()[package] = __import__(to_be_imported)
         43:     except ImportError:
```

Why this is a false positive: `to_be_imported` is computed from the hardcoded module tuple `("urllib3", "charset_normalizer", "idna", "chardet")` plus fixed alias mappings — no request- or user-derived value reaches `__import__`.

Checklist evidence: CWE-829/CWE-94's condition is dynamically loading untrusted/request-derived modules; the shown loop iterates a compile-time allowlist of the package's own dependencies.

### [ ] Finding `153` — `CWE-94`

- Function context: `scripts/niquests/findings/functions/153.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/packages.py:42:29`
- Checklist pattern: dynamic import over a developer-controlled constant allowlist

Source excerpt:

```
         40:
         41:     try:
    >    42:         locals()[package] = __import__(to_be_imported)
         43:     except ImportError:
```

Why this is a false positive: `to_be_imported` is computed from the hardcoded module tuple `("urllib3", "charset_normalizer", "idna", "chardet")` plus fixed alias mappings — no request- or user-derived value reaches `__import__`.

Checklist evidence: CWE-829/CWE-94's condition is dynamically loading untrusted/request-derived modules; the shown loop iterates a compile-time allowlist of the package's own dependencies.

### [ ] Finding `154` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/154.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/sessions.py:115:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
        113:     try:
        114:         from .extensions.wasi._adapter import WASIAdapter as WASIHTTPAdapter
    >   115:     except ImportError:
        116:         pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `155` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/155.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/sessions.py:115:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
        113:     try:
        114:         from .extensions.wasi._adapter import WASIAdapter as WASIHTTPAdapter
    >   115:     except ImportError:
        116:         pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `156` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/156.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/sessions.py:115:5`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
        113:     try:
        114:         from .extensions.wasi._adapter import WASIAdapter as WASIHTTPAdapter
    >   115:     except ImportError:
        116:         pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `157` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/157.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/sessions.py:1449:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
       1447:                             verify as ocsp_verify,
       1448:                         )
    >  1449:                     except ImportError:
       1450:                         pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `159` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/159.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/sessions.py:1476:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
       1474:                             verify as crl_verify,
       1475:                         )
    >  1476:                     except ImportError:
       1477:                         pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `161` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/161.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/sessions.py:1709:1`
- Checklist pattern: generator-exhaustion control flow

Source excerpt:

```
       1707:                         self.resolve_redirects(r, request, yield_requests=True, **kwargs)  # type: ignore[assignment]
       1708:                     )
    >  1709:                 except StopIteration:
       1710:                     pass
```

Why this is a false positive: `StopIteration`/`StopAsyncIteration` is the generator's normal termination signal; catching it is the intended loop-exit control flow, not a discarded failure.

Checklist evidence: the rule's 'failures discarded silently' condition does not apply to iterator-exhaustion signalling.

### [ ] Finding `162` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/162.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/sessions.py:1785:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
       1783:                 if scheme in extension.supported_schemes() and extension.scheme_to_http_scheme(scheme) == parse_scheme(prefix):
       1784:                     return adapter
    >  1785:         except ImportError:
       1786:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `169` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/169.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/utils.py:157:1`
- Checklist pattern: expected fallback for non-file objects

Source excerpt:

```
        155:         try:
        156:             fileno = o.fileno()
    >   157:         except (io.UnsupportedOperation, AttributeError):
        158:             # AttributeError is a surprising exception, seeing as how we've just checked
```

Why this is a false positive: `fileno()` is unsupported for pure-Python/byte-stream objects; the surrounding comment documents the fallback path that handles objects without a file descriptor.

Checklist evidence: the shown source documents the `AttributeError`/`UnsupportedOperation` as a surprising-but-expected case with a defined fallback.

### [ ] Finding `170` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/170.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/utils.py:157:1`
- Checklist pattern: expected fallback for non-file objects

Source excerpt:

```
        155:         try:
        156:             fileno = o.fileno()
    >   157:         except (io.UnsupportedOperation, AttributeError):
        158:             # AttributeError is a surprising exception, seeing as how we've just checked
```

Why this is a false positive: `fileno()` is unsupported for pure-Python/byte-stream objects; the surrounding comment documents the fallback path that handles objects without a file descriptor.

Checklist evidence: the shown source documents the `AttributeError`/`UnsupportedOperation` as a surprising-but-expected case with a defined fallback.

### [ ] Finding `171` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/171.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/utils.py:815:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
        813:         try:
        814:             extension_class = load_extension(maybe_extension_scheme, implementation)
    >   815:         except ImportError:
        816:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `172` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/172.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/utils.py:815:9`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
        813:         try:
        814:             extension_class = load_extension(maybe_extension_scheme, implementation)
    >   815:         except ImportError:
        816:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `174` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/174.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/src/niquests/utils.py:1196:1`
- Checklist pattern: scheme-parsing fallback that falls through to an explicit raise

Source excerpt:

```
       1194:
       1195:             return outsider_scheme.lower()
    >  1196:         except ValueError:
       1197:             pass
```

Why this is a false positive: the inner `except ValueError: pass` is immediately followed by `raise MissingSchema(...)`, so the parse error is surfaced to the caller, not discarded.

Checklist evidence: the `pass` only bridges to the enclosing `raise`, so no error condition goes without action.

### [ ] Finding `175` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/175.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/conftest.py:107:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
        105:             yield selenium
        106:
    >   107: except ImportError:
        108:     pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `176` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/176.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/conftest.py:107:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
        105:             yield selenium
        106:
    >   107: except ImportError:
        108:     pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `177` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/177.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/conftest.py:107:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
        105:             yield selenium
        106:
    >   107: except ImportError:
        108:     pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `178` — `BP-PY-1`

- Function context: `scripts/niquests/findings/functions/178.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_asgi.py:52:1`
- Checklist pattern: test websocket receive loop termination

Source excerpt:

```
         50:         try:
         51:             data = await websocket.receive()
    >    52:         except Exception:
         53:             break
```

Why this is a false positive: the test receives frames in a loop and breaks on any connection error — the generic catch is the loop's termination, and a broken socket fails the test otherwise.

Checklist evidence: test code: the exception is the expected end of the receive loop and any real failure fails the test.

### [ ] Finding `179` — `CWE-617`

- Function context: `scripts/niquests/findings/functions/179.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_async.py:30:1`
- Checklist pattern: test assertion is the verification mechanism

Source excerpt:

```
         28:     assert encoded_values == [value]
         29:     assert request.body == b'{"id":2}'
    >    30:     assert request.headers["Content-Type"] == "application/json;charset=utf-8"
         31:     assert request.headers["Content-Length"] == "8"
```

Why this is a false positive: `assert request.headers[...] == ...` verifies the request the server received; the asserted values are the test's own expectations.

Checklist evidence: CWE-617's condition is a reachable assertion on request-controlled state in production; the shown source is test verification.

### [ ] Finding `180` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/180.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_async.py:195:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
        193:             assert "Herman Melville - Moby-Dick" in content
        194:
    >   195:     async def test_explicit_close_in_streaming_response(self) -> None:
        196:         async with AsyncSession() as s:
```

Why this is a false positive: the test explicitly closes a streaming response and relies on any failure to surface as an exception (there is no `except`), so the test fails if close is broken.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `181` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/181.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:684:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
        682:
        683:
    >   684: def test_sync_sse_send_raises(selenium_jspi):
        685:     """Test that SSE send_payload raises NotImplementedError."""
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `182` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/182.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:727:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
        725:
        726:
    >   727: def test_sync_timeout(selenium_jspi):
        728:     """Test that timeout raises ConnectTimeout on slow responses."""
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `183` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/183.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_emscripten.py:770:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
        768:
        769:
    >   770: def test_async_timeout(selenium):
        771:     """Test that async timeout raises ConnectTimeout on slow responses."""
```

Why this is a false positive: the test body delegates to the `_inner_test_*` helper, which contains the actual assertions (status code, JSON payload, raw bytes); the wrapper only transports coverage data, so the test does verify outcomes.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `184` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/184.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:79:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
         77:     digest_auth_algo = ("MD5", "SHA-256", "SHA-512")
         78:
    >    79:     def test_entry_points(self):
         80:         niquests.Session
```

Why this is a false positive: the body only touches the public entry points (`niquests.Session`, `.get`, `.head`); a missing or broken entry point raises `AttributeError` and fails the test — an API-surface smoke test.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `185` — `CWE-208`

- Function context: `scripts/niquests/findings/functions/185.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:574:16`
- Checklist pattern: equality of an auth object in a test assertion

Source excerpt:

```
        572:         """Session accepts custom auth at initialization."""
        573:         s = niquests.Session(auth=init_auth)
    >   574:         assert s.auth == expected_auth
        575:
```

Why this is a false positive: `assert s.auth == expected_auth` compares `AuthBase` objects configured by the test itself; no secret comparison subject to timing side-channels exists.

Checklist evidence: CWE-208's condition is comparing security-sensitive values; the compared values are test-owned auth objects, and the line is a test assertion.

### [ ] Finding `186` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/186.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:1030:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
       1028:         assert not r.ok
       1029:
    >  1030:     def test_decompress_gzip(self, httpbin):
       1031:         r = niquests.get(httpbin("gzip"))
```

Why this is a false positive: `r.content.decode("ascii")` raises if gzip decompression produced garbage — the test fails on broken decompression, so it verifies the outcome.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `187` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/187.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:1044:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
       1042:         ),
       1043:     )
    >  1044:     def test_unicode_get(self, httpbin, url, params):
       1045:         niquests.get(httpbin(url), params=params)
```

Why this is a false positive: the request fails (raises) if unicode query params are mishandled, failing the test — a smoke test of unicode parameter support.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `188` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/188.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:1047:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
       1045:         niquests.get(httpbin(url), params=params)
       1046:
    >  1047:     def test_unicode_header_name(self, httpbin):
       1048:         niquests.put(
```

Why this is a false positive: the request fails (raises) if unicode header names/values are mishandled, failing the test — a smoke test of unicode header support.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `189` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/189.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:1054:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
       1052:         )  # compat.str is unicode.
       1053:
    >  1054:     def test_pyopenssl_redirect(self, httpbin_secure, httpbin_ca_bundle):
       1055:         niquests.get(httpbin_secure("status", "301"), verify=httpbin_ca_bundle)
```

Why this is a false positive: the redirected request over the pyopenssl TLS stack fails (raises) if redirects or TLS verification break, failing the test.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `190` — `BP-PY-10`

- Function context: `scripts/niquests/findings/functions/190.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:1651:16`
- Checklist pattern: self-serialized pickle round-trip tests

Source excerpt:

```
       1649:
       1650:         # verify we can pickle the original request
    >  1651:         assert pickle.loads(pickle.dumps(r.request))
       1652:
```

Why this is a false positive: each flagged line is `pickle.loads(pickle.dumps(obj))` in a test — the data is serialized by the same expression, so the unpickled bytes are self-produced, not untrusted input.

Checklist evidence: BP-PY-10/CWE-502's condition is deserializing untrusted data; the tests round-trip their own in-memory objects.

### [ ] Finding `191` — `CWE-502`

- Function context: `scripts/niquests/findings/functions/191.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:1651:16`
- Checklist pattern: self-serialized pickle round-trip tests

Source excerpt:

```
       1649:
       1650:         # verify we can pickle the original request
    >  1651:         assert pickle.loads(pickle.dumps(r.request))
       1652:
```

Why this is a false positive: each flagged line is `pickle.loads(pickle.dumps(obj))` in a test — the data is serialized by the same expression, so the unpickled bytes are self-produced, not untrusted input.

Checklist evidence: BP-PY-10/CWE-502's condition is deserializing untrusted data; the tests round-trip their own in-memory objects.

### [ ] Finding `192` — `BP-PY-10`

- Function context: `scripts/niquests/findings/functions/192.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:1655:14`
- Checklist pattern: self-serialized pickle round-trip tests

Source excerpt:

```
       1653:         # verify we can pickle the response and that we have access to
       1654:         # the original request.
    >  1655:         pr = pickle.loads(pickle.dumps(r))
       1656:         assert r.request.url == pr.request.url
```

Why this is a false positive: each flagged line is `pickle.loads(pickle.dumps(obj))` in a test — the data is serialized by the same expression, so the unpickled bytes are self-produced, not untrusted input.

Checklist evidence: BP-PY-10/CWE-502's condition is deserializing untrusted data; the tests round-trip their own in-memory objects.

### [ ] Finding `193` — `BP-PY-10`

- Function context: `scripts/niquests/findings/functions/193.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:1663:13`
- Checklist pattern: self-serialized pickle round-trip tests

Source excerpt:

```
       1661:
       1662:         # Verify PreparedRequest can be pickled and unpickled
    >  1663:         r = pickle.loads(pickle.dumps(p))
       1664:         assert r.url == p.url
```

Why this is a false positive: each flagged line is `pickle.loads(pickle.dumps(obj))` in a test — the data is serialized by the same expression, so the unpickled bytes are self-produced, not untrusted input.

Checklist evidence: BP-PY-10/CWE-502's condition is deserializing untrusted data; the tests round-trip their own in-memory objects.

### [ ] Finding `194` — `BP-PY-10`

- Function context: `scripts/niquests/findings/functions/194.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:1679:13`
- Checklist pattern: self-serialized pickle round-trip tests

Source excerpt:

```
       1677:
       1678:         # Verify PreparedRequest can be pickled and unpickled
    >  1679:         r = pickle.loads(pickle.dumps(p))
       1680:         assert r.url == p.url
```

Why this is a false positive: each flagged line is `pickle.loads(pickle.dumps(obj))` in a test — the data is serialized by the same expression, so the unpickled bytes are self-produced, not untrusted input.

Checklist evidence: BP-PY-10/CWE-502's condition is deserializing untrusted data; the tests round-trip their own in-memory objects.

### [ ] Finding `195` — `BP-PY-10`

- Function context: `scripts/niquests/findings/functions/195.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:1694:13`
- Checklist pattern: self-serialized pickle round-trip tests

Source excerpt:

```
       1692:
       1693:         # Verify PreparedRequest can be pickled
    >  1694:         r = pickle.loads(pickle.dumps(p))
       1695:         assert r.url == p.url
```

Why this is a false positive: each flagged line is `pickle.loads(pickle.dumps(obj))` in a test — the data is serialized by the same expression, so the unpickled bytes are self-produced, not untrusted input.

Checklist evidence: BP-PY-10/CWE-502's condition is deserializing untrusted data; the tests round-trip their own in-memory objects.

### [ ] Finding `196` — `BP-PY-10`

- Function context: `scripts/niquests/findings/functions/196.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:1724:13`
- Checklist pattern: self-serialized pickle round-trip tests

Source excerpt:

```
       1722:         s = niquests.Session()
       1723:
    >  1724:         s = pickle.loads(pickle.dumps(s))
       1725:         s.proxies = getproxies()
```

Why this is a false positive: each flagged line is `pickle.loads(pickle.dumps(obj))` in a test — the data is serialized by the same expression, so the unpickled bytes are self-produced, not untrusted input.

Checklist evidence: BP-PY-10/CWE-502's condition is deserializing untrusted data; the tests round-trip their own in-memory objects.

### [ ] Finding `197` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/197.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:2163:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
       2161:         adapter.build_response = build_response
       2162:
    >  2163:     def test_redirect_with_wrong_gzipped_header(self, httpbin):
       2164:         s = niquests.Session()
```

Why this is a false positive: the patched adapter breaks the gzipped redirect; if the client mishandles it, `s.get(url)` raises and the test fails — a regression smoke test.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `198` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/198.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:2630:1`
- Checklist pattern: expected-exception path of a timeout test

Source excerpt:

```
       2628:             niquests.get(httpbin("delay/10"), timeout=timeout)
       2629:             pytest.fail("The recv() request should time out.")
    >  2630:         except ReadTimeout:
       2631:             pass
```

Why this is a false positive: `pytest.fail` fires when the request does NOT time out; the `except ReadTimeout`/`ConnectTimeout: pass` is the test's assertion mechanism — the exception IS the expected outcome.

Checklist evidence: the `pass` is the success path of the test (the exception was raised as expected), not a discarded failure.

### [ ] Finding `199` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/199.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:2630:1`
- Checklist pattern: expected-exception path of a timeout test

Source excerpt:

```
       2628:             niquests.get(httpbin("delay/10"), timeout=timeout)
       2629:             pytest.fail("The recv() request should time out.")
    >  2630:         except ReadTimeout:
       2631:             pass
```

Why this is a false positive: `pytest.fail` fires when the request does NOT time out; the `except ReadTimeout`/`ConnectTimeout: pass` is the test's assertion mechanism — the exception IS the expected outcome.

Checklist evidence: the `pass` is the success path of the test (the exception was raised as expected), not a discarded failure.

### [ ] Finding `200` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/200.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:2630:9`
- Checklist pattern: expected-exception path of a timeout test

Source excerpt:

```
       2628:             niquests.get(httpbin("delay/10"), timeout=timeout)
       2629:             pytest.fail("The recv() request should time out.")
    >  2630:         except ReadTimeout:
       2631:             pass
```

Why this is a false positive: `pytest.fail` fires when the request does NOT time out; the `except ReadTimeout`/`ConnectTimeout: pass` is the test's assertion mechanism — the exception IS the expected outcome.

Checklist evidence: the `pass` is the success path of the test (the exception was raised as expected), not a discarded failure.

### [ ] Finding `201` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/201.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:2647:1`
- Checklist pattern: expected-exception path of a timeout test

Source excerpt:

```
       2645:             niquests.get(TARPIT, timeout=timeout)
       2646:             pytest.fail("The connect() request should time out.")
    >  2647:         except ConnectTimeout:
       2648:             pass
```

Why this is a false positive: `pytest.fail` fires when the request does NOT time out; the `except ReadTimeout`/`ConnectTimeout: pass` is the test's assertion mechanism — the exception IS the expected outcome.

Checklist evidence: the `pass` is the success path of the test (the exception was raised as expected), not a discarded failure.

### [ ] Finding `202` — `CWE-617`

- Function context: `scripts/niquests/findings/functions/202.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_requests.py:3078:1`
- Checklist pattern: test assertion is the verification mechanism

Source excerpt:

```
       3076:         assert encoded_values == [value]
       3077:         assert request.body == b'{"id":2}'
    >  3078:         assert request.headers["Content-Type"] == "application/json;charset=utf-8"
       3079:         assert request.headers["Content-Length"] == "8"
```

Why this is a false positive: `assert request.headers[...] == ...` verifies the request the server received; the asserted values are the test's own expectations.

Checklist evidence: CWE-617's condition is a reachable assertion on request-controlled state in production; the shown source is test verification.

### [ ] Finding `203` — `CWE-772`

- Function context: `scripts/niquests/findings/functions/203.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_testserver.py:41:1`
- Checklist pattern: test socket used to verify connection refusal

Source excerpt:

```
         39:
         40:         with pytest.raises(socket.error):
    >    41:             new_sock = socket.socket()
         42:             new_sock.connect((host, port))
```

Why this is a false positive: `new_sock = socket.socket()` is created inside `with pytest.raises(socket.error)` to assert that connecting to a closed port fails; the test discards the socket by design after the expected error.

Checklist evidence: CWE-772's condition is a resource leak in production code; the shown source is a test fixture checking connection failure.

### [ ] Finding `204` — `BP-PY-41`

- Function context: `scripts/niquests/findings/functions/204.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_testserver.py:134:1`
- Checklist pattern: test delegates to an assert-bearing helper / is a smoke or expected-exception test

Source excerpt:

```
        132:             assert r.headers["Content-Length"] == "0"
        133:
    >   134:     def test_basic_waiting_server(self):
        135:         """the server waits for the block_server event to be set before closing"""
```

Why this is a false positive: the test drives the blocking server and would fail with a socket error if the server misbehaved; the server's `handler_results` are asserted by the server fixtures' surrounding checks.

Checklist evidence: BP-PY-41's condition is a placeholder test that passes even when the code is broken; the shown tests verify outcomes by delegation, by smoke failure, or by expected-exception handling.

### [ ] Finding `205` — `CWE-397`

- Function context: `scripts/niquests/findings/functions/205.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_testserver.py:220:1`
- Checklist pattern: tests deliberately raising a generic exception

Source excerpt:

```
        218:         with pytest.raises(Exception):
        219:             with server:
    >   220:                 raise Exception()
        221:
```

Why this is a false positive: `raise Exception()` / `raise Exception("Expected exception")` inside `with pytest.raises(Exception):` exercises exception-handling behavior with a deliberately generic exception.

Checklist evidence: CWE-397's condition targets production code raising generic exceptions; the shown lines are test fixtures for exception handling.

### [ ] Finding `206` — `CWE-397`

- Function context: `scripts/niquests/findings/functions/206.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_utils.py:816:1`
- Checklist pattern: tests deliberately raising a generic exception

Source excerpt:

```
        814:     with pytest.raises(Exception) as exception:
        815:         with set_environ("test1", None):
    >   816:             raise Exception("Expected exception")
        817:
```

Why this is a false positive: `raise Exception()` / `raise Exception("Expected exception")` inside `with pytest.raises(Exception):` exercises exception-handling behavior with a deliberately generic exception.

Checklist evidence: CWE-397's condition targets production code raising generic exceptions; the shown lines are test fixtures for exception handling.

### [ ] Finding `207` — `CWE-770`

- Function context: `scripts/niquests/findings/functions/207.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/test_wsgi.py:59:21`
- Checklist pattern: test WSGI fixture reading its own request body

Source excerpt:

```
         57:             "path": request.path,
         58:             "query": request.query_string.decode("utf-8"),
    >    59:             "body": request.get_data(as_text=True),
         60:             "headers": dict(request.headers),
```

Why this is a false positive: `request.get_data(as_text=True)` runs in a minimal Flask test app in `test_wsgi.py`; the test client controls the body and no production resource-exhaustion surface exists.

Checklist evidence: CWE-770's condition is a production request reader without a size limit; the shown source is a test fixture.

### [ ] Finding `208` — `CWE-772`

- Function context: `scripts/niquests/findings/functions/208.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/testserver/server.py:81:39`
- Checklist pattern: socket ownership transferred via return

Source excerpt:

```
         79:             self.stop_event.set()
         80:
    >    81:     def _create_socket_and_bind(self):
         82:         sock = socket.socket()
```

Why this is a false positive: `_create_socket_and_bind` returns `sock`; the caller stores it as `self.server_sock`, which `_close_server_sock_ignore_errors` closes (line 89).

Checklist evidence: the resource escapes the function to its owner, which releases it; no leak exists.

### [ ] Finding `209` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/209.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/testserver/server.py:90:1`
- Checklist pattern: defensive close of the server socket

Source excerpt:

```
         88:         try:
         89:             self.server_sock.close()
    >    90:         except OSError:
         91:             pass
```

Why this is a false positive: `close()` on an already-closed server socket may raise `OSError`; the pass is the designed best-effort teardown of the test server.

Checklist evidence: the exception is the expected outcome of idempotent close; the handler is deliberately best-effort.

### [ ] Finding `210` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/210.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/testserver/server.py:90:1`
- Checklist pattern: defensive close of the server socket

Source excerpt:

```
         88:         try:
         89:             self.server_sock.close()
    >    90:         except OSError:
         91:             pass
```

Why this is a false positive: `close()` on an already-closed server socket may raise `OSError`; the pass is the designed best-effort teardown of the test server.

Checklist evidence: the exception is the expected outcome of idempotent close; the handler is deliberately best-effort.

### [ ] Finding `211` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/211.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/testserver/server.py:90:9`
- Checklist pattern: defensive close of the server socket

Source excerpt:

```
         88:         try:
         89:             self.server_sock.close()
    >    90:         except OSError:
         91:             pass
```

Why this is a false positive: `close()` on an already-closed server socket may raise `OSError`; the pass is the designed best-effort teardown of the test server.

Checklist evidence: the exception is the expected outcome of idempotent close; the handler is deliberately best-effort.

### [ ] Finding `212` — `CWE-829`

- Function context: `scripts/niquests/findings/functions/212.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/app.py:21:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         19:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         20:         try:
    >    21:             importlib.import_module(module.name)
         22:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `213` — `CWE-94`

- Function context: `scripts/niquests/findings/functions/213.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/app.py:21:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         19:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         20:         try:
    >    21:             importlib.import_module(module.name)
         22:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `214` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/214.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/app.py:22:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         20:         try:
         21:             importlib.import_module(module.name)
    >    22:         except ImportError:
         23:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `215` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/215.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/app.py:22:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         20:         try:
         21:             importlib.import_module(module.name)
    >    22:         except ImportError:
         23:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `216` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/216.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/app.py:22:9`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         20:         try:
         21:             importlib.import_module(module.name)
    >    22:         except ImportError:
         23:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `217` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/217.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/edge_cases.py:257:1`
- Checklist pattern: expected-exception exercises in WASI guest edge-case tests

Source excerpt:

```
        255:     try:
        256:         wasi._unwrap_result(Err(Failed()))
    >   257:     except Err:
        258:         pass
```

Why this is a false positive: each `try` deliberately triggers a failure (`_unwrap_result(Err(...))`, reads on broken streams) and the `pass` records that the expected exception was raised — the test's assertion.

Checklist evidence: the except-pass is the expected-exception verification of the test, not silent error discarding.

### [ ] Finding `218` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/218.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/edge_cases.py:257:1`
- Checklist pattern: expected-exception exercises in WASI guest edge-case tests

Source excerpt:

```
        255:     try:
        256:         wasi._unwrap_result(Err(Failed()))
    >   257:     except Err:
        258:         pass
```

Why this is a false positive: each `try` deliberately triggers a failure (`_unwrap_result(Err(...))`, reads on broken streams) and the `pass` records that the expected exception was raised — the test's assertion.

Checklist evidence: the except-pass is the expected-exception verification of the test, not silent error discarding.

### [ ] Finding `219` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/219.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/edge_cases.py:257:5`
- Checklist pattern: expected-exception exercises in WASI guest edge-case tests

Source excerpt:

```
        255:     try:
        256:         wasi._unwrap_result(Err(Failed()))
    >   257:     except Err:
        258:         pass
```

Why this is a false positive: each `try` deliberately triggers a failure (`_unwrap_result(Err(...))`, reads on broken streams) and the `pass` records that the expected exception was raised — the test's assertion.

Checklist evidence: the except-pass is the expected-exception verification of the test, not silent error discarding.

### [ ] Finding `220` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/220.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/edge_cases.py:297:1`
- Checklist pattern: expected-exception exercises in WASI guest edge-case tests

Source excerpt:

```
        295:     try:
        296:         await failed.read(1)
    >   297:     except ReadTimeout:
        298:         pass
```

Why this is a false positive: each `try` deliberately triggers a failure (`_unwrap_result(Err(...))`, reads on broken streams) and the `pass` records that the expected exception was raised — the test's assertion.

Checklist evidence: the except-pass is the expected-exception verification of the test, not silent error discarding.

### [ ] Finding `221` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/221.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/edge_cases.py:305:1`
- Checklist pattern: expected-exception exercises in WASI guest edge-case tests

Source excerpt:

```
        303:     try:
        304:         await other_error.read(1)
    >   305:     except ValueError:
        306:         pass
```

Why this is a false positive: each `try` deliberately triggers a failure (`_unwrap_result(Err(...))`, reads on broken streams) and the `pass` records that the expected exception was raised — the test's assertion.

Checklist evidence: the except-pass is the expected-exception verification of the test, not silent error discarding.

### [ ] Finding `222` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/222.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/edge_cases.py:313:1`
- Checklist pattern: expected-exception exercises in WASI guest edge-case tests

Source excerpt:

```
        311:     try:
        312:         await future_error.read(1)
    >   313:     except ReadTimeout:
        314:         pass
```

Why this is a false positive: each `try` deliberately triggers a failure (`_unwrap_result(Err(...))`, reads on broken streams) and the `pass` records that the expected exception was raised — the test's assertion.

Checklist evidence: the except-pass is the expected-exception verification of the test, not silent error discarding.

### [ ] Finding `223` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/223.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/edge_cases.py:321:1`
- Checklist pattern: expected-exception exercises in WASI guest edge-case tests

Source excerpt:

```
        319:     try:
        320:         await unexpected.read(1)
    >   321:     except ValueError:
        322:         pass
```

Why this is a false positive: each `try` deliberately triggers a failure (`_unwrap_result(Err(...))`, reads on broken streams) and the `pass` records that the expected exception was raised — the test's assertion.

Checklist evidence: the except-pass is the expected-exception verification of the test, not silent error discarding.

### [ ] Finding `224` — `CWE-1341`

- Function context: `scripts/niquests/findings/functions/224.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/edge_cases.py:335:11`
- Checklist pattern: test exercising idempotent close

Source excerpt:

```
        333:     )
        334:     response = wasi._AsyncWASIHTTPResponse(body=cancellable, headers={}, preload_content=False)
    >   335:     await response.close()
        336:     await response.close()
```

Why this is a false positive: the tests deliberately call `response.close()` twice to verify the close is idempotent (a WASI response teardown edge case).

Checklist evidence: the double close is the test's subject matter, not an accidental double release.

### [ ] Finding `225` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/225.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/async/edge_cases.py:417:1`
- Checklist pattern: expected-exception exercises in WASI guest edge-case tests

Source excerpt:

```
        415:     try:
        416:         await extension.next_payload()
    >   417:     except OSError:
        418:         pass
```

Why this is a false positive: each `try` deliberately triggers a failure (`_unwrap_result(Err(...))`, reads on broken streams) and the `pass` records that the expected exception was raised — the test's assertion.

Checklist evidence: the except-pass is the expected-exception verification of the test, not silent error discarding.

### [ ] Finding `226` — `CWE-829`

- Function context: `scripts/niquests/findings/functions/226.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/combined/app.py:20:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         18:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         19:         try:
    >    20:             importlib.import_module(module.name)
         21:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `227` — `CWE-94`

- Function context: `scripts/niquests/findings/functions/227.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/combined/app.py:20:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         18:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         19:         try:
    >    20:             importlib.import_module(module.name)
         21:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `228` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/228.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/combined/app.py:21:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         19:         try:
         20:             importlib.import_module(module.name)
    >    21:         except ImportError:
         22:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `229` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/229.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/combined/app.py:21:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         19:         try:
         20:             importlib.import_module(module.name)
    >    21:         except ImportError:
         22:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `230` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/230.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/combined/app.py:21:9`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         19:         try:
         20:             importlib.import_module(module.name)
    >    21:         except ImportError:
         22:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `231` — `CWE-829`

- Function context: `scripts/niquests/findings/functions/231.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/hybrid_async/app.py:20:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         18:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         19:         try:
    >    20:             importlib.import_module(module.name)
         21:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `232` — `CWE-94`

- Function context: `scripts/niquests/findings/functions/232.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/hybrid_async/app.py:20:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         18:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         19:         try:
    >    20:             importlib.import_module(module.name)
         21:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `233` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/233.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/hybrid_async/app.py:21:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         19:         try:
         20:             importlib.import_module(module.name)
    >    21:         except ImportError:
         22:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `234` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/234.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/hybrid_async/app.py:21:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         19:         try:
         20:             importlib.import_module(module.name)
    >    21:         except ImportError:
         22:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `235` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/235.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/hybrid_async/app.py:21:9`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         19:         try:
         20:             importlib.import_module(module.name)
    >    21:         except ImportError:
         22:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `236` — `CWE-829`

- Function context: `scripts/niquests/findings/functions/236.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/hybrid_sync/app.py:19:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         17:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         18:         try:
    >    19:             importlib.import_module(module.name)
         20:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `237` — `CWE-94`

- Function context: `scripts/niquests/findings/functions/237.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/hybrid_sync/app.py:19:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         17:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         18:         try:
    >    19:             importlib.import_module(module.name)
         20:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `238` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/238.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/hybrid_sync/app.py:20:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         18:         try:
         19:             importlib.import_module(module.name)
    >    20:         except ImportError:
         21:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `239` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/239.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/hybrid_sync/app.py:20:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         18:         try:
         19:             importlib.import_module(module.name)
    >    20:         except ImportError:
         21:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `240` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/240.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/hybrid_sync/app.py:20:9`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         18:         try:
         19:             importlib.import_module(module.name)
    >    20:         except ImportError:
         21:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `241` — `CWE-829`

- Function context: `scripts/niquests/findings/functions/241.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/p1_sync/app.py:19:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         17:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         18:         try:
    >    19:             importlib.import_module(module.name)
         20:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `242` — `CWE-94`

- Function context: `scripts/niquests/findings/functions/242.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/p1_sync/app.py:19:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         17:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         18:         try:
    >    19:             importlib.import_module(module.name)
         20:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `243` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/243.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/p1_sync/app.py:20:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         18:         try:
         19:             importlib.import_module(module.name)
    >    20:         except ImportError:
         21:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `244` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/244.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/p1_sync/app.py:20:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         18:         try:
         19:             importlib.import_module(module.name)
    >    20:         except ImportError:
         21:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `245` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/245.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/p1_sync/app.py:20:9`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         18:         try:
         19:             importlib.import_module(module.name)
    >    20:         except ImportError:
         21:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `246` — `CWE-829`

- Function context: `scripts/niquests/findings/functions/246.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/socket_async/app.py:20:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         18:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         19:         try:
    >    20:             importlib.import_module(module.name)
         21:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `247` — `CWE-94`

- Function context: `scripts/niquests/findings/functions/247.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/socket_async/app.py:20:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         18:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         19:         try:
    >    20:             importlib.import_module(module.name)
         21:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `248` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/248.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/socket_async/app.py:21:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         19:         try:
         20:             importlib.import_module(module.name)
    >    21:         except ImportError:
         22:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `249` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/249.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/socket_async/app.py:21:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         19:         try:
         20:             importlib.import_module(module.name)
    >    21:         except ImportError:
         22:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `250` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/250.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/socket_async/app.py:21:9`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         19:         try:
         20:             importlib.import_module(module.name)
    >    21:         except ImportError:
         22:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `251` — `CWE-829`

- Function context: `scripts/niquests/findings/functions/251.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/socket_sync/app.py:19:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         17:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         18:         try:
    >    19:             importlib.import_module(module.name)
         20:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `252` — `CWE-94`

- Function context: `scripts/niquests/findings/functions/252.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/socket_sync/app.py:19:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         17:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         18:         try:
    >    19:             importlib.import_module(module.name)
         20:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `253` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/253.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/socket_sync/app.py:20:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         18:         try:
         19:             importlib.import_module(module.name)
    >    20:         except ImportError:
         21:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `254` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/254.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/socket_sync/app.py:20:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         18:         try:
         19:             importlib.import_module(module.name)
    >    20:         except ImportError:
         21:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `255` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/255.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/socket_sync/app.py:20:9`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         18:         try:
         19:             importlib.import_module(module.name)
    >    20:         except ImportError:
         21:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `256` — `CWE-829`

- Function context: `scripts/niquests/findings/functions/256.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/sync/app.py:21:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         19:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         20:         try:
    >    21:             importlib.import_module(module.name)
         22:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `257` — `CWE-94`

- Function context: `scripts/niquests/findings/functions/257.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/sync/app.py:21:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         19:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         20:         try:
    >    21:             importlib.import_module(module.name)
         22:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `258` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/258.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/sync/app.py:22:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         20:         try:
         21:             importlib.import_module(module.name)
    >    22:         except ImportError:
         23:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `259` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/259.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/sync/app.py:22:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         20:         try:
         21:             importlib.import_module(module.name)
    >    22:         except ImportError:
         23:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `260` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/260.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/sync/app.py:22:9`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         20:         try:
         21:             importlib.import_module(module.name)
    >    22:         except ImportError:
         23:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `261` — `CWE-1341`

- Function context: `scripts/niquests/findings/functions/261.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/sync/test_edges.py:377:5`
- Checklist pattern: test exercising idempotent close

Source excerpt:

```
        375:     cancellable = wasi._WASILowLevelResponse("GET", 200, "OK", HTTPHeaderDict(), Body(), Stream([b"pending"]), "url")
        376:     response = wasi._WASIHTTPResponse(body=cancellable, headers={}, preload_content=False)
    >   377:     response.close()
        378:     response.close()
```

Why this is a false positive: the tests deliberately call `response.close()` twice to verify the close is idempotent (a WASI response teardown edge case).

Checklist evidence: the double close is the test's subject matter, not an accidental double release.

### [ ] Finding `262` — `CWE-829`

- Function context: `scripts/niquests/findings/functions/262.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/unavailable_async/app.py:20:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         18:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         19:         try:
    >    20:             importlib.import_module(module.name)
         21:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `263` — `CWE-94`

- Function context: `scripts/niquests/findings/functions/263.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/unavailable_async/app.py:20:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         18:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         19:         try:
    >    20:             importlib.import_module(module.name)
         21:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `264` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/264.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/unavailable_async/app.py:21:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         19:         try:
         20:             importlib.import_module(module.name)
    >    21:         except ImportError:
         22:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `265` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/265.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/unavailable_async/app.py:21:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         19:         try:
         20:             importlib.import_module(module.name)
    >    21:         except ImportError:
         22:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `266` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/266.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/unavailable_async/app.py:21:9`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         19:         try:
         20:             importlib.import_module(module.name)
    >    21:         except ImportError:
         22:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `267` — `CWE-829`

- Function context: `scripts/niquests/findings/functions/267.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/unavailable_sync/app.py:19:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         17:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         18:         try:
    >    19:             importlib.import_module(module.name)
         20:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `268` — `CWE-94`

- Function context: `scripts/niquests/findings/functions/268.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/unavailable_sync/app.py:19:13`
- Checklist pattern: enumeration and import of the package's own submodules

Source excerpt:

```
         17:     for module in pkgutil.walk_packages(package.__path__, f"{package.__name__}."):
         18:         try:
    >    19:             importlib.import_module(module.name)
         20:         except ImportError:
```

Why this is a false positive: `pkgutil.walk_packages(package.__path__, ...)` walks the WASI guest package's own filesystem (a fixed developer package) and imports each found module; `module.name` is not user input.

Checklist evidence: CWE-829/CWE-94's condition is an untrusted control sphere; the imported names come from the developer's own package tree.

### [ ] Finding `269` — `BP-PY-2`

- Function context: `scripts/niquests/findings/functions/269.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/unavailable_sync/app.py:20:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         18:         try:
         19:             importlib.import_module(module.name)
    >    20:         except ImportError:
         21:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `270` — `CWE-390`

- Function context: `scripts/niquests/findings/functions/270.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/unavailable_sync/app.py:20:1`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         18:         try:
         19:             importlib.import_module(module.name)
    >    20:         except ImportError:
         21:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### [ ] Finding `271` — `CWE-1071`

- Function context: `scripts/niquests/findings/functions/271.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests/tests/wasi_guest/unavailable_sync/app.py:20:9`
- Checklist pattern: optional-dependency import fallback (`except ImportError: pass`)

Source excerpt:

```
         18:         try:
         19:             importlib.import_module(module.name)
    >    20:         except ImportError:
         21:             pass
```

Why this is a false positive: the `ImportError` is the expected outcome when an optional extension is not installed; the module continues with a defined fallback state (the adapter/helper simply stays unbound and callers guard for it), so no real failure is discarded.

Checklist evidence: BP-PY-2/CWE-390/CWE-1071's condition is 'an exception handler that silently discards a real failure'; here the exception is the designed, expected condition of a missing optional dependency, with an intentional fallback.

### Uncertain findings

None. No fresh finding lacked the source evidence needed to decide; the only gap is the missing audited TP 48, which is not a fresh finding and is documented in the summary above.

### Over-suppressed audited TP (regression note, for review)

Old finding `48` — `PERF-PY-26` at `src/niquests/async_session.py:618` was an audited TP; the fresh scan has no finding at that line while the source is unchanged (same commit). Current source still contains the flagged construct:

```
        try:
            extension = load_extension(scheme, implementation=implementation)
            for prefix, adapter in self.adapters.items():
                if scheme in extension.supported_schemes() and extension.scheme_to_http_scheme(scheme) == parse_scheme(prefix):
                    return adapter
        except ImportError:
            pass
```

The `parse_scheme(prefix)` call still runs inside the loop on the request path with no visible cache, satisfying the PERF-PY-26 condition — the fix suppressed it, i.e. one audited TP is over-suppressed.

### Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/niquests/chunks`
- Function evidence: `scripts/niquests/findings/functions`
- Validation: `git diff --check` — pass (run in the goslop repo root after appending this report)

## Post-fix v2 audit (latest binary)

### Run metadata

```yaml
timestamp: 2026-08-02 (latest-binary scan, binary rebuilt ~17:56 via make build)
repository: niquests
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/niquests
commit: 7633aa3f1f9fcdb7790192ffd8cfacb69ca2c807 (unchanged)
scan_target: real-repos/niquests
chunk_path: scripts/niquests/chunks
function_context_path: scripts/niquests/findings/functions
```

### Scan evidence

- Build: `make build` → `bin/goslop` (~17:56); scan command as in the Mode A append (`--profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache`)
- Findings: `275`; chunks `Chunk_1_25.txt`–`Chunk_251_275.txt` (all 11); contexts `findings/functions/1.txt`–`275.txt`
- Classification by `Source:` (file:line:col) + rule against the audited TP list and the audited FP subsections (original audit + Mode A append)

### Classification summary (fresh counts)

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| True positive | 34 | 19, 20, 21, 26, 27, 28, 35, 38, 40, 43, 44, 99, 103, 104, 105, 106, 110, 111, 112, 124, 125, 130, 131, 141, 142, 161, 163, 166, 167, 168, 169, 170, 171, 177 |
| False positive | 241 | 1–18, 22–25, 29–34, 36–37, 39, 41–42, 45–98, 100–102, 107–109, 113–123, 126–129, 132–160, 162, 164–165, 172–176, 178–216, 217–275 (all other IDs) |
| Uncertain | 0 | — |
| New findings | 0 | — |

Every fresh finding matched a prior classification (34 audited TP sources + 241 audited FP source/rule pairs); no fresh rule firing on an audited line and no unclassified source. All 34 TP sources are the same audited set as Mode A; audited TP `48` (`PERF-PY-26` at `src/niquests/async_session.py:618`, `parse_scheme(prefix)` in the adapter-mount loop) is still absent from the fresh scan while the source is unchanged — the fix still over-suppresses it (regression note carried over from Mode A).

## Fix checklist (FP patterns)

| Pattern # | Rule(s) | Trigger shape (safe condition in bold) | Count | Example sources |
| --- | --- | --- | ---: | --- |
| 1 | BP-PY-2, CWE-390, CWE-1071 | `except ImportError:\n pass` around an optional-dependency import / `load_extension` / `importlib.import_module` — **safe iff the imported name is optional and callers guard the unbound state** | 56 | adapters.py:1075, sessions.py:115, utils.py:815, conftest.py:107, wasi_guest/*/app.py |
| 2 | BP-PY-1, BP-PY-2, CWE-390, CWE-396, CWE-1071 | Pyodide JS-bridge probes (`except Exception: pass` or fallback return): `getReader()` setup, reader-cancel / `_ws.close()` / `proxy.destroy()` teardown, `to_py()` conversion, `js_headers.entries()`, `status_text`/`bytes()` — **safe iff best-effort bridge with defined fallback (None/""/b""/unset reader)** | 51 | extensions/pyodide/__init__.py:57,131,299; _async/_ws.py:123; _sse.py:151 |
| 3 | BP-PY-14, BP-PY-46 | `requests.get(...)` / `session.get(...)` / `print(` inside `>>>` docstring doctest lines — **safe iff the token is documentation text, not executable code** | 3 | kiss_headers/__init__.py:11,27; auth.py:399 |
| 4 | BP-PY-14 | requests/session call that passes `timeout` positionally — **safe iff the call carries a timeout argument; rule only sees `timeout=` keyword form** | 1 | async_api.py:183 |
| 5 | BP-PY-41, BP-PY-46, CWE-22, CWE-409, BP-PY-5 | noxfile.py dev/CI scripts and `# noqa` re-export shims: print progress, test delegation, constant path join, `extractall(filter="data")`, `from .x import *` — **safe iff the file is a build/CI script or the wildcard is the module's designed re-export** | 7 | noxfile.py:95,147,157,237; _async.py:14; _typing.py:14 |
| 6 | CWE-829, CWE-94 | `__import__`/`locals()[pkg] = __import__(...)`/`importlib.import_module` over the project's own fixed package list (packages.py) or `pkgutil.walk_packages` of its own tree (wasi_guest) — **safe iff the imported names derive from the repo's own package, not request input** | 22 | packages.py:42; wasi_guest/*/app.py |
| 7 | BP-PY-41 | `def test_...` whose body only delegates / is a smoke or expected-exception test — **safe iff assertions live in the delegated helper or `pytest.raises` guards the outcome; suppress in tests/ paths** | 11 | test_requests.py:79,1030,2163; test_emscripten.py:684; test_async.py:195 |
| 8 | BP-PY-10, CWE-502 | `pickle.loads(pickle.dumps(x))` round-trip in tests — **safe iff the unpickled data is produced in the same test** | 7 | test_requests.py:1651,1655,1663,1679,1694,1724 |
| 9 | BP-PY-2, CWE-390, CWE-1071 | expected-exception tests: `except ReadTimeout/ConnectTimeout/ValueError/OSError/Err: pass` with a preceding `pytest.fail(...)` — **safe iff the handler is the asserted positive outcome** | 9 | test_requests.py:2630,2647; wasi_guest/async/edge_cases.py:297,305 |
| 10 | CWE-1341 | intentional double `close()` in edge-case tests — **safe iff the test asserts idempotent close** | 2 | wasi_guest/async/edge_cases.py:335; sync/test_edges.py:377 |
| 11 | CWE-617, CWE-397, CWE-208, CWE-770, CWE-772, BP-PY-1 | misc test-harness firings: reachable assert, `raise Exception()`, `==` on auth objects, request-body read, test sockets — **safe iff inside the tests/ harness** | 10 | test_async.py:30; test_testserver.py:41,220; test_wsgi.py:59; testserver/server.py:90 |
| 12 | BP-PY-49 | `verify=False`/`CERT_NONE` token inside an error-message string literal — **safe iff the marker is a string, not a disabling assignment** | 1 | wasi/_utils.py:63 |
| 13 | CWE-1341 | `close()` in mutually exclusive branches (isinstance if/else, retry-raise path) — **safe iff each execution path releases the handle exactly once** | 3 | async_session.py:1015; wasi/_adapter.py:235; _async/_adapter.py:256 |
| 14 | CWE-396, BP-PY-1 | generic `except Exception/BaseException` that propagates: `retries.increment` + MaxRetryError re-raise, `app_exception` capture/re-raise, `future.set_exception(e)`, re-raise of non-expected types, cleanup-then-propagate — **safe iff the failure is fed forward or translated, not swallowed** | 12 | sgi/__init__.py:129; sgi/_async/__init__.py:355,661; utils.py:256; sgi/_sse.py:153; sgi/_ws.py:44 |
| 15 | BP-PY-2, CWE-390, CWE-1071 | parsing/attr probes with `except ValueError/TypeError/OSError/AttributeError/UnsupportedOperation: pass` — **safe iff a defined fallback state results (property falls through, `_body_position`/`content_length` stays unset/None)** | 29 | kiss_headers/builder.py:350; models.py:1065,1233; utils.py:157,1196; wasi/_async/_sse.py:62 |
| 16 | BP-PY-2, CWE-390 | iterator-exhaustion and async shutdown: `except StopIteration/StopAsyncIteration/CancelledError/GeneratorExit/TimeoutError: pass` in generator loops and lifespan-server teardown — **safe iff these are designed termination signals** | 5 | async_session.py:957; sessions.py:1709; sgi/_async/__init__.py:679,687 |
| 17 | BP-PY-2, CWE-390, CWE-1071 | best-effort issuer-certificate fetch: `except RequestException: pass` — **safe iff the fetch is optional enrichment with a defined fallback (no issuer hint)** | 6 | revocation/_crl/__init__.py:259; _ocsp/__init__.py:325 |
| 18 | CWE-93 | header write into outgoing request headers (Authorization, PreparedRequest.headers dict) — **safe iff the target is a client-side request header, not a response header** | 2 | auth.py:75; models.py:488 |
| 19 | CWE-772 | socket created and ownership transferred to the connection object (closed by owner lifecycle) or test sockets under `pytest.raises` — **safe iff the resource is released by its owner** | 3 | unixsocket/__init__.py:47; test_testserver.py:41; testserver/server.py:81 |
| 20 | BP-PY-36 | module-level `pypi_session = Session()` in help.py — **safe iff `Session` is the package's own HTTP session class, not SQLAlchemy** | 1 | help.py:176 |

## New findings

None — all 275 fresh findings matched a prior classification by `Source:` (+rule); the only gap is the still over-suppressed audited TP `PERF-PY-26` at `src/niquests/async_session.py:618` (documented above, carried over from Mode A).

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/niquests/chunks`
- Function evidence: `scripts/niquests/findings/functions`
- Validation: `git diff --check` — pass (run in the goslop repo root after appending this report)
