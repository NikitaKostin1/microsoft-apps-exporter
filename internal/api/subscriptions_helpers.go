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
			isOutdated, err := g.deleteOutdatedSubscription(sub, webhookResourceEndpoint)
			if err != nil {
				return nil, err
			}
			if isOutdated {
				continue
			}
			return sub, nil
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

// deleteOutdatedSubscription checks if the given subscription has outdated webhook URLs.
// If the subscription is outdated, it is deleted, and the function returns true.
func (g *GraphHelper) deleteOutdatedSubscription(subscription gmodels.Subscriptionable, webhookResourceEndpoint string) (bool, error) {
	config := configuration.GetConfig()
	webhookBaseURL := config.WEBHOOK_BASE_URL

	expectedNotificationURL := webhookBaseURL + webhookResourceEndpoint
	expectedLifecycleURL := webhookBaseURL + webhookSubscriptionEndpoint

	actualNotificationURL := *subscription.GetNotificationUrl()
	actualLifecycleURL := *subscription.GetLifecycleNotificationUrl()

	if actualNotificationURL != expectedNotificationURL || actualLifecycleURL != expectedLifecycleURL {
		if err := g.DeleteSubscription(*subscription.GetId()); err != nil {
			return true, fmt.Errorf("failed to delete outdated subscription: %w", err)
		}
		return true, nil
	}
	return false, nil
}
