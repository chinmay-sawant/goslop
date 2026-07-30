# Batch 03 — Mass Assign Deser

> **Parent:** `plans/v0.0.2/heuristics/cwe-plans/README.md` — v0.0.2 Python CWE remaining heuristics
> **Epic / issue:** [#51](https://github.com/chinmay-sawant/goslop/issues/51) · [#52](https://github.com/chinmay-sawant/goslop/issues/52) expansion
> **Status:** **pending** — not started
> **Wave:** P0
> **IDs (3):** CWE-915, CWE-914, CWE-916
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
| Target file(s) | `rules_mass_assign.go` |
| Validation | `make lint` + `make test` unchecked until green |
| Fixtures (text) | **Required:** `tests/fixtures/python/cwe/CWE-N-{vulnerable,safe}.txt` (BP-style; see Testing section) |
| Unit tests | **Required:** hit/miss in `internal/lang/python/detectors/cwe/*_test.go` |
| Integration | **Required:** `tests/integration/python/cwe_matrix_test.go` / `make integration-python` |

## Overview

CWE-915 High mass assignment + dynamic attrs / related deser-adjacent

| Rule | Name | Tier | Category | Relevance | Chunk |
|------|------|------|----------|-----------|-------|
| `CWE-915` | Improperly Controlled Modification of Dynamically-Determined Object Attributes | A | General | High | `cwe-901-950.json` |
| `CWE-914` | Improper Control of Dynamically-Identified Variables | A | General | Medium | `cwe-901-950.json` |
| `CWE-916` | Use of Password Hash With Insufficient Computational Effort | A | Security | Medium | `cwe-901-950.json` |

## Executive Summary

Ship **3** remaining CWE heuristics in this theme. Prefer high-signal, low-FP patterns; same-file heuristics only (no interproc taint). Extend registration catalogue tests and fixture matrix.

## Phase 0: Placement + budget

- [ ] Record baseline: `wc -l internal/lang/python/detectors/cwe/*.go`
- [ ] Ensure target `rules_mass_assign.go` exists (create if needed); keep ≤1500 lines projected
- [ ] Prefer domain file over growing `rules.go` past soft cap
- [ ] Append FN-safe needles to `needles.go` / per-rule gates
- [ ] No changes to `detectors/all.go` unless constructor rename (should be none)

## Phase 1: `CWE-915` — Improperly Controlled Modification of Dynamically-Determined Object Attributes

> Catalogue: `ruleset/python/chunks/cwe-901-950.json` · tier **A** · relevance **High**

### Register + meta

- [ ] Add `MetaCWE915` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-915", detectCWE915, &Meta…, gates...)` in `rules_mass_assign.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-915`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `Model(**request.POST/data/json)`, `obj.__dict__.update(request.*)`, `setattr loop over untrusted keys`, `django update(**request.data)`, `pydantic/extra mass bind without exclude`
- [ ] Implement `detectCWE915` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-915`
- [ ] Unit miss: safe pattern → no `CWE-915`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-915-vulnerable.txt` + `CWE-915-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [ ] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [ ] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 2: `CWE-914` — Improper Control of Dynamically-Identified Variables

> Catalogue: `ruleset/python/chunks/cwe-901-950.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE914` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-914", detectCWE914, &Meta…, gates...)` in `rules_mass_assign.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-914`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `globals()[user]`, `locals()[user]`, `vars(obj)[k]=`, `eval/exec building names`, `setattr(obj, user_key, v) without allowlist`
- [ ] Implement `detectCWE914` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-914`
- [ ] Unit miss: safe pattern → no `CWE-914`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-914-vulnerable.txt` + `CWE-914-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [ ] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [ ] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 3: `CWE-916` — Use of Password Hash With Insufficient Computational Effort

> Catalogue: `ruleset/python/chunks/cwe-901-950.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE916` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-916", detectCWE916, &Meta…, gates...)` in `rules_mass_assign.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-916`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis - Binary or Bytecode, Manual Static Analysis - Binary or Bytecode, Manual Static Analysis - Source Code, Automated Static Analysis - Source Code, Automated Static Analys…
- Suggested sinks/patterns: `hashlib.md5/sha1(.*password`, `hashlib.md5(password.encode)`, `crypt.crypt weak`, `passlib with md5_crypt`
- [ ] Implement `detectCWE916` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [ ] **Unit hit:** vulnerable snippet → finding `CWE-916`
- [ ] Unit miss: safe pattern → no `CWE-916`
- [ ] **Fixtures:** `tests/fixtures/python/cwe/CWE-916-vulnerable.txt` + `CWE-916-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [ ] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [ ] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 4: Batch validation

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

