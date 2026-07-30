# Reporting formats (text · JSON · SARIF)

goslop writes finding reports to **stdout** and the product **scan summary** to **stderr**. Choose a format with `--format`:

| Format | Flag | Primary use |
|--------|------|-------------|
| **text** | `--format text` (default) | Terminal / CI logs / editor problem matchers |
| **json** | `--format json` | Scripts, dashboards, automation |
| **sarif** | `--format sarif` | GitHub Code Scanning, SARIF viewers, IDEs |

```sh
./bin/goslop .
./bin/goslop --format json .
./bin/goslop --format sarif . > goslop.sarif
```

Invalid formats are rejected (exit 2): only `text` | `json` | `sarif`.

---

## Dual streams

```text
stderr ──► scan summary (files, lines, wall time, cache hits, severity histogram, top rules)
stdout ──► findings (text lines | JSON envelope | SARIF log)
```

This keeps JSON/SARIF **pipe-clean**:

```sh
# Pure SARIF file (discard summary)
./bin/goslop --format sarif . > goslop.sarif 2>/dev/null

# Keep summary visible while writing SARIF
./bin/goslop --format sarif . > goslop.sarif
```

### `--no-terminal`

With **text** format, `--no-terminal` prints **summary only** (no per-finding dump).  
With **json** / **sarif**, machine output still goes to stdout:

```sh
./bin/goslop --no-terminal .                         # summary only
./bin/goslop --no-terminal --format json .           # summary + JSON findings
```

`make run` uses `--no-terminal` for product summary scans. See [make-run.md](./make-run.md).

---

## Text format

### Purpose

Human-readable one-liner per finding, friendly to `file:line:col` tooling.

### Shape

```text
RULE_ID file:line:column message
```

- Line and column clamp to **1** if missing/zero.
- No severity, fingerprint, or CWE on the line (use JSON for full fields).

### Example

```text
PERF-1 a.go:2:4 hot path alloc
CWE-78 cmd.go:12:5 unsanitized command
PERF-6 loop.go:3:1 fmt in loop
```

### With summary (illustrative)

```text
# stderr
scanned 12 files (3400 lines) in 42.1ms
  cache: 3 hits, 9 misses
12 findings
  severity: 2 high, 1 info, 4 low, 5 medium
  top rules: PERF-6 ×3, CWE-89 ×2, ...

# stdout
PERF-6 loop.go:3:1 fmt in loop
CWE-89 db.go:40:8 string-built query
...
```

`NO_COLOR` is recognized for future styling; the current text reporter emits plain text.

---

## JSON format

### Purpose

Stable machine envelope for scripts, custom gates, and tooling that needs full finding fields.

### Envelope

```json
{
  "findings": [ /* Finding objects */ ],
  "version": "0.1.0-dev"
}
```

- Empty scan → `"findings": []` (never null).
- Each finding gets a **fingerprint** if missing before encode.

### Finding fields

| Field | Type | Description |
|-------|------|-------------|
| `rule_id` | string | e.g. `PERF-6`, `CWE-78`, `BP-57` |
| `rule_title` | string | Human title |
| `file` | string | Path as reported by the analyzer |
| `line` / `column` | int | 1-indexed location |
| `message` | string | Finding message |
| `severity` | string | `info` \| `low` \| `medium` \| `high` \| `critical` |
| `cwe` | array | CWE refs (`id`, optional `name`, `url`) |
| `fingerprint` | string | Stable id (see below) |
| `fix` | string | Optional fix hint |
| `snippet` | string | Optional source excerpt |
| `evidence` | object | Optional structured data (e.g. taint paths) |
| `confidence` | number | Optional |
| `suppressed` | bool | Optional |
| `remediation` | string | Optional |
| `tags` | string[] | Optional |

### Fingerprint (v2)

```text
goslop:2:{ruleID}:{file-with-forward-slashes}:{first-16-hex-of-sha256(message)}
```

Content-stable across line shifts when rule + file + message stay the same. Used by baseline matching and SARIF `partialFingerprints`.

### Full example

```json
{
  "findings": [
    {
      "rule_id": "CWE-89",
      "rule_title": "SQL injection",
      "file": "main.go",
      "line": 10,
      "column": 3,
      "message": "unsanitized query",
      "severity": "high",
      "cwe": [
        {
          "id": "CWE-89",
          "url": "https://cwe.mitre.org/data/definitions/89.html"
        }
      ],
      "fingerprint": "goslop:2:CWE-89:main.go:a1b2c3d4e5f67890",
      "fix": "Use parameterized queries or a query builder."
    },
    {
      "rule_id": "PERF-6",
      "rule_title": "Fmt In Loop",
      "file": "loop.go",
      "line": 3,
      "column": 1,
      "message": "fmt in loop",
      "severity": "medium",
      "cwe": [],
      "fingerprint": "goslop:2:PERF-6:loop.go:1122334455667788"
    }
  ],
  "version": "0.1.0-dev"
}
```

### Empty findings

```json
{"findings":[],"version":"0.1.0-dev"}
```

### Usage

```sh
./bin/goslop --format json . > findings.json
./bin/goslop --profile all --format json --no-cache . | jq '.findings | length'
./bin/goslop --format json . | jq '[.findings[] | select(.severity=="high")]'
```

---

## SARIF 2.1.0

### Purpose

Industry-standard static-analysis log for:

- **GitHub Code Scanning** (`upload-sarif`)
- VS Code / IDE **SARIF Viewer** extensions
- Azure DevOps and other SARIF 2.1.0 consumers

