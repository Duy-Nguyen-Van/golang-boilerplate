package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang-boilerplate/internal/config"
	"golang-boilerplate/internal/errors"

	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseManager handles database connections with advanced features
type DatabaseManager struct {
	db           *gorm.DB
	config       *config.Config
	healthStatus HealthStatus
	mu           sync.RWMutex
	metrics      *ConnectionMetrics
}

// HealthStatus represents the health status of the database connection
type HealthStatus struct {
	IsHealthy    bool          `json:"is_healthy"`
	LastCheck    time.Time     `json:"last_check"`
	LastError    string        `json:"last_error,omitempty"`
	ResponseTime time.Duration `json:"response_time"`
	RetryCount   int           `json:"retry_count"`
}

// ConnectionMetrics tracks database connection metrics
type ConnectionMetrics struct {
	TotalConnections   int           `json:"total_connections"`
	OpenConnections    int           `json:"open_connections"`
	IdleConnections    int           `json:"idle_connections"`
	InUseConnections   int           `json:"in_use_connections"`
	WaitCount          int64         `json:"wait_count"`
	WaitDuration       time.Duration `json:"wait_duration"`
	MaxOpenConnections int           `json:"max_open_connections"`
	MaxIdleConnections int           `json:"max_idle_connections"`
}

// NewDatabaseManager creates a new database manager with advanced connection management
func NewDatabaseManager(cfg *config.Config) (*DatabaseManager, error) {
	manager := &DatabaseManager{
		config:  cfg,
		metrics: &ConnectionMetrics{},
		healthStatus: HealthStatus{
			IsHealthy: false,
		},
	}

	// Initialize database connection with retry logic
	if err := manager.connectWithRetry(); err != nil {
		return nil, errors.DatabaseError("Failed to initialize database manager", err).
			WithOperation("initialize_database_manager").
			WithResource("database")
	}

	// Start health check routine
	go manager.startHealthCheck()

	// Start metrics collection routine
	go manager.startMetricsCollection()

	return manager, nil
}

// connectWithRetry attempts to connect to the database with retry logic
func (dm *DatabaseManager) connectWithRetry() error {
	var lastErr error

	for attempt := 1; attempt <= dm.config.DatabaseRetryAttempts; attempt++ {
		log.WithFields(log.Fields{
			"attempt":      attempt,
			"max_attempts": dm.config.DatabaseRetryAttempts,
		}).Info("Attempting to connect to database")

		if err := dm.connect(); err != nil {
			lastErr = err
			dm.healthStatus.RetryCount = attempt

			if hub := sentry.GetHubFromContext(context.Background()); hub != nil {
				hub.WithScope(func(scope *sentry.Scope) {
					scope.SetTag("operation", "database_connection")
					scope.SetTag("attempt", fmt.Sprintf("%d", attempt))
					scope.SetExtra("retry_count", attempt)
					hub.CaptureException(err)
				})
			}

			log.WithFields(log.Fields{
				"attempt": attempt,
				"error":   err.Error(),
			}).Error("Database connection attempt failed")

			if attempt < dm.config.DatabaseRetryAttempts {
				time.Sleep(dm.config.DatabaseRetryDelay)
			}
		} else {
			dm.healthStatus.IsHealthy = true
			dm.healthStatus.RetryCount = 0
			log.Info("Database connection established successfully")
			return nil
		}
	}

	return errors.DatabaseError("Failed to connect to database after retries", lastErr).
		WithOperation("connect_database").
		WithResource("database").
		WithContext("retry_attempts", dm.config.DatabaseRetryAttempts)
}

// connect establishes the actual database connection
func (dm *DatabaseManager) connect() error {
	logLevel := logger.Error
	if dm.config.IsDebugMode() {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dm.config.ConnectionString()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// Set connection pool parameters
	sqlDB.SetMaxOpenConns(dm.config.DatabaseMaxOpenConns)
	sqlDB.SetMaxIdleConns(dm.config.DatabaseMaxIdleConns)
	sqlDB.SetConnMaxLifetime(dm.config.DatabaseConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(dm.config.DatabaseConnMaxIdleTime)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), dm.config.DatabaseConnectTimeout)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return err
	}

	dm.db = db
	dm.updateMetrics()

	return nil
}

