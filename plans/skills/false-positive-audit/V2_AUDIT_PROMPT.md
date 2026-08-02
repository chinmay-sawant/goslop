# Post-fix v2 audit instructions (read first)

You are auditing ONE repo's findings from the **latest binary** (bin/goslop rebuilt ~2026-08-02 17:56). The previous audit reports (Mode A/B appends from the b5b8fde binary) already classify most of these sources. Your job: classify every fresh finding, reuse prior classifications by `Source:` match, and produce a compact fix checklist.

## Evidence

- Fresh chunks (latest binary): `scripts/<REPO>/chunks/Chunk_*.txt` — READ ALL
- Fresh contexts: `scripts/<REPO>/findings/functions/<id>.txt`
- Prior classifications: `plans/fp-validations/reports/<REPO>.md` — old audit FP/TP lists (by Source) in the original audit + Mode A/B appends
- Repo source: `real-repos/<REPO>/...` (follow `Source:` paths)

## Workflow

1. Read every fresh chunk. For each fresh finding:
   - If its `Source:` (file:line:col) matches an audited TP → True positive (TP).
   - If it matches an audited FP → False positive (FP) — reuse the prior reason.
   - If unmatched → classify fresh against the rule condition (`./bin/goslop -explain <RULE>` when unsure); read its function context + enclosing source.
2. Count: fresh TP / fresh FP / fresh Uncertain. Note any NEW findings (not in old audit at all).
3. Group the FPs into **patterns** (rule id + exact trigger shape + count + first-seen sources). These patterns are the fix checklist for the detector fix phase — be precise about the *condition* that should distinguish safe from vulnerable (e.g. "shebang + library module + no __main__ guard", "raise RuntimeError from e where caught type is Exception", "examples/ path + no guard").
4. APPEND to `plans/fp-validations/reports/<REPO>.md`:

```
## Post-fix v2 audit (latest binary)

Run metadata (build: make build ~17:56), scan evidence, classification summary (fresh counts), then:
## Fix checklist (FP patterns)
| Pattern # | Rule | Trigger shape | Count | Example sources |
one row per grouped FP pattern (keep the shape description tight),
then ## New findings (any fresh finding with no prior classification; classify each).
```

Keep it compact — this is a working checklist, not a full template report. Do not edit prior sections. Run `git diff --check` (goslop root) and report pass/fail.

## Return value

Repo, fresh findings count, fresh TP/FP/U, number of grouped FP patterns (with rule ids), number of new findings. Under 12 lines.
