# Architecture post-remediation review — phase-wise checklist

**Review date:** 2026-07-30

**Comparison:** `main...aaf6880`

**Scope:** architecture depth, interface, seam, leverage, locality, and deletion-test review.

**Context:** no `CONTEXT.md` or ADR directory exists.

**Rating:** **8.3 / 10**.
**Review boundary:** source review only; no new tests were run for this document.

## Status legend

- [x] Confirmed architectural improvement.
- [ ] Concrete, evidence-backed architectural work.
- [~] Conditional opportunity; wait for the stated real variation before changing the design.

## Phase 0 — confirm the modules that gained depth

- [x] **`scanPlan`:** one interface resolves CLI/TOML precedence before engine construction, concentrating configuration knowledge: `internal/app/scan_plan.go:12`.
- [x] **`scanSession`:** one scan owns fresh detector lifecycle state while the registry remains a catalogue for listing/discovery: `internal/engine/registry.go:30`, `internal/engine/registry.go:136`.
- [x] **`ScanScope`:** cache-key and display-path identity live behind one value interface: `internal/engine/scan_scope.go:10`.
- [x] **Reporting/export:** output modules own normalization, owned-file cleanup, and write errors: `internal/export/export.go:157`.
- [x] **BP facts:** project/package knowledge is scan-local instead of process-global: `internal/lang/go/detectors/bad_practices/project.go:20`.

## Phase 1 — make parser quality an explicit seam

### A-1: parser-quality outcome

- [x] **Strong — return parser quality as an explicit module outcome.** `ParseSource` now returns a `ParseResult` with unit, quality, and non-fatal diagnostic information: `internal/core/plugin.go`, `internal/lang/go/plugin.go`.
- [x] Preserve the current default behaviour: text detectors continue to run on malformed Go.
- [x] Let the engine retain source-only analysis while exposing a `ScanDiagnostic` to stderr; cache entries preserve that diagnostic on warm hits: `internal/engine/analyzer.go`, `internal/engine/cache/types.go`, `internal/app/run.go`.
- [x] Add an end-to-end malformed-Go JSON-output test that proves the warning is deterministic and non-fatal: `internal/app/run_test.go`.
- [x] **Deletion test:** without this module, parse-policy knowledge spreads into each language adapter and output path; it therefore increases leverage and locality.

## Phase 2 — conditionally deepen detector and registry interfaces

### A-2: optional detector capabilities

- [~] **Worth exploring — narrow the broad detector interface.** Stateless detectors currently cross lifecycle, cache-state, finalization, reset, metadata, and run behaviour, while base implementations provide no-ops: `internal/core/detector.go`.
- [x] Add the Ponytail registry catalogue/factory parity regression before considering detector-interface changes: `internal/engine/registry_session_test.go`.
- [~] When a second concrete detector shape proves the variation, keep a narrow detector module and move lifecycle/cache/finalization into optional capability interfaces at the engine seam.
- [x] **Deletion test:** optional capabilities could increase test leverage because simple detector fakes would implement only needed behavior; there is not enough present variation to justify the change now.

### Catalogue decision

- [x] Do not add a compiled multi-language catalogue now. There is one concrete Go adapter and the existing registry/catalogue split is consistent; a second adapter is still hypothetical.

## Phase 3 — stabilize multi-root cache identity

### A-3: canonical key space

- [~] **Worth exploring — make multi-root cache keys order-invariant.** `ScanScope` selects the first requested root (`internal/engine/scan_scope.go:19`), and analyzer setup repeats the same policy (`internal/engine/analyzer.go:185`). Equivalent root lists can therefore produce different persisted identities and `../` keys.
- [ ] When multi-root cache reuse becomes a real workflow concern, make `ScanScope` the only owner of a canonical multi-root key space and pass it through the analyzer builder instead of a raw root string.
- [ ] Before implementation, add acceptance tests for root-order invariance, mixed file/directory roots, and persisted-cache reuse.
- [x] **Deletion test:** centralizing this policy prevents cache maintenance and scan execution from relearning ordering rules, improving locality.

## Phase 4 — stage CLI orchestration only on real variation

### A-4: execution/output staging

- [x] Resolved-plan creation and scan scope have already moved out of the CLI entry point.
- [~] **Worth exploring — retain one `Run(args)` interface, but later split private scan-execution and output/exit-policy modules.** `run` still coordinates cache, baseline, analyzer construction, reporting, export, and exit policy: `internal/app/run.go:31`.
- [ ] Revisit only when a second output or exit-policy variation appears. Do not introduce an adapter for a hypothetical variation.
- [x] **Deletion test:** if output policy begins to vary, staging will concentrate change and verification. Current tests are sufficient, so immediate extraction is not warranted.

## Phase 5 — rating and delivery order

- [x] **Architecture score: 8.3 / 10.** `scanPlan`, `scanSession`, `ScanScope`, output ownership, and scan-local facts materially improve depth, leverage, and locality.
- [x] A-1 is the only strong candidate because it makes an existing intentional policy observable without removing text-detector value.
- [x] A-2, A-3, and A-4 are deliberately conditional; they are not generic refactoring debt.

### Recommended delivery order

1. [x] A-1 parser-quality outcome and malformed-Go diagnostic test.
2. [x] Registry catalogue/factory parity validation from the Ponytail review.
3. [~] Revisit A-2 only with a second detector shape or language adapter.
4. [~] Revisit A-3 only with observed multi-root cache reuse friction.
5. [~] Revisit A-4 only with a second output/exit-policy variation.

## Review conclusion

- [x] The remediation created deeper modules at the configuration, scan-lifetime, path-identity, output, and project-facts seams.
- [x] The parser-quality and catalogue/factory seams are now explicit and validation-backed.
- [x] The remaining architecture ledger is phased, evidence-backed, and intentionally avoids a generic clean/hexagonal rewrite.
