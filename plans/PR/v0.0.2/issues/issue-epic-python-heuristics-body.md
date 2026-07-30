## Context

Epic **#39** landed the Python language foundation: config `languages`, source-only `LanguagePlugin` stub, and catalogue JSON under `ruleset/python/` (CWE chunks from `699.csv`, `BP-PY-*` bad practices). Those catalogues are **metadata only** — no Python detectors run yet.

Go already implements three detector families with heuristics driven by rule metadata / AST walkers:

| Family | Go home | Python catalogue today |
|--------|---------|------------------------|
| **CWE** | `internal/lang/go/detectors/cwe/` | `ruleset/python/chunks/cwe-*.json` (~344 rules) |
| **Bad practices** | `internal/lang/go/detectors/bad_practices/` | `ruleset/python/bad-practices.json` (50 `BP-PY-*`) |
| **PERF** | `internal/lang/go/detectors/perf/` | **missing** — need Python PERF seeds + heuristics |

This epic tracks **implementing heuristic detectors** for Python: encode `detection_notes` / catalogue text into real scan logic, register under the Python plugin, and cover CWE, BP, and PERF.

## Goal

Ship a working Python heuristic surface (incremental OK) so `languages = ["python"]` (or `["go","python"]`) produces real findings — not only empty walks.

## Sub-issues

- [ ] #52 — python(cwe): implement CWE heuristic detectors from catalogue
- [ ] #53 — python(bp): implement BP-PY bad-practice heuristics from catalogue
- [ ] #54 — python(perf): seed PERF catalogue and implement performance heuristics

## Out of scope (epic-level)

- Full parity with all Go PERF/CWE/BP detectors in one PR
- Tree-sitter / CGO Python parsers (prefer pure-Go or source-pattern heuristics first)
- Re-opening #39 (foundation already shipped in #50)

## Success criteria

- [ ] Three child issues linked as sub-issues of this epic
- [ ] Each family has detectors registered on `internal/lang/python`
- [ ] Fixture-backed heuristics for initial batches per family
- [ ] Default remains Go-only; Python stays opt-in via `languages`

## Plan

- Foundation: PR #50 / epic #39
- Catalogues: `ruleset/python/`
- Mapping: `plans/v0.0.2/python-cwe-from-699-mapping.md`
- Seam: `internal/core/plugin.go`, `internal/lang/python/`
- Bodies: `plans/PR/v0.0.2/issue-python-*-heuristics-body.md`

## References

- Continues from #39 (closed)
- PRs: #50
- Docs: `ruleset/python/README.md`, `documents/architecture-performance.md`
