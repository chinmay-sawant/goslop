# Python ruleset catalogues (WIP)

> **Status:** seed only — **no** Python AST detectors in this tree  
> **Issue:** #43 · epic #39 · plan: `plans/v0.0.2/phase-4-rulesets.md`  
> **Audit:** `plans/v0.0.2/ruleset-reuse-audit.md`

## Layout

```text
ruleset/python/
  README.md              # this file
  bad-practices.json     # empty object for now (do not bulk-copy Go BP)
  chunks/
    cwe-seed.json        # portable CWE shells (CWE-22, 78, 79, 89)
```

Sibling Go catalogues remain the default product source of truth:

- `ruleset/golang/bad-practices.json`
- `ruleset/golang/chunks/*.json`

## Reuse policy

1. **Same JSON shape** as golang CWE/PERF chunks and BP maps. Do not invent a second schema.
2. **Portable CWE shells (class A):** keep `id`, `name`, `original_description`, `category`, `weakness_abstraction`. Rewrite `description` and `detection_notes` for Python sinks. Set `applicable_to` to include `"python"`.
3. **Do not bulk-copy** Go PERF (`PERF-*`) or Go bad practices (`BP-*`). Those are runtime- and idiom-specific (see audit classes C/D).
4. **Do not** point Go generators (`metadata_gen.go`, `metadata_generated.go`) at this directory. Go loaders stay on `ruleset/golang/**`.
5. `go_relevance` is retained for **schema parity** with golang chunks; it is a historical field name. Prefer language-neutral severity of the weakness, not “only relevant to Go.”

## Seeds in this PR

| File | Contents |
|------|----------|
| `chunks/cwe-seed.json` | CWE-22, CWE-78, CWE-79, CWE-89 with `applicable_to: ["python"]` and Python-oriented `detection_notes` |
| `bad-practices.json` | `{}` — Python BP idioms (exceptions, typing, packaging) are a later phase |

## What this is not

- Not a claim of detector implementation or CLI Python scanning  
- Not catalogue parity with the full 175 CWE / 242 PERF golang JSON sets  
- Not consumed by `goslop --list-rules` until a Python plugin embeds/loads these files  

## Validation

JSON must parse as objects. Optional: `go test ./ruleset/python/`.
