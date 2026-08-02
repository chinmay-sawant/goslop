# False-positive audit: tenso

## Run metadata

```yaml
timestamp: 2026-08-02T00:00:00Z
repository: tenso
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso
branch: main
commit: ee5d6eb7baba8aca90b1d63a5a176b0a7d37692e
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso
chunk_path: scripts/tenso/chunks
function_context_path: scripts/tenso/findings/functions
```

## Scan evidence

- Build command: `make build` (local `./bin/goslop` binary, see `Makefile`)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/tenso/chunks -context-dir scripts/tenso/findings/functions real-repos/tenso`
- Findings: `125`
- Chunks reviewed: `scripts/tenso/chunks/Chunk_1_25.txt`, `Chunk_26_50.txt`, `Chunk_51_75.txt`, `Chunk_76_100.txt`, `Chunk_101_125.txt`
- Function contexts reviewed: `scripts/tenso/findings/functions/1.txt`, `8.txt`, `9.txt`, `23.txt`, `29.txt`, `61.txt`, `80.txt`, `84.txt`, `92.txt`, `93.txt`, `98.txt`, `114.txt`, `115.txt`, `117.txt`, `118.txt`, `120.txt`, `125.txt`, plus the enclosing sources for all other findings

## Audit checklist

- [x] Read every assigned chunk under `scripts/tenso/chunks`.
- [x] Read `scripts/tenso/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 54 | 1, 2, 8, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 53, 54, 59, 60, 61, 80, 83, 84, 92, 93, 98, 104, 105, 114, 118, 125 |
| True positive | 71 | 3, 4, 5, 6, 7, 9, 10, 11, 12, 13, 14, 15, 16, 55, 56, 57, 58, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 81, 82, 85, 86, 87, 88, 89, 90, 91, 94, 95, 96, 97, 99, 100, 101, 102, 103, 106, 107, 108, 109, 110, 111, 112, 113, 115, 116, 117, 119, 120, 121, 122, 123, 124 |
| Uncertain | 0 | — |

## False positives

### [ ] Finding `1` — `BP-PY-10`

- Function context: `scripts/tenso/findings/functions/1.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/benchmark.py:111:21`
- Checklist pattern: rule condition not met — no untrusted-data sink

Source excerpt:

```
def bench_pickle(data):
    ...
    enc = lambda x: pickle.dumps(x, protocol=pickle.HIGHEST_PROTOCOL)  # noqa
    dec = lambda x: pickle.loads(x)  # noqa
    return enc, dec
```

Why this is a false positive: the decoder is only ever fed bytes produced by the encoder in the same process on `np.random`-generated data; the rule targets `pickle.loads` on attacker-controlled sources (request body, user-path file, cache), and none of the rule's non-constant source categories applies to a benchmark harness that round-trips its own data.

Checklist evidence: the source is a benchmark of serialization formats; `dec` has no call path from any external input in the shown source.

### [ ] Finding `2` — `CWE-502`

- Function context: `scripts/tenso/findings/functions/1.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/benchmark.py:111:21`
- Checklist pattern: rule condition not met — data is not untrusted

Source excerpt:

```
def bench_pickle(data):
    ...
    enc = lambda x: pickle.dumps(x, protocol=pickle.HIGHEST_PROTOCOL)  # noqa
    dec = lambda x: pickle.loads(x)  # noqa
    return enc, dec
```

Why this is a false positive: the detector flags every `pickle.loads` call ("any call is treated as unsafe … conservative"), but the rule condition is deserialization of *untrusted* data; here the unpickled bytes are produced and consumed within the same benchmark process.

Checklist evidence: the deserialized value flows from `enc` only; no trust boundary exists in the shown source.

### [ ] Finding `8` — `BP-PY-10`

- Function context: `scripts/tenso/findings/functions/8.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/benchmark.py:409:13`
- Checklist pattern: rule condition not met — self-written temp file, not user input

Source excerpt:

```
        t0 = time.perf_counter()
        with open(path, "wb") as f:
            pickle.dump(data, f)          # same path written moments earlier
        ...
        with open(path, "rb") as f:
            t0 = time.perf_counter()
            pickle.load(f)
```

Why this is a false positive: `path` is a `tempfile.NamedTemporaryFile` written by this same benchmark function immediately before the load; the file is not a user-path or request input, so the "file from user path" source category in the rule's detection notes does not apply.

Checklist evidence: the file handle passed to `pickle.load` is opened from a path created and written by the same function.

### [ ] Finding `17` — `BP-PY-45`

- Function context: `scripts/tenso/findings/functions/17.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/docs/source/conf.py:7:1`
- Checklist pattern: excluded construct — packaging/docs bootstrap

Source excerpt:

```
# If you installed the package in editable mode (`uv sync`), Sphinx can often find it.
# However, pointing explicitly to 'src' is the most reliable method.
sys.path.insert(0, os.path.abspath("../../src"))
```

Why this is a false positive: this is the Sphinx documentation build bootstrap (`docs/source/conf.py`), executed only at docs-build time; the rule condition explicitly excludes "packaging bootstrap" (`sys.path` mutation outside tests and packaging bootstrap).

Checklist evidence: the mutation is the standard docs-build path setup, not application runtime code.

### [ ] Findings `18`, `19`, `20`, `21`, `22` — `BP-PY-46` (examples/client.py)

