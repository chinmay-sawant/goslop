# Architecture and performance notes

## Scan pipeline

```text
CLI and goslop.toml resolution
  -> ScanPlan and ScanScope
  -> file collection
  -> per-scan detector session
  -> bounded parallel read, parse, cache, and detection
  -> finalize stateful detectors
  -> baseline/filtering, cache prune and flush
  -> text, JSON, SARIF, or context export
```

`Analyzer.AnalyzePaths` owns scanning and cache persistence. Every call creates
a private detector session through `LanguagePlugin.NewDetectors`, so detector
lifecycle state cannot cross concurrent scans. The registry retains a separate
immutable catalogue for rule listing and file-language resolution.

## Cache

- **Directory:** `.goslop-cache/` by default, or `--cache-dir` / `[goslop.cache].path`.
- **Key:** cache entries use a project-relative, normalized source path plus a
  content hash and rule-configuration fingerprint.
- **Hit behavior:** ordinary detectors reuse stored findings. Detectors that
  require project state, including taint analysis, receive a parsed unit and
  `AccumulateState` on a hit without rerunning normal detection.
- **Persistence:** scan-time prune and flush failures fail the scan. The CLI
  then reconciles orphaned cache files.
- **Sizing:** `[goslop.cache]` supports `max_size_mb`, `max_file_size_mb`, and
  `evict_target_ratio` (0.10 through 0.99).

## Language seam

Go is the only shipped language. Its plugin uses the standard library parser
and keeps Go AST details out of the engine. A future language plugin must:

1. implement `LanguagePlugin` for its extensions and parser;
2. return immutable catalogue detectors from `Detectors`;
3. return fresh session-local instances from `NewDetectors`;
4. keep language-specific AST values behind `ParsedUnit.Tree`.

This preserves the engine's language-neutral interface while making stateful
detectors safe under bounded file parallelism.

## Performance choices

| Area | Current approach |
|---|---|
| File work | `errgroup` with a worker limit from `ScanContext` / `GOMAXPROCS` |
| Go parsing | `go/parser` and shared Go AST facts per parsed unit |
| Detector dispatch | language and enabled-rule filtering before rule bodies run |
| Cache | content-addressed findings with bounded file eligibility and eviction |
| Taint | project facts accumulated across cached and fresh files before finalization |
| Export | owned numeric context files and `Chunk_*.txt` files only; unrelated caller files are preserved |

## Validation

The standard local validation gate is:

```sh
gofmt -w <changed Go files>
make test
```

For concurrency-sensitive changes, also run `go test -race ./...`. Benchmark
claims should use a representative external Go project and record separate
cold-cache and warm-cache results.
