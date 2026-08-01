# Python multi-file corpus (slim)

In-repo multi-module snippets for precision/regression checks before promoting
`PERF-PY-*` / blocking `BP-PY-*` maturity.

| App | Path | Intent |
|-----|------|--------|
| Django ORM | `django_orm/` | N+1, claim/lock, index covering |
| FastAPI | `fastapi_app/` | middleware, async sleep, globals |
| Flask | `flask_app/` | `send_file` vs `send_from_directory` |
| SQLAlchemy | `sqlalchemy_app/` | lazy relation / session patterns |

Expected findings live in `expected.json` beside each app when added.
External canary apps remain the primary FP triage source until this corpus expands.