- Function context: `scripts/tenso/findings/functions/18.txt` … `22.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/examples/client.py:13:1` (and 21, 26, 27, 29)
- Checklist pattern: rule condition not met — script module, not library code

Source excerpt:

```
# Connect
client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
client.connect(("localhost", 9999))
...
print(f"Sending Tensor: {data.shape} ({data.nbytes / 1024 / 1024:.2f} MB)")
...
print("Receiving response...")
...
if result is not None:
    print(f"Got Result in {time.time() - t0:.4f}s")
    print(f"Result Shape: {result.shape} | Mean: {result.mean():.4f}")
else:
    print("Server disconnected.")
```

Why this is a false positive: `examples/client.py` is a standalone demo script with top-level socket setup and no classes; the rule condition is "print used for operational logging in non-script modules" — these prints are the script's program output.

Checklist evidence: module is a runnable example script, not an importable library module; prints are not debugging leftovers.

### [ ] Finding `23` — `PERF-PY-28`

- Function context: `scripts/tenso/findings/functions/23.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/examples/grpc/bench_server.py:29:44`
- Checklist pattern: rule condition not met — executor is process-lifetime

Source excerpt:

```
def serve():
    ...
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=1), options=options)
    benchmark_msg_pb2_grpc.add_BenchmarkerServicer_to_server(Benchmarker(), server)
    server.add_insecure_port("[::]:50051")
    server.start()
    server.wait_for_termination()
```

Why this is a false positive: the executor is constructed once at server startup and handed to `grpc.server`, which owns it for the server's lifetime — it is not "created per unit of work" as the rule condition requires.

Checklist evidence: the construction site is the top-level `serve()` bootstrap, executed once per process.

### [ ] Findings `24`, `25`, `26`, `27` — `BP-PY-46` (examples/grpc/client_grpc.py)

- Function context: `scripts/tenso/findings/functions/24.txt` … `27.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/examples/grpc/client_grpc.py:21:9` (and 37, 38, 39)
- Checklist pattern: rule condition not met — example script CLI

Source excerpt:

```
def run():
    ...
    print(f"Sending tensor: {data.shape} ({data.nbytes / 1024 / 1024:.2f} MB)")
    ...
    print(f"Response status: {response.status}")
    print(f"Result mean: {result.mean():.4f}")
    print(f"Roundtrip + Serialization time: {t_total:.2f} ms")

if __name__ == "__main__":
    run()
```

Why this is a false positive: the module is a runnable gRPC client example whose entry point is invoked from the `__main__` guard; the rule's fix explicitly allows print for CLIs, and this is a script module, not "library code".

Checklist evidence: the file is an example script whose prints are user-facing output; the detector only skips functions literally named `main`, missing the `run()` CLI entry.

### [ ] Finding `28` — `BP-PY-46`

- Function context: `scripts/tenso/findings/functions/28.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/examples/grpc/server_grpc.py:13:9`
- Checklist pattern: rule condition not met — example script

Source excerpt:

```
class TensorInferenceServicer(tenso_msg_pb2_grpc.TensorInferenceServicer):
    def Predict(self, request, context):
        input_tensor = tenso.loads(request.tensor_packet)
        print(f"Received {request.model_name} request. Shape: {input_tensor.shape}")
```

Why this is a false positive: this is the example gRPC server demo script; the print is example-server logging, and the module is a script, not library code.

Checklist evidence: file lives under `examples/` and is executed as a demo server; no importable library surface.

### [ ] Finding `29` — `PERF-PY-28`

- Function context: `scripts/tenso/findings/functions/29.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/examples/grpc/server_grpc.py:37:44`
- Checklist pattern: rule condition not met — executor is process-lifetime

Source excerpt:

```
def serve():
    ...
    # Pass the options to the server
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10), options=options)
    tenso_msg_pb2_grpc.add_TensorInferenceServicer_to_server(TensorInferenceServicer(), server)
    server.add_insecure_port("[::]:50051")
    print("Tenso gRPC Server starting on port 50051...")
    server.start()
    server.wait_for_termination()
```

Why this is a false positive: the `ThreadPoolExecutor` is created once in `serve()` and owned by `grpc.server` for the whole process lifetime; it is not constructed per unit of work.

Checklist evidence: construction site is the server bootstrap executed once; the detector only checks that the executor line is inside a function (indent > 0).

### [ ] Finding `30` — `BP-PY-46`

- Function context: `scripts/tenso/findings/functions/30.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/examples/grpc/server_grpc.py:43:5`
- Checklist pattern: rule condition not met — example script

Source excerpt:

```
    server.add_insecure_port("[::]:50051")
    print("Tenso gRPC Server starting on port 50051...")
    server.start()
    server.wait_for_termination()
```

Why this is a false positive: startup banner of the example server script; not operational logging in library code.

Checklist evidence: same script module as finding 28; prints are the demo program's output.

### [ ] Findings `31`, `32`, `33`, `34`, `35`, `36`, `37`, `38`, `39`, `40`, `41`, `42`, `43`, `44`, `45`, `46` — `BP-PY-46` (examples/ray_example.py)

