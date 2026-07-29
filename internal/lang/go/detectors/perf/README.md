# Go PERF detectors

Performance heuristics for Go hot paths and stdlib misuse.

## Architecture

- **`GoPerfScan`** (`scan.go`) — single `core.Detector`; builds facts once, dispatches enabled rules
- **`RegisterRule`** (`register.go`) — `init()`-based registration so domain batches can land independently
- **`BuildFacts`** (`facts.go`) — tree-sitter fused walk: calls, assignments, conversions, defer/go/for, var kinds
- **`common.go`** — hot-path / handler-shaped / loop helpers

## Coverage

| Source | Rules |
|--------|------:|
| Registry TOMLs | **239** |
| Registered in Go | **239** (100%) |

### Seed (Phase 6a)

`PERF-1`…`8`, `PERF-32`, `PERF-50`, `PERF-116`, `PERF-230` — `rules_loop.go` / `rules_hot.go` / `seed_register.go`

### Batch ports (Phase 6b)

| Batch | File prefix | Rule range (approx) | Count |
|------:|-------------|---------------------|------:|
| 1 | `*_batch1.go` | PERF-9…60 (ex 32, 50) | 50 |
| 2 | `*_batch2.go` | PERF-61…111 | 50 |
| 3 | `*_batch3.go` | PERF-112…163 (ex 116) | 50 |
| 4 | `*_batch4.go` | PERF-164…214 | 50 |
| 5 | `*_batch5.go` | PERF-215…242 (ex 230) | 27 |

## Registry domains

| Registry file | Status |
|---------------|--------|
| `registry.loop_allocations.toml` | done |
| `registry.general_perf.toml` | done (heuristic ports) |
| `registry.data_access.toml` | done |
| `registry.gin_framework.toml` | done |
| `registry.request_path.toml` | done |
| `registry.protocols.toml` | done |
| `registry.parsing_in_loops.toml` | done |

## Quality note

Batch ports are **behavioral heuristics** aligned with Rust detectors and fixtures, not a mechanical transpile. Some rules use source-token guards where shared `perfNeedles` is incomplete (shared `facts.go` intentionally frozen during parallel landings). Tighten against the Rust oracle in follow-ups; full §12.4 gate still requires BP + packs + export.

## Proof

```bash
go test ./internal/lang/go/detectors/perf/ ./...
go build -o bin/codehound ./cmd/codehound
./bin/codehound --list-rules | wc -l   # expect 239 PERF + CWE seeds
```
