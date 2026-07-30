## Summary

Fix an infinite loop in BP-28/BP-29 text scanners that hung `make run` when scanning detector sources (and any file containing an unbalanced `"interface {"` string). Also clear `make lint-all` findings in export/config tests so CI-style lint can pass.

---

## Motivation / context

- Self-scan / dogfood of `internal/lang/go/detectors` (or `SCAN_PATH` pointed at goslop) never finished: BP-29 burned 99%+ CPU forever on `rules_api.go`.
- Root cause: brace matching started on a string-literal needle and never returned to depth 0, then `start = end` with `end == abs` looped forever.
- Fixture tests never caught this because BP-28/29 fixtures use real, balanced interfaces—not self-referential string needles.
- Plans: hang diagnosis from product `make run` debugging; no dedicated plan doc.

---

## Changes

### BP-28 / BP-29 hang fix

- Extract `matchBraceBlock` and always advance past the needle when braces do not balance.
- Unit regression: `TestBP28BP29UnbalancedBraceNeedleDoesNotHang` (2s wall).
- Integration dogfood: `TestDogfoodDetectorSourcesNoHang` runs profile-all on `internal/lang/go/detectors` with a 10s wall.

### lint-all

- Drop named returns in `enclosingFunctionLines` (`nonamedreturns`).
- Guard `strings.Index` before slicing in export tests (`gocritic offBy1`).
- Rename shadowed `err` in config test (`govet shadow`).

---

## Impact

| Area | Impact |
|------|--------|
| **Performance** | Self-scan / dogfood completes in ~200ms instead of hanging. |
| **Memory** | None |
| **Behavior / correctness** | Unbalanced string-literal `"interface {"` no longer hangs; may skip false positive interface parses (correct). |
| **API / CLI** | None |
| **Dependencies** | None |
| **Binary size / build time** | Negligible |

---

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| None | - |

---

## Test plan

- [x] `make lint-all`
- [x] `go test ./internal/lang/go/detectors/bad_practices/ -run TestBP28BP29`
- [x] `go test ./tests/integration/ -run TestDogfoodDetectorSourcesNoHang`
- [x] `go test ./internal/export/ ./internal/config/`
- [ ] `make test` (full suite recommended before merge)
- [ ] `make run` on gopdfsuit (oracle path)

### Commands

```sh
make lint-all
go test ./internal/lang/go/detectors/bad_practices/ ./tests/integration/ ./internal/export/ ./internal/config/ -count=1 -timeout 60s
make run
```

---

## Screenshots / sample output

```
# before: BP-29 alone on rules_api.go → timeout 30s+
# after:
scanned 1 files (2091 lines) in 42.0ms
dogfood detectors: files=58 findings=512 (wall limit 10s)  PASS
make lint-all → ok
```

---

## Related issues

- Relates to product `make run` hang when scanning goslop / detector sources
- No open GitHub issue ID

---

## PR metadata checklist (author)

- [x] Self-assigned (`--assignee @me`)
- [x] Labels applied
- [x] Related issues filled (no ticket ID)
- [x] Filled body under `plans/PR/pr-bp29-infinite-loop.md`

---

## Follow-ups (out of scope)

- Broader dogfood of full repo (not only detectors)
- AST-based interface detection to ignore string/comment needles entirely

---

## Reviewer checklist

- [ ] Behavior matches summary and test plan
- [ ] No unrelated changes in diff
- [ ] Public API / CLI changes documented
- [ ] New rules have fixture coverage when applicable
- [ ] PR has assignee and labels
- [ ] Related issues use correct Closes/Relates keywords
- [ ] No secrets or generated artifacts committed
