# False-positive audit: astroz

## Run metadata

```yaml
timestamp: 2026-08-02T07:24:02Z
repository: astroz
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/astroz
branch: main
commit: d558933ec3a9c9ee826eb8de665b6e5d229ebecb
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/astroz
chunk_path: scripts/astroz/chunks
function_context_path: scripts/astroz/findings/functions
```

## Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/astroz/chunks -context-dir scripts/astroz/findings/functions real-repos/astroz`
- Findings: `21`
- Chunks reviewed: `scripts/astroz/chunks/Chunk_1_21.txt`
- Function contexts reviewed: `scripts/astroz/findings/functions/1.txt` .. `scripts/astroz/findings/functions/21.txt` (all 21)

## Audit checklist

- [x] Read every assigned chunk under `scripts/astroz/chunks`.
- [x] Read `scripts/astroz/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 6 | 5, 6, 13, 14, 15, 21 |
| True positive | 15 | 1, 2, 3, 4, 7, 8, 9, 10, 11, 12, 16, 17, 18, 19, 20 |
| Uncertain | 0 | — |

## False positives

### [ ] Finding `5` — `CWE-772`

- Function context: `scripts/astroz/findings/functions/5.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/astroz/benchmarks/sgp4_compat_test.py:81:1`
- Checklist pattern: the URL response is not assigned to a variable — the assignment binds a string, so no unclosed resource handle exists

Source excerpt:

```
79:     url = "https://celestrak.org/NORAD/elements/gp.php?GROUP=active&FORMAT=tle"
80:     req = urllib.request.Request(url, headers={"User-Agent": "astroz-test"})
81:     tle_text = urllib.request.urlopen(req, timeout=60).read().decode("utf-8")
```

Why this is a false positive: The rule's condition is "a file, socket, or URL response is assigned without same-function close"; here the `urlopen` response is an intermediate expression consumed immediately by `.read()` — the assigned variable `tle_text` holds a `str`, not a resource, so no response handle is ever bound or leaked.

Checklist evidence: the regex trigger `\w+ = urllib\.request\.urlopen(` matched the assignment line, but the response object is never assigned to `tle_text`; only the decoded text is assigned, so the "resource is assigned without release" condition is not satisfied.

### [ ] Finding `6` — `BP-PY-42`

- Function context: `scripts/astroz/findings/functions/6.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/astroz/benchmarks/sgp4_compat_test.py:103:1`
- Checklist pattern: rule condition not satisfied — the try/except is tolerant data filtering, not an expect-failure assertion

Source excerpt:

```
101:     satrecs = []
102:     for line1, line2 in tle_pairs:
103:         try:
104:             sat = Satrec.twoline2rv(line1, line2, WGS72)
105:             if sat.error == 0:
106:                 satrecs.append(sat)
107:         except:
108:             pass
109:     num_sats = len(satrecs)
```

Why this is a false positive: The rule's condition is a test using try/except to *expect failure* instead of `assertRaises`/`pytest.raises`; here the loop silently skips malformed TLE records while building benchmark input, the function never asserts that an exception occurs, and the file contains no assertions at all.

Checklist evidence: the except suite is `pass` inside a data-loading loop — there is no failure expectation and no assertion intent, so the "instead of assertRaises/pytest.raises" clause of the rule condition is unmet.

### [ ] Finding `13` — `BP-PY-46`

- Function context: `scripts/astroz/findings/functions/13.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/astroz/bindings/python/astroz/__init__.py:598:9`
- Checklist pattern: `print(` token inside a docstring doctest example, not executable code

Source excerpt:

```
538:     """Screen a constellation for conjunction events.
...
587:     Examples
588:     --------
589:     Single-target screening (fastest):
590:
591:     >>> min_dist, min_t = screen("starlink", range(1440), threshold=5.0, target=0)
...
597:     >>> pairs, times = screen("starlink", range(1440), threshold=10.0)
598:     >>> print(f"Found {len(pairs)} conjunction events")
...
618:     """
```

Why this is a false positive: Line 598 is a doctest line inside the `screen()` docstring (opened at line 538, closed at line 618); the `print(` is documentation text, not an executable statement, so it cannot pollute a library consumer's stdout.

Checklist evidence: BP-PY-46's condition is `print` in executable non-script code; the flagged line is inside the multi-line `"""..."""` docstring, which the per-line string check fails to track across lines.

### [ ] Finding `14` — `BP-PY-46`

- Function context: `scripts/astroz/findings/functions/14.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/astroz/bindings/python/setup.py:49:9`
- Checklist pattern: `print` in a packaging script, not in a library module

Source excerpt:

```
15:     def build_extension(self, ext):
...
49:         print(f"Building with: {' '.join(cmd)}")
50:         subprocess.check_call(cmd, cwd=project_root)
```

Why this is a false positive: The rule targets `print` "used for operational logging in non-script modules"; `setup.py` is the packaging/build script executed by the toolchain, not a library module, and its prints are the script's own build-progress output — the same output class the rule's fix exempts for CLIs.

Checklist evidence: `setup.py` is a setup/packaging script whose prints are its operational output, so the "non-script modules" scope of the rule condition is not satisfied.

### [ ] Finding `15` — `BP-PY-46`

- Function context: `scripts/astroz/findings/functions/15.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/astroz/bindings/python/setup.py:90:9`
- Checklist pattern: `print` in a packaging script, not in a library module

Source excerpt:

```
86:         # Copy to target
87:         import shutil
88:
89:         target = ext_dir / f"_astroz{sysconfig.get_config_var('EXT_SUFFIX')}"
90:         print(f"Copying {built_lib} -> {target}")
91:         shutil.copy2(built_lib, target)
```

Why this is a false positive: Same as finding 14 — the print is inside the `ZigBuildExt.build_extension` method of the `setup.py` build script, reporting copy progress; `setup.py` is not a library module, so the rule's scope condition is not satisfied.

Checklist evidence: `setup.py` is a packaging script (not library code), so the "non-script modules" clause of the rule condition is unmet.

### [ ] Finding `21` — `PERF-PY-25`

- Function context: `scripts/astroz/findings/functions/21.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/astroz/examples/conjunction_screening.py:56:1`
- Checklist pattern: the lambda is constructed once per `sorted()` call, not once per loop element

Source excerpt:

```
54:     print(f"\n{len(pair_times)} unique satellite pairs with close approaches:")
55:     for (sat_i, sat_j), times_list in sorted(
56:         pair_times.items(), key=lambda x: len(x[1]), reverse=True
57:     )[:20]:
58:         print(f"  sat {sat_i:>5d} & {sat_j:>5d}: {len(times_list)} events")
```

Why this is a false positive: The rule's condition is a "heavy object or lambda constructed per homogeneous loop element"; here the lambda is a single argument of one `sorted()` call, evaluated exactly once when the loop header's iterable is built — the key function is then only *invoked* per element, not constructed per element.

Checklist evidence: the flagged line is the continuation of the `for` header's single `sorted(...)` call; the `lambda` object is created once before iteration, so the "constructed per homogeneous loop element" condition is not satisfied.

## True positives

### CWE-1121 — Excessive McCabe Cyclomatic Complexity (Findings 1, 2, 12)

The detector counts `if `/`elif `/`for `/`while `/`except ` tokens in the masked function body and fires at ≥ 12 (`minimumRouteBranches`). Each shown function genuinely exceeds the threshold:

| Finding id | Source | Reason |
| --- | --- | --- |
| 1 | `benchmarks/jax_gpu_bench.py:51` | `main()` contains 14 `for` loops plus `if num_gpus > 1:` — 15 counted branch tokens ≥ 12; the benchmark harness is genuinely deeply loop-structured. |
| 2 | `benchmarks/python_sgp4_bench.py:39` | `main()` contains 8 `for` loops plus 5 comprehension `for` clauses (`[i * step for i in range(points)]`, chunk splits) — 13 tokens ≥ 12. |
| 12 | `bindings/python/astroz/__init__.py:203` | `_omm_to_tle_pairs` contains the `for rec in data` loop, 2 `if` statements, and 9 ternary `if` expressions inside the checksum generators (`int(c) if c.isdigit() else ... for c in l1`) — 12 tokens, at threshold. |

### BP-PY-41 — pytest assert With Side Effects Only (Findings 3, 4, 11)

| Finding id | Source | Reason |
| --- | --- | --- |
| 3 | `benchmarks/sgp4_compat_test.py:18` | `test_single_satellite_api` only calls production code and `print`s results; it contains no `assert`, `pytest.raises`, or `self.assert` anywhere — side-effect-only test, matching the rule condition. |
| 4 | `benchmarks/sgp4_compat_test.py:66` | `test_batch_api_performance` loads TLEs and benchmarks `sat_array.sgp4` but never asserts; it prints PASS/FAIL and returns a bool — no assertion, matching the rule condition. |
| 11 | `benchmarks/sgp4_compat_test.py:165` | `test_output_validity` checks NaN/Inf via `if` and `print` only; no assertion statement exists in the body, matching the rule condition. |

### BP-PY-1 — Bare Except Clause (Findings 7, 16)

| Finding id | Source | Reason |
| --- | --- | --- |
| 7 | `benchmarks/sgp4_compat_test.py:107:1` | `except:` with a `pass` suite — a bare except that swallows every exception without handling or re-raise; the rule's condition is met. |
| 16 | `examples/cesium_fast.py:146:1` | `except Exception:` with a `pass` suite — broad except, no handling, no re-raise, not exempted as a test file; the rule's condition is met. |

### BP-PY-2 — Except Pass (Findings 8, 17)

| Finding id | Source | Reason |
| --- | --- | --- |
| 8 | `benchmarks/sgp4_compat_test.py:107:1` | Except suite is solely `pass` (suite length 1), exactly the rule's condition; failures are discarded silently. |
| 17 | `examples/cesium_fast.py:146:1` | Except suite is solely `pass`, exactly the rule's condition. |

### CWE-390 — Detection of Error Condition Without Action (Findings 9, 18)

| Finding id | Source | Reason |
| --- | --- | --- |
| 9 | `benchmarks/sgp4_compat_test.py:107:1` | `exceptPassStart` matches: except clause whose direct body is `pass` — error detected, no action taken. |
| 18 | `examples/cesium_fast.py:146:1` | Same construct: except clause whose direct body is `pass`, matching the rule condition. |

### CWE-396 — Declaration of Catch for Generic Exception (Finding 19)

| Finding id | Source | Reason |
| --- | --- | --- |
| 19 | `examples/cesium_fast.py:146:1` | `except Exception:` in a non-test module; the generic handler silently swallows all failures (no logging/re-raise), genuinely hiding distinct failure conditions — the rule's condition is met. |

### CWE-1071 — Empty Code Block (Findings 10, 20)

| Finding id | Source | Reason |
| --- | --- | --- |
| 10 | `benchmarks/sgp4_compat_test.py:107:9` | `except:\n    pass` matches the `pyTierBEmptyExceptRE` pattern (except clause containing only pass) — an empty exception handler. |
| 20 | `examples/cesium_fast.py:146:9` | `except Exception:\n    pass` matches the same empty-handler pattern. |

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none (single-reviewer audit)
- Chunk evidence: `scripts/astroz/chunks`
- Function evidence: `scripts/astroz/findings/functions`
- Validation: `git diff --check` — pass
