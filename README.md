# goslop

**goslop** is a **static analysis tool (SAT)** with a **multi-language engine** (`LanguagePlugin`). **Go** is the production default today: it inspects Go source with pure Go parsing (`go/parser` + `go/ast`, no CGO) and reports issues across three main families. **Python** is reserved / WIP (fixtures + `LanguagePython`; no full detector catalogue yet)—see epic [#39](https://github.com/chinmay-sawant/goslop/issues/39) and [`plans/v0.0.2/python-support.md`](./plans/v0.0.2/python-support.md).

| Family | What it finds (Go catalogues) | Examples |
|--------|-------------------------------|----------|
| **Performance (`PERF-*`)** | Hot-path and runtime footguns | Allocations in loops, missing HTTP timeouts, N+1 patterns, framework misuse |
| **Bad practices (`BP-*`)** | Style, hygiene, and project-level issues | Error handling, server shutdown, go.mod hygiene, noisy API patterns |
| **CWE (`CWE-*`)** | Security heuristics (structural + optional taint) | Path traversal, command injection, XSS, SQL injection (CWE-22/78/79/89 and more) |

Use it as a local linter, a CI gate, or a triage pipeline that exports findings for human or agent review (JSON, SARIF, and on-disk context chunks).

---

## Highlights

- **Three Go detector catalogues** - **239** Go PERF, **175** Go CWE, plus a full Go bad-practice catalogue (`BP-*`)
- **Multi-language engine** - language-agnostic `LanguagePlugin` seam; Go shipped; Python WIP (#39)
- **Product profiles** - `recommended`, `perf`, `security`, `style`, `all` (curated packs + fail policies)
- **Reporters** - human **text**, machine **JSON**, and **SARIF 2.1.0** (GitHub Code Scanning-ready)
- **Optional taint** - experimental inter-procedural graph for high-signal injection CWEs
- **Incremental cache** - `.goslop-cache/` for fast re-scans
- **Baseline & ignore** - ship with known debt; suppress with `// goslop-ignore`
- **Export for agents** - per-finding refs under `scripts/findings/functions/`, batched **chunks** under `scripts/chunks/` (Context defaults to the **whole enclosing function**; `[goslop.export] whole_function`)
- **Pure Go binary** - `CGO_ENABLED=0` by default; easy cross-compile

---

## Requirements

- **Go 1.26.4** (the version declared in `go.mod`)
- Linux / macOS / Windows
- No C toolchain required for the default pure-Go build

---

## Install / build

```sh
# From repo root (CGO_ENABLED=0 by default via Makefile)
make build

# Or:
CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop
```

Binary: `./bin/goslop`

Optional multi-arch packaging: [`.goreleaser.stub.yml`](./.goreleaser.stub.yml).

---

## Quick start

```sh
# Scan current directory (default profile: recommended)
./bin/goslop .

# Performance-focused pack
./bin/goslop --profile perf .

# Security pack (taint on by default)
./bin/goslop --profile security ./cmd

# Bad practices / style
./bin/goslop --profile style .

# Full catalogue
./bin/goslop --profile all .

# Starter config
./bin/goslop init
```

### Machine output

```sh
./bin/goslop --format json .
./bin/goslop --format sarif . > goslop.sarif
```

### Export findings for review / agent delegation

```sh
# Per-finding refs → scripts/findings/functions/N.txt
# Combined batches → scripts/chunks/Chunk_START_END.txt (default 25 findings each)
./bin/goslop --profile all --export-context --export-chunks --no-cache .
```

- **Chunks** = combined findings for **delegating** work to agents  
- **Functions** = individual finding **refs** for single-issue deep dives  

Product-style summary scan (build + export defaults):

```sh
make run SCAN_PATH=./your/go/project
```

**Note:** Scanning **this** repository under `--profile all` yields a large finding count. That is expected for a SAT that ships **pattern needles**, **detector source**, and **intentional fixtures**—the same class of noise you get when tools like **Semgrep** (or other SATs) analyze their own rule packs and test corpora. Use a real app (or the `gopdfsuit` baseline via default `make run`) to judge product signal. See [documents/overview.md](./documents/overview.md#scanning-the-goslop-repo-itself-high-finding-counts).

### Rule discovery

```sh
./bin/goslop --list-rules
./bin/goslop --explain PERF-6
./bin/goslop --explain CWE-89
./bin/goslop --version
```

### Exit codes

| Code | Meaning |
|------|---------|
| **0** | Clean, or `--no-fail` |
| **1** | Findings violate the active fail policy |
| **2** | Usage / configuration error |
| **3** | Internal error |

---

## What goslop analyzes

Shipped detectors today target **Go**. Catalogue counts below are the **Go** surface (not multi-language totals).

### Performance (`PERF-*`)

Heuristic detectors for expensive patterns on request paths and loops: regex compile-in-loop, `defer` in loops, HTTP server timeouts, body close issues, framework-specific hot paths, and more. **239** Go rules registered.

Human notes (partial): [`documents/perf-rules.md`](./documents/perf-rules.md).

### Bad practices (`BP-*`)

Style and engineering hygiene: missing package docs, error handling habits, testing smells, public HTTP without rate limiting, graceful shutdown, go.mod / dependency hygiene, and related project-level checks. Enabled under **`style`** and **`all`** profiles (and when BP is turned on in config).

### CWE / security (`CWE-*`)

**175** Go structural CWE heuristics. For injection-class issues, enable experimental **taint** tracking:

```sh
./bin/goslop --taint --taint-depth 3 --taint-show-paths .
# or
./bin/goslop --profile security .
```

Taint model and limits: [`documents/taint.md`](./documents/taint.md).

---

## Profiles (packs)

| Profile | Aliases | Focus | Default fail |
|---------|---------|--------|--------------|
| `recommended` | `ci`, `default` | High-signal PERF + core CWE allow-list | high / critical |
| `perf` | `performance` | Broader PERF tiers | high / critical |
| `security` | `sec` | Security CWEs + **taint on** | high / critical |
| `style` | `bp`, `bad-practices` | `BP-*` catalogue | none (advisory) |
| `all` | `full` | Full PERF + CWE + BP | medium and above |

Details and exact allow-lists: [`documents/go-recommended-pack.md`](./documents/go-recommended-pack.md).  
Full flag list: [`documents/cli-reference.md`](./documents/cli-reference.md).

---

## Configuration

```sh
./bin/goslop init   # writes goslop.toml
```

Optional `goslop.toml` supports rule filters (`only` / `skip`), `fail_on`, include/exclude globs, cache, baseline, taint, and bad-practice severity overrides.

- Template: [`templates/goslop.toml`](./templates/goslop.toml)  
- Schema: [`goslop.schema.json`](./goslop.schema.json)  
- Merge rules: [`documents/cli-reference.md`](./documents/cli-reference.md#config-file-and-cli-merge)

Suppress inline with:

```go
// goslop-ignore: PERF-101
srv := &http.Server{Addr: ":8080"}
```

---

## Reporting formats

| Format | Flag | Use |
|--------|------|-----|
| **text** | `--format text` (default) | Terminal / CI logs |
| **json** | `--format json` | Scripts and automation |
| **sarif** | `--format sarif` | GitHub Code Scanning, IDE SARIF viewers |

Scan **summary** always goes to **stderr**; finding payloads go to **stdout** so JSON/SARIF stay pipe-clean.

Examples and GitHub upload notes: [`documents/reporting-formats.md`](./documents/reporting-formats.md).

---

## Develop / test

```sh
make test          # go test ./...
make integration   # fixture harness under tests/integration
make lint          # go vet + gofmt check
make lint-all      # golangci-lint (when configured)
make ci            # lint + test + build
make run           # product summary scan + optional exports
make reference-metrics        # large-corpus metrics gate (see Makefile SCAN_PATH)
make help          # list targets
```

CI: [`.github/workflows/ci.yml`](./.github/workflows/ci.yml) runs `go vet`, `go test ./...`, and `go build` with `CGO_ENABLED=0`.

---

## Documentation

Detailed guides live under [`documents/`](./documents/):

| Document | Contents |
|----------|----------|
| [`documents/README.md`](./documents/README.md) | Documentation index |
| [`documents/overview.md`](./documents/overview.md) | Features, profiles, cache, baseline, ignore |
| [`documents/cli-reference.md`](./documents/cli-reference.md) | Every CLI flag, exit codes, config merge |
| [`documents/make-run.md`](./documents/make-run.md) | Product `make run` / `make reference-metrics` workflow |
| [`documents/reporting-formats.md`](./documents/reporting-formats.md) | Text, JSON, and **SARIF** with examples |
| [`documents/export-context-and-chunks.md`](./documents/export-context-and-chunks.md) | Function refs vs chunks for agent delegation |
| [`documents/go-recommended-pack.md`](./documents/go-recommended-pack.md) | Recommended pack and profile tables |
| [`documents/perf-rules.md`](./documents/perf-rules.md) | PERF rule notes |
| [`documents/taint.md`](./documents/taint.md) | Taint engine |
| [`documents/suppressions-and-baselines.md`](./documents/suppressions-and-baselines.md) | Inline ignores and baseline rollout |
| [`documents/rule-catalog-and-maturity.md`](./documents/rule-catalog-and-maturity.md) | Live catalog discovery and maturity labels |
| [`documents/development.md`](./documents/development.md) | Build, validation, CI, scan, and benchmark commands |
| [`documents/architecture-performance.md`](./documents/architecture-performance.md) | Engine pipeline and performance design |

---

## Repository layout

| Path | Role |
|------|------|
| `cmd/goslop` | CLI entrypoint |
| `internal/` | Core engine, detectors, reporting, export |
| `documents/` | User-facing documentation |
| `ruleset/` | Rule metadata (PERF / CWE chunks, BP catalogue) |
| `templates/` | Starter `goslop.toml` |
| `tests/fixtures` | Detector fixtures |
| `tests/integration` | Integration harness |
| `scripts/findings/functions` | Per-finding export (`--export-context`) |
| `scripts/chunks` | Batched findings for delegation (`--export-chunks`) |
| `.github/workflows/` | CI |

---
