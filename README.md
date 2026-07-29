# CodeHound (Go port)

Go reimplementation of [CodeHound](https://github.com/chinmay-sawant/goslop) — a static analyzer for **Go PERF** hot paths, framework footguns, curated **CWE** heuristics, and **bad-practice** rules.

This repository ports the Rust product at `../codehound` while keeping:

- **Fixtures / heuristics text as-is** under `tests/fixtures/` and `ruleset/`
- **Detector logic rewritten in Go**

## Status

**Phases 0–12 landed.** Product surface, full detector catalogues, packs, CI, and the **§12.4** hard ship gate on `gopdfsuit` are locked on `main` (PR #16).

| Area | Status |
|------|--------|
| Engine / CLI / reporters (text, JSON, SARIF) | ✅ MVP |
| PERF detectors | ✅ **239/239** registered (heuristic) |
| CWE structural | ✅ **175/175** registered (+ taint-lite seeds) |
| Bad practices (BP) | ✅ Phase 8 catalogue + project-level rules |
| Taint graph | ✅ Phase 9 taint-core (CWE-22/78/79/89) |
| Cache / baseline / ignore | ✅ Phase 10 |
| Packs / maturity | ✅ Phase 11 (recommended / security / style / all) |
| CI + integration harness | ✅ Phase 12 |
| §12.4 export/scan oracle (`gopdfsuit`) | ✅ **915** findings; severity **10h/197i/312l/396m**; top-five exact; export **915+37**; skipped **383**; pure-Go wall **&lt;400ms** |
| Parse path | ✅ **Pure Go** (`go/parser` + `go/ast`); **no CGO** / no tree-sitter; `LanguagePlugin` seam for a second language |
| Config | ✅ `codehound.toml` discover/load + CLI merge (`only`/`skip` additive, include/exclude, cache, taint, …) |

Residual (severity-neutral / soft): partial per-site rule swaps (e.g. PERF-119↔128); skipped-file absolute may differ from Rust **383**; pure-FP museums suppressed on real project scans only. Details: [`plans/port-phasewise-checklist.md`](./plans/port-phasewise-checklist.md) §12.4.

Canonical ledger:

- [`plans/port-phasewise-checklist.md`](./plans/port-phasewise-checklist.md)
- [`plans/architecture-go.md`](./plans/architecture-go.md)
- [`plans/parity-matrix.md`](./plans/parity-matrix.md)
- [`plans/phase-parallel/README.md`](./plans/phase-parallel/README.md)

## Requirements

- **Go 1.22+** (module currently declares a newer toolchain — see `go.mod`)
- **Pure Go** — parsing uses `go/parser` + `go/ast` (**no CGO**, no tree-sitter)
- Linux / macOS / Windows

## Install / build

```sh
# From repo root (CGO_ENABLED=0 by default)
make build
# or:
CGO_ENABLED=0 go build -o bin/codehound ./cmd/codehound
```

Binary: `./bin/codehound`

Optional multi-arch releases: see [`.goreleaser.stub.yml`](./.goreleaser.stub.yml). Pure Go enables simple cross-compile without a C toolchain.

## Usage

```sh
# Scan current directory (default profile: recommended)
./bin/codehound .

# Profiles: recommended | perf | security | style | all
./bin/codehound --profile all --format json .
./bin/codehound --profile security ./path

# Export context + chunks (default dirs: scripts/findings/functions, scripts/chunks)
./bin/codehound --profile all --export-context --export-chunks --no-cache .

# Machine formats
./bin/codehound --format json .
./bin/codehound --format sarif .

# Rule surface
./bin/codehound --list-rules
./bin/codehound --explain PERF-6
./bin/codehound --version

# Write a starter config
./bin/codehound init

# §12.4 product oracle (SCAN_PATH defaults to gopdfsuit in Makefile)
make oracle
```

Exit codes: `0` clean, `1` findings at fail policy, `2` usage/config, `3` internal.

## Develop / test

```sh
make test          # go test ./...
make integration   # seed fixture harness under tests/integration
make lint          # go vet + gofmt check
make lint-all      # golangci-lint when configured
make ci            # lint + test + build (local CI parity)
make oracle        # gopdfsuit-scale §12.4 hard metrics (915 / sev / top-five)
```

GitHub Actions: [`.github/workflows/ci.yml`](./.github/workflows/ci.yml) runs `go vet`, `go test ./...`, and `go build` (pure Go / `CGO_ENABLED=0`).

Integration harness: `tests/integration` materializes a small fixture seed and asserts fire/silent behavior. Full fixture matrix is optional polish; **§12.4** hard counts are locked on the `gopdfsuit` corpus via `make oracle`.

## Documentation

Detailed user docs live under [`documents/`](./documents/):

| Doc | Contents |
|-----|----------|
| [`documents/README.md`](./documents/README.md) | Documentation index |
| [`documents/overview.md`](./documents/overview.md) | Features, profiles, cache, baseline, ignore |
| [`documents/cli-reference.md`](./documents/cli-reference.md) | All CLI flags, exit codes, config merge |
| [`documents/make-run.md`](./documents/make-run.md) | Product `make run` / `make oracle` workflow |
| [`documents/reporting-formats.md`](./documents/reporting-formats.md) | Text, JSON, **SARIF** examples |
| [`documents/export-context-and-chunks.md`](./documents/export-context-and-chunks.md) | `scripts/chunks` (delegate) + `scripts/findings/functions` (refs) |

Topic guides: [recommended pack](./documents/go-recommended-pack.md), [PERF notes](./documents/perf-rules.md), [taint](./documents/taint.md), [architecture](./documents/architecture-performance.md).

## Layout

| Path | Role |
|------|------|
| `cmd/codehound` | CLI entry |
| `internal/` | Core, engine, lang/go detectors, reporting |
| `documents/` | User-facing documentation |
| `tests/fixtures` | Oracle fixtures (copied from Rust) |
| `tests/integration` | Materialized fixture harness |
| `ruleset/` | Rule metadata JSON (copied) |
| `scripts/findings/functions` | Per-finding export (with `--export-context`) |
| `scripts/chunks` | Batched findings for agent delegation (`--export-chunks`) |
| `plans/` | Phase-wise port ledger |
| `.github/workflows/` | CI |

## License

MIT (same as upstream CodeHound).
