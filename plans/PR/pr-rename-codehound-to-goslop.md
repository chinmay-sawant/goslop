## Summary

Rename the product from codehound to **goslop** across the Go module, CLI binary, config files, ignore directives, cache/baseline paths, fingerprints, SARIF driver metadata, Makefile, and CI. Removes leftover codehound naming from the runtime surface so the repo matches the goslop product brand.

---

## Motivation / context

- Docs already brand the tool as goslop; the code still used codehound paths and identifiers.
- Plans: product rebrand follow-up after #27
- Issues: see **Related issues**

---

## Changes

### Module and CLI

- Module path: `github.com/chinmay/codehound` → `github.com/chinmay/goslop`
- `cmd/codehound` → `cmd/goslop`
- Binary: `bin/goslop` (`make build`, CI, Makefile run/oracle)

### Config and on-disk artifacts

- `goslop.toml` (was `codehound.toml`); TOML table `[goslop]`
- `goslop.schema.json` and `templates/goslop.toml`
- Cache dir: `.goslop-cache`
- Baseline: `.goslop-baseline.json`
- Path ignore file: `.goslopignore`

### Suppressions and identity

- Directives: `// goslop-ignore`, `goslop-ignore-file`, start/end blocks
- Fingerprint prefix: `goslop:2:…`
- SARIF driver name / information URI updated to goslop

### Tooling

- Makefile, `.github/workflows/ci.yml`, `.goreleaser.stub.yml`, `.golangci.yml`, tests/fixtures imports and comments

---

## Impact

| Area | Impact |
|------|--------|
| **Performance** | None |
| **Memory** | None |
| **Behavior / correctness** | Same detectors; renames only. Existing baselines/caches/ignore comments using codehound names will not apply until migrated |
| **API / CLI** | Binary and config filenames change to goslop |
| **Dependencies** | Module path change for importers |
| **Binary size / build time** | None |

---

## Breaking changes / migration

| Item | Migration |
|------|-----------|
| Binary name | Use `./bin/goslop` instead of `./bin/codehound` |
| Config file | Rename `codehound.toml` → `goslop.toml` and `[codehound]` → `[goslop]` |
| Cache | Delete `.codehound-cache` or rename; new default `.goslop-cache` |
| Baseline | Rename `.codehound-baseline.json` → `.goslop-baseline.json`; fingerprints change prefix |
| Ignore comments | `// codehound-ignore` → `// goslop-ignore` |
| Go imports | `github.com/chinmay/codehound/...` → `github.com/chinmay/goslop/...` |

---

## Test plan

- [x] `CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop`
- [x] `./bin/goslop --version` / `--list-rules`
- [x] `CGO_ENABLED=0 go test ./...`
- [x] `rg -i codehound` over source (excluding generated scripts/cache) empty
- [ ] `make lint` optional

### Commands

```sh
make build
make test
./bin/goslop --version
```

---

## Screenshots / sample output

```
0.1.0-dev
ok  github.com/chinmay/goslop/... (all packages)
```

---

## Related issues

- Relates to #27 (docs rebrand to goslop)
- Refs full removal of codehound product naming from runtime

---

## PR metadata checklist (author)

- [x] Self-assigned (`--assignee @me`)
- [x] Labels applied
- [x] Related issues filled
- [x] Filled body under `plans/PR/pr-rename-codehound-to-goslop.md`

---

## Follow-ups (out of scope)

- Dual-read of old `codehound.toml` / `// codehound-ignore` for migration (not implemented; clean break)
- Publish module path on public proxy if needed

---

## Reviewer checklist

- [ ] Behavior matches summary and test plan
- [ ] No unrelated changes in diff
- [ ] Public API / CLI changes documented
- [ ] PR has assignee and labels
- [ ] Related issues use correct Closes/Relates keywords
- [ ] No secrets or generated artifacts committed
