//go:build testing

// Exports internal functions for testing purposes.
// This file is only included in builds with the "testing" tag.
package api

import (
	gmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
)

func (g *GraphHelper) SerializeFields(fields gmodels.FieldValueSetable) ([]byte, error) {
	return g.serializeFields(fields)
}

func DeserializeFields(serializedFields []byte) (map[string]interface{}, error) {
	return deserializeFields(serializedFields)
}
