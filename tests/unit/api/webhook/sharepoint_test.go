//go:build testing && unit

package webhook_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"math"
	"microsoft-apps-exporter/internal/api/webhook"
	"microsoft-apps-exporter/internal/configuration"
	"microsoft-apps-exporter/internal/sync"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// setupTestResourcesYaml configures Viper to load test YAML resources.
func setupTestResourcesYaml() {
	configuration.ResetConfig()
	viper.Reset()

	viper.SetConfigName("resources")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../..")
}

// TestNewSharepointHandler tests the Sharepoint webhook handler for various scenarios.
func TestNewSharepointHandler(t *testing.T) {
	slog.SetLogLoggerLevel(math.MaxInt) // Disable logging
	setupTestResourcesYaml()

	syncer := &sync.Syncer{}
	handler := webhook.NewSharepointHandler(syncer)

	tests := []struct {
		name           string
		method         string
		urlQuery       string
		body           []byte
		expectedStatus int
	}{
		{
			name:           "Invalid method",
			method:         http.MethodGet,
			body:           nil,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Token validation",
			method:         http.MethodPost,
			urlQuery:       "?validationToken=test-token",
			body:           nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid request body",
			method:         http.MethodPost,
			body:           []byte(`{"invalid": "data"}`),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid resource format",
			method:         http.MethodPost,
			body:           []byte(`{"value": [{"resource": "invalid/resource"}]}`),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method,
				strings.Join([]string{"/webhook/sharepoint", tt.urlQuery}, ""),
				bytes.NewReader(tt.body))

			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

// TestExctractListReference tests the extraction of ListReference based on site and list IDs.
func TestExctractListReference(t *testing.T) {
	slog.SetLogLoggerLevel(math.MaxInt) // Disable logging
	setupTestResourcesYaml()

	tests := []struct {
		name        string
		siteID      string
		listID      string
		expectFound bool
	}{
		{"Valid list - Found", "site_id1", "list_id1", true},
		{"Valid list - Found", "site_id2", "list_id2", true},
		{"Invalid list - Not Found", "site999", "listX", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, found := webhook.ExctractListReference(tt.siteID, tt.listID)
			assert.Equal(t, tt.expectFound, found)
			if tt.expectFound {
				assert.Equal(t, tt.siteID, list.SiteID)
				assert.Equal(t, tt.listID, list.ListID)
			}
		})
	}
}

// TestExtractResourceUpdateData tests the extraction of resource update data from the request body.
func TestExtractResourceUpdateData(t *testing.T) {
	validBody := webhook.ResourceUpdateBody{
		Value: []struct {
			Resource     string `json:"resource"`
			ResourceData struct {
				OdataType string `json:"@odata.type"`
			} `json:"resourceData"`
		}{
			{
				Resource: "sites/site123/lists/listA",
				ResourceData: struct {
					OdataType string `json:"@odata.type"`
				}{
					OdataType: "#Microsoft.Graph.ListItem",
				},
			},
		},
	}

	invalidBodyType := webhook.ResourceUpdateBody{
		Value: []struct {
			Resource     string `json:"resource"`
			ResourceData struct {
				OdataType string `json:"@odata.type"`
			} `json:"resourceData"`
		}{
			{
				Resource: "sites/site123/lists/listA",
				ResourceData: struct {
					OdataType string `json:"@odata.type"`
				}{
					OdataType: "invalidType",
				},
			},
		},
	}

	invalidResource := webhook.ResourceUpdateBody{
		Value: []struct {
			Resource     string `json:"resource"`
			ResourceData struct {
				OdataType string `json:"@odata.type"`
			} `json:"resourceData"`
		}{
			{
				Resource: "invalid/resource/format",
				ResourceData: struct {
					OdataType string `json:"@odata.type"`
				}{
					OdataType: "#Microsoft.Graph.ListItem",
				},
			},
		},
	}

	tests := []struct {
		name        string
		body        webhook.ResourceUpdateBody
		expectError bool
	}{
		{"Valid Resource", validBody, false},
		{"Invalid Data Type", invalidBodyType, true},
		{"Invalid Resource Format", invalidResource, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyBytes))

			siteID, listID, err := webhook.ExtractResourceUpdateData(req)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "site123", siteID)
				assert.Equal(t, "listA", listID)
			}
		})
	}
}

// TestParseSharepointResource tests the parsing of Sharepoint resource strings.
func TestParseSharepointResource(t *testing.T) {
	tests := []struct {
		name        string
		resource    string
		expectSite  string
		expectList  string
		expectError bool
	}{
		{"Valid Resource", "sites/site123/lists/listA", "site123", "listA", false},
		{"Invalid Format", "invalid/resource/format", "", "", true},
		{"Missing Parts", "sites/site123/list/listA", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			siteID, listID, err := webhook.ParseSharepointResource(tt.resource)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectSite, siteID)
				assert.Equal(t, tt.expectList, listID)
			}
		})
	}
}
