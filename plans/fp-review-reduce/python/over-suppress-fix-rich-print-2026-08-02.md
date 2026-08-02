# Over-suppression fix — BP-PY-46 rich-print module skip (2026-08-02)

## Finding IDs / rule

| Repo | Old finding IDs | Rule | Sources |
| --- | --- | --- | --- |
| caniscrape | 290, 291 | BP-PY-46 | `caniscrape/telemetry.py:243`, `:245` (`contribute_scan`) |
| caniscrape | 318, 319, 320, 322, 323, 324, 325 | BP-PY-46 | `caniscrape/upload_handler.py:58–71` (`try_upload_scan`) |

Checklist: `plans/fp-validations/reports/caniscrape.md` Post-fix over-suppression
(9 BP-PY-46 rich-print cases) + batch-1 ledger
`plans/fp-review-reduce/python/batch-1-bp-py-46-reduction-2026-08-02.md`.

## Root cause

Batch-1 added a **module-wide early-return** in `detectBPPY46` when
`pythonUsesRichPrint(unit.Source)` (`from rich import print`). That silenced
every `print(` in `telemetry.py` / `upload_handler.py`, including operational
worker printers in `contribute_scan` / `try_upload_scan` (audited TPs), not only
CLI presentation helpers (`prompt_*`, `show_*`, `enable_*`, `disable_*`,
`display_*`).

## Detector change

File: `internal/lang/python/detectors/bad_practices/rules_observability.go`

1. **Removed** the `pythonUsesRichPrint` whole-module early-return from
   `detectBPPY46`.
2. **Kept** `pythonUsesRichPrint` as a signal inside `pythonHasClickishCLI`
   (with Click/Typer/Cyclopts / `rich.console`).
3. **Expanded** `pythonCLIPrintSkipFunc` when `hasClickishCLI`: skip presentation
   name prefixes `prompt_`, `show_`, `enable_`, `disable_`, `display_`.
   Explicitly does **not** skip `contribute_*`, `try_*`, `upload_*`, `worker_*`.

Shebang / script-path edits from parallel agents (`pythonShebangIsEntryScript`,
narrowed `scriptPathDirNames`) were left intact.

## Fixtures

| Path | Role |
| --- | --- |
| `tests/fixtures/python/bp/BP-PY-46-rich-print-safe.txt` | `from rich import print` + `prompt_*` / `show_*` / `enable_*` / `disable_*` / `display_*` — **silent** |
| `tests/fixtures/python/bp/BP-PY-46-rich-print-vulnerable.txt` | `from rich import print` + prints in `contribute_scan` / `try_upload_scan` — **must fire** |
| `tests/fixtures/python/bp/BP-PY-46-vulnerable.txt` | Plain no-rich library `print` (pre-existing) — **must fire** |

Wired in `rules_observability_test.go` (`TestBPPY46PrintInLibrary`) and
`audit_variants_test.go` (`BP-PY-46-rich-print`). No inline source in tests.

## Tests

```text
go test ./internal/lang/python/detectors/bad_practices/ -count=1 \
  -run 'TestBPPY46PrintInLibrary|TestBPFalsePositiveAuditFixtureVariants/BP-PY-46'
→ PASS
```

## Rescan (post `go build -o bin/goslop`, `--no-cache`)

```text
./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml \
  --no-cache --only BP-PY-46 \
  real-repos/caniscrape/caniscrape/telemetry.py \
  real-repos/caniscrape/caniscrape/upload_handler.py
```

| Check | Result |
| --- | --- |
| `contribute_scan` lines 243, 245 | **present** (restored TPs 290–291) |
| `try_upload_scan` lines 58, 59, 65, 66, 68, 69, 71 | **present** (restored TPs 318–325) |
| `prompt_*` / `show_*` / `enable_*` / `disable_*` presentation lines | **absent** |

## Residual risks

1. **`request_data_deletion` prints** in `telemetry.py` (~298–334) still fire —
   name is not under the presentation prefixes. Interactive CLI UX, not in the
   Mode-B over-suppressed TP list; narrowing further would need an explicit
   `request_*` (or similar) skip with care not to hide workers.
2. Any non-presentation function in a rich/clickish module that lacks a skipped
   name prefix will still be flagged (intentional for operational printers).
3. Parallel agents continue to edit shared `rules_observability.go` /
   `common.go`; re-verify shebang/path guards stay compatible with this change.
