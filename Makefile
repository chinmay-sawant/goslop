.PHONY: build test vet fmt lint lint-all ci integration version help reference-metrics run bench
# Pure Go by default (go/ast parse path — no tree-sitter / CGO).
export CGO_ENABLED ?= 0

# Default scan path for product-style runs (Rust makefile parity).
SCAN_PATH ?= /home/chinmay/ChinmayPersonalProjects/gopdfsuit
# Rust: make run RUN_ARGS="--export-context --export-chunks --no-cache"
# Defaults match product reference-export surface; override with RUN_ARGS=... as needed.
RUN_ARGS ?= --export-context --export-chunks --no-cache

# Product benches (ns/op, B/op, allocs/op). Override path/time as needed.
#   make bench
#   make bench BENCHTIME=20x
BENCHTIME ?= 100x
BENCHPKG ?= ./internal/bench/
GOSLOP_BENCH_SCAN_PATH ?= $(SCAN_PATH)
export GOSLOP_BENCH_SCAN_PATH

build:
	go build -o bin/goslop ./cmd/goslop

test:
	go test ./...

# Focused integration harness (materialized fixture seed expectations).
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
	./bin/goslop --version

help:
	@echo "Targets: build test integration vet fmt lint lint-all ci version run reference-metrics bench"
	@echo "CGO_ENABLED=$(CGO_ENABLED) (0 = pure Go / go/ast; default)"
	@echo "run: product summary scan (profile all, --no-fail --no-terminal + RUN_ARGS)"
	@echo "  SCAN_PATH=$(SCAN_PATH)"
	@echo "  RUN_ARGS=$(RUN_ARGS)"
	@echo "reference-metrics: §12.4 parity baseline on REFERENCE_PATH (testing expected metrics; not a product name)"
	@echo "  REFERENCE_PATH=$(REFERENCE_PATH)  REFERENCE_PROFILE=$(REFERENCE_PROFILE)"
	@echo "bench: go test -bench (ns/op B/op allocs/op); BENCHTIME=$(BENCHTIME) GOSLOP_BENCH_SCAN_PATH=$(GOSLOP_BENCH_SCAN_PATH)"


# Comprehensive lint with all practical linters enabled. Keep in sync with the
# .golangci.yml if one exists. Run this before pushing to catch pre-existing
# issues that the default lint target does not cover.
lint-all:
	golangci-lint run -c .golangci.yml ./...

# Product-style scan (Rust makefile: --no-fail --no-terminal --profile all + RUN_ARGS).
# Default RUN_ARGS exports context/chunks and disables cache (reference-export surface).
#   make run
#   make run RUN_ARGS="--export-context --export-chunks --no-cache"
#   make run SCAN_PATH=./some/project
# Timing loop (optional): for i in {1..20}; do make run || break; done > 1.txt 2>&1
run: build
	@mkdir -p scripts/findings/functions scripts/chunks
	./bin/goslop --profile all --no-fail --no-terminal $(RUN_ARGS) $(SCAN_PATH)

# §12.4 export/scan parity gate on a fixed reference corpus (default: gopdfsuit).
# "reference" = testing jargon for expected baseline metrics — not a product name.
# Hard expected baseline (Rust 2026-07-29): 915 findings; 10/197/312/396 sev; top BP-1×181…
# Current Go delta is tracked in plans/port-phasewise-checklist.md §12.4.
REFERENCE_PATH ?= $(SCAN_PATH)
REFERENCE_PROFILE ?= all
reference-metrics: build
	@rm -rf scripts/findings/functions scripts/chunks
	@mkdir -p scripts/findings/functions scripts/chunks
	-./bin/goslop --profile $(REFERENCE_PROFILE) --format json --export-context --export-chunks --no-cache \
		--context-dir scripts/findings/functions --chunks-dir scripts/chunks \
		$(REFERENCE_PATH) > /tmp/goslop-reference-metrics.json
	@echo "--- summary (stderr above) ---"
	@python3 -c "import json,collections; fs=json.load(open('/tmp/goslop-reference-metrics.json')).get('findings',[]); print('findings',len(fs)); print('severity',dict(collections.Counter(f['severity'] for f in fs))); print('top',collections.Counter(f['rule_id'] for f in fs).most_common(5))"
	@echo -n "context files: "; ls scripts/findings/functions 2>/dev/null | wc -l
	@echo -n "chunk files: "; ls scripts/chunks 2>/dev/null | wc -l
	@echo "Reference hard targets: findings=915 sev=10h/197i/312l/396m top=BP-1×181,PERF-6×94,PERF-32×59,BP-5×50,PERF-230×44 exports=915+37"

# Product benchmarks via stdlib testing.B (no external harness).
#   make bench
#   make bench BENCHTIME=20x
bench:
	go test -run='^$$' -bench=. -benchmem -benchtime=$(BENCHTIME) $(BENCHPKG)
