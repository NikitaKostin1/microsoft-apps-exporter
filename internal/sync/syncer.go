package sync

import (
	"fmt"
	"log/slog"
	"microsoft-apps-exporter/internal/api"
	"microsoft-apps-exporter/internal/configuration"
	"microsoft-apps-exporter/internal/database"
)

// Syncer is responsible for synchronizing data between the database and MS Graph API.
type Syncer struct {
	Graph    *api.GraphHelper
	Database *database.Database
}

// NewSyncer creates a new Syncer with the provided database and API clients.
func NewSyncer(graph *api.GraphHelper, db *database.Database) *Syncer {
	return &Syncer{Graph: graph, Database: db}
}

// SyncResources synchronizes all resources from config between the database and the API.
func (s *Syncer) SyncResources() error {
	config := configuration.GetConfig()
	slog.Info("Starting resource synchronization", "operation", "sync")

	if config.Sharepoint != nil {
		slog.Debug("SharePoint resource found in config", "database_table", config.Sharepoint.DbTableName, "operation", "sync")

		for _, list := range config.Sharepoint.Lists {
			if err := s.SyncSharepoint(list); err != nil {
				return fmt.Errorf("failed to sync SharePoint resource: %w", err)
			}
		}
	}

	slog.Info("Initial resource synchronization completed. Further sync will occur on webhook Change Notificaiton.", "operation", "sync")

	return nil
}
