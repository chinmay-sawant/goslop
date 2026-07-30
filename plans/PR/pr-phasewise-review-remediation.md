## Summary

Completes the phase-wise scan-contract, configuration, cache, export, and Go-practices remediation for Goslop. The branch makes failure and ownership boundaries explicit, preserves existing CLI/output behavior, and closes every concrete item in the post-remediation Go review.

---

## Motivation / context

- Plans: `plans/v0.0.1/reviews/ponytail-code-and-architecture-review.md`, `plans/v0.0.1/reviews/ponytail-post-remediation-review.md`, `plans/v0.0.1/reviews/architecture-post-remediation-review.md`, and `plans/v0.0.1/reviews/go-code-style-and-design-review.md`
- Issues: none supplied; this PR is driven by the phase-wise review ledgers.

---

## Changes

### Scan lifecycle, configuration, and parser-quality contracts

- Create detector instances per scan, validate registry catalogue/factory parity, and retain correct detector finalization on cache hits.
- Centralize resolved CLI/TOML scan policy in `scanPlan`, remove unsupported configuration promises, and add observable configuration-effect coverage.
- Preserve source-only Go analysis for malformed input while exposing parse quality and a non-fatal diagnostic through cold and cached paths.
- Make BP project/package facts scan-owned and single-flight, preventing cross-root state leakage and duplicate cold reads.

### Output, cache, and failure handling

- Restrict export cleanup to exporter-owned files, return cleanup/write errors, and guarantee temporary matrix cleanup on failures.
- Make context/chunk directory collision checks read-only so rejected configurations do not create directories.
- Propagate cache tool-version invalidation cleanup failures instead of silently retaining stale entry files.
- Consolidate per-export rendering dependencies into private `renderState`; dual-surface output is regression-tested against each independent output mode.

### Go style, quality, and review artifacts

- Decode TOML directly from bytes and make the intentional partial-AST parser policy clear at the call site.
- Apply the lint-all cleanup and add deterministic export/cache failure regressions.
- Add separate Ponytail, architecture, and Go-style/design phase-wise Markdown and HTML review reports; the Go review is now fully checked off and re-rated at 9.0/10.

---

## Impact

| Area | Impact |
|------|--------|
| **Performance** | No detector algorithm change. TOML parsing avoids an unnecessary byte-to-string conversion. |
| **Memory** | No material increase; export caches remain scoped to one export call. |
| **Behavior / correctness** | Invalid export configuration no longer creates directories; cache invalidation and output cleanup failures are observable. |
| **API / CLI** | No intended public CLI or API break. Unsupported inert configuration fields were removed from the documented/schema contract. |
| **Dependencies** | None. |
| **Binary size / build time** | No material expected change. |

---

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| Removed inert `languages` and `typed.enabled` configuration promises | Remove those unsupported settings from `goslop.toml`; active Go scanning and supported configuration behavior are unchanged. |

---

## Test plan

- [x] `make test`
- [x] `make lint-all`
- [x] `go test -race ./...`
- [x] `git diff --check`
- [ ] `CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop` (not run separately; no CGO-related change)
- [ ] `make run` wall-time comparison (not rerun; detector surface did not change in this closure batch)
- [ ] `make reference-metrics` / gopdfsuit hard metrics (not applicable; no detector-rule surface changed)

### Commands

```sh
make lint-all
make test
go test -race ./...
git diff --check
```

---

## Screenshots / sample output

No user-interface change. Focused regression coverage proves export collision rejection leaves no directory behind, dual output matches independent context/chunk output, and cache cleanup errors propagate.

---

## Related issues

- None supplied or linked; this is review-led remediation work.

---

## PR metadata checklist (author)

- [x] Self-assigned (`--assignee @me`)
- [x] Labels applied (`enhancement`, `documentation`)
- [x] Related issues reviewed; no matching ticket was supplied
- [x] Filled body committed under `plans/PR/pr-phasewise-review-remediation.md`

---

## Follow-ups (out of scope)

- Freeze or synchronize late default-plugin registration only if runtime plugin mutation becomes supported.
- Define `ScanContext` immutability only if external callers need concurrent reconfiguration.
- Revisit multi-root cache identity only with an observed reuse/order issue and an acceptance test.

---

## Reviewer checklist

- [ ] Behavior matches summary and test plan
- [ ] No unrelated changes in diff
- [ ] Public API / CLI changes documented
- [ ] New rules have fixture coverage when applicable
- [ ] PR has assignee and labels
- [ ] Related issues use correct Closes/Relates keywords when applicable
- [ ] No secrets or generated artifacts committed
