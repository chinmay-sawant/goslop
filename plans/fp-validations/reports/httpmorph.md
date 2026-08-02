# False-positive audit — httpmorph

## Run metadata

```yaml
timestamp: 2026-08-02T07:58:32Z
repository: httpmorph
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph
branch: main
commit: 608decb1f5b82da9d13f01920bbd50d7e1a2a196
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph
chunk_path: ./scripts/chunks
function_context_path: ./scripts/findings/functions
```

## Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir real-repos/httpmorph/scripts/chunks -context-dir real-repos/httpmorph/scripts/findings/functions real-repos/httpmorph`
- Findings: `714`
- Chunks reviewed: `./scripts/chunks/Chunk_1_25.txt` … `./scripts/chunks/Chunk_701_714.txt` (all 29 chunk files)
- Function contexts reviewed: `./scripts/findings/functions/<id>.txt` for every proposed false positive

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
| False positive | 278 | 225, 226, 227, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 247, 248, 249, 250, 251, 252, 253, 254, 255, 256, 257, 258, 259, 260, 262, 263, 264, 265, 266, 267, 268, 269, 270, 275, 276, 277, 278, 279, 280, 281, 282, 283, 284, 285, 286, 289, 290, 291, 292, 293, 294, 295, 296, 297, 298, 299, 300, 301, 303, 304, 305, 306, 307, 308, 309, 310, 312, 313, 314, 315, 316, 317, 318, 319, 320, 321, 322, 323, 324, 325, 326, 327, 328, 329, 330, 331, 332, 333, 335, 336, 337, 338, 339, 340, 341, 342, 343, 344, 345, 346, 348, 349, 350, 352, 353, 354, 356, 357, 358, 359, 360, 361, 364, 366, 367, 368, 369, 370, 372, 373, 374, 375, 376, 378, 380, 381, 382, 383, 384, 385, 386, 387, 388, 389, 390, 391, 393, 395, 396, 397, 398, 399, 400, 401, 402, 404, 406, 407, 408, 409, 410, 412, 413, 415, 416, 417, 418, 419, 420, 421, 422, 423, 424, 425, 426, 427, 428, 429, 430, 431, 432, 433, 434, 435, 436, 437, 438, 439, 440, 442, 444, 445, 446, 447, 448, 449, 454, 456, 458, 459, 461, 465, 466, 467, 468, 469, 470, 471, 472, 475, 476, 477, 478, 479, 480, 481, 500, 501, 502, 503, 504, 505, 506, 507, 508, 509, 510, 511, 512, 513, 517, 518, 519, 520, 521, 522, 523, 524, 525, 526, 527, 529, 530, 531, 532, 533, 534, 538, 539, 540, 541, 542, 543, 544, 545, 546, 547, 548, 549, 550, 551, 552, 553, 554, 555, 556, 557, 558, 559, 560, 561, 562, 563, 570, 594, 680, 684, 697 |
| True positive | 436 | 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 228, 246, 261, 271, 272, 273, 274, 287, 288, 302, 311, 334, 347, 351, 355, 362, 363, 365, 371, 377, 379, 392, 394, 403, 405, 411, 414, 441, 443, 450, 451, 452, 453, 455, 457, 460, 462, 463, 464, 473, 474, 482, 483, 484, 485, 486, 487, 488, 489, 490, 491, 492, 493, 494, 495, 496, 497, 498, 499, 514, 515, 516, 528, 535, 536, 537, 564, 565, 566, 567, 568, 569, 571, 572, 573, 574, 575, 576, 577, 578, 579, 580, 581, 582, 583, 584, 585, 586, 587, 588, 589, 590, 591, 592, 593, 595, 596, 597, 598, 599, 600, 601, 602, 603, 604, 605, 606, 607, 608, 609, 610, 611, 612, 613, 614, 615, 616, 617, 618, 619, 620, 621, 622, 623, 624, 625, 626, 627, 628, 629, 630, 631, 632, 633, 634, 635, 636, 637, 638, 639, 640, 641, 642, 643, 644, 645, 646, 647, 648, 649, 650, 651, 652, 653, 654, 655, 656, 657, 658, 659, 660, 661, 662, 663, 664, 665, 666, 667, 668, 669, 670, 671, 672, 673, 674, 675, 676, 677, 678, 679, 681, 682, 683, 685, 686, 687, 688, 689, 690, 691, 692, 693, 694, 695, 696, 698, 699, 700, 701, 702, 703, 704, 705, 706, 707, 708, 709, 710, 711, 712, 713, 714 |
| Uncertain | 0 | — |

## False positives

The 272 BP-PY-46 findings below share one checklist pattern: the flagged `print` is user-facing output of a script module (example/demo script, setuptools packaging script, or per-platform build helper), not operational logging in a non-script (library) module. BP-PY-46's own condition is "`print` is used for operational logging in non-script modules", and its fix explicitly permits print for CLIs ("keep print under `if __name__ == "__main__"` for CLIs"). The detector's exemptions only cover test files, benchmark files, prints textually inside a `__main__` guard, and argparse/Click-style CLI entrypoints — it does not recognize prints in the demo/build script modules below, so every print line in them is flagged.

### [ ] Finding 225 — BP-PY-46

- Function context: `./scripts/findings/functions/225.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:19:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\n" + "=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 226 — BP-PY-46

- Function context: `./scripts/findings/functions/226.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:20:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("1. Browser Profiles")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 227 — BP-PY-46

- Function context: `./scripts/findings/functions/227.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:21:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 229 — BP-PY-46

- Function context: `./scripts/findings/functions/229.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:29:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"\n{browser.upper()}:")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 230 — BP-PY-46

- Function context: `./scripts/findings/functions/230.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:30:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"  Status:      {response.status_code}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 231 — BP-PY-46

- Function context: `./scripts/findings/functions/231.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:31:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"  User-Agent:  {session.user_agent[:50]}...")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 232 — BP-PY-46

- Function context: `./scripts/findings/functions/232.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:32:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"  JA3 Hash:    {response.ja3_fingerprint}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 233 — BP-PY-46

- Function context: `./scripts/findings/functions/233.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:33:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"  TLS Version: {response.tls_version}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 234 — BP-PY-46

- Function context: `./scripts/findings/functions/234.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:34:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"  TLS Cipher:  {response.tls_cipher}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 235 — BP-PY-46

- Function context: `./scripts/findings/functions/235.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:39:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\n" + "=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 236 — BP-PY-46

- Function context: `./scripts/findings/functions/236.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:40:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("2. TLS Information Extraction")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 237 — BP-PY-46

- Function context: `./scripts/findings/functions/237.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:41:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 238 — BP-PY-46

- Function context: `./scripts/findings/functions/238.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:45:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\nTLS Details:")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 239 — BP-PY-46

- Function context: `./scripts/findings/functions/239.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:46:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print(f"  Version:     {response.tls_version}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 240 — BP-PY-46

- Function context: `./scripts/findings/functions/240.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:47:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print(f"  Cipher:      {response.tls_cipher}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 241 — BP-PY-46

- Function context: `./scripts/findings/functions/241.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:48:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print(f"  JA3:         {response.ja3_fingerprint}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 242 — BP-PY-46

- Function context: `./scripts/findings/functions/242.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:49:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print(f"  HTTP Ver:    {response.http_version}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 243 — BP-PY-46

- Function context: `./scripts/findings/functions/243.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:54:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\n" + "=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 244 — BP-PY-46

- Function context: `./scripts/findings/functions/244.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:55:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("3. Session Management")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 245 — BP-PY-46

- Function context: `./scripts/findings/functions/245.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:56:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 247 — BP-PY-46

- Function context: `./scripts/findings/functions/247.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:64:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print("\nGitHub API Request:")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 248 — BP-PY-46

- Function context: `./scripts/findings/functions/248.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:65:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"  Status:      {response.status_code}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 249 — BP-PY-46

- Function context: `./scripts/findings/functions/249.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:66:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"  Body length: {len(response.body)} bytes")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 250 — BP-PY-46

- Function context: `./scripts/findings/functions/250.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:67:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"  Time taken:  {response.total_time_us / 1000:.2f}ms")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 251 — BP-PY-46

- Function context: `./scripts/findings/functions/251.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:72:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\n" + "=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 252 — BP-PY-46

- Function context: `./scripts/findings/functions/252.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:73:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("4. HTTP Methods")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 253 — BP-PY-46

- Function context: `./scripts/findings/functions/253.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:74:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 254 — BP-PY-46

- Function context: `./scripts/findings/functions/254.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:78:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\nGET Request:")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 255 — BP-PY-46

- Function context: `./scripts/findings/functions/255.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:79:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print(f"  Status: {response.status_code}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 256 — BP-PY-46

- Function context: `./scripts/findings/functions/256.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:80:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print(f"  Time:   {response.total_time_us / 1000:.2f}ms")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 257 — BP-PY-46

- Function context: `./scripts/findings/functions/257.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:89:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\n" + "=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 258 — BP-PY-46

- Function context: `./scripts/findings/functions/258.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:90:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("5. Performance Metrics")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 259 — BP-PY-46

- Function context: `./scripts/findings/functions/259.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:91:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 260 — BP-PY-46

- Function context: `./scripts/findings/functions/260.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:100:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\nSequential Requests:")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 262 — BP-PY-46

- Function context: `./scripts/findings/functions/262.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:106:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"  {url:30} - {response.status_code} ({elapsed:.2f}ms)")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 263 — BP-PY-46

- Function context: `./scripts/findings/functions/263.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:111:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\n" + "=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 264 — BP-PY-46

- Function context: `./scripts/findings/functions/264.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:112:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("6. Automatic Gzip Decompression")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 265 — BP-PY-46

- Function context: `./scripts/findings/functions/265.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:113:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 266 — BP-PY-46

- Function context: `./scripts/findings/functions/266.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:117:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\nResponse:")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 267 — BP-PY-46

- Function context: `./scripts/findings/functions/267.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:118:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print(f"  Encoding:    {response.headers.get('Content-Encoding', 'none')}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 268 — BP-PY-46

- Function context: `./scripts/findings/functions/268.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:119:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print(f"  Body length: {len(response.body)} bytes")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 269 — BP-PY-46

- Function context: `./scripts/findings/functions/269.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:120:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print(f"  Text length: {len(response.text)} chars")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 270 — BP-PY-46

- Function context: `./scripts/findings/functions/270.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/advanced_features.py:121:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print(f"  Decoded:     {'Yes' if response.text else 'No'}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/advanced_features.py guarded by `if __name__ == "__main__":` at line 150); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/advanced_features.py is a script module (guarded by `if __name__ == "__main__":` at line 150); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 275 — BP-PY-46

- Function context: `./scripts/findings/functions/275.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/async_example.py:69:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\n=== Concurrent Requests Demo ===\n")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/async_example.py guarded by `if __name__ == "__main__":` at line 109); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/async_example.py is a script module (guarded by `if __name__ == "__main__":` at line 109); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 276 — BP-PY-46

- Function context: `./scripts/findings/functions/276.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/async_example.py:82:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"🚀 Making {len(urls)} concurrent requests...")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/async_example.py guarded by `if __name__ == "__main__":` at line 109); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/async_example.py is a script module (guarded by `if __name__ == "__main__":` at line 109); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 277 — BP-PY-46

- Function context: `./scripts/findings/functions/277.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/async_example.py:95:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"\n✅ All requests completed in {total_time:.2f}s\n")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/async_example.py guarded by `if __name__ == "__main__":` at line 109); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/async_example.py is a script module (guarded by `if __name__ == "__main__":` at line 109); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 278 — BP-PY-46

