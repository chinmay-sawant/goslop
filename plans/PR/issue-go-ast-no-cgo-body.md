## Context

goslop currently parses Go sources with **tree-sitter** (`internal/lang/go/tsparse` + `github.com/tree-sitter/go-tree-sitter` / `tree-sitter-go`). That path **requires CGO** and a C toolchain, which adds build friction, complicates multi-arch releases, and was called out in review as overhead we should avoid.

The product surface is otherwise landed: §12.4 hard metrics on `gopdfsuit` are locked (PR #16), and product `make run` summary/export is in place (PR #18).

**Reference benchmark** (host baseline to preserve):

```text
scanned 78 files (28042 lines) in 295.7ms
  cache: 0 hits, 78 misses (full re-analysis)
  skipped 551 files
915 findings
  severity: 10 high, 197 info, 312 low, 396 medium
  top rules: BP-1 ×181, PERF-6 ×94, PERF-32 ×59, BP-5 ×50, PERF-230 ×44
  example findings: 63 (of 915 total)
exported 915 context file(s) to scripts/findings/functions; exported 37 chunk file(s) to scripts/chunks
```

- Soft wall budget: **295.7ms ± 50ms**
- Hard wall budget: **&lt; 400ms** on the same host/corpus (`make run` / full re-analysis)
- Hard correctness: §12.4 findings **915**, severity **10/197/312/396**, top-five multiset, export **915+37**

## Scope (in)

1. Remove tree-sitter and **all CGO** requirements from the default build.
2. Parse Go with **`go/parser` + `go/ast`** (stdlib) behind a clear language plugin seam.
3. Keep the engine **language-agnostic / pluggable**: `LanguagePlugin` + opaque unit tree so Python/other languages can plug in later without CGO.
4. Rewrite consumers that walk tree-sitter CSTs (PERF facts, BP facts, taint extract) to use the new Go AST (or pure source facts where AST is unnecessary).
5. Update Makefile/README/CI so `CGO_ENABLED=0` builds and tests work.
6. Re-baseline performance and keep wall under hard/soft budgets above.
7. Preserve `make run` product summary surface and §12.4 oracle hard metrics.

## Out of scope

- Implementing Python or other languages (only the plugin seam).
- Rewriting pure-FP museum rule bodies for site-for-site residual swaps.
- Changing §12.4 oracle numbers or fail policy.
- Windows CGO matrix (N/A once pure Go).

## Success criteria

- [x] `CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop` succeeds
- [x] `CGO_ENABLED=0 go test ./...` passes
- [x] No `github.com/tree-sitter/*` deps in `go.mod`
- [x] `make run` wall time **&lt; 400ms** (scan ~170-220ms; process ~0.23-0.30s warm; faster than documented **295.7ms ±50ms** band on this host)
- [x] §12.4 hard metrics still hold: **915** findings; sev **10/197/312/396**; top-five exact; export **915+37**
- [x] Language plugin interface documented for adding a second language without CGO
- [x] README/checklist updated for pure-Go parse path

## Plan

- Checklist: `plans/port-phasewise-checklist.md` (engine parse / Phase 2 notes)
- Architecture: `plans/architecture-go.md`, `documents/architecture-performance.md`
- Branch: `feat/go-ast-no-cgo` (after issue number exists)
- PR: `plans/PR/PR_TEMPLATE.md` → `plans/PR/pr-go-ast-no-cgo.md`

## References

- Relates to #8 (§12.4 gate)
- PRs: #16 (oracle lock), #18 (`make run` summary)
- Docs: `internal/lang/go/goparse/parse.go` (pure-Go parse), `internal/core/plugin.go` (plugin seam + second-language guide)
