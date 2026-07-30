# v0.0.2 — Python PERF heuristics (DEFERRED)

> **Parent:** `plans/v0.0.2/python-heuristics.md` — epic [#51](https://github.com/chinmay-sawant/goslop/issues/51); issue body `plans/PR/v0.0.2/issue-python-perf-heuristics-body.md`  
> **Issue:** [#54](https://github.com/chinmay-sawant/goslop/issues/54) — python(perf): seed PERF catalogue and implement performance heuristics  
> **Status:** deferred  
> **Estimated effort:** medium–large once unblocked (catalogue seed + initial detector batch; not full Go PERF parity)  
> **Ledger rule:** this file is the **canonical execution ledger for #54**. Every implementation row is `[~]` until the catalogue gate is met. Do not flip rows to `[ ]` / `[x]` without evidence.

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
| **PERF catalogue** | `ruleset/python/chunks/perf-*.json` | **absent** (confirmed) |
| PERF detectors | under `internal/lang/python` | **none** (plugin returns empty `Detectors()`) |

**#54 is intentionally deferred.** Seeding catalogue JSON is the first gate; detector implementation must not start as active work until that gate clears.

---

## Executive Summary

### Why deferred

1. **No Python PERF catalogue** — `ruleset/python/chunks/perf-*` does not exist. Reuse audit (`plans/v0.0.2/ruleset-reuse-audit.md`) classifies Go PERF as **Class C — Go-only mechanics**; bulk-copy is out of policy.  
2. **Issue body order** — #54 explicitly sequences: (1) seed JSON, (2) implement heuristics from `detection_notes`.  
3. **Epic priority** — parent ledger prioritizes CWE (#52) then BP (#53); PERF stays backlog until seeds land.  
4. **Avoid double-firing** — several candidate PERF themes (N+1, HTTP timeouts) may overlap BP; ownership must be decided in catalogue notes before detectors.

### What “unblocked” means

A minimal, schema-valid `ruleset/python/chunks/perf-*.json` set exists (golang-style range names, ≤50 entries/file, `applicable_to: ["python"]`, `python_relevance`, Python-oriented `detection_notes`), catalogue tests pass, and README/mapping notes document the layout. **Only then** promote detector rows from `[~]` to active `[ ]`.

### Non-goals (remain out of scope even after unblock)

- Porting all ~239 / 242 Go PERF rules  
- Micro-benchmark CI gates  
- CWE / BP family work (siblings #52 / #53) except ownership coordination  
- Tree-sitter / CGO Python parsers as a hard requirement for the first PERF batch  

---

## Phase 1: Catalogue seed (deferred)

> **Gate for all later phases.** No detector code should land before this phase has real files and parse tests.

### 1.1 Evidence of absence (baseline)

- [~] Confirm no `ruleset/python/chunks/perf-*.json` on the implementation branch before seeding — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  
- [~] Confirm Python plugin still has zero PERF detectors (`internal/lang/python/plugin.go` empty `Detectors` / `NewDetectors`) — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  

### 1.2 Layout + naming (match Go / Python CWE convention)

- [~] Add range-named chunk files under `ruleset/python/chunks/` (e.g. `perf-001-050.json`, …) — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  
- [~] Enforce **≤50** rule objects per file — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  
- [~] Use map keys / IDs as `PERF-*` (zero-padded or consistent with Go catalogue style chosen for Python; document choice in README) — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  

### 1.3 Schema (align with existing Python CWE chunks)

Per-rule fields (mirror golang PERF + Python CWE shape):

`id`, `name`, `original_description`, `description`, `detection_notes`,  
`category` (`Performance`), `status`, `weakness_abstraction`, `python_relevance`, `applicable_to`

- [~] Every seeded rule has `applicable_to` including `"python"` (not Go tags alone) — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  
- [~] Every seeded rule has `python_relevance` (`High` / `Medium` / `Low`) — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  
- [~] `description` / `detection_notes` rewritten for CPython / asyncio / Django / Flask / FastAPI / DB-API — **not** gin/gorm/goroutine prose — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  
- [~] Do **not** bulk-copy `ruleset/golang/chunks/perf-*.json` (reuse audit Class C) — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  

### 1.4 Initial seed themes (examples from issue #54 — finalize at seed time)

Target an **initial** set only (not Go parity). Candidate domains (IDs assigned when seeding):

| Theme | Python-oriented signals (illustrative) | Go domain lesson (do not copy IDs) |
|-------|----------------------------------------|-------------------------------------|
| Async / request-path blocking | `time.sleep` / sync I/O / `requests.*` inside `async def` | request_path / general_perf |
| Expensive ops in loops | `re.compile` in loop; open/read file in loop; repeated JSON parse | loop_allocations / parsing_in_loops |
| N+1 / ORM cost | Django ORM access-in-loop; missing `select_related` / `prefetch_related`; SQLAlchemy lazy load in loop | data_access |
| Unbounded growth | unbounded list append / unbounded queryset without pagination on hot path | data_access / general_perf |
| HTTP client timeouts | `requests.get` / `httpx` / `urllib` without timeout (if not BP-owned) | general_perf stdlib_misuse |

- [~] Draft seed rule list (theme → provisional `PERF-*` ID → `detection_notes` sketch) — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  
- [~] Coordinate N+1 and HTTP-timeout ownership with BP (`BP-PY-*` in `bad-practices.json`): prefer PERF for **cost**, BP for **style** when already listed; document in rule notes to avoid double-firing — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  
- [~] Write first chunk file(s) with the initial seed set only — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  

### 1.5 Docs + catalogue validation

- [~] Update `ruleset/python/README.md` PERF layout section (currently CWE/BP only) — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  
- [~] Extend `ruleset/python/catalogue_test.go` (or equivalent) so PERF chunks parse and respect ≤50/file — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  
- [~] Validation after seed: `go test ./ruleset/python/` green — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  

### 1.6 Phase 1 closure gate

- [~] Catalogue seed committed and documented; parent ledger #54 row can note “catalogue unblocked” — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  

**Proof when unblocked (record outcomes then):**

```sh
go test ./ruleset/python/
# list chunk files: ruleset/python/chunks/perf-*.json
# assert each file has ≤50 entries
```

---

## Phase 2: Detectors (deferred)

> **Blocked on Phase 1.** Reference architecture only: `internal/lang/go/detectors/perf/` (`scan.go`, `facts.go`, `register.go`, batch rules, fixtures).

### 2.1 Package / registration skeleton

- [~] Add Python PERF detector package under `internal/lang/python/` (e.g. `detectors/perf/` or family-consistent layout chosen with #52/#53) — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  
- [~] Wire `Register` / `init` so seeded `PERF-*` IDs can register independently — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  
- [~] Expose detectors via Python plugin `Detectors()` / `NewDetectors()` without enabling Python by default — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  
- [~] Ensure pack classification treats `PERF-*` as performance (existing language-agnostic pack prefix rules in `internal/rules/pack.go`) — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  

### 2.2 Heuristic implementation (seeded IDs only)

Driven by catalogue `detection_notes`; prefer pure-Go source-pattern / light-parse heuristics consistent with sibling CWE/BP approach.

- [~] Implement detectors for initial seed IDs only (async blocking, loop-expensive ops, N+1/ORM cost if PERF-owned, unbounded growth, HTTP timeout if PERF-owned) — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  
- [~] Suppress paths documented in `detection_notes` (hoisted compile, timeouts set, pagination present, etc.) — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  
- [~] Avoid double-firing with BP detectors for shared themes (shared suppress or single-owner rule) — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  

### 2.3 Fixtures

- [~] Add hit/miss fixtures per seeded `PERF-*` under the repo’s Python fixture layout (e.g. `tests/fixtures/python/…`) — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  
- [~] Table tests: vulnerable → finding; safe → no finding — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  

### 2.4 Integration smoke

- [~] With `languages = ["python"]` (or multi-lang including python), CLI can emit `PERF-*` findings on fixtures — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  
- [~] Default Go-only behavior unchanged (no accidental Python enable) — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  

### 2.5 Phase 2 validation (when active)

```sh
gofmt -w <changed Go files>
make lint
make test
CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop
# smoke with languages including python against PERF fixtures
```

- [~] Record `make lint` + `make test` outcomes on the implement branch — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  

---

## Phase 3: Docs, mapping, closure (deferred)

- [~] README / mapping notes updated for PERF layout (`ruleset/python/README.md`, any maturity / parity pointers) — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  
- [~] Parent ledger `plans/v0.0.2/python-heuristics.md` PERF rollup updated from deferred → partial/complete with evidence — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  
- [~] Filled PR body under `plans/PR/v0.0.2/` when implementation lands; `Closes #54` / `Relates to #51` — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  
- [~] Issue #54 success criteria checked against evidence (catalogue, fixtures, plugin emission, lint/test) — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  

---

## When unblocked — backlog only (still not ready-to-start)

> Keep these as a **reminder list**, not an active sprint. Promote into Phase 1–3 checkboxes (and change `[~]` → `[ ]`) only after the catalogue gate is met.

1. Seed `ruleset/python/chunks/perf-001-050.json` (and further ranges if needed) with initial Python PERF themes.  
2. Document ownership vs `BP-PY-*` for N+1 and HTTP timeouts.  
3. Extend `ruleset/python/catalogue_test.go` + README.  
4. Implement first detector batch + hit/miss fixtures for seeded IDs.  
5. Register on Python plugin; smoke CLI `PERF-*` findings.  
6. `make lint` + `make test`; close #54 with evidence.

---

## Unblock criteria

All of the following must be true before treating detector work as active `[ ]` rows:

| # | Criterion | Proof |
|---|-----------|-------|
| 1 | At least one `ruleset/python/chunks/perf-*.json` exists | `ls ruleset/python/chunks/perf-*.json` |
| 2 | Each file ≤50 rules; golang-style range names | scripted `len(json)` per file |
| 3 | Rules use `applicable_to` ⊇ `python`, `python_relevance`, Python `detection_notes` | sample + test |
| 4 | Catalogue tests pass | `go test ./ruleset/python/` |
| 5 | README documents PERF layout | `ruleset/python/README.md` |
| 6 | Explicit decision logged: “catalogue gate met; promote Phase 2 rows” | note in this ledger + parent `python-heuristics.md` |

Until then: **status remains deferred**; do not open “implement PERF detectors” PRs as primary scope.

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

## Success criteria (issue #54 — all deferred until unblocked)

- [~] Python PERF JSON chunks exist (≤50 rules/file, range names) — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  
- [~] Initial PERF heuristics registered and fixture-tested — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  
- [~] `languages = ["python"]` can emit `PERF-*` findings — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  
- [~] README / mapping notes updated for PERF layout — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  
- [~] `make lint` + `make test` green on implement branch — **reason:** no Python PERF catalogue yet; **owner:** epic #51 / issue #54; **next gate:** seed `ruleset/python/chunks/perf-*.json` (≤50/file, golang-style names) before detector work  

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
