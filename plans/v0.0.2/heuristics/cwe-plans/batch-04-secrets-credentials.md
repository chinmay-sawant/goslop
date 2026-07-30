# Batch 04 — Secrets Credentials

> **Parent:** `plans/v0.0.2/heuristics/cwe-plans/README.md` — v0.0.2 Python CWE remaining heuristics
> **Epic / issue:** [#51](https://github.com/chinmay-sawant/goslop/issues/51) · [#52](https://github.com/chinmay-sawant/goslop/issues/52) expansion
> **Status:** **pending** — not started
> **Wave:** P0
> **IDs (8):** CWE-798, CWE-256, CWE-260, CWE-261, CWE-312, CWE-319, CWE-547, CWE-523
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
| Target file(s) | `rules_secrets.go` |
| Validation | `make lint` + `make test` unchecked until green |
| Fixtures | `tests/fixtures/python/cwe/CWE-N-{{vulnerable,safe}}.txt` |

## Overview

Hardcoded credentials, cleartext secrets, password storage, weak transport of secrets

| Rule | Name | Tier | Category | Relevance | Chunk |
|------|------|------|----------|-----------|-------|
| `CWE-798` | Use of Hard-coded Credentials | A | Security | Medium | `cwe-751-800.json` |
| `CWE-256` | Plaintext Storage of a Password | A | Security | Medium | `cwe-251-300.json` |
| `CWE-260` | Password in Configuration File | A | Security | Medium | `cwe-251-300.json` |
| `CWE-261` | Weak Encoding for Password | A | Security | Medium | `cwe-251-300.json` |
| `CWE-312` | Cleartext Storage of Sensitive Information | A | Information Disclosure | Medium | `cwe-301-350.json` |
| `CWE-319` | Cleartext Transmission of Sensitive Information | A | Information Disclosure | Medium | `cwe-301-350.json` |
| `CWE-547` | Use of Hard-coded, Security-relevant Constants | A | Configuration | Medium | `cwe-501-550.json` |
| `CWE-523` | Unprotected Transport of Credentials | A | Security | Medium | `cwe-501-550.json` |

## Executive Summary

Ship **8** remaining CWE heuristics in this theme. Prefer high-signal, low-FP patterns; same-file heuristics only (no interproc taint). Extend registration catalogue tests and fixture matrix.

## Phase 0: Placement + budget

- [ ] Record baseline: `wc -l internal/lang/python/detectors/cwe/*.go`
- [ ] Ensure target `rules_secrets.go` exists (create if needed); keep ≤1500 lines projected
- [ ] Prefer domain file over growing `rules.go` past soft cap
- [ ] Append FN-safe needles to `needles.go` / per-rule gates
- [ ] No changes to `detectors/all.go` unless constructor rename (should be none)

## Phase 1: `CWE-798` — Use of Hard-coded Credentials

> Catalogue: `ruleset/python/chunks/cwe-751-800.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE798` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-798", detectCWE798, &Meta…, gates...)` in `rules_secrets.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-798`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Black Box, Automated Static Analysis, Manual Static Analysis, Manual Dynamic Analysis, Automated Static Analysis - Binary or Bytecode, Manual Static Analysis - Binary or Bytecode. Credential stor…
- Suggested sinks/patterns: `password\s*=\s*['\"][^'"]{3,}['\"]`, `api_key\s*=\s*['\"]`, `secret_key\s*=\s*['\"]`, `AWS_SECRET`, `PRIVATE_KEY\s*=\s*['\"]-----BEGIN`, `token\s*=\s*['\"][A-Za-z0-9_\-]{16,}`
- [ ] Implement `detectCWE798` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-798`
- [ ] Unit miss: safe pattern → no `CWE-798`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-798-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 2: `CWE-256` — Plaintext Storage of a Password

> Catalogue: `ruleset/python/chunks/cwe-251-300.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE256` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-256", detectCWE256, &Meta…, gates...)` in `rules_secrets.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-256`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `password\s*=\s*['\"][^'"]+['\"]`, `PASSWORD\s*=`, `passwd\s*=\s*literal`, `settings.DATABASES PASSWORD literals`
- [ ] Implement `detectCWE256` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-256`
- [ ] Unit miss: safe pattern → no `CWE-256`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-256-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 3: `CWE-260` — Password in Configuration File

> Catalogue: `ruleset/python/chunks/cwe-251-300.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE260` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-260", detectCWE260, &Meta…, gates...)` in `rules_secrets.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-260`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `settings.py PASSWORD=`, `config['password']='…'`, `.ini/.yaml password keys as literals`, `django DATABASES password string`
- [ ] Implement `detectCWE260` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-260`
- [ ] Unit miss: safe pattern → no `CWE-260`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-260-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 4: `CWE-261` — Weak Encoding for Password

> Catalogue: `ruleset/python/chunks/cwe-251-300.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE261` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-261", detectCWE261, &Meta…, gates...)` in `rules_secrets.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-261`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `base64.b64encode(.*password`, `base64.b64encode(password`, `codecs.encode(...,'rot_13') password`, `binascii.hexlify(password)`
- [ ] Implement `detectCWE261` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-261`
- [ ] Unit miss: safe pattern → no `CWE-261`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-261-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 5: `CWE-312` — Cleartext Storage of Sensitive Information

> Catalogue: `ruleset/python/chunks/cwe-301-350.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE312` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-312", detectCWE312, &Meta…, gates...)` in `rules_secrets.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-312`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `SECRET_KEY = '…'`, `API_KEY = '…'`, `aws_secret_access_key literal`, `token = 'sk-…'`
- [ ] Implement `detectCWE312` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-312`
- [ ] Unit miss: safe pattern → no `CWE-312`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-312-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 6: `CWE-319` — Cleartext Transmission of Sensitive Information

> Catalogue: `ruleset/python/chunks/cwe-301-350.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE319` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-319", detectCWE319, &Meta…, gates...)` in `rules_secrets.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-319`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Black Box, Automated Static Analysis. Use monitoring tools that examine the software's process as it interacts with the operating system and the network. This technique is useful in cases when so…
- Suggested sinks/patterns: `http:// with basic auth or password query`, `smtplib.SMTP without starttls`, `ftplib.FTP(`, `IMAP4 without SSL for creds`
- [ ] Implement `detectCWE319` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-319`
- [ ] Unit miss: safe pattern → no `CWE-319`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-319-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 7: `CWE-547` — Use of Hard-coded, Security-relevant Constants

> Catalogue: `ruleset/python/chunks/cwe-501-550.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE547` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-547", detectCWE547, &Meta…, gates...)` in `rules_secrets.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-547`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `SECURE_SSL_REDIRECT\s*=\s*False`, `SESSION_COOKIE_SECURE\s*=\s*False`, `CSRF_COOKIE_SECURE\s*=\s*False`, `SESSION_COOKIE_HTTPONLY\s*=\s*False`, `SECURE_HSTS_SECONDS\s*=\s*0`, `ALLOWED_HOSTS\s*=\s*\[\s*['\"]\*['\"]`
- [ ] Implement `detectCWE547` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-547`
- [ ] Unit miss: safe pattern → no `CWE-547`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-547-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 8: `CWE-523` — Unprotected Transport of Credentials

> Catalogue: `ruleset/python/chunks/cwe-501-550.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE523` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-523", detectCWE523, &Meta…, gates...)` in `rules_secrets.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-523`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `verify\s*=\s*False`, `ssl._create_unverified_context`, `CERT_NONE`, `check_hostname\s*=\s*False`, `http://`, `urllib3.disable_warnings`, `requests.get/post(... verify=False)`
- [ ] Implement `detectCWE523` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-523`
- [ ] Unit miss: safe pattern → no `CWE-523`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-523-vulnerable.txt` + `-safe.txt`
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

