# Golang Boilerplate (Echo + FX)

A production-ready Go web application built on Echo, featuring clean architecture, dependency injection (Uber FX), Keycloak auth integration, Redis caching, PostgreSQL, structured logging, and observability (New Relic + Sentry).

## Features

- **Clean Architecture**: Layered modules and clear separation of concerns
- **Dependency Injection**: Uber FX for modular dependency management
- **DTO & Model Layers**: Separation between API DTOs and domain models
- **Comprehensive Error Handling**: Structured error system with context, logging, and monitoring
- **Authentication**: JWT-based authentication with Keycloak integration
- **Caching**: Redis cache provider
- **Database**: PostgreSQL with migrations (goose)
- **Email**: AWS SES integration
- **Logging**: Structured logging with Logrus
- **Observability**: New Relic APM + Sentry error tracking
- **Docker**: Dockerfile and Compose services for Postgres/Redis
- **Middleware**: Auth, CORS, logging, rate limiting, error handling
- **Health Checks**: Built-in health endpoint

## Project Structure

```text
golang-boilerplate/
├─ cmd/
│  ├─ migrations/
│  │  └─ sql/
│  │     └─ 20250911012728_init_database.sql
│  └─ server/
│     ├─ main.go                 # Application entrypoint + FX wiring
│     └─ routes/
│        └─ router.go            # Echo routes and middleware
│
├─ docs/                         # Swagger docs output (generated)
│  ├─ docs.go
│  ├─ swagger.json
│  └─ swagger.yaml
│
├─ internal/
│  ├─ cache/                     # Cache abstraction + Redis
│  │  ├─ cache.go
│  │  └─ redis.go
│  ├─ config/                    # Config loader and env bindings
│  │  └─ config.go
│  ├─ constants/                 # Error codes, pagination, providers
│  │  ├─ error_codes.go
│  │  ├─ pagination.go
│  │  └─ third_party_provider.go
│  ├─ db/                        # Database connection management
│  │  ├─ manager.go              # Database manager with connection pooling
│  │  └─ postgres.go             # Postgres connection wrapper
│  ├─ dtos/                      # API DTOs
│  │  ├─ common.go
│  │  ├─ company.go
│  │  ├─ email.go
│  │  ├─ health.go
│  │  └─ user.go
│  ├─ errors/                    # Comprehensive error handling system
│  │  ├─ app_error.go            # Custom error types and structures
│  │  ├─ handler.go              # Error handler utilities
│  │  └─ middleware.go           # Error middleware for panic recovery
│  ├─ handlers/                  # Echo handlers
│  │  ├─ base.go                 # Base handler with error handling
│  │  ├─ company.go              # Company management endpoints
│  │  ├─ health.go               # Health check endpoints
│  │  └─ user.go                 # User management endpoints
│  ├─ httpclient/                # Outbound HTTP client (Resty)
│  │  └─ resty.go
│  ├─ integration/               # External integrations
│  │  ├─ auth/
│  │  │  ├─ auth.go
│  │  │  └─ keycloak.go
│  │  └─ email/
│  │     ├─ email.go
│  │     └─ ses.go
│  ├─ logger/
│  │  └─ logger.go
│  ├─ middlewares/
│  │  ├─ auth.go
│  │  ├─ basic_auth.go
│  │  ├─ cors.go
│  │  ├─ logging.go
│  │  └─ rate_limiter.go
│  ├─ models/
│  │  ├─ auth.go
│  │  ├─ base.go
│  │  ├─ company.go
│  │  ├─ email.go
│  │  └─ user.go
│  ├─ monitoring/
│  │  ├─ logrus.go
│  │  ├─ new_relic.go
│  │  └─ sentry.go
│  ├─ repositories/
│  │  ├─ abstract.go
│  │  ├─ company.go
│  │  └─ user.go
│  ├─ services/
│  │  ├─ auth.go
│  │  ├─ company.go
│  │  ├─ email.go
│  │  └─ user.go
│  └─ utils/
│     ├─ accent.go
│     ├─ date.go
│     └─ i18n/
│        └─ translator.go
│
├─ Dockerfile
├─ docker-compose.yml
├─ go.mod
├─ go.sum
├─ Makefile
└─ README.md
```

