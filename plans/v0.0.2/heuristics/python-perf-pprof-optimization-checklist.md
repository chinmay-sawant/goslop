# v0.0.2 — Python performance optimization checklist (from pprof)

> **Parent:** [`python-heuristics.md`](./python-heuristics.md) — Python heuristic detectors  
> **Status:** **P0–P2 implemented + verified**; **Phase 5 verified (r2→after)** — engine **132.9 → 78.5 ms**/op · product mean **147.9 → 93.3 ms** / best **81.0** — closure **≤90 bench / ≤100 product met**; residual: backtrack share, CWE-1050, PERF/BP (see Phase 5)  
> **Branch:** `feat/python-perf-ruleset-plan`  
> **Evidence:** `/tmp/goslop-python-pprof/*` · [`python-perf-pprof-measurement.md`](./python-perf-pprof-measurement.md)  
> **Analysis:** [CWE hotspots](e668df02-2031-4b0d-8a90-c00390c64e81) · [PERF/BP secondary](4d88f7be-b6eb-4879-a9d0-d996b5540c54)  
> **Corpus:** `pythoncoreengine` · 63 files / 20945 lines · `--no-cache` · profile `all`  
> **Benches (pre → after, `benchtime=5s`, 2026-08-02):**  
> - `BenchmarkPythonScanProfileAll` ~**728.5 → 146.3 ms**/op · **553 → 24.5 MB** · **208k → 70.5k** allocs (**~5.0×** wall, **~22.5×** bytes)  
> - `BenchmarkPythonScanAndExport` ~**627.4 ms**/op · 557 MB · 213k allocs *(pre only; not re-run after)*  
> - Go compare `BenchmarkScanProfileAll` ~**125.1 ms**/op · 55 MB · 439k allocs (gopdfsuit)  
> - Product **`make run-python`**: pre ~**1.08–1.21 s** → after 10× mean **157.5 ms** / best **133.4 ms** (export on); Go **`make run`**: ~**94–112 ms**  
> - Findings after: **102** · 38h/56i/0l/8m · top BP-PY-46×32, BP-PY-41×24, CWE-328×14, CWE-90×12, CWE-22×6 (**parity**)

Legend: `[ ]` not started · `[~]` partial · `[x]` done with evidence

---

## Executive Summary

Baseline: **CWE pack (~80% CPU)** remasked each file **~175×** (`pytext.Mask` ≈ **86% alloc_space**) and ran **`FindAllStringIndex` via `firstCodeMatchStart` (~48% CPU)** across **154/159 ungated** rules. P0 cached Mask/Funcs + iterative first-match; P1 wired SourceIndex gates (**5 → 65** gated) + hoisted hot `MustCompile`; P2 precomputed PERF `inLoop` + fixed BP-PY-13 alloc. After verify: Mask alloc **~6%**, `firstCodeMatchStart` **~25%** cum, `PyCweScan` **~71%** cum — absolute time down **~5×**, but residual regex in CWE still keeps product/bench **above 100 ms**.

**Trajectory vs plan:** P0+P1+P2 landed product in the **~130–190 ms** band (plan’s P1 band), not yet **&lt;100 ms**.

---

## Gates (prove after each batch)

- [x] Duration benches `benchtime=5s` — `/tmp/goslop-python-pprof/bench-scan-after-5s.txt` (**146.3 ms**/op)
- [x] Fixed-iter confirmation (10× product) — mean **157.5 ms**, median **154.6 ms**, best **133.4 ms**, worst **189.8 ms** (`product-run-python-after-10x.txt`)
- [x] `go tool pprof -top -cum` on post `scan-cpu-after.prof` — Mask gone from top; `firstCodeMatchStart` **48% → 24.8%**; `PyCweScan.Run` **80.1% → 71.4%** (share of a much smaller total)
- [x] Mem: `go tool pprof -top -sample_index=alloc_space` — `pytext.Mask` flat **85.8% → 6.07%** (target was &lt;5%; close)
- [x] Correctness: findings **102**; 38h/56i/0l/8m; top BP-PY-46×32, BP-PY-41×24, CWE-328×14, CWE-90×12, CWE-22×6 — unchanged
- [x] Tests + lint after verify: `go1.26.4 test` python packages green; `go1.26.4 vet` clean; `gofmt -l` clean (fixed `rules_tier_b_resource.go`, `rules_hotpath.go`)
- [ ] Product target: `make run-python` scan line **&lt; 100 ms** (cold `--no-cache`) — **not met** (mean 157.5 / best 133.4)

