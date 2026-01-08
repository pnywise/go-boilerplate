package middlewares

import (
	"net/http"
	"strings"
	"go-boilerplate/internal/configs"

	"github.com/gin-gonic/gin"
)
// BasicAuthMiddleware returns a middleware that enforces HTTP Basic Auth using
// credentials from the provided Config.
// If either credential in cfg is empty, the middleware is a no-op.
func BasicAuthMiddleware(cfg configs.Config) gin.HandlerFunc {
    user := strings.TrimSpace(cfg.BasicAuthUser)
    pass := cfg.BasicAuthPass

    if user == "" || pass == "" {
        return func(c *gin.Context) { c.Next() }
    }

    return func(c *gin.Context) {
        username, password, ok := c.Request.BasicAuth()
        if !ok || username != user || password != pass {
            c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }
        c.Next()
    }
}