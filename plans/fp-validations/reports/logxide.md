# False-positive audit: logxide

## Run metadata

```yaml
timestamp: 2026-08-02T07:58:24Z
repository: logxide
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide
branch: main
commit: 136f7a4c3bc593488cd1e2c62bd74956265533d6
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide
chunk_path: scripts/logxide/chunks
function_context_path: scripts/logxide/findings/functions
```

## Scan evidence

- Build command: `n/a` (scan artifacts pre-generated; no rebuild performed)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/logxide/chunks -context-dir scripts/logxide/findings/functions real-repos/logxide`
- Findings: `503`
- Chunks reviewed: `scripts/logxide/chunks/Chunk_1_25.txt` .. `scripts/logxide/chunks/Chunk_501_503.txt` (all 21 chunk files)
- Function contexts reviewed: `scripts/logxide/findings/functions/<id>.txt` for every proposed false positive (chunk `Context:` excerpt for all 503; individual context files and enclosing sources followed up where the excerpt was insufficient)

## Audit checklist

- [x] Read every assigned chunk under `scripts/logxide/chunks`.
- [x] Read `scripts/logxide/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 186 | 2, 4, 5, 14, 16, 18, 19, 21, 22, 25, 26, 30, 31, 33, 52, 57, 60, 115, 124, 129, 146, 151, 154, 168, 169, 170, 171, 173, 174, 175, 177, 179, 180, 182, 183, 185, 186, 188, 189, 191, 192, 194, 196, 197, 199, 201, 203, 204, 215, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 230, 231, 232, 234, 235, 236, 237, 238, 239, 242, 243, 244, 247, 248, 249, 250, 251, 252, 253, 254, 255, 256, 260, 261, 263, 270, 279, 284, 285, 286, 287, 288, 289, 290, 291, 292, 293, 294, 295, 297, 298, 303, 304, 307, 308, 309, 312, 313, 314, 315, 316, 317, 318, 319, 320, 321, 322, 323, 324, 325, 326, 327, 328, 329, 330, 331, 332, 333, 334, 335, 336, 337, 341, 342, 343, 346, 348, 349, 351, 353, 354, 355, 357, 359, 361, 362, 364, 366, 368, 370, 372, 373, 375, 377, 382, 397, 405, 411, 412, 414, 416, 417, 418, 419, 420, 423, 426, 427, 431, 439, 441, 442, 446, 457, 460, 461, 467, 490, 491, 492, 498, 500, 501, 502, 503 |
| True positive | 317 | 1, 3, 6, 7, 8, 9, 10, 11, 12, 13, 15, 17, 20, 23, 24, 27, 28, 29, 32, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 53, 54, 55, 56, 58, 59, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 116, 117, 118, 119, 120, 121, 122, 123, 125, 126, 127, 128, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 147, 148, 149, 150, 152, 153, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 172, 176, 178, 181, 184, 187, 190, 193, 195, 198, 200, 202, 205, 206, 207, 208, 209, 210, 211, 212, 213, 214, 216, 217, 233, 240, 241, 245, 246, 257, 258, 259, 262, 264, 265, 266, 267, 268, 269, 271, 272, 273, 274, 275, 276, 277, 278, 280, 281, 282, 283, 296, 299, 300, 301, 302, 305, 306, 310, 311, 338, 339, 340, 344, 345, 347, 350, 352, 356, 358, 360, 363, 365, 367, 369, 371, 374, 376, 378, 379, 380, 381, 383, 384, 385, 386, 387, 388, 389, 390, 391, 392, 393, 394, 395, 396, 398, 399, 400, 401, 402, 403, 404, 406, 407, 408, 409, 410, 413, 415, 421, 422, 424, 425, 428, 429, 430, 432, 433, 434, 435, 436, 437, 438, 440, 443, 444, 445, 447, 448, 449, 450, 451, 452, 453, 454, 455, 456, 458, 459, 462, 463, 464, 465, 466, 468, 469, 470, 471, 472, 473, 474, 475, 476, 477, 478, 479, 480, 481, 482, 483, 484, 485, 486, 487, 488, 489, 493, 494, 495, 496, 497, 499 |
| Uncertain | 0 | — |

## False positives

One subsection per finding; the excerpt is the smallest snippet from the function-context file or referenced source that proves the decision.

### [ ] Finding 2 — BP-PY-7

- Function context: `scripts/logxide/findings/functions/2.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/_bench_common.py:237:25`
- Checklist pattern: `os.open()` is not the `open()` builtin

Source excerpt:

```
        235: 
        236:     def __enter__(self) -> RedirectedFD:
    >   237:         self._capfd = os.open(
        238:             self.capture_path, os.O_WRONLY | os.O_CREAT | os.O_TRUNC, 0o644
```

Why this is a false positive: `os.open()` is the low-level fd API, not the `open()` builtin, and never supports `with`; the fd is owned by a context-manager class (`RedirectedFD.__enter__`) and closed by its `__exit__`.

Checklist evidence: the flagged call is `os.open` (returns a raw fd), not the builtin `open()`.

### [ ] Finding 4 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/4.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/basic_handlers_benchmark.py:117:1`
- Checklist pattern: catch-all converts the failure into a recorded benchmark result

Source excerpt:

```
        115:         result.skipped = True
        116:         result.error = f"missing library: {e}"
    >   117:     except Exception as e:  # noqa: BLE001
        118:         import traceback
```

Why this is a false positive: the handler body handles the exception by recording it into the benchmark `result` (marked `# noqa: BLE001`): `result.ok = False` and the traceback are stored in the result payload — the failure is not swallowed, so the rule's “without handling” condition is not met.

Checklist evidence: the handler body assigns `result.ok`/`result.error`, explicitly modeling the failure outcome.

### [ ] Finding 5 — CWE-396

- Function context: `scripts/logxide/findings/functions/5.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/basic_handlers_benchmark.py:117:1`
- Checklist pattern: catch-all converts the failure into a recorded benchmark result

Source excerpt:

```
        115:         result.skipped = True
        116:         result.error = f"missing library: {e}"
    >   117:     except Exception as e:  # noqa: BLE001
        118:         import traceback
```

Why this is a false positive: same handler as the paired BP-PY-1 finding: the exception is deliberately caught and captured into the benchmark result structure, so the failure is processed rather than hidden.

Checklist evidence: the generic catch exists only to record any benchmark failure into the result.

### [ ] Finding 14 — CWE-117