Implementation: `internal/reporting/sarif.go` - **minimal valid** SARIF 2.1.0 subset.

### Schema metadata

| Property | Value |
|----------|--------|
| `$schema` | `https://json.schemastore.org/sarif-2.1.0.json` |
| `version` | `2.1.0` |
| `runs[0].tool.driver.name` | `goslop` |
| `informationUri` | `https://github.com/chinmay-sawant/goslop` |
| `version` (driver) | tool version, e.g. `0.1.0-dev` |

### Severity → SARIF `level`

| goslop severity | SARIF level |
|--------------------|-------------|
| `info` | `note` |
| `low` | `warning` |
| `medium` | `warning` |
| `high` | `error` |
| `critical` | `error` |

### What is included

| SARIF field | Source |
|-------------|--------|
| `driver.rules[].id` | Unique `RuleID` (first-seen order) |
| `driver.rules[].name` / `shortDescription` | `RuleTitle` (or id) |
| `results[].ruleId` | `RuleID` |
| `results[].level` | Mapped severity |
| `results[].message.text` | Message |
| `locations[].physicalLocation.artifactLocation.uri` | File path |
| `region.startLine` / `startColumn` | Line / column (min 1) |
| `partialFingerprints.primaryLocationLineHash` | Finding fingerprint |

### What is not included (minimal emitter)

- Full rule help / CWE on rule objects  
- `endLine` / `endColumn`  
- Code flows / related locations  
- Fix objects  
- Full five-level severity (collapsed to note/warning/error)

For richest fields, use **`--format json`** in parallel.

### Full example SARIF

```json
{
  "$schema": "https://json.schemastore.org/sarif-2.1.0.json",
  "version": "2.1.0",
  "runs": [
    {
      "tool": {
        "driver": {
          "name": "goslop",
          "informationUri": "https://github.com/chinmay-sawant/goslop",
          "version": "0.1.0-dev",
          "rules": [
            {
              "id": "CWE-78",
              "name": "OS Command Injection",
              "shortDescription": {
                "text": "OS Command Injection"
              }
            },
            {
              "id": "PERF-6",
              "name": "PERF-6",
              "shortDescription": {
                "text": "PERF-6"
              }
            }
          ]
        }
      },
      "results": [
        {
          "ruleId": "CWE-78",
          "level": "error",
          "message": {
            "text": "unsanitized command"
          },
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": {
                  "uri": "cmd.go"
                },
                "region": {
                  "startLine": 12,
                  "startColumn": 5
                }
              }
            }
          ],
          "partialFingerprints": {
            "primaryLocationLineHash": "goslop:2:CWE-78:cmd.go:abcdef0123456789"
          }
        },
        {
          "ruleId": "PERF-6",
          "level": "warning",
          "message": {
            "text": "fmt in loop"
          },
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": {
                  "uri": "loop.go"
                },
                "region": {
                  "startLine": 3,
                  "startColumn": 1
                }
              }
            }
          ],
          "partialFingerprints": {
            "primaryLocationLineHash": "goslop:2:PERF-6:loop.go:1122334455667788"
          }
        }
      ]
    }
  ]
}
```

### Empty findings

Still emits one `run` with driver metadata; `results` is `[]`. Rule array may be omitted when empty.

### Generate SARIF

```sh
./bin/goslop --profile recommended --format sarif . > goslop.sarif

# Advisory CI (do not fail the job before upload)
./bin/goslop --profile recommended --no-fail --format sarif . > goslop.sarif
```

### GitHub Code Scanning (example workflow step)

```yaml
- name: Build goslop
  run: CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop

- name: Run goslop (SARIF)
  run: ./bin/goslop --profile recommended --no-fail --format sarif . > goslop.sarif

- name: Upload SARIF
  uses: github/codeql-action/upload-sarif@v3
  with:
    sarif_file: goslop.sarif
  # permissions: security-events: write
```

**Compatibility notes**

1. Prefer scanning from the **repo root** so `uri` paths match the GitHub workspace.  
2. Driver `rules` only lists rules that **fired** (not the full pack).  
3. Low and medium both map to SARIF `warning`.  
4. `partialFingerprints` use goslop’s content fingerprint (stable when message/path/rule hold).  
5. Exit codes are orthogonal to format: use `--no-fail` for advisory generation.

### VS Code / other tools

| Tool | How |
|------|-----|
| SARIF Viewer (VS Code) | Open `goslop.sarif`; resolve paths against the workspace |
| GitHub Security tab | After `upload-sarif` |
| `jq` / scripts | Prefer **JSON** for field-rich automation |
| Problem matchers | Prefer **text** (`RULE file:line:col message`) |

---

## Format comparison

| Concern | text | json | sarif |
|---------|------|------|-------|
| Human scan | ✅ | summary on stderr | summary on stderr |
| Scripts | weak | ✅ native | parse SARIF |
| GitHub Code Scanning | ❌ | ❌ | ✅ |
| Full finding fields | ❌ | ✅ | partial |
| Fingerprints | ❌ | `fingerprint` | `partialFingerprints` |
| Severity fidelity | none on line | 5 levels | 3 levels |

---

## Exit codes (all formats)

| Code | Meaning |
|------|---------|
| 0 | Clean / `--no-fail` |
| 1 | Findings vs fail policy |
| 2 | Usage / config |
| 3 | Internal |

---

## Related

- [cli-reference.md](./cli-reference.md)  
- [make-run.md](./make-run.md)  
- [export-context-and-chunks.md](./export-context-and-chunks.md) - disk exports are **independent** of `--format`; Context defaults to the **whole enclosing function** (`[goslop.export] whole_function`)
