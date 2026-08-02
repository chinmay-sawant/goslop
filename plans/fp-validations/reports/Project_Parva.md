# False-positive audit — Project_Parva

## Run metadata

```yaml
timestamp: 2026-08-02T07:10:00Z
repository: Project_Parva
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva
branch: main
commit: d05f6111bb0a39ce8dc3c82330297b60ae82c7c5
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva
chunk_path: ./scripts/chunks
function_context_path: ./scripts/findings/functions
```

## Scan evidence

- Build command: `make build` (goslop binary at `./bin/goslop`)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir real-repos/Project_Parva/scripts/chunks -context-dir real-repos/Project_Parva/scripts/findings/functions real-repos/Project_Parva`
- Findings: `412`
- Chunks reviewed: `./scripts/chunks/Chunk_1_25.txt` … `Chunk_401_412.txt` (17 files, findings 1–412)
- Function contexts reviewed: `./scripts/findings/functions/<id>.txt` for every proposed false positive; enclosing source read when the exported context was insufficient (billing `service.py`/`storage.py`/`migrations.py`, `middleware.py`, `app_factory.py`, `rate_limit.py`, `public_artifacts_routes.py`, `rulelang_service.py`, `triad_pipeline.py`, `week3_ground_truth_pipeline.py`, `harness.py`, `ingest_moha_pdfs.py`, `verify_public.py`, `client.py`, `run_browser_smoke.py`, `check_public_claims.py`, `verify_clean_clone_assumptions.py`, `test_temporal_trust_tools.py`, `test_final_artifacts_exist.py`, `tithi.py`, `validate_schemas.py`).

## Audit checklist

- [x] Read every assigned chunk under `./scripts/chunks` (all 17 chunks, findings 1–412).
- [x] Read `./scripts/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition (`-explain` metadata + catalogue `detection_notes` in `ruleset/python/`) and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient (none required).
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 401 | 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 63, 64, 65, 66, 67, 68, 69, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 125, 126, 127, 128, 129, 130, 131, 132, 135, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 255, 256, 257, 258, 259, 260, 261, 262, 263, 264, 265, 266, 267, 268, 269, 270, 271, 272, 274, 275, 276, 277, 278, 279, 280, 281, 282, 283, 284, 285, 286, 287, 288, 289, 290, 291, 292, 293, 294, 295, 296, 297, 298, 299, 300, 301, 302, 303, 304, 305, 306, 307, 308, 309, 310, 311, 312, 313, 314, 316, 317, 318, 319, 320, 321, 322, 323, 324, 325, 326, 327, 328, 329, 330, 331, 332, 333, 334, 335, 336, 337, 338, 339, 340, 341, 342, 343, 344, 345, 346, 347, 348, 349, 350, 351, 352, 353, 354, 355, 356, 357, 358, 359, 360, 361, 362, 363, 364, 365, 366, 367, 368, 369, 370, 371, 372, 373, 374, 375, 376, 377, 378, 379, 380, 381, 382, 383, 384, 385, 386, 387, 388, 389, 390, 391, 392, 393, 394, 395, 396, 397, 398, 399, 400, 401, 402, 403, 404, 405, 406, 407, 408, 409, 410, 411, 412 |
| True positive | 11 | 41, 62, 70, 88, 124, 133, 134, 136, 137, 273, 315 |
| Uncertain | 0 | — |

## False positives

One subsection per finding. No grouping was applied: no two false-positive findings reference the exact same source construct (same file, same line, same rule, same trigger).

### [ ] Finding 1 — BP-PY-13

- Function context: `./scripts/findings/functions/1.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/api/billing_routes.py:327:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
raise HTTPException(status_code=403, detail="Webhook notifications require Professional or Enterprise.")
secret = f"parva_whsec_{secrets.token_urlsafe(32)}"
return await _billing_call(
```

Why this is a false positive: the 'secret' is generated at runtime with `secrets.token_urlsafe(32)`; no secret-like string literal is assigned, so the hardcoded-literal rule condition is not satisfied.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 2 — CWE-396

- Function context: `./scripts/findings/functions/2.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/api/engine_routes.py:278:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
value = date(year, month, day)
except Exception as exc:
raise HTTPException(status_code=400, detail="Invalid date format. Use YYYY-MM-DD") from exc
```

Why this is a false positive: the handler converts the failure into an explicit HTTP 400 and re-raises via `from exc`; the failure is deliberately handled, not hidden.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 3 — BP-PY-32

- Function context: `./scripts/findings/functions/3.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/api/public_artifacts_routes.py:100:12`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
raise HTTPException(status_code=404, detail="Artifact not found")
return FileResponse(path, media_type="application/json", filename=filename)

```

Why this is a false positive: the user-supplied `filename` is confined by a basename-only policy: `/` and `..` are rejected and a `.json` suffix enforced before joining under the constant `PRECOMPUTED_DIR`, satisfying the confinement requirement.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 4 — BP-PY-32

- Function context: `./scripts/findings/functions/4.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/api/public_artifacts_routes.py:110:12`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
raise HTTPException(status_code=404, detail="Dashboard artifact not generated yet")
return FileResponse(path, media_type="application/json", filename=path.name)

```

Why this is a false positive: path built from module constants (`REPORTS_DIR`/`PUBLIC_ARTIFACTS_DIR`) plus a fixed filename; no user input reaches `FileResponse`.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 5 — BP-PY-32

- Function context: `./scripts/findings/functions/5.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/api/public_artifacts_routes.py:118:12`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
raise HTTPException(status_code=404, detail="Source review queue artifact not generated yet")
return FileResponse(path, media_type="application/json", filename=path.name)

```

Why this is a false positive: constant directory plus fixed filename; no user input.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 6 — BP-PY-32

- Function context: `./scripts/findings/functions/6.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/api/public_artifacts_routes.py:126:12`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
raise HTTPException(status_code=404, detail="Boundary suite artifact not generated yet")
return FileResponse(path, media_type="application/json", filename=path.name)

```

Why this is a false positive: constant directory plus fixed filename; no user input.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 7 — BP-PY-32

- Function context: `./scripts/findings/functions/7.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/api/public_artifacts_routes.py:134:12`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
raise HTTPException(status_code=404, detail="Differential artifact not generated yet")
return FileResponse(path, media_type="application/json", filename=path.name)

```

Why this is a false positive: constant `data_dir()`/`differential` path plus fixed filename; no user input.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 8 — BP-PY-41

- Function context: `./scripts/findings/functions/8.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/api/rules_routes.py:112:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
@router.post("/{rule_id}/test")
async def test_rule(rule_id: str, request: Request) -> dict[str, Any]:
try:
```

Why this is a false positive: `test_rule` is a FastAPI route handler, not a pytest test function; the placeholder-test condition does not apply.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 9 — CWE-89

- Function context: `./scripts/findings/functions/9.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/billing/migrations.py:224:10`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
schema = SQLITE_SCHEMA[0] if store.config.dialect == "sqlite" else POSTGRES_SCHEMA[0]
store.execute(schema)
for migration in MIGRATIONS:
```

Why this is a false positive: the SQL passed to `execute` is a static schema constant selected at migration time; no data is interpolated.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 10 — BP-PY-37

