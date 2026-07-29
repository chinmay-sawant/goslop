package app

// Safe.Open shares the bare method name with Sink.Open but is not a sink.
type Safe struct{}

func (s *Safe) Open(p string) {
	_ = p
}
