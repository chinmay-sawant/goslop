# Phase 5 — Integration gates (epic #39)

> **Parent:** `plans/v0.0.2/python-support.md`  
> **Issue:** #39 (epic)  
> **Status:** blocked until Phases 1–4 PRs exist  
> **Branch:** `chore/epic-39-integration`  
> **Base branch:** `main`

---

## Overview

Merge parallel workstream heads into one integration branch, validate composition, open a single integration PR to `main` per `plans/PR/PR_TEMPLATE.md` multi-workstream section and `plans/skills/worktree-deleation/SKILL.md` (use **`main`** not `master` for this repo).

---

## 5.1 Merge order (recommended)

1. `docs/python-wip-refs` (Phase 1) — low conflict  
2. `feat/python-rulesets` (Phase 4) — mostly `ruleset/` + docs  
3. `feat/python-models` (Phase 2) — core types  
4. `feat/config-languages` (Phase 3) — config + engine filter (may conflict with 2)

After each merge: resolve conflicts, re-run tests when code present.

## 5.2 Combined validation

- [ ] `make lint`  
- [ ] `make test`  
- [ ] `CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop`  
- [ ] Smoke: default scan still Go-oriented; with config `languages = ["go"]` behaves  
- [ ] Record commands + outcomes in integration PR body  

## 5.3 Integration PR

- [ ] Filled body `plans/PR/v0.0.2/pr-epic-39-integration.md`  
- [ ] Table of child issue / branch / standalone PR  
- [ ] `Closes #40` `#41` `#42` `#43` and `Closes #39` when foundation complete  
- [ ] Note child PRs superseded by integration  
- [ ] Comment on each child PR pointing at integration  

## 5.4 Parent ledger sync

- [ ] Update `python-support.md` rollup checkboxes with evidence  
- [ ] Update phase file rows that shipped  

## 5.5 Land (only when user requests merge)

- [ ] Merge **integration PR only** into `main`  
- [ ] Close child PRs without merging (if still open)  
- [ ] Delete all non-`main` local/remote branches  
- [ ] `git pull origin main` — clean status  

---

## Success criteria

- One green integration PR composes all foundation streams  
- Default behavior remains Go-first  
- Plan ledger matches code  

## Dependencies

Requires Phase 1–4 branches pushed.
