# Batch 08 — Path Fs Temp

> **Parent:** `plans/v0.0.2/heuristics/cwe-plans/README.md` — v0.0.2 Python CWE remaining heuristics
> **Epic / issue:** [#51](https://github.com/chinmay-sawant/goslop/issues/51) · [#52](https://github.com/chinmay-sawant/goslop/issues/52) expansion
> **Status:** **complete** — implemented and validated 2026-07-31
> **Validation evidence:** `gofmt`, Python CWE package tests, Python CWE fixture matrix, `make lint`, `make test`, and `git diff --check` passed.
> **Wave:** P1
> **IDs (8):** CWE-73, CWE-59, CWE-41, CWE-276, CWE-378, CWE-426, CWE-250, CWE-494
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
| Target file(s) | `rules_path_fs.go` |
| Validation | `make lint` + `make test` unchecked until green |
| Fixtures (text) | **Required:** `tests/fixtures/python/cwe/CWE-N-{vulnerable,safe}.txt` (BP-style; see Testing section) |
| Unit tests | **Required:** hit/miss in `internal/lang/python/detectors/cwe/*_test.go` |
| Integration | **Required:** `tests/integration/python/cwe_matrix_test.go` / `make integration-python` |

## Overview

Path residuals beyond CWE-22, permissions, temp files, search path, untrusted search/update

| Rule | Name | Tier | Category | Relevance | Chunk |
|------|------|------|----------|-----------|-------|
| `CWE-73` | External Control of File Name or Path | B | Path Traversal | Medium | `cwe-051-100.json` |
| `CWE-59` | Improper Link Resolution Before File Access ('Link Following') | B | File Processing | Medium | `cwe-051-100.json` |
| `CWE-41` | Improper Resolution of Path Equivalence | B | File Processing | Medium | `cwe-001-050.json` |
| `CWE-276` | Incorrect Default Permissions | A | Security | Medium | `cwe-251-300.json` |
| `CWE-378` | Creation of Temporary File With Insecure Permissions | A | Security | Medium | `cwe-351-400.json` |
| `CWE-426` | Untrusted Search Path | A | Program Invocation | Medium | `cwe-401-450.json` |
| `CWE-250` | Execution with Unnecessary Privileges | A | Security | Medium | `cwe-201-250.json` |
| `CWE-494` | Download of Code Without Integrity Check | A | General | Medium | `cwe-451-500.json` |

## Executive Summary

Ship **8** remaining CWE heuristics in this theme. Prefer high-signal, low-FP patterns; same-file heuristics only (no interproc taint). Extend registration catalogue tests and fixture matrix.

## Phase 0: Placement + budget

- [x] Record baseline: `wc -l internal/lang/python/detectors/cwe/*.go`
- [x] Ensure target `rules_path_fs.go` exists (create if needed); keep ≤1500 lines projected
- [x] Prefer domain file over growing `rules.go` past soft cap
- [x] Append FN-safe needles to `needles.go` / per-rule gates
- [x] No changes to `detectors/all.go` unless constructor rename (should be none)

## Phase 1: `CWE-73` — External Control of File Name or Path

> Catalogue: `ruleset/python/chunks/cwe-051-100.json` · tier **B** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE73` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-73", detectCWE73, &Meta…, gates...)` in `rules_path_fs.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-73`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. The external control or influence of filenames can often be detected using automated static analysis that models data flow within the product. Automated static analysis…
- Suggested sinks/patterns: `open(request.`, `open(f"`, `Path(user)`, `os.path.join(base, request`
- [x] Implement `detectCWE73` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-73`
- [x] Unit miss: safe pattern → no `CWE-73`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-73-vulnerable.txt` + `CWE-73-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 2: `CWE-59` — Improper Link Resolution Before File Access ('Link Following')

> Catalogue: `ruleset/python/chunks/cwe-051-100.json` · tier **B** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE59` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-59", detectCWE59, &Meta…, gates...)` in `rules_path_fs.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-59`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis - Binary or Bytecode, Manual Static Analysis - Binary or Bytecode, Dynamic Analysis with Automated Results Interpretation, Dynamic Analysis with Manual Results Interpret…
- Suggested sinks/patterns: `open(`, `os.remove`, `shutil.move`, `pathlib`, `missing os.lstat/follow_symlinks=False`
- [x] Implement `detectCWE59` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-59`
- [x] Unit miss: safe pattern → no `CWE-59`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-59-vulnerable.txt` + `CWE-59-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 3: `CWE-41` — Improper Resolution of Path Equivalence

> Catalogue: `ruleset/python/chunks/cwe-001-050.json` · tier **B** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE41` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-41", detectCWE41, &Meta…, gates...)` in `rules_path_fs.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-41`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis - Binary or Bytecode, Manual Static Analysis - Binary or Bytecode, Dynamic Analysis with Automated Results Interpretation, Dynamic Analysis with Manual Results Interpret…
- Suggested sinks/patterns: `os.path.join`, `Path /`, `..`, `%2e`, `os.path.normpath without resolve`
- [x] Implement `detectCWE41` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-41`
- [x] Unit miss: safe pattern → no `CWE-41`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-41-vulnerable.txt` + `CWE-41-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 4: `CWE-276` — Incorrect Default Permissions

> Catalogue: `ruleset/python/chunks/cwe-251-300.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE276` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-276", detectCWE276, &Meta…, gates...)` in `rules_path_fs.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-276`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis - Binary or Bytecode, Manual Static Analysis - Binary or Bytecode, Dynamic Analysis with Automated Results Interpretation, Dynamic Analysis with Manual Results Interpret…
- Suggested sinks/patterns: `os.chmod(..., 0o777)`, `os.chmod(..., 0o666)`, `os.umask(0)`, `open mode world-writable`, `Path.chmod(0o777)`
- [x] Implement `detectCWE276` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-276`
- [x] Unit miss: safe pattern → no `CWE-276`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-276-vulnerable.txt` + `CWE-276-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 5: `CWE-378` — Creation of Temporary File With Insecure Permissions

> Catalogue: `ruleset/python/chunks/cwe-351-400.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE378` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-378", detectCWE378, &Meta…, gates...)` in `rules_path_fs.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-378`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `tempfile.mktemp(`, `open('/tmp/…'+name)`, `NamedTemporaryFile(delete=False) + chmod widen`
- [x] Implement `detectCWE378` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-378`
- [x] Unit miss: safe pattern → no `CWE-378`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-378-vulnerable.txt` + `CWE-378-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 6: `CWE-426` — Untrusted Search Path

> Catalogue: `ruleset/python/chunks/cwe-401-450.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE426` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-426", detectCWE426, &Meta…, gates...)` in `rules_path_fs.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-426`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Black Box, Automated Static Analysis, Manual Analysis. Use monitoring tools that examine the software's process as it interacts with the operating system and the network. This technique is useful…
- Suggested sinks/patterns: `sys.path.insert(0, '.')`, `sys.path.insert(0, os.getcwd())`, `sys.path.append(user_input)`, `PATH prepend user-controlled`
- [x] Implement `detectCWE426` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-426`
- [x] Unit miss: safe pattern → no `CWE-426`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-426-vulnerable.txt` + `CWE-426-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 7: `CWE-250` — Execution with Unnecessary Privileges

> Catalogue: `ruleset/python/chunks/cwe-201-250.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE250` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-250", detectCWE250, &Meta…, gates...)` in `rules_path_fs.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-250`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Manual Analysis, Black Box, Automated Static Analysis - Binary or Bytecode, Manual Static Analysis - Binary or Bytecode, Dynamic Analysis with Automated Results Interpretation, Dynamic Analysis w…
- Suggested sinks/patterns: `os.chmod(..., 0o777)`, `os.umask(0)`, `stat.S_IWOTH`, `os.setuid`, `os.setgid`, `0o666 world`
- [x] Implement `detectCWE250` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-250`
- [x] Unit miss: safe pattern → no `CWE-250`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-250-vulnerable.txt` + `CWE-250-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 8: `CWE-494` — Download of Code Without Integrity Check

> Catalogue: `ruleset/python/chunks/cwe-451-500.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE494` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-494", detectCWE494, &Meta…, gates...)` in `rules_path_fs.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-494`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Manual Analysis, Black Box, Automated Static Analysis. This weakness can be detected using tools and techniques that require manual (human) analysis, such as penetration testing, threat modeling,…
- Suggested sinks/patterns: `exec(urlopen(...).read())`, `eval(requests.get(...).text)`, `compile(download)`, `importlib from URL without hash`
- [x] Implement `detectCWE494` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-494`
- [x] Unit miss: safe pattern → no `CWE-494`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-494-vulnerable.txt` + `CWE-494-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
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
