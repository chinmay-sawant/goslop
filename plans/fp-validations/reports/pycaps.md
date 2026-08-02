# False-positive audit — pycaps

## Run metadata

```yaml
timestamp: 2026-08-02T07:23:46Z
repository: pycaps
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps
branch: main
commit: 68a6843bbd3f7d0ea33047ffb1045d978a7671c4
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps
chunk_path: scripts/pycaps/chunks
function_context_path: scripts/pycaps/findings/functions
```

## Scan evidence

- Build command: `(pre-built) ./bin/goslop`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/pycaps/chunks -context-dir scripts/pycaps/findings/functions real-repos/pycaps`
- Findings: `51`
- Chunks reviewed: `scripts/pycaps/chunks/Chunk_1_25.txt`, `scripts/pycaps/chunks/Chunk_26_50.txt`, `scripts/pycaps/chunks/Chunk_51_51.txt`
- Function contexts reviewed: `scripts/pycaps/findings/functions/1.txt … 51.txt`

## Audit checklist

- [x] Read every assigned chunk under `scripts/pycaps/chunks`.
- [x] Read `scripts/pycaps/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 23 | 1, 10, 13, 15, 16, 20, 23, 25, 28, 30, 31, 35, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 51 |
| True positive | 28 | 2, 3, 4, 5, 6, 7, 8, 9, 11, 12, 14, 17, 18, 19, 21, 22, 24, 26, 27, 29, 32, 33, 34, 36, 37, 38, 49, 50 |
| Uncertain | 0 | — |

## False positives

### [ ] Finding 1 — BP-PY-13

- Function context: `scripts/pycaps/findings/functions/1.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/ai/gpt.py:6:1`
- Checklist pattern: the flagged string is an environment-variable *name* used as the lookup key for `os.getenv`, not a credential value.

Source excerpt:

```
    OPENAI_API_KEY_NAME = "PYCAPS_OPENAI_API_KEY"

    def is_enabled(self) -> bool:
        return os.getenv(self.OPENAI_API_KEY_NAME) is not None

    ...
            self._client = OpenAI(api_key=os.getenv(self.OPENAI_API_KEY_NAME))
```

Why this is a false positive: the heuristic fires on the identifier containing `api_key` and a non-empty string RHS, but the literal is the name of the environment variable the code reads the key from (`os.getenv(self.OPENAI_API_KEY_NAME)`) — exactly the fix the rule recommends; no credential value is hardcoded.

Checklist evidence: the rule condition "secret-like name assigned a string literal" is matched only on the identifier pattern; the value is a lookup key for `os.getenv`, not a loaded-from-environment-able secret.

### [ ] Finding 10 — CWE-409

- Function context: `scripts/pycaps/findings/functions/10.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/effect/clip/animate_segment_emojis_effect.py:71:18`
- Checklist pattern: the archive is fetched from a hardcoded, version-pinned URL constant of the project's own release, so no untrusted party supplies the compressed data.

Source excerpt:

```
    ASSETS_ZIP_URL = "https://github.com/francozanardi/pycaps/releases/download/emoji-assets-v2/animated_emojis.zip"

    ...
            file_buffer.seek(0)
            with zipfile.ZipFile(file_buffer) as z:
                z.extractall(self.CACHE_DIR)
```

Why this is a false positive: the detector reports any `.extractall` call, but the archive content is pinned to the project's own HTTPS release asset via a source constant; the rule's own fix says to "validate the relevant trust boundary", and no user-influenced input reaches the extraction sink.

Checklist evidence: the zip bomb trust boundary (untrusted compressed data) is absent — the archive source is a compile-time constant, not user input.

### [ ] Finding 13 — BP-PY-7

- Function context: `scripts/pycaps/findings/functions/13.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/pipeline/caps_pipeline.py:78:23`
- Checklist pattern: `.open(` is a method call on the `SubtitleRenderer` object, not Python's builtin `open()`.

Source excerpt:

```
        resources_dir = Path(self._resources_dir) if self._resources_dir else None
        self._renderer.open(self._video_width, self._video_height, resources_dir, self._cache_strategy)
```

Why this is a false positive: `self._renderer.open(...)` is an application lifecycle method (defined in `subtitle_renderer.py:15`), not a builtin `open()` call, so no file handle is leaked by omitting `with`.

Checklist evidence: the detector's `.open(` heuristic matches method calls; the shown call is on `self._renderer`, whose `open` is a `def open(self, ...)` method, not the file-opening builtin.

### [ ] Finding 15 — CWE-1341

- Function context: `scripts/pycaps/findings/functions/15.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/pipeline/caps_pipeline.py:237:14`
- Checklist pattern: two `.close()` calls target two *different* objects, each released once.

Source excerpt:

```
        logger().debug("Cleaning up pipeline resources...")
        self._video_generator.close()
        self._renderer.close()
```

