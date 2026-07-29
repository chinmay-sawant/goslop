package app

import (
	"net/http"
	"os"
)

// Sink.Open is a FileOpen sink. A second Safe.Open exists in the package so
// bare-name resolution alone is ambiguous; the call below is on the enclosing
// method's own receiver and must still resolve to *Sink.
type Sink struct{}

func (s *Sink) Open(p string) {
	os.Open(p)
}

func (s *Sink) Serve(r *http.Request) {
	x := r.URL.Query().Get("input")
	s.Open(x)
}
