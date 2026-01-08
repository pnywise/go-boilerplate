package middewares

import (
    "encoding/base64"
    "net/http"
    "os"
    "strings"

    "github.com/gin-gonic/gin"
)

// BasicAuthMiddleware returns a middleware that enforces HTTP Basic Auth using
// credentials from BASIC_AUTH_USER and BASIC_AUTH_PASS environment variables.
//
// Behavior:
// - If both env vars are set: requests must present valid Basic auth credentials.
// - If either env var is empty: middleware is a no-op (does not block requests).
func BasicAuthMiddleware() gin.HandlerFunc {
    user := strings.TrimSpace(os.Getenv("BASIC_AUTH_USER"))
    pass := os.Getenv("BASIC_AUTH_PASS")

    // If not configured, don't enforce auth.
    if user == "" || pass == "" {
        return func(c *gin.Context) { c.Next() }
    }

    return func(c *gin.Context) {
        auth := c.GetHeader("Authorization")
        if auth == "" || !strings.HasPrefix(auth, "Basic ") {
            c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        payload := strings.TrimPrefix(auth, "Basic ")
        decoded, err := base64.StdEncoding.DecodeString(payload)
        if err != nil {
            c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        parts := strings.SplitN(string(decoded), ":", 2)
        if len(parts) != 2 || parts[0] != user || parts[1] != pass {
            c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        c.Next()
    }
}