package models

import (
	"fmt"
)

const SharepointResourceSignature string = "sites/%s/lists/%s"
const WebhookSharepointEndpoint string = "/webhook/sharepoint-notification"

func GenerateSharepointResourceString(siteID, listID string) string {
	return fmt.Sprintf(SharepointResourceSignature, siteID, listID)
}

type SharepointResource struct {
	DbTableName string          `mapstructure:"database_table"`
	Lists       []ListReference `mapstructure:"lists"`
}

type ListReference struct {
	SiteID      string            `mapstructure:"site_id"`
	ListID      string            `mapstructure:"list_id"`
	DbTableName string            `mapstructure:"database_table"`
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
	MappedFields ListItemMappedFields
}

type ListItemMetadata struct {
	ID     string `json:"id"`
	ListID string `json:"list_id"`
	SiteID string `json:"site_id"`
	ETag   string `json:"etag"`
}

type ListItemMappedFields map[string]any

func (m *ListItemMetadata) AsArray() []interface{} {
	return []interface{}{m.ID, m.ListID, m.SiteID, m.ETag}
}

func (m ListItemMetadata) DbColumns() []string {
	return []string{"id", "list_id", "site_id", "etag"}
}
