# Batch 06 — Ssrf Redirect

> **Parent:** `plans/v0.0.2/heuristics/cwe-plans/README.md` — v0.0.2 Python CWE remaining heuristics
> **Epic / issue:** [#51](https://github.com/chinmay-sawant/goslop/issues/51) · [#52](https://github.com/chinmay-sawant/goslop/issues/52) expansion
> **Status:** **complete** — implemented and validated 2026-07-31
> **Validation evidence:** `gofmt`, Python CWE package tests, Python CWE fixture matrix, `make lint`, `make test`, and `git diff --check` passed.
> **Wave:** P1
> **IDs (6):** CWE-918, CWE-601, CWE-605, CWE-924, CWE-940, CWE-941
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
| Target file(s) | `rules_ssrf.go` |
| Validation | `make lint` + `make test` unchecked until green |
| Fixtures (text) | **Required:** `tests/fixtures/python/cwe/CWE-N-{vulnerable,safe}.txt` (BP-style; see Testing section) |
| Unit tests | **Required:** hit/miss in `internal/lang/python/detectors/cwe/*_test.go` |
| Integration | **Required:** `tests/integration/python/cwe_matrix_test.go` / `make integration-python` |

## Overview

SSRF, open redirect, channel/host integrity

| Rule | Name | Tier | Category | Relevance | Chunk |
|------|------|------|----------|-----------|-------|
| `CWE-918` | Server-Side Request Forgery (SSRF) | A | General | Medium | `cwe-901-950.json` |
| `CWE-601` | URL Redirection to Untrusted Site ('Open Redirect') | A | Web | Medium | `cwe-601-650.json` |
| `CWE-605` | Multiple Binds to the Same Port | A | General | Medium | `cwe-601-650.json` |
| `CWE-924` | Improper Enforcement of Message Integrity During Transmission in a Communication Channel | A | General | Medium | `cwe-901-950.json` |
| `CWE-940` | Improper Verification of Source of a Communication Channel | A | General | Medium | `cwe-901-950.json` |
| `CWE-941` | Incorrectly Specified Destination in a Communication Channel | A | General | Medium | `cwe-901-950.json` |

## Executive Summary

Ship **6** remaining CWE heuristics in this theme. Prefer high-signal, low-FP patterns; same-file heuristics only (no interproc taint). Extend registration catalogue tests and fixture matrix.

## Phase 0: Placement + budget

- [x] Record baseline: `wc -l internal/lang/python/detectors/cwe/*.go`
- [x] Ensure target `rules_ssrf.go` exists (create if needed); keep ≤1500 lines projected
- [x] Prefer domain file over growing `rules.go` past soft cap
- [x] Append FN-safe needles to `needles.go` / per-rule gates
- [x] No changes to `detectors/all.go` unless constructor rename (should be none)

## Phase 1: `CWE-918` — Server-Side Request Forgery (SSRF)

> Catalogue: `ruleset/python/chunks/cwe-901-950.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE918` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-918", detectCWE918, &Meta…, gates...)` in `rules_ssrf.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-918`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `requests.get/post(url)`, `urllib.request.urlopen(url)`, `httpx.get(url)`, `aiohttp session.get(url)`, `urllib3`
- [x] Implement `detectCWE918` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-918`
- [x] Unit miss: safe pattern → no `CWE-918`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-918-vulnerable.txt` + `CWE-918-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 2: `CWE-601` — URL Redirection to Untrusted Site ('Open Redirect')

> Catalogue: `ruleset/python/chunks/cwe-601-650.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE601` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-601", detectCWE601, &Meta…, gates...)` in `rules_ssrf.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-601`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Manual Static Analysis, Automated Dynamic Analysis, Automated Static Analysis, Automated Static Analysis - Binary or Bytecode, Dynamic Analysis with Automated Results Interpretation, Dynamic Anal…
- Suggested sinks/patterns: `redirect(`, `HttpResponseRedirect(`, `RedirectResponse(`, `flask.redirect`, `django.shortcuts.redirect`, `url_for with _external + user next=`
- [x] Implement `detectCWE601` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-601`
- [x] Unit miss: safe pattern → no `CWE-601`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-601-vulnerable.txt` + `CWE-601-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 3: `CWE-605` — Multiple Binds to the Same Port

> Catalogue: `ruleset/python/chunks/cwe-601-650.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE605` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-605", detectCWE605, &Meta…, gates...)` in `rules_ssrf.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-605`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `SO_REUSEADDR`, `setsockopt(.*SO_REUSEADDR`, `socket.bind((\"0.0.0.0\"`, `bind((\"::\"`, `app.run(host=\"0.0.0.0\")`
- [x] Implement `detectCWE605` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-605`
- [x] Unit miss: safe pattern → no `CWE-605`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-605-vulnerable.txt` + `CWE-605-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 4: `CWE-924` — Improper Enforcement of Message Integrity During Transmission in a Communication Channel

> Catalogue: `ruleset/python/chunks/cwe-901-950.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE924` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-924", detectCWE924, &Meta…, gates...)` in `rules_ssrf.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-924`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `webhook handler decode body without hmac/signature`, `stripe/github webhook without X-*-Signature verify`
- [x] Implement `detectCWE924` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-924`
- [x] Unit miss: safe pattern → no `CWE-924`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-924-vulnerable.txt` + `CWE-924-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 5: `CWE-940` — Improper Verification of Source of a Communication Channel

> Catalogue: `ruleset/python/chunks/cwe-901-950.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE940` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-940", detectCWE940, &Meta…, gates...)` in `rules_ssrf.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-940`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `OAuth callback without state/nonce check`, `request.args user_id trusted as identity`
- [x] Implement `detectCWE940` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-940`
- [x] Unit miss: safe pattern → no `CWE-940`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-940-vulnerable.txt` + `CWE-940-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 6: `CWE-941` — Incorrectly Specified Destination in a Communication Channel

> Catalogue: `ruleset/python/chunks/cwe-901-950.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE941` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-941", detectCWE941, &Meta…, gates...)` in `rules_ssrf.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-941`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `smtp/send_mail to request-controlled address`, `webhook URL from client without allowlist`
- [x] Implement `detectCWE941` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-941`
- [x] Unit miss: safe pattern → no `CWE-941`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-941-vulnerable.txt` + `CWE-941-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 7: Batch validation

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
