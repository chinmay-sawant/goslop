# Over-suppression fix — BP-PY-46 shebang entry-script gate (2026-08-02)

## Finding IDs / rule

- Rule: **BP-PY-46**
- Over-suppressed TPs (WHEN-Language Mode B): old findings **1, 2, 5, 7, 9, 10, 11, 12, 14, 15, 16** — all `print` operational logging in `hot_reload.py`
- Retained FP pattern: `when.py` CLI prints (shebang + `__main__`) must stay silent

## FP / over-suppression mechanism

Prior reduction used `isPythonShebangScript` as a whole-file early return in `detectBPPY46`. Any source starting with `#!` was skipped. `hot_reload.py` has a vestigial `#!/usr/bin/env python3` but is an importable library (no `__main__`, no CLI parser) — so 11 true positives were suppressed.

## Detector condition changed

| Before | After |
| --- | --- |
| Early-return when `isPythonShebangScript(src)` (`#!` alone) | Early-return when `pythonShebangIsEntryScript(src)` |

`pythonShebangIsEntryScript` (in `common.go`) requires:

1. shebang **AND**
2. (`__name__` + `__main__` present) **OR** argparse (`pythonHasArgparseCLI`) **OR** click/typer/cyclopts imports

Explicitly **does not** fold `rich.print` into this AND (rich remains a separate path via clickish helpers).

Call site: `rules_observability.go` `detectBPPY46` — only the shebang gate was touched.

## Fixtures

| Path | Expectation |
| --- | --- |
| `tests/fixtures/python/bp/BP-PY-46-shebang-script-safe.txt` | shebang + `__main__` + print → **silent** (when.py-shaped) |
| `tests/fixtures/python/bp/BP-PY-46-shebang-script-vulnerable.txt` | shebang + library print, no entry signals → **fires** |
| `tests/fixtures/python/bp/BP-PY-46-shebang-vestigial-library-vulnerable.txt` | shebang + `HotReloader`-shaped class prints → **fires** |
| `tests/fixtures/python/bp/BP-PY-46-shebang-vestigial-library-safe.txt` | shebang + argparse, no `__main__` → **silent** |

Wired in:

- `rules_observability_test.go` (`TestBPPY46PrintInLibrary`) — loads fixtures by path/name only
- `audit_variants_test.go` — `BP-PY-46-shebang-script` + `BP-PY-46-shebang-vestigial-library`
- Integration matrix auto-discovers the new pair stem

## Test results

```
go test ./internal/lang/python/detectors/bad_practices/ -count=1 \
  -run 'TestBPPY46PrintInLibrary|TestBPFalsePositiveAuditFixtureVariants/BP-PY-46'
→ PASS
```

Verbose shebang subcases (`BP-PY-46-shebang-script`, `BP-PY-46-shebang-vestigial-library`) also PASS.

## Residual risks

1. Shebang + CLI framework imports without any user-facing `print` path still whole-file skips (same as before for argparse/click modules) — acceptable for entry scripts.
2. A library that imports click/typer but only uses them as a dependency (not as an entry CLI) and has a vestigial shebang would be silenced — rare; not seen in WHEN-Language.
3. Concurrent agents may still edit `rules_observability.go` / path-exemption helpers; shebang gate is isolated to `pythonShebangIsEntryScript`.
4. Full WHEN-Language rescan not re-run in this pass; fixture coverage encodes the hot_reload vs when.py distinction.
