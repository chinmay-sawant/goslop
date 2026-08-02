# movielite — false-positive audit report

## Run metadata

```yaml
timestamp: 2026-08-02T00:00:00Z
repository: movielite
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite
branch: main
commit: b7cd21f75a7102fa3578f6f568cd7d344aa0f958
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite
chunk_path: ./scripts/chunks
function_context_path: ./scripts/findings/functions
```

## Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir real-repos/movielite/scripts/chunks -context-dir real-repos/movielite/scripts/findings/functions real-repos/movielite`
- Findings: `154`
- Chunks reviewed: `./scripts/chunks/Chunk_1_25.txt`, `Chunk_26_50.txt`, `Chunk_51_75.txt`, `Chunk_76_100.txt`, `Chunk_101_125.txt`, `Chunk_126_150.txt`, `Chunk_151_154.txt`
- Function contexts reviewed: `./scripts/findings/functions/<id>.txt` for every false positive (11, 26, 43, 49–138, 139, 143, 150, 154) plus the enclosing source files (benchmarks/compare_moviepy.py, examples/*.py, src/movielite/core/video_writer.py, tests/e2e/helpers.py, tests/e2e/test_core_video.py, docs/icon/movielite_gif.py)

## Audit checklist

- [x] Read every assigned chunk under `./scripts/chunks`.
- [x] Read `./scripts/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

Rule conditions verified with `./bin/goslop -explain <RULE-ID> --config templates/goslop-python.toml` (BP-PY-1/2/7/9/41/45/46, CWE-78/252/367/390/396/584/772/1046/1071/1124/1341, PERF-PY-25) and, for the narrow-pattern rules, against the detector implementations in `internal/lang/python/detectors/` (CWE-1341 `pyTierBDoubleCloseRE`, CWE-1046 `textAccumulatorEvidence`, CWE-367 `pyTierBTOCTOURE`, CWE-1124 control-frame counter, CWE-772 `pyResourceAssignmentRE`, BP-PY-46 skip logic, BP-PY-41 test-scope logic).

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 96 | 11, 26, 43, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 143, 150, 154 |
| True positive | 58 | 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 44, 45, 46, 47, 48, 128, 140, 141, 142, 144, 145, 146, 147, 148, 149, 151, 152, 153 |
| Uncertain | 0 | — |

## False positives

### [ ] Finding `11` — `CWE-1341`

- Function context: `./scripts/findings/functions/11.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/benchmarks/compare_moviepy.py:167:5`
- Checklist pattern: two `.close()` calls target two *different* objects, each released once.

Source excerpt:

```
    video.close()
    final.close()
```

Why this is a false positive: the detector regex `\w+\.close\(\s*\)[\s\S]{0,180}\w+\.close\(\s*\)` matches any two adjacent `.close()` calls; here `video` (a movielite `VideoClip`) and `final` (the moviepy writer result) are distinct handles closed once each — the same resource is not released twice.

Checklist evidence: rule condition "same resource handle is released twice" — the two receivers are different variables with independent lifecycles, so no double release exists.

### [ ] Finding `26` — `CWE-1341`

- Function context: `./scripts/findings/functions/26.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/benchmarks/compare_moviepy_v1.py:171:5`
- Checklist pattern: two `.close()` calls target two *different* objects, each released once.

Source excerpt:

```
    video.close()
    final.close()
```

Why this is a false positive: same construct as finding 11 in the v1 benchmark copy — `video` and `final` are distinct handles closed once each; the regex fires on the two close calls without checking the receivers are identical.

Checklist evidence: rule condition "same resource handle is released twice" — the two receivers are different variables, each released once.

### [ ] Finding `43` — `CWE-1046`

- Function context: `./scripts/findings/functions/43.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/docs/icon/movielite_gif.py:55:1`
- Checklist pattern: the `+=` accumulator is a numeric pixel-width sum, not immutable text.

Source excerpt:

```
    for i, letter in enumerate(word):
        ...
        clip = TextClip(letter, canvas=canvas, duration=DURATION)
        letter_clips.append(clip)
        total_text_width += clip.size[0]
```

Why this is a false positive: `total_text_width += clip.size[0]` accumulates integer pixel widths (`clip.size[0]` is an int); the rule's name heuristic fires because the name contains "text", but no string is concatenated and no immutable text is built in the loop.

Checklist evidence: `textAccumulatorEvidence` returns true on the name heuristic alone; the flagged RHS is `clip.size[0]` (int width), and the prior assignment search finds no string literal/str( initializer for `total_text_width` — the rule's "creation of immutable text using string concatenation" condition is not met.

### [ ] Finding `49` — `BP-PY-46`

- Function context: `./scripts/findings/functions/49.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/basic_editing.py:12:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_extract_subclip():
    """Extract a 10-second segment from a video."""
    print("Example 1: Extract subclip")
```

Why this is a false positive: `examples/` is a directory of runnable demo scripts (each has an `if __name__ == "__main__":` block); the print is intentional user-facing demo output, not operational logging in library code, so the rule condition "print is used for operational logging in non-script modules" is not satisfied.

Checklist evidence: the file is a standalone runnable example script, not an importable non-script module; the rule's stated condition ("non-script modules") does not apply.

### [ ] Finding `50` — `BP-PY-46`

- Function context: `./scripts/findings/functions/50.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/basic_editing.py:22:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    clip.close()
    print("Created output_subclip.mp4\n")
```

Why this is a false positive: same construct as finding 49 — user-facing progress output inside a demo function of a runnable example script, not library-code logging.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `51` — `BP-PY-46`

- Function context: `./scripts/findings/functions/51.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/basic_editing.py:27:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_concatenate_clips():
    """Concatenate multiple video clips sequentially."""
    print("Example 2: Concatenate clips")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `52` — `CWE-1341`

- Function context: `./scripts/findings/functions/52.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/basic_editing.py:44:5`
- Checklist pattern: two `.close()` calls target two *different* objects, each released once.

Source excerpt:

```
    clip1.close()
    clip2.close()
```

Why this is a false positive: `clip1` and `clip2` are separate `VideoClip` objects with independent lifecycles, each closed once; the regex matches the adjacent close calls without the receivers being identical.

Checklist evidence: rule condition "same resource handle is released twice" — the two receivers are different variables, each released once.

### [ ] Finding `53` — `BP-PY-46`

- Function context: `./scripts/findings/functions/53.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/basic_editing.py:47:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_concat.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `54` — `BP-PY-46`

- Function context: `./scripts/findings/functions/54.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/basic_editing.py:52:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_resize_video():
    """Resize a video to 720p."""
    print("Example 3: Resize video")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `55` — `BP-PY-46`

- Function context: `./scripts/findings/functions/55.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/basic_editing.py:62:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_720p.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `56` — `BP-PY-46`

- Function context: `./scripts/findings/functions/56.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/basic_editing.py:67:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_adjust_opacity():
    """Create a semi-transparent overlay effect."""
    print("Example 4: Adjust opacity")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `57` — `BP-PY-46`

- Function context: `./scripts/findings/functions/57.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/basic_editing.py:81:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_overlay.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `58` — `BP-PY-46`

- Function context: `./scripts/findings/functions/58.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/basic_editing.py:86:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_position_clip():
    """Position a clip at a specific location."""
    print("Example 5: Position clip")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `59` — `BP-PY-46`

- Function context: `./scripts/findings/functions/59.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/basic_editing.py:100:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_positioned.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `60` — `BP-PY-46`

- Function context: `./scripts/findings/functions/60.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/basic_editing.py:105:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_loop_video():
    """Loop a short video clip."""
    print("Example 7: Loop video")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `61` — `BP-PY-46`

- Function context: `./scripts/findings/functions/61.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/basic_editing.py:115:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_looped.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `62` — `BP-PY-46`

- Function context: `./scripts/findings/functions/62.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/composite_clips.py:16:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Example 1: Group logo and text to transform together")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `63` — `BP-PY-46`

- Function context: `./scripts/findings/functions/63.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/composite_clips.py:47:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_branding.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `64` — `BP-PY-46`

- Function context: `./scripts/findings/functions/64.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/composite_clips.py:54:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Example 2: Reusable composite element")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `65` — `BP-PY-46`

- Function context: `./scripts/findings/functions/65.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/composite_clips.py:87:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_captions.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `66` — `BP-PY-46`

- Function context: `./scripts/findings/functions/66.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/composite_clips.py:94:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Example 3: Transparent watermark with alpha composite")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `67` — `BP-PY-46`

- Function context: `./scripts/findings/functions/67.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/composite_clips.py:125:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_watermark.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `68` — `BP-PY-46`

- Function context: `./scripts/findings/functions/68.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/composite_clips.py:132:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Example 4: Animated composite with relative timing")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `69` — `BP-PY-46`

- Function context: `./scripts/findings/functions/69.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/composite_clips.py:163:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_intro.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `70` — `BP-PY-46`

- Function context: `./scripts/findings/functions/70.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/composite_clips.py:172:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Example 5: When NOT to use CompositeClip")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `71` — `BP-PY-46`

- Function context: `./scripts/findings/functions/71.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/composite_clips.py:175:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("  ❌ Wrong approach (using CompositeClip unnecessarily):")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `72` — `BP-PY-46`

- Function context: `./scripts/findings/functions/72.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/composite_clips.py:186:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("  ✓ Correct approach (using add_clips directly):")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `73` — `BP-PY-46`

- Function context: `./scripts/findings/functions/73.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/composite_clips.py:193:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_simple.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `74` — `BP-PY-46`

- Function context: `./scripts/findings/functions/74.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/composite_clips.py:194:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("  Use CompositeClip only when you need to transform clips as a group!")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `75` — `BP-PY-46`

- Function context: `./scripts/findings/functions/75.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/effects_showcase.py:12:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_fade_effects():
    """Demonstrate fade in and fade out effects."""
    print("Example 1: Fade effects")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `76` — `BP-PY-46`

- Function context: `./scripts/findings/functions/76.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/effects_showcase.py:23:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_fade.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `77` — `BP-PY-46`

- Function context: `./scripts/findings/functions/77.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/effects_showcase.py:28:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_blur_effects():
    """Demonstrate blur, blur-in, and blur-out effects."""
    print("Example 2: Blur effects")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `78` — `BP-PY-46`

- Function context: `./scripts/findings/functions/78.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/effects_showcase.py:57:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created blur effect videos\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `79` — `BP-PY-46`

- Function context: `./scripts/findings/functions/79.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/effects_showcase.py:62:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_color_effects():
    """Demonstrate color adjustment effects."""
    print("Example 3: Color effects")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `80` — `BP-PY-46`

- Function context: `./scripts/findings/functions/80.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/effects_showcase.py:109:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created color effect videos\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `81` — `BP-PY-46`

- Function context: `./scripts/findings/functions/81.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/effects_showcase.py:114:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_zoom_effects():
    """Demonstrate zoom and Ken Burns effects."""
    print("Example 4: Zoom effects")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `82` — `BP-PY-46`

- Function context: `./scripts/findings/functions/82.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/effects_showcase.py:148:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created zoom effect videos\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `83` — `BP-PY-46`

- Function context: `./scripts/findings/functions/83.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/effects_showcase.py:153:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_glitch_effects():
    """Demonstrate glitch and distortion effects."""
    print("Example 5: Glitch effects")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `84` — `BP-PY-46`

- Function context: `./scripts/findings/functions/84.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/effects_showcase.py:187:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created glitch effect videos\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `85` — `BP-PY-46`

- Function context: `./scripts/findings/functions/85.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/effects_showcase.py:192:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_vignette_effect():
    """Demonstrate vignette effect."""
    print("Example 6: Vignette effect")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `86` — `BP-PY-46`

- Function context: `./scripts/findings/functions/86.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/effects_showcase.py:202:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_vignette.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `87` — `BP-PY-46`

- Function context: `./scripts/findings/functions/87.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/effects_showcase.py:207:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_combined_effects():
    """Demonstrate combining multiple effects."""
    print("Example 7: Combined effects")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `88` — `BP-PY-46`

- Function context: `./scripts/findings/functions/88.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/effects_showcase.py:223:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_combined.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `89` — `BP-PY-46`

- Function context: `./scripts/findings/functions/89.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/effects_showcase.py:228:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_audio_effects():
    """Demonstrate audio fade effects."""
    print("Example 8: Audio effects")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `90` — `BP-PY-46`

- Function context: `./scripts/findings/functions/90.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/effects_showcase.py:247:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_audio_effects.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `91` — `BP-PY-46`

- Function context: `./scripts/findings/functions/91.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/masking_effects.py:14:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_simple_image_mask():
    """Apply a simple image mask to a video."""
    print("Example 1: Simple image mask")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `92` — `BP-PY-46`

- Function context: `./scripts/findings/functions/92.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/masking_effects.py:26:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_image_mask.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `93` — `BP-PY-46`

- Function context: `./scripts/findings/functions/93.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/masking_effects.py:31:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_text_mask_static():
    """Video visible only through text shape."""
    print("Example 2: Static text mask")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `94` — `BP-PY-46`

- Function context: `./scripts/findings/functions/94.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/masking_effects.py:48:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_text_mask.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `95` — `BP-PY-46`

- Function context: `./scripts/findings/functions/95.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/masking_effects.py:53:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_animated_text_mask():
    """Animated text mask with moving position and scaling."""
    print("Example 3: Animated text mask")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `96` — `BP-PY-46`

- Function context: `./scripts/findings/functions/96.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/masking_effects.py:79:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_animated_text_mask.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `97` — `BP-PY-46`

- Function context: `./scripts/findings/functions/97.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/masking_effects.py:84:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_shape_mask():
    """Create a custom shape mask using solid colors."""
    print("Example 4: Shape mask")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `98` — `BP-PY-46`

- Function context: `./scripts/findings/functions/98.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/masking_effects.py:117:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_shape_mask.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `99` — `BP-PY-46`

- Function context: `./scripts/findings/functions/99.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/text_animations.py:14:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_simple_text():
    """Add simple static text overlay."""
    print("Example 1: Simple text")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `100` — `BP-PY-46`

- Function context: `./scripts/findings/functions/100.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/text_animations.py:27:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_simple_text.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `101` — `BP-PY-46`

- Function context: `./scripts/findings/functions/101.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/text_animations.py:32:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_styled_text():
    """Add styled text with gradient and shadow."""
    print("Example 2: Styled text")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `102` — `BP-PY-46`

- Function context: `./scripts/findings/functions/102.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/text_animations.py:55:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_styled_text.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `103` — `BP-PY-46`

- Function context: `./scripts/findings/functions/103.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/text_animations.py:60:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_fade_text():
    """Text with fade in and fade out."""
    print("Example 3: Fade text")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `104` — `BP-PY-46`

- Function context: `./scripts/findings/functions/104.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/text_animations.py:75:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_fade_text.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `105` — `BP-PY-46`

- Function context: `./scripts/findings/functions/105.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/text_animations.py:80:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_animated_position():
    """Text with animated position (sliding in from left)."""
    print("Example 4: Animated position")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `106` — `BP-PY-46`

- Function context: `./scripts/findings/functions/106.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/text_animations.py:106:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_animated_pos.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `107` — `BP-PY-46`

- Function context: `./scripts/findings/functions/107.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/text_animations.py:111:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_bouncing_text():
    """Text with bouncing animation."""
    print("Example 5: Bouncing text")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `108` — `BP-PY-46`

- Function context: `./scripts/findings/functions/108.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/text_animations.py:136:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_bounce.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `109` — `BP-PY-46`

- Function context: `./scripts/findings/functions/109.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/text_animations.py:141:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_scaling_text():
    """Text with animated scale (zooming in)."""
    print("Example 6: Scaling text")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `110` — `BP-PY-46`

- Function context: `./scripts/findings/functions/110.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/text_animations.py:157:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_zoom_text.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `111` — `BP-PY-46`

- Function context: `./scripts/findings/functions/111.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/text_animations.py:162:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_multiple_text_layers():
    """Multiple text layers with different timings."""
    print("Example 7: Multiple text layers")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `112` — `BP-PY-46`

- Function context: `./scripts/findings/functions/112.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/text_animations.py:192:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_multi_text.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `113` — `BP-PY-46`

- Function context: `./scripts/findings/functions/113.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/text_animations.py:197:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_lower_third():
    """Create a lower third text overlay."""
    print("Example 8: Lower third")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `114` — `BP-PY-46`

- Function context: `./scripts/findings/functions/114.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/text_animations.py:230:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_lower_third.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `115` — `BP-PY-46`

- Function context: `./scripts/findings/functions/115.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/transitions.py:12:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_crossfade():
    """Simple crossfade transition between two clips."""
    print("Example 1: Crossfade transition")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `116` — `CWE-1341`

- Function context: `./scripts/findings/functions/116.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/transitions.py:27:5`
- Checklist pattern: two `.close()` calls target two *different* objects, each released once.

Source excerpt:

```
    clip1.close()
    clip2.close()
```

Why this is a false positive: `clip1` and `clip2` are separate `VideoClip` objects, each closed once; the regex matches the adjacent close calls without the receivers being identical.

Checklist evidence: rule condition "same resource handle is released twice" — the two receivers are different variables, each released once.

### [ ] Finding `117` — `BP-PY-46`

- Function context: `./scripts/findings/functions/117.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/transitions.py:29:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_crossfade.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `118` — `BP-PY-46`

- Function context: `./scripts/findings/functions/118.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/transitions.py:34:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_blur_dissolve():
    """Blur dissolve transition (blurs during transition)."""
    print("Example 2: Blur dissolve transition")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `119` — `BP-PY-46`

- Function context: `./scripts/findings/functions/119.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/transitions.py:50:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_blur_dissolve.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `120` — `BP-PY-46`

- Function context: `./scripts/findings/functions/120.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/transitions.py:55:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_multiple_transitions():
    """Multiple clips with transitions between each."""
    print("Example 3: Multiple transitions")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `121` — `BP-PY-46`

- Function context: `./scripts/findings/functions/121.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/transitions.py:75:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_multi_transitions.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `122` — `BP-PY-46`

- Function context: `./scripts/findings/functions/122.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/transitions.py:80:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_long_crossfade():
    """Longer crossfade for smoother transition."""
    print("Example 4: Long crossfade")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `123` — `BP-PY-46`

- Function context: `./scripts/findings/functions/123.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/transitions.py:96:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_long_crossfade.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `124` — `BP-PY-46`

- Function context: `./scripts/findings/functions/124.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/transitions.py:101:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_transition_with_effects():
    """Combine transitions with other effects."""
    print("Example 5: Transition with effects")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `125` — `BP-PY-46`

- Function context: `./scripts/findings/functions/125.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/transitions.py:123:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_transition_effects.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `126` — `BP-PY-46`

- Function context: `./scripts/findings/functions/126.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/transitions.py:128:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
def example_slideshow():
    """Create a slideshow with transitions."""
    print("Example 6: Slideshow with transitions")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `127` — `BP-PY-46`

- Function context: `./scripts/findings/functions/127.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/transitions.py:150:5`
- Checklist pattern: print in a runnable example script; the file is not a non-script library module.

Source excerpt:

```
    print("Created output_slideshow.mp4\n")
```

Why this is a false positive: user-facing demo output in a runnable example script, not operational logging in a non-script module.

Checklist evidence: the file is a standalone runnable example script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `129` — `BP-PY-46`

- Function context: `./scripts/findings/functions/129.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/video_audio_control.py:12:1`
- Checklist pattern: module-scope print in a standalone straight-line demo script.

Source excerpt:

```
# Example 1: Modify video audio independently
print("Example 1: Fade in/out on video audio")
```

Why this is a false positive: the file is a top-to-bottom runnable demo script (no importable library API); the module-level prints are the script's user-facing output, not operational logging in a non-script module.

Checklist evidence: zero-indent, module-level `print(` in a standalone script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `130` — `BP-PY-46`

- Function context: `./scripts/findings/functions/130.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/video_audio_control.py:23:1`
- Checklist pattern: module-scope print in a standalone straight-line demo script.

Source excerpt:

```
print("\nExample 2: Offset audio from video")
```

Why this is a false positive: module-level print in a straight-line demo script, intentional user-facing output, not library-code logging.

Checklist evidence: zero-indent, module-level `print(` in a standalone script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `131` — `BP-PY-46`

- Function context: `./scripts/findings/functions/131.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/video_audio_control.py:34:1`
- Checklist pattern: module-scope print in a standalone straight-line demo script.

Source excerpt:

```
print("\nExample 3: Video with background music")
```

Why this is a false positive: module-level print in a straight-line demo script, intentional user-facing output, not library-code logging.

Checklist evidence: zero-indent, module-level `print(` in a standalone script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `132` — `BP-PY-46`

- Function context: `./scripts/findings/functions/132.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/video_audio_control.py:50:1`
- Checklist pattern: module-scope print in a standalone straight-line demo script.

Source excerpt:

```
print("\nExample 4: Video timing changes sync to audio")
```

Why this is a false positive: module-level print in a straight-line demo script, intentional user-facing output, not library-code logging.

Checklist evidence: zero-indent, module-level `print(` in a standalone script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `133` — `BP-PY-46`

- Function context: `./scripts/findings/functions/133.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/video_audio_control.py:63:1`
- Checklist pattern: module-scope print in a standalone straight-line demo script.

Source excerpt:

```
print("\nAll examples completed!")
```

Why this is a false positive: module-level print in a straight-line demo script, intentional user-facing output, not library-code logging.

Checklist evidence: zero-indent, module-level `print(` in a standalone script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `134` — `BP-PY-46`

- Function context: `./scripts/findings/functions/134.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/video_audio_control.py:64:1`
- Checklist pattern: module-scope print in a standalone straight-line demo script.

Source excerpt:

```
print("\nKey features:")
```

Why this is a false positive: module-level print in a straight-line demo script, intentional user-facing output, not library-code logging.

Checklist evidence: zero-indent, module-level `print(` in a standalone script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `135` — `BP-PY-46`

- Function context: `./scripts/findings/functions/135.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/video_audio_control.py:65:1`
- Checklist pattern: module-scope print in a standalone straight-line demo script.

Source excerpt:

```
print("- video.audio: Access audio track")
```

Why this is a false positive: module-level print in a straight-line demo script, intentional user-facing output, not library-code logging.

Checklist evidence: zero-indent, module-level `print(` in a standalone script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `136` — `BP-PY-46`

- Function context: `./scripts/findings/functions/136.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/video_audio_control.py:66:1`
- Checklist pattern: module-scope print in a standalone straight-line demo script.

Source excerpt:

```
print("- video.audio.fade_in/fade_out/set_volume: Apply audio effects")
```

Why this is a false positive: module-level print in a straight-line demo script, intentional user-facing output, not library-code logging.

Checklist evidence: zero-indent, module-level `print(` in a standalone script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `137` — `BP-PY-46`

- Function context: `./scripts/findings/functions/137.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/video_audio_control.py:67:1`
- Checklist pattern: module-scope print in a standalone straight-line demo script.

Source excerpt:

```
print("- video.audio._start: Move audio independently")
```

Why this is a false positive: module-level print in a straight-line demo script, intentional user-facing output, not library-code logging.

Checklist evidence: zero-indent, module-level `print(` in a standalone script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `138` — `BP-PY-46`

- Function context: `./scripts/findings/functions/138.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/examples/video_audio_control.py:68:1`
- Checklist pattern: module-scope print in a standalone straight-line demo script.

Source excerpt:

```
print("- video.set_start/set_duration: Automatically syncs audio")
```

Why this is a false positive: module-level print in a straight-line demo script, intentional user-facing output, not library-code logging.

Checklist evidence: zero-indent, module-level `print(` in a standalone script; the "non-script modules" condition of the rule is not satisfied.

### [ ] Finding `139` — `CWE-1341`

- Function context: `./scripts/findings/functions/139.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/src/movielite/core/video_writer.py:253:17`
- Checklist pattern: two `.close()` calls target two *different* objects, each released once.

Source excerpt:

```
                for clip in clips_to_close:
                    clip.close()
                    remaining_clips_to_process.remove(clip)
                ...
        process.stdin.close()
        process.wait()
```

Why this is a false positive: the regex pair is `clip.close()` (clip handles) and `process.stdin.close()` (the ffmpeg subprocess stdin) — different resources, each released once; `process.stdin` is then waited on, not closed again.

Checklist evidence: rule condition "same resource handle is released twice" — the two receivers in the matched pair are different variables, each released once.

### [ ] Finding `143` — `CWE-367`

- Function context: `./scripts/findings/functions/143.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/src/movielite/core/video_writer.py:403:20`
- Checklist pattern: the exists-check is a teardown guard on the function's own temp file; the race has no consequence.

Source excerpt:

```
            finally:
                if os.path.exists(temp_audio_path):
                    os.remove(temp_audio_path)
```

Why this is a false positive: the exists-check guards `os.remove` of the function's own `tempfile.NamedTemporaryFile` path in a single-threaded `finally` cleanup; if the file vanished between check and remove the cleanup intent is unchanged and the failure is benign, so the TOCTOU has no exploitable consequence.

Checklist evidence: the checked path is a locally created temp file in a `finally` teardown (same shape as the teardown-guard false positives recorded in the pycaps audit); no attacker-controlled path is checked before a security-relevant use.

### [ ] Finding `150` — `CWE-1341`

- Function context: `./scripts/findings/functions/150.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/tests/e2e/helpers.py:64:18`
- Checklist pattern: two `.close()` calls target two *different* objects, each released once.

Source excerpt:

```
    def close(self):
        try:
            self._capture.close()
        finally:
            return self._real.close()
```

Why this is a false positive: `self._capture` (the capture file) and `self._real` (the ffmpeg stdin wrapper target) are distinct handles, each closed once; the regex matches the two close calls without the receivers being identical.

Checklist evidence: rule condition "same resource handle is released twice" — the two receivers are different attributes, each released once.

### [ ] Finding `154` — `CWE-1341`

- Function context: `./scripts/findings/functions/154.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/movielite/tests/e2e/test_core_video.py:29:5`
- Checklist pattern: two `.close()` calls target two *different* objects, each released once.

Source excerpt:

```
def test_passthrough(bg_video, output):
    clip = VideoClip(str(bg_video))
    _write(clip, output.mp4)
    clip.close()
    assert_matches_golden(output, "passthrough.mp4")

def test_subclip(bg_video, output):
    clip = VideoClip(str(bg_video)).subclip(0.5, 1.5)
    _write(clip, output.mp4)
    clip.close()
```

Why this is a false positive: the regex pair spans `clip.close()` in `test_passthrough` and `clip.close()` in `test_subclip` — each function constructs and closes its own `clip` object once; no handle is closed twice.

Checklist evidence: rule condition "same resource handle is released twice" — the two close calls operate on different objects created in different functions.

## Uncertain findings

None.

## True positives

### `BP-PY-1` — Bare Except Clause (7)

| Finding | Source | Reason |
| --- | --- | --- |
| 1 | `benchmarks/compare_moviepy.py:36:1` | `except Exception as e:` records the error but does not re-raise; broad catch per rule condition. |
| 16 | `benchmarks/compare_moviepy_v1.py:40:1` | Same `except Exception as e:` construct as finding 1. |
| 31 | `benchmarks/read_backends/backends.py:190:1` | `except Exception: pass` around `self._proc.stdout.close()`; broad catch, no re-raise. |
| 40 | `benchmarks/read_backends/bench.py:297:1` | `except Exception as e:` prints and returns; broad catch, no re-raise. |
| 141 | `src/movielite/core/video_writer.py:351:1` | `except Exception as e:` logs a warning; broad catch, no re-raise. |
| 144 | `src/movielite/core/video_writer.py:406:1` | `except Exception as e:` logs an error; broad catch, no re-raise. |
| 145 | `src/movielite/video/readers/ffmpeg_reader.py:147:1` | `except Exception: pass` around `self._proc.stdout.close()`; broad catch, no re-raise. |

### `BP-PY-2` — Except Pass (4)

| Finding | Source | Reason |
| --- | --- | --- |
| 32 | `benchmarks/read_backends/backends.py:190:1` | Handler body is solely `pass`. |
| 37 | `benchmarks/read_backends/bench.py:194:1` | `except (FileNotFoundError, ProcessLookupError, PermissionError): pass` — suite is solely `pass`. |
| 39 | `benchmarks/read_backends/bench.py:210:1` | `except (FileNotFoundError, ProcessLookupError): pass` — suite is solely `pass`. |
| 146 | `src/movielite/video/readers/ffmpeg_reader.py:147:1` | Handler body is solely `pass`. |

### `BP-PY-7` — open Without Context Manager (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 153 | `tests/e2e/helpers.py:89:23` | `capture = open(capture_path, "ab")` — `open(` not preceded by `with`; no same-function close (release is delegated to the `_StdinTee` wrapper). |

### `BP-PY-9` — os.system Or os.popen (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 46 | `docs/icon/movielite_gif.py:78:5` | `os.system(f'ffmpeg ...')` — direct `os.system` call per rule condition. |

### `BP-PY-41` — pytest assert With Side Effects Only (24)

| Finding | Source | Reason |
| --- | --- | --- |
| 3 | `benchmarks/compare_moviepy.py:43:1` | `def test_no_processing_movielite` calls production code (`VideoClip`, `writer.write()`), contains no assert. |
| 4 | `benchmarks/compare_moviepy.py:54:1` | `def test_no_processing_moviepy` — side-effect calls, no assert. |
| 5 | `benchmarks/compare_moviepy.py:68:1` | `def test_video_zoom_movielite` — side-effect calls, no assert. |
| 6 | `benchmarks/compare_moviepy.py:86:1` | `def test_video_zoom_moviepy` — side-effect calls, no assert. |
| 7 | `benchmarks/compare_moviepy.py:108:1` | `def test_fade_movielite` — side-effect calls, no assert. |
| 8 | `benchmarks/compare_moviepy.py:120:1` | `def test_fade_moviepy` — side-effect calls, no assert. |
| 9 | `benchmarks/compare_moviepy.py:136:1` | `def test_text_overlay_movielite` — side-effect calls, no assert. |
| 10 | `benchmarks/compare_moviepy.py:151:1` | `def test_text_overlay_moviepy` — side-effect calls, no assert. |
| 12 | `benchmarks/compare_moviepy.py:172:1` | `def test_video_overlay_movielite` — side-effect calls, no assert. |
| 13 | `benchmarks/compare_moviepy.py:188:1` | `def test_video_overlay_moviepy` — side-effect calls, no assert. |
| 14 | `benchmarks/compare_moviepy.py:212:1` | `def test_alpha_overlay_movielite` — side-effect calls, no assert. |
| 15 | `benchmarks/compare_moviepy.py:232:1` | `def test_alpha_overlay_moviepy` — side-effect calls, no assert. |
| 18 | `benchmarks/compare_moviepy_v1.py:47:1` | `def test_no_processing_movielite` — side-effect calls, no assert. |
| 19 | `benchmarks/compare_moviepy_v1.py:58:1` | `def test_no_processing_moviepy` — side-effect calls, no assert. |
| 20 | `benchmarks/compare_moviepy_v1.py:72:1` | `def test_video_zoom_movielite` — side-effect calls, no assert. |
| 21 | `benchmarks/compare_moviepy_v1.py:90:1` | `def test_video_zoom_moviepy` — side-effect calls, no assert. |
| 22 | `benchmarks/compare_moviepy_v1.py:112:1` | `def test_fade_movielite` — side-effect calls, no assert. |
| 23 | `benchmarks/compare_moviepy_v1.py:124:1` | `def test_fade_moviepy` — side-effect calls, no assert. |
| 24 | `benchmarks/compare_moviepy_v1.py:140:1` | `def test_text_overlay_movielite` — side-effect calls, no assert. |
| 25 | `benchmarks/compare_moviepy_v1.py:155:1` | `def test_text_overlay_moviepy` — side-effect calls, no assert. |
| 27 | `benchmarks/compare_moviepy_v1.py:176:1` | `def test_video_overlay_movielite` — side-effect calls, no assert. |
| 28 | `benchmarks/compare_moviepy_v1.py:192:1` | `def test_video_overlay_moviepy` — side-effect calls, no assert. |
| 29 | `benchmarks/compare_moviepy_v1.py:216:1` | `def test_alpha_overlay_movielite` — side-effect calls, no assert. |
| 30 | `benchmarks/compare_moviepy_v1.py:236:1` | `def test_alpha_overlay_moviepy` — side-effect calls, no assert. |

### `BP-PY-45` — sys.path Mutation At Runtime (2)

| Finding | Source | Reason |
| --- | --- | --- |
| 36 | `benchmarks/read_backends/bench.py:48:1` | `sys.path.insert(0, str(Path(__file__).resolve().parent))` — runtime `sys.path` mutation per rule condition. |
| 128 | `examples/video_audio_control.py:7:1` | `sys.path.insert(0, 'src')` — runtime `sys.path` mutation per rule condition. |

### `CWE-78` — OS Command Injection (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 48 | `docs/icon/movielite_gif.py:78:5` | `os.system(f'ffmpeg ... {VIDEO_FILENAME} ... {GIF_FILENAME}')` — f-string is a dynamic expression reaching the `os.system` sink per rule condition. |

### `CWE-252` — Unchecked Return Value (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 47 | `docs/icon/movielite_gif.py:78:5` | `os.system(...)` used as a standalone statement; return status discarded without check. |

### `CWE-390` — Detection of Error Condition Without Action (3)

| Finding | Source | Reason |
| --- | --- | --- |
| 33 | `benchmarks/read_backends/backends.py:190:1` | `except Exception: pass` — error detected, no action. |
| 38 | `benchmarks/read_backends/bench.py:194:1` | `except (FileNotFoundError, ProcessLookupError, PermissionError): pass` — error detected, no action. |
| 147 | `src/movielite/video/readers/ffmpeg_reader.py:147:1` | `except Exception: pass` — error detected, no action. |

### `CWE-396` — Declaration of Catch for Generic Exception (6)

| Finding | Source | Reason |
| --- | --- | --- |
| 2 | `benchmarks/compare_moviepy.py:36:1` | `except Exception as e:` — generic catch per rule condition. |
| 17 | `benchmarks/compare_moviepy_v1.py:40:1` | `except Exception as e:` — generic catch per rule condition. |
| 34 | `benchmarks/read_backends/backends.py:190:1` | `except Exception:` — generic catch per rule condition. |
| 41 | `benchmarks/read_backends/bench.py:297:1` | `except Exception as e:` — generic catch per rule condition. |
| 142 | `src/movielite/core/video_writer.py:351:1` | `except Exception as e:` — generic catch per rule condition. |
| 148 | `src/movielite/video/readers/ffmpeg_reader.py:147:1` | `except Exception:` — generic catch per rule condition. |

### `CWE-584` — Return Inside Finally Block (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 151 | `tests/e2e/helpers.py:66:1` | `return self._real.close()` is a direct return in the `finally` suite — a `close()` error from the protected block is suppressed per rule condition. |

### `CWE-772` — Missing Release of Resource after Effective Lifetime (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 152 | `tests/e2e/helpers.py:89:1` | `capture = open(capture_path, "ab")` assigned with no same-function `capture.close(` and no context manager; the close is only delegated to the `_StdinTee` wrapper, so the rule's same-function condition is met. |

### `CWE-1071` — Empty Code Block (2)

| Finding | Source | Reason |
| --- | --- | --- |
| 35 | `benchmarks/read_backends/backends.py:190:1` | `except Exception: pass` — handler silently contains only `pass`. |
| 149 | `src/movielite/video/readers/ffmpeg_reader.py:147:1` | `except Exception: pass` — handler silently contains only `pass`. |

### `CWE-1124` — Excessively Deep Nesting (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 140 | `src/movielite/core/video_writer.py:315:1` | `pbar.update(1)` sits at six control-flow levels (`try` → `with` → `for` → `try` → `for` → `if`), meeting the ≥6-level condition. |

### `PERF-PY-25` — Heavy Object Construction Per Homogeneous Element (3)

| Finding | Source | Reason |
| --- | --- | --- |
| 42 | `docs/icon/movielite_gif.py:33:1` | `particle.set_position(lambda t, ...)` — lambda constructed inside a 80-iteration `for` loop per rule condition. |
| 44 | `docs/icon/movielite_gif.py:64:1` | `clip.set_scale(lambda t, p=phase: ...)` — lambda constructed per letter in the `for i, clip in enumerate(...)` loop. |
| 45 | `docs/icon/movielite_gif.py:65:1` | `clip.set_opacity(lambda t, p=phase: ...)` — lambda constructed per letter in the same loop. |

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `./scripts/chunks`
- Function evidence: `./scripts/findings/functions`
- Validation: `git diff --check` — pass
