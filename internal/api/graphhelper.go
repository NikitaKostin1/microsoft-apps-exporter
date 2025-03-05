package api

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

// GraphClient handles authentication and communication with Microsoft Graph API.
type GraphHelper struct {
	Ctx        context.Context
	Credential *azidentity.ClientSecretCredential
	Adapter    *msgraphsdk.GraphRequestAdapter
	Client     *msgraphsdk.GraphServiceClient
	AppScopes  []string
}

// NewGraphClient initializes and authenticates a new GraphClient instance.
func NewGraphHelper(ctx context.Context) (*GraphHelper, error) {
	g := &GraphHelper{Ctx: ctx}

	if err := g.AuthenticateGraphHelper(); err != nil {
		return nil, fmt.Errorf("failed to authenticate GraphHelper: %w", err)
	}
	return g, nil
}