## Quick Start

### Prerequisites

- Go 1.25+
- Docker and Docker Compose
- Make (optional)

#### Developer tools

```bash
# Database migrations
go install github.com/pressly/goose/v3/cmd/goose@latest

# Swagger generator
go install github.com/swaggo/swag/cmd/swag@latest
```

### Setup

#### 1. Clone and setup

```bash
git clone <repository-url>
cd golang-boilerplate
go mod tidy
```

#### 2. Configure environment

```bash
# Copy example and adjust as needed
cp .env.example .env
```

#### 2.1. Configure migrations environment (goose)

```bash
cp cmd/migrations/.env.example cmd/migrations/.env
```

#### 2.2. Configure server env for container (optional)

```bash
cp cmd/server/.env.example cmd/server/.env
```

#### 3. Start dependencies with Docker

```bash
make container-up
```

#### 4. Run the application locally

```bash
# Run DB migrations (configure cmd/migrations/.env as needed for goose)
make migration-up

# Start the server
make up
```

### API Endpoints

#### Public Endpoints

- `GET /v1/` - Health check
- `GET /v1/swagger/*` - Swagger UI (non-production, protected by basic auth)

#### Health Check Endpoints

- `GET /v1/health/database` - Database health status with connection metrics
- `GET /v1/health/metrics` - Comprehensive database metrics and configuration

#### Protected Endpoints (require JWT)

**User Management:**

- `POST /v1/users` - Create user
- `GET /v1/users/{id}` - Get user by ID
- `PUT /v1/users/{id}` - Update user
- `DELETE /v1/users/{id}` - Delete user
- `GET /v1/users` - Get users list
- `GET /v1/users/test-rest-client` - Demo endpoint to test outbound REST client

**Company Management:**

- `POST /v1/companies` - Create new company
- `GET /v1/companies/{id}` - Get company by ID
- `PUT /v1/companies/{id}` - Update company
- `DELETE /v1/companies/{id}` - Delete company
- `GET /v1/companies` - Get companies list

### Example API Usage

#### Create user (with JWT token)

```bash
curl -X POST http://localhost:8080/v1/users \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'
```

#### Get user by ID (with JWT token)

```bash
curl -X GET http://localhost:8080/v1/users/123 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Create company (with JWT token)

```bash
curl -X POST http://localhost:8080/v1/companies \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Acme Corp", "description": "A great company"}'
```

#### Database health check

```bash
curl -X GET http://localhost:8080/v1/health/database
```

#### Database metrics

```bash
curl -X GET http://localhost:8080/v1/health/metrics
```

## Development

### Available Make Commands

```bash
make container-up         # Start Docker services (postgres, redis)
make container-down       # Stop Docker services
make up                   # Run the server (cmd/server)
make build cmd=server service_name=main   # Build linux binary for server
make dep                  # go mod tidy
make tests                # Run tests
make lint                 # Run golangci-lint

# Migrations (goose, reads env from cmd/migrations/.env)
make migration-status
make migration-up
make migration-down
make migration-create name=add_table

# Seeds (optional, if configured)
make seed-up
make seed-down

