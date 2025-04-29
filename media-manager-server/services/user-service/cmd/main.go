package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "shared/proto/api" // Import your generated proto package

	"google.golang.org/grpc"
)

// Implement the UserServiceServer interface
type userServiceServer struct {
	pb.UnimplementedUserServiceServer
}

// Implement the GetUser RPC
func (s *userServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	log.Printf("Received GetUser request for user_id: %s", req.GetUserId())

	// Mock user data (normally you would query a database)
	return &pb.GetUserResponse{
		UserId:   req.GetUserId(),
		Username: "testuser_" + req.GetUserId(),
	}, nil
}

func main() {
	// Create a TCP listener
	lis, err := net.Listen("tcp", ":12000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer lis.Close()
	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the UserService server
	pb.RegisterUserServiceServer(grpcServer, &userServiceServer{})

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()

	log.Println("UserService gRPC server is listening on :12000...")

	// Start serving
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
