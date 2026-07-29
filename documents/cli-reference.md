# CLI reference

Complete command-line reference for the Go CodeHound binary (`./bin/codehound`).

```text
Usage: codehound [flags] [PATH...]
       codehound init
```

- Flags accept both `-flag` and `--flag` forms (stdlib `flag` package).
- If no `PATH` is given, the scan root defaults to **`.`**.
- Scan **summary** always prints to **stderr**. Finding payloads go to **stdout** (except `--no-terminal` with text format).

---

## Subcommands and meta modes

| Mode | Invocation | Behavior |
|------|------------|----------|
| **Scan** | `codehound [flags] [PATH...]` | Analyze paths (default) |
| **`init`** | `codehound init` | Write starter `codehound.toml` in cwd; **no flags parsed after `init`** |
| **Help** | `-h` / `-help` | Print usage; exit 0 |
| **Version** | `--version` | Print version (e.g. `0.1.0-dev`); exit 0 |
| **List rules** | `--list-rules` | Print registered rules; exit 0 |
| **Explain rule** | `--explain RULE_ID` | Print catalogue metadata; exit 0 (unknown → exit 2) |
| **Prune cache** | `--prune-cache` | Prune stale cache for PATHS and exit |

### `init`

```sh
./bin/codehound init
# wrote starter codehound.toml to /abs/path/codehound.toml
```

- Fails with exit **2** if `codehound.toml` already exists.
- Content matches [`templates/codehound.toml`](../templates/codehound.toml).

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
./bin/codehound --only "CWE-78,PERF-6,CWE-89" .
./bin/codehound --profile style --only BP-28 .
./bin/codehound --skip PERF-1,PERF-7 .
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
./bin/codehound --format text .
./bin/codehound --format json . > findings.json
./bin/codehound --format sarif . > codehound.sarif
./bin/codehound --no-terminal .                    # summary only
./bin/codehound --no-terminal --format json .      # summary on stderr + JSON on stdout
```

Full format schemas and SARIF examples: [reporting-formats.md](./reporting-formats.md).

---

### Rule discovery

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--list-rules` | bool | `false` | List registered rules (`[maturity]  id  title`) and exit |
| `--explain` | string | _(empty)_ | Print catalogue metadata for one rule id and exit |

```sh
./bin/codehound --list-rules
./bin/codehound --explain PERF-6
./bin/codehound --explain CWE-89
```

---

### Tests, paths, and config

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--include-tests` | bool | `false` | Include test files (`*_test.*`) in analysis |
| `--config` | string | _(discover)_ | Path to `codehound.toml` (default: walk upward from first PATH) |
| `--version` | bool | `false` | Print version and exit |

Positional **PATH…** arguments are scan roots (files or directories).

---

### Cache

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--no-cache` | bool | `false` | Disable the incremental analysis cache |
| `--cache-dir` | string | `.codehound-cache` | Incremental cache directory |
| `--rebuild-cache` | bool | `false` | Purge the cache directory before scanning |
| `--prune-cache` | bool | `false` | Prune stale cache entries for PATHS and exit |

```sh
./bin/codehound --no-cache .
./bin/codehound --rebuild-cache .
./bin/codehound --cache-dir /tmp/ch-cache .
./bin/codehound --prune-cache .
```

---

### Baseline and ignore visibility

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--no-baseline` | bool | `false` | Ignore any existing `.codehound-baseline.json` |
| `--baseline-file` | string | _(discover)_ | Explicit baseline path |
| `--show-ignored` | bool | `false` | Report findings suppressed by `codehound-ignore` |
| `--show-baselined` | bool | `false` | Report findings present in the baseline |

---

### Taint

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--taint` | bool | `false` | Enable experimental taint tracking (CWE-22/78/79/89) |
| `--no-taint` | bool | `false` | Disable taint even under `security` profile |
| `--taint-depth` | int | `0` | Inter-procedural hops **1–4**; `0` = profile default (security → 4, else 1) |
| `--taint-show-paths` | bool | `false` | Attach taint hop evidence to findings |

```sh
./bin/codehound --taint --taint-depth 3 --taint-show-paths .
./bin/codehound --profile security .
./bin/codehound --profile security --no-taint .
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
./bin/codehound --profile all --export-context --export-chunks --no-cache .
./bin/codehound --export-chunks --chunk-size 50 --chunks-dir /tmp/chunks .
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
2. Else walk upward from the first scan path (or `.`) for `codehound.toml`  
3. Missing file is fine (CLI-only defaults)

### Common TOML keys (`[codehound]`)

| Key | Effect |
|-----|--------|
| `only` / `skip` | Union with CLI `--only` / `--skip` |
| `fail_on` | Sets fail policy unless `--no-fail` |
| `include` / `exclude` | Walk globs (config only; no CLI flag) |
| `exclude_tests` | Default true; set `false` to include tests (CLI `--include-tests` also includes) |
| `[codehound.baseline]` | `enabled`, `path` |
| `[codehound.cache]` | `enabled`, `path`, `max_size_mb`, `max_file_size_mb`, … |
| `[codehound.taint]` | `enabled`, `show_paths` |
| `[codehound.bad_practices]` | `enabled`, `severity`, `severity_overrides` |

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

Starter template with comments: [`templates/codehound.toml`](../templates/codehound.toml).

---

## Example command gallery

```sh
# Default recommended gate
./bin/codehound .

# Full catalogue, JSON, no cache
./bin/codehound --profile all --format json --no-cache ./cmd

# Security + taint paths
./bin/codehound --profile security --taint-show-paths .

# Style-only (BP)
./bin/codehound --profile style .

# Export both surfaces for agent workflows
./bin/codehound --profile all --no-fail --export-context --export-chunks --no-cache .

# Product Makefile equivalent (see make-run.md)
./bin/codehound --profile all --no-fail --no-terminal \
  --export-context --export-chunks --no-cache /path/to/project

# SARIF for GitHub Code Scanning
./bin/codehound --profile recommended --no-fail --format sarif . > codehound.sarif

# Narrow run
./bin/codehound --only PERF-101,CWE-78 --format text ./internal
```

---

## Related docs

- [overview.md](./overview.md) — product features  
- [make-run.md](./make-run.md) — `make run` / `make oracle`  
- [reporting-formats.md](./reporting-formats.md) — text / JSON / SARIF  
- [export-context-and-chunks.md](./export-context-and-chunks.md) — exports  
- [go-recommended-pack.md](./go-recommended-pack.md) — pack contents  
