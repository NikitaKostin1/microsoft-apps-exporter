//go:build testing && unit

package webhook_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"microsoft-apps-exporter/internal/api/webhook"
	"microsoft-apps-exporter/internal/models"
	"microsoft-apps-exporter/internal/sync"

	"github.com/stretchr/testify/assert"
)

func TestNewWebhookServer(t *testing.T) {
	setupTestResourcesYaml()
	syncer := &sync.Syncer{}

	testPort := "8080"
	testListenIp := "0.0.0.0"
	testAddr := strings.Join([]string{testListenIp, ":", testPort}, "")
	os.Setenv("WEBHOOK_LISTEN_PORT", testPort)
	os.Setenv("WEBHOOK_LISTEN_IP", testListenIp)

	webhookServer := webhook.NewWebhookServer(syncer)
	assert.NotNil(t, webhookServer)
	assert.Equal(t, testAddr, webhookServer.Server.Addr)

	tests := []struct {
		name       string
		path       string
		method     string
		wantStatus int
	}{
		{
			name:       "Subscription Notification Handler",
			path:       "/webhook/subscription-notification",
			method:     http.MethodPost,
			wantStatus: http.StatusBadRequest, // Assuming no valid request body is sent
		},
		{
			name:       "SharePoint Webhook Handler",
			path:       models.WebhookSharepointEndpoint,
			method:     http.MethodPost,
			wantStatus: http.StatusBadRequest, // Assuming no valid request body is sent
		},
		{
			name:       "Ping Handler",
			path:       "/webhook/ping",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid Endpoint",
			path:       "/webhook/invalid",
			method:     http.MethodGet,
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			webhookServer.Server.Handler.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
