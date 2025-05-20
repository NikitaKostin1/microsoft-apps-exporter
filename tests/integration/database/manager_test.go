//go:build testing && integration

package database_test

import (
	"database/sql"
	"errors"
	"log/slog"
	"math"
	"microsoft-apps-exporter/internal/database"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTransactionCommit ensures that a transaction successfully commits
// and data persists in the database.
func TestTransactionCommit(t *testing.T) {
	slog.SetLogLoggerLevel(math.MaxInt) // Disable logging

	db, err := database.NewDatabase()
	assert.NoError(t, err, "Database should initialize correctly")
	require.NotNil(t, db, "Database failed to initialize")
	defer db.Close()

	// Begin transaction and commit
	err = db.WithTransaction(func(tx *sql.Tx) error {
		_, execErr := tx.Exec("CREATE TABLE IF NOT EXISTS commit_test_table (id SERIAL PRIMARY KEY, value TEXT)")
		if execErr != nil {
			return execErr
		}
		_, execErr = tx.Exec("INSERT INTO commit_test_table (value) VALUES ('commit_test')")
		return execErr
	})
	assert.NoError(t, err, "Transaction should commit successfully")

	// Verify that the data exists after the transaction commits
	var count int
	err = db.Connection.QueryRow("SELECT COUNT(*) FROM commit_test_table WHERE value = 'commit_test'").Scan(&count)
	assert.NoError(t, err, "Query execution should not fail")
	assert.Equal(t, 1, count, "Transaction should have committed and inserted data")

	// Cleanup test data
	_, _ = db.Connection.Exec("DELETE FROM commit_test_table")
}

// TestTransactionRollback ensures that when an error occurs within a transaction,
// all changes are rolled back and no data persists in the database.
func TestTransactionRollback(t *testing.T) {
	slog.SetLogLoggerLevel(math.MaxInt) // Disable logging

	db, err := database.NewDatabase()
	require.NotNil(t, db, "Database failed to initialize")
	assert.NoError(t, err, "Database should initialize correctly")
	defer db.Close()

	_, err = db.Connection.Exec("CREATE TABLE IF NOT EXISTS rollback_test_table (id SERIAL PRIMARY KEY, value TEXT)")
	assert.NoError(t, err, "Table creation should not fail")

	// Start transaction and attempt rollback
	err = db.WithTransaction(func(tx *sql.Tx) error {
		_, execErr := tx.Exec("INSERT INTO rollback_test_table (value) VALUES ('rollback_test')")
		if execErr != nil {
			return execErr
		}
		return errors.New("forced rollback") // Simulates error - forces a rollback
	})
	assert.Error(t, err, "Transaction should fail and rollback")

	// Verify rollback by checking if the inserted data is gone
	var count int
	err = db.Connection.QueryRow("SELECT COUNT(*) FROM rollback_test_table WHERE value = 'rollback_test'").Scan(&count)
	assert.NoError(t, err, "Query execution should not fail")
	assert.Equal(t, 0, count, "Transaction should have rolled back and not inserted data")

	// Cleanup test data
	_, _ = db.Connection.Exec("DELETE FROM rollback_test_table")
}

// TestRollbackFailure simulates a failure when rolling back an already closed transaction.
// It ensures that the rollback failure is handled gracefully.
func TestRollbackFailure(t *testing.T) {
	slog.SetLogLoggerLevel(math.MaxInt) // Disable logging

	db, err := database.NewDatabase()
	require.NotNil(t, db, "Database failed to initialize")
	assert.NoError(t, err, "Database should initialize correctly")
	defer db.Close()

	err = db.WithTransaction(func(tx *sql.Tx) error {
		tx.Rollback()                            // First rollback happens here
		return errors.New("some rollback error") // Forces rollback again
	})
	assert.Error(t, err, "Transaction should fail due to rollback attempt on an already closed transaction")
}

// TestTransactionPanicHandling ensures that if a panic occurs within a transaction,
// the transaction is rolled back before the panic is rethrown.
func TestTransactionPanicHandling(t *testing.T) {
	slog.SetLogLoggerLevel(math.MaxInt) // Disable logging

	db, err := database.NewDatabase()
	require.NotNil(t, db, "Database failed to initialize")
	assert.NoError(t, err, "Database should initialize correctly")
	defer db.Close()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected a panic, but none occurred")
		}
	}()

	err = db.WithTransaction(func(tx *sql.Tx) error {
		_, execErr := tx.Exec("CREATE TABLE IF NOT EXISTS panic_test_table (id SERIAL PRIMARY KEY, value TEXT)")
		if execErr != nil {
			return execErr
		}
		_, execErr = tx.Exec("INSERT INTO panic_test_table (value) VALUES ('panic_test')")
		if execErr != nil {
			return execErr
		}
		panic("unexpected panic") // Simulate a panic
	})
	assert.Error(t, err, "Transaction should fail due to panic")

	// Verify rollback by checking if the inserted data is gone
	var count int
	err = db.Connection.QueryRow("SELECT COUNT(*) FROM panic_test_table WHERE value = 'panic_test'").Scan(&count)
	assert.NoError(t, err, "Query execution should not fail")
	assert.Equal(t, 0, count, "Transaction should have rolled back due to panic")

	// Cleanup test data
	_, _ = db.Connection.Exec("DELETE FROM panic_test_table")
}
