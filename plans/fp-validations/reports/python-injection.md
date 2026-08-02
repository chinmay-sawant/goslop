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

## Post-fix over-suppression audit (2026-08-02)

### Run metadata

```yaml
timestamp: 2026-08-02T11:08:51Z
repository: python-injection
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/python-injection
branch: prod
commit: 4f36e8faedd41272e91dbcb493c87abfc93c3b66 (unchanged from pre-fix audit)
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/python-injection
chunk_path: scripts/python-injection/chunks
function_context_path: scripts/python-injection/findings/functions
scanner: ./bin/goslop rebuilt from fix commit b5b8fde (2026-08-02 16:29)
```

### Scan evidence

- Build command: `make build` (`go build -o bin/goslop ./cmd/goslop`)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/python-injection/chunks -context-dir scripts/python-injection/findings/functions real-repos/python-injection`
- Findings: `1` (pre-fix audit had 2)
- Chunks reviewed: `scripts/python-injection/chunks/Chunk_1_1.txt`
- Function contexts reviewed: `scripts/python-injection/findings/functions/1.txt`

### Audit checklist

- [x] Collected the full audited TP list from the pre-fix audit (True positives table).
- [x] Collected the sources present in the fresh scan from `scripts/python-injection/chunks`.
- [x] Compared each audited TP `Source:` against fresh-scan sources.
- [x] Read the current source file for every audited TP absent from the fresh scan and confirmed whether the construct still exists.
- [x] Verified the rule condition in `internal/lang/python/detectors/bad_practices/rules_testing.go:32` (`detectBPPY41`) and with a targeted re-scan of the file.
- [x] Ran `git diff --check` after updating this report.

### Over-suppression table

| Old finding ID | Rule | Source | One-line reason (from old audit) | Current status |
| --- | --- | --- | --- | --- |
| 1 | BP-PY-41 | `.../tests/loaders/test_profile_loader.py:111` | Body calls `mod()`, `uuid4().hex`, `ProfileLoader(...)` and enters `loader.load(...)` but contains no assert — a placeholder "do nothing" test that passes even if loading breaks. | still flagged (present in fresh scan) |
| 2 | BP-PY-41 | `.../tests/test_inject.py:56` | Body only calls `my_function()` (an `@inject`-decorated no-parameter function); no assert — if `@inject` became a no-op passthrough the test would still pass. | **suppressed-but-present** |

Counts: over-suppressed TPs = 1 (old finding 2); fixed-removed = 0.

### [x] Finding 2 (old audit ID) — `test_inject_with_no_parameter` at tests/test_inject.py:56 — BP-PY-41 (over-suppressed)

- Current source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/python-injection/tests/test_inject.py:56`
- Not present in fresh scan; construct still present in source (verified `git log` — repo commit unchanged since pre-fix audit).

Source excerpt (`tests/test_inject.py:56-60`):

```python
    def test_inject_with_no_parameter(self):
        @inject
        def my_function(): ...

        my_function()
```

Why this satisfies the rule condition: `detectBPPY41` scans the body of `def test_inject_with_no_parameter` for a side-effect call and no assertion. The body contains `my_function()` — a statement with `(` and `)` that does not start with any of the `looksLikeSideEffectCall` exclusions (control-flow/def/with/import prefixes) — so `hasCall` is set; the body has no `assert`/`pytest.raises`/`self.assert*`/assertion-helper call, so `hasAssert` stays false; `hasCall && !hasAssert` fires. The `def my_function(): ...` line is skipped by the `def ` exclusion, and `@inject` matches no call pattern.

Verification: a targeted post-fix re-scan `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --no-cache real-repos/python-injection/tests/test_inject.py` returns 0 findings (`no slop detected`), while the pre-fix scan (same repo commit 4f36e8f) reported finding 2 at line 56 — the fix commit b5b8fde suppressed this true positive (its BP-PY-41 changes added assertion-helper crediting, `pytest.fail(`, regression/benchmark fixtures and a triple-quote body scan, none of which this body uses). Recommend review: the guardrail suppressed a genuine placeholder test.

### Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/python-injection/chunks`
- Function evidence: `scripts/python-injection/findings/functions`
- Validation: `git diff --check` — pass
