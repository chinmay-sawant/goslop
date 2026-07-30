# Batch 00 — Shipped Priority

> **Parent:** `plans/v0.0.2/heuristics/cwe-plans/README.md` — v0.0.2 Python CWE remaining heuristics
> **Epic / issue:** [#51](https://github.com/chinmay-sawant/goslop/issues/51) · [#52](https://github.com/chinmay-sawant/goslop/issues/52) expansion
> **Status:** **complete** — already on `main` (ledger only)
> **IDs (5):** CWE-22, CWE-78, CWE-79, CWE-89, CWE-502
> **PR policy:** one PR for this batch only — do not mix other wave IDs

---

## Architecture constraints

| Constraint | Detail |
|------------|--------|
| Package | `internal/lang/python/detectors/cwe/` only |
| Registration | `RegisterRule("CWE-*", detect…, &Meta…, gates…)` from `init()` |
| Scan | Existing `PyCweScan` — no new detector type |
| Detection | Pure-Go source patterns + needles / `PyCweFacts` |
| Language | `LanguagePython` gate already in scan |
| Plugin | **Do NOT invent a second plugin.** `detectors.All()` already wires CWE scan |
| IDs | Always `CWE-*`; metadata aligned with chunk catalogue |
| **File size policy** | Target max **1500** / hard max **2000** per Go file |
| Target file(s) | `rules.go (split to domain files over time)` |
| Validation | `make lint` + `make test` unchecked until green |
| Fixtures | `tests/fixtures/python/cwe/CWE-N-{{vulnerable,safe}}.txt` |

## Overview

Priority #52 batch already RegisterRule on main

| Rule | Name | Tier | Category | Relevance | Chunk |
|------|------|------|----------|-----------|-------|
| `CWE-22` | Improper Limitation of a Pathname to a Restricted Directory ('Path Traversal') | shipped | File Processing | Medium | `cwe-001-050.json` |
| `CWE-78` | Improper Neutralization of Special Elements used in an OS Command ('OS Command Injection') | shipped | Program Invocation | Medium | `cwe-051-100.json` |
| `CWE-79` | Improper Neutralization of Input During Web Page Generation ('Cross-site Scripting') | shipped | Injection | Medium | `cwe-051-100.json` |
| `CWE-89` | Improper Neutralization of Special Elements used in an SQL Command ('SQL Injection') | shipped | Injection | Medium | `cwe-051-100.json` |
| `CWE-502` | Deserialization of Untrusted Data | shipped | Deserialization | High | `cwe-501-550.json` |

## Executive Summary

These five rules are **already implemented** on `main` (priority #52 batch). This file is a **closed ledger** — do not re-open as pending work. Expand via sibling batches only.

## Shipped proof

- [x] `CWE-22` registered in `rules.go` `init()` + fixtures under `tests/fixtures/python/cwe/`
- [x] `CWE-78` registered in `rules.go` `init()` + fixtures under `tests/fixtures/python/cwe/`
- [x] `CWE-79` registered in `rules.go` `init()` + fixtures under `tests/fixtures/python/cwe/`
- [x] `CWE-89` registered in `rules.go` `init()` + fixtures under `tests/fixtures/python/cwe/`
- [x] `CWE-502` registered in `rules.go` `init()` + fixtures under `tests/fixtures/python/cwe/`
- [x] `go test ./internal/lang/python/detectors/cwe/` + integration matrix

## Do not implement here

Any remaining catalogue IDs belong to batch-01+ or `batch-deferred.md`.

## Shipped fixtures + tests (reference pattern)

These are the gold standard for expansion batches (same layout as BP under `tests/fixtures/python/bp/`).

### Fixture pairs (`tests/fixtures/python/cwe/`)

| Rule | Vulnerable | Safe |
|------|------------|------|
| CWE-22 | `CWE-22-vulnerable.txt` | `CWE-22-safe.txt` |
| CWE-78 | `CWE-78-vulnerable.txt` | `CWE-78-safe.txt` |
| CWE-79 | `CWE-79-vulnerable.txt` | `CWE-79-safe.txt` |
| CWE-89 | `CWE-89-vulnerable.txt` | `CWE-89-safe.txt` |
| CWE-502 | `CWE-502-vulnerable.txt` | `CWE-502-safe.txt` |

### Unit tests

- [x] `internal/lang/python/detectors/cwe/rules_test.go` (+ `scan_test.go`) hit/miss for each of the five

### Integration

- [x] `tests/integration/python/cwe_matrix_test.go` — `TestPythonCWEFixturesMatrix`
- [x] Discovery via `DiscoverPythonCWECases()` over `tests/fixtures/python/cwe/`
- [x] `make integration-python` / `go test ./tests/integration/python/`

