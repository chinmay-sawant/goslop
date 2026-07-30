# Ruleset reuse audit — Go catalogues → Python

> **Parent:** `plans/v0.0.2/phase-4-rulesets.md` / epic #39  
> **Status:** evidence complete (Phase 4 / #43)  
> **Issue:** #43  
> **Branch:** `feat/python-rulesets`  
> **Inventory date:** 2026-07-30 on `main` + this PR’s Python seeds  

---

## Purpose

Decide which existing `ruleset/golang/**` entries are **portable shells** vs **Go-only**, so Python catalogues reuse structure without copying framework noise.

---

## Inventory (baseline on main)

### Layout

| Path | Role | Count (scripted) |
|------|------|-----------------:|
| `ruleset/golang/bad-practices.json` | BP-* metadata | **135** |
| `ruleset/golang/chunks/cwe-*.json` | CWE catalogue chunks | **175** |
| `ruleset/golang/chunks/perf-*.json` | PERF catalogue chunks | **242** |
| `ruleset/golang/go_module_advisories.csv` | Go module advisories (Go-only) | 1 data row |
| `ruleset/python/**` on main | — | **absent** (confirmed before seed create) |

### Exact rule counts per chunk file

Script: load each `ruleset/golang/chunks/*.json` and `len(object)`.

| File | Rules |
|------|------:|
| `cwe-001-050.json` | 3 |
| `cwe-051-100.json` | 8 |
| `cwe-101-150.json` | 2 |
| `cwe-151-200.json` | 5 |
| `cwe-201-9999.json` | 157 |
| **CWE total** | **175** |
| `perf-001-050.json` | 50 |
| `perf-051-100.json` | 50 |
| `perf-101-150.json` | 50 |
| `perf-151-200.json` | 50 |
| `perf-201-224.json` | 24 |
| `perf-225-232.json` | 8 |
| `perf-232-241.json` | 10 |
| **PERF total** | **242** |

Note: runtime registration may expose **239** PERF detectors (`documents/rule-catalog-and-maturity.md`); catalogue JSON holds **242** PERF entries. Counts above are **JSON catalogue**, not binary registration.

### Schema — CWE/PERF chunk entry

Every CWE and PERF object uses the same key set:

| Field | Type (observed) | Notes |
|-------|-----------------|-------|
| `id` | number | Numeric CWE/PERF id (e.g. `22`, not `"CWE-22"`) |
| `name` | string | Title |
| `description` | string | Often Go/Gin/GORM-flavoured prose |
| `original_description` | string | Closer to MITRE / language-neutral text |
| `detection_notes` | string | Today often a **shared boilerplate** string about taint / `.Raw(` / `fmt.Sprintf` |
| `category` | string | e.g. Injection, Path Traversal |
| `status` | string | `Draft` / `Implemented` / `Incomplete` / `Stable` |
| `weakness_abstraction` | string | e.g. Base, Variant |
| `go_relevance` | string | `High` / `Medium` / `Low` |
| `applicable_to` | string[] | Ecosystem tags |

Object map key form: `"CWE-22"`, `"PERF-101"`, etc.

### Schema — BP entry

| Field | Type | Notes |
|-------|------|-------|
| `id` | string | e.g. `"BP-1"` |
| `name` | string | |
| `description` | string | Go-idiomatic |
| `detection_notes` | string | AST / Go patterns |
| `severity` | string | `info` / `low` / `medium` / `high` |
| `category` | string | See BP categories below |

**No** `applicable_to`, **no** `go_relevance`, **no** language field on BP today.

### `applicable_to` distribution (golang chunks)

**Scripted tag occurrence counts** (a rule with N tags increments N counters):

| Tag | Occurrences (CWE+PERF) |
|-----|-----------------------:|
| `golang` | 417 |
| `go-web` | 319 |
| `gin` | 264 |
| `sqlx` | 252 |
| `gorm` | 245 |
| `echo` | 24 |
| `fiber` | 21 |
| `grpc` | 2 |
| `protobuf` | 2 |
| `go-redis` | 2 |
| `cli` | 2 |
| `prometheus` | 1 |
| `cobra` | 1 |
| **`python`** | **0** |

- **Rules with `python` in `applicable_to`:** **0** (confirmed on main; all 175 CWE + 242 PERF).
- **CWE tag sets:** **all 175** CWE entries share exactly  
  `["golang", "gin", "gorm", "sqlx", "go-web"]` (order may vary in file; set is identical).
- **PERF tag sets (top):** 96 `["golang"]` only; 55 full gin/gorm/sqlx/go-web set; remainder mixed go-web / echo / fiber combinations. **No** non-Go tags.

### `go_relevance` distribution

| Value | CWE | PERF | Total |
|-------|----:|-----:|------:|
| High | 9 | 57 | 66 |
| Medium | 166 | 141 | 307 |
| Low | 0 | 44 | 44 |

### `status` distribution

| Status | CWE | PERF | Total |
|--------|----:|-----:|------:|
| Draft | 93 | 118 | 211 |
| Implemented | 0 | 124 | 124 |
| Incomplete | 77 | 0 | 77 |
| Stable | 5 | 0 | 5 |

**Stable CWE (all High or Medium relevance):** CWE-22, CWE-78, CWE-79, CWE-89, CWE-426.

### BP severity + categories

| Severity | Count |
|----------|------:|
| low | 69 |
| medium | 43 |
| high | 12 |
| info | 11 |
| **total** | **135** |

| Category | Count | Portability note |
|----------|------:|------------------|
| Concurrency | 20 | **Go-only** (goroutines, WaitGroup, mutex-by-value) |
| HTTP Frameworks | 15 | **Go-only** (net/http, Gin, Echo, Fiber idioms) |
| Database | 14 | **Mostly Go-only** (GORM/sqlx); idea of raw SQL is portable as CWE-89 shell, not BP text |
| Testing | 13 | **Go-only** (`t.Helper`, `t.Parallel`, `_test.go`) |
| Core Language | 12 | **Go-only** (interfaces, defer, zero values) |
| API Design | 11 | **Go-only** context/API shape |
| Code Organization | 10 | **Go-only** package layout |
| Production Hardening | 10 | **Mostly Go-only**; a few ideas (timeouts) could inspire Python later |
| Dependency Hygiene | 10 | **Go-only** (`go.mod`, module path) |
| Error Handling | 6 | **Go-only** (`_ = err`, `%w` wrapping) |
| Observability | 5 | **Go-only** slog/prometheus patterns today |
| Resource Management | 4 | **Mostly Go-only** (Close err); conceptual “close resource” is portable as new Python BP later |
| Panics | 3 | **Go-only** (`panic`/`recover`) |
| Data Handling | 1 | Review case-by-case; do not bulk-copy |
| Configuration | 1 | Review case-by-case; do not bulk-copy |

**Conclusion:** do **not** bulk-copy Go `bad-practices.json` into Python. Purpose-written **`BP-PY-*`** seeds now live in `ruleset/python/bad-practices.json` (core idioms + Flask/Django/FastAPI/SQLAlchemy/Jinja2/production).

### Go code consumption / hard-coded `ruleset/golang` paths

| Consumer | Path hard-coded? | Notes |
|----------|------------------|-------|
| `internal/lang/go/detectors/bad_practices/metadata_gen.go` | **Yes (comment + embed source)** | Header: `// Code generated from ruleset/golang/bad-practices.json`. Map built in `loadEmbeddedMetadata()`. |
| `internal/lang/go/detectors/cwe/metadata_generated.go` | **Yes (generation source)** | Titles/descriptions from `ruleset/golang/chunks/cwe-*.json` (see package README). |
| `internal/lang/go/detectors/bad_practices/README.md` | Documents `ruleset/golang/bad-practices.json` | |
| `internal/lang/go/detectors/cwe/README.md` | Documents `ruleset/golang/chunks/cwe-*.json` | |
| `internal/rules/pack.go` | ID-prefix only | Classifies `BP-` / `PERF-` / `CWE-`; language-agnostic pack keys |
| Makefile / `scripts/` | **No** live generator checked in | No Makefile target regenerates metadata; generated Go files are committed artifacts |
| Docs (`documents/overview.md`, `perf-rules.md`) | Path references | Claim catalogue root is `ruleset/golang/` (updated for multi-root in Phase 4 docs touch) |

**Implication for Python:** new catalogues under `ruleset/python/` are **not** consumed by current generators. Future Python plugin must load its own tree (or a shared loader); do not point Go `metadata_gen.go` at Python paths.

---

## Classification (working policy + evidence)

| Class | Definition | Count / examples | Reuse for Python? |
|-------|------------|------------------|-------------------|
| **A — Portable CWE shell** | Weakness is language-neutral; MITRE id/name/`original_description` hold; only detection notes / sinks differ | Seed set: **CWE-22, 78, 79, 89** (+ table of 20 candidates below). Many of the 175 CWEs are shells underneath Go-flavoured `description` | **Yes** — copy `id`/`name`/`original_description`/`category`/`weakness_abstraction`; rewrite `description` + `detection_notes`; set `applicable_to` to include `python`; retain `go_relevance` field for schema parity (value may stay High/Medium as “weakness relevance”, not “Go only”) |
| **B — Go framework CWE** | Notes/APIs tied to gin/echo/gorm/sqlx/etc. | **All 175** CWE rows currently tag gin/gorm/sqlx/go-web; ~20 `description` strings name Gin/GORM/sqlx explicitly (e.g. CWE-22/79/89) | **Partial** — keep CWE id if weakness is generic; strip framework tags; rewrite notes (same as A once rewritten) |
| **C — Go PERF** | Allocations, goroutines, net/http, etc. | **242** PERF JSON entries; tags are golang / go-web / frameworks only | **No** as-is — separate Python PERF catalogue later |
| **D — Go BP** | go.mod, goroutines, `t.Helper`, error wrapping `%w` | **135** BP rules; all categories Go-idiomatic | **No** as-is — Python BP needs idioms (exceptions, venv, typing) |
| **E — Advisories CSV** | Go modules | `go_module_advisories.csv` (`github.com/dgrijalva/jwt-go`) | **No** |

### High-value portable CWE shells (seed + backlog)

Prefer first Python seeds from taint-relevant CWEs already special-cased in Go detectors/fixtures:

| CWE ID | Name (short) | Go status | Rationale for Python portability |
|--------|--------------|-----------|----------------------------------|
| **CWE-22** | Path traversal | Stable / High | Language-neutral FS weakness; Python sinks: `open`, `pathlib`, `send_file`, upload joins |
| **CWE-78** | OS command injection | Stable / High | `subprocess(..., shell=True)`, `os.system`, `os.popen` |
| **CWE-79** | XSS | Stable / High | Jinja2/Django autoescape off, `|safe`, `Markup()` |
| **CWE-89** | SQL injection | Stable / High | DB-API f-strings, SQLAlchemy `text()` + format, Django raw |
| CWE-90 | LDAP injection | Draft / Medium | `ldap3` / `python-ldap` filter construction |
| CWE-91 | XML / XPath injection | Draft / Medium | `lxml`, ElementTree query building |
| CWE-93 | CRLF injection | Draft / Medium | Header injection in WSGI/ASGI responses |
| CWE-502 | Unsafe deserialization | Draft / High | `pickle`, `yaml.load` (unsafe), `marshal` |
| CWE-611 | XXE | Draft / Medium | `lxml` / recommend `defusedxml` |
| CWE-918 | SSRF | Incomplete / High | `requests` / `urllib` with user-controlled URL |
| CWE-256 | Plaintext password storage | Incomplete / Medium | Secrets in repo / settings |
| CWE-260 | Password in config file | Incomplete / Medium | `settings.py`, `.env` committed patterns |
| CWE-306 | Missing authentication | Draft / Medium | Unprotected Flask/Django views (future) |
| CWE-338 | Weak PRNG | Draft / Medium | `random` for security tokens |
| CWE-601 | Open redirect | Draft / Medium | `redirect(request.args[...])` |
| CWE-798 | Hard-coded credentials | Draft / Medium | Constants in source |
| CWE-915 | Mass assignment | Incomplete / Medium | Django forms / model construct from request |
| CWE-15 | External control of config | Incomplete / High | Env/config injection into critical settings |
| CWE-426 | Untrusted search path | Stable / Medium | `PATH` / plugin path manipulation |
| CWE-59 | Link following | Draft / Medium | Symlink races on file open |

**Not in current golang CWE catalogue (do not invent IDs in seeds):** CWE-20, CWE-94, CWE-116, CWE-200, CWE-287, CWE-327, CWE-352 — useful later if catalogue expands; verify before adding.

### Explicitly non-portable examples

- Most `PERF-*` (Go runtime / stdlib / goroutines) — e.g. PERF entries about inlining, escape analysis, `net/http` server tuning  
- BP concurrency / testing / go.mod hygiene — e.g. BP-6 WaitGroup Add inside goroutine, BP-3 panic outside main  
- Go module advisories CSV  
- Rules whose **description** is framework-bound *and* left unrewritten (class B until rewritten)

---

## Recommended layout

```text
ruleset/
  golang/          # existing; unchanged ownership for Go plugin
  python/          # WIP Python catalogues (this phase)
    README.md
    bad-practices.json   # BP-PY-* seeds (framework-targeted; no bulk Go BP copy)
    chunks/
      cwe-seed.json      # portable CWE shells only (CWE-22/78/79/89)
  shared/          # OPTIONAL later: pure CWE ID→name map without language notes
```

**Reuse mechanism:** same JSON object shape; per-language `detection_notes` + `applicable_to`. Do not invent a second schema.

**Phase 4 seed deliverable:** 4 portable CWE shells under `ruleset/python/chunks/cwe-seed.json`; empty BP object; README documents policy.

---

## How Go loads BP metadata today (loader note)

1. Source of truth: `ruleset/golang/bad-practices.json`  
2. Generated (committed) code: `internal/lang/go/detectors/bad_practices/metadata_gen.go`  
3. At init: `loadEmbeddedMetadata()` fills `metaByID` with `rules.RuleMetadata` (ID, Title, Description, Severity, Fix from detection_notes, Pack)  
4. Detectors register against that map; packs use ID prefix classification in `internal/rules/pack.go`  
5. CWE parallel path: `metadata_generated.go` from `ruleset/golang/chunks/cwe-*.json`

Python will need an analogous loader under a future `internal/lang/python` plugin — **out of scope** for Phase 4 beyond JSON seeds + schema parity.

---

## Scan agent checklist

- [x] Exact rule counts per chunk file (table above)  
- [x] Table of top 20 portable CWE IDs with rationale  
- [x] List of BP categories with portability note  
- [x] Confirm no existing `ruleset/python` on main before create  
- [x] Note any generator scripts that hardcode `ruleset/golang` paths  
- [x] Inventory schema fields for CWE/PERF chunks and BP JSON  
- [x] Count rules by `applicable_to` / `go_relevance` / `status`  
- [x] Classify portable CWE shells vs Go-framework CWE vs Go-PERF vs Go-BP  
- [x] Confirm **0** existing `python` `applicable_to` tags on main  

---

## References

- `ruleset/golang/chunks/cwe-001-050.json` (CWE-22 sample)  
- `ruleset/python/chunks/cwe-seed.json` (Phase 4 seeds)  
- `internal/lang/go/detectors/bad_practices/metadata_gen.go`  
- `documents/rule-catalog-and-maturity.md`  
- Issue #43 · Relates to epic #39  
