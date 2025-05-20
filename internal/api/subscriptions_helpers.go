package api

import (
	"fmt"
	"microsoft-apps-exporter/internal/configuration"

	gmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
)

// EnsureSubscription checks if a subscription for the specified SharePoint resource exists.
// If an existing subscription is found and has the correct webhook URL, it is returned.
// Otherwise, it creates a new subscription after deleting outdated ones.
func (g *GraphHelper) ensureResourceSubscription(resource, webhookResourceEndpoint string) (gmodels.Subscriptionable, error) {
	subscriptions, err := g.GetSubscriptions()
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}

	for _, sub := range subscriptions {
		if *sub.GetResource() == resource {

			if g.isSubscriptionWebhookURLsMatch(sub, webhookResourceEndpoint) {
				return sub, nil
			} else {
				if err := g.DeleteSubscription(*sub.GetId()); err != nil {
					return nil, fmt.Errorf("failed to delete subscription with wrong mismatched webhook URL: %w", err)
				}
			}

			break
		}
	}

	return g.CreateResourceSubscription(resource, webhookResourceEndpoint)
}

// deleteInactiveSubscriptions removes subscriptions for resources that are no longer active.
func (g *GraphHelper) deleteInactiveSubscriptions(activeResources map[string]struct{}) error {
	existingSubscriptions, err := g.GetSubscriptions()
	if err != nil {
		return fmt.Errorf("failed to get existing subscriptions: %w", err)
	}

	for _, sub := range existingSubscriptions {
		resource := *sub.GetResource()

		if _, exists := activeResources[resource]; !exists { // If not marked active, delete
			if err := g.DeleteSubscription(*sub.GetId()); err != nil {
				return fmt.Errorf("failed to delete inactive subscription for resource %s: %w", resource, err)
			}
		}
	}
	return nil
}

// isSubscriptionWebhookURLsMatch checks if the given subscription has outdated webhook URLs.
func (g *GraphHelper) isSubscriptionWebhookURLsMatch(subscription gmodels.Subscriptionable, webhookResourceEndpoint string) bool {
	config := configuration.GetConfig()
	webhookBaseURL := config.WEBHOOK_EXTERNAL_BASE_URL

	expectedNotificationURL := webhookBaseURL + webhookResourceEndpoint
	expectedLifecycleURL := webhookBaseURL + webhookSubscriptionEndpoint

	actualNotificationURL := *subscription.GetNotificationUrl()
	actualLifecycleURL := *subscription.GetLifecycleNotificationUrl()

	return actualNotificationURL == expectedNotificationURL && actualLifecycleURL == expectedLifecycleURL
}
