# False-positive audit — wse

## Run metadata

```yaml
timestamp: 2026-08-02T07:54:19Z
repository: wse
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse
branch: main
commit: 90ebca1d1dfb9b9ae9efc03f55738bb6a3b42444
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse
chunk_path: scripts/wse/chunks
function_context_path: scripts/wse/findings/functions
```

## Scan evidence

- Build command: `make build` (goslop binary at `./bin/goslop`)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/wse/chunks -context-dir scripts/wse/findings/functions real-repos/wse`
- Findings: `206`
- Chunks reviewed: `scripts/wse/chunks/Chunk_1_25.txt`, `Chunk_26_50.txt`, `Chunk_51_75.txt`, `Chunk_76_100.txt`, `Chunk_101_125.txt`, `Chunk_126_150.txt`, `Chunk_151_175.txt`, `Chunk_176_200.txt`, `Chunk_201_206.txt`
- Function contexts reviewed: `scripts/wse/findings/functions/<id>.txt` for every finding; source files read at the `Source:` path for every proposed false positive.

## Audit checklist

- [x] Read every assigned chunk under `scripts/wse/chunks`.
- [x] Read `scripts/wse/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient (no delegated reviews; no disagreements).
- [x] Ran `git diff --check` after updating this report.

Rule conditions were taken from `ruleset/python/bad-practices.json` (detection notes) and the detector implementations in `internal/lang/python/detectors/` (the `-explain` flag only covers the default Go/pack rules; the Python rules are registered only in Python scans).

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 49 | 1, 2, 3, 19, 41, 42, 43, 48, 72, 73, 74, 75, 77, 78, 79, 80, 81, 82, 84, 85, 86, 87, 91, 92, 94, 95, 96, 97, 98, 99, 100, 101, 102, 104, 105, 107, 108, 110, 111, 112, 114, 118, 119, 120, 122, 128, 172, 197, 198 |
| True positive | 157 | 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 44, 45, 46, 47, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 76, 83, 88, 89, 90, 93, 103, 106, 109, 113, 115, 116, 117, 121, 123, 124, 125, 126, 127, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 173, 174, 175, 176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194, 195, 196, 199, 200, 201, 202, 203, 204, 205, 206 |
| Uncertain | 0 | — |

## False positives

### [ ] Finding 1 — BP-PY-13

- Function context: `scripts/wse/findings/functions/1.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/benchmarks/bench_battle_server.py:55:1`
- Checklist pattern: rule excludes test fixtures and obvious placeholders (BP-PY-13 detection notes)

Source excerpt:

```
JWT_SECRET = b"bench-secret-key-for-testing-only"
JWT_ISSUER = "wse-bench"
```

Why this is a false positive: the value is a self-describing test fixture ("for-testing-only") in a benchmark harness, exactly the "test fixtures and obvious placeholders" case the rule documents as excluded.

Checklist evidence: the assignment is a pure string literal to `JWT_SECRET`, but the value and file location match the rule's documented exclusion for test fixtures/obvious placeholders.

### [ ] Finding 2 — CWE-88

- Function context: `scripts/wse/findings/functions/2.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/benchmarks/bench_battle_server.py:86:5`
- Checklist pattern: interpolated argv segments derive from a module constant, not from any input

Source excerpt:

```
    ca_key = f"{TLS_CERT_DIR}/ca.key"
    ca_crt = f"{TLS_CERT_DIR}/ca.crt"
    ...
    subprocess.run([
        "openssl", "req", "-x509", "-newkey", "rsa:2048",
        "-keyout", ca_key, "-out", ca_crt,
        "-days", "1", "-nodes", "-subj", "/CN=wse-test-ca",
    ], check=True, capture_output=True)
```

Why this is a false positive: the only non-literal argv elements are f-strings over `TLS_CERT_DIR` (a module constant `"/tmp/wse_tls_test"`); no input can reach the argument vector, so no untrusted argument can become an unintended option.

Checklist evidence: every argv element is a literal or an f-string of a fixed module-level constant; no user-controlled value exists in this function.

### [ ] Finding 3 — BP-PY-40

- Function context: `scripts/wse/findings/functions/3.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/benchmarks/bench_battle_server.py:129:6`
- Checklist pattern: rule is a heuristic for non-daemon workers; thread is daemon

Source excerpt:

```
    t = threading.Thread(target=httpd.serve_forever, daemon=True)
    t.start()
```

Why this is a false positive: the thread is created with an explicit `daemon=True` policy, which the rule's condition ("review-only heuristic for non-daemon workers"; skip daemon-only policy) intends to exempt; the detector only checks the `.start()` line, which is on the next line.

Checklist evidence: the `Thread` construction immediately above the flagged `.start()` declares `daemon=True`.

### [ ] Finding 19 — BP-PY-42

- Function context: `scripts/wse/findings/functions/19.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/benchmarks/bench_brutal.py:224:1`
- Checklist pattern: not a failure-expectation test; benchmark measurement helper

Source excerpt:

```
            async def _connect_timed():
                t0 = time.perf_counter()
                try:
                    ws = await connect_one(token)
                    lat = (time.perf_counter() - t0) * 1000
                    return ws, lat
                except Exception:
                    return None, None
```

Why this is a false positive: the try/except is a benchmark harness helper that counts connection failures (`errors += 1`); it is not a test expecting failure, so the `pytest.raises` recommendation does not apply.

Checklist evidence: the exception path returns a sentinel that the caller tallies as a measurement, not an assertion of expected failure.

### [ ] Finding 41 — BP-PY-13

- Function context: `scripts/wse/findings/functions/41.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/benchmarks/bench_fanout_server.py:42:1`
- Checklist pattern: rule excludes test fixtures and obvious placeholders

Source excerpt:

```
JWT_SECRET = b"bench-secret-key-for-testing-only"
JWT_ISSUER = "wse-bench"
```

Why this is a false positive: identical placeholder fixture pattern to finding 1, in a benchmark server; the value is self-describing test-only material.

Checklist evidence: assignment to a secret-named variable whose literal is an explicit testing placeholder in a benchmark script.

### [ ] Finding 42 — CWE-88

- Function context: `scripts/wse/findings/functions/42.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/benchmarks/bench_fanout_server.py:60:5`
- Checklist pattern: interpolated argv segments derive from a module constant, not from any input

Source excerpt:

```
    ca_key = f"{TLS_CERT_DIR}/ca.key"
    ca_crt = f"{TLS_CERT_DIR}/ca.crt"
    ...
    subprocess.run([
        "openssl", "req", "-x509", "-newkey", "rsa:2048",
        "-keyout", ca_key, "-out", ca_crt,
        "-days", "1", "-nodes", "-subj", "/CN=wse-test-ca",
    ], check=True, capture_output=True)
```

