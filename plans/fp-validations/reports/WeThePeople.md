# False-positive audit: WeThePeople

## Run metadata

```yaml
timestamp: 2026-08-02T07:55:45Z
repository: WeThePeople
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople
branch: main
commit: 6acbd5b2a67d4499ed17a05ea48cf9aebd3d1da0
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople
chunk_path: ./scripts/chunks
function_context_path: ./scripts/findings/functions
```

## Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop` (binary present as `./bin/goslop`)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir real-repos/WeThePeople/scripts/chunks -context-dir real-repos/WeThePeople/scripts/findings/functions real-repos/WeThePeople`
- Findings: `1492`
- Chunks reviewed: `./scripts/chunks/Chunk_1_25.txt` … `./scripts/chunks/Chunk_1476_1492.txt` (60 files)
- Function contexts reviewed: `./scripts/findings/functions/<finding-id>.txt` for every proposed false positive

Note: several rule IDs emitted by the scan (BP-PY-*, CWE-1084, CWE-1121, CWE-1124, CWE-117, CWE-396, CWE-390, CWE-1071, CWE-1341, CWE-93, CWE-94, CWE-88, CWE-1046, CWE-186, CWE-260, CWE-215, CWE-359, CWE-290, CWE-829) are not present in the `-explain` catalogue of the checked-out `./bin/goslop` build (the scan was produced by a newer build). For those rules the decision is based on the rule title/message and the shown source excerpt, as documented per finding.

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
| False positive | 70 | 2, 193, 197, 206, 216, 235, 237, 248, 252, 258, 259, 378, 386, 400, 436, 484, 485, 489, 490, 493, 495, 496, 511, 515, 563, 577, 599, 619, 638, 648, 670, 685, 762, 764, 790, 872, 873, 895, 896, 897, 899, 905, 906, 912, 924, 925, 969, 974, 1033, 1059, 1071, 1076, 1093, 1094, 1096, 1097, 1112, 1160, 1164, 1165, 1206, 1279, 1359, 1377, 1413, 1455, 1460, 1470, 1473, 1474 |
| True positive | 1422 | every remaining finding (1–1492 minus the FP IDs above); enumerated per rule in the `## True positives` tables below (includes former uncertain CWE-1084 IDs 174, 285, 440, 451, 977, 1110, 1174, 1204) |
| Uncertain | 0 | — |

## False positives

One subsection per finding (or per exact-same-construct group). All 70 false positives are on distinct source lines, so no grouping was applied.

### [ ] Finding 2 — CWE-89

- Function context: `./scripts/findings/functions/2.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/alembic/versions/20260505_anomaly_dedupe_index.py:52:7`
- Checklist pattern: `Static SQL literal (no interpolation) reaching execute`

Source excerpt:

```
op.execute("CREATE INDEX IF NOT EXISTS ix_anomalies_dedupe ...")
```

Why this is a false positive: The SQL passed to op.execute is a static string literal — no interpolation or runtime data is involved.

Checklist evidence: Static SQL literal (no interpolation) reaching execute — verified against the shown source.

### [ ] Finding 193 — CWE-1341

- Function context: `./scripts/findings/functions/193.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/backfill_company_logos.py:270:13`
- Checklist pattern: `Single close per control-flow path`

Source excerpt:

```
except ImportError:
            log.error("anthropic SDK not installed; pip install anthropic")
            conn.close()
            return 1
```

Why this is a false positive: conn.close() at this line is on the ImportError return path; the other close() calls (lines 275, 335) sit on disjoint branches that return first — no path releases the handle twice.

Checklist evidence: Single close per control-flow path — verified against the shown source.

### [ ] Finding 197 — PERF-PY-26

- Function context: `./scripts/findings/functions/197.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/backfill_company_logos.py:346:1`
- Checklist pattern: `No expensive decode/parse at the flagged site`

Source excerpt:

```
args = parser.parse_args()
```

Why this is a false positive: The flagged statement is argparse parse_args() in main() — CLI argument parsing, not an expensive decode/parse on a hot path.

Checklist evidence: No expensive decode/parse at the flagged site — verified against the shown source.

### [ ] Finding 206 — PERF-PY-26

- Function context: `./scripts/findings/functions/206.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/backfill_logos_wikidata.py:237:1`
- Checklist pattern: `No expensive decode/parse at the flagged site`

Source excerpt:

```
args = p.parse_args()
```

Why this is a false positive: Same: parse_args() at CLI entry, not a hot path.

Checklist evidence: No expensive decode/parse at the flagged site — verified against the shown source.

### [ ] Finding 216 — PERF-PY-26

- Function context: `./scripts/findings/functions/216.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/backfill_logos_wikipedia.py:221:1`
- Checklist pattern: `No expensive decode/parse at the flagged site`

Source excerpt:

```
args = p.parse_args()
```

Why this is a false positive: Same: parse_args() at CLI entry.

Checklist evidence: No expensive decode/parse at the flagged site — verified against the shown source.

### [ ] Finding 235 — PERF-PY-26

- Function context: `./scripts/findings/functions/235.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/backfill_sanctions_status.py:380:1`
- Checklist pattern: `No expensive decode/parse at the flagged site`

Source excerpt:

```
args = parser.parse_args()
```

Why this is a false positive: Same: parse_args() at CLI entry.

Checklist evidence: No expensive decode/parse at the flagged site — verified against the shown source.

### [ ] Finding 237 — CWE-89

- Function context: `./scripts/findings/functions/237.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/backfill_stock_fundamentals.py:68:15`
- Checklist pattern: `Static SQL with bound parameters reaching execute`

Source excerpt:

```
conn.execute("SELECT 1 FROM stock_fundamentals " "WHERE entity_type = ? AND entity_id = ? AND snapshot_date = ? LIMIT 1", (entity_type, entity_id, _today_iso()))
```

Why this is a false positive: The SQL string is a static literal with '?' placeholders bound to a parameter tuple; no dynamic text is interpolated.

Checklist evidence: Static SQL with bound parameters reaching execute — verified against the shown source.

### [ ] Finding 248 — CWE-89

- Function context: `./scripts/findings/functions/248.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/backfill_verification_tier.py:70:19`
- Checklist pattern: `Static SQL with bind parameters reaching execute`

Source excerpt:

```
text("UPDATE stories SET verification_tier = :tier, verification_score = :score, updated_at = :now WHERE id = :id")
```

Why this is a false positive: The statement is a static literal with :tier/:score/:now/:id bind params; no interpolation.

Checklist evidence: Static SQL with bind parameters reaching execute — verified against the shown source.

### [ ] Finding 252 — CWE-89

- Function context: `./scripts/findings/functions/252.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/correct_lobby_double_count_stories.py:149:13`
- Checklist pattern: `Static SQL with bind parameters reaching execute`

Source excerpt:

```
text("SELECT id, title, body, evidence, correction_history FROM stories WHERE id = :id")
```

Why this is a false positive: Static SELECT with a single :id bind parameter; nothing is formatted into the string.

Checklist evidence: Static SQL with bind parameters reaching execute — verified against the shown source.

### [ ] Finding 258 — CWE-909

- Function context: `./scripts/findings/functions/258.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/detect_stories.py:175:38`
- Checklist pattern: `DB resource injected as a parameter`

Source excerpt:

```
def _verify_story_numbers(db, story):
```

Why this is a false positive: db is a function parameter supplied by the caller; the resource is initialized before use, so 'missing initialization' does not hold.

Checklist evidence: DB resource injected as a parameter — verified against the shown source.

### [ ] Finding 259 — CWE-89

- Function context: `./scripts/findings/functions/259.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/detect_stories.py:198:25`
- Checklist pattern: `Static SQL with bind parameters reaching execute`

Source excerpt:

```
text("SELECT COUNT(*) FROM congressional_trades WHERE person_id = :eid")
```

Why this is a false positive: Static SELECT with :eid bind parameter; no dynamic text.

Checklist evidence: Static SQL with bind parameters reaching execute — verified against the shown source.

### [ ] Finding 378 — PERF-PY-26

- Function context: `./scripts/findings/functions/378.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/dump_public_snapshot.py:152:1`
- Checklist pattern: `No expensive decode/parse at the flagged site`

Source excerpt:

```
args = ap.parse_args()
```

Why this is a false positive: Same: parse_args() at CLI entry.

Checklist evidence: No expensive decode/parse at the flagged site — verified against the shown source.

### [ ] Finding 386 — PERF-PY-26

- Function context: `./scripts/findings/functions/386.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/evaluate_legislative_claims.py:224:1`
- Checklist pattern: `No expensive decode/parse at the flagged site`

Source excerpt:

```
args = parser.parse_args()
```

Why this is a false positive: Same: parse_args() at CLI entry.

Checklist evidence: No expensive decode/parse at the flagged site — verified against the shown source.

### [ ] Finding 400 — CWE-186

- Function context: `./scripts/findings/functions/400.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/generate_under_standards.py:314:19`
- Checklist pattern: `Extraction regex, not a validation regex`

Source excerpt:

```
_TIME_WINDOW_RE = re.compile(
    r"\b("
```

Why this is a false positive: The regex is an extraction pattern applied to story text via re.search, not an input-validation regex; the 'overly restrictive validation' condition does not apply.

Checklist evidence: Extraction regex, not a validation regex — verified against the shown source.

### [ ] Finding 436 — CWE-89

- Function context: `./scripts/findings/functions/436.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/migrate_add_indexes.py:131:17`
- Checklist pattern: `Static SQL list executed verbatim`

Source excerpt:

```
for sql in indexes:
    conn.execute(sql)
```

Why this is a false positive: sql is a static string literal drawn from the module-level 'indexes' list; conn.execute(sql) executes a constant, not dynamic text.

Checklist evidence: Static SQL list executed verbatim — verified against the shown source.

### [ ] Finding 484 — CWE-829

- Function context: `./scripts/findings/functions/484.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/rebuild_search_index.py:124:19`
- Checklist pattern: `Dynamic import fed from a static allowlist`

Source excerpt:

```
for module_path, cls_name, ... in _COMPANY_TABLES:
            mod = importlib.import_module(module_path)
```

Why this is a false positive: import_module is fed from _COMPANY_TABLES, a static tuple of package module names defined in the same file — no untrusted input reaches the sink.

Checklist evidence: Dynamic import fed from a static allowlist — verified against the shown source.

### [ ] Finding 485 — CWE-94

- Function context: `./scripts/findings/functions/485.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/rebuild_search_index.py:124:19`
- Checklist pattern: `Dynamic-import sink fed from a static allowlist`

Source excerpt:

```
for module_path, cls_name, ... in _COMPANY_TABLES:
            mod = importlib.import_module(module_path)
```

Why this is a false positive: Same static allowlist source; no request-derived value reaches import_module.

Checklist evidence: Dynamic-import sink fed from a static allowlist — verified against the shown source.

### [ ] Finding 489 — CWE-89

- Function context: `./scripts/findings/functions/489.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/rebuild_search_index.py:159:15`
- Checklist pattern: `Static SQL literal reaching execute`

Source excerpt:

```
db.execute(text("DELETE FROM entity_search"))
```

Why this is a false positive: text("DELETE FROM entity_search") is a fully static literal.

Checklist evidence: Static SQL literal reaching execute — verified against the shown source.

### [ ] Finding 490 — CWE-909

- Function context: `./scripts/findings/functions/490.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/rebuild_search_index.py:163:37`
- Checklist pattern: `DB resource initialized in enclosing scope`

Source excerpt:

```
def _ingest(name: str, gen):
```

Why this is a false positive: db is captured from the enclosing scope where it is created by SessionLocal() before _ingest runs.

Checklist evidence: DB resource initialized in enclosing scope — verified against the shown source.

### [ ] Finding 493 — CWE-89

- Function context: `./scripts/findings/functions/493.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/repair_correction_notice_substitution.py:101:13`
- Checklist pattern: `Static SQL with bind parameters reaching execute`

Source excerpt:

```
text("SELECT id, title, body, correction_history ...")
```

Why this is a false positive: Static SELECT with :id bind parameter; no interpolation.

Checklist evidence: Static SQL with bind parameters reaching execute — verified against the shown source.

### [ ] Finding 495 — CWE-909

- Function context: `./scripts/findings/functions/495.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/retract_misattributed_stories.py:105:23`
- Checklist pattern: `DB resource injected as a parameter`

Source excerpt:

```
def ensure_tables(db):
```

Why this is a false positive: db is a parameter; tables are created via the parameterized resource.

Checklist evidence: DB resource injected as a parameter — verified against the shown source.

### [ ] Finding 496 — CWE-89

- Function context: `./scripts/findings/functions/496.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/retract_misattributed_stories.py:108:11`
- Checklist pattern: `Static SQL literal reaching execute`

Source excerpt:

```
text("CREATE TABLE IF NOT EXISTS story_corrections (")
```

Why this is a false positive: CREATE TABLE statement is a static literal with no interpolation.

Checklist evidence: Static SQL literal reaching execute — verified against the shown source.

### [ ] Finding 511 — CWE-89

- Function context: `./scripts/findings/functions/511.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/run_pipeline.py:77:17`
- Checklist pattern: `Static SQL literal reaching execute`

Source excerpt:

```
text("SELECT COALESCE(published_at, created_at) AS ts FROM stories WHERE status IN ('draft', 'published') ...")
```

Why this is a false positive: Static SELECT with no interpolated values.

Checklist evidence: Static SQL literal reaching execute — verified against the shown source.

### [ ] Finding 515 — BP-PY-7

- Function context: `./scripts/findings/functions/515.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/scheduler.py:731:22`
- Checklist pattern: `os.open() with explicit fd close, not builtin open()`

Source excerpt:

```
self._fd = os.open(str(self.lock_path), os.O_CREAT | os.O_RDWR)
```

Why this is a false positive: os.open() returns a raw fd that is explicitly closed via os.close() in release(); the rule's 'open without with' pattern targets the builtin open() file object.

Checklist evidence: os.open() with explicit fd close, not builtin open() — verified against the shown source.

### [ ] Finding 563 — BP-PY-36

- Function context: `./scripts/findings/functions/563.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/sync_agriculture_enforcement.py:246:1`
- Checklist pattern: `Session explicitly closed before exit`

Source excerpt:

```
session = Session()
    ...
    session.close()
```

Why this is a false positive: main() creates session = Session() but closes it with session.close() at the end of the function.

Checklist evidence: Session explicitly closed before exit — verified against the shown source.

### [ ] Finding 577 — BP-PY-36

- Function context: `./scripts/findings/functions/577.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/sync_chemicals_enforcement.py:246:1`
- Checklist pattern: `Session explicitly closed before exit`

Source excerpt:

```
session = Session()
    ...
    session.close()
```

Why this is a false positive: Same pattern: session.close() is called at the end of main().

Checklist evidence: Session explicitly closed before exit — verified against the shown source.

### [ ] Finding 599 — BP-PY-36

- Function context: `./scripts/findings/functions/599.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/sync_defense_enforcement.py:310:1`
- Checklist pattern: `Session explicitly closed before exit`

Source excerpt:

```
session = Session()
    ...
    session.close()
```

Why this is a false positive: Same pattern: session.close() at end of main().

Checklist evidence: Session explicitly closed before exit — verified against the shown source.

### [ ] Finding 619 — BP-PY-36

- Function context: `./scripts/findings/functions/619.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/sync_education_enforcement.py:265:1`
- Checklist pattern: `Session explicitly closed before exit`

Source excerpt:

```
session = Session()
    ...
    session.close()
```

Why this is a false positive: Same pattern: session.close() at end of main().

Checklist evidence: Session explicitly closed before exit — verified against the shown source.

### [ ] Finding 638 — BP-PY-36

- Function context: `./scripts/findings/functions/638.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/sync_energy_enforcement.py:246:1`
- Checklist pattern: `Session explicitly closed before exit`

Source excerpt:

```
session = Session()
    ...
    session.close()
```

Why this is a false positive: Same pattern: session.close() at end of main().

Checklist evidence: Session explicitly closed before exit — verified against the shown source.

### [ ] Finding 648 — BP-PY-36

- Function context: `./scripts/findings/functions/648.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/sync_finance_enforcement.py:249:1`
- Checklist pattern: `Session explicitly closed before exit`

Source excerpt:

```
session = Session()
    ...
    session.close()
```

Why this is a false positive: Same pattern: session.close() at end of main().

Checklist evidence: Session explicitly closed before exit — verified against the shown source.

### [ ] Finding 670 — BP-PY-36

- Function context: `./scripts/findings/functions/670.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/sync_health_enforcement.py:246:1`
- Checklist pattern: `Session explicitly closed before exit`

Source excerpt:

```
session = Session()
    ...
    session.close()
```

Why this is a false positive: Same pattern: session.close() at end of main().

Checklist evidence: Session explicitly closed before exit — verified against the shown source.

### [ ] Finding 685 — CWE-89

- Function context: `./scripts/findings/functions/685.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/sync_insider_trades.py:284:24`
- Checklist pattern: `Static SQL with bind parameters reaching execute`

Source excerpt:

```
text("SELECT 1 FROM sec_insider_trades WHERE dedupe_hash=:h")
```

Why this is a false positive: Static SELECT with :h bind parameter.

Checklist evidence: Static SQL with bind parameters reaching execute — verified against the shown source.

### [ ] Finding 762 — BP-PY-36

- Function context: `./scripts/findings/functions/762.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/sync_telecom_enforcement.py:265:1`
- Checklist pattern: `Session explicitly closed before exit`

Source excerpt:

```
session = Session()
    ...
    session.close()
```

Why this is a false positive: Same pattern: session.close() at end of main().

Checklist evidence: Session explicitly closed before exit — verified against the shown source.

### [ ] Finding 764 — CWE-186

- Function context: `./scripts/findings/functions/764.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/sync_trades_from_disclosures.py:71:17`
- Checklist pattern: `Extraction regex, not a validation regex`

Source excerpt:

```
ASSET_TYPE_RE = re.compile(r"\[([A-Z]{2})\]")
```

Why this is a false positive: ASSET_TYPE_RE is used with .search() to parse asset-type codes out of disclosure text — a parser, not a validator.

Checklist evidence: Extraction regex, not a validation regex — verified against the shown source.

### [ ] Finding 790 — BP-PY-36

- Function context: `./scripts/findings/functions/790.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/jobs/sync_transportation_enforcement.py:318:1`
- Checklist pattern: `Session explicitly closed before exit`

Source excerpt:

```
session = Session()
    ...
    session.close()
```

Why this is a false positive: Same pattern: session.close() at end of main().

Checklist evidence: Session explicitly closed before exit — verified against the shown source.

### [ ] Finding 872 — CWE-93

- Function context: `./scripts/findings/functions/872.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/middleware/rate_limit_headers.py:71:17`
- Checklist pattern: `Static header value (no dynamic input)`

Source excerpt:

```
response.headers["RateLimit-Limit"] = str(_GLOBAL_LIMIT)
```

Why this is a false positive: The header value is str(_GLOBAL_LIMIT), a static integer constant — not a dynamic/attacker-controllable string.

Checklist evidence: Static header value (no dynamic input) — verified against the shown source.

### [ ] Finding 873 — CWE-93

- Function context: `./scripts/findings/functions/873.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/middleware/security.py:121:17`
- Checklist pattern: `Static header value (no dynamic input)`

Source excerpt:

```
response.headers["Permissions-Policy"] = ("camera=(), microphone=(), geolocation=()"
```

Why this is a false positive: The header value is a fixed string literal.

Checklist evidence: Static header value (no dynamic input) — verified against the shown source.

### [ ] Finding 895 — BP-PY-14

- Function context: `./scripts/findings/functions/895.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/routers/auth.py:1314:23`
- Checklist pattern: `No requests call at the flagged site`

Source excerpt:

```
user_id = session.get("client_reference_id") or session.get("metadata", {}).get("user_id")
```

Why this is a false positive: The flagged expression is a dict .get() on a Stripe session object in the webhook handler; this file contains no requests call at all.

Checklist evidence: No requests call at the flagged site — verified against the shown source.

### [ ] Finding 896 — BP-PY-14

- Function context: `./scripts/findings/functions/896.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/routers/auth.py:1314:61`
- Checklist pattern: `No requests call at the flagged site`

Source excerpt:

```
user_id = session.get("client_reference_id") or session.get("metadata", {}).get("user_id")
```

Why this is a false positive: Same dict .get() expression, misattributed to requests.

Checklist evidence: No requests call at the flagged site — verified against the shown source.

### [ ] Finding 897 — BP-PY-14

- Function context: `./scripts/findings/functions/897.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/routers/auth.py:1320:17`
- Checklist pattern: `No requests call at the flagged site`

Source excerpt:

```
target_role = (
                session.get("metadata", {}).get("role")
                or "enterprise")
```

Why this is a false positive: Same: session.get(...) dict access.

Checklist evidence: No requests call at the flagged site — verified against the shown source.

### [ ] Finding 899 — BP-PY-14

- Function context: `./scripts/findings/functions/899.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/routers/auth.py:1333:50`
- Checklist pattern: `No requests call at the flagged site`

Source excerpt:

```
logger.info("User %s upgraded to %s via Stripe (session=%s)",
                        user.email, target_role, session.get("id"),)
```

Why this is a false positive: Same: logger call with dict access, no requests call.

Checklist evidence: No requests call at the flagged site — verified against the shown source.

### [ ] Finding 905 — CWE-909

- Function context: `./scripts/findings/functions/905.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/routers/civic.py:121:67`
- Checklist pattern: `DB resource injected as a parameter`

Source excerpt:

```
def _update_scores(db: Session, target_type: str, target_id: int):
```

Why this is a false positive: db is a typed parameter (db: Session).

Checklist evidence: DB resource injected as a parameter — verified against the shown source.

### [ ] Finding 906 — CWE-89

- Function context: `./scripts/findings/functions/906.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/routers/civic.py:149:15`
- Checklist pattern: `Static SQL literal reaching execute`

Source excerpt:

```
db.execute(text("BEGIN IMMEDIATE"))
```

Why this is a false positive: BEGIN IMMEDIATE is a static SQL keyword string.

Checklist evidence: Static SQL literal reaching execute — verified against the shown source.

### [ ] Finding 912 — CWE-1121

- Function context: `./scripts/findings/functions/912.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/routers/civic.py:188:67`
- Checklist pattern: `Branch count below the ≥12 threshold`

Source excerpt:

```
def _check_badge_progress(db: Session, user_id: int, action: str):
    ...
    for slug in slugs:
        badge = db.query(Badge)...
        if not badge: continue
        if existing: ... continue
        if action == "vote": ... elif ... 
```

Why this is a false positive: The function body contains only 10 control-flow branches (1 for, 2 if, 6 if/elif, 1 if), below the 'at least twelve' threshold.

Checklist evidence: Branch count below the ≥12 threshold — verified against the shown source.

### [ ] Finding 924 — CWE-909

- Function context: `./scripts/findings/functions/924.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/routers/common.py:54:49`
- Checklist pattern: `DB resource provided by dependency injection`

Source excerpt:

