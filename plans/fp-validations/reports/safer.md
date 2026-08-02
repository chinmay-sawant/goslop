# False-positive audit: safer

## Run metadata

```yaml
timestamp: 2026-08-02T07:17:41Z
repository: safer
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer
branch: main
commit: eae83f7df824752540ad1e67d50099e13c86a647
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer
chunk_path: scripts/safer/chunks
function_context_path: scripts/safer/findings/functions
```

## Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/safer/chunks -context-dir scripts/safer/findings/functions real-repos/safer`
- Findings: `40`
- Chunks reviewed: `scripts/safer/chunks/Chunk_1_25.txt`, `scripts/safer/chunks/Chunk_26_40.txt`
- Function contexts reviewed: `scripts/safer/findings/functions/1.txt` .. `scripts/safer/findings/functions/40.txt` (all 40)

## Audit checklist

- [x] Read every assigned chunk under `scripts/safer/chunks`.
- [x] Read `scripts/safer/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 33 | 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 24, 26, 28, 29, 30, 32, 33, 34, 35, 37, 39, 40 |
| True positive | 7 | 22, 23, 25, 27, 31, 36, 38 |
| Uncertain | 0 | — |

## False positives

### [ ] Finding 1 — BP-PY-7

- Function context: `scripts/safer/findings/functions/1.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:6:30`
- Checklist pattern: `open(` token inside module docstring prose, not a call

Source excerpt:

```
1: """# 🧿 `safer`: A safer writer 🧿
...
5: `safer` wraps file streams, sockets, or a callable, and offers a drop-in
6: replacement for regular old `open()`.
```

Why this is a false positive: Line 6 is prose inside the module docstring (lines 1–149); the `open()` token is documentation text, so there is no `open` call to leak a resource.

Checklist evidence: The rule fires on the substring `open(` on a line without `with`, but the flagged line is inside the module docstring `"""..."""` — no executable `open` call exists.

### [ ] Finding 2 — BP-PY-46

- Function context: `scripts/safer/findings/functions/2.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:16:9`
- Checklist pattern: `print(` token inside module docstring example code, not executable

Source excerpt:

```
12:     import safer
13:
14:     with safer.open(filename, 'w') as fp:
15:         fp.write('one')
16:         print('two', file=fp)
17:         raise ValueError
```

Why this is a false positive: This is the docstring's "tiny example" (lines 12–18), indented example code inside the module docstring, not executable library code.

Checklist evidence: The `print(` occurs in a code sample embedded in the module docstring (line 1–149); BP-PY-46 targets `print` in executable library code, and no such statement exists.

### [ ] Finding 3 — BP-PY-7

- Function context: `scripts/safer/findings/functions/3.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:51:9`
- Checklist pattern: `open(` token inside module docstring prose

Source excerpt:

```
48: * `safer.writer()` wraps an existing writer, socket or stream and writes a
49:   whole response or nothing
50:
51: * `safer.open()` is a drop-in replacement for built-in `open` that
52:   writes a whole file or nothing
```

Why this is a false positive: Line 51 is prose inside the module docstring describing the `safer.open` API; there is no `open` call.

Checklist evidence: `open(` appears in docstring prose; the rule requires an `open(` call outside a `with`, and no call exists on this line.

### [ ] Finding 4 — BP-PY-7

- Function context: `scripts/safer/findings/functions/4.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:61:30`
- Checklist pattern: `open(` token inside module docstring prose

Source excerpt:

```
59:   write to file streams or any other callable.
60:
61: * `safer.printer()` is `safer.open()` except that it yields a
62:   a function that prints to the stream.
```

Why this is a false positive: Line 61 is docstring prose; the `safer.open()` token is documentation, not a call.

Checklist evidence: Same as finding 3 — the `open(` token is inside the module docstring (lines 1–149), so no resource leak is possible.

### [ ] Finding 5 — BP-PY-7

- Function context: `scripts/safer/findings/functions/5.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:67:29`
- Checklist pattern: `open(` token inside module docstring prose

Source excerpt:

```
65: or `io.BytesIO`.
66:
67: For very large files, `safer.open()` has a `temp_file` argument which
68: writes the data to a temporary file on disk, which is moved over using
```

Why this is a false positive: Line 67 is docstring prose describing the `temp_file` argument; no `open` call occurs.

Checklist evidence: `open(` is inside the module docstring; no executable `open(` call is present on the flagged line.

### [ ] Finding 6 — BP-PY-1

- Function context: `scripts/safer/findings/functions/6.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:87:1`
- Checklist pattern: `except Exception:` token inside module docstring example, not code

Source excerpt:

```
81: The old, dangerous way goes like this.
82:
83:     try:
84:         write_header(sock)
85:         write_body(sock)   # Exception is thrown here
86:         write_footer(sock)
87:      except Exception:
88:         write_error(sock)  # Oops, the header was already written
```

Why this is a false positive: This is the docstring's "old, dangerous way" example (lines 83–88) — indented sample code inside the module docstring, not an executable `except` clause.

Checklist evidence: The flagged `except Exception:` line is text inside the module docstring (lines 1–149); no real exception handler exists at this location.

### [ ] Finding 7 — BP-PY-1

- Function context: `scripts/safer/findings/functions/7.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:97:1`
- Checklist pattern: `except Exception:` token inside module docstring example, not code

Source excerpt:

```
92:     try:
93:         with safer.writer(sock) as s:
94:             write_header(s)
95:             write_body(s)  # Exception is thrown here
96:             write_footer(s)
97:      except Exception:
98:         write_error(sock)  # Nothing has been written
```

Why this is a false positive: This is the docstring's "With `safer`" example (lines 92–98) — sample code inside the module docstring, not an executable handler.

Checklist evidence: Same as finding 6 — the `except Exception:` line is docstring text; BP-PY-1's condition (broad except in code with non-re-raising suite) cannot be satisfied by documentation prose.

### [ ] Finding 8 — BP-PY-7

- Function context: `scripts/safer/findings/functions/8.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:102:7`
- Checklist pattern: `open(` token inside module docstring prose

Source excerpt:

```
100: ### Example: `safer.open()` and json
101:
102: `safer.open()` is a a drop-in replacement for built-in `open()` except that
103: when used as a context, it leaves the original file unchanged on failure.
```

Why this is a false positive: Line 102 is docstring prose; `safer.open()` / `open()` are documentation tokens, not calls.

Checklist evidence: `open(` appears in the module docstring; the rule's call condition is not met.

### [ ] Finding 9 — BP-PY-7

- Function context: `scripts/safer/findings/functions/9.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:117:7`
- Checklist pattern: `open(` token inside module docstring prose

Source excerpt:

```
115:         # If an exception is raised, the file is unchanged.
116:
117: `safer.open(filename)` returns a file stream `fp` like `open(filename)`
118: would, except that `fp` writes to memory stream or a temporary file in the
```

Why this is a false positive: Line 117 is docstring prose describing the return value of `safer.open`; no call is made.

Checklist evidence: `open(` is inside the module docstring; no executable `open` call on this line.

### [ ] Finding 10 — BP-PY-7

- Function context: `scripts/safer/findings/functions/10.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:129:39`
- Checklist pattern: `open(` token inside module docstring prose

Source excerpt:

```
127: ### Example: `safer.printer()`
128:
129: `safer.printer()` is similar to `safer.open()` except it yields a function
130: that prints to the open file - it's very convenient for printing text.
```

Why this is a false positive: Line 129 is docstring prose; `safer.open()` is a documentation token, not a call.

Checklist evidence: `open(` occurs inside the module docstring; the rule's call condition is not met.

### [ ] Finding 11 — BP-PY-7

- Function context: `scripts/safer/findings/functions/11.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:132:12`
- Checklist pattern: `open(` token inside module docstring prose

Source excerpt:

```
130: that prints to the open file - it's very convenient for printing text.
131:
132: Like `safer.open()`, if an exception is raised within its context manager,
133: the original file is left unchanged.
```

Why this is a false positive: Line 132 is docstring prose; `safer.open()` is a documentation token, not a call.

Checklist evidence: `open(` is inside the module docstring; no executable `open` call exists.

### [ ] Finding 12 — BP-PY-46

- Function context: `scripts/safer/findings/functions/12.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:139:13`
- Checklist pattern: `print(` token inside module docstring example, not code

Source excerpt:

```
135: Before.
136:
137:     with open(file, 'w') as fp:
138:         for item in items:
139:             print(item, file=fp)
140:         # Prints lines until the first exception
```

Why this is a false positive: This is the docstring's "Before" example (lines 137–140) — sample code inside the module docstring, not executable library code.

Checklist evidence: `print(` occurs in docstring example text (module docstring lines 1–149); BP-PY-46's condition (print in executable library code outside a main guard) is not met.

### [ ] Finding 13 — BP-PY-46

- Function context: `scripts/safer/findings/functions/13.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:146:13`
- Checklist pattern: `print(` token inside module docstring example, not code

Source excerpt:

```
142: With `safer`
143:
144:     with safer.printer(file) as print:
145:         for item in items:
146:             print(item)
147:         # Either the whole file is written, or nothing
```

Why this is a false positive: This is the docstring's "With `safer`" example (lines 144–147) — sample code inside the module docstring, not executable library code.

Checklist evidence: Same as finding 12 — the `print(` is docstring example text, not executable library code.

### [ ] Finding 14 — BP-PY-7

- Function context: `scripts/safer/findings/functions/14.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:225:16`
- Checklist pattern: the flagged `open(` call targets the module's own `open` (proven by non-builtin kwargs) and its return value is returned to the caller

Source excerpt:

```
224:     if isinstance(stream, (str, Path)):
225:         return open(
226:             stream,
227:             'wb' if is_binary else 'w',
228:             delete_failures=delete_failures,
229:             dry_run=dry_run,
230:             enabled=enabled,
231:         )
```

Why this is a false positive: The callee is the module's own `safer.open` (line 323) — the arguments `delete_failures`, `dry_run`, `enabled` are not builtin `open` parameters, so this can only be the library's context-manager-compatible replacement. Its returned stream is the function's return value, handed to the caller for management, so no resource is leaked by `writer()` itself.

Checklist evidence: The rule flags any `open(` outside `with`, but the shown source proves the call is to the library's own `open` (builtin `open` would raise TypeError on these kwargs) and the stream is returned, not leaked.

### [ ] Finding 15 — CWE-396

- Function context: `scripts/safer/findings/functions/15.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:312:1`
- Checklist pattern: generic handler suite re-raises, so no failure condition is hidden

Source excerpt:

```
312:     except Exception:
313:         if close_on_exit:
314:             getattr(stream, 'close', lambda: None)()
315:         raise
```

Why this is a false positive: The generic handler only performs cleanup (closing the stream) and immediately re-raises with `raise`; it cannot hide distinct failure conditions, which is the weakness CWE-396 describes. The sibling rule BP-PY-1 explicitly skips broad excepts whose suite re-raises (`suiteReraises`), confirming this shape is not a weakness.

Checklist evidence: CWE-396's message is "generic Exception handler can hide distinct failure conditions"; the shown source's handler re-raises, so the condition (failures hidden) is not satisfied.

### [ ] Finding 16 — BP-PY-7

- Function context: `scripts/safer/findings/functions/16.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:323:5`
- Checklist pattern: `open(` token is a `def` statement, not a call

Source excerpt:

```
322:
323: def open(
324:     name: Path | str,
325:     mode: str = 'r',
```

Why this is a false positive: Line 323 is the definition of the module's `open` function (`def open(`); no `open` call occurs, so no resource can be leaked.

Checklist evidence: BP-PY-7 flags `open(` outside `with`, but the token is a function definition header, not a call — the rule's call condition is not satisfied.

### [ ] Finding 17 — BP-PY-7

- Function context: `scripts/safer/findings/functions/17.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:358:20`
- Checklist pattern: `open(` token inside `safer.open` docstring prose

Source excerpt:

```
356:       enabled:
357:          If `enabled` is falsey, safer is entirely bypassed, and
358:          built-in `open()` is used instead.
359:
```

Why this is a false positive: Line 358 is prose inside the `open()` function's docstring (lines 338–387); `open()` is documentation text, not a call.

Checklist evidence: `open(` is inside a docstring; no executable `open` call exists on the flagged line.

### [ ] Finding 18 — BP-PY-7

- Function context: `scripts/safer/findings/functions/18.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:360:59`
- Checklist pattern: `open(` token inside `safer.open` docstring prose

Source excerpt:

```
358:          built-in `open()` is used instead.
359:
360:     The remaining arguments are the same as for built-in `open()`.
361:
```

Why this is a false positive: Line 360 is prose inside the `open()` function's docstring; `open()` is a documentation token, not a call.

Checklist evidence: `open(` is inside a docstring; no executable call exists.

### [ ] Finding 19 — BP-PY-7

- Function context: `scripts/safer/findings/functions/19.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:362:11`
- Checklist pattern: `open(` token inside `safer.open` docstring prose

Source excerpt:

```
360:     The remaining arguments are the same as for built-in `open()`.
361:
362:     `safer.open() is a drop-in replacement for built-in`open()`. It returns a
363:     stream which only overwrites the original file when close() is called, and
```

Why this is a false positive: Line 362 is prose inside the `open()` function's docstring; the `open()` tokens are documentation text, not calls.

Checklist evidence: `open(` is inside a docstring; no executable `open` call exists.

### [ ] Finding 20 — BP-PY-7

- Function context: `scripts/safer/findings/functions/20.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:368:40`
- Checklist pattern: `open(` token inside `safer.open` docstring prose

Source excerpt:

```
366:     It works as follows:
367:
368:     If a stream `fp` return from `safer.open()` is used as a context
369:     manager and an exception is raised, the property `fp.safer_failed` is
```

Why this is a false positive: Line 368 is prose inside the `open()` function's docstring; `safer.open()` is a documentation token, not a call.

Checklist evidence: `open(` is inside a docstring; no executable `open` call exists.

### [ ] Finding 21 — CWE-367

- Function context: `scripts/safer/findings/functions/21.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:402:12`
- Checklist pattern: benign directory pre-check; the matched "use" token is a substring artifact, not a use of the checked path

Source excerpt:

```
400:     name = os.path.realpath(name)
401:     parent = os.path.dirname(os.path.abspath(name))
402:     if not os.path.exists(parent):
403:         if not make_parents:
404:             raise OSError('Directory does not exist')
405:         os.makedirs(parent)
...
407:     def simple_open():
408:         return __builtins__['open'](name, mode, buffering, **kwargs)
```

Why this is a false positive: The check `os.path.exists(parent)` is a convenience pre-check whose result is only used to produce a clearer error or to call `os.makedirs(parent)`; the "later separate use" the rule matched is the substring `open(` inside `simple_open()`/`simple_open()`, not an `open`/`os.remove`/`os.unlink` call on the checked path. The actual open uses the realpath'd `name` with a fresh `os.makedirs`-safe flow; no security decision depends on the stale existence value.

Checklist evidence: CWE-367's regex requires a path existence check followed by `open(`/`os.remove(`/`os.unlink(` on the same construct; here the match is the word `open(` embedded in `simple_open()`, and the checked path `parent` is never the operand of the later use — the TOCTOU condition is not satisfied.

### [ ] Finding 24 — BP-PY-7

- Function context: `scripts/safer/findings/functions/24.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:554:25`
- Checklist pattern: `open(` token inside `printer` docstring prose

Source excerpt:

```
552:
553:     ARGUMENTS
554:       Same as for `safer.open()`
555:     """
```

Why this is a false positive: Line 554 is the last line of the `printer()` function's docstring (lines 548–555); `safer.open()` is documentation text, not a call.

Checklist evidence: `open(` is inside a docstring; no executable `open` call exists.

### [ ] Finding 26 — CWE-459

- Function context: `scripts/safer/findings/functions/26.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:621:29`
- Checklist pattern: temp-file cleanup exists in the owning class's lifecycle methods, so the file does not persist

Source excerpt:

```
618: class _FileCloser(_Closer):
619:     def __init__(self, temp_file, delete_failures, parent=None):
620:         if temp_file is True:
621:             fd, temp_file = tempfile.mkstemp(dir=parent)
622:             os.close(fd)
...
627:     def _failure(self):
628:         if self.delete_failures:
629:             os.remove(self.temp_file)
630:         else:
631:             print('Temp_file saved:', self.temp_file, file=sys.stderr)
...
679:     def _success(self):
680:         if not self.dry_run:
681:             if os.path.exists(self.target_file):
682:                 shutil.copymode(self.target_file, self.temp_file)
683:             os.replace(self.temp_file, self.target_file)
```

Why this is a false positive: The mkstemp file is a staging file whose lifecycle is owned by the `_FileCloser` object: on failure it is removed by `_failure()` (`os.remove(self.temp_file)`, line 629), and on success it is renamed over the target by `_success()` (`os.replace`, line 683). The file never persists, so there is no incomplete cleanup; the rule's same-function heuristic simply cannot see cleanup implemented in sibling methods of the same class.

Checklist evidence: CWE-459's condition is "persistent temporary file has no same-function unlink cleanup", but the shown source demonstrates the file is always removed or replaced by the owning class's close path — the "incomplete cleanup" condition is not satisfied.

### [ ] Finding 28 — BP-PY-41

- Function context: `scripts/safer/findings/functions/28.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/test/test_dump.py:25:1`
- Checklist pattern: the test delegates to a same-file helper that performs assertions

Source excerpt:

```
25:     def test_dump(self):
26:         _test()
...
78: def _test(load=json.load, dumps=DUMPS, tests=TESTS):
79:     for dump in dumps:
80:         for data in tests:
81:             dump(data, 'one')
82:
83:             with open('one') as fp:
84:                 assert load(fp) == data
```

Why this is a false positive: `test_dump` delegates to the same-file helper `_test`, whose body ends with `assert load(fp) == data` (line 84) — the outcome is verified. The rule's helper heuristic only credits `self.assert`/`AssertionError` in helpers, missing bare `assert`, so it misclassifies a thin delegating test as assertion-less.

Checklist evidence: BP-PY-41's condition is "test function performs side effects without assertions"; the shown source's helper contains an `assert` on the round-tripped data, so the test does verify outcomes.

### [ ] Finding 29 — BP-PY-41

- Function context: `scripts/safer/findings/functions/29.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/test/test_dump.py:28:1`
- Checklist pattern: the test delegates to a same-file helper that performs assertions

Source excerpt:

```
28:     def test_yaml(self):
29:         dumps = yaml, yaml.dump, 'yaml', 'yaml.dump'
30:         dumps = (dumper(i) for i in dumps)
31:
32:         _test(yaml.safe_load, dumps)
...
83:             with open('one') as fp:
84:                 assert load(fp) == data
```

Why this is a false positive: `test_yaml` delegates to the same-file helper `_test`, which asserts `load(fp) == data` (line 84) using the passed loader — the round-trip is verified.

Checklist evidence: Same as finding 28 — the delegated helper performs the assertion, so the test is not assertion-less.

### [ ] Finding 30 — BP-PY-41

- Function context: `scripts/safer/findings/functions/30.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/test/test_dump.py:34:1`
- Checklist pattern: the test delegates to a same-file helper that performs assertions

Source excerpt:

```
34:     def test_toml(self):
35:         dumps = toml, toml.dump, 'toml', 'toml.dump'
36:         dumps = (dumper(i) for i in dumps)
37:         tests = ({}, {'a': 1}, {'a': [1, 2]})
38:
39:         _test(toml.load, dumps, tests)
...
83:             with open('one') as fp:
84:                 assert load(fp) == data
```

Why this is a false positive: `test_toml` delegates to the same-file helper `_test`, which asserts `load(fp) == data` (line 84) using the passed loader — the round-trip is verified.

Checklist evidence: Same as finding 28 — the delegated helper performs the assertion, so the test is not assertion-less.

### [ ] Finding 32 — BP-PY-7

- Function context: `scripts/safer/findings/functions/32.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/test/test_files.py:28:19`
- Checklist pattern: the opened stream is explicitly closed in the same function

Source excerpt:

```
28:         fp = safer.open(FILENAME, 'w', temp_file=True)
29:         fp.write('OK!')
30:         assert FILENAME.read_text() == 'hello'
...
36:         fp.close()
```

Why this is a false positive: The stream opened at line 28 is explicitly released by `fp.close()` in the same function (line 36) — the test is literally `test_explicit_close`. The resource-leak risk BP-PY-7 warns about does not occur.

Checklist evidence: The rule's condition is "open without `with` risks resource leaks", but the shown source closes the stream at line 36, so the leak condition is not satisfied.

### [ ] Finding 33 — CWE-367

- Function context: `scripts/safer/findings/functions/33.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/test/test_files.py:47:20`
- Checklist pattern: test assertion on a local fixture path; the later "use" operates on a different path

Source excerpt:

```
45:         with safer.open(FILENAME, 'w', temp_file=temp_file) as fp:
46:             assert temp_file.exists()
47:             assert os.path.exists(temp_file)
48:             fp.write('hello')
...
57:         with safer.open(FILENAME, 'w', temp_file=temp_file) as fp:
```

Why this is a false positive: `os.path.exists(temp_file)` is a test assertion on a local fixture path in a single-threaded test; the "later use" the rule matched is `safer.open(FILENAME, ...)` (line 57), which operates on a different path and is not a sensitive open/remove/unlink of the checked path. No trust boundary or security decision is involved.

Checklist evidence: CWE-367's condition requires a check followed by a use of the same path; here the use is on `FILENAME` (a different path) and the check is a test assertion — the TOCTOU condition is not satisfied.

### [ ] Finding 34 — BP-PY-7

- Function context: `scripts/safer/findings/functions/34.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/test/test_files.py:97:18`
- Checklist pattern: the call always raises before creating any stream, so no resource exists

Source excerpt:

```
96:         with self.assertRaises(ValueError):
97:             safer.open(FILENAME, 'r+')
```

Why this is a false positive: `safer.open(FILENAME, 'r+')` without `temp_file` raises `ValueError('+ mode requires a temp_file argument')` at `__init__.py:419` before any stream or file is created, so there is no resource to leak. The call is wrapped in `assertRaises` to verify exactly that.

Checklist evidence: The rule's leak condition requires an opened resource; the shown source shows a call that fails before creating any stream (mode validation precedes closer construction), so no leak is possible.

### [ ] Finding 35 — BP-PY-7

- Function context: `scripts/safer/findings/functions/35.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/test/test_files.py:109:18`
- Checklist pattern: the call always raises before creating any stream, so no resource exists

Source excerpt:

```
107:     def _error(self, mode='w', **kwds):
108:         with self.assertRaises(ValueError) as e:
109:             safer.open(FILENAME, mode, temp_file=True, **kwds)
110:         return e.exception.args[0]
```

Why this is a false positive: Every `_error(...)` invocation passes an invalid combination (`closefd=False`, `'bt'`, `newline`, `encoding`, `errors`) that raises `ValueError` during argument validation (`__init__.py:430-441`) before `_FileRenameCloser` is constructed — no temp file or stream is ever created, so there is no resource to leak.

Checklist evidence: The rule's leak condition requires an opened resource; the shown source shows a call that fails during validation before any resource is created.

### [ ] Finding 37 — CWE-367

- Function context: `scripts/safer/findings/functions/37.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/test/test_open.py:193:16`
- Checklist pattern: test assertion after a completed write; no trust boundary

Source excerpt:

```
191:         with safer_open(FILENAME, 'w'):
192:             pass
193:         assert os.path.exists(FILENAME), FILENAME
```

Why this is a false positive: `os.path.exists(FILENAME)` is an assertion verifying the effect of the preceding completed write in a single-threaded test; the "later use" the rule matched is a subsequent test's `safer_open(FILENAME, 'x')` call on the same fixture path — there is no check-then-use security decision and no attacker-controlled path.

Checklist evidence: CWE-367's condition requires a trust-boundary check followed by a separate sensitive use; the shown source is a post-condition assertion in a test with no security decision depending on it.

### [ ] Finding 39 — CWE-772

- Function context: `scripts/safer/findings/functions/39.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/test/test_writer.py:191:1`
- Checklist pattern: the assigned resource is released by a context-manager in the same function

Source excerpt:

```
192:     fp = open(FILENAME, 'w')
193:     with safer.writer(fp, close_on_exit=True):
194:         fp.write('hello, world')
```

Why this is a false positive: The rule's own message exempts "context-manager release" — `fp` is passed to `safer.writer(fp, close_on_exit=True)`, whose write wrapper executes `with stream:` (safer/__init__.py:258) so the stream is closed when the writer context exits. The resource is released within its effective lifetime.

Checklist evidence: CWE-772's condition is "resource assigned without a same-function close or context-manager release"; the shown source has a context-manager release (`with safer.writer(fp, close_on_exit=True)`), so the condition is not satisfied.

### [ ] Finding 40 — BP-PY-7

- Function context: `scripts/safer/findings/functions/40.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/test/test_writer.py:192:10`
- Checklist pattern: the opened stream is released by the writer wrapper in the next line

Source excerpt:

```
192:     fp = open(FILENAME, 'w')
193:     with safer.writer(fp, close_on_exit=True):
194:         fp.write('hello, world')
```

Why this is a false positive: The `fp` opened at line 192 is closed by the writer wrapper: with `close_on_exit=True`, the writer's write path runs `with stream:` (safer/__init__.py:258), releasing the stream. The leak BP-PY-7 warns about does not occur.

Checklist evidence: The rule's condition is "open without `with` risks resource leaks"; the shown source's very next line hands the stream to a context manager that closes it, so no leak occurs.

## True positives

### CWE-829 — dynamically selected module import (Finding 22)

| Finding id | Source | Reason |
| --- | --- | --- |
| 22 | `safer/__init__.py:528:20` | `dump = __import__(dump)` imports a module whose name is a function parameter derived from the caller's `dump` argument or the target file suffix (line 510) — a dynamic expression reaching a dynamic-import sink, exactly the rule condition (`__import__(` with non-literal first arg). |

### CWE-94 — dynamic import sink (Finding 23)

| Finding id | Source | Reason |
| --- | --- | --- |
| 23 | `safer/__init__.py:528:20` | Same construct: `__import__(dump)` with a non-literal module name; CWE-94 fires on dynamic-import sinks with a non-literal argument, which the shown source satisfies. |

### BP-PY-1 — broad except without re-raise (Findings 25, 36)

| Finding id | Source | Reason |
| --- | --- | --- |
| 25 | `safer/__init__.py:573:1` | `except Exception:` whose suite is only `traceback.print_exc()` (no re-raise) — matches the rule condition (broad except, suite does not re-raise). |
| 36 | `test/test_open.py:21:1` | `except Exception: return False` swallows every error from the `ctypes` call — broad except with no re-raise, matching the rule condition. |

### BP-PY-46 — print in library code (Finding 27)

| Finding id | Source | Reason |
| --- | --- | --- |
| 27 | `safer/__init__.py:631:13` | Real executable `print(...)` in library code (`_FileCloser._failure`), not under a `__main__` guard, not a test/benchmark — matches the rule condition. |

### BP-PY-7 — open without with in test code (Finding 31)

| Finding id | Source | Reason |
| --- | --- | --- |
| 31 | `test/test_dump.py:51:32` | `json.dump(data, one.open('w'))` — a `Path.open('w')` call not wrapped in `with`; the handle is never closed (the call raises inside `assertRaises`, leaking the already-opened fd). |

### BP-PY-41 — side-effect-only test (Finding 38)

| Finding id | Source | Reason |
| --- | --- | --- |
| 38 | `test/test_writer.py:126:1` | `test_none` only writes to stdout inside the writer context with no assertion or `pytest.raises` anywhere — matches the rule's condition of a side-effect-only test. |

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none (single-reviewer audit)
- Chunk evidence: `scripts/safer/chunks`
- Function evidence: `scripts/safer/findings/functions`
- Validation: `git diff --check` — pass

## Post-fix remaining-FP audit (2026-08-02)

### Run metadata

```yaml
timestamp: 2026-08-02T11:08:51Z
repository: safer
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer
branch: main
commit: eae83f7df824752540ad1e67d50099e13c86a647
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer
chunk_path: scripts/safer/chunks
function_context_path: scripts/safer/findings/functions
scanner_binary: bin/goslop (rebuilt 2026-08-02 16:29 +0530, post-fix b5b8fde)
```

### Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/safer/chunks -context-dir scripts/safer/findings/functions real-repos/safer`
- Findings: `13`
- Chunks reviewed: `scripts/safer/chunks/Chunk_1_13.txt`
- Function contexts reviewed: `scripts/safer/findings/functions/1.txt` .. `scripts/safer/findings/functions/13.txt` (all 13)

### Audit checklist

- [x] Read every assigned chunk under `scripts/safer/chunks`.
- [x] Read `scripts/safer/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

### Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 6 | 1, 2, 3, 8, 12, 13 |
| True positive | 7 | 4, 5, 6, 7, 9, 10, 11 |
| Uncertain | 0 | — |

### False positives

All six remaining false positives re-appear at source locations that the original audit already classified as false positives (matched by `Source:` path). Each was re-verified against its rule condition on the fresh scan.

#### [ ] Finding 1 — BP-PY-1

- Function context: `scripts/safer/findings/functions/1.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:87:1`
- Checklist pattern: `except Exception:` token inside module docstring example, not code (same construct as audited FP 6)

Source excerpt:

```
83:     try:
84:         write_header(sock)
85:         write_body(sock)   # Exception is thrown here
86:         write_footer(sock)
87:      except Exception:
88:         write_error(sock)  # Oops, the header was already written
```

Why this is a false positive: The flagged line is the docstring's "old, dangerous way" example (module docstring lines 1–149), not an executable handler, so BP-PY-1's condition (broad except in executable code without re-raise) cannot be met.

Checklist evidence: The source line is unchanged from audited FP 6 (`safer/__init__.py:87:1`); it remains text inside the module docstring — no real exception handler exists.

#### [ ] Finding 2 — BP-PY-1

- Function context: `scripts/safer/findings/functions/2.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:97:1`
- Checklist pattern: `except Exception:` token inside module docstring example, not code (same construct as audited FP 7)

Source excerpt:

```
92:     try:
93:         with safer.writer(sock) as s:
94:             write_header(s)
95:             write_body(s)  # Exception is thrown here
96:             write_footer(s)
97:      except Exception:
98:         write_error(sock)  # Nothing has been written
```

Why this is a false positive: Same docstring-example shape as finding 1 ("With `safer`" example, lines 92–98 of the module docstring) — documentation text, not an executable handler.

Checklist evidence: The source line is unchanged from audited FP 7 (`safer/__init__.py:97:1`); the `except Exception:` token is inside the module docstring, so no executable broad-except handler exists.

#### [ ] Finding 3 — BP-PY-7

- Function context: `scripts/safer/findings/functions/3.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:225:16`
- Checklist pattern: the flagged `open(` call targets the module's own `open` (proven by non-builtin kwargs) and its return value is returned to the caller (same construct as audited FP 14)

Source excerpt:

```
224:     if isinstance(stream, (str, Path)):
225:         return open(
226:             stream,
227:             'wb' if is_binary else 'w',
228:             delete_failures=delete_failures,
229:             dry_run=dry_run,
230:             enabled=enabled,
231:         )
```

Why this is a false positive: The callee is the library's own `safer.open` — the `delete_failures`/`dry_run`/`enabled` keyword arguments are not builtin `open` parameters, so this cannot be the builtin. The returned stream is the function's return value, handed to the caller, so nothing is leaked by `writer()` itself.

Checklist evidence: Source unchanged from audited FP 14; the rule's leak condition requires an unmanaged builtin `open`, and here the call is the library's context-manager-compatible `open` whose result is returned, not abandoned.

#### [ ] Finding 8 — CWE-459

- Function context: `scripts/safer/findings/functions/8.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/safer/__init__.py:621:29`
- Checklist pattern: temp-file cleanup exists in the owning class's lifecycle methods, so the file does not persist (same construct as audited FP 26)

Source excerpt:

```
618: class _FileCloser(_Closer):
619:     def __init__(self, temp_file, delete_failures, parent=None):
620:         if temp_file is True:
621:             fd, temp_file = tempfile.mkstemp(dir=parent)
622:             os.close(fd)
...
627:     def _failure(self):
628:         if self.delete_failures:
629:             os.remove(self.temp_file)
...
679:     def _success(self):
680:         if not self.dry_run:
681:             if os.path.exists(self.target_file):
682:                 shutil.copymode(self.target_file, self.temp_file)
683:             os.replace(self.temp_file, self.target_file)
```

Why this is a false positive: The mkstemp file's lifecycle is owned by `_FileCloser`: removed on failure by `_failure()` (`os.remove`, line 629) and renamed over the target on success by `_success()` (`os.replace`, line 683). The file never persists, so there is no incomplete cleanup; the rule's same-function heuristic cannot see cleanup in sibling class methods.

Checklist evidence: Source unchanged from audited FP 26; the "persistent temporary file without cleanup" condition is not satisfied because both lifecycle outcomes dispose of the file.

#### [ ] Finding 12 — CWE-772

- Function context: `scripts/safer/findings/functions/12.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/test/test_writer.py:191:1`
- Checklist pattern: the assigned resource is released by a context manager in the same function (same construct as audited FP 39)

Source excerpt:

```
192:     fp = open(FILENAME, 'w')
193:     with safer.writer(fp, close_on_exit=True):
194:         fp.write('hello, world')
```

Why this is a false positive: `fp` is passed to `safer.writer(fp, close_on_exit=True)`, whose wrapper executes `with stream:` (safer/__init__.py:258), so the stream is closed when the writer context exits — the rule's own message exempts context-manager release.

Checklist evidence: Source unchanged from audited FP 39 (`test_writer.py:191:1`); the "resource without same-function close or context-manager release" condition is not met because the next line hands the stream to a context manager that releases it.

#### [ ] Finding 13 — BP-PY-7

- Function context: `scripts/safer/findings/functions/13.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer/test/test_writer.py:192:10`
- Checklist pattern: the opened stream is released by the writer wrapper on the very next line (same construct as audited FP 40)

Source excerpt:

```
192:     fp = open(FILENAME, 'w')
193:     with safer.writer(fp, close_on_exit=True):
194:         fp.write('hello, world')
```

Why this is a false positive: The stream opened at line 192 is closed by the writer wrapper (with `close_on_exit=True` the write path runs `with stream:` at safer/__init__.py:258), so the leak BP-PY-7 warns about does not occur.

Checklist evidence: Source unchanged from audited FP 40 (`test_writer.py:192:10`); the same construct that satisfied CWE-772's context-manager exemption (finding 12) also defeats BP-PY-7's leak condition.

### True positives

Seven fresh findings satisfy their rule conditions and match audited true positives by `Source:` (or, for finding 7, satisfy the rule condition at a source the old scan did not report):

| Finding id | Rule | Source | Match / basis |
| --- | --- | --- | --- |
| 4 | CWE-829 | `safer/__init__.py:528:20` | audited TP 22 — `__import__(dump)` with non-literal module name |
| 5 | CWE-94 | `safer/__init__.py:528:20` | audited TP 23 — dynamic import sink with non-literal argument |
| 6 | BP-PY-1 | `safer/__init__.py:573:1` | audited TP 25 — broad except, suite only `traceback.print_exc()` (no re-raise) |
| 7 | CWE-396 | `safer/__init__.py:573:1` | new source for CWE-396; suite is `traceback.print_exc()` with no re-raise, which the post-fix detector (`rules_platform.go:107`) does not treat as surfacing the failure, so the handler hides the failure of `_close(True)` — satisfies the rule condition |
| 9 | BP-PY-46 | `safer/__init__.py:631:13` | audited TP 27 — real executable `print(...)` in library code |
| 10 | BP-PY-1 | `test/test_open.py:21:1` | audited TP 36 — `except Exception: return False` swallows errors |
| 11 | BP-PY-41 | `test/test_writer.py:126:1` | audited TP 38 — side-effect-only test with no assertion |

For finding 7 the source is:

```
570:         except Exception:  # pragma: no cover
571:             try:
572:                 self._close(True)
573:             except Exception:
574:                 traceback.print_exc()
575:             raise
```

Why finding 7 is a true positive: the inner handler at 573:1 catches the failure of the cleanup call `self._close(True)` and only prints its traceback; the exception is not propagated (`raise` at line 575 re-raises only the outer `parent_close` failure). The cleanup failure condition is therefore hidden, matching CWE-396's condition. The old scan reported CWE-396 at the re-raising handler (`__init__.py:312:1`, audited FP 15), which the post-fix detector now skips; the non-re-raising handler at 573:1 is its replacement, and it is a genuine match.

### Uncertain findings

None.

### Final evidence

- Delegated reviewers: none (single-reviewer audit)
- Chunk evidence: `scripts/safer/chunks/Chunk_1_13.txt`
- Function evidence: `scripts/safer/findings/functions/1.txt` .. `13.txt`
- Validation: `git diff --check` — pass
- Note: audited TP 31 (`test_dump.py:51:32`, BP-PY-7) did not re-fire in the fresh scan; out of Mode A scope (potential over-suppression, not assessed here).

## Post-fix v2 audit (latest binary)

### Run metadata

```yaml
timestamp: 2026-08-02 (v2 audit)
repository: safer
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/safer
branch: main
commit: eae83f7df824752540ad1e67d50099e13c86a647 (unchanged)
chunk_path: scripts/safer/chunks
function_context_path: scripts/safer/findings/functions
scanner_binary: bin/goslop (rebuilt ~2026-08-02 17:56)
```

### Scan evidence

- Findings: `13` (chunk `scripts/safer/chunks/Chunk_1_13.txt`; contexts 1..13 all read; enclosing source re-checked at `safer/__init__.py:312,573` for the re-raise shapes)
- All 13 findings match a prior audit classification by `Source:` — no fresh classification needed beyond reuse; re-verified each source is unchanged from the audited construct.

### Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 7 | 1, 2, 3, 4, 8, 12, 13 |
| True positive | 6 | 5, 6, 7, 9, 10, 11 |
| Uncertain | 0 | — |

Matches: 1→FP6, 2→FP7, 3→FP14, 4→FP15, 5→TP22, 6→TP23, 7→TP25, 8→FP26, 9→TP27, 10→TP36, 11→TP38, 12→FP39, 13→FP40.

## Fix checklist (FP patterns)

| Pattern # | Rule | Trigger shape | Count | Example sources |
| --- | --- | --- | ---: | --- |
| 1 | BP-PY-1 | `except Exception:` token inside module docstring example — `detectBPPY1` scans unmasked source (`codeLinesFacts(..., unit.Source)`), only strips comments; fix: mask string/docstring content like BP-PY-7 does (`pytext.Mask`, rules_core.go:356) so a literal-in-string cannot be a handler | 2 | `safer/__init__.py:87:1`, `:97:1` |
| 2 | BP-PY-7 | bare `open(` that resolves to the module's own `open` replacement — non-builtin kwargs (`delete_failures`/`dry_run`/`enabled`) prove callee is `def open` at `__init__.py:323`, and its result is returned to the caller, not leaked; fix: skip bare `open(` when the callee is module-defined or kwargs are non-builtin, or when the result is a return value owned by the caller | 1 | `safer/__init__.py:225:16` |
| 3 | CWE-396 | generic `except Exception:` whose suite re-raises (cleanup + bare `raise`) — current `detectCWE396` (rules_platform.go:41-60) has no suite/re-raise exemption and Mode A's b5b8fde build skipped this handler, so this re-fire is a regression; fix: restore exemption for suites containing a bare `raise` | 1 | `safer/__init__.py:312:1` |
| 4 | CWE-459 | `tempfile.mkstemp` in class `__init__` whose file is disposed by sibling class methods (`_failure` → `os.remove`, `_success` → `os.replace`) — same-function heuristic can't see same-class lifecycle; fix: skip when the owning class has both removal and replace paths for the temp file | 1 | `safer/__init__.py:621:29` |
| 5 | CWE-772 + BP-PY-7 | resource assigned (`fp = open(...)`) then handed to a `with`-managed wrapper on the next statement (`with safer.writer(fp, close_on_exit=True)` closes it via `with stream:`) — neither rule looks ahead to the downstream context-manager release; fix: treat a same-function context-manager release as satisfying cleanup/leak conditions (one shared shape, two rules) | 2 | `test/test_writer.py:191:1`, `:192:10` |

## New findings

None — every fresh finding (13/13) had a prior audit classification at the identical `Source:`; no unclassified constructs.

Coverage observations for the fix phase (not new findings):

- Mode A TP 7 (CWE-396 @ `573:1`, non-re-raising handler) is no longer emitted: `detectCWE396` returns after the first match per module (rules_platform.go:59), and 312 precedes 573 — the earlier FP location shadows the real one.
- Audited TP 31 (`test_dump.py:51:32`, BP-PY-7) still does not re-fire; out of v2 scope.
