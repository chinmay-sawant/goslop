# Batch 16 — Tier B Expansion

> **Parent:** `plans/v0.0.2/heuristics/cwe-plans/README.md` — v0.0.2 Python CWE remaining heuristics
> **Epic / issue:** [#51](https://github.com/chinmay-sawant/goslop/issues/51) · [#52](https://github.com/chinmay-sawant/goslop/issues/52) expansion
> **Status:** **pending** — not started
> **Wave:** P3
> **IDs (60):** CWE-66, CWE-76, CWE-178, CWE-179, CWE-182, CWE-184, CWE-186, CWE-257, CWE-272, CWE-279, CWE-289, CWE-290, CWE-323, CWE-331, CWE-334, CWE-367, CWE-403, CWE-409, CWE-454, CWE-472, CWE-521, CWE-524, CWE-538, CWE-552, CWE-617, CWE-641, CWE-648, CWE-779, CWE-836, CWE-838, CWE-908, CWE-909, CWE-910, CWE-911, CWE-920, CWE-939, CWE-1007, CWE-1021, CWE-1046, CWE-1050…
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
| Target file(s) | `rules_tier_b_*.go (split by theme when implementing)` |
| Validation | `make lint` + `make test` unchecked until green |
| Fixtures (text) | **Required:** `tests/fixtures/python/cwe/CWE-N-{vulnerable,safe}.txt` (BP-style; see Testing section) |
| Unit tests | **Required:** hit/miss in `internal/lang/python/detectors/cwe/*_test.go` |
| Integration | **Required:** `tests/integration/python/cwe_matrix_test.go` / `make integration-python` |

## Overview

All remaining tier-B — optional expansion after A waves; split into sub-PRs if >15 IDs at implement time

| Rule | Name | Tier | Category | Relevance | Chunk |
|------|------|------|----------|-----------|-------|
| `CWE-66` | Improper Handling of File Names that Identify Virtual Resources | B | File Processing | Medium | `cwe-051-100.json` |
| `CWE-76` | Improper Neutralization of Equivalent Special Elements | B | General | Medium | `cwe-051-100.json` |
| `CWE-178` | Improper Handling of Case Sensitivity | B | File Processing | Medium | `cwe-151-200.json` |
| `CWE-179` | Incorrect Behavior Order: Early Validation | B | General | Medium | `cwe-151-200.json` |
| `CWE-182` | Collapse of Data into Unsafe Value | B | General | Medium | `cwe-151-200.json` |
| `CWE-184` | Incomplete List of Disallowed Inputs | B | General | Medium | `cwe-151-200.json` |
| `CWE-186` | Overly Restrictive Regular Expression | B | General | Medium | `cwe-151-200.json` |
| `CWE-257` | Storing Passwords in a Recoverable Format | B | Security | Medium | `cwe-251-300.json` |
| `CWE-272` | Least Privilege Violation | B | Security | Medium | `cwe-251-300.json` |
| `CWE-279` | Incorrect Execution-Assigned Permissions | B | Security | Medium | `cwe-251-300.json` |
| `CWE-289` | Authentication Bypass by Alternate Name | B | Security | Medium | `cwe-251-300.json` |
| `CWE-290` | Authentication Bypass by Spoofing | B | Security | Medium | `cwe-251-300.json` |
| `CWE-323` | Reusing a Nonce, Key Pair in Encryption | B | General | Medium | `cwe-301-350.json` |
| `CWE-331` | Insufficient Entropy | B | General | Medium | `cwe-301-350.json` |
| `CWE-334` | Small Space of Random Values | B | Security | Medium | `cwe-301-350.json` |
| `CWE-367` | Time-of-check Time-of-use (TOCTOU) Race Condition | B | Resource Management | Medium | `cwe-351-400.json` |
| `CWE-403` | Exposure of File Descriptor to Unintended Control Sphere ('File Descriptor Leak') | B | Program Invocation | Medium | `cwe-401-450.json` |
| `CWE-409` | Improper Handling of Highly Compressed Data (Data Amplification) | B | General | Medium | `cwe-401-450.json` |
| `CWE-454` | External Initialization of Trusted Variables or Data Stores | B | General | Medium | `cwe-451-500.json` |
| `CWE-472` | External Control of Assumed-Immutable Web Parameter | B | General | Medium | `cwe-451-500.json` |
| `CWE-521` | Weak Password Requirements | B | Security | Medium | `cwe-501-550.json` |
| `CWE-524` | Use of Cache Containing Sensitive Information | B | Information Disclosure | Medium | `cwe-501-550.json` |
| `CWE-538` | Insertion of Sensitive Information into Externally-Accessible File or Directory | B | Path Traversal | Medium | `cwe-501-550.json` |
| `CWE-552` | Files or Directories Accessible to External Parties | B | General | Medium | `cwe-551-600.json` |
| `CWE-617` | Reachable Assertion | B | General | Medium | `cwe-601-650.json` |
| `CWE-641` | Improper Restriction of Names for Files and Other Resources | B | Resource Management | Medium | `cwe-601-650.json` |
| `CWE-648` | Incorrect Use of Privileged APIs | B | Security | Medium | `cwe-601-650.json` |
| `CWE-779` | Logging of Excessive Data | B | Information Disclosure | Medium | `cwe-751-800.json` |
| `CWE-836` | Use of Password Hash Instead of Password for Authentication | B | Security | Medium | `cwe-801-850.json` |
| `CWE-838` | Inappropriate Encoding for Output Context | B | General | Medium | `cwe-801-850.json` |
| `CWE-908` | Use of Uninitialized Resource | B | Resource Management | Medium | `cwe-901-950.json` |
| `CWE-909` | Missing Initialization of Resource | B | Resource Management | Medium | `cwe-901-950.json` |
| `CWE-910` | Use of Expired File Descriptor | B | Input Validation | Medium | `cwe-901-950.json` |
| `CWE-911` | Improper Update of Reference Count | B | General | Medium | `cwe-901-950.json` |
| `CWE-920` | Improper Restriction of Power Consumption | B | General | Medium | `cwe-901-950.json` |
| `CWE-939` | Improper Authorization in Handler for Custom URL Scheme | B | Security | Medium | `cwe-901-950.json` |
| `CWE-1007` | Insufficient Visual Distinction of Homoglyphs Presented to User | B | General | Medium | `cwe-1001-1050.json` |
| `CWE-1021` | Improper Restriction of Rendered UI Layers or Frames | B | General | Medium | `cwe-1001-1050.json` |
| `CWE-1046` | Creation of Immutable Text Using String Concatenation | B | General | Medium | `cwe-1001-1050.json` |
| `CWE-1050` | Excessive Platform Resource Consumption within a Loop | B | Resource Management | Medium | `cwe-1001-1050.json` |
| `CWE-1060` | Excessive Number of Inefficient Server-Side Data Accesses | B | General | Medium | `cwe-1051-1100.json` |
| `CWE-1067` | Excessive Execution of Sequential Searches of Data Resource | B | Resource Management | Medium | `cwe-1051-1100.json` |
| `CWE-1071` | Empty Code Block | B | Resource Management | Medium | `cwe-1051-1100.json` |
| `CWE-1072` | Data Resource Access without Use of Connection Pooling | B | Resource Management | Medium | `cwe-1051-1100.json` |
| `CWE-1084` | Invokable Control Element with Excessive File or Data Access Operations | B | General | Medium | `cwe-1051-1100.json` |
| `CWE-1104` | Use of Unmaintained Third Party Components | B | General | Medium | `cwe-1101-1150.json` |
| `CWE-1106` | Insufficient Use of Symbolic Constants | B | General | Medium | `cwe-1101-1150.json` |
| `CWE-1108` | Excessive Reliance on Global Variables | B | General | Medium | `cwe-1101-1150.json` |
| `CWE-1121` | Excessive McCabe Cyclomatic Complexity | B | General | Medium | `cwe-1101-1150.json` |
| `CWE-1123` | Excessive Use of Self-Modifying Code | B | General | Medium | `cwe-1101-1150.json` |
| `CWE-1124` | Excessively Deep Nesting | B | General | Medium | `cwe-1101-1150.json` |
| `CWE-1220` | Insufficient Granularity of Access Control | B | General | Medium | `cwe-1201-1250.json` |
| `CWE-1265` | Unintended Reentrant Invocation of Non-reentrant Code Via Nested Calls | B | General | Medium | `cwe-1251-1300.json` |
| `CWE-1284` | Improper Validation of Specified Quantity in Input | B | General | Medium | `cwe-1251-1300.json` |
| `CWE-1285` | Improper Validation of Specified Index, Position, or Offset in Input | B | General | Medium | `cwe-1251-1300.json` |
| `CWE-1287` | Improper Validation of Specified Type of Input | B | General | Medium | `cwe-1251-1300.json` |
| `CWE-1288` | Improper Validation of Consistency within Input | B | General | Medium | `cwe-1251-1300.json` |
| `CWE-1322` | Use of Blocking Code in Single-threaded, Non-blocking Context | B | Resource Management | Medium | `cwe-1301-1350.json` |
| `CWE-1339` | Insufficient Precision or Accuracy of a Real Number | B | General | Medium | `cwe-1301-1350.json` |
| `CWE-1341` | Multiple Releases of Same Resource or Handle | B | Resource Management | Medium | `cwe-1301-1350.json` |

## Executive Summary

Ship **60** remaining CWE heuristics in this theme. Prefer high-signal, low-FP patterns; same-file heuristics only (no interproc taint). Extend registration catalogue tests and fixture matrix.

> **Note:** This batch is large (all remaining tier-B). At implement time, **split into multiple PRs** of ≤15 IDs while keeping this ledger as the ownership index (or create batch-16a/b siblings).

## Phase 0: Placement + budget

- [ ] Record baseline: `wc -l internal/lang/python/detectors/cwe/*.go`
- [ ] Ensure target `rules_tier_b_*.go (split by theme when implementing)` exists (create if needed); keep ≤1500 lines projected
- [ ] Prefer domain file over growing `rules.go` past soft cap
- [ ] Append FN-safe needles to `needles.go` / per-rule gates
- [ ] No changes to `detectors/all.go` unless constructor rename (should be none)

## Phase 1: Bulk register plan (large batch)

Implement in **sub-PRs of ≤15 IDs**. For each ID below:

1. Meta + `RegisterRule` + gates
2. Detector with **individual** hit/miss unit tests (`*_test.go`)
3. Fixture pair `tests/fixtures/python/cwe/CWE-N-{vulnerable,safe}.txt` (BP-style `.txt` under `python/cwe/`, not `python/bp/`)
4. Confirm integration matrix picks up pairs (`TestPythonCWEFixturesMatrix`)

### ID checklist

- [ ] `CWE-66` — Improper Handling of File Names that Identify Virtual Resources — tier B — CON, NUL, COM1, \\.\
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-66-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-66`
- [ ] `CWE-76` — Improper Neutralization of Equivalent Special Elements — tier B — str.replace('<') without html.escape, manual denylist strip, missing html.escape/bleach
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-76-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-76`
- [ ] `CWE-178` — Improper Handling of Case Sensitivity — tier B — == on usernames/paths, dict lookup case, missing casefold/lower on security compare
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-178-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-178`
- [ ] `CWE-179` — Incorrect Behavior Order: Early Validation — tier B — validate then unquote/decode, re.match before urllib.parse.unquote
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-179-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-179`
- [ ] `CWE-182` — Collapse of Data into Unsafe Value — tier B — re.sub strip then auth check, filter then membership
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-182-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-182`
- [ ] `CWE-184` — Incomplete List of Disallowed Inputs — tier B — for bad in denylist: if bad in s, manual blocklists without escape
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-184-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-184`
- [ ] `CWE-186` — Overly Restrictive Regular Expression — tier B — re.compile for validation that is too narrow
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-186-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-186`
- [ ] `CWE-257` — Storing Passwords in a Recoverable Format — tier B — Fernet.encrypt(password), AES encrypt password for storage, base64 of password stored as 'hash'
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-257-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-257`
- [ ] `CWE-272` — Least Privilege Violation — tier B — os.setuid(0), run as root, sudo patterns, capability flags
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-272-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-272`
- [ ] `CWE-279` — Incorrect Execution-Assigned Permissions — tier B — os.chmod while running to widen access, tempfile then chmod 777
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-279-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-279`
- [ ] `CWE-289` — Authentication Bypass by Alternate Name — tier B — case-insensitive host/user compare missing, email.lower() inconsistently
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-289-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-289`
- [ ] `CWE-290` — Authentication Bypass by Spoofing — tier B — X-Forwarded-For trust without proxy config, REMOTE_USER trust, request.headers host trust
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-290-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-290`
- [ ] `CWE-323` — Reusing a Nonce, Key Pair in Encryption — tier B — fixed nonce/IV constant to AESGCM, iv = b'\x00'*12
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-323-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-323`
- [ ] `CWE-331` — Insufficient Entropy — tier B — token length < 16, secrets.token_bytes(2)
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-331-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-331`
- [ ] `CWE-334` — Small Space of Random Values — tier B — randint(0, 9999) for OTP/session, choice from tiny set
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-334-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-334`
- [ ] `CWE-367` — Time-of-check Time-of-use (TOCTOU) Race Condition — tier B — os.path.exists then open, lexists then remove
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-367-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-367`
- [ ] `CWE-403` — Exposure of File Descriptor to Unintended Control Sphere ('File Descriptor Leak') — tier B — subprocess with close_fds=False, pass open fd to child
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-403-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-403`
- [ ] `CWE-409` — Improper Handling of Highly Compressed Data (Data Amplification) — tier B — zipfile.extractall unrestricted, gzip without size cap, tarfile.extractall
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-409-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-409`
- [ ] `CWE-454` — External Initialization of Trusted Variables or Data Stores — tier B — load pickle/config from user path into globals, exec user config, importlib load user module as settings
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-454-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-454`
- [ ] `CWE-472` — External Control of Assumed-Immutable Web Parameter — tier B — trust hidden form field, request.POST used as price/role, mass assignment Model(**request.data)
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-472-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-472`
- [ ] `CWE-521` — Weak Password Requirements — tier B — min_length\s*=\s*[1-5]\b, PASSWORD_MIN_LENGTH\s*=\s*[1-5], validators.MinLengthValidator\(\s*[1-5]\s*\), len\(password\)\s*[<>=]=?\s*[1-5]
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-521-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-521`
- [ ] `CWE-524` — Use of Cache Containing Sensitive Information — tier B — cache.set(, lru_cache, functools.cache, django.core.cache
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-524-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-524`
- [ ] `CWE-538` — Insertion of Sensitive Information into Externally-Accessible File or Directory — tier B — open(...'w') to /tmp, static/, media/, ., logging to world paths, Path.write_text with secret names
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-538-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-538`
- [ ] `CWE-552` — Files or Directories Accessible to External Parties — tier B — send_file(, send_from_directory(, FileResponse(, StaticFilesHandler
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-552-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-552`
- [ ] `CWE-617` — Reachable Assertion — tier B — assert user.is_authenticated, assert permission, assert token
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-617-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-617`
- [ ] `CWE-641` — Improper Restriction of Names for Files and Other Resources — tier B — open/join with user filename without sanitization
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-641-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-641`
- [ ] `CWE-648` — Incorrect Use of Privileged APIs — tier B — os.setuid, os.seteuid, os.setgid, ctypes.CDLL
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-648-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-648`
- [ ] `CWE-779` — Logging of Excessive Data — tier B — log.*(password|secret|token|api_key|credit_card), logger.info(request.headers), print(password)
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-779-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-779`
- [ ] `CWE-836` — Use of Password Hash Instead of Password for Authentication — tier B — compare client-supplied hash to stored hash, authenticate(password_hash=request...)
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-836-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-836`
- [ ] `CWE-838` — Inappropriate Encoding for Output Context — tier B — urllib.parse.quote used as HTML escape, html.escape missing before Markup, json.dumps into HTML
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-838-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-838`
- [ ] `CWE-908` — Use of Uninitialized Resource — tier B — use before assign of optional globals, None dereference patterns hard in pure source
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-908-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-908`
- [ ] `CWE-909` — Missing Initialization of Resource — tier B — global db/engine used without nil/None check
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-909-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-909`
- [ ] `CWE-910` — Use of Expired File Descriptor — tier B — f.close(); f.read/write, double close after with-block misuse
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-910-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-910`
- [ ] `CWE-911` — Improper Update of Reference Count — tier B — ctypes/manual refcount, Py_INCREF style via cffi
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-911-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-911`
- [ ] `CWE-920` — Improper Restriction of Power Consumption — tier B — tight busy loops without sleep, crypto mining-like patterns
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-920-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-920`
- [ ] `CWE-939` — Improper Authorization in Handler for Custom URL Scheme — tier B — custom protocol handlers, webbrowser / desktop deep links
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-939-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-939`
- [ ] `CWE-1007` — Insufficient Visual Distinction of Homoglyphs Presented to User — tier B — username display without NFKC/confusable check
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-1007-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-1007`
- [ ] `CWE-1021` — Improper Restriction of Rendered UI Layers or Frames — tier B — missing X-Frame-Options / CSP frame-ancestors in Flask/Django responses, clickjacking
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-1021-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-1021`
- [ ] `CWE-1046` — Creation of Immutable Text Using String Concatenation — tier B — s += in loop
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-1046-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-1046`
- [ ] `CWE-1050` — Excessive Platform Resource Consumption within a Loop — tier B — open/connect/sleep inside for/while without bound
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-1050-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-1050`
- [ ] `CWE-1060` — Excessive Number of Inefficient Server-Side Data Accesses — tier B — N+1 ORM in loop: for x in qs: x.related
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-1060-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-1060`
- [ ] `CWE-1067` — Excessive Execution of Sequential Searches of Data Resource — tier B — LIKE %term%, filter(field__contains=), icontains without prefix
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-1067-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-1067`
- [ ] `CWE-1071` — Empty Code Block — tier B — except: pass, except Exception: pass, empty def/if body
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-1071-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-1071`
- [ ] `CWE-1072` — Data Resource Access without Use of Connection Pooling — tier B — psycopg2.connect per request without pool, create_engine without pool_size
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-1072-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-1072`
- [ ] `CWE-1084` — Invokable Control Element with Excessive File or Data Access Operations — tier B — many open/execute in one function
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-1084-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-1084`
- [ ] `CWE-1104` — Use of Unmaintained Third Party Components — tier B — requirements/import of known-abandoned packages
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-1104-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-1104`
- [ ] `CWE-1106` — Insufficient Use of Symbolic Constants — tier B — magic numbers in security-sensitive sites
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-1106-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-1106`
- [ ] `CWE-1108` — Excessive Reliance on Global Variables — tier B — global keyword heavy use, module mutables as request state
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-1108-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-1108`
- [ ] `CWE-1121` — Excessive McCabe Cyclomatic Complexity — tier B — branch count per function
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-1121-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-1121`
- [ ] `CWE-1123` — Excessive Use of Self-Modifying Code — tier B — types.FunctionType, code object rewrite, heavy monkeypatch
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-1123-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-1123`
- [ ] `CWE-1124` — Excessively Deep Nesting — tier B — indent depth > N
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-1124-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-1124`
- [ ] `CWE-1220` — Insufficient Granularity of Access Control — tier B — fetch by id without owner check
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-1220-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-1220`
- [ ] `CWE-1265` — Unintended Reentrant Invocation of Non-reentrant Code Via Nested Calls — tier B — threading.Lock re-entry patterns
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-1265-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-1265`
- [ ] `CWE-1284` — Improper Validation of Specified Quantity in Input — tier B — int(request) used as size/limit without bounds
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-1284-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-1284`
- [ ] `CWE-1285` — Improper Validation of Specified Index, Position, or Offset in Input — tier B — list[int(user)] without bounds
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-1285-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-1285`
- [ ] `CWE-1287` — Improper Validation of Specified Type of Input — tier B — no isinstance / schema type checks
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-1287-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-1287`
- [ ] `CWE-1288` — Improper Validation of Consistency within Input — tier B — count field vs list length unchecked
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-1288-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-1288`
- [ ] `CWE-1322` — Use of Blocking Code in Single-threaded, Non-blocking Context — tier B — time.sleep in async def, requests.get inside async without to_thread, open() sync in asyncio path
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-1322-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-1322`
- [ ] `CWE-1339` — Insufficient Precision or Accuracy of a Real Number — tier B — float money/price/balance arithmetic
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-1339-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-1339`
- [ ] `CWE-1341` — Multiple Releases of Same Resource or Handle — tier B — f.close(); f.close(), conn.close twice
  - [ ] unit hit/miss
  - [ ] `tests/fixtures/python/cwe/CWE-1341-vulnerable.txt` + `-safe.txt`
  - [ ] matrix green for `CWE-1341`

## Phase 2: Batch validation

- [ ] `gofmt -w` on touched files
- [ ] `make lint`
- [ ] `make test`
- [ ] `go test ./internal/lang/python/detectors/cwe/ -count=1` (individual unit tests)
- [ ] `go test ./tests/integration/python/ -count=1` **or** `make integration-python` (CWE matrix over all `python/cwe` pairs)
- [ ] Confirm every batch ID has both `CWE-N-vulnerable.txt` and `CWE-N-safe.txt` under `tests/fixtures/python/cwe/`
- [ ] Fixture count: `DiscoverPythonCWECases` includes all new IDs (pair discovery, not a manual allowlist)
- [ ] Update `_inventory.json`: move batch IDs from `missing` → `implemented`
- [ ] Update this ledger statuses to `[x]` with evidence
- [ ] Package files still ≤2000 lines (split if not)


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

- [ ] Unit hit/miss for `CWE-N`
- [ ] `tests/fixtures/python/cwe/CWE-N-vulnerable.txt`
- [ ] `tests/fixtures/python/cwe/CWE-N-safe.txt`
- [ ] Matrix auto-discovers pair; vulnerable asserts finding; safe asserts absence
- [ ] `make lint` + `make test` + `make integration-python` green before PR merge

## Dependencies

| Depends on | Note |
|------------|------|
| batch-00 | Framework + priority rules already present |
| Catalogue chunks | IDs must exist in `ruleset/python/chunks/` |
| Parent README | ownership + PR policy |