Why this is a false positive: the detector regex `\w+\.close\(\s*\)[\s\S]{0,180}\w+\.close\(\s*\)` matches any two adjacent `.close()` calls; here `_video_generator` and `_renderer` are distinct handles closed once each — the same resource is not released twice.

Checklist evidence: rule condition "same resource handle is released twice" — the two receivers are different attributes with independent lifecycle, so no double release exists.

### [ ] Finding 16 — CWE-367

- Function context: `scripts/pycaps/findings/functions/16.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/pipeline/caps_pipeline_builder.py:49:16`
- Checklist pattern: the `exists` check only feeds a user-facing `ValueError`; the subsequent use fails benignly.

Source excerpt:

```
    def add_css(self, css_file_path: str) -> "CapsPipelineBuilder":
        if not os.path.exists(css_file_path):
            raise ValueError(f"CSS file not found: {css_file_path}")
        css_content = open(css_file_path, "r", encoding="utf-8").read()
```

Why this is a false positive: the check-then-use race has no privilege or integrity consequence — if the file vanishes between check and `open`, `open()` simply raises `FileNotFoundError` with the same benign outcome; no attacker-controlled state is used and no security boundary is crossed.

Checklist evidence: rule condition "path checked before later separate use" is matched literally, but the check exists only to produce a friendly error and the race's failure mode is a benign exception on the user's own file.

### [ ] Finding 20 — CWE-478

- Function context: `scripts/pycaps/findings/functions/20.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/pipeline/json_config_loader.py:93:1`
- Checklist pattern: the match domain is a pydantic-validated `Literal` union whose every member has a case branch.

Source excerpt:

```
            match effect.type:
                case "emoji_in_segment":
                    ...
                case "emoji_in_word":
                    ...
                case "remove_punctuation_marks":
                    ...
                case "typewriting":
                    ...
                case "animate_segment_emojis":
                    ...
```

Why this is a false positive: `effect.type` is a pydantic `Literal["emoji_in_segment", "emoji_in_word", "remove_punctuation_marks", "typewriting", "animate_segment_emojis"]` (`json_schema.py:51-74`), and all five members are handled by the five case branches; any other value raises `ValidationError` before the match is reached, so no unhandled value can occur.

Checklist evidence: rule condition "match with two or more cases lacks a wildcard" is met syntactically, but the matched value is exhaustively constrained by the schema validation, making a default branch unreachable dead code.

### [ ] Finding 23 — BP-PY-7

- Function context: `scripts/pycaps/findings/functions/23.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/renderer/css_subtitle_renderer.py:54:9`
- Checklist pattern: `def open(self, ...)` is a method *definition*, not an `open()` call.

Source excerpt:

```
    def open(self, video_width: int, video_height: int, resources_dir: Optional[Path] = None, cache_strategy: CacheStrategy = CacheStrategy.CSS_CLASSES_AWARE):
        """Initializes Playwright and loads the base HTML page."""
```

Why this is a false positive: the line defines a renderer lifecycle method named `open`; the detector's `open(` identifier match fires on the definition signature, and no file is opened anywhere on the line.

Checklist evidence: rule condition "a file is opened without a `with` statement" — no `open()` call exists on the flagged line; it is a method declaration.

### [ ] Finding 25 — BP-PY-7

- Function context: `scripts/pycaps/findings/functions/25.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/renderer/pictex_subtitle_renderer.py:39:9`
- Checklist pattern: `def open(self, ...)` is a method *definition*, not an `open()` call.

Source excerpt:

```
    def open(self, video_width: int, video_height: int, resources_dir: Optional[Path] = None, cache_strategy: CacheStrategy = CacheStrategy.CSS_CLASSES_AWARE):
        scale_modifier = self._calculate_scale_modifier(video_height)
```

Why this is a false positive: the flagged line declares the renderer's `open` lifecycle method; no file is opened, so there is no handle to leak.

Checklist evidence: rule condition "a file is opened without a `with` statement" — the line is a method definition, not a file-open call.

### [ ] Finding 28 — BP-PY-7

- Function context: `scripts/pycaps/findings/functions/28.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/renderer/playwright_screenshot_capturer.py:37:22`
- Checklist pattern: `Image.open` is a Pillow library method on an in-memory `BytesIO`, not the builtin `open()`.

Source excerpt:

```
        png_bytes = page.screenshot(omit_background=True, type="png", animations="disabled", scale="device", clip=clip)
        image = Image.open(io.BytesIO(png_bytes)).convert("RGBA")
```

Why this is a false positive: `Image.open` is Pillow's decoder entry point fed by an in-memory `io.BytesIO` (never a path), so the rule's `with open(...)` fix does not apply; no filesystem handle is left open.

