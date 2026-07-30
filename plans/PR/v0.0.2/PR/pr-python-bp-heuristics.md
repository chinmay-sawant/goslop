## Summary

Implements **Python BP-PY bad-practice heuristics** (issue #53): Phase 1 detector scaffold under `internal/lang/python/detectors/bad_practices/`, Batch A core rules, Batch B security hygiene, and Batch C high-signal Flask/Django DEBUG/SECRET_KEY rules. Wires the Python plugin `Detectors()` / `NewDetectors()` so `languages = ["python"]` surfaces `BP-PY-*` findings. Pure-Go source-pattern heuristics (no Python AST / CGO). Remaining catalogue IDs deferred as `[~]`.

---

## Motivation / context

- Plans: `plans/v0.0.2/heuristics/python-heuristics-bp.md`, parent `plans/v0.0.2/heuristics/python-heuristics.md`
- Issues: see **Related issues**
- Catalogue already shipped in epic #39 / PR #50 (`ruleset/python/bad-practices.json`); detectors were zero

---

## Changes

### Scaffold + plugin wire

- `internal/lang/python/detectors/all.go` — additive `All()` returning BP scan (CWE may append later)
- `internal/lang/python/detectors/bad_practices/` — `register`, `scan`, `common`, `facts`, `metadata` (hand-written metas for implemented IDs; validates against catalogue JSON in tests)
- `internal/lang/python/plugin.go` — `Detectors` / `NewDetectors` call `detectors.All()`; package doc no longer “zero detectors”
- Plugin / registry / `--list-rules` tests updated for non-empty `BP-PY-*`

### Batch A — Core

| ID | Heuristic |
|----|-----------|
| BP-PY-1 | bare `except:` / broad `Exception`/`BaseException` without re-raise |
| BP-PY-2 | except suite solely `pass` |
| BP-PY-4 | mutable default `[]`/`{}`/`set()`/`list()`/`dict()` |
| BP-PY-6 | `assert` in non-test modules |
| BP-PY-7 | `open(` / `.open(` without `with` |

### Batch B — Security hygiene

| ID | Heuristic |
|----|-----------|
| BP-PY-8 | `subprocess.*(… shell=True …)` |
| BP-PY-9 | `os.system` / `os.popen` |
| BP-PY-10 | `pickle`/`cloudpickle`/`_pickle` load(s) |
| BP-PY-11 | `yaml.load` without SafeLoader |
| BP-PY-12 | `eval`/`exec`/`compile(...,'exec')` on non-literal args |
| BP-PY-13 | hardcoded secret-like string assigns (placeholder-skipping) |

### Batch C — Framework high-signal

| ID | Heuristic |
|----|-----------|
| BP-PY-16 | Flask `debug=True` / DEBUG True |
| BP-PY-17 | Flask `secret_key` / SECRET_KEY string literals |
| BP-PY-21 | Django `DEBUG = True` in settings |

### Deferred (`[~]`)

- BP-PY-3, 5 (optional core)
- Rest of Batch C/D/E (20, 22, 24–26, 30, 32, 33, 35, 14–15, 18–19, … 50)

---

## Impact

| Area | Impact |
|------|--------|
| **Performance** | Light source scans over Python units only when Python plugin enabled |
| **Memory** | Negligible (line tables + needle index per file) |
| **Behavior / correctness** | DefaultRegistry remains Go-only; Python opt-in surfaces BP-PY findings |
| **API / CLI** | `--list-rules` with `languages=["python"]` lists BP-PY-* |
| **Dependencies** | None |
| **Binary size / build time** | Small pure-Go package addition |

---

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| None for Go-only default | - |
| Tests / docs that assumed Python had zero detectors | Updated (`plugin_test`, `registry_test`, `run_test`) |

---

## Test plan

- [x] `make test`
- [x] `make lint` / `go vet`
- [x] `CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop`
- [x] `go test ./internal/lang/python/...`
- [ ] `make run` wall time vs baseline (hard < 400ms; soft ±50ms of reference) — N/A for Python-only change on Go default path
- [ ] `make reference-metrics` — not required (Go detector surface unchanged)

### Commands

```sh
make lint
make test
CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop
go test ./internal/lang/python/...
```

---

## Screenshots / sample output

```
go test ./internal/lang/python/...
ok  github.com/chinmay-sawant/goslop/internal/lang/python
ok  github.com/chinmay-sawant/goslop/internal/lang/python/detectors/bad_practices

# --list-rules with languages=["python"] surfaces BP-PY-1 … BP-PY-21 (shipped subset)
```

---

## Related issues

- Closes #53
- Relates to #51

---

## PR metadata checklist (author)

- [x] Self-assigned (`--assignee @me`)
- [x] Labels applied (`enhancement`, `documentation`)
- [x] Related issues filled with real ticket IDs
- [x] Filled body committed under `plans/PR/v0.0.2/PR/pr-python-bp-heuristics.md`

---

## Follow-ups (out of scope)

- Remaining Batch C/D/E BP-PY rules
- CWE detectors (#52)
- PERF catalogue + detectors (#54)
- Optional on-disk fixtures under `tests/fixtures/python/bp/` + manifest entries
- Full Python AST path (still source-only parse)

---

## Reviewer checklist

- [ ] Behavior matches summary and test plan
- [ ] No unrelated changes in diff
- [ ] Public API / CLI changes documented
- [ ] New rules have fixture coverage when applicable (inline unit tests)
- [ ] PR has assignee and labels
- [ ] Related issues use correct Closes/Relates keywords
- [ ] No secrets or generated artifacts committed
- [ ] IDs are `BP-PY-*` only (no bare `BP-*` on Python path)
- [ ] Go BP package untouched
