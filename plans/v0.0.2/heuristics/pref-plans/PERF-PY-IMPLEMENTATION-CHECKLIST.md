# v0.0.2 / #54 — PERF-PY implementation checklist

> **Parent:** `plans/v0.0.2/heuristics/python-heuristics-perf.md` — catalogue and #54 ledger
> **Source evaluation:** `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/plans/perf-evaluation/`
> **Status:** implementation complete — 22 experimental `PERF-PY-*` detectors are registered in the opt-in Python plugin, with paired fixtures and passing repository validation; corpus canary and maturity promotion remain pending
> **Estimated effort:** 5 focused PRs; medium–large (source-only heuristics, 22 hit/miss pairs, matrix and CLI proof)
> **Branch sequence:** `feat/python-perf-scaffold` → `feat/python-perf-local` → `feat/python-perf-orm` → `feat/python-perf-runtime` → `feat/python-perf-integration`
> **Ledger rule:** mark a row `[x]` only after the matching detector/fixture proof is green. For every implementation PR, record `make lint` and `make test` here; catalogue validation alone never proves rule emission.

---

## Overview

This is the single execution ledger for the 22 seeded `PERF-PY-*` catalogue entries in:

- `ruleset/python/chunks/perf-py-001-014.json`
- `ruleset/python/chunks/perf-py-015-022.json`

The implementation follows the shipped Python BP architecture rather than the Go AST detector directly:

| Concern | Required shape |
|---|---|
| Package | New `internal/lang/python/detectors/perf/` package |
| Detector | One `PythonPerfScan`, containing a registered, immutable rule catalogue |
| Parsing | Pure-Go source patterns over `ParsedUnit.Source`; no Python AST, CGO, or type checker |
| Facts | One per-file `pyPerfFacts` bag: source index, comment-stripped lines, function/loop spans, and small assignment/query windows |
| Registration | `RegisterRule("PERF-PY-N", fn, meta, gates...)` from `init()` in domain files |
| Metadata | Hand-written metadata matching the two Python PERF JSON chunks; all rules start `experimental` until fixture and canary evidence supports a maturity change |
| Plugin | Add exactly one PERF detector to `internal/lang/python/detectors/all.go`; Python remains opt-in |
| Fixtures | Paired `tests/fixtures/python/perf/PERF-PY-N-{vulnerable,safe}.txt` files |
| Integration | Add Python PERF discovery/matrix support; each fixture scan uses `languages=[python]`, `--only PERF-PY-N`, and `ProfileAll` |

### Source-evaluation coverage

