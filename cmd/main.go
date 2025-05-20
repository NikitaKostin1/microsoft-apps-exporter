package main

import (
	"context"
	"log/slog"
	"microsoft-apps-exporter/internal/api"
	"microsoft-apps-exporter/internal/api/webhook"
	"microsoft-apps-exporter/internal/database"
	"microsoft-apps-exporter/internal/logging"
	"microsoft-apps-exporter/internal/sync"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := makeSignalNotify(cancel)

	logging.ConfigureSlog()
	slog.Info("Application starting")

	// Establish database connection.
	db, err := database.NewDatabase()
	if err != nil {
		slog.Error("Failed to create Database instance", "exception", err)
		return
	}
	defer db.Close()

	// Initiate API client.
	graphHelper, err := api.NewGraphHelper(ctx)
	if err != nil {
		slog.Error("Failed to create GraphHelper instance", "exception", err)
		return
	}

	syncer := sync.NewSyncer(graphHelper, db)

	// Start Webhook Server to listen for Change Notifications.
	webhookServer := webhook.NewWebhookServer(syncer)
	if err := webhookServer.RunAsync(); err != nil {
		slog.Error("Failed to run webhook server", "exception", err)
		return
	}
	defer webhookServer.Shutdown(ctx)

	// Ensures subscribed to Microsoft Graph API notificaitions.
	if _, err := graphHelper.EnsureResourcesSubscriptions(); err != nil {
		slog.Error("Failed to validate MS Graph API subscriptions", "exception", err)
		return
	}

	// Start synchronization.
	if err := syncer.SyncResources(); err != nil {
		slog.Error("Failed to sync resources", "exception", err)
		return
	}

	<-stop
}

// Channel to listen for OS signals and gracefully exit.
func makeSignalNotify(cancel context.CancelFunc) <-chan struct{} {
	stop := make(chan struct{})
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChannel
		slog.Info("Received termination signal. Shutting down gracefully...")
		cancel()
		close(stop)
	}()

	return stop
}
