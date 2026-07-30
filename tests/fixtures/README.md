# Test fixtures (`.txt` text format)

Fixtures are **plain text** (`.txt`), not committed as `.go` / `.py` / `.rs`.

## Format

```text
# optional comment
lang: go | python | rust
file: sample.go          # optional; default: <stem>.<ext>
---
<language source body>
```

The `---` line separates header from source.

## Materialization

Before tests run:

1. Read `*.txt` under `tests/fixtures/<lang>/`
2. Write source to `target/goslop-fixtures/<lang>/<file>`
3. Run goslop on the generated paths

Integration tests call `goslop::fixture::materialize_fixture` automatically.

## Layout

| Language | Text fixture | Generated (gitignored) |
|----------|--------------|-------------------------|
| Go (stdlib / BP / CWE / PERF) | `go/**/*.txt` | temp dirs via integration harness |
| Python BP-PY | `python/bp/BP-PY-N-{vulnerable,safe}.txt` | temp dirs via `tests/integration/python` |
| Python CWE | `python/cwe/CWE-N-{vulnerable,safe}.txt` | temp dirs via `tests/integration/python` |
| Python samples | `python/sample.txt`, `python/safe.txt` | temp dirs |

### Integration test packages

| Package | Command | Coverage |
|---------|---------|----------|
| `tests/integration` | `make integration-go` / `go test ./tests/integration/` | **Go** fixtures only (`DefaultRegistry`) |
| `tests/integration/python` | `make integration-python` / `go test ./tests/integration/python/` | **Python** fixtures (`NewRegistryWithLanguages(LanguagePython)`) |

Do **not** add `.go`, `.py`, or `.rs` sources here — only `.txt` fixtures (header + `---` + body).
