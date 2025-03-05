package logging

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"strings"
)

// ConfigureSlog sets the default slog logger with custom attributes.
func ConfigureSlog() {
	slog.SetDefault(NewLogger())
}

// GoroutineLoggerHandler is a custom slog handler that injects goroutine ID and caller name into log records.
type GoroutineLoggerHandler struct {
	handler slog.Handler
}

// Handle enriches the log record with additional attributes.
func (h *GoroutineLoggerHandler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(
		slog.String("thread", getGoroutineID()),
		slog.String("caller", getCallerName(4)),
	)
	return h.handler.Handle(ctx, r)
}

// Enabled reports whether the handler handles records at the given level.
func (h *GoroutineLoggerHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// WithAttrs returns a new Handler whose attributes consist of both the receiver's attributes and the arguments.
func (h *GoroutineLoggerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &GoroutineLoggerHandler{handler: h.handler.WithAttrs(attrs)}
}

// WithGroup returns a new Handler with the given group appended to the receiver's existing groups.
func (h *GoroutineLoggerHandler) WithGroup(name string) slog.Handler {
	return &GoroutineLoggerHandler{handler: h.handler.WithGroup(name)}
}

// NewLogger initializes and returns a new slog.Logger with custom attributes.
func NewLogger() *slog.Logger {
	baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})
	customHandler := &GoroutineLoggerHandler{handler: baseHandler}
	return slog.New(customHandler)
}

// getGoroutineID returns the current goroutine ID as a string.
func getGoroutineID() string {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	fields := strings.Fields(string(buf[:n]))
	if len(fields) > 1 {
		return fields[1] // Goroutine ID is the second field in the stack trace
	}
	return "unknown"
}

// getCallerName extracts the caller name from the call stack.
func getCallerName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	fn := runtime.FuncForPC(pc)
	return fn.Name()
}
