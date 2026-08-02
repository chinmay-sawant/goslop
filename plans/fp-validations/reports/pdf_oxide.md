# False-positive audit: pdf_oxide

## Run metadata

```yaml
timestamp: 2026-08-02T07:59:27Z
repository: pdf_oxide
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide
branch: main
commit: 10b87f153200cd5c4d4a4defee471757091e6559
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide
chunk_path: ./scripts/chunks
function_context_path: ./scripts/findings/functions
```

## Scan evidence

- Build command: `cargo build --release` (prebuilt `./bin/goslop`)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir real-repos/pdf_oxide/scripts/chunks -context-dir real-repos/pdf_oxide/scripts/findings/functions real-repos/pdf_oxide`
- Findings: `636`
- Chunks reviewed: `./scripts/chunks/Chunk_1_25.txt .. Chunk_626_636.txt` (all 26 chunk files)
- Function contexts reviewed: `./scripts/findings/functions/<finding-id>.txt` for every false positive and for the uncertain finding; enclosing source read where the exported context was insufficient (findings 19, 20, 141, 142, 144, 150, 276, 291, 403, 464, 488, 510, 562, 572, 592, 601, 604, 607, 632)

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
| False positive | 447 | 6, 7, 8, 9, 10, 13, 14, 15, 16, 17, 18, 19, 28, 29, 30, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 130, 137, 138, 141, 142, 143, 144, 149, 150, 151, 152, 153, 155, 158, 160, 161, 162, 163, 164, 166, 168, 169, 170, 171, 172, 174, 176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 193, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207, 208, 209, 213, 214, 215, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 258, 268, 269, 270, 273, 274, 275, 276, 277, 278, 279, 280, 281, 282, 283, 284, 285, 286, 288, 289, 290, 291, 292, 293, 294, 295, 296, 297, 298, 299, 300, 301, 302, 303, 305, 306, 307, 308, 309, 310, 311, 312, 313, 314, 315, 316, 317, 318, 319, 320, 321, 322, 323, 324, 325, 326, 327, 330, 331, 332, 333, 334, 335, 336, 337, 338, 339, 340, 341, 342, 343, 344, 345, 346, 347, 348, 349, 350, 351, 357, 362, 363, 365, 366, 367, 368, 369, 370, 371, 372, 373, 374, 375, 376, 377, 378, 379, 380, 381, 382, 383, 384, 385, 386, 387, 388, 389, 390, 391, 392, 393, 394, 395, 396, 397, 398, 399, 400, 401, 402, 404, 405, 406, 407, 408, 409, 412, 413, 414, 415, 416, 417, 418, 421, 425, 427, 428, 429, 430, 433, 434, 435, 436, 437, 438, 439, 442, 443, 444, 445, 446, 449, 453, 456, 457, 462, 463, 464, 465, 466, 467, 468, 469, 470, 473, 475, 476, 477, 478, 479, 480, 481, 484, 487, 488, 489, 490, 491, 494, 497, 498, 499, 501, 504, 505, 510, 511, 512, 513, 516, 517, 518, 519, 520, 521, 522, 523, 524, 525, 526, 527, 528, 529, 530, 531, 532, 533, 534, 535, 536, 537, 538, 539, 542, 543, 544, 545, 546, 547, 548, 549, 550, 551, 552, 553, 554, 555, 556, 557, 558, 559, 560, 561, 562, 572, 586, 587, 590, 592, 600, 601, 603, 604, 606, 607, 608, 609, 610, 611, 612, 613, 614 |
| True positive | 188 | 1, 2, 3, 4, 5, 11, 12, 20, 21, 22, 23, 24, 25, 26, 27, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 113, 128, 129, 131, 132, 133, 134, 135, 136, 139, 140, 145, 146, 147, 148, 154, 156, 157, 159, 165, 167, 173, 175, 191, 192, 194, 210, 211, 212, 216, 254, 255, 256, 257, 259, 260, 261, 262, 263, 264, 265, 266, 267, 271, 272, 287, 304, 328, 329, 352, 353, 354, 355, 356, 358, 359, 360, 361, 364, 403, 410, 411, 419, 420, 422, 423, 424, 426, 431, 432, 440, 441, 447, 448, 450, 451, 452, 454, 455, 458, 459, 460, 461, 471, 472, 474, 482, 483, 485, 486, 492, 493, 495, 496, 500, 502, 503, 506, 507, 508, 509, 514, 515, 540, 541, 563, 564, 565, 566, 567, 568, 569, 570, 571, 573, 574, 575, 576, 577, 578, 579, 580, 581, 582, 583, 584, 585, 588, 589, 591, 593, 594, 595, 596, 597, 598, 599, 602, 605, 615, 616, 617, 618, 619, 620, 621, 622, 623, 624, 625, 626, 627, 628, 629, 630, 631, 633, 634, 635, 636 |
| Uncertain | 1 | 632 |

## False positives

One subsection per finding. Each excerpt is the smallest source excerpt proving the decision; all excerpts are copied from the function-context files (identical to the chunk context blocks).

### [ ] Finding 6 — BP-PY-46

- Function context: `./scripts/findings/functions/6.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/examples/python/03-create-pdf/main.py:23:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("Saved: from_markdown.pdf")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 7 — BP-PY-46

- Function context: `./scripts/findings/functions/7.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/examples/python/03-create-pdf/main.py:34:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("Saved: from_html.pdf")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 8 — BP-PY-46

- Function context: `./scripts/findings/functions/8.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/examples/python/03-create-pdf/main.py:40:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("Saved: from_text.pdf")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 9 — BP-PY-46

- Function context: `./scripts/findings/functions/9.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/examples/python/03-create-pdf/main.py:42:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("Done. 3 PDFs created.")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 10 — PERF-PY-28

- Function context: `./scripts/findings/functions/10.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/examples/python/08-batch-processing/main.py:31:21`
- Checklist pattern: `not per unit of work — executor created once per batch/object lifetime`

Source excerpt:

```
    with ProcessPoolExecutor() as pool:
        futures = {pool.submit(process_pdf, p): p for p in paths}
```

Why this is a false positive: The executor is created once at the top of the batch function (or once per long-lived wrapper object in `__init__`); it is not created per unit of work.

Checklist evidence: not per unit of work — executor created once per batch/object lifetime — The executor is created once at the top of the batch function (or once per long-lived wrapper object in `__init__`); it is not created per unit of work.

### [ ] Finding 13 — BP-PY-46

- Function context: `./scripts/findings/functions/13.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/examples/python/09-new-features/image_embedding/main.py:196:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print(f"Written: {out_path}")
assert os.path.getsize(out_path) > 0, "output PDF must be non-empty"
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 14 — BP-PY-6

- Function context: `./scripts/findings/functions/14.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/examples/python/09-new-features/image_embedding/main.py:197:1`
- Checklist pattern: `not input/security validation — assert on own-output sanity check in an example`

Source excerpt:

```
assert os.path.getsize(out_path) > 0, "output PDF must be non-empty"
print("All image embedding checks passed.")
```

Why this is a false positive: The assert checks the size of a file the example script itself just wrote (own-output sanity check), not request/CLI/authz/path input validation; the needle hit is coincidental (`os.path` in `os.path.getsize`). The rule condition (assert used for input or security checks) is not met.

Checklist evidence: not input/security validation — assert on own-output sanity check in an example — The assert checks the size of a file the example script itself just wrote (own-output sanity check), not request/CLI/authz/path input validation; the needle hit is coincidental (`os.path` in `os.path.getsize`). The rule condition (assert used for input or security checks) is not met.

### [ ] Finding 15 — BP-PY-46

- Function context: `./scripts/findings/functions/15.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/examples/python/09-new-features/image_embedding/main.py:198:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print("All image embedding checks passed.")

```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 16 — BP-PY-46

- Function context: `./scripts/findings/functions/16.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/examples/python/09-new-features/ocr_scanned_pdf/main.py:76:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("OcrEngine and OcrConfig are available in this build.")
    except AttributeError as exc:
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 17 — BP-PY-46

- Function context: `./scripts/findings/functions/17.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/examples/python/09-new-features/ocr_scanned_pdf/main.py:78:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"OCR NOT available: {exc}")
        print("Rebuild with: maturin develop --features python,ocr")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 18 — BP-PY-46

- Function context: `./scripts/findings/functions/18.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/examples/python/09-new-features/ocr_scanned_pdf/main.py:79:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("Rebuild with: maturin develop --features python,ocr")
        sys.exit(1)
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 19 — CWE-22

- Function context: `./scripts/findings/functions/19.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/examples/python/09-new-features/office_conversion/main.py:37:10`
- Checklist pattern: `constant path — no external input reaches the open() sink`

Source excerpt:

```
    with open(os.path.join(OUT_DIR, "output.docx"), "wb") as f:
        f.write(docx_bytes)
```

Why this is a false positive: The path is `os.path.join(OUT_DIR, "output.docx")` where OUT_DIR is a module constant (`OUT_DIR = "output"`) and the second segment is a string literal; no external/user-influenced segment reaches the open() sink, so the traversal condition is not met.

Checklist evidence: constant path — no external input reaches the open() sink — The path is `os.path.join(OUT_DIR, "output.docx")` where OUT_DIR is a module constant (`OUT_DIR = "output"`) and the second segment is a string literal; no external/user-influenced segment reaches the open() sink, so the traversal condition is not met.

### [ ] Finding 28 — BP-PY-7

- Function context: `./scripts/findings/functions/28.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/python/pdf_oxide/_async.py:23:37`
- Checklist pattern: `wrong-construct match — library method/module function or a definition, not builtin open()`

Source excerpt:

```
        doc = await AsyncPdfDocument.open("report.pdf")
        text = await doc.extract_text(0)
```

Why this is a false positive: The matched `open` is not the builtin file-open function: it is a library/module method call (`fitz.open`, `pymupdf.open`, `AsyncPdfDocument.open`) or a function definition (`def open`). No builtin file handle is opened without a `with` statement, so the rule condition is not met.

Checklist evidence: wrong-construct match — library method/module function or a definition, not builtin open() — The matched `open` is not the builtin file-open function: it is a library/module method call (`fitz.open`, `pymupdf.open`, `AsyncPdfDocument.open`) or a function definition (`def open`). No builtin file handle is opened without a `with` statement, so the rule condition is not met.

### [ ] Finding 29 — PERF-PY-28

- Function context: `./scripts/findings/functions/29.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/python/pdf_oxide/_async.py:118:36`
- Checklist pattern: `not per unit of work — executor created once per batch/object lifetime`

Source excerpt:

```
        self._executor = ThreadPoolExecutor(max_workers=1)

```

Why this is a false positive: The executor is created once at the top of the batch function (or once per long-lived wrapper object in `__init__`); it is not created per unit of work.

Checklist evidence: not per unit of work — executor created once per batch/object lifetime — The executor is created once at the top of the batch function (or once per long-lived wrapper object in `__init__`); it is not created per unit of work.

### [ ] Finding 30 — BP-PY-7

- Function context: `./scripts/findings/functions/30.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/python/pdf_oxide/_async.py:123:15`
- Checklist pattern: `wrong-construct match — library method/module function or a definition, not builtin open()`

Source excerpt:

```
    async def open(path: str, password: str | None = None) -> AsyncPdfDocument:
        """Open a PDF file.  The document is created on the background thread."""
```

Why this is a false positive: The matched `open` is not the builtin file-open function: it is a library/module method call (`fitz.open`, `pymupdf.open`, `AsyncPdfDocument.open`) or a function definition (`def open`). No builtin file handle is opened without a `with` statement, so the rule condition is not met.

Checklist evidence: wrong-construct match — library method/module function or a definition, not builtin open() — The matched `open` is not the builtin file-open function: it is a library/module method call (`fitz.open`, `pymupdf.open`, `AsyncPdfDocument.open`) or a function definition (`def open`). No builtin file handle is opened without a `with` statement, so the rule condition is not met.

### [ ] Finding 44 — BP-PY-46

- Function context: `./scripts/findings/functions/44.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:18:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"\n{'=' * 70}")
    print(f"Layout Analysis: {pdf_path}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 45 — BP-PY-46

- Function context: `./scripts/findings/functions/45.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:19:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"Layout Analysis: {pdf_path}")
    print(f"{'=' * 70}\n")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 46 — BP-PY-46

- Function context: `./scripts/findings/functions/46.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:20:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"{'=' * 70}\n")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 47 — BP-PY-46

- Function context: `./scripts/findings/functions/47.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:29:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("No text found!")
        return
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 48 — BP-PY-46

- Function context: `./scripts/findings/functions/48.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:39:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"Total characters: {len(chars)}")
    print("\nPage Dimensions:")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 49 — BP-PY-46

- Function context: `./scripts/findings/functions/49.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:40:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("\nPage Dimensions:")
    print(f"  X range: {min(xs):.1f} to {max(xs):.1f} (width: {max(xs) - min(xs):.1f})")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 50 — BP-PY-46

- Function context: `./scripts/findings/functions/50.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:41:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  X range: {min(xs):.1f} to {max(xs):.1f} (width: {max(xs) - min(xs):.1f})")
    print(f"  Y range: {min(ys):.1f} to {max(ys):.1f} (height: {max(ys) - min(ys):.1f})")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 51 — BP-PY-46

- Function context: `./scripts/findings/functions/51.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:42:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  Y range: {min(ys):.1f} to {max(ys):.1f} (height: {max(ys) - min(ys):.1f})")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 52 — BP-PY-46

- Function context: `./scripts/findings/functions/52.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:44:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("\nFont Statistics:")
    print(
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 53 — BP-PY-46

- Function context: `./scripts/findings/functions/53.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:45:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(
        f"  Font sizes: min={min(font_sizes):.1f}, median={sorted(font_sizes)[len(font_sizes) // 2]:.1f}, max={max(font_sizes):.1f}"
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 54 — BP-PY-46

- Function context: `./scripts/findings/functions/54.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:48:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(
        f"  Char widths: min={min(widths):.1f}, median={sorted(widths)[len(widths) // 2]:.1f}, max={max(widths):.1f}"
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 55 — BP-PY-46

- Function context: `./scripts/findings/functions/55.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:51:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(
        f"  Char heights: min={min(heights):.1f}, median={sorted(heights)[len(heights) // 2]:.1f}, max={max(heights):.1f}"
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 56 — BP-PY-46

- Function context: `./scripts/findings/functions/56.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:56:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"\n{'=' * 70}")
    print("COLUMN DETECTION ANALYSIS")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 57 — BP-PY-46

- Function context: `./scripts/findings/functions/57.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:57:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("COLUMN DETECTION ANALYSIS")
    print(f"{'=' * 70}\n")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 58 — BP-PY-46

- Function context: `./scripts/findings/functions/58.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:58:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"{'=' * 70}\n")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 59 — BP-PY-46