- Function context: `scripts/logxide/findings/functions/14.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/compare_loggers.py:138:26`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
        136:         return lambda i: logger.info("Simple log message")
        137:     if scenario == "structured":
    >   138:         return lambda i: logger.info(
        139:             f"User action - user_id: {i}, action: login, status: success"
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 16 — CWE-88

- Function context: `scripts/logxide/findings/functions/16.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/compare_loggers.py:162:12`
- Checklist pattern: dynamic value bound as named-option argument in a no-shell argv

Source excerpt:

```
        160:     env = dict(os.environ)
        161:     env["PYTHONUTF8"] = "1"
    >   162:     proc = subprocess.run(
        163:         [
```

Why this is a false positive: the dynamic values are passed as separate argv elements consumed as the values of their named options (`--library`, `--scenario`, `-n`); they are delimited arguments in a no-shell argv list and cannot be parsed as unintended switches.

Checklist evidence: each value follows its named flag as its own argv element; no shell parsing is involved.

### [ ] Finding 18 — CWE-88

- Function context: `scripts/logxide/findings/functions/18.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/gil_benchmark.py:90:12`
- Checklist pattern: dynamic value bound as named-option argument in a no-shell argv

Source excerpt:

```
         88:     env = dict(os.environ)
         89:     env["PYTHONUTF8"] = "1"
    >    90:     proc = subprocess.run(
         91:         [sys.executable, THIS, "--worker", "--library", library, "-n", str(n)],
```

Why this is a false positive: the dynamic values are passed as separate argv elements consumed as the values of their named options (`--library`, `--scenario`, `-n`); they are delimited arguments in a no-shell argv list and cannot be parsed as unintended switches.

Checklist evidence: each value follows its named flag as its own argv element; no shell parsing is involved.

### [ ] Finding 19 — PERF-PY-26

- Function context: `scripts/logxide/findings/functions/19.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/gil_benchmark.py:109:1`
- Checklist pattern: `argparse.parse_args()` runs once at process start, not on a hot path

Source excerpt:

```
        107:     parser.add_argument("--worker", action="store_true", help=argparse.SUPPRESS)
        108:     parser.add_argument("--library", default="", help=argparse.SUPPRESS)
    >   109:     args = parser.parse_args()
        110:
```

Why this is a false positive: `parser.parse_args()` runs a single time at program entry before the benchmark loop; there is no hot path and no parse to cache.

Checklist evidence: the flagged parse executes once per process at startup.

### [ ] Finding 21 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/21.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/perf_micro.py:79:1`
- Checklist pattern: catch-all converts the failure into a recorded benchmark result

Source excerpt:

```
         77:     try:
         78:         _dispatch(result, scenario, n, warmup, args.runs, args.threads)
    >    79:     except Exception as e:  # noqa: BLE001
         80:         import traceback
```

Why this is a false positive: the handler body handles the exception by recording it into the benchmark `result` (marked `# noqa: BLE001`): `result.ok = False` and the traceback are stored in the result payload — the failure is not swallowed, so the rule's “without handling” condition is not met.

Checklist evidence: the handler body assigns `result.ok`/`result.error`, explicitly modeling the failure outcome.

### [ ] Finding 22 — CWE-396

- Function context: `scripts/logxide/findings/functions/22.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/perf_micro.py:79:1`
- Checklist pattern: catch-all converts the failure into a recorded benchmark result

Source excerpt:

```
         77:     try:
         78:         _dispatch(result, scenario, n, warmup, args.runs, args.threads)
    >    79:     except Exception as e:  # noqa: BLE001
         80:         import traceback
```

Why this is a false positive: same handler as the paired BP-PY-1 finding: the exception is deliberately caught and captured into the benchmark result structure, so the failure is processed rather than hidden.

Checklist evidence: the generic catch exists only to record any benchmark failure into the result.

### [ ] Finding 25 — CWE-117

- Function context: `scripts/logxide/findings/functions/25.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/perf_vs_stdlib.py:73:26`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
         71:         call = lambda i: logger.info("Simple log message")  # noqa: E731
         72:     elif scenario == "structured":
    >    73:         call = lambda i: logger.info(f"User action - user_id: {i}, action: login")  # noqa: E731
         74:     else:
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 26 — CWE-88

- Function context: `scripts/logxide/findings/functions/26.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/perf_vs_stdlib.py:90:12`
- Checklist pattern: dynamic value bound as named-option argument in a no-shell argv

Source excerpt:

```
         88:     env = dict(os.environ)
         89:     env["PYTHONUTF8"] = "1"
    >    90:     proc = subprocess.run(
         91:         [
```

Why this is a false positive: the dynamic values are passed as separate argv elements consumed as the values of their named options (`--library`, `--scenario`, `-n`); they are delimited arguments in a no-shell argv list and cannot be parsed as unintended switches.

Checklist evidence: each value follows its named flag as its own argv element; no shell parsing is involved.

### [ ] Finding 30 — CWE-88

- Function context: `scripts/logxide/findings/functions/30.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/real_handlers_comparison.py:154:12`
- Checklist pattern: dynamic value bound as named-option argument in a no-shell argv

Source excerpt:

```
        152:     env = dict(os.environ)
        153:     env["PYTHONUTF8"] = "1"
    >   154:     proc = subprocess.run(
        155:         [
```

Why this is a false positive: the dynamic values are passed as separate argv elements consumed as the values of their named options (`--library`, `--scenario`, `-n`); they are delimited arguments in a no-shell argv list and cannot be parsed as unintended switches.

Checklist evidence: each value follows its named flag as its own argv element; no shell parsing is involved.

### [ ] Finding 31 — PERF-PY-26

- Function context: `scripts/logxide/findings/functions/31.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/real_handlers_comparison.py:184:1`
- Checklist pattern: `argparse.parse_args()` runs once at process start, not on a hot path

Source excerpt:

```
        182:     parser.add_argument("--library", default="", help=argparse.SUPPRESS)
        183:     parser.add_argument("--handler", default="", help=argparse.SUPPRESS)
    >   184:     args = parser.parse_args()
        185:
```

Why this is a false positive: `parser.parse_args()` runs a single time at program entry before the benchmark loop; there is no hot path and no parse to cache.

Checklist evidence: the flagged parse executes once per process at startup.

### [ ] Finding 33 — CWE-117

- Function context: `scripts/logxide/findings/functions/33.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/basic_usage.py:64:5`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
         62:     # 5. String formatting
         63:     print("\n5. String Formatting:")
    >    64:     logger.info(f"User {'alice'} logged in from {'192.168.1.100'}")
         65:     logger.warning(f"High memory usage: {85}% ({1024} MB)")
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 52 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/52.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:317:1`
- Checklist pattern: catch-all converts the failure into an HTTP error response

Source excerpt:

```
        315:             return JsonResponse(user.to_dict(), status=201)
        316: 
    >   317:         except Exception as e:
        318:             app_logger.error(f"Failed to create user: {str(e)}")
```

Why this is a false positive: the handler body handles the exception by converting it into an HTTP 500 error response after logging — the failure outcome is modeled, so the “without handling” condition fails.

Checklist evidence: the handler returns an error response to the caller after logging.

### [ ] Finding 57 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/57.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:347:1`
- Checklist pattern: health-check catch-all updates service status

Source excerpt:

```
        345:             cursor.execute("SELECT 1")
        346:         db_status = "healthy"
    >   347:     except Exception as e:
        348:         db_logger.error(f"Database health check failed: {str(e)}")
```

Why this is a false positive: the catch-all is the deliberate health-check contract: any failure is handled by logging and reporting `db_status = "unhealthy"` — the failure is processed, not swallowed.

Checklist evidence: the exception is processed into the health-check status result.

### [ ] Finding 60 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/60.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:405:1`
- Checklist pattern: catch-all converts the failure into an HTTP error response

Source excerpt:

```
        403:         return JsonResponse(metrics_data)
        404: 
    >   405:     except Exception as e:
        406:         app_logger.error(f"Failed to retrieve metrics: {str(e)}")
```

Why this is a false positive: the handler body handles the exception by converting it into an HTTP 500 error response after logging — the failure outcome is modeled, so the “without handling” condition fails.

Checklist evidence: the handler returns an error response to the caller after logging.

### [ ] Finding 115 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/115.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:549:1`
- Checklist pattern: health-check catch-all updates service status

Source excerpt:

```
        547:         db.execute("SELECT 1")
        548:         db_status = "healthy"
    >   549:     except Exception as e:
        550:         db_logger.error(f"Database health check failed: {str(e)}")
```

Why this is a false positive: the catch-all is the deliberate health-check contract: any failure is handled by logging and reporting `db_status = "unhealthy"` — the failure is processed, not swallowed.

Checklist evidence: the exception is processed into the health-check status result.

### [ ] Finding 124 — CWE-1046

- Function context: `scripts/logxide/findings/functions/124.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:656:1`
- Checklist pattern: flagged line is an integer increment, not string concatenation

Source excerpt:

```
        654:     while time.time() - start_time < duration:
        655:         performance_logger.info(f"Stress test message {message_count + 1}")
    >   656:         message_count += 1
        657:
```

Why this is a false positive: the flagged construct `message_count += 1` is an integer counter increment; no immutable text is concatenated, so the rule condition (“creation of immutable text using string concatenation”) is not met.

Checklist evidence: the flagged expression is `int += 1`, not `str += str`.

### [ ] Finding 129 — CWE-117

- Function context: `scripts/logxide/findings/functions/129.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_demo.py:69:9`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
         67:     # Log many messages using standard logging
         68:     for i in range(count):
    >    69:         logger.info(f"Performance test message {i + 1}")
         70:
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 146 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/146.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:219:1`
- Checklist pattern: catch-all rolls back the transaction and returns an error response

Source excerpt:

```
        217:             return jsonify(user.to_dict()), 201
        218: 
    >   219:         except Exception as e:
        220:             db.session.rollback()
```

Why this is a false positive: the handler body handles the failure with `db.session.rollback()` plus an error response — the transaction outcome is explicitly handled.

Checklist evidence: the exception triggers `db.session.rollback()` and a 500 response.

### [ ] Finding 151 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/151.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:249:1`
- Checklist pattern: health-check catch-all updates service status

Source excerpt:

```
        247:         db.session.execute("SELECT 1")
        248:         db_status = "healthy"
    >   249:     except Exception as e:
        250:         db_logger.error(f"Database health check failed: {str(e)}")
```

Why this is a false positive: the catch-all is the deliberate health-check contract: any failure is handled by logging and reporting `db_status = "unhealthy"` — the failure is processed, not swallowed.

Checklist evidence: the exception is processed into the health-check status result.

### [ ] Finding 154 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/154.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:304:1`
- Checklist pattern: catch-all converts the failure into an HTTP error response

Source excerpt:

```
        302:         return jsonify(metrics)
        303: 
    >   304:     except Exception as e:
        305:         error_logger.error(f"Failed to retrieve metrics: {str(e)}")
```

Why this is a false positive: the handler body handles the exception by converting it into an HTTP 500 error response after logging — the failure outcome is modeled, so the “without handling” condition fails.

Checklist evidence: the handler returns an error response to the caller after logging.

### [ ] Finding 168 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/168.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/format_detailed.py:3:1`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
          1: from logxide import logging
          2: 
    >     3: print("=== Detailed Format with Thread Info ===")
          4: detailed_format = (
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 169 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/169.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/format_json.py:3:1`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
          1: from logxide import logging
          2: 
    >     3: print("=== JSON-like Structured Format ===")
          4: json_format = '{"timestamp":"%(asctime)s","level":"%(levelname)s","logger":"%(name)s","thread":%(thread)d,"process":%(process)d,"message":"%(message)s"}'
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 170 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/170.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/format_minimal.py:3:1`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
          1: from logxide import logging
          2: 
    >     3: print("=== Minimal Format ===")
          4: logging.basicConfig(format="%(message)s")
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 171 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/171.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/format_production.py:3:1`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
          1: from logxide import logging
          2: 
    >     3: print("=== Production Format ===")
          4: prod_format = "%(asctime)s [%(process)d:%(thread)d] %(levelname)s %(name)s: %(message)s"
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 173 — CWE-117

- Function context: `scripts/logxide/findings/functions/173.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/format_production.py:19:5`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
         17:     logger = logging.getLogger(component)
         18:     logger.setLevel(logging.INFO)
    >    19:     logger.info(f"{component} initialized successfully")
         20:     if component == "db.connection":
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 174 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/174.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/format_simple.py:3:1`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
          1: from logxide import logging
          2: 
    >     3: print("=== Simple Format Example ===")
          4: logging.basicConfig(format="%(levelname)s: %(name)s - %(message)s")
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 175 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/175.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/format_threaded.py:6:1`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
          4: from logxide import logging
          5: 
    >     6: print("=== Multi-threaded Logging ===")
          7: thread_format = (
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 177 — CWE-117

- Function context: `scripts/logxide/findings/functions/177.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/format_threaded.py:21:9`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
         19: 
         20:     for i in range(3):
    >    21:         logger.info(f"Processing task {i + 1}")
         22:         time.sleep(0.1)  # Simulate work
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 179 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/179.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:15:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         13: def test_default_format():
         14:     """Test 1: Default Python logging format"""
    >    15:     print("=== Test 1: Default Format ===")
         16:     logging.basicConfig()
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 180 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/180.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:35:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         33:     root_logger.info("Root logger message.")
         34:     logging.flush()  # Ensure all log messages are processed
    >    35:     print()
         36:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 182 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/182.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:40:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         38: def test_simple_format():
         39:     """Test 2: Simple format"""
    >    40:     print("=== Test 2: Simple Format ===")
         41:     logging.basicConfig(format="%(levelname)s: %(name)s - %(message)s")
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 183 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/183.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:49:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         47:     logger.error("Error in simple format")
         48:     logging.flush()  # Ensure all log messages are processed
    >    49:     print()
         50:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 185 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/185.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:54:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         52: def test_detailed_format():
         53:     """Test 3: Detailed format with timestamp and thread info"""
    >    54:     print("=== Test 3: Detailed Format with Thread Info ===")
         55:     detailed_format = (
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 186 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/186.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:67:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         65:     logger.error("Error message with detailed format")
         66:     logging.flush()  # Ensure all log messages are processed
    >    67:     print()
         68:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 188 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/188.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:72:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         70: def test_json_format():
         71:     """Test 4: JSON-like structured format"""
    >    72:     print("=== Test 4: JSON-like Structured Format ===")
         73:     json_format = '{"timestamp":"%(asctime)s","level":"%(levelname)s","logger":"%(name)s","thread":%(thread)d,"process":%(process)d,"message":"%(message)s"}'
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 189 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/189.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:82:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         80:     logger.error("Database connection timeout")
         81:     logging.flush()  # Ensure all log messages are processed
    >    82:     print()
         83:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 191 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/191.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:87:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         85: def test_debug_format():
         86:     """Test 5: Development/Debug format with all available fields"""
    >    87:     print("=== Test 5: Development/Debug Format ===")
         88:     debug_format = (
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 192 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/192.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:101:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         99:     logger.critical("Critical system failure")
        100:     logging.flush()  # Ensure all log messages are processed
    >   101:     print()
        102:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 194 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/194.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:106:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        104: def test_multithreaded_format():
        105:     """Test 6: Multi-threaded logging with thread names"""
    >   106:     print("=== Test 6: Multi-threaded Logging ===")
        107:     thread_format = (
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 196 — CWE-117

- Function context: `scripts/logxide/findings/functions/196.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:120:13`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
        118: 
        119:         for i in range(3):
    >   120:             logger.info(f"Processing task {i + 1}")
        121:             time.sleep(0.1)  # Simulate work
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 197 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/197.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:146:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        144:     main_logger.info("All worker threads completed")
        145:     logging.flush()  # Ensure all log messages are processed
    >   146:     print()
        147:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 199 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/199.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:151:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        149: def test_production_format():
        150:     """Test 7: Production-ready format"""
    >   151:     print("=== Test 7: Production Format ===")
        152:     prod_format = (
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 201 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/201.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:175:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        173:             logger.error("Cache miss rate high")
        174:     logging.flush()  # Ensure all log messages are processed
    >   175:     print()
        176:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 203 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/203.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:180:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        178: def test_minimal_format():
        179:     """Test 8: Minimal format for clean output"""
    >   180:     print("=== Test 8: Minimal Format ===")
        181:     logging.basicConfig(format="%(message)s")
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 204 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/204.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:189:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        187:     logger.error("Error: Something went wrong")
        188:     logging.flush()  # Ensure all log messages are processed
    >   189:     print()
        190:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 215 — CWE-117

- Function context: `scripts/logxide/findings/functions/215.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/performance_demo.py:145:9`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
        143:         logger = logging.getLogger(logger_name)
        144:         loggers.append(logger)
    >   145:         logger.info(f"Logger {logger_name} initialized")
        146:
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 218 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/218.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:35:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         33: except ImportError:
         34:     HAS_SENTRY_SDK = False
    >    35:     print("⚠️  sentry-sdk is not installed.")
         36:     print("   Install it with: pip install logxide[sentry]")
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 219 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/219.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:36:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         34:     HAS_SENTRY_SDK = False
         35:     print("⚠️  sentry-sdk is not installed.")
    >    36:     print("   Install it with: pip install logxide[sentry]")
         37:     print("   Exiting example.")
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 220 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/220.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:37:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         35:     print("⚠️  sentry-sdk is not installed.")
         36:     print("   Install it with: pip install logxide[sentry]")
    >    37:     print("   Exiting example.")
         38:     sys.exit(1)
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 221 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/221.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:44:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         42: 
         43: if not SENTRY_DSN:
    >    44:     print("⚠️  No SENTRY_DSN environment variable found.")
         45:     print(
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 222 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/222.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:45:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         43: if not SENTRY_DSN:
         44:     print("⚠️  No SENTRY_DSN environment variable found.")
    >    45:     print(
         46:         "   Set SENTRY_DSN to your actual Sentry project DSN to see events in Sentry."
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 223 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/223.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:48:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         46:         "   Set SENTRY_DSN to your actual Sentry project DSN to see events in Sentry."
         47:     )
    >    48:     print("   Running example without Sentry to demonstrate LogXide functionality.")
         49:     print()
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 224 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/224.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:49:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         47:     )
         48:     print("   Running example without Sentry to demonstrate LogXide functionality.")
    >    49:     print()
         50:     # Don't configure Sentry if no valid DSN
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 225 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/225.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:53:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         51:     SENTRY_CONFIGURED = False
         52: else:
    >    53:     print(f"✅ Using Sentry DSN: {SENTRY_DSN[:20]}...")
         54:     SENTRY_CONFIGURED = True
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 226 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/226.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:89:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         87: def demo_basic_logging():
         88:     """Demonstrate basic logging with automatic Sentry integration."""
    >    89:     print("=== Basic Logging Demo ===")
         90:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 227 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/227.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:100:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         98:     app_logger.critical("This is critical and will appear in Sentry")
         99: 
    >   100:     print("✓ Basic logging messages sent (check Sentry for WARNING/ERROR/CRITICAL)")
        101:     print()
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 228 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/228.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:101:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         99: 
        100:     print("✓ Basic logging messages sent (check Sentry for WARNING/ERROR/CRITICAL)")
    >   101:     print()
        102:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 229 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/229.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:106:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        104: def demo_exception_handling():
        105:     """Demonstrate exception logging with Sentry integration."""
    >   106:     print("=== Exception Handling Demo ===")
        107:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 230 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/230.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:111:1`
- Checklist pattern: catch-all delegates to full exception reporting

Source excerpt:

```
        109:         # This will cause a ZeroDivisionError
        110:         result = 10 / 0
    >   111:     except Exception as e:
        112:         # LogXide will automatically send this exception to Sentry
```

Why this is a false positive: the handler body reports the (deliberately raised) exception through the full exception-reporting path (`logger.exception` / `error(..., exc_info=True)`, forwarded to Sentry) — handled, not swallowed.

Checklist evidence: the handler body is a full exception-reporting call (`logger.exception`/`exc_info=True`).

### [ ] Finding 231 — CWE-396

- Function context: `scripts/logxide/findings/functions/231.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:111:1`
- Checklist pattern: catch-all delegates to full exception reporting

Source excerpt:

```
        109:         # This will cause a ZeroDivisionError
        110:         result = 10 / 0
    >   111:     except Exception as e:
        112:         # LogXide will automatically send this exception to Sentry
```

Why this is a false positive: same handler as finding 230: the demo deliberately raises the exception and fully reports it through the Sentry integration — the failure is not hidden.

Checklist evidence: the exception is routed to the error-reporting contract, not hidden.

### [ ] Finding 232 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/232.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:119:1`
- Checklist pattern: catch-all delegates to full exception reporting

Source excerpt:

```
        117:         data = {"name": "Alice"}
        118:         value = data["age"]  # Key doesn't exist
    >   119:     except Exception as e:
        120:         # Another exception with additional context
```

Why this is a false positive: the handler body reports the (deliberately raised) exception through the full exception-reporting path (`logger.exception` / `error(..., exc_info=True)`, forwarded to Sentry) — handled, not swallowed.

Checklist evidence: the handler body is a full exception-reporting call (`logger.exception`/`exc_info=True`).

### [ ] Finding 234 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/234.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:123:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        121:         app_logger.error(f"Failed to access user data: {str(e)}", exc_info=True)
        122: 
    >   123:     print("✓ Exceptions sent to Sentry with full stack traces")
        124:     print()
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 235 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/235.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:124:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        122: 
        123:     print("✓ Exceptions sent to Sentry with full stack traces")
    >   124:     print()
        125:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 236 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/236.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:129:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        127: def demo_structured_logging():
        128:     """Demonstrate structured logging with extra context."""
    >   129:     print("=== Structured Logging Demo ===")
        130:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 237 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/237.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:149:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        147:         db_logger.error("Database connection timeout during user query")
        148: 
    >   149:     print("✓ Structured logs with user and operation context sent to Sentry")
        150:     print()
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 238 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/238.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:150:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        148: 
        149:     print("✓ Structured logs with user and operation context sent to Sentry")
    >   150:     print()
        151:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 239 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/239.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:155:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        153: def demo_api_error_tracking():
        154:     """Demonstrate API error tracking simulation."""
    >   155:     print("=== API Error Tracking Demo ===")
        156:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 242 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/242.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:179:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        177:                 )
        178: 
    >   179:     print("✓ API errors with endpoint context sent to Sentry")
        180:     print()
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 243 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/243.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:180:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        178: 
        179:     print("✓ API errors with endpoint context sent to Sentry")
    >   180:     print()
        181:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 244 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/244.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:185:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        183: def demo_performance_monitoring():
        184:     """Demonstrate performance monitoring with logging."""
    >   185:     print("=== Performance Monitoring Demo ===")
        186:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 247 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/247.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:213:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        211:                 app_logger.warning(f"Operation '{op['name']}' is slow: {duration:.1f}s")
        212: 
    >   213:     print("✓ Performance issues logged to Sentry")
        214:     print()
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 248 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/248.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:214:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        212: 
        213:     print("✓ Performance issues logged to Sentry")
    >   214:     print()
        215:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 249 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/249.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:219:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        217: def demo_manual_sentry_handler():
        218:     """Demonstrate manual Sentry handler configuration."""
    >   219:     print("=== Manual Sentry Handler Demo ===")
        220:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 250 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/250.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:241:9`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        239:         custom_logger.error("This error will go to Sentry")
        240: 
    >   241:         print("✓ Manual Sentry handler configured and tested")
        242:     else:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 251 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/251.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:243:9`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        241:         print("✓ Manual Sentry handler configured and tested")
        242:     else:
    >   243:         print("❌ SentryHandler not available (sentry-sdk not installed)")
        244:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 252 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/252.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:245:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        243:         print("❌ SentryHandler not available (sentry-sdk not installed)")
        244: 
    >   245:     print()
        246:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 253 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/253.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:250:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        248: def demo_framework_integration():
        249:     """Demonstrate integration with web frameworks."""
    >   250:     print("=== Framework Integration Demo ===")
        251:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 254 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/254.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:283:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        281:         fastapi_logger.error("FastAPI dependency injection failed")
        282: 
    >   283:     print("✓ Framework-specific errors logged with context")
        284:     print()
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 255 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/255.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:284:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        282: 
        283:     print("✓ Framework-specific errors logged with context")
    >   284:     print()
        285:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 256 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/256.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:289:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        287: def demo_batch_processing():
        288:     """Demonstrate batch processing with error tracking."""
    >   289:     print("=== Batch Processing Demo ===")
        290:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 260 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/260.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:325:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        323:         )
        324: 
    >   325:     print("✓ Batch processing errors tracked with item-level context")
        326:     print()
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 261 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/261.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:326:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        324: 
        325:     print("✓ Batch processing errors tracked with item-level context")
    >   326:     print()
        327:
```

Why this is a false positive: the print is output of a standalone example/demo script under `examples/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 263 — CWE-117

- Function context: `scripts/logxide/findings/functions/263.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/simple_demo.py:39:9`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
         37:     start = time.time()
         38:     for i in range(10000):
    >    39:         logger.info(f"Performance test message {i}")
         40:
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 270 — CWE-89

- Function context: `scripts/logxide/findings/functions/270.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/third_party_integration.py:71:17`
- Checklist pattern: SQL string is a constant literal with no interpolation

Source excerpt:

```
         69:         with engine.connect() as conn:
         70:             # Create a table
    >    71:             conn.execute(text("CREATE TABLE users (id INTEGER, name TEXT)"))
         72:
```

Why this is a false positive: the SQL string is a static literal (`text("CREATE TABLE users (id INTEGER, name TEXT)")`) with no f-string, `%`-formatting, or `.format` interpolation; the “dynamic SQL reaches execute” condition is not met.

Checklist evidence: the executed SQL contains no placeholders or interpolated expressions.

### [ ] Finding 279 — CWE-1333

- Function context: `scripts/logxide/findings/functions/279.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/compat_handlers.py:65:18`
- Checklist pattern: nested repetition alternatives are disjoint and anchored by literals (linear-time)

Source excerpt:

```
         63:         r"^(.?[<>=^])?[+ -]?#?0?(\d+|{\w+})?[,_]?(\.(\d+|{\w+}))?[bcdefgnosx%]?$", re.I
         64:     )
    >    65:     field_spec = re.compile(r"^(\d+|\w+)(\.\w+|\[[^]]+\])*$")
         66:
```

Why this is a false positive: in `(\.\w+|\[[^]]+\])*` every repetition is anchored by a distinct literal delimiter (`.` or `[`) that the inner character classes cannot match, so the alternatives are disjoint and matching is linear — no catastrophic backtracking.

Checklist evidence: the repeated group's alternatives start with disjoint literal delimiters, eliminating ambiguity between iterations.

### [ ] Finding 284 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/284.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/compat_handlers.py:361:1`
- Checklist pattern: catch-all delegates to the logging error contract

Source excerpt:

```
        359:         except RecursionError:
        360:             raise
    >   361:         except Exception:
        362:             self.handleError(record)
```

Why this is a false positive: the `except Exception` is the mandated logging `emit()` contract (emit must never propagate) and the body routes the failure to `self.handleError(record)` — the exception is handled, not swallowed.

Checklist evidence: the handler body delegates to `handleError`, the stdlib emit error path.

### [ ] Finding 285 — CWE-396

- Function context: `scripts/logxide/findings/functions/285.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/compat_handlers.py:361:1`
- Checklist pattern: catch-all delegates to the logging error contract

Source excerpt:

```
        359:         except RecursionError:
        360:             raise
    >   361:         except Exception:
        362:             self.handleError(record)
```

Why this is a false positive: same handler as the paired BP-PY-1 finding: the generic catch is the required `emit()` contract and the error is handled via `self.handleError(record)`.

Checklist evidence: the exception is routed to `handleError`, not hidden.

### [ ] Finding 287 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/287.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/fast_logger_wrapper.py:55:1`
- Checklist pattern: catch-all recovers with a fallback value

Source excerpt:

```
         53:             self._effective_level = self._rust_logger.getEffectiveLevel()
         54:             self._name = self._rust_logger.name
    >    55:         except Exception:
         56:             # Fallback to safe defaults
```

Why this is a false positive: the handler body recovers with an explicit fallback value (safe defaults, `str(record.msg)`, `old_level = None`) — the failure is handled, satisfying the rule's “without handling” clause.

Checklist evidence: the handler body performs explicit fallback recovery.

### [ ] Finding 288 — CWE-396