// GetDB returns the GORM database instance
func (dm *DatabaseManager) GetDB() *gorm.DB {
	return dm.db
}

// FastHealthCheck performs a quick health check without database ping
func (dm *DatabaseManager) FastHealthCheck() HealthStatus {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	// Return cached health status if it's recent (within last 10 seconds)
	if time.Since(dm.healthStatus.LastCheck) < 10*time.Second {
		return dm.healthStatus
	}

	// If health status is stale, trigger a full health check
	go dm.HealthCheck()
	return dm.healthStatus
}

// HealthCheck performs a health check on the database connection
func (dm *DatabaseManager) HealthCheck() HealthStatus {
	start := time.Now()

	// Use configurable timeout for health checks (default 5 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), dm.config.DatabaseHealthTimeout)
	defer cancel()

	// Get SQL DB without holding locks
	sqlDB, err := dm.db.DB()
	if err != nil {
		dm.updateHealthStatus(false, err.Error(), time.Since(start))
		return dm.getHealthStatus()
	}

	// Perform ping without holding locks
	if err := sqlDB.PingContext(ctx); err != nil {
		dm.updateHealthStatus(false, err.Error(), time.Since(start))
		return dm.getHealthStatus()
	}

	dm.updateHealthStatus(true, "", time.Since(start))
	return dm.getHealthStatus()
}

// getHealthStatus returns a copy of the current health status
func (dm *DatabaseManager) getHealthStatus() HealthStatus {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.healthStatus
}

// updateHealthStatus updates the health status
func (dm *DatabaseManager) updateHealthStatus(isHealthy bool, lastError string, responseTime time.Duration) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	dm.healthStatus.IsHealthy = isHealthy
	dm.healthStatus.LastCheck = time.Now()
	dm.healthStatus.LastError = lastError
	dm.healthStatus.ResponseTime = responseTime
}

// startHealthCheck starts a routine to periodically check database health
func (dm *DatabaseManager) startHealthCheck() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		dm.HealthCheck()
	}
}

// startMetricsCollection starts a routine to collect connection metrics
func (dm *DatabaseManager) startMetricsCollection() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		dm.updateMetrics()
	}
}

// updateMetrics updates the connection metrics
func (dm *DatabaseManager) updateMetrics() {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	sqlDB, err := dm.db.DB()
	if err != nil {
		log.WithError(err).Error("Failed to get SQL DB for metrics")
		return
	}

	stats := sqlDB.Stats()
	dm.metrics.TotalConnections = stats.OpenConnections
	dm.metrics.OpenConnections = stats.OpenConnections
	dm.metrics.IdleConnections = stats.Idle
	dm.metrics.InUseConnections = stats.InUse
	dm.metrics.WaitCount = stats.WaitCount
	dm.metrics.WaitDuration = stats.WaitDuration
	dm.metrics.MaxOpenConnections = dm.config.DatabaseMaxOpenConns
	dm.metrics.MaxIdleConnections = dm.config.DatabaseMaxIdleConns
}

// GetMetrics returns the current connection metrics
func (dm *DatabaseManager) GetMetrics() ConnectionMetrics {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return *dm.metrics
}

// Close gracefully closes the database connection
func (dm *DatabaseManager) Close() error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.db != nil {
		sqlDB, err := dm.db.DB()
		if err != nil {
			return errors.DatabaseError("Failed to get SQL DB for closing", err).
				WithOperation("close_database").
				WithResource("database")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := sqlDB.Close(); err != nil {
			log.WithContext(ctx).Error("Failed to close database connection")
			return errors.DatabaseError("Failed to close database connection", err).
				WithOperation("close_database").
				WithResource("database")
		}

		log.Info("Database connection closed gracefully")
	}

	return nil
}

// IsHealthy returns whether the database connection is healthy
func (dm *DatabaseManager) IsHealthy() bool {
	return dm.getHealthStatus().IsHealthy
}

// GetHealthStatus returns the detailed health status
func (dm *DatabaseManager) GetHealthStatus() HealthStatus {
	return dm.getHealthStatus()
}
