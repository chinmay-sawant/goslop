# False-positive audit report

## Run metadata

```yaml
timestamp: 2026-07-31T18:23:54Z
repository: goslop
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop
branch: feat/python-perf-ruleset-plan
commit: 25e0d61fb89e3173f72c2919823e2849b2f2c149
scan_target: /home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets
chunk_path: ./scripts/chunks
function_context_path: ./scripts/findings/functions
```

## Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config "templates/goslop-python.toml" --export-context --export-chunks --no-cache "/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets"`
- Findings: 56
- Chunks reviewed: `./scripts/chunks/Chunk_1_25.txt`, `./scripts/chunks/Chunk_26_50.txt`, `./scripts/chunks/Chunk_51_56.txt`
- Function contexts reviewed for every proposed false positive: `./scripts/findings/functions/1.txt`, `2.txt`, `3.txt`, `41.txt`–`46.txt`, and `50.txt`–`56.txt`

## Audit checklist

- [x] Read every assigned chunk under `./scripts/chunks`.
- [x] Read `./scripts/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed each proposed false positive's `Source:` path and inspected the enclosing source when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based each decision on the rule condition and the shown source, not on application-specific assumptions.
- [x] Reconciled the chunk review with a full 56-finding cross-check; no delegated-review disagreement remained.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 16 | 1–3, 41–46, 50–56 |
| True positive | 40 | 4–40, 47–49 |
| Uncertain | 0 | None |

## False positives

### [x] Finding 1 — CWE-89

- Function context: `./scripts/findings/functions/1.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/django-flash-sale-inventory/bench/reserve_bench.py:24:16`
- Checklist pattern: SQL construction versus passthrough

Source excerpt:

```
def __call__(self, execute, sql, params, many, context):
    self.count += 1
    return execute(sql, params, many, context)
```

Why this is a false positive: The benchmark wrapper forwards a callback-provided SQL/params pair; it does not construct or interpolate SQL text.

Checklist evidence: No formatting, concatenation, or dynamic SQL construction appears in the reported statement.

### [x] Finding 2 — BP-PY-46

- Function context: `./scripts/findings/functions/2.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/django-flash-sale-inventory/bench/reserve_bench.py:37:5`
- Checklist pattern: Intentional benchmark output

Source excerpt:

```
print(f"{label:<24}: {statistics.median(times) * 1000:9.1f} ms/op  queries={last_queries}  (median of {N_REPEATS})")
```

Why this is a false positive: This print reports benchmark measurements from the `bench/` harness; it is intentional user-facing benchmark output, not library debugging.

Checklist evidence: The message is explicitly a timing/query summary, and the enclosing file is invoked through its benchmark `main()` entry point.

### [x] Finding 3 — BP-PY-46

- Function context: `./scripts/findings/functions/3.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/django-flash-sale-inventory/bench/reserve_bench.py:45:5`
- Checklist pattern: Intentional benchmark output

Source excerpt:

```
def main():
    ...
    print("== Django service-layer benchmarks (2026-07-31, sqlite file DB) ==")

if __name__ == "__main__":
    main()
```

Why this is a false positive: The print is the banner for a benchmark CLI whose entry point is protected by the `__main__` guard.

Checklist evidence: This is deliberate command output from a benchmark driver, not operational logging in reusable application code.

### [x] Finding 41 — BP-PY-46

- Function context: `./scripts/findings/functions/41.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/fastapi-live-metrics-ingest/bench/microbench.py:25:1`
- Checklist pattern: Intentional benchmark output

Source excerpt:

```
print("== FastAPI CPU hot-path microbenchmarks (2026-07-31) ==")
```

Why this is a false positive: The source explicitly labels this as a microbenchmark banner; the output is the harness result stream, not debugging in library code.

Checklist evidence: The reported file is under `bench/` and the message names the benchmark being run.

### [x] Finding 42 — BP-PY-46

- Function context: `./scripts/findings/functions/42.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/fastapi-live-metrics-ingest/bench/microbench.py:26:1`
- Checklist pattern: Intentional benchmark output

Source excerpt:

```
print(f"normalize_route 100k labels : {bench('route', lambda: [normalize_route(l) for l in LABELS], 100_000)}")
```

Why this is a false positive: This is a named benchmark measurement emitted by the microbenchmark harness, not a debug print from application runtime code.

Checklist evidence: The output identifies the measured operation and sample size.

### [x] Finding 43 — BP-PY-46

- Function context: `./scripts/findings/functions/43.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/fastapi-live-metrics-ingest/bench/microbench.py:47:1`
- Checklist pattern: Intentional benchmark output

Source excerpt:

```
print(f"pydantic validate 100-sample : {bench('pyd', validate, 1)}")
```

Why this is a false positive: The print emits the measured Pydantic benchmark result and is part of the intended harness report.

Checklist evidence: The message names the benchmark operation and its sample count.

### [x] Finding 44 — BP-PY-46

- Function context: `./scripts/findings/functions/44.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/fastapi-live-metrics-ingest/bench/microbench.py:56:5`
- Checklist pattern: Intentional benchmark output

Source excerpt:

```
print(f"app-side percentile sort {n:<9}: {bench('sort', sort_it, 1)}")
```

Why this is a false positive: This is a benchmark result for the explicitly named app-side sort operation, not reusable-code logging or debugging.

Checklist evidence: The source names the measured operation and prints its timing result.

### [x] Finding 45 — CWE-89

- Function context: `./scripts/findings/functions/45.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/fastapi-live-metrics-ingest/bench/seed_db.py:46:26`
- Checklist pattern: ORM/query-builder expression

Source excerpt:

```
from sqlalchemy import insert
...
await session.execute(insert(MetricSample), rows)
```

Why this is a false positive: `insert(MetricSample)` is a SQLAlchemy Core statement and `rows` is a bound mapping sequence, not interpolated SQL command text.

Checklist evidence: The execute call receives a typed query-builder expression with separate row parameters.

### [x] Finding 46 — BP-PY-46

- Function context: `./scripts/findings/functions/46.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/fastapi-live-metrics-ingest/bench/seed_db.py:49:5`
- Checklist pattern: Intentional seed/benchmark output

Source excerpt:

```
async def main():
    ...
    print(f"seeded {SAMPLE_ROWS} metric samples for tenant_id={tenant_id}")

