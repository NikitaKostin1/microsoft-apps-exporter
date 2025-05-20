//go:build testing && e2e

package api_test

import (
	"context"
	"testing"

	"microsoft-apps-exporter/internal/api"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGraphHelper_Initialize validates that GraphHelper initializes properly
func TestGraphHelper_Initialize(t *testing.T) {
	ctx := context.Background()
	graphHelper, err := api.NewGraphHelper(ctx)

	assert.NoError(t, err, "GraphHelper should initialize without errors")
	require.NotNil(t, graphHelper, "GraphHelper should not be nil")
	assert.NotNil(t, graphHelper.Client, "GraphHelper should have a valid Graph client")
}

// TestGraphHelper_Authentication ensures authentication succeeds with real credentials
func TestGraphHelper_Authentication(t *testing.T) {
	graphHelper := &api.GraphHelper{}
	err := graphHelper.AuthenticateGraphHelper()

	assert.NoError(t, err, "Authentication should succeed")
	assert.NotNil(t, graphHelper.Credential, "Credential should not be nil")
	assert.NotNil(t, graphHelper.Adapter, "Adapter should be initialized")
	assert.NotNil(t, graphHelper.Client, "Client should be initialized")
}
