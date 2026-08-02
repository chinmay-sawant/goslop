# False-positive audit — onlymaps

## Run metadata

```yaml
timestamp: 2026-08-02T07:36:17Z
repository: onlymaps
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps
branch: main
commit: 6444a59
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps
chunk_path: scripts/onlymaps/chunks
function_context_path: scripts/onlymaps/findings/functions
```

## Scan evidence

- Build command: `(pre-built) ./bin/goslop`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/onlymaps/chunks -context-dir scripts/onlymaps/findings/functions real-repos/onlymaps`
- Findings: `99`
- Chunks reviewed: `scripts/onlymaps/chunks/Chunk_1_25.txt`, `scripts/onlymaps/chunks/Chunk_26_50.txt`, `scripts/onlymaps/chunks/Chunk_51_75.txt`, `scripts/onlymaps/chunks/Chunk_76_99.txt`
- Function contexts reviewed: `scripts/onlymaps/findings/functions/1.txt` … `50.txt` (all exported contexts read; contexts for findings 51–99 taken from the chunk files and the enclosing source read directly for every finding)

## Audit checklist

- [x] Read every assigned chunk under `scripts/onlymaps/chunks`.
- [x] Read `scripts/onlymaps/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 63 | 2, 3, 4, 5, 6, 7, 8, 9, 15, 19, 20, 21, 22, 23, 24, 25, 26, 27, 32, 33, 34, 35, 36, 37, 38, 39, 41, 43, 44, 45, 46, 47, 48, 49, 50, 55, 56, 57, 58, 64, 65, 66, 67, 69, 70, 71, 72, 73, 74, 75, 83, 84, 87, 88, 89, 90, 92, 93, 94, 95, 96, 97, 98 |
| True positive | 36 | 1, 10, 11, 12, 13, 14, 16, 17, 18, 28, 29, 30, 31, 40, 42, 51, 52, 53, 54, 59, 60, 61, 62, 63, 68, 76, 77, 78, 79, 80, 81, 82, 85, 86, 91, 99 |
| Uncertain | 0 | — |

## False positives

### [ ] Finding `2` — `BP-PY-7`

- Function context: `scripts/onlymaps/findings/functions/2.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_connection.py:136:13`
- Checklist pattern: `.open(` is a call on a custom lifecycle method, not the builtin file-opening `open()`

Source excerpt:

```
    def __enter__(self) -> Self:  # <async>
        """
        Opens a connection to the database and returns
        the instance itself.
        """
        self.open()  # <await>
        return self
```

Why this is a false positive: `self.open()` invokes `Connection.open` — a database-connection lifecycle method defined at `_connection.py:382` — from inside `__enter__`; no file handle is created, so no `with`-managed file resource exists, and the instance is itself used under `with`.

Checklist evidence: BP-PY-7's condition (a file opened without `with`, risking resource leaks) is unmet — the flagged call opens a database connection, not a file.

### [ ] Finding `3` — `BP-PY-1`

- Function context: `scripts/onlymaps/findings/functions/3.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_connection.py:215:1`
- Checklist pattern: handler's failure is immediately re-raised, so it is not swallowed

Source excerpt:

```
            try:
                if not self.__in_transaction:
                    self.__driver.init_transaction(self.__conn)
                yield cursor
                is_query_successful = True
            except:
                is_query_successful = False
                if not self.__in_transaction:
                    self.__conn.rollback()  # <await>
                raise
```

Why this is a false positive: the bare `except:` only rolls back and then re-raises (`raise`), so the exception propagates to the caller; nothing is swallowed.

Checklist evidence: the shown source's suite ends in `raise`, so the "swallows all exceptions" predicate of BP-PY-1 is not satisfied.

### [ ] Finding `4` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/4.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_connection.py:232:9`
- Checklist pattern: `exec` token is a method *definition*, not a call of the builtin `exec()`

Source excerpt:

```
    @require_open
    @__require_not_iter_same_ctx
    def exec(self, sql: str, /, *args: Any, **kwargs: Any) -> None:  # <async>
```

Why this is a false positive: the flagged token is the declaration of `Connection.exec`, the library's SQL-execution method (DB-API style naming); no Python `eval`/`exec`/`compile` builtin is called anywhere on the line.

Checklist evidence: BP-PY-12's condition "eval/exec on dynamic input" requires an actual `exec(` call with a non-literal argument; the line is a `def` signature.

### [ ] Finding `5` — `CWE-94`

- Function context: `scripts/onlymaps/findings/functions/5.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_connection.py:232:9`
- Checklist pattern: `exec` token is a method definition; no code-generation sink exists

Source excerpt:

```
    @require_open
    @__require_not_iter_same_ctx
    def exec(self, sql: str, /, *args: Any, **kwargs: Any) -> None:  # <async>
```

Why this is a false positive: CWE-94's condition is a code-generation/dynamic-import sink (`eval`/`exec`/`compile`/`__import__`) reached by a non-literal expression; the shown line declares a method named `exec` and never invokes a code-generation API.

Checklist evidence: the "dynamic value reaches a Python code-generation sink" condition is unmet — there is no call, and the first "argument" (`self, sql: str`) is a parameter list of a definition, not a dynamic expression.

### [ ] Finding `6` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/6.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_connection.py:244:22`
- Checklist pattern: `exec` token is a method call on the library's own `Query` object, not the builtin

Source excerpt:

```
        self.__query.exec(sql, args, kwargs)  # <await>
```

Why this is a false positive: `self.__query.exec(...)` calls `Query.exec` (defined at `_query.py:52`), the library's SQL-query runner; Python's `exec()` builtin is never invoked.

Checklist evidence: BP-PY-12's condition requires the `eval`/`exec` builtin; the call target is a user-defined method on `self.__query`.

### [ ] Finding `7` — `BP-PY-1`

- Function context: `scripts/onlymaps/findings/functions/7.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_connection.py:376:1`
- Checklist pattern: handler rolls back and re-raises; failure is propagated

Source excerpt:

```
            try:
                self.__transaction_id = self.__context_id
                self.__driver.init_transaction(self.__conn)
                yield
                self.__conn.commit()  # <await>
            except:
                self.__conn.rollback()  # <await>
                raise
```

Why this is a false positive: the bare `except:` performs the transaction rollback and re-raises, so the failure is neither swallowed nor hidden.

Checklist evidence: the suite ends with `raise` — the "bare except swallows all exceptions" condition is not satisfied.

### [ ] Finding `8` — `BP-PY-7`

- Function context: `scripts/onlymaps/findings/functions/8.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_connection.py:382:9`
- Checklist pattern: `def open(self)` is a method definition, not an `open()` call

Source excerpt:

```
    def open(self) -> None:  # <async>
        """
        Establishes a connection to the database.
```

Why this is a false positive: the flagged line declares the `Connection.open` lifecycle method; no file is opened on the line, so there is no file handle to leak.

Checklist evidence: BP-PY-7's condition "a file is opened without a `with` statement" — the line is a method declaration, not a file-open call.

### [ ] Finding `9` — `BP-PY-1`

- Function context: `scripts/onlymaps/findings/functions/9.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_connection.py:408:1`
- Checklist pattern: handler clears state and re-raises

Source excerpt:

```
            try:
                self.__conn = self.__driver.init_connection(
                    self.__conn_factory()  # <await>
                )
                self.__is_open = True
            except:
                self.__is_open = False
                raise
```

Why this is a false positive: the bare `except:` resets `__is_open` and re-raises; the failure propagates to the caller and is not hidden.

Checklist evidence: suite ends with `raise` — the "swallows all exceptions" predicate of BP-PY-1 is not met.

### [ ] Finding `15` — `BP-PY-7`

- Function context: `scripts/onlymaps/findings/functions/15.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_connection.py:464:25`
- Checklist pattern: `.open(` is the custom connection lifecycle method, not file open

Source excerpt:

```
                finally:
                    self.open()  # <await>
                    test_connection = False
```

Why this is a false positive: `self.open()` re-establishes the database connection (method defined at `_connection.py:382`); no file resource exists.

Checklist evidence: BP-PY-7 targets the builtin file-opening `open()`; the call is `self.open()` on the connection object.

### [ ] Finding `19` — `CWE-478`

- Function context: `scripts/onlymaps/findings/functions/19.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_drivers.py:403:1`
- Checklist pattern: unmatched values are handled by an unconditional fallback statement immediately after the match

