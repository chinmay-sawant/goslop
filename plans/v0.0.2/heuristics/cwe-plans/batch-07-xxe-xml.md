# Batch 07 — Xxe Xml

> **Parent:** `plans/v0.0.2/heuristics/cwe-plans/README.md` — v0.0.2 Python CWE remaining heuristics
> **Epic / issue:** [#51](https://github.com/chinmay-sawant/goslop/issues/51) · [#52](https://github.com/chinmay-sawant/goslop/issues/52) expansion
> **Status:** **pending** — not started
> **Wave:** P1
> **IDs (3):** CWE-611, CWE-776, CWE-112
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
| Target file(s) | `rules_xml.go` |
| Validation | `make lint` + `make test` unchecked until green |
| Fixtures | `tests/fixtures/python/cwe/CWE-N-{{vulnerable,safe}}.txt` |

## Overview

XXE and unsafe XML entity expansion / missing validation

| Rule | Name | Tier | Category | Relevance | Chunk |
|------|------|------|----------|-----------|-------|
| `CWE-611` | Improper Restriction of XML External Entity Reference | A | XML | Medium | `cwe-601-650.json` |
| `CWE-776` | Improper Restriction of Recursive Entity References in DTDs ('XML Entity Expansion') | A | XML | Medium | `cwe-751-800.json` |
| `CWE-112` | Missing XML Validation | A | XML | Medium | `cwe-101-150.json` |

## Executive Summary

Ship **3** remaining CWE heuristics in this theme. Prefer high-signal, low-FP patterns; same-file heuristics only (no interproc taint). Extend registration catalogue tests and fixture matrix.

## Phase 0: Placement + budget

- [ ] Record baseline: `wc -l internal/lang/python/detectors/cwe/*.go`
- [ ] Ensure target `rules_xml.go` exists (create if needed); keep ≤1500 lines projected
- [ ] Prefer domain file over growing `rules.go` past soft cap
- [ ] Append FN-safe needles to `needles.go` / per-rule gates
- [ ] No changes to `detectors/all.go` unless constructor rename (should be none)

## Phase 1: `CWE-611` — Improper Restriction of XML External Entity Reference

> Catalogue: `ruleset/python/chunks/cwe-601-650.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE611` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-611", detectCWE611, &Meta…, gates...)` in `rules_xml.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-611`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `xml.etree.ElementTree.parse/fromstring`, `lxml.etree.parse/fromstring/XML`, `xml.dom.minidom.parse`, `xml.sax.parse`, `pulldom`, `BeautifulSoup with lxml xml`, `resolve_entities=True`
- [ ] Implement `detectCWE611` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-611`
- [ ] Unit miss: safe pattern → no `CWE-611`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-611-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 2: `CWE-776` — Improper Restriction of Recursive Entity References in DTDs ('XML Entity Expansion')

> Catalogue: `ruleset/python/chunks/cwe-751-800.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE776` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-776", detectCWE776, &Meta…, gates...)` in `rules_xml.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-776`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `same as 611: etree/lxml/minidom without defusedxml`, `DTD enabled parsers`
- [ ] Implement `detectCWE776` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-776`
- [ ] Unit miss: safe pattern → no `CWE-776`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-776-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 3: `CWE-112` — Missing XML Validation

> Catalogue: `ruleset/python/chunks/cwe-101-150.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE112` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-112", detectCWE112, &Meta…, gates...)` in `rules_xml.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-112`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `xml.etree.ElementTree.parse`, `lxml.etree.parse`, `XMLParser(resolve_entities=True)`, `no defusedxml`
- [ ] Implement `detectCWE112` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-112`
- [ ] Unit miss: safe pattern → no `CWE-112`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-112-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 4: Batch validation

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

