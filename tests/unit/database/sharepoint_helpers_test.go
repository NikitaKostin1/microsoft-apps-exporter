//go:build testing && unit

package database_test

import (
	"microsoft-apps-exporter/internal/database"
	"microsoft-apps-exporter/internal/models"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestExtractKeys verifies that ExtractKeys correctly extracts map keys into a slice.
func TestExtractKeys(t *testing.T) {
	input := map[string]string{"a": "1", "b": "2", "c": "3"}
	expected := []string{"a", "b", "c"}

	result := database.ExtractKeys(input)

	assert.ElementsMatch(t, expected, result)
}

// TestGeneratePlaceholders ensures placeholders are correctly generated for SQL queries.
func TestGeneratePlaceholders(t *testing.T) {
	expected := []string{"$1", "$2", "$3"}

	result := database.GeneratePlaceholders(3)

	assert.ElementsMatch(t, expected, result)
}

// TestMapFieldValues checks if fields are correctly mapped based on the provided column mappings.
func TestMapFieldValues(t *testing.T) {
	mappedFields := models.ListItemMappedFields{"col1": 10, "col2": 20}
	columnsMap := map[string]string{"dbCol1": "col1", "dbCol2": "col2"}
	fieldsColumns := []string{"dbCol1", "dbCol2"}

	expected := []interface{}{10, 20}

	result := database.MapFieldValues(mappedFields, columnsMap, fieldsColumns)

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestBuildUpdateClauses verifies that the update SQL clauses and values are built correctly.
func TestBuildUpdateClauses(t *testing.T) {
	metadataColumns := []string{"id", "list_id", "site_id", "etag"}
	columnsMap := map[string]string{"db_col": "api_col"}
	listItem := models.ListItem{
		Metadata:     models.ListItemMetadata{ID: "1", ListID: "2", SiteID: "3", ETag: "etag_val"},
		MappedFields: models.ListItemMappedFields{"api_col": "value"},
	}

	expectedClauses := []string{"list_id = $2", "site_id = $3", "etag = $4", "db_col = $5"}
	expectedValues := []interface{}{"2", "3", "etag_val", "value"}

	clauses, values := database.BuildUpdateClauses(metadataColumns, columnsMap, listItem)

	if !reflect.DeepEqual(expectedClauses, clauses) {
		t.Errorf("Expected %v, got %v", expectedClauses, clauses)
	}
	if !reflect.DeepEqual(expectedValues, values) {
		t.Errorf("Expected %v, got %v", expectedValues, values)
	}
}

// TestSegregateColumns ensures metadata and custom fields are separated correctly.
func TestSegregateColumns(t *testing.T) {
	columns := []string{"id", "list_id", "custom"}
	values := []interface{}{"1", "2", "custom_value"}
	metadataFields := []string{"id", "list_id"}

	expectedMetadata := map[string]interface{}{"id": "1", "list_id": "2"}
	expectedFields := map[string]interface{}{"custom": "custom_value"}

	metadata, fields := database.SegregateColumns(columns, values, metadataFields)

	if !reflect.DeepEqual(expectedMetadata, metadata) {
		t.Errorf("Expected %v, got %v", expectedMetadata, metadata)
	}
	if !reflect.DeepEqual(expectedFields, fields) {
		t.Errorf("Expected %v, got %v", expectedFields, fields)
	}
}

// TestContains checks if an element exists within a slice.
func TestContains(t *testing.T) {
	slice := []string{"a", "b", "c"}

	if !database.Contains(slice, "b") {
		t.Errorf("Expected true, got false")
	}
	if database.Contains(slice, "d") {
		t.Errorf("Expected false, got true")
	}
}

// TestMarshalJSON verifies that JSON marshaling returns the expected JSON string.
func TestMarshalJSON(t *testing.T) {
	data := map[string]string{"key": "value"}
	expected := `{"key":"value"}`

	result := string(database.MarshalJSON(data))

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
