# Taint tracking

goslop provides experimental Go taint tracking for triaging data flow into
four security rule families: `CWE-22` (path traversal), `CWE-78` (command
injection), `CWE-79` (XSS-related output), and `CWE-89` (SQL injection).

It is an AST- and name-string-based analysis with bounded same-package
inter-procedural summaries. It is useful evidence for review, not a replacement
for code review or a full-program security analysis.

## Enable it

```sh
# Explicitly enable experimental taint tracking.
./bin/goslop --taint .

# The security profile enables taint and uses a deeper default.
./bin/goslop --profile security .

# Include a bounded path explanation in findings.
./bin/goslop --taint --taint-depth 3 --taint-show-paths .

# Disable it even when the security profile is selected.
./bin/goslop --profile security --no-taint .
```

Configuration equivalents:

```toml
[goslop.taint]
enabled = true
show_paths = true
```

`--taint-depth` accepts 1 through 4. The normal default is 1; the `security`
profile defaults to 4. The `recommended` profile allow-lists the related CWE
IDs but does not enable taint by default.

## What it models

The detector extracts source and sink facts from Go source, constructs a
per-unit flow graph, and then combines bounded same-package call summaries at
scan finalization. It retains the necessary scan state on warm cache hits, so
taint-enabled warm scans still perform project-state work.

Recognized source categories include request and framework parameters,
arguments, environment variables, file reads, and network input. Supported
reported sink categories are command execution, SQL queries, file operations,
templates, and HTTP response writes. The implementation also recognizes some
deserialization, LDAP, and XML sink shapes internally, but those do not expand
the four taint rule IDs reported by the current detector.

Common sanitizer recognition includes `filepath.Base`, HTML and URL escaping,
numeric parsing, selected LDAP/XML escaping helpers, regular-expression
validation, and functions whose names begin with `sanitize`, `escape`,
`validate`, or `purify`. This is intentionally heuristic. In particular,
`filepath.Clean` and `path.Clean` alone are not treated as safe path
confinement.

## Evidence output

Use `--taint-show-paths` to attach flow evidence to a finding. The text
reporter includes the source kind, sink kind, and hop count; JSON and SARIF
carry the same finding evidence for automation.

```text
taint flow UserInput.r.URL.Query -> CommandExec.exec.Command across 1 hop
  hop: runCommand(cmd) at handler.go:42
```

## Important limits

- Call resolution is bounded to the same package and the configured depth.
- Imported calls and ambiguous method receivers can be declined rather than
  guessed.
- It uses AST/name matching rather than SSA, Go type information, or a
  whole-program call graph.
- Interface dispatch, aliasing, deep recursion, and complex concurrent flows
  can produce false negatives; heuristic names can produce false positives.
- Taint output should be reviewed with the surrounding code before it becomes
  a security decision.

For all taint flags and precedence rules, see
[cli-reference.md](./cli-reference.md). For the supported profile selection,
see [go-recommended-pack.md](./go-recommended-pack.md).
