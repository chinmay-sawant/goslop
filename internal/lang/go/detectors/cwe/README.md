# Go CWE detectors

Structural and taint-lite security heuristics for Go.

## Seed (implemented)

| Rule ID | Kind | Notes |
|---------|------|-------|
| `CWE-78` | taint-lite / structural | `exec.Command` / `CommandContext` + request-derived args / shell `-c` |
| `CWE-89` | structural | dynamic SQL into `.Query`/`.Exec`/… (concat, `fmt.Sprintf`, tainted idents) |

## Registry inventory (remaining)

Source TOMLs under `registry/` (copied from Rust). Counts are `[[detector]]` rows:

| Registry file | Rows | Status |
|---------------|-----:|--------|
| `registry.access_control.toml` | 30 | pending |
| `registry.concurrency.toml` | 6 | pending |
| `registry.configuration.toml` | 6 | pending |
| `registry.credentials_and_secrets.toml` | 15 | pending |
| `registry.cryptography.toml` | 9 | pending |
| `registry.deserialization.toml` | 3 | pending |
| `registry.file_handling.toml` | 1 | pending |
| `registry.general_security.toml` | 75 | pending |
| `registry.information_exposure.toml` | 12 | pending |
| `registry.injection.toml` | 7 | seed: 78, 89 |
| `registry.input_validation.toml` | 5 | pending |
| `registry.input_validation_redos.toml` | 2 | pending |
| `registry.network_binding.toml` | 1 | pending |
| `registry.path_traversal.toml` | 1 | pending |
| `registry.request_handling.toml` | 2 | pending |
| **Total** | **175** | |

Full taint graph (Phase 9) will replace/augment seed CWE-78/89 heuristics.
