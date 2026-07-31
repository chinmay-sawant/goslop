## Summary

This follow-up closes the remaining Python false-positive audit gaps. CWE-1124 now measures executable control-flow depth without counting class or function declarations, and CWE-924 recognizes an explicit module-level API-key authentication boundary.

## Motivation / context

- Plans: `plans/fp-review-reduce/python/false-positive-audit-2026-07-31.md`
- Implementation evidence: `plans/fp-review-reduce/python/review-and-reduce-implementation-2026-07-31.md`
- Related work: `Relates to #71`

## Changes

### Detector guardrails

- Refined CWE-1124 to exclude lexical class/function scopes and avoid counting exception-handler clauses as extra control-flow depth.
- Refined CWE-924 to recognize module-level header authentication for authenticated webhook relays.
- Updated the audit classification to 29 false positives, 20 true positives, and no remaining uncertain findings.

### Fixture coverage

- Added dedicated text-fixture pairs for the CWE-1124 declaration-scope case.
- Added authenticated-route safe/vulnerable CWE-924 text fixtures.
- Wired both variants into per-heuristic and Python integration fixture tests.

## Impact

| Area | Impact |
|------|--------|
| **Performance** | No algorithmic or scan-pipeline changes; detector scans remain linear over source text. |
| **Memory** | No material change. |
| **Behavior / correctness** | Removes declaration-scope CWE-1124 noise and API-key-authenticated CWE-924 noise while preserving vulnerable fixture coverage. |
| **API / CLI** | None. |
| **Dependencies** | None. |
| **Binary size / build time** | No material change. |

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| None | - |

## Test plan

- [x] `make test`
- [x] `make lint-all`
- [x] `go build -o bin/goslop ./cmd/goslop`
- [x] Python detector and integration fixture suites
- [x] Full no-cache corpus scan: 20 findings after 29 confirmed false positives were removed

### Commands

```sh
GOCACHE=/tmp/goslop-review-reduce-gocache go test ./internal/lang/python/detectors/cwe/... -count=1
GOCACHE=/tmp/goslop-review-reduce-gocache go test ./tests/integration/python -count=1
GOCACHE=/tmp/goslop-review-reduce-gocache GOLANGCI_LINT_CACHE=/tmp/goslop-review-reduce-golangci make lint-all
GOCACHE=/tmp/goslop-review-reduce-gocache make test
GOCACHE=/tmp/goslop-review-reduce-gocache go build -o bin/goslop ./cmd/goslop
```

## Screenshots / sample output

```text
scanned 52 files (2386 lines)
20 findings
severity: 7 high, 0 info, 0 low, 13 medium
```

## Related issues

- Relates to #71

## PR metadata checklist (author)

- [ ] Self-assigned (`--assignee @me`)
- [ ] Labels applied
- [x] Related work filled (`Relates to #71`)
- [x] Filled body committed under `plans/PR/pr-python-fp-audit-finalize.md`

## Follow-ups (out of scope)

- None.

## Reviewer checklist

- [ ] Behavior matches summary and test plan
- [ ] No unrelated changes in diff
- [ ] Public API / CLI changes documented
- [x] New rules have fixture coverage when applicable
- [ ] PR has assignee and labels
- [x] Related work uses the `Relates to` keyword
- [ ] No secrets or generated artifacts committed
