# False-positive audit — among-llms

## Run metadata

```yaml
timestamp: 2026-08-02T07:22:52Z
repository: among-llms
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/among-llms
branch: develop
commit: a6d4f57
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/among-llms
chunk_path: scripts/among-llms/chunks
function_context_path: scripts/among-llms/findings/functions
```

## Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/among-llms/chunks -context-dir scripts/among-llms/findings/functions real-repos/among-llms`
- Findings: `72`
- Chunks reviewed: `scripts/among-llms/chunks/Chunk_1_25.txt`, `Chunk_26_50.txt`, `Chunk_51_72.txt`
- Function contexts reviewed: `scripts/among-llms/findings/functions/1.txt` … `72.txt` (all 72 exist; sources re-read directly for every finding)

## Audit checklist

- [x] Read every assigned chunk under `scripts/among-llms/chunks`.
- [x] Read `scripts/among-llms/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 0 | — |
| True positive | 72 | 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72 |
| Uncertain | 0 | — |

## False positives

No false positives found. Every finding satisfies the literal rule condition of the reported rule, verified against the detectors in `internal/lang/python/detectors/bad_practices/rules_core.go`, `rules_core_extra.go`, `rules_observability.go` and `internal/lang/python/detectors/cwe/rules_injection.go`, `rules_platform.go`, `rules_secrets.go`, `rules_tier_b_runtime.go`. No built-in exemption (re-raise, test module, fixture assignment, `__init__.py` re-export) applies to any shown source.

## True positives

### BP-PY-47 — logging With String Format Before Logger (54 findings)

Rule condition (`rules_observability.go` `detectBPPY47`): text match of `logger.<method>(` / `logging.<method>(` (methods include `log`, `critical`, …) whose first argument is an f-string (any `f`/`fr`/`rf` prefix) or a `.format(`/`%`-preformatted expression (`isEagerLogFormat`). All 54 findings are `AppConfiguration.logger.log(f"…")`, `self._logger.log(f"…")` or `logging.critical(f"…")` calls with an f-string first argument. `AppConfiguration.logger` is an `AppLogger` instance (`allms/utils/logger.py`) whose `log(msg, level)` forwards the already-formatted message — the f-string is eagerly evaluated before the call.

| Finding | Source | Reason |
| --- | --- | --- |
| 1 | `allms/__main__.py:53` | `logging.critical(f"…{e}")` — `logging.critical(` needle, f-string first arg. |
| 3 | `allms/__main__.py:54` | `logging.critical(f"Exiting the app ...")` — f-string first arg; rule matches any `f`-prefixed literal. |
| 4 | `allms/__main__.py:60` | `AppConfiguration.logger.log(f"…", level=…)` — `logger.log(` needle, f-string first arg. |
| 5 | `allms/cli/screens/main.py:83` | `AppConfiguration.logger.log(f"…")` — f-string first arg. |
| 8 | `allms/cli/widgets/load.py:107` | `AppConfiguration.logger.log(f"…{err}'", …)` — f-string first arg. |
| 9 | `allms/cli/widgets/new.py:110` | `AppConfiguration.logger.log(f"…")` — f-string first arg. |
| 10 | `allms/cli/widgets/type.py:35` | `AppConfiguration.logger.log(f"…{self._are_typing}")` — f-string first arg. |
| 11 | `allms/cli/widgets/vote.py:126` | `AppConfiguration.logger.log(f"You …", f"…")` — f-string first arg. |
| 14 | `allms/core/agents.py:56` | `AppConfiguration.logger.log(f"…{msg}")` — f-string first arg. |
| 15 | `allms/core/chat/history.py:28` | `AppConfiguration.logger.log(f"…{message}")` — f-string first arg. |
| 16 | `allms/core/chat/history.py:33` | `AppConfiguration.logger.log(f"…{msg_id}…")` — f-string first arg. |
| 17 | `allms/core/chat/history.py:42` | `AppConfiguration.logger.log(f"…{msg_id}…")` — f-string first arg. |
| 21 | `allms/core/llm/loop.py:45` | `AppConfiguration.logger.log(f"…{agent_id}…")` — f-string first arg. |
| 22 | `allms/core/llm/loop.py:116` | `AppConfiguration.logger.log(f"…{agent_id}…")` — f-string first arg. |
| 23 | `allms/core/llm/loop.py:131` | `AppConfiguration.logger.log(f"…{model_response}")` — f-string first arg. |
| 24 | `allms/core/llm/loop.py:163` | `AppConfiguration.logger.log(f"…{agent_id}…")` — f-string first arg. |
| 25 | `allms/core/llm/loop.py:177` | `AppConfiguration.logger.log(f"…{agent_id}…")` — f-string first arg. |
| 26 | `allms/core/llm/manager.py:72` | `AppConfiguration.logger.log(f"[{tries}] …", level=…)` — f-string first arg. |
| 27 | `allms/core/llm/manager.py:90` | `AppConfiguration.logger.log(f"…" + f"…", level=…)` — first arg is an f-string. |
| 28 | `allms/core/llm/manager.py:103` | `AppConfiguration.logger.log(f"…")` — f-string first arg. |
| 32 | `allms/core/state/manager.py:94` | `logger.log(f"…{e} ")` — `logger.log(` needle, f-string first arg. |
| 33 | `allms/core/state/manager.py:117` | `self._logger.log(f"…{your_id}")` — f-string first arg. |
| 34 | `allms/core/state/manager.py:140` | `self._logger.log(f"…", level=…)` — f-string first arg. |
| 35 | `allms/core/state/manager.py:168` | `self._logger.log(f"…{genre}' …")` — f-string first arg. |
| 36 | `allms/core/state/manager.py:241` | `self._logger.log(f"…'{genre}' …")` — f-string first arg. |
| 37 | `allms/core/state/manager.py:249` | `self._logger.log(f"…'{scenario}' …")` — f-string first arg. |
| 38 | `allms/core/state/manager.py:313` | `self._logger.log(f"…{started_by}…")` — f-string first arg. |
| 39 | `allms/core/state/manager.py:336` | `self._logger.log(f"…{end_ts_iso}")` — f-string first arg. |
| 40 | `allms/core/state/manager.py:349` | `self._logger.log(f"Result: {vote_conclusion}")` — f-string first arg. |
| 41 | `allms/core/state/manager.py:375` | `self._logger.log(f"…{by_agent}…{for_agent}…")` — f-string first arg. |
| 42 | `allms/core/state/manager.py:399` | `self._logger.log(f"{agent_id} terminated", level=…)` — f-string first arg. |
| 43 | `allms/core/state/manager.py:412` | `self._logger.log(f"…{inform_msg}")` — f-string first arg. |
| 44 | `allms/core/state/manager.py:446` | `self._logger.log(f"…")` — f-string first arg. |
| 45 | `allms/core/state/manager.py:458` | `self._logger.log(f"…")` — f-string first arg. |
| 46 | `allms/core/state/manager.py:462` | `self._logger.log(f"…", level=…)` — f-string first arg. |
| 47 | `allms/core/state/manager.py:466` | `self._logger.log(f"Stopping background worker")` — f-string first arg. |
| 48 | `allms/core/state/manager.py:467` | `self._logger.log(f"Elapsed …: {duration} {duration_unit}")` — f-string first arg. |
| 49 | `allms/core/state/manager.py:498` | `AppConfiguration.logger.log(f"…{chat_msg}")` — f-string first arg. |
| 50 | `allms/core/state/state.py:193` | `AppConfiguration.logger.log(\n f"…{agent_id}…" + …)` — first arg is an f-string. |
| 51 | `allms/core/state/state.py:271` | `AppConfiguration.logger.log(f"…{agent_ids}…")` — f-string first arg. |
| 52 | `allms/core/state/state.py:296` | `AppConfiguration.logger.log(f"…", level=…)` — f-string first arg. |
| 53 | `allms/core/state/state.py:319` | `AppConfiguration.logger.log(f"…{started_by}…")` — f-string first arg. |
| 54 | `allms/core/state/state.py:333` | `AppConfiguration.logger.log(f"{by_agent} …", level=…)` — f-string first arg. |
| 55 | `allms/core/state/state.py:336` | `AppConfiguration.logger.log(f"…{for_agent} …", level=…)` — f-string first arg. |
| 56 | `allms/core/state/state.py:357` | `AppConfiguration.logger.log(f"…")` — f-string first arg. |
| 57 | `allms/core/state/state.py:424` | `AppConfiguration.logger.log(f"…{notify_msg}", level=…)` — f-string first arg. |
| 58 | `allms/core/vote.py:27` | `AppConfiguration.logger.log(f"…{started_by}…")` — f-string first arg. |
| 59 | `allms/core/vote.py:33` | `AppConfiguration.logger.log(f"…" + f"…", level=…)` — first arg is an f-string. |
| 60 | `allms/core/vote.py:39` | `AppConfiguration.logger.log(f"…{by_agent}…{for_agent}")` — f-string first arg. |
| 61 | `allms/core/vote.py:77` | `AppConfiguration.logger.log(f"-- Voting process has ended --")` — f-string first arg. |
| 62 | `allms/core/vote.py:88` | `AppConfiguration.logger.log(f"…{by_agent}…")` — f-string first arg. |
| 63 | `allms/utils/callbacks.py:26` | `AppConfiguration.logger.log(f"…{callback.__name__}…" + …, …)` — first arg is an f-string. |
| 66 | `allms/utils/save.py:33` | `AppConfiguration.logger.log(f"…{type(obj)}…", level=…)` — f-string first arg. |
| 67 | `allms/utils/save.py:120` | `AppConfiguration.logger.log(f"…{f.name}: {field_value}")` — f-string first arg. |

### BP-PY-1 — Bare Except Clause

Rule condition (`rules_core.go` `detectBPPY1`): `except Exception` / `except BaseException` (or bare `except:`) whose suite does not clearly re-raise (`suiteReraises`), outside test modules with evidence-collection suites. The three findings all match the broad-except branch and none re-raises.

| Finding | Source | Reason |
| --- | --- | --- |
| 6 | `allms/cli/widgets/load.py:106` | `except Exception as err:` with `logger.log(...)` + `self.notify(...)` suite; no re-raise, non-test file. |
| 30 | `allms/core/state/manager.py:93` | `except Exception as e:` with `logger.log(...); return None`; no re-raise, non-test file. |
| 68 | `allms/utils/save.py:134` | `except Exception:` with `pass` suite; no re-raise, non-test file. |

### BP-PY-2 — Except Pass

Rule condition (`rules_core.go` `detectBPPY2`): an except clause whose direct suite consists solely of `pass`. The rule has no test-file exemption.

| Finding | Source | Reason |
| --- | --- | --- |
| 69 | `allms/utils/save.py:134` | `except Exception:` suite is only `pass` (`suite[0] == "pass"`). |

### BP-PY-5 — Wildcard Import

Rule condition (`rules_core_extra.go` `detectBPPY5`): `from <mod> import *` in a non-`__init__.py` module (`fromImportStarRe`). Neither file is `__init__.py`.

| Finding | Source | Reason |
| --- | --- | --- |
| 12 | `allms/config/app.py:6` | `from allms.config.models import *` — matches `fromImportStarRe`, file is `app.py` not `__init__.py`. |
| 20 | `allms/core/llm/factory.py:4` | `from allms.core.llm.client import *` — matches `fromImportStarRe`, file is `factory.py` not `__init__.py`. |

### BP-PY-6 — assert Used For Runtime Validation

Rule condition (`rules_core.go` `detectBPPY6`): `assert <expr>` line where the expression contains runtime-validation needles (here `\bpath\b`), in a non-test module.

| Finding | Source | Reason |
| --- | --- | --- |
| 64 | `allms/utils/logger.py:52` | `assert log_dir.is_dir(), f"Provided path must be a valid directory"` — contains `path` token; non-test file. |
| 65 | `allms/utils/parser.py:14` | `assert self._file_path.exists(), f"…Make sure the path is correct"` — contains `path` token; non-test file. |

### CWE-117 — Improper Output Neutralization for Logs

Rule condition (`rules_injection.go` `detectCWE117`): `logging.debug/info/warning/error/critical` (or `logger.*`/`log.*`) call whose first argument `looksLogFormatted` (f-string / `.format(` / `%`) and is not CRLF-sanitized (`headerValueLooksSanitized`).

| Finding | Source | Reason |
| --- | --- | --- |
| 2 | `allms/__main__.py:53` | `logging.critical(f"…invalid configuration: {e}")` — f-string first arg with dynamic `{e}`; no `replace("\r"/"\n")` neutralization present. |

### CWE-312 — Cleartext Storage of Sensitive Information

Rule condition (`rules_secrets.go` `detectCWE312`): `pyCleartextSecretRE` — `api[_-]?key|secret…|…access_token|…auth_token|private_key = "<literal ≥3 chars>"` in non-test, non-benchmark source. No fixture exclusion applies.

| Finding | Source | Reason |
| --- | --- | --- |
| 18 | `allms/core/llm/client/ollama.py:20` | `api_key="ollama"` matches `pyCleartextSecretRE`; `ollama.py` is not a test/benchmark/fixtures file. |

### CWE-798 — Use of Hard-coded Credentials

Rule condition (`rules_secrets.go` `detectCWE798`): `pyHardcodedCredentialRE` — credential-shaped name (`password`/`api_key`/`secret`/…) `= "<literal ≥3 chars>"`; excluded only when `isFixtureCredentialAssignment` (fixture-password literal, `fixtures.py`, or fixture-named function) — none apply here.

| Finding | Source | Reason |
| --- | --- | --- |
| 19 | `allms/core/llm/client/ollama.py:20` | `api_key="ollama"` matches `pyHardcodedCredentialRE` (6 chars ≥ 3); not a fixture assignment (no fixture literal, not `fixtures.py`, function `create_client` has no fixture name). |

### CWE-390 — Detection of Error Condition Without Action

Rule condition (`rules_platform.go` `detectCWE390`): an except clause whose direct body is `pass` (`exceptPassStart`).

| Finding | Source | Reason |
| --- | --- | --- |
| 70 | `allms/utils/save.py:134` | `except Exception:` followed only by `pass` — handler takes no action. |

### CWE-396 — Declaration of Catch for Generic Exception

Rule condition (`rules_platform.go` `detectCWE396`): first `except Exception` / `except BaseException` (`pyGenericExceptRE`) in a non-test module.

| Finding | Source | Reason |
| --- | --- | --- |
| 7 | `allms/cli/widgets/load.py:106` | `except Exception as err:` matches `pyGenericExceptRE`; `load.py` is not a test module. |
| 13 | `allms/config/runtime.py:159` | `except Exception as e:` matches `pyGenericExceptRE`; `runtime.py` is not a test module. |
| 29 | `allms/core/llm/parser.py:54` | `except Exception as ex:` matches `pyGenericExceptRE`; `parser.py` is not a test module. |
| 31 | `allms/core/state/manager.py:93` | `except Exception as e:` matches `pyGenericExceptRE`; `manager.py` is not a test module. |
| 71 | `allms/utils/save.py:134` | `except Exception:` matches `pyGenericExceptRE`; `save.py` is not a test module. |

### CWE-1071 — Empty Code Block

Rule condition (`rules_tier_b_runtime.go` `detectCWE1071`): `pyTierBEmptyExceptRE` — an exception handler containing only `pass`.

| Finding | Source | Reason |
| --- | --- | --- |
| 72 | `allms/utils/save.py:134` | `except Exception:` handler body is only `pass` — empty code block. |

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/among-llms/chunks/Chunk_1_25.txt`, `Chunk_26_50.txt`, `Chunk_51_72.txt`
- Function evidence: `scripts/among-llms/findings/functions/1.txt` … `72.txt`
- Validation: `git diff --check` — pass

## Post-fix over-suppression audit (2026-08-02)

### Run metadata

```yaml
timestamp: 2026-08-02T16:38:00Z
repository: among-llms
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/among-llms
branch: develop
commit: a6d4f57
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/among-llms
chunk_path: scripts/among-llms/chunks
function_context_path: scripts/among-llms/findings/functions
scanner_binary: bin/goslop rebuilt 2026-08-02 16:29 from b5b8fde (FP-reduction fix)
```

Fresh scan produced 70 findings (old audit: 72). Fresh finding IDs were mapped to audited TPs by `Source:` path (fresh IDs 1–70 = old IDs 1–12, 14–28, 30–72 in order). Two audited TP sources are absent from the fresh scan: `allms/config/runtime.py:159` and `allms/core/llm/parser.py:54`, both CWE-396. Both constructs were verified still present at those locations in the current source (same commit `a6d4f57`).

### Over-suppressed audit table

| Old finding ID | Rule | Source | One-line reason (old audit) | Current status |
| --- | --- | --- | --- | --- |
| 13 | CWE-396 | `allms/config/runtime.py:159` | `except Exception as e:` matches `pyGenericExceptRE`; `runtime.py` is not a test module. | suppressed-but-present |
| 29 | CWE-396 | `allms/core/llm/parser.py:54` | `except Exception as ex:` matches `pyGenericExceptRE`; `parser.py` is not a test module. | suppressed-but-present |

Over-suppressed TPs found: 2. Fixed/removed: 0.

### Suppressed-but-present TP — old finding 13 — CWE-396

- Source: `allms/config/runtime.py:159`
- Rule condition: `detectCWE396` (`rules_platform.go`) matches the first `except Exception` / `except BaseException` (`pyGenericExceptRE`) in a non-test module.

Source excerpt (`allms/config/runtime.py:156-160`):

```
        try:
            if not path.exists():
                path.mkdir(parents=True, exist_ok=True)
        except Exception as e:
            raise RuntimeError(f"There was an issue creating the given directory: {str(path)}: {e}")
```

Why this is over-suppressed: the generic `except Exception as e:` at line 159 is unchanged (same commit `a6d4f57`, non-test module, still matches `pyGenericExceptRE`), so the old audited TP construct remains. The FP-reduction fix (`b5b8fde`) added `suiteSurfacesFailureMasked` to `detectCWE396` (rules_platform.go:107), which skips handlers whose suite re-raises — this suite starts with `raise RuntimeError(...)` (`suiteLineSurfacesFailure` returns true on the `raise ` prefix) — so the previously-reported TP is now silently dropped.

### Suppressed-but-present TP — old finding 29 — CWE-396

- Source: `allms/core/llm/parser.py:54`
- Rule condition: `detectCWE396` (`rules_platform.go`) matches the first `except Exception` / `except BaseException` (`pyGenericExceptRE`) in a non-test module.

Source excerpt (`allms/core/llm/parser.py:51-55`):

```
            except KeyError as ke:
                supp_keys = " ".join(list(parser_map.keys()))
                raise ValueError(f"{ke}. Supported keys: {supp_keys}. DOES NOT MATCH THE OUTPUT SCHEMA")
            except Exception as ex:
                raise ValueError(f"{ex}. Doesn't match the requested output schema")
```

Why this is over-suppressed: the generic `except Exception as ex:` at line 54 is unchanged (same commit `a6d4f57`, non-test module, still matches `pyGenericExceptRE`), so the old audited TP construct remains. The `suiteSurfacesFailureMasked` exemption added in `b5b8fde` skips this handler because its suite starts with `raise ValueError(...)`, dropping the previously-reported TP.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/among-llms/chunks/Chunk_1_25.txt`, `Chunk_26_50.txt`, `Chunk_51_70.txt` (fresh, post-fix)
- Function evidence: `scripts/among-llms/findings/functions/1.txt` … `70.txt` (fresh)
- Validation: `git diff --check` — pass
