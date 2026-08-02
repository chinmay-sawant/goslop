# False-positive audit: pictex

## Run metadata

```yaml
timestamp: 2026-08-02T07:41:44Z
repository: pictex
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex
branch: main
commit: 07072311d51917679fe41a827c7136195fb0dcf5
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex
chunk_path: ./scripts/chunks
function_context_path: ./scripts/findings/functions
```

## Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir real-repos/pictex/scripts/chunks -context-dir real-repos/pictex/scripts/findings/functions real-repos/pictex`
- Findings: `208`
- Chunks reviewed: `./scripts/chunks/Chunk_1_25.txt`, `Chunk_26_50.txt`, `Chunk_51_75.txt`, `Chunk_76_100.txt`, `Chunk_101_125.txt`, `Chunk_126_150.txt`, `Chunk_151_175.txt`, `Chunk_176_200.txt`, `Chunk_201_208.txt`
- Function contexts reviewed: `./scripts/findings/functions/<id>.txt` for every proposed false positive (2, 3, 19, 23, 39, 41-208 excluding 56, 57, 70, 77, 94, 113, 150, 167); enclosing source read when the exported context was insufficient

## Audit checklist

- [x] Read every assigned chunk under `./scripts/chunks`.
- [x] Read `./scripts/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 166 | 2, 3, 19, 23, 39, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 71, 72, 73, 74, 75, 76, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207, 208 |
| True positive | 42 | 1, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 20, 21, 22, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 40, 56, 57, 70, 77, 94, 113, 150 |
| Uncertain | 0 | — |

## False positives

### [ ] Finding `2` — BP-PY-7

- Function context: `./scripts/findings/functions/2.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/examples/code_to_image/code_to_image.py:32:16`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
CODE_SNIPPET = open("code_snippet.py", encoding="utf-8").read()
```

Why this is a false positive: The call is a temporary `open(...).read()` one-liner whose file object becomes unreferenced as soon as `.read()` returns and is closed immediately by refcounting; the rule's own detection notes say to ignore such one-liners ("Ignore temporary open().read() one-liners if closed by GC policy").

Checklist evidence: the `open` result is consumed inline by `.read()` in the same expression, so no open file handle outlives the statement.

### [ ] Finding `3` — PERF-PY-26

- Function context: `./scripts/findings/functions/3.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/examples/code_to_image/code_to_image.py:69:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
for i, line in enumerate(CODE_SNIPPET.strip().split('\n'), 1):
    line_number = Text(f"{i: >2}").color(COLORS["line_number"]).margin(0, 15, 0, 0)
    parsed_code_row = parse_python_line(line)
```

Why this is a false positive: The parse sits in a module-level loop of a one-shot example script that executes exactly once per process, and each call lexes a distinct line of a static snippet — there is no repeated decode of the same payload on a request/job hot path. The rule's detection notes explicitly say to "Suppress one-shot CLI tools".

Checklist evidence: the decode is executed once per script run on distinct inputs; there is no handler/job loop and no reused payload, so no cache is warranted.

### [ ] Finding `19` — BP-PY-7

- Function context: `./scripts/findings/functions/19.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/src/pictex/__skia_init.py:27:24`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
devnull_fd = os.open(os.devnull, os.O_WRONLY)
os.dup2(devnull_fd, original_stderr_fd)
...
    if 'devnull_fd' in locals():
        os.close(devnull_fd)
```

Why this is a false positive: `os.open` is the low-level file-descriptor API, not the `open` builtin or `Path.open` that the rule matches ("Call to open/Path.open"); it returns an integer fd (not a file object/context manager) that is explicitly released via `os.close` in the same function's `finally`.

Checklist evidence: the callee is `os.open`, which is neither `open` nor `Path.open`, and the returned fd is deterministically closed.

### [ ] Finding `23` — BP-PY-7

- Function context: `./scripts/findings/functions/23.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/src/pictex/models/public/background.py:22:46`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
self._skia_image = skia.Image.open(self.path)
```

Why this is a false positive: `skia.Image.open` is a third-party method of the skia library, not the `open` builtin or `Path.open`; it decodes the image into an in-memory `skia.Image` (no Python file handle survives the call, so the "open without `with`" resource-leak condition does not apply).

Checklist evidence: the callee is `skia.Image.open`, which the rule's condition (`open`/`Path.open`) does not cover.

### [ ] Finding `39` — BP-PY-46

