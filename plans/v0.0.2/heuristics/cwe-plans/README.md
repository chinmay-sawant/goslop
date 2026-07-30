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
| Fixtures | `tests/fixtures/python/cwe/CWE-N-{vulnerable,safe}.txt` + matrix `tests/integration/python/cwe_matrix_test.go` |
| Product default | Go-only registry; Python opt-in via `languages = ["python"]` |
| Go CWE | **read-only** inspiration — never modify Go detectors for Python work |

**Do not** reopen shipped IDs as pending — see [batch-00-shipped.md](./batch-00-shipped.md).

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
4. Add hit/miss unit tests + fixture pair + matrix discovery for each new ID.
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