Source excerpt:

```
            match value:
                # Force `BINARY_DOUBLE` TYPE even if the decimal part is zero.
                case float():
                    return cursor.var(
                        oracledb.DB_TYPE_BINARY_DOUBLE, arraysize=num_elements
                    )
                # Force `TIMESTAMP` type if datetime has microseconds.
                case datetime() if value.microsecond > 0:
                    return cursor.var(
                        oracledb.DB_TYPE_TIMESTAMP, arraysize=num_elements
                    )
            return None
```

Why this is a false positive: the missing-default concern is addressed — every value outside the two cases falls through to the unconditional `return None`, which is exactly the behavior a `case _:` would provide.

Checklist evidence: the rule condition "match without a default branch" is met syntactically, but the source's immediately following `return None` safely handles values outside the explicitly supported cases.

### [ ] Finding `20` — `BP-PY-7`

- Function context: `scripts/onlymaps/findings/functions/20.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_pool.py:169:13`
- Checklist pattern: `.open(` is a custom pool lifecycle method inside `__enter__`

Source excerpt:

```
        self.open()  # <await>
        return self
```

Why this is a false positive: `self.open()` invokes the pool's own connection-opening method (defined at `_pool.py:429`) from `__enter__`; no file handle is involved.

Checklist evidence: the call target is a user-defined method, not the builtin `open()`; BP-PY-7's file-resource condition is unmet.

### [ ] Finding `21` — `BP-PY-7`

- Function context: `scripts/onlymaps/findings/functions/21.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_pool.py:192:13`
- Checklist pattern: `.open(` is a call on a `Connection` object, not file open

Source excerpt:

```
        conn = Connection(
            self.__conn_factory,
            self.__driver,
            handle_broken_conn=True,
        )
        conn.open()  # <await>
        return conn
```

Why this is a false positive: `conn.open()` opens a database connection via `Connection.open`; there is no file being opened without `with`.

Checklist evidence: the receiver is a `Connection` instance with a custom `open` method; BP-PY-7's builtin-`open()` condition is unmet.

### [ ] Finding `22` — `BP-PY-1`

- Function context: `scripts/onlymaps/findings/functions/22.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_pool.py:215:1`
- Checklist pattern: handler decrements the counter and re-raises

Source excerpt:

```
        try:
            return self.__create_connection()  # <await>
        except:
            # Not sure whether a lock is needed here,
            # though better safe than sorry!
            with self.__lock:  # <async>
                self.__current_connections -= 1
            raise
```

Why this is a false positive: the bare `except:` only rolls back the pool's connection count and re-raises; the failure is propagated, not swallowed.

Checklist evidence: the suite ends with `raise` — the "bare except swallows all exceptions" condition is not satisfied.

### [ ] Finding `23` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/23.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_pool.py:280:9`
- Checklist pattern: `exec` token is a method definition, not the builtin

Source excerpt:

```
    @require_open
    def exec(self, sql: str, /, *args: Any, **kwargs: Any) -> None:  # <async>
```

Why this is a false positive: the line declares `ConnectionPool.exec`, the pool's SQL-execution method; no Python `eval`/`exec` builtin is called.

Checklist evidence: BP-PY-12 requires an `eval`/`exec`/`compile` call with a non-literal first argument; the flagged line is a `def` signature.

### [ ] Finding `24` — `CWE-94`

- Function context: `scripts/onlymaps/findings/functions/24.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_pool.py:280:9`
- Checklist pattern: `exec` token is a method definition; no code-generation sink exists

Source excerpt:

```
    @require_open
    def exec(self, sql: str, /, *args: Any, **kwargs: Any) -> None:  # <async>
```

Why this is a false positive: the line declares a method named `exec`; no code-generation or dynamic-import API is invoked with a dynamic argument.

Checklist evidence: CWE-94's "dynamic value reaches a Python code-generation sink" condition is unmet — there is no call to `eval`/`exec`/`compile`/`__import__`.

### [ ] Finding `25` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/25.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_pool.py:293:18`
- Checklist pattern: `exec` token is a method call on a pool-checked-out `Connection`

Source excerpt:

```
        with self.__connection() as conn:  # <async>
            conn.exec(sql, *args, **kwargs)  # <await>
```

Why this is a false positive: `conn.exec(...)` calls the `Connection.exec` SQL method; the builtin `exec()` is never invoked.

Checklist evidence: the call target is the user-defined `Connection.exec`; BP-PY-12's builtin-exec condition is unmet.

### [ ] Finding `26` — `BP-PY-7`

- Function context: `scripts/onlymaps/findings/functions/26.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_pool.py:429:9`
- Checklist pattern: `def open(self)` is a method definition, not an `open()` call

Source excerpt:

```
    def open(self) -> None:  # <async>
        """
        Establishes a connection to the database.
```

Why this is a false positive: the flagged line declares the pool's `open` lifecycle method; no file is opened.

Checklist evidence: BP-PY-7's "file opened without `with`" condition is unmet — this is a method declaration.

### [ ] Finding `27` — `BP-PY-1`

- Function context: `scripts/onlymaps/findings/functions/27.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_pool.py:443:1`
- Checklist pattern: handler closes the pool and re-raises

Source excerpt:

```
            try:
                conn = self.__create_connection()  # <await>
            except:
                self.close()  # <await>
                raise
```

Why this is a false positive: the bare `except:` closes the pool and re-raises; the failure propagates to the caller.

Checklist evidence: the suite ends with `raise` — the "swallows all exceptions" predicate is not satisfied.

### [ ] Finding `32` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/32.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_query.py:52:9`
- Checklist pattern: `exec` token is a method definition, not the builtin

Source excerpt:

```
    def exec(  # <async>
        self,
        sql: str,
        params: tuple[Any, ...],
        kwparams: dict[str, Any],
        /,
    ) -> None:
```

Why this is a false positive: the flagged token is the definition of `Query.exec`, the library's SQL-query runner; no `eval`/`exec` builtin call exists.

Checklist evidence: BP-PY-12's condition requires an actual `exec(` call; the line is a `def` signature whose "arguments" are the method's parameter list.

### [ ] Finding `33` — `CWE-94`

- Function context: `scripts/onlymaps/findings/functions/33.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_query.py:52:9`
- Checklist pattern: `exec` token is a method definition; no code-generation sink exists

Source excerpt:

```
    def exec(  # <async>
        self,
        sql: str,
        params: tuple[Any, ...],
        kwparams: dict[str, Any],
        /,
    ) -> None:
```

Why this is a false positive: the line declares the `Query.exec` SQL method; no code-generation/dynamic-import API is called with a dynamic expression.

Checklist evidence: CWE-94's code-generation-sink condition is unmet — there is no call, only a definition.

### [ ] Finding `34` — `CWE-89`

- Function context: `scripts/onlymaps/findings/functions/34.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_query.py:257:23`
- Checklist pattern: the sink already receives bound parameters; SQL text is never interpolated here

Source excerpt:

```
            if is_bulk:
                cursor.executemany(sql, sql_params)  # <await>
            else:
                cursor.execute(sql, sql_params)  # <await>
```

Why this is a false positive: `cursor.execute` receives the query text as an untouched variable plus a separate bound-parameter argument (`sql_params`) — exactly the parameterized-query shape the rule's fix recommends; no data is interpolated into the SQL string in this file.

Checklist evidence: CWE-89's condition "dynamic SQL string reaches execute" is unmet — the first argument is passed through verbatim and all data travels through the bound `sql_params` argument.

### [ ] Finding `35` — `CWE-89`

- Function context: `scripts/onlymaps/findings/functions/35.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_spec.py:38:9`
- Checklist pattern: `execute` token is a protocol method *definition* (PEP-249), not a call

Source excerpt:

```
    def execute(  # <async>
        self,
        operation: str,
        params: dict[str, Any] | list[Any] | tuple[Any, ...] = ...,
        /,
        *args: Any,
        **kwargs: Any,
    ) -> Self | None:
```

Why this is a false positive: the line declares the PEP-249 `Cursor.execute` protocol method; no SQL string — dynamic or static — is passed to any `execute` sink.

Checklist evidence: CWE-89's condition "dynamic SQL string reaches execute/executemany" is unmet — the flagged token is a method declaration, not a call.

### [ ] Finding `36` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/36.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_spec.py:153:9`
- Checklist pattern: `exec` token is a `@overload` method definition, not the builtin

