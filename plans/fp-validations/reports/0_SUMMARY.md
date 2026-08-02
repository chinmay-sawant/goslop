# False-positive audit — summary of results

Source: full audit of goslop scan findings on 42 real-world Python repos (r/Python top showcase projects).
Per-repo detail: see `reports/<name>.md`. Aggregated detail: see `reports/MASTER.md`.

## Session token usage (estimate)

Session window: 2026-08-02 ~11:00 IST → 14:04 IST (approx. 3h), incl. repo cloning (42 repos), 42 goslop scans, 41 parallel audit sub-agents, report writing, commit + PR #74.

| Component | Measured chars | Est. tokens (≈ chars/4) |
| --- | ---: | ---: |
| Report output written (43 files) | 3,149,362 | ~790k |
| Chunks read by auditors (all 6,964 findings) | 5,922,125 | ~1,480k |
| FP/uncertain function contexts read | ~2,150k (41%) | ~540k |
| Source files read beyond contexts (unmeasured) | — | ~400k |
| Auditor reasoning/generation overhead (~2–3× final output) | — | ~1,400k |
| Main session (dispatches, results, scans, PR) | — | ~350k |
| **Total (estimate)** | | **~4.9M ± 25%** |

Method: measured bytes of written reports and read evidence on disk; code-heavy markdown/source ≈ 4 chars/token; auditor reasoning multiplier and main-session overhead are experience-based estimates, not metered. Exact numbers depend on model context/tokenizer; treat as ±25%.

### Agent count

| Agent | Count | Work |
| --- | ---: | --- |
| Main session (opencode, deepseek-v4-flash) | 1 | cloning, scans, orchestration, reports aggregation, commit, PR |
| Parallel audit sub-agents (general) | 41 | one per repo with findings; each read all chunks + FP contexts + source, wrote per-repo report (waves of 10/10/10/11) |
| **Total agents** | **42** | (polygon-screenshot-tool had 0 findings — no agent; scanned by main session) |

### Estimated cost

Assumed list pricing (deepseek-v4-flash-class, approx.): $0.28/M input tokens, $0.42/M output tokens. Split of the ~4.9M tokens: ~2.6M input (evidence reading, contexts, source) / ~2.4M output (reports, reasoning, orchestration).

| Metric | Value |
| --- | --- |
| Input tokens | ~2.6M × $0.28/M ≈ **$0.73** |
| Output tokens | ~2.4M × $0.42/M ≈ **$1.01** |
| **Total (mid estimate)** | **≈ $1.74** (≈ ₹145) |
| Range (±25% tokens, ±30% pricing) | **$1.2 – $2.4** |

Pricing varies by provider/plan (API list price vs subscription); treat as an order-of-magnitude estimate, recompute with your actual rates.

## Deliverables

- 41 per-repo reports at `plans/fp-validations/reports/<name>.md` (~59,400 lines), each following the false-positive-audit template: run metadata, classification summary, per-finding false-positive subsections with source excerpts, uncertain findings, `git diff --check` validation (all pass)
- Master report `reports/MASTER.md` with cross-repo root-cause analysis
- Evidence per repo under `scripts/<name>/{chunks,findings/functions}` (re-scanned with `-chunks-dir`/`-context-dir` after the first run's exports collided in the root)

## Verdict across 6,964 findings

| Classification | Count | Share |
| --- | ---: | ---: |
| False positive | 2,839 | 40.8% |
| True positive | 4,125 | 59.2% |
| Uncertain | 0 | 0.0% |

All 12 former uncertains reclassified (WeThePeople 8×CWE-1084 → TP; pycaps CWE-22 → TP; violit CWE-829 → TP; logxide CWE-1121 → FP; pdf_oxide CWE-88 → FP). See `MASTER.md` § Uncertain findings.

## Top detector weaknesses (fixable)

1. **BP-PY-46** (~1,400 FP) — no script-module detection (examples/, scripts/, CLI commands) and no docstring masking
2. **BP-PY-41** (~250 FP) — misses pytest-regressions `check_func`, pytest-benchmark, test-helper delegation; body-scan aborts on dedented string lines
3. **Broad-except family** (~350 FP) — doesn't credit re-raises/error recording/logging in the handler suite (BP-PY-1, CWE-396, BP-PY-2, CWE-390, CWE-1071)
4. **Identifier collisions** (~150 FP) — `exec`/`open`/`execute` name triggers on ORM/framework methods (onlymaps, violit, pdf_oxide)
5. Smaller: CWE-117 non-string interpolations, CWE-1341 different-handle pairs, PERF-PY-26 name-substring "hot path", CWE-88 constant-only argv

## FP rate by repo

Worst FP rates: Project_Parva 97%, httptap 100%, pictex 80%, niquests 90%, pdf_oxide 70%.
Cleanest (high TP): among-llms 0% FP, WeThePeople 95% TP, CourtScrapper, graphzero, pytogether.
