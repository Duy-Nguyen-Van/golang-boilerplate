package errors

import (
	"fmt"
	"runtime/debug"
	"strings"

	"golang-boilerplate/internal/config"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// RecoveryMiddleware provides panic recovery with proper error handling
func RecoveryMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}

					// Create a panic error
					panicErr := InternalError("Application panic occurred", err).
						WithOperation("panic_recovery")

					// Only expose panic details to client in development
					if cfg != nil && cfg.AppEnv == "development" {
						panicErr = panicErr.
							WithContext("panic_value", r).
							WithContext("stack_trace", string(debug.Stack()))
					} else {
						// Always capture stack trace internally for logs/Sentry, but don't expose in response
						panicErr.StackTrace = string(debug.Stack())
					}

					// Log the panic
					log.WithFields(log.Fields{
						"error_code":  panicErr.Code,
						"error_type":  panicErr.Type,
						"http_status": panicErr.HTTPStatus,
						"operation":   panicErr.Operation,
						"path":        c.Request().URL.Path,
						"method":      c.Request().Method,
						"query":       c.QueryParams(),
						"user_agent":  c.Request().UserAgent(),
						"ip":          c.RealIP(),
						"stack_trace": panicErr.StackTrace,
						"panic_value": r,
					}).Error("Application panic recovered")

					// Report to Sentry
					if hub := sentry.GetHubFromContext(c.Request().Context()); hub != nil {
						hub.WithScope(func(scope *sentry.Scope) {
							scope.SetTag("error_code", panicErr.Code)
							scope.SetTag("error_type", string(panicErr.Type))
							scope.SetTag("operation", panicErr.Operation)
							scope.SetExtra("panic_value", r)
							scope.SetExtra("stack_trace", panicErr.StackTrace)
							scope.SetExtra("path", c.Request().URL.Path)
							scope.SetExtra("method", c.Request().Method)
							scope.SetExtra("query", c.QueryParams())
							scope.SetExtra("user_agent", c.Request().UserAgent())
							scope.SetExtra("ip", c.RealIP())
							hub.CaptureException(panicErr)
						})
					}

					// Return error response
					errorHandler := NewErrorHandler()
					errorHandler.errorResponse(c, panicErr)
				}
			}()

			return next(c)
		}
	}
}

// ErrorMiddleware provides centralized error handling
func ErrorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				// Skip error handling for swagger and favicon routes
				path := c.Request().URL.Path
				if isSwaggerRoute(path) || isFaviconRoute(path) {
					return err
				}

				errorHandler := NewErrorHandler()
				return errorHandler.HandleError(c, err)
			}
			return nil
		}
	}
}

// isSwaggerRoute checks if the path is a swagger-related route
func isSwaggerRoute(path string) bool {
	return strings.Contains(path, "/swagger") || strings.Contains(path, "/docs")
}

// isFaviconRoute checks if the path is a favicon request
func isFaviconRoute(path string) bool {
	return path == "/favicon.ico" || strings.Contains(path, "favicon")
}