**Measured (pre → after):** engine scan **728.5 → 146.3 ms**; product wall **~1.1 s → ~157 ms** mean; Go product ~100 ms.

---

## Phase 0 — Measurement harness (done for baseline)

### 0.1 Product Python benches + pprof dumps

Package: `internal/bench`

- [x] Add `BenchmarkPythonScanProfileAll` / `BenchmarkPythonScanAndExport` using `NewRegistryWithLanguages(LanguagePython)` — `make_run_python_bench_test.go`
- [x] Capture CPU + mem profiles — `/tmp/goslop-python-pprof/scan-cpu.prof`, `scan-mem.prof`, `scan-export-*.prof`
- [x] Record tops — `scan-cpu-top-cum.txt`, `scan-mem-alloc-space.txt`, `python-perf-pprof-measurement.md`
- [x] Go compare on same machine — `bench-go-scan-5s-compare.txt`, `product-run-go-3x.txt`

---

## Phase 1 (P0) — highest leverage (CWE mask + first-match) — ~55–75% scan

### 1.1 Cache `pytext.Mask` once per file in CWE facts (~86% alloc_space; est. **~40–60%** wall)

Package: `internal/lang/python/detectors/cwe` (`facts.go`, `common.go`, callers)

- [x] Extend `PyCweFacts` with `Masked string` (and optionally `MaskedLines []pyMaskedLine`) built **once** in `BuildFacts`
- [x] Stop ignoring facts: rule fns use `facts *PyCweFacts` (0 remaining `_ *PyCweFacts`)
- [x] Change `findCalls`, `firstCodeMatchStart`, `pythonFunctions`, `maskedPythonLines`, and direct `pythonCodeMask(unit.Source)` sites to use `facts.Masked` — target **1 Mask/file** (was ~175+)
- [x] Fix nested remask: `matchWithoutDefaultStart` uses `facts.MaskedLines()` (`rules_platform.go`)
- [x] Keep offset stability: mask remains byte-aligned with `unit.Source` (`pytext.Mask` contract) — proven by findings parity
- [~] Re-profile: `pytext.Mask` alloc_space flat **85.8% → 6.07%** (target &lt;5% not quite); Mask CPU off the top; B/op **553 → 24.5 MB** (&lt;100 MB met)
- [x] Proof: `go test ./internal/lang/python/detectors/cwe/ -count=1` + parity findings on pythoncoreengine (**102**)

### 1.2 Fix `firstCodeMatchStart` algorithm (~48% CPU; est. **~20–35%** additional)

Package: `internal/lang/python/detectors/cwe/common.go` (+ `firstMatchStart` in `rules_secrets.go`)

- [x] Replace `FindAllStringIndex(source, -1)` with iterative `FindStringIndex` from advancing offset until masked span is non-blank
- [x] Prefer matching on **cached masked** text (stacks with 1.1) — `firstCodeMatchStartMasked`
- [x] Keep comment/string false-positive filtering semantics identical (`common_test.go` / fixtures)
- [~] Re-profile: `FindAllStringIndex` under helper near-gone from this path; helper cum **48.0% → 24.8%** (target mid/low teens not met — residual regex still dominates)
- [x] Proof: no FN/FP delta on CWE unit + integration matrices + product parity

### 1.3 Memoize `pythonFunctions` once per file (~13% CPU)

Package: `internal/lang/python/detectors/cwe` (`rules_ssrf.go`, facts)

- [x] Compute function spans **once** in `BuildFacts` (`facts.Funcs` / `buildPythonFunctions`); reuse across callers
- [x] Eliminate per-body remask of `fn.body` where full-file Mask is already cached (slice from `Masked`)
- [~] Ensure `findCalls` never remasks when facts hold `Masked` — full-file path uses cache; some fragment/`fragStart&lt;0` fallbacks still remask
- [~] Re-profile: `buildPythonFunctions` cum **~3.1%** (old `pythonFunctions` **13%** gone; target &lt;2% not quite)

---

## Phase 2 (P1) — CWE rule scheduling / gates — push toward ~100–150 ms

### 2.1 SourceIndex gates for ungated catalogue (154/159 ungated; est. **~15–30%**)

Package: `needles.go` + `RegisterRule` call sites

**After:** **65 gated / 94 ungated** (was 5 gated). Inventory below was wired.

**Wire needles that already exist but are ungated:**

- [x] Secrets: CWE-798, 256, 260, 312, 547, 523, 319, 261
- [x] Injection: CWE-90, 91, 93, 94, 88, 117 (`rules_injection.go`)
- [x] Code-dynamic: CWE-749, 829, 695, 214, 215 (`rules_code_dynamic.go`)
- [x] `detectCWE208` / `pyTimingCompareRE` — FN-safe token needles (`password`/`token`/`secret`/…)
- [x] Expand needles for Tier-B `firstCodeMatchStart` rules only with FN-safe tokens (inventory rows gated)
- [x] Re-validate intentional ungated sets (`rules_ssrf.go`, `rules_auth.go`, `rules_path_fs.go`, `rules_info_exposure.go` multi-needle FN comments) — parity held at **102**
- [~] Re-profile: fewer entries into `firstCodeMatchStart` / `findCalls` on non-matching files — helper still **~25%** cum; pack still **~71%**
- [x] Proof: fixture matrices + pythoncoreengine finding parity (**102**, no intentional delta)

#### Formerly-ungated `firstCodeMatchStart` / `firstMatchStart` inventory (now gated)

| File | Rules | Status |
|------|-------|--------|
| `rules_info_exposure.go` | CWE-208 | gated |
| `rules_secrets.go` | CWE-798, 256, 260, 312, 547 (+ CWE-523 via firstMatchStart) | gated |
| `rules_platform.go` | CWE-396, 397 | gated |
| `rules_auth.go` | CWE-613 | gated |
| `rules_crypto.go` | CWE-1392 | gated |
| `rules_web_config.go` | CWE-15, 1052, 1188, 921 | gated |
| `rules_validation.go` | CWE-1230 | gated |
| `rules_tier_b_quality.go` | CWE-1104, 1106, 1220, 1265, 1284, 1285, 1287, 1288, 1322, 1339, 1341 | gated (subset via RegisterRule needles) |
| `rules_tier_b_runtime.go` | CWE-920, 1007, 1021, 1050, 1060, 1071, 1072 | gated |
| `rules_tier_b_resource.go` | CWE-367, 472, 521, 524, 538, 617, 641, 779 | gated |
| `rules_tier_b_access.go` | CWE-66, 76, 179, 182, 184, 323 | gated |

### 2.2 Kill per-call `regexp.MustCompile` on hot paths

Package: tier-b + secrets + web_config

- [x] Hoist inline `MustCompile` in `detectCWE1104`, `detectCWE472`, `detectCWE641`, `detectCWE523` / CWE-295 `check_hostname` to package `var`
- [ ] Soften `pyTimingCompareRE` if still &gt;2% after gates (two-pass: find `==` lines, then identifier check) — not done; CWE-208 gated instead
- [x] Re-profile: `regexp.MustCompile` alloc_space **~0.72%** (compile noise collapsed vs pre)
- [x] Proof: unit tests unchanged (python packages green)

### 2.3 `maskedPythonLines` once on facts (~1–3%)

Package: `rules_platform.go`

- [x] Cache lines on facts after 1.1 (`MaskedLines()`); close nested remask in `matchWithoutDefaultStart`

**P0+P1 validation target:** product wall **~150–250 ms** after P0; **&lt;150 ms** after P1; findings **102**.  
**Observed after P0–P2:** product mean **157.5 ms** / best **133.4 ms**; findings **102**. Mean still slightly above the &lt;150 ms P1 goal.

---

## Phase 3 (P2) — PERF / BP polish (after CWE ~&lt;150 ms) — close &lt;100 ms

PERF and BP already build line views **once per file** (`buildCodeLines`); they do **not** use `pytext.Mask`. Do **not** share Mask across packs (PERF/BP keep string literals for heuristics like `status="pending"`).

### 3.1 Precompute `inLoop` membership (~7.6% cum; TrimSpace-dominated)

Package: `internal/lang/python/detectors/perf/facts.go`

