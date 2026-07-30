# Go code style and design-patterns review — phase-wise checklist

**Review date:** 2026-07-30  
**Scope:** Current `fix/phasewise-review-remediation` application state, assessed with the `golang-code-style` and `golang-design-patterns` practices.  
**Baseline considered:** the earlier Ponytail and architecture checklists, including remediation commits `0a5de57` and `7f64ffa`.

## Rating

| Area | Rating | Assessment |
| --- | ---: | --- |
| Go code style | **8.4 / 10** | Clear control flow and consistently formatted code; a few export/config seams can become more explicit and cohesive. |
| Go design patterns | **8.2 / 10** | Scan-local state, bounded concurrency, and error propagation are strong; cache invalidation currently masks filesystem failures. |
| **Overall** | **8.3 / 10** | Mature, maintainable Go application with two concrete medium-priority improvements and no evidence for a broad redesign. |

## Status legend

- [x] Confirmed good practice or completed earlier remediation.
- [ ] Concrete improvement warranted now.
- [~] Conditional follow-up; do not implement without the stated trigger.

## Phase 0 — review boundary and carried-forward evidence

- [x] The earlier remediation remains present: session-local detectors, registry catalogue/factory parity checks, single-flight project facts, parser-quality diagnostics, ownership-safe export cleanup, and resolved-config tests.
- [x] `7f64ffa` records the latest lint-all cleanup; the prior remediation ledger records successful `make lint`, `make test`, and `go test -race ./...`.
- [x] This review performed a source/style pass only. It does not claim a fresh full test run.
- [x] No current high-severity regression was confirmed.

## Phase 1 — Go code-style improvements

### Export validation and rendering cohesion

- [ ] **G-STYLE-1 — make directory comparison read-only.** **Priority: Medium.** `dirsEqual` creates both directories before reporting that context and chunk output locations collide. An invalid configuration can therefore leave a directory behind. Compare cleaned/absolute paths without writing; retain `MkdirAll` in the export branch that actually writes. Evidence: `internal/export/export.go:84`, `internal/export/export.go:139`.
- [ ] **G-STYLE-2 — group export rendering dependencies.** **Priority: Medium.** `writeChunkFiles` accepts eight positional parameters and `formatFindingBlock` accepts seven, including the related file/source/AST caches and rendering policy. Introduce one unexported rendering state struct so output-mode changes cannot miswire correlated arguments. Evidence: `internal/export/export.go:130`, `internal/export/export.go:202`, `internal/export/export.go:262`.
- [ ] **G-STYLE-3 — use a byte reader for TOML input.** **Priority: Low.** `config.Parse` converts `[]byte` to `string` solely to create an `io.Reader`. Replace `strings.NewReader(string(data))` with `bytes.NewReader(data)` to state intent and avoid the conversion. Evidence: `internal/config/config.go:98`.
- [ ] **G-STYLE-4 — state the intentional partial-AST policy directly.** **Priority: Low.** `getParsed` binds `err` from `parser.ParseFile` only to discard it when no AST exists. Use `file, _ := ...` (with the existing partial-AST comment) or retain the error as diagnostic data; either makes the selected policy obvious. Evidence: `internal/export/export.go:386`.

### Confirmed style strengths

- [x] Configuration validation uses small, direct guard clauses and wrapped errors rather than panics or package initialization work: `internal/config/config.go`.
- [x] `AnalyzerBuilder` is an appropriate fluent builder: it configures genuinely optional collaborators (registry, policy, cache, baseline, walk options, root) without inventing one-off functional options: `internal/engine/analyzer.go:32`.
- [x] Error ownership is generally explicit at scan boundaries; cache `Prune`/`Flush` failures propagate through the analyzer instead of being silently ignored: `internal/engine/analyzer.go:304`.

## Phase 2 — Go design-pattern improvements

### Error flow and mutable state

- [ ] **G-DESIGN-1 — surface cache invalidation cleanup failures.** **Priority: Medium.** On tool-version mismatch, cache opening resets the manifest but ignores `clearFilesDir` errors; `clearFilesDir` then ignores each `RemoveAll` error. The cache may claim old entries were discarded while stale files remain. Return a wrapped cache error, or explicitly preserve/report an orphaned-cleanup state. Evidence: `internal/engine/cache/store.go:94`, `internal/engine/cache/store.go:473`.
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

1. [ ] **Medium:** Make output-directory collision checking side-effect free; add a regression proving rejected equal paths do not create either directory.
2. [ ] **Medium:** Propagate or explicitly report cache version-invalidation cleanup failures; add a deterministic filesystem-failure regression.
3. [ ] **Medium:** Consolidate export rendering arguments behind an unexported state object; retain context/chunk output byte-for-byte tests.
4. [ ] **Low:** Use `bytes.NewReader` and make partial-AST error handling self-documenting; run targeted config/export tests.
5. [~] Revisit conditional rows only when their triggers occur.

## Conclusion

The completed prior checklist work materially improved the Go application: scan state is localized, error behavior is more deliberate, and configuration/output contracts are tested. The current review finds improvements, but they are focused—not a reason to redesign the application. Address the two medium error/side-effect issues first; defer the conditional architecture work until the product has a real second use case.
