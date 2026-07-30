# Batch 00 — Shipped BP-PY ledger (do not re-open)

> **Parent:** `plans/v0.0.2/heuristics/bp-plans/README.md` — batch index  
> **Canonical #53 ledger:** `plans/v0.0.2/heuristics/python-heuristics-bp.md`  
> **Issue:** [#53](https://github.com/chinmay-sawant/goslop/issues/53) under epic [#51](https://github.com/chinmay-sawant/goslop/issues/51)  
> **Status:** **closed inventory** — implemented `RegisterRule` subset; remaining work is **other batches only**  
> **Estimated effort:** none (ledger / evidence only)

---

## Overview

This file is a **short evidence ledger** of `BP-PY-*` rules already registered and unit-tested. It is **not** an implementation backlog.

- Do **not** re-open these IDs as `[ ]` pending work in sibling batches.
- Do **not** re-implement detectors for these IDs unless fixing a proven bug (separate PR, not batch-06/07/08 scope).
- Remaining catalogue coverage lives in batches 01–08+ (async, testing/deps/obs, prod, framework leftovers, etc.).

### Evidence baseline (package)

| Surface | Path |
|---------|------|
| Register API | `internal/lang/python/detectors/bad_practices/register.go` |
| Scan driver | `internal/lang/python/detectors/bad_practices/scan.go` |
| Metadata | `internal/lang/python/detectors/bad_practices/metadata.go` |
| Facts / needles | `internal/lang/python/detectors/bad_practices/facts.go` |
| Helpers | `internal/lang/python/detectors/bad_practices/common.go` |
| Unit proofs | `internal/lang/python/detectors/bad_practices/scan_test.go` (`TestBPRulesRegistered`, per-rule hit/miss) |
| Inventory JSON | `plans/v0.0.2/heuristics/bp-plans/_inventory.json` → `implemented` |

---

## Executive Summary

| ID | Name | File (RegisterRule + detect) | Status |
|----|------|------------------------------|--------|
| BP-PY-1 | Bare Except Clause | `rules_core.go` | [x] |
| BP-PY-2 | Except Pass | `rules_core.go` | [x] |
| BP-PY-4 | Mutable Default Argument | `rules_core.go` | [x] |
| BP-PY-6 | assert Used For Runtime Validation | `rules_core.go` | [x] |
| BP-PY-7 | open Without Context Manager | `rules_core.go` | [x] |
| BP-PY-8 | subprocess With shell=True | `rules_security.go` | [x] |
| BP-PY-9 | os.system Or os.popen | `rules_security.go` | [x] |
| BP-PY-10 | pickle Loads Untrusted Data | `rules_security.go` | [x] |
| BP-PY-11 | yaml.load Without SafeLoader | `rules_security.go` | [x] |
| BP-PY-12 | eval Or exec On Dynamic Input | `rules_security.go` | [x] |
| BP-PY-13 | Hardcoded Secret In Source | `rules_security.go` | [x] |
| BP-PY-16 | Flask DEBUG True In Production Code | `rules_framework.go` | [x] |
| BP-PY-17 | Flask SECRET_KEY Hardcoded | `rules_framework.go` | [x] |
| BP-PY-21 | Django DEBUG True In Settings | `rules_framework.go` | [x] |

**Count:** 14 rules shipped in this ledger.

---

## Phase 1: Core — `rules_core.go`

### 1.1 Registered rules (evidence)

- [x] `BP-PY-1` Bare Except — `RegisterRule` + `detectBPPY1` in `internal/lang/python/detectors/bad_practices/rules_core.go`; hit/miss in `scan_test.go` (`TestBPPY1BareExcept`, `TestBPPY1BroadException`)
- [x] `BP-PY-2` Except Pass — `rules_core.go` + `TestBPPY2ExceptPass`
- [x] `BP-PY-4` Mutable Default Argument — `rules_core.go` + `TestBPPY4…` (severity high from metadata)
- [x] `BP-PY-6` assert validation — `rules_core.go`; skips test paths via `isPythonTestFile` in `common.go`
- [x] `BP-PY-7` open without `with` — `rules_core.go` + hit/miss open tests

### 1.2 Explicitly **not** in this shipped set

- [x] Document: `BP-PY-3`, `BP-PY-5` were optional/deferred in parent Phase 2 — owned by other batches if scheduled; **not** batch-00 work

---

## Phase 2: Security — `rules_security.go`

### 2.1 Registered rules (evidence)

- [x] `BP-PY-8` subprocess `shell=True` — `rules_security.go` + tests
- [x] `BP-PY-9` `os.system` / `os.popen` — `rules_security.go` + tests
- [x] `BP-PY-10` pickle loads — `rules_security.go` + tests
- [x] `BP-PY-11` `yaml.load` without SafeLoader — `rules_security.go` + tests
- [x] `BP-PY-12` eval/exec — `rules_security.go` + tests
- [x] `BP-PY-13` hardcoded secret — `rules_security.go` + tests

---

## Phase 3: Framework high-signal — `rules_framework.go`

### 3.1 Registered rules (evidence)

- [x] `BP-PY-16` Flask DEBUG — `rules_framework.go` + tests (skips test files)
- [x] `BP-PY-17` Flask SECRET_KEY literal — `rules_framework.go` + tests
- [x] `BP-PY-21` Django DEBUG in settings — `rules_framework.go` + tests

### 3.2 Not shipped here (other batches)

- [x] Note only: `BP-PY-18`–`20`, `22`–`37`, etc. remain missing per `_inventory.json` — implement under their owner batch plans, **not** by reopening batch-00

---

## Phase 4: Registration smoke (already green)

- [x] `TestBPRulesRegistered` expects at least: `1,2,4,6,7,8–13,16,17,21` — `scan_test.go`
- [x] Collision guard: no bare `BP-<n>` Go IDs registered
- [x] Pack = `PackBadPractice` for each registered id
- [x] Prior validation on implement branch: `make lint` / `make test` green (see parent `python-heuristics-bp.md` Phase 6 evidence, 2026-07-31)

---

## Remaining work (pointer only)

| Next | Plan | IDs |
|------|------|-----|
| Async | [batch-06-async.md](./batch-06-async.md) | 38, 39, 40 |
| Testing / deps / obs | [batch-07-testing-deps-obs.md](./batch-07-testing-deps-obs.md) | 41–47 |
| Prod hardening | [batch-08-prod-hardening.md](./batch-08-prod-hardening.md) | 48, 49, 50 (**not** 14) |
| HTTP timeout etc. | sibling batch-01 (if present) | **BP-PY-14** lives there |
| Framework / DB / etc. | sibling batches 02–05 | remaining missing IDs |

---

## Dependencies

None — ledger only. Implementers of other batches depend on this inventory for **exclusion**.

---

## ID list (complete for this file)

**Shipped [x]:** BP-PY-1, BP-PY-2, BP-PY-4, BP-PY-6, BP-PY-7, BP-PY-8, BP-PY-9, BP-PY-10, BP-PY-11, BP-PY-12, BP-PY-13, BP-PY-16, BP-PY-17, BP-PY-21
