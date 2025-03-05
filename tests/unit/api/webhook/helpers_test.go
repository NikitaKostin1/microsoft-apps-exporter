//go:build testing && unit

package webhook_test

import (
	"microsoft-apps-exporter/internal/api/webhook"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// TestHandleValidationToken verifies the handling of validation tokens in webhook requests.
func TestHandleValidationToken(t *testing.T) {
	w := httptest.NewRecorder()

	t.Run("Valid token", func(t *testing.T) {
		u, _ := url.Parse("http://example.com?validationToken=test-token")
		validated, err := webhook.HandleValidationToken(w, u)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !validated {
			t.Fatal("expected validation token to be handled")
		}
		if w.Code != http.StatusOK {
			t.Fatalf("expected status code 200, got %d", w.Code)
		}
	})

	t.Run("No Token", func(t *testing.T) {
		u, _ := url.Parse("http://example.com")
		validated, err := webhook.HandleValidationToken(w, u)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if validated {
			t.Fatal("expected validation token to be false when no token is provided")
		}
	})
}
