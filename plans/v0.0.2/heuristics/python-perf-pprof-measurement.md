# v0.0.2 — Python scan pprof baseline + after (2026-08-02)

> **Parent:** [`python-perf-pprof-optimization-checklist.md`](./python-perf-pprof-optimization-checklist.md)  
> **Status:** pre + after captured (P0 Mask/firstMatch/Funcs · P1 gates/MustCompile · P2 inLoop + BP-PY-13)  
> **Artifacts:** `/tmp/goslop-python-pprof/*`  
> **Corpus:** `PYTHON_SCAN_PATH=/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/pythoncoreengine` (63 files / 20945 lines)  
> **Compare:** Go `SCAN_PATH=gopdfsuit` (78 files / 28042 lines)

---

## Product wall (cold cache, export on)

| Command | Wall (this machine) | Findings |
|---------|---------------------|----------|
| `make run` (Go, 3×) | **93.5 / 100.5 / 111.9 ms** | 915 |
| `make run-python` (pre, 5×) | **1.08 / 1.10 / 1.12 / 1.13 / 1.21 s** | 102 |
| `./bin/goslop` Python (after, 10× `--no-cache`) | mean **157.5** / median **154.6** / best **133.4** / worst **189.8** ms | **102** (parity) |

After severity / top rules (unchanged): 38 high, 56 info, 0 low, 8 medium · BP-PY-46×32, BP-PY-41×24, CWE-328×14, CWE-90×12, CWE-22×6.

Artifact: `product-run-python-after-10x.txt`.

---

## Engine benches (`benchtime=5s`, go1.26.4, `CGO_ENABLED=0`)

| Bench | ns/op | B/op | allocs/op | Artifact |
|-------|------:|-----:|----------:|----------|
| `BenchmarkScanProfileAll` (Go) | **125.1 ms** | 54.6 MB | 439k | `bench-go-scan-5s-compare.txt` |
| `BenchmarkPythonScanProfileAll` **pre** | **728.5 ms** | **552.9 MB** | 208k | `bench-scan-5s.txt` |
| `BenchmarkPythonScanProfileAll` **after** | **146.3 ms** | **24.5 MB** | **70.5k** | `bench-scan-after-5s.txt` |
| `BenchmarkPythonScanAndExport` pre | **627.4 ms** | 556.5 MB | 213k | `bench-scan-export-5s.txt` |

**Delta (ProfileAll):** **~5.0×** faster · **~22.5×** fewer bytes/op · **~3.0×** fewer allocs. Still ~**1.17×** slower than Go’s 125 ms on a *larger* corpus; **not yet &lt;100 ms**.

---

## CPU profile — pre (`scan-cpu.prof`, `-top -cum`)

Dominant stack (of 181.54s samples across 8 iters):

| Symbol | cum% | Notes |
|--------|-----:|-------|
| `cwe.(*PyCweScan).Run` | **80.1%** | entire CWE pack |
| `regexp.(*Regexp).doExecute` | **70.4%** | regex engine |
| `cwe.firstCodeMatchStart` | **48.0%** | `FindAllStringIndex` + per-call Mask |
| `cwe.firstMatchStart` | **14.7%** | secrets / wrappers |
| `cwe.pythonFunctions` | **13.0%** | remask + `FindAllStringSubmatchIndex` |
| `perf.(*PythonPerfScan).Run` | **11.2%** | secondary |
| `pytext.Mask` | **11.0%** | CPU share of remasking |
| `cwe.detectCWE208` | **10.7%** | single-rule standout |
| `cwe.findCalls` | **8.7%** | remask per call-site scan |
| `perf.inLoop` | **7.6%** | TrimSpace walk |
| `bad_practices.(*PythonBadPracticeScan).Run` | **1.8%** | already cheap |

---

## CPU profile — after (`scan-cpu-after.prof`, `-top -cum`)

Total samples ~92s over ~6s wall (38 iters):

| Symbol | cum% | vs pre |
|--------|-----:|--------|
| `cwe.(*PyCweScan).Run` | **71.4%** | 80.1% → 71.4% (absolute time down ~5× with bench) |
| `regexp.(*Regexp).doExecute` | **78.4%** | still dominant work |
| `cwe.firstCodeMatchStart` / `…Masked` | **24.8%** | **48.0% → 24.8%** |
| `cwe.firstMatchStart` | **12.8%** | 14.7% → 12.8% |
| `perf.(*PythonPerfScan).Run` | **16.5%** | share rose as CWE absolute cost fell |
| `bad_practices.(*PythonBadPracticeScan).Run` | **7.6%** | share rose; still small absolute |
| `cwe.buildPythonFunctions` | **3.1%** | replaces old `pythonFunctions` **13%** |
| `cwe.findCalls` / `…Masked` | **3.2%** | was **8.7%** |
| `pytext.Mask` | *(off top)* | was **11.0%** CPU |
| `perf.inLoop` | *(off top)* | was **7.6%** |

