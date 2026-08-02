# False-positive audit: html2pic

## Run metadata

```yaml
timestamp: 2026-08-02T07:23:48Z
repository: html2pic
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/html2pic
branch: main
commit: 27c292f1f2afd8975b9cb58470b9e7469df52dad
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/html2pic
chunk_path: ./scripts/chunks
function_context_path: ./scripts/findings/functions
```

## Scan evidence

- Build command: `(pre-built) ./bin/goslop`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir real-repos/html2pic/scripts/chunks -context-dir real-repos/html2pic/scripts/findings/functions real-repos/html2pic`
- Findings: `36`
- Chunks reviewed: `./scripts/chunks/Chunk_1_25.txt`, `./scripts/chunks/Chunk_26_36.txt`
- Function contexts reviewed: `./scripts/findings/functions/1.txt` .. `./scripts/findings/functions/36.txt` (all 36)

## Audit checklist

- [x] Read every assigned chunk under `./scripts/chunks`.
- [x] Read `./scripts/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [ ] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 4 | 6, 7, 10, 32 |
| True positive | 32 | 1, 2, 3, 4, 5, 8, 9, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 33, 34, 35, 36 |
| Uncertain | 0 | — |

## False positives

All four false positives are CWE-396 findings whose handler propagates the
exception to the caller (wrapped as a domain error with the original chained
via `from e`) instead of consuming it, so no distinct failure condition is
hidden — the weakness the rule describes does not occur.

### [ ] Finding `6` — CWE-396

- Function context: `./scripts/findings/functions/6.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/html2pic/src/html2pic/html2pic.py:87:1`
- Checklist pattern: generic handler suite propagates the exception (wrapped, `from e`), so no failure condition is hidden

Source excerpt:

```
    def render(self, crop_mode: pictex.CropMode = pictex.CropMode.CONTENT_BOX) -> pictex.BitmapImage:
        """Render to a bitmap image."""
        try:
            canvas, root = self._translator.translate(self.styled_tree, self.font_registry)
            ...
        except Exception as e:
            self._print_warnings()
            raise RenderError(f"Failed to render: {e}") from e
```

Why this is a false positive: The handler only prints accumulated warnings and then re-raises the failure as `RenderError` with the original exception chained via `from e`; the failure conditions reach the caller and are not hidden, so CWE-396's condition (a generic handler that hides distinct failure conditions) is not satisfied.

Checklist evidence: The rule's message is "generic Exception or BaseException handler can hide distinct failure conditions"; the shown source's handler propagates the exception (`raise RenderError(...) from e`), so the "hides failures" condition is unmet.

### [ ] Finding `7` — CWE-396

- Function context: `./scripts/findings/functions/7.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/html2pic/src/html2pic/parsing/css_parser.py:56:1`
- Checklist pattern: generic handler suite re-raises as a domain error with the original exception chained, so no failure condition is hidden

Source excerpt:

```
        try:
            self.font_registry.clear()
            stylesheet = tinycss2.parse_stylesheet(css_content)
            ...
            return rules, self.font_registry
        except Exception as e:
            raise ParseError(f"Failed to parse CSS: {e}") from e
```

Why this is a false positive: The handler immediately wraps and re-raises (`raise ParseError(...) from e`); the original exception remains attached as `__cause__`, so the failure is propagated to the caller, not hidden by the generic catch.

Checklist evidence: CWE-396's condition is that the generic handler "can hide distinct failure conditions"; here the suite re-raises with the original exception chained, so the condition is not satisfied.

### [ ] Finding `10` — CWE-396

- Function context: `./scripts/findings/functions/10.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/html2pic/src/html2pic/parsing/html_parser.py:24:1`
- Checklist pattern: generic handler suite re-raises as a domain error with the original exception chained, so no failure condition is hidden

Source excerpt:

```
        try:
            soup = BeautifulSoup(html_content, 'html.parser')
            return self._create_root(soup)
        except Exception as e:
            raise ParseError(f"Failed to parse HTML: {e}") from e
```

Why this is a false positive: The handler wraps and re-raises the exception (`raise ParseError(...) from e`); the failure propagates to the caller with the original exception chained, so nothing is hidden.

