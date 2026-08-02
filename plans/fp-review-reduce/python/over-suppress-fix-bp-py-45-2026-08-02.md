# Over-suppression fix — BP-PY-45 `__file__` / guarded bootstrap (2026-08-02)

## Finding IDs / rule

| Repo | Old finding IDs (representative) | Rule |
| --- | --- | --- |
| WeThePeople | 135, 173, 179, … (64 module-level); 502, 861 (guarded); 1343 (in-function) | BP-PY-45 |
| html2pic | 1–5 (`examples/*` `__file__` inserts) | BP-PY-45 |
| violit | Mode-A notes: old 25–27, 130 (`sys.path.insert` bootstrap in examples/cli) suppressed-but-present | BP-PY-45 |
| httpmorph | 1, 224, 347 (benchmarks / docs conf / examples) — docs conf stays silent; others must fire | BP-PY-45 |

Checklist sources:
- `plans/fp-validations/reports/WeThePeople.md` (Post-fix over-suppression, 67 BP-PY-45)
- `plans/fp-validations/reports/html2pic.md` (5 BP-PY-45)
- `plans/fp-validations/reports/violit.md` Mode-A notes
- Prior reduction: `plans/fp-review-reduce/python/batch-2-bp-py-41-45-reduction-2026-08-02.md`

## Root cause

Batch-2 `isSysPathBootstrap` treated **any** `sys.path.insert/append/extend` line containing
`__file__` as a safe import-time bootstrap. That silenced real TPs:

1. WeThePeople `jobs/*` / `scripts/*` module-level
   `sys.path.insert(0, os.path.dirname(...abspath(__file__)))`
2. Module-level `if … not in sys.path:` guarded inserts that still mutate path on import
3. In-function inserts that happen to mention `__file__` (despite the comment claiming
   in-function stays reportable)

The only validated FP cut from that batch is Sphinx `docs/**/conf.py`, already handled by
`isSphinxDocsConfPath`.

## Detector change

File: `internal/lang/python/detectors/bad_practices/rules_deps.go`

- Removed blanket `strings.Contains(t, "__file__")` short-circuit
- Removed module-level `not in sys.path` guarded-bootstrap skip (WeThePeople uses that as TP)
- Deleted `isSysPathBootstrap` and the scope stack that existed only to support it
- **Kept** `isSphinxDocsConfPath` / `docs/**/conf.py` path skip
- **Kept** test-file, requirements, and `sitecustomize` / `usercustomize` / `conftest` skips

## Fixtures

| Path | Role |
| --- | --- |
| `tests/fixtures/python/bp/BP-PY-45-file-bootstrap-vulnerable.txt` | WeThePeople jobs-style `__file__` module insert — **must fire** |
| `tests/fixtures/python/bp/BP-PY-45-file-bootstrap-safe.txt` | Sphinx `docs/source/conf.py` — **silent** (path exemption) |
| `tests/fixtures/python/bp/BP-PY-45-guarded-bootstrap-vulnerable.txt` | In-function `__file__` insert (finding 1343 shape) — **must fire** |
| `tests/fixtures/python/bp/BP-PY-45-guarded-bootstrap-safe.txt` | Packaged import, no path mutation — **silent** |
| `tests/fixtures/python/bp/BP-PY-45-guarded-module-vulnerable.txt` | WeThePeople guarded `not in sys.path` module bootstrap — **must fire** |
| `tests/fixtures/python/bp/BP-PY-45-guarded-module-safe.txt` | Same guarded shape under `docs/conf.py` — **silent** |
| `tests/fixtures/python/bp/BP-PY-45-docs-conf-{safe,vulnerable}.txt` | Unchanged Sphinx vs library-vendor pair |

## Tests

Updated `rules_deps_test.go` and `audit_variants_test.go` to load the new
`BP-PY-45-guarded-module` pair (fixtures only; no inline Python).

```text
go test ./internal/lang/python/detectors/bad_practices/ -count=1 \
  -run 'TestBPPY45SysPathMutation|TestBPFalsePositiveAuditFixtureVariants/BP-PY-45'
ok  	github.com/chinmay-sawant/goslop/internal/lang/python/detectors/bad_practices
```

All `BP-PY-45-*` audit subtests PASS (docs-conf, file-bootstrap, guarded-bootstrap, guarded-module).

## Residual risks

- Project_Parva-style module `__file__` bootstraps that batch-2 treated as FPs will fire again.
  Per WeThePeople / html2pic / violit Mode-A evidence those are TPs for this rule.
- Sphinx exemption is path-based (`docs/**/conf.py` only). A Sphinx conf outside `docs/`
  still reports.
- Corpus rescan of `real-repos/*` not run (focused fixture tests only).
- BP-PY-46 / CWE-396 / rich-print intentionally untouched.
