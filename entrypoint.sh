#!/bin/sh
set -e

POSTGRES_DSN="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"

# Run migrations with .env file parametres
echo "📦 Running $GOOSE_DRIVER DB migrations: GOOSE_MIGRATION_DIR=$GOOSE_MIGRATION_DIR | GOOSE_TABLE=$GOOSE_TABLE"
goose $POSTGRES_DSN up

# Start the application
./app