Why this is a false positive: same construct as finding 2 — f-string arguments interpolate only the fixed module constant `TLS_CERT_DIR = "/tmp/wse_tls_test"`; there is no untrusted input to inject an option.

Checklist evidence: no runtime/input-derived value appears in the argv; all dynamic segments reference a module-level constant.

### [ ] Finding 43 — BP-PY-40

- Function context: `scripts/wse/findings/functions/43.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/benchmarks/bench_fanout_server.py:113:6`
- Checklist pattern: rule is a heuristic for non-daemon workers; thread is daemon

Source excerpt:

```
    t = threading.Thread(target=httpd.serve_forever, daemon=True)
    t.start()
```

Why this is a false positive: explicit `daemon=True` policy on the construction line; only the separate `.start()` line tripped the line-based heuristic.

Checklist evidence: the `Thread` construction immediately above the flagged `.start()` declares `daemon=True`.

### [ ] Finding 48 — BP-PY-13

- Function context: `scripts/wse/findings/functions/48.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/benchmarks/bench_server.py:21:1`
- Checklist pattern: rule excludes test fixtures and obvious placeholders

Source excerpt:

```
# JWT config — matches what the benchmark client will use
JWT_SECRET = b"bench-secret-key-for-testing-only"
JWT_ISSUER = "wse-bench"
```

Why this is a false positive: benchmark fixture secret, explicitly self-described as testing-only.

Checklist evidence: secret-named assignment whose literal is an explicit testing placeholder in a benchmark script.

### [ ] Finding 72 — BP-PY-13

- Function context: `scripts/wse/findings/functions/72.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_basic.py:20:1`
- Checklist pattern: rule excludes obvious placeholders (`changeme` family)

Source excerpt:

```
JWT_SECRET = b"change-me-to-a-secure-32-byte-key!"
JWT_ISSUER = "my-app"
```

Why this is a false positive: the literal is an obvious `change-me` placeholder — the canonical placeholder family the rule documents as excluded; the detector's narrow list only matches the exact token `change-me`, not this expanded form.

Checklist evidence: the value starts with the documented placeholder token "change-me".

### [ ] Finding 73 — BP-PY-46

- Function context: `scripts/wse/findings/functions/73.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_basic.py:67:17`
- Checklist pattern: script (CLI) output, not library logging

Source excerpt:

```
if __name__ == "__main__":
    main()
```

and (inside `event_loop`, called from `main()`):

```
            if event_type == "auth_connect":
                user_id = event[2]
                connections[conn_id] = user_id
                print(f"[+] {user_id} connected ({conn_id})")
```

Why this is a false positive: `standalone_basic.py` is a runnable demo script whose prints are user-facing CLI output on the `__main__`-guard call path; the rule targets print-based logging in importable library modules.

Checklist evidence: the module is a `__main__`-guarded script (fix text: "keep print under if __name__ == '__main__' for CLIs"), and the print is in a function invoked from `main()` inside that guard.

### [ ] Finding 74 — BP-PY-46

- Function context: `scripts/wse/findings/functions/74.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_basic.py:74:17`
- Checklist pattern: script (CLI) output, not library logging

Source excerpt:

```
            elif event_type == "msg":
                data = event[2]  # parsed dict
                user = connections.get(conn_id, "?")
                print(f"[msg] {user}: {data}")
```

Why this is a false positive: demo-script presentation output inside `event_loop`, invoked from `main()` under the `__main__` guard of the same standalone example script.

Checklist evidence: module is a `__main__`-guarded example script; print is on the CLI call path.

### [ ] Finding 75 — BP-PY-46

- Function context: `scripts/wse/findings/functions/75.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_basic.py:80:17`
- Checklist pattern: script (CLI) output, not library logging

Source excerpt:

```
            elif event_type == "disconnect":
                user = connections.pop(conn_id, "?")
                print(f"[-] {user} disconnected")
```

Why this is a false positive: demo-script presentation output in the same example script, on the `__main__`-guard call path.

Checklist evidence: module is a `__main__`-guarded example script; print is on the CLI call path.

### [ ] Finding 77 — BP-PY-40

- Function context: `scripts/wse/findings/functions/77.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_basic.py:111:11`
- Checklist pattern: rule is a heuristic for non-daemon workers; thread is daemon

Source excerpt:

```
    thread = threading.Thread(target=event_loop, args=(server, stop), daemon=True)
    thread.start()
```

Why this is a false positive: explicit `daemon=True` policy; `.start()` on the next line trips the line-scoped daemon check.

Checklist evidence: the `Thread` construction immediately above the flagged `.start()` declares `daemon=True`.

### [ ] Finding 78 — BP-PY-13

- Function context: `scripts/wse/findings/functions/78.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_broadcast.py:24:1`
- Checklist pattern: rule excludes obvious placeholders (`changeme` family)

Source excerpt:

```
JWT_SECRET = b"change-me-to-a-secure-32-byte-key!"
JWT_ISSUER = "my-app"
```

Why this is a false positive: obvious `change-me` placeholder literal in a demo example script.

Checklist evidence: the value starts with the documented placeholder token "change-me".

### [ ] Finding 79 — BP-PY-46

- Function context: `scripts/wse/findings/functions/79.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_broadcast.py:52:17`
- Checklist pattern: script (CLI) output, not library logging

Source excerpt:

```
            if event_type == "auth_connect":
                user_id = event[2]
                print(f"[+] {user_id} connected")
```

Why this is a false positive: demo-script presentation output in `event_loop`, on the `main()`/`__main__` call path of the example script.

Checklist evidence: module is a `__main__`-guarded example script; print is on the CLI call path.

### [ ] Finding 80 — BP-PY-46

- Function context: `scripts/wse/findings/functions/80.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_broadcast.py:63:25`
- Checklist pattern: script (CLI) output, not library logging

Source excerpt:

```
                    if topics:
                        server.subscribe_connection(conn_id, topics)
                        print(f"[sub] {conn_id} subscribed to {topics}")
```

Why this is a false positive: demo-script presentation output on the `__main__`-guard call path of the example script.

Checklist evidence: module is a `__main__`-guarded example script; print is on the CLI call path.

### [ ] Finding 81 — BP-PY-46

- Function context: `scripts/wse/findings/functions/81.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_broadcast.py:66:17`
- Checklist pattern: script (CLI) output, not library logging

Source excerpt:

```
            elif event_type == "disconnect":
                print(f"[-] {conn_id} disconnected")
```

