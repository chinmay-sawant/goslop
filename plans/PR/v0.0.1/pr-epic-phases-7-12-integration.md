## Summary

Integration branch that merges the parallel Phase 7-12 workstreams into one tree for combined validation. **Do not merge child PRs separately** - prefer this integration PR once green.

---

## Motivation / context

- Plans: [`plans/phase-parallel/README.md`](../phase-parallel/README.md), [`plans/port-phasewise-checklist.md`](../port-phasewise-checklist.md)
- Child workstreams were developed in isolated worktrees; this PR proves they compose.
- Issues: see **Related issues**

---

## Child PRs (superseded by this integration)

| Phase | PR | Branch | Role |
|------:|----|--------|------|
| 12a CI | #10 | `feat/phase-12-ci` | GitHub Actions, integration harness, README |
| 10 cache | #11 | `feat/phase-10-cache` | Incremental cache, baseline, ignores |
| 9 taint | #12 | `feat/phase-9-taint` | Taint graph CWE-22/78/79/89 |
| 8 BP | #13 | `feat/phase-8-bp` | Bad-practice detectors (~127 rules) |
| 7 CWE | #14 | `feat/phase-7-cwe` | Structural CWE catalogue (175/175) |

Direct links:

- https://github.com/chinmay-sawant/goslop/pull/10
- https://github.com/chinmay-sawant/goslop/pull/11
- https://github.com/chinmay-sawant/goslop/pull/12
- https://github.com/chinmay-sawant/goslop/pull/13
- https://github.com/chinmay-sawant/goslop/pull/14

Merge order used: **12a → 10 → 7 → 8 → 9** (then conflict resolution for CLI/cache/taint + detector wiring).

---

## Changes

### Integrated product surface

- **PERF** (already on `main`): 239/239
- **CWE structural**: full registry via `GoCweScan` + taint-lite seeds
- **BP**: `GoBadPracticeScan` (~127 rules); recommended pack keeps BP off
- **Taint**: experimental graph detector gated by `--taint` / security profile
- **Cache / baseline / ignore**: `.goslop-cache`, baseline file, `goslop-ignore`, walk ignores
- **CI**: `.github/workflows/ci.yml` + `tests/integration` seed harness

### Integration-specific conflict resolutions

- CLI `Options` / flags: **both** cache/baseline **and** taint flags
- `detectors.All()`: `GoCweScan` + `taint.NewDetector` + PERF + BP
- Skip taint-lite CWE-22/78/79/89 when `TaintEnabled` (no double findings)
- Dropped duplicate `cwe/metadata.go` (use `metadata_generated.go`)

### Not in this PR (follow-ups)

- **Phase 11** packs/maturity full fidelity (issue #7)
- **§12.4** hard expected baseline (915 findings / exports) - still blocked
- Export-context / export-chunks full product path

---

## Impact

| Area | Impact |
|------|--------|
| **Behavior** | Full multi-pack detector surface; taint opt-in; cache on by default |
| **CLI** | New flags from #11 and #12 combined |
| **CI** | CGO build + test workflow from #10 |
| **Risk** | Heuristic detectors may FN/FP vs Rust; needs corpus review before ship |

---

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| None for Go port consumers | Greenfield; recommended pack still BP-off |

---

## Test plan

- [x] `CGO_ENABLED=1 go test ./...` on integration branch
- [x] Child packages: cwe, taint, bad_practices, perf, engine/cache, tests/integration
- [ ] Manual: `go build -o bin/goslop ./cmd/goslop && ./bin/goslop --list-rules | wc -l`
- [ ] Manual: `--taint` on taint fixtures; cache hit/miss smoke
- [ ] Optional: close/supersede child PRs **without merging** them after this lands

### Commands

```sh
CGO_ENABLED=1 go test ./...
go build -o bin/goslop ./cmd/goslop
./bin/goslop --list-rules | head
./bin/goslop --profile all --taint --no-cache .
```

---

## Related issues

- Relates to #3 (via #14)
- Relates to #4 (via #13)
- Relates to #5 (via #12)
- Relates to #6 (via #11)
- Relates to #8 (via #10; §12.4 still open)
- Relates to #7 (Phase 11 still deferred)

---

## PR metadata checklist (author)

- [x] Self-assigned (`--assignee @me`)
- [x] Labels applied (`enhancement`, `documentation`)
- [x] Related issues + child PR links filled
- [x] Body under `plans/PR/pr-epic-phases-7-12-integration.md`

---

## Follow-ups (out of scope)

- Phase 11 recommended/security/all pack fidelity + maturity tags (#7)
- §12.4 parity baseline on reference corpus (#8)
- Supersede child PRs #10-#14 without double-merge once this is accepted
- Residual FN/FP polish vs Rust

---

## Reviewer checklist

- [ ] Combined tree builds and tests green
- [ ] Child PR scope is fully present (no silent drops)
- [ ] Taint + structural CWE do not double-fire when `--taint`
- [ ] Cache flags coexist with taint flags
- [ ] **Do not merge child PRs #10-#14 into main separately**

---

## Release notes (if user-facing)

Integration of CWE, BP, taint, cache/baseline, and CI scaffolding for the goslop Go port (pre-§12.4).
