//go:build testing && unit

package api_test

import (
	"microsoft-apps-exporter/internal/api"
	"testing"

	"github.com/stretchr/testify/assert"
)

// NewListItemsWithDeltaOptions creates a delta query configuration for list items with optional field expansion and result limiting.
func TestNewListItemsWithDeltaOptions(t *testing.T) {
	tests := []struct {
		name           string
		expandFields   []string
		top            *int32
		expectedExpand string
		expectedTop    *int32
	}{
		{
			name:           "No expand fields and nil top",
			expandFields:   []string{},
			top:            nil,
			expectedExpand: "fields",
			expectedTop:    nil,
		},
		{
			name:           "Single expand field and nil top",
			expandFields:   []string{"Title"},
			top:            nil,
			expectedExpand: "fields($select=Title)",
			expectedTop:    nil,
		},
		{
			name:           "Multiple expand fields and nil top",
			expandFields:   []string{"Title", "Modified"},
			top:            nil,
			expectedExpand: "fields($select=Title,Modified)",
			expectedTop:    nil,
		},
		{
			name:           "No expand fields with top value",
			expandFields:   []string{},
			top:            int32ToPtr(10),
			expectedExpand: "fields",
			expectedTop:    int32ToPtr(10),
		},
		{
			name:           "Multiple expand fields with top value",
			expandFields:   []string{"Title", "Modified"},
			top:            int32ToPtr(5),
			expectedExpand: "fields($select=Title,Modified)",
			expectedTop:    int32ToPtr(5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := api.NewListItemsWithDeltaOptions(tt.expandFields, tt.top)

			// Then: The configuration should be properly constructed
			assert.NotNil(t, options)
			assert.NotNil(t, options.QueryParameters)

			// Verify expand parameter
			assert.Equal(t, []string{tt.expectedExpand}, options.QueryParameters.Expand)

			// Verify top parameter
			assert.Equal(t, tt.expectedTop, options.QueryParameters.Top)
		})
	}
}

// Helper function to create int32 pointers
func int32ToPtr(i int32) *int32 {
	return &i
}
