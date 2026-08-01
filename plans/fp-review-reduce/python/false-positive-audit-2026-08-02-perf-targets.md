# False-positive audit report

## Run metadata

```yaml
timestamp: 2026-08-01T21:05:00Z
repository: goslop
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop
branch: feat/python-perf-ruleset-plan
commit: 5217d65
scan_target: /home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets
chunk_path: ./scripts/chunks
function_context_path: ./scripts/findings/functions
```

## Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop`
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config "templates/goslop-python.toml" --export-context --export-chunks --no-cache "/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets"` (via `make run-python`)
- Findings: `45` (116 files scanned; ~125ms)
- Chunks reviewed: `./scripts/chunks/Chunk_1_25.txt`, `./scripts/chunks/Chunk_26_45.txt`
- Function contexts reviewed: `./scripts/findings/functions/1.txt` … `45.txt` (full read for FP/Uncertain; chunk context for TP clusters; source follow-up for 27, 37, 38, 39, 42)

## Audit checklist

- [x] Read every assigned chunk under `./scripts/chunks`.
- [x] Read `./scripts/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews ([1–25](fc71b9d2-5061-4d07-a24c-082b59c3b940) all TP; [26–45](f86145a8-ab9f-4320-b9e5-ce6702bb8970) FP on 38+39). Parent initially TP’d 38; reconciled to FP on stated CWE-91 condition (parse ≠ format-into-XPath). Finding 42 (CWE-328) resolved Uncertain→TP: password MD5 meets narrowed `securityHashContext` rule condition; PDF Algorithm 3 mandate is protocol knowledge outside FP criteria.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 2 | 38, 39 |
| True positive | 43 | 1–37, 40–45 |
| Uncertain | 0 | — |

User belief that most are true positives is confirmed: **43/45 TP**, with two FPs (CWE-91 parse-only; PERF-PY-26 CLI one-shot). Former uncertain 42 (CWE-328) resolved to TP.

## Classification table (all 45)

| ID | Rule | Classification | One-line rationale |
| --- | --- | --- | --- |
| 1 | CWE-1052 | True positive | Hard-coded `SECRET_KEY` literal in settings init |
| 2 | CWE-312 | True positive | Cleartext secret literal in source |
| 3 | CWE-798 | True positive | Hard-coded credential string |
| 4 | BP-PY-13 | True positive | Hardcoded secret-like string |
| 5 | BP-PY-22 | True positive | Hardcoded Django `SECRET_KEY` |
| 6 | BP-PY-21 | True positive | `DEBUG = True` in settings |
| 7 | CWE-489 | True positive | Debug mode enabled in application source |
| 8 | CWE-756 | True positive | Debug error output enabled |
| 9 | CWE-1188 | True positive | `ALLOWED_HOSTS = ['*']` wildcard default |
| 10 | CWE-547 | True positive | Hard-coded insecure host validation |
| 11 | BP-PY-23 | True positive | `ALLOWED_HOSTS` uses `'*'` |
| 12 | PERF-PY-2 | True positive | `Sku.objects.get` inside `for item in items` |
| 13 | PERF-PY-3 | True positive | `ReservationLine.objects.create` per batch item |
| 14 | PERF-PY-3 | True positive | `StockLedger.objects.create` per batch item |
| 15 | PERF-PY-2 | True positive | `WarehouseStock.objects.get` inside line loop |
| 16 | PERF-PY-4 | True positive | RMW `reserved_quantity -=` then save |
| 17 | PERF-PY-4 | True positive | RMW `quantity -=` then save |
| 18 | PERF-PY-3 | True positive | `StockLedger.objects.create` per line |
| 19 | PERF-PY-2 | True positive | ORM get inside reservation line loop |
| 20 | PERF-PY-4 | True positive | RMW reserved counter then save |
| 21 | PERF-PY-3 | True positive | Per-row ledger create |
| 22 | PERF-PY-2 | True positive | ORM get inside nested line loop |
| 23 | PERF-PY-4 | True positive | RMW reserved counter then save |
| 24 | PERF-PY-3 | True positive | Per-row ledger create |
| 25 | PERF-PY-12 | True positive | `json.loads(request.body)` with no size bound |
| 26 | CWE-396 | True positive | `except Exception:` before HTTPException rewrite |
| 27 | PERF-PY-1 | True positive | `.all()` then `sorted(rows)` for percentiles |
| 28 | PERF-PY-18 | True positive | Two `.sub` rewrites of `label` in `normalize_route` |
| 29 | BP-PY-1 | True positive | Broad `except Exception as e` |
| 30 | CWE-396 | True positive | Same generic Exception handler |
| 31 | BP-PY-1 | True positive | Outer broad `except Exception` |
| 32 | BP-PY-1 | True positive | Broad except in TTL cleanup |
| 33 | CWE-396 | True positive | Same generic Exception handler |
| 34 | BP-PY-1 | True positive | Broad except around export loop tick |
| 35 | BP-PY-1 | True positive | Broad except in webhook deliver |
| 36 | CWE-396 | True positive | Same generic Exception handler |
| 37 | PERF-PY-5 | True positive | `for outbox in items: self.deliver(outbox)` after claim |
| 38 | CWE-91 | False positive | `ET.fromstring` parses XML; does not format into XPath/XML |
| 39 | PERF-PY-26 | False positive | One-shot CLI `main()` parse; not a hot path |
| 40 | CWE-88 | True positive | Dynamic `flavour`/`pdf` in subprocess argv |
| 41 | BP-PY-46 | True positive | Library `print` under `ENGINE_DEBUG_BUFFERS` |
| 42 | CWE-328 | True positive | `hashlib.md5` on `owner_password` in key-derivation context |
| 43 | PERF-PY-24 | True positive | Duplicate `wrap_text` measure before draw |
| 44 | BP-PY-1 | True positive | Bare `except Exception:` swallows decode failures |
| 45 | CWE-396 | True positive | Generic Exception catch hides distinct failures |

