# v0.0.2 — Python precision hardening checklist

> **Parent:** [python-heuristics.md](./python-heuristics.md) — CWE / BP / PERF ledger  
> **Related:** [pref-plans/PERF-PY-IMPLEMENTATION-CHECKLIST.md](./pref-plans/PERF-PY-IMPLEMENTATION-CHECKLIST.md), [python-heuristics-bp.md](./python-heuristics-bp.md), [python-heuristics-perf.md](./python-heuristics-perf.md)  
> **Trigger:** Luna five-pass live review (2026-08-01) — Python PERF/BP heuristics useful but not yet generally correct or project-agnostic  
> **Status:** **complete** — Phases 0–4 closed 2026-08-01 (canary + maturity decision: stay experimental; no pack promotion)  
> **Scope:** Python-only (`PERF-PY-*`, `BP-PY-*`, Python plugin/facts). Go detectors out of scope.  
> **Ledger rule:** mark `[x]` only after detector + fixture/unit proof is green. Record `go test` evidence on phase gates. Do not promote rules into `recommended`/`perf` packs from this checklist alone.

---

## Overview

Validation baseline (keep green while hardening):

| Suite | Evidence |
|-------|----------|
| Python PERF registration + fixtures | 30/30 fixture-backed (integration matrix) |
| Python BP registration + fixtures | 50/50 registered; BP fixture matrix green |
| Slim multi-file corpus | `tests/fixtures/python/corpus/` + `TestPythonSlimCorpusExpectedFindings` |
| Full unit | `go test ./... -count=1` green (2026-08-01) |

Product stance (do not reverse without canary):

- Keep PERF-PY **experimental / review-oriented**; not in `PerfTierSRules` / `PerfTierARules`.
- Do **not** treat all BP-PY rules as blocking defaults.
- Keep high-signal direct checks: `shell=True`, unsafe `yaml.load`, dynamic `eval`/`exec`, explicit async `time.sleep`.

### Luna claim disposition

| # | Claim | Disposition |
|---|-------|-------------|
| 1 | PERF-PY missing from `recommended`/`perf` | **Correct** — intentional; contract test added |
| 2a | PERF-PY-20 requires `pk=` | **Incorrect** — detector suppresses on `pk=`/`id=`; regression test added |
| 2b | Index check ignores columns | **Correct** — column-aware covering-prefix checks shipped |
| 3 | PERF-PY-6/9/14 function-wide joins | **Correct** — name/key binding shipped |
| 4 | FastAPI/Flask gates too broad | **Correct** — gates + `send_from_directory` policy fixed |
| 5 | BP-PY-10 / BP-PY-6 trust/alias blind | **Correct** — narrowed to untrusted/runtime-validation signals |
| 6 | Source-only limits reliability | **Correct** — shared `pytext.Mask` + triple-quote blanking in PERF facts |

---

## Executive Summary

```text
Phase 0: Contract & documentation freeze
    └─ Phase 1: P0 precision (association, index, framework, security FP)
             └─ Phase 2: P1 noise reduction + catalogue honesty + maturity labels
                      └─ Phase 3: Shared Python facts / string masking
                               └─ Phase 4: Corpus, benchmarks, promotion gates
```

| Phase | Goal | Depends | Status |
|------:|------|---------|--------|
| 0 | Lock stance, profile contract tests, stale comment cleanup | — | done |
| 1 | P0 correctness / FP fixes that poison blocking use | 0 | done |
| 2 | P1 high-noise PERF + SA honesty + BP maturity alignment | 1 | done |
| 3 | Shared facts bag + triple-quote masking (still source-only OK) | 1–2 | done (light) |
| 4 | Multi-file corpus, microbenchmarks, canary before pack promotion | 3 | done |

---

## Phase 0: Contract and documentation freeze

### 0.1 Product stance

- [x] **PERF-PY stays experimental** — not added to `PerfTierSRules` / `PerfTierARules`.
- [x] **Profile contract test** — `ProfileRecommended` / `ProfilePerf` deny `PERF-PY-1` / `PERF-PY-20` (`internal/core/profile_contract_test.go`).
- [x] **Explicit maturity on PERF-PY metadata** — `MaturityExperimental` in `perfMeta`; asserted in `metadata_test.go`.
- [x] **Stale plugin comment** — `plugin.go` documents shipped experimental PERF-PY.
- [x] **Parent ledger pointer** — linked from `python-heuristics.md`.

### 0.2 Phase 0 gate

- [x] Profile contract test green; plugin comment accurate; parent ledger updated.

---

## Phase 1: P0 precision fixes

### 1.1 PERF association binding (PERF-PY-6 / 9 / 14)

- [x] **PERF-PY-6** — bind `.first()` result; require `name.status =` / `name.save(`.
- [x] **PERF-PY-9** — require `json.dumps(payload)` (or dump assign) to appear in create/add args.
- [x] **PERF-PY-14** — require shared idempotency/key token (and model when parseable).
- [x] **Adversarial fixtures/tests** — unrelated status, dump-for-log, lookup-A/create-B cases.

### 1.2 PERF index column matching (PERF-PY-15 / 16 / 20)

- [x] **Column-aware index helper** — `models.Index(fields=[...])` / `db_index=True`; covering prefix.
- [x] **PERF-PY-20 `pk=` regression** — miss `filter(pk=...)`; hit unrelated `db_index` on other column.
- [x] **FN-safe fallback** — unparseable `__table_args__` / opaque Index markers suppress.
- [x] **Catalogue/notes honesty** — Django-focused notes; SQLAlchemy removed from applicable_to for 15/16/20.

### 1.3 FastAPI framework gates (BP-PY-29…32)

- [x] **Tighten `looksFastAPIish`** — require fastapi/FastAPI/APIRouter/starlette markers.
- [x] **BP-PY-29** — `global`/`nonlocal` requires `routeOrDep`.
- [x] **Negative tests** — Flask `@app.route` does not arm BP-PY-29.
- [x] **BP-PY-32** — dropped `FileResponse ⇒ fastapiish` tautology.

### 1.4 Flask path rule (BP-PY-20)

- [x] **`send_from_directory` policy** — miss fixed-root + request name; hit request-derived root / `send_file(request…)`.
- [x] **Update fixtures / batch notes** — unit test + `batch-02-flask.md` + catalogue notes.

### 1.5 Security / assert BP precision (BP-PY-6 / 10)

- [x] **BP-PY-6** — only request/authz/path/CLI-ish asserts; invariant asserts missed.
- [x] **BP-PY-10** — untrusted-arg heuristic; miss `LOCAL_CACHE_BYTES` / `cache_*` style.
- [x] **Keep blockers** — BP-PY-8/9/11/12 unchanged.

### 1.6 Phase 1 gate

- [x] Focused package tests + Python PERF/BP matrices green.
- [x] `go test ./internal/lang/python/... ./tests/integration/python/... ./internal/core/ -count=1` green.
- [x] **Evidence (2026-08-01):** `go test ./... -count=1` green.

---

## Phase 2: P1 noise reduction and catalogue honesty

### 2.1 High-noise PERF rules

- [x] **PERF-PY-25** — denylist lightweight ctors; require lambda or alloc context.
- [x] **PERF-PY-27** — load must be in-loop or in a function directly called from a loop (fixture updated).
- [x] **PERF-PY-5** — loop iterable must come from a claim assignment/call.
- [x] **PERF-PY-21** — drop tautological `"delete"` body needle; require purge/cleanup/retention/expire/maintenance.
- [x] **PERF-PY-2** — require loop variable appears in the ORM call line.

### 2.2 SQLAlchemy / framework honesty

- [x] **PERF-PY-15/16/20 SA** — removed from `applicable_to` / notes until implemented.

### 2.3 BP maturity alignment

- [x] **Document style-pack coupling** — style `BP-*` still includes `BP-PY-*`; BP-PY-6/10/20/29 precision improved for safer advisory/blocking use.
- [x] **Optional gate hygiene** — FastAPI/Flask gates tightened as above.

### 2.4 Phase 2 gate

- [x] Adversarial/unit fixtures green; SA honesty landed; BP docs updated.

---

## Phase 3: Shared facts and masking

### 3.1 Reliability infrastructure (still source-only OK)

- [x] **Unify string/comment masking** — `internal/lang/python/pytext.Mask`; CWE `pythonCodeMask` delegates to it.
- [x] **PERF triple-quote blanking** — `stripPyLineForFacts` blanks docstring bodies while keeping ordinary string keywords (`pending`, etc.).
- [~] **Shared facts bag (imports/aliases/spans)** — not a full import-graph yet; facts still one per `Run` with improved line prep. Full alias graph deferred.
- [~] **Optional light alias resolution** — not shipped for BP-PY-10; direct `pickle.loads(` remains the v0 seam.

### 3.2 Phase 3 gate

- [x] Masking unit tests green (`pytext`, PERF `facts_test`).
- [~] Full shared facts bag beyond line masking — follow-up.

---

## Phase 4: Corpus, benchmarks, promotion

### 4.1 Evidence before maturity promotion

- [x] **In-repo slim multi-file corpus** — Django/Flask/FastAPI/SQLAlchemy under `tests/fixtures/python/corpus/` with `expected.json` + integration test.
- [x] **External canary** — three eval apps scanned; report + FP triage in [pref-plans/PERF-PY-CANARY-2026-08-01.md](./pref-plans/PERF-PY-CANARY-2026-08-01.md) (17 PERF-PY findings; no project-wide spray).
- [x] **Detector microbenchmarks** — `BenchmarkPythonPerfScanSynthetic` in `perf/bench_test.go`.
- [x] **Promotion rule (decision)** — **defer** `PerfTierSPY` / `PerfTierAPY`; keep all PERF-PY experimental; recommended/perf packs unchanged (contract test enforces deny).

### 4.2 Phase 4 gate

- [x] Canary FP table recorded; no PERF-PY in recommended/perf; maturity stays experimental (2026-08-01).

---

## Safe keepers (do not “fix” into silence)

| Rule | Stance |
|------|--------|
| BP-PY-8 | Keep blocking — `shell=True` |
| BP-PY-9 | Keep blocking — `os.system` / `os.popen` |
| BP-PY-11 | Keep blocking — unsafe `yaml.load` |
| BP-PY-12 | Keep blocking — dynamic `eval`/`exec` |
| PERF-PY-7 / 10 / 18 / 28 / 30 | Keep — tight predicates |

---

## Proof commands

```bash
go test ./internal/lang/python/... -count=1
go test ./tests/integration/python/... -count=1
go test ./... -count=1
go test ./internal/lang/python/detectors/perf/ -bench=BenchmarkPythonPerfScanSynthetic -benchmem
```

---

## Progress log

| Date | Phase | Notes |
|------|-------|-------|
| 2026-08-01 | — | Checklist created from Luna review + explore-agent validation |
| 2026-08-01 | 0–4 | Implemented P0/P1 precision, pytext mask, slim corpus, microbench; `go test ./... -count=1` green. External canary / pack promotion still open. |
| 2026-08-01 | 4 | External canary on django/fastapi/flask eval apps; FP triage + maturity defer recorded; CLI smoke PERF-PY-6; neighborhood ownership test; checklist closed. |
