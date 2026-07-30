# Profiles and recommended pack

The default scan uses `--profile recommended`. Profiles select a useful rule
set and its default failure policy; they do not change the supported command
syntax. Select a profile with the CLI flag or use `goslop.toml` for the other
scan settings.

```sh
./bin/goslop .
./bin/goslop --profile perf .
./bin/goslop --profile security .
./bin/goslop --profile style .
./bin/goslop --profile all .
```

There is no `GOSLOP_PROFILE` environment variable and no `goslop rules`
subcommand. Use `--list-rules` and `--explain RULE_ID` for catalog discovery.

## Profile behavior

| Profile | Aliases | Rule selection | Taint | Bad practices | Default failure policy |
|---|---|---|---|---|---|
| `recommended` | `ci`, `default` | S-tier PERF plus taint-core CWE IDs | off | off | high / critical |
| `perf` | `performance` | S-tier and A-tier PERF | off | off | high / critical |
| `security` | `sec` | Security CWE pack | on | off | high / critical |
| `style` | `bp`, `bad-practices`, `bad_practices` | `BP-*` except default skips | off | on | none (advisory) |
| `all` | `full` | Full runtime catalog | off | on | medium and above |

`recommended`, `perf`, and `security` fail only on high and critical findings.
Use `fail_on = "medium"` in `goslop.toml` to make medium findings fail too, or
use `--no-fail` for an advisory scan.

## Recommended rule selection

The recommended pack includes these S-tier performance rules:

| Rule | Purpose |
|---|---|
| `PERF-1` | Regex compilation in a loop |
| `PERF-7` | `defer` in a loop |
| `PERF-50` | `regexp.MatchString` in a loop |
| `PERF-58` | Gin request body not closed |
| `PERF-71` | GORM N+1 query pattern |
| `PERF-101` | `http.Server` missing timeouts |
| `PERF-103` | HTTP response body not closed |
| `PERF-116` | Seeded hot-path performance detector |
| `PERF-189` | HTTP response body not drained before close |
| `PERF-190` | HTTP client missing a timeout |

It also allow-lists `CWE-22`, `CWE-78`, `CWE-79`, `CWE-89`, `CWE-90`, and
`CWE-91`. The experimental taint detector currently emits flows for the first
four IDs only; the other security rules remain structural catalog entries. See
[taint.md](./taint.md) for the precise taint scope.

## Security and style details

The security pack uses this exact allow-list:

```text
CWE-22, CWE-41, CWE-59, CWE-78, CWE-79, CWE-89, CWE-90, CWE-91, CWE-93
```

The `style` profile allows all `BP-*` rules but skips `BP-21`, `BP-28`, and
`BP-30` by default because they are intentionally opinionated. Add an explicit
`--only BP-28` (or configure `only`) when a skipped rule is wanted for a scan.

## CI examples

```sh
# Default high-signal gate.
./bin/goslop --profile recommended .

# Security triage with machine-readable results without blocking the job.
./bin/goslop --profile security --no-fail --format sarif . > goslop.sarif

# Fail a recommended scan on medium and higher findings.
./bin/goslop --profile recommended --format sarif . > goslop.sarif
```

For the final command, configure the threshold:

```toml
[goslop]
fail_on = "medium"
```

Use [reporting-formats.md](./reporting-formats.md) for SARIF upload details and
[suppressions-and-baselines.md](./suppressions-and-baselines.md) when adopting
the tool in an existing repository.

## Discovering the live catalog

```sh
./bin/goslop --list-rules
./bin/goslop --explain PERF-101
./bin/goslop --explain CWE-89
```

The list labels rules as `production`, `experimental`, `fixture-only`, or
`quarantined`. The labels help with rollout decisions; they do not add extra
CLI modes. See [rule-catalog-and-maturity.md](./rule-catalog-and-maturity.md).
