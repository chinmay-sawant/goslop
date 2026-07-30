# goslop documentation

User-facing documentation for the **goslop** static analysis tool (SAT) for Go.

goslop scans Go source for **PERF** hot-path issues, structural **CWE** heuristics, experimental **taint** flows, and **bad-practice** (`BP-*`) style rules. Output formats: **text**, **JSON**, and **SARIF 2.1.0**. Optional exports write per-finding context under `scripts/findings/functions/` and batched chunks under `scripts/chunks/` for agent delegation; by default each finding’s **Context** is the **whole enclosing function** (`[goslop.export] whole_function = true`).

## Start here

| Document | What it covers |
|----------|----------------|
| [**overview.md**](./overview.md) | Product features, detector families, profiles, suppressions, cache, baseline; **why self-scan of this repo is noisy** (SAT self-host) |
| [**cli-reference.md**](./cli-reference.md) | Every CLI flag, subcommands, exit codes, config merge |
| [**make-run.md**](./make-run.md) | Product `make run` / `make reference-metrics` / `make bench`; reference corpus vs self-scan; perf numbers |
| [**reporting-formats.md**](./reporting-formats.md) | Text, JSON, and **SARIF** with full examples |
| [**export-context-and-chunks.md**](./export-context-and-chunks.md) | Function refs vs chunks; whole-function Context (default); AI delegation |

## Deeper topic guides (existing)

| Document | What it covers |
|----------|----------------|
| [**go-recommended-pack.md**](./go-recommended-pack.md) | Default `recommended` pack rule lists and CI guidance |
| [**perf-rules.md**](./perf-rules.md) | Human notes for many `PERF-*` rules (partial catalog) |
| [**taint.md**](./taint.md) | Experimental taint engine (CWE-22/78/79/89) |
| [**architecture-performance.md**](./architecture-performance.md) | Engine pipeline, cache layout, performance design |

## Quick commands

```sh
# Build (pure Go, no CGO)
make build

# Everyday scan (default profile: recommended)
./bin/goslop .

# Full catalog + exports (product-style)
make run SCAN_PATH=./your/project

# Machine formats
./bin/goslop --format json .
./bin/goslop --format sarif . > goslop.sarif

# List / explain rules
./bin/goslop --list-rules
./bin/goslop --explain PERF-6

# Starter config
./bin/goslop init
```

## Related repo paths

| Path | Role |
|------|------|
| [`README.md`](../README.md) | Project status and install |
| [`templates/goslop.toml`](../templates/goslop.toml) | Config template |
| [`goslop.schema.json`](../goslop.schema.json) | Config JSON Schema |
| [`plans/`](../plans/) | Port phase ledger and architecture notes |
| [`scripts/findings/functions/`](../scripts/findings/functions/) | Generated per-finding context (after export; whole function by default) |
| [`scripts/chunks/`](../scripts/chunks/) | Generated finding batches (after export; same Context mode) |
