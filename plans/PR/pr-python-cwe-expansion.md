## Summary

Expands Python CWE coverage from the initial five priority rules to 159 source-only heuristics. Every implemented rule is covered by the same paired text fixtures in both package-level and Python integration tests.

---

## Motivation / context

- Plans: `plans/v0.0.2/heuristics/cwe-plans/`
- Parent ledger: `plans/v0.0.2/heuristics/python-heuristics-cwe.md`
- This is one explicitly approved integration PR for batches 01–14 and 16.
- Issues: see **Related issues**

---

## Changes

### Python CWE detectors

- Add 154 approved high-signal source-pattern heuristics across injection, dynamic code, credentials, crypto, SSRF, XML, filesystem, configuration, authorization, information exposure, validation, platform, resource, and tier-B domains.
- Register metadata for every rule and keep all detection within the existing `PyCweScan` source-only detector surface.
- Harden several patterns against cross-function and mismatched-resource false positives.

### Fixture-backed verification and planning

- Add 154 canonical safe/vulnerable Python text-fixture pairs; the package now has 159 pairs for 159 registered rules.
- Replace inline Python source in per-rule Go tests with shared paired-fixture assertions, so unit and integration tests execute the same corpus.
- Reconcile the batch inventory and parent execution ledger: 159 implemented, 0 approved implementation IDs pending, and 185 catalogue IDs explicitly deferred.

---

## Impact

| Area | Impact |
|------|--------|
| **Performance** | No measured regression; the Go-only reference corpus remains at its 915-finding hard baseline. |
| **Memory** | No material expected change beyond source-pattern scanning for opted-in Python files. |
| **Behavior / correctness** | Adds Python-only CWE findings when `languages = ["python"]` is enabled; detector coverage is substantially expanded. |
| **API / CLI** | No new flags or schema changes. |
| **Dependencies** | None. |
| **Binary size / build time** | Pure-Go build remains valid. |

---

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| None | - |

---

## Test plan

- [x] `make test`
- [x] `make lint` / `go vet`
- [x] `CGO_ENABLED=0 go build -o /tmp/goslop-cwe-pr ./cmd/goslop`
- [x] `make reference-metrics REFERENCE_PATH=/home/chinmay/ChinmayPersonalProjects/gopdfsuit`
- [x] `make integration-python`
- [x] `git diff --check`

### Commands

```sh
make lint
make test
make integration-python
CGO_ENABLED=0 go build -o /tmp/goslop-cwe-pr ./cmd/goslop
make reference-metrics REFERENCE_PATH=/home/chinmay/ChinmayPersonalProjects/gopdfsuit
```

---

## Screenshots / sample output

```text
reference-metrics: 915 findings
severity: 10 high, 197 info, 312 low, 396 medium
exports: 915 context files and 37 chunk files
```

---

## Related issues

- Relates to #52
- Relates to #51

---

## PR metadata checklist (author)

- [x] Self-assigned (`--assignee @me`)
- [x] Labels applied
- [x] Related issues filled with real ticket IDs
- [x] Filled body committed under `plans/PR/pr-python-cwe-expansion.md`

---

## Follow-ups (out of scope)

- The remaining 185 Tier-C catalogue IDs stay deferred because they do not have a trustworthy pure-Go source-pattern heuristic in v0.

---

## Reviewer checklist

- [x] Behavior matches summary and test plan
- [x] No unrelated changes in diff
- [x] Public API / CLI changes documented
- [x] New rules have fixture coverage
- [x] PR has assignee and labels
- [x] Related issues use correct Closes/Relates keywords
- [x] No secrets or generated artifacts committed
