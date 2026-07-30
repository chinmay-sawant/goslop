# v0.0.2 / #52 — Python CWE heuristic detectors

> **Parent:** `plans/v0.0.2/heuristics/python-heuristics.md` — epic #51 rollup; issue body `plans/PR/v0.0.2/issues/issue-python-cwe-heuristics-body.md`
> **Status:** priority batch **shipped** on `main`; remaining catalogue expansion planned in [cwe-plans/](./cwe-plans/) (2026-07-31)
> **Estimated effort:** multi-PR expansion (see cwe-plans waves P0–P3); priority 5 already done
> **Batch index:** [cwe-plans/README.md](./cwe-plans/README.md) · inventory [cwe-plans/_inventory.json](./cwe-plans/_inventory.json)

---

## Overview

Turn `ruleset/python/chunks/cwe-*.json` (~344 CWE rules, `python_relevance`, `detection_notes`) into **runnable** Python CWE heuristics that emit `CWE-*` findings.

| Current fact | Evidence |
|--------------|----------|
| Python plugin source-only, **zero** detectors | `internal/lang/python/plugin.go` — `Detectors()` / `NewDetectors()` return `nil` |
| Plugin tests assert empty catalogue | `internal/lang/python/plugin_test.go` |
| CWE catalogue present | `ruleset/python/chunks/cwe-*.json` (~344); mapping `plans/v0.0.2/python-cwe-from-699-mapping.md` |
| Priority IDs in catalogue | CWE-22 (`cwe-001-050.json`), CWE-78/79/89 (`cwe-051-100.json`), CWE-502 (`cwe-501-550.json`, `python_relevance: High`) |
| Go reference (do not modify) | `internal/lang/go/detectors/cwe/` — `GoCweScan`, `RegisterRule`, `SourceIndex` facts, taint-lite for 22/78/79/89 |
| Fixtures today | `tests/fixtures/python/{sample,safe}.txt` only (no CWE hit/miss) |
| Product default | Go-only; Python opt-in via `languages = ["python"]` |