Checklist evidence: rule condition "file opened without context manager" targets the builtin `open()`; the flagged call is `Image.open(BytesIO(...))` — a library method with no filesystem resource.

### [ ] Finding 30 — BP-PY-7

- Function context: `scripts/pycaps/findings/functions/30.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/renderer/subtitle_renderer.py:15:9`
- Checklist pattern: an `@abstractmethod def open(self, ...)` declaration, not an `open()` call.

Source excerpt:

```
    @abstractmethod
    def open(self, video_width: int, video_height: int, resources_dir: Optional[Path] = None, cache_strategy: CacheStrategy = CacheStrategy.CSS_CLASSES_AWARE):
        pass
```

Why this is a false positive: the line is an abstract interface declaration whose body is `pass`; there is no file operation at all.

Checklist evidence: rule condition "a file is opened without a `with` statement" — no `open()` call; the identifier match lands on an abstract method signature.

### [ ] Finding 31 — CWE-478

- Function context: `scripts/pycaps/findings/functions/31.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/selector/time_event_selector.py:28:1`
- Checklist pattern: the matched value is a Python `Enum` whose members are all covered by the case branches.

Source excerpt:

```
        match self._element_type:
            case ElementType.WORD:
                return self.__filter_by_words(clips)
            case ElementType.LINE:
                return self.__filter_by_lines(clips)
            case ElementType.SEGMENT:
                return self.__filter_by_segments(clips)
```

Why this is a false positive: `ElementType` is a `str, Enum` with exactly `WORD`, `LINE`, `SEGMENT` members (`common/types.py:21-24`), all three handled; invalid values cannot construct the enum, so no value reaches the match unhandled.

Checklist evidence: rule condition "match with multiple cases and no wildcard" is met syntactically, but the enum domain is fully covered — a default branch would be unreachable.

### [ ] Finding 35 — CWE-186

- Function context: `scripts/pycaps/findings/functions/35.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/transcriber/transcript_loader.py:14:23`
- Checklist pattern: the regex parses a fixed standard format (VTT inline timestamp), not a validator that rejects valid identifiers.

Source excerpt:

```
_TIME_EPSILON = 0.01
_VTT_INLINE_TIME_RE = re.compile(r"<((?:\d{2}:)?\d{2}:\d{2}\.\d{3})>")
```

Why this is a false positive: the detector fires on the `\d{3}` quantifier, but this regex extracts WebVTT inline timestamps, whose format is exactly `[HH:]MM:SS.mmm` per the VTT spec; the millisecond part is always three digits, so the pattern is precise, not overly restrictive — no valid input is rejected.

Checklist evidence: rule condition "validation regex accepts only a fixed narrow identifier shape" — this is a parser for a standards-defined fixed format (used to extract, not validate user identifiers), so narrowness is inherent to the format, not an over-restriction.

### [ ] Finding 39 — PERF-PY-25

- Function context: `scripts/pycaps/findings/functions/39.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/video/subtitle_clips_generator.py:27:1`
- Checklist pattern: the lambda closes over per-iteration state and cannot be hoisted out of the loop.

Source excerpt:

```
            for segment in document.segments:
                for line in segment.lines:
                    self.__generate_word_clips_for_line(
                        line,
                        ElementState.LINE_NOT_NARRATED_YET,
                        ElementState.WORD_NOT_NARRATED_YET,
                        lambda _: segment.time.start,
                        lambda _: line.time.start,
                        pbar
                    )
```

Why this is a false positive: each lambda captures the current `segment`/`line` iteration values (`segment.time.start`, `line.time.start`), so constructing them per iteration is semantically required; no "fixed-schema path" exists and there is no avoidable per-element work.

Checklist evidence: rule condition "heavy object or lambda is constructed per homogeneous loop element" — the closures are data-dependent on loop variables, so pre-construction outside the loop is impossible; the allocation is necessary.

### [ ] Finding 40 — PERF-PY-25

- Function context: `scripts/pycaps/findings/functions/40.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/video/subtitle_clips_generator.py:28:1`
- Checklist pattern: same per-iteration closure as Finding 39.

Source excerpt:

```
                    self.__generate_word_clips_for_line(
                        line,
                        ElementState.LINE_NOT_NARRATED_YET,
                        ElementState.WORD_NOT_NARRATED_YET,
                        lambda _: segment.time.start,
                        lambda _: line.time.start,
                        pbar
                    )
```

Why this is a false positive: `lambda _: line.time.start` captures the current `line` from the enclosing loop; the value differs every iteration, so the lambda cannot be hoisted and the per-iteration construction is required.

Checklist evidence: rule condition "lambda constructed per homogeneous loop element" — closure over a mutable loop variable prevents hoisting; no avoidable allocation.

### [ ] Finding 41 — PERF-PY-25

