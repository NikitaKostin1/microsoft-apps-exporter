package api

import (
	"encoding/json"
	"fmt"
	"microsoft-apps-exporter/internal/models"

	gmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	graphsites "github.com/microsoftgraph/msgraph-sdk-go/sites"
)

// requestList retrieves a SharePoint list using its site and list IDs.
func (g *GraphHelper) requestList(siteID, listID string) (gmodels.Listable, error) {
	return g.Client.Sites().BySiteId(siteID).Lists().ByListId(listID).Get(g.Ctx, nil)
}

// requestListItemsWithDelta retrieves paginated list items, updating delta links as needed.
func (g *GraphHelper) requestListItemsWithDelta(
	siteID, listID string, deltaLink *string,
	options *graphsites.ItemListsItemItemsDeltaRequestBuilderGetRequestConfiguration,
) (*string, []gmodels.ListItemable, error) {
	var (
		collectionResponse graphsites.ItemListsItemItemsDeltaGetResponseable
		err                error
	)

	// Use delta link if available
	req := g.Client.Sites().BySiteId(siteID).Lists().ByListId(listID).Items()
	if deltaLink != nil {
		collectionResponse, err = req.WithUrl(*deltaLink).Delta().GetAsDeltaGetResponse(g.Ctx, options)
	} else {
		collectionResponse, err = req.Delta().GetAsDeltaGetResponse(g.Ctx, options)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch list items: %w", err)
	}

	listItems := collectionResponse.GetValue()

	// Extract optional Top parameter safely
	var topLimit int
	if options != nil && options.QueryParameters.Top != nil {
		topLimit = int(*options.QueryParameters.Top)
	}

	// Handle pagination and delta link updates
	for {
		if delta := collectionResponse.GetOdataDeltaLink(); delta != nil {
			return delta, listItems, nil
		}

		nextLink := collectionResponse.GetOdataNextLink()
		if nextLink == nil || (topLimit > 0 && len(listItems) >= topLimit) {
			break
		}

		collectionResponse, err = req.WithUrl(*nextLink).Delta().GetAsDeltaGetResponse(g.Ctx, options)
		if err != nil {
			return nil, nil, fmt.Errorf("error fetching next page: %w", err)
		}

		listItems = append(listItems, collectionResponse.GetValue()...)
	}

	return nil, listItems, nil
}

// parseListItemResponse extracts and deserializes list item fields into a structured format.
func (g *GraphHelper) parseListItemResponse(itemResponse gmodels.ListItemable) (*models.ListItem, error) {
	fields := itemResponse.GetFields()

	// Serialize fields to JSON
	serializedFields, err := g.serializeFields(fields)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize fields: %w", err)
	}

	// Parse JSON into structured format
	mappedFields, err := deserializeFields(serializedFields)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize fields: %w", err)
	}

	return &models.ListItem{MappedFields: mappedFields}, nil
}

// serializeFields serializes SharePoint list item fields to JSON.
func (g *GraphHelper) serializeFields(fields gmodels.FieldValueSetable) ([]byte, error) {
	writer, err := g.Adapter.GetSerializationWriterFactory().GetSerializationWriter("application/json")
	if err != nil {
		return nil, fmt.Errorf("failed to create serialization writer: %w", err)
	}

	if err := writer.WriteObjectValue("", fields); err != nil {
		return nil, fmt.Errorf("failed to serialize object value: %w", err)
	}

	return writer.GetSerializedContent()
}

// deserializeFields converts serialized JSON data into a map of list item fields.
func deserializeFields(serializedFields []byte) (models.ListItemMappedFields, error) {
	var mappedData models.ListItemMappedFields
	if err := json.Unmarshal(serializedFields, &mappedData); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	delete(mappedData, "@odata.etag") // Remove irrelevant field

	return mappedData, nil
}

func safeString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
