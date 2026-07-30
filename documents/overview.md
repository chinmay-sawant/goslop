# goslop product overview

goslop is a **pure-Go** static analysis tool (SAT) for Go codebases. It uses `go/parser` + `go/ast` (no CGO, no tree-sitter) to find performance issues, bad practices, and CWE-class security problems.

## What it does

| Capability | Description |
|------------|-------------|
| **PERF rules** | Hot-path heuristics (loops, allocations, HTTP timeouts, framework footguns). **239** registered. |
| **CWE structural** | Catalogue of structural security heuristics. **175** registered. |
| **Taint (experimental)** | Inter-procedural graph for **CWE-22 / 78 / 79 / 89** (path traversal, command injection, XSS, SQLi). |
| **Bad practices (BP)** | Style / hygiene / project-level rules (`BP-*`). On for `style` and `all` profiles. |
| **Profiles (packs)** | `recommended`, `perf`, `security`, `style`, `all` - curated rule surfaces and defaults. |
| **Reporters** | `text` (default), `json`, `sarif` (SARIF 2.1.0). Summary always on **stderr**. |
| **Exports** | Per-finding context (`scripts/findings/functions`) and batched chunks (`scripts/chunks`) for agent work. **Context defaults to the whole enclosing function** (`[goslop.export] whole_function = true`). |
| **Cache** | Incremental per-file cache under `.goslop-cache/`. |
| **Baseline** | Suppress known findings via `.goslop-baseline.json`. See [suppressions-and-baselines.md](./suppressions-and-baselines.md). |
| **Ignore** | Inline `// goslop-ignore` directives and path ignores. See [suppressions-and-baselines.md](./suppressions-and-baselines.md). |

## Status (shipped)

Phases **0-12** landed. High-level status:

| Area | Status |
|------|--------|
| Engine / CLI / text · JSON · SARIF | ✅ |
| PERF detectors | ✅ 239/239 |
| CWE structural | ✅ 175/175 |
| Bad practices | ✅ catalogue + project-level rules |
| Taint graph | ✅ CWE-22/78/79/89 |
| Cache / baseline / ignore | ✅ |
| Packs | ✅ recommended / perf / security / style / all |
| §12.4 product parity baseline (`gopdfsuit`) | ✅ 915 findings; exports 915 context + 37 chunks |

Details: root [`README.md`](../README.md) and [`plans/port-phasewise-checklist.md`](../plans/port-phasewise-checklist.md).

## Architecture (product view)

```text
CLI flags + goslop.toml + profile
        │
        ▼
  ScanContext (only/skip, fail policy, taint, BP, cache, baseline)
        │
        ▼
  Analyzer (walk → parse → detectors → finalize)
        │
        ▼
  ignore directives → baseline filter
        │
        ├── stderr: scan summary
        ├── stdout: text | json | sarif findings
        └── optional disk export: functions + chunks
```

| Layer | Path |
|-------|------|
| CLI entry | `cmd/goslop` |
| App orchestration | `internal/app` |
| Core (profiles, fail, plugins) | `internal/core` |
| Engine (walk, cache, baseline, ignore) | `internal/engine` |
| Detectors | `internal/lang/go/detectors` |
| Reporting | `internal/reporting` |
| Export | `internal/export` |
| Rule metadata | `ruleset/golang/` (default); WIP Python seeds in `ruleset/python/` |

## Profiles (packs)

Default: **`--profile recommended`**.

| Profile | Aliases | Contents | Taint | BP | Default fail |
|---------|---------|----------|-------|----|--------------|
| `recommended` | `ci`, `default` | S-tier PERF + taint-core CWE IDs | off | off | **high** |
| `perf` | `performance` | S + A tier PERF | off | off | **high** |
| `security` | `sec` | Security CWE pack | **on** (depth 4) | off | **high** |
| `style` | `bp`, `bad-practices` | `BP-*` (skips BP-21/28/30 by default) | off | **on** | **none** |
| `all` | `full` | Full catalogue | off | **on** | **medium** |

Exact allow-lists and rationale: [go-recommended-pack.md](./go-recommended-pack.md).  
Full flag list: [cli-reference.md](./cli-reference.md).

