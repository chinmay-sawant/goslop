# v0.0.2 — Python heuristic detectors (CWE, BP, PERF)

> **Parent:** GitHub epic [#51](https://github.com/chinmay-sawant/goslop/issues/51) — Epic: Python heuristic detectors (CWE, BP, PERF)  
> **Status:** planning (foundation #39 / PR #50 shipped; detectors not started)  
> **Estimated effort:** multi-PR (CWE + BP active; PERF deferred)  
> **Ledger rule:** this file is the **canonical execution ledger** for heuristic implementation. Detail lives in sibling checklists. Update checkboxes only with evidence (`make lint` / `make test` for code).

---

## Overview

| Stream | Issue | Detail plan | Status |
|--------|------:|-------------|--------|
| CWE heuristics | #52 | [python-heuristics-cwe.md](./python-heuristics-cwe.md) | active plan |
| Bad-practice heuristics | #53 | [python-heuristics-bp.md](./python-heuristics-bp.md) | active plan |
| PERF catalogue + heuristics | #54 | [python-heuristics-perf.md](./python-heuristics-perf.md) | **deferred** (`[~]`) |

**Foundation already shipped** (do not re-plan): epic #39 / PR #50 — `languages` config, Python plugin stub, `ruleset/python/` CWE + BP catalogues. See [python-support.md](./python-support.md).

**Product default:** Go-only. Python remains opt-in via `languages = ["python"]` or `["go","python"]`.

---

## Executive Summary

### Current evidence (main)

| Fact | Evidence |
|------|----------|
| Python plugin source-only + **BP-PY detectors** (A+B+C subset) | `internal/lang/python/` + `detectors/bad_practices/` |
| CWE catalogue ~344 rules | `ruleset/python/chunks/cwe-*.json` (`python_relevance`) |
| BP catalogue 50 `BP-PY-*` | `ruleset/python/bad-practices.json` |
| **No** Python PERF JSON | no `ruleset/python/chunks/perf-*` |
| Go detector families | `internal/lang/go/detectors/{cwe,bad_practices,perf}/` |

### Strategy

1. **CWE first** (security correctness) — priority batch CWE-22 / 78 / 79 / 89 / 502, then expand.  
2. **BP second** — core + security hygiene + high-signal Flask/Django/FastAPI.  
3. **PERF deferred** — no catalogue yet; issue #54 stays open as backlog until PERF JSON seeds exist.

### Non-goals

- Full parity with Go’s entire PERF/CWE/BP surface in one epic  
- CGO / tree-sitter requirement for v0.0.2 heuristics (prefer pure-Go pattern / light parse)  
- Re-doing foundation config/docs from #39  

### Dependency graph

```text
#39 foundation (done)
    ├─ #52 CWE heuristics  ──► fixtures + plugin Detectors()
    ├─ #53 BP heuristics   ──► same plugin registration surface
    └─ #54 PERF            ──► [~] blocked: no perf-*.json yet
```

---

## Phase index (rollup)

### CWE — #52

- [x] Detail checklist complete in [python-heuristics-cwe.md](./python-heuristics-cwe.md) (docs-only plan, 2026-07-31)
- [ ] Priority CWE batch green on fixtures
- [ ] Validation: `make lint` + `make test` recorded on implement branch

### BP — #53

- [x] Detail checklist complete in [python-heuristics-bp.md](./python-heuristics-bp.md) (docs-only plan, 2026-07-31)
- [x] Priority `BP-PY-*` batch green (A+B + C high-signal; inline unit tests; D/E deferred)
- [x] Validation: `make lint` + `make test` green on `feat/python-bp-heuristics` (2026-07-31)

### PERF — #54 (deferred)

- [x] Deferred ledger written in [python-heuristics-perf.md](./python-heuristics-perf.md) (docs-only, 2026-07-31)
- [~] All **implementation** deferred — no `ruleset/python/chunks/perf-*` yet  
  **Reason:** no Python PERF catalogue.  
  **Owner:** issue #54 / epic #51.  
  **Next gate:** seed PERF JSON (≤50 rules/file, golang-style range names) before any detector work.

### Epic closure

- [ ] #52 and #53 have shipped initial heuristic batches (or explicit partial `[~]` with remaining IDs listed)
- [ ] #54 remains deferred or unblocked with catalogue + first detectors
- [ ] Epic #51 success criteria reviewed against evidence

---

## Validation gates (when implementing)

```sh
gofmt -w <changed Go files>
make lint
make test
CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop
# optional smoke
./bin/goslop --format text --no-cache <python-fixture-root>  # with languages=["python"]
```

Docs-only plan edits: no lint/test required.

---

## Issues map

| Issue | Title | Plan |
|------:|-------|------|
| #51 | Epic: Python heuristic detectors (CWE, BP, PERF) | this file |
| #52 | python(cwe): implement CWE heuristic detectors | [python-heuristics-cwe.md](./python-heuristics-cwe.md) |
| #53 | python(bp): implement BP-PY heuristics | [python-heuristics-bp.md](./python-heuristics-bp.md) |
| #54 | python(perf): seed PERF + heuristics | [python-heuristics-perf.md](./python-heuristics-perf.md) (deferred) |

---

## Dependencies

| Depends on | Note |
|------------|------|
| #39 / PR #50 | Language seam, config, catalogues |
| `core.LanguagePlugin` | Detectors must register on Python plugin |
| `ruleset/python/**` | Metadata source of truth for IDs |

---

## References

- Issue bodies: `plans/PR/v0.0.2/issue-epic-python-heuristics-body.md`, `issue-python-*-heuristics-body.md`  
- Skill: `plans/skills/phase-wise-checklist/SKILLS.md`  
- CWE mapping: `plans/v0.0.2/python-cwe-from-699-mapping.md`  
- Foundation ledger: `plans/v0.0.2/python-support.md`  
