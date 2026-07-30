## Summary

Integration branch for epic **#39** (Python language support foundation). Merges the parallel worktree streams (docs, models, config languages, rulesets, Python plugin stub) into one tree with conflict resolution and combined validation. Prefer reviewing and merging **this** PR; child PRs are superseded.

---

## Motivation / context

- Plans: [`plans/v0.0.2/python-support.md`](../../v0.0.2/python-support.md), [`plans/v0.0.2/phase-5-integration.md`](../../v0.0.2/phase-5-integration.md)
- Skills: `plans/skills/worktree-deleation/SKILL.md`, `plans/skills/phase-wise-checklist/SKILLS.md`
- Issues: see **Related issues**

---

## Child PRs (superseded by this integration)

| Phase | Issue | PR | Branch | Role |
|------:|------:|----|--------|------|
| Plan ledger | #39 | [#44](https://github.com/chinmay-sawant/goslop/pull/44) | `docs/v0.0.2-python-plan` | Phase-wise checklists |
| Plugin stub | #39/#41/#42 | [#45](https://github.com/chinmay-sawant/goslop/pull/45) | `feat/python-plugin-stub` | `internal/lang/python` + `NewRegistryWithLanguages` |
| Phase 1 docs | #40 | [#46](https://github.com/chinmay-sawant/goslop/pull/46) | `docs/python-wip-refs` | Go-only → multi-language WIP docs |
| Phase 2 models | #41 | [#47](https://github.com/chinmay-sawant/goslop/pull/47) | `feat/python-models` | `ParseLanguages`, defaults, language tests |
| Phase 3 config | #42 | [#48](https://github.com/chinmay-sawant/goslop/pull/48) | `feat/config-languages` | `languages` TOML + schema + filter |
| Phase 4 rulesets | #43 | [#49](https://github.com/chinmay-sawant/goslop/pull/49) | `feat/python-rulesets` | Python JSON seeds + reuse audit |

Merge order used: **plan/skills → docs → rulesets → plugin stub → models → config**, then integration fixes.

---

## Changes

### Integrated product surface

- **Docs:** multi-language engine + Python WIP wording (README, architecture, parity)
- **Models:** `ParseLanguages`, `DefaultEnabledLanguages`, `LanguageError`, extension/alias helpers
- **Config:** `languages = ["go"]` / `["go","python"]`; reject unknown / empty list; schema + templates
- **Engine:** `NewRegistryWithLanguages`, `RegistryForLanguages`, Go-only `DefaultRegistry`
- **App:** compose scan/`--list-rules` registry via `NewRegistryWithLanguages(merged.Languages...)` so enable/disable is real (not filter-only on Go-only default)
- **Python plugin stub:** source-only `LanguagePlugin`, zero detectors, `.py` when enabled
- **Rulesets:** `ruleset/python/` CWE-22/78/79/89 seeds + empty BP + reuse audit
- **Skills:** phase-wise checklist + worktree multi-stream skill docs under `plans/skills/`

### Integration-specific conflict resolutions

- Unified `internal/lang/python` from #45 and #47
- Unified `ParseLanguages` empty-list policy with config (error, not nil)
- Combined registry tests (Go-only default, multi-lang helpers, filter helpers)
- App path uses **built-in** composition so `languages = ["python"]` scans `.py` (not zero plugins)

### Not in this PR (follow-ups)

- Full Python PERF/CWE/BP detector port
- Python AST / pure-Go parser beyond source-only
- CLI `--languages` flag (config-first in this epic)

---

## Impact

| Area | Impact |
|------|--------|
| **Performance** | Negligible; default remains Go-only walk |
| **Memory** | Unchanged for default scans |
| **Behavior / correctness** | Opt-in Python language enable; default Go-only |
| **API / CLI** | New TOML `languages`; no new required flags |
| **Dependencies** | None |
| **Binary size / build time** | Small stub package only |

---

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| Previously rejected `languages` key | Now accepted when valid; invalid tokens fail load |
| None for default scans | Unset `languages` → Go only (same as before) |

---

## Test plan

- [x] `make lint` — green
- [x] `make test` — green (all packages + integration)
- [x] `CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop` — green
- [x] Smoke: default scan ignores `.py`; `languages = ["python"]` scans only `.py`

### Commands

```sh
make lint
make test
CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop
```

---

## Screenshots / sample output

```
# default (Go-only): scanned 1 files (.go), skipped .py
# languages = ["python"]: scanned 1 files (.py), skipped .go
```

---

## Related issues

- Closes #40
- Closes #41
- Closes #42
- Closes #43
- Closes #39

---

## PR metadata checklist (author)

- [x] Self-assigned (`--assignee @me`)
- [x] Labels applied
- [x] Related issues filled with real ticket IDs
- [x] Filled body committed under `plans/PR/v0.0.2/pr-epic-39-integration.md`

---

## Follow-ups (out of scope)

- Python detector catalogue implementation
- Optional CLI `--languages`
- Shared `ruleset/shared/` pure CWE shells (audit recommended; not required)

---

## Reviewer checklist

- [ ] Behavior matches summary and test plan
- [ ] No unrelated changes in diff
- [ ] Public API / CLI changes documented
- [ ] Default remains Go-only without config
- [ ] PR has assignee and labels
- [ ] Related issues use correct Closes keywords
- [ ] No secrets or generated artifacts committed
