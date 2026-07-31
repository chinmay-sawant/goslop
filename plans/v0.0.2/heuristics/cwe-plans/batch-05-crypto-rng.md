# Batch 05 — Crypto Rng

> **Parent:** `plans/v0.0.2/heuristics/cwe-plans/README.md` — v0.0.2 Python CWE remaining heuristics
> **Epic / issue:** [#51](https://github.com/chinmay-sawant/goslop/issues/51) · [#52](https://github.com/chinmay-sawant/goslop/issues/52) expansion
> **Status:** **complete** — implemented and validated 2026-07-31
> **Validation evidence:** `gofmt`, Python CWE package tests, Python CWE fixture matrix, `make lint`, `make test`, and `git diff --check` passed.
> **Wave:** P1
> **IDs (9):** CWE-295, CWE-328, CWE-335, CWE-338, CWE-347, CWE-1204, CWE-1240, CWE-1241, CWE-1392
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
| Target file(s) | `rules_crypto.go` |
| Validation | `make lint` + `make test` unchecked until green |
| Fixtures (text) | **Required:** `tests/fixtures/python/cwe/CWE-N-{vulnerable,safe}.txt` (BP-style; see Testing section) |
| Unit tests | **Required:** hit/miss in `internal/lang/python/detectors/cwe/*_test.go` |
| Integration | **Required:** `tests/integration/python/cwe_matrix_test.go` / `make integration-python` |

## Overview

TLS verify disable, weak hashes, insecure RNG, password hashing mistakes

| Rule | Name | Tier | Category | Relevance | Chunk |
|------|------|------|----------|-----------|-------|
| `CWE-295` | Improper Certificate Validation | A | Security | Medium | `cwe-251-300.json` |
| `CWE-328` | Use of Weak Hash | A | General | Medium | `cwe-301-350.json` |
| `CWE-335` | Incorrect Usage of Seeds in Pseudo-Random Number Generator (PRNG) | A | Security | Medium | `cwe-301-350.json` |
| `CWE-338` | Use of Cryptographically Weak Pseudo-Random Number Generator (PRNG) | A | Security | Medium | `cwe-301-350.json` |
| `CWE-347` | Improper Verification of Cryptographic Signature | A | Security | Medium | `cwe-301-350.json` |
| `CWE-1204` | Generation of Weak Initialization Vector (IV) | A | Cryptography | Medium | `cwe-1201-1250.json` |
| `CWE-1240` | Use of a Cryptographic Primitive with a Risky Implementation | A | Security | Medium | `cwe-1201-1250.json` |
| `CWE-1241` | Use of Predictable Algorithm in Random Number Generator | A | Security | Medium | `cwe-1201-1250.json` |
| `CWE-1392` | Use of Default Credentials | A | Security | Medium | `cwe-1351-1400.json` |

## Executive Summary

Ship **9** remaining CWE heuristics in this theme. Prefer high-signal, low-FP patterns; same-file heuristics only (no interproc taint). Extend registration catalogue tests and fixture matrix.

## Phase 0: Placement + budget

- [x] Record baseline: `wc -l internal/lang/python/detectors/cwe/*.go`
- [x] Ensure target `rules_crypto.go` exists (create if needed); keep ≤1500 lines projected
- [x] Prefer domain file over growing `rules.go` past soft cap
- [x] Append FN-safe needles to `needles.go` / per-rule gates
- [x] No changes to `detectors/all.go` unless constructor rename (should be none)

## Phase 1: `CWE-295` — Improper Certificate Validation

> Catalogue: `ruleset/python/chunks/cwe-251-300.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE295` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-295", detectCWE295, &Meta…, gates...)` in `rules_crypto.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-295`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis - Binary or Bytecode, Manual Static Analysis - Binary or Bytecode, Dynamic Analysis with Automated Results Interpretation, Dynamic Analysis with Manual Results Interpret…
- Suggested sinks/patterns: `verify=False`, `ssl.CERT_NONE`, `check_hostname=False`, `ssl._create_unverified_context`, `urllib3 disable_warnings + no verify`
- [x] Implement `detectCWE295` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-295`
- [x] Unit miss: safe pattern → no `CWE-295`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-295-vulnerable.txt` + `CWE-295-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 2: `CWE-328` — Use of Weak Hash

> Catalogue: `ruleset/python/chunks/cwe-301-350.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE328` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-328", detectCWE328, &Meta…, gates...)` in `rules_crypto.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-328`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `hashlib.md5(`, `hashlib.sha1(`, `hashlib.new('md5'`, `hashlib.new('sha1'`, `Crypto.Hash.MD5`
- [x] Implement `detectCWE328` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-328`
- [x] Unit miss: safe pattern → no `CWE-328`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-328-vulnerable.txt` + `CWE-328-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 3: `CWE-335` — Incorrect Usage of Seeds in Pseudo-Random Number Generator (PRNG)

> Catalogue: `ruleset/python/chunks/cwe-301-350.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE335` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-335", detectCWE335, &Meta…, gates...)` in `rules_crypto.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-335`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `random.seed(0)`, `random.seed(constant)`, `numpy.random.seed(fixed) for security tokens`
- [x] Implement `detectCWE335` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-335`
- [x] Unit miss: safe pattern → no `CWE-335`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-335-vulnerable.txt` + `CWE-335-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 4: `CWE-338` — Use of Cryptographically Weak Pseudo-Random Number Generator (PRNG)

> Catalogue: `ruleset/python/chunks/cwe-301-350.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE338` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-338", detectCWE338, &Meta…, gates...)` in `rules_crypto.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-338`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `random.random/randint/choice/getrandbits used for token/secret/password/key/session/nonce/otp/csrf`
- [x] Implement `detectCWE338` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-338`
- [x] Unit miss: safe pattern → no `CWE-338`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-338-vulnerable.txt` + `CWE-338-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 5: `CWE-347` — Improper Verification of Cryptographic Signature

> Catalogue: `ruleset/python/chunks/cwe-301-350.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE347` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-347", detectCWE347, &Meta…, gates...)` in `rules_crypto.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-347`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `jwt.decode(..., verify=False)`, `options={'verify_signature': False}`, `PyJWT without key`, `hmac.compare_digest missing / == on mac`
- [x] Implement `detectCWE347` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-347`
- [x] Unit miss: safe pattern → no `CWE-347`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-347-vulnerable.txt` + `CWE-347-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 6: `CWE-1204` — Generation of Weak Initialization Vector (IV)

> Catalogue: `ruleset/python/chunks/cwe-1201-1250.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE1204` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-1204", detectCWE1204, &Meta…, gates...)` in `rules_crypto.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-1204`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `AES.new(..., MODE_CBC, iv=b'0'*16)`, `fixed IV literal`, `Crypto.Cipher with constant iv`
- [x] Implement `detectCWE1204` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-1204`
- [x] Unit miss: safe pattern → no `CWE-1204`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-1204-vulnerable.txt` + `CWE-1204-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 7: `CWE-1240` — Use of a Cryptographic Primitive with a Risky Implementation

> Catalogue: `ruleset/python/chunks/cwe-1201-1250.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE1240` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-1240", detectCWE1240, &Meta…, gates...)` in `rules_crypto.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-1240`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Architecture or Design Review, Manual Analysis, Dynamic Analysis with Manual Results Interpretation. Review requirements, documentation, and product design to ensure that primitives are consisten…
- Suggested sinks/patterns: `custom xorCipher`, `homegrown encrypt without AES-GCM`, `roll-your-own crypto helpers`
- [x] Implement `detectCWE1240` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-1240`
- [x] Unit miss: safe pattern → no `CWE-1240`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-1240-vulnerable.txt` + `CWE-1240-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 8: `CWE-1241` — Use of Predictable Algorithm in Random Number Generator

> Catalogue: `ruleset/python/chunks/cwe-1201-1250.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE1241` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-1241", detectCWE1241, &Meta…, gates...)` in `rules_crypto.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-1241`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `random.random/randint for tokens/secrets`, `random.choice for session ids`, `numpy.random for security`
- [x] Implement `detectCWE1241` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-1241`
- [x] Unit miss: safe pattern → no `CWE-1241`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-1241-vulnerable.txt` + `CWE-1241-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 9: `CWE-1392` — Use of Default Credentials

> Catalogue: `ruleset/python/chunks/cwe-1351-1400.json` · tier **A** · relevance **Medium**

### Register + meta

- [x] Add `MetaCWE1392` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [x] `RegisterRule("CWE-1392", detectCWE1392, &Meta…, gates...)` in `rules_crypto.go` `init()`
- [x] Confirm ID exists in chunk JSON key `CWE-1392`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `password='admin'`, `USERNAME/PASSWORD default pairs`, `Django createsuperuser hardcode`, `MQTT/redis default passwords`
- [x] Implement `detectCWE1392` with needle prefilter
- [x] Prefer high-signal positive; document safe suppressions
- [x] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures (triad)

- [x] **Unit hit:** vulnerable snippet → finding `CWE-1392`
- [x] Unit miss: safe pattern → no `CWE-1392`
- [x] **Fixtures:** `tests/fixtures/python/cwe/CWE-1392-vulnerable.txt` + `CWE-1392-safe.txt` (`lang: python` header; same `.txt` format as `tests/fixtures/python/bp/`)
- [x] **Integration:** pair auto-discovered by `DiscoverPythonCWECases` → `TestPythonCWEFixturesMatrix` in `tests/integration/python/cwe_matrix_test.go`
- [x] Run: `go test ./internal/lang/python/detectors/cwe/ -count=1` and `go test ./tests/integration/python/ -count=1` (or `make integration-python`)

### Proof

- [x] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 10: Batch validation

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