- [x] Precompute loop-membership bitset once in `buildFacts` (`computeInLoop`); rules query O(1) via `lineInLoop`
- [x] Cache trimmed / empty on `codeLine` at build time (`trim` field)
- [x] Cheap gates **before** `inLoop` on hot rules (`.objects.` / delivery needle / ctor hint first) — PY-2/5/3/25
- [x] Re-profile: `perf.inLoop` gone from CPU top (was **7.6%**); PERF pack **~16.5%** cum of smaller total
- [x] Proof: PERF fixture matrix + canary unchanged (perf package tests green; product PERF findings stable in parity)

### 3.2 BP-PY-13 alloc fix (small CPU, large alloc)

Package: `internal/lang/python/detectors/bad_practices`

- [x] Hoist per-line `MustCompile` in `detectBPPY13` to package `var` (`secretExactAssignRe` / `secretNameRe`)
- [x] Avoid up to **4× `ToLower(unit.Source)`** on gate-miss path — single `ToLower` then contains checks
- [x] Re-profile: BPPY13 no longer dominates mem; BP pack **~7.6%** CPU cum; `detectBPPY13` alloc cum **~3.3%**
- [x] Proof: BP unit + integration tests green

### 3.3 PERF residual alloc / domain skips

Package: `internal/lang/python/detectors/perf`

- [ ] Audit hot-path dynamic `MustCompile` (`rules_orm`, `rules_local`, `rules_batch` — e.g. PYPERF10 / PERFPY18)
- [ ] Optional SourceIndex / needle gates so non-Django/hotpath files skip full line walks
- [ ] Optional shared stripped lines across PERF+BP only if still measurable after 3.1–3.2 (small)

### 3.4 Export path (only if scan &lt;80 ms but product wall still &gt;100 ms)

Package: `internal/export` + `BenchmarkPythonScanAndExport`

- [~] Deferred: engine scan **~146 ms** (not &lt;80 ms); product mean **~157 ms** — residual is still detector/regex time, not export-dominated. Revisit only after scan &lt;80 ms.

---

## Phase 4 — Closure gates

- [ ] `BenchmarkPythonScanProfileAll` **&lt; 100 ms**/op (same corpus, `benchtime=5s`) — **146.3 ms** after
- [ ] Product `make run-python` scan line **&lt; 100 ms** mean over ≥10 runs (`--no-cache`) — mean **157.5 ms**
- [x] Findings parity gate recorded (**102**; severity + top rules match baseline)
- [x] `go1.26.4 vet` + `go1.26.4 test` python packages + `gofmt -l` green (verify session)
- [x] Update [`python-perf-pprof-measurement.md`](./python-perf-pprof-measurement.md) with after numbers + new pprof tops
- [ ] Link PR(s) to closed checklist rows

---

## Suggested PR order

| PR | Checklist items | Theme | Status |
|----|-----------------|--------|--------|
| 1 | 0.1 (harness on branch) + 1.1 Mask cache in facts | alloc_space collapse | **implemented / verified** |
| 2 | 1.2 firstCodeMatchStart + 1.3 function memo | CPU regex collapse | **implemented / verified** (`[~]` on &lt;teens cum%) |
| 3 | 2.1–2.3 gates + MustCompile + CWE-208 | catalogue scheduling | **implemented / verified** |
| 4 | 3.1–3.3 PERF inLoop + BP-PY-13 | polish to &lt;100 ms | **3.1–3.2 done**; 3.3 open; &lt;100 ms not met |
| 5 | 3.4 (if needed) + Phase 4 measurement | closure | **open** (3.4 deferred; &lt;100 ms fail) |

---

## Out of scope (this checklist)

- [ ] Changing detector *semantics* / finding counts for speed (unless gated FN risk documented)
- [ ] Sharing `pytext.Mask` across CWE/PERF/BP (PERF/BP intentionally keep literals)
- [ ] Dropping whole-function export default
- [ ] Committing raw `.prof` blobs to git
- [ ] Turning off experimental PERF-PY to “win” the bench without fixing CWE
- [ ] CGO / tree-sitter Python AST (source-only remains)

---

## References

- pprof data: `/tmp/goslop-python-pprof/`
- Measurement write-up: [`python-perf-pprof-measurement.md`](./python-perf-pprof-measurement.md)
- Go precedent: [`../v0.0.1/perf-pprof-optimization-checklist.md`](../v0.0.1/perf-pprof-optimization-checklist.md)
- Bench harness: `internal/bench/make_run_python_bench_test.go`
- Hot code: `internal/lang/python/detectors/cwe/{common.go,facts.go,scan.go,rules_*.go}`
- Mask helper: `internal/lang/python/pytext/mask.go`

