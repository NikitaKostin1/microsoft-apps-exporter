//go:build testing && integration

package database_test

import (
	"database/sql"
	"log/slog"
	"math"
	"microsoft-apps-exporter/internal/database"
	"microsoft-apps-exporter/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTempTable initializes a temporary table with hardcoded test data.
func createTempTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TEMP TABLE list_items (
			id TEXT PRIMARY KEY,
			list_id TEXT,
			site_id TEXT,
			etag TEXT,
			field1 TEXT,
			field2 INTEGER
		);

		-- Insert hardcoded test data
		INSERT INTO list_items (id, list_id, site_id, etag, field1, field2)
		VALUES (
			'123',
			'list-001',
			'site-001',
			'etag-001',
			'Test Item',
			42
		);
	`)
	return err
}

func TestScanListItem_Success(t *testing.T) {
	slog.SetLogLoggerLevel(math.MaxInt) // Disable logging

	db, err := database.NewDatabase()
	assert.NoError(t, err, "Database should initialize correctly")
	require.NotNil(t, db, "Database instance should not be nil")
	defer db.Close()

	err = db.WithTransaction(func(tx *sql.Tx) error {
		// Create temp table with hardcoded test data
		assert.NoError(t, createTempTable(tx), "Temporary table should be created")

		// Query and scan the first row
		rows, queryErr := tx.Query("SELECT * FROM list_items")
		assert.NoError(t, queryErr, "Query should execute successfully")
		defer rows.Close()

		// Ensure there is at least one row
		assert.True(t, rows.Next(), "There should be at least one row")

		// Scan the row into a ListItem
		listItem, scanErr := database.ScanListItem(rows)
		assert.NoError(t, scanErr, "Scan should not return an error")

		// Verify extracted data
		expected := models.ListItem{
			Metadata: models.ListItemMetadata{
				ID:     "123",
				ListID: "list-001",
				SiteID: "site-001",
				ETag:   "etag-001",
			},
			MappedFields: map[string]interface{}{
				"field1": "Test Item",
				"field2": float64(42), // json.Unmarshal function treat any number type as float64
			},
		}

		assert.Equal(t, expected, listItem, "Extracted ListItem should match expected")
		return nil
	})

	assert.NoError(t, err, "Transaction should complete successfully")
}
