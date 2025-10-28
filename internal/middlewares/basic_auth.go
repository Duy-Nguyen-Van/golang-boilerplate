package middlewares

import (
	"crypto/subtle"
	"golang-boilerplate/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func BasicAuthMiddleware(cfg config.Config) echo.MiddlewareFunc {
	return middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if subtle.ConstantTimeCompare([]byte(username), []byte(cfg.BasicAuthUsername)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(cfg.BasicAuthPassword)) == 1 {
			return true, nil
		}
		return false, nil
	})
}
