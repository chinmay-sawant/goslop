## Summary

Adds Python ruleset catalogue seeds under `ruleset/python/` (portable CWE-22/78/79/89 shells + empty BP map) and expands the Go→Python ruleset reuse audit with scripted counts and classification evidence. No Python detectors; Go catalogues and generators unchanged.

---

## Motivation / context

- Plans: `plans/v0.0.2/phase-4-rulesets.md`, `plans/v0.0.2/ruleset-reuse-audit.md`, `plans/v0.0.2/python-support.md`
- Issues: see **Related issues**
- Main had **0** `applicable_to: python` tags and no `ruleset/python/` tree; Phase 4 establishes layout + reuse policy before any AST work.

---

## Changes

### Reuse audit (evidence)

- Exact chunk counts (CWE 175 / PERF 242 / BP 135)
- `applicable_to` / `go_relevance` / `status` distributions
- Classification A–E with portable CWE table (20 IDs) and BP category portability notes
- Confirmed hard-coded Go metadata paths (`metadata_gen.go`, `metadata_generated.go`)

### Python catalogue seeds

- `ruleset/python/README.md` — WIP + reuse policy
- `ruleset/python/bad-practices.json` — `{}`
- `ruleset/python/chunks/cwe-seed.json` — CWE-22, 78, 79, 89 with `applicable_to: ["python"]` and rewritten Python `detection_notes`
- `ruleset/python/catalogue_test.go` — parse/shape smoke test

### Docs

- `documents/rule-catalog-and-maturity.md` — multi-language catalogue roots
- `documents/overview.md` — metadata path note
- `plans/parity-matrix.md` — python fixtures / lang row → #43
- Phase 4 checklist marked complete with evidence

---

## Impact

| Area | Impact |
|------|--------|
| **Performance** | None |
| **Memory** | None |
| **Behavior / correctness** | None for Go scan path; Python JSON not loaded by runtime |
| **API / CLI** | None |
| **Dependencies** | None |
| **Binary size / build time** | Negligible (optional test package only) |

---

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| None | - |

---

## Test plan

- [x] `go test ./ruleset/python/` — Python catalogue JSON parses
- [x] Confirmed `ruleset/golang/**` untouched
- [ ] Full `make test` / `make lint` — not required for docs/JSON-only scope; tiny Go test only added under `ruleset/python/`
- [ ] No detector surface change → no `make run` / reference-metrics delta expected

### Commands

```sh
go test ./ruleset/python/
# optional full suite if desired:
# make lint && make test
```

### Validation note

**Mostly docs + JSON.** One small pure-Go test for JSON parse/shape. No Go catalogue or detector registration changes.

---

## Screenshots / sample output

```
go test ./ruleset/python/
# ok  .../ruleset/python  0.00Xs
```

Portable seed rule count: **4** CWE shells (22, 78, 79, 89); BP **0**.

---

## Related issues

Closes #43  
Relates to #39  
