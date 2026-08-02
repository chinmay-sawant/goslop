# False-positive audit — pyauto-desktop

## Run metadata

```yaml
timestamp: 2026-08-02T07:40:34Z
repository: pyauto-desktop
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pyauto-desktop
branch: main
commit: 2448ba4a843cde40697a289895b7cb2e399d5afa
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pyauto-desktop
chunk_path: ./scripts/chunks
function_context_path: ./scripts/findings/functions
```

## Scan evidence

- Build command: `n/a` (prebuilt `./bin/goslop` in goslop repo root)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir real-repos/pyauto-desktop/scripts/chunks -context-dir real-repos/pyauto-desktop/scripts/findings/functions real-repos/pyauto-desktop`
- Findings: `92`
- Chunks reviewed: `./scripts/chunks/Chunk_1_25.txt`, `./scripts/chunks/Chunk_26_50.txt`, `./scripts/chunks/Chunk_51_75.txt`, `./scripts/chunks/Chunk_76_92.txt`
- Function contexts reviewed: `./scripts/findings/functions/1.txt`, `63.txt`, `66.txt`, `74.txt` (all proposed false positives); full enclosing source read for all FP candidates and for every CWE-1121 / BP-PY-7 / BP-PY-12 / BP-PY-36 / PERF-PY-26 / CWE-1341 trigger.

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
| False positive | 4 | 1, 63, 66, 74 |
| True positive | 88 | 2–62 (except none), 64, 65, 67–73, 75–92 |
| Uncertain | 0 | — |

## False positives

### [ ] Finding `1` — `CWE-1341`

- Function context: `./scripts/findings/functions/1.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pyauto-desktop/pyauto_desktop/capture_tool.py:91:13`
- Checklist pattern: rule condition (two releases of the same resource handle) is not met by the shown source.

Source excerpt:

```
        try:
            cropped_pixmap = self.original_pixmap.copy(phys_x, phys_y, phys_w, phys_h)
            self.snipped.emit(cropped_pixmap, phys_coords, self.target_screen)
            self.close()
        except Exception as e:
            print(f"Crop failed: {e}")
            self.closed.emit()
            self.close()
```

Why this is a false positive: the regex trigger matches two `self.close()` calls (lines 91 and 95), but they sit in mutually exclusive try/except branches so at most one executes per run, and `close()` is the QWidget method, not a resource-handle release.

Checklist evidence: `self.close()` is never invoked twice in a single execution path; `close()` on `self` (a `QWidget`) hides the window rather than releasing a handle.

### [ ] Finding `63` — `PERF-PY-26`

- Function context: `./scripts/findings/functions/63.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pyauto-desktop/pyauto_desktop/main.py:342:1`
- Checklist pattern: the decode is not on a hot path (no loop, no render/job path); the trigger only matches the function name.

Source excerpt:

```
    def process_loaded_image(self, path, mode='target'):
        try:
            img = Image.open(path)
            edited_img = self.open_editor(img)
            if not edited_img: return
```
```
        fname, _ = QFileDialog.getOpenFileName(self, f"Open {label} Image", "", "Images (*.png *.jpg *.bmp)")
        if fname:
            self.process_loaded_image(fname, mode)
```

Why this is a false positive: `Image.open` runs once per user-initiated image load (invoked only from the file-dialog handler `request_upload_image` and the drag-drop handler `handle_dropped_image`); the "hot path" signal comes solely from the function name containing "process".

Checklist evidence: the call site is not inside a loop and no render/job/per-frame path reaches it; a one-shot user action cannot satisfy "expensive decode/parse runs on a hot path".

### [ ] Finding `66` — `BP-PY-12`

- Function context: `./scripts/findings/functions/66.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pyauto-desktop/pyauto_desktop/main.py:384:19`
- Checklist pattern: the `exec` identifier is an attribute (method) reference, not the builtin exec.

Source excerpt:

```
    def open_editor(self, pil_img):
        editor = MagicWandEditor(pil_img, self)
        if editor.exec():
            return editor.get_result()
```

Why this is a false positive: `editor.exec()` is the PyQt `QDialog.exec()` modal-dialog method (runs the dialog's event loop); it is not the `exec` builtin and cannot execute arbitrary code.

Checklist evidence: the ident `exec` is preceded by a `.` attribute access on a `MagicWandEditor` instance, so the "eval/exec on dynamic input" condition is not satisfied.

### [ ] Finding `74` — `BP-PY-12`

- Function context: `./scripts/findings/functions/74.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pyauto-desktop/pyauto_desktop/main.py:890:9`
- Checklist pattern: the `exec` identifier is an attribute (method) reference, not the builtin exec.

Source excerpt:

```
def run_inspector():
    app = QApplication(sys.argv)
    window = MainWindow()
    window.show()
    app.exec()