- Function context: `./scripts/findings/functions/59.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:77:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"Average density per bin: {avg_density:.1f} chars")
    print(f"Threshold (10%): {threshold_10pct:.1f}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 60 — BP-PY-46

- Function context: `./scripts/findings/functions/60.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:78:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"Threshold (10%): {threshold_10pct:.1f}")
    print(f"Threshold (15%): {threshold_15pct:.1f}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 61 — BP-PY-46

- Function context: `./scripts/findings/functions/61.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:79:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"Threshold (15%): {threshold_15pct:.1f}")
    print(f"Threshold (20%): {threshold_20pct:.1f}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 62 — BP-PY-46

- Function context: `./scripts/findings/functions/62.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:80:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"Threshold (20%): {threshold_20pct:.1f}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 63 — BP-PY-46

- Function context: `./scripts/findings/functions/63.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:96:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("\nValleys found:")
    print(f"  With 10% threshold: {len(valleys_10)} valleys")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 64 — BP-PY-46

- Function context: `./scripts/findings/functions/64.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:97:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  With 10% threshold: {len(valleys_10)} valleys")
    print(f"  With 15% threshold: {len(valleys_15)} valleys")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 65 — BP-PY-46

- Function context: `./scripts/findings/functions/65.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:98:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  With 15% threshold: {len(valleys_15)} valleys")
    print(f"  With 20% threshold: {len(valleys_20)} valleys")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 66 — BP-PY-46

- Function context: `./scripts/findings/functions/66.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:99:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  With 20% threshold: {len(valleys_20)} valleys")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 67 — BP-PY-46

- Function context: `./scripts/findings/functions/67.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:103:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("\n  Likely column gaps (15% threshold):")
        # Group consecutive valleys
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 68 — BP-PY-46

- Function context: `./scripts/findings/functions/68.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:118:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print(
                    f"    Gap at X={gap_start:.1f}-{x_pos:.1f} (width: {gap_width:.1f}, avg density: {avg_gap_density:.1f})"
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 69 — BP-PY-46

- Function context: `./scripts/findings/functions/69.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:125:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("\nX Distribution (simplified histogram):")
    print(
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 70 — BP-PY-46

- Function context: `./scripts/findings/functions/70.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:126:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(
        "  "
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 71 — BP-PY-46

- Function context: `./scripts/findings/functions/71.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:151:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(line)

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 72 — BP-PY-46

- Function context: `./scripts/findings/functions/72.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:153:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("  " + "^" * 100)
    print(f"  0{' ' * 48}{page_width / 2:.0f}{' ' * 47}{page_width:.0f}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 73 — BP-PY-46

- Function context: `./scripts/findings/functions/73.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:154:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  0{' ' * 48}{page_width / 2:.0f}{' ' * 47}{page_width:.0f}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 74 — BP-PY-46

- Function context: `./scripts/findings/functions/74.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:157:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"\n{'=' * 70}")
    print("RECOMMENDATIONS FOR XY-CUT")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 75 — BP-PY-46

- Function context: `./scripts/findings/functions/75.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:158:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("RECOMMENDATIONS FOR XY-CUT")
    print(f"{'=' * 70}\n")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 76 — BP-PY-46

- Function context: `./scripts/findings/functions/76.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:159:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"{'=' * 70}\n")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 77 — BP-PY-46

- Function context: `./scripts/findings/functions/77.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:164:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("Document properties:")
    print(f"  Median font size: {median_font_size:.1f}pt")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 78 — BP-PY-46

- Function context: `./scripts/findings/functions/78.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:165:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  Median font size: {median_font_size:.1f}pt")
    print(f"  Median char height: {median_char_height:.1f}pt")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 79 — BP-PY-46

- Function context: `./scripts/findings/functions/79.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:166:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  Median char height: {median_char_height:.1f}pt")
    print(f"  Page width: {page_width:.1f}pt")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 80 — BP-PY-46

- Function context: `./scripts/findings/functions/80.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:167:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  Page width: {page_width:.1f}pt")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 81 — BP-PY-46

- Function context: `./scripts/findings/functions/81.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:169:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("\nRecommended XY-Cut parameters:")
    print(f"  min_region_size: {3 * median_char_height:.1f}  (3× char height)")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 82 — BP-PY-46

- Function context: `./scripts/findings/functions/82.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:170:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  min_region_size: {3 * median_char_height:.1f}  (3× char height)")
    print("  valley_threshold: Use 15-20% of average (not 10%)")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 83 — BP-PY-46

- Function context: `./scripts/findings/functions/83.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:171:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("  valley_threshold: Use 15-20% of average (not 10%)")
    print(f"    → Absolute threshold: {threshold_15pct:.1f} chars per bin")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 84 — BP-PY-46

- Function context: `./scripts/findings/functions/84.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:172:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"    → Absolute threshold: {threshold_15pct:.1f} chars per bin")
    print(f"  projection_bins: {int(page_width):.0f}  (1pt per bin, not width/2)")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 85 — BP-PY-46

- Function context: `./scripts/findings/functions/85.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:173:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  projection_bins: {int(page_width):.0f}  (1pt per bin, not width/2)")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 86 — BP-PY-46

- Function context: `./scripts/findings/functions/86.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:175:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("\nWhy current parameters fail:")
    print(f"  Current valley threshold (10%): {threshold_10pct:.1f}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 87 — BP-PY-46

- Function context: `./scripts/findings/functions/87.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:176:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  Current valley threshold (10%): {threshold_10pct:.1f}")
    print(f"  → Finds {len(valleys_10)} valleys (too few if < 1)")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 88 — BP-PY-46

- Function context: `./scripts/findings/functions/88.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:177:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  → Finds {len(valleys_10)} valleys (too few if < 1)")
    print(f"  Better threshold (15%): {threshold_15pct:.1f}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 89 — BP-PY-46

- Function context: `./scripts/findings/functions/89.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:178:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  Better threshold (15%): {threshold_15pct:.1f}")
    print(f"  → Finds {len(valleys_15)} valleys")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 90 — BP-PY-46

- Function context: `./scripts/findings/functions/90.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_layout.py:179:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  → Finds {len(valleys_15)} valleys")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 91 — BP-PY-7

- Function context: `./scripts/findings/functions/91.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:23:15`
- Checklist pattern: `wrong-construct match — library method/module function or a definition, not builtin open()`

Source excerpt:

```
    doc = fitz.open(pdf_path)
    page = doc[0]
```

Why this is a false positive: The matched `open` is not the builtin file-open function: it is a library/module method call (`fitz.open`, `pymupdf.open`, `AsyncPdfDocument.open`) or a function definition (`def open`). No builtin file handle is opened without a `with` statement, so the rule condition is not met.

Checklist evidence: wrong-construct match — library method/module function or a definition, not builtin open() — The matched `open` is not the builtin file-open function: it is a library/module method call (`fitz.open`, `pymupdf.open`, `AsyncPdfDocument.open`) or a function definition (`def open`). No builtin file handle is opened without a `with` statement, so the rule condition is not met.

### [ ] Finding 92 — BP-PY-46

- Function context: `./scripts/findings/functions/92.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:27:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("=" * 80)
    print("RAW CONTENT STREAM ANALYSIS")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 93 — BP-PY-46

- Function context: `./scripts/findings/functions/93.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:28:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("RAW CONTENT STREAM ANALYSIS")
    print("=" * 80)
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 94 — BP-PY-46

- Function context: `./scripts/findings/functions/94.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:29:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("=" * 80)

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 95 — BP-PY-46

- Function context: `./scripts/findings/functions/95.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:34:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"\nFound {len(blocks['blocks'])} blocks")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 96 — BP-PY-46

- Function context: `./scripts/findings/functions/96.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:38:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"\n{'=' * 80}")
            print(f"BLOCK {block_idx}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 97 — BP-PY-46

- Function context: `./scripts/findings/functions/97.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:39:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"BLOCK {block_idx}")
            print(f"{'=' * 80}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 98 — BP-PY-46

- Function context: `./scripts/findings/functions/98.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:40:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"{'=' * 80}")
            print(f"Bbox: {block['bbox']}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 99 — BP-PY-46

- Function context: `./scripts/findings/functions/99.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:41:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"Bbox: {block['bbox']}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 100 — BP-PY-46

- Function context: `./scripts/findings/functions/100.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:44:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print(f"\n  LINE {line_idx}:")
                print(f"  Bbox: {line['bbox']}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 101 — BP-PY-46

- Function context: `./scripts/findings/functions/101.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:45:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print(f"  Bbox: {line['bbox']}")
                print(f"  Direction: {line['dir']}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 102 — BP-PY-46

- Function context: `./scripts/findings/functions/102.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:46:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print(f"  Direction: {line['dir']}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 103 — BP-PY-46

- Function context: `./scripts/findings/functions/103.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:49:21`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                    print(f"\n    SPAN {span_idx}:")
                    print(f"    Text: '{span['text']}'")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 104 — BP-PY-46

- Function context: `./scripts/findings/functions/104.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:50:21`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                    print(f"    Text: '{span['text']}'")
                    print(f"    Bbox: {span['bbox']}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 105 — BP-PY-46

- Function context: `./scripts/findings/functions/105.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:51:21`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                    print(f"    Bbox: {span['bbox']}")
                    print(f"    Font: {span['font']}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 106 — BP-PY-46

- Function context: `./scripts/findings/functions/106.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:52:21`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                    print(f"    Font: {span['font']}")
                    print(f"    Size: {span['size']}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 107 — BP-PY-46

- Function context: `./scripts/findings/functions/107.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:53:21`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                    print(f"    Size: {span['size']}")
                    print(f"    Origin: {span['origin']}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 108 — BP-PY-46

- Function context: `./scripts/findings/functions/108.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:54:21`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                    print(f"    Origin: {span['origin']}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 109 — BP-PY-46

- Function context: `./scripts/findings/functions/109.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:63:25`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                        print(f"    Span width: {span_width:.2f}pt")
                        print(f"    Chars: {char_count}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 110 — BP-PY-46

- Function context: `./scripts/findings/functions/110.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:64:25`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                        print(f"    Chars: {char_count}")
                        print(f"    Avg char width: {avg_char_width:.2f}pt")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 111 — BP-PY-46

- Function context: `./scripts/findings/functions/111.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:65:25`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                        print(f"    Avg char width: {avg_char_width:.2f}pt")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 112 — BP-PY-46

- Function context: `./scripts/findings/functions/112.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:71:25`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                        print(f"    Gap to next: {gap:.2f}pt")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 114 — BP-PY-46

- Function context: `./scripts/findings/functions/114.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:75:29`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                            print(f"    Gap ratio: {gap_ratio:.2f}x char width")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 115 — BP-PY-46

- Function context: `./scripts/findings/functions/115.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:78:33`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                                print("    → Should MERGE (gap < 0.5x char width)")
                            else:
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 116 — BP-PY-46

- Function context: `./scripts/findings/functions/116.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:80:33`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                                print("    → Should INSERT SPACE (gap >= 0.5x char width)")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 117 — BP-PY-7

- Function context: `./scripts/findings/functions/117.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:87:15`
- Checklist pattern: `wrong-construct match — library method/module function or a definition, not builtin open()`

Source excerpt:

```
    doc = fitz.open(pdf_path)
    page = doc[0]
```

Why this is a false positive: The matched `open` is not the builtin file-open function: it is a library/module method call (`fitz.open`, `pymupdf.open`, `AsyncPdfDocument.open`) or a function definition (`def open`). No builtin file handle is opened without a `with` statement, so the rule condition is not met.

Checklist evidence: wrong-construct match — library method/module function or a definition, not builtin open() — The matched `open` is not the builtin file-open function: it is a library/module method call (`fitz.open`, `pymupdf.open`, `AsyncPdfDocument.open`) or a function definition (`def open`). No builtin file handle is opened without a `with` statement, so the rule condition is not met.

### [ ] Finding 118 — BP-PY-46

- Function context: `./scripts/findings/functions/118.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:90:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("\n" + "=" * 80)
    print("RAWDICT ANALYSIS (Low-level)")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 119 — BP-PY-46

- Function context: `./scripts/findings/functions/119.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:91:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("RAWDICT ANALYSIS (Low-level)")
    print("=" * 80)
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 120 — BP-PY-46

- Function context: `./scripts/findings/functions/120.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:92:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("=" * 80)

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 121 — BP-PY-46

- Function context: `./scripts/findings/functions/121.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:99:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"\nBlock bbox: {block['bbox']}")
            for line in block["lines"][:3]:  # First 3 lines
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 122 — BP-PY-46

- Function context: `./scripts/findings/functions/122.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:101:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print(f"  Line: {line['bbox']}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 123 — BP-PY-46

- Function context: `./scripts/findings/functions/123.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:104:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print(f"  {len(spans)} spans:")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 124 — BP-PY-46

- Function context: `./scripts/findings/functions/124.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:112:21`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                    print(f"    [{i}] '{text}' @ x={bbox[0]:.1f} w={width:.1f} chars={char_count}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 125 — BP-PY-46

- Function context: `./scripts/findings/functions/125.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:117:25`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                        print(f"         gap={gap:.1f}pt", end="")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 126 — BP-PY-46

- Function context: `./scripts/findings/functions/126.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:121:29`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                            print(f" ({gap / avg_cw:.2f}x)", end="")
                        print()
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 127 — BP-PY-46

- Function context: `./scripts/findings/functions/127.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:122:25`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                        print()

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 130 — BP-PY-46

- Function context: `./scripts/findings/functions/130.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/analyze_pdf_spacing.py:155:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"Error: {e}")
        import traceback
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 137 — BP-PY-7

- Function context: `./scripts/findings/functions/137.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/bench_pymupdf.py:47:18`
- Checklist pattern: `wrong-construct match — library method/module function or a definition, not builtin open()`

Source excerpt:

```
    doc = pymupdf.open(pdf_path)
    if doc.needs_pass:
```

Why this is a false positive: The matched `open` is not the builtin file-open function: it is a library/module method call (`fitz.open`, `pymupdf.open`, `AsyncPdfDocument.open`) or a function definition (`def open`). No builtin file handle is opened without a `with` statement, so the rule condition is not met.

Checklist evidence: wrong-construct match — library method/module function or a definition, not builtin open() — The matched `open` is not the builtin file-open function: it is a library/module method call (`fitz.open`, `pymupdf.open`, `AsyncPdfDocument.open`) or a function definition (`def open`). No builtin file handle is opened without a `with` statement, so the rule condition is not met.

### [ ] Finding 138 — PERF-PY-26

- Function context: `./scripts/findings/functions/138.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/bench_pymupdf.py:65:1`
- Checklist pattern: `not a hot path — one-time CLI parse at startup`

Source excerpt:

```
    args = parser.parse_args()

```

Why this is a false positive: `parser.parse_args()` is a one-time CLI configuration parse at program start, not an expensive decode/parse on a hot path.

Checklist evidence: not a hot path — one-time CLI parse at startup — `parser.parse_args()` is a one-time CLI configuration parse at program start, not an expensive decode/parse on a hot path.

### [ ] Finding 141 — CWE-829

- Function context: `./scripts/findings/functions/141.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/benchmark_all_libraries.py:50:13`
- Checklist pattern: `no untrusted source — module names are hardcoded literals in the source`

Source excerpt:

```
            __import__(import_name)
            AVAILABLE_LIBRARIES[name] = True
```

Why this is a false positive: The dynamically imported module names come from a hardcoded literal dict in the source; there is no untrusted control sphere selecting what executes.

Checklist evidence: no untrusted source — module names are hardcoded literals in the source — The dynamically imported module names come from a hardcoded literal dict in the source; there is no untrusted control sphere selecting what executes.

### [ ] Finding 142 — CWE-94

- Function context: `./scripts/findings/functions/142.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/benchmark_all_libraries.py:50:13`
- Checklist pattern: `no untrusted source — module names are hardcoded literals in the source`

Source excerpt:

```
            __import__(import_name)
            AVAILABLE_LIBRARIES[name] = True
```

Why this is a false positive: The value reaching the dynamic-import sink is a hardcoded module name from a source-literal dict, not externally influenced text.

Checklist evidence: no untrusted source — module names are hardcoded literals in the source — The value reaching the dynamic-import sink is a hardcoded module name from a source-literal dict, not externally influenced text.

### [ ] Finding 143 — BP-PY-7

- Function context: `./scripts/findings/functions/143.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/benchmark_all_libraries.py:87:15`
- Checklist pattern: `wrong-construct match — library method/module function or a definition, not builtin open()`

Source excerpt:

```
    doc = fitz.open(str(pdf_path))
    text_parts = []
```

Why this is a false positive: The matched `open` is not the builtin file-open function: it is a library/module method call (`fitz.open`, `pymupdf.open`, `AsyncPdfDocument.open`) or a function definition (`def open`). No builtin file handle is opened without a `with` statement, so the rule condition is not met.

Checklist evidence: wrong-construct match — library method/module function or a definition, not builtin open() — The matched `open` is not the builtin file-open function: it is a library/module method call (`fitz.open`, `pymupdf.open`, `AsyncPdfDocument.open`) or a function definition (`def open`). No builtin file handle is opened without a `with` statement, so the rule condition is not met.

### [ ] Finding 144 — CWE-1341

- Function context: `./scripts/findings/functions/144.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/benchmark_all_libraries.py:239:9`
- Checklist pattern: `distinct handles — two different objects are closed, no double release`

Source excerpt:

```
        textpage.close()
        page.close()
```

Why this is a false positive: The two adjacent `.close()` calls release two different handles (textpage/page, tp/page); no single resource handle is released twice.

Checklist evidence: distinct handles — two different objects are closed, no double release — The two adjacent `.close()` calls release two different handles (textpage/page, tp/page); no single resource handle is released twice.

### [ ] Finding 149 — BP-PY-46

- Function context: `./scripts/findings/functions/149.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/browser_download_pdfs.py:25:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("Playwright not found. Installing...")
    subprocess.check_call([sys.executable, "-m", "pip", "install", "playwright"])
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 150 — CWE-88

- Function context: `./scripts/findings/functions/150.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/browser_download_pdfs.py:26:5`
- Checklist pattern: `no externally influenced input — fixed literals / internal repo paths`

Source excerpt:

```
    subprocess.check_call([sys.executable, "-m", "pip", "install", "playwright"])
    subprocess.check_call([sys.executable, "-m", "playwright", "install", "chromium"])
```

Why this is a false positive: All argv elements are fixed literals except `sys.executable` (the interpreter path, always position 0 and never an option); no externally influenced input is embedded, so the rule condition is not met.

Checklist evidence: no externally influenced input — fixed literals / internal repo paths — All argv elements are fixed literals except `sys.executable` (the interpreter path, always position 0 and never an option); no externally influenced input is embedded, so the rule condition is not met.

### [ ] Finding 151 — BP-PY-46

- Function context: `./scripts/findings/functions/151.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/browser_download_pdfs.py:33:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("\n=== Downloading SEC Filings ===")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 152 — BP-PY-46

- Function context: `./scripts/findings/functions/152.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/browser_download_pdfs.py:54:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"\n{name} ({ticker})...")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 153 — BP-PY-46

- Function context: `./scripts/findings/functions/153.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/browser_download_pdfs.py:75:21`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                    print(f"  Downloading {filename}...")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 155 — BP-PY-46

- Function context: `./scripts/findings/functions/155.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/browser_download_pdfs.py:85:25`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                        print(f"    ✓ Downloaded {filename}")
                        time.sleep(1)
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 158 — BP-PY-46

- Function context: `./scripts/findings/functions/158.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/browser_download_pdfs.py:88:25`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                        print(f"    ✗ Failed: {e}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 160 — BP-PY-46

- Function context: `./scripts/findings/functions/160.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/browser_download_pdfs.py:91:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"  Error accessing {name}: {e}")
            continue
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 161 — BP-PY-46

- Function context: `./scripts/findings/functions/161.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/browser_download_pdfs.py:99:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("\n=== Downloading Government PDFs ===")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 162 — BP-PY-46

- Function context: `./scripts/findings/functions/162.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/browser_download_pdfs.py:125:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"\n{source['name']}...")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 163 — BP-PY-46

- Function context: `./scripts/findings/functions/163.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/browser_download_pdfs.py:156:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print(f"  Downloading {filename}...")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 164 — BP-PY-46

- Function context: `./scripts/findings/functions/164.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/browser_download_pdfs.py:165:21`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                    print("    ✓ Downloaded")
                    time.sleep(2)
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 166 — BP-PY-46

- Function context: `./scripts/findings/functions/166.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/browser_download_pdfs.py:168:21`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                    print(f"    ✗ Failed: {e}")
                    continue
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 168 — BP-PY-46

- Function context: `./scripts/findings/functions/168.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/browser_download_pdfs.py:172:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"  Error with {source['name']}: {e}")
            continue
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 169 — BP-PY-46

