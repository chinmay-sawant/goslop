# Batch 07 — Testing, dependency hygiene, observability

> **Parent:** `plans/v0.0.2/heuristics/bp-plans/README.md` — batch index  
> **Canonical #53 ledger:** `plans/v0.0.2/heuristics/python-heuristics-bp.md`  
> **Issue:** [#53](https://github.com/chinmay-sawant/goslop/issues/53) expansion under epic [#51](https://github.com/chinmay-sawant/goslop/issues/51)  
> **Status:** **complete** — shipped on `main` (PR #65 / merge `2b3e635`); 50/50 catalogue coverage
> **Estimated effort:** 1 PR (medium)  
> **PR policy:** single batchwise PR; title e.g. `python(bp): batch-07 testing/deps/obs (BP-PY-41..47)`

---

## Overview

| Group | IDs | Target Go file (new) |
|-------|-----|----------------------|
| Testing | **BP-PY-41**, **BP-PY-42** | `internal/lang/python/detectors/bad_practices/rules_testing.go` |
| Dependency hygiene | **BP-PY-43**, **BP-PY-44**, **BP-PY-45** | `…/rules_deps.go` |
| Observability | **BP-PY-46**, **BP-PY-47** | `…/rules_observability.go` |

**All seven IDs ship in this batch** — no “later later” deferrals. If a rule needs path/engine support, the support is **in-scope rows** below, not a silent skip.

| Constraint | Detail |
|------------|--------|
| Package | `internal/lang/python/detectors/bad_practices/` |
| Line budget | **1500 soft / 2000 hard** per file — three domain files keep each small |
| Do not dump into | `rules_core.go` / `rules_security.go` / `rules_framework.go` / `rules_async.go` / `rules_prod.go` |
| Registration | each file’s `init()` → `RegisterRule` |
| Tests | hit/miss per ID; split `rules_*_test.go` if `scan_test.go` grows past soft cap |

### Catalogue contract

| ID | Name | Sev | Category |
|----|------|-----|----------|
| BP-PY-41 | pytest assert With Side Effects Only | info | Testing |
| BP-PY-42 | unittest Assert Without Context On Raises | low | Testing |
| BP-PY-43 | requirements Without Pins | low | Dependency Hygiene |
| BP-PY-44 | Import Deprecated stdlib Module | low | Dependency Hygiene |
| BP-PY-45 | sys.path Mutation At Runtime | low | Dependency Hygiene |
| BP-PY-46 | print Debugging In Library Code | info | Observability |
| BP-PY-47 | logging With String Format Before Logger | info | Observability |

---

## Executive Summary

```text
rules_testing.go        → 41, 42
rules_deps.go           → 43, 44, 45
rules_observability.go  → 46, 47
facts.go                → needles for print, logging, sys.path, deprecated imports, etc.
plugin / scan gating    → BP-PY-43 may need requirements*.txt units (see Phase 4)
```

Info-severity rules (41, 46, 47) are still **required** detectors with hit/miss proofs; severity comes from catalogue metadata.

---

## Phase 1: Scaffold domain files

### 1.1 Create three rule files

- [x] Create `rules_testing.go` with `init()` registering **BP-PY-41**, **BP-PY-42**
- [x] Create `rules_deps.go` with `init()` registering **BP-PY-43**, **BP-PY-44**, **BP-PY-45**
- [x] Create `rules_observability.go` with `init()` registering **BP-PY-46**, **BP-PY-47**
- [x] Each file: `package badpractices`; use `MetadataForID` + `pushAt`
- [x] Confirm no file exceeds 1500 soft after implementation

### 1.2 Facts / needles

- [x] Extend `bpNeedles` in `facts.go` for: `assert`, `pytest`, `sys.path`, `print(`, `logging`, `logger.`, deprecated module names as needed (`imp`, `asyncore`, `asynchat`, `cgi`, `telnetlib`, …)
- [x] Prefer `facts.has` / `hasAny` early-outs in each detector

### 1.3 Metadata smoke

- [x] Metadata non-nil for 41–47; packs bad-practice; severities match catalogue (info/low as above)

---

## Phase 2: Testing — `BP-PY-41`, `BP-PY-42` (`rules_testing.go`)

### 2.1 `BP-PY-41` — test functions with side effects only (no asserts)

- [x] Implement `detectBPPY41`
- [x] Scope: functions named `test_*` **or** methods in test files (`isPythonTestFile` / path `test_*.py`, `*_test.py`, `tests/`)
- [x] Hit: test function body calls something (name call, attribute call) but contains **no** `assert `, no `pytest.raises`, no `self.assert*`, no `unittest` assert helpers
- [x] Miss: same shape with at least one assert / `pytest.raises` / `self.assertEqual`
- [x] Low confidence / style signal — message should say heuristic / info
- [x] Severity **info** from metadata

### 2.2 `BP-PY-41` proof

- [x] Hit fixture source: `def test_api():\n    client.get('/x')\n` in `test_api.py`
- [x] Miss: `def test_api():\n    r = client.get('/x')\n    assert r.status_code == 200\n`
- [x] Register + `RuleIDs` contains `BP-PY-41`
- [x] `TestBPRulesRegistered` want-list includes `BP-PY-41`

### 2.3 `BP-PY-42` — bare try/except instead of assertRaises / pytest.raises

- [x] Implement `detectBPPY42`
- [x] Scope: test functions / test files
- [x] Hit: `try:` … `except AssertionError` or broad `except:` / `except Exception` used to “expect” failure without `assertRaises` / `pytest.raises` / `raises` context manager
- [x] Miss: `with self.assertRaises(...)` / `with pytest.raises(...)`
- [x] Severity **low**

### 2.4 `BP-PY-42` proof

- [x] Hit: try/except AssertionError pattern in `test_foo.py`
- [x] Miss: `with pytest.raises(ValueError): ...`
- [x] Register + tests green for `BP-PY-42`

---

## Phase 3: Dependencies — `BP-PY-44`, `BP-PY-45` (pure `.py` first)

### 3.1 `BP-PY-44` — deprecated stdlib imports

- [x] Implement `detectBPPY44` in `rules_deps.go`
- [x] Maintain an explicit allowlist/denylist of deprecated modules (comment source: PEP / 3.11–3.13 removals), minimum set: `imp`, `asyncore`, `asynchat`, `cgi`, `telnetlib`, `uu`, `xdrlib`, `aifc`, `audioop`, `chunk`, `msilib`, `nis`, `ossaudiodev`, `pipes`, `sunau` (trim to a documented subset if noise, but **ship non-empty detector**)
- [x] Hit: `import imp`, `from asyncore import ...`, etc.
- [x] Miss: modern replacements (`importlib`, `asyncio`, …)
- [x] Severity **low**; message names preferred replacement when known

### 3.2 `BP-PY-44` proof

- [x] Hit: `import imp\n`
- [x] Miss: `import importlib\n`
- [x] Register + `RuleIDs` includes `BP-PY-44`

### 3.3 `BP-PY-45` — sys.path mutation at runtime

- [x] Implement `detectBPPY45`
- [x] Hit: `sys.path.insert(`, `sys.path.append(`, `sys.path.extend(` outside packaging bootstrap
- [x] Skip / miss: test files via `isPythonTestFile` (document); optional miss for files that are clearly `sitecustomize` / install scripts if path basename matches known set
- [x] Miss: reading `sys.path` without mutation
- [x] Severity **low**

### 3.4 `BP-PY-45` proof

- [x] Hit: `import sys\nsys.path.insert(0, './lib')\n` in `app.py`
- [x] Miss: same code under `tests/test_path.py` (if skip policy enabled)
- [x] Register + tests for `BP-PY-45`

---

## Phase 4: `BP-PY-43` — requirements without pins (**must ship**, not deferred)

Catalogue expects parsing `requirements*.txt` bare package lines without `==` / `~=` / `>=` (exclude `-e`, `-r`, blank, comments).

### 4.1 Engine / unit surface (choose and implement one path — both rows tracked)

- [x] **Path A (preferred v0):** allow detector to run when `unit.Path` / display path basename matches `requirements.txt`, `requirements-*.txt`, `requirements/*.txt` even if body is not Python syntax; gate in `detectBPPY43` on path, not on Python grammar
- [x] **Path B (if A needs product change):** extend Python plugin `Extensions()` **or** walk/scan plan so `requirements*.txt` becomes a `LanguagePython` (or dedicated) unit when `languages` includes python — document choice in PR; implement minimal change under `internal/lang/python/` and/or engine only if required for integration tests
- [x] Unit tests **must** exercise the detector with path `requirements.txt` and source lines (do not leave 43 untested because default Extensions is `py` only)

### 4.2 Detector behavior

- [x] Implement `detectBPPY43` in `rules_deps.go`
- [x] Hit lines: `requests`, `Django`, `flask` (bare name) without version operators
- [x] Miss: `requests==2.31.0`, `django>=4.2`, `foo~=1.0`, `-r other.txt`, `-e .`, `# comment`, empty lines
- [x] Optional miss: VCS lines (`git+…`) documented
- [x] Severity **low**; message prefers pins / lockfiles for apps

### 4.3 `BP-PY-43` proof

- [x] Hit: path `requirements.txt`, body `requests\nflask\n`
- [x] Miss: path `requirements.txt`, body `requests==2.31.0\n`
- [x] Miss: path `app.py` with unrelated content does not false-fire (path gate)
- [x] Register + `RuleIDs` includes `BP-PY-43`
- [x] If plugin Extensions changed: `plugin_test.go` updated for new extensions list

---

## Phase 5: Observability — `BP-PY-46`, `BP-PY-47` (`rules_observability.go`)

### 5.1 `BP-PY-46` — print debugging in library code

- [x] Implement `detectBPPY46`
- [x] Hit: `print(` in modules that are **not** under `if __name__ == "__main__":` guard region (heuristic: flag prints not indented under a detected main guard) and **not** test files
- [x] Miss: prints only under `__main__` guard; miss `logging.info`
- [x] Miss: `isPythonTestFile` true
- [x] Severity **info**

### 5.2 `BP-PY-46` proof

- [x] Hit: `def f():\n    print('debug')\n` in `lib.py`
- [x] Miss: `if __name__ == '__main__':\n    print('cli')\n`
- [x] Miss: same print in `test_lib.py`
- [x] Register + tests for `BP-PY-46`

### 5.3 `BP-PY-47` — logging f-string / format before logger

- [x] Implement `detectBPPY47`
- [x] Hit: `logger.debug/info/warning/error/critical/exception(f"...")` or `.format(` as **eager** first arg; also `logging.info(f"...")`
- [x] Miss: `logger.info("x %s", val)` / `logger.info("x {}", …)` lazy styles; miss non-f-string literals without `.format` on the call arg
- [x] Severity **info**; message prefers lazy `%s` / structured logging

### 5.4 `BP-PY-47` proof

- [x] Hit: `logger.info(f"user={user}")\n`
- [x] Miss: `logger.info("user=%s", user)\n`
- [x] Register + tests for `BP-PY-47`

---

## Phase 6: Package integration

### 6.1 Registration table

- [x] `RuleIDs()` contains **all** of: 41, 42, 43, 44, 45, 46, 47
- [x] `TestBPRulesRegistered` want-list updated for all seven
- [x] Collision guard still forbids bare `BP-<n>`
- [x] No accidental registration of batch-06/08 IDs in these files

### 6.2 Shared helpers

- [x] Reuse `isPythonTestFile`, `pushAt`, line scanners from `common.go`
- [x] Keep requirements line parser local to `rules_deps.go` (or tiny helper) — do not bloat `common.go` past soft cap

### 6.3 Line-limit gate

- [x] `wc -l` each of `rules_testing.go`, `rules_deps.go`, `rules_observability.go`, `facts.go`, `common.go`, test files — all under **2000** hard; soft **1500** preferred
- [x] Split tests into domain `*_test.go` if needed

---

## Phase 7: Validation gates (required for code)

- [x] `gofmt -w` on all touched Go files
- [x] `go test ./internal/lang/python/... -count=1` green
- [x] `make lint` — green; record: ________
- [x] `make test` — green; record: ________
- [x] Optional build: `CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop`
- [x] Parent inventory: mark 41–47 `[x]` only with evidence; update `_inventory.json`
- [x] PR: `Relates to #53`, `Relates to #51` (batch PR does not alone close #53)

---

## Dependencies

| Depends on | Note |
|------------|------|
| Batch 00 | register/scan/metadata/common shipped |
| Catalogue | all 41–47 keys present in `ruleset/python/bad-practices.json` |
| Optional engine tweak | only if BP-PY-43 Path B chosen |
| Not blocked on | batch-06 async or batch-08 prod (parallel-safe if no file conflicts) |

---

## Non-goals

- Full pytest AST semantic “was this assertion meaningful?”  
- Lockfile (`poetry.lock` / `uv.lock`) completeness analysis  
- Structured logging library-specific APIs beyond stdlib `logging` / common `logger.*`  
- BP-PY-38–40 or BP-PY-48–50 in this PR  

---

## Complete ID checklist (none deferred)

- [x] **BP-PY-41** — `rules_testing.go` + hit/miss + RuleIDs
- [x] **BP-PY-42** — `rules_testing.go` + hit/miss + RuleIDs
- [x] **BP-PY-43** — `rules_deps.go` + path/unit support + hit/miss + RuleIDs
- [x] **BP-PY-44** — `rules_deps.go` + hit/miss + RuleIDs
- [x] **BP-PY-45** — `rules_deps.go` + hit/miss + RuleIDs
- [x] **BP-PY-46** — `rules_observability.go` + hit/miss + RuleIDs
- [x] **BP-PY-47** — `rules_observability.go` + hit/miss + RuleIDs

**Batch-07 ID list:** BP-PY-41, BP-PY-42, BP-PY-43, BP-PY-44, BP-PY-45, BP-PY-46, BP-PY-47

---

## Completion stamp

- [x] Batch ledger synchronized to code on `main` (PR #65, `2b3e635`, 2026-07-31)
- [x] All catalogue IDs in this batch have `RegisterRule` + hit/miss tests (or prior shipped evidence)
- [x] File size policy observed (≤1500 soft / 2000 hard per Go domain file)
- [x] Validation: `make lint` + `make test` green on integration before merge
