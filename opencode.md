# opencode session log — showcase-corpus FP validation & reduction (2026-08-02)

This file records what was done in this session, end to end: commands executed, how findings were produced and classified, and where the evidence lives. It references the other markdown documents in this repo rather than repeating them.

## Purpose

Validate goslop's Python detector precision against 42 real-world Python repositories (r/Python top showcase projects, 2025–2026), classify every finding as false positive / true positive / uncertain, then iteratively narrow the detectors until the scan output matches the audited true positives only.

Ground truth of the whole exercise is the **false-positive audit template**:

- `plans/skills/false-positive-audit/SKILLS.md` — per-finding audit report template
- `plans/skills/false-positive-audit/AUDIT_PROMPT.md` — instructions used by the 41 audit agents
- `plans/skills/false-positive-audit/REAUDIT_PROMPT.md` — instructions for the Mode A/B re-audits
- `plans/skills/false-positive-audit/V2_AUDIT_PROMPT.md` — instructions for the post-fix-v2 audits
- `plans/skills/review-and-reduce/SKILLS.md` — the rule-narrowing workflow used for every detector fix

## Phase 0 — Corpus setup (cloning)

Cloned the 42 repos listed in the r/Python showcase roundup into `real-repos/` (41 GitHub + `enso` from GitLab). The repo list and upvote data are referenced in the roundup; the clones are gitignored (`real-repos` in `.gitignore`, commit `5a6f4aa`).

```sh
xargs -P 8 -I{} sh -c 'git clone --quiet "{}"' < /tmp/opencode/repos.txt
```

## Phase 1 — Scanning

The canonical scan command (per project, exporting chunks + per-finding function contexts):

```sh
./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml \
  --export-context --export-chunks --no-cache \
  -chunks-dir "scripts/$p/chunks" \
  -context-dir "scripts/$p/findings/functions" \
  "real-repos/$p"
```

Run in a loop over all 42 projects. Config: `templates/goslop-python.toml` (languages = ["python"], fail_on = "none").

History / mistakes worth knowing:

1. **First run (no `-chunks-dir`/`-context-dir`)**: all exports collided into the goslop root `scripts/` — only the last repo's evidence survived. Root cause: default export dirs are relative to the working directory.
2. **Per-repo export (run 2)**: `-chunks-dir real-repos/<name>/scripts/chunks` — worked, but…
3. **Tracked-scripts incident**: cleanup `rm -rf real-repos/<name>/scripts` deleted **tracked** `scripts/` dirs that 7 repos ship in their own git history (Project_Parva ×212, pdf_oxide ×60, WeThePeople ×30, rendercv ×28, httpmorph ×14, logxide ×4, sync-with-uv ×1 = 349 files). All restored via `git -C real-repos/<name> checkout -- scripts`; the 7 repos were re-scanned and counts verified identical.
4. **Final layout (current)**: evidence lives centrally under `scripts/<project>/{chunks,findings/functions}` (gitignored). 42 project folders; counts table for every project is in `plans/fp-validations/reports/MASTER.md`.

## Phase 2 — Full false-positive audit (the 6,964-finding classification)

41 parallel audit sub-agents (one per repo with findings; polygon-screenshot-tool had 0 findings — its only code file is `.pyw` which the scanner skips) classified **every finding** per the template: chunks + function contexts + enclosing source + rule conditions (`./bin/goslop -explain <rule>`).

Results (aggregated in `plans/fp-validations/reports/MASTER.md` and `0_SUMMARY.md`):

| Classification | Count | Share |
| --- | ---: | ---: |
| False positive | 2,837 | 40.7% |
| True positive | 4,115 | 59.1% |
| Uncertain | 12 | 0.2% |

Cross-repo root causes (each with per-repo counts and fix recommendations) are documented in `MASTER.md`:

1. **BP-PY-46** print-in-library (~1,400 FP) — no script-module detection, no docstring masking
2. **BP-PY-41** test-with-no-assertions (~250 FP) — misses `check_func`/pytest-benchmark/helper delegation; body-scan aborts on dedented string lines
3. **Broad-except family** (~350 FP) — doesn't credit re-raises/error recording/logging (BP-PY-1, CWE-396, BP-PY-2, CWE-390, CWE-1071)
4. **Identifier collisions** (~150 FP) — `exec`/`open`/`execute` name triggers on ORM/framework methods
5. Smaller: CWE-117 non-string interpolations, CWE-1341 different-handle pairs, PERF-PY-26 name-substring "hot path", CWE-88 constant-only argv

Deliverables: 41 per-repo reports (template-structured, ~59,400 lines) under `plans/fp-validations/reports/<repo>.md`, plus `MASTER.md` and `0_SUMMARY.md` (renamed from `SUMMARY.md` to sort first; includes the session token/cost estimate and agent counts).

