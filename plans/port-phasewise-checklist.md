# CodeHound Go Port — Phase-Wise Checklist

> **Parent:** Rust product at `/home/chinmay/ChinmayPersonalProjects/codehound`
> **Status:** Phase 6 PERF **239/239**; Phase 7 CWE structural **175/175 registered** (SI museums + taint-lite); Phase 9 full taint still open
> **Estimated effort:** multi-session full port (400+ Rust modules → Go packages)
> **Canonical ledger:** this file only

---

## Overview

Port **CodeHound** from Rust to Go while keeping:

1. **Heuristic oracles as-is** — `tests/fixtures/**/*.txt`, `ruleset/`, PERF/CWE registry TOMLs, findings/chunks corpora.
2. **Heuristic logic rewritten in Go** — no blind transpile; preserve rule IDs, messages, and fixture expectations.
3. **Product surface** — CLI, profiles, cache, baseline, reporters, incremental analysis.

Supporting docs: [`architecture-go.md`](./architecture-go.md), [`parity-matrix.md`](./parity-matrix.md).

---

## Executive Summary

| Phase | Name | Goal | Depends |
|------:|------|------|---------|
| 0 | Bootstrap | Module, layout, copied assets, README | — |
| 1 | Core contracts | Finding, Severity, Detector, ScanContext, profiles | 0 |
| 2 | Fixture + parse | Materializer, tree-sitter Go, ParsedUnit, SourceIndex | 1 |
| 3 | Engine MVP | Walk, registry, parallel scan, one end-to-end rule | 1–2 |
| 4 | CLI / app | Flags, profiles, exit codes, `init`, list/explain | 3 |
| 5 | Reporting | text / JSON / SARIF (+ export stubs) | 3–4 |
| 6 | PERF detectors | All PERF domains from registries | 3 |
| 7 | CWE structural | All CWE domains | 3 |
| 8 | Bad practices | BP rule modules | 3 |
| 9 | Taint | Intra/inter-procedural graph + CWE-22/78/79/89 | 7 |
| 10 | Cache / baseline / ignore | Parity with Rust incremental + ignore model | 3–5 |
| 11 | Packs / maturity / quarantine | Recommended pack, metadata, --list-rules fidelity | 6–8 |
| 12 | Parity gates | Fixture oracle suites, smoke budgets, CI, **§12.4 final export/scan oracle** | 6–11 |

**Proof commands (repo root):**

```bash
go test ./...
go build -o bin/codehound ./cmd/codehound
./bin/codehound --list-rules | head
./bin/codehound --profile recommended .
```

---

## Phase 0: Bootstrap

### 0.1 Repository skeleton
- [x] Create `codehound-go/` empty project root
- [x] Create `plans/`, `cmd/codehound/`, `internal/`, `tests/`, `ruleset/`, `scripts/`, `documents/`, `templates/`
- [x] Copy `tests/fixtures/` as-is from Rust (1746 files)
- [x] Copy `ruleset/` as-is
- [x] Copy PERF (+ CWE) detector registries into `internal/lang/go/detectors/*/registry/`
- [x] Copy `scripts/findings` and `scripts/chunks` corpora
- [x] Copy selected product docs + `codehound.schema.json` + `templates/codehound.toml`
- [x] `go mod init github.com/chinmay/codehound` (module path: `github.com/chinmay/codehound`)
- [x] Root `README.md` describing Go port status and how to run
- [x] Root `Makefile` with `build`, `test`, `lint` targets
- [x] `.gitignore` for `bin/`, materialize dirs, cache

### 0.2 Tooling baseline
- [x] Go 1.22+ toolchain documented (module uses Go 1.23; README says 1.22+)
- [x] `gofmt` / `go vet` clean on scaffold — `go test ./...` PASS (2026-07-29)
- [ ] Optional: `golangci-lint` config (minimal)

---

## Phase 1: Core contracts

