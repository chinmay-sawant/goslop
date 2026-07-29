// Package main logs request traffic without request-id propagation.
package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		w.WriteHeader(http.StatusOK)
	})
	if err := http.ListenAndServe(":8080", mux); err != nil {
		return
	}
}
