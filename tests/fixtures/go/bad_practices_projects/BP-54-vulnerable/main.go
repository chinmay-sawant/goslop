// Package main demonstrates a public HTTP endpoint without rate limiting.
package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5,
		WriteTimeout: 5,
	}
	if err := server.ListenAndServe(); err != nil {
		return
	}
}