Checklist evidence: Same as finding 7 — the suite re-raises with `from e`, so the "hides distinct failure conditions" condition of CWE-396 is not satisfied.

### [ ] Finding `32` — CWE-396

- Function context: `./scripts/findings/functions/32.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/html2pic/src/html2pic/translation/translator.py:49:1`
- Checklist pattern: generic handler suite re-raises as a domain error with the original exception chained, so no failure condition is hidden

Source excerpt:

```
        try:
            canvas = Canvas()
            root_element = self._translate_node(styled_dom)
            return canvas, root_element
        except Exception as e:
            raise RenderError(f"Translation failed: {e}") from e
```

Why this is a false positive: The handler wraps and re-raises (`raise RenderError(...) from e`); the original exception is chained and the failure reaches the caller, so the generic catch does not hide distinct failure conditions.

Checklist evidence: Same as findings 6/7/10 — the suite re-raises with `from e`, so CWE-396's "hides distinct failure conditions" condition is not satisfied.

## True positives

### BP-PY-45 — sys.path Mutation At Runtime (5 findings)

The rule flags any `sys.path.insert(`/`append(`/`extend(` line except in test files, requirements paths, and the bootstrap basenames `sitecustomize.py`, `usercustomize.py`, `conftest.py`. Each example script performs a literal `sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'src'))` at module scope to fix imports; examples are not on the exemption list.

| Finding id | Source | Reason |
| --- | --- | --- |
| 1 | `examples/01_quick_start.py:12:1` | `sys.path.insert(0, ...)` at module scope; not a test file / requirements path / bootstrap basename, no exemption applies. |
| 2 | `examples/02_flexbox_card.py:12:1` | Same literal construct as finding 1, distinct file. |
| 3 | `examples/03_typography_showcase.py:12:1` | Same literal construct, distinct file. |
| 4 | `examples/04_shadows_and_effects.py:12:1` | Same literal construct, distinct file. |
| 5 | `examples/05_background_images.py:12:1` | Same literal construct, distinct file. |

### BP-PY-1 — Bare Except Clause (11 findings)

The rule flags a bare `except:` unconditionally, and `except Exception`/`except BaseException` unless the suite re-raises (`suiteReraises`) or a test file collects evidence. Every flagged handler here reports the stringified error via `warn_unexpected_error(...)` and consumes the exception — the suite never re-raises, so no built-in exemption applies.

| Finding id | Source | Reason |
| --- | --- | --- |
| 8 | `src/html2pic/parsing/css_parser.py:128:1` | `except Exception as e:` in `_calculate_specificity`; suite only warns and `return 1`, no re-raise. |
| 9 | `src/html2pic/parsing/css_parser.py:166:1` | `except Exception as e:` in `_process_font_face`; suite only warns, no re-raise. |
| 12 | `src/html2pic/styling/cascade_resolver.py:35:1` | `except Exception as e:`; suite warns and `return 1`, no re-raise. |
| 14 | `src/html2pic/styling/cascade_resolver.py:55:1` | `except Exception as e:`; suite warns and `return False`, no re-raise. |
| 16 | `src/html2pic/translation/shadow_parser.py:83:1` | `except Exception as e:`; suite warns and `return None`, no re-raise. |
| 18 | `src/html2pic/translation/style_applicators/background_applicator.py:28:1` | `except Exception as e:`; suite warns and continues, no re-raise. |
| 20 | `src/html2pic/translation/style_applicators/background_applicator.py:88:1` | `except Exception as e:` around color-stop parse; suite warns and continues, no re-raise. |
| 21 | `src/html2pic/translation/style_applicators/background_applicator.py:93:1` | `except Exception as e:` around gradient parse; suite warns and `return None`, no re-raise. |
| 22 | `src/html2pic/translation/style_applicators/background_applicator.py:110:1` | Bare `except:` (no type) around angle conversion; bare except is flagged unconditionally — catches `BaseException` including `KeyboardInterrupt`/`SystemExit`. |
| 24 | `src/html2pic/translation/style_applicators/border_applicator.py:38:1` | `except Exception as e:`; suite warns and continues, no re-raise. |
| 30 | `src/html2pic/translation/style_applicators/typography_applicator.py:83:1` | `except Exception as e:`; suite warns and continues, no re-raise. |

