# False-positive audit report — pytogether

## Run metadata

```yaml
timestamp: 2026-08-02T07:39:53Z
repository: pytogether
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pytogether
branch: main
commit: 2004decbafddca7342699318a2d8e50ba788177a
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pytogether
chunk_path: scripts/pytogether/chunks
function_context_path: scripts/pytogether/findings/functions
```

## Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/pytogether/chunks -context-dir scripts/pytogether/findings/functions real-repos/pytogether`
- Findings: `71`
- Chunks reviewed: `scripts/pytogether/chunks/Chunk_1_25.txt`, `scripts/pytogether/chunks/Chunk_26_50.txt`, `scripts/pytogether/chunks/Chunk_51_71.txt`
- Function contexts reviewed: `scripts/pytogether/findings/functions/1.txt` … `scripts/pytogether/findings/functions/71.txt` (all 71)

## Audit checklist

- [x] Read every assigned chunk under `scripts/pytogether/chunks`.
- [x] Read `scripts/pytogether/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 1 | 12 |
| True positive | 70 | 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71 |
| Uncertain | 0 | — |

## False positives

### [x] Finding `12` — `BP-PY-46`

- Function context: `scripts/pytogether/findings/functions/12.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pytogether/backend/backend/settings/base.py:224:1`
- Checklist pattern: trigger token is inside a string literal, not executable code.

Source excerpt:

```
# ------------------ CODE TEMPLATES ---------------------

