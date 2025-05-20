package api

import (
	"fmt"
	"log/slog"
	"microsoft-apps-exporter/internal/configuration"
	"microsoft-apps-exporter/internal/models"
	"time"

	gmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
)

const (
	webhookSubscriptionEndpoint = "/webhook/subscription-notification"
	defaultSubscriptionExpiry   = 48 * time.Hour
	subscriptionUpdateExpiry    = 72 * time.Hour
)

// EnsureResourcesSubscriptions ensures that subscriptions exist for all configured resources.
// It returns a slice of active subscriptions or an error if the process fails.
func (g *GraphHelper) EnsureResourcesSubscriptions() ([]gmodels.Subscriptionable, error) {
	var subscriptions []gmodels.Subscriptionable
	activeResources := make(map[string]struct{})

	slog.Info("Ensuring MS Graph API resources subscriptions are active", "operation", "subscriptions")

	config := configuration.GetConfig()
	if config.Sharepoint != nil {
		// Ensure subscriptions exist for all configured SharePoint resources
		for _, list := range config.Sharepoint.Lists {
			resource := models.GenerateSharepointResourceString(list.SiteID, list.ListID)
			activeResources[resource] = struct{}{} // Mark as active

			subscription, err := g.ensureResourceSubscription(resource, models.WebhookSharepointEndpoint)
			if err != nil {
				return nil, fmt.Errorf("failed to ensure subscription for resource %s: %w", resource, err)
			}
			subscriptions = append(subscriptions, subscription)
		}
	}

	err := g.deleteInactiveSubscriptions(activeResources)
	if err != nil {
		return nil, err
	}

	slog.Info("Active MS Graph API subscriptions validated", "count", len(subscriptions), "operation", "subscriptions")
	return subscriptions, nil
}

// CreateSharepointSubscription creates a new subscription for the specified SharePoint resource.
func (g *GraphHelper) CreateResourceSubscription(resource, webhookResourceEndpoint string) (gmodels.Subscriptionable, error) {
	config := configuration.GetConfig()
	webhookBaseURL := config.WEBHOOK_EXTERNAL_BASE_URL

	requestBody := gmodels.NewSubscription()
	changeType := "updated"
	notificationUrl := webhookBaseURL + webhookResourceEndpoint
	lifecycleNotificationUrl := webhookBaseURL + webhookSubscriptionEndpoint
	expirationDateTime := time.Now().Add(defaultSubscriptionExpiry)
	latestSupportedTlsVersion := "v1_2"

	requestBody.SetChangeType(&changeType)
	requestBody.SetNotificationUrl(&notificationUrl)
	requestBody.SetLifecycleNotificationUrl(&lifecycleNotificationUrl)
	requestBody.SetResource(&resource)
	requestBody.SetExpirationDateTime(&expirationDateTime)
	requestBody.SetLatestSupportedTlsVersion(&latestSupportedTlsVersion)

	subscription, err := g.Client.Subscriptions().Post(g.Ctx, requestBody, nil)
	if err != nil {
		slog.Debug("Failed to create subscription",
			slog.Group("requestBody",
				"changeType", changeType,
				"notificationUrl", notificationUrl,
				"lifecycleNotificationUrl", lifecycleNotificationUrl,
				"resource", resource,
				"expirationDateTime", expirationDateTime,
				"latestSupportedTlsVersion", latestSupportedTlsVersion),
			"exception", err, "operation", "subscriptions")

		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	slog.Info("Subscription created successfully", "subscription_id", *subscription.GetId(),
		"resource", resource, "operation", "subscriptions")
	return subscription, nil
}

// GetSubscriptions retrieves all active subscriptions.
func (g *GraphHelper) GetSubscriptions() ([]gmodels.Subscriptionable, error) {
	subscriptions, err := g.Client.Subscriptions().Get(g.Ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to request subscriptions: %w", err)
	}

	slog.Debug("Client requested subscriptions", "count", len(subscriptions.GetValue()), "operation", "subscriptions")
	return subscriptions.GetValue(), nil
}

// UpdateSubscription updates the expiration time of a specific subscription.
func (g *GraphHelper) UpdateSubscription(subscriptionID string) (gmodels.Subscriptionable, error) {
	requestBody := gmodels.NewSubscription()
	expirationDateTime := time.Now().Add(subscriptionUpdateExpiry)
	requestBody.SetExpirationDateTime(&expirationDateTime)

	subscription, err := g.Client.Subscriptions().BySubscriptionId(subscriptionID).Patch(g.Ctx, requestBody, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to update subscription %s: %w", subscriptionID, err)
	}

	slog.Info("Subscription updated successfully", "subscription_id", subscriptionID, "operation", "subscriptions")
	return subscription, nil
}

// DeleteSubscription deletes a specific subscription by its ID.
func (g *GraphHelper) DeleteSubscription(subscriptionID string) error {
	err := g.Client.Subscriptions().BySubscriptionId(subscriptionID).Delete(g.Ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to delete subscription %s: %w", subscriptionID, err)
	}

	slog.Info("Subscription deleted successfully", "subscription_id", subscriptionID, "operation", "subscriptions")
	return nil
}

/* // ReauthorizeSubscription reauthorizes a specific subscription by its ID.
func (g *GraphHelper) ReauthorizeSubscription(subscriptionID string) error {
	err := g.Client.Subscriptions().BySubscriptionId(subscriptionID).Reauthorize().Post(g.Ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to reauthorize subscription %s: %w", subscriptionID, err)
	}

	slog.Info("Subscription reauthorized successfully", "subscription_id", subscriptionID, "operation", "subscriptions")
	return nil
} */
