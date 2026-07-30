# Suppressions and baselines

Use a narrow inline suppression for a reviewed local exception. Use a baseline
for existing repository debt that should remain visible but not fail an initial
rollout. Neither mechanism replaces fixing a confirmed issue.

## Inline suppressions

goslop recognizes directives only in real `//` or `#` line comments, not in
string literals or block comments.

```go
// Suppress the next code line.
// goslop-ignore: PERF-101
srv := &http.Server{Addr: ":8080"}

// Suppress an issue reported on this line.
runSlowPath() // goslop-ignore: PERF-1

// Suppress a selected block.
// goslop-ignore-start: PERF-1,PERF-7
for _, value := range values {
	defer release(value)
}
// goslop-ignore-end

// Suppress selected rules for a file when this directive is in its first 20 lines.
// goslop-ignore-file: BP-21,BP-28
```

Use `all` in place of a list only when the entire scope has a documented,
short-lived reason:

```go
// goslop-ignore-start: all
// goslop-ignore-end
```

Inspect suppressed findings during review with:

```sh
./bin/goslop --show-ignored .
```

## Baseline behavior

By default, goslop looks upward from the scan root for
`.goslop-baseline.json`, stopping at the repository `.git` directory. Supply a
specific file with `--baseline-file`; disable loading with `--no-baseline`.

```sh
./bin/goslop --baseline-file ci/goslop-baseline.json .
./bin/goslop --no-baseline .
./bin/goslop --show-baselined .
```

Baselined findings are removed by default. With `--show-baselined`, they are
retained as suppressed informational findings. Matching prefers a finding
fingerprint and falls back to its file, line, and column. An entry with an
`expires` value earlier than the current UTC timestamp no longer suppresses a
finding.

The CLI loads and filters baselines; it does not currently provide a baseline
creation, update, diff, or prune subcommand. Author a baseline with a reviewed
script or tool that produces this wire shape:

```json
{
  "version": "1",
  "generated_at": "2026-07-30T00:00:00Z",
  "tool_version": "0.1.0-dev",
  "entries": {
    "PERF-101": [
      {
        "file": "server.go",
        "line": 12,
        "column": 2,
        "fingerprint": "reviewed-finding-fingerprint",
        "reason": "planned timeout migration",
        "expires": "2026-12-31T00:00:00Z"
      }
    ]
  }
}
```

Keep baseline entries specific, include a reason and expiration when possible,
and remove each entry when the underlying code is fixed.
