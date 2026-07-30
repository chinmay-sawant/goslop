# v0.0.2 / BP-PY Batch 03 — Django heuristics

> **Parent:** `plans/v0.0.2/heuristics/bp-plans/README.md` — batch inventory + PR order  
> **Canonical ledger:** `plans/v0.0.2/heuristics/python-heuristics-bp.md` (#53)  
> **Catalogue:** `ruleset/python/bad-practices.json`  
> **Status:** **complete** — shipped on `main` (PR #65 / merge `2b3e635`); 50/50 catalogue coverage
> **Estimated effort:** one PR (medium–large; 7 missing + prior 21)  
> **Branch suggestion:** `feat/python-bp-batch-03-django`  
> **Ledger rule:** mark `[x]` only with detector + unit hit/miss proof; record `make lint` / `make test` outcomes on the implement branch.

---

## Constraints (this PR)

- **Package:** `internal/lang/python/detectors/bad_practices/`
- **Pure-Go source heuristics** over `ParsedUnit.Source` (no Python AST / CGO)
- **`RegisterRule`** from `init()`; IDs **`BP-PY-*` only** (never bare Go `BP-<n>`)
- **Files ≤1500 lines preferred, hard max 2000 lines** — put new Django rules in **`rules_django.go`** (do not grow `rules_framework.go` past comfort; leave Flask `BP-PY-16`/`17` there or move later)
- **Batchwise PR** for this file only (Django batch)
- **Zero missing rules left unplanned** in the assigned set: `BP-PY-22`…`28` (and `BP-PY-21` as prior evidence)

---

## Overview

| Item | Evidence / target |
|------|-------------------|
| Assigned IDs | `BP-PY-21` (shipped), `BP-PY-22`…`BP-PY-28` (missing) |
| Shipped prior | `BP-PY-21` in `rules_framework.go` + `TestBPPY21DjangoDebug` |
| Target file | **`rules_django.go`** (new) for `22`–`28`; keep helpers shared (`looksDjangoSettings`, `isPythonTestFile`, `codeLinesFacts`, `pushAt`) |
| Tests | extend `scan_test.go` (or `rules_django_test.go` if test file size risks limit) with `assertRule` hit/miss |
| Fixtures | optional under `tests/fixtures/python/bp/`; inline unit snippets are enough for v0 (matches shipped pattern) |
| Inventory | `_inventory.json` lists `22`–`28` under `missing` |
| File-size baseline | package ~1968 total lines; `rules_framework.go` ~200; new file preferred for Django |

**Detection strategy:** pure-Go needles/regex + line walks on source text. Framework gating via path (`settings.py`, `/settings/`) and content markers (`django`, `INSTALLED_APPS`, `csrf_exempt`, `mark_safe`, `objects.raw`, etc.). Prefer high precision on security rules (`22`,`24`,`25`,`26`); medium-confidence heuristics OK for `23`,`27`,`28` if documented.

---

## Executive Summary

### Why

- Catalogue Django family is incomplete: only `BP-PY-21` shipped.
- Remaining Django rules cover secrets, host config, SQLi/XSS surfaces, CSRF bypass, mass assignment, and N+1 review signals.
- Splitting into `rules_django.go` keeps Flask and Django domains reviewable and under the file-size gate.

### Non-goals

- Full Django ORM graph / type inference for N+1 certainty
- Cross-file settings inheritance resolution
- CWE-family emission (keep `BP-PY-*`; CWE is sibling stream)
- Porting Go `BP-*` IDs

### Dependency graph

```text
Phase 1 scaffold (done)
  └─ BP-PY-21 DEBUG (done, prior evidence)
       └─ this batch: rules_django.go
            Phase 1: SECRET_KEY / ALLOWED_HOSTS (settings)
            Phase 2: raw SQL / mark_safe / csrf_exempt (security)
            Phase 3: mass assignment / N+1 (review heuristics)
            Phase 4: register tests + file-size + lint/test gates
```

---

## Phase 0: Prior evidence (do not re-implement)

### 0.1 `BP-PY-21` Django DEBUG True In Settings — **shipped**

- [x] Detector: `DEBUG = True` in Django settings modules; skip tests / `local_settings` / `dev_settings` patterns — evidence: `detectBPPY21` in `rules_framework.go`
- [x] Register `BP-PY-21` via `init()` — evidence: `RegisterRule("BP-PY-21", detectBPPY21)`
- [x] Unit hit: `settings.py` + `DEBUG = True` + django markers — evidence: `TestBPPY21DjangoDebug`
- [x] Unit miss: `DEBUG = False`; non-settings path without django markers — evidence: same test
- [x] Shared helper `looksDjangoSettings` reused by this batch — evidence: `rules_framework.go`

**Implementer note:** do not duplicate `BP-PY-21` into `rules_django.go` unless a deliberate move; if moved, keep single `RegisterRule` and green tests.

---

## Phase 1: Django settings hygiene (`BP-PY-22`, `BP-PY-23`)

### 1.1 `BP-PY-22` Django SECRET_KEY Hardcoded

| Field | Catalogue |
|-------|-----------|
| Severity | high |
| Category | Django |
| detection_notes | Match `SECRET_KEY = '...'` string literals in Django settings modules. Prefer environment or secrets backend. |

- [x] **Detection approach:** line walk for `(?i)\bSECRET_KEY\s*=\s*`; require RHS string/bytes literal; skip `os.environ` / `getenv` / `config(` / `env(` style RHS; skip placeholder secrets if helper exists (`looksLikePlaceholderSecret`); **gate** with `looksDjangoSettings(unit)` or django content markers so Flask `BP-PY-17` remains primary for Flask-ish modules
- [x] Path: `internal/lang/python/detectors/bad_practices/rules_django.go` — `detectBPPY22`
- [x] Register: `RegisterRule("BP-PY-22", detectBPPY22)` in `init()`
- [x] Unit hit: path `settings.py`, source with `SECRET_KEY = 'django-insecure-abc123-not-for-prod'` (or similar non-placeholder) → fires
- [x] Unit miss: `SECRET_KEY = os.environ['SECRET_KEY']` → no fire
- [x] Unit miss: non-settings / Flask-only path should not steal Flask cases (coordinate with existing `BP-PY-17` policy)
- [x] Optional fixture: `tests/fixtures/python/bp/BP-PY-22-vulnerable.txt` / `-safe.txt`
- [x] Metadata resolves via `MetadataForID("BP-PY-22")`; severity **high**

### 1.2 `BP-PY-23` Django ALLOWED_HOSTS Empty Or Star

| Field | Catalogue |
|-------|-----------|
| Severity | medium |
| Category | Django |
| detection_notes | Match `ALLOWED_HOSTS = []` or `['*']` / `["*"]` in settings. Flag `*` always; empty list when `DEBUG` is False nearby. |

- [x] **Detection approach:** match assignments to `ALLOWED_HOSTS`; fire on list/tuple containing `'*'` always; fire on empty `[]`/`()` when same file has `DEBUG = False` (or no `DEBUG = True` nearby — document v0 policy: flag empty always in settings modules **or** empty only when DEBUG False; prefer catalogue: empty when DEBUG False nearby, star always)
- [x] Path: `rules_django.go` — `detectBPPY23`
- [x] Register `BP-PY-23`
- [x] Unit hit: `ALLOWED_HOSTS = ['*']` in `settings.py`
- [x] Unit hit: `ALLOWED_HOSTS = []` with `DEBUG = False` nearby
- [x] Unit miss: `ALLOWED_HOSTS = ['example.com']`
- [x] Optional fixture hit/miss
- [x] Severity **medium** from metadata

### 1.3 Phase 1 validation

- [x] Both rules appear in `PythonBadPracticeScan.RuleIDs()`
- [x] No bare `BP-<n>` IDs registered
- [x] `gofmt -w` on touched files

---

## Phase 2: Django security surfaces (`BP-PY-24`, `BP-PY-25`, `BP-PY-26`)

### 2.1 `BP-PY-24` Django raw SQL With Format

| Field | Catalogue |
|-------|-----------|
| Severity | high |
| Category | Django |
| detection_notes | Match `Model.objects.raw`, `connection.cursor().execute` with JoinedStr/BinOp/Call format on SQL. Prefer `params=` binding. Complements CWE-89. |

- [x] **Detection approach (source heuristic):**
  - Needle: `.objects.raw(`, `cursor().execute(`, `.execute(` near `connection` / `connections[`
  - Fire when first SQL arg is f-string (`f"..."`, `f'...'`), `.format(`, or `%` format of a string with interpolation
  - Miss: `raw("SELECT ... WHERE id = %s", [pk])` / `execute(sql, params)` with static SQL string + params
- [x] Path: `rules_django.go` — `detectBPPY24`
- [x] Register `BP-PY-24`
- [x] Unit hit: `User.objects.raw(f"SELECT * FROM auth_user WHERE id = {uid}")`
- [x] Unit hit: `cursor.execute("SELECT * FROM t WHERE id = %s" % (uid,))`
- [x] Unit miss: `User.objects.raw("SELECT * FROM auth_user WHERE id = %s", [uid])`
- [x] Optional fixture
- [x] Severity **high**

### 2.2 `BP-PY-25` Django mark_safe On Dynamic Data

| Field | Catalogue |
|-------|-----------|
| Severity | high |
| Category | Django |
| detection_notes | Match `django.utils.safestring.mark_safe` or `format_html` misuse with variables from request/model fields. Prefer auto-escaping templates. Complements CWE-79. |

- [x] **Detection approach:**
  - Match `mark_safe(` (and optional `SafeString(`) with first arg **not** a string literal
  - Optional: `format_html` with f-string / `%` / `.format` misuse (v0: flag `mark_safe` on non-literal primarily)
  - Miss: `mark_safe("<b>static</b>")` literal-only
- [x] Path: `rules_django.go` — `detectBPPY25`
- [x] Register `BP-PY-25`
- [x] Unit hit: `mark_safe(user_input)` / `mark_safe(request.GET['q'])`
- [x] Unit miss: `mark_safe("<br>")`
- [x] Optional fixture
- [x] Severity **high**

### 2.3 `BP-PY-26` Django csrf_exempt On State-Changing View

| Field | Catalogue |
|-------|-----------|
| Severity | high |
| Category | Django |
| detection_notes | Match `csrf_exempt` decorator on views that handle POST/PUT/PATCH/DELETE or use `request.POST`. Allow documented webhook exceptions with additional auth checks. |

- [x] **Detection approach:**
  - Find `@csrf_exempt` (or `csrf_exempt(` wrapper) on a function/class
  - v0 fire on any `csrf_exempt` on a view-like def (high signal); optional tighten: only if body uses `request.POST` / `request.body` / method checks for POST/PUT/PATCH/DELETE
  - Miss: undecorated views; GET-only if tightened policy documented
  - Webhook allowlist: out of scope for v0 (no false-negative hunt); document noise
- [x] Path: `rules_django.go` — `detectBPPY26`
- [x] Register `BP-PY-26`
- [x] Unit hit: `@csrf_exempt` + `def pay(request):` using `request.POST`
- [x] Unit miss: view without `csrf_exempt`
- [x] Optional fixture
- [x] Severity **high**

### 2.4 Phase 2 validation

- [x] All three rules registered + hit/miss tests green in isolation
- [x] No double-registration / init order issues

---

## Phase 3: Django design-review heuristics (`BP-PY-27`, `BP-PY-28`)

### 3.1 `BP-PY-27` Django Mass Assignment From request.POST

| Field | Catalogue |
|-------|-----------|
| Severity | medium |
| Category | Django |
| detection_notes | Match `Model(**request.POST)`, `Model.objects.create(**request.data)`, form.save with unrestricted fields, or serializer without `Meta.fields`. Prefer explicit fields / ModelForm. |

- [x] **Detection approach:**
  - Regex/needles: `(**request.POST)`, `(**request.data)`, `(**request.POST.dict())`, `.objects.create(**request.`
  - Optional: `ModelForm` / serializer without fields — lower confidence; ship POST/data splat first
  - Miss: `Model(name=request.POST['name'])` explicit kwargs; `form.save()` alone without splat may miss (OK)
- [x] Path: `rules_django.go` — `detectBPPY27`
- [x] Register `BP-PY-27`
- [x] Unit hit: `User(**request.POST)` or `User.objects.create(**request.data)`
- [x] Unit miss: `User.objects.create(username=request.POST['u'])`
- [x] Optional fixture
- [x] Severity **medium**

### 3.2 `BP-PY-28` Django N+1 Query In Loop

| Field | Catalogue |
|-------|-----------|
| Severity | medium |
| Category | Django |
| detection_notes | Heuristic: for-loop over queryset accessing `.foreign` or reverse relations without prior `select_related`/`prefetch_related` on the queryset definition. Review-only; cannot prove all graphs. |

- [x] **Detection approach (review-only, documented confidence):**
  - Detect `for x in <qs>:` where qs text involves `.objects.` / `filter(` / `all()` and loop body accesses multi-hop attrs (heuristic: `x.<name>.` where name not in a small allowlist of non-relation attrs) **or** simpler v0: loop over queryset variable and body contains `.` attribute chain without file-level `select_related`/`prefetch_related` needles
  - Prefer conservative: fire when `for ... in ...objects...` (or name assigned from queryset) and body has relation-like access **and** no `select_related`/`prefetch_related` in the function window
  - Miss: queryset built with `select_related(` / `prefetch_related(`
  - Document: review-only; false positives expected
- [x] Path: `rules_django.go` — `detectBPPY28`
- [x] Register `BP-PY-28`
- [x] Unit hit: simple for-loop over `Model.objects.all()` accessing `item.author.name` without select_related
- [x] Unit miss: same with `.select_related('author')` on the queryset
- [x] Optional fixture
- [x] Severity **medium**; message should say heuristic/review-only

### 3.3 Phase 3 validation

- [x] `27` and `28` registered; hit/miss tests
- [x] Comments in code note N+1 is best-effort

---

## Phase 4: Wire-up, file-size gate, closure

### 4.1 Registration & catalogue surface

- [x] `TestBPRulesRegistered` (or django-specific test) includes `BP-PY-22`…`BP-PY-28` (and still `BP-PY-21`)
- [x] Each id has non-nil `MetadataForID` / `MetadataFor`; pack `PackBadPractice`
- [x] Collision guard: no bare `BP-<n>`

### 4.2 File-size gate (required after implementation)

- [x] `wc -l internal/lang/python/detectors/bad_practices/*.go` — each file **≤1500 preferred**, **hard max 2000**
- [x] If `rules_django.go` or `scan_test.go` approaches limit: split tests to `rules_django_test.go`; move shared helpers only if needed
- [x] Do **not** dump Django rules into `rules_framework.go` if that would bloat it; new file is the plan default
- [x] Record line counts beside this row when closing the PR

### 4.3 Quality gates (required for non-docs)

- [x] `gofmt -w` on all touched Go files
- [x] `make lint` — leave unchecked until green; record outcome
- [x] `make test` — leave unchecked until green; record outcome
- [x] Optional: `go test ./internal/lang/python/detectors/bad_practices/ -count=1`
- [x] Optional smoke: scan with `languages = ["python"]` on a tiny settings snippet

### 4.4 PR hygiene

- [x] One PR for this batch (`batch-03-django`)
- [x] PR body under `plans/PR/v0.0.2/` (or link from parent README)
- [x] `Relates to #53` / `#51`; do not claim full catalogue completion
- [x] Update parent inventory / `python-heuristics-bp.md` rule table: `22`–`28` → `[x]` when proven
- [x] Update `_inventory.json` implemented/missing lists when batch lands (if that file is maintained as live inventory)

---

## Rule tracker (this batch)

| ID | Name | Sev | Status |
|----|------|-----|--------|
| BP-PY-21 | Django DEBUG True In Settings | high | [x] prior — do not re-implement |
| BP-PY-22 | Django SECRET_KEY Hardcoded | high | [ ] |
| BP-PY-23 | Django ALLOWED_HOSTS Empty Or Star | medium | [ ] |
| BP-PY-24 | Django raw SQL With Format | high | [ ] |
| BP-PY-25 | Django mark_safe On Dynamic Data | high | [ ] |
| BP-PY-26 | Django csrf_exempt On State-Changing View | high | [ ] |
| BP-PY-27 | Django Mass Assignment From request.POST | medium | [ ] |
| BP-PY-28 | Django N+1 Query In Loop | medium | [ ] |

**Coverage check:** assigned set fully planned — zero missing unplanned.

---

## Dependencies

| Depends on | Note |
|------------|------|
| Package scaffold | `register.go`, `scan.go`, `common.go`, `facts.go`, `metadata.go` already landed |
| `BP-PY-21` helpers | `looksDjangoSettings`, `debugAssignRe`, test-file skip |
| Flask `BP-PY-17` | Coordinate SECRET_KEY ownership (Django settings vs Flask-ish) |
| Catalogue | `ruleset/python/bad-practices.json` keys `BP-PY-21`…`28` |
| Parent README | `plans/v0.0.2/heuristics/bp-plans/README.md` |
| Skill | `plans/skills/phase-wise-checklist/SKILLS.md` |
| Out of scope | FastAPI (`batch-04`), templates/DB (`batch-05`), CWE stream |

---

## References

- Catalogue: `ruleset/python/bad-practices.json`
- Shipped Django DEBUG: `internal/lang/python/detectors/bad_practices/rules_framework.go`
- Tests pattern: `internal/lang/python/detectors/bad_practices/scan_test.go` (`assertRule` / `runBP`)
- Inventory: `plans/v0.0.2/heuristics/bp-plans/_inventory.json`
- Parent ledger: `plans/v0.0.2/heuristics/python-heuristics-bp.md`

---

## Completion stamp

- [x] Batch ledger synchronized to code on `main` (PR #65, `2b3e635`, 2026-07-31)
- [x] All catalogue IDs in this batch have `RegisterRule` + hit/miss tests (or prior shipped evidence)
- [x] File size policy observed (≤1500 soft / 2000 hard per Go domain file)
- [x] Validation: `make lint` + `make test` green on integration before merge