```
def health_check(db: Session = Depends(get_db)):
```

Why this is a false positive: db comes from FastAPI dependency injection (Depends(get_db)) which initializes it.

Checklist evidence: DB resource provided by dependency injection — verified against the shown source.

### [ ] Finding 925 — CWE-89

- Function context: `./scripts/findings/functions/925.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/routers/common.py:67:11`
- Checklist pattern: `Static SQL literal reaching execute`

Source excerpt:

```
db.execute(text("SELECT 1"))
```

Why this is a false positive: SELECT 1 is a static literal.

Checklist evidence: Static SQL literal reaching execute — verified against the shown source.

### [ ] Finding 969 — CWE-88

- Function context: `./scripts/findings/functions/969.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/routers/ops.py:325:20`
- Checklist pattern: `Subprocess argv from static registry with existence check`

Source excerpt:

```
proc = subprocess.run(
                [sys.executable, str(script_path)],
```

Why this is a false positive: script_path is ROOT / script where script comes from the static JOB_REGISTRY map, plus an explicit path-existence check; the fixed vector [sys.executable, path] cannot become an unintended option.

Checklist evidence: Subprocess argv from static registry with existence check — verified against the shown source.

### [ ] Finding 974 — CWE-909

- Function context: `./scripts/findings/functions/974.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/routers/ops.py:379:53`
- Checklist pattern: `DB resource provided by dependency injection`

Source excerpt:

```
def pipeline_quality(db: Session = Depends(get_db)):
```

Why this is a false positive: db comes from FastAPI dependency injection.

Checklist evidence: DB resource provided by dependency injection — verified against the shown source.

### [ ] Finding 1033 — CWE-89

- Function context: `./scripts/findings/functions/1033.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/routers/search.py:84:18`
- Checklist pattern: `Static SQL with bind parameters reaching execute`

Source excerpt:

```
_sa_text("SELECT entity_type, entity_id, sector, rank FROM entity_search WHERE entity_search MATCH :q ...")
```

Why this is a false positive: Static FTS SELECT with :q bind parameter.

Checklist evidence: Static SQL with bind parameters reaching execute — verified against the shown source.

### [ ] Finding 1059 — CWE-186

- Function context: `./scripts/findings/functions/1059.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/scripts/audit_published_stories.py:96:9`
- Checklist pattern: `Extraction regex, not a validation regex`

Source excerpt:

```
re.compile(r"\b(compared\s+to|versus|vs\.?|prior\s+year|...)", re.IGNORECASE)
```

Why this is a false positive: The regexes are content-detection patterns applied to story bodies, not validation of input.

Checklist evidence: Extraction regex, not a validation regex — verified against the shown source.

### [ ] Finding 1071 — CWE-215

- Function context: `./scripts/findings/functions/1071.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/scripts/diagnose_uspto_odp.py:87:5`
- Checklist pattern: `Sensitive value already masked before logging`

Source excerpt:

```
masked = API_KEY[:6] + "..." + API_KEY[-4:] if len(API_KEY) > 10 else "<short>"
    print(f"  USPTO_API_KEY:  {masked}  (length {len(API_KEY)})")
```

Why this is a false positive: The debug output is already redacted: masked = API_KEY[:6] + "..." + API_KEY[-4:]; the raw key is not printed.

Checklist evidence: Sensitive value already masked before logging — verified against the shown source.

### [ ] Finding 1076 — CWE-909

- Function context: `./scripts/findings/functions/1076.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/scripts/enrich_stories_with_lobbying_issues.py:71:49`
- Checklist pattern: `DB resource injected as a parameter`

Source excerpt:

```
def get_lobbying_data_for_entity(db, entity_id):
```

Why this is a false positive: db is a parameter; the resource is created by the caller.

Checklist evidence: DB resource injected as a parameter — verified against the shown source.

### [ ] Finding 1093 — PERF-PY-28

- Function context: `./scripts/findings/functions/1093.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/scripts/exhaustive_profile_audit.py:291:28`
- Checklist pattern: `Executor created once per run, not per unit of work`

Source excerpt:

```
with ThreadPoolExecutor(max_workers=args.concurrency) as pool:
                futures = {pool.submit(_check_person, pid): pid for pid in people_ids}
```

Why this is a false positive: The ThreadPoolExecutor is created once per main() branch (people audit), not per unit of work; there is no pool-per-item construction.

Checklist evidence: Executor created once per run, not per unit of work — verified against the shown source.

### [ ] Finding 1094 — PERF-PY-28

- Function context: `./scripts/findings/functions/1094.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/scripts/exhaustive_profile_audit.py:301:28`
- Checklist pattern: `Executor created once per run, not per unit of work`

Source excerpt:

```
with ThreadPoolExecutor(max_workers=args.concurrency) as pool:
                futures = {pool.submit(_check_company, sec, slug): (sec, slug)
```

Why this is a false positive: Same: created once per main() branch (companies audit).

Checklist evidence: Executor created once per run, not per unit of work — verified against the shown source.

### [ ] Finding 1096 — CWE-909

- Function context: `./scripts/findings/functions/1096.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/scripts/fix_finance_audit_20260417.py:115:25`
- Checklist pattern: `DB resource injected as a parameter`

Source excerpt:

```
def existing_tables(db):
```

Why this is a false positive: db is a parameter.

Checklist evidence: DB resource injected as a parameter — verified against the shown source.

### [ ] Finding 1097 — CWE-89

- Function context: `./scripts/findings/functions/1097.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/scripts/fix_finance_audit_20260417.py:116:14`
- Checklist pattern: `Static SQL literal reaching execute`

Source excerpt:

```
db.execute(text("SELECT name FROM sqlite_master WHERE type='table'"))
```

Why this is a false positive: Static SELECT against sqlite_master with no interpolation.

Checklist evidence: Static SQL literal reaching execute — verified against the shown source.

### [ ] Finding 1112 — CWE-909

- Function context: `./scripts/findings/functions/1112.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/scripts/generate_lobbying_breakdown_stories.py:163:52`
- Checklist pattern: `DB resource injected as a parameter`

Source excerpt:

```
def generate_sector_breakdown(db, sector_key, cfg):
```

Why this is a false positive: db is a parameter.

Checklist evidence: DB resource injected as a parameter — verified against the shown source.

### [ ] Finding 1160 — PERF-PY-26

- Function context: `./scripts/findings/functions/1160.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/scripts/generate_retraction_patches.py:174:1`
- Checklist pattern: `No expensive decode/parse at the flagged site`

Source excerpt:

```
args = parser.parse_args()
```

Why this is a false positive: Same: parse_args() at CLI entry.

Checklist evidence: No expensive decode/parse at the flagged site — verified against the shown source.

### [ ] Finding 1164 — CWE-909

- Function context: `./scripts/findings/functions/1164.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/scripts/generate_tech_stories.py:68:26`
- Checklist pattern: `DB resource injected as a parameter`

Source excerpt:

```
def generate_stories(db):
```

Why this is a false positive: db is a parameter.

Checklist evidence: DB resource injected as a parameter — verified against the shown source.

### [ ] Finding 1165 — CWE-89

- Function context: `./scripts/findings/functions/1165.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/scripts/generate_tech_stories.py:73:23`
- Checklist pattern: `Static SQL literal reaching execute`

Source excerpt:

```
db.execute(text("""
        SELECT ticker, amount_range, transaction_date, transaction_type
```

Why this is a false positive: The triple-quoted SELECT is a static literal.

Checklist evidence: Static SQL literal reaching execute — verified against the shown source.

### [ ] Finding 1206 — CWE-909

- Function context: `./scripts/findings/functions/1206.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/scripts/remediate_published_stories.py:107:63`
- Checklist pattern: `DB resource injected as a parameter`

Source excerpt:

```
def _derive_range_from_family(db, family, sector, entity_ids):
```

Why this is a false positive: db is a parameter.

Checklist evidence: DB resource injected as a parameter — verified against the shown source.

### [ ] Finding 1279 — CWE-290

- Function context: `./scripts/findings/functions/1279.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/services/auth.py:81:21`
- Checklist pattern: `Header gated behind trusted-proxy network check`

Source excerpt:

```
direct_is_trusted = any(direct_ip in net for net in _TRUSTED_PROXY_NETWORKS)
    if direct_is_trusted:
        forwarded = request.headers.get("x-forwarded-for")
```

Why this is a false positive: X-Forwarded-For is only honored when the immediate connection peer is inside _TRUSTED_PROXY_NETWORKS; the header is not trusted directly.

Checklist evidence: Header gated behind trusted-proxy network check — verified against the shown source.

### [ ] Finding 1359 — BP-PY-13

- Function context: `./scripts/findings/functions/1359.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/services/privacy.py:340:6`
- Checklist pattern: `Sentinel literal, not a real secret`

Source excerpt:

```
user.hashed_password = "ANONYMIZED"
```

Why this is a false positive: The string 'ANONYMIZED' is a sentinel written over a hashed password during account anonymization — not a real credential.

Checklist evidence: Sentinel literal, not a real secret — verified against the shown source.

### [ ] Finding 1377 — CWE-89

- Function context: `./scripts/findings/functions/1377.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/services/research_pipeline/dedup_gate.py:81:14`
- Checklist pattern: `Static SQL with bind parameters reaching execute`

Source excerpt:

```
text("SELECT id, slug, ... WHERE COALESCE(published_at, created_at) >= :cutoff AND status IN ('draft', 'published') AND LOWER(category) = :pattern")
```

Why this is a false positive: Static SELECT with :cutoff/:pattern bind parameters.

Checklist evidence: Static SQL with bind parameters reaching execute — verified against the shown source.

### [ ] Finding 1413 — CWE-89

- Function context: `./scripts/findings/functions/1413.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/services/research_pipeline/rotating_selector.py:71:14`
- Checklist pattern: `Static SQL with bind parameters reaching execute`

Source excerpt:

```
text("SELECT sector, category, entity_ids, published_at, created_at FROM stories WHERE COALESCE(published_at, created_at) >= :cutoff ...")
```

Why this is a false positive: Static SELECT with :cutoff bind parameter.

Checklist evidence: Static SQL with bind parameters reaching execute — verified against the shown source.

### [ ] Finding 1455 — CWE-186

- Function context: `./scripts/findings/functions/1455.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/services/story_validators.py:136:11`
- Checklist pattern: `Extraction regex, not a validation regex`

Source excerpt:

```
YEAR_RE = re.compile(r"\b(20\d{2})\b")
```

Why this is a false positive: YEAR_RE is used with finditer() to detect year claims in story text — extraction, not validation.

Checklist evidence: Extraction regex, not a validation regex — verified against the shown source.

### [ ] Finding 1460 — BP-PY-41

- Function context: `./scripts/findings/functions/1460.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/tests/chaos/test_db_resilience.py:43:1`
- Checklist pattern: `Fixture, not a test function`

Source excerpt:

```
@pytest.fixture(scope="module")
def test_client(working_engine):
```

Why this is a false positive: The flagged symbol is a pytest fixture that returns a TestClient, not a test function; the 'test function without assertions' condition does not apply.

Checklist evidence: Fixture, not a test function — verified against the shown source.

### [ ] Finding 1470 — CWE-89

- Function context: `./scripts/findings/functions/1470.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/tests/performance/query_analysis.py:264:17`
- Checklist pattern: `Static SQL literal reaching execute`

Source excerpt:

```
conn.execute(text("SELECT 1"))
```

Why this is a false positive: SELECT 1 is a static literal in a connectivity check.

Checklist evidence: Static SQL literal reaching execute — verified against the shown source.

### [ ] Finding 1473 — CWE-260

- Function context: `./scripts/findings/functions/1473.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/tests/test_auth.py:8:9`
- Checklist pattern: `Test fixture literal, not a configuration credential`

Source excerpt:

```
r = client.post("/auth/register", json={
        "email": "testuser@example.com",
        "password": "securepassword123",
```

Why this is a false positive: The literal password is a throwaway test-fixture value in a request payload, not a configuration file.

Checklist evidence: Test fixture literal, not a configuration credential — verified against the shown source.

### [ ] Finding 1474 — CWE-89

- Function context: `./scripts/findings/functions/1474.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WeThePeople/utils/db_compat.py:158:19`
- Checklist pattern: `Static SQL literal reaching execute`

Source excerpt:

```
connection.execute(text("PRAGMA journal_mode=WAL"))
```

Why this is a false positive: PRAGMA statements are static literals.

Checklist evidence: Static SQL literal reaching execute — verified against the shown source.

## Uncertain findings

None. Former uncertain CWE-1084 findings (174, 285, 440, 451, 977, 1110, 1174, 1204) reclassified as true positives after confirming detector threshold `len(open/.execute) >= 3` in `detectCWE1084`.


## True positives

### BP-PY-1 — Bare Except Clause (494)

