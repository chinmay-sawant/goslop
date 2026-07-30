## Context

Parent epic: Python heuristic detectors (CWE / BP / PERF).

Go ships a large **PERF** catalogue (`ruleset/golang/chunks/perf-*.json` + `internal/lang/go/detectors/perf/`). Python has **no PERF JSON tree yet** and no performance detectors.

This issue covers:

1. Seeding a Python PERF catalogue (JSON, ≤50 rules per chunk file, golang-style naming)
2. Implementing **performance heuristics** for Python hot paths (stdlib / web / DB patterns)

Heuristics should be driven by catalogue `detection_notes` (and, where useful, lessons from Go PERF domains) but rewritten for Python runtimes (CPython, asyncio, Django/Flask/FastAPI, DB-API).

## Scope (in)

1. Add `ruleset/python/chunks/perf-*.json` (range-named, **max 50 entries per file**) with `python_relevance`, `applicable_to: ["python"]`
2. Seed an initial PERF set (examples — finalize in implementation):
   - Work in request/async loops (sync I/O in `async def`, `time.sleep` in async)
   - Repeated expensive ops in loops (compile regex in loop, re-open files in loop)
   - N+1-style ORM access patterns (coordinate with BP N+1 where overlap is clear — prefer PERF for cost, BP for style if already listed)
   - Unbounded growth (unbounded list append in hot path, missing pagination)
   - HTTP client without timeout (if not already BP-only; avoid double-firing — document ownership)
3. Implement detectors for the seeded PERF IDs under `internal/lang/python`
4. Fixtures for hit/miss
5. Register on Python plugin; pack classification as performance (`PERF-*` prefix)

## Out of scope

- Porting all ~239 Go PERF rules
- CWE / BP families (siblings) except coordination to avoid duplicate findings
- Micro-benchmark CI gates

## Success criteria

- [ ] Python PERF JSON chunks exist (≤50 rules/file, range names)
- [ ] Initial PERF heuristics registered and fixture-tested
- [ ] `languages = ["python"]` can emit `PERF-*` findings
- [ ] README / mapping notes updated for PERF layout
- [ ] `make lint` + `make test` green

## Plan

- Go reference: `ruleset/golang/chunks/perf-*.json`, `internal/lang/go/detectors/perf/`
- Python layout: `ruleset/python/chunks/` (alongside existing `cwe-*.json`)
- Plugin: `internal/lang/python/`

## References

- Parent: #51
- Continues from #39 / PR #50
- Relates to sibling CWE and BP heuristic issues
