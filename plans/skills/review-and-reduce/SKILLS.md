---
name: review-and-reduce
description: Reduce confirmed static-analysis false positives without suppressing valid findings. Use only when a user supplies a false-positive review checklist path and asks to narrow detectors, correct anchors, add guardrails, add variant text fixtures, update fixture-driven tests, or validate a noise-reduction change.
---

# Review and reduce false positives

Require a review checklist path from the user. This path is the source of truth for scope.

If the user has not supplied that path, do not inspect, modify, or validate rules. Ask for the checklist path and wait.

Read the supplied checklist fully before editing. Treat its `False positives`, `Uncertain findings`, and `True positives` sections as distinct sets.

## Workflow

1. Read each selected false-positive entry, its `./scripts/findings/functions/<id>.txt` context, and the referenced source block. Identify the source language from the source path and fixture metadata.
2. Explain the exact false-positive mechanism from the full function context before changing code. Do not infer a guardrail from the rule title alone.
3. Locate the rule implementation, registration, canonical text fixtures, package tests, and integration coverage. Prefer the codebase graph for code discovery; use text search for Markdown, fixtures, and literal rule IDs.
4. State the missing rule precondition and add the narrowest guardrail that excludes the audited function pattern while preserving the actual unsafe sink or condition.
5. Preserve valid findings. Do not delete a rule, lower severity, add a broad ignore, or weaken matching globally merely to remove the reviewed hit.
6. Add the same source pattern that caused the false positive as a language-appropriate **safe variant text fixture**. Preserve the relevant data shape, call form, and surrounding context from the exported function file.
7. Add or retain a matching **vulnerable variant text fixture** for the same rule. The safe variant must not emit; the vulnerable variant must still emit.
8. Update the per-rule test and integration matrix to reference the canonical text fixtures. Do not embed replacement source snippets in test code when the repository uses fixture files.
9. Run the focused detector tests, affected fixture matrix, required lint, and full test commands. Re-run the audit scan and verify the reviewed false-positive pattern is absent while the vulnerable variant still fires.

## Uncertain and true-positive handling

- Do not change a rule for an entry under `Uncertain findings` without additional source evidence and user direction.
- Use true-positive entries as non-regression cases. Add a fixture or test when the detector change could affect the same rule family.
- Re-run the audit scan after the fix. Confirm that the reviewed false positive is gone and the retained true-positive case still fires.

## Completion record

For each repaired finding, record:

- finding ID and rule ID;
- detector condition changed;
- safe-variant fixture and vulnerable-variant fixture paths;
- per-rule test, integration fixture-matrix test, lint, and full-test results;
- rescan result and any remaining uncertainty.

Keep the review checklist's `[ ]` entries unchanged unless the user explicitly asks to update that audit artifact. Record implementation progress in the requested implementation plan, PR body, or a separate reduction ledger.
