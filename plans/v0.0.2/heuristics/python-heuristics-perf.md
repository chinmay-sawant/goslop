# v0.0.2 — Python PERF heuristics (22/22 experimental detectors implemented)

> **Parent:** `plans/v0.0.2/python-heuristics.md` — epic [#51](https://github.com/chinmay-sawant/goslop/issues/51); issue body `plans/PR/v0.0.2/issue-python-perf-heuristics-body.md`  
> **Issue:** [#54](https://github.com/chinmay-sawant/goslop/issues/54) — python(perf): seed PERF catalogue and implement performance heuristics  
> **Status:** **complete** — experimental detectors + fixtures; reference-corpus canary and maturity decision recorded 2026-08-01 (stay experimental; no pack promotion)
> **Estimated effort:** closed for #54 scope; future tier promotion needs a fresh canary
> **Ledger rule:** canonical evidence is `pref-plans/PERF-PY-IMPLEMENTATION-CHECKLIST.md` + `pref-plans/PERF-PY-CANARY-2026-08-01.md`; do not promote maturity without a fresh canary.

---

## Overview

Go ships a full **PERF** surface:

| Layer | Path | Scale |
|-------|------|------:|
| Catalogue JSON | `ruleset/golang/chunks/perf-*.json` | **242** rules (≤50/file; range names) |
| Detectors | `internal/lang/go/detectors/perf/` | **239** registered heuristics |
| Domains | `…/perf/registry/*.toml` | loop_allocations, parsing_in_loops, request_path, general_perf, gin_framework, data_access, protocols |

Python today:

| Layer | Path | Status |
|-------|------|--------|
| CWE catalogue | `ruleset/python/chunks/cwe-*.json` | present (~344) |
| BP catalogue | `ruleset/python/bad-practices.json` | present (50 `BP-PY-*`) |
| **PERF catalogue** | `ruleset/python/chunks/perf-py-001-014.json` + `perf-py-015-022.json` | **seeded** (22 Python-only rules) |
| PERF detectors | `internal/lang/python/detectors/perf/` | **22** experimental `PERF-PY-*` rules; Python remains opt-in |

**#54 is closed experimentally.** Seeded IDs have source-only detection, paired fixtures, repository validation, and a 2026-08-01 reference-corpus canary with an explicit “stay experimental” maturity decision.

---

## Executive Summary

### Implemented detector boundary

1. **Python PERF remains source-only** — the 22 rules cover the defensible static subset; no Go AST/CGO port was introduced.
2. **Experimental maturity is intentional** — fixture parity proves detector behavior, not production signal.
3. **Review-level inference is bounded** — rules 15, 16, and 20 require visible local declarations; rules 17 and 22 require explicit production evidence.
4. **BP ownership remains separate** — timeout and blocking-async patterns stay with BP-PY; the PERF rules use cost-specific conditions.

### What “unblocked” means

A schema-valid `ruleset/python/chunks/perf-py-001-014.json` + `perf-py-015-022.json` set exists (≤50 entries per file, `PERF-PY-*` IDs, `applicable_to` including `python`, `python_relevance`, Python-oriented `detection_notes`), catalogue tests pass, and README documents the layout. The catalogue gate is met; detector rows can now be promoted to active `[ ]` work.

### Non-goals (remain out of scope even after unblock)

- Porting all ~239 / 242 Go PERF rules  
- Micro-benchmark CI gates  
- CWE / BP family work (siblings #52 / #53) except ownership coordination  
- Tree-sitter / CGO Python parsers as a hard requirement for the first PERF batch  

---

## Phase 1: Catalogue seed (complete)

> **Gate for all later phases.** No detector code should land before this phase has real files and parse tests.

### 1.1 Evidence of absence (baseline)

- [x] Confirm the pre-seed branch had no Python PERF chunks — **evidence:** no `ruleset/python/chunks/perf-*.json` before this change (2026-07-31)  
- [x] Confirm Python plugin still has zero PERF detectors — **evidence:** `internal/lang/python/plugin.go` + `internal/lang/python/detectors/all.go` contain no PERF detector (2026-07-31); this remains the next implementation slice  

### 1.2 Layout + naming (match Go / Python CWE convention)

- [x] Add range-named Python chunk files — **evidence:** `perf-py-001-014.json` and `perf-py-015-022.json` (2026-07-31)  
- [x] Enforce **≤50** rule objects per file — **evidence:** catalogue test checks the limit; seed contains 14 + 8 entries  
- [x] Use Python-scoped map keys / IDs — **evidence:** `PERF-PY-1` … `PERF-PY-22`; documented in `ruleset/python/README.md`  

### 1.3 Schema (align with existing Python CWE chunks)

Per-rule fields (mirror golang PERF + Python CWE shape):

`id`, `name`, `original_description`, `description`, `detection_notes`,  
`category` (`Performance`), `status`, `weakness_abstraction`, `python_relevance`, `applicable_to`

- [x] Every seeded rule has `applicable_to` including `"python"` — **evidence:** `TestPythonPERFSeedChunks`  
- [x] Every seeded rule has `python_relevance` (`High` / `Medium` / `Low`) — **evidence:** `TestPythonPERFSeedChunks`  
- [x] Descriptions and notes are Python/framework oriented — **evidence:** catalogue review; no Go-only fields or prose  
- [x] Do **not** bulk-copy Go PERF — **evidence:** independently authored Python seed with `PERF-PY-*` namespace  

### 1.4 Initial seed themes (examples from issue #54 — finalize at seed time)

Target an **initial** set only (not Go parity). Candidate domains (IDs assigned when seeding):

| Theme | Python-oriented signals (illustrative) | Go domain lesson (do not copy IDs) |
|-------|----------------------------------------|-------------------------------------|
| Async / request-path blocking | `time.sleep` / sync I/O / `requests.*` inside `async def` | request_path / general_perf |
| Expensive ops in loops | `re.compile` in loop; open/read file in loop; repeated JSON parse | loop_allocations / parsing_in_loops |
| N+1 / ORM cost | Django ORM access-in-loop; missing `select_related` / `prefetch_related`; SQLAlchemy lazy load in loop | data_access |
| Unbounded growth | unbounded list append / unbounded queryset without pagination on hot path | data_access / general_perf |
| HTTP client timeouts | `requests.get` / `httpx` / `urllib` without timeout (if not BP-owned) | general_perf stdlib_misuse |

- [x] Draft seed rule list with stable `PERF-PY-*` IDs — **evidence:** `PERF-PY-1` … `PERF-PY-22` across the two chunks  
- [x] Coordinate ownership with BP — **evidence:** async blocking and HTTP timeout remain BP-owned and are explicitly excluded/documented  
- [x] Write seed chunks with the expanded static set — **evidence:** `ruleset/python/chunks/perf-py-001-014.json` and `perf-py-015-022.json`  

### 1.5 Docs + catalogue validation

- [x] Update `ruleset/python/README.md` PERF layout section — **evidence:** Performance section documents the seed and namespace  
- [x] Extend `ruleset/python/catalogue_test.go` for PERF chunks and ≤50/file — **evidence:** `TestPythonPERFSeedChunks`  
- [x] Validate catalogue seed — **evidence:** `go test ./ruleset/python/` passes (2026-07-31)  

### 1.6 Phase 1 closure gate

- [x] Catalogue seed documented and unblocked — **evidence:** both PERF-PY chunks, README, catalogue test, and parent rollup updated (2026-07-31)  

### 1.7 Evaluation coverage map

The expanded seed maps the source evaluation to static, source-verifiable PERF ownership:

| Evaluation rows | PERF-PY rules | Decision |
|---|---|---|
| FA-1, FA-4, FA-7 | 1, 7, 14 | seeded |
| FA-2, FA-3, FA-5, FA-6 | 15, 16, 17, 18 | seeded |
| DJ-1, DJ-2, DJ-3, DJ-5, DJ-6 | 2, 3, 4, 12, 13, 14 | seeded |
| DJ-4, DJ-7, DJ-8 | 19, 22, 20 | seeded |
| FL-1, FL-2, FL-3, FL-4, FL-6, FL-8 | 5, 6, 8, 9, 10, 11 | seeded |
| FL-5, FL-7 | 16, 21, 17, 22 | seeded |
| FA-8 | — | correctness/operability; remains outside PERF |
| XC-1, XC-3, XC-4, XC-5 | — | benchmark, deployment, observability, or security work; not static PERF |
| XC-2 | 17 | absorbed into the database lifecycle rule |

**Proof when unblocked (record outcomes then):**

```sh
go test ./ruleset/python/
# list chunk files: ruleset/python/chunks/perf-*.json
# assert each file has ≤50 entries
```

---

## Phase 2: Detectors (implemented; canary pending)

> **Implemented:** `internal/lang/python/detectors/perf/`; `make lint`, `make test`, and `CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop` passed on 2026-07-31.

> **Canonical implementation checklist:** `plans/v0.0.2/heuristics/pref-plans/PERF-PY-IMPLEMENTATION-CHECKLIST.md`. Keep all 22 detector, fixture, canary, and quality-gate status there; this parent remains the #54 rollup.

### 2.1 Package / registration skeleton

- [x] Python PERF package, `RegisterRule` contract, metadata, plugin wire-up, and `PERF-PY-*` pack/discovery identity — **evidence:** canonical checklist §§1.1–1.2 (2026-07-31).

### 2.2 Heuristic implementation (seeded IDs only)

Driven by catalogue `detection_notes`; prefer pure-Go source-pattern / light-parse heuristics consistent with sibling CWE/BP approach.

- [x] Detector implementation and rule-specific suppressions for all 22 seeded IDs — **evidence:** canonical checklist §§2–4 and Python PERF matrix (2026-07-31).

### 2.3 Fixtures

- [x] Per-rule vulnerable/safe fixtures and the Python PERF fixture matrix — **evidence:** 22 paired files and `go test ./tests/integration/python -run TestPythonPERF -count=1` (2026-07-31).

### 2.4 Integration smoke

- [x] Python-enabled detector emission and unchanged default Go-only behavior — **evidence:** Python integration matrix plus `internal/lang/python/plugin_test.go` (2026-07-31); standalone binary smoke remains pending.

### 2.5 Phase 2 validation (when active)

```sh
gofmt -w <changed Go files>
make lint
make test
CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop
# smoke with languages including python against PERF fixtures
```

- [x] Record implementation-branch `make lint` and `make test` outcomes — **evidence:** both passed 2026-07-31; canonical checklist §5.3.

---

## Phase 3: Docs, mapping, closure (pending detector implementation)

- [~] README/maturity mapping, parent rollup, PR body, and #54 closure evidence — **moved to canonical atomic rows:** `pref-plans/PERF-PY-IMPLEMENTATION-CHECKLIST.md` §5.5; **owner:** epic #51 / issue #54.  

---

## Catalogue gate — completed

The former unblock backlog is complete: the Python-scoped `PERF-PY-1` through `PERF-PY-22` catalogue, ownership notes, and ruleset documentation now exist. The active implementation sequence is deliberately maintained only in the canonical checklist: `pref-plans/PERF-PY-IMPLEMENTATION-CHECKLIST.md`.

---

## Catalogue-gate record

The following catalogue criteria were required before opening the canonical implementation checklist:

| # | Criterion | Proof |
|---|-----------|-------|
| 1 | At least one `ruleset/python/chunks/perf-py-*.json` exists | `ls ruleset/python/chunks/perf-py-*.json` |
| 2 | Each file ≤50 rules; Python-scoped range names | `go test ./ruleset/python/` |
| 3 | Rules use `PERF-PY-*`, `applicable_to` ⊇ `python`, `python_relevance`, Python `detection_notes` | sample + test |
| 4 | Catalogue tests pass | `go test ./ruleset/python/` |
| 5 | README documents PERF layout | `ruleset/python/README.md` |
| 6 | Explicit decision logged: “catalogue gate met; create canonical implementation checklist” | this ledger + `pref-plans/PERF-PY-IMPLEMENTATION-CHECKLIST.md` |

The catalogue gate is met. Detector work is planned in the canonical checklist and no detector row is complete from catalogue evidence alone.

---

## Dependencies

| Depends on | Note |
|------------|------|
| Epic #51 | Parent; sibling streams #52 CWE, #53 BP |
| Foundation #39 / PR #50 | Python plugin stub, `languages` config, `ruleset/python/` tree |
| `ruleset/python/chunks/` CWE layout | Naming + ≤50/file precedent for PERF files |
| Go reference (read-only) | `ruleset/golang/chunks/perf-*.json`, `internal/lang/go/detectors/perf/` |
| BP ownership (#53) | Coordinate N+1 / timeout themes to avoid double findings |
| Reuse policy | `plans/v0.0.2/ruleset-reuse-audit.md` — do not bulk-copy Go PERF |

### Does **not** block (siblings may proceed)

- #52 CWE heuristics  
- #53 BP heuristics  
- Epic #51 partial ship of CWE/BP without PERF  

---

## Success criteria (issue #54 — catalogue criteria met; detector criteria pending)

- [x] Python PERF JSON chunks exist (≤50 rules/file, Python-scoped range names) — **evidence:** two chunks, 22 rules; `go test ./ruleset/python/` passes (2026-07-31)  
- [~] Initial PERF heuristics registered and fixture-tested — **reason:** detector implementation has not started; **canonical next gate:** `pref-plans/PERF-PY-IMPLEMENTATION-CHECKLIST.md` §§1–5  
- [~] `languages = ["python"]` can emit `PERF-PY-*` findings — **reason:** detector implementation has not started; **canonical next gate:** `pref-plans/PERF-PY-IMPLEMENTATION-CHECKLIST.md` §§1–5  
- [x] README / mapping notes updated for PERF layout — **evidence:** `ruleset/python/README.md` and this ledger document `PERF-PY-*`  
- [~] `make lint` + `make test` green on implement branch — **reason:** no detector implementation branch yet; **canonical next gate:** `pref-plans/PERF-PY-IMPLEMENTATION-CHECKLIST.md` §5.3  

---

## References

- Issue body: `plans/PR/v0.0.2/issue-python-perf-heuristics-body.md`  
- Epic body: `plans/PR/v0.0.2/issue-epic-python-heuristics-body.md`  
- Parent ledger: `plans/v0.0.2/python-heuristics.md`  
- Foundation: `plans/v0.0.2/python-support.md`, `ruleset/python/README.md`  
- Reuse audit (Go PERF = Class C): `plans/v0.0.2/ruleset-reuse-audit.md`  
- Go PERF catalogue: `ruleset/golang/chunks/perf-001-050.json` … `perf-232-241.json`  
- Go PERF detectors: `internal/lang/go/detectors/perf/` (+ `registry/*.toml`)  
- Skill: `plans/skills/phase-wise-checklist/SKILLS.md`  
