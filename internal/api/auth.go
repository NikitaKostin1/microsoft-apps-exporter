package api

import (
	"fmt"
	"strings"

	"microsoft-apps-exporter/internal/configuration"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	auth "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

// InitializeGraphForUserAuth sets up authentication for Microsoft Graph API.
func (g *GraphHelper) AuthenticateGraphHelper() error {
	config := configuration.GetConfig()
	appScopes := config.GRAPH_APP_SCOPES
	tenantID := config.GRAPH_TENANT_ID
	clientID := config.GRAPH_CLIENT_ID
	clientSecret := config.GRAPH_CLIENT_SECRET

	g.AppScopes = strings.Split(appScopes, ",")

	// Create Azure credentials
	credential, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		return fmt.Errorf("failed to create credentials: %w", err)
	}
	g.Credential = credential

	// Create authentication provider
	authProvider, err := auth.NewAzureIdentityAuthenticationProviderWithScopes(g.Credential, g.AppScopes)
	if err != nil {
		return fmt.Errorf("failed to create authentication provider: %w", err)
	}

	// Create request adapter
	g.Adapter, err = msgraphsdk.NewGraphRequestAdapter(authProvider)
	if err != nil {
		return fmt.Errorf("failed to create Graph request adapter: %w", err)
	}

	// Initialize Graph client
	g.Client, err = msgraphsdk.NewGraphServiceClientWithCredentials(g.Credential, g.AppScopes)
	if err != nil {
		return fmt.Errorf("failed to initialize Graph client: %w", err)
	}

	return nil
}
