package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"microsoft-apps-exporter/internal/configuration"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type Database struct {
	Connection *sql.DB
}

// NewDatabase initializes a new Database instance and establishes a connection to the PostgreSQL database.
// It also sets up the database schema if it doesn't already exist.
func NewDatabase() (*Database, error) {
	config := configuration.GetConfig()

	db, err := sql.Open("postgres", config.DB_DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &Database{Connection: db}

	slog.Info("Database connection established", "operation", "database")
	return database, nil
}

// Close terminates the database connection and releases associated resources.
func (db *Database) Close() {
	if db.Connection != nil {
		if err := db.Connection.Close(); err != nil {
			slog.Error("Failed to close database connection", "exception", err, "operation", "database")
		} else {
			slog.Info("Database connection closed successfully", "operation", "database")
		}
	}
}
