# v0.0.2 / #54 — PERF-PY implementation checklist

> **Parent:** `plans/v0.0.2/heuristics/python-heuristics-perf.md` — catalogue and #54 ledger
> **Source evaluation:** `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/plans/perf-evaluation/`
> **Status:** planned — catalogue seeded (`PERF-PY-1` through `PERF-PY-22`); no Python PERF detector, fixture, or CLI emission exists yet
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

- [ ] **PERF-PY foundation** — add `internal/lang/python/detectors/perf/scan.go` with `PythonPerfScan`, stable sorted `RuleIDs`, `MetadataFor`, `Run`, context allow-list handling, and immutable post-`init()` catalogue snapshot. Expected: one Python PERF detector owns all implemented PERF-PY IDs. Proof: unit test verifies fresh scan instances expose the same catalogue.
- [ ] **PERF-PY registration** — add `register.go` with idempotent `RegisterRule`, metadata map updates, and optional safe source-needle gates. Expected: duplicate init registration replaces one entry, never emits two findings. Proof: registration collision unit test.
- [ ] **PERF-PY metadata** — add `metadata.go` for `PERF-PY-1`…`PERF-PY-22`, matching title, description, severity, and `PackPerformance`; set explicit experimental maturity until a later certification decision. Proof: metadata parity test against both JSON chunks.
- [ ] **PERF-PY facts** — add `facts.go` and `common.go` with a single source index, comment-stripped code lines, indentation-aware function/loop windows, and bounded forward/backward line searches. Expected: rules do not independently rescan or parse full source. Proof: focused facts tests for nested loops and comments/strings.
- [ ] **Python plugin wire-up** — add `perf.NewPythonPerfScan()` to `internal/lang/python/detectors/all.go`; preserve Go-only `DefaultRegistry`. Proof: `internal/lang/python/plugin_test.go` sees `PERF-PY-1` and `PERF-PY-22` only when Python is explicitly enabled.

### 1.2 Fixture and integration plumbing

- [ ] **Fixture directory** — add `tests/fixtures/python/perf/` using the existing `lang: python` / `file:` fixture format. Expected: every implemented PERF-PY rule receives one vulnerable and one safe fixture. Proof: pair-discovery test rejects missing or duplicate pairs.
- [ ] **Python PERF discovery** — extend `tests/integration/discover.go` with `DiscoverPythonPERFCases`, `PythonPERFRuleID`, `PythonPERFFixtureRel`, and numeric sort for `PERF-PY-N` stems. Expected: a Python ID is never normalized into Go `PERF-N`. Proof: discovery unit tests for `PERF-PY-1` and `PERF-PY-22` variants.
- [ ] **Python PERF matrix** — add `tests/integration/python/perf_matrix_test.go`. Expected: each vulnerable fixture emits its exact rule, and each safe fixture stays silent for that exact rule under `ProfileAll`. Proof: matrix runs with `languages=[python]` and `Only=[PERF-PY-N]`.
- [ ] **Raw CLI smoke seam** — add one integration test that materializes a Python PERF fixture with `languages=["python"]` and asserts `PERF-PY-*` JSON/text output. Expected: plugin registration is observable outside in-process unit tests. Proof: one vulnerable fixture scan.

### 1.3 Phase 1 gate

- [ ] **Foundation gate** — no rule-specific implementation starts until scaffold, metadata parity, and empty-matrix plumbing compile. Proof: `go test ./internal/lang/python/... ./tests/integration/python -run 'Python.*PERF|PythonPlugin'`.

---

## Phase 2: Local request, batch, and idempotency rules

These rules are intraprocedural and should ship as the first detector PR after foundation. Target file: `internal/lang/python/detectors/perf/rules_local.go`; split into `rules_batch.go` before either file exceeds 1,500 lines.

### 2.1 `PERF-PY-1` — Full result set materialized before app-side sort

- [ ] **Detector** — detect `.scalars().all()` / `.all()` assigned to a name followed in the same function by `sorted(name)`, `name.sort()`, or percentile-like index access. Suppress bounded export/admin paths and SQL aggregate/window queries. Path: `rules_local.go`. Proof: unit vulnerable/miss snippets.
- [ ] **Fixtures** — add `PERF-PY-1-{vulnerable,safe}.txt`; safe case uses `func.percentile_cont`/aggregate or a bounded explicit export. Proof: Python PERF matrix.

