# Batch 11 — Info Exposure

> **Parent:** `plans/v0.0.2/heuristics/cwe-plans/README.md` — v0.0.2 Python CWE remaining heuristics
> **Epic / issue:** [#51](https://github.com/chinmay-sawant/goslop/issues/51) · [#52](https://github.com/chinmay-sawant/goslop/issues/52) expansion
> **Status:** **pending** — not started
> **Wave:** P2
> **IDs (8):** CWE-209, CWE-208, CWE-201, CWE-204, CWE-212, CWE-213, CWE-497, CWE-488
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
| Target file(s) | `rules_info_exposure.go` |
| Validation | `make lint` + `make test` unchecked until green |
| Fixtures (text) | **Required:** `tests/fixtures/python/cwe/CWE-N-{vulnerable,safe}.txt` (BP-style; see Testing section) |
| Unit tests | **Required:** hit/miss in `internal/lang/python/detectors/cwe/*_test.go` |
| Integration | **Required:** `tests/integration/python/cwe_matrix_test.go` / `make integration-python` |

## Overview

Error messages, timing, debug, sensitive data in responses/logs

| Rule | Name | Tier | Category | Relevance | Chunk |
|------|------|------|----------|-----------|-------|
| `CWE-209` | Generation of Error Message Containing Sensitive Information | A | Information Disclosure | Medium | `cwe-201-250.json` |
| `CWE-208` | Observable Timing Discrepancy | A | Cryptography | Medium | `cwe-201-250.json` |
| `CWE-201` | Insertion of Sensitive Information Into Sent Data | B | Information Disclosure | Medium | `cwe-201-250.json` |
| `CWE-204` | Observable Response Discrepancy | B | General | Medium | `cwe-201-250.json` |
| `CWE-212` | Improper Removal of Sensitive Information Before Storage or Transfer | B | Information Disclosure | Medium | `cwe-201-250.json` |
| `CWE-213` | Exposure of Sensitive Information Due to Incompatible Policies | B | Information Disclosure | Medium | `cwe-201-250.json` |
| `CWE-497` | Exposure of Sensitive System Information to an Unauthorized Control Sphere | B | Information Disclosure | Medium | `cwe-451-500.json` |
| `CWE-488` | Exposure of Data Element to Wrong Session | B | Security | Medium | `cwe-451-500.json` |

## Executive Summary

Ship **8** remaining CWE heuristics in this theme. Prefer high-signal, low-FP patterns; same-file heuristics only (no interproc taint). Extend registration catalogue tests and fixture matrix.

## Phase 0: Placement + budget

- [ ] Record baseline: `wc -l internal/lang/python/detectors/cwe/*.go`
- [ ] Ensure target `rules_info_exposure.go` exists (create if needed); keep ≤1500 lines projected
- [ ] Prefer domain file over growing `rules.go` past soft cap
- [ ] Append FN-safe needles to `needles.go` / per-rule gates
- [ ] No changes to `detectors/all.go` unless constructor rename (should be none)

## Phase 1: `CWE-209` — Generation of Error Message Containing Sensitive Information

> Catalogue: `ruleset/python/chunks/cwe-201-250.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE209` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-209", detectCWE209, &Meta…, gates...)` in `rules_info_exposure.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-209`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Manual Analysis, Automated Analysis, Automated Dynamic Analysis, Manual Dynamic Analysis, Automated Static Analysis. This weakness generally requires domain-specific interpretation using manual a…
- Suggested sinks/patterns: `traceback.format_exc()`, `str(e) in Response`, `Flask debug`, `django.views.debug`, `return jsonify(error=str(e))`
- [ ] Implement `detectCWE209` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-209`
- [ ] Unit miss: safe pattern → no `CWE-209`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-209-vulnerable.txt` + `CWE-209-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [ ] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [ ] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 2: `CWE-208` — Observable Timing Discrepancy

> Catalogue: `ruleset/python/chunks/cwe-201-250.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE208` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-208", detectCWE208, &Meta…, gates...)` in `rules_info_exposure.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-208`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `password == expected`, `secret == token`, `hmac.compare_digest missing`, `secrets.compare_digest`
- [ ] Implement `detectCWE208` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-208`
- [ ] Unit miss: safe pattern → no `CWE-208`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-208-vulnerable.txt` + `CWE-208-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [ ] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [ ] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 3: `CWE-201` — Insertion of Sensitive Information Into Sent Data

> Catalogue: `ruleset/python/chunks/cwe-201-250.json` · tier **B** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE201` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-201", detectCWE201, &Meta…, gates...)` in `rules_info_exposure.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-201`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `jsonify(user_model)`, `password field in response`, `serialize full ORM object`
- [ ] Implement `detectCWE201` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-201`
- [ ] Unit miss: safe pattern → no `CWE-201`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-201-vulnerable.txt` + `CWE-201-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [ ] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [ ] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 4: `CWE-204` — Observable Response Discrepancy

> Catalogue: `ruleset/python/chunks/cwe-201-250.json` · tier **B** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE204` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-204", detectCWE204, &Meta…, gates...)` in `rules_info_exposure.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-204`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `different auth error messages`, `user not found vs wrong password`
- [ ] Implement `detectCWE204` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-204`
- [ ] Unit miss: safe pattern → no `CWE-204`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-204-vulnerable.txt` + `CWE-204-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [ ] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [ ] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 5: `CWE-212` — Improper Removal of Sensitive Information Before Storage or Transfer

> Catalogue: `ruleset/python/chunks/cwe-201-250.json` · tier **B** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE212` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-212", detectCWE212, &Meta…, gates...)` in `rules_info_exposure.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-212`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Tools are available to analyze documents (such as PDF, Word, etc.) to look for private information such as names, addresses, etc. Python-oriented (future detector): pre…
- Suggested sinks/patterns: `export full model with PAN/ssn`, `redact missing before dump`
- [ ] Implement `detectCWE212` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-212`
- [ ] Unit miss: safe pattern → no `CWE-212`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-212-vulnerable.txt` + `CWE-212-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [ ] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [ ] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 6: `CWE-213` — Exposure of Sensitive Information Due to Incompatible Policies

> Catalogue: `ruleset/python/chunks/cwe-201-250.json` · tier **B** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE213` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-213", detectCWE213, &Meta…, gates...)` in `rules_info_exposure.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-213`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `serialize salary/PII to guest API`
- [ ] Implement `detectCWE213` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-213`
- [ ] Unit miss: safe pattern → no `CWE-213`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-213-vulnerable.txt` + `CWE-213-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [ ] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [ ] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 7: `CWE-497` — Exposure of Sensitive System Information to an Unauthorized Control Sphere

> Catalogue: `ruleset/python/chunks/cwe-451-500.json` · tier **B** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE497` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-497", detectCWE497, &Meta…, gates...)` in `rules_info_exposure.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-497`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `traceback in HTTP response`, `app.config expose`, `sys.path / env dumped`, `Flask debug error pages in prod`
- [ ] Implement `detectCWE497` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-497`
- [ ] Unit miss: safe pattern → no `CWE-497`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-497-vulnerable.txt` + `CWE-497-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [ ] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [ ] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 8: `CWE-488` — Exposure of Data Element to Wrong Session

> Catalogue: `ruleset/python/chunks/cwe-451-500.json` · tier **B** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE488` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-488", detectCWE488, &Meta…, gates...)` in `rules_info_exposure.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-488`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `mutable module-level user state`, `global current_user`, `Flask g misuse across requests`
- [ ] Implement `detectCWE488` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-488`
- [ ] Unit miss: safe pattern → no `CWE-488`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-488-vulnerable.txt` + `CWE-488-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
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

