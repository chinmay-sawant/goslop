package perf

import "github.com/chinmay/codehound/internal/rules"

// Batch 4 metadata: PERF-164..214 (excluding PERF-208).

var (
	MetaPERF164 = rules.Meta(
		"PERF-164",
		"Missing Context In Database Calls",
		"Uses database/sql methods like db.Query, db.Exec, or db.Prepare without their Context variants in request-serving code. The non-context variants lose cancellation propagation and request-level timeout enforcement.",
		rules.SeverityMedium,
		nil,
		"",
	)
	MetaPERF165 = rules.Meta(
		"PERF-165",
		"Not Implementing sql.Scanner For Custom Types",
		"Scans database values into intermediate string or []byte variables and then parses them into a custom type, instead of implementing the sql.Scanner interface. The Scanner interface avoids the per-row intermediate allocation.",
		rules.SeverityMedium,
		nil,
		"Implement sql.Scanner on the custom type so rows.Scan can decode directly without manual extraction.",
	)
	MetaPERF166 = rules.Meta(
		"PERF-166",
		"database/sql Null Handling Without sql.Null Types",
		"Scans nullable database columns into pointer types (*string, *int) and dereferences them on every row. The sql.NullString, sql.NullInt64 types avoid heap-escaping and provide explicit Valid fields for null checking.",
		rules.SeverityLow,
		nil,
		"Use sql.NullString / sql.NullInt64 in the scan target; the database/sql package exposes typed null wrappers.",
	)
	MetaPERF167 = rules.Meta(
		"PERF-167",
		"WaitGroup.Add Inside Goroutine",
		"Calls wg.Add(1) inside the goroutine body instead of before the go statement. This creates a race condition where wg.Wait() can return before Add is executed, causing zero or negative WaitGroup counters.",
		rules.SeverityHigh,
		nil,
		"",
	)
	MetaPERF168 = rules.Meta(
		"PERF-168",
		"Large Struct Sent By Value Over Channel",
		"Sends large structs by value through channels instead of using pointers. Each channel operation copies the full struct; for structs with many fields, strings, or slices, the copy cost can be significant under throughput.",
		rules.SeverityMedium,
		nil,
		"Send a pointer to the struct over the channel: ch <- &Large{...} to avoid copying every field.",
	)
	MetaPERF169 = rules.Meta(
		"PERF-169",
		"atomic.Value Frequent Store Allocation",
		"Uses sync/atomic.Value for values frequently updated in hot paths. Each Store boxes the value into an interface{}, allocating on the heap. Typed atomic types (atomic.Int64, atomic.Pointer) or mutex are cheaper for frequent updates.",
		rules.SeverityMedium,
		nil,
		"",
	)
	MetaPERF170 = rules.Meta(
		"PERF-170",
		"sync.Once In Hot Function Path",
		"Places sync.Once.Do() inside a function called on every request or in a loop, adding an atomic check overhead even after initialization completes. Move the sync.Once to package level so the hot path never touches it after the first call.",
		rules.SeverityMedium,
		nil,
		"Hoist sync.Once out of hot paths; use a sync/atomic.Bool or a plain package-level var for cheap one-time init.",
	)
	MetaPERF171 = rules.Meta(
		"PERF-171",
		"Channel Used As Mutex",
		"Uses a buffered channel with capacity 1 as mutual exclusion (send-to-acquire, receive-to-release) instead of sync.Mutex. Channel-based mutual exclusion adds goroutine scheduling overhead and is harder to make deadlock-free.",
		rules.SeverityMedium,
		nil,
		"Use a sync.Mutex for acquire/release; channels used as mutexes add an extra scheduling hop and a heap-allocated channel struct.",
	)
	MetaPERF172 = rules.Meta(
		"PERF-172",
		"WaitGroup.Wait Blocking Serving Goroutine",
		"Calls sync.WaitGroup.Wait() in a request handler to wait for spawned goroutines, blocking the serving goroutine until all sub-tasks complete. Use a result channel or errgroup with context instead.",
		rules.SeverityHigh,
		nil,
		"",
	)
	MetaPERF173 = rules.Meta(
		"PERF-173",
		"time.Tick Not Stopped Causing Goroutine Leak",
		"Uses time.NewTicker in a long-running function without calling ticker.Stop(). The underlying timer goroutine continues firing after the ticker goes out of scope, leaking a goroutine. Always defer ticker.Stop().",
		rules.SeverityMedium,
		nil,
		"",
	)
	MetaPERF174 = rules.Meta(
		"PERF-174",
		"Closing Channel By Receiver",
		"Closes a channel from the receiver side rather than the sender. If the sender tries to send after the close, it panics with 'send on closed channel'. The sender should own the channel lifecycle.",
		rules.SeverityHigh,
		nil,
		"",
	)
	MetaPERF175 = rules.Meta(
		"PERF-175",
		"Buffered Channel Spinning On Receive",
		"Tight-loops on receiving from a buffered channel without a select or context deadline, consuming CPU when the channel is empty. A select with default + small sleep, or blocking receive with context cancellation, avoids CPU spin.",
		rules.SeverityMedium,
		nil,
		"",
	)
	MetaPERF176 = rules.Meta(
		"PERF-176",
		"io.Copy Without Buffer Reuse",
		"Uses io.Copy without providing a reusable buffer via io.CopyBuffer. io.Copy allocates a 32KB buffer internally on every call. When copying in a loop or per-request, providing a shared buffer avoids this repeated allocation.",
		rules.SeverityMedium,
		nil,
		"Use io.CopyBuffer with a pooled *[]byte; io.Copy allocates a 32 KiB buffer per call.",
	)
	MetaPERF177 = rules.Meta(
		"PERF-177",
		"os.File.Readdir Instead Of os.ReadDir",
		"Uses (*os.File).Readdir to list directory entries, which returns []os.FileInfo and stats every entry. os.ReadDir (Go 1.16+) returns []os.DirEntry, deferring the stat call and avoiding unnecessary filesystem syscalls.",
		rules.SeverityLow,
		nil,
		"Call os.ReadDir(name) to get []os.DirEntry; switch to os.ReadDir for new code.",
	)
	MetaPERF178 = rules.Meta(
		"PERF-178",
		"time.Format Instead Of time.AppendFormat",
		"Uses time.Format which allocates a new string on every call. When formatting many timestamps in a loop, time.AppendFormat with a reusable byte slice avoids per-call string allocation, reducing GC pressure.",
		rules.SeverityLow,
		nil,
		"",
	)
	MetaPERF179 = rules.Meta(
		"PERF-179",
		"strings.Replacer Not Used For Repeated Replace",
		"Calls strings.Replace or strings.ReplaceAll repeatedly in a loop with the same old/new pairs on different inputs. strings.NewReplacer compiles the replacement set once and amortizes the cost across calls.",
		rules.SeverityLow,
		nil,
		"",
	)
	MetaPERF180 = rules.Meta(
		"PERF-180",
		"encoding/csv Reader Per Row",
		"Creates a new csv.NewReader for each row or batch in a hot parsing path instead of reusing one reader instance. Each csv.NewReader allocates internal buffers; reusing the reader avoids repeated allocation.",
		rules.SeverityMedium,
		nil,
		"",
	)
	MetaPERF181 = rules.Meta(
		"PERF-181",
		"json.Decoder UseNumber Missing",
		"Uses json.NewDecoder without calling .UseNumber() when decoded JSON contains numeric values that will be immediately used as integers. Without UseNumber(), the decoder parses all numbers as float64, requiring lossy float-to-int conversion.",
		rules.SeverityMedium,
		nil,
		"Call decoder.UseNumber() before Decode when the target struct has int/int64 fields, to avoid silent float64 precision loss for big numbers.",
	)
	MetaPERF182 = rules.Meta(
		"PERF-182",
		"bufio.Writer Default Buffer Undersized",
		"Uses bufio.NewWriter with the default 4096-byte buffer when typical writes are significantly larger. This causes excessive flush operations and small syscalls. Setting a buffer size matching expected write chunks with bufio.NewWriterSize reduces overhead.",
		rules.SeverityMedium,
		nil,
		"Pass an explicit buffer size to bufio.NewWriter(w, size) when the downstream Write calls are larger than the default 4 KiB buffer.",
	)
	MetaPERF183 = rules.Meta(
		"PERF-183",
		"context.WithTimeout Inside Loop",
		"Creates a new context.WithTimeout or context.WithCancel inside a loop body, allocating a new timer and context value chain on every iteration. Hoist the context derivation outside the loop or use a single long-lived context.",
		rules.SeverityMedium,
		nil,
		"",
	)
	MetaPERF184 = rules.Meta(
		"PERF-184",
		"mime.TypeByExtension In Hot Path",
		"Calls mime.TypeByExtension in request-serving paths to determine Content-Type headers. This function performs a map lookup and string manipulation per call. Cache the result or use a constant map for known MIME types.",
		rules.SeverityMedium,
		nil,
		"",
	)
	MetaPERF185 = rules.Meta(
		"PERF-185",
		"http.DetectContentType In Request Handler",
		"Uses http.DetectContentType in request-serving handlers to sniff file content types. This function reads up to 512 bytes and performs pattern matching on every call. Cache the detected type or use explicit MIME type mappings.",
		rules.SeverityMedium,
		nil,
		"",
	)
	MetaPERF186 = rules.Meta(
		"PERF-186",
		"strings.Fields In Hot Parsing Path",
		"Uses strings.Fields or strings.FieldsFunc in hot parsing paths, allocating a slice of all substrings at once. For large inputs, a bufio.Scanner or manual index-based parsing avoids allocating all fields simultaneously.",
		rules.SeverityMedium,
		nil,
		"",
	)
	MetaPERF187 = rules.Meta(
		"PERF-187",
		"template.HTMLEscaper In Hot Path",
		"Calls template.HTMLEscaper or template.JSEscaper repeatedly in hot paths that build HTML or JavaScript fragments. These functions allocate a new string per call. Pre-escaping at startup or using html/template automatic escaping is more efficient.",
		rules.SeverityMedium,
		nil,
		"",
	)
	MetaPERF188 = rules.Meta(
		"PERF-188",
		"fmt.Sscanf In Hot Path",
		"Uses fmt.Sscanf or fmt.Fscanf in performance-sensitive parsing paths. These functions use reflection to parse according to format strings and are orders of magnitude slower than direct strconv or manual parsing.",
		rules.SeverityInfo, // B-tier
		nil,
		"",
	)
	MetaPERF189 = rules.Meta(
		"PERF-189",
		"HTTP Response Body Not Drained Before Close",
		"Calls resp.Body.Close() without fully consuming the response body first. HTTP/1.1 persistent connections cannot be reused unless the entire body is read. This causes new TCP connections and TLS handshakes for every request.",
		rules.SeverityHigh,
		nil,
		"",
	)
	MetaPERF190 = rules.Meta(
		"PERF-190",
		"HTTP Client Missing Timeout",
		"Creates an http.Client without setting the Timeout field, leaving outbound HTTP requests susceptible to hanging indefinitely. A hung request ties up a goroutine, and under load this cascade-exhausts goroutines and memory.",
		rules.SeverityHigh,
		nil,
		"",
	)
	MetaPERF191 = rules.Meta(
		"PERF-191",
		"Slice Of Pointers For Small Structs",
		"Uses []*SmallStruct when []SmallStruct would provide better cache locality and avoid pointer-chasing overhead. Each pointer dereference causes a cache miss when iterating, and pointed-to values are individually heap-allocated and GC-scanned.",
		rules.SeverityLow,
		nil,
		"",
	)
	MetaPERF192 = rules.Meta(
		"PERF-192",
		"Map Without Size Hint",
		"Creates a map with make(map[K]V) without a capacity hint when the approximate number of entries is known at allocation time. Without a hint, the map resizes its hash table multiple times as entries are added, causing rehashing and allocation churn.",
		rules.SeverityMedium, // unclassified PERF → Medium
		nil,
		"Pass the expected size: make(map[K]V, len(src)) before the population loop to avoid map growth.",
	)
	MetaPERF193 = rules.Meta(
		"PERF-193",
		"Not Resetting Timer In Loop",
		"Creates a new time.NewTimer on every loop iteration instead of calling .Reset() on an existing timer. Even though Go 1.23+ improved Reset safety, creating a new timer per iteration allocates and registers a new runtime timer object.",
		rules.SeverityMedium,
		nil,
		"",
	)
	MetaPERF194 = rules.Meta(
		"PERF-194",
		"Using time.Sleep For Polling",
		"Uses time.Sleep in a for loop as a polling mechanism instead of event-driven patterns like time.Ticker, fsnotify, or channel-based notifications. Sleep-based polling wastes CPU on wake-ups and adds latency equal to the sleep interval.",
		rules.SeverityMedium,
		nil,
		"",
	)
	MetaPERF195 = rules.Meta(
		"PERF-195",
		"log.Fatal Or log.Panic In Goroutine",
		"Uses log.Fatal or log.Panic inside a goroutine that is not the main goroutine. log.Fatal calls os.Exit(1), terminating the entire process. log.Panic unwinds only the calling goroutine, leaving others in an undefined state.",
		rules.SeverityMedium, // unclassified PERF → Medium
		nil,
		"Return the error from the goroutine instead of calling log.Fatal; the caller decides whether to terminate the process.",
	)
	MetaPERF196 = rules.Meta(
		"PERF-196",
		"JWT Token Parsing Per Handler",
		"Parses and validates a JWT token in individual route handlers instead of in a shared authentication middleware. This duplicates cryptographic signature verification work on every handler in the chain.",
		rules.SeverityMedium,
		nil,
		"",
	)
	MetaPERF197 = rules.Meta(
		"PERF-197",
		"Multiple io.ReadAll On Request Body",
		"Reads the HTTP request body with io.ReadAll or json.NewDecoder in multiple places (middleware then handler). The body is a read-once stream; the second read gets EOF. Re-reading requires restoring via io.NopCloser with bytes.Buffer.",
		rules.SeverityMedium,
		nil,
		"",
	)
	MetaPERF198 = rules.Meta(
		"PERF-198",
		"Content-Type Check With strings.Contains",
		"Checks the Content-Type header using strings.Contains(header, \"json\") instead of a direct comparison or switch. strings.Contains performs a full substring scan and may false-match (e.g., application/vnd.api+json matches both 'json' and 'xml').",
		rules.SeverityLow,
		nil,
		"",
	)
	MetaPERF199 = rules.Meta(
		"PERF-199",
		"Session Store Lookup Per Handler",
		"Looks up the user session from a session store (Redis, database, cookie) on every individual route handler, instead of extracting it once in middleware. This multiplies network round-trips by the number of handlers in the chain.",
		rules.SeverityMedium,
		nil,
		"",
	)
	MetaPERF200 = rules.Meta(
		"PERF-200",
		"Middleware Ordering Penalty",
		"Places expensive middleware (authentication, body parsing, rate limiting) before cheap early-reject middleware (CORS preflight, request size check, IP allowlist). An early-rejected request still pays the cost of the expensive middleware.",
		rules.SeverityMedium,
		nil,
		"",
	)
	MetaPERF201 = rules.Meta(
		"PERF-201",
		"CORS Preflight Handler Allocation",
		"Implements a custom CORS preflight handler that allocates response headers or performs work for OPTIONS requests. Using a well-known CORS middleware that short-circuits on preflight is cheaper and well-tested.",
		rules.SeverityMedium,
		nil,
		"",
	)
	MetaPERF202 = rules.Meta(
		"PERF-202",
		"json.Marshal Indent In Production Handler",
		"Uses json.MarshalIndent or json.Encoder.SetIndent in production handlers to pretty-print JSON responses. Indentation adds whitespace (increasing response size) and requires extra allocation. Pretty-printing should be development-only.",
		rules.SeverityMedium,
		nil,
		"",
	)
	MetaPERF203 = rules.Meta(
		"PERF-203",
		"net.IP.String Repeated In Hot Path",
		"Calls net.IP.String() repeatedly on the same IP address in hot paths. Each call allocates a new string. Cache the string form once or use a formatting helper that writes to a reusable buffer.",
		rules.SeverityLow,
		nil,
		"",
	)
	MetaPERF204 = rules.Meta(
		"PERF-204",
		"GORM Updates With Map Without Select",
		"Uses db.Model(&record).Updates(map[string]interface{}{...}) without a preceding .Select() clause. This updates all fields in the map unconditionally, including zero-value fields, generating larger UPDATE statements than needed.",
		rules.SeverityMedium,
		nil,
		"Add a .Select(\"col1\", \"col2\") call before db.Updates(map) so GORM only writes the columns you intend.",
	)
	MetaPERF205 = rules.Meta(
		"PERF-205",
		"GORM Pagination Without Count Optimization",
		"Calls db.Count(&total) followed by db.Offset(offset).Limit(limit).Find(&records) on every paginated list request. For large tables, OFFSET-based pagination degrades linearly; keyset pagination avoids this performance cliff.",
		rules.SeverityMedium,
		nil,
		"",
	)
	MetaPERF206 = rules.Meta(
		"PERF-206",
		"sqlx Unsafe Without Known Input",
		"Uses sqlx.Unsafe to disable compile-time SQL name binding resolution when the query string is built from runtime input. Unsafe mode should only be used with static queries.",
		rules.SeverityMedium,
		nil,
		"",
	)
	MetaPERF207 = rules.Meta(
		"PERF-207",
		"Fiber ctx.SendFile Without Caching",
		"Uses c.SendFile in Fiber handlers without setting cache headers (Cache-Control, ETag, Last-Modified). Each request re-reads and re-transmits the file, wasting disk I/O and network bandwidth for static assets.",
		rules.SeverityMedium,
		nil,
		"",
	)
	MetaPERF209 = rules.Meta(
		"PERF-209",
		"Cobra PersistentPreRun In Every Command",
		"Defines heavy PersistentPreRun or PersistentPostRun hooks on parent Cobra commands that execute for every subcommand invocation, even when the subcommand does not need the setup. Use subcommand-specific hooks or lazy initialization.",
		rules.SeverityMedium,
		nil,
		"Move shared initialization out of PersistentPreRunE into a sync.Once in the parent command, or into a setup function called once at startup.",
	)
	MetaPERF210 = rules.Meta(
		"PERF-210",
		"go-redis KEYS Command In Application Code",
		"Uses the Redis KEYS command via go-redis in application logic, which scans the entire keyspace and blocks the Redis server. For production, SCAN (rdb.Scan) iterates incrementally without blocking.",
		rules.SeverityMedium,
		nil,
		"",
	)
	MetaPERF211 = rules.Meta(
		"PERF-211",
		"GORM Not In Select Clause",
		"Uses GORM's Not() condition on indexed queries where NOT IN, NOT LIKE, or != conditions often defeat index usage. GORM translates these without warning, potentially causing full table scans on large tables.",
		rules.SeverityMedium,
		nil,
		"Replace db.Not() / NOT IN with an explicit positive list (WHERE id IN (?)) so the query planner can use the index.",
	)
	MetaPERF212 = rules.Meta(
		"PERF-212",
		"GORM Find Without Limit On Large Table",
		"Calls db.Find(&records) without Limit() on a table that may contain many rows. This loads all rows into memory, potentially causing OOM for unbounded tables. Always paginate or limit queries on tables that can grow large.",
		rules.SeverityMedium,
		nil,
		"",
	)
	MetaPERF213 = rules.Meta(
		"PERF-213",
		"Cache Without Eviction or Bounding",
		"A map, sync.Map, or other cache data structure is allocated at package scope or reused across requests without any eviction mechanism (entry count cap, byte size cap, or TTL). Under concurrent load, an unbounded cache grows until it exhausts available memory, causing process OOM or server hang. Every cache persisted across requests must have a reachable eviction boundary.",
		rules.SeverityMedium,
		nil,
		"Add an eviction boundary to the cache: limit entries (max N), limit retained bytes (max M), or add TTL-based expiry so the cache cannot grow unbounded under load.",
	)
	MetaPERF214 = rules.Meta(
		"PERF-214",
		"Cache Key Includes Volatile Request-Scoped Fields",
		"A cache key incorporates fields that change on every request or invocation (e.g., pointer addresses, request IDs, coordinates, iteration variables). When the key changes every call, the cache has zero hit rate while still accumulating entries, wasting memory. Every cache key must consist only of fields that are stable across repeated identical calls.",
		rules.SeverityMedium,
		nil,
		"Remove volatile fields (pointer addresses, request IDs, coordinates) from the cache key. Key only on fields that are stable across repeated identical calls.",
	)
)
