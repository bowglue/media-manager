FROM media-manager-server-base:latest

# # Initialize the Go module (if not already initialized)
RUN go mod init proto-gen || true  # Skip if already initialized

# Install protoc and required tools
RUN apt-get update && apt-get install -y protobuf-compiler

# Install Go protobuf plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

# WORKDIR /proto

# # Create directories
# RUN mkdir -p definitions generated