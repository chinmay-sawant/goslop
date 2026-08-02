# PERF-PY reference-corpus canary — 2026-08-01

> **Parent:** [PERF-PY-IMPLEMENTATION-CHECKLIST.md](./PERF-PY-IMPLEMENTATION-CHECKLIST.md) §5.2  
> **Also closes:** [PYTHON-PRECISION-HARDENING-CHECKLIST.md](../PYTHON-PRECISION-HARDENING-CHECKLIST.md) Phase 4 external canary  
> **goslop revision:** `e28dfde` (+ precision-hardening WIP on top)  
> **Corpus revision:** `codehound-python-perf-targets@82e07b9`

## Command

```bash
CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop

GOSLOP=./bin/goslop
CFG=templates/goslop-python.toml   # languages=["python"]
ROOT=/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets

for app in django-flash-sale-inventory fastapi-live-metrics-ingest flask-partner-webhook-relay; do
  "$GOSLOP" --config "$CFG" --profile all --only 'PERF-PY-*' \
    --format json --no-fail --no-cache --no-baseline "$ROOT/$app"
done
```

## Summary counts

| App | Files scanned (stderr) | PERF-PY findings | Distinct rules |
|-----|------------------------:|-----------------:|----------------|
| django-flash-sale-inventory | 21 | 14 | PERF-PY-2,3,4,12 |
| fastapi-live-metrics-ingest | (app tree) | 2 | PERF-PY-1,18 |
| flask-partner-webhook-relay | (app tree) | 1 | PERF-PY-5 |
| **Total** | | **17** | |

No project-wide spray across unrelated modules. Findings concentrate on the designed hot paths (reservation service, aggregation/ingest, delivery worker).

## CLI smoke (Phase 5.3)

```bash
# materialized tests/fixtures/python/perf/PERF-PY-6-vulnerable.txt → worker.py
./bin/goslop --config /tmp/goslop-py-smoke.toml --profile all --only PERF-PY-6 \
  --format text --no-fail --no-cache --no-baseline /tmp/goslop-py-fixture-smoke
```

**Result:** exactly one finding — `PERF-PY-6` on `worker.py:2`.

## False-positive triage

Disposition legend: **expected** (corpus designed smell) · **accepted-review** (true-ish but may need identity/dependency depth) · **false-positive** (should not fire) · **fixed** (detector changed)

### django-flash-sale-inventory (14)

| Rule | Path:line | Disposition | Reason |
|------|-----------|-------------|--------|
| PERF-PY-2 | `inventory/services/reservation.py:38` | expected | `Sku.objects.get` inside `for item in items` |
| PERF-PY-2 | `…/reservation.py:99` | expected | `WarehouseStock.objects.get` inside confirm loop |
| PERF-PY-2 | `…/reservation.py:120` | expected | cancel loop stock lookup |
| PERF-PY-2 | `…/reservation.py:146` | expected | expire loop stock lookup |
| PERF-PY-3 | `…/reservation.py:68` | accepted-review | per-line `ReservationLine.objects.create` after parent create; parent-link dependency not always visible to source heuristic |
| PERF-PY-3 | `…/reservation.py:77` | accepted-review | `StockLedger.objects.create` in same batch loop |
| PERF-PY-3 | `…/reservation.py:103` | accepted-review | confirm-path ledger create |
| PERF-PY-3 | `…/reservation.py:123` | accepted-review | cancel-path ledger create |
| PERF-PY-3 | `…/reservation.py:149` | accepted-review | expire-path ledger create |
| PERF-PY-4 | `…/reservation.py:100` | expected | `reserved_quantity -=` then save |
| PERF-PY-4 | `…/reservation.py:101` | expected | `quantity -=` then save |
| PERF-PY-4 | `…/reservation.py:121` | expected | cancel RMW |
| PERF-PY-4 | `…/reservation.py:147` | expected | expire RMW |
| PERF-PY-12 | `inventory/views.py:23` | expected | `json.loads(request.body)` without visible body bound |

### fastapi-live-metrics-ingest (2)

| Rule | Path:line | Disposition | Reason |
|------|-----------|-------------|--------|
| PERF-PY-1 | `app/services/aggregation.py:18` | expected | `.scalars().all()` then `sorted(rows)` for percentiles |
| PERF-PY-18 | `app/services/ingest.py:19` | accepted-review | intentional two-stage route normalize (`re.sub` ×2); useful review signal, not CI-hard-fail grade |

### flask-partner-webhook-relay (1)

| Rule | Path:line | Disposition | Reason |
|------|-----------|-------------|--------|
| PERF-PY-5 | `app/services/delivery.py:117` | expected | `claim_work` batch then synchronous `deliver` in loop |

## Maturity decision (2026-08-01)

| Decision | Rationale |
|----------|-----------|
| Keep **all** `PERF-PY-*` at `MaturityExperimental` | Canary is focused and useful, but PERF-PY-3 dependent-create and PERF-PY-18 staged-regex remain review-level noise |
| **Do not** add `PerfTierSPY` / `PerfTierAPY` | No rule cleared “reliable production / CI-hard-fail” bar |
| **Do not** widen `recommended` / `perf` packs to `PERF-*` or `PERF-PY-*` | Profile contract tests continue to deny PERF-PY under those packs |
| Promote surface | Opt-in: `languages=["python"]` + `--profile all` + `--only PERF-PY-*` (or slim in-repo corpus) |

## Follow-ups (not blocking this closure)

1. Tighten PERF-PY-3 when create uses an outer-scope parent object created before the loop (reduce accepted-review ledger creates).
2. Optional deliberate-stage suppress for PERF-PY-18 when consecutive `re.sub` clearly rewrite toward a normalized route template.
3. Re-run this canary after those tightenings before any tier-list promotion.
