package webhook

import (
	"context"
	"fmt"
	"log/slog"
	"microsoft-apps-exporter/internal/configuration"
	"microsoft-apps-exporter/internal/models"
	"microsoft-apps-exporter/internal/sync"
	"net/http"
	"strings"
	"time"
)

// WebhookServer manages the HTTP server for handling webhook notifications.
type WebhookServer struct {
	Port   string
	Server *http.Server
}

var pingEndpoint string = "/webhook/ping"

// NewWebhookServer initializes and configures a new WebhookServer instance.
func NewWebhookServer(syncer *sync.Syncer) *WebhookServer {
	config := configuration.GetConfig()
	port := config.WEBHOOK_LISTEN_PORT
	ip := config.WEBHOOK_LISTEN_IP

	// Define webhook routes and their handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook/subscription-notification", newSubscriptionHandler(syncer))
	mux.HandleFunc(models.WebhookSharepointEndpoint, newSharepointHandler(syncer))
	mux.HandleFunc(pingEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return &WebhookServer{
		Port: port,
		Server: &http.Server{
			Addr:    strings.Join([]string{ip, ":", port}, ""),
			Handler: mux,
		},
	}
}

// RunAsync starts the HTTP server and waits for it to be reachable.
func (ws *WebhookServer) RunAsync() error {
	errChan := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		slog.Info("Starting webhook server", "listen_address", ws.Server.Addr, "operation", "webhook")

		if err := ws.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("webhook server encountered an error: %w", err)
		}
		close(errChan)
	}()

	// Wait for the server to be available
	if err := ws.PingWebhookServer(); err != nil {
		return fmt.Errorf("webhook server failed to start in time: %w", err)
	}

	// Check if an error occurred during startup
	select {
	case err := <-errChan:
		if err != nil {
			return err
		}
	default:
	}

	slog.Info("Webhook server started", "listen_address", ws.Server.Addr, "operation", "webhook")
	return nil
}

// PingWebhookServer attempts to reach the webhook server until it responds or times out.
func (ws *WebhookServer) PingWebhookServer() error {
	config := configuration.GetConfig()
	url := config.WEBHOOK_EXTERNAL_BASE_URL + pingEndpoint
	slog.Info("Webhook server healtheck", "url", url, "operation", "webhook")

	timeout := 10 * time.Second
	interval := 500 * time.Millisecond
	deadline := time.Now().Add(timeout)

	var err error
	var resp *http.Response
	for time.Now().Before(deadline) {
		resp, err = http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			return nil
		}
		time.Sleep(interval)
	}

	slog.Debug("Failed to reach webhook server after multiple attempts",
		slog.Group("healthcheck", "url", url, "response", resp, "exception", err), "operation", "webhook")
	return fmt.Errorf("webhook server healtheck did not respond on ping: %v", err)
}

// Shutdown gracefully stops the webhook server, allowing ongoing requests to complete.
func (ws *WebhookServer) Shutdown(ctx context.Context) {
	// Allow up to 5 seconds for graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := ws.Server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Failed to shut down webhook server", "exception", err, "operation", "database")
	} else {
		slog.Info("Webhook server is shut down successfully", "operation", "webhook")
	}
}
