# False-positive audit report

## Run metadata

```yaml
timestamp: 2026-07-31T14:26:39Z
repository: goslop
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop
branch: main
commit: 2f4dd861613deae37966f688a7fb3605627b0619
scan_target: /home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets
chunk_path: ./scripts/chunks
function_context_path: ./scripts/findings/functions
```

## Scan evidence

- Build command: Not run in this audit; existing exported evidence was reviewed.
- Scan command: Not run in this audit; existing exported evidence was reviewed.
- Findings: 49
- Chunks reviewed: `./scripts/chunks/Chunk_1_25.txt`, `./scripts/chunks/Chunk_26_49.txt`
- Function contexts reviewed: `./scripts/findings/functions/<finding-id>.txt` for every false-positive candidate and both uncertain findings.

## Audit checklist

- [x] Read both assigned chunks under `./scripts/chunks`.
- [x] Read `./scripts/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed each `Source:` path and inspected enclosing source where the exported context was insufficient.
- [x] Classified all 49 findings as false positive, true positive, or uncertain.
- [x] Used general software and security knowledge only; no application-specific assumptions were used.
- [x] Reconciled independent reviews from a first-chunk auditor, second-chunk auditor, and full-set cross-checker.
- [x] Kept unresolved reviewer disagreement as `Uncertain`.
- [x] `git diff --check` passes after writing this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 29 | 12–19, 21–25, 29–32, 36–43, 46–49 |
| True positive | 20 | 1–11, 20, 26–28, 33–35, 44–45 |
| Uncertain | 0 | None |

## False positives

### [x] Finding 12 — CWE-1124

- Function context: `./scripts/findings/functions/12.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/django-flash-sale-inventory/flash_sale/settings.py:13:1`
- Checklist pattern: Indentation/nesting

Source excerpt:

```
INSTALLED_APPS = [
    'django.contrib.admin',
]
```

Why this is a false positive: The reported line is a configuration-list element, not an executable statement nested through control flow.

Checklist evidence: The indentation belongs to collection layout, so the nesting precondition is absent.

### [x] Finding 13 — CWE-1124

- Function context: `./scripts/findings/functions/13.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/django-flash-sale-inventory/inventory/services/availability.py:42:1`
- Checklist pattern: Indentation/nesting

Source excerpt:

```
for ws in qs:
    result.append({
        'sku_code': ws.sku.sku_code,
    })
```

Why this is a false positive: The additional indentation is a dictionary field inside one loop, not another executable control-flow level.

Checklist evidence: Data-literal layout is being counted as nesting.

### [x] Finding 14 — CWE-1124

- Function context: `./scripts/findings/functions/14.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/django-flash-sale-inventory/inventory/services/reservation.py:56:1`
- Checklist pattern: Indentation/nesting

Source excerpt:

```
item_details.append({
    'sku': sku,
    'stock': stock,
    'quantity': item['quantity'],
})
```

Why this is a false positive: The finding points at a dictionary member; the indentation represents data layout, not deep control flow.

Checklist evidence: The reported construct is a multiline literal.

### [x] Findings 15–17 — BP-PY-42, BP-PY-1, CWE-396

- Function contexts: `./scripts/findings/functions/15.txt`, `./scripts/findings/functions/16.txt`, `./scripts/findings/functions/17.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/django-flash-sale-inventory/inventory/tests.py:379:1`
- Checklist pattern: Test-only behavior

Source excerpt:

```
try:
    res = svc.reserve(sale, user_id, [{'sku_code': 'SKU001', 'quantity': 5}])
    results.append(user_id)
except Exception as e:
    errors.append((user_id, type(e).__name__, str(e)))
...
self.assertEqual(len(errors), 1, ...)
self.assertTrue(has_expected, ...)
```

Why these are false positives: The worker captures its cross-thread outcome so the parent test can assert it after joining. The exception is retained and verified, not used as a direct assertion substitute or silently discarded.

Checklist evidence: This is intentional concurrent-test coordination and error collection.

### [x] Finding 18 — BP-PY-26

- Function context: `./scripts/findings/functions/18.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/django-flash-sale-inventory/inventory/views.py:20:1`
- Checklist pattern: Rule precondition

Source excerpt:

```
@csrf_exempt
@require_POST
def batch_availability(request):
    result = svc.get_batch_availability(sku_codes, region=region)
    return JsonResponse(result)