- Function context: `./scripts/findings/functions/10.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/billing/migrations.py:235:14`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
store.execute(statement)
store.execute(
f"INSERT INTO billing_schema_migrations (version, name) VALUES ({store.param()}, {store.param()})",
```

Why this is a false positive: the f-string interpolates only `store.param()` placeholder tokens (`?` / `%s`); every value is passed as a bound parameter in the `execute(sql, params)` call — exactly the rule's prescribed fix, so no data is `%`-formatted into the SQL.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 11 — BP-PY-37

- Function context: `./scripts/findings/functions/11.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/billing/service.py:112:19`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
scrubbed_metadata = scrub_structured_trace(metadata or {})
self.store.execute(
f"""
```

Why this is a false positive: the f-string interpolates only `store.param()` placeholder tokens (`?` / `%s`); every value is passed as a bound parameter in the `execute(sql, params)` call — exactly the rule's prescribed fix, so no data is `%`-formatted into the SQL.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 12 — CWE-89

- Function context: `./scripts/findings/functions/12.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/billing/service.py:112:19`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
scrubbed_metadata = scrub_structured_trace(metadata or {})
self.store.execute(
f"""
```

Why this is a false positive: no value is interpolated into the SQL text: either a static schema constant or placeholder-only f-strings with all values passed as bound parameters through the `(sql, params)` wrapper.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 13 — BP-PY-37

- Function context: `./scripts/findings/functions/13.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/billing/service.py:193:23`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if existing:
self.store.execute(
f"""
```

Why this is a false positive: the f-string interpolates only `store.param()` placeholder tokens (`?` / `%s`); every value is passed as a bound parameter in the `execute(sql, params)` call — exactly the rule's prescribed fix, so no data is `%`-formatted into the SQL.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 14 — BP-PY-37

- Function context: `./scripts/findings/functions/14.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/billing/service.py:208:19`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
customer_id = _new_id("cus")
self.store.execute(
f"""
```

Why this is a false positive: the f-string interpolates only `store.param()` placeholder tokens (`?` / `%s`); every value is passed as a bound parameter in the `execute(sql, params)` call — exactly the rule's prescribed fix, so no data is `%`-formatted into the SQL.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 15 — BP-PY-37

- Function context: `./scripts/findings/functions/15.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/billing/service.py:252:19`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
due_at = (utc_now() + timedelta(days=7)).isoformat()
self.store.execute(
f"""
```

Why this is a false positive: the f-string interpolates only `store.param()` placeholder tokens (`?` / `%s`); every value is passed as a bound parameter in the `execute(sql, params)` call — exactly the rule's prescribed fix, so no data is `%`-formatted into the SQL.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 16 — BP-PY-37

- Function context: `./scripts/findings/functions/16.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/billing/service.py:259:19`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
)
self.store.execute(
f"""
```

Why this is a false positive: the f-string interpolates only `store.param()` placeholder tokens (`?` / `%s`); every value is passed as a bound parameter in the `execute(sql, params)` call — exactly the rule's prescribed fix, so no data is `%`-formatted into the SQL.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 17 — BP-PY-37

- Function context: `./scripts/findings/functions/17.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/billing/service.py:276:19`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
)
self.store.execute(
f"""
```

Why this is a false positive: the f-string interpolates only `store.param()` placeholder tokens (`?` / `%s`); every value is passed as a bound parameter in the `execute(sql, params)` call — exactly the rule's prescribed fix, so no data is `%`-formatted into the SQL.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 18 — BP-PY-37

- Function context: `./scripts/findings/functions/18.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/billing/service.py:382:19`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
payment_status = "failed"
self.store.execute(
f"""
```

Why this is a false positive: the f-string interpolates only `store.param()` placeholder tokens (`?` / `%s`); every value is passed as a bound parameter in the `execute(sql, params)` call — exactly the rule's prescribed fix, so no data is `%`-formatted into the SQL.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 19 — BP-PY-37

- Function context: `./scripts/findings/functions/19.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/billing/service.py:397:23`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
elif payment_status == "failed":
self.store.execute(
f"UPDATE invoices SET status = 'failed', updated_at = {self.store.param()} WHERE id = {self.store.param()}",
```

Why this is a false positive: the f-string interpolates only `store.param()` placeholder tokens (`?` / `%s`); every value is passed as a bound parameter in the `execute(sql, params)` call — exactly the rule's prescribed fix, so no data is `%`-formatted into the SQL.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 20 — BP-PY-37

- Function context: `./scripts/findings/functions/20.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/billing/service.py:420:19`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
renews_at = (now_dt + timedelta(days=30)).isoformat()
self.store.execute(
f"""
```

Why this is a false positive: the f-string interpolates only `store.param()` placeholder tokens (`?` / `%s`); every value is passed as a bound parameter in the `execute(sql, params)` call — exactly the rule's prescribed fix, so no data is `%`-formatted into the SQL.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 21 — BP-PY-37

- Function context: `./scripts/findings/functions/21.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/billing/service.py:429:19`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
)
self.store.execute(
f"""
```

Why this is a false positive: the f-string interpolates only `store.param()` placeholder tokens (`?` / `%s`); every value is passed as a bound parameter in the `execute(sql, params)` call — exactly the rule's prescribed fix, so no data is `%`-formatted into the SQL.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 22 — BP-PY-37

- Function context: `./scripts/findings/functions/22.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/billing/service.py:480:19`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
key_id = _new_id("key")
self.store.execute(
f"""
```

Why this is a false positive: the f-string interpolates only `store.param()` placeholder tokens (`?` / `%s`); every value is passed as a bound parameter in the `execute(sql, params)` call — exactly the rule's prescribed fix, so no data is `%`-formatted into the SQL.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 23 — BP-PY-37

- Function context: `./scripts/findings/functions/23.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/billing/service.py:554:19`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
raise BillingAuthError(403, "Subscription is not active.")
self.store.execute(
f"UPDATE api_keys SET last_used_at = {self.store.param()} WHERE id = {self.store.param()}",
```

Why this is a false positive: the f-string interpolates only `store.param()` placeholder tokens (`?` / `%s`); every value is passed as a bound parameter in the `execute(sql, params)` call — exactly the rule's prescribed fix, so no data is `%`-formatted into the SQL.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 24 — BP-PY-37

- Function context: `./scripts/findings/functions/24.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/billing/service.py:593:23`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if existing:
self.store.execute(
f"""
```

Why this is a false positive: the f-string interpolates only `store.param()` placeholder tokens (`?` / `%s`); every value is passed as a bound parameter in the `execute(sql, params)` call — exactly the rule's prescribed fix, so no data is `%`-formatted into the SQL.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 25 — BP-PY-37

- Function context: `./scripts/findings/functions/25.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/billing/service.py:602:23`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
else:
self.store.execute(
f"""
```

Why this is a false positive: the f-string interpolates only `store.param()` placeholder tokens (`?` / `%s`); every value is passed as a bound parameter in the `execute(sql, params)` call — exactly the rule's prescribed fix, so no data is `%`-formatted into the SQL.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 26 — BP-PY-37

- Function context: `./scripts/findings/functions/26.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/billing/service.py:654:19`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
now = iso_now()
self.store.execute(
f"UPDATE api_keys SET active = {self.store.param()}, revoked_at = {self.store.param()} WHERE id = {self.store.param()}",
```

Why this is a false positive: the f-string interpolates only `store.param()` placeholder tokens (`?` / `%s`); every value is passed as a bound parameter in the `execute(sql, params)` call — exactly the rule's prescribed fix, so no data is `%`-formatted into the SQL.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 27 — BP-PY-37

- Function context: `./scripts/findings/functions/27.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/billing/service.py:678:19`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
secret_hash = hash_api_key_secret(secret, self.settings.api_key_pepper)
self.store.execute(
f"""
```

Why this is a false positive: the f-string interpolates only `store.param()` placeholder tokens (`?` / `%s`); every value is passed as a bound parameter in the `execute(sql, params)` call — exactly the rule's prescribed fix, so no data is `%`-formatted into the SQL.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 28 — BP-PY-37

- Function context: `./scripts/findings/functions/28.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/billing/service.py:732:19`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
now = iso_now()
self.store.execute(
f"""
```

Why this is a false positive: the f-string interpolates only `store.param()` placeholder tokens (`?` / `%s`); every value is passed as a bound parameter in the `execute(sql, params)` call — exactly the rule's prescribed fix, so no data is `%`-formatted into the SQL.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 29 — BP-PY-37

- Function context: `./scripts/findings/functions/29.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/billing/service.py:741:23`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if row["payment_id"]:
self.store.execute(
f"""
```

Why this is a false positive: the f-string interpolates only `store.param()` placeholder tokens (`?` / `%s`); every value is passed as a bound parameter in the `execute(sql, params)` call — exactly the rule's prescribed fix, so no data is `%`-formatted into the SQL.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 30 — BP-PY-37

- Function context: `./scripts/findings/functions/30.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/billing/service.py:788:19`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
renews_at = (base + timedelta(days=max(1, min(days, 366)))).isoformat()
self.store.execute(
f"""
```

Why this is a false positive: the f-string interpolates only `store.param()` placeholder tokens (`?` / `%s`); every value is passed as a bound parameter in the `execute(sql, params)` call — exactly the rule's prescribed fix, so no data is `%`-formatted into the SQL.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 31 — CWE-89

- Function context: `./scripts/findings/functions/31.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/billing/storage.py:108:17`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
with self.connect() as conn:
conn.execute(sql, params)

```

Why this is a false positive: this is the parameterized-query wrapper itself: it forwards both SQL and bound params, and every caller passes values via params, so no unbound data reaches the cursor.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 32 — PERF-PY-23

- Function context: `./scripts/findings/functions/32.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/billing/storage.py:131:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
for plan in PLAN_DEFINITIONS:
features_json = json.dumps(list(plan.features), separators=(",", ":"))
if self.config.dialect == "sqlite":
```

Why this is a false positive: serialization runs in a startup seeding loop over a small constant plan list, or per-row in an offline research script producing report rows — not a hot service loop.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 33 — BP-PY-32

- Function context: `./scripts/findings/functions/33.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/bootstrap/app_factory.py:77:16`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
def _frontend_response(path, *, immutable: bool = False) -> FileResponse:
response = FileResponse(path)
response.headers["Cache-Control"] = (
```

Why this is a false positive: the helper serves the frontend build directory; callers pass `settings.frontend_dist`-rooted paths, not request input.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 34 — CWE-93

- Function context: `./scripts/findings/functions/34.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/bootstrap/app_factory.py:78:13`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
response = FileResponse(path)
response.headers["Cache-Control"] = (
"public, max-age=31536000, immutable" if immutable else "no-cache"
```

Why this is a false positive: the header value is a fixed string literal chosen by a boolean; nothing an attacker could inject CRLF into.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 35 — CWE-93

- Function context: `./scripts/findings/functions/35.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/bootstrap/middleware.py:61:13`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
existing = response.headers.get("Link")
response.headers["Link"] = f"{existing}, {value}" if existing else value

```

Why this is a false positive: callers pass `source_url` from `settings.source_url` (server configuration) or the constant `'</v3/docs>; rel="successor-version"'`; no request-controlled value reaches the header.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 36 — CWE-290

- Function context: `./scripts/findings/functions/36.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/bootstrap/middleware.py:181:21`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
remote_host = request.client.host if request.client and request.client.host else ""
forwarded_for = request.headers.get("x-forwarded-for")
trusted_proxy_ips = settings.trusted_proxy_ips
```

Why this is a false positive: X-Forwarded-For is honored only when `trust_proxy_headers` is true (direct peer in the configured trusted-proxy allowlist or `*`); the header is not trusted directly.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 37 — CWE-396

- Function context: `./scripts/findings/functions/37.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/bootstrap/middleware.py:204:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
response = await call_next(request)
except Exception as exc:
latency_ms = round((time.perf_counter() - started) * 1000.0, 2)
```

Why this is a false positive: the handler logs the exception with full context and re-raises (`logger.exception(...)` + `raise`); the failure is surfaced.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 38 — BP-PY-1

- Function context: `./scripts/findings/functions/38.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/bootstrap/middleware.py:794:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
ephemeris_value = "jpl-de440-lahiri-sidereal"
except Exception as exc:  # noqa: BLE001
logger.warning("Unable to resolve active future-BS ephemeris label: %s", exc)
```

Why this is a false positive: the exception object is recorded in `logger.warning(..., exc)`; the handler deliberately degrades a header value and the failure is not silently discarded.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 39 — BP-PY-12

- Function context: `./scripts/findings/functions/39.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/bootstrap/rate_limit.py:166:23`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
client = self._get_client()
return client.eval(
self._ATOMIC_CHECK_SCRIPT,
```

Why this is a false positive: `client.eval` is the redis Lua evaluation API invoked with the module-constant `_ATOMIC_CHECK_SCRIPT`, not Python `eval`/`exec` on dynamic input.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 40 — CWE-396

- Function context: `./scripts/findings/functions/40.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/bootstrap/rate_limit.py:197:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
)
except Exception as exc:  # noqa: BLE001 - fail closed when the shared limiter is unavailable.
raise RateLimiterUnavailable(
```

Why this is a false positive: deliberate fail-closed conversion with `raise RateLimiterUnavailable(...) from exc`; the failure is propagated.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 42 — BP-PY-46

- Function context: `./scripts/findings/functions/42.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/calendar/calculator.py:17:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
>>> dashain_2026 = calculate_festival_date("dashain", 2026)
>>> print(dashain_2026)

```

Why this is a false positive: the `print` lives in a standalone script / CLI module (scripts, tools, examples, CLI entry) whose output mechanism is stdout, or inside a docstring doctest example — not operational logging in a non-script library module.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 43 — BP-PY-46

- Function context: `./scripts/findings/functions/43.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/calendar/calculator.py:347:13`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
>>> dashain = calculate_festival_date("dashain", 2026)
>>> print(dashain.start)  # Returns October 2, 2026

```

Why this is a false positive: the `print` lives in a standalone script / CLI module (scripts, tools, examples, CLI entry) whose output mechanism is stdout, or inside a docstring doctest example — not operational logging in a non-script library module.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 44 — BP-PY-46

- Function context: `./scripts/findings/functions/44.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/calendar/calculator.py:350:13`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
>>> tihar = calculate_festival_date("tihar", 2026)
>>> print(tihar.duration_days)  # Returns 5
"""
```

Why this is a false positive: the `print` lives in a standalone script / CLI module (scripts, tools, examples, CLI entry) whose output mechanism is stdout, or inside a docstring doctest example — not operational logging in a non-script library module.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 45 — BP-PY-2

- Function context: `./scripts/findings/functions/45.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/calendar/calculator.py:366:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
return DateRange(start=override_date, end=end_date, year=gregorian_year)
except (ImportError, OSError, JSONDecodeError, TypeError, ValueError, KeyError):
# If overrides aren't available, fall back to calculations
```

Why this is a false positive: typed-exception tuple with a comment documenting the fallback to algorithmic calculation when override data is unavailable; the pass is the intended fallback.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 46 — CWE-390

- Function context: `./scripts/findings/functions/46.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/calendar/calculator.py:366:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
return DateRange(start=override_date, end=end_date, year=gregorian_year)
except (ImportError, OSError, JSONDecodeError, TypeError, ValueError, KeyError):
# If overrides aren't available, fall back to calculations
```

Why this is a false positive: the handler's purpose is the documented fallback to algorithmic calculation; recovery exists, so 'error condition without action' does not hold.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 47 — BP-PY-46

- Function context: `./scripts/findings/functions/47.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/calendar/calculator.py:511:17`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
>>> for festival_id, dates in upcoming:
...     print(f"{festival_id}: {dates.start}")
"""
```

Why this is a false positive: the `print` lives in a standalone script / CLI module (scripts, tools, examples, CLI entry) whose output mechanism is stdout, or inside a docstring doctest example — not operational logging in a non-script library module.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 48 — BP-PY-46

- Function context: `./scripts/findings/functions/48.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/calendar/calculator.py:592:13`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
>>> next_dashain = get_next_occurrence("dashain")
>>> print(f"Dashain starts: {next_dashain.start}")
"""
```

Why this is a false positive: the `print` lives in a standalone script / CLI module (scripts, tools, examples, CLI entry) whose output mechanism is stdout, or inside a docstring doctest example — not operational logging in a non-script library module.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 49 — BP-PY-2

- Function context: `./scripts/findings/functions/49.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/calendar/calculator_v2.py:173:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
)
except (ImportError, OSError, JSONDecodeError, TypeError, ValueError, KeyError):
# Fall back to algorithmic calculation
```

Why this is a false positive: documented deliberate fallback to algorithmic calculation when override data is unavailable.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 50 — CWE-390

- Function context: `./scripts/findings/functions/50.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/calendar/calculator_v2.py:173:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
)
except (ImportError, OSError, JSONDecodeError, TypeError, ValueError, KeyError):
# Fall back to algorithmic calculation
```

Why this is a false positive: the handler triggers the documented fallback recovery path.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 51 — BP-PY-46

- Function context: `./scripts/findings/functions/51.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/calendar/panchanga.py:81:13`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
>>> panchanga = get_panchanga(date(2026, 2, 6))
>>> print(panchanga['tithi']['name'])
'Panchami'
```

Why this is a false positive: the `print` lives in a standalone script / CLI module (scripts, tools, examples, CLI entry) whose output mechanism is stdout, or inside a docstring doctest example — not operational logging in a non-script library module.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 52 — BP-PY-5

- Function context: `./scripts/findings/functions/52.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/calendar/tithi.py:9:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

from app.calendar.tithi import *  # noqa: F401,F403

```

Why this is a false positive: the module is an explicitly documented compatibility re-export stub for the canonical `tithi` package (`# noqa: F401,F403`), the documented re-export pattern the rule's note allows.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 53 — CWE-186

- Function context: `./scripts/findings/functions/53.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/compliance/notice_ingestion.py:14:14`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

BS_DATE_RE = re.compile(r"(?P<year>20\d{2})-(?P<month>\d{2})-(?P<day>\d{2})")
FIELD_RE = re.compile(r"^(?P<key>issuer|published|effective|deadline|action|affected_party|jurisdiction):\s*(?P<value>.+)$", re.I)
```

Why this is a false positive: the regex is a parser for the fixed BS date format of the ingestion sources; the narrow shape is the intended format, not an over-restriction.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 54 — CWE-208

- Function context: `./scripts/findings/functions/54.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/membranes/source_resolution.py:98:59`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
review_required=True,
source_docket_ids=resolution.source_docket_ids if resolution.authority == AuthorityTaint.STATIC_REFERENCE else (),
source_refs=resolution.source_refs if resolution.authority == AuthorityTaint.STATIC_REFERENCE else (),
```

Why this is a false positive: the `==` compares enum members (`AuthorityTaint`), not secrets or authentication values.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 55 — CWE-208

- Function context: `./scripts/findings/functions/55.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/policy/vm.py:64:21`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
static_rule == "explicit_mode_or_compare_branch_only"
and candidate.authority == AuthorityTaint.STATIC_REFERENCE
and len(candidates) > 1
```

Why this is a false positive: enum member comparison, not a secret comparison.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 56 — PERF-PY-26

- Function context: `./scripts/findings/functions/56.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/research/future_bs/accuracy_architecture.py:219:15`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
hard_cases = build_hard_case_benchmark()
lattice = decode_month_start_lattice()
green = certify_green_predictions()
```

Why this is a false positive: the decode runs once per report build in offline research tooling, not on a hot request path.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 57 — BP-PY-2

- Function context: `./scripts/findings/functions/57.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/research/future_bs/backtest.py:282:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
mismatch_by_ingress_hour[hour] += 1
except ValueError:
pass
```

Why this is a false positive: the handler tolerates malformed diagnostic rows in an offline analytics loop; the unparseable row is deliberately skipped without breaking the aggregate.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 58 — CWE-390

- Function context: `./scripts/findings/functions/58.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/research/future_bs/backtest.py:282:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
mismatch_by_ingress_hour[hour] += 1
except ValueError:
pass
```

Why this is a false positive: skipping the malformed row is the defined handling; the pass is the intended action.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 59 — CWE-1071

- Function context: `./scripts/findings/functions/59.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/research/future_bs/backtest.py:282:21`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
mismatch_by_ingress_hour[hour] += 1
except ValueError:
pass
```

Why this is a false positive: the empty block implements the deliberate malformed-row skip in batch analytics.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 60 — PERF-PY-25

- Function context: `./scripts/findings/functions/60.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/research/future_bs/hamropatro_shadow.py:176:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
disagreements,
key=lambda item: (item["model"], int(item["bs_year"]), int(item["bs_month"])),
):
```

Why this is a false positive: sort-key lambda with per-element `int()` conversion in a one-shot research script; the work is inherent to sorting, not a hot path.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 61 — PERF-PY-26

- Function context: `./scripts/findings/functions/61.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/research/future_bs/shadow_residual_correction.py:96:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
continue
month, year_mod4, base_days = parse_residual_key(key)
rules[key] = {
```

Why this is a false positive: per-element key parsing in an offline model-calibration script; each key is a distinct element and the parse is required.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 63 — PERF-PY-27

- Function context: `./scripts/findings/functions/63.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/rules/triad_pipeline.py:125:12`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
validation_path = triad_paths(rule.festival_id)["validation"]
payload = json.loads(validation_path.read_text(encoding="utf-8"))
if any(case.get("status") == "passed" for case in payload.get("cases", [])):
```

Why this is a false positive: the path is derived per rule (`rule.festival_id`), so each loop iteration loads a different file; the rule targets repeated loads of the same path.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 64 — PERF-PY-27

- Function context: `./scripts/findings/functions/64.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/rules/triad_pipeline.py:150:23`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

validation_payload = json.loads(paths["validation"].read_text(encoding="utf-8"))
if any(case.get("status") == "passed" for case in validation_payload.get("cases", [])):
```

Why this is a false positive: same as 63: the path varies per rule iteration.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 65 — CWE-186

- Function context: `./scripts/findings/functions/65.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/services/agent_service.py:71:18`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

BS_AD_CLAIM_RE = re.compile(
r"(?P<bs>\d{4}-\d{2}-\d{2})\s*(?:BS|B\.S\.|Bikram Sambat)?\s*"
```

Why this is a false positive: the regex parses the fixed BS/AD date-claim format the feature documents; the narrow shape is the intended grammar.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 66 — BP-PY-1

- Function context: `./scripts/findings/functions/66.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/services/rulelang_service.py:415:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
evidence_packet_id = packet["packet_id"]
except Exception as exc:  # noqa: BLE001
context.warnings.append(f"rule_evidence_packet_unavailable: {exc}")
```

Why this is a false positive: the exception text is recorded in `context.warnings`; the failure is surfaced to the caller.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 67 — CWE-396

- Function context: `./scripts/findings/functions/67.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/services/rulelang_service.py:415:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
evidence_packet_id = packet["packet_id"]
except Exception as exc:  # noqa: BLE001
context.warnings.append(f"rule_evidence_packet_unavailable: {exc}")
```

Why this is a false positive: the failure is recorded in the evaluation context warnings.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 68 — BP-PY-1

- Function context: `./scripts/findings/functions/68.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/services/rulelang_service.py:523:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
)
except Exception as exc:  # noqa: BLE001
results.append(
```

Why this is a false positive: the failure is recorded per case with `str(exc)` in the results the caller consumes.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 69 — BP-PY-1

- Function context: `./scripts/findings/functions/69.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/services/rulelang_service.py:761:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
fact_ids = [bs_ad_fact_id(year, month, day)]
except Exception as exc:  # noqa: BLE001
value = {"valid": False, "error": str(exc)}
```

Why this is a false positive: the exception text is captured into the returned error value (`{"valid": False, "error": str(exc)}`).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 71 — BP-PY-2

- Function context: `./scripts/findings/functions/71.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/services/rulelang_service.py:1384:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
return True
except ValueError:
pass
```

Why this is a false positive: unparseable AD date strings are deliberately skipped while scanning candidate dates.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 72 — CWE-390

- Function context: `./scripts/findings/functions/72.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/services/rulelang_service.py:1384:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
return True
except ValueError:
pass
```

Why this is a false positive: skipping the malformed date is the defined handling.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 73 — CWE-1071

- Function context: `./scripts/findings/functions/73.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/services/rulelang_service.py:1384:13`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
return True
except ValueError:
pass
```

Why this is a false positive: the empty handler implements the intended skip of unparseable dates.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 74 — BP-PY-2

- Function context: `./scripts/findings/functions/74.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/services/timegraph_fact_links.py:39:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
fact_ids.append(fiscal_period_fact_id(year, month, day))
except (KeyError, TypeError, ValueError):
pass
```

Why this is a false positive: typed-exception tuple; malformed date tuples are deliberately skipped while building fact-id lists.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 75 — CWE-390

- Function context: `./scripts/findings/functions/75.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/services/timegraph_fact_links.py:39:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
fact_ids.append(fiscal_period_fact_id(year, month, day))
except (KeyError, TypeError, ValueError):
pass
```

Why this is a false positive: skip of malformed rows is the defined handling.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 76 — BP-PY-45

- Function context: `./scripts/findings/functions/76.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/benchmark_runner.py:19:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

sys.path.insert(0, str(Path(__file__).parent.parent))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 77 — BP-PY-1

- Function context: `./scripts/findings/functions/77.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/build_baseline_supplement.py:157:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
converted_dates.append(g_date.isoformat())
except Exception as exc:  # noqa: BLE001
err = str(exc)
```

Why this is a false positive: the exception text is captured into the conversion error record for reporting.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 78 — CWE-396

- Function context: `./scripts/findings/functions/78.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/build_baseline_supplement.py:157:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
converted_dates.append(g_date.isoformat())
except Exception as exc:  # noqa: BLE001
err = str(exc)
```

Why this is a false positive: the failure is recorded per row.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 79 — PERF-PY-25

- Function context: `./scripts/findings/functions/79.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/build_baseline_supplement.py:195:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
records.append(
BaselineRecord(
bs_year=bs_year,
```

Why this is a false positive: per-row record construction is inherent to the batch conversion tool; not a hot path.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 80 — PERF-PY-26

- Function context: `./scripts/findings/functions/80.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/build_baseline_supplement.py:238:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
parser.add_argument("--output", default=None)
args = parser.parse_args()

```

Why this is a false positive: argparse runs once at CLI startup.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 81 — BP-PY-45

- Function context: `./scripts/findings/functions/81.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/build_boundary_suite.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 82 — BP-PY-46

- Function context: `./scripts/findings/functions/82.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/build_overrides.py:63:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
save_json(OVERRIDES_PATH, overrides)
print(f"✅ Updated overrides: {OVERRIDES_PATH}")

```

Why this is a false positive: standalone CLI tool in `backend/tools/`; print is the tool's CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 83 — BP-PY-45

- Function context: `./scripts/findings/functions/83.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/build_source_review_queue.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 84 — BP-PY-46

- Function context: `./scripts/findings/functions/84.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/evaluate.py:281:13`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if verbose:
print(f"[{i+1}/{len(test_cases)}] Evaluating {fid}...", end=" ")

```

Why this is a false positive: standalone CLI evaluation tool; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 85 — BP-PY-46

- Function context: `./scripts/findings/functions/85.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/evaluate.py:291:17`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if result.calculated_date:
print(
f"{status} (expected {expected}, got {result.calculated_date}, Δ={result.variance_days}d)"
```

Why this is a false positive: standalone CLI evaluation tool; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 86 — BP-PY-46

- Function context: `./scripts/findings/functions/86.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/evaluate.py:295:17`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
else:
print(f"{status} ERROR: {result.error}")

```

Why this is a false positive: standalone CLI evaluation tool; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 87 — BP-PY-46

- Function context: `./scripts/findings/functions/87.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/evaluate.py:333:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

print(f"Results saved to {output_path}")

```

Why this is a false positive: standalone CLI evaluation tool; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 89 — BP-PY-46

- Function context: `./scripts/findings/functions/89.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/evaluate_v2.py:114:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
"""Print evaluation results."""
print("\n" + "=" * 70)
print("FESTIVAL DATE CALCULATOR V2 - EVALUATION")
```

Why this is a false positive: standalone CLI tool; `print_results`' stated purpose is printing results.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 90 — BP-PY-46

- Function context: `./scripts/findings/functions/90.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/evaluate_v2.py:115:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
print("\n" + "=" * 70)
print("FESTIVAL DATE CALCULATOR V2 - EVALUATION")
print("(Using CORRECT Lunar Month Model)")
```

Why this is a false positive: CLI evaluation tool; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 91 — BP-PY-46

- Function context: `./scripts/findings/functions/91.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/evaluate_v2.py:116:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
print("FESTIVAL DATE CALCULATOR V2 - EVALUATION")
print("(Using CORRECT Lunar Month Model)")
print("=" * 70 + "\n")
```

Why this is a false positive: CLI evaluation tool; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 92 — BP-PY-46

- Function context: `./scripts/findings/functions/92.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/evaluate_v2.py:117:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
print("(Using CORRECT Lunar Month Model)")
print("=" * 70 + "\n")

```

Why this is a false positive: CLI evaluation tool; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 93 — BP-PY-46

- Function context: `./scripts/findings/functions/93.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/evaluate_v2.py:127:13`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if r.error:
print(f"{status} {r.festival_id:20} | expected {r.expected} | ERROR: {r.error}")
else:
```

Why this is a false positive: CLI evaluation tool; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 94 — BP-PY-46

- Function context: `./scripts/findings/functions/94.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/evaluate_v2.py:129:13`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
else:
print(
f"{status} {r.festival_id:20} | expected {r.expected} | got {calc} | Δ={r.variance}d | {r.method}"
```

Why this is a false positive: CLI evaluation tool; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 95 — BP-PY-46

- Function context: `./scripts/findings/functions/95.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/evaluate_v2.py:133:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

print("\n" + "-" * 70)
print(f"RESULT: {passed}/{total} passed ({100*passed/total:.0f}%)")
```

Why this is a false positive: CLI evaluation tool; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 96 — BP-PY-46

- Function context: `./scripts/findings/functions/96.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/evaluate_v2.py:134:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
print("\n" + "-" * 70)
print(f"RESULT: {passed}/{total} passed ({100*passed/total:.0f}%)")

```

Why this is a false positive: CLI evaluation tool; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 97 — BP-PY-46

- Function context: `./scripts/findings/functions/97.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/evaluate_v2.py:137:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if passed == total:
print("🎉 ALL TESTS PASSED!")
print("-" * 70)
```

Why this is a false positive: CLI evaluation tool; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 98 — BP-PY-46

- Function context: `./scripts/findings/functions/98.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/evaluate_v2.py:138:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
print("🎉 ALL TESTS PASSED!")
print("-" * 70)

```

Why this is a false positive: CLI evaluation tool; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 99 — BP-PY-46

- Function context: `./scripts/findings/functions/99.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/evaluate_v3.py:27:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

print("=" * 70)
print("FESTIVAL DATE CALCULATOR V2 - EPHEMERIS VERIFICATION")
```

Why this is a false positive: CLI verification tool; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 100 — BP-PY-46

- Function context: `./scripts/findings/functions/100.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/evaluate_v3.py:28:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
print("=" * 70)
print("FESTIVAL DATE CALCULATOR V2 - EPHEMERIS VERIFICATION")
print("=" * 70 + "\n")
```

Why this is a false positive: CLI verification tool; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 101 — BP-PY-46

- Function context: `./scripts/findings/functions/101.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/evaluate_v3.py:29:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
print("FESTIVAL DATE CALCULATOR V2 - EPHEMERIS VERIFICATION")
print("=" * 70 + "\n")

```

Why this is a false positive: CLI verification tool; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 102 — BP-PY-46

- Function context: `./scripts/findings/functions/102.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/evaluate_v3.py:45:13`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

print(f"✅ {fid:22} -> {d} | {tithi_info:30} | {result.method}")
results.append(
```

Why this is a false positive: CLI verification tool; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 103 — BP-PY-46

- Function context: `./scripts/findings/functions/103.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/evaluate_v3.py:55:13`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
else:
print(f"❌ {fid:22} -> NOT CALCULATED")
results.append(
```

Why this is a false positive: CLI verification tool; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 104 — BP-PY-46

- Function context: `./scripts/findings/functions/104.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/evaluate_v3.py:67:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
calculated = sum(1 for r in results if r.calculated)
print("\n" + "-" * 70)
print(f"CALCULATED: {calculated}/{len(results)} festivals")
```

Why this is a false positive: CLI verification tool; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 105 — BP-PY-46

- Function context: `./scripts/findings/functions/105.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/evaluate_v3.py:68:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
print("\n" + "-" * 70)
print(f"CALCULATED: {calculated}/{len(results)} festivals")
print("-" * 70)
```

Why this is a false positive: CLI verification tool; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 106 — BP-PY-46

- Function context: `./scripts/findings/functions/106.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/evaluate_v3.py:69:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
print(f"CALCULATED: {calculated}/{len(results)} festivals")
print("-" * 70)

```

Why this is a false positive: CLI verification tool; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 107 — BP-PY-2

- Function context: `./scripts/findings/functions/107.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/generate_beta_dashboard.py:33:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
return float(metrics["uptime_percent"])
except (TypeError, ValueError):
pass
```

Why this is a false positive: missing/invalid metric values fall through to a default return; the typed except implements a documented default.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 108 — CWE-390

- Function context: `./scripts/findings/functions/108.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/generate_beta_dashboard.py:33:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
return float(metrics["uptime_percent"])
except (TypeError, ValueError):
pass
```

Why this is a false positive: the handler's action is falling back to the default metric value.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 109 — BP-PY-45

- Function context: `./scripts/findings/functions/109.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/generate_offline.py:18:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

sys.path.insert(0, str(Path(__file__).parent.parent))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 110 — BP-PY-45

- Function context: `./scripts/findings/functions/110.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/generate_snapshot.py:16:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

sys.path.insert(0, str(Path(__file__).parent.parent))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 111 — CWE-88

- Function context: `./scripts/findings/functions/111.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/ingest_moha_pdfs.py:178:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
prefix = tmp_dir / "page"
subprocess.run(
["pdftoppm", "-r", "300", "-png", str(pdf_path), str(prefix)],
```

Why this is a false positive: `pdf_path`/`prefix` are paths from a local filesystem scan in an offline ingestion tool with fixed flags; no externally influenced value can become an unintended option.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 112 — PERF-PY-27

- Function context: `./scripts/findings/functions/112.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/ingest_moha_pdfs.py:193:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if txt_path.exists():
parts.append(txt_path.read_text(encoding="utf-8", errors="ignore"))
return "\n".join(parts)
```

Why this is a false positive: each iteration reads a different `txt_path` derived from the per-image base name; the same path is not reloaded.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 113 — PERF-PY-27

- Function context: `./scripts/findings/functions/113.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/ingest_moha_pdfs.py:386:22`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
) -> int:
overrides = json.loads(overrides_path.read_text())
# Remove known bad OCR-derived entries
```

Why this is a false positive: the path is read and parsed exactly once per function call.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 114 — PERF-PY-26

- Function context: `./scripts/findings/functions/114.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/ingest_moha_pdfs.py:466:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
for pdf in pdfs:
bs_year = parse_bs_year_from_filename(pdf)
if not bs_year:
```

Why this is a false positive: per-file parsing in an offline batch ingest; each element is distinct and the parse is required.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 115 — CWE-772

- Function context: `./scripts/findings/functions/115.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/ingest_pradhanlaw_2082.py:113:43`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

def ingest(url: str, force: bool = False):
html = urllib.request.urlopen(url).read().decode("utf-8", errors="ignore")
```

Why this is a false positive: the `urlopen` response is a temporary consumed inline via `.read()` and becomes garbage immediately; no persistent handle is retained.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 116 — BP-PY-46

- Function context: `./scripts/findings/functions/116.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/ingest_pradhanlaw_2082.py:156:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
save_overrides(overrides)
print(f"✅ Ingested {added} override(s), skipped {skipped} existing.")

```

Why this is a false positive: standalone CLI ingest tool; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 117 — BP-PY-2

- Function context: `./scripts/findings/functions/117.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/load_test.py:163:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
await make_request(session, f"{API_BASE}{ep}")
except (OSError, RuntimeError, TimeoutError, ValueError):
pass
```

Why this is a false positive: a load-test harness intentionally tolerates per-endpoint request failures to keep exercising the remaining endpoints.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 118 — CWE-390

- Function context: `./scripts/findings/functions/118.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/load_test.py:163:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
await make_request(session, f"{API_BASE}{ep}")
except (OSError, RuntimeError, TimeoutError, ValueError):
pass
```

Why this is a false positive: continue-on-failure is the defined load-test behavior.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 119 — BP-PY-46

- Function context: `./scripts/findings/functions/119.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/profile_ephemeris.py:60:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
out.write_text(report, encoding="utf-8")
print(out)

```

Why this is a false positive: standalone profiling CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 120 — BP-PY-46

- Function context: `./scripts/findings/functions/120.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/profile_tithi.py:60:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
out.write_text(text, encoding="utf-8")
print(out)

```

Why this is a false positive: standalone profiling CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 121 — BP-PY-1

- Function context: `./scripts/findings/functions/121.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/week3_ground_truth_pipeline.py:174:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
converted_dates.append(g_date.isoformat())
except Exception as e:  # noqa: BLE001
err = str(e)
```

Why this is a false positive: the exception text is captured into the row's error record.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 122 — CWE-396

- Function context: `./scripts/findings/functions/122.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/week3_ground_truth_pipeline.py:174:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
converted_dates.append(g_date.isoformat())
except Exception as e:  # noqa: BLE001
err = str(e)
```

Why this is a false positive: the failure is recorded per row.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 123 — PERF-PY-25

- Function context: `./scripts/findings/functions/123.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/week3_ground_truth_pipeline.py:216:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
baseline_records.append(
BaselineRecord(
bs_year=bs_year,
```

Why this is a false positive: per-row record construction is inherent to the batch pipeline; not a hot path.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 125 — CWE-22

- Function context: `./scripts/findings/functions/125.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/benchmark/harness.py:163:73`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

out_path = Path(args.out) if args.out else Path("benchmark/results") / f"{pack['pack_id']}_report.json"
out_path.parent.mkdir(parents=True, exist_ok=True)
```

Why this is a false positive: a local benchmark CLI: `pack_id` comes from a validated pack file and `args.out` is an operator CLI argument, not request-controlled input.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 126 — CWE-73

- Function context: `./scripts/findings/functions/126.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/benchmark/validate_pack.py:83:12`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

path = Path(sys.argv[1])
if not path.exists():
```

Why this is a false positive: operator-supplied CLI argument to a local validation tool; no request boundary.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 127 — BP-PY-46

- Function context: `./scripts/findings/functions/127.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/examples/python/convert.py:15:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

print(json.dumps(payload, indent=2, sort_keys=True))

```

Why this is a false positive: example script whose purpose is printing the converted payload.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 128 — BP-PY-46

- Function context: `./scripts/findings/functions/128.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/examples/python/holidays.py:16:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

print(json.dumps(payload, indent=2, sort_keys=True))

```

Why this is a false positive: example script whose purpose is printing output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 129 — BP-PY-46

- Function context: `./scripts/findings/functions/129.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/examples/python/preflight.py:11:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
base_url = os.environ.get("PARVA_API_BASE", DEFAULT_EXAMPLE_API_BASE).rstrip("/")
print(f"[parva-example] Using API base: {base_url}", file=sys.stderr)
return base_url
```

Why this is a false positive: example CLI script; intentional diagnostic output to stderr.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 130 — BP-PY-46

- Function context: `./scripts/findings/functions/130.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/examples/python/verify_bundle.py:16:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

print(json.dumps(payload, indent=2, sort_keys=True))

```

Why this is a false positive: example script whose purpose is printing output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 131 — BP-PY-45

- Function context: `./scripts/findings/functions/131.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/integrations/mcp/server.py:11:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(PROJECT_ROOT) not in sys.path:
sys.path.insert(0, str(PROJECT_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 132 — BP-PY-45

- Function context: `./scripts/findings/functions/132.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/packages/parva-agent-tools/examples/basic_usage.py:11:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(PACKAGE_ROOT) not in sys.path:
sys.path.insert(0, str(PACKAGE_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 135 — BP-PY-13

- Function context: `./scripts/findings/functions/135.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/packages/parva-agent-tools/parva_tools/langchain.py:100:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
for key, value in list(payload.items()):
token = "{" + key + "}"
if token in bound:
```

Why this is a false positive: the 'secret-like' variable `token` is a template placeholder built from a dict key (`"{" + key + "}"`), not a credential.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 138 — BP-PY-2

- Function context: `./scripts/findings/functions/138.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/packages/parva-mcp-server/src/parva_mcp_server/client.py:192:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
)
except ValueError:
pass
```

Why this is a false positive: an unparseable `content-length` header falls through to the streaming size check; the ValueError triggers the robust path.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 139 — CWE-390

- Function context: `./scripts/findings/functions/139.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/packages/parva-mcp-server/src/parva_mcp_server/client.py:192:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
)
except ValueError:
pass
```

Why this is a false positive: the handler's action is the fall-through to streaming validation.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 140 — CWE-1071

- Function context: `./scripts/findings/functions/140.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/packages/parva-mcp-server/src/parva_mcp_server/client.py:192:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
)
except ValueError:
pass
```

Why this is a false positive: the empty handler implements the intended fallback to the streaming check.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 141 — BP-PY-1

- Function context: `./scripts/findings/functions/141.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/packages/parva-mcp-server/src/parva_mcp_server/server.py:341:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
return _tool_error_result(exc, code="MCP_ADAPTER_ERROR")
except Exception as exc:  # pragma: no cover - defensive protocol boundary
LOGGER.error("Unexpected MCP tool failure: %s", type(exc).__name__)
```

Why this is a false positive: the exception type is logged and an error result returned at a defensive protocol boundary; the failure is surfaced.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 142 — CWE-396

- Function context: `./scripts/findings/functions/142.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/packages/parva-mcp-server/src/parva_mcp_server/server.py:341:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
return _tool_error_result(exc, code="MCP_ADAPTER_ERROR")
except Exception as exc:  # pragma: no cover - defensive protocol boundary
LOGGER.error("Unexpected MCP tool failure: %s", type(exc).__name__)
```

Why this is a false positive: the handler logs and returns a structured error result.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 143 — BP-PY-2

- Function context: `./scripts/findings/functions/143.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/packages/parva-mcp-server/tests/test_client.py:69:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
validate_public_origin(origin)
except ValueError:
pass
```

Why this is a false positive: the test intentionally invokes the validator expecting it to raise for an invalid origin; the pass is the expected-rejection assertion.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 144 — CWE-390

- Function context: `./scripts/findings/functions/144.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/packages/parva-mcp-server/tests/test_client.py:69:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
validate_public_origin(origin)
except ValueError:
pass
```

Why this is a false positive: the expected exception is the test outcome.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 145 — CWE-1071

- Function context: `./scripts/findings/functions/145.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/packages/parva-mcp-server/tests/test_client.py:69:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
validate_public_origin(origin)
except ValueError:
pass
```

Why this is a false positive: the empty handler is the expected-rejection assertion.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 146 — BP-PY-46

- Function context: `./scripts/findings/functions/146.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/packages/parva-python/parva/cli.py:18:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
def _print(payload: dict[str, Any]) -> int:
print(json.dumps(payload, indent=2, sort_keys=True))
return 0
```

Why this is a false positive: `cli.py` is the SDK's command-line entry module; printing JSON is its stated purpose.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 147 — PERF-PY-26

- Function context: `./scripts/findings/functions/147.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/packages/parva-python/parva/cli.py:199:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
parser = build_parser()
args = parser.parse_args(argv)
return args.func(args)
```

Why this is a false positive: argparse runs once at CLI startup.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 148 — CWE-829

- Function context: `./scripts/findings/functions/148.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/parva_mcp_server/__init__.py:30:14`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
raise AttributeError(name)
module = __import__(f"{__name__}.{module_name}", fromlist=[name])
return getattr(module, name)
```

Why this is a false positive: package `__getattr__` lazy-import shim; `module_name` is derived from a Python attribute name, not an untrusted control sphere.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 149 — CWE-94

- Function context: `./scripts/findings/functions/149.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/parva_mcp_server/__init__.py:30:14`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
raise AttributeError(name)
module = __import__(f"{__name__}.{module_name}", fromlist=[name])
return getattr(module, name)
```

Why this is a false positive: the dynamic import is a lazy-attribute shim whose argument is an attribute name; no externally influenced text reaches it.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 150 — BP-PY-45

- Function context: `./scripts/findings/functions/150.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/public-benchmark/runners/run_against_parva.py:19:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(ROOT) not in sys.path:
sys.path.insert(0, str(ROOT))
PROJECT_ROOT = ROOT.parent
```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 151 — BP-PY-45

- Function context: `./scripts/findings/functions/151.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/public-benchmark/runners/run_against_parva.py:23:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 152 — BP-PY-1

- Function context: `./scripts/findings/functions/152.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/public-benchmark/runners/run_against_parva.py:316:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
result = _evaluate_task(task, status, payload)
except Exception as exc:  # noqa: BLE001 - benchmark output should record every task.
result = {
```

Why this is a false positive: the exception is recorded per task in the benchmark result row (`"error": str(exc)`), which the comment explicitly requires.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 153 — CWE-396

- Function context: `./scripts/findings/functions/153.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/public-benchmark/runners/run_against_parva.py:316:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
result = _evaluate_task(task, status, payload)
except Exception as exc:  # noqa: BLE001 - benchmark output should record every task.
result = {
```

Why this is a false positive: the failure is recorded per task.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 154 — BP-PY-45

- Function context: `./scripts/findings/functions/154.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/public-benchmark/runners/run_against_static_baseline.py:12:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(ROOT) not in sys.path:
sys.path.insert(0, str(ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 155 — BP-PY-45

- Function context: `./scripts/findings/functions/155.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/backtest_future_bs_model.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 156 — CWE-88

- Function context: `./scripts/findings/functions/156.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/benchmark/generate_benchmark_badge.py:46:12`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
def _source_timestamp() -> str:
proc = subprocess.run(
[
```

Why this is a false positive: the only dynamic operand is a repo-relative constant path passed after `--`; no externally influenced input.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 157 — BP-PY-45

- Function context: `./scripts/findings/functions/157.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/calibrate_future_bs_rules.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 158 — BP-PY-45

- Function context: `./scripts/findings/functions/158.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/check_future_bs_public_leakage.py:20:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 159 — BP-PY-45

- Function context: `./scripts/findings/functions/159.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/check_maturity_lanes.py:17:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 160 — BP-PY-46

- Function context: `./scripts/findings/functions/160.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/check_path_leaks.py:80:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
try:
print(text)
except UnicodeEncodeError:
```

Why this is a false positive: standalone checker script; the helper's purpose is printing lines.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 161 — BP-PY-45

- Function context: `./scripts/findings/functions/161.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/claims/compile_public_claims.py:11:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(PROJECT_ROOT) not in sys.path:
sys.path.insert(0, str(PROJECT_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 162 — BP-PY-45

- Function context: `./scripts/findings/functions/162.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/compare_external_sheet.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 163 — BP-PY-45

- Function context: `./scripts/findings/functions/163.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/conformance/generate_conformance_report.py:15:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(SCRIPT_ROOT) not in sys.path:
sys.path.insert(0, str(SCRIPT_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 164 — BP-PY-2

- Function context: `./scripts/findings/functions/164.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/download_jpl_kernel.py:51:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
)
except (OSError, SpkValidationError):
pass
```

Why this is a false positive: the pass is inside a download/verify loop that retries; a failed verify deliberately falls through to the download path.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 165 — CWE-390

- Function context: `./scripts/findings/functions/165.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/download_jpl_kernel.py:51:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
)
except (OSError, SpkValidationError):
pass
```

Why this is a false positive: the handler's action is proceeding to the download/retry path.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 166 — BP-PY-46

- Function context: `./scripts/findings/functions/166.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/download_jpl_kernel.py:55:17`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if not quiet:
print(f"JPL kernel already present and verified: {output}")
return output
```

Why this is a false positive: standalone downloader script; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 167 — BP-PY-46

- Function context: `./scripts/findings/functions/167.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/download_jpl_kernel.py:62:13`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if not quiet:
print(f"Downloading {url} -> {output}")
with urllib.request.urlopen(url, timeout=120) as response, tmp_path.open("wb") as fh:
```

Why this is a false positive: standalone downloader script; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 168 — BP-PY-46

- Function context: `./scripts/findings/functions/168.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/download_jpl_kernel.py:88:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if not quiet:
print(f"Verified JPL kernel {output} sha256={expected_sha256}")
return output
```

Why this is a false positive: standalone downloader script; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 169 — BP-PY-45

- Function context: `./scripts/findings/functions/169.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/ephemeris/generate_solar_ingress_differential.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(PROJECT_ROOT / "backend") not in sys.path:
sys.path.insert(0, str(PROJECT_ROOT / "backend"))
if str(PROJECT_ROOT) not in sys.path:
```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 170 — BP-PY-45

- Function context: `./scripts/findings/functions/170.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/ephemeris/generate_solar_ingress_differential.py:16:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(PROJECT_ROOT) not in sys.path:
sys.path.insert(0, str(PROJECT_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 171 — BP-PY-45

- Function context: `./scripts/findings/functions/171.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/evidence/generate_evidence_packet.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND) not in sys.path:
sys.path.insert(0, str(BACKEND))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 172 — BP-PY-45

- Function context: `./scripts/findings/functions/172.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/evidence/ingest_source_record.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND) not in sys.path:
sys.path.insert(0, str(BACKEND))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 173 — BP-PY-45

- Function context: `./scripts/findings/functions/173.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/evidence/normalize_source_rows.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND) not in sys.path:
sys.path.insert(0, str(BACKEND))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 174 — BP-PY-45

- Function context: `./scripts/findings/functions/174.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/evidence/promote_evidence_to_benchmark_candidate.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND) not in sys.path:
sys.path.insert(0, str(BACKEND))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 175 — BP-PY-45

- Function context: `./scripts/findings/functions/175.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/export_future_bs_predictions.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 176 — BP-PY-45

- Function context: `./scripts/findings/functions/176.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/forge/build_bitplanes.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 177 — BP-PY-45

- Function context: `./scripts/findings/functions/177.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/forge/build_manifest.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 178 — BP-PY-45

- Function context: `./scripts/findings/functions/178.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/forge/build_static_index.py:12:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 179 — BP-PY-45

- Function context: `./scripts/findings/functions/179.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/forge/verify_manifest.py:12:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 180 — BP-PY-45

- Function context: `./scripts/findings/functions/180.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/accuracy_lab.py:16:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 181 — CWE-88

- Function context: `./scripts/findings/functions/181.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/accuracy_lab.py:28:21`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
python_executable = os.getenv("PARVA_PYTHON", sys.executable)
completed = subprocess.run([python_executable, *sys.argv], env=env)
raise SystemExit(completed.returncode)
```

Why this is a false positive: the script re-executes itself with its own already-parsed `sys.argv` under a re-exec guard; no new command is constructed from untrusted input.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 182 — BP-PY-45

- Function context: `./scripts/findings/functions/182.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/audit_external_bs_sheet.py:16:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(path) not in sys.path:
sys.path.insert(0, str(path))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 183 — BP-PY-45

- Function context: `./scripts/findings/functions/183.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/audit_witness_corpus.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 184 — BP-PY-45

- Function context: `./scripts/findings/functions/184.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/build_source_agreement_graph.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 185 — BP-PY-45

- Function context: `./scripts/findings/functions/185.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/check_data_target.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 186 — BP-PY-45

- Function context: `./scripts/findings/functions/186.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/collect_witness_corpus.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 187 — BP-PY-45

- Function context: `./scripts/findings/functions/187.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/compare_solar_civil_before_after.py:18:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 188 — BP-PY-45

- Function context: `./scripts/findings/functions/188.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/generate_30_year_past_corpus_report.py:15:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 189 — BP-PY-45

- Function context: `./scripts/findings/functions/189.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/generate_all_final_artifacts.py:18:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(path) not in sys.path:
sys.path.insert(0, str(path))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 190 — CWE-88

- Function context: `./scripts/findings/functions/190.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/generate_all_final_artifacts.py:30:21`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
python_executable = os.getenv("PARVA_PYTHON", sys.executable)
completed = subprocess.run([python_executable, *sys.argv], env=env)
raise SystemExit(completed.returncode)
```

Why this is a false positive: self re-execution with the process's own argv under a guard (same pattern as 181).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 191 — BP-PY-45

- Function context: `./scripts/findings/functions/191.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/generate_calendar_var_report.py:16:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(path) not in sys.path:
sys.path.insert(0, str(path))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 192 — BP-PY-45

- Function context: `./scripts/findings/functions/192.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/generate_claim_readiness_report.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 193 — BP-PY-45

- Function context: `./scripts/findings/functions/193.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/generate_human_review_promotion_plan.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 194 — BP-PY-45

- Function context: `./scripts/findings/functions/194.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/generate_human_review_queue.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 195 — BP-PY-45

- Function context: `./scripts/findings/functions/195.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/generate_residual_report.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 196 — BP-PY-45

- Function context: `./scripts/findings/functions/196.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/merge_high_trust_witnesses.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 197 — BP-PY-45

- Function context: `./scripts/findings/functions/197.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/optimize_regime_aware_accuracy_loop.py:17:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 198 — BP-PY-45

- Function context: `./scripts/findings/functions/198.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/optimize_solar_civil_rules_loop.py:23:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 199 — PERF-PY-23

- Function context: `./scripts/findings/functions/199.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/optimize_solar_civil_rules_loop.py:281:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
"risk_flags": ";".join(risk_cache.get(year, [])),
"applied_rules": json.dumps(applied_for_month, ensure_ascii=False, sort_keys=True),
}
```

Why this is a false positive: per-month serialization in an offline rule-optimization research script; the serialization is the output being produced.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 200 — BP-PY-45

- Function context: `./scripts/findings/functions/200.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/parse_archive_panchanga.py:17:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 201 — BP-PY-45

- Function context: `./scripts/findings/functions/201.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/parse_gorkhapatra_witnesses.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 202 — BP-PY-45

- Function context: `./scripts/findings/functions/202.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/parse_open_source_converter_tables.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 203 — BP-PY-45

- Function context: `./scripts/findings/functions/203.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/parse_rat32_calendar.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 204 — BP-PY-45

- Function context: `./scripts/findings/functions/204.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/reconstruct_month_starts.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 205 — BP-PY-45

- Function context: `./scripts/findings/functions/205.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/replay_2083_ashwin.py:16:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 206 — CWE-88

- Function context: `./scripts/findings/functions/206.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/replay_2083_ashwin.py:28:21`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
python_executable = os.getenv("PARVA_PYTHON", sys.executable)
completed = subprocess.run([python_executable, *sys.argv], env=env)
raise SystemExit(completed.returncode)
```

Why this is a false positive: self re-execution with the process's own argv under a guard (same pattern as 181).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 207 — BP-PY-45

- Function context: `./scripts/findings/functions/207.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/research_and_collect_high_trust_sources.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 208 — BP-PY-45

- Function context: `./scripts/findings/functions/208.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/run_data_acquisition_loop.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 209 — BP-PY-45

- Function context: `./scripts/findings/functions/209.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/run_hamropatro_shadow_evaluation.py:12:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND) not in sys.path:
sys.path.insert(0, str(BACKEND))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 210 — BP-PY-45

- Function context: `./scripts/findings/functions/210.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/run_model_search.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 211 — BP-PY-45

- Function context: `./scripts/findings/functions/211.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/run_month_start_inversion_workbench.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 212 — BP-PY-45

- Function context: `./scripts/findings/functions/212.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/run_time_travel_backtest.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 213 — BP-PY-45

- Function context: `./scripts/findings/functions/213.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/train_solar_civil_with_reconstructed.py:17:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 214 — BP-PY-45

- Function context: `./scripts/findings/functions/214.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/future_bs/tune_risk_thresholds.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 215 — BP-PY-45

- Function context: `./scripts/findings/functions/215.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/generate_accuracy_report.py:29:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
"""Compare V2 calculator output against ground truth for all years."""
sys.path.insert(0, str(PROJECT_ROOT / "backend"))

```

Why this is a false positive: the script needs the in-tree `backend` package importable before running; module-top bootstrap.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 216 — BP-PY-45

- Function context: `./scripts/findings/functions/216.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/generate_authority_dashboard.py:16:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 217 — BP-PY-45

- Function context: `./scripts/findings/functions/217.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/generate_residual_report.py:15:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 218 — CWE-88

- Function context: `./scripts/findings/functions/218.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/node_runtime.py:49:14`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

result = subprocess.run(
[resolved, "--version"],
```

Why this is a false positive: `resolved` is the node executable from `shutil.which` and the flag is a fixed constant; no externally influenced argument.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 219 — BP-PY-45

- Function context: `./scripts/findings/functions/219.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_agent_benchmark.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 220 — BP-PY-45

- Function context: `./scripts/findings/functions/220.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_agent_verify.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 221 — BP-PY-45

- Function context: `./scripts/findings/functions/221.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_conformance.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 222 — BP-PY-45

- Function context: `./scripts/findings/functions/222.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_credential_issue.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 223 — BP-PY-45

- Function context: `./scripts/findings/functions/223.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_credential_verify.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 224 — CWE-73

- Function context: `./scripts/findings/functions/224.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_credential_verify.py:22:29`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
return 0 if len(sys.argv) == 2 and sys.argv[1] in {"-h", "--help"} else 2
credential = json.loads(Path(sys.argv[1]).read_text(encoding="utf-8"))
result = verify_calendar_credential_payload(credential)
```

Why this is a false positive: operator-supplied CLI argument to a local verification tool; no request boundary.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 225 — BP-PY-45

- Function context: `./scripts/findings/functions/225.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_evidence_packet.py:15:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(PROJECT_ROOT) not in sys.path:
sys.path.insert(0, str(PROJECT_ROOT))
if str(BACKEND_ROOT) not in sys.path:
```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 226 — BP-PY-45

- Function context: `./scripts/findings/functions/226.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_evidence_packet.py:17:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 227 — BP-PY-45

- Function context: `./scripts/findings/functions/227.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_impact_verify.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 228 — BP-PY-45

- Function context: `./scripts/findings/functions/228.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_offline_bundle.py:15:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 229 — CWE-73

- Function context: `./scripts/findings/functions/229.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_offline_verify.py:20:12`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
return 2
root = Path(sys.argv[1])
manifest_path = root / "bundle-manifest.json"
```

Why this is a false positive: operator-supplied CLI argument to a local verification tool; no request boundary.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 230 — BP-PY-45

- Function context: `./scripts/findings/functions/230.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_protocol_verify.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 231 — BP-PY-45

- Function context: `./scripts/findings/functions/231.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_release_diff.py:15:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(PROJECT_ROOT) not in sys.path:
sys.path.insert(0, str(PROJECT_ROOT))
if str(BACKEND_ROOT) not in sys.path:
```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 232 — BP-PY-45

- Function context: `./scripts/findings/functions/232.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_release_diff.py:17:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 233 — BP-PY-45

- Function context: `./scripts/findings/functions/233.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_rulelang_verify.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(PROJECT_ROOT) not in sys.path:
sys.path.insert(0, str(PROJECT_ROOT))
if str(BACKEND_ROOT) not in sys.path:
```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 234 — BP-PY-45

- Function context: `./scripts/findings/functions/234.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_rulelang_verify.py:16:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 235 — BP-PY-1

- Function context: `./scripts/findings/functions/235.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_rulelang_verify.py:169:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
_assert_disputed_fact_blocks()
except Exception as exc:  # noqa: BLE001
print(f"Project Parva RuleLang verification failed: {exc}", file=sys.stderr)
```

Why this is a false positive: the handler records the exception (log message with `exc`, result/error field with `str(exc)`, warning text, or stderr output) and/or re-raises, so the failure is surfaced — it does not pass or continue without recording the exception type.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 236 — CWE-396

- Function context: `./scripts/findings/functions/236.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_rulelang_verify.py:169:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
_assert_disputed_fact_blocks()
except Exception as exc:  # noqa: BLE001
print(f"Project Parva RuleLang verification failed: {exc}", file=sys.stderr)
```

Why this is a false positive: the handler deliberately handles the failure: it records the exception in logs/results/warnings or re-raises a specific error; nothing is hidden.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 237 — BP-PY-45

- Function context: `./scripts/findings/functions/237.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_timegraph_verify.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(PROJECT_ROOT) not in sys.path:
sys.path.insert(0, str(PROJECT_ROOT))
if str(BACKEND_ROOT) not in sys.path:
```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 238 — BP-PY-45

- Function context: `./scripts/findings/functions/238.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_timegraph_verify.py:16:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 239 — BP-PY-1

- Function context: `./scripts/findings/functions/239.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_timegraph_verify.py:41:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
raise TimeGraphError("sample date query returned no facts")
except Exception as exc:  # noqa: BLE001
print(f"Project Parva TimeGraph verification failed: {exc}", file=sys.stderr)
```

Why this is a false positive: the handler records the exception (log message with `exc`, result/error field with `str(exc)`, warning text, or stderr output) and/or re-raises, so the failure is surfaced — it does not pass or continue without recording the exception type.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 240 — CWE-396

- Function context: `./scripts/findings/functions/240.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_timegraph_verify.py:41:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
raise TimeGraphError("sample date query returned no facts")
except Exception as exc:  # noqa: BLE001
print(f"Project Parva TimeGraph verification failed: {exc}", file=sys.stderr)
```

Why this is a false positive: the handler deliberately handles the failure: it records the exception in logs/results/warnings or re-raises a specific error; nothing is hidden.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 241 — BP-PY-45

- Function context: `./scripts/findings/functions/241.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_trust_verify.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(PROJECT_ROOT) not in sys.path:
sys.path.insert(0, str(PROJECT_ROOT))
if str(BACKEND_ROOT) not in sys.path:
```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 242 — BP-PY-45

- Function context: `./scripts/findings/functions/242.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_trust_verify.py:16:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 243 — BP-PY-1

- Function context: `./scripts/findings/functions/243.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_trust_verify.py:75:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
raise TrustInfrastructureError("; ".join(trust["issues"]))
except Exception as exc:  # noqa: BLE001
print(f"Project Parva trust verification failed: {exc}", file=sys.stderr)
```

Why this is a false positive: the handler records the exception (log message with `exc`, result/error field with `str(exc)`, warning text, or stderr output) and/or re-raises, so the failure is surfaced — it does not pass or continue without recording the exception type.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 244 — CWE-396

- Function context: `./scripts/findings/functions/244.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/parva_trust_verify.py:75:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
raise TrustInfrastructureError("; ".join(trust["issues"]))
except Exception as exc:  # noqa: BLE001
print(f"Project Parva trust verification failed: {exc}", file=sys.stderr)
```

Why this is a false positive: the handler deliberately handles the failure: it records the exception in logs/results/warnings or re-raises a specific error; nothing is hidden.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 245 — BP-PY-45

- Function context: `./scripts/findings/functions/245.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/perf/route_latency_smoke.py:21:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 246 — BP-PY-45

- Function context: `./scripts/findings/functions/246.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/precompute/loadtest_cache.py:19:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 247 — BP-PY-46

- Function context: `./scripts/findings/functions/247.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/precompute/precompute_all.py:15:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
def _run(cmd: list[str]) -> None:
print(f"$ {' '.join(cmd)}")
subprocess.run(cmd, cwd=PROJECT_ROOT, check=True)
```

Why this is a false positive: standalone precompute script; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 248 — BP-PY-45

- Function context: `./scripts/findings/functions/248.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/precompute/precompute_festivals.py:15:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 249 — BP-PY-45

- Function context: `./scripts/findings/functions/249.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/precompute/precompute_panchanga.py:15:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 250 — BP-PY-45

- Function context: `./scripts/findings/functions/250.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/precompute/profile_endpoints.py:17:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 251 — BP-PY-45

- Function context: `./scripts/findings/functions/251.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/precompute_future_bs_predictions.py:15:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 252 — BP-PY-45

- Function context: `./scripts/findings/functions/252.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/precompute_solar_ingress_events.py:16:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 253 — BP-PY-46

- Function context: `./scripts/findings/functions/253.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/prepare_temple_data.py:23:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if not filepath.exists():
print(f"Warning: {filepath} not found")
return {"features": []}
```

Why this is a false positive: standalone data-prep script; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 254 — BP-PY-45

- Function context: `./scripts/findings/functions/254.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/audit_ceiling_depth.py:15:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 255 — BP-PY-1

- Function context: `./scripts/findings/functions/255.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/audit_ceiling_depth.py:42:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
detail = fn()
except Exception as exc:  # noqa: BLE001 - audit scripts report exact failure text.
return AuditCheck(name, "fail", f"{type(exc).__name__}: {exc}")
```

Why this is a false positive: the handler records the exception (log message with `exc`, result/error field with `str(exc)`, warning text, or stderr output) and/or re-raises, so the failure is surfaced — it does not pass or continue without recording the exception type.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 256 — CWE-396

- Function context: `./scripts/findings/functions/256.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/audit_ceiling_depth.py:42:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
detail = fn()
except Exception as exc:  # noqa: BLE001 - audit scripts report exact failure text.
return AuditCheck(name, "fail", f"{type(exc).__name__}: {exc}")
```

Why this is a false positive: the handler deliberately handles the failure: it records the exception in logs/results/warnings or re-raises a specific error; nothing is hidden.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 257 — BP-PY-45

- Function context: `./scripts/findings/functions/257.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/check_backend_smoke.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 258 — BP-PY-45

- Function context: `./scripts/findings/functions/258.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/check_ceiling_depth_semantics.py:15:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 259 — BP-PY-1

- Function context: `./scripts/findings/functions/259.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/check_ceiling_depth_semantics.py:219:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
failures.extend(_verify_replay_artifacts())
except Exception as exc:  # noqa: BLE001
failures.append(f"proof artifact evidence verification failed: {exc}")
```

Why this is a false positive: the handler records the exception (log message with `exc`, result/error field with `str(exc)`, warning text, or stderr output) and/or re-raises, so the failure is surfaced — it does not pass or continue without recording the exception type.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 260 — CWE-396

- Function context: `./scripts/findings/functions/260.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/check_ceiling_depth_semantics.py:219:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
failures.extend(_verify_replay_artifacts())
except Exception as exc:  # noqa: BLE001
failures.append(f"proof artifact evidence verification failed: {exc}")
```

Why this is a false positive: the handler deliberately handles the failure: it records the exception in logs/results/warnings or re-raises a specific error; nothing is hidden.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 261 — BP-PY-45

- Function context: `./scripts/findings/functions/261.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/check_ceiling_phase_requirements.py:18:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 262 — BP-PY-45

- Function context: `./scripts/findings/functions/262.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/check_contract_freeze.py:24:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 263 — BP-PY-46

- Function context: `./scripts/findings/functions/263.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/check_contract_freeze.py:97:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if not snapshot_path.exists():
print(f"[{track}] Missing snapshot: {snapshot_path}")
return 2
```

Why this is a false positive: release-check script; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 264 — BP-PY-46

- Function context: `./scripts/findings/functions/264.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/check_contract_freeze.py:105:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if frozen == current:
print(f"[{track}] Contract freeze check passed.")
return 0
```

Why this is a false positive: release-check script; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 265 — BP-PY-46

- Function context: `./scripts/findings/functions/265.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/check_contract_freeze.py:108:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

print(f"[{track}] Contract freeze check failed: schema differs from snapshot.")
return 1
```

Why this is a false positive: release-check script; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 266 — BP-PY-45

- Function context: `./scripts/findings/functions/266.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/check_documented_routes.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 267 — BP-PY-45

- Function context: `./scripts/findings/functions/267.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/check_jpl_lane.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(PROJECT_ROOT / "backend") not in sys.path:
sys.path.insert(0, str(PROJECT_ROOT / "backend"))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 268 — BP-PY-45

- Function context: `./scripts/findings/functions/268.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/check_local_kernel_package.py:16:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(SCRIPTS_ROOT) not in sys.path:
sys.path.insert(0, str(SCRIPTS_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 269 — BP-PY-45

- Function context: `./scripts/findings/functions/269.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/check_mcp_registry_metadata.py:11:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
MCP_SRC = ROOT / "packages" / "parva-mcp-server" / "src"
sys.path.insert(0, str(MCP_SRC))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 270 — CWE-88

- Function context: `./scripts/findings/functions/270.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/check_package_readiness.py:39:14`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
def _git_files(package_dir: Path) -> list[str]:
result = subprocess.run(
["git", "ls-files", "--", str(package_dir.relative_to(ROOT)).replace("\\", "/")],
```

Why this is a false positive: `package_dir` is derived from the repo root constant and passed after `--`; no externally influenced input.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 271 — BP-PY-45

- Function context: `./scripts/findings/functions/271.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/check_production_preflight.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 272 — BP-PY-45

- Function context: `./scripts/findings/functions/272.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/check_provenance_readiness.py:16:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_DIR) not in sys.path:
sys.path.insert(0, str(BACKEND_DIR))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 274 — CWE-88

- Function context: `./scripts/findings/functions/274.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/check_public_openapi_drift.py:37:18`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
env.setdefault("PARVA_RATE_LIMIT_ENABLED", "false")
result = subprocess.run(
[sys.executable, "scripts/release/generate_public_demo_openapi.py"],
```

Why this is a false positive: fixed argument vector of `sys.executable` plus a constant script path.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 275 — BP-PY-45

- Function context: `./scripts/findings/functions/275.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/check_public_safety_gate.py:17:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(PROJECT_ROOT) not in sys.path:
sys.path.insert(0, str(PROJECT_ROOT))
if str(BACKEND_ROOT) not in sys.path:
```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 276 — BP-PY-45

- Function context: `./scripts/findings/functions/276.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/check_public_safety_gate.py:19:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 277 — CWE-88

- Function context: `./scripts/findings/functions/277.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/check_public_safety_gate.py:143:14`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
def _check_protocol_schema_validation() -> list[str]:
result = subprocess.run(
[sys.executable, "tools/validate_schemas.py"],
```

Why this is a false positive: fixed argument vector of `sys.executable` plus a constant script path.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 278 — BP-PY-45

- Function context: `./scripts/findings/functions/278.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/check_route_inventory.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 279 — CWE-88

- Function context: `./scripts/findings/functions/279.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/check_sdk_examples.py:30:18`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
for rel in PYTHON_EXAMPLES:
result = subprocess.run([sys.executable, rel], cwd=PROJECT_ROOT, text=True, capture_output=True, check=False)
if result.returncode != 0:
```

Why this is a false positive: `rel` iterates the module-constant `PYTHON_EXAMPLES` list; no externally influenced input.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 280 — BP-PY-45

- Function context: `./scripts/findings/functions/280.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/check_sdk_install.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(SDK_ROOT) not in sys.path:
sys.path.insert(0, str(SDK_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 281 — CWE-88

- Function context: `./scripts/findings/functions/281.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/clean_local_artifacts.py:36:14`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
def _git_lines(*args: str) -> set[str]:
result = subprocess.run(
["git", *args],
```

Why this is a false positive: `args` come from internal callers with fixed repository paths.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 282 — BP-PY-46

- Function context: `./scripts/findings/functions/282.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/clean_local_artifacts.py:66:13`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if not _is_safe_target(target, tracked):
print(f"SKIP unsafe or tracked path: {relative}")
skipped += 1
```

Why this is a false positive: release-cleanup CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 283 — BP-PY-46

- Function context: `./scripts/findings/functions/283.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/clean_local_artifacts.py:70:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
action = "REMOVE" if apply else "DRY-RUN"
print(f"{action} {relative}")
if apply:
```

Why this is a false positive: release-cleanup CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 284 — BP-PY-46

- Function context: `./scripts/findings/functions/284.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/clean_local_artifacts.py:77:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
removed += 1
print(f"Summary: candidates={removed}, skipped={skipped}, apply={apply}")
return 1 if skipped else 0
```

Why this is a false positive: release-cleanup CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 285 — BP-PY-45

- Function context: `./scripts/findings/functions/285.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/generate_dashboard_metrics.py:19:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 286 — BP-PY-45

- Function context: `./scripts/findings/functions/286.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/generate_openapi_profiles.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(Path(__file__).resolve().parent) not in sys.path:
sys.path.insert(0, str(Path(__file__).resolve().parent))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 287 — CWE-215

- Function context: `./scripts/findings/functions/287.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/generate_partner_api_key.py:83:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
print(f"Scopes: {', '.join(record.scopes)}")
print(f"Secret: {record.secret}")
print("")
```

Why this is a false positive: printing the generated secret is the provisioning CLI's primary output (the operator must capture the one-time secret), not debug output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 288 — BP-PY-45

- Function context: `./scripts/findings/functions/288.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/generate_phase_00_trust_arrest_report.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 289 — BP-PY-46

- Function context: `./scripts/findings/functions/289.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/generate_public_beta_artifacts.py:15:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
def _run_step(label: str, args: list[str]) -> None:
print(f"[artifacts] {label}")
completed = subprocess.run(args, cwd=PROJECT_ROOT, check=False)
```

Why this is a false positive: artifact-generation CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 290 — BP-PY-45

- Function context: `./scripts/findings/functions/290.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/generate_public_beta_dossier.py:15:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 291 — BP-PY-45

- Function context: `./scripts/findings/functions/291.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/generate_public_demo_openapi.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(Path(__file__).resolve().parent) not in sys.path:
sys.path.insert(0, str(Path(__file__).resolve().parent))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 292 — BP-PY-46

- Function context: `./scripts/findings/functions/292.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/generate_public_demo_openapi.py:44:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
rendered_path = output_path
print(f"Wrote {rendered_path} with {len(schema.get('paths', {}))} paths.")
return 0
```

Why this is a false positive: OpenAPI generation CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 293 — BP-PY-45

- Function context: `./scripts/findings/functions/293.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/generate_release_candidate_dossier.py:18:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(SCRIPTS_DIR) not in sys.path:
sys.path.insert(0, str(SCRIPTS_DIR))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 294 — CWE-88

- Function context: `./scripts/findings/functions/294.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/generate_release_candidate_dossier.py:74:21`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
try:
completed = subprocess.run(
["git", *args],
```

Why this is a false positive: `args` come from internal callers with fixed git subcommands.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 295 — BP-PY-45

- Function context: `./scripts/findings/functions/295.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/generate_trust_status_report.py:19:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(PROJECT_ROOT) not in sys.path:
sys.path.insert(0, str(PROJECT_ROOT))
if str(BACKEND_ROOT) not in sys.path:
```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 296 — BP-PY-45

- Function context: `./scripts/findings/functions/296.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/generate_trust_status_report.py:21:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 297 — BP-PY-1

- Function context: `./scripts/findings/functions/297.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/generate_trust_status_report.py:92:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
validate_schema_file(path)
except Exception as exc:  # noqa: BLE001
failures.append(f"{path.relative_to(PROJECT_ROOT)}: {exc}")
```

Why this is a false positive: the handler records the exception (log message with `exc`, result/error field with `str(exc)`, warning text, or stderr output) and/or re-raises, so the failure is surfaced — it does not pass or continue without recording the exception type.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 298 — CWE-396

- Function context: `./scripts/findings/functions/298.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/generate_trust_status_report.py:92:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
validate_schema_file(path)
except Exception as exc:  # noqa: BLE001
failures.append(f"{path.relative_to(PROJECT_ROOT)}: {exc}")
```

Why this is a false positive: the handler deliberately handles the failure: it records the exception in logs/results/warnings or re-raises a specific error; nothing is hidden.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 299 — BP-PY-45

- Function context: `./scripts/findings/functions/299.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/refresh_trust_artifacts.py:18:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 300 — CWE-88

- Function context: `./scripts/findings/functions/300.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/refresh_trust_artifacts.py:24:17`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
def _run(script: Path, *, optional: bool = False, reason: str | None = None) -> None:
completed = subprocess.run(
[sys.executable, str(script)],
```

Why this is a false positive: `script` comes from internal constant paths; the vector is `sys.executable` + path.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 301 — BP-PY-46

- Function context: `./scripts/findings/functions/301.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/refresh_trust_artifacts.py:33:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
message = reason or f"Skipping optional step {script.name}"
print(f"{message}: {completed.stderr.strip() or completed.stdout.strip()}")
return
```

Why this is a false positive: trust-artifact refresh CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 302 — BP-PY-46

- Function context: `./scripts/findings/functions/302.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/refresh_trust_artifacts.py:36:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if completed.stdout.strip():
print(completed.stdout.strip())

```

Why this is a false positive: trust-artifact refresh CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 303 — BP-PY-45

- Function context: `./scripts/findings/functions/303.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/regenerate_public_release_hashes.py:16:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(PROJECT_ROOT) not in sys.path:
sys.path.insert(0, str(PROJECT_ROOT))
if str(BACKEND_ROOT) not in sys.path:
```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 304 — BP-PY-45

- Function context: `./scripts/findings/functions/304.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/regenerate_public_release_hashes.py:18:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 305 — CWE-208

- Function context: `./scripts/findings/functions/305.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/regenerate_public_release_hashes.py:113:20`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
next_signature = expected_signature(manifest_path, signature_path)
signature_ok = current_signature == next_signature

```

Why this is a false positive: offline release-verification script comparing JSON dicts of release signatures; no online timing attack surface and no secret equality comparison.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 306 — BP-PY-45

- Function context: `./scripts/findings/functions/306.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/reviewer_dry_run.py:17:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(SCRIPTS_ROOT) not in sys.path:
sys.path.insert(0, str(SCRIPTS_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 307 — BP-PY-45

- Function context: `./scripts/findings/functions/307.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/run_ceiling_climax_demos.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 308 — CWE-208

- Function context: `./scripts/findings/functions/308.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/run_ceiling_climax_demos.py:61:61`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
else None,
flags=frozenset({TaintFlag.REVIEW_REQUIRED}) if authority == AuthorityTaint.STATIC_REFERENCE else frozenset(),
)
```

Why this is a false positive: enum member comparison (`AuthorityTaint.STATIC_REFERENCE`), not a secret comparison.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 309 — BP-PY-46

- Function context: `./scripts/findings/functions/309.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/run_public_coverage.py:28:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
raise SystemExit("Coverage JSON is missing totals.percent_covered")
print(
json.dumps(
```

Why this is a false positive: coverage-report CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 310 — BP-PY-45

- Function context: `./scripts/findings/functions/310.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/run_release_candidate_gates.py:15:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 311 — BP-PY-46

- Function context: `./scripts/findings/functions/311.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/run_release_candidate_gates.py:33:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
def _run_step(index: int, total: int, label: str, command: list[str], env: dict[str, str] | None = None) -> None:
print(f"[RC {index}/{total}] {label}")
resolved_command = list(command)
```

Why this is a false positive: release-candidate gates CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 312 — BP-PY-46

- Function context: `./scripts/findings/functions/312.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/run_release_candidate_gates.py:61:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
)
print(
"[RC] Using provisioned local provenance signing profile "
```

Why this is a false positive: release-candidate gates CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 313 — BP-PY-45

- Function context: `./scripts/findings/functions/313.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/snapshot_openapi.py:16:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 314 — CWE-88

- Function context: `./scripts/findings/functions/314.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/verify_clean_clone_assumptions.py:87:14`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
def _is_ignored(path: str) -> bool:
result = subprocess.run(
["git", "check-ignore", "-q", path],
```

Why this is a false positive: `path` is a repo-relative path from internal iteration passed as an operand; no externally influenced input.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 316 — BP-PY-45

- Function context: `./scripts/findings/functions/316.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/verify_public.py:16:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(SCRIPTS_ROOT) not in sys.path:
sys.path.insert(0, str(SCRIPTS_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 317 — CWE-88

- Function context: `./scripts/findings/functions/317.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/verify_public.py:22:14`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
def _python_version(command: list[str]) -> tuple[int, int] | None:
result = subprocess.run(
[
```

Why this is a false positive: `command` is an internal python invocation list and the `-c` code string is fixed; no externally influenced input.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 318 — BP-PY-46

- Function context: `./scripts/findings/functions/318.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/verify_public.py:96:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
def _run(label: str, command: list[str], env: dict[str, str]) -> bool:
print(f"\n[verify-public] {label}")
print("[verify-public] " + " ".join(command))
```

Why this is a false positive: release verification CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 319 — BP-PY-46

- Function context: `./scripts/findings/functions/319.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/verify_public.py:97:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
print(f"\n[verify-public] {label}")
print("[verify-public] " + " ".join(command))
result = subprocess.run(command, cwd=PROJECT_ROOT, env=env, check=False)
```

Why this is a false positive: release verification CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 320 — BP-PY-46

- Function context: `./scripts/findings/functions/320.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/verify_public.py:100:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if result.returncode == 0:
print(f"[verify-public] PASS: {label}")
return True
```

Why this is a false positive: release verification CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 321 — BP-PY-46

- Function context: `./scripts/findings/functions/321.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/verify_public.py:102:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
return True
print(f"[verify-public] FAIL: {label} exited {result.returncode}")
return False
```

Why this is a false positive: release verification CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 322 — BP-PY-45

- Function context: `./scripts/findings/functions/322.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/resolve_npm_command.py:12:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(SCRIPTS_ROOT) not in sys.path:
sys.path.insert(0, str(SCRIPTS_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 323 — BP-PY-45

- Function context: `./scripts/findings/functions/323.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/rules/generate_rule_triads.py:14:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 324 — BP-PY-45

- Function context: `./scripts/findings/functions/324.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/rules/ingest_rule_sources.py:21:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 325 — PERF-PY-26

- Function context: `./scripts/findings/functions/325.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/rules/ingest_rule_sources.py:119:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
)
args = parser.parse_args()

```

Why this is a false positive: the parse/decode runs once at CLI startup (`parse_args`), once per report build, or once per distinct file element in an offline batch — not a repeated hot-path parse without cache.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 326 — BP-PY-45

- Function context: `./scripts/findings/functions/326.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/rules/migrate_provisional_templates.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 327 — BP-PY-45

- Function context: `./scripts/findings/functions/327.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/rules/migrate_rules_v4.py:17:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 328 — BP-PY-46

- Function context: `./scripts/findings/functions/328.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/run_browser_smoke.py:50:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if node_path is None:
print("node is required to run the browser smoke suite.")
return 2
```

Why this is a false positive: browser-smoke CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 329 — BP-PY-46

- Function context: `./scripts/findings/functions/329.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/run_browser_smoke.py:57:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if not playwright_cli.exists():
print("Playwright is not installed in frontend/. Run `npm --prefix frontend install` first.")
return 2
```

Why this is a false positive: browser-smoke CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 330 — CWE-88

- Function context: `./scripts/findings/functions/330.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/run_browser_smoke.py:60:24`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

executable_probe = subprocess.run(
[
```

Why this is a false positive: `node_path` comes from `shutil.which` and the `-e` script is a fixed constant; no externally influenced argument.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 331 — BP-PY-46

- Function context: `./scripts/findings/functions/331.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/run_browser_smoke.py:73:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if executable_probe.returncode != 0:
print(executable_probe.stderr.strip() or "Unable to resolve Playwright browser path.")
return executable_probe.returncode or 1
```

Why this is a false positive: browser-smoke CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 332 — CWE-396

- Function context: `./scripts/findings/functions/332.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/run_browser_smoke.py:195:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
return result
except Exception as exc:
_write_report(
```

Why this is a false positive: the exception text is recorded in the failure report (`_write_report(..., str(exc))`); the failure is surfaced.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 333 — BP-PY-46

- Function context: `./scripts/findings/functions/333.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/run_frontend_accessibility.py:55:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if node_path is None:
print("node is required to run the accessibility walkthrough.")
return 2
```

Why this is a false positive: accessibility CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 334 — BP-PY-46

- Function context: `./scripts/findings/functions/334.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/run_frontend_accessibility.py:62:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if not playwright_cli.exists():
print("Playwright is not installed in frontend/. Run `npm --prefix frontend install` first.")
return 2
```

Why this is a false positive: accessibility CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 335 — CWE-88

- Function context: `./scripts/findings/functions/335.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/run_frontend_accessibility.py:65:24`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

executable_probe = subprocess.run(
[
```

Why this is a false positive: `node_path` from `shutil.which` plus a fixed `-e` constant script.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 336 — BP-PY-46

- Function context: `./scripts/findings/functions/336.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/run_frontend_accessibility.py:78:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if executable_probe.returncode != 0:
print(executable_probe.stderr.strip() or "Unable to resolve Playwright browser path.")
return executable_probe.returncode or 1
```

Why this is a false positive: accessibility CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 337 — CWE-396

- Function context: `./scripts/findings/functions/337.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/run_frontend_accessibility.py:214:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
return result
except Exception as exc:
_write_report(
```

Why this is a false positive: the exception text is recorded in the failure report.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 338 — BP-PY-46

- Function context: `./scripts/findings/functions/338.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/run_frontend_performance.py:59:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if node_path is None:
print("node is required to run the performance walkthrough.")
return 2
```

Why this is a false positive: performance CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 339 — BP-PY-46

- Function context: `./scripts/findings/functions/339.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/run_frontend_performance.py:66:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if not playwright_cli.exists():
print("Playwright is not installed in frontend/. Run `npm --prefix frontend install` first.")
return 2
```

Why this is a false positive: performance CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 340 — CWE-88

- Function context: `./scripts/findings/functions/340.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/run_frontend_performance.py:69:13`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

probe = subprocess.run(
[
```

Why this is a false positive: `node_path` from `shutil.which` plus a fixed `-e` constant script.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 341 — BP-PY-46

- Function context: `./scripts/findings/functions/341.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/run_frontend_performance.py:82:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if probe.returncode != 0:
print(probe.stderr.strip() or "Unable to resolve Playwright browser path.")
return probe.returncode or 1
```

Why this is a false positive: performance CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 342 — CWE-396

- Function context: `./scripts/findings/functions/342.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/run_frontend_performance.py:187:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
return _run_performance(base_url, report_path)
except Exception as exc:
_write_failure(report_path, base_url, str(exc))
```

Why this is a false positive: the exception text is recorded in the failure report.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 343 — BP-PY-46

- Function context: `./scripts/findings/functions/343.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/run_golden_journeys.py:55:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if node_path is None:
print("node is required to run the golden browser journeys.")
return 2
```

Why this is a false positive: golden-journeys CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 344 — BP-PY-46

- Function context: `./scripts/findings/functions/344.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/run_golden_journeys.py:62:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if not playwright_cli.exists():
print("Playwright is not installed in frontend/. Run `npm --prefix frontend install` first.")
return 2
```

Why this is a false positive: golden-journeys CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 345 — CWE-88

- Function context: `./scripts/findings/functions/345.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/run_golden_journeys.py:65:24`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

executable_probe = subprocess.run(
[
```

Why this is a false positive: `node_path` from `shutil.which` plus a fixed `-e` constant script.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 346 — BP-PY-46

- Function context: `./scripts/findings/functions/346.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/run_golden_journeys.py:78:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if executable_probe.returncode != 0:
print(executable_probe.stderr.strip() or "Unable to resolve Playwright browser path.")
return executable_probe.returncode or 1
```

Why this is a false positive: golden-journeys CLI; print is CLI output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 347 — CWE-396

- Function context: `./scripts/findings/functions/347.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/run_golden_journeys.py:211:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
return result
except Exception as exc:
_write_report(
```

Why this is a false positive: the exception text is recorded in the failure report.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 348 — BP-PY-45

- Function context: `./scripts/findings/functions/348.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/sources/build_source_snapshot.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 349 — BP-PY-45

- Function context: `./scripts/findings/functions/349.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/sources/validate_dockets.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 350 — BP-PY-45

- Function context: `./scripts/findings/functions/350.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/spec/run_conformance_tests.py:16:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 351 — BP-PY-46

- Function context: `./scripts/findings/functions/351.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/telegram/bot_poc.py:99:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
offset = 0
print("Starting Telegram polling loop...")
while True:
```

Why this is a false positive: Telegram POC script; print is the script's output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 352 — BP-PY-46

- Function context: `./scripts/findings/functions/352.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/telegram/bot_poc.py:117:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
def run_demo(api_base: str) -> None:
print("[DEMO] /panchanga")
print(process_command("/panchanga", api_base))
```

Why this is a false positive: Telegram POC script; print is the script's output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 353 — BP-PY-46

- Function context: `./scripts/findings/functions/353.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/telegram/bot_poc.py:118:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
print("[DEMO] /panchanga")
print(process_command("/panchanga", api_base))
print("\n[DEMO] /upcoming 14")
```

Why this is a false positive: Telegram POC script; print is the script's output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 354 — BP-PY-46

- Function context: `./scripts/findings/functions/354.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/telegram/bot_poc.py:119:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
print(process_command("/panchanga", api_base))
print("\n[DEMO] /upcoming 14")
print(process_command("/upcoming 14", api_base))
```

Why this is a false positive: Telegram POC script; print is the script's output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 355 — BP-PY-46

- Function context: `./scripts/findings/functions/355.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/telegram/bot_poc.py:120:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
print("\n[DEMO] /upcoming 14")
print(process_command("/upcoming 14", api_base))

```

Why this is a false positive: Telegram POC script; print is the script's output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 356 — BP-PY-45

- Function context: `./scripts/findings/functions/356.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/transparency/append_release_log.py:12:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 357 — BP-PY-45

- Function context: `./scripts/findings/functions/357.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/transparency/verify_log.py:12:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 358 — BP-PY-45

- Function context: `./scripts/findings/functions/358.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/vendor_audit/run_vendor_date_risk_audit.py:17:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 359 — CWE-22

- Function context: `./scripts/findings/functions/359.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/accuracy/test_prediction_run_immutability.py:28:45`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
def test_model_run_file_name_matches_run_id():
path = Path("data/future_bs/model_runs") / f"{DEFAULT_RUN_ID}.json"

```

Why this is a false positive: test constructs the path from a constant run id; no dynamic or untrusted segment.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 360 — CWE-22

- Function context: `./scripts/findings/functions/360.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/accuracy/test_source_policy_metrics_exist.py:12:51`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
]:
path = Path("data/future_bs/accuracy_lab") / name
assert path.exists(), path
```

Why this is a false positive: test iterates a module-constant list of filenames; no untrusted segment.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 361 — CWE-829

- Function context: `./scripts/findings/functions/361.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/architecture/test_canonical_runtime_registry.py:9:8`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

spec = importlib.util.spec_from_file_location("check_canonical_runtime", SCRIPT_PATH)
assert spec is not None and spec.loader is not None
```

Why this is a false positive: test loads a hard-coded repo path constant (`PROJECT_ROOT / "scripts" / "check_canonical_runtime.py"`).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 362 — CWE-829

- Function context: `./scripts/findings/functions/362.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/architecture/test_deprecated_modules_not_imported_by_public_routes.py:9:8`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

spec = importlib.util.spec_from_file_location("check_canonical_runtime", SCRIPT_PATH)
assert spec is not None and spec.loader is not None
```

Why this is a false positive: test loads a hard-coded repo path constant.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 363 — CWE-829

- Function context: `./scripts/findings/functions/363.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/architecture/test_no_public_route_imports_research_private.py:9:8`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

spec = importlib.util.spec_from_file_location("check_canonical_runtime", SCRIPT_PATH)
assert spec is not None and spec.loader is not None
```

Why this is a false positive: test loads a hard-coded repo path constant.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 364 — CWE-829

- Function context: `./scripts/findings/functions/364.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/architecture/test_no_runtime_tests_fixture_dependency.py:9:8`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

spec = importlib.util.spec_from_file_location("check_canonical_runtime", SCRIPT_PATH)
assert spec is not None and spec.loader is not None
```

Why this is a false positive: test loads a hard-coded repo path constant.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 365 — CWE-829

- Function context: `./scripts/findings/functions/365.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/architecture/test_runtime_does_not_depend_on_tests_fixtures.py:9:8`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

spec = importlib.util.spec_from_file_location("check_canonical_runtime", SCRIPT_PATH)
assert spec is not None and spec.loader is not None
```

Why this is a false positive: test loads a hard-coded repo path constant.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 366 — BP-PY-41

- Function context: `./scripts/findings/functions/366.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/artifacts/test_final_artifacts_exist.py:8:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

def test_verify_final_artifacts_script_passes():
subprocess.run(["python", "scripts/future_bs/generate_all_final_artifacts.py"], check=True)
```

Why this is a false positive: `subprocess.run(..., check=True)` raises on non-zero exit, so the test verifies through the raise; not a side-effect-only placeholder.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 367 — CWE-829

- Function context: `./scripts/findings/functions/367.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/benchmark/test_benchmark_schema.py:8:8`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
MODULE_PATH = ROOT / "public-benchmark" / "validate_benchmark.py"
spec = importlib.util.spec_from_file_location("validate_benchmark", MODULE_PATH)
assert spec and spec.loader
```

Why this is a false positive: test loads a hard-coded repo path constant.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 368 — CWE-829

- Function context: `./scripts/findings/functions/368.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/benchmark/test_runners.py:11:12`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
def _load(name: str, path: Path):
spec = importlib.util.spec_from_file_location(name, path)
assert spec and spec.loader
```

Why this is a false positive: test loads module-constant paths for repo runners.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 369 — CWE-93

- Function context: `./scripts/findings/functions/369.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/integration/test_engine_e2e.py:11:20`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
assert response.status_code == 200
assert response.headers["X-Parva-Engine"] == "v3"
assert response.headers["X-Parva-License"] == "AGPL-3.0-or-later"
```

Why this is a false positive: the test asserts on a response header value; it reads the header, it does not write one.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 370 — CWE-93

- Function context: `./scripts/findings/functions/370.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/integration/test_personal_stack_v3.py:55:16`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
assert resp.request.url.query == b""
assert resp.headers["Cache-Control"] == "no-store"

```

Why this is a false positive: the test asserts on a response header value; it reads the header.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 371 — CWE-93

- Function context: `./scripts/findings/functions/371.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/integration/test_private_surface_proof_capsules.py:23:20`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
assert response.request.url.query == b""
assert response.headers["Cache-Control"] == "no-store"
body = response.json()
```

Why this is a false positive: the test asserts on a response header value; it reads the header.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 372 — CWE-93

- Function context: `./scripts/findings/functions/372.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/integration/test_request_guards.py:47:20`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
assert response.status_code == 200
assert response.headers["Cache-Control"] == "no-store"
assert response.headers["Pragma"] == "no-cache"
```

Why this is a false positive: the test asserts on a response header value; it reads the header.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 373 — CWE-93

- Function context: `./scripts/findings/functions/373.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/integration/test_temporal_compass_api.py:14:20`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
assert response.status_code == 200
assert response.headers["Cache-Control"] == "no-store"
body = response.json()
```

Why this is a false positive: the test asserts on a response header value; it reads the header.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 374 — CWE-93

- Function context: `./scripts/findings/functions/374.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/integration/test_v3_envelope_opt_in.py:10:20`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
assert response.status_code == 200
assert response.headers["X-Parva-Envelope"] == "data-meta"
assert "X-Parva-Envelope" in response.headers.get("Vary", "")
```

Why this is a false positive: the test asserts on a response header value; it reads the header.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 375 — BP-PY-7

- Function context: `./scripts/findings/functions/375.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/release/test_archive_hygiene.py:41:16`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
with (
tarfile.open(archive_path) as archive,
pytest.raises(RuntimeError, match="escapes extraction root"),
```

Why this is a false positive: `tarfile.open(...) as archive` is already the context expression of a `with (...)` block.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 376 — BP-PY-7

- Function context: `./scripts/findings/functions/376.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/release/test_archive_hygiene.py:58:16`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
with (
tarfile.open(archive_path) as archive,
pytest.raises(RuntimeError, match="not a regular file"),
```

Why this is a false positive: `tarfile.open(...) as archive` is already the context expression of a `with (...)` block.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 377 — CWE-829

- Function context: `./scripts/findings/functions/377.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/runtime/test_canonical_runtime_imports.py:10:8`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

spec = importlib.util.spec_from_file_location("check_canonical_runtime", SCRIPT_PATH)
assert spec is not None and spec.loader is not None
```

Why this is a false positive: test loads a hard-coded repo path constant.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 378 — BP-PY-12

- Function context: `./scripts/findings/functions/378.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/security/test_embed_security.py:44:13`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
assert "createElement('script')" not in source
assert "eval(" not in source
assert "Function(" not in source
```

Why this is a false positive: the flagged token is a string-containment assertion (`assert "eval(" not in source`), not an eval call.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 379 — CWE-829

- Function context: `./scripts/findings/functions/379.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/unit/benchmark/test_public_benchmark_runners.py:12:12`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
path = ROOT / "public-benchmark" / "runners" / name
spec = importlib.util.spec_from_file_location(name.removesuffix(".py"), path)
module = importlib.util.module_from_spec(spec)
```

Why this is a false positive: test loads module-constant paths for repo runners.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 380 — BP-PY-13

- Function context: `./scripts/findings/functions/380.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/unit/billing/test_api_key_hashing.py:11:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
def test_api_key_hash_uses_salted_versioned_pbkdf2():
secret = "sample-secret"
pepper = "sample-pepper"
```

Why this is a false positive: test fixture placeholder values (`"sample-secret"`, `"sample-pepper"`), which the rule's note excludes.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 381 — CWE-93

- Function context: `./scripts/findings/functions/381.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/unit/bootstrap/test_middleware.py:44:17`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if forwarded_for is not None:
self.headers["x-forwarded-for"] = forwarded_for

```

Why this is a false positive: the fake request in the test stores a test-supplied `forwarded_for` value into its own header dict; not a server writing a response header.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 382 — BP-PY-12

- Function context: `./scripts/findings/functions/382.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/unit/bootstrap/test_rate_limit.py:17:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

def eval(self, script, numkeys, *args):
self.calls.append((script, numkeys, args))
```

Why this is a false positive: `def eval(self, script, numkeys, *args)` is a method definition on a fake redis client used to record calls; it is not an eval call.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 383 — CWE-94

- Function context: `./scripts/findings/functions/383.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/unit/bootstrap/test_rate_limit.py:17:9`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

def eval(self, script, numkeys, *args):
self.calls.append((script, numkeys, args))
```

Why this is a false positive: the flagged `eval` is the fake redis client's method definition (same construct as 382); no code generation occurs.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 384 — CWE-93

- Function context: `./scripts/findings/functions/384.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/unit/bootstrap/test_security_hardening.py:15:20`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
assert response.status_code == 200
assert response.headers["X-Content-Type-Options"] == "nosniff"
assert response.headers["X-Frame-Options"] == "DENY"
```

Why this is a false positive: the test asserts on a response header value; it reads the header.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 385 — BP-PY-2

- Function context: `./scripts/findings/functions/385.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/unit/calendar/test_calculator.py:109:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
assert result is None or hasattr(result, "start")
except (KeyError, RuntimeError, ValueError):
pass  # Expected - unsupported historical range can reject the request
```

Why this is a false positive: the test deliberately tolerates the expected rejection of unsupported historical ranges (comment documents it); the pass is the expectation.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 386 — CWE-390

- Function context: `./scripts/findings/functions/386.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/unit/calendar/test_calculator.py:109:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
assert result is None or hasattr(result, "start")
except (KeyError, RuntimeError, ValueError):
pass  # Expected - unsupported historical range can reject the request
```

Why this is a false positive: the expected rejection is the test outcome; recovery (assert before the handler) exists.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 387 — CWE-208

- Function context: `./scripts/findings/functions/387.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/unit/policy/test_policy_vm.py:18:20`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
flags=frozenset({TaintFlag.REVIEW_REQUIRED})
if authority == AuthorityTaint.STATIC_REFERENCE
else frozenset(),
```

Why this is a false positive: enum member comparison in a unit test; not a secret comparison.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 388 — CWE-93

- Function context: `./scripts/findings/functions/388.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/unit/reliability/test_request_id_and_trace_id.py:14:20`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
assert response.status_code == 200
assert response.headers["X-Request-ID"] == "phase08-trace"
assert response.headers["X-Trace-ID"] == "phase08-trace"
```

Why this is a false positive: the test asserts on a response header value; it reads the header.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 389 — CWE-208

- Function context: `./scripts/findings/functions/389.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/unit/trust/test_field_provenance.py:37:12`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
boundary = BoundaryVector.from_provenance(provenance)
assert boundary.authority == AuthorityTaint.STATIC_REFERENCE
assert boundary.review_state == "required"
```

Why this is a false positive: enum member comparison in a unit test; not a secret comparison.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 390 — CWE-208

- Function context: `./scripts/findings/functions/390.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/unit/trust/test_taint_algebra.py:33:84`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
)
assert apply_review_upgrade(value, AuthorityTaint.COMPUTED_CERTIFIED, witness).authority == AuthorityTaint.COMPUTED_CERTIFIED

```

Why this is a false positive: enum member comparison in a unit test; not a secret comparison.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 391 — BP-PY-41

- Function context: `./scripts/findings/functions/391.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tests/unit/trust/test_temporal_trust_tools.py:53:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```

def test_trust_schemas_parse_and_validate_examples():
for relative in [
```

Why this is a false positive: `validate_schema_file` raises `SchemaValidationError` on invalid schemas (verified in `tools/validate_schemas.py`), so the test verifies through the raise.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 392 — BP-PY-45

- Function context: `./scripts/findings/functions/392.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tools/conformance_runner/run.py:22:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 393 — BP-PY-1

- Function context: `./scripts/findings/functions/393.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tools/conformance_runner/run.py:557:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
results.append(CaseResult(case_id, True, "passed", _display_path(path)))
except Exception as exc:  # noqa: BLE001
results.append(CaseResult(case_id, False, str(exc), _display_path(path)))
```

Why this is a false positive: the handler records the exception (log message with `exc`, result/error field with `str(exc)`, warning text, or stderr output) and/or re-raises, so the failure is surfaced — it does not pass or continue without recording the exception type.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 394 — CWE-396

- Function context: `./scripts/findings/functions/394.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tools/conformance_runner/run.py:557:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
results.append(CaseResult(case_id, True, "passed", _display_path(path)))
except Exception as exc:  # noqa: BLE001
results.append(CaseResult(case_id, False, str(exc), _display_path(path)))
```

Why this is a false positive: the handler deliberately handles the failure: it records the exception in logs/results/warnings or re-raises a specific error; nothing is hidden.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 395 — PERF-PY-25

- Function context: `./scripts/findings/functions/395.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tools/conformance_runner/run.py:584:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if not case["executable"]:
results.append(CaseResult(case_id, True, "skipped: documented-only public issue", _display_path(path)))
continue
```

Why this is a false positive: per-case result record construction is inherent to the conformance runner's output.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 396 — BP-PY-1

- Function context: `./scripts/findings/functions/396.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tools/conformance_runner/run.py:588:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
results.append(CaseResult(case_id, True, "passed", _display_path(path)))
except Exception as exc:  # noqa: BLE001
results.append(CaseResult(case_id, False, str(exc), _display_path(path)))
```

Why this is a false positive: the handler records the exception (log message with `exc`, result/error field with `str(exc)`, warning text, or stderr output) and/or re-raises, so the failure is surfaced — it does not pass or continue without recording the exception type.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 397 — BP-PY-1

- Function context: `./scripts/findings/functions/397.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tools/conformance_runner/run.py:893:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
failures, results = run(Path(args.case_root), api_base_url=api_base_url)
except Exception as exc:  # noqa: BLE001
print(f"Conformance setup failed: {exc}", file=sys.stderr)
```

Why this is a false positive: the handler records the exception (log message with `exc`, result/error field with `str(exc)`, warning text, or stderr output) and/or re-raises, so the failure is surfaced — it does not pass or continue without recording the exception type.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 398 — BP-PY-45

- Function context: `./scripts/findings/functions/398.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tools/future_bs_audit/blinded_audit.py:15:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 399 — BP-PY-45

- Function context: `./scripts/findings/functions/399.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tools/parva-cli/parva_cli.py:13:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(PYTHON_SDK) not in sys.path:
sys.path.insert(0, str(PYTHON_SDK))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 400 — PERF-PY-26

- Function context: `./scripts/findings/functions/400.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tools/parva-cli/parva_cli.py:98:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
parser = build_parser()
args = parser.parse_args(argv)
try:
```

Why this is a false positive: the parse/decode runs once at CLI startup (`parse_args`), once per report build, or once per distinct file element in an offline batch — not a repeated hot-path parse without cache.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 401 — BP-PY-45

- Function context: `./scripts/findings/functions/401.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tools/release/verify_release.py:15:5`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if str(BACKEND_ROOT) not in sys.path:
sys.path.insert(0, str(BACKEND_ROOT))

```

Why this is a false positive: guarded module-top bootstrap (`if <root> not in sys.path`) with a constant repo path in a standalone in-tree script whose purpose is to import the uninstalled package; this is bootstrap usage the rule's own note excludes (it targets mutation outside tests and bootstrap).

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 402 — PERF-PY-25

- Function context: `./scripts/findings/functions/402.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tools/release/verify_release.py:246:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if not path.exists():
raise ReleaseVerificationError(f"missing schema listed in schemas_used: {schema_path}")
load_json(path)
```

Why this is a false positive: the flagged statement is a `raise ReleaseVerificationError(...)` on a missing-schema error path inside a loop — exception construction on the error path, not hot-path work.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 403 — BP-PY-1

- Function context: `./scripts/findings/functions/403.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tools/release/verify_release.py:271:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
messages = verify_release(manifest_path.resolve())
except Exception as exc:  # noqa: BLE001
print(f"release verification failed: {exc}", file=sys.stderr)
```

Why this is a false positive: the handler records the exception (log message with `exc`, result/error field with `str(exc)`, warning text, or stderr output) and/or re-raises, so the failure is surfaced — it does not pass or continue without recording the exception type.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 404 — CWE-396

- Function context: `./scripts/findings/functions/404.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tools/release/verify_release.py:271:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
messages = verify_release(manifest_path.resolve())
except Exception as exc:  # noqa: BLE001
print(f"release verification failed: {exc}", file=sys.stderr)
```

Why this is a false positive: the handler deliberately handles the failure: it records the exception in logs/results/warnings or re-raises a specific error; nothing is hidden.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 405 — BP-PY-7

- Function context: `./scripts/findings/functions/405.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tools/trust/append_log_entry.py:83:25`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
try:
lock_fd = os.open(lock_path, os.O_CREAT | os.O_EXCL | os.O_WRONLY)
except FileExistsError as exc:
```

Why this is a false positive: `os.open` returns a deliberately held lock-file descriptor (created with `O_EXCL` and closed after the lock is released); it is not `open()` without a context manager.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 406 — PERF-PY-25

- Function context: `./scripts/findings/functions/406.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tools/validate_schemas.py:186:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if key not in instance:
raise SchemaValidationError(f"{path}: missing required key {key!r}")

```

Why this is a false positive: `raise SchemaValidationError(...)` on a missing-key error path inside a validation loop; not hot-path object construction.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 407 — PERF-PY-25

- Function context: `./scripts/findings/functions/407.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tools/validate_schemas.py:195:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if additional is False:
raise SchemaValidationError(f"{path}: unexpected key {key!r}")
if isinstance(additional, dict):
```

Why this is a false positive: `raise SchemaValidationError(...)` on an unexpected-key error path inside a validation loop.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 408 — PERF-PY-25

- Function context: `./scripts/findings/functions/408.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tools/validate_schemas.py:204:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if pattern.search(text):
raise SchemaValidationError(f"{path}: public-safety pattern matched: {pattern.pattern}")

```

Why this is a false positive: `raise SchemaValidationError(...)` on a pattern-match error path inside a loop.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 409 — PERF-PY-25

- Function context: `./scripts/findings/functions/409.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tools/validate_schemas.py:209:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if example.get("publication_status") != "computed_prediction_not_official":
raise SchemaValidationError(f"{path}: future-risk example must be computed_prediction_not_official")
if example.get("corrected_value_included") is not False:
```

Why this is a false positive: `raise SchemaValidationError(...)` on an invalid-example error path inside a loop.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 410 — PERF-PY-25

- Function context: `./scripts/findings/functions/410.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tools/validate_schemas.py:211:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
if example.get("corrected_value_included") is not False:
raise SchemaValidationError(f"{path}: public future-risk example must not include corrected values")

```

Why this is a false positive: `raise SchemaValidationError(...)` on an invalid-example error path inside a loop.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 411 — BP-PY-1

- Function context: `./scripts/findings/functions/411.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tools/validate_schemas.py:245:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
print(f"ok: {path.relative_to(ROOT)}")
except Exception as exc:  # noqa: BLE001
failures.append(f"{path.relative_to(ROOT)}: {exc}")
```

Why this is a false positive: the handler records the exception (log message with `exc`, result/error field with `str(exc)`, warning text, or stderr output) and/or re-raises, so the failure is surfaced — it does not pass or continue without recording the exception type.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


### [ ] Finding 412 — CWE-396

- Function context: `./scripts/findings/functions/412.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/tools/validate_schemas.py:245:1`
- Checklist pattern: Based the decision on the rule condition and the shown source, not on application-specific knowledge.

Source excerpt:

```
print(f"ok: {path.relative_to(ROOT)}")
except Exception as exc:  # noqa: BLE001
failures.append(f"{path.relative_to(ROOT)}: {exc}")
```

Why this is a false positive: the handler deliberately handles the failure: it records the exception in logs/results/warnings or re-raises a specific error; nothing is hidden.

Checklist evidence: the flagged source does not satisfy the rule condition as stated in the rule metadata/detection notes.


## True positives

### CWE-1124 — Excessively Deep Nesting

| Finding | Source | Reason |
| --- | --- | --- |
| 41 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/calendar/bikram_sambat.py:614` | executable statement is nested at least six control-flow levels (verified in the enclosing function) |
| 273 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/check_public_claims.py:164` | executable statement is nested at least six control-flow levels (verified in the enclosing function) |
| 315 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/scripts/release/verify_clean_clone_assumptions.py:257` | executable statement is nested at least six control-flow levels (verified in the enclosing function) |

### CWE-1121 — Excessive McCabe Cyclomatic Complexity

| Finding | Source | Reason |
| --- | --- | --- |
| 62 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/rules/service.py:37` | `upcoming()` has at least twelve visible control-flow branches |

### BP-PY-1 — Bare Except Clause

| Finding | Source | Reason |
| --- | --- | --- |
| 70 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/app/services/rulelang_service.py:1362` | bare `except Exception:` continues with only a constant warning; the exception type is never recorded |
| 133 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/packages/parva-agent-tools/parva_tools/langchain.py:67` | optional-import failure falls back to descriptors without recording the exception |
| 136 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/packages/parva-agent-tools/parva_tools/llamaindex.py:16` | optional-import failure falls back to descriptors without recording the exception |

### CWE-396 — Declaration of Catch for Generic Exception

| Finding | Source | Reason |
| --- | --- | --- |
| 134 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/packages/parva-agent-tools/parva_tools/langchain.py:67` | same construct as 133 — catch-and-fallback without recording |
| 137 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/packages/parva-agent-tools/parva_tools/llamaindex.py:16` | same construct as 136 — catch-and-fallback without recording |

### CWE-1046 — Creation of Immutable Text Using String Concatenation

| Finding | Source | Reason |
| --- | --- | --- |
| 88 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/evaluate.py:383` | `md += f"..."` repeatedly concatenates immutable text inside a loop |

### PERF-PY-27 — Repeated Load Of Same Filesystem Path

| Finding | Source | Reason |
| --- | --- | --- |
| 124 | `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/Project_Parva/backend/tools/week3_ground_truth_pipeline.py:263` | `_rule_group()` re-reads and re-parses the same constant `festival_rules_v3.json` path for every record in the loop |

## Uncertain findings

None — every finding was classifiable from the rule condition and the shown source.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `./scripts/chunks`
- Function evidence: `./scripts/findings/functions`
- Validation: `git diff --check` — pass (exit 0)