### 2.2 `PERF-PY-3` — Django per-row create in batch loop

- [ ] **Detector** — detect `Model.objects.create(...)` in an item loop, while preserving cases where the new object is immediately needed for a dependent operation. Do not recommend `bulk_create` when save hooks or generated IDs are required. Path: `rules_batch.go`. Proof: hit/miss unit tests.
- [ ] **Fixtures** — add paired fixtures for a loop of independent creates and a dependent-create safe case. Proof: matrix target rule only.

### 2.3 `PERF-PY-4` — Django read-modify-write counter update

- [ ] **Detector** — find numeric-looking model field `+=`/`-=` followed by `.save()` in the same function or loop. Suppress existing `F()`/`QuerySet.update` use and documented required model hooks. Path: `rules_batch.go`. Proof: hit/miss unit tests.
- [ ] **Fixtures** — add paired `stock/reserved/quantity` source cases; safe uses `models.F(...)` in `update`. Proof: matrix target rule only.

### 2.4 `PERF-PY-9` — JSON payload parse then re-serialize

- [ ] **Detector** — match request JSON parse into a value followed by `json.dumps(value)` or `str(value)` before ORM persistence in the same function. Suppress redaction, normalization, and intentional projection. Path: `rules_local.go`. Proof: hit/miss unit tests.
- [ ] **Fixtures** — add Flask/FastAPI parse-dump vulnerable case and raw-byte/explicit-transform safe cases. Proof: matrix target rule only.

### 2.5 `PERF-PY-10` — Worker sleeps after successful batch

- [ ] **Detector** — detect a worker loop that assigns a processed count, conditionally handles a truthy result, then unconditionally calls `time.sleep` or `await asyncio.sleep`. Suppress `continue`/`return` success paths and explicit rate limiting. Path: `rules_batch.go`. Proof: hit/miss unit tests.
- [ ] **Fixtures** — add paired polling-loop cases. Proof: matrix target rule only.

### 2.6 `PERF-PY-11` — Per-row ORM mutation instead of set-based update

- [ ] **Detector** — match full result hydration followed by uniform field assignments and one commit/save sequence. Suppress per-row derived values, external effects, and model-hook-required code. Path: `rules_batch.go`. Proof: hit/miss unit tests.
- [ ] **Fixtures** — add dead-letter redrive vulnerable case and `QuerySet.update` safe case. Proof: matrix target rule only.

### 2.7 `PERF-PY-12` — Unbounded JSON request body parse

- [ ] **Detector** — detect request-body JSON parsing in route/view code without a preceding content-length/body-length limit or configured parser cap in the same function. Suppress streamed parsers and visible framework limits. Path: `rules_local.go`. Proof: hit/miss unit tests.
- [ ] **Fixtures** — add paired Django/Flask request-body cases. Proof: matrix target rule only.

### 2.8 `PERF-PY-14` — Select-then-insert idempotency check

- [ ] **Detector** — match an idempotency/request/event key lookup followed by create/insert of the same model. Suppress dialect upserts, conflict handling, and `get_or_create` protected by a unique constraint. Path: `rules_local.go`. Proof: hit/miss unit tests.
- [ ] **Fixtures** — add paired Django and SQLAlchemy idempotency cases. Proof: matrix target rule only.

### 2.9 Phase 2 gate

- [ ] **Local-rule gate** — confirm all Phase 2 IDs are registered, metadata-backed, fixture-covered, and do not fire on their safe controls. Proof: targeted detector package test plus Python PERF fixture matrix.

---

## Phase 3: ORM, locking, indexes, and maintenance rules

These rules use wider source windows. Rules `15`, `16`, and `20` must be review-level: only report a missing index when the parsed declaration scope is sufficient; otherwise remain silent.

### 3.1 `PERF-PY-2` — Django ORM lookup inside item loop

- [ ] **Detector** — match terminal Django ORM lookups (`get`, `first`, evaluating `all`) in an item loop, including repeated loop-invariant lookup. Suppress a tiny explicit collection, paged management command, and lazy query construction alone. Path: `rules_orm.go`. Proof: hit/miss unit tests.
- [ ] **Fixtures** — add N+1 lookup vulnerable case and hoisted/prefetched safe case. Proof: matrix target rule only; coordinate suppression with BP-PY-28 if both would report the same source span.

### 3.2 `PERF-PY-6` — ORM work claim without row lock

- [ ] **Detector** — match pending/status selection, later status assignment, and commit/save without `select_for_update`, `with_for_update`, `skip_locked`, or atomic update. Suppress documented single-process worker mode. Path: `rules_orm.go`. Proof: hit/miss unit tests.
- [ ] **Fixtures** — add paired queue-claim cases with and without locking. Proof: matrix target rule only.

### 3.3 `PERF-PY-8` — SQLAlchemy lazy relationship access in batch loop

- [ ] **Detector** — match a loop over query results with relation-like attribute access and no visible `joinedload`, `selectinload`, or `contains_eager` in the query/function window. Keep it review-level and avoid generic scalar attributes. Path: `rules_orm.go`. Proof: hit/miss unit tests.
- [ ] **Fixtures** — add paired eager-load and lazy-loop fixtures. Proof: matrix target rule only.

### 3.4 `PERF-PY-13` — Full ORM hydration for projection read

- [ ] **Detector** — match an ORM iteration using only one/two scalar fields without methods/relations; suppress `values`, `values_list`, `only`, `defer`, `load_only`, and explicit column selects. Path: `rules_orm.go`. Proof: hit/miss unit tests.
- [ ] **Fixtures** — add paired read-projection cases. Proof: matrix target rule only.

### 3.5 `PERF-PY-15` — ORM composite filter without supporting index

- [ ] **Detector** — correlate visible model declarations with a tenant/owner plus time-range ORM filter. Emit only when no local `Index`, `db_index`, `Meta.indexes`, or `__table_args__` covers the query shape. Path: `rules_orm_index.go`. Proof: cross-file unit hit/miss snippets.
- [ ] **Fixtures** — add paired model-plus-query fixture cases, including an external-migration/comment-safe control that must remain silent. Proof: matrix target rule only.

### 3.6 `PERF-PY-16` — Retention timestamp predicate without index

- [ ] **Detector** — match cleanup/purge/expiry predicate on timestamp against cutoff with no visible matching index/partition declaration. Do not conflate this with transaction batching (`PERF-PY-21`). Path: `rules_orm_index.go`. Proof: hit/miss unit tests.
- [ ] **Fixtures** — add paired TTL/retention fixtures with indexed and unindexed model declarations. Proof: matrix target rule only.

### 3.7 `PERF-PY-19` — Unbounded ORM locking sweep

- [ ] **Detector** — detect `transaction.atomic`/transaction scope with locking query iteration and no slice, `limit`, keyset loop, batch helper, or per-batch commit. Path: `rules_orm.go`. Proof: hit/miss unit tests.
- [ ] **Fixtures** — add paired expiry-sweep fixtures. Proof: matrix target rule only.

### 3.8 `PERF-PY-20` — ORM sort without supporting composite index

- [ ] **Detector** — correlate equality filter plus `order_by`/`order_by('-quantity')` with visible model index declarations. Suppress primary-key singleton and external migration/config cases. Path: `rules_orm_index.go`. Proof: hit/miss unit tests.
- [ ] **Fixtures** — add paired filter/order fixtures. Proof: matrix target rule only.

### 3.9 `PERF-PY-21` — Unbounded bulk delete in maintenance path

- [ ] **Detector** — detect retention/status `QuerySet.delete`/SQLAlchemy `Query.delete` without visible limit, key-range loop, chunk helper, or per-batch commit. Suppress one-row deletes and partition drops. Path: `rules_orm.go`. Proof: hit/miss unit tests.
- [ ] **Fixtures** — add paired purge fixtures. Proof: matrix target rule only.

### 3.10 Phase 3 gate

- [ ] **ORM-rule gate** — run all Phase 3 fixtures under `--only PERF-PY-N`; review all vulnerable message text to retain “review-level” wording for inference-dependent rules. Proof: detector unit suite + Python PERF matrix.

---

## Phase 4: Framework, runtime, and configuration rules

### 4.1 `PERF-PY-5` — Sequential blocking delivery over claimed batch

- [ ] **Detector** — match worker-like functions that claim a batch and synchronously call delivery/network work in a for-loop. Suppress bounded executor, semaphore, `gather` with a limit, and explicit per-endpoint limiter paths. Path: `rules_runtime.go`. Proof: hit/miss unit tests.
- [ ] **Fixtures** — add paired sequential and bounded-fan-out worker fixtures. Proof: matrix target rule only.

### 4.2 `PERF-PY-7` — BaseHTTPMiddleware on FastAPI request path

- [ ] **Detector** — require a project-defined `BaseHTTPMiddleware` subclass plus visible app registration; do not flag third-party middleware imports alone. Path: `rules_runtime.go`. Proof: hit/miss unit tests.
- [ ] **Fixtures** — add paired BaseHTTPMiddleware and pure-ASGI middleware fixtures. Proof: matrix target rule only.

### 4.3 `PERF-PY-17` — Database connection reuse or timeout controls missing

- [ ] **Detector** — inspect SQLAlchemy engine kwargs, Django `DATABASES`, and Flask engine options for missing pool lifecycle/health/timeout controls. Suppress tests, local SQLite, and externally imported settings. Path: `rules_config.go`. Proof: framework-specific hit/miss unit tests.
- [ ] **Fixtures** — add FastAPI, Django, and Flask configuration pairs. Proof: matrix target rule only.

### 4.4 `PERF-PY-18` — Repeated regex rewrites on the same input

- [ ] **Detector** — detect sequential `re.sub`/compiled `.sub` passes that reassign the same hot-path variable; suppress deliberately staged/overlapping transformations. Path: `rules_runtime.go`. Proof: hit/miss unit tests.
- [ ] **Fixtures** — add paired one-pass and repeated-pass normalization fixtures. Proof: matrix target rule only.

### 4.5 `PERF-PY-22` — SQLite backend for concurrent service writes

- [ ] **Detector** — match Django SQLite/Flask SQLite default configuration only when the same scope indicates workers, queues, threads, or process concurrency. Suppress tests, fixtures, and explicitly local-only config. Path: `rules_config.go`. Proof: hit/miss unit tests.
- [ ] **Fixtures** — add paired concurrent-service and test/local SQLite configurations. Proof: matrix target rule only.

### 4.6 Phase 4 gate

- [ ] **Runtime-rule gate** — ensure configuration rules do not fire on test fixtures or plain SQLite examples lacking concurrent service evidence. Proof: targeted unit tests plus the complete Python PERF matrix.

---

## Phase 5: Integration, canary, and closure gates

### 5.1 Catalogue, detector, and fixture parity

- [ ] **22/22 registration parity** — assert every JSON `PERF-PY-*` key has exactly one registered detector rule and non-nil metadata; assert no bare Go `PERF-N` is registered by the Python detector. Path: `internal/lang/python/detectors/perf/scan_test.go`. Proof: unit test enumerates 22 IDs.
- [ ] **22/22 fixture parity** — enforce one vulnerable and one safe fixture per registered Python PERF rule. Path: `tests/integration/python/perf_matrix_test.go`. Proof: inventory test expects 22 cases.
- [ ] **No BP duplicate policy** — add regression cases for the shared ORM/N+1, async, and timeout neighborhoods. Expected: a `PERF-PY-*` result is emitted only for the cost-specific condition; BP-PY continues to own blocking async and HTTP-timeout style rules. Proof: targeted `Only` scans and documented suppressions.

### 5.2 Corpus canary and false-positive review

