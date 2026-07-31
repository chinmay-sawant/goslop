package perf

import "github.com/chinmay-sawant/goslop/internal/rules"

func perfMeta(id, title string) *rules.RuleMetadata {
	return &rules.RuleMetadata{
		ID: id, Title: title, Description: title + " can add avoidable work on a Python service path.",
		Severity: rules.SeverityMedium, Pack: rules.PackPerformance,
	}
}

var metaByID = map[string]*rules.RuleMetadata{
	"PERF-PY-1":  perfMeta("PERF-PY-1", "Full Result Set Materialized Before App-Side Sort"),
	"PERF-PY-2":  perfMeta("PERF-PY-2", "Django ORM Lookup Inside Item Loop"),
	"PERF-PY-3":  perfMeta("PERF-PY-3", "Django Per-Row Create In Batch Loop"),
	"PERF-PY-4":  perfMeta("PERF-PY-4", "Django Read-Modify-Write Counter Update"),
	"PERF-PY-5":  perfMeta("PERF-PY-5", "Sequential Blocking Delivery Over Claimed Batch"),
	"PERF-PY-6":  perfMeta("PERF-PY-6", "ORM Work Claim Without Row Lock"),
	"PERF-PY-7":  perfMeta("PERF-PY-7", "BaseHTTPMiddleware On FastAPI Request Path"),
	"PERF-PY-8":  perfMeta("PERF-PY-8", "SQLAlchemy Lazy Relationship Access In Batch Loop"),
	"PERF-PY-9":  perfMeta("PERF-PY-9", "JSON Payload Parse Then Re-Serialize"),
	"PERF-PY-10": perfMeta("PERF-PY-10", "Worker Sleeps After Successful Batch"),
	"PERF-PY-11": perfMeta("PERF-PY-11", "Per-Row ORM Mutation Instead Of Set-Based Update"),
	"PERF-PY-12": perfMeta("PERF-PY-12", "Unbounded JSON Request Body Parse"),
	"PERF-PY-13": perfMeta("PERF-PY-13", "Full ORM Hydration For Projection Read"),
	"PERF-PY-14": perfMeta("PERF-PY-14", "Select-Then-Insert Idempotency Check"),
	"PERF-PY-15": perfMeta("PERF-PY-15", "ORM Composite Filter Without Supporting Index"),
	"PERF-PY-16": perfMeta("PERF-PY-16", "Retention Timestamp Predicate Without Index"),
	"PERF-PY-17": perfMeta("PERF-PY-17", "Database Connection Reuse Or Timeout Controls Missing"),
	"PERF-PY-18": perfMeta("PERF-PY-18", "Repeated Regex Rewrites On The Same Input"),
	"PERF-PY-19": perfMeta("PERF-PY-19", "Unbounded ORM Locking Sweep"),
	"PERF-PY-20": perfMeta("PERF-PY-20", "ORM Sort Without Supporting Composite Index"),
	"PERF-PY-21": perfMeta("PERF-PY-21", "Unbounded Bulk Delete In Maintenance Path"),
	"PERF-PY-22": perfMeta("PERF-PY-22", "SQLite Backend For Concurrent Service Writes"),
}

// MetadataForID returns metadata for an implemented PERF-PY rule.
func MetadataForID(id string) *rules.RuleMetadata { return metaByID[id] }

// CatalogueSize returns the PERF-PY metadata entry count.
func CatalogueSize() int { return len(metaByID) }