Why this is a false positive: demo-script presentation output on the `__main__`-guard call path of the example script.

Checklist evidence: module is a `__main__`-guarded example script; print is on the CLI call path.

### [ ] Finding 82 — BP-PY-46

- Function context: `scripts/wse/findings/functions/82.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_broadcast.py:95:17`
- Checklist pattern: script (CLI) output, not library logging

Source excerpt:

```
            subs = server.get_topic_subscriber_count("prices")
            if conns > 0:
                print(f"[pub] {seq} ticks, {conns} connections, {subs} subscribed to prices")
```

Why this is a false positive: publisher presentation output in the example script, on the `__main__`-guard call path.

Checklist evidence: module is a `__main__`-guarded example script; print is on the CLI call path.

### [ ] Finding 84 — BP-PY-40

- Function context: `scripts/wse/findings/functions/84.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_broadcast.py:123:7`
- Checklist pattern: rule is a heuristic for non-daemon workers; threads are daemon

Source excerpt:

```
    t1 = threading.Thread(target=event_loop, args=(server, stop), daemon=True)
    t2 = threading.Thread(target=publisher, args=(server, stop), daemon=True)
    t1.start()
    t2.start()
```

Why this is a false positive: both threads carry explicit `daemon=True`; `.start()` lines are separate, tripping the line-scoped daemon check.

Checklist evidence: every `Thread` construction in the file declares `daemon=True`.

### [ ] Finding 85 — BP-PY-13

- Function context: `scripts/wse/findings/functions/85.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_cluster.py:25:1`
- Checklist pattern: rule excludes obvious placeholders (`changeme` family)

Source excerpt:

```
JWT_SECRET = b"change-me-to-a-secure-32-byte-key!"
JWT_ISSUER = "my-app"
```

Why this is a false positive: obvious `change-me` placeholder literal in a demo example script.

Checklist evidence: the value starts with the documented placeholder token "change-me".

### [ ] Finding 86 — BP-PY-46

- Function context: `scripts/wse/findings/functions/86.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_cluster.py:52:17`
- Checklist pattern: script (CLI) output, not library logging

Source excerpt:

```
            if event_type == "auth_connect":
                user_id = event[2]
                print(f"[+] {user_id} connected")
                server.subscribe_connection(conn_id, ["notifications"])
```

Why this is a false positive: demo-script presentation output on the `main()`/`__main__` call path of the example script.

Checklist evidence: module is a `__main__`-guarded example script; print is on the CLI call path.

### [ ] Finding 87 — BP-PY-46

- Function context: `scripts/wse/findings/functions/87.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_cluster.py:56:17`
- Checklist pattern: script (CLI) output, not library logging

Source excerpt:

```
            elif event_type == "disconnect":
                print(f"[-] {conn_id} disconnected")
```

Why this is a false positive: demo-script presentation output on the `__main__`-guard call path of the example script.

Checklist evidence: module is a `__main__`-guarded example script; print is on the CLI call path.

### [ ] Finding 91 — BP-PY-40

- Function context: `scripts/wse/findings/functions/91.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_cluster.py:110:7`
- Checklist pattern: rule is a heuristic for non-daemon workers; threads are daemon

Source excerpt:

```
    t1 = threading.Thread(target=event_loop, args=(server, stop), daemon=True)
    t2 = threading.Thread(target=publisher, args=(server, stop, args.port), daemon=True)
    t1.start()
    t2.start()
```

Why this is a false positive: both threads carry explicit `daemon=True`; `.start()` lines are separate, tripping the line-scoped daemon check.

Checklist evidence: every `Thread` construction in the file declares `daemon=True`.

### [ ] Finding 92 — BP-PY-13

- Function context: `scripts/wse/findings/functions/92.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_presence.py:20:1`
- Checklist pattern: rule excludes obvious placeholders (`changeme` family)

Source excerpt:

```
JWT_SECRET = b"change-me-to-a-secure-32-byte-key!"
JWT_ISSUER = "my-app"
```

Why this is a false positive: obvious `change-me` placeholder literal in a demo example script.

Checklist evidence: the value starts with the documented placeholder token "change-me".

### [ ] Finding 94 — BP-PY-46

- Function context: `scripts/wse/findings/functions/94.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_presence.py:61:17`
- Checklist pattern: script (CLI) output, not library logging

Source excerpt:

```
                user_id = event[2]
                connections[conn_id] = user_id
                print(f"[+] {user_id} connected ({conn_id})")
```

Why this is a false positive: demo-script presentation output on the `main()`/`__main__` call path of the example script.

Checklist evidence: module is a `__main__`-guarded example script; print is on the CLI call path.

### [ ] Finding 95 — BP-PY-46

- Function context: `scripts/wse/findings/functions/95.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_presence.py:75:17`
- Checklist pattern: script (CLI) output, not library logging

Source excerpt:

```
                print(
                    f"[presence] {payload['user_id']} joined "
```

Why this is a false positive: demo-script presentation output on the `__main__`-guard call path of the example script.

Checklist evidence: module is a `__main__`-guarded example script; print is on the CLI call path.

### [ ] Finding 96 — BP-PY-46

- Function context: `scripts/wse/findings/functions/96.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_presence.py:82:17`
- Checklist pattern: script (CLI) output, not library logging

Source excerpt:

```
                members = server.presence("lobby")
                print(f"  Members in lobby: {len(members)}")
```

Why this is a false positive: demo-script presentation output on the `__main__`-guard call path of the example script.

Checklist evidence: module is a `__main__`-guarded example script; print is on the CLI call path.

### [ ] Finding 97 — BP-PY-46

- Function context: `scripts/wse/findings/functions/97.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_presence.py:84:21`
- Checklist pattern: script (CLI) output, not library logging

Source excerpt:

```
                for uid, info in members.items():
                    print(f"    {uid}: {info['data']} ({info['connections']} conn)")
```

Why this is a false positive: demo-script presentation output on the `__main__`-guard call path of the example script.

Checklist evidence: module is a `__main__`-guarded example script; print is on the CLI call path.

### [ ] Finding 98 — BP-PY-46

- Function context: `scripts/wse/findings/functions/98.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_presence.py:88:17`
- Checklist pattern: script (CLI) output, not library logging

Source excerpt:

```
                print(
                    f"[presence] {payload['user_id']} left "
```

Why this is a false positive: demo-script presentation output on the `__main__`-guard call path of the example script.

Checklist evidence: module is a `__main__`-guarded example script; print is on the CLI call path.

### [ ] Finding 99 — BP-PY-46

