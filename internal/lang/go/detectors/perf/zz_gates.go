package perf

// Wire SourceIndex gates for high-cost / high-hit PERF rules after all
// RegisterRule inits (file name zz_ ensures late init order).
// Gates are any-of and must appear in perfNeedles. FN-safe only when the
// detector truly requires those substrings (or APIs) to fire.
func init() {
	// Top product / seed rules with clear API requirements.
	RegisterRule("PERF-1", detectPERF1, &MetaPERF1, "regexp.MustCompile", "regexp.Compile", "regexp.MatchString")
	RegisterRule("PERF-5", detectPERF5, &MetaPERF5, "json.Marshal", "json.Unmarshal", "json.NewEncoder", "json.NewDecoder")
	RegisterRule("PERF-6", detectPERF6, &MetaPERF6, "fmt.Sprintf(", "fmt.Fprintf(")
	RegisterRule("PERF-32", detectPERF32, &MetaPERF32, "[]byte(", "string(")

	// Batch API gates (match detector early Contains / callee checks).
	// PERF-21 = io.ReadAll on request body (not reflect — that was a mis-map).
	RegisterRule("PERF-16", detectPERF16, &MetaPERF16, "bytes.Buffer{", "new(bytes.Buffer)")
	RegisterRule("PERF-21", detectPERF21, &MetaPERF21, "io.ReadAll(")
	RegisterRule("PERF-22", detectPERF22, &MetaPERF22, "os.ReadFile(", "ioutil.ReadFile(")
	RegisterRule("PERF-23", detectPERF23, &MetaPERF23, "bytes.NewReader(", "bytes.NewBuffer(")
	RegisterRule("PERF-28", detectPERF28, &MetaPERF28, "sync.Mutex", "sync.RWMutex")
	RegisterRule("PERF-35", detectPERF35, &MetaPERF35, "fmt.Sprintf(", "fmt.Errorf(", "fmt.Fprintf(")
	RegisterRule("PERF-46", detectPERF46, &MetaPERF46, "strings.Trim", "TrimSpace", "TrimPrefix")
	RegisterRule("PERF-47", detectPERF47, &MetaPERF47, "strings.Split")
	RegisterRule("PERF-101", detectPERF101, &MetaPERF101, "http.Server{", "&http.Server{", "ListenAndServe")
	RegisterRule("PERF-122", detectPERF122, &MetaPERF122, "HasPrefix", "TrimPrefix", "strings.HasPrefix")
	RegisterRule("PERF-186", detectPERF186, &MetaPERF186, "strings.Fields")
}
