package perf

import "github.com/chinmay/goslop/internal/rules"

// Catalogue metadata for the PERF rules implemented in this package.
// Titles/descriptions align with Rust registry / product docs; fix strings
// mirror metadata_overrides where present.

var (
	MetaPERF1 = rules.Meta(
		"PERF-1",
		"Regexp Compiled Inside Loop",
		"regexp.MustCompile / regexp.Compile is invoked inside a loop body, recompiling the pattern every iteration.",
		rules.SeverityMedium,
		nil,
		"Hoist regexp.MustCompile outside the loop and reuse the *Regexp.",
	)
	MetaPERF2 = rules.Meta(
		"PERF-2",
		"String Concatenation In Loop",
		"A string is built by repeated concatenation (+=) inside a loop body.",
		rules.SeverityInfo, // C-tier (Rust severity_for_tier)
		nil,
		"Use strings.Builder, bytes.Buffer, or strings.Join instead of repeated +=.",
	)
	MetaPERF3 = rules.Meta(
		"PERF-3",
		"Slice Rebuilt Inside Loop",
		"A working slice is reallocated with make inside a loop body.",
		rules.SeverityInfo, // C-tier
		nil,
		"Hoist the make outside the loop or preallocate once with a capacity hint.",
	)
	MetaPERF4 = rules.Meta(
		"PERF-4",
		"Map Allocated Inside Loop",
		"A map is allocated with make inside a loop body without a size hint.",
		rules.SeverityInfo, // C-tier
		nil,
		"Hoist the map allocation outside the loop or pre-size with make(map[K]V, n).",
	)
	MetaPERF5 = rules.Meta(
		"PERF-5",
		"JSON Conversion Inside Loop",
		"json.Marshal / Unmarshal / NewEncoder / NewDecoder is used inside a loop body.",
		rules.SeverityMedium,
		nil,
		"Hoist encoder/decoder construction or batch the values before a single marshal.",
	)
	MetaPERF6 = rules.Meta(
		"PERF-6",
		"Fmt Formatting Inside Loop",
		"fmt.Sprintf / fmt.Fprintf is used for string construction inside a loop body.",
		rules.SeverityInfo, // C-tier (Rust severity_for_tier)
		nil,
		"Use a bytes.Buffer, strings.Builder, or pool of buffers to avoid repeated fmt allocations.",
	)
	MetaPERF7 = rules.Meta(
		"PERF-7",
		"Defer Inside Loop",
		"A defer statement is placed inside a loop body, delaying cleanup until the function returns.",
		rules.SeverityMedium,
		nil,
		"Replace the defer with explicit close calls inside the loop or move the work into a helper function.",
	)
	MetaPERF8 = rules.Meta(
		"PERF-8",
		"time.Parse Inside Loop",
		"time.Parse / time.ParseInLocation is called inside a loop with a literal layout string.",
		rules.SeverityLow,
		nil,
		"Parse once outside the loop when the layout is fixed, or use a precompiled layout path.",
	)
	MetaPERF32 = rules.Meta(
		"PERF-32",
		"String Byte Conversion On Hot Path",
		"string <-> []byte conversion copies the underlying data on a hot path or in a loop.",
		rules.SeverityMedium,
		nil,
		"Use unsafe conversions only in measured hot paths, or hoist the conversion outside the loop with a pooled buffer.",
	)
	MetaPERF50 = rules.Meta(
		"PERF-50",
		"Regexp Match Inside Loop",
		"regexp.MatchString / Match / MatchReader is invoked inside a loop, recompiling each time.",
		rules.SeverityMedium,
		nil,
		"Compile the regex once with regexp.MustCompile and call re.MatchString in the loop.",
	)
	MetaPERF116 = rules.Meta(
		"PERF-116",
		"strings.Index Used For Contains Check",
		"Uses strings.Index(s, sub) != -1 to check substring containment. strings.Contains(s, sub) communicates intent more clearly and has equivalent performance, but Index comparison patterns are often vestigial and obscure the boolean check.",
		rules.SeverityInfo,
		nil,
		"Replace strings.Index(s, sub) != -1 with strings.Contains(s, sub) (or == -1 with !strings.Contains).",
	)
	MetaPERF230 = rules.Meta(
		"PERF-230",
		"Loop-Invariant Pure Call",
		"A pure/helper function is re-evaluated every loop iteration with stable arguments.",
		rules.SeverityMedium, // unclassified PERF → Medium (Rust)
		nil,
		"Hoist the pure call before the loop or cache its result when arguments do not change across iterations.",
	)
)

// metaByID indexes all implemented PERF metadata by rule id (filled by RegisterRule).
var metaByID = map[string]*rules.RuleMetadata{}
