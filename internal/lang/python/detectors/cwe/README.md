# Python CWE detectors

Pure-Go **source-pattern** security heuristics for Python (issue #52).

## Scope

| Piece | Role |
|-------|------|
| `PyCweScan` | Unified detector (`LanguagePython`, many rules) |
| `RegisterRule` / `init` | Catalogue registration with optional needle gates |
| `PyCweFacts` + `pyCweNeedles` | Per-file `ast.SourceIndex` prefilter |
| `metadata*.go` | Hand-authored metadata, split by detector domain |
| `rules*.go` | Heuristics, split by detector domain to stay within the file-size budget |

**No CGO / tree-sitter.** Detectors operate on `unit.Source` only (source-only `ParsedUnit`).

## Implemented rules (v0.0.2 / #52)

| Rule | Sinks (high signal) | Safe suppressions |
|------|---------------------|-------------------|
| CWE-502 | `pickle.loads/load/Unpickler`, `yaml.load` without SafeLoader, `yaml.unsafe_load` | `yaml.safe_load`, `json.loads` |
| CWE-78 | `os.system` / `os.popen` dynamic; `subprocess.*(shell=True)` dynamic | list argv + `shell=False` / default, pure string literals |
| CWE-89 | `execute` / `executemany` with f-string / `%` / `.format` SQL | parameterized `?` / `%s` + bound args tuple |
| CWE-22 | `open(os.path.join(root, user))`, `Path(root) / user` without confinement | `os.path.basename` only; `resolve` + `startswith` root |
| CWE-79 | `mark_safe`, `Markup(`, `render_template_string` dynamic | plain `render_template("…html", …)` |
| CWE-88 / 90 / 91 / 93 / 94 / 117 | Dynamic subprocess arguments, LDAP/XPath/header/code/log sinks | Literal, escaped, parameterized, or sanitized expressions |
| CWE-214 / 215 / 695 / 749 / 829 | Sensitive process args, debug output, low-level APIs, dynamic route/code loading | Password-file options, redacted logs, package-controlled imports |
| CWE-256 / 260 / 261 / 312 / 319 / 523 / 547 / 798 | Password/secret storage and transport, weak encoding, insecure security settings | Environment/secret-provider values, TLS, default verification, secure settings |
| CWE-914 / 915 / 916 | Request-controlled names/attributes and fast password hashing | Allowlisted attributes and modern password hashes |
| CWE-295 / 328 / 335 / 338 / 347 / 1204 / 1240 / 1241 / 1392 | Certificate/signature bypasses, weak hashes/PRNG/IVs, risky crypto, default passwords | Verified TLS/signatures, `secrets`, runtime IVs, vetted crypto, non-default credentials |
| CWE-601 / 605 / 918 / 924 / 940 / 941 | Request-controlled redirects/HTTP/mail destinations and channel integrity | Allowlisted destinations, signature/state checks, non-wildcard bindings |
| CWE-112 / 611 / 776 | Request XML parsing and explicit entity-expansion configuration | Schema-aware parsing and disabled entity resolution |
| CWE-41 / 59 / 73 / 250 / 276 / 378 / 426 / 494 | Filesystem path/link/permission/temp/search-path and downloaded-code hazards | Canonical paths, private modes, secure temp APIs, trusted import/download flows |
| CWE-15 / 489 / 756 / 921 / 1051 / 1052 / 1125 / 1188 | Debug endpoints, insecure defaults, and web-server configuration | Production-safe configuration and disabled diagnostics |
| CWE-306 / 307 / 346 / 359 / 565 / 613 / 698 / 807 | Authentication, session, authorization, and privacy boundaries | Authentication, throttling, state validation, and protected session data |
| CWE-201 / 204 / 208 / 209 / 212 / 213 / 488 / 497 | Information exposure through errors, timing, logging, and diagnostics | Generic errors, constant-time checks, and redacted output |
| CWE-140 / 1173 / 1230 / 1236 / 1286 / 1289 / 1333 / 1389 | Input validation, parser limits, and export controls | Explicit allowlists, bounded parsing, and validated output |
| CWE-252 / 390 / 396 / 397 / 478 / 584 | Platform API error handling and language-quality hazards | Checked return values, explicit errors, and safe API usage |
| CWE-379 / 427 / 434 / 459 / 477 / 708 / 770 / 772 | Upload, resource, temporary-file, and deprecated API hazards | Trusted paths, bounded resources, secure upload and cleanup flows |
| 60 Tier-B CWE rules | Narrow same-source patterns for residual parser, resource, path, and runtime hazards | Explicit safe alternatives in paired fixtures |

## Non-goals

- Full 344-rule catalogue registration
- BP-PY / PERF detectors (siblings #53 / #54)
- Inter-procedural taint graph
- Modifying Go CWE detectors

## Adding a rule later

1. Ensure the CWE id exists in `ruleset/python/chunks/cwe-*.json`.
2. Add `MetaCWEN` in `metadata.go` or the matching domain metadata file (titles aligned with catalogue).
3. Add a needle to `pyCweNeedles` only when it is proven FN-safe for the detector gate.
4. `RegisterRule("CWE-N", detect, &MetaCWEN, gates...)` in the matching domain file; omit gates rather than risk suppressing a valid finding.
5. Unit hit/miss tests in the matching domain test file.

## Tests

Per-rule Python CWE tests use the same paired `.txt` fixtures as the integration
matrix; do not embed Python source snippets in Go test files.

```sh
go test ./internal/lang/python/detectors/cwe/ -count=1
```

Wired via `internal/lang/python/detectors.All()` → Python plugin `Detectors` / `NewDetectors`.