- Function context: `scripts/tenso/findings/functions/31.txt` … `46.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/examples/ray_example.py:22:1` (and 26, 34, 35, 36, 37, 41, 65, 66, 67, 71, 99, 103, 135, 136, 139)
- Checklist pattern: rule condition not met — script module

Source excerpt:

```
ray.init()
register()
print("Tenso registered as Ray serializer\n")

# --- Example 1: Basic put/get ---
print("--- Example 1: Object Store ---")
...
print(f"  Shape: {tensor.shape}, Dtype: {tensor.dtype}")
print(f"  Size: {tensor.nbytes / 1e6:.1f} MB")
print(f"  Put + Get: {elapsed:.1f} ms")
print(f"  Match: {np.array_equal(tensor, result)}\n")
```

Why this is a false positive: `ray_example.py` is a top-level tutorial script (no `__main__` guard, no library exports); every print is example-program output, so the "non-script modules" rule condition is not met.

Checklist evidence: module-level executable script; prints are the demo's report output, not debug logging in library code.

### [ ] Findings `47`, `48`, `49`, `50`, `53` — `BP-PY-46` (examples/server.py)

- Function context: `scripts/tenso/findings/functions/47.txt` … `50.txt`, `53.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/examples/server.py:9:1` (and 12, 23, 29, 33)
- Checklist pattern: rule condition not met — script module

Source excerpt:

```
print("Tenso Inference Server Waiting...")
conn, addr = server.accept()
print(f"Connected by {addr}")
...
print(f"Received Input: {tensor.shape} | Mean: {tensor.mean():.4f}")
...
print("Sending response...")
...
except Exception as e:
    print(f"Error: {e}")
```

Why this is a false positive: `examples/server.py` is a standalone socket demo script; the prints are the server program's console output, not library-code logging.

Checklist evidence: top-level script with module-level socket setup; no classes or importable API.

### [ ] Finding `54` — `BP-PY-46`

- Function context: `scripts/tenso/findings/functions/54.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/pyfuzz/_make_seeds.py:28:5`
- Checklist pattern: rule condition not met — dev-tool script

Source excerpt:

```
def _seed(name: str, payload: bytes) -> None:
    SEED_DIR.mkdir(parents=True, exist_ok=True)
    out = SEED_DIR / name
    out.write_bytes(payload)
    print(f"  wrote {out.relative_to(SEED_DIR.parent)} ({len(payload)} bytes)")
```

Why this is a false positive: `_make_seeds.py` is a standalone CLI tool (`#!/usr/bin/env python3`, "Run from anywhere: python pyfuzz/_make_seeds.py"); the print is the tool's progress output, not library-code logging.

Checklist evidence: script module invoked directly; the rule targets non-script modules.

### [ ] Finding `59` — `BP-PY-46`

- Function context: `scripts/tenso/findings/functions/59.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/src/tenso/cache.py:242:13`
- Checklist pattern: detector over-match — print inside docstring

Source excerpt:

```
    Example::

        import numpy as np
        from tenso import TensoCache

        with TensoCache("64MB") as cache:
            cache.put("weights", np.random.randn(1000, 1000).astype(np.float32))
            arr = cache.get("weights")        # zero-copy view into SHM
            print(cache.info("weights"))       # metadata without deserialization
            print(cache.stats)                 # hit/miss counts, memory usage
    """
```

Why this is a false positive: the flagged `print(` lines are inside the class docstring's example block (the `"""` closes at line 244); the BP-PY-46 detector is line-based and strips only `#` comments, so it matches `print(` inside string literals.

Checklist evidence: the shown lines are docstring content, not executable code.

### [ ] Finding `60` — `BP-PY-46`

- Function context: `scripts/tenso/findings/functions/60.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/src/tenso/cache.py:243:13`
- Checklist pattern: detector over-match — print inside docstring

Source excerpt:

```
            arr = cache.get("weights")        # zero-copy view into SHM
            print(cache.info("weights"))       # metadata without deserialization
            print(cache.stats)                 # hit/miss counts, memory usage
    """
```

Why this is a false positive: same class-docstring example block as finding 59; the line is not executable code.

Checklist evidence: `print(cache.stats)` sits between the docstring opening `Example::` block and its closing `"""`.

### [ ] Finding `61` — `BP-PY-7`

- Function context: `scripts/tenso/findings/functions/61.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/src/tenso/cache.py:393:27`
- Checklist pattern: detector over-match — `os.open` is not builtin `open`

Source excerpt:

```
        filename = f".tenso_cache_{uid}_{self._name}.lock"
        self._lock_file_path = os.path.join(tempfile.gettempdir(), filename)
        self._lock_fd = os.open(self._lock_file_path, os.O_CREAT | os.O_RDWR, 0o600)
```

Why this is a false positive: the rule matches calls to `open`/`Path.open` assigned to a name; here the call is `os.open`, the low-level syscall wrapper returning a raw fd that is deliberately held for the process-lifetime file lock and closed explicitly in `close()` (cache.py:1362–1367). The detector's identifier-boundary check (`!isIdentByte('.')` before `open(`) makes it match `os.open` too.

Checklist evidence: no `with`-able builtin `open`/`Path.open` call exists on this line; the fd lifecycle is managed by `close()`.

### [ ] Finding `80` — `CWE-1341`

- Function context: `scripts/tenso/findings/functions/80.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/src/tenso/cache.py:1385:9`
- Checklist pattern: detector over-match — idempotent guarded close, not double release

