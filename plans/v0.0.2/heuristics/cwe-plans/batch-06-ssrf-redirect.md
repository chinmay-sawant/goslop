# Batch 06 — Ssrf Redirect

> **Parent:** `plans/v0.0.2/heuristics/cwe-plans/README.md` — v0.0.2 Python CWE remaining heuristics
> **Epic / issue:** [#51](https://github.com/chinmay-sawant/goslop/issues/51) · [#52](https://github.com/chinmay-sawant/goslop/issues/52) expansion
> **Status:** **pending** — not started
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
| Fixtures | `tests/fixtures/python/cwe/CWE-N-{{vulnerable,safe}}.txt` |

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

- [ ] Record baseline: `wc -l internal/lang/python/detectors/cwe/*.go`
- [ ] Ensure target `rules_ssrf.go` exists (create if needed); keep ≤1500 lines projected
- [ ] Prefer domain file over growing `rules.go` past soft cap
- [ ] Append FN-safe needles to `needles.go` / per-rule gates
- [ ] No changes to `detectors/all.go` unless constructor rename (should be none)

## Phase 1: `CWE-918` — Server-Side Request Forgery (SSRF)

> Catalogue: `ruleset/python/chunks/cwe-901-950.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE918` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-918", detectCWE918, &Meta…, gates...)` in `rules_ssrf.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-918`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `requests.get/post(url)`, `urllib.request.urlopen(url)`, `httpx.get(url)`, `aiohttp session.get(url)`, `urllib3`
- [ ] Implement `detectCWE918` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-918`
- [ ] Unit miss: safe pattern → no `CWE-918`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-918-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 2: `CWE-601` — URL Redirection to Untrusted Site ('Open Redirect')

> Catalogue: `ruleset/python/chunks/cwe-601-650.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE601` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-601", detectCWE601, &Meta…, gates...)` in `rules_ssrf.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-601`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Manual Static Analysis, Automated Dynamic Analysis, Automated Static Analysis, Automated Static Analysis - Binary or Bytecode, Dynamic Analysis with Automated Results Interpretation, Dynamic Anal…
- Suggested sinks/patterns: `redirect(`, `HttpResponseRedirect(`, `RedirectResponse(`, `flask.redirect`, `django.shortcuts.redirect`, `url_for with _external + user next=`
- [ ] Implement `detectCWE601` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-601`
- [ ] Unit miss: safe pattern → no `CWE-601`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-601-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 3: `CWE-605` — Multiple Binds to the Same Port

> Catalogue: `ruleset/python/chunks/cwe-601-650.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE605` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-605", detectCWE605, &Meta…, gates...)` in `rules_ssrf.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-605`

### Detect heuristic

- Cite **detection_notes** (abbrev): MITRE detection methods: Automated Static Analysis. Automated static analysis, commonly referred to as Static Application Security Testing (SAST), can find some instances of this weakness by analyzing source code (or bin…
- Suggested sinks/patterns: `SO_REUSEADDR`, `setsockopt(.*SO_REUSEADDR`, `socket.bind((\"0.0.0.0\"`, `bind((\"::\"`, `app.run(host=\"0.0.0.0\")`
- [ ] Implement `detectCWE605` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-605`
- [ ] Unit miss: safe pattern → no `CWE-605`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-605-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 4: `CWE-924` — Improper Enforcement of Message Integrity During Transmission in a Communication Channel

> Catalogue: `ruleset/python/chunks/cwe-901-950.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE924` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-924", detectCWE924, &Meta…, gates...)` in `rules_ssrf.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-924`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `webhook handler decode body without hmac/signature`, `stripe/github webhook without X-*-Signature verify`
- [ ] Implement `detectCWE924` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-924`
- [ ] Unit miss: safe pattern → no `CWE-924`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-924-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 5: `CWE-940` — Improper Verification of Source of a Communication Channel

> Catalogue: `ruleset/python/chunks/cwe-901-950.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE940` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-940", detectCWE940, &Meta…, gates...)` in `rules_ssrf.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-940`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `OAuth callback without state/nonce check`, `request.args user_id trusted as identity`
- [ ] Implement `detectCWE940` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-940`
- [ ] Unit miss: safe pattern → no `CWE-940`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-940-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 6: `CWE-941` — Incorrectly Specified Destination in a Communication Channel

> Catalogue: `ruleset/python/chunks/cwe-901-950.json` · tier **A** · relevance **Medium**

### Register + meta

- [ ] Add `MetaCWE941` in `metadata.go` (or domain metadata file) — title/name match catalogue
- [ ] `RegisterRule("CWE-941", detectCWE941, &Meta…, gates...)` in `rules_ssrf.go` `init()`
- [ ] Confirm ID exists in chunk JSON key `CWE-941`

### Detect heuristic

- Cite **detection_notes** (abbrev): Python-oriented (future detector): prefer stdlib/framework sinks (pathlib/os, subprocess, django/flask/fastapi templates & ORMs, pickle/yaml, hashlib/secrets, urllib) over C memory APIs; rewrite detectors per language pl…
- Suggested sinks/patterns: `smtp/send_mail to request-controlled address`, `webhook URL from client without allowlist`
- [ ] Implement `detectCWE941` with needle prefilter
- [ ] Prefer high-signal positive; document safe suppressions
- [ ] Message catalogue-aligned; confidence documented

### Hit / miss tests + fixtures

- [ ] Unit hit: vulnerable snippet → finding `CWE-941`
- [ ] Unit miss: safe pattern → no `CWE-941`
- [ ] Fixtures: `tests/fixtures/python/cwe/CWE-941-vulnerable.txt` + `-safe.txt`
- [ ] Matrix discovers pair via integration python CWE suite

### Proof

- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` covers this ID

## Phase 7: Batch validation

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

