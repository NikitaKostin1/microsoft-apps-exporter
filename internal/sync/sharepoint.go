package sync

import (
	"fmt"
	"log/slog"
	"microsoft-apps-exporter/internal/api"
	"microsoft-apps-exporter/internal/models"
)

// SyncSharepoint synchronizes a SharePoint list and its items for a given site and list ID.
func (s *Syncer) SyncSharepoint(list models.ListReference) error {

	slog.Info("Syncing SharePoint list", "site_id", list.SiteID, "list_id", list.ListID,
		"database_table", list.DbTableName, "operation", "sync")

	if err := s.syncList(list); err != nil {
		return fmt.Errorf("failed to sync list: %w", err)
	}

	if err := s.syncListItems(list); err != nil {
		if cleanupErr := s.Database.DeleteDeltaLink(list.ListID); cleanupErr != nil {
			return fmt.Errorf("failed to sync list items: %w; cleanup failed: %v", err, cleanupErr)
		}

		slog.Info("Delta link cleared due to encountered exception during sync operation",
			"site_id", list.SiteID, "list_id", list.ListID, "operation", "sync")

		return fmt.Errorf("failed to sync list items: %w", err)
	}

	return nil
}

// syncList synchronizes the metadata of a SharePoint list.
func (s *Syncer) syncList(list models.ListReference) error {
	dbList, err := s.Database.GetList(list.ListID)
	if err != nil {
		return fmt.Errorf("failed to retrieve list from database: %w", err)
	}

	apiList, err := s.Graph.GetList(list.SiteID, list.ListID)
	if err != nil {
		return fmt.Errorf("failed to retrieve list from API: %w", err)
	}

	// Calculate the differences between the database and API lists
	toInsert, toUpdate, toDelete := diffFull(
		dbList, apiList,
		func(l models.ListMetadata) string { return l.ID },
		func(l models.ListMetadata) string { return l.ETag },
	)

	slog.Info("Syncing SharePoint list metadata", "site_id", list.SiteID, "list_id", list.ListID,
		slog.Group("changes", "to_insert", len(toInsert), "to_update", len(toUpdate), "to_delete", len(toDelete)),
		"operation", "sync")

	if err := s.Database.InsertLists(&toInsert); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}

	for _, list := range toUpdate {
		if err := s.Database.UpdateListIgnoreDelta(list); err != nil {
			return fmt.Errorf("failed to update: %w", err)
		}
	}

	for _, id := range toDelete {
		if err := s.Database.DeleteList(id); err != nil {
			return fmt.Errorf("failed to delete: %w", err)
		}
	}
	return nil
}

// syncListItems synchronizes SharePoint list items.
func (s *Syncer) syncListItems(list models.ListReference) error {
	dbTable, columnsMap := list.DbTableName, list.ColumnsMap
	deltaLink, err := s.Database.GetDeltaLink(list.ListID)
	if err != nil {
		return fmt.Errorf("failed to retrieve delta link: %w", err)
	}

	dbItems, err := s.Database.GetListItems(dbTable, list.SiteID, list.ListID)
	if err != nil {
		return fmt.Errorf("failed to retrieve list items from database: %w", err)
	}

	expandFields := make([]string, 0, len(columnsMap))
	for _, column := range columnsMap {
		expandFields = append(expandFields, column)
	}

	options := api.NewListItemsWithDeltaOptions(expandFields, nil)
	newDeltaLink, apiItems, err := s.Graph.GetListItemsWithDelta(list.SiteID, list.ListID, deltaLink, options)
	if err != nil {
		return fmt.Errorf("failed to retrieve list items from API: %w", err)
	}

	if newDeltaLink != nil {
		if err := s.Database.SaveDeltaLink(list.ListID, *newDeltaLink); err != nil {
			return fmt.Errorf("failed to save delta link: %w", err)
		}
	}

	// Determine insert, update, delete actions
	var toInsert, toUpdate []models.ListItem
	var toDelete []string
	if deltaLink != nil { // Delta synchronization
		toInsert, toUpdate, toDelete = diffDelta(*dbItems, *apiItems,
			func(li models.ListItem) string { return li.Metadata.ID },
			func(li models.ListItem) string { return li.Metadata.ETag },
		)
	} else { // Full synchronization
		toInsert, toUpdate, toDelete = diffFull(*dbItems, *apiItems,
			func(li models.ListItem) string { return li.Metadata.ID },
			func(li models.ListItem) string { return li.Metadata.ETag },
		)
	}

	slog.Info("Syncing SharePoint list items",
		"with_delta", deltaLink != nil, "site_id", list.SiteID, "list_id", list.ListID,
		slog.Group("changes", "to_insert", len(toInsert), "to_update", len(toUpdate), "to_delete", len(toDelete)),
		"operation", "sync")

	if err := s.Database.InsertListItems(dbTable, columnsMap, &toInsert); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}

	for _, item := range toUpdate {
		if err := s.Database.UpdateListItem(dbTable, columnsMap, item); err != nil {
			return fmt.Errorf("failed to update: %w", err)
		}
	}

	for _, id := range toDelete {
		if err := s.Database.DeleteListItem(dbTable, id); err != nil {
			return fmt.Errorf("failed to delete: %w", err)
		}
	}

	return nil
}
