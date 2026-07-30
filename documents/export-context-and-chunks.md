# Export context and chunks

After a scan, goslop can write **human-readable text dossiers** for every finding. There are two complementary export surfaces:

| Export | CLI flag | Default directory | Unit of work |
|--------|----------|-------------------|--------------|
| **Function / context refs** | `--export-context` | `scripts/findings/functions` | **One finding → one file** |
| **Chunks** | `--export-chunks` | `scripts/chunks` | **N findings → one file** (default N = **25**) |

Implementation: `internal/export/export.go`.

---

## When to use which

### `scripts/findings/functions/` - individual finding refs

Use for:

- Looking up **one** issue by index  
- Linking a ticket/PR comment to a single dossier  
- Agent prompts that should focus on a **single** finding  
- Comparing two specific findings side by side  

Files: `1.txt`, `2.txt`, … `N.txt` (1-based scan order).

> The directory name is `functions` for historical product parity. By default each file’s **Context** is the **whole enclosing function** (outermost `FuncDecl` preferred; otherwise `FuncLit`). Set `[goslop.export] whole_function = false` for the previous nearby ~4-line window. After upgrading, run `make build` (or `make run`) so regenerated files replace stale short windows.

### `scripts/chunks/` - combined findings for delegation

Use for:

- **Delegating batches of work to agents** (parallel review sessions)  
- Bulk triage (“fix this pack of 25”)  
- Multi-finding remediation passes  

Each `Chunk_START_END.txt` concatenates the **same content** as the numbered function files for that index range, with a batch header and separators.

**Rule of thumb**

| Goal | Use |
|------|-----|
| Hand a batch to an agent / subagent | **`scripts/chunks/Chunk_*.txt`** |
| Open or cite one finding | **`scripts/findings/functions/N.txt`** |

A **chunk is the combined functions** (findings) for that range. Function files remain the stable per-finding refs.

---

## Workflow: scan → export → delegate

```text
1. Scan (+ export flags)
   ./bin/goslop --profile all --export-context --export-chunks --no-cache /path

2. Artifacts on disk
   scripts/findings/functions/{1..N}.txt     # per-finding refs
   scripts/chunks/Chunk_{start}_{end}.txt   # batches of combined findings

3. Delegate
 - Assign Chunk_1_25.txt to agent A
 - Assign Chunk_26_50.txt to agent B
 - When an agent needs one issue, open functions/42.txt
```

Product helpers:

```sh
make run SCAN_PATH=/path/to/project
# defaults: --export-context --export-chunks --no-cache
# (plus --profile all --no-fail --no-terminal)

make oracle   # wipe + re-export + count checks (915 + 37 on gopdfsuit)
```

Stderr after a successful export:

```text
exported 915 context file(s) to scripts/findings/functions; exported 37 chunk file(s) to scripts/chunks
```

---

## CLI flags

| Flag | Default | Description |
|------|---------|-------------|
| `--export-context` | off | Write per-finding files |
| `--export-chunks` | off | Write chunk files |
| `--context-dir` | `scripts/findings/functions` | Override context directory |
| `--chunks-dir` | `scripts/chunks` | Override chunks directory |
| `--chunk-size` | **25** | Findings per chunk |

### Config: whole-function vs line window

| Key | Default | Description |
|-----|---------|-------------|
| `[goslop.export] whole_function` | **`true`** | When true, Context includes the full enclosing function. When false, uses `line-2 … line+1` around the hit. |

```toml
# goslop.toml
[goslop.export]
whole_function = true   # default; omit for the same effect
# whole_function = false  # nearby ~4-line window only
```

Both **function refs** and **chunks** share this setting (same `formatFindingBlock` path).

```sh
# Both surfaces (typical)
./bin/goslop --profile all --export-context --export-chunks --no-cache .

# Only individual refs
./bin/goslop --export-context .

# Only batches for agents
./bin/goslop --export-chunks --chunk-size 25 .

# Custom dirs / larger batches
./bin/goslop --export-context --export-chunks \
  --context-dir /tmp/ch-ctx --chunks-dir /tmp/ch-chunks \
  --chunk-size 50 /path/to/project
```

**Constraints**

- Context and chunk directories **must differ** when both exports are enabled.  
- Before writing: context dir deletes all `*.txt`; chunks dir deletes `Chunk_*.txt`.  
- Sources are retained in memory only when an export flag is set (`RetainSources`).

---

## File naming

### Function refs

```text
scripts/findings/functions/{index}.txt
```

- `index` is 1-based position in the post-scan finding list.  
- Header inside the file: `Finding {index}/{total}`.

### Chunks

```text
scripts/chunks/Chunk_{start}_{end}.txt
```

- Inclusive 1-based range of finding indices.  
- With chunk size 25 and 915 findings → **37** files:

| File | Findings |
|------|----------|
| `Chunk_1_25.txt` | 1-25 |
| `Chunk_26_50.txt` | 26-50 |
| … | … |
| `Chunk_876_900.txt` | 876-900 |
| `Chunk_901_915.txt` | 901-915 (partial last batch) |

### Cross-reference

| Finding index | Function file | Chunk (size 25) |
|---------------|---------------|-----------------|
| 1 | `1.txt` | `Chunk_1_25.txt` |
| 25 | `25.txt` | `Chunk_1_25.txt` |
| 26 | `26.txt` | `Chunk_26_50.txt` |
| 915 | `915.txt` | `Chunk_901_915.txt` |

Indices in headers stay **global** (`Finding i/total`) even inside a chunk, so refs stay consistent across formats.

---

## Finding block format

Shared by function files and chunks:

