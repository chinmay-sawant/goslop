# CLI reference

Complete command-line reference for the Go goslop binary (`./bin/goslop`).

```text
Usage: goslop [flags] [PATH...]
       goslop init
```

- Flags accept both `-flag` and `--flag` forms (stdlib `flag` package).
- If no `PATH` is given, the scan root defaults to **`.`**.
- Scan **summary** always prints to **stderr**. Finding payloads go to **stdout** (except `--no-terminal` with text format).

---

## Subcommands and meta modes

| Mode | Invocation | Behavior |
|------|------------|----------|
| **Scan** | `goslop [flags] [PATH...]` | Analyze paths (default) |
| **`init`** | `goslop init` | Write starter `goslop.toml` in cwd; **no flags parsed after `init`** |
| **Help** | `-h` / `-help` | Print usage; exit 0 |
| **Version** | `--version` | Print version (e.g. `0.1.0-dev`); exit 0 |
| **List rules** | `--list-rules` | Print registered rules; exit 0 |
| **Explain rule** | `--explain RULE_ID` | Print catalogue metadata; exit 0 (unknown → exit 2) |
| **Prune cache** | `--prune-cache` | Prune stale cache for PATHS and exit |

### `init`

```sh
./bin/goslop init
# wrote starter goslop.toml to /abs/path/goslop.toml
```

- Fails with exit **2** if `goslop.toml` already exists.
- Content matches [`templates/goslop.toml`](../templates/goslop.toml).

---

## All flags

### Profile and rule filters

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--profile` | string | `recommended` | Product pack: `recommended` \| `perf` \| `security` \| `style` \| `all` (see aliases below) |
| `--only` | CSV | _(empty)_ | Only run these rule IDs (comma-separated). Union with profile allow-list and config `only`. |
| `--skip` | CSV | _(empty)_ | Skip these rule IDs. Union with profile default skips and config `skip`. |

**Profile aliases**

| Canonical | Aliases |
|-----------|---------|
| `recommended` | `default`, `ci` |
| `perf` | `performance` |
| `security` | `sec` |
| `style` | `bp`, `bad-practices`, `bad_practices` |
| `all` | `full` |

**CSV examples**

```sh
./bin/goslop --only "CWE-78,PERF-6,CWE-89" .
./bin/goslop --profile style --only BP-28 .
./bin/goslop --skip PERF-1,PERF-7 .
```

Patterns support exact IDs and prefix wildcards such as `PERF-*`, `CWE-*`, `BP-*` in the filter engine.

---

### Output format

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--format` | string | `text` | Output format: `text` \| `json` \| `sarif` |
| `--no-terminal` | bool | `false` | Print product scan summary only (skip per-finding **text** dump). JSON/SARIF still emit on stdout when requested. |
| `--no-fail` | bool | `false` | Always exit **0** even when findings match the fail policy (still prints summary / exports). |

```sh
./bin/goslop --format text .
./bin/goslop --format json . > findings.json
./bin/goslop --format sarif . > goslop.sarif
./bin/goslop --no-terminal .                    # summary only
./bin/goslop --no-terminal --format json .      # summary on stderr + JSON on stdout
```

Full format schemas and SARIF examples: [reporting-formats.md](./reporting-formats.md).

---

### Rule discovery

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--list-rules` | bool | `false` | List registered rules (`[maturity]  id  title`) and exit |
| `--explain` | string | _(empty)_ | Print catalogue metadata for one rule id and exit |

```sh
./bin/goslop --list-rules
./bin/goslop --explain PERF-6
./bin/goslop --explain CWE-89
```

---

### Tests, paths, and config

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--include-tests` | bool | `false` | Include test files (`*_test.*`) in analysis |
| `--config` | string | _(discover)_ | Path to `goslop.toml` (default: walk upward from first PATH) |
| `--version` | bool | `false` | Print version and exit |

Positional **PATH…** arguments are scan roots (files or directories).

---

### Cache

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--no-cache` | bool | `false` | Disable the incremental analysis cache |
| `--cache-dir` | string | `.goslop-cache` | Incremental cache directory |
| `--rebuild-cache` | bool | `false` | Purge the cache directory before scanning |
| `--prune-cache` | bool | `false` | Prune stale cache entries for PATHS and exit |

```sh
./bin/goslop --no-cache .
./bin/goslop --rebuild-cache .
./bin/goslop --cache-dir /tmp/ch-cache .
./bin/goslop --prune-cache .
```

---

### Baseline and ignore visibility

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--no-baseline` | bool | `false` | Ignore any existing `.goslop-baseline.json` |
| `--baseline-file` | string | _(discover)_ | Explicit baseline path |
| `--show-ignored` | bool | `false` | Report findings suppressed by `goslop-ignore` |
| `--show-baselined` | bool | `false` | Report findings present in the baseline |

---

