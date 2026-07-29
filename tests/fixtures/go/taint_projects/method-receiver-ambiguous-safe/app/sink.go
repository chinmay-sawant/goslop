package app

import "os"

// Sink.Open shares the bare method name with Safe.Open but is a FileOpen sink.
type Sink struct{}

func (s *Sink) Open(p string) {
	os.Open(p)
}
