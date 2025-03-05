package webhook

import (
	"log/slog"
	"net/http"
)

// handleMethodNotAllowed responds with a 405 Method Not Allowed status code.
func handleMethodNotAllowed(w http.ResponseWriter, method string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("405 - Method Not Allowed: Only POST method is allowed"))
	slog.Error("Only POST method is allowed", "attempted_method", method, "operation", "webhook")
}

// handleBadRequest responds with a 400 Bad Request status code and a custom message.
func handleBadRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("400 - Bad Request: " + message))
	slog.Error("BadRequest respond with message", "message", message, "operation", "webhook")
}

// handleInternalError responds with a 500 Internal Server Error status code and a custom message.
func handleInternalError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 - Internal Server Error: " + message))
	slog.Error("InternalError respond with message", "message", message, "operation", "webhook")
}
