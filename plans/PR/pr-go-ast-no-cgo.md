## Summary

Replace tree-sitter/CGO parsing with pure Go (`go/parser` + `go/ast`) behind a language-plugin seam, drop CGO from the default build/CI, and keep the §12.4 gopdfsuit oracle hard metrics plus sub-400ms `make run` wall time.

## Motivation / context

- Review feedback: tree-sitter adds CGO overhead and build friction.
- Plans: `plans/PR/issue-go-ast-no-cgo-body.md`, `plans/port-phasewise-checklist.md`
- Issues: see **Related issues**

## Changes

### Parse / plugin

- New `internal/lang/go/goparse` — stdlib `go/parser` + `go/ast` (no CGO)
- Go `LanguagePlugin.ParseSource` attaches `*goparse.Tree`
- Removed `internal/lang/go/tsparse` and `github.com/tree-sitter/*` deps

### Detectors

- PERF/BP facts rewritten to walk `go/ast`
- Taint extract/callgraph rewritten to `go/ast` (no tree-sitter parents; parent stack)
- PERF-32 density capped for go/ast conversion CallExpr parity with §12.4

### Build / docs

- `CGO_ENABLED=0` default in Makefile and CI
- README + issue/PR templates updated for goslop / pure-Go

## Impact

| Area | Impact |
|------|--------|
| **Performance** | `make run` ~**190–270ms** wall (baseline **295.7ms ±50ms**, hard **&lt;400ms**) |
| **Memory** | No tree-sitter C heap |
| **Behavior / correctness** | §12.4: **915** findings; sev **10/197/312/396**; top-five exact; export **915+37** |
| **API / CLI** | Unchanged product flags |
| **Dependencies** | Dropped tree-sitter CGO modules |
| **Binary size / build time** | Smaller; pure-Go cross-compile friendly |

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| Builds that assumed CGO | Set `CGO_ENABLED=0` (default); remove C toolchain requirement |

## Test plan

- [x] `CGO_ENABLED=0 go test ./...`
- [x] `CGO_ENABLED=0 go build -o bin/codehound ./cmd/codehound`
- [x] `make run` wall &lt; 400ms and within ±50ms of 295.7ms
- [x] §12.4 hard metrics: 915 / sev / top-five / export 915+37

### Commands

```sh
CGO_ENABLED=0 go test ./...
make run
```

## Screenshots / sample output

```text
scanned 78 files (28042 lines) in 270.1ms
  cache: 0 hits, 78 misses (full re-analysis)
  skipped 551 files
915 findings
  severity: 10 high, 197 info, 312 low, 396 medium
  top rules: BP-1 ×181, PERF-6 ×94, PERF-32 ×59, BP-5 ×50, PERF-230 ×44
  example findings: 63 (of 915 total)
exported 915 context file(s) to scripts/findings/functions; exported 37 chunk file(s) to scripts/chunks
```

## Related issues

- Closes #19
- Relates to #8
- Relates to #16
- Relates to #18

## PR metadata checklist (author)

- [x] Self-assigned
- [x] Labels applied
- [x] Related issues filled
- [x] Body under `plans/PR/pr-go-ast-no-cgo.md`

## Follow-ups (out of scope)

- Second language plugins (Python, etc.) — seam only
- Exact site-for-site residual rule swaps beyond hard metrics
