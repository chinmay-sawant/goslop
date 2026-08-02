# Post-fix TP-match validation report

**Date:** 2026-08-02 (post `b5b8fde fix(python): reduce showcase-corpus false positives across BP/CWE/PERF`, binary rebuilt 16:29)

**Goal:** after the FP-reduction fix, each repo's scan findings should equal its audited **true positive** count from the false-positive audit. Deltas are either *remaining false positives* (found > TP) or *over-suppressed true positives* (found < TP).

**Scan command:**

```sh
for d in real-repos/*/; do name="${d%/}"; p="${name##*/}"; \
  ./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml \
    --export-context --export-chunks --no-cache \
    -chunks-dir "scripts/$p/chunks" \
    -context-dir "scripts/$p/findings/functions" \
    "$name" >/dev/null 2>&1; \
done
```

**Reference totals (0_SUMMARY):** TP 4,125 / FP 2,839 / U 0. Fresh scan total: 4,574 findings.

## Per-repo comparison (scan findings vs audited TP)

| Repo | Found (fresh scan) | TP (audit) | Delta | Interpretation |
| --- | ---: | ---: | ---: | --- |
| Ai-copypaste-insult | 2 | 2 | 0 | match |
| PyDepends | 4 | 4 | 0 | match |
| numeth-Numerical-Methods-Library | 4 | 4 | 0 | match |
| pingram | 6 | 6 | 0 | match |
| polygon-screenshot-tool | 0 | 0 | 0 | match (no code scanned) |
| CourtScrapper | 337 | 328 | +9 | 9 remaining FPs |
| Cronboard | 66 | 65 | +1 | 1 remaining FP |
| FlashySurf | 5 | 1 | +4 | 4 remaining FPs |
| FuncToWeb | 32 | 38 | -6 | 6 TPs over-suppressed |
| Project_Parva | 171 | 11 | +160 | 160 remaining FPs |
| WHEN-Language | 42 | 46 | -4 | 4 TPs over-suppressed |
| WeThePeople | 1298 | 1422 | -124 | 124 TPs over-suppressed |
| among-llms | 70 | 72 | -2 | 2 TPs over-suppressed |
| astroz | 17 | 15 | +2 | 2 remaining FPs |
| calgebra | 34 | 19 | +15 | 15 remaining FPs |
| caniscrape | 67 | 79 | -12 | 12 TPs over-suppressed |
| cylinder | 12 | 6 | +6 | 6 remaining FPs |
| enso | 134 | 127 | +7 | 7 remaining FPs |
| graphzero | 6 | 7 | -1 | 1 TP over-suppressed |
| html2pic | 28 | 32 | -4 | 4 TPs over-suppressed |
| httpmorph | 427 | 436 | -9 | 9 TPs over-suppressed |
| httptap | 45 | 0 | +45 | 45 remaining FPs (was 103 FPs) |
| logxide | 377 | 317 | +60 | 60 remaining FPs |
| movielite | 59 | 58 | +1 | 1 remaining FP |
| niquests | 271 | 35 | +236 | 236 remaining FPs |
| onlymaps | 60 | 36 | +24 | 24 remaining FPs |
| pdf_oxide | 202 | 188 | +14 | 14 remaining FPs |
| pictex | 41 | 42 | -1 | 1 TP over-suppressed |
| pyauto-desktop | 86 | 88 | -2 | 2 TPs over-suppressed |
| pycaps | 29 | 28 | +1 | 1 remaining FP |
| pyhash-complete | 9 | 10 | -1 | 1 TP over-suppressed |
| python-injection | 1 | 2 | -1 | 1 TP over-suppressed |
| pytogether | 68 | 70 | -2 | 2 TPs over-suppressed |
| rendercv | 36 | 27 | +9 | 9 remaining FPs |
| requestSpeedTest | 11 | 9 | +2 | 2 remaining FPs |
| safer | 13 | 7 | +6 | 6 remaining FPs |
| sync-with-uv | 7 | 5 | +2 | 2 remaining FPs |
| tenso | 75 | 71 | +4 | 4 remaining FPs |
| violit | 224 | 219 | +5 | 5 remaining FPs |
| voicetag | 6 | 13 | -7 | 7 TPs over-suppressed |
| whatsapp-wrapped | 23 | 23 | 0 | match |
| wse | 169 | 157 | +12 | 12 remaining FPs |

## Totals

| Group | Repos | Sum |
| --- | ---: | ---: |
| Match (found == TP) | 6 | — |
| Remaining FPs (found > TP) | 22 | +625 |
| Over-suppressed TPs (found < TP) | 14 | -176 |

## Remaining-FP hotspots (biggest, for next fix round)

- niquests +236 (was 311 FP) — largest residual
- Project_Parva +160 (was 401 FP)
- logxide +60 (was 185 FP)
- httptap +45 (was 103 FP — 100% FP repo, still not clean)
- onlymaps +24, calgebra +15, pdf_oxide +14, wse +12

## Over-suppression (true positives removed by the fix — regression risk)

- WeThePeople -124 (TP fell 1,414 → likely over-broad suppression)
- caniscrape -12, httpmorph -9, voicetag -7, FuncToWeb -6
- WHEN-Language -4, html2pic -4, others -1/-2 each

**Evidence to check:** each repo's fresh findings are in `scripts/<project>/chunks/` and `scripts/<project>/findings/functions/`; audited TP/FP lists are in `reports/<project>.md` classification summary.