- [ ] **Reference-corpus canary** — scan the three source evaluation projects with Python enabled and `ProfileAll`; record `PERF-PY-*` counts and finding paths in this checklist. Expected: planned source examples fire; unexpected broad project-wide noise blocks promotion. Proof: exact command, revision, and output summary.
- [ ] **False-positive triage** — classify every canary finding as expected, fixed, accepted review-level, or false positive. Tighten a rule/fixture before changing maturity. Proof: table appended to this checklist with path, rule, disposition, and reason.
- [ ] **Maturity decision** — keep all rules experimental unless per-rule canary evidence demonstrates reliable production signal. Proof: metadata/maturity test and recorded rationale; do not add `PERF-PY-*` IDs to recommended/perf tier lists without this evidence.

### 5.3 Required implementation validation

- [ ] **Formatting** — run `gofmt -w` on every touched Go file. Proof: command and clean diff outcome recorded before merge.
- [ ] **Lint** — run `make lint`. Proof: exact command and passing outcome recorded here.
- [ ] **Test suite** — run `make test`. Proof: exact command and passing outcome recorded here.
- [ ] **Build** — run `CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop`. Proof: exact command and passing outcome recorded here.
- [ ] **CLI smoke** — run the built binary against one materialized Python PERF fixture using a config with `languages = ["python"]`. Proof: output contains only the selected `PERF-PY-N` finding.

### 5.4 Ledger closure

- [ ] **Parent ledger** — update `plans/v0.0.2/heuristics/python-heuristics-perf.md` only after each batch’s detector, fixture, and quality proof is current; record completed IDs and leave unimplemented IDs unchecked.
- [ ] **Ruleset docs** — update `ruleset/python/README.md` from “catalogue seeded” to “detectors implemented” only when all 22 IDs emit and the full validation gate passes.
- [ ] **PR handoff** — prepare one PR per phase/batch with `Relates to #54` and `Relates to #51`; do not claim `Closes #54` until every Phase 5 gate is checked.

---

## Rule tracker

| ID | Name | Phase | Status |
|---|---|---:|---|
| PERF-PY-1 | Full Result Set Materialized Before App-Side Sort | 2 | [ ] |
| PERF-PY-2 | Django ORM Lookup Inside Item Loop | 3 | [ ] |
| PERF-PY-3 | Django Per-Row Create In Batch Loop | 2 | [ ] |
| PERF-PY-4 | Django Read-Modify-Write Counter Update | 2 | [ ] |
| PERF-PY-5 | Sequential Blocking Delivery Over Claimed Batch | 4 | [ ] |
| PERF-PY-6 | ORM Work Claim Without Row Lock | 3 | [ ] |
| PERF-PY-7 | BaseHTTPMiddleware On FastAPI Request Path | 4 | [ ] |
| PERF-PY-8 | SQLAlchemy Lazy Relationship Access In Batch Loop | 3 | [ ] |
| PERF-PY-9 | JSON Payload Parse Then Re-Serialize | 2 | [ ] |
| PERF-PY-10 | Worker Sleeps After Successful Batch | 2 | [ ] |
| PERF-PY-11 | Per-Row ORM Mutation Instead Of Set-Based Update | 2 | [ ] |
| PERF-PY-12 | Unbounded JSON Request Body Parse | 2 | [ ] |
| PERF-PY-13 | Full ORM Hydration For Projection Read | 3 | [ ] |
| PERF-PY-14 | Select-Then-Insert Idempotency Check | 2 | [ ] |
| PERF-PY-15 | ORM Composite Filter Without Supporting Index | 3 | [ ] |
| PERF-PY-16 | Retention Timestamp Predicate Without Index | 3 | [ ] |
| PERF-PY-17 | Database Connection Reuse Or Timeout Controls Missing | 4 | [ ] |
| PERF-PY-18 | Repeated Regex Rewrites On The Same Input | 4 | [ ] |
| PERF-PY-19 | Unbounded ORM Locking Sweep | 3 | [ ] |
| PERF-PY-20 | ORM Sort Without Supporting Composite Index | 3 | [ ] |
| PERF-PY-21 | Unbounded Bulk Delete In Maintenance Path | 3 | [ ] |
| PERF-PY-22 | SQLite Backend For Concurrent Service Writes | 4 | [ ] |

**Coverage check:** all 22 seeded catalogue IDs have one owner phase. No detector work is marked implemented by this plan.

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