- Function context: `./scripts/findings/functions/169.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/browser_download_pdfs.py:180:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("\n=== Downloading Academic PDFs ===")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 170 — BP-PY-46

- Function context: `./scripts/findings/functions/170.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/browser_download_pdfs.py:201:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"\n{source['name']}...")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 171 — BP-PY-46

- Function context: `./scripts/findings/functions/171.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/browser_download_pdfs.py:226:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print(f"  Downloading {filename}...")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 172 — BP-PY-46

- Function context: `./scripts/findings/functions/172.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/browser_download_pdfs.py:235:21`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                    print("    ✓ Downloaded")
                    time.sleep(2)
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 174 — BP-PY-46

- Function context: `./scripts/findings/functions/174.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/browser_download_pdfs.py:238:21`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                    print(f"    ✗ Failed: {e}")
                    continue
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 176 — BP-PY-46

- Function context: `./scripts/findings/functions/176.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/browser_download_pdfs.py:242:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"  Error with {source['name']}: {e}")
            continue
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 177 — BP-PY-46

- Function context: `./scripts/findings/functions/177.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_pdf_structure.py:14:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"\n{'=' * 70}")
    print(f"Analyzing: {pdf_path}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 178 — BP-PY-46

- Function context: `./scripts/findings/functions/178.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_pdf_structure.py:15:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"Analyzing: {pdf_path}")
    print(f"{'=' * 70}\n")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 179 — BP-PY-46