Source excerpt:

```
    @overload
    def exec(  # <async>
        self,
        sql: str,
        /,
        *args: Any,
    ) -> None:
```

Why this is a false positive: the flagged token is a type-hint `@overload` declaration of the `exec` method; no builtin `exec()` call exists.

Checklist evidence: BP-PY-12's eval/exec-call condition is unmet — the line is an overload signature.

### [ ] Finding `37` — `CWE-94`

- Function context: `scripts/onlymaps/findings/functions/37.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_spec.py:153:9`
- Checklist pattern: `exec` token is a `@overload` method definition; no code-generation sink

Source excerpt:

```
    @overload
    def exec(  # <async>
        self,
        sql: str,
        /,
        *args: Any,
    ) -> None:
```

Why this is a false positive: the line is a type-hint overload declaration of the library's `exec` method; no code-generation API is invoked.

Checklist evidence: CWE-94's code-generation-sink condition is unmet — no call exists on the line.

### [ ] Finding `38` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/38.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_spec.py:170:9`
- Checklist pattern: `exec` token is a `@overload` method definition, not the builtin

Source excerpt:

```
    @overload
    def exec(  # <async>
        self,
        sql: str,
        /,
        **kwargs: Any,
    ) -> None:
```

Why this is a false positive: the flagged token is the second `@overload` signature of the `exec` method; no `eval`/`exec` builtin is called.

Checklist evidence: BP-PY-12's condition requires an actual `exec(` call; the line is a definition signature.

### [ ] Finding `39` — `BP-PY-7`

- Function context: `scripts/onlymaps/findings/functions/39.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_spec.py:538:9`
- Checklist pattern: `def open(self)` is a protocol method definition, not an `open()` call

Source excerpt:

```
    def open(self) -> None:  # <async>
        """
        Establishes a connection to the database.
```

Why this is a false positive: the line declares the protocol's `open` lifecycle method; no file is opened.

Checklist evidence: BP-PY-7's file-resource condition is unmet — the line is a method declaration.

### [ ] Finding `41` — `CWE-478`

- Function context: `scripts/onlymaps/findings/functions/41.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_types.py:307:1`
- Checklist pattern: unmatched values are handled by an unconditional fallback statement after the match

Source excerpt:

```
        match value:
            case str():
                dt: datetime = OnlymapsDatetime.parse_impl(value)
                return dt.date()
            case datetime():
                return value.date()
        return value
```

Why this is a false positive: every value outside the two cases falls through to the unconditional `return value` directly after the match — the same behavior a `case _:` default would provide.

Checklist evidence: the rule condition "match with no default branch" is met syntactically, but the trailing `return value` safely handles values outside the explicitly supported cases.

### [ ] Finding `43` — `BP-PY-46`

- Function context: `scripts/onlymaps/findings/functions/43.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tasks.py:6:5`
- Checklist pattern: `print` in an invoke task-runner CLI file, user-facing step output, not library logging

Source excerpt:

```
@task
def asyncio(c):
    format(c)
    print("Generating async source code...")
    c.run("python gen_async.py")
    format(c)
```

Why this is a false positive: `tasks.py` is the project's `invoke` CLI task file (all functions are `@task`-decorated entrypoints); the prints are the CLI's progress output, so the rule's "operational logging in non-script modules" condition is not satisfied.

Checklist evidence: BP-PY-46 targets library modules; the flagged module is a CLI task runner whose only execution path is `invoke <task>`, matching the rule's own "keep print … for CLIs" guidance.

### [ ] Finding `44` — `BP-PY-46`

- Function context: `scripts/onlymaps/findings/functions/44.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tasks.py:12:5`
- Checklist pattern: same CLI task-file user output as Finding 43

Source excerpt:

```
@task
def format(c):
    print("Formatting source code...")
    c.run("black onlymaps")
```

Why this is a false positive: the print is user-facing output of the `invoke format` CLI task in the project's task file, not library debug logging.

Checklist evidence: the module is a CLI task file; the print is inside a `@task` entrypoint, so the "non-script modules" condition is unmet.

### [ ] Finding `45` — `BP-PY-46`

- Function context: `scripts/onlymaps/findings/functions/45.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tasks.py:20:5`
- Checklist pattern: same CLI task-file user output as Finding 43

Source excerpt:

```
@task
def check(c):
    print("Checking for formatting issues...")
    c.run("black --check onlymaps tests")
```

Why this is a false positive: the print is progress output of the `invoke check` CLI task, not operational logging in a library module.

Checklist evidence: the flagged file is a CLI task runner (`@task` entrypoints), so BP-PY-46's library-code condition is unmet.

### [ ] Finding `46` — `BP-PY-46`

- Function context: `scripts/onlymaps/findings/functions/46.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tasks.py:23:5`
- Checklist pattern: same CLI task-file user output as Finding 43

Source excerpt:

```
    print("Running mypy...")
    c.run("mypy onlymaps tests")
    print("Running pylint...")
```

Why this is a false positive: the print announces the next step of the `invoke check` CLI task; it is user-facing CLI output, not library logging.

Checklist evidence: the module is the project's CLI task file — the "print used in library code" condition is unmet.

### [ ] Finding `47` — `BP-PY-46`

- Function context: `scripts/onlymaps/findings/functions/47.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tasks.py:25:5`
- Checklist pattern: same CLI task-file user output as Finding 43

Source excerpt:

```
    print("Running pylint...")
    c.run("pylint onlymaps")
```

Why this is a false positive: the print is step output of the `invoke check` CLI task in the project's task file, not library debug logging.

Checklist evidence: the flagged file is a CLI task runner, so the "non-script modules" condition of BP-PY-46 is unmet.

### [ ] Finding `48` — `BP-PY-7`

- Function context: `scripts/onlymaps/findings/functions/48.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/fixtures/connections.py:49:9`
- Checklist pattern: `.open(` is a call on a `Connection` fixture object, not file open

Source excerpt:

```
    conn = Connection(conn_factory)

    conn.open()  # <await>

    yield conn
```

Why this is a false positive: `conn.open()` opens the database connection via `Connection.open`; no file handle is involved.

Checklist evidence: the receiver is a `Connection` with a custom `open` method; BP-PY-7's builtin-`open()` condition is unmet.

### [ ] Finding `49` — `BP-PY-7`

- Function context: `scripts/onlymaps/findings/functions/49.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/fixtures/connections.py:83:9`
- Checklist pattern: `.open(` is a call on a `ConnectionPool` fixture object, not file open

Source excerpt:

```
    conn = ConnectionPool(conn_factory, **pool_kwargs)

    conn.open()  # <await>

    yield conn
```

Why this is a false positive: `conn.open()` opens the pool's database connections via `ConnectionPool.open`; no file resource is created.

Checklist evidence: the call target is a user-defined method; the "file opened without `with`" condition is unmet.

### [ ] Finding `50` — `BP-PY-42`

- Function context: `scripts/onlymaps/findings/functions/50.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/fixtures/connections.py:142:1`
- Checklist pattern: try/except is fixture teardown cleanup, not a test function expecting failure

Source excerpt:

```
@pytest.fixture(scope="function")
def dbapiv2(connection: Connection) -> Iterator[PyDbAPIv2Connection]:  # <async>
    ...
    yield conn

    try:
        conn.close()  # <replace:await co_exec(conn.close)>
    except:
        pass
```

Why this is a false positive: BP-PY-42's condition is a *test function* using try/except instead of `assertRaises`/`pytest.raises`; the flagged block is teardown cleanup in a `@pytest.fixture` (no `def test_*` anywhere in the file) that ignores close errors.

Checklist evidence: the rule restricts to test-function bodies; the flagged try/except sits in a fixture's teardown, not in a test that "expects failure".

### [ ] Finding `55` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/55.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/fixtures/containers.py:139:18`
- Checklist pattern: `exec` token is a method call on a testcontainers container object, not the builtin

Source excerpt:

```
        oracledb.exec(
            [
                "bash",
                "-c",
                (...),
            ]
        )
```

Why this is a false positive: `oracledb` is the `OracleDbContainer` fixture object (from `testcontainers.oracle`) and `oracledb.exec([...])` runs a command inside the Docker container; Python's `exec()` builtin is never invoked.

Checklist evidence: BP-PY-12's condition requires the `eval`/`exec` builtin; the call target is a user-defined container method.

