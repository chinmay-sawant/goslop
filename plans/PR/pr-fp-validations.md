# goslop - Pull Request

---

## Summary

Adds the complete false-positive audit deliverable: 41 per-repo audit reports, a master report with cross-repo root-cause analysis, and an executive summary (renamed `0_SUMMARY.md` to sort first) under `plans/fp-validations/`, plus the reusable audit instructions (`AUDIT_PROMPT.md`) that produced them. The audit classified **6,964 findings** from scans of 42 real-world Python repos into **2,837 false positives (40.7%) / 4,115 true positives (59.1%) / 12 uncertain**, and distills the top detector weaknesses into concrete, fixable root causes for the Python rule set.

---

## Motivation / context

- Plans: `plans/fp-validations/`, `plans/skills/false-positive-audit/`
- Issues: none (validation corpus, no ticket)

---

## Changes

### False-positive audit reports (43 files under `plans/fp-validations/reports/`)

One report per scanned repo, following the false-positive-audit template (run metadata, classification summary, per-finding false-positive subsections with source excerpts, uncertain findings, `git diff --check` validation):

- Ai-copypaste-insult, CourtScrapper, Cronboard, FlashySurf, FuncToWeb, Project_Parva, PyDepends, WHEN-Language, WeThePeople, among-llms, astroz, calgebra, caniscrape, cylinder, enso, graphzero, html2pic, httpmorph, httptap, logxide, movielite, niquests, numeth-Numerical-Methods-Library, onlymaps, pdf_oxide, pictex, pingram, polygon-screenshot-tool, pyauto-desktop, pycaps, pyhash-complete, python-injection, pytogether, rendercv, requestSpeedTest, safer, sync-with-uv, tenso, violit, voicetag, whatsapp-wrapped, wse

(polygon-screenshot-tool is included with 0 findings — its only code file is `.pyw`, which the scanner skips.)

### Summary documents

- `0_SUMMARY.md` — executive verdict and top detector weaknesses (renamed from `SUMMARY.md` so it sorts first)
- `MASTER.md` — per-repo breakdown table, cross-repo root-cause analysis, fix recommendations per rule

### Audit instructions

- `plans/skills/false-positive-audit/AUDIT_PROMPT.md` — shared instructions for the 41 parallel audit agents (evidence locations, scan command, classification workflow, report output path)

---

## Summary of audit results

| Classification | Count | Share |
| --- | ---: | ---: |
| False positive | 2,837 | 40.7% |
| True positive | 4,115 | 59.1% |
| Uncertain | 12 | 0.2% |

### Top detector weaknesses found (fixable)

1. **BP-PY-46** (~1,400 FP) — no script-module detection (examples/, scripts/, CLI commands) and no docstring masking
2. **BP-PY-41** (~250 FP) — misses pytest-regressions `check_func`, pytest-benchmark, test-helper delegation; body-scan aborts on dedented string lines
3. **Broad-except family** (~350 FP) — doesn't credit re-raises/error recording/logging in the handler suite (BP-PY-1, CWE-396, BP-PY-2, CWE-390, CWE-1071)
4. **Identifier collisions** (~150 FP) — `exec`/`open`/`execute` name triggers on ORM/framework methods (onlymaps, violit, pdf_oxide)
5. Smaller: CWE-117 non-string interpolations, CWE-1341 different-handle pairs, PERF-PY-26 name-substring "hot path", CWE-88 constant-only argv

### Highest FP-rate repos (targets for rule fixes)

Project_Parva 97%, httptap 100%, pictex 80%, niquests 90%, pdf_oxide 70%.
Cleanest (high TP): among-llms 0% FP, WeThePeople 95% TP, CourtScrapper, graphzero, pytogether.

---

## Impact

| Area | Impact |
|------|--------|
| **Performance** | None (docs/plans only) |
| **Memory** | None |
| **Behavior / correctness** | None — no code changed |
| **API / CLI** | None |
| **Dependencies** | None |
| **Binary size / build time** | None |

---

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| None | - |

---

## Test plan

- [ ] `git diff --check` (run on all reports — pass)

### Commands

```sh
git diff --check
```

---

## Screenshots / sample output

```
safer: 40 findings — 33 FP / 7 TP / 0 U
WeThePeople: 1,492 findings — 70 FP / 1,414 TP / 8 U
httpmorph: 714 findings — 278 FP / 436 TP / 0 U
```

---

## Related issues

- None

---

## PR metadata checklist (author)

- [x] Self-assigned (`--assignee @me`)
- [x] Labels applied
- [ ] Related issues filled with real ticket IDs (none exist)
- [x] Filled body committed under `plans/PR/pr-fp-validations.md`

---

## Follow-ups (out of scope)

- Fix BP-PY-46 script-module/docstring detection
- Fix BP-PY-41 assertion-idiom recognition and body-scan bug
- Fix broad-except suite analysis (re-raise/record/log credit)
- Fix `exec`/`open`/`execute` identifier collisions

---

## Reviewer checklist

- [ ] Behavior matches summary and test plan
- [ ] No unrelated changes in diff
- [ ] Public API / CLI changes documented (none)
- [ ] New rules have fixture coverage when applicable (none — no rule changes)
- [ ] PR has assignee and labels
- [ ] Related issues use correct Closes/Relates keywords
- [ ] No secrets or generated artifacts committed
