.PHONY: build test test-go test-python vet fmt lint lint-all ci integration integration-go integration-python version help reference-metrics run run-python bench
# Pure Go by default (go/ast parse path — no tree-sitter / CGO).
export CGO_ENABLED ?= 0

# Default scan path for product-style runs (Rust makefile parity).
SCAN_PATH ?= /home/chinmay/ChinmayPersonalProjects/gopdfsuit
# Python-only product scan (opt-in languages = ["python"]).
PYTHON_SCAN_PATH ?= /home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/
PYTHON_CONFIG ?= templates/goslop-python.toml
# Rust: make run RUN_ARGS="--export-context --export-chunks --no-cache"
# Defaults match product reference-export surface; override with RUN_ARGS=... as needed.
RUN_ARGS ?= --export-context --export-chunks --no-cache
# Python runs default to no-cache text summary; override with PYTHON_RUN_ARGS=...
PYTHON_RUN_ARGS ?= --export-context --export-chunks --no-cache

# Product benches (ns/op, B/op, allocs/op). Override path/time as needed.
#   make bench
#   make bench BENCHTIME=20x
BENCHTIME ?= 100x
BENCHPKG ?= ./internal/bench/
GOSLOP_BENCH_SCAN_PATH ?= $(SCAN_PATH)
export GOSLOP_BENCH_SCAN_PATH

# Python Go packages (detectors, ruleset, fixture matrix + slim corpus).
PYTHON_TEST_PKGS ?= ./internal/lang/python/... ./ruleset/python/ ./tests/integration/python/

build:
	go build -o bin/goslop ./cmd/goslop

# Full suite (Go + Python packages). Prefer test-go / test-python when iterating.
test:
	go test ./...

# Go-only: everything except the Python language plugin, Python ruleset, and
# Python integration matrices.
test-go:
	go test $$(go list ./... | grep -v -E '/internal/lang/python(/|$$)|/ruleset/python(/|$$)|/tests/integration/python(/|$$)')

# Python-only: plugin/detectors, Python ruleset validation, BP/CWE/PERF matrices + corpus.
#   make test-python
#   make test-python PYTHON_TEST_PKGS='./internal/lang/python/detectors/perf/ ./tests/integration/python/'
test-python:
	go test $(PYTHON_TEST_PKGS)

# Focused integration harness (materialized fixture seed expectations).
# Go and Python matrices live in separate packages so DefaultRegistry (Go-only)
# and LanguagePython scans do not share one suite surface.
integration: integration-go integration-python

# Go fixture matrices (BP/CWE/PERF under tests/fixtures/go).
integration-go:
	go test ./tests/integration/

# Python fixture matrices + slim corpus (BP-PY / CWE / PERF-PY under tests/fixtures/python).
integration-python:
	go test ./tests/integration/python/

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
	@echo "Targets: build test test-go test-python integration integration-go integration-python vet fmt lint lint-all ci version run run-python reference-metrics bench"
	@echo "CGO_ENABLED=$(CGO_ENABLED) (0 = pure Go / go/ast; default)"
	@echo "test: full go test ./..."
	@echo "test-go: Go packages only (excludes internal/lang/python, ruleset/python, tests/integration/python)"
	@echo "test-python: Python detectors + ruleset + integration matrices/corpus"
	@echo "  PYTHON_TEST_PKGS=$(PYTHON_TEST_PKGS)"
	@echo "run: product summary scan (profile all, --no-fail --no-terminal + RUN_ARGS)"
	@echo "  SCAN_PATH=$(SCAN_PATH)"
	@echo "  RUN_ARGS=$(RUN_ARGS)"
	@echo "run-python: Python-only scan (languages=[\"python\"] via PYTHON_CONFIG)"
	@echo "  PYTHON_SCAN_PATH=$(PYTHON_SCAN_PATH)"
	@echo "  PYTHON_CONFIG=$(PYTHON_CONFIG)"
	@echo "  PYTHON_RUN_ARGS=$(PYTHON_RUN_ARGS)"
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

# Python-only product scan (epic #51 heuristics; opt-in languages via config).
#   make run-python
#   make run-python PYTHON_SCAN_PATH=/path/to/project
#   make run-python PYTHON_RUN_ARGS="--format json --no-cache"
run-python: build
	@test -f "$(PYTHON_CONFIG)" || (echo "missing PYTHON_CONFIG=$(PYTHON_CONFIG)" >&2; exit 1)
	@test -d "$(PYTHON_SCAN_PATH)" || (echo "missing PYTHON_SCAN_PATH=$(PYTHON_SCAN_PATH)" >&2; exit 1)
	./bin/goslop --profile all --no-fail --no-terminal --config "$(PYTHON_CONFIG)" $(PYTHON_RUN_ARGS) "$(PYTHON_SCAN_PATH)"

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
