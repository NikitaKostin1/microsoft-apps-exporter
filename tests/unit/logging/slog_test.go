//go:build testing && unit

package logging_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"microsoft-apps-exporter/internal/logging"
	"os"
	"strings"
	"testing"
)

// TestNewLogger ensures the logger uses JSONHandler and is correctly initialized.
func TestNewLogger(t *testing.T) {
	logger := logging.NewLogger()
	if logger == nil {
		t.Fatal("Expected NewLogger to return a non-nil logger")
	}

	// Check if the logger handler is JSONHandler wrapped in GoroutineLoggerHandler
	handler := logger.Handler()
	if _, ok := handler.(*logging.GoroutineLoggerHandler); !ok {
		t.Fatalf("Expected handler to be GoroutineLoggerHandler, got %T", handler)
	}
}

// TestGoroutineLoggerHandler_HandlesLogRecord verifies the custom handler enriches logs.
func TestGoroutineLoggerHandler_HandlesLogRecord(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{})
	handler := logging.NewGoroutineLoggerHandler(baseHandler)
	logger := slog.New(handler)

	logger.Info("Test message")

	var logOutput map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logOutput); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	if logOutput["thread"] == "" {
		t.Error("Expected thread attribute, got empty value")
	}

	if logOutput["caller"] == "" {
		t.Error("Expected caller attribute, got empty value")
	}
}

// TestGoroutineLoggerHandler_Enabled verifies log level filtering.
func TestGoroutineLoggerHandler_Enabled(t *testing.T) {
	baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})
	handler := logging.NewGoroutineLoggerHandler(baseHandler)

	if !handler.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("Expected handler to be enabled for INFO level")
	}
}

// TestGoroutineLoggerHandler_WithAttrs ensures attributes are passed through correctly.
func TestGoroutineLoggerHandler_WithAttrs(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{})
	handler := logging.NewGoroutineLoggerHandler(baseHandler).WithAttrs([]slog.Attr{
		slog.String("key", "value"),
	})
	logger := slog.New(handler)
	logger.Info("Test with attrs")

	if !strings.Contains(buf.String(), "\"key\":\"value\"") {
		t.Error("Expected log output to contain added attributes")
	}
}

// TestGoroutineLoggerHandler_WithGroup verifies grouping behavior.
func TestGoroutineLoggerHandler_WithGroup(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{})
	handler := logging.NewGoroutineLoggerHandler(baseHandler).WithGroup("group1")
	logger := slog.New(handler)
	logger.Info("Test with group", slog.String("nestedKey", "nestedValue"))

	if !strings.Contains(buf.String(), "\"group1\":{\"nestedKey\":\"nestedValue\"") {
		t.Error("Expected grouped attributes in log output")
	}
}

// TestLogAttributes test log attributes include thread ID and caller name.
func TestLogAttributes(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{})
	handler := logging.NewGoroutineLoggerHandler(baseHandler)

	record := slog.Record{}
	handler.Handle(context.Background(), record)

	var logOutput map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logOutput); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	if _, exists := logOutput["thread"]; !exists {
		t.Fatal("Expected 'thread' attribute in log output")
	}
	if _, exists := logOutput["caller"]; !exists {
		t.Fatal("Expected 'caller' attribute in log output")
	}
}

// TestGetGoroutineID ensures the function extracts the goroutine ID correctly.
func TestGetGoroutineID(t *testing.T) {
	id := logging.GetGoroutineID()
	if id == "unknown" {
		t.Fatal("Expected a valid goroutine ID, got 'unknown'")
	}
}

// TestGetCallerName ensures the function extracts the correct caller function name.
func TestGetCallerName(t *testing.T) {
	name := logging.GetCallerName(2)

	if !strings.Contains(name, "TestGetCallerName") {
		t.Fatalf("Expected function name containing 'TestGetCallerName', got '%s'", name)
	}
}
