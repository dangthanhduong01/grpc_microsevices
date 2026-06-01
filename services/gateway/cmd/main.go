package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"services/gateway/config"
	httpadapter "services/gateway/internal/adapters/http"
)

func main() {
	gatewayAdapter := httpadapter.NewAdapter(
		config.GetAuthServiceURL(),
		config.GetOrderServiceURL(),
	)

	// Start server in goroutine
	go func() {
		log.Printf("Starting HTTP gateway on port %d", config.GetApplicationPort())
		if err := gatewayAdapter.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := gatewayAdapter.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
