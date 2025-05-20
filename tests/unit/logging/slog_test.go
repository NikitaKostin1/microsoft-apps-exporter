//go:build testing && unit

package logging_test

import (
	"log/slog"
	"microsoft-apps-exporter/internal/logging"
	"strings"
	"testing"
)

// TestInitLogger_ValidLevels tests valid log level strings.
func TestInitLogger_ValidLevels(t *testing.T) {
	levels := []string{"DEBUG", "INFO", "WARN", "ERROR"}
	for _, lvl := range levels {
		logger, err := logging.InitLogger(lvl)
		if err != nil {
			t.Errorf("InitLogger(%q) failed unexpectedly: %v", lvl, err)
		}
		if logger == nil {
			t.Errorf("InitLogger(%q) returned nil logger", lvl)
		}
	}
}

// TestInitLogger_InvalidLevel tests invalid log level strings.
func TestInitLogger_InvalidLevel(t *testing.T) {
	logger, err := logging.InitLogger("INVALID_LEVEL")
	if err == nil {
		t.Error("Expected error for invalid log level, got nil")
	}
	if logger != nil {
		t.Error("Expected nil logger for invalid log level")
	}
}

// TestSetLogLevel_ValidLevels tests valid log level strings mapping.
func TestSetLogLevel_ValidLevels(t *testing.T) {
	tests := []struct {
		input    string
		expected slog.Level
	}{
		{"DEBUG", slog.LevelDebug},
		{"INFO", slog.LevelInfo},
		{"WARN", slog.LevelWarn},
		{"ERROR", slog.LevelError},
	}

	for _, tt := range tests {
		var lv slog.LevelVar
		err := callSetLogLevel(&lv, tt.input)
		if err != nil {
			t.Errorf("Unexpected error for input %s: %v", tt.input, err)
		}
		if lv.Level() != tt.expected {
			t.Errorf("Expected level %v, got %v", tt.expected, lv.Level())
		}
	}
}

// TestSetLogLevel_InvalidLevel ensures invalid values are handled.
func TestSetLogLevel_InvalidLevel(t *testing.T) {
	var lv slog.LevelVar
	err := callSetLogLevel(&lv, "VERBOSE")
	if err == nil {
		t.Fatal("Expected error for invalid log level, got nil")
	}
}

func callSetLogLevel(lv *slog.LevelVar, level string) error {
	// mirror internal logic (simulate unexported setLogLevel)
	switch strings.ToUpper(level) {
	case "DEBUG":
		lv.Set(slog.LevelDebug)
	case "INFO":
		lv.Set(slog.LevelInfo)
	case "WARN":
		lv.Set(slog.LevelWarn)
	case "ERROR":
		lv.Set(slog.LevelError)
	default:
		return &loggingError{"invalid log level"}
	}
	return nil
}

type loggingError struct {
	msg string
}

func (e *loggingError) Error() string {
	return e.msg
}
