# CodeHound Go Port — Plans

Canonical execution ledger for porting the Rust CodeHound application to Go.

| Document | Role |
|----------|------|
| [`port-phasewise-checklist.md`](./port-phasewise-checklist.md) | **Canonical** phase-wise checklist (live ledger) |
| [`parity-matrix.md`](./parity-matrix.md) | Rust → Go package map + rule-pack parity tracking |
| [`architecture-go.md`](./architecture-go.md) | Target Go layout and design decisions |

## Source of truth

- **Rust product:** `/home/chinmay/ChinmayPersonalProjects/codehound`
- **Go port:** `/home/chinmay/ChinmayPersonalProjects/codehound-go`
- **Fixtures / heuristics text:** copied **as-is** under `tests/fixtures/` and `ruleset/`
- **Detector logic:** rewritten in Go (not a mechanical transpile)

## Status legend

- `[ ]` not started / not proven
- `[x]` implemented and validated with evidence
- `[~]` deferred/partial — reason + next gate required
