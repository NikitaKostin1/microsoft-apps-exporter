//go:build testing && integration

package webhook_test

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"microsoft-apps-exporter/internal/api/webhook"
	"microsoft-apps-exporter/internal/configuration"
	"microsoft-apps-exporter/internal/sync"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWebhookServerIntegration(t *testing.T) {
	slog.SetLogLoggerLevel(math.MaxInt) // Disable logging

	// Helper to create a test server with config overrides
	createTestServer := func(ip, port string) *webhook.WebhookServer {
		os.Setenv("WEBHOOK_LISTEN_IP", ip)
		os.Setenv("WEBHOOK_LISTEN_PORT", port)
		os.Setenv("WEBHOOK_EXTERNAL_BASE_URL", fmt.Sprintf("http://%s:%s", ip, port))

		configuration.ResetConfig()
		return webhook.NewWebhookServer(&sync.Syncer{})
	}

	t.Run("RunAsync success case", func(t *testing.T) {
		server := createTestServer("localhost", "1026")
		defer server.Shutdown(context.Background())

		err := server.RunAsync()
		assert.NoError(t, err, "RunAsync should succeed")
	})

	t.Run("RunAsync with startup error (port in use)", func(t *testing.T) {
		// Create a blocker server first
		blocker := createTestServer("localhost", "9876")
		go blocker.Server.ListenAndServe()
		defer blocker.Server.Close()
		time.Sleep(300 * time.Millisecond)

		server := createTestServer("localhost", "9876") // Same port as blocker
		defer server.Shutdown(context.Background())

		err := server.RunAsync()
		assert.Error(t, err, "RunAsync should fail when port is in use")
		assert.Contains(t, err.Error(), "address already in use", "Error should indicate port conflict")
	})

	t.Run("PingWebhookServer success case", func(t *testing.T) {
		server := createTestServer("localhost", "1026")
		defer server.Shutdown(context.Background())

		// Start server manually for this test
		go server.Server.ListenAndServe()
		time.Sleep(100 * time.Millisecond)

		err := server.PingWebhookServer()
		assert.NoError(t, err, "Ping should succeed when server is running")
	})

	t.Run("PingWebhookServer failure case", func(t *testing.T) {
		server := createTestServer("localhost", "1026")
		defer server.Shutdown(context.Background())

		// Don't start the server - should fail to ping
		err := server.PingWebhookServer()
		assert.Error(t, err, "Ping should fail when server isn't running")
	})

	t.Run("RunAsync with ping timeout", func(t *testing.T) {
		// Create server that will never respond to pings
		server := createTestServer("invalid-host", "1234")
		defer server.Shutdown(context.Background())

		err := server.RunAsync()
		assert.Error(t, err, "Should fail when ping times out")
	})

	t.Run("RunAsync error channel receives startup failure", func(t *testing.T) {
		// Create invalid server config (invalid port)
		server := createTestServer("localhost", "99999")
		defer server.Shutdown(context.Background())

		err := server.RunAsync()
		assert.Error(t, err, "Should propagate startup errors from error channel")
	})

	configuration.ResetConfig()
}
