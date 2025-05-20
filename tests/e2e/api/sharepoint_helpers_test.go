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

// TestRequestList_Success checks if a valid SharePoint list can be fetched
func TestRequestList_Success(t *testing.T) {
	setupProdResourcesYaml()

	graphHelper, err := api.NewGraphHelper(context.Background())
	require.NoError(t, err, "Failed to initialize GraphHelper")

	config := configuration.GetConfig()
	if config.Sharepoint == nil {
		t.Fatalf("Sharepoint resource is expected to be specified")
	}
	if len(config.Sharepoint.Lists) == 0 {
		t.Fatalf("No SharePoint lists configured. Expected at least one.")
	}

	list := config.Sharepoint.Lists[0]
	siteID, listID := list.SiteID, list.ListID

	listable, err := graphHelper.RequestList(siteID, listID)

	assert.NoError(t, err, "Fetching SharePoint list should succeed")
	assert.NotNil(t, listable, "List should not be nil")
}

// TestRequestListItemsWithDelta checks list item retrieval with options but without delta link (actual request would take minutes to accomplish)
func TestRequestListItemsWithDelta(t *testing.T) {
	setupProdResourcesYaml()

	graphHelper, err := api.NewGraphHelper(context.Background())
	require.NoError(t, err, "Failed to initialize GraphHelper")

	config := configuration.GetConfig()
	if config.Sharepoint == nil {
		t.Fatalf("Sharepoint resource is expected to be specified")
	}
	if len(config.Sharepoint.Lists) == 0 {
		t.Fatalf("No SharePoint lists configured. Expected at least one.")
	}

	list := config.Sharepoint.Lists[0]
	siteID, listID := list.SiteID, list.ListID

	var top int32 = 3000
	options := api.NewListItemsWithDeltaOptions(nil, &top)
	_, items, err := graphHelper.RequestListItemsWithDelta(siteID, listID, nil, options)

	assert.NoError(t, err, "Fetching list items should succeed")
	assert.NotNil(t, items, "Items should not be nil")
	assert.Equal(t, len(items), int(top), "Should return requested amount in options (if the resource actualy has this amount)")
	// assert.NotNil(t, deltaLink, "Delta link should not be nil")
}

/* // TestRequestListItems_WithDelta checks incremental updates using a delta link
func TestRequestListItems_WithDelta(t *testing.T) {
	setupProdResourcesYaml()

	graphHelper, err := api.NewGraphHelper(context.Background())
	require.NoError(t, err, "Failed to initialize GraphHelper")

	config := configuration.GetConfig()
	if len(config.Sharepoint.Lists) == 0 {
		t.Fatalf("No SharePoint lists configured. Expected at least one.")
	}

	list := config.Sharepoint.Lists[0]
	siteID, listID := list.SiteID, list.ListID

	var top int32 = 10
	options := api.NewListItemsWithDeltaOptions(nil, &top)

	// First request to get a valid delta link
	deltaLink, _, err := graphHelper.RequestListItemsWithDelta(siteID, listID, nil, options)
	assert.NoError(t, err, "Initial request should succeed")
	assert.NotNil(t, deltaLink, "Initial delta link should not be nil")

	// Second request using delta link
	newDeltaLink, items, err := graphHelper.RequestListItemsWithDelta(siteID, listID, deltaLink, options)

	assert.NoError(t, err, "Fetching incremental updates should succeed")
	assert.NotNil(t, items, "Incremental update items should not be nil")
	assert.NotNil(t, newDeltaLink, "New delta link should be returned")
} */

// TestParseListItemResponse checks if a SharePoint item can be parsed correctly
func TestParseListItemResponse(t *testing.T) {
	setupProdResourcesYaml()

	graphHelper, err := api.NewGraphHelper(context.Background())
	require.NoError(t, err, "Failed to initialize GraphHelper")

	config := configuration.GetConfig()
	if config.Sharepoint == nil {
		t.Fatalf("Sharepoint resource is expected to be specified")
	}
	if len(config.Sharepoint.Lists) == 0 {
		t.Fatalf("No SharePoint lists configured. Expected at least one.")
	}

	list := config.Sharepoint.Lists[0]
	siteID, listID := list.SiteID, list.ListID

	var top int32 = 10
	options := api.NewListItemsWithDeltaOptions(nil, &top)

	_, items, err := graphHelper.RequestListItemsWithDelta(siteID, listID, nil, options)
	assert.NoError(t, err, "Fetching list items should succeed")
	assert.Greater(t, len(items), 0, "Should return at least one item")
	assert.Equal(t, len(items), int(top), "Should return requested amount in options (if the resource actualy has this amount)")

	parsedItem, err := graphHelper.ParseListItemResponse(items[0])

	assert.NoError(t, err, "Parsing list item should succeed")
	assert.NotNil(t, parsedItem, "Parsed item should not be nil")
	assert.Greater(t, len(parsedItem.MappedFields), 0, "Parsed item should have fields")
}
