# Phase 7-12 parallel workstreams

## Dependency graph (from checklist)

```
main (0-6 done)
  ├─ Phase 7 CWE structural  ──► Phase 9 taint
  ├─ Phase 8 BP
  ├─ Phase 10 cache / baseline / ignore
  └─ Phase 12 scaffolding (CI, harness)  [partial; §12.4 blocked]
         │
  6+7+8 ─► Phase 11 packs / maturity
  6..11 ─► Phase 12 full parity / §12.4 oracle
```

## Parallel wave 1 (unlocked) - PRs open, **not merged**

| Phase | Issue | Branch | PR | Status |
|------:|------:|--------|-----|--------|
| 7 | #3 | `feat/phase-7-cwe` | [#14](https://github.com/chinmay-sawant/goslop/pull/14) | OPEN - 175/175 CWE |
| 8 | #4 | `feat/phase-8-bp` | [#13](https://github.com/chinmay-sawant/goslop/pull/13) | OPEN - ~127/135 BP |
| 9 | #5 | `feat/phase-9-taint` | [#12](https://github.com/chinmay-sawant/goslop/pull/12) | OPEN - CWE-22/78/79/89 taint |
| 10 | #6 | `feat/phase-10-cache` | [#11](https://github.com/chinmay-sawant/goslop/pull/11) | OPEN - cache/baseline/ignore |
| 12a | #8 | `feat/phase-12-ci` | [#10](https://github.com/chinmay-sawant/goslop/pull/10) | OPEN - CI + harness (not §12.4) |

## Deferred

| Phase | Issue | Why deferred |
|------:|------:|--------------|
| 11 | #7 | Needs integrated 6+7+8 detector surface for pack fidelity |
| 12.4 | #8 | Needs 6-11 + export/cache full surface; PR #10 is scaffolding only |

## Integration rule

- Child PRs stay **open / unmerged**
- A later **integration PR** merges heads (`feat/phase-7-cwe` … `feat/phase-12-ci`) and runs full validation
- Preferred merge order into integration: **12a CI → 10 cache → 7 CWE → 8 BP → 9 taint → 11 packs** (then §12.4)
