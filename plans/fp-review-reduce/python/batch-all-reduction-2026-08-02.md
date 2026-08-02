# Showcase-corpus FP reduction — integrated batches (2026-08-02)

- Skill: `plans/skills/review-and-reduce/SKILLS.md`
- Checklists: `plans/fp-validations/reports/*.md` (41 repos) + `MASTER.md`
- Evidence override: per-report `scripts/<repo>/chunks` and `scripts/<repo>/findings/functions/<id>.txt` (not generic `scripts/chunks`)
- Audit report `[ ]` boxes left unchanged

## Parallel batches

| Batch | Agent scope | Ledger |
| --- | --- | --- |
| 1 | BP-PY-46 | [batch-1-bp-py-46-reduction-2026-08-02.md](./batch-1-bp-py-46-reduction-2026-08-02.md) |
| 2 | BP-PY-41, BP-PY-45 | [batch-2-bp-py-41-45-reduction-2026-08-02.md](./batch-2-bp-py-41-45-reduction-2026-08-02.md) |
| 3 | BP-PY-1, CWE-396, BP-PY-42 | [batch-3-broad-except-reduction-2026-08-02.md](./batch-3-broad-except-reduction-2026-08-02.md) |
| 4 | BP-PY-12, BP-PY-7, CWE-89, CWE-94 | [batch-4-identifier-collision-reduction-2026-08-02.md](./batch-4-identifier-collision-reduction-2026-08-02.md) |
| 5 | CWE-117/1341/367/88/93/215/502, PERF-PY-25/26, BP-PY-11/13/49 | [batch-5-misc-reduction-2026-08-02.md](./batch-5-misc-reduction-2026-08-02.md) |

## Integration validation (parent session)

```text
go test ./internal/lang/python/detectors/bad_practices/... \
        ./internal/lang/python/detectors/cwe/... \
        ./internal/lang/python/detectors/perf/... \
        ./tests/integration/python -count=1
```

**Result: pass** after merge cleanup (gofmt; named-return / nestif / unused-param fixes in batch-touched files).

- `make lint` still fails on `real-repos/pdf_oxide/...` gofmt (third-party clones; unrelated).
- Detector packages under `internal/lang/python/detectors/{bad_practices,cwe,perf}` gofmt-clean.
- Full `make test` / corpus-wide rescan not run in this integration pass.
- Zero-FP repos confirmed skipped by batch 5: among-llms, graphzero, pingram, PyDepends, python-injection, Ai-copypaste-insult.

## Residual notes (from batch ledgers)

- FlashySurf BP-PY-46 residuals without path/shebang signal
- Some CWE-88 harnesses outside bench paths
- BP-PY-2 / CWE-390 / CWE-1071 left unchanged (except-pass TPs)
