# Batch 01 — Injection Expand

> **Parent:** `plans/v0.0.2/heuristics/cwe-plans/README.md` — v0.0.2 Python CWE remaining heuristics
> **Epic / issue:** [#51](https://github.com/chinmay-sawant/goslop/issues/51) · [#52](https://github.com/chinmay-sawant/goslop/issues/52) expansion
> **Status:** **complete** — implemented and validated 2026-07-31
> **Validation evidence:** `gofmt`, Python CWE package tests, Python CWE fixture matrix, `make lint`, `make test`, and `git diff --check` passed.
> **Wave:** P0
> **IDs (6):** CWE-90, CWE-91, CWE-93, CWE-94, CWE-88, CWE-117
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
| Target file(s) | `rules_injection.go` |
| Validation | `make lint` + `make test` unchecked until green |
| Fixtures (text) | **Required:** `tests/fixtures/python/cwe/CWE-N-{vulnerable,safe}.txt` (BP-style; see Testing section) |
| Unit tests | **Required:** hit/miss in `internal/lang/python/detectors/cwe/*_test.go` |
| Integration | **Required:** `tests/integration/python/cwe_matrix_test.go` / `make integration-python` |

## Overview

Query/command-adjacent injection beyond 78/89: LDAP, XPath, CRLF, code inj, log/arg injection

| Rule | Name | Tier | Category | Relevance | Chunk |
|------|------|------|----------|-----------|-------|
| `CWE-90` | Improper Neutralization of Special Elements used in an LDAP Query ('LDAP Injection') | A | Injection | Medium | `cwe-051-100.json` |
| `CWE-91` | XML Injection (aka Blind XPath Injection) | A | Path Traversal | Medium | `cwe-051-100.json` |
| `CWE-93` | Improper Neutralization of CRLF Sequences ('CRLF Injection') | A | Injection | Medium | `cwe-051-100.json` |
| `CWE-94` | Improper Control of Generation of Code ('Code Injection') | A | Injection | Medium | `cwe-051-100.json` |
| `CWE-88` | Improper Neutralization of Argument Delimiters in a Command ('Argument Injection') | B | Program Invocation | Medium | `cwe-051-100.json` |
| `CWE-117` | Improper Output Neutralization for Logs | B | Information Disclosure | Medium | `cwe-101-150.json` |

## Executive Summary

Ship **6** remaining CWE heuristics in this theme. Prefer high-signal, low-FP patterns; same-file heuristics only (no interproc taint). Extend registration catalogue tests and fixture matrix.

## Phase 0: Placement + budget

- [x] Record baseline: `wc -l internal/lang/python/detectors/cwe/*.go`
- [x] Ensure target `rules_injection.go` exists (create if needed); keep ≤1500 lines projected
- [x] Prefer domain file over growing `rules.go` past soft cap
- [x] Append FN-safe needles to `needles.go` / per-rule gates
- [x] No changes to `detectors/all.go` unless constructor rename (should be none)

## Phase 1: `CWE-90` — Improper Neutralization of Special Elements used in an LDAP Query ('LDAP Injection')

> Catalogue: `ruleset/python/chunks/cwe-051-100.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE90` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-90", detectCWE90, &Meta…, gates...)` in `rules_injection.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-90`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `ldap3`, `ldap.initialize`, `.search(`, `f"(uid={`, `filter % user`, `ldap.filter`
- [x] Implement `detectCWE90` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-90`
- [x] Unit miss: safe pattern → no `CWE-90`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-90-vulnerable.txt` + `CWE-90-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 2: `CWE-91` — XML Injection (aka Blind XPath Injection)

> Catalogue: `ruleset/python/chunks/cwe-051-100.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE91` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-91", detectCWE91, &Meta…, gates...)` in `rules_injection.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-91`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `.xpath(f"`, `etree.XPath(`, `lxml`, `xml.etree`, `format into XML string + parse`
- [x] Implement `detectCWE91` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-91`
- [x] Unit miss: safe pattern → no `CWE-91`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-91-vulnerable.txt` + `CWE-91-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 3: `CWE-93` — Improper Neutralization of CRLF Sequences ('CRLF Injection')

> Catalogue: `ruleset/python/chunks/cwe-051-100.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE93` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-93", detectCWE93, &Meta…, gates...)` in `rules_injection.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-93`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `response.headers[`, `set_header`, `HttpResponse`, `Location:`, `Set-Cookie`, `f-string header value`
- [x] Implement `detectCWE93` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-93`
- [x] Unit miss: safe pattern → no `CWE-93`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-93-vulnerable.txt` + `CWE-93-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 4: `CWE-94` — Improper Control of Generation of Code ('Code Injection')

> Catalogue: `ruleset/python/chunks/cwe-051-100.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE94` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-94", detectCWE94, &Meta…, gates...)` in `rules_injection.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-94`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `eval(`, `exec(`, `compile(`, `__import__(`, `importlib.import_module(dynamic)`
- [x] Implement `detectCWE94` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-94`
- [x] Unit miss: safe pattern → no `CWE-94`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-94-vulnerable.txt` + `CWE-94-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 5: `CWE-88` — Improper Neutralization of Argument Delimiters in a Command ('Argument Injection')

> Catalogue: `ruleset/python/chunks/cwe-051-100.json` · tier **B** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE88` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-88", detectCWE88, &Meta…, gates...)` in `rules_injection.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-88`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `subprocess.run([.., user])`, `subprocess.Popen list argv`, `shell=False with dynamic args`
- [x] Implement `detectCWE88` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-88`
- [x] Unit miss: safe pattern → no `CWE-88`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-88-vulnerable.txt` + `CWE-88-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 6: `CWE-117` — Improper Output Neutralization for Logs

> Catalogue: `ruleset/python/chunks/cwe-101-150.json` · tier **B** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE117` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-117", detectCWE117, &Meta…, gates...)` in `rules_injection.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-117`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `logging.`, `logger.`, `log.info(f"...{request`, `print to log without \n strip`
- [x] Implement `detectCWE117` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-117`
- [x] Unit miss: safe pattern → no `CWE-117`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-117-vulnerable.txt` + `CWE-117-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 7: Batch validation

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
