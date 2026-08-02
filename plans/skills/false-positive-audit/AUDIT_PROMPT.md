# Shared audit instructions (read first)

You are auditing goslop scan findings for a single repository. Follow the template in `plans/skills/false-positive-audit/SKILLS.md` exactly and write the report to `plans/fp-validations/reports/<REPO_NAME>.md`.

## Evidence locations (per repo)

- Chunks (every finding, with source excerpt): `real-repos/<REPO_NAME>/scripts/chunks/Chunk_*.txt` — READ ALL OF THEM
- Per-finding function contexts: `real-repos/<REPO_NAME>/scripts/findings/functions/<id>.txt`
- Actual source: follow the `Source:` path from each finding.

## Scan command (for run metadata)

```
./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir real-repos/<REPO_NAME>/scripts/chunks -context-dir real-repos/<REPO_NAME>/scripts/findings/functions real-repos/<REPO_NAME>
```

## Workflow

1. Read the template `plans/skills/false-positive-audit/SKILLS.md`.
2. Read every chunk file under `real-repos/<REPO_NAME>/scripts/chunks/`. Each chunk lists findings with: index (Finding N/M), Source path:line:col, Rule id, title, severity, message, fix, context excerpt. Collect the full list of finding IDs.
3. For a rule you are unsure about, run `./bin/goslop -explain <RULE-ID>` to get the exact rule condition. The decision must be based on the rule condition and the shown source, never on application-specific knowledge.
4. Classify EVERY finding as `False positive`, `True positive`, or `Uncertain`. Work through findings in order (1..M).
5. For every finding you classify as a false positive: read its function context file `real-repos/<REPO_NAME>/scripts/findings/functions/<id>.txt` (and the enclosing source when the context is insufficient), then add a subsection per the template with the smallest source excerpt that proves the decision.
6. Group findings ONLY when they reference the exact same source construct (same file, same line, same rule, same trigger). When grouping, list every grouped ID explicitly in the subsection heading and keep a single excerpt. Do not group distinct lines.
7. True positives: list them compactly in one table per rule (finding id, source, one-line reason) — the template does not have a true-positive section, so append a `## True positives` section before Final evidence.
8. Uncertain: one subsection per finding with the template's structure.
9. Fill the classification summary table with exact counts and all finding IDs grouped by classification.
10. Run `git diff --check` in the goslop repo root after writing the report and record pass/fail in Final evidence.

## Return value

Reply with: repo name, total findings, counts (FP / TP / Uncertain), and the top 5 most common false-positive patterns (rule id + reason) with counts. Keep it under 20 lines.
