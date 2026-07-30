## Summary

Implements **batch 08 production hardening** Python bad-practice heuristics **BP-PY-48**, **BP-PY-49**, and **BP-PY-50** as pure-Go source-pattern detectors under `internal/lang/python/detectors/bad_practices/`. **BP-PY-14** (requests without timeout) is intentionally **not** included — owned by batch-01.

Relates to #53, Relates to #51

---

## Motivation / context

- Plan: `plans/v0.0.2/heuristics/bp-plans/batch-08-prod-hardening.md`
- Parent ledger: `plans/v0.0.2/heuristics/python-heuristics-bp.md` (#53 under epic #51)
- Catalogue rows already present in `ruleset/python/bad-practices.json`

---

## Changes

### New detectors (`rules_prod.go`)

| Rule | Sev | Detection |
|------|-----|-----------|
| **BP-PY-48** | high | CORS `*` origins with credentials: FastAPI `CORSMiddleware` (`allow_origins=["*"]` + `allow_credentials=True`), flask-cors (`CORS(..., supports_credentials=True)` + star origins), django-cors-headers (`CORS_ALLOW_ALL_ORIGINS=True` + `CORS_ALLOW_CREDENTIALS=True`) |
| **BP-PY-49** | high | TLS verify disabled: `verify=False`, `ssl._create_unverified_context`, `CERT_NONE`, optional `assert_hostname=False` |
| **BP-PY-50** | medium | Insecure cookie flags: `SESSION_COOKIE_SECURE` / `CSRF_COOKIE_SECURE` / `SESSION_COOKIE_HTTPONLY = False`; Flask `app.config[...] = False` |

### Supporting

- `facts.go` — needles for early-out (`CORSMiddleware`, `verify=False`, cookie flag names, …)
- `metaByID` populated in `rules_prod.go` `init()` before `RegisterRule`
- Tests: `rules_prod_test.go` hit/miss + registration; `scan_test.go` want-list includes 48–50
- Inventory / plan ledgers updated; **BP-PY-14** remains missing

### Explicit non-goals

- **BP-PY-14** requests timeout (batch-01)
- Full multi-file CORS config resolution
- Inferring missing Secure flags from `DEBUG=False` alone

---

## Impact

| Area | Impact |
|------|--------|
| **Performance** | Negligible; Python opt-in path only. Needle-gated source scans. |
| **Memory** | Small fact bag reuse; no new AST |
| **Behavior** | New findings for Python when bad-practices + Python language enabled |
| **API / CLI** | No new flags |
| **Dependencies** | None |
| **File size** | `rules_prod.go` ~286 lines (soft 1500 / hard 2000) |

---

## Breaking changes / migration

None. Additive detectors only.

---

## Test plan

- [x] `gofmt` on touched Go files
- [x] `go test ./internal/lang/python/detectors/bad_practices/ -count=1`
- [x] `make lint`
- [x] `make test`
- [x] Line budget: `rules_prod.go` ≤ 2000

---

## Related issues

- Relates to #53
- Relates to #51
