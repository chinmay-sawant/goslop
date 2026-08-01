# v0.0.2 — Python heuristic detectors (CWE, BP, PERF)

> **Parent:** GitHub epic [#51](https://github.com/chinmay-sawant/goslop/issues/51) — Epic: Python heuristic detectors (CWE, BP, PERF)  
> **Status:** CWE #52 priority batch shipped; CWE expansion planned in `cwe-plans/`; BP #53 **complete** 50/50; PERF #54 **complete** (experimental detectors + 2026-08-01 canary; no pack promotion)
> **Estimated effort:** multi-PR (CWE expansion waves P0–P3 remain)  
> **Ledger rule:** this file is the **canonical execution ledger** for heuristic implementation. Detail lives in sibling checklists. Update checkboxes only with evidence (`make lint` / `make test` for code).

---

## Overview

| Stream | Issue | Detail plan | Status |
|--------|------:|-------------|--------|
| CWE heuristics | #52 | [python-heuristics-cwe.md](./python-heuristics-cwe.md) + **[cwe-plans/](./cwe-plans/)** | priority 5 shipped; **154** implement-owned + **185** deferred |
| Bad-practice heuristics | #53 | [python-heuristics-bp.md](./python-heuristics-bp.md) + [bp-plans/](./bp-plans/) | **complete** 50/50 on main (PR #65) |
| PERF catalogue + heuristics | #54 | [python-heuristics-perf.md](./python-heuristics-perf.md) | **complete** — experimental detectors + canary ([PERF-PY-CANARY-2026-08-01.md](./pref-plans/PERF-PY-CANARY-2026-08-01.md)) |

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
| Python PERF catalogue + detectors | `ruleset/python/chunks/perf-py-001-014.json` + `perf-py-015-022.json`; `internal/lang/python/detectors/perf/` (22/22 fixture-backed) |
| Go detector families | `internal/lang/go/detectors/{cwe,bad_practices,perf}/` |

### Strategy

1. **CWE first** (security correctness) — priority batch CWE-22 / 78 / 79 / 89 / 502, then expand.  
2. **BP second** — core + security hygiene + high-signal Flask/Django/FastAPI.  
3. **PERF experimental batch** — issue #54 closed experimentally: fixture-backed rules + canary triage; maturity stays experimental (not in recommended/perf).

### Non-goals

- Full parity with Go’s entire PERF/CWE/BP surface in one epic  
- CGO / tree-sitter requirement for v0.0.2 heuristics (prefer pure-Go pattern / light parse)  
- Re-doing foundation config/docs from #39  

### Dependency graph

```text
#39 foundation (done)
    ├─ #52 CWE heuristics  ──► fixtures + plugin Detectors()
    ├─ #53 BP heuristics   ──► same plugin registration surface
    └─ #54 PERF            ──► [x] experimental detector batch + canary closed (no pack promotion)
```

---

## Phase index (rollup)

### CWE — #52

- [x] Detail checklist complete in [python-heuristics-cwe.md](./python-heuristics-cwe.md) (docs-only plan, 2026-07-31)
- [x] Priority CWE batch green on fixtures  
  **Evidence:** `go test ./internal/lang/python/detectors/cwe/` hit/miss for CWE-22/78/79/89/502; fixtures under `tests/fixtures/python/cwe/` (2026-07-31, `feat/python-cwe-heuristics`).
- [x] Validation: `make lint` + `make test` recorded on implement branch  
  **Evidence:** `make lint` / `make test` green; `CGO_ENABLED=0 go build` OK (2026-07-31, `feat/python-cwe-heuristics`).
- [x] **Remaining CWE batch plans** under [cwe-plans/](./cwe-plans/) after chunk scan (2026-07-31, post-PR #67 on `main`)  
  **Evidence:** `cwe-plans/README.md` + `_inventory.json` — 344 partition (5 shipped / 154 implement / 185 deferred); waves P0–P3.  
  **Testing triad:** unit tests + `tests/fixtures/python/cwe/CWE-N-{vulnerable,safe}.txt` (BP-style, parallel to `python/bp/`) + `make integration-python` / `TestPythonCWEFixturesMatrix`.
- [ ] Execute cwe-plans implement waves (start **P0** batches 01–04; each ID ships unit + fixtures + matrix)

### BP — #53

- [x] Detail checklist complete in [python-heuristics-bp.md](./python-heuristics-bp.md) (docs-only plan, 2026-07-31)
- [x] Priority `BP-PY-*` batch green (A+B + C high-signal; inline unit tests; D/E deferred)
- [x] Validation: `make lint` + `make test` green on `feat/python-bp-heuristics` (2026-07-31)

### PERF — #54 (complete — experimental; canary closed)

- [x] PERF ledger updated in [python-heuristics-perf.md](./python-heuristics-perf.md) (catalogue phase, 2026-07-31)
- [x] Python PERF catalogue seed expanded (`PERF-PY-1` … `PERF-PY-30`)  
  **Evidence:** `ruleset/python/chunks/perf-py-*.json`; `go test ./ruleset/python/` passes.  
  **Evidence:** `internal/lang/python/detectors/perf/`; paired fixtures under `tests/fixtures/python/perf/`.
- [x] **Precision hardening** — [PYTHON-PRECISION-HARDENING-CHECKLIST.md](./PYTHON-PRECISION-HARDENING-CHECKLIST.md) Phases 0–4 complete (2026-08-01).
- [x] **Reference-corpus canary + maturity** — [pref-plans/PERF-PY-CANARY-2026-08-01.md](./pref-plans/PERF-PY-CANARY-2026-08-01.md); stay experimental; no recommended/perf promotion.

### Epic closure

- [x] #52 and #53 have shipped initial heuristic batches on `chore/epic-51-integration`  
  **Evidence:** CWE-22/78/79/89/502 + BP-PY A/B + C high-signal; remaining CWE expansion deferred.  
  **Validation (integration):** `make lint` + `make test` green 2026-07-31.
- [x] #54 experimental detector batch + canary/maturity closure (2026-08-01)
- [ ] Epic #51 closed when CWE expansion waves finish (or left open for remaining expansion)

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
| #54 | python(perf): seed PERF + heuristics | [python-heuristics-perf.md](./python-heuristics-perf.md) (complete — experimental + canary 2026-08-01) |

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
