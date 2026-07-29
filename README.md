# CodeHound (Go port)

Go reimplementation of [CodeHound](https://github.com/) — a static analyzer for **Go PERF** hot paths, framework footguns, curated **CWE** heuristics, and **bad-practice** rules.

This repository ports the Rust product at `../codehound` while keeping:

- **Fixtures / heuristics text as-is** under `tests/fixtures/` and `ruleset/`
- **Detector logic rewritten in Go**

## Status

**Active port.** See the canonical checklist:

- [`plans/port-phasewise-checklist.md`](./plans/port-phasewise-checklist.md)
- [`plans/architecture-go.md`](./plans/architecture-go.md)
- [`plans/parity-matrix.md`](./plans/parity-matrix.md)

## Build

```sh
go build -o bin/codehound ./cmd/codehound
make test
```

Requires **Go 1.22+**.

## Usage (as phases land)

```sh
./bin/codehound .
./bin/codehound --profile recommended --format json ./...
./bin/codehound --list-rules
```

## Layout

| Path | Role |
|------|------|
| `cmd/codehound` | CLI entry |
| `internal/` | Core, engine, lang/go detectors, reporting |
| `tests/fixtures` | Oracle fixtures (copied from Rust) |
| `ruleset/` | Rule metadata JSON (copied) |
| `plans/` | Phase-wise port ledger |

## License

MIT (same as upstream CodeHound).