Source excerpt:

```
    def __exit__(self, exc_type, exc_val, exc_tb):
        self.close()
        if self._owns:
            self.unlink()

    def __del__(self):
        try:
            self.close()
        except Exception:
            pass
```

Why this is a false positive: the regex `\w+\.close\s*\([\s\S]{0,180}\w+\.close\s*\(` pairs `__exit__`'s `self.close()` with `__del__`'s `self.close()`, but `close()` is guarded by `if not self._closed` (cache.py:1352) and sets `_closed = True`, so the shared-memory handle is released exactly once.

Checklist evidence: both calls are on the same object, but the guard makes the second call a no-op — no second release of the handle occurs.

### [ ] Finding `83` — `BP-PY-46`

- Function context: `scripts/tenso/findings/functions/83.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/src/tenso/client.py:26:9`
- Checklist pattern: detector over-match — print inside docstring

Source excerpt:

```
    Example::

        client = TensoFastAPIClient("http://localhost:8000")
        result = client.predict("/infer", np.random.randn(1, 224, 224, 3).astype(np.float32))
        print(result.shape)

        # Async usage
        result = await client.apredict("/infer", tensor)
    """
```

Why this is a false positive: the `print(result.shape)` line is inside the `TensoFastAPIClient` class docstring's example block (docstring closes at line 30); it is not executable code.

Checklist evidence: line 26 sits between `Example::` and the closing `"""`; line-based detector does not mask docstrings.

### [ ] Finding `84` — `CWE-1341`

- Function context: `scripts/tenso/findings/functions/84.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/src/tenso/client.py:168:9`
- Checklist pattern: detector over-match — distinct lifecycle methods, idempotent close

Source excerpt:

```
    async def aclose(self):
        """Close both sync and async clients."""
        if self._async_client is not None:
            await self._async_client.aclose()
            self._async_client = None
        self.close()
```

Why this is a false positive: the regex pairs `self.close()` (line 168) with the `self.close()` inside `__exit__` (line 174); `close()` only releases the sync client when `self._sync_client is not None` and then sets it to `None` (client.py:159–161), so the httpx handle is closed at most once.

Checklist evidence: two `.close()` calls within 180 characters, but they are guarded, idempotent calls in different lifecycle methods — no duplicate release of the same handle.

### [ ] Finding `92` — `CWE-695`

- Function context: `scripts/tenso/findings/functions/92.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/src/tenso/core.py:692:14`
- Checklist pattern: rule condition not met — mmap is the product's intended API, not prohibited low-level functionality

Source excerpt:

```
    if mmap_mode:
        mm = mmap.mmap(fp.fileno(), 0, access=mmap.ACCESS_READ)
        return loads(mm, copy=copy)
```

Why this is a false positive: the rule condition is "low-level functionality explicitly prohibited by the framework or specification under which the product is supposed to operate"; tenso is a zero-copy serialization library whose public `load(mmap_mode=True)` API is precisely this stdlib `mmap` use — it is the documented platform API, not a bypass of higher-level safety controls.

Checklist evidence: the call is the opt-in `mmap_mode` branch of the library's own public load API; the detector mechanically matches any `mmap.mmap(` call.

### [ ] Finding `93` — `CWE-93`

- Function context: `scripts/tenso/findings/functions/93.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/src/tenso/fastapi.py:54:17`
- Checklist pattern: rule condition not met — value cannot contain CRLF

Source excerpt:

```
        self.headers["X-Tenso-Version"] = "4"
        if hasattr(tensor, "shape"):
            self.headers["X-Tenso-Shape"] = str(tensor.shape)
        if hasattr(tensor, "dtype"):
```

Why this is a false positive: `tensor.shape` is always a tuple of integers for every type the protocol can deserialize (numpy arrays, scipy sparse matrices), so `str(tensor.shape)` can never contain CR or LF characters; the detector already excludes numeric-derived values (`str(int(...)`, `str(round(...)`), and `str(tensor.shape)` is the same class of non-injectable value, so the CRLF-neutralization condition is vacuous.

Checklist evidence: the header value is the string form of an integer tuple — CR/LF cannot appear regardless of input.

### [ ] Finding `98` — `CWE-1341`

- Function context: `scripts/tenso/findings/functions/98.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/src/tenso/gpu.py:315:21`
- Checklist pattern: detector over-match — two distinct handles

Source excerpt:

```
            fobj = os.fdopen(fd, 'rb', closefd=False)
            try:
                f = self._kvikio.CuFile(fobj)
                try:
                    bytes_read = f.pread(gpu_buf, nbytes, file_offset=offset)
                finally:
                    f.close()
            finally:
                fobj.close()
```

Why this is a false positive: `f.close()` closes the kvikio `CuFile` wrapper and `fobj.close()` closes the `os.fdopen` file object; they are two different handles (with `closefd=False` the underlying fd is owned by neither), so no resource is released twice.

Checklist evidence: the two `.close()` calls within 180 characters target distinct objects.

### [ ] Finding `104` — `BP-PY-46`

- Function context: `scripts/tenso/findings/functions/104.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/src/tenso/shm.py:39:13`
- Checklist pattern: detector over-match — print inside docstring

Source excerpt:

