# Rule catalog and maturity

The runtime catalog currently exposes 239 `PERF-*`, 175 `CWE-*`, and 135
`BP-*` rules. The command output is authoritative because it reflects the
rules registered in the binary you are running.

```sh
./bin/goslop --list-rules
./bin/goslop --list-rules | rg 'PERF-'
./bin/goslop --explain PERF-101
./bin/goslop --explain CWE-89
```

There is no category filter or `rules` subcommand. Filter the plain-text list
in the shell when needed.

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
