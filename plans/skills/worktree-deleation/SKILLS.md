---
name: worktree-deleation
description: >
  Execute a phased markdown checklist plan via parallel git worktree subagents
  on Grok 4.5 (cursor-grok-4.5-high-fast): one branch/PR per workstream, filled
  PR bodies from plans/PR/PR_TEMPLATE.md, self-assign + labels, optional epic
  integration PR, merge integration only, close child PRs, delete all
  non-master local/remote branches, pull master. Use when the user points at a
  plans/*.md checklist and wants parallel worktrees, multi-PR ship, or "same
  as last architecture batch". Triggers: worktree, parallel agents, phase
  checklist, integrate then delete branches, /worktree-deleation.
---

# Worktree multi-stream plan execution

> **Subagent model (required):** launch every workstream / integration subagent
> with **Grok 4.5** in high + fast mode — Cursor model slug
> `cursor-grok-4.5-high-fast`. Do not use another model unless the user
> explicitly overrides.

Run a **markdown checklist plan** end-to-end the same way as the architecture
review batches: parallel worktrees → tested child PRs → one integration PR →
merge only the integration → close children → wipe non-`master` branches →
pull `master`.

## Input

**Required:** path to a markdown plan with phases / workstreams and checkboxes.

Example:

```text
plans/v0.0.5/rust-architecture-review.md
```

If the user did not pass a path, ask for it once. Do not invent workstreams.

## Output contract

1. GitHub **issues** (epic + one child per implementable workstream).
2. **One branch + one PR per child** (worktree-isolated subagent).
3. Each PR body **filled** (not empty template copy) under `plans/` and committed on the branch.
4. Optional **integration branch + PR** when 2+ children ship.
5. After user asks to land: **merge integration only**, close children, **delete all non-master branches** (local + remote), **pull master**.

---

## Phase A — Read the plan and split work

1. Read the plan file fully.
2. Identify **implementable** workstreams (unchecked `[ ]` items grouped by phase/section).
   - Skip pure evidence/ledger sections already `[x]`.
   - Skip exit-gate-only phases until children exist (those become the integration validation).
3. Cap parallelism at **4–5 subagents** unless the user asks for more.
4. Prefer independent seams to reduce merge conflicts; if two streams touch the same trait/engine file, still launch in parallel and resolve on the integration branch.

Record for each stream:

| Field | Rule |
|-------|------|
| Issue title | Imperative / scoped (`fix(taint): …`, `refactor(go): …`) |
| Branch | `fix/…`, `refactor/…`, `chore/…` from `origin/master` |
| Scope | Checkbox list + success condition from the plan |
| Out of scope | Sibling streams (do not steal work) |
| Labels | See Labels below |

---

## Phase B — Issues first

Follow `plans/PR/ISSUE_TEMPLATE.md`.

1. Create **one epic** for the plan (or parent phase set).
2. Create **one child issue per workstream**.
3. Always: `--assignee "@me"` and at least one label.
4. Issue bodies: Context, Scope, Out of scope, Success criteria, Plan path, References.
5. Never invent issue numbers. Create issues before PRs.

```sh
gh issue create --title "…" --assignee "@me" --label enhancement --body-file …
```

---

## Phase C — Parallel worktree subagents

Launch **one subagent per workstream** with:

- `model: cursor-grok-4.5-high-fast` (Grok 4.5 high + fast — required)
- `isolation: worktree`
- `subagent_type: general-purpose`
- `background: true`

Each agent prompt must include:

1. Issue number (`Closes #N`, `Relates to #EPIC`).
2. Branch name from `origin/master`.
3. Exact checklist items and success condition (paste plan text; do not rely on uncommitted plan edits on master).
4. **Out of scope** (sibling workstreams).
5. Validation commands (`make lint`, focused tests, `make test` when feasible).
6. PR process:
   - Write **filled** body to `plans/v0.0.x/pr-<slug>.md` (or `plans/PR/pr-<slug>.md`) using sections from `plans/PR/PR_TEMPLATE.md`.
   - **Fill data** — real Summary, Changes, Test plan with results, Related issues, Impact. Do not leave template placeholders.
   - Commit body + code on the branch.
   - Push and `gh pr create` with metadata below.
7. Return: branch, PR URL, test summary, blockers, conflict notes.

### PR open metadata (required)

```sh
gh pr create \
  --base master \
  --head "$(git branch --show-current)" \
  --title "<type>(<scope>): <short imperative description>" \
  --body-file plans/…/pr-<slug>.md \
  --assignee "@me" \
  --label <label> \
  [--label <label>]
```

### Labels

| Title type | Labels |
|------------|--------|
| `feat` / `refactor` / `chore` product work | `enhancement` |
| `fix` correctness | `bug` (+ `enhancement` if hardening) |
| docs / plan records | `documentation` |
| mixed code + plans | `documentation` + `enhancement` (or `bug`) |

### Title convention

```text
<type>: <short imperative description>
```

Optional scope: `fix(go): …`, `refactor(engine): …`.

---

## Phase D — Integration branch (2+ children)

Follow multi-workstream rules in `plans/PR/PR_TEMPLATE.md`.

```sh
git fetch origin master
git checkout -B chore/epic-N-integration origin/master

# Prefer docs/quality first, then detectors, engine/taint last
for b in origin/child-a origin/child-b origin/child-c; do
  git merge "$b" -m "merge: integrate $b into epic-N integration"
done

make lint
make test
# plus area-focused tests

git push -u origin HEAD
gh pr create \
  --base master \
  --title "chore: integrate epic #N workstreams" \
  --body-file plans/…/pr-epic-N-integration.md \
  --assignee "@me" \
  --label enhancement \
  --label documentation
```

Integration PR body must:

- Table of child issue / branch / standalone PR
- Combined validation results
- `Closes` every child issue (+ epic when complete)
- Note that **child PRs are superseded** by the integration PR

Comment on each child PR:

```markdown
## Integration

This branch is also merged into `chore/epic-N-integration` for combined validation.
Prefer reviewing/merging the integration PR when present.
```

Resolve conflicts on the integration branch; re-run full tests after every resolution.

Update plan checkboxes for completed workstreams on the integration branch.

---

## Phase E — Land (when user asks to merge)

Default: **merge only the integration PR** into `master`.

```sh
gh pr merge <INTEGRATION_PR> --merge --delete-branch=false
```

Then:

1. Close open child PRs **without merging** (if still open):

```sh
gh pr close <N> --comment "Superseded by integration PR #<I> (merged). Closing without merge."
```

2. If GitHub already marked children `MERGED` because commits landed via integration, leave them; do not re-merge.

3. **Delete all branches except `master`** (local and remote):

```sh
git checkout master

# Local
git branch --format='%(refname:short)' | grep -v '^master$' | xargs -r git branch -D

# Remote: only heads that still exist
git ls-remote --heads origin | awk '{print $2}' | sed 's|refs/heads/||' | grep -v '^master$' \
  | while read -r b; do git push origin --delete "$b" || true; done

git fetch origin --prune
git pull origin master
```

4. Confirm:

```text
* master
  remotes/origin/master
```

Remote heads must list **master only** (unless the user exempted branches).

---

## Phase F — Single-stream shortcut

If the plan has **one** workstream only:

- Skip integration branch.
- One child PR targeting `master`.
- On land: merge that PR, close issue (via `Closes #N`), delete that feature branch, pull `master`.

---

## Validation gates (every stream + integration)

Prefer:

```sh
make lint
make test
```

Also when the plan requires docs:

```sh
RUSTDOCFLAGS='-D warnings' cargo doc --all-features --no-deps --locked
```

Focused tests for the area touched before the full suite when the suite is long.

Do not open a PR with a red suite unless the user explicitly accepts known flakes and documents them in the PR body.

---

## Agent prompt skeleton (copy into each worktree subagent)

```text
You implement issue #N (plan PATH §X) in an isolated git worktree.

Branch: <name> from origin/master.
Closes #N · Relates to #EPIC.

Requirements:
- [ ] …paste checkboxes…

Out of scope: sibling streams …

Validation:
make lint
<focused cargo test …>
make test

PR:
1. Filled body → plans/…/pr-<slug>.md (PR_TEMPLATE sections, real content)
2. Commit message: <type>(scope): …
3. git push -u origin HEAD
4. gh pr create --base master --assignee "@me" --label … --body-file …

Return: branch, PR URL, tests, blockers.
```

---

## Hard rules

- **Input is the plan path** — all workstreams come from that file.
- **Self-assign** every issue and PR (`--assignee "@me"`).
- **Labels required** on every PR/issue.
- **Filled PR bodies only** — no empty template paste.
- **Never invent issue numbers.**
- **Prefer merge integration only** for multi-stream epics.
- **Delete every non-master branch** local + remote after land (unless user exempts).
- **Pull master** last and report clean status.
- Respect `action_safety`: confirm only when the user has not already ordered merge/delete; when they say merge/delete/pull, execute.

---

## Done criteria

- [ ] Epic + child issues exist and are assigned
- [ ] Each stream: branch, green tests, filled PR, assignee, labels
- [ ] Integration PR green (if multi-stream)
- [ ] Plan checkboxes updated for completed work
- [ ] User land request: integration merged, children closed/superseded
- [ ] Only `master` remains local and remote
- [ ] Local `master` matches `origin/master`
