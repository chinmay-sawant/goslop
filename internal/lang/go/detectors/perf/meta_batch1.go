package perf

import "github.com/chinmay/goslop/internal/rules"

// Batch 1 metadata (PERF-9..PERF-60 excluding 32 and 50 which are seed rules).
// Titles/descriptions align with Rust go_perf_metadata; fixes from metadata_overrides.
// Severity: Medium default; Info for B/C micro-opts (15–19, 35, 42, 46, 16).

var (
	MetaPERF9 = rules.Meta(
		"PERF-9",
		"URL Parsing Inside Loop",
		"Repeatedly parses URLs inside a tight loop rather than moving parsing or validation to a cheaper boundary.",
		rules.SeverityMedium,
		nil,
		"Parse each URL once outside the loop and store the *url.URL value.",
	)
	MetaPERF10 = rules.Meta(
		"PERF-10",
		"Template Parsing On Request Path",
		"Parses templates on the request path instead of compiling them once at process start or initialization time.",
		rules.SeverityMedium,
		nil,
		"Compile templates at process start with template.Must and reuse the parsed template.",
	)
	MetaPERF11 = rules.Meta(
		"PERF-11",
		"HTTP Client Allocation Per Request",
		"Allocates a new http.Client on request-serving paths where a shared client should usually be reused.",
		rules.SeverityMedium,
		nil,
		"Reuse a single http.Client (and Transport) declared at package scope.",
	)
	MetaPERF12 = rules.Meta(
		"PERF-12",
		"Prepared Statement Setup On Request Path",
		"Creates prepared statements on request-serving paths instead of caching or initializing them once.",
		rules.SeverityMedium,
		nil,
		"Prepare the statement once at startup and reuse it, or use a connection-pooled helper.",
	)
	MetaPERF13 = rules.Meta(
		"PERF-13",
		"time.After Allocation In Long-Running Loop",
		"Creates time.After values repeatedly inside long-running loops where reusable timers or tickers would avoid churn.",
		rules.SeverityMedium,
		nil,
		"Use a reusable *time.Timer with Stop+Reset, or a single time.Ticker, instead of time.After.",
	)
	MetaPERF14 = rules.Meta(
		"PERF-14",
		"Directory Scan Inside Loop",
		"Performs repeated filesystem directory scans inside loops when the scanned set is stable and could be hoisted or cached.",
		rules.SeverityMedium,
		nil,
		"Hoist filepath.Glob / os.ReadDir out of the loop and cache the directory listing.",
	)
	MetaPERF15 = rules.Meta(
		"PERF-15",
		"strconv Formatting Inside Loop",
		"Performs repeated numeric-to-string conversion inside loops where buffered writes, cached conversion, or batching would be cheaper.",
		rules.SeverityInfo,
		nil,
		"Use a strings.Builder or strconv.Append* to avoid repeated allocations.",
	)
	MetaPERF16 = rules.Meta(
		"PERF-16",
		"Buffer Reallocation Inside Loop",
		"Allocates a fresh bytes.Buffer inside an inner loop instead of reusing a resettable buffer across iterations.",
		rules.SeverityInfo,
		nil,
		"Reuse a single bytes.Buffer by calling Reset at the start of each iteration.",
	)
	MetaPERF17 = rules.Meta(
		"PERF-17",
		"String Concatenation With Plus Inside Loop",
		"Builds a string with the + operator inside a loop, causing repeated allocation and copy churn. Use strings.Builder or a pre-sized buffer.",
		rules.SeverityInfo,
		nil,
		"Hoist a strings.Builder outside the loop or use strings.Join to avoid repeated concatenation.",
	)
	MetaPERF18 = rules.Meta(
		"PERF-18",
		"Unnecessary Slice Copy On Reslice",
		"Copies a slice header into a function that only needs a sub-range, causing the full backing array to remain live and increase GC pressure.",
		rules.SeverityInfo,
		nil,
		"Pass a reslice of the original slice instead of copying when the callee does not mutate.",
	)
	MetaPERF19 = rules.Meta(
		"PERF-19",
		"Large Value Copy In Range Loop",
		"Iterates over a slice of large structs by value, copying each element. Range by index or pointer to avoid per-iteration copies.",
		rules.SeverityInfo,
		nil,
		"Range by index (&slice[i]) or pointer to avoid copying each struct value.",
	)
	MetaPERF20 = rules.Meta(
		"PERF-20",
		"Reflection Usage In Hot Path",
		"Uses the reflect package on a hot request path. Reflection is significantly slower than typed or codegen-based dispatch.",
		rules.SeverityMedium,
		nil,
		"Cache reflect.Type / reflect.Value at startup, or use code generation to avoid hot-path reflection.",
	)
	MetaPERF21 = rules.Meta(
		"PERF-21",
		"io.ReadAll On Large Request Body",
		"Uses io.ReadAll to fully buffer a request body that may be large. This materializes the whole body in memory and can amplify memory pressure.",
		rules.SeverityMedium,
		nil,
		"Stream the request body via json.NewDecoder or io.Copy instead of fully buffering with io.ReadAll.",
	)
	MetaPERF22 = rules.Meta(
		"PERF-22",
		"os.ReadFile Inside Handler",
		"Reads an entire file inside a request handler on every call. Large or remote files cause repeated disk I/O on a hot path.",
		rules.SeverityMedium,
		nil,
		"Load the file once at startup, or stream it, instead of reading on the request path.",
	)
	MetaPERF23 = rules.Meta(
		"PERF-23",
		"bytes.NewReader Allocation Per Request",
		"Allocates a new bytes.Reader or bytes.Buffer for each request when a sync.Pool or reusable buffer would be cheaper.",
		rules.SeverityMedium,
		nil,
		"Use a sync.Pool of *bytes.Reader or reuse a buffer across requests.",
	)
	MetaPERF24 = rules.Meta(
		"PERF-24",
		"Crypto Hashing In Tight Loop",
		"Creates a new hash instance (sha256, sha1, md5) inside a tight loop where the hasher could be reset or hoisted.",
		rules.SeverityMedium,
		nil,
		"Hoist the hasher out of the loop and call h.Reset() each iteration instead of allocating a new hasher.",
	)
	MetaPERF25 = rules.Meta(
		"PERF-25",
		"RSA Key Generation Per Request",
		"Generates RSA keys inside request handlers or per-call paths. RSA key generation is expensive and should be cached.",
		rules.SeverityMedium,
		nil,
		"Generate the key pair once at startup and reuse the private key across requests.",
	)
	MetaPERF26 = rules.Meta(
		"PERF-26",
		"Base64 Encode Or Decode In Loop",
		"Performs base64 encoding or decoding repeatedly inside a loop. Encoders and decoders can be reused, and streaming is often cheaper.",
		rules.SeverityMedium,
		nil,
		"Hoist base64.NewEncoder/Decoder outside the loop, or reuse a single encoder with a pooled bytes.Buffer.",
	)
	MetaPERF27 = rules.Meta(
		"PERF-27",
		"Missed sync.Pool Reuse Opportunity",
		"Allocates short-lived buffers or objects on hot paths without using sync.Pool, increasing GC pressure.",
		rules.SeverityMedium,
		nil,
		"Wrap the buffer in sync.Pool and call Get/Put around the hot section.",
	)
	MetaPERF28 = rules.Meta(
		"PERF-28",
		"Sync Mutex Allocation Per Request",
		"Embeds or constructs a sync.Mutex per request or per struct instance where a shared mutex or atomic primitive would be cheaper.",
		rules.SeverityMedium,
		nil,
		"Embed the mutex in a long-lived struct (package or request-scoped pool) instead of a per-request literal.",
	)
	MetaPERF29 = rules.Meta(
		"PERF-29",
		"Unbounded Goroutine Spawn",
		"Spawns goroutines without a worker pool, semaphore, or bounded concurrency, which can exhaust memory and file descriptors under load.",
		rules.SeverityMedium,
		nil,
		"Bound goroutines with a worker pool, semaphore channel, or errgroup with SetLimit.",
	)
	MetaPERF30 = rules.Meta(
		"PERF-30",
		"context.Background In Request Goroutine",
		"Creates a new context.Background() inside a goroutine launched from a request handler, breaking cancellation and request-scoped values.",
		rules.SeverityMedium,
		nil,
		"Propagate c.Request.Context() or the caller's context to the goroutine instead of context.Background().",
	)
	MetaPERF31 = rules.Meta(
		"PERF-31",
		"Defer In Hot Function",
		"Uses defer in performance-sensitive hot functions where the per-call overhead is measurable. defer has a small but non-zero cost in Go 1.13 and earlier.",
		rules.SeverityMedium,
		nil,
		"Move the cleanup into a helper function that returns a Close() method called from a single defer outside the loop.",
	)
	MetaPERF33 = rules.Meta(
		"PERF-33",
		"Range Over Large Slice With Early Break Missing",
		"Iterates over a large slice without using indexed access, missing opportunities for early termination or two-pointer patterns.",
		rules.SeverityMedium,
		nil,
		"Use an indexed scan with an explicit break when you only need the first match, or stream the slice.",
	)
	MetaPERF34 = rules.Meta(
		"PERF-34",
		"Map Iteration With Append",
		"Appends to a slice inside a map iteration, which is safe but often indicates accidental quadratic copying when the slice grows large.",
		rules.SeverityMedium,
		nil,
		"Preallocate the destination slice with make([]T, 0, len(m)) before the range loop.",
	)
	MetaPERF35 = rules.Meta(
		"PERF-35",
		"Interface Boxing On Hot Path",
		"Passes concrete values into interface{} parameters on hot paths, causing escape-to-heap and dynamic dispatch overhead.",
		rules.SeverityInfo,
		nil,
		"Cast non-string args to a concrete type or use strconv/strings builders to avoid interface boxing.",
	)
	MetaPERF36 = rules.Meta(
		"PERF-36",
		"Loop Variable Capture In Goroutine",
		"Captures a loop variable by reference inside a goroutine (pre-Go 1.22 semantics), causing all goroutines to see the final value.",
		rules.SeverityMedium,
		nil,
		"Copy the loop variable into a per-iteration local (v := v) before launching the goroutine.",
	)
	MetaPERF37 = rules.Meta(
		"PERF-37",
		"Slice Growth Miscalculation",
		"Repeatedly appends to a slice without preallocation or with insufficient capacity hint, causing repeated reallocation and copy.",
		rules.SeverityMedium,
		nil,
		"Replace the var declaration with make([]T, 0, hint) to give the runtime a growth target.",
	)
	MetaPERF38 = rules.Meta(
		"PERF-38",
		"Unbuffered Channel In Producer Consumer",
		"Uses an unbuffered channel for high-throughput producer/consumer patterns, causing unnecessary blocking and synchronization.",
		rules.SeverityMedium,
		nil,
		"Use make(chan T, N) with a buffer sized to expected concurrency.",
	)
	MetaPERF39 = rules.Meta(
		"PERF-39",
		"Select With Default In Tight Loop",
		"Uses select with a default case in a tight loop, which becomes a busy-wait pattern that pins a CPU core.",
		rules.SeverityMedium,
		nil,
		"Use a time.Sleep backoff or remove the default branch to avoid busy-looping.",
	)
	MetaPERF40 = rules.Meta(
		"PERF-40",
		"Repeated time.Now On Hot Path",
		"Calls time.Now() repeatedly inside hot paths where a single timestamp at function entry would suffice.",
		rules.SeverityMedium,
		nil,
		"Hoist time.Now() outside the function, or cache the value in a struct field, when measuring a single event.",
	)
	MetaPERF41 = rules.Meta(
		"PERF-41",
		"Stdlib log On Hot Path",
		"Uses the standard log package or log.Println/Printf on production hot paths, where structured, leveled, or async loggers are cheaper.",
		rules.SeverityMedium,
		nil,
		"Route logs through a structured logger (slog/zap/zerolog) and gate debug levels with build tags.",
	)
	MetaPERF42 = rules.Meta(
		"PERF-42",
		"fmt.Errorf Without Format Verbs",
		"Uses fmt.Errorf for error wrapping without format verbs, which is slower than errors.New and requires extra parsing.",
		rules.SeverityInfo,
		nil,
		"Use errors.New or a sentinel error when the message has no format verbs.",
	)
	MetaPERF43 = rules.Meta(
		"PERF-43",
		"Panic Recovery In Hot Path",
		"Recovers from panics in tight loops or hot functions, which adds defer setup and runtime cost on every iteration.",
		rules.SeverityMedium,
		nil,
		"Move recover() to a middleware boundary instead of a per-request defer.",
	)
	MetaPERF44 = rules.Meta(
		"PERF-44",
		"Repeated Type Assertion On Same Interface",
		"Performs the same type assertion on an interface value multiple times, where a single assertion with reuse would be cheaper.",
		rules.SeverityMedium,
		nil,
		"Bind the asserted value to a local once (v, ok := x.(T)) and reuse the binding.",
	)
	MetaPERF45 = rules.Meta(
		"PERF-45",
		"Append Without Capacity Hint In Loop",
		"Appends to a slice inside a loop without setting an initial capacity, leading to repeated reallocation as the slice grows geometrically.",
		rules.SeverityMedium,
		nil,
		"Preallocate with make([]T, 0, hint) before the loop so append does not reallocate.",
	)
	MetaPERF46 = rules.Meta(
		"PERF-46",
		"String Trimming With Allocations",
		"Uses strings.TrimSpace or similar trimming functions on hot paths where the result is the same as the input, allocating a new string.",
		rules.SeverityInfo,
		nil,
		"Advisory micro-opt: guard before TrimSpace/Trim when a cheap length or edge-byte check avoids the alloc; intentional header/value trimming is often fine.",
	)
	MetaPERF47 = rules.Meta(
		"PERF-47",
		"strings.Split Allocation In Loop",
		"Uses strings.Split or strings.SplitN in hot paths where a streaming scanner or precomputed index would be cheaper.",
		rules.SeverityMedium,
		nil,
		"Use strings.SplitSeq or a manual index loop over strings.IndexByte to avoid the []string allocation.",
	)
	MetaPERF48 = rules.Meta(
		"PERF-48",
		"Equality Without Length Precheck",
		"Compares long byte slices or strings without short-circuit optimizations, where length or hash precheck would help.",
		rules.SeverityMedium,
		nil,
		"Add an early length-mismatch or prefix check before the comparison.",
	)
	MetaPERF49 = rules.Meta(
		"PERF-49",
		"Copy With Mismatched Length",
		"Uses copy(dst, src) where the destination length is not checked, leading to silent truncation or wasted capacity.",
		rules.SeverityMedium,
		nil,
		"Validate the payload length and size the destination buffer to the exact count before copy().",
	)
	MetaPERF51 = rules.Meta(
		"PERF-51",
		"unsafe.Pointer In Request Handler",
		"Uses unsafe.Pointer for type punning or zero-copy conversion in ways that hurt readability and may cause subtle performance issues when miscompiled by future Go versions.",
		rules.SeverityMedium,
		nil,
		"Replace the unsafe conversion with strconv.Quote/Unquote or a measured []byte(s) outside the loop.",
	)
	MetaPERF52 = rules.Meta(
		"PERF-52",
		"Manual runtime.GC Call",
		"Calls runtime.GC() in production code, forcing a stop-the-world GC pause and undermining the runtime's tuning.",
		rules.SeverityMedium,
		nil,
		"Remove runtime.GC; the runtime already manages collection, and manual calls slow the allocator.",
	)
	MetaPERF53 = rules.Meta(
		"PERF-53",
		"Package-Level math/rand Contention",
		"Uses the package-level math/rand functions, which are guarded by a global mutex and can become a contention hotspot.",
		rules.SeverityMedium,
		nil,
		"Use a per-goroutine rand.NewSource(rand.New(rand.NewSource(time.Now().UnixNano()))) or math/rand/v2.",
	)
	MetaPERF54 = rules.Meta(
		"PERF-54",
		"strings.Builder Allocated Per Request",
		"Allocates a new strings.Builder per request instead of calling Reset on a pooled or reused builder.",
		rules.SeverityMedium,
		nil,
		"Hoist a *strings.Builder at package scope and call b.Reset() before each reuse.",
	)
	MetaPERF55 = rules.Meta(
		"PERF-55",
		"bufio.Scanner Without Buffer Sizing",
		"Uses bufio.Scanner with the default 64KiB max token size on potentially large input, which silently truncates or errors out.",
		rules.SeverityMedium,
		nil,
		"Call scanner.Buffer(make([]byte, initial), max) before Scan to size the token buffer.",
	)
	MetaPERF56 = rules.Meta(
		"PERF-56",
		"Gin c.JSON With Marshal In Loop",
		"Calls c.JSON inside a loop with on-the-fly marshaling per iteration. Marshaling per iteration defeats batch response strategies and can balloon CPU under load.",
		rules.SeverityMedium,
		nil,
		"Collect the response items into a slice and call c.JSON once, or stream with c.Stream.",
	)
	MetaPERF57 = rules.Meta(
		"PERF-57",
		"Gin Middleware Heavy Allocation",
		"Defines middleware that allocates large buffers, parses bodies, or builds heavy structures for every request even when downstream handlers do not need them.",
		rules.SeverityMedium,
		nil,
		"Move heavy parsing (io.ReadAll/json.Unmarshal) out of middleware into the handler, or cache the result.",
	)
	MetaPERF58 = rules.Meta(
		"PERF-58",
		"Gin c.Request.Body Not Closed",
		"Reads c.Request.Body without deferring Close, causing connection leaks and request body buffer retention.",
		rules.SeverityMedium,
		nil,
		"Add defer c.Request.Body.Close() after every access, or drain with io.Copy(io.Discard, body).",
	)
	MetaPERF59 = rules.Meta(
		"PERF-59",
		"Gin Binding Reuse Missed",
		"Uses c.ShouldBindJSON, c.ShouldBind, or c.Bind on every request even when the same DTO is reused, missing typed handler or struct reuse opportunities.",
		rules.SeverityMedium,
		nil,
		"Bind once into a pre-validated struct or share a custom binding.Validator across handlers.",
	)
	MetaPERF60 = rules.Meta(
		"PERF-60",
		"Gin Render Allocation Per Request",
		"Allocates a new render.Render or per-request render struct instead of relying on the default render pool Gin's c.Render manages.",
		rules.SeverityMedium,
		nil,
		"Reuse a single render.JSON / render.HTML instance created at startup, or use c.Render / c.HTML.",
	)
)