- Function context: `./scripts/findings/functions/39.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/src/pictex/vector_image.py:109:13`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
    def __str__(self) -> str:
        """Returns the SVG content as a string.
        ...
        Example:
            vector_image = canvas.render_as_svg("Hi")
            print(vector_image)
            '<svg>...</svg>'
        """
        return self._svg_content
```

Why this is a false positive: The flagged `print(vector_image)` line is inside the `__str__` docstring — documentation text, not executable code; the module contains no actual `print` statement, so the rule's condition ("print used in library code") is not satisfied.

Checklist evidence: the `print(` token occurs inside a multi-line string literal (the docstring between `"""` delimiters).

### [ ] Finding `167` — CWE-367

- Function context: `./scripts/findings/functions/167.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_svg.py:51:16`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
font_basename = os.path.basename(STATIC_FONT_PATH)
copied_font_path = os.path.join(fonts_dir, font_basename)
assert os.path.exists(copied_font_path), f"Font file {font_basename} should be copied"

with open(output_path, 'r', encoding='utf-8') as f:
    saved_svg = f.read()
```

Why this is a false positive: The trigger regex matched `os.path.exists(...)` followed by any `open(` within 300 bytes, but the existence check (`copied_font_path`) and the later `open()` (`output_path`) refer to different files, so there is no check-then-use of the same resource. The check is additionally a test assertion (the statement's purpose), not a guard preceding a use of the checked path.

Checklist evidence: the checked path and the used path are distinct (`copied_font_path` vs `output_path`), so the TOCTOU condition (same resource checked, then used) is not satisfied.

### BP-PY-41 — pytest assert With Side Effects Only (160 findings)

### [ ] Finding `41` — BP-PY-41

- Function context: `./scripts/findings/functions/41.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_aspect_ratio.py:3:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_aspect_ratio_16_9_with_width(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `42` — BP-PY-41

- Function context: `./scripts/findings/functions/42.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_aspect_ratio.py:18:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_aspect_ratio_1_1_square(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `43` — BP-PY-41

- Function context: `./scripts/findings/functions/43.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_aspect_ratio.py:32:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_aspect_ratio_9_16_portrait(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `44` — BP-PY-41

- Function context: `./scripts/findings/functions/44.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_aspect_ratio.py:47:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_aspect_ratio_string_format(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `45` — BP-PY-41

- Function context: `./scripts/findings/functions/45.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_aspect_ratio.py:61:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_aspect_ratio_with_multiple_elements(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `46` — BP-PY-41

- Function context: `./scripts/findings/functions/46.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_aspect_ratio.py:77:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_aspect_ratio_golden_ratio(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `47` — BP-PY-41

- Function context: `./scripts/findings/functions/47.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_aspect_ratio.py:91:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_aspect_ratio_with_flex_grow(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `48` — BP-PY-41

- Function context: `./scripts/findings/functions/48.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_background.py:4:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_render_with_solid_background(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `49` — BP-PY-41

- Function context: `./scripts/findings/functions/49.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_background.py:21:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_render_with_gradient_background(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `50` — BP-PY-41

- Function context: `./scripts/findings/functions/50.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_background.py:44:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_background_without_padding(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `51` — BP-PY-41

- Function context: `./scripts/findings/functions/51.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_background.py:61:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_background_image_with_contain_mode(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `52` — BP-PY-41

- Function context: `./scripts/findings/functions/52.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_background.py:79:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_background_image_with_tile_mode(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `53` — BP-PY-41

- Function context: `./scripts/findings/functions/53.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_border.py:3:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_border_style_dashed(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `54` — BP-PY-41

- Function context: `./scripts/findings/functions/54.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_border.py:22:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_border_style_dotted(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `55` — BP-PY-41

- Function context: `./scripts/findings/functions/55.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_border.py:41:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_border_with_gradient_color(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `58` — BP-PY-41

- Function context: `./scripts/findings/functions/58.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_decorations.py:5:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_render_with_underline(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `59` — BP-PY-41

- Function context: `./scripts/findings/functions/59.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_decorations.py:20:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_render_with_strikethrough_custom_color(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `60` — BP-PY-41

- Function context: `./scripts/findings/functions/60.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_decorations.py:35:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_render_with_multiple_decorations(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `61` — BP-PY-41

- Function context: `./scripts/findings/functions/61.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_decorations.py:51:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_render_with_gradient_decoration(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `62` — BP-PY-41

- Function context: `./scripts/findings/functions/62.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_direction.py:3:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_rtl_text_rendering(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `63` — BP-PY-41

- Function context: `./scripts/findings/functions/63.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_direction.py:19:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_ltr_text_rendering(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `64` — BP-PY-41

- Function context: `./scripts/findings/functions/64.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_direction.py:34:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_direction_inheritance(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `65` — BP-PY-41

- Function context: `./scripts/findings/functions/65.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_direction.py:54:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_flex_row_direction_rtl(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `66` — BP-PY-41

- Function context: `./scripts/findings/functions/66.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_direction.py:76:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_mixed_direction_explicit(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `67` — BP-PY-41

- Function context: `./scripts/findings/functions/67.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_direction.py:96:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_bidi_punctuation(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `68` — BP-PY-41

- Function context: `./scripts/findings/functions/68.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_direction.py:110:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_automatic_text_align_resolution(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `69` — BP-PY-41

- Function context: `./scripts/findings/functions/69.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_direction.py:124:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_complex_mixed_bidi(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `71` — BP-PY-41

- Function context: `./scripts/findings/functions/71.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_flex_properties.py:4:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_flex_grow_basic(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `72` — BP-PY-41

- Function context: `./scripts/findings/functions/72.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_flex_properties.py:18:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_flex_shrink_prevents_shrinking(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `73` — BP-PY-41

- Function context: `./scripts/findings/functions/73.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_flex_properties.py:31:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_align_self_overrides_container(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `74` — BP-PY-41

- Function context: `./scripts/findings/functions/74.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_flex_properties.py:46:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_flex_wrap_creates_grid(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `75` — BP-PY-41

- Function context: `./scripts/findings/functions/75.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_flex_properties.py:61:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_flex_wrap_with_align_items(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `76` — BP-PY-41

- Function context: `./scripts/findings/functions/76.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_flex_properties.py:79:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_combined_flex_properties(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `78` — BP-PY-41

- Function context: `./scripts/findings/functions/78.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_fonts.py:5:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_render_with_custom_static_font(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `79` — BP-PY-41

- Function context: `./scripts/findings/functions/79.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_fonts.py:18:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_render_with_variable_font_weight(file_regression, render_engine, weight, expected_style):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `80` — BP-PY-41

- Function context: `./scripts/findings/functions/80.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_fonts.py:31:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_render_with_font_fallback_for_emoji(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `81` — BP-PY-41

- Function context: `./scripts/findings/functions/81.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_fonts.py:46:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_render_with_system_font_fallback(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `82` — BP-PY-41

- Function context: `./scripts/findings/functions/82.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_fonts.py:66:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_render_basic_text_and_alignment(file_regression, render_engine, text, align):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `83` — BP-PY-41

- Function context: `./scripts/findings/functions/83.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_fonts.py:73:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_render_with_default_font(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `84` — BP-PY-41

- Function context: `./scripts/findings/functions/84.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_fonts.py:86:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_render_with_invalid_fonts(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `85` — BP-PY-41

- Function context: `./scripts/findings/functions/85.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_gradients.py:3:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_linear_gradient_on_text_fill(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `86` — BP-PY-41

- Function context: `./scripts/findings/functions/86.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_gradients.py:23:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_linear_gradient_direction_vertical(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `87` — BP-PY-41

- Function context: `./scripts/findings/functions/87.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_gradients.py:43:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_linear_gradient_with_custom_stops(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `88` — BP-PY-41

- Function context: `./scripts/findings/functions/88.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_gradients.py:65:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_radial_gradient_on_background(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `89` — BP-PY-41

- Function context: `./scripts/findings/functions/89.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_gradients.py:87:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_radial_gradient_with_custom_center(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `90` — BP-PY-41

- Function context: `./scripts/findings/functions/90.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_gradients.py:110:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_sweep_gradient_color_wheel(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `91` — BP-PY-41

- Function context: `./scripts/findings/functions/91.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_gradients.py:131:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_sweep_gradient_with_stops(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `92` — BP-PY-41

- Function context: `./scripts/findings/functions/92.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_gradients.py:153:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_two_point_conical_gradient(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `93` — BP-PY-41

- Function context: `./scripts/findings/functions/93.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_gradients.py:178:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_two_point_conical_gradient_spotlight(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `95` — BP-PY-41

- Function context: `./scripts/findings/functions/95.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_integration.py:4:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_kitchen_sink_all_features_combined(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `96` — BP-PY-41

- Function context: `./scripts/findings/functions/96.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_layout.py:16:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_row_default_layout(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `97` — BP-PY-41

- Function context: `./scripts/findings/functions/97.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_layout.py:25:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_row_horizontal_distribution(file_regression, render_engine, mode):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `98` — BP-PY-41

- Function context: `./scripts/findings/functions/98.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_layout.py:37:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_row_vertical_alignment(file_regression, render_engine, mode):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `99` — BP-PY-41

- Function context: `./scripts/findings/functions/99.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_layout.py:48:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_row_with_gap(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `100` — BP-PY-41

- Function context: `./scripts/findings/functions/100.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_layout.py:54:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_row_with_gap_and_distribution(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `101` — BP-PY-41

- Function context: `./scripts/findings/functions/101.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_layout.py:66:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_column_default_layout(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `102` — BP-PY-41

- Function context: `./scripts/findings/functions/102.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_layout.py:75:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_column_vertical_distribution(file_regression, render_engine, mode):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `103` — BP-PY-41

- Function context: `./scripts/findings/functions/103.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_layout.py:87:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_column_horizontal_alignment(file_regression, render_engine, mode):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `104` — BP-PY-41

- Function context: `./scripts/findings/functions/104.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_layout.py:98:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_column_with_gap(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `105` — BP-PY-41

- Function context: `./scripts/findings/functions/105.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_layout.py:104:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_column_with_gap_and_distribution(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `106` — BP-PY-41

- Function context: `./scripts/findings/functions/106.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_letter_spacing.py:79:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_letter_spacing_percent(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `107` — BP-PY-41

- Function context: `./scripts/findings/functions/107.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_letter_spacing.py:93:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_letter_spacing_comparison(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `108` — BP-PY-41

- Function context: `./scripts/findings/functions/108.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_letter_spacing.py:116:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_letter_spacing_override_inheritance(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `109` — BP-PY-41

- Function context: `./scripts/findings/functions/109.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_letter_spacing.py:138:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_letter_spacing_multiline(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `110` — BP-PY-41

- Function context: `./scripts/findings/functions/110.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_letter_spacing.py:152:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_letter_spacing_with_multirun_text(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `111` — BP-PY-41

- Function context: `./scripts/findings/functions/111.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_letter_spacing.py:170:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_letter_spacing_disables_optional_ligatures(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `112` — BP-PY-41

- Function context: `./scripts/findings/functions/112.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_letter_spacing.py:192:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_letter_spacing_is_ignored_on_arabic_script(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `114` — BP-PY-41

- Function context: `./scripts/findings/functions/114.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_line_height.py:4:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_line_height_auto_default(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `115` — BP-PY-41

- Function context: `./scripts/findings/functions/115.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_line_height.py:21:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_line_height_explicit_multiplier(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `116` — BP-PY-41

- Function context: `./scripts/findings/functions/116.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_outline.py:4:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_render_with_basic_outline(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `117` — BP-PY-41

- Function context: `./scripts/findings/functions/117.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_outline.py:20:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_render_with_gradient_outline(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `118` — BP-PY-41

- Function context: `./scripts/findings/functions/118.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_outline.py:38:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_outline_without_fill(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `119` — BP-PY-41

- Function context: `./scripts/findings/functions/119.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_outline.py:54:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_stroke_mode_center(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `120` — BP-PY-41

- Function context: `./scripts/findings/functions/120.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_outline.py:70:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_stroke_mode_outline(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `121` — BP-PY-41

- Function context: `./scripts/findings/functions/121.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_outline.py:86:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_stroke_mode_inline(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `122` — BP-PY-41

- Function context: `./scripts/findings/functions/122.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_overflow.py:14:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_overflow_hidden_row_clips_children(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `123` — BP-PY-41

- Function context: `./scripts/findings/functions/123.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_overflow.py:32:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_overflow_hidden_column_clips_children(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `124` — BP-PY-41

- Function context: `./scripts/findings/functions/124.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_overflow.py:50:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_overflow_visible_row_does_not_clip(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `125` — BP-PY-41

- Function context: `./scripts/findings/functions/125.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_overflow.py:72:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_overflow_hidden_text_clips_content(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `126` — BP-PY-41

- Function context: `./scripts/findings/functions/126.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_overflow.py:87:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_overflow_visible_text_does_not_clip(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `127` — BP-PY-41

- Function context: `./scripts/findings/functions/127.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_overflow.py:106:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_overflow_hidden_preserves_border(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `128` — BP-PY-41

- Function context: `./scripts/findings/functions/128.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_overflow.py:123:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_overflow_hidden_preserves_border_radius(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `129` — BP-PY-41

- Function context: `./scripts/findings/functions/129.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_positioning.py:4:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_place_absolute_is_removed_from_row_flow(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `130` — BP-PY-41

- Function context: `./scripts/findings/functions/130.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_positioning.py:27:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_place_absolute_is_removed_from_column_flow(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `131` — BP-PY-41

- Function context: `./scripts/findings/functions/131.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_positioning.py:50:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_place_anchors_and_offsets(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `132` — BP-PY-41

- Function context: `./scripts/findings/functions/132.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_positioning.py:76:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_place_with_percentages(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `133` — BP-PY-41

- Function context: `./scripts/findings/functions/133.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_positioning.py:95:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_place_with_mixed_anchors(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `134` — BP-PY-41

- Function context: `./scripts/findings/functions/134.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_positioning.py:118:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_place_with_nested_children(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `135` — BP-PY-41

- Function context: `./scripts/findings/functions/135.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_positioning.py:138:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_container_with_place(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `136` — BP-PY-41

- Function context: `./scripts/findings/functions/136.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_positioning.py:147:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_absolute_position_with_inset(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `137` — BP-PY-41

- Function context: `./scripts/findings/functions/137.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_positioning.py:162:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_translate_for_centering(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `138` — BP-PY-41

- Function context: `./scripts/findings/functions/138.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_positioning.py:179:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_translate_with_nested_children(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `139` — BP-PY-41

- Function context: `./scripts/findings/functions/139.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_positioning.py:197:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_fixed_position_canvas_relative(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `140` — BP-PY-41

- Function context: `./scripts/findings/functions/140.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_positioning.py:214:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_fixed_vs_absolute_positioning(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `141` — BP-PY-41

- Function context: `./scripts/findings/functions/141.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_positioning.py:227:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_relative_position_offset(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `142` — BP-PY-41

- Function context: `./scripts/findings/functions/142.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_positioning.py:244:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_fixed_with_canvas_box_model(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `143` — BP-PY-41

- Function context: `./scripts/findings/functions/143.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_positioning.py:264:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_multiple_fixed_elements(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `144` — BP-PY-41

- Function context: `./scripts/findings/functions/144.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_shadows.py:3:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_render_with_simple_text_shadow(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `145` — BP-PY-41

- Function context: `./scripts/findings/functions/145.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_shadows.py:19:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_render_with_simple_box_shadow(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `146` — BP-PY-41

- Function context: `./scripts/findings/functions/146.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_shadows.py:37:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_render_with_multiple_text_shadows(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `147` — BP-PY-41

- Function context: `./scripts/findings/functions/147.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_shadows.py:58:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_render_with_multiple_box_shadows(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `148` — BP-PY-41

- Function context: `./scripts/findings/functions/148.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_shadows.py:79:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_text_and_box_shadows_together(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `149` — BP-PY-41

- Function context: `./scripts/findings/functions/149.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_shadows.py:99:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_hard_shadow_without_blur(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `151` — BP-PY-41

- Function context: `./scripts/findings/functions/151.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_size_constraints.py:3:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_max_width_prevents_overflow(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `152` — BP-PY-41

- Function context: `./scripts/findings/functions/152.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_size_constraints.py:24:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_min_width_prevents_collapse(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `153` — BP-PY-41

- Function context: `./scripts/findings/functions/153.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_size_constraints.py:40:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_max_height_limits_vertical_growth(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `154` — BP-PY-41

- Function context: `./scripts/findings/functions/154.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_size_constraints.py:59:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_min_height_ensures_minimum_space(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `155` — BP-PY-41

- Function context: `./scripts/findings/functions/155.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_size_constraints.py:76:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_max_width_with_flex_grow(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `156` — BP-PY-41

- Function context: `./scripts/findings/functions/156.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_size_constraints.py:102:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_min_width_with_flex_shrink(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `157` — BP-PY-41

- Function context: `./scripts/findings/functions/157.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_size_constraints.py:129:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_size_constraints_with_percentages(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `158` — BP-PY-41

- Function context: `./scripts/findings/functions/158.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_sizing.py:5:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_size_absolute(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `159` — BP-PY-41

- Function context: `./scripts/findings/functions/159.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_sizing.py:18:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_size_fit_content_on_text(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `160` — BP-PY-41

- Function context: `./scripts/findings/functions/160.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_sizing.py:32:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_size_fit_content_on_row(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `161` — BP-PY-41

- Function context: `./scripts/findings/functions/161.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_sizing.py:48:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_size_fit_background_image(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `162` — BP-PY-41

- Function context: `./scripts/findings/functions/162.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_sizing.py:65:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_size_percent_width_and_fixed_height(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `163` — BP-PY-41

- Function context: `./scripts/findings/functions/163.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_sizing.py:84:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_size_percent_height(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `164` — BP-PY-41

- Function context: `./scripts/findings/functions/164.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_sizing.py:102:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_size_flex_grow_single_child(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `165` — BP-PY-41

- Function context: `./scripts/findings/functions/165.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_sizing.py:124:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_size_flex_grow_multiple_children(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `166` — BP-PY-41

- Function context: `./scripts/findings/functions/166.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_sizing.py:148:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_size_flex_grow_parent_with_percent_child(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `168` — BP-PY-41

- Function context: `./scripts/findings/functions/168.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_box_edge.py:43:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_text_box_edge_all_modes(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `169` — BP-PY-41

- Function context: `./scripts/findings/functions/169.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_box_edge.py:68:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_text_box_edge_glyphs_with_line_height(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `170` — BP-PY-41

- Function context: `./scripts/findings/functions/170.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_box_edge.py:89:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_text_box_edge_with_emojis(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `171` — BP-PY-41

- Function context: `./scripts/findings/functions/171.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_shaping.py:5:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_ligature_rendering(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `172` — BP-PY-41

- Function context: `./scripts/findings/functions/172.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_shaping.py:22:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_kerning_support(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `173` — BP-PY-41

- Function context: `./scripts/findings/functions/173.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_shaping.py:40:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_complex_emoji_rendering(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `174` — BP-PY-41

- Function context: `./scripts/findings/functions/174.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_shaping.py:56:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_arabic_text_shaping(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `175` — BP-PY-41

- Function context: `./scripts/findings/functions/175.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_shaping.py:72:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_arabic_text_shaping_with_diacritics(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `176` — BP-PY-41

- Function context: `./scripts/findings/functions/176.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_shaping.py:89:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_family_emoji_with_zwj(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `177` — BP-PY-41

- Function context: `./scripts/findings/functions/177.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_shaping.py:115:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_emoji_with_modifiers(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `178` — BP-PY-41

- Function context: `./scripts/findings/functions/178.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_shaping.py:135:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_multi_font_metrics_and_decorations(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `179` — BP-PY-41

- Function context: `./scripts/findings/functions/179.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_shaping.py:156:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_contextual_joining_across_font_runs(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `180` — BP-PY-41

- Function context: `./scripts/findings/functions/180.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_spans.py:3:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_span_basic_color(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `181` — BP-PY-41

- Function context: `./scripts/findings/functions/181.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_spans.py:16:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_span_nested_inheritance(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `182` — BP-PY-41

- Function context: `./scripts/findings/functions/182.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_spans.py:30:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_span_gradient_across_subspans(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `183` — BP-PY-41

- Function context: `./scripts/findings/functions/183.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_spans.py:46:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_span_gradient_with_implicit_text(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `184` — BP-PY-41

- Function context: `./scripts/findings/functions/184.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_spans.py:60:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_span_multiline_gradient(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `185` — BP-PY-41

- Function context: `./scripts/findings/functions/185.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_spans.py:76:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_span_continuous_underline(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `186` — BP-PY-41

- Function context: `./scripts/findings/functions/186.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_spans.py:92:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_span_distinct_underlines(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `187` — BP-PY-41

- Function context: `./scripts/findings/functions/187.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_spans.py:104:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_span_text_stroke_and_shadow(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `188` — BP-PY-41

- Function context: `./scripts/findings/functions/188.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_spans.py:121:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_span_complex_nesting(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `189` — BP-PY-41

- Function context: `./scripts/findings/functions/189.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_spans.py:143:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_span_letter_spacing(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `190` — BP-PY-41

- Function context: `./scripts/findings/functions/190.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_spans.py:156:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_text_block_multiline_gradient(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `191` — BP-PY-41

- Function context: `./scripts/findings/functions/191.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_spans.py:170:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_span_wrap_multiline_gradient_with_width(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `192` — BP-PY-41

- Function context: `./scripts/findings/functions/192.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_spans.py:191:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_span_inline_emoji_gradient(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `193` — BP-PY-41

- Function context: `./scripts/findings/functions/193.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_spans.py:209:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_span_inherited_gradient_underline_with_runs(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `194` — BP-PY-41

- Function context: `./scripts/findings/functions/194.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_spans.py:228:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_span_arabic_contextual_joining(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `195` — BP-PY-41

- Function context: `./scripts/findings/functions/195.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_spans.py:244:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_span_variable_font_sizes(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `196` — BP-PY-41

- Function context: `./scripts/findings/functions/196.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_spans.py:263:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_span_bidi_reordering_with_gradient(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `197` — BP-PY-41

- Function context: `./scripts/findings/functions/197.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_wrap.py:8:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_text_wrap_normal_vs_nowrap(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `198` — BP-PY-41

- Function context: `./scripts/findings/functions/198.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_wrap.py:26:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_text_wrap_inherited_from_parent(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `199` — BP-PY-41

- Function context: `./scripts/findings/functions/199.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_wrap.py:51:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_text_wrap_override_inheritance(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `200` — BP-PY-41

- Function context: `./scripts/findings/functions/200.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_wrap.py:69:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_text_wrap_different_widths(file_regression, render_engine, width):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `201` — BP-PY-41

- Function context: `./scripts/findings/functions/201.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_wrap.py:81:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_text_wrap_with_multiline_text(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `202` — BP-PY-41

- Function context: `./scripts/findings/functions/202.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_wrap.py:101:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_text_wrap_with_percentage_width(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `203` — BP-PY-41

- Function context: `./scripts/findings/functions/203.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_wrap.py:115:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_text_wrap_nested_containers(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `204` — BP-PY-41

- Function context: `./scripts/findings/functions/204.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_wrap.py:132:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_text_wrap_with_styling(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `205` — BP-PY-41

- Function context: `./scripts/findings/functions/205.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_wrap.py:150:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_no_false_wrap_without_width_constraint(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `206` — BP-PY-41

- Function context: `./scripts/findings/functions/206.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_wrap.py:172:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_cjk_text_wraps_within_fixed_width(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `207` — BP-PY-41

- Function context: `./scripts/findings/functions/207.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_wrap.py:193:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_two_sibling_texts_using_width_limited_by_ancestor(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

### [ ] Finding `208` — BP-PY-41

- Function context: `./scripts/findings/functions/208.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pictex/tests/test_text_wrap.py:209:1`
- Checklist pattern: based the decision on the rule condition and the shown source

Source excerpt:

```
def test_empty_lines(file_regression, render_engine):
    ...
    check_func(file_regression, image)
```

Why this is a false positive: The test is not a placeholder — its final statement calls `check_func(file_regression, image)`, the `render_engine` fixture's check helper (conftest.py) that performs `file_regression.check(...)` (pytest-regressions baseline assertion), so the test verifies output and fails on mismatch. The BP-PY-41 heuristic only recognizes `assert`/`pytest.raises`/`assert*` callees and misses this helper call.

Checklist evidence: the rule condition requires no assertions; the shown source contains an assertion-equivalent regression check call (`check_func(file_regression, image)`) in the test body.

## True positives

### BP-PY-5 — Wildcard Import (19)

| Finding | Source | Reason |
| --- | --- | --- |
| 1 | examples/code_to_image/code_to_image.py:1 | `from pictex import *` — wildcard ImportFrom in non-`__init__.py` |
| 4 | examples/examples_made_by_ai/certificate_template.py:6 | `from pictex import *` — wildcard ImportFrom |
| 5 | examples/examples_made_by_ai/dashboard_metrics.py:6 | `from pictex import *` — wildcard ImportFrom |
| 6 | examples/examples_made_by_ai/infographic_stats.py:1 | `from pictex import *` — wildcard ImportFrom |
| 7 | examples/examples_made_by_ai/product_showcase.py:6 | `from pictex import *` — wildcard ImportFrom |
| 8 | examples/github_card/github_card.py:9 | `from pictex import *` — wildcard ImportFrom |
| 12 | examples/table/table.py:1 | `from pictex import *` — wildcard ImportFrom |
| 13 | examples/tweet_card/tweet_card.py:1 | `from pictex import *` — wildcard ImportFrom |
| 20 | src/pictex/builders/canvas.py:6 | `from ..models import *` — wildcard ImportFrom in library code |
| 21 | src/pictex/builders/inline_styleable.py:4 | `from ..models import *` — wildcard ImportFrom in library code |
| 22 | src/pictex/builders/stylable.py:4 | `from ..models import *` — wildcard ImportFrom in library code |
| 40 | tests/test_aspect_ratio.py:1 | `from pictex import *` — wildcard ImportFrom |
| 56 | tests/test_builders_api.py:2 | `from pictex import *` — wildcard ImportFrom |
| 57 | tests/test_crop.py:1 | `from pictex import *` — wildcard ImportFrom |
| 70 | tests/test_flex_properties.py:1 | `from pictex import *` — wildcard ImportFrom |
| 77 | tests/test_fonts.py:1 | `from pictex import *` — wildcard ImportFrom |
| 94 | tests/test_integration.py:1 | `from pictex import *` — wildcard ImportFrom |
| 113 | tests/test_line_height.py:1 | `from pictex import *` — wildcard ImportFrom |
| 150 | tests/test_size_constraints.py:1 | `from pictex import *` — wildcard ImportFrom |

### BP-PY-14 — requests Without Timeout (2)

| Finding | Source | Reason |
| --- | --- | --- |
| 9 | examples/github_card/github_card.py:22 | `requests.get(base_url)` without `timeout=` |
| 10 | examples/github_card/github_card.py:23 | `requests.get(f"{base_url}/languages")` without `timeout=` |

### BP-PY-46 — print Debugging In Library Code (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 11 | examples/github_card/github_card.py:143 | module-level `print(...)` not under `__main__` guard, not a test |

### BP-PY-1 — Bare Except Clause (5)

| Finding | Source | Reason |
| --- | --- | --- |
| 14 | src/pictex/__skia_init.py:17 | `except Exception:` with `pass` suite |
| 35 | src/pictex/text/typeface_loader.py:31 | bare `except:` (no type) |
| 36 | src/pictex/text/typeface_loader.py:45 | bare `except:` (no type) |
| 37 | src/pictex/text/typeface_loader.py:66 | bare `except:` (no type) |
| 38 | src/pictex/utils/font.py:6 | bare `except:` (no type) |

### BP-PY-2 — Except Pass (2)

| Finding | Source | Reason |
| --- | --- | --- |
| 15 | src/pictex/__skia_init.py:17 | except suite is solely `pass` |
| 33 | src/pictex/text/harfbuzz_shaper.py:216 | `except (AttributeError, Exception):` suite is solely `pass` |

### BP-PY-3 — Raise Generic Exception (2)

| Finding | Source | Reason |
| --- | --- | --- |
| 29 | src/pictex/renderer/renderer.py:25 | `raise Exception(...)` in library code |
| 31 | src/pictex/renderer/renderer.py:33 | `raise Exception("Failed to create surface")` |

### CWE-390 — Detection of Error Condition Without Action (2)

| Finding | Source | Reason |
| --- | --- | --- |
| 16 | src/pictex/__skia_init.py:17 | exception detected, handler body is only `pass` |
| 34 | src/pictex/text/harfbuzz_shaper.py:216 | exception detected, handler body is only `pass` |

### CWE-396 — Declaration of Catch for Generic Exception (2)

| Finding | Source | Reason |
| --- | --- | --- |
| 17 | src/pictex/__skia_init.py:17 | `except Exception:` handler |
| 24 | src/pictex/models/public/background.py:23 | `except Exception:` handler |

### CWE-397 — Declaration of Throws for Generic Exception (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 30 | src/pictex/renderer/renderer.py:25 | `raise Exception(...)` directly raised |

### CWE-1071 — Empty Code Block (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 18 | src/pictex/__skia_init.py:17 | exception handler body is empty (`pass` only) |

### CWE-1121 — Excessive McCabe Cyclomatic Complexity (4)

| Finding | Source | Reason |
| --- | --- | --- |
| 25 | src/pictex/models/public/linear_gradient.py:67 | `__post_init__` has 18 control-flow branches (threshold 12) |
| 26 | src/pictex/models/public/radial_gradient.py:68 | `__post_init__` has 16 control-flow branches (threshold 12) |
| 27 | src/pictex/models/public/sweep_gradient.py:52 | `__post_init__` has 15 control-flow branches (threshold 12) |
| 28 | src/pictex/models/public/two_point_conical_gradient.py:83 | `__post_init__` has 20 control-flow branches (threshold 12) |

### CWE-1046 — Creation of Immutable Text Using String Concatenation (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 32 | src/pictex/renderer/vector_image_processor.py:122 | `css += f"""@font-face ..."""` inside `for typeface in typefaces:` loop |

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `./scripts/chunks`
- Function evidence: `./scripts/findings/functions`
- Validation: `git diff --check` — `pass`
