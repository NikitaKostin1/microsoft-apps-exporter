#!/bin/sh

# Export variables from .env file
export $(grep -v '^#' .env | xargs)

# Construct the DSN from individual environment variables
POSTGRES_DSN="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"

# Run migrations using goose
goose postgres $POSTGRES_DSN up --dir ./migrations --table sharepoint_goose_db_version

# Start the main application
./cache/main
