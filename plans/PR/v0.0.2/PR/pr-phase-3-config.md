## Summary

Reintroduce a **real, wired** `languages` configuration key so analysis language selection is validated and honored end-to-end (schema → parse → merge → registry filter → walk / `--list-rules`). Previously `languages = ["go"]` was rejected as unsupported; that rejection is replaced with validated behavior — never a silent no-op.

---

## Motivation / context

- Plans: `plans/v0.0.2/phase-3-config-languages.md`, parent `plans/v0.0.2/python-support.md`
- Issues: see **Related issues**
- Unblocks multi-language opt-in (Python WIP) without changing the Go-only default

---

## Changes

### Schema and templates

- `goslop.schema.json`: `languages` array (`go` | `python` | `py`), `minItems: 1`
- `templates/goslop.toml` and `internal/app/init_template.toml`: commented examples

### Config model

- `Section.Languages []string`
- Validate via `core.ParseLanguages` (unknown tokens rejected; **explicit empty list rejected**)
- `Merged.Languages []core.LanguageID` defaults to `[LanguageGo]` when unset / no config
- No CLI `--languages` override in this phase (config-first)
- `typed.enabled` still rejected

### Core helpers (Phase 2 minimal surface, self-contained)

- `DefaultEnabledLanguages()`, `ParseLanguages()`, `LanguageEnabled()`

### Engine / app

- `engine.RegistryForLanguages(base, enabled)` filters plugins; missing plugins skipped (no crash)
- Scan and `--list-rules` / `--explain` use the filtered registry
- Walk extensions come from the filtered registry `ExtensionMap`

### Docs

- `CONTRIBUTING.md` migration note updated for restored `languages` key
- Phase 3 checklist marked complete

---

## Impact

| Area | Impact |
|------|--------|
| **Performance** | None material (extra registry filter on startup only) |
| **Memory** | Negligible (filtered registry copy of plugin list) |
| **Behavior / correctness** | `languages` now affects plugins, file discovery, and list-rules; default remains Go-only |
| **API / CLI** | New config key; no new CLI flags |
| **Dependencies** | None |
| **Binary size / build time** | Negligible |

---

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| Configs that relied on `languages` being rejected as unknown | Valid tokens (`go`, `python`, `py`) now accepted and applied; unknown tokens still fail load |
| Explicit `languages = []` | Rejected with clear error (use unset for default Go, or list at least one language) |
| `typed.enabled` | Still unsupported / rejected (unchanged) |

---

## Empty-list policy

**Reject** an explicit empty `languages` array at parse time. When the key is **unset**, merge defaults to `[LanguageGo]`. Documented in schema, templates, and package comments.

---

## How language filtering works

1. Config load validates language tokens → `[]LanguageID` on `Merged`.
2. `RegistryForLanguages(DefaultRegistry(), merged.Languages)` keeps only plugins whose `ID()` is enabled.
3. Languages enabled in config but **not** registered (e.g. `python` today) are skipped — no error, no crash.
4. Analyzer walk uses the filtered registry’s `ExtensionMap()`, so only enabled-language extensions are collected.
5. `--list-rules` loads/merges config first and lists rules from the same filtered registry.

---

## Test plan

- [x] `make lint`
- [x] `make test`
- [x] `CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop`
- [x] Unit: valid languages accepted; unknown / empty rejected; typed still rejected
- [x] Unit: merge default Go; explicit `["python","go"]` ordered merge
- [x] Unit: `RegistryForLanguages` filter + missing-plugin skip
- [x] App: `languages=["go"]` ignores `.py`; `languages=["python"]` no plugin → 0 files, clean exit
- [x] App: `--list-rules --config` respects languages; invalid languages → exit config

### Commands

```sh
make lint
make test
CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop
```

---

## Screenshots / sample output

```
# default list-rules still shows Go catalogue
./bin/goslop --list-rules | head -3
# [fixture-only] CWE-...

# python-only with no plugin
./bin/goslop --list-rules --config <toml with languages=["python"]>
# no rules registered

# unknown language
./bin/goslop --list-rules --config <toml with languages=["ruby"]>
# goslop.languages: unknown language "ruby"  (exit 2)
```

---

## Related issues

- Closes #42
- Relates to #39

---

## PR metadata checklist (author)

- [x] Self-assigned (`--assignee @me`)
- [x] Labels applied (`enhancement`, `documentation`)
- [x] Related issues filled with real ticket IDs
- [x] Filled body committed under `plans/PR/v0.0.2/pr-phase-3-config.md`

---

## Follow-ups (out of scope)

- Full Python plugin / detectors (later phases)
- CLI `--languages` override
- Ruleset JSON catalogues (Phase 4)
- Product README marketing rewrite (Phase 1)

---

## Reviewer checklist

- [ ] Behavior matches summary and test plan
- [ ] No unrelated changes in diff
- [ ] Public API / CLI changes documented
- [ ] PR has assignee and labels
- [ ] Related issues use correct Closes/Relates keywords
- [ ] No secrets or generated artifacts committed
