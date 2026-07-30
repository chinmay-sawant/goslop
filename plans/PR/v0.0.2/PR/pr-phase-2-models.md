## Summary

Harden multi-language domain types for Python WIP: language list parsing helpers, Go-only default enabled set, documented `ParsedUnit.Language` / registry indexing, and a tests-only source-only Python plugin stub that is **not** registered in production `DefaultRegistry`.

## Motivation / context

- Plans: `plans/v0.0.2/python-support.md`, `plans/v0.0.2/phase-2-models.md`
- Issues: see **Related issues**
- Unblocks Phase 3 (`languages` config merge → `[]LanguageID`) without wiring TOML yet

## Changes

### Core language models (`internal/core`)

- Audit `LanguageID` / `ParseLanguage` / `LanguageFromExtension` (aliases `python`/`py`; extension `py`; optional leading-dot ext)
- Add `ParseLanguages([]string) ([]LanguageID, error)` — stable first-seen order, alias dedupe, clear unknown-token errors, blank skip, empty → `(nil, nil)`
- Add `DefaultEnabledLanguages()` → `[LanguageGo]` (documented production default)
- Document `BasePlugin.ParseSource` forces `LanguageGo`; non-Go plugins must override
- Unit tests: `internal/core/language_test.go`

### Registry docs / tests (`internal/engine`)

- Document `Registry` indexes plugins and detectors by `LanguageID` (`byID`, `byLanguage`, `byExtension`)
- Document `DefaultRegistry` is Go-only (Python deferred until language filter)
- Tests: `TestDefaultRegistryGoOnly`, `TestRegistryIndexesDetectorsByLanguageID`

### Optional Python stub (`internal/lang/python`)

- Minimal `LanguagePlugin`: `Extensions() == ["py"]`, zero detectors, source-only parse with `LanguagePython`
- **Not** composed into `DefaultRegistry` (tests register via `NewRegistry` only)

### Finding / cache

- Audited only: no concrete gap for fingerprint/cache language fields; walk/session already key detectors by `LanguageID`. No schema change.

## Impact

| Area | Impact |
|------|--------|
| **Performance** | None (helpers only; no scan path change) |
| **Memory** | Negligible (small slices) |
| **Behavior / correctness** | Production scan path unchanged (Go-only registry) |
| **API / CLI** | No CLI flags; internal helpers for Phase 3 |
| **Dependencies** | None |
| **Binary size / build time** | Trivial (`internal/lang/python` not linked into default app path unless imported) |

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| None | - |

## Test plan

- [x] `make test` — all packages ok including `internal/core`, `internal/engine`, `internal/lang/python`, `internal/fixture`, `tests/integration`
- [x] `make lint` — `go vet ./...` + `gofmt -l` clean
- [x] `CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop`
- [x] `TestMaterializeTree_IncludesPythonTagged` still passes
- [x] `TestDefaultRegistryGoOnly` asserts no Python plugin / no `.py` in default extension map
- [ ] `make run` wall time vs baseline (not required — no detector/scan path change)

### Commands

```sh
gofmt -w internal/core/language.go internal/core/language_test.go internal/core/plugin.go \
  internal/engine/registry.go internal/engine/registry_test.go \
  internal/lang/python/plugin.go internal/lang/python/plugin_test.go
make lint
make test
CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop
```

### Recorded outcomes (2026-07-30)

```text
make lint  → go vet ./... + gofmt clean (exit 0)
make test  → go test ./... all ok (internal/lang/python ok; integration ok)
CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop → success
```

## Screenshots / sample output

N/A (domain helpers; no CLI UX change). Default product behavior remains Go-only.

## Related issues

- Closes #41
- Relates to #39

## PR metadata checklist (author)

- [x] Self-assigned (`--assignee @me`)
- [x] Labels applied (`enhancement`, `documentation`)
- [x] Related issues filled with real ticket IDs
- [x] Filled body committed under `plans/PR/v0.0.2/pr-phase-2-models.md`

## Follow-ups (out of scope)

- Phase 3: TOML `languages` key, schema, merge, registry filter (`#42`)
- Phase 4: Python ruleset layout / JSON catalogues (`#43`)
- Full Python detectors / pure-Go AST parse
- Registering Python in `DefaultRegistry` (blocked on Phase 3 filter)

## Reviewer checklist

- [ ] Behavior matches summary and test plan
- [ ] No unrelated changes in diff
- [ ] Public API / CLI changes documented
- [ ] New rules have fixture coverage when applicable (N/A)
- [ ] PR has assignee and labels
- [ ] Related issues use correct Closes/Relates keywords
- [ ] No secrets or generated artifacts committed
- [ ] `DefaultRegistry` remains Go-only; Python stub tests-only
