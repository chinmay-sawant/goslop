# Contributing to Goslop

Thanks for helping improve Goslop. This guide covers the local workflow, how to keep scan behavior trustworthy, and how we document review-led work.

## Before you start

- Use the Go toolchain declared in [`go.mod`](go.mod) (currently Go 1.26.4). The project builds as pure Go; `CGO_ENABLED=0` is the default Makefile setting.
- Start from an up-to-date `main` branch and use a focused branch for one coherent change.
- Keep generated scan exports, binaries, and local cache data out of commits. In particular, do not commit `bin/`, `.goslop-cache/`, `scripts/findings/`, or `scripts/chunks/`.

```sh
git switch main
git pull --ff-only
git switch -c <type>/<short-description>
go mod download
make build
```

## Choose the right size of change

Small documentation, test, or localized implementation changes usually need only focused tests and a concise PR description.

For a cross-cutting change—scan lifecycle, cache behavior, configuration, report formats, rule registries, or a remediation batch—create and maintain one phase-wise checklist. The checklist should state the problem, intended behavior, acceptance evidence, and whether each row is complete, pending, or intentionally deferred:

- `[x]` — implemented and backed by source, tests, or benchmark evidence.
- `[ ]` — concrete work that is approved and still required.
- `[~]` — conditional follow-up; do not implement it until its stated trigger is real.

The current review-led examples live in [`plans/v0.0.1/reviews/`](plans/v0.0.1/reviews/):

- [`ponytail-code-and-architecture-review.md`](plans/v0.0.1/reviews/ponytail-code-and-architecture-review.md)
- [`ponytail-post-remediation-review.md`](plans/v0.0.1/reviews/ponytail-post-remediation-review.md)
- [`architecture-post-remediation-review.md`](plans/v0.0.1/reviews/architecture-post-remediation-review.md)
- [`go-code-style-and-design-review.md`](plans/v0.0.1/reviews/go-code-style-and-design-review.md)

When a review is intended for wider discussion, include a companion HTML report beside the Markdown ledger. Keep the Markdown checklist and HTML summary synchronized after evidence closes a row.

## Implementation expectations

- Preserve user-facing scan, JSON, SARIF, and export contracts unless the change explicitly documents a migration.
- Prefer scan-local state over mutable package globals. Keep worker concurrency bounded and make cleanup/error behavior observable.
- Add a regression for every correctness fix. Exercise both the success path and the failure path when ownership, cleanup, caching, or partial parsing is involved.
- Keep configuration honest: do not document or retain inert settings. Configuration changes need a test at the resolved runtime behavior seam, not only a parser test.
- For detector or rule changes, follow the existing rule metadata, fixture, registry, and integration-test conventions in the affected package. Do not add blanket lint exclusions or ignore directives to hide a real issue.

## Validate your change

Run the narrowest useful command while developing, then run the applicable final checks before requesting review.

```sh
# Formatting, vet, and the repository's enabled linters.
make lint-all

# Unit and integration suites.
make test

# Required when shared scan state, caches, or concurrency changes.
go test -race ./...

# Pure-Go build contract.
CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop
```

Documentation-only changes normally need `git diff --check`; do not claim a full test run unless it actually ran for the final diff.

```sh
git diff --check
```

## Product scan and benchmark baseline

`make run` performs the product-style scan: profile `all`, no terminal output, no failure exit, context/chunk export, and a cold cache. Its default corpus is GopdfSuit through `SCAN_PATH`; override it for a different real application.

```sh
make run
make run SCAN_PATH=/path/to/go-project
```

The output has this shape:

```text
scanned <files> files (<lines> lines) in <scan-time>
  cache: <hits> hits, <misses> misses (...)
  skipped <files> files
<findings> findings
  severity: <high> high, <info> info, <low> low, <medium> medium
  top rules: <rule counts>
exported <context> context file(s) to scripts/findings/functions; exported <chunks> chunk file(s) to scripts/chunks
```

The current GopdfSuit reference run is a behavioral baseline, not a cross-machine timing promise:

```text
scanned 78 files (28042 lines) in 227.3ms
915 findings
  severity: 10 high, 197 info, 312 low, 396 medium
  top rules: BP-1 ×181, PERF-6 ×94, PERF-32 ×59, BP-5 ×50, PERF-230 ×44
exported 915 context file(s); exported 37 chunk file(s)
```

For detector, cache, parser, export, or performance work, compare against the same corpus with the hard parity target:

```sh
make reference-metrics REFERENCE_PATH=/home/chinmay/ChinmayPersonalProjects/gopdfsuit
```

The scanner may return a nonzero status because findings exist; the `reference-metrics` target intentionally records the metrics and compares the expected counts. A valid reference run matches 915 findings, the severity distribution above, the listed top rules, and 915 context plus 37 chunk exports. Do not report a benchmark as publishable without the exact command, corpus, cache mode, and output.

Use `make bench` for Go benchmark allocations and operation timings when you change a hot path:

```sh
make bench
make bench BENCHTIME=20x
```

## Breaking changes and migrations

Call out every compatibility change in the PR body and in the affected documentation or schema. Include the old behavior, new behavior, migration step, and validation proof.

Configuration notes:

- `languages` is a supported key again (v0.0.2 Phase 3). Default when unset is `["go"]`. Unknown tokens and an explicit empty list are rejected. Only registered language plugins contribute detectors and file extensions; enabling `python` before a Python plugin is registered is a no-op for that language (no crash).
- `typed.enabled` remains unsupported and is still rejected at load time.

If you change a CLI flag, config field, output schema, rule identifier, baseline fingerprint, or exported file format, add migration guidance and compatibility tests before marking the checklist row complete.

## Pull requests

Use [`plans/PR/PR_TEMPLATE.md`](plans/PR/PR_TEMPLATE.md) as the base for every PR:

1. Fill a copy under `plans/PR/pr-<short-slug>.md` before opening the PR.
2. Use a conventional title such as `fix(export): return cache cleanup errors`.
3. Summarize the behavior, motivation, impact, breaking changes, exact validation commands, and benchmark results where applicable.
4. Self-assign, apply appropriate labels, and link a real issue when one exists.
5. Keep the committed PR body file and the GitHub PR description synchronized.
6. Open a draft PR unless it is explicitly ready for review. Contributors do not merge their own PRs; wait for the repository maintainer's review and merge decision.

Before opening the PR, inspect the final diff for unrelated changes and secrets:

```sh
git status -sb
git diff --check
git diff --stat origin/main...HEAD
```

Thank you for keeping Goslop's findings, performance claims, and review evidence reproducible.
