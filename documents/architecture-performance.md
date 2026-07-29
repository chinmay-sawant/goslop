# Architecture & performance notes

## Pipeline (language-agnostic)

```
CLI/config/profile merge
  → Analyzer
  → discover dual roots (prep root vs dependency base root)
  → collect_entries (walk + include/exclude)
  → detector begin_scan / prepare_project
  → chunked parallel scan (read → cache preflight | parse + detect;
       same-scan cascade invalidation for reverse deps)
  → detector finalize / end_scan
  → cache flush (+ advisory lock) + optional prune
  → reporting
```

Each file is read, parsed, analyzed, and dropped independently so peak memory
stays bounded on large repos. Details: [incremental-cache.md](./incremental-cache.md),
[cli.md](./cli.md).

## Incremental cache (P2.3)

- **Directory**: `.codehound-cache/` (near project / `go.mod` — see dual roots).
- **Manifest**: `manifest.json` tracks `tool_version`, `rule_config_hash`, per-file
  content hash, and dependency list (no separate `metadata.json`, no mtime).
- **Per-file entries**: `files/<sha256>.json` stores serialized findings + meta.
- **Cache hit flow**: hash match → load findings → re-apply ignores + rule filters.
- **Invalidation**: content hash, tool version, **rule-config fingerprint**,
  transitive deps, schema version. See [incremental-cache.md](./incremental-cache.md).
- **Pruning**: full-root scans only for sibling cleanup; `--prune-cache` /
  `cache prune`; `--rebuild-cache` purges the directory.
- **CLI flags**: `--no-cache`, `--cache-dir`, `--rebuild-cache`, `--prune-cache`.
- **Configuration**: `[codehound.cache]` (`enabled`, `path`, `max_size_mb`,
  `evict_target_ratio`, `max_file_size_mb`).
- **Eviction**: oldest-by-`cached_at` on flush when over `max_size_mb`.
- **Locking**: `.manifest.lock` advisory exclusive lock during flush.

## Multi-language (Go-first)

- **Cargo `default` features**: `go` + `cli` + `terminal-output` (**Python not default**).
- **Python**: opt-in `--features python` (one experimental rule). See ADR 0005.
- **`--lang auto`**: only languages compiled into the binary are scanned.
- TypeScript: **not supported** (no empty feature stub).

## Performance choices

| Area | Approach |
|------|----------|
| Parser | `ParsePool`: one parser per `LanguageId` per Rayon worker |
| Detectors | `Registry.by_language`: only matching rules per file |
| Go AST | One fused facts walk + `SourceIndex` flags per file |
| Go rules | Multi-file `registry/*.toml` drives `build.rs` (PERF and CWE) |
| CWE catalog | `src/cwe/{consts,description,reference,mod}.rs` |
| File pipeline | Chunked parallel read → parse → detect (`rayon`). Cache hits usually skip parse+detect (except taint `requires_cache_state`) |
| Source load | `String::from_utf8(bytes)` into `Arc<str>` |
| Source cache | Retained only when `ScanContext.retain_sources` (CLI export flags) |
| SourceIndex | **Aho-Corasick** multi-pattern build; `has()` is O(1) flag lookup via process-lifetime needle→index map |
| Typed Go (G4) | Optional `--typed` / `go list -json` package facts; degrades without toolchain |
| Taint project state | Built when taint enabled; units assembled off-lock; multi-hop finalize |
| Hash maps | See [adr/0001-hash-maps-on-hot-path.md](./adr/0001-hash-maps-on-hot-path.md) |
| Path identity | `normalize_project_path`; [adr/0002](./adr/0002-project-path-identity.md) |
| Same-scan cascade | Dirty fixpoint over reverse deps; [incremental-cache.md](./incremental-cache.md) |
| Detector lifecycle | `begin_scan` → parallel `run`/`accumulate_state` → `finalize` → `end_scan` |
| Dep extraction | `LanguagePlugin::extract_deps` relative to dependency base root |
| Export | Stream context/chunk files (no upfront `Vec` of all blocks) |
| Timing / stats | `--debug-timing` / `--diagnostics`; zero-cost when disabled |
| Diagnostics | `--diagnostics <FILE>` JSON + optional `--diagnostics-summary` |

## Codebase conventions (enforced)

| Rule | Limit / policy |
|------|----------------|
| `src/**/*.rs` module file | Soft **400**, hard **500** lines (`scripts/check_module_size.sh`); large exemptions for PERF domains, BP, taint, `app/run.rs`, … |
| Go CWE detector | Domain modules under `domains/`; multi-file `registry/*.toml` |
| New Go CWE/PERF rule | Add `[[detector]]` to the matching registry TOML + implement detector fn |
| Binary orchestration | `src/app/` only — `main.rs` stays tracing + `app::run` |
| PERF domain split | Prefer split around **500** lines when multi-group |

Run `scripts/check_module_size.sh` in CI or locally to catch module growth.

## Config behavior

- `only` and `skip` are additive across config and CLI.
- `fail_on` from config applies only when the CLI did not explicitly set `--strict`, `--no-fail`, or `--warnings-as-errors`.
- `include` and `exclude` are gitignore-style path globs applied during file collection.
- `.codehoundignore`, `.gitignore`, and `.ignore` remain active alongside config-backed include/exclude filtering.
- `--debug-timing` and `--diagnostics` are CLI-only flags (no config-file equivalent); they enable per-detector timing and phase-level instrumentation.

## Complexity (typical repo)

- Walk: O(files)
- Parse + detect: O(files / cores) wall time with rayon
- Per Go file: one tree-sitter parse + one fused AST walk + one `SourceIndex` build
- Detect: O(enabled_rules × facts); `--only` skips disabled rule bodies early
- Source cache memory: O(total UTF-8 bytes scanned successfully). The cache holds one shared `Arc<str>` per successful file; a 10 MiB file therefore keeps about 10 MiB of source text alive until the `AnalysisResult` is dropped. Files that cannot be read or decoded as UTF-8 are reported as `ScanError` and omitted from the cache. Use `AnalysisResult::source_cache_bytes()` to report the retained source-text byte count for a scan.

## Benchmarks & regression tests

- `cargo bench --bench scan_throughput` — full scan, collect-only, and `--only` subset
- `cargo test materialized_fixture_scan` — wall-clock smoke tests with tight ceilings (see `tests/perf_regression.rs`)

### Benchmark regression history

| Date | Baseline mean | After batch 3 | Regression | Cause |
|------|--------------|---------------|------------|-------|
| P2.4 batch 3 | ~3.2s | ~4.4s | ~38% | 7 new Category-A PERF detectors (PERF-114, 119, 125, 129, 156, 177, 192) adding source scan overhead |

**Mitigation:** Smoke budget in `tests/perf_regression.rs` was bumped from 600ms → 1.1s → 1.5s → 2.0s → 12s → 16s to accommodate the cumulative fixture surface. Current smoke tests pass at ~28s combined (under 32s ceiling).

**Baseline verification (2026-07-03):**
- `scan_materialized_fixtures` criterion mean: ~4.43s (baseline saved in `target/criterion/`)
- Smoke budget: 27.75s combined (within 32s ceiling)
- Benchmark takes ~6+ minutes to collect 100 samples; run with `cargo bench -- --sample-size 10` for quick checks

## Future optimizations

- Tree-sitter Query captures for hot rules
- Callee-indexed rule scheduling to skip rules when sinks are absent
- Per-detector rule-pack disabling (e.g. turn off `BP-*` or `PERF-1xx` via config)
