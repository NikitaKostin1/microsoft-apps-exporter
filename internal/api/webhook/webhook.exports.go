//go:build testing

// Exports internal functions for testing purposes.
// This file is only included in builds with the "testing" tag.
package webhook

import (
	"net/http"
	"net/url"
)

func HandleValidationToken(w http.ResponseWriter, requestURL *url.URL) (bool, error) {
	return handleValidationToken(w, requestURL)
}

func HandleMethodNotAllowed(w http.ResponseWriter, method string) {
	handleMethodNotAllowed(w, method)
}

func HandleBadRequest(w http.ResponseWriter, message string) {
	handleBadRequest(w, message)
}

func HandleInternalError(w http.ResponseWriter, message string) {
	handleInternalError(w, message)
}

func ExtractSubscriptionLifecycleData(r *http.Request) (string, error) {
	return extractSubscriptionLifecycleData(r)
}
