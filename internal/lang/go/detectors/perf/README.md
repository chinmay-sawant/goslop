# Go PERF detectors

Performance heuristics for Go hot paths and stdlib misuse.

## Seed (implemented)

| Rule ID | Domain | Notes |
|---------|--------|-------|
| `PERF-116` | general_perf / strings | `strings.Index(s, sub) != -1` → prefer `strings.Contains` |

## Registry inventory (remaining)

Source TOMLs under `registry/`:

| Registry file | Rows | Status |
|---------------|-----:|--------|
| `registry.general_perf.toml` | 164 | seed: PERF-116 |
| `registry.data_access.toml` | 20 | pending |
| `registry.gin_framework.toml` | 20 | pending |
| `registry.request_path.toml` | 9 | pending |
| `registry.protocols.toml` | 10 | pending |
| `registry.parsing_in_loops.toml` | 8 | pending |
| `registry.loop_allocations.toml` | 8 | pending |
| **Total** | **239** | |

Port Rust domains under `codehound/src/lang/go/detectors/perf/domains/` into sibling packages; keep `PERF-N` ids stable.
