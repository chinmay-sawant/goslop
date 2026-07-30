# Batch 14 — Resource Upload

> **Parent:** `plans/v0.0.2/heuristics/cwe-plans/README.md` — v0.0.2 Python CWE remaining heuristics
> **Epic / issue:** [#51](https://github.com/chinmay-sawant/goslop/issues/51) · [#52](https://github.com/chinmay-sawant/goslop/issues/52) expansion
> **Status:** **pending** — not started
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
| Fixtures | `tests/fixtures/python/cwe/CWE-N-{{vulnerable,safe}}.txt` |

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

- [ ] Record baseline: `wc -l internal/lang/python/detectors/cwe/*.go`
- [ ] Ensure target `rules_resource.go` exists (create if needed); keep ≤1500 lines projected
- [ ] Prefer domain file over growing `rules.go` past soft cap
- [ ] Append FN-safe needles to `needles.go` / per-rule gates
- [ ] No changes to `detectors/all.go` unless constructor rename (should be none)

## Phase 1: `CWE-434` — Unrestricted Upload of File with Dangerous Type

> Catalogue: `ruleset/python/chunks/cwe-401-450.json` · tier **B** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE434` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-434", detectCWE434, &Meta…, gates...)` in `rules_resource.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-434`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Dynamic Analysis with Automated Results Interpretation, Dynamic Analysis with Manual Results Interpretation, Manual Static Analysis - Source Code, Automated Static Analysis - Source Code, Archite…
- Suggested sinks/patterns: `request.files save without allowlist`, `secure_filename only`, `.html/.php/.svg upload`, `werkzeug FileStorage.save`
- [ ] Implement `detectCWE434` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-434`
- [ ] Unit miss: safe pattern → no `CWE-434`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-434-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 2: `CWE-427` — Uncontrolled Search Path Element

> Catalogue: `ruleset/python/chunks/cwe-401-450.json` · tier **B** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE427` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-427", detectCWE427, &Meta…, gates...)` in `rules_resource.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-427`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `LD_LIBRARY_PATH set`, `PYTHONPATH mutation`, `dll load from cwd`
- [ ] Implement `detectCWE427` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-427`
- [ ] Unit miss: safe pattern → no `CWE-427`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-427-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 3: `CWE-379` — Creation of Temporary File in Directory with Insecure Permissions

> Catalogue: `ruleset/python/chunks/cwe-351-400.json` · tier **B** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE379` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-379", detectCWE379, &Meta…, gates...)` in `rules_resource.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-379`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `hardcoded /tmp writes`, `tempfile.tempdir = '/tmp'`
- [ ] Implement `detectCWE379` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-379`
- [ ] Unit miss: safe pattern → no `CWE-379`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-379-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 4: `CWE-459` — Incomplete Cleanup

> Catalogue: `ruleset/python/chunks/cwe-451-500.json` · tier **B** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE459` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-459", detectCWE459, &Meta…, gates...)` in `rules_resource.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-459`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `tempfile without delete`, `mkstemp never unlink`
- [ ] Implement `detectCWE459` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-459`
- [ ] Unit miss: safe pattern → no `CWE-459`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-459-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 5: `CWE-772` — Missing Release of Resource after Effective Lifetime

> Catalogue: `ruleset/python/chunks/cwe-751-800.json` · tier **B** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE772` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-772", detectCWE772, &Meta…, gates...)` in `rules_resource.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-772`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `open( without with/`, `socket.socket without close/context`, `urllib urlopen without with`
- [ ] Implement `detectCWE772` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-772`
- [ ] Unit miss: safe pattern → no `CWE-772`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-772-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 6: `CWE-770` — Allocation of Resources Without Limits or Throttling

> Catalogue: `ruleset/python/chunks/cwe-751-800.json` · tier **B** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE770` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-770", detectCWE770, &Meta…, gates...)` in `rules_resource.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-770`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Manual Static Analysis, Fuzzing, Automated Dynamic Analysis, Automated Static Analysis. Manual static analysis can be useful for finding this weakness, but it might not achieve desired code cover…
- Suggested sinks/patterns: `request.get_data(cache=True) no max`, `MAX_CONTENT_LENGTH missing`, `read() unbounded`, `while True: accept/socket`
- [ ] Implement `detectCWE770` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-770`
- [ ] Unit miss: safe pattern → no `CWE-770`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-770-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 7: `CWE-708` — Incorrect Ownership Assignment

> Catalogue: `ruleset/python/chunks/cwe-701-750.json` · tier **B** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE708` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-708", detectCWE708, &Meta…, gates...)` in `rules_resource.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-708`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Analysis. Use automated tools to check for privilege settings. Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates &…
- Suggested sinks/patterns: `os.chown`, `os.fchown`, `pathlib chmod 0o777`, `os.chmod(.*0o777`
- [ ] Implement `detectCWE708` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-708`
- [ ] Unit miss: safe pattern → no `CWE-708`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-708-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 8: `CWE-477` — Use of Obsolete Function

> Catalogue: `ruleset/python/chunks/cwe-451-500.json` · tier **B** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE477` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-477", detectCWE477, &Meta…, gates...)` in `rules_resource.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-477`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis - Binary or Bytecode, Manual Static Analysis - Binary or Bytecode, Dynamic Analysis with Manual Results Interpretation, Manual Static Analysis - Source Code, Automated S…
- Suggested sinks/patterns: `tempfile.mktemp`, `cgi.escape`, `asyncore`, `imp module`, `platform-specific obsolete`
- [ ] Implement `detectCWE477` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-477`
- [ ] Unit miss: safe pattern → no `CWE-477`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-477-vulnerable.txt` + `-safe.txt`
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

