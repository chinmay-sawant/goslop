package perf

import "github.com/chinmay/codehound/internal/rules"

// Catalogue metadata for PERF batch 3 (PERF-112..163 excluding PERF-116).

var (
	MetaPERF112 = rules.Meta(
		"PERF-112",
		"strings.ToLower Before Comparison Instead Of EqualFold",
		"Compares strings with strings.ToLower(s1) == strings.ToLower(s2), allocating two intermediate strings even if they differ on the first byte. strings.EqualFold compares character-by-character without allocation and returns early on mismatch, up to 100x faster.",
		rules.SeverityMedium,
		nil,
		"Prefer strings.EqualFold for case-insensitive equality instead of ToLower/ToUpper comparisons.",
	)
	MetaPERF113 = rules.Meta(
		"PERF-113",
		"Single-Case Select Statement Instead Of Channel Op",
		"Uses a select statement with a single case instead of a plain channel send or receive. The select adds goroutine scheduling overhead and a runtime.selectgo call. A direct channel operation compiles to fewer instructions.",
		rules.SeverityLow,
		nil,
		"Replace a single-case select with a direct channel send or receive.",
	)
	MetaPERF114 = rules.Meta(
		"PERF-114",
		"Manual Loop Copy Instead Of copy() Builtin",
		"Uses a for loop to copy elements from one slice to another index-by-index instead of the copy() builtin. copy() uses optimized memory routines (often memmove in assembly) that are 3-10x faster and handle memory overlap detection.",
		rules.SeverityMedium,
		nil,
		"Replace the manual for-range copy with the copy() builtin.",
	)
	MetaPERF115 = rules.Meta(
		"PERF-115",
		"strings.Compare Used For Equality Check",
		"Uses strings.Compare(a, b) == 0 to check string equality. strings.Compare performs a full three-way comparison before the result is checked; a direct == comparison short-circuits at the first differing byte and avoids the function call overhead.",
		rules.SeverityMedium,
		nil,
		"Use a == b (or !=) instead of strings.Compare(a, b) == 0.",
	)
	MetaPERF117 = rules.Meta(
		"PERF-117",
		"bytes.Compare Used For Equality Check",
		"Uses bytes.Compare(a, b) == 0 to check byte slice equality. bytes.Equal(a, b) returns a boolean directly after the first difference, short-circuiting faster than a full three-way comparison.",
		rules.SeverityMedium,
		nil,
		"Use bytes.Equal(a, b) instead of bytes.Compare(a, b) == 0.",
	)
	MetaPERF118 = rules.Meta(
		"PERF-118",
		"Unnecessary http.NewRequest For Simple Methods",
		"Uses http.NewRequest for simple GET/HEAD/POST requests when http.Get, http.Head, or http.Post would suffice. The convenience functions avoid manual request construction overhead when custom headers or body handling are not needed.",
		rules.SeverityLow,
		nil,
		"Use http.Get/http.Head only for simple no-header GET/HEAD; keep NewRequest when",
	)
	MetaPERF119 = rules.Meta(
		"PERF-119",
		"Multiple Separate Appends Instead Of Spread Concatenation",
		"Uses multiple sequential append calls to the same slice instead of a single variadic append. Each separate append may trigger independent capacity growth and reallocation; combining them into append(s, a, b, c) grows once if needed.",
		rules.SeverityMedium,
		nil,
		"Merge the consecutive append calls into a single variadic append, e.g. s = appen",
	)
	MetaPERF120 = rules.Meta(
		"PERF-120",
		"time.Now().Sub Instead Of time.Since",
		"Uses time.Now().Sub(t) instead of time.Since(t). Both have identical behavior, but time.Since reads more naturally. When time.Now() has already been called in the function, using the cached value avoids a second syscall.",
		rules.SeverityLow,
		nil,
		"Prefer time.Since(t) over time.Now().Sub(t).",
	)
	MetaPERF121 = rules.Meta(
		"PERF-121",
		"Struct Literal Instead Of Direct Type Conversion",
		"Manually copies struct fields between types with identical underlying structures instead of using a direct type conversion T(v). The conversion is compile-time checked and zero-cost, while manual copying risks bugs from missed fields.",
		rules.SeverityLow,
		nil,
		"Use a direct type conversion (T(x)) when the source and target structs have iden",
	)
	MetaPERF122 = rules.Meta(
		"PERF-122",
		"HasPrefix Followed By Slice Instead Of TrimPrefix",
		"Checks strings.HasPrefix(s, prefix) and then slices s[len(prefix):] to extract the remainder. strings.TrimPrefix(s, prefix) returns the string with the prefix removed if present in a single call with clearer intent.",
		rules.SeverityInfo, // B-tier (Rust severity_for_tier)
		nil,
		"Use strings.TrimPrefix(s, p) instead of HasPrefix + s[len(p):].",
	)
	MetaPERF123 = rules.Meta(
		"PERF-123",
		"Redundant make Argument With Zero Value",
		"Calls make([]T, 0, cap) or make(map[K]V, 0) with an explicit zero length or capacity argument. The zero value is the default; omitting it produces identical behavior with less code.",
		rules.SeverityMedium, // unclassified PERF → Medium (Rust)
		nil,
		"Omit redundant zero length/capacity arguments to make.",
	)
	MetaPERF124 = rules.Meta(
		"PERF-124",
		"strings.Replace With -1 Instead Of ReplaceAll",
		"Uses strings.Replace(s, old, new, -1) with the magic number -1 to replace all occurrences. Since Go 1.12, strings.ReplaceAll(s, old, new) is available and communicates intent more clearly with identical performance.",
		rules.SeverityLow,
		nil,
		"Use strings.ReplaceAll instead of strings.Replace(..., -1).",
	)
	MetaPERF125 = rules.Meta(
		"PERF-125",
		"Redundant nil Check Before append",
		"Checks if a slice is nil before appending to it. append works correctly on nil slices (treating them as zero-length), so the nil check is unnecessary and adds a branch the compiler may not optimize away.",
		rules.SeverityLow,
		nil,
		"Drop the `if s != nil` guard; append handles a nil slice by allocating a new bac",
	)
	MetaPERF126 = rules.Meta(
		"PERF-126",
		"Redundant http.CanonicalHeaderKey Call",
		"Calls http.CanonicalHeaderKey on a header key that is already in canonical form. Methods like Header.Get and Header.Set already canonicalize internally, making the explicit call redundant.",
		rules.SeverityLow,
		nil,
		"Skip http.CanonicalHeaderKey for already-canonical header names; Header.Get/Set canonicalize.",
	)
	MetaPERF127 = rules.Meta(
		"PERF-127",
		"Unnecessary fmt.Sprintf In Log Call",
		"Wraps a static string in fmt.Sprintf before passing it to a log function. The log package Printf-style methods (log.Printf, log.Fatalf) handle formatting internally; wrapping in Sprintf first allocates an intermediate string.",
		rules.SeverityMedium,
		nil,
		"Pass the format string directly to log.Printf; do not wrap with fmt.Sprintf first.",
	)
	MetaPERF128 = rules.Meta(
		"PERF-128",
		"Multiple Independent Appends Can Be Combined",
		"Performs append on a slice in multiple consecutive statements when all appended values are available at once. Each separate append may independently grow the backing array. Combining into a single variadic append reduces amortized growth overhead.",
		rules.SeverityMedium,
		nil,
		"Merge the 3+ consecutive append calls into a single variadic append; each separa",
	)
	MetaPERF129 = rules.Meta(
		"PERF-129",
		"Range Loop Copies Value When Only Index Needed",
		"Ranges over a slice with for _, v := range and uses only the index or ignores v entirely. Even when the value is discarded, the runtime still copies each element. For large structs, the copy cost adds up. Use for i := range to iterate indices only.",
		rules.SeverityMedium,
		nil,
		"Use `for i := range xs` to skip copying the value when the loop body only needs",
	)
	MetaPERF130 = rules.Meta(
		"PERF-130",
		"Unnecessary Function Wrapper Adding Call Overhead",
		"Wraps a function call in an anonymous function literal for no reason: func() { f(args) }(). This adds a closure allocation and indirect call overhead. If wrapping is needed for argument adaptation, use a minimal wrapper.",
		rules.SeverityMedium,
		nil,
		"Inline the call: drop the `func() { ... }()` wrapper when the body is a single c",
	)
	MetaPERF131 = rules.Meta(
		"PERF-131",
		"sync.Mutex Used Where sync/atomic Suffices",
		"Uses sync.Mutex to protect a simple integer counter incremented in concurrent goroutines. atomic.AddInt64 is lock-free and significantly faster for simple arithmetic on integral types, avoiding lock/unlock overhead and goroutine blocking.",
		rules.SeverityMedium,
		nil,
		"Replace the mutex with sync/atomic for simple counter-style mutations; atomics c",
	)
	MetaPERF132 = rules.Meta(
		"PERF-132",
		"Goroutine Spawned Without Context Propagation",
		"Launches a goroutine that performs work (HTTP calls, DB queries) without receiving a context.Context from the parent. Without context propagation, the goroutine cannot be cancelled when the parent request times out, leading to wasted work and goroutine leaks.",
		rules.SeverityHigh,
		nil,
		"Pass the request context into the goroutine and use *Context APIs.",
	)
	MetaPERF133 = rules.Meta(
		"PERF-133",
		"sort.Slice Closure Allocation Inside Loop",
		"Uses sort.Slice with a closure inside a loop or repeatedly-called function, causing a heap allocation for the closure on every call. For basic types, sort.Ints/sort.Strings/sort.Float64s avoid the closure allocation entirely.",
		rules.SeverityMedium,
		nil,
		"Hoist sort.Slice out of the loop, or use sort.Sort with a sort.Interface type th",
	)
	MetaPERF134 = rules.Meta(
		"PERF-134",
		"Manual io.Read/Write Loop Instead Of io.Copy",
		"Implements a manual read-write loop with io.Reader and io.Writer instead of using io.Copy. io.Copy uses an internal 32KB buffer and handles EOF, short reads, and errors correctly. Manual implementations often use suboptimal buffer sizes.",
		rules.SeverityMedium,
		nil,
		"Use io.Copy instead of a manual Read/Write loop.",
	)
	MetaPERF135 = rules.Meta(
		"PERF-135",
		"encoding/gob Encoder Or Decoder Not Reused",
		"Creates a new gob.NewEncoder or gob.NewDecoder for every message on a persistent connection instead of reusing one encoder/decoder pair per stream. gob encoders maintain internal type info caches and buffers; reusing avoids re-transmitting type definitions.",
		rules.SeverityMedium,
		nil,
		"Hoist gob.NewEncoder/Decoder to a single instance created at startup; the constr",
	)
	MetaPERF137 = rules.Meta(
		"PERF-137",
		"runtime.Caller Used In Hot Path",
		"Calls runtime.Caller on hot paths to determine the calling function. runtime.Caller walks the stack, reading stack frames and program counter tables. On hot paths this adds measurable overhead; cache the result or use a compile-time constant.",
		rules.SeverityMedium,
		nil,
		"Avoid runtime.Caller on the hot path.",
	)
	MetaPERF138 = rules.Meta(
		"PERF-138",
		"runtime.Stack Used In Hot Path",
		"Calls runtime.Stack on hot paths to capture goroutine stack traces. This allocates a buffer and formats the stack, which is expensive. Stack traces should be captured only for error reporting or debugging, not on production hot paths.",
		rules.SeverityMedium,
		nil,
		"Avoid runtime.Stack on the hot path; capture stacks only for error paths.",
	)
	MetaPERF139 = rules.Meta(
		"PERF-139",
		"Closure Allocates Due To Variable Escape",
		"A closure captures multiple large outer variables, causing each to escape to the heap as separate allocations. Grouping captured variables into a single struct reduces N allocations to 1, since the struct itself escapes but fields are stored contiguously.",
		rules.SeverityMedium,
		nil,
		"Group captured variables into a struct to reduce escape allocations, or avoid the closure.",
	)
	MetaPERF140 = rules.Meta(
		"PERF-140",
		"debug.SetGCPercent Misuse Or Tuning In Production",
		"Calls debug.SetGCPercent with aggressive values (very low or -1 to disable) in production code without corresponding GOMEMLIMIT. Setting GCPercent too low increases GC frequency; setting to -1 disables GC entirely and risks OOM.",
		rules.SeverityHigh,
		nil,
		"Remove the debug.SetGCPercent(-1) call (it disables the GC assist entirely) or s",
	)
	MetaPERF141 = rules.Meta(
		"PERF-141",
		"URL.Query() Called Repeatedly Without Caching",
		"Calls r.URL.Query() or r.URL.Query().Get(key) multiple times in a single handler. Each call to URL.Query() parses the raw query string and allocates a new url.Values map. Cache the result in a local variable to parse once.",
		rules.SeverityMedium,
		nil,
		"Cache r.URL.Query() in a local variable at the top of the handler; subsequent ca",
	)
	MetaPERF142 = rules.Meta(
		"PERF-142",
		"http.MaxBytesReader Not Used For Untrusted Body",
		"Reads from an HTTP request body (r.Body) without wrapping it in http.MaxBytesReader. A malicious client can send unbounded data, causing the server to allocate memory until OOM. MaxBytesReader limits the body to a configured maximum.",
		rules.SeverityHigh,
		nil,
		"Wrap r.Body with http.MaxBytesReader before reading untrusted bodies.",
	)
	MetaPERF143 = rules.Meta(
		"PERF-143",
		"http.TimeoutHandler Not Used For Route-Level Timeouts",
		"Exposes HTTP endpoints without per-route timeout enforcement via http.TimeoutHandler, relying solely on server-level timeouts. Per-route timeouts allow different deadlines for fast endpoints vs slow ones.",
		rules.SeverityMedium,
		nil,
		"Wrap slow routes with http.TimeoutHandler for per-route deadlines.",
	)
	MetaPERF144 = rules.Meta(
		"PERF-144",
		"Content-Length Not Set In HTTP Response",
		"Serves HTTP responses with known body sizes but does not set Content-Length header. Without Content-Length, Go uses chunked transfer encoding, which adds per-chunk overhead. Setting Content-Length enables client progress bars and connection reuse.",
		rules.SeverityMedium,
		nil,
		"Set Content-Length when the response body size is known before Write.",
	)
	MetaPERF145 = rules.Meta(
		"PERF-145",
		"http.Request.WithContext Allocation On Hot Path",
		"Calls r.WithContext(ctx) in middleware or request chains, which shallow-copies the entire http.Request struct (~2KB). Each WithContext call in the chain allocates a fresh copy. Prefer passing context values through other mechanisms.",
		rules.SeverityMedium,
		nil,
		"Advisory micro-opt: r.WithContext allocates a new *http.Request by design; hoist",
	)
	MetaPERF146 = rules.Meta(
		"PERF-146",
		"fmt.Sprintf With Single String And No Verbs",
		"Uses fmt.Sprintf(\"%s\", s) where s is already a string. This parses the format string, builds an argument list, and copies s to a new string. The fmt overhead is 20-50x more expensive than using s directly.",
		rules.SeverityMedium,
		nil,
		"Use the string directly instead of fmt.Sprintf(\"%s\", s).",
	)
	MetaPERF147 = rules.Meta(
		"PERF-147",
		"strings.Replace Call Where ReplaceAll Suffices",
		"Uses strings.Replace with a negative count to replace all occurrences instead of strings.ReplaceAll. While performance is identical, the magic-number pattern obscures the 'replace all' intention and is a common source of bugs.",
		rules.SeverityLow,
		nil,
		"Use strings.ReplaceAll instead of strings.Replace with a negative count.",
	)
	MetaPERF148 = rules.Meta(
		"PERF-148",
		"Goroutine Leak Via Channel Send Without Guaranteed Receiver",
		"Sends to an unbuffered channel from a goroutine without guaranteeing a corresponding receive on all code paths. If the receiver exits early via error return or context cancellation, the sender blocks forever on the unbuffered channel, leaking the goroutine.",
		rules.SeverityMedium, // unclassified PERF → Medium (Rust)
		nil,
		"Ensure channel sends have a guaranteed receiver or use buffered/select with default.",
	)
	MetaPERF149 = rules.Meta(
		"PERF-149",
		"net.Conn Deadlines Not Set For Network Operations",
		"Performs network I/O (Read, Write) on net.Conn without setting deadlines. Without deadlines, a stalled peer can block the goroutine indefinitely. SetReadDeadline and SetWriteDeadline ensure I/O operations time out rather than hanging.",
		rules.SeverityHigh,
		nil,
		"Set a deadline before conn.Read / conn.Write with conn.SetReadDeadline / SetWrit",
	)
	MetaPERF150 = rules.Meta(
		"PERF-150",
		"Large Stack Frame From Local Variables",
		"Declares large arrays or many local variables that exceed the compiler stack frame threshold, preventing the function from being inlined. Functions that cannot be inlined add call overhead at every invocation.",
		rules.SeverityMedium,
		nil,
		"Avoid large stack arrays on hot paths; allocate on the heap or pool buffers.",
	)
	MetaPERF151 = rules.Meta(
		"PERF-151",
		"Non-Inlinable Function On Hot Path Due To Complexity",
		"A hot-path function exceeds the Go compiler inlining budget (~80 cost units as of Go 1.22). Each defer, for loop, and select adds to the budget. Functions that cross the threshold add ~5-15ns call overhead per invocation.",
		rules.SeverityMedium, // unclassified PERF → Medium (Rust)
		nil,
		"Simplify hot-path handlers to stay within the compiler inlining budget.",
	)
	MetaPERF152 = rules.Meta(
		"PERF-152",
		"Header Copy Via Manual Loop Instead Of Clone",
		"Copies HTTP headers by iterating over the source header map with a for-range loop and setting each key individually. http.Header.Clone() does a shallow copy in bulk and is significantly faster for header forwarding patterns.",
		rules.SeverityMedium,
		nil,
		"Use http.Header.Clone() instead of manually ranging and Set-ing headers.",
	)
	MetaPERF153 = rules.Meta(
		"PERF-153",
		"http.Cookie.String Called Repeatedly",
		"Calls http.Cookie.String() multiple times for the same cookie in a request handler. The String method re-serializes the cookie each call. Cache the serialized string or extract only the needed field.",
		rules.SeverityMedium,
		nil,
		"Cache cookie.String() when the serialized form is needed more than once.",
	)
	MetaPERF154 = rules.Meta(
		"PERF-154",
		"Unnecessary http.HandlerFunc Type Conversion",
		"Converts a function matching the http.HandlerFunc signature to http.HandlerFunc type and then passes it to http.Handle. Both Handle and HandleFunc accept compatible function types; the explicit conversion is a no-op.",
		rules.SeverityLow,
		nil,
		"Pass the handler function to HandleFunc, or register without redundant HandlerFunc conversion on Handle.",
	)
	MetaPERF155 = rules.Meta(
		"PERF-155",
		"http.ServeMux Pattern Without Method Restriction",
		"Registers route handlers that internally check r.Method with if/switch statements instead of using method-aware mux patterns or Go 1.22+ method routing. This duplicates HTTP method dispatching logic in every handler.",
		rules.SeverityLow,
		nil,
		"Prefer method-aware routing (Go 1.22+ patterns) over manual r.Method checks.",
	)
	MetaPERF156 = rules.Meta(
		"PERF-156",
		"Ranging Over String With Only Index Usage",
		"Ranges over a string with for i, ch := range s when only the index i is used (ch is underscore). The runtime still decodes UTF-8 runes into ch on every iteration. Use a standard for i := 0; i < len(s); i++ loop to iterate bytes only.",
		rules.SeverityMedium,
		nil,
		"Use `for i := range s` to skip UTF-8 decoding; the rune binding is only useful w",
	)
	MetaPERF157 = rules.Meta(
		"PERF-157",
		"Unnecessary Use Of fmt.Sprint With Single String",
		"Calls fmt.Sprint(x) where x is already a string. fmt.Sprint parses the format, boxes arguments, and copies the result, returning x unchanged. The call is 100x+ overhead compared to using x directly.",
		rules.SeverityMedium,
		nil,
		"Use the string directly instead of fmt.Sprint(string).",
	)
	MetaPERF158 = rules.Meta(
		"PERF-158",
		"Sorting Slice Of Basic Types With Closure",
		"Uses sort.Slice with a closure to sort a slice of basic types (int, string, float64) instead of using sort.Ints, sort.Strings, or sort.Float64s. sort.Slice allocates a closure per call; the type-specific functions are closure-free.",
		rules.SeverityMedium,
		nil,
		"Use slices.Sort for []int / []string / []float64; sort.Slice allocates a closure",
	)
	MetaPERF159 = rules.Meta(
		"PERF-159",
		"Using json.NewDecoder Instead Of json.Unmarshal For Buffered Data",
		"Uses io.ReadAll + json.Unmarshal to parse JSON from a reader, buffering the entire body before deserialization. json.NewDecoder(r).Decode(&v) streams directly from the reader, avoiding the intermediate byte slice allocation.",
		rules.SeverityMedium,
		nil,
		"Stream with json.NewDecoder(r).Decode instead of buffering then decoding a bytes.Reader.",
	)
	MetaPERF160 = rules.Meta(
		"PERF-160",
		"sql.Open Inside Request Handler",
		"Calls sql.Open inside a request handler or loop, creating a new connection pool per invocation instead of sharing a single *sql.DB singleton. Each sql.Open creates new connections, bypasses pool tuning, and leaks idle connections under load.",
		rules.SeverityHigh,
		nil,
		"Open the *sql.DB once at startup; do not call sql.Open inside handlers.",
	)
	MetaPERF161 = rules.Meta(
		"PERF-161",
		"rows.Err Not Checked After Iteration",
		"Iterates over sql.Rows with for rows.Next() but never checks rows.Err() after the loop exits. This silently drops scan errors and row-level failures, potentially missing database errors that occurred mid-iteration.",
		rules.SeverityHigh,
		nil,
		"Call rows.Err() after the rows.Next() loop.",
	)
	MetaPERF162 = rules.Meta(
		"PERF-162",
		"db.Ping Inside Request Handler",
		"Calls db.Ping or db.PingContext on every request to verify database liveness, adding an unnecessary round-trip per inbound call. Use a health-check endpoint or periodic goroutine with pings instead.",
		rules.SeverityMedium,
		nil,
		"Avoid db.Ping on every request; use a health-check path or background probe.",
	)
	MetaPERF163 = rules.Meta(
		"PERF-163",
		"db.Query Instead Of QueryRow For Single Row",
		"Uses db.Query or db.QueryContext for queries that return at most one row, then iterates with rows.Next() and rows.Scan(). db.QueryRow or db.QueryRowContext avoids the rowset iteration overhead for single-row queries.",
		rules.SeverityMedium,
		nil,
		"Use db.QueryRow for single-row queries; it handles rows.Close() for you.",
	)
)
