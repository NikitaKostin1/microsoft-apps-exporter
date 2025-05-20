//go:build testing && unit

package api_test

import (
	"testing"

	"microsoft-apps-exporter/internal/api"
	"microsoft-apps-exporter/internal/models"

	"github.com/stretchr/testify/assert"
)

// TestDeserializeFields tests the DeserializeFields function from the api package.
func TestDeserializeFields(t *testing.T) {
	validJSON := []byte(`{"key1": "value1", "@odata.etag": "ignore"}`)
	invalidJSON := []byte(`{"key1": "value1", "@odata.etag":}`)

	t.Run("Valid JSON", func(t *testing.T) {
		result, err := api.DeserializeFields(validJSON)
		assert.NoError(t, err)
		assert.Equal(t, models.ListItemMappedFields{"key1": "value1"}, result)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		_, err := api.DeserializeFields(invalidJSON)
		assert.Error(t, err)
	})
}
