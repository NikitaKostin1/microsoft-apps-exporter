//go:build testing && e2e

package api_test

import (
	"context"
	"microsoft-apps-exporter/internal/api"
	"microsoft-apps-exporter/internal/configuration"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetList(t *testing.T) {
	setupProdResourcesYaml()

	graphHelper, err := api.NewGraphHelper(context.Background())
	require.NoError(t, err, "Failed to initialize GraphHelper")

	config := configuration.GetConfig()
	if config.Sharepoint == nil {
		t.Fatalf("Sharepoint resource is expected to be specified")
	}
	if len(config.Sharepoint.Lists) == 0 {
		t.Fatalf("No SharePoint lists configured")
	}

	list := config.Sharepoint.Lists[0]
	siteID, listID := list.SiteID, list.ListID

	metadata, err := graphHelper.GetList(siteID, listID)

	assert.NoError(t, err, "Fetching list metadata should succeed")
	assert.Len(t, metadata, 1, "Expected exactly one list metadata object")

	md := metadata[0]
	assert.Equal(t, siteID, md.SiteID, "SiteID should match input")
	assert.Equal(t, listID, md.ID, "ListID should match input")
	assert.NotEmpty(t, md.ETag, "ETag should not be empty")
	assert.NotEmpty(t, md.Name, "List name should not be empty")
	assert.NotEmpty(t, md.DisplayName, "Display name should not be empty")
}

func TestGetListItemsWithDelta(t *testing.T) {
	setupProdResourcesYaml()

	graphHelper, err := api.NewGraphHelper(context.Background())
	require.NoError(t, err, "Failed to initialize GraphHelper")

	config := configuration.GetConfig()
	if config.Sharepoint == nil {
		t.Fatalf("Sharepoint resource is expected to be specified")
	}
	if len(config.Sharepoint.Lists) == 0 {
		t.Fatalf("No SharePoint lists configured")
	}

	list := config.Sharepoint.Lists[0]
	siteID, listID := list.SiteID, list.ListID

	var top int32 = 10
	options := api.NewListItemsWithDeltaOptions(nil, &top)

	// First call to retrieve data and initial delta
	deltaLink, items, err := graphHelper.GetListItemsWithDelta(siteID, listID, nil, options)
	assert.NoError(t, err, "Fetching list items with delta should succeed")
	assert.Nil(t, deltaLink, "Delta link should be nil")
	assert.NotNil(t, items, "List items should not be nil")
	assert.Greater(t, len(*items), 0, "Expected at least one list item")

	item := (*items)[0]
	assert.NotEmpty(t, item.Metadata.ID, "Item ID should not be empty")
	assert.Equal(t, listID, item.Metadata.ListID, "ListID in metadata should match")
	assert.Equal(t, siteID, item.Metadata.SiteID, "SiteID in metadata should match")
	assert.NotEmpty(t, item.Metadata.ETag, "ETag should not be empty")
	assert.Greater(t, len(item.MappedFields), 0, "Item fields should not be empty")
	/*
		// Second call with delta (usually returns 0 unless items changed recently)
		newDeltaLink, newItems, err := graphHelper.GetListItemsWithDelta(siteID, listID, deltaLink, nil)
		assert.NoError(t, err, "Fetching with delta link should not fail")
		assert.NotNil(t, newDeltaLink, "New delta link should be returned")
		assert.NotNil(t, newItems, "Items should not be nil")
	*/
}
