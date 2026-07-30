# Batch 07 — Testing, dependency hygiene, observability

> **Parent:** `plans/v0.0.2/heuristics/bp-plans/README.md` — batch index  
> **Canonical #53 ledger:** `plans/v0.0.2/heuristics/python-heuristics-bp.md`  
> **Issue:** [#53](https://github.com/chinmay-sawant/goslop/issues/53) expansion under epic [#51](https://github.com/chinmay-sawant/goslop/issues/51)  
> **Status:** implemented — `rules_testing.go` / `rules_deps.go` / `rules_observability.go` (BP-PY-41..47)  
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

- [ ] Create `rules_testing.go` with `init()` registering **BP-PY-41**, **BP-PY-42**
- [ ] Create `rules_deps.go` with `init()` registering **BP-PY-43**, **BP-PY-44**, **BP-PY-45**
- [ ] Create `rules_observability.go` with `init()` registering **BP-PY-46**, **BP-PY-47**
- [ ] Each file: `package badpractices`; use `MetadataForID` + `pushAt`
- [ ] Confirm no file exceeds 1500 soft after implementation

### 1.2 Facts / needles

- [ ] Extend `bpNeedles` in `facts.go` for: `assert`, `pytest`, `sys.path`, `print(`, `logging`, `logger.`, deprecated module names as needed (`imp`, `asyncore`, `asynchat`, `cgi`, `telnetlib`, …)
- [ ] Prefer `facts.has` / `hasAny` early-outs in each detector

### 1.3 Metadata smoke

- [ ] Metadata non-nil for 41–47; packs bad-practice; severities match catalogue (info/low as above)

---

## Phase 2: Testing — `BP-PY-41`, `BP-PY-42` (`rules_testing.go`)

### 2.1 `BP-PY-41` — test functions with side effects only (no asserts)

- [ ] Implement `detectBPPY41`
- [ ] Scope: functions named `test_*` **or** methods in test files (`isPythonTestFile` / path `test_*.py`, `*_test.py`, `tests/`)
- [ ] Hit: test function body calls something (name call, attribute call) but contains **no** `assert `, no `pytest.raises`, no `self.assert*`, no `unittest` assert helpers
- [ ] Miss: same shape with at least one assert / `pytest.raises` / `self.assertEqual`
- [ ] Low confidence / style signal — message should say heuristic / info
- [ ] Severity **info** from metadata

### 2.2 `BP-PY-41` proof

- [ ] Hit fixture source: `def test_api():\n    client.get('/x')\n` in `test_api.py`
- [ ] Miss: `def test_api():\n    r = client.get('/x')\n    assert r.status_code == 200\n`
- [ ] Register + `RuleIDs` contains `BP-PY-41`
- [ ] `TestBPRulesRegistered` want-list includes `BP-PY-41`

### 2.3 `BP-PY-42` — bare try/except instead of assertRaises / pytest.raises

- [ ] Implement `detectBPPY42`
- [ ] Scope: test functions / test files
- [ ] Hit: `try:` … `except AssertionError` or broad `except:` / `except Exception` used to “expect” failure without `assertRaises` / `pytest.raises` / `raises` context manager
- [ ] Miss: `with self.assertRaises(...)` / `with pytest.raises(...)`
- [ ] Severity **low**

### 2.4 `BP-PY-42` proof

- [ ] Hit: try/except AssertionError pattern in `test_foo.py`
- [ ] Miss: `with pytest.raises(ValueError): ...`
- [ ] Register + tests green for `BP-PY-42`

---

## Phase 3: Dependencies — `BP-PY-44`, `BP-PY-45` (pure `.py` first)

### 3.1 `BP-PY-44` — deprecated stdlib imports

- [ ] Implement `detectBPPY44` in `rules_deps.go`
- [ ] Maintain an explicit allowlist/denylist of deprecated modules (comment source: PEP / 3.11–3.13 removals), minimum set: `imp`, `asyncore`, `asynchat`, `cgi`, `telnetlib`, `uu`, `xdrlib`, `aifc`, `audioop`, `chunk`, `msilib`, `nis`, `ossaudiodev`, `pipes`, `sunau` (trim to a documented subset if noise, but **ship non-empty detector**)
- [ ] Hit: `import imp`, `from asyncore import ...`, etc.
- [ ] Miss: modern replacements (`importlib`, `asyncio`, …)
- [ ] Severity **low**; message names preferred replacement when known

### 3.2 `BP-PY-44` proof

- [ ] Hit: `import imp\n`
- [ ] Miss: `import importlib\n`
- [ ] Register + `RuleIDs` includes `BP-PY-44`

### 3.3 `BP-PY-45` — sys.path mutation at runtime

- [ ] Implement `detectBPPY45`
- [ ] Hit: `sys.path.insert(`, `sys.path.append(`, `sys.path.extend(` outside packaging bootstrap
- [ ] Skip / miss: test files via `isPythonTestFile` (document); optional miss for files that are clearly `sitecustomize` / install scripts if path basename matches known set
- [ ] Miss: reading `sys.path` without mutation
- [ ] Severity **low**

### 3.4 `BP-PY-45` proof

- [ ] Hit: `import sys\nsys.path.insert(0, './lib')\n` in `app.py`
- [ ] Miss: same code under `tests/test_path.py` (if skip policy enabled)
- [ ] Register + tests for `BP-PY-45`

---

## Phase 4: `BP-PY-43` — requirements without pins (**must ship**, not deferred)

Catalogue expects parsing `requirements*.txt` bare package lines without `==` / `~=` / `>=` (exclude `-e`, `-r`, blank, comments).

### 4.1 Engine / unit surface (choose and implement one path — both rows tracked)

- [ ] **Path A (preferred v0):** allow detector to run when `unit.Path` / display path basename matches `requirements.txt`, `requirements-*.txt`, `requirements/*.txt` even if body is not Python syntax; gate in `detectBPPY43` on path, not on Python grammar
- [ ] **Path B (if A needs product change):** extend Python plugin `Extensions()` **or** walk/scan plan so `requirements*.txt` becomes a `LanguagePython` (or dedicated) unit when `languages` includes python — document choice in PR; implement minimal change under `internal/lang/python/` and/or engine only if required for integration tests
- [ ] Unit tests **must** exercise the detector with path `requirements.txt` and source lines (do not leave 43 untested because default Extensions is `py` only)

### 4.2 Detector behavior

- [ ] Implement `detectBPPY43` in `rules_deps.go`
- [ ] Hit lines: `requests`, `Django`, `flask` (bare name) without version operators
- [ ] Miss: `requests==2.31.0`, `django>=4.2`, `foo~=1.0`, `-r other.txt`, `-e .`, `# comment`, empty lines
- [ ] Optional miss: VCS lines (`git+…`) documented
- [ ] Severity **low**; message prefers pins / lockfiles for apps

### 4.3 `BP-PY-43` proof

- [ ] Hit: path `requirements.txt`, body `requests\nflask\n`
- [ ] Miss: path `requirements.txt`, body `requests==2.31.0\n`
- [ ] Miss: path `app.py` with unrelated content does not false-fire (path gate)
- [ ] Register + `RuleIDs` includes `BP-PY-43`
- [ ] If plugin Extensions changed: `plugin_test.go` updated for new extensions list

---

## Phase 5: Observability — `BP-PY-46`, `BP-PY-47` (`rules_observability.go`)

### 5.1 `BP-PY-46` — print debugging in library code

- [ ] Implement `detectBPPY46`
- [ ] Hit: `print(` in modules that are **not** under `if __name__ == "__main__":` guard region (heuristic: flag prints not indented under a detected main guard) and **not** test files
- [ ] Miss: prints only under `__main__` guard; miss `logging.info`
- [ ] Miss: `isPythonTestFile` true
- [ ] Severity **info**

### 5.2 `BP-PY-46` proof

- [ ] Hit: `def f():\n    print('debug')\n` in `lib.py`
- [ ] Miss: `if __name__ == '__main__':\n    print('cli')\n`
- [ ] Miss: same print in `test_lib.py`
- [ ] Register + tests for `BP-PY-46`

### 5.3 `BP-PY-47` — logging f-string / format before logger

- [ ] Implement `detectBPPY47`
- [ ] Hit: `logger.debug/info/warning/error/critical/exception(f"...")` or `.format(` as **eager** first arg; also `logging.info(f"...")`
- [ ] Miss: `logger.info("x %s", val)` / `logger.info("x {}", …)` lazy styles; miss non-f-string literals without `.format` on the call arg
- [ ] Severity **info**; message prefers lazy `%s` / structured logging

### 5.4 `BP-PY-47` proof

- [ ] Hit: `logger.info(f"user={user}")\n`
- [ ] Miss: `logger.info("user=%s", user)\n`
- [ ] Register + tests for `BP-PY-47`

---

## Phase 6: Package integration

### 6.1 Registration table

- [ ] `RuleIDs()` contains **all** of: 41, 42, 43, 44, 45, 46, 47
- [ ] `TestBPRulesRegistered` want-list updated for all seven
- [ ] Collision guard still forbids bare `BP-<n>`
- [ ] No accidental registration of batch-06/08 IDs in these files

### 6.2 Shared helpers

- [ ] Reuse `isPythonTestFile`, `pushAt`, line scanners from `common.go`
- [ ] Keep requirements line parser local to `rules_deps.go` (or tiny helper) — do not bloat `common.go` past soft cap

### 6.3 Line-limit gate

- [ ] `wc -l` each of `rules_testing.go`, `rules_deps.go`, `rules_observability.go`, `facts.go`, `common.go`, test files — all under **2000** hard; soft **1500** preferred
- [ ] Split tests into domain `*_test.go` if needed

---

## Phase 7: Validation gates (required for code)

- [ ] `gofmt -w` on all touched Go files
- [ ] `go test ./internal/lang/python/... -count=1` green
- [ ] `make lint` — green; record: ________
- [ ] `make test` — green; record: ________
- [ ] Optional build: `CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop`
- [ ] Parent inventory: mark 41–47 `[x]` only with evidence; update `_inventory.json`
- [ ] PR: `Relates to #53`, `Relates to #51` (batch PR does not alone close #53)

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

- [ ] **BP-PY-41** — `rules_testing.go` + hit/miss + RuleIDs
- [ ] **BP-PY-42** — `rules_testing.go` + hit/miss + RuleIDs
- [ ] **BP-PY-43** — `rules_deps.go` + path/unit support + hit/miss + RuleIDs
- [ ] **BP-PY-44** — `rules_deps.go` + hit/miss + RuleIDs
- [ ] **BP-PY-45** — `rules_deps.go` + hit/miss + RuleIDs
- [ ] **BP-PY-46** — `rules_observability.go` + hit/miss + RuleIDs
- [ ] **BP-PY-47** — `rules_observability.go` + hit/miss + RuleIDs

**Batch-07 ID list:** BP-PY-41, BP-PY-42, BP-PY-43, BP-PY-44, BP-PY-45, BP-PY-46, BP-PY-47