### [ ] Finding `56` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/56.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/fixtures/containers.py:152:18`
- Checklist pattern: same container-method `exec` call as Finding 55

Source excerpt:

```
        # Create any necessary tables.
        oracledb.exec(
            [
                "bash",
                "-c",
                (...),
            ]
        )
```

Why this is a false positive: `oracledb.exec(...)` is the testcontainers container's command runner, not Python's `exec()` builtin.

Checklist evidence: the call target is the `OracleDbContainer` method; the eval/exec-builtin condition is unmet.

### [ ] Finding `57` — `CWE-89`

- Function context: `scripts/onlymaps/findings/functions/57.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/fixtures/containers.py:180:13`
- Checklist pattern: the first argument is a compile-time SQL constant, not a dynamic/interpolated string

Source excerpt:

```
    with sqlite3.connect(database=db, isolation_level="DEFERRED") as conn:
        conn.execute(SQL.CREATE_TEST_TABLE)
        yield SqliteContainer(dbname=db)
```

Why this is a false positive: `conn.execute` receives `SQL.CREATE_TEST_TABLE`, a module-level constant SQL string; no dynamic expression or interpolation is involved, so no special elements can be injected.

Checklist evidence: CWE-89's "dynamic SQL string reaches execute" condition is unmet — the argument is a static constant reference.

### [ ] Finding `58` — `BP-PY-42`

- Function context: `scripts/onlymaps/findings/functions/58.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/fixtures/executors.py:77:1`
- Checklist pattern: try/except is an executor-drain helper, not a test function expecting failure

Source excerpt:

```
    async def wait(self) -> None:
        """
        Waits for all submitted coroutine to finish execution.
        """
        for task in self.__tasks:
            try:
                await task
            except:
                pass
```

Why this is a false positive: the flagged try/except is inside `TaskExecutor.wait()`, a fixture helper that drains submitted tasks; it is not a test function using try/except as a failure assertion.

Checklist evidence: BP-PY-42 restricts to test-function bodies (`def test_*`); `wait` is a helper method, so the condition is unmet.

### [ ] Finding `64` — `BP-PY-7`

- Function context: `scripts/onlymaps/findings/functions/64.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_connect.py:40:11`
- Checklist pattern: `.open(` is a call on a `Database` object, not file open

Source excerpt:

```
        db = connect(conn_str, **kwargs)

        db.open()  # <await>

        db.close()  # <await>
```

Why this is a false positive: `db.open()` invokes the library's connection-opening method; no file handle is created.

Checklist evidence: the receiver is a `Database` object with a custom `open` method; BP-PY-7's builtin-`open()` condition is unmet.

### [ ] Finding `65` — `BP-PY-7`

- Function context: `scripts/onlymaps/findings/functions/65.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_connect.py:97:11`
- Checklist pattern: `.open(` is a call on a `Database` object, not file open

Source excerpt:

```
        db = connect(conn_factory)

        db.open()  # <await>

        db.close()  # <await>
```

Why this is a false positive: `db.open()` calls the custom connection-opening method; no file resource is involved.

Checklist evidence: the call target is a user-defined method, not the builtin `open()`.

### [ ] Finding `66` — `BP-PY-7`

- Function context: `scripts/onlymaps/findings/functions/66.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_connection.py:172:27`
- Checklist pattern: `.open(` is a call on a `Connection` object under `pytest.raises`, not file open

Source excerpt:

```
            with pytest.raises(RuntimeError) as exc_info:
                connection.open()  # <await>
```

Why this is a false positive: `connection.open()` is the custom database-open method invoked to assert the expected `RuntimeError`; no file is opened.

Checklist evidence: the receiver is a `Connection` with a custom `open`; BP-PY-7's file-resource condition is unmet.

### [ ] Finding `67` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/67.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_connection.py:234:30`
- Checklist pattern: `exec` token is a method call on a `Connection` object, not the builtin

Source excerpt:

```
                connection_B.exec(  # <await>
                    SQL.INSERT_INTO_TEST_TABLE(connection.driver), uid
                )
```

Why this is a false positive: `connection_B.exec(...)` calls the library's SQL-execution method (`Connection.exec`); Python's `exec()` builtin is never invoked.

Checklist evidence: BP-PY-12 requires the `eval`/`exec` builtin; the call target is the user-defined `Connection.exec`.

### [ ] Finding `69` — `BP-PY-7`

- Function context: `scripts/onlymaps/findings/functions/69.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_database.py:52:15`
- Checklist pattern: `.open(` is a call on a `Database` object under `pytest.raises`, not file open

Source excerpt:

```
        with pytest.raises(RuntimeError) as exc_info:
            db.open()  # <await>
```

Why this is a false positive: `db.open()` is the custom database-open method used to trigger the expected `RuntimeError`; no file handle is involved.

Checklist evidence: the call target is a user-defined method; BP-PY-7's "file opened without `with`" condition is unmet.

### [ ] Finding `70` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/70.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_database.py:94:19`
- Checklist pattern: `exec` token is a method call on a `Database` object, not the builtin

Source excerpt:

```
        assert db.exec( # type: ignore # <await>
            SQL.SELECT_SINGLE_SCALAR(db.driver), 1
        ) is None
```

Why this is a false positive: `db.exec(...)` calls the library's SQL-execution method; the builtin `exec()` is not invoked.

Checklist evidence: the call target is the user-defined `exec` method of the `Database` API; the eval/exec-builtin condition is unmet.

### [ ] Finding `71` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/71.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_database.py:108:12`
- Checklist pattern: same `Database.exec` method call as Finding 70

Source excerpt:

```
        db.exec(SQL.INSERT_INTO_TEST_TABLE(db.driver), batch)  # <await>
```

Why this is a false positive: `db.exec(...)` is a call to the library's SQL-execution method, not Python's `exec()` builtin.

Checklist evidence: BP-PY-12's builtin-exec condition is unmet — the call target is the `Database.exec` method.

### [ ] Finding `72` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/72.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_database.py:274:16`
- Checklist pattern: same `Database.exec` method call as Finding 70

Source excerpt:

```
        with db.transaction():  # <async>
            db.exec(SQL.INSERT_INTO_TEST_TABLE(db.driver), uid)  # <await>
```

Why this is a false positive: `db.exec(...)` is a call to the library's SQL method; no `eval`/`exec` builtin is invoked.

Checklist evidence: the call target is the user-defined `Database.exec` method, so BP-PY-12's condition is unmet.

### [ ] Finding `73` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/73.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_database.py:290:16`
- Checklist pattern: same `Database.exec` method call as Finding 70

Source excerpt:

```
        with db.transaction():  # <async>
            db.exec(SQL.INSERT_INTO_TEST_TABLE(db.driver), uid)  # <await>
```

Why this is a false positive: `db.exec(...)` invokes the library's SQL-execution method, not the builtin `exec()`.

Checklist evidence: the call target is the `Database.exec` method; the eval/exec-builtin condition is unmet.

### [ ] Finding `74` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/74.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_database.py:319:16`
- Checklist pattern: same `Database.exec` method call as Finding 70

Source excerpt:

```
        with db.transaction():  # <async>
            db.exec(SQL.INSERT_INTO_TEST_TABLE(db.driver), uid)  # <await>
```

Why this is a false positive: `db.exec(...)` is a call to the library's SQL method, not Python's `exec()`.

Checklist evidence: the call target is the user-defined method; BP-PY-12's condition is unmet.

### [ ] Finding `75` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/75.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_database.py:337:20`
- Checklist pattern: same `Database.exec` method call as Finding 70

Source excerpt:

```
            with db.transaction():  # <async>
                db.exec(SQL.INSERT_INTO_TEST_TABLE(db.driver), uid)  # <await>
                raise Exception()
```

Why this is a false positive: `db.exec(...)` calls the library's SQL-execution method inside a test; the `exec` token is not the Python builtin.

Checklist evidence: the call target is the `Database.exec` method; the eval/exec-builtin condition is unmet.

### [ ] Finding `83` — `BP-PY-7`

- Function context: `scripts/onlymaps/findings/functions/83.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_database.py:452:11`
- Checklist pattern: `.open(` is a call on a `Database` object under `pytest.raises`, not file open

Source excerpt:

```
    with pytest.raises(PyDbFactoryError):
        db.open()  # <await>
```

