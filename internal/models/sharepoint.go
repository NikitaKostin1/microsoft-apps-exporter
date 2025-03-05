package models

import (
	"fmt"
)

const SharepointResourceSignature string = "sites/%s/lists/%s"
const WebhookSharepointEndpoint string = "/webhook/sharepoint-notification"

func GenerateSharepointResourceString(siteID, listID string) string {
	return fmt.Sprintf(SharepointResourceSignature, siteID, listID)
}

type List struct {
	Config   ListConfig
	Metadata ListMetadata
}

// SharepointList represents a SharePoint list with its associated Site and List IDs
type ListConfig struct {
	SiteID      string            `mapstructure:"site_id"`
	ListID      string            `mapstructure:"list_id"`
	DbTableName string            `mapstructure:"table_name"`
	ColumnsMap  map[string]string `mapstructure:"columns_map"`
}

type ListMetadata struct {
	ID          string  `json:"id"`
	SiteID      string  `json:"site_id"`
	ETag        string  `json:"etag"`
	Name        string  `json:"name"`
	DisplayName string  `json:"display_name"`
	DeltaLink   *string `json:"delta_link"`
}

type ListItem struct {
	Metadata     ListItemMetadata
	MappedFields map[string]interface{}
}

type ListItemMetadata struct {
	ID     string `json:"id"`
	ListID string `json:"list_id"`
	SiteID string `json:"site_id"`
	ETag   string `json:"etag"`
}

func (m *ListItemMetadata) AsArray() []interface{} {
	return []interface{}{m.ID, m.ListID, m.SiteID, m.ETag}
}

func (m ListItemMetadata) DbColumns() []string {
	return []string{"id", "list_id", "site_id", "etag"}
}
