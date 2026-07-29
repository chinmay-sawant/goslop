# Go CWE detectors

Structural and taint-lite security heuristics for Go.

## Architecture

| Piece | Role |
|-------|------|
| `GoCweScan` | Unified detector (one `core.Detector`, many rules) |
| `RegisterRule` / `init` | Catalogue registration (PERF-style) |
| `GoCweFacts` + `cweNeedles` | Per-file `SourceIndex` fact bag |
| `rules_table.go` | SI / museum needle rules for structural domains |
| `rules_taintlite.go` | Same-file heuristics for taint-core IDs (22/78/79/89/90/91) |
| `metadata_generated.go` | Titles/descriptions from `ruleset/golang/chunks/cwe-*.json` |
| `registry/*.toml` | Domain inventory (175 rows; source of truth for IDs) |

Wired via `detectors.All()` → `cwe.NewGoCweScan()`.

## Coverage

**175 / 175** registry IDs registered.

| Registry file | Rows | Status |
|---------------|-----:|--------|
| `registry.access_control.toml` | 30 | structural SI ports |
| `registry.concurrency.toml` | 6 | structural SI ports |
| `registry.configuration.toml` | 6 | structural SI ports |
| `registry.credentials_and_secrets.toml` | 15 | structural SI ports |
| `registry.cryptography.toml` | 9 | structural SI ports |
| `registry.deserialization.toml` | 3 | structural SI ports |
| `registry.file_handling.toml` | 1 | structural SI ports |
| `registry.general_security.toml` | 75 | structural SI ports |
| `registry.information_exposure.toml` | 12 | structural SI ports |
| `registry.injection.toml` | 7 | 78/89/90/91 taint-lite; 93/619/917 structural |
| `registry.input_validation.toml` | 5 | structural SI ports |
| `registry.input_validation_redos.toml` | 2 | structural SI ports |
| `registry.network_binding.toml` | 1 | structural SI ports |
| `registry.path_traversal.toml` | 1 | CWE-22 taint-lite |
| `registry.request_handling.toml` | 2 | structural SI ports |
| **Total** | **175** | |

### Taint-core (Phase 7 lite / Phase 9 full)

| Rule | Phase 7 | Phase 9 |
|------|---------|---------|
| CWE-22 | same-file path sink + request taint; Base/HasPrefix safe | full taint graph |
| CWE-78 | seed: `exec.Command` + request/shell | full taint graph |
| CWE-79 | text/template + request; `html.EscapeString` safe | full taint graph |
| CWE-89 | seed: dynamic Query/Exec SQL | full taint graph |
| CWE-90 | LDAP filter + request, no EscapeFilter | full taint graph |
| CWE-91 | XML parse/build + request | full taint graph |

### Disposition note

Many structural rows are **fixture-shaped SourceIndex museums** (parity with Rust trust freezes). They register and fire on the stdlib/frameworks corpus; confidence is set to `0.55` for table rules. Production-shaped call-fact rewrites and full taint remain follow-ups (Phase 9 / trust hardening).

## Tests

```sh
go test ./internal/lang/go/detectors/cwe/ -count=1
```

- `registry_test.go` — 175 IDs registered, metadata present, structural sample matrix, taint-lite fixtures
- `cwe78_test.go` / `cwe89_test.go` — seed adapters

## Regeneration

Needle tables / metadata were produced from Rust `domains/` + ruleset chunks via an offline generator (see session notes). Prefer hand-edits for trust freezes; re-gen only when bulk-syncing from Rust.
