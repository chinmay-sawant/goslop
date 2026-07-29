# Recommended pack

Default CLI profile: `--profile recommended` (also `CODEHOUND_PROFILE=recommended`).

Aliases: `ci`, `default` → `recommended`.

## Goal

A small, high-signal CI gate: **PERF footguns** teams actually fix, plus
**taint-core CWE IDs** on the allow-list (taint engine off until you enable it).

Bad practices (`BP-*`) are **off**. Fail policy defaults to **strict**
(high / critical only) unless you set `--warnings-as-errors` / `--no-fail`.

Full CLI: [cli.md](./cli.md). Rule navigation: [rule-catalog.md](./rule-catalog.md).

---

## What actually fails under recommended

| Finding class | Typical severity | Fails recommended default? |
|---------------|------------------|----------------------------|
| S-tier PERF | **medium** | **No** — visible only |
| Taint-core CWE (when fired) | **high** | **Yes** |
| BP | n/a (off) | — |

To fail CI on medium PERF (timeouts, body close, …):

```bash
codehound --profile recommended --warnings-as-errors .
```

Taint is **off** under recommended. The CWE IDs stay allow-listed so
`codehound --taint --profile recommended .` works without switching packs.
Use `--profile security` to turn taint on by default.

---

## Rules (exact list)

### PERF (S-tier)

| Rule | Why |
|------|-----|
| `PERF-1` | Regex compilation inside a loop |
| `PERF-7` | `defer` inside a loop (same function scope) |
| `PERF-50` | `regexp.MatchString` inside a loop |
| `PERF-58` | Gin `c.Request.Body` not closed |
| `PERF-71` | GORM N+1 query pattern |
| `PERF-101` | `http.Server` missing timeouts |
| `PERF-103` | HTTP response body not closed |
| `PERF-189` | HTTP response body not drained before close |
| `PERF-190` | HTTP client missing timeout |

See [perf-tiers.md](./perf-tiers.md) for S/A/B/C policy.

### CWE (taint-core allow-list)

| Rule | Why |
|------|-----|
| `CWE-22` | Path traversal (taint) |
| `CWE-78` | OS command injection (taint) |
| `CWE-79` | XSS / template+HTTP write (taint) |
| `CWE-89` | SQL injection heuristic (taint) |
| `CWE-90` | LDAP injection (taint) |
| `CWE-91` | XML injection (taint) |

---

## Other profiles

| Profile | Aliases | Contents | Taint | BP | Default fail |
|---------|---------|----------|-------|----|--------------|
| `recommended` | `ci`, `default` | S PERF + taint-core CWE | off | off | strict |
| `perf` | — | S+A PERF (`PERF_TIER_S` + `PERF_TIER_A`) | off | off | strict |
| `security` | `sec` | Security pack CWEs (below) | **on** | off | strict |
| `style` | `bp`, `bad-practices` | `BP-*` | off | on | no-fail |
| `all` | `full` | Full catalog | off | on | medium-as-errors |

### Security pack rule IDs

Exact allow-list (`SECURITY_PACK_RULES` in `src/rules/pack.rs`):

`CWE-22`, `CWE-41`, `CWE-59`, `CWE-78`, `CWE-79`, `CWE-89`, `CWE-90`, `CWE-91`, `CWE-93`

Taint auto-enables under `--profile security`.

### Style pack defaults

- Allow pattern: all `BP-*`
- Default skip (opinionated / noisy): **`BP-21`**, **`BP-28`**, **`BP-30`**
  - BP-21: missing `t.Parallel` (policy)
  - BP-28: single-method interfaces (API style)
  - BP-30: external / capability interfaces
- Opt back in with `--only BP-28` (etc.) under `--profile style`.

### Perf profile

S-tier + A-tier numeric lists: [perf-tiers.md](./perf-tiers.md).

---

## Fixture-only quarantine

Rules tagged `fixture-only` are **never** in recommended/security/perf default
packs. They are available under `--profile all` (or explicit `--only <id>`), but
that does **not** mean they are production-certified for CI hard-fail gates.

Examples include long-tail CWE PRNG corpus patterns and other museum entries
audited in
[`plans/v0.0.5/cwe-catalog-trust-audit.md`](../plans/v0.0.5/cwe-catalog-trust-audit.md)
(archive note — historical audit).

Reserved rules (today: `BP-63`) follow the same quarantine: available under
`--profile all`, not for production CI packs until completed.

See `src/rules/maturity.rs` and `src/core/profile.rs`.

---

## Rule explainability

```bash
codehound rules --explain CWE-334
codehound --explain CWE-89
codehound rules --category security
codehound --list-rules --rule-category performance
```

The explain surface reuses the single maturity registry (`RuleMaturity` /
`maturity_for`); it does not invent a second rule-status model.

---

## CI one-liner

```bash
codehound --profile recommended --format sarif . > codehound.sarif
```

Fail on medium PERF too:

```bash
codehound --profile recommended --warnings-as-errors --format sarif . > codehound.sarif
```

Sample workflow: [ci-integration.md](./ci-integration.md),
[`.github/workflows/codehound.yml`](../.github/workflows/codehound.yml).

---

## Brownfield

1. Start advisory: `--no-fail` or upload SARIF without blocking merge.
2. Save a baseline: `codehound --profile recommended --baseline .`
3. Suppress known local noise: `// codehound-ignore: PERF-101` (or file-level).
4. Gate on new fingerprints; periodically `baseline diff|prune|update`.

Full guide: [finding-identity.md](./finding-identity.md),
[ci-integration.md](./ci-integration.md).