- Function context: `scripts/wse/findings/functions/99.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_presence.py:95:17`
- Checklist pattern: script (CLI) output, not library logging

Source excerpt:

```
                print(
                    f"  Lobby: {stats['num_users']} users, "
```

Why this is a false positive: demo-script presentation output on the `__main__`-guard call path of the example script.

Checklist evidence: module is a `__main__`-guarded example script; print is on the CLI call path.

### [ ] Finding 100 — BP-PY-46

- Function context: `scripts/wse/findings/functions/100.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_presence.py:104:17`
- Checklist pattern: script (CLI) output, not library logging

Source excerpt:

```
                data = event[2]
                user = connections.get(conn_id, "?")
                print(f"[msg] {user}: {data}")
```

Why this is a false positive: demo-script presentation output on the `__main__`-guard call path of the example script.

Checklist evidence: module is a `__main__`-guarded example script; print is on the CLI call path.

### [ ] Finding 101 — BP-PY-46

- Function context: `scripts/wse/findings/functions/101.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_presence.py:110:21`
- Checklist pattern: script (CLI) output, not library logging

Source excerpt:

```
                    server.update_presence(conn_id, {"status": new_status, "name": user})
                    print(f"  Updated {user} status to {new_status}")
```

Why this is a false positive: demo-script presentation output on the `__main__`-guard call path of the example script.

Checklist evidence: module is a `__main__`-guarded example script; print is on the CLI call path.

### [ ] Finding 102 — BP-PY-46

- Function context: `scripts/wse/findings/functions/102.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_presence.py:115:17`
- Checklist pattern: script (CLI) output, not library logging

Source excerpt:

```
                user = connections.pop(conn_id, "?")
                print(f"[-] {user} disconnected")
```

Why this is a false positive: demo-script presentation output on the `__main__`-guard call path of the example script.

Checklist evidence: module is a `__main__`-guarded example script; print is on the CLI call path.

### [ ] Finding 104 — BP-PY-40

- Function context: `scripts/wse/findings/functions/104.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_presence.py:145:11`
- Checklist pattern: rule is a heuristic for non-daemon workers; thread is daemon

Source excerpt:

```
    thread = threading.Thread(target=event_loop, args=(server, stop), daemon=True)
    thread.start()
```

Why this is a false positive: explicit `daemon=True` policy; `.start()` on the next line trips the line-scoped daemon check.

Checklist evidence: the `Thread` construction immediately above the flagged `.start()` declares `daemon=True`.

### [ ] Finding 105 — BP-PY-13

- Function context: `scripts/wse/findings/functions/105.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_recovery.py:26:1`
- Checklist pattern: rule excludes obvious placeholders (`changeme` family)

Source excerpt:

```
JWT_SECRET = b"change-me-to-a-secure-32-byte-key!"
JWT_ISSUER = "my-app"
```

Why this is a false positive: obvious `change-me` placeholder literal in a demo example script.

Checklist evidence: the value starts with the documented placeholder token "change-me".

### [ ] Finding 107 — BP-PY-46

- Function context: `scripts/wse/findings/functions/107.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_recovery.py:54:17`
- Checklist pattern: script (CLI) output, not library logging

Source excerpt:

```
            if event_type == "auth_connect":
                user_id = event[2]
                print(f"[+] {user_id} connected ({conn_id})")
```

Why this is a false positive: demo-script presentation output on the `main()`/`__main__` call path of the example script.

Checklist evidence: module is a `__main__`-guarded example script; print is on the CLI call path.

### [ ] Finding 108 — BP-PY-46

- Function context: `scripts/wse/findings/functions/108.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_recovery.py:59:17`
- Checklist pattern: script (CLI) output, not library logging

Source excerpt:

```
                print(f"    subscribe result: {result}")
```

Why this is a false positive: demo-script presentation output on the `__main__`-guard call path of the example script.

Checklist evidence: module is a `__main__`-guarded example script; print is on the CLI call path.

### [ ] Finding 110 — BP-PY-46

- Function context: `scripts/wse/findings/functions/110.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_recovery.py:117:29`
- Checklist pattern: script (CLI) output, not library logging

Source excerpt:

```
                            server.send(conn_id, response)
                            print(f"    [{conn_id}] subscribe {topics} recover={recover} -> {result}")
```

Why this is a false positive: demo-script presentation output on the `__main__`-guard call path of the example script.

Checklist evidence: module is a `__main__`-guarded example script; print is on the CLI call path.

### [ ] Finding 111 — BP-PY-46

- Function context: `scripts/wse/findings/functions/111.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_recovery.py:120:17`
- Checklist pattern: script (CLI) output, not library logging

Source excerpt:

```
            elif event_type == "disconnect":
                print(f"[-] {conn_id} disconnected")
```

Why this is a false positive: demo-script presentation output on the `__main__`-guard call path of the example script.

Checklist evidence: module is a `__main__`-guarded example script; print is on the CLI call path.

### [ ] Finding 112 — BP-PY-46

- Function context: `scripts/wse/findings/functions/112.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_recovery.py:146:13`
- Checklist pattern: script (CLI) output, not library logging

Source excerpt:

```
            print(
                f"[pub] {seq} ticks | {conns} conns | "
```

Why this is a false positive: publisher presentation output in the example script, on the `main()`/`__main__` call path.

Checklist evidence: module is a `__main__`-guarded example script; print is on the CLI call path.

### [ ] Finding 114 — BP-PY-40

- Function context: `scripts/wse/findings/functions/114.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_recovery.py:190:7`
- Checklist pattern: rule is a heuristic for non-daemon workers; threads are daemon

Source excerpt:

```
    t1 = threading.Thread(target=event_loop, args=(server, stop), daemon=True)
    t2 = threading.Thread(target=publisher, args=(server, stop), daemon=True)
    t1.start()
    t2.start()
```

Why this is a false positive: both threads carry explicit `daemon=True`; `.start()` lines are separate, tripping the line-scoped daemon check.

Checklist evidence: every `Thread` construction in the file declares `daemon=True`.

### [ ] Finding 118 — BP-PY-46

- Function context: `scripts/wse/findings/functions/118.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/python-client/wse_client/__init__.py:10:13`
- Checklist pattern: print is inside a docstring, not executable code

Source excerpt:

```
    async with connect("ws://localhost:5007/wse", token="your-jwt") as client:
        await client.subscribe(["notifications"])
        async for event in client:
            print(event.type, event.payload)
```

Why this is a false positive: the "print" is documentation example code inside the module docstring (`"""..."""`); the library contains no executable `print`.

Checklist evidence: the flagged line is inside a triple-quoted module docstring; the detector's per-line string check does not track docstring spans.