## Scanning the goslop repo itself (high finding counts)

If you point goslop at **this repository** (`SCAN_PATH=.` or the goslop tree), expect a **large** finding count. That is expected and is a poor measure of product false-positive rate.

### Why (SAT self-host, not “broken rules”)

goslop is a static analysis tool (SAT). Like **Semgrep**, **CodeQL**, or similar engines run **on their own source**, a self-scan fires heavily because the tree is full of **what detectors look for**:

| In this repo | What happens |
|--------------|--------------|
| **Pattern catalogues** (`needles.go`, rule tables, metadata) | Rules match **string literals** that encode vulnerable or sloppy snippets used as search needles—not live application code. |
| **Detector implementations** | Source contains the same substrings the scanners search for (`"http.Server{"`, `"interface {"`, Gin/SQL examples, etc.). |
| **Test fixtures** (`tests/fixtures/...`) | Deliberate **vulnerable** and **safe** samples for the matrix; many findings are **by design**. |
| **Product code** (`cmd/`, most of `internal/` outside detectors) | The smaller slice that resembles a normal Go project. |

So a self-scan looks “noisy” for the same reason Semgrep (or any pattern-based SAT) flags its own rule packs and fixture corpora: **the tool is matching the patterns it was built to detect**, which live in the repo as data and tests.

### What to do instead

| Goal | Approach |
|------|----------|
| **Product / baseline metrics** | Scan a real application tree (default `make run` → **gopdfsuit**). See [make-run.md](./make-run.md). |
| **Judge FPs on real apps** | Use application code and CI profiles (`recommended`, etc.), not full self-scan under `--profile all`. |
| **Dogfood hang / smoke** | Optional short scan of detector packages with a wall clock; do not use finding count as quality score. |
| **Clean product debt in goslop** | Scan `cmd/` + `internal/` while **excluding** `internal/lang/go/detectors/**` and `tests/fixtures/**`, then triage. |

True false positives on **normal application code** still deserve detector improvements. Self-host and fixture volume do **not** mean “fix every finding in this repo until zero.”

### Recommended pack (summary)

**PERF S-tier:** `PERF-1`, `PERF-7`, `PERF-50`, `PERF-58`, `PERF-71`, `PERF-101`, `PERF-103`, `PERF-116`, `PERF-189`, `PERF-190`

**Taint-core CWE allow-list (taint engine off until enabled):**  
`CWE-22`, `CWE-78`, `CWE-79`, `CWE-89`, `CWE-90`, `CWE-91`

Under default fail **high**, S-tier PERF is typically **medium** (visible, does not fail CI). High/critical CWE findings **do** fail.

## Detector families

### PERF (`PERF-*`)

- **239** rules across loop allocations, HTTP, frameworks (Gin, etc.), parsing, data access.
- Metadata chunks: `ruleset/golang/chunks/perf-*.json`
- Human notes (partial): [perf-rules.md](./perf-rules.md)
- List live: `./bin/goslop --list-rules` and filter for `PERF-`

### CWE (`CWE-*`)

- **175** structural rules.
- Taint-lite same-file heuristics for injection-class CWEs when full taint is off.
- Full taint graph: enable with `--taint` or `--profile security`. See [taint.md](./taint.md).

### Bad practices (`BP-*`)

- Style, error handling, testing, HTTP/server, go.mod hygiene, resources.
- Enabled only for **`style`** and **`all`** (unless overridden in config).
- Style profile skips noisy rules `BP-21`, `BP-28`, `BP-30` by default; re-enable with `--only BP-28`.

### Taint

```sh
# Explicit
./bin/goslop --taint --taint-depth 3 --taint-show-paths .

# Security profile turns taint on (default depth 4)
./bin/goslop --profile security .

# Force off even under security
./bin/goslop --profile security --no-taint .
```

Honesty bar and source/sink model: [taint.md](./taint.md).

## Suppressions

### Inline ignore directives

Comment-only (not inside strings):