---

## Phase 5 — Round-2 pprof → target ~90 ms (2026-08-02 r2 dump)

> **Parent phases:** 0–4 above remain the history; **do not remove**. This phase is the next execution ledger toward **&lt;100 ms / ~90 ms**.  
> **Status:** **verified after 5.1–5.4** — bench **132.9 → 78.5 ms**/op · product mean **147.9 → 93.3 ms** / best **81.0** / worst **116.7** · findings **102** parity · `FindAllStringIndex` **gone**; residual: `backtrack` **~31%** cum, `firstCodeMatchStart` **~18%**, `detectCWE1050` **~8.8%**, PERF/BP share  
> **Goal:** `BenchmarkPythonScanProfileAll` **~90 ms**/op · product `make run-python` mean **~90–100 ms** · findings parity **102** — **bench + product-mean gates met**  
> **Need (met):** ~**32%** further engine cut from r2 (133 → 90) — achieved **~41%** (133 → 78.5); product mean ≤100 — **93.3 ms**

### 5.0 Round-2 measurement (captured)

Artifacts under `/tmp/goslop-python-pprof/`:

| Artifact | Content |
|----------|---------|
| `bench-scan-r2-5s.txt` | `BenchmarkPythonScanProfileAll` **132.9 ms**/op · **24.7 MB** · **70.6k** allocs |
| `scan-cpu-r2.prof` / `scan-cpu-r2-top-cum.txt` | CPU tops |
| `scan-mem-r2.prof` / `scan-mem-r2-alloc-space.txt` | alloc_space tops |
| `scan-cpu-r2-focus.txt` / `scan-cpu-r2-detectCWE.txt` | focused tops |
| `product-run-python-r2-10x.txt` | 10× product wall |
| `bench-scan-r2-after-5s.txt` | after 5.1–5.4: **78.5 ms**/op · **22.8 MB** · **54.2k** allocs |
| `scan-cpu-r2-after.prof` / `scan-cpu-r2-after-top-cum.txt` | after CPU tops |
| `scan-mem-r2-after.prof` / `scan-mem-r2-after-alloc-space.txt` | after alloc_space |
| `product-run-python-r2-after-10x.txt` | after 10× product wall |

**Product wall r2 (10×, export on, `--no-cache`):** mean **147.9 ms** · median **141.8** · best **139.0** · worst **171.5** · findings **102**.

**Product wall after-r2 (10×):** mean **93.3 ms** · median **88.3** · best **81.0** · worst **116.7** · findings **102**.

**vs Phase-4 verify:** bench **146.3 → 132.9 ms** (noise/warmup); product mean **157.5 → 147.9 ms**. Still **well above 100 ms** at r2 baseline.

**vs after-r2:** bench **132.9 → 78.5 ms** (~**1.7×**); product mean **147.9 → 93.3 ms**.

#### r2 CPU hotspots (`-top -cum`)

| Symbol | cum% | Action theme |
|--------|-----:|--------------|
| `regexp.doExecute` | **77.9%** | still owns the scan |
| `cwe.(*PyCweScan).Run` | **71.0%** | CWE pack |
| `FindStringIndex` | **35.3%** | mostly via `firstCodeMatchStart*` |
| `firstCodeMatchStart` / `…Masked` | **24.5%** | match on **unmasked** `source` then filter — backtrack-prone |
| `regexp.backtrack` / `tryBacktrack` | **21.0% / 17.3%** | catastrophic / ambiguous REs |
| `FindAllStringIndex` | **15.0%** | auth + `routeHandlerBodies` |
| `firstMatchStart` | **12.7%** | secrets / web_config / platform wrappers |
| `perf.(*PythonPerfScan).Run` | **16.6%** | secondary after CWE absolute drop |
| `debugEnabledStart` | **9.8%** | heavy `(?im)` DEBUG=True RE (CWE-756/489) |
| `routeHandlerBodies` | **8.7%** | `FindAllStringIndex` on decorator RE; **not memoized on facts** |
| `bad_practices.Run` | **7.4%** | share rose; absolute still modest |
| `BuildFacts` | **6.6%** | Mask + Funcs + Index once/file |
| `detectCWE1050` | **5.8%** | Tier-B runtime |
| `detectCWE915` | **5.7%** | mass assign |
| `detectCWE756` | **5.1%** | via `debugEnabledStart` |
| `detectCWE489` | **4.8%** | via `debugEnabledStart` |
| `detectCWE914` / `detectCWE908` | **4.2% / 3.9%** | residual |