## False positives

### [ ] Finding 38 — CWE-91

- Function context: `./scripts/findings/functions/38.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/pythoncoreengine/compliance/verapdf/report.py:17:12`
- Checklist pattern: dynamic value formatted into an XML or XPath expression

Source excerpt:

```
def parse_report(xml_data: str) -> tuple[bool, int, int, list[str]]:
    root = ET.fromstring(xml_data)
```

Why this is a false positive: Rule condition / message is XML or XPath **injection** (value formatted into markup or an XPath query; fix = builders / bound variables). This site **parses** a document with `fromstring`; it does not interpolate into XML or XPath text.

Checklist evidence: Parse API with a whole-document argument, not `f"…"` / concat into `.xpath(...)` or XML construction.

### [ ] Finding 39 — PERF-PY-26

- Function context: `./scripts/findings/functions/39.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/pythoncoreengine/compliance/verapdf/report.py:63:1`
- Checklist pattern: expensive parse not on request/job/render loop path

Source excerpt:

```
def main() -> int:
    xml_data = sys.stdin.read()
    if not xml_data.strip():
        print("FAILED: no XML data on stdin", file=sys.stderr)
        return 1

    try:
        compliant, passed, failed, errors = parse_report(xml_data)
    except ET.ParseError as exc:
        ...
```

Why this is a false positive: Rule requires an expensive decode/parse on a hot path (`handle_job` / `handle_request` / `render` / `build_` / `process`, or inside a loop). This site is a one-shot CLI `main()` that parses stdin once; the surrounding function has none of those hot-path markers and is not in a loop.

Checklist evidence: No hot-path function name and no in-function loop around `parse_report` (earlier loops live in a different function and do not enclose this call).

## Uncertain findings

None.

## True positives (compact by rule)

| Rule | IDs | Pattern |
| --- | --- | --- |
| Django settings secrets/debug/hosts | 1–11 | Hardcoded `SECRET_KEY`, `DEBUG=True`, `ALLOWED_HOSTS=['*']` |
| PERF-PY-2 | 12, 15, 19, 22 | ORM `.get` inside item/line loops |
| PERF-PY-3 | 13–14, 18, 21, 24 | Per-row `.create` in batch loops |
| PERF-PY-4 | 16–17, 20, 23 | Read-modify-write counters then `.save()` |
| PERF-PY-12 | 25 | Unbounded `json.loads(request.body)` |
| CWE-396 / BP-PY-1 | 26, 29–36, 44–45 | Broad `except Exception` |
| PERF-PY-1 | 27 | `.all()` then Python `sorted` for percentiles |
| PERF-PY-18 | 28 | Chained regex `.sub` on same `label` |
| PERF-PY-5 | 37 | Sequential `deliver` over claimed batch |
| CWE-88 | 40 | Dynamic subprocess argv (`flavour`, `pdf`) |
| BP-PY-46 | 41 | Library debug `print` |
| CWE-328 | 42 | `hashlib.md5(_pad_password_32(owner_password))` in `_compute_o_r4` — password/key-derivation security context (narrowed rule condition). PDF Algorithm 3 mandating MD5 is protocol knowledge, not an FP under the audit skill. |
| PERF-PY-24 | 43 | Duplicate `wrap_text` measure |

## Final evidence

- Delegated reviewers: [1–25](fc71b9d2-5061-4d07-a24c-082b59c3b940), [26–45](f86145a8-ab9f-4320-b9e5-ce6702bb8970); finding 38 reconciled to FP after 26–45 review
- Chunk evidence: `./scripts/chunks`
- Function evidence: `./scripts/findings/functions`
- Validation: `git diff --check` — pass