```

Why this is a false positive: `app.exec()` is the PyQt `QApplication.exec()` application event-loop method, not the `exec` builtin; no dynamic code can be executed through it.

Checklist evidence: the ident `exec` is an attribute call on the `QApplication` instance, so the "eval/exec on dynamic input" condition is not satisfied.

## Uncertain findings

None.

## True positives

All remaining findings satisfy their rule's condition; tables per rule.

### BP-PY-1 (bare `except` / `except Exception`)

| ID | Source | Reason |
| --- | --- | --- |
| 2 | pyauto_desktop/capture_tool.py:92 | broad `except Exception as e`, failure only printed |
| 5 | pyauto_desktop/detection.py:119 | broad handler converts OCR error to text |
| 7 | pyauto_detection/detection.py:124 | broad handler, traceback only |
| 10 | pyauto_desktop/detection.py:236 | broad worker handler, no re-raise |
| 14 | pyauto_desktop/dpi_manager.py:20 | broad handler in DPI fallback chain |
| 16 | pyauto_desktop/dpi_manager.py:23 | broad handler in DPI fallback chain |
| 17 | pyauto_desktop/dpi_manager.py:26 | broad handler in DPI fallback chain |
| 21 | pyauto_desktop/dpi_manager.py:44 | broad handler, silent fallback |
| 23 | pyauto_desktop/editor.py:544 | broad handler, only printed |
| 26 | pyauto_desktop/functions.py:64 | broad handler around `mss.mss()` |
| 30 | pyauto_desktop/functions.py:107 | broad handler, silent pass |
| 34 | pyauto_desktop/functions.py:124 | broad handler, silent `return -1, -1` |
| 35 | pyauto_desktop/functions.py:173 | broad handler, only debug-printed |
| 41 | pyauto_desktop/functions.py:342 | bare `except:` swallows everything |
| 45 | pyauto_desktop/functions.py:659 | broad MSS capture handler |
| 49 | pyauto_desktop/functions.py:673 | broad ImageGrab fallback handler |
| 53 | pyauto_desktop/functions.py:843 | broad save handler |
| 57 | pyauto_desktop/main.py:27 | broad handler, only printed |
| 60 | pyauto_desktop/main.py:109 | broad handler, only printed |
| 65 | pyauto_desktop/main.py:372 | broad handler around image load |
| 68 | pyauto_desktop/main.py:493 | broad handler, silent pass |
| 71 | pyauto_desktop/main.py:692 | broad save handler |
| 72 | pyauto_desktop/main.py:873 | broad handler, only printed |
| 75 | pyauto_desktop/overlay.py:136 | broad paint handler |
| 79 | pyauto_desktop/text_recognition.py:148 | broad OCR handler |
| 82 | pyauto_desktop/window_control.py:20 | broad handler, silent None |
| 84 | pyauto_desktop/window_control.py:61 | broad handler, only printed |
| 87 | pyauto_desktop/window_control.py:77 | broad handler, only printed |
| 89 | pyauto_desktop/window_control.py:96 | broad handler, only printed |
| 91 | pyauto_desktop/window_control.py:149 | broad handler, only printed |

### CWE-396 (generic catch)

| ID | Source | Reason |
| --- | --- | --- |
| 3 | pyauto_desktop/capture_tool.py:92 | generic `Exception` handler |
| 6 | pyauto_desktop/detection.py:119 | generic handler around OCR |
| 15 | pyauto_desktop/dpi_manager.py:20 | generic handler |
| 24 | pyauto_desktop/editor.py:544 | generic handler |
| 27 | pyauto_desktop/functions.py:64 | generic handler |
| 58 | pyauto_desktop/main.py:27 | generic handler |
| 76 | pyauto_desktop/overlay.py:136 | generic handler |
| 80 | pyauto_desktop/text_recognition.py:148 | generic handler |
| 83 | pyauto_desktop/window_control.py:20 | generic handler |

### BP-PY-46 (print in library code)

| ID | Source | Reason |
| --- | --- | --- |
| 4 | pyauto_desktop/capture_tool.py:93 | `print` outside main guard / CLI entrypoints |
| 9 | pyauto_desktop/detection.py:144 | `print` in worker method |
| 11 | pyauto_desktop/detection.py:237 | `print` in worker method |
| 25 | pyauto_desktop/editor.py:545 | `print` in editor method |
| 28 | pyauto_desktop/functions.py:65 | `print` debug log |
| 29 | pyauto_desktop/functions.py:90 | `print` perf log |
| 36 | pyauto_desktop/functions.py:174 | `print` debug log |
| 38 | pyauto_desktop/functions.py:221 | `print` cache log |
| 42 | pyauto_desktop/functions.py:439 | `print` perf log |
| 43 | pyauto_desktop/functions.py:524 | `print` error log |
| 44 | pyauto_desktop/functions.py:649 | `print` error log |
| 46 | pyauto_desktop/functions.py:661 | `print` error log |
| 47 | pyauto_desktop/functions.py:662 | `print` error log |
| 48 | pyauto_desktop/functions.py:663 | `print` error log |
| 50 | pyauto_desktop/functions.py:674 | `print` error log |
| 51 | pyauto_desktop/functions.py:747 | `print` debug log |
| 52 | pyauto_desktop/functions.py:828 | `print` error log |
| 54 | pyauto_desktop/functions.py:844 | `print` error log |
| 55 | pyauto_desktop/functions.py:855 | `print` empty line |
| 56 | pyauto_desktop/functions.py:869 | `print` debug log |
| 59 | pyauto_desktop/main.py:28 | `print` outside main guard |
| 61 | pyauto_desktop/main.py:110 | `print` in slot |
| 73 | pyauto_desktop/main.py:874 | `print` error log |
| 77 | pyauto_desktop/overlay.py:137 | `print` in paint handler |
| 81 | pyauto_desktop/text_recognition.py:149 | `print` error log |
| 85 | pyauto_desktop/window_control.py:62 | `print` error log |
| 86 | pyauto_desktop/window_control.py:64 | `print` error log |
| 88 | pyauto_desktop/window_control.py:78 | `print` error log |
| 90 | pyauto_desktop/window_control.py:97 | `print` error log |
| 92 | pyauto_desktop/window_control.py:150 | `print` error log |

### BP-PY-2 (except pass)

| ID | Source | Reason |
| --- | --- | --- |
| 18 | pyauto_desktop/dpi_manager.py:26 | handler body is only `pass` |
| 31 | pyauto_desktop/functions.py:107 | handler body is only `pass` |
| 69 | pyauto_desktop/main.py:493 | handler body is only `pass` |

### CWE-390 (error condition without action)

| ID | Source | Reason |
| --- | --- | --- |
| 19 | pyauto_desktop/dpi_manager.py:26 | detected error, no action |
| 32 | pyauto_desktop/functions.py:107 | detected error, no action |
| 70 | pyauto_desktop/main.py:493 | detected error, no action |

### CWE-1071 (empty code block)

| ID | Source | Reason |
| --- | --- | --- |
| 20 | pyauto_desktop/dpi_manager.py:26 | handler contains only `pass` |
| 33 | pyauto_desktop/functions.py:107 | handler contains only `pass` |

### BP-PY-3 (raise generic exception)

| ID | Source | Reason |
| --- | --- | --- |
| 12 | pyauto_desktop/dpi_manager.py:19 | `raise Exception("V2 Failed")` |

### CWE-397 (throws generic exception)

| ID | Source | Reason |
| --- | --- | --- |
| 13 | pyauto_desktop/dpi_manager.py:19 | generic `Exception` raised directly |

### BP-PY-7 (open without context manager)

| ID | Source | Reason |
| --- | --- | --- |
| 37 | pyauto_desktop/functions.py:185 | `return Image.open(img)` — `.open(` without `with` |
| 39 | pyauto_desktop/functions.py:223 | `img = Image.open(image_path)` without `with` |
| 64 | pyauto_desktop/main.py:342 | `img = Image.open(path)` without `with` |

### BP-PY-36 (SQLAlchemy session not closed)

| ID | Source | Reason |
| --- | --- | --- |
| 67 | pyauto_desktop/main.py:476 | `session = Session(...)` in `update_live_preview`; no `session.close()` anywhere in scope (except handler is `pass`) |

### CWE-1121 (excessive cyclomatic complexity)

| ID | Source | Reason |
| --- | --- | --- |
| 8 | pyauto_desktop/detection.py:128 | `run_image_detection`: 13 branch keywords (`if`/`elif`/`for`/`except`) |
| 22 | pyauto_desktop/editor.py:160 | `mouseMoveEvent`: 12+ branch keywords |
| 40 | pyauto_desktop/functions.py:227 | `_process_needle_to_cv2`: 17 branch keywords |
| 62 | pyauto_desktop/main.py:196 | `on_snip_finished`: 12+ branch keywords |
| 78 | pyauto_desktop/text_recognition.py:18 | `preprocess_image`: 12 branch-keyword matches per the rule's counting (incl. `elif` and a `for` in a comprehension) |

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `./scripts/chunks` (4 chunk files, findings 1–92)
- Function evidence: `./scripts/findings/functions` (92 per-finding contexts)
- Validation: `git diff --check` — pass
