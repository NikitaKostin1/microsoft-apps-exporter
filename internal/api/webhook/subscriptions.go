package webhook

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"microsoft-apps-exporter/internal/sync"
)

// Request body
type SubscriptionLifecycleBody struct {
	Value []struct {
		SubscriptionId string `json:"subscriptionId"`
	} `json:"value"`
}

// newSubscriptionHandler handles subscription-related webhook notifications.
func newSubscriptionHandler(syncer *sync.Syncer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Received subscription lifecycle notification", "operation", "webhook")
		defer r.Body.Close()

		if r.Method != http.MethodPost {
			handleMethodNotAllowed(w, fmt.Sprintf("Only POST method allowed, got: "+r.Method))
			return
		}

		// Validate subscription creation
		if validated, err := handleValidationToken(w, r.URL); err != nil {
			handleBadRequest(w, fmt.Sprintf("failed to validate token: %v", err))
			return
		} else if validated { // subscription creation doesnt include resource update
			return
		}

		// Extract subscription ID from request payload
		subscriptionID, err := extractSubscriptionLifecycleData(r)
		if err != nil {
			handleBadRequest(w, fmt.Sprintf("invalid subscription payload: %v", err))
			return
		}

		// Reauthorize the subscription
		_, err = syncer.Graph.UpdateSubscription(subscriptionID)
		if err != nil {
			handleInternalError(w, fmt.Sprintf("failed to reauthorize subscription: %v", err))
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

// extractSubscriptionLifecycleData extracts the subscription ID from the JSON payload.
func extractSubscriptionLifecycleData(r *http.Request) (string, error) {
	var lifecycleBody SubscriptionLifecycleBody
	if err := json.NewDecoder(r.Body).Decode(&lifecycleBody); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	if len(lifecycleBody.Value) == 0 || lifecycleBody.Value[0].SubscriptionId == "" {
		return "", fmt.Errorf("missing subscriptionId in the lifecycleBody")
	}

	return lifecycleBody.Value[0].SubscriptionId, nil
}
