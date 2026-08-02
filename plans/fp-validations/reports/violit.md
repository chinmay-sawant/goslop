# False-positive audit — violit

## Run metadata

```yaml
timestamp: 2026-08-02T10:00:00Z
repository: violit
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit
branch: main
commit: 8fae080f49f374b062172ed6ac71042539ad1f7a
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit
chunk_path: scripts/violit/chunks
function_context_path: scripts/violit/findings/functions
```

## Scan evidence

- Build command: `make build` (`go build -o bin/goslop ./cmd/goslop`)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/violit/chunks -context-dir scripts/violit/findings/functions real-repos/violit`
- Findings: `248`
- Chunks reviewed: `scripts/violit/chunks/Chunk_1_25.txt` … `scripts/violit/chunks/Chunk_226_248.txt` (all 10 chunk files)
- Function contexts reviewed: `scripts/violit/findings/functions/2.txt`, `5.txt`, `8.txt`, `12.txt`, `13.txt`, `14.txt`, `28.txt`, `29.txt`, `30.txt`, `67.txt`, `86.txt`, `90.txt`, `91.txt`, `111.txt`, `124.txt`, `131.txt`, `147.txt`, `148.txt`, `157.txt` … `164.txt`, `179.txt`, `213.txt`, `226.txt`, `235.txt` (for every proposed false positive / uncertain finding)

## Audit checklist

- [x] Read every assigned chunk under `scripts/violit/chunks`.
- [x] Read `scripts/violit/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 29 | 2, 5, 8, 12, 13, 14, 28, 29, 30, 67, 86, 90, 91, 111, 124, 147, 148, 157, 158, 159, 160, 161, 162, 163, 164, 179, 213, 226, 235 |
| True positive | 219 | 1, 3, 4, 6, 7, 9, 10, 11, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 87, 88, 89, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 149, 150, 151, 152, 153, 154, 155, 156, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 227, 228, 229, 230, 231, 232, 233, 234, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248 |
| Uncertain | 0 | — |

## False positives

### [ ] Finding `2` — `CWE-1121`

- Function context: `scripts/violit/findings/functions/2.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/examples/01_demo_showcase/old_archive_demo_showcase.py:178:38`
- Checklist pattern: rule condition “function has at least twelve visible control-flow branches” (counter counts `if `/`elif `/`for `/`while `/`except ` substrings)

Source excerpt:

```
def _submit_chat_prompt(prompt: str):
    if isinstance(prompt, dict):
        prompt_payload = {
            "text": str(prompt.get("text") or "").strip(),
            "files": [entry for entry in list(prompt.get("files") or []) if entry is not None],
            ...
    if not cleaned or chat_busy.value:
        return
    ...
    def run():
        result = _pseudo_chat_reply(cleaned)
        if isinstance(result, str):
```

Why this is a false positive: the function contains only 7 real `if` statements and, in a nested helper, 1 `if` + 1 `for` — at most 9 control-flow statements. The substring counter only reaches 12 by counting `for `/`if ` inside list comprehensions (`[entry for entry in ... if entry is not None]`, `... for entry in ...`) which are expressions, not branches.

Checklist evidence: counted the function body from the source file; real control-flow statements < 12.

### [ ] Finding `5` — `CWE-89`

- Function context: `scripts/violit/findings/functions/5.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/examples/02_violit_blog/violit_advanced_blog.py:43:6`
- Checklist pattern: rule condition “dynamic SQL string reaches execute/executemany (use bound parameters)”

Source excerpt:

```
def db_query(query, params=(), fetch=False):
    conn = sqlite3.connect(DB_NAME)
    conn.row_factory = sqlite3.Row
    c = conn.cursor()
    c.execute(query, params)
```

Why this is a false positive: the flagged call already supplies a bound-params argument (`params`), which is exactly the rule's prescribed fix; `query` is a pass-through parameter of a helper and every caller in the same file passes a static literal SQL string with `?` placeholders (e.g. `db_query("SELECT * FROM posts WHERE id = ?", (post_id,), fetch=True)`). No untrusted data is interpolated into SQL.

Checklist evidence: checked all `db_query(...)` call sites in the enclosing source — all literal SQL + separate params tuple.

### [ ] Finding `8` — `CWE-89`

- Function context: `scripts/violit/findings/functions/8.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/examples/02_violit_blog/violit_blog.py:40:6`
- Checklist pattern: rule condition “dynamic SQL string reaches execute/executemany (use bound parameters)”

Source excerpt:

```
def db_query(query, params=(), fetch=False):
    conn = sqlite3.connect(DB_NAME)
    conn.row_factory = sqlite3.Row
    c = conn.cursor()
    c.execute(query, params)
```

Why this is a false positive: identical helper-wrapper pattern as finding 5 — bound parameters are passed on the flagged line and every caller passes literal placeholder SQL (`"SELECT id, username FROM users WHERE username = ? AND password = ?"`, etc.); no interpolation of untrusted data.

Checklist evidence: checked all `db_query(...)` call sites in the enclosing source — all literal SQL + separate params tuple.

### [ ] Finding `12` — `BP-PY-37`

- Function context: `scripts/violit/findings/functions/12.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/examples/03_multi_db_sqlite_editor/demo_multi_db_sqlite_editor.py:230:17`
- Checklist pattern: rule condition “DB-API `cursor.execute` builds SQL with f-string or % format; use bound parameters”

Source excerpt:

```
def persist_editor_change(event: dict) -> str:
    field = event.get("field")
    if field not in HR_FIELDS | PROJECT_FIELDS:
        raise ValueError(f"Unsupported field: {field}")
    ...
        with _connect(HR_DB) as conn:
            conn.execute(f"UPDATE employees SET {field} = ? WHERE emp_id = ?", (stored_value, emp_id))
```

Why this is a false positive: the f-string interpolates only the column identifier, which the same function restricts to the developer-defined allowlists `HR_FIELDS = {"name", "email", "department", "is_active"}` and `PROJECT_FIELDS = {"project_name", "status", "deadline", "allocation_pct"}` before the call; all values are bound with `?` placeholders, so no untrusted characters can reach the SQL string.

Checklist evidence: the interpolation site is unreachable for any `field` outside the fixed allowlist, and the value payload is bound via placeholders.

### [ ] Finding `13` — `CWE-89`

- Function context: `scripts/violit/findings/functions/13.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/examples/03_multi_db_sqlite_editor/demo_multi_db_sqlite_editor.py:230:17`
- Checklist pattern: rule condition “dynamic SQL string reaches execute/executemany (use bound parameters)”

Source excerpt:

```
    field = event.get("field")
    if field not in HR_FIELDS | PROJECT_FIELDS:
        raise ValueError(f"Unsupported field: {field}")
    ...
    conn.execute(f"UPDATE employees SET {field} = ? WHERE emp_id = ?", (stored_value, emp_id))
```

Why this is a false positive: same construct as finding 12 — the only interpolated token is the column name, constrained by an in-function allowlist check to the fixed sets `HR_FIELDS` / `PROJECT_FIELDS`; the data values travel via `?` bound parameters, so the “Improper Neutralization” condition (untrusted data in the SQL command) is not met.

Checklist evidence: the SQL identifier is allowlist-constrained in the same function; values are bound parameters.

### [ ] Finding `14` — `BP-PY-37`

- Function context: `scripts/violit/findings/functions/14.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/examples/03_multi_db_sqlite_editor/demo_multi_db_sqlite_editor.py:234:13`
- Checklist pattern: rule condition “DB-API `cursor.execute` builds SQL with f-string or % format; use bound parameters”

Source excerpt:

```
    with _connect(PROJECT_DB) as conn:
        conn.execute(f"UPDATE assignments SET {field} = ? WHERE emp_id = ?", (value, emp_id))
```

Why this is a false positive: same allowlist-constrained identifier interpolation as finding 12 (second statement in `persist_editor_change`, guarded by the same `field not in HR_FIELDS | PROJECT_FIELDS` check); values are bound with `?` placeholders, so no untrusted data is interpolated.

Checklist evidence: identifier is allowlist-constrained; values bound.

### [ ] Finding `28` — `BP-PY-46`

