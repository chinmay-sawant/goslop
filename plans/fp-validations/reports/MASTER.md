# False-positive audit — master report (42 real repos)

## Run metadata

```yaml
timestamp: 2026-08-02
audit_corpus: r/Python top showcase projects (2025-2026)
repositories: 42 clones under real-repos/ (41 GitHub + 1 GitLab: enso)
evidence_root: scripts/<name>/chunks + scripts/findings/functions
reports: plans/skills/false-positive-audit/reports/<name>.md
```

## Scan evidence

- Build command: `go build` (pre-built `./bin/goslop`)
- Scan command (per repo):

```
./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml \
  --export-context --export-chunks --no-cache \
  -chunks-dir scripts/<name>/chunks \
  -context-dir scripts/<name>/findings/functions real-repos/<name>
```

- Config: `templates/goslop-python.toml` (languages = ["python"], fail_on = "none")
- Findings: **6,964** across 41 repos (1 repo, polygon-screenshot-tool, has 0 findings — its only code file is `.pyw` which the scanner skips)
- Chunks + function contexts: per-project under `scripts/<name>/`
- Audit method: 41 parallel sub-agent audits, each reading every chunk, every function context for classified false positives, and the enclosing source; rule conditions verified via `./bin/goslop -explain <rule>` and detector sources.

## Classification summary (all 41 repos)

| Classification | Count | Repos (FP/TP/U) |
| --- | ---: | --- |
| False positive | **2,839 (40.8%)** | see table below |
| True positive | **4,125 (59.2%)** | see table below |
| Uncertain | **0 (0.0%)** | all 12 former uncertains reclassified (see note below) |

### Per-repo breakdown

| Repo | Findings | FP | TP | U | Dominant FP pattern |
| --- | ---: | ---: | ---: | ---: | --- |
| WeThePeople | 1,492 | 70 | 1,422 | 0 | CWE-89 static/parameterized SQL; CWE-909 `db` param; BP-PY-36 closed session |
| httpmorph | 714 | 278 | 436 | 0 | BP-PY-46 print in script modules (272) |
| pdf_oxide | 636 | 448 | 188 | 0 | BP-PY-46 print in scripts/examples (414) |
| logxide | 503 | 186 | 317 | 0 | BP-PY-46 scripts/benchmarks (105); BP-PY-1 handled errors (23) |
| Project_Parva | 412 | 401 | 11 | 0 | BP-PY-45 sys.path bootstrap (133); BP-PY-46 CLI prints (77) |
| CourtScrapper | 340 | 12 | 328 | 0 | CWE-117 constant/numeric log interpolations |
| niquests | 346 | 311 | 35 | 0 | except-pass family (132); BP-PY-41 test helpers (52) |
| caniscrape | 338 | 259 | 79 | 0 | BP-PY-46 CLI command output (255) |
| violit | 248 | 29 | 219 | 0 | CWE-1121 branch counting; BP-PY-12 `session.exec` (5) |
| pictex | 208 | 166 | 42 | 0 | BP-PY-41 `check_func` golden-file assertions (160) |
| wse | 206 | 49 | 157 | 0 | BP-PY-46 scripts+docstrings (28); BP-PY-13 test fixtures (9) |
| movielite | 154 | 96 | 58 | 0 | BP-PY-46 examples/ scripts (87) |
| enso | 138 | 11 | 127 | 0 | BP-PY-46 test-file detection gap (8) |
| tenso | 125 | 54 | 71 | 0 | BP-PY-46 examples/docstrings (38); CWE-1341 distinct handles (6) |
| httptap | 103 | 103 | 0 | 0 | BP-PY-46 docstring+CLI prints (30); BP-PY-41 benchmarks (28) |
| onlymaps | 99 | 63 | 36 | 0 | BP-PY-12 `exec` name collision (25); BP-PY-7 custom `open` (15) |
| pyauto-desktop | 92 | 4 | 88 | 0 | BP-PY-12 PyQt `.exec()`; PERF-PY-26 name substring |
| pycaps | 51 | 23 | 28 | 0 | PERF-PY-25 lambda closure semantics (10); BP-PY-7 non-builtin open (5) |
| rendercv | 73 | 46 | 27 | 0 | BP-PY-46 CLI (15); BP-PY-11/CWE-502 ruamel.yaml safe (10) |
| among-llms | 72 | 0 | 72 | 0 | none (BP-PY-47 ×54 all true) |
| pytogether | 71 | 1 | 70 | 0 | BP-PY-46 print inside string template |
| Cronboard | 69 | 4 | 65 | 0 | CWE-215 static word "password" in message |
| FuncToWeb | 57 | 19 | 38 | 0 | CWE-93 header reads as writes (4); CWE-89 `execute()` collision (3) |
| calgebra | 42 | 23 | 19 | 0 | BP-PY-1 handlers that re-raise/record (11) |
| voicetag | 42 | 29 | 13 | 0 | BP-PY-46 Typer `@app.command` CLI (29) |
| whatsapp-wrapped | 40 | 17 | 23 | 0 | BP-PY-46 `__main__`-guarded prints (15) |
| safer | 40 | 33 | 7 | 0 | BP-PY-7 open in docstrings (13); CWE-367 benign checks (3) |
| html2pic | 36 | 4 | 32 | 0 | CWE-396 re-raise with `from e` (4) |
| sync-with-uv | 20 | 15 | 5 | 0 | BP-PY-46 CLI output (14) |
| cylinder | 17 | 11 | 6 | 0 | BP-PY-7 `.open` HTTP method (4); BP-PY-6 framework flag (3) |
| requestSpeedTest | 13 | 4 | 9 | 0 | BP-PY-42 re-raise/error-accounting (2) |
| pyhash-complete | 11 | 1 | 10 | 0 | BP-PY-12 `app.exec()` Qt method |
| graphzero | 7 | 0 | 7 | 0 | none |
| pingram | 6 | 0 | 6 | 0 | none |
| numeth | 5 | 1 | 4 | 0 | PERF-PY-25 lambda in early-return loop |
| FlashySurf | 5 | 4 | 1 | 0 | BP-PY-46 module-level prints in scripts |
| PyDepends | 4 | 0 | 4 | 0 | none (except-pass on StopIteration teardown, structural) |
| python-injection | 2 | 0 | 2 | 0 | none |
| Ai-copypaste-insult | 2 | 0 | 2 | 0 | none |
| polygon-screenshot-tool | 0 | — | — | — | no Python files scanned (.pyw skipped) |