- Function context: `scripts/pycaps/findings/functions/41.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/video/subtitle_clips_generator.py:36:1`
- Checklist pattern: same per-iteration closure as Finding 39.

Source excerpt:

```
                    self.__generate_word_clips_for_line(
                        line,
                        ElementState.LINE_BEING_NARRATED,
                        ElementState.WORD_NOT_NARRATED_YET,
                        lambda _: line.time.start,
                        lambda word: word.time.start,
                        pbar
                    )
```

Why this is a false positive: `lambda _: line.time.start` must be created per line because `line` changes on every iteration; there is no hoistable construction to eliminate.

Checklist evidence: rule condition "lambda constructed per homogeneous loop element" — the captured `line` is loop-dependent; allocation cannot be moved out of the loop.

### [ ] Finding 42 — PERF-PY-25

- Function context: `scripts/pycaps/findings/functions/42.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/video/subtitle_clips_generator.py:37:1`
- Checklist pattern: same per-iteration closure as Finding 39.

Source excerpt:

```
                    self.__generate_word_clips_for_line(
                        line,
                        ElementState.LINE_BEING_NARRATED,
                        ElementState.WORD_NOT_NARRATED_YET,
                        lambda _: line.time.start,
                        lambda word: word.time.start,
                        pbar
                    )
```

Why this is a false positive: `lambda word: word.time.start` is a per-word selector that must exist per call; together with the per-`line` closure above, the lambda arguments are data-dependent and cannot be pre-built.

Checklist evidence: rule condition "lambda constructed per homogeneous loop element" — closures depend on loop variables (`line`, `word`); no fixed-schema path exists.

### [ ] Finding 43 — PERF-PY-25

- Function context: `scripts/pycaps/findings/functions/43.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/video/subtitle_clips_generator.py:45:1`
- Checklist pattern: same per-iteration closure as Finding 39.

Source excerpt:

```
                    self.__generate_word_clips_for_line(
                        line,
                        ElementState.LINE_BEING_NARRATED,
                        ElementState.WORD_BEING_NARRATED,
                        lambda word: word.time.start,
                        lambda word: word.time.end,
                        pbar
                    )
```

Why this is a false positive: both lambdas select per-word time values and are passed as callbacks; their captured values are only available per element, so per-element construction is semantically required.

Checklist evidence: rule condition "lambda constructed per homogeneous loop element" — the closures are per-word accessors with loop-dependent state; no hoisting possible.

### [ ] Finding 44 — PERF-PY-25

- Function context: `scripts/pycaps/findings/functions/44.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/video/subtitle_clips_generator.py:46:1`
- Checklist pattern: same per-iteration closure as Finding 39.

Source excerpt:

```
                    self.__generate_word_clips_for_line(
                        line,
                        ElementState.LINE_BEING_NARRATED,
                        ElementState.WORD_BEING_NARRATED,
                        lambda word: word.time.start,
                        lambda word: word.time.end,
                        pbar
                    )
```

Why this is a false positive: `lambda word: word.time.end` is a per-word accessor required at call time; the allocation is necessary and has no hoistable alternative.

Checklist evidence: rule condition "lambda constructed per homogeneous loop element" — the lambda is a data-dependent callback, not a heavy object with a fixed-schema path.

### [ ] Finding 45 — PERF-PY-25

- Function context: `scripts/pycaps/findings/functions/45.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/video/subtitle_clips_generator.py:54:1`
- Checklist pattern: same per-iteration closure as Finding 39.

Source excerpt:

```
                    self.__generate_word_clips_for_line(
                        line,
                        ElementState.LINE_BEING_NARRATED,
                        ElementState.WORD_ALREADY_NARRATED,
                        lambda word: word.time.end,
                        lambda _: line.time.end,
                        pbar
                    )
```

Why this is a false positive: `lambda word: word.time.end` depends on the per-word iteration value and `lambda _: line.time.end` on the per-line value; neither can be hoisted, so the construction is unavoidable.

Checklist evidence: rule condition "lambda constructed per homogeneous loop element" — loop-variable-dependent closures; no avoidable per-element work.

### [ ] Finding 46 — PERF-PY-25

- Function context: `scripts/pycaps/findings/functions/46.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/video/subtitle_clips_generator.py:55:1`
- Checklist pattern: same per-iteration closure as Finding 39.

Source excerpt:

```
                    self.__generate_word_clips_for_line(
                        line,
                        ElementState.LINE_BEING_NARRATED,
                        ElementState.WORD_ALREADY_NARRATED,
                        lambda word: word.time.end,
                        lambda _: line.time.end,
                        pbar
                    )
```

Why this is a false positive: `lambda _: line.time.end` captures the current `line`; since `line` changes each iteration, pre-construction outside the loop is impossible and the allocation is required.

Checklist evidence: rule condition "lambda constructed per homogeneous loop element" — closure over a loop variable; no fixed-schema path exists.

