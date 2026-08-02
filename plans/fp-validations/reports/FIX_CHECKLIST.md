# Round-2 fix checklist (aggregated FP patterns, latest binary)

Source: post-fix v2 audits (28 repos, appended to each report). Every excess finding matched a prior audited FP by Source — these are re-appeared FPs after the local detector changes removed guardrails. Fix per rule family, following `plans/skills/review-and-reduce/SKILLS.md`. Preserve the audited TP lists (per-repo reports) — do not reintroduce over-suppression.

## F1. BP-PY-46 print (≈260 findings) — rules_observability.go

| Repo | Count | Trigger shape |
| --- | ---: | --- |
| movielite | 87 | prints in `examples/` runnable demo scripts |
| pdf_oxide | 65 | prints in `scripts/`, `examples/`, `tools/` CLIs |
| tenso | 32 | prints in `examples/*.py` + `_make_seeds.py` script |
| wse | 23 | prints in `examples/standalone_*.py` invoked from `main()` under `__main__` guard |
| rendercv | 11 | CLI presentation (typer callbacks, rich panels) |
| caniscrape | 9 | telemetry command impl prints (click `@cli.command`) |
| enso | 8+1 | test-file exemption misses `Test.py`/`Tests/` (case); `def print(...)` header matched as call |
| FlashySurf | 4 | module-scope prints in standalone scripts, no `__main__` guard, no importable API |
| Project_Parva / logxide / others | few | CLI/script module prints |
| sync-with-uv | 1 | module-scope `print(file=fd)` in docs script `gen_ref_pages.py` |

Narrow guardrails: script-module detection (path components examples/scripts/tools/demos + `__main__`-guarded invocation + CLI decorators click/typer/cyclopts + setup.py), docstring masking, `def print` definition skip, test-file exemption case-insensitive (`Test.py`, `Tests/`). Do NOT suppress: library-module prints (audited TPs), print inside string templates handled already.

## F2. Broad-except family BP-PY-1 / CWE-396 / BP-PY-2 / CWE-390 / CWE-1071 (≈180) — rules_core.go, cwe/rules_platform.go, fp_guards.go

| Pattern | Repos (count) |
| --- | --- |
| Optional-dependency `ImportError` fallback (`except ImportError: pass/fallback`) | niquests 56, tenso, graphzero |
| Parsing fallback with defined fallback state (return default) | niquests 29 |
| Generic-catch propagation: re-raise wrapped/chained (`raise X from e`), record into result/error field, `set_exception`, log + continue | niquests 12+6, calgebra (re-raise, error-result, stderr+continue), html2pic 4 (`raise DomainError from e`), cylinder (unconditional re-raise), onlymaps (re-raise, recorded), safer, requestSpeedTest, sync-with-uv, FuncToWeb, logxide, httptap, CourtScrapper |
| `pass` fall-through to explicit `raise`; `pass` in test expected-exception assert | calgebra (pass→raise, test assert), tenso |
| Test-result recording then assert | onlymaps, httptap |

Caution: audits CONFLICT on `raise X from e` (voicetag 7 audited TP; html2pic 4 audited FP). Default: if the raised exception type is more specific than caught (subclass) AND chains with `from e`, exempt; if it is the SAME type or the audit marked TP, keep firing. ImportError-fallback: exempt only for `ImportError`/`ModuleNotFoundError` + fallback assignment in suite.

## F3. Identifier collisions BP-PY-12/CWE-94 exec/eval, CWE-89 execute (≈80) — rules_core.go, rules_security.go, fp_guards.go

