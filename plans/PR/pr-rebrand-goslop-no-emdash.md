## Summary

Rebrand all markdown documentation from CodeHound/codehound to **goslop**, and remove em dashes (and related unicode dashes) from every `.md` file so prose matches the goslop product name and plain ASCII hyphen style.

---

## Motivation / context

- Product and GitHub repository are **goslop**; leftover CodeHound branding in docs was confusing.
- Plans: `plans/PR/PR_TEMPLATE.md` (process)
- Issues: see **Related issues**

---

## Changes

### Branding

- Replace `CodeHound` / `codehound` with `goslop` across all markdown (README, `documents/`, `plans/`, detector READMEs, fixture README).
- Soften rewrite-oriented wording in user docs (overview presents goslop as a Go SAT for PERF, BP, and CWE).

### Typography

- Remove Unicode em dashes (`—`), en dashes (`–`), and similar minus/dash characters from all markdown; use ASCII `-` spacing.

### Process

- Filled PR body under `plans/PR/pr-rebrand-goslop-no-emdash.md`.

---

## Impact

| Area | Impact |
|------|--------|
| **Performance** | None |
| **Memory** | None |
| **Behavior / correctness** | Docs only; runtime CLI/binary paths in code unchanged in this PR |
| **API / CLI** | Documentation commands use `goslop` branding |
| **Dependencies** | None |
| **Binary size / build time** | None |

---

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| None (docs only) | - |

---

## Test plan

- [x] Grep all `*.md` for `codehound` / `CodeHound` (expect zero)
- [x] Grep all `*.md` for em dash `—` (expect zero)
- [ ] Spot-check README and `documents/` for readable prose after dash replacement
- [ ] Docs-only; no `make test` required for correctness of product binary

### Commands

```sh
rg -i 'codehound' --glob '*.md' || true
rg $'\u2014' --glob '*.md' || true
```

---

## Screenshots / sample output

```
Remaining codehound in markdown: 0
Remaining em dashes in markdown: 0
```

---

## Related issues

- Relates to documentation polish after #26
- Refs product branding as goslop (repo `chinmay-sawant/goslop`)

---

## PR metadata checklist (author)

- [x] Self-assigned (`--assignee @me`)
- [x] Labels applied (`documentation`)
- [x] Related issues filled
- [x] Filled body under `plans/PR/pr-rebrand-goslop-no-emdash.md`

---

## Follow-ups (out of scope)

- Align runtime identifiers (`cmd/codehound`, `codehound.toml`, fingerprints, ignore directives) with goslop if a full code rename is desired later.

---

## Reviewer checklist

- [ ] Behavior matches summary and test plan
- [ ] No unrelated changes in diff
- [ ] Public API / CLI changes documented
- [ ] PR has assignee and labels
- [ ] Related issues use correct Closes/Relates keywords
- [ ] No secrets or generated artifacts committed
