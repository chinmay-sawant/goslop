# Batch 12 — Validation Export Redos

> **Parent:** `plans/v0.0.2/heuristics/cwe-plans/README.md` — v0.0.2 Python CWE remaining heuristics
> **Epic / issue:** [#51](https://github.com/chinmay-sawant/goslop/issues/51) · [#52](https://github.com/chinmay-sawant/goslop/issues/52) expansion
> **Status:** **pending** — not started
> **Wave:** P2
> **IDs (8):** CWE-1173, CWE-1230, CWE-1236, CWE-1286, CWE-1289, CWE-1333, CWE-1389, CWE-140
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
| Target file(s) | `rules_validation.go` |
| Validation | `make lint` + `make test` unchecked until green |
| Fixtures (text) | **Required:** `tests/fixtures/python/cwe/CWE-N-{vulnerable,safe}.txt` (BP-style; see Testing section) |
| Unit tests | **Required:** hit/miss in `internal/lang/python/detectors/cwe/*_test.go` |
| Integration | **Required:** `tests/integration/python/cwe_matrix_test.go` / `make integration-python` |

## Overview

Input validation gaps, CSV formula injection, ReDoS, numeric parse

| Rule | Name | Tier | Category | Relevance | Chunk |
|------|------|------|----------|-----------|-------|
| `CWE-1173` | Improper Use of Validation Framework | A | General | Medium | `cwe-1151-1200.json` |
| `CWE-1230` | Exposure of Sensitive Information Through Metadata | A | Information Disclosure | Medium | `cwe-1201-1250.json` |
| `CWE-1236` | Improper Neutralization of Formula Elements in a CSV File | A | General | Medium | `cwe-1201-1250.json` |
| `CWE-1286` | Improper Validation of Syntactic Correctness of Input | A | General | Medium | `cwe-1251-1300.json` |
| `CWE-1289` | Improper Validation of Unsafe Equivalence in Input | A | General | Medium | `cwe-1251-1300.json` |
| `CWE-1333` | Inefficient Regular Expression Complexity | A | General | Medium | `cwe-1301-1350.json` |
| `CWE-1389` | Incorrect Parsing of Numbers with Different Radices | A | General | Medium | `cwe-1351-1400.json` |
| `CWE-140` | Improper Neutralization of Delimiters | B | General | Medium | `cwe-101-150.json` |

## Executive Summary

Ship **8** remaining CWE heuristics in this theme. Prefer high-signal, low-FP patterns; same-file heuristics only (no interproc taint). Extend registration catalogue tests and fixture matrix.

## Phase 0: Placement + budget

- [ ] Record baseline: `wc -l internal/lang/python/detectors/cwe/*.go`
- [ ] Ensure target `rules_validation.go` exists (create if needed); keep ≤1500 lines projected
- [ ] Prefer domain file over growing `rules.go` past soft cap
- [ ] Append FN-safe needles to `needles.go` / per-rule gates
- [ ] No changes to `detectors/all.go` unless constructor rename (should be none)

## Phase 1: `CWE-1173` — Improper Use of Validation Framework

> Catalogue: `ruleset/python/chunks/cwe-1151-1200.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE1173` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-1173", detectCWE1173, &Meta…, gates...)` in `rules_validation.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-1173`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Some instances of improper input validation can be detected using automated static analysis. A static analysis tool might allow the user to specify which application-sp…
- Suggested sinks/patterns: `request.get_json / request.json used without schema`, `dict from body instead of pydantic/marshmallow/DRF serializer`, `Django form skipped for Model.save`
- [ ] Implement `detectCWE1173` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-1173`
- [ ] Unit miss: safe pattern → no `CWE-1173`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-1173-vulnerable.txt` + `CWE-1173-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [ ] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [ ] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 2: `CWE-1230` — Exposure of Sensitive Information Through Metadata

> Catalogue: `ruleset/python/chunks/cwe-1201-1250.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE1230` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-1230", detectCWE1230, &Meta…, gates...)` in `rules_validation.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-1230`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `Content-Disposition filename=user`, `headers leaking original names/sizes`, `Server/X-Powered-By verbose`
- [ ] Implement `detectCWE1230` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-1230`
- [ ] Unit miss: safe pattern → no `CWE-1230`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-1230-vulnerable.txt` + `CWE-1230-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [ ] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [ ] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 3: `CWE-1236` — Improper Neutralization of Formula Elements in a CSV File

> Catalogue: `ruleset/python/chunks/cwe-1201-1250.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE1236` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-1236", detectCWE1236, &Meta…, gates...)` in `rules_validation.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-1236`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `csv.writer writerow user fields`, `f-string CSV without stripping =+-@`
- [ ] Implement `detectCWE1236` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-1236`
- [ ] Unit miss: safe pattern → no `CWE-1236`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-1236-vulnerable.txt` + `CWE-1236-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [ ] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [ ] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 4: `CWE-1286` — Improper Validation of Syntactic Correctness of Input

> Catalogue: `ruleset/python/chunks/cwe-1251-1300.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE1286` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-1286", detectCWE1286, &Meta…, gates...)` in `rules_validation.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-1286`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `json.loads then use URLs without urlparse validation`, `yaml.safe_load config without schema`
- [ ] Implement `detectCWE1286` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-1286`
- [ ] Unit miss: safe pattern → no `CWE-1286`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-1286-vulnerable.txt` + `CWE-1286-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [ ] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [ ] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 5: `CWE-1289` — Improper Validation of Unsafe Equivalence in Input

> Catalogue: `ruleset/python/chunks/cwe-1251-1300.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE1289` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-1289", detectCWE1289, &Meta…, gates...)` in `rules_validation.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-1289`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `blocklist path == 'private/x' before normalize`, `string equality deny-list for paths`
- [ ] Implement `detectCWE1289` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-1289`
- [ ] Unit miss: safe pattern → no `CWE-1289`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-1289-vulnerable.txt` + `CWE-1289-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [ ] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [ ] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 6: `CWE-1333` — Inefficient Regular Expression Complexity

> Catalogue: `ruleset/python/chunks/cwe-1301-1350.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE1333` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-1333", detectCWE1333, &Meta…, gates...)` in `rules_validation.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-1333`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `re.compile nested quantifiers (a+)+`, `evil regex on user input`, `django URL/validators catastrophic patterns`
- [ ] Implement `detectCWE1333` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-1333`
- [ ] Unit miss: safe pattern → no `CWE-1333`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-1333-vulnerable.txt` + `CWE-1333-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [ ] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [ ] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 7: `CWE-1389` — Incorrect Parsing of Numbers with Different Radices

> Catalogue: `ruleset/python/chunks/cwe-1351-1400.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE1389` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-1389", detectCWE1389, &Meta…, gates...)` in `rules_validation.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-1389`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `int(s, 0)`, `int(user_input, 0)`
- [ ] Implement `detectCWE1389` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-1389`
- [ ] Unit miss: safe pattern → no `CWE-1389`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-1389-vulnerable.txt` + `CWE-1389-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [ ] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [ ] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 8: `CWE-140` — Improper Neutralization of Delimiters

> Catalogue: `ruleset/python/chunks/cwe-101-150.json` · tier **B** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE140` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-140", detectCWE140, &Meta…, gates...)` in `rules_validation.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-140`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `",".join(`, `csv without csv.writer`, `manual delimiter concat`
- [ ] Implement `detectCWE140` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-140`
- [ ] Unit miss: safe pattern → no `CWE-140`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-140-vulnerable.txt` + `CWE-140-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [ ] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [ ] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 9: Batch validation

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