- Function context: `./scripts/findings/functions/179.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_pdf_structure.py:16:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"{'=' * 70}\n")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 180 — BP-PY-46

- Function context: `./scripts/findings/functions/180.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_pdf_structure.py:24:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print(f"PDF Version: {reader.pdf_header}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 181 — BP-PY-46

- Function context: `./scripts/findings/functions/181.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_pdf_structure.py:32:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print(f"✅ MarkInfo present: {mark_info}")
                if "/Marked" in mark_info:
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 182 — BP-PY-46

- Function context: `./scripts/findings/functions/182.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_pdf_structure.py:36:25`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                        print(f"   → Marked: {marked} (Tagged PDF!)")
                    else:
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 183 — BP-PY-46

- Function context: `./scripts/findings/functions/183.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_pdf_structure.py:38:25`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                        print(f"   → Marked: {marked} (Not tagged)")
            else:
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 184 — BP-PY-46

- Function context: `./scripts/findings/functions/184.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_pdf_structure.py:40:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print("❌ No MarkInfo (likely not tagged)")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 185 — BP-PY-46

- Function context: `./scripts/findings/functions/185.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_pdf_structure.py:44:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print("✅ StructTreeRoot present (Tagged PDF!)")
                struct_root = catalog["/StructTreeRoot"]
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 186 — BP-PY-46

- Function context: `./scripts/findings/functions/186.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_pdf_structure.py:46:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print(f"   Structure Tree Root: {struct_root}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 187 — BP-PY-46

- Function context: `./scripts/findings/functions/187.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_pdf_structure.py:52:25`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                        print(f"   → Has {len(k)} top-level structure elements")
                    else:
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 188 — BP-PY-46

- Function context: `./scripts/findings/functions/188.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_pdf_structure.py:54:25`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                        print(f"   → Has structure elements: {k}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 189 — BP-PY-46

- Function context: `./scripts/findings/functions/189.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_pdf_structure.py:58:21`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                    print("   → Has ParentTree (maps marked content to structure)")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 190 — BP-PY-46

- Function context: `./scripts/findings/functions/190.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_pdf_structure.py:62:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print("❌ No StructTreeRoot (not a Tagged PDF)")
                return False
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 193 — BP-PY-46

- Function context: `./scripts/findings/functions/193.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_pdf_structure.py:66:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"❌ Error analyzing PDF: {e}")
        return False
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 195 — BP-PY-46

- Function context: `./scripts/findings/functions/195.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_span_spacing.py:8:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print("=" * 80)
print("TESTING PROBLEMATIC PDF")
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 196 — BP-PY-46

- Function context: `./scripts/findings/functions/196.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_span_spacing.py:9:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print("TESTING PROBLEMATIC PDF")
print("=" * 80)
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 197 — BP-PY-46

- Function context: `./scripts/findings/functions/197.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_span_spacing.py:10:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print("=" * 80)

```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 198 — BP-PY-46

- Function context: `./scripts/findings/functions/198.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_span_spacing.py:17:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print(f"\n Number of spans: {len(spans)}")

```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 199 — BP-PY-46

- Function context: `./scripts/findings/functions/199.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_span_spacing.py:21:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print(f"Total text length: {len(text)} characters")
print("\nFirst 300 chars:")
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 200 — BP-PY-46

- Function context: `./scripts/findings/functions/200.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_span_spacing.py:22:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print("\nFirst 300 chars:")
print(text[:300])
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 201 — BP-PY-46

- Function context: `./scripts/findings/functions/201.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_span_spacing.py:23:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print(text[:300])
print()
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 202 — BP-PY-46

- Function context: `./scripts/findings/functions/202.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_span_spacing.py:24:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print()

```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 203 — BP-PY-46

- Function context: `./scripts/findings/functions/203.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_span_spacing.py:28:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("❌ STILL BROKEN - text has spacing between characters")
    print("   Pattern found: ' F i s c a l '")
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 204 — BP-PY-46

- Function context: `./scripts/findings/functions/204.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_span_spacing.py:29:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("   Pattern found: ' F i s c a l '")
elif "Fiscal Year" in text[:300]:
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 205 — BP-PY-46

- Function context: `./scripts/findings/functions/205.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_span_spacing.py:31:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("✅ FIXED - text is properly formatted")
    print("   Pattern found: 'Fiscal Year'")
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 206 — BP-PY-46

- Function context: `./scripts/findings/functions/206.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_span_spacing.py:32:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("   Pattern found: 'Fiscal Year'")
else:
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 207 — BP-PY-46

- Function context: `./scripts/findings/functions/207.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_span_spacing.py:34:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("⚠️  UNEXPECTED - neither pattern found")
    print("   Text sample:", text[:100])
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 208 — BP-PY-46

- Function context: `./scripts/findings/functions/208.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_span_spacing.py:35:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("   Text sample:", text[:100])

```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 209 — BP-PY-46

- Function context: `./scripts/findings/functions/209.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/check_span_spacing.py:37:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print("\n" + "=" * 80)

```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 213 — BP-PY-46

- Function context: `./scripts/findings/functions/213.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:112:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"Found {len(pairs)} matching file pairs")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 214 — BP-PY-46

- Function context: `./scripts/findings/functions/214.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:116:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print(f"Processing {i}/{len(pairs)}...")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 215 — BP-PY-46

- Function context: `./scripts/findings/functions/215.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:125:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"Completed analyzing {len(self.results)} files")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 217 — BP-PY-46

- Function context: `./scripts/findings/functions/217.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:216:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"Saved results to {output_file}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 218 — BP-PY-46

- Function context: `./scripts/findings/functions/218.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:222:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("\n" + "=" * 70)
        print("COMPREHENSIVE COMPARISON SUMMARY")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 219 — BP-PY-46

- Function context: `./scripts/findings/functions/219.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:223:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("COMPREHENSIVE COMPARISON SUMMARY")
        print("=" * 70)
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 220 — BP-PY-46

- Function context: `./scripts/findings/functions/220.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:224:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("=" * 70)
        print(f"\nTotal files analyzed: {stats['total_files']}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 221 — BP-PY-46

- Function context: `./scripts/findings/functions/221.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:225:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"\nTotal files analyzed: {stats['total_files']}")
        print("\nSize Analysis:")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 222 — BP-PY-46

- Function context: `./scripts/findings/functions/222.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:226:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("\nSize Analysis:")
        print(f"  Average size ratio (ours/pymupdf): {stats['avg_size_ratio']:.3f}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 223 — BP-PY-46

- Function context: `./scripts/findings/functions/223.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:227:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  Average size ratio (ours/pymupdf): {stats['avg_size_ratio']:.3f}")
        print(f"  Median size ratio: {stats['median_size_ratio']:.3f}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 224 — BP-PY-46

- Function context: `./scripts/findings/functions/224.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:228:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  Median size ratio: {stats['median_size_ratio']:.3f}")
        print(f"  Files with excellent size match (95-105%): {stats['size_match_excellent']}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 225 — BP-PY-46

- Function context: `./scripts/findings/functions/225.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:229:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  Files with excellent size match (95-105%): {stats['size_match_excellent']}")
        print(f"  Files with good size match (90-110%): {stats['size_match_good']}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 226 — BP-PY-46

- Function context: `./scripts/findings/functions/226.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:230:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  Files with good size match (90-110%): {stats['size_match_good']}")
        print(f"  Files significantly smaller: {stats['size_smaller']}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 227 — BP-PY-46

- Function context: `./scripts/findings/functions/227.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:231:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  Files significantly smaller: {stats['size_smaller']}")
        print(f"  Files significantly larger: {stats['size_larger']}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 228 — BP-PY-46

- Function context: `./scripts/findings/functions/228.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:232:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  Files significantly larger: {stats['size_larger']}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 229 — BP-PY-46

- Function context: `./scripts/findings/functions/229.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:234:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("\nTotal Output Sizes:")
        print(
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 230 — BP-PY-46

- Function context: `./scripts/findings/functions/230.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:235:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(
            f"  Our library: {stats['our_total_size']:,} bytes ({stats['our_total_size'] / 1024 / 1024:.2f} MB)"
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 231 — BP-PY-46

- Function context: `./scripts/findings/functions/231.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:238:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(
            f"  PyMuPDF4LLM: {stats['pymupdf_total_size']:,} bytes ({stats['pymupdf_total_size'] / 1024 / 1024:.2f} MB)"
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 232 — BP-PY-46

- Function context: `./scripts/findings/functions/232.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:241:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  Difference: {stats['our_total_size'] - stats['pymupdf_total_size']:,} bytes")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 233 — BP-PY-46

- Function context: `./scripts/findings/functions/233.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:243:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("\nFeature Comparison:")
        print("  <br> tags:")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 234 — BP-PY-46

- Function context: `./scripts/findings/functions/234.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:244:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("  <br> tags:")
        print(f"    Our library: {stats['total_br_tags_ours']:,}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 235 — BP-PY-46

- Function context: `./scripts/findings/functions/235.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:245:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"    Our library: {stats['total_br_tags_ours']:,}")
        print(f"    PyMuPDF4LLM: {stats['total_br_tags_pymupdf']:,}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 236 — BP-PY-46

- Function context: `./scripts/findings/functions/236.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:246:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"    PyMuPDF4LLM: {stats['total_br_tags_pymupdf']:,}")
        print("  Bold markers:")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 237 — BP-PY-46

- Function context: `./scripts/findings/functions/237.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:247:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("  Bold markers:")
        print(f"    Our library: {stats['total_bold_ours']:,}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 238 — BP-PY-46

- Function context: `./scripts/findings/functions/238.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:248:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"    Our library: {stats['total_bold_ours']:,}")
        print(f"    PyMuPDF4LLM: {stats['total_bold_pymupdf']:,}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 239 — BP-PY-46

- Function context: `./scripts/findings/functions/239.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:249:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"    PyMuPDF4LLM: {stats['total_bold_pymupdf']:,}")
        print(f"    Ratio: 1:{stats['total_bold_pymupdf'] / max(stats['total_bold_ours'], 1):.1f}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 240 — BP-PY-46

- Function context: `./scripts/findings/functions/240.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:250:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"    Ratio: 1:{stats['total_bold_pymupdf'] / max(stats['total_bold_ours'], 1):.1f}")
        print("  Form fields:")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 241 — BP-PY-46

- Function context: `./scripts/findings/functions/241.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:251:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("  Form fields:")
        print(f"    Our library: {stats['files_with_forms_ours']} files")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 242 — BP-PY-46

- Function context: `./scripts/findings/functions/242.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:252:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"    Our library: {stats['files_with_forms_ours']} files")
        print(f"    PyMuPDF4LLM: {stats['files_with_forms_pymupdf']} files")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 243 — BP-PY-46

- Function context: `./scripts/findings/functions/243.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:253:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"    PyMuPDF4LLM: {stats['files_with_forms_pymupdf']} files")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 244 — BP-PY-46

- Function context: `./scripts/findings/functions/244.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:255:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("\nText Quality:")
        print("  Files with potential garbled text:")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 245 — BP-PY-46

- Function context: `./scripts/findings/functions/245.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:256:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("  Files with potential garbled text:")
        print(f"    Our library: {stats['files_with_garbled_text_ours']}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 246 — BP-PY-46

- Function context: `./scripts/findings/functions/246.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:257:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"    Our library: {stats['files_with_garbled_text_ours']}")
        print(f"    PyMuPDF4LLM: {stats['files_with_garbled_text_pymupdf']}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 247 — BP-PY-46

- Function context: `./scripts/findings/functions/247.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:258:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"    PyMuPDF4LLM: {stats['files_with_garbled_text_pymupdf']}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 248 — BP-PY-46

- Function context: `./scripts/findings/functions/248.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:260:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("\nBy Category:")
        for cat, cat_stats in sorted(stats["by_category"].items()):
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 249 — BP-PY-46

- Function context: `./scripts/findings/functions/249.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:262:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"  {cat}:")
            print(f"    Files: {cat_stats['count']}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 250 — BP-PY-46

- Function context: `./scripts/findings/functions/250.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:263:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"    Files: {cat_stats['count']}")
            print(f"    Avg size ratio: {cat_stats['avg_size_ratio']:.3f}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 251 — BP-PY-46

- Function context: `./scripts/findings/functions/251.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:264:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"    Avg size ratio: {cat_stats['avg_size_ratio']:.3f}")
            print(f"    With forms: {cat_stats['with_forms']}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 252 — BP-PY-46

- Function context: `./scripts/findings/functions/252.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:265:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"    With forms: {cat_stats['with_forms']}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 253 — BP-PY-46

- Function context: `./scripts/findings/functions/253.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_all_outputs.py:267:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("\n" + "=" * 70)

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 258 — BP-PY-7

- Function context: `./scripts/findings/functions/258.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/compare_extractors.py:45:19`
- Checklist pattern: `wrong-construct match — library method/module function or a definition, not builtin open()`

Source excerpt:

```
        doc = fitz.open(PDF_PATH)
        page = doc[0]
```

Why this is a false positive: The matched `open` is not the builtin file-open function: it is a library/module method call (`fitz.open`, `pymupdf.open`, `AsyncPdfDocument.open`) or a function definition (`def open`). No builtin file handle is opened without a `with` statement, so the rule condition is not met.

Checklist evidence: wrong-construct match — library method/module function or a definition, not builtin open() — The matched `open` is not the builtin file-open function: it is a library/module method call (`fitz.open`, `pymupdf.open`, `AsyncPdfDocument.open`) or a function definition (`def open`). No builtin file handle is opened without a `with` statement, so the rule condition is not met.

### [ ] Finding 268 — BP-PY-46

- Function context: `./scripts/findings/functions/268.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:26:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("Converting LayoutReader model...")
    print("=" * 60)
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 269 — BP-PY-46

- Function context: `./scripts/findings/functions/269.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:27:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("=" * 60)

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 270 — BP-PY-46

- Function context: `./scripts/findings/functions/270.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:31:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"Loading model: {model_name}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 273 — BP-PY-46

- Function context: `./scripts/findings/functions/273.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:37:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"⚠ Failed to download model: {e}")
        print("  Note: This requires internet connection and ~400MB download")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 274 — BP-PY-46

- Function context: `./scripts/findings/functions/274.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:38:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("  Note: This requires internet connection and ~400MB download")
        print("  For testing purposes, you can skip this and use simplified heuristics")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 275 — BP-PY-46

- Function context: `./scripts/findings/functions/275.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:39:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("  For testing purposes, you can skip this and use simplified heuristics")
        return False
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 276 — BP-PY-12

- Function context: `./scripts/findings/functions/276.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:43:11`
- Checklist pattern: `wrong-construct match — zero-argument method call named eval, not builtin eval on dynamic input`

Source excerpt:

```
    model.eval()

```

Why this is a false positive: `model.eval()` is a zero-argument method call on the PyTorch model object, not the builtin `eval`; no dynamic input reaches eval/exec, so the rule condition (eval/exec on dynamic input) is not met.

Checklist evidence: wrong-construct match — zero-argument method call named eval, not builtin eval on dynamic input — `model.eval()` is a zero-argument method call on the PyTorch model object, not the builtin `eval`; no dynamic input reaches eval/exec, so the rule condition (eval/exec on dynamic input) is not met.

### [ ] Finding 277 — BP-PY-46

- Function context: `./scripts/findings/functions/277.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:56:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("Exporting to ONNX...")
    torch.onnx.export(
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 278 — BP-PY-46

- Function context: `./scripts/findings/functions/278.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:72:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"✓ Model exported to {output_path}")
    print(f"  Size: {output_path.stat().st_size / 1024 / 1024:.1f} MB")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 279 — BP-PY-46

- Function context: `./scripts/findings/functions/279.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:73:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  Size: {output_path.stat().st_size / 1024 / 1024:.1f} MB")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 280 — BP-PY-46

- Function context: `./scripts/findings/functions/280.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:76:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("Quantizing to INT8...")
    quantized_path = Path("models/layout_reader_int8.onnx")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 281 — BP-PY-46

- Function context: `./scripts/findings/functions/281.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:80:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"✓ Quantized model saved to {quantized_path}")
    print(f"  Size: {quantized_path.stat().st_size / 1024 / 1024:.1f} MB")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 282 — BP-PY-46

- Function context: `./scripts/findings/functions/282.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:81:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  Size: {quantized_path.stat().st_size / 1024 / 1024:.1f} MB")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 283 — BP-PY-46

- Function context: `./scripts/findings/functions/283.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:85:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("✓ Tokenizer saved to models/tokenizer")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 284 — BP-PY-46

- Function context: `./scripts/findings/functions/284.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:92:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("\nConverting heading classifier...")
    print("=" * 60)
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 285 — BP-PY-46

- Function context: `./scripts/findings/functions/285.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:93:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("=" * 60)

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 286 — BP-PY-46

- Function context: `./scripts/findings/functions/286.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:98:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"Loading model: {model_name}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 288 — BP-PY-46

- Function context: `./scripts/findings/functions/288.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:107:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"⚠ Failed to download model: {e}")
        print("  Note: This requires internet connection and ~250MB download")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 289 — BP-PY-46

- Function context: `./scripts/findings/functions/289.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:108:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("  Note: This requires internet connection and ~250MB download")
        print("  For testing purposes, you can skip this and use rule-based classification")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 290 — BP-PY-46

- Function context: `./scripts/findings/functions/290.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:109:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("  For testing purposes, you can skip this and use rule-based classification")
        return False
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 291 — BP-PY-12

- Function context: `./scripts/findings/functions/291.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:112:11`
- Checklist pattern: `wrong-construct match — zero-argument method call named eval, not builtin eval on dynamic input`

Source excerpt:

```
    model.eval()

```

Why this is a false positive: `model.eval()` is a zero-argument method call on the PyTorch model object, not the builtin `eval`; no dynamic input reaches eval/exec, so the rule condition (eval/exec on dynamic input) is not met.

Checklist evidence: wrong-construct match — zero-argument method call named eval, not builtin eval on dynamic input — `model.eval()` is a zero-argument method call on the PyTorch model object, not the builtin `eval`; no dynamic input reaches eval/exec, so the rule condition (eval/exec on dynamic input) is not met.

### [ ] Finding 292 — BP-PY-46

- Function context: `./scripts/findings/functions/292.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:119:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("Exporting to ONNX...")
    torch.onnx.export(
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 293 — BP-PY-46

- Function context: `./scripts/findings/functions/293.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:133:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"✓ Model exported to {output_path}")
    print(f"  Size: {output_path.stat().st_size / 1024 / 1024:.1f} MB")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 294 — BP-PY-46

- Function context: `./scripts/findings/functions/294.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:134:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  Size: {output_path.stat().st_size / 1024 / 1024:.1f} MB")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 295 — BP-PY-46

- Function context: `./scripts/findings/functions/295.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:137:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("Quantizing to INT8...")
    quantized_path = Path("models/heading_classifier_int8.onnx")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 296 — BP-PY-46

- Function context: `./scripts/findings/functions/296.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:141:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"✓ Quantized model saved to {quantized_path}")
    print(f"  Size: {quantized_path.stat().st_size / 1024 / 1024:.1f} MB")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 297 — BP-PY-46

- Function context: `./scripts/findings/functions/297.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:142:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  Size: {quantized_path.stat().st_size / 1024 / 1024:.1f} MB")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 298 — BP-PY-46

- Function context: `./scripts/findings/functions/298.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:145:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("✓ Tokenizer saved to models/heading_tokenizer")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 299 — BP-PY-46

- Function context: `./scripts/findings/functions/299.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:152:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("\nVerifying models...")
    print("=" * 60)
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 300 — BP-PY-46

- Function context: `./scripts/findings/functions/300.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:153:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("=" * 60)

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 301 — BP-PY-46

- Function context: `./scripts/findings/functions/301.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:164:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print(f"✓ {model_path} verified")
                print(f"  Inputs: {[i.name for i in session.get_inputs()]}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 302 — BP-PY-46

- Function context: `./scripts/findings/functions/302.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:165:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print(f"  Inputs: {[i.name for i in session.get_inputs()]}")
                print(f"  Outputs: {[o.name for o in session.get_outputs()]}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 303 — BP-PY-46

- Function context: `./scripts/findings/functions/303.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:166:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print(f"  Outputs: {[o.name for o in session.get_outputs()]}")
            except Exception as e:
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 305 — BP-PY-46

- Function context: `./scripts/findings/functions/305.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:168:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print(f"✗ {model_path} failed verification: {e}")
                success = False
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 306 — BP-PY-46

- Function context: `./scripts/findings/functions/306.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:171:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"⚠ {model_path} not found (skipped)")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 307 — BP-PY-46

- Function context: `./scripts/findings/functions/307.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:204:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("\n" + "=" * 60)
    print("PDF Library - Model Conversion Script")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 308 — BP-PY-46

- Function context: `./scripts/findings/functions/308.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:205:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("PDF Library - Model Conversion Script")
    print("=" * 60)
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 309 — BP-PY-46

- Function context: `./scripts/findings/functions/309.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:206:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("=" * 60)

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 310 — BP-PY-46

- Function context: `./scripts/findings/functions/310.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:219:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("\n" + "=" * 60)
    if success:
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 311 — BP-PY-46

- Function context: `./scripts/findings/functions/311.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:221:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("✓ Model conversion complete!")
        print("\nNext steps:")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 312 — BP-PY-46

- Function context: `./scripts/findings/functions/312.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:222:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("\nNext steps:")
        print("  1. Run: cargo build --features ml")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 313 — BP-PY-46

- Function context: `./scripts/findings/functions/313.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:223:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("  1. Run: cargo build --features ml")
        print("  2. Test: cargo test --features ml")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 314 — BP-PY-46

- Function context: `./scripts/findings/functions/314.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:224:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("  2. Test: cargo test --features ml")
        print("\nThe ML models are ready for use.")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 315 — BP-PY-46

- Function context: `./scripts/findings/functions/315.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:225:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("\nThe ML models are ready for use.")
    else:
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 316 — BP-PY-46

- Function context: `./scripts/findings/functions/316.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:227:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("⚠ Model conversion completed with warnings")
        print("\nThe library will fall back to rule-based algorithms.")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 317 — BP-PY-46

- Function context: `./scripts/findings/functions/317.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:228:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("\nThe library will fall back to rule-based algorithms.")
        print("For full ML functionality, ensure:")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 318 — BP-PY-46

- Function context: `./scripts/findings/functions/318.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:229:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("For full ML functionality, ensure:")
        print("  - Internet connection is available")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 319 — BP-PY-46

- Function context: `./scripts/findings/functions/319.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:230:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("  - Internet connection is available")
        print("  - pip install -r scripts/requirements.txt is run")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 320 — BP-PY-46

- Function context: `./scripts/findings/functions/320.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:231:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("  - pip install -r scripts/requirements.txt is run")
        print("  - Sufficient disk space (~1GB for downloads)")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 321 — BP-PY-46

- Function context: `./scripts/findings/functions/321.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:232:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("  - Sufficient disk space (~1GB for downloads)")
    print("=" * 60)
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 322 — BP-PY-46

- Function context: `./scripts/findings/functions/322.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_models.py:233:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("=" * 60)

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 323 — BP-PY-46

- Function context: `./scripts/findings/functions/323.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_pdf_spec.py:7:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print("Converting PDF specification to markdown...")
print("This may take a few minutes for a 750+ page document...")
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 324 — BP-PY-46

- Function context: `./scripts/findings/functions/324.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_pdf_spec.py:8:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print("This may take a few minutes for a 750+ page document...")

```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 325 — BP-PY-46

- Function context: `./scripts/findings/functions/325.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_pdf_spec.py:19:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("✅ Conversion complete!")
    print("   Output: docs/spec/pdf.md")
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 326 — BP-PY-46

- Function context: `./scripts/findings/functions/326.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_pdf_spec.py:20:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("   Output: docs/spec/pdf.md")
    print(f"   Size: {len(md_text):,} characters")
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 327 — BP-PY-46

- Function context: `./scripts/findings/functions/327.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_pdf_spec.py:21:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"   Size: {len(md_text):,} characters")

```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 330 — BP-PY-46

- Function context: `./scripts/findings/functions/330.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/convert_pdf_spec.py:24:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"❌ Error: {e}")
    import traceback
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 331 — BP-PY-7

- Function context: `./scripts/findings/functions/331.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/debug_font_encoding.py:11:11`
- Checklist pattern: `wrong-construct match — library method/module function or a definition, not builtin open()`

Source excerpt:

```
doc = fitz.open(PDF_PATH)
page = doc[0]
```

Why this is a false positive: The matched `open` is not the builtin file-open function: it is a library/module method call (`fitz.open`, `pymupdf.open`, `AsyncPdfDocument.open`) or a function definition (`def open`). No builtin file handle is opened without a `with` statement, so the rule condition is not met.

Checklist evidence: wrong-construct match — library method/module function or a definition, not builtin open() — The matched `open` is not the builtin file-open function: it is a library/module method call (`fitz.open`, `pymupdf.open`, `AsyncPdfDocument.open`) or a function definition (`def open`). No builtin file handle is opened without a `with` statement, so the rule condition is not met.

### [ ] Finding 332 — BP-PY-46

- Function context: `./scripts/findings/functions/332.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/debug_font_encoding.py:14:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print("=" * 80)
print("FONT ENCODING ANALYSIS")
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 333 — BP-PY-46

- Function context: `./scripts/findings/functions/333.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/debug_font_encoding.py:15:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print("FONT ENCODING ANALYSIS")
print("=" * 80)
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 334 — BP-PY-46

- Function context: `./scripts/findings/functions/334.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/debug_font_encoding.py:16:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print("=" * 80)

```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 335 — BP-PY-46

- Function context: `./scripts/findings/functions/335.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/debug_font_encoding.py:20:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print(f"\nFonts on page 1: {len(fonts)}")
for i, font in enumerate(fonts):
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 336 — BP-PY-46

- Function context: `./scripts/findings/functions/336.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/debug_font_encoding.py:22:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"\n  Font #{i + 1}:")
    print(f"    Name: {font[3]}")  # Font name
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 337 — BP-PY-46

- Function context: `./scripts/findings/functions/337.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/debug_font_encoding.py:23:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"    Name: {font[3]}")  # Font name
    print(f"    Type: {font[1]}")  # Font type
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 338 — BP-PY-46

- Function context: `./scripts/findings/functions/338.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/debug_font_encoding.py:24:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"    Type: {font[1]}")  # Font type
    print(f"    Encoding: {font[2]}")  # Encoding
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 339 — BP-PY-46

- Function context: `./scripts/findings/functions/339.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/debug_font_encoding.py:25:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"    Encoding: {font[2]}")  # Encoding
    print(f"    Reference: {font[0]}")  # Font reference
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 340 — BP-PY-46

- Function context: `./scripts/findings/functions/340.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/debug_font_encoding.py:26:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"    Reference: {font[0]}")  # Font reference

```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 341 — BP-PY-46

- Function context: `./scripts/findings/functions/341.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/debug_font_encoding.py:34:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print("\n" + "=" * 80)
print("SEARCHING FOR 'Fiscal' IN CONTENT STREAM")
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 342 — BP-PY-46

- Function context: `./scripts/findings/functions/342.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/debug_font_encoding.py:35:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print("SEARCHING FOR 'Fiscal' IN CONTENT STREAM")
print("=" * 80)
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 343 — BP-PY-46

- Function context: `./scripts/findings/functions/343.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/debug_font_encoding.py:36:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print("=" * 80)

```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 344 — BP-PY-46

- Function context: `./scripts/findings/functions/344.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/debug_font_encoding.py:45:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print(f"\nFound {len(tj_arrays)} TJ arrays")

```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 345 — BP-PY-46

- Function context: `./scripts/findings/functions/345.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/debug_font_encoding.py:50:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"\n--- TJ Array #{i + 1} ---")
    print("Raw content (first 200 chars):")
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 346 — BP-PY-46

- Function context: `./scripts/findings/functions/346.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/debug_font_encoding.py:51:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("Raw content (first 200 chars):")
    print(array_content[:200])
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 347 — BP-PY-46

- Function context: `./scripts/findings/functions/347.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/debug_font_encoding.py:52:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(array_content[:200])

```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 348 — BP-PY-46

- Function context: `./scripts/findings/functions/348.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/debug_font_encoding.py:58:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"\nHex strings found: {len(hex_strings)}")
        for j, hex_str in enumerate(hex_strings[:3]):
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 349 — BP-PY-46

- Function context: `./scripts/findings/functions/349.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/debug_font_encoding.py:60:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"  Hex #{j + 1}: {hex_str[:40]}...")
            # Convert hex to bytes
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 350 — BP-PY-46

- Function context: `./scripts/findings/functions/350.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/debug_font_encoding.py:64:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print(f"    Bytes: {list(byte_array[:20])}")
                # Try UTF-16 BE decode
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 351 — BP-PY-46

- Function context: `./scripts/findings/functions/351.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/debug_font_encoding.py:68:21`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                    print(f"    UTF-16 BE: '{text[:50]}'")
                except Exception:
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 357 — BP-PY-46

- Function context: `./scripts/findings/functions/357.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/debug_font_encoding.py:74:21`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                    print(f"    Latin-1: '{text[:50]}'")
                except Exception:
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 362 — BP-PY-46

- Function context: `./scripts/findings/functions/362.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/debug_font_encoding.py:84:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"\nParenthesized strings: {len(paren_strings)}")
        for j, s in enumerate(paren_strings[:3]):
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 363 — BP-PY-46

- Function context: `./scripts/findings/functions/363.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/debug_font_encoding.py:86:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"  String #{j + 1}: '{s[:50]}'")

```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 365 — BP-PY-46

- Function context: `./scripts/findings/functions/365.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:25:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("=== Diagnosing Character Order ===")
    print(f"PDF: {pdf_path}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 366 — BP-PY-46

- Function context: `./scripts/findings/functions/366.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:26:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"PDF: {pdf_path}")
    print(f"Page: {page_num}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 367 — BP-PY-46

- Function context: `./scripts/findings/functions/367.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:27:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"Page: {page_num}")
    print()
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 368 — BP-PY-46

- Function context: `./scripts/findings/functions/368.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:28:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print()

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 369 — BP-PY-46

- Function context: `./scripts/findings/functions/369.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:34:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"Total characters extracted: {len(chars)}")
    print()
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 370 — BP-PY-46

- Function context: `./scripts/findings/functions/370.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:35:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print()

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 371 — BP-PY-46

- Function context: `./scripts/findings/functions/371.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:38:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("=" * 80)
    print("EXTRACTION ORDER (as they appear in PDF content stream):")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 372 — BP-PY-46

- Function context: `./scripts/findings/functions/372.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:39:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("EXTRACTION ORDER (as they appear in PDF content stream):")
    print("=" * 80)
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 373 — BP-PY-46

- Function context: `./scripts/findings/functions/373.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:40:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("=" * 80)
    for i, char in enumerate(chars[:max_chars]):
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 374 — BP-PY-46

- Function context: `./scripts/findings/functions/374.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:49:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"{i:3d}: '{display_char}' at X={x:6.1f} Y={y:6.1f} size={font_size:.1f}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 375 — BP-PY-46

- Function context: `./scripts/findings/functions/375.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:51:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print()

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 376 — BP-PY-46

- Function context: `./scripts/findings/functions/376.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:57:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("=" * 80)
    print("SPATIAL ORDER (top-to-bottom, left-to-right):")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 377 — BP-PY-46

- Function context: `./scripts/findings/functions/377.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:58:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("SPATIAL ORDER (top-to-bottom, left-to-right):")
    print("=" * 80)
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 378 — BP-PY-46

- Function context: `./scripts/findings/functions/378.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:59:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("=" * 80)
    for i, char in enumerate(sorted_chars):
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 379 — BP-PY-46

- Function context: `./scripts/findings/functions/379.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:66:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"{i:3d}: '{display_char}' at X={x:6.1f} Y={y:6.1f} size={font_size:.1f}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 380 — BP-PY-46

- Function context: `./scripts/findings/functions/380.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:68:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print()

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 381 — BP-PY-46

- Function context: `./scripts/findings/functions/381.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:74:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("=" * 80)
    print("TEXT COMPARISON:")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 382 — BP-PY-46

- Function context: `./scripts/findings/functions/382.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:75:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("TEXT COMPARISON:")
    print("=" * 80)
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 383 — BP-PY-46

- Function context: `./scripts/findings/functions/383.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:76:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("=" * 80)
    print(f"\nExtraction order text (first {max_chars} chars):")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 384 — BP-PY-46

- Function context: `./scripts/findings/functions/384.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:77:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"\nExtraction order text (first {max_chars} chars):")
    print(f"{extraction_order_text[:200]}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 385 — BP-PY-46

- Function context: `./scripts/findings/functions/385.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:78:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"{extraction_order_text[:200]}")
    print()
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 386 — BP-PY-46

- Function context: `./scripts/findings/functions/386.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:79:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print()
    print(f"Spatial order text (first {max_chars} chars):")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 387 — BP-PY-46

- Function context: `./scripts/findings/functions/387.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:80:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"Spatial order text (first {max_chars} chars):")
    print(f"{spatial_order_text[:200]}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 388 — BP-PY-46

- Function context: `./scripts/findings/functions/388.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:81:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"{spatial_order_text[:200]}")
    print()
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 389 — BP-PY-46

- Function context: `./scripts/findings/functions/389.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:82:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print()

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 390 — BP-PY-46

- Function context: `./scripts/findings/functions/390.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:86:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("❌ ORDER MISMATCH DETECTED!")
        print("Characters in content stream are NOT in spatial left-to-right order.")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 391 — BP-PY-46

- Function context: `./scripts/findings/functions/391.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:87:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("Characters in content stream are NOT in spatial left-to-right order.")
        print("This is the ROOT CAUSE of column mixing.")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 392 — BP-PY-46

- Function context: `./scripts/findings/functions/392.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:88:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("This is the ROOT CAUSE of column mixing.")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 393 — BP-PY-46

- Function context: `./scripts/findings/functions/393.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:93:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print(f"\nFirst difference at position {i}:")
                print(f"  Extraction order: '{c1}'")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 394 — BP-PY-46

- Function context: `./scripts/findings/functions/394.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:94:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print(f"  Extraction order: '{c1}'")
                print(f"  Spatial order:    '{c2}'")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 395 — BP-PY-46

- Function context: `./scripts/findings/functions/395.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:95:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print(f"  Spatial order:    '{c2}'")
                break
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 396 — BP-PY-46

- Function context: `./scripts/findings/functions/396.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:98:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("✅ Orders match - content stream is already in spatial order")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 397 — BP-PY-46

- Function context: `./scripts/findings/functions/397.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:100:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print()

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 398 — BP-PY-46

- Function context: `./scripts/findings/functions/398.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:103:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("=" * 80)
    print("Y-COORDINATE ANALYSIS:")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 399 — BP-PY-46

- Function context: `./scripts/findings/functions/399.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:104:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("Y-COORDINATE ANALYSIS:")
    print("=" * 80)
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 400 — BP-PY-46

- Function context: `./scripts/findings/functions/400.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:105:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("=" * 80)

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 401 — BP-PY-46

- Function context: `./scripts/findings/functions/401.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:122:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"\nFound {len(y_groups)} distinct Y positions (±5 units)")
    print("\nLines with wide X range (potential multi-column):")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 402 — BP-PY-46

- Function context: `./scripts/findings/functions/402.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:123:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("\nLines with wide X range (potential multi-column):")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 404 — BP-PY-46

- Function context: `./scripts/findings/functions/404.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:133:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"  Y={y:6.1f}: X range {x_min:6.1f} - {x_max:6.1f} (width: {x_range:6.1f})")
            print(f"            Text: {text[:80]}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 405 — BP-PY-46

- Function context: `./scripts/findings/functions/405.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/diagnose_character_order.py:134:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"            Text: {text[:80]}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 406 — BP-PY-46

- Function context: `./scripts/findings/functions/406.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_arxiv.py:30:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"Error fetching from ArXiv: {e}")
        return None
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 407 — BP-PY-46

- Function context: `./scripts/findings/functions/407.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_arxiv.py:59:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"Error parsing XML: {e}")
        return []
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 408 — BP-PY-46

