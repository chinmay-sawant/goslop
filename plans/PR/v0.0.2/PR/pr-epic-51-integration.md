## Summary

Integration branch for epic **#51** (Python heuristic detectors). Merges CWE priority heuristics (#52) and BP-PY priority heuristics (#53). PERF (#54) stays deferred (no catalogue). Prefer reviewing/merging **this** PR; child PRs are superseded.

---

## Motivation / context

- Plans: `plans/v0.0.2/heuristics/python-heuristics.md` and sibling CWE/BP/PERF ledgers
- Foundation: epic #39 / PR #50 (plugin + catalogues)
- Issues: see **Related issues**

---

## Child PRs (superseded by this integration)

| Stream | Issue | PR | Branch |
|--------|------:|----|--------|
| CWE priority | #52 | #56 | `feat/python-cwe-heuristics` |
| BP priority | #53 | #57 | `feat/python-bp-heuristics` |
| PERF | #54 | — | deferred (no PR) |

---

## Changes

### CWE (#52)

- `internal/lang/python/detectors/cwe/` — `PyCweScan`, RegisterRule, pure-Go source patterns
- Rules: **CWE-22, 78, 79, 89, 502**
- Fixtures under `tests/fixtures/python/cwe/`
- Unit hit/miss tables

### BP (#53)

- `internal/lang/python/detectors/bad_practices/` — priority subset
- Rules: **BP-PY-1,2,4,6,7** (core) + **8–13** (security) + **16,17,21** (Flask/Django sample)
- Inline unit hit/miss tests

### Integration

- `detectors/all.go` registers **both** CWE and BP scans
- Plugin docs/tests assert combined catalogue
- `list-rules` with `languages=["python"]` shows CWE + BP-PY; no Go PERF

### Still deferred

- Remaining CWE catalogue (beyond 5)
- Remaining BP-PY (batches C/D/E rest)
- All PERF (#54)

---

## Impact

| Area | Impact |
|------|--------|
| **Behavior** | Opt-in Python scans can emit CWE + BP-PY findings |
| **Default** | Go-only unchanged |
| **API / CLI** | No new flags; uses existing `languages` |

---

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| None for default scans | Unset languages → Go only |

---

## Test plan

- [x] `make lint` — green on integration
- [x] `make test` — green (incl. `internal/lang/python/...`)
- [x] `CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop`

### Commands

```sh
make lint
make test
```

---

## Related issues

- Closes #52
- Closes #53
- Relates to #51 (epic remains open while PERF #54 and expansions continue)
- Relates to #54 (deferred)

---

## Follow-ups

- Expand CWE beyond priority 5
- Remaining BP-PY rules
- PERF catalogue seed + detectors (#54)

---

## Reviewer checklist

- [ ] Default remains Go-only
- [ ] Python opt-in surfaces both CWE and BP-PY
- [ ] No PERF claim without catalogue
- [ ] Child PRs superseded by this integration
