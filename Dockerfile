# Stage 1: Build the Application
FROM golang:1.24.2-alpine3.21 AS builder

WORKDIR /app

# Install Goose with pinned version and override insecure x/crypto dependency
WORKDIR /tmp/goosebuild
RUN go mod init temp && \
    go get github.com/pressly/goose/v3/cmd/goose@v3.24.1 && \
    go get golang.org/x/crypto@v0.36.0 && \
    GOFLAGS="-tags=no_mysql,no_sqlite3,no_mssql,no_libsql,no_vertica,no_ydb -ldflags=-s -w" \
    go build -o /app/goose github.com/pressly/goose/v3/cmd/goose && \
    rm -rf /tmp/goosebuild

WORKDIR /app

# Download dependencies with cache mount
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download -x

# Copy only necessary source files
COPY cmd/ ./cmd/ 
COPY internal/ ./internal/
COPY migrations/ ./migrations/
COPY entrypoint.sh .

# Build the application with optimizations
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux \
    go build -ldflags="-w -s -extldflags=-static" -o=app ./cmd/main.go

# Stage 2: Create the Final Runtime Image
FROM alpine:3.21 AS final

WORKDIR /app

# Copy only the necessary files from builder
COPY --from=builder /app/app .
COPY --from=builder /app/goose /usr/local/bin/goose
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/entrypoint.sh .
COPY resources.yaml .

# Expose the application port
EXPOSE $WEBHOOK_LISTEN_PORT

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN chown -R appuser:appgroup /app
USER appuser

# Set the entrypoint script to start the application
RUN chmod +x /app/entrypoint.sh
ENTRYPOINT ["./entrypoint.sh"]
