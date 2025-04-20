package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	// "time"

	// "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gateway/config"
	"gateway/handlers"

	"google.golang.org/grpc"

	// Import gRPC client stubs
	userpb "proto/generated/user-service"
)

// Main function to run the API Gateway
func main() {
	// Load config (could be port, service URLs, etc.)
	cfg := config.LoadConfig()

	// Set up the gRPC server
	grpcServer := grpc.NewServer()

	// Register gRPC services
	userpb.RegisterUserServiceServer(grpcServer, &handlers.UserService{})

	// Start the gRPC server in a goroutine (to run alongside HTTP server)
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Printf("gRPC server started on port %d", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC server: %v", err)
		}
	}()

	// Set up the HTTP server for the gRPC-Gateway (REST API)
	mux := http.NewServeMux()

	// Register handlers
	mux.HandleFunc("/users", handlers.HandleUserRequest)

	// Start the HTTP server (this is the API Gateway's external-facing endpoint)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler: mux,
	}

	log.Printf("HTTP server started on port %d", cfg.HTTPPort)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("failed to serve HTTP server: %v", err)
	}
}