```

Why this is a false positive: The handler reads request data and returns availability data; it does not change server state.

Checklist evidence: The rule requires a state-changing view, which the enclosing function does not contain.

### [x] Finding 19 — CWE-93

- Function context: `./scripts/findings/functions/19.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/fastapi-live-metrics-ingest/app/middleware.py:13:17`
- Checklist pattern: Header/output encoding

Source excerpt:

```
duration_ms = (time.monotonic() - start) * 1000
response.headers["X-Request-Duration-Ms"] = str(int(duration_ms))
```

Why this is a false positive: The header value is an internally generated integer and cannot contain request-controlled CR or LF characters.

Checklist evidence: No untrusted control-character path reaches the header.

### [x] Finding 21 — CWE-89

- Function context: `./scripts/findings/functions/21.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/fastapi-live-metrics-ingest/app/services/aggregation.py:17:36`
- Checklist pattern: Injection

Source excerpt:

```
query = select(MetricSample.latency_ms).where(
    MetricSample.tenant_id == tenant_id,
)
result = await self.session.execute(query)
```

Why this is a false positive: The sink receives an ORM expression with bound values, not an interpolated SQL string.

Checklist evidence: The source uses a typed query builder rather than dynamic command text.

### [x] Finding 22 — CWE-1124

- Function context: `./scripts/findings/functions/22.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/fastapi-live-metrics-ingest/app/services/aggregation.py:30:1`
- Checklist pattern: Indentation/nesting

Source excerpt:

```
return {
    "p95": round(p95, 2),
    "total_samples": total_samples,
}
```

Why this is a false positive: The flagged line is a dictionary value in a return literal, not executable control-flow nesting.

Checklist evidence: Data layout accounts for the indentation.

### [x] Finding 23 — CWE-89

- Function context: `./scripts/findings/functions/23.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/fastapi-live-metrics-ingest/app/services/ingest.py:80:36`
- Checklist pattern: Injection

Source excerpt:

```
stmt = select(IngestBatch).where(IngestBatch.idempotency_key == idempotency_key)
result = await self.session.execute(stmt)
```

Why this is a false positive: The value is part of an ORM predicate rather than text interpolated into SQL.

Checklist evidence: The execute call receives a bound query expression.

### [x] Finding 24 — CWE-89

- Function context: `./scripts/findings/functions/24.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/fastapi-live-metrics-ingest/app/services/vendor_export.py:35:40`
- Checklist pattern: Injection

Source excerpt:

```
stmt = select(WindowAggregate).where(
    WindowAggregate.tenant_id == tenant_id,
    WindowAggregate.window_start >= window_start,
)
result = await self.session.execute(stmt)
```

Why this is a false positive: This is a typed query expression with comparisons, not dynamically assembled SQL text.

Checklist evidence: No concatenated or interpolated command string reaches the sink.

### [x] Finding 25 — CWE-1124

- Function context: `./scripts/findings/functions/25.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/fastapi-live-metrics-ingest/app/services/vendor_export.py:41:1`
- Checklist pattern: Indentation/nesting

Source excerpt:

```
for agg in aggregates:
    routes.append({
        "route_label": agg.route_label,
        "p50": agg.p50,
    })
```

Why this is a false positive: The reported member is part of a dictionary literal in a single loop, not excessive executable nesting.

Checklist evidence: Collection formatting creates the reported indentation.

### [x] Findings 29–30 — BP-PY-38

- Function contexts: `./scripts/findings/functions/29.txt`, `./scripts/findings/functions/30.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/fastapi-live-metrics-ingest/app/tasks.py:22:21`
- Checklist pattern: Async/task lifecycle

Source excerpt:

```
self._tasks = [
    asyncio.create_task(self._ttl_cleanup_loop(session_factory)),
    asyncio.create_task(self._vendor_export_loop(session_factory)),
]
...
await asyncio.gather(*self._tasks, return_exceptions=True)
```

Why these are false positives: Both task references are retained and later gathered during shutdown.

Checklist evidence: The task lifecycle is explicitly managed rather than discarded.

### [x] Finding 31 — CWE-1124

- Function context: `./scripts/findings/functions/31.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/fastapi-live-metrics-ingest/app/tasks.py:38:1`
- Checklist pattern: Indentation/nesting

Source excerpt:

```
while not self._shutdown_event.is_set():
    try:
        async with session_factory() as session:
            while True:
                stmt = delete(MetricSample).where(...).limit(500)
