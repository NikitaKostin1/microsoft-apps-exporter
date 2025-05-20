//go:build testing && integration

package database_test

import (
	"log/slog"
	"math"
	"microsoft-apps-exporter/internal/database"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDatabase_ConnectionSuccess(t *testing.T) {
	slog.SetLogLoggerLevel(math.MaxInt) // Disable logging

	db, err := database.NewDatabase()
	assert.NoError(t, err, "Expected no error when initializing the database")
	require.NotNil(t, db, "Database instance should not be nil")
	assert.NotNil(t, db.Connection, "Database connection should not be nil")

	err = db.Connection.Ping()
	assert.NoError(t, err, "Database ping should succeed with valid connection")

	db.Close()
}

func TestDatabase_Close(t *testing.T) {
	slog.SetLogLoggerLevel(math.MaxInt) // Disable logging

	db, err := database.NewDatabase()
	assert.NoError(t, err, "Expected no error when initializing the database")
	require.NotNil(t, db, "Database instance should not be nil")

	// Close database and ensure no panic
	assert.NotPanics(t, func() { db.Close() }, "Database close should not panic")
}