```text
Finding {i}/{total}
Source: {file}:{line}:{column}
Rule: {rule_id}
Fingerprint: {fingerprint}
Rule title: {title}
Severity: {severity}
Message: {message}
[CWEs: ...]
[Fix: ...]
[Confidence: ...]
[Tags: ...]
[Remediation: ...]
Context:
    {context lines}
```

### Context lines

With **`whole_function = true`** (default):

1. Resolve the enclosing function via pure-Go `go/parser` + `go/ast`: prefer the **outermost `FuncDecl`** containing the hit (so a hit inside `defer func() { ... }` still exports the full named method); otherwise the outermost `FuncLit`.  
2. Emit that full line span with `>` on the hit line.  
3. If no enclosing function (package-level finding, non-Go source, parse failure) → fall back to the nearby window.  
4. Missing source → `<context unavailable>` (or detector `Snippet` if present).

With **`whole_function = false`**:

1. Prefer `Finding.Snippet` when present.  
2. Else a window from source cache / disk: **line-2 … line+1**, with `>` on the hit line.  
3. Missing source → `<context unavailable>`.

### Chunk wrapper

```text
Findings {start}-{end} of {total}

{finding block 1}

====================================================================================================

{finding block 2}
...
```

Separator is **100** `=` characters between findings.

---

## Generated snippet examples

Examples below match the §12.4 **gopdfsuit** oracle export tree that may already exist under the repo after `make run` / `make oracle`.

### Function ref - `scripts/findings/functions/1.txt`

```text
Finding 1/915
Source: /home/chinmay/.../gopdfsuit/bindings/python/cgo/exports.go:1:1
Rule: BP-57
Fingerprint: goslop:2:BP-57:.../exports.go:fe9363bf4dff0127
Rule title: Stale Go Version In go.mod
Severity: low
Message: go.mod targets an out-of-support Go major release; update to a currently supported baseline
Fix: Project-level go.mod parse comparing the declared Go version with policy.
Context:
    >     1: // Package main provides CGO exports for the Python bindings.
          2: package main
```

### Function ref with whole function - `scripts/findings/functions/915.txt`

```text
Finding 915/915
Source: .../gopdfsuit/typstsyntax/renderer.go:1256:9
Rule: PERF-35
Fingerprint: goslop:2:PERF-35:.../renderer.go:619456b9f1ed1271
Rule title: Interface Boxing On Hot Path
Severity: info
Message: fmt.Sprintf / Errorf boxes arguments through interface{} on a hot path
Fix: Cast non-string args to a concrete type or use strconv/strings builders to avoid interface boxing.
Context:
       1255: func fmtFloat(f float64) string {
    >  1256: 	return fmt.Sprintf("%.2f", f)
       1257: }
```

(With default `whole_function = true`, larger functions include every line of the enclosing body, not only a 4-line window.)

### Chunk (abbreviated) - `scripts/chunks/Chunk_1_25.txt`

```text
Findings 1-25 of 915

Finding 1/915
Source: .../exports.go:1:1
Rule: BP-57
...
Context:
    >     1: // Package main provides CGO exports for the Python bindings.
          2: package main
====================================================================================================

Finding 2/915
Source: .../exports.go:1:1
Rule: BP-60
...
====================================================================================================

... Findings 3-24 ...

====================================================================================================

Finding 25/915
...
```

### Last partial chunk - `scripts/chunks/Chunk_901_915.txt`

```text
Findings 901-915 of 915

Finding 901/915
Source: .../typstsyntax/parser.go:258:11
Rule: PERF-230
...
====================================================================================================

Finding 902/915
...
```

---

## Oracle counts (reference)

| Artifact | Count | Location |
|----------|------:|----------|
| Context / function files | **915** | `scripts/findings/functions/` |
| Chunk files | **37** | `scripts/chunks/` |
| Default chunk size | **25** | `export.DefaultChunkSize` |

Math: `ceil(915 / 25) = 37`.

---

## Agent / AI delegation recipes

### Parallel batch review

```sh
make run SCAN_PATH=/path/to/project

# Agent A: scripts/chunks/Chunk_1_25.txt
# Agent B: scripts/chunks/Chunk_26_50.txt
# Agent C: scripts/chunks/Chunk_51_75.txt
# ...
```

Prompt sketch:

> You are reviewing goslop findings. The attached chunk contains up to 25 findings. Each finding’s Context block is the full enclosing Go function (unless the project set `whole_function = false`). For each finding: confirm true/false positive, propose a minimal fix, and note the rule id and file:line. If you need a single finding alone, ask for `scripts/findings/functions/{i}.txt`.

### Single-finding deep dive

```sh
# After export, open one ref
less scripts/findings/functions/42.txt
```

### Larger batches for fewer agents

```sh
./bin/goslop --profile all --export-chunks --chunk-size 50 --no-cache /path
# → Chunk_1_50.txt, Chunk_51_100.txt, ...
```

### Keep JSON machine output + disk dossiers

```sh
./bin/goslop --profile all --format json \
  --export-context --export-chunks --no-cache \
  /path > findings.json
```

Stdout JSON is independent of the text files under `scripts/`.

---

## Defaults (code constants)

```go
DefaultContextDir      = "scripts/findings/functions"
DefaultChunksDir       = "scripts/chunks"
DefaultChunkSize       = 25
DefaultWholeFunction   = true
```

Paths are relative to the **process working directory** unless absolute paths are passed via `--context-dir` / `--chunks-dir`.

---

## Related docs

- [make-run.md](./make-run.md) - product `make run` that generates these files  
- [cli-reference.md](./cli-reference.md) - flag details  
- [reporting-formats.md](./reporting-formats.md) - stdout formats (separate from disk export)  
- [overview.md](./overview.md) - product surface  