```

Why this is a false positive: CWE-1124 measures executable control-flow depth; this source has four control constructs. Class and method declarations are lexical scopes, not executable nesting, so counting them inflates the finding.

Checklist evidence: The detector must count only control-flow headers inside the function body, not class or function declarations.

### [x] Finding 32 — CWE-89

- Function context: `./scripts/findings/functions/32.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/fastapi-live-metrics-ingest/app/tasks.py:41:47`
- Checklist pattern: Injection

Source excerpt:

```
stmt = delete(MetricSample).where(
    MetricSample.created_at < cutoff
).limit(500)
result = await session.execute(stmt)
```

Why this is a false positive: The executed value is a SQLAlchemy delete expression, not a dynamically built SQL command string.

Checklist evidence: The query builder binds the comparison value.

### [x] Finding 36 — CWE-1124

- Function context: `./scripts/findings/functions/36.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/fastapi-live-metrics-ingest/tests/test_api.py:25:1`
- Checklist pattern: Indentation/nesting

Source excerpt:

```
payload = {
    "samples": [
        {"timestamp": "2024-01-01T00:00:00Z"},
    ]
}
```

Why this is a false positive: The reported indentation is nested test data, not nested executable control flow.

Checklist evidence: The construct is a test payload literal.

### [x] Finding 37 — BP-PY-46

- Function context: `./scripts/findings/functions/37.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/flask-partner-webhook-relay/app/cli.py:49:9`
- Checklist pattern: Console output

Source excerpt:

```
@app.cli.command("purge-old-data")
def purge_old_data():
    print(f"Purged {old_attempts} attempts, {old_outbox} outbox rows, {old_events} events")
```

Why this is a false positive: The print is deliberate user-facing output from a registered command.

Checklist evidence: It is not debug output in reusable runtime code.

### [x] Finding 38 — CWE-1046

- Function context: `./scripts/findings/functions/38.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/flask-partner-webhook-relay/app/cli.py:56:9`
- Checklist pattern: String construction

Source excerpt:

```
for item in items:
    item.status = "PENDING"
    count += 1
```

Why this is a false positive: The loop mutates records and increments a counter; it does not concatenate immutable text.

Checklist evidence: No string construction occurs in the loop body.

### [x] Finding 39 — BP-PY-46

- Function context: `./scripts/findings/functions/39.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/flask-partner-webhook-relay/app/cli.py:63:9`
- Checklist pattern: Console output

Source excerpt:

```
@app.cli.command("redrive-dead-letter")
def redrive_dead_letter():
    print(f"Redrove {count} items")
```

Why this is a false positive: This is a command callback reporting its result to the command user.

Checklist evidence: The output is intentional CLI behavior.

### [x] Finding 40 — CWE-1124

- Function context: `./scripts/findings/functions/40.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/flask-partner-webhook-relay/app/cli.py:76:1`
- Checklist pattern: Indentation/nesting

Source excerpt:

```
stats = {
    "in_flight": DeliveryOutbox.query.filter_by(status="IN_FLIGHT").count(),
    "dead_letter": DeliveryOutbox.query.filter_by(status="DEAD_LETTER").count(),
}
```

Why this is a false positive: The indentation is dictionary layout for command statistics, not executable nesting.

Checklist evidence: The reported line is a collection member.

### [x] Finding 41 — BP-PY-46

- Function context: `./scripts/findings/functions/41.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/flask-partner-webhook-relay/app/cli.py:79:9`
- Checklist pattern: Console output

Source excerpt:

```
@app.cli.command("queue-depth")
def queue_depth():
    print(f"Queue depth: {stats}")
```

Why this is a false positive: The print reports a command result to its user rather than debugging a library path.

Checklist evidence: It is explicit command-line output.

### [x] Finding 43 — CWE-1124

- Function context: `./scripts/findings/functions/43.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/flask-partner-webhook-relay/app/services/delivery.py:72:1`
- Checklist pattern: Indentation/nesting

Source excerpt:

```
headers={
    "Content-Type": "application/json",
    "X-Signature-256": signature,
}
```

Why this is a false positive: The reported indentation belongs to an HTTP-header dictionary passed to a call, not deep control flow.

Checklist evidence: This is a multiline literal.

### [x] Finding 46 — CWE-1046

- Function context: `./scripts/findings/functions/46.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/flask-partner-webhook-relay/app/services/delivery.py:119:13`
- Checklist pattern: String construction

Source excerpt:

```
for outbox in items:
    self.deliver(outbox)
    delivered += 1
