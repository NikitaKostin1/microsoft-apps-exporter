package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"microsoft-apps-exporter/internal/models"
)

// extractKeys retrieves all keys from the given map as a slice.
// The order is not persisted!
func extractKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// generatePlaceholders creates a slice of PostgreSQL-style placeholders (e.g., $1, $2, ...).
func generatePlaceholders(count int) []string {
	placeholders := make([]string, count)
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	return placeholders
}

// mapFieldValues extracts the values of mapped fields based on a given mapping of column names.
func mapFieldValues(mappedFields models.ListItemMappedFields, columnsMap map[string]string, fieldsColumns []string) []interface{} {
	values := make([]interface{}, len(fieldsColumns))
	for i, dbColumn := range fieldsColumns {
		values[i] = mappedFields[columnsMap[dbColumn]]
	}
	return values
}

// buildUpdateClauses constructs SQL update statements dynamically for metadata and mapped fields.
func buildUpdateClauses(metadataColumns []string, columnsMap map[string]string, listItem models.ListItem) ([]string, []interface{}) {
	setClauses := make([]string, 0, len(metadataColumns)+len(columnsMap))
	values := make([]interface{}, 0, len(metadataColumns)+len(columnsMap))

	// Process metadata fields (excluding the primary key at index 0)
	for i, col := range metadataColumns[1:] {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, i+2))
		values = append(values, listItem.Metadata.AsArray()[i+1])
	}

	offset := len(metadataColumns)
	// Process mapped fields
	for dbColumn, apiColumn := range columnsMap {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", dbColumn, offset+1))
		values = append(values, listItem.MappedFields[apiColumn])
		offset++
	}

	return setClauses, values
}

// scanListItem extracts a ListItem from a database query result.
func scanListItem(rows *sql.Rows) (models.ListItem, error) {
	var listItem models.ListItem
	columns, err := rows.Columns()
	if err != nil {
		return models.ListItem{}, err
	}

	// Retrieve metadata field names for identification
	metadataFields := models.ListItemMetadata{}.DbColumns()
	values := make([]interface{}, len(columns))
	pointers := make([]interface{}, len(columns))
	for i := range values {
		pointers[i] = &values[i]
	}

	// Scan the row values into the pointers
	if err := rows.Scan(pointers...); err != nil {
		return models.ListItem{}, err
	}

	// Separate metadata and mapped fields
	metadataMap, fieldsMap := segregateColumns(columns, values, metadataFields)

	// Unmarshal JSON data into the respective struct fields
	json.Unmarshal(marshalJSON(metadataMap), &listItem.Metadata)
	json.Unmarshal(marshalJSON(fieldsMap), &listItem.MappedFields)

	return listItem, nil
}

// segregateColumns classifies database columns into metadata and mapped fields.
func segregateColumns(columns []string, values []interface{}, metadataFields []string) (map[string]interface{}, map[string]interface{}) {
	metadataMap := make(map[string]interface{})
	fieldsMap := make(map[string]interface{})

	for i, column := range columns {
		if contains(metadataFields, column) {
			metadataMap[column] = values[i]
		} else {
			fieldsMap[column] = values[i]
		}
	}
	return metadataMap, fieldsMap
}

// marshalJSON safely converts an interface to a JSON byte slice.
func marshalJSON(data interface{}) []byte {
	jsonData, _ := json.Marshal(data)
	return jsonData
}

// contains checks if a slice contains a given string.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
