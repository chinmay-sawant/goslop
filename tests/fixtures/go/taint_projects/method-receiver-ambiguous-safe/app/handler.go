package app

import "net/http"

// Call Open on a Safe receiver from a free function so the call-site receiver
// type cannot be inferred from an enclosing method parameter. With two
// same-package Open methods, resolution must decline rather than inherit
// Sink.Open's summary.
func handler(r *http.Request) {
	x := r.URL.Query().Get("input")
	s := &Safe{}
	s.Open(x)
}