#### r2 alloc_space (no longer Mask-dominated)

| Symbol | flat% |
|--------|------:|
| `bytealg.MakeNoZero` | 23.1% |
| `regexp.(*bitState).reset` | 11.8% |
| `strings.genSplit` | 10.4% |
| `regexp.get` | 9.1% |
| `pytext.Mask` | **7.2%** (was 86% pre-P0) |
| `perf.buildCodeLines` | 5.0% |
| `bp.buildCodeLines` | 4.7% |

**Gate coverage now:** **128 gated / 31 ungated** of 159 CWE rules (was 66/93 at r2; Phase 5.3–5.4).

- [x] Capture r2 CPU + mem profiles + 10× product wall — artifacts above
- [x] Record r2 tops in this section + [`python-perf-pprof-measurement.md`](./python-perf-pprof-measurement.md) (r2 appendix)

---

### 5.1 P0 — match on masked text + kill backtracking (~20–35% engine est.)

Package: `internal/lang/python/detectors/cwe/common.go`, `rules_secrets.go` (`firstMatchStart` → `firstLiteralMatchStart*`)

- [x] Change first-hit helpers to run `FindStringIndex` on **`masked`** (byte-aligned) when the pattern is code-oriented; keep source-based path only where string-literal evidence is required (`firstLiteralMatchStart*`)
- [~] Soften / split worst backtracking REs used via first-match helpers — some prefilters landed; **`backtrack` still ~30.9% cum** (share rose vs r2 21% as CWE absolute fell; absolute ~28→~24 ms)
- [x] Rewrite `debugEnabledStart`: Index gate + masked-line scan + `pyDebugLineRE` (off top; was ~10% via `pyDebugAssignmentRE`)
- [~] Re-profile: `backtrack` cum **21% → 30.9%** (target &lt;5% **not met**); `firstCodeMatchStart` cum **24.5% → 18.2%** (target &lt;12% **not met**)
- [x] Proof: `go1.26.4 test` python packages green + product findings **102**

**Hypothesis:** largest remaining wall cut toward ~100 ms — **wall target met anyway** via 5.1–5.4 combo; backtrack share still the main residual lever.

---

### 5.2 P0 — memoize `routeHandlerBodies` + stop hot `FindAllStringIndex` (~8–15% est.)

Package: `rules_code_dynamic.go`, `rules_auth.go`, consumers in info_exposure / validation / code_dynamic

- [x] Compute route-handler spans **once** lazy on `PyCweFacts.RouteHandlers` (like `Funcs`); `routeHandlerBodies` returns cached slice
- [x] Auth / related: replace `FindAllStringIndex` with masked iterative find / cached route index + needle gates — **`FindAllStringIndex` focus = 0 samples**
- [~] Re-profile: `FindAllStringIndex` cum **15% → 0%** (**met**); `routeHandlerBodies` **8.7% → 4.6%** (target &lt;1% **not met**)
- [x] Proof: auth / code-dynamic / info-exposure fixtures + parity **102**

---

### 5.3 P1 — hot single-rule passes (CWE-1050 / 915 / 756 / 489 / 914 / 908)

Package: tier_b_runtime, mass_assign, web_config, related

- [x] Profile each named detect after 5.1–5.2; rewrite top survivors with needle + linear / line scan
  - CWE-1050: line-oriented `openInsideShortLoopStart` (no multi-line catastrophic RE)
  - CWE-908/910: iterative `FindStringSubmatchIndex` + `= None` / `.close(` Contains
  - CWE-915/914: Contains prefilters before REs; match on masked; RegisterRule gates
  - CWE-756/489: gated (debugEnabledStart rewrite owned by 5.1)
- [x] Ensure FN-safe gates present and tight enough that cold files skip these bodies
- [~] Re-profile: `detectCWE1050` still **8.8%** cum; `detectCWE908` **3.5%** (target no single `detectCWE*` above **~2%** — **not met** for 1050/908)
- [x] Proof: unit + integration matrices; findings **102** — `go1.26.4 test ./internal/lang/python/...` green

---

### 5.4 P1 — more SourceIndex gates on remaining 93 ungated rules

Package: `needles.go` + `RegisterRule` sites still without gates

