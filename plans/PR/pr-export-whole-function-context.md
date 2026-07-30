## Summary

Export context and chunks now include the **full enclosing Go function** by default (configurable via `[goslop.export] whole_function`), instead of a ~4-line nearby window. Hits inside tiny closures (e.g. `defer func()`) expand to the outer named `FuncDecl`.

## Motivation / context

- Agents and humans triaging `scripts/findings/functions/` and `scripts/chunks/` need the whole function body, not 3–4 surrounding lines.
- Plans: `plans/PR/issue-export-whole-function-body.md`
- Issues: see **Related issues**

## Changes

### Export context

- `export.Options.WholeFunction *bool` — `nil` / unset → **default true**
- Whole-function span via pure-Go `go/parser` + `go/ast`
- Prefer outermost **`FuncDecl`** over inner **`FuncLit`** so defer closures still export the surrounding method
- Fall back to `line-2…line+1` when no enclosing function / non-Go / parse failure
- Shared path for both `--export-context` and `--export-chunks`

### Config / schema / templates

- New `[goslop.export] whole_function` in `goslop.toml` (default true when omitted)
- `goslop.schema.json`, `templates/goslop.toml`, `init_template.toml`

### Docs / tests

- `documents/export-context-and-chunks.md` updated for whole-function default
- Unit tests: whole vs window, outer FuncDecl over closure, chunk export, config merge

## Impact

| Area | Impact |
|------|--------|
| **Performance** | Extra pure-Go parse per unique source file during export only (cached per file); scan path unchanged |
| **Memory** | Larger export text files when functions are long |
| **Behavior / correctness** | Context shape changes (full function); finding counts / fingerprints / §12.4 metrics unchanged |
| **API / CLI** | No new CLI flags; config-only toggle |
| **Dependencies** | None (stdlib `go/parser`) |
| **Binary size / build time** | Negligible |

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| Export Context text shape | Full function by default. Set `whole_function = false` under `[goslop.export]` for the old ~4-line window |

## Test plan

- [x] `go test ./internal/export/ ./internal/config/`
- [x] `CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop`
- [x] Re-export gopdfsuit: **915** context + **37** chunks; context medians far above 4 lines
- [x] Spot-check: BP-5 inside `defer func()` exports full outer handler

### Commands

```sh
make build
go test ./internal/export/ ./internal/config/
make run   # rebuilds + re-exports with default whole_function=true
```

### Config

```toml
[goslop.export]
whole_function = true   # default; omit for same effect
# whole_function = false  # nearby ~4-line window
```

## Screenshots / sample output

After rebuild + re-export (gopdfsuit), finding with hit inside defer:

```text
Context:
        423: func handlerSplitPDF(c *gin.Context) {
        ...
        430: defer func() {
    >   431: _ = pdfFile.Close()
        432: }()
        ...
        # full remainder of handlerSplitPDF
```

Previously this was only the 3-line closure / 4-line window.

## Related issues

Closes #29