| Finding | Source | Reason |
| --- | --- | --- |
| 3 | `connectors/alpha_vantage.py:75:1` | bare `except Exception`/`BaseException` handler |
| 5 | `connectors/alpha_vantage.py:141:1` | bare `except Exception`/`BaseException` handler |
| 6 | `connectors/cfpb_complaints.py:63:1` | bare `except Exception`/`BaseException` handler |
| 8 | `connectors/clinicaltrials.py:121:1` | bare `except Exception`/`BaseException` handler |
| 10 | `connectors/cms_payments.py:125:1` | bare `except Exception`/`BaseException` handler |
| 12 | `connectors/cms_payments.py:169:1` | bare `except Exception`/`BaseException` handler |
| 13 | `connectors/college_scorecard.py:89:1` | bare `except Exception`/`BaseException` handler |
| 15 | `connectors/college_scorecard.py:127:1` | bare `except Exception`/`BaseException` handler |
| 16 | `connectors/college_scorecard.py:174:1` | bare `except Exception`/`BaseException` handler |
| 17 | `connectors/congress.py:68:1` | bare `except Exception`/`BaseException` handler |
| 19 | `connectors/congress.py:184:1` | bare `except Exception`/`BaseException` handler |
| 24 | `connectors/congress.py:595:1` | bare `except Exception`/`BaseException` handler |
| 25 | `connectors/congress_votes.py:208:1` | bare `except Exception`/`BaseException` handler |
| 27 | `connectors/earmarks.py:168:1` | bare `except Exception`/`BaseException` handler |
| 29 | `connectors/earmarks.py:269:1` | bare `except Exception`/`BaseException` handler |
| 30 | `connectors/epa_envirofacts.py:98:1` | bare `except Exception`/`BaseException` handler |
| 32 | `connectors/everypolitician.py:75:1` | bare `except Exception`/`BaseException` handler |
| 34 | `connectors/everypolitician.py:112:1` | bare `except Exception`/`BaseException` handler |
| 35 | `connectors/fcc_complaints.py:120:1` | bare `except Exception`/`BaseException` handler |
| 37 | `connectors/fcc_complaints.py:159:1` | bare `except Exception`/`BaseException` handler |
| 38 | `connectors/fcc_ecfs.py:63:1` | bare `except Exception`/`BaseException` handler |
| 40 | `connectors/fcc_ecfs.py:100:1` | bare `except Exception`/`BaseException` handler |
| 41 | `connectors/fcc_license.py:92:1` | bare `except Exception`/`BaseException` handler |
| 43 | `connectors/fcc_license.py:149:1` | bare `except Exception`/`BaseException` handler |
| 44 | `connectors/fdic_bankfind.py:85:1` | bare `except Exception`/`BaseException` handler |
| 46 | `connectors/fed_press.py:81:1` | bare `except Exception`/`BaseException` handler |
| 48 | `connectors/federal_register.py:508:1` | bare `except Exception`/`BaseException` handler |
| 51 | `connectors/followthemoney.py:76:1` | bare `except Exception`/`BaseException` handler |
| 53 | `connectors/followthemoney.py:120:1` | bare `except Exception`/`BaseException` handler |
| 54 | `connectors/followthemoney.py:167:1` | bare `except Exception`/`BaseException` handler |
| 55 | `connectors/followthemoney.py:216:1` | bare `except Exception`/`BaseException` handler |
| 56 | `connectors/fred.py:94:1` | bare `except Exception`/`BaseException` handler |
| 61 | `connectors/fred.py:115:1` | bare `except Exception`/`BaseException` handler |
| 62 | `connectors/fred.py:177:1` | bare `except Exception`/`BaseException` handler |
| 63 | `connectors/ftc_cases.py:64:1` | bare `except Exception`/`BaseException` handler |
| 65 | `connectors/fueleconomy.py:86:1` | bare `except Exception`/`BaseException` handler |
| 67 | `connectors/fueleconomy.py:123:1` | bare `except Exception`/`BaseException` handler |
| 68 | `connectors/fueleconomy.py:146:1` | bare `except Exception`/`BaseException` handler |
| 69 | `connectors/fueleconomy.py:168:1` | bare `except Exception`/`BaseException` handler |
| 70 | `connectors/fueleconomy.py:216:1` | bare `except Exception`/`BaseException` handler |
| 71 | `connectors/fueleconomy.py:239:1` | bare `except Exception`/`BaseException` handler |
| 72 | `connectors/fueleconomy.py:261:1` | bare `except Exception`/`BaseException` handler |
| 73 | `connectors/google_civic.py:81:1` | bare `except Exception`/`BaseException` handler |
| 78 | `connectors/govinfo.py:92:1` | bare `except Exception`/`BaseException` handler |
| 80 | `connectors/govinfo.py:380:1` | bare `except Exception`/`BaseException` handler |
| 81 | `connectors/grants_gov.py:70:1` | bare `except Exception`/`BaseException` handler |
| 83 | `connectors/grants_gov.py:112:1` | bare `except Exception`/`BaseException` handler |
| 84 | `connectors/gsa_site_scanning.py:124:1` | bare `except Exception`/`BaseException` handler |
| 86 | `connectors/it_dashboard.py:68:1` | bare `except Exception`/`BaseException` handler |
| 90 | `connectors/news_feed.py:76:1` | bare `except Exception`/`BaseException` handler |
| 92 | `connectors/nhtsa.py:64:1` | bare `except Exception`/`BaseException` handler |
| 94 | `connectors/nhtsa.py:111:1` | bare `except Exception`/`BaseException` handler |
| 95 | `connectors/nhtsa.py:179:1` | bare `except Exception`/`BaseException` handler |
| 96 | `connectors/nhtsa.py:237:1` | bare `except Exception`/`BaseException` handler |
| 97 | `connectors/nhtsa.py:255:1` | bare `except Exception`/`BaseException` handler |
| 98 | `connectors/opencorporates.py:60:1` | bare `except Exception`/`BaseException` handler |
| 100 | `connectors/opencorporates.py:104:1` | bare `except Exception`/`BaseException` handler |
| 101 | `connectors/opencorporates.py:144:1` | bare `except Exception`/`BaseException` handler |
| 102 | `connectors/opencorporates.py:193:1` | bare `except Exception`/`BaseException` handler |
| 103 | `connectors/openfda.py:123:1` | bare `except Exception`/`BaseException` handler |
| 105 | `connectors/openfda.py:207:1` | bare `except Exception`/`BaseException` handler |
| 106 | `connectors/openfda.py:272:1` | bare `except Exception`/`BaseException` handler |
| 107 | `connectors/openstates.py:274:1` | bare `except Exception`/`BaseException` handler |
| 109 | `connectors/samgov.py:95:1` | bare `except Exception`/`BaseException` handler |
| 111 | `connectors/samgov.py:185:1` | bare `except Exception`/`BaseException` handler |
| 112 | `connectors/sec_edgar.py:63:1` | bare `except Exception`/`BaseException` handler |
| 114 | `connectors/senate_lda.py:96:1` | bare `except Exception`/`BaseException` handler |
| 117 | `connectors/treasury_fiscal.py:51:1` | bare `except Exception`/`BaseException` handler |
| 119 | `connectors/urban_institute.py:51:1` | bare `except Exception`/`BaseException` handler |
| 121 | `connectors/usajobs.py:76:1` | bare `except Exception`/`BaseException` handler |
| 123 | `connectors/usaspending.py:343:1` | bare `except Exception`/`BaseException` handler |
| 125 | `connectors/usaspending.py:486:1` | bare `except Exception`/`BaseException` handler |
| 126 | `connectors/usaspending.py:505:1` | bare `except Exception`/`BaseException` handler |
| 127 | `connectors/usaspending.py:525:1` | bare `except Exception`/`BaseException` handler |
| 128 | `connectors/yahoo_finance.py:87:1` | bare `except Exception`/`BaseException` handler |
| 130 | `connectors/yahoo_finance.py:121:1` | bare `except Exception`/`BaseException` handler |
| 131 | `connectors/yahoo_finance.py:165:1` | bare `except Exception`/`BaseException` handler |
| 140 | `jobs/ai_summarize.py:240:1` | bare `except Exception`/`BaseException` handler |
| 180 | `jobs/backfill_bill_ai_summaries.py:145:1` | bare `except Exception`/`BaseException` handler |
| 183 | `jobs/backfill_claims_from_actions.py:205:1` | bare `except Exception`/`BaseException` handler |
| 186 | `jobs/backfill_company_logos.py:124:1` | bare `except Exception`/`BaseException` handler |
| 188 | `jobs/backfill_company_logos.py:147:1` | bare `except Exception`/`BaseException` handler |
| 189 | `jobs/backfill_company_logos.py:177:1` | bare `except Exception`/`BaseException` handler |
| 190 | `jobs/backfill_company_logos.py:219:1` | bare `except Exception`/`BaseException` handler |
| 192 | `jobs/backfill_company_logos.py:246:1` | bare `except Exception`/`BaseException` handler |
| 199 | `jobs/backfill_logos_wikidata.py:88:1` | bare `except Exception`/`BaseException` handler |
| 201 | `jobs/backfill_logos_wikidata.py:114:1` | bare `except Exception`/`BaseException` handler |
| 202 | `jobs/backfill_logos_wikidata.py:141:1` | bare `except Exception`/`BaseException` handler |
| 204 | `jobs/backfill_logos_wikidata.py:186:1` | bare `except Exception`/`BaseException` handler |
| 209 | `jobs/backfill_logos_wikipedia.py:107:1` | bare `except Exception`/`BaseException` handler |
| 211 | `jobs/backfill_logos_wikipedia.py:115:1` | bare `except Exception`/`BaseException` handler |
| 212 | `jobs/backfill_logos_wikipedia.py:145:1` | bare `except Exception`/`BaseException` handler |
| 214 | `jobs/backfill_logos_wikipedia.py:171:1` | bare `except Exception`/`BaseException` handler |
| 219 | `jobs/backfill_sanctions_global.py:284:1` | bare `except Exception`/`BaseException` handler |
| 226 | `jobs/backfill_sanctions_status.py:224:1` | bare `except Exception`/`BaseException` handler |
| 228 | `jobs/backfill_sanctions_status.py:254:1` | bare `except Exception`/`BaseException` handler |
| 229 | `jobs/backfill_sanctions_status.py:286:1` | bare `except Exception`/`BaseException` handler |
| 238 | `jobs/backfill_stock_fundamentals.py:94:1` | bare `except Exception`/`BaseException` handler |
| 242 | `jobs/backfill_story_simplified.py:110:1` | bare `except Exception`/`BaseException` handler |
| 246 | `jobs/backfill_verification_tier.py:65:1` | bare `except Exception`/`BaseException` handler |
| 261 | `jobs/detect_stories.py:209:1` | bare `except Exception`/`BaseException` handler |
| 263 | `jobs/detect_stories.py:421:1` | bare `except Exception`/`BaseException` handler |
| 264 | `jobs/detect_stories.py:449:1` | bare `except Exception`/`BaseException` handler |
| 265 | `jobs/detect_stories.py:538:1` | bare `except Exception`/`BaseException` handler |
| 267 | `jobs/detect_stories.py:584:1` | bare `except Exception`/`BaseException` handler |
| 269 | `jobs/detect_stories.py:674:1` | bare `except Exception`/`BaseException` handler |
| 271 | `jobs/detect_stories.py:731:1` | bare `except Exception`/`BaseException` handler |
| 275 | `jobs/detect_stories.py:1051:1` | bare `except Exception`/`BaseException` handler |
| 277 | `jobs/detect_stories.py:1152:1` | bare `except Exception`/`BaseException` handler |
| 278 | `jobs/detect_stories.py:1163:1` | bare `except Exception`/`BaseException` handler |
| 280 | `jobs/detect_stories.py:1212:1` | bare `except Exception`/`BaseException` handler |
| 287 | `jobs/detect_stories.py:1349:1` | bare `except Exception`/`BaseException` handler |
| 288 | `jobs/detect_stories.py:1360:1` | bare `except Exception`/`BaseException` handler |
| 290 | `jobs/detect_stories.py:1387:1` | bare `except Exception`/`BaseException` handler |
| 292 | `jobs/detect_stories.py:1407:1` | bare `except Exception`/`BaseException` handler |
| 295 | `jobs/detect_stories.py:1519:1` | bare `except Exception`/`BaseException` handler |
| 297 | `jobs/detect_stories.py:1539:1` | bare `except Exception`/`BaseException` handler |
| 298 | `jobs/detect_stories.py:1551:1` | bare `except Exception`/`BaseException` handler |
| 300 | `jobs/detect_stories.py:1638:1` | bare `except Exception`/`BaseException` handler |
| 301 | `jobs/detect_stories.py:1653:1` | bare `except Exception`/`BaseException` handler |
| 302 | `jobs/detect_stories.py:1663:1` | bare `except Exception`/`BaseException` handler |
| 303 | `jobs/detect_stories.py:1674:1` | bare `except Exception`/`BaseException` handler |
| 304 | `jobs/detect_stories.py:1688:1` | bare `except Exception`/`BaseException` handler |
| 306 | `jobs/detect_stories.py:1779:1` | bare `except Exception`/`BaseException` handler |
| 308 | `jobs/detect_stories.py:1797:1` | bare `except Exception`/`BaseException` handler |
| 311 | `jobs/detect_stories.py:1815:1` | bare `except Exception`/`BaseException` handler |
| 313 | `jobs/detect_stories.py:1883:1` | bare `except Exception`/`BaseException` handler |
| 315 | `jobs/detect_stories.py:1931:1` | bare `except Exception`/`BaseException` handler |
| 317 | `jobs/detect_stories.py:1963:1` | bare `except Exception`/`BaseException` handler |
| 319 | `jobs/detect_stories.py:2000:1` | bare `except Exception`/`BaseException` handler |
| 321 | `jobs/detect_stories.py:2021:1` | bare `except Exception`/`BaseException` handler |
| 322 | `jobs/detect_stories.py:2104:1` | bare `except Exception`/`BaseException` handler |
| 323 | `jobs/detect_stories.py:2173:1` | bare `except Exception`/`BaseException` handler |
| 326 | `jobs/detect_stories.py:2255:1` | bare `except Exception`/`BaseException` handler |
| 328 | `jobs/detect_stories.py:2281:1` | bare `except Exception`/`BaseException` handler |
| 330 | `jobs/detect_stories.py:2304:1` | bare `except Exception`/`BaseException` handler |
| 332 | `jobs/detect_stories.py:2395:1` | bare `except Exception`/`BaseException` handler |
| 334 | `jobs/detect_stories.py:2409:1` | bare `except Exception`/`BaseException` handler |
| 336 | `jobs/detect_stories.py:2427:1` | bare `except Exception`/`BaseException` handler |
| 337 | `jobs/detect_stories.py:2524:1` | bare `except Exception`/`BaseException` handler |
| 338 | `jobs/detect_stories.py:2604:1` | bare `except Exception`/`BaseException` handler |
| 340 | `jobs/detect_stories.py:2643:1` | bare `except Exception`/`BaseException` handler |
| 341 | `jobs/detect_stories.py:2661:1` | bare `except Exception`/`BaseException` handler |
| 342 | `jobs/detect_stories.py:2741:1` | bare `except Exception`/`BaseException` handler |
| 344 | `jobs/detect_stories.py:2759:1` | bare `except Exception`/`BaseException` handler |
| 345 | `jobs/detect_stories.py:2771:1` | bare `except Exception`/`BaseException` handler |
| 346 | `jobs/detect_stories.py:2825:1` | bare `except Exception`/`BaseException` handler |
| 348 | `jobs/detect_stories.py:2868:1` | bare `except Exception`/`BaseException` handler |
| 350 | `jobs/detect_stories.py:2897:1` | bare `except Exception`/`BaseException` handler |
| 352 | `jobs/detect_stories.py:2910:1` | bare `except Exception`/`BaseException` handler |
| 354 | `jobs/detect_stories.py:2949:1` | bare `except Exception`/`BaseException` handler |
| 355 | `jobs/detect_stories.py:2994:1` | bare `except Exception`/`BaseException` handler |
| 356 | `jobs/detect_stories.py:3012:1` | bare `except Exception`/`BaseException` handler |
| 357 | `jobs/detect_stories.py:3062:1` | bare `except Exception`/`BaseException` handler |
| 358 | `jobs/detect_stories.py:3083:1` | bare `except Exception`/`BaseException` handler |
| 359 | `jobs/detect_stories.py:3202:1` | bare `except Exception`/`BaseException` handler |
| 360 | `jobs/detect_stories.py:3209:1` | bare `except Exception`/`BaseException` handler |
| 361 | `jobs/detect_stories.py:3284:1` | bare `except Exception`/`BaseException` handler |
| 362 | `jobs/detect_stories.py:3307:1` | bare `except Exception`/`BaseException` handler |
| 367 | `jobs/detect_story_outcomes.py:206:1` | bare `except Exception`/`BaseException` handler |
| 371 | `jobs/detect_story_outcomes.py:302:1` | bare `except Exception`/`BaseException` handler |
| 372 | `jobs/detect_story_outcomes.py:386:1` | bare `except Exception`/`BaseException` handler |
| 384 | `jobs/evaluate_legislative_claims.py:210:1` | bare `except Exception`/`BaseException` handler |
| 387 | `jobs/fetch_photos.py:76:1` | bare `except Exception`/`BaseException` handler |
| 393 | `jobs/generate_digest.py:465:1` | bare `except Exception`/`BaseException` handler |
| 396 | `jobs/generate_digest.py:533:1` | bare `except Exception`/`BaseException` handler |
| 397 | `jobs/generate_digest.py:571:1` | bare `except Exception`/`BaseException` handler |
| 404 | `jobs/generate_under_standards.py:533:1` | bare `except Exception`/`BaseException` handler |
| 406 | `jobs/generate_under_standards.py:593:1` | bare `except Exception`/`BaseException` handler |
| 413 | `jobs/import_openstates_people.py:146:1` | bare `except Exception`/`BaseException` handler |
| 437 | `jobs/migrate_add_indexes.py:133:1` | bare `except Exception`/`BaseException` handler |
| 454 | `jobs/migrate_add_specific_issues.py:30:1` | bare `except Exception`/`BaseException` handler |
| 459 | `jobs/monitor_pipeline.py:171:1` | bare `except Exception`/`BaseException` handler |
| 464 | `jobs/monitor_pipeline.py:241:1` | bare `except Exception`/`BaseException` handler |
| 477 | `jobs/monitor_pipeline.py:322:1` | bare `except Exception`/`BaseException` handler |
| 480 | `jobs/publish_huggingface_dataset.py:213:1` | bare `except Exception`/`BaseException` handler |
| 486 | `jobs/rebuild_search_index.py:126:1` | bare `except Exception`/`BaseException` handler |
| 488 | `jobs/rebuild_search_index.py:146:1` | bare `except Exception`/`BaseException` handler |
| 491 | `jobs/rebuild_search_index.py:193:1` | bare `except Exception`/`BaseException` handler |
| 497 | `jobs/retract_misattributed_stories.py:119:1` | bare `except Exception`/`BaseException` handler |
| 500 | `jobs/retract_misattributed_stories.py:135:1` | bare `except Exception`/`BaseException` handler |
| 501 | `jobs/retract_misattributed_stories.py:194:1` | bare `except Exception`/`BaseException` handler |
| 503 | `jobs/retry_wayback_snapshots.py:146:1` | bare `except Exception`/`BaseException` handler |
| 505 | `jobs/retry_wayback_snapshots.py:167:1` | bare `except Exception`/`BaseException` handler |
| 506 | `jobs/retry_wayback_snapshots.py:170:1` | bare `except Exception`/`BaseException` handler |
| 512 | `jobs/run_pipeline.py:86:1` | bare `except Exception`/`BaseException` handler |
| 519 | `jobs/scheduler.py:835:1` | bare `except Exception`/`BaseException` handler |
| 521 | `jobs/scheduler.py:870:1` | bare `except Exception`/`BaseException` handler |
| 541 | `jobs/send_alerts.py:351:1` | bare `except Exception`/`BaseException` handler |
| 543 | `jobs/send_alerts.py:374:1` | bare `except Exception`/`BaseException` handler |
| 544 | `jobs/send_alerts.py:433:1` | bare `except Exception`/`BaseException` handler |
| 546 | `jobs/story_review_digest.py:186:1` | bare `except Exception`/`BaseException` handler |
| 552 | `jobs/sync_agriculture_data.py:173:1` | bare `except Exception`/`BaseException` handler |
| 554 | `jobs/sync_agriculture_data.py:244:1` | bare `except Exception`/`BaseException` handler |
| 556 | `jobs/sync_agriculture_data.py:324:1` | bare `except Exception`/`BaseException` handler |
| 558 | `jobs/sync_agriculture_data.py:430:1` | bare `except Exception`/`BaseException` handler |
| 560 | `jobs/sync_agriculture_enforcement.py:146:1` | bare `except Exception`/`BaseException` handler |
| 566 | `jobs/sync_chemicals_data.py:226:1` | bare `except Exception`/`BaseException` handler |
| 568 | `jobs/sync_chemicals_data.py:297:1` | bare `except Exception`/`BaseException` handler |
| 570 | `jobs/sync_chemicals_data.py:377:1` | bare `except Exception`/`BaseException` handler |
| 572 | `jobs/sync_chemicals_data.py:483:1` | bare `except Exception`/`BaseException` handler |
| 574 | `jobs/sync_chemicals_enforcement.py:146:1` | bare `except Exception`/`BaseException` handler |
| 580 | `jobs/sync_congressional_trades.py:162:1` | bare `except Exception`/`BaseException` handler |
| 584 | `jobs/sync_congressional_trades.py:225:1` | bare `except Exception`/`BaseException` handler |
| 587 | `jobs/sync_defense_data.py:168:1` | bare `except Exception`/`BaseException` handler |
| 589 | `jobs/sync_defense_data.py:255:1` | bare `except Exception`/`BaseException` handler |
| 591 | `jobs/sync_defense_data.py:338:1` | bare `except Exception`/`BaseException` handler |
| 593 | `jobs/sync_defense_data.py:444:1` | bare `except Exception`/`BaseException` handler |
| 596 | `jobs/sync_defense_enforcement.py:184:1` | bare `except Exception`/`BaseException` handler |
| 604 | `jobs/sync_donations.py:475:1` | bare `except Exception`/`BaseException` handler |
| 608 | `jobs/sync_education_data.py:99:1` | bare `except Exception`/`BaseException` handler |
| 610 | `jobs/sync_education_data.py:170:1` | bare `except Exception`/`BaseException` handler |
| 612 | `jobs/sync_education_data.py:250:1` | bare `except Exception`/`BaseException` handler |
| 614 | `jobs/sync_education_data.py:348:1` | bare `except Exception`/`BaseException` handler |
| 616 | `jobs/sync_education_enforcement.py:158:1` | bare `except Exception`/`BaseException` handler |
| 622 | `jobs/sync_emissions.py:195:1` | bare `except Exception`/`BaseException` handler |
| 624 | `jobs/sync_emissions.py:230:1` | bare `except Exception`/`BaseException` handler |
| 627 | `jobs/sync_energy_data.py:172:1` | bare `except Exception`/`BaseException` handler |
| 629 | `jobs/sync_energy_data.py:243:1` | bare `except Exception`/`BaseException` handler |
| 631 | `jobs/sync_energy_data.py:323:1` | bare `except Exception`/`BaseException` handler |
| 633 | `jobs/sync_energy_data.py:431:1` | bare `except Exception`/`BaseException` handler |
| 635 | `jobs/sync_energy_enforcement.py:146:1` | bare `except Exception`/`BaseException` handler |
| 642 | `jobs/sync_finance_data.py:387:1` | bare `except Exception`/`BaseException` handler |
| 645 | `jobs/sync_finance_enforcement.py:149:1` | bare `except Exception`/`BaseException` handler |
| 651 | `jobs/sync_finance_political_data.py:121:1` | bare `except Exception`/`BaseException` handler |
| 655 | `jobs/sync_finance_political_data.py:211:1` | bare `except Exception`/`BaseException` handler |
| 656 | `jobs/sync_finance_political_data.py:296:1` | bare `except Exception`/`BaseException` handler |
| 658 | `jobs/sync_fuel_economy.py:90:1` | bare `except Exception`/`BaseException` handler |
| 664 | `jobs/sync_health_data.py:333:1` | bare `except Exception`/`BaseException` handler |
| 667 | `jobs/sync_health_enforcement.py:141:1` | bare `except Exception`/`BaseException` handler |
| 673 | `jobs/sync_health_political_data.py:121:1` | bare `except Exception`/`BaseException` handler |
| 677 | `jobs/sync_health_political_data.py:203:1` | bare `except Exception`/`BaseException` handler |
| 678 | `jobs/sync_health_political_data.py:280:1` | bare `except Exception`/`BaseException` handler |
| 680 | `jobs/sync_insider_trades.py:80:1` | bare `except Exception`/`BaseException` handler |
| 682 | `jobs/sync_insider_trades.py:100:1` | bare `except Exception`/`BaseException` handler |
| 683 | `jobs/sync_insider_trades.py:139:1` | bare `except Exception`/`BaseException` handler |
| 684 | `jobs/sync_insider_trades.py:164:1` | bare `except Exception`/`BaseException` handler |
| 686 | `jobs/sync_insider_trades.py:308:1` | bare `except Exception`/`BaseException` handler |
| 687 | `jobs/sync_insider_trades.py:344:1` | bare `except Exception`/`BaseException` handler |
| 689 | `jobs/sync_it_dashboard.py:126:1` | bare `except Exception`/`BaseException` handler |
| 692 | `jobs/sync_member_actions.py:54:1` | bare `except Exception`/`BaseException` handler |
| 695 | `jobs/sync_nhtsa_data.py:127:1` | bare `except Exception`/`BaseException` handler |
| 697 | `jobs/sync_nhtsa_data.py:174:1` | bare `except Exception`/`BaseException` handler |
| 698 | `jobs/sync_nhtsa_data.py:246:1` | bare `except Exception`/`BaseException` handler |
| 702 | `jobs/sync_regulatory_comments.py:93:1` | bare `except Exception`/`BaseException` handler |
| 705 | `jobs/sync_samgov.py:83:1` | bare `except Exception`/`BaseException` handler |
| 707 | `jobs/sync_samgov.py:123:1` | bare `except Exception`/`BaseException` handler |
| 725 | `jobs/sync_senate_votes.py:561:1` | bare `except Exception`/`BaseException` handler |
| 736 | `jobs/sync_site_scanning.py:72:1` | bare `except Exception`/`BaseException` handler |
| 741 | `jobs/sync_state_data_all.py:70:1` | bare `except Exception`/`BaseException` handler |
| 747 | `jobs/sync_tech_data.py:385:1` | bare `except Exception`/`BaseException` handler |
| 751 | `jobs/sync_telecom_data.py:99:1` | bare `except Exception`/`BaseException` handler |
| 753 | `jobs/sync_telecom_data.py:170:1` | bare `except Exception`/`BaseException` handler |
| 755 | `jobs/sync_telecom_data.py:250:1` | bare `except Exception`/`BaseException` handler |
| 757 | `jobs/sync_telecom_data.py:348:1` | bare `except Exception`/`BaseException` handler |
| 759 | `jobs/sync_telecom_enforcement.py:158:1` | bare `except Exception`/`BaseException` handler |
| 766 | `jobs/sync_trades_from_disclosures.py:331:1` | bare `except Exception`/`BaseException` handler |
| 768 | `jobs/sync_trades_from_disclosures.py:360:1` | bare `except Exception`/`BaseException` handler |
| 774 | `jobs/sync_trades_from_disclosures.py:888:1` | bare `except Exception`/`BaseException` handler |
| 775 | `jobs/sync_trades_from_disclosures.py:898:1` | bare `except Exception`/`BaseException` handler |
| 778 | `jobs/sync_transportation_data.py:126:1` | bare `except Exception`/`BaseException` handler |
| 780 | `jobs/sync_transportation_data.py:196:1` | bare `except Exception`/`BaseException` handler |
| 782 | `jobs/sync_transportation_data.py:313:1` | bare `except Exception`/`BaseException` handler |
| 784 | `jobs/sync_transportation_data.py:416:1` | bare `except Exception`/`BaseException` handler |
| 787 | `jobs/sync_transportation_enforcement.py:188:1` | bare `except Exception`/`BaseException` handler |
| 795 | `jobs/sync_votes.py:248:1` | bare `except Exception`/`BaseException` handler |
| 805 | `jobs/twitter_bot.py:617:1` | bare `except Exception`/`BaseException` handler |
| 818 | `jobs/twitter_monitor.py:175:1` | bare `except Exception`/`BaseException` handler |
| 821 | `jobs/twitter_monitor.py:263:1` | bare `except Exception`/`BaseException` handler |
| 823 | `jobs/twitter_monitor.py:275:1` | bare `except Exception`/`BaseException` handler |
| 825 | `jobs/twitter_monitor.py:290:1` | bare `except Exception`/`BaseException` handler |
| 826 | `jobs/twitter_monitor.py:302:1` | bare `except Exception`/`BaseException` handler |
| 827 | `jobs/twitter_monitor.py:314:1` | bare `except Exception`/`BaseException` handler |
| 828 | `jobs/twitter_monitor.py:582:1` | bare `except Exception`/`BaseException` handler |
| 829 | `jobs/twitter_monitor.py:591:1` | bare `except Exception`/`BaseException` handler |
| 834 | `jobs/twitter_monitor.py:678:1` | bare `except Exception`/`BaseException` handler |
| 835 | `jobs/twitter_monitor.py:713:1` | bare `except Exception`/`BaseException` handler |
| 859 | `jobs/twitter_reply.py:577:1` | bare `except Exception`/`BaseException` handler |
| 863 | `jobs/warm_politician_cache.py:60:1` | bare `except Exception`/`BaseException` handler |
| 869 | `main.py:71:1` | bare `except Exception`/`BaseException` handler |
| 871 | `main.py:75:1` | bare `except Exception`/`BaseException` handler |
| 882 | `models/database.py:85:1` | bare `except Exception`/`BaseException` handler |
| 885 | `models/database.py:97:1` | bare `except Exception`/`BaseException` handler |
| 889 | `routers/auth.py:506:1` | bare `except Exception`/`BaseException` handler |
| 891 | `routers/auth.py:587:1` | bare `except Exception`/`BaseException` handler |
| 894 | `routers/auth.py:1249:1` | bare `except Exception`/`BaseException` handler |
| 900 | `routers/auth.py:1370:1` | bare `except Exception`/`BaseException` handler |
| 904 | `routers/chat.py:276:1` | bare `except Exception`/`BaseException` handler |
| 907 | `routers/civic.py:150:1` | bare `except Exception`/`BaseException` handler |
| 913 | `routers/civic.py:809:1` | bare `except Exception`/`BaseException` handler |
| 914 | `routers/civic.py:827:1` | bare `except Exception`/`BaseException` handler |
| 915 | `routers/civic.py:839:1` | bare `except Exception`/`BaseException` handler |
| 918 | `routers/claims.py:271:1` | bare `except Exception`/`BaseException` handler |
| 919 | `routers/common.py:31:1` | bare `except Exception`/`BaseException` handler |
| 926 | `routers/common.py:69:1` | bare `except Exception`/`BaseException` handler |
| 934 | `routers/influence.py:129:1` | bare `except Exception`/`BaseException` handler |
| 936 | `routers/influence.py:304:1` | bare `except Exception`/`BaseException` handler |
| 938 | `routers/influence.py:431:1` | bare `except Exception`/`BaseException` handler |
| 943 | `routers/influence.py:729:1` | bare `except Exception`/`BaseException` handler |
| 945 | `routers/influence.py:973:1` | bare `except Exception`/`BaseException` handler |
| 947 | `routers/lookup.py:74:1` | bare `except Exception`/`BaseException` handler |
| 949 | `routers/lookup.py:100:1` | bare `except Exception`/`BaseException` handler |
| 956 | `routers/metrics.py:141:1` | bare `except Exception`/`BaseException` handler |
| 960 | `routers/metrics.py:160:1` | bare `except Exception`/`BaseException` handler |
| 963 | `routers/og.py:181:1` | bare `except Exception`/`BaseException` handler |
| 965 | `routers/og.py:198:1` | bare `except Exception`/`BaseException` handler |
| 970 | `routers/ops.py:349:1` | bare `except Exception`/`BaseException` handler |
| 976 | `routers/ops.py:406:1` | bare `except Exception`/`BaseException` handler |
| 978 | `routers/ops.py:513:1` | bare `except Exception`/`BaseException` handler |
| 979 | `routers/ops.py:516:1` | bare `except Exception`/`BaseException` handler |
| 980 | `routers/ops.py:529:1` | bare `except Exception`/`BaseException` handler |
| 981 | `routers/ops.py:545:1` | bare `except Exception`/`BaseException` handler |
| 983 | `routers/ops.py:937:1` | bare `except Exception`/`BaseException` handler |
| 984 | `routers/ops.py:974:1` | bare `except Exception`/`BaseException` handler |
| 985 | `routers/ops.py:1003:1` | bare `except Exception`/`BaseException` handler |
| 986 | `routers/ops.py:1006:1` | bare `except Exception`/`BaseException` handler |
| 988 | `routers/ops.py:1041:1` | bare `except Exception`/`BaseException` handler |
| 989 | `routers/ops.py:1090:1` | bare `except Exception`/`BaseException` handler |
| 991 | `routers/ops.py:1310:1` | bare `except Exception`/`BaseException` handler |
| 992 | `routers/ops.py:1425:1` | bare `except Exception`/`BaseException` handler |
| 993 | `routers/ops.py:1463:1` | bare `except Exception`/`BaseException` handler |
| 994 | `routers/ops.py:1515:1` | bare `except Exception`/`BaseException` handler |
| 997 | `routers/ops.py:2143:1` | bare `except Exception`/`BaseException` handler |
| 998 | `routers/politics.py:105:1` | bare `except Exception`/`BaseException` handler |
| 1000 | `routers/politics.py:296:1` | bare `except Exception`/`BaseException` handler |
| 1001 | `routers/politics_bills.py:58:1` | bare `except Exception`/`BaseException` handler |
| 1006 | `routers/politics_bills.py:65:1` | bare `except Exception`/`BaseException` handler |
| 1009 | `routers/politics_bills.py:278:1` | bare `except Exception`/`BaseException` handler |
| 1011 | `routers/politics_people.py:53:1` | bare `except Exception`/`BaseException` handler |
| 1013 | `routers/politics_people.py:136:1` | bare `except Exception`/`BaseException` handler |
| 1014 | `routers/politics_people.py:146:1` | bare `except Exception`/`BaseException` handler |
| 1016 | `routers/politics_people.py:1304:1` | bare `except Exception`/`BaseException` handler |
| 1020 | `routers/politics_people.py:1354:1` | bare `except Exception`/`BaseException` handler |
| 1023 | `routers/research_tools.py:473:1` | bare `except Exception`/`BaseException` handler |
| 1027 | `routers/research_tools.py:769:1` | bare `except Exception`/`BaseException` handler |
| 1028 | `routers/research_tools.py:844:1` | bare `except Exception`/`BaseException` handler |
| 1029 | `routers/research_tools.py:914:1` | bare `except Exception`/`BaseException` handler |
| 1030 | `routers/research_tools.py:987:1` | bare `except Exception`/`BaseException` handler |
| 1031 | `routers/research_tools.py:1100:1` | bare `except Exception`/`BaseException` handler |
| 1032 | `routers/research_tools.py:1197:1` | bare `except Exception`/`BaseException` handler |
| 1034 | `routers/search.py:93:1` | bare `except Exception`/`BaseException` handler |
| 1036 | `routers/search.py:463:1` | bare `except Exception`/`BaseException` handler |
| 1037 | `routers/stories.py:144:1` | bare `except Exception`/`BaseException` handler |
| 1039 | `routers/stories.py:212:1` | bare `except Exception`/`BaseException` handler |
| 1040 | `routers/stories.py:246:1` | bare `except Exception`/`BaseException` handler |
| 1041 | `routers/stories.py:295:1` | bare `except Exception`/`BaseException` handler |
| 1042 | `routers/stories.py:369:1` | bare `except Exception`/`BaseException` handler |
| 1043 | `routers/stories.py:388:1` | bare `except Exception`/`BaseException` handler |
| 1044 | `routers/stories.py:427:1` | bare `except Exception`/`BaseException` handler |
| 1045 | `routers/stories.py:647:1` | bare `except Exception`/`BaseException` handler |
| 1047 | `routers/tech.py:129:1` | bare `except Exception`/`BaseException` handler |
| 1049 | `routers/tips.py:149:1` | bare `except Exception`/`BaseException` handler |
| 1051 | `routers/tips.py:216:1` | bare `except Exception`/`BaseException` handler |
| 1053 | `scripts/apply_retraction_patches.py:143:1` | bare `except Exception`/`BaseException` handler |
| 1055 | `scripts/apply_retraction_patches.py:148:1` | bare `except Exception`/`BaseException` handler |
| 1063 | `scripts/audit_published_stories.py:413:1` | bare `except Exception`/`BaseException` handler |
| 1067 | `scripts/diagnose_usajobs_auth.py:48:1` | bare `except Exception`/`BaseException` handler |
| 1069 | `scripts/diagnose_usajobs_auth.py:57:1` | bare `except Exception`/`BaseException` handler |
| 1079 | `scripts/enrich_stories_with_lobbying_issues.py:93:1` | bare `except Exception`/`BaseException` handler |
| 1082 | `scripts/enrich_stories_with_lobbying_issues.py:146:1` | bare `except Exception`/`BaseException` handler |
| 1085 | `scripts/exhaustive_profile_audit.py:102:1` | bare `except Exception`/`BaseException` handler |
| 1087 | `scripts/exhaustive_profile_audit.py:146:1` | bare `except Exception`/`BaseException` handler |
| 1090 | `scripts/exhaustive_profile_audit.py:176:1` | bare `except Exception`/`BaseException` handler |
| 1117 | `scripts/generate_lobbying_breakdown_stories.py:220:1` | bare `except Exception`/`BaseException` handler |
| 1123 | `scripts/generate_lobbying_breakdown_stories.py:238:1` | bare `except Exception`/`BaseException` handler |
| 1126 | `scripts/generate_lobbying_breakdown_stories.py:394:1` | bare `except Exception`/`BaseException` handler |
| 1129 | `scripts/generate_lobbying_breakdown_stories.py:489:1` | bare `except Exception`/`BaseException` handler |
| 1131 | `scripts/generate_lobbying_breakdown_stories.py:503:1` | bare `except Exception`/`BaseException` handler |
| 1134 | `scripts/generate_lobbying_breakdown_stories.py:517:1` | bare `except Exception`/`BaseException` handler |
| 1137 | `scripts/generate_lobbying_breakdown_stories.py:643:1` | bare `except Exception`/`BaseException` handler |
| 1139 | `scripts/generate_lobbying_breakdown_stories.py:656:1` | bare `except Exception`/`BaseException` handler |
| 1142 | `scripts/generate_lobbying_breakdown_stories.py:670:1` | bare `except Exception`/`BaseException` handler |
| 1145 | `scripts/generate_lobbying_breakdown_stories.py:782:1` | bare `except Exception`/`BaseException` handler |
| 1147 | `scripts/generate_lobbying_breakdown_stories.py:793:1` | bare `except Exception`/`BaseException` handler |
| 1150 | `scripts/generate_lobbying_breakdown_stories.py:807:1` | bare `except Exception`/`BaseException` handler |
| 1153 | `scripts/generate_lobbying_breakdown_stories.py:901:1` | bare `except Exception`/`BaseException` handler |
| 1155 | `scripts/generate_lobbying_breakdown_stories.py:912:1` | bare `except Exception`/`BaseException` handler |
| 1158 | `scripts/generate_lobbying_breakdown_stories.py:926:1` | bare `except Exception`/`BaseException` handler |
| 1169 | `scripts/generate_tech_stories.py:525:1` | bare `except Exception`/`BaseException` handler |
| 1172 | `scripts/migrate_twitter_models.py:29:1` | bare `except Exception`/`BaseException` handler |
| 1180 | `scripts/migrate_twitter_models.py:70:1` | bare `except Exception`/`BaseException` handler |
| 1183 | `scripts/migrate_twitter_models.py:79:1` | bare `except Exception`/`BaseException` handler |
| 1189 | `scripts/migrate_twitter_models.py:105:1` | bare `except Exception`/`BaseException` handler |
| 1201 | `scripts/regenerate_stories_under_new_standards.py:62:1` | bare `except Exception`/`BaseException` handler |
| 1209 | `scripts/remediate_published_stories.py:150:1` | bare `except Exception`/`BaseException` handler |
| 1211 | `scripts/remediate_published_stories.py:163:1` | bare `except Exception`/`BaseException` handler |
| 1212 | `scripts/remediate_published_stories.py:176:1` | bare `except Exception`/`BaseException` handler |
| 1213 | `scripts/remediate_published_stories.py:190:1` | bare `except Exception`/`BaseException` handler |
| 1216 | `scripts/remediate_short_stories.py:154:1` | bare `except Exception`/`BaseException` handler |
| 1221 | `scripts/remediate_short_stories.py:179:1` | bare `except Exception`/`BaseException` handler |
| 1226 | `scripts/remediate_short_stories.py:319:1` | bare `except Exception`/`BaseException` handler |
| 1232 | `scripts/run_story_gates_audit.py:87:1` | bare `except Exception`/`BaseException` handler |
| 1234 | `scripts/run_story_gates_audit.py:98:1` | bare `except Exception`/`BaseException` handler |
| 1236 | `scripts/seed_story_actions.py:869:1` | bare `except Exception`/`BaseException` handler |
| 1243 | `scripts/veritas_v4_patch.py:54:1` | bare `except Exception`/`BaseException` handler |
| 1245 | `scripts/veritas_v4_patch.py:76:1` | bare `except Exception`/`BaseException` handler |
| 1254 | `scripts/wtp_database.py:132:1` | bare `except Exception`/`BaseException` handler |
| 1259 | `scripts/wtp_database.py:201:1` | bare `except Exception`/`BaseException` handler |
| 1264 | `scripts/wtp_database.py:224:1` | bare `except Exception`/`BaseException` handler |
| 1267 | `scripts/wtp_database.py:249:1` | bare `except Exception`/`BaseException` handler |
| 1269 | `scripts/wtp_database.py:269:1` | bare `except Exception`/`BaseException` handler |
| 1271 | `scripts/wtp_database.py:291:1` | bare `except Exception`/`BaseException` handler |
| 1273 | `scripts/wtp_database.py:309:1` | bare `except Exception`/`BaseException` handler |
| 1281 | `services/auth.py:216:1` | bare `except Exception`/`BaseException` handler |
| 1284 | `services/auth.py:372:1` | bare `except Exception`/`BaseException` handler |
| 1286 | `services/auth.py:396:1` | bare `except Exception`/`BaseException` handler |
| 1290 | `services/bill_ai_summary.py:183:1` | bare `except Exception`/`BaseException` handler |
| 1292 | `services/bill_ai_summary.py:211:1` | bare `except Exception`/`BaseException` handler |
| 1293 | `services/bill_ai_summary.py:215:1` | bare `except Exception`/`BaseException` handler |
| 1296 | `services/bill_text.py:88:1` | bare `except Exception`/`BaseException` handler |
| 1302 | `services/budget.py:85:1` | bare `except Exception`/`BaseException` handler |
| 1304 | `services/budget.py:125:1` | bare `except Exception`/`BaseException` handler |
| 1306 | `services/claims/veritas_bridge.py:52:1` | bare `except Exception`/`BaseException` handler |
| 1308 | `services/claims/veritas_bridge.py:97:1` | bare `except Exception`/`BaseException` handler |
| 1311 | `services/claims/veritas_bridge.py:208:1` | bare `except Exception`/`BaseException` handler |
| 1313 | `services/claims/veritas_bridge.py:279:1` | bare `except Exception`/`BaseException` handler |
| 1315 | `services/claims/veritas_bridge.py:301:1` | bare `except Exception`/`BaseException` handler |
| 1316 | `services/claims/veritas_bridge.py:321:1` | bare `except Exception`/`BaseException` handler |
| 1317 | `services/claims/veritas_bridge.py:343:1` | bare `except Exception`/`BaseException` handler |
| 1318 | `services/claims/veritas_bridge.py:363:1` | bare `except Exception`/`BaseException` handler |
| 1319 | `services/claims/veritas_bridge.py:399:1` | bare `except Exception`/`BaseException` handler |
| 1320 | `services/claims/veritas_bridge.py:464:1` | bare `except Exception`/`BaseException` handler |
| 1321 | `services/claims/veritas_bridge.py:532:1` | bare `except Exception`/`BaseException` handler |
| 1322 | `services/claims/veritas_bridge.py:724:1` | bare `except Exception`/`BaseException` handler |
| 1323 | `services/closed_loop_detection.py:92:1` | bare `except Exception`/`BaseException` handler |
| 1329 | `services/closed_loop_detection.py:159:1` | bare `except Exception`/`BaseException` handler |
| 1330 | `services/closed_loop_detection.py:176:1` | bare `except Exception`/`BaseException` handler |
| 1332 | `services/closed_loop_detection.py:303:1` | bare `except Exception`/`BaseException` handler |
| 1333 | `services/closed_loop_detection.py:326:1` | bare `except Exception`/`BaseException` handler |
| 1335 | `services/closed_loop_detection.py:332:1` | bare `except Exception`/`BaseException` handler |
| 1340 | `services/data_retention.py:195:1` | bare `except Exception`/`BaseException` handler |
| 1342 | `services/data_retention.py:296:1` | bare `except Exception`/`BaseException` handler |
| 1344 | `services/email.py:45:1` | bare `except Exception`/`BaseException` handler |
| 1348 | `services/llm/client.py:178:1` | bare `except Exception`/`BaseException` handler |
| 1356 | `services/pipeline_reliability.py:230:1` | bare `except Exception`/`BaseException` handler |
| 1360 | `services/privacy.py:415:1` | bare `except Exception`/`BaseException` handler |
| 1362 | `services/rate_limit_store.py:115:1` | bare `except Exception`/`BaseException` handler |
| 1364 | `services/rate_limit_store.py:179:1` | bare `except Exception`/`BaseException` handler |
| 1368 | `services/research_pipeline/black_swan.py:134:1` | bare `except Exception`/`BaseException` handler |
| 1370 | `services/research_pipeline/black_swan.py:174:1` | bare `except Exception`/`BaseException` handler |
| 1371 | `services/research_pipeline/black_swan.py:233:1` | bare `except Exception`/`BaseException` handler |
| 1373 | `services/research_pipeline/black_swan.py:281:1` | bare `except Exception`/`BaseException` handler |
| 1375 | `services/research_pipeline/black_swan.py:329:1` | bare `except Exception`/`BaseException` handler |
| 1376 | `services/research_pipeline/black_swan.py:371:1` | bare `except Exception`/`BaseException` handler |
| 1378 | `services/research_pipeline/implication_review.py:143:1` | bare `except Exception`/`BaseException` handler |
| 1384 | `services/research_pipeline/orchestrator.py:150:1` | bare `except Exception`/`BaseException` handler |
| 1386 | `services/research_pipeline/orchestrator.py:200:1` | bare `except Exception`/`BaseException` handler |
| 1387 | `services/research_pipeline/orchestrator.py:209:1` | bare `except Exception`/`BaseException` handler |
| 1391 | `services/research_pipeline/orchestrator.py:219:1` | bare `except Exception`/`BaseException` handler |
| 1392 | `services/research_pipeline/orchestrator.py:228:1` | bare `except Exception`/`BaseException` handler |
| 1393 | `services/research_pipeline/orchestrator.py:238:1` | bare `except Exception`/`BaseException` handler |
| 1395 | `services/research_pipeline/orchestrator.py:439:1` | bare `except Exception`/`BaseException` handler |
| 1396 | `services/research_pipeline/orchestrator.py:446:1` | bare `except Exception`/`BaseException` handler |
| 1397 | `services/research_pipeline/orchestrator.py:835:1` | bare `except Exception`/`BaseException` handler |
| 1398 | `services/research_pipeline/orchestrator.py:919:1` | bare `except Exception`/`BaseException` handler |
| 1399 | `services/research_pipeline/orchestrator.py:996:1` | bare `except Exception`/`BaseException` handler |
| 1400 | `services/research_pipeline/orchestrator.py:1038:1` | bare `except Exception`/`BaseException` handler |
| 1401 | `services/research_pipeline/orchestrator.py:1112:1` | bare `except Exception`/`BaseException` handler |
| 1402 | `services/research_pipeline/orchestrator.py:1141:1` | bare `except Exception`/`BaseException` handler |
| 1403 | `services/research_pipeline/orchestrator.py:1187:1` | bare `except Exception`/`BaseException` handler |
| 1404 | `services/research_pipeline/orchestrator.py:1219:1` | bare `except Exception`/`BaseException` handler |
| 1405 | `services/research_pipeline/orchestrator.py:1229:1` | bare `except Exception`/`BaseException` handler |
| 1406 | `services/research_pipeline/orchestrator.py:1243:1` | bare `except Exception`/`BaseException` handler |
| 1408 | `services/research_pipeline/orchestrator.py:1252:1` | bare `except Exception`/`BaseException` handler |
| 1411 | `services/research_pipeline/orphan_check.py:85:1` | bare `except Exception`/`BaseException` handler |
| 1414 | `services/sector_queries.py:123:1` | bare `except Exception`/`BaseException` handler |
| 1416 | `services/sector_queries.py:309:1` | bare `except Exception`/`BaseException` handler |
| 1417 | `services/sector_queries.py:320:1` | bare `except Exception`/`BaseException` handler |
| 1418 | `services/sector_queries.py:331:1` | bare `except Exception`/`BaseException` handler |
| 1421 | `services/story_data_gates.py:216:1` | bare `except Exception`/`BaseException` handler |
| 1424 | `services/story_data_gates.py:229:1` | bare `except Exception`/`BaseException` handler |
| 1426 | `services/story_data_gates.py:243:1` | bare `except Exception`/`BaseException` handler |
| 1428 | `services/story_data_gates.py:257:1` | bare `except Exception`/`BaseException` handler |
| 1430 | `services/story_data_gates.py:272:1` | bare `except Exception`/`BaseException` handler |
| 1432 | `services/story_data_gates.py:293:1` | bare `except Exception`/`BaseException` handler |
| 1435 | `services/story_fact_checker.py:149:1` | bare `except Exception`/`BaseException` handler |
| 1437 | `services/story_fact_checker.py:284:1` | bare `except Exception`/`BaseException` handler |
| 1439 | `services/story_fact_checker.py:295:1` | bare `except Exception`/`BaseException` handler |
| 1441 | `services/story_fact_checker.py:326:1` | bare `except Exception`/`BaseException` handler |
| 1443 | `services/story_fact_checker.py:380:1` | bare `except Exception`/`BaseException` handler |
| 1445 | `services/story_fact_checker.py:400:1` | bare `except Exception`/`BaseException` handler |
| 1447 | `services/story_fact_checker.py:418:1` | bare `except Exception`/`BaseException` handler |
| 1448 | `services/story_simplified_summary.py:130:1` | bare `except Exception`/`BaseException` handler |
| 1450 | `services/story_simplified_summary.py:197:1` | bare `except Exception`/`BaseException` handler |
| 1451 | `services/story_simplified_summary.py:201:1` | bare `except Exception`/`BaseException` handler |
| 1456 | `services/wayback_archive.py:61:1` | bare `except Exception`/`BaseException` handler |
| 1471 | `tests/performance/query_analysis.py:265:1` | bare `except Exception`/`BaseException` handler |
| 1472 | `tests/performance/query_analysis.py:294:1` | bare `except Exception`/`BaseException` handler |
| 1479 | `utils/metrics_hooks.py:33:1` | bare `except Exception`/`BaseException` handler |
| 1484 | `utils/twitter_helpers.py:99:1` | bare `except Exception`/`BaseException` handler |
| 1486 | `utils/twitter_helpers.py:117:1` | bare `except Exception`/`BaseException` handler |
| 1490 | `utils/twitter_helpers.py:127:1` | bare `except Exception`/`BaseException` handler |
| 1491 | `utils/twitter_helpers.py:195:1` | bare `except Exception`/`BaseException` handler |
| 1492 | `utils/twitter_helpers.py:230:1` | bare `except Exception`/`BaseException` handler |