### [ ] Finding 119 — BP-PY-46

- Function context: `scripts/wse/findings/functions/119.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/python-client/wse_client/__init__.py:80:17`
- Checklist pattern: print is inside a docstring, not executable code

Source excerpt:

```
    Example::

        async with connect("ws://localhost:5007/wse", token="jwt") as client:
            await client.subscribe(["notifications"])
            async for event in client:
                print(event.type, event.payload)
    """
```

Why this is a false positive: documentation example inside the `connect()` docstring, not executable code.

Checklist evidence: the flagged line is inside a triple-quoted docstring of a library function.

### [ ] Finding 120 — CWE-89

- Function context: `scripts/wse/findings/functions/120.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/python-client/wse_client/circuit_breaker.py:77:15`
- Checklist pattern: no SQL exists; the trigger is a method named `execute`

Source excerpt:

```
    async def execute(self, fn: Callable[[], Awaitable[T]]) -> T:
        if not self.can_execute():
            raise WSECircuitBreakerError("Circuit breaker is open")
        try:
            result = await fn()
            self.record_success()
            return result
        except Exception:
            self.record_failure()
            raise
```

Why this is a false positive: the match is the bare `execute(` in a method definition of a circuit breaker; the file contains no SQL statement, cursor, or database call, so no dynamic SQL reaches any `execute`/`executemany`.

Checklist evidence: the first "argument" of the matched call is `self` (a method definition), and no SQL string exists in the unit.

### [ ] Finding 122 — BP-PY-46

- Function context: `scripts/wse/findings/functions/122.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/python-client/wse_client/client.py:90:17`
- Checklist pattern: print is inside a docstring, not executable code

Source excerpt:

```
    Example::

        async with AsyncWSEClient("ws://localhost:5006/wse", token="jwt") as client:
            await client.subscribe(["notifications"])
            async for event in client:
                print(event.type, event.payload)
    """
```

Why this is a false positive: documentation example inside the `AsyncWSEClient` docstring, not executable code.

Checklist evidence: the flagged line is inside a triple-quoted docstring of a library class.

### [ ] Finding 128 — BP-PY-46

- Function context: `scripts/wse/findings/functions/128.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/python-client/wse_client/client.py:663:17`
- Checklist pattern: print is inside a docstring, not executable code

Source excerpt:

```
        Example::

            @client.on("notifications")
            async def handle(event: WSEEvent):
                print(event.payload)
        """
```

Why this is a false positive: documentation example inside a decorator's docstring, not executable code.

Checklist evidence: the flagged line is inside a triple-quoted docstring of a library method.

### [ ] Finding 172 — BP-PY-46

- Function context: `scripts/wse/findings/functions/172.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/python-client/wse_client/sync_client.py:53:13`
- Checklist pattern: print is inside a docstring, not executable code

Source excerpt:

```
    Example (callbacks)::

        @client.on("notifications")
        def handle(event):
            print(event.payload)

        client.run_forever()
    """
```

Why this is a false positive: documentation example inside a class docstring, not executable code.

Checklist evidence: the flagged line is inside a triple-quoted docstring of a library class.

### [ ] Finding 197 — BP-PY-13

- Function context: `scripts/wse/findings/functions/197.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/tests/conftest.py:13:1`
- Checklist pattern: rule excludes test fixtures

Source excerpt:

```
# Shared JWT config
JWT_SECRET = b"test-secret-key-for-integration!"
JWT_ISSUER = "wse-test"
```

Why this is a false positive: a test-fixture secret in `tests/conftest.py`, exactly the "test fixtures" case the rule documents as excluded.

Checklist evidence: the value is an explicit test fixture literal (`test-...-for-integration`) in the tests directory.

### [ ] Finding 198 — CWE-1341

- Function context: `scripts/wse/findings/functions/198.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/tests/test_integration.py:134:15`
- Checklist pattern: two distinct handles, each closed once

Source excerpt:

```
        await ws1.close()
        drain_until(server, "disconnect")
        assert server.get_connection_count() == 1

        await ws2.close()
        drain_until(server, "disconnect")
        assert server.get_connection_count() == 0
```

Why this is a false positive: the rule's regex (`\w+\.close(...)` twice within 180 chars) matched `ws1.close()` followed by `ws2.close()` — two different connection handles, each released exactly once; no handle is released twice.

Checklist evidence: the two matched `.close()` calls act on different identifiers (`ws1`, `ws2`), so the "same resource handle released twice" condition is not met.

## True positives

All 157 true positives satisfy their rule condition on the shown source; grouped per rule below.

### BP-PY-1 — broad `except Exception` without re-raise (59)

