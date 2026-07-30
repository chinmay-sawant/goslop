# Parity Matrix - Rust → Go

> **Parent:** `plans/port-phasewise-checklist.md`
> **Status:** living map

## Package map

| Rust (`goslop/src`) | Go (`goslop-go/internal`) | Phase | Status |
|------------------------|------------------------------|-------|--------|
| `main.rs` / `bin/` | `cmd/goslop` | 0 | [ ] |
| `app/` | `app/` | 4 | [ ] |
| `cli/` | `cli/` | 4 | [ ] |
| `core/` | `core/` | 1 | [ ] |
| `rules/` | `rules/` | 1 | [ ] |
| `engine/` | `engine/` | 2-3, 10 | [ ] |
| `ast/` | `ast/` | 2 | [ ] |
| `lang/go/` | `lang/go/` | 2, 6-9 | [ ] |
| `lang/python/` | - | v0.0.2 epic #39 | [~] foundation WIP; catalogues → #43 `ruleset/python/` |
| `reporting/` | `reporting/` | 5 | [ ] |
| `export/` | `export/` | 5 | [ ] |
| `fixture/` | `fixture/` | 2 | [ ] |
| `cwe/` | `cwe/` | 1 | [ ] |
| `error.rs` | `internal/errcode` or `engine` errors | 1 | [ ] |

## Rule pack counts (from Rust registries / product docs)

| Pack | Approximate count | Go port target | Status |
|------|------------------:|----------------|--------|
| PERF | ~239 | all domains | [x] Phase 6 - **239/239** registered (heuristic batch ports) |
| CWE structural | ~175 | all domains | [x] Phase 7 - **175/175** (integration) |
| BP | ~135 (rules files ~39 modules) | all | [x] Phase 8 (127 detectors registered, heuristic) |
| Taint CWE-22/78/79/89 | experimental | port graph + rules | [x] Phase 9 (integration) |

Update counts from `goslop --list-rules` when Rust binary available; registries under `internal/lang/go/detectors/*/registry/` are source for wiring.

## Fixture surface

| Path | Files (copied) | Notes |
|------|---------------:|-------|
| `tests/fixtures/go/**` | bulk of 1746 | primary expectations |
| `tests/fixtures/python/**` | small | [~] still skip for detectors; catalogue seeds in `ruleset/python/` (#43) |
| `tests/fixtures/manifest.toml` | 1 | materializer must read |

## CLI flag parity (minimum)

| Flag / concept | Phase | Status |
|----------------|-------|--------|
| path args / `.` default | 4 | [ ] |
| `--profile recommended\|security\|all` | 4, 11 | [ ] |
| `--only` / `--skip` | 4 | [ ] |
| `--format text\|json\|sarif` | 5 | [ ] |
| `--list-rules` / `--explain` | 4 | [ ] |
| `--include-tests` | 3 | [ ] |
| `--no-cache` / `--cache-dir` / `--rebuild-cache` / `--prune-cache` | 10 | [x] |
| baseline / ignore | 10 | [x] |
| `--taint` / taint depth | 9 | [x] |
| `--typed` | 9+ | [~] after core |
| `init` | 4 | [ ] |
| `--export-context` / `--export-chunks` | 5, **12.4** | [ ] |
| `--no-cache` full re-analysis | 10, **12.4** | [ ] |

## Final validation baseline (Phase 12.4)

Rust reference (do not change without updating the checklist):

```bash
make run RUN_ARGS="--export-context --export-chunks --no-cache"
```

| Metric | Value |
|--------|------:|
| Files scanned | 78 |
| Lines | 28120 |
| Cache | 0 hits / 78 misses |
| Skipped | 383 |
| Findings | **915** |
| Severity | 10 high, 197 info, 312 low, 396 medium |
| Top rules | BP-1×181, PERF-6×94, PERF-32×59, BP-5×50, PERF-230×44 |
| Context exports | 915 → `scripts/findings/functions` |
| Chunk exports | 37 → `scripts/chunks` |
| Wall (ref host) | ~479.5 ms (soft) |
