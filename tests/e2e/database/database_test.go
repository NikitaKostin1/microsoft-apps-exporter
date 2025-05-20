//go:build testing && e2e

package database_test

/*
import (
	"fmt"
	"log/slog"
	"math"
	"testing"

	"microsoft-apps-exporter/internal/configuration"
	"microsoft-apps-exporter/internal/database"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupEnv() {
	godotenv.Load("../../../.env")
}

func setupProdResourcesYaml() {
	viper.Reset()
	viper.AddConfigPath("../../..")
}

func TestDatabaseConnectionAndPing(t *testing.T) {
	slog.SetLogLoggerLevel(math.MaxInt) // Disable logging
	setupEnv()
	setupProdResourcesYaml()

	db, err := database.NewDatabase()
	require.NoError(t, err, "Database connection should succeed")
	require.NotNil(t, db.Connection, "DB connection object must not be nil")

	// Sanity ping
	err = db.Connection.Ping()
	require.NoError(t, err, "Ping to database should succeed")

	db.Close()
}

func TestExpectedSchemaPresent(t *testing.T) {
	slog.SetLogLoggerLevel(math.MaxInt) // Disable logging
	setupEnv()

	config := configuration.GetConfig()
	if config.Sharepoint == nil || len(config.Sharepoint.Lists) == 0 {
		t.Fatal("Sharepoint resource with at least one list must be configured")
	}

	db, err := database.NewDatabase()
	require.NoError(t, err)
	defer db.Close()

	// Test sharepoint resource table if specified
	if config.Sharepoint.DbTableName != "" {
		t.Run("Sharepoint resource table schema", func(t *testing.T) {
			testTableSchema(t, db, config.Sharepoint.DbTableName)
		})
	}

	// Test each list's table schema
	for i, list := range config.Sharepoint.Lists {
		t.Run(fmt.Sprintf("List %d table schema", i+1), func(t *testing.T) {
			if list.DbTableName == "" {
				t.Fatal("List must specify database_table")
			}
			testTableSchema(t, db, list.DbTableName)
		})
	}
}

// testTableSchema checks if a table exists and contains required columns
func testTableSchema(t *testing.T, db *database.Database, tableName string) {
	// Verify table exists
	var tableExists bool
	err := db.Connection.QueryRow(`
        SELECT EXISTS (
            SELECT FROM information_schema.tables
            WHERE table_name = $1
        )`, tableName).Scan(&tableExists)
	require.NoError(t, err)
	require.True(t, tableExists, "Table %s should exist", tableName)

	// Verify required columns exist
	requiredColumns := []string{"etag", "id"} // Add other required columns as needed
	rows, err := db.Connection.Query(`
        SELECT column_name
        FROM information_schema.columns
        WHERE table_name = $1
    `, tableName)
	require.NoError(t, err)
	defer rows.Close()

	foundColumns := make(map[string]bool)
	for rows.Next() {
		var column string
		require.NoError(t, rows.Scan(&column))
		foundColumns[column] = true
	}
	require.NoError(t, rows.Err())

	for _, col := range requiredColumns {
		assert.True(t, foundColumns[col], "Table %s missing required column: %s", tableName, col)
	}
}
*/