- Function context: `scripts/logxide/findings/functions/288.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/fast_logger_wrapper.py:55:1`
- Checklist pattern: catch-all recovers with a fallback value

Source excerpt:

```
         53:             self._effective_level = self._rust_logger.getEffectiveLevel()
         54:             self._name = self._rust_logger.name
    >    55:         except Exception:
         56:             # Fallback to safe defaults
```

Why this is a false positive: same recovery handler as the paired BP-PY-1 finding: a deliberate fallback converts the exception into a safe default state.

Checklist evidence: the exception is converted into a fallback value/state.

### [ ] Finding 289 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/289.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/handlers.py:191:1`
- Checklist pattern: catch-all delegates to the logging error contract

Source excerpt:

```
        189:                     record.args = None
        190:                 self._inner.emit(_prepare_record_for_rust(record))
    >   191:         except Exception:
        192:             self.handleError(record)
```

Why this is a false positive: the `except Exception` is the mandated logging `emit()` contract (emit must never propagate) and the body routes the failure to `self.handleError(record)` — the exception is handled, not swallowed.

Checklist evidence: the handler body delegates to `handleError`, the stdlib emit error path.

### [ ] Finding 290 — CWE-396

- Function context: `scripts/logxide/findings/functions/290.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/handlers.py:191:1`
- Checklist pattern: catch-all delegates to the logging error contract

Source excerpt:

```
        189:                     record.args = None
        190:                 self._inner.emit(_prepare_record_for_rust(record))
    >   191:         except Exception:
        192:             self.handleError(record)
```

Why this is a false positive: same handler as the paired BP-PY-1 finding: the generic catch is the required `emit()` contract and the error is handled via `self.handleError(record)`.

Checklist evidence: the exception is routed to `handleError`, not hidden.

### [ ] Finding 291 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/291.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/handlers.py:260:1`
- Checklist pattern: catch-all delegates to the logging error contract

Source excerpt:

```
        258:                     record.args = None
        259:                 self._inner.emit(_prepare_record_for_rust(record))
    >   260:         except Exception:
        261:             self.handleError(record)
```

Why this is a false positive: the `except Exception` is the mandated logging `emit()` contract (emit must never propagate) and the body routes the failure to `self.handleError(record)` — the exception is handled, not swallowed.

Checklist evidence: the handler body delegates to `handleError`, the stdlib emit error path.

### [ ] Finding 292 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/292.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/handlers.py:325:1`
- Checklist pattern: catch-all delegates to the logging error contract

Source excerpt:

```
        323:                     record.args = None
        324:                 self._inner.emit(_prepare_record_for_rust(record))
    >   325:         except Exception:
        326:             self.handleError(record)
```

Why this is a false positive: the `except Exception` is the mandated logging `emit()` contract (emit must never propagate) and the body routes the failure to `self.handleError(record)` — the exception is handled, not swallowed.

Checklist evidence: the handler body delegates to `handleError`, the stdlib emit error path.

### [ ] Finding 293 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/293.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/handlers.py:406:1`
- Checklist pattern: catch-all delegates to the logging error contract

Source excerpt:

```
        404:             rust_record = _prepare_record_for_rust(record)
        405:             self._inner.emit(rust_record)
    >   406:         except Exception:
        407:             self.handleError(record)
```

Why this is a false positive: the `except Exception` is the mandated logging `emit()` contract (emit must never propagate) and the body routes the failure to `self.handleError(record)` — the exception is handled, not swallowed.

Checklist evidence: the handler body delegates to `handleError`, the stdlib emit error path.

### [ ] Finding 294 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/294.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/handlers.py:477:1`
- Checklist pattern: catch-all delegates to the logging error contract

Source excerpt:

```
        475:             rust_record = _prepare_record_for_rust(record)
        476:             self._inner.emit(rust_record)
    >   477:         except Exception:
        478:             self.handleError(record)
```

Why this is a false positive: the `except Exception` is the mandated logging `emit()` contract (emit must never propagate) and the body routes the failure to `self.handleError(record)` — the exception is handled, not swallowed.

Checklist evidence: the handler body delegates to `handleError`, the stdlib emit error path.

### [ ] Finding 295 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/295.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/handlers.py:519:1`
- Checklist pattern: catch-all delegates to the logging error contract

Source excerpt:

```
        517:             # MemoryHandler is always native: forward raw; caplog reads _inner.
        518:             self._inner.emit(_prepare_record_for_rust(record, native=True))
    >   519:         except Exception:
        520:             self.handleError(record)
```

Why this is a false positive: the `except Exception` is the mandated logging `emit()` contract (emit must never propagate) and the body routes the failure to `self.handleError(record)` — the exception is handled, not swallowed.

Checklist evidence: the handler body delegates to `handleError`, the stdlib emit error path.

### [ ] Finding 297 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/297.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/interceptor.py:42:1`
- Checklist pattern: catch-all recovers with a fallback value

Source excerpt:

```
         40:             try:
         41:                 message = record.getMessage()
    >    42:             except Exception:
         43:                 message = str(record.msg)
```

Why this is a false positive: the handler body recovers with an explicit fallback value (safe defaults, `str(record.msg)`, `old_level = None`) — the failure is handled, satisfying the rule's “without handling” clause.

Checklist evidence: the handler body performs explicit fallback recovery.

### [ ] Finding 298 — CWE-396

- Function context: `scripts/logxide/findings/functions/298.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/interceptor.py:42:1`
- Checklist pattern: catch-all recovers with a fallback value

Source excerpt:

```
         40:             try:
         41:                 message = record.getMessage()
    >    42:             except Exception:
         43:                 message = str(record.msg)
```

Why this is a false positive: same recovery handler as the paired BP-PY-1 finding: a deliberate fallback converts the exception into a safe default state.

Checklist evidence: the exception is converted into a fallback value/state.

### [ ] Finding 303 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/303.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/sentry_integration.py:85:1`
- Checklist pattern: catch-all delegates to the error handler contract

Source excerpt:

```
         83:                 self._add_breadcrumb(record, level_name, message, logger_name)
         84: 
    >    85:         except Exception as e:
         86:             # Prevent logging errors from causing infinite loops
```

Why this is a false positive: the handler body routes the exception to `self._handle_error(e)` (guarding against recursive logging errors) — handled, not swallowed.

Checklist evidence: the handler body delegates to the dedicated `_handle_error` path.

### [ ] Finding 304 — CWE-396

- Function context: `scripts/logxide/findings/functions/304.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/sentry_integration.py:85:1`
- Checklist pattern: catch-all delegates to the error handler contract

Source excerpt:

```
         83:                 self._add_breadcrumb(record, level_name, message, logger_name)
         84: 
    >    85:         except Exception as e:
         86:             # Prevent logging errors from causing infinite loops
```

Why this is a false positive: same handler as finding 303: deliberate catch-all in the Sentry emit path with explicit `_handle_error` handling.

Checklist evidence: the exception is routed to `_handle_error`, not hidden.

### [ ] Finding 307 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/307.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/testing.py:163:1`
- Checklist pattern: catch-all recovers with a fallback value

Source excerpt:

```
        161:         try:
        162:             old_level = handler._inner.level if hasattr(handler, "_inner") else None
    >   163:         except Exception:
        164:             old_level = None
```

Why this is a false positive: the handler body recovers with an explicit fallback value (safe defaults, `str(record.msg)`, `old_level = None`) — the failure is handled, satisfying the rule's “without handling” clause.

Checklist evidence: the handler body performs explicit fallback recovery.

### [ ] Finding 308 — CWE-396

- Function context: `scripts/logxide/findings/functions/308.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/testing.py:163:1`
- Checklist pattern: catch-all recovers with a fallback value

Source excerpt:

```
        161:         try:
        162:             old_level = handler._inner.level if hasattr(handler, "_inner") else None
    >   163:         except Exception:
        164:             old_level = None
```

Why this is a false positive: same recovery handler as the paired BP-PY-1 finding: a deliberate fallback converts the exception into a safe default state.

Checklist evidence: the exception is converted into a fallback value/state.

### [ ] Finding 309 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/309.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:18:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         16: def run_command(cmd, check=True):
         17:     """Run a command and return the result."""
    >    18:     print(f"Running: {cmd}")
         19:     result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 312 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/312.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:21:9`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         19:     result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
         20:     if check and result.returncode != 0:
    >    21:         print(f"Error: {result.stderr}")
         22:         sys.exit(1)
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 313 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/313.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:28:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         26: def check_version_consistency():
         27:     """Check that versions are consistent across files."""
    >    28:     print("Checking version consistency...")
         29:
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 314 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/314.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:35:13`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         33:         pyproject_match = re.search(r'version = "([^"]+)"', content)
         34:         if not pyproject_match:
    >    35:             print("Error: Could not find version in pyproject.toml")
         36:             sys.exit(1)
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 315 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/315.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:44:13`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         42:         init_match = re.search(r'__version__ = "([^"]+)"', content)
         43:         if not init_match:
    >    44:             print("Error: Could not find __version__ in logxide/__init__.py")
         45:             sys.exit(1)
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 316 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/316.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:50:9`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         48:     # Check consistency
         49:     if pyproject_version != init_version:
    >    50:         print("Error: Version mismatch!")
         51:         print(f"  pyproject.toml: {pyproject_version}")
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 317 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/317.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:51:9`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         49:     if pyproject_version != init_version:
         50:         print("Error: Version mismatch!")
    >    51:         print(f"  pyproject.toml: {pyproject_version}")
         52:         print(f"  __init__.py: {init_version}")
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 318 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/318.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:52:9`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         50:         print("Error: Version mismatch!")
         51:         print(f"  pyproject.toml: {pyproject_version}")
    >    52:         print(f"  __init__.py: {init_version}")
         53:         sys.exit(1)
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 319 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/319.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:55:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         53:         sys.exit(1)
         54: 
    >    55:     print(f"✓ Version consistency check passed: {pyproject_version}")
         56:     return pyproject_version
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 320 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/320.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:61:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         59: def check_git_status():
         60:     """Check git status to ensure clean working directory."""
    >    61:     print("Checking git status...")
         62:     result = run_command("git status --porcelain", check=False)
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 321 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/321.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:64:9`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         62:     result = run_command("git status --porcelain", check=False)
         63:     if result.stdout.strip():
    >    64:         print("Warning: Working directory is not clean:")
         65:         print(result.stdout)
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 322 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/322.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:65:9`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         63:     if result.stdout.strip():
         64:         print("Warning: Working directory is not clean:")
    >    65:         print(result.stdout)
         66:         response = input("Continue anyway? (y/N): ")
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 323 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/323.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:70:9`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         68:             sys.exit(1)
         69:     else:
    >    70:         print("✓ Working directory is clean")
         71:
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 324 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/324.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:75:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         73: def run_tests():
         74:     """Run the test suite."""
    >    75:     print("Running tests...")
         76:
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 325 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/325.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:79:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         77:     # Run Rust tests
         78:     run_command("cargo test")
    >    79:     print("✓ Rust tests passed")
         80:
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 326 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/326.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:83:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         81:     # Build Python extension
         82:     run_command("maturin develop")
    >    83:     print("✓ Python extension built")
         84:
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 327 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/327.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:87:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         85:     # Run Python tests
         86:     run_command("python -m pytest tests/ -v")
    >    87:     print("✓ Python tests passed")
         88:
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 328 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/328.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:92:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         90: def build_package():
         91:     """Build the package."""
    >    92:     print("Building package...")
         93:
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 329 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/329.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:100:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         98:     # Build wheels
         99:     run_command("maturin build --release")
    >   100:     print("✓ Package built")
        101:
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 330 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/330.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:105:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        103: def check_package():
        104:     """Check the built package."""
    >   105:     print("Checking package...")
        106:
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 331 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/331.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:112:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        110:     # Check package
        111:     run_command("twine check target/wheels/*")
    >   112:     print("✓ Package check passed")
        113:
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 332 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/332.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:117:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        115: def upload_to_test_pypi():
        116:     """Upload to Test PyPI."""
    >   117:     print("Uploading to Test PyPI...")
        118:
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 333 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/333.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:121:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        119:     # Upload to Test PyPI
        120:     run_command("twine upload --repository testpypi target/wheels/*")
    >   121:     print("✓ Uploaded to Test PyPI")
        122:
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 334 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/334.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:126:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        124: def upload_to_pypi():
        125:     """Upload to PyPI."""
    >   126:     print("Uploading to PyPI...")
        127:
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 335 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/335.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:130:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        128:     # Upload to PyPI
        129:     run_command("twine upload target/wheels/*")
    >   130:     print("✓ Uploaded to PyPI")
        131:
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 336 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/336.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:135:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        133: def create_git_tag(version):
        134:     """Create a git tag for the version."""
    >   135:     print(f"Creating git tag v{version}...")
        136:
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 337 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/337.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:142:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        140:     # Push tag
        141:     run_command(f"git push origin v{version}")
    >   142:     print(f"✓ Created and pushed tag v{version}")
        143:
```

Why this is a false positive: the print is output of a maintenance CLI script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 341 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/341.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:14:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         12: def test_import():
         13:     """Test basic import functionality."""
    >    14:     print("Testing basic import...")
         15:
```

Why this is a false positive: the print is output of a standalone verification script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 342 — BP-PY-42

- Function context: `scripts/logxide/findings/functions/342.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:16:1`
- Checklist pattern: try/except is a defensive guard, not a failure expectation

Source excerpt:

```
         14:     print("Testing basic import...")
         15: 
    >    16:     try:
         17:         import logxide
```

Why this is a false positive: the try/except is a defensive guard that reports failure to the console so the smoke script can continue; the test does not “expect failure”, so the rule condition (try/except instead of `pytest.raises`) is not met.

Checklist evidence: the except branch reports the failure rather than asserting a raise.

### [ ] Finding 343 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/343.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:19:9`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         17:         import logxide
         18: 
    >    19:         print(f"✓ Successfully imported logxide version {logxide.__version__}")
         20:         return True
```

Why this is a false positive: the print is output of a standalone verification script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 346 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/346.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:22:9`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         20:         return True
         21:     except Exception as e:
    >    22:         print(f"❌ Failed to import logxide: {e}")
         23:         traceback.print_exc()
```

Why this is a false positive: the print is output of a standalone verification script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 348 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/348.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:29:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         27: def test_logging_module():
         28:     """Test logging module functionality."""
    >    29:     print("Testing logging module...")
         30:
```

Why this is a false positive: the print is output of a standalone verification script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 349 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/349.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:42:9`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         40:         logger.error("Test error message")
         41: 
    >    42:         print("✓ Basic logging functionality works")
         43:         return True
```

Why this is a false positive: the print is output of a standalone verification script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 351 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/351.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:45:9`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         43:         return True
         44:     except Exception as e:
    >    45:         print(f"❌ Failed logging test: {e}")
         46:         traceback.print_exc()
```

Why this is a false positive: the print is output of a standalone verification script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 353 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/353.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:52:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         50: def test_drop_in_replacement():
         51:     """Test drop-in replacement functionality."""
    >    52:     print("Testing drop-in replacement...")
         53:
```

Why this is a false positive: the print is output of a standalone verification script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 354 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/354.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:68:13`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         66:             logger.info("Drop-in replacement test")
         67:             logging.flush()
    >    68:             print("✓ Drop-in replacement functionality works")
         69:             return True
```

Why this is a false positive: the print is output of a standalone verification script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 355 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/355.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:71:13`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         69:             return True
         70:         else:
    >    71:             print("❌ Drop-in replacement failed: missing flush method")
         72:             return False
```

Why this is a false positive: the print is output of a standalone verification script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 357 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/357.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:74:9`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         72:             return False
         73:     except Exception as e:
    >    74:         print(f"❌ Failed drop-in replacement test: {e}")
         75:         traceback.print_exc()
```

Why this is a false positive: the print is output of a standalone verification script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 359 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/359.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:81:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         79: def test_performance():
         80:     """Test basic performance (simple benchmark)."""
    >    81:     print("Testing performance...")
         82:
```

Why this is a false positive: the print is output of a standalone verification script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 361 — CWE-117

- Function context: `scripts/logxide/findings/functions/361.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:94:13`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
         92:         start_time = time.time()
         93:         for i in range(1000):
    >    94:             logger.info(f"Performance test message {i}")
         95:         logging.flush()
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 362 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/362.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:101:9`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         99:         messages_per_second = 1000 / duration
        100: 
    >   101:         print(
        102:             f"✓ Performance test completed: {messages_per_second:.0f} messages/second"
```

Why this is a false positive: the print is output of a standalone verification script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 364 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/364.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:106:9`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        104:         return True
        105:     except Exception as e:
    >   106:         print(f"❌ Failed performance test: {e}")
        107:         traceback.print_exc()
```

Why this is a false positive: the print is output of a standalone verification script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 366 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/366.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:113:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        111: def test_thread_safety():
        112:     """Test thread safety."""
    >   113:     print("Testing thread safety...")
        114:
```

Why this is a false positive: the print is output of a standalone verification script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 368 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/368.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:141:9`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        139: 
        140:         logging.flush()
    >   141:         print("✓ Thread safety test completed")
        142:         return True
```

Why this is a false positive: the print is output of a standalone verification script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 370 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/370.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:144:9`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        142:         return True
        143:     except Exception as e:
    >   144:         print(f"❌ Failed thread safety test: {e}")
        145:         traceback.print_exc()
