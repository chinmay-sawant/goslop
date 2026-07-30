# Development commands

Goslop builds as pure Go. The module declares Go **1.26.4** and the Makefile
defaults `CGO_ENABLED=0`.

```sh
go mod download
make build
./bin/goslop --version
```

## Everyday checks

| Command | What it does |
|---|---|
| `make test` | Runs `go test ./...` |
| `make integration` | Runs the fixture integration harness |
| `make vet` | Runs `go vet ./...` |
| `make lint` | Runs vet and checks `gofmt` without changing files |
| `make lint-all` | Runs configured `golangci-lint` checks |
| `make ci` | Runs lint, test, and build, matching the local CI contract |
| `make version` | Builds and prints the binary version |

The GitHub Actions workflow runs Go module download, vet, tests, a pure-Go
build, and smoke checks for `--version` and `--list-rules`.

## Product scans and performance commands

```sh
# Build, scan a real application with profile all, and export review files.
make run SCAN_PATH=/path/to/go-project

# Compare a fixed corpus with the repository's expected export/count baseline.
make reference-metrics REFERENCE_PATH=/path/to/go-project

# Run the Go benchmark package with allocation metrics.
make bench
make bench BENCHTIME=20x
```

`make run` defaults to the local GopdfSuit path configured in the Makefile.
Override `SCAN_PATH` instead of relying on that machine-specific default.
`make reference-metrics` removes and recreates generated export directories;
use it only when those generated outputs are disposable. It intentionally
permits a nonzero scan exit when findings exist so that its metrics are still
reported.

Do not treat a self-scan finding count as a product precision metric: the
repository contains detector needles and intentional fixtures. See
[make-run.md](./make-run.md) for the product-scan contract and
[overview.md](./overview.md) for the self-scan explanation.

## Documentation-only changes

For a documentation-only update, verify Markdown links and run:

```sh
git diff --check
git diff --stat origin/main...HEAD
```

For code changes, run the narrowest relevant test while developing and the
applicable final commands above before opening a pull request. The full
contributor workflow lives in [`CONTRIBUTING.md`](../CONTRIBUTING.md).