| Finding | Source | Reason |
| --- | --- | --- |
| 6 | benchmarks/bench_battle_server.py:172 | `except Exception: pass` swallows update_presence failures |
| 11 | benchmarks/bench_battle_server.py:272 | broad except, no re-raise, handler is pass |
| 13 | benchmarks/bench_battle_server.py:283 | broad except with fallback assignment |
| 14 | benchmarks/bench_battle_server.py:299 | broad except, only prints |
| 15 | benchmarks/bench_battle_server.py:339 | broad except, only sleeps |
| 20 | benchmarks/bench_brutal.py:228 | broad except returning sentinel, no re-raise |
| 23 | benchmarks/bench_brutal.py:288 | broad except breaking loop |
| 24 | benchmarks/bench_brutal.py:329 | broad except reporting error, no re-raise |
| 25 | benchmarks/bench_brutal.py:348 | broad except building error string |
| 26 | benchmarks/bench_brutal.py:354 | `except Exception: pass` on ws.close |
| 31 | benchmarks/bench_brutal.py:418 | broad except, pass |
| 33 | benchmarks/bench_brutal.py:430 | broad except, pass |
| 37 | benchmarks/bench_brutal.py:653 | broad except returning early |
| 44 | benchmarks/bench_fanout_server.py:132 | broad except, only prints |
| 46 | benchmarks/bench_fanout_server.py:161 | broad except, only sleeps |
| 50 | benchmarks/bench_server.py:69 | broad except, only prints |
| 52 | benchmarks/bench_wse.py:109 | broad except, only prints |
| 54 | benchmarks/bench_wse.py:147 | broad except, only prints |
| 55 | benchmarks/bench_wse.py:192 | broad except, only prints |
| 56 | benchmarks/bench_wse.py:216 | broad except, only prints |
| 57 | benchmarks/bench_wse.py:267 | broad except, only prints |
| 65 | benchmarks/bench_wse_multiprocess.py:197 | broad except, pass |
| 67 | benchmarks/bench_wse_multiprocess.py:259 | broad except, pass |
| 69 | benchmarks/bench_wse_multiprocess.py:353 | broad except reporting error |
| 88 | examples/standalone_cluster.py:94 | broad except, only prints |
| 129 | python-client/wse_client/client.py:801 | broad except, logs at debug |
| 132 | python-client/wse_client/client.py:846 | broad except logging handler error |
| 135 | python-client/wse_client/connection.py:232 | broad except, pass |
| 139 | python-client/wse_client/connection.py:242 | broad except, pass |
| 141 | python-client/wse_client/connection.py:287 | broad except, pass |
| 143 | python-client/wse_client/connection.py:293 | broad except, pass |
| 145 | python-client/wse_client/connection.py:329 | broad except, only logs |
| 146 | python-client/wse_client/connection.py:352 | broad except returning False |
| 147 | python-client/wse_client/connection.py:380 | broad except, pass |
| 149 | python-client/wse_client/connection.py:389 | broad except, only logs |
| 151 | python-client/wse_client/connection.py:476 | broad except, fallback timestamp |
| 152 | python-client/wse_client/connection.py:683 | broad except, only logs |
| 153 | python-client/wse_client/connection.py:720 | broad except, pass |
| 155 | python-client/wse_client/connection.py:727 | broad except, pass |
| 157 | python-client/wse_client/protocol.py:200 | broad except, logs warning, no re-raise |
| 159 | python-client/wse_client/protocol.py:229 | broad except, logs warning, no re-raise |
| 160 | python-client/wse_client/protocol.py:240 | broad except, logs warning, no re-raise |
| 164 | python-client/wse_client/protocol.py:264 | broad except, pass |
| 167 | python-client/wse_client/security.py:185 | broad except returning False |
| 171 | python-client/wse_client/security.py:250 | broad except, failure counter |
| 173 | python-client/wse_client/sync_client.py:119 | broad except, pass |
| 178 | python-client/wse_client/sync_client.py:154 | broad except, pass |
| 181 | python-client/wse_client/sync_client.py:186 | broad except returning False |
| 182 | python-client/wse_client/sync_client.py:206 | broad except returning False |
| 183 | python-client/wse_client/sync_client.py:225 | broad except returning False |
| 184 | python-client/wse_client/sync_client.py:247 | broad except returning False |
| 185 | python-client/wse_client/sync_client.py:259 | broad except returning False |
| 186 | python-client/wse_client/sync_client.py:278 | broad except returning False |
| 187 | python-client/wse_client/sync_client.py:290 | broad except, pass |
| 189 | python-client/wse_client/sync_client.py:302 | broad except, pass |
| 191 | python-client/wse_client/sync_client.py:447 | broad except, only logs |
| 193 | python-client/wse_client/sync_client.py:485 | broad except storing error |
| 194 | python-client/wse_client/sync_client.py:492 | broad except, pass |
| 196 | python-client/wse_client/sync_client.py:504 | broad except, only logs |

### BP-PY-2 — except handler body is solely `pass` (33)

| Finding | Source | Reason |
| --- | --- | --- |
| 7 | benchmarks/bench_battle_server.py:172 | suite is only `pass` |
| 12 | benchmarks/bench_battle_server.py:272 | suite is only `pass` |
| 17 | benchmarks/bench_battle_server.py:489 | `except RuntimeError: pass` |
| 27 | benchmarks/bench_brutal.py:354 | suite is only `pass` |
| 32 | benchmarks/bench_brutal.py:418 | suite is only `pass` |
| 34 | benchmarks/bench_brutal.py:430 | suite is only `pass` |
| 60 | benchmarks/bench_wse_multiprocess.py:108 | `except TimeoutError: pass` |
| 63 | benchmarks/bench_wse_multiprocess.py:121 | `except TimeoutError: pass` |
| 66 | benchmarks/bench_wse_multiprocess.py:197 | suite is only `pass` |
| 68 | benchmarks/bench_wse_multiprocess.py:259 | suite is only `pass` |
| 123 | python-client/wse_client/client.py:126 | `except ImportError: pass` |
| 126 | python-client/wse_client/client.py:296 | `except (QueueEmpty, QueueFull): pass` |
| 127 | python-client/wse_client/client.py:613 | `except QueueFull: pass` |
| 131 | python-client/wse_client/client.py:829 | `except (QueueEmpty, QueueFull): pass` |
| 133 | python-client/wse_client/client.py:1211 | `except (QueueEmpty, QueueFull): pass` |
| 136 | python-client/wse_client/connection.py:232 | suite is only `pass` |
| 140 | python-client/wse_client/connection.py:242 | suite is only `pass` |
| 142 | python-client/wse_client/connection.py:287 | suite is only `pass` |
| 144 | python-client/wse_client/connection.py:293 | suite is only `pass` |
| 148 | python-client/wse_client/connection.py:380 | suite is only `pass` |
| 150 | python-client/wse_client/connection.py:459 | `except (ValueError, TypeError): pass` |
| 154 | python-client/wse_client/connection.py:720 | suite is only `pass` |
| 156 | python-client/wse_client/connection.py:727 | suite is only `pass` |
| 161 | python-client/wse_client/protocol.py:256 | `except UnicodeDecodeError: pass` |
| 165 | python-client/wse_client/protocol.py:264 | suite is only `pass` |
| 168 | python-client/wse_client/security.py:224 | `except RuntimeError: pass` (comment-only suite) |
| 174 | python-client/wse_client/sync_client.py:119 | suite is only `pass` |
| 179 | python-client/wse_client/sync_client.py:154 | suite is only `pass` |
| 180 | python-client/wse_client/sync_client.py:159 | `except queue.Full: pass` |
| 188 | python-client/wse_client/sync_client.py:290 | suite is only `pass` |
| 190 | python-client/wse_client/sync_client.py:302 | suite is only `pass` |
| 192 | python-client/wse_client/sync_client.py:482 | `except (Empty, Full): pass` |
| 195 | python-client/wse_client/sync_client.py:492 | suite is only `pass` |

### CWE-390 — error condition detected without action (8)

| Finding | Source | Reason |
| --- | --- | --- |
| 8 | benchmarks/bench_battle_server.py:172 | handler takes no action |
| 28 | benchmarks/bench_brutal.py:354 | handler takes no action |
| 61 | benchmarks/bench_wse_multiprocess.py:108 | `TimeoutError: pass`, no action |
| 124 | python-client/wse_client/client.py:126 | `ImportError: pass`, no action |
| 137 | python-client/wse_client/connection.py:232 | handler takes no action |
| 162 | python-client/wse_client/protocol.py:256 | `UnicodeDecodeError: pass`, no action |
| 169 | python-client/wse_client/security.py:224 | `RuntimeError: pass`, no action |
| 175 | python-client/wse_client/sync_client.py:119 | handler takes no action |

