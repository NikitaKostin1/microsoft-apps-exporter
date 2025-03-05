//go:build testing

// Exports internal functions for testing purposes.
// This file is only included in builds with the "testing" tag.
package database

import (
	"microsoft-apps-exporter/internal/models"
)

func ExtractKeys(m map[string]string) []string {
	return extractKeys(m)
}

func GeneratePlaceholders(count int) []string {
	return generatePlaceholders(count)
}

func MapFieldValues(mappedFields map[string]interface{}, columnsMap map[string]string, fieldsColumns []string) []interface{} {
	return mapFieldValues(mappedFields, columnsMap, fieldsColumns)
}

func BuildUpdateClauses(metadataColumns []string, columnsMap map[string]string, listItem models.ListItem) ([]string, []interface{}) {
	return buildUpdateClauses(metadataColumns, columnsMap, listItem)
}

func SegregateColumns(columns []string, values []interface{}, metadataFields []string) (map[string]interface{}, map[string]interface{}) {
	return segregateColumns(columns, values, metadataFields)
}

func MarshalJSON(data interface{}) []byte {
	return marshalJSON(data)
}

func Contains(slice []string, item string) bool {
	return contains(slice, item)
}
