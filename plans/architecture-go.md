# Architecture — CodeHound Go Port

> **Parent:** `plans/port-phasewise-checklist.md`
> **Status:** scaffold target (Phase 0–1)

---

## Goals (parity with Rust product)

1. Single static binary CLI: scan Go trees for **PERF / CWE / BP** findings.
2. Same **fixture-backed** rule oracle: `tests/fixtures/**/*.txt` materialize → scan → expected rule IDs.
3. Offline, no network; optional typed path via `go list` later (`--typed`).
4. Profiles: `recommended` / `security` / `all` (+ only/skip, baseline, cache, SARIF/JSON).
5. Complements golangci-lint — does not replace it.

## Non-goals (v0 of port)

- Python plugin / `SLOP101` (Rust opt-in only) — defer until Go core is green.
- Perfect binary-identical JSON field ordering vs Rust (semantic wire parity is enough).
- Criterion-class microbenches day one — wall-time smoke budgets first.

## Package layout

```
codehound-go/
  cmd/codehound/          # main
  internal/
    app/                  # run orchestration (mirrors src/app)
    cli/                  # flags / enums
    core/                 # Detector, ScanContext, LanguagePlugin, ParsedUnit, profiles
    rules/                # Finding, Severity, emit, maturity, packs
    engine/               # Analyzer, walk, parse pool, registry, cache, baseline, ignore
    ast/                  # tree-sitter walk helpers, function spans, line/col
    lang/go/              # Go plugin
      detectors/
        perf/             # facts + domains (logic in Go)
        cwe/              # structural domains + taint
        bad_practices/    # BP rules
      sinks/
      typed/              # optional --typed
    reporting/            # text, json, sarif
    export/               # context/chunks (opt-in)
    fixture/              # .txt materializer (parity with Rust fixture format)
    cwe/                  # catalog consts / refs
  ruleset/                # copied as-is from Rust
  tests/
    fixtures/             # copied as-is from Rust
    integration/          # Go tests driving materializer + analyzer
  scripts/                # findings corpus, chunks (copied)
  plans/                  # this folder
  documents/              # selected product docs for agents
```

## Parse strategy

| Choice | Decision |
|--------|----------|
| Parser | **tree-sitter Go** via `github.com/tree-sitter/go-tree-sitter` + `tree-sitter-go` grammar (behavioral parity with Rust CST heuristics) |
| Alternative (later) | `go/parser` + `go/types` for typed accuracy on a subset of rules |
| Parallelism | worker pool over files (`errgroup` / channel workers), one parser per worker |
| Facts | One fused AST walk per file → `GoPerfFacts` / `GoCweFacts` + `SourceIndex` (Aho-Corasick or multi-string index) |

## Core contracts (mirrors Rust traits)

```go
type Detector interface {
    Language() LanguageID
    RuleIDs() []string
    BeginScan(ctx *ScanContext)
    EndScan()
    Run(ctx *ScanContext, unit *ParsedUnit, out *[]Finding)
    AccumulateState(ctx *ScanContext, unit *ParsedUnit)
    RequiresCacheState(ctx *ScanContext) bool
    Finalize(ctx *ScanContext, out *[]Finding)
    ResetState()
}
```

## Pipeline

```
CLI/config/profile
  → Analyzer
  → collect entries (walk + ignore)
  → begin_scan
  → parallel: read → cache? → parse → detect
  → finalize / end_scan
  → report (text|json|sarif)
```

## Testing contract

1. **Do not rewrite fixture bodies** under `tests/fixtures/**` unless fixing a documented Rust bug.
2. Materializer must honor Rust header format:
   - `# comment`
   - `lang: go`
   - `file: name.go`
   - optional `variant:`, `expect:`, etc.
   - `---` then source
3. Integration tests assert **rule ID presence/absence** per fixture (same oracle as Rust).
4. Prefer table-driven Go tests under `tests/integration/` rather than porting every Rust `#[test]` file 1:1; keep **oracle fidelity**.

## Dependency policy (Go)

Keep the module lean:

| Need | Package |
|------|---------|
| CLI | `github.com/spf13/cobra` or std `flag` (prefer cobra if subcommands grow) |
| Config TOML | `github.com/pelletier/go-toml/v2` |
| JSON/SARIF | `encoding/json` |
| Tree-sitter | go-tree-sitter + grammar |
| Parallel walk | `golang.org/x/sync/errgroup` |
| Ignore globs | `github.com/sabhiram/go-gitignore` or reimplementation of gitignore rules |
| Hash | `crypto/sha256` |
| Color | optional `github.com/fatih/color` behind `NO_COLOR` |

Avoid large frameworks. Prefer stdlib.

## Proof of phase completion

For every non-doc phase: `go test ./...` (and later `make test` once Makefile exists) must pass for that phase’s scoped packages. Record command + outcome in the checklist.
