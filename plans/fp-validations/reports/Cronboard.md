# False-positive audit — Cronboard

## Run metadata

```yaml
timestamp: 2026-08-02T07:51:50Z
repository: Cronboard
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Cronboard
branch: main
commit: 0fa5f0d
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Cronboard
chunk_path: ./scripts/chunks
function_context_path: ./scripts/findings/functions
```

## Scan evidence

- Build command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir real-repos/Cronboard/scripts/chunks -context-dir real-repos/Cronboard/scripts/findings/functions real-repos/Cronboard`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir real-repos/Cronboard/scripts/chunks -context-dir real-repos/Cronboard/scripts/findings/functions real-repos/Cronboard`
- Findings: 69
- Chunks reviewed: `./scripts/chunks/Chunk_1_25.txt`, `./scripts/chunks/Chunk_26_50.txt`, `./scripts/chunks/Chunk_51_69.txt`
- Function contexts reviewed: `./scripts/findings/functions/30.txt`, `40.txt`, `60.txt`, `61.txt`, `62.txt`, `69.txt`, plus the full enclosing source files for all of them (`cron_encrypt.py`, `cron_logger.py`, `cron_servers.py`, `test_cron_ssh_modal.py`)

## Audit checklist

- [x] Read every assigned chunk under `./scripts/chunks`.
- [x] Read `./scripts/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 4 | 30, 40, 62, 69 |
| True positive | 65 | 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 31, 32, 33, 34, 35, 36, 37, 38, 39, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 63, 64, 65, 66, 67, 68 |
| Uncertain | 0 | — |

## False positives

### [ ] Finding `30` — `CWE-215`

- Function context: `./scripts/findings/functions/30.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Cronboard/src/cronboard/screens/cron_servers.py:236:33`
- Checklist pattern: debug sink arguments mention the word "password" but interpolate no sensitive value

Source excerpt:

```python
                            except Exception as e:
                                print(
                                    f"❌ Failed to decrypt password for {server_id}: {e}"
                                )
                                server_info["password"] = None
```

Why this is a false positive: the only interpolated values are `server_id` and the exception object `e` (raised by `fernet.decrypt`, whose messages never embed the token); the word "password" appears only as literal message text, so no password or token value reaches the output.

Checklist evidence: CWE-215's condition is "debug output includes a sensitive value"; the detector matches the arg text `password` in a non-literal f-string, but the f-string's interpolations are `server_id` and `e`, neither of which is the sensitive value.

### [ ] Finding `40` — `PERF-PY-25`

- Function context: `./scripts/findings/functions/40.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Cronboard/src/cronboard/services/cron_logging/cron_logger.py:31:1`
- Checklist pattern: lambda is constructed once as a `sorted()` key, not per loop element

Source excerpt:

```python
        return {
            p.stem: str(p)
            for p in sorted(
                log_dir.glob(f"{identificator}_*.log"), key=lambda p: p.stem
            )
        }
```

Why this is a false positive: the `key=lambda p: p.stem` expression is evaluated exactly once when `sorted()` is called — the lambda object is never constructed per comprehension element; `sorted` only invokes the single key function per element.

Checklist evidence: PERF-PY-25's condition is "heavy object or lambda is constructed per homogeneous loop element"; here the lambda is constructed once at the `sorted()` call site and only *called* per element, so the condition is not satisfied.

### [ ] Finding `62` — `CWE-88`

- Function context: `./scripts/findings/functions/62.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Cronboard/src/cronboard/services/encryption/cron_encrypt.py:74:12`
- Checklist pattern: the untrusted value is passed via stdin; the argv's only dynamic segment is a constant path

Source excerpt:

```python
    return subprocess.run(
        [
            "openssl",
            "enc",
            "-aes-256-cbc",
            "-salt",
            "-pbkdf2",
            "-pass",
            "file:" + str(KEY_FILE),
            "-base64",
            "-A",
        ],
        input=token.encode(),
        capture_output=True,
        check=False,
    ).stdout.decode()
