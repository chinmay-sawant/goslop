package perf

import "github.com/chinmay/goslop/internal/rules"

// Batch 5 metadata: PERF-215..242 (excluding PERF-230, which lives in seed).
var (
	MetaPERF215 = rules.Meta(
		"PERF-215",
		"bytes.Buffer or strings.Builder Without Pre-Sizing",
		"A bytes.Buffer or strings.Builder is created or reused without calling Grow() when the final content size is knowable or estimable upfront. Without pre-sizing, the underlying buffer grows geometrically via reallocation+copy (bytes.growSlice), causing excessive memory churn and CPU time in memmove. Pre-size with Grow(expectedSize) when the output size is known or can be bounded.",
		rules.SeverityMedium, // unclassified PERF → Medium
		nil,
		"Pre-size the bytes.Buffer or strings.Builder by calling Grow(expectedSize) when the final content size is known or can be estimated from input parameters.",
	)
	MetaPERF216 = rules.Meta(
		"PERF-216",
		"Hot-Path Struct Allocation Without Slab Arena",
		"Hot-path structures (e.g., tree nodes, table rows, linked-list elements) are individually heap-allocated inside tight loops or high-frequency code paths. Each allocation adds GC scan pressure and allocator overhead. For structures with bounded lifetime and known type, use a slab or arena allocator (pre-allocated contiguous blocks) to batch allocations and reduce GC work.",
		rules.SeverityMedium,
		nil,
		"Replace individual per-element heap allocations in the hot path with a slab or arena allocator that pre-allocates contiguous blocks for the struct type.",
	)
	MetaPERF217 = rules.Meta(
		"PERF-217",
		"Static Computation Rebuilt Per Operation",
		"A deterministic computation (ICC profile generation, metadata serialization, color space construction, font object generation) that always produces the same output is executed on every request or operation instead of once at startup. Hoist the computation to package init time or cache the result for the process lifetime. Each rebuild wastes CPU and allocation for identical results.",
		rules.SeverityMedium,
		nil,
		"Cache the deterministic computation result in a package-level variable computed at init(); the output never changes per request.",
	)
	MetaPERF218 = rules.Meta(
		"PERF-218",
		"sync.Pool or Cache Without Per-CPU Sharding",
		"A single sync.Pool or shared cache structure is contended across many goroutines (e.g., 48+ concurrent workers on 24 cores). Under high concurrency, a single pool becomes a contention bottleneck as goroutines compete for the same underlying data structure. Use per-CPU or per-P sharding, or partition the pool by capacity class, to reduce lock contention.",
		rules.SeverityMedium,
		nil,
		"Replace the single contended sync.Pool with per-CPU shards (e.g., a [runtime.NumCPU()]sync.Pool array) to reduce lock contention under high concurrency.",
	)
	MetaPERF219 = rules.Meta(
		"PERF-219",
		"Oversized Object Returned to sync.Pool",
		"A buffer or large object is returned to sync.Pool after use without checking whether its capacity exceeds a sensible threshold. When oversized objects (e.g., 8 MB buffers) are returned to the pool, subsequent acquire calls from other goroutines receive these oversized objects, wasting memory across workers. Discard objects that exceed a capacity threshold instead of returning them to the pool.",
		rules.SeverityMedium,
		nil,
		"Guard the Put call: if cap(obj) > maxSize { return } to discard oversized buffers instead of returning them to the pool.",
	)
	MetaPERF220 = rules.Meta(
		"PERF-220",
		"Sequential Scans Over Identical Data",
		"The same data structure (slice, map, tree) is iterated multiple times sequentially in the same function or call chain, performing the same or related work on each pass. Repeated scanning multiplies CPU time proportional to the data size times the number of passes. Merge the passes into a single iteration that does all required work in one go.",
		rules.SeverityLow,
		nil,
		"Merge the consecutive loops over the same data into a single pass that does all required work, eliminating the redundant iteration overhead.",
	)
	MetaPERF221 = rules.Meta(
		"PERF-221",
		"map[int]T for Dense Sequential Integer Keys",
		"A map[int]T or map[int64]T is used to store values keyed by dense sequential integers (e.g. array indices or counters). For sequential integer keys, a []T slice provides O(1) access with no hashing overhead, better cache locality, and no per-entry heap allocation for key/value pairs. Replace map[int]T with []T when keys are sequential or can be made sequential.",
		rules.SeverityMedium, // unclassified PERF → Medium
		nil,
		"Replace map[int]T with []T when the integer keys are dense and sequential (e.g. indices, counters). Use make([]T, maxKey+1) and direct index access.",
	)
	MetaPERF222 = rules.Meta(
		"PERF-222",
		"Generic Function on Measured Hot Path",
		"A generic function (func[T any] or func[T constraint]) is called on a measured hot path. Go generics use shape-based dispatch: each unique type instantiation generates a separate function, but the call still goes through a dictionary, defeating inlining and adding call overhead. For measured hot paths, prefer concrete types or code generation over generics.",
		rules.SeverityLow,
		nil,
		"Replace the generic function on the measured hot path with a concrete type or use code generation. Shape-based dispatch prevents inlining and adds call overhead.",
	)
	MetaPERF223 = rules.Meta(
		"PERF-223",
		"sync.Pool Backing Array Discarded on Return",
		"A slice or struct containing a backing array is returned to sync.Pool after being set to nil or zeroed, discarding the pre-allocated backing array capacity. On subsequent acquire, the pool returns an object with no backing capacity, forcing a fresh allocation. When objects returned to a pool are expected to be reused, retain the backing array by using Reset() (s = s[:0]) instead of Release() (s = nil).",
		rules.SeverityLow,
		nil,
		"Retain backing array capacity on pool return: use obj.Reset() (obj.Slice = obj.Slice[:0]) instead of obj.Slice = nil so the backing array is reused on next acquire.",
	)
	MetaPERF224 = rules.Meta(
		"PERF-224",
		"Recursive Tree Walk on Hot Execution Path",
		"A deeply nested tree or graph is traversed using recursive function calls (walk, visit, assign) on a hot execution path. Recursive traversal adds per-call overhead for function frames, parameter passing, and return handling. For trees where the structure is already available in a flat pre-ordered slice, replace the recursive walk with an iterative loop over the flat representation.",
		rules.SeverityLow,
		nil,
		"Replace the recursive tree walk with an iterative loop over the existing flat pre-ordered representation of the same data.",
	)
	MetaPERF225 = rules.Meta(
		"PERF-225",
		"Redundant Large Slice Clone",
		"A large slice is fully cloned more than once in the same function (slices.Clone, append([]T(nil), src...), or equivalent), discarding ownership of an intermediate buffer. Prefer a single clone or in-place mutation when the source is already exclusively owned.",
		rules.SeverityMedium,
		nil,
		"Keep a single owned buffer: clone once, or mutate in place when the source is exclusive. Avoid chaining slices.Clone / append([]T(nil), â¦) on the same data.",
	)
	MetaPERF226 = rules.Meta(
		"PERF-226",
		"Post-Producer Buffer Re-Copy",
		"Immediately after a producer fills or returns a buffer (bytes.Buffer.Bytes, compress writer Close, similar), the code allocates a second slice and copy/clones the result without the source being returned to a pool. Prefer taking ownership of the producer buffer or copying only when the source must be recycled.",
		rules.SeverityMedium,
		nil,
		"Return or use the producer buffer directly after Bytes()/Close(). Only copy when the source must go back to a pool, and copy before Put â not after an exclusive local.",
	)
	MetaPERF227 = rules.Meta(
		"PERF-227",
		"Compress Writer Allocated Without Pool",
		"A flate, zlib, or gzip writer is constructed with NewWriter / NewWriterLevel on a hot path or inside a loop without reusing a pooled writer via Reset. Writer construction allocates deflate state; pool + Reset is the standard reuse pattern.",
		rules.SeverityMedium,
		nil,
		"Pool flate/zlib/gzip writers and call Reset(dst) on each use instead of NewWriter on every encode.",
	)
	MetaPERF228 = rules.Meta(
		"PERF-228",
		"Parallel Fan-Out For Tiny Workset",
		"A tiny workset (composite slice literal with 1–2 elements, or a for-range over such a slice) is parallelized with errgroup.Go, WaitGroup + go, or bare go func per item. Spawn and synchronization cost often exceeds the work for N≤2; prefer a serial path for small fixed fan-out.",
		rules.SeverityMedium,
		nil,
		"For worksets of 1â2 items, run the work serially instead of errgroup/WaitGroup/go fan-out; spawn cost usually dominates.",
	)
	MetaPERF229 = rules.Meta(
		"PERF-229",
		"Intermediate String On Byte Append Path",
		"A temporary string is built with strconv.Itoa, FormatInt, or fmt.Sprintf and then immediately appended or written into a []byte / bytes.Buffer sink. Prefer strconv.AppendInt / AppendFormat-style APIs that write into the destination buffer without an intermediate string.",
		rules.SeverityMedium,
		nil,
		"Write numbers/text into the destination with strconv.AppendInt / AppendUint / AppendFloat (or Builder) instead of Itoa/Sprintf then append([]byte(s)).",
	)
	MetaPERF231 = rules.Meta(
		"PERF-231",
		"PEM Or Key Material Parsed On Hot Path",
		"PEM blocks, X.509 certificates, or private keys are parsed on a hot path (encode/sign/serve/loop) instead of once at process start. Parsing is expensive relative to using a cached *rsa.PrivateKey / certificate chain. Distinct from key generation (PERF-025).",
		rules.SeverityMedium,
		nil,
		"Parse PEM/keys once at process start (package var or sync.Once) and reuse *rsa.PrivateKey / certificates on the hot path.",
	)
	MetaPERF232 = rules.Meta(
		"PERF-232",
		"Unbounded Parallel Fan-Out",
		"Parallel work is fanned out from a loop with errgroup without a concurrency bound such as SetLimit or a semaphore. Unbounded fan-out can overwhelm CPU, memory, and pooled resources.",
		rules.SeverityMedium,
		nil,
		"Cap fan-out with errgroup.SetLimit or a semaphore before spawning per-item work.",
	)
	MetaPERF233 = rules.Meta(
		"PERF-233",
		"Slow Compress Level On Hot Encode Path",
		"A flate/zlib/gzip writer is created with DefaultCompression, BestCompression, or the default NewWriter level on a hot encode path. For bulk stream compression where output size budgets allow, BestSpeed (or level 1) typically cuts CPU substantially while remaining compatible with standard flate decoders.",
		rules.SeverityMedium,
		nil,
		"Use flate/zlib BestSpeed (or level 1) for hot stream compression when size budgets allow. Reserve Default/BestCompression for cold or archival paths.",
	)
	MetaPERF234 = rules.Meta(
		"PERF-234",
		"Bulk Buffer Without Workload Sizing",
		"A bulk output buffer uses a large fixed Grow capacity or is taken from a pool, reset, and written without an estimated Grow. Prefer sizing from a known workload estimate when possible.",
		rules.SeverityMedium,
		nil,
		"Grow bulk buffers from a workload estimate (payload length, item count), not a large fixed default or bare pool Reset.",
	)
	MetaPERF235 = rules.Meta(
		"PERF-235",
		"Intermediate Strings Builder Bridge",
		"A temporary strings.Builder is filled and then converted with .String() into a WriteString, Write, or append sink. Prefer writing operators/bytes into the destination []byte or *bytes.Buffer directly.",
		rules.SeverityMedium,
		nil,
		"Write into the destination *bytes.Buffer or []byte directly instead of building a temporary strings.Builder and flushing with .String().",
	)
	MetaPERF236 = rules.Meta(
		"PERF-236",
		"Full Buffer Clone On Signing Path",
		"A signing helper clones an entire buffer before patching fixed-width fields. Prefer an owned writable buffer or in-place mutation when the caller already owns the data.",
		rules.SeverityMedium,
		nil,
		"Keep an owned writable buffer with reserved holes, or mutate in place, instead of bytes.Clone of the entire document on the signing path.",
	)
	MetaPERF237 = rules.Meta(
		"PERF-237",
		"Always-Parallel Fan-Out Without Tiny-N Serial Path",
		"errgroup fans out over a dynamic range with no serial short-circuit for tiny worksets (len <= 2). Spawn cost often dominates for one or two items; run those serially first.",
		rules.SeverityMedium,
		nil,
		"Before errgroup fan-out, handle len(items) <= 2 on a serial path; spawn only when concurrency pays off.",
	)
	MetaPERF238 = rules.Meta(
		"PERF-238",
		"Rune Membership Map In Loop",
		"Rune membership is tracked with map[rune]bool and updated inside a loop. For BMP-heavy text, a compact bitset or denser set reduces hashing and allocation overhead.",
		rules.SeverityMedium, // unclassified PERF → Medium
		nil,
		"Replace map[rune]bool membership on hot paths with a bitset or denser set when the code-point domain is bounded.",
	)
	MetaPERF239 = rules.Meta(
		"PERF-239",
		"Dense Integer Map Write Churn",
		"An integer-keyed map is written many times (>=6) in one function after make(map[int]…). Dense integer keys are usually better served by a slice or append-only offset records with one final index pass.",
		rules.SeverityMedium, // unclassified PERF → Medium
		nil,
		"Prefer a slice (or append-only {id,offset} records + one index pass) instead of many map[int] writes for dense keys.",
	)
	MetaPERF240 = rules.Meta(
		"PERF-240",
		"Unpooled Len-Sized Byte Scratch",
		"A hot path allocates make([]byte, len(src)) (or equivalent) without reusing a pooled scratch buffer. Prefer sync.Pool + Reset for large temporary copies of source tables or payloads.",
		rules.SeverityMedium,
		nil,
		"Pool a []byte scratch (sync.Pool + Reset/[:0]) instead of make([]byte, len(src)) on every hot call.",
	)
	MetaPERF241 = rules.Meta(
		"PERF-241",
		"ASN1 Remarshal With Fresh Time On Sign Path",
		"A signing/CMS-style helper calls asn1.Marshal while also using time.Now, typically re-serializing mostly-static structure every call. Cache immutable DER components and only re-marshal time-varying attributes.",
		rules.SeverityMedium,
		nil,
		"Pre-marshal immutable ASN.1/CMS components; only re-marshal time-varying authenticated attributes each sign.",
	)
	MetaPERF242 = rules.Meta(
		"PERF-242",
		"Per-Iteration Encode Scratch Allocation",
		"A loop allocates make([]byte, … len(x)*N …) every iteration for encoding or buffering. Reuse one scratch buffer outside the loop with [:0] growth.",
		rules.SeverityMedium,
		nil,
		"Hoist make([]byte, len(x)*N) out of the loop and reuse one buffer with append/[:0].",
	)
)
