## Summary

Adds a **minimal source-only Python `LanguagePlugin` stub** under `internal/lang/python/` so multi-language enable/disable workstreams have a real second plugin to filter against. Production `DefaultRegistry` remains **Go-only**; Python is composed via `engine.NewRegistryWithLanguages` or `python.Register`.

---

## Motivation / context

- Plans: `plans/v0.0.2/python-support.md`, `plans/v0.0.2/phase-2-models.md` (§2.3 optional stub seam)
- Issues: see **Related issues**
- Complementary workstream for epic #39: Phase 3 config filtering needs more than Go in the plugin catalogue without turning Python on by default

---

## Changes

### `internal/lang/python/`

- New package implementing `core.LanguagePlugin`:
  - `ID()` → `LanguagePython`
  - `Extensions()` → `["py"]`
  - `Detectors()` / `NewDetectors()` → empty (no rules yet)
  - `ParseSource` → source-only `ParsedUnit` with **LanguagePython** (overrides `BasePlugin` Go default)
- `Register(reg any)` hook mirrors `golang.Register`
- Package doc marks WIP / no detectors

### `internal/engine/registry.go`

- `builtInPlugin` + `NewRegistryWithLanguages(...LanguageID)` for merge-friendly multi-language composition
- `DefaultRegistry` still constructs **Go only** (explicit `golang.NewPlugin()`, not via the multi-lang helper)

### Tests

- Plugin ID/extensions/empty detectors
- Parse sets `LanguagePython` + source-only quality
- Registry with Go+Python resolves `.py` / `.go`
- Default registry remains Go-only
- `python.Register` hook
- Unknown language id error path for `NewRegistryWithLanguages`

---

## Impact

| Area | Impact |
|------|--------|
| **Performance** | None (Python not in default registry) |
| **Memory** | None on default path |
| **Behavior / correctness** | Unchanged production scans; opt-in registries can resolve `.py` with zero detectors |
| **API / CLI** | No CLI surface; new engine helper `NewRegistryWithLanguages` |
| **Dependencies** | None |
| **Binary size / build time** | Negligible (empty stub package) |

---

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| None | - |

---

## Test plan

- [x] `make test`
- [x] `make lint` / `go vet`
- [x] `CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop`
- [ ] `make run` wall time vs baseline (hard &lt; 400ms; soft ±50ms of reference) — N/A product surface unchanged
- [ ] `make reference-metrics` / gopdfsuit hard metrics if detector surface changed — N/A

### Commands

```sh
make lint
make test
CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop
```

### Recorded outcome

- `make lint` green
- `make test` green (includes `internal/lang/python`, `internal/engine`)
- pure-Go build OK

---

## Screenshots / sample output

```
ok  github.com/chinmay-sawant/goslop/internal/lang/python
ok  github.com/chinmay-sawant/goslop/internal/engine
# DefaultRegistry: Go only; NewRegistryWithLanguages(go, python) resolves .py
```

---

## Related issues

- Relates to #39
- Relates to #41
- Relates to #42

---

## PR metadata checklist (author)

- [x] Self-assigned (`--assignee @me`)
- [x] Labels applied (`enhancement`, `documentation`)
- [x] Related issues filled with real ticket IDs (Relates only for #41/#42; those close via sibling PRs)
- [x] Filled body under `plans/PR/v0.0.2/pr-python-plugin-stub.md`

---

## Follow-ups (out of scope)

- Detectors / AST parse for Python
- Config TOML `languages` key (Phase 3 / #42)
- Ruleset JSON catalogues (Phase 4 / #43)
- Docs marketing reframe (Phase 1 / #40)
- Domain helpers `DefaultEnabledLanguages` / list parse if Phase 2 models PR does not land them

---

## Reviewer checklist

- [ ] Behavior matches summary and test plan
- [ ] No unrelated changes in diff
- [ ] Public API / CLI changes documented
- [ ] New rules have fixture coverage when applicable (N/A — zero detectors)
- [ ] PR has assignee and labels
- [ ] Related issues use correct Closes/Relates keywords
- [ ] No secrets or generated artifacts committed