```
    Example::

        # Writer
        ary = np.random.rand(100, 100)
        with TensoShm.create_from("my_tensor", ary) as shm:
            print("Wrote to SHM")
            input("Press enter to cleanup...")
```

Why this is a false positive: the `print("Wrote to SHM")` line is inside the `TensoShm` class docstring example (docstring closes at line 46); it is not executable code.

Checklist evidence: line 39 sits between `Example::` and the closing `"""`.

### [ ] Finding `105` — `BP-PY-46`

- Function context: `scripts/tenso/findings/functions/105.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/src/tenso/shm.py:45:13`
- Checklist pattern: detector over-match — print inside docstring

Source excerpt:

```
        # Reader
        with TensoShm("my_tensor") as shm:
            ary = shm.get()
            print(ary.shape)
    """
```

Why this is a false positive: same class-docstring example block as finding 104; the `print(ary.shape)` line is docstring content, not executable code.

Checklist evidence: line 45 precedes the docstring closing `"""` at line 46.

### [ ] Finding `114` — `CWE-1341`

- Function context: `scripts/tenso/findings/functions/114.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/tests/test_cache.py:61:5`
- Checklist pattern: detector over-match — distinct cache instances

Source excerpt:

```
@pytest.fixture
def file_lock_cache(monkeypatch):
    ...
    c = TensoCache("4MB")
    assert not c._use_robust_mutex, "fixture failed to disable robust mutex"
    yield c
    c.close()
    c.unlink()
```

Why this is a false positive: the regex pairs `c.close()` in the `file_lock_cache` fixture (line 61) with `c.close()` in the `small_cache` fixture (line 70); the two calls operate on different `TensoCache` instances, so the same handle is never released twice.

Checklist evidence: the matched `.close()` calls are in separate fixtures on separate `c` objects; `unlink()` (line 62) is not a close call.

### [ ] Finding `118` — `CWE-93`

- Function context: `scripts/tenso/findings/functions/118.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/tests/test_fastapi.py:29:20`
- Checklist pattern: detector over-match — header read in assertion, not a write

Source excerpt:

```
    assert response.status_code == 200
    assert response.headers["content-type"] == "application/octet-stream"
    # Header advertises the wire protocol version (v4), matching what's emitted.
```

Why this is a false positive: the line is an assertion that *reads* the `content-type` header and compares it to a constant; nothing is written to a response header, so the "dynamic value written to an HTTP response header" condition is not met. The detector's line-based `.headers[` branch does not verify that the line is an assignment.

Checklist evidence: `==` comparison of an existing header against a string literal — no header mutation sink.

### [ ] Finding `125` — `CWE-1341`

- Function context: `scripts/tenso/findings/functions/125.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/tests/test_shm.py:56:9`
- Checklist pattern: detector over-match — distinct handles

Source excerpt:

```
    shm = multiprocessing.shared_memory.SharedMemory(name=name, create=True, size=1024)
    ts = TensoShm(name)
    ...
    finally:
        ts.close()
        shm.close()
        shm.unlink()
```

Why this is a false positive: `ts.close()` closes the `TensoShm` wrapper and `shm.close()` closes the underlying `multiprocessing.shared_memory.SharedMemory`; they are two distinct handles, each released once.

Checklist evidence: two `.close()` calls within 180 characters target different objects.

## True positives

### CWE-1121 — Excessive McCabe Cyclomatic Complexity

