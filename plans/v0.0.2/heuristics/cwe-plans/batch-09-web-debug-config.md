# Batch 09 — Web Debug Config

> **Parent:** `plans/v0.0.2/heuristics/cwe-plans/README.md` — v0.0.2 Python CWE remaining heuristics
> **Epic / issue:** [#51](https://github.com/chinmay-sawant/goslop/issues/51) · [#52](https://github.com/chinmay-sawant/goslop/issues/52) expansion
> **Status:** **pending** — not started
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
| Fixtures | `tests/fixtures/python/cwe/CWE-N-{{vulnerable,safe}}.txt` |

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

- [ ] Record baseline: `wc -l internal/lang/python/detectors/cwe/*.go`
- [ ] Ensure target `rules_web_config.go` exists (create if needed); keep ≤1500 lines projected
- [ ] Prefer domain file over growing `rules.go` past soft cap
- [ ] Append FN-safe needles to `needles.go` / per-rule gates
- [ ] No changes to `detectors/all.go` unless constructor rename (should be none)

## Phase 1: `CWE-756` — Missing Custom Error Page

> Catalogue: `ruleset/python/chunks/cwe-751-800.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE756` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-756", detectCWE756, &Meta…, gates...)` in `rules_web_config.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-756`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `DEBUG\s*=\s*True`, `app.debug\s*=\s*True`, `app.config\[\"DEBUG\"\]\s*=\s*True`, `Flask(__name__).run(debug=True)`, `django.views.debug`
- [ ] Implement `detectCWE756` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-756`
- [ ] Unit miss: safe pattern → no `CWE-756`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-756-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 2: `CWE-489` — Active Debug Code

> Catalogue: `ruleset/python/chunks/cwe-451-500.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE489` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-489", detectCWE489, &Meta…, gates...)` in `rules_web_config.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-489`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `DEBUG = True`, `app.run(debug=True)`, `Flask(__name__, debug=True)`, `django settings DEBUG True`, `pdb.set_trace()`, `breakpoint() in non-test`
- [ ] Implement `detectCWE489` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-489`
- [ ] Unit miss: safe pattern → no `CWE-489`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-489-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 3: `CWE-15` — External Control of System or Configuration Setting

> Catalogue: `ruleset/python/chunks/cwe-001-050.json` · tier **B** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE15` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-15", detectCWE15, &Meta…, gates...)` in `rules_web_config.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-15`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `os.environ[`, `os.putenv`, `settings.`, `config[`, `django.conf.settings`, `request.* → conf`
- [ ] Implement `detectCWE15` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-15`
- [ ] Unit miss: safe pattern → no `CWE-15`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-15-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 4: `CWE-1051` — Initialization with Hard-Coded Network Resource Configuration Data

> Catalogue: `ruleset/python/chunks/cwe-1051-1100.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE1051` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-1051", detectCWE1051, &Meta…, gates...)` in `rules_web_config.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-1051`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `http://IP:port literals`, `hardcoded internal hosts in requests/urlopen`
- [ ] Implement `detectCWE1051` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-1051`
- [ ] Unit miss: safe pattern → no `CWE-1051`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-1051-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 5: `CWE-1052` — Excessive Use of Hard-Coded Literals in Initialization

> Catalogue: `ruleset/python/chunks/cwe-1051-1100.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE1052` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-1052", detectCWE1052, &Meta…, gates...)` in `rules_web_config.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-1052`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `sqlalchemy/create_engine DSN with password=`, `psycopg2.connect password= literal`, `SECRET_KEY = '...'`
- [ ] Implement `detectCWE1052` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-1052`
- [ ] Unit miss: safe pattern → no `CWE-1052`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-1052-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 6: `CWE-1125` — Excessive Attack Surface

> Catalogue: `ruleset/python/chunks/cwe-1101-1150.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE1125` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-1125", detectCWE1125, &Meta…, gates...)` in `rules_web_config.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-1125`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `/debug`, `/admin`, `django.contrib.admin + DEBUG`, `Flask debug toolbar routes`, `werkzeug debugger`, `app.run(debug=True)`
- [ ] Implement `detectCWE1125` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-1125`
- [ ] Unit miss: safe pattern → no `CWE-1125`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-1125-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 7: `CWE-1188` — Initialization of a Resource with an Insecure Default

> Catalogue: `ruleset/python/chunks/cwe-1151-1200.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE1188` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-1188", detectCWE1188, &Meta…, gates...)` in `rules_web_config.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-1188`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `DEBUG=True in prod settings`, `ALLOWED_HOSTS=['*']`, `SECRET_KEY default`, `ssl_context check_hostname=False`, `verify=False`
- [ ] Implement `detectCWE1188` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-1188`
- [ ] Unit miss: safe pattern → no `CWE-1188`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-1188-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 8: `CWE-921` — Storage of Sensitive Data in a Mechanism without Access Control

> Catalogue: `ruleset/python/chunks/cwe-901-950.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE921` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-921", detectCWE921, &Meta…, gates...)` in `rules_web_config.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-921`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `open('/tmp/....key|pem|token|secret')`, `NamedTemporaryFile delete=False with secrets`, `world-writable path literals`
- [ ] Implement `detectCWE921` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-921`
- [ ] Unit miss: safe pattern → no `CWE-921`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-921-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 9: Batch validation

- [ ] `gofmt -w` on touched files
- [ ] `make lint`
- [ ] `make test`
- [ ] `go test ./tests/integration/python/ -count=1` (CWE matrix)
- [ ] Update `_inventory.json`: move batch IDs from `missing` → `implemented`
- [ ] Update this ledger statuses to `[x]` with evidence
- [ ] Package files still ≤2000 lines (split if not)

## Dependencies

| Depends on | Note |
|------------|------|
| batch-00 | Framework + priority rules already present |
| Catalogue chunks | IDs must exist in `ruleset/python/chunks/` |
| Parent README | ownership + PR policy |

