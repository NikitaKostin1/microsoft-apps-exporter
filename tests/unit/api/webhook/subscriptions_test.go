//go:build testing && unit

package webhook_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"math"
	"microsoft-apps-exporter/internal/api/webhook"
	"microsoft-apps-exporter/internal/sync"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSubscriptionHandler(t *testing.T) {
	syncer := &sync.Syncer{}
	slog.SetLogLoggerLevel(math.MaxInt) // Disable logging
	handler := webhook.NewSubscriptionHandler(syncer)

	invalidPayload := []byte(`{"invalid": "data"}`)
	missingSubIDPayload := []byte(`{"value": [{}]}`)

	tests := []struct {
		name           string
		method         string
		urlQuery       string
		body           io.Reader
		expectedStatus int
	}{
		{
			name:           "Invalid method",
			method:         http.MethodGet,
			body:           nil,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Token validation",
			method:         http.MethodPost,
			urlQuery:       "?validationToken=abc123",
			body:           nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid JSON payload",
			method:         http.MethodPost,
			body:           bytes.NewReader(invalidPayload),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Missing subscription ID",
			method:         http.MethodPost,
			body:           bytes.NewReader(missingSubIDPayload),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/subscription" + tt.urlQuery
			req := httptest.NewRequest(tt.method, url, tt.body)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

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
