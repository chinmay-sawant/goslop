## Summary

Standardize product and plan wording on testing conventions: **reference corpus**, **expected baseline**, and **parity metrics**. Rename the Makefile metrics gate to `make reference-metrics` and align detector helper names accordingly.

---

## Motivation / context

- Informal “expected-answer” jargon was easy to misread as a product or vendor name.
- **gopdfsuit** is only a **reference corpus** (sample app tree used for §12.4 parity), not special software.
- Plans: naming cleanup for docs, Makefile, and detector helpers.

---

## Changes

### Makefile

- Metrics gate target: **`make reference-metrics`**
- Env: **`REFERENCE_PATH` / `REFERENCE_PROFILE`**
- JSON output: `/tmp/goslop-reference-metrics.json`

### Code

- `reference_parity.go` (BP + PERF)
- `referenceSkip` / `referenceSkipUnit`
- `isReferencePureFP`
- `TestSeedFixtureMatrix`

### Docs / plans

- README, `documents/*`, port checklist, PR template: use reference / baseline / fixture expectations

---

## Impact

| Area | Impact |
|------|--------|
| **Behavior / correctness** | None (rename + docs) |
| **API / CLI** | Makefile target rename only |
| **Dependencies** | None |

---

## Breaking changes / migration

| Former | Use instead |
|--------|-------------|
| Legacy Makefile metrics target | `make reference-metrics` |
| Legacy path/profile env vars | `REFERENCE_PATH` / `REFERENCE_PROFILE` |

---

## Test plan

- [x] `rg -i` clean for the retired informal metrics jargon (sources/docs)
- [x] `make help` shows `reference-metrics`

---

## Related issues

- Relates to product docs clarity

---

## PR metadata checklist (author)

- [x] Self-assigned
- [x] Labels applied
- [x] Body under `plans/PR/pr-rename-to-reference-metrics.md`
