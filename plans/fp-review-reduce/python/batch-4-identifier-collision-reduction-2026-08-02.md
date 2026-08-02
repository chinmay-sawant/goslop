# Batch 4 — Identifier collision FP reduction (2026-08-02)

## Scope

MASTER.md root cause #4: name-based triggers without callee disambiguation for
`exec` / `open` / `execute`.

| Rule | Detector file | FP mechanism addressed |
| --- | --- | --- |
| BP-PY-12 | `internal/lang/python/detectors/bad_practices/rules_security.go` | Attribute `.exec(` (SQLModel/Qt/testcontainers) and `exec(` inside string literals |
| BP-PY-7 | `internal/lang/python/detectors/bad_practices/rules_core.go` | Attribute `.open(` (fitz/Image/HTTP/self), `def open(`, docstring prose |
| CWE-89 | `internal/lang/python/detectors/cwe/rules.go` | Adjacent string literals, `text("…")`, non-SQL `execute`, `def execute(` |
| CWE-94 | `internal/lang/python/detectors/cwe/rules_injection.go` | Defense-in-depth attribute `.exec` skip (findCalls already left-boundary rejects `.`) |

Audit checklists under `plans/fp-validations/reports/` were **not** modified.

## Detector changes

### BP-PY-12
- Scan `pytext.Mask(source)` so docstring/string `exec(` tokens cannot match.
- Skip attribute callees (`session.exec`, `app.exec`, `db.exec`) when the byte before `exec`/`eval`/`compile` is `.`.
- Bare builtin `exec(payload)` / `eval(user_code)` still report.

### BP-PY-7
- Scan masked lines; only bare builtin `open(` (not `.open(`).
- Skip `def open(` / `async def open(`.
- Real `f = open(path)` without `with` still reports.

### CWE-89
- Skip `def execute(` / `async def execute(` headers.
- Treat adjacent string literals and static `text("…")` / `sqlalchemy.text("…")` as static SQL (suppress).
- For dynamic first args, require SQL evidence: formatted SQL, local assignment of dynamic SQL (unmasked), or expression carrying SQL keywords/quotes. HTTP/`web_function` execute wrappers stay silent.
- Retained TPs: f-string/`%` SQL, ORM stmt assigned from f-string, execute-wrapper that rewrites `sql`.

### CWE-94
- Explicit skip when call start is preceded by `.` (locks attribute non-match alongside findCalls boundary).

## Fixtures (unique batch-4 names)

| Rule | Safe | Vulnerable |
| --- | --- | --- |
| BP-PY-12 | `tests/fixtures/python/bp/BP-PY-12-attr-method-exec-safe.txt` | `…-vulnerable.txt` |
| BP-PY-12 | `tests/fixtures/python/bp/BP-PY-12-string-literal-exec-safe.txt` | `…-vulnerable.txt` |
| BP-PY-7 | `tests/fixtures/python/bp/BP-PY-7-attr-method-open-safe.txt` | `…-vulnerable.txt` |
| BP-PY-7 | `tests/fixtures/python/bp/BP-PY-7-def-open-safe.txt` | `…-vulnerable.txt` |
| CWE-89 | `tests/fixtures/python/cwe/CWE-89-adjacent-literal-safe.txt` | `…-vulnerable.txt` |
| CWE-89 | `tests/fixtures/python/cwe/CWE-89-sqlalchemy-text-safe.txt` | `…-vulnerable.txt` |
| CWE-89 | `tests/fixtures/python/cwe/CWE-89-non-sql-execute-safe.txt` | `…-vulnerable.txt` |
| CWE-89 | `tests/fixtures/python/cwe/CWE-89-def-execute-safe.txt` | `…-vulnerable.txt` |
| CWE-94 | `tests/fixtures/python/cwe/CWE-94-attr-method-exec-safe.txt` | `…-vulnerable.txt` |

## Tests updated (append-only)

- `internal/lang/python/detectors/bad_practices/audit_variants_test.go` — BP-PY-7/12 batch-4 cases
- `internal/lang/python/detectors/bad_practices/scan_test.go` — `TestBPPY7` / `TestBPPY12`
- `internal/lang/python/detectors/cwe/audit_variants_test.go` — CWE-89/94 batch-4 cases
- Integration matrices auto-discover new `*-{safe,vulnerable}.txt` pairs

## Sampled FP evidence (scripts/`<repo>` paths)

| Repo / finding | Rule | Pattern |
| --- | --- | --- |
| onlymaps / MASTER | BP-PY-12 | `db.exec` / `oracledb.exec` |
| violit 157–164 | BP-PY-12 | `session.exec(stmt)` |
| pyauto-desktop 66, 74; pyhash-complete 11 | BP-PY-12 | Qt `app.exec()` / `editor.exec()` |
| WHEN-Language 24 | BP-PY-12 | `"exec() is not supported…"` string |
| pdf_oxide / pycaps / cylinder / safer | BP-PY-7 | `fitz.open`, `Image.open`, Client`.open`, docstring/`def open` |
| WeThePeople 2, 237, 248 | CWE-89 | adjacent literals; `text("…" )` + binds |
| FuncToWeb / httptap / enso | CWE-89 | non-SQL / `def execute` |

## Validation

| Command | Result |
| --- | --- |
| `go test ./internal/lang/python/detectors/cwe/... ./internal/lang/python/detectors/bad_practices/... -count=1` | **PASS** (full packages) |
| Batch-4 audit subtests (CWE-89/94 + BP-PY-7/12 cases above) | **PASS** |
| Prior CWE-89 ORM/execute-wrapper variants | **PASS** (non-regression) |
| `go test ./tests/integration/python -count=1` | Partial fail from **other parallel batches** (`CWE-117-constant-log`, `CWE-367-different-path`, `BP-PY-49-fingerprint-pin` safe fixtures still emit). Batch-4 pairs are clean when run alone. |
| `make test` / `make lint-all` | Not run (per assignment) |
| Audit checklist `[ ]` boxes | Unchanged |

## Completion record

| Finding family | Rule | Condition changed | Safe fixture | Vulnerable fixture | Tests |
| --- | --- | --- | --- | --- | --- |
| Attribute / string `exec` | BP-PY-12 | Mask + reject `.exec` | `BP-PY-12-attr-method-exec-safe`, `BP-PY-12-string-literal-exec-safe` | matching `-vulnerable` | scan + audit_variants + BP matrix discovery |
| Attribute / def / docstring `open` | BP-PY-7 | Mask + bare `open(` only + skip `def open` | `BP-PY-7-attr-method-open-safe`, `BP-PY-7-def-open-safe` | matching `-vulnerable` | same |
| Static / non-SQL `execute` | CWE-89 | Static adjacent/`text()`; SQL evidence; skip `def` | `CWE-89-adjacent-literal-safe`, `CWE-89-sqlalchemy-text-safe`, `CWE-89-non-sql-execute-safe`, `CWE-89-def-execute-safe` | matching `-vulnerable` | audit_variants + CWE matrix discovery |
| Attribute `.exec` | CWE-94 | Explicit `.` skip | `CWE-94-attr-method-exec-safe` | `CWE-94-attr-method-exec-vulnerable` | audit_variants |

## Remaining uncertainty

- Full integration matrix may stay red until peer batches finish CWE-117 / CWE-367 / BP-PY-49 guardrails.
- `Path.open` without `with` is no longer reported (MASTER: ignore attribute `.open` unless builtin). Builtin `open` remains covered.
- Corpus rescan of `real-repos/*` not re-run in this batch (focused fixture validation only).
