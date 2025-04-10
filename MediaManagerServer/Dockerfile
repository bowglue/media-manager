# Use the latest official Go image
FROM golang:1.24.1

WORKDIR /app

# Install FFmpeg and other dependencies
RUN apt-get update && apt-get install -y \
    ffmpeg \
    sqlite3 \
    libsqlite3-dev \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Install development tools
RUN go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install github.com/codegangsta/gin@latest


# Expose the application port
EXPOSE 8080