```

Why this is a false positive: The loop processes work items and increments a number; it never concatenates text.

Checklist evidence: No immutable string construction appears in the loop.

### [x] Findings 47–48 — CWE-312, CWE-798

- Function contexts: `./scripts/findings/functions/47.txt`, `./scripts/findings/functions/48.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/flask-partner-webhook-relay/tests/conftest.py:36:1`
- Checklist pattern: Credential/secret value; Test-only behavior

Source excerpt:

```
@pytest.fixture
def sample_partner(db):
    ...
    secret="secret123",
```

Why these are false positives: The literal is synthetic data in a test fixture, not evidence of a deployed credential.

Checklist evidence: The enclosing file and fixture establish non-production test-only usage.

### [x] Finding 49 — CWE-1124

- Function context: `./scripts/findings/functions/49.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/flask-partner-webhook-relay/tests/test_api.py:23:1`
- Checklist pattern: Indentation/nesting; Test-only behavior

Source excerpt:

```
payload = {
    "event_type": "order.created",
    "payload": {"order_id": "123"},
}
```

Why this is a false positive: The highlighted line is nested test data, not an executable statement with excessive control-flow nesting.

Checklist evidence: The construct is a test payload literal.

### [x] Finding 42 — CWE-924

- Function context: `./scripts/findings/functions/42.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/flask-partner-webhook-relay/app/routes.py:19:1`
- Checklist pattern: Message integrity; authentication boundary

Source excerpt:

```
@bp.before_request
def verify_auth():
    api_key = request.headers.get("X-Api-Key")
    if api_key != current_app.config["INGEST_API_KEY"]:
        return jsonify({"error": "unauthorized"}), 401
...
def ingest_webhook():
    data = request.get_json(silent=True)
```

Why this is a false positive: The module authenticates callers with an API-key gate, while requiring an HMAC payload signature is an architectural and protocol decision. CWE-924 cannot infer that this authenticated route is an untrusted public webhook whose body must be independently signed.

Checklist evidence: A module-level authentication boundary is present; TLS and API-key enforcement can provide the intended transport and caller protection without a second application-level signature.

## Uncertain findings

No uncertain findings remain.

## True positives

| Finding IDs | Rule(s) | Source | Short evidence |
| ---: | --- | --- | --- |
| 1–5 | CWE-1052, CWE-312, CWE-798, BP-PY-13, BP-PY-22 | `django-flash-sale-inventory/flash_sale/settings.py:5–6` | A signing secret is hard-coded as a cleartext source literal. |
| 6–8 | CWE-489, CWE-756, BP-PY-21 | `django-flash-sale-inventory/flash_sale/settings.py:7–8` | Debug mode is unconditionally enabled. |
| 9–11 | CWE-1188, CWE-547, BP-PY-23 | `django-flash-sale-inventory/flash_sale/settings.py:9–10` | Host validation accepts a wildcard. |
| 20 | CWE-396 | `fastapi-live-metrics-ingest/app/routers/dashboard.py:30` | The endpoint catches every exception and returns a generic failure without preserving distinct handling. |
| 26–27 | BP-PY-1, CWE-396 | `fastapi-live-metrics-ingest/app/services/vendor_export.py:71` | A retry loop catches all exception types, including unexpected programming failures. |
| 28 | BP-PY-1 | `fastapi-live-metrics-ingest/app/services/vendor_export.py:86` | An outer generic handler converts arbitrary export failures into a failed status. |
| 33–34 | BP-PY-1, CWE-396 | `fastapi-live-metrics-ingest/app/tasks.py:46` | The maintenance worker catches every exception, logs it, and continues. |
| 35 | BP-PY-1 | `fastapi-live-metrics-ingest/app/tasks.py:54` | A generic exception handler is present around the background-loop action. |
| 44–45 | BP-PY-1, CWE-396 | `flask-partner-webhook-relay/app/services/delivery.py:93` | The final generic handler converts unexpected delivery errors into retry state. |

## Final evidence

- Delegated reviewers: `current_chunk_one`, `current_chunk_two`, `current_crosscheck`
- Chunk evidence: `./scripts/chunks`
- Function evidence: `./scripts/findings/functions`
- Validation: `git diff --check` — passed
