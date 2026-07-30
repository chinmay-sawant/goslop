## Summary

Reframe public and plan docs that permanently claimed Go-only product scope: Go remains the shipped production language and detector catalogue surface; the engine is multi-language (`LanguagePlugin`); Python is reserved / WIP under epic #39 with no claim of full detector support.

## Motivation / context

- Plans: `plans/v0.0.2/phase-1-docs.md`, `plans/v0.0.2/python-support.md`
- Issues: see **Related issues**
- Product README / overview and architecture “Language seam” still read as forever-Go-only while `LanguagePython`, python fixtures, and the plugin seam already exist.

## Changes

### Product docs

- `README.md` — multi-language engine + Go default + Python WIP (#39); Go-labeled PERF/CWE/BP catalogue counts
- `documents/overview.md` — same accuracy pass; status table notes Python foundation WIP
- `documents/README.md` — index no longer claims permanent “SAT for Go only”
- `documents/architecture-performance.md` § Language seam — shipped Go vs Python WIP; pure-Go / no-CGO plugin rule; cross-link `internal/core/plugin.go`

### Plans / parity

- `plans/parity-matrix.md` — `lang/python/` and python fixtures point at epic #39 / `plans/v0.0.2/python-support.md`
- `plans/architecture-go.md` — Python deferral points at #39 / v0.0.2 (not “defer forever”)
- `plans/port-phasewise-checklist.md` — Python reserved bullets + next-action note for v0.0.2 foundation

### v0.0.2 ledger (from plan branch)

- Adds `plans/v0.0.2/**` phase-wise checklist (docs, models, config, rulesets, integration, audit)
- Filled PR body under `plans/PR/v0.0.2/pr-phase-1-docs.md`
- Phase 1 checkboxes marked complete in `phase-1-docs.md` and parent rollup

### Explicit non-changes

- No `languages=` config keys (Phase 3 owns schema)
- No claim of full Python detector support
- `templates/goslop.toml` unchanged (no invented `languages` key)
- `plans/v0.0.1/reviews/**` untouched

## Impact

| Area | Impact |
|------|--------|
| **Performance** | None (docs only) |
| **Memory** | None |
| **Behavior / correctness** | None — documentation accuracy only |
| **API / CLI** | None |
| **Dependencies** | None |
| **Binary size / build time** | None |

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| None | Docs reframe only; default product remains Go scanning |

## Test plan

Docs-only stream (per `plans/v0.0.2/phase-1-docs.md`): no `make lint` / `make test` required.

- [x] README / overview no longer assert permanent Go-only product scope
- [x] Go catalogue counts remain explicitly Go-labeled
- [x] Architecture language seam: shipped Go vs Python WIP; pure-Go / no-CGO; link to `LanguagePlugin`
- [x] Parity matrix + architecture-go + port checklist point Python work at #39 / v0.0.2
- [x] Scanned `documents/*.md` for permanent Go-only product claims; fixed product-facing hits
- [x] Did not invent `languages` config keys or full Python detector support
- [x] Left `plans/v0.0.1/reviews/**` untouched

### Commands

```sh
# Docs-only; optional smoke that tree still builds:
# CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop
rg -n 'only shipped language|static analysis tool \(SAT\) for \*\*Go\*\*|out of scope v0' README.md documents plans --glob '!plans/v0.0.1/reviews/**' || true
```

## Screenshots / sample output

N/A (documentation only).

## Related issues

- Closes #40
- Relates to #39

## PR metadata checklist (author)

- [x] Self-assigned (`--assignee @me`)
- [x] Labels applied (`documentation`, `enhancement`)
- [x] Related issues filled with real ticket IDs
- [x] Filled body committed under `plans/PR/v0.0.2/pr-phase-1-docs.md`

## Follow-ups (out of scope)

- Phase 2 models (#41)
- Phase 3 config `languages` key + schema (#42)
- Phase 4 Python rulesets + reuse audit (#43)
- Phase 5 epic integration (#39)
- Full Python detector catalogue parity (explicit non-goal of foundation epic)

## Reviewer checklist

- [ ] Behavior matches summary and test plan (docs accuracy)
- [ ] No unrelated code/config changes in diff
- [ ] No invented `languages=` keys or full Python support claims
- [ ] PR has assignee and labels
- [ ] Related issues use Closes #40 / Relates to #39
- [ ] No secrets or generated artifacts committed
