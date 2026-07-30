# Phase 1 — Docs: Go-only → multi-language + Python WIP

> **Parent:** `plans/v0.0.2/python-support.md`  
> **Issue:** #40  
> **Status:** not started  
> **Estimated effort:** small (docs only)  
> **Branch:** `docs/python-wip-refs`  
> **Out of scope:** code/config/ruleset implementation (Phases 2–4)

---

## Overview

Public and internal docs still claim permanent Go-only product scope. Reframe as:

- **Go** = production default, fully shipped  
- **Engine** = multi-language (`LanguagePlugin`)  
- **Python** = reserved / WIP (fixtures + `LanguagePython`; no full plugin yet)

---

## Executive Summary

Accuracy fix only. Do not invent Python detector parity. Link epic #39 where plans said “defer forever”.

---

## 1.1 Product README / overview

- [ ] `README.md` — first paragraphs no longer assert permanent Go-only product; mention multi-language engine + Python WIP if language support is described
- [ ] `documents/overview.md` (if present) — same accuracy pass
- [ ] Keep Go detector catalogue numbers (PERF/CWE/BP) clearly **Go** catalogue counts

## 1.2 Architecture language seam

- [ ] `documents/architecture-performance.md` § Language seam — replace “Go is the only shipped language” with shipped-vs-WIP wording; keep pure-Go / no-CGO requirement for default plugins
- [ ] Cross-link `internal/core/plugin.go` “Adding a second language” guidance

## 1.3 Plans / parity ledger

- [ ] `plans/parity-matrix.md` — `lang/python/` row points to epic #39 / this v0.0.2 plan instead of “out of scope v0” forever
- [ ] `plans/architecture-go.md` — Python deferral points at #39
- [ ] `plans/port-phasewise-checklist.md` — Python reserved bullets note foundation work under v0.0.2

## 1.4 Templates / user-facing strings (docs only)

- [ ] `templates/goslop.toml` comments remain accurate (do **not** invent `languages` key here if Phase 3 owns schema — optional note “see Phase 3” only if needed)
- [ ] Scan `documents/*.md` for “Go only” / “Go-only” claims; fix product claims, leave historical review HTML alone

## 1.5 Closure

- [ ] Filled PR body at `plans/PR/v0.0.2/pr-phase-1-docs.md` (PR_TEMPLATE sections, real content)
- [ ] `gh pr create` with `--assignee @me`, labels `documentation` (+ `enhancement` if mixed), `Closes #40`, `Relates to #39`
- [ ] Parent ledger Phase 1 rollup updated when PR open

---

## Success criteria

- Readers understand Go is default shipped language and Python is WIP under #39  
- No claim of full Python detector support  
- Historical `plans/v0.0.1/reviews/**` left untouched  

## Validation

Docs-only: no `make lint` / `make test` required.

## Dependencies

None (can ship in parallel with Phases 2–4).
