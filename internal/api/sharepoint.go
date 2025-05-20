package api

import (
	"fmt"
	"microsoft-apps-exporter/internal/models"
	"strings"

	graphsites "github.com/microsoftgraph/msgraph-sdk-go/sites"
)

// GetListMetadata retrieves metadata of a SharePoint list by site ID and list ID.
func (g *GraphHelper) GetList(siteID, listID string) ([]models.ListMetadata, error) {
	listResponse, err := g.requestList(siteID, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch list metadata: %w", err)
	}

	return []models.ListMetadata{
		{
			ID:          listID,
			SiteID:      siteID,
			ETag:        *listResponse.GetETag(),
			Name:        *listResponse.GetName(),
			DisplayName: *listResponse.GetDisplayName(),
		},
	}, nil
}

// GetListItemsWithDelta retrieves SharePoint list items using Delta Query for tracking changes.
func (g *GraphHelper) GetListItemsWithDelta(
	siteID, listID string, deltaLink *string, options *graphsites.ItemListsItemItemsDeltaRequestBuilderGetRequestConfiguration,
) (*string, *[]models.ListItem, error) {
	newDeltaLink, itemsResponse, err := g.requestListItemsWithDelta(siteID, listID, deltaLink, options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve list items: %w", err)
	}

	listItems := make([]models.ListItem, 0, len(itemsResponse))
	for _, itemResponse := range itemsResponse {
		listItem, err := g.parseListItemResponse(itemResponse)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse list item fields: %w", err)
		}

		listItem.Metadata = models.ListItemMetadata{
			ID:     *itemResponse.GetId(),
			ListID: listID,
			SiteID: siteID,
			ETag:   safeString(itemResponse.GetETag()),
		}
		listItems = append(listItems, *listItem)
	}

	return newDeltaLink, &listItems, nil
}

// NewListItemsWithDeltaOptions generates request configuration for delta-tracked list item retrieval.
func NewListItemsWithDeltaOptions(expandFields []string, top *int32) *graphsites.ItemListsItemItemsDeltaRequestBuilderGetRequestConfiguration {
	expandString := "fields"
	if len(expandFields) > 0 {
		expandString += fmt.Sprintf("($select=%s)", strings.Join(expandFields, ","))
	}

	return &graphsites.ItemListsItemItemsDeltaRequestBuilderGetRequestConfiguration{
		QueryParameters: &graphsites.ItemListsItemItemsDeltaRequestBuilderGetQueryParameters{
			Expand: []string{expandString}, // Expand related entities
			Top:    top,                    // Show only the first n items
		},
	}
}