### BP-PY-46 — print Debugging In Library Code (8 findings)

The rule flags `print(` in executable code outside a `__main__` guard, CLI-decorated functions, or argparse CLI entrypoint functions. Neither `transform_applicator.py` nor `warning_collector.py` contains a `__main__` guard, `argparse`, or `main()`; the `print_summary` method name (`print_` prefix) is only exempt when `hasArgparse && mainInvoked`, which is false here.

| Finding id | Source | Reason |
| --- | --- | --- |
| 26 | `src/html2pic/translation/style_applicators/transform_applicator.py:29:9` | `print(f"element: {builder}")` in library `apply()`; leftover debug output, no guard. |
| 27 | `src/html2pic/translation/style_applicators/transform_applicator.py:30:9` | `print(f"children: {builder._children}")` — same debug block, distinct line. |
| 28 | `src/html2pic/translation/style_applicators/transform_applicator.py:31:9` | `print(f"translate_x: {translate_x}")` — same debug block, distinct line. |
| 29 | `src/html2pic/translation/style_applicators/transform_applicator.py:32:9` | `print(f"translate_y: {translate_y}")` — same debug block, distinct line. |
| 33 | `src/html2pic/warnings/warning_collector.py:142:9` | `print(f"\nCompleted with {summary['total_warnings']} warnings:")` in library method `print_summary()`. |
| 34 | `src/html2pic/warnings/warning_collector.py:145:13` | `print(f"  {category.replace('_', ' ').title()}: {count}")` in the same method, distinct line. |
| 35 | `src/html2pic/warnings/warning_collector.py:147:9` | `print("\nDetailed warnings:")` in the same method, distinct line. |
| 36 | `src/html2pic/warnings/warning_collector.py:149:13` | `print(f"  {i}. [{warning['category'].upper()}] {warning['message']}")` in the same method, distinct line. |

### CWE-1046 — Creation of Immutable Text Using String Concatenation (3 findings)

The rule fires when a loop body contains `name += rhs` where `name` was previously assigned a string literal (or the name/RHS look like text). In each flagged function, `current_part`/`current` is initialized to `''` and accumulated one character at a time inside `for char in ...:` — the exact accumulator pattern the rule targets.

| Finding id | Source | Reason |
| --- | --- | --- |
| 11 | `src/html2pic/parsing/shorthand_expander.py:64:1` | `current_part = ''` then `current_part += char` inside the `for char in value:` loop of `_split_preserving_functions`. |
| 15 | `src/html2pic/translation/shadow_parser.py:30:1` | `current = ''` then `current += char` inside the `for char in shadow_str:` loop of `_split_shadows`. |
| 23 | `src/html2pic/translation/style_applicators/background_applicator.py:121:1` | `current = ''` then `current += char` inside the `for char in content:` loop of `_split_gradient_parts`. |

### CWE-396 — Declaration of Catch for Generic Exception (5 findings)

These handlers consume the exception — recording only the stringified message and degrading to a default (warning + `return`/continue) — so distinct failure conditions are merged and never reach the caller; the broad catch genuinely hides what failed.

| Finding id | Source | Reason |
| --- | --- | --- |
| 13 | `src/html2pic/styling/cascade_resolver.py:35:1` | `except Exception as e:` swallows selector-parse failures, warns, and `return 1` — the caller cannot distinguish failure causes. |
| 17 | `src/html2pic/translation/shadow_parser.py:83:1` | `except Exception as e:` swallows shadow-parse failures, warns, and `return None`. |
| 19 | `src/html2pic/translation/style_applicators/background_applicator.py:28:1` | `except Exception as e:` swallows color-parse failures and continues with the property unapplied. |
| 25 | `src/html2pic/translation/style_applicators/border_applicator.py:38:1` | `except Exception as e:` swallows border-parse failures and continues with the property unapplied. |
| 31 | `src/html2pic/translation/style_applicators/typography_applicator.py:83:1` | `except Exception as e:` swallows color-parse failures and continues with the property unapplied. |

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none (single-reviewer audit)
- Chunk evidence: `./scripts/chunks`
- Function evidence: `./scripts/findings/functions`
- Validation: `git diff --check` — pass
