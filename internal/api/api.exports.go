//go:build testing

// Exports internal functions for testing purposes.
// This file is only included in builds with the "testing" tag.
package api

import (
	"microsoft-apps-exporter/internal/models"

	gmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	graphsites "github.com/microsoftgraph/msgraph-sdk-go/sites"
)

var (
	SubscriptionUpdateExpiry    = subscriptionUpdateExpiry
	WebhookSubscriptionEndpoint = webhookSubscriptionEndpoint
)

func (g *GraphHelper) RequestList(siteID, listID string) (gmodels.Listable, error) {
	return g.requestList(siteID, listID)
}

func (g *GraphHelper) RequestListItemsWithDelta(siteID, listID string, deltaLink *string,
	options *graphsites.ItemListsItemItemsDeltaRequestBuilderGetRequestConfiguration) (*string, []gmodels.ListItemable, error) {
	return g.requestListItemsWithDelta(siteID, listID, deltaLink, options)
}

func (g *GraphHelper) ParseListItemResponse(itemResponse gmodels.ListItemable) (*models.ListItem, error) {
	return g.parseListItemResponse(itemResponse)
}

func DeserializeFields(serializedFields []byte) (models.ListItemMappedFields, error) {
	return deserializeFields(serializedFields)
}
