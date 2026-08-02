# False-positive audit: python-injection

## Run metadata

```yaml
timestamp: 2026-08-02T07:13:28Z
repository: python-injection
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/python-injection
branch: prod
commit: 4f36e8faedd41272e91dbcb493c87abfc93c3b66
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/python-injection
chunk_path: scripts/python-injection/chunks
function_context_path: scripts/python-injection/findings/functions
```

## Scan evidence

- Build command: `make build` (`go build -o bin/goslop ./cmd/goslop`)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/python-injection/chunks -context-dir scripts/python-injection/findings/functions real-repos/python-injection`
- Findings: `2`
- Chunks reviewed: `scripts/python-injection/chunks/Chunk_1_2.txt`
- Function contexts reviewed: `scripts/python-injection/findings/functions/1.txt`, `scripts/python-injection/findings/functions/2.txt`

## Audit checklist

- [x] Read every assigned chunk under `scripts/python-injection/chunks`.
- [x] Read `scripts/python-injection/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 0 | — |
| True positive | 2 | 1, 2 |
| Uncertain | 0 | — |

## True positives

### BP-PY-41 — pytest assert With Side Effects Only

Rule condition (from `detectBPPY41` in `internal/lang/python/detectors/bad_practices/rules_testing.go`): for each `def test_*` line, scan the body; flag when the body contains at least one side-effect call (a statement containing `(` and `)`, excluding control-flow/def/with/import prefixes) and no assertion (`assert`, `pytest.raises`, `self.assert*`, or a call to an assertion helper). Both findings satisfy `hasCall && !hasAssert`.

| Finding id | Source | One-line reason |
| --- | --- | --- |
| 1 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/python-injection/tests/loaders/test_profile_loader.py:111` | Body calls `mod()`, `uuid4().hex`, `ProfileLoader(...)` and enters `loader.load(...)` but contains no assert — a placeholder "do nothing" test that passes even if loading breaks. |
| 2 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/python-injection/tests/test_inject.py:56` | Body only calls `my_function()` (an `@inject`-decorated no-parameter function); no assert — if `@inject` became a no-op passthrough the test would still pass. |

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/python-injection/chunks`
- Function evidence: `scripts/python-injection/findings/functions`
- Validation: `git diff --check` — pass
