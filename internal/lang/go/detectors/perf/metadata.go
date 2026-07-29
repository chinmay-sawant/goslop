package perf

import "github.com/chinmay/codehound/internal/rules"

// MetaPERF116 is catalogue metadata for strings.Index used as Contains.
// Registry: general_perf / PERF-116. Severity Info (micro-opt tier B/C class).
var MetaPERF116 = rules.Meta(
	"PERF-116",
	"strings.Index Used For Contains Check",
	"Uses strings.Index(s, sub) != -1 to check substring containment. strings.Contains(s, sub) communicates intent more clearly and has equivalent performance, but Index comparison patterns are often vestigial and obscure the boolean check.",
	rules.SeverityInfo,
	nil,
	"Replace strings.Index(s, sub) != -1 with strings.Contains(s, sub) (or == -1 with !strings.Contains).",
)
