## Summary

Implements **priority Python CWE heuristics** (CWE-22, 78, 79, 89, 502) as pure-Go source-pattern detectors under `internal/lang/python/detectors/cwe`, wired into the Python language plugin. No CGO/tree-sitter. DefaultRegistry remains Go-only; Python stays opt-in via `languages = ["python"]`.

---

## Motivation / context

- Plans: `plans/v0.0.2/heuristics/python-heuristics-cwe.md`, parent `plans/v0.0.2/heuristics/python-heuristics.md`
- Issues: see **Related issues**
- Foundation: epic #39 / PR #50 shipped Python plugin stub + catalogues; this PR turns the priority CWE batch into runnable detectors

---

## Changes

### Framework (`internal/lang/python/detectors/cwe/`)

- `PyCweScan` unified detector (`LanguagePython`, `RuleIDs`, `MetadataFor`, `Run`)
- `RegisterRule(id, fn, meta, gates...)` with needle prefilters via pack-local `pyCweNeedles` + `ast.SourceIndex`
- Hand-authored `MetaCWE22` / `78` / `79` / `89` / `502` (`rules.Meta`, Severity High, PackSecurity)
- README documents scope, safe suppressions, and how to add rules later

### Priority rules

| Rule | Hits | Safe suppressions |
|------|------|-------------------|
| **CWE-502** | `pickle.loads/load/Unpickler`, `yaml.load` w/o SafeLoader, `yaml.unsafe_load` | `yaml.safe_load`, `Loader=yaml.SafeLoader`, `json.loads` |
| **CWE-78** | `os.system`/`os.popen` dynamic; `subprocess.*(shell=True)` dynamic | list argv, pure string literals, `shell=False` |
| **CWE-89** | `execute`/`executemany` with f-string / `%` / `.format` SQL | parameterized `?`/`%s` + bound args |
| **CWE-22** | `open(os.path.join(root, user))`, `Path(root) / user` without confinement | `os.path.basename`, `resolve`+`startswith` |
| **CWE-79** | `mark_safe`/`Markup`/`render_template_string` dynamic | plain `render_template`, pure literal mark_safe |

### Wiring

- `internal/lang/python/detectors/all.go` → `[]core.Detector{cwe.NewPyCweScan()}` (additive surface for BP later)
- `plugin.go` Detectors/NewDetectors → `detectors.All()`; package comment updated
- Tests updated: plugin, engine registry, app `list-rules` for `languages=["python"]`

### Fixtures

- Hit/miss corpus under `tests/fixtures/python/cwe/` for all five IDs

### Docs / plans

- `ruleset/python/README.md` notes runnable priority CWE detectors
- Plan ledgers updated with evidence checkboxes

---

## Impact

| Area | Impact |
|------|--------|
| **Performance** | Negligible on default path (Python opt-in). SourceIndex + string scan over `.py` when enabled. |
| **Memory** | Small per-file fact bag when Python is scanned |
| **Behavior / correctness** | New findings for Python sources when `languages` includes `python`. Go catalogue unchanged. |
| **API / CLI** | No new flags. `languages=["python"]` + `--list-rules` now lists CWE-22/78/79/89/502. |
| **Dependencies** | None |
| **Binary size / build time** | Small pure-Go package; `CGO_ENABLED=0` build OK |

---

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| None for Go default path | - |
| Tests that assumed zero Python detectors | Updated in this PR |

---

## Test plan

- [x] `make test`
- [x] `make lint` / `go vet`
- [x] `CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop`
- [x] `go test ./internal/lang/python/...`
- [ ] `make run` wall time vs baseline — N/A (Python opt-in; Go detectors unchanged)
- [ ] `make reference-metrics` / gopdfsuit — N/A (Go CWE surface unchanged)

### Commands

```sh
gofmt -w internal/lang/python/**/*.go
make lint
make test
CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop
go test ./internal/lang/python/... -count=1
```

### Recorded outcome (2026-07-31, `feat/python-cwe-heuristics`)

- `make lint` green
- `make test` green (all packages including integration)
- pure-Go build OK
- Unit table tests: 5 rules × hit/miss green under `./internal/lang/python/detectors/cwe/`
- Smoke: `languages=["python"]` + `os.system(cmd)` → CWE-78; `pickle.loads` with `--profile all` / `--only CWE-502` → CWE-502

---

## Screenshots / sample output

```
$ languages=["python"] scan of dynamic os.system
1 findings
  severity: 1 high
CWE-78 ... dynamic input reaches an OS command sink (os.system/os.popen)

$ go test ./internal/lang/python/detectors/cwe/ -count=1
ok  github.com/chinmay-sawant/goslop/internal/lang/python/detectors/cwe
```

---

## Related issues

- Closes #52
- Relates to #51

---

## PR metadata checklist (author)

- [x] Self-assigned (`--assignee @me`)
- [x] Labels applied (`enhancement`, `documentation`)
- [x] Related issues filled with real ticket IDs
- [x] Filled body committed under `plans/PR/v0.0.2/PR/pr-python-cwe-heuristics.md`

---

## Follow-ups (out of scope)

- Full 344 CWE catalogue registration
- BP-PY heuristics (#53)
- PERF catalogue + heuristics (#54 deferred)
- Inter-procedural Python taint
- Do not modify Go detectors

---

## Reviewer checklist

- [ ] Behavior matches summary and test plan
- [ ] No unrelated changes in diff
- [ ] Public API / CLI changes documented
- [ ] New rules have fixture coverage when applicable
- [ ] PR has assignee and labels
- [ ] Related issues use correct Closes/Relates keywords
- [ ] No secrets or generated artifacts committed
