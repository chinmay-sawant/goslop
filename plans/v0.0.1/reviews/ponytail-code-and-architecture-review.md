# Ponytail code and architecture review — phase-wise improvement checklist

**Reviewed:** 2026-07-30
**Scope:** Whole Go application: CLI, scan engine, configuration, cache, Go detector catalogues, reporting, export, and integration harness. This is a source review, not a diff review.
**Method:** Four independent local review passes: two Ponytail code-quality passes and two `improve-codebase-architecture` exploration passes. No MCP tools were used.
**Repository context:** No `CONTEXT.md` or `docs/adr/` directory exists, so there is no project glossary or recorded architecture decision to constrain the recommendations.

## Ratings

| Report | Rating | Basis |
|---|---:|---|
| Ponytail code-quality review | **6.5 / 10** | Two high-risk scan-trust defects, destructive export semantics, incomplete configuration contracts, and scan-lifetime global state outweigh an otherwise well-covered, idiomatic core. |
| Architecture-deepening review | **6.8 / 10** | The language-plugin, analyzer, cache, and fixture modules have useful depth, but scan state, configuration resolution, project facts, and output ownership leak across seams. |
| Combined current rating | **6.6 / 10** | The product has a strong foundation and a passing suite, but a static analyzer cannot score higher while cache state and partial-scan failures can silently change its security result. |

## Evidence and validation boundary

- [x] Four source passes completed locally; no MCP used and no production code changed.
- [x] `make test` passed: `go test ./...`, including the normal integration package.
- [x] The detector/reporting pass also ran targeted `go test ./internal/rules ./internal/reporting ./internal/export ./internal/lang/go/...` plus seed/dogfood integration coverage.
- [~] Full fixture matrices, concurrency/race checks, and the proposed cold-to-warm cache regressions were not run; every unchecked item below needs its listed acceptance test before closure.

## Report 1 — Ponytail code-quality findings

### Confirmed defects

1. **P0 — warm-cache taint scans can lose inter-procedural state and miss findings.** `Analyzer.scanOne` returns a cache hit before parsing or calling `AccumulateState` ([`internal/engine/analyzer.go`](../../internal/engine/analyzer.go):350), while the taint detector requires that state before finalization. A cold scan and an otherwise identical warm scan can therefore disagree on CWE-22/78/79/89 findings.
2. **P0 — partial scans can be reported as clean and exit successfully.** File failures are collected into `AnalysisResult.Errors` ([`internal/engine/analyzer.go`](../../internal/engine/analyzer.go):235) but `app.run` prints the normal summary without reporting them or choosing a failing exit path ([`internal/app/run.go`](../../internal/app/run.go):209). CI can accept a scan which skipped unreadable, invalid-UTF-8, or invalid-Go files.
3. **P1 — caller-selected context output can delete unrelated `*.txt` files.** `ExportFindings` unconditionally cleans every text file in `--context-dir` ([`internal/export/export.go`](../../internal/export/export.go):107); the preceding numbered-file cleanup is immediately superseded. Deletion failures are ignored ([`internal/export/export.go`](../../internal/export/export.go):171).
4. **P1 — several accepted configuration controls are inert.** `evict_target_ratio` is parsed but omitted from `Merged` and cache options ([`internal/config/config.go`](../../internal/config/config.go):49, 183, 267); `languages` and `typed.enabled` likewise never affect the scan. A valid configuration can silently fail to honor the user’s requested behavior.
5. **P1 — BP scan facts are process-global, not scan-scoped.** The BP module installs mutable `activeCaches`; overlapping analyzers can replace or clear each other’s project state ([`internal/lang/go/detectors/bad_practices/project.go`](../../internal/lang/go/detectors/bad_practices/project.go):181; [`scan.go`](../../internal/lang/go/detectors/bad_practices/scan.go):89).
6. **P1 — reporter behavior is inconsistent.** JSON copies findings before fingerprints; SARIF writes fingerprints into the caller-owned slice ([`internal/reporting/json.go`](../../internal/reporting/json.go):28; [`sarif.go`](../../internal/reporting/sarif.go):95). Reusing one result across formats has an undocumented side effect and can race.
7. **P1 — the integration harness leaves temporary fixture trees behind.** Its temporary materialization directory has no `RemoveAll`, multiplying stale `goslop-fixture-*` directories across fixture runs.
8. **P1 — copied documentation/template controls are unsupported.** The init template/docs advertise a `warnings` fail policy and flags such as `--warnings-as-errors` / `--strict`, which the parser/fail-policy implementation does not accept.
9. **P2 — BP-30/BP-31 repeatedly read and parse whole package directories.** Per-source execution redoes package facts, creating avoidable quadratic I/O on large packages.
10. **P2 — normal cache flush/prune failures are discarded.** The analyzer ignores their errors ([`internal/engine/analyzer.go`](../../internal/engine/analyzer.go):299), unlike the explicit prune command.
11. **P2 — CLI output injection is incomplete.** `run` accepts writers, but `init` writes to process stdout directly, defeating capture at the seam.

