# Over-suppression fix — CWE-396 suite-surface skip (2026-08-02)

## Finding IDs / rule

| Repo | Old finding IDs | Rule |
| --- | --- | --- |
| caniscrape | 54, 245, 329 | CWE-396 |
| voicetag | 30, 35–39, 42 | CWE-396 |
| among-llms | 13, 29 | CWE-396 |
| FuncToWeb | 19 | CWE-396 |

Checklist sources: `plans/fp-validations/reports/{caniscrape,voicetag,among-llms,FuncToWeb}.md`
(Post-fix over-suppression audit sections) + batch-3 plan
`plans/fp-review-reduce/python/batch-3-broad-except-reduction-2026-08-02.md`.

## Root cause

Batch-3 incorrectly copied BP-PY-1’s `suiteSurfacesFailure` exemption into
`detectCWE396`. CWE-396 is declaration-of-generic-catch; re-raise / `exc_info` /
`set_exception` / `result.error =` must not silence it. BP-PY-1 correctly keeps
suite-surface skipping.

## Detector change

File: `internal/lang/python/detectors/cwe/rules_platform.go`

- Removed `suiteSurfacesFailureMasked` skip from `detectCWE396`
- Deleted `suiteSurfacesFailureMasked` and `suiteLineSurfacesFailure` from the
  `cwe` package (unused after the skip removal)
- Restored comment to declaration-only semantics

**Not touched:** BP-PY-1 / `internal/lang/python/detectors/bad_practices/rules_core.go`
`suiteSurfacesFailure`.

## Fixtures

| Path | Role |
| --- | --- |
| `tests/fixtures/python/cwe/CWE-396-batch3-reraise-from-vulnerable.txt` | Wrap-and-re-raise (voicetag shape: `except Exception as exc: raise DomainError(...) from exc`) — **must fire** |
| `tests/fixtures/python/cwe/CWE-396-batch3-reraise-from-safe.txt` | Specific `except OSError` wrap-and-re-raise — **silent** |

Deleted (wrong CWE-396 “safe silence” pairs that still declared `except Exception`):

- `CWE-396-batch3-exc-info-{safe,vulnerable}.txt`
- `CWE-396-batch3-set-exception-{safe,vulnerable}.txt`
- `CWE-396-batch3-result-error-{safe,vulnerable}.txt`

BP-PY-1-batch3-* fixtures left unchanged.

## Tests

`audit_variants_test.go`: kept only `CWE-396-batch3-reraise-from` among batch3 cases.

```text
go test ./internal/lang/python/detectors/cwe/ -count=1 -run 'CWE.?396|Audit|Fixture'
ok  	github.com/chinmay-sawant/goslop/internal/lang/python/detectors/cwe	0.092s
```

## Residual risks

- BP-PY-1 still exempts surfacing suites; that is intentional and separate from CWE-396.
- Other audited over-suppression families in the same reports (e.g. BP-PY-46 Rich-print /
  examples path) are out of scope for this fix.
