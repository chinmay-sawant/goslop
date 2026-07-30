# Go code style and design-patterns review — phase-wise checklist

**Review date:** 2026-07-30  
**Scope:** Current `fix/phasewise-review-remediation` application state, assessed with the `golang-code-style` and `golang-design-patterns` practices.  
**Baseline considered:** the earlier Ponytail and architecture checklists, including remediation commits `0a5de57` and `7f64ffa`.

## Rating

| Area | Rating | Assessment |
| --- | ---: | --- |
| Go code style | **9.1 / 10** | The export and config seams are now explicit, cohesive, and behavior-preserving. |
| Go design patterns | **9.0 / 10** | Scan-local state, bounded concurrency, and cache error propagation are now consistently deliberate. |
| **Overall** | **9.0 / 10** | All concrete review work is complete and fully validated; the remaining rows are conditional rather than current defects. |

## Status legend

- [x] Confirmed good practice or completed earlier remediation.
- [ ] Concrete improvement warranted now.
- [~] Conditional follow-up; do not implement without the stated trigger.

## Phase 0 — review boundary and carried-forward evidence

- [x] The earlier remediation remains present: session-local detectors, registry catalogue/factory parity checks, single-flight project facts, parser-quality diagnostics, ownership-safe export cleanup, and resolved-config tests.
- [x] `7f64ffa` records the earlier lint-all cleanup; the prior remediation ledger records earlier `make lint`, `make test`, and `go test -race ./...` successes.
- [x] The implementation batch closed below passed fresh `make lint-all`, `make test`, `go test -race ./...`, and `git diff --check`.
- [x] No current high-severity regression was confirmed.

## Phase 1 — Go code-style improvements

### Export validation and rendering cohesion

- [x] **G-STYLE-1 — make directory comparison read-only.** `dirsEqual` now compares absolute paths without creating either directory; the collision regression proves no output directory is left behind. Evidence: `internal/export/export.go`, `internal/export/export_test.go`.
- [x] **G-STYLE-2 — group export rendering dependencies.** An unexported per-export `renderState` owns the caches, rendering policy, and reusable blocks. The dual-surface regression proves its context and chunk output exactly matches independent single-surface output. Evidence: `internal/export/export.go`, `internal/export/export_test.go`.
- [x] **G-STYLE-3 — use a byte reader for TOML input.** `config.Parse` now passes TOML bytes directly through `bytes.NewReader`. Evidence: `internal/config/config.go`.
- [x] **G-STYLE-4 — state the intentional partial-AST policy directly.** `getParsed` discards the parser error at the call that deliberately accepts a partial AST; the existing policy comment remains adjacent. Evidence: `internal/export/export.go`.

### Confirmed style strengths

- [x] Configuration validation uses small, direct guard clauses and wrapped errors rather than panics or package initialization work: `internal/config/config.go`.
- [x] `AnalyzerBuilder` is an appropriate fluent builder: it configures genuinely optional collaborators (registry, policy, cache, baseline, walk options, root) without inventing one-off functional options: `internal/engine/analyzer.go:32`.
- [x] Error ownership is generally explicit at scan boundaries; cache `Prune`/`Flush` failures propagate through the analyzer instead of being silently ignored: `internal/engine/analyzer.go:304`.

## Phase 2 — Go design-pattern improvements

### Error flow and mutable state

- [x] **G-DESIGN-1 — surface cache invalidation cleanup failures.** Tool-version invalidation now stops with a wrapped `cache.Error` when old entry cleanup fails, and `clearFilesDir` returns the first removal error. A deterministic injected-removal regression covers that path. Evidence: `internal/engine/cache/store.go`, `internal/engine/cache/store_internal_test.go`.
- [~] **G-DESIGN-2 — freeze or fully synchronize late default-plugin registration only if it is supported at runtime.** **Priority: Low, conditional.** The default-registry pointer is protected, but `RegisterDefaultPlugin` mutates the selected registry after that lock is released. Startup-only registration is currently safe; add registry-level synchronization or a post-startup freeze only when late registration is a supported use case. Evidence: `internal/engine/registry.go:324`, `internal/engine/registry.go:358`.
- [~] **G-DESIGN-3 — document `ScanContext` as immutable during a scan only if external callers can reconfigure it concurrently.** **Priority: Low, conditional.** Workers share one context pointer, whose slices/maps are public. The current application constructs the plan before scanning and has no observed race. Snapshotting every field now would add unnecessary copying. Evidence: `internal/core/context.go:14`, `internal/engine/analyzer.go:149`.
- [~] **G-DESIGN-4 — replace the narrow export test seam only when export tests need parallel filesystem fault injection.** **Priority: Low, conditional.** `removeOwnedFile` is mutable package state, but it is private, documented, restored with `t.Cleanup`, and covers a real regression. Do not introduce a public filesystem abstraction prematurely. Evidence: `internal/export/export.go:182`, `internal/export/export_cleanup_test.go`.

### Confirmed design strengths

- [x] Detector state is created per scan; the catalogue/session parity check prevents a registry from listing a different rule multiset than it can execute: `internal/engine/registry.go`, `internal/engine/registry_session_test.go`.
- [x] BP project/package facts are scan-owned and single-flight, avoiding global cache contamination and duplicate parallel cold work: `internal/lang/go/detectors/bad_practices/project.go`, `internal/lang/go/detectors/bad_practices/project_cache_test.go`.
- [x] The analyzer has bounded parallelism and per-scan lifecycle boundaries rather than unbounded goroutines or package-global work state: `internal/engine/analyzer.go`.
- [x] Parser fallback is represented as data (`ParseQuality` and a diagnostic), allowing source-only analysis without hiding degraded parsing from callers: `internal/core/plugin.go`, `internal/lang/go/plugin.go`.

## Phase 3 — deliberately deferred non-changes

- [x] Do **not** split the broad `core.Detector` interface yet. Its no-op lifecycle methods are a manageable trade-off while there is only one concrete detector shape; revisit on the second materially different detector type.
- [x] Do **not** replace the analyzer builder with functional options. It already expresses optional dependencies clearly and has no required cross-field validation contract.
- [x] Do **not** split `app.run` merely for line count. `scanPlan` and `ScanScope` already own meaningful seams; stage further only when a second output or exit-policy flow appears.
- [~] Canonicalize multi-root cache identity only with observed cache-reuse/order friction and an acceptance test; the current first-root policy is intentional and has no reported failure.

## Phase 4 — implementation order and acceptance criteria

1. [x] **Medium:** Output-directory collision checking is side-effect free, with a no-directory regression.
2. [x] **Medium:** Cache version-invalidation cleanup failures propagate with deterministic failure coverage.
3. [x] **Medium:** Export rendering arguments are consolidated behind unexported state, with dual/single-surface byte-for-byte output coverage.
4. [x] **Low:** TOML uses `bytes.NewReader`, partial-AST error handling is self-documenting, and focused plus full validation passed.
5. [~] Revisit conditional rows only when their triggers occur.

## Conclusion

The completed prior and current checklist work materially improved the Go application: scan state is localized, error behavior is deliberate, and configuration/output contracts are tested. All concrete rows from this review are closed. The remaining architecture rows are conditional and should stay deferred until their documented product triggers occur.