- Function context: `./scripts/findings/functions/278.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/async_example.py:99:17`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
                print(f"   Request {i + 1}: ❌ {result}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/async_example.py guarded by `if __name__ == "__main__":` at line 109); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/async_example.py is a script module (guarded by `if __name__ == "__main__":` at line 109); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 279 — BP-PY-46

- Function context: `./scripts/findings/functions/279.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/async_example.py:101:17`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
                print(
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/async_example.py guarded by `if __name__ == "__main__":` at line 109); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/async_example.py is a script module (guarded by `if __name__ == "__main__":` at line 109); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 280 — BP-PY-46

- Function context: `./scripts/findings/functions/280.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/async_example.py:105:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"\n💡 Note: With thread pool, this would take ~{sum([1, 2, 1])}s")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/async_example.py guarded by `if __name__ == "__main__":` at line 109); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/async_example.py is a script module (guarded by `if __name__ == "__main__":` at line 109); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 281 — BP-PY-46

- Function context: `./scripts/findings/functions/281.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/async_example.py:106:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"   With async I/O, it took {total_time:.2f}s (concurrent!)")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/async_example.py guarded by `if __name__ == "__main__":` at line 109); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/async_example.py is a script module (guarded by `if __name__ == "__main__":` at line 109); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 282 — BP-PY-46

- Function context: `./scripts/findings/functions/282.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:14:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("=== Simple GET Request ===")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 283 — BP-PY-46

- Function context: `./scripts/findings/functions/283.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:18:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"Status: {response.status_code}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 284 — BP-PY-46

- Function context: `./scripts/findings/functions/284.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:19:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"Body: {response.body[:100]}...")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 285 — BP-PY-46

- Function context: `./scripts/findings/functions/285.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:21:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print("Not yet implemented - C extension needs to be built")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 286 — BP-PY-46

- Function context: `./scripts/findings/functions/286.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:26:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\n=== Session with Persistent Fingerprint ===")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 289 — BP-PY-46

- Function context: `./scripts/findings/functions/289.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:36:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"Request 1 - Status: {response1.status_code}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 290 — BP-PY-46

- Function context: `./scripts/findings/functions/290.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:37:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"Request 2 - Status: {response2.status_code}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 291 — BP-PY-46

- Function context: `./scripts/findings/functions/291.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:38:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print("Same TLS fingerprint maintained across requests")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 292 — BP-PY-46

- Function context: `./scripts/findings/functions/292.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:40:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print("Not yet implemented - C extension needs to be built")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 293 — BP-PY-46

- Function context: `./scripts/findings/functions/293.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:45:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\n=== Chrome Browser Fingerprint ===")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 294 — BP-PY-46

- Function context: `./scripts/findings/functions/294.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:49:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print("Created Chrome session")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 295 — BP-PY-14

- Function context: `./scripts/findings/functions/295.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:50:22`
- Checklist pattern: flagged HTTP call is inside a comment, not an executable statement

Source excerpt:

```
        # response = session.get("https://httpbin.org/user-agent")
```

Why this is a false positive: The line is commented out; no HTTP call exists, so there is no call that could omit `timeout=`.

Checklist evidence: `requestsCallRe`/`sessionHTTPCallRe` match raw source text including comments; the leading `#` shows this is not an executable call.

### [ ] Finding 296 — BP-PY-46

- Function context: `./scripts/findings/functions/296.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:53:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print("  Chrome: Not yet implemented")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 297 — BP-PY-46

- Function context: `./scripts/findings/functions/297.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:58:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\n=== POST Request with JSON ===")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 298 — BP-PY-46

- Function context: `./scripts/findings/functions/298.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:67:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"Status: {response.status_code}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 299 — BP-PY-46

- Function context: `./scripts/findings/functions/299.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:68:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"Response: {response.body[:200]}...")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 300 — BP-PY-46

- Function context: `./scripts/findings/functions/300.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:70:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print("Not yet implemented - C extension needs to be built")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 301 — BP-PY-46

- Function context: `./scripts/findings/functions/301.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:75:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\n=== Custom Headers ===")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 303 — BP-PY-46

- Function context: `./scripts/findings/functions/303.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:87:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"Status: {response.status_code}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 304 — BP-PY-46

- Function context: `./scripts/findings/functions/304.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:89:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print("Not yet implemented - C extension needs to be built")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 305 — BP-PY-46

- Function context: `./scripts/findings/functions/305.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:94:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\n=== Rotating Fingerprints ===")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 306 — BP-PY-46

- Function context: `./scripts/findings/functions/306.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:102:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 307 — BP-PY-46

- Function context: `./scripts/findings/functions/307.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:106:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print("Not yet implemented - C extension needs to be built")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 308 — BP-PY-46

- Function context: `./scripts/findings/functions/308.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:111:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\n=== Performance Test ===")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 309 — BP-PY-46

- Function context: `./scripts/findings/functions/309.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:125:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"httpmorph: {iterations} requests in {fast_time:.2f}s")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 310 — BP-PY-46

- Function context: `./scripts/findings/functions/310.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:126:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"Average: {fast_time / iterations * 1000:.1f}ms per request")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 312 — BP-PY-46

- Function context: `./scripts/findings/functions/312.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:137:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"requests: {iterations} requests in {requests_time:.2f}s")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 313 — BP-PY-46

- Function context: `./scripts/findings/functions/313.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:138:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"Average: {requests_time / iterations * 1000:.1f}ms per request")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 314 — BP-PY-46

- Function context: `./scripts/findings/functions/314.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:139:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"Speedup: {requests_time / fast_time:.2f}x faster")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 315 — BP-PY-46

- Function context: `./scripts/findings/functions/315.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:141:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print("requests library not installed for comparison")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 316 — BP-PY-46

- Function context: `./scripts/findings/functions/316.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:144:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print("Not yet implemented - C extension needs to be built")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 317 — BP-PY-46

- Function context: `./scripts/findings/functions/317.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:149:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("=== httpmorph Library Information ===")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 318 — BP-PY-46

- Function context: `./scripts/findings/functions/318.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:150:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print(f"Version: {httpmorph.version()}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 319 — BP-PY-46

- Function context: `./scripts/findings/functions/319.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:151:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("Features:")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 320 — BP-PY-46

- Function context: `./scripts/findings/functions/320.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:152:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("  - High-performance C implementation")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 321 — BP-PY-46

- Function context: `./scripts/findings/functions/321.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:153:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("  - io_uring support (Linux 5.1+)")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 322 — BP-PY-46

- Function context: `./scripts/findings/functions/322.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:154:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("  - BoringSSL for TLS (Chrome-compatible)")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 323 — BP-PY-46

- Function context: `./scripts/findings/functions/323.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:155:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("  - Dynamic JA3/JA4 fingerprinting")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 324 — BP-PY-46

- Function context: `./scripts/findings/functions/324.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:156:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("  - HTTP/2 and HTTP/3 support")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 325 — BP-PY-46

- Function context: `./scripts/findings/functions/325.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:157:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("  - Zero-copy I/O operations")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 326 — BP-PY-46

