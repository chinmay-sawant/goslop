// Code generated from ruleset/golang/bad-practices.json. DO NOT EDIT BY HAND.
package badpractices

import "github.com/chinmay-sawant/goslop/internal/rules"

// metaByID is the BP catalogue loaded from bad-practices.json at init.
var metaByID = map[string]*rules.RuleMetadata{}

func init() {
	loadEmbeddedMetadata()
}

func loadEmbeddedMetadata() {
	metaByID["BP-1"] = &rules.RuleMetadata{
		ID:          "BP-1",
		Title:       "Discarded Error Return",
		Description: "A returned error is assigned to `_`, suppressing error handling.",
		Severity:    rules.SeverityLow,
		Fix:         "Match assignment_statement or short_var_declaration with `_ =` on an error-producing call.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-2"] = &rules.RuleMetadata{
		ID:          "BP-2",
		Title:       "Naked Error Return",
		Description: "An error is returned without wrapping contextual information.",
		Severity:    rules.SeverityLow,
		Fix:         "Match `return err` sites that lose operation-specific context.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-3"] = &rules.RuleMetadata{
		ID:          "BP-3",
		Title:       "Panic Outside Main Or Test",
		Description: "panic is called outside main() or test files.",
		Severity:    rules.SeverityLow,
		Fix:         "Match panic() outside main() and `_test.go` files.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-4"] = &rules.RuleMetadata{
		ID:          "BP-4",
		Title:       "Recover Without Logging",
		Description: "recover() is used without recording the recovered panic.",
		Severity:    rules.SeverityLow,
		Fix:         "Match recover() handlers that do not log, report, or rethrow the panic.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-5"] = &rules.RuleMetadata{
		ID:          "BP-5",
		Title:       "Ignored Close Error",
		Description: "Close() is called without checking its returned error.",
		Severity:    rules.SeverityLow,
		Fix:         "Match `.Close()` calls whose error result is ignored.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-6"] = &rules.RuleMetadata{
		ID:          "BP-6",
		Title:       "WaitGroup Add Inside Goroutine",
		Description: "sync.WaitGroup.Add is called inside the goroutine it tracks.",
		Severity:    rules.SeverityMedium,
		Fix:         "Review-only AST match for `go func` bodies that call `.Add(`; nested goroutines are scoped independently, but receiver type is not proven and atomic counters can match.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-7"] = &rules.RuleMetadata{
		ID:          "BP-7",
		Title:       "Mutex Passed By Value",
		Description: "sync.Mutex is passed by value, copying lock state.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match function parameters that take `sync.Mutex` by value.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-8"] = &rules.RuleMetadata{
		ID:          "BP-8",
		Title:       "Unlock Deferred On Mutex Copy",
		Description: "defer mu.Unlock() is used on a mutex value copy.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match deferred unlocks on mutex values that have been copied.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-9"] = &rules.RuleMetadata{
		ID:          "BP-9",
		Title:       "Blocking Select Without Timeout",
		Description: "select waits without a default branch, timeout, or context cancellation.",
		Severity:    rules.SeverityLow,
		Fix:         "Match `select {}` blocks with no escape hatch.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-10"] = &rules.RuleMetadata{
		ID:          "BP-10",
		Title:       "time.After Inside Loop",
		Description: "time.After is called inside a loop, allocating a timer per iteration.",
		Severity:    rules.SeverityLow,
		Fix:         "Match `time.After()` in `for` or `range` bodies.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-11"] = &rules.RuleMetadata{
		ID:          "BP-11",
		Title:       "Defer Inside Loop",
		Description: "defer is used inside a loop body.",
		Severity:    rules.SeverityLow,
		Fix:         "Match defer statements nested in loop bodies.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-12"] = &rules.RuleMetadata{
		ID:          "BP-12",
		Title:       "Unbuffered Channel Send From Multiple Goroutines",
		Description: "Multiple goroutines send on an unbuffered channel without coordinated receivers.",
		Severity:    rules.SeverityLow,
		Fix:         "Review-only same-file heuristic; it cannot prove channel ownership, all receiver paths, or helper-mediated coordination.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-13"] = &rules.RuleMetadata{
		ID:          "BP-13",
		Title:       "Background Context In Library Function",
		Description: "context.Background() is used outside main or initialization code.",
		Severity:    rules.SeverityLow,
		Fix:         "Match `context.Background()` in library functions and helper code.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-14"] = &rules.RuleMetadata{
		ID:          "BP-14",
		Title:       "Goroutine Without Context Cancellation",
		Description: "A goroutine is launched without listening for ctx.Done() or another shutdown signal.",
		Severity:    rules.SeverityLow,
		Fix:         "Review-only same-file heuristic; it cannot prove ownership, alternative shutdown channels, or helper-mediated lifecycle control.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-15"] = &rules.RuleMetadata{
		ID:          "BP-15",
		Title:       "Recursive sync.Once.Do",
		Description: "sync.Once.Do invokes a closure that recursively calls the same Once.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match direct or same-file helper call chains that recurse into the same `sync.Once.Do`.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-16"] = &rules.RuleMetadata{
		ID:          "BP-16",
		Title:       "time.Sleep In Test",
		Description: "A test uses time.Sleep instead of deterministic synchronization or polling.",
		Severity:    rules.SeverityLow,
		Fix:         "Match `time.Sleep()` in `_test.go` files outside controlled retry loops.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-17"] = &rules.RuleMetadata{
		ID:          "BP-17",
		Title:       "Redundant t.Error Before t.Fatal",
		Description: "A test calls t.Error and then immediately escalates to t.Fatal.",
		Severity:    rules.SeverityLow,
		Fix:         "Match consecutive `t.Error` or `t.Errorf` followed by `t.Fatal` or `t.Fatalf`.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-18"] = &rules.RuleMetadata{
		ID:          "BP-18",
		Title:       "t.Error Without Early Exit",
		Description: "A test records an error and continues execution without returning or failing now.",
		Severity:    rules.SeverityLow,
		Fix:         "Match `t.Error` or `t.Errorf` paths that keep executing the same block.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-19"] = &rules.RuleMetadata{
		ID:          "BP-19",
		Title:       "Missing t.Helper In Test Helper",
		Description: "A helper used by tests omits t.Helper(), reducing useful failure locations.",
		Severity:    rules.SeverityLow,
		Fix:         "Match helper functions called from tests whose first statement is not `t.Helper()`.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-20"] = &rules.RuleMetadata{
		ID:          "BP-20",
		Title:       "Table Test Without t.Run",
		Description: "A table-driven test loops over cases without subtests.",
		Severity:    rules.SeverityLow,
		Fix:         "Match table-driven loops in `_test.go` files that do not call `t.Run`.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-21"] = &rules.RuleMetadata{
		ID:          "BP-21",
		Title:       "Subtest Missing t.Parallel",
		Description: "A table-driven subtest omits t.Parallel where the pattern expects isolated parallel cases.",
		Severity:    rules.SeverityInfo,
		Fix:         "Match `t.Run` bodies in table tests that omit `t.Parallel()`.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-22"] = &rules.RuleMetadata{
		ID:          "BP-22",
		Title:       "TestMain Without os.Exit",
		Description: "TestMain does not exit with the result of m.Run().",
		Severity:    rules.SeverityLow,
		Fix:         "Match `func TestMain(m *testing.M)` bodies missing `os.Exit(m.Run())`.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-23"] = &rules.RuleMetadata{
		ID:          "BP-23",
		Title:       "Long Test Without testing.Short Guard",
		Description: "A long-running test does not respect testing.Short().",
		Severity:    rules.SeverityLow,
		Fix:         "Heuristic for large `_test.go` bodies with no `testing.Short()` guard.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-24"] = &rules.RuleMetadata{
		ID:          "BP-24",
		Title:       "Test File Without Test Functions",
		Description: "A `_test.go` file defines no `Test*` functions.",
		Severity:    rules.SeverityLow,
		Fix:         "Match `_test.go` files with zero `func Test*` declarations.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-25"] = &rules.RuleMetadata{
		ID:          "BP-25",
		Title:       "Test Helper Returns Error Instead Of Failing Test",
		Description: "A test helper returns an error that the caller always converts into t.Fatal-style handling.",
		Severity:    rules.SeverityLow,
		Fix:         "Match helper functions in tests that return `error` instead of owning the test failure path.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-26"] = &rules.RuleMetadata{
		ID:          "BP-26",
		Title:       "Context Not First Parameter",
		Description: "A public Go API accepts context.Context but not as its first parameter.",
		Severity:    rules.SeverityLow,
		Fix:         "Match parameter lists where `context.Context` is present but not first.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-27"] = &rules.RuleMetadata{
		ID:          "BP-27",
		Title:       "Exported Function Returns Unexported Type",
		Description: "An exported function returns a package-private concrete type.",
		Severity:    rules.SeverityLow,
		Fix:         "Match exported functions whose return types start with a lowercase identifier.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-28"] = &rules.RuleMetadata{
		ID:          "BP-28",
		Title:       "Single Method Interface",
		Description: "An interface with one method could often be represented as a function type.",
		Severity:    rules.SeverityInfo,
		Fix:         "Match interface declarations with exactly one method.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-29"] = &rules.RuleMetadata{
		ID:          "BP-29",
		Title:       "Interface Bloat",
		Description: "An interface declares too many methods and becomes difficult to satisfy or mock.",
		Severity:    rules.SeverityInfo,
		Fix:         "Match interfaces with more than five methods.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-30"] = &rules.RuleMetadata{
		ID:          "BP-30",
		Title:       "Exported Interface Without Documented Implementation",
		Description: "An exported interface is exposed without clear evidence of a same-package implementation; this is opt-in style advice because external implementations are common.",
		Severity:    rules.SeverityInfo,
		Fix:         "Match exported interfaces whose implementing types are not evident in the package. Review external implementations and intentional capability boundaries manually.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-31"] = &rules.RuleMetadata{
		ID:          "BP-31",
		Title:       "Function Returns Concrete Type Instead Of Interface",
		Description: "A function exposes a concrete implementation where an interface boundary would be more stable.",
		Severity:    rules.SeverityInfo,
		Fix:         "Heuristic for exported constructors or factories returning concrete implementation types.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-32"] = &rules.RuleMetadata{
		ID:          "BP-32",
		Title:       "String Alias Error Type",
		Description: "An error type is modeled as a string alias instead of a structured type.",
		Severity:    rules.SeverityLow,
		Fix:         "Match `type X string` definitions that implement `Error() string`.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-33"] = &rules.RuleMetadata{
		ID:          "BP-33",
		Title:       "Sentinel Error Without Is Method",
		Description: "A sentinel-like error type lacks an Is method for stable comparisons.",
		Severity:    rules.SeverityLow,
		Fix:         "Match custom error types intended as sentinels without `Is(error) bool`.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-34"] = &rules.RuleMetadata{
		ID:          "BP-34",
		Title:       "Error Wrapping Without Percent-W",
		Description: "fmt.Errorf wraps an error using `%v` or `%s` instead of `%w`.",
		Severity:    rules.SeverityLow,
		Fix:         "Match `fmt.Errorf` calls that mention an error argument without `%w`.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-36"] = &rules.RuleMetadata{
		ID:          "BP-36",
		Title:       "init With Side Effects",
		Description: "An init() function performs work beyond simple variable setup.",
		Severity:    rules.SeverityLow,
		Fix:         "Match `init()` bodies containing non-trivial side-effect statements.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-37"] = &rules.RuleMetadata{
		ID:          "BP-37",
		Title:       "Package-Level Mutable Global",
		Description: "A package exposes mutable global state outside tests.",
		Severity:    rules.SeverityLow,
		Fix:         "Match package-level mutable variables in non-test files.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-38"] = &rules.RuleMetadata{
		ID:          "BP-38",
		Title:       "Unused Unexported Helper",
		Description: "An unexported helper has no same-package callers.",
		Severity:    rules.SeverityLow,
		Fix:         "Match unexported helper functions with zero internal call sites.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-39"] = &rules.RuleMetadata{
		ID:          "BP-39",
		Title:       "Exported Function Without Doc Comment",
		Description: "An exported API surface is missing its doc comment.",
		Severity:    rules.SeverityInfo,
		Fix:         "Match exported functions without preceding package comments.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-40"] = &rules.RuleMetadata{
		ID:          "BP-40",
		Title:       "Unrelated Constants In One Block",
		Description: "A package-level const block groups unrelated constants together.",
		Severity:    rules.SeverityInfo,
		Fix:         "Heuristic for const blocks whose names do not share a prefix or domain.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-41"] = &rules.RuleMetadata{
		ID:          "BP-41",
		Title:       "Missing Package Doc Comment",
		Description: "A package lacks a package-level documentation comment.",
		Severity:    rules.SeverityInfo,
		Fix:         "Match Go files in a package set with no package doc comment.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-42"] = &rules.RuleMetadata{
		ID:          "BP-42",
		Title:       "One-Off Import Alias",
		Description: "An import alias is introduced without consistent package-wide benefit.",
		Severity:    rules.SeverityInfo,
		Fix:         "Match aliased imports that appear only once or without repeated local use.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-43"] = &rules.RuleMetadata{
		ID:          "BP-43",
		Title:       "Dot Import Outside Tests",
		Description: "A dot import is used in production code.",
		Severity:    rules.SeverityLow,
		Fix:         "Match `import . \"pkg\"` outside `_test.go` files.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-44"] = &rules.RuleMetadata{
		ID:          "BP-44",
		Title:       "Blank Import Without Justification",
		Description: "A blank import appears outside conventional driver or image registration cases.",
		Severity:    rules.SeverityLow,
		Fix:         "Match blank imports without a nearby justification or standard registration pattern.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-45"] = &rules.RuleMetadata{
		ID:          "BP-45",
		Title:       "Inconsistent Receiver Name",
		Description: "Methods on the same receiver type use inconsistent receiver names.",
		Severity:    rules.SeverityInfo,
		Fix:         "Match receiver-name drift across methods on the same type.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-46"] = &rules.RuleMetadata{
		ID:          "BP-46",
		Title:       "HTTP Server Without Timeouts",
		Description: "An HTTP server is configured without read or write timeouts.",
		Severity:    rules.SeverityLow,
		Fix:         "Match `http.Server` literals missing timeout fields.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-47"] = &rules.RuleMetadata{
		ID:          "BP-47",
		Title:       "Missing Graceful Shutdown",
		Description: "A long-running server starts without a graceful shutdown path.",
		Severity:    rules.SeverityLow,
		Fix:         "Match server startup code with no `Shutdown` handling.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-48"] = &rules.RuleMetadata{
		ID:          "BP-48",
		Title:       "log.Fatal Or os.Exit In Library Code",
		Description: "A non-main function exits the process directly.",
		Severity:    rules.SeverityLow,
		Fix:         "Match `log.Fatal*` or `os.Exit` outside main/test entry points.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-49"] = &rules.RuleMetadata{
		ID:          "BP-49",
		Title:       "Deferred Cleanup Without Error Handling",
		Description: "A deferred cleanup drops an error that can affect correctness.",
		Severity:    rules.SeverityLow,
		Fix:         "Match deferred functions that ignore cleanup errors.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-50"] = &rules.RuleMetadata{
		ID:          "BP-50",
		Title:       "Missing Signal Handling In Long-Running Process",
		Description: "A long-running process does not watch SIGTERM or SIGINT for shutdown.",
		Severity:    rules.SeverityLow,
		Fix:         "Match service-style main functions with no signal handling.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-51"] = &rules.RuleMetadata{
		ID:          "BP-51",
		Title:       "Recover Without Re-Panic In Library Code",
		Description: "Library code recovers a panic and continues without preserving failure semantics.",
		Severity:    rules.SeverityLow,
		Fix:         "Match recover handlers in library code that suppress the panic entirely.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-52"] = &rules.RuleMetadata{
		ID:          "BP-52",
		Title:       "Unchecked Integer Multiplication",
		Description: "Arithmetic multiplies values without obvious overflow bounds checks.",
		Severity:    rules.SeverityLow,
		Fix:         "Heuristic for multiplication-heavy code with no guard or size validation.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-53"] = &rules.RuleMetadata{
		ID:          "BP-53",
		Title:       "Gob Registration Mismatch",
		Description: "encoding/gob registration appears inconsistent with the encoded types.",
		Severity:    rules.SeverityLow,
		Fix:         "Match gob.Register patterns that do not line up with later encoder or decoder usage.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-54"] = &rules.RuleMetadata{
		ID:          "BP-54",
		Title:       "Public HTTP Endpoint Without Rate Limiting",
		Description: "A public HTTP endpoint is exposed without any rate-limiting guard.",
		Severity:    rules.SeverityLow,
		Fix:         "Heuristic for public router setup with no rate-limiter middleware.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-55"] = &rules.RuleMetadata{
		ID:          "BP-55",
		Title:       "Missing Request ID Propagation",
		Description: "Middleware chains do not propagate a request identifier through request handling.",
		Severity:    rules.SeverityLow,
		Fix:         "Match middleware stacks that log or branch on requests without request-id propagation.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-56"] = &rules.RuleMetadata{
		ID:          "BP-56",
		Title:       "Deprecated Standard Library Package",
		Description: "Code imports a deprecated package or legacy context path.",
		Severity:    rules.SeverityLow,
		Fix:         "Match imports such as `io/ioutil` or `golang.org/x/net/context`.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-57"] = &rules.RuleMetadata{
		ID:          "BP-57",
		Title:       "Stale Go Version In go.mod",
		Description: "go.mod targets a Go version materially behind the current supported baseline.",
		Severity:    rules.SeverityLow,
		Fix:         "Project-level go.mod parse comparing the declared Go version with policy.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-58"] = &rules.RuleMetadata{
		ID:          "BP-58",
		Title:       "Unpinned Dependency Version",
		Description: "A dependency declaration is too loose to be reproducible.",
		Severity:    rules.SeverityLow,
		Fix:         "Project-level go.mod parse looking for non-specific versions.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-59"] = &rules.RuleMetadata{
		ID:          "BP-59",
		Title:       "Unused Direct Dependency",
		Description: "A direct dependency is listed but not imported anywhere in the project.",
		Severity:    rules.SeverityLow,
		Fix:         "Project-wide dependency-to-import reconciliation.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-60"] = &rules.RuleMetadata{
		ID:          "BP-60",
		Title:       "Test Dependency In Main Module",
		Description: "A dependency only needed for tests lives in the main module dependency set.",
		Severity:    rules.SeverityLow,
		Fix:         "Project-level analysis of go.mod usage versus test-only imports.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-61"] = &rules.RuleMetadata{
		ID:          "BP-61",
		Title:       "Missing Indirect Annotation",
		Description: "An indirect dependency is present without the expected `// indirect` marker.",
		Severity:    rules.SeverityLow,
		Fix:         "Project-level go.mod parse comparing the dependency graph with annotations.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-62"] = &rules.RuleMetadata{
		ID:          "BP-62",
		Title:       "Dependency Used In One File",
		Description: "A dependency is only used once and may not justify an external module.",
		Severity:    rules.SeverityInfo,
		Fix:         "Project-wide import frequency heuristic for low-use dependencies.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-63"] = &rules.RuleMetadata{
		ID:          "BP-63",
		Title:       "Dependency With Known CVE Not Updated",
		Description: "A dependency version remains on a release associated with a known vulnerability.",
		Severity:    rules.SeverityLow,
		Fix:         "Reserved for dependency metadata and vulnerability feed integration.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-64"] = &rules.RuleMetadata{
		ID:          "BP-64",
		Title:       "Local Filesystem Replace Directive",
		Description: "go.mod contains a replace directive that points at the local filesystem.",
		Severity:    rules.SeverityLow,
		Fix:         "Project-level go.mod parse for `replace` directives with local paths.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-65"] = &rules.RuleMetadata{
		ID:          "BP-65",
		Title:       "Missing go.sum Entries",
		Description: "A module references dependencies without complete go.sum coverage.",
		Severity:    rules.SeverityLow,
		Fix:         "Project-level verification that go.mod and go.sum stay in sync.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-66"] = &rules.RuleMetadata{
		ID:          "BP-66",
		Title:       "Wrapped Sentinel Compared Directly",
		Description: "A wrapped sentinel error is compared with == or !=, bypassing the wrapped error chain.",
		Severity:    rules.SeverityLow,
		Fix:         "Match same-function fmt.Errorf calls with %w around an Err-like sentinel followed by direct comparison of the error value; use errors.Is instead.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-67"] = &rules.RuleMetadata{
		ID:          "BP-67",
		Title:       "errors.As Target Not Passed By Address",
		Description: "errors.As receives a target value instead of an addressable pointer and can panic at runtime.",
		Severity:    rules.SeverityHigh,
		Fix:         "Match errors.As calls whose second argument is not visibly address-prefixed; keep the rule limited to the stdlib errors.As call shape.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-68"] = &rules.RuleMetadata{
		ID:          "BP-68",
		Title:       "Discarded errors.Join Result",
		Description: "The combined error returned by errors.Join is discarded, losing the joined failure information.",
		Severity:    rules.SeverityLow,
		Fix:         "Match stdlib errors.Join calls used as expression statements or assigned only to `_`; do not duplicate ordinary discarded-error assignment findings.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-70"] = &rules.RuleMetadata{
		ID:          "BP-70",
		Title:       "Logged Error Then Continued",
		Description: "An error is logged but the error path continues without an explicit exit.",
		Severity:    rules.SeverityLow,
		Fix:         "Match a local `err != nil` branch containing an error logger and no bare return; do not infer the wider function contract.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-72"] = &rules.RuleMetadata{
		ID:          "BP-72",
		Title:       "Typed Nil Interface Return",
		Description: "A typed nil pointer is returned through an error or interface result, producing a non-nil interface value.",
		Severity:    rules.SeverityHigh,
		Fix:         "Match `var p *T = nil` followed by direct `return p` from a visible error or anonymous interface result.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-73"] = &rules.RuleMetadata{
		ID:          "BP-73",
		Title:       "Nil Map Write Without Initialization",
		Description: "A function-local map is indexed before it is initialized with make, which will panic at runtime.",
		Severity:    rules.SeverityHigh,
		Fix:         "Match a local `var name map[...]` declaration followed by an index assignment in the same function with no visible make assignment before the write.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-75"] = &rules.RuleMetadata{
		ID:          "BP-75",
		Title:       "Copy Into Zero-Length Slice",
		Description: "copy writes into a statically zero-length destination, so no source element can be copied.",
		Severity:    rules.SeverityLow,
		Fix:         "Match copy calls whose destination is a local zero-value or make(..., 0) slice and whose source is a non-empty slice literal.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-76"] = &rules.RuleMetadata{
		ID:          "BP-76",
		Title:       "Map Range Used For Ordered Output",
		Description: "Values collected directly from map iteration are used as ordered output even though map iteration is nondeterministic.",
		Severity:    rules.SeverityLow,
		Fix:         "Match a local map ranged into a slice that is later passed to strings.Join without an intervening sort call.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-79"] = &rules.RuleMetadata{
		ID:          "BP-79",
		Title:       "Context Cancellation Not Released",
		Description: "A locally bound context cancellation function has no visible call or defer in the same function.",
		Severity:    rules.SeverityMedium,
		Fix:         "Review-only heuristic: match WithCancel, WithTimeout, or WithDeadline bindings whose cancel identifier is never visibly called or deferred locally.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-80"] = &rules.RuleMetadata{
		ID:          "BP-80",
		Title:       "context.TODO In Production Code",
		Description: "context.TODO leaves production context ownership unresolved outside test files.",
		Severity:    rules.SeverityLow,
		Fix:         "Match exact context.TODO calls outside test files; do not treat this as a proof of a runtime defect.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-81"] = &rules.RuleMetadata{
		ID:          "BP-81",
		Title:       "Repeated time.Now In Condition",
		Description: "A condition reads the clock more than once and can compare values from different instants.",
		Severity:    rules.SeverityLow,
		Fix:         "Match multiple time.Now calls in one boolean condition; avoid standalone timestamp assignments.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-82"] = &rules.RuleMetadata{
		ID:          "BP-82",
		Title:       "Time Parsed Without Location",
		Description: "time.Parse uses the local default location when the application may require an explicit location.",
		Severity:    rules.SeverityLow,
		Fix:         "Match exact time.Parse calls outside test files; this is advisory and does not infer the application's timezone policy.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-83"] = &rules.RuleMetadata{
		ID:          "BP-83",
		Title:       "Sleep Used For Synchronization",
		Description: "time.Sleep is used as a synchronization boundary instead of explicit coordination.",
		Severity:    rules.SeverityLow,
		Fix:         "Match production time.Sleep calls while excluding obvious retry, backoff, jitter, and delay-shaped lines.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-84"] = &rules.RuleMetadata{
		ID:          "BP-84",
		Title:       "Integer Percentage Truncation",
		Description: "Integer division occurs before multiplication by 100 in a percentage-shaped calculation.",
		Severity:    rules.SeverityLow,
		Fix:         "Match simple `a / b * 100` expressions only when the destination or function name visibly indicates a percentage.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-85"] = &rules.RuleMetadata{
		ID:          "BP-85",
		Title:       "Unchecked Request Context Assertion",
		Description: "A net/http handler asserts a request context value without checking whether the value exists and has the expected type.",
		Severity:    rules.SeverityHigh,
		Fix:         "Match single-value Context.Value type assertions inside functions with typed http.ResponseWriter and *http.Request parameters.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-86"] = &rules.RuleMetadata{
		ID:          "BP-86",
		Title:       "Mutex Lock Without Unlock",
		Description: "A locally declared mutex is locked without a visible matching Unlock in the same function.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match sync.Mutex or sync.RWMutex receivers declared locally or as pointer parameters with Lock and no same-receiver Unlock; branch completeness is out of scope.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-87"] = &rules.RuleMetadata{
		ID:          "BP-87",
		Title:       "RLock Held Across Blocking Call",
		Description: "A read lock remains held while an obvious blocking operation executes.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match a declared sync.RWMutex receiver with RLock, a blocking call or channel receive, and a later RUnlock in the same function.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-88"] = &rules.RuleMetadata{
		ID:          "BP-88",
		Title:       "Nil Channel Send Or Receive",
		Description: "A local zero-value channel is used directly outside a select and will block indefinitely.",
		Severity:    rules.SeverityHigh,
		Fix:         "Match direct send/receive operations on a local var ch chan T before a visible make assignment; ignore nil-channel select cases.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-89"] = &rules.RuleMetadata{
		ID:          "BP-89",
		Title:       "Repeated Unconditional Channel Close",
		Description: "The same channel is closed more than once by unconditional statements in one function.",
		Severity:    rules.SeverityHigh,
		Fix:         "Match repeated builtin close calls on the same simple identifier outside conditionals and loops; ownership and reachability beyond this local shape are out of scope.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-90"] = &rules.RuleMetadata{
		ID:          "BP-90",
		Title:       "Channel Receive Loop Without Exit",
		Description: "An infinite channel receive loop has no visible local cancellation or exit path.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match `for {}` loops containing a bare channel receive without select, break, or return.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-91"] = &rules.RuleMetadata{
		ID:          "BP-91",
		Title:       "Data-Bearing Notification Channel",
		Description: "A boolean or integer channel carries a constant notification value while receivers discard the value.",
		Severity:    rules.SeverityLow,
		Fix:         "Match notification-shaped channels that send constants and receive without using the value.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-92"] = &rules.RuleMetadata{
		ID:          "BP-92",
		Title:       "errgroup Without Context",
		Description: "An errgroup is used without the cancellation-aware WithContext form.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match locally declared errgroup.Group values used with Go or Wait and no errgroup.WithContext call.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-93"] = &rules.RuleMetadata{
		ID:          "BP-93",
		Title:       "errgroup Closure Discards Error",
		Description: "An errgroup closure explicitly discards an operation error instead of returning it to the group.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match blank assignment of an error-producing call inside an errgroup Go closure.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-94"] = &rules.RuleMetadata{
		ID:          "BP-94",
		Title:       "Goroutine Map Write Without Sync",
		Description: "A goroutine writes to a shared map without a visible synchronization boundary.",
		Severity:    rules.SeverityHigh,
		Fix:         "Match a goroutine map-index assignment in a function with multiple goroutines and no visible mutex or channel handoff.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-95"] = &rules.RuleMetadata{
		ID:          "BP-95",
		Title:       "HTTP Response Body Without Close",
		Description: "An HTTP response body is acquired without a visible Close call.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match local client Do results without Body.Close; this provides zero-dependency coverage and overlaps bodyclose/sqlclosecheck.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-96"] = &rules.RuleMetadata{
		ID:          "BP-96",
		Title:       "sql.Rows Without Close",
		Description: "A database rows value is acquired without a visible Close call.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match Query or QueryContext results assigned to a rows-like local without a later Close in the same function.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-97"] = &rules.RuleMetadata{
		ID:          "BP-97",
		Title:       "Writer Not Flushed Before Read",
		Description: "A buffered or compressed writer is read through its underlying buffer before being flushed.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match writes through bufio or gzip writers followed by a read of the underlying buffer without Flush or Close.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-98"] = &rules.RuleMetadata{
		ID:          "BP-98",
		Title:       "Opened File Without Close Or Transfer",
		Description: "An os.Open or os.OpenFile result is neither closed nor visibly transferred to the caller.",
		Severity:    rules.SeverityMedium,
		Fix:         "Review-only same-function heuristic for local os.Open/OpenFile bindings with no visible Close or return transfer.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-99"] = &rules.RuleMetadata{
		ID:          "BP-99",
		Title:       "Cond Wait Without Locker",
		Description: "sync.Cond.Wait is used without a visible Lock or RLock on its associated locker.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match locally created sync.Cond values whose same-function body has Wait but no visible Lock/RLock on the constructor locker.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-100"] = &rules.RuleMetadata{
		ID:          "BP-100",
		Title:       "Unbounded Goroutine Fan-Out",
		Description: "A loop launches one goroutine per item without a visible concurrency bound.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match goroutine launches inside range loops without a semaphore, worker pool, or errgroup limit.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-101"] = &rules.RuleMetadata{
		ID:          "BP-101",
		Title:       "HTTP Header Written After Body",
		Description: "A net/http handler writes a response body before setting its intended status header.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match net/http handler functions with ResponseWriter and Request parameters where a writer body call precedes WriteHeader in the same block.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-102"] = &rules.RuleMetadata{
		ID:          "BP-102",
		Title:       "HTTP Error Path Without Response",
		Description: "A net/http handler returns from an error path without writing an error response or status.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match typed net/http handlers with a local err != nil guard and direct bare return when the guarded branch has no response action.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-104"] = &rules.RuleMetadata{
		ID:          "BP-104",
		Title:       "Duplicate ServeMux Pattern",
		Description: "The same literal pattern is registered more than once on a net/http ServeMux.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match repeated literal Handle or HandleFunc registrations on a locally created http.ServeMux.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-105"] = &rules.RuleMetadata{
		ID:          "BP-105",
		Title:       "Sensitive Cookie Missing Security Flags",
		Description: "A sensitive cookie is issued without Secure and HttpOnly protections.",
		Severity:    rules.SeverityHigh,
		Fix:         "Match sensitive cookie literals in net/http handlers missing either Secure or HttpOnly.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-107"] = &rules.RuleMetadata{
		ID:          "BP-107",
		Title:       "HTTP Middleware Missing Next",
		Description: "A net/http middleware returns a handler that neither delegates to the next handler nor writes a terminal response.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match typed net/http middleware closures with no next call and no explicit response action.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-109"] = &rules.RuleMetadata{
		ID:          "BP-109",
		Title:       "Gin Error Response Without Abort",
		Description: "A Gin handler writes an error response and continues without aborting or returning.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match imported Gin handlers with typed *gin.Context parameters where an error JSON response is not followed by Abort or return in its block.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-110"] = &rules.RuleMetadata{
		ID:          "BP-110",
		Title:       "Gin Bind Error Ignored",
		Description: "A Gin binding call is used as a bare statement and its returned error is discarded.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match imported Gin handlers with typed *gin.Context parameters where ShouldBind or ShouldBindJSON is a bare expression statement.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-111"] = &rules.RuleMetadata{
		ID:          "BP-111",
		Title:       "Gin Context In Goroutine",
		Description: "A Gin context is used in a goroutine without a c.Copy lifetime boundary.",
		Severity:    rules.SeverityMedium,
		Fix:         "Require the Gin import and match c. use inside a goroutine without Copy in the same local block; PERF overlap requires review.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-116"] = &rules.RuleMetadata{
		ID:          "BP-116",
		Title:       "Echo Error Response And Raw Error",
		Description: "An Echo handler writes an error response and returns the raw error, creating two response-handling paths.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match imported Echo handlers with typed echo.Context parameters where an error JSON response is followed by a raw error return.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-117"] = &rules.RuleMetadata{
		ID:          "BP-117",
		Title:       "Echo Bind Error Ignored",
		Description: "An Echo binding call is used as a bare statement and its returned error is discarded.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match imported Echo handlers with typed echo.Context parameters where Bind is a bare expression statement.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-119"] = &rules.RuleMetadata{
		ID:          "BP-119",
		Title:       "Fiber Context In Goroutine",
		Description: "A Fiber context is captured in a goroutine without a visible lifetime boundary.",
		Severity:    rules.SeverityMedium,
		Fix:         "Require the Fiber import and match c. or ctx. use inside a goroutine without Immutable or Copy; PERF overlap requires review.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-120"] = &rules.RuleMetadata{
		ID:          "BP-120",
		Title:       "Fiber BodyParser Error Ignored",
		Description: "A Fiber body parser is used as a bare statement and its returned error is discarded.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match imported Fiber handlers with typed *fiber.Ctx parameters where BodyParser is a bare expression statement.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-122"] = &rules.RuleMetadata{
		ID:          "BP-122",
		Title:       "Chi Middleware Missing Next",
		Description: "A Chi middleware returns without calling the next handler or writing a terminal response.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match imported Chi middleware closures with no next.ServeHTTP call and no explicit response action.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-126"] = &rules.RuleMetadata{
		ID:          "BP-126",
		Title:       "Transaction Without Completion",
		Description: "A database transaction has no visible Commit or Rollback path.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match a local database/sql Begin or BeginTx call whose containing function has no visible Commit or Rollback; ownership transfer is out of scope.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-128"] = &rules.RuleMetadata{
		ID:          "BP-128",
		Title:       "QueryRow Scan Without ErrNoRows Handling",
		Description: "QueryRow.Scan errors are handled without distinguishing sql.ErrNoRows.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match typed database/sql QueryRow calls whose Scan error path has no ErrNoRows handling.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-131"] = &rules.RuleMetadata{
		ID:          "BP-131",
		Title:       "Query Used For DML Without Rows",
		Description: "database/sql Query is used for a literal DML statement that does not return rows.",
		Severity:    rules.SeverityLow,
		Fix:         "Review-only import-gated heuristic for literal INSERT, UPDATE, or DELETE passed to Query or QueryContext without RETURNING.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-132"] = &rules.RuleMetadata{
		ID:          "BP-132",
		Title:       "Update Without RowsAffected Check",
		Description: "An optimistic-lock-shaped update ignores RowsAffected and cannot detect a lost update.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match UPDATE execution with a version predicate and no visible RowsAffected check.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-133"] = &rules.RuleMetadata{
		ID:          "BP-133",
		Title:       "GORM Chain Error Ignored",
		Description: "A GORM chain result is used without checking its Error field.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match typed GORM chains whose result is consumed without a visible Error check.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-134"] = &rules.RuleMetadata{
		ID:          "BP-134",
		Title:       "GORM First Without Not Found Handling",
		Description: "A GORM First or Take query does not distinguish gorm.ErrRecordNotFound.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match typed GORM First or Take error paths with no ErrRecordNotFound handling.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-135"] = &rules.RuleMetadata{
		ID:          "BP-135",
		Title:       "GORM Global Without Session",
		Description: "A package-level GORM handle is chained directly in a request path without a session boundary.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match imported GORM package globals used in typed request handlers without Session or WithContext.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-136"] = &rules.RuleMetadata{
		ID:          "BP-136",
		Title:       "GORM AutoMigrate In Request Path",
		Description: "GORM schema migration is executed from a request handler instead of startup or a migration command.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match AutoMigrate on a typed *gorm.DB parameter in a function that also accepts net/http or recognized framework request parameters.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-138"] = &rules.RuleMetadata{
		ID:          "BP-138",
		Title:       "External I/O In GORM Hook",
		Description: "A GORM lifecycle hook performs direct external HTTP or SMTP I/O before the database operation is safely committed.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match recognized GORM hook methods with a *gorm.DB parameter and direct imported http or smtp calls; helper ownership is out of scope.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-140"] = &rules.RuleMetadata{
		ID:          "BP-140",
		Title:       "sqlx Retrieval Error Ignored",
		Description: "A sqlx StructScan or Get retrieval call is used as a bare expression and its error is discarded.",
		Severity:    rules.SeverityLow,
		Fix:         "Match imported sqlx retrieval calls used as expression statements.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-141"] = &rules.RuleMetadata{
		ID:          "BP-141",
		Title:       "sqlx Named Placeholder Missing Struct Tag",
		Description: "A sqlx named query uses a snake_case placeholder that an untagged Go field will not map to under the default mapper.",
		Severity:    rules.SeverityLow,
		Fix:         "Match typed sqlx named execution with a same-file struct/value binding, snake_case placeholders, and no matching db tag; custom mapper configuration is not inferred.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-142"] = &rules.RuleMetadata{
		ID:          "BP-142",
		Title:       "sqlx.In Query Without Rebind",
		Description: "A query expanded by sqlx.In is executed without rebinding placeholders for the database driver.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match same-function sqlx.In output reaching a typed *sqlx.DB or *sqlx.Tx execution call without an intervening receiver Rebind call.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-143"] = &rules.RuleMetadata{
		ID:          "BP-143",
		Title:       "Redis Result Error Ignored",
		Description: "A go-redis command result is used as a bare expression and its error is discarded.",
		Severity:    rules.SeverityLow,
		Fix:         "Match imported go-redis command calls used as expression statements without checking Result().Err().",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-145"] = &rules.RuleMetadata{
		ID:          "BP-145",
		Title:       "pgx Pool Connection Not Released",
		Description: "A pgx pool connection is acquired without a visible Release or Close path.",
		Severity:    rules.SeverityMedium,
		Fix:         "Review-only import-gated same-function heuristic for pgxpool Acquire results with no later Release or Close.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-146"] = &rules.RuleMetadata{
		ID:          "BP-146",
		Title:       "Sensitive Fields Logged",
		Description: "A logger receives a sensitive field or value without an obvious redaction step.",
		Severity:    rules.SeverityHigh,
		Fix:         "Match sensitive field names or environment values passed to log, slog, or zap calls without redaction.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-147"] = &rules.RuleMetadata{
		ID:          "BP-147",
		Title:       "Unstructured Service Logging",
		Description: "A non-main service package mixes the standard logger with a structured logger.",
		Severity:    rules.SeverityLow,
		Fix:         "Match log.Print calls in non-main packages that also import log/slog or zap.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-149"] = &rules.RuleMetadata{
		ID:          "BP-149",
		Title:       "Error Log Without Attribute",
		Description: "An error-level logger call inside an error branch omits the error as a structured attribute.",
		Severity:    rules.SeverityLow,
		Fix:         "Match err != nil branches with error-level slog or zap calls that do not include err.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-151"] = &rules.RuleMetadata{
		ID:          "BP-151",
		Title:       "Secret Environment Value Logged",
		Description: "A sensitive environment value is passed directly to a logger and may expose a secret in logs.",
		Severity:    rules.SeverityHigh,
		Fix:         "Match direct os.Getenv calls for secret-like names passed to stdlib log, log/slog, or explicit zap.L/zap.S logger calls.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-154"] = &rules.RuleMetadata{
		ID:          "BP-154",
		Title:       "Discarded json.Unmarshal Error",
		Description: "json.Unmarshal is used without handling its returned error.",
		Severity:    rules.SeverityLow,
		Fix:         "Match exact json.Unmarshal expression statements or blank assignments outside test files; document overlap with generic ignored-error rules.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-155"] = &rules.RuleMetadata{
		ID:          "BP-155",
		Title:       "Unbounded JSON Request Body",
		Description: "A JSON decoder reads an HTTP request body without a visible size limit.",
		Severity:    rules.SeverityHigh,
		Fix:         "Match json.Decoder on http.Request.Body without http.MaxBytesReader or an equivalent bounded reader.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-156"] = &rules.RuleMetadata{
		ID:          "BP-156",
		Title:       "Security-Sensitive JSON Omitempty",
		Description: "A security-sensitive JSON field relies on omitempty and can silently disappear when zero-valued.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match sensitive struct fields with json omitempty tags in request or security-shaped types.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-158"] = &rules.RuleMetadata{
		ID:          "BP-158",
		Title:       "Naked gRPC Error Return",
		Description: "A gRPC error is returned naked despite importing status helpers.",
		Severity:    rules.SeverityMedium,
		Fix:         "Require google.golang.org/grpc/status and match return err lines when status helpers are present; do not infer service-wide error policy.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-159"] = &rules.RuleMetadata{
		ID:          "BP-159",
		Title:       "Flag Value Read Before Parse",
		Description: "A flag pointer is dereferenced before flag.Parse processes command-line arguments.",
		Severity:    rules.SeverityLow,
		Fix:         "Match flag pointer declarations and dereferences before the first flag.Parse call in the same function.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-160"] = &rules.RuleMetadata{
		ID:          "BP-160",
		Title:       "Cobra Run Without RunE",
		Description: "A Cobra command uses Run instead of RunE and may swallow command errors.",
		Severity:    rules.SeverityLow,
		Fix:         "Require the Cobra import and match command literals with Run but no RunE; keep advisory because command policy may be intentional.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-161"] = &rules.RuleMetadata{
		ID:          "BP-161",
		Title:       "Test Uses Production DSN",
		Description: "A test opens a literal database target containing an explicit production marker.",
		Severity:    rules.SeverityHigh,
		Fix:         "Match test files with imported database/sql or GORM calls whose literal sql.Open or gorm.Open text contains a standalone prod or production marker; local/container targets are ignored.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-162"] = &rules.RuleMetadata{
		ID:          "BP-162",
		Title:       "Parallel Test Mutates Shared State",
		Description: "A parallel test mutates package-level state that can race with other tests.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match Test, Benchmark, or Fuzz functions calling t.Parallel and assigning to package-level variables in the same test file.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-163"] = &rules.RuleMetadata{
		ID:          "BP-163",
		Title:       "Golden Update Without Short Guard",
		Description: "A golden-file update path writes output without skipping short tests.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match test update flags and golden writes in an update branch when the test function has no testing.Short call.",
		Pack:        rules.PackBadPractice,
	}
	metaByID["BP-164"] = &rules.RuleMetadata{
		ID:          "BP-164",
		Title:       "Functional Option Mutates Global Default",
		Description: "An exported functional option changes package-level default state instead of only configuring its supplied instance.",
		Severity:    rules.SeverityMedium,
		Fix:         "Match exported With* functions returning an Option whose body assigns to a package-level variable.",
		Pack:        rules.PackBadPractice,
	}
}

// MetadataForID returns catalogue metadata for a BP rule id.
func MetadataForID(ruleID string) *rules.RuleMetadata {
	return metaByID[ruleID]
}

// CatalogueSize returns the number of BP rules in the embedded catalogue.
func CatalogueSize() int { return len(metaByID) }
