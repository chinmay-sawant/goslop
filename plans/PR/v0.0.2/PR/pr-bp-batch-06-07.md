## Summary

Implements **Python BP-PY batches 06–07** (async, testing, dependency hygiene, observability): ten catalogue rules **BP-PY-38..47** as pure-Go source heuristics under `internal/lang/python/detectors/bad_practices/`, with domain-split rule files under the soft 1500 / hard 2000 line budget.

---

## Motivation / context

- Plans: `plans/v0.0.2/heuristics/bp-plans/batch-06-async.md`, `plans/v0.0.2/heuristics/bp-plans/batch-07-testing-deps-obs.md`
- Parent ledger: `plans/v0.0.2/heuristics/python-heuristics-bp.md`
- Issues: see **Related issues**

---

## Changes

### Batch 06 — Async (`rules_async.go`)

| ID | Heuristic |
|----|-----------|
| BP-PY-38 | bare `create_task` / `ensure_future` statement (return unused) |
| BP-PY-39 | `time.sleep` inside `async def` body (prefer `await asyncio.sleep`) |
| BP-PY-40 | `threading.Thread` + `.start()` without `.join` in same unit (review-only; skip `daemon=True` start lines) |

### Batch 07 — Testing / deps / observability

| ID | File | Heuristic |
|----|------|-----------|
| BP-PY-41 | `rules_testing.go` | `test_*` body with calls but no assert / `pytest.raises` / `self.assert*` |
| BP-PY-42 | `rules_testing.go` | try/except AssertionError/Exception instead of raises context |
| BP-PY-43 | `rules_deps.go` | bare package lines in `requirements*.txt` without pins (path-gated) |
| BP-PY-44 | `rules_deps.go` | import deprecated/removed stdlib (`imp`, `asyncore`, `cgi`, …) |
| BP-PY-45 | `rules_deps.go` | `sys.path.insert/append/extend` (skip test files) |
| BP-PY-46 | `rules_observability.go` | `print(` outside `__main__` guard / tests |
| BP-PY-47 | `rules_observability.go` | eager f-string / `.format` log messages |

### Support

- `facts.go` — needles for async/testing/deps/obs fast-paths
- `metaByID` populated in each domain file `init()` before `RegisterRule`
- Domain tests: `rules_async_test.go`, `rules_testing_test.go`, `rules_deps_test.go`, `rules_observability_test.go`
- Inventory / ledger: `_inventory.json`, `python-heuristics-bp.md` marked for 38–47

### BP-PY-43 path note

Detector is **path-gated** on `requirements.txt` / `requirements-*.txt` / `requirements/*.txt`. Unit tests construct `LanguagePython` units with those paths. Plugin `Extensions()` remains `.py`-only — scanning real `requirements.txt` on disk still needs a future walk/extension change (Path B); unit coverage is complete via Path A.

---

## Impact

| Area | Impact |
|------|--------|
| **Performance** | Light source scans; needle prefilters on new IDs |
| **Memory** | Negligible |
| **Behavior / correctness** | Python opt-in (`languages = ["python"]`) surfaces BP-PY-38..47 |
| **API / CLI** | `--list-rules` lists new IDs when Python enabled |
| **Dependencies** | None |
| **Binary size / build time** | Small pure-Go addition |

---

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| None | - |

---

## Test plan

- [x] `make test`
- [x] `make lint` / `go vet`
- [x] Package tests hit/miss for BP-PY-38..47
- [x] Line counts ≤2000 hard (all new domain files ≪1500 soft)

### Commands

```sh
make lint
make test
wc -l internal/lang/python/detectors/bad_practices/rules_{async,testing,deps,observability}.go
```

---

## Related issues

- Relates to #53
- Relates to #51

---

## PR metadata checklist (author)

- [x] Self-assigned (`--assignee @me`)
- [x] Labels applied (`enhancement`, `documentation`)
- [x] Related issues filled
- [x] Filled body under `plans/PR/v0.0.2/PR/pr-bp-batch-06-07.md`

---

## Follow-ups (out of scope)

- BP-PY-30 FastAPI-scoped sleep (distinct from 39)
- Plugin extension so production walks feed `requirements*.txt` units to BP-PY-43
- Batches 01–05, 08 remaining catalogue IDs