- Function context: `./scripts/findings/functions/326.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/basic_usage.py:158:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print()
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/basic_usage.py guarded by `if __name__ == "__main__":` at line 161); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/basic_usage.py is a script module (guarded by `if __name__ == "__main__":` at line 161); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 327 — BP-PY-46

- Function context: `./scripts/findings/functions/327.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/http2_example.py:11:1`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
print("Example 1: Client with HTTP/2")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/http2_example.py a module-level demo script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/http2_example.py is a script module (a module-level demo script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 328 — BP-PY-46

- Function context: `./scripts/findings/functions/328.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/http2_example.py:12:1`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
print("-" * 40)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/http2_example.py a module-level demo script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/http2_example.py is a script module (a module-level demo script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 329 — BP-PY-46

- Function context: `./scripts/findings/functions/329.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/http2_example.py:17:1`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
print(f"Status: {response.status_code}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/http2_example.py a module-level demo script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/http2_example.py is a script module (a module-level demo script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 330 — BP-PY-46

- Function context: `./scripts/findings/functions/330.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/http2_example.py:18:1`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
print(f"HTTP Version: {response.http_version}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/http2_example.py a module-level demo script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/http2_example.py is a script module (a module-level demo script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 331 — BP-PY-46

- Function context: `./scripts/findings/functions/331.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/http2_example.py:19:1`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
print()
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/http2_example.py a module-level demo script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/http2_example.py is a script module (a module-level demo script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 332 — BP-PY-46

- Function context: `./scripts/findings/functions/332.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/http2_example.py:23:1`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
print("Example 2: Session with HTTP/2")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/http2_example.py a module-level demo script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/http2_example.py is a script module (a module-level demo script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 333 — BP-PY-46

- Function context: `./scripts/findings/functions/333.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/http2_example.py:24:1`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
print("-" * 40)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/http2_example.py a module-level demo script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/http2_example.py is a script module (a module-level demo script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 335 — BP-PY-46

- Function context: `./scripts/findings/functions/335.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/http2_example.py:28:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print(f"Status: {response.status_code}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/http2_example.py a module-level demo script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/http2_example.py is a script module (a module-level demo script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 336 — BP-PY-46

- Function context: `./scripts/findings/functions/336.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/http2_example.py:29:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print(f"HTTP Version: {response.http_version}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/http2_example.py a module-level demo script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/http2_example.py is a script module (a module-level demo script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 337 — BP-PY-46

- Function context: `./scripts/findings/functions/337.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/http2_example.py:30:1`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
print()
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/http2_example.py a module-level demo script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/http2_example.py is a script module (a module-level demo script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 338 — BP-PY-46

- Function context: `./scripts/findings/functions/338.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/http2_example.py:34:1`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
print("Example 3: Per-request override")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/http2_example.py a module-level demo script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/http2_example.py is a script module (a module-level demo script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 339 — BP-PY-46

- Function context: `./scripts/findings/functions/339.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/http2_example.py:35:1`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
print("-" * 40)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/http2_example.py a module-level demo script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/http2_example.py is a script module (a module-level demo script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 340 — BP-PY-46

- Function context: `./scripts/findings/functions/340.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/http2_example.py:42:1`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
print(f"HTTP Version: {response.http_version}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/http2_example.py a module-level demo script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/http2_example.py is a script module (a module-level demo script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 341 — BP-PY-46

- Function context: `./scripts/findings/functions/341.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/http2_example.py:43:1`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
print()
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/http2_example.py a module-level demo script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/http2_example.py is a script module (a module-level demo script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 342 — BP-PY-46

- Function context: `./scripts/findings/functions/342.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/http2_example.py:47:1`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
print("Example 4: httpx API compatibility")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/http2_example.py a module-level demo script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/http2_example.py is a script module (a module-level demo script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 343 — BP-PY-46

- Function context: `./scripts/findings/functions/343.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/http2_example.py:48:1`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
print("-" * 40)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/http2_example.py a module-level demo script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/http2_example.py is a script module (a module-level demo script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 344 — BP-PY-46

- Function context: `./scripts/findings/functions/344.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/http2_example.py:49:1`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
print("""
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/http2_example.py a module-level demo script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/http2_example.py is a script module (a module-level demo script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 345 — BP-PY-46

- Function context: `./scripts/findings/functions/345.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:31:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("ERROR: python-dotenv not installed")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 346 — BP-PY-46

- Function context: `./scripts/findings/functions/346.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:32:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("Install with: pip install python-dotenv")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 348 — BP-PY-46

- Function context: `./scripts/findings/functions/348.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:57:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"Loading proxy config from: {env_file}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 349 — BP-PY-46

- Function context: `./scripts/findings/functions/349.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:60:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print("No .env file found. Create one with PROXY_URL, PROXY_USERNAME, PROXY_PASSWORD")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 350 — BP-PY-46

- Function context: `./scripts/findings/functions/350.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:61:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"Searched in: {[str(p) for p in env_paths]}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 352 — BP-PY-46

- Function context: `./scripts/findings/functions/352.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:72:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\n" + "=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 353 — BP-PY-46

- Function context: `./scripts/findings/functions/353.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:73:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("TEST 1: Direct Connection (no proxy)")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 354 — BP-PY-46

- Function context: `./scripts/findings/functions/354.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:74:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 356 — BP-PY-46

- Function context: `./scripts/findings/functions/356.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:81:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"✓ Status: {response.status_code}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 357 — BP-PY-46

- Function context: `./scripts/findings/functions/357.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:82:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"✓ Response size: {len(response.body)} bytes")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 358 — BP-PY-46

- Function context: `./scripts/findings/functions/358.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:83:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print("✓ Direct connection working")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 359 — BP-PY-46

- Function context: `./scripts/findings/functions/359.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:86:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print("✗ Connection failed")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 360 — BP-PY-46

- Function context: `./scripts/findings/functions/360.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:88:17`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
                print(f"  Error: {response.error_message}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 361 — BP-PY-46

- Function context: `./scripts/findings/functions/361.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:91:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"✗ Failed with status: {response.status_code}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 364 — BP-PY-46

- Function context: `./scripts/findings/functions/364.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:94:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"✗ Error: {e}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 366 — BP-PY-46

- Function context: `./scripts/findings/functions/366.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:100:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\n" + "=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 367 — BP-PY-46

- Function context: `./scripts/findings/functions/367.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:101:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("TEST 2: HTTP Request via Proxy")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 368 — BP-PY-46

- Function context: `./scripts/findings/functions/368.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:102:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 369 — BP-PY-46

- Function context: `./scripts/findings/functions/369.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:103:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print(f"Proxy: {proxy_url}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 370 — BP-PY-46

- Function context: `./scripts/findings/functions/370.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:105:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"Auth: {username}:{'*' * len(password) if password else ''}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 372 — BP-PY-46

- Function context: `./scripts/findings/functions/372.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:115:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"✓ Status: {response.status_code}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 373 — BP-PY-46

- Function context: `./scripts/findings/functions/373.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:117:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"✓ Your IP via proxy: {ip}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 374 — BP-PY-46

- Function context: `./scripts/findings/functions/374.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:120:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print("✗ Connection failed")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 375 — BP-PY-46

- Function context: `./scripts/findings/functions/375.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:122:17`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
                print(f"  Error: {response.error_message}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 376 — BP-PY-46

- Function context: `./scripts/findings/functions/376.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:125:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"✗ HTTP {response.status_code}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 378 — BP-PY-46

- Function context: `./scripts/findings/functions/378.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:128:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"✗ Exception: {e}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 380 — BP-PY-46

- Function context: `./scripts/findings/functions/380.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:134:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\n" + "=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 381 — BP-PY-46

- Function context: `./scripts/findings/functions/381.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:135:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("TEST 3: HTTPS Request via Proxy (CONNECT tunnel)")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 382 — BP-PY-46

- Function context: `./scripts/findings/functions/382.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:136:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 383 — BP-PY-46

- Function context: `./scripts/findings/functions/383.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:137:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print(f"Proxy: {proxy_url}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 384 — BP-PY-46

- Function context: `./scripts/findings/functions/384.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:139:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"Auth: {username}:{'*' * len(password) if password else ''}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 385 — BP-PY-46

- Function context: `./scripts/findings/functions/385.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:149:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"✓ Status: {response.status_code}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 386 — BP-PY-46

- Function context: `./scripts/findings/functions/386.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:151:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"✓ Your IP via proxy: {ip}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 387 — BP-PY-46

- Function context: `./scripts/findings/functions/387.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:152:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"✓ TLS Version: {response.tls_version}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 388 — BP-PY-46

- Function context: `./scripts/findings/functions/388.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:153:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"✓ TLS Cipher: {response.tls_cipher}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 389 — BP-PY-46

- Function context: `./scripts/findings/functions/389.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:156:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print("✗ Connection failed")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 390 — BP-PY-46

- Function context: `./scripts/findings/functions/390.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:158:17`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
                print(f"  Error: {response.error_message}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 391 — BP-PY-46

- Function context: `./scripts/findings/functions/391.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:161:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"✗ HTTP {response.status_code}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 393 — BP-PY-46

- Function context: `./scripts/findings/functions/393.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:164:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"✗ Exception: {e}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 395 — BP-PY-46

- Function context: `./scripts/findings/functions/395.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:170:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\n" + "=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 396 — BP-PY-46

- Function context: `./scripts/findings/functions/396.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:171:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("TEST 4: Proxy Dict Format (requests-compatible)")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 397 — BP-PY-46

- Function context: `./scripts/findings/functions/397.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:172:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 398 — BP-PY-46

- Function context: `./scripts/findings/functions/398.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:175:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print(f"Proxies: {proxies}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 399 — BP-PY-46

- Function context: `./scripts/findings/functions/399.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:185:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"✓ Status: {response.status_code}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 400 — BP-PY-46

- Function context: `./scripts/findings/functions/400.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:187:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"✓ Your IP via proxy: {ip}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 401 — BP-PY-46

- Function context: `./scripts/findings/functions/401.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:190:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print("✗ Connection failed")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 402 — BP-PY-46

- Function context: `./scripts/findings/functions/402.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:193:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"✗ HTTP {response.status_code}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 404 — BP-PY-46

- Function context: `./scripts/findings/functions/404.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:196:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"✗ Exception: {e}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 406 — BP-PY-46

- Function context: `./scripts/findings/functions/406.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:202:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\n" + "=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 407 — BP-PY-46

- Function context: `./scripts/findings/functions/407.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:203:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("TEST 5: Multiple Requests via Proxy")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 408 — BP-PY-46

- Function context: `./scripts/findings/functions/408.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:204:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 409 — BP-PY-46

- Function context: `./scripts/findings/functions/409.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:217:17`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
                print(f"✓ {url}: {response.status_code}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 410 — BP-PY-46

- Function context: `./scripts/findings/functions/410.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:220:17`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
                print(f"✗ {url}: {response.status_code}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 412 — BP-PY-46

- Function context: `./scripts/findings/functions/412.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:222:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"✗ {url}: {e}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 413 — BP-PY-46

- Function context: `./scripts/findings/functions/413.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:224:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print(f"\nResults: {success_count}/{len(urls)} successful")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 415 — BP-PY-46

- Function context: `./scripts/findings/functions/415.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:230:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\n" + "=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 416 — BP-PY-46

- Function context: `./scripts/findings/functions/416.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:231:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("TEST 6: IP Comparison (Direct vs Proxy)")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 417 — BP-PY-46

- Function context: `./scripts/findings/functions/417.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:232:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 418 — BP-PY-46

- Function context: `./scripts/findings/functions/418.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:240:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print("\n[HTTPS Test]")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 419 — BP-PY-46

- Function context: `./scripts/findings/functions/419.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:241:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print("Getting IP without proxy (HTTPS)...")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 420 — BP-PY-46

- Function context: `./scripts/findings/functions/420.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:245:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print("✗ Failed to get direct IP (HTTPS)")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 421 — BP-PY-46

- Function context: `./scripts/findings/functions/421.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:249:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"  Direct IP: {direct_ip_https}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 422 — BP-PY-46

- Function context: `./scripts/findings/functions/422.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:251:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print("Getting IP via proxy (HTTPS)...")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 423 — BP-PY-46

- Function context: `./scripts/findings/functions/423.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:255:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print("✗ Failed to get proxy IP (HTTPS)")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 424 — BP-PY-46

- Function context: `./scripts/findings/functions/424.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:259:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"  Proxy IP: {proxy_ip_https}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 425 — BP-PY-46

- Function context: `./scripts/findings/functions/425.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:263:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print("✓ HTTPS: IPs are different (proxy working)")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 426 — BP-PY-46

- Function context: `./scripts/findings/functions/426.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:265:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print("⚠️  HTTPS: IPs are the same")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 427 — BP-PY-46

- Function context: `./scripts/findings/functions/427.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:268:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print("\n[HTTP Test]")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 428 — BP-PY-46

- Function context: `./scripts/findings/functions/428.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:269:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print("Getting IP without proxy (HTTP)...")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 429 — BP-PY-46

- Function context: `./scripts/findings/functions/429.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:273:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print("✗ Failed to get direct IP (HTTP)")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 430 — BP-PY-46

- Function context: `./scripts/findings/functions/430.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:277:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"  Direct IP: {direct_ip_http}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 431 — BP-PY-46

- Function context: `./scripts/findings/functions/431.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:279:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print("Getting IP via proxy (HTTP)...")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 432 — BP-PY-46

- Function context: `./scripts/findings/functions/432.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:283:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print("✗ Failed to get proxy IP (HTTP)")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 433 — BP-PY-46

- Function context: `./scripts/findings/functions/433.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:287:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"  Proxy IP: {proxy_ip_http}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 434 — BP-PY-46

- Function context: `./scripts/findings/functions/434.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:291:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print("✓ HTTP: IPs are different (proxy working)")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 435 — BP-PY-46

- Function context: `./scripts/findings/functions/435.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:293:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print("⚠️  HTTP: IPs are the same")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 436 — BP-PY-46

- Function context: `./scripts/findings/functions/436.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:297:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print("\n✓ Both HTTP and HTTPS proxy working correctly")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 437 — BP-PY-46

- Function context: `./scripts/findings/functions/437.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:298:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"  Direct: {direct_ip_https}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 438 — BP-PY-46

- Function context: `./scripts/findings/functions/438.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:299:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"  Via Proxy: {proxy_ip_https}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 439 — BP-PY-46

- Function context: `./scripts/findings/functions/439.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:302:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print("\n⚠️  Warning: Only one protocol showing different IP")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 440 — BP-PY-46

- Function context: `./scripts/findings/functions/440.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:305:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print("\n⚠️  Warning: IPs are the same for both protocols (proxy may not be working)")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 442 — BP-PY-46

- Function context: `./scripts/findings/functions/442.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:309:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"✗ Exception: {e}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 444 — BP-PY-46

- Function context: `./scripts/findings/functions/444.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:315:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\n" + "=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 445 — BP-PY-46

- Function context: `./scripts/findings/functions/445.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:316:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("TEST 7: POST Request via Proxy")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 446 — BP-PY-46

- Function context: `./scripts/findings/functions/446.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:317:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("=" * 60)
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 447 — BP-PY-46

- Function context: `./scripts/findings/functions/447.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:329:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"✓ Status: {response.status_code}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 448 — BP-PY-46

- Function context: `./scripts/findings/functions/448.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:330:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print("✓ Posted data via proxy successfully")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 449 — BP-PY-46

- Function context: `./scripts/findings/functions/449.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:337:21`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
                    print(f"✓ Server received: {data['json']}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 454 — BP-PY-46

- Function context: `./scripts/findings/functions/454.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:342:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"✗ HTTP {response.status_code}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 456 — BP-PY-46

- Function context: `./scripts/findings/functions/456.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/examples/proxy_example.py:345:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"✗ Exception: {e}")
```

Why this is a false positive: The print is user-facing output of a demo script module (examples/proxy_example.py guarded by `if __name__ == "__main__":` at line 418); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: examples/proxy_example.py is a script module (guarded by `if __name__ == "__main__":` at line 418); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 458 — BP-PY-46

- Function context: `./scripts/findings/functions/458.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/darwin/setup.py:31:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print("\n=== BoringSSL build directory contents ===")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/darwin/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/darwin/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 459 — BP-PY-46

- Function context: `./scripts/findings/functions/459.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/darwin/setup.py:35:17`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
                print(f"  DIR:  {item}/")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/darwin/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/darwin/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 461 — BP-PY-46

- Function context: `./scripts/findings/functions/461.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/darwin/setup.py:40:29`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
                            print(f"    LIB: {subitem}")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/darwin/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/darwin/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 465 — BP-PY-46

- Function context: `./scripts/findings/functions/465.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/darwin/setup.py:44:17`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
                print(f"  LIB:  {item}")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/darwin/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/darwin/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 466 — BP-PY-46

- Function context: `./scripts/findings/functions/466.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/darwin/setup.py:45:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print("=" * 43 + "\n")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/darwin/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/darwin/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 467 — BP-PY-46

- Function context: `./scripts/findings/functions/467.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/darwin/setup.py:50:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"Using BoringSSL from: {boringssl_lib}")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/darwin/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/darwin/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 468 — BP-PY-46

- Function context: `./scripts/findings/functions/468.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/darwin/setup.py:53:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"Using BoringSSL from: {boringssl_lib}")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/darwin/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/darwin/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 469 — BP-PY-46

