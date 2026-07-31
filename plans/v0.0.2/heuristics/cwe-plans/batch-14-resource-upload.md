# Batch 14 — Resource Upload

> **Parent:** `plans/v0.0.2/heuristics/cwe-plans/README.md` — v0.0.2 Python CWE remaining heuristics
> **Epic / issue:** [#51](https://github.com/chinmay-sawant/goslop/issues/51) · [#52](https://github.com/chinmay-sawant/goslop/issues/52) expansion
> **Status:** **complete** — implemented and validated 2026-07-31
> **Validation evidence:** fixture-only per-rule tests, Python CWE fixture matrix, `gofmt`, `make lint`, `make test`, and `git diff --check` passed.
> **Wave:** P3
> **IDs (8):** CWE-434, CWE-427, CWE-379, CWE-459, CWE-772, CWE-770, CWE-708, CWE-477
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
| Target file(s) | `rules_resource.go` |
| Validation | `make lint` + `make test` unchecked until green |
| Fixtures (text) | **Required:** `tests/fixtures/python/cwe/CWE-N-{vulnerable,safe}.txt` (BP-style; see Testing section) |
| Unit tests | **Required:** hit/miss in `internal/lang/python/detectors/cwe/*_test.go` |
| Integration | **Required:** `tests/integration/python/cwe_matrix_test.go` / `make integration-python` |

## Overview

Upload, resource lifecycle, secondary FS hygiene

| Rule | Name | Tier | Category | Relevance | Chunk |
|------|------|------|----------|-----------|-------|
| `CWE-434` | Unrestricted Upload of File with Dangerous Type | B | File Processing | Medium | `cwe-401-450.json` |
| `CWE-427` | Uncontrolled Search Path Element | B | Program Invocation | Medium | `cwe-401-450.json` |
| `CWE-379` | Creation of Temporary File in Directory with Insecure Permissions | B | Path Traversal | Medium | `cwe-351-400.json` |
| `CWE-459` | Incomplete Cleanup | B | File Processing | Medium | `cwe-451-500.json` |
| `CWE-772` | Missing Release of Resource after Effective Lifetime | B | Resource Management | Medium | `cwe-751-800.json` |
| `CWE-770` | Allocation of Resources Without Limits or Throttling | B | Resource Management | Medium | `cwe-751-800.json` |
| `CWE-708` | Incorrect Ownership Assignment | B | General | Medium | `cwe-701-750.json` |
| `CWE-477` | Use of Obsolete Function | B | General | Medium | `cwe-451-500.json` |

## Executive Summary

Ship **8** remaining CWE heuristics in this theme. Prefer high-signal, low-FP patterns; same-file heuristics only (no interproc taint). Extend registration catalogue tests and fixture matrix.

## Phase 0: Placement + budget

- [x] Record baseline: `wc -l internal/lang/python/detectors/cwe/*.go`
- [x] Ensure target `rules_resource.go` exists (create if needed); keep ≤1500 lines projected
- [x] Prefer domain file over growing `rules.go` past soft cap
- [x] Append FN-safe needles to `needles.go` / per-rule gates
- [x] No changes to `detectors/all.go` unless constructor rename (should be none)

## Phase 1: `CWE-434` — Unrestricted Upload of File with Dangerous Type

> Catalogue: `ruleset/python/chunks/cwe-401-450.json` · tier **B** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE434` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-434", detectCWE434, &Meta…, gates...)` in `rules_resource.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-434`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Dynamic Analysis with Automated Results Interpretation, Dynamic Analysis with Manual Results Interpretation, Manual Static Analysis - Source Code, Automated Static Analysis - Source Code, Archite…
- Suggested sinks/patterns: `request.files save without allowlist`, `secure_filename only`, `.html/.php/.svg upload`, `werkzeug FileStorage.save`
- [x] Implement `detectCWE434` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-434`
- [x] Unit miss: safe pattern → no `CWE-434`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-434-vulnerable.txt` + `CWE-434-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 2: `CWE-427` — Uncontrolled Search Path Element

> Catalogue: `ruleset/python/chunks/cwe-401-450.json` · tier **B** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE427` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-427", detectCWE427, &Meta…, gates...)` in `rules_resource.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-427`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `LD_LIBRARY_PATH set`, `PYTHONPATH mutation`, `dll load from cwd`
- [x] Implement `detectCWE427` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-427`
- [x] Unit miss: safe pattern → no `CWE-427`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-427-vulnerable.txt` + `CWE-427-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 3: `CWE-379` — Creation of Temporary File in Directory with Insecure Permissions

> Catalogue: `ruleset/python/chunks/cwe-351-400.json` · tier **B** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE379` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-379", detectCWE379, &Meta…, gates...)` in `rules_resource.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-379`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `hardcoded /tmp writes`, `tempfile.tempdir = '/tmp'`
- [x] Implement `detectCWE379` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-379`
- [x] Unit miss: safe pattern → no `CWE-379`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-379-vulnerable.txt` + `CWE-379-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 4: `CWE-459` — Incomplete Cleanup

> Catalogue: `ruleset/python/chunks/cwe-451-500.json` · tier **B** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE459` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-459", detectCWE459, &Meta…, gates...)` in `rules_resource.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-459`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `tempfile without delete`, `mkstemp never unlink`
- [x] Implement `detectCWE459` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-459`
- [x] Unit miss: safe pattern → no `CWE-459`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-459-vulnerable.txt` + `CWE-459-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 5: `CWE-772` — Missing Release of Resource after Effective Lifetime

> Catalogue: `ruleset/python/chunks/cwe-751-800.json` · tier **B** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE772` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-772", detectCWE772, &Meta…, gates...)` in `rules_resource.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-772`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `open( without with/`, `socket.socket without close/context`, `urllib urlopen without with`
- [x] Implement `detectCWE772` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-772`
- [x] Unit miss: safe pattern → no `CWE-772`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-772-vulnerable.txt` + `CWE-772-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 6: `CWE-770` — Allocation of Resources Without Limits or Throttling

> Catalogue: `ruleset/python/chunks/cwe-751-800.json` · tier **B** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE770` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-770", detectCWE770, &Meta…, gates...)` in `rules_resource.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-770`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Manual Static Analysis, Fuzzing, Automated Dynamic Analysis, Automated Static Analysis. Manual static analysis can be useful for finding this weakness, but it might not achieve desired code cover…
- Suggested sinks/patterns: `request.get_data(cache=True) no max`, `MAX_CONTENT_LENGTH missing`, `read() unbounded`, `while True: accept/socket`
- [x] Implement `detectCWE770` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-770`
- [x] Unit miss: safe pattern → no `CWE-770`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-770-vulnerable.txt` + `CWE-770-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 7: `CWE-708` — Incorrect Ownership Assignment

> Catalogue: `ruleset/python/chunks/cwe-701-750.json` · tier **B** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE708` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-708", detectCWE708, &Meta…, gates...)` in `rules_resource.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-708`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Analysis. Use automated tools to check for privilege settings. Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates &…
- Suggested sinks/patterns: `os.chown`, `os.fchown`, `pathlib chmod 0o777`, `os.chmod(.*0o777`
- [x] Implement `detectCWE708` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-708`
- [x] Unit miss: safe pattern → no `CWE-708`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-708-vulnerable.txt` + `CWE-708-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 8: `CWE-477` — Use of Obsolete Function

> Catalogue: `ruleset/python/chunks/cwe-451-500.json` · tier **B** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE477` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-477", detectCWE477, &Meta…, gates...)` in `rules_resource.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-477`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis - Binary or Bytecode, Manual Static Analysis - Binary or Bytecode, Dynamic Analysis with Manual Results Interpretation, Manual Static Analysis - Source Code, Automated S…
- Suggested sinks/patterns: `tempfile.mktemp`, `cgi.escape`, `asyncore`, `imp module`, `platform-specific obsolete`
- [x] Implement `detectCWE477` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-477`
- [x] Unit miss: safe pattern → no `CWE-477`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-477-vulnerable.txt` + `CWE-477-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 9: Batch validation

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