NONE_TEMPLATE = """name = input("Whats your name? ")
print(f"Hello from PyTogether, {name}!")"""
```

Why this is a false positive: the flagged `print(...)` is text inside the triple-quoted string assigned to `NONE_TEMPLATE` — a code template for generated starter files, not an executable statement in the module. The rule condition (`print(` call in library code) is not met because there is no call.

Checklist evidence: the detector's `printCallOutsideString` "cheap check" does not track triple-quoted multiline strings, so it reports the first `print(` token of the line even though the whole line is inside a string literal; the enclosing source shows the print is template content, not code.

## True positives

### Rule `BP-PY-46` — print Debugging In Library Code (29)

| Finding | Source | Reason |
| --- | --- | --- |
| 1 | backend/backend/jwt_auth_middleware.py:50 | `print(f"Share Token Validated: {share_data}")` in Django middleware module, not under `__main__` guard, not a test. |
| 3 | backend/backend/jwt_auth_middleware.py:61 | `print(f"JWT Error: {e}")` in library code. |
| 6 | backend/backend/jwt_auth_middleware.py:65 | `print(f"Generic JWT Error: {e}")` in library code. |
| 7 | backend/backend/jwt_auth_middleware.py:80 | `print("Share link expired")` in library code. |
| 8 | backend/backend/jwt_auth_middleware.py:83 | `print("Invalid share link signature")` in library code. |
| 9 | backend/backend/settings/__init__.py:11 | Module-level `print("--- PRODUCTION (VPS) MODE DETECTED ---")` at import time. |
| 10 | backend/backend/settings/__init__.py:14 | Module-level `print("--- SELF-HOSTING MODE DETECTED ---")`. |
| 11 | backend/backend/settings/__init__.py:17 | Module-level `print("--- DEVELOPMENT MODE DETECTED ---")`. |
| 14 | backend/backend/settings/dev.py:3 | Module-level `print("starting dev....")`. |
| 25 | backend/backend/settings/prod.py:7 | Module-level `print("setting up......")`. |
| 29 | backend/codes/consumers.py:35 | `print(self.group_id, ..., self.user.email, "is_member:", ...)` in Channels consumer. |
| 31 | backend/codes/consumers.py:45 | `print(f"Connection rejected: User {self.user.email} ...")`. |
| 35 | backend/codes/consumers.py:139 | `print(f"Error during disconnect cleanup: {e}")`. |
| 37 | backend/codes/consumers.py:178 | `print(f"Error in users_changed: {e}")`. |
| 40 | backend/codes/consumers.py:270 | `print(f"Error processing message: {e}")`. |
| 42 | backend/codes/consumers.py:317 | `print(f"Error sending voice room update: {e}")`. |
| 44 | backend/codes/consumers.py:347 | `print(f"Poison update rejected for project {project_id}: {e}")`. |
| 45 | backend/codes/consumers.py:358 | `print(f"Skipping update for project {project_id}: size exceeds limit")`. |
| 47 | backend/codes/consumers.py:368 | `print(f"Failed to acquire lock or write to Redis for project {project_id}: {e}")`. |
| 50 | backend/codes/tasks.py:36 | `print(f"Error persisting project {pid}: {e}")`. |
| 56 | backend/codes/tasks.py:77 | `print(f"Error cleaning up ghost project {pid}: {e}")`. |
| 59 | backend/usergroups/views.py:77 | `print("deleted the group")` in view handler. |
| 60 | backend/users/tokens.py:45 | `print("ok!")` in token refresh handler. |
| 64 | backend/utils/redis_helpers.py:45 | `print(f"YDoc size: {byte_size} bytes")` in helper module. |
| 65 | backend/utils/redis_helpers.py:47 | `print(f"Skipping save: codetext too thicc ({byte_size} bytes)")`. |
| 66 | backend/utils/redis_helpers.py:59 | `print(f"Saved {len(text)} chars to DB")`. |
| 67 | backend/utils/redis_helpers.py:61 | `print("Skipped DB save: No changes detected")`. |
| 68 | backend/utils/redis_helpers.py:69 | `print(f"Project {project_id} not found. Cleaning up Redis keys.")`. |
| 71 | backend/utils/redis_helpers.py:75 | `print(f"Error persisting YDoc to DB for project {project_id}: {e}")`. |

### Rule `CWE-215` — Insertion of Sensitive Information Into Debugging Code (2)

| Finding | Source | Reason |
| --- | --- | --- |
| 2 | backend/backend/jwt_auth_middleware.py:50 | print args contain `Token` (sensitive identifier) and are an interpolated f-string, not a pure literal; prints the share-token payload. |
| 32 | backend/codes/consumers.py:45 | print args contain `token` identifier and interpolate user data; not a pure string literal, so the rule condition matches. |

### Rule `BP-PY-1` — Bare Except Clause (14)

| Finding | Source | Reason |
| --- | --- | --- |
| 4 | backend/backend/jwt_auth_middleware.py:64 | `except Exception as e:` suite prints and returns; no re-raise. |
| 33 | backend/codes/consumers.py:138 | `except Exception as e:` with print handling only. |
| 36 | backend/codes/consumers.py:177 | `except Exception as e:` with print handling only. |
| 38 | backend/codes/consumers.py:191 | `except Exception:` suite is bare `return`. |
| 39 | backend/codes/consumers.py:269 | `except Exception as e:` with print handling only. |
| 41 | backend/codes/consumers.py:316 | `except Exception as e:` with print handling only. |
| 43 | backend/codes/consumers.py:346 | `except Exception as e:` suite prints, sends force_disconnect, returns; verified no re-raise in source. |
| 46 | backend/codes/consumers.py:367 | `except Exception as e:` with print handling only. |
| 48 | backend/codes/tasks.py:35 | `except Exception as e:` with print handling only. |
| 51 | backend/codes/tasks.py:43 | `except Exception:` with bare pass. |
| 55 | backend/codes/tasks.py:76 | `except Exception as e:` with print handling only. |
| 57 | backend/codes/tasks.py:81 | `except Exception:` with bare pass. |
| 61 | backend/users/views.py:153 | `except Exception as e:` returns error Response. |
| 69 | backend/utils/redis_helpers.py:74 | `except Exception as e:` with print handling only. |

### Rule `CWE-396` — Declaration of Catch for Generic Exception (5)

| Finding | Source | Reason |
| --- | --- | --- |
| 5 | backend/backend/jwt_auth_middleware.py:64 | `except Exception` root class handler, non-test module. |
| 34 | backend/codes/consumers.py:138 | `except Exception` root class handler. |
| 49 | backend/codes/tasks.py:35 | `except Exception` root class handler. |
| 62 | backend/users/views.py:153 | `except Exception` root class handler. |
| 70 | backend/utils/redis_helpers.py:74 | `except Exception` root class handler. |

### Rule `BP-PY-5` — Wildcard Import (3)

| Finding | Source | Reason |
| --- | --- | --- |
| 13 | backend/backend/settings/dev.py:1 | `from .base import *` in non-`__init__.py` file. |
| 24 | backend/backend/settings/prod.py:1 | `from .base import *` in non-`__init__.py` file. |
| 27 | backend/backend/settings/selfhost.py:1 | `from .base import *` in non-`__init__.py` file. |

### Rule `BP-PY-21` — Django DEBUG True In Settings (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 15 | backend/backend/settings/dev.py:5 | `DEBUG = True` in Django settings file; basename "dev" is not in the rule's `local_settings`/`dev_settings` skip list. |

### Rule `CWE-489` — Active Debug Code (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 16 | backend/backend/settings/dev.py:5 | `DEBUG = True` matches `pyDebugLineRE`; no debugger-call requirement. |

### Rule `CWE-756` — Missing Custom Error Page (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 17 | backend/backend/settings/dev.py:5 | `debugEnabledStart` matches `DEBUG = True`; explicit debug enablement is the whole rule condition. |

### Rule `BP-PY-23` — Django ALLOWED_HOSTS Empty Or Star (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 18 | backend/backend/settings/dev.py:6 | `ALLOWED_HOSTS = ["*"]` matches `allowedHostsHasStar`. |

### Rule `CWE-1188` — Initialization of a Resource with an Insecure Default (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 19 | backend/backend/settings/dev.py:6 | `ALLOWED_HOSTS = ["*"]` matches `pyWildcardHostsRE`. |

### Rule `CWE-547` — Use of Hard-coded, Security-relevant Constants (3)

| Finding | Source | Reason |
| --- | --- | --- |
| 20 | backend/backend/settings/dev.py:6 | `ALLOWED_HOSTS = ["*"]` matches `pyWeakSecuritySettingRE`. |
| 26 | backend/backend/settings/prod.py:10 | `SECURE_SSL_REDIRECT = False` matches `pyWeakSecuritySettingRE`. |
| 28 | backend/backend/settings/selfhost.py:29 | `SECURE_HSTS_SECONDS = 0` matches `pyWeakSecuritySettingRE`. |

### Rule `BP-PY-48` — CORS Allow Origins Star With Credentials (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 21 | backend/backend/settings/dev.py:9 | django-cors pair `CORS_ALLOW_ALL_ORIGINS = True` + `CORS_ALLOW_CREDENTIALS = True` both match. |

### Rule `BP-PY-50` — Django/Flask CSRF Or Session Cookie Insecure Flags (2)

| Finding | Source | Reason |
| --- | --- | --- |
| 22 | backend/backend/settings/dev.py:13 | `CSRF_COOKIE_SECURE = False` matches `cookieFlagFalseRe`. |
| 23 | backend/backend/settings/dev.py:14 | `SESSION_COOKIE_SECURE = False` matches `cookieFlagFalseRe`. |

### Rule `CWE-359` — Exposure of Private Personal Information to an Unauthorized Actor (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 30 | backend/codes/consumers.py:35 | print arguments contain `self.user.email`, matching `pyPersonalFieldRE` `user\.email`. |

### Rule `BP-PY-2` — Except Pass (2)

| Finding | Source | Reason |
| --- | --- | --- |
| 52 | backend/codes/tasks.py:43 | `except Exception:` suite is solely `pass`. |
| 58 | backend/codes/tasks.py:81 | `except Exception:` suite is solely `pass`. |

### Rule `CWE-390` — Detection of Error Condition Without Action (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 53 | backend/codes/tasks.py:43 | `exceptPassStart` matches the `except Exception: pass` handler. |

### Rule `CWE-1071` — Empty Code Block (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 54 | backend/codes/tasks.py:43 | `pyTierBEmptyExceptRE` matches the exception handler containing only `pass`. |

### Rule `CWE-209` — Generation of Error Message Containing Sensitive Information (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 63 | backend/users/views.py:154 | `Response({"detail": str(e)}, ...)` matches `pyExceptionStringRE` (`str(e)`), returning exception details in an HTTP response. |

## Uncertain findings

None. Every finding was resolvable from the rule condition and the shown source; no finding needed deployment assumptions to classify.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/pytogether/chunks/Chunk_1_25.txt`, `scripts/pytogether/chunks/Chunk_26_50.txt`, `scripts/pytogether/chunks/Chunk_51_71.txt`
- Function evidence: `scripts/pytogether/findings/functions/1.txt` … `scripts/pytogether/findings/functions/71.txt`
- Validation: `git diff --check` — pass

## Post-fix over-suppression audit (2026-08-02)

Mode B: fresh findings (68) are fewer than the audited TP count (70).

### Run metadata (fresh scan)

```yaml
timestamp: 2026-08-02 (post-fix binary rebuilt 2026-08-02 16:29, commit b5b8fde)
repository: pytogether
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pytogether
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pytogether
build_command: go build -o bin/goslop ./cmd/goslop
scan_command: ./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/pytogether/chunks -context-dir scripts/pytogether/findings/functions real-repos/pytogether
findings: 68
chunk_path: scripts/pytogether/chunks/Chunk_1_25.txt, Chunk_26_50.txt, Chunk_51_68.txt
function_context_path: scripts/pytogether/findings/functions/1.txt … 68.txt
```

### Over-suppressed true positives

Matching fresh findings to audited TPs by `Source:` (file:line), 68 of the 70 audited TPs are present in the fresh scan. The two audited TPs below are absent from the fresh scan (their sources appear only under BP-PY-46, not CWE-215). Both constructs still exist in the current source, so neither is fixed-removed; both were suppressed by the CWE-215 guard added in `b5b8fde` (`sensitiveIdentOutsideLiterals`, `internal/lang/python/detectors/cwe/fp_guards.go:337`, wired in `detectCWE215`, `rules_code_dynamic.go:132`). The pre-fix condition `sensitiveValueRE.MatchString(call.ArgsText) && !isPureStringLiteral(call.ArgsText)` matched the sensitive word anywhere in the args; the post-fix guard masks the whole f-string literal and only re-injects interpolation bodies, so the sensitive word inside the message text no longer counts. The old audited FP (finding 12, BP-PY-46, `base.py:224`) no longer appears in the fresh scan.

| Old finding ID | Rule | Source | One-line reason (from old audit) | Current status |
| --- | --- | --- | --- | --- |
| 2 | CWE-215 | backend/backend/jwt_auth_middleware.py:50 | print args contain `Token` (sensitive identifier) and are an interpolated f-string, not a pure literal; prints the share-token payload | suppressed-but-present |
| 32 | CWE-215 | backend/codes/consumers.py:45 | print args contain `token` identifier and interpolate user data; not a pure string literal, so the rule condition matches | suppressed-but-present |

### [ ] Old finding `2` — `CWE-215`

- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pytogether/backend/backend/jwt_auth_middleware.py:50`

Source excerpt (current file, line 50):

```
        if share_token and ":" in share_token:
            share_data = await self.validate_share_token(share_token)
            if share_data:
                # Add the share data to scope so consumers can use it
                scope['share_context'] = share_data
                print(f"Share Token Validated: {share_data}")
```

Why it is over-suppressed: the `print(f"Share Token Validated: {share_data}")` call still exists at line 50, so the construct was not removed. Under the pre-fix rule condition it matched because the args contain the sensitive identifier `Token` (`sensitiveValueRE`) and the f-string is not a pure literal (`isPureStringLiteral`), and the call prints the share-token payload — the sensitive data itself. Post-fix, `sensitiveIdentOutsideLiterals` masks the entire literal and re-injects only the interpolation body `share_data`, which does not match `sensitiveValueRE`, so the finding disappears even though sensitive payload data is still written to the console.

### [ ] Old finding `32` — `CWE-215`

- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/pytogether/backend/codes/consumers.py:45`

Source excerpt (current file, lines 44-47):

```
            if not self._validate_share_token(share_token, self.group_id, self.project_id):
                print(f"Connection rejected: User {self.user.email} is not a member and invalid token.")
                await self.close(code=4003)
                return
```

Why it is over-suppressed: the `print(...)` call still exists at line 45 and still interpolates user data, so the construct was not removed. Under the pre-fix rule condition it matched because the args contain the sensitive identifier `token` and are not a pure string literal. Post-fix, `sensitiveIdentOutsideLiterals` masks the message text (where "token" appears) and re-injects only `self.user.email`, which does not match `sensitiveValueRE`, so the finding is dropped (the email exposure itself remains covered separately by CWE-359 at `consumers.py:35`).

### Summary

- Over-suppressed TPs found: 2 (old IDs 2, 32 — both CWE-215)
- Fixed-removed: 0
- Validation: `git diff --check` — pass