## Cross-repo false-positive patterns (root causes)

Ranked by findings count. Each is a detector weakness, not repo-specific noise.

### 1. BP-PY-46 — "print for logging in non-script modules" (~1,400 FP)

The single largest FP source. The detector has no script/library distinction beyond `__main__`-guard lexical checks, and no docstring masking.

- Script modules (examples/, scripts/, tools/, setup.py, CLI command bodies): pdf_oxide 414, httpmorph 272, caniscrape 255, movielite 87, Project_Parva 77, logxide 105, tenso 33, voicetag 29, wse 23, sync-with-uv 14, rendercv 15, WHEN-Language 15, whatsapp-wrapped 15, astroz 3, FlashySurf 4.
- Docstring/doctest/triple-quoted-string prints: safer 3, tenso 5, wse 5, httpmorph 1, pytogether 1, pictex 1, astroz 1, calgebra 1, violit 1.
- Missed CLI recognition: `@cli.command()` (caniscrape click, 252), `@app.command()` (voicetag Typer, 29), cyclopts `@app.default()` (sync-with-uv, 14), `main()` invoked under `__main__` (whatsapp-wrapped 15, wse 23).

**Fix:** script-module detection (dir names examples/scripts/tools/benchmarks, setup.py, CLI decorators, `__main__` guard including indented call to main) + docstring span masking.

### 2. BP-PY-41 — "test with no assertions" (~250 FP)

- Unrecognized assert forms: pytest-regressions `check_func(file_regression, ...)` (pictex 160), pytest-benchmark `benchmark(...)` fixtures (httptap 25), delegation to assert-bearing `_inner_*` helpers (niquests 52), `pytest.fail`/raises-based smoke tests.
- Body-scan bug: dedented lines inside triple-quoted strings abort the scan before the real asserts (whatsapp-wrapped 2, rendercv 3, safer 3).

**Fix:** recognize more assertion idioms + fix the string-literal indentation bug in the body scan.

### 3. Broad-except family — BP-PY-1, CWE-396, BP-PY-2, CWE-390, CWE-1071 (~350 FP)

Handlers that actually handle: re-raise (`raise`, `raise … from e`), record into result/error/health fields, log with `exc_info`, retry-and-re-raise, or set `future.set_exception`. Rule conditions require "swallows/hides failures" but the detectors don't inspect the suite for these actions.

- Biggest: niquests 48, logxide 32, calgebra 13, Project_Parva 45, onlymaps 8, html2pic 4, FuncToWeb 4, CourtScrapper 3, requestSpeedTest 2, tenso (BP-PY-10 pickle bench).

### 4. Identifier collisions — BP-PY-12/CWE-94 `exec`, BP-PY-7 `open`, CWE-89 `execute` (~150 FP)

Name-based triggers with no callee disambiguation:

- `exec`: onlymaps 25 (`db.exec` SQL method), violit 5 (`session.exec` SQLModel), pyauto-desktop 2 (PyQt `app.exec()`), pyhash-complete 1, WHEN-Language 1 (string literal).
- `open`: onlymaps 15 (`self.open()` connection method), pdf_oxide 10 (`fitz.open`), pycaps 5 (`Image.open`, `def open`), cylinder 4 (werkzeug `.open` HTTP), safer 19 (docstring matches + `def open`), logxide (os.open).
- `execute`: violit 4, FuncToWeb 3, WeThePeople 20 (parameterized SQL misread as dynamic), httptap 4 (HTTP protocol method), Project_Parva 3, enso 2, wse 1.

**Fix:** only flag the real builtins/known callees; ignore attribute calls unless receiver is the stdlib module; for CWE-89 require SQL context or data-flow into the string.

### 5. CWE-117 — log injection (~25 FP)