| Finding | Source | Reason |
| --- | --- | --- |
| 3 | benchmark.py:262 | `run_serialization()` counts 17 `if/elif/for/while/except` branches on masked code (verified with the detector's counting logic); ≥ 12 threshold met. |
| 70 | cache.py:1055 | `TensoCache.get()` counts exactly 12 branches on masked code; ≥ 12 threshold met. |

### CWE-1084 — Excessive File or Data Access Operations

| Finding | Source | Reason |
| --- | --- | --- |
| 6 | benchmark.py:349 | `run_io()` contains ≥ 3 `open(` calls (four `with open(path, ...)` statements), meeting the rule's ≥ 3 threshold. |

### CWE-1124 — Excessively Deep Nesting

| Finding | Source | Reason |
| --- | --- | --- |
| 109 | utils.py:86 | assignment is nested under `if HAS_RUST` → `try` → `if "dtype" not in info` → `if dc in _QDTYPE_NAMES` → `if "total_elements" in info` → `if dc in (18, 19)`: six control-flow levels, meeting the ≥ 6 threshold. |

### CWE-367 — TOCTOU

| Finding | Source | Reason |
| --- | --- | --- |
| 7 | benchmark.py:379 | `if os.path.exists(path): os.remove(path)` — check-then-use of the same path within the regex's 300-char window; condition satisfied. |
| 117 | test_core.py:157 | `if os.path.exists(tmp_path): os.remove(tmp_path)` — same check-then-use pattern; condition satisfied. |

### BP-PY-1 — Bare Except Clause

| Finding | Source | Reason |
| --- | --- | --- |
| 4 | benchmark.py:344 | `except Exception as e:` logs only, no re-raise. |
| 10 | benchmark.py:487 | `except Exception: pass` in sink server. |
| 14 | benchmark.py:517 | `except Exception as e:` logs only. |
| 16 | benchmark.py:752 | `except Exception as e:` logs only. |
| 51 | examples/server.py:32 | `except Exception as e:` logs only. |
| 57 | src/tenso/__init__.py:37 | `except Exception:` with fallback assignment, no re-raise. |
| 62 | cache.py:436 | `except Exception as exc:` prints, no re-raise. |
| 68 | cache.py:1050 | `except Exception as exc:` prints, no re-raise. |
| 72 | cache.py:1218 | `except Exception:` returns default dict. |
| 73 | cache.py:1360 | `except Exception: pass` in `close()`. |
| 77 | cache.py:1373 | `except Exception: pass` in `unlink()`. |
| 81 | cache.py:1392 | `except Exception: pass` in `__del__`. |
| 88 | core.py:609 | `except Exception:` sets `_hdr = None`, no re-raise. |
| 99 | ray.py:207 | `except Exception: pass` in `unregister()`. |
| 119 | test_gpu.py:39 | `except Exception:` sets `HAS_CUDA_CUPY = False`. |
| 121 | test_gpu.py:168 | `except Exception: pass` in test. |

### BP-PY-2 — Except Pass

| Finding | Source | Reason |
| --- | --- | --- |
| 11 | benchmark.py:487 | handler body is only `pass`. |
| 15 | benchmark.py:718 | `except ImportError: pass` (optional bfloat16). |
| 66 | cache.py:872 | `except (TypeError, NotImplementedError): pass` before fallback. |
| 71 | cache.py:1120 | `except (ValueError, TypeError): pass` before fallback. |
| 74 | cache.py:1360 | `except Exception: pass` in `close()`. |
| 76 | cache.py:1365 | `except OSError: pass` around `os.close`. |
| 78 | cache.py:1373 | `except Exception: pass` in `unlink()`. |
| 79 | cache.py:1378 | `except OSError: pass` around `os.unlink`. |
| 82 | cache.py:1392 | `except Exception: pass` in `__del__`. |
| 85 | config.py:84 | `except ImportError: pass` (optional bfloat16). |
| 90 | core.py:664 | `except (ValueError, TypeError, AttributeError, OSError): pass` before fallback. |
| 95 | gpu.py:259 | `except ImportError: pass` (optional kvikio). |
| 100 | ray.py:207 | `except Exception: pass` in `unregister()`. |
| 106 | shm.py:116 | `except (NotImplementedError, TypeError, AttributeError): pass` before fallback. |
| 110 | tests/fixtures/_generate.py:65 | `except ImportError: pass` (optional bfloat16). |
| 113 | tests/fixtures/_generate.py:102 | `except ImportError: pass` (optional scipy). |
| 116 | test_cache.py:735 | `except MemoryError: pass` in test loop. |
| 122 | test_gpu.py:168 | `except Exception: pass` in test. |

### CWE-390 — Detection of Error Condition Without Action

| Finding | Source | Reason |
| --- | --- | --- |
| 12 | benchmark.py:487 | exception detected, handler only `pass`. |
| 67 | cache.py:872 | exception detected, handler only `pass`. |
| 86 | config.py:84 | exception detected, handler only `pass`. |
| 91 | core.py:664 | exception detected, handler only `pass`. |
| 96 | gpu.py:259 | exception detected, handler only `pass`. |
| 101 | ray.py:207 | exception detected, handler only `pass`. |
| 107 | shm.py:116 | exception detected, handler only `pass`. |
| 111 | tests/fixtures/_generate.py:65 | exception detected, handler only `pass`. |
| 123 | test_gpu.py:168 | exception detected, handler only `pass`. |

### CWE-396 — Declaration of Catch for Generic Exception

| Finding | Source | Reason |
| --- | --- | --- |
| 5 | benchmark.py:344 | `except Exception as e:`. |
| 52 | examples/server.py:32 | `except Exception as e:`. |
| 55 | pyfuzz/fuzz_get_packet_info.py:34 | `except Exception as exc:` (re-raised as AssertionError). |
| 56 | pyfuzz/fuzz_loads.py:42 | `except Exception as exc:` (re-raised as AssertionError). |
| 58 | src/tenso/__init__.py:37 | `except Exception:`. |
| 63 | cache.py:436 | `except Exception as exc:`. |
| 89 | core.py:609 | `except Exception:`. |
| 94 | fastapi.py:89 | `except Exception as e:` → HTTPException. |
| 102 | ray.py:207 | `except Exception:`. |
| 108 | shm.py:180 | `except Exception:` → `shm.close()`. |

### CWE-1071 — Empty Code Block

| Finding | Source | Reason |
| --- | --- | --- |
| 13 | benchmark.py:487 | handler contains only `pass`. |
| 75 | cache.py:1360 | handler contains only `pass`. |
| 87 | config.py:84 | handler contains only `pass`. |
| 97 | gpu.py:259 | handler contains only `pass`. |
| 103 | ray.py:207 | handler contains only `pass`. |
| 112 | tests/fixtures/_generate.py:65 | handler contains only `pass`. |
| 124 | test_gpu.py:168 | handler contains only `pass`. |

### BP-PY-46 — print Debugging In Library Code

| Finding | Source | Reason |
| --- | --- | --- |
| 64 | cache.py:437 | executable print in library error path (stderr diagnostic). |
| 65 | cache.py:451 | executable print in library lock-recovery path. |
| 69 | cache.py:1051 | executable print in library metadata fallback. |

### BP-PY-41 — pytest assert With Side Effects Only

| Finding | Source | Reason |
| --- | --- | --- |
| 115 | test_cache.py:682 | `test_acquire_release` calls `_shm_lock_acquire`/`_shm_lock_release` with no assert — heuristic placeholder-test signal matches. |

### BP-PY-42 — unittest Assert Without Context On Raises

| Finding | Source | Reason |
| --- | --- | --- |
| 120 | test_gpu.py:166 | test uses bare `try/except` to expect a failure instead of `pytest.raises`. |

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none (single-reviewer audit)
- Chunk evidence: `scripts/tenso/chunks`
- Function evidence: `scripts/tenso/findings/functions`
- Validation: `git diff --check` — `pass`

## Post-fix remaining-FP audit (2026-08-02)

### Run metadata

```yaml
timestamp: 2026-08-02T16:29:00Z
repository: tenso
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso
branch: main
commit: ee5d6eb7baba8aca90b1d63a5a176b0a7d37692e
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso
chunk_path: scripts/tenso/chunks
function_context_path: scripts/tenso/findings/functions
fix_commit: b5b8fde (FP-reduction fix; binary rebuilt 2026-08-02 16:29)
```

### Scan evidence

- Build command: `make build` (local `./bin/goslop` binary rebuilt post-fix)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/tenso/chunks -context-dir scripts/tenso/findings/functions real-repos/tenso`
- Findings: `75`
- Chunks reviewed: `scripts/tenso/chunks/Chunk_1_25.txt`, `Chunk_26_50.txt`, `Chunk_51_75.txt`
- Function contexts reviewed: `scripts/tenso/findings/functions/1.txt`, `2.txt`, `8.txt`, `16.txt`, `17.txt`, `40.txt`, `43.txt`, `51.txt`, `52.txt`, plus the enclosing sources for those findings

### Audit checklist

- [x] Read every assigned chunk under `scripts/tenso/chunks`.
- [x] Read `scripts/tenso/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

### Classification summary (fresh run)

Fresh finding IDs do not correspond to old IDs; matched by `Source:` path (file:line:col) against the audited TP/FP lists. All 75 fresh findings matched an audited source; no new findings appeared. 66 match audited TPs, 9 re-appearing audited FPs, 0 Uncertain.

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 9 | 1, 2, 8, 16, 17, 40, 43, 51, 52 |
| True positive | 66 | 3, 4, 5, 6, 7, 9, 10, 11, 12, 13, 14, 15, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 41, 42, 44, 45, 46, 47, 48, 49, 50, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75 |
| Uncertain | 0 | — |

Remaining false positives after the fix: 9 (down from 54). The fix suppressed the audited FPs not present in this run: BP-PY-45 (docs/source/conf.py), BP-PY-46 example-script prints (examples/client.py, client_grpc.py, server_grpc.py, ray_example.py, server.py, pyfuzz/_make_seeds.py), BP-PY-46 docstring prints (cache.py, client.py, shm.py), BP-PY-7 (cache.py `os.open`), CWE-1341 (gpu.py:315, test_cache.py:61, test_shm.py:56) and CWE-93 (test_fastapi.py:29).

### False positives

### [ ] Findings `1`, `2` — `BP-PY-10` / `CWE-502`

- Function context: `scripts/tenso/findings/functions/1.txt`, `2.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/benchmark.py:111:21`
- Checklist pattern: rule condition not met — no untrusted-data sink (re-appearing audited FPs 1 and 2)

Source excerpt:

```
    enc = lambda x: pickle.dumps(x, protocol=pickle.HIGHEST_PROTOCOL)  # noqa
    dec = lambda x: pickle.loads(x)  # noqa
    return enc, dec
```

Why this is a false positive: `bench_pickle` is a serialization-format benchmark; `dec` round-trips only bytes produced by `enc` in the same process from `np.random`-generated data, so no rule source category (request body, user-path file, cache) applies.

Checklist evidence: no external-input call path into `pickle.loads` exists in the shown source; same construct as audited FPs 1/2.

### [ ] Finding `8` — `BP-PY-10`

- Function context: `scripts/tenso/findings/functions/8.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/benchmark.py:409:13`
- Checklist pattern: rule condition not met — self-written temp file, not user input (re-appearing audited FP 8)

Source excerpt:

```
    with tempfile.NamedTemporaryFile(delete=False) as f:
        path = f.name
    try:
        t0 = time.perf_counter()
        with open(path, "wb") as f:
            pickle.dump(data, f)
        ...
        with open(path, "rb") as f:
            t0 = time.perf_counter()
            pickle.load(f)
```

Why this is a false positive: `path` is the benchmark's own `NamedTemporaryFile` written by this same function immediately before the load; it is not a user-path file, so the rule's "file from user path" source category does not apply.

Checklist evidence: the handle passed to `pickle.load` was opened from a path created and written by the same function.

### [ ] Finding `16` — `PERF-PY-28`

- Function context: `scripts/tenso/findings/functions/16.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/examples/grpc/bench_server.py:29:44`
- Checklist pattern: rule condition not met — executor is process-lifetime (re-appearing audited FP 23)

Source excerpt:

```
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=1), options=options)
    benchmark_msg_pb2_grpc.add_BenchmarkerServicer_to_server(Benchmarker(), server)
    server.add_insecure_port("[::]:50051")
    print("Tenso Benchmark Server starting on port 50051...")
    server.start()
    server.wait_for_termination()
```

Why this is a false positive: the executor is constructed once in `serve()` and handed to `grpc.server`, which owns it for the server's lifetime — it is not created per unit of work as the rule condition requires.

Checklist evidence: construction site is the server bootstrap executed once per process.

### [ ] Finding `17` — `PERF-PY-28`

- Function context: `scripts/tenso/findings/functions/17.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/examples/grpc/server_grpc.py:37:44`
- Checklist pattern: rule condition not met — executor is process-lifetime (re-appearing audited FP 29)

Source excerpt:

```
    # Pass the options to the server
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10), options=options)
    ...
    server.add_insecure_port("[::]:50051")
    server.start()
    server.wait_for_termination()
```

Why this is a false positive: same pattern as finding 16 — the executor is created once in `serve()` and owned by `grpc.server` for the whole process lifetime, not per unit of work.

Checklist evidence: construction site is the server bootstrap executed once; the detector only checks that the executor line is inside a function (indent > 0).

### [ ] Finding `40` — `CWE-1341`

- Function context: `scripts/tenso/findings/functions/40.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/src/tenso/cache.py:1385:9`
- Checklist pattern: detector over-match — idempotent guarded close, not double release (re-appearing audited FP 80)

Source excerpt:

```
    def close(self):
        """Close access to the shared memory pool."""
        if not self._closed:
            self._closed = True
            try:
                self._shm.close()
            ...

    def __exit__(self, exc_type, exc_val, exc_tb):
        self.close()
        if self._owns:
            self.unlink()

    def __del__(self):
        try:
            self.close()
        except Exception:
            pass
```

Why this is a false positive: the regex pairs `__exit__`'s `self.close()` with `__del__`'s `self.close()`, but `close()` is guarded by `if not self._closed` and sets `_closed = True`, so the shared-memory handle is released exactly once.

Checklist evidence: both calls target the same object, but the guard makes the second call a no-op — no second release of the handle occurs.

### [ ] Finding `43` — `CWE-1341`

- Function context: `scripts/tenso/findings/functions/43.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/src/tenso/client.py:168:9`
- Checklist pattern: detector over-match — distinct lifecycle methods, idempotent close (re-appearing audited FP 84)

Source excerpt:

```
    def close(self):
        """Close sync client. For async client, use ``async with`` or ``await client.aclose()``."""
        if self._sync_client is not None:
            self._sync_client.close()
            self._sync_client = None

    async def aclose(self):
        """Close both sync and async clients."""
        if self._async_client is not None:
            await self._async_client.aclose()
            self._async_client = None
        self.close()

    def __exit__(self, *args):
        self.close()
```

Why this is a false positive: the regex pairs `self.close()` in `aclose()` (line 168) with `self.close()` in `__exit__`; `close()` only releases the sync client when `self._sync_client is not None` and then sets it to `None`, so the httpx handle is closed at most once.

Checklist evidence: two `.close()` calls within 180 characters, but they are guarded, idempotent calls in different lifecycle methods — no duplicate release of the same handle.

### [ ] Finding `51` — `CWE-695`

- Function context: `scripts/tenso/findings/functions/51.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/src/tenso/core.py:692:14`
- Checklist pattern: rule condition not met — mmap is the product's intended API, not prohibited low-level functionality (re-appearing audited FP 92)

Source excerpt:

```
    if mmap_mode:
        mm = mmap.mmap(fp.fileno(), 0, access=mmap.ACCESS_READ)
        return loads(mm, copy=copy)
```

Why this is a false positive: the rule condition is "low-level functionality explicitly prohibited by the framework or specification under which the product is supposed to operate"; tenso is a zero-copy serialization library whose public `load(mmap_mode=True)` API is precisely this stdlib `mmap` use — the documented platform API, not a bypass of higher-level safety controls.

Checklist evidence: the call is the opt-in `mmap_mode` branch of the library's own public load API; the detector mechanically matches any `mmap.mmap(` call.

### [ ] Finding `52` — `CWE-93`

- Function context: `scripts/tenso/findings/functions/52.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/tenso/src/tenso/fastapi.py:54:17`
- Checklist pattern: rule condition not met — value cannot contain CRLF (re-appearing audited FP 93)

Source excerpt:

```
        self.headers["X-Tenso-Version"] = "4"
        if hasattr(tensor, "shape"):
            self.headers["X-Tenso-Shape"] = str(tensor.shape)
        if hasattr(tensor, "dtype"):
            self.headers["X-Tenso-Dtype"] = str(tensor.dtype)
```

Why this is a false positive: `tensor.shape` is always a tuple of integers for every type the protocol can deserialize, so `str(tensor.shape)` can never contain CR or LF; like the detector's already-excluded `str(int(...))`/`str(round(...))` forms, the CRLF-neutralization condition is vacuous.

Checklist evidence: the header value is the string form of an integer tuple — CR/LF cannot appear regardless of input.

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none (single-reviewer audit)
- Chunk evidence: `scripts/tenso/chunks`
- Function evidence: `scripts/tenso/findings/functions`
- Validation: `git diff --check` — `pass`
