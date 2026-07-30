# Phase 2 — Models: multi-language domain types

> **Parent:** `plans/v0.0.2/python-support.md`  
> **Issue:** #41  
> **Status:** not started  
> **Estimated effort:** medium  
> **Branch:** `feat/python-models`  
> **Out of scope:** TOML schema / `languages` config (Phase 3); JSON ruleset files (Phase 4); full Python detectors

---

## Overview

Make language a first-class, consistent domain concept so Phase 3 can filter plugins and Phase 4 can attach catalogues without hacks.

Existing anchors:

- `internal/core/language.go` — `LanguageGo`, `LanguagePython`, `ParseLanguage`, `LanguageFromExtension`
- `internal/core/plugin.go` — `LanguagePlugin` contract
- `internal/fixture/format.go` — fixture `lang: go | python`
- `BasePlugin.ParseSource` defaults language to **Go** (plugins must override)

---

## Executive Summary

Harden models and small helpers for multi-language readiness. Prefer minimal surface: helpers for parsing language lists, default language set, and finding/unit language consistency. Avoid speculative catalogues until a second plugin is real.

---

## 2.1 Language identity helpers

- [ ] Audit `LanguageID` / `ParseLanguage` / `LanguageFromExtension` for python aliases (`python`, `py`) — keep behavior, add tests if gaps
- [ ] Add helper to parse a list of language names → `[]LanguageID` with stable order and clear errors (for Phase 3 reuse)
- [ ] Document default enabled set: Go only (`DefaultEnabledLanguages` or equivalent)

## 2.2 ParsedUnit / detector / session consistency

- [ ] Ensure `ParsedUnit.Language` is always set correctly by plugins (document BasePlugin Go default)
- [ ] Finding / fingerprint paths: confirm language is available where cache keys or reporting need it (extend only if a concrete gap exists)
- [ ] Registry: document that detectors are indexed by `LanguageID` (`internal/engine/registry.go`)

## 2.3 Optional stub plugin seam (minimal)

- [ ] If useful for tests: empty or source-only Python plugin **stub** under `internal/lang/python/` implementing `LanguagePlugin` with `Extensions() == ["py"]`, **zero detectors**, pure source-only parse — **only if** it does not force Phase 3/4 scope creep
- [ ] Do **not** register Python in `DefaultRegistry` until Phase 3 language filter exists (or register behind tests only)

## 2.4 Tests

- [ ] Unit tests for language list parsing / defaults
- [ ] Fixture materialize python path still works (`internal/fixture`)
- [ ] No regression: `DefaultRegistry` still Go-only for production default

## 2.5 Closure

- [ ] `make lint` green — record outcome  
- [ ] `make test` green — record outcome  
- [ ] Filled PR body `plans/PR/v0.0.2/pr-phase-2-models.md`  
- [ ] `Closes #41`, `Relates to #39`  

---

## Success criteria

- Domain types represent Go + Python without hacks  
- Helpers ready for config merge  
- Go remains production default in registry  

## Dependencies

Unblocks Phase 3. Parallel with Phase 1 and Phase 4 (careful if both touch `DefaultRegistry`).