### CWE-396 — Declaration of Catch for Generic Exception (172)

| Finding | Source | Reason |
| --- | --- | --- |
| 4 | `connectors/alpha_vantage.py:75:1` | generic Exception handler |
| 7 | `connectors/cfpb_complaints.py:63:1` | generic Exception handler |
| 9 | `connectors/clinicaltrials.py:121:1` | generic Exception handler |
| 11 | `connectors/cms_payments.py:125:1` | generic Exception handler |
| 14 | `connectors/college_scorecard.py:89:1` | generic Exception handler |
| 18 | `connectors/congress.py:68:1` | generic Exception handler |
| 26 | `connectors/congress_votes.py:208:1` | generic Exception handler |
| 28 | `connectors/earmarks.py:168:1` | generic Exception handler |
| 31 | `connectors/epa_envirofacts.py:98:1` | generic Exception handler |
| 33 | `connectors/everypolitician.py:75:1` | generic Exception handler |
| 36 | `connectors/fcc_complaints.py:120:1` | generic Exception handler |
| 39 | `connectors/fcc_ecfs.py:63:1` | generic Exception handler |
| 42 | `connectors/fcc_license.py:92:1` | generic Exception handler |
| 45 | `connectors/fdic_bankfind.py:85:1` | generic Exception handler |
| 47 | `connectors/fed_press.py:81:1` | generic Exception handler |
| 49 | `connectors/federal_register.py:508:1` | generic Exception handler |
| 50 | `connectors/finnhub.py:81:1` | generic Exception handler |
| 52 | `connectors/followthemoney.py:76:1` | generic Exception handler |
| 59 | `connectors/fred.py:94:1` | generic Exception handler |
| 64 | `connectors/ftc_cases.py:64:1` | generic Exception handler |
| 66 | `connectors/fueleconomy.py:86:1` | generic Exception handler |
| 76 | `connectors/google_civic.py:81:1` | generic Exception handler |
| 79 | `connectors/govinfo.py:92:1` | generic Exception handler |
| 82 | `connectors/grants_gov.py:70:1` | generic Exception handler |
| 85 | `connectors/gsa_site_scanning.py:124:1` | generic Exception handler |
| 87 | `connectors/it_dashboard.py:68:1` | generic Exception handler |
| 91 | `connectors/news_feed.py:76:1` | generic Exception handler |
| 93 | `connectors/nhtsa.py:64:1` | generic Exception handler |
| 99 | `connectors/opencorporates.py:60:1` | generic Exception handler |
| 104 | `connectors/openfda.py:123:1` | generic Exception handler |
| 108 | `connectors/openstates.py:274:1` | generic Exception handler |
| 110 | `connectors/samgov.py:95:1` | generic Exception handler |
| 113 | `connectors/sec_edgar.py:63:1` | generic Exception handler |
| 115 | `connectors/senate_lda.py:96:1` | generic Exception handler |
| 118 | `connectors/treasury_fiscal.py:51:1` | generic Exception handler |
| 120 | `connectors/urban_institute.py:51:1` | generic Exception handler |
| 122 | `connectors/usajobs.py:76:1` | generic Exception handler |
| 124 | `connectors/usaspending.py:343:1` | generic Exception handler |
| 129 | `connectors/yahoo_finance.py:87:1` | generic Exception handler |
| 141 | `jobs/ai_summarize.py:240:1` | generic Exception handler |
| 181 | `jobs/backfill_bill_ai_summaries.py:145:1` | generic Exception handler |
| 184 | `jobs/backfill_claims_from_actions.py:205:1` | generic Exception handler |
| 187 | `jobs/backfill_company_logos.py:124:1` | generic Exception handler |
| 200 | `jobs/backfill_logos_wikidata.py:88:1` | generic Exception handler |
| 210 | `jobs/backfill_logos_wikipedia.py:107:1` | generic Exception handler |
| 220 | `jobs/backfill_sanctions_global.py:284:1` | generic Exception handler |
| 227 | `jobs/backfill_sanctions_status.py:224:1` | generic Exception handler |
| 239 | `jobs/backfill_stock_fundamentals.py:94:1` | generic Exception handler |
| 243 | `jobs/backfill_story_simplified.py:110:1` | generic Exception handler |
| 247 | `jobs/backfill_verification_tier.py:65:1` | generic Exception handler |
| 262 | `jobs/detect_stories.py:209:1` | generic Exception handler |
| 368 | `jobs/detect_story_outcomes.py:206:1` | generic Exception handler |
| 385 | `jobs/evaluate_legislative_claims.py:210:1` | generic Exception handler |
| 388 | `jobs/fetch_photos.py:76:1` | generic Exception handler |
| 394 | `jobs/generate_digest.py:465:1` | generic Exception handler |
| 405 | `jobs/generate_under_standards.py:533:1` | generic Exception handler |
| 414 | `jobs/import_openstates_people.py:146:1` | generic Exception handler |
| 438 | `jobs/migrate_add_indexes.py:133:1` | generic Exception handler |
| 455 | `jobs/migrate_add_specific_issues.py:30:1` | generic Exception handler |
| 461 | `jobs/monitor_pipeline.py:171:1` | generic Exception handler |
| 481 | `jobs/publish_huggingface_dataset.py:213:1` | generic Exception handler |
| 487 | `jobs/rebuild_search_index.py:126:1` | generic Exception handler |
| 498 | `jobs/retract_misattributed_stories.py:119:1` | generic Exception handler |
| 504 | `jobs/retry_wayback_snapshots.py:146:1` | generic Exception handler |
| 513 | `jobs/run_pipeline.py:86:1` | generic Exception handler |
| 520 | `jobs/scheduler.py:835:1` | generic Exception handler |
| 542 | `jobs/send_alerts.py:351:1` | generic Exception handler |
| 547 | `jobs/story_review_digest.py:186:1` | generic Exception handler |
| 553 | `jobs/sync_agriculture_data.py:173:1` | generic Exception handler |
| 561 | `jobs/sync_agriculture_enforcement.py:146:1` | generic Exception handler |
| 567 | `jobs/sync_chemicals_data.py:226:1` | generic Exception handler |
| 575 | `jobs/sync_chemicals_enforcement.py:146:1` | generic Exception handler |
| 579 | `jobs/sync_congressional_trades.py:125:1` | generic Exception handler |
| 588 | `jobs/sync_defense_data.py:168:1` | generic Exception handler |
| 597 | `jobs/sync_defense_enforcement.py:184:1` | generic Exception handler |
| 605 | `jobs/sync_donations.py:475:1` | generic Exception handler |
| 609 | `jobs/sync_education_data.py:99:1` | generic Exception handler |
| 617 | `jobs/sync_education_enforcement.py:158:1` | generic Exception handler |
| 623 | `jobs/sync_emissions.py:195:1` | generic Exception handler |
| 628 | `jobs/sync_energy_data.py:172:1` | generic Exception handler |
| 636 | `jobs/sync_energy_enforcement.py:146:1` | generic Exception handler |
| 640 | `jobs/sync_fara_data.py:185:1` | generic Exception handler |
| 643 | `jobs/sync_finance_data.py:387:1` | generic Exception handler |
| 646 | `jobs/sync_finance_enforcement.py:149:1` | generic Exception handler |
| 652 | `jobs/sync_finance_political_data.py:121:1` | generic Exception handler |
| 659 | `jobs/sync_fuel_economy.py:90:1` | generic Exception handler |
| 665 | `jobs/sync_health_data.py:333:1` | generic Exception handler |
| 668 | `jobs/sync_health_enforcement.py:141:1` | generic Exception handler |
| 674 | `jobs/sync_health_political_data.py:121:1` | generic Exception handler |
| 681 | `jobs/sync_insider_trades.py:80:1` | generic Exception handler |
| 690 | `jobs/sync_it_dashboard.py:126:1` | generic Exception handler |
| 693 | `jobs/sync_member_actions.py:54:1` | generic Exception handler |
| 696 | `jobs/sync_nhtsa_data.py:127:1` | generic Exception handler |
| 700 | `jobs/sync_opensanctions.py:127:1` | generic Exception handler |
| 703 | `jobs/sync_regulatory_comments.py:93:1` | generic Exception handler |
| 706 | `jobs/sync_samgov.py:83:1` | generic Exception handler |
| 726 | `jobs/sync_senate_votes.py:561:1` | generic Exception handler |
| 737 | `jobs/sync_site_scanning.py:72:1` | generic Exception handler |
| 742 | `jobs/sync_state_data_all.py:70:1` | generic Exception handler |
| 748 | `jobs/sync_tech_data.py:385:1` | generic Exception handler |
| 752 | `jobs/sync_telecom_data.py:99:1` | generic Exception handler |
| 760 | `jobs/sync_telecom_enforcement.py:158:1` | generic Exception handler |
| 767 | `jobs/sync_trades_from_disclosures.py:331:1` | generic Exception handler |
| 779 | `jobs/sync_transportation_data.py:126:1` | generic Exception handler |
| 788 | `jobs/sync_transportation_enforcement.py:188:1` | generic Exception handler |
| 796 | `jobs/sync_votes.py:248:1` | generic Exception handler |
| 806 | `jobs/twitter_bot.py:617:1` | generic Exception handler |
| 819 | `jobs/twitter_monitor.py:175:1` | generic Exception handler |
| 860 | `jobs/twitter_reply.py:577:1` | generic Exception handler |
| 864 | `jobs/warm_politician_cache.py:60:1` | generic Exception handler |
| 870 | `main.py:71:1` | generic Exception handler |
| 875 | `middleware/tracing.py:76:1` | generic Exception handler |
| 883 | `models/database.py:85:1` | generic Exception handler |
| 890 | `routers/auth.py:506:1` | generic Exception handler |
| 903 | `routers/chat.py:256:1` | generic Exception handler |
| 910 | `routers/civic.py:150:1` | generic Exception handler |
| 917 | `routers/claims.py:101:1` | generic Exception handler |
| 923 | `routers/common.py:31:1` | generic Exception handler |
| 929 | `routers/digest.py:276:1` | generic Exception handler |
| 931 | `routers/events.py:110:1` | generic Exception handler |
| 935 | `routers/influence.py:129:1` | generic Exception handler |
| 948 | `routers/lookup.py:74:1` | generic Exception handler |
| 958 | `routers/metrics.py:141:1` | generic Exception handler |
| 964 | `routers/og.py:181:1` | generic Exception handler |
| 971 | `routers/ops.py:349:1` | generic Exception handler |
| 999 | `routers/politics.py:105:1` | generic Exception handler |
| 1004 | `routers/politics_bills.py:58:1` | generic Exception handler |
| 1012 | `routers/politics_people.py:53:1` | generic Exception handler |
| 1021 | `routers/research_tools.py:100:1` | generic Exception handler |
| 1035 | `routers/search.py:93:1` | generic Exception handler |
| 1038 | `routers/stories.py:144:1` | generic Exception handler |
| 1048 | `routers/tech.py:129:1` | generic Exception handler |
| 1050 | `routers/tips.py:149:1` | generic Exception handler |
| 1054 | `scripts/apply_retraction_patches.py:143:1` | generic Exception handler |
| 1064 | `scripts/audit_published_stories.py:413:1` | generic Exception handler |
| 1068 | `scripts/diagnose_usajobs_auth.py:48:1` | generic Exception handler |
| 1080 | `scripts/enrich_stories_with_lobbying_issues.py:93:1` | generic Exception handler |
| 1086 | `scripts/exhaustive_profile_audit.py:102:1` | generic Exception handler |
| 1120 | `scripts/generate_lobbying_breakdown_stories.py:220:1` | generic Exception handler |
| 1170 | `scripts/generate_tech_stories.py:525:1` | generic Exception handler |
| 1173 | `scripts/migrate_twitter_models.py:29:1` | generic Exception handler |
| 1202 | `scripts/regenerate_stories_under_new_standards.py:62:1` | generic Exception handler |
| 1210 | `scripts/remediate_published_stories.py:150:1` | generic Exception handler |
| 1219 | `scripts/remediate_short_stories.py:154:1` | generic Exception handler |
| 1233 | `scripts/run_story_gates_audit.py:87:1` | generic Exception handler |
| 1237 | `scripts/seed_story_actions.py:869:1` | generic Exception handler |
| 1255 | `scripts/wtp_database.py:132:1` | generic Exception handler |
| 1275 | `services/audit.py:45:1` | generic Exception handler |
| 1283 | `services/auth.py:216:1` | generic Exception handler |
| 1291 | `services/bill_ai_summary.py:183:1` | generic Exception handler |
| 1297 | `services/bill_text.py:88:1` | generic Exception handler |
| 1303 | `services/budget.py:85:1` | generic Exception handler |
| 1305 | `services/circuit_breaker.py:142:1` | generic Exception handler |
| 1307 | `services/claims/veritas_bridge.py:52:1` | generic Exception handler |
| 1326 | `services/closed_loop_detection.py:92:1` | generic Exception handler |
| 1341 | `services/data_retention.py:195:1` | generic Exception handler |
| 1345 | `services/email.py:45:1` | generic Exception handler |
| 1347 | `services/llm/client.py:130:1` | generic Exception handler |
| 1357 | `services/pipeline_reliability.py:230:1` | generic Exception handler |
| 1361 | `services/privacy.py:415:1` | generic Exception handler |
| 1363 | `services/rate_limit_store.py:115:1` | generic Exception handler |
| 1369 | `services/research_pipeline/black_swan.py:134:1` | generic Exception handler |
| 1379 | `services/research_pipeline/implication_review.py:143:1` | generic Exception handler |
| 1385 | `services/research_pipeline/orchestrator.py:150:1` | generic Exception handler |
| 1412 | `services/research_pipeline/orphan_check.py:85:1` | generic Exception handler |
| 1415 | `services/sector_queries.py:123:1` | generic Exception handler |
| 1422 | `services/story_data_gates.py:216:1` | generic Exception handler |
| 1436 | `services/story_fact_checker.py:149:1` | generic Exception handler |
| 1449 | `services/story_simplified_summary.py:130:1` | generic Exception handler |
| 1457 | `services/wayback_archive.py:61:1` | generic Exception handler |
| 1482 | `utils/metrics_hooks.py:33:1` | generic Exception handler |
| 1485 | `utils/twitter_helpers.py:99:1` | generic Exception handler |

