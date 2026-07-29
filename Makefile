.PHONY: build test vet fmt lint lint-all ci integration version help oracle
# CGO is required for tree-sitter Go bindings.
export CGO_ENABLED ?= 1

build:
	go build -o bin/codehound ./cmd/codehound

test:
	go test ./...

# Focused integration harness (materialized fixture seed oracle).
integration:
	go test ./tests/integration/...

vet:
	go vet ./...

fmt:
	gofmt -w .

lint: vet
	gofmt -l . | tee /tmp/gofmt.out; test ! -s /tmp/gofmt.out

# Local CI parity: vet + test + build (matches .github/workflows/ci.yml).
ci: lint test build

version: build
	./bin/codehound --version

help:
	@echo "Targets: build test integration vet fmt lint ci version"
	@echo "CGO_ENABLED=$(CGO_ENABLED) (required for tree-sitter)"


# Comprehensive lint with all practical linters enabled. Keep in sync with the
# .golangci.yml if one exists. Run this before pushing to catch pre-existing
# issues that the default lint target does not cover.
lint-all:
	golangci-lint run -c .golangci.yml ./...
# §12.4-style export gate. Set ORACLE_PATH to the reference corpus for hard
# counts (915 findings / 915 context / 37 chunks). Default exercises wiring only.
ORACLE_PATH ?= .
ORACLE_PROFILE ?= all
oracle: build
	@rm -rf scripts/findings/functions scripts/chunks
	-./bin/codehound --profile $(ORACLE_PROFILE) --export-context --export-chunks --no-cache $(ORACLE_PATH)
	@echo "--- export counts ---"
	@echo -n "context files: "; ls scripts/findings/functions 2>/dev/null | wc -l
	@echo -n "chunk files: "; ls scripts/chunks 2>/dev/null | wc -l
	@echo "See plans/port-phasewise-checklist.md §12.4 for hard oracle targets."