### CWE-396 — generic Exception catch (13)

| Finding | Source | Reason |
| --- | --- | --- |
| 9 | benchmarks/bench_battle_server.py:172 | generic `except Exception` |
| 21 | benchmarks/bench_brutal.py:228 | generic `except Exception` |
| 45 | benchmarks/bench_fanout_server.py:132 | generic `except Exception` |
| 51 | benchmarks/bench_server.py:69 | generic `except Exception` |
| 53 | benchmarks/bench_wse.py:109 | generic `except Exception` |
| 59 | benchmarks/bench_wse_multiprocess.py:96 | generic `except Exception` |
| 89 | examples/standalone_cluster.py:94 | generic `except Exception` |
| 121 | python-client/wse_client/circuit_breaker.py:84 | generic `except Exception` |
| 130 | python-client/wse_client/client.py:801 | generic `except Exception` |
| 134 | python-client/wse_client/connection.py:193 | generic `except Exception` |
| 158 | python-client/wse_client/protocol.py:200 | generic `except Exception` |
| 166 | python-client/wse_client/security.py:147 | generic `except Exception` |
| 176 | python-client/wse_client/sync_client.py:119 | generic `except Exception` |

### CWE-1071 — empty code block (8)

| Finding | Source | Reason |
| --- | --- | --- |
| 10 | benchmarks/bench_battle_server.py:172 | handler contains only `pass` |
| 29 | benchmarks/bench_brutal.py:354 | handler contains only `pass` |
| 62 | benchmarks/bench_wse_multiprocess.py:108 | handler contains only `pass` |
| 125 | python-client/wse_client/client.py:126 | handler contains only `pass` |
| 138 | python-client/wse_client/connection.py:232 | handler contains only `pass` |
| 163 | python-client/wse_client/protocol.py:256 | handler contains only `pass` |
| 170 | python-client/wse_client/security.py:224 | handler contains only `pass` |
| 177 | python-client/wse_client/sync_client.py:119 | handler contains only `pass` |

### CWE-1121 — excessive cyclomatic complexity (5)

| Finding | Source | Reason |
| --- | --- | --- |
| 4 | benchmarks/bench_battle_server.py:134 (`drain_loop`) | 34 counted control-flow branches ≥ 12 |
| 30 | benchmarks/bench_brutal.py:362 (`_throughput_worker_async`) | 15 counted branches ≥ 12 |
| 70 | benchmarks/bench_wse_multiprocess.py:561 (`run_test_7`) | 13 counted branches ≥ 12 |
| 93 | examples/standalone_presence.py:40 (`event_loop`) | 13 counted branches ≥ 12 |
| 106 | examples/standalone_recovery.py:44 (`event_loop`) | 14 counted branches ≥ 12 |

### CWE-1124 — excessively deep nesting (3)

| Finding | Source | Reason |
| --- | --- | --- |
| 5 | benchmarks/bench_battle_server.py:146 | executable statement at 6 control-flow levels |
| 64 | benchmarks/bench_wse_multiprocess.py:195 | executable statement at 6 control-flow levels |
| 109 | examples/standalone_recovery.py:75 | executable statement at 6 control-flow levels |

### CWE-215 — sensitive value in debug output (10)

| Finding | Source | Reason |
| --- | --- | --- |
| 16 | benchmarks/bench_battle_server.py:455 | prints full JWT token |
| 47 | benchmarks/bench_fanout_server.py:250 | prints full JWT token |
| 49 | benchmarks/bench_server.py:52 | prints full JWT token |
| 58 | benchmarks/bench_wse.py:287 | prints token prefix |
| 71 | benchmarks/bench_wse_multiprocess.py:730 | prints token prefix |
| 76 | examples/standalone_basic.py:102 | prints token |
| 83 | examples/standalone_broadcast.py:113 | prints token |
| 90 | examples/standalone_cluster.py:101 | prints token |
| 103 | examples/standalone_presence.py:136 | prints token |
| 113 | examples/standalone_recovery.py:172 | prints token |

### BP-PY-41 — test function with side effects, no assertions (13)

| Finding | Source | Reason |
| --- | --- | --- |
| 18 | benchmarks/bench_brutal.py:201 (`test_connection_storm`) | `test_*` fn, calls, no assert |
| 22 | benchmarks/bench_brutal.py:265 (`test_echo_latency`) | `test_*` fn, calls, no assert |
| 35 | benchmarks/bench_brutal.py:605 (`test_throughput`) | `test_*` fn, calls, no assert |
| 36 | benchmarks/bench_brutal.py:631 (`test_sustained`) | `test_*` fn, calls, no assert |
| 38 | benchmarks/bench_brutal.py:698 (`test_json_comparison`) | `test_*` fn, calls, no assert |
| 39 | benchmarks/bench_brutal.py:772 (`test_message_sizes`) | `test_*` fn, calls, no assert |
| 40 | benchmarks/bench_brutal.py:809 (`test_format_throughput`) | `test_*` fn, calls, no assert |
| 115 | python-client/tests/test_client.py:177 | calls `_handle_error`, no assert |
| 116 | python-client/tests/test_client.py:185 | calls `_handle_error`, no assert |
| 117 | python-client/tests/test_sync_client.py:49 | calls `disconnect`, no assert |
| 200 | tests/test_integration.py:441 | calls `disconnect`, no assert |
| 201 | tests/test_integration.py:445 | calls `send`, no assert |
| 202 | tests/test_integration.py:449 | calls broadcast APIs, no assert |

### BP-PY-39 — `time.sleep` inside `async def` (5)

| Finding | Source | Reason |
| --- | --- | --- |
| 199 | tests/test_integration.py:327 | `time.sleep` in async test body |
| 203 | tests/test_integration.py:532 | `time.sleep` in async test body |
| 204 | tests/test_integration.py:539 | `time.sleep` in async test body |
| 205 | tests/test_integration.py:556 | `time.sleep` in async test body |
| 206 | tests/test_integration.py:818 | `time.sleep` in async test body |

## Uncertain findings

