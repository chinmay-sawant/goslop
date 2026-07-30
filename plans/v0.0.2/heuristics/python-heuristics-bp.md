# v0.0.2 / #53 — Python BP-PY bad-practice heuristics

> **Parent:** `plans/v0.0.2/heuristics/python-heuristics.md` — epic #51 rollup (BP stream)  
> **Issue:** [#53](https://github.com/chinmay-sawant/goslop/issues/53) — `python(bp): implement BP-PY bad-practice heuristics from catalogue`  
> **Status:** priority subset **shipped** (14 rules via PR #58); batches 06–07 add **10** more (38–47); remaining **26** planned as batchwise PRs under `bp-plans/`  
> **Estimated effort:** multi-PR batches 01–08  
> **Expansion ledgers:** [`bp-plans/README.md`](./bp-plans/README.md) (all remaining IDs; **no catalogue gap**)  
> **Ledger rule:** Mark `[x]` only with evidence; for code batches run `make lint` + `make test`.  
> **File size:** prefer ≤**1500** lines per Go file; **hard max 2000** — split domain `rules_*.go` files instead of growing past the cap.

---

## Overview

| Item | Current evidence |
|------|------------------|
| Catalogue | `ruleset/python/bad-practices.json` — **50** map keys `BP-PY-1` … `BP-PY-50` |
| Implemented (`RegisterRule`) | **14:** 1,2,4,6,7,8–13,16,17,21 — see `bp-plans/batch-00-shipped.md` |
| Missing (planned) | **36:** 3,5,14,15,18–20,22–50 — see `bp-plans/batch-01` … `batch-08` |
| Package | `internal/lang/python/detectors/bad_practices/` — scan + rules_core/security/framework |
| Parse quality | source-only (`ParseQualitySourceOnly`); pure-Go source patterns |
| File size policy | ≤1500 preferred / **2000 hard max** per Go file (split domain files) |
| Inventory | `bp-plans/_inventory.json` |
| Batch index | [`bp-plans/README.md`](./bp-plans/README.md) |
| Parent epic | #51 |

**ID policy:** keep catalogue IDs as **`BP-PY-*`**. Do **not** emit Go `BP-*` for Python. Do **not** point Go `metadata_gen.go` at `ruleset/python/`.

**Detection strategy (v0):** pure-Go **source-pattern heuristics** over `ParsedUnit.Source` (mirror Go text fallbacks / needles). Full Python AST is out of scope for the first landings; refine later if a pure-Go parse path appears.

**Priority subsets (issue #53 / product signal):**

| Batch | Rules | Theme |
|------:|-------|-------|
| A — Core | `BP-PY-1`, `2`, `4`, `6`, `7` (+ optional `3`, `5`) | bare except / pass, mutable defaults, assert validation, `open` without `with` |
| B — Security hygiene | `BP-PY-8` … `13` | `shell=True`, `os.system`, pickle, `yaml.load`, eval/exec, hardcoded secrets |
| C — Framework high-signal | Flask `16`,`17`,`20`; Django `21`,`22`,`24`,`25`,`26`; FastAPI `30`,`32` | DEBUG/SECRET_KEY/`send_file`/raw SQL/`mark_safe`/`csrf_exempt`/blocking I/O/user path |
| D — Templates / DB | `BP-PY-33`, `35` (+ optional `34`,`37`) | Jinja2 autoescape off; SQLAlchemy `text()` f-strings |
| E — Rest | remaining of 50 | second+ PR; may ship as `[~]` deferred after A–D |

Avoid shipping all **50** in the first PR if size risks review; first PR success = **A+B+C+D** (or A+B minimum with C/D follow-up), registered, fixture-backed.

---

## Executive Summary

### Why this work

- Epic #39 / PR #50 shipped catalogues + plugin stub; `BP-PY-*` never run.
- Go already proves the product shape: one `*BadPracticeScan` detector, many registered rule fns, catalogue metadata, hit/miss fixtures.
- Python must stay **opt-in** (`languages = ["python"]`); default registry remains Go-only.

### Non-goals

- CWE / PERF families (siblings #52 / #54)
- Porting all **135** Go BP rules
- Full framework type inference or cross-file taint
- Tree-sitter / CGO Python parse requirement for first ship
- Claiming catalogue coverage without fixtures

### Dependency graph

```text
#39 foundation (done)
    └─ #53 BP heuristics
         Phase 1 scaffold + metadata + plugin wire
              ├─ Phase 2 core (A)
              ├─ Phase 3 security (B)
              ├─ Phase 4 framework + templates/DB (C+D)
              └─ Phase 5 remaining (E) optional / multi-PR
         Phase 6 closure: fixtures manifest, list-rules, lint/test
```

Sibling #52 (CWE) may land package layout under `internal/lang/python/detectors/` first — **share** package tree; do not fork a second plugin.

---

## Phase 1: Scaffold, metadata, plugin registration

### 1.1 Package layout (mirror Go, Python-local)

- [x] Create `internal/lang/python/detectors/` package root (or co-own if #52 already created it)
- [x] Create `internal/lang/python/detectors/bad_practices/` with files modeled on Go:
  - [x] `register.go` — `RegisterRule(id, fn)`, snapshot catalogue (`ruleEntry` + mutex or init-only append)
  - [x] `scan.go` — `PythonBadPracticeScan` implementing `core.Detector` (`Language() == LanguagePython`, `RuleIDs`, `Run`, lifecycle via `core.BaseDetector` unless project cache needed)
  - [x] `common.go` — finding push helpers, path display, test-file skip (`*_test.py`, `test_*.py`, `tests/`)
  - [x] `facts.go` — optional light fact bag (line index, needle index via `internal/ast.Build` if useful); **no** `go/ast` dependency
  - [x] `metadata.go` — load `BP-PY-*` → `rules.RuleMetadata` from Python catalogue (hand-written map **or** `//go:embed` of `ruleset/python/bad-practices.json` decoded at init)
- [x] Explicitly **do not** import or regenerate `internal/lang/go/detectors/bad_practices/metadata_gen.go`
- [x] Export constructor `NewPythonBadPracticeScan() *PythonBadPracticeScan` (name bikeshed OK; keep stable in tests)

### 1.2 Metadata contract

- [x] Every registered rule id resolves non-nil metadata via `MetadataFor` / `MetadataForID`
- [x] Severity string mapping matches catalogue (`info|low|medium|high|critical` → `rules.Severity*`)
- [x] `Pack` set to `rules.PackBadPractice` (or rely on `PackFromRuleID("BP-PY-…")` which already returns bad-practice)
- [x] Title/name from catalogue `name`; description from catalogue `description`; fix/detection notes optional in Fix field
- [x] Unit test: catalogue size ≥ 50 keys parseable; implemented subset has metadata for each registered id
- [x] Unit test: no registered id is bare `BP-<n>` (collision guard)

### 1.3 Plugin wire-up

- [x] Add `internal/lang/python/detectors/all.go` (or equivalent) returning `[]core.Detector{ badpractices.New…() }` (CWE may append later)
- [x] Change `internal/lang/python/plugin.go` `Detectors()` / `NewDetectors()` to return session detectors (same pattern as `internal/lang/go/plugin.go`)
- [x] Update package doc on `plugin.go`: no longer “zero detectors”
- [x] Update `internal/lang/python/plugin_test.go`: empty-catalogue assertions become **non-empty** when BP is registered; assert at least one detector with `BP-PY-*` in `RuleIDs()`
- [x] Keep `DefaultRegistry` Go-only; Python still via `NewRegistryWithLanguages` / config `languages`

### 1.4 Scan gating

- [x] `Run` early-returns when unit language ≠ Python or source empty
- [x] Respect `ScanContext.Allows(ruleID)` and `BadPracticesEnabled` (engine already filters; still skip disabled packs if detectors are invoked broadly)
- [x] Needle / cheap prefilter before expensive windows (copy Go BP style)
- [x] No hang on unbalanced braces/strings (add regression if scanners use brace matching)

### 1.5 Phase 1 validation

- [x] `gofmt -w` on touched Go files
- [x] `make lint` — green 2026-07-31 on feat/python-bp-heuristics
- [x] `make test` — green 2026-07-31 on feat/python-bp-heuristics
- [x] Proof: `go test ./internal/lang/python/...` shows registered BP detector without requiring all 50 rules yet

---

## Phase 2: Core language & error handling (Batch A)

Priority IDs from issue body. Implement + fixture per rule.

### 2.1 `BP-PY-1` Bare Except Clause

- [x] Detector: match `except:` or broad `except Exception` / `except BaseException` with weak handling (pass / bare continue / log-only without re-raise — start with `except:` and `except Exception:` + `pass`)
- [ ] Path: `internal/lang/python/detectors/bad_practices/` (e.g. `rules_core.go` / `rules_error.go`)
- [x] Register `BP-PY-1` in `init()`
- [~] Fixture hit: `tests/fixtures/python/bp/BP-PY-1-vulnerable.txt` — expects finding `BP-PY-1` (inline unit tests instead)
- [~] Fixture miss: `tests/fixtures/python/bp/BP-PY-1-safe.txt` — specific `except ValueError` / re-raise; **no** `BP-PY-1` (inline unit tests instead)
- [x] Test: table or `runFixtureRule`-style helper asserts true positive / true negative

### 2.2 `BP-PY-2` Except Pass

- [x] Detector: except suite is solely `pass` (optional comment)
- [x] Register `BP-PY-2`
- [~] Fixtures: `BP-PY-2-vulnerable.txt` / `BP-PY-2-safe.txt` (inline unit tests)
- [x] Proof: test green for hit/miss

### 2.3 `BP-PY-4` Mutable Default Argument

- [x] Detector: `def` / `async def` defaults that are `[]`, `{}`, `set()`, or bare list/dict/set literals in signature text
- [x] Register `BP-PY-4` (high severity per catalogue)
- [x] Fixtures hit/miss; miss uses `None` default pattern (inline tests)
- [x] Proof: severity from metadata is **high**

### 2.4 `BP-PY-6` assert Used For Runtime Validation

- [x] Detector: `assert` in non-test modules (skip `test_*.py`, `*_test.py`, `tests/`); optional needle for request/auth patterns later
- [x] Register `BP-PY-6`
- [x] Fixtures: library module hit; test file miss (inline)
- [x] Proof: test file does not fire

### 2.5 `BP-PY-7` open Without Context Manager

- [x] Detector: `open(` / `.open(` assigned or used without surrounding `with`
- [x] Register `BP-PY-7`
- [x] Fixtures: bare `f = open(...)` hit; `with open(...) as f` miss (inline)
- [x] Proof: hit/miss tests

### 2.6 Optional core (same PR if small)

- [ ] `[~]` `BP-PY-3` Raise Generic Exception — defer if PR size; reason: lower severity / noise; next gate Phase 5
- [ ] `[~]` `BP-PY-5` Wildcard Import — defer if PR size; allowlist `__init__.py` when implemented

### 2.7 Phase 2 validation

- [x] All Phase 2 registered IDs appear in `RuleIDs()`
- [ ] `make lint` — unchecked until green
- [ ] `make test` — unchecked until green

---

## Phase 3: Security hygiene (Batch B)

### 3.1 `BP-PY-8` subprocess With shell=True

- [x] Detector: `subprocess.(run|Popen|call|check_output|check_call)` with `shell=True`
- [x] Register + fixtures hit/miss (`shell=False` / argv list miss) (inline)
- [x] Proof: finding severity **high**

### 3.2 `BP-PY-9` os.system Or os.popen

- [x] Detector: `os.system(` / `os.popen(`
- [x] Register + fixtures hit/miss (inline)
- [x] Proof: true positive on vulnerable fixture

### 3.3 `BP-PY-10` pickle Loads Untrusted Data

- [x] Detector: `pickle.load` / `pickle.loads` / `_pickle` / `cloudpickle` calls
- [x] Register + fixtures (non-constant source preferred; literal-only may still flag at v0 with note) (inline)
- [x] Proof: hit/miss tests

### 3.4 `BP-PY-11` yaml.load Without SafeLoader

- [x] Detector: `yaml.load(` without `Loader=yaml.SafeLoader` / `CSafeLoader`; prefer flagging bare `yaml.load`
- [x] Miss: `yaml.safe_load` / explicit SafeLoader
- [x] Register + fixtures (inline)
- [x] Proof: hit/miss tests

### 3.5 `BP-PY-12` eval Or exec On Dynamic Input

- [x] Detector: `eval(` / `exec(` / `compile(..., 'exec')` with non-literal args (v0: any non-string-literal call site)
- [x] Register + fixtures (inline)
- [x] Proof: hit/miss tests

### 3.6 `BP-PY-13` Hardcoded Secret In Source

- [x] Detector: assignments to names matching `password|secret|api_key|token|private_key` (case-insensitive) with non-empty string literals
- [x] Skip obvious placeholders (`changeme`, empty, env-lookup patterns) only if documented false-positive policy
- [x] Register + fixtures (inline)
- [x] Proof: hit/miss tests

### 3.7 Phase 3 validation

- [x] Batch B all registered and fixture-backed (inline unit tests)
- [ ] `make lint` — unchecked until green
- [ ] `make test` — unchecked until green

---

## Phase 4: Framework high-signal + templates/DB (Batches C + D)

### 4.1 Flask (`BP-PY-16`, `17`, `20`)

- [x] `BP-PY-16` — `app.run(debug=True)`, `DEBUG = True` / `app.config['DEBUG'] = True` in non-test modules
- [x] `BP-PY-17` — `app.secret_key = '...'` / `SECRET_KEY = '...'` string literals in Flask-ish config
- [~] `BP-PY-20` — `send_file` / `send_from_directory` with path from `request.args` / `request.form` / view args (heuristic) — deferred
- [ ] Register each; fixtures under `tests/fixtures/python/bp/`
- [ ] `[~]` `BP-PY-18`, `19` — lower signal / harder heuristics; defer to Phase 5 unless free

### 4.2 Django (`BP-PY-21`, `22`, `24`, `25`, `26`)

- [x] `BP-PY-21` — `DEBUG = True` in `settings.py` / settings modules (skip tests / `local_settings` patterns if documented)
- [~] `BP-PY-22` — `SECRET_KEY = '...'` literals in Django settings — deferred
- [~] `BP-PY-24` — deferred
- [~] `BP-PY-25` — deferred
- [~] `BP-PY-26` — deferred
- [ ] Register each + hit/miss fixtures
- [ ] `[~]` `BP-PY-23`, `27`, `28` — ALLOWED_HOSTS / mass assignment / N+1; defer (noise or multi-line)

### 4.3 FastAPI / Starlette (`BP-PY-30`, `32`)

- [ ] `BP-PY-30` — `async def` route/dependency bodies calling `time.sleep`, `requests.`, sync SQLAlchemy session, blocking `open` (heuristic windows)
- [ ] `BP-PY-32` — `FileResponse` / path from path params / query without sanitization needle
- [ ] Register + fixtures
- [ ] `[~]` `BP-PY-29`, `31` — mutable global Depends / response_model; defer

### 4.4 Templates / DB (`BP-PY-33`, `35`)

- [ ] `BP-PY-33` — `jinja2.Environment(autoescape=False)` (and obvious autoescape off)
- [ ] `BP-PY-35` — `text(f"...")` / `text("...".format` / concatenated SQL into `text(`
- [ ] Register + fixtures
- [ ] `[~]` `BP-PY-34` Markup/`|safe` — templates may be non-`.py`; defer or limited Python-side Markup()
- [ ] `[~]` `BP-PY-36`, `37` — session close / DB-API `%` format; optional same PR if small

### 4.5 Phase 4 validation

- [ ] All Phase 4 shipped IDs fixture-backed with true positives
- [x] Severity/category match catalogue metadata for sampled rules
- [ ] `make lint` — unchecked until green
- [ ] `make test` — unchecked until green

---

## Phase 5: Remaining catalogue (Batch E) — optional follow-up PR

Ship only after A–D stable. Prefer one rules file per domain (async, testing, deps, observability, production).

### 5.1 Production hardening

- [ ] `BP-PY-14` requests without timeout
- [ ] `BP-PY-48` CORS `*` with credentials
- [ ] `BP-PY-49` TLS verify disabled (`verify=False`, etc.)
- [ ] `BP-PY-50` insecure cookie flags
- [ ] Each: register + hit/miss fixtures + proof

### 5.2 Resource / HTTP client

- [ ] `BP-PY-15` httpx AsyncClient not closed
- [ ] Fixtures + proof

### 5.3 Async

- [x] `BP-PY-38` create_task without reference
- [x] `BP-PY-39` time.sleep in async def (overlaps `30`; keep distinct rule ids)
- [x] `BP-PY-40` threading without join — review-only / low confidence OK if documented
- [x] Fixtures + proof (`rules_async.go` / `rules_async_test.go`)

### 5.4 Testing / deps / observability

- [x] `BP-PY-41`, `42` testing hygiene
- [x] `BP-PY-43` requirements without pins — path-gated on `requirements*.txt` unit path (Path A; plugin Extensions stay `.py`-only)
- [x] `BP-PY-44` deprecated stdlib imports
- [x] `BP-PY-45` `sys.path` mutation
- [x] `BP-PY-46` print debugging (info)
- [x] `BP-PY-47` logging f-string before logger
- [x] Fixtures + proof per shipped rule (`rules_testing.go`, `rules_deps.go`, `rules_observability.go` + matching `*_test.go`)

### 5.5 Deferred leftovers (explicit)

- [ ] List any unshipped `BP-PY-*` here as `[~]` with reason, owner (#53 follow-up), next gate
- [ ] Do not claim “50/50 implemented” until every id has detector + fixture or an explicit `[~]` row

### 5.6 Phase 5 validation

- [ ] `make lint` — unchecked until green
- [ ] `make test` — unchecked until green

---

## Phase 6: Integration, product surface, closure

### 6.1 Fixtures manifest

- [ ] Add `[[fixture]]` entries to `tests/fixtures/manifest.toml` for each shipped vulnerable fixture (`lang = "python"`, `required_rules = ["BP-PY-N"]`) and safe fixtures (`required_rules = []`)
- [ ] Materialize path remains `.txt` → `.py` via `internal/fixture` (already supports `LangPython`)

### 6.2 Engine / CLI smoke

- [ ] With `languages = ["python"]` (or registry including Python), scan of vulnerable fixture root surfaces `BP-PY-*` findings
- [ ] Severity/category match catalogue for at least one rule per phase batch
- [ ] Go-only default still produces **no** Python findings on pure-Go scans
- [ ] Optional: `./bin/goslop --list-rules` (or equivalent) lists registered `BP-PY-*` when Python enabled — only if product already lists by plugin

### 6.3 Docs / README honesty

- [ ] Update `ruleset/python/README.md` “What this is not” / status: detectors exist for **implemented subset**; list shipped IDs or point at this ledger
- [ ] Cross-link parent `plans/v0.0.2/python-heuristics.md` BP rollup rows when first PR lands
- [ ] Issue body success criteria in `plans/PR/v0.0.2/issue-python-bp-heuristics-body.md` remain the bar

### 6.4 PR hygiene

- [ ] Branch e.g. `feat/python-bp-heuristics` (or batch branches `…-core`, `…-security`, `…-framework`)
- [ ] Filled PR body under `plans/PR/v0.0.2/` from `plans/PR/PR_TEMPLATE.md`
- [ ] `Closes #53` only when priority success criteria met; else `Relates to #53` + partial `[~]` inventory
- [ ] `Relates to #51`

### 6.5 Final validation gates (required for non-docs code)

- [x] `gofmt -w` on all touched Go files
- [x] `make lint` — green 2026-07-31 feat/python-bp-heuristics
- [x] `make test` — green 2026-07-31 feat/python-bp-heuristics
- [x] `CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop`
- [ ] Optional smoke: `./bin/goslop --format text --no-cache <python-fixture-root>` with Python enabled

### 6.6 Success criteria (issue #53)

- [x] Priority BP heuristics (Batches **A+B** + Batch C high-signal `16`,`17`,`21`; C remainder + D/E `[~]`) implemented and registered on Python plugin
- [x] Fixtures prove true positives for each **shipped** `BP-PY-*` (inline unit-test snippets)
- [x] `languages = ["python"]` surfaces `BP-PY-*` findings (`--list-rules` + scan path)
- [ ] Severity/category match catalogue metadata
- [x] `make lint` + `make test` green (2026-07-31)

---

## Rule inventory (all 50)

Use as batch tracker. Status starts `[ ]`; move to `[x]` only with detector + fixture proof; `[~]` needs reason.

### Batch A — Core / error (priority)

| ID | Name | Sev | Category | Status |
|----|------|-----|----------|--------|
| BP-PY-1 | Bare Except Clause | medium | Error Handling | [x] |
| BP-PY-2 | Except Pass | medium | Error Handling | [x] |
| BP-PY-3 | Raise Generic Exception | low | Error Handling | [~] optional deferred |
| BP-PY-4 | Mutable Default Argument | high | Core Language | [x] |
| BP-PY-5 | Wildcard Import | low | Core Language | [~] optional deferred |
| BP-PY-6 | assert Used For Runtime Validation | high | Core Language | [x] |
| BP-PY-7 | open Without Context Manager | medium | Resource Management | [x] |

### Batch B — Security hygiene (priority)

| ID | Name | Sev | Category | Status |
|----|------|-----|----------|--------|
| BP-PY-8 | subprocess With shell=True | high | Security Hygiene | [x] |
| BP-PY-9 | os.system Or os.popen | high | Security Hygiene | [x] |
| BP-PY-10 | pickle Loads Untrusted Data | high | Security Hygiene | [x] |
| BP-PY-11 | yaml.load Without SafeLoader | high | Security Hygiene | [x] |
| BP-PY-12 | eval Or exec On Dynamic Input | high | Security Hygiene | [x] |
| BP-PY-13 | Hardcoded Secret In Source | high | Security Hygiene | [x] |

### Batch C — Framework high-signal (priority)

| ID | Name | Sev | Category | Status |
|----|------|-----|----------|--------|
| BP-PY-16 | Flask DEBUG True In Production Code | high | Flask | [x] |
| BP-PY-17 | Flask SECRET_KEY Hardcoded | high | Flask | [x] |
| BP-PY-20 | Flask send_file User Path | high | Flask | [~] deferred |
| BP-PY-21 | Django DEBUG True In Settings | high | Django | [x] |
| BP-PY-22 | Django SECRET_KEY Hardcoded | high | Django | [~] deferred |
| BP-PY-24 | Django raw SQL With Format | high | Django | [~] deferred |
| BP-PY-25 | Django mark_safe On Dynamic Data | high | Django | [~] deferred |
| BP-PY-26 | Django csrf_exempt On State-Changing View | high | Django | [~] deferred |
| BP-PY-30 | FastAPI Blocking I/O In Async Route | high | FastAPI | [~] deferred |
| BP-PY-32 | Starlette FileResponse User Path | high | FastAPI | [~] deferred |

### Batch D — Templates / DB (priority)

| ID | Name | Sev | Category | Status |
|----|------|-----|----------|--------|
| BP-PY-33 | Jinja2 autoescape Disabled | high | Templates | [~] deferred |
| BP-PY-35 | SQLAlchemy text With F-String | high | Database | [~] deferred |

### Batch E — Remaining (follow-up)

| ID | Name | Sev | Category | Status |
|----|------|-----|----------|--------|
| BP-PY-14 | requests Without Timeout | medium | Production Hardening | [~] deferred |
| BP-PY-15 | httpx Async Client Not Closed | medium | Resource Management | [~] deferred |
| BP-PY-18 | Flask Route Missing Methods Restriction | low | Flask | [~] deferred |
| BP-PY-19 | Flask jsonify Error Leaks Exception | medium | Flask | [~] deferred |
| BP-PY-23 | Django ALLOWED_HOSTS Empty Or Star | medium | Django | [~] deferred |
| BP-PY-27 | Django Mass Assignment From request.POST | medium | Django | [~] deferred |
| BP-PY-28 | Django N+1 Query In Loop | medium | Django | [~] deferred |
| BP-PY-29 | FastAPI Depends On Mutable Global | medium | FastAPI | [~] deferred |
| BP-PY-31 | FastAPI response_model Disabled Unsafely | medium | FastAPI | [~] deferred |
| BP-PY-34 | Jinja2 Markup Or safe Filter On Variables | high | Templates | [~] deferred |
| BP-PY-36 | SQLAlchemy Session Not Closed | medium | Database | [~] deferred |
| BP-PY-37 | DB-API Cursor Execute With Percent Format | high | Database | [~] deferred |
| BP-PY-38 | asyncio create_task Without Reference | medium | Async | [x] `rules_async.go` |
| BP-PY-39 | time.sleep In Async Function | high | Async | [x] `rules_async.go` |
| BP-PY-40 | threading Without Join Or Shutdown | low | Async | [x] `rules_async.go` |
| BP-PY-41 | pytest assert With Side Effects Only | info | Testing | [x] `rules_testing.go` |
| BP-PY-42 | unittest Assert Without Context On Raises | low | Testing | [x] `rules_testing.go` |
| BP-PY-43 | requirements Without Pins | low | Dependency Hygiene | [x] `rules_deps.go` (path-gated) |
| BP-PY-44 | Import Deprecated stdlib Module | low | Dependency Hygiene | [x] `rules_deps.go` |
| BP-PY-45 | sys.path Mutation At Runtime | low | Dependency Hygiene | [x] `rules_deps.go` |
| BP-PY-46 | print Debugging In Library Code | info | Observability | [x] `rules_observability.go` |
| BP-PY-47 | logging With String Format Before Logger | info | Observability | [x] `rules_observability.go` |
| BP-PY-48 | CORS Allow Origins Star With Credentials | high | Production Hardening | [~] deferred |
| BP-PY-49 | TLS Verification Disabled | high | Production Hardening | [~] deferred |
| BP-PY-50 | Django/Flask CSRF Or Session Cookie Insecure Flags | medium | Production Hardening | [~] deferred |

---

## Suggested first PR slice

Minimal shippable PR (closes or substantially advances #53):

1. Phase 1 scaffold + metadata + plugin registration  
2. Phase 2 Batch A (`1`,`2`,`4`,`6`,`7`)  
3. Phase 3 Batch B (`8`–`13`)  
4. Phase 4 subset if review budget allows (`16`,`17`,`21`,`22`,`30`,`33`,`35` high-signal)  
5. Phase 6 gates: fixtures + `make lint` + `make test`  

Second PR: remainder of C/D + Batch E.

---

## Dependencies

| Depends on | Note |
|------------|------|
| Epic #39 / PR #50 | Plugin stub, `languages`, `ruleset/python/bad-practices.json` |
| Parent epic #51 | Rollup ledger `plans/v0.0.2/python-heuristics.md` |
| `core.LanguagePlugin` / `core.Detector` | Registration surface on Python plugin |
| `internal/rules` | `RuleMetadata`, `Finding`, `PackFromRuleID` (`BP-` prefix already OK for `BP-PY-*`) |
| Go BP package | **Pattern reference only** — `internal/lang/go/detectors/bad_practices/` |
| Sibling #52 | May create `internal/lang/python/detectors/` first; coordinate package root, do not block on full CWE |
| Out of scope | #54 PERF; tree-sitter; Go `BP-*` IDs |

---

## References

- Issue body: `plans/PR/v0.0.2/issue-python-bp-heuristics-body.md`  
- Epic body: `plans/PR/v0.0.2/issue-epic-python-heuristics-body.md`  
- Catalogue: `ruleset/python/bad-practices.json`  
- README: `ruleset/python/README.md`  
- Skill: `plans/skills/phase-wise-checklist/SKILLS.md`  
- Foundation: `plans/v0.0.2/python-support.md`  
- Go detectors: `internal/lang/go/detectors/bad_practices/{register,scan,rules_core,common,facts,metadata_gen}.go`  
- Plugin stub: `internal/lang/python/plugin.go`  
