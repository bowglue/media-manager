package main

import (
	"context"
	"mms/common/logger"
	"mms/core/server"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log := logger.NewLogger()

	config := server.ServerConfig{
		GatewayPort: ":8080", // Port for the HTTP server
		Log:         log,
	}

	// Initialize the server with the configuration
	srv := server.NewServer(config)

	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Error("Server error: ", err)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Failed to shutdown server cleanly: ", err)
	}
}
