.PHONY: build test vet fmt lint ci integration version help

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
