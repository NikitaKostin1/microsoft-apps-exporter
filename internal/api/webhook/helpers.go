package webhook

import (
	"net/http"
	"net/url"
)

// handleValidationToken extracts and responds with the validation token if present in the request URL.
// This is used for webhook subscription validation.
func handleValidationToken(w http.ResponseWriter, requestURL *url.URL) (bool, error) {
	validationToken := requestURL.Query().Get("validationToken")
	if validationToken == "" {
		return false, nil
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(validationToken))

	return true, err
}