### BP-PY-2 — Except Pass (105)

| Finding | Source | Reason |
| --- | --- | --- |
| 21 | `connectors/congress.py:471:1` | except handler body is only `pass` |
| 23 | `connectors/congress.py:534:1` | except handler body is only `pass` |
| 57 | `connectors/fred.py:94:1` | except handler body is only `pass` |
| 74 | `connectors/google_civic.py:81:1` | except handler body is only `pass` |
| 88 | `connectors/news_feed.py:52:1` | except handler body is only `pass` |
| 132 | `jobs/ai_summarize.py:56:1` | except handler body is only `pass` |
| 136 | `jobs/ai_summarize.py:110:1` | except handler body is only `pass` |
| 253 | `jobs/correct_lobby_double_count_stories.py:193:1` | except handler body is only `pass` |
| 281 | `jobs/detect_stories.py:1212:1` | except handler body is only `pass` |
| 293 | `jobs/detect_stories.py:1407:1` | except handler body is only `pass` |
| 299 | `jobs/detect_stories.py:1551:1` | except handler body is only `pass` |
| 324 | `jobs/detect_stories.py:2173:1` | except handler body is only `pass` |
| 390 | `jobs/generate_digest.py:226:1` | except handler body is only `pass` |
| 401 | `jobs/generate_under_standards.py:453:1` | except handler body is only `pass` |
| 457 | `jobs/monitor_pipeline.py:113:1` | except handler body is only `pass` |
| 460 | `jobs/monitor_pipeline.py:171:1` | except handler body is only `pass` |
| 463 | `jobs/monitor_pipeline.py:193:1` | except handler body is only `pass` |
| 478 | `jobs/monitor_pipeline.py:322:1` | except handler body is only `pass` |
| 507 | `jobs/retry_wayback_snapshots.py:170:1` | except handler body is only `pass` |
| 516 | `jobs/scheduler.py:761:1` | except handler body is only `pass` |
| 539 | `jobs/send_alerts.py:89:1` | except handler body is only `pass` |
| 548 | `jobs/story_review_digest.py:212:1` | except handler body is only `pass` |
| 661 | `jobs/sync_health_data.py:53:1` | except handler body is only `pass` |
| 663 | `jobs/sync_health_data.py:59:1` | except handler body is only `pass` |
| 719 | `jobs/sync_senate_votes.py:326:1` | except handler body is only `pass` |
| 744 | `jobs/sync_tech_data.py:61:1` | except handler body is only `pass` |
| 746 | `jobs/sync_tech_data.py:66:1` | except handler body is only `pass` |
| 769 | `jobs/sync_trades_from_disclosures.py:641:1` | except handler body is only `pass` |
| 772 | `jobs/sync_trades_from_disclosures.py:651:1` | except handler body is only `pass` |
| 803 | `jobs/twitter_bot.py:143:1` | except handler body is only `pass` |
| 866 | `main.py:63:1` | except handler body is only `pass` |
| 876 | `middleware/tracing.py:93:1` | except handler body is only `pass` |
| 879 | `models/database.py:83:1` | except handler body is only `pass` |
| 884 | `models/database.py:95:1` | except handler body is only `pass` |
| 892 | `routers/auth.py:681:1` | except handler body is only `pass` |
| 908 | `routers/civic.py:150:1` | except handler body is only `pass` |
| 920 | `routers/common.py:31:1` | except handler body is only `pass` |
| 927 | `routers/common.py:91:1` | except handler body is only `pass` |
| 928 | `routers/common.py:99:1` | except handler body is only `pass` |
| 939 | `routers/influence.py:431:1` | except handler body is only `pass` |
| 944 | `routers/influence.py:729:1` | except handler body is only `pass` |
| 946 | `routers/influence.py:973:1` | except handler body is only `pass` |
| 950 | `routers/lookup.py:152:1` | except handler body is only `pass` |
| 953 | `routers/metrics.py:123:1` | except handler body is only `pass` |
| 957 | `routers/metrics.py:141:1` | except handler body is only `pass` |
| 961 | `routers/metrics.py:160:1` | except handler body is only `pass` |
| 967 | `routers/ops.py:194:1` | except handler body is only `pass` |
| 972 | `routers/ops.py:372:1` | except handler body is only `pass` |
| 982 | `routers/ops.py:545:1` | except handler body is only `pass` |
| 987 | `routers/ops.py:1006:1` | except handler body is only `pass` |
| 1002 | `routers/politics_bills.py:58:1` | except handler body is only `pass` |
| 1007 | `routers/politics_bills.py:65:1` | except handler body is only `pass` |
| 1017 | `routers/politics_people.py:1304:1` | except handler body is only `pass` |
| 1024 | `routers/research_tools.py:708:1` | except handler body is only `pass` |
| 1056 | `scripts/apply_retraction_patches.py:148:1` | except handler body is only `pass` |
| 1060 | `scripts/audit_published_stories.py:344:1` | except handler body is only `pass` |
| 1118 | `scripts/generate_lobbying_breakdown_stories.py:220:1` | except handler body is only `pass` |
| 1124 | `scripts/generate_lobbying_breakdown_stories.py:238:1` | except handler body is only `pass` |
| 1132 | `scripts/generate_lobbying_breakdown_stories.py:503:1` | except handler body is only `pass` |
| 1135 | `scripts/generate_lobbying_breakdown_stories.py:517:1` | except handler body is only `pass` |
| 1140 | `scripts/generate_lobbying_breakdown_stories.py:656:1` | except handler body is only `pass` |
| 1143 | `scripts/generate_lobbying_breakdown_stories.py:670:1` | except handler body is only `pass` |
| 1148 | `scripts/generate_lobbying_breakdown_stories.py:793:1` | except handler body is only `pass` |
| 1151 | `scripts/generate_lobbying_breakdown_stories.py:807:1` | except handler body is only `pass` |
| 1156 | `scripts/generate_lobbying_breakdown_stories.py:912:1` | except handler body is only `pass` |
| 1159 | `scripts/generate_lobbying_breakdown_stories.py:926:1` | except handler body is only `pass` |
| 1198 | `scripts/regenerate_stories_under_new_standards.py:43:1` | except handler body is only `pass` |
| 1217 | `scripts/remediate_short_stories.py:154:1` | except handler body is only `pass` |
| 1244 | `scripts/veritas_v4_patch.py:54:1` | except handler body is only `pass` |
| 1246 | `scripts/veritas_v4_patch.py:76:1` | except handler body is only `pass` |
| 1260 | `scripts/wtp_database.py:201:1` | except handler body is only `pass` |
| 1265 | `scripts/wtp_database.py:224:1` | except handler body is only `pass` |
| 1268 | `scripts/wtp_database.py:249:1` | except handler body is only `pass` |
| 1270 | `scripts/wtp_database.py:269:1` | except handler body is only `pass` |
| 1272 | `scripts/wtp_database.py:291:1` | except handler body is only `pass` |
| 1274 | `scripts/wtp_database.py:309:1` | except handler body is only `pass` |
| 1276 | `services/auth.py:54:1` | except handler body is only `pass` |
| 1280 | `services/auth.py:90:1` | except handler body is only `pass` |
| 1282 | `services/auth.py:216:1` | except handler body is only `pass` |
| 1285 | `services/auth.py:372:1` | except handler body is only `pass` |
| 1287 | `services/auth.py:396:1` | except handler body is only `pass` |
| 1288 | `services/bill_ai_summary.py:126:1` | except handler body is only `pass` |
| 1294 | `services/bill_ai_summary.py:215:1` | except handler body is only `pass` |
| 1299 | `services/budget.py:36:1` | except handler body is only `pass` |
| 1324 | `services/closed_loop_detection.py:92:1` | except handler body is only `pass` |
| 1328 | `services/closed_loop_detection.py:148:1` | except handler body is only `pass` |
| 1331 | `services/closed_loop_detection.py:176:1` | except handler body is only `pass` |
| 1334 | `services/closed_loop_detection.py:326:1` | except handler body is only `pass` |
| 1336 | `services/closed_loop_detection.py:332:1` | except handler body is only `pass` |
| 1349 | `services/llm/client.py:178:1` | except handler body is only `pass` |
| 1380 | `services/research_pipeline/implication_review.py:170:1` | except handler body is only `pass` |
| 1383 | `services/research_pipeline/implication_review.py:177:1` | except handler body is only `pass` |
| 1388 | `services/research_pipeline/orchestrator.py:209:1` | except handler body is only `pass` |
| 1394 | `services/research_pipeline/orchestrator.py:238:1` | except handler body is only `pass` |
| 1407 | `services/research_pipeline/orchestrator.py:1243:1` | except handler body is only `pass` |
| 1452 | `services/story_simplified_summary.py:201:1` | except handler body is only `pass` |
| 1458 | `tests/chaos/test_circuit_breakers.py:359:1` | except handler body is only `pass` |
| 1461 | `tests/chaos/test_db_resilience.py:147:1` | except handler body is only `pass` |
| 1464 | `tests/chaos/test_db_resilience.py:166:1` | except handler body is only `pass` |
| 1466 | `tests/chaos/test_external_api_failure.py:326:1` | except handler body is only `pass` |
| 1468 | `tests/chaos/test_external_api_failure.py:332:1` | except handler body is only `pass` |
| 1475 | `utils/http_client.py:143:1` | except handler body is only `pass` |
| 1477 | `utils/logging.py:190:1` | except handler body is only `pass` |
| 1480 | `utils/metrics_hooks.py:33:1` | except handler body is only `pass` |
| 1487 | `utils/twitter_helpers.py:117:1` | except handler body is only `pass` |

### BP-PY-45 — sys.path Mutation At Runtime (104)