### 1.1 `internal/rules`
- [x] `Severity` enum + parse/string (Info/Low/Medium/High/Critical) — `internal/rules/severity.go`; JSON lowercase; tests in `severity_test.go`
- [x] `Finding` struct with wire fields: `rule_id`, `rule_title`, `file`, `line`, `column`, `message`, `severity`, `cwe`, `fingerprint`, optional evidence/snippet — `finding.go` (`NewFinding` / `NewFindingFromMeta`)
- [x] Fingerprint v2: `codehound:2:{rule}:{file}:{msg_hash16}` (sha256 lower hex) — match Rust — `Fingerprint` / `FingerprintV2` in `fingerprint.go`
- [x] `emit.PushFinding` / metadata helper — `emit.go` + `Meta`/`NewRuleMetadata`
- [x] `RulePack` (Security / Performance / BadPractice / General) — `pack.go` + tier allow-lists
- [x] `RuleMetadata` stub (title, severity, pack, maturity) — maturity deferred; pack/severity/cwe/fix present
- [x] Unit tests for fingerprint stability — `go test ./internal/rules/...` PASS

### 1.2 `internal/core`
- [x] `LanguageID` (Go; Python reserved/deferred) — `language.go` (`LanguageGo` / `LanguagePython`)
- [x] `ScanContext` (only/skip, fail policy, taint/typed flags, BP flags, retain_sources, …) — `context.go`
- [x] `ScanContext.Allows(ruleID string) bool` — only/skip/BP semantics (exact + `*` prefix) — Rust parity
- [x] `FailPolicy` / fail-on severity — `FailNone`/`FailHigh`/`FailMedium` (`FailNever` alias)
- [x] `ScanProfile` (Recommended / Security / All) + apply to context — also Perf/Style; `NewScanContext` / `ApplyBase`
- [x] `ParsedUnit` (path, display path, source, tree handle, line_col helper) — `unit.go`
- [x] `Detector` interface (BeginScan/Run/AccumulateState/Finalize/EndScan/ResetState) — `detector.go` + `BaseDetector`
- [x] `LanguagePlugin` interface (extensions, detectors, configure parser, extract deps, prepare project) — `plugin.go` + `BasePlugin`
- [x] Unit tests for `Allows` matrix — `go test ./internal/core/...` PASS

### 1.3 `internal/cwe`
- [x] `CweRef` type + catalog stubs used by findings — `ref.go` (`CweRef`/`Ref`, `New`, `NewFromID`, `FormatList`)
- [x] Load or embed minimal catalog from ruleset/Rust consts as needed — minimal helpers only (full catalog later)

---

## Phase 2: Fixture materializer + Go parse path

### 2.1 Fixture format (parity)
- [x] Parser for `#` headers + `key: value` + `---` body (Rust `fixture` module behavior) — `internal/fixture/format.go`
- [x] Materialize to temp dir preserving `file:` names — `MaterializeFixture` / `MaterializeTree`
- [x] Respect `lang: go` (python still materialized under `python/` tag; detectors may skip)
- [ ] Manifest support (`tests/fixtures/manifest.toml`) if required by integration tests
- [x] Unit tests: round-trip sample fixtures (`CWE-22-vulnerable.txt`, PERF samples)

### 2.2 Tree-sitter Go integration
- [x] Add tree-sitter Go bindings + grammar dependency (`go-tree-sitter` + `tree-sitter-go`; **CGO required**)
- [x] `ParsePool`: one parser per worker — `tsparse.NewParser` / `Parser.Parse` (pool wiring deferred to engine)
- [x] Build `ParsedUnit` from path + source bytes — `core.NewParsedUnit` / `NewParsedUnitWithTree` ready; tsparse wiring in engine Phase 3
- [x] Line/col from byte offset (UTF-8 safe byte offsets) — `internal/ast` + `tsparse.Tree.LineCol`
- [x] Iterative tree walk helpers in `internal/ast`
- [x] Smoke: parse a materialized fixture without panic

### 2.3 SourceIndex
- [x] Multi-needle index (strings.Contains build; Aho-Corasick upgrade path documented in `source_index.go`)
- [x] `Has(needle string) bool` O(1) after build
- [ ] Shared needle table for PERF/CWE fast paths (needs detector needle tables)

---

## Phase 3: Engine MVP (end-to-end one rule)