- Function context: `./scripts/findings/functions/408.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_arxiv.py:68:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  Skipping {filename} (already exists)")
        return True
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 409 — BP-PY-46

- Function context: `./scripts/findings/functions/409.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_arxiv.py:72:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  Downloading {filename}...")
        urllib.request.urlretrieve(url, output_file)
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 412 — BP-PY-46

- Function context: `./scripts/findings/functions/412.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_arxiv.py:76:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  Error downloading {filename}: {e}")
        return False
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 413 — BP-PY-46

- Function context: `./scripts/findings/functions/413.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_cfr_regulations.py:67:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  Skipping Title {title_num} (already exists)")
        return True
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 414 — BP-PY-46

- Function context: `./scripts/findings/functions/414.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_cfr_regulations.py:71:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  Downloading Title {title_num}: {title_name}...")
        headers = {"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) PDF Library Testing"}
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 415 — BP-PY-46

- Function context: `./scripts/findings/functions/415.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_cfr_regulations.py:80:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print("    Warning: Not a PDF file")
            return False
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 416 — BP-PY-46

- Function context: `./scripts/findings/functions/416.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_cfr_regulations.py:86:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"    Downloaded {len(data) // 1024 // 1024} MB")
        return True
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 417 — BP-PY-46

- Function context: `./scripts/findings/functions/417.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_cfr_regulations.py:91:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"    Not found for {year} (may not have vol1)")
        else:
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 418 — BP-PY-46

- Function context: `./scripts/findings/functions/418.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_cfr_regulations.py:93:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"    HTTP Error {e.code}")
        return False
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 421 — BP-PY-46

- Function context: `./scripts/findings/functions/421.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_cfr_regulations.py:96:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"    Error: {e}")
        return False
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 425 — BP-PY-46

