# CodeHound (Go port)

Go reimplementation of [CodeHound](https://github.com/chinmay-sawant/goslop) — a static analyzer for **Go PERF** hot paths, framework footguns, curated **CWE** heuristics, and **bad-practice** rules.

This repository ports the Rust product at `../codehound` while keeping:

- **Fixtures / heuristics text as-is** under `tests/fixtures/` and `ruleset/`
- **Detector logic rewritten in Go**

## Status

**Active port.** Phases **0–6** landed (MVP engine/CLI + full **239/239** PERF catalogue). Phases **7–12** are in parallel workstreams; **§12.4 product oracle is not complete**.

| Area | Status |
|------|--------|
| Engine / CLI / reporters (text, JSON, SARIF) | ✅ MVP |
| PERF detectors | ✅ 239/239 registered (heuristic) |
| CWE structural | 🌱 seed CWE-78 / CWE-89 only |
| Bad practices (BP) | ⏳ Phase 8 |
| Taint graph | ⏳ Phase 9 |
| Cache / baseline / ignore | ⏳ Phase 10 |
| Packs / maturity | ⏳ Phase 11 |
| CI + integration harness scaffolding | ✅ Phase 12 partial (this repo) |
| §12.4 export/scan oracle (915 findings) | ❌ blocked on 7–11 |

Canonical ledger:

- [`plans/port-phasewise-checklist.md`](./plans/port-phasewise-checklist.md)
- [`plans/architecture-go.md`](./plans/architecture-go.md)
- [`plans/parity-matrix.md`](./plans/parity-matrix.md)
- [`plans/phase-parallel/README.md`](./plans/phase-parallel/README.md)

## Requirements

- **Go 1.22+** (module currently declares a newer toolchain — see `go.mod`)
- **CGO enabled** — tree-sitter Go bindings need a C compiler (`build-essential` / Xcode CLT)
- Linux / macOS (Windows CGO not validated yet)

## Install / build

```sh
# From repo root
make build
# or:
CGO_ENABLED=1 go build -o bin/codehound ./cmd/codehound
```

Binary: `./bin/codehound`

Optional future multi-arch releases: see [`.goreleaser.stub.yml`](./.goreleaser.stub.yml) (not wired into CI yet). Prefer native runners with `CGO_ENABLED=1` over pure-Go cross-compile.

## Usage

```sh
# Scan current directory (default profile: recommended)
./bin/codehound .

# Profiles: recommended | perf | security | style | all
./bin/codehound --profile all --format json .
./bin/codehound --profile security --only CWE-78,CWE-89 ./path

# Machine formats
./bin/codehound --format json .
./bin/codehound --format sarif .

# Rule surface
./bin/codehound --list-rules
./bin/codehound --version

# Write a starter config
./bin/codehound init
```

Exit codes: `0` clean, `1` findings at fail policy, `2` usage/config, `3` internal.

## Develop / test

```sh
make test          # go test ./...
make integration   # seed fixture harness under tests/integration
make lint          # go vet + gofmt check
make ci            # lint + test + build (local CI parity)
```

GitHub Actions: [`.github/workflows/ci.yml`](./.github/workflows/ci.yml) runs `go vet`, `go test ./...`, and `go build` with CGO + `build-essential`.

Integration harness scaffolding: `tests/integration` materializes a small fixture seed (CWE-78/89, PERF-6) and asserts fire/silent behavior. Full fixture matrix and **§12.4** Rust oracle remain deferred.

## Layout

| Path | Role |
|------|------|
| `cmd/codehound` | CLI entry |
| `internal/` | Core, engine, lang/go detectors, reporting |
| `tests/fixtures` | Oracle fixtures (copied from Rust) |
| `tests/integration` | Materialized fixture harness (Phase 12 scaffolding) |
| `ruleset/` | Rule metadata JSON (copied) |
| `plans/` | Phase-wise port ledger |
| `.github/workflows/` | CI |

## License

MIT (same as upstream CodeHound).
