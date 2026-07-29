// Package main demonstrates signal handling for server shutdown.
package main

import (
	"context"
	"net/http"
	"os/signal"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background())
	defer stop()

	server := &http.Server{
		Addr:         ":8080",
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
