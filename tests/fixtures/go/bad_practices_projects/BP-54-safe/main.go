// Package main demonstrates a rate-limited HTTP endpoint with shutdown handling.
package main

import (
	"context"
	"net/http"
	"os/signal"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background())
	defer stop()

	limiter := rate.NewLimiter(1, 4)
	mux := http.NewServeMux()
	mux.HandleFunc("/status", func(w http.ResponseWriter, _ *http.Request) {
		status := http.StatusOK
		if !limiter.Allow() {
			status = http.StatusTooManyRequests
		}
		w.WriteHeader(status)
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	go func() {
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()
	if err := server.ListenAndServe(); err != nil {
		return
	}
}
