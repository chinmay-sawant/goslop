# Batch 13 — Platform Quality

> **Parent:** `plans/v0.0.2/heuristics/cwe-plans/README.md` — v0.0.2 Python CWE remaining heuristics
> **Epic / issue:** [#51](https://github.com/chinmay-sawant/goslop/issues/51) · [#52](https://github.com/chinmay-sawant/goslop/issues/52) expansion
> **Status:** **pending** — not started
> **Wave:** P3
> **IDs (6):** CWE-396, CWE-397, CWE-478, CWE-252, CWE-390, CWE-584
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
| Target file(s) | `rules_platform.go` |
| Validation | `make lint` + `make test` unchecked until green |
| Fixtures (text) | **Required:** `tests/fixtures/python/cwe/CWE-N-{vulnerable,safe}.txt` (BP-style; see Testing section) |
| Unit tests | **Required:** hit/miss in `internal/lang/python/detectors/cwe/*_test.go` |
| Integration | **Required:** `tests/integration/python/cwe_matrix_test.go` / `make integration-python` |

## Overview

High python_relevance error-handling (396/397) + quality sinks

| Rule | Name | Tier | Category | Relevance | Chunk |
|------|------|------|----------|-----------|-------|
| `CWE-396` | Declaration of Catch for Generic Exception | B | Error Handling | High | `cwe-351-400.json` |
| `CWE-397` | Declaration of Throws for Generic Exception | B | Error Handling | High | `cwe-351-400.json` |
| `CWE-478` | Missing Default Case in Multiple Condition Expression | C | Configuration | High | `cwe-451-500.json` |
| `CWE-252` | Unchecked Return Value | B | General | Medium | `cwe-251-300.json` |
| `CWE-390` | Detection of Error Condition Without Action | B | Error Handling | Medium | `cwe-351-400.json` |
| `CWE-584` | Return Inside Finally Block | A | Resource Management | Medium | `cwe-551-600.json` |

## Executive Summary

Ship **6** remaining CWE heuristics in this theme. Prefer high-signal, low-FP patterns; same-file heuristics only (no interproc taint). Extend registration catalogue tests and fixture matrix.

## Phase 0: Placement + budget

- [ ] Record baseline: `wc -l internal/lang/python/detectors/cwe/*.go`
- [ ] Ensure target `rules_platform.go` exists (create if needed); keep ≤1500 lines projected
- [ ] Prefer domain file over growing `rules.go` past soft cap
- [ ] Append FN-safe needles to `needles.go` / per-rule gates
- [ ] No changes to `detectors/all.go` unless constructor rename (should be none)

## Phase 1: `CWE-396` — Declaration of Catch for Generic Exception

> Catalogue: `ruleset/python/chunks/cwe-351-400.json` · tier **B** · relevance **High**

### Register + meta

- [ ] Add `MetaCWE396` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-396", detectCWE396, &Meta…, gates...)` in `rules_platform.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-396`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `except Exception:`, `except BaseException:`
- [ ] Implement `detectCWE396` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-396`
- [ ] Unit miss: safe pattern → no `CWE-396`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-396-vulnerable.txt` + `CWE-396-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [ ] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [ ] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 2: `CWE-397` — Declaration of Throws for Generic Exception

> Catalogue: `ruleset/python/chunks/cwe-351-400.json` · tier **B** · relevance **High**

### Register + meta

- [ ] Add `MetaCWE397` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-397", detectCWE397, &Meta…, gates...)` in `rules_platform.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-397`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `raise Exception(`, `raise BaseException`
- [ ] Implement `detectCWE397` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-397`
- [ ] Unit miss: safe pattern → no `CWE-397`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-397-vulnerable.txt` + `CWE-397-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [ ] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [ ] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 3: `CWE-478` — Missing Default Case in Multiple Condition Expression

> Catalogue: `ruleset/python/chunks/cwe-451-500.json` · tier **C** · relevance **High**

### Register + meta

- [ ] Add `MetaCWE478` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-478", detectCWE478, &Meta…, gates...)` in `rules_platform.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-478`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `match/case without _`
- [ ] Implement `detectCWE478` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-478`
- [ ] Unit miss: safe pattern → no `CWE-478`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-478-vulnerable.txt` + `CWE-478-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [ ] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [ ] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 4: `CWE-252` — Unchecked Return Value

> Catalogue: `ruleset/python/chunks/cwe-251-300.json` · tier **B** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE252` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-252", detectCWE252, &Meta…, gates...)` in `rules_platform.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-252`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `ignored os.* returns`, `pathlib ops without try`, `subprocess without check`
- [ ] Implement `detectCWE252` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-252`
- [ ] Unit miss: safe pattern → no `CWE-252`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-252-vulnerable.txt` + `CWE-252-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [ ] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [ ] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 5: `CWE-390` — Detection of Error Condition Without Action

> Catalogue: `ruleset/python/chunks/cwe-351-400.json` · tier **B** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE390` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-390", detectCWE390, &Meta…, gates...)` in `rules_platform.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-390`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `except Exception: pass`, `except: pass`
- [ ] Implement `detectCWE390` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-390`
- [ ] Unit miss: safe pattern → no `CWE-390`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-390-vulnerable.txt` + `CWE-390-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [ ] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [ ] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 6: `CWE-584` — Return Inside Finally Block

> Catalogue: `ruleset/python/chunks/cwe-551-600.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE584` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-584", detectCWE584, &Meta…, gates...)` in `rules_platform.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-584`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `finally:\n    ... return `, `finally: return`
- [ ] Implement `detectCWE584` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-584`
- [ ] Unit miss: safe pattern → no `CWE-584`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-584-vulnerable.txt` + `CWE-584-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [ ] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [ ] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 7: Batch validation

- [ ] `gofmt -w` on touched files
- [ ] `make lint`
- [ ] `make test`
- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` (individual unit tests)
- [ ] `go test ./tests/integration/python/ -count=1` **or** `make integration-python` (CWE matrix over all `python/cwe` pairs)
- [ ] Confirm every batch ID has both `CWE-N-vulnerable.txt` and `CWE-N-safe.txt` under `tests/fixtures/python/cwe/`
- [ ] Fixture count: `DiscoverPythonCWECases` includes all new IDs (pair discovery, not a manual allowlist)
- [ ] Update `_inventory.json`: move batch IDs from `missing` → `implemented`
- [ ] Update this ledger statuses to `[x]` with evidence
- [ ] Package files still ≤2000 lines (split if not)


## Testing requirements (fixtures + unit + integration)

> Same bar as BP-PY: detector alone is **not** enough. Each owned `CWE-N` needs the triad below.

### Fixture text files (required)

| Path | Purpose |
|------|---------|
| `tests/fixtures/python/cwe/CWE-N-vulnerable.txt` | Hit corpus — scan must emit `CWE-N` |
| `tests/fixtures/python/cwe/CWE-N-safe.txt` | Miss corpus — scan must **not** emit `CWE-N` |

**Format** (see `tests/fixtures/README.md` and BP examples under `tests/fixtures/python/bp/`):

```text
# Fixture for CWE-N (vulnerable|safe)
lang: python
file: CWE-N-vulnerable.py   # or -safe.py
---
# python source body
```

- Keep fixtures under **`python/cwe/`** only (parallel to **`python/bp/`** for BP-PY).
- Do **not** commit `.py` sources; discovery is pair-based (`DiscoverPythonCWECases`).

### Individual unit tests (required)

| Path | Purpose |
|------|---------|
| `internal/lang/python/detectors/cwe/rules_test.go` (or domain `rules_*_test.go` / `scan_test.go`) | Hit + miss table tests per ID |
| Registration catalogue test | Want-list includes every new `CWE-N` from this batch |

```sh
go test ./internal/lang/python/detectors/cwe/ -count=1
```

### Integration matrix (required)

| Path | Purpose |
|------|---------|
| `tests/integration/python/cwe_matrix_test.go` | `TestPythonCWEFixturesMatrix` — scans every discovered pair |
| Helpers | `integration.DiscoverPythonCWECases`, `PythonCWEFixtureRel` in `tests/integration/discover.go` |

```sh
go test ./tests/integration/python/ -count=1
# or
make integration-python
```

BP analogue (do not mix): `tests/fixtures/python/bp/` + `bp_matrix_test.go` + `DiscoverPythonBPCases`.

### Per-ID checklist (repeat for every rule in this batch)

- [ ] Unit hit/miss for `CWE-N`
- [ ] `tests/fixtures/python/cwe/CWE-N-vulnerable.txt`
- [ ] `tests/fixtures/python/cwe/CWE-N-safe.txt`
- [ ] Matrix auto-discovers pair; vulnerable asserts finding; safe asserts absence
- [ ] `make lint` + `make test` + `make integration-python` green before PR merge

## Dependencies

| Depends on | Note |
|------------|------|
| batch-00 | Framework + priority rules already present |
| Catalogue chunks | IDs must exist in `ruleset/python/chunks/` |
| Parent README | ownership + PR policy |