Why this is a false positive: `db.open()` is the custom database-open method used to trigger the expected factory error; no file handle exists.

Checklist evidence: the call target is a user-defined method; BP-PY-7's file-resource condition is unmet.

### [ ] Finding `84` — `BP-PY-7`

- Function context: `scripts/onlymaps/findings/functions/84.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_database.py:480:15`
- Checklist pattern: `.open(` is a call on a `Database` object in a concurrency test, not file open

Source excerpt:

```
        nonlocal error_counter
        try:
            db.open()  # <await>
        except RuntimeError:
```

Why this is a false positive: `db.open()` invokes the custom connection-opening method inside a worker function; no file is opened.

Checklist evidence: the receiver is a `Database` with a custom `open`; the "file opened without `with`" condition is unmet.

### [ ] Finding `87` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/87.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_pool.py:49:26`
- Checklist pattern: `exec` token is a method call on a `ConnectionPool` object, not the builtin

Source excerpt:

```
            with pool.transaction():  # <async>
                try:
                    pool.exec(SQL.INSERT_INTO_TEST_TABLE(pool.driver), uid)  # <await>
```

Why this is a false positive: `pool.exec(...)` calls the library's SQL-execution method (`ConnectionPool.exec`); Python's `exec()` builtin is never invoked.

Checklist evidence: the call target is the user-defined `ConnectionPool.exec`; BP-PY-12's builtin-exec condition is unmet.

### [ ] Finding `88` — `BP-PY-1`

- Function context: `scripts/onlymaps/findings/functions/88.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_pool.py:63:1`
- Checklist pattern: the handler records the exception into the test's result variable, which is asserted later

Source excerpt:

```
            try:
                result = pool.fetch_one_or_none(  # <await>
                    int, SQL.SELECT_FROM_TEST_TABLE(pool.driver), uid
                )
            except Exception as exc:
                result = exc
            finally:
                continue_counter = 2
```

Why this is a false positive: the handler stores the exception in `result`, and the test later asserts on it (`assert isinstance(result, Exception)` at line 76) — the failure is captured as test evidence, not hidden.

Checklist evidence: BP-PY-1's "broad except hides failures" predicate is not satisfied — the exception is recorded into a variable that the test asserts on.

### [ ] Finding `89` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/89.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_pool.py:176:22`
- Checklist pattern: same `ConnectionPool.exec` method call as Finding 87

Source excerpt:

```
                continue_flag = True
                sleep(2)  # <await>
                pool.exec(SQL.INSERT_INTO_TEST_TABLE(pool.driver), uid)  # <await>
```

Why this is a false positive: `pool.exec(...)` is a call to the library's SQL-execution method, not the builtin `exec()`.

Checklist evidence: the call target is the user-defined method; the eval/exec-builtin condition is unmet.

### [ ] Finding `90` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/90.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_pool.py:218:16`
- Checklist pattern: same `ConnectionPool.exec` method call as Finding 87

Source excerpt:

```
    small_pool.exec(SQL.KILL_CONNECTION(small_pool.driver))  # <await>
```

Why this is a false positive: `small_pool.exec(...)` invokes the library's SQL method; no `eval`/`exec` builtin is called.

Checklist evidence: the call target is `ConnectionPool.exec`; BP-PY-12's condition is unmet.

### [ ] Finding `92` — `BP-PY-1`

- Function context: `scripts/onlymaps/findings/functions/92.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_query.py:78:1`
- Checklist pattern: the handler records the exception into the test's result variable for later assertion

Source excerpt:

```
        try:
            result: Any = db.fetch_one(t, query, param)  # <await>
        except Exception as e:
            result = e
```

Why this is a false positive: the handler stores the exception in `result`, which the test subsequently asserts against (`assert result == scalar` / type checks below) — the failure is captured as test evidence, not hidden.

Checklist evidence: BP-PY-1's "broad except hides failures" predicate is not satisfied — the exception is recorded into the asserted result variable.

### [ ] Finding `93` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/93.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_query.py:343:16`
- Checklist pattern: `exec` token is a method call on a `Database` object, not the builtin

Source excerpt:

```
        with pytest.raises(ValueError):
            db.exec(  # <await>
                SQL.SELECT_NONE, 1, param=1
            )  # type: ignore
```

Why this is a false positive: `db.exec(...)` calls the library's SQL-execution method; the builtin `exec()` is not invoked.

Checklist evidence: the call target is the user-defined `Database.exec`; BP-PY-12's builtin-exec condition is unmet.

