# v0.0.2 / #52 — Python CWE batch plans (execution index)

> **Parent:** `plans/v0.0.2/heuristics/python-heuristics-cwe.md` — canonical #52 ledger  
> **Epic:** [#51](https://github.com/chinmay-sawant/goslop/issues/51)  
> **Issue:** [#52](https://github.com/chinmay-sawant/goslop/issues/52) expansion — batchwise PRs for remaining catalogue CWEs  
> **Status:** **planning** — priority 5 CWEs shipped; 339 pending partitioned into implement + deferred  
> **Inventory snapshot:** `plans/v0.0.2/heuristics/cwe-plans/_inventory.json`  
> **Catalogue:** `ruleset/python/chunks/cwe-*.json` (**344** rules)  
> **Scan method:** 5 parallel agents (2026-07-31) over chunks + Go CWE parity hints

---

## Overview

Batch plans under this directory are **live execution ledgers** for #52 follow-up work (same shape as `bp-plans/`). Each batch file follows `plans/skills/phase-wise-checklist/SKILLS.md`:

- Atomic `[ ]` / `[x]` / `[~]` rows with path, rule ID, expected behavior, proof.
- Mark `[x]` only after detector + hit/miss proof and (for code) `make lint` + `make test`.
- One **batchwise PR** per implement batch file (do not merge unrelated waves unless tiny).

### Coverage honesty

| Bucket | Count | Meaning |
|--------|------:|---------|
| Implemented (batch-00) | **5** | CWE-22, 78, 79, 89, 502 on `main` |
| Missing / implement batches | **154** | Owned by batch-01..14 + 16 (tier A+B) |
| Deferred (batch-deferred) | **185** | No pure-Go sink heuristic in v0 — catalogue-only |
| **Catalogue total** | **344** | Partition is complete |

| Tier | Count | Definition |
|------|------:|------------|
| **A** | **59** | High-signal pure-Go source heuristic — implement first |
| **B** | **94** | Feasible later; higher FP / niche |
| **C** | **186** | Design/process / no clear line-level sink → deferred |

**Success criterion for expansion waves:** ship tier-A (and valuable B) with fixtures — **not** full 344 `RegisterRule`.

### Architecture (all batches)

| Constraint | Detail |
|------------|--------|
| Package | `internal/lang/python/detectors/cwe/` |
| Registration | `RegisterRule(id, fn, meta, gates...)` from `init()` in domain rule files |
| Detector surface | single `PyCweScan` (`scan.go`); **do not** add a second `core.Detector` |
| Detection v0 | pure-Go **source-pattern** heuristics over `ParsedUnit.Source` (+ `PyCweFacts` / needles) |
| Metadata | hand-authored `MetaCWE*` aligned with chunk `name` / `detection_notes` |
| File size | **1500 soft / 2000 hard** lines per Go file — split before growing |
| Fixtures (text) | `tests/fixtures/python/cwe/CWE-N-{vulnerable,safe}.txt` — same `.txt` contract as `tests/fixtures/python/bp/` |
| Unit tests | `internal/lang/python/detectors/cwe/*_test.go` hit/miss per ID |
| Integration | `tests/integration/python/cwe_matrix_test.go` (`make integration-python`) — **not** Go matrix |
| Product default | Go-only registry; Python opt-in via `languages = ["python"]` |
| Go CWE | **read-only** inspiration — never modify Go detectors for Python work |

**Do not** reopen shipped IDs as pending — see [batch-00-shipped.md](./batch-00-shipped.md).

---

## Testing contract (required for every implement batch)

Mirror the BP-PY layout under `tests/fixtures/python/bp/` and `tests/integration/python/`.  
**A rule is not done until all three layers below are green for that ID.**

### 1. Fixture text files (hit / miss corpus)

| Item | Requirement |
|------|-------------|
| Directory | `tests/fixtures/python/cwe/` (**not** under `python/bp/` or `python/` root) |
| Naming | **Exactly** `CWE-<N>-vulnerable.txt` and `CWE-<N>-safe.txt` (same stem pattern as `BP-PY-N-{vulnerable,safe}.txt`) |
| Format | Plain `.txt` fixtures only — see `tests/fixtures/README.md` |
| Header | `lang: python` + optional `file: CWE-<N>-….py` + `---` + Python body |
| Vulnerable | Must cause `RegisterRule("CWE-<N>")` to fire when scanned with `languages=["python"]` and `--only CWE-<N>` |
| Safe | Must **not** emit `CWE-<N>` (parameterized / safe API / confined path / etc.) |
| No committed `.py` | Only `.txt`; materialization is test-time |

**Example** (same shape as BP fixtures):

```text
# Fixture for CWE-90 (vulnerable)
lang: python
file: CWE-90-vulnerable.py
---
import ldap3
# ... sink that must fire CWE-90 ...
```

Shipped reference pairs (batch-00):

```text
tests/fixtures/python/cwe/CWE-22-{vulnerable,safe}.txt
tests/fixtures/python/cwe/CWE-78-{vulnerable,safe}.txt
tests/fixtures/python/cwe/CWE-79-{vulnerable,safe}.txt
tests/fixtures/python/cwe/CWE-89-{vulnerable,safe}.txt
tests/fixtures/python/cwe/CWE-502-{vulnerable,safe}.txt
```

BP analogue for comparison:

```text
tests/fixtures/python/bp/BP-PY-1-{vulnerable,safe}.txt
… (50 pairs) …
```

### 2. Individual (unit) tests — package-local

| Item | Requirement |
|------|-------------|
| Package | `internal/lang/python/detectors/cwe/` |
| Files | extend `rules_test.go` / `scan_test.go`, **or** domain `rules_<theme>_test.go` when approaching 1500-line soft cap |
| Per ID | At least one **hit** snippet and one **miss** snippet asserting `RuleID == "CWE-<N>"` presence/absence |
| Catalogue | Keep / extend registration want-list (all registered IDs after the batch) |
| Command | `go test ./internal/lang/python/detectors/cwe/ -count=1` |

### 3. Integration matrix tests — Python-only package

| Item | Requirement |
|------|-------------|
| Package | `tests/integration/python/` (**separate** from Go `tests/integration/`) |
| Test | `TestPythonCWEFixturesMatrix` in `cwe_matrix_test.go` |
| Discovery | `integration.DiscoverPythonCWECases()` — auto-picks every paired stem under `tests/fixtures/python/cwe/` |
| Paths | `integration.PythonCWEFixtureRel(case, vulnerable\|safe)` |
| Scan opts | `Languages = [LanguagePython]`, `Only = [CWE-N]`, typically `ProfileAll` |
| Registry | Python opt-in registry (not `DefaultRegistry` / Go-only) |
| Command | `go test ./tests/integration/python/ -count=1` **or** `make integration-python` |
| BP analogue | `TestPythonBPFixturesMatrix` + `DiscoverPythonBPCases` over `tests/fixtures/python/bp/` |

**Do not** put Python CWE pairs into `tests/fixtures/go/` or the Go CWE matrix (`tests/integration/cwe_matrix_test.go`).

### 4. Definition of done (per CWE ID in an implement batch)

- [ ] Detector + `RegisterRule` + meta + gates
- [ ] Unit hit/miss in `internal/lang/python/detectors/cwe/*_test.go`
- [ ] Fixture pair `tests/fixtures/python/cwe/CWE-N-{vulnerable,safe}.txt`
- [ ] Integration matrix picks up the new pair (no hard-coded allowlist required if discovery is glob-based — verify count increases)
- [ ] `make lint` + `make test` green (includes unit + usually integration via `./...` if configured; always run `make integration-python` before merge)
- [ ] Inventory / ledger checkboxes updated

### 5. What deferred IDs skip

IDs in [batch-deferred.md](./batch-deferred.md) need **no** fixtures, unit tests, or matrix rows until promoted into an implement batch.


---

## Waves (security-first, not numeric ID order)

| Wave | Batches | Focus |
|------|---------|-------|
| **P0** | 01–04 | Injection expand, dynamic code, mass-assign (915 High), secrets |
| **P1** | 05–08 | Crypto/RNG, SSRF/redirect, XXE/XML, path/fs |
| **P2** | 09–12 | Web/debug config, authz/session, info exposure, validation/ReDoS |
| **P3** | 13–14, 16 | Platform quality High, resources/upload, full tier-B expansion |
| **Deferred** | deferred | Tier-C remainder — honesty ledger only |

---

## Batch index

| Batch | File | #IDs | Wave | Theme | Target files |
|------:|------|-----:|------|-------|--------------|
| 00 | [batch-00-shipped.md](./batch-00-shipped.md) | 5 | — | shipped-priority | `rules.go (split to domain files over time)` |
| 01 | [batch-01-injection-expand.md](./batch-01-injection-expand.md) | 6 | P0 | injection-expand | `rules_injection.go` |
| 02 | [batch-02-code-dynamic.md](./batch-02-code-dynamic.md) | 5 | P0 | code-dynamic | `rules_code_dynamic.go` |
| 03 | [batch-03-mass-assign-deser.md](./batch-03-mass-assign-deser.md) | 3 | P0 | mass-assign-deser | `rules_mass_assign.go` |
| 04 | [batch-04-secrets-credentials.md](./batch-04-secrets-credentials.md) | 8 | P0 | secrets-credentials | `rules_secrets.go` |
| 05 | [batch-05-crypto-rng.md](./batch-05-crypto-rng.md) | 9 | P1 | crypto-rng | `rules_crypto.go` |
| 06 | [batch-06-ssrf-redirect.md](./batch-06-ssrf-redirect.md) | 6 | P1 | ssrf-redirect | `rules_ssrf.go` |
| 07 | [batch-07-xxe-xml.md](./batch-07-xxe-xml.md) | 3 | P1 | xxe-xml | `rules_xml.go` |
| 08 | [batch-08-path-fs-temp.md](./batch-08-path-fs-temp.md) | 8 | P1 | path-fs-temp | `rules_path_fs.go` |
| 09 | [batch-09-web-debug-config.md](./batch-09-web-debug-config.md) | 8 | P2 | web-debug-config | `rules_web_config.go` |
| 10 | [batch-10-authz-session.md](./batch-10-authz-session.md) | 8 | P2 | authz-session | `rules_auth.go` |
| 11 | [batch-11-info-exposure.md](./batch-11-info-exposure.md) | 8 | P2 | info-exposure | `rules_info_exposure.go` |
| 12 | [batch-12-validation-export.md](./batch-12-validation-export.md) | 8 | P2 | validation-export-redos | `rules_validation.go` |
| 13 | [batch-13-platform-quality.md](./batch-13-platform-quality.md) | 6 | P3 | platform-quality | `rules_platform.go` |
| 14 | [batch-14-resource-upload.md](./batch-14-resource-upload.md) | 8 | P3 | resource-upload | `rules_resource.go` |
| 16 | [batch-16-tier-b-expansion.md](./batch-16-tier-b-expansion.md) | 60 | P3 | tier-b-expansion | `rules_tier_b_*.go (split by theme when implementing)` |
| deferred | [batch-deferred.md](./batch-deferred.md) | 185 | defer | catalogue-deferred | `(none — no RegisterRule)` |

### ID ownership (no double-booking)

Every catalogue ID appears in **exactly one** batch (00, 01–16, or deferred). See `_inventory.json` `batches.*.ids`.

---

## PR policy

1. One PR per implement batch (`python(cwe): batch-NN …` / `Relates to #52` / `Relates to #51`).
2. Optional **infra PR first**: split shipped detectors out of monolith `rules.go` into domain files before P0 flood.
3. Before PR: `gofmt`, `make lint`, `make test`; record outcomes in batch Phase validation.
4. **Testing triad for each new ID** (required):
   - unit hit/miss in `internal/lang/python/detectors/cwe/*_test.go`
   - fixture pair under `tests/fixtures/python/cwe/CWE-N-{vulnerable,safe}.txt` (BP-style `.txt`)
   - green `TestPythonCWEFixturesMatrix` / `make integration-python`
5. Update `_inventory.json` `implemented` / `missing` when a batch lands.
6. Do **not** `Closes #52` until deferred inventory is honest in parent `python-heuristics-cwe.md`.
7. Batch-16 (tier-B) may be split into multiple PRs at implement time if still >15 IDs.

---

## Current package size (baseline)

From `_inventory.json` / `wc` at plan time:

| `common.go` | 530 |
| `facts.go` | 26 |
| `metadata.go` | 62 |
| `needles.go` | 37 |
| `rules.go` | 463 |
| `rules_test.go` | 324 |
| `scan.go` | 166 |
| `scan_test.go` | 86 |

| **Package total** | **1694** |

`rules.go` (~463) holds all five shipped rules — **split before large expansion**.

---

## Dependencies

| Depends on | Note |
|------------|------|
| Parent ledger | `plans/v0.0.2/heuristics/python-heuristics-cwe.md` |
| Epic rollup | `plans/v0.0.2/heuristics/python-heuristics.md` |
| Scaffold | `PyCweScan` + `RegisterRule` already on `main` |
| Catalogue | `ruleset/python/chunks/cwe-*.json` (344 keys) |
| Skill | `plans/skills/phase-wise-checklist/SKILLS.md` |
| Template | `plans/v0.0.2/heuristics/bp-plans/` |

---

## References

- Catalogue README: `ruleset/python/README.md`
- Detector package README: `internal/lang/python/detectors/cwe/README.md`
- Go pattern reference (do not modify): `internal/lang/go/detectors/cwe/`