### Taint

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--taint` | bool | `false` | Enable experimental taint tracking (CWE-22/78/79/89) |
| `--no-taint` | bool | `false` | Disable taint even under `security` profile |
| `--taint-depth` | int | `0` | Inter-procedural hops **1-4**; `0` = profile default (security → 4, else 1) |
| `--taint-show-paths` | bool | `false` | Attach taint hop evidence to findings |

```sh
./bin/goslop --taint --taint-depth 3 --taint-show-paths .
./bin/goslop --profile security .
./bin/goslop --profile security --no-taint .
```

Details: [taint.md](./taint.md).

---

### Export context and chunks

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--export-context` | bool | `false` | Write per-finding context files (default dir `scripts/findings/functions`) |
| `--export-chunks` | bool | `false` | Write chunked finding files (default dir `scripts/chunks`) |
| `--context-dir` | string | `scripts/findings/functions` | Override context export directory |
| `--chunks-dir` | string | `scripts/chunks` | Override chunks export directory |
| `--chunk-size` | int | **25** | Findings per chunk file |

```sh
./bin/goslop --profile all --export-context --export-chunks --no-cache .
./bin/goslop --export-chunks --chunk-size 50 --chunks-dir /tmp/chunks .
```

**Use chunks to delegate work** (each chunk is a batch of combined findings).  
**Use `scripts/findings/functions/N.txt` for individual finding refs.**

Full guide: [export-context-and-chunks.md](./export-context-and-chunks.md).

---

## Exit codes

| Code | Constant | When |
|------|----------|------|
| **0** | clean | No failing findings; help/version/list/explain success; `--no-fail` |
| **1** | failing | Findings violate active fail policy and `--no-fail` is not set |
| **2** | config/usage | Bad flags, unknown profile/format/rule, invalid config, `init` when file exists |
| **3** | internal | Analyze failure, export I/O, hard cache errors, `init` write errors |

Fail policies (from profile / config `fail_on`):

| Policy | Names | Fails on |
|--------|-------|----------|
| none | `none`, `never`, `no_fail`, `nofail` | never |
| high | `high`, `strict` | high + critical |
| medium | `medium`, `medium_as_errors`, `default` | medium and above |

---

## Config file and CLI merge

Discovery order:

1. `--config PATH` if set  
2. Else walk upward from the first scan path (or `.`) for `goslop.toml`  
3. Missing file is fine (CLI-only defaults)

### Common TOML keys (`[goslop]`)

| Key | Effect |
|-----|--------|
| `only` / `skip` | Union with CLI `--only` / `--skip` |
| `fail_on` | Sets fail policy unless `--no-fail` |
| `include` / `exclude` | Walk globs (config only; no CLI flag) |
| `exclude_tests` | Default true; set `false` to include tests (CLI `--include-tests` also includes) |
| `[goslop.baseline]` | `enabled`, `path` |
| `[goslop.cache]` | `enabled`, `path`, `max_size_mb`, `max_file_size_mb`, … |
| `[goslop.taint]` | `enabled`, `show_paths` |
| `[goslop.bad_practices]` | `enabled`, `severity`, `severity_overrides` |

### Precedence (summary)

| Concern | Winner |
|---------|--------|
| only / skip | Config ∪ CLI |
| include tests | CLI `--include-tests` or config `exclude_tests=false` |
| cache off | CLI `--no-cache` or config `enabled=false` |
| cache dir | CLI `--cache-dir` else config |
| baseline off | CLI `--no-baseline` or config `enabled=false` |
| taint on/off | CLI `--taint` / `--no-taint` else config |
| profile / format / export / no-terminal | CLI only |

Starter template with comments: [`templates/goslop.toml`](../templates/goslop.toml).

---

## Example command gallery

```sh
# Default recommended gate
./bin/goslop .

# Full catalogue, JSON, no cache
./bin/goslop --profile all --format json --no-cache ./cmd

# Security + taint paths
./bin/goslop --profile security --taint-show-paths .

# Style-only (BP)
./bin/goslop --profile style .

# Export both surfaces for agent workflows
./bin/goslop --profile all --no-fail --export-context --export-chunks --no-cache .

# Product Makefile equivalent (see make-run.md)
./bin/goslop --profile all --no-fail --no-terminal \
  --export-context --export-chunks --no-cache /path/to/project

# SARIF for GitHub Code Scanning
./bin/goslop --profile recommended --no-fail --format sarif . > goslop.sarif

# Narrow run
./bin/goslop --only PERF-101,CWE-78 --format text ./internal
```

---

## Related docs

- [overview.md](./overview.md) - product features  
- [make-run.md](./make-run.md) - `make run` / `make oracle`  
- [reporting-formats.md](./reporting-formats.md) - text / JSON / SARIF  
- [export-context-and-chunks.md](./export-context-and-chunks.md) - exports  
- [go-recommended-pack.md](./go-recommended-pack.md) - pack contents  
