//go:build testing && e2e

package api_test

import (
	"context"
	"log/slog"
	"math"
	"testing"

	"microsoft-apps-exporter/internal/api"
	"microsoft-apps-exporter/internal/configuration"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupProdResourcesYaml() {
	viper.Reset()
	viper.AddConfigPath("../../..")
}

func TestNewGraphHelperInitialization(t *testing.T) {
	slog.SetLogLoggerLevel(math.MaxInt) // Disable logging

	setupProdResourcesYaml()

	tests := []struct {
		name           string
		setup          func()
		testApiRequest func(*api.GraphHelper) error
	}{
		{
			name: "Authenticate and fetch SharePoint resource",
			testApiRequest: func(g *api.GraphHelper) error {
				// Ensure SharePoint lists exist in configuration
				config := configuration.GetConfig()
				if config.Sharepoint == nil {
					t.Fatalf("Sharepoint resource is expected to be specified")
				}
				if len(config.Sharepoint.Lists) == 0 {
					t.Fatalf("No SharePoint lists configured. Expected at least one.")
				}

				list := config.Sharepoint.Lists[0]
				_, err := g.Client.Sites().BySiteId(list.SiteID).Lists().ByListId(list.ListID).Get(context.Background(), nil)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configuration.ResetConfig()
			config := configuration.GetConfig()

			// Ensure required credentials are present
			if config.GRAPH_TENANT_ID == "" || config.GRAPH_CLIENT_ID == "" || config.GRAPH_CLIENT_SECRET == "" {
				t.Fatalf("Missing required Microsoft Graph API credentials in configuration.")
			}

			// Initialize GraphHelper
			ctx := context.Background()
			graph, err := api.NewGraphHelper(ctx)
			require.NoError(t, err, "Failed to initialize GraphHelper")
			require.NotNil(t, graph, "GraphHelper should be instantiated")

			// Validate that authentication and initialization succeeded
			assert.NotNil(t, graph.Ctx, "GraphHelper context should be set")
			assert.NotNil(t, graph.Credential, "GraphHelper credential should be initialized")
			assert.NotNil(t, graph.Adapter, "GraphHelper adapter should be initialized")
			assert.NotNil(t, graph.Client, "GraphHelper client should be initialized")
			assert.NotEmpty(t, graph.AppScopes, "GraphHelper scopes should not be empty")

			// Execute API request test
			if tt.testApiRequest != nil {
				err := tt.testApiRequest(graph)
				if err != nil {
					t.Fatalf("API request failed: %v", err)
				}
			}
		})
	}
}