- Function context: `./scripts/findings/functions/425.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_chronicling_america.py:81:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"    Error fetching pages: {e}")
        return []
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 427 — BP-PY-46

- Function context: `./scripts/findings/functions/427.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_diverse_pdfs.py:106:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  Skipping {source['name']} (already exists)")
        return True
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 428 — BP-PY-46

- Function context: `./scripts/findings/functions/428.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_diverse_pdfs.py:110:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  Downloading {source['name']} ({source['type']})...")
        headers = {"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36"}
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 429 — BP-PY-46

- Function context: `./scripts/findings/functions/429.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_diverse_pdfs.py:119:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print("    Warning: Not a PDF file")
            return False
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 430 — BP-PY-46

- Function context: `./scripts/findings/functions/430.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_diverse_pdfs.py:125:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"    Downloaded {len(data) // 1024} KB")
        return True
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 433 — BP-PY-46

- Function context: `./scripts/findings/functions/433.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_diverse_pdfs.py:129:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"    Error: {e}")
        return False
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 434 — BP-PY-46

- Function context: `./scripts/findings/functions/434.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_financial.py:148:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  Skipping {source['name']} (already exists)")
        return True
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 435 — BP-PY-46

- Function context: `./scripts/findings/functions/435.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_financial.py:152:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  Downloading {source['company']} {source['type']}...")
        headers = {"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) PDF Library Testing"}
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 436 — BP-PY-46

- Function context: `./scripts/findings/functions/436.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_financial.py:161:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print("    Warning: Not a PDF file, skipping")
            return False
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 437 — BP-PY-46

- Function context: `./scripts/findings/functions/437.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_financial.py:167:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"    Downloaded {len(data) // 1024} KB")
        return True
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 438 — BP-PY-46

- Function context: `./scripts/findings/functions/438.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_financial.py:171:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"    HTTP Error {e.code}: {e.reason}")
        if e.code == 404:
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 439 — BP-PY-46

- Function context: `./scripts/findings/functions/439.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_financial.py:173:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print("    (File may have been moved or archived)")
        return False
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 442 — BP-PY-46

- Function context: `./scripts/findings/functions/442.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_financial.py:176:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"    Error: {e}")
        if output_file.exists():
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 443 — BP-PY-46

- Function context: `./scripts/findings/functions/443.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_government.py:77:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  Skipping {source['name']} (already exists)")
        return True
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 444 — BP-PY-46

- Function context: `./scripts/findings/functions/444.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_government.py:81:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  Downloading {source['name']} ({source['category']})...")
        headers = {"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36"}
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 445 — BP-PY-46

- Function context: `./scripts/findings/functions/445.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_government.py:93:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"    Warning: {source['name']} may not be a valid PDF")
            output_file.unlink()
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 446 — BP-PY-46

- Function context: `./scripts/findings/functions/446.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_government.py:97:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"    Downloaded {len(data) // 1024} KB")
        return True
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 449 — BP-PY-46

- Function context: `./scripts/findings/functions/449.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_government.py:101:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"    Error: {e}")
        if output_file.exists():
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 453 — BP-PY-46

- Function context: `./scripts/findings/functions/453.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_govinfo_policies.py:71:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"    Error fetching packages: {e}")
        return []
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 456 — BP-PY-46

- Function context: `./scripts/findings/functions/456.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_govinfo_policies.py:134:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("\nTrying direct Federal Register download...")
    successful = 0
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 457 — BP-PY-46

- Function context: `./scripts/findings/functions/457.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_govinfo_policies.py:170:21`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                    print(f"  ✓ Federal Register {year} Issue {issue_num} ({size_mb} MB)")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 462 — BP-PY-46

- Function context: `./scripts/findings/functions/462.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_govinfo_policies.py:177:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print(f"  ✗ Error: {e}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 463 — BP-PY-46

- Function context: `./scripts/findings/functions/463.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_internet_archive_newspapers.py:21:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("Installing internetarchive...")
    import subprocess
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 464 — CWE-88

- Function context: `./scripts/findings/functions/464.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_internet_archive_newspapers.py:24:5`
- Checklist pattern: `no externally influenced input — fixed literals / internal repo paths`

Source excerpt:

```
    subprocess.check_call([sys.executable, "-m", "pip", "install", "internetarchive"])
    import internetarchive as ia
```

Why this is a false positive: All argv elements are fixed literals except `sys.executable` (the interpreter path, always position 0 and never an option); no externally influenced input is embedded, so the rule condition is not met.

Checklist evidence: no externally influenced input — fixed literals / internal repo paths — All argv elements are fixed literals except `sys.executable` (the interpreter path, always position 0 and never an option); no externally influenced input is embedded, so the rule condition is not met.

### [ ] Finding 465 — BP-PY-46

- Function context: `./scripts/findings/functions/465.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_internet_archive_newspapers.py:39:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"\nSearching: '{query}'")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 466 — BP-PY-46

- Function context: `./scripts/findings/functions/466.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_internet_archive_newspapers.py:52:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"  Checking: {title}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 467 — BP-PY-46

- Function context: `./scripts/findings/functions/467.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_internet_archive_newspapers.py:65:21`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                    print("    No PDFs found")
                    continue
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 468 — BP-PY-46

- Function context: `./scripts/findings/functions/468.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_internet_archive_newspapers.py:77:21`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                    print("    Already exists")
                    continue
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 469 — BP-PY-46

- Function context: `./scripts/findings/functions/469.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_internet_archive_newspapers.py:80:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print(f"    Downloading {filename}...")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 470 — BP-PY-46

- Function context: `./scripts/findings/functions/470.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_internet_archive_newspapers.py:91:21`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                    print(f"    ✓ Downloaded ({size_mb} MB)")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 473 — BP-PY-46

- Function context: `./scripts/findings/functions/473.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_internet_archive_newspapers.py:94:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print(f"    ✗ Error: {e}")
                continue
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 475 — BP-PY-46

- Function context: `./scripts/findings/functions/475.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_internet_archive_newspapers.py:100:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  Search error: {e}")
        return 0
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 476 — BP-PY-46

- Function context: `./scripts/findings/functions/476.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_irs_forms.py:95:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  Skipping {form_name} (already exists)")
        return True
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 477 — BP-PY-46

- Function context: `./scripts/findings/functions/477.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_irs_forms.py:99:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  Downloading Form {form_name}...")
        headers = {"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) PDF Library Testing"}
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 478 — BP-PY-46

- Function context: `./scripts/findings/functions/478.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_irs_forms.py:108:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print("    Warning: Not a PDF file")
            return False
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 479 — BP-PY-46

- Function context: `./scripts/findings/functions/479.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_irs_forms.py:114:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"    Downloaded {len(data) // 1024} KB")
        return True
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 480 — BP-PY-46

- Function context: `./scripts/findings/functions/480.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_irs_forms.py:119:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"    Not found (form may not exist for {year})")
        else:
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 481 — BP-PY-46

- Function context: `./scripts/findings/functions/481.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_irs_forms.py:121:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"    HTTP Error {e.code}")
        return False
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 484 — BP-PY-46

- Function context: `./scripts/findings/functions/484.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_irs_forms.py:124:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"    Error: {e}")
        return False
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 487 — BP-PY-46

- Function context: `./scripts/findings/functions/487.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_sec_edgar_bulk.py:21:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("Installing sec-edgar-downloader...")
    import subprocess
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 488 — CWE-88

- Function context: `./scripts/findings/functions/488.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_sec_edgar_bulk.py:24:5`
- Checklist pattern: `no externally influenced input — fixed literals / internal repo paths`

Source excerpt:

```
    subprocess.check_call([sys.executable, "-m", "pip", "install", "sec-edgar-downloader"])
    from sec_edgar_downloader import Downloader
```

Why this is a false positive: All argv elements are fixed literals except `sys.executable` (the interpreter path, always position 0 and never an option); no externally influenced input is embedded, so the rule condition is not met.

Checklist evidence: no externally influenced input — fixed literals / internal repo paths — All argv elements are fixed literals except `sys.executable` (the interpreter path, always position 0 and never an option); no externally influenced input is embedded, so the rule condition is not met.

### [ ] Finding 489 — BP-PY-46

- Function context: `./scripts/findings/functions/489.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_sec_edgar_bulk.py:83:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"\n{company_name} ({ticker})")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 490 — BP-PY-46

- Function context: `./scripts/findings/functions/490.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_sec_edgar_bulk.py:92:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"  Downloading {filing_type} filings...")
            # Download recent filings (after 2020 for more PDFs)
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 491 — BP-PY-46

- Function context: `./scripts/findings/functions/491.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_sec_edgar_bulk.py:96:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"    ✓ Downloaded up to {limit_per_type} {filing_type} filings")
        except Exception as e:
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 494 — BP-PY-46

- Function context: `./scripts/findings/functions/494.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_sec_edgar_bulk.py:98:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"    ✗ Error: {e}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 497 — BP-PY-46

- Function context: `./scripts/findings/functions/497.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_sec_filings.py:85:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"Error fetching filings for CIK {cik}: {e}")
        return []
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 498 — BP-PY-46

- Function context: `./scripts/findings/functions/498.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_sec_filings.py:97:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  Skipping {filename} (already exists)")
        return True
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 499 — BP-PY-46

- Function context: `./scripts/findings/functions/499.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_sec_filings.py:101:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  Downloading {filename}...")
        headers = {"User-Agent": "PDF Library Testing bot@example.com"}
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 501 — BP-PY-46

- Function context: `./scripts/findings/functions/501.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/download_sec_filings.py:114:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  Error downloading {filename}: {e}")
        return False
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 504 — BP-PY-46

- Function context: `./scripts/findings/functions/504.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/export_text_comparison.py:48:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"WARNING: corpus '{corpus_name}' not found at {corpus_path}", file=sys.stderr)
            continue
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 505 — BP-PY-7

- Function context: `./scripts/findings/functions/505.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/export_text_comparison.py:78:18`
- Checklist pattern: `wrong-construct match — library method/module function or a definition, not builtin open()`

Source excerpt:

```
    doc = pymupdf.open(pdf_path)
    if doc.needs_pass:
```

Why this is a false positive: The matched `open` is not the builtin file-open function: it is a library/module method call (`fitz.open`, `pymupdf.open`, `AsyncPdfDocument.open`) or a function definition (`def open`). No builtin file handle is opened without a `with` statement, so the rule condition is not met.

Checklist evidence: wrong-construct match — library method/module function or a definition, not builtin open() — The matched `open` is not the builtin file-open function: it is a library/module method call (`fitz.open`, `pymupdf.open`, `AsyncPdfDocument.open`) or a function definition (`def open`). No builtin file handle is opened without a `with` statement, so the rule condition is not met.

### [ ] Finding 510 — PERF-PY-27

- Function context: `./scripts/findings/functions/510.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/fetch_real_fixtures.py:33:1`
- Checklist pattern: `distinct paths per iteration — same path is not loaded repeatedly`

Source excerpt:

```
    if out.is_file() and out.read_bytes()[:5] == b"%PDF-":
        print(f"already present: {out}")
```

Why this is a false positive: Each iteration loads a distinct path derived from the loop element (`out`/`out_p`/`path` vary per item); the same path is not loaded repeatedly, so the condition is not met.

Checklist evidence: distinct paths per iteration — same path is not loaded repeatedly — Each iteration loads a distinct path derived from the loop element (`out`/`out_p`/`path` vary per item); the same path is not loaded repeatedly, so the condition is not met.

### [ ] Finding 511 — BP-PY-46

- Function context: `./scripts/findings/functions/511.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/fetch_real_fixtures.py:34:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"already present: {out}")
        return
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 512 — BP-PY-46

- Function context: `./scripts/findings/functions/512.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/fetch_real_fixtures.py:37:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"fetching {pmcid} -> {out}")
    req = urllib.request.Request(url, headers={"User-Agent": "pdf_oxide-fixture-fetch"})
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 513 — BP-PY-7

- Function context: `./scripts/findings/functions/513.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/inspect_tj_array.py:13:11`
- Checklist pattern: `wrong-construct match — library method/module function or a definition, not builtin open()`

Source excerpt:

```
doc = fitz.open(PDF_PATH)
page = doc[0]
```

Why this is a false positive: The matched `open` is not the builtin file-open function: it is a library/module method call (`fitz.open`, `pymupdf.open`, `AsyncPdfDocument.open`) or a function definition (`def open`). No builtin file handle is opened without a `with` statement, so the rule condition is not met.

Checklist evidence: wrong-construct match — library method/module function or a definition, not builtin open() — The matched `open` is not the builtin file-open function: it is a library/module method call (`fitz.open`, `pymupdf.open`, `AsyncPdfDocument.open`) or a function definition (`def open`). No builtin file handle is opened without a `with` statement, so the rule condition is not met.

### [ ] Finding 516 — BP-PY-46

- Function context: `./scripts/findings/functions/516.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/inspect_tj_array.py:34:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print("=" * 80)
print("SEARCHING FOR TJ ARRAYS IN CONTENT STREAM")
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 517 — BP-PY-46

- Function context: `./scripts/findings/functions/517.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/inspect_tj_array.py:35:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print("SEARCHING FOR TJ ARRAYS IN CONTENT STREAM")
print("=" * 80)
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 518 — BP-PY-46

- Function context: `./scripts/findings/functions/518.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/inspect_tj_array.py:36:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print("=" * 80)

```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 519 — BP-PY-46

- Function context: `./scripts/findings/functions/519.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/inspect_tj_array.py:42:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print(f"\nFound {len(matches)} TJ arrays\n")

```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 520 — BP-PY-46

- Function context: `./scripts/findings/functions/520.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/inspect_tj_array.py:46:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"TJ Array #{i + 1}:")
    print(f"  Content: [{match[:200]}{'...' if len(match) > 200 else ''}]")
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 521 — BP-PY-46

- Function context: `./scripts/findings/functions/521.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/inspect_tj_array.py:47:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  Content: [{match[:200]}{'...' if len(match) > 200 else ''}]")

```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 522 — BP-PY-46

- Function context: `./scripts/findings/functions/522.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/inspect_tj_array.py:54:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  Strings: {len(strings)}")
    if strings:
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 523 — BP-PY-46

- Function context: `./scripts/findings/functions/523.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/inspect_tj_array.py:56:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  First 5 strings: {strings[:5]}")
    print(f"  Numbers: {len(numbers)}")
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 524 — BP-PY-46

- Function context: `./scripts/findings/functions/524.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/inspect_tj_array.py:57:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  Numbers: {len(numbers)}")
    if numbers:
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 525 — BP-PY-46

- Function context: `./scripts/findings/functions/525.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/inspect_tj_array.py:59:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  First 5 numbers: {numbers[:5]}")
    print()
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 526 — BP-PY-46

- Function context: `./scripts/findings/functions/526.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/inspect_tj_array.py:60:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print()

```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 527 — BP-PY-46

- Function context: `./scripts/findings/functions/527.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/move_large_pdfs.py:61:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"Moving {pdf_file} -> {dest}")
            shutil.move(str(pdf_file), str(dest))
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 528 — BP-PY-46

- Function context: `./scripts/findings/functions/528.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/move_large_pdfs.py:67:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print(f"\n✓ Moved {moved_count} large PDFs to {large_dir}")

```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 529 — BP-PY-46

