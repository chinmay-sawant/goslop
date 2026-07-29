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
2. Write source to `target/codehound-fixtures/<lang>/<file>`
3. Run CodeHound on the generated paths

Integration tests call `codehound::fixture::materialize_fixture` automatically.

## Layout

| Language | Text fixture | Generated (gitignored) |
|----------|--------------|-------------------------|
| Go (stdlib) | `go/stdlib/*.txt` | `target/codehound-fixtures/go/*.pure.go` |
| Python | `python/sample.txt` | `target/codehound-fixtures/python/sample.py` |
| Rust | `rust/sample.txt` | `target/codehound-fixtures/rust/sample.rs` (plugin TBD) |

Do **not** add `.go`, `.py`, or `.rs` files here — only `.txt`.
