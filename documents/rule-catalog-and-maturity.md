# Rule catalog and maturity

The runtime catalog currently exposes 239 `PERF-*`, 175 `CWE-*`, and 135
`BP-*` rules. The command output is authoritative because it reflects the
rules registered in the binary you are running (Go plugin only today).

```sh
./bin/goslop --list-rules
./bin/goslop --list-rules | rg 'PERF-'
./bin/goslop --explain PERF-101
./bin/goslop --explain CWE-89
```

There is no category filter or `rules` subcommand. Filter the plain-text list
in the shell when needed.

## Catalogue roots (JSON)

| Language | Path | Notes |
|----------|------|--------|
| **Go** (default, loaded by generators / runtime metadata) | `ruleset/golang/` | `bad-practices.json`, `chunks/cwe-*.json`, `chunks/perf-*.json` |
| **Python** (WIP seeds; not registered in the binary yet) | `ruleset/python/` | Portable CWE shells + empty BP map; see README and `plans/v0.0.2/ruleset-reuse-audit.md` |

Reuse policy: same JSON field shapes across languages; rewrite
`detection_notes` / `applicable_to` per language. Do not bulk-copy Go PERF or
BP into Python.

## Maturity labels

Every listed rule has one of these labels:

| Label | Meaning |
|---|---|
| `production` | Eligible for default packs and intended for normal CI use |
| `experimental` | Available for review but not certified as a hard CI gate |
| `fixture-only` | Long-tail or museum coverage; opt in deliberately |
| `quarantined` | Reserved or incomplete; available only through an explicit/full selection |

Maturity describes rollout confidence, not severity. A high-severity finding
and a production label are separate properties.

## Finding a rule’s documentation

- [go-recommended-pack.md](./go-recommended-pack.md) explains profile
  selection and the recommended pack.
- [perf-rules.md](./perf-rules.md) contains human notes for part of the PERF
  catalog.
- [taint.md](./taint.md) documents the experimental taint model.
- `--explain RULE_ID` prints the rule title, severity, pack, maturity,
  description, and suggested fix from the live metadata.

For detector authors, metadata and pack selection live under `internal/rules/`;
the Go detector registries are under `internal/lang/go/detectors/`.
Python catalogue seeds live under `ruleset/python/` (no detector package yet).