| Finding | Source | Reason |
| --- | --- | --- |
| 1 | `alembic/env.py:19:1` | sys.path mutated at runtime |
| 135 | `jobs/ai_summarize.py:64:1` | sys.path mutated at runtime |
| 173 | `jobs/audit_orphan_lobby_company_ids.py:28:1` | sys.path mutated at runtime |
| 179 | `jobs/backfill_bill_ai_summaries.py:52:1` | sys.path mutated at runtime |
| 182 | `jobs/backfill_claims_from_actions.py:38:1` | sys.path mutated at runtime |
| 185 | `jobs/backfill_company_logos.py:43:1` | sys.path mutated at runtime |
| 198 | `jobs/backfill_logos_wikidata.py:45:1` | sys.path mutated at runtime |
| 207 | `jobs/backfill_logos_wikipedia.py:45:1` | sys.path mutated at runtime |
| 217 | `jobs/backfill_sanctions_global.py:55:1` | sys.path mutated at runtime |
| 223 | `jobs/backfill_sanctions_status.py:52:1` | sys.path mutated at runtime |
| 236 | `jobs/backfill_stock_fundamentals.py:39:1` | sys.path mutated at runtime |
| 241 | `jobs/backfill_story_simplified.py:49:1` | sys.path mutated at runtime |
| 244 | `jobs/backfill_verification_tier.py:28:1` | sys.path mutated at runtime |
| 249 | `jobs/check_corrected_bodies.py:16:1` | sys.path mutated at runtime |
| 250 | `jobs/cleanup_story_html_comments.py:34:1` | sys.path mutated at runtime |
| 251 | `jobs/correct_lobby_double_count_stories.py:55:1` | sys.path mutated at runtime |
| 255 | `jobs/detect_anomalies.py:16:1` | sys.path mutated at runtime |
| 256 | `jobs/detect_stories.py:53:1` | sys.path mutated at runtime |
| 363 | `jobs/detect_story_outcomes.py:54:1` | sys.path mutated at runtime |
| 379 | `jobs/enforce_retention.py:18:1` | sys.path mutated at runtime |
| 380 | `jobs/evaluate_legislative_claims.py:48:1` | sys.path mutated at runtime |
| 389 | `jobs/generate_digest.py:22:1` | sys.path mutated at runtime |
| 407 | `jobs/import_congress_legislators.py:34:1` | sys.path mutated at runtime |
| 411 | `jobs/import_openstates_people.py:37:1` | sys.path mutated at runtime |
| 417 | `jobs/migrate_add_ai_summaries.py:17:1` | sys.path mutated at runtime |
| 434 | `jobs/migrate_add_indexes.py:9:1` | sys.path mutated at runtime |
| 439 | `jobs/migrate_add_sanctions.py:20:1` | sys.path mutated at runtime |
| 450 | `jobs/migrate_add_specific_issues.py:8:1` | sys.path mutated at runtime |
| 456 | `jobs/monitor_pipeline.py:26:1` | sys.path mutated at runtime |
| 483 | `jobs/rebuild_search_index.py:32:1` | sys.path mutated at runtime |
| 492 | `jobs/repair_correction_notice_substitution.py:57:1` | sys.path mutated at runtime |
| 494 | `jobs/retract_misattributed_stories.py:23:1` | sys.path mutated at runtime |
| 502 | `jobs/retry_wayback_snapshots.py:40:5` | sys.path mutated at runtime |
| 510 | `jobs/run_pipeline.py:56:1` | sys.path mutated at runtime |
| 514 | `jobs/scheduler.py:46:1` | sys.path mutated at runtime |
| 522 | `jobs/seed_badges.py:8:1` | sys.path mutated at runtime |
| 524 | `jobs/seed_education_companies.py:20:1` | sys.path mutated at runtime |
| 526 | `jobs/seed_promises.py:29:1` | sys.path mutated at runtime |
| 529 | `jobs/seed_telecom_companies.py:20:1` | sys.path mutated at runtime |
| 531 | `jobs/seed_tracked_companies.py:24:1` | sys.path mutated at runtime |
| 538 | `jobs/send_alerts.py:34:1` | sys.path mutated at runtime |
| 545 | `jobs/story_review_digest.py:29:1` | sys.path mutated at runtime |
| 550 | `jobs/sync_agriculture_data.py:31:1` | sys.path mutated at runtime |
| 559 | `jobs/sync_agriculture_enforcement.py:25:1` | sys.path mutated at runtime |
| 564 | `jobs/sync_chemicals_data.py:31:1` | sys.path mutated at runtime |
| 573 | `jobs/sync_chemicals_enforcement.py:25:1` | sys.path mutated at runtime |
| 578 | `jobs/sync_congressional_trades.py:29:1` | sys.path mutated at runtime |
| 585 | `jobs/sync_defense_data.py:31:1` | sys.path mutated at runtime |
| 594 | `jobs/sync_defense_enforcement.py:27:1` | sys.path mutated at runtime |
| 600 | `jobs/sync_donations.py:37:1` | sys.path mutated at runtime |
| 606 | `jobs/sync_education_data.py:31:1` | sys.path mutated at runtime |
| 615 | `jobs/sync_education_enforcement.py:25:1` | sys.path mutated at runtime |
| 620 | `jobs/sync_emissions.py:27:1` | sys.path mutated at runtime |
| 625 | `jobs/sync_energy_data.py:32:1` | sys.path mutated at runtime |
| 634 | `jobs/sync_energy_enforcement.py:25:1` | sys.path mutated at runtime |
| 639 | `jobs/sync_fara_data.py:24:1` | sys.path mutated at runtime |
| 641 | `jobs/sync_finance_data.py:22:1` | sys.path mutated at runtime |
| 644 | `jobs/sync_finance_enforcement.py:25:1` | sys.path mutated at runtime |
| 649 | `jobs/sync_finance_political_data.py:32:1` | sys.path mutated at runtime |
| 657 | `jobs/sync_fuel_economy.py:18:1` | sys.path mutated at runtime |
| 660 | `jobs/sync_health_data.py:20:1` | sys.path mutated at runtime |
| 666 | `jobs/sync_health_enforcement.py:25:1` | sys.path mutated at runtime |
| 671 | `jobs/sync_health_political_data.py:32:1` | sys.path mutated at runtime |
| 679 | `jobs/sync_insider_trades.py:31:1` | sys.path mutated at runtime |
| 688 | `jobs/sync_it_dashboard.py:19:1` | sys.path mutated at runtime |
| 691 | `jobs/sync_member_actions.py:37:1` | sys.path mutated at runtime |
| 694 | `jobs/sync_nhtsa_data.py:20:1` | sys.path mutated at runtime |
| 699 | `jobs/sync_opensanctions.py:16:1` | sys.path mutated at runtime |
| 701 | `jobs/sync_regulatory_comments.py:23:1` | sys.path mutated at runtime |
| 704 | `jobs/sync_samgov.py:26:1` | sys.path mutated at runtime |
| 710 | `jobs/sync_senate_votes.py:39:1` | sys.path mutated at runtime |
| 735 | `jobs/sync_site_scanning.py:20:1` | sys.path mutated at runtime |
| 738 | `jobs/sync_state_data.py:26:1` | sys.path mutated at runtime |
| 740 | `jobs/sync_state_data_all.py:35:1` | sys.path mutated at runtime |
| 743 | `jobs/sync_tech_data.py:19:1` | sys.path mutated at runtime |
| 749 | `jobs/sync_telecom_data.py:31:1` | sys.path mutated at runtime |
| 758 | `jobs/sync_telecom_enforcement.py:25:1` | sys.path mutated at runtime |
| 763 | `jobs/sync_trades_from_disclosures.py:40:1` | sys.path mutated at runtime |
| 776 | `jobs/sync_transportation_data.py:31:1` | sys.path mutated at runtime |
| 785 | `jobs/sync_transportation_enforcement.py:27:1` | sys.path mutated at runtime |
| 791 | `jobs/sync_votes.py:26:1` | sys.path mutated at runtime |
| 802 | `jobs/twitter_bot.py:42:1` | sys.path mutated at runtime |
| 815 | `jobs/twitter_monitor.py:39:1` | sys.path mutated at runtime |
| 841 | `jobs/twitter_reply.py:46:1` | sys.path mutated at runtime |
| 861 | `jobs/warm_closed_loop_cache.py:54:5` | sys.path mutated at runtime |
| 862 | `jobs/warm_politician_cache.py:34:1` | sys.path mutated at runtime |
| 1052 | `scripts/apply_retraction_patches.py:37:1` | sys.path mutated at runtime |
| 1065 | `scripts/backfill_normalize_congress_urls.py:45:1` | sys.path mutated at runtime |
| 1066 | `scripts/diagnose_usajobs_auth.py:30:1` | sys.path mutated at runtime |
| 1070 | `scripts/diagnose_uspto_odp.py:34:1` | sys.path mutated at runtime |
| 1073 | `scripts/dump_published_stories.py:38:1` | sys.path mutated at runtime |
| 1074 | `scripts/enrich_stories_with_lobbying_issues.py:22:1` | sys.path mutated at runtime |
| 1084 | `scripts/exhaustive_profile_audit.py:39:1` | sys.path mutated at runtime |
| 1095 | `scripts/fix_finance_audit_20260417.py:32:1` | sys.path mutated at runtime |
| 1109 | `scripts/generate_lobbying_breakdown_stories.py:23:1` | sys.path mutated at runtime |
| 1161 | `scripts/generate_tech_stories.py:22:1` | sys.path mutated at runtime |
| 1171 | `scripts/migrate_twitter_models.py:14:1` | sys.path mutated at runtime |
| 1197 | `scripts/regenerate_stories_under_new_standards.py:37:1` | sys.path mutated at runtime |
| 1203 | `scripts/remediate_published_stories.py:31:1` | sys.path mutated at runtime |
| 1215 | `scripts/remediate_short_stories.py:19:1` | sys.path mutated at runtime |
| 1227 | `scripts/retract_and_correct_stories.py:33:1` | sys.path mutated at runtime |
| 1229 | `scripts/run_story_gates_audit.py:25:1` | sys.path mutated at runtime |
| 1235 | `scripts/seed_story_actions.py:31:1` | sys.path mutated at runtime |
| 1343 | `services/data_retention.py:326:5` | sys.path mutated at runtime |

### BP-PY-35 — SQLAlchemy text With F-String (104)

| Finding | Source | Reason |
| --- | --- | --- |
| 176 | `jobs/audit_orphan_lobby_company_ids.py:56:17` | sqlalchemy.text() built with f-string/% format |
| 177 | `jobs/audit_orphan_lobby_company_ids.py:59:17` | sqlalchemy.text() built with f-string/% format |
| 178 | `jobs/audit_orphan_lobby_company_ids.py:83:31` | sqlalchemy.text() built with f-string/% format |
| 266 | `jobs/detect_stories.py:580:26` | sqlalchemy.text() built with f-string/% format |
| 268 | `jobs/detect_stories.py:668:27` | sqlalchemy.text() built with f-string/% format |
| 270 | `jobs/detect_stories.py:716:30` | sqlalchemy.text() built with f-string/% format |
| 272 | `jobs/detect_stories.py:749:36` | sqlalchemy.text() built with f-string/% format |
| 273 | `jobs/detect_stories.py:753:26` | sqlalchemy.text() built with f-string/% format |
| 274 | `jobs/detect_stories.py:1046:31` | sqlalchemy.text() built with f-string/% format |
| 276 | `jobs/detect_stories.py:1141:27` | sqlalchemy.text() built with f-string/% format |
| 279 | `jobs/detect_stories.py:1206:29` | sqlalchemy.text() built with f-string/% format |
| 286 | `jobs/detect_stories.py:1343:27` | sqlalchemy.text() built with f-string/% format |
| 289 | `jobs/detect_stories.py:1381:38` | sqlalchemy.text() built with f-string/% format |
| 291 | `jobs/detect_stories.py:1400:35` | sqlalchemy.text() built with f-string/% format |
| 294 | `jobs/detect_stories.py:1510:27` | sqlalchemy.text() built with f-string/% format |
| 296 | `jobs/detect_stories.py:1533:34` | sqlalchemy.text() built with f-string/% format |
| 305 | `jobs/detect_stories.py:1759:27` | sqlalchemy.text() built with f-string/% format |
| 307 | `jobs/detect_stories.py:1793:34` | sqlalchemy.text() built with f-string/% format |
| 310 | `jobs/detect_stories.py:1810:34` | sqlalchemy.text() built with f-string/% format |
| 312 | `jobs/detect_stories.py:1878:27` | sqlalchemy.text() built with f-string/% format |
| 314 | `jobs/detect_stories.py:1924:25` | sqlalchemy.text() built with f-string/% format |
| 316 | `jobs/detect_stories.py:1958:27` | sqlalchemy.text() built with f-string/% format |
| 318 | `jobs/detect_stories.py:1996:29` | sqlalchemy.text() built with f-string/% format |
| 320 | `jobs/detect_stories.py:2013:25` | sqlalchemy.text() built with f-string/% format |
| 325 | `jobs/detect_stories.py:2250:31` | sqlalchemy.text() built with f-string/% format |
| 327 | `jobs/detect_stories.py:2276:33` | sqlalchemy.text() built with f-string/% format |
| 329 | `jobs/detect_stories.py:2299:41` | sqlalchemy.text() built with f-string/% format |
| 331 | `jobs/detect_stories.py:2386:27` | sqlalchemy.text() built with f-string/% format |
| 333 | `jobs/detect_stories.py:2404:29` | sqlalchemy.text() built with f-string/% format |
| 335 | `jobs/detect_stories.py:2422:35` | sqlalchemy.text() built with f-string/% format |
| 339 | `jobs/detect_stories.py:2636:32` | sqlalchemy.text() built with f-string/% format |
| 343 | `jobs/detect_stories.py:2752:31` | sqlalchemy.text() built with f-string/% format |
| 347 | `jobs/detect_stories.py:2861:31` | sqlalchemy.text() built with f-string/% format |
| 349 | `jobs/detect_stories.py:2891:38` | sqlalchemy.text() built with f-string/% format |
| 351 | `jobs/detect_stories.py:2905:39` | sqlalchemy.text() built with f-string/% format |
| 353 | `jobs/detect_stories.py:2942:33` | sqlalchemy.text() built with f-string/% format |
| 365 | `jobs/detect_story_outcomes.py:193:21` | sqlalchemy.text() built with f-string/% format |
| 366 | `jobs/detect_story_outcomes.py:200:21` | sqlalchemy.text() built with f-string/% format |
| 369 | `jobs/detect_story_outcomes.py:289:17` | sqlalchemy.text() built with f-string/% format |
| 370 | `jobs/detect_story_outcomes.py:296:17` | sqlalchemy.text() built with f-string/% format |
| 499 | `jobs/retract_misattributed_stories.py:132:24` | sqlalchemy.text() built with f-string/% format |
| 817 | `jobs/twitter_monitor.py:140:36` | sqlalchemy.text() built with f-string/% format |
| 820 | `jobs/twitter_monitor.py:258:45` | sqlalchemy.text() built with f-string/% format |
| 822 | `jobs/twitter_monitor.py:268:46` | sqlalchemy.text() built with f-string/% format |
| 824 | `jobs/twitter_monitor.py:284:39` | sqlalchemy.text() built with f-string/% format |
| 887 | `models/government_data_models.py:42:13` | sqlalchemy.text() built with f-string/% format |
| 1078 | `scripts/enrich_stories_with_lobbying_issues.py:87:31` | sqlalchemy.text() built with f-string/% format |
| 1081 | `scripts/enrich_stories_with_lobbying_issues.py:141:30` | sqlalchemy.text() built with f-string/% format |
| 1088 | `scripts/exhaustive_profile_audit.py:171:20` | sqlalchemy.text() built with f-string/% format |
| 1098 | `scripts/fix_finance_audit_20260417.py:133:13` | sqlalchemy.text() built with f-string/% format |
| 1099 | `scripts/fix_finance_audit_20260417.py:139:17` | sqlalchemy.text() built with f-string/% format |
| 1101 | `scripts/fix_finance_audit_20260417.py:181:17` | sqlalchemy.text() built with f-string/% format |
| 1102 | `scripts/fix_finance_audit_20260417.py:190:17` | sqlalchemy.text() built with f-string/% format |
| 1103 | `scripts/fix_finance_audit_20260417.py:201:17` | sqlalchemy.text() built with f-string/% format |
| 1104 | `scripts/fix_finance_audit_20260417.py:208:17` | sqlalchemy.text() built with f-string/% format |
| 1105 | `scripts/fix_finance_audit_20260417.py:267:25` | sqlalchemy.text() built with f-string/% format |
| 1106 | `scripts/fix_finance_audit_20260417.py:346:25` | sqlalchemy.text() built with f-string/% format |
| 1107 | `scripts/fix_finance_audit_20260417.py:350:21` | sqlalchemy.text() built with f-string/% format |
| 1114 | `scripts/generate_lobbying_breakdown_stories.py:171:23` | sqlalchemy.text() built with f-string/% format |
| 1116 | `scripts/generate_lobbying_breakdown_stories.py:212:36` | sqlalchemy.text() built with f-string/% format |
| 1122 | `scripts/generate_lobbying_breakdown_stories.py:233:32` | sqlalchemy.text() built with f-string/% format |
| 1125 | `scripts/generate_lobbying_breakdown_stories.py:390:31` | sqlalchemy.text() built with f-string/% format |
| 1128 | `scripts/generate_lobbying_breakdown_stories.py:484:31` | sqlalchemy.text() built with f-string/% format |
| 1130 | `scripts/generate_lobbying_breakdown_stories.py:498:36` | sqlalchemy.text() built with f-string/% format |
| 1133 | `scripts/generate_lobbying_breakdown_stories.py:509:32` | sqlalchemy.text() built with f-string/% format |
| 1136 | `scripts/generate_lobbying_breakdown_stories.py:637:31` | sqlalchemy.text() built with f-string/% format |
| 1138 | `scripts/generate_lobbying_breakdown_stories.py:651:36` | sqlalchemy.text() built with f-string/% format |
| 1141 | `scripts/generate_lobbying_breakdown_stories.py:662:32` | sqlalchemy.text() built with f-string/% format |
| 1144 | `scripts/generate_lobbying_breakdown_stories.py:776:31` | sqlalchemy.text() built with f-string/% format |
| 1146 | `scripts/generate_lobbying_breakdown_stories.py:788:40` | sqlalchemy.text() built with f-string/% format |
| 1149 | `scripts/generate_lobbying_breakdown_stories.py:799:36` | sqlalchemy.text() built with f-string/% format |
| 1152 | `scripts/generate_lobbying_breakdown_stories.py:895:31` | sqlalchemy.text() built with f-string/% format |
| 1154 | `scripts/generate_lobbying_breakdown_stories.py:907:40` | sqlalchemy.text() built with f-string/% format |
| 1157 | `scripts/generate_lobbying_breakdown_stories.py:918:36` | sqlalchemy.text() built with f-string/% format |
| 1166 | `scripts/generate_tech_stories.py:157:34` | sqlalchemy.text() built with f-string/% format |
| 1167 | `scripts/generate_tech_stories.py:247:26` | sqlalchemy.text() built with f-string/% format |
| 1168 | `scripts/generate_tech_stories.py:430:38` | sqlalchemy.text() built with f-string/% format |
| 1176 | `scripts/migrate_twitter_models.py:57:34` | sqlalchemy.text() built with f-string/% format |
| 1185 | `scripts/migrate_twitter_models.py:91:34` | sqlalchemy.text() built with f-string/% format |
| 1208 | `scripts/remediate_published_stories.py:145:38` | sqlalchemy.text() built with f-string/% format |
| 1225 | `scripts/remediate_short_stories.py:249:27` | sqlalchemy.text() built with f-string/% format |
| 1310 | `services/claims/veritas_bridge.py:203:35` | sqlalchemy.text() built with f-string/% format |
| 1312 | `services/claims/veritas_bridge.py:265:42` | sqlalchemy.text() built with f-string/% format |
| 1314 | `services/claims/veritas_bridge.py:288:36` | sqlalchemy.text() built with f-string/% format |
| 1337 | `services/closed_loop_detection.py:438:15` | sqlalchemy.text() built with f-string/% format |
| 1354 | `services/pipeline_reliability.py:221:25` | sqlalchemy.text() built with f-string/% format |
| 1355 | `services/pipeline_reliability.py:228:20` | sqlalchemy.text() built with f-string/% format |
| 1358 | `services/pipeline_reliability.py:354:41` | sqlalchemy.text() built with f-string/% format |
| 1367 | `services/research_pipeline/black_swan.py:120:17` | sqlalchemy.text() built with f-string/% format |
| 1372 | `services/research_pipeline/black_swan.py:272:17` | sqlalchemy.text() built with f-string/% format |
| 1374 | `services/research_pipeline/black_swan.py:320:17` | sqlalchemy.text() built with f-string/% format |
| 1410 | `services/research_pipeline/orphan_check.py:81:13` | sqlalchemy.text() built with f-string/% format |
| 1420 | `services/story_data_gates.py:214:26` | sqlalchemy.text() built with f-string/% format |
| 1423 | `services/story_data_gates.py:226:13` | sqlalchemy.text() built with f-string/% format |
| 1425 | `services/story_data_gates.py:239:13` | sqlalchemy.text() built with f-string/% format |
| 1427 | `services/story_data_gates.py:253:13` | sqlalchemy.text() built with f-string/% format |
| 1429 | `services/story_data_gates.py:270:27` | sqlalchemy.text() built with f-string/% format |
| 1431 | `services/story_data_gates.py:287:26` | sqlalchemy.text() built with f-string/% format |
| 1434 | `services/story_fact_checker.py:141:21` | sqlalchemy.text() built with f-string/% format |
| 1438 | `services/story_fact_checker.py:291:21` | sqlalchemy.text() built with f-string/% format |
| 1440 | `services/story_fact_checker.py:321:21` | sqlalchemy.text() built with f-string/% format |
| 1442 | `services/story_fact_checker.py:374:13` | sqlalchemy.text() built with f-string/% format |
| 1444 | `services/story_fact_checker.py:395:13` | sqlalchemy.text() built with f-string/% format |
| 1446 | `services/story_fact_checker.py:412:13` | sqlalchemy.text() built with f-string/% format |

### BP-PY-46 — print Debugging In Library Code (102)

