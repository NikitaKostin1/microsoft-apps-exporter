//go:build testing && e2e

package api_test

import (
	"context"
	"log/slog"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"microsoft-apps-exporter/internal/api"
	"microsoft-apps-exporter/internal/api/webhook"
	"microsoft-apps-exporter/internal/configuration"
	"microsoft-apps-exporter/internal/models"
	"microsoft-apps-exporter/internal/sync"
)

func setupTest(t *testing.T) (context.Context, *api.GraphHelper, *webhook.WebhookServer) {
	slog.SetLogLoggerLevel(math.MaxInt) // Disable logging
	t.Helper()

	setupProdResourcesYaml()

	ctx := context.Background()

	// Initialize Graph client
	graph, err := api.NewGraphHelper(ctx)
	require.NoError(t, err, "failed to initialize GraphHelper")

	// Start webhook server
	server := webhook.NewWebhookServer(&sync.Syncer{Graph: graph})
	err = server.RunAsync()
	require.NoError(t, err, "failed to start webhook server")

	t.Cleanup(func() {
		subscriptions, _ := graph.GetSubscriptions()

		for _, sub := range subscriptions {
			if err := graph.DeleteSubscription(*sub.GetId()); err != nil {
				require.NoError(t, err, "DeleteSubscription failed")
			}
		}
		server.Shutdown(ctx)
	})

	return ctx, graph, server
}

func TestEnsureResourcesSubscriptions(t *testing.T) {
	_, graph, _ := setupTest(t)

	config := configuration.GetConfig()
	if config.Sharepoint == nil {
		t.Fatalf("Sharepoint resource is expected to be specified")
	}

	expectedCount := len(config.Sharepoint.Lists)
	subscriptions, err := graph.GetSubscriptions()

	for _, sub := range subscriptions {
		if err := graph.DeleteSubscription(*sub.GetId()); err != nil {
			require.NoError(t, err, "DeleteSubscription failed")
		}
	}

	// Subscriptions dont exist case
	subscriptions, err = graph.EnsureResourcesSubscriptions()
	require.NoError(t, err, "EnsureResourcesSubscriptions failed")
	assert.Equal(t, expectedCount, len(subscriptions), "expected same number of subscriptions as resources")

	//  Subscriptions exist case
	subscriptions, err = graph.EnsureResourcesSubscriptions()
	require.NoError(t, err, "EnsureResourcesSubscriptions failed")
	assert.Equal(t, expectedCount, len(subscriptions), "expected same number of subscriptions as resources")

}

func TestCreateResourceSubscription(t *testing.T) {
	_, graph, _ := setupTest(t)

	subscriptions, _ := graph.GetSubscriptions()
	for _, sub := range subscriptions {
		if err := graph.DeleteSubscription(*sub.GetId()); err != nil {
			require.NoError(t, err, "DeleteSubscription failed")
		}
	}

	config := configuration.GetConfig()
	if config.Sharepoint == nil {
		t.Fatalf("Sharepoint resource is expected to be specified")
	}

	for _, list := range config.Sharepoint.Lists {
		resource := models.GenerateSharepointResourceString(list.SiteID, list.ListID)

		sub, err := graph.CreateResourceSubscription(resource, api.WebhookSubscriptionEndpoint)
		require.NoError(t, err, "CreateResourceSubscription failed")
		assert.Equal(t, resource, *sub.GetResource())

		subscriptions, _ := graph.GetSubscriptions()
		found := false
		for _, existingSub := range subscriptions {
			if *existingSub.GetId() == *sub.GetId() {
				found = true
				break
			}
		}
		assert.True(t, found, "Created subscription not found in the list of subscriptions")
	}
}

func TestGetSubscriptions(t *testing.T) {
	_, graph, _ := setupTest(t)

	expectedSubscriptions, err := graph.EnsureResourcesSubscriptions()
	require.NoError(t, err, "EnsureResourcesSubscriptions failed")

	subscriptions, err := graph.GetSubscriptions()
	require.NoError(t, err, "GetSubscriptions failed")
	assert.Equal(t, len(expectedSubscriptions), len(subscriptions), "expected same number of subscriptions as resources")
}

func TestUpdateSubscription(t *testing.T) {
	_, graph, _ := setupTest(t)

	expectedSubscriptions, err := graph.EnsureResourcesSubscriptions()
	require.NoError(t, err, "EnsureResourcesSubscriptions failed")

	for _, sub := range expectedSubscriptions {
		updated, err := graph.UpdateSubscription(*sub.GetId())
		require.NoError(t, err, "UpdateSubscription failed")

		expectedExpiry := time.Now().Add(api.SubscriptionUpdateExpiry)
		actualExpiry := *updated.GetExpirationDateTime()
		assert.WithinDuration(t, expectedExpiry, actualExpiry, 5*time.Minute)
	}
}

func TestDeleteSubscription(t *testing.T) {
	_, graph, _ := setupTest(t)

	expectedSubscriptions, err := graph.EnsureResourcesSubscriptions()
	require.NoError(t, err, "EnsureResourcesSubscriptions failed")

	for _, sub := range expectedSubscriptions {
		err = graph.DeleteSubscription(*sub.GetId())
		require.NoError(t, err, "DeleteSubscription failed")

		// Validate that the subscription was deleted
		subscriptions, err := graph.GetSubscriptions()
		require.NoError(t, err, "GetSubscriptions failed")

		for _, existingSub := range subscriptions {
			assert.NotEqual(t, *sub.GetId(), *existingSub.GetId(), "Deleted subscription still exists")
		}
	}
}

// func TestReauthorizeSubscription(t *testing.T) {
// 	_, graph, _ := setupTest(t)

// 	expectedSubscriptions, err := graph.EnsureResourcesSubscriptions()
// 	require.NoError(t, err, "EnsureResourcesSubscriptions failed")

// 	for _, sub := range expectedSubscriptions {
// 		err = graph.ReauthorizeSubscription(*sub.GetId())
// 		require.NoError(t, err, "ReauthorizeSubscription failed")
// 	}
// }