| Source finding(s) | Rule(s) | Implementation phase |
|---|---|---|
| FA-1 | PERF-PY-1 | 2 |
| FA-2, FA-3, FA-5, FA-6 | PERF-PY-15, 16, 17, 18 | 3 / 4 |
| FA-4, FA-7 | PERF-PY-7, 14 | 4 / 2 |
| DJ-1, DJ-2, DJ-3 | PERF-PY-2, 3, 4 | 2 / 3 |
| DJ-4, DJ-6, DJ-7, DJ-8 | PERF-PY-19, 12, 13, 22, 20 | 2 / 3 / 4 |
| FL-1 through FL-8 (except no separate rule is needed for FL-5's index half) | PERF-PY-5, 6, 8, 9, 10, 11, 16, 21, 22 | 2 / 3 / 4 |
| XC-2 | PERF-PY-17 | 4 |

FA-8 is correctness/operability; XC-1/XC-3/XC-4 are benchmark, deployment, and observability work; XC-5 is security. They remain outside this source PERF detector checklist.

## Executive Summary

Implement in dependency order:

```text
Phase 1: PERF package + metadata + fixture discovery
    ├─ Phase 2: local/request/batch data-flow rules (high signal)
    ├─ Phase 3: ORM/query/index/maintenance rules (review-level where cross-file)
    └─ Phase 4: framework/runtime/config rules
             └─ Phase 5: fixture matrix, raw CLI smoke, canary and closure gates
```

Rules that need absent-index or deployment inference (`15`, `16`, `17`, `20`, `22`) are **review-level experimental**. They must suppress when the relevant migration/settings module is unavailable; they must never claim the database lacks an external index.

---

## Phase 1: Detector foundation and test contract

### 1.1 Package and immutable catalogue

- [x] **PERF-PY foundation** — implemented in `internal/lang/python/detectors/perf/scan.go`; `scan_test.go` verifies the 22-ID immutable catalogue (2026-07-31).
- [x] **PERF-PY registration** — implemented in `register.go`; each rule is registered once at package init and exercised by the package suite (2026-07-31).
- [x] **PERF-PY metadata** — implemented in `metadata.go`; `metadata_test.go` verifies all 22 IDs/titles against both JSON chunks (2026-07-31).
- [x] **PERF-PY facts** — implemented in `facts.go` and `common.go`; all rules consume the shared per-file source/line facts (2026-07-31).
- [x] **Python plugin wire-up** — `perf.NewPythonPerfScan()` is in `internal/lang/python/detectors/all.go`; plugin test asserts `PERF-PY-1` and `PERF-PY-22` under explicit Python enablement (2026-07-31).

### 1.2 Fixture and integration plumbing

- [x] **Fixture directory** — added `tests/fixtures/python/perf/` with 22 vulnerable/safe pairs in the established fixture format (2026-07-31).
- [x] **Python PERF discovery** — added `DiscoverPythonPERFCases`, `PythonPERFRuleID`, and `PythonPERFFixtureRel`; the 22-case inventory proves Python IDs do not enter Go PERF normalization (2026-07-31).
- [x] **Python PERF matrix** — added `tests/integration/python/perf_matrix_test.go`; all 22 vulnerable fixtures emit their exact selected ID and all safe controls are silent (2026-07-31).
- [~] **Raw CLI smoke seam** — integration materialization proves the Python-enabled analyzer boundary; a standalone binary smoke command remains a Phase 5 closure gate.

### 1.3 Phase 1 gate

- [x] **Foundation gate** — `go test ./internal/lang/python/detectors/perf ./tests/integration/python -count=1` passed (2026-07-31).

---

## Phase 2: Local request, batch, and idempotency rules

These rules are intraprocedural and should ship as the first detector PR after foundation. Target file: `internal/lang/python/detectors/perf/rules_local.go`; split into `rules_batch.go` before either file exceeds 1,500 lines.

### 2.1 `PERF-PY-1` — Full result set materialized before app-side sort

- [x] **Detector** — implemented in `rules_local.go`; focused hit/miss test and matrix pass (2026-07-31).
- [x] **Fixtures** — `PERF-PY-1-{vulnerable,safe}.txt` pass the Python PERF matrix (2026-07-31).

### 2.2 `PERF-PY-3` — Django per-row create in batch loop

- [x] **Detector** — implemented in `rules_batch.go` with dependent-create suppression; focused hit/miss test passes (2026-07-31).
- [x] **Fixtures** — paired independent/dependent create fixtures pass the matrix (2026-07-31).

### 2.3 `PERF-PY-4` — Django read-modify-write counter update

- [x] **Detector** — implemented in `rules_batch.go` with `F()`/update and hook suppression; focused hit/miss test passes (2026-07-31).
- [x] **Fixtures** — paired counter fixtures pass the matrix (2026-07-31).

### 2.4 `PERF-PY-9` — JSON payload parse then re-serialize

- [x] **Detector** — implemented in `rules_local.go` with raw/intentional-transform controls; focused hit/miss test passes (2026-07-31).
- [x] **Fixtures** — paired parse/dump fixtures pass the matrix (2026-07-31).

### 2.5 `PERF-PY-10` — Worker sleeps after successful batch

- [x] **Detector** — implemented in `rules_batch.go` with success-path `continue` suppression; focused hit/miss test passes (2026-07-31).
- [x] **Fixtures** — paired polling fixtures pass the matrix (2026-07-31).

### 2.6 `PERF-PY-11` — Per-row ORM mutation instead of set-based update

- [x] **Detector** — implemented in `rules_batch.go` with set-based update control; focused hit/miss test passes (2026-07-31).
- [x] **Fixtures** — paired redrive fixtures pass the matrix (2026-07-31).

### 2.7 `PERF-PY-12` — Unbounded JSON request body parse

- [x] **Detector** — implemented in `rules_local.go` with visible body-limit suppression; focused hit/miss test passes (2026-07-31).
- [x] **Fixtures** — paired request-body fixtures pass the matrix (2026-07-31).

### 2.8 `PERF-PY-14` — Select-then-insert idempotency check

- [x] **Detector** — implemented in `rules_local.go` with `get_or_create`/conflict-safe suppression; focused hit/miss test passes (2026-07-31).
- [x] **Fixtures** — paired idempotency fixtures pass the matrix (2026-07-31).

### 2.9 Phase 2 gate

- [x] **Local-rule gate** — targeted detector package and 22-case Python PERF matrix pass (2026-07-31).

---

## Phase 3: ORM, locking, indexes, and maintenance rules

These rules use wider source windows. Rules `15`, `16`, and `20` must be review-level: only report a missing index when the parsed declaration scope is sufficient; otherwise remain silent.

### 3.1 `PERF-PY-2` — Django ORM lookup inside item loop

- [x] **Detector** — implemented in `rules_orm.go`; loop lookup hit/miss test passes (2026-07-31).
- [x] **Fixtures** — paired N+1/hoisted fixtures pass the matrix (2026-07-31).

### 3.2 `PERF-PY-6` — ORM work claim without row lock

- [x] **Detector** — implemented in `rules_orm.go` with row-lock/atomic/single-process suppression; focused hit/miss test passes (2026-07-31).
- [x] **Fixtures** — paired queue-claim fixtures pass the matrix (2026-07-31).

### 3.3 `PERF-PY-8` — SQLAlchemy lazy relationship access in batch loop

- [x] **Detector** — implemented as a narrow relation-shaped, review-level rule in `rules_orm.go`; eager-load control passes (2026-07-31).
- [x] **Fixtures** — paired lazy/eager fixtures pass the matrix (2026-07-31).

### 3.4 `PERF-PY-13` — Full ORM hydration for projection read

- [x] **Detector** — implemented in `rules_orm.go` with projection controls; focused hit/miss test passes (2026-07-31).
- [x] **Fixtures** — paired projection fixtures pass the matrix (2026-07-31).

### 3.5 `PERF-PY-15` — ORM composite filter without supporting index

- [x] **Detector** — implemented in `rules_orm_index.go`; requires a visible local model declaration and suppresses index/external-migration evidence (2026-07-31).
- [x] **Fixtures** — paired indexed/unindexed fixtures pass the matrix (2026-07-31).

### 3.6 `PERF-PY-16` — Retention timestamp predicate without index

- [x] **Detector** — implemented in `rules_orm_index.go` with visible-index/partition guards; focused hit/miss test passes (2026-07-31).
- [x] **Fixtures** — paired retention fixtures pass the matrix (2026-07-31).

### 3.7 `PERF-PY-19` — Unbounded ORM locking sweep

- [x] **Detector** — implemented in `rules_orm.go` with slice/limit/batch guards; focused hit/miss test passes (2026-07-31).
- [x] **Fixtures** — paired locking-sweep fixtures pass the matrix (2026-07-31).

### 3.8 `PERF-PY-20` — ORM sort without supporting composite index

- [x] **Detector** — implemented in `rules_orm_index.go` with visible-index/primary-key/external-migration guards (2026-07-31).
- [x] **Fixtures** — paired filter/order fixtures pass the matrix (2026-07-31).

### 3.9 `PERF-PY-21` — Unbounded bulk delete in maintenance path

- [x] **Detector** — implemented in `rules_orm.go` with batch/partition guards; focused hit/miss test passes (2026-07-31).
- [x] **Fixtures** — paired purge fixtures pass the matrix (2026-07-31).

### 3.10 Phase 3 gate

- [x] **ORM-rule gate** — `TestORMRules`, conservative-index suppression tests, and the Python PERF matrix pass; inference-dependent findings retain review-level wording (2026-07-31).

---

## Phase 4: Framework, runtime, and configuration rules

### 4.1 `PERF-PY-5` — Sequential blocking delivery over claimed batch

- [x] **Detector** — implemented in `rules_runtime.go` with bounded-fan-out suppression; focused hit/miss test passes (2026-07-31).
- [x] **Fixtures** — paired sequential/bounded-fan-out fixtures pass the matrix (2026-07-31).

### 4.2 `PERF-PY-7` — BaseHTTPMiddleware on FastAPI request path

- [x] **Detector** — implemented in `rules_runtime.go` with class-and-registration requirement; focused hit/miss test passes (2026-07-31).
- [x] **Fixtures** — paired BaseHTTPMiddleware/pure-ASGI fixtures pass the matrix (2026-07-31).

### 4.3 `PERF-PY-17` — Database connection reuse or timeout controls missing

- [x] **Detector** — implemented in `rules_runtime.go` with production-evidence and test/local/dev guards; focused hit/miss test passes (2026-07-31).
- [x] **Fixtures** — paired production/control configuration fixtures pass the matrix (2026-07-31).

### 4.4 `PERF-PY-18` — Repeated regex rewrites on the same input

- [x] **Detector** — implemented in `rules_runtime.go` with staged-transformation suppression; focused hit/miss test passes (2026-07-31).
- [x] **Fixtures** — paired repeated/one-pass fixtures pass the matrix (2026-07-31).

### 4.5 `PERF-PY-22` — SQLite backend for concurrent service writes

- [x] **Detector** — implemented in `rules_runtime.go` with explicit production + concurrency evidence and test/local/dev guards; focused hit/miss test passes (2026-07-31).
- [x] **Fixtures** — paired concurrent-service/control SQLite fixtures pass the matrix (2026-07-31).

### 4.6 Phase 4 gate

- [x] **Runtime-rule gate** — targeted runtime tests and the complete Python PERF matrix pass, including test/local configuration controls (2026-07-31).

---

## Phase 5: Integration, canary, and closure gates

### 5.1 Catalogue, detector, and fixture parity

- [x] **22/22 registration parity** — `scan_test.go` and `metadata_test.go` enumerate every JSON ID, require performance metadata, and reject bare Go IDs (2026-07-31).
- [x] **22/22 fixture parity** — `perf_matrix_test.go` requires exactly 22 numerically ordered vulnerable/safe pairs (2026-07-31).
- [ ] **No BP duplicate policy** — add regression cases for the shared ORM/N+1, async, and timeout neighborhoods. Expected: a `PERF-PY-*` result is emitted only for the cost-specific condition; BP-PY continues to own blocking async and HTTP-timeout style rules. Proof: targeted `Only` scans and documented suppressions.

### 5.2 Corpus canary and false-positive review

- [ ] **Reference-corpus canary** — scan the three source evaluation projects with Python enabled and `ProfileAll`; record `PERF-PY-*` counts and finding paths in this checklist. Expected: planned source examples fire; unexpected broad project-wide noise blocks promotion. Proof: exact command, revision, and output summary.
- [ ] **False-positive triage** — classify every canary finding as expected, fixed, accepted review-level, or false positive. Tighten a rule/fixture before changing maturity. Proof: table appended to this checklist with path, rule, disposition, and reason.
- [ ] **Maturity decision** — keep all rules experimental unless per-rule canary evidence demonstrates reliable production signal. Proof: metadata/maturity test and recorded rationale; do not add `PERF-PY-*` IDs to recommended/perf tier lists without this evidence.

### 5.3 Required implementation validation

- [x] **Formatting** — `gofmt -w` ran on every touched Go file; `git diff --check` passed (2026-07-31).
- [x] **Lint** — `make lint` and the repository’s stricter `make lint-all` passed (2026-07-31).
- [x] **Test suite** — `make test` passed (2026-07-31).
- [x] **Build** — `CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop` passed (2026-07-31).
- [ ] **CLI smoke** — run the built binary against one materialized Python PERF fixture using a config with `languages = ["python"]`. Proof: output contains only the selected `PERF-PY-N` finding.

### 5.4 Ledger closure

- [x] **Parent ledger** — updated parent #54 rollup and this canonical tracker after detector, fixture, lint, test, and build proof (2026-07-31).
- [x] **Ruleset docs** — `ruleset/python/README.md` now describes the 22 experimental, fixture-backed PERF-PY detectors (2026-07-31).
- [ ] **PR handoff** — prepare one PR per phase/batch with `Relates to #54` and `Relates to #51`; do not claim `Closes #54` until every Phase 5 gate is checked.

---

## Rule tracker

| ID | Name | Phase | Status |
|---|---|---:|---|
| PERF-PY-1 | Full Result Set Materialized Before App-Side Sort | 2 | [x] |
| PERF-PY-2 | Django ORM Lookup Inside Item Loop | 3 | [x] |
| PERF-PY-3 | Django Per-Row Create In Batch Loop | 2 | [x] |
| PERF-PY-4 | Django Read-Modify-Write Counter Update | 2 | [x] |
| PERF-PY-5 | Sequential Blocking Delivery Over Claimed Batch | 4 | [x] |
| PERF-PY-6 | ORM Work Claim Without Row Lock | 3 | [x] |
| PERF-PY-7 | BaseHTTPMiddleware On FastAPI Request Path | 4 | [x] |
| PERF-PY-8 | SQLAlchemy Lazy Relationship Access In Batch Loop | 3 | [x] |
| PERF-PY-9 | JSON Payload Parse Then Re-Serialize | 2 | [x] |
| PERF-PY-10 | Worker Sleeps After Successful Batch | 2 | [x] |
| PERF-PY-11 | Per-Row ORM Mutation Instead Of Set-Based Update | 2 | [x] |
| PERF-PY-12 | Unbounded JSON Request Body Parse | 2 | [x] |
| PERF-PY-13 | Full ORM Hydration For Projection Read | 3 | [x] |
| PERF-PY-14 | Select-Then-Insert Idempotency Check | 2 | [x] |
| PERF-PY-15 | ORM Composite Filter Without Supporting Index | 3 | [x] |
| PERF-PY-16 | Retention Timestamp Predicate Without Index | 3 | [x] |
| PERF-PY-17 | Database Connection Reuse Or Timeout Controls Missing | 4 | [x] |
| PERF-PY-18 | Repeated Regex Rewrites On The Same Input | 4 | [x] |
| PERF-PY-19 | Unbounded ORM Locking Sweep | 3 | [x] |
| PERF-PY-20 | ORM Sort Without Supporting Composite Index | 3 | [x] |
| PERF-PY-21 | Unbounded Bulk Delete In Maintenance Path | 3 | [x] |
| PERF-PY-22 | SQLite Backend For Concurrent Service Writes | 4 | [x] |

**Coverage check:** all 22 seeded catalogue IDs have one implemented owner phase and one passing fixture pair; corpus canary/maturity closure remains pending.

## Dependencies

| Dependency | Why it matters |
|---|---|
| `ruleset/python/chunks/perf-py-*.json` | Source of truth for ID/title/description/relevance |
| `internal/lang/python/detectors/bad_practices/` | Registration, facts, line-walk, metadata, and test architecture reference |
| `internal/lang/python/detectors/cwe/` | PERF scanner can reuse its immutable-catalogue/needle-gate lifecycle shape |
| `internal/lang/python/detectors/all.go` | Required Python plugin aggregation point |
| `tests/integration/discover.go` | Needs Python PERF-specific fixture discovery; Go PERF normalizers must not consume `PERF-PY-*` |
| `tests/integration/python/` | Python-only fixture matrix location; keeps Go PERF parity suite unchanged |
| Source evaluation corpus | Canary and acceptance examples; it is not proof of detector correctness without paired safe fixtures |

## Deferred boundaries

- [~] FA-8 vendor export path — owner: application correctness/operability; next gate: separate correctness rule or app fix, not PERF-PY.
- [~] XC-1 benchmark harness, XC-3 server docs, XC-4 observability, XC-5 security — owner: application performance/operations/security delivery; next gate: their source evaluation ledger, not static Python PERF detection.
- [~] Recommended/perf profile promotion — owner: maturity policy; next gate: Phase 5 canary evidence and explicit tier decision.

## References

- Parent #54 ledger: `plans/v0.0.2/heuristics/python-heuristics-perf.md`
- BP execution model: `plans/v0.0.2/heuristics/bp-plans/README.md`
- Checklist procedure: `plans/skills/phase-wise-checklist/SKILLS.md`
- Catalogue: `ruleset/python/chunks/perf-py-001-014.json`, `ruleset/python/chunks/perf-py-015-022.json`
