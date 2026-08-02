# Post-fix re-audit instructions (read first)

You are re-auditing ONE repo after the FP-reduction fix (`b5b8fde`, binary rebuilt 2026-08-02 16:29). Read `plans/skills/false-positive-audit/SKILLS.md` for the template. Work only with your assigned repo and mode.

## Mode A — remaining false positives (repo has fresh findings beyond its audited TP count)

Goal: classify EVERY fresh-scan finding as FP/TP/Uncertain and document the remaining false positives.

1. Fresh evidence (from the NEW scan, post-fix):
   - Chunks: `scripts/<REPO>/chunks/Chunk_*.txt` — READ ALL
   - Contexts: `scripts/<REPO>/findings/functions/<id>.txt`
   - Old audit: `plans/fp-validations/reports/<REPO>.md` — its classification summary lists audited TP finding IDs (with sources in the True positives tables) and FP IDs.
2. Fresh finding IDs do NOT correspond to old IDs. Match by `Source:` path (file:line:col) — a fresh finding whose source matches an audited TP is a true positive; everything else is either a re-appearing audited FP or a new finding.
3. Classify every fresh finding: TP (matches audited TP by source) / FP (matches audited FP source or new finding that fails its rule condition) / Uncertain.
4. For every fresh finding classified FP, read its function context (`scripts/<REPO>/findings/functions/<id>.txt`) and the enclosing source; verify against the rule condition (`./bin/goslop -explain <RULE>` when unsure).
5. APPEND (do not edit existing sections) to `plans/fp-validations/reports/<REPO>.md`:

```
## Post-fix remaining-FP audit (2026-08-02)

Run metadata (fresh scan), scan evidence (fresh findings count), audit checklist, classification summary of the fresh run (FP/TP/Uncertain counts with fresh finding IDs), then template False positives subsections for each remaining FP (group only identical source constructs, list every grouped ID), Uncertain findings if any, and Final evidence.
```

## Mode B — over-suppressed true positives (repo's fresh findings are FEWER than its audited TP count)

Goal: identify which audited true positives were suppressed by the fix and document them for review.

1. Old audit: `plans/fp-validations/reports/<REPO>.md` — collect the full audited TP list (finding ID, Source, rule, one-line reason from the True positives tables).
2. Fresh evidence: `scripts/<REPO>/chunks/` — collect the sources present in the fresh scan.
3. A TP is over-suppressed if its old `Source:` (file:line) is NOT present in the fresh scan. For each such TP, read the current source file at that location to confirm the construct still exists (suppressed, not fixed) — if the code was actually removed, mark it as fixed/removed instead.
4. APPEND to `plans/fp-validations/reports/<REPO>.md`:

```
## Post-fix over-suppression audit (2026-08-02)

Run metadata, then a table: old finding ID | rule | Source | one-line reason (from old audit) | current status (suppressed-but-present / fixed-removed). Then one short subsection per suppressed-but-present TP with the smallest source excerpt proving the construct is still there and why it satisfies the rule condition.
```

## Return value

Repo, mode, fresh findings count, counts: (Mode A) fresh FP/TP/Uncertain with remaining-FP count; (Mode B) over-suppressed TPs found, fixed-removed count. Under 15 lines. Also run `git diff --check` in the goslop repo root and report pass/fail.
