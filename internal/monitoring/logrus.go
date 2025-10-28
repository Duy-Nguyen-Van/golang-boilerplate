package monitoring

import (
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/ory/viper"
	"github.com/sirupsen/logrus"
)

type NRHook struct {
	app *newrelic.Application
}

func NewNRHook(app *newrelic.Application) *NRHook {
	return &NRHook{app: app}
}

func (hook *NRHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *NRHook) Fire(entry *logrus.Entry) error {
	// Map logrus levels to NewRelic severity levels
	severity := "INFO"
	switch entry.Level {
	case logrus.DebugLevel:
		severity = "DEBUG"
	case logrus.InfoLevel:
		severity = "INFO"
	case logrus.WarnLevel:
		severity = "WARN"
	case logrus.ErrorLevel:
		severity = "ERROR"
	case logrus.FatalLevel:
		severity = "CRITICAL"
	case logrus.PanicLevel:
		severity = "CRITICAL"
	}

	// Create a copy of the entry data to avoid modifying the original
	data := make(map[string]interface{})
	for k, v := range entry.Data {
		data[k] = v
	}

	// Add standard fields
	data["app"] = viper.GetString("NEWRELIC_APP_NAME")
	data["level"] = entry.Level.String()

	// Convert latency to microseconds if present
	if latency, ok := data["latency"].(time.Duration); ok {
		data["latency"] = latency.Microseconds()
	}

	// Record the log to NewRelic
	hook.app.RecordLog(newrelic.LogData{
		Timestamp:  entry.Time.UnixMilli(),
		Severity:   severity,
		Message:    entry.Message,
		Attributes: data,
	})

	return nil
}