```

Why this is a false positive: the print is output of a standalone verification script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 372 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/372.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:151:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        149: def test_formatting():
        150:     """Test formatting capabilities."""
    >   151:     print("Testing formatting...")
        152:
```

Why this is a false positive: the print is output of a standalone verification script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 373 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/373.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:170:9`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        168: 
        169:         logging.flush()
    >   170:         print("✓ Formatting test completed")
        171:         return True
```

Why this is a false positive: the print is output of a standalone verification script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 375 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/375.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:173:9`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
        171:         return True
        172:     except Exception as e:
    >   173:         print(f"❌ Failed formatting test: {e}")
        174:         traceback.print_exc()
```

Why this is a false positive: the print is output of a standalone verification script under `scripts/`, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 377 — BP-PY-42

- Function context: `scripts/logxide/findings/functions/377.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/conftest.py:28:1`
- Checklist pattern: flagged try/except is a helper, not a test expecting failure

Source excerpt:

```
         26: 
         27:     def _do_flush():
    >    28:         try:
         29:             logging.flush()
```

Why this is a false positive: `_do_flush` is a helper inside a timeout guard in `conftest.py`, not a test using try/except to expect a raise; the rule condition (“test uses try/except to expect failure”) is not met.

Checklist evidence: the construct is a defensive flush helper in a fixture, not a failure-expecting test.

### [ ] Finding 382 — BP-PY-40

- Function context: `scripts/logxide/findings/functions/382.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/conftest.py:36:6`
- Checklist pattern: daemon thread with explicit synchronization, no join needed

Source excerpt:

```
         34: 
         35:     t = threading.Thread(target=_do_flush, daemon=True)
    >    36:     t.start()
         37:     done.wait(timeout=timeout_seconds)
```

Why this is a false positive: the thread is `daemon=True` (documented as a deliberate deadlock-bailout) and the caller synchronizes with `done.wait(timeout)`; the rule's own fix targets “fire-and-forget non-daemon threads”, so the condition is not met.

Checklist evidence: daemon thread plus explicit `Event.wait` synchronization protocol.

### [ ] Finding 397 — BP-PY-42

- Function context: `scripts/logxide/findings/functions/397.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_basic_logging.py:97:1`
- Checklist pattern: try/except is a crash-guard, not a failure expectation

Source excerpt:

```
         95: 
         96:         # Test various format configurations - just make sure they don't crash
    >    97:         try:
         98:             logging.basicConfig(format="%(levelname)s: %(message)s")
```