Committed as `65eae01` and PR'd: **https://github.com/chinmay-sawant/goslop/pull/74** (body: `plans/PR/pr-fp-validations.md`, template: `plans/PR/PR_TEMPLATE.md`). Later reclassification of the 12 uncertains (`105ec5a`) moved totals to FP 2,839 / TP 4,125 / U 0.

## Phase 3 — Post-fix validation (commit `b5b8fde`)

The user's fix `fix(python): reduce showcase-corpus false positives across BP/CWE/PERF` added guardrails + fixtures. Re-scanned with the rebuilt binary and compared per-repo findings vs audited TP counts → `plans/fp-validations/reports/TP_MATCH.md`:

- 6 repos exact, 22 with remaining FPs (+625), 14 with over-suppressed TPs (−176).
- **Over-suppression root causes found via 10 explore agents** (shebang-on-library-module, type-agnostic re-raise exemption, path-only examples/ exemption, `__file__` token-only check, record-into-field handlers) — documented in the TP_MATCH report and in the Mode B appends.

Re-audits appended to the 36 affected reports:

- `## Post-fix remaining-FP audit (2026-08-02)` (22 repos, Mode A, 662 remaining FPs)
- `## Post-fix over-suppression audit (2026-08-02)` (14 repos, Mode B, 255 suppressed-but-present TPs, 0 fixed-removed)

Committed as `3407305`.

## Phase 4 — Fix iterations (the loop you asked for)

Goal: findings == audited TPs per repo. Loop per round: audit → aggregate → fix detectors (review-and-reduce) → `make build` → rescan → compare.

Round markers (all reproducible — scans write to `scripts/<project>/`, logs were kept in `/tmp/opencode/scan*.log`):

| Round | Binary state | Total | Excess | Under |
| --- | --- | ---: | ---: | ---: |
| 0 (original) | pre-fix | 6,964 | — | — |
| 1 (your binary) | `b5b8fde`-reverted guardrails | 5,197 | 1,072 | 0 |
| 2 (first fix wave, F1–F11) | 7 agents | 4,671 | 716 | 170 |
| 3 (reconciliation round) | 4 agents (BP-PY-45/46/41, except-family) | 4,575 | 486 | 36 |

Aggregated fix checklist with per-rule trigger shapes and counts: `plans/fp-validations/reports/FIX_CHECKLIST.md`.

Fix waves:

- **Wave 1 (F1–F11)**: 7 agents, disjoint file partitions — BP-PY-46 script/docstring/def-print/test-case, broad-except family, identifier collisions (BP-PY-12/CWE-89/94/829), CWE-117/88/93, perf (PERF-PY-25/26/27/28), resource rules (CWE-1341/367/478/260/256/459/772/779/186/409/695), misc BP (BP-PY-40/13/44/42/41/14/37/36/49/32, CWE-1121).
- **Wave 2 (reconciliation)**: BP-PY-45 `isSysPathBootstrap` (module-level guarded bootstrap; Project_Parva 133 → 0 while WeThePeople/html2pic TPs kept firing), BP-PY-46 rework (examples/ prints only when the module runs itself), except-family reversal (restored WeThePeople 494 / Cronboard 21 / pyauto-desktop 30 / pytogether 14 / httpmorph 151 BP-PY-1 TPs; kept calgebra/onlymaps/httptap FPs suppressed), BP-PY-41 (movielite ×24 verified as audited TPs; `raise AssertionError` counts as assertion; pytest.fail only inside suites).

Every guardrail change followed `review-and-reduce/SKILLS.md`: safe + vulnerable text fixture pairs under `tests/fixtures/python/{bp,cwe,perf}/`, wired into `audit_variants_test.go` / rule test matrices (fixtures only, never inline snippets), focused tests green (`go test ./...` — 446 tests, 0 failures).

## Phase 5 — Historical checkpoint (Round 5, ~20:30)

Earlier rescan (Round 5) vs audited TP:

- **Total: 4,575 findings** (from 6,964 original — 34% reduction with 100% TP preservation rate 4,125/4,125 minus 36)
- **Excess: 486** across 24 repos — top: niquests +138, Project_Parva +118, pdf_oxide +61, WeThePeople +42, logxide +40, httptap +14
- **Under: 36** across 11 repos (mostly −1/−2; worst voicetag −7, caniscrape −12)

Later intermediate markers (same day, continued waves on `fix/python-fp-reduction`):

| Marker | NOW | Excess | Under | Net (NOW−TP) | Notes |
| --- | ---: | ---: | ---: | ---: | --- |
| Round 5 | 4,575 | 486 | 36 | +450 | early post-wave |
| goal15 | ~4,235 | ~132 | ~22 | ~+110 | pre last parallel wave |
| goal16 | 4,196 | 97 | 26 | +71 | before wave below |
| **goal17 (latest)** | **4,110** | **36** | **51** | **−15** | after 4 parallel agents |

