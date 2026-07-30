## Context

Parent epic: Python heuristic detectors (CWE / BP / PERF).

Python bad-practice **catalogue** exists at `ruleset/python/bad-practices.json` (**50** rules, IDs `BP-PY-*`) targeting core Python plus Flask, Django, FastAPI/Starlette, Jinja2, SQLAlchemy, and production hygiene. Detection notes describe sinks; **no detectors** implement them yet.

This issue implements **bad-practice heuristic logic** for those catalogue entries (and registers them on the Python plugin).

## Scope (in)

1. Consume `ruleset/python/bad-practices.json` metadata for rule ID / severity / name (embed or load; do not reuse Go `metadata_gen.go` paths)
2. Implement heuristics for a **priority subset**, then expand:
   - Core: bare except, mutable defaults, `assert` for validation, `open` without `with`
   - Security hygiene: `shell=True`, pickle/yaml.load/eval, hardcoded secrets
   - Framework: Flask DEBUG/SECRET_KEY/`send_file`; Django DEBUG/raw SQL/`mark_safe`/`csrf_exempt`; FastAPI blocking I/O in `async def`
   - Templates / DB: Jinja2 autoescape off; SQLAlchemy `text()` f-strings
3. Register BP detectors on `internal/lang/python` session catalogue
4. Fixtures for hit/miss per implemented `BP-PY-*`
5. Keep IDs as `BP-PY-*` (do not collide with Go `BP-*`)

## Out of scope

- CWE and PERF families (sibling issues)
- Porting all 135 Go BP rules
- Full framework type inference

## Success criteria

- [ ] Priority BP heuristics implemented and registered
- [ ] Fixtures prove true positives for each shipped rule
- [ ] `languages = ["python"]` surfaces `BP-PY-*` findings
- [ ] Severity/category match catalogue metadata
- [ ] `make lint` + `make test` green

## Plan

- Catalogue: `ruleset/python/bad-practices.json`
- README: `ruleset/python/README.md`
- Go reference: `internal/lang/go/detectors/bad_practices/`
- Plugin: `internal/lang/python/`

## References

- Parent: #51
- Continues from #39 / PR #50
- Relates to sibling CWE and PERF heuristic issues
