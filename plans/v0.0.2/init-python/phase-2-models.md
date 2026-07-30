# Phase 2 — Models: multi-language domain types

> **Parent:** `plans/v0.0.2/python-support.md`  
> **Issue:** #41  
> **Status:** complete (implementation on `feat/python-models`)  
> **Estimated effort:** medium  
> **Branch:** `feat/python-models`  
> **Out of scope:** TOML schema / `languages` config (Phase 3); JSON ruleset files (Phase 4); full Python detectors

---

## Overview

Make language a first-class, consistent domain concept so Phase 3 can filter plugins and Phase 4 can attach catalogues without hacks.

Existing anchors:

- `internal/core/language.go` — `LanguageGo`, `LanguagePython`, `ParseLanguage`, `LanguageFromExtension`, `ParseLanguages`, `DefaultEnabledLanguages`
- `internal/core/plugin.go` — `LanguagePlugin` contract; `BasePlugin.ParseSource` → Go
- `internal/fixture/format.go` — fixture `lang: go | python`
- `BasePlugin.ParseSource` defaults language to **Go** (plugins must override)
- `internal/lang/python/` — optional source-only stub (not in `DefaultRegistry`)

---

## Executive Summary

Harden models and small helpers for multi-language readiness. Prefer minimal surface: helpers for parsing language lists, default language set, and finding/unit language consistency. Avoid speculative catalogues until a second plugin is real.

---

## 2.1 Language identity helpers

- [x] Audit `LanguageID` / `ParseLanguage` / `LanguageFromExtension` for python aliases (`python`, `py`) — keep behavior, add tests if gaps — evidence: `internal/core/language.go` + `language_test.go` (`TestParseLanguage`, `TestLanguageFromExtension`; also accepts leading-dot ext)
- [x] Add helper to parse a list of language names → `[]LanguageID` with stable order and clear errors (for Phase 3 reuse) — evidence: `core.ParseLanguages`
- [x] Document default enabled set: Go only (`DefaultEnabledLanguages` or equivalent) — evidence: `core.DefaultEnabledLanguages()` → `[LanguageGo]`

## 2.2 ParsedUnit / detector / session consistency

- [x] Ensure `ParsedUnit.Language` is always set correctly by plugins (document BasePlugin Go default) — evidence: `BasePlugin` package docs; python stub overrides `ParseSource` with `LanguagePython`
- [x] Finding / fingerprint paths: confirm language is available where cache keys or reporting need it (extend only if a concrete gap exists) — audit: `Finding` / fingerprint / cache entry are path+rule keyed; walk uses `LanguageID` via registry extension map; no concrete gap → no schema change
- [x] Registry: document that detectors are indexed by `LanguageID` (`internal/engine/registry.go`) — evidence: `Registry` struct doc + `TestRegistryIndexesDetectorsByLanguageID`

## 2.3 Optional stub plugin seam (minimal)

- [x] If useful for tests: empty or source-only Python plugin **stub** under `internal/lang/python/` implementing `LanguagePlugin` with `Extensions() == ["py"]`, **zero detectors**, pure source-only parse — evidence: `internal/lang/python/plugin.go`
- [x] Do **not** register Python in `DefaultRegistry` until Phase 3 language filter exists (or register behind tests only) — evidence: `DefaultRegistry` comment + `TestDefaultRegistryGoOnly` / python package tests

## 2.4 Tests

- [x] Unit tests for language list parsing / defaults — evidence: `internal/core/language_test.go`
- [x] Fixture materialize python path still works (`internal/fixture`) — evidence: `TestMaterializeTree_IncludesPythonTagged` still green
- [x] No regression: `DefaultRegistry` still Go-only for production default — evidence: `TestDefaultRegistryGoOnly`

## 2.5 Closure

- [x] `make lint` green — record outcome — 2026-07-30: `go vet ./...` + `gofmt -l` clean
- [x] `make test` green — record outcome — 2026-07-30: `go test ./...` all packages ok (incl. `internal/lang/python`, integration)
- [x] Filled PR body `plans/PR/v0.0.2/pr-phase-2-models.md`
- [x] `Closes #41`, `Relates to #39`
- [x] `CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop` — green

---

## Success criteria

- Domain types represent Go + Python without hacks  
- Helpers ready for config merge  
- Go remains production default in registry  

## Dependencies

Unblocks Phase 3. Parallel with Phase 1 and Phase 4 (careful if both touch `DefaultRegistry`).
