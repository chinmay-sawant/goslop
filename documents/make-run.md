# Product run guide (`make run`)

This document is the **markdown run guide**: how to run goslop in the product style used for full-catalogue scans, export generation, and the §12.4 oracle gate.

> There is **no** `--format markdown`. Machine formats are `text` | `json` | `sarif`.  
> “Markdown run” here means **documenting and operating the product run** (Makefile + CLI).

---

## Quick start

```sh
# From the goslop-go repo root
make build
make run SCAN_PATH=/path/to/your/go/project
```

Default without overrides scans `SCAN_PATH` (repo default points at the `gopdfsuit` oracle tree) and writes exports under `scripts/`.

---

## What `make run` does

From the root [`Makefile`](../Makefile):

```make
run: build
	@mkdir -p scripts/findings/functions scripts/chunks
	./bin/goslop --profile all --no-fail --no-terminal $(RUN_ARGS) $(SCAN_PATH)
```

### Defaults

| Variable | Default | Meaning |
|----------|---------|---------|
| `CGO_ENABLED` | `0` | Pure Go build (no CGO) |
| `SCAN_PATH` | `gopdfsuit` absolute path (see Makefile) | Tree to analyze |
| `RUN_ARGS` | `--export-context --export-chunks --no-cache` | Extra CLI flags |

### Expanded equivalent

```sh
CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop
mkdir -p scripts/findings/functions scripts/chunks
./bin/goslop \
  --profile all \
  --no-fail \
  --no-terminal \
  --export-context \
  --export-chunks \
  --no-cache \
  "$SCAN_PATH"
```

### Flag semantics

| Flag | Why `make run` uses it |
|------|-------------------------|
| `--profile all` | Full catalogue: PERF + CWE + BP |
| `--no-fail` | Always exit 0 so exports finish even with many findings |
| `--no-terminal` | **Summary only** on stderr (no multi-thousand-line text dump) |
| `--export-context` | Write one file per finding → `scripts/findings/functions/` |
| `--export-chunks` | Write batches of findings → `scripts/chunks/` |
| `--no-cache` | Full re-analysis (oracle-style, reproducible) |

**Context shape:** both export surfaces default to the **whole enclosing function**
in each finding’s `Context:` block (`[goslop.export] whole_function = true` when
unset). Rebuild with `make build` / `make run` after upgrading so stale short
windows are rewritten. Set `whole_function = false` in `goslop.toml` for the old
nearby ~4-line window. See [export-context-and-chunks.md](./export-context-and-chunks.md).

---

## Overrides

```sh
# Scan a different project
make run SCAN_PATH=./my-service

# Change extra flags (replace RUN_ARGS entirely)
make run RUN_ARGS="--export-context --export-chunks --no-cache"
make run RUN_ARGS="" SCAN_PATH=.                 # summary only, no export, cache allowed
make run RUN_ARGS="--format json --no-cache" SCAN_PATH=.

# Custom chunk size / dirs via RUN_ARGS
make run SCAN_PATH=. RUN_ARGS="--export-context --export-chunks --chunk-size 50 --no-cache"
```

---

## What you see

### Stderr - product summary (always)

Example shape:

```text
scanned 78 files (28120 lines) in 479.5ms
  cache: 0 hits, 78 misses (full re-analysis)
915 findings
  severity: 10 high, 197 info, 312 low, 396 medium
  top rules: BP-1 ×181, PERF-6 ×94, PERF-32 ×59, BP-5 ×50, PERF-230 ×44
exported 915 context file(s) to scripts/findings/functions; exported 37 chunk file(s) to scripts/chunks
```

With `--no-terminal` and default `text` format, **no per-finding lines** are printed to stdout.

### Disk - exports (when `RUN_ARGS` includes export flags)

| Artifact | Path | Count (gopdfsuit oracle) |
|----------|------|---------------------------|
| Per-finding refs | `scripts/findings/functions/1.txt` … `N.txt` | **915** |
| Chunks for delegation | `scripts/chunks/Chunk_1_25.txt` … | **37** (size 25) |

