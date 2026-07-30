## Context

Export context (`scripts/findings/functions/N.txt`) and chunks (`scripts/chunks/Chunk_*.txt`) previously showed a **~4-line window** around the hit (`line-2 … line+1`) via `lineWindow` in `internal/export/export.go`. That is often too little for agents and humans to understand the issue; the directory is even named `functions`, but the body was not the enclosing function.

Product need: make **whole enclosing function** the default context for both export surfaces, with a config flag to fall back to the nearby-line window.

**Issue:** https://github.com/chinmay-sawant/goslop/issues/29

## Scope (in)

1. Config flag under `[goslop.export]`:
   - `whole_function = true` (**default** when unset)
   - `whole_function = false` restores nearby ~4-line window
2. Shared context builder used by **both** per-finding refs and chunks (`formatFindingBlock` / `findingContextLines`)
3. Resolve enclosing function with pure-Go `go/parser` + `go/ast`:
   - Prefer **outermost `FuncDecl`** containing the hit (named method/function)
   - Else outermost `FuncLit` (standalone closures)
   - Fall back to line window when outside any function / non-Go / parse failure
4. Schema (`goslop.schema.json`), starter templates (`templates/goslop.toml`, `init_template.toml`), and docs (`documents/export-context-and-chunks.md`)
5. Unit tests covering whole-function vs window, outer FuncDecl over defer closure, chunks

## Out of scope

- Changing stdout reporters (text / JSON / SARIF snippets)
- Changing finding fingerprints or detector emit sites
- Multi-language whole-function extraction (non-Go → keep line window)
- Committing regenerated §12.4 corpora (content shape changes; file counts stay 915+37)

## Success criteria

- [x] Default (no config / `whole_function` unset): export Context includes the full enclosing function with `>` on the hit line
- [x] `whole_function = false`: Context uses the previous nearby-line window
- [x] Chunks and function refs use the same context builder
- [x] Hits inside `defer func() { ... }` export the **outer named function**, not only the tiny closure
- [x] Config rejected for unknown fields still works; schema documents the new table
- [x] Tests pass for both modes

## Plan

- [x] Config + schema + merge wiring (`[goslop.export] whole_function`)
- [x] `export.Options.WholeFunction` + function-span extraction (FuncDecl preferred)
- [x] Docs + tests
- [x] Rebuilt binary + re-export on gopdfsuit verified (median context ~70 lines; large funcs hundreds of lines)

## References

- Implementation: `internal/export/export.go` (`functionWindow`, `enclosingFunctionLines`, `lineWindow`)
- Docs: `documents/export-context-and-chunks.md`
- Config: `internal/config/config.go`, `templates/goslop.toml`
- Issue: #29