Why this is a false positive: the try/except (“just make sure they don't crash”) guards a smoke test against unsupported formats rather than expecting a raise; `pytest.raises` is the wrong construct, so the condition is not met.

Checklist evidence: the except branch tolerates unsupported formats instead of asserting an exception.

### [ ] Finding 405 — CWE-117

- Function context: `scripts/logxide/findings/functions/405.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_basic_logging.py:147:17`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
        145: 
        146:             for i in range(5):
    >   147:                 logger.info(f"Message {i}")
        148:                 time.sleep(0.001)  # Small delay
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 411 — BP-PY-42

- Function context: `scripts/logxide/findings/functions/411.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_compatibility.py:10:1`
- Checklist pattern: try/except deliberately creates the exception to exercise `logging.exception`

Source excerpt:

```
          8: 
          9: def test_exception_logging():
    >    10:     try:
         11:         raise ValueError("test")
```

Why this is a false positive: the test deliberately raises its own controlled exception to exercise the `logging.exception` path; it is not expecting failure of a call under test, so `pytest.raises` is not the applicable construct.

Checklist evidence: the try body raises a controlled exception to drive the logging exception handler.

### [ ] Finding 412 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/412.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_compatibility.py:12:1`
- Checklist pattern: catch-all delegates to full exception reporting

Source excerpt:

```
         10:     try:
         11:         raise ValueError("test")
    >    12:     except:
         13:         logging.exception("test_exception_logging")
```

Why this is a false positive: the handler body reports the (deliberately raised) exception through the full exception-reporting path (`logger.exception` / `error(..., exc_info=True)`, forwarded to Sentry) — handled, not swallowed.

Checklist evidence: the handler body is a full exception-reporting call (`logger.exception`/`exc_info=True`).

### [ ] Finding 414 — BP-PY-41

- Function context: `scripts/logxide/findings/functions/414.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_doc_codeblocks.py:65:1`
- Checklist pattern: test verifies outcomes via `pytest.fail`, not side-effects-only

Source excerpt:

```
         63:     ids=[b[0] for b in ALL_CODEBLOCKS],
         64: )
    >    65: def test_doc_codeblock_syntax(block_id: str, code: str):
         66:     """Every Python code block in docs must be syntactically valid."""
```

Why this is a false positive: the test verifies outcomes via `pytest.fail` (and `pytest.skip`) instead of only performing side effects; it is not a placeholder test, so the heuristic's condition is not met.

Checklist evidence: the test contains explicit verification mechanisms (`pytest.fail`/`pytest.skip`).

### [ ] Finding 416 — CWE-94

- Function context: `scripts/logxide/findings/functions/416.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_doc_codeblocks.py:68:9`
- Checklist pattern: compile input is the repo's own committed docs, not externally influenced

Source excerpt:

```
         66:     """Every Python code block in docs must be syntactically valid."""
         67:     try:
    >    68:         compile(code, block_id, "exec")
         69:     except SyntaxError as exc:
```

Why this is a false positive: the compiled text comes from the repository's own documentation files and the resulting code object is discarded after syntax validation (never executed); it is not externally influenced input, so the CWE-94 condition is not met.

Checklist evidence: input source is repo-committed docs; the compile result is never executed.

### [ ] Finding 417 — BP-PY-12

- Function context: `scripts/logxide/findings/functions/417.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_doc_codeblocks.py:74:33`
- Checklist pattern: flagged line is a comment, not an eval/exec sink

Source excerpt:

```
         72: 
         73: # ---------------------------------------------------------------------------
    >    74: # 2. Execution verification — exec() standalone blocks
         75: # ---------------------------------------------------------------------------
```

Why this is a false positive: the flagged line 74 is a comment describing the section (“Execution verification — exec() standalone blocks”); there is no `exec`/`eval` call at that location — blocks are written to files and run in subprocesses — so the rule condition is not met.

Checklist evidence: the flagged token is prose in a comment; no dynamic eval/exec sink exists at the location.

### [ ] Finding 418 — BP-PY-41

- Function context: `scripts/logxide/findings/functions/418.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_doc_codeblocks.py:201:1`
- Checklist pattern: test verifies outcomes via `pytest.fail`, not side-effects-only

Source excerpt:

```
        199:     ids=[b[0] for b in EXEC_CODEBLOCKS],
        200: )
    >   201: def test_doc_codeblock_exec(block_id: str, code: str):
        202:     """Executable doc code blocks should run without errors.
```

Why this is a false positive: the test verifies outcomes via `pytest.fail` (and `pytest.skip`) instead of only performing side effects; it is not a placeholder test, so the heuristic's condition is not met.

Checklist evidence: the test contains explicit verification mechanisms (`pytest.fail`/`pytest.skip`).

### [ ] Finding 419 — CWE-829

- Function context: `scripts/logxide/findings/functions/419.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_doc_codeblocks.py:215:13`
- Checklist pattern: dynamic import selects from repo-derived framework names, not request input

Source excerpt:

```
        213:     for framework in needed:
        214:         try:
    >   215:             __import__(framework)
        216:         except ImportError:
```

Why this is a false positive: `needed` is computed by scanning the repository's own committed doc code blocks (`_needs_framework(code)`); the imported names come from package-controlled content, not from request-derived input, so the untrusted-control-sphere condition is not met.

Checklist evidence: module names originate from repo-committed code content, not user/request input.

### [ ] Finding 420 — CWE-829

- Function context: `scripts/logxide/findings/functions/420.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_examples.py:25:16`
- Checklist pattern: file path comes from the repo's own examples directory, not request input

Source excerpt:

```
         23:         # Load the module
         24:         module_name = example_path.stem
    >    25:         spec = importlib.util.spec_from_file_location(module_name, example_path)
         26:         if spec is None:
```

Why this is a false positive: `example_path` is iterated from the repository's own `examples/` directory; the rule's own fix (“never pass request-derived module names or paths”) does not apply to package-controlled repo files.

Checklist evidence: the path is derived from a local directory listing of repo files, not external input.

### [ ] Finding 423 — BP-PY-42

- Function context: `scripts/logxide/findings/functions/423.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_examples.py:55:1`
- Checklist pattern: try/except is a defensive guard, not a failure expectation

Source excerpt:

```
         53:     print("Testing LogXide basic functionality...")
         54: 
    >    55:     try:
         56:         # Use auto-install pattern
```

Why this is a false positive: the try/except is a defensive guard that reports failure to the console so the smoke script can continue; the test does not “expect failure”, so the rule condition (try/except instead of `pytest.raises`) is not met.

Checklist evidence: the except branch reports the failure rather than asserting a raise.

### [ ] Finding 426 — CWE-489

- Function context: `scripts/logxide/findings/functions/426.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_examples.py:137:1`
- Checklist pattern: debug flag set inside a test, not application source

Source excerpt:

```
        135:         if not settings.configured:
        136:             settings.configure(
    >   137:                 DEBUG=True,
        138:                 SECRET_KEY="test-key-for-integration-test",
```

Why this is a false positive: the flagged debug flag configures Django/Sentry inside a pytest test under `tests/`, where debug mode is standard test scaffolding; the rule condition (“debug mode is enabled in application source”) is not met.

Checklist evidence: the construct is test code configuring a fake settings/SDK object, not application source.

### [ ] Finding 427 — CWE-756

- Function context: `scripts/logxide/findings/functions/427.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_examples.py:137:1`
- Checklist pattern: debug flag set inside a test, not a deployed application

Source excerpt:

```
        135:         if not settings.configured:
        136:             settings.configure(
    >   137:                 DEBUG=True,
        138:                 SECRET_KEY="test-key-for-integration-test",
```

Why this is a false positive: the same test scaffolding as the paired CWE-489 finding; the “application explicitly enables debug error output” condition targets deployed application source, not test code.

Checklist evidence: the construct is test code, not application source.

### [ ] Finding 431 — CWE-117

- Function context: `scripts/logxide/findings/functions/431.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_examples.py:237:13`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
        235: 
        236:         for i in range(1000):
    >   237:             logger.info(f"Performance test message {i}")
        238:
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 439 — CWE-117

- Function context: `scripts/logxide/findings/functions/439.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_gil_scaling.py:123:9`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
        121:     n = 25
        122:     for i in range(n):
    >   123:         logger.info(f"secret-{i}")
        124:
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 441 — BP-PY-40

- Function context: `scripts/logxide/findings/functions/441.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_http_auth.py:41:11`
- Checklist pattern: daemon server thread with explicit shutdown protocol

Source excerpt:

```
         39:     port = server.server_address[1]
         40:     thread = threading.Thread(target=server.serve_forever, daemon=True)
    >    41:     thread.start()
         42:     yield port
```

Why this is a false positive: the fixture thread is `daemon=True` and teardown calls `server.shutdown()` — a clear shutdown protocol, satisfying the rule's own fix (avoid non-daemon fire-and-forget).

Checklist evidence: daemon thread plus explicit `shutdown()` on teardown.

### [ ] Finding 442 — BP-PY-40

- Function context: `scripts/logxide/findings/functions/442.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_http_handler_features.py:43:11`
- Checklist pattern: daemon server thread with explicit shutdown protocol

Source excerpt:

```
         41:     port = server.server_address[1]
         42:     thread = threading.Thread(target=server.serve_forever, daemon=True)
    >    43:     thread.start()
         44:     yield port
```

Why this is a false positive: the fixture thread is `daemon=True` and teardown calls `server.shutdown()` — a clear shutdown protocol, satisfying the rule's own fix (avoid non-daemon fire-and-forget).

Checklist evidence: daemon thread plus explicit `shutdown()` on teardown.

### [ ] Finding 446 — CWE-117

- Function context: `scripts/logxide/findings/functions/446.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration.py:105:17`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
        103:             logger = logging.getLogger("app.database")
        104:             for i in range(10):
    >   105:                 logger.info(f"Database query {i}")
        106:                 time.sleep(0.001)
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 457 — BP-PY-42

- Function context: `scripts/logxide/findings/functions/457.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:20:1`
- Checklist pattern: try/except is a defensive guard, not a failure expectation

Source excerpt:

```
         18:     print("Testing Flask integration...")
         19: 
    >    20:     try:
         21:         # Import and test basic functionality
```

Why this is a false positive: the try/except is a defensive guard that reports failure to the console so the smoke script can continue; the test does not “expect failure”, so the rule condition (try/except instead of `pytest.raises`) is not met.

Checklist evidence: the except branch reports the failure rather than asserting a raise.

### [ ] Finding 460 — CWE-489

- Function context: `scripts/logxide/findings/functions/460.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:69:1`
- Checklist pattern: debug flag set inside a test, not application source

Source excerpt:

```
         67:         if not settings.configured:
         68:             settings.configure(
    >    69:                 DEBUG=True,
         70:                 SECRET_KEY="test-key",
```

Why this is a false positive: the flagged debug flag configures Django/Sentry inside a pytest test under `tests/`, where debug mode is standard test scaffolding; the rule condition (“debug mode is enabled in application source”) is not met.

Checklist evidence: the construct is test code configuring a fake settings/SDK object, not application source.

### [ ] Finding 461 — CWE-756

- Function context: `scripts/logxide/findings/functions/461.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:69:1`
- Checklist pattern: debug flag set inside a test, not a deployed application

Source excerpt:

```
         67:         if not settings.configured:
         68:             settings.configure(
    >    69:                 DEBUG=True,
         70:                 SECRET_KEY="test-key",
```

Why this is a false positive: the same test scaffolding as the paired CWE-489 finding; the “application explicitly enables debug error output” condition targets deployed application source, not test code.

Checklist evidence: the construct is test code, not application source.

### [ ] Finding 467 — CWE-117

- Function context: `scripts/logxide/findings/functions/467.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:159:13`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
        157: 
        158:         for i in range(1000):
    >   159:             logger.info(f"Performance test message {i}")
        160:
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 490 — BP-PY-7

- Function context: `scripts/logxide/findings/functions/490.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_python315_compat.py:187:51`
- Checklist pattern: `open(` token appears in docstring prose, not a call

Source excerpt:

```
        185:         """Test that default encoding works regardless of Python version.
        186: 
    >   187:         In Python 3.15+, the default encoding for open() is UTF-8 (PEP 686).
        188:         In Python 3.12-3.14, it depends on locale.
```

Why this is a false positive: the flagged `open()` token is prose inside a docstring; the actual file opens in the function use `with open(...)` context managers.

Checklist evidence: the flagged token is documentation text, not an `open` call.

### [ ] Finding 491 — CWE-489

- Function context: `scripts/logxide/findings/functions/491.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_sentry_integration.py:353:1`
- Checklist pattern: debug flag set inside a test, not application source

Source excerpt:

```
        351:         sentry_sdk.init(
        352:             dsn="https://test@example.com/1",  # Valid format but fake
    >   353:             debug=True,
        354:             # Disable default integrations to avoid noise
```

Why this is a false positive: the flagged debug flag configures Django/Sentry inside a pytest test under `tests/`, where debug mode is standard test scaffolding; the rule condition (“debug mode is enabled in application source”) is not met.

Checklist evidence: the construct is test code configuring a fake settings/SDK object, not application source.

### [ ] Finding 492 — CWE-756

- Function context: `scripts/logxide/findings/functions/492.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_sentry_integration.py:353:1`
- Checklist pattern: debug flag set inside a test, not a deployed application

Source excerpt:

```
        351:         sentry_sdk.init(
        352:             dsn="https://test@example.com/1",  # Valid format but fake
    >   353:             debug=True,
        354:             # Disable default integrations to avoid noise
```

Why this is a false positive: the same test scaffolding as the paired CWE-489 finding; the “application explicitly enables debug error output” condition targets deployed application source, not test code.

Checklist evidence: the construct is test code, not application source.

### [ ] Finding 498 — CWE-117

- Function context: `scripts/logxide/findings/functions/498.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_timed_rotation.py:70:13`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
         68:         # Create multiple rotations
         69:         for i in range(4):
    >    70:             logger.info(f"message {i}")
         71:             handler.flush()
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 500 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/500.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tmp_test_intercept.py:10:1`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
          8: lock = threading.RLock()
          9: 
    >    10: print("Calling outside lock...")
         11: lx_logger.log(logging.INFO, "Outside")
```

Why this is a false positive: the print is output of a throwaway scratch script at the repo root, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 501 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/501.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tmp_test_intercept.py:13:1`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         11: lx_logger.log(logging.INFO, "Outside")
         12: 
    >    13: print("Acquiring lock...")
         14: with lock:
```

Why this is a false positive: the print is output of a throwaway scratch script at the repo root, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 502 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/502.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tmp_test_intercept.py:15:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         13: print("Acquiring lock...")
         14: with lock:
    >    15:     print("Calling inside lock...")
         16:     lx_logger.log(logging.INFO, "Inside lock!")
```

Why this is a false positive: the print is output of a throwaway scratch script at the repo root, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 503 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/503.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tmp_test_intercept.py:18:1`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         16:     lx_logger.log(logging.INFO, "Inside lock!")
         17: 
    >    18: print("FINISHED!")
         19:
```

Why this is a false positive: the print is output of a throwaway scratch script at the repo root, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 286 — CWE-1121

- Function context: `scripts/logxide/findings/functions/286.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/config.py:21:24`
- Checklist pattern: branch count below threshold of 12

Source excerpt:

```
def dictConfig(config):
    cfg = copy.deepcopy(config)
    if "handlers" in cfg and isinstance(cfg["handlers"], dict):
        for _name, handler_config in cfg["handlers"].items():
            if not isinstance(handler_config, dict):
                continue
            class_name = handler_config.get("class")
            if class_name == "logxide.FileHandler":
                ...
            elif class_name == "logxide.StreamHandler":
                ...
            elif class_name == "logxide.RotatingFileHandler":
                ...
            elif class_name == "logxide.HTTPHandler":
                ...
            elif class_name == "logxide.OTLPHandler":
                ...
            elif class_name in HANDLER_MAP:
                ...
    logging.config.dictConfig(cfg)
```

Why this is a false positive: CWE-1121 requires `branches >= 12` (`minimumRouteBranches` in `internal/lang/python/detectors/cwe/common.go`). Independent count of `if `/`elif `/`for `/`while `/`except ` tokens in `dictConfig` is **9** (3 `if` + 5 `elif` + 1 `for`) after docstring masking — below the threshold. Reclassified from Uncertain.

Checklist evidence: rule condition "at least twelve visible control-flow branches"; source has 9.

## True positives

The following findings satisfy the rule condition as shown in the chunk context and are reported compactly per rule.

### BP-PY-47 — 153 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 13 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/compare_loggers.py:138:26` | f-string message passed eagerly to a logging call |
| 15 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/compare_loggers.py:141:22` | f-string message passed eagerly to a logging call |
| 24 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/perf_vs_stdlib.py:73:26` | f-string message passed eagerly to a logging call |
| 32 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/basic_usage.py:64:5` | f-string message passed eagerly to a logging call |
| 34 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/basic_usage.py:65:5` | f-string message passed eagerly to a logging call |
| 35 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/basic_usage.py:66:5` | f-string message passed eagerly to a logging call |
| 40 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:190:17` | f-string message passed eagerly to a logging call |
| 41 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:200:17` | f-string message passed eagerly to a logging call |
| 42 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:214:16` | f-string message passed eagerly to a logging call |
| 45 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:216:16` | f-string message passed eagerly to a logging call |
| 46 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:222:17` | f-string message passed eagerly to a logging call |
| 47 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:250:25` | f-string message passed eagerly to a logging call |
| 48 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:254:25` | f-string message passed eagerly to a logging call |
| 50 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:290:13` | f-string message passed eagerly to a logging call |
| 51 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:314:17` | f-string message passed eagerly to a logging call |
| 53 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:318:17` | f-string message passed eagerly to a logging call |
| 54 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:325:9` | f-string message passed eagerly to a logging call |
| 55 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:329:13` | f-string message passed eagerly to a logging call |
| 56 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:332:13` | f-string message passed eagerly to a logging call |
| 58 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:348:12` | f-string message passed eagerly to a logging call |
| 59 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:400:13` | f-string message passed eagerly to a logging call |
| 61 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:406:13` | f-string message passed eagerly to a logging call |
| 62 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:413:9` | f-string message passed eagerly to a logging call |
| 63 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:419:21` | f-string message passed eagerly to a logging call |
| 64 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:428:9` | f-string message passed eagerly to a logging call |
| 65 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:447:9` | f-string message passed eagerly to a logging call |
| 66 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:516:21` | f-string message passed eagerly to a logging call |
| 68 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:529:17` | f-string message passed eagerly to a logging call |
| 69 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:562:21` | f-string message passed eagerly to a logging call |
| 70 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:575:13` | f-string message passed eagerly to a logging call |
| 71 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/drop_in_replacement.py:65:10` | f-string message passed eagerly to a logging call |
| 72 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/drop_in_replacement.py:66:10` | f-string message passed eagerly to a logging call |
| 73 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/drop_in_replacement.py:67:10` | f-string message passed eagerly to a logging call |
| 74 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/drop_in_replacement.py:68:10` | f-string message passed eagerly to a logging call |
| 75 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:139:19` | f-string message passed eagerly to a logging call |
| 76 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:145:19` | f-string message passed eagerly to a logging call |
| 77 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:153:19` | f-string message passed eagerly to a logging call |
| 80 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:160:27` | f-string message passed eagerly to a logging call |
| 81 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:206:10` | f-string message passed eagerly to a logging call |
| 82 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:213:18` | f-string message passed eagerly to a logging call |
| 83 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:237:22` | f-string message passed eagerly to a logging call |
| 85 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:247:14` | f-string message passed eagerly to a logging call |
| 86 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:278:29` | f-string message passed eagerly to a logging call |
| 88 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:286:25` | f-string message passed eagerly to a logging call |
| 91 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:342:13` | f-string message passed eagerly to a logging call |
| 92 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:355:13` | f-string message passed eagerly to a logging call |
| 94 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:374:12` | f-string message passed eagerly to a logging call |
| 95 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:383:9` | f-string message passed eagerly to a logging call |
| 96 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:390:9` | f-string message passed eagerly to a logging call |
| 97 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:426:9` | f-string message passed eagerly to a logging call |
| 98 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:434:9` | f-string message passed eagerly to a logging call |
| 99 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:439:13` | f-string message passed eagerly to a logging call |
| 100 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:444:13` | f-string message passed eagerly to a logging call |
| 101 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:449:13` | f-string message passed eagerly to a logging call |
| 102 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:456:9` | f-string message passed eagerly to a logging call |
| 103 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:460:13` | f-string message passed eagerly to a logging call |
| 104 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:463:9` | f-string message passed eagerly to a logging call |
| 105 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:473:9` | f-string message passed eagerly to a logging call |
| 106 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:485:9` | f-string message passed eagerly to a logging call |
| 107 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:497:13` | f-string message passed eagerly to a logging call |
| 108 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:501:13` | f-string message passed eagerly to a logging call |
| 109 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:508:9` | f-string message passed eagerly to a logging call |
| 110 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:512:13` | f-string message passed eagerly to a logging call |
| 111 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:515:9` | f-string message passed eagerly to a logging call |
| 112 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:528:23` | f-string message passed eagerly to a logging call |
| 114 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:536:19` | f-string message passed eagerly to a logging call |
| 116 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:550:12` | f-string message passed eagerly to a logging call |
| 117 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:606:13` | f-string message passed eagerly to a logging call |
| 118 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:610:13` | f-string message passed eagerly to a logging call |
| 119 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:617:9` | f-string message passed eagerly to a logging call |
| 120 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:623:21` | f-string message passed eagerly to a logging call |
| 121 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:632:9` | f-string message passed eagerly to a logging call |
| 122 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:649:9` | f-string message passed eagerly to a logging call |
| 123 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:655:21` | f-string message passed eagerly to a logging call |
| 125 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:666:9` | f-string message passed eagerly to a logging call |
| 127 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:695:13` | f-string message passed eagerly to a logging call |
| 128 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_demo.py:69:9` | f-string message passed eagerly to a logging call |
| 130 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_demo.py:77:5` | f-string message passed eagerly to a logging call |
| 131 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_demo.py:90:9` | f-string message passed eagerly to a logging call |
| 132 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:81:17` | f-string message passed eagerly to a logging call |
| 134 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:85:19` | f-string message passed eagerly to a logging call |
| 136 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:107:12` | f-string message passed eagerly to a logging call |
| 137 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:120:12` | f-string message passed eagerly to a logging call |
| 138 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:136:12` | f-string message passed eagerly to a logging call |
| 140 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:138:12` | f-string message passed eagerly to a logging call |
| 141 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:147:11` | f-string message passed eagerly to a logging call |
| 142 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:154:11` | f-string message passed eagerly to a logging call |
| 143 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:164:11` | f-string message passed eagerly to a logging call |
| 144 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:197:13` | f-string message passed eagerly to a logging call |
| 145 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:216:17` | f-string message passed eagerly to a logging call |
| 147 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:221:19` | f-string message passed eagerly to a logging call |
| 148 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:229:9` | f-string message passed eagerly to a logging call |
| 149 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:233:13` | f-string message passed eagerly to a logging call |
| 150 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:236:9` | f-string message passed eagerly to a logging call |
| 152 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:250:12` | f-string message passed eagerly to a logging call |
| 153 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:299:13` | f-string message passed eagerly to a logging call |
| 155 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:305:15` | f-string message passed eagerly to a logging call |
| 156 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:313:9` | f-string message passed eagerly to a logging call |
| 157 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:319:13` | f-string message passed eagerly to a logging call |
| 158 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:328:9` | f-string message passed eagerly to a logging call |
| 159 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:348:9` | f-string message passed eagerly to a logging call |
| 160 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:380:17` | f-string message passed eagerly to a logging call |
| 162 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:399:13` | f-string message passed eagerly to a logging call |
| 163 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:416:9` | f-string message passed eagerly to a logging call |
| 172 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/format_production.py:19:5` | f-string message passed eagerly to a logging call |
| 176 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/format_threaded.py:21:9` | f-string message passed eagerly to a logging call |
| 195 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:120:13` | f-string message passed eagerly to a logging call |
| 200 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:169:9` | f-string message passed eagerly to a logging call |
| 205 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/performance_demo.py:47:16` | f-string message passed eagerly to a logging call |
| 206 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/performance_demo.py:53:10` | f-string message passed eagerly to a logging call |
| 207 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/performance_demo.py:70:20` | f-string message passed eagerly to a logging call |
| 208 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/performance_demo.py:92:10` | f-string message passed eagerly to a logging call |
| 209 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/performance_demo.py:109:15` | f-string message passed eagerly to a logging call |
| 210 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/performance_demo.py:112:10` | f-string message passed eagerly to a logging call |
| 211 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/performance_demo.py:123:16` | f-string message passed eagerly to a logging call |
| 212 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/performance_demo.py:126:16` | f-string message passed eagerly to a logging call |
| 213 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/performance_demo.py:133:10` | f-string message passed eagerly to a logging call |
| 214 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/performance_demo.py:145:9` | f-string message passed eagerly to a logging call |
| 216 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/performance_demo.py:161:19` | f-string message passed eagerly to a logging call |
| 217 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/performance_demo.py:164:14` | f-string message passed eagerly to a logging call |
| 233 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:121:13` | f-string message passed eagerly to a logging call |
| 240 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:171:21` | f-string message passed eagerly to a logging call |
| 241 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:175:21` | f-string message passed eagerly to a logging call |
| 245 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:207:21` | f-string message passed eagerly to a logging call |
| 246 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:211:21` | f-string message passed eagerly to a logging call |
| 257 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:311:27` | f-string message passed eagerly to a logging call |
| 258 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:316:27` | f-string message passed eagerly to a logging call |
| 259 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/sentry_integration.py:321:15` | f-string message passed eagerly to a logging call |
| 262 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/simple_demo.py:39:9` | f-string message passed eagerly to a logging call |
| 265 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/third_party_integration.py:44:13` | f-string message passed eagerly to a logging call |
| 266 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/third_party_integration.py:48:13` | f-string message passed eagerly to a logging call |
| 269 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/third_party_integration.py:51:13` | f-string message passed eagerly to a logging call |
| 271 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/third_party_integration.py:80:17` | f-string message passed eagerly to a logging call |
| 272 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/third_party_integration.py:85:13` | f-string message passed eagerly to a logging call |
| 274 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/third_party_integration.py:88:13` | f-string message passed eagerly to a logging call |
| 275 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/third_party_integration.py:116:14` | f-string message passed eagerly to a logging call |
| 360 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:94:13` | f-string message passed eagerly to a logging call |
| 367 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:126:17` | f-string message passed eagerly to a logging call |
| 404 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_basic_logging.py:147:17` | f-string message passed eagerly to a logging call |
| 430 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_examples.py:237:13` | f-string message passed eagerly to a logging call |
| 438 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_gil_scaling.py:123:9` | f-string message passed eagerly to a logging call |
| 445 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration.py:105:17` | f-string message passed eagerly to a logging call |
| 447 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration.py:112:17` | f-string message passed eagerly to a logging call |
| 448 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration.py:119:17` | f-string message passed eagerly to a logging call |
| 449 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration.py:150:13` | f-string message passed eagerly to a logging call |
| 450 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration.py:152:13` | f-string message passed eagerly to a logging call |
| 452 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration.py:188:13` | f-string message passed eagerly to a logging call |
| 454 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration.py:228:13` | f-string message passed eagerly to a logging call |
| 455 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration.py:248:13` | f-string message passed eagerly to a logging call |
| 466 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:159:13` | f-string message passed eagerly to a logging call |
| 474 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:210:17` | f-string message passed eagerly to a logging call |
| 479 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:268:13` | f-string message passed eagerly to a logging call |
| 497 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_timed_rotation.py:70:13` | f-string message passed eagerly to a logging call |

### BP-PY-41 — 50 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 67 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:520:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 161 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:392:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 178 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:13:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 181 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:38:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 184 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:52:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 187 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:70:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 190 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:85:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 193 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:104:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 198 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:149:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 202 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:178:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 340 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:12:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 347 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:27:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 352 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:50:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 358 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:79:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 365 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:111:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 371 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:149:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 392 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_basic_logging.py:36:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 393 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_basic_logging.py:42:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 394 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_basic_logging.py:52:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 395 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_basic_logging.py:67:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 396 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_basic_logging.py:92:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 402 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_basic_logging.py:109:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 403 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_basic_logging.py:136:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 406 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_color_formatter.py:200:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 407 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_color_formatter.py:232:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 408 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_compat_functions.py:284:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 409 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_compatibility.py:5:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 410 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_compatibility.py:9:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 413 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_compatibility.py:16:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 422 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_examples.py:51:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 443 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration.py:45:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 444 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration.py:66:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 451 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration.py:176:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 453 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration.py:199:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 456 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:16:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 459 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:54:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 463 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:99:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 465 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:137:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 473 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:190:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 478 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:248:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 481 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:286:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 484 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_simple.py:13:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 485 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_simple.py:31:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 488 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_python315_compat.py:79:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 489 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_python315_compat.py:93:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 493 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_sentry_integration.py:743:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 494 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_simple.py:29:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 495 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_simple.py:34:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 496 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_simple.py:88:1` | test function performs side effects only, without assertions or pytest.fail/raises |
| 499 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_timed_rotation.py:137:1` | test function performs side effects only, without assertions or pytest.fail/raises |

### BP-PY-1 — 42 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 43 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:215:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 78 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:159:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 84 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:246:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 87 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:285:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 93 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:373:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 113 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:535:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 126 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:694:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 139 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:137:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 267 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/third_party_integration.py:50:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 273 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/third_party_integration.py:87:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 338 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:212:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 344 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:21:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 350 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:44:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 356 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:73:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 363 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:105:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 369 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:143:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 374 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:172:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 376 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:202:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 378 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/conftest.py:30:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 384 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/conftest.py:73:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 386 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/conftest.py:83:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 388 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/conftest.py:120:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 390 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/conftest.py:130:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 398 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_basic_logging.py:105:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 421 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_examples.py:44:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 424 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_examples.py:77:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 425 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_examples.py:118:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 428 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_examples.py:168:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 429 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_examples.py:210:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 432 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_examples.py:259:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 436 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_examples.py:266:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 437 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_examples.py:304:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 458 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:49:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 462 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:94:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 464 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:132:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 468 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:178:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 472 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:185:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 475 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:231:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 477 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:243:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 480 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:281:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 482 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:322:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |
| 483 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:350:1` | broad `except Exception`/bare `except:` whose body only logs/reports and continues (or is empty) |

### BP-PY-2 — 15 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 276 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/__init__.py:156:1` | exception handler body is only `pass` |
| 280 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/compat_handlers.py:150:1` | exception handler body is only `pass` |
| 283 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/compat_handlers.py:296:1` | exception handler body is only `pass` |
| 299 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/module_system.py:99:1` | exception handler body is only `pass` |
| 305 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/sentry_integration.py:256:1` | exception handler body is only `pass` |
| 379 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/conftest.py:30:1` | exception handler body is only `pass` |
| 385 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/conftest.py:73:1` | exception handler body is only `pass` |
| 387 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/conftest.py:83:1` | exception handler body is only `pass` |
| 389 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/conftest.py:120:1` | exception handler body is only `pass` |
| 391 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/conftest.py:130:1` | exception handler body is only `pass` |
| 399 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_basic_logging.py:105:1` | exception handler body is only `pass` |
| 433 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_examples.py:259:1` | exception handler body is only `pass` |
| 469 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:178:1` | exception handler body is only `pass` |
| 476 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:231:1` | exception handler body is only `pass` |
| 486 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_overflow_flush.py:44:1` | exception handler body is only `pass` |

### CWE-390 — 9 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 278 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/__init__.py:156:1` | exception detected but handler takes no action (pass-only) |
| 281 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/compat_handlers.py:150:1` | exception detected but handler takes no action (pass-only) |
| 300 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/module_system.py:99:1` | exception detected but handler takes no action (pass-only) |
| 306 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/sentry_integration.py:256:1` | exception detected but handler takes no action (pass-only) |
| 380 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/conftest.py:30:1` | exception detected but handler takes no action (pass-only) |
| 400 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_basic_logging.py:105:1` | exception detected but handler takes no action (pass-only) |
| 434 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_examples.py:259:1` | exception detected but handler takes no action (pass-only) |
| 470 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:178:1` | exception detected but handler takes no action (pass-only) |
| 487 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_overflow_flush.py:44:1` | exception detected but handler takes no action (pass-only) |

### BP-PY-7 — 8 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 6 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/basic_handlers_benchmark.py:366:18` | `open(...)` without a context manager; stream held for handler lifetime |
| 7 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/basic_handlers_benchmark.py:431:18` | `open(...)` without a context manager; stream held for handler lifetime |
| 8 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/basic_handlers_benchmark.py:496:18` | `open(...)` without a context manager; stream held for handler lifetime |
| 9 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/basic_handlers_benchmark.py:538:18` | `open(...)` without a context manager; stream held for handler lifetime |
| 10 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/basic_handlers_benchmark.py:586:18` | `open(...)` without a context manager; stream held for handler lifetime |
| 12 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/compare_loggers.py:108:18` | `open(...)` without a context manager; stream held for handler lifetime |
| 28 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/real_handlers_comparison.py:105:22` | `open(...)` without a context manager; stream held for handler lifetime |
| 29 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/real_handlers_comparison.py:126:18` | `open(...)` without a context manager; stream held for handler lifetime |

### CWE-1071 — 7 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 277 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/__init__.py:156:1` | exception handler silently contains only `pass` |
| 282 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/compat_handlers.py:150:13` | exception handler silently contains only `pass` |
| 301 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/module_system.py:99:5` | exception handler silently contains only `pass` |
| 381 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/conftest.py:30:9` | exception handler silently contains only `pass` |
| 401 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_basic_logging.py:105:9` | exception handler silently contains only `pass` |
| 435 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_examples.py:259:17` | exception handler silently contains only `pass` |
| 471 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:178:17` | exception handler silently contains only `pass` |

### BP-PY-45 — 6 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 3 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/basic_handlers_benchmark.py:51:1` | `sys.path.insert` used to fix local imports instead of packaging |
| 11 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/compare_loggers.py:27:1` | `sys.path.insert` used to fix local imports instead of packaging |
| 17 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/gil_benchmark.py:27:1` | `sys.path.insert` used to fix local imports instead of packaging |
| 20 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/perf_micro.py:39:1` | `sys.path.insert` used to fix local imports instead of packaging |
| 23 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/perf_vs_stdlib.py:28:1` | `sys.path.insert` used to fix local imports instead of packaging |
| 27 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/real_handlers_comparison.py:33:1` | `sys.path.insert` used to fix local imports instead of packaging |

### CWE-396 — 6 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 44 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:215:1` | generic `Exception` handler whose body only logs and continues |
| 79 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:159:1` | generic `Exception` handler whose body only logs and continues |
| 133 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:83:1` | generic `Exception` handler whose body only logs and continues |
| 268 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/third_party_integration.py:50:1` | generic `Exception` handler whose body only logs and continues |
| 339 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:212:1` | generic `Exception` handler whose body only logs and continues |
| 345 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:21:1` | generic `Exception` handler whose body only logs and continues |

### CWE-367 — 2 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 1 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/_bench_common.py:187:12` | filesystem path checked (`os.path.exists`) before a later separate use |
| 440 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_handler_output.py:85:16` | filesystem path checked (`os.path.exists`) before a later separate use |

### CWE-489 — 2 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 36 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:28:1` | debug mode enabled in application source (example apps) |
| 166 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:420:12` | debug mode enabled in application source (example apps) |

### CWE-756 — 2 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 37 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:28:1` | debug error output enabled instead of a custom error page (example apps) |
| 167 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:420:12` | debug error output enabled instead of a custom error page (example apps) |

### BP-PY-16 — 2 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 164 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:420:1` | `app.run(..., debug=True)` in non-test application code |
| 165 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:420:12` | `app.run(..., debug=True)` in non-test application code |

### CWE-1121 — 2 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 302 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/module_system.py:240:27` | function has at least twelve visible control-flow branches |
| 383 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/conftest.py:41:23` | function has at least twelve visible control-flow branches |

### CWE-312 — 1 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 38 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:29:1` | security-sensitive key stored as a cleartext source literal |

### CWE-798 — 1 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 39 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:29:1` | credential assigned directly to a Python source literal |

### BP-PY-26 — 1 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 49 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:282:1` | `@csrf_exempt` on a view that mutates state (POST user creation) |

### CWE-346 — 1 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 89 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:326:1` | CORS allows every origin while credentialed requests are enabled |

### BP-PY-48 — 1 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 90 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:328:5` | `allow_origins=["*"]` combined with `allow_credentials=True` |

### CWE-290 — 1 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 135 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:95:8` | client-provided `X-Forwarded-For` header trusted directly |

### CWE-1084 — 1 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 264 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/third_party_integration.py:13:12` | single `main()` performs many file/data-access operations |

### BP-PY-46 — 1 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 296 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/interceptor.py:31:13` | print used in library code (`logxide/interceptor.py`) |

### BP-PY-8 — 1 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 310 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:19:14` | `subprocess.run(cmd, shell=True)` with a dynamic command parameter |

### CWE-78 — 1 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 311 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/publish.py:19:14` | dynamic command string executed via `subprocess` with `shell=True` |

### BP-PY-12 — 1 findings

| Finding | Source | Reason |
| --- | --- | --- |
| 415 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_doc_codeblocks.py:68:9` | `compile(code, ..., "exec")` on dynamic input (syntax-validation test) |

## Uncertain findings

None. Finding 286 (CWE-1121 on `dictConfig`) reclassified as a false positive: measured branch count is 9 < threshold 12.


## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/logxide/chunks`
- Function evidence: `scripts/logxide/findings/functions`
- Validation: `git diff --check` — pass (exit 0, no whitespace errors)

## Post-fix remaining-FP audit (2026-08-02)

### Run metadata (fresh scan)

```yaml
timestamp: 2026-08-02T16:38:00Z
repository: logxide
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide
branch: main
commit: 136f7a4c3bc593488cd1e2c62bd74956265533d6
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide
chunk_path: scripts/logxide/chunks
function_context_path: scripts/logxide/findings/functions
binary: ./bin/goslop rebuilt from b5b8fde on 2026-08-02 16:29 (post FP-reduction fix)
```

### Scan evidence (fresh run)

- Build command: `n/a` (pre-built binary from fix commit `b5b8fde`, rebuilt 2026-08-02 16:29)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/logxide/chunks -context-dir scripts/logxide/findings/functions real-repos/logxide`
- Findings: `377`
- Chunks reviewed: `scripts/logxide/chunks/Chunk_1_25.txt` .. `scripts/logxide/chunks/Chunk_376_377.txt` (all 15 chunk files)
- Function contexts reviewed: `scripts/logxide/findings/functions/<id>.txt` for all 67 proposed false positives and the 1 new finding

### Audit checklist (fresh run)

- [x] Read every assigned chunk under `scripts/logxide/chunks`.
- [x] Read `scripts/logxide/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

### Classification summary (fresh run)

Fresh findings were matched to the audited run by `Source:` path (`file:line:col`) and rule. A fresh finding whose source matches an audited true positive is a true positive; a fresh finding whose source matches an audited false positive is a remaining false positive. The fix suppressed 186 audited FPs down to 67 remaining.

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 67 | 9, 11, 12, 14, 15, 18, 38, 43, 46, 101, 110, 115, 132, 137, 140, 155, 157, 165, 179, 191, 206, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 224, 225, 230, 231, 234, 235, 241, 250, 262, 277, 285, 291, 294, 295, 296, 299, 302, 303, 307, 315, 316, 317, 321, 332, 335, 336, 342, 365, 366, 372, 374, 375, 376, 377 |
| True positive | 310 | 1, 2, 3, 4, 5, 6, 7, 8, 10, 13, 16, 17, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 39, 40, 41, 42, 44, 45, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 102, 103, 104, 105, 106, 107, 108, 109, 111, 112, 113, 114, 116, 117, 118, 119, 120, 121, 122, 123, 124, 126, 127, 128, 129, 130, 131, 133, 134, 135, 136, 138, 139, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 156, 158, 159, 160, 161, 162, 163, 164, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 207, 208, 209, 210, 223, 226, 227, 228, 229, 232, 233, 236, 237, 238, 239, 240, 242, 243, 244, 245, 246, 247, 248, 249, 251, 252, 253, 254, 255, 256, 257, 258, 259, 260, 261, 263, 264, 265, 266, 267, 268, 269, 270, 271, 272, 273, 274, 275, 276, 278, 279, 280, 281, 282, 283, 284, 286, 287, 288, 289, 290, 292, 293, 297, 298, 300, 301, 304, 305, 306, 308, 309, 310, 311, 312, 313, 314, 318, 319, 320, 322, 323, 324, 325, 326, 327, 328, 329, 330, 331, 333, 334, 337, 338, 339, 340, 341, 343, 344, 345, 346, 347, 348, 349, 350, 351, 352, 353, 354, 355, 356, 357, 358, 359, 360, 361, 362, 363, 364, 367, 368, 369, 370, 371, 373 and 125 |
| Uncertain | 0 | — |

## False positives (remaining, fresh run)

One subsection per remaining false positive (each is a re-appearing audited FP at an unchanged source; fresh finding IDs do not correspond to old IDs). All 67 reference distinct source constructs, so no grouping applies. Source excerpts are taken from the fresh function-context files and match the audited excerpts at identical lines.

### [ ] Finding 9 — CWE-117

- Function context: `scripts/logxide/findings/functions/9.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/compare_loggers.py:138:26`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
        136:         return lambda i: logger.info("Simple log message")
        137:     if scenario == "structured":
    >   138:         return lambda i: logger.info(
        139:             f"User action - user_id: {i}, action: login, status: success"
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 11 — CWE-88

- Function context: `scripts/logxide/findings/functions/11.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/compare_loggers.py:162:12`
- Checklist pattern: dynamic value bound as named-option argument in a no-shell argv

Source excerpt:

```
        160:     env = dict(os.environ)
        161:     env["PYTHONUTF8"] = "1"
    >   162:     proc = subprocess.run(
        163:         [
```

Why this is a false positive: the dynamic values are passed as separate argv elements consumed as the values of their named options (`--library`, `--scenario`, `-n`); they are delimited arguments in a no-shell argv list and cannot be parsed as unintended switches.

Checklist evidence: each value follows its named flag as its own argv element; no shell parsing is involved.

### [ ] Finding 12 — CWE-88

- Function context: `scripts/logxide/findings/functions/12.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/gil_benchmark.py:90:12`
- Checklist pattern: dynamic value bound as named-option argument in a no-shell argv

Source excerpt:

```
         88:     env = dict(os.environ)
         89:     env["PYTHONUTF8"] = "1"
    >    90:     proc = subprocess.run(
         91:         [sys.executable, THIS, "--worker", "--library", library, "-n", str(n)],
```

Why this is a false positive: the dynamic values are passed as separate argv elements consumed as the values of their named options (`--library`, `--scenario`, `-n`); they are delimited arguments in a no-shell argv list and cannot be parsed as unintended switches.

Checklist evidence: each value follows its named flag as its own argv element; no shell parsing is involved.

### [ ] Finding 14 — CWE-117

- Function context: `scripts/logxide/findings/functions/14.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/perf_vs_stdlib.py:73:26`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
         71:         call = lambda i: logger.info("Simple log message")  # noqa: E731
         72:     elif scenario == "structured":
    >    73:         call = lambda i: logger.info(f"User action - user_id: {i}, action: login")  # noqa: E731
         74:     else:
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 15 — CWE-88

- Function context: `scripts/logxide/findings/functions/15.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/perf_vs_stdlib.py:90:12`
- Checklist pattern: dynamic value bound as named-option argument in a no-shell argv

Source excerpt:

```
         88:     env = dict(os.environ)
         89:     env["PYTHONUTF8"] = "1"
    >    90:     proc = subprocess.run(
         91:         [
```

Why this is a false positive: the dynamic values are passed as separate argv elements consumed as the values of their named options (`--library`, `--scenario`, `-n`); they are delimited arguments in a no-shell argv list and cannot be parsed as unintended switches.

Checklist evidence: each value follows its named flag as its own argv element; no shell parsing is involved.

### [ ] Finding 18 — CWE-88

- Function context: `scripts/logxide/findings/functions/18.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/benchmark/real_handlers_comparison.py:154:12`
- Checklist pattern: dynamic value bound as named-option argument in a no-shell argv

Source excerpt:

```
        152:     env = dict(os.environ)
        153:     env["PYTHONUTF8"] = "1"
    >   154:     proc = subprocess.run(
        155:         [
```

Why this is a false positive: the dynamic values are passed as separate argv elements consumed as the values of their named options (`--library`, `--scenario`, `-n`); they are delimited arguments in a no-shell argv list and cannot be parsed as unintended switches.

Checklist evidence: each value follows its named flag as its own argv element; no shell parsing is involved.

### [ ] Finding 38 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/38.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:317:1`
- Checklist pattern: catch-all converts the failure into an HTTP error response

Source excerpt:

```
        315:             return JsonResponse(user.to_dict(), status=201)
        316:
    >   317:         except Exception as e:
        318:             app_logger.error(f"Failed to create user: {str(e)}")
```

Why this is a false positive: the handler body handles the exception by converting it into an HTTP 500 error response after logging — the failure outcome is modeled, so the “without handling” condition fails.

Checklist evidence: the handler returns an error response to the caller after logging.

### [ ] Finding 43 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/43.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:347:1`
- Checklist pattern: health-check catch-all updates service status

Source excerpt:

```
        345:             cursor.execute("SELECT 1")
        346:         db_status = "healthy"
    >   347:     except Exception as e:
        348:         db_logger.error(f"Database health check failed: {str(e)}")
```

Why this is a false positive: the catch-all is the deliberate health-check contract: any failure is handled by logging and reporting `db_status = "unhealthy"` — the failure is processed, not swallowed.

Checklist evidence: the exception is processed into the health-check status result.

### [ ] Finding 46 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/46.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/django_integration.py:405:1`
- Checklist pattern: catch-all converts the failure into an HTTP error response

Source excerpt:

```
        403:         return JsonResponse(metrics_data)
        404:
    >   405:     except Exception as e:
        406:         app_logger.error(f"Failed to retrieve metrics: {str(e)}")
```

Why this is a false positive: the handler body handles the exception by converting it into an HTTP 500 error response after logging — the failure outcome is modeled, so the “without handling” condition fails.

Checklist evidence: the handler returns an error response to the caller after logging.

### [ ] Finding 101 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/101.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:549:1`
- Checklist pattern: health-check catch-all updates service status

Source excerpt:

```
        547:         db.execute("SELECT 1")
        548:         db_status = "healthy"
    >   549:     except Exception as e:
        550:         db_logger.error(f"Database health check failed: {str(e)}")
```

Why this is a false positive: the catch-all is the deliberate health-check contract: any failure is handled by logging and reporting `db_status = "unhealthy"` — the failure is processed, not swallowed.

Checklist evidence: the exception is processed into the health-check status result.

### [ ] Finding 110 — CWE-1046

- Function context: `scripts/logxide/findings/functions/110.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_advanced.py:656:1`
- Checklist pattern: flagged line is an integer increment, not string concatenation

Source excerpt:

```
        654:     while time.time() - start_time < duration:
        655:         performance_logger.info(f"Stress test message {message_count + 1}")
    >   656:         message_count += 1
        657:
```

Why this is a false positive: the flagged construct `message_count += 1` is an integer counter increment; no immutable text is concatenated, so the rule condition (“creation of immutable text using string concatenation”) is not met.

Checklist evidence: the flagged expression is `int += 1`, not `str += str`.

### [ ] Finding 115 — CWE-117

- Function context: `scripts/logxide/findings/functions/115.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/fastapi_demo.py:69:9`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
         67:     # Log many messages using standard logging
         68:     for i in range(count):
    >    69:         logger.info(f"Performance test message {i + 1}")
         70:
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 132 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/132.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:219:1`
- Checklist pattern: catch-all rolls back the transaction and returns an error response

Source excerpt:

```
        217:             return jsonify(user.to_dict()), 201
        218:
    >   219:         except Exception as e:
        220:             db.session.rollback()
```

Why this is a false positive: the handler body handles the failure with `db.session.rollback()` plus an error response — the transaction outcome is explicitly handled.

Checklist evidence: the exception triggers `db.session.rollback()` and a 500 response.

### [ ] Finding 137 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/137.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:249:1`
- Checklist pattern: health-check catch-all updates service status

Source excerpt:

```
        247:         db.session.execute("SELECT 1")
        248:         db_status = "healthy"
    >   249:     except Exception as e:
        250:         db_logger.error(f"Database health check failed: {str(e)}")
```

Why this is a false positive: the catch-all is the deliberate health-check contract: any failure is handled by logging and reporting `db_status = "unhealthy"` — the failure is processed, not swallowed.

Checklist evidence: the exception is processed into the health-check status result.

### [ ] Finding 140 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/140.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/flask_integration.py:304:1`
- Checklist pattern: catch-all converts the failure into an HTTP error response

Source excerpt:

```
        302:         return jsonify(metrics)
        303:
    >   304:     except Exception as e:
        305:         error_logger.error(f"Failed to retrieve metrics: {str(e)}")
```

Why this is a false positive: the handler body handles the exception by converting it into an HTTP 500 error response after logging — the failure outcome is modeled, so the “without handling” condition fails.

Checklist evidence: the handler returns an error response to the caller after logging.

### [ ] Finding 155 — CWE-117

- Function context: `scripts/logxide/findings/functions/155.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/format_production.py:19:5`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
         17:     logger = logging.getLogger(component)
         18:     logger.setLevel(logging.INFO)
    >    19:     logger.info(f"{component} initialized successfully")
         20:     if component == "db.connection":
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 157 — CWE-117

- Function context: `scripts/logxide/findings/functions/157.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/format_threaded.py:21:9`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
         19:
         20:     for i in range(3):
    >    21:         logger.info(f"Processing task {i + 1}")
         22:         time.sleep(0.1)  # Simulate work
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 165 — CWE-117

- Function context: `scripts/logxide/findings/functions/165.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/minimal_dropin.py:120:13`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
        118:
        119:         for i in range(3):
    >   120:             logger.info(f"Processing task {i + 1}")
        121:             time.sleep(0.1)  # Simulate work
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 179 — CWE-117

- Function context: `scripts/logxide/findings/functions/179.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/performance_demo.py:145:9`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
        143:         logger = logging.getLogger(logger_name)
        144:         loggers.append(logger)
    >   145:         logger.info(f"Logger {logger_name} initialized")
        146:
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 191 — CWE-117

- Function context: `scripts/logxide/findings/functions/191.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/examples/simple_demo.py:39:9`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
         37:     start = time.time()
         38:     for i in range(10000):
    >    39:         logger.info(f"Performance test message {i}")
         40:
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 206 — CWE-1333

- Function context: `scripts/logxide/findings/functions/206.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/compat_handlers.py:65:18`
- Checklist pattern: nested repetition alternatives are disjoint and anchored by literals (linear-time)

Source excerpt:

```
         63:         r"^(.?[<>=^])?[+ -]?#?0?(\d+|{\w+})?[,_]?(\.(\d+|{\w+}))?[bcdefgnosx%]?$", re.I
         64:     )
    >    65:     field_spec = re.compile(r"^(\d+|\w+)(\.\w+|\[[^]]+\])*$")
         66:
```

Why this is a false positive: in `(\.\w+|\[[^]]+\])*` every repetition is anchored by a distinct literal delimiter (`.` or `[`) that the inner character classes cannot match, so the alternatives are disjoint and matching is linear — no catastrophic backtracking.

Checklist evidence: the repeated group's alternatives start with disjoint literal delimiters, eliminating ambiguity between iterations.

### [ ] Finding 211 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/211.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/compat_handlers.py:361:1`
- Checklist pattern: catch-all delegates to the logging error contract

Source excerpt:

```
        359:         except RecursionError:
        360:             raise
    >   361:         except Exception:
        362:             self.handleError(record)
```

Why this is a false positive: the `except Exception` is the mandated logging `emit()` contract (emit must never propagate) and the body routes the failure to `self.handleError(record)` — the exception is handled, not swallowed.

Checklist evidence: the handler body delegates to `handleError`, the stdlib emit error path.

### [ ] Finding 212 — CWE-396

- Function context: `scripts/logxide/findings/functions/212.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/compat_handlers.py:361:1`
- Checklist pattern: catch-all delegates to the logging error contract

Source excerpt:

```
        359:         except RecursionError:
        360:             raise
    >   361:         except Exception:
        362:             self.handleError(record)
```

Why this is a false positive: same handler as the paired BP-PY-1 finding: the generic catch is the required `emit()` contract and the error is handled via `self.handleError(record)`.

Checklist evidence: the exception is routed to `handleError`, not hidden.

### [ ] Finding 213 — CWE-1121

- Function context: `scripts/logxide/findings/functions/213.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/config.py:21:24`
- Checklist pattern: branch count below threshold of 12

Source excerpt:

```
def dictConfig(config):
    cfg = copy.deepcopy(config)
    if "handlers" in cfg and isinstance(cfg["handlers"], dict):
        for _name, handler_config in cfg["handlers"].items():
            if not isinstance(handler_config, dict):
                continue
            class_name = handler_config.get("class")
            if class_name == "logxide.FileHandler":
                ...
            elif class_name == "logxide.StreamHandler":
                ...
            elif class_name == "logxide.RotatingFileHandler":
                ...
            elif class_name == "logxide.HTTPHandler":
                ...
            elif class_name == "logxide.OTLPHandler":
                ...
            elif class_name in HANDLER_MAP:
                ...
    logging.config.dictConfig(cfg)
```

Why this is a false positive: CWE-1121 requires `branches >= 12` (`minimumRouteBranches` in `internal/lang/python/detectors/cwe/common.go`). Independent count of `if `/`elif `/`for `/`while `/`except ` tokens in `dictConfig` is **9** (3 `if` + 5 `elif` + 1 `for`) after docstring masking — below the threshold. Reclassified from Uncertain.

Checklist evidence: rule condition "at least twelve visible control-flow branches"; source has 9.

### [ ] Finding 214 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/214.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/fast_logger_wrapper.py:55:1`
- Checklist pattern: catch-all recovers with a fallback value

Source excerpt:

```
         53:             self._effective_level = self._rust_logger.getEffectiveLevel()
         54:             self._name = self._rust_logger.name
    >    55:         except Exception:
         56:             # Fallback to safe defaults
```

Why this is a false positive: the handler body recovers with an explicit fallback value (safe defaults, `str(record.msg)`, `old_level = None`) — the failure is handled, satisfying the rule's “without handling” clause.

Checklist evidence: the handler body performs explicit fallback recovery.

### [ ] Finding 215 — CWE-396

- Function context: `scripts/logxide/findings/functions/215.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/fast_logger_wrapper.py:55:1`
- Checklist pattern: catch-all recovers with a fallback value

Source excerpt:

```
         53:             self._effective_level = self._rust_logger.getEffectiveLevel()
         54:             self._name = self._rust_logger.name
    >    55:         except Exception:
         56:             # Fallback to safe defaults
```

Why this is a false positive: same recovery handler as the paired BP-PY-1 finding: a deliberate fallback converts the exception into a safe default state.

Checklist evidence: the exception is converted into a fallback value/state.

### [ ] Finding 216 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/216.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/handlers.py:191:1`
- Checklist pattern: catch-all delegates to the logging error contract

Source excerpt:

```
        189:                     record.args = None
        190:                 self._inner.emit(_prepare_record_for_rust(record))
    >   191:         except Exception:
        192:             self.handleError(record)
```

Why this is a false positive: the `except Exception` is the mandated logging `emit()` contract (emit must never propagate) and the body routes the failure to `self.handleError(record)` — the exception is handled, not swallowed.

Checklist evidence: the handler body delegates to `handleError`, the stdlib emit error path.

### [ ] Finding 217 — CWE-396

- Function context: `scripts/logxide/findings/functions/217.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/handlers.py:191:1`
- Checklist pattern: catch-all delegates to the logging error contract

Source excerpt:

```
        189:                     record.args = None
        190:                 self._inner.emit(_prepare_record_for_rust(record))
    >   191:         except Exception:
        192:             self.handleError(record)
```

Why this is a false positive: same handler as the paired BP-PY-1 finding: the generic catch is the required `emit()` contract and the error is handled via `self.handleError(record)`.

Checklist evidence: the exception is routed to `handleError`, not hidden.

### [ ] Finding 218 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/218.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/handlers.py:260:1`
- Checklist pattern: catch-all delegates to the logging error contract

Source excerpt:

```
        258:                     record.args = None
        259:                 self._inner.emit(_prepare_record_for_rust(record))
    >   260:         except Exception:
        261:             self.handleError(record)
```

Why this is a false positive: the `except Exception` is the mandated logging `emit()` contract (emit must never propagate) and the body routes the failure to `self.handleError(record)` — the exception is handled, not swallowed.

Checklist evidence: the handler body delegates to `handleError`, the stdlib emit error path.

### [ ] Finding 219 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/219.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/handlers.py:325:1`
- Checklist pattern: catch-all delegates to the logging error contract

Source excerpt:

```
        323:                     record.args = None
        324:                 self._inner.emit(_prepare_record_for_rust(record))
    >   325:         except Exception:
        326:             self.handleError(record)
```

Why this is a false positive: the `except Exception` is the mandated logging `emit()` contract (emit must never propagate) and the body routes the failure to `self.handleError(record)` — the exception is handled, not swallowed.

Checklist evidence: the handler body delegates to `handleError`, the stdlib emit error path.

### [ ] Finding 220 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/220.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/handlers.py:406:1`
- Checklist pattern: catch-all delegates to the logging error contract

Source excerpt:

```
        404:             rust_record = _prepare_record_for_rust(record)
        405:             self._inner.emit(rust_record)
    >   406:         except Exception:
        407:             self.handleError(record)
```

Why this is a false positive: the `except Exception` is the mandated logging `emit()` contract (emit must never propagate) and the body routes the failure to `self.handleError(record)` — the exception is handled, not swallowed.

Checklist evidence: the handler body delegates to `handleError`, the stdlib emit error path.

### [ ] Finding 221 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/221.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/handlers.py:477:1`
- Checklist pattern: catch-all delegates to the logging error contract

Source excerpt:

```
        475:             rust_record = _prepare_record_for_rust(record)
        476:             self._inner.emit(rust_record)
    >   477:         except Exception:
        478:             self.handleError(record)
```

Why this is a false positive: the `except Exception` is the mandated logging `emit()` contract (emit must never propagate) and the body routes the failure to `self.handleError(record)` — the exception is handled, not swallowed.

Checklist evidence: the handler body delegates to `handleError`, the stdlib emit error path.

### [ ] Finding 222 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/222.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/handlers.py:519:1`
- Checklist pattern: catch-all delegates to the logging error contract

Source excerpt:

```
        517:             # MemoryHandler is always native: forward raw; caplog reads _inner.
        518:             self._inner.emit(_prepare_record_for_rust(record, native=True))
    >   519:         except Exception:
        520:             self.handleError(record)
```

Why this is a false positive: the `except Exception` is the mandated logging `emit()` contract (emit must never propagate) and the body routes the failure to `self.handleError(record)` — the exception is handled, not swallowed.

Checklist evidence: the handler body delegates to `handleError`, the stdlib emit error path.

### [ ] Finding 224 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/224.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/interceptor.py:42:1`
- Checklist pattern: catch-all recovers with a fallback value

Source excerpt:

```
         40:             try:
         41:                 message = record.getMessage()
    >    42:             except Exception:
         43:                 message = str(record.msg)
```

Why this is a false positive: the handler body recovers with an explicit fallback value (safe defaults, `str(record.msg)`, `old_level = None`) — the failure is handled, satisfying the rule's “without handling” clause.

Checklist evidence: the handler body performs explicit fallback recovery.

### [ ] Finding 225 — CWE-396

- Function context: `scripts/logxide/findings/functions/225.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/interceptor.py:42:1`
- Checklist pattern: catch-all recovers with a fallback value

Source excerpt:

```
         40:             try:
         41:                 message = record.getMessage()
    >    42:             except Exception:
         43:                 message = str(record.msg)
```

Why this is a false positive: same recovery handler as the paired BP-PY-1 finding: a deliberate fallback converts the exception into a safe default state.

Checklist evidence: the exception is converted into a fallback value/state.

### [ ] Finding 230 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/230.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/sentry_integration.py:85:1`
- Checklist pattern: catch-all delegates to the error handler contract

Source excerpt:

```
         83:                 self._add_breadcrumb(record, level_name, message, logger_name)
         84:
    >    85:         except Exception as e:
         86:             # Prevent logging errors from causing infinite loops
```

Why this is a false positive: the handler body routes the exception to `self._handle_error(e)` (guarding against recursive logging errors) — handled, not swallowed.

Checklist evidence: the handler body delegates to the dedicated `_handle_error` path.

### [ ] Finding 231 — CWE-396

- Function context: `scripts/logxide/findings/functions/231.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/sentry_integration.py:85:1`
- Checklist pattern: catch-all delegates to the error handler contract

Source excerpt:

```
         83:                 self._add_breadcrumb(record, level_name, message, logger_name)
         84:
    >    85:         except Exception as e:
         86:             # Prevent logging errors from causing infinite loops
```

Why this is a false positive: same handler as finding 303: deliberate catch-all in the Sentry emit path with explicit `_handle_error` handling.

Checklist evidence: the exception is routed to `_handle_error`, not hidden.

### [ ] Finding 234 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/234.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/testing.py:163:1`
- Checklist pattern: catch-all recovers with a fallback value

Source excerpt:

```
        161:         try:
        162:             old_level = handler._inner.level if hasattr(handler, "_inner") else None
    >   163:         except Exception:
        164:             old_level = None
```

Why this is a false positive: the handler body recovers with an explicit fallback value (safe defaults, `str(record.msg)`, `old_level = None`) — the failure is handled, satisfying the rule's “without handling” clause.

Checklist evidence: the handler body performs explicit fallback recovery.

### [ ] Finding 235 — CWE-396

- Function context: `scripts/logxide/findings/functions/235.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/logxide/testing.py:163:1`
- Checklist pattern: catch-all recovers with a fallback value

Source excerpt:

```
        161:         try:
        162:             old_level = handler._inner.level if hasattr(handler, "_inner") else None
    >   163:         except Exception:
        164:             old_level = None
```

Why this is a false positive: same recovery handler as the paired BP-PY-1 finding: a deliberate fallback converts the exception into a safe default state.

Checklist evidence: the exception is converted into a fallback value/state.

### [ ] Finding 241 — BP-PY-42

- Function context: `scripts/logxide/findings/functions/241.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:16:1`
- Checklist pattern: try/except is a defensive guard, not a failure expectation

Source excerpt:

```
         14:     print("Testing basic import...")
         15:
    >    16:     try:
         17:         import logxide
```

Why this is a false positive: the try/except is a defensive guard that reports failure to the console so the smoke script can continue; the test does not “expect failure”, so the rule condition (try/except instead of `pytest.raises`) is not met.

Checklist evidence: the except branch reports the failure rather than asserting a raise.

### [ ] Finding 250 — CWE-117

- Function context: `scripts/logxide/findings/functions/250.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/scripts/verify_package.py:94:13`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
         92:         start_time = time.time()
         93:         for i in range(1000):
    >    94:             logger.info(f"Performance test message {i}")
         95:         logging.flush()
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 262 — BP-PY-40

- Function context: `scripts/logxide/findings/functions/262.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/conftest.py:36:6`
- Checklist pattern: daemon thread with explicit synchronization, no join needed

Source excerpt:

```
         34:
         35:     t = threading.Thread(target=_do_flush, daemon=True)
    >    36:     t.start()
         37:     done.wait(timeout=timeout_seconds)
```

Why this is a false positive: the thread is `daemon=True` (documented as a deliberate deadlock-bailout) and the caller synchronizes with `done.wait(timeout)`; the rule's own fix targets “fire-and-forget non-daemon threads”, so the condition is not met.

Checklist evidence: daemon thread plus explicit `Event.wait` synchronization protocol.

### [ ] Finding 277 — BP-PY-42

- Function context: `scripts/logxide/findings/functions/277.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_basic_logging.py:97:1`
- Checklist pattern: try/except is a crash-guard, not a failure expectation

Source excerpt:

```
         95:
         96:         # Test various format configurations - just make sure they don't crash
    >    97:         try:
         98:             logging.basicConfig(format="%(levelname)s: %(message)s")
```

Why this is a false positive: the try/except (“just make sure they don't crash”) guards a smoke test against unsupported formats rather than expecting a raise; `pytest.raises` is the wrong construct, so the condition is not met.

Checklist evidence: the except branch tolerates unsupported formats instead of asserting an exception.

### [ ] Finding 285 — CWE-117

- Function context: `scripts/logxide/findings/functions/285.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_basic_logging.py:147:17`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
        145:
        146:             for i in range(5):
    >   147:                 logger.info(f"Message {i}")
        148:                 time.sleep(0.001)  # Small delay
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 291 — BP-PY-1

- Function context: `scripts/logxide/findings/functions/291.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_compatibility.py:12:1`
- Checklist pattern: catch-all delegates to full exception reporting

Source excerpt:

```
         10:     try:
         11:         raise ValueError("test")
    >    12:     except:
         13:         logging.exception("test_exception_logging")
```

Why this is a false positive: the handler body reports the (deliberately raised) exception through the full exception-reporting path (`logger.exception` / `error(..., exc_info=True)`, forwarded to Sentry) — handled, not swallowed.

Checklist evidence: the handler body is a full exception-reporting call (`logger.exception`/`exc_info=True`).

### [ ] Finding 294 — CWE-94

- Function context: `scripts/logxide/findings/functions/294.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_doc_codeblocks.py:68:9`
- Checklist pattern: compile input is the repo's own committed docs, not externally influenced

Source excerpt:

```
         66:     """Every Python code block in docs must be syntactically valid."""
         67:     try:
    >    68:         compile(code, block_id, "exec")
         69:     except SyntaxError as exc:
```

Why this is a false positive: the compiled text comes from the repository's own documentation files and the resulting code object is discarded after syntax validation (never executed); it is not externally influenced input, so the CWE-94 condition is not met.

Checklist evidence: input source is repo-committed docs; the compile result is never executed.

### [ ] Finding 295 — CWE-829

- Function context: `scripts/logxide/findings/functions/295.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_doc_codeblocks.py:215:13`
- Checklist pattern: dynamic import selects from repo-derived framework names, not request input

Source excerpt:

```
        213:     for framework in needed:
        214:         try:
    >   215:             __import__(framework)
        216:         except ImportError:
```

Why this is a false positive: `needed` is computed by scanning the repository's own committed doc code blocks (`_needs_framework(code)`); the imported names come from package-controlled content, not from request-derived input, so the untrusted-control-sphere condition is not met.

Checklist evidence: module names originate from repo-committed code content, not user/request input.

### [ ] Finding 296 — CWE-829

- Function context: `scripts/logxide/findings/functions/296.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_examples.py:25:16`
- Checklist pattern: file path comes from the repo's own examples directory, not request input

Source excerpt:

```
         23:         # Load the module
         24:         module_name = example_path.stem
    >    25:         spec = importlib.util.spec_from_file_location(module_name, example_path)
         26:         if spec is None:
```

Why this is a false positive: `example_path` is iterated from the repository's own `examples/` directory; the rule's own fix (“never pass request-derived module names or paths”) does not apply to package-controlled repo files.

Checklist evidence: the path is derived from a local directory listing of repo files, not external input.

### [ ] Finding 299 — BP-PY-42

- Function context: `scripts/logxide/findings/functions/299.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_examples.py:55:1`
- Checklist pattern: try/except is a defensive guard, not a failure expectation

Source excerpt:

```
         53:     print("Testing LogXide basic functionality...")
         54:
    >    55:     try:
         56:         # Use auto-install pattern
```

Why this is a false positive: the try/except is a defensive guard that reports failure to the console so the smoke script can continue; the test does not “expect failure”, so the rule condition (try/except instead of `pytest.raises`) is not met.

Checklist evidence: the except branch reports the failure rather than asserting a raise.

### [ ] Finding 302 — CWE-489

- Function context: `scripts/logxide/findings/functions/302.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_examples.py:137:1`
- Checklist pattern: debug flag set inside a test, not application source

Source excerpt:

```
        135:         if not settings.configured:
        136:             settings.configure(
    >   137:                 DEBUG=True,
        138:                 SECRET_KEY="test-key-for-integration-test",
```

Why this is a false positive: the flagged debug flag configures Django/Sentry inside a pytest test under `tests/`, where debug mode is standard test scaffolding; the rule condition (“debug mode is enabled in application source”) is not met.

Checklist evidence: the construct is test code configuring a fake settings/SDK object, not application source.

### [ ] Finding 303 — CWE-756

- Function context: `scripts/logxide/findings/functions/303.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_examples.py:137:1`
- Checklist pattern: debug flag set inside a test, not a deployed application

Source excerpt:

```
        135:         if not settings.configured:
        136:             settings.configure(
    >   137:                 DEBUG=True,
        138:                 SECRET_KEY="test-key-for-integration-test",
```

Why this is a false positive: the same test scaffolding as the paired CWE-489 finding; the “application explicitly enables debug error output” condition targets deployed application source, not test code.

Checklist evidence: the construct is test code, not application source.

### [ ] Finding 307 — CWE-117

- Function context: `scripts/logxide/findings/functions/307.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_examples.py:237:13`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
        235:
        236:         for i in range(1000):
    >   237:             logger.info(f"Performance test message {i}")
        238:
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 315 — CWE-117

- Function context: `scripts/logxide/findings/functions/315.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_gil_scaling.py:123:9`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
        121:     n = 25
        122:     for i in range(n):
    >   123:         logger.info(f"secret-{i}")
        124:
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 316 — BP-PY-40

- Function context: `scripts/logxide/findings/functions/316.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_http_auth.py:41:11`
- Checklist pattern: daemon server thread with explicit shutdown protocol

Source excerpt:

```
         39:     port = server.server_address[1]
         40:     thread = threading.Thread(target=server.serve_forever, daemon=True)
    >    41:     thread.start()
         42:     yield port
```

Why this is a false positive: the fixture thread is `daemon=True` and teardown calls `server.shutdown()` — a clear shutdown protocol, satisfying the rule's own fix (avoid non-daemon fire-and-forget).

Checklist evidence: daemon thread plus explicit `shutdown()` on teardown.

### [ ] Finding 317 — BP-PY-40

- Function context: `scripts/logxide/findings/functions/317.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_http_handler_features.py:43:11`
- Checklist pattern: daemon server thread with explicit shutdown protocol

Source excerpt:

```
         41:     port = server.server_address[1]
         42:     thread = threading.Thread(target=server.serve_forever, daemon=True)
    >    43:     thread.start()
         44:     yield port
```

Why this is a false positive: the fixture thread is `daemon=True` and teardown calls `server.shutdown()` — a clear shutdown protocol, satisfying the rule's own fix (avoid non-daemon fire-and-forget).

Checklist evidence: daemon thread plus explicit `shutdown()` on teardown.

### [ ] Finding 321 — CWE-117

- Function context: `scripts/logxide/findings/functions/321.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration.py:105:17`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
        103:             logger = logging.getLogger("app.database")
        104:             for i in range(10):
    >   105:                 logger.info(f"Database query {i}")
        106:                 time.sleep(0.001)
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 332 — BP-PY-42

- Function context: `scripts/logxide/findings/functions/332.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:20:1`
- Checklist pattern: try/except is a defensive guard, not a failure expectation

Source excerpt:

```
         18:     print("Testing Flask integration...")
         19:
    >    20:     try:
         21:         # Import and test basic functionality
```

Why this is a false positive: the try/except is a defensive guard that reports failure to the console so the smoke script can continue; the test does not “expect failure”, so the rule condition (try/except instead of `pytest.raises`) is not met.

Checklist evidence: the except branch reports the failure rather than asserting a raise.

### [ ] Finding 335 — CWE-489

- Function context: `scripts/logxide/findings/functions/335.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:69:1`
- Checklist pattern: debug flag set inside a test, not application source

Source excerpt:

```
         67:         if not settings.configured:
         68:             settings.configure(
    >    69:                 DEBUG=True,
         70:                 SECRET_KEY="test-key",
```

Why this is a false positive: the flagged debug flag configures Django/Sentry inside a pytest test under `tests/`, where debug mode is standard test scaffolding; the rule condition (“debug mode is enabled in application source”) is not met.

Checklist evidence: the construct is test code configuring a fake settings/SDK object, not application source.

### [ ] Finding 336 — CWE-756

- Function context: `scripts/logxide/findings/functions/336.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:69:1`
- Checklist pattern: debug flag set inside a test, not a deployed application

Source excerpt:

```
         67:         if not settings.configured:
         68:             settings.configure(
    >    69:                 DEBUG=True,
         70:                 SECRET_KEY="test-key",
```

Why this is a false positive: the same test scaffolding as the paired CWE-489 finding; the “application explicitly enables debug error output” condition targets deployed application source, not test code.

Checklist evidence: the construct is test code, not application source.

### [ ] Finding 342 — CWE-117

- Function context: `scripts/logxide/findings/functions/342.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_integration_examples.py:159:13`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
        157:
        158:         for i in range(1000):
    >   159:             logger.info(f"Performance test message {i}")
        160:
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 365 — CWE-489

- Function context: `scripts/logxide/findings/functions/365.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_sentry_integration.py:353:1`
- Checklist pattern: debug flag set inside a test, not application source

Source excerpt:

```
        351:         sentry_sdk.init(
        352:             dsn="https://test@example.com/1",  # Valid format but fake
    >   353:             debug=True,
        354:             # Disable default integrations to avoid noise
```

Why this is a false positive: the flagged debug flag configures Django/Sentry inside a pytest test under `tests/`, where debug mode is standard test scaffolding; the rule condition (“debug mode is enabled in application source”) is not met.

Checklist evidence: the construct is test code configuring a fake settings/SDK object, not application source.

### [ ] Finding 366 — CWE-756

- Function context: `scripts/logxide/findings/functions/366.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_sentry_integration.py:353:1`
- Checklist pattern: debug flag set inside a test, not a deployed application

Source excerpt:

```
        351:         sentry_sdk.init(
        352:             dsn="https://test@example.com/1",  # Valid format but fake
    >   353:             debug=True,
        354:             # Disable default integrations to avoid noise
```

Why this is a false positive: the same test scaffolding as the paired CWE-489 finding; the “application explicitly enables debug error output” condition targets deployed application source, not test code.

Checklist evidence: the construct is test code, not application source.

### [ ] Finding 372 — CWE-117

- Function context: `scripts/logxide/findings/functions/372.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tests/test_timed_rotation.py:70:13`
- Checklist pattern: interpolated value is not externally influenced input

Source excerpt:

```
         68:         # Create multiple rotations
         69:         for i in range(4):
    >    70:             logger.info(f"message {i}")
         71:             handler.flush()
```

Why this is a false positive: the value formatted into the log message is a local loop counter, an inline literal, or an internally generated name; there is no externally influenced value, so the CWE-117 “externally influenced input” condition is not met.

Checklist evidence: the interpolated expression derives from a local counter/literal/internal variable, not request, network, or environment data.

### [ ] Finding 374 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/374.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tmp_test_intercept.py:10:1:tmp_test_intercept.py:10:1`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
          8: lock = threading.RLock()
          9:
    >    10: print("Calling outside lock...")
         11: lx_logger.log(logging.INFO, "Outside")
```

Why this is a false positive: the print is output of a throwaway scratch script at the repo root, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 375 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/375.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tmp_test_intercept.py:13:1:tmp_test_intercept.py:13:1`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         11: lx_logger.log(logging.INFO, "Outside")
         12:
    >    13: print("Acquiring lock...")
         14: with lock:
```

Why this is a false positive: the print is output of a throwaway scratch script at the repo root, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 376 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/376.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tmp_test_intercept.py:15:5:tmp_test_intercept.py:15:5`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         13: print("Acquiring lock...")
         14: with lock:
    >    15:     print("Calling inside lock...")
         16:     lx_logger.log(logging.INFO, "Inside lock!")
```

Why this is a false positive: the print is output of a throwaway scratch script at the repo root, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

### [ ] Finding 377 — BP-PY-46

- Function context: `scripts/logxide/findings/functions/377.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide/tmp_test_intercept.py:18:1:tmp_test_intercept.py:18:1`
- Checklist pattern: print in a script module, not a library module

Source excerpt:

```
         16:     lx_logger.log(logging.INFO, "Inside lock!")
         17:
    >    18: print("FINISHED!")
         19:
```

Why this is a false positive: the print is output of a throwaway scratch script at the repo root, not operational logging in a library module; the rule condition (“print used for operational logging in non-script modules”) is not met.

Checklist evidence: the file is a runnable script (the rule's own exception for CLI/demo scripts), not library code.

## Uncertain findings (fresh run)

None.

Note: fresh finding 125 (`examples/flask_integration.py:137:1`, CWE-396) is a new finding not present in the audited run. It is a true positive: the generic `except Exception` body only calls `db_logger.error(...)` and continues, matching the audited CWE-396 true-positive pattern (e.g., audited findings 44, 79) — the same handler the audited run classified as a BP-PY-1 true positive at the identical source.

### Final evidence (fresh run)

- Delegated reviewers: none
- Chunk evidence: `scripts/logxide/chunks` (all 15 files)
- Function evidence: `scripts/logxide/findings/functions` (context files for all 67 FPs and finding 125)
- Validation: `git diff --check` — pass (exit 0, no whitespace errors)

## Post-fix v2 audit (latest binary)

### Run metadata (fresh scan, latest binary)

```yaml
timestamp: 2026-08-02T18:00:00Z
repository: logxide
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide
branch: main
commit: 136f7a4c3bc593488cd1e2c62bd74956265533d6
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/logxide
chunk_path: scripts/logxide/chunks
function_context_path: scripts/logxide/findings/functions
binary: ./bin/goslop rebuilt 2026-08-02 17:56 (latest, post-b5b8fde; detector worktree changes uncommitted)
```

### Scan evidence (fresh run)

- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/logxide/chunks -context-dir scripts/logxide/findings/functions real-repos/logxide`
- Findings: `408`
- Chunks reviewed: `scripts/logxide/chunks/Chunk_1_25.txt` .. `scripts/logxide/chunks/Chunk_401_408.txt` (all 17 chunk files)
- Function contexts reviewed: `scripts/logxide/findings/functions/<id>.txt` (408 files; source files followed up where the excerpt was insufficient)

### Classification summary (fresh run)

Fresh findings matched to prior audits by `Source:` (`file:line:col`) + rule: audited TP → TP, audited FP (original 503-run and/or b5b8fde append) → FP reusing the prior reason. 408 fresh = 316 TP + 92 FP + 0 Uncertain + 0 new (every fresh (rule, source) pair has a prior classification; note fresh finding 133 `examples/flask_integration.py:137:1` BP-PY-1 matches the audited TP — the CWE-396 twin from the b5b8fde run no longer fires, so no new finding).

| Classification | Count |
| --- | ---: |
| False positive | 92 |
| True positive | 316 |
| Uncertain | 0 |
| New (no prior classification) | 0 |

## Fix checklist (FP patterns)

| Pattern # | Rule | Trigger shape | Count | Example sources |
| --- | --- | --- | ---: | --- |
| 1 | BP-PY-46 | `print(...)` in a runnable script module (examples/ demos, repo-root scratch `tmp_test_intercept.py`); condition: file is a script/demo (not an imported library module) | 26 | examples/minimal_dropin.py:15:5, examples/format_minimal.py:3:1, tmp_test_intercept.py:10:1 |
| 2 | BP-PY-1 (7) + CWE-396 (2) | `except Exception:` in `Handler.emit()` whose body is exactly `self.handleError(record)` (stdlib emit contract, often after `except RecursionError: raise`); condition: body routes to handleError vs empty/pass | 9 | logxide/handlers.py:191:1, logxide/compat_handlers.py:361:1 |
| 3 | BP-PY-1 (3) + CWE-396 (3) | `except Exception:` with explicit fallback assignment (`old_level = None`, `message = str(record.msg)`, safe defaults); condition: fallback value assigned vs bare swallow | 6 | logxide/interceptor.py:42:1, logxide/fast_logger_wrapper.py:55:1, logxide/testing.py:163:1 |
| 4 | BP-PY-1 | catch-all in web handler converting failure into a modeled outcome: HTTP 500 response (3), health-check status (3), `db.session.rollback()` (1); condition: handler body models the failure outcome | 7 | examples/django_integration.py:317:1, examples/flask_integration.py:219:1 |
| 5 | BP-PY-1 (2) + CWE-396 (2) | catch-all delegates to full exception reporting: `logger.exception(...)`/`exc_info=True` or `self._handle_error(e)` | 4 | examples/sentry_integration.py:111:1, logxide/sentry_integration.py:85:1, tests/test_compatibility.py:12:1 |
| 6 | CWE-396 | catch-all converts failure into a recorded benchmark result (`result.ok = False`/`result.error`/traceback) | 2 | benchmark/basic_handlers_benchmark.py:117:1, benchmark/perf_micro.py:79:1 |
| 7 | CWE-117 | f-string interpolated into log message where the value is a local loop counter / inline literal / internally generated name (no request/network/env input); condition: interpolated expression is externally influenced | 15 | benchmark/compare_loggers.py:138:26, tests/test_integration.py:105:17 |
| 8 | CWE-88 | dynamic value passed as its own argv element as the value of a named option (`--library`, `--scenario`, `-n`) in list-form `subprocess.run` (no shell); condition: named-option value vs free operand | 4 | benchmark/compare_loggers.py:162:12 |
| 9 | BP-PY-42 | try/except guarding an import/smoke check in test scripts where the except branch reports failure instead of asserting a raise (defensive guard, not expected-failure); condition: `pytest.raises`-shaped expectation | 4 | tests/test_examples.py:55:1, scripts/verify_package.py:16:1 |
| 10 | BP-PY-40 | `Thread(..., daemon=True).start()` with explicit synchronization (`Event.wait(timeout)`, `server.shutdown()` on teardown); condition: non-daemon fire-and-forget | 3 | tests/conftest.py:36:6, tests/test_http_auth.py:41:11 |
| 11 | CWE-489 (3) + CWE-756 (3) | `DEBUG=True` inside pytest test scaffolding (`settings.configure(...)`, `sentry_sdk.init(debug=True)`); condition: deployed application source, not test code | 6 | tests/test_examples.py:137:1, tests/test_sentry_integration.py:353:1 |
| 12 | CWE-829 | dynamic `__import__`/importlib load where the name/path derives from repo-committed content (docs-code scan, `examples/` listing); condition: untrusted control sphere | 2 | tests/test_doc_codeblocks.py:215:13, tests/test_examples.py:25:16 |
| 13 | CWE-94 | `compile(code, ...)` where input is the repo's own committed docs and the code object is never executed | 1 | tests/test_doc_codeblocks.py:68:9 |
| 14 | CWE-1046 | `message_count += 1` integer increment misread as immutable string concatenation | 1 | examples/fastapi_advanced.py:656:1 |
| 15 | CWE-1121 | `dictConfig` measured at 9 control-flow branches, below the rule threshold of 12 | 1 | logxide/config.py:21:24 |
| 16 | CWE-1333 | `(\.\w+|\[[^]]+\])*` repetition anchored by disjoint literal delimiters (linear-time); condition: ambiguous overlapping repetition | 1 | logxide/compat_handlers.py:65:18 |

## New findings

None — all 408 fresh (rule, source) pairs match a prior audit classification (316 TP, 92 FP). No fresh finding lacks a prior classification.

### Final evidence (v2)

- Delegated reviewers: none
- Chunk evidence: `scripts/logxide/chunks` (all 17 files)
- Function evidence: `scripts/logxide/findings/functions` (all 408 context files)
- Validation: `git diff --check` — pass (exit 0, no whitespace errors)