```

Why this is a false positive: the attacker-influenced `token` is delivered through `input=` (stdin), not through the argument vector, so it can never be parsed as an option; the only non-literal argv segment is `"file:" + str(KEY_FILE)`, where `KEY_FILE` is a fixed config constant (`Path.home() / ".config/cronboard/secret.key"`) and cannot carry attacker-controlled option text.

Checklist evidence: CWE-88's condition is "dynamic value is embedded in a subprocess argument vector and can become an unintended option"; the dynamic user value is not in the argv at all, and the flagged `str(KEY_FILE)` segment is a static configuration path, not an external input.

### [ ] Finding `69` — `CWE-260`

- Function context: `./scripts/findings/functions/69.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Cronboard/tests/screens/test_cron_ssh_modal.py:96:13`
- Checklist pattern: literal `"password": "password"` in a test assertion fixture, not a configuration map

Source excerpt:

```python
    modal.dismiss.assert_called_once_with(
        {
            "hostname": "node9",
            "port": 2222,
            "username": "test",
            "password": "password",
            "ssh_key": False,
            "crontab_user": "root",
        }
    )
```

Why this is a false positive: the dict is the expected-call argument of a unit-test mock assertion in `tests/`, not a configuration file or runtime config map consumed by the product; the value is the placeholder string "password", not a credential that could leak to unauthorized actors.

Checklist evidence: CWE-260's condition is "configuration contains a literal password"; the flagged literal is test fixture data in a test module (the detector's `pyConfigPasswordRE` matches any `"password": "…"` dict entry without a test-scope exclusion), so no product configuration stores a password.

## True positives

### BP-PY-1 — Bare Except Clause (21)

| Finding ID | Source | Reason |
| --- | --- | --- |
| 1 | app.py:142 | `except Exception as e:` swallows all config-load failures |
| 4 | app.py:159 | `except Exception as e:` swallows all config-save failures |
| 6 | cron_creator.py:331 | `except Exception:` masks every validation error |
| 9 | cron_creator.py:435 | `except Exception as e:` swallows remote sync failure |
| 11 | cron_delete_confirmation.py:169 | `except Exception as e:` swallows remote sync failure |
| 17 | cron_delete_confirmation.py:210 | `except Exception as e:` swallows crontab write failure |
| 19 | cron_servers.py:117 | bare `except:` around ssh close |
| 23 | cron_servers.py:131 | `except Exception as e:` after a specific handler |
| 25 | cron_servers.py:198 | bare `except:` around ssh close |
| 28 | cron_servers.py:235 | `except Exception as e:` around password decryption |
| 31 | cron_servers.py:249 | `except Exception as e:` swallows config-load failure |
| 34 | cron_servers.py:281 | `except Exception as e:` swallows config-save failure |
| 35 | cron_servers.py:384 | bare `except:` around `_focus_tree()` |
| 36 | cron_settings.py:71 | `except Exception as e:` swallows settings-save failure |
| 38 | cron_settings.py:86 | `except Exception as e:` swallows settings-load failure |
| 43 | cron_wrapper.py:49 | `except Exception:` swallows fetch error |
| 51 | cron_wrapper.py:82 | `except Exception as e:` swallows connection failure |
| 54 | cron_wrapper.py:241 | `except Exception as e:` swallows file-write failure |
| 56 | cron_wrapper.py:260 | `except Exception as e:` swallows config-serialize failure |
| 58 | cron_wrapper.py:282 | `except Exception as e:` swallows config-serialize failure |
| 66 | cron_table.py:556 | `except Exception as e:` swallows crontab-write failure |

### BP-PY-2 — Except Pass (3)

| Finding ID | Source | Reason |
| --- | --- | --- |
| 20 | cron_servers.py:117 | handler body is only `pass` |
| 26 | cron_servers.py:198 | handler body is only `pass` |
| 44 | cron_wrapper.py:49 | handler body is only `pass` |

### BP-PY-46 — print Debugging In Library Code (27)

| Finding ID | Source | Reason |
| --- | --- | --- |
| 3 | app.py:143 | `print` in non-script module method |
| 5 | app.py:160 | `print` in non-script module method |
| 8 | cron_creator.py:358 | `print` in non-script module method |
| 10 | cron_creator.py:436 | `print` in non-script module method |
| 13 | cron_delete_confirmation.py:170 | `print` in non-script module method |
| 14 | cron_delete_confirmation.py:200 | `print` in non-script module method |
| 15 | cron_delete_confirmation.py:204 | `print` in non-script module method |
| 16 | cron_delete_confirmation.py:207 | `print` in non-script module method |
| 18 | cron_delete_confirmation.py:211 | `print` in non-script module method |
| 29 | cron_servers.py:236 | `print` in non-script module method |
| 32 | cron_servers.py:250 | `print` in non-script module method |
| 33 | cron_servers.py:252 | `print` in non-script module method |
| 39 | cron_settings.py:87 | `print` in non-script module method |
| 41 | cron_logger.py:46 | `print` in module-level service function |
| 42 | cron_logger.py:86 | `print` in module-level service function |
| 48 | cron_wrapper.py:70 | `print` in module-level service function |
| 49 | cron_wrapper.py:74 | `print` in module-level service function |
| 50 | cron_wrapper.py:80 | `print` in module-level service function |
| 52 | cron_wrapper.py:83 | `print` in module-level service function |
| 53 | cron_wrapper.py:162 | `print` in module-level service function |
| 55 | cron_wrapper.py:242 | `print` in module-level service function |
| 57 | cron_wrapper.py:261 | `print` in module-level service function |
| 59 | cron_wrapper.py:283 | `print` in module-level service function |
| 63 | cron_table.py:546 | `print` in non-script module method |
| 64 | cron_table.py:550 | `print` in non-script module method |
| 65 | cron_table.py:553 | `print` in non-script module method |
| 68 | cron_table.py:557 | `print` in non-script module method |

### CWE-396 — Declaration of Catch for Generic Exception (7)

| Finding ID | Source | Reason |
| --- | --- | --- |
| 2 | app.py:142 | generic `except Exception` handler |
| 7 | cron_creator.py:331 | generic `except Exception` handler |
| 12 | cron_delete_confirmation.py:169 | generic `except Exception` handler |
| 24 | cron_servers.py:131 | generic `except Exception` handler |
| 37 | cron_settings.py:71 | generic `except Exception` handler |
| 46 | cron_wrapper.py:49 | generic `except Exception` handler |
| 67 | cron_table.py:556 | generic `except Exception` handler |

### CWE-390 — Detection of Error Condition Without Action (2)

| Finding ID | Source | Reason |
| --- | --- | --- |
| 21 | cron_servers.py:117 | handler body is only `pass`, no action taken |
| 45 | cron_wrapper.py:49 | handler body is only `pass`, no action taken |

### CWE-1071 — Empty Code Block (2)

| Finding ID | Source | Reason |
| --- | --- | --- |
| 22 | cron_servers.py:117 | exception handler contains only `pass` |
| 47 | cron_wrapper.py:49 | exception handler contains only `pass` |

### CWE-1124 — Excessively Deep Nesting (1)

| Finding ID | Source | Reason |
| --- | --- | --- |
| 27 | cron_servers.py:232 | assignment nested 6 control-flow levels (`if` → `try` → `for` → `if` → `if` → `try`) |

### CWE-257 — Storing Passwords in a Recoverable Format (1)

| Finding ID | Source | Reason |
| --- | --- | --- |
| 61 | cron_encrypt.py:43 | `fernet.encrypt(password.encode())` stores the password in a symmetric-key-recoverable form |

### CWE-367 — TOCTOU Race Condition (1)

| Finding ID | Source | Reason |
| --- | --- | --- |
| 60 | cron_encrypt.py:17 | `os.path.exists(KEY_FILE)` check is followed by a separate `open(KEY_FILE, "wb")` use |

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `./scripts/chunks`
- Function evidence: `./scripts/findings/functions`
- Validation: `git diff --check` — pass
