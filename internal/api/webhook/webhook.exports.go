//go:build testing

// Exports internal functions for testing purposes.
// This file is only included in builds with the "testing" tag.
package webhook

import (
	"microsoft-apps-exporter/internal/models"
	"microsoft-apps-exporter/internal/sync"
	"net/http"
	"net/url"
)

func NewSharepointHandler(syncer *sync.Syncer) http.HandlerFunc {
	return newSharepointHandler(syncer)
}

func NewSubscriptionHandler(syncer *sync.Syncer) http.HandlerFunc {
	return newSubscriptionHandler(syncer)
}

func HandleValidationToken(w http.ResponseWriter, requestURL *url.URL) (bool, error) {
	return handleValidationToken(w, requestURL)
}

func HandleMethodNotAllowed(w http.ResponseWriter, message string) {
	handleMethodNotAllowed(w, message)
}

func HandleBadRequest(w http.ResponseWriter, message string) {
	handleBadRequest(w, message)
}

func HandleInternalError(w http.ResponseWriter, message string) {
	handleInternalError(w, message)
}

func ExctractListReference(siteID, listID string) (models.ListReference, bool) {
	return exctractListReference(siteID, listID)
}

func ExtractResourceUpdateData(r *http.Request) (string, string, error) {
	return extractResourceUpdateData(r)
}

func ParseSharepointResource(resource string) (string, string, error) {
	return parseSharepointResource(resource)
}

func ExtractSubscriptionLifecycleData(r *http.Request) (string, error) {
	return extractSubscriptionLifecycleData(r)
}