### [ ] Finding 47 — PERF-PY-25

- Function context: `scripts/pycaps/findings/functions/47.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/video/subtitle_clips_generator.py:63:1`
- Checklist pattern: same per-iteration closure as Finding 39.

Source excerpt:

```
                    self.__generate_word_clips_for_line(
                        line,
                        ElementState.LINE_ALREADY_NARRATED,
                        ElementState.WORD_ALREADY_NARRATED,
                        lambda _: line.time.end,
                        lambda _: segment.time.end,
                        pbar
                    )
```

Why this is a false positive: `lambda _: segment.time.end` captures the current `segment` from the outer loop; the value differs per iteration, so the lambda cannot be hoisted and the construction is necessary.

Checklist evidence: rule condition "lambda constructed per homogeneous loop element" — data-dependent closure over `segment`; no avoidable allocation.

### [ ] Finding 48 — PERF-PY-25

- Function context: `scripts/pycaps/findings/functions/48.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/video/subtitle_clips_generator.py:64:1`
- Checklist pattern: same per-iteration closure as Finding 39.

Source excerpt:

```
                    self.__generate_word_clips_for_line(
                        line,
                        ElementState.LINE_ALREADY_NARRATED,
                        ElementState.WORD_ALREADY_NARRATED,
                        lambda _: line.time.end,
                        lambda _: segment.time.end,
                        pbar
                    )
```

Why this is a false positive: `lambda _: segment.time.end` must be constructed per segment/line because its captured value changes every iteration; no hoistable "fixed-schema" alternative exists.

Checklist evidence: rule condition "lambda constructed per homogeneous loop element" — the closure depends on loop state, so per-element construction is semantically required, not avoidable work.

### [ ] Finding 51 — CWE-367

- Function context: `scripts/pycaps/findings/functions/51.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/tests/test_transcript_loader.py:69:33`
- Checklist pattern: the exists-check is a teardown guard inside a single-threaded test; the race has no security consequence.

Source excerpt:

```
        srt = "1\n00:00:01,000 --> 00:00:02,000\nHello\n"
        path = self._write_temp_file(".srt", srt)
        self.addCleanup(lambda: os.path.exists(path) and os.remove(path))
```

Why this is a false positive: in a single-threaded unit test teardown the check only prevents `FileNotFoundError` when the file was already removed by the test; if the file vanished between check and remove, the worst outcome is that same benign error in cleanup — no privilege boundary or attacker-controlled state exists.

Checklist evidence: rule condition "path checked before later separate use" is matched, but the guard's only purpose is avoiding a benign error in test teardown; no TOCTOU weakness with security impact is present.

## Uncertain findings

None. Finding 22 (CWE-22) reclassified as a true positive: `open(os.path.join(self._base_path, rule.filename), ...)` matches the rule condition (joined path segment reaches `open` without basename/resolve confinement).


## True positives

### BP-PY-1 — Bare Except Clause

| Finding | Source | Reason |
| --- | --- | --- |
| 4 | `src/pycaps/api/emoji_in_segments_api.py:26` | `except Exception as e:` logs and swallows; no re-raise |
| 6 | `src/pycaps/api/pycaps_tagger_api.py:22` | `except Exception as e:` logs and swallows; no re-raise |
| 11 | `src/pycaps/effect/clip/animate_segment_emojis_effect.py:76` | `except Exception as e:` logs and disables feature; no re-raise |
| 26 | `src/pycaps/renderer/pictex_subtitle_renderer.py:76` | bare `except:` swallows all exceptions |
| 27 | `src/pycaps/renderer/pictex_subtitle_renderer.py:106` | bare `except:` swallows all exceptions |
| 32 | `src/pycaps/template/template_service.py:20` | bare `except:` returns False |
| 37 | `src/pycaps/video/audio_utils.py:28` | `except Exception as e:` logs only |
| 50 | `src/pycaps/video/video_generator.py:127` | `except Exception as e:` logs only |

### CWE-396 — Declaration of Catch for Generic Exception

| Finding | Source | Reason |
| --- | --- | --- |
| 2 | `src/pycaps/ai/gpt.py:31` | `except Exception as e:` wraps into RuntimeError |
| 5 | `src/pycaps/api/emoji_in_segments_api.py:26` | `except Exception as e:` hides distinct failures |
| 7 | `src/pycaps/api/pycaps_tagger_api.py:22` | `except Exception as e:` hides distinct failures |
| 12 | `src/pycaps/effect/clip/animate_segment_emojis_effect.py:76` | `except Exception as e:` hides distinct failures |
| 14 | `src/pycaps/pipeline/caps_pipeline.py:224` | `except Exception as e:` logs and re-raises raw |
| 24 | `src/pycaps/renderer/css_subtitle_renderer.py:73` | `except Exception as e:` wraps into RuntimeError |
| 34 | `src/pycaps/transcriber/google_audio_transcriber.py:99` | `except Exception as e:` wraps into RuntimeError |
| 36 | `src/pycaps/transcriber/whisper_audio_transcriber.py:84` | `except Exception as e:` wraps into RuntimeError |
| 38 | `src/pycaps/video/audio_utils.py:28` | `except Exception as e:` hides distinct failures |
| 49 | `src/pycaps/video/video_generator.py:77` | `except Exception as e:` logs only |

