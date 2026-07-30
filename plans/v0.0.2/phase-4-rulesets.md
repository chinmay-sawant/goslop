# Phase 4 — Rulesets: Python JSON catalogues + reuse

> **Parent:** `plans/v0.0.2/python-support.md`  
> **Issue:** #43  
> **Status:** not started  
> **Estimated effort:** medium  
> **Branch:** `feat/python-rulesets`  
> **Out of scope:** Implementing Python AST detectors; changing Go detector registration counts; Phase 3 config wiring (consume audit only)

---

## Overview

Today all catalogue JSON lives under `ruleset/golang/`:

- `bad-practices.json` — BP metadata (135 rules), Go-idiomatic
- `chunks/*.json` — CWE + PERF metadata with `applicable_to` **entirely Go ecosystem tags** (`golang`, `gin`, `gorm`, …) — **0 rules** list `python`

Phase 4 establishes:

1. **Reuse audit** → [ruleset-reuse-audit.md](./ruleset-reuse-audit.md)  
2. **Layout** for Python (and optional shared) catalogues  
3. **Stub / seed JSON** that reuses schema fields and portable CWE *identity* where sensible  
4. Loader notes so a future plugin can embed/load without forking formats  

---

## Executive Summary

Do not duplicate 400+ Go-tagged rules into Python blindly. Classify **portable weakness IDs** (CWE-22, CWE-78, CWE-79, CWE-89, …) vs **Go-only mechanics** (most PERF, most BP). Seed a minimal Python catalogue + README.

---

## 4.1 Ruleset reuse audit (scan agent deliverable)

- [ ] Inventory schema fields for CWE/PERF chunks and BP JSON  
- [ ] Count rules by `applicable_to` / `go_relevance` / `status`  
- [ ] Classify into: **portable CWE shell**, **Go-framework-specific**, **Go-PERF**, **Go-BP**, **shared metadata only**  
- [ ] Write findings to [ruleset-reuse-audit.md](./ruleset-reuse-audit.md) with concrete rule ID examples  
- [ ] Recommend layout: e.g. `ruleset/python/` + optional `ruleset/shared/` vs overlay files  

## 4.2 Directory + seed catalogues

- [ ] Create `ruleset/python/README.md` explaining WIP + reuse policy  
- [ ] Seed minimal Python catalogues, reusing field names from golang where possible:
  - [ ] `ruleset/python/bad-practices.json` — empty `{}` or small portable subset with clear python notes  
  - [ ] `ruleset/python/chunks/` — at least one CWE seed file for high-value portable IDs (e.g. CWE-22/78/79/89) with `applicable_to` including `python` and python-oriented `detection_notes`  
- [ ] Do **not** break `ruleset/golang/**` or generated Go metadata  

## 4.3 Loader / discovery (minimal code if needed)

- [ ] Document how Go loads BP metadata today (`metadata_gen.go` from bad-practices.json)  
- [ ] If a shared path helper is needed, keep it tiny; prefer docs + file layout over big refactors  
- [ ] Optional: small test that Python JSON files parse as valid JSON objects  

## 4.4 Docs cross-links

- [ ] `documents/rule-catalog-and-maturity.md` — note multi-language catalogue roots if it claims golang-only paths  
- [ ] Point parity matrix python fixtures row at this phase / #43  

## 4.5 Closure

- [ ] Audit markdown committed  
- [ ] Seed rulesets committed  
- [ ] Validation: if only JSON/docs → no full suite required; if Go code touched → `make lint` + `make test`  
- [ ] Filled PR body `plans/PR/v0.0.2/pr-phase-4-rulesets.md`  
- [ ] `Closes #43`, `Relates to #39`  

---

## Success criteria

- Reuse policy is explicit and evidence-based  
- Python ruleset tree exists and loads as JSON  
- Go catalogues unchanged and still default  
- No fake claim of detector implementation  

## Dependencies

Independent of Phases 1–3; integration resolves doc overlaps.
