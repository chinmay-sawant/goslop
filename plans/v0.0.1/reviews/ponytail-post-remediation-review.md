# Ponytail post-remediation review — phase-wise checklist

**Review date:** 2026-07-30

**Comparison:** `main...aaf6880`

**Scope:** Go standards and the prior [phase-wise checklist](ponytail-code-and-architecture-review.md).

**Rating:** **8.2 / 10** — up from 6.6 after the evidence-backed remediation.
**Review boundary:** source review only; no new tests were run for this document.

## Status legend

- [x] Confirmed complete by source and recorded validation evidence.
- [ ] Concrete remaining implementation or test work.
- [~] Intentional policy or conditional follow-up; do not treat as a current defect.

## Phase 0 — validation and scan-result trust

### Validation evidence

- [x] `make lint`, `make test`, and `go test -race ./...` passed for the remediation branch, as recorded in the prior checklist.
- [x] The GopdfSuit no-cache benchmark retained 915 findings; the recorded cache-enabled warm pass was 21.5 ms with 78 hits.
- [x] Cold, warm, and mixed cache scans preserve the inter-file `CWE-22` finding fingerprint: `internal/engine/cache_state_integration_test.go:59`.

### Correctness outcomes

- [x] Construct detector lifecycle state per `AnalyzePaths` call, rather than sharing mutable registry detectors: `internal/engine/registry.go:30`, `internal/engine/registry.go:136`.
- [x] Keep ordinary detector execution off cache hits while accumulating required finalization state: `internal/engine/cache_state_integration_test.go:15`.
- [x] Return `ExitInternal` and write diagnostics to stderr before JSON/SARIF output for invalid encoding and unreadable inputs: `internal/app/run.go:173`, `internal/app/run_test.go:71`.
- [~] Keep Go parse fallback as the current product policy: source-only analysis continues for text detectors when no AST is available (`internal/lang/go/plugin.go:61`). Do not change this merely to satisfy a generic error-handling preference.
- [x] Expose parser quality as a non-fatal `ScanDiagnostic`, retain source-only analysis, preserve the diagnostic through cache hits, and prove malformed-Go JSON output remains usable: `internal/core/plugin.go`, `internal/engine/analyzer.go`, `internal/app/run_test.go`.
- [~] If CI later requires parse completeness, consider a strict mode; the current documented default remains warning-only fallback.

## Phase 1 — preserve output and configuration contracts

### Ownership-safe export

- [x] Delete only exporter-owned numeric context files and `Chunk_*.txt` files; preserve caller files and return cleanup/write errors: `internal/export/export.go:107`, `internal/export/export.go:157`.
- [x] Cover unowned-file preservation, output-directory collision, and context-write failure: `internal/export/export_test.go:79`.
- [x] Add deterministic owned-file removal failure coverage through the narrow `removeOwnedFile` filesystem seam: `internal/export/export_cleanup_test.go`.

### Resolved configuration

- [x] Concentrate CLI/TOML precedence, walk options, cache options, baseline policy, and scan context in `scanPlan`: `internal/app/scan_plan.go:12`.
- [x] Remove the unsupported `languages` and `typed.enabled` configuration promises rather than leaving inert settings.
- [x] Prove that configured `fail_on` has an observable CLI effect: `internal/app/run_test.go:122`.
- [x] Add resolved runtime-field coverage plus an end-to-end `export.whole_function=false` observable-effect test: `internal/app/scan_plan_test.go`, `internal/app/run_test.go`.

### Reporting

- [x] Normalize/copy findings before format-specific reporter work, preserving caller-owned values.
- [x] Keep JSON/SARIF parity and concurrent-reporter coverage recorded in the prior checklist P1.3 evidence.

## Phase 2 — Go module contracts and cold-scan efficiency

### Registry/session contract

- [x] Retain a catalogue for rule listing/discovery while creating session-local detectors for scanning.
- [x] **P-1 — validate catalogue/factory parity.** The registry now rejects a session whose detector count or rule-ID multiplicity differs from its registered catalogue: `internal/engine/registry.go`.
- [x] Add divergent-multiplicity rejection and fresh-instance acceptance coverage: `internal/engine/registry_session_test.go`.

### Project facts

- [x] Keep BP project/package facts scan-local, synchronized, and isolated across roots: `internal/lang/go/detectors/bad_practices/project.go:20`.
- [x] Keep package type facts single-flight per directory using a per-entry `sync.Once`: `internal/lang/go/detectors/bad_practices/package_facts.go:38`.
- [x] **P-2 — make package-doc snapshots single-flight.** Package-doc entries now use per-directory `sync.Once` rather than allowing duplicate cold parallel reads: `internal/lang/go/detectors/bad_practices/project.go`.
- [x] Add a 16-worker regression proving one build and one shared documented snapshot: `internal/lang/go/detectors/bad_practices/project_cache_test.go`.

### Temporary fixtures

- [x] Cover cleanup after materialization failure in the fixture matrix.
- [x] Add an injected analyzer-failure cleanup regression for the temporary matrix directory: `tests/integration/matrix_cleanup_test.go`.

## Phase 3 — scoped maintainability work

- [x] `ScanScope` now owns project-root resolution, cache keys, and display paths: `internal/engine/scan_scope.go:10`.
- [x] The concurrent session regression proves detector state does not leak between analyzers: `internal/engine/scan_session_integration_test.go:13`.
- [~] `app.run` still coordinates cache, baseline, analyzer construction, reporting, export, and exit policy (`internal/app/run.go:31`). This is a valid follow-up, but do not split it until another output/exit-policy variation makes the seam real.
- [x] Do not add a compiled multi-language catalogue yet: Go is the only concrete language adapter, so a second adapter remains hypothetical.

## Phase 4 — rating and next delivery batch

- [x] Re-rate after Phases 0–2 implementation and recorded lint, test, race, and benchmark evidence.
- [x] **Standards score: 8.1 / 10.** Errors, cleanup ownership, scan-local state, and plan/scope locality are materially improved.
- [x] **Checklist/spec score: 8.4 / 10.** Phases 0–2 are substantively complete; remaining rows are narrow follow-ups.
- [x] **Overall score: 8.2 / 10.** There are no confirmed high-severity regressions in the reviewed branch.

### Next batch, in order

1. [x] Implement registry catalogue/factory parity validation and its regression.
2. [x] Make package-doc snapshots single-flight and test cold parallel behavior.
3. [x] Close deterministic export-removal and analyzer-failure fixture-cleanup tests.
4. [x] Broaden observable configuration-effect coverage.
5. [~] Revisit CLI staging only when the second output/exit variation is real.

## Review conclusion

- [x] The earlier scan-trust, destructive-output, configuration-honesty, cache-parity, and shared-state risks are addressed with source and test evidence.
- [x] The completed follow-up batch passed `make lint`, `make test`, and `go test -race ./...`.
- [x] The remaining backlog is phase-scoped and checklist-trackable; it does not justify a broad rewrite.
