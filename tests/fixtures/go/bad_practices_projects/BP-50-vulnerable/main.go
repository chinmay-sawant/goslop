package main

import "net/http"

func main() {
	server := &http.Server{Addr: ":8080"}
	if err := server.ListenAndServe(); err != nil {
		return
	}
}