Artifact: `scan-cpu-after-top-cum.txt`.

---

## Memory profile — pre (`scan-mem.prof`)

### `alloc_space`

| Symbol | flat% | Notes |
|--------|------:|-------|
| `pytext.Mask` | **85.8%** | full-file string copy per call |
| `cwe.pythonCodeMask` (inline) | cum **85.8%** | every Mask site |
| `cwe.findCalls` | cum **60.5%** | remasks `unit.Source` repeatedly |
| `cwe.firstCodeMatchStart` | cum **20.0%** | Mask + regex |
| `cwe.(*PyCweScan).Run` | cum **94.8%** | pack total |

### `alloc_objects` (secondary)

| Symbol | flat% |
|--------|------:|
| `pytext.Mask` | 14.8% |
| `regexp.(*machine).alloc` | 10.1% |
| `strings.genSplit` | 7.3% |
| `perf.stripPyLineForFacts` | 5.0% |

---

## Memory profile — after (`scan-mem-after.prof`, `alloc_space`)

| Symbol | flat% | Notes |
|--------|------:|-------|
| `internal/bytealg.MakeNoZero` | **23.7%** | string builder growth |
| `strings.genSplit` | **11.4%** | line splits |
| `regexp.(*bitState).reset` | **10.8%** | regex scratch |
| `regexp.(*Regexp).get` | **9.2%** | regex machines |
| `pytext.Mask` | **6.07%** | **85.8% → 6.07%** (target was &lt;5%) |
| `perf.buildCodeLines` | **4.8%** | once-per-file PERF strip |
| `bad_practices.buildCodeLines` | **4.3%** | once-per-file BP strip |
| `cwe.(*PyCweScan).Run` | cum **40.2%** | was cum **94.8%** |
| `cwe.firstCodeMatchStart` | cum **8.7%** | was cum **20.0%** |
| `regexp.MustCompile` | **0.72%** | hot-path compiles collapsed |

Artifact: `scan-mem-after-alloc-space.txt`.

---

## Code-level root causes (validated against source)

1. **`BuildFacts` is index-only; detectors ignore facts** — ~~all 159 CWE rules take `_ *PyCweFacts`~~ **fixed**: `Masked` + `Funcs` once; rules take `facts *PyCweFacts`.
2. **`firstCodeMatchStart` uses `FindAllStringIndex`** — ~~fixed~~: iterative `FindStringIndex` on advancing offset with cached mask filter.
3. **159 CWE rules, only 5 with SourceIndex gates** — ~~fixed~~: **65 gated / 94 ungated** after P1 wiring (secrets/injection/code-dynamic/208/tier-b inventory).
4. **No triple Mask across packs** — PERF/BP each strip lines **once** via their own `buildCodeLines` (keep literals); only CWE calls `pytext.Mask`.
5. **PERF secondary:** ~~`inLoop` rescans with `TrimSpace`~~ **fixed**: `computeInLoop` + cached `trim` + cheap needles before membership.
6. **BP:** ~~`detectBPPY13` per-line `MustCompile`~~ **fixed**: package `secretExactAssignRe` + single gate-miss `ToLower`.

**Residual (blocks &lt;100 ms):** CWE regex still ~70%+ of CPU (`doExecute` / `firstCodeMatchStart` ~25%). Next levers: soften remaining heavy regexes (e.g. timing compare), more aggressive FN-safe gates on ungated 94, PERF 3.3 residual, only then export (3.4).

---

## Reproduce

```sh
export CGO_ENABLED=0 PATH="$HOME/go/bin:$PATH"
export GOSLOP_BENCH_PYTHON_SCAN_PATH=/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/pythoncoreengine
mkdir -p /tmp/goslop-python-pprof

go1.26.4 test -run='^$' -bench=BenchmarkPythonScanProfileAll -benchmem -benchtime=5s \
  -cpuprofile=/tmp/goslop-python-pprof/scan-cpu-after.prof \
  -memprofile=/tmp/goslop-python-pprof/scan-mem-after.prof \
  ./internal/bench/

go1.26.4 tool pprof -top -cum /tmp/goslop-python-pprof/scan-cpu-after.prof
go1.26.4 tool pprof -top -sample_index=alloc_space /tmp/goslop-python-pprof/scan-mem-after.prof
```

Product wall:

```sh
go1.26.4 build -o bin/goslop ./cmd/goslop
./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml \
  --export-context --export-chunks --no-cache \
  /home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/pythoncoreengine
```

Harness: `internal/bench/make_run_python_bench_test.go`.
