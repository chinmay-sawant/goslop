# Batch 02 — Flask remaining (`BP-PY-18`, `19`, `20`)

> **Parent:** `plans/v0.0.2/heuristics/bp-plans/README.md` — v0.0.2 BP-PY remaining heuristics  
> **Epic / issue:** [#51](https://github.com/chinmay-sawant/goslop/issues/51) · [#53](https://github.com/chinmay-sawant/goslop/issues/53) expansion  
> **Status:** planning  
> **Estimated effort:** 1 PR (`feat/python-bp-batch-02-flask` or similar)  
> **PR policy:** one PR for this batch only — Flask IDs only (`18`–`20`); `16`/`17` already shipped

---

## Architecture constraints

| Constraint | Detail |
|------------|--------|
| Package | `internal/lang/python/detectors/bad_practices/` only |
| Registration | `RegisterRule("BP-PY-*", detect…)` from `init()` |
| Scan | Existing `PythonBadPracticeScan` — no new detector type |
| Detection | Pure-Go source patterns + `bpFacts` (mirror `detectBPPY16` / `17`) |
| Language | `LanguagePython` |
| Plugin | **Do NOT invent a second plugin.** Wire remains `internal/lang/python/detectors/all.go` → `badpractices.NewPythonBadPracticeScan()` |
| IDs | Always `BP-PY-*`; metadata from `ruleset/python/bad-practices.json` |
| **File size policy** | Target max **1500** lines / hard max **2000** per Go source file. If `rules_framework.go` would exceed 1500 after this batch, **split** Flask rules into `rules_flask.go` (Django stays in framework or moves later to `rules_django.go`). |
| Validation | `make lint` + `make test` unchecked until green |

### File placement (pre-implementation)

| File | Current lines (inventory) | Plan |
|------|--------------------------:|------|
| `rules_framework.go` | 200 | **Default** extend: already hosts Flask `16`/`17` + Django `21`. Adding 18–20 keeps Flask co-located. |
| `rules_flask.go` | — (new) | **If** post-batch `rules_framework.go` projected >1500 **or** framework file becomes hard to review — move **all** Flask rules (`16`–`20`) here and leave Django `21` in framework (or later batch-03 `rules_django.go`) |
| `scan_test.go` | 296 | Add tests for 18–20; split test file if approaching 1500 |
| `common.go` | 446 | Only shared Flask helpers if truly reused |
| `all.go` | stable | no second plugin |

**Already implemented (do not re-implement):**

| ID | Name | File |
|----|------|------|
| `BP-PY-16` | Flask DEBUG True In Production Code | `rules_framework.go` |
| `BP-PY-17` | Flask SECRET_KEY Hardcoded | `rules_framework.go` |

**This batch:**

| ID | Name | Severity | Category |
|----|------|----------|----------|
| `BP-PY-18` | Flask Route Missing Methods Restriction | low | Flask |
| `BP-PY-19` | Flask jsonify Error Leaks Exception | medium | Flask |
| `BP-PY-20` | Flask send_file User Path | high | Flask |

---

## Overview

| Rule | detection_notes (catalogue) |
|------|-----------------------------|
| `BP-PY-18` | Heuristic: route handlers that call `request.form` / `get_json` / POST-like logic without `methods=['POST',...]` or methods list. Flag CSRF-sensitive form posts especially. |
| `BP-PY-19` | Match Flask `errorhandler` / `register_error_handler` bodies returning `str(exc)`, `traceback.format_exc()`, or `repr(e)` in JSON/HTML responses. |
| `BP-PY-20` | Match `flask.send_file` / `send_from_directory` where path/filename comes from `request.args` / form / view args without `safe_join` + root check. Complements CWE-22. |

Historical note: `python-heuristics-bp.md` Phase 4 listed `BP-PY-20` as `[~]` deferred and 18/19 as lower signal — this batch **owns** shipping them with explicit heuristics and tests.

---

## Executive Summary

Complete the Flask surface after DEBUG/SECRET_KEY. Prefer extending `rules_framework.go` while under budget; split to `rules_flask.go` rather than exceeding 1500 lines. Keep Django work for batch 03.

---

## Phase 1: Placement + budget

### 1.1 File budget check (before edit)

- [ ] Baseline: `wc -l internal/lang/python/detectors/bad_practices/rules_framework.go scan_test.go`
- [ ] Decision:
  - [ ] Extend `rules_framework.go` if post-batch estimate ≤1500
  - [ ] Else create `rules_flask.go`, relocate Flask `init` registrations (`16`–`20`) carefully (single registration per id)
- [ ] Do not add Django IDs in this PR

### 1.2 Flask context helpers

- [ ] Reuse `looksFlaskish` / path helpers already in `rules_framework.go` where useful
- [ ] Reuse `callArgsRegion`, `pushAt`, `isPythonTestFile` from package
- [ ] Skip test files for production-flavored Flask findings (consistent with BP-PY-16)

---

## Phase 2: `BP-PY-18` Flask Route Missing Methods Restriction

### 2.1 Register

- [ ] `RegisterRule("BP-PY-18", detectBPPY18)` in framework or flask rules `init()`

### 2.2 Detect heuristic

Cite **detection_notes:** route handlers using `request.form` / `get_json` / POST-like logic without `methods=['POST', ...]` (or methods list). Flag CSRF-sensitive form posts especially.

- [ ] Implement `detectBPPY18` as **heuristic** (document low severity / false-positive rate):
  - Find `@app.route` / `@blueprint.route` / `add_url_rule` nearby route defs
  - If decorator/call has no `methods=` **or** methods is default GET-only
  - And function body references `request.form` / `request.get_json` / `request.json` / `request.files` (POST-like)
  - Then flag
  - Miss: `@app.route(..., methods=["POST"])` or `methods=['GET', 'POST']` with form use
  - Miss: GET-only handlers that only use `request.args`
  - v0 may use same-function window rather than full CFG
- [ ] Optional: only flag when `request.form` present (stronger CSRF signal) — document if narrowed

### 2.3 Hit / miss tests

- [ ] Hit:
  ```python
  @app.route("/login")
  def login():
      user = request.form["user"]
      return user
  ```
- [ ] Miss:
  ```python
  @app.route("/login", methods=["POST"])
  def login():
      user = request.form["user"]
      return user
  ```
- [ ] Miss: route using only `request.args` without form/json
- [ ] Want-list includes `BP-PY-18`

### 2.4 Proof

- [ ] Unit test green; severity low from metadata

---

## Phase 3: `BP-PY-19` Flask jsonify Error Leaks Exception

### 3.1 Register

- [ ] `RegisterRule("BP-PY-19", detectBPPY19)`

### 3.2 Detect heuristic

Cite **detection_notes:** `errorhandler` / `register_error_handler` bodies returning `str(exc)`, `traceback.format_exc()`, or `repr(e)` in JSON/HTML responses.

- [ ] Implement `detectBPPY19`:
  - Needle: `errorhandler`, `register_error_handler`, `traceback`, `jsonify`
  - Match decorated error handlers or `app.register_error_handler` callbacks
  - Flag body patterns: `str(e)` / `str(exc)` / `repr(e)` / `traceback.format_exc()` returned or passed to `jsonify` / response
  - Miss: generic message constants (`jsonify({"error": "internal"})`) without exception text
  - Miss: logging exception without returning traceback to client
- [ ] Medium severity from metadata

### 3.3 Hit / miss tests

- [ ] Hit:
  ```python
  @app.errorhandler(Exception)
  def handle(e):
      return jsonify(error=str(e)), 500
  ```
- [ ] Hit: `return traceback.format_exc()` in errorhandler
- [ ] Miss: `return jsonify(error="internal"), 500` without `str(e)`
- [ ] Want-list includes `BP-PY-19`

### 3.4 Proof

- [ ] Unit test green

---

## Phase 4: `BP-PY-20` Flask send_file User Path

### 4.1 Register

- [ ] `RegisterRule("BP-PY-20", detectBPPY20)`

### 4.2 Detect heuristic

Cite **detection_notes:** `flask.send_file` / `send_from_directory` where path/filename comes from `request.args` / form / view args without `safe_join` + root check. Complements CWE-22.

- [ ] Implement `detectBPPY20`:
  - Needle: `send_file`, `send_from_directory`
  - Hit when first path/filename argument is clearly request-derived: `request.args[...]`, `request.args.get`, `request.form`, `request.files` name, or path param name flowing in simple assignment window
  - Miss: constant path literals
  - Miss: `send_from_directory(SAFE_ROOT, safe_join(...))` / explicit `safe_join` on user segment (best-effort)
  - High severity from metadata
  - Align message with path-traversal hygiene (no need to implement full CWE-22)
- [ ] Skip test fixtures that intentionally demo attacks only if path is under tests — prefer still flagging demo code unless `isPythonTestFile`

### 4.3 Hit / miss tests

- [ ] Hit:
  ```python
  from flask import send_file, request
  @app.route("/dl")
  def dl():
      return send_file(request.args["path"])
  ```
- [ ] Hit: `send_from_directory("/var/data", request.args.get("f"))` without safe_join
- [ ] Miss: `send_file("/var/data/report.pdf")`
- [ ] Miss (if implemented): `send_from_directory(ROOT, safe_join(ROOT, name))` pattern
- [ ] Want-list includes `BP-PY-20`
- [ ] Assert severity **high** on a hit finding

### 4.4 Proof

- [ ] Unit test green; high severity asserted

---

## Phase 5: Registration surface + size check

### 5.1 Catalogue registration tests

- [ ] `TestBPRulesRegistered` want-list includes `BP-PY-18`, `BP-PY-19`, `BP-PY-20` (and still `16`, `17`, `21` + prior)
- [ ] No bare `BP-*` IDs
- [ ] `MetadataFor` non-nil; pack bad-practice
- [ ] Confirm `all.go` unchanged (still one BP scan)

### 5.2 File size check after batch

- [ ] `wc -l internal/lang/python/detectors/bad_practices/*.go`
- [ ] No file > **2000** (hard max)
- [ ] No file > **1500** (target); if `rules_framework.go` exceeds target, **split to `rules_flask.go` before merge**
- [ ] Record line counts in PR body

### 5.3 Interaction with prior Flask rules

- [ ] Existing `TestBPPY16FlaskDebug` / `TestBPPY17FlaskSecretKey` still pass
- [ ] New detectors do not double-report DEBUG/SECRET under 18–20 IDs

---

## Phase 6: Validation gates (batch PR)

> Per skill: leave unchecked until green on implement branch.

- [ ] `gofmt -w` on all touched Go files
- [ ] `make lint` — unchecked until green  
  **Evidence:** _(command + date + branch)_
- [ ] `make test` — unchecked until green  
  **Evidence:** _(command + date + branch)_
- [ ] Optional: `go test ./internal/lang/python/detectors/bad_practices/ -count=1`
- [ ] PR: `Relates to #53` / `Relates to #51`
- [ ] Update parent [README.md](./README.md) Batch 02 rollup to `[x]` only after gates green
- [ ] Optionally mark historical `[~]` rows for 18–20 in `python-heuristics-bp.md` as deferred-to this ledger with pointer

---

## Dependencies

| Depends on | Notes |
|------------|--------|
| Batch 01 | **Not required** (different files preferred); may land in either order |
| Existing Flask 16/17 detectors | patterns/helpers to mirror |
| CWE-22 Python detector | complementary only; BP-PY-20 stays BP id |
| Batch 03 Django | do not implement 22–28 here |

## Out of scope

- Django `BP-PY-21` changes (already shipped) / 22–28
- FastAPI FileResponse (`BP-PY-32`) — batch 04
- Full CSRF framework analysis beyond methods restriction heuristic

---

## References

- Inventory: `_inventory.json` rules `BP-PY-18`, `19`, `20`
- Existing: `internal/lang/python/detectors/bad_practices/rules_framework.go`
- Skill: `plans/skills/phase-wise-checklist/SKILLS.md`
- Parent: `plans/v0.0.2/heuristics/bp-plans/README.md`