None — every finding could be classified from the rule condition and the shown source.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/wse/chunks`
- Function evidence: `scripts/wse/findings/functions`
- Validation: `git diff --check` — pass
## Post-fix remaining-FP audit (2026-08-02)

### Run metadata

```yaml
timestamp: 2026-08-02T11:08:00Z
repository: wse
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse
branch: main
commit: 90ebca1d1dfb9b9ae9efc03f55738bb6a3b42444 (unchanged since first audit)
goslop_binary: ./bin/goslop rebuilt 2026-08-02 16:29 IST from commit b5b8fde (FP-reduction fix)
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse
chunk_path: scripts/wse/chunks
function_context_path: scripts/wse/findings/functions
```

### Scan evidence

- Build command: `make build` (goslop binary at `./bin/goslop`)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/wse/chunks -context-dir scripts/wse/findings/functions real-repos/wse`
- Findings: `169` (previous run: 206)
- Chunks reviewed: `scripts/wse/chunks/Chunk_1_25.txt`, `Chunk_26_50.txt`, `Chunk_51_75.txt`, `Chunk_76_100.txt`, `Chunk_101_125.txt`, `Chunk_126_150.txt`, `Chunk_151_169.txt` (all 7)
- Function contexts reviewed: `scripts/wse/findings/functions/<id>.txt` for every proposed false positive (1, 17, 39, 67, 69, 70, 72, 73, 77, 78, 81, 82, 86) and for the three newly-appearing CWE-396 findings (61, 103, 133); enclosing source read for all of the above.

### Audit checklist

- [x] Read every assigned chunk under `scripts/wse/chunks`.
- [x] Read `scripts/wse/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient (no delegated reviews; no disagreements).
- [x] Ran `git diff --check` after updating this report.

### Classification summary

Fresh finding IDs are per the fresh scan; matching to the audited run is by `Source:` path (file:line).

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 13 | 1, 17, 39, 67, 69, 70, 72, 73, 77, 78, 81, 82, 86 |
| True positive | 156 | 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 68, 71, 74, 75, 76, 79, 80, 83, 84, 85, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169 |
| Uncertain | 0 | — |

Notes: the 156 true positives are the 153 fresh findings whose `Source:` matches an audited TP (all still present in the source, unchanged), plus 3 newly-appearing CWE-396 findings (61, 103, 133) at `bench_wse_multiprocess.py:197`, `connection.py:232`, `security.py:185` — generic `except Exception:` handlers without re-raise, which satisfy the CWE-396 rule condition (the fix dropped the re-raise variants at `bench_wse_multiprocess.py:96`, `circuit_breaker.py:84`, `connection.py:193`, `security.py:147`). The 13 remaining false positives all match audited FP sources; the fix removed the other 36 audited FPs (BP-PY-46 docstring/CLI prints, BP-PY-13 bench/test fixtures, CWE-88/89/1341, and the two remaining BP-PY-40 sites' siblings were not re-flagged).

## False positives

Two of the three rule types repeat the exact same source construct across files and are grouped; the third is a single finding. Each grouped finding ID is listed with its source.

### [ ] Findings 1, 39, 69, 72, 77, 81, 86 — BP-PY-40

- Function context: `./scripts/wse/findings/functions/<id>.txt` for 1, 39, 69, 72, 77, 81, 86
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/benchmarks/bench_battle_server.py:129:6` (finding 1), `bench_fanout_server.py:113:6` (finding 39), `examples/standalone_basic.py:111:11` (finding 69), `examples/standalone_broadcast.py:123:7` (finding 72), `examples/standalone_cluster.py:110:7` (finding 77), `examples/standalone_presence.py:145:11` (finding 81), `examples/standalone_recovery.py:190:7` (finding 86)
- Checklist pattern: rule is a review-only heuristic for non-daemon workers; every thread is created with an explicit `daemon=True` policy

Source excerpt (identical construct at every listed source; finding 1 shown):

```
    t = threading.Thread(target=httpd.serve_forever, daemon=True)
    t.start()
```

and (finding 72, two threads):

```
    t1 = threading.Thread(target=event_loop, args=(server, stop), daemon=True)
    t2 = threading.Thread(target=publisher, args=(server, stop), daemon=True)
    t1.start()
    t2.start()
```

Why these are false positives: each finding flags a `.start()` line whose `Thread` construction (the immediately preceding line) declares an explicit `daemon=True` policy, so none of the threads is a non-daemon fire-and-forget worker; the line-scoped detector trips on the separate `.start()` statement without consulting the constructor.

Checklist evidence: in every flagged function, the `threading.Thread(...)` call carries `daemon=True`; the rule's own detection notes scope it to non-daemon workers ("review-only heuristic for non-daemon workers").

### [ ] Findings 67, 70, 73, 78, 82 — BP-PY-13

- Function context: `./scripts/wse/findings/functions/<id>.txt` for 67, 70, 73, 78, 82
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/examples/standalone_basic.py:20:1` (finding 67), `standalone_broadcast.py:24:1` (finding 70), `standalone_cluster.py:25:1` (finding 73), `standalone_presence.py:20:1` (finding 78), `standalone_recovery.py:26:1` (finding 82)
- Checklist pattern: rule excludes obvious placeholders (changeme, xxx)

Source excerpt (identical construct at every listed source):

```
JWT_SECRET = b"change-me-to-a-secure-32-byte-key!"
JWT_ISSUER = "my-app"
```

Why these are false positives: the literal starts with the canonical `change-me` placeholder token, which the rule's detection notes explicitly exclude ("obvious placeholders (changeme, xxx)"); each is a demo-script fixture in `examples/`, not a real credential.

Checklist evidence: the flagged assignment is a secret-named variable whose value is an obvious `change-me` placeholder in an example script — exactly the documented exclusion case.

### [ ] Finding 17 — BP-PY-42

- Function context: `./scripts/wse/findings/functions/17.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/wse/benchmarks/bench_brutal.py:224:1`
- Checklist pattern: not a failure-expectation test; benchmark measurement helper

Source excerpt:

```
            async def _connect_timed():
                t0 = time.perf_counter()
                try:
                    ws = await connect_one(token)
                    lat = (time.perf_counter() - t0) * 1000
                    return ws, lat
                except Exception:
                    return None, None
```

Why this is a false positive: the try/except lives in a benchmark harness helper whose exception path returns a `(None, None)` sentinel that the caller tallies as a connection failure count; it is not a test expecting failure, so the `pytest.raises` recommendation does not apply.

Checklist evidence: no assertion or `pytest.raises` context is intended; the handler's purpose is measurement of failures, which the rule's condition ("try/except used to 'expect' failures" in tests) does not cover.

## Uncertain findings

None — every fresh finding could be classified from the rule condition and the shown source.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/wse/chunks`
- Function evidence: `scripts/wse/findings/functions`
- Validation: `git diff --check` — pass