- Function context: `./scripts/findings/functions/529.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/move_large_pdfs.py:70:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"\n⚠ Could not find {len(not_found)} files:")
    for filename in not_found:
```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 530 — BP-PY-46

- Function context: `./scripts/findings/functions/530.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/move_large_pdfs.py:72:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"  - {filename}")

```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 531 — BP-PY-46

- Function context: `./scripts/findings/functions/531.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/move_large_pdfs.py:76:1`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
print(f"\n📊 Remaining PDFs in corpus: {len(remaining)}")

```

Why this is a false positive: Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Top-level standalone script; the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 532 — BP-PY-46

- Function context: `./scripts/findings/functions/532.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:27:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"Backing up old results to: {backup_dir}")
        shutil.copytree(pdf_lib_dir, backup_dir)
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 533 — BP-PY-46

- Function context: `./scripts/findings/functions/533.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:37:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("Cleaning old pdf_oxide results...")
        shutil.rmtree(pdf_lib_dir)
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 534 — BP-PY-46

- Function context: `./scripts/findings/functions/534.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:40:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"✓ Cleaned: {pdf_lib_dir}")
    else:
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 535 — BP-PY-46

- Function context: `./scripts/findings/functions/535.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:43:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"✓ Created: {pdf_lib_dir}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 536 — BP-PY-46

- Function context: `./scripts/findings/functions/536.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:61:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"\n{'=' * 60}")
    print("Benchmarking: pdf_oxide (UPDATED)")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 537 — BP-PY-46

- Function context: `./scripts/findings/functions/537.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:62:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("Benchmarking: pdf_oxide (UPDATED)")
    print(f"{'=' * 60}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 538 — BP-PY-46

- Function context: `./scripts/findings/functions/538.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:63:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"{'=' * 60}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 539 — BP-PY-46

- Function context: `./scripts/findings/functions/539.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:102:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(
                f"  [{i}/{len(pdf_files)}] ✓ {pdf_path.name} ({elapsed:.3f}s, {output_size:,} bytes)"
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 542 — BP-PY-46

- Function context: `./scripts/findings/functions/542.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:110:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"  [{i}/{len(pdf_files)}] ✗ {pdf_path.name} - {str(e)[:100]}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 543 — BP-PY-46

- Function context: `./scripts/findings/functions/543.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:122:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("\n✅ Results for pdf_oxide (UPDATED):")
    print(
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 544 — BP-PY-46

- Function context: `./scripts/findings/functions/544.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:123:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(
        f"  Success: {results['successful']}/{results['total_pdfs']} ({results['success_rate']:.1f}%)"
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 545 — BP-PY-46

- Function context: `./scripts/findings/functions/545.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:126:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  Total time: {results['total_time']:.2f}s")
    print(f"  Avg time/PDF: {results['avg_time'] * 1000:.1f}ms")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 546 — BP-PY-46

- Function context: `./scripts/findings/functions/546.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:127:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  Avg time/PDF: {results['avg_time'] * 1000:.1f}ms")
    print(f"  Total output: {results['total_output_size']:,} bytes")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 547 — BP-PY-46

- Function context: `./scripts/findings/functions/547.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:128:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"  Total output: {results['total_output_size']:,} bytes")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 548 — BP-PY-46

- Function context: `./scripts/findings/functions/548.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:136:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"No existing results found at {summary_file}")
        return []
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 549 — BP-PY-46

- Function context: `./scripts/findings/functions/549.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:145:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"\n📊 Loaded existing results for {len(other_results)} libraries:")
    for r in other_results:
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 550 — BP-PY-46

- Function context: `./scripts/findings/functions/550.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:147:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(
            f"  - {r['library']}: {r['successful']}/{r['total_pdfs']} PDFs, {r['avg_time'] * 1000:.1f}ms avg"
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 551 — BP-PY-46

- Function context: `./scripts/findings/functions/551.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:163:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"\n✅ Saved merged results to: {summary_file}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 552 — BP-PY-46

- Function context: `./scripts/findings/functions/552.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:166:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"\n{'=' * 60}")
    print("UPDATED BENCHMARK COMPARISON")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 553 — BP-PY-46

- Function context: `./scripts/findings/functions/553.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:167:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print("UPDATED BENCHMARK COMPARISON")
    print(f"{'=' * 60}\n")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 554 — BP-PY-46

- Function context: `./scripts/findings/functions/554.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:168:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"{'=' * 60}\n")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 555 — BP-PY-46

- Function context: `./scripts/findings/functions/555.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:173:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"{'Library':<20} {'Success':<12} {'Total Time':<12} {'Avg/PDF':<12} {'Output Size':<15}")
    print(f"{'-' * 20} {'-' * 12} {'-' * 12} {'-' * 12} {'-' * 15}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 556 — BP-PY-46

- Function context: `./scripts/findings/functions/556.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:174:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"{'-' * 20} {'-' * 12} {'-' * 12} {'-' * 12} {'-' * 15}")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 557 — BP-PY-46

- Function context: `./scripts/findings/functions/557.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:181:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(
            f"{lib_name:<20} "
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 558 — BP-PY-46

- Function context: `./scripts/findings/functions/558.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:192:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"\n{'=' * 60}")
        print(f"RELATIVE PERFORMANCE (vs {baseline['library']})")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 559 — BP-PY-46

- Function context: `./scripts/findings/functions/559.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:193:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"RELATIVE PERFORMANCE (vs {baseline['library']})")
        print(f"{'=' * 60}\n")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 560 — BP-PY-46

- Function context: `./scripts/findings/functions/560.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:194:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"{'=' * 60}\n")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 561 — BP-PY-46

