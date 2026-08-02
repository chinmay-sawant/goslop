# False-positive audit — summary of results

Source: full audit of goslop scan findings on 42 real-world Python repos (r/Python top showcase projects).
Per-repo detail: see `reports/<name>.md`. Aggregated detail: see `reports/MASTER.md`.

## Deliverables

- 41 per-repo reports at `plans/skills/false-positive-audit/reports/<name>.md` (~59,400 lines), each following the false-positive-audit template: run metadata, classification summary, per-finding false-positive subsections with source excerpts, uncertain findings, `git diff --check` validation (all pass)
- Master report `reports/MASTER.md` with cross-repo root-cause analysis
- Evidence per repo under `real-repos/<name>/scripts/{chunks,findings/functions}` (re-scanned with `-chunks-dir`/`-context-dir` after the first run's exports collided in the root)

## Verdict across 6,964 findings

| Classification | Count | Share |
| --- | ---: | ---: |
| False positive | 2,837 | 40.7% |
| True positive | 4,115 | 59.1% |
| Uncertain | 12 | 0.2% |

## Top detector weaknesses (fixable)

1. **BP-PY-46** (~1,400 FP) — no script-module detection (examples/, scripts/, CLI commands) and no docstring masking
2. **BP-PY-41** (~250 FP) — misses pytest-regressions `check_func`, pytest-benchmark, test-helper delegation; body-scan aborts on dedented string lines
3. **Broad-except family** (~350 FP) — doesn't credit re-raises/error recording/logging in the handler suite (BP-PY-1, CWE-396, BP-PY-2, CWE-390, CWE-1071)
4. **Identifier collisions** (~150 FP) — `exec`/`open`/`execute` name triggers on ORM/framework methods (onlymaps, violit, pdf_oxide)
5. Smaller: CWE-117 non-string interpolations, CWE-1341 different-handle pairs, PERF-PY-26 name-substring "hot path", CWE-88 constant-only argv

## FP rate by repo

Worst FP rates: Project_Parva 97%, httptap 100%, pictex 80%, niquests 90%, pdf_oxide 70%.
Cleanest (high TP): among-llms 0% FP, WeThePeople 95% TP, CourtScrapper, graphzero, pytogether.
