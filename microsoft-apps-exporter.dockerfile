FROM golang:1.20-alpine

WORKDIR /app

# Copy the Go application source code
COPY . .

# Build the Go application
RUN go build -o app .

# Copy migration files into the image
COPY migrations /app/migrations

# Copy the entrypoint script
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
