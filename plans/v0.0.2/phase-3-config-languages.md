# Phase 3 — Config: enable / disable languages

> **Parent:** `plans/v0.0.2/python-support.md`  
> **Issue:** #42  
> **Status:** not started  
> **Estimated effort:** medium–large  
> **Branch:** `feat/config-languages`  
> **Out of scope:** Full Python detectors; JSON catalogue content (Phase 4); docs-only product marketing (Phase 1)

---

## Overview

Re-introduce a **real** `languages` configuration that the engine honors. Previously `languages = ["go"]` was rejected as an unknown/inert field (`TestUnsupportedLanguageAndTypedConfigurationRejected`). That rejection must be replaced with validated, wired behavior — never a silent no-op.

### Proposed surface

```toml
[goslop]
# Enabled analysis languages. Default when unset: ["go"].
# Unknown tokens rejected at load time.
languages = ["go"]
# Future opt-in:
# languages = ["go", "python"]
```

CLI (optional in this phase): document as follow-up if not implemented — e.g. `--languages go,python`. Prefer config-first if timeboxed.

---

## Executive Summary

Wire `languages` through schema → parse → merge → registry/scan so only enabled language plugins contribute extensions and detectors. Default remains Go-only. Invalid tokens fail load.

---

## 3.1 Schema and templates

- [ ] `goslop.schema.json` — add `languages` array of strings (enum or free string validated in Go)
- [ ] `templates/goslop.toml` — commented example for `languages`
- [ ] `internal/app/init_template.toml` — same if used by `goslop init`

## 3.2 Config model + validation

- [ ] `internal/config.Section` — `Languages []string \`toml:"languages"\``
- [ ] Validate each token via `core.ParseLanguage` (or Phase 2 helper); reject unknowns
- [ ] Empty list: reject **or** treat as default Go — pick one, test it, document (recommend: empty rejected or falls back to default `["go"]` explicitly)
- [ ] Rewrite `TestUnsupportedLanguageAndTypedConfigurationRejected`: **languages accepted when valid**; still reject `typed.enabled`

## 3.3 Merge

- [ ] `Merged.Languages []core.LanguageID` (or string form with parse at use site)
- [ ] Default when config missing / field unset: `[LanguageGo]`
- [ ] Document CLI override policy if any

## 3.4 Engine / app application

- [ ] Build registry from **enabled** plugins only (filter `DefaultRegistry` plugins or construct filtered registry)
- [ ] Walk / collect files only for enabled extensions
- [ ] `--list-rules` respects enabled languages when config present
- [ ] Tests: with `languages = ["go"]` python files not scanned; with `["python"]` only (if stub plugin registered) go files skipped

## 3.5 Python enable path (WIP)

- [ ] When Python is enabled and a stub/plugin exists, `.py` files enter walk; zero findings OK
- [ ] When Python disabled (default), `.py` files ignored even if present under scan root

## 3.6 Closure

- [ ] `make lint` green — record  
- [ ] `make test` green — record  
- [ ] Filled PR body `plans/PR/v0.0.2/pr-phase-3-config.md`  
- [ ] `Closes #42`, `Relates to #39`  

---

## Success criteria

| Check | Expected |
|-------|----------|
| Default scan | Go only |
| Valid `languages = ["go","python"]` | accepted; engine filters plugins |
| Invalid `languages = ["ruby"]` | load error |
| No silent inert key | schema + merge + scan agree |

## Dependencies

- Prefer Phase 2 helpers merged first on integration; may land in parallel and resolve conflicts on integration branch.
- Phase 4 independent except list-rules catalogue paths.
