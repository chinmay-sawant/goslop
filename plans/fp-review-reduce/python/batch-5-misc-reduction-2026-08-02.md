# Batch 5 — misc MASTER pattern FP reduction (2026-08-02)

- Agent: Batch 5 of 5 (parallel)
- Skill: `plans/skills/review-and-reduce/SKILLS.md`
- Scope: CWE-117, CWE-1341, PERF-PY-25, PERF-PY-26, CWE-367, CWE-88, BP-PY-11/CWE-502 (ruamel), CWE-93, BP-PY-49, BP-PY-13, CWE-215; zero-FP confirmation repos
- Evidence roots: `scripts/<repo>/…`, `real-repos/<repo>/`
- Checklist boxes left unchanged

## Zero-FP repos (confirmed, no detector changes)

| Report | FP count | Action |
| --- | ---: | --- |
| among-llms.md | 0 | skipped / confirmed |
| graphzero.md | 0 | skipped / confirmed |
| pingram.md | 0 | skipped / confirmed |
| PyDepends.md | 0 | skipped / confirmed |
| python-injection.md | 0 | skipped / confirmed |
| Ai-copypaste-insult.md | 0 | skipped / confirmed |

## Root-cause fixes (#5–#11)

| Pattern | Mechanism | Guardrail |
| --- | --- | --- |
| CWE-117 | Constant / `len(...)` / numeric f-string interpolations treated as CRLF sinks | Require at least one CRLF-capable interpolation |
| CWE-1341 | Adjacent `.close()` on different receivers or lifecycle `__exit__`/`__del__` | Same receiver, same function body; skip lifecycle hooks |
| PERF-PY-25 | `key=lambda`, light attribute lambdas, construct-then-`return` | Skip sort keys, non-allocating bodies, early-return loop arms |
| PERF-PY-26 | `parse_*` + name substrings (`render`/`process`) / token loops | `parse_*` only on `handle_job`/`handle_request`; decode/Image/zlib still loop-sensitive |
| CWE-367 | `exists(A)` then `open(B)` / name collision with `simple_open` | Same path identifier between check and use (Go-side, no RE2 backrefs) |
| CWE-88 | Module constants / `str(KEY_FILE)` / bench harness argv | Trusted UPPER_SNAKE + literal concat; skip test/benchmark modules |
| BP-PY-11 / CWE-502 | `yaml = YAML(); yaml.load(f)` (ruamel) matched as PyYAML | Skip ruamel / `YAML()` loaders without `Loader=` |
| CWE-93 | `.headers[` inside `==` / asserts treated as writes | Assignment `=` only; skip test modules |
| CWE-215 | English word `password` inside message string | Sensitive identifier outside string literals |
| BP-PY-49 | Fingerprint pin `verify = False` / `!= CERT_NONE` | Skip local assigns + fingerprint window; skip comparisons; fix empty-seen fallback |
| BP-PY-13 | Bench/test placeholder secrets | Skip test/benchmark paths; expand placeholder prefixes (`bench-`, `testing-only`, …) |

## Completion record

