# Batch 08 — Production hardening BP-PY heuristics

> **Parent:** `plans/v0.0.2/heuristics/bp-plans/README.md` — batch index  
> **Canonical #53 ledger:** `plans/v0.0.2/heuristics/python-heuristics-bp.md`  
> **Issue:** [#53](https://github.com/chinmay-sawant/goslop/issues/53) expansion under epic [#51](https://github.com/chinmay-sawant/goslop/issues/51)  
> **Status:** implemented — BP-PY-48/49/50 in `rules_prod.go` + `rules_prod_test.go` 
> **Estimated effort:** 1 PR (small–medium)  
> **PR policy:** single batchwise PR; title e.g. `python(bp): batch-08 prod hardening (BP-PY-48..50)`

---

## Overview

| Item | Value |
|------|-------|
| IDs **in this batch** | **BP-PY-48**, **BP-PY-49**, **BP-PY-50** |
| Explicitly **out of this batch** | **BP-PY-14** (`requests` without timeout) — owned by **batch-01** (or sibling core/HTTP batch). **Do not implement or checklist-duplicate BP-PY-14 here.** |
| Category | Production Hardening |
| Target rule file | **`internal/lang/python/detectors/bad_practices/rules_prod.go`** (create if missing; **append** if batch-01 already added `BP-PY-14`/`15` there) |
| Line budget | **1500 soft / 2000 hard** per Go file — if shared with 14/15 and file approaches cap, split e.g. `rules_prod_cors.go` only if needed |
| Detection | pure-Go source patterns + needles; no framework type inference |

### Catalogue contract

| ID | Name | Sev | Detection notes (v0) |
|----|------|-----|----------------------|
| BP-PY-48 | CORS Allow Origins Star With Credentials | high | CORSMiddleware / flask-cors / django-cors-headers: `allow_origins=['*']` (or True) **with** `allow_credentials=True` |
| BP-PY-49 | TLS Verification Disabled | high | `verify=False`, `ssl._create_unverified_context`, `CERT_NONE` on requests/httpx/urllib3 paths |
| BP-PY-50 | Django/Flask CSRF Or Session Cookie Insecure Flags | medium | `SESSION_COOKIE_SECURE` / `CSRF_COOKIE_SECURE` / `SESSION_COOKIE_HTTPONLY` = `False`; Flask session cookie secure false |

### Related ID ownership (do not steal)

| ID | Owner | Note |
|----|-------|------|
| BP-PY-14 | batch-01 (sibling) | also Production Hardening category — **different PR** |
| BP-PY-16/17/21 | batch-00 shipped | DEBUG/SECRET already done; do not re-open |
| BP-PY-48–50 | **this batch** | only |

---

## Executive Summary

```text
rules_prod.go init()
  ├─ RegisterRule("BP-PY-48", detectBPPY48)
  ├─ RegisterRule("BP-PY-49", detectBPPY49)
  └─ RegisterRule("BP-PY-50", detectBPPY50)
# BP-PY-14 is NOT registered in this file
```

Ship three high/medium production hardening detectors with hit/miss unit tests; validate with `make lint` + `make test`.

---

## Phase 1: Scaffold `rules_prod.go`

### 1.1 File + registration

- [x] Create or open `internal/lang/python/detectors/bad_practices/rules_prod.go` (batch-01 may already own 14/15 in this file — **do not remove** those registrations)
- [x] This batch’s `RegisterRule` additions: **BP-PY-48**, **BP-PY-49**, **BP-PY-50** only (this PR must not add or rework 14)
- [x] Assert in review: this PR does **not** introduce `RegisterRule("BP-PY-14", …)` (14 stays batch-01)
- [x] Keep file under 1500 soft line budget; split only if shared file hits hard cap

### 1.2 Facts / needles

- [x] Extend `bpNeedles` in `facts.go` for e.g. `CORSMiddleware`, `allow_origins`, `allow_credentials`, `verify=False`, `_create_unverified_context`, `CERT_NONE`, `SESSION_COOKIE_SECURE`, `CSRF_COOKIE_SECURE`, `SESSION_COOKIE_HTTPONLY` (subset used by detectors)
- [x] Early-out with `facts.has` / `hasAny` where applicable

### 1.3 Metadata

- [x] Metadata non-nil: 48 **high**, 49 **high**, 50 **medium**
- [x] Pack bad-practice; IDs remain `BP-PY-*`

---

## Phase 2: `BP-PY-48` — CORS `*` with credentials

### 2.1 Detector behavior

- [x] Implement `detectBPPY48` in `rules_prod.go`
- [x] Hit patterns (any one combination in same unit is enough for v0):
  - FastAPI/Starlette: `CORSMiddleware` (or import) with `allow_origins=["*"]` / `['*']` / `allow_origins='*'` **and** `allow_credentials=True`
  - flask-cors: `CORS(..., supports_credentials=True)` with origins `*` / `resources` star heuristic
  - django-cors-headers: `CORS_ALLOW_ALL_ORIGINS = True` (or `CORS_ORIGIN_ALLOW_ALL`) **with** `CORS_ALLOW_CREDENTIALS = True`
- [x] Miss: star origins **without** credentials true; miss explicit origin list with credentials
- [x] Skip test files if noise is high (`isPythonTestFile`) — document
- [x] Message: star + credentials is unsafe for browsers; name concrete origins

### 2.2 Proof

- [x] Hit: FastAPI-style middleware snippet with `allow_origins=["*"]` and `allow_credentials=True`
- [x] Miss: `allow_origins=["https://app.example"]` + `allow_credentials=True`
- [x] Miss: `allow_origins=["*"]` without credentials true
- [x] Optional hit: Django settings pair `CORS_ALLOW_ALL_ORIGINS = True` + `CORS_ALLOW_CREDENTIALS = True`
- [x] Register + `RuleIDs` includes `BP-PY-48`
- [x] Severity high from metadata

---

## Phase 3: `BP-PY-49` — TLS verification disabled

### 3.1 Detector behavior

- [x] Implement `detectBPPY49`
- [x] Hit: `verify=False` on HTTP client calls (`requests.*`, `httpx.*`, session methods) — keyword in call window
- [x] Hit: `ssl._create_unverified_context`, `ssl.CERT_NONE`, `CERT_NONE` assigned into SSL context usage
- [x] Hit: urllib3 `assert_hostname=False` only if cheap and low FP; optional
- [x] Miss: `verify=True`, default verify omitted, `verify="/path/to/ca.pem"`
- [x] Skip test files / fixtures paths when basename suggests cert test vectors (optional; prefer skip `isPythonTestFile`)
- [x] Severity **high**; message: do not disable TLS verify in production

### 3.2 Proof

- [x] Hit: `requests.get(url, verify=False)\n`
- [x] Hit: `httpx.get(url, verify=False)\n` (if pattern shared)
- [x] Miss: `requests.get(url, timeout=5)\n` (no verify=False) — must **not** fire 49; also must **not** fire 14 from this detector
- [x] Miss: `requests.get(url, verify=True)\n`
- [x] Optional hit: `ssl._create_unverified_context()`
- [x] Register + tests for `BP-PY-49`
- [x] Explicit test or review note: this PR does **not** add BP-PY-14

---

## Phase 4: `BP-PY-50` — insecure session/CSRF cookie flags

### 4.1 Detector behavior

- [x] Implement `detectBPPY50`
- [x] Hit (Django settings-style assignments):
  - `SESSION_COOKIE_SECURE = False`
  - `CSRF_COOKIE_SECURE = False`
  - `SESSION_COOKIE_HTTPONLY = False`
- [x] Hit (Flask): `SESSION_COOKIE_SECURE = False` or `app.config['SESSION_COOKIE_SECURE'] = False`
- [x] Prefer flagging in settings-like paths (`settings.py`, `settings/`, flask config modules) when path heuristics are easy; if path unknown, still flag explicit `= False` assignments (document FP policy)
- [x] Miss: `= True` for the same flags
- [x] Optional: when `DEBUG = False` nearby and Secure flags missing entirely — **only if** low FP; otherwise stick to explicit False assignments for v0
- [x] Severity **medium**

### 4.2 Proof

- [x] Hit: `settings.py` with `SESSION_COOKIE_SECURE = False\n`
- [x] Hit: `CSRF_COOKIE_SECURE = False\n`
- [x] Miss: `SESSION_COOKIE_SECURE = True\n`
- [x] Optional Flask config hit/miss pair
- [x] Register + tests for `BP-PY-50`

---

## Phase 5: Package integration

### 5.1 Registration

- [x] `RuleIDs()` contains 48, 49, 50
- [x] `TestBPRulesRegistered` want-list includes 48–50
- [x] Confirm `RuleIDs()` does **not** gain BP-PY-14 from this PR
- [x] No bare `BP-*` IDs

### 5.2 Shared helpers

- [x] Reuse `pushAt`, `isPythonTestFile`, settings path helpers if already present in `rules_framework.go` — **extract to `common.go` only if needed by two domains**; watch line budgets on `common.go` (~446 already)
- [x] Prefer local helpers inside `rules_prod.go` when single-use

### 5.3 Line-limit gate

- [x] `wc -l rules_prod.go` under 1500 soft / 2000 hard
- [x] Touched test files under soft cap or split to `rules_prod_test.go`

---

## Phase 6: Validation gates (required for code)

- [x] `gofmt -w` on all touched Go files
- [x] `go test ./internal/lang/python/detectors/bad_practices/ -count=1` green
- [x] `make lint` — green; record: go vet + gofmt clean
- [x] `make test` — green; record: `go test ./...` ok
- [ ] Optional: `CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop`
- [x] Parent `python-heuristics-bp.md`: mark 48–50 `[x]` with evidence only
- [x] `_inventory.json`: move 48–50 to `implemented`; leave **14** missing until batch-01
- [x] PR: `Relates to #53`, `Relates to #51`

---

## Dependencies

| Depends on | Note |
|------------|------|
| Batch 00 | scaffold + metadata for catalogue IDs |
| Catalogue | BP-PY-48–50 in `ruleset/python/bad-practices.json` |
| Not this PR | BP-PY-14 (batch-01), async 38–40, testing/deps 41–47 |
| Framework overlap | Django settings path heuristics may mirror `rules_framework.go` BP-PY-21 — reuse carefully |

---

## Non-goals

- Implementing **BP-PY-14** requests timeout (batch-01)  
- Full CORS middleware AST config resolution across multiple files  
- Certificate pinning recommendations beyond verify=False detection  
- Changing default Go-only registry  

---

## Complete ID checklist (none deferred; 14 excluded by ownership)

- [x] **BP-PY-48** — detector + register + hit/miss + RuleIDs
- [x] **BP-PY-49** — detector + register + hit/miss + RuleIDs
- [x] **BP-PY-50** — detector + register + hit/miss + RuleIDs
- [x] **BP-PY-14** — **not in this batch** (owned by batch-01 / sibling); no pending row here

**Batch-08 ID list:** BP-PY-48, BP-PY-49, BP-PY-50  

**Explicit non-list:** BP-PY-14 → batch-01
