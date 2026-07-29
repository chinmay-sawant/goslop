# Go Performance Rules (PERF)

This document is a **partial human-notes catalog** for Go performance (`PERF-*`)
detectors — not a complete listing of all **239** shipped rules.

| Truth source | How to use |
|--------------|------------|
| Live inventory | `codehound --list-rules --rule-category performance` |
| Titles / detection notes | `ruleset/golang/chunks/perf-*.json` + `codehound --explain PERF-N` |
| Pack / tier policy | [perf-tiers.md](./perf-tiers.md), [go-recommended-pack.md](./go-recommended-pack.md) |
| Catalog hub | [rule-catalog.md](./rule-catalog.md) |

Canonical CLI IDs are unpadded (`PERF-1`, not `PERF-001`). Prefer the CLI and
ruleset JSON when prose here conflicts.

## 1 — Idiomatic Usage

- **PERF-101** flags `http.Server` values without `ReadTimeout` or `WriteTimeout` because unbounded connections can exhaust server resources. *Fix:* set both timeout fields to a reasonable duration (e.g. `5s` / `10s`).
- **PERF-102** flags HTTP `Transport` not shared across requests (new `http.Client` / transport per call thrash). *Fix:* reuse a shared `http.Client` / `http.Transport`.
- **PERF-103** flags HTTP response bodies that are never closed (resource leak). *Fix:* `defer resp.Body.Close()` (and drain when required — see PERF-189).
- **PERF-104** flags multiple calls to `w.WriteHeader` in the same handler because only the first call sets the status code. *Fix:* call `WriteHeader` once.
- **PERF-105** flags `runtime.SetFinalizer` on hot-path objects because finalizers add GC cost and surprising lifetime semantics. *Fix:* avoid finalizers on request-scoped objects; prefer explicit cleanup.
- **PERF-107** flags `encoding/binary.Read` or `binary.Write` inside a loop because the reflection-based API allocates per call. *Fix:* call `binary.Read`/`Write` once or batch the data outside the loop.
- **PERF-108** flags `sort.Search` inside a loop because binary search is already O(log n) and iterating it removes the benefit. *Fix:* compute the search key and call `sort.Search` once.
- **PERF-111** flags `fmt.Sprintf` with `%s` and a string argument because it is equivalent to a direct string concatenation or `string()` conversion. *Fix:* replace with `+` or the appropriate conversion.
- **PERF-112** flags `fmt.Sprintf` with `%d` and an integer argument in hot paths because `strconv.Itoa` / `strconv.FormatInt` is faster and allocation-free when the buffer is reused. *Fix:* use `strconv.Itoa` or `strconv.FormatInt`.
- **PERF-113** flags `fmt.Sprintf` with `%f` and a float argument because `strconv.FormatFloat` is faster in hot paths. *Fix:* use `strconv.FormatFloat`.
- **PERF-114** flags `fmt.Sprintf` with `%t` and a bool argument because `strconv.FormatBool` is faster. *Fix:* use `strconv.FormatBool`.
- **PERF-115** flags `fmt.Sprintf` with `%x` / `%X` and a byte slice because `hex.EncodeToString` is faster. *Fix:* use `encoding/hex.EncodeToString`.
- **PERF-116** flags `fmt.Sprintf` with `%q` because `strconv.Quote` is faster and more explicit. *Fix:* use `strconv.Quote`.
- **PERF-117** flags `fmt.Sprintf` with `%v` and a known primitive type because the explicit verb is faster and self-documenting. *Fix:* use the type-specific verb or `strconv` function.
- **PERF-118** flags `fmt.Sprintf` with `%p` because `reflect.ValueOf` (used by `%p`) is unnecessary for pointer formatting; use `fmt.Sprintf("%x", …)` or `%x` with a `uintptr`. *Fix:* use a more specific formatting approach.
- **PERF-119** flags `fmt.Sprintf` with a single format argument and no dynamic values because a plain string is simpler. *Fix:* use a plain string literal.
- **PERF-120** flags `time.Now().Sub(t)` where `time.Since(t)` is clearer and preferred. *Fix:* use `time.Since`.
- **PERF-121** flags allocations inside request-scoped loops where the total allocation is proportional to request concurrency. *Fix:* move the allocation outside the loop or use a `sync.Pool`.
- **PERF-122** flags `HasPrefix` + slice where `TrimPrefix` is the clearer API. *Fix:* use `strings.TrimPrefix` / `bytes.TrimPrefix`.
- **PERF-123** flags `json.Marshal` inside loops because marshalling allocates new memory per call. *Fix:* marshal outside the loop or use `json.Encoder` with a pooled buffer.
- **PERF-124** flags `strings.Replace(s, old, new, -1)` because the magic `-1` is less readable than `strings.ReplaceAll`. *Fix:* use `strings.ReplaceAll`.
- **PERF-125** flags `strings.TrimSuffix` / `strings.TrimPrefix` when a simple suffix/prefix check and slice would suffice, because the allocation is unnecessary in hot paths. *Fix:* use `strings.HasPrefix` + slice or `strings.Cut`.
- **PERF-126** flags `w.Header().Set(key, val)` inside a loop when the key is loop-invariant, because redundant header writes are wasteful. *Fix:* hoist the header set outside the loop.
- **PERF-127** flags unnecessary `fmt.Sprintf` inside log calls. *Fix:* pass structured fields / format args to the logger instead of pre-sprintf.
- **PERF-128** flags `for i, c := range s` on a string where only the rune count or length is needed, because decoding each rune is wasted work. *Fix:* use `utf8.RuneCountInString` or `len` instead.
- **PERF-129** flags `for range s` on a string where only the byte length matters, because the range decodes runes unnecessarily. *Fix:* use `len(s)`.
- **PERF-130** flags `for range m` on a map where only the key or value is needed, because map iteration with undesired value extraction wastes cycles. *Fix:* iterate only the needed fields.
- **PERF-131** flags empty `select {}` without a default case because it blocks forever. *Fix:* add a `default` case or a timeout.
- **PERF-132** flags calling `recover()` outside a deferred function because it is a no-op. *Fix:* move `recover()` into a `defer` block.
- **PERF-133** flags `sort.Slice` inside a loop because the less-function closure is re-allocated per iteration. *Fix:* hoist the sort outside the loop or use `sort.SliceStable` with a prepared less function.
- **PERF-135** flags redundant `strings.ToLower` / `strings.ToUpper` calls because the result is often compared or used immediately and a cached lower/upper variant may be reused. *Fix:* compute the lower/upper form once and reuse.
- **PERF-137** flags `runtime.Caller` in request handlers because it is expensive (stack-walking) and debug information is rarely needed on every request. *Fix:* remove `runtime.Caller` from the hot path or use structured logging.
- **PERF-140** flags creating a new `regexp.Regexp` inside a function called repeatedly, because regex compilation is expensive. *Fix:* use `regexp.MustCompile` at package level or `sync.Once`.
- **PERF-141** flags calling `r.URL.Query()` more than once in the same handler because each call parses the query string from scratch. *Fix:* call `r.URL.Query()` once and reuse the returned `Values`.
- **PERF-145** flags `http.Request.WithContext` allocations on hot middleware paths (advisory). *Fix:* avoid redundant `WithContext` when the request already carries the needed context.
- **PERF-146** flags `fmt.Sprintf` with a single string and no verbs. *Fix:* use the string directly.
- **PERF-147** flags `bytes.Equal` with a single-byte slice where a byte comparison is simpler. *Fix:* compare the single byte directly.
- **PERF-149** flags `net.Conn` `Read` / `Write` without a deadline because a missing deadline can hang the connection indefinitely. *Fix:* set `SetReadDeadline` / `SetWriteDeadline` before each I/O call.
- **PERF-153** flags package-level maps without an eviction strategy because unbounded caches grow indefinitely under load. *Fix:* add an eviction policy (e.g. `sync.Map` with periodic cleanup, or an LRU wrapper).
- **PERF-156** flags `copy` with overlapping source and destination slices because the behavior is undefined. *Fix:* ensure the slices do not overlap or use a temporary buffer.
- **PERF-157** flags unnecessary `fmt.Sprint` with a single string. *Fix:* use the string directly.
- **PERF-158** flags `sync.Pool` objects that are returned via deferred `Put` in a hot loop, because defer overhead adds up. *Fix:* call `Put` explicitly after each use rather than deferred.
- **PERF-161** flags `database/sql` rows iteration without checking `.Err()` after the loop. *Fix:* check `rows.Err()`.
- **PERF-163** flags `db.Query` for a query that returns at most one row, because `QueryRow` is simpler and avoids the boilerplate. *Fix:* use `db.QueryRow`.
- **PERF-165** flags `http.Handler` functions that use `r.URL.Query()` when the query parameter has a single expected value, because the `Values` map is overkill. *Fix:* extract the value directly.
- **PERF-166** flags `http.Handler` functions that call `fmt.Fprintf(w, …)` with a static string, because `w.Write` or `io.WriteString` is faster. *Fix:* use `io.WriteString` or `w.Write`.
- **PERF-167** flags `http.Redirect` with a relative path because the Go HTTP server will resolve it against the request host, which may not be the intended target. *Fix:* use an absolute URL or construct the redirect target server-side.
- **PERF-168** flags handler middleware that allocates a `log.Logger` per request, because logger creation is expensive. *Fix:* use a single package-level logger.
- **PERF-170** flags `sync.Once.Do` inside request handlers because `sync.Once` is designed for package-level initialization, not per-request use. *Fix:* move the initialization to package level.
- **PERF-171** flags `io.Copy` without checking the returned error, because partial writes can silently lose data. *Fix:* check the error from `io.Copy`.
- **PERF-176** flags `io.Copy` inside a loop because the buffer allocation is repeated. *Fix:* hoist the copy outside the loop or use `io.CopyBuffer` with a reused buffer.
- **PERF-177** flags `io.ReadAll` on an `*http.Response.Body` without a `LimitReader`, because unbounded reads can exhaust memory. *Fix:* use `io.LimitReader` to cap the response body size.
- **PERF-178** flags `fmt.Errorf` with `%s` / `%v` and `err.Error()` because the error value can be wrapped with `%w` instead. *Fix:* use `%w` to propagate the error.
- **PERF-179** flags repeated `time.Now()` calls in a single scope because each call incurs a syscall. *Fix:* call `time.Now()` once and reuse.
- **PERF-181** flags creating a `context.Context` with `context.WithTimeout` / `context.WithCancel` inside a loop without deferring the cancel, because resources leak until the context is garbage-collected. *Fix:* defer or call `cancel()` explicitly each iteration.
- **PERF-182** flags `time.Timer` without a `Stop` call in a loop, because the timer's channel and goroutine are not freed until it fires. *Fix:* call `timer.Stop()` and drain the channel if needed.
- **PERF-190** flags HTTP clients missing `Timeout` (or equivalent deadline). *Fix:* set `http.Client.Timeout` or per-request contexts.
- **PERF-192** flags creating a `*template.Template` from string inside a function called repeatedly, because template parsing is expensive. *Fix:* parse the template at package level with `template.Must`.
- **PERF-195** flags `log.Fatal` inside a goroutine because it calls `os.Exit(1)`, which does not give other goroutines a chance to clean up. *Fix:* return the error to the caller or signal the main goroutine.
- **PERF-198** flags `panic` in library code outside initialization, because panics in library code abort the caller's process without recovery options. *Fix:* return an error instead.
- **PERF-203** flags repeated `net.IP.String` on hot paths. *Fix:* cache the string form or avoid formatting per request.
- **PERF-204** flags `http.Server` without `ReadHeaderTimeout` set, because slow headers can hold connections open indefinitely. *Fix:* set `ReadHeaderTimeout` alongside `ReadTimeout`.
- **PERF-209** flags `json.Decoder` use without `DisallowUnknownFields` or `UseNumber` when the input schema is strictly controlled, because silent unknown fields mask data quality issues. *Fix:* enable strict decoding options.
- **PERF-211** flags `Request.ParseForm` called implicitly by `r.Form` / `r.PostForm` without checking the returned error, because parse errors can silently produce partial results. *Fix:* call `r.ParseForm` explicitly and check the error.

