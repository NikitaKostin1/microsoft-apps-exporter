//go:build testing && unit

package webhook_test

import (
	"fmt"
	"log/slog"
	"math"
	"microsoft-apps-exporter/internal/api/webhook"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHandleMethodNotAllowed verifies that requests with unsupported methods return a 405 response.
func TestHandleMethodNotAllowed(t *testing.T) {
	slog.SetLogLoggerLevel(math.MaxInt) // Disable logging
	w := httptest.NewRecorder()
	webhook.HandleMethodNotAllowed(w, fmt.Sprintf("Only POST method allowed, got: "+"GET"))

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status code 405, got %d", w.Code)
	}
}

// TestHandleBadRequest verifies that bad requests return a 400 response.
func TestHandleBadRequest(t *testing.T) {
	slog.SetLogLoggerLevel(math.MaxInt) // Disable logging
	w := httptest.NewRecorder()
	webhook.HandleBadRequest(w, "Invalid request")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status code 400, got %d", w.Code)
	}
}

// TestHandleInternalError verifies that internal errors return a 500 response.
func TestHandleInternalError(t *testing.T) {
	slog.SetLogLoggerLevel(math.MaxInt) // Disable logging
	w := httptest.NewRecorder()
	webhook.HandleInternalError(w, "Internal failure")

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status code 500, got %d", w.Code)
	}
}
