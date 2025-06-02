# FROM golang:1.24.1

# COPY go.mod go.sum ./
# RUN go mod download

# RUN apt-get update && apt-get install -y \
#     ffmpeg \
#     sqlite3 \
#     libsqlite3-dev \
#     && apt-get clean && rm -rf /var/lib/apt/lists/*

# RUN go install github.com/cosmtrek/air@v1.52.0
# ENV PATH="/go/bin:$PATH"

# RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
# RUN go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
# WORKDIR /tmp/gqlgen-install
# RUN go mod init temp && \
#     go install github.com/99designs/gqlgen@v0.17.73

FROM golang:1.24.1

# Set PATH early
ENV PATH="/go/bin:$PATH"

# Copy module files and download deps
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Install OS dependencies
RUN apt-get update && apt-get install -y \
    ffmpeg \
    sqlite3 \
    libsqlite3-dev \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# Install dev tools
RUN go install github.com/cosmtrek/air@v1.52.0
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest


# Install gqlgen using a temp module to avoid go.sum issues
WORKDIR /tmp/gqlgen-install
RUN go mod init temp && \
    go install github.com/99designs/gqlgen@v0.17.73

# Return to working directory for app code
WORKDIR /app
COPY . .