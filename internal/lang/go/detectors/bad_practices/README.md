# Go bad-practice (BP) detectors

Phase 8 implementation of Go bad-practice heuristics.

## Layout

| File | Role |
|------|------|
| `scan.go` | `GoBadPracticeScan` runner — per-rule enable via `ScanContext.Allows`, severity overrides |
| `register.go` | `RegisterRule` catalogue (init-time registration) |
| `facts.go` | Shared AST/source facts (tree-sitter when available) |
| `project.go` | Project snapshot + caches for server/module rules (BP-47/50/54/55, BP-57+) |
| `metadata_gen.go` | Catalogue metadata from `ruleset/golang/bad-practices.json` (135 rules) |
| `rules_*.go` | Heuristic detectors by domain |

## Ruleset

- Source of truth: `ruleset/golang/bad-practices.json` (~135 BP rules)
- Fixtures: `tests/fixtures/go/bad_practices/` and `tests/fixtures/go/bad_practices_projects/`

## Profile parity

| Profile | BP enabled? |
|---------|-------------|
| `recommended` | **off** |
| `perf` | off |
| `security` | off |
| `style` | **on** (`BP-*`) |
| `all` | **on** |

`ScanContext.BadPracticesEnabled` and pack `Only` filters enforce this. Global BP severity override: `ScanContext.BadPracticeSeverity`.

## Project-level rules

Server-policy rules emit **once** on the server entrypoint (`requiresServerAnchor`):

- **BP-47** — missing graceful shutdown
- **BP-50** — missing signal handling
- **BP-54** — public endpoint without rate limiting
- **BP-55** — missing request-id propagation

Module-hygiene rules emit on the project anchor (`requiresProjectAnchor`): BP-57…BP-65 (subset implemented with go.mod heuristics).

## Coverage notes

Heuristic ports of the Rust modules under `codehound/src/lang/go/detectors/bad_practices/`.
Priority fidelity: **BP-1**, **BP-5** (top oracle frequency), then broad catalogue registration.

Not every fixture is expected to pass; AST-precise rules from Rust may still need tightening. Gaps are listed in the Phase 8 PR body.

## Tests

```sh
go test ./internal/lang/go/detectors/bad_practices/ -count=1
```
