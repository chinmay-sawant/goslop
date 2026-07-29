package perf

import "github.com/chinmay/codehound/internal/rules"

// Catalogue metadata for PERF batch 2 (PERF-61..PERF-111, excluding PERF-104).

var (
	MetaPERF61 = rules.Meta(
		"PERF-61",
		"Gin Static Handler No Cache Headers",
		"Serves static files via gin.Static or router.StaticFS without Cache-Control, ETag, or Last-Modified headers, forcing redundant client fetches.",
		rules.SeverityMedium,
		nil,
		"Add explicit Cache-Control / ETag headers via c.Header() before serving static content.",
	)
	MetaPERF62 = rules.Meta(
		"PERF-62",
		"Gin Route Parameter Parsing In Middleware",
		"Parses complex route parameters (e.g., UUIDs, JSON) inside middleware that runs for every request, even for routes that do not need them.",
		rules.SeverityMedium,
		nil,
		"Parse path parameters once at registration time, not in middleware on every request.",
	)
	MetaPERF63 = rules.Meta(
		"PERF-63",
		"Gin Custom Validator Allocation",
		"Registers or invokes custom validator factories on every request instead of initializing them once via binding.Validator.Engine().",
		rules.SeverityMedium,
		nil,
		"Cache binding.Validator.Engine() in a package-level variable and reuse it across requests.",
	)
	MetaPERF64 = rules.Meta(
		"PERF-64",
		"Gin c.Copy Not Called In Goroutine",
		"Spawns a goroutine that uses *gin.Context directly. The context's underlying objects are pooled and will be recycled, causing data races and corruption.",
		rules.SeverityMedium,
		nil,
		"Call c.Copy() before launching the goroutine so the context survives the request lifetime.",
	)
	MetaPERF65 = rules.Meta(
		"PERF-65",
		"Gin ShouldBind In Middleware Chain",
		"Calls c.ShouldBind in a middleware that runs before route dispatch, performing the same parse work for routes that may not need it.",
		rules.SeverityMedium,
		nil,
		"Move c.ShouldBindJSON out of shared middleware into the leaf handler that owns the payload.",
	)
	MetaPERF66 = rules.Meta(
		"PERF-66",
		"Gin Deep Middleware Chain",
		"Builds a deep middleware chain per route group where many middlewares can be merged or composed, adding per-request frame overhead.",
		rules.SeverityMedium,
		nil,
		"Flatten the middleware chain or merge small middlewares into a single function.",
	)
	MetaPERF67 = rules.Meta(
		"PERF-67",
		"Gin Recovery Middleware Disabled",
		"Removes gin.Recovery() in production, allowing panics to crash workers. While strictly a reliability issue, repeated panics under load have severe performance impact.",
		rules.SeverityMedium,
		nil,
		"Add gin.Recovery() (or gin.RecoveryWithWriter) to gin.New() to avoid panics taking down the process.",
	)
	MetaPERF68 = rules.Meta(
		"PERF-68",
		"Gin Logger Middleware In Production Hot Path",
		"Leaves gin.Logger() or gin.LoggerWithFormatter in production, writing a log line per request with synchronous I/O on the hot path.",
		rules.SeverityMedium,
		nil,
		"Replace gin.Logger() with gin.LoggerWithConfig(gin.LoggerConfig{Output: io.Discard, ...}) in production.",
	)
	MetaPERF69 = rules.Meta(
		"PERF-69",
		"Gin c.Writer Buffer Flush Missed",
		"Streams a large response with c.Writer without periodic flush, increasing time-to-first-byte for clients.",
		rules.SeverityMedium,
		nil,
		"Call c.Writer.Flush() (or c.Writer.Flush()) after c.Writer.Write to flush the response in chunks.",
	)
	MetaPERF70 = rules.Meta(
		"PERF-70",
		"Gin Handler Goroutine Leak",
		"Spawns goroutines inside a Gin handler without ensuring they finish before the response is returned, leading to goroutine leaks and resource retention.",
		rules.SeverityMedium,
		nil,
		"Tie the goroutine to a sync.WaitGroup, a done channel, or c.Request.Context() so it cannot outlive the request.",
	)
	MetaPERF71 = rules.Meta(
		"PERF-71",
		"GORM N+1 Query Pattern",
		"Loads a list of parent records and then queries related records one by one inside a loop, producing N+1 round-trips to the database.",
		rules.SeverityMedium,
		nil,
		"Use db.Preload(\"Orders\") (or db.Joins) before the iteration, or batch the query with WHERE user_id IN (?).",
	)
	MetaPERF72 = rules.Meta(
		"PERF-72",
		"GORM Transaction Per Request",
		"Wraps every request in a database transaction, adding BEGIN/COMMIT round-trip overhead even for read-only or single-statement paths.",
		rules.SeverityMedium,
		nil,
		"Drop the transaction wrapper when the work is a single read or single statement.",
	)
	MetaPERF73 = rules.Meta(
		"PERF-73",
		"GORM Preload Missing For Relations",
		"Queries a parent record and accesses a relation in code paths, triggering lazy loads on the hot path. Preload or Joins would be cheaper.",
		rules.SeverityMedium,
		nil,
		"Call db.Preload(\"Relation\") on the parent query so GORM hydrates the relation in one round trip.",
	)
	MetaPERF74 = rules.Meta(
		"PERF-74",
		"GORM Select All Columns",
		"Uses db.Find or unscoped queries without Select(), pulling all columns including large text or blob fields that the response does not need.",
		rules.SeverityMedium,
		nil,
		"Add .Select(\"id, name, email\") to project only the columns the handler actually returns.",
	)
	MetaPERF75 = rules.Meta(
		"PERF-75",
		"GORM Session Not Reused",
		"Configures GORM session options (logger, dry run, prepared statement) per call instead of once via Session or a shared session config.",
		rules.SeverityMedium,
		nil,
		"Hoist the gorm.Session config to a package-level var and reuse it via db.WithContext(...).",
	)
	MetaPERF76 = rules.Meta(
		"PERF-76",
		"GORM Create In Loop",
		"Inserts records one at a time with db.Create inside a loop. Bulk insert with db.CreateInBatches or a single multi-value INSERT is significantly faster.",
		rules.SeverityMedium,
		nil,
		"Use db.CreateInBatches(rows, 100) (or insertBuilder) instead of calling db.Create per row.",
	)
	MetaPERF77 = rules.Meta(
		"PERF-77",
		"GORM Save Vs Update Misuse",
		"Uses db.Save for partial updates, which performs an UPSERT (UPDATE all columns or INSERT). db.Update or db.Updates with explicit fields is more efficient for partial updates.",
		rules.SeverityMedium,
		nil,
		"Use db.Update(\"field\", value) or db.Updates(map[string]any{...}) when only a subset of fields changes.",
	)
	MetaPERF78 = rules.Meta(
		"PERF-78",
		"GORM Raw Query Missing Index Hint",
		"Builds raw SQL with db.Raw or db.Exec on columns that lack indexes, leading to sequential scans on large tables.",
		rules.SeverityMedium,
		nil,
		"Add the corresponding index in the migration file (or use FORCE INDEX in the query) to back the WHERE/ORDER BY clause.",
	)
	MetaPERF79 = rules.Meta(
		"PERF-79",
		"GORM Connection Pool Exhaustion",
		"Misses db.SetMaxOpenConns, db.SetMaxIdleConns, or db.SetConnMaxLifetime configuration, leading to pool exhaustion or stale connections.",
		rules.SeverityMedium,
		nil,
		"Call db.SetMaxOpenConns / SetMaxIdleConns / SetConnMaxLifetime once at startup.",
	)
	MetaPERF80 = rules.Meta(
		"PERF-80",
		"GORM Pluck Returns Large Slice",
		"Uses db.Pluck to load a single column from a large table into a slice in memory, which can balloon memory usage.",
		rules.SeverityMedium,
		nil,
		"Add .Limit(N) (or chunk via FindInBatches) to bound the slice Pluck returns.",
	)
	MetaPERF81 = rules.Meta(
		"PERF-81",
		"sqlx SelectIn Unbounded Slice",
		"Builds an IN clause with a large slice in sqlx, generating a query so large it hits DB parameter limits or performs poorly.",
		rules.SeverityMedium,
		nil,
		"Chunk the slice into groups of 100â500 and run multiple db.Select queries.",
	)
	MetaPERF82 = rules.Meta(
		"PERF-82",
		"sqlx StructScan In Loop",
		"Calls rows.StructScan inside a loop without amortizing reflection. Caching struct scan targets or using Get/SELECT INTO is cheaper.",
		rules.SeverityMedium,
		nil,
		"Use rows.StructScan with a small destination struct, or use sqlx.Select to scan into a preallocated slice.",
	)
	MetaPERF83 = rules.Meta(
		"PERF-83",
		"sqlx MapScan Allocation",
		"Uses MapScan in hot paths where typed scanning or reuse would be cheaper, since MapScan always allocates a new map per row.",
		rules.SeverityMedium,
		nil,
		"Pre-declare a map[string]any with the expected columns, or use rows.Scan with explicit destinations.",
	)
	MetaPERF84 = rules.Meta(
		"PERF-84",
		"sqlx Transaction Misuse",
		"Mixes Begin/Commit patterns with sqlx incorrectly, holding a transaction longer than needed or nesting transactions without savepoints.",
		rules.SeverityMedium,
		nil,
		"Drop the transaction wrapper for single-statement work, or batch the work inside a shorter transaction.",
	)
	MetaPERF85 = rules.Meta(
		"PERF-85",
		"sqlx Named Query Binding Cost",
		"Uses sqlx.Named or sqlx.In repeatedly with large parameter sets, which builds query strings and argument lists per call.",
		rules.SeverityMedium,
		nil,
		"Pre-build the named query once and reuse it inside the loop (sqlx.NamedExec with a cached statement).",
	)
	MetaPERF86 = rules.Meta(
		"PERF-86",
		"Echo c.JSON Marshal In Hot Path",
		"Calls c.JSON inside Echo handlers with on-the-fly marshaling per request, missing pooled encoder opportunities.",
		rules.SeverityMedium,
		nil,
		"Batch responses with c.Stream, or reuse a json.Encoder with a sync.Pool to avoid reallocation.",
	)
	MetaPERF87 = rules.Meta(
		"PERF-87",
		"Echo Bind Without Validation Reuse",
		"Calls c.Bind on every request with full validation enabled, even for trusted internal routes that only need parsing.",
		rules.SeverityMedium,
		nil,
		"Skip full validation in trusted paths, or share a custom echo.Binder across handlers.",
	)
	MetaPERF88 = rules.Meta(
		"PERF-88",
		"Echo Static Handler Missing Cache",
		"Serves static files via echo.Static without Cache-Control, ETag, or Last-Modified headers, forcing redundant client fetches.",
		rules.SeverityMedium,
		nil,
		"Add Cache-Control / ETag headers to e.Static or the static middleware.",
	)
	MetaPERF89 = rules.Meta(
		"PERF-89",
		"Echo Middleware Allocation",
		"Defines middleware that allocates large structures on every request, even when the downstream handler does not need them.",
		rules.SeverityMedium,
		nil,
		"Move heavy parsing (make / json.Unmarshal) out of middleware or wrap the middleware behind sync.Once.",
	)
	MetaPERF90 = rules.Meta(
		"PERF-90",
		"Echo Context Store Growth",
		"Stores unbounded keys in c.Set() across the request lifetime, retaining memory until the context is recycled.",
		rules.SeverityMedium,
		nil,
		"Use c.Set with small scalar values (user_id, request_id, trace_id) and propagate context instead of large payloads.",
	)
	MetaPERF91 = rules.Meta(
		"PERF-91",
		"Fiber fasthttp Pooling Miss",
		"Allocates per-request buffers or objects in Fiber handlers without using fasthttp's pooling primitives like *fasthttp.RequestCtx byte slices or sync.Pool.",
		rules.SeverityMedium,
		nil,
		"Use fasthttp's bytebufferpool or c.Request.BodyStream() to avoid per-request allocations.",
	)
	MetaPERF92 = rules.Meta(
		"PERF-92",
		"Fiber Ctx Allocation Misuse",
		"Stores request-scoped data in a *fiber.Ctx field or returns the ctx from a goroutine, missing the fact that fiber.Ctx is reused per request.",
		rules.SeverityMedium,
		nil,
		"Copy needed fields out of c before launching the goroutine, or use c.UserContext() to scope the lifetime.",
	)
	MetaPERF93 = rules.Meta(
		"PERF-93",
		"Fiber JSON Encoder Reuse Missed",
		"Allocates a new json.Encoder or json.Marshal call per response when a pooled encoder would be cheaper.",
		rules.SeverityMedium,
		nil,
		"Reuse a json.Encoder via sync.Pool, or stream responses with c.SendStream.",
	)
	MetaPERF94 = rules.Meta(
		"PERF-94",
		"Fiber Request Body Read Pattern",
		"Reads the Fiber request body in a way that triggers buffer copies, missing fasthttp's zero-copy PostArgs or BodyStream access.",
		rules.SeverityMedium,
		nil,
		"Use c.Body() / c.RequestBodyStream() as a []byte and avoid io.ReadAll on fasthttp streams.",
	)
	MetaPERF95 = rules.Meta(
		"PERF-95",
		"Fiber Middleware Chain Growth",
		"Builds a deep middleware chain per route group in Fiber, where fasthttp per-request overhead multiplies with chain depth.",
		rules.SeverityMedium,
		nil,
		"Consolidate overlapping middlewares or attach them once at the app level instead of nested groups.",
	)
	MetaPERF96 = rules.Meta(
		"PERF-96",
		"gRPC Stream Allocation Per Recv",
		"Allocates a new message struct per Recv in a gRPC client stream, where reusing a single message would reduce GC pressure.",
		rules.SeverityMedium,
		nil,
		"Allocate the message once outside the loop and call msg.Reset() between stream.RecvMsg calls.",
	)
	MetaPERF97 = rules.Meta(
		"PERF-97",
		"protobuf Marshal In Loop",
		"Marshals protobuf messages inside a loop without reusing marshaling buffers or using proto.MarshalOptions with deterministic state.",
		rules.SeverityMedium,
		nil,
		"Reuse a MarshalOptions or a pooled bytes.Buffer across iterations to avoid per-call allocation.",
	)
	MetaPERF98 = rules.Meta(
		"PERF-98",
		"go-redis Pipeline Missed",
		"Issues sequential redis commands in a loop, missing the chance to pipeline them into a single round-trip.",
		rules.SeverityMedium,
		nil,
		"Wrap the redis calls in rdb.Pipeline().Exec(ctx) or rdb.Pipelined(ctx, ...) to batch the round trips.",
	)
	MetaPERF99 = rules.Meta(
		"PERF-99",
		"Prometheus Label High Cardinality",
		"Defines a counter, gauge, or histogram with a high-cardinality label (user ID, request ID, full URL) that explodes time series storage.",
		rules.SeverityMedium,
		nil,
		"Replace high-cardinality labels (user_id, uuid, path) with low-cardinality aggregates (status, method, route).",
	)
	MetaPERF100 = rules.Meta(
		"PERF-100",
		"Cobra Command Subcommand Allocation",
		"Registers cobra subcommands with expensive Run hooks, large init logic, or per-invocation flag parsing that could be amortized.",
		rules.SeverityMedium,
		nil,
		"Move heavy RunE work into a function and reuse flag registration via a pre-built flag.FlagSet.",
	)
	MetaPERF101 = rules.Meta(
		"PERF-101",
		"HTTP Server Timeouts Not Configured",
		"An http.Server is created without ReadTimeout, WriteTimeout, or IdleTimeout configured. Without these timeouts, slow clients can hold connections indefinitely, exhausting goroutines and file descriptors under load. This is a Slowloris-style DoS vector in production Go services.",
		rules.SeverityMedium,
		nil,
		"Set ReadTimeout, WriteTimeout, and IdleTimeout on http.Server to bound request lifetimes.",
	)
	MetaPERF102 = rules.Meta(
		"PERF-102",
		"HTTP WriteHeader Called Multiple Times",
		"w.WriteHeader is called more than once on the same ResponseWriter in a single handler; only the first call takes effect.",
		rules.SeverityMedium,
		nil,
		"WriteHeader can only be called once per response; set the status via WriteHeader(status) before the first Write.",
	)
	MetaPERF103 = rules.Meta(
		"PERF-103",
		"HTTP Response Body Not Closed",
		"Failing to close http.Response.Body leaks TCP connections from the connection pool. Even if all bytes are read, the connection cannot be reused for keep-alive unless the body is fully consumed AND closed. Over time this causes connection pool exhaustion and file descriptor leaks.",
		rules.SeverityMedium,
		nil,
		"Always defer resp.Body.Close() after a successful client call (and drain the body when reusing connections).",
	)
	MetaPERF105 = rules.Meta(
		"PERF-105",
		"runtime.SetFinalizer On Hot Path Object",
		"Uses runtime.SetFinalizer on objects created frequently in hot paths. Objects with finalizers require an extra GC cycle to collect. Each finalizer-eligible object is resurrected for the finalizer call, then collected on the next cycle, doubling GC work.",
		rules.SeverityMedium,
		nil,
		"Prefer explicit Close/Release methods over runtime.SetFinalizer for hot-path objects.",
	)
	MetaPERF106 = rules.Meta(
		"PERF-106",
		"sync.Map Used For Write-Heavy Workload",
		"Uses sync.Map for workloads dominated by writes (Store) and updates. sync.Map is optimized for read-heavy, key-stable workloads with its internal read/dirty dual-map structure. Under frequent writes, a plain map with sync.Mutex outperforms sync.Map. Additionally, any map or sync.Map used as a cache without eviction bounds (entry cap, byte cap, or TTL) can cause unbounded memory growth under concurrent load, leading to process OOM. Always pair caches with an eviction strategy and validate memory usage under concurrent load, not only in-process micro-benchmarks.",
		rules.SeverityMedium,
		nil,
		"Replace sync.Map with a plain map guarded by a sync.Mutex when the workload is write-heavy; sync.Map's read/dirty dual-map structure only pays off for read-heavy, key-stable workloads. If used as a cache, add eviction bounds: entry cap, byte cap, or TTL to prevent unbounded growth under concurrent load.",
	)
	MetaPERF107 = rules.Meta(
		"PERF-107",
		"encoding/binary Write Or Read Inside Loop",
		"Uses encoding/binary.Write or binary.Read inside a loop. These functions use reflection to determine the size and layout of the target type on every call. For fixed-size types, manual BigEndian methods avoid reflection and are 10-20x faster.",
		rules.SeverityMedium,
		nil,
		"Hoist encoding/binary out of the loop or hand-encode fixed-size types with binary.BigEndian.",
	)
	MetaPERF108 = rules.Meta(
		"PERF-108",
		"sort.Search Repeated In Loop",
		"Calls sort.Search inside a loop with a function literal closure, allocating a new closure on every iteration. For repeated binary searches on the same sorted data, extract the predicate as a named function or pre-build a lookup structure.",
		rules.SeverityMedium,
		nil,
		"Hoist sort.Search out of the loop; if the search space changes per iteration, cache the index instead.",
	)
	MetaPERF109 = rules.Meta(
		"PERF-109",
		"Map Key Recomputed In Loop Without Caching",
		"Recomputes a map key inside a loop by calling an expensive function on every iteration when the key could be computed once and stored. Each recomputation adds CPU overhead proportional to the key computation cost times the loop count.",
		rules.SeverityMedium,
		nil,
		"Cache expensive map keys outside or once per iteration instead of recomputing them.",
	)
	MetaPERF110 = rules.Meta(
		"PERF-110",
		"sync.Pool Element Type Causes Allocation On Put",
		"Stores non-pointer value types in sync.Pool (e.g., bytes.Buffer, not *bytes.Buffer). Each Get and Put operation boxes the value into interface{}, causing a heap allocation. Store pointer types instead to avoid the interface boxing allocation per pool operation.",
		rules.SeverityMedium,
		nil,
		"Return a pointer type from sync.Pool's New function (e.g. *Foo) so Put does not box the value back into the pool.",
	)
	MetaPERF111 = rules.Meta(
		"PERF-111",
		"Range Over String Produces Rune Allocation",
		"Converting a string to []rune(s) before ranging allocates a full rune slice on the heap. Ranging over the string directly with for _, r := range s decodes UTF-8 incrementally without allocating, cutting memory usage by 4x.",
		rules.SeverityMedium,
		nil,
		"Range over the string directly (for _, r := range s) instead of []rune(s).",
	)
)