| Finding | Source | Reason |
| --- | --- | --- |
| 245 | `jobs/backfill_verification_tier.py:38:5` | print() used outside `if __name__ == "__main__"` |
| 376 | `jobs/dump_public_snapshot.py:117:17` | print() used outside `if __name__ == "__main__"` |
| 377 | `jobs/dump_public_snapshot.py:131:17` | print() used outside `if __name__ == "__main__"` |
| 398 | `jobs/generate_og_images.py:35:5` | print() used outside `if __name__ == "__main__"` |
| 399 | `jobs/generate_og_images.py:169:5` | print() used outside `if __name__ == "__main__"` |
| 408 | `jobs/import_congress_legislators.py:39:5` | print() used outside `if __name__ == "__main__"` |
| 412 | `jobs/import_openstates_people.py:42:5` | print() used outside `if __name__ == "__main__"` |
| 419 | `jobs/migrate_add_ai_summaries.py:54:9` | print() used outside `if __name__ == "__main__"` |
| 423 | `jobs/migrate_add_ai_summaries.py:70:13` | print() used outside `if __name__ == "__main__"` |
| 426 | `jobs/migrate_add_ai_summaries.py:73:13` | print() used outside `if __name__ == "__main__"` |
| 429 | `jobs/migrate_add_ai_summaries.py:80:13` | print() used outside `if __name__ == "__main__"` |
| 432 | `jobs/migrate_add_ai_summaries.py:83:13` | print() used outside `if __name__ == "__main__"` |
| 433 | `jobs/migrate_add_ai_summaries.py:88:5` | print() used outside `if __name__ == "__main__"` |
| 441 | `jobs/migrate_add_sanctions.py:43:9` | print() used outside `if __name__ == "__main__"` |
| 445 | `jobs/migrate_add_sanctions.py:59:17` | print() used outside `if __name__ == "__main__"` |
| 448 | `jobs/migrate_add_sanctions.py:62:17` | print() used outside `if __name__ == "__main__"` |
| 449 | `jobs/migrate_add_sanctions.py:67:5` | print() used outside `if __name__ == "__main__"` |
| 465 | `jobs/monitor_pipeline.py:280:9` | print() used outside `if __name__ == "__main__"` |
| 466 | `jobs/monitor_pipeline.py:283:9` | print() used outside `if __name__ == "__main__"` |
| 467 | `jobs/monitor_pipeline.py:284:9` | print() used outside `if __name__ == "__main__"` |
| 468 | `jobs/monitor_pipeline.py:285:9` | print() used outside `if __name__ == "__main__"` |
| 469 | `jobs/monitor_pipeline.py:286:9` | print() used outside `if __name__ == "__main__"` |
| 470 | `jobs/monitor_pipeline.py:289:13` | print() used outside `if __name__ == "__main__"` |
| 471 | `jobs/monitor_pipeline.py:291:13` | print() used outside `if __name__ == "__main__"` |
| 472 | `jobs/monitor_pipeline.py:293:13` | print() used outside `if __name__ == "__main__"` |
| 473 | `jobs/monitor_pipeline.py:302:13` | print() used outside `if __name__ == "__main__"` |
| 474 | `jobs/monitor_pipeline.py:305:17` | print() used outside `if __name__ == "__main__"` |
| 475 | `jobs/monitor_pipeline.py:306:13` | print() used outside `if __name__ == "__main__"` |
| 476 | `jobs/monitor_pipeline.py:308:9` | print() used outside `if __name__ == "__main__"` |
| 479 | `jobs/publish_huggingface_dataset.py:135:1` | print() used outside `if __name__ == "__main__"` |
| 482 | `jobs/publish_huggingface_dataset.py:214:9` | print() used outside `if __name__ == "__main__"` |
| 523 | `jobs/seed_badges.py:44:9` | print() used outside `if __name__ == "__main__"` |
| 525 | `jobs/seed_education_companies.py:87:5` | print() used outside `if __name__ == "__main__"` |
| 527 | `jobs/seed_promises.py:253:17` | print() used outside `if __name__ == "__main__"` |
| 528 | `jobs/seed_promises.py:287:9` | print() used outside `if __name__ == "__main__"` |
| 530 | `jobs/seed_telecom_companies.py:87:5` | print() used outside `if __name__ == "__main__"` |
| 532 | `jobs/seed_tracked_companies.py:729:5` | print() used outside `if __name__ == "__main__"` |
| 533 | `jobs/seed_tracked_companies.py:747:5` | print() used outside `if __name__ == "__main__"` |
| 534 | `jobs/seed_tracked_companies.py:765:5` | print() used outside `if __name__ == "__main__"` |
| 535 | `jobs/seed_tracked_companies.py:783:5` | print() used outside `if __name__ == "__main__"` |
| 536 | `jobs/seed_tracked_companies.py:801:5` | print() used outside `if __name__ == "__main__"` |
| 537 | `jobs/seed_tracked_companies.py:819:5` | print() used outside `if __name__ == "__main__"` |
| 808 | `jobs/twitter_bot.py:704:13` | print() used outside `if __name__ == "__main__"` |
| 809 | `jobs/twitter_bot.py:706:17` | print() used outside `if __name__ == "__main__"` |
| 810 | `jobs/twitter_bot.py:708:17` | print() used outside `if __name__ == "__main__"` |
| 811 | `jobs/twitter_bot.py:710:13` | print() used outside `if __name__ == "__main__"` |
| 812 | `jobs/twitter_bot.py:711:13` | print() used outside `if __name__ == "__main__"` |
| 813 | `jobs/twitter_bot.py:713:17` | print() used outside `if __name__ == "__main__"` |
| 814 | `jobs/twitter_bot.py:714:9` | print() used outside `if __name__ == "__main__"` |
| 830 | `jobs/twitter_monitor.py:647:17` | print() used outside `if __name__ == "__main__"` |
| 831 | `jobs/twitter_monitor.py:648:17` | print() used outside `if __name__ == "__main__"` |
| 832 | `jobs/twitter_monitor.py:649:17` | print() used outside `if __name__ == "__main__"` |
| 833 | `jobs/twitter_monitor.py:650:17` | print() used outside `if __name__ == "__main__"` |
| 836 | `jobs/twitter_monitor.py:737:9` | print() used outside `if __name__ == "__main__"` |
| 837 | `jobs/twitter_monitor.py:738:9` | print() used outside `if __name__ == "__main__"` |
| 838 | `jobs/twitter_monitor.py:739:9` | print() used outside `if __name__ == "__main__"` |
| 839 | `jobs/twitter_monitor.py:740:9` | print() used outside `if __name__ == "__main__"` |
| 840 | `jobs/twitter_monitor.py:743:9` | print() used outside `if __name__ == "__main__"` |
| 842 | `jobs/twitter_reply.py:395:9` | print() used outside `if __name__ == "__main__"` |
| 843 | `jobs/twitter_reply.py:396:9` | print() used outside `if __name__ == "__main__"` |
| 844 | `jobs/twitter_reply.py:397:9` | print() used outside `if __name__ == "__main__"` |
| 845 | `jobs/twitter_reply.py:398:9` | print() used outside `if __name__ == "__main__"` |
| 846 | `jobs/twitter_reply.py:399:9` | print() used outside `if __name__ == "__main__"` |
| 847 | `jobs/twitter_reply.py:400:9` | print() used outside `if __name__ == "__main__"` |
| 848 | `jobs/twitter_reply.py:401:9` | print() used outside `if __name__ == "__main__"` |
| 849 | `jobs/twitter_reply.py:402:9` | print() used outside `if __name__ == "__main__"` |
| 850 | `jobs/twitter_reply.py:411:9` | print() used outside `if __name__ == "__main__"` |
| 852 | `jobs/twitter_reply.py:552:17` | print() used outside `if __name__ == "__main__"` |
| 853 | `jobs/twitter_reply.py:553:17` | print() used outside `if __name__ == "__main__"` |
| 854 | `jobs/twitter_reply.py:554:17` | print() used outside `if __name__ == "__main__"` |
| 855 | `jobs/twitter_reply.py:555:17` | print() used outside `if __name__ == "__main__"` |
| 856 | `jobs/twitter_reply.py:556:17` | print() used outside `if __name__ == "__main__"` |
| 857 | `jobs/twitter_reply.py:557:17` | print() used outside `if __name__ == "__main__"` |
| 858 | `jobs/twitter_reply.py:558:17` | print() used outside `if __name__ == "__main__"` |
| 1091 | `scripts/exhaustive_profile_audit.py:177:13` | print() used outside `if __name__ == "__main__"` |
| 1092 | `scripts/exhaustive_profile_audit.py:269:5` | print() used outside `if __name__ == "__main__"` |
| 1177 | `scripts/migrate_twitter_models.py:60:21` | print() used outside `if __name__ == "__main__"` |
| 1178 | `scripts/migrate_twitter_models.py:62:21` | print() used outside `if __name__ == "__main__"` |
| 1179 | `scripts/migrate_twitter_models.py:69:17` | print() used outside `if __name__ == "__main__"` |
| 1181 | `scripts/migrate_twitter_models.py:71:17` | print() used outside `if __name__ == "__main__"` |
| 1182 | `scripts/migrate_twitter_models.py:78:17` | print() used outside `if __name__ == "__main__"` |
| 1184 | `scripts/migrate_twitter_models.py:80:17` | print() used outside `if __name__ == "__main__"` |
| 1186 | `scripts/migrate_twitter_models.py:94:21` | print() used outside `if __name__ == "__main__"` |
| 1187 | `scripts/migrate_twitter_models.py:96:21` | print() used outside `if __name__ == "__main__"` |
| 1188 | `scripts/migrate_twitter_models.py:104:17` | print() used outside `if __name__ == "__main__"` |
| 1190 | `scripts/migrate_twitter_models.py:106:17` | print() used outside `if __name__ == "__main__"` |
| 1191 | `scripts/migrate_twitter_models.py:108:13` | print() used outside `if __name__ == "__main__"` |
| 1192 | `scripts/migrate_twitter_models.py:110:13` | print() used outside `if __name__ == "__main__"` |
| 1193 | `scripts/migrate_twitter_models.py:114:13` | print() used outside `if __name__ == "__main__"` |
| 1194 | `scripts/migrate_twitter_models.py:116:13` | print() used outside `if __name__ == "__main__"` |
| 1195 | `scripts/migrate_twitter_models.py:118:13` | print() used outside `if __name__ == "__main__"` |
| 1196 | `scripts/migrate_twitter_models.py:121:9` | print() used outside `if __name__ == "__main__"` |
| 1238 | `scripts/send_welcome.py:13:5` | print() used outside `if __name__ == "__main__"` |
| 1240 | `scripts/send_welcome.py:40:1` | print() used outside `if __name__ == "__main__"` |
| 1241 | `scripts/send_welcome.py:41:1` | print() used outside `if __name__ == "__main__"` |
| 1242 | `scripts/veritas_v4_patch.py:11:5` | print() used outside `if __name__ == "__main__"` |
| 1247 | `scripts/veritas_v4_patch.py:131:1` | print() used outside `if __name__ == "__main__"` |
| 1248 | `scripts/veritas_v5_patch.py:21:5` | print() used outside `if __name__ == "__main__"` |
| 1249 | `scripts/veritas_v5_patch.py:61:1` | print() used outside `if __name__ == "__main__"` |
| 1250 | `scripts/veritas_v7_patch.py:23:5` | print() used outside `if __name__ == "__main__"` |
| 1251 | `scripts/veritas_v7_patch.py:71:1` | print() used outside `if __name__ == "__main__"` |
| 1298 | `services/bill_text.py:89:9` | print() used outside `if __name__ == "__main__"` |

### CWE-390 — Detection of Error Condition Without Action (53)

| Finding | Source | Reason |
| --- | --- | --- |
| 22 | `connectors/congress.py:471:1` | exception caught with no action taken |
| 58 | `connectors/fred.py:94:1` | exception caught with no action taken |
| 75 | `connectors/google_civic.py:81:1` | exception caught with no action taken |
| 89 | `connectors/news_feed.py:52:1` | exception caught with no action taken |
| 133 | `jobs/ai_summarize.py:56:1` | exception caught with no action taken |
| 254 | `jobs/correct_lobby_double_count_stories.py:193:1` | exception caught with no action taken |
| 282 | `jobs/detect_stories.py:1212:1` | exception caught with no action taken |
| 391 | `jobs/generate_digest.py:226:1` | exception caught with no action taken |
| 402 | `jobs/generate_under_standards.py:453:1` | exception caught with no action taken |
| 458 | `jobs/monitor_pipeline.py:113:1` | exception caught with no action taken |
| 508 | `jobs/retry_wayback_snapshots.py:170:1` | exception caught with no action taken |
| 517 | `jobs/scheduler.py:761:1` | exception caught with no action taken |
| 540 | `jobs/send_alerts.py:89:1` | exception caught with no action taken |
| 549 | `jobs/story_review_digest.py:212:1` | exception caught with no action taken |
| 662 | `jobs/sync_health_data.py:53:1` | exception caught with no action taken |
| 720 | `jobs/sync_senate_votes.py:326:1` | exception caught with no action taken |
| 745 | `jobs/sync_tech_data.py:61:1` | exception caught with no action taken |
| 770 | `jobs/sync_trades_from_disclosures.py:641:1` | exception caught with no action taken |
| 804 | `jobs/twitter_bot.py:143:1` | exception caught with no action taken |
| 867 | `main.py:63:1` | exception caught with no action taken |
| 877 | `middleware/tracing.py:93:1` | exception caught with no action taken |
| 880 | `models/database.py:83:1` | exception caught with no action taken |
| 893 | `routers/auth.py:681:1` | exception caught with no action taken |
| 909 | `routers/civic.py:150:1` | exception caught with no action taken |
| 922 | `routers/common.py:31:1` | exception caught with no action taken |
| 940 | `routers/influence.py:431:1` | exception caught with no action taken |
| 951 | `routers/lookup.py:152:1` | exception caught with no action taken |
| 955 | `routers/metrics.py:123:1` | exception caught with no action taken |
| 968 | `routers/ops.py:194:1` | exception caught with no action taken |
| 1003 | `routers/politics_bills.py:58:1` | exception caught with no action taken |
| 1018 | `routers/politics_people.py:1304:1` | exception caught with no action taken |
| 1025 | `routers/research_tools.py:708:1` | exception caught with no action taken |
| 1057 | `scripts/apply_retraction_patches.py:148:1` | exception caught with no action taken |
| 1061 | `scripts/audit_published_stories.py:344:1` | exception caught with no action taken |
| 1119 | `scripts/generate_lobbying_breakdown_stories.py:220:1` | exception caught with no action taken |
| 1200 | `scripts/regenerate_stories_under_new_standards.py:43:1` | exception caught with no action taken |
| 1218 | `scripts/remediate_short_stories.py:154:1` | exception caught with no action taken |
| 1261 | `scripts/wtp_database.py:201:1` | exception caught with no action taken |
| 1277 | `services/auth.py:54:1` | exception caught with no action taken |
| 1289 | `services/bill_ai_summary.py:126:1` | exception caught with no action taken |
| 1300 | `services/budget.py:36:1` | exception caught with no action taken |
| 1325 | `services/closed_loop_detection.py:92:1` | exception caught with no action taken |
| 1350 | `services/llm/client.py:178:1` | exception caught with no action taken |
| 1381 | `services/research_pipeline/implication_review.py:170:1` | exception caught with no action taken |
| 1389 | `services/research_pipeline/orchestrator.py:209:1` | exception caught with no action taken |
| 1453 | `services/story_simplified_summary.py:201:1` | exception caught with no action taken |
| 1459 | `tests/chaos/test_circuit_breakers.py:359:1` | exception caught with no action taken |
| 1462 | `tests/chaos/test_db_resilience.py:147:1` | exception caught with no action taken |
| 1467 | `tests/chaos/test_external_api_failure.py:326:1` | exception caught with no action taken |
| 1476 | `utils/http_client.py:143:1` | exception caught with no action taken |
| 1478 | `utils/logging.py:190:1` | exception caught with no action taken |
| 1481 | `utils/metrics_hooks.py:33:1` | exception caught with no action taken |
| 1488 | `utils/twitter_helpers.py:117:1` | exception caught with no action taken |

### BP-PY-47 — logging With String Format Before Logger (53)

| Finding | Source | Reason |
| --- | --- | --- |
| 137 | `jobs/ai_summarize.py:180:5` | f-string/.format log message instead of lazy %s |
| 139 | `jobs/ai_summarize.py:213:9` | f-string/.format log message instead of lazy %s |
| 142 | `jobs/ai_summarize.py:242:5` | f-string/.format log message instead of lazy %s |
| 146 | `jobs/ai_summarize.py:337:5` | f-string/.format log message instead of lazy %s |
| 147 | `jobs/ai_summarize.py:364:13` | f-string/.format log message instead of lazy %s |
| 148 | `jobs/ai_summarize.py:366:13` | f-string/.format log message instead of lazy %s |
| 149 | `jobs/ai_summarize.py:367:13` | f-string/.format log message instead of lazy %s |
| 150 | `jobs/ai_summarize.py:389:13` | f-string/.format log message instead of lazy %s |
| 151 | `jobs/ai_summarize.py:392:9` | f-string/.format log message instead of lazy %s |
| 152 | `jobs/ai_summarize.py:418:17` | f-string/.format log message instead of lazy %s |
| 153 | `jobs/ai_summarize.py:420:17` | f-string/.format log message instead of lazy %s |
| 154 | `jobs/ai_summarize.py:442:13` | f-string/.format log message instead of lazy %s |
| 155 | `jobs/ai_summarize.py:445:9` | f-string/.format log message instead of lazy %s |
| 156 | `jobs/ai_summarize.py:471:17` | f-string/.format log message instead of lazy %s |
| 157 | `jobs/ai_summarize.py:473:17` | f-string/.format log message instead of lazy %s |
| 158 | `jobs/ai_summarize.py:495:13` | f-string/.format log message instead of lazy %s |
| 159 | `jobs/ai_summarize.py:498:9` | f-string/.format log message instead of lazy %s |
| 160 | `jobs/ai_summarize.py:526:17` | f-string/.format log message instead of lazy %s |
| 161 | `jobs/ai_summarize.py:528:17` | f-string/.format log message instead of lazy %s |
| 162 | `jobs/ai_summarize.py:547:5` | f-string/.format log message instead of lazy %s |
| 163 | `jobs/ai_summarize.py:614:9` | f-string/.format log message instead of lazy %s |
| 164 | `jobs/ai_summarize.py:667:13` | f-string/.format log message instead of lazy %s |
| 165 | `jobs/ai_summarize.py:670:9` | f-string/.format log message instead of lazy %s |
| 169 | `jobs/ai_summarize.py:728:13` | f-string/.format log message instead of lazy %s |
| 170 | `jobs/ai_summarize.py:755:5` | f-string/.format log message instead of lazy %s |
| 171 | `jobs/ai_summarize.py:799:9` | f-string/.format log message instead of lazy %s |
| 172 | `jobs/ai_summarize.py:800:5` | f-string/.format log message instead of lazy %s |
| 711 | `jobs/sync_senate_votes.py:135:5` | f-string/.format log message instead of lazy %s |
| 713 | `jobs/sync_senate_votes.py:193:17` | f-string/.format log message instead of lazy %s |
| 714 | `jobs/sync_senate_votes.py:195:5` | f-string/.format log message instead of lazy %s |
| 715 | `jobs/sync_senate_votes.py:211:17` | f-string/.format log message instead of lazy %s |
| 716 | `jobs/sync_senate_votes.py:213:5` | f-string/.format log message instead of lazy %s |
| 717 | `jobs/sync_senate_votes.py:248:5` | f-string/.format log message instead of lazy %s |
| 718 | `jobs/sync_senate_votes.py:266:9` | f-string/.format log message instead of lazy %s |
| 722 | `jobs/sync_senate_votes.py:333:5` | f-string/.format log message instead of lazy %s |
| 723 | `jobs/sync_senate_votes.py:547:13` | f-string/.format log message instead of lazy %s |
| 724 | `jobs/sync_senate_votes.py:553:13` | f-string/.format log message instead of lazy %s |
| 727 | `jobs/sync_senate_votes.py:563:9` | f-string/.format log message instead of lazy %s |
| 728 | `jobs/sync_senate_votes.py:592:9` | f-string/.format log message instead of lazy %s |
| 729 | `jobs/sync_senate_votes.py:630:13` | f-string/.format log message instead of lazy %s |
| 730 | `jobs/sync_senate_votes.py:645:17` | f-string/.format log message instead of lazy %s |
| 731 | `jobs/sync_senate_votes.py:685:5` | f-string/.format log message instead of lazy %s |
| 732 | `jobs/sync_senate_votes.py:694:9` | f-string/.format log message instead of lazy %s |
| 733 | `jobs/sync_senate_votes.py:702:9` | f-string/.format log message instead of lazy %s |
| 734 | `jobs/sync_senate_votes.py:707:5` | f-string/.format log message instead of lazy %s |
| 792 | `jobs/sync_votes.py:48:9` | f-string/.format log message instead of lazy %s |
| 794 | `jobs/sync_votes.py:244:9` | f-string/.format log message instead of lazy %s |
| 797 | `jobs/sync_votes.py:250:9` | f-string/.format log message instead of lazy %s |
| 798 | `jobs/sync_votes.py:260:5` | f-string/.format log message instead of lazy %s |
| 799 | `jobs/sync_votes.py:262:5` | f-string/.format log message instead of lazy %s |
| 800 | `jobs/sync_votes.py:305:5` | f-string/.format log message instead of lazy %s |
| 801 | `jobs/sync_votes.py:311:9` | f-string/.format log message instead of lazy %s |
| 1346 | `services/llm/client.py:117:23` | f-string/.format log message instead of lazy %s |

### CWE-1121 — Excessive McCabe Cyclomatic Complexity (39)

| Finding | Source | Reason |
| --- | --- | --- |
| 20 | `connectors/congress.py:416:64` | ≥12 control-flow branches in function |
| 257 | `jobs/detect_stories.py:175:38` | ≥12 control-flow branches in function |
| 395 | `jobs/generate_digest.py:470:12` | ≥12 control-flow branches in function |
| 409 | `jobs/import_congress_legislators.py:199:71` | ≥12 control-flow branches in function |
| 415 | `jobs/import_openstates_people.py:174:87` | ≥12 control-flow branches in function |
| 555 | `jobs/sync_agriculture_data.py:297:65` | ≥12 control-flow branches in function |
| 569 | `jobs/sync_chemicals_data.py:350:62` | ≥12 control-flow branches in function |
| 582 | `jobs/sync_congressional_trades.py:170:12` | ≥12 control-flow branches in function |
| 590 | `jobs/sync_defense_data.py:311:61` | ≥12 control-flow branches in function |
| 595 | `jobs/sync_defense_enforcement.py:147:95` | ≥12 control-flow branches in function |
| 602 | `jobs/sync_donations.py:348:12` | ≥12 control-flow branches in function |
| 611 | `jobs/sync_education_data.py:223:63` | ≥12 control-flow branches in function |
| 630 | `jobs/sync_energy_data.py:296:60` | ≥12 control-flow branches in function |
| 650 | `jobs/sync_finance_political_data.py:93:55` | ≥12 control-flow branches in function |
| 672 | `jobs/sync_health_political_data.py:93:54` | ≥12 control-flow branches in function |
| 708 | `jobs/sync_samgov.py:134:12` | ≥12 control-flow branches in function |
| 754 | `jobs/sync_telecom_data.py:223:61` | ≥12 control-flow branches in function |
| 773 | `jobs/sync_trades_from_disclosures.py:739:12` | ≥12 control-flow branches in function |
| 781 | `jobs/sync_transportation_data.py:249:68` | ≥12 control-flow branches in function |
| 786 | `jobs/sync_transportation_enforcement.py:148:95` | ≥12 control-flow branches in function |
| 807 | `jobs/twitter_bot.py:644:70` | ≥12 control-flow branches in function |
| 851 | `jobs/twitter_reply.py:487:51` | ≥12 control-flow branches in function |
| 930 | `routers/digest.py:331:66` | ≥12 control-flow branches in function |
| 932 | `routers/finance.py:316:129` | ≥12 control-flow branches in function |
| 933 | `routers/health.py:47:199` | ≥12 control-flow branches in function |
| 942 | `routers/influence.py:440:52` | ≥12 control-flow branches in function |
| 952 | `routers/lookup.py:161:62` | ≥12 control-flow branches in function |
| 966 | `routers/ops.py:159:23` | ≥12 control-flow branches in function |
| 1010 | `routers/politics_bills.py:360:28` | ≥12 control-flow branches in function |
| 1015 | `routers/politics_people.py:847:39` | ≥12 control-flow branches in function |
| 1046 | `routers/tech.py:44:98` | ≥12 control-flow branches in function |
| 1075 | `scripts/enrich_stories_with_lobbying_issues.py:71:49` | ≥12 control-flow branches in function |
| 1111 | `scripts/generate_lobbying_breakdown_stories.py:163:52` | ≥12 control-flow branches in function |
| 1163 | `scripts/generate_tech_stories.py:68:26` | ≥12 control-flow branches in function |
| 1205 | `scripts/remediate_published_stories.py:107:63` | ≥12 control-flow branches in function |
| 1223 | `scripts/remediate_short_stories.py:228:12` | ≥12 control-flow branches in function |
| 1228 | `scripts/retract_and_correct_stories.py:159:12` | ≥12 control-flow branches in function |
| 1230 | `scripts/run_story_gates_audit.py:46:12` | ≥12 control-flow branches in function |
| 1365 | `services/rbac.py:133:31` | ≥12 control-flow branches in function |

### CWE-1071 — Empty Code Block (37)

