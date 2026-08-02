# Over-suppression fix — BP-PY-46 script-path dirs (2026-08-02)

Skill: `plans/skills/review-and-reduce/SKILLS.md`

Checklists:

- `plans/fp-validations/reports/FuncToWeb.md` (Post-fix over-suppression — 14× BP-PY-46 under `examples/`)
- `plans/fp-validations/reports/WeThePeople.md` (scripts/ BP-PY-46 cluster)
- Ledger context: `plans/fp-review-reduce/python/batch-1-bp-py-46-reduction-2026-08-02.md`

## Problem (validated)

`isPythonScriptModule` whole-file-skipped any path containing
`examples|example|scripts|script|tools|tool|demos|demo|commands`.

That silenced:

- **FuncToWeb** `examples/**` helper `print()` calls (14 audited TPs) — prints live in
  library-style functions, not under `if __name__ == "__main__":`; path skip was the sole cause.
- **WeThePeople** `scripts/**` prints (27 in the scripts/ cluster of the BP-PY-46 post-fix set).

httpmorph `examples/` mostly stay quiet via shebang + `__main__` / CLI without needing a path skip.

## Detector condition changed

In `internal/lang/python/detectors/bad_practices/common.go` → `scriptPathDirNames`:

| Action | Dir names |
| --- | --- |
| **REMOVED** | `examples`, `example`, `scripts`, `script`, `tools`, `tool`, `demos`, `demo` |
| **KEPT** | `commands` (Click/Typer subcommand packages; caniscrape relies on `commands/` + decorator) |
| **KEPT** (basenames in `isPythonScriptModule`) | `setup.py`, `__main__.py`, `cli.py` |

Shebang / `__main__` / CLI-decorator / rich-print guardrails (owned by parallel agents in
`rules_observability.go`) continue to silence real entry scripts.

## Fixtures

| Role | Path |
| --- | --- |
| Rewritten safe (was `examples/`-only silence) | `tests/fixtures/python/bp/BP-PY-46-script-path-safe.txt` — shebang entry + `__main__` (no path-dir skip) |
| Unchanged vulnerable (library) | `tests/fixtures/python/bp/BP-PY-46-script-path-vulnerable.txt` |
| **New** FuncToWeb-style vulnerable | `tests/fixtures/python/bp/BP-PY-46-examples-library-print-vulnerable.txt` — `examples/...` helper `print()` **must fire** |
| **New** safe sibling | `tests/fixtures/python/bp/BP-PY-46-examples-library-print-safe.txt` — same path, `print` under `__main__` (silent) |

Wired in:

- `rules_observability_test.go` (`TestBPPY46PrintInLibrary`)
- `audit_variants_test.go` (append-only `BP-PY-46-examples-library-print`)
- Integration matrix auto-discovers `*-{safe,vulnerable}.txt` pairs

## Tests

```text
go test ./internal/lang/python/detectors/bad_practices/ -count=1 \
  -run 'TestBPPY46PrintInLibrary|TestBPFalsePositiveAuditFixtureVariants/BP-PY-46'
→ PASS
```

`TestPythonBPFixturesMatrix`: no BP-PY-46 failures. Residual matrix failures are **BP-PY-45**
bootstrap safe fixtures (parallel agent), outside this change.

## Completion record

| Finding samples | Rule | Detector change | Safe fixture | Vulnerable fixture |
| --- | --- | --- | --- | --- |
| FuncToWeb 1,2,4,5,8–17 | BP-PY-46 | Drop examples/scripts/tools/demos from `scriptPathDirNames` | `BP-PY-46-examples-library-print-safe.txt` | `BP-PY-46-examples-library-print-vulnerable.txt` |
| WeThePeople scripts/ 1091+ | BP-PY-46 | Same (scripts/ no longer whole-file skip) | `BP-PY-46-script-path-safe.txt` (shebang) | `BP-PY-46-script-path-vulnerable.txt` |

- Checklist `[ ]` boxes left unchanged.
- No inline Python in Go tests.
- DO NOT COMMIT (per task).

## Residual risks

1. **True scripts under `scripts/` / `tools/` / `examples/` without shebang, `__main__`, CLI decorator, or rich-print** will now fire BP-PY-46. That restores the audited FuncToWeb/WeThePeople TPs; some genuine CLI noise may return in repos that relied only on the path skip (e.g. pdf_oxide `examples/` without shebang/main — accept and rely on other guardrails).
2. **`commands/` path skip remains** — dual-covered with Click/Typer decorators for caniscrape; a non-CLI library accidentally under `package/commands/` would still be silenced.
3. Parallel agents editing shebang / rich-print / BP-PY-45 in the same package; this change only touched `scriptPathDirNames` + BP-PY-46 fixtures/tests.
