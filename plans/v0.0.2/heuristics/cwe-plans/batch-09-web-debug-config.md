# Batch 09 — Web Debug Config

> **Parent:** `plans/v0.0.2/heuristics/cwe-plans/README.md` — v0.0.2 Python CWE remaining heuristics
> **Epic / issue:** [#51](https://github.com/chinmay-sawant/goslop/issues/51) · [#52](https://github.com/chinmay-sawant/goslop/issues/52) expansion
> **Status:** **complete** — implemented and validated 2026-07-31
> **Validation evidence:** fixture-only per-rule tests, Python CWE fixture matrix, `gofmt`, `make lint`, `make test`, and `git diff --check` passed.
> **Wave:** P2
> **IDs (8):** CWE-756, CWE-489, CWE-15, CWE-1051, CWE-1052, CWE-1125, CWE-1188, CWE-921
> **PR policy:** one PR for this batch only — do not mix other wave IDs

---

## Architecture constraints

| Constraint | Detail |
|------------|--------|
| Package | `internal/lang/python/detectors/cwe/` only |
| Registration | `RegisterRule("CWE-*", detect…, &Meta…, gates…)` from `init()` |
| Scan | Existing `PyCweScan` — no new detector type |
| Detection | Pure-Go source patterns + needles / `PyCweFacts` |
| Language | `LanguagePython` gate already in scan |
| Plugin | **Do NOT invent a second plugin.** `detectors.All()` already wires CWE scan |
| IDs | Always `CWE-*`; metadata aligned with chunk catalogue |
| **File size policy** | Target max **1500** / hard max **2000** per Go file |
| Target file(s) | `rules_web_config.go` |
| Validation | `make lint` + `make test` unchecked until green |
| Fixtures (text) | **Required:** `tests/fixtures/python/cwe/CWE-N-{vulnerable,safe}.txt` (BP-style; see Testing section) |
| Unit tests | **Required:** hit/miss in `internal/lang/python/detectors/cwe/*_test.go` |
| Integration | **Required:** `tests/integration/python/cwe_matrix_test.go` / `make integration-python` |

## Overview

DEBUG flags, insecure framework security settings, hardcoded config surface

| Rule | Name | Tier | Category | Relevance | Chunk |
|------|------|------|----------|-----------|-------|
| `CWE-756` | Missing Custom Error Page | A | Error Handling | Medium | `cwe-751-800.json` |
| `CWE-489` | Active Debug Code | A | General | Medium | `cwe-451-500.json` |
| `CWE-15` | External Control of System or Configuration Setting | B | Configuration | Medium | `cwe-001-050.json` |
| `CWE-1051` | Initialization with Hard-Coded Network Resource Configuration Data | A | Resource Management | Medium | `cwe-1051-1100.json` |
| `CWE-1052` | Excessive Use of Hard-Coded Literals in Initialization | A | Configuration | Medium | `cwe-1051-1100.json` |
| `CWE-1125` | Excessive Attack Surface | A | General | Medium | `cwe-1101-1150.json` |
| `CWE-1188` | Initialization of a Resource with an Insecure Default | A | Resource Management | Medium | `cwe-1151-1200.json` |
| `CWE-921` | Storage of Sensitive Data in a Mechanism without Access Control | A | Information Disclosure | Medium | `cwe-901-950.json` |

## Executive Summary

Ship **8** remaining CWE heuristics in this theme. Prefer high-signal, low-FP patterns; same-file heuristics only (no interproc taint). Extend registration catalogue tests and fixture matrix.

## Phase 0: Placement + budget

- [x] Record baseline: `wc -l internal/lang/python/detectors/cwe/*.go`
- [x] Ensure target `rules_web_config.go` exists (create if needed); keep ≤1500 lines projected
- [x] Prefer domain file over growing `rules.go` past soft cap
- [x] Append FN-safe needles to `needles.go` / per-rule gates
- [x] No changes to `detectors/all.go` unless constructor rename (should be none)

## Phase 1: `CWE-756` — Missing Custom Error Page

> Catalogue: `ruleset/python/chunks/cwe-751-800.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE756` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-756", detectCWE756, &Meta…, gates...)` in `rules_web_config.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-756`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `DEBUG\s*=\s*True`, `app.debug\s*=\s*True`, `app.config\[\"DEBUG\"\]\s*=\s*True`, `Flask(__name__).run(debug=True)`, `django.views.debug`
- [x] Implement `detectCWE756` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-756`
- [x] Unit miss: safe pattern → no `CWE-756`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-756-vulnerable.txt` + `CWE-756-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 2: `CWE-489` — Active Debug Code

> Catalogue: `ruleset/python/chunks/cwe-451-500.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE489` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-489", detectCWE489, &Meta…, gates...)` in `rules_web_config.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-489`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `DEBUG = True`, `app.run(debug=True)`, `Flask(__name__, debug=True)`, `django settings DEBUG True`, `pdb.set_trace()`, `breakpoint() in non-test`
- [x] Implement `detectCWE489` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-489`
- [x] Unit miss: safe pattern → no `CWE-489`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-489-vulnerable.txt` + `CWE-489-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 3: `CWE-15` — External Control of System or Configuration Setting

> Catalogue: `ruleset/python/chunks/cwe-001-050.json` · tier **B** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE15` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-15", detectCWE15, &Meta…, gates...)` in `rules_web_config.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-15`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `os.environ[`, `os.putenv`, `settings.`, `config[`, `django.conf.settings`, `request.* → conf`
- [x] Implement `detectCWE15` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-15`
- [x] Unit miss: safe pattern → no `CWE-15`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-15-vulnerable.txt` + `CWE-15-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 4: `CWE-1051` — Initialization with Hard-Coded Network Resource Configuration Data

> Catalogue: `ruleset/python/chunks/cwe-1051-1100.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE1051` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-1051", detectCWE1051, &Meta…, gates...)` in `rules_web_config.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-1051`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `http://IP:port literals`, `hardcoded internal hosts in requests/urlopen`
- [x] Implement `detectCWE1051` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-1051`
- [x] Unit miss: safe pattern → no `CWE-1051`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-1051-vulnerable.txt` + `CWE-1051-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 5: `CWE-1052` — Excessive Use of Hard-Coded Literals in Initialization

> Catalogue: `ruleset/python/chunks/cwe-1051-1100.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE1052` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-1052", detectCWE1052, &Meta…, gates...)` in `rules_web_config.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-1052`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `sqlalchemy/create_engine DSN with password=`, `psycopg2.connect password= literal`, `SECRET_KEY = '...'`
- [x] Implement `detectCWE1052` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-1052`
- [x] Unit miss: safe pattern → no `CWE-1052`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-1052-vulnerable.txt` + `CWE-1052-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 6: `CWE-1125` — Excessive Attack Surface

> Catalogue: `ruleset/python/chunks/cwe-1101-1150.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE1125` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-1125", detectCWE1125, &Meta…, gates...)` in `rules_web_config.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-1125`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `/debug`, `/admin`, `django.contrib.admin + DEBUG`, `Flask debug toolbar routes`, `werkzeug debugger`, `app.run(debug=True)`
- [x] Implement `detectCWE1125` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-1125`
- [x] Unit miss: safe pattern → no `CWE-1125`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-1125-vulnerable.txt` + `CWE-1125-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 7: `CWE-1188` — Initialization of a Resource with an Insecure Default

> Catalogue: `ruleset/python/chunks/cwe-1151-1200.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE1188` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-1188", detectCWE1188, &Meta…, gates...)` in `rules_web_config.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-1188`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `DEBUG=True in prod settings`, `ALLOWED_HOSTS=['*']`, `SECRET_KEY default`, `ssl_context check_hostname=False`, `verify=False`
- [x] Implement `detectCWE1188` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-1188`
- [x] Unit miss: safe pattern → no `CWE-1188`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-1188-vulnerable.txt` + `CWE-1188-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 8: `CWE-921` — Storage of Sensitive Data in a Mechanism without Access Control

> Catalogue: `ruleset/python/chunks/cwe-901-950.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE921` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-921", detectCWE921, &Meta…, gates...)` in `rules_web_config.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-921`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `open('/tmp/....key|pem|token|secret')`, `NamedTemporaryFile delete=False with secrets`, `world-writable path literals`
- [x] Implement `detectCWE921` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-921`
- [x] Unit miss: safe pattern → no `CWE-921`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-921-vulnerable.txt` + `CWE-921-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 9: Batch validation

- [x] `gofmt -w` on touched files
- [x] `make lint`
- [x] `make test`
- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` (individual unit tests)
- [x] `go test ./tests/integration/python/ -count=1` **or** `make integration-python` (CWE matrix over all `python/cwe` pairs)
- [x] Confirm every batch ID has both `CWE-N-vulnerable.txt` and `CWE-N-safe.txt` under `tests/fixtures/python/cwe/`
- [x] Fixture count: `DiscoverPythonCWECases` includes all new IDs (pair discovery, not a manual allowlist)
- [x] Update `_inventory.json`: move batch IDs from `missing` → `implemented`
- [x] Update this ledger statuses to `[x]` with evidence
- [x] Package files still ≤2000 lines (split if not)


## Testing requirements (fixtures + unit + integration)

> Same bar as BP-PY: detector alone is **not** enough. Each owned `CWE-N` needs the triad below.

### Fixture text files (required)

| Path | Purpose |
|------|---------|
| `tests/fixtures/python/cwe/CWE-N-vulnerable.txt` | Hit corpus — scan must emit `CWE-N` |
| `tests/fixtures/python/cwe/CWE-N-safe.txt` | Miss corpus — scan must **not** emit `CWE-N` |

**Format** (see `tests/fixtures/README.md` and BP examples under `tests/fixtures/python/bp/`):

```text
# Fixture for CWE-N (vulnerable|safe)
lang: python
file: CWE-N-vulnerable.py   # or -safe.py
---
# python source body
```

- Keep fixtures under **`python/cwe/`** only (parallel to **`python/bp/`** for BP-PY).
- Do **not** commit `.py` sources; discovery is pair-based (`DiscoverPythonCWECases`).

### Individual unit tests (required)

| Path | Purpose |
|------|---------|
| `internal/lang/python/detectors/cwe/rules_test.go` (or domain `rules_*_test.go` / `scan_test.go`) | Hit + miss table tests per ID |
| Registration catalogue test | Want-list includes every new `CWE-N` from this batch |

```sh
go test ./internal/lang/python/detectors/cwe/ -count=1
```

### Integration matrix (required)

| Path | Purpose |
|------|---------|
| `tests/integration/python/cwe_matrix_test.go` | `TestPythonCWEFixturesMatrix` — scans every discovered pair |
| Helpers | `integration.DiscoverPythonCWECases`, `PythonCWEFixtureRel` in `tests/integration/discover.go` |

```sh
go test ./tests/integration/python/ -count=1
# or
make integration-python
```

BP analogue (do not mix): `tests/fixtures/python/bp/` + `bp_matrix_test.go` + `DiscoverPythonBPCases`.

### Per-ID checklist (repeat for every rule in this batch)

- [x] Unit hit/miss for `CWE-N`
- [x] `tests/fixtures/python/cwe/CWE-N-vulnerable.txt`
- [x] `tests/fixtures/python/cwe/CWE-N-safe.txt`
- [x] Matrix auto-discovers pair; vulnerable asserts finding; safe asserts absence
- [x] `make lint` + `make test` + `make integration-python` green before PR merge

## Dependencies

| Depends on | Note |
|------------|------|
| batch-00 | Framework + priority rules already present |
| Catalogue chunks | IDs must exist in `ruleset/python/chunks/` |
| Parent README | ownership + PR policy |