| Finding samples (repo / id) | Rule | Detector change | Safe fixture | Vulnerable fixture |
| --- | --- | --- | --- | --- |
| CourtScrapper 102, 127, 298; logxide cluster | CWE-117 | `logMessageHasCRLFCapableValue` | `tests/fixtures/python/cwe/CWE-117-constant-log-safe.txt` | `…/CWE-117-constant-log-vulnerable.txt` |
| tenso 80/84/98/114/125; CourtScrapper teardown | CWE-1341 | `sameHandleDoubleCloseStart` | `…/CWE-1341-distinct-handles-safe.txt` | `…/CWE-1341-distinct-handles-vulnerable.txt` |
| Cronboard 40; pycaps 39–48; numeth 1 | PERF-PY-25 | key/light/early-return guards | `…/perf/PERF-PY-25-key-lambda-safe.txt`, `PERF-PY-25-early-return-safe.txt` | matching `*-vulnerable.txt` |
| WHEN-Language 48+; rendercv 45–55 | PERF-PY-26 | narrow hot markers; parse_* gate | `…/perf/PERF-PY-26-parser-descent-safe.txt` | `…/PERF-PY-26-parser-descent-vulnerable.txt` |
| safer 21/33/37; pycaps/movielite | CWE-367 | same-path TOCTOU | `…/cwe/CWE-367-different-path-safe.txt` | `…/CWE-367-different-path-vulnerable.txt` |
| Cronboard 62; wse 2; rendercv 15/19 | CWE-88 | trusted argv + bench skip | `…/cwe/CWE-88-constant-argv-safe.txt` | `…/CWE-88-constant-argv-vulnerable.txt` |
| rendercv 5/6/9–11/… | BP-PY-11, CWE-502 | ruamel `YAML().load` skip | `…/bp/BP-PY-11-ruamel-safe.txt`, `…/cwe/CWE-502-ruamel-safe.txt` | matching `*-vulnerable.txt` |
| FuncToWeb 44–57; tenso 118 | CWE-93 | assignment-only header writes | `…/cwe/CWE-93-header-read-safe.txt` | `…/CWE-93-header-read-vulnerable.txt` |
| Cronboard 30 | CWE-215 | sensitive ident outside literals | `…/cwe/CWE-215-literal-password-word-safe.txt` | `…/CWE-215-literal-password-word-vulnerable.txt` |
| niquests 21–35; httptap | BP-PY-49 | fingerprint / compare / fallback | `…/bp/BP-PY-49-fingerprint-pin-safe.txt` | `…/BP-PY-49-fingerprint-pin-vulnerable.txt` |
| wse 1/41/48/… | BP-PY-13 | bench path + placeholder expand | `…/bp/BP-PY-13-bench-secret-safe.txt` | `…/BP-PY-13-bench-secret-vulnerable.txt` |

## Code touchpoints

- `internal/lang/python/detectors/cwe/fp_guards.go` (new helpers)
- `internal/lang/python/detectors/cwe/rules_injection.go` (117/88/93)
- `internal/lang/python/detectors/cwe/rules_tier_b_quality.go` (1341)
- `internal/lang/python/detectors/cwe/rules_tier_b_resource.go` (367)
- `internal/lang/python/detectors/cwe/rules.go` (502)
- `internal/lang/python/detectors/cwe/rules_code_dynamic.go` (215)
- `internal/lang/python/detectors/perf/rules_hotpath.go` (25/26)
- `internal/lang/python/detectors/bad_practices/rules_security.go` (11/13)
- `internal/lang/python/detectors/bad_practices/rules_prod.go` (49)
- `internal/lang/python/detectors/bad_practices/common.go` (placeholder prefixes)
- Append-only cases in `cwe`/`perf`/`bad_practices` `audit_variants_test.go`

## Validation

```text
go test ./internal/lang/python/detectors/cwe/... \
        ./internal/lang/python/detectors/perf/... \
        ./internal/lang/python/detectors/bad_practices/... \
        ./tests/integration/python -count=1
```

**Result: pass** (2026-08-02).

- No inline Python in Go tests; all variants are text fixtures.
- Safe fixtures silent; vulnerable fixtures still fire for each rule.
- Audit report `[ ]` checkboxes not modified.
- Full `make test` / `make lint-all` / corpus rescan not run (per batch instructions).

## Remaining uncertainty

- CWE-88 cases whose dynamic segments are lower-case locals from `mkdtemp`/`getsockname` outside `benchmarks/` may still fire; CourtScrapper/pdf_oxide harnesses that are not under a bench path may need a follow-up trusted-path heuristic.
- PERF-PY-25 light-lambda skip may miss a true “lambda allocates via call” pattern that only closes over loop state without a ctor on the same line — retained TPs still use `WorkItem(...)` + `finalize(...)`.
- BP-PY-49 files that both implement fingerprint pinning and expose a public `verify=False` kwarg on the same source line window rely on nearby `assert_fingerprint`; kwarg-form `verify=False` remains reportable.