if __name__ == "__main__":
    asyncio.run(main())
```

Why this is a false positive: The print reports completion of the local benchmark database seed and is emitted from a guarded CLI entry point.

Checklist evidence: The message is a seed-completion status line, not library debug logging.

### [x] Finding 50 — BP-PY-46

- Function context: `./scripts/findings/functions/50.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/flask-partner-webhook-relay/bench/delivery_bench.py:47:9`
- Checklist pattern: Intentional benchmark setup output

Source excerpt:

```
def seed_outbox(app, n, round_id):
    ...
    print(f"seeded {n} outbox items -> {endpoint.url}")
```

Why this is a false positive: The helper reports benchmark fixture setup; it is not operational logging from the Flask application library.

Checklist evidence: The message explicitly says that benchmark outbox items were seeded.

### [x] Finding 51 — BP-PY-46

- Function context: `./scripts/findings/functions/51.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/flask-partner-webhook-relay/bench/delivery_bench.py:68:5`
- Checklist pattern: Intentional benchmark output

Source excerpt:

```
def main():
    ...
    print("== Flask delivery worker benchmark (2026-07-31, mock partner 200ms) ==")

if __name__ == "__main__":
    main()
```

Why this is a false positive: This is the banner for a guarded delivery benchmark CLI.

Checklist evidence: The message names the benchmark and its mock partner latency.

### [x] Finding 52 — BP-PY-46

- Function context: `./scripts/findings/functions/52.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/flask-partner-webhook-relay/bench/delivery_bench.py:69:5`
- Checklist pattern: Intentional benchmark output

Source excerpt:

```
print(f"run_once({N_ITEMS} items, sequential) : {statistics.median(times):.2f}s median  "
      f"(per-item ≈ {statistics.median(times) / N_ITEMS * 1000:.0f} ms; partner latency 200ms)")
```

Why this is a false positive: The line is an intentional benchmark measurement of sequential delivery and its per-item latency.

Checklist evidence: The output reports the exact measured operation, item count, and latency.

### [x] Finding 53 — BP-PY-46

- Function context: `./scripts/findings/functions/53.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/flask-partner-webhook-relay/bench/delivery_bench.py:77:13`
- Checklist pattern: Intentional benchmark output

Source excerpt:

```
print(f"single deliver()                 : {(time.perf_counter() - t0) * 1000:.0f} ms")
```

Why this is a false positive: The print reports a single benchmark sample from the delivery harness, not a library debug message.

Checklist evidence: The output is explicitly labeled `single deliver()` and contains a measured duration.

### [x] Findings 54–55 — CWE-312, CWE-798

- Function contexts: `./scripts/findings/functions/54.txt`, `./scripts/findings/functions/55.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/flask-partner-webhook-relay/bench/seed.py:18:1`
- Checklist pattern: Synthetic benchmark credential

Source excerpt:

```
endpoint = PartnerEndpoint(
    partner_id=partner.id,
    url="http://127.0.0.1:8200/webhook",
    secret="bench-secret",
    is_active=True,
)
```

Why these are false positives: The value is an explicitly synthetic secret for a localhost benchmark endpoint, not a deployment credential stored in application code.

Checklist evidence: Both the `bench-secret` name and loopback endpoint establish benchmark-only fixture data.

### [x] Finding 56 — BP-PY-46

- Function context: `./scripts/findings/functions/56.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/flask-partner-webhook-relay/bench/seed.py:23:5`
- Checklist pattern: Intentional seed/benchmark output

Source excerpt:

```
db.session.commit()
print(f"seeded partner={partner.id}, endpoint={endpoint.id} -> {endpoint.url}")
```

Why this is a false positive: The print reports completion of a local benchmark seed script and is not library debugging output.

Checklist evidence: The message is explicitly a `seeded` status line and the source file is the benchmark seed harness.

## Uncertain findings

No uncertain findings remain.

## Final evidence

- Delegated reviewers: none; local chunk-by-chunk review plus full-set cross-check
- Chunk evidence: `./scripts/chunks`
- Function evidence: `./scripts/findings/functions`
- Validation: `git diff --check` — pass
