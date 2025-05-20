package webhook

import (
	"fmt"
	"log/slog"

	"encoding/json"
	"net/http"
	"strings"

	"microsoft-apps-exporter/internal/configuration"
	"microsoft-apps-exporter/internal/models"
	"microsoft-apps-exporter/internal/sync"
)

// Request body
type ResourceUpdateBody struct {
	Value []struct {
		Resource     string `json:"resource"`
		ResourceData struct {
			OdataType string `json:"@odata.type"`
		} `json:"resourceData"`
	} `json:"value"`
}

// newSharepointHandler returns an HTTP handler for processing SharePoint webhook notifications.
func newSharepointHandler(syncer *sync.Syncer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Received SharePoint update notification", "operation", "webhook")
		defer r.Body.Close()

		if r.Method != http.MethodPost {
			handleMethodNotAllowed(w, fmt.Sprintf("Only POST method allowed, got: "+r.Method))
			return
		}

		// Validate subscription creation
		if validated, err := handleValidationToken(w, r.URL); err != nil {
			handleBadRequest(w, fmt.Sprintf("failed to validate token: %v", err))
			return
		} else if validated { // subscription creation doesnt include resource update
			return
		}

		siteID, listID, err := extractResourceUpdateData(r)
		if err != nil {
			handleBadRequest(w, err.Error())
			return
		}

		list, found := exctractListReference(siteID, listID)

		// Validate if the struct was found
		if !found {
			handleBadRequest(w, "The resource doesn't exist")
			return
		}

		// Perform sync operation in a separate goroutine
		go func() {
			if err := syncer.SyncSharepoint(list); err != nil {
				slog.Error("failed to sync SharePoint resource: %w", "exception", err, "operation", "webhook")
			}
		}()

		w.WriteHeader(http.StatusOK)
	}
}

func exctractListReference(siteID, listID string) (models.ListReference, bool) {
	config := configuration.GetConfig()

	for _, list := range config.Sharepoint.Lists {
		if list.SiteID == siteID && list.ListID == listID {
			return list, true
		}
	}
	return models.ListReference{}, false
}

// extractResourceUpdateData validates the request body and extracts siteID, listID, and resource details.
func extractResourceUpdateData(r *http.Request) (string, string, error) {
	var updateBody ResourceUpdateBody
	var siteID, listID string
	var err error

	if err := json.NewDecoder(r.Body).Decode(&updateBody); err != nil {
		return "", "", fmt.Errorf("invalid request body: %w", err)
	}

	for _, updateUnit := range updateBody.Value {
		if updateUnit.ResourceData.OdataType != "#Microsoft.Graph.ListItem" {
			return "", "", fmt.Errorf("invalid data type: '%s', expected: '#Microsoft.Graph.ListItem'", updateUnit.ResourceData.OdataType)
		}

		siteID, listID, err = parseSharepointResource(updateUnit.Resource)
		if err != nil {
			return "", "", fmt.Errorf("invalid resource format: %s", err)
		}

	}
	return siteID, listID, nil
}

// parseResourceString ensures models.SharepointResourceFormat signature and returns siteID and listID.
func parseSharepointResource(resource string) (string, string, error) {
	parts := strings.Split(resource, "/")
	if len(parts) < 4 || parts[0] != "sites" || parts[2] != "lists" {
		return "", "", fmt.Errorf("expected resource format: '%s', expecterd signature: '%s'", resource, models.SharepointResourceSignature)
	}
	return parts[1], parts[3], nil
}