### 3.1 Walk + collect
- [x] Walk project roots; include `*.go` — `engine.CollectFiles` / `CollectGoFiles`
- [ ] Honor `.gitignore` / `.codehoundignore` / `.ignore` (minimal) — vendor/VCS skipped; gitignore later
- [x] Exclude `*_test.go` by default; `--include-tests` override — `WalkOptions.IncludeTests` + CLI
- [ ] Include/exclude globs from config (stub OK if wired)

### 3.2 Registry + analyzer
- [x] `Registry` holds language plugins + detectors — `engine.NewRegistry` / `DefaultRegistry`
- [x] Register Go plugin with detector list — `internal/lang/go` + `detectors.All()`
- [x] `Analyzer.AnalyzePaths(paths, ctx)` → `AnalysisResult{Findings, Errors, Stats}`
- [x] Detector lifecycle: begin → run per file → finalize → end
- [x] Chunked parallel file processing with bounded workers — `errgroup` + worker limit

### 3.3 First detector proof
- [x] Port high-signal structural rules end-to-end: **CWE-78**, **CWE-89**, **PERF-116** (seed heuristics)
- [x] Fixture-shaped unit tests: vulnerable fires; safe silent (`detectors/cwe/*_test.go`, `detectors/perf/*_test.go`, plugin e2e)
- [x] `go test ./internal/...` green (also `go test ./...`)

---

## Phase 4: CLI / app

### 4.1 CLI surface
- [x] `cmd/codehound` → `app.Run`
- [x] Path args, default `.`
- [x] `--profile`, `--only`, `--skip`, `--format`
- [x] `--list-rules` (seed rules + titles); `--explain` / `rules` subcommand deferred
- [x] `--include-tests`; `--verbose` / `--debug-timing` deferred (stub OK later)
- [x] Exit codes: 0 clean, 1 findings-at-fail-policy (`ExitCodeError`), 2 usage/error, 3 internal

### 4.2 Config
- [ ] Load `codehound.toml` / nested tables (schema-compatible subset)
- [ ] Merge CLI over config (only/skip additive; fail_on rules documented)
- [x] `codehound init` writes `templates/codehound.toml` (embedded)

### 4.3 Rule listing
- [x] Enumerate registered rule IDs + metadata (from engine registry / detector MetadataFor)
- [ ] Filter by category when flag present

---

## Phase 5: Reporting + export

### 5.1 Reporters
- [x] Text reporter (one line: `rule file:line:col message`; summary polish later)
- [x] JSON reporter (`{"findings":[...],"version":"..."}`)
- [x] SARIF 2.1.0 reporter (minimal valid log + results)
- [x] `NO_COLOR` honored for text (env checked; text currently uncolored)

### 5.2 Export (optional path)
- [ ] `--export-context` / `--export-chunks` full port (not stub-only)
- [ ] Default dirs: `scripts/findings/functions` (context), `scripts/chunks` (chunks)
- [ ] `retain_sources` only when export enabled
- [ ] Final gate: export counts in **§12.4** (915 context + 37 chunks on reference corpus)

---

## Phase 6: PERF detectors (logic in Go)

> Source registries: `internal/lang/go/detectors/perf/registry/*.toml`  
> Rust logic: `codehound/src/lang/go/detectors/perf/`  
> Keep rule IDs `PERF-N` stable.

### 6.0 PERF infrastructure
- [x] `GoPerfFacts` fused walk (calls, loops, defers, go stmts, conversions, …) — `facts.go` / `BuildFacts`
- [x] Shared `perf/common` helpers (hot path, handler shaped, in loop, …) — `common.go`
- [x] Metadata for implemented rules — `metadata.go` (12 rules; full catalogue later)
- [x] Domain dispatcher: `GoPerfScan` runs enabled PERF rules only — `scan.go`
- [x] Skeleton + registry inventory README — `detectors/perf/README.md` (239 rows; 12 implemented)
- [x] Go plugin `ParseSource` wires tree-sitter; engine closes trees after scan