Fires on interpolations that cannot carry CRLF: ints, loop counters, `len(...)`, `"="*60` constants, internally generated names (CourtScrapper 5, logxide 16, cylinder 1, tenso).

**Fix:** only flag tainted/non-numeric interpolations.

### 6. CWE-1341 — double close (~20 FP)

Regex pairs adjacent `.close()` on *different* handles (`clip1.close(); clip2.close()`, `page/context/browser/playwright`, `ws1/ws2`), or idempotent/guarded closes in `__exit__`/`__del__`.

**Fix:** track handle identity; require same identifier.

### 7. PERF-PY-25 — lambda in loop (~15 FP)

Lambdas that close over per-iteration variables (semantically required, pycaps 10) or are built once per call inside an early-return loop (numeth 1, Cronboard 1, astroz 1).

### 8. PERF-PY-26 — "hot path" (~50 FP)

Triggered by name substrings, not loops: `render` in `rendercv_model`/`run_rendercv` (rendercv 5), recursive-descent `parse_*` steps consuming distinct tokens (WHEN-Language 35), one-shot CLIs (pictex 1, caniscrape 2, pyauto-desktop 1, WeThePeople 7).

### 9. CWE-367 — check-then-use (~6 FP)

Benign guards: test teardown `exists→remove` on own temp file, ValueError-only existence checks (safer 3, pycaps 2, movielite 1).

### 10. CWE-88 — argument injection (~10 FP)

argv segments interpolate only fixed constants/internal paths (`TLS_CERT_DIR`, mkdtemp paths, socket-assigned ports); untrusted values reach stdin not argv (wse 2, pdf_oxide 6, CourtScrapper 1, Cronboard 1, rendercv 2).

### 11. Misc singletons (~40 FP)

- BP-PY-13 hardcoded secrets: test/bench fixtures and "change-me" placeholders (wse 9, pycaps 1, calgebra 1).
- BP-PY-11/CWE-502 `yaml.load` on ruamel `YAML()` (rendercv 10) — ruamel is not PyYAML.
- CWE-93 header reads/`==` comparisons/numeric-only values as CRLF sinks (FuncToWeb 4, cylinder 2, tenso 2, violit 1).
- BP-PY-49 TLS markers where verification isn't disabled (niquests 11, httptap 7, wse).
- BP-PY-36/BPPY-45/BPPY-37/BPPY-32/BPPY-14 etc. — see per-repo reports.

## True positives (4,125)

Verified against rule conditions and source. Notable clusters:

- Bare/broad excepts that genuinely swallow: BP-PY-1 (WeThePeople 494, httpmorph 151, pdf_oxide 91, logxide 65, wse 59, violit 56…), CWE-396, except-pass family.
- BP-PY-47 eager f-string logging (among-llms 54, CourtScrapper 192, logxide 153) — eager evaluation before formatting.
- BP-PY-46 prints in actual library modules (httpmorph 163, WeThePeople, logxide 1 real hit in interceptor.py).
- Complexity/nesting heuristics CWE-1121/1124 (violit, niquests 10, pdf_oxide 16).
- CWE-829/CWE-94 dynamic imports over non-literal args (voicetag 2, niquests 22 — some FP; violit 131 `runpy.run_path` reclassified TP).
- CWE-1084 ≥3 `open`/`.execute` in one function (WeThePeople: 13 total after reclassifying 8 former uncertains).
- CWE-22 path join without confinement (pycaps 22).
- Missing timeouts on `session.get` (niquests 4, BP-PY-14), CWE-502 pickle of potentially untrusted data, CWE-295 CERT_NONE, CWE-489/756 in logxide tests.

## Uncertain findings (0)

All 12 former uncertains reclassified after verifying detector thresholds / rule conditions in source:

| Repo | ID(s) | Rule | New class | Rationale |
| --- | --- | --- | --- | --- |
| WeThePeople | 174, 285, 440, 451, 977, 1110, 1174, 1204 | CWE-1084 | **TP** | threshold is `>= 3` `open`/`.execute` calls (`detectCWE1084`); each function has 3–5 `.execute` |
| pycaps | 22 | CWE-22 | **TP** | `open(os.path.join(base, rule.filename))` without basename/resolve confinement |
| logxide | 286 | CWE-1121 | **FP** | threshold is 12 branches; `dictConfig` has 9 (`if`/`elif`/`for`) |
| pdf_oxide | 632 | CWE-88 | **FP** | local olmOCR bench harness; argv is `sys.executable`-class + internal corpus paths (same pattern as sibling CWE-88 FPs) |
| violit | 131 | CWE-829 | **TP** | `runpy.run_path(script_path)` with dynamic path matches `detectCWE829` |

## Final evidence

- Delegated reviewers: 41 parallel audit agents (one per repo), instructed per `plans/skills/false-positive-audit/AUDIT_PROMPT.md`
- Chunk evidence: `scripts/<name>/chunks/`
- Function evidence: `scripts/<name>/findings/functions/`
- Per-repo reports: `plans/skills/false-positive-audit/reports/<name>.md` (41 files, ~59,400 lines)
- Validation: `git diff --check` — **pass** (all reports + this master)
