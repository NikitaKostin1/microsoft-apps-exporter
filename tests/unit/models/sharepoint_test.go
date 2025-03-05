//go:build testing && unit

package models_test

import (
	"testing"

	"microsoft-apps-exporter/internal/models"

	"github.com/stretchr/testify/assert"
)

// TestGenerateSharepointResourceString verifies that the function correctly formats a SharePoint resource string.
func TestGenerateSharepointResourceString(t *testing.T) {
	siteID := "site123"
	listID := "list456"
	expected := "sites/site123/lists/list456"
	actual := models.GenerateSharepointResourceString(siteID, listID)

	assert.Equal(t, expected, actual, "Generated SharePoint resource string is incorrect")
}

// TestListItemMetadata_AsArray checks if ListItemMetadata is correctly converted to an array.
func TestListItemMetadata_AsArray(t *testing.T) {
	metadata := models.ListItemMetadata{
		ID:     "item123",
		ListID: "list456",
		SiteID: "site789",
		ETag:   "etag_001",
	}

	expected := []interface{}{"item123", "list456", "site789", "etag_001"}
	actual := metadata.AsArray()

	assert.Equal(t, expected, actual, "AsArray() did not return the expected result")
}

// TestListItemMetadata_DbColumns ensures that the correct database column names are returned.
func TestListItemMetadata_DbColumns(t *testing.T) {
	metadata := models.ListItemMetadata{}
	expected := []string{"id", "list_id", "site_id", "etag"}
	actual := metadata.DbColumns()

	assert.Equal(t, expected, actual, "DbColumns() did not return expected column names")
}
