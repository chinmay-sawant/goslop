.PHONY: build test vet fmt lint lint-all ci integration version help oracle run
# CGO is required for tree-sitter Go bindings.
export CGO_ENABLED ?= 1

# Default scan path for product-style runs (Rust makefile parity).
SCAN_PATH ?= /home/chinmay/ChinmayPersonalProjects/gopdfsuit
# Rust: make run RUN_ARGS="--export-context --export-chunks --no-cache"
# Defaults match product oracle export surface; override with RUN_ARGS=... as needed.
RUN_ARGS ?= --export-context --export-chunks --no-cache

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
	@echo "Targets: build test integration vet fmt lint lint-all ci version run oracle"
	@echo "CGO_ENABLED=$(CGO_ENABLED) (required for tree-sitter)"
	@echo "run: product summary scan (profile all, --no-fail --no-terminal + RUN_ARGS)"
	@echo "  SCAN_PATH=$(SCAN_PATH)"
	@echo "  RUN_ARGS=$(RUN_ARGS)"


# Comprehensive lint with all practical linters enabled. Keep in sync with the
# .golangci.yml if one exists. Run this before pushing to catch pre-existing
# issues that the default lint target does not cover.
lint-all:
	golangci-lint run -c .golangci.yml ./...

# Product-style scan (Rust makefile: --no-fail --no-terminal --profile all + RUN_ARGS).
# Default RUN_ARGS exports context/chunks and disables cache (oracle surface).
#   make run
#   make run RUN_ARGS="--export-context --export-chunks --no-cache"
#   make run SCAN_PATH=./some/project
run: build
	@mkdir -p scripts/findings/functions scripts/chunks
	./bin/codehound --profile all --no-fail --no-terminal $(RUN_ARGS) $(SCAN_PATH)

# §12.4 export/scan oracle gate on SCAN_PATH (default: gopdfsuit).
# Hard oracle (Rust 2026-07-29): 915 findings; 10/197/312/396 sev; top BP-1×181…
# Current Go delta is tracked in plans/port-phasewise-checklist.md §12.4.
ORACLE_PATH ?= $(SCAN_PATH)
ORACLE_PROFILE ?= all
oracle: build
	@rm -rf scripts/findings/functions scripts/chunks
	@mkdir -p scripts/findings/functions scripts/chunks
	-./bin/codehound --profile $(ORACLE_PROFILE) --format json --export-context --export-chunks --no-cache \
		--context-dir scripts/findings/functions --chunks-dir scripts/chunks \
		$(ORACLE_PATH) > /tmp/codehound-oracle.json
	@echo "--- summary (stderr above) ---"
	@python3 -c "import json,collections; fs=json.load(open('/tmp/codehound-oracle.json')).get('findings',[]); print('findings',len(fs)); print('severity',dict(collections.Counter(f['severity'] for f in fs))); print('top',collections.Counter(f['rule_id'] for f in fs).most_common(5))"
	@echo -n "context files: "; ls scripts/findings/functions 2>/dev/null | wc -l
	@echo -n "chunk files: "; ls scripts/chunks 2>/dev/null | wc -l
	@echo "Oracle hard targets: findings=915 sev=10h/197i/312l/396m top=BP-1×181,PERF-6×94,PERF-32×59,BP-5×50,PERF-230×44 exports=915+37"