Pushed commits on branch (prior waves): `5a1bb3b`, `6fe9aae`, `93370ca`. **goal17 detector edits are uncommitted** on the working tree.

## Phase 6 — Latest state (goal17, 2026-08-02 late)

**Goal:** findings == audited TP **4,125** (per-repo preferred).

**Latest full corpus rescan:** `/tmp/opencode/goal17/*.json`  
**TP table:** `/tmp/opencode/compare.txt` (col2 = audited TP)  
**Baseline before this wave:** `/tmp/opencode/goal16/` (NOW 4,196 / excess 97 / under 26)

| Metric | Value |
| --- | ---: |
| NOW | **4,110** |
| TP | **4,125** |
| Excess (found > TP) | **36** |
| Under (found < TP) | **51** |
| Net | **−15** (slightly under total) |

### Exact repos (delta 0)

Ai-copypaste-insult, FlashySurf, PyDepends, among-llms, caniscrape, httptap, numeth-Numerical-Methods-Library, pingram, polygon-screenshot-tool, requestSpeedTest, voicetag, whatsapp-wrapped.

### Remaining excess (+36)

| Repo | now | tp | excess |
| --- | ---: | ---: | ---: |
| CourtScrapper | 336 | 328 | **+8** |
| WHEN-Language | 51 | 46 | **+5** |
| FuncToWeb | 43 | 38 | **+5** |
| pycaps | 31 | 28 | **+3** |
| Project_Parva | 13 | 11 | **+2** |
| Cronboard | 67 | 65 | **+2** |
| tenso | 73 | 71 | **+2** |
| sync-with-uv | 7 | 5 | **+2** |
| safer | 9 | 7 | **+2** |
| pictex / onlymaps / movielite / enso / cylinder | — | — | **+1 each** |

### Remaining under (−51) — **priority**

This wave cut excess hard but overshot several TPs. Restore before more FP kills.

| Repo | now | tp | under | Wave lost shapes (goal16→goal17) |
| --- | ---: | ---: | ---: | --- |
| **httpmorph** | 424 | 436 | **−12** | CWE-390/1071 tests+examples; CWE-396 examples; CWE-215 example |
| **pdf_oxide** | 180 | 188 | **−8** | CWE-396 examples; CWE-390/1071 test; PERF-PY-27 scripts/regtest |
| **violit** | 211 | 219 | **−8** | CWE-396 examples + data_widgets (was −6; −2 more this wave) |
| **wse** | 153 | 157 | **−4** | CWE-215 benchmarks/examples overshot (was +8 excess) |
| **logxide** | 314 | 317 | **−3** | BP-PY-1 handleError/HTTP/health (−15) + CWE-396 examples (−7) overshot |
| WeThePeople | 1420 | 1422 | **−2** | CWE-396 `jobs/send_alerts.py:351`; CWE-93 middleware |
| astroz | 13 | 15 | **−2** | examples/cesium_fast CWE-396/390/1071 |
| niquests / pytogether / rendercv | — | — | **−2 each** | rendercv: PERF-PY-27 skill script; niquests pre-existing |
| calgebra / graphzero / html2pic / pyauto-desktop / pyhash-complete / python-injection | — | — | **−1 each** | html2pic soft-warn CWE-396 may be audited FP; check before restore |

### What this wave already shipped (uncommitted)

Four parallel agents (non-overlapping file ownership):

1. **CWE** (`rules_platform` / resource / tier_b + cwe fixtures) — Parva −25 CWE; soft-warn; under-restore run_story_gates + caniscrape wrap-after-unlink; many CWE-396 guards
2. **BP-PY-1/2** (`rules_core`) — logxide −15 BP-PY-1 (`handleError`, framework HTTP return, `"unhealthy"`)
3. **BP-41/45/46** (`rules_testing` / prod / observability) — Parva BP-PY-41×3 + BP-PY-45×1; light BP-46
4. **Path/misc** (path_fs / injection / info_exposure / perf hotpath) — Parva CWE-22/73/88/93/215 + PERF-23/27

Agent reports: `/tmp/opencode/agent_{cwe,bp12,bp414647,path_misc}_report.md`  
Residual brief: `/tmp/opencode/residual_goal16_brief.md`

Detector packages currently **green** (`go test ./internal/lang/python/detectors/...`).

---

## Pending (do next)

### P0 — Under-restore (must fix first; net is −15)

Policy remains **TP-safe**: prefer re-firing lost audited TPs over chasing the last FPs.

