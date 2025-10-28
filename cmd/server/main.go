package main

import (
	"context"
	"errors"
	"fmt"
	"golang-boilerplate/docs"
	"golang-boilerplate/internal/cache"
	"golang-boilerplate/internal/config"
	"golang-boilerplate/internal/handlers"
	"golang-boilerplate/internal/httpclient"
	"golang-boilerplate/internal/integration/auth"
	"golang-boilerplate/internal/integration/email"
	"golang-boilerplate/internal/integration/payment"
	"golang-boilerplate/internal/integration/storage"
	"golang-boilerplate/internal/logger"
	"golang-boilerplate/internal/monitoring"
	"golang-boilerplate/internal/repositories"
	"golang-boilerplate/internal/services"
	"net"
	"net/http"
	"os"
	"time"

	"golang-boilerplate/cmd/server/routes"

	"golang-boilerplate/internal/db"

	"github.com/go-playground/validator/v10"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/ory/viper"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func NewHTTPServer(lc fx.Lifecycle,
	healthHandler *handlers.HealthHandler,
	userHandler *handlers.UserHandler,
	companyHandler *handlers.CompanyHandler,
	authProvider auth.AuthService,
	nrApp *newrelic.Application,
	cfg *config.Config,
	db *db.PostgresDB,
) *http.Server {
	handler := routes.Router(userHandler, companyHandler, healthHandler, authProvider, nrApp, cfg).Server.Handler

	srv := &http.Server{
		Addr: viper.GetString("APP_HTTP_SERVER"), Handler: handler,
		ReadHeaderTimeout: time.Second * time.Duration(viper.GetInt("APP_REQUEST_TIMEOUT")),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			log.Println("Starting HTTP server at", srv.Addr)
			go func() {
				err := srv.Serve(ln)
				if err != nil {
					log.Panic(err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down HTTP server...")

			// Graceful shutdown with timeout
			shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			// Shutdown HTTP server
			if err := srv.Shutdown(shutdownCtx); err != nil {
				log.Errorf("HTTP server shutdown error: %v", err)
				return err
			}

			// Close database connections
			if err := db.Close(); err != nil {
				log.Errorf("Database shutdown error: %v", err)
				return err
			}

			log.Println("Server shutdown completed")
			return nil
		},
	})

	return srv
}

// @title Golang Boilerplate API
// @version 1.0
// @description This is a backend API for Golang Boilerplate
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.basic  BasicAuth
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description Bearer Token Authentication. Use "Bearer {token}" as the value.
func main() {
	InitConfig(".env")
	// Ensure Swagger spec is registered and optionally override fields at runtime
	docs.SwaggerInfo.BasePath = "/api/v1"
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	// Initialize global logger before any middleware uses it
	logger.Init(cfg.LogLevel)
	nrApp := monitoring.InitNewRelic(*cfg)
	monitoring.InitSentry(*cfg)

	// Add Sentry hook to logrus with background context
	log.AddHook(monitoring.NewSentryHook(context.Background(), []log.Level{
		log.ErrorLevel,
		log.FatalLevel,
		log.PanicLevel,
	}))

	// Ensure all events are flushed before the program exits
	defer monitoring.FlushSentry()

	// Set application timezone from environment variable, default to UTC if not specified
	timezone := cfg.Timezone
	if timezone == "" {
		timezone = "UTC"
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Warnf("Invalid timezone %s, falling back to UTC", timezone)
		loc = time.UTC
	}
	time.Local = loc

	fx.New(
		fx.Supply(cfg),
		fx.Supply(nrApp),
		fx.Provide(
			NewHTTPServer,
			ProvideGormPostgres,
			ProvideValidator,
			httpclient.ProvideRestClient,
			auth.ProvideAuth,
			cache.ProvideCache,
			email.ProvideEmailSender,
			payment.ProvidePaymentAdapter,
			storage.ProvideStorageAdapter,
			repositories.ProvideUserRepository,
			repositories.ProvideCompanyRepository,
			services.ProvideCompanyService,
			services.ProvideEmailService,
			services.ProvideUserService,
			services.ProvideAuthService,
			handlers.ProvideHealthHandler,
			handlers.ProvideUserHandler,
			handlers.ProvideCompanyHandler,
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}

func ProvideValidator() *validator.Validate {
	return validator.New()
}

func InitConfig(path string) {
	viper.AutomaticEnv()
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		log.Warnf("no config file '%s' not found. Using default values", path)
	} else if err != nil { // Handle other errors that occurred while reading the config file
		panic(fmt.Errorf("fatal error while reading the config file: %w", err))
	}
}

func ProvideGormPostgres(cfg *config.Config) *db.PostgresDB {
	appDB := &db.PostgresDB{}
	err := appDB.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Connecting to Database: %v", err)
	}
	return appDB
}
