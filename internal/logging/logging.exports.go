//go:build testing

// Exports internal functions for testing purposes.
// This file is only included in builds with the "testing" tag.
package logging

import (
	"log/slog"
)

func NewGoroutineLoggerHandler(handler slog.Handler) *GoroutineLoggerHandler {
	return &GoroutineLoggerHandler{handler: handler}
}

func SetLogLevel(levelVar *slog.LevelVar, logLevel string) error {
	return setLogLevel(levelVar, logLevel)
}

func GetGoroutineID() string {
	return getGoroutineID()
}

func GetCallerName(skip int) string {
	return getCallerName(skip)
}
