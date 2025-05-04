FROM golang:1.24.1

# COPY go.mod go.sum ./
# RUN go mod download

RUN apt-get update && apt-get install -y \
    ffmpeg \
    sqlite3 \
    libsqlite3-dev \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

RUN go install github.com/cosmtrek/air@v1.52.0
ENV PATH="/go/bin:$PATH"

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
