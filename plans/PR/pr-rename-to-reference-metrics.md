## Summary

Replace ambiguous “oracle” wording with standard testing language: **reference corpus**, **expected baseline**, and **parity metrics**. Rename `make oracle` to `make reference-metrics` and align docs, plans, and detector helper names.

---

## Motivation / context

- “Oracle” was testing jargon (expected baseline) but read like a product name or vendor.
- Clarify that **gopdfsuit** is only a **reference corpus** (sample app tree), not “oracle software.”
- Self-scan noise guidance remains: SAT self-host is like Semgrep on its own rules/fixtures.

---

## Changes

### Makefile

- `make oracle` → **`make reference-metrics`**
- `ORACLE_PATH` / `ORACLE_PROFILE` → **`REFERENCE_PATH` / `REFERENCE_PROFILE`**
- Output: `/tmp/goslop-reference-metrics.json`
- Comments document “reference = testing expected metrics”

### Code

- `oracle_parity.go` → `reference_parity.go` (BP + PERF)
- `oracleSkip` / `oracleSkipUnit` → `referenceSkip` / `referenceSkipUnit`
- `isOraclePureFP` → `isReferencePureFP`
- `TestSeedFixtureOracle` → `TestSeedFixtureMatrix`

### Docs / plans

- README, `documents/*`, port checklist, PR template, historical plans: oracle → reference / baseline / fixture expectations

---

## Impact

| Area | Impact |
|------|--------|
| **Behavior / correctness** | None (rename + docs only) |
| **API / CLI** | Makefile target rename only (`make reference-metrics`) |
| **Dependencies** | None |

---

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| `make oracle` | Use `make reference-metrics` |
| `ORACLE_PATH` / `ORACLE_PROFILE` | Use `REFERENCE_PATH` / `REFERENCE_PROFILE` |

---

## Test plan

- [x] `go test` detectors + integration after renames
- [x] `rg -i oracle` clean (repo sources/docs)
- [x] `make help` shows `reference-metrics`

```sh
make help
go test ./internal/lang/go/detectors/... ./tests/integration/ -count=1 -timeout 60s
```

---

## Related issues

- Relates to product docs clarity (no GitHub issue ID)

---

## PR metadata checklist (author)

- [x] Self-assigned
- [x] Labels applied
- [x] Body under `plans/PR/pr-rename-oracle-to-reference-metrics.md`