### 6.1 Domains (check off as ported + fixture-backed)
- [x] `general_perf` (~164 registry rows) — batches 1/3/4/5 + seed (heuristic ports)
- [x] `data_access` (20) — batch 2/4 GORM/sqlx/etc.
- [x] `gin_framework` (20) — batch 2
- [x] `request_path` (9) — batch ports
- [x] `protocols` (10) — gRPC/redis/prom/Fiber (batch 2)
- [x] `parsing_in_loops` (8) — batch ports
- [x] `loop_allocations` (8) — **PERF-1…PERF-8**
- [~] Integration: seed + batch fixture samples pass; full 490-fixture matrix not yet a CI gate

### 6.2 PERF closure gate
- [x] All registry `[[detector]]` entries have a Go function (`RegisterRule` → 239/239)
- [~] `go test` PERF package green with batch samples; full materialized fixture suite still optional
- [ ] Record residual FN/FP vs Rust oracle in notes (tighten heuristics)

---

## Phase 7: CWE structural detectors

> Registries: `internal/lang/go/detectors/cwe/registry/*.toml`  
> Rust: `codehound/src/lang/go/detectors/cwe/domains/`

### 7.0 CWE infrastructure
- [x] `GoCweFacts` + `cweNeedles` SourceIndex (`detectors/cwe/facts.go`, `needles.go`)
- [x] Unified `GoCweScan` + `RegisterRule` (`scan.go`) — PERF-style catalogue
- [x] Full metadata table from ruleset chunks (`metadata_generated.go`) — all 175 IDs
- [x] Skeleton + registry inventory README — `detectors/cwe/README.md` (175 rows)
- [x] Shared source heuristics — `detectors/sourceutil` (tainted idents, call scan)
- [~] Domain packages as separate Go pkgs — **deferred**; single `cwe` package with table + taint-lite (same API surface)

### 7.1 Domains
- [x] `access_control` — SI needle ports (30)
- [x] `credentials_and_secrets` — SI needle ports (15)
- [x] `cryptography` — SI needle ports (9)
- [x] `injection` — **CWE-78/89/90/91** taint-lite; **93/619/917** structural SI
- [x] `information_exposure` — SI needle ports (12)
- [x] `input_validation` (+ redos) — SI needle ports (5+2)
- [x] `configuration` — SI needle ports (6)
- [x] `concurrency` — SI needle ports (6)
- [x] `deserialization` — SI needle ports (3)
- [x] `general_security` clusters — SI needle ports (75)
- [x] file/path/network/request registry rows — 434, 22 (taint-lite), 1327, 601/918

### 7.2 CWE closure gate
- [x] Registry complete: **175/175** IDs registered via `RegisterRule`
- [x] Unit + sample fixture suite green (`go test ./internal/lang/go/detectors/cwe/`)
- [~] Full 350-fixture matrix / call-fact structural parity — optional CI; many rules remain fixture-shaped museums (Rust trust freezes)
- [ ] Phase 9 full taint graph upgrades 22/78/79/89/90/91 beyond taint-lite

---

## Phase 8: Bad practices (BP)

> Rust: `codehound/src/lang/go/detectors/bad_practices/`  
> Ruleset: `ruleset/golang/bad-practices.json`

- [ ] BP detector runner with per-rule enable + severity override
- [ ] Port rule modules under `rules/` (error handling, concurrency, testing, API, prod hardening, …)
- [ ] Project-level rules + prewarm cache (BP-47/50/54/55 class)
- [ ] Default **recommended** pack keeps BP **off** (parity)
- [ ] Fixture / integration coverage for BP positives/negatives

---

## Phase 9: Taint tracking

> Docs: `documents/taint.md`  
> Rust: `codehound/src/lang/go/detectors/cwe/taint/`

- [ ] Taint graph extract (assignments, calls, returns)
- [ ] Sources / sinks / sanitizers tables
- [ ] Intra-procedural BFS path finding
- [ ] Bounded inter-procedural hops (`taint_max_depth`)
- [ ] Rules: CWE-22, CWE-78, CWE-79, CWE-89 (name-string model honesty preserved)
- [ ] `requires_cache_state` + accumulate/finalize lifecycle
- [ ] Fixture suite under `tests/fixtures/go/taint/` (including quarantined honest FNs)
- [ ] CLI `--taint` / config `[codehound.taint]`

