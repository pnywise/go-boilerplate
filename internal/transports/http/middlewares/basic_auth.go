package middlewares

import (
	"go-boilerplate/internal/configs"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
)

// BasicAuthMiddleware returns a middleware that enforces HTTP Basic Auth using
// credentials from the provided Config.
// If either credential in cfg is empty, the middleware is a no-op.
func BasicAuthMiddleware(cfg configs.Config) echo.MiddlewareFunc {
	user := strings.TrimSpace(cfg.BasicAuthUser)
	pass := cfg.BasicAuthPass
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if user == "" || pass == "" {
				return next(c)
			}

			u, p, ok := c.Request().BasicAuth()
			if !ok || u != user || p != pass {
				return c.String(http.StatusUnauthorized, "Unauthorized")
			}
			return next(c)
		}
	}
}