- Function context: `scripts/violit/findings/functions/28.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/examples/99_github_issues/github_issue_99_codeblock.py:27:1`
- Checklist pattern: rule condition “print( in library code (not under __main__ guard, not tests)” with `printCallOutsideString` per-line check

Source excerpt:

```
PYTHON_SAMPLE = '''from dataclasses import dataclass
...
def greet(user: User) -> str:
    status = "active" if user.active else "inactive"
    return f"Hello, {user.name} ({status})"


print(greet(User("Violit")))
'''
```

Why this is a false positive: the `print(...)` is inside the triple-quoted string constant `PYTHON_SAMPLE` — sample source text displayed by the demo, never executed. The rule's per-line `printCallOutsideString` check cannot see multi-line string state, so it flags a `print(` that is inside a string literal.

Checklist evidence: line 27 lies inside the `'''` block that starts at line 13 and ends at line 28; the `print(` occurrence is string content, not code.

### [ ] Finding `29` — `CWE-93`

- Function context: `scripts/violit/findings/functions/29.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/app.py:85:21`
- Checklist pattern: rule condition “externally influenced value is written into an HTTP response header without CRLF neutralization”

Source excerpt:

```
    def _cache_control_for_scope(self, scope) -> str:
        raw_path = str(scope.get("path") or "").replace("\\", "/").lower()
        ...
        if "/static/runtime/" in raw_path:
            if query.get("v"):
                return "public, max-age=31536000, immutable"
            return "public, max-age=0, must-revalidate"
        ...
        return "public, max-age=0, must-revalidate"

    def file_response(self, full_path, stat_result, scope, status_code=200):
        response = super().file_response(full_path, stat_result, scope, status_code)
        if not response.headers.get("Cache-Control"):
            response.headers["Cache-Control"] = self._cache_control_for_scope(scope)
```

Why this is a false positive: the header value is the return of an internal helper whose every branch returns one of four fixed string literals; no request data is embedded in the value, so CR/LF injection is impossible. The rule's “dynamic” classification is only triggered because the RHS is a call expression, not because the value is externally influenced.

Checklist evidence: `_cache_control_for_scope` returns only hardcoded literal strings; the header value cannot contain CR or LF.

### [ ] Finding `30` — `CWE-1121`

- Function context: `scripts/violit/findings/functions/30.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/app.py:124:434`
- Checklist pattern: rule condition “function has at least twelve visible control-flow branches”

Source excerpt:

```
    def __init__(self, mode: Optional[str] = None, title="Violit App", ...):
        self._mode_is_explicit = mode is not None
        self.mode = (mode or 'ws').strip().lower()
        ...
        self.db = None
        if db is not None:
            from .db import ViolItDB, normalize_db_url
            self.db = ViolItDB(normalize_db_url(db), migrate=migrate)
        ...
        if static_path.exists():
            self.fastapi.mount("/static", ...)
        ...
        self.debug_mode = '--debug' in sys.argv
```

Why this is a false positive: the constructor spans ~370 lines but contains only 7 real `if` statements (the rest are attribute assignments). The substring counter reaches 30 `if ` hits by counting conditional expressions (`mode or 'ws'`, ternaries) and `if` text inside string literals/nested definitions — none are control-flow branches.

Checklist evidence: AST count of `If`/`For`/`While`/`Try` statements directly in `App.__init__` (lines 124–494) = 7 < 12.

### [ ] Finding `67` — `BP-PY-40`

- Function context: `scripts/violit/findings/functions/67.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/app_launcher.py:298:15`
- Checklist pattern: rule condition “threading.Thread started without .join … avoid fire-and-forget non-daemon threads” (the detector explicitly skips lines with `daemon=True`)

Source excerpt:

```
        thread = threading.Thread(target=server_manager, daemon=True)
        thread.start()
```

Why this is a false positive: the thread is a daemon thread — the detector's own exemption (“Daemon-only policy: if every Thread construction uses daemon=True, miss”). The exemption only checks the `.start(` line for `daemon=True`, and here the constructor with `daemon=True` is on the preceding line, so the same-line check misses it.

Checklist evidence: the thread is created with `daemon=True` (line 297) immediately before the flagged `.start()` (line 298).

### [ ] Finding `86` — `BP-PY-13`

- Function context: `scripts/violit/findings/functions/86.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/app_runtime.py:414:32`
- Checklist pattern: rule condition “A secret-like name is assigned a non-empty string literal in source”

Source excerpt:

```
            csrf_token = self._generate_csrf_token(sid) if sid and self.csrf_enabled else ""
            csrf_script = f'<script>window._csrf_token = "{csrf_token}";</script>' if csrf_token else ""
```

Why this is a false positive: the regex matched the JavaScript identifier `window._csrf_token = ` *inside* an HTML/JS string template, not a Python assignment. The “value” is a runtime-generated CSRF token interpolated into an f-string; no literal secret is assigned in source.

Checklist evidence: the matched construct is a JS variable name inside a string literal; the RHS is an f-string with a runtime-generated value, not a hardcoded secret.

### [ ] Finding `90` — `CWE-215`

- Function context: `scripts/violit/findings/functions/90.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/app_runtime.py:420:17`
- Checklist pattern: rule condition “debug output includes a sensitive value; redact it before logging”

Source excerpt:

```
                print(f"[DEBUG] Session ID: {sid[:8] if sid else 'None'}...")
                print(f"[DEBUG] CSRF enabled: {self.csrf_enabled}")
                print(f"[DEBUG] CSRF token generated: {bool(csrf_token)}")
```

Why this is a false positive: the flagged debug line prints `bool(csrf_token)` — only whether a token exists. The token value itself is never included, so no sensitive value is inserted into debug output.

Checklist evidence: the printed expression is `bool(csrf_token)`, which can only be `True`/`False`.

### [ ] Finding `91` — `BP-PY-32`

- Function context: `scripts/violit/findings/functions/91.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/app_runtime.py:553:24`
- Checklist pattern: rule condition “FileResponse / static file helpers use a path from user input without confinement”

Source excerpt:

```
                store = get_session_store()
                registry = store.get("_vl_media_sources") or {}
                payload = registry.get(media_id)
                if not isinstance(payload, dict):
                    return Response(status_code=404)

                media_path = str(payload.get("path") or "").strip()
                if not media_path or not os.path.exists(media_path) or not os.path.isfile(media_path):
                    return Response(status_code=404)
                ...
                return FileResponse(
                    path=media_path, ...
```

Why this is a false positive: the path is not user input — it comes from the application's own session-store media registry, populated by the framework when the developer creates media widgets. The only user-controlled value is `media_id`, which selects among app-registered entries; the request cannot name an arbitrary filesystem path.

Checklist evidence: `media_path` is read from `registry` (app-populated `_vl_media_sources`), not derived from request parameters.

### [ ] Finding `111` — `CWE-1121`

- Function context: `scripts/violit/findings/functions/111.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/background.py:188:52`
- Checklist pattern: rule condition “function has at least twelve visible control-flow branches”

Source excerpt:

```
    def _run(self, sid: str, current_view_id: str):
        t = session_ctx.set(sid) if sid else None
        view_token = view_ctx.set(current_view_id) if current_view_id else None
        ...
        try:
            self._result = self._fn()
            if self._cancel_event.is_set():
                ...
            if self._on_complete:
                ...
        except CancelledError:
            ...
        except Exception as e:
            ...
```

Why this is a false positive: the function contains 11 real control-flow statements (AST count excluding nested definitions), one short of the 12-branch threshold. The substring counter reaches ≥12 by counting ternary `if` expressions (`session_ctx.set(sid) if sid else None`, `view_ctx.set(...) if ... else ...`) as branches.

Checklist evidence: AST count of `If`/`For`/`While`/`Try` in `BackgroundTask._run` = 11 < 12.

### [ ] Finding `124` — `CWE-117`

- Function context: `scripts/violit/findings/functions/124.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/background.py:312:17`
- Checklist pattern: rule condition “externally influenced input … without neutralizing line-breaking control characters”

Source excerpt:

```
                logger.debug(f"[background] Pushed {len(dirty)} updates to session {sid[:8]}...")
```

Why this is a false positive: both interpolated values are internally generated and cannot carry CR/LF: `len(dirty)` is an integer, and `sid[:8]` is a truncated prefix of a server-generated session id (not attacker-controlled text).

Checklist evidence: the formatted values are an `int` and a truncated internally generated session id; no externally controlled string is written.

### [ ] Finding `147` — `CWE-117`

- Function context: `scripts/violit/findings/functions/147.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/db.py:129:17`
- Checklist pattern: rule condition “externally influenced input … without neutralizing line-breaking control characters”

Source excerpt:

```
            # New table -> already created by create_all, so skip
            if table_name not in existing_tables:
                logger.info(f"[violit:db] ✅ New table created: {table_name}")
                continue
```

Why this is a false positive: `table_name` iterates SQLModel model metadata (`SQLModel.metadata.tables`), i.e. developer-defined class names, not externally controlled input — it cannot contain CR/LF from an attacker.

Checklist evidence: the interpolated value comes from the developer-defined SQLModel schema, not from external input.

### [ ] Finding `148` — `CWE-89`

- Function context: `scripts/violit/findings/functions/148.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/db.py:151:33`
- Checklist pattern: rule condition “dynamic SQL string reaches execute/executemany (use bound parameters)”

Source excerpt:

```
                    for col in table.columns:
                        if col.name not in existing_cols:
                            col_type_str = col.type.compile(dialect=self._engine.dialect)
                            default_clause = self._build_default_clause(col)
                            sql = (
                                f'ALTER TABLE "{table_name}" '
                                f'ADD COLUMN "{col.name}" {col_type_str}{default_clause}'
                            )
                            conn.execute(text(sql))
```

Why this is a false positive: the SQL is a DDL statement built exclusively from SQLModel metadata — `table_name`, `col.name`, and the compiled `col.type` all originate from developer-defined model classes. No user or request data reaches the statement; it is a schema-migration path with no injection surface.

Checklist evidence: all interpolated values derive from `table.columns` / `col.name` / `col.type.compile(...)` (model metadata), not from external input.

### [ ] Finding `157` — `BP-PY-12`

- Function context: `scripts/violit/findings/functions/157.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/db.py:234:29`
- Checklist pattern: rule condition “eval/exec on dynamic input enables arbitrary code execution”

Source excerpt:

```
            from sqlmodel import select
            with app.db.session() as s:
                results = s.exec(
                    select(Task).where(Task.done == False).limit(10)
                ).all()
```

Why this is a false positive: `s.exec(...)` is a method call on a `sqlmodel.Session` object (a typed query-execution API), not the Python builtin `exec`. The detector matches any identifier `exec(` without verifying the receiver.

Checklist evidence: the receiver `s` is a SQLModel session; `exec` is a method name, not the builtin.

### [ ] Finding `158` — `BP-PY-14`

- Function context: `scripts/violit/findings/functions/158.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/db.py:279:20`
- Checklist pattern: rule condition “`requests` HTTP calls omit `timeout=`” (detector regex `\b(session|sess|req_session)\.(get|post|…)\s*\(`)

Source excerpt:

```
        _check_sqlmodel()
        with Session(self._engine) as session:
            return session.get(model, pk)
```

Why this is a false positive: `session` here is a SQLModel `Session` (database session), and `session.get(model, pk)` is an ORM primary-key lookup, not an HTTP request. The detector flags any `session.get(` by variable name without checking the receiver type.

Checklist evidence: `Session` is imported from `sqlmodel` and wraps `self._engine`; no `requests` import exists in the file.

### [ ] Finding `159` — `BP-PY-12`

- Function context: `scripts/violit/findings/functions/159.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/db.py:294:28`
- Checklist pattern: rule condition “eval/exec on dynamic input enables arbitrary code execution”

Source excerpt:

```
        with Session(self._engine) as session:
            stmt = select(model)
            for cond in conditions:
                stmt = stmt.where(cond)
            return session.exec(stmt).first()
```

Why this is a false positive: `session.exec(stmt)` is the SQLModel Session query-execution method (executes a typed `select` statement), not the Python builtin `exec`.

Checklist evidence: the receiver is a SQLModel `Session`; the argument is a typed statement object.

### [ ] Finding `160` — `BP-PY-12`

- Function context: `scripts/violit/findings/functions/160.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/db.py:306:33`
- Checklist pattern: rule condition “eval/exec on dynamic input enables arbitrary code execution”

Source excerpt:

```
        _check_sqlmodel()
        with Session(self._engine) as session:
            return list(session.exec(select(model)).all())
```

Why this is a false positive: `session.exec(select(model))` is the SQLModel typed query-execution method, not the `exec` builtin; no dynamic code is executed.

Checklist evidence: receiver is a SQLModel `Session`; argument is a `select(...)` statement.

### [ ] Finding `161` — `BP-PY-12`

- Function context: `scripts/violit/findings/functions/161.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/db.py:339:33`
- Checklist pattern: rule condition “eval/exec on dynamic input enables arbitrary code execution”

Source excerpt:

```
            if limit:
                stmt = stmt.limit(limit)
            return list(session.exec(stmt).all())
```

Why this is a false positive: `session.exec(stmt)` is the SQLModel Session query-execution method, not the `exec` builtin.

Checklist evidence: receiver is a SQLModel `Session`; argument is a typed statement object.

### [ ] Finding `162` — `BP-PY-14`

- Function context: `scripts/violit/findings/functions/162.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/db.py:353:13`
- Checklist pattern: rule condition “`requests` HTTP calls omit `timeout=`” (detector regex matches `session.delete(` by name)

Source excerpt:

```
            # use merge to safely delete objects fetched from other sessions
            managed = session.merge(obj)
            session.delete(managed)
            session.commit()
```

Why this is a false positive: `session.delete(managed)` is a SQLModel ORM delete on the database session, not a `requests` HTTP call; the file contains no `requests` usage.

Checklist evidence: `session` is a SQLModel `Session` bound to `self._engine`; `delete`/`commit` are ORM methods.

### [ ] Finding `163` — `BP-PY-14`

- Function context: `scripts/violit/findings/functions/163.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/db.py:369:17`
- Checklist pattern: rule condition “`requests` HTTP calls omit `timeout=`” (detector regex matches `session.delete(` by name)

Source excerpt:

```
            for row in rows:
                managed = session.merge(row)
                session.delete(managed)
            session.commit()
```

Why this is a false positive: same construct as finding 162 — ORM delete on a SQLModel database session, not an HTTP call.

Checklist evidence: `session` is a SQLModel `Session`; `delete`/`commit` are ORM methods.

### [ ] Finding `164` — `BP-PY-12`

- Function context: `scripts/violit/findings/functions/164.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/db.py:387:28`
- Checklist pattern: rule condition “eval/exec on dynamic input enables arbitrary code execution”

Source excerpt:

```
            for cond in conditions:
                stmt = stmt.where(cond)
            return session.exec(stmt).one()
```

Why this is a false positive: `session.exec(stmt)` is the SQLModel Session query-execution method, not the Python `exec` builtin.

Checklist evidence: receiver is a SQLModel `Session`; argument is a typed statement object.

### [ ] Finding `179` — `CWE-1121`

- Function context: `scripts/violit/findings/functions/179.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/widgets/chart_widgets.py:19:27`
- Checklist pattern: rule condition “function has at least twelve visible control-flow branches”

Source excerpt:

```
    def _normalize(value):
        if isinstance(value, np.ndarray):
            return [_normalize(item) for item in value.tolist()]
        if isinstance(value, np.generic):
            return value.item()
        if isinstance(value, dict):
            if set(value.keys()) == {"dtype", "bdata"}:
                try:
                    array = np.frombuffer(...)
                    return [_normalize(item) for item in array.tolist()]
                except Exception:
                    return value
            return {key: _normalize(item) for key, item in value.items()}
        if isinstance(value, list):
            return [_normalize(item) for item in value]
        if isinstance(value, tuple):
            return [_normalize(item) for item in value]
        return value
```

Why this is a false positive: the function has 6 `if` statements and 1 `except` (7 real branches). The substring count reaches 12 only by counting `for `/`if ` inside list comprehensions (`[... for item in ...]`, `{key: ... for key, item in ...}`), which are expressions, not control-flow branches.

Checklist evidence: AST count of `If`/`For`/`While`/`Try` statements = 7 < 12.

### [ ] Finding `213` — `CWE-1121`

- Function context: `scripts/violit/findings/functions/213.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/widgets/form_widgets.py:67:23`
- Checklist pattern: rule condition “function has at least twelve visible control-flow branches”

Source excerpt:

```
        def builder():
            token = rendering_ctx.set(cid)
            bt = text() if callable(text) else text
            ...
            icon_html = f'...' if icon and not any(ord(c) > 127 for c in str(icon)) else ''
            ...
            if use_container_width:
                host_style = merge_style(host_style, "width:100%;")
            if height not in (None, "", "auto"):
                if height == "fill":
                    ...
                elif isinstance(height, (int, float)):
                    ...
                else:
                    ...
            if disabled:
                inner_html = inner_html.replace(...)
            return Component(None, id=cid, content=inner_html)
```

Why this is a false positive: the closure contains only 6 real control-flow statements (3 `if`/`elif` blocks). The substring count reaches 12 by counting ternary `... if ... else ...` expressions and `for` inside generator expressions (`any(ord(c) > 127 for c in str(icon))`), none of which are branches.

Checklist evidence: AST count of real control-flow statements in the `builder()` at line 67 = 6 < 12.

### [ ] Finding `226` — `CWE-1121`

- Function context: `scripts/violit/findings/functions/226.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/widgets/input_widgets.py:415:23`
- Checklist pattern: rule condition “function has at least twelve visible control-flow branches”

Source excerpt:

```
        def builder():
            token = rendering_ctx.set(cid)
            cv = s.value
            ...
            opts_html = ""
            for i_opt, opt in enumerate(options):
                sel = 'checked' if opt == cv else ''
                escaped_opt = html_lib.escape(str(opt), quote=True)
                ...
                if captions and i_opt < len(captions) and captions[i_opt]:
                    caption_html = f'...'
                    if horizontal:
                        option_style = ...
                elif horizontal:
                    option_style = ...
                opts_html += f'<wa-radio ...>'
```

Why this is a false positive: the closure contains 1 real `for` loop and 6 `if` statements (7 real branches). The substring count reaches 12 by counting ternary `if` expressions (`'checked' if opt == cv else ''`, `f'...' if _part_cls else ''`, etc.) as branches.

Checklist evidence: AST count of real control-flow statements in the `builder()` at line 415 = 7 < 12.

### [ ] Finding `235` — `CWE-1121`

- Function context: `scripts/violit/findings/functions/235.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/widgets/layout_widgets.py:326:107`
- Checklist pattern: rule condition “function has at least twelve visible control-flow branches”

Source excerpt:

```
    def expander(self, label, expanded=False, icon=None, cls: str = "", style: str = "", key: Any = None):
        cid = self._resolve_widget_cid("expander", key)
        summary_cid = f"{cid}__summary"
        details_cid = f"{cid}__details"

        class ExpanderContext:
            def __init__(self, app, expander_id, label, expanded, icon=None, user_cls="", user_style=""):
                ...
            def __enter__(self):
                def summary_builder():
                    ...
                    if isinstance(self.label, (State, ComputedState)):
                        resolved_label = str(self.label.value)
                    elif callable(self.label):
                        resolved_label = str(self.label())
                    else:
                        resolved_label = str(self.label)
```

Why this is a false positive: `expander` itself is a thin wrapper that defines a context-manager class; it contains only 8 real control-flow statements (AST count excluding the nested class/function bodies). The substring counter reaches 12 by including branches of the nested `summary_builder`/`__init__`/`__enter__` definitions, which are not executed as part of `expander`'s own flow.

Checklist evidence: AST count of real control-flow statements in `expander` (excluding nested definitions) = 8 < 12.

## Uncertain findings

None. Finding 131 (CWE-829) reclassified as a true positive: `runpy.run_path(script_path)` with a dynamic path matches `detectCWE829` (`isDynamicExpr` on the first argument).


## True positives

### BP-PY-1 — Bare Except Clause

| Finding | Source | Reason |
| --- | --- | --- |
| 17 | examples/04_llm_chat_application/demo_app_gemini_agent_highlevel.py:673 | `except Exception as exc:` in `reply()` generator |
| 21 | examples/04_llm_chat_application/demo_app_gemini_agent_primitive.py:677 | `except Exception as exc:` in `reply()` generator |
| 35 | src/violit/app.py:1537 | `except Exception as e:` around builder().render() |
| 37 | src/violit/app.py:1636 | `except Exception as e:` around builder() |
| 38 | src/violit/app.py:1872 | `except Exception:` silently discards icon-set failure |
| 42 | src/violit/app.py:1879 | `except Exception:` silently discards patching failure |
| 44 | src/violit/app.py:2068 | `except Exception:` + pass around click dispatch |
| 46 | src/violit/app_launcher.py:41 | `except Exception:` on socket reachability probe |
| 63 | src/violit/app_launcher.py:198 | `except Exception as exc:` around uvicorn start |
| 66 | src/violit/app_launcher.py:288 | `except Exception as exc:` around webview reload |
| 68 | src/violit/app_launcher.py:321 | `except Exception:` + pass on window switch |
| 71 | src/violit/app_launcher.py:330 | `except Exception:` + pass on terminate |
| 73 | src/violit/app_launcher.py:388 | `except Exception:` + pass on logging filter setup |
| 75 | src/violit/app_launcher.py:394 | `except Exception as db_exc:` on startup migration |
| 78 | src/violit/app_launcher.py:489 | `except Exception:` + pass on window switch |
| 82 | src/violit/app_runtime.py:149 | `except Exception as exc:` on queue put |
| 99 | src/violit/app_support.py:169 | `except Exception:` + pass on stdout reconfigure |
| 103 | src/violit/auth.py:67 | `except Exception:` returns False from bcrypt check |
| 112 | src/violit/background.py:214 | `except Exception:` + pass on cancel handler |
| 117 | src/violit/background.py:238 | `except Exception:` + pass on cancel handler |
| 119 | src/violit/background.py:243 | `except Exception as e:` captures task error |
| 121 | src/violit/background.py:252 | `except Exception:` + pass on error handler |
| 127 | src/violit/background.py:325 | `except Exception as e:` on push failure |
| 135 | src/violit/cli.py:37 | `except Exception as e:` on runpy failure |
| 143 | src/violit/db.py:120 | `except Exception as e:` on schema inspect |
| 155 | src/violit/db.py:201 | `except Exception as e:` on per-table migration |
| 171 | src/violit/engine.py:75 | `except Exception:` drops socket on send error |
| 176 | src/violit/state.py:488 | `except Exception as e:` around callbacks |
| 180 | src/violit/widgets/chart_widgets.py:29 | `except Exception:` returns raw value |
| 182 | src/violit/widgets/chart_widgets.py:613 | `except Exception:` + pass on rcParams |
| 186 | src/violit/widgets/chart_widgets.py:628 | `except Exception:` + pass on font set |
| 188 | src/violit/widgets/chat_widgets.py:87 | `except Exception:` + pass on setattr |
| 193 | src/violit/widgets/chat_widgets.py:99 | `except Exception:` on json.loads |
| 194 | src/violit/widgets/chat_widgets.py:132 | `except Exception:` on int cast |
| 196 | src/violit/widgets/custom_widgets.py:61 | `except Exception:` on json.loads |
| 198 | src/violit/widgets/data_widgets.py:44 | `except Exception:` on pandas import |
| 200 | src/violit/widgets/data_widgets.py:57 | `except Exception:` on pandas import |
| 201 | src/violit/widgets/data_widgets.py:74 | `except Exception:` returns False |
| 202 | src/violit/widgets/data_widgets.py:531 | `except Exception:` on json.loads |
| 204 | src/violit/widgets/data_widgets.py:559 | bare `except:` on DataFrame cast |
| 206 | src/violit/widgets/data_widgets.py:738 | bare `except:` on DataFrame cast |
| 207 | src/violit/widgets/data_widgets.py:881 | `except Exception as exc:` on validator |
| 208 | src/violit/widgets/data_widgets.py:898 | `except Exception as exc:` on on_change |
| 209 | src/violit/widgets/data_widgets.py:901 | `except Exception:` + pass |
| 218 | src/violit/widgets/form_widgets.py:183 | `except Exception as e:` on native save |
| 222 | src/violit/widgets/form_widgets.py:556 | `except Exception as e:` on file save |
| 224 | src/violit/widgets/form_widgets.py:618 | `except Exception as e:` on file save |
| 225 | src/violit/widgets/input_widgets.py:29 | bare `except:` on b64decode |
| 230 | src/violit/widgets/input_widgets.py:933 | bare `except:` on json.loads |
| 231 | src/violit/widgets/input_widgets.py:953 | `except Exception as e:` on upload |
| 234 | src/violit/widgets/input_widgets.py:1320 | `except Exception:` on json.loads |
| 236 | src/violit/widgets/layout_widgets.py:580 | `except Exception:` + pass on render |
| 241 | src/violit/widgets/layout_widgets.py:586 | `except Exception:` + pass on render |
| 243 | src/violit/widgets/media_widgets.py:32 | `except Exception:` returns path |
| 245 | src/violit/widgets/media_widgets.py:96 | `except Exception:` returns str(image) |
| 246 | src/violit/widgets/media_widgets.py:219 | `except Exception:` returns str(audio) |

### BP-PY-2 — Except Pass

| Finding | Source | Reason |
| --- | --- | --- |
| 39 | src/violit/app.py:1872 | `except Exception: pass` on icon set |
| 43 | src/violit/app.py:1879 | `except Exception: pass` on patching |
| 45 | src/violit/app.py:2068 | `except Exception: pass` on click dispatch |
| 58 | src/violit/app_launcher.py:107 | `except (ProcessLookupError, OSError): pass` on terminate |
| 60 | src/violit/app_launcher.py:117 | `except (ProcessLookupError, OSError): pass` on kill |
| 64 | src/violit/app_launcher.py:228 | `except (ProcessLookupError, OSError): pass` on terminate |
| 65 | src/violit/app_launcher.py:238 | `except (ProcessLookupError, OSError): pass` on kill |
| 69 | src/violit/app_launcher.py:321 | `except Exception: pass` on window switch |
| 72 | src/violit/app_launcher.py:330 | `except Exception: pass` on terminate |
| 74 | src/violit/app_launcher.py:388 | `except Exception: pass` on logging filter |
| 79 | src/violit/app_launcher.py:489 | `except Exception: pass` on window switch |
| 93 | src/violit/app_support.py:35 | `except OSError: pass` on mtime stat |
| 98 | src/violit/app_support.py:60 | `except OSError: pass` on watch poll |
| 100 | src/violit/app_support.py:169 | `except Exception: pass` on stdout reconfigure |
| 113 | src/violit/background.py:214 | `except Exception: pass` on cancel handler |
| 118 | src/violit/background.py:238 | `except Exception: pass` on cancel handler |
| 122 | src/violit/background.py:252 | `except Exception: pass` on error handler |
| 132 | src/violit/cli.py:35 | `except KeyboardInterrupt: pass` on runpy |
| 173 | src/violit/state.py:426 | `except ValueError: pass` on list.remove |
| 183 | src/violit/widgets/chart_widgets.py:613 | `except Exception: pass` on rcParams |
| 187 | src/violit/widgets/chart_widgets.py:628 | `except Exception: pass` on font set |
| 189 | src/violit/widgets/chat_widgets.py:87 | `except Exception: pass` on setattr |
| 210 | src/violit/widgets/data_widgets.py:901 | `except Exception: pass` |
| 214 | src/violit/widgets/form_widgets.py:147 | `except ImportError: pass` on webview probe |
| 223 | src/violit/widgets/form_widgets.py:592 | `except ImportError: pass` on webview probe |
| 228 | src/violit/widgets/input_widgets.py:847 | `except (ValueError, TypeError): pass` on number parse |
| 237 | src/violit/widgets/layout_widgets.py:580 | `except Exception: pass` on render |
| 242 | src/violit/widgets/layout_widgets.py:586 | `except Exception: pass` on render |

### CWE-396 — Declaration of Catch for Generic Exception

| Finding | Source | Reason |
| --- | --- | --- |
| 18 | examples/04_llm_chat_application/demo_app_gemini_agent_highlevel.py:673 | generic `except Exception` in generator |
| 22 | examples/04_llm_chat_application/demo_app_gemini_agent_primitive.py:677 | generic `except Exception` in generator |
| 36 | src/violit/app.py:1537 | generic `except Exception` |
| 47 | src/violit/app_launcher.py:41 | generic `except Exception` on probe |
| 83 | src/violit/app_runtime.py:149 | generic `except Exception` on queue put |
| 101 | src/violit/app_support.py:169 | generic `except Exception` |
| 104 | src/violit/auth.py:67 | generic `except Exception` on bcrypt |
| 115 | src/violit/background.py:214 | generic `except Exception` |
| 136 | src/violit/cli.py:37 | generic `except Exception` |
| 144 | src/violit/db.py:120 | generic `except Exception` |
| 172 | src/violit/engine.py:75 | generic `except Exception` on send |
| 177 | src/violit/state.py:488 | generic `except Exception` on callbacks |
| 181 | src/violit/widgets/chart_widgets.py:29 | generic `except Exception` |
| 191 | src/violit/widgets/chat_widgets.py:87 | generic `except Exception` |
| 197 | src/violit/widgets/custom_widgets.py:61 | generic `except Exception` |
| 199 | src/violit/widgets/data_widgets.py:44 | generic `except Exception` on import |
| 219 | src/violit/widgets/form_widgets.py:183 | generic `except Exception` |
| 232 | src/violit/widgets/input_widgets.py:953 | generic `except Exception` |
| 239 | src/violit/widgets/layout_widgets.py:580 | generic `except Exception` |
| 244 | src/violit/widgets/media_widgets.py:32 | generic `except Exception` |

### CWE-390 — Detection of Error Condition Without Action

| Finding | Source | Reason |
| --- | --- | --- |
| 40 | src/violit/app.py:1872 | handler is only `pass` |
| 59 | src/violit/app_launcher.py:107 | handler is only `pass` |
| 94 | src/violit/app_support.py:35 | handler is only `pass` |
| 114 | src/violit/background.py:214 | handler is only `pass` |
| 133 | src/violit/cli.py:35 | handler is only `pass` |
| 174 | src/violit/state.py:426 | handler is only `pass` |
| 184 | src/violit/widgets/chart_widgets.py:613 | handler is only `pass` |
| 190 | src/violit/widgets/chat_widgets.py:87 | handler is only `pass` |
| 211 | src/violit/widgets/data_widgets.py:901 | handler is only `pass` |
| 215 | src/violit/widgets/form_widgets.py:147 | handler is only `pass` |
| 229 | src/violit/widgets/input_widgets.py:847 | handler is only `pass` |
| 238 | src/violit/widgets/layout_widgets.py:580 | handler is only `pass` |

### CWE-1071 — Empty Code Block

| Finding | Source | Reason |
| --- | --- | --- |
| 41 | src/violit/app.py:1872 | handler body only `pass` |
| 70 | src/violit/app_launcher.py:321 | handler body only `pass` |
| 95 | src/violit/app_support.py:35 | handler body only `pass` |
| 116 | src/violit/background.py:214 | handler body only `pass` |
| 134 | src/violit/cli.py:35 | handler body only `pass` |
| 175 | src/violit/state.py:426 | handler body only `pass` |
| 185 | src/violit/widgets/chart_widgets.py:613 | handler body only `pass` |
| 192 | src/violit/widgets/chat_widgets.py:87 | handler body only `pass` |
| 212 | src/violit/widgets/data_widgets.py:901 | handler body only `pass` |
| 216 | src/violit/widgets/form_widgets.py:147 | handler body only `pass` |
| 240 | src/violit/widgets/layout_widgets.py:580 | handler body only `pass` |

### BP-PY-46 — print Debugging In Library Code

| Finding | Source | Reason |
| --- | --- | --- |
| 7 | examples/02_violit_blog/violit_advanced_blog.py:349 | module-level print, no `__main__` guard |
| 10 | examples/02_violit_blog/violit_blog.py:299 | module-level print, no `__main__` guard |
| 31 | src/violit/app.py:177 | `print` in library code (debug block) |
| 33 | src/violit/app.py:179 | `print` in library code |
| 34 | src/violit/app.py:635 | `print` in `debug_print` helper |
| 48 | src/violit/app_launcher.py:51 | `print` in launcher method |
| 49 | src/violit/app_launcher.py:52 | `print` in launcher method |
| 50 | src/violit/app_launcher.py:53 | `print` in launcher method |
| 51 | src/violit/app_launcher.py:54 | `print` in launcher method |
| 52 | src/violit/app_launcher.py:57 | `print` in launcher method |
| 53 | src/violit/app_launcher.py:58 | `print` in launcher method |
| 56 | src/violit/app_launcher.py:90 | `print` in reload fallback |
| 57 | src/violit/app_launcher.py:91 | `print` in reload fallback |
| 76 | src/violit/app_launcher.py:449 | `print` in launcher method |
| 77 | src/violit/app_launcher.py:477 | `print` in launcher method |
| 80 | src/violit/app_launcher.py:494 | `print` in launcher method |
| 81 | src/violit/app_launcher.py:505 | `print` in launcher method |
| 87 | src/violit/app_runtime.py:418 | `print` in debug block |
| 88 | src/violit/app_runtime.py:419 | `print` in debug block |
| 89 | src/violit/app_runtime.py:420 | `print` in debug block |
| 96 | src/violit/app_support.py:52 | `print` in watcher |
| 97 | src/violit/app_support.py:58 | `print` in watcher |
| 102 | src/violit/app_support.py:183 | `print(splash)` in library code |
| 129 | src/violit/cli.py:23 | `print` to stderr in CLI library module |
| 137 | src/violit/cli.py:38 | `print` to stderr in CLI library module |
| 138 | src/violit/cli.py:51 | `print` in scaffold |
| 139 | src/violit/cli.py:124 | `print` in scaffold |
| 140 | src/violit/cli.py:125 | `print` in scaffold |
| 141 | src/violit/cli.py:126 | `print` in scaffold |
| 142 | src/violit/cli.py:127 | `print` in scaffold |
| 165 | src/violit/db.py:427 | `print` in migration helper |
| 166 | src/violit/db.py:446 | `print` in migration helper |
| 167 | src/violit/db.py:466 | `print` in migration helper |
| 168 | src/violit/db.py:490 | `print` in migration helper |
| 169 | src/violit/db.py:503 | `print` in migration helper |
| 217 | src/violit/widgets/form_widgets.py:175 | `print` in save handler |
| 220 | src/violit/widgets/form_widgets.py:184 | `print` in save handler |
| 233 | src/violit/widgets/input_widgets.py:954 | `print` in upload handler |

### BP-PY-47 — logging With String Format Before Logger

| Finding | Source | Reason |
| --- | --- | --- |
| 105 | src/violit/auth.py:176 | `logger.warning(f"...")` |
| 106 | src/violit/auth.py:225 | `logger.info(f"... {username}")` |
| 108 | src/violit/auth.py:329 | `logger.info(f"... {username}")` |
| 109 | src/violit/auth.py:357 | `logger.info(f"... id={uid}")` |
| 110 | src/violit/auth.py:390 | `logger.warning(f"...")` |
| 120 | src/violit/background.py:245 | `logger.error(f"... {e}")` |
| 123 | src/violit/background.py:312 | `logger.debug(f"... {len(dirty)} ...")` |
| 125 | src/violit/background.py:320 | `logger.debug(f"...")` |
| 126 | src/violit/background.py:323 | `logger.debug(f"...")` |
| 128 | src/violit/background.py:326 | `logger.debug(f"... {e}")` |
| 145 | src/violit/db.py:121 | `logger.warning(f"... {e}")` |
| 146 | src/violit/db.py:129 | `logger.info(f"... {table_name}")` |
| 149 | src/violit/db.py:152 | `logger.info(f"...")` |
| 151 | src/violit/db.py:170 | `logger.warning(f"...")` |
| 152 | src/violit/db.py:175 | `logger.warning(f"...")` |
| 153 | src/violit/db.py:187 | `logger.warning(f"...")` |
| 154 | src/violit/db.py:195 | `logger.warning(f"...")` |
| 156 | src/violit/db.py:202 | `logger.error(f"...")` |

### BP-PY-45 — sys.path Mutation At Runtime

| Finding | Source | Reason |
| --- | --- | --- |
| 25 | examples/99_github_issues/github_issue_22_reactive_slider_text_and_expander_title.py:4 | `sys.path.insert` for local src import |
| 26 | examples/99_github_issues/github_issue_24_function_vs_main_rendering.py:4 | `sys.path.insert` for local src import |
| 27 | examples/99_github_issues/github_issue_25_styling_over_multicolumn.py:4 | `sys.path.insert` for local src import |
| 62 | src/violit/app_launcher.py:174 | `sys.path.insert(0, reload_dir)` |
| 130 | src/violit/cli.py:30 | `sys.path.insert(0, script_dir)` |

### CWE-117 — Improper Output Neutralization for Logs

| Finding | Source | Reason |
| --- | --- | --- |
| 107 | src/violit/auth.py:225 | externally supplied `username` formatted into log message |

### CWE-215 — Insertion of Sensitive Information Into Debugging Code

| Finding | Source | Reason |
| --- | --- | --- |
| 32 | src/violit/app.py:177 | debug print exposes 20 chars of native token |

### CWE-208 — Observable Timing Discrepancy

| Finding | Source | Reason |
| --- | --- | --- |
| 85 | src/violit/app_runtime.py:279 | native auth token compared with `==` (not `compare_digest`) in access decision |

### CWE-89 — Improper Neutralization of Special Elements used in an SQL Command

*(no true positives; all 4 findings are false positives — 5, 8, 13, 148)*

### CWE-1121 — Excessive McCabe Cyclomatic Complexity

| Finding | Source | Reason |
| --- | --- | --- |
| 11 | examples/03_multi_db_sqlite_editor/demo_multi_db_sqlite_editor.py:171 | `validate_editor_change` — 15 real branches |
| 15 | examples/04_llm_chat_application/demo_app_gemini_agent_highlevel.py:577 | `reply` — 12 real branches |
| 19 | examples/04_llm_chat_application/demo_app_gemini_agent_primitive.py:578 | `reply` — 12 real branches |
| 54 | src/violit/app_launcher.py:60 | `_run_web_reload` — 25 real branches |
| 84 | src/violit/app_runtime.py:264 | `_setup_routes` — 82 real branches |
| 178 | src/violit/style_utils.py:386 | `split_utility_tokens` — 14 real branches |
| 195 | src/violit/widgets/chat_widgets.py:932 | `_apply_agent_event` — 26 real branches |
| 203 | src/violit/widgets/data_widgets.py:535 | `builder` — 15 real branches |
| 247 | src/violit/widgets/text_widgets.py:630 | `write` — 18 real branches |

### CWE-1124 — Excessively Deep Nesting

| Finding | Source | Reason |
| --- | --- | --- |
| 55 | src/violit/app_launcher.py:78 | assignment nested in try→while→if→if→for→if (6 levels) |
| 92 | src/violit/app_runtime.py:673 | `stream_components.insert` nested ≥6 levels in event handler |
| 150 | src/violit/db.py:165 | DDL build nested try→with→for→if→if→if→if (7 levels) |

### CWE-1046 — Creation of Immutable Text Using String Concatenation

| Finding | Source | Reason |
| --- | --- | --- |
| 170 | src/violit/engine.py:17 | `html += rendered[...] + ...` inside loop |
| 227 | src/violit/widgets/input_widgets.py:450 | `opts_html += f'<wa-radio ...>'` inside option loop |
| 248 | src/violit/widgets/text_widgets.py:742 | `text += str(cursor)` inside loop |

### PERF-PY-25 — Heavy Object Construction Per Homogeneous Element

| Finding | Source | Reason |
| --- | --- | --- |
| 1 | examples/01_demo_showcase/demo_showcase.py:521 | `on_click=lambda ...` constructed per loop element |
| 3 | examples/01_demo_showcase/old_archive_demo_showcase.py:346 | `on_click=lambda k=key: ...` per loop element |
| 6 | examples/02_violit_blog/violit_advanced_blog.py:128 | `lambda` per loop element (view handler) |
| 9 | examples/02_violit_blog/violit_blog.py:117 | `lambda` per loop element (view handler) |
| 23 | examples/05_orm/demo_orm_app.py:121 | `on_click=lambda todo_id=...` per loop element |
| 24 | examples/05_orm/demo_orm_app.py:123 | `on_click=lambda todo_id=...` per loop element |

### PERF-PY-23 — Polymorphic Serialize Or Encode Inside Hot Loop

| Finding | Source | Reason |
| --- | --- | --- |
| 16 | examples/04_llm_chat_application/demo_app_gemini_agent_highlevel.py:618 | `json.dumps(action_input)` inside loop over actions |
| 20 | examples/04_llm_chat_application/demo_app_gemini_agent_primitive.py:622 | `json.dumps(action_input)` inside loop over actions |

### CWE-1084 — Invokable Control Element with Excessive File or Data Access Operations

| Finding | Source | Reason |
| --- | --- | --- |
| 4 | examples/02_violit_blog/violit_advanced_blog.py:15 | `init_db()` performs ≥3 `execute`/`open` operations |

### CWE-328 — Use of Weak Hash

| Finding | Source | Reason |
| --- | --- | --- |
| 205 | src/violit/widgets/data_widgets.py:637 | `hashlib.sha1(...)` used for grid structure hash |

### CWE-367 — Time-of-check Time-of-use (TOCTOU) Race Condition

| Finding | Source | Reason |
| --- | --- | --- |
| 221 | src/violit/widgets/form_widgets.py:543 | `os.path.exists(directory)` check before `os.makedirs` |


### CWE-829 — Inclusion of Functionality from Untrusted Control Sphere (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 131 | `src/violit/cli.py:34` | `runpy.run_path(script_path, run_name="__main__")` loads a dynamically selected file path for execution — rule condition met (reclassified from Uncertain) |

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/violit/chunks`
- Function evidence: `scripts/violit/findings/functions`
- Validation: `git diff --check` — `pass`

## Post-fix remaining-FP audit (2026-08-02)

Mode A — remaining false positives. Fresh scan run with the post-fix binary (b5b8fde, rebuilt 2026-08-02 16:29); the fresh finding IDs (1–224) do not correspond to the old IDs, so every fresh finding was matched to the old audit by `Source:` path (file:line:col).

### Run metadata

```yaml
timestamp: 2026-08-02T16:38:51+05:30
repository: violit
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit
branch: main
commit: 8fae080f49f374b062172ed6ac71042539ad1f7a
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit
chunk_path: scripts/violit/chunks
function_context_path: scripts/violit/findings/functions
```

### Scan evidence

- Build command: `make build` (`go build -o bin/goslop ./cmd/goslop`), binary rebuilt at commit `b5b8fde` (2026-08-02 16:29)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/violit/chunks -context-dir scripts/violit/findings/functions real-repos/violit`
- Findings: `224` (was `248` pre-fix; repo commit unchanged)
- Chunks reviewed: `scripts/violit/chunks/Chunk_1_25.txt` … `scripts/violit/chunks/Chunk_201_224.txt` (all 9 chunk files)
- Function contexts reviewed: `scripts/violit/findings/functions/2.txt`, `8.txt`, `9.txt`, `10.txt`, `21.txt`, `22.txt`, `58.txt`, `77.txt`, `81.txt`, `101.txt`, `114.txt`, `130.txt`, `139.txt`, `140.txt`, `141.txt`, `156.txt`, `190.txt`, `202.txt`, `211.txt`, `52.txt`
- Source verification: re-read the enclosing source for every candidate FP at the flagged line and confirmed the construct is unchanged from the pre-fix audit (same line numbers)
- Observations: 10 audited FPs no longer fire (old 5, 8, 28, 90, 148, 157, 159, 160, 161, 164 — CWE-89 db_query helper/DDL, BP-PY-46 string-literal print, CWE-215 bool print, BP-PY-12 `session.exec`), i.e. fixed by the reducer. 14 audited TPs no longer fire (old 7, 10, 25, 26, 27, 32, 129, 137, 138, 139, 140, 141, 142, 221 — module-level example prints, `sys.path.insert` bootstrap in examples/cli, CWE-215 native-token debug print at app.py:177, cli.py prints, CWE-367 TOCTOU at form_widgets.py:543); every one of these constructs was verified still present in source, i.e. suppressed-but-present (not fixed) — candidates for Mode-B review.

### Audit checklist

- [x] Read every assigned chunk under `scripts/violit/chunks` (all 9 files).
- [x] Read `scripts/violit/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

### Classification summary (fresh run)

Fresh findings were matched to the old audit by `Source:` path. Every fresh finding whose source matches an audited TP was classified TP (including fresh 52, which matches audited TP 61 by source: `app_launcher.py:159:17` — the `print("INFO:     Stopping reloader process")` line, construct verified at line 159). Every fresh finding matching an audited FP source was re-classified FP after re-verifying the construct against the rule condition. No genuinely new findings (no fresh finding without an old-audit source match).

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 19 | 2, 8, 9, 10, 21, 22, 58, 77, 81, 101, 114, 130, 139, 140, 141, 156, 190, 202, 211 |
| True positive | 205 | 1, 3, 4, 5, 6, 7, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 78, 79, 80, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 131, 132, 133, 134, 135, 136, 137, 138, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 191, 192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 203, 204, 205, 206, 207, 208, 209, 210, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224 |
| Uncertain | 0 | — |

### False positives (remaining)

All 19 remaining false positives re-appear at the exact `Source:` locations audited as false positives in the pre-fix run; the flagged constructs were re-verified in the current source (unchanged) and each still fails to satisfy its rule condition.

### [ ] Finding `2` — `CWE-1121`

- Function context: `scripts/violit/findings/functions/2.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/examples/01_demo_showcase/old_archive_demo_showcase.py:178:38`
- Checklist pattern: rule condition “function has at least twelve visible control-flow branches”

Source excerpt:

```
def _submit_chat_prompt(prompt: str):
    if isinstance(prompt, dict):
        prompt_payload = {
            "text": str(prompt.get("text") or "").strip(),
            "files": [entry for entry in list(prompt.get("files") or []) if entry is not None],
            ...
```

Why this is a false positive: identical to audited FP 2 — at most 9 real control-flow statements; the substring counter reaches 12 only by counting `if `/`for ` inside list comprehensions, which are expressions, not branches.

Checklist evidence: re-counted the function body from the current source; real control-flow statements < 12.

### [ ] Finding `8` and `9` — `BP-PY-37` / `CWE-89`

- Function context: `scripts/violit/findings/functions/8.txt`, `9.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/examples/03_multi_db_sqlite_editor/demo_multi_db_sqlite_editor.py:230:17` (same construct, two rules)
- Checklist pattern: rule conditions “DB-API execute builds SQL with f-string/%; use bound parameters” and “dynamic SQL string reaches execute/executemany”

Source excerpt:

```
    if field not in HR_FIELDS | PROJECT_FIELDS:
        raise ValueError(f"Unsupported field: {field}")
    ...
        conn.execute(f"UPDATE employees SET {field} = ? WHERE emp_id = ?", (stored_value, emp_id))
```

Why this is a false positive: identical to audited FPs 12/13 — the only interpolated token is the column identifier, constrained by the same-function allowlist check to `HR_FIELDS | PROJECT_FIELDS`; all values travel via `?` bound parameters.

Checklist evidence: field is allowlist-constrained at demo_multi_db_sqlite_editor.py:223-224 (verified in current source); values are bound parameters.

### [ ] Finding `10` — `BP-PY-37`

- Function context: `scripts/violit/findings/functions/10.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/examples/03_multi_db_sqlite_editor/demo_multi_db_sqlite_editor.py:234:13`
- Checklist pattern: rule condition “DB-API `cursor.execute` builds SQL with f-string or % format; use bound parameters”

Source excerpt:

```
    with _connect(PROJECT_DB) as conn:
        conn.execute(f"UPDATE assignments SET {field} = ? WHERE emp_id = ?", (value, emp_id))
```

Why this is a false positive: identical to audited FP 14 — second statement in `persist_editor_change`, guarded by the same `field not in HR_FIELDS | PROJECT_FIELDS` allowlist check; values are bound with `?` placeholders.

Checklist evidence: identifier is allowlist-constrained in the same function; values bound.

### [ ] Finding `21` — `CWE-93`

- Function context: `scripts/violit/findings/functions/21.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/app.py:85:21`
- Checklist pattern: rule condition “externally influenced value is written into an HTTP response header without CRLF neutralization”

Source excerpt:

```
    def file_response(self, full_path, stat_result, scope, status_code=200):
        response = super().file_response(full_path, stat_result, scope, status_code)
        if not response.headers.get("Cache-Control"):
            response.headers["Cache-Control"] = self._cache_control_for_scope(scope)
```

Why this is a false positive: identical to audited FP 29 — the header value is the return of `_cache_control_for_scope`, whose every branch returns one of four fixed string literals; no request data can reach the value.

Checklist evidence: verified `_cache_control_for_scope` at app.py:89 ff. returns only hardcoded literal strings; the RHS being a call expression is what triggers the “dynamic” classification.

### [ ] Finding `22` — `CWE-1121`

- Function context: `scripts/violit/findings/functions/22.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/app.py:124:434`
- Checklist pattern: rule condition “function has at least twelve visible control-flow branches”

Source excerpt:

```
    def __init__(self, mode: Optional[str] = None, title="Violit App", ...):
        self._mode_is_explicit = mode is not None
        self.mode = (mode or 'ws').strip().lower()
        ...
        if db is not None:
            from .db import ViolItDB, normalize_db_url
            self.db = ViolItDB(normalize_db_url(db), migrate=migrate)
        ...
        self.debug_mode = '--debug' in sys.argv
```

Why this is a false positive: identical to audited FP 30 — the constructor contains only 7 real `if` statements (AST count over lines 124–494); the substring counter reaches 30 `if ` hits by counting ternary/`or` expressions and `if` text inside string literals.

Checklist evidence: AST count of `If`/`For`/`While`/`Try` statements directly in `App.__init__` = 7 < 12.

### [ ] Finding `58` — `BP-PY-40`

- Function context: `scripts/violit/findings/functions/58.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/app_launcher.py:298:15`
- Checklist pattern: rule condition “threading.Thread started without .join … avoid fire-and-forget non-daemon threads” (detector explicitly skips lines with `daemon=True`)

Source excerpt:

```
        thread = threading.Thread(target=server_manager, daemon=True)
        thread.start()
```

Why this is a false positive: identical to audited FP 67 — the thread is constructed with `daemon=True` on the preceding line, but the exemption only inspects the `.start(` line, so the same-line check misses it.

Checklist evidence: verified `threading.Thread(target=server_manager, daemon=True)` at app_launcher.py:297 immediately before the flagged `.start()` at line 298.

### [ ] Finding `77` — `BP-PY-13`

- Function context: `scripts/violit/findings/functions/77.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/app_runtime.py:414:32`
- Checklist pattern: rule condition “A secret-like name is assigned a non-empty string literal in source”

Source excerpt:

```
            csrf_token = self._generate_csrf_token(sid) if sid and self.csrf_enabled else ""
            csrf_script = f'<script>window._csrf_token = "{csrf_token}";</script>' if csrf_token else ""
```

Why this is a false positive: identical to audited FP 86 — the matched `window._csrf_token = ` is JavaScript text inside an HTML/JS string template, not a Python assignment; the value is a runtime-generated CSRF token, not a hardcoded secret.

Checklist evidence: the matched construct is a JS variable name inside a string literal with a runtime-generated f-string value.

### [ ] Finding `81` — `BP-PY-32`

- Function context: `scripts/violit/findings/functions/81.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/app_runtime.py:553:24`
- Checklist pattern: rule condition “FileResponse / static file helpers use a path from user input without confinement”

Source excerpt:

```
                media_path = str(payload.get("path") or "").strip()
                if not media_path or not os.path.exists(media_path) or not os.path.isfile(media_path):
                    return Response(status_code=404)
                ...
                return FileResponse(
                    path=media_path, ...
```

Why this is a false positive: identical to audited FP 91 — the path comes from the application's own session-store media registry (`_vl_media_sources`), populated by the framework for developer-created media widgets; only `media_id` is user-controlled and it selects among registered entries.

Checklist evidence: verified `media_path` is read from the app-populated `registry` (app_runtime.py:548-550), not derived from request parameters.

### [ ] Finding `101` — `CWE-1121`

- Function context: `scripts/violit/findings/functions/101.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/background.py:188:52`
- Checklist pattern: rule condition “function has at least twelve visible control-flow branches”

Source excerpt:

```
    def _run(self, sid: str, current_view_id: str):
        t = session_ctx.set(sid) if sid else None
        view_token = view_ctx.set(current_view_id) if current_view_id else None
        ...
        try:
            self._result = self._fn()
            if self._cancel_event.is_set():
                ...
```

Why this is a false positive: identical to audited FP 111 — 11 real control-flow statements (AST count excluding nested definitions); the counter reaches ≥12 by counting ternary `if` expressions as branches.

Checklist evidence: AST count of `If`/`For`/`While`/`Try` in `BackgroundTask._run` = 11 < 12.

### [ ] Finding `114` — `CWE-117`

- Function context: `scripts/violit/findings/functions/114.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/background.py:312:17`
- Checklist pattern: rule condition “externally influenced input … without neutralizing line-breaking control characters”

Source excerpt:

```
                logger.debug(f"[background] Pushed {len(dirty)} updates to session {sid[:8]}...")
```

Why this is a false positive: identical to audited FP 124 — the interpolated values are an integer and a truncated server-generated session-id prefix; no externally controlled string is written to the log.

Checklist evidence: the formatted values are `len(dirty)` (int) and `sid[:8]` (truncated internally generated id).

### [ ] Finding `130` — `CWE-117`

- Function context: `scripts/violit/findings/functions/130.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/db.py:129:17`
- Checklist pattern: rule condition “externally influenced input … without neutralizing line-breaking control characters”

Source excerpt:

```
            if table_name not in existing_tables:
                logger.info(f"[violit:db] ✅ New table created: {table_name}")
                continue
```

Why this is a false positive: identical to audited FP 147 — `table_name` iterates SQLModel model metadata (developer-defined class names), not external input; it cannot carry CR/LF from an attacker.

Checklist evidence: the interpolated value derives from `SQLModel.metadata.tables`, not from external input.

### [ ] Finding `139` — `BP-PY-14`

- Function context: `scripts/violit/findings/functions/139.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/db.py:279:20`
- Checklist pattern: rule condition “`requests` HTTP calls omit `timeout=`” (detector regex `\b(session|sess|req_session)\.(get|post|…)\s*\(`)

Source excerpt:

```
        _check_sqlmodel()
        with Session(self._engine) as session:
            return session.get(model, pk)
```

Why this is a false positive: identical to audited FP 158 — `session` is a SQLModel `Session` and `session.get(model, pk)` is an ORM primary-key lookup, not an HTTP request; the file imports no `requests`.

Checklist evidence: `Session` is imported from `sqlmodel` and wraps `self._engine`; verified in current db.py:277-279.

### [ ] Finding `140` — `BP-PY-14`

- Function context: `scripts/violit/findings/functions/140.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/db.py:353:13`
- Checklist pattern: rule condition “`requests` HTTP calls omit `timeout=`” (detector regex matches `session.delete(` by name)

Source excerpt:

```
            managed = session.merge(obj)
            session.delete(managed)
            session.commit()
```

Why this is a false positive: identical to audited FP 162 — ORM delete on a SQLModel database session, not a `requests` call.

Checklist evidence: `session` is a SQLModel `Session` bound to `self._engine`; `delete`/`commit` are ORM methods.

### [ ] Finding `141` — `BP-PY-14`

- Function context: `scripts/violit/findings/functions/141.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/db.py:369:17`
- Checklist pattern: rule condition “`requests` HTTP calls omit `timeout=`” (detector regex matches `session.delete(` by name)

Source excerpt:

```
            for row in rows:
                managed = session.merge(row)
                session.delete(managed)
            session.commit()
```

Why this is a false positive: identical to audited FP 163 — same ORM delete construct on a SQLModel database session.

Checklist evidence: `session` is a SQLModel `Session`; `delete`/`commit` are ORM methods.

### [ ] Finding `156` — `CWE-1121`

- Function context: `scripts/violit/findings/functions/156.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/widgets/chart_widgets.py:19:27`
- Checklist pattern: rule condition “function has at least twelve visible control-flow branches”

Source excerpt:

```
    def _normalize(value):
        if isinstance(value, np.ndarray):
            return [_normalize(item) for item in value.tolist()]
        if isinstance(value, np.generic):
            return value.item()
        if isinstance(value, dict):
            ...
        if isinstance(value, list):
            return [_normalize(item) for item in value]
        if isinstance(value, tuple):
            return [_normalize(item) for item in value]
        return value
```

Why this is a false positive: identical to audited FP 179 — 6 `if` statements + 1 `except` (7 real branches); the count reaches 12 only via `for `/`if ` inside list/dict comprehensions.

Checklist evidence: AST count of `If`/`For`/`While`/`Try` statements = 7 < 12.

### [ ] Finding `190` — `CWE-1121`

- Function context: `scripts/violit/findings/functions/190.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/widgets/form_widgets.py:67:23`
- Checklist pattern: rule condition “function has at least twelve visible control-flow branches”

Source excerpt:

```
        def builder():
            token = rendering_ctx.set(cid)
            bt = text() if callable(text) else text
            ...
            icon_html = f'...' if icon and not any(ord(c) > 127 for c in str(icon)) else ''
            ...
            if use_container_width:
                host_style = merge_style(host_style, "width:100%;")
            if height not in (None, "", "auto"):
                ...
```

Why this is a false positive: identical to audited FP 213 — 6 real control-flow statements in the `builder()` closure; the count reaches 12 via ternary `... if ... else ...` expressions and `for` inside generator expressions.

Checklist evidence: AST count of real control-flow statements in `builder()` at line 67 = 6 < 12.

### [ ] Finding `202` — `CWE-1121`

- Function context: `scripts/violit/findings/functions/202.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/widgets/input_widgets.py:415:23`
- Checklist pattern: rule condition “function has at least twelve visible control-flow branches”

Source excerpt:

```
        def builder():
            token = rendering_ctx.set(cid)
            cv = s.value
            ...
            opts_html = ""
            for i_opt, opt in enumerate(options):
                sel = 'checked' if opt == cv else ''
                escaped_opt = html_lib.escape(str(opt), quote=True)
                ...
                if captions and i_opt < len(captions) and captions[i_opt]:
                    ...
```

Why this is a false positive: identical to audited FP 226 — 1 real `for` loop + 6 `if` statements (7 real branches); the count reaches 12 via ternary `if` expressions.

Checklist evidence: AST count of real control-flow statements in the `builder()` at line 415 = 7 < 12.

### [ ] Finding `211` — `CWE-1121`

- Function context: `scripts/violit/findings/functions/211.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit/src/violit/widgets/layout_widgets.py:326:107`
- Checklist pattern: rule condition “function has at least twelve visible control-flow branches”

Source excerpt:

```
    def expander(self, label, expanded=False, icon=None, cls: str = "", style: str = "", key: Any = None):
        cid = self._resolve_widget_cid("expander", key)
        summary_cid = f"{cid}__summary"
        details_cid = f"{cid}__details"

        class ExpanderContext:
            def __init__(self, app, expander_id, label, expanded, icon=None, user_cls="", user_style=""):
                ...
            def __enter__(self):
                def summary_builder():
                    ...
```

Why this is a false positive: identical to audited FP 235 — `expander` contains only 8 real control-flow statements; the counter reaches 12 by including branches of the nested `ExpanderContext` class and `summary_builder` definitions, which are not part of `expander`'s own flow.

Checklist evidence: AST count of real control-flow statements in `expander` (excluding nested definitions) = 8 < 12.

### Uncertain findings

None — every fresh finding was matched to an audited TP or an audited FP by exact `Source:` path, and the one finding not listed in the old tables (fresh 52 = `app_launcher.py:159:17`, BP-PY-46) matches audited TP 61 by source with the construct verified present.

### Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/violit/chunks` (fresh scan, 9 chunk files)
- Function evidence: `scripts/violit/findings/functions` (fresh scan)
- Validation: `git diff --check` — `pass`

## Post-fix v2 audit (latest binary)

### Run metadata

```yaml
timestamp: 2026-08-02T18:00:00+05:30 (latest binary rebuild ~17:56)
repository: violit
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit
branch: main
commit: 8fae080f49f374b062172ed6ac71042539ad1f7a (unchanged across all audits)
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/violit
chunk_path: scripts/violit/chunks (10 files, fresh 17:58)
function_context_path: scripts/violit/findings/functions (fresh 17:58)
```

### Scan evidence

- Build command: `make build` (`go build -o bin/goslop ./cmd/goslop`), binary rebuilt ~2026-08-02 17:56
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/violit/chunks -context-dir scripts/violit/findings/functions real-repos/violit`
- Findings: `229` (was `224` in Mode A; repo commit unchanged)
- Chunks reviewed: `scripts/violit/chunks/Chunk_1_25.txt` … `Chunk_226_229.txt` (all 10)
- Function contexts reviewed: `scripts/violit/findings/functions/` for every fresh finding matched to a prior audited FP (2, 10, 11, 12, 26, 27, 63, 82, 86, 106, 119, 135, 144, 145, 146, 161, 195, 207, 216) plus spot-checks on returned TPs
- Matching: every fresh finding was matched to the prior audits by `Source:` (file:line[:col]); a source path is present in at least one prior classification for all 229 fresh findings

### Classification summary (fresh run)

All 229 fresh findings matched an audited classification by exact `Source:`; the flagged constructs were re-verified in the current source (line numbers unchanged) and re-checked against the rule condition. The 19 fresh findings matching audited FPs remain false positives; the 210 matching audited TPs remain true positives.

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 19 | 2, 10, 11, 12, 26, 27, 63, 82, 86, 106, 119, 135, 144, 145, 146, 161, 195, 207, 216 |
| True positive | 210 | all other fresh finding IDs (1–229) |
| Uncertain | 0 | — |

Observations vs Mode A:
- 5 Mode-A-suppressed TPs are firing again in the latest binary (constructs unchanged in source): old 7/10 (`print` at module level in `violit_advanced_blog.py:349` / `violit_blog.py:299`) and old 25/26/27 (`sys.path.insert` bootstrap in the three `examples/99_github_issues/*.py:4`) — re-classified TP.
- Still suppressed (construct present, no finding): old 32 (CWE-215 at `app.py:177`), old 129/137–142 (`cli.py` prints), old 221 (CWE-367 at `form_widgets.py:543`).
- All 10 Mode-A-fixed FPs stay fixed (no fresh finding at old 5/8/28/90/148/157/159/160/161/164 sources).
- All 19 Mode-A remaining FPs re-fire at identical `Source:` locations; each was re-verified against its rule condition and remains a false positive.

### False positives (remaining, 19)

Every one re-appears at the exact `Source:` audited as FP in the pre-fix run and re-verified unchanged; the reasons are identical to the prior audit entries (old 2, 12/13/14, 29, 30, 67, 86, 91, 111, 124, 147, 158/162/163, 179, 213, 226, 235).

## Fix checklist (FP patterns)

| Pattern # | Rule | Trigger shape | Count | Example sources |
| --- | --- | --- | ---: | --- |
| 1 | CWE-1121 | substring `if `/`for `/`while `/`except ` counter treats ternaries (`x if y else z`), comprehension/generator clauses (`for item in …`, `if …` inside `[...]`/`{...}`/`any(...)`), and nested def/class bodies as branches; real AST branch count < 12 | 7 | old_archive_demo_showcase.py:178, app.py:124 (`__init__`), background.py:188 (`_run`), chart_widgets.py:19 (`_normalize`), form_widgets.py:67, input_widgets.py:415, layout_widgets.py:326 (`builder` closures) |
| 2 | BP-PY-37 + CWE-89 | `conn.execute(f"UPDATE … SET {field} = ? …", (v, id))` where the interpolated token `field` is allowlist-constrained by a same-function `if field not in HR_FIELDS \| PROJECT_FIELDS: raise`; values travel via bound params | 3 | demo_multi_db_sqlite_editor.py:230:17 (×2 rules), :234:13 |
| 3 | CWE-93 | header assignment `response.headers["Cache-Control"] = self._cache_control_for_scope(scope)` where the callee returns only fixed string literals (trigger is the call expression, not external data) | 1 | src/violit/app.py:85:21 |
| 4 | BP-PY-40 | `thread = threading.Thread(…, daemon=True)` on line N, `thread.start()` on line N+1; the `daemon=True` exemption inspects only the `.start(` line | 1 | app_launcher.py:298:15 (constructor at :297) |
| 5 | BP-PY-13 | regex matches `window._csrf_token = ` — JavaScript identifier text *inside* an f-string/HTML template, whose RHS is a runtime-generated value, not a hardcoded literal | 1 | app_runtime.py:414:32 |
| 6 | BP-PY-32 | `FileResponse(path=media_path)` where `media_path` comes from the app-populated session/media registry (`registry.get(media_id)`), not from request params | 1 | app_runtime.py:553:24 |
| 7 | CWE-117 | log f-strings interpolating only internally generated values — `len(dirty)` (int), `sid[:8]` (truncated server id), `table_name` (SQLModel metadata) — which cannot carry CR/LF | 2 | background.py:312:17, db.py:129:17 |
| 8 | BP-PY-14 | `session.get(model, pk)` / `session.delete(managed)` where `session` is a SQLModel `Session` (imported from `sqlmodel`, bound to `self._engine`), not a `requests.Session`; name-based regex match | 3 | db.py:279:20, :353:13, :369:17 |

Condition to distinguish safe from vulnerable (per pattern): 1 — count AST `If/For/While/Try` statements of the function body only, excluding ternaries, comprehensions, and nested defs; 2 — require the interpolated SQL token to pass an in-scope allowlist membership check before the call (else flag); 3 — only flag header writes whose value can actually contain externally influenced data (resolve callee literal-only returns as constant); 4 — resolve the `Thread(...)` constructor (previous line / AST) for the `daemon=True` exemption; 5 — skip matches whose text lies inside a string literal and whose RHS is not a literal; 6 — only flag FileResponse paths derived from request input, not framework registry lookups; 7 — only flag log f-strings that interpolate externally controlled string values; 8 — resolve the receiver type/import (SQLModel `Session` = miss), flag only `requests` sessions.

## New findings

None — every fresh finding (229/229) matched a prior audited classification by exact `Source:`; no fresh finding lacked a prior classification.

### Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/violit/chunks` (fresh scan, 10 chunk files)
- Function evidence: `scripts/violit/findings/functions` (fresh scan)
- Validation: `git diff --check` — `pass`
