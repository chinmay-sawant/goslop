# False-positive audit report — CourtScrapper

## Run metadata

```yaml
timestamp: 2026-08-02T07:53:28Z
repository: CourtScrapper
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/CourtScrapper
branch: master
commit: c1377ced4655f51eb59d685959080fd1a55f03af
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/CourtScrapper
chunk_path: ./scripts/chunks
function_context_path: ./scripts/findings/functions
```

## Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop` (implied by `./bin/goslop`)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir real-repos/CourtScrapper/scripts/chunks -context-dir real-repos/CourtScrapper/scripts/findings/functions real-repos/CourtScrapper`
- Findings: `340`
- Chunks reviewed: `./scripts/chunks/Chunk_1_25.txt` … `Chunk_326_340.txt` (14 files, all findings 1–340)
- Function contexts reviewed: `./scripts/findings/functions/26.txt, 102.txt, 127.txt, 149.txt, 152.txt, 264.txt, 276.txt, 280.txt, 285.txt, 298.txt, 302.txt, 305.txt`

## Audit checklist

- [x] Read every assigned chunk under `./scripts/chunks`.
- [x] Read `./scripts/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient (e.g. `scraper.py` cleanup block for CWE-1341, `utils.py` port/temp-dir derivation for CWE-88, `scraper.py:67-69` re-raise for CWE-396, `main.py` executor call site for PERF-PY-28).
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient. (No delegated reviews; no disagreements.)
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 12 | 26, 102, 127, 149, 152, 264, 276, 280, 285, 298, 302, 305 |
| True positive | 328 | 1-25 (minus 26), 27-101 (minus 102), 103-126 (minus 127), 128-148 (minus 149), 150-151 (minus 152), 153-263 (minus 264), 265-275 (minus 276), 277-279 (minus 280), 281-284 (minus 285), 286-297 (minus 298), 299-301 (minus 302), 303-304 (minus 305), 306-340 |
| Uncertain | 0 | — |

Breakdown of true positives by rule: BP-PY-1 (74), BP-PY-2 (11), BP-PY-46 (26), BP-PY-47 (192), CWE-1071 (4), CWE-1121 (6), CWE-1124 (4), CWE-117 (3: 41, 62, 140), CWE-390 (4), CWE-396 (4: 9, 60, 122, 129).

## False positives

### [ ] Finding `26` — `CWE-779`

- Function context: `./scripts/findings/functions/26.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/CourtScrapper/captcha_handler.py:229:17`
- Checklist pattern: sensitive value is **not** logged; only a non-sensitive derived quantity is.

Source excerpt:

```
    logger.info("Captcha solved! Token received (length: %d)", len(token))
```

Why this is a false positive: the rule condition requires a "sensitive value … logged without a visible redaction", but the only data written is `len(token)` — the captcha token itself is never logged, only its length.

Checklist evidence: the call's argument is `len(token)` (an integer length), so the sensitive value (the token string) is not logged at all; the regex signal is the literal word "Token" in the static message text.

### [ ] Finding `102` — `CWE-117`

- Function context: `./scripts/findings/functions/102.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/CourtScrapper/inspect_website.py:14:5`
- Checklist pattern: no dynamic value is interpolated at all.

Source excerpt:

```
def inspect_page(page, page_name):
    logger.info(f"\n{'='*60}")
```

Why this is a false positive: the rule condition requires a "dynamic value … without CRLF neutralization", but the f-string interpolates only the constant expression `'='*60` — there is no runtime value that could carry CR or LF.

Checklist evidence: the entire interpolated content is a constant (string multiplication of a literal), so no externally controlled CRLF-capable value reaches the log.

### [ ] Finding `127` — `CWE-117`

- Function context: `./scripts/findings/functions/127.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/CourtScrapper/main.py:30:9`
- Checklist pattern: interpolated value is a numeric length of a static config constant.

Source excerpt:

```
logger.info(f"Starting concurrent scraping for {len(ATTORNEYS)} attorney(s)...")
```

Why this is a false positive: the only interpolated value is `len(ATTORNEYS)`, an integer derived from a local config list; an integer cannot contain CRLF and is not externally controlled.

Checklist evidence: the dynamic segment is a `len()` call over a developer-owned constant; no CR/LF-carrying input is formatted into the message.

### [ ] Finding `149` — `CWE-396`

- Function context: `./scripts/findings/functions/149.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/CourtScrapper/scraper.py:67:1`
- Checklist pattern: handler does not hide failures — it logs and re-raises.

Source excerpt:

```
        except Exception as e:
            logger.error(f"Error setting up browser: {e}")
            raise
