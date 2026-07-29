# Phase 7–12 parallel workstreams

## Dependency graph (from checklist)

```
main (0–6 done)
  ├─ Phase 7 CWE structural  ──► Phase 9 taint
  ├─ Phase 8 BP
  ├─ Phase 10 cache / baseline / ignore
  └─ Phase 12 scaffolding (CI, harness)  [partial; §12.4 blocked]
         │
  6+7+8 ─► Phase 11 packs / maturity
  6..11 ─► Phase 12 full parity / §12.4 oracle
```

## Parallel wave 1 (unlocked)

| Phase | Issue | Branch | Unlocked? |
|------:|------:|--------|-----------|
| 7 | #3 | `feat/phase-7-cwe` | yes (deps 0–3) |
| 8 | #4 | `feat/phase-8-bp` | yes (deps 0–3) |
| 9 | #5 | `feat/phase-9-taint` | yes-ish (seed CWE + fixtures; full structural CWE optional) |
| 10 | #6 | `feat/phase-10-cache` | yes (deps 3–5) |
| 12a | #8 | `feat/phase-12-ci` | partial (CI/tests/docs only; not §12.4) |

## Deferred

| Phase | Issue | Why deferred |
|------:|------:|--------------|
| 11 | #7 | Needs integrated 6+7+8 detector surface for pack fidelity |
| 12.4 | #8 | Needs 6–11 + export/cache full surface |

## Integration rule

- Child PRs stay **open / unmerged**
- A later **integration PR** merges heads and runs full validation
