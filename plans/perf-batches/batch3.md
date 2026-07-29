# PERF batch 3
Rules: 50

| ID | Rust file | Fix hint | Fixtures |
|---:|---|---|---|
| PERF-112 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/strings_bytes.rs` |  | safe:PERF-112-safe.txt,vulnerable:PERF-112-vulnerable.txt |
| PERF-113 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/ranges_and_types.rs` |  | safe:PERF-113-safe.txt,vulnerable:PERF-113-vulnerable.txt |
| PERF-114 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/ranges_and_types.rs` | Replace the manual for-range copy with the copy() builtin; copy() uses memmove a | safe:PERF-114-safe.txt,vulnerable:PERF-114-vulnerable.txt |
| PERF-115 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/strings_bytes.rs` |  | safe:PERF-115-safe.txt,vulnerable:PERF-115-vulnerable.txt |
| PERF-117 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/strings_bytes.rs` |  | safe:PERF-117-safe.txt,vulnerable:PERF-117-vulnerable.txt |
| PERF-118 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/http_client.rs` | Use http.Get/http.Head only for simple no-header GET/HEAD; keep NewRequest when  | safe:PERF-118-safe.txt,vulnerable:PERF-118-vulnerable.txt |
| PERF-119 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/strings_bytes.rs` | Merge the consecutive append calls into a single variadic append, e.g. s = appen | safe:PERF-119-safe.txt,vulnerable:PERF-119-vulnerable.txt |
| PERF-120 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/http_server.rs` |  | safe:PERF-120-safe.txt,vulnerable:PERF-120-vulnerable.txt |
| PERF-121 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/ranges_and_types.rs` | Use a direct type conversion (T(x)) when the source and target structs have iden | safe:PERF-121-safe.txt,vulnerable:PERF-121-vulnerable.txt |
| PERF-122 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/http_server.rs` |  | safe:PERF-122-safe.txt,vulnerable:PERF-122-vulnerable.txt |
| PERF-123 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/maps_and_slices.rs` |  | safe:PERF-123-safe.txt,vulnerable:PERF-123-vulnerable.txt |
| PERF-124 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/strings_bytes.rs` |  | safe:PERF-124-safe.txt,vulnerable:PERF-124-vulnerable.txt |
| PERF-125 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/strings_bytes.rs` | Drop the `if s != nil` guard; append handles a nil slice by allocating a new bac | safe:PERF-125-safe.txt,vulnerable:PERF-125-vulnerable.txt |
| PERF-126 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/http_server.rs` |  | safe:PERF-126-safe.txt,vulnerable:PERF-126-vulnerable.txt |
| PERF-127 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/http_server.rs` |  | safe:PERF-127-safe.txt,vulnerable:PERF-127-vulnerable.txt |
| PERF-128 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/maps_and_slices.rs` | Merge the 3+ consecutive append calls into a single variadic append; each separa | safe:PERF-128-safe.txt,vulnerable:PERF-128-vulnerable.txt |
| PERF-129 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/maps_and_slices.rs` | Use `for i := range xs` to skip copying the value when the loop body only needs  | safe:PERF-129-safe.txt,vulnerable:PERF-129-vulnerable.txt |
| PERF-130 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/strings_bytes.rs` | Inline the call: drop the `func() { ... }()` wrapper when the body is a single c | safe:PERF-130-safe.txt,vulnerable:PERF-130-vulnerable.txt |
| PERF-131 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/sync_and_locks.rs` | Replace the mutex with sync/atomic for simple counter-style mutations; atomics c | safe:PERF-131-safe.txt,vulnerable:PERF-131-vulnerable.txt |
| PERF-132 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/sync_and_locks.rs` | Pass the request context into the goroutine: go func() { db.QueryContext(ctx, .. | safe:PERF-132-safe.txt,vulnerable:PERF-132-vulnerable.txt |
| PERF-133 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/sort_and_search.rs` | Hoist sort.Slice out of the loop, or use sort.Sort with a sort.Interface type th | safe:PERF-133-safe.txt,vulnerable:PERF-133-vulnerable.txt |
| PERF-134 | `goslop/src/lang/go/detectors/perf/domains/memory_gc.rs` |  | safe:PERF-134-safe.txt,vulnerable:PERF-134-vulnerable.txt |
| PERF-135 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/sync_and_locks.rs` | Hoist gob.NewEncoder/Decoder to a single instance created at startup; the constr | safe:PERF-135-safe.txt,vulnerable:PERF-135-vulnerable.txt |
| PERF-137 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/io_runtime.rs` | Avoid runtime.Caller on the hot path; pass a stack index as a constant or use a  | safe:PERF-137-safe.txt,vulnerable:PERF-137-vulnerable.txt |
| PERF-138 | `goslop/src/lang/go/detectors/perf/domains/memory_gc.rs` |  | safe:PERF-138-safe.txt,vulnerable:PERF-138-vulnerable.txt |
| PERF-139 | `goslop/src/lang/go/detectors/perf/domains/memory_gc.rs` |  | safe:PERF-139-safe.txt,vulnerable:PERF-139-vulnerable.txt |
| PERF-140 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/sync_and_locks.rs` | Remove the debug.SetGCPercent(-1) call (it disables the GC assist entirely) or s | safe:PERF-140-safe.txt,vulnerable:PERF-140-vulnerable.txt |
| PERF-141 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/io_runtime.rs` | Cache r.URL.Query() in a local variable at the top of the handler; subsequent ca | safe:PERF-141-safe.txt,vulnerable:PERF-141-vulnerable.txt |
| PERF-142 | `goslop/src/lang/go/detectors/perf/domains/stdlib_optimization/handler_limits.rs` |  | safe:PERF-142-safe.txt,vulnerable:PERF-142-vulnerable.txt |
| PERF-143 | `goslop/src/lang/go/detectors/perf/domains/stdlib_optimization/handler_limits.rs` |  | safe:PERF-143-safe.txt,vulnerable:PERF-143-vulnerable.txt |
| PERF-144 | `goslop/src/lang/go/detectors/perf/domains/stdlib_optimization/handler_limits.rs` |  | safe:PERF-144-safe.txt,vulnerable:PERF-144-vulnerable.txt |
| PERF-145 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/http_client.rs` | Advisory micro-opt: r.WithContext allocates a new *http.Request by design; hoist | safe:PERF-145-safe.txt,vulnerable:PERF-145-vulnerable.txt |
| PERF-146 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/ranges_and_types.rs` |  | safe:PERF-146-safe.txt,vulnerable:PERF-146-vulnerable.txt |
| PERF-147 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/ranges_and_types.rs` |  | safe:PERF-147-safe.txt,vulnerable:PERF-147-vulnerable.txt |
| PERF-148 | `goslop/src/lang/go/detectors/perf/domains/concurrency.rs` |  | safe:PERF-148-safe.txt,vulnerable:PERF-148-vulnerable.txt |
| PERF-149 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/io_runtime.rs` | Set a deadline before conn.Read / conn.Write with conn.SetReadDeadline / SetWrit | safe:PERF-149-safe.txt,vulnerable:PERF-149-vulnerable.txt |
| PERF-150 | `goslop/src/lang/go/detectors/perf/domains/memory_gc.rs` |  | safe:PERF-150-safe.txt,vulnerable:PERF-150-vulnerable.txt |
| PERF-151 | `goslop/src/lang/go/detectors/perf/domains/memory_gc.rs` |  | safe:PERF-151-safe.txt,vulnerable:PERF-151-vulnerable.txt |
| PERF-152 | `goslop/src/lang/go/detectors/perf/domains/stdlib_optimization/handler_limits.rs` |  | safe:PERF-152-safe.txt,vulnerable:PERF-152-vulnerable.txt |
| PERF-153 | `goslop/src/lang/go/detectors/perf/domains/stdlib_optimization/handler_limits.rs` |  | safe:PERF-153-safe.txt,vulnerable:PERF-153-vulnerable.txt |
| PERF-154 | `goslop/src/lang/go/detectors/perf/domains/stdlib_optimization/handler_limits.rs` |  | safe:PERF-154-safe.txt,vulnerable:PERF-154-vulnerable.txt |
| PERF-155 | `goslop/src/lang/go/detectors/perf/domains/stdlib_optimization/handler_limits.rs` |  | safe:PERF-155-safe.txt,vulnerable:PERF-155-vulnerable.txt |
| PERF-156 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/sort_and_search.rs` | Use `for i := range s` to skip UTF-8 decoding; the rune binding is only useful w | safe:PERF-156-safe.txt,vulnerable:PERF-156-vulnerable.txt |
| PERF-157 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/ranges_and_types.rs` |  | safe:PERF-157-safe.txt,vulnerable:PERF-157-vulnerable.txt |
| PERF-158 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/sort_and_search.rs` | Use slices.Sort for []int / []string / []float64; sort.Slice allocates a closure | safe:PERF-158-safe.txt,vulnerable:PERF-158-vulnerable.txt |
| PERF-159 | `goslop/src/lang/go/detectors/perf/domains/string_bytes.rs` |  | safe:PERF-159-safe.txt,vulnerable:PERF-159-vulnerable.txt |
| PERF-160 | `goslop/src/lang/go/detectors/perf/domains/stdlib_optimization/handler_limits.rs` |  | safe:PERF-160-safe.txt,vulnerable:PERF-160-vulnerable.txt |
| PERF-161 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/io_runtime.rs` | Call rows.Err() after the rows.Next() loop to distinguish 'no more rows' from a  | safe:PERF-161-safe.txt,vulnerable:PERF-161-vulnerable.txt |
| PERF-162 | `goslop/src/lang/go/detectors/perf/domains/stdlib_optimization/handler_limits.rs` |  | safe:PERF-162-safe.txt,vulnerable:PERF-162-vulnerable.txt |
| PERF-163 | `goslop/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/io_runtime.rs` | Use db.QueryRow for single-row queries; it handles rows.Close() for you. | safe:PERF-163-safe.txt,vulnerable:PERF-163-vulnerable.txt |
