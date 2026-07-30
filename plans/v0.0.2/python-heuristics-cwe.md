# v0.0.2 / #52 — Python CWE heuristic detectors

> **Parent:** `plans/v0.0.2/python-heuristics.md` — epic #51 rollup; issue body `plans/PR/v0.0.2/issue-python-cwe-heuristics-body.md`
> **Status:** `[ ]` not started (foundation #39 / PR #50 shipped; Python plugin still zero detectors)
> **Estimated effort:** medium (1 PR for framework + priority batch; expand later)

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

- [ ] Create `internal/lang/python/detectors/cwe/` package with package comment stating pure-Go / source-pattern scope (no CGO tree-sitter).
- [ ] Add `internal/lang/python/detectors/all.go` (or equivalent) that will return `[]core.Detector` for the Python plugin (initially CWE-only).
- [ ] Document non-goals in `internal/lang/python/detectors/cwe/README.md` (one short file): priority batch IDs; full 344 deferred; Go CWE untouched.

### 1.2 Unified scan + registration API

- [ ] Implement `PyCweScan` (`core.Detector`) with `Language() == core.LanguagePython`, `RuleIDs()`, `MetadataFor()`, `Run()`.
- [ ] Implement `RegisterRule(id, fn, meta, gates...)` (or init-table equivalent) so new CWEs register without changing `Run` control flow — proof: adding a no-op rule increases `RuleIDs()` by one in unit test.
- [ ] `Run` honors `ctx.Allows(ruleID)` and skips when unit/source empty (nil-safe).
- [ ] Build facts once per unit: `BuildFacts(unit)` using `ast.Build(source, pyCweNeedles)` into a pack-local fact bag (do not reuse Go `cweNeedles`).

### 1.3 Metadata for priority batch

- [ ] Hand-authored (or small local embed) `MetaCWE22`, `MetaCWE78`, `MetaCWE79`, `MetaCWE89`, `MetaCWE502` via `rules.Meta(...)` with titles/names matching `ruleset/python/chunks` for those IDs.
- [ ] Assert IDs exist in catalogue JSON: CWE-22 in `cwe-001-050.json`; CWE-78/79/89 in `cwe-051-100.json`; CWE-502 in `cwe-501-550.json` (`python_relevance` High for 502).
- [ ] Confirm `PackFromRuleID("CWE-*") == PackSecurity` for each meta (no custom pack hacks).
- [ ] Do **not** wire full 344 metadata generation in this phase (`[~]` full embed of all chunks deferred — framework must accept more metas later).

### 1.4 Needle / prefilter table (Python sinks)

- [ ] Define `pyCweNeedles` covering gates for priority rules (examples: `subprocess.`, `os.system`, `os.popen`, `pickle.loads`, `pickle.load`, `yaml.load`, `cursor.execute`, `execute(`, `open(`, `pathlib`, `os.path.join`, `mark_safe`, `Markup(`, `render_template_string`, `HttpResponse(`).
- [ ] Gate each registered rule with FN-safe any-of needles (skip file when no sink tokens present).

### 1.5 Infrastructure unit tests

- [ ] `go test ./internal/lang/python/detectors/cwe/` — empty source emits zero findings.
- [ ] `go test ./internal/lang/python/detectors/cwe/` — catalogue lists exactly the registered priority IDs after Phase 2 registration (update expected set when rules land).
- [ ] Language ID on detector is `LanguagePython` (not Go).

---

## Phase 2: Priority CWE heuristics (security correctness)

> Prefer high-signal, low-FP patterns. Same-file heuristics only (no interproc taint graph).  
> Emit via `rules.PushFindingWithConfidence` / `NewFindingFromMeta` with stable messages.

### 2.1 CWE-502 — Deserialization of Untrusted Data

> Catalogue: `ruleset/python/chunks/cwe-501-550.json` · `python_relevance: High` · Python-platform CWE.

- [ ] Register `CWE-502` on `PyCweScan` with meta + gates (`pickle.`, `yaml.load`, etc.).
- [ ] Positive heuristic: flag `pickle.loads(` / `pickle.load(` / `pickle.Unpickler(` on non-constant / request-like data (or any non-literal arg when conservative mode matches Go museum style — document choice).
- [ ] Positive heuristic: flag `yaml.load(` without `Loader=yaml.SafeLoader` / `FullLoader` safe pattern (and/or `yaml.unsafe_load`).
- [ ] Negative: `yaml.safe_load(` and `json.loads(` must not fire CWE-502.
- [ ] Unit test: vulnerable snippet → ≥1 finding `RuleID == "CWE-502"`; safe snippet → 0 for CWE-502.
- [ ] Finding severity High; message names unsafe deserialization sink.

### 2.2 CWE-78 — OS command injection

> Catalogue: `cwe-051-100.json`. Go reference: `detectCWE78` / `exec.Command` — reimplement for Python sinks only.

- [ ] Register `CWE-78` with gates on `subprocess`, `os.system`, `os.popen`, `commands.`.
- [ ] Positive: `os.system(` with non-literal / formatted command; `subprocess.*(..., shell=True)` with dynamic command; `subprocess.call/run/Popen` with string cmd + shell.
- [ ] Negative: `subprocess.run(["ls", "-l"], shell=False)` / list argv with literals only → no CWE-78.
- [ ] Unit test hit + miss for `CWE-78`.
- [ ] Message: user-controlled / dynamic input reaches OS command sink.

### 2.3 CWE-89 — SQL injection

> Catalogue: `cwe-051-100.json`. Go reference: `detectCWE89` — Python DB-API / driver sinks.

- [ ] Register `CWE-89` with gates on `execute(`, `executemany(`, `raw(`, SQLAlchemy text patterns as needed.
- [ ] Positive: `cursor.execute("... %s ..." % var)` / f-string / `.format` SQL; `execute(f"SELECT ... {x}")`.
- [ ] Negative: parameterized `cursor.execute("SELECT ... WHERE id = ?", (x,))` or `%s` with bound args tuple only → no CWE-89.
- [ ] Unit test hit + miss for `CWE-89`.

### 2.4 CWE-22 — Path traversal

> Catalogue: `cwe-001-050.json`. Go reference: `detectCWE22` file sinks without confinement.

- [ ] Register `CWE-22` with gates on `open(`, `pathlib`, `os.path.join`, `os.remove`, etc.
- [ ] Positive: `open(os.path.join(root, user_path))` / `Path(root) / request_path` without `resolve`+prefix check or `basename`-only confinement.
- [ ] Negative: path built only from literals, or explicit confinement (`startswith` root after `resolve`, `os.path.basename` only policy) — document which safe patterns suppress.
- [ ] Unit test hit + miss for `CWE-22`.

### 2.5 CWE-79 — XSS

> Catalogue: `cwe-051-100.json`. Framework sinks (Flask/Django/Jinja), not browser DOM.

- [ ] Register `CWE-79` with gates on `mark_safe`, `Markup(`, `render_template_string`, `|safe`, `HttpResponse(`.
- [ ] Positive: `mark_safe(user)` / `Markup(request...)` / `render_template_string` with dynamic HTML; unescaped response body from request data.
- [ ] Negative: autoescaped `render_template("x.html", name=name)` without `|safe` / `mark_safe` → no CWE-79.
- [ ] Unit test hit + miss for `CWE-79`.

### 2.6 Batch integrity

- [ ] All five IDs registered exactly once (no double `RegisterRule` for same ID).
- [ ] Each rule FN-safe: files without gate needles produce zero findings for that rule.
- [ ] No findings use Go language ID or Go-only messages (gin/gorm wording).
- [ ] Table test matrix: 5 rules × {vulnerable, safe} snippets in `*_test.go` under `detectors/cwe/`.

---

## Phase 3: Plugin wiring (API / data contracts)

### 3.1 Python plugin catalogue

- [ ] `internal/lang/python/plugin.go` — `Detectors()` returns CWE scan detector(s) via `detectors.All()` (non-empty).
- [ ] `NewDetectors()` returns **fresh** instances matching catalogue length (session-local; no shared mutable detector state).
- [ ] Keep `ParseSource` source-only / `LanguagePython` unless a later pure-Go parse is explicitly chosen; detectors must not require `unit.Tree`.
- [ ] Update package comment on `plugin.go` (remove “zero detectors” claim once wired).

### 3.2 Tests that currently require empty catalogue

- [ ] Update `internal/lang/python/plugin_test.go` — stop requiring `len(Detectors()) == 0`; assert priority CWE IDs present instead.
- [ ] Update `internal/engine/registry_test.go` expectations that `DetectorsForLanguage(LanguagePython) == 0` when Python plugin is registered.
- [ ] Update `internal/app/run_test.go` `TestRunListRulesRespectsLanguagesConfig` — `languages=["python"]` must list `CWE-22`/`CWE-78`/… (not “no rules registered”).
- [ ] DefaultRegistry remains Go-only (Python still opt-in) — existing `TestDefaultRegistryRemainsGoOnly` stays green.

### 3.3 Finding language / severity smoke

- [ ] Unit or small integration: run `PyCweScan` on a Python `ParsedUnit` → findings have `RuleID` prefix `CWE-` and severity ≥ Medium for priority batch.
- [ ] Confirm multi-language registry: `NewRegistryWithLanguages(go, python)` still resolves `.py` to Python plugin with non-empty detectors.

---

## Phase 4: Fixtures (hit / miss corpus)

> Layout: `tests/fixtures/python/` text fixtures (`lang: python`, materialize to `.py`) — match existing `sample.txt` / `safe.txt` header style.

### 4.1 Per-rule fixtures

- [ ] `tests/fixtures/python/CWE-502-vulnerable.txt` — pickle/yaml unsafe load fires CWE-502.
- [ ] `tests/fixtures/python/CWE-502-safe.txt` — safe_load / json only; no CWE-502.
- [ ] `tests/fixtures/python/CWE-78-vulnerable.txt` — shell/command injection pattern.
- [ ] `tests/fixtures/python/CWE-78-safe.txt` — list argv / no shell.
- [ ] `tests/fixtures/python/CWE-89-vulnerable.txt` — string-built SQL execute.
- [ ] `tests/fixtures/python/CWE-89-safe.txt` — parameterized query.
- [ ] `tests/fixtures/python/CWE-22-vulnerable.txt` — path join + open without confinement.
- [ ] `tests/fixtures/python/CWE-22-safe.txt` — confined or literal paths.
- [ ] `tests/fixtures/python/CWE-79-vulnerable.txt` — mark_safe / Markup / template_string XSS.
- [ ] `tests/fixtures/python/CWE-79-safe.txt` — autoescaped template render.

### 4.2 Fixture runner tests

- [ ] Detector or package test materializes / reads each fixture and asserts rule ID presence/absence (mirror Go `registry_test.go` vulnerable/safe pairs).
- [ ] Fixtures stay under `tests/fixtures/python/` (or documented sibling); paths included in FP-suppression exemptions if any Go-style real-project filters are copied (prefer not copying FP museum filters into Python yet).

---

## Phase 5: End-to-end scan + expansion framework

### 5.1 CLI / config path

- [ ] With `languages = ["python"]` and fixture tree, `goslop --format json --no-cache <dir>` reports ≥1 CWE finding for vulnerable fixtures.
- [ ] `languages = ["go"]` on a pure-Python tree still scans 0 Python files (unchanged product default).
- [ ] `languages = ["go","python"]` scans both; Python CWE IDs do not break Go catalogue counts.

### 5.2 Expansion framework (not full 344)

- [ ] Document how to add CWE-N: catalogue row exists → meta + `RegisterRule` + needles + tests (README in detectors/cwe).
- [ ] `[~]` Implement remaining high-`python_relevance` IDs beyond the five — **deferred** to follow-up PRs; owner #52 backlog / epic #51; next gate: pick next batch from `python_relevance` + `detection_notes` (e.g. CWE-798, CWE-327/328 class crypto, CWE-611) with fixtures.
- [ ] `[~]` Full 344-rule registration — **deferred** (issue scope); do not block ship of priority batch.
- [ ] `[~]` Inter-procedural taint for Python 22/78/79/89 — **deferred** (Go has separate taint package; not required for #52).

### 5.3 Docs touchpoints (minimal)

- [ ] `ruleset/python/README.md` — note that priority CWE detectors exist under `internal/lang/python/detectors/cwe` (catalogue still source of ID truth).
- [ ] Parent ledger `plans/v0.0.2/python-heuristics.md` CWE rows flipped only after Phase 6 evidence.

---

## Phase 6: Validation gates + issue closure

> Per `plans/skills/phase-wise-checklist/SKILLS.md`: leave unchecked until commands succeed on the implement branch; record outcomes beside the row.

### 6.1 Local gates

- [ ] `gofmt -w` on all changed Go files under `internal/lang/python/**`.
- [ ] `make lint` — green (record date / branch in PR or ledger when done).
- [ ] `make test` — green (record date / branch when done).
- [ ] `CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop` — succeeds (proves no CGO requirement).

### 6.2 Success criteria (#52)

- [ ] Priority batch CWE-22, 78, 79, 89, 502 registered and fire on fixtures.
- [ ] `languages = ["python"]` scan reports CWE findings for fixture tree.
- [ ] Tests cover each implemented rule (positive + ≥1 negative).
- [ ] Catalogue IDs aligned with `ruleset/python/chunks` for implemented rules.
- [ ] Go CWE detectors / `ruleset/golang/**` unchanged.
- [ ] Issue #52 closed or partial `[~]` list of deferred IDs recorded in this ledger + PR body.

### 6.3 Handoff

- [ ] Sync checkboxes in `plans/v0.0.2/python-heuristics.md` CWE rollup with evidence.
- [ ] PR body under `plans/PR/v0.0.2/` (when opened): `Closes #52`, `Relates to #51`.
- [ ] No claim of full 344 coverage or BP/PERF in the CWE PR.

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
