# Batch 01 — Injection Expand

> **Parent:** `plans/v0.0.2/heuristics/cwe-plans/README.md` — v0.0.2 Python CWE remaining heuristics
> **Epic / issue:** [#51](https://github.com/chinmay-sawant/goslop/issues/51) · [#52](https://github.com/chinmay-sawant/goslop/issues/52) expansion
> **Status:** **pending** — not started
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
| Fixtures | `tests/fixtures/python/cwe/CWE-N-{{vulnerable,safe}}.txt` |

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

- [ ] Record baseline: `wc -l internal/lang/python/detectors/cwe/*.go`
- [ ] Ensure target `rules_injection.go` exists (create if needed); keep ≤1500 lines projected
- [ ] Prefer domain file over growing `rules.go` past soft cap
- [ ] Append FN-safe needles to `needles.go` / per-rule gates
- [ ] No changes to `detectors/all.go` unless constructor rename (should be none)

## Phase 1: `CWE-90` — Improper Neutralization of Special Elements used in an LDAP Query ('LDAP Injection')

> Catalogue: `ruleset/python/chunks/cwe-051-100.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE90` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-90", detectCWE90, &Meta…, gates...)` in `rules_injection.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-90`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `ldap3`, `ldap.initialize`, `.search(`, `f"(uid={`, `filter % user`, `ldap.filter`
- [ ] Implement `detectCWE90` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-90`
- [ ] Unit miss: safe pattern → no `CWE-90`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-90-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 2: `CWE-91` — XML Injection (aka Blind XPath Injection)

> Catalogue: `ruleset/python/chunks/cwe-051-100.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE91` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-91", detectCWE91, &Meta…, gates...)` in `rules_injection.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-91`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `.xpath(f"`, `etree.XPath(`, `lxml`, `xml.etree`, `format into XML string + parse`
- [ ] Implement `detectCWE91` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-91`
- [ ] Unit miss: safe pattern → no `CWE-91`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-91-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 3: `CWE-93` — Improper Neutralization of CRLF Sequences ('CRLF Injection')

> Catalogue: `ruleset/python/chunks/cwe-051-100.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE93` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-93", detectCWE93, &Meta…, gates...)` in `rules_injection.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-93`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `response.headers[`, `set_header`, `HttpResponse`, `Location:`, `Set-Cookie`, `f-string header value`
- [ ] Implement `detectCWE93` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-93`
- [ ] Unit miss: safe pattern → no `CWE-93`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-93-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 4: `CWE-94` — Improper Control of Generation of Code ('Code Injection')

> Catalogue: `ruleset/python/chunks/cwe-051-100.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE94` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-94", detectCWE94, &Meta…, gates...)` in `rules_injection.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-94`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `eval(`, `exec(`, `compile(`, `__import__(`, `importlib.import_module(dynamic)`
- [ ] Implement `detectCWE94` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-94`
- [ ] Unit miss: safe pattern → no `CWE-94`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-94-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 5: `CWE-88` — Improper Neutralization of Argument Delimiters in a Command ('Argument Injection')

> Catalogue: `ruleset/python/chunks/cwe-051-100.json` · tier **B** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE88` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-88", detectCWE88, &Meta…, gates...)` in `rules_injection.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-88`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `subprocess.run([.., user])`, `subprocess.Popen list argv`, `shell=False with dynamic args`
- [ ] Implement `detectCWE88` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-88`
- [ ] Unit miss: safe pattern → no `CWE-88`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-88-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 6: `CWE-117` — Improper Output Neutralization for Logs

> Catalogue: `ruleset/python/chunks/cwe-101-150.json` · tier **B** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE117` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-117", detectCWE117, &Meta…, gates...)` in `rules_injection.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-117`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `logging.`, `logger.`, `log.info(f"...{request`, `print to log without \n strip`
- [ ] Implement `detectCWE117` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-117`
- [ ] Unit miss: safe pattern → no `CWE-117`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-117-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 7: Batch validation

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