# Swagger docs
make swagger-load
```

### Swagger

- In non-production, Swagger is available at `GET /v1/swagger/*` and protected by basic auth.
- Set `BASIC_AUTH_USER` and `BASIC_AUTH_SECRET` in `.env`.
- Regenerate docs after handler changes:

```bash
make swagger-load
```

Example request (non-production):

```bash
curl -u "$BASIC_AUTH_USER:$BASIC_AUTH_SECRET" \
  http://localhost:8080/v1/swagger/index.html
```

### Testing

```bash
# Run all tests
make tests

# Run specific package tests
go test ./...
```

## Architecture

### Dependency Injection with Uber FX

Uber FX wires the application graph (config, logger, monitoring, db, cache, repositories, services, handlers, HTTP server). See `cmd/server/main.go` for providers and lifecycle hooks.

### DTO & Model Layers

The application separates domain models and API DTOs:

- **Models** (`internal/models/`): Domain entities
- **DTOs** (`internal/dtos/`): Request/response structs

Benefits:

- **Type Safety** and **stability** across API boundaries
- **Security**: sensitive fields are not exposed via DTOs
- **Consistency**: standardized response envelope

## Error Handling System

The application features a comprehensive error handling system that provides consistent error responses, structured logging, and monitoring integration.

### Error Types

The system defines several error types for different scenarios:

- `ErrorTypeValidation` - Input validation errors
- `ErrorTypeNotFound` - Resource not found errors
- `ErrorTypeUnauthorized` - Authentication errors
- `ErrorTypeForbidden` - Authorization errors
- `ErrorTypeConflict` - Resource conflict errors
- `ErrorTypeInternal` - Internal server errors
- `ErrorTypeExternal` - External service errors
- `ErrorTypeDatabase` - Database errors
- `ErrorTypeCache` - Cache errors
- `ErrorTypeTimeout` - Timeout errors

### Error Structure

All errors follow a consistent structure with rich context:

```go
type AppError struct {
    Code       string                 // Error code from constants
    Message    string                 // Human-readable error message
    Type       ErrorType              // Error category
    HTTPStatus int                    // HTTP status code
    Cause      error                  // Underlying error
    Context    map[string]interface{} // Additional context
    Timestamp  time.Time              // When error occurred
    StackTrace string                 // Stack trace for debugging
    Operation  string                 // Operation being performed
    Resource   string                 // Resource being accessed
}
```

### Usage Examples

#### Creating Errors

```go
// Simple error creation
err := errors.ValidationError("Invalid email format", nil)

// Error with context
err := errors.DatabaseError("Failed to create user", dbErr).
    WithOperation("create_user").
    WithResource("user").
    WithContext("user_id", userID).
    WithContext("email", email)
```

#### Handler Usage

```go
func (h *UserHandler) CreateUser(c echo.Context) error {
    // Validation errors
    if err := h.validator.Struct(requestDto); err != nil {
        return h.HandleError(c, errors.ValidationError("Validation failed", err))
    }

    // Service errors
    user, err := h.userService.Create(c.Request().Context(), &requestDto)
    if err != nil {
        return h.HandleError(c, err) // Error is already wrapped in service
    }

    return h.SuccessResponse(c, "User created successfully", user, nil)
}
```

### Error Response Format

All errors are returned in a consistent format:

```json
{
  "meta": {
    "error_code": "VALIDATION_ERROR",
    "message": "Invalid email format",
    "code": 400
  },
  "data": null
}
```

### Middleware Integration

The error handling system includes middleware for:

1. **Panic Recovery** - Catches panics and converts them to structured errors
2. **Centralized Error Handling** - Processes all errors and returns consistent responses

### Monitoring and Logging

- **Structured Logging**: All errors are logged with context fields
- **Sentry Integration**: Errors are automatically reported to Sentry with context
- **Stack Traces**: Internal errors include stack traces for debugging

## Database Connection Management

The application features a comprehensive database connection management system that provides enterprise-grade reliability, monitoring, and performance.

### Database Features

- **Advanced Connection Pooling** with configurable parameters
- **Health Monitoring** with automatic checks and retry logic
- **Graceful Shutdown** handling
- **Connection Metrics** and monitoring
- **Error Handling** with structured error reporting
- **Automatic Reconnection** on failures

### Database Architecture

The system consists of:

1. **DatabaseManager** (`internal/db/manager.go`) - Core connection management, health monitoring, metrics collection, and retry logic
2. **PostgresDB** (`internal/db/postgres.go`) - Wrapper around GORM with integration to DatabaseManager
3. **Configuration** (`internal/config/config.go`) - Database connection parameters, pool settings, and timeout configurations

### Connection Pooling

The system implements advanced connection pooling with:

- **Configurable Pool Size**: Set maximum open and idle connections
- **Connection Lifecycle Management**: Automatic cleanup of old connections
- **Idle Connection Management**: Efficient handling of unused connections

```go
// Example usage
db := &db.PostgresDB{}
err := db.NewPostgresDB(cfg)
if err != nil {
    log.Fatal("Failed to connect to database:", err)
}

// Get connection metrics
metrics := db.GetMetrics()
log.Printf("Open connections: %d", metrics.OpenConnections)
```

### Database Health Monitoring

Automatic health checks every 30 seconds:

- **Connection Validation**: Ping database to verify connectivity
- **Response Time Tracking**: Monitor query response times
- **Error Tracking**: Log and report connection issues
- **Retry Logic**: Automatic reconnection on failures

```go
// Manual health check
healthStatus := db.HealthCheck()
if !healthStatus.IsHealthy {
    log.Errorf("Database unhealthy: %s", healthStatus.LastError)
}
```

### Connection Metrics

Real-time monitoring of connection statistics:

- **Pool Statistics**: Open, idle, and in-use connections
- **Wait Metrics**: Connection wait times and counts
- **Configuration**: Current pool settings

## Configuration

Set via `.env` (loaded by viper and godotenv):

- **Server**: `APP_ENV`, `APP_NAME`, `APP_VERSION`, `TIMEZONE`, `APP_HTTP_SERVER` (e.g. `:8080`)
- **Database**: `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `DATABASE_DEBUG`
- **Database Connection Pool**: `DATABASE_MAX_OPEN_CONNS` (default: 25), `DATABASE_MAX_IDLE_CONNS` (default: 5), `DATABASE_CONN_MAX_LIFETIME` (default: 5m), `DATABASE_CONN_MAX_IDLE_TIME` (default: 1m)
- **Database Timeouts**: `DATABASE_CONNECT_TIMEOUT` (default: 30s), `DATABASE_QUERY_TIMEOUT` (default: 30s)
- **Database Retry**: `DATABASE_RETRY_ATTEMPTS` (default: 3), `DATABASE_RETRY_DELAY` (default: 1s)
- **Database Health**: `DATABASE_HEALTH_TIMEOUT` (default: 5s)
- **Database SSL**: `DATABASE_SSL_MODE` (default: disable), `DATABASE_TIMEZONE` (default: UTC)
- **Cache**: `CACHE_PROVIDER` (default: redis), `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB`, `REDIS_POOL_SIZE`, `REDIS_DIAL_TIMEOUT`, `REDIS_READ_TIMEOUT`, `REDIS_WRITE_TIMEOUT`, `REDIS_POOL_TIMEOUT`, `REDIS_MAX_RETRIES`, `REDIS_MIN_RETRY_BACKOFF`, `REDIS_MAX_RETRY_BACKOFF`
- **Authentication**: `AUTH_PROVIDER`, `KEYCLOAK_URL`, `KEYCLOAK_REALM`, `KEYCLOAK_CLIENT_ID`, `KEYCLOAK_CLIENT_SECRET`, `KEY_CLAIMS`, `KEYCLOAK_REDIRECT_URI`
- **Email**: `EMAIL_PROVIDER` (ses), `AWS_SES_REGION`, `AWS_SES_ACCESS_KEY`, `AWS_SES_SECRET_KEY`
- **Rate Limiting**: `DEFAULT_RATE_LIMIT`, `AUTH_RATE_LIMIT`, `PUBLIC_RATE_LIMIT`, `RATE_LIMIT`, `RATE_LIMIT_DURATION`
- **Observability**: `NEWRELIC_APP_NAME`, `NEWRELIC_LICENSE`, `SENTRY_DSN`
- **Swagger Basic Auth**: `BASIC_AUTH_USER`, `BASIC_AUTH_SECRET`

### Database Configuration Parameters

| Parameter                     | Default | Description                        |
| ----------------------------- | ------- | ---------------------------------- |
| `DATABASE_MAX_OPEN_CONNS`     | 25      | Maximum number of open connections |
| `DATABASE_MAX_IDLE_CONNS`     | 5       | Maximum number of idle connections |
| `DATABASE_CONN_MAX_LIFETIME`  | 5m      | Maximum lifetime of a connection   |
| `DATABASE_CONN_MAX_IDLE_TIME` | 1m      | Maximum idle time of a connection  |
| `DATABASE_CONNECT_TIMEOUT`    | 30s     | Connection timeout                 |
| `DATABASE_QUERY_TIMEOUT`      | 30s     | Query timeout                      |
| `DATABASE_HEALTH_TIMEOUT`     | 5s      | Health check timeout               |
| `DATABASE_RETRY_ATTEMPTS`     | 3       | Number of retry attempts           |
| `DATABASE_RETRY_DELAY`        | 1s      | Delay between retries              |

### Rate Limiting

Default limits are configurable via env. Middleware is applied globally in `router.go`.

## Docker

### Build and Run

```bash
# Build image
docker build -t golang-boilerplate .

# Run container (ensure APP_HTTP_SERVER is set to :3000 in container env)
docker run -p 3000:3000 --env-file .env golang-boilerplate
```

### Services Included

- **PostgreSQL**: Database (5432)
- **Redis**: Cache (6379)
  - App service is commented out in `docker-compose.yml`. Run the app locally with `make up` or create your own app service.

## Database Monitoring and Troubleshooting

### Metrics to Monitor

1. **Connection Pool Utilization**

   - Open connections vs. max connections
   - Idle connections
   - Wait times

2. **Health Status**

   - Connection health
   - Response times
   - Error rates

3. **Performance Metrics**
   - Query response times
   - Connection establishment time
   - Retry attempts

### Recommended Alerts

- Connection pool utilization > 80%
- Health check failures
- Response times > 1 second
- Retry attempts > 2

### Common Issues and Solutions

1. **Connection Pool Exhaustion**

   - Increase `DATABASE_MAX_OPEN_CONNS`
   - Check for connection leaks
   - Monitor connection usage patterns

2. **Slow Queries**

   - Check `DATABASE_QUERY_TIMEOUT`
   - Monitor query performance
   - Optimize database queries

3. **Connection Failures**

   - Check network connectivity
   - Verify database credentials
   - Monitor database server health

4. **High Response Times**
   - Check connection pool settings
   - Monitor database server performance
   - Optimize query performance

### Debugging

1. **Enable Debug Logging**

   ```bash
   DATABASE_DEBUG=true
   LOG_LEVEL=debug
   ```

2. **Check Health Endpoints**

   ```bash
   curl http://localhost:8080/v1/health/database
   curl http://localhost:8080/v1/health/metrics
   ```

3. **Monitor Logs**
   - Check application logs for database errors
   - Monitor Sentry for error reports
   - Use database monitoring tools

### Performance Considerations

- **Small Applications**: 5-10 connections
- **Medium Applications**: 10-25 connections
- **Large Applications**: 25-50 connections
- **High Traffic**: 50+ connections

## Database Migrations

```bash
# Migration status
make migration-status

# Run migrations up/down (configure cmd/migrations/.env for goose)
make migration-up
make migration-down

# Create a new migration
make migration-create name=add_users
```

## Production Deployment

1. Build and push Docker image or deploy the binary built from `cmd/server`.
2. Set `ENVIRONMENT=production` and all required env vars.
3. Expose the port configured by `APP_HTTP_SERVER` (e.g. `:3000`).

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run tests and linting
6. Submit a pull request

## License

This project is licensed under the MIT License.
