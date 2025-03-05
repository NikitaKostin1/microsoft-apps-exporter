//go:build testing && unit

package webhook_test

import (
	"bytes"
	"encoding/json"
	"io"
	"microsoft-apps-exporter/internal/api/webhook"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestExtractSubscriptionLifecycleData verifies the extraction of subscription lifecycle data from HTTP requests.
func TestExtractSubscriptionLifecycleData(t *testing.T) {
	t.Run("Valid payload", func(t *testing.T) {
		payload := webhook.SubscriptionLifecycleBody{
			Value: []struct {
				SubscriptionId string `json:"subscriptionId"`
			}{
				{SubscriptionId: "sub-123"},
			},
		}
		body, _ := json.Marshal(payload)
		r := httptest.NewRequest(http.MethodPost, "http://example.com", bytes.NewReader(body))
		id, err := webhook.ExtractSubscriptionLifecycleData(r)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if id != "sub-123" {
			t.Fatalf("expected subscriptionId 'sub-123', got '%s'", id)
		}
	})

	t.Run("Invalid payload", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "http://example.com", io.NopCloser(bytes.NewReader([]byte("invalid"))))
		_, err := webhook.ExtractSubscriptionLifecycleData(r)

		if err == nil {
			t.Fatal("expected an error due to invalid JSON payload")
		}
	})

	t.Run("Missing", func(t *testing.T) {
		payload := webhook.SubscriptionLifecycleBody{
			Value: []struct {
				SubscriptionId string `json:"subscriptionId"`
			}{
				{},
			},
		}
		body, _ := json.Marshal(payload)
		r := httptest.NewRequest(http.MethodPost, "http://example.com", bytes.NewReader(body))
		_, err := webhook.ExtractSubscriptionLifecycleData(r)

		if err == nil {
			t.Fatal("expected an error due to missing subscriptionId")
		}
	})

}
