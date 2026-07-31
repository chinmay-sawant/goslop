# False-positive audit template

## Run metadata

```yaml
timestamp: YYYY-MM-DDTHH:MM:SSZ
repository: repository-name
repository_path: /absolute/path/to/repository
branch: branch-name
commit: git-revision
scan_target: /absolute/path/to/scan-target
chunk_path: ./scripts/chunks
function_context_path: ./scripts/findings/functions
```

## Scan evidence

- Build command: `<exact build command>`
- Scan command: `<exact scan command>`
- Findings: `<count>`
- Chunks reviewed: `./scripts/chunks/<chunk-file(s)>`
- Function contexts reviewed: `./scripts/findings/functions/<finding-id>.txt`

## Audit checklist

- [ ] Read every assigned chunk under `./scripts/chunks`.
- [ ] Read `./scripts/findings/functions/<finding-id>.txt` for every proposed false positive.
- [ ] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [ ] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [ ] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [ ] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [ ] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | `<count>` | `<ids>` |
| True positive | `<count>` | `<ids>` |
| Uncertain | `<count>` | `<ids>` |

## False positives

Create one subsection per finding. The source excerpt must be copied from the function-context file or the referenced source file and must be the smallest excerpt that proves the decision. Keep the explanation to one or two lines.

### [ ] Finding `<id>` — `<rule-id>`

- Function context: `./scripts/findings/functions/<id>.txt`
- Source: `/absolute/path/to/source:<line>:<column>`
- Checklist pattern: `<pattern from the checklist below>`

Source excerpt:

```
<smallest source excerpt proving the finding is a false positive>
```

Why this is a false positive: `<one or two lines explaining why this source does not satisfy the rule condition.>`

Checklist evidence: `<state the exact checked condition, for example: “the query is a bound ORM expression, not interpolated command text.”>`

Repeat this section for every false positive. Only group findings when they reference the exact same source construct; list every grouped ID explicitly.

## Uncertain findings

### [ ] Finding `<id>` — `<rule-id>`

- Function context: `./scripts/findings/functions/<id>.txt`
- Source: `/absolute/path/to/source:<line>:<column>`

Source excerpt:

```
<smallest relevant source excerpt>
```

Why it is uncertain: `<one or two lines describing the missing context, deployment assumption, or threat-model dependency.>`

## True positives

| Finding ID | Rule | Source | Short evidence |
| ---: | --- | --- | --- |
| `<id>` | `<rule-id>` | `<path:line>` | `<one-line reason the rule condition is present>` |

## Final evidence

- Delegated reviewers: `<names/tasks>`
- Chunk evidence: `./scripts/chunks`
- Function evidence: `./scripts/findings/functions`
- Validation: `git diff --check` — `<pass/fail>`
