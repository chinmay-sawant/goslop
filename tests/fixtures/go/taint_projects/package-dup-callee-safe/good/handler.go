package good

import "net/http"

// openPath shares the bare name with package bad but is not a sink.
func openPath(s string) {
	_ = s
}

func handler(r *http.Request) {
	x := r.URL.Query().Get("input")
	openPath(x)
}