## 2 — Cache & Memory

- **PERF-213** flags package-level caches (maps) with reads and writes but no eviction bound, because unbounded caches grow until memory exhaustion. *Fix:* add an eviction policy (LRU, TTL, or size cap).
- **PERF-214** flags volatile or dynamic keys used in caching maps, because keys that include request-scoped values (pointers, timestamps) prevent cache hits. *Fix:* use stable, comparable keys.
- **PERF-215** flags missing buffer pre-sizing (`make([]byte, 0, expected)`), because append-based growth re-allocates. *Fix:* pre-size the buffer to the expected capacity.
- **PERF-216** flags struct allocation on hot paths where the struct is only used temporarily, because heap-allocated structs pressure the GC. *Fix:* reuse the struct or allocate on the stack by returning it by value.
- **PERF-217** flags static or quasi-static computation (e.g. deriving a configuration value) repeated per operation, because the result can be computed once at startup. *Fix:* pre-compute and store the result.
- **PERF-218** flags per-operation pool allocation where a single unsharded pool suffices, because per-goroutine pools increase memory pressure. *Fix:* use a single shared `sync.Pool`.
- **PERF-219** flags oversized pool returns where a pooled object is much larger than needed, because large retained objects waste memory. *Fix:* return a correctly-sized object or reset the oversized one.
- **PERF-220** flags repeated scans (linear searches) over the same data within a single operation, because each scan is O(n). *Fix:* build an index or hash map for lookups.
- **PERF-221** flags dense integer maps where a slice or array would be more efficient, because map overhead (hashing, bucket chaining) is unnecessary for compact integer keys. *Fix:* use a `[]T` indexed by the integer value.
- **PERF-222** flags generic function calls in hot paths where the type parameter causes a new monomorphized copy per type, increasing binary size and icache pressure. *Fix:* use concrete types or inline the implementation.
- **PERF-223** flags pool backing-array discard where a pooled `[]byte` is truncated and re-grown instead of reused, because the backing array is thrown away. *Fix:* `Reset()` the slice length without discarding the capacity.
- **PERF-224** flags recursive tree walks on hot paths that can overflow the stack on deep inputs. *Fix:* convert recursion to iterative traversal with an explicit stack.

## 3 — Concurrency & Goroutines

- PERF-131, PERF-158, PERF-170, PERF-181, PERF-182, PERF-195, PERF-198 above also cover concurrency patterns.

## Per-Detector Development

See [`documents/perf-detector-development.md`](./perf-detector-development.md) for the step-by-step guide to adding a new PERF detector.