### [ ] Finding `94` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/94.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_query.py:372:12`
- Checklist pattern: same `Database.exec` method call as Finding 93 (SQL text passed as the method's argument)

Source excerpt:

```
        db.exec(  # <await>
            f"""
            CREATE TEMPORARY TABLE {tmp_table} (
                {c1} VARCHAR(100),
                {c2} VARCHAR(100)
            )
            """
        )
```

Why this is a false positive: `db.exec(...)` is the library's SQL-execution method; the f-string is its SQL argument, not a code string passed to Python's `exec()`.

Checklist evidence: the call target is `Database.exec` (a user-defined method), so the eval/exec-builtin condition is unmet regardless of the argument.

### [ ] Finding `95` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/95.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_query.py:394:12`
- Checklist pattern: same `Database.exec` method call as Finding 93

Source excerpt:

```
        db.exec(  # <await>
            f"INSERT INTO {tmp_table}({c1}, {c2}) VALUES({plchld_1}, {plchld_2})",
            Bulk(params),
        )
```

Why this is a false positive: `db.exec(...)` invokes the library's SQL method with an SQL string and bulk parameters; Python's `exec()` builtin is not called.

Checklist evidence: the call target is the user-defined method; BP-PY-12's condition is unmet.

### [ ] Finding `96` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/96.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_query.py:415:12`
- Checklist pattern: same `Database.exec` method call as Finding 93

Source excerpt:

```
        db.exec(  # <await>
            f"""
            CREATE TEMPORARY TABLE {tmp_table} (
                {c1} VARCHAR(100),
                {c2} VARCHAR(100)
            )
            """
        )
```

Why this is a false positive: `db.exec(...)` is the library's SQL-execution method; the f-string is the SQL query argument, not code passed to the `exec()` builtin.

Checklist evidence: the call target is `Database.exec`, so the eval/exec-builtin condition is unmet.

### [ ] Finding `97` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/97.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_query.py:439:12`
- Checklist pattern: same `Database.exec` method call as Finding 93

Source excerpt:

```
        db.exec(  # <await>
            f"INSERT INTO {tmp_table}({c1}, {c2}) VALUES({plchld_1}, {plchld_2})",
            Bulk(params),
        )
```

Why this is a false positive: `db.exec(...)` calls the library's SQL method; no `eval`/`exec` builtin is invoked.

Checklist evidence: the call target is the user-defined `Database.exec` method; BP-PY-12's condition is unmet.

### [ ] Finding `98` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/98.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_query.py:475:16`
- Checklist pattern: same `Database.exec` method call as Finding 93

Source excerpt:

```
            bulk = Bulk([[1]])
            # fmt: off
            db.exec( # <await>
                "SELECT 1", bulk, param=1
            )  #  type: ignore
```

Why this is a false positive: `db.exec("SELECT 1", ...)` invokes the library's SQL-execution method with a static SQL string; the `exec` token is not Python's builtin.

Checklist evidence: the call target is the user-defined `Database.exec`; the eval/exec-builtin condition is unmet.

## Uncertain findings

None.

## True positives

### CWE-1046 — Creation of Immutable Text Using String Concatenation

Rule condition (`rules_tier_b_runtime.go` `detectCWE1046`): `+=` inside a `for`/`while` loop where the accumulator has string-literal initialization evidence (`textAccumulatorEvidence`).

| Finding | Source | Reason |
| --- | --- | --- |
| 1 | `gen_async.py:212` | `async_file += line` inside `for line in file.readlines()`; `async_file` initialized to `""` at line 198 — string accumulator built in a loop. |
| 99 | `tests/utils.py:404` | `query += f"…"` inside `for i, field_name in enumerate(...)`; `query` initialized to `"SELECT "` at line 402. |

### BP-PY-1 — Bare Except Clause

Rule condition (`rules_core.go` `detectBPPY1`): bare `except:` or broad `except Exception`/`BaseException` whose suite neither re-raises nor records the failure.

| Finding | Source | Reason |
| --- | --- | --- |
| 10 | `onlymaps/_connection.py:455` | bare `except:` in `__cursor(test_connection=True)`; the exception is discarded (recovery path runs regardless) with no re-raise or record. |
| 11 | `onlymaps/_connection.py:461` | bare `except:` whose suite is only `pass` — swallows the `__close()` failure. |
| 28 | `onlymaps/_pool.py:476` | bare `except:` with only `pass` while draining connections in `close()`. |
| 40 | `onlymaps/_types.py:189` | bare `except:` returns the raw `value`, silently discarding the parse failure with no record. |
| 51 | `tests/fixtures/connections.py:144` | bare `except:` with only `pass` in fixture teardown. |
| 59 | `tests/fixtures/executors.py:79` | bare `except:` with only `pass` in `TaskExecutor.wait()`. |
| 77 | `tests/test_database.py:339` | bare `except:` with only `pass` after the deliberate `raise Exception()`. |
| 81 | `tests/test_database.py:413` | bare `except:` with only `pass`. |
| 85 | `tests/test_database.py:516` | bare `except:` with only `pass`. |

### BP-PY-2 — Except Pass

Rule condition (`rules_core.go` `detectBPPY2`): an except clause whose direct suite consists solely of `pass`. The rule has no test-file or optional-import exemption.

| Finding | Source | Reason |
| --- | --- | --- |
| 12 | `onlymaps/_connection.py:461` | `except:` suite is only `pass`. |
| 16 | `onlymaps/_drivers.py:33` | `except ImportError:` suite is only `pass` — the rule has no exemption for optional-import guards. |
| 29 | `onlymaps/_pool.py:476` | `except:` suite is only `pass`. |
| 52 | `tests/fixtures/connections.py:144` | `except:` suite is only `pass`. |
| 60 | `tests/fixtures/executors.py:79` | `except:` suite is only `pass`. |
| 78 | `tests/test_database.py:339` | `except:` suite is only `pass`. |
| 82 | `tests/test_database.py:413` | `except:` suite is only `pass`. |
| 86 | `tests/test_database.py:516` | `except:` suite is only `pass`. |

### CWE-390 — Detection of Error Condition Without Action

Rule condition (`rules_platform.go` `detectCWE390`): an except clause whose direct body is `pass` (`exceptPassStart`). The detector deliberately does not judge whether documented handling elsewhere is adequate.

| Finding | Source | Reason |
| --- | --- | --- |
| 13 | `onlymaps/_connection.py:461` | `except:` followed only by `pass`. |
| 17 | `onlymaps/_drivers.py:33` | `except ImportError:` followed only by `pass`. |
| 30 | `onlymaps/_pool.py:476` | `except:` followed only by `pass`. |
| 53 | `tests/fixtures/connections.py:144` | `except:` followed only by `pass`. |
| 61 | `tests/fixtures/executors.py:79` | `except:` followed only by `pass`. |
| 79 | `tests/test_database.py:339` | `except:` followed only by `pass`. |

### CWE-1071 — Empty Code Block

Rule condition (`rules_tier_b_runtime.go` `detectCWE1071`): `pyTierBEmptyExceptRE` — an exception handler containing only `pass`.

| Finding | Source | Reason |
| --- | --- | --- |
| 14 | `onlymaps/_connection.py:461` | handler body is only `pass`. |
| 18 | `onlymaps/_drivers.py:33` | handler body is only `pass`. |
| 31 | `onlymaps/_pool.py:476` | handler body is only `pass`. |
| 54 | `tests/fixtures/connections.py:144` | handler body is only `pass`. |
| 62 | `tests/fixtures/executors.py:79` | handler body is only `pass`. |
| 80 | `tests/test_database.py:339` | handler body is only `pass`. |

### CWE-397 — Declaration of Throws for Generic Exception

Rule condition (`rules_platform.go` `detectCWE397`): `pyGenericRaiseRE` — `raise Exception` / `raise BaseException` directly constructed. No test-file exemption exists.

| Finding | Source | Reason |
| --- | --- | --- |
| 68 | `tests/test_connection.py:290` | `raise Exception()` in the chaos-monkey mock — matches the generic-raise pattern. |
| 76 | `tests/test_database.py:338` | `raise Exception()` inside the transaction under test. |
| 91 | `tests/test_pool.py:271` | `raise Exception()` in the `Connection.__init__` mock. |

### CWE-478 — Missing Default Case in Multiple Condition Expression

Rule condition (`rules_platform.go` `detectCWE478`): a `match` with ≥2 immediate case branches and no wildcard (`matchWithoutDefaultStart`).

| Finding | Source | Reason |
| --- | --- | --- |
| 42 | `onlymaps/_utils.py:283` | `match driver:` covers 7 of the 8 `Driver` enum members with no `case _` and no fallback statement after the match — an unmatched driver leaves `module` unbound and is not handled. |

### BP-PY-41 — pytest assert With Side Effects Only

Rule condition (`rules_testing.go` `detectBPPY41`): a `test_*` function body containing side-effect calls but no assertion.

| Finding | Source | Reason |
| --- | --- | --- |
| 63 | `tests/test_connect.py:30` | `test_connect_via_conn_str` only assigns, opens and closes — no `assert`/`pytest.raises` anywhere in the body. |

## Final evidence

- Delegated reviewers: none (single-reviewer audit)
- Chunk evidence: `scripts/onlymaps/chunks`
- Function evidence: `scripts/onlymaps/findings/functions`
- Validation: `git diff --check` — pass

## Post-fix remaining-FP audit (2026-08-02)

### Run metadata

```yaml
timestamp: 2026-08-02T16:38:00Z
repository: onlymaps
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps
branch: main
commit: 6444a59
scanner_commit: b5b8fde (FP-reduction fix, binary rebuilt 2026-08-02 16:29)
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps
chunk_path: scripts/onlymaps/chunks
function_context_path: scripts/onlymaps/findings/functions
```

### Scan evidence

- Build command: `(pre-built) ./bin/goslop`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/onlymaps/chunks -context-dir scripts/onlymaps/findings/functions real-repos/onlymaps`
- Findings: `60` (down from `99` pre-fix)
- Chunks reviewed: `scripts/onlymaps/chunks/Chunk_1_25.txt`, `scripts/onlymaps/chunks/Chunk_26_50.txt`, `scripts/onlymaps/chunks/Chunk_51_60.txt`
- Function contexts reviewed: `scripts/onlymaps/findings/functions/<id>.txt` for every proposed false positive (2, 3, 4, 5, 6, 15, 16, 17, 18, 19, 24, 25, 26, 27, 28, 30, 32, 33, 34, 35, 36, 41, 57, 59)
- Fix effectiveness: all 36 audited true positives still fire; 39 of the 63 audited false positives (all BP-PY-7 `open()`/`def open` cases, all `exec`/`execute` method-call and protocol-definition cases, both BP-PY-42 cases, and the bound-parameter `execute` case) are gone; 24 audited FPs re-appear.

### Audit checklist

- [x] Read every assigned chunk under `scripts/onlymaps/chunks`.
- [x] Read `scripts/onlymaps/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

### Classification summary (fresh run)

Matching by `Source:` path against the audited TP/FP lists of the pre-fix audit.

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 24 | 2, 3, 4, 5, 6, 15, 16, 17, 18, 19, 24, 25, 26, 27, 28, 30, 32, 33, 34, 35, 36, 41, 57, 59 |
| True positive | 36 | 1, 7, 8, 9, 10, 11, 12, 13, 14, 20, 21, 22, 23, 29, 31, 37, 38, 39, 40, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 58, 60 |
| Uncertain | 0 | — |

Every remaining false positive is a re-appearing audited FP; there are no new findings.

## False positives (remaining)

### [ ] Findings `3`, `4` — `BP-PY-12`, `CWE-94`

- Function context: `scripts/onlymaps/findings/functions/3.txt`, `scripts/onlymaps/findings/functions/4.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_connection.py:232:9`
- Checklist pattern: the `exec` token is a method *definition*, not a call of the builtin

Source excerpt:

```
    @require_open
    @__require_not_iter_same_ctx
    def exec(self, sql: str, /, *args: Any, **kwargs: Any) -> None:  # <async>
```

Why this is a false positive: the line declares the library's `Connection.exec` SQL method; no `eval`/`exec` builtin is invoked, so neither BP-PY-12 nor CWE-94 has a call reaching a code-generation sink.

Checklist evidence: both rules require an actual `exec(`/`eval(` call; the flagged token is the `def` signature of a method named `exec`.

### [ ] Finding `2` — `BP-PY-1`

- Function context: `scripts/onlymaps/findings/functions/2.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_connection.py:215:1`
- Checklist pattern: handler rolls back and re-raises; the failure is propagated

Source excerpt:

```
            except:
                is_query_successful = False
                if not self.__in_transaction:
                    self.__conn.rollback()  # <await>
                raise
```

Why this is a false positive: the bare `except:` resets the success flag, rolls back, and re-raises (`_connection.py:215-219`); nothing is swallowed.

Checklist evidence: the suite ends in `raise`, so the "neither re-raises nor records the failure" predicate of BP-PY-1 is unmet.

### [ ] Finding `5` — `BP-PY-1`

- Function context: `scripts/onlymaps/findings/functions/5.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_connection.py:376:1`
- Checklist pattern: handler rolls back and re-raises; the failure is propagated

Source excerpt:

```
            except:
                self.__conn.rollback()  # <await>
                raise
```

Why this is a false positive: the transaction handler rolls back and re-raises (`_connection.py:376-378`); the failure propagates to the caller.

Checklist evidence: the suite ends in `raise` — the "hides failures" predicate is unmet.

### [ ] Finding `6` — `BP-PY-1`

- Function context: `scripts/onlymaps/findings/functions/6.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_connection.py:408:1`
- Checklist pattern: handler resets state and re-raises

Source excerpt:

```
            except:
                self.__is_open = False
                raise
```

Why this is a false positive: the handler clears the open flag and re-raises (`_connection.py:408-410`); the failure is not hidden.

Checklist evidence: the suite ends in `raise` — BP-PY-1's condition is unmet.

### [ ] Finding `15` — `CWE-478`

- Function context: `scripts/onlymaps/findings/functions/15.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_drivers.py:403:1`
- Checklist pattern: unmatched values are handled by an unconditional fallback statement immediately after the match

Source excerpt:

```
            match value:
                # Force `BINARY_DOUBLE` TYPE even if the decimal part is zero.
                case float():
                    return cursor.var(
                        oracledb.DB_TYPE_BINARY_DOUBLE, arraysize=num_elements
                    )
                # Force `TIMESTAMP` type if datetime has microseconds.
                case datetime() if value.microsecond > 0:
                    return cursor.var(
                        oracledb.DB_TYPE_TIMESTAMP, arraysize=num_elements
                    )
            return None
```

Why this is a false positive: every value outside the two cases falls through to the unconditional `return None` directly after the match — the same behavior a `case _:` would provide.

Checklist evidence: the "no default and no fallback handling" condition of CWE-478 is unmet — the match is immediately followed by an unconditional return.

### [ ] Findings `17`, `18` — `BP-PY-12`, `CWE-94`

- Function context: `scripts/onlymaps/findings/functions/17.txt`, `scripts/onlymaps/findings/functions/18.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_pool.py:280:9`
- Checklist pattern: the `exec` token is a method *definition*, not the builtin

Source excerpt:

```
    @require_open
    def exec(self, sql: str, /, *args: Any, **kwargs: Any) -> None:  # <async>
```

Why this is a false positive: the line declares the pool's `ConnectionPool.exec` SQL method; no `eval`/`exec` builtin call exists.

Checklist evidence: BP-PY-12/CWE-94 require a code-generation call; the flagged line is a `def` signature.

### [ ] Finding `16` — `BP-PY-1`

- Function context: `scripts/onlymaps/findings/functions/16.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_pool.py:215:1`
- Checklist pattern: handler decrements the connection counter and re-raises

Source excerpt:

```
        try:
            return self.__create_connection()  # <await>
        except:
            # Not sure whether a lock is needed here,
            # though better safe than sorry!
            with self.__lock:  # <async>
                self.__current_connections -= 1
            raise
```

Why this is a false positive: the handler only rolls back the pool's connection count and re-raises (`_pool.py:215-220`); the failure propagates.

Checklist evidence: the suite ends in `raise` — the "swallows all exceptions" predicate is unmet.

### [ ] Finding `19` — `BP-PY-1`

- Function context: `scripts/onlymaps/findings/functions/19.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_pool.py:443:1`
- Checklist pattern: handler closes the pool and re-raises

Source excerpt:

```
            except:
                self.close()  # <await>
                raise
```

Why this is a false positive: the handler closes the pool and re-raises (`_pool.py:443-445`); the failure propagates to the caller.

Checklist evidence: the suite ends in `raise` — BP-PY-1's condition is unmet.

### [ ] Findings `24`, `25` — `BP-PY-12`, `CWE-94`

- Function context: `scripts/onlymaps/findings/functions/24.txt`, `scripts/onlymaps/findings/functions/25.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_query.py:52:9`
- Checklist pattern: the `exec` token is a method *definition*, not the builtin

Source excerpt:

```
    def exec(  # <async>
        self,
        sql: str,
        params: tuple[Any, ...],
        kwparams: dict[str, Any],
        /,
    ) -> None:
```

Why this is a false positive: the line declares the library's `Query.exec` SQL-query runner; no `eval`/`exec` builtin is invoked.

Checklist evidence: the rules require an actual call to a code-generation sink; the flagged token is a `def` signature.

### [ ] Findings `26`, `27` — `BP-PY-12`, `CWE-94`

- Function context: `scripts/onlymaps/findings/functions/26.txt`, `scripts/onlymaps/findings/functions/27.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_spec.py:153:9`
- Checklist pattern: the `exec` token is a `@overload` method *definition*, not the builtin

Source excerpt:

```
    @overload
    def exec(  # <async>
        self,
        sql: str,
        /,
        *args: Any,
    ) -> None:
```

Why this is a false positive: the line is a type-hint `@overload` signature of the `exec` method; no code-generation API is called.

Checklist evidence: BP-PY-12/CWE-94 require a `eval`/`exec`/`compile` call; the line is a declaration.

### [ ] Finding `28` — `BP-PY-12`

- Function context: `scripts/onlymaps/findings/functions/28.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_spec.py:170:9`
- Checklist pattern: the `exec` token is a `@overload` method *definition*, not the builtin

Source excerpt:

```
    @overload
    def exec(  # <async>
        self,
        sql: str,
        /,
        **kwargs: Any,
    ) -> None:
```

Why this is a false positive: the line is the second `@overload` signature of the `exec` method; no `eval`/`exec` builtin call exists.

Checklist evidence: BP-PY-12 requires an actual `exec(` call; the line is a definition signature.

### [ ] Finding `30` — `CWE-478`

- Function context: `scripts/onlymaps/findings/functions/30.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/onlymaps/_types.py:307:1`
- Checklist pattern: unmatched values are handled by an unconditional fallback statement after the match

Source excerpt:

```
        match value:
            case str():
                dt: datetime = OnlymapsDatetime.parse_impl(value)
                return dt.date()
            case datetime():
                return value.date()
        return value
```

Why this is a false positive: every value outside the two cases falls through to the unconditional `return value` directly after the match — the same behavior a `case _:` default would provide.

Checklist evidence: the "no default and no fallback handling" condition is unmet — the match is immediately followed by an unconditional return.

### [ ] Finding `32` — `BP-PY-46`

- Function context: `scripts/onlymaps/findings/functions/32.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tasks.py:6:5`
- Checklist pattern: `print` in an invoke task-runner CLI file, user-facing step output, not library logging

Source excerpt:

```
@task
def asyncio(c):
    format(c)
    print("Generating async source code...")
    c.run("python gen_async.py")
```

Why this is a false positive: `tasks.py` is the project's `invoke` CLI task file; the print is progress output of a `@task` entrypoint, matching the rule's own "keep print … for CLIs" guidance.

Checklist evidence: BP-PY-46 targets library modules; the flagged file is a CLI task runner.

### [ ] Finding `33` — `BP-PY-46`

- Function context: `scripts/onlymaps/findings/functions/33.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tasks.py:12:5`
- Checklist pattern: same CLI task-file user output as Finding 32

Source excerpt:

```
@task
def format(c):
    print("Formatting source code...")
    c.run("black onlymaps")
```

Why this is a false positive: the print is user-facing output of the `invoke format` CLI task, not library debug logging.

Checklist evidence: the flagged module is a CLI task file; the "non-script modules" condition is unmet.

### [ ] Finding `34` — `BP-PY-46`

- Function context: `scripts/onlymaps/findings/functions/34.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tasks.py:20:5`
- Checklist pattern: same CLI task-file user output as Finding 32

Source excerpt:

```
@task
def check(c):
    print("Checking for formatting issues...")
    c.run("black --check onlymaps tests")
```

Why this is a false positive: the print is progress output of the `invoke check` CLI task, not operational logging in a library module.

Checklist evidence: the flagged file is a CLI task runner, so BP-PY-46's library-code condition is unmet.

### [ ] Finding `35` — `BP-PY-46`

- Function context: `scripts/onlymaps/findings/functions/35.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tasks.py:23:5`
- Checklist pattern: same CLI task-file user output as Finding 32

Source excerpt:

```
    print("Running mypy...")
    c.run("mypy onlymaps tests")
```

Why this is a false positive: the print announces the next step of the `invoke check` CLI task; it is user-facing CLI output.

Checklist evidence: the module is the project's CLI task file — the "print used in library code" condition is unmet.

### [ ] Finding `36` — `BP-PY-46`

- Function context: `scripts/onlymaps/findings/functions/36.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tasks.py:25:5`
- Checklist pattern: same CLI task-file user output as Finding 32

Source excerpt:

```
    print("Running pylint...")
    c.run("pylint onlymaps")
```

Why this is a false positive: the print is step output of the `invoke check` CLI task in the project's task file, not library debug logging.

Checklist evidence: the flagged file is a CLI task runner, so the "non-script modules" condition of BP-PY-46 is unmet.

### [ ] Finding `41` — `CWE-89`

- Function context: `scripts/onlymaps/findings/functions/41.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/fixtures/containers.py:180:13`
- Checklist pattern: the first argument is a compile-time SQL constant, not a dynamic/interpolated string

Source excerpt:

```
    with sqlite3.connect(database=db, isolation_level="DEFERRED") as conn:
        conn.execute(SQL.CREATE_TEST_TABLE)
        yield SqliteContainer(dbname=db)
```

Why this is a false positive: `conn.execute` receives `SQL.CREATE_TEST_TABLE`, a module-level constant SQL string; no dynamic expression or interpolation is involved, so no special elements can be injected.

Checklist evidence: CWE-89's "dynamic SQL string reaches execute" condition is unmet — the argument is a static constant reference.

### [ ] Finding `57` — `BP-PY-1`

- Function context: `scripts/onlymaps/findings/functions/57.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_pool.py:63:1`
- Checklist pattern: the handler records the exception into the test's result variable, which is asserted later

Source excerpt:

```
            except Exception as exc:
                result = exc
            finally:
                continue_counter = 2

        executor.submit(fn_1)
        executor.submit(fn_2)
        executor.wait()  # <await>

        # SQL Server implements READ COMMITTED isolation via
        # locking, not MVCC, therefore the query should fail
        # instead due to the lock held on the table.
        if pool.driver == Driver.SQL_SERVER:
            assert isinstance(result, Exception)
        else:
            assert result is None
```

Why this is a false positive: the handler stores the exception in `result`, and the test asserts on it at line 76 (`assert isinstance(result, Exception)`); the failure is captured as test evidence, not hidden.

Checklist evidence: BP-PY-1's "broad except hides failures" predicate is not satisfied — the exception is recorded into a variable that the test asserts on.

### [ ] Finding `59` — `BP-PY-1`

- Function context: `scripts/onlymaps/findings/functions/59.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps/tests/test_query.py:78:1`
- Checklist pattern: the handler records the exception into the test's result variable for later assertion

Source excerpt:

```
        try:
            result: Any = db.fetch_one(t, query, param)  # <await>
        except Exception as e:
            result = e

        # A value of type `T` should always be castable to type `T`.
        if t is type(scalar):
            assert result == scalar
        else:
            match scalar:
```

Why this is a false positive: the handler stores the exception in `result`, which the test subsequently asserts against (`assert result == scalar` / type checks below); the failure is captured as test evidence, not hidden.

Checklist evidence: BP-PY-1's "broad except hides failures" predicate is not satisfied — the exception is recorded into the asserted result variable.

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none (single-reviewer audit)
- Chunk evidence: `scripts/onlymaps/chunks`
- Function evidence: `scripts/onlymaps/findings/functions`
- Validation: `git diff --check` — pass

## Post-fix v2 audit (latest binary)

### Run metadata

```yaml
timestamp: 2026-08-02T18:05:00Z
repository: onlymaps
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps
branch: main
commit: 6444a59
scanner_commit: 3407305 (binary rebuilt 2026-08-02 17:56)
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/onlymaps
chunk_path: scripts/onlymaps/chunks (exported 17:58)
function_context_path: scripts/onlymaps/findings/functions (60 contexts)
```

### Scan evidence

- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/onlymaps/chunks -context-dir scripts/onlymaps/findings/functions real-repos/onlymaps`
- Findings: `60` (unchanged from the b5b8fde run)
- Chunks reviewed: `scripts/onlymaps/chunks/Chunk_1_25.txt`, `Chunk_26_50.txt`, `Chunk_51_60.txt`
- Matching: every fresh `Source:` (file:line:col) matched an audited TP or FP entry from the Mode A/B audit above; no fresh classification decisions were needed beyond reuse.

### Classification summary (fresh run)

| Classification | Count | Fresh finding IDs |
| --- | ---: | --- |
| False positive | 24 | 2, 3, 4, 5, 6, 15, 16, 17, 18, 19, 24, 25, 26, 27, 28, 30, 32, 33, 34, 35, 36, 41, 57, 59 |
| True positive | 36 | 1, 7, 8, 9, 10, 11, 12, 13, 14, 20, 21, 22, 23, 29, 31, 37, 38, 39, 40, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 58, 60 |
| Uncertain | 0 | — |

New findings (never classified before): none.

## Fix checklist (FP patterns)

| Pattern # | Rule | Trigger shape | Count | Example sources |
| --- | --- | --- | ---: | --- |
| 1 | BP-PY-12, CWE-94 | `exec` token inside a `def`/`@overload` signature (method named `exec`, incl. type-only overloads) flagged as code-exec sink | 9 | `_connection.py:232`, `_pool.py:280`, `_query.py:52`, `_spec.py:153`, `_spec.py:170` |
| 2 | BP-PY-1 | bare `except:` whose suite ends with bare `raise` (rollback/counter/close first) flagged as swallowing | 5 | `_connection.py:215`, `_connection.py:376`, `_connection.py:408`, `_pool.py:215`, `_pool.py:443` |
| 3 | BP-PY-46 | `print` in `tasks.py` (invoke `@task` CLI runner) flagged as library logging | 5 | `tasks.py:6`, `tasks.py:12`, `tasks.py:20`, `tasks.py:23`, `tasks.py:25` |
| 4 | BP-PY-1 | `except Exception as exc:` in a test that assigns the exception to a result variable later asserted — flagged as hiding failures | 2 | `tests/test_pool.py:63`, `tests/test_query.py:78` |
| 5 | CWE-478 | `match` with no `case _` but an unconditional fallback statement immediately following the match | 2 | `_drivers.py:403` (`return None`), `_types.py:307` (`return value`) |
| 6 | CWE-89 | `conn.execute(<module-level SQL constant attr>)` — non-dynamic first argument flagged as injection | 1 | `tests/fixtures/containers.py:180` |

Distinguishing condition per pattern:
1. Only flag `exec`/`eval` when it is an actual call expression whose callee is the builtin — never a `def`/`@overload` parameter list.
2. A bare/broad except is not "swallowing" when its suite's final statement is `raise` (failure propagates).
3. Treat invoke-style CLI task files (all module functions `@task`-decorated) as CLI/script modules, like the existing `if __name__ == "__main__"` exemption.
4. In test files, a broad except that binds the exception into a name later asserted is evidence capture, not swallowing.
5. CWE-478: treat a match as having a default when an unconditional return/fallback statement immediately follows it (no intervening control flow).
6. CWE-89: only flag when the `execute`/`executemany` first argument is dynamic (interpolation, concatenation, or non-constant expression).

## New findings

None — every fresh finding matched a prior audited classification; no new sources, no Uncertain.

## Final evidence

- Delegated reviewers: none (single-reviewer audit)
- Validation: `git diff --check` — pass
