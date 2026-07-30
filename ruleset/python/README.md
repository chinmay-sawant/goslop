# Python ruleset catalogues (WIP)

> **Status:** seed only — **no** Python AST detectors in this tree  
> **Issue:** #43 · epic #39 · plan: `plans/v0.0.2/phase-4-rulesets.md`  
> **Audit:** `plans/v0.0.2/ruleset-reuse-audit.md`

## Layout

```text
ruleset/python/
  README.md              # this file
  bad-practices.json     # BP-PY-* seeds (core + Flask/Django/FastAPI/SQLAlchemy/Jinja2)
  chunks/
    cwe-seed.json        # portable CWE shells (CWE-22, 78, 79, 89)
```

Sibling Go catalogues remain the default product source of truth:

- `ruleset/golang/bad-practices.json`
- `ruleset/golang/chunks/*.json`

## Reuse policy

1. **Same JSON shape** as golang CWE/PERF chunks and BP maps. Do not invent a second schema.
2. **Portable CWE shells (class A):** keep `id`, `name`, `original_description`, `category`, `weakness_abstraction`. Rewrite `description` and `detection_notes` for Python sinks. Set `applicable_to` to include `"python"`.
3. **Do not bulk-copy** Go PERF (`PERF-*`) or Go bad practices (`BP-*`). Those are runtime- and idiom-specific (see audit classes C/D). Python BP IDs use the **`BP-PY-*`** prefix so they never collide with Go `BP-*` when both catalogues are listed.
4. **Do not** point Go generators (`metadata_gen.go`, `metadata_generated.go`) at this directory. Go loaders stay on `ruleset/golang/**`.
5. `go_relevance` is retained for **schema parity** with golang chunks; it is a historical field name. Prefer language-neutral severity of the weakness, not “only relevant to Go.”
6. Optional `applicable_to` on BP entries tags stacks: `python`, `flask`, `django`, `fastapi`, `starlette`, `jinja2`, `sqlalchemy`, `requests`, `httpx`.

## Target frameworks (BP seeds)

| Stack | Example rule themes |
|-------|---------------------|
| **Core Python** | bare except, mutable defaults, `assert` for validation, `open` without `with` |
| **Security hygiene** | `shell=True`, pickle/yaml.load/eval, hardcoded secrets |
| **Flask** | DEBUG/SECRET_KEY, `send_file` user paths, error leakage |
| **Django** | DEBUG/SECRET_KEY/ALLOWED_HOSTS, raw SQL, `mark_safe`, `csrf_exempt`, N+1 |
| **FastAPI / Starlette** | blocking I/O in `async def`, FileResponse paths, response_model leaks |
| **Jinja2** | autoescape off, Markup/\|safe on variables |
| **SQLAlchemy / DB-API** | `text()` f-strings, unclosed sessions, %-format execute |
| **Production** | CORS `*` + credentials, `verify=False`, insecure cookie flags |

## Seeds

| File | Contents |
|------|----------|
| `chunks/cwe-seed.json` | CWE-22, CWE-78, CWE-79, CWE-89 with `applicable_to: ["python"]` and Python-oriented `detection_notes` |
| `bad-practices.json` | **50** `BP-PY-*` catalogue entries (metadata only; detectors not implemented) |

## What this is not

- Not a claim of detector implementation or CLI Python scanning  
- Not catalogue parity with the full 175 CWE / 242 PERF golang JSON sets  
- Not consumed by `goslop --list-rules` until a Python plugin embeds/loads these files  

## Validation

JSON must parse as objects. Optional: `go test ./ruleset/python/`.
