## Context

Parent epic: Python heuristic detectors (CWE / BP / PERF).

Python CWE **catalogue JSON** exists under `ruleset/python/chunks/cwe-*.json` (~344 rules from `699.csv`, field `python_relevance`, `detection_notes`). The Python plugin is still source-only with **zero detectors**.

This issue is to implement **CWE heuristic detection logic** for Python: turn catalogue / detection notes (and any related text chunk material) into scanners that emit findings with stable rule IDs (`CWE-*`).

## Scope (in)

1. Load or embed Python CWE metadata needed for registration (mirror Go pattern where useful; do not point Go generators at `ruleset/python`)
2. Implement an initial **high-signal CWE batch** first (minimum target set):
   - CWE-22 path traversal  
   - CWE-78 OS command injection  
   - CWE-79 XSS  
   - CWE-89 SQL injection  
   - CWE-502 deserialization (pickle/yaml unsafe load)  
   Expand from `detection_notes` for additional high `python_relevance` IDs as capacity allows
3. Wire detectors into `internal/lang/python` (`Detectors` / `NewDetectors`) so enabled-language scans report findings
4. Fixtures under `tests/fixtures/python/` (or sibling) for hit/miss cases per implemented rule
5. Prefer pure-Go / pattern / light AST approaches; document parser choice
6. Ensure findings carry language Python, correct severity/pack classification

## Out of scope

- Full 344-rule CWE implementation in this ticket (batch + framework for more)
- Bad-practice and PERF families (sibling issues)
- Changing Go CWE detectors

## Success criteria

- [ ] At least the priority CWE batch above is registered and can fire on fixtures
- [ ] `languages = ["python"]` scan reports CWE findings for fixture tree
- [ ] Tests cover each implemented rule (positive + at least one negative)
- [ ] Catalogue IDs stay aligned with `ruleset/python/chunks` for implemented rules
- [ ] `make lint` + `make test` green

## Plan

- Catalogue: `ruleset/python/chunks/cwe-*.json`
- Mapping: `plans/v0.0.2/python-cwe-from-699-mapping.md`
- Go reference: `internal/lang/go/detectors/cwe/`
- Plugin: `internal/lang/python/`

## References

- Parent: #51
- Continues from #39 / PR #50
- Relates to sibling BP and PERF heuristic issues
