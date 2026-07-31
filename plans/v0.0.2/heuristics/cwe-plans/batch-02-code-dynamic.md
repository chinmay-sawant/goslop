# Batch 02 — Code Dynamic

> **Parent:** `plans/v0.0.2/heuristics/cwe-plans/README.md` — v0.0.2 Python CWE remaining heuristics
> **Epic / issue:** [#51](https://github.com/chinmay-sawant/goslop/issues/51) · [#52](https://github.com/chinmay-sawant/goslop/issues/52) expansion
> **Status:** **complete** — implemented and validated 2026-07-31
> **Validation evidence:** `gofmt`, Python CWE package tests, Python CWE fixture matrix, `make lint`, `make test`, and `git diff --check` passed.
> **Wave:** P0
> **IDs (5):** CWE-749, CWE-829, CWE-695, CWE-214, CWE-215
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
| Target file(s) | `rules_code_dynamic.go` |
| Validation | `make lint` + `make test` unchecked until green |
| Fixtures (text) | **Required:** `tests/fixtures/python/cwe/CWE-N-{vulnerable,safe}.txt` (BP-style; see Testing section) |
| Unit tests | **Required:** hit/miss in `internal/lang/python/detectors/cwe/*_test.go` |
| Integration | **Required:** `tests/integration/python/cwe_matrix_test.go` / `make integration-python` |

## Overview

eval/exec/importlib/dynamic code inclusion; dangerous reflection APIs

| Rule | Name | Tier | Category | Relevance | Chunk |
|------|------|------|----------|-----------|-------|
| `CWE-749` | Exposed Dangerous Method or Function | A | General | Medium | `cwe-701-750.json` |
| `CWE-829` | Inclusion of Functionality from Untrusted Control Sphere | A | General | Medium | `cwe-801-850.json` |
| `CWE-695` | Use of Low-Level Functionality | B | General | Medium | `cwe-651-700.json` |
| `CWE-214` | Invocation of Process Using Visible Sensitive Information | A | Program Invocation | Medium | `cwe-201-250.json` |
| `CWE-215` | Insertion of Sensitive Information Into Debugging Code | A | Information Disclosure | Medium | `cwe-201-250.json` |

## Executive Summary

Ship **5** remaining CWE heuristics in this theme. Prefer high-signal, low-FP patterns; same-file heuristics only (no interproc taint). Extend registration catalogue tests and fixture matrix.

## Phase 0: Placement + budget

- [x] Record baseline: `wc -l internal/lang/python/detectors/cwe/*.go`
- [x] Ensure target `rules_code_dynamic.go` exists (create if needed); keep ≤1500 lines projected
- [x] Prefer domain file over growing `rules.go` past soft cap
- [x] Append FN-safe needles to `needles.go` / per-rule gates
- [x] No changes to `detectors/all.go` unless constructor rename (should be none)

## Phase 1: `CWE-749` — Exposed Dangerous Method or Function

> Catalogue: `ruleset/python/chunks/cwe-701-750.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE749` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-749", detectCWE749, &Meta…, gates...)` in `rules_code_dynamic.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-749`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `eval(`, `exec(`, `compile(`, `__import__(`, `importlib.import_module(dynamic)`, `getattr(__builtins__`, `os.system exposed via HTTP handlers`
- [x] Implement `detectCWE749` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-749`
- [x] Unit miss: safe pattern → no `CWE-749`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-749-vulnerable.txt` + `CWE-749-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 2: `CWE-829` — Inclusion of Functionality from Untrusted Control Sphere

> Catalogue: `ruleset/python/chunks/cwe-801-850.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE829` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-829", detectCWE829, &Meta…, gates...)` in `rules_code_dynamic.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-829`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis - Binary or Bytecode, Manual Static Analysis - Binary or Bytecode, Dynamic Analysis with Manual Results Interpretation, Manual Static Analysis - Source Code, Automated S…
- Suggested sinks/patterns: `__import__(dynamic)`, `importlib.import_module(dynamic)`, `importlib.util.spec_from_file_location`, `runpy`, `exec(open(...).read())`, `eval(`, `plugin load from user path`
- [x] Implement `detectCWE829` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-829`
- [x] Unit miss: safe pattern → no `CWE-829`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-829-vulnerable.txt` + `CWE-829-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 3: `CWE-695` — Use of Low-Level Functionality

> Catalogue: `ruleset/python/chunks/cwe-651-700.json` · tier **B** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE695` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-695", detectCWE695, &Meta…, gates...)` in `rules_code_dynamic.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-695`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `ctypes.`, `cffi`, `mmap.`, `os.system already CWE-78`
- [x] Implement `detectCWE695` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-695`
- [x] Unit miss: safe pattern → no `CWE-695`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-695-vulnerable.txt` + `CWE-695-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 4: `CWE-214` — Invocation of Process Using Visible Sensitive Information

> Catalogue: `ruleset/python/chunks/cwe-201-250.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE214` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-214", detectCWE214, &Meta…, gates...)` in `rules_code_dynamic.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-214`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `subprocess.*(password|token|secret|api_key in argv)`, `os.system with secret`, `Popen env secrets`
- [x] Implement `detectCWE214` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-214`
- [x] Unit miss: safe pattern → no `CWE-214`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-214-vulnerable.txt` + `CWE-214-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 5: `CWE-215` — Insertion of Sensitive Information Into Debugging Code

> Catalogue: `ruleset/python/chunks/cwe-201-250.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE215` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-215", detectCWE215, &Meta…, gates...)` in `rules_code_dynamic.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-215`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `print(password`, `logging.debug(...secret`, `pdb`, `print(SECRET_KEY`, `logger.debug api_key`
- [x] Implement `detectCWE215` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-215`
- [x] Unit miss: safe pattern → no `CWE-215`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-215-vulnerable.txt` + `CWE-215-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 6: Batch validation

- [x] `gofmt -w` on touched files
- [x] `make lint`
- [x] `make test`
- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` (individual unit tests)
- [x] `go test ./tests/integration/python/ -count=1` **or** `make integration-python` (CWE matrix over all `python/cwe` pairs)
- [x] Confirm every batch ID has both `CWE-N-vulnerable.txt` and `CWE-N-safe.txt` under `tests/fixtures/python/cwe/`
- [x] Fixture count: `DiscoverPythonCWECases` includes all new IDs (pair discovery, not a manual allowlist)
- [x] Update `_inventory.json`: move batch IDs from `missing` → `implemented`
- [x] Update this ledger statuses to `[x]` with evidence
- [x] Package files still ≤2000 lines (split if not)


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

- [x] Unit hit/miss for `CWE-N`
- [x] `tests/fixtures/python/cwe/CWE-N-vulnerable.txt`
- [x] `tests/fixtures/python/cwe/CWE-N-safe.txt`
- [x] Matrix auto-discovers pair; vulnerable asserts finding; safe asserts absence
- [x] `make lint` + `make test` + `make integration-python` green before PR merge

## Dependencies

| Depends on | Note |
|------------|------|
| batch-00 | Framework + priority rules already present |
| Catalogue chunks | IDs must exist in `ruleset/python/chunks/` |
| Parent README | ownership + PR policy |