**In scope (#52):** priority batch **CWE-22, 78, 79, 89, 502** + detector framework for more later.  
**Out of scope:** full 344-rule ship; BP (#53) / PERF (#54); changing Go CWE detectors; CGO / tree-sitter requirement.

---

## Executive Summary

### Strategy

1. **Correctness first** — pure-Go pattern / light source analysis (mirror Go `SourceIndex` + sink heuristics), no CGO tree-sitter.
2. **Priority security batch** — CWE-502 (Python-platform / High relevance), then 78 / 89 / 22 / 79 with Python stdlib/framework sinks.
3. **Wire** into `internal/lang/python` so enabled-language scans emit findings.
4. **Fixtures** — positive + negative per implemented rule under `tests/fixtures/python/`.
5. **Gates** — `make lint` / `make test` remain unchecked until green on the implement branch.

### Parser / approach decision (document in code README when implementing)

| Approach | Use for |
|----------|---------|
| `ast.Build` / `SourceIndex` needles + string/call heuristics | Default for v0.0.2 batch (same idea as `internal/lang/go/detectors/cwe/facts.go`) |
| Light pure-Go token/line scan (no full Python AST) | Call-site sinks: `subprocess.*`, `os.system`, `pickle.loads`, `cursor.execute`, `open(`, Jinja/Flask markers |
| Full Python AST / tree-sitter CGO | **Not required** for #52; defer |

### Finding contract (must hold for every rule)

- Rule ID stable: `CWE-22` … `CWE-502` (aligned with catalogue keys).
- `rules.Meta` / `PackFromRuleID` → `PackSecurity` for `CWE-*`.
- Finding language context is Python unit (`ParsedUnit.Language == LanguagePython`); plugin must keep source-only parse if no tree.
- Severity: High for priority batch (match Go meta defaults for these IDs unless catalogue later overrides).
- Do **not** point Go `metadata_gen` at `ruleset/python/`.

### Dependency order

```text
Phase 1 framework  →  Phase 2 detectors  →  Phase 3 plugin wire
                                           →  Phase 4 fixtures
                                           →  Phase 5 e2e / expand hooks
                                           →  Phase 6 validation gates
```

---

## Phase 1: CWE detector infrastructure (correctness foundation)

> Target package (proposed): `internal/lang/python/detectors/cwe/`  
> Mirror shape of `internal/lang/go/detectors/cwe/{scan,facts,register,metadata}.go` without importing Go detector packages as a runtime dependency.

### 1.1 Package skeleton

- [x] Create `internal/lang/python/detectors/cwe/` package with package comment stating pure-Go / source-pattern scope (no CGO tree-sitter).  
  **Evidence:** package doc on `scan.go`; README non-goals (2026-07-31, `feat/python-cwe-heuristics`).
- [x] Add `internal/lang/python/detectors/all.go` (or equivalent) that will return `[]core.Detector` for the Python plugin (initially CWE-only).  
  **Evidence:** `detectors.All()` returns `cwe.NewPyCweScan()` only.
- [x] Document non-goals in `internal/lang/python/detectors/cwe/README.md` (one short file): priority batch IDs; full 344 deferred; Go CWE untouched.

### 1.2 Unified scan + registration API

- [x] Implement `PyCweScan` (`core.Detector`) with `Language() == core.LanguagePython`, `RuleIDs()`, `MetadataFor()`, `Run()`.
- [x] Implement `RegisterRule(id, fn, meta, gates...)` (or init-table equivalent) so new CWEs register without changing `Run` control flow — proof: catalogue test asserts 5 priority IDs.
- [x] `Run` honors `ctx.Allows(ruleID)` and skips when unit/source empty (nil-safe).  
  **Evidence:** `TestAllowsSkipsDisallowedRules`, `TestEmptySourceEmitsZeroFindings`.
- [x] Build facts once per unit: `BuildFacts(unit)` using `ast.Build(source, pyCweNeedles)` into a pack-local fact bag (do not reuse Go `cweNeedles`).

### 1.3 Metadata for priority batch

- [x] Hand-authored (or small local embed) `MetaCWE22`, `MetaCWE78`, `MetaCWE79`, `MetaCWE89`, `MetaCWE502` via `rules.Meta(...)` with titles/names matching `ruleset/python/chunks` for those IDs.
- [x] Assert IDs exist in catalogue JSON: CWE-22 in `cwe-001-050.json`; CWE-78/79/89 in `cwe-051-100.json`; CWE-502 in `cwe-501-550.json` (`python_relevance` High for 502).
- [x] Confirm `PackFromRuleID("CWE-*") == PackSecurity` for each meta (no custom pack hacks).  
  **Evidence:** `TestPyCweScanLanguageAndCatalogue`.
- [x] Do **not** wire full 344 metadata generation in this phase (`[~]` full embed of all chunks deferred — framework must accept more metas later).

### 1.4 Needle / prefilter table (Python sinks)

- [x] Define `pyCweNeedles` covering gates for priority rules.
- [x] Gate each registered rule with FN-safe any-of needles (skip file when no sink tokens present).

### 1.5 Infrastructure unit tests

- [x] `go test ./internal/lang/python/detectors/cwe/` — empty source emits zero findings.
- [x] `go test ./internal/lang/python/detectors/cwe/` — catalogue lists exactly the registered priority IDs after Phase 2 registration.
- [x] Language ID on detector is `LanguagePython` (not Go).

---

## Phase 2: Priority CWE heuristics (security correctness)

> Prefer high-signal, low-FP patterns. Same-file heuristics only (no interproc taint graph).  
> Emit via `rules.PushFindingWithConfidence` / `NewFindingFromMeta` with stable messages.

### 2.1 CWE-502 — Deserialization of Untrusted Data

> Catalogue: `ruleset/python/chunks/cwe-501-550.json` · `python_relevance: High` · Python-platform CWE.

- [x] Register `CWE-502` on `PyCweScan` with meta + gates (`pickle.`, `yaml.load`, etc.).
- [x] Positive heuristic: flag `pickle.loads(` / `pickle.load(` / `pickle.Unpickler(` on non-constant / request-like data (or any non-literal arg when conservative mode matches Go museum style — document choice).
- [x] Positive heuristic: flag `yaml.load(` without `Loader=yaml.SafeLoader` / `FullLoader` safe pattern (and/or `yaml.unsafe_load`).
- [x] Negative: `yaml.safe_load(` and `json.loads(` must not fire CWE-502.
- [x] Unit test: vulnerable snippet → ≥1 finding `RuleID == "CWE-502"`; safe snippet → 0 for CWE-502.
- [x] Finding severity High; message names unsafe deserialization sink.

### 2.2 CWE-78 — OS command injection

> Catalogue: `cwe-051-100.json`. Go reference: `detectCWE78` / `exec.Command` — reimplement for Python sinks only.

- [x] Register `CWE-78` with gates on `subprocess`, `os.system`, `os.popen`, `commands.`.
- [x] Positive: `os.system(` with non-literal / formatted command; `subprocess.*(..., shell=True)` with dynamic command; `subprocess.call/run/Popen` with string cmd + shell.
- [x] Negative: `subprocess.run(["ls", "-l"], shell=False)` / list argv with literals only → no CWE-78.
- [x] Unit test hit + miss for `CWE-78`.
- [x] Message: user-controlled / dynamic input reaches OS command sink.

### 2.3 CWE-89 — SQL injection

> Catalogue: `cwe-051-100.json`. Go reference: `detectCWE89` — Python DB-API / driver sinks.

- [x] Register `CWE-89` with gates on `execute(`, `executemany(`, `raw(`, SQLAlchemy text patterns as needed.
- [x] Positive: `cursor.execute("... %s ..." % var)` / f-string / `.format` SQL; `execute(f"SELECT ... {x}")`.
- [x] Negative: parameterized `cursor.execute("SELECT ... WHERE id = ?", (x,))` or `%s` with bound args tuple only → no CWE-89.
- [x] Unit test hit + miss for `CWE-89`.

### 2.4 CWE-22 — Path traversal

> Catalogue: `cwe-001-050.json`. Go reference: `detectCWE22` file sinks without confinement.

- [x] Register `CWE-22` with gates on `open(`, `pathlib`, `os.path.join`, `os.remove`, etc.
- [x] Positive: `open(os.path.join(root, user_path))` / `Path(root) / request_path` without `resolve`+prefix check or `basename`-only confinement.
- [x] Negative: path built only from literals, or explicit confinement (`startswith` root after `resolve`, `os.path.basename` only policy) — document which safe patterns suppress.
- [x] Unit test hit + miss for `CWE-22`.

### 2.5 CWE-79 — XSS

> Catalogue: `cwe-051-100.json`. Framework sinks (Flask/Django/Jinja), not browser DOM.

- [x] Register `CWE-79` with gates on `mark_safe`, `Markup(`, `render_template_string`, `|safe`, `HttpResponse(`.
- [x] Positive: `mark_safe(user)` / `Markup(request...)` / `render_template_string` with dynamic HTML; unescaped response body from request data.
- [x] Negative: autoescaped `render_template("x.html", name=name)` without `|safe` / `mark_safe` → no CWE-79.
- [x] Unit test hit + miss for `CWE-79`.

### 2.6 Batch integrity

- [x] All five IDs registered exactly once (no double `RegisterRule` for same ID).
- [x] Each rule FN-safe: files without gate needles produce zero findings for that rule.
- [x] No findings use Go language ID or Go-only messages (gin/gorm wording).
- [x] Table test matrix: 5 rules × {vulnerable, safe} snippets in `*_test.go` under `detectors/cwe/`.

---

## Phase 3: Plugin wiring (API / data contracts)

### 3.1 Python plugin catalogue

- [x] `internal/lang/python/plugin.go` — `Detectors()` returns CWE scan detector(s) via `detectors.All()` (non-empty).
- [x] `NewDetectors()` returns **fresh** instances matching catalogue length (session-local; no shared mutable detector state).
- [x] Keep `ParseSource` source-only / `LanguagePython` unless a later pure-Go parse is explicitly chosen; detectors must not require `unit.Tree`.
- [x] Update package comment on `plugin.go` (remove “zero detectors” claim once wired).

### 3.2 Tests that currently require empty catalogue

- [x] Update `internal/lang/python/plugin_test.go` — stop requiring `len(Detectors()) == 0`; assert priority CWE IDs present instead.
- [x] Update `internal/engine/registry_test.go` expectations that `DetectorsForLanguage(LanguagePython) == 0` when Python plugin is registered.
- [x] Update `internal/app/run_test.go` `TestRunListRulesRespectsLanguagesConfig` — `languages=["python"]` must list `CWE-22`/`CWE-78`/… (not “no rules registered”).
- [x] DefaultRegistry remains Go-only (Python still opt-in) — existing `TestDefaultRegistryRemainsGoOnly` stays green.

### 3.3 Finding language / severity smoke

- [x] Unit or small integration: run `PyCweScan` on a Python `ParsedUnit` → findings have `RuleID` prefix `CWE-` and severity ≥ Medium for priority batch.
- [x] Confirm multi-language registry: `NewRegistryWithLanguages(go, python)` still resolves `.py` to Python plugin with non-empty detectors.

---

## Phase 4: Fixtures (hit / miss corpus)

> Layout: `tests/fixtures/python/` text fixtures (`lang: python`, materialize to `.py`) — match existing `sample.txt` / `safe.txt` header style.

### 4.1 Per-rule fixtures

- [x] `tests/fixtures/python/CWE-502-vulnerable.txt` — pickle/yaml unsafe load fires CWE-502.
- [x] `tests/fixtures/python/CWE-502-safe.txt` — safe_load / json only; no CWE-502.
- [x] `tests/fixtures/python/CWE-78-vulnerable.txt` — shell/command injection pattern.
- [x] `tests/fixtures/python/CWE-78-safe.txt` — list argv / no shell.
- [x] `tests/fixtures/python/CWE-89-vulnerable.txt` — string-built SQL execute.
- [x] `tests/fixtures/python/CWE-89-safe.txt` — parameterized query.
- [x] `tests/fixtures/python/CWE-22-vulnerable.txt` — path join + open without confinement.
- [x] `tests/fixtures/python/CWE-22-safe.txt` — confined or literal paths.
- [x] `tests/fixtures/python/CWE-79-vulnerable.txt` — mark_safe / Markup / template_string XSS.
- [x] `tests/fixtures/python/CWE-79-safe.txt` — autoescaped template render.

### 4.2 Fixture runner tests

- [x] Detector or package test materializes / reads each fixture and asserts rule ID presence/absence (mirror Go `registry_test.go` vulnerable/safe pairs).
- [x] Fixtures stay under `tests/fixtures/python/` (or documented sibling); paths included in FP-suppression exemptions if any Go-style real-project filters are copied (prefer not copying FP museum filters into Python yet).

---

## Phase 5: End-to-end scan + expansion framework

### 5.1 CLI / config path

- [x] With `languages = ["python"]` and fixture tree, `goslop --format json --no-cache <dir>` reports ≥1 CWE finding for vulnerable fixtures.
- [x] `languages = ["go"]` on a pure-Python tree still scans 0 Python files (unchanged product default).
- [x] `languages = ["go","python"]` scans both; Python CWE IDs do not break Go catalogue counts.

### 5.2 Expansion framework (not full 344)

- [x] Document how to add CWE-N: catalogue row exists → meta + `RegisterRule` + needles + tests (README in detectors/cwe).
- [x] **Batch expansion plans** written under [cwe-plans/](./cwe-plans/) (mirror of `bp-plans/`) after 5-agent scan of `ruleset/python/chunks/` (2026-07-31 on `main` post-PR #67).  
  **Partition:** 5 shipped · **154** implement-owned (batches 01–14, 16) · **185** deferred (tier-C / no pure-Go sink).  
  **Tiers:** A≈59 · B≈94 · C≈186 — success ≠ full 344 `RegisterRule`.  
  **Testing contract (every implement ID):** unit hit/miss in `internal/lang/python/detectors/cwe/*_test.go` + fixture pair `tests/fixtures/python/cwe/CWE-N-{vulnerable,safe}.txt` (same `.txt` layout as `tests/fixtures/python/bp/`) + `TestPythonCWEFixturesMatrix` / `make integration-python` — see [cwe-plans/README.md § Testing contract](./cwe-plans/README.md).
- [ ] Execute **wave P0** implement batches: [batch-01](./cwe-plans/batch-01-injection-expand.md) … [batch-04](./cwe-plans/batch-04-secrets-credentials.md) (one PR each; each ID ships unit + fixtures + matrix).
- [ ] Execute **wave P1** (crypto / SSRF / XXE / path-fs): batches 05–08.
- [ ] Execute **wave P2** (web config / authz / info / validation): batches 09–12.
- [ ] Execute **wave P3** (platform quality / resources / tier-B expansion): batches 13–14, 16 (split 16 if large).
- [~] Full 344-rule registration — **not a goal**; remainder lives in [batch-deferred.md](./cwe-plans/batch-deferred.md).
- [~] Inter-procedural taint for Python 22/78/79/89 — **deferred** (Go has separate taint package; not required for expansion waves).

### 5.3 Docs touchpoints (minimal)

- [x] `ruleset/python/README.md` — note that priority CWE detectors exist under `internal/lang/python/detectors/cwe` (catalogue still source of ID truth).
- [x] Parent ledger `plans/v0.0.2/heuristics/python-heuristics.md` CWE rows flipped with Phase 6 evidence.

---

## Phase 6: Validation gates + issue closure

> Per `plans/skills/phase-wise-checklist/SKILLS.md`: leave unchecked until commands succeed on the implement branch; record outcomes beside the row.

### 6.1 Local gates

- [x] `gofmt -w` on all changed Go files under `internal/lang/python/**`.
- [x] `make lint` — green (record date / branch in PR or ledger when done).
- [x] `make test` — green (record date / branch when done).
- [x] `CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop` — succeeds (proves no CGO requirement).

### 6.2 Success criteria (#52)

- [x] Priority batch CWE-22, 78, 79, 89, 502 registered and fire on fixtures.
- [x] `languages = ["python"]` scan reports CWE findings for fixture tree.
- [x] Tests cover each implemented rule (positive + ≥1 negative).
- [x] Catalogue IDs aligned with `ruleset/python/chunks` for implemented rules.
- [x] Go CWE detectors / `ruleset/golang/**` unchanged.
- [x] Issue #52 closed or partial `[~]` list of deferred IDs recorded in this ledger + PR body.  
  **Deferred:** remaining high-`python_relevance` CWEs beyond the five; full 344; interproc taint; BP/PERF.

### 6.3 Handoff

- [x] Sync checkboxes in `plans/v0.0.2/python-heuristics.md` CWE rollup with evidence.
- [x] PR body under `plans/PR/v0.0.2/PR/pr-python-cwe-heuristics.md`: `Closes #52`, `Relates to #51`.
- [x] No claim of full 344 coverage or BP/PERF in the CWE PR.

---

## Dependencies

| Depends on | Note |
|------------|------|
| Epic #39 / PR #50 | Python plugin stub, `languages` config, catalogues |
| `plans/v0.0.2/python-heuristics.md` | Parent epic #51 ledger |
| `internal/core.LanguagePlugin` | `Detectors` / `NewDetectors` / `ParseSource` |
| `internal/ast.SourceIndex` | Needle prefilter without tree-sitter |
| `internal/rules` | `Meta`, `Finding`, `PackSecurity` via `CWE-*` prefix |
| `ruleset/python/chunks/cwe-*.json` | ID / name / relevance source of truth |
| Go CWE package | **Reference only** — do not modify for #52 |
| Sibling #53 BP | Shares plugin `detectors.All()` surface; coordinate registration order when both land |
| Sibling #54 PERF | Deferred catalogue; independent of this checklist |

### Explicit non-dependencies

- Tree-sitter / CGO Python parser  
- Go taint package port  
- Full 344 CWE implementation  
- BP / PERF heuristics  
