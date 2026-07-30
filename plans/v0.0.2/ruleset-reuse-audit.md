# Ruleset reuse audit — Go catalogues → Python

> **Parent:** `plans/v0.0.2/phase-4-rulesets.md` / epic #39  
> **Status:** draft seed (scan agent should expand with full evidence)  
> **Issue:** #43  

---

## Purpose

Decide which existing `ruleset/golang/**` entries are **portable shells** vs **Go-only**, so Python catalogues reuse structure without copying framework noise.

---

## Inventory (baseline on main)

### Layout

| Path | Role |
|------|------|
| `ruleset/golang/bad-practices.json` | BP-* metadata (~135 rules) |
| `ruleset/golang/chunks/cwe-*.json` | CWE catalogue chunks |
| `ruleset/golang/chunks/perf-*.json` | PERF catalogue chunks |
| `ruleset/golang/go_module_advisories.csv` | Go module advisories (Go-only) |

### Schema — CWE/PERF chunk entry (sample keys)

`id`, `name`, `description`, `original_description`, `detection_notes`, `category`, `status`, `weakness_abstraction`, `go_relevance`, `applicable_to`

### Schema — BP entry

`id`, `name`, `description`, `detection_notes`, `severity`, `category`  
(No language field today.)

### `applicable_to` distribution (golang chunks)

Observed tags (all Go ecosystem): `golang`, `go-web`, `gin`, `sqlx`, `gorm`, `echo`, `fiber`, `grpc`, `protobuf`, `go-redis`, `cli`, `prometheus`, `cobra`.

- **Rules with `python` in `applicable_to`:** **0**  
- Status mix (approx): Draft / Implemented / Incomplete / Stable  

### Go code consumption

- BP: generated into `internal/lang/go/detectors/bad_practices/metadata_gen.go` from JSON  
- PERF/CWE: detectors + generated metadata under `internal/lang/go/detectors/**`  
- Packs: `internal/rules/pack.go` classifies by ID prefix (`BP-`, `PERF-`, `CWE-`)

---

## Classification (working policy)

| Class | Definition | Reuse for Python? |
|-------|------------|-------------------|
| **A — Portable CWE shell** | Weakness is language-neutral; only detection notes / sinks differ | **Yes** — copy ID/name/original_description; rewrite `detection_notes`, set `applicable_to` to include `python`, drop or replace `go_relevance` |
| **B — Go framework CWE** | Notes/APIs tied to gin/echo/gorm/sqlx/etc. | **Partial** — keep CWE id if weakness is generic; strip framework tags; rewrite notes |
| **C — Go PERF** | Allocations, goroutines, net/http, etc. | **No** as-is — separate Python PERF catalogue later |
| **D — Go BP** | go.mod, goroutines, `t.Helper`, error wrapping `%w` | **No** as-is — Python BP needs idioms (exceptions, venv, typing) |
| **E — Advisories CSV** | Go modules | **No** |

### High-value portable CWE shells (seed candidates)

Prefer first Python seeds from taint-relevant CWEs already special-cased in Go:

- CWE-22 Path traversal  
- CWE-78 OS command injection  
- CWE-79 XSS  
- CWE-89 SQL injection  

Also consider other language-neutral bases once seeds land: CWE-20, CWE-90, CWE-94, CWE-116, CWE-200, CWE-287, CWE-352, CWE-502, CWE-611, CWE-918 (verify per entry).

### Explicitly non-portable examples

- Most `PERF-*` (Go runtime / stdlib)  
- BP concurrency / testing / go.mod hygiene  
- Rules whose `applicable_to` is only gin/gorm-specific *and* description is framework-bound  

---

## Recommended layout

```text
ruleset/
  golang/          # existing; unchanged ownership for Go plugin
  python/          # WIP Python catalogues (this phase)
    README.md
    bad-practices.json
    chunks/
      cwe-seed.json   # portable CWE shells only
  shared/          # OPTIONAL later: pure CWE ID→name map without language notes
```

**Reuse mechanism:** same JSON object shape; per-language `detection_notes` + `applicable_to`. Do not invent a second schema.

---

## Scan agent checklist (expand this file)

- [ ] Exact rule counts per chunk file  
- [ ] Table of top 20 portable CWE IDs with rationale  
- [ ] List of BP categories with portability note  
- [ ] Confirm no existing `ruleset/python` on main before create  
- [ ] Note any generator scripts that hardcode `ruleset/golang` paths  

---

## References

- `ruleset/golang/chunks/cwe-001-050.json` (CWE-22 sample)  
- `internal/lang/go/detectors/bad_practices/metadata_gen.go`  
- `documents/rule-catalog-and-maturity.md`  
- Issue #43  