### BP-PY-14 — requests Without Timeout

| Finding | Source | Reason |
| --- | --- | --- |
| 3 | `src/pycaps/api/api_sender.py:31` | `requests.post(_PYCAPS_API_URL, json=body)` has no `timeout=` |
| 9 | `src/pycaps/effect/clip/animate_segment_emojis_effect.py:49` | `requests.get(self.ASSETS_ZIP_URL, stream=True)` has no `timeout=` |

### BP-PY-7 — open Without Context Manager

| Finding | Source | Reason |
| --- | --- | --- |
| 8 | `src/pycaps/cli/preview_styles_cli.py:22` | builtin `open(css, "r", ...).read()` without `with` |
| 17 | `src/pycaps/pipeline/caps_pipeline_builder.py:51` | builtin `open(css_file_path, "r", ...).read()` without `with` |
| 21 | `src/pycaps/pipeline/json_config_loader.py:227` | builtin `open(os.path.join(...), "r", ...).read()` without `with` |
| 33 | `src/pycaps/transcriber/editor/transcription_editor.py:37` | builtin `open(html_file_path, 'r', ...).read()` without `with` |

### BP-PY-5 — Wildcard Import

| Finding | Source | Reason |
| --- | --- | --- |
| 18 | `src/pycaps/pipeline/json_config_loader.py:6` | `from pycaps.effect import *` |
| 19 | `src/pycaps/pipeline/json_config_loader.py:7` | `from pycaps.animation import *` |

### BP-PY-4 — Mutable Default Argument

| Finding | Source | Reason |
| --- | --- | --- |
| 29 | `src/pycaps/renderer/renderer_page.py:9` | `get_html` defaults `segment_tags=[], line_tags=[], words=[], word_tags=[], word_states=[]` are shared mutable lists |


### CWE-22 — Path Traversal (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 22 | `src/pycaps/pipeline/json_config_loader.py:227` | `open(os.path.join(self._base_path, rule.filename), ...)` joins a config-controlled filename onto a base path without basename/resolve confinement — rule condition met (reclassified from Uncertain) |

## Final evidence

- Delegated reviewers: none (single-reviewer audit)
- Chunk evidence: `scripts/pycaps/chunks`
- Function evidence: `scripts/pycaps/findings/functions`
- Validation: `git diff --check` — pass

## Post-fix remaining-FP audit (2026-08-02)

Mode A — remaining false positives. Fresh scan run with the FP-reduction binary (rebuilt 2026-08-02 16:29, fix commit `b5b8fde`).

### Run metadata (fresh scan)

```yaml
timestamp: 2026-08-02
repository: pycaps
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps
branch: main
commit: 68a6843bbd3f7d0ea33047ffb1045d978a7671c4
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps
chunk_path: scripts/pycaps/chunks
function_context_path: scripts/pycaps/findings/functions
```

### Scan evidence (fresh scan)

- Build command: `(pre-built) ./bin/goslop`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/pycaps/chunks -context-dir scripts/pycaps/findings/functions real-repos/pycaps`
- Findings: `29` (down from 51 pre-fix)
- Chunks reviewed: `scripts/pycaps/chunks/Chunk_1_25.txt`, `scripts/pycaps/chunks/Chunk_26_29.txt`
- Function contexts reviewed: `scripts/pycaps/findings/functions/1.txt … 29.txt`

### Audit checklist

- [x] Read every assigned chunk under `scripts/pycaps/chunks`.
- [x] Read `scripts/pycaps/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

### Classification summary (fresh run)

Fresh IDs were matched to old IDs by `Source:` (file:line:col). Fresh findings whose source matches an audited TP are true positives; sources matching audited FPs are re-appearing false positives; unmatched sources were classified on the rule condition.

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 6 | 1, 9, 12, 16, 22, 25 |
| True positive | 23 | 2, 3, 4, 5, 6, 7, 8, 10, 11, 13, 14, 15, 17, 18, 19, 20, 21, 23, 24, 26, 27, 28, 29 |
| Uncertain | 0 | — |

Remaining false positives: 6 (fresh IDs 1, 9, 12, 16, 22, 25), corresponding to old IDs 1, 10, 16, 20, 31, 35. All other 17 pre-fix FPs (old IDs 13, 15, 23, 25, 28, 30, 39–48, 51) no longer appear in the fresh scan — the fix suppressed them.

