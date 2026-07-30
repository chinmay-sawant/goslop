## Summary

Adds a contributor guide that turns the repository's recently established review, checklist, validation, benchmark, migration, and pull-request practices into a discoverable, repeatable workflow.

---

## Motivation / context

- Plan and review examples: `plans/v0.0.1/reviews/`
- Pull-request process: `plans/PR/PR_TEMPLATE.md`
- Issues: none supplied; this is documentation/process work.

---

## Changes

### Contributor workflow

- Add `CONTRIBUTING.md` with pure-Go setup, focused branch guidance, change-sizing criteria, and implementation expectations.
- Explain when a phase-wise checklist is useful and how `[x]`, `[ ]`, and `[~]` states must be maintained.
- Link the existing Ponytail, architecture, and Go-style/design review reports as concrete examples.

### Validation and reproducibility

- Document lint, tests, race tests, pure-Go build verification, product scans, and Go benchmark commands.
- Record the GopdfSuit `make run` output shape and current 915-finding reference baseline while keeping machine-dependent timing separate from parity requirements.
- Explain the `reference-metrics` findings exit behavior and its hard count/severity/export targets.

### Compatibility and PR hygiene

- Document the intentional removal of unsupported `languages` and `typed.enabled` configuration promises and the required `goslop.toml` migration.
- Require template-backed PR bodies, self-assignment, labels, accurate validation reporting, and maintainer-controlled merges.

---

## Impact

| Area | Impact |
|------|--------|
| **Performance** | No runtime change; benchmark requirements are documented. |
| **Memory** | No change. |
| **Behavior / correctness** | No product behavior change; contributors receive reproducible scan and validation guidance. |
| **API / CLI** | No change. |
| **Dependencies** | None. |
| **Binary size / build time** | No change. |

---

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| None introduced by this documentation PR | The guide documents the existing removal of inert `languages` and `typed.enabled` settings. |

---

## Test plan

- [x] `git diff --check`
- [x] Review all repository-relative links and commands against the current Makefile, PR template, and review artifacts.
- [x] `make run` — current GopdfSuit reference output matches 915 findings, the expected severity distribution, top rules, and 915 + 37 exports.
- [ ] Code tests not run; this PR changes documentation only.

### Commands

```sh
git diff --check
make run
```

---

## Screenshots / sample output

No UI change. `CONTRIBUTING.md` includes the expected `make run` and reference-metrics output format and GopdfSuit parity baseline.

---

## Related issues

- None supplied or linked.

---

## PR metadata checklist (author)

- [x] Self-assigned (`--assignee @me`)
- [x] Labels applied (`documentation`)
- [x] Related issues reviewed; no matching ticket was supplied
- [x] Filled body committed under `plans/PR/pr-contributing-guide.md`

---

## Follow-ups (out of scope)

- Add a dedicated security disclosure policy if maintainers choose a reporting channel.

---

## Reviewer checklist

- [ ] Guidance matches the Makefile, PR template, and current review artifacts.
- [ ] Commands and expected baseline output are clearly labeled as reproducible evidence rather than universal timing guarantees.
- [ ] No unrelated repository behavior changes.
- [ ] PR has assignee and labels.