- Function context: `./scripts/findings/functions/561.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/rebenchmark_pdf_library.py:202:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print(f"{lib_name:<20} {speedup:>6.2f}× slower")

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 562 — CWE-88

- Function context: `./scripts/findings/functions/562.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/regression_harness.py:111:16`
- Checklist pattern: `no externally influenced input — fixed literals / internal repo paths`

Source excerpt:

```
        proc = subprocess.run(
            [str(binary), str(pdf)],
```

Why this is a false positive: The argv operands are internal repository paths (built binary path + corpus scan of the repo's own test files); no externally influenced input reaches the argument vector.

Checklist evidence: no externally influenced input — fixed literals / internal repo paths — The argv operands are internal repository paths (built binary path + corpus scan of the repo's own test files); no externally influenced input reaches the argument vector.

### [ ] Finding 572 — CWE-1341

- Function context: `./scripts/findings/functions/572.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/regression_harness.py:250:21`
- Checklist pattern: `distinct handles — two different objects are closed, no double release`

Source excerpt:

```
                    tp.close()
                with contextlib.suppress(Exception):
```

Why this is a false positive: The two adjacent `.close()` calls release two different handles (textpage/page, tp/page); no single resource handle is released twice.

Checklist evidence: distinct handles — two different objects are closed, no double release — The two adjacent `.close()` calls release two different handles (textpage/page, tp/page); no single resource handle is released twice.

### [ ] Finding 586 — BP-PY-46

- Function context: `./scripts/findings/functions/586.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/regression_harness.py:657:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print(f"[warn] unavailable for {fmt}: {', '.join(unavailable)}", file=sys.stderr)

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 587 — BP-PY-46

- Function context: `./scripts/findings/functions/587.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/regression_harness.py:662:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print(f"[{idx:>3}/{total}] {bucket:<16} {pdf}", flush=True)

```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 590 — BP-PY-46

- Function context: `./scripts/findings/functions/590.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/regression_harness.py:761:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"[done] manifest: {manifest_path}")
    return manifest
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 592 — CWE-88

- Function context: `./scripts/findings/functions/592.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/regtest_branch_vs_main.py:193:13`
- Checklist pattern: `no externally influenced input — fixed literals / internal repo paths`

Source excerpt:

```
        r = subprocess.run(
            [bin_path, pdf, fmt],
```

Why this is a false positive: The argv operands are internal repository paths (built binary path + corpus scan of the repo's own test files); no externally influenced input reaches the argument vector.

Checklist evidence: no externally influenced input — fixed literals / internal repo paths — The argv operands are internal repository paths (built binary path + corpus scan of the repo's own test files); no externally influenced input reaches the argument vector.

### [ ] Finding 600 — PERF-PY-28

- Function context: `./scripts/findings/functions/600.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/regtest_branch_vs_main.py:270:53`
- Checklist pattern: `not per unit of work — executor created once per batch/object lifetime`

Source excerpt:

```
    with open(manifest_path, "a") as mf, ProcessPoolExecutor(max_workers=jobs) as ex:
        futs = {ex.submit(_worker, t): t for t in tasks}
```

Why this is a false positive: The executor is created once at the top of the batch function (or once per long-lived wrapper object in `__init__`); it is not created per unit of work.

Checklist evidence: not per unit of work — executor created once per batch/object lifetime — The executor is created once at the top of the batch function (or once per long-lived wrapper object in `__init__`); it is not created per unit of work.

### [ ] Finding 601 — PERF-PY-23

- Function context: `./scripts/findings/functions/601.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/regtest_branch_vs_main.py:284:1`
- Checklist pattern: `unavoidable per-item encode — data produced inside the loop`

Source excerpt:

```
            mf.write(json.dumps(row) + "\n")
            mf.flush()
```

Why this is a false positive: Each iteration serializes freshly produced per-task result data with `json.dumps`; the encode cannot be hoisted out of the loop, so no avoidable polymorphic-encode work exists.

Checklist evidence: unavoidable per-item encode — data produced inside the loop — Each iteration serializes freshly produced per-task result data with `json.dumps`; the encode cannot be hoisted out of the loop, so no avoidable polymorphic-encode work exists.

### [ ] Finding 603 — PERF-PY-27

- Function context: `./scripts/findings/functions/603.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/regtest_branch_vs_main.py:404:1`
- Checklist pattern: `distinct paths per iteration — same path is not loaded repeatedly`

Source excerpt:

```
                return out_p.read_text(errors="replace")
            return ""
```

Why this is a false positive: Each iteration loads a distinct path derived from the loop element (`out`/`out_p`/`path` vary per item); the same path is not loaded repeatedly, so the condition is not met.

Checklist evidence: distinct paths per iteration — same path is not loaded repeatedly — Each iteration loads a distinct path derived from the loop element (`out`/`out_p`/`path` vary per item); the same path is not loaded repeatedly, so the condition is not met.

### [ ] Finding 604 — CWE-1046

- Function context: `./scripts/findings/functions/604.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/regtest_branch_vs_main.py:473:1`
- Checklist pattern: `no accumulating concatenation — accumulator re-initialized per iteration`

Source excerpt:

```
        flag += " ← expected improvement" if st in ("CRASH_MAIN", "NEW_OUTPUT") else ""
        md_lines.append(f"| {st} | {cols} | {tc} |{flag}")
```

Why this is a false positive: `flag` is re-initialized at the top of each loop iteration and the `+=` appends a constant suffix at most once; the accumulator does not grow by repeated concatenation inside the loop.

Checklist evidence: no accumulating concatenation — accumulator re-initialized per iteration — `flag` is re-initialized at the top of each loop iteration and the `+=` appends a constant suffix at most once; the accumulator does not grow by repeated concatenation inside the loop.

### [ ] Finding 606 — PERF-PY-27

- Function context: `./scripts/findings/functions/606.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/regtest_branch_vs_main.py:568:1`
- Checklist pattern: `distinct paths per iteration — same path is not loaded repeatedly`

Source excerpt:

```
            print(out_p.read_text(errors="replace")[:3000])

```

Why this is a false positive: Each iteration loads a distinct path derived from the loop element (`out`/`out_p`/`path` vary per item); the same path is not loaded repeatedly, so the condition is not met.

Checklist evidence: distinct paths per iteration — same path is not loaded repeatedly — Each iteration loads a distinct path derived from the loop element (`out`/`out_p`/`path` vary per item); the same path is not loaded repeatedly, so the condition is not met.

### [ ] Finding 607 — PERF-PY-27

- Function context: `./scripts/findings/functions/607.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/sync_version.py:120:1`
- Checklist pattern: `distinct paths per iteration — same path is not loaded repeatedly`

Source excerpt:

```
        text = path.read_text()
        matches = list(re.finditer(pat, text))
```

Why this is a false positive: Each iteration loads a distinct path derived from the loop element (`out`/`out_p`/`path` vary per item); the same path is not loaded repeatedly, so the condition is not met.

Checklist evidence: distinct paths per iteration — same path is not loaded repeatedly — Each iteration loads a distinct path derived from the loop element (`out`/`out_p`/`path` vary per item); the same path is not loaded repeatedly, so the condition is not met.

### [ ] Finding 608 — BP-PY-46

- Function context: `./scripts/findings/functions/608.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/sync_version.py:144:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"canonical version (Cargo.toml): {ver}")
    if missing:
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 609 — BP-PY-46

- Function context: `./scripts/findings/functions/609.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/sync_version.py:146:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("skipped (absent): " + ", ".join(missing))
    if check:
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 610 — BP-PY-46

- Function context: `./scripts/findings/functions/610.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/sync_version.py:149:13`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
            print("OUT OF SYNC:")
            for d in drift:
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 611 — BP-PY-46

- Function context: `./scripts/findings/functions/611.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/sync_version.py:151:17`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
                print("  " + d)
            return 1
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 612 — BP-PY-46

- Function context: `./scripts/findings/functions/612.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/sync_version.py:153:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("all bindings in sync ✓")
        return 0
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 613 — BP-PY-46

- Function context: `./scripts/findings/functions/613.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/sync_version.py:156:9`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
        print("  " + c)
    print(f"synced {len(changed)} file(s) to {ver}")
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

### [ ] Finding 614 — BP-PY-46

- Function context: `./scripts/findings/functions/614.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/scripts/sync_version.py:157:5`
- Checklist pattern: `CLI script output — print in a standalone script, not a non-script module`

Source excerpt:

```
    print(f"synced {len(changed)} file(s) to {ver}")
    return 0
```

Why this is a false positive: Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

Checklist evidence: CLI script output — print in a standalone script, not a non-script module — Standalone CLI script (has `if __name__ == "__main__"` guard and/or argparse entrypoint); the print is user-facing CLI output, not library logging, so the rule condition (`print` used for operational logging in non-script modules) is not met.

## Uncertain findings

### [ ] Finding 632 — CWE-88

- Function context: `./scripts/findings/functions/632.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pdf_oxide/tools/benchmark-harness/olmocr/run_olmocr_bench.py:49:13`

Source excerpt:

```
def extract_text(python_bin, pdf_path, timeout):
    code = ("import sys,pdf_oxide;")
    r = subprocess.run(
        [python_bin, "-c", code, str(pdf_path)], capture_output=True, timeout=timeout
```

Why it is uncertain: `python_bin` is supplied via the `--python` CLI flag (default `sys.executable`) and `pdf_path` derives from a locally downloaded JSONL corpus; whether a CLI flag of a local benchmark harness counts as "externally influenced input" reaching an argument vector depends on the deployment assumption, which the shown source does not settle.

## True positives

Findings whose shown source satisfies the rule condition; listed compactly per rule.

### BP-PY-1 — 91 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 2 | `examples/ocr_example.py:109:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 4 | `examples/ocr_example.py:119:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 5 | `examples/ocr_example.py:146:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 11 | `examples/python/08-batch-processing/main.py:37:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 23 | `python/pdf_oxide/__init__.py:109:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 32 | `python/tests/test_api_coverage.py:472:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 41 | `python/tests/test_feature_guard.py:107:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 128 | `scripts/analyze_pdf_spacing.py:154:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 131 | `scripts/bench_compare.py:43:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 139 | `scripts/bench_pymupdf.py:100:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 145 | `scripts/benchmark_all_libraries.py:323:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 147 | `scripts/benchmark_pdf_library.py:127:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 156 | `scripts/browser_download_pdfs.py:87:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 159 | `scripts/browser_download_pdfs.py:90:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 165 | `scripts/browser_download_pdfs.py:167:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 167 | `scripts/browser_download_pdfs.py:171:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 173 | `scripts/browser_download_pdfs.py:237:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 175 | `scripts/browser_download_pdfs.py:241:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 191 | `scripts/check_pdf_structure.py:65:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 194 | `scripts/check_pdf_structure.py:88:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 210 | `scripts/compare_all_outputs.py:32:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 212 | `scripts/compare_all_outputs.py:38:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 255 | `scripts/compare_extractors.py:36:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 259 | `scripts/compare_extractors.py:52:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 261 | `scripts/compare_extractors.py:67:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 263 | `scripts/compare_extractors.py:82:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 265 | `scripts/compare_extractors.py:101:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 271 | `scripts/convert_models.py:36:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 287 | `scripts/convert_models.py:106:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 304 | `scripts/convert_models.py:167:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 328 | `scripts/convert_pdf_spec.py:23:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 352 | `scripts/debug_font_encoding.py:69:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 358 | `scripts/debug_font_encoding.py:75:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 360 | `scripts/debug_font_encoding.py:77:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 410 | `scripts/download_arxiv.py:75:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 419 | `scripts/download_cfr_regulations.py:95:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 423 | `scripts/download_chronicling_america.py:80:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 426 | `scripts/download_chronicling_america.py:121:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 431 | `scripts/download_diverse_pdfs.py:128:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 440 | `scripts/download_financial.py:175:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 447 | `scripts/download_government.py:100:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 451 | `scripts/download_govinfo_policies.py:70:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 454 | `scripts/download_govinfo_policies.py:122:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 455 | `scripts/download_govinfo_policies.py:125:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 461 | `scripts/download_govinfo_policies.py:176:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 471 | `scripts/download_internet_archive_newspapers.py:93:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 474 | `scripts/download_internet_archive_newspapers.py:99:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 482 | `scripts/download_irs_forms.py:123:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 485 | `scripts/download_oatd_theses.py:102:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 492 | `scripts/download_sec_edgar_bulk.py:97:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 495 | `scripts/download_sec_filings.py:84:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 500 | `scripts/download_sec_filings.py:113:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 502 | `scripts/download_who_reports.py:90:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 506 | `scripts/export_text_comparison.py:100:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 508 | `scripts/export_text_comparison.py:118:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 514 | `scripts/inspect_tj_array.py:31:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 540 | `scripts/rebenchmark_pdf_library.py:106:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 563 | `scripts/regression_harness.py:118:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 565 | `scripts/regression_harness.py:123:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 566 | `scripts/regression_harness.py:160:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 567 | `scripts/regression_harness.py:173:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 568 | `scripts/regression_harness.py:195:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 569 | `scripts/regression_harness.py:206:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 570 | `scripts/regression_harness.py:214:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 571 | `scripts/regression_harness.py:233:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 573 | `scripts/regression_harness.py:256:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 574 | `scripts/regression_harness.py:263:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 575 | `scripts/regression_harness.py:271:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 576 | `scripts/regression_harness.py:282:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 577 | `scripts/regression_harness.py:296:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 578 | `scripts/regression_harness.py:303:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 579 | `scripts/regression_harness.py:311:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 580 | `scripts/regression_harness.py:576:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 584 | `scripts/regression_harness.py:581:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 588 | `scripts/regression_harness.py:703:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 589 | `scripts/regression_harness.py:726:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 591 | `scripts/regression_harness.py:801:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 593 | `scripts/regtest_branch_vs_main.py:223:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 596 | `scripts/regtest_branch_vs_main.py:253:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 602 | `scripts/regtest_branch_vs_main.py:378:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 605 | `scripts/regtest_branch_vs_main.py:559:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 617 | `scripts/test_adaptive_heuristics.py:38:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 618 | `scripts/test_adaptive_heuristics.py:55:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 619 | `scripts/test_adaptive_heuristics.py:76:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 620 | `scripts/test_arxiv_ligatures.py:77:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 624 | `scripts/test_structure_tree.py:35:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 625 | `scripts/test_structure_tree.py:61:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 628 | `scripts/test_structure_tree_integration.py:48:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 629 | `scripts/test_structure_tree_integration.py:73:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 633 | `tools/benchmark-harness/olmocr/run_olmocr_bench.py:53:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |
| 636 | `tools/benchmark-harness/olmocr/run_olmocr_bench.py:105:1` | broad `except Exception`/`except BaseException` handler that does not re-raise |

### CWE-396 — 34 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 3 | `examples/ocr_example.py:109:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 12 | `examples/python/08-batch-processing/main.py:37:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 26 | `python/pdf_oxide/__init__.py:109:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 129 | `scripts/analyze_pdf_spacing.py:154:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 134 | `scripts/bench_compare.py:43:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 140 | `scripts/bench_pymupdf.py:100:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 146 | `scripts/benchmark_all_libraries.py:323:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 148 | `scripts/benchmark_pdf_library.py:127:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 157 | `scripts/browser_download_pdfs.py:87:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 192 | `scripts/check_pdf_structure.py:65:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 211 | `scripts/compare_all_outputs.py:32:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 256 | `scripts/compare_extractors.py:36:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 272 | `scripts/convert_models.py:36:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 329 | `scripts/convert_pdf_spec.py:23:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 355 | `scripts/debug_font_encoding.py:69:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 411 | `scripts/download_arxiv.py:75:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 420 | `scripts/download_cfr_regulations.py:95:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 424 | `scripts/download_chronicling_america.py:80:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 432 | `scripts/download_diverse_pdfs.py:128:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 441 | `scripts/download_financial.py:175:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 448 | `scripts/download_government.py:100:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 452 | `scripts/download_govinfo_policies.py:70:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 472 | `scripts/download_internet_archive_newspapers.py:93:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 483 | `scripts/download_irs_forms.py:123:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 486 | `scripts/download_oatd_theses.py:102:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 493 | `scripts/download_sec_edgar_bulk.py:97:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 496 | `scripts/download_sec_filings.py:84:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 503 | `scripts/download_who_reports.py:90:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 507 | `scripts/export_text_comparison.py:100:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 515 | `scripts/inspect_tj_array.py:31:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 541 | `scripts/rebenchmark_pdf_library.py:106:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 564 | `scripts/regression_harness.py:118:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 594 | `scripts/regtest_branch_vs_main.py:223:1` | generic `except Exception:`/`except BaseException` handler in non-test code |
| 634 | `tools/benchmark-harness/olmocr/run_olmocr_bench.py:53:1` | generic `except Exception:`/`except BaseException` handler in non-test code |

### CWE-1121 — 12 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 1 | `examples/ocr_example.py:43:12` | function body counts >= 12 control-flow branch tokens (verified per function) |
| 22 | `examples/python_example.py:26:12` | function body counts >= 12 control-flow branch tokens (verified per function) |
| 43 | `scripts/analyze_pdf_layout.py:16:30` | function body counts >= 12 control-flow branch tokens (verified per function) |
| 136 | `scripts/bench_compare.py:63:70` | function body counts >= 12 control-flow branch tokens (verified per function) |
| 216 | `scripts/compare_all_outputs.py:127:35` | function body counts >= 12 control-flow branch tokens (verified per function) |
| 266 | `scripts/compare_extractors.py:124:12` | function body counts >= 12 control-flow branch tokens (verified per function) |
| 267 | `scripts/compare_results.py:106:12` | function body counts >= 12 control-flow branch tokens (verified per function) |
| 364 | `scripts/diagnose_character_order.py:16:67` | function body counts >= 12 control-flow branch tokens (verified per function) |
| 509 | `scripts/export_text_comparison.py:122:12` | function body counts >= 12 control-flow branch tokens (verified per function) |
| 595 | `scripts/regtest_branch_vs_main.py:234:19` | function body counts >= 12 control-flow branch tokens (verified per function) |
| 621 | `scripts/test_arxiv_ligatures.py:137:32` | function body counts >= 12 control-flow branch tokens (verified per function) |
| 635 | `tools/benchmark-harness/olmocr/run_olmocr_bench.py:73:12` | function body counts >= 12 control-flow branch tokens (verified per function) |

### BP-PY-41 — 12 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 36 | `python/tests/test_api_coverage.py:635:1` | test function (`def test_*`) calls code without any assertion |
| 37 | `python/tests/test_feature_guard.py:93:1` | test function (`def test_*`) calls code without any assertion |
| 254 | `scripts/compare_extractors.py:25:1` | test function (`def test_*`) calls code without any assertion |
| 257 | `scripts/compare_extractors.py:40:1` | test function (`def test_*`) calls code without any assertion |
| 260 | `scripts/compare_extractors.py:56:1` | test function (`def test_*`) calls code without any assertion |
| 262 | `scripts/compare_extractors.py:71:1` | test function (`def test_*`) calls code without any assertion |
| 264 | `scripts/compare_extractors.py:86:1` | test function (`def test_*`) calls code without any assertion |
| 615 | `scripts/test_adaptive_heuristics.py:17:1` | test function (`def test_*`) calls code without any assertion |
| 622 | `scripts/test_structure_tree.py:14:1` | test function (`def test_*`) calls code without any assertion |
| 626 | `scripts/test_structure_tree_integration.py:17:1` | test function (`def test_*`) calls code without any assertion |
| 630 | `tests/test_python.py:1333:1` | test function (`def test_*`) calls code without any assertion |
| 631 | `tests/test_python.py:1373:1` | test function (`def test_*`) calls code without any assertion |

### BP-PY-2 — 11 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 24 | `python/pdf_oxide/__init__.py:109:1` | except handler body is solely `pass` |
| 33 | `python/tests/test_api_coverage.py:472:1` | except handler body is solely `pass` |
| 38 | `python/tests/test_feature_guard.py:105:1` | except handler body is solely `pass` |
| 132 | `scripts/bench_compare.py:43:1` | except handler body is solely `pass` |
| 353 | `scripts/debug_font_encoding.py:69:1` | except handler body is solely `pass` |
| 359 | `scripts/debug_font_encoding.py:75:1` | except handler body is solely `pass` |
| 361 | `scripts/debug_font_encoding.py:77:1` | except handler body is solely `pass` |
| 458 | `scripts/download_govinfo_policies.py:174:1` | except handler body is solely `pass` |
| 581 | `scripts/regression_harness.py:576:1` | except handler body is solely `pass` |
| 585 | `scripts/regression_harness.py:581:1` | except handler body is solely `pass` |
| 597 | `scripts/regtest_branch_vs_main.py:253:1` | except handler body is solely `pass` |

### CWE-390 — 8 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 25 | `python/pdf_oxide/__init__.py:109:1` | except clause whose direct body is only `pass` |
| 34 | `python/tests/test_api_coverage.py:472:1` | except clause whose direct body is only `pass` |
| 39 | `python/tests/test_feature_guard.py:105:1` | except clause whose direct body is only `pass` |
| 133 | `scripts/bench_compare.py:43:1` | except clause whose direct body is only `pass` |
| 354 | `scripts/debug_font_encoding.py:69:1` | except clause whose direct body is only `pass` |
| 459 | `scripts/download_govinfo_policies.py:174:1` | except clause whose direct body is only `pass` |
| 582 | `scripts/regression_harness.py:576:1` | except clause whose direct body is only `pass` |
| 598 | `scripts/regtest_branch_vs_main.py:253:1` | except clause whose direct body is only `pass` |

### CWE-1071 — 8 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 27 | `python/pdf_oxide/__init__.py:109:5` | exception handler contains only `pass` (empty handler block) |
| 35 | `python/tests/test_api_coverage.py:472:9` | exception handler contains only `pass` (empty handler block) |
| 40 | `python/tests/test_feature_guard.py:105:5` | exception handler contains only `pass` (empty handler block) |
| 135 | `scripts/bench_compare.py:43:9` | exception handler contains only `pass` (empty handler block) |
| 356 | `scripts/debug_font_encoding.py:69:17` | exception handler contains only `pass` (empty handler block) |
| 460 | `scripts/download_govinfo_policies.py:174:13` | exception handler contains only `pass` (empty handler block) |
| 583 | `scripts/regression_harness.py:576:5` | exception handler contains only `pass` (empty handler block) |
| 599 | `scripts/regtest_branch_vs_main.py:253:13` | exception handler contains only `pass` (empty handler block) |

### BP-PY-42 — 4 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 31 | `python/tests/test_api_coverage.py:468:1` | test uses bare try/except to expect failure instead of assertRaises/pytest.raises |
| 616 | `scripts/test_adaptive_heuristics.py:29:1` | test uses bare try/except to expect failure instead of assertRaises/pytest.raises |
| 623 | `scripts/test_structure_tree.py:24:1` | test uses bare try/except to expect failure instead of assertRaises/pytest.raises |
| 627 | `scripts/test_structure_tree_integration.py:29:1` | test uses bare try/except to expect failure instead of assertRaises/pytest.raises |

### CWE-1124 — 4 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 113 | `scripts/analyze_pdf_spacing.py:74:1` | statement nested at >= 6 control-flow levels |
| 154 | `scripts/browser_download_pdfs.py:80:1` | statement nested at >= 6 control-flow levels |
| 422 | `scripts/download_chronicling_america.py:71:1` | statement nested at >= 6 control-flow levels |
| 450 | `scripts/download_govinfo_policies.py:66:1` | statement nested at >= 6 control-flow levels |

### CWE-367 — 1 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 20 | `examples/python/09-new-features/pkcs12_signing/main.py:22:12` | `os.path.exists(P12_PATH)` check followed within 300 chars by `open(P12_PATH, "rb")` |

### CWE-1084 — 1 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 21 | `examples/python_example.py:26:12` | function contains >= 3 `open(`/`.execute` calls (5 counted in main()) |

### BP-PY-45 — 1 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 42 | `scripts/analyze_pdf_layout.py:11:1` | `sys.path.insert` mutates sys.path at runtime |

### PERF-PY-25 — 1 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 403 | `scripts/diagnose_character_order.py:132:1` | sort-key lambda is constructed per loop iteration inside the loop |

## Final evidence

- Delegated reviewers: none (single-reviewer audit)
- Chunk evidence: `./scripts/chunks` (26 chunk files, findings 1-636)
- Function evidence: `./scripts/findings/functions` (per-finding context files)
- Validation: `git diff --check` — pass
