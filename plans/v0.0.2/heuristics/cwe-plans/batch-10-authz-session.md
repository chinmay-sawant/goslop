# Batch 10 — Authz Session

> **Parent:** `plans/v0.0.2/heuristics/cwe-plans/README.md` — v0.0.2 Python CWE remaining heuristics
> **Epic / issue:** [#51](https://github.com/chinmay-sawant/goslop/issues/51) · [#52](https://github.com/chinmay-sawant/goslop/issues/52) expansion
> **Status:** **complete** — implemented and validated 2026-07-31
> **Validation evidence:** fixture-only per-rule tests, Python CWE fixture matrix, `gofmt`, `make lint`, `make test`, and `git diff --check` passed.
> **Wave:** P2
> **IDs (8):** CWE-306, CWE-307, CWE-346, CWE-359, CWE-613, CWE-565, CWE-807, CWE-698
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
| Target file(s) | `rules_auth.go` |
| Validation | `make lint` + `make test` unchecked until green |
| Fixtures (text) | **Required:** `tests/fixtures/python/cwe/CWE-N-{vulnerable,safe}.txt` (BP-style; see Testing section) |
| Unit tests | **Required:** hit/miss in `internal/lang/python/detectors/cwe/*_test.go` |
| Integration | **Required:** `tests/integration/python/cwe_matrix_test.go` / `make integration-python` |

## Overview

Missing auth, origin checks, session hygiene (many B-tier — careful FP)

| Rule | Name | Tier | Category | Relevance | Chunk |
|------|------|------|----------|-----------|-------|
| `CWE-306` | Missing Authentication for Critical Function | B | Security | Medium | `cwe-301-350.json` |
| `CWE-307` | Improper Restriction of Excessive Authentication Attempts | B | Security | Medium | `cwe-301-350.json` |
| `CWE-346` | Origin Validation Error | B | Error Handling | Medium | `cwe-301-350.json` |
| `CWE-359` | Exposure of Private Personal Information to an Unauthorized Actor | B | General | Medium | `cwe-351-400.json` |
| `CWE-613` | Insufficient Session Expiration | B | Security | Medium | `cwe-601-650.json` |
| `CWE-565` | Reliance on Cookies without Validation and Integrity Checking | B | General | Medium | `cwe-551-600.json` |
| `CWE-807` | Reliance on Untrusted Inputs in a Security Decision | B | General | Medium | `cwe-801-850.json` |
| `CWE-698` | Execution After Redirect (EAR) | B | Web | Medium | `cwe-651-700.json` |

## Executive Summary

Ship **8** remaining CWE heuristics in this theme. Prefer high-signal, low-FP patterns; same-file heuristics only (no interproc taint). Extend registration catalogue tests and fixture matrix.

## Phase 0: Placement + budget

- [x] Record baseline: `wc -l internal/lang/python/detectors/cwe/*.go`
- [x] Ensure target `rules_auth.go` exists (create if needed); keep ≤1500 lines projected
- [x] Prefer domain file over growing `rules.go` past soft cap
- [x] Append FN-safe needles to `needles.go` / per-rule gates
- [x] No changes to `detectors/all.go` unless constructor rename (should be none)

## Phase 1: `CWE-306` — Missing Authentication for Critical Function

> Catalogue: `ruleset/python/chunks/cwe-301-350.json` · tier **B** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE306` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-306", detectCWE306, &Meta…, gates...)` in `rules_auth.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-306`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Manual Analysis, Automated Static Analysis, Manual Static Analysis - Binary or Bytecode, Dynamic Analysis with Automated Results Interpretation, Dynamic Analysis with Manual Results Interpretatio…
- Suggested sinks/patterns: `Flask/Django view without login_required/@auth`, `FastAPI route without Depends(security)`
- [x] Implement `detectCWE306` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-306`
- [x] Unit miss: safe pattern → no `CWE-306`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-306-vulnerable.txt` + `CWE-306-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 2: `CWE-307` — Improper Restriction of Excessive Authentication Attempts

> Catalogue: `ruleset/python/chunks/cwe-301-350.json` · tier **B** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE307` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-307", detectCWE307, &Meta…, gates...)` in `rules_auth.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-307`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Dynamic Analysis with Automated Results Interpretation, Dynamic Analysis with Manual Results Interpretation, Manual Static Analysis - Source Code, Automated Static Analysis - Source Code, Automat…
- Suggested sinks/patterns: `login view without rate limit/throttle`, `django-axes missing`
- [x] Implement `detectCWE307` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-307`
- [x] Unit miss: safe pattern → no `CWE-307`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-307-vulnerable.txt` + `CWE-307-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 3: `CWE-346` — Origin Validation Error

> Catalogue: `ruleset/python/chunks/cwe-301-350.json` · tier **B** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE346` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-346", detectCWE346, &Meta…, gates...)` in `rules_auth.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-346`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `Access-Control-Allow-Origin: * with credentials`, `CORS(app, origins='*')`, `csrf_exempt overuse`
- [x] Implement `detectCWE346` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-346`
- [x] Unit miss: safe pattern → no `CWE-346`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-346-vulnerable.txt` + `CWE-346-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 4: `CWE-359` — Exposure of Private Personal Information to an Unauthorized Actor

> Catalogue: `ruleset/python/chunks/cwe-351-400.json` · tier **B** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE359` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-359", detectCWE359, &Meta…, gates...)` in `rules_auth.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-359`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Architecture or Design Review, Automated Static Analysis. Private personal data can enter a program in a variety of ways: Directly from the user in the form of a password or personal information …
- Suggested sinks/patterns: `log ssn/email/password fields`, `print(user.password)`, `response includes full user model`
- [x] Implement `detectCWE359` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-359`
- [x] Unit miss: safe pattern → no `CWE-359`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-359-vulnerable.txt` + `CWE-359-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 5: `CWE-613` — Insufficient Session Expiration

> Catalogue: `ruleset/python/chunks/cwe-601-650.json` · tier **B** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE613` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-613", detectCWE613, &Meta…, gates...)` in `rules_auth.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-613`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `SESSION_COOKIE_AGE\s*=\s*0`, `permanent_session_lifetime missing`, `SESSION_EXPIRE_AT_BROWSER_CLOSE\s*=\s*False with huge AGE`
- [x] Implement `detectCWE613` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-613`
- [x] Unit miss: safe pattern → no `CWE-613`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-613-vulnerable.txt` + `CWE-613-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 6: `CWE-565` — Reliance on Cookies without Validation and Integrity Checking

> Catalogue: `ruleset/python/chunks/cwe-551-600.json` · tier **B** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE565` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-565", detectCWE565, &Meta…, gates...)` in `rules_auth.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-565`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `request.COOKIES[`, `request.cookies.get(`, `session id from cookie used as authz key without HMAC/signing`
- [x] Implement `detectCWE565` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-565`
- [x] Unit miss: safe pattern → no `CWE-565`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-565-vulnerable.txt` + `CWE-565-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 7: `CWE-807` — Reliance on Untrusted Inputs in a Security Decision

> Catalogue: `ruleset/python/chunks/cwe-801-850.json` · tier **B** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE807` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-807", detectCWE807, &Meta…, gates...)` in `rules_auth.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-807`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Manual Static Analysis, Automated Static Analysis - Binary or Bytecode, Manual Static Analysis - Binary or Bytecode, Dynamic Analysis with Automated Results Interpretation, Dynamic Analysis with …
- Suggested sinks/patterns: `if request.headers.get(\"X-Admin\")`, `if request.GET.get(\"role\")==`, `if request.cookies.get(\"is_admin\")`, `HTTP_X_FORWARDED_* for auth`
- [x] Implement `detectCWE807` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-807`
- [x] Unit miss: safe pattern → no `CWE-807`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-807-vulnerable.txt` + `CWE-807-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 8: `CWE-698` — Execution After Redirect (EAR)

> Catalogue: `ruleset/python/chunks/cwe-651-700.json` · tier **B** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE698` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-698", detectCWE698, &Meta…, gates...)` in `rules_auth.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-698`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Black Box. This issue might not be detected if testing is performed using a web browser, because the browser might obey the redirect and move the user to a different page before the application h…
- Suggested sinks/patterns: `return redirect missing: redirect(...); code continues`, `HttpResponseRedirect without return`
- [x] Implement `detectCWE698` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-698`
- [x] Unit miss: safe pattern → no `CWE-698`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-698-vulnerable.txt` + `CWE-698-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
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
