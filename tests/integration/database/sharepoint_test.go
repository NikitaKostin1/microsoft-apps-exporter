//go:build testing && integration

package database_test

import (
	"context"
	"database/sql"
	"microsoft-apps-exporter/internal/database"
	"microsoft-apps-exporter/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDatabase initializes the test database and creates the required tables.
func setupTestDatabase(t *testing.T) *database.Database {
	// Connect to the test database
	db, err := database.NewDatabase()
	require.NoError(t, err, "Failed to connect to the test database")
	require.NotNil(t, db, "Database instance should not be nil")

	// Create the sharepoint_lists table
	_, err = db.Connection.ExecContext(context.Background(), `
		CREATE TEMP TABLE sharepoint_lists (
			id TEXT PRIMARY KEY,
			site_id TEXT,
			etag TEXT,
			name TEXT,
			display_name TEXT,
			delta_link TEXT
		);
	`)
	require.NoError(t, err, "Failed to create sharepoint_lists table")

	// Create a dynamic table for list items
	_, err = db.Connection.ExecContext(context.Background(), `
		CREATE TEMP TABLE list_items (
			id TEXT PRIMARY KEY,
			list_id TEXT,
			site_id TEXT,
			etag TEXT,
			field1 TEXT,
			field2 INTEGER
		);
	`)
	require.NoError(t, err, "Failed to create list_items table")

	return db
}

// teardownTestDatabase cleans up the test database.
func teardownTestDatabase(db *database.Database) {
	db.Connection.Close()
}

/*
Lists
*/

// TestGetList tests the GetList function.
func TestGetList(t *testing.T) {
	db := setupTestDatabase(t)
	defer teardownTestDatabase(db)

	// Insert test data
	_, err := db.Connection.ExecContext(context.Background(), `
		INSERT INTO sharepoint_lists (id, site_id, etag, name, display_name, delta_link)
		VALUES ('list-001', 'site-001', 'etag-001', 'Test List', 'Test Display Name', 'delta-link-001');
	`)
	require.NoError(t, err, "Failed to insert test data")

	// Test GetList
	lists, err := db.GetList("list-001")
	assert.NoError(t, err, "GetList should not return an error")
	assert.Len(t, lists, 1, "Expected one list to be returned")

	expected := models.ListMetadata{
		ID:          "list-001",
		SiteID:      "site-001",
		ETag:        "etag-001",
		Name:        "Test List",
		DisplayName: "Test Display Name",
		DeltaLink:   stringPtr("delta-link-001"),
	}
	assert.Equal(t, expected, lists[0], "Returned list metadata should match expected")
}

// TestInsertLists tests the InsertLists function.
func TestInsertLists(t *testing.T) {
	db := setupTestDatabase(t)
	defer teardownTestDatabase(db)

	// Test data
	lists := []models.ListMetadata{
		{
			ID:          "list-001",
			SiteID:      "site-001",
			ETag:        "etag-001",
			Name:        "Test List 1",
			DisplayName: "Test Display Name 1",
			DeltaLink:   stringPtr("delta-link-001"),
		},
		{
			ID:          "list-002",
			SiteID:      "site-002",
			ETag:        "etag-002",
			Name:        "Test List 2",
			DisplayName: "Test Display Name 2",
			DeltaLink:   stringPtr("delta-link-002"),
		},
	}

	// Test InsertLists
	err := db.InsertLists(&lists)
	assert.NoError(t, err, "InsertLists should not return an error")

	// Verify data was inserted
	var count int
	err = db.Connection.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM sharepoint_lists;").Scan(&count)
	assert.NoError(t, err, "Failed to query row count")
	assert.Equal(t, 2, count, "Expected 2 rows to be inserted")
}

// TestUpdateListIgnoreDelta tests the UpdateListIgnoreDelta function.
func TestUpdateListIgnoreDelta(t *testing.T) {
	db := setupTestDatabase(t)
	defer teardownTestDatabase(db)

	// Insert test data
	_, err := db.Connection.ExecContext(context.Background(), `
		INSERT INTO sharepoint_lists (id, site_id, etag, name, display_name, delta_link)
		VALUES ('list-001', 'site-001', 'etag-001', 'Test List', 'Test Display Name', 'delta-link-001');
	`)
	require.NoError(t, err, "Failed to insert test data")

	// Test UpdateListIgnoreDelta
	updatedMetadata := models.ListMetadata{
		ID:          "list-001",
		SiteID:      "site-002",
		ETag:        "etag-002",
		Name:        "Updated List",
		DisplayName: "Updated Display Name",
		DeltaLink:   stringPtr("delta-link-001"),
	}
	err = db.UpdateListIgnoreDelta(updatedMetadata)
	assert.NoError(t, err, "UpdateListIgnoreDelta should not return an error")

	// Verify data was updated
	var metadata models.ListMetadata
	err = db.Connection.QueryRowContext(context.Background(), `
		SELECT id, site_id, etag, name, display_name, delta_link
		FROM sharepoint_lists
		WHERE id = 'list-001';
	`).Scan(
		&metadata.ID,
		&metadata.SiteID,
		&metadata.ETag,
		&metadata.Name,
		&metadata.DisplayName,
		&metadata.DeltaLink,
	)
	assert.NoError(t, err, "Failed to query updated data")
	assert.Equal(t, updatedMetadata, metadata, "Updated metadata should match expected")
}

// TestDeleteList tests the DeleteList function.
func TestDeleteList(t *testing.T) {
	db := setupTestDatabase(t)
	defer teardownTestDatabase(db)

	// Insert test data
	_, err := db.Connection.ExecContext(context.Background(), `
		INSERT INTO sharepoint_lists (id, site_id, etag, name, display_name, delta_link)
		VALUES ('list-001', 'site-001', 'etag-001', 'Test List', 'Test Display Name', 'delta-link-001');
	`)
	require.NoError(t, err, "Failed to insert test data")

	// Test DeleteList
	err = db.DeleteList("list-001")
	assert.NoError(t, err, "DeleteList should not return an error")

	// Verify data was deleted
	var count int
	err = db.Connection.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM sharepoint_lists;").Scan(&count)
	assert.NoError(t, err, "Failed to query row count")
	assert.Equal(t, 0, count, "Expected 0 rows after deletion")
}

// TestGetDeltaLink tests the GetDeltaLink function.
func TestGetDeltaLink(t *testing.T) {
	db := setupTestDatabase(t)
	defer teardownTestDatabase(db)

	// Insert test data
	_, err := db.Connection.ExecContext(context.Background(), `
		INSERT INTO sharepoint_lists (id, site_id, etag, name, display_name, delta_link)
		VALUES ('list-001', 'site-001', 'etag-001', 'Test List', 'Test Display Name', 'delta-link-001');
	`)
	require.NoError(t, err, "Failed to insert test data")

	// Test GetDeltaLink
	deltaLink, err := db.GetDeltaLink("list-001")
	assert.NoError(t, err, "GetDeltaLink should not return an error")
	assert.Equal(t, "delta-link-001", *deltaLink, "Delta link should match expected")
}

// TestSaveDeltaLink tests the SaveDeltaLink function.
func TestSaveDeltaLink(t *testing.T) {
	db := setupTestDatabase(t)
	defer teardownTestDatabase(db)

	// Insert test data
	_, err := db.Connection.ExecContext(context.Background(), `
		INSERT INTO sharepoint_lists (id, site_id, etag, name, display_name, delta_link)
		VALUES ('list-001', 'site-001', 'etag-001', 'Test List', 'Test Display Name', NULL);
	`)
	require.NoError(t, err, "Failed to insert test data")

	// Test SaveDeltaLink
	err = db.SaveDeltaLink("list-001", "delta-link-001")
	assert.NoError(t, err, "SaveDeltaLink should not return an error")

	// Verify delta link was saved
	var deltaLink sql.NullString
	err = db.Connection.QueryRowContext(context.Background(), `
		SELECT delta_link FROM sharepoint_lists WHERE id = 'list-001';
	`).Scan(&deltaLink)
	assert.NoError(t, err, "Failed to query delta link")
	assert.True(t, deltaLink.Valid, "Delta link should not be NULL")
	assert.Equal(t, "delta-link-001", deltaLink.String, "Delta link should match expected")
}

// TestDeleteDeltaLink tests the DeleteDeltaLink function.
func TestDeleteDeltaLink(t *testing.T) {
	db := setupTestDatabase(t)
	defer teardownTestDatabase(db)

	// Insert test data
	_, err := db.Connection.ExecContext(context.Background(), `
		INSERT INTO sharepoint_lists (id, site_id, etag, name, display_name, delta_link)
		VALUES ('list-001', 'site-001', 'etag-001', 'Test List', 'Test Display Name', 'delta-link-001');
	`)
	require.NoError(t, err, "Failed to insert test data")

	// Test DeleteDeltaLink
	err = db.DeleteDeltaLink("list-001")
	assert.NoError(t, err, "DeleteDeltaLink should not return an error")

	// Verify delta link was deleted
	var deltaLink sql.NullString
	err = db.Connection.QueryRowContext(context.Background(), `
		SELECT delta_link FROM sharepoint_lists WHERE id = 'list-001';
	`).Scan(&deltaLink)
	assert.NoError(t, err, "Failed to query delta link")
	assert.False(t, deltaLink.Valid, "Delta link should be NULL")
}

// Helper function to create a string pointer.
func stringPtr(s string) *string {
	return &s
}

/*
List Items
*/

// TestGetListItems tests the GetListItems function.
func TestGetListItems(t *testing.T) {
	db := setupTestDatabase(t)
	defer teardownTestDatabase(db)

	// Insert test data
	_, err := db.Connection.ExecContext(context.Background(), `
		INSERT INTO list_items (id, list_id, site_id, etag, field1, field2)
		VALUES 
			('item-001', 'list-001', 'site-001', 'etag-001', 'Test Item 1', 42),
			('item-002', 'list-001', 'site-001', 'etag-002', 'Test Item 2', 99);
	`)
	require.NoError(t, err, "Failed to insert test data")

	// Test GetListItems
	listItems, err := db.GetListItems("list_items", "site-001", "list-001")
	assert.NoError(t, err, "GetListItems should not return an error")
	assert.Len(t, *listItems, 2, "Expected 2 list items to be returned")

	// Verify the first list item
	expectedItem1 := models.ListItem{
		Metadata: models.ListItemMetadata{
			ID:     "item-001",
			ListID: "list-001",
			SiteID: "site-001",
			ETag:   "etag-001",
		},
		MappedFields: map[string]interface{}{
			"field1": "Test Item 1",
			"field2": float64(42), // JSON numbers are unmarshaled as float64
		},
	}
	assert.Equal(t, expectedItem1, (*listItems)[0], "First list item should match expected")

	// Verify the second list item
	expectedItem2 := models.ListItem{
		Metadata: models.ListItemMetadata{
			ID:     "item-002",
			ListID: "list-001",
			SiteID: "site-001",
			ETag:   "etag-002",
		},
		MappedFields: map[string]interface{}{
			"field1": "Test Item 2",
			"field2": float64(99), // JSON numbers are unmarshaled as float64
		},
	}
	assert.Equal(t, expectedItem2, (*listItems)[1], "Second list item should match expected")
}

// TestInsertListItems tests the InsertListItems function.
func TestInsertListItems(t *testing.T) {
	db := setupTestDatabase(t)
	defer teardownTestDatabase(db)

	// Test data
	listItems := []models.ListItem{
		{
			Metadata: models.ListItemMetadata{
				ID:     "item-001",
				ListID: "list-001",
				SiteID: "site-001",
				ETag:   "etag-001",
			},
			MappedFields: map[string]interface{}{
				"field1": "Test Item 1",
				"field2": 42,
			},
		},
		{
			Metadata: models.ListItemMetadata{
				ID:     "item-002",
				ListID: "list-001",
				SiteID: "site-001",
				ETag:   "etag-002",
			},
			MappedFields: map[string]interface{}{
				"field1": "Test Item 2",
				"field2": 99,
			},
		},
	}

	// Columns mapping
	columnsMap := map[string]string{
		"field1": "field1",
		"field2": "field2",
	}

	// Test InsertListItems
	err := db.InsertListItems("list_items", columnsMap, &listItems)
	assert.NoError(t, err, "InsertListItems should not return an error")

	// Verify data was inserted
	var count int
	err = db.Connection.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM list_items;").Scan(&count)
	assert.NoError(t, err, "Failed to query row count")
	assert.Equal(t, 2, count, "Expected 2 rows to be inserted")
}

// TestUpdateListItem tests the UpdateListItem function.
func TestUpdateListItem(t *testing.T) {
	db := setupTestDatabase(t)
	defer teardownTestDatabase(db)

	// Insert test data
	_, err := db.Connection.ExecContext(context.Background(), `
		INSERT INTO list_items (id, list_id, site_id, etag, field1, field2)
		VALUES ('item-001', 'list-001', 'site-001', 'etag-001', 'Test Item 1', 42);
	`)
	require.NoError(t, err, "Failed to insert test data")

	// Test UpdateListItem
	updatedItem := models.ListItem{
		Metadata: models.ListItemMetadata{
			ID:     "item-001",
			ListID: "list-001",
			SiteID: "site-001",
			ETag:   "etag-002", // Updated ETag
		},
		MappedFields: map[string]interface{}{
			"field1": "Updated Item 1", // Updated field1
			"field2": 100,              // Updated field2
		},
	}

	columnsMap := map[string]string{
		"field1": "field1",
		"field2": "field2",
	}

	err = db.UpdateListItem("list_items", columnsMap, updatedItem)
	assert.NoError(t, err, "UpdateListItem should not return an error")

	// Verify data was updated
	var field1 string
	var field2 int
	err = db.Connection.QueryRowContext(context.Background(), `
		SELECT field1, field2 FROM list_items WHERE id = 'item-001';
	`).Scan(&field1, &field2)

	assert.NoError(t, err, "Failed to query updated data")
	assert.Equal(t, "Updated Item 1", field1, "field1 should be updated")
	assert.Equal(t, 100, field2, "field2 should be updated")
}

// TestDeleteListItem tests the DeleteListItem function.
func TestDeleteListItem(t *testing.T) {
	db := setupTestDatabase(t)
	defer teardownTestDatabase(db)

	// Insert test data
	_, err := db.Connection.ExecContext(context.Background(), `
		INSERT INTO list_items (id, list_id, site_id, etag, field1, field2)
		VALUES ('item-001', 'list-001', 'site-001', 'etag-001', 'Test Item 1', 42);
	`)
	require.NoError(t, err, "Failed to insert test data")

	// Test DeleteListItem
	err = db.DeleteListItem("list_items", "item-001")
	assert.NoError(t, err, "DeleteListItem should not return an error")

	// Verify data was deleted
	var count int
	err = db.Connection.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM list_items;").Scan(&count)
	assert.NoError(t, err, "Failed to query row count")
	assert.Equal(t, 0, count, "Expected 0 rows after deletion")
}