---

## Phase 10: Cache, baseline, ignore

### 10.1 Incremental cache
- [ ] `.codehound-cache/` layout: `manifest.json`, `files/<sha>.json`
- [ ] Content hash + tool version + rule-config fingerprint invalidation
- [ ] Warm hit skips parse+detect (except taint accumulate)
- [ ] `--no-cache`, `--cache-dir`, `--rebuild-cache`, `--prune-cache`
- [ ] Eviction by `max_size_mb`

### 10.2 Baseline + ignore
- [ ] Inline `// codehound-ignore: RULE` (comment-only)
- [ ] File / config ignores
- [ ] Baseline store + filter (new vs baselined)
- [ ] Fixture: `tests/fixtures/go/baseline/suppressed_inline.txt`

---

## Phase 11: Packs, maturity, quarantine

- [ ] Recommended pack = S-tier PERF + taint-core CWEs; BP off; fail high (match product docs)
- [ ] Security pack enables taint
- [ ] `all` pack full catalog
- [ ] Maturity tags + quarantine reasons in `--list-rules` / `--explain`
- [ ] Wire metadata from ruleset JSON / generated tables

---

## Phase 12: Parity gates & delivery

### 12.1 Test suites
- [ ] Materialized fixture integration suite (Go)
- [ ] Profile / CLI contract tests
- [ ] JSON / SARIF snapshot or schema tests
- [ ] Perf smoke budget (generous initially; tighten later)
- [ ] Optional: compare finding sets Rust vs Go on `gopdfsuit` or fixed corpus

### 12.2 Delivery
- [ ] Multi-arch build notes (or GoReleaser config)
- [ ] CI workflow: `go test ./...`, `go vet`, build
- [ ] Version / `--version` matches release process
- [ ] Update root README: install, usage, status of port

### 12.3 Final closure
- [ ] All Phase 6–9 registry rows implemented or explicitly `[~]` with issue link
- [ ] No silent empty detector stubs in default packs
- [ ] Checklist statuses reconciled with `go test` evidence

### 12.4 Final product validation (Rust oracle → Go parity)

> **Ship gate.** Do not mark the Go port complete until this command (or its Go equivalent) matches the oracle below on the same corpus.
>
> **Corpus:** the default `make run` target tree used by the Rust product (gopdfsuit-scale / project configured in Rust `makefile` — typically the real-repo scan that yields **78** analyzed Go files).  
> **Rust reference command** (recorded 2026-07-29 on `codehound` @ `docs/expand-documents-and-frontend`):

```bash
# From the Rust codehound repo:
make run RUN_ARGS="--export-context --export-chunks --no-cache"
```

**Reference output (oracle — hold finding counts & top rules):**

```text
scanned 78 files (28120 lines) in 479.5ms
  cache: 0 hits, 78 misses (full re-analysis)
  skipped 383 files
915 findings
  severity: 10 high, 197 info, 312 low, 396 medium
  top rules: BP-1 ×181, PERF-6 ×94, PERF-32 ×59, BP-5 ×50, PERF-230 ×44
  example findings: 63 (of 915 total)
exported 915 context file(s) to scripts/findings/functions; exported 37 chunk file(s) to scripts/chunks
```

| Check | Expected (oracle) | Go port proof |
|------:|-------------------|---------------|
| Files scanned | **78** | [ ] |
| Lines scanned | **28120** | [ ] |
| Cache | **0 hits, 78 misses** (`--no-cache` / full re-analysis) | [ ] |
| Files skipped | **383** | [ ] |
| **Total findings** | **915** | [ ] |
| Severity breakdown | **10 high, 197 info, 312 low, 396 medium** | [ ] |
| Top rules (order + counts) | **BP-1×181, PERF-6×94, PERF-32×59, BP-5×50, PERF-230×44** | [ ] |
| Export context | **915** files under `scripts/findings/functions` | [ ] |
| Export chunks | **37** files under `scripts/chunks` | [ ] |
| Wall time | ~**480 ms** class on reference host (budget: document host; do not fail solely on ms until Phase 12 perf gate is set) | [ ] |

