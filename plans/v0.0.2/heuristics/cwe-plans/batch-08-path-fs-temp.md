# Batch 08 — Path Fs Temp

> **Parent:** `plans/v0.0.2/heuristics/cwe-plans/README.md` — v0.0.2 Python CWE remaining heuristics
> **Epic / issue:** [#51](https://github.com/chinmay-sawant/goslop/issues/51) · [#52](https://github.com/chinmay-sawant/goslop/issues/52) expansion
> **Status:** **pending** — not started
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
| Fixtures | `tests/fixtures/python/cwe/CWE-N-{{vulnerable,safe}}.txt` |

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

- [ ] Record baseline: `wc -l internal/lang/python/detectors/cwe/*.go`
- [ ] Ensure target `rules_path_fs.go` exists (create if needed); keep ≤1500 lines projected
- [ ] Prefer domain file over growing `rules.go` past soft cap
- [ ] Append FN-safe needles to `needles.go` / per-rule gates
- [ ] No changes to `detectors/all.go` unless constructor rename (should be none)

## Phase 1: `CWE-73` — External Control of File Name or Path

> Catalogue: `ruleset/python/chunks/cwe-051-100.json` · tier **B** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE73` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-73", detectCWE73, &Meta…, gates...)` in `rules_path_fs.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-73`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. The external control or influence of filenames can often be detected using automated static analysis that models data flow within the product. Automated static analysis…
- Suggested sinks/patterns: `open(request.`, `open(f"`, `Path(user)`, `os.path.join(base, request`
- [ ] Implement `detectCWE73` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-73`
- [ ] Unit miss: safe pattern → no `CWE-73`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-73-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 2: `CWE-59` — Improper Link Resolution Before File Access ('Link Following')

> Catalogue: `ruleset/python/chunks/cwe-051-100.json` · tier **B** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE59` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-59", detectCWE59, &Meta…, gates...)` in `rules_path_fs.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-59`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis - Binary or Bytecode, Manual Static Analysis - Binary or Bytecode, Dynamic Analysis with Automated Results Interpretation, Dynamic Analysis with Manual Results Interpret…
- Suggested sinks/patterns: `open(`, `os.remove`, `shutil.move`, `pathlib`, `missing os.lstat/follow_symlinks=False`
- [ ] Implement `detectCWE59` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-59`
- [ ] Unit miss: safe pattern → no `CWE-59`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-59-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 3: `CWE-41` — Improper Resolution of Path Equivalence

> Catalogue: `ruleset/python/chunks/cwe-001-050.json` · tier **B** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE41` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-41", detectCWE41, &Meta…, gates...)` in `rules_path_fs.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-41`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis - Binary or Bytecode, Manual Static Analysis - Binary or Bytecode, Dynamic Analysis with Automated Results Interpretation, Dynamic Analysis with Manual Results Interpret…
- Suggested sinks/patterns: `os.path.join`, `Path /`, `..`, `%2e`, `os.path.normpath without resolve`
- [ ] Implement `detectCWE41` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-41`
- [ ] Unit miss: safe pattern → no `CWE-41`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-41-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 4: `CWE-276` — Incorrect Default Permissions

> Catalogue: `ruleset/python/chunks/cwe-251-300.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE276` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-276", detectCWE276, &Meta…, gates...)` in `rules_path_fs.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-276`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis - Binary or Bytecode, Manual Static Analysis - Binary or Bytecode, Dynamic Analysis with Automated Results Interpretation, Dynamic Analysis with Manual Results Interpret…
- Suggested sinks/patterns: `os.chmod(..., 0o777)`, `os.chmod(..., 0o666)`, `os.umask(0)`, `open mode world-writable`, `Path.chmod(0o777)`
- [ ] Implement `detectCWE276` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-276`
- [ ] Unit miss: safe pattern → no `CWE-276`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-276-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 5: `CWE-378` — Creation of Temporary File With Insecure Permissions

> Catalogue: `ruleset/python/chunks/cwe-351-400.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE378` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-378", detectCWE378, &Meta…, gates...)` in `rules_path_fs.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-378`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `tempfile.mktemp(`, `open('/tmp/…'+name)`, `NamedTemporaryFile(delete=False) + chmod widen`
- [ ] Implement `detectCWE378` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-378`
- [ ] Unit miss: safe pattern → no `CWE-378`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-378-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 6: `CWE-426` — Untrusted Search Path

> Catalogue: `ruleset/python/chunks/cwe-401-450.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE426` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-426", detectCWE426, &Meta…, gates...)` in `rules_path_fs.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-426`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Black Box, Automated Static Analysis, Manual Analysis. Use monitoring tools that examine the software's process as it interacts with the operating system and the network. This technique is useful…
- Suggested sinks/patterns: `sys.path.insert(0, '.')`, `sys.path.insert(0, os.getcwd())`, `sys.path.append(user_input)`, `PATH prepend user-controlled`
- [ ] Implement `detectCWE426` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-426`
- [ ] Unit miss: safe pattern → no `CWE-426`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-426-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 7: `CWE-250` — Execution with Unnecessary Privileges

> Catalogue: `ruleset/python/chunks/cwe-201-250.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE250` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-250", detectCWE250, &Meta…, gates...)` in `rules_path_fs.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-250`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Manual Analysis, Black Box, Automated Static Analysis - Binary or Bytecode, Manual Static Analysis - Binary or Bytecode, Dynamic Analysis with Automated Results Interpretation, Dynamic Analysis w…
- Suggested sinks/patterns: `os.chmod(..., 0o777)`, `os.umask(0)`, `stat.S_IWOTH`, `os.setuid`, `os.setgid`, `0o666 world`
- [ ] Implement `detectCWE250` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-250`
- [ ] Unit miss: safe pattern → no `CWE-250`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-250-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 8: `CWE-494` — Download of Code Without Integrity Check

> Catalogue: `ruleset/python/chunks/cwe-451-500.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE494` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-494", detectCWE494, &Meta…, gates...)` in `rules_path_fs.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-494`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Manual Analysis, Black Box, Automated Static Analysis. This weakness can be detected using tools and techniques that require manual (human) analysis, such as penetration testing, threat modeling,…
- Suggested sinks/patterns: `exec(urlopen(...).read())`, `eval(requests.get(...).text)`, `compile(download)`, `importlib from URL without hash`
- [ ] Implement `detectCWE494` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-494`
- [ ] Unit miss: safe pattern → no `CWE-494`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-494-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 9: Batch validation

- [ ] `gofmt -w` on touched files
- [ ] `make lint`
- [ ] `make test`
- [ ] `go test ./tests/integration/python/ -count=1` (CWE matrix)
- [ ] Update `_inventory.json`: move batch IDs from `missing` → `implemented`
- [ ] Update this ledger statuses to `[x]` with evidence
- [ ] Package files still ≤2000 lines (split if not)

## Dependencies

| Depends on | Note |
|------------|------|
| batch-00 | Framework + priority rules already present |
| Catalogue chunks | IDs must exist in `ruleset/python/chunks/` |
| Parent README | ownership + PR policy |