- Function context: `./scripts/findings/functions/469.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/darwin/setup.py:57:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"WARNING: BoringSSL libraries not found, using default: {boringssl_lib}")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/darwin/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/darwin/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 470 — BP-PY-46

- Function context: `./scripts/findings/functions/470.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/darwin/setup.py:62:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"Using vendor nghttp2 from: {vendor_nghttp2}")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/darwin/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/darwin/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 471 — BP-PY-46

- Function context: `./scripts/findings/functions/471.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/linux/setup.py:38:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print(f"Using vendor BoringSSL from: {vendor_boringssl}")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/linux/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/linux/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 472 — BP-PY-46

- Function context: `./scripts/findings/functions/472.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/linux/setup.py:45:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"Using vendor nghttp2 from: {vendor_nghttp2}")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/linux/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/linux/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 475 — BP-PY-46

- Function context: `./scripts/findings/functions/475.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/linux/setup.py:93:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"Using vendor liburing from: {vendor_liburing}")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/linux/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/linux/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 476 — BP-PY-46

- Function context: `./scripts/findings/functions/476.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/linux/setup.py:118:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"io_uring support: enabled (kernel {platform.release()})")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/linux/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/linux/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 477 — BP-PY-46

- Function context: `./scripts/findings/functions/477.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/linux/setup.py:120:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"io_uring support: disabled (kernel {platform.release()})")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/linux/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/linux/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 478 — BP-PY-46

- Function context: `./scripts/findings/functions/478.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/linux/setup.py:153:17`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
                print(f"Static linking liburing from: {liburing_static}")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/linux/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/linux/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 479 — BP-PY-46

- Function context: `./scripts/findings/functions/479.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/linux/setup.py:155:17`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
                print("WARNING: liburing.a not found, falling back to dynamic linking")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/linux/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/linux/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 480 — BP-PY-46

- Function context: `./scripts/findings/functions/480.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/linux/setup.py:158:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print("Static linking BoringSSL libraries with --whole-archive")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/linux/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/linux/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 481 — BP-PY-46

- Function context: `./scripts/findings/functions/481.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/linux/setup.py:161:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"WARNING: BoringSSL static libraries not found in {vendor_boringssl_build}")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/linux/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/linux/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 500 — BP-PY-46

- Function context: `./scripts/findings/functions/500.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/windows/setup.py:31:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"Using BoringSSL from: {vendor_boringssl / 'build'} (Visual Studio layout)")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/windows/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/windows/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 501 — BP-PY-46

- Function context: `./scripts/findings/functions/501.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/windows/setup.py:35:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"Using BoringSSL from: {boringssl_lib}")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/windows/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/windows/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 502 — BP-PY-46

- Function context: `./scripts/findings/functions/502.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/windows/setup.py:39:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"Using BoringSSL from: {boringssl_lib}")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/windows/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/windows/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 503 — BP-PY-46

- Function context: `./scripts/findings/functions/503.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/windows/setup.py:46:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"WARNING: BoringSSL not found, using default: {vendor_boringssl / 'build'}")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/windows/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/windows/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 504 — BP-PY-46

- Function context: `./scripts/findings/functions/504.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/windows/setup.py:47:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print("Please run: bash scripts/windows/setup_vendors.sh")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/windows/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/windows/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 505 — BP-PY-46

- Function context: `./scripts/findings/functions/505.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/windows/setup.py:55:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"Using vendor nghttp2 from: {vendor_nghttp2}")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/windows/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/windows/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 506 — BP-PY-46

- Function context: `./scripts/findings/functions/506.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/windows/setup.py:60:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"Using vendor nghttp2 from: {vendor_nghttp2}")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/windows/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/windows/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 507 — BP-PY-46

- Function context: `./scripts/findings/functions/507.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/windows/setup.py:65:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"Using vendor nghttp2 from: {vendor_nghttp2}")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/windows/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/windows/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 508 — BP-PY-46

- Function context: `./scripts/findings/functions/508.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/windows/setup.py:74:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"Using vcpkg nghttp2 from: {vcpkg_installed}")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/windows/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/windows/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 509 — BP-PY-46

- Function context: `./scripts/findings/functions/509.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/windows/setup.py:78:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print("Using MSYS2 nghttp2 from: /mingw64")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/windows/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/windows/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 510 — BP-PY-46

- Function context: `./scripts/findings/functions/510.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/windows/setup.py:83:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print("WARNING: nghttp2 not found. Run: bash scripts/windows/setup_vendors.sh")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/windows/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/windows/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 511 — BP-PY-46

- Function context: `./scripts/findings/functions/511.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/windows/setup.py:93:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"Using vendor zlib from: {vendor_zlib}")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/windows/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/windows/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 512 — BP-PY-46

- Function context: `./scripts/findings/functions/512.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/windows/setup.py:98:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print(f"Using vcpkg zlib from: {vcpkg_installed}")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/windows/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/windows/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 513 — BP-PY-46

- Function context: `./scripts/findings/functions/513.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/scripts/windows/setup.py:102:9`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
        print("WARNING: zlib not found. Install via vcpkg: vcpkg install zlib:x64-windows")
```

Why this is a false positive: The print is build-status output of a build helper script (scripts/windows/setup.py a build helper script with no library exports); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: scripts/windows/setup.py is a script module (a build helper script with no library exports); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 517 — BP-PY-46

- Function context: `./scripts/findings/functions/517.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:75:1`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
print(f"Building for platform: {platform.system()}")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 518 — BP-PY-46

- Function context: `./scripts/findings/functions/518.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:76:1`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
print(f"io_uring support: {HAS_IO_URING}")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 519 — BP-PY-46

- Function context: `./scripts/findings/functions/519.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:156:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\n" + "=" * 70)
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 520 — BP-PY-46

- Function context: `./scripts/findings/functions/520.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:157:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("WARNING: Vendor dependencies not found!")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 521 — BP-PY-46

- Function context: `./scripts/findings/functions/521.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:158:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("=" * 70)
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 522 — BP-PY-46

- Function context: `./scripts/findings/functions/522.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:159:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\nPlease run: make setup")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 523 — BP-PY-46

- Function context: `./scripts/findings/functions/523.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:160:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\nThis will download and build:")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 524 — BP-PY-46

- Function context: `./scripts/findings/functions/524.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:161:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("  • BoringSSL (TLS library)")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 525 — BP-PY-46

- Function context: `./scripts/findings/functions/525.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:162:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("  • liburing (io_uring - Linux only)")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 526 — BP-PY-46

- Function context: `./scripts/findings/functions/526.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:163:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("  • nghttp2 (HTTP/2 library)")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 527 — BP-PY-46

- Function context: `./scripts/findings/functions/527.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:164:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\n" + "=" * 70 + "\n")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 529 — BP-PY-46

- Function context: `./scripts/findings/functions/529.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:185:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"Using vendor BoringSSL from: {vendor_boringssl}")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 530 — BP-PY-46

- Function context: `./scripts/findings/functions/530.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:200:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"Using Homebrew OpenSSL from: {openssl_prefix}")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 531 — BP-PY-46

- Function context: `./scripts/findings/functions/531.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:207:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"Using vendor nghttp2 from: {vendor_nghttp2}")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 532 — BP-PY-46

- Function context: `./scripts/findings/functions/532.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:222:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"Using Homebrew nghttp2 from: {nghttp2_prefix}")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 533 — BP-PY-46

- Function context: `./scripts/findings/functions/533.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:240:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"Using vendor BoringSSL from: {vendor_boringssl}")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 534 — BP-PY-46

- Function context: `./scripts/findings/functions/534.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:250:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"Using vendor nghttp2 from: {vendor_nghttp2}")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 538 — BP-PY-46

- Function context: `./scripts/findings/functions/538.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:347:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"Using vendor liburing from: {vendor_liburing_install}")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 539 — BP-PY-46

- Function context: `./scripts/findings/functions/539.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:352:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"Using vendor liburing from: {vendor_dir / 'liburing'}")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 540 — BP-PY-46

- Function context: `./scripts/findings/functions/540.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:377:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"WARNING: BoringSSL not found at {vendor_boringssl}")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 541 — BP-PY-46

- Function context: `./scripts/findings/functions/541.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:378:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print("Please run: make setup")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 542 — BP-PY-46

- Function context: `./scripts/findings/functions/542.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:385:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"Using vendor nghttp2 from: {vendor_nghttp2}")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 543 — BP-PY-46

- Function context: `./scripts/findings/functions/543.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:390:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"Using vendor nghttp2 from: {vendor_nghttp2}")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 544 — BP-PY-46

- Function context: `./scripts/findings/functions/544.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:399:17`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
                print(f"Using vcpkg nghttp2 from: {vcpkg_installed}")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 545 — BP-PY-46

- Function context: `./scripts/findings/functions/545.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:403:17`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
                print("Using MSYS2 nghttp2 from: /mingw64")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 546 — BP-PY-46

- Function context: `./scripts/findings/functions/546.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:408:17`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
                print("WARNING: nghttp2 not found. Install via vcpkg or MSYS2")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 547 — BP-PY-46