### Review judgement calls, not hard defects

- `internal/app/run.go` is an orchestration hotspot: config, cache, baseline, scan, output, export, and exit policy live in one long function. This is possible **Divergent Change** / **Shotgun Surgery**, not a violation by itself.
- Multi-language and typed-analysis scaffolding is possible **Speculative Generality** until an active language implementation needs it.
- `sourceutil` has no direct tests despite parser-like helpers used by CWE-78 and taint-lite; malformed-source, comments, raw strings, nested arguments, and selectors need table/fuzz coverage.

## Report 2 — architecture-deepening findings

### Material strengths

- [x] `Analyzer.AnalyzePaths` is a deep module: one interface hides discovery, bounded concurrency, filtering, ordering, cache use, baseline filtering, and result aggregation.
- [x] The language-plugin seam has leverage: the engine does not depend on Go AST types, while one parsed unit is shared across the detector catalogue.
- [x] Cache persistence/invalidation and fixture parsing/materialization keep substantial implementation behind small interfaces.

### Candidates

1. **Strong — make detector state a per-scan `ScanSession` module.** Registry-held detector instances are reused across scans while `AnalyzePaths` starts lifecycle hooks and runs file work concurrently. Taint state and BP global caches show that the current interface leaks scan lifetime. Keep an immutable catalogue in `Registry`; let each scan session own detector instances, project facts, cache-hit state accumulation, and finalization.
2. **Strong — deepen configuration into a resolved `ScanPlan` module.** TOML, merged config, app field mutations, and cache/walk structures collectively define behavior. The inert fields prove the interface is distributed. Resolve CLI + config once into a plan that owns policy, walk options, cache options, baseline choice, and roots.
3. **Strong — make BP project/package facts a scan-scoped `ProjectFacts` module.** It should own root discovery, package snapshots, package docs, synchronization, and reuse. Its small interface can expose `ForProject`, `ForPackage`, and `ForPackageDocs`; no global adapter is justified.
4. **Strong — deepen export through ownership-aware output targets.** A target should synchronize only files it owns, return cleanup errors, and use a manifest or staging-and-rename result. Context and chunk modes are two real strategies with different owned-file predicates.
5. **Worth exploring — compile each language rule catalogue into immutable entries.** An entry should own ID, metadata, gate, and detector function; use a concrete catalogue with `Entries` and `ByID`, then verify ruleset IDs, compiled descriptors, and execution stay in parity.
6. **Worth exploring — normalize findings once at the reporting seam.** Return an owned normalized slice before text, JSON, or SARIF rendering so the `Reporter` interface no longer leaks format-dependent mutation.
7. **Worth exploring — introduce a `ScanScope` module.** Resolve requested paths, display paths, and cache-relative keys once. Keep BP Go-module discovery separate because it has a different language-specific interface and semantics.

## Phase 0 — restore scan-result trust

### P0.1 Cache-state parity

- [x] Create per-scan detector state; cache hits populate every detector state required for finalization.
- [x] Preserve cache performance: state-only parsing/accumulation does not rerun ordinary detector rules on a hit.
- [x] Add a multi-file taint fixture where an inter-file source/sink is found on a cold run and the exact same result is found after a warm run.
- [x] Add a regression for mixed cache hit/miss scans and concurrent analyzer isolation.
- [x] Acceptance: `make test` proves cold, warm, and mixed results have identical taint finding fingerprints.

### P0.2 Explicit partial-scan outcome

- [x] Define the product policy: eligible per-file analysis failures fail by default with `ExitInternal`.
- [x] Emit every `AnalysisResult.Errors` entry to stderr without contaminating JSON/SARIF stdout.
- [~] Test unreadable files, invalid UTF-8, and parse failures at the CLI boundary. Invalid UTF-8 is covered; add the remaining two variants.
- [x] Acceptance: no failed input file can produce an unqualified clean scan or exit code 0.

