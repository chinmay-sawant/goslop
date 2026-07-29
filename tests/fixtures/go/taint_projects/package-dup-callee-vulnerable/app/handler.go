package app

import "net/http"

func handler(r *http.Request) {
	x := r.URL.Query().Get("input")
	openPath(x)
}