- [x] Inventory ungated rules that still enter `firstCodeMatchStart` / `FindAll*` / heavy helpers on pythoncoreengine
- [x] Add FN-safe any-of gates; expand `pyCweNeedles` only with required tokens
  - Gate coverage: **65 → 128** gated / **31** intentionally ungated (of 159)
- [x] Leave intentionally ungated only where multi-shape FN risk is documented (ssrf / path_fs / xml / info_exposure / auth CORS+cookie / quality structural 1108/1121/1124 / validation 1289)
- [x] Re-profile: `PyCweScan.Run` cum **71.0% → 51.5%**; absolute ~**94 → ~40 ms**/op-equivalent (bench 133→78.5)
- [x] Proof: fixtures + product parity — cwe + integration python green

---

### 5.5 P2 — PERF / BP share after CWE &lt;~100 ms engine

Package: `perf/`, `bad_practices/`

- [ ] After 5.1–5.4, if `PythonPerfScan` still ≥10% cum: domain needle skip for non-ORM/hotpath files (extends open **3.3**) — **still 29.0% cum** after-r2
- [ ] If BP share still ≥5% absolute: audit `buildCodeLines` / BPPY4 alloc — **`PythonBadPracticeScan` 12.8% cum**; `detectBPPY4` still in alloc_space
- [ ] Optional: avoid duplicate `strings.Split` across PERF+BP only if measured (still out of scope to share Mask)

---

### 5.6 Closure gates for ~90 ms

- [x] `BenchmarkPythonScanProfileAll` **≤ 90 ms**/op — **78.5 ms**/op · 22.8 MB · 54.2k allocs (`bench-scan-r2-after-5s.txt`)
- [~] Product `make run-python` mean **≤ 100 ms** over ≥10 runs — **mean 93.3 / median 88.3 / best 81.0 / worst 116.7** (`product-run-python-r2-after-10x.txt`); stretch mean **~90 ms** close but not quite
- [x] Findings parity **102** — 38h/56i/0l/8m · BP-PY-46×32, BP-PY-41×24, CWE-328×14, CWE-90×12, CWE-22×6
- [~] `go tool pprof -top -cum` on post profile: `FindAllStringIndex` **gone**; `firstCodeMatchStart` still **18.2%**; `backtrack` still **30.9%** (not cleared as top residual helpers)
- [x] Tests + vet + gofmt: `go1.26.4 test` python packages green; `go1.26.4 vet ./internal/lang/python/...` clean; `gofmt -l …/cwe` empty
- [x] Append after-r2 numbers to [`python-perf-pprof-measurement.md`](./python-perf-pprof-measurement.md); keep this Phase 5 history

### Reproduce r2 dump / after-r2 verify

```sh
export CGO_ENABLED=0 PATH="$HOME/go/bin:$PATH"
export GOSLOP_BENCH_PYTHON_SCAN_PATH=/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/pythoncoreengine
mkdir -p /tmp/goslop-python-pprof

go1.26.4 test -run='^$' -bench=BenchmarkPythonScanProfileAll -benchmem -benchtime=5s \
  -cpuprofile=/tmp/goslop-python-pprof/scan-cpu-r2-after.prof \
  -memprofile=/tmp/goslop-python-pprof/scan-mem-r2-after.prof \
  ./internal/bench/ | tee /tmp/goslop-python-pprof/bench-scan-r2-after-5s.txt

go1.26.4 tool pprof -top -cum /tmp/goslop-python-pprof/scan-cpu-r2-after.prof
go1.26.4 tool pprof -top -sample_index=alloc_space /tmp/goslop-python-pprof/scan-mem-r2-after.prof
```

### Suggested next PR order (Phase 5)

| PR | Items | Theme | Status |
|----|-------|--------|--------|
| 6 | 5.1 masked first-match + debugEnabled soften | kill backtrack | **implemented** (backtrack/firstCode targets **[~]**) |
| 7 | 5.2 routeHandlerBodies memo + auth FindAll | FindAll collapse | **implemented** (`FindAll` **[x]**; route **[~]** 4.6%) |
| 8 | 5.3–5.4 hot rules + more gates | catalogue polish | **implemented** (1050 still **8.8%** cum) |
| 9 | 5.5–5.6 PERF/BP + closure ~90 ms | finish | **partial** — timing closure **[x]/; 5.5 PERF/BP **open** |