```go
// goslop-ignore: CWE-78
exec.Command("sh", "-c", cmd)

// goslop-ignore: PERF-101, PERF-190
srv := &http.Server{Addr: ":8080"}

// goslop-ignore-file
package generated

// goslop-ignore-start: BP-1
// … noisy region …
// goslop-ignore-end
```

- Default: suppressed findings are **dropped**.
- `--show-ignored`: keep them (severity forced to **info**, marked suppressed).

### Path ignores

Walk respects `.gitignore`, `.goslopignore`, `.ignore`, plus config `include` / `exclude` globs. Test files (`*_test.*`) are excluded by default (`--include-tests` to include).

## Baseline

File: **`.goslop-baseline.json`** (discovered upward from cwd unless `--baseline-file` / config path).

| Flag | Effect |
|------|--------|
| (default) | Drop findings that match baseline fingerprints |
| `--no-baseline` | Ignore baseline entirely |
| `--baseline-file PATH` | Explicit baseline path |
| `--show-baselined` | Keep baselined findings as info / suppressed |

Match priority: fingerprint (`goslop:2:{rule}:{file}:{msgHash16}`), else file/line/column.

> **Note:** The Go CLI loads and filters baselines; it does not expose a baseline-generation subcommand. See [suppressions-and-baselines.md](./suppressions-and-baselines.md) for the file format and rollout guidance.

## Incremental cache

| Item | Default |
|------|---------|
| Directory | `.goslop-cache/` |
| Enabled | yes |
| Disable | `--no-cache` or config `cache.enabled = false` |
| Rebuild | `--rebuild-cache` |
| Prune | `--prune-cache` (then exit) |

Cache keys on content hash + tool/rule-config fingerprint. Pipeline details: [architecture-performance.md](./architecture-performance.md).

## Exit codes

| Code | Meaning |
|------|---------|
| **0** | Clean, or meta-command success, or `--no-fail` |
| **1** | Findings violate the active fail policy |
| **2** | Usage / config / unknown profile or rule |
| **3** | Internal (engine, I/O, cache hard failure) |

## Configuration

```sh
./bin/goslop init   # writes ./goslop.toml
```

Template: [`templates/goslop.toml`](../templates/goslop.toml).  
Schema: [`goslop.schema.json`](../goslop.schema.json).  
Merge rules: [cli-reference.md](./cli-reference.md#config-file-and-cli-merge).

### Export context mode

When using `--export-context` / `--export-chunks`, each finding’s **Context** block
defaults to the **full enclosing Go function** (outermost `FuncDecl` when the hit
sits inside a nested closure such as `defer func()`). Opt out with:

```toml
[goslop.export]
whole_function = false   # nearby ~4-line window instead
```

Details: [export-context-and-chunks.md](./export-context-and-chunks.md).

## Outputs at a glance

| Output | How | Docs |
|--------|-----|------|
| Terminal findings | default `--format text` | [reporting-formats.md](./reporting-formats.md) |
| JSON envelope | `--format json` | [reporting-formats.md](./reporting-formats.md) |
| SARIF 2.1.0 | `--format sarif` | [reporting-formats.md](./reporting-formats.md#sarif-210) |
| Per-finding refs (whole function by default) | `--export-context` → `scripts/findings/functions/N.txt` | [export-context-and-chunks.md](./export-context-and-chunks.md) |
| Batches for agents (same Context) | `--export-chunks` → `scripts/chunks/Chunk_A_B.txt` | [export-context-and-chunks.md](./export-context-and-chunks.md) |
| Product run | `make run` | [make-run.md](./make-run.md) |

## Typical workflows

```sh
# 1) CI gate (high-signal)
./bin/goslop --profile recommended --format sarif . > goslop.sarif

# 2) Full triage + agent batches
./bin/goslop --profile all --no-fail --export-context --export-chunks --no-cache .
# then open scripts/chunks/Chunk_1_25.txt for a batch,
# or scripts/findings/functions/42.txt for one finding

# 3) Security-focused with taint
./bin/goslop --profile security --taint-show-paths ./cmd

# 4) Makefile product scan (summary + exports)
make run SCAN_PATH=./your/go/project
```