1. **httpmorph −12** — reverse over-broad CWE-390/1071/396 guards that hit `tests/` + `examples/` (proxy/async/advanced). Prefer path-narrow or shape-narrow restores that still keep Project_Parva load_test/timegraph FPs dead.
2. **pdf_oxide −8** — restore CWE-396 example sites, test_api_coverage silence, PERF-PY-27 `scripts/regtest_branch_vs_main.py` without reopening Parva PERF FPs.
3. **violit −8** — restore CWE-396 on demo_app_* and `data_widgets.py` if audited TPs; do not re-inflate soft-warn-only html2pic FPs unless required.
4. **wse −4** — CWE-215 benchmark/example prints over-suppressed by offline/release-style guards; narrow so wse TPs fire while Parva `generate_partner_api_key` FP stays suppressed.
5. **logxide −3** — decide which of the 15 BP-PY-1 / 7 CWE-396 kills were audited TPs vs FPs; restore only TPs (handleError contract is likely correct FP kill — check report).
6. **WeThePeople −2** — restore `jobs/send_alerts.py:351` CWE-396 and `middleware/security.py:121` CWE-93 if still audited TPs (attr-cache / header guards too broad).
7. **Small under (−1/−2)** — niquests, pytogether, rendercv, astroz, calgebra, graphzero, pyauto-desktop, pyhash-complete, python-injection; html2pic −1 only if wrap/soft-warn was TP.

### P1 — Remaining excess (+36) after under is healthy

1. **CourtScrapper +8** — residual not BP-PY-47 bulk (mostly TPs); likely BP-PY-1/46/CWE mix — need site-level FP list.
2. **WHEN-Language +5 / FuncToWeb +5** — BP-PY-46 examples/prints; audit conflict with FuncToWeb examples as TP — careful.
3. **pycaps +3** — residual CWE-396 wrap-raise / other (left wrap-raise on purpose for voicetag).
4. **Project_Parva +2** — down from +27; finish last two (likely CWE-186×2 or CWE-290 / CWE-93 residual / rulelang CWE-396 f-string wrap).
5. **Cronboard +2, tenso +2, safer +2, sync-with-uv +2** and five **+1** repos — low priority once under ≈ 0.

### P2 — Hygiene

1. **Commit + push** uncommitted wave once under-restore lands and scorecard improves (detailed message, same style as `93370ca`).
2. Full rescan → `/tmp/opencode/goal18/` + update this section.
3. Stop when **NOW ≈ 4,125** with excess and under both near 0 (accept audit-conflict residuals under TP-safe policy).

### Still open blockers (unchanged)

1. **Audit conflicts** — same shape, opposite labels:
   - `raise X from e`: voicetag TP vs html2pic FP
   - examples/ prints: movielite FP vs FuncToWeb TP
   - niquests pyodide probes vs wse probes
2. Suggested default: **TP-safe** (keep firing when ambiguous).
3. Do **not** wholesale-suppress BP-PY-47 on logxide/CourtScrapper (bulk audited TPs).

### Suggested next agent partition (under-restore)

| Agent | Owns | Focus |
| --- | --- | --- |
| A CWE platform | `cwe/rules_platform.go` + 390/1071/396 fixtures | httpmorph, pdf_oxide, violit, logxide CWE, WeThePeople send_alerts |
| B CWE info/path | `cwe/rules_info_exposure.go`, `rules_injection.go` | wse CWE-215, WeThePeople CWE-93 |
| C BP-1 | `bad_practices/rules_core.go` | logxide BP-PY-1 TP restore only if needed |
| D PERF | `perf/rules_hotpath.go` | pdf_oxide / rendercv PERF-PY-27 restore |

## Where everything lives

| Artifact | Location |
| --- | --- |
| Cloned corpus | `real-repos/<project>/` (gitignored) |
| Scan evidence (chunks + function contexts) | `scripts/<project>/{chunks,findings/functions}` (gitignored) |
| Per-repo audit reports (+ Mode A/B + v2 appends) | `plans/fp-validations/reports/<project>.md` |
| Aggregated classification + session stats | `plans/fp-validations/reports/0_SUMMARY.md` |
| Cross-repo root causes | `plans/fp-validations/reports/MASTER.md` |
| Post-fix validation | `plans/fp-validations/reports/TP_MATCH.md` |
| Fix checklist (per-rule trigger shapes) | `plans/fp-validations/reports/FIX_CHECKLIST.md` |
| Detector fixtures | `tests/fixtures/python/{bp,cwe,perf}/` |
| Audit/re-audit/fix instructions | `plans/skills/false-positive-audit/*.md`, `plans/skills/review-and-reduce/SKILLS.md` |
| PR | https://github.com/chinmay-sawant/goslop/pull/74 |
