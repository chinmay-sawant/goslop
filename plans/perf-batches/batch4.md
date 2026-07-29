# PERF batch 4
Rules: 50

| ID | Rust file | Fix hint | Fixtures |
|---:|---|---|---|
| PERF-164 | `codehound/src/lang/go/detectors/perf/domains/stdlib_optimization/handler_limits.rs` |  | safe:PERF-164-safe.txt,vulnerable:PERF-164-vulnerable.txt |
| PERF-165 | `codehound/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/db_and_sql.rs` | Implement sql.Scanner on the custom type so rows.Scan can decode directly withou | safe:PERF-165-safe.txt,vulnerable:PERF-165-vulnerable.txt |
| PERF-166 | `codehound/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/db_and_sql.rs` | Use sql.NullString / sql.NullInt64 in the scan target; the database/sql package  | safe:PERF-166-safe.txt,vulnerable:PERF-166-vulnerable.txt |
| PERF-167 | `codehound/src/lang/go/detectors/perf/domains/concurrency.rs` |  | safe:PERF-167-safe.txt,vulnerable:PERF-167-vulnerable.txt |
| PERF-168 | `codehound/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/sync_and_locks.rs` | Send a pointer to the struct over the channel: ch <- &Large{...} to avoid copyin | safe:PERF-168-safe.txt,vulnerable:PERF-168-vulnerable.txt |
| PERF-169 | `codehound/src/lang/go/detectors/perf/domains/memory_gc.rs` |  | safe:PERF-169-safe.txt,vulnerable:PERF-169-vulnerable.txt |
| PERF-170 | `codehound/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/io_runtime.rs` | Hoist sync.Once out of hot paths; use a sync/atomic.Bool or a plain package-leve | safe:PERF-170-safe.txt,vulnerable:PERF-170-vulnerable.txt |
| PERF-171 | `codehound/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/sync_and_locks.rs` | Use a sync.Mutex for acquire/release; channels used as mutexes add an extra sche | safe:PERF-171-safe.txt,vulnerable:PERF-171-vulnerable.txt |
| PERF-172 | `codehound/src/lang/go/detectors/perf/domains/concurrency.rs` |  | safe:PERF-172-safe.txt,vulnerable:PERF-172-vulnerable.txt |
| PERF-173 | `codehound/src/lang/go/detectors/perf/domains/concurrency.rs` |  | safe:PERF-173-safe.txt,vulnerable:PERF-173-vulnerable.txt |
| PERF-174 | `codehound/src/lang/go/detectors/perf/domains/concurrency.rs` |  | safe:PERF-174-safe.txt,vulnerable:PERF-174-vulnerable.txt |
| PERF-175 | `codehound/src/lang/go/detectors/perf/domains/concurrency.rs` |  | safe:PERF-175-safe.txt,vulnerable:PERF-175-vulnerable.txt |
| PERF-176 | `codehound/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/io_runtime.rs` | Use io.CopyBuffer with a pooled *[]byte; io.Copy allocates a 32 KiB buffer per c | safe:PERF-176-safe.txt,vulnerable:PERF-176-vulnerable.txt |
| PERF-177 | `codehound/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/sort_and_search.rs` | Call os.ReadDir(name) to get []os.DirEntry; switch to os.ReadDir for new code. | safe:PERF-177-safe.txt,vulnerable:PERF-177-vulnerable.txt |
| PERF-178 | `codehound/src/lang/go/detectors/perf/domains/string_bytes.rs` |  | safe:PERF-178-safe.txt,vulnerable:PERF-178-vulnerable.txt |
| PERF-179 | `codehound/src/lang/go/detectors/perf/domains/string_bytes.rs` |  | safe:PERF-179-safe.txt,vulnerable:PERF-179-vulnerable.txt |
| PERF-180 | `codehound/src/lang/go/detectors/perf/domains/stdlib_optimization/io_and_runtime.rs` |  | safe:PERF-180-safe.txt,vulnerable:PERF-180-vulnerable.txt |
| PERF-181 | `codehound/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/db_and_sql.rs` | Call decoder.UseNumber() before Decode when the target struct has int/int64 fiel | safe:PERF-181-safe.txt,vulnerable:PERF-181-vulnerable.txt |
| PERF-182 | `codehound/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/db_and_sql.rs` | Pass an explicit buffer size to bufio.NewWriter(w, size) when the downstream Wri | safe:PERF-182-safe.txt,vulnerable:PERF-182-vulnerable.txt |
| PERF-183 | `codehound/src/lang/go/detectors/perf/domains/concurrency.rs` |  | safe:PERF-183-safe.txt,vulnerable:PERF-183-vulnerable.txt |
| PERF-184 | `codehound/src/lang/go/detectors/perf/domains/stdlib_optimization/io_and_runtime.rs` |  | safe:PERF-184-safe.txt,vulnerable:PERF-184-vulnerable.txt |
| PERF-185 | `codehound/src/lang/go/detectors/perf/domains/stdlib_optimization/io_and_runtime.rs` |  | safe:PERF-185-safe.txt,vulnerable:PERF-185-vulnerable.txt |
| PERF-186 | `codehound/src/lang/go/detectors/perf/domains/string_bytes.rs` |  | safe:PERF-186-safe.txt,vulnerable:PERF-186-vulnerable.txt |
| PERF-187 | `codehound/src/lang/go/detectors/perf/domains/stdlib_optimization/io_and_runtime.rs` |  | safe:PERF-187-safe.txt,vulnerable:PERF-187-vulnerable.txt |
| PERF-188 | `codehound/src/lang/go/detectors/perf/domains/stdlib_optimization/io_and_runtime.rs` |  | safe:PERF-188-safe.txt,vulnerable:PERF-188-vulnerable.txt |
| PERF-189 | `codehound/src/lang/go/detectors/perf/domains/stdlib_optimization/io_and_runtime.rs` |  | safe:PERF-189-safe.txt,vulnerable:PERF-189-vulnerable.txt |
| PERF-190 | `codehound/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/http_client.rs` |  | safe:PERF-190-safe.txt,vulnerable:PERF-190-vulnerable.txt |
| PERF-191 | `codehound/src/lang/go/detectors/perf/domains/memory_gc.rs` |  | safe:PERF-191-safe.txt,vulnerable:PERF-191-vulnerable.txt |
| PERF-192 | `codehound/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/maps_and_slices.rs` | Pass the expected size: make(map[K]V, len(src)) before the population loop to av | safe:PERF-192-safe.txt,vulnerable:PERF-192-vulnerable.txt |
| PERF-193 | `codehound/src/lang/go/detectors/perf/domains/concurrency.rs` |  | safe:PERF-193-safe.txt,vulnerable:PERF-193-vulnerable.txt |
| PERF-194 | `codehound/src/lang/go/detectors/perf/domains/concurrency.rs` |  | safe:PERF-194-safe.txt,vulnerable:PERF-194-vulnerable.txt |
| PERF-195 | `codehound/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/io_runtime.rs` | Return the error from the goroutine instead of calling log.Fatal; the caller dec | safe:PERF-195-safe.txt,vulnerable:PERF-195-vulnerable.txt |
| PERF-196 | `codehound/src/lang/go/detectors/perf/domains/stdlib_optimization/io_and_runtime.rs` |  | safe:PERF-196-safe.txt,vulnerable:PERF-196-vulnerable.txt |
| PERF-197 | `codehound/src/lang/go/detectors/perf/domains/stdlib_optimization/io_and_runtime.rs` |  | safe:PERF-197-safe.txt,vulnerable:PERF-197-vulnerable.txt |
| PERF-198 | `codehound/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/http_client.rs` |  | safe:PERF-198-safe.txt,vulnerable:PERF-198-vulnerable.txt |
| PERF-199 | `codehound/src/lang/go/detectors/perf/domains/stdlib_optimization/io_and_runtime.rs` |  | safe:PERF-199-safe.txt,vulnerable:PERF-199-vulnerable.txt |
| PERF-200 | `codehound/src/lang/go/detectors/perf/domains/stdlib_optimization/io_and_runtime.rs` |  | safe:PERF-200-safe.txt,vulnerable:PERF-200-vulnerable.txt |
| PERF-201 | `codehound/src/lang/go/detectors/perf/domains/stdlib_optimization/io_and_runtime.rs` |  | safe:PERF-201-safe.txt,vulnerable:PERF-201-vulnerable.txt |
| PERF-202 | `codehound/src/lang/go/detectors/perf/domains/stdlib_optimization/io_and_runtime.rs` |  | safe:PERF-202-safe.txt,vulnerable:PERF-202-vulnerable.txt |
| PERF-203 | `codehound/src/lang/go/detectors/perf/domains/string_bytes.rs` |  | safe:PERF-203-safe.txt,vulnerable:PERF-203-vulnerable.txt |
| PERF-204 | `codehound/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/cli_and_orm.rs` |  | safe:PERF-204-safe.txt,vulnerable:PERF-204-vulnerable.txt |
| PERF-205 | `codehound/src/lang/go/detectors/perf/domains/stdlib_optimization/io_and_runtime.rs` |  | safe:PERF-205-safe.txt,vulnerable:PERF-205-vulnerable.txt |
| PERF-206 | `codehound/src/lang/go/detectors/perf/domains/stdlib_optimization/io_and_runtime.rs` |  | safe:PERF-206-safe.txt,vulnerable:PERF-206-vulnerable.txt |
| PERF-207 | `codehound/src/lang/go/detectors/perf/domains/stdlib_optimization/io_and_runtime.rs` |  | safe:PERF-207-safe.txt,vulnerable:PERF-207-vulnerable.txt |
| PERF-209 | `codehound/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/cli_and_orm.rs` | Move shared initialization out of PersistentPreRunE into a sync.Once in the pare | safe:PERF-209-safe.txt,vulnerable:PERF-209-vulnerable.txt |
| PERF-210 | `codehound/src/lang/go/detectors/perf/domains/stdlib_optimization/io_and_runtime.rs` |  | safe:PERF-210-safe.txt,vulnerable:PERF-210-vulnerable.txt |
| PERF-211 | `codehound/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/cli_and_orm.rs` | Replace db.Not() / NOT IN with an explicit positive list (WHERE id IN (?)) so th | safe:PERF-211-safe.txt,vulnerable:PERF-211-vulnerable.txt |
| PERF-212 | `codehound/src/lang/go/detectors/perf/domains/stdlib_optimization/io_and_runtime.rs` |  | safe:PERF-212-safe.txt,vulnerable:PERF-212-vulnerable.txt |
| PERF-213 | `codehound/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/caching_and_allocation.rs` | Add an eviction boundary to the cache: limit entries (max N), limit retained byt | safe:PERF-213-safe.txt,vulnerable:PERF-213-vulnerable.txt |
| PERF-214 | `codehound/src/lang/go/detectors/perf/domains/general_perf/stdlib_misuse/caching_and_allocation.rs` | Remove volatile fields (pointer addresses, request IDs, coordinates) from the ca | safe:PERF-214-safe.txt,vulnerable:PERF-214-vulnerable.txt |
