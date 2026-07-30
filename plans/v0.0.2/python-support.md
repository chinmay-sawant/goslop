# v0.0.2 — Introduce Python language support (foundation)

> **Parent:** GitHub epic [#39](https://github.com/chinmay-sawant/goslop/issues/39) — Epic: Introduce Python language support  
> **Status:** planning → parallel worktrees (foundation only; not full detector parity)  
> **Estimated effort:** multi-PR epic (4 workstreams + integration)  
> **Ledger rule:** this file is the **canonical execution ledger**. Phase detail lives in sibling markdowns; update checkboxes only with evidence (`make lint` / `make test` when code changes).

---

## Overview

goslop ships **Go** only today. The engine already has a language-agnostic seam (`core.LanguagePlugin`, `LanguagePython` reserved, python fixtures). This epic makes multi-language **honest and operable**:

| Stream | Issue | Detail plan | Branch |
|--------|------:|-------------|--------|
| Phase 1 — Docs / references | #40 | [phase-1-docs.md](./phase-1-docs.md) | `docs/python-wip-refs` |
| Phase 2 — Domain models | #41 | [phase-2-models.md](./phase-2-models.md) | `feat/python-models` |
| Phase 3 — Config languages | #42 | [phase-3-config-languages.md](./phase-3-config-languages.md) | `feat/config-languages` |
| Phase 4 — Rulesets + reuse | #43 | [phase-4-rulesets.md](./phase-4-rulesets.md) | `feat/python-rulesets` |
| Phase 5 — Integration gates | #39 | [phase-5-integration.md](./phase-5-integration.md) | `chore/epic-39-integration` |

Supporting research (ruleset audit): [ruleset-reuse-audit.md](./ruleset-reuse-audit.md) (owned by Phase 4; may be drafted by a dedicated scan agent).

**Default product behavior after this epic:** `languages = ["go"]` (or equivalent). Python is **opt-in / WIP**.

---

## Executive Summary

### Why now

- Epic #39 + children #40–#43 exist as placeholders.
- `LanguagePython` and `tests/fixtures/python/` already exist; docs still say Go-only.
- Inert `languages` config was **rejected** on purpose; reintroduce only with real wiring.
- `ruleset/golang/` is entirely Go-tagged (`applicable_to` has no `python`); need layout + reuse policy before detectors.

### Non-goals (all phases)

- Full PERF/CWE/BP Python detector port
- Tree-sitter / CGO Python parse
- Production-ready Python catalogue parity
- Re-introducing dead `typed.enabled` without a real backend

### Dependency graph

```text
Phase 1 docs ─────────────────────────────┐
Phase 2 models ──► Phase 3 config ────────┤──► Phase 5 integration
Phase 4 rulesets (audit + stubs) ─────────┘
```

Phases 1 and 4 are mostly independent of 2/3. Phase 3 depends on model decisions for `[]LanguageID` / merge fields. Integration merges all heads.

### Parallel worktree policy

Follow `plans/skills/worktree-deleation/SKILL.md` adapted to this repo’s default branch **`main`** (not `master`):

1. One branch + PR per Phase 1–4 workstream  
2. Filled PR bodies under `plans/PR/v0.0.2/` from `plans/PR/PR_TEMPLATE.md`  
3. Single integration branch `chore/epic-39-integration` → PR to `main`  
4. Prefer merge **integration only**; close child PRs as superseded; delete non-`main` branches on land  

---

## Phase index (rollup)

### Phase 1 — Docs / multi-language WIP — #40

- [x] Ledger rows complete in [phase-1-docs.md](./phase-1-docs.md)
- [x] PR open with filled body `plans/PR/v0.0.2/pr-phase-1-docs.md`

### Phase 2 — Models — #41

- [ ] Ledger rows complete in [phase-2-models.md](./phase-2-models.md)
- [ ] PR open with filled body `plans/PR/v0.0.2/pr-phase-2-models.md`
- [ ] Validation: `make lint` + `make test` recorded

### Phase 3 — Config enable/disable languages — #42

- [ ] Ledger rows complete in [phase-3-config-languages.md](./phase-3-config-languages.md)
- [ ] PR open with filled body `plans/PR/v0.0.2/pr-phase-3-config.md`
- [ ] Validation: `make lint` + `make test` recorded

### Phase 4 — Rulesets + generic reuse audit — #43

- [ ] [ruleset-reuse-audit.md](./ruleset-reuse-audit.md) written with generic vs Go-only classification
- [ ] Ledger rows complete in [phase-4-rulesets.md](./phase-4-rulesets.md)
- [ ] PR open with filled body `plans/PR/v0.0.2/pr-phase-4-rulesets.md`
- [ ] Validation: `make lint` + `make test` if Go loaders change; else note docs-only

### Phase 5 — Integration — #39

- [ ] Child branches merged into `chore/epic-39-integration`
- [ ] Combined `make lint` + `make test` green
- [ ] Integration PR body `plans/PR/v0.0.2/pr-epic-39-integration.md`
- [ ] Plan checkboxes synchronized to evidence

---

## Issues map

| Issue | Title | Closes via |
|------:|-------|------------|
| #39 | Epic: Introduce Python language support | Integration PR when foundation done |
| #40 | docs: reframe Go-only references as multi-language with Python WIP | Phase 1 PR / integration |
| #41 | models: extend domain types for multi-language / Python support | Phase 2 PR / integration |
| #42 | config: enable/disable languages (languages key + schema) | Phase 3 PR / integration |
| #43 | ruleset: add Python JSON catalogues and reuse shared rules | Phase 4 PR / integration |

---

## Dependencies

| Depends on | Note |
|------------|------|
| #19 (closed) | Pluggable languages / pure-Go parse seam |
| `internal/core/plugin.go` | LanguagePlugin contract already documented |
| `internal/core/language.go` | `LanguagePython` reserved |
| Rejected inert config | `TestUnsupportedLanguageAndTypedConfigurationRejected` must be **rewritten**, not ignored |

---

## Validation gates (every code stream + integration)

```sh
gofmt -w <changed Go files>
make lint
make test
CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop
```

Docs-only streams: no lint/test required per phase-wise-checklist skill; still build if unsure.

---

## References

- Epic: https://github.com/chinmay-sawant/goslop/issues/39  
- Skills: `plans/skills/worktree-deleation/SKILL.md`, `plans/skills/phase-wise-checklist/SKILLS.md`  
- PR process: `plans/PR/PR_TEMPLATE.md`  
- Architecture seam: `documents/architecture-performance.md`  
- Parity deferral: `plans/parity-matrix.md` (`lang/python/`)  
- Prior integration example: `plans/PR/v0.0.1/pr-epic-phases-7-12-integration.md`  
