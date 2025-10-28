package monitoring

import (
	"context"
	"fmt"
	"golang-boilerplate/internal/config"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

// SentryHook is a logrus hook that sends logs to Sentry using the official Sentry logger
type SentryHook struct {
	levels []logrus.Level
	ctx    context.Context
}

func InitSentry(config config.Config) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              config.SentryDSN,
		Environment:      config.AppEnv,
		Release:          config.AppName + "@" + config.AppVersion,
		Debug:            config.AppEnv == "development",
		AttachStacktrace: true,
		EnableTracing:    true,
		EnableLogs:       true,
		TracesSampleRate: 1.0,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			// Filter out sensitive information
			if event.Request != nil {
				// Remove sensitive headers
				if event.Request.Headers != nil {
					delete(event.Request.Headers, "Authorization")
					delete(event.Request.Headers, "Cookie")
				}
			}
			return event
		},
	})

	if err != nil {
		log.Fatalf("sentry.Init failed: %v", err)
	}
}

// Flush buffered events before shutdown
func FlushSentry() {
	sentry.Flush(2 * time.Second)
}

// GetSentryHub returns a Sentry hub from the context if available, otherwise falls back to the current hub.
// This makes error reporting more robust by ensuring we always try to report errors when possible.
func GetSentryHub(ctx context.Context) *sentry.Hub {
	// Try to get hub from context first
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		// Fallback to current hub if context hub is not available
		hub = sentry.CurrentHub()
	}
	return hub
}

// NewSentryHook creates a new hook for logrus that sends logs to Sentry
func NewSentryHook(ctx context.Context, levels []logrus.Level) *SentryHook {
	if levels == nil {
		levels = logrus.AllLevels
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return &SentryHook{
		levels: levels,
		ctx:    ctx,
	}
}

// Levels returns the logging levels this hook will be fired for
func (h *SentryHook) Levels() []logrus.Level {
	return h.levels
}

// Fire sends the log entry to Sentry using the official Sentry logger
func (h *SentryHook) Fire(entry *logrus.Entry) error {
	// Create a Sentry logger with the current context
	sentryLogger := sentry.NewLogger(h.ctx)

	// Create a standard logger that writes to Sentry
	logger := log.New(sentryLogger, "", log.LstdFlags)

	// Configure Sentry scope with all fields
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		// Set log level
		scope.SetLevel(getSentryLevel(entry.Level))

		// Add all fields as extra data
		for k, v := range entry.Data {
			// Special handling for error field
			if k == "error" {
				if err, ok := v.(error); ok {
					scope.SetExtra("error_details", err.Error())
					continue
				}
			}
			// Add other fields
			scope.SetExtra(k, v)
		}

		// Add standard fields
		scope.SetTag("logger", "logrus")
		scope.SetTag("log_level", entry.Level.String())

		// Add timestamp
		scope.SetExtra("timestamp", entry.Time.Format(time.RFC3339))

		// Add caller information if available
		if entry.HasCaller() {
			scope.SetExtra("caller_file", entry.Caller.File)
			scope.SetExtra("caller_line", entry.Caller.Line)
			scope.SetExtra("caller_function", entry.Caller.Function)
		}
	})

	// If there's an error in the fields, capture it as an exception
	if err, ok := entry.Data["error"].(error); ok {
		hub.CaptureException(err)
		return nil
	}

	// Format the message with fields that aren't already in the scope
	msg := entry.Message
	if len(entry.Data) > 0 {
		// Add fields to the message that aren't already in the scope
		for k, v := range entry.Data {
			if k != "error" { // Skip error as it's handled separately
				msg = msg + " " + k + "=" + formatValue(v)
			}
		}
	}

	// Log the message
	logger.Println(msg)

	return nil
}

// getSentryLevel converts logrus level to sentry level
func getSentryLevel(level logrus.Level) sentry.Level {
	switch level {
	case logrus.DebugLevel:
		return sentry.LevelDebug
	case logrus.InfoLevel:
		return sentry.LevelInfo
	case logrus.WarnLevel:
		return sentry.LevelWarning
	case logrus.ErrorLevel:
		return sentry.LevelError
	case logrus.FatalLevel:
		return sentry.LevelFatal
	case logrus.PanicLevel:
		return sentry.LevelFatal
	default:
		return sentry.LevelInfo
	}
}

// formatValue formats a value for logging
func formatValue(v interface{}) string {
	switch val := v.(type) {
	case error:
		return val.Error()
	case string:
		return val
	default:
		return fmt.Sprintf("%v", val)
	}
}

// SentryWriter implements io.Writer to write logs to Sentry
type SentryWriter struct {
	ctx context.Context
}

// NewSentryWriter creates a new writer that writes to Sentry
func NewSentryWriter(ctx context.Context) *SentryWriter {
	if ctx == nil {
		ctx = context.Background()
	}
	return &SentryWriter{ctx: ctx}
}

// Write implements io.Writer
func (w *SentryWriter) Write(p []byte) (n int, err error) {
	sentryLogger := sentry.NewLogger(w.ctx)
	logger := log.New(sentryLogger, "", log.LstdFlags)
	logger.Print(string(p))
	return len(p), nil
}