- Function context: `./scripts/findings/functions/547.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:418:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"Using vendor zlib from: {vendor_zlib}")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 548 — BP-PY-46

- Function context: `./scripts/findings/functions/548.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:423:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"Using vcpkg zlib from: {vcpkg_installed}")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 549 — BP-PY-46

- Function context: `./scripts/findings/functions/549.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:427:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print("WARNING: zlib not found. Install via vcpkg: vcpkg install zlib:x64-windows")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 550 — BP-PY-46

- Function context: `./scripts/findings/functions/550.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:433:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"Using vcpkg brotli from: {vcpkg_installed}")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 551 — BP-PY-46

- Function context: `./scripts/findings/functions/551.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:437:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print("WARNING: brotli not found. Install via vcpkg: vcpkg install brotli:x64-windows")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 552 — BP-PY-46

- Function context: `./scripts/findings/functions/552.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:465:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("\nLibrary paths detected:")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 553 — BP-PY-46

- Function context: `./scripts/findings/functions/553.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:466:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print(f"  BoringSSL include: {LIB_PATHS['openssl_include']}")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 554 — BP-PY-46

- Function context: `./scripts/findings/functions/554.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:467:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print(f"  BoringSSL lib: {LIB_PATHS['openssl_lib']}")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 555 — BP-PY-46

- Function context: `./scripts/findings/functions/555.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:468:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print(f"  nghttp2 include: {LIB_PATHS['nghttp2_include']}")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 556 — BP-PY-46

- Function context: `./scripts/findings/functions/556.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:469:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print(f"  nghttp2 lib: {LIB_PATHS['nghttp2_lib']}")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 557 — BP-PY-46

- Function context: `./scripts/findings/functions/557.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:470:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print()
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 558 — BP-PY-46

- Function context: `./scripts/findings/functions/558.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:569:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print(f"Using vendor brotli from: {vendor_dir / 'brotli'}")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 559 — BP-PY-46

- Function context: `./scripts/findings/functions/559.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:579:17`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
                print(f"Using Homebrew brotli from: {homebrew_prefix}")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 560 — BP-PY-46

- Function context: `./scripts/findings/functions/560.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:648:13`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
            print("\nUsing static libraries:")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 561 — BP-PY-46

- Function context: `./scripts/findings/functions/561.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:650:17`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
                print(f"  {lib}")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 562 — BP-PY-46

- Function context: `./scripts/findings/functions/562.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/setup.py:755:5`
- Checklist pattern: print statement is user-facing/CLI or build-status output in a script module, not operational logging in a non-script (library) module

Source excerpt:

```
    print("Skipping C extension build on Read the Docs")
```

Why this is a false positive: The print is build-status output of the setuptools packaging script (setup.py guarded by `if __name__ == "__main__":` at line 771); BP-PY-46's condition limits it to operational logging in non-script modules, and this is a script module, not a library module.

Checklist evidence: setup.py is a script module (guarded by `if __name__ == "__main__":` at line 771); the flagged line is a `print(` call executed as script output, not library debug logging.

### [ ] Finding 563 — BP-PY-46

- Function context: `./scripts/findings/functions/563.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/src/httpmorph/_async_client.py:145:13`
- Checklist pattern: print call appears inside a docstring, not in executable code

Source excerpt:

```
            print(response.status_code)
    """
```

Why this is a false positive: The flagged `print` is the example body inside the method's docstring (opened at line 130, closed at line 146); it is never executed.

Checklist evidence: Line 145 sits between the docstring opener (line 130) and closer (line 146); `printCallOutsideString` inspects only the single line and cannot see the enclosing multi-line string, so the print inside the docstring is flagged.

### [ ] Finding 570 — CWE-829

- Function context: `./scripts/findings/functions/570.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/src/httpmorph/_client_c.py:61:12`
- Checklist pattern: module path is selected from a package-controlled glob of the package's own directory

Source excerpt:

```
    current_dir = os.path.dirname(os.path.abspath(__file__))
    so_files = glob.glob(os.path.join(current_dir, "_httpmorph*.so"))
    spec = importlib.util.spec_from_file_location("_httpmorph", so_files[0])
```

Why this is a false positive: The path is derived from `__file__`'s own directory with a fixed `_httpmorph*.so` pattern — package-controlled, so no untrusted control sphere can select what executes; CWE-829's own scope suppresses package-controlled paths.

Checklist evidence: `so_files` is the result of a literal glob in the package's own directory (plus a `.pyd` fallback of the same directory); only `so_files[0]` (a subscript of a package-owned list) is dynamic, which `isDynamicExpr` treats as untrusted.

### [ ] Finding 594 — CWE-256

- Function context: `./scripts/findings/functions/594.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/tests/test_edge_cases_security.py:208:1`
- Checklist pattern: assignment is synthetic test input, not a stored password literal

Source excerpt:

```
    # Long credentials should be truncated, not cause buffer overflow
    long_username = "a" * 300
    long_password = "b" * 300
```

Why this is a false positive: `long_password` is a 300-char generated test buffer for overflow/truncation testing, not a stored password; `pyPlaintextPasswordRE` matches only the `long_password = "b"` prefix and cannot tell it is a multiplied fixture.

Checklist evidence: The value is the expression `"b" * 300`, not a plaintext literal, and the surrounding test (`test_edge_cases_security.py`) verifies credential truncation behavior; no credential is stored.

### [ ] Finding 680 — BP-PY-42

- Function context: `./scripts/findings/functions/680.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/tests/test_proxy_server.py:31:1`
- Checklist pattern: try/except is server-implementation error handling, not a test asserting failure

Source excerpt:

```
            try:
                auth_type, credentials = auth_header.split(" ", 1)
                if auth_type == "Basic":
                    decoded = base64.b64decode(credentials).decode("utf-8")
                    ...
            except Exception:
                self.send_response(403)
                self.end_headers()
                return
```

Why this is a false positive: The try/except is the CONNECT proxy's credential-parsing error handling (it replies 403), not a test using try/except to assert an expected failure — BP-PY-42's condition is "Tests use a bare try/except instead of assertRaises/pytest.raises".

Checklist evidence: `test_proxy_server.py` contains no `def test_*` functions, so `detectBPPY42` falls back to whole-file scanning and flags handler code in `do_CONNECT`; `isExpectFailureExcept` matches the `except Exception:` clause even though the region is not an assertion.

### [ ] Finding 684 — CWE-1341

- Function context: `./scripts/findings/functions/684.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/tests/test_proxy_server.py:87:29`
- Checklist pattern: each close() terminates its own execution path — the handle is never released twice

Source excerpt:

```
                            if not data:
                                target_sock.close()
                                return
                        except Exception:
                            target_sock.close()
                            return
                except Exception:
                    break

            target_sock.close()
```

Why this is a false positive: Every `target_sock.close()` is immediately followed by `return` (or is the single close after the relay loop ends), so at most one close executes per path; the handle is never released twice.

Checklist evidence: `pyTierBDoubleCloseRE` matches any two `.close()` calls within 180 chars of source regardless of control flow; the shown source proves each close is on a distinct exit path.

### [ ] Finding 697 — BP-PY-42

- Function context: `./scripts/findings/functions/697.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/httpmorph/tests/test_server.py:83:1`
- Checklist pattern: try/except is server-implementation error handling, not a test asserting failure

Source excerpt:

```
        elif path_without_query.startswith("/delay/"):
            try:
                delay = int(path_without_query.split("/")[-1])
                ...
            except Exception:
                self.send_response(400)
```

Why this is a false positive: The try/except is the mock HTTP server's `/delay/` route parsing (it replies 400 on malformed input), not a test asserting failure — BP-PY-42's condition is "Tests use a bare try/except instead of assertRaises/pytest.raises".

Checklist evidence: `test_server.py` defines no `def test_*` functions, so `detectBPPY42` whole-file scans the test-named module and flags the `do_GET` handler implementation, which is not an assertion substitute.

## True positives

### BP-PY-45 — [1, 224, 347]

| Finding id | Source | Reason |
| --- | --- | --- |
| 1 | `benchmarks/benchmark.py:47:1` | `sys.path.insert` executed at runtime — rule condition met |
| 224 | `docs/source/conf.py:9:1` | `sys.path.insert` executed at runtime — rule condition met |
| 347 | `examples/proxy_example.py:36:1` | `sys.path.insert` executed at runtime — rule condition met |

### BP-PY-2 — [2, 12, 21, 23, 25, 31, 33, 54, 62, 70, 77, 81, 89, 93, 97, 108, 112, 116, 138, 145, 149, 152, 158, 165, 172, 179, 183, 193, 206, 213, 218, 222, 451, 462, 473, 514, 516, 535, 537, 566, 571, 575, 577, 578, 583, 586, 589, 596, 609, 613, 615, 617, 619, 621, 623, 625, 627, 629, 631, 633, 636, 639, 642, 645, 646, 649, 652, 654, 656, 658, 660, 662, 664, 670, 675, 678, 711, 713, 714]

