package monitoring

import (
	"golang-boilerplate/internal/config"

	"github.com/newrelic/go-agent/v3/integrations/nrlogrus"
	"github.com/newrelic/go-agent/v3/newrelic"
	log "github.com/sirupsen/logrus"
)

func InitNewRelic(config config.Config) *newrelic.Application {
	appName := config.NewRelicAppName
	license := config.NewRelicLicense

	// Skip NewRelic initialization if license is not provided
	if license == "" {
		log.Warn("NewRelic license not provided, skipping NewRelic initialization")
		return nil
	}

	// Use default app name if not provided
	if appName == "" {
		appName = "golang-boilerplate"
	}

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(appName),
		newrelic.ConfigLicense(license),
		newrelic.ConfigAppLogForwardingEnabled(true),
		newrelic.ConfigDistributedTracerEnabled(true),
		func(config *newrelic.Config) {
			log.SetLevel(log.InfoLevel)
			config.Logger = nrlogrus.StandardLogger()
		},
	)

	if err != nil {
		log.Warnf("Failed to initialize NewRelic: %v", err)
		return nil
	}

	log.AddHook(NewNRHook(app))
	log.Info("NewRelic initialized successfully")

	return app
}
