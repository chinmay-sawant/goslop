# Python CWE detectors

Pure-Go **source-pattern** security heuristics for Python (issue #52).

## Scope

| Piece | Role |
|-------|------|
| `PyCweScan` | Unified detector (`LanguagePython`, many rules) |
| `RegisterRule` / `init` | Catalogue registration with optional needle gates |
| `PyCweFacts` + `pyCweNeedles` | Per-file `ast.SourceIndex` prefilter |
| `metadata.go` | Hand-authored meta for priority batch |
| `rules.go` | CWE-22 / 78 / 79 / 89 / 502 heuristics |

**No CGO / tree-sitter.** Detectors operate on `unit.Source` only (source-only `ParsedUnit`).

## Priority batch (v0.0.2 / #52)

| Rule | Sinks (high signal) | Safe suppressions |
|------|---------------------|-------------------|
| CWE-502 | `pickle.loads/load/Unpickler`, `yaml.load` without SafeLoader, `yaml.unsafe_load` | `yaml.safe_load`, `json.loads` |
| CWE-78 | `os.system` / `os.popen` dynamic; `subprocess.*(shell=True)` dynamic | list argv + `shell=False` / default, pure string literals |
| CWE-89 | `execute` / `executemany` with f-string / `%` / `.format` SQL | parameterized `?` / `%s` + bound args tuple |
| CWE-22 | `open(os.path.join(root, user))`, `Path(root) / user` without confinement | `os.path.basename` only; `resolve` + `startswith` root |
| CWE-79 | `mark_safe`, `Markup(`, `render_template_string` dynamic | plain `render_template("…html", …)` |

## Non-goals

- Full 344-rule catalogue registration
- BP-PY / PERF detectors (siblings #53 / #54)
- Inter-procedural taint graph
- Modifying Go CWE detectors

## Adding a rule later

1. Ensure the CWE id exists in `ruleset/python/chunks/cwe-*.json`.
2. Add `MetaCWEN` in `metadata.go` (titles aligned with catalogue).
3. Add needles to `pyCweNeedles` for prefilter gates.
4. `RegisterRule("CWE-N", detect, &MetaCWEN, gates...)` in `init`.
5. Unit hit/miss tests in `rules_test.go`.

## Tests

```sh
go test ./internal/lang/python/detectors/cwe/ -count=1
```

Wired via `internal/lang/python/detectors.All()` → Python plugin `Detectors` / `NewDetectors`.