## False positives (fresh run)

### [ ] Finding 1 — BP-PY-13

- Function context: `scripts/pycaps/findings/functions/1.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/ai/gpt.py:6:1`
- Checklist pattern: the flagged string is an environment-variable *name* used as the lookup key for `os.getenv`, not a credential value (matches audited FP 1).

Source excerpt:

```
    OPENAI_API_KEY_NAME = "PYCAPS_OPENAI_API_KEY"
    ...
        return os.getenv(self.OPENAI_API_KEY_NAME) is not None
    ...
            self._client = OpenAI(api_key=os.getenv(self.OPENAI_API_KEY_NAME))
```

Why this is a false positive: the literal is the name of the environment variable the code reads the key from via `os.getenv(self.OPENAI_API_KEY_NAME)` (gpt.py:15, 24) — the exact fix the rule recommends; no credential value is hardcoded.

Checklist evidence: the rule condition "secret-like name assigned a string literal" fires only on the identifier pattern; the value is a lookup key for `os.getenv`, not a loaded-from-environment-able secret.

### [ ] Finding 9 — CWE-409

- Function context: `scripts/pycaps/findings/functions/9.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/effect/clip/animate_segment_emojis_effect.py:71:18`
- Checklist pattern: the archive is fetched from a hardcoded, version-pinned URL constant of the project's own release, so no untrusted party supplies the compressed data (matches audited FP 10).

Source excerpt:

```
    ASSETS_ZIP_URL = "https://github.com/francozanardi/pycaps/releases/download/emoji-assets-v2/animated_emojis.zip"
    ...
            file_buffer.seek(0)
            with zipfile.ZipFile(file_buffer) as z:
                z.extractall(self.CACHE_DIR)
```

Why this is a false positive: the detector reports any `.extractall` call, but the archive bytes come solely from `requests.get(self.ASSETS_ZIP_URL, ...)` where the URL is a compile-time constant pinned to the project's own HTTPS release asset; no user-influenced input reaches the extraction sink.

Checklist evidence: the zip-bomb trust boundary (untrusted compressed data) is absent — the archive source is a source constant, not user input.

### [ ] Finding 12 — CWE-367

- Function context: `scripts/pycaps/findings/functions/12.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/pipeline/caps_pipeline_builder.py:49:16`
- Checklist pattern: the `exists` check only feeds a user-facing `ValueError`; the subsequent use fails benignly (matches audited FP 16).

Source excerpt:

```
    def add_css(self, css_file_path: str) -> "CapsPipelineBuilder":
        if not os.path.exists(css_file_path):
            raise ValueError(f"CSS file not found: {css_file_path}")
        css_content = open(css_file_path, "r", encoding="utf-8").read()
```

Why this is a false positive: the check-then-use race has no privilege or integrity consequence — if the file vanishes between check and `open`, `open()` raises `FileNotFoundError` with the same benign outcome; no attacker-controlled state is used and no security boundary is crossed.

Checklist evidence: rule condition "path checked before later separate use" is matched literally, but the check exists only to produce a friendly error and the race's failure mode is a benign exception on the user's own file.

### [ ] Finding 16 — CWE-478

- Function context: `scripts/pycaps/findings/functions/16.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/pipeline/json_config_loader.py:93:1`
- Checklist pattern: the match domain is a pydantic-validated `Literal` union whose every member has a case branch (matches audited FP 20).

Source excerpt:

```
        for effect in self._config.effects:
            match effect.type:
                case "emoji_in_segment":
                case "emoji_in_word":
                case "remove_punctuation_marks":
                case "typewriting":
                case "animate_segment_emojis":
```

Why this is a false positive: `effect.type` is a pydantic `Literal["emoji_in_segment", "emoji_in_word", "remove_punctuation_marks", "typewriting", "animate_segment_emojis"]` (`json_schema.py:51-74`, verified unchanged), and all five members are handled by the five case branches; any other value raises `ValidationError` before the match is reached, so no unhandled value can occur.

Checklist evidence: rule condition "match with two or more cases lacks a wildcard" is met syntactically, but the matched value is exhaustively constrained by schema validation, making a default branch unreachable dead code.

### [ ] Finding 22 — CWE-478

- Function context: `scripts/pycaps/findings/functions/22.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/selector/time_event_selector.py:28:1`
- Checklist pattern: the matched value is a Python `Enum` whose members are all covered by the case branches (matches audited FP 31).

Source excerpt:

```
        match self._element_type:
            case ElementType.WORD:
                return self.__filter_by_words(clips)
            case ElementType.LINE:
                return self.__filter_by_lines(clips)
            case ElementType.SEGMENT:
                return self.__filter_by_segments(clips)
```