**Go equivalent (once CLI/export/cache land):**

```bash
# From codehound-go (adjust binary + default path to match makefile corpus):
./bin/codehound --profile all --export-context --export-chunks --no-cache .
# or, when Makefile exists:
# make run RUN_ARGS="--export-context --export-chunks --no-cache"
```

**Pass criteria**

- [ ] Same command surface: `--export-context`, `--export-chunks`, `--no-cache` (profile **all** implied by Rust `make run` product defaults if that is what `makefile` uses — confirm in Rust `makefile` when wiring Go Makefile)
- [ ] **Finding oracle locked:** total **915**; severity **10 / 197 / 312 / 396**; top-five rules and counts exact
- [ ] **Export oracle locked:** **915** context files + **37** chunk files to the same default dirs (`scripts/findings/functions`, `scripts/chunks`)
- [ ] **Scan shape:** **78** files scanned, **383** skipped, full re-analysis under `--no-cache`
- [ ] Record actual Go wall time + host next to this section when first green (timing is soft; counts are hard)
- [ ] Diff any residual FN/FP vs Rust only with an explicit `[~]` row and reason — do not silently change the oracle

**Hard vs soft**

| Kind | Fields |
|------|--------|
| **Hard (must match)** | findings total, severity histogram, top-rule multiset, export file counts, scanned/skipped file counts |
| **Soft (record, budget later)** | wall-clock ms (479.5ms is a reference, not a CI fail until budgets are agreed) |

---

## Dependencies (phase graph)

```
0 → 1 → 2 → 3 → 4 → 5
              ↘ 6
              ↘ 7 → 9
              ↘ 8
              ↘ 10
         6+7+8 → 11 → 12
```

---

## Working agreements

1. **Fixtures and ruleset are sacred** — copy/update only with cross-repo intent.
2. **One checklist** — update this file when closing work; do not fork status docs.
3. **Atomic rows** — check a box only with path + command evidence.
4. **Quarantine honestly** — if a rule cannot meet oracle yet, mark `[~]` with reason (like Rust taint FNs).
5. **Prefer small PRs/phases** — land engine + one rule before mass detector ports.

---

## Session log

| Date | Session | Outcome |
|------|---------|---------|
| 2026-07-29 | Init | Plans + asset copy; scaffold; 5 parallel subagents |
| 2026-07-29 | Phase 1 | `internal/{cwe,rules,core}` — findings, fingerprint v2, Detector, profiles |
| 2026-07-29 | Phase 2 | `fixture` materializer + `ast` + `tsparse` (CGO tree-sitter) |
| 2026-07-29 | Phase 3 | Engine walk/registry/analyzer; seed **CWE-78 / CWE-89 / PERF-116** |
| 2026-07-29 | Phase 4–5 | CLI + text/JSON/SARIF reporters + `init` |
| 2026-07-29 | Integrate | `go test ./...` PASS; smoke scan fires CWE-78 (exit 1) |
| 2026-07-29 | Phase 6a | PERF infra (`GoPerfFacts`, `GoPerfScan`); loop_allocations 1–8; PERF-32/50/116/230; plugin tsparse |
| 2026-07-29 | Phase 6b | 5 parallel batches → **239/239** PERF rules registered; `go test ./...` PASS |
| 2026-07-29 | Phase 7 | CWE infra + **175/175** structural/taint-lite registered; `go test ./...` PASS |

---

## Next actions (immediate)

1. **Phase 6 polish** — residual FN/FP vs Rust; expand `perfNeedles`; full fixture matrix optional CI.
2. **Phase 7 polish** — full stdlib CWE fixture matrix in CI; promote high-signal SI museums toward call-facts.
3. **Phase 8** — bad-practice modules.
4. **Phase 9** — full taint graph (upgrade CWE-22/78/79/89/90/91 toward Rust oracle).
5. **Phase 10–12** — cache/baseline/ignore, pack fidelity, fixture suite, then **§12.4** final validation (`--export-context --export-chunks --no-cache` → 915 findings / 915+37 exports).