## Phase 1 — make destructive and configuration contracts honest

### P1.1 Ownership-safe export

- [x] Replace broad context-directory cleanup with explicit owned-file tracking: only numeric context files and `Chunk_*.txt` chunk files are owned.
- [x] Preserve unrelated files in a caller-selected output directory; return cleanup and write errors.
- [x] Remove the redundant cleanup pass.
- [~] Test unrelated text preservation, empty-result reconciliation, partial-write behavior, cleanup failure, and context/chunk collision. Unrelated-file and empty-result reconciliation are covered; add the failure/collision cases.

### P1.2 Resolved `ScanPlan`

- [x] Add a plan-resolution module that resolves CLI and TOML precedence into effective scan, cache, baseline, and walk settings.
- [x] Wire and validate `evict_target_ratio`; add configuration coverage for the value.
- [x] Reject/remove inert `languages` and `typed.enabled` keys and every template/schema/document promise that they work.
- [x] Remove stale `warnings`, `--warnings-as-errors`, and `--strict` promises rather than advertise unsupported controls.
- [~] Acceptance: precedence coverage exists for the changed configuration controls; broaden observable-effect coverage across every supported field.

### P1.3 Reporting contract

- [x] Normalize/copy findings before any format-specific work so machine reporters are read-only from the caller’s perspective.
- [~] Test fingerprint/location parity across text, JSON, and SARIF, input immutability, and concurrent reporter use. Input immutability is covered; add parity/concurrency cases.

## Phase 2 — concentrate scan lifetime and project facts

### P2.1 `ScanSession`

- [x] Keep the registry catalogue immutable and construct detector/session state per `AnalyzePaths` invocation.
- [x] Move detector lifecycle, run/finalization, and required cache-state accumulation behind the session interface.
- [x] Remove mutable process-global detector/project state from the scan path.
- [~] Acceptance: the concurrent-scanner and cache hit/miss regressions pass under `make test`; run the same concurrency coverage with `-race`.

### P2.2 BP `ProjectFacts`

- [x] Make root/project/package/package-doc facts owned by the BP scan, not a package-global pointer.
- [x] Cache package type facts once per directory with scan-local synchronization.
- [x] Update BP-30/BP-31 to consume the shared package snapshot.
- [~] Add a benchmark over a multi-file package and a concurrent distinct-root regression.
- [~] Acceptance: scan-local cache identity is tested; add the benchmark and race-enabled distinct-root proof.

### P2.3 Temporary fixture lifecycle

- [x] Ensure every temporary materialization directory is removed on success and all failure paths.
- [~] Add a test seam that proves cleanup after materialize and analyze errors.

## Phase 3 — deepen catalogue, scope, and CLI orchestration

- [ ] Introduce an immutable compiled language catalogue with descriptors for ID, metadata, gate, and execution; retain one concrete implementation until a second adapter is real.
- [ ] Add ruleset-to-catalogue-to-execution parity tests and make late registration behavior explicit.
- [~] Introduce `ScanScope` to own root interpretation, display paths, and cache-relative keys. Directory/file cache roots are covered; display paths plus relative/multi-root matrices remain in the engine.
- [~] Split `app.run` around resolved-plan creation, scan execution, output, and exit policy while preserving one CLI interface. Resolution and scope are extracted; scan/output staging remains a follow-up.
- [x] Route `init` through the supplied writer and make cache flush/prune failure visible according to an explicit policy.

## Phase 4 — harden maintainability and test surfaces

- [x] Add direct table and fuzz tests for `sourceutil` token/argument helpers.
- [x] Keep one interface as the test surface for each new deep module; no hypothetical adapter was introduced.
- [~] Run `go test -race ./...`, the complete fixture matrices, formatting/vet/lint gates used by CI, and a representative cold/warm cache performance comparison. `gofmt` and `make test` pass; the remaining gates are outstanding.
- [~] Re-rate only after Phases 0–2 are complete and the final diff passes the expanded validation set.

## Recommended delivery order

1. P0.2 partial-scan failure policy — shortest route to trustworthy CI results.
2. P0.1 cache-state parity — closes cache-dependent security false negatives.
3. P1.1 export ownership — stops user-data loss risk.
4. P2.1/P2.2 scan-session and BP project facts — removes the root lifecycle cause and enables safe cache reuse.
5. P1.2, then Phases 3–4 — configuration honesty and structural leverage after product correctness is protected.