Why this is a false positive: `ElementType` is a `str, Enum` with exactly `WORD`, `LINE`, `SEGMENT` members (`common/types.py:21-24`, verified unchanged), all three handled; invalid values cannot construct the enum, so no value reaches the match unhandled.

Checklist evidence: rule condition "match with multiple cases and no wildcard" is met syntactically, but the enum domain is fully covered — a default branch would be unreachable.

### [ ] Finding 25 — CWE-186

- Function context: `scripts/pycaps/findings/functions/25.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pycaps/src/pycaps/transcriber/transcript_loader.py:14:23`
- Checklist pattern: the regex parses a fixed standard format (VTT inline timestamp), not a validator that rejects valid identifiers (matches audited FP 35).

Source excerpt:

```
_TIME_EPSILON = 0.01
_VTT_INLINE_TIME_RE = re.compile(r"<((?:\d{2}:)?\d{2}:\d{2}\.\d{3})>")
```

Why this is a false positive: the detector fires on the `\d{3}` quantifier, but this regex extracts WebVTT inline timestamps, whose format is exactly `[HH:]MM:SS.mmm` per the VTT spec; the millisecond part is always three digits, so the pattern is precise, not overly restrictive — no valid input is rejected.

Checklist evidence: rule condition "validation regex accepts only a fixed narrow identifier shape" — this is a parser for a standards-defined fixed format (used to extract, not validate user identifiers), so narrowness is inherent to the format, not an over-restriction.

## Uncertain findings (fresh run)

None.

## True positives (fresh run)

Fresh findings whose `Source:` matches an audited TP (rule, source, old ID):

| Fresh finding | Rule | Source (matches old) | Old finding |
| --- | --- | --- | --- |
| 2 | BP-PY-14 | `src/pycaps/api/api_sender.py:31` | 3 |
| 3 | BP-PY-1 | `src/pycaps/api/emoji_in_segments_api.py:26` | 4 |
| 4 | CWE-396 | `src/pycaps/api/emoji_in_segments_api.py:26` | 5 |
| 5 | BP-PY-1 | `src/pycaps/api/pycaps_tagger_api.py:22` | 6 |
| 6 | CWE-396 | `src/pycaps/api/pycaps_tagger_api.py:22` | 7 |
| 7 | BP-PY-7 | `src/pycaps/cli/preview_styles_cli.py:22` | 8 |
| 8 | BP-PY-14 | `src/pycaps/effect/clip/animate_segment_emojis_effect.py:49` | 9 |
| 10 | BP-PY-1 | `src/pycaps/effect/clip/animate_segment_emojis_effect.py:76` | 11 |
| 11 | CWE-396 | `src/pycaps/effect/clip/animate_segment_emojis_effect.py:76` | 12 |
| 13 | BP-PY-7 | `src/pycaps/pipeline/caps_pipeline_builder.py:51` | 17 |
| 14 | BP-PY-5 | `src/pycaps/pipeline/json_config_loader.py:6` | 18 |
| 15 | BP-PY-5 | `src/pycaps/pipeline/json_config_loader.py:7` | 19 |
| 17 | BP-PY-7 | `src/pycaps/pipeline/json_config_loader.py:227` | 21 |
| 18 | CWE-22 | `src/pycaps/pipeline/json_config_loader.py:227` | 22 |
| 19 | BP-PY-1 | `src/pycaps/renderer/pictex_subtitle_renderer.py:76` | 26 |
| 20 | BP-PY-1 | `src/pycaps/renderer/pictex_subtitle_renderer.py:106` | 27 |
| 21 | BP-PY-4 | `src/pycaps/renderer/renderer_page.py:9` | 29 |
| 23 | BP-PY-1 | `src/pycaps/template/template_service.py:20` | 32 |
| 24 | BP-PY-7 | `src/pycaps/transcriber/editor/transcription_editor.py:37` | 33 |
| 26 | BP-PY-1 | `src/pycaps/video/audio_utils.py:28` | 37 |
| 27 | CWE-396 | `src/pycaps/video/audio_utils.py:28` | 38 |
| 28 | BP-PY-1 | `src/pycaps/video/video_generator.py:127` | 50 |

Finding 29 (CWE-396, `src/pycaps/video/video_generator.py:127:1`) is a new true positive: `except Exception as e:` at video_generator.py:127 logs a warning and does not re-raise, so the generic handler hides the distinct failure conditions of `os.remove` — the rule condition is satisfied. It does not match any audited source (the audited CWE-396 at this file was line 77, old finding 49), so it was classified on the rule condition.

## Final evidence (fresh run)

- Delegated reviewers: none (single-reviewer audit)
- Chunk evidence: `scripts/pycaps/chunks`
- Function evidence: `scripts/pycaps/findings/functions`
- Validation: `git diff --check` — pass