How to use them:

- **Chunks** = combined findings for **delegating** work to agents (one batch per file).  
- **Functions** = individual finding **refs** for single-issue deep dives.  
- **Context** inside each file = **full enclosing function by default** (not only the hit line ± a few lines).

See [export-context-and-chunks.md](./export-context-and-chunks.md).

---

## `make oracle` - hard metrics gate

The oracle target rebuilds exports, writes JSON, and prints count checks used for §12.4 parity:

```make
oracle: build
	@rm -rf scripts/findings/functions scripts/chunks
	@mkdir -p scripts/findings/functions scripts/chunks
	-./bin/goslop --profile $(ORACLE_PROFILE) --format json \
		--export-context --export-chunks --no-cache \
		--context-dir scripts/findings/functions --chunks-dir scripts/chunks \
		$(ORACLE_PATH) > /tmp/goslop-oracle.json
	# Python: findings count, severity histogram, top-5 rules
	# Also: context file count, chunk file count
```

```sh
make oracle
make oracle ORACLE_PATH=./some/project ORACLE_PROFILE=all
```

### Hard targets (default `gopdfsuit` corpus)

| Metric | Target |
|--------|--------|
| Findings | **915** |
| Severity | **10** high / **197** info / **312** low / **396** medium |
| Top rules | BP-1×181, PERF-6×94, PERF-32×59, BP-5×50, PERF-230×44 |
| Exports | **915** context + **37** chunks |
| Wall (pure Go) | **&lt;400ms** (checklist gate) |

The leading `-` on the scan command allows non-zero exit under fail policy so the summary still runs. Ledger: [`plans/port-phasewise-checklist.md`](../plans/port-phasewise-checklist.md) §12.4.

---

## Other Makefile targets

| Target | Purpose |
|--------|---------|
| `build` | `go build -o bin/goslop ./cmd/goslop` |
| `test` | `go test ./...` |
| `integration` | Fixture harness under `tests/integration` |
| `vet` / `fmt` / `lint` / `lint-all` | Static checks |
| `ci` | `lint` + `test` + `build` (local CI parity) |
| `version` | Build + `./bin/goslop --version` |
| `help` | Print targets and current `SCAN_PATH` / `RUN_ARGS` |

```sh
make help
make ci
make test
make integration
```

---

## End-to-end workflows

### A. Product-style full scan + agent handoff

```sh
make run SCAN_PATH=/path/to/project

# Delegate a batch of 25 findings to an agent
# (open or paste scripts/chunks/Chunk_1_25.txt)

# Deep-dive a single finding by index
# (open scripts/findings/functions/12.txt)
```

### B. Summary-only without wiping your head

```sh
./bin/goslop --profile all --no-fail --no-terminal --no-cache /path/to/project
```

### C. JSON + exports (oracle-like, custom path)

```sh
./bin/goslop --profile all --format json \
  --export-context --export-chunks --no-cache \
  --context-dir scripts/findings/functions \
  --chunks-dir scripts/chunks \
  /path/to/project > /tmp/goslop.json
```

### D. SARIF for CI (not the same as `make run`)

```sh
./bin/goslop --profile recommended --no-fail --format sarif . > goslop.sarif
```

See [reporting-formats.md](./reporting-formats.md).

### E. Everyday developer loop

```sh
make lint && make test && make build
./bin/goslop --profile recommended .
```

---

## CI note

GitHub Actions (`.github/workflows/ci.yml`) runs **vet / test / build** and smoke `--version` / `--list-rules`. It does **not** run `make run` or `make oracle` against external corpora. Use those locally or in a dedicated workflow.

---

## Related docs

- [cli-reference.md](./cli-reference.md) - every flag  
- [export-context-and-chunks.md](./export-context-and-chunks.md) - using generated snippets  
- [reporting-formats.md](./reporting-formats.md) - text / JSON / SARIF  
- [overview.md](./overview.md) - product surface  