| Finding id | Source | Reason |
| --- | --- | --- |
| 2 | `benchmarks/benchmark.py:56:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 12 | `benchmarks/benchmark.py:493:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 21 | `benchmarks/benchmark.py:1843:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 23 | `benchmarks/benchmark.py:1850:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 25 | `benchmarks/libs/aiohttp_bench.py:81:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 31 | `benchmarks/libs/aiohttp_bench.py:106:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 33 | `benchmarks/libs/base.py:69:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 54 | `benchmarks/libs/curl_cffi_bench.py:177:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 62 | `benchmarks/libs/curl_cffi_bench.py:213:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 70 | `benchmarks/libs/httpmorph_bench.py:103:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 77 | `benchmarks/libs/httpmorph_bench.py:130:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 81 | `benchmarks/libs/httpmorph_bench.py:159:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 89 | `benchmarks/libs/httpmorph_bench.py:229:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 93 | `benchmarks/libs/httpmorph_bench.py:251:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 97 | `benchmarks/libs/httpmorph_bench.py:276:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 108 | `benchmarks/libs/httpmorph_bench.py:384:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 112 | `benchmarks/libs/httpmorph_bench.py:415:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 116 | `benchmarks/libs/httpmorph_bench.py:448:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 138 | `benchmarks/libs/httpx_bench.py:236:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 145 | `benchmarks/libs/httpx_bench.py:256:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 149 | `benchmarks/libs/httpx_bench.py:276:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 152 | `benchmarks/libs/pycurl_bench.py:277:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 158 | `benchmarks/libs/pycurl_bench.py:332:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 165 | `benchmarks/libs/requests_bench.py:74:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 172 | `benchmarks/libs/requests_bench.py:99:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 179 | `benchmarks/libs/requests_bench.py:154:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 183 | `benchmarks/libs/requests_bench.py:173:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 193 | `benchmarks/libs/urllib3_bench.py:125:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 206 | `benchmarks/libs/urllib_bench.py:76:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 213 | `benchmarks/libs/urllib_bench.py:99:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 218 | `benchmarks/libs/urllib_bench.py:152:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 222 | `benchmarks/libs/urllib_bench.py:175:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 451 | `examples/proxy_example.py:338:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 462 | `scripts/darwin/setup.py:41:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 473 | `scripts/linux/setup.py:70:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 514 | `setup.py:68:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 516 | `setup.py:72:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 535 | `setup.py:277:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 537 | `setup.py:316:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 566 | `src/httpmorph/_client_c.py:39:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 571 | `src/httpmorph/_client_c.py:379:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 575 | `src/httpmorph/_client_c.py:1035:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 577 | `src/httpmorph/_client_c.py:1042:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 578 | `tests/conftest.py:23:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 583 | `tests/test_buffer_reallocation.py:20:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 586 | `tests/test_connection_pool.py:18:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 589 | `tests/test_edge_cases_security.py:19:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 596 | `tests/test_edge_cases_security.py:222:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 609 | `tests/test_proxy.py:40:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 613 | `tests/test_proxy.py:50:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 615 | `tests/test_proxy.py:67:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 617 | `tests/test_proxy.py:106:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 619 | `tests/test_proxy.py:117:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 621 | `tests/test_proxy.py:157:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 623 | `tests/test_proxy.py:174:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 625 | `tests/test_proxy.py:184:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 627 | `tests/test_proxy.py:200:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 629 | `tests/test_proxy.py:214:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 631 | `tests/test_proxy.py:227:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 633 | `tests/test_proxy.py:260:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 636 | `tests/test_proxy.py:304:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 639 | `tests/test_proxy.py:378:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 642 | `tests/test_proxy.py:434:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 645 | `tests/test_proxy.py:462:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 646 | `tests/test_proxy.py:512:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 649 | `tests/test_proxy.py:561:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 652 | `tests/test_proxy.py:667:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 654 | `tests/test_proxy.py:682:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 656 | `tests/test_proxy.py:704:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 658 | `tests/test_proxy.py:759:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 660 | `tests/test_proxy.py:775:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 662 | `tests/test_proxy.py:851:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 664 | `tests/test_proxy.py:868:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 670 | `tests/test_proxy.py:980:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 675 | `tests/test_proxy.py:1145:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 678 | `tests/test_proxy.py:1190:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 711 | `tests/test_unicode.py:154:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 713 | `tests/test_unicode.py:294:1` | except suite is solely `pass` (comments stripped) — rule condition met |
| 714 | `tests/test_unicode.py:308:1` | except suite is solely `pass` (comments stripped) — rule condition met |

### CWE-1071 — [3, 28, 57, 73, 141, 155, 168, 196, 209, 453, 464, 569, 579, 584, 587, 590, 611]

| Finding id | Source | Reason |
| --- | --- | --- |
| 3 | `benchmarks/benchmark.py:56:1` | exception handler body is empty (`pass` only) — rule condition met |
| 28 | `benchmarks/libs/aiohttp_bench.py:81:13` | exception handler body is empty (`pass` only) — rule condition met |
| 57 | `benchmarks/libs/curl_cffi_bench.py:177:13` | exception handler body is empty (`pass` only) — rule condition met |
| 73 | `benchmarks/libs/httpmorph_bench.py:103:13` | exception handler body is empty (`pass` only) — rule condition met |
| 141 | `benchmarks/libs/httpx_bench.py:236:13` | exception handler body is empty (`pass` only) — rule condition met |
| 155 | `benchmarks/libs/pycurl_bench.py:277:13` | exception handler body is empty (`pass` only) — rule condition met |
| 168 | `benchmarks/libs/requests_bench.py:74:13` | exception handler body is empty (`pass` only) — rule condition met |
| 196 | `benchmarks/libs/urllib3_bench.py:125:13` | exception handler body is empty (`pass` only) — rule condition met |
| 209 | `benchmarks/libs/urllib_bench.py:76:13` | exception handler body is empty (`pass` only) — rule condition met |
| 453 | `examples/proxy_example.py:338:13` | exception handler body is empty (`pass` only) — rule condition met |
| 464 | `scripts/darwin/setup.py:41:17` | exception handler body is empty (`pass` only) — rule condition met |
| 569 | `src/httpmorph/_client_c.py:39:5` | exception handler body is empty (`pass` only) — rule condition met |
| 579 | `tests/conftest.py:23:1` | exception handler body is empty (`pass` only) — rule condition met |
| 584 | `tests/test_buffer_reallocation.py:20:1` | exception handler body is empty (`pass` only) — rule condition met |
| 587 | `tests/test_connection_pool.py:18:1` | exception handler body is empty (`pass` only) — rule condition met |
| 590 | `tests/test_edge_cases_security.py:19:1` | exception handler body is empty (`pass` only) — rule condition met |
| 611 | `tests/test_proxy.py:40:9` | exception handler body is empty (`pass` only) — rule condition met |

### CWE-390 — [4, 26, 34, 55, 71, 139, 153, 166, 194, 207, 452, 463, 474, 515, 567, 580, 585, 588, 591, 610, 712]

| Finding id | Source | Reason |
| --- | --- | --- |
| 4 | `benchmarks/benchmark.py:56:1` | except handler takes no action (only `pass`) — rule condition met |
| 26 | `benchmarks/libs/aiohttp_bench.py:81:1` | except handler takes no action (only `pass`) — rule condition met |
| 34 | `benchmarks/libs/base.py:69:1` | except handler takes no action (only `pass`) — rule condition met |
| 55 | `benchmarks/libs/curl_cffi_bench.py:177:1` | except handler takes no action (only `pass`) — rule condition met |
| 71 | `benchmarks/libs/httpmorph_bench.py:103:1` | except handler takes no action (only `pass`) — rule condition met |
| 139 | `benchmarks/libs/httpx_bench.py:236:1` | except handler takes no action (only `pass`) — rule condition met |
| 153 | `benchmarks/libs/pycurl_bench.py:277:1` | except handler takes no action (only `pass`) — rule condition met |
| 166 | `benchmarks/libs/requests_bench.py:74:1` | except handler takes no action (only `pass`) — rule condition met |
| 194 | `benchmarks/libs/urllib3_bench.py:125:1` | except handler takes no action (only `pass`) — rule condition met |
| 207 | `benchmarks/libs/urllib_bench.py:76:1` | except handler takes no action (only `pass`) — rule condition met |
| 452 | `examples/proxy_example.py:338:1` | except handler takes no action (only `pass`) — rule condition met |
| 463 | `scripts/darwin/setup.py:41:1` | except handler takes no action (only `pass`) — rule condition met |
| 474 | `scripts/linux/setup.py:70:1` | except handler takes no action (only `pass`) — rule condition met |
| 515 | `setup.py:68:1` | except handler takes no action (only `pass`) — rule condition met |
| 567 | `src/httpmorph/_client_c.py:39:1` | except handler takes no action (only `pass`) — rule condition met |
| 580 | `tests/conftest.py:23:1` | except handler takes no action (only `pass`) — rule condition met |
| 585 | `tests/test_buffer_reallocation.py:20:1` | except handler takes no action (only `pass`) — rule condition met |
| 588 | `tests/test_connection_pool.py:18:1` | except handler takes no action (only `pass`) — rule condition met |
| 591 | `tests/test_edge_cases_security.py:19:1` | except handler takes no action (only `pass`) — rule condition met |
| 610 | `tests/test_proxy.py:40:1` | except handler takes no action (only `pass`) — rule condition met |
| 712 | `tests/test_unicode.py:154:1` | except handler takes no action (only `pass`) — rule condition met |

### BP-PY-1 — [5, 7, 9, 11, 13, 15, 16, 17, 18, 20, 22, 24, 29, 30, 32, 36, 38, 39, 40, 41, 42, 53, 58, 61, 63, 69, 74, 76, 78, 80, 82, 88, 90, 92, 94, 96, 98, 99, 101, 103, 105, 107, 109, 111, 113, 115, 117, 137, 142, 144, 146, 148, 150, 151, 156, 157, 159, 164, 169, 171, 173, 178, 180, 182, 184, 192, 197, 205, 210, 212, 214, 217, 219, 221, 223, 271, 273, 362, 377, 392, 403, 411, 441, 450, 455, 485, 487, 489, 491, 492, 494, 496, 497, 498, 565, 574, 576, 593, 595, 597, 602, 608, 612, 614, 616, 618, 620, 622, 624, 626, 628, 630, 632, 635, 638, 641, 644, 648, 650, 651, 653, 655, 657, 659, 661, 663, 665, 666, 667, 668, 669, 671, 672, 674, 677, 681, 683, 685, 686, 687, 698, 699, 700, 701, 702, 703, 704, 705, 707, 708, 709]

| Finding id | Source | Reason |
| --- | --- | --- |
| 5 | `benchmarks/benchmark.py:86:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 7 | `benchmarks/benchmark.py:107:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 9 | `benchmarks/benchmark.py:375:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 11 | `benchmarks/benchmark.py:493:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 13 | `benchmarks/benchmark.py:516:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 15 | `benchmarks/benchmark.py:635:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 16 | `benchmarks/benchmark.py:1023:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 17 | `benchmarks/benchmark.py:1367:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 18 | `benchmarks/benchmark.py:1456:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 20 | `benchmarks/benchmark.py:1843:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 22 | `benchmarks/benchmark.py:1850:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 24 | `benchmarks/libs/aiohttp_bench.py:81:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 29 | `benchmarks/libs/aiohttp_bench.py:86:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 30 | `benchmarks/libs/aiohttp_bench.py:106:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 32 | `benchmarks/libs/aiohttp_bench.py:111:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 36 | `benchmarks/libs/base.py:118:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 38 | `benchmarks/libs/base.py:146:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 39 | `benchmarks/libs/base.py:190:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 40 | `benchmarks/libs/base.py:219:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 41 | `benchmarks/libs/base.py:264:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 42 | `benchmarks/libs/base.py:293:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 53 | `benchmarks/libs/curl_cffi_bench.py:177:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 58 | `benchmarks/libs/curl_cffi_bench.py:182:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 61 | `benchmarks/libs/curl_cffi_bench.py:213:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 63 | `benchmarks/libs/curl_cffi_bench.py:218:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 69 | `benchmarks/libs/httpmorph_bench.py:103:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 74 | `benchmarks/libs/httpmorph_bench.py:109:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 76 | `benchmarks/libs/httpmorph_bench.py:130:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 78 | `benchmarks/libs/httpmorph_bench.py:136:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 80 | `benchmarks/libs/httpmorph_bench.py:159:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 82 | `benchmarks/libs/httpmorph_bench.py:164:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 88 | `benchmarks/libs/httpmorph_bench.py:229:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 90 | `benchmarks/libs/httpmorph_bench.py:235:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 92 | `benchmarks/libs/httpmorph_bench.py:251:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 94 | `benchmarks/libs/httpmorph_bench.py:257:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 96 | `benchmarks/libs/httpmorph_bench.py:276:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 98 | `benchmarks/libs/httpmorph_bench.py:282:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 99 | `benchmarks/libs/httpmorph_bench.py:303:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 101 | `benchmarks/libs/httpmorph_bench.py:322:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 103 | `benchmarks/libs/httpmorph_bench.py:341:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 105 | `benchmarks/libs/httpmorph_bench.py:362:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 107 | `benchmarks/libs/httpmorph_bench.py:384:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 109 | `benchmarks/libs/httpmorph_bench.py:392:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 111 | `benchmarks/libs/httpmorph_bench.py:415:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 113 | `benchmarks/libs/httpmorph_bench.py:423:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 115 | `benchmarks/libs/httpmorph_bench.py:448:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 117 | `benchmarks/libs/httpmorph_bench.py:456:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 137 | `benchmarks/libs/httpx_bench.py:236:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 142 | `benchmarks/libs/httpx_bench.py:241:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 144 | `benchmarks/libs/httpx_bench.py:256:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 146 | `benchmarks/libs/httpx_bench.py:261:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 148 | `benchmarks/libs/httpx_bench.py:276:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 150 | `benchmarks/libs/httpx_bench.py:281:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 151 | `benchmarks/libs/pycurl_bench.py:277:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 156 | `benchmarks/libs/pycurl_bench.py:282:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 157 | `benchmarks/libs/pycurl_bench.py:332:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 159 | `benchmarks/libs/pycurl_bench.py:337:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 164 | `benchmarks/libs/requests_bench.py:74:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 169 | `benchmarks/libs/requests_bench.py:79:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 171 | `benchmarks/libs/requests_bench.py:99:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 173 | `benchmarks/libs/requests_bench.py:104:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 178 | `benchmarks/libs/requests_bench.py:154:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 180 | `benchmarks/libs/requests_bench.py:159:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 182 | `benchmarks/libs/requests_bench.py:173:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 184 | `benchmarks/libs/requests_bench.py:178:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 192 | `benchmarks/libs/urllib3_bench.py:125:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 197 | `benchmarks/libs/urllib3_bench.py:130:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 205 | `benchmarks/libs/urllib_bench.py:76:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 210 | `benchmarks/libs/urllib_bench.py:81:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 212 | `benchmarks/libs/urllib_bench.py:99:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 214 | `benchmarks/libs/urllib_bench.py:104:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 217 | `benchmarks/libs/urllib_bench.py:152:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 219 | `benchmarks/libs/urllib_bench.py:157:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 221 | `benchmarks/libs/urllib_bench.py:175:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 223 | `benchmarks/libs/urllib_bench.py:180:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 271 | `examples/advanced_features.py:143:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 273 | `examples/async_example.py:60:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 362 | `examples/proxy_example.py:93:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 377 | `examples/proxy_example.py:127:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 392 | `examples/proxy_example.py:163:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 403 | `examples/proxy_example.py:195:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 411 | `examples/proxy_example.py:221:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 441 | `examples/proxy_example.py:308:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 450 | `examples/proxy_example.py:338:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 455 | `examples/proxy_example.py:344:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 485 | `scripts/test_local_build.py:50:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 487 | `scripts/test_local_build.py:78:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 489 | `scripts/test_local_build.py:99:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 491 | `scripts/test_local_build.py:122:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 492 | `scripts/test_local_build.py:126:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 494 | `scripts/test_local_build.py:170:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 496 | `scripts/test_local_build.py:198:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 497 | `scripts/test_local_build.py:202:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 498 | `scripts/test_local_build.py:237:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 565 | `src/httpmorph/_client_c.py:39:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 574 | `src/httpmorph/_client_c.py:1035:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 576 | `src/httpmorph/_client_c.py:1042:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 593 | `tests/test_edge_cases_security.py:54:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 595 | `tests/test_edge_cases_security.py:222:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 597 | `tests/test_edge_cases_security.py:355:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 602 | `tests/test_http2.py:49:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 608 | `tests/test_proxy.py:40:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 612 | `tests/test_proxy.py:50:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 614 | `tests/test_proxy.py:67:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 616 | `tests/test_proxy.py:106:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 618 | `tests/test_proxy.py:117:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 620 | `tests/test_proxy.py:157:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 622 | `tests/test_proxy.py:174:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 624 | `tests/test_proxy.py:184:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 626 | `tests/test_proxy.py:200:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 628 | `tests/test_proxy.py:214:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 630 | `tests/test_proxy.py:227:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 632 | `tests/test_proxy.py:260:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 635 | `tests/test_proxy.py:304:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 638 | `tests/test_proxy.py:378:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 641 | `tests/test_proxy.py:434:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 644 | `tests/test_proxy.py:462:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 648 | `tests/test_proxy.py:561:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 650 | `tests/test_proxy.py:622:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 651 | `tests/test_proxy.py:667:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 653 | `tests/test_proxy.py:682:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 655 | `tests/test_proxy.py:704:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 657 | `tests/test_proxy.py:759:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 659 | `tests/test_proxy.py:775:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 661 | `tests/test_proxy.py:851:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 663 | `tests/test_proxy.py:868:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 665 | `tests/test_proxy.py:885:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 666 | `tests/test_proxy.py:909:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 667 | `tests/test_proxy.py:939:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 668 | `tests/test_proxy.py:951:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 669 | `tests/test_proxy.py:980:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 671 | `tests/test_proxy.py:1020:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 672 | `tests/test_proxy.py:1102:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 674 | `tests/test_proxy.py:1145:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 677 | `tests/test_proxy.py:1190:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 681 | `tests/test_proxy_server.py:40:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 683 | `tests/test_proxy_server.py:85:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 685 | `tests/test_proxy_server.py:89:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 686 | `tests/test_proxy_server.py:93:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 687 | `tests/test_proxy_server.py:118:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 698 | `tests/test_server.py:91:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 699 | `tests/test_server.py:140:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 700 | `tests/test_server.py:156:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 701 | `tests/test_server.py:245:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 702 | `tests/test_server.py:261:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 703 | `tests/test_server.py:310:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 704 | `tests/test_server.py:326:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 705 | `tests/test_server.py:340:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 707 | `tests/test_server.py:504:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 708 | `tests/test_server.py:564:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |
| 709 | `tests/test_server.py:617:1` | broad `except Exception`/`except BaseException` whose suite does not re-raise (`suiteReraises` false) — rule condition met |

### CWE-396 — [6, 27, 37, 56, 72, 140, 154, 167, 195, 208, 272, 274, 363, 568]

| Finding id | Source | Reason |
| --- | --- | --- |
| 6 | `benchmarks/benchmark.py:86:1` | generic `except Exception` handler can hide distinct failure conditions — rule condition met |
| 27 | `benchmarks/libs/aiohttp_bench.py:81:1` | generic `except Exception` handler can hide distinct failure conditions — rule condition met |
| 37 | `benchmarks/libs/base.py:118:1` | generic `except Exception` handler can hide distinct failure conditions — rule condition met |
| 56 | `benchmarks/libs/curl_cffi_bench.py:177:1` | generic `except Exception` handler can hide distinct failure conditions — rule condition met |
| 72 | `benchmarks/libs/httpmorph_bench.py:103:1` | generic `except Exception` handler can hide distinct failure conditions — rule condition met |
| 140 | `benchmarks/libs/httpx_bench.py:236:1` | generic `except Exception` handler can hide distinct failure conditions — rule condition met |
| 154 | `benchmarks/libs/pycurl_bench.py:277:1` | generic `except Exception` handler can hide distinct failure conditions — rule condition met |
| 167 | `benchmarks/libs/requests_bench.py:74:1` | generic `except Exception` handler can hide distinct failure conditions — rule condition met |
| 195 | `benchmarks/libs/urllib3_bench.py:125:1` | generic `except Exception` handler can hide distinct failure conditions — rule condition met |
| 208 | `benchmarks/libs/urllib_bench.py:76:1` | generic `except Exception` handler can hide distinct failure conditions — rule condition met |
| 272 | `examples/advanced_features.py:143:1` | generic `except Exception` handler can hide distinct failure conditions — rule condition met |
| 274 | `examples/async_example.py:60:1` | generic `except Exception` handler can hide distinct failure conditions — rule condition met |
| 363 | `examples/proxy_example.py:93:1` | generic `except Exception` handler can hide distinct failure conditions — rule condition met |
| 568 | `src/httpmorph/_client_c.py:39:1` | generic `except Exception` handler can hide distinct failure conditions — rule condition met |

### CWE-1121 — [8, 35, 457, 499, 528, 572, 679, 696]

| Finding id | Source | Reason |
| --- | --- | --- |
| 8 | `benchmarks/benchmark.py:270:64` | function body has ≥12 branch keywords (if/elif/for/while/except) — rule condition met |
| 35 | `benchmarks/libs/base.py:83:52` | function body has ≥12 branch keywords (if/elif/for/while/except) — rule condition met |
| 457 | `scripts/darwin/setup.py:16:25` | function body has ≥12 branch keywords (if/elif/for/while/except) — rule condition met |
| 499 | `scripts/windows/setup.py:15:25` | function body has ≥12 branch keywords (if/elif/for/while/except) — rule condition met |
| 528 | `setup.py:168:25` | function body has ≥12 branch keywords (if/elif/for/while/except) — rule condition met |
| 572 | `src/httpmorph/_client_c.py:395:46` | function body has ≥12 branch keywords (if/elif/for/while/except) — rule condition met |
| 679 | `tests/test_proxy_server.py:19:26` | function body has ≥12 branch keywords (if/elif/for/while/except) — rule condition met |
| 696 | `tests/test_server.py:24:22` | function body has ≥12 branch keywords (if/elif/for/while/except) — rule condition met |

### CWE-1084 — [10]

| Finding id | Source | Reason |
| --- | --- | --- |
| 10 | `benchmarks/benchmark.py:477:54` | single function performs ≥3 `open(`/`.execute(` file or data-access operations — rule condition met |

### CWE-1124 — [14, 460, 536, 682, 706]

| Finding id | Source | Reason |
| --- | --- | --- |
| 14 | `benchmarks/benchmark.py:592:1` | statement is nested ≥6 control-flow levels — rule condition met |
| 460 | `scripts/darwin/setup.py:40:1` | statement is nested ≥6 control-flow levels — rule condition met |
| 536 | `setup.py:289:1` | statement is nested ≥6 control-flow levels — rule condition met |
| 682 | `tests/test_proxy_server.py:77:1` | statement is nested ≥6 control-flow levels — rule condition met |
| 706 | `tests/test_server.py:482:1` | statement is nested ≥6 control-flow levels — rule condition met |

### CWE-1046 — [19]

| Finding id | Source | Reason |
| --- | --- | --- |
| 19 | `benchmarks/benchmark.py:1555:1` | immutable text is concatenated with `+=` inside a loop — rule condition met |

### BP-PY-49 — [43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 59, 60, 64, 65, 66, 67, 68, 75, 79, 83, 84, 85, 86, 87, 91, 95, 100, 102, 104, 106, 110, 114, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 143, 147, 160, 161, 162, 163, 170, 174, 175, 176, 177, 181, 185, 186, 187, 188, 189, 190, 191, 198, 199, 200, 201, 204, 211, 215, 216, 220]

| Finding id | Source | Reason |
| --- | --- | --- |
| 43 | `benchmarks/libs/curl_cffi_bench.py:41:42` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 44 | `benchmarks/libs/curl_cffi_bench.py:61:49` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 45 | `benchmarks/libs/curl_cffi_bench.py:72:42` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 46 | `benchmarks/libs/curl_cffi_bench.py:82:42` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 47 | `benchmarks/libs/curl_cffi_bench.py:102:49` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 48 | `benchmarks/libs/curl_cffi_bench.py:112:42` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 49 | `benchmarks/libs/curl_cffi_bench.py:127:42` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 50 | `benchmarks/libs/curl_cffi_bench.py:142:43` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 51 | `benchmarks/libs/curl_cffi_bench.py:158:43` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 52 | `benchmarks/libs/curl_cffi_bench.py:174:42` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 59 | `benchmarks/libs/curl_cffi_bench.py:194:43` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 60 | `benchmarks/libs/curl_cffi_bench.py:210:43` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 64 | `benchmarks/libs/httpmorph_bench.py:51:35` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 65 | `benchmarks/libs/httpmorph_bench.py:61:41` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 66 | `benchmarks/libs/httpmorph_bench.py:71:42` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 67 | `benchmarks/libs/httpmorph_bench.py:85:35` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 68 | `benchmarks/libs/httpmorph_bench.py:99:52` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 75 | `benchmarks/libs/httpmorph_bench.py:126:54` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 79 | `benchmarks/libs/httpmorph_bench.py:155:54` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 83 | `benchmarks/libs/httpmorph_bench.py:178:35` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 84 | `benchmarks/libs/httpmorph_bench.py:188:41` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 85 | `benchmarks/libs/httpmorph_bench.py:198:42` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 86 | `benchmarks/libs/httpmorph_bench.py:211:35` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 87 | `benchmarks/libs/httpmorph_bench.py:225:52` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 91 | `benchmarks/libs/httpmorph_bench.py:247:54` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 95 | `benchmarks/libs/httpmorph_bench.py:272:54` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 100 | `benchmarks/libs/httpmorph_bench.py:314:47` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 102 | `benchmarks/libs/httpmorph_bench.py:333:48` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 104 | `benchmarks/libs/httpmorph_bench.py:354:41` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 106 | `benchmarks/libs/httpmorph_bench.py:379:1` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 110 | `benchmarks/libs/httpmorph_bench.py:410:1` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 114 | `benchmarks/libs/httpmorph_bench.py:443:1` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 118 | `benchmarks/libs/httpx_bench.py:46:19` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 119 | `benchmarks/libs/httpx_bench.py:55:19` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 120 | `benchmarks/libs/httpx_bench.py:64:19` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 121 | `benchmarks/libs/httpx_bench.py:73:31` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 122 | `benchmarks/libs/httpx_bench.py:83:50` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 123 | `benchmarks/libs/httpx_bench.py:95:51` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 124 | `benchmarks/libs/httpx_bench.py:107:63` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 125 | `benchmarks/libs/httpx_bench.py:118:19` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 126 | `benchmarks/libs/httpx_bench.py:127:19` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 127 | `benchmarks/libs/httpx_bench.py:136:19` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 128 | `benchmarks/libs/httpx_bench.py:145:31` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 129 | `benchmarks/libs/httpx_bench.py:155:50` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 130 | `benchmarks/libs/httpx_bench.py:167:51` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 131 | `benchmarks/libs/httpx_bench.py:179:63` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 132 | `benchmarks/libs/httpx_bench.py:190:30` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 133 | `benchmarks/libs/httpx_bench.py:199:30` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 134 | `benchmarks/libs/httpx_bench.py:208:30` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 135 | `benchmarks/libs/httpx_bench.py:217:42` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 136 | `benchmarks/libs/httpx_bench.py:231:28` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 143 | `benchmarks/libs/httpx_bench.py:251:29` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 147 | `benchmarks/libs/httpx_bench.py:271:41` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 160 | `benchmarks/libs/requests_bench.py:37:36` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 161 | `benchmarks/libs/requests_bench.py:47:42` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 162 | `benchmarks/libs/requests_bench.py:57:43` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 163 | `benchmarks/libs/requests_bench.py:71:44` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 170 | `benchmarks/libs/requests_bench.py:96:45` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 174 | `benchmarks/libs/requests_bench.py:117:36` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 175 | `benchmarks/libs/requests_bench.py:127:42` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 176 | `benchmarks/libs/requests_bench.py:137:43` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 177 | `benchmarks/libs/requests_bench.py:151:44` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 181 | `benchmarks/libs/requests_bench.py:170:45` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 185 | `benchmarks/libs/urllib3_bench.py:49:39` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 186 | `benchmarks/libs/urllib3_bench.py:60:39` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 187 | `benchmarks/libs/urllib3_bench.py:71:39` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 188 | `benchmarks/libs/urllib3_bench.py:82:39` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 189 | `benchmarks/libs/urllib3_bench.py:93:39` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 190 | `benchmarks/libs/urllib3_bench.py:104:39` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 191 | `benchmarks/libs/urllib3_bench.py:117:64` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 198 | `benchmarks/libs/urllib3_bench.py:136:65` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 199 | `benchmarks/libs/urllib3_bench.py:149:64` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 200 | `benchmarks/libs/urllib3_bench.py:162:65` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 201 | `benchmarks/libs/urllib_bench.py:35:19` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 204 | `benchmarks/libs/urllib_bench.py:56:19` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 211 | `benchmarks/libs/urllib_bench.py:92:19` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 215 | `benchmarks/libs/urllib_bench.py:111:19` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 216 | `benchmarks/libs/urllib_bench.py:132:19` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |
| 220 | `benchmarks/libs/urllib_bench.py:168:19` | `verify=False` / `cert_reqs="CERT_NONE"` disables TLS verification on an HTTP client call — rule condition met |

### CWE-295 — [202, 605]

| Finding id | Source | Reason |
| --- | --- | --- |
| 202 | `benchmarks/libs/urllib_bench.py:35:31` | `ssl.CERT_NONE` / `check_hostname = False` explicitly disables certificate validation — rule condition met |
| 605 | `tests/test_mock_server.py:184:35` | `ssl.CERT_NONE` / `check_hostname = False` explicitly disables certificate validation — rule condition met |

### CWE-523 — [203, 606]

| Finding id | Source | Reason |
| --- | --- | --- |
| 203 | `benchmarks/libs/urllib_bench.py:35:31` | certificate validation disabled while communicating over HTTPS transport — rule condition met |
| 606 | `tests/test_mock_server.py:184:35` | certificate validation disabled while communicating over HTTPS transport — rule condition met |

### BP-PY-14 — [228, 246, 261, 287, 288, 302, 311, 334]

| Finding id | Source | Reason |
| --- | --- | --- |
| 228 | `examples/advanced_features.py:27:20` | `session.`/`requests.` HTTP call without `timeout=` keyword — rule condition met |
| 246 | `examples/advanced_features.py:62:20` | `session.`/`requests.` HTTP call without `timeout=` keyword — rule condition met |
| 261 | `examples/advanced_features.py:103:20` | `session.`/`requests.` HTTP call without `timeout=` keyword — rule condition met |
| 287 | `examples/basic_usage.py:33:21` | `session.`/`requests.` HTTP call without `timeout=` keyword — rule condition met |
| 288 | `examples/basic_usage.py:34:21` | `session.`/`requests.` HTTP call without `timeout=` keyword — rule condition met |
| 302 | `examples/basic_usage.py:85:20` | `session.`/`requests.` HTTP call without `timeout=` keyword — rule condition met |
| 311 | `examples/basic_usage.py:134:17` | `session.`/`requests.` HTTP call without `timeout=` keyword — rule condition met |
| 334 | `examples/http2_example.py:27:16` | `session.`/`requests.` HTTP call without `timeout=` keyword — rule condition met |

### BP-PY-41 — [351, 365, 379, 394, 405, 414, 443, 482, 483, 486, 488, 490, 493, 495, 581, 582, 598, 599, 600, 603, 634, 637, 640, 643, 647, 673, 676, 688, 689, 690, 691, 692, 693, 694, 695]

| Finding id | Source | Reason |
| --- | --- | --- |
| 351 | `examples/proxy_example.py:70:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 365 | `examples/proxy_example.py:98:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 379 | `examples/proxy_example.py:132:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 394 | `examples/proxy_example.py:168:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 405 | `examples/proxy_example.py:200:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 414 | `examples/proxy_example.py:228:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 443 | `examples/proxy_example.py:313:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 482 | `scripts/test_local_build.py:17:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 483 | `scripts/test_local_build.py:34:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 486 | `scripts/test_local_build.py:56:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 488 | `scripts/test_local_build.py:84:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 490 | `scripts/test_local_build.py:105:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 493 | `scripts/test_local_build.py:132:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 495 | `scripts/test_local_build.py:176:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 581 | `tests/test_basic.py:26:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 582 | `tests/test_basic.py:79:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 598 | `tests/test_errors.py:94:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 599 | `tests/test_errors.py:100:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 600 | `tests/test_errors.py:106:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 603 | `tests/test_integration.py:202:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 634 | `tests/test_proxy.py:285:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 637 | `tests/test_proxy.py:354:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 640 | `tests/test_proxy.py:411:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 643 | `tests/test_proxy.py:439:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 647 | `tests/test_proxy.py:528:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 673 | `tests/test_proxy.py:1113:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 676 | `tests/test_proxy.py:1151:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 688 | `tests/test_requests_advanced.py:246:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 689 | `tests/test_requests_advanced.py:251:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 690 | `tests/test_requests_advanced.py:259:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 691 | `tests/test_requests_advanced.py:264:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 692 | `tests/test_requests_advanced.py:268:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 693 | `tests/test_requests_advanced.py:272:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 694 | `tests/test_requests_advanced.py:352:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |
| 695 | `tests/test_requests_compat.py:95:1` | test function performs calls with no `assert`/`pytest.raises`/`self.assert*` in its body — rule condition met |

### BP-PY-42 — [355, 484, 592, 601, 607]

| Finding id | Source | Reason |
| --- | --- | --- |
| 355 | `examples/proxy_example.py:76:1` | test function wraps a call in bare `try/except Exception` instead of `pytest.raises` — rule condition met |
| 484 | `scripts/test_local_build.py:39:1` | test function wraps a call in bare `try/except Exception` instead of `pytest.raises` — rule condition met |
| 592 | `tests/test_edge_cases_security.py:46:1` | test function wraps a call in bare `try/except Exception` instead of `pytest.raises` — rule condition met |
| 601 | `tests/test_http2.py:44:1` | test function wraps a call in bare `try/except Exception` instead of `pytest.raises` — rule condition met |
| 607 | `tests/test_proxy.py:34:1` | test function wraps a call in bare `try/except Exception` instead of `pytest.raises` — rule condition met |

### CWE-215 — [371]

| Finding id | Source | Reason |
| --- | --- | --- |
| 371 | `examples/proxy_example.py:105:9` | debug sink prints a sensitive identifier (`password`) — rule condition met |

### BP-PY-46 — [564]

| Finding id | Source | Reason |
| --- | --- | --- |
| 564 | `src/httpmorph/_async_client.py:313:21` | `print` in a library module (`src/httpmorph/`) outside any `__main__` guard — operational print in non-script module code, rule condition met |

### BP-PY-36 — [573]

| Finding id | Source | Reason |
| --- | --- | --- |
| 573 | `src/httpmorph/_client_c.py:974:1` | `Session()` created without `with`/`.close()` in scope — rule condition met |

### CWE-772 — [604]

| Finding id | Source | Reason |
| --- | --- | --- |
| 604 | `tests/test_mock_server.py:33:1` | `urlopen` response assigned with no same-function `.close(` — rule condition met |

### CWE-459 — [710]

| Finding id | Source | Reason |
| --- | --- | --- |
| 710 | `tests/test_server.py:739:18` | `NamedTemporaryFile(delete=False)` without same-function unlink — rule condition met |

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `./scripts/chunks`
- Function evidence: `./scripts/findings/functions`
- Validation: `git diff --check` — `pass`
