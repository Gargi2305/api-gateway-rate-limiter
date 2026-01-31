package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-gateway/internal/server"
	"api-gateway/internal/redisclient"
	"api-gateway/internal/metrics"
)

func main() {

	redisclient.Init()
	metrics.Init()

	srv := server.New()

	log.Println("API Gateway starting on port 8080...")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    <-quit

	log.Println("Shutting down gateway...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Gateway forced to shutdown: %v", err)
	}

	log.Println("Gateway exited cleanly")

}
