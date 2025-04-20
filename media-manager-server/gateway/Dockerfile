# Use the latest official Go image
FROM media-manager-server-base:latest

# WORKDIR /app

# # Install FFmpeg, SQLite, Protobuf compiler, and other dependencies
# RUN apt-get update && apt-get install -y \
#     ffmpeg \
#     sqlite3 \
#     libsqlite3-dev \
#     protobuf-compiler \
#     && apt-get clean \
#     && rm -rf /var/lib/apt/lists/*

# # Copy go.mod and go.sum files
# COPY go.mod go.sum ./
# RUN go mod download

# # Install development tools
# RUN go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
# RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
# RUN go install github.com/codegangsta/gin@latest

# # Install Protobuf and gRPC tools
# RUN go install google.golang.org/grpc/cmd/protoc-gen-go@v1.36.0
# RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
# RUN go install github.com/golang/protobuf/protoc-gen-go@v1.5.2

# # Expose the application port
# EXPOSE 8080