```

Why this is a false positive: the rule condition is a generic handler that "can hide distinct failure conditions"; here the handler logs the exception and immediately re-raises it, so no failure condition is hidden or swallowed.

Checklist evidence: the handler body ends in a bare `raise` (scraper.py:69), propagating the original exception — the generic catch is a log-and-re-raise wrapper.

### [ ] Finding `152` — `CWE-117`

- Function context: `./scripts/findings/functions/152.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/CourtScrapper/scraper.py:82:13`
- Checklist pattern: interpolated values are a numeric delay and an internal literal description.

Source excerpt:

```
        if description:
            logger.debug(f"Waiting {delay}s before {description}")
        time.sleep(delay)
```

Why this is a false positive: `delay` is a numeric duration and `description` is a developer-supplied literal passed by callers (e.g. "navigating to search page"); neither is externally controlled or CRLF-capable.

Checklist evidence: both interpolated segments are internal scalar values; there is no input-derived string that could inject log record separators.

### [ ] Finding `264` — `CWE-1341`

- Function context: `./scripts/findings/functions/264.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/CourtScrapper/scraper.py:881:26`
- Checklist pattern: distinct handles, each released exactly once.

Source excerpt:

```
            if self.page:
                try:
                    self.page.close()
                except:
                    pass
                self.page = None
            if self.context:
                try:
                    self.context.close()
                except:
                    pass
                self.context = None
```

Why this is a false positive: the rule condition is "same resource handle is released twice", but the matched `close()` calls (lines 881, 887, 893) release four **different** handles — `page`, `context`, `browser` — each exactly once, which is normal teardown.

Checklist evidence: the adjacent `\w+\.close()` pairs are `self.page.close()` vs `self.context.close()`; no two calls operate on the same handle.

### [ ] Finding `276` — `CWE-117`

- Function context: `./scripts/findings/functions/276.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/CourtScrapper/scraper_pool.py:51:5`
- Checklist pattern: interpolated values are an internally generated thread name and a static config constant.

Source excerpt:

```
    logger.info(f"{thread_name} starting for attorney: {first_name} {last_name}")
```

Why this is a false positive: `thread_name` is generated by the thread pool and `first_name`/`last_name` come from the developer-owned `ATTORNEYS` list in `config.py`; no externally controlled, CRLF-capable value is formatted into the message.

Checklist evidence: all interpolated segments trace to internal configuration or pool-generated strings, not to external input.

### [ ] Finding `280` — `CWE-396`

- Function context: `./scripts/findings/functions/280.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/CourtScrapper/scraper_pool.py:74:1`
- Checklist pattern: handler surfaces the failure — logs with traceback and returns the exception to the caller.

Source excerpt:

```
    except Exception as e:
        logger.error(f"{thread_name} error: {e}", exc_info=True)
        return (attorney_index, [], False, e)
```

Why this is a false positive: the rule condition is a generic handler that "can hide distinct failure conditions"; this handler logs the exception with `exc_info=True` and returns the exception object to the caller, which records and reports it — nothing is hidden.

Checklist evidence: the exception value is propagated in the return tuple (`(…, False, e)`) and the caller aggregates it into `thread_exceptions` for reporting.

### [ ] Finding `285` — `PERF-PY-28`

- Function context: `./scripts/findings/functions/285.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/CourtScrapper/scraper_pool.py:112:24`
- Checklist pattern: the executor is created once for the entire batch, not per unit of work.

Source excerpt:

```
def run_all_attorneys_concurrent(attorneys):
    ...
    with ThreadPoolExecutor(max_workers=num_workers, thread_name_prefix="ScraperWorker") as executor:
        # Submit all attorney scraping tasks to the thread pool
        future_to_attorney = {executor.submit(...): ...}
        ...
        logger.info(f"Submitted {len(future_to_attorney)} task(s) to thread pool")
```

Why this is a false positive: the rule condition is an executor "created per unit of work"; here the pool is constructed once at the top of the single batch function, all attorney tasks are submitted to it at once, and it lives for the whole run (the function is invoked exactly once from `main.py`). This already matches the recommended process-lifetime-pool pattern.

Checklist evidence: the executor is created once per program run and handles the complete workload in one `with` block; no per-task or per-work-item pool creation exists.

### [ ] Finding `298` — `CWE-117`

- Function context: `./scripts/findings/functions/298.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/CourtScrapper/utils.py:90:13`
- Checklist pattern: interpolated values are numeric counters.

Source excerpt:

```
            logger.debug(f"Chrome debug endpoint ready after {elapsed:.2f}s ({attempt} attempts)")
            return
```

Why this is a false positive: the only interpolated values are `elapsed` (a float formatted to two decimals) and `attempt` (an int counter) — internal numerics that cannot carry CRLF, so no CRLF neutralization concern exists.

Checklist evidence: both dynamic segments are numeric (float/int) values; no string input is formatted into the log message.

### [ ] Finding `302` — `CWE-396`

- Function context: `./scripts/findings/functions/302.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/CourtScrapper/utils.py:136:1`
- Checklist pattern: handler does not hide failures — it logs and re-raises as a chained error.

Source excerpt:

```
            except Exception as e:
                if temp_socket:
                    temp_socket.close()
                logger.error(f"Failed to get OS-assigned port: {e}")
                raise RuntimeError(f"Unable to obtain a free port for Chrome debugging: {e}") from e
```

Why this is a false positive: the rule condition is a generic handler that "can hide distinct failure conditions"; here the handler logs the failure and re-raises it as a chained `RuntimeError ... from e`, so the original failure is preserved and surfaced.

Checklist evidence: the handler ends with `raise RuntimeError(...) from e`, explicitly chaining and propagating the original exception rather than swallowing it.

### [ ] Finding `305` — `CWE-88`

- Function context: `./scripts/findings/functions/305.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/CourtScrapper/utils.py:148:24`
- Checklist pattern: interpolated argv segments are OS-generated or config values, not attacker-controlled.

Source excerpt:

```
            temp_socket.bind(('', 0))
            debug_port = temp_socket.getsockname()[1]
            ...
            temp_dir = tempfile.mkdtemp(prefix="playwright-debug-")
            ...
            proc = subprocess.Popen([
                chrome_path,
                f"--remote-debugging-port={debug_port}",
                f"--user-data-dir={temp_dir}"
            ], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
```

Why this is a false positive: the rule condition is an "untrusted argument [that] can become a tool option"; the dynamic segments are an OS-assigned integer port (`getsockname()[1]`), a `mkdtemp`-generated directory path, and a config constant — none attacker-controlled, and the port is numeric so it cannot introduce option delimiters.

Checklist evidence: all interpolated argv values originate from the OS (`socket` bind, `tempfile.mkdtemp`) or local config; no externally supplied text reaches the argument vector.

## True positives

### BP-PY-1 — Bare Except Clause (74 findings)

All findings show a bare `except:` or broad `except Exception` that swallows failures without re-raising; per the rule condition ("catch specific types or re-raise"), each is a genuine match.

| Finding | Source | Reason |
| --- | --- | --- |
| 2 | captcha_handler.py:44 | bare `except:` + `continue` |
| 3 | captcha_handler.py:55 | bare `except:` + `pass` |
| 7 | captcha_handler.py:76 | `except Exception:` + `pass` |
| 12 | captcha_handler.py:104 | bare `except:` + `pass` |
| 14 | captcha_handler.py:109 | `except Exception as e:` + log + return None |
| 32 | captcha_handler.py:264 | `except Exception as e:` + log + return False |
| 34 | captcha_handler.py:303 | `except Exception as e:` + log |
| 37 | captcha_handler.py:355 | `except Exception as e:` + log |
| 39 | captcha_handler.py:377 | `except Exception as e:` + log |
| 43 | captcha_handler.py:404 | `except Exception as e:` + log + return False |
| 47 | captcha_handler.py:434 | `except Exception as e:` + log + return None |
| 51 | captcha_handler.py:502 | `except Exception as e:` + log + return False |
| 55 | captcha_handler.py:552 | `except Exception as e:` + log + fallback |
| 59 | case_data_extractor.py:113 | `except Exception as exc:` + log + return "" |
| 63 | case_data_extractor.py:130 | `except Exception:` + `pass` |
| 67 | case_data_extractor.py:138 | `except Exception:` + `pass` |
| 69 | case_data_extractor.py:170 | `except Exception as exc:` + log + continue |
| 72 | case_data_extractor.py:273 | `except Exception as exc:` + log |
| 75 | case_data_extractor.py:339 | `except Exception as e:` + log |
| 80 | case_data_extractor.py:359 | `except Exception as exc:` + log |
| 82 | case_data_extractor.py:374 | `except Exception:` + return "" |
| 83 | case_data_extractor.py:383 | `except Exception as exc:` + log |
| 85 | case_data_extractor.py:429 | `except Exception as exc:` + log |
| 87 | case_data_extractor.py:475 | `except Exception as exc:` + log + return "" |
| 90 | case_data_extractor.py:586 | `except Exception as exc:` + log |
| 95 | case_data_extractor.py:657 | `except Exception as exc:` + log + return default |
| 99 | case_data_extractor.py:699 | `except Exception as e:` + log + return False |
| 121 | inspect_website.py:115 | `except Exception as e:` + log + next selector |
| 124 | inspect_website.py:122 | `except Exception as e:` + log |
| 128 | main.py:42 | `except Exception as e:` + log (top-level catch-all) |
| 154 | scraper.py:109 | `except Exception as e:` + log + return False |
| 159 | scraper.py:132 | `except Exception as e:` + log + return False |
| 162 | scraper.py:161 | `except Exception as e:` + log + continue |
| 164 | scraper.py:166 | `except Exception as e:` + log + return False |
| 166 | scraper.py:199 | `except Exception as e:` + log + continue |
| 168 | scraper.py:215 | `except Exception as e:` + log + return False |
| 171 | scraper.py:247 | `except Exception as e:` + log + return False |
| 174 | scraper.py:272 | `except Exception as e:` + log + continue |
| 178 | scraper.py:310 | `except Exception as e:` + log + continue |
| 187 | scraper.py:354 | `except Exception as e:` + log + return True |
| 190 | scraper.py:373 | `except Exception:` + log, fallback path |
| 191 | scraper.py:410 | bare `except:` + nested fallback |
| 192 | scraper.py:418 | bare `except:` + nested fallback |
| 195 | scraper.py:431 | `except Exception as e:` + log + continue |
| 197 | scraper.py:437 | `except Exception as e:` + log + return False |
| 199 | scraper.py:476 | bare `except:` + `continue` |
| 203 | scraper.py:490 | `except Exception as e:` + log + return [] |
| 205 | scraper.py:506 | `except Exception:` + `pass` |
| 210 | scraper.py:536 | `except Exception as e:` + log + continue |
| 212 | scraper.py:542 | `except Exception as e:` + log + return False |
| 216 | scraper.py:571 | `except Exception as e:` + log + return None |
| 220 | scraper.py:598 | `except Exception as count_exc:` + log + count=0 |
| 228 | scraper.py:627 | `except Exception as e:` + log + continue |
| 231 | scraper.py:637 | `except Exception as e:` + log + return False |
| 235 | scraper.py:659 | `except Exception as exc:` + log + return False |
| 242 | scraper.py:723 | `except Exception as e:` + log + recovery |
| 244 | scraper.py:741 | `except Exception as e:` + log |
| 247 | scraper.py:763 | `except Exception as e:` + log + continue |
| 256 | scraper.py:810 | `except Exception as e:` + log + continue |
| 262 | scraper.py:872 | `except Exception as e:` + log + return False |
| 265 | scraper.py:882 | bare `except:` + `pass` |
| 267 | scraper.py:888 | bare `except:` + `pass` |
| 269 | scraper.py:894 | bare `except:` + `pass` |
| 271 | scraper.py:900 | bare `except:` + `pass` |
| 273 | scraper.py:905 | `except Exception as e:` + log |
| 279 | scraper_pool.py:74 | `except Exception as e:` + log + error tuple |
| 289 | scraper_pool.py:141 | `except Exception as e:` + log + append error |
| 295 | scraper_pool.py:164 | `except Exception as e:` + log + return [] |
| 310 | utils.py:173 | `except Exception:` + `pass` |
| 317 | utils.py:192 | `except Exception as e:` + log + fallback |
| 319 | utils.py:210 | `except Exception as e:` + log + kill fallback |
| 321 | utils.py:216 | `except Exception as kill_error:` + log |
| 323 | utils.py:224 | `except Exception as e:` + log |
| 326 | utils.py:233 | `except Exception as e:` + log |

### BP-PY-2 — Except Pass (11 findings)

All show an exception handler whose body is only `pass`; per the rule condition ("failures are discarded silently"), each is genuine.

| Finding | Source | Reason |
| --- | --- | --- |
| 4 | captcha_handler.py:55 | `except:` + `pass` |
| 8 | captcha_handler.py:76 | `except Exception:` + `pass` |
| 13 | captcha_handler.py:104 | `except:` + `pass` |
| 64 | case_data_extractor.py:130 | `except Exception:` + `pass` |
| 68 | case_data_extractor.py:138 | `except Exception:` + `pass` |
| 206 | scraper.py:506 | `except Exception:` + `pass` |
| 266 | scraper.py:882 | `except:` + `pass` |
| 268 | scraper.py:888 | `except:` + `pass` |
| 270 | scraper.py:894 | `except:` + `pass` |
| 272 | scraper.py:900 | `except:` + `pass` |
| 311 | utils.py:173 | `except Exception:` + `pass` |

### BP-PY-46 — print Debugging In Library Code (26 findings)

Neither `result_exporter.py` nor `utils.py` has an `if __name__ == "__main__"` guard or argparse CLI; all flagged `print(` calls sit in module-level functions, matching the rule condition.

| Finding | Source | Reason |
| --- | --- | --- |
| 131 | result_exporter.py:22 | `print("No results to export")` in module function |
| 132 | result_exporter.py:87 | `print(f"Unknown output format: …")` |
| 133 | result_exporter.py:98 | `print(f"Results exported to: …")` in `export_csv` |
| 134 | result_exporter.py:105 | `print(f"Results exported to: …")` in `export_json` |
| 135 | result_exporter.py:200 | `print(f"Results exported to: …")` in `export_excel` |
| 136 | result_exporter.py:204 | `print(f"Created {…} sheets …")` |
| 141 | result_exporter.py:275 | `print(f"\nTotal cases found: …")` in `print_export_summary` |
| 142 | result_exporter.py:279 | `print("\nCases by Attorney:")` |
| 143 | result_exporter.py:282 | `print(f"  - {attorney}: …")` |
| 144 | result_exporter.py:284 | `print("\nColumn titles:")` |
| 145 | result_exporter.py:286 | `print(f"  - {col}")` |
| 146 | result_exporter.py:287 | `print("\nFirst few results:")` |
| 147 | result_exporter.py:288 | `print(df.head().to_string())` |
| 328 | utils.py:353 | `print(f"\nAttorneys to search …")` in `display_config` |
| 329 | utils.py:355 | `print(f"  {idx}. {attorney…}")` |
| 330 | utils.py:358 | `print(f"\nCase type filter:")` |
| 331 | utils.py:360 | `print("  (None - will get all case types)")` |
| 332 | utils.py:363 | `print("  (Empty list - will get all case types)")` |
| 333 | utils.py:365 | `print(f"  {len(CASE_TYPE)} case type(s):")` |
| 334 | utils.py:367 | `print(f"    {idx}. {case_type}")` |
| 335 | utils.py:369 | `print(f"  {CASE_TYPE}")` |
| 336 | utils.py:372 | `print(f"\nCharge keywords to filter:")` |
| 337 | utils.py:374 | `print("  (None - will get all cases regardless of charge)")` |
| 338 | utils.py:376 | `print("  (Empty list - will get all cases regardless of charge)")` |
| 339 | utils.py:378 | `print(f"  {len(CHARGE_KEYWORDS)} keyword(s):")` |
| 340 | utils.py:380 | `print(f"    {idx}. {keyword}")` |

### BP-PY-47 — logging With String Format Before Logger (192 findings)

Every finding is a `logger.debug/info/warning/error(f"…")` call whose first argument is an eagerly formatted f-string, exactly the rule condition ("eager f-string or .format( as first arg").

| Finding | Source | Reason |
| --- | --- | --- |
| 1 | captcha_handler.py:42 | f-string log arg |
| 10 | captcha_handler.py:93 | f-string log arg |
| 11 | captcha_handler.py:102 | f-string log arg |
| 15 | captcha_handler.py:110 | f-string log arg |
| 17 | captcha_handler.py:129 | f-string log arg |
| 18 | captcha_handler.py:133 | f-string log arg |
| 19 | captcha_handler.py:134 | f-string log arg |
| 20 | captcha_handler.py:157 | f-string log arg |
| 21 | captcha_handler.py:159 | f-string log arg |
| 22 | captcha_handler.py:182 | f-string log arg |
| 23 | captcha_handler.py:193 | f-string log arg |
| 24 | captcha_handler.py:217 | f-string log arg |
| 25 | captcha_handler.py:226 | f-string log arg |
| 27 | captcha_handler.py:240 | f-string log arg |
| 28 | captcha_handler.py:251 | f-string log arg |
| 29 | captcha_handler.py:254 | f-string log arg |
| 30 | captcha_handler.py:258 | f-string log arg |
| 31 | captcha_handler.py:262 | f-string log arg |
| 33 | captcha_handler.py:265 | f-string log arg |
| 35 | captcha_handler.py:304 | f-string log arg |
| 36 | captcha_handler.py:354 | f-string log arg |
| 38 | captcha_handler.py:356 | f-string log arg |
| 40 | captcha_handler.py:378 | f-string log arg |
| 42 | captcha_handler.py:394 | f-string log arg |
| 44 | captcha_handler.py:405 | f-string log arg |
| 45 | captcha_handler.py:427 | f-string log arg |
| 46 | captcha_handler.py:431 | f-string log arg |
| 48 | captcha_handler.py:435 | f-string log arg |
| 49 | captcha_handler.py:486 | f-string log arg |
| 50 | captcha_handler.py:494 | f-string log arg |
| 52 | captcha_handler.py:503 | f-string log arg |
| 53 | captcha_handler.py:526 | f-string log arg |
| 54 | captcha_handler.py:548 | f-string log arg |
| 56 | captcha_handler.py:553 | f-string log arg |
| 57 | case_data_extractor.py:48 | f-string log arg |
| 58 | case_data_extractor.py:52 | f-string log arg |
| 61 | case_data_extractor.py:114 | f-string log arg |
| 70 | case_data_extractor.py:171 | f-string log arg |
| 73 | case_data_extractor.py:274 | f-string log arg |
| 74 | case_data_extractor.py:318 | f-string log arg |
| 76 | case_data_extractor.py:340 | f-string log arg |
| 77 | case_data_extractor.py:342 | f-string log arg |
| 78 | case_data_extractor.py:344 | f-string log arg |
| 81 | case_data_extractor.py:360 | f-string log arg |
| 84 | case_data_extractor.py:384 | f-string log arg |
| 86 | case_data_extractor.py:430 | f-string log arg |
| 88 | case_data_extractor.py:476 | f-string log arg |
| 89 | case_data_extractor.py:584 | f-string log arg |
| 91 | case_data_extractor.py:587 | f-string log arg |
| 92 | case_data_extractor.py:625 | f-string log arg |
| 93 | case_data_extractor.py:629 | f-string log arg |
| 94 | case_data_extractor.py:651 | f-string log arg |
| 96 | case_data_extractor.py:658 | f-string log arg |
| 97 | case_data_extractor.py:694 | f-string log arg |
| 98 | case_data_extractor.py:697 | f-string log arg |
| 100 | case_data_extractor.py:700 | f-string log arg |
| 101 | inspect_website.py:14 | f-string log arg |
| 103 | inspect_website.py:15 | f-string log arg |
| 104 | inspect_website.py:16 | f-string log arg |
| 105 | inspect_website.py:19 | f-string log arg |
| 106 | inspect_website.py:20 | f-string log arg |
| 107 | inspect_website.py:30 | f-string log arg |
| 108 | inspect_website.py:38 | f-string log arg |
| 109 | inspect_website.py:40 | f-string log arg |
| 110 | inspect_website.py:43 | f-string log arg |
| 111 | inspect_website.py:52 | f-string log arg |
| 112 | inspect_website.py:60 | f-string log arg |
| 113 | inspect_website.py:65 | f-string log arg |
| 114 | inspect_website.py:68 | f-string log arg |
| 115 | inspect_website.py:71 | f-string log arg |
| 116 | inspect_website.py:74 | f-string log arg |
| 117 | inspect_website.py:105 | f-string log arg |
| 118 | inspect_website.py:108 | f-string log arg |
| 119 | inspect_website.py:109 | f-string log arg |
| 120 | inspect_website.py:110 | f-string log arg |
| 123 | inspect_website.py:116 | f-string log arg |
| 125 | inspect_website.py:123 | f-string log arg |
| 126 | main.py:30 | f-string log arg |
| 130 | main.py:43 | f-string log arg |
| 139 | result_exporter.py:260 | f-string log arg |
| 148 | scraper.py:66 | f-string log arg |
| 150 | scraper.py:68 | f-string log arg |
| 151 | scraper.py:82 | f-string log arg |
| 153 | scraper.py:107 | f-string log arg |
| 155 | scraper.py:110 | f-string log arg |
| 156 | scraper.py:116 | f-string log arg |
| 157 | scraper.py:127 | f-string log arg |
| 158 | scraper.py:130 | f-string log arg |
| 160 | scraper.py:133 | f-string log arg |
| 161 | scraper.py:157 | f-string log arg |
| 163 | scraper.py:162 | f-string log arg |
| 165 | scraper.py:167 | f-string log arg |
| 167 | scraper.py:200 | f-string log arg |
| 169 | scraper.py:216 | f-string log arg |
| 170 | scraper.py:225 | f-string log arg |
| 172 | scraper.py:248 | f-string log arg |
| 173 | scraper.py:268 | f-string log arg |
| 175 | scraper.py:273 | f-string log arg |
| 176 | scraper.py:307 | f-string log arg |
| 177 | scraper.py:308 | f-string log arg |
| 179 | scraper.py:311 | f-string log arg |
| 180 | scraper.py:325 | f-string log arg |
| 181 | scraper.py:328 | f-string log arg |
| 182 | scraper.py:331 | f-string log arg |
| 183 | scraper.py:340 | f-string log arg |
| 184 | scraper.py:343 | f-string log arg |
| 185 | scraper.py:346 | f-string log arg |
| 186 | scraper.py:351 | f-string log arg |
| 188 | scraper.py:355 | f-string log arg |
| 194 | scraper.py:425 | f-string log arg |
| 196 | scraper.py:432 | f-string log arg |
| 198 | scraper.py:438 | f-string log arg |
| 200 | scraper.py:481 | f-string log arg |
| 201 | scraper.py:486 | f-string log arg |
| 202 | scraper.py:488 | f-string log arg |
| 204 | scraper.py:491 | f-string log arg |
| 209 | scraper.py:531 | f-string log arg |
| 211 | scraper.py:537 | f-string log arg |
| 213 | scraper.py:543 | f-string log arg |
| 214 | scraper.py:552 | f-string log arg |
| 215 | scraper.py:564 | f-string log arg (multi-line) |
| 217 | scraper.py:572 | f-string log arg |
| 218 | scraper.py:587 | f-string log arg |
| 219 | scraper.py:591 | f-string log arg (multi-line) |
| 221 | scraper.py:599 | f-string log arg |
| 222 | scraper.py:603 | f-string log arg |
| 223 | scraper.py:607 | f-string log arg |
| 224 | scraper.py:610 | f-string log arg |
| 225 | scraper.py:613 | f-string log arg (multi-line) |
| 226 | scraper.py:623 | f-string log arg |
| 227 | scraper.py:625 | f-string log arg |
| 229 | scraper.py:628 | f-string log arg |
| 230 | scraper.py:635 | f-string log arg |
| 232 | scraper.py:638 | f-string log arg |
| 233 | scraper.py:646 | f-string log arg (multi-line) |
| 234 | scraper.py:657 | f-string log arg |
| 236 | scraper.py:660 | f-string log arg |
| 237 | scraper.py:677 | f-string log arg |
| 238 | scraper.py:681 | f-string log arg |
| 239 | scraper.py:685 | f-string log arg |
| 240 | scraper.py:703 | f-string log arg |
| 241 | scraper.py:715 | f-string log arg |
| 243 | scraper.py:724 | f-string log arg |
| 245 | scraper.py:742 | f-string log arg |
| 246 | scraper.py:753 | f-string log arg |
| 248 | scraper.py:764 | f-string log arg |
| 249 | scraper.py:769 | f-string log arg |
| 250 | scraper.py:777 | f-string log arg |
| 251 | scraper.py:782 | f-string log arg |
| 252 | scraper.py:792 | f-string log arg |
| 253 | scraper.py:797 | f-string log arg |
| 254 | scraper.py:806 | f-string log arg |
| 255 | scraper.py:807 | f-string log arg |
| 257 | scraper.py:811 | f-string log arg |
| 258 | scraper.py:814 | f-string log arg |
| 259 | scraper.py:854 | f-string log arg |
| 260 | scraper.py:862 | f-string log arg |
| 261 | scraper.py:868 | f-string log arg |
| 263 | scraper.py:873 | f-string log arg |
| 274 | scraper.py:906 | f-string log arg |
| 275 | scraper_pool.py:51 | f-string log arg |
| 277 | scraper_pool.py:65 | f-string log arg (multi-line) |
| 278 | scraper_pool.py:71 | f-string log arg |
| 281 | scraper_pool.py:75 | f-string log arg |
| 282 | scraper_pool.py:83 | f-string log arg |
| 284 | scraper_pool.py:107 | f-string log arg (multi-line) |
| 286 | scraper_pool.py:119 | f-string log arg |
| 287 | scraper_pool.py:121 | f-string log arg (multi-line) |
| 290 | scraper_pool.py:142 | f-string log arg |
| 291 | scraper_pool.py:149 | f-string log arg |
| 292 | scraper_pool.py:151 | f-string log arg |
| 293 | scraper_pool.py:159 | f-string log arg |
| 294 | scraper_pool.py:161 | f-string log arg |
| 296 | scraper_pool.py:165 | f-string log arg |
| 297 | utils.py:90 | f-string log arg |
| 299 | utils.py:102 | f-string log arg |
| 301 | utils.py:135 | f-string log arg |
| 303 | utils.py:139 | f-string log arg |
| 304 | utils.py:144 | f-string log arg |
| 306 | utils.py:154 | f-string log arg |
| 307 | utils.py:164 | f-string log arg |
| 308 | utils.py:166 | f-string log arg |
| 309 | utils.py:172 | f-string log arg |
| 314 | utils.py:183 | f-string log arg |
| 315 | utils.py:187 | f-string log arg |
| 316 | utils.py:190 | f-string log arg |
| 318 | utils.py:193 | f-string log arg |
| 320 | utils.py:211 | f-string log arg |
| 322 | utils.py:217 | f-string log arg |
| 324 | utils.py:225 | f-string log arg |
| 325 | utils.py:232 | f-string log arg |
| 327 | utils.py:234 | f-string log arg |

### CWE-1071 — Empty Code Block (4 findings)

| Finding | Source | Reason |
| --- | --- | --- |
| 6 | captcha_handler.py:55 | `except:` + `pass` only |
| 66 | case_data_extractor.py:130 | `except Exception:` + `pass` only |
| 208 | scraper.py:506 | `except Exception:` + `pass` only |
| 313 | utils.py:173 | `except Exception:` + `pass` only |

### CWE-1121 — Excessive McCabe Cyclomatic Complexity (6 findings)

Branch counts verified against the rule's counting (≥12 of `if/elif/for/while/except` tokens in the function body):

| Finding | Source | Reason |
| --- | --- | --- |
| 16 | captcha_handler.py:114 (`solve_recaptcha_v2_with_2captcha`) | 24 branch tokens |
| 79 | case_data_extractor.py:351 (`extract_charge_description`) | 15 branch tokens |
| 137 | result_exporter.py:211 (`format_excel_sheets`) | 13 branch tokens |
| 189 | scraper.py:359 (`set_items_per_page`) | 13 branch tokens |
| 283 | scraper_pool.py:92 (`run_all_attorneys_concurrent`) | 14 branch tokens |
| 300 | utils.py:108 (`get_chrome_user_agent`) | 23 branch tokens |

### CWE-1124 — Excessively Deep Nesting (4 findings)

Nesting depth verified at the flagged statements (6 control-flow levels):

| Finding | Source | Reason |
| --- | --- | --- |
| 71 | case_data_extractor.py:272 | statement at depth 6 |
| 138 | result_exporter.py:237 | statement at depth 6 |
| 193 | scraper.py:423 | statement at depth 6 |
| 288 | scraper_pool.py:135 | statement at depth 6 |

### CWE-117 — Improper Output Neutralization for Logs (3 findings)

| Finding | Source | Reason |
| --- | --- | --- |
| 41 | captcha_handler.py:378 | exception `e` from `page.evaluate` (page-derived value) interpolated un-neutralized |
| 62 | case_data_extractor.py:114 | DOM-derived `heading_text` interpolated un-neutralized |
| 140 | result_exporter.py:260 | `col` (scraped-data column name) interpolated un-neutralized |

### CWE-390 — Detection of Error Condition Without Action (4 findings)

| Finding | Source | Reason |
| --- | --- | --- |
| 5 | captcha_handler.py:55 | exception detected, handler only `pass` |
| 65 | case_data_extractor.py:130 | exception detected, handler only `pass` |
| 207 | scraper.py:506 | exception detected, handler only `pass` |
| 312 | utils.py:173 | exception detected, handler only `pass` |

### CWE-396 — Declaration of Catch for Generic Exception (4 findings)

| Finding | Source | Reason |
| --- | --- | --- |
| 9 | captcha_handler.py:76 | `except Exception:` + `pass` hides failures |
| 60 | case_data_extractor.py:113 | `except Exception as exc:` collapses all failures to "" |
| 122 | inspect_website.py:115 | `except Exception as e:` collapses all failures to "not found" |
| 129 | main.py:42 | top-level `except Exception as e:` catch-all |

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `./scripts/chunks`
- Function evidence: `./scripts/findings/functions`
- Validation: `git diff --check` — `pass` (see below)

```
$ git diff --check   # run in goslop repo root after writing this report
<no output — pass>
```
