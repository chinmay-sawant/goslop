package perf

// Seed rules (Phase 6a) register via init so batch files can append freely.
func init() {
	RegisterRule("PERF-1", detectPERF1, &MetaPERF1)
	RegisterRule("PERF-2", detectPERF2, &MetaPERF2)
	RegisterRule("PERF-3", detectPERF3, &MetaPERF3)
	RegisterRule("PERF-4", detectPERF4, &MetaPERF4)
	RegisterRule("PERF-5", detectPERF5, &MetaPERF5)
	RegisterRule("PERF-6", detectPERF6, &MetaPERF6)
	RegisterRule("PERF-7", detectPERF7, &MetaPERF7)
	RegisterRule("PERF-8", detectPERF8, &MetaPERF8)
	RegisterRule("PERF-32", detectPERF32, &MetaPERF32)
	RegisterRule("PERF-50", detectPERF50, &MetaPERF50)
	RegisterRule("PERF-116", detectPERF116, &MetaPERF116)
	RegisterRule("PERF-230", detectPERF230, &MetaPERF230)
}