| Pattern | Repos (count) |
| --- | --- |
| Attribute calls `.exec(...)` / `.execute(...)` on non-stdlib receivers | onlymaps 25, violit, tenso (pyodide JS bridge 51!), enso, wse, FuncToWeb, Project_Parva |
| `def exec(...)`/`def execute(...)` definitions | onlymaps, enso, tenso, violit |
| `__import__`/importlib over developer-controlled constants (own-pkg enumeration, hardcoded allowlist) | niquests 22, rendercv, logxide, Project_Parva |
| Static/bound SQL misread as dynamic (text()/parameterized) | violit 3, onlymaps 1, Project_Parva |
| pyodide `import js; js.<x>.exec` JS-bridge — evaluate only actual Python eval/exec builtins | niquests 51 |

## F4. CWE-117 log injection (≈20) — rules_security.go / fp_guards.go

Fires on non-CRLF-capable interpolations: ints, loop counters, `len(...)`, `"="*60`, internally generated names. Repos: logxide 16, CourtScrapper 2, cylinder 1, tenso 1.

## F5. CWE-88 argument injection (≈15) — fp_guards.go

argv f-strings interpolating only fixed constants/internal paths (sys.executable, mkdtemp, TLS_CERT_DIR, socket-assigned ports). Repos: pdf_oxide 3, Project_Parva, logxide, CourtScrapper, rendercv.

## F6. PERF-PY-28 executor lifetime (≈7) — perf rules

Executor created once per run/process (recommended pattern) flagged as per-unit. Repos: pdf_oxide 3, tenso 2, CourtScrapper 1, Project_Parva.

## F7. PERF-PY-27 distinct-path loads (≈6) — perf rules

Each loop iteration loads a distinct path derived from the loop element. Repos: pdf_oxide 4, rendercv, Project_Parva.

## F8. CWE-1341 double-close on different/idempotent handles (≈4) — fp_guards.go

Adjacent `.close()` on different handles or guarded idempotent closes. Repos: tenso 2, httpmorph 1, movielite (already fixed 7), pyauto-desktop.

## F9. BP-PY-40 daemon threads (≈14) — rules_prod.go

Threads constructed with explicit `daemon=True`; `.start()` on next line trips line-scoped check; threads with shutdown/wait protocol. Repos: wse 7, violit 1, logxide, rendercv, FuncToWeb, Project_Parva.

## F10. BP-PY-13 hardcoded secrets (≈8) — rules_security.go

`change-me` placeholder prefixes and self-describing test/bench fixture secrets. Repos: wse 5, pycaps 1, calgebra 1, violit 1.

## F11. Misc singletons (≈40) — various

CWE-367 benign TOCTOU (movielite 1, pycaps 1); CWE-478 exhaustive match (pycaps 2, onlymaps 1); CWE-1121 branch counting (violit 7, logxide); CWE-117/1046/1333/695/260/256/459/772/779/186/409/22, BP-PY-7 (safer own-open, pictex `open(...).read()` one-liner), BP-PY-6 (cylinder framework flag 3, pdf_oxide 1), BP-PY-44 (WHEN-Language 6, first-party `parser` module name), CWE-208 (WHEN-Language, Project_Parva), CWE-489/756 (logxide tests), BP-PY-37 (violit 3, Project_Parva), BP-PY-14 (violit 3, httpmorph 1), BP-PY-42 (wse 1, logxide, astroz, httptap), BP-PY-49 (httptap, niquests 1), BP-PY-32 (violit, FuncToWeb, Project_Parva), BP-PY-36 (niquests 1), CWE-93 (tenso 1, violit 1, cylinder 1, httptap), CWE-829/CWE-94 (niquests 22 own-pkg, logxide, rendercv, Project_Parva).

## Over-suppression residual (must NOT regress; these are currently below TP)

WeThePeople −14 (BP-PY-45/BP-PY-1/BP-PY-46), pyauto-desktop −2, pytogether −2, graphzero −1, python-injection −1. Fixes for F1/F2 must not suppress these sources (all listed in the per-repo over-suppression appends).

## Validation gate

After fixes: `make build`, rescan, per-repo found count must land between TP and TP+small residual (target ≤ 5% over TP per repo, 0 under).
