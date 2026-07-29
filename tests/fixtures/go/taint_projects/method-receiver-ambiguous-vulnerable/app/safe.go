package app

// Safe.Open shares the bare method name with Sink.Open but is not a sink.
// Present so the package is intentionally ambiguous without receiver typing.
type Safe struct{}

func (s *Safe) Open(p string) {
	_ = p
}