| Finding | Source | Reason |
| --- | --- | --- |
| 60 | `connectors/fred.py:94:5` | handler contains only `pass` |
| 77 | `connectors/google_civic.py:81:13` | handler contains only `pass` |
| 134 | `jobs/ai_summarize.py:56:13` | handler contains only `pass` |
| 283 | `jobs/detect_stories.py:1212:9` | handler contains only `pass` |
| 403 | `jobs/generate_under_standards.py:453:13` | handler contains only `pass` |
| 462 | `jobs/monitor_pipeline.py:171:5` | handler contains only `pass` |
| 509 | `jobs/retry_wayback_snapshots.py:170:13` | handler contains only `pass` |
| 518 | `jobs/scheduler.py:761:17` | handler contains only `pass` |
| 721 | `jobs/sync_senate_votes.py:326:5` | handler contains only `pass` |
| 771 | `jobs/sync_trades_from_disclosures.py:641:9` | handler contains only `pass` |
| 868 | `main.py:63:17` | handler contains only `pass` |
| 881 | `models/database.py:83:5` | handler contains only `pass` |
| 911 | `routers/civic.py:150:9` | handler contains only `pass` |
| 921 | `routers/common.py:31:1` | handler contains only `pass` |
| 941 | `routers/influence.py:431:9` | handler contains only `pass` |
| 954 | `routers/metrics.py:123:1` | handler contains only `pass` |
| 973 | `routers/ops.py:372:5` | handler contains only `pass` |
| 1005 | `routers/politics_bills.py:58:9` | handler contains only `pass` |
| 1019 | `routers/politics_people.py:1304:17` | handler contains only `pass` |
| 1026 | `routers/research_tools.py:708:9` | handler contains only `pass` |
| 1058 | `scripts/apply_retraction_patches.py:148:13` | handler contains only `pass` |
| 1062 | `scripts/audit_published_stories.py:344:13` | handler contains only `pass` |
| 1121 | `scripts/generate_lobbying_breakdown_stories.py:220:5` | handler contains only `pass` |
| 1199 | `scripts/regenerate_stories_under_new_standards.py:43:1` | handler contains only `pass` |
| 1220 | `scripts/remediate_short_stories.py:154:9` | handler contains only `pass` |
| 1262 | `scripts/wtp_database.py:201:21` | handler contains only `pass` |
| 1278 | `services/auth.py:54:9` | handler contains only `pass` |
| 1295 | `services/bill_ai_summary.py:215:9` | handler contains only `pass` |
| 1301 | `services/budget.py:36:13` | handler contains only `pass` |
| 1327 | `services/closed_loop_detection.py:92:5` | handler contains only `pass` |
| 1351 | `services/llm/client.py:178:9` | handler contains only `pass` |
| 1382 | `services/research_pipeline/implication_review.py:170:9` | handler contains only `pass` |
| 1390 | `services/research_pipeline/orchestrator.py:209:13` | handler contains only `pass` |
| 1454 | `services/story_simplified_summary.py:201:9` | handler contains only `pass` |
| 1463 | `tests/chaos/test_db_resilience.py:147:13` | handler contains only `pass` |
| 1483 | `utils/metrics_hooks.py:33:9` | handler contains only `pass` |
| 1489 | `utils/twitter_helpers.py:117:5` | handler contains only `pass` |

### CWE-89 — Improper Neutralization of Special Elements used in an SQL Command ('SQL Injection') (34)

| Finding | Source | Reason |
| --- | --- | --- |
| 143 | `jobs/ai_summarize.py:285:16` | dynamic SQL string reaches execute |
| 175 | `jobs/audit_orphan_lobby_company_ids.py:55:23` | dynamic SQL string reaches execute |
| 191 | `jobs/backfill_company_logos.py:245:16` | dynamic SQL string reaches execute |
| 203 | `jobs/backfill_logos_wikidata.py:185:16` | dynamic SQL string reaches execute |
| 213 | `jobs/backfill_logos_wikipedia.py:170:16` | dynamic SQL string reaches execute |
| 218 | `jobs/backfill_sanctions_global.py:283:16` | dynamic SQL string reaches execute |
| 225 | `jobs/backfill_sanctions_status.py:223:16` | dynamic SQL string reaches execute |
| 364 | `jobs/detect_story_outcomes.py:192:28` | dynamic SQL string reaches execute |
| 374 | `jobs/dump_public_snapshot.py:102:21` | dynamic SQL string reaches execute |
| 422 | `jobs/migrate_add_ai_summaries.py:69:19` | dynamic SQL string reaches execute |
| 444 | `jobs/migrate_add_sanctions.py:58:23` | dynamic SQL string reaches execute |
| 453 | `jobs/migrate_add_specific_issues.py:28:17` | dynamic SQL string reaches execute |
| 816 | `jobs/twitter_monitor.py:140:27` | dynamic SQL string reaches execute |
| 886 | `models/government_data_models.py:41:28` | dynamic SQL string reaches execute |
| 888 | `routers/anomalies.py:111:17` | dynamic SQL string reaches execute |
| 901 | `routers/bulk.py:158:20` | dynamic SQL string reaches execute |
| 975 | `routers/ops.py:394:24` | dynamic SQL string reaches execute |
| 1022 | `routers/research_tools.py:432:18` | dynamic SQL string reaches execute |
| 1077 | `scripts/enrich_stories_with_lobbying_issues.py:87:22` | dynamic SQL string reaches execute |
| 1089 | `scripts/exhaustive_profile_audit.py:172:26` | dynamic SQL string reaches execute |
| 1113 | `scripts/generate_lobbying_breakdown_stories.py:171:14` | dynamic SQL string reaches execute |
| 1175 | `scripts/migrate_twitter_models.py:57:25` | dynamic SQL string reaches execute |
| 1207 | `scripts/remediate_published_stories.py:145:29` | dynamic SQL string reaches execute |
| 1224 | `scripts/remediate_short_stories.py:249:18` | dynamic SQL string reaches execute |
| 1253 | `scripts/wtp_database.py:120:27` | dynamic SQL string reaches execute |
| 1309 | `services/claims/veritas_bridge.py:203:26` | dynamic SQL string reaches execute |
| 1338 | `services/closed_loop_detection.py:450:24` | dynamic SQL string reaches execute |
| 1339 | `services/data_retention.py:167:25` | dynamic SQL string reaches execute |
| 1352 | `services/lobby_spend.py:147:13` | dynamic SQL string reaches execute |
| 1353 | `services/pipeline_reliability.py:221:16` | dynamic SQL string reaches execute |
| 1366 | `services/research_pipeline/black_swan.py:119:22` | dynamic SQL string reaches execute |
| 1409 | `services/research_pipeline/orphan_check.py:80:17` | dynamic SQL string reaches execute |
| 1419 | `services/story_data_gates.py:214:17` | dynamic SQL string reaches execute |
| 1433 | `services/story_fact_checker.py:140:25` | dynamic SQL string reaches execute |

### BP-PY-37 — DB-API Cursor Execute With Percent Format (31)

| Finding | Source | Reason |
| --- | --- | --- |
| 144 | `jobs/ai_summarize.py:292:9` | cursor/conn.execute with f-string/% formatted SQL |
| 145 | `jobs/ai_summarize.py:304:13` | cursor/conn.execute with f-string/% formatted SQL |
| 166 | `jobs/ai_summarize.py:685:23` | cursor/conn.execute with f-string/% formatted SQL |
| 167 | `jobs/ai_summarize.py:696:23` | cursor/conn.execute with f-string/% formatted SQL |
| 168 | `jobs/ai_summarize.py:700:23` | cursor/conn.execute with f-string/% formatted SQL |
| 194 | `jobs/backfill_company_logos.py:306:21` | cursor/conn.execute with f-string/% formatted SQL |
| 195 | `jobs/backfill_company_logos.py:314:17` | cursor/conn.execute with f-string/% formatted SQL |
| 196 | `jobs/backfill_company_logos.py:319:17` | cursor/conn.execute with f-string/% formatted SQL |
| 205 | `jobs/backfill_logos_wikidata.py:212:17` | cursor/conn.execute with f-string/% formatted SQL |
| 215 | `jobs/backfill_logos_wikipedia.py:197:17` | cursor/conn.execute with f-string/% formatted SQL |
| 221 | `jobs/backfill_sanctions_global.py:327:16` | cursor/conn.execute with f-string/% formatted SQL |
| 222 | `jobs/backfill_sanctions_global.py:338:16` | cursor/conn.execute with f-string/% formatted SQL |
| 230 | `jobs/backfill_sanctions_status.py:322:25` | cursor/conn.execute with f-string/% formatted SQL |
| 232 | `jobs/backfill_sanctions_status.py:328:25` | cursor/conn.execute with f-string/% formatted SQL |
| 233 | `jobs/backfill_sanctions_status.py:347:25` | cursor/conn.execute with f-string/% formatted SQL |
| 234 | `jobs/backfill_sanctions_status.py:353:25` | cursor/conn.execute with f-string/% formatted SQL |
| 240 | `jobs/backfill_stock_fundamentals.py:130:9` | cursor/conn.execute with f-string/% formatted SQL |
| 373 | `jobs/dump_public_snapshot.py:102:21` | cursor/conn.execute with f-string/% formatted SQL |
| 375 | `jobs/dump_public_snapshot.py:107:21` | cursor/conn.execute with f-string/% formatted SQL |
| 421 | `jobs/migrate_add_ai_summaries.py:69:19` | cursor/conn.execute with f-string/% formatted SQL |
| 425 | `jobs/migrate_add_ai_summaries.py:72:19` | cursor/conn.execute with f-string/% formatted SQL |
| 428 | `jobs/migrate_add_ai_summaries.py:79:19` | cursor/conn.execute with f-string/% formatted SQL |
| 431 | `jobs/migrate_add_ai_summaries.py:82:19` | cursor/conn.execute with f-string/% formatted SQL |
| 443 | `jobs/migrate_add_sanctions.py:58:23` | cursor/conn.execute with f-string/% formatted SQL |
| 447 | `jobs/migrate_add_sanctions.py:61:23` | cursor/conn.execute with f-string/% formatted SQL |
| 452 | `jobs/migrate_add_specific_issues.py:28:17` | cursor/conn.execute with f-string/% formatted SQL |
| 1252 | `scripts/wtp_database.py:120:27` | cursor/conn.execute with f-string/% formatted SQL |
| 1256 | `scripts/wtp_database.py:181:33` | cursor/conn.execute with f-string/% formatted SQL |
| 1258 | `scripts/wtp_database.py:188:46` | cursor/conn.execute with f-string/% formatted SQL |
| 1263 | `scripts/wtp_database.py:209:33` | cursor/conn.execute with f-string/% formatted SQL |
| 1266 | `scripts/wtp_database.py:232:33` | cursor/conn.execute with f-string/% formatted SQL |

### CWE-117 — Improper Output Neutralization for Logs (30)

| Finding | Source | Reason |
| --- | --- | --- |
| 138 | `jobs/ai_summarize.py:180:5` | dynamic value formatted into log message without CRLF neutralization |
| 410 | `jobs/import_congress_legislators.py:236:17` | dynamic value formatted into log message without CRLF neutralization |
| 416 | `jobs/import_openstates_people.py:206:9` | dynamic value formatted into log message without CRLF neutralization |
| 551 | `jobs/sync_agriculture_data.py:153:5` | dynamic value formatted into log message without CRLF neutralization |
| 562 | `jobs/sync_agriculture_enforcement.py:225:5` | dynamic value formatted into log message without CRLF neutralization |
| 565 | `jobs/sync_chemicals_data.py:206:5` | dynamic value formatted into log message without CRLF neutralization |
| 576 | `jobs/sync_chemicals_enforcement.py:225:5` | dynamic value formatted into log message without CRLF neutralization |
| 581 | `jobs/sync_congressional_trades.py:163:13` | dynamic value formatted into log message without CRLF neutralization |
| 586 | `jobs/sync_defense_data.py:149:5` | dynamic value formatted into log message without CRLF neutralization |
| 598 | `jobs/sync_defense_enforcement.py:289:5` | dynamic value formatted into log message without CRLF neutralization |
| 601 | `jobs/sync_donations.py:242:5` | dynamic value formatted into log message without CRLF neutralization |
| 607 | `jobs/sync_education_data.py:87:9` | dynamic value formatted into log message without CRLF neutralization |
| 618 | `jobs/sync_education_enforcement.py:244:5` | dynamic value formatted into log message without CRLF neutralization |
| 621 | `jobs/sync_emissions.py:132:5` | dynamic value formatted into log message without CRLF neutralization |
| 626 | `jobs/sync_energy_data.py:151:5` | dynamic value formatted into log message without CRLF neutralization |
| 637 | `jobs/sync_energy_enforcement.py:225:5` | dynamic value formatted into log message without CRLF neutralization |
| 647 | `jobs/sync_finance_enforcement.py:228:5` | dynamic value formatted into log message without CRLF neutralization |
| 654 | `jobs/sync_finance_political_data.py:178:5` | dynamic value formatted into log message without CRLF neutralization |
| 669 | `jobs/sync_health_enforcement.py:225:5` | dynamic value formatted into log message without CRLF neutralization |
| 676 | `jobs/sync_health_political_data.py:170:5` | dynamic value formatted into log message without CRLF neutralization |
| 712 | `jobs/sync_senate_votes.py:135:5` | dynamic value formatted into log message without CRLF neutralization |
| 739 | `jobs/sync_state_data.py:73:5` | dynamic value formatted into log message without CRLF neutralization |
| 750 | `jobs/sync_telecom_data.py:87:9` | dynamic value formatted into log message without CRLF neutralization |
| 761 | `jobs/sync_telecom_enforcement.py:244:5` | dynamic value formatted into log message without CRLF neutralization |
| 765 | `jobs/sync_trades_from_disclosures.py:302:9` | dynamic value formatted into log message without CRLF neutralization |
| 777 | `jobs/sync_transportation_data.py:114:9` | dynamic value formatted into log message without CRLF neutralization |
| 789 | `jobs/sync_transportation_enforcement.py:297:5` | dynamic value formatted into log message without CRLF neutralization |
| 793 | `jobs/sync_votes.py:48:9` | dynamic value formatted into log message without CRLF neutralization |
| 1083 | `scripts/enrich_stories_with_lobbying_issues.py:245:13` | dynamic value formatted into log message without CRLF neutralization |
| 1115 | `scripts/generate_lobbying_breakdown_stories.py:178:9` | dynamic value formatted into log message without CRLF neutralization |

### CWE-1124 — Excessively Deep Nesting (15)

| Finding | Source | Reason |
| --- | --- | --- |
| 116 | `connectors/senate_lda.py:121:1` | statement nested ≥6 control-flow levels |
| 260 | `jobs/detect_stories.py:206:1` | statement nested ≥6 control-flow levels |
| 557 | `jobs/sync_agriculture_data.py:348:1` | statement nested ≥6 control-flow levels |
| 571 | `jobs/sync_chemicals_data.py:401:1` | statement nested ≥6 control-flow levels |
| 583 | `jobs/sync_congressional_trades.py:219:1` | statement nested ≥6 control-flow levels |
| 592 | `jobs/sync_defense_data.py:362:1` | statement nested ≥6 control-flow levels |
| 603 | `jobs/sync_donations.py:447:1` | statement nested ≥6 control-flow levels |
| 613 | `jobs/sync_education_data.py:274:1` | statement nested ≥6 control-flow levels |
| 632 | `jobs/sync_energy_data.py:347:1` | statement nested ≥6 control-flow levels |
| 653 | `jobs/sync_finance_political_data.py:144:1` | statement nested ≥6 control-flow levels |
| 675 | `jobs/sync_health_political_data.py:144:1` | statement nested ≥6 control-flow levels |
| 709 | `jobs/sync_samgov.py:176:1` | statement nested ≥6 control-flow levels |
| 756 | `jobs/sync_telecom_data.py:274:1` | statement nested ≥6 control-flow levels |
| 783 | `jobs/sync_transportation_data.py:330:1` | statement nested ≥6 control-flow levels |
| 1257 | `scripts/wtp_database.py:185:1` | statement nested ≥6 control-flow levels |

### BP-PY-24 — Django raw SQL With Format (6)

| Finding | Source | Reason |
| --- | --- | --- |
| 420 | `jobs/migrate_add_ai_summaries.py:69:19` | cursor.execute with f-string formatted SQL |
| 424 | `jobs/migrate_add_ai_summaries.py:72:19` | cursor.execute with f-string formatted SQL |
| 427 | `jobs/migrate_add_ai_summaries.py:79:19` | cursor.execute with f-string formatted SQL |
| 430 | `jobs/migrate_add_ai_summaries.py:82:19` | cursor.execute with f-string formatted SQL |
| 442 | `jobs/migrate_add_sanctions.py:58:23` | cursor.execute with f-string formatted SQL |
| 446 | `jobs/migrate_add_sanctions.py:61:23` | cursor.execute with f-string formatted SQL |

### CWE-1084 — Invokable Control Element with Excessive File or Data Access Operations (13)

| Finding | Source | Reason |
| --- | --- | --- |
| 174 | `jobs/audit_orphan_lobby_company_ids.py:48:12` | `main()` has 3× `db.execute` (≥3 open/.execute threshold) — reclassified from Uncertain |
| 285 | `jobs/detect_stories.py:1336:51` | `detect_contract_windfall` has 3× `db.execute` (≥3 threshold) — reclassified from Uncertain |
| 418 | `jobs/migrate_add_ai_summaries.py:52:15` | function performs many data-access operations |
| 435 | `jobs/migrate_add_indexes.py:116:12` | function performs many data-access operations |
| 440 | `jobs/migrate_add_sanctions.py:41:15` | `migrate()` has 4× `.execute` (≥3 threshold) — reclassified from Uncertain |
| 451 | `jobs/migrate_add_specific_issues.py:14:12` | `main()` has 3× `conn.execute` (≥3 threshold) — reclassified from Uncertain |
| 977 | `routers/ops.py:483:45` | `db_stats` has 3× `db.execute` (≥3 threshold) — reclassified from Uncertain |
| 1100 | `scripts/fix_finance_audit_20260417.py:147:42` | function performs many data-access operations |
| 1110 | `scripts/generate_lobbying_breakdown_stories.py:163:52` | `generate_sector_breakdown` has 3× `db.execute` (≥3 threshold) — reclassified from Uncertain |
| 1162 | `scripts/generate_tech_stories.py:68:26` | function performs many data-access operations |
| 1174 | `scripts/migrate_twitter_models.py:38:15` | `migrate()` has 5× `conn.execute` (≥3 threshold) — reclassified from Uncertain |
| 1204 | `scripts/remediate_published_stories.py:107:63` | `_derive_range_from_family` has 4× `db.execute` (≥3 threshold) — reclassified from Uncertain |
| 1222 | `scripts/remediate_short_stories.py:228:12` | function performs many data-access operations |

### CWE-1046 — Creation of Immutable Text Using String Concatenation (5)

| Finding | Source | Reason |
| --- | --- | --- |
| 284 | `jobs/detect_stories.py:1247:1` | string concatenated repeatedly in a loop |
| 392 | `jobs/generate_digest.py:384:1` | string concatenated repeatedly in a loop |
| 962 | `routers/og.py:136:1` | string concatenated repeatedly in a loop |
| 990 | `routers/ops.py:1122:1` | string concatenated repeatedly in a loop |
| 1127 | `scripts/generate_lobbying_breakdown_stories.py:435:1` | string concatenated repeatedly in a loop |

### PERF-PY-23 — Polymorphic Serialize Or Encode Inside Hot Loop (3)

| Finding | Source | Reason |
| --- | --- | --- |
| 231 | `jobs/backfill_sanctions_status.py:325:1` | json.dumps serialization inside a hot loop |
| 381 | `jobs/evaluate_legislative_claims.py:182:1` | json.dumps serialization inside a hot loop |
| 382 | `jobs/evaluate_legislative_claims.py:183:1` | json.dumps serialization inside a hot loop |

### PERF-PY-25 — Heavy Object Construction Per Homogeneous Element (3)

| Finding | Source | Reason |
| --- | --- | --- |
| 309 | `jobs/detect_stories.py:1806:1` | heavy object/lambda per homogeneous loop element |
| 383 | `jobs/evaluate_legislative_claims.py:195:1` | heavy object/lambda per homogeneous loop element |
| 1072 | `scripts/diagnose_uspto_odp.py:130:1` | heavy object/lambda per homogeneous loop element |

### PERF-PY-18 — Repeated Regex Rewrites On The Same Input (2)

| Finding | Source | Reason |
| --- | --- | --- |
| 208 | `jobs/backfill_logos_wikipedia.py:92:5` | repeated regex rewrites on the same input |
| 224 | `jobs/backfill_sanctions_status.py:90:5` | repeated regex rewrites on the same input |

### BP-PY-36 — SQLAlchemy Session Not Closed (2)

| Finding | Source | Reason |
| --- | --- | --- |
| 1214 | `scripts/remediate_published_stories.py:257:1` | SQLAlchemy Session created without with/close |
| 1231 | `scripts/run_story_gates_audit.py:56:1` | SQLAlchemy Session created without with/close |

### CWE-93 — Improper Neutralization of CRLF Sequences ('CRLF Injection') (2)

| Finding | Source | Reason |
| --- | --- | --- |
| 874 | `middleware/tracing.py:74:21` | dynamic value written to HTTP response header without CRLF neutralization |
| 916 | `routers/claims.py:31:13` | dynamic value written to HTTP response header without CRLF neutralization |

### CWE-367 — Time-of-check Time-of-use (TOCTOU) Race Condition (2)

| Finding | Source | Reason |
| --- | --- | --- |
| 902 | `routers/bulk.py:261:12` | filesystem path checked before later separate use |
| 959 | `routers/metrics.py:153:12` | filesystem path checked before later separate use |

### BP-PY-30 — FastAPI Blocking I/O In Async Route (2)

| Finding | Source | Reason |
| --- | --- | --- |
| 995 | `routers/ops.py:2067:1` | sync ORM call in async FastAPI route |
| 996 | `routers/ops.py:2101:1` | sync ORM call in async FastAPI route |

### BP-PY-42 — unittest Assert Without Context On Raises (2)

| Finding | Source | Reason |
| --- | --- | --- |
| 1465 | `tests/chaos/test_external_api_failure.py:174:1` | try/except used to expect failure in test |
| 1469 | `tests/performance/query_analysis.py:263:1` | try/except used to expect failure in test |

### BP-PY-7 — open Without Context Manager (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 865 | `main.py:54:27` | open() without context manager |

### CWE-290 — Authentication Bypass by Spoofing (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 878 | `middleware/tracing.py:97:25` | client-provided X-Forwarded-For trusted directly |

### BP-PY-14 — requests Without Timeout (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 1239 | `scripts/send_welcome.py:29:5` | requests call missing timeout= |

### CWE-359 — Exposure of Private Personal Information to an Unauthorized Actor (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 898 | `routers/auth.py:1331:21` | personal data field written to log sink |

### BP-PY-29 — FastAPI Depends On Mutable Global (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 937 | `routers/influence.py:399:1` | FastAPI route mutates module-level global |

### BP-PY-40 — threading Without Join Or Shutdown (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 1008 | `routers/politics_bills.py:72:6` | threading.Thread started without join |

### PERF-PY-27 — Repeated Load Of Same Filesystem Path (1)

| Finding | Source | Reason |
| --- | --- | --- |
| 1108 | `scripts/gen_mobile_sector_screens.py:385:1` | same path re-read/parsed without cache |

## Final evidence

- Delegated reviewers: none (single-reviewer audit)
- Chunk evidence: `./scripts/chunks` (60 chunk files, findings 1–1492)
- Function evidence: `./scripts/findings/functions` (1492 context files)
- Validation: `git diff --check` — pass