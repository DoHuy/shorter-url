package middleware

import (
	"shorter-rest-api/internal/config"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware creates a middleware for handling CORS
func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	// Parse allowed origins only once
	allowedOrigins := strings.Split(cfg.Server.AllowOrigins, ",")
	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if the request's Origin is in the allowed list
		originAllowed := false
		for _, o := range allowedOrigins {
			if o == origin {
				originAllowed = true
				break
			}
		}

		// Set headers only if origin is allowed
		if originAllowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		}

		// Handle preflight
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RegisterMiddlewares registers all middlewares
func RegisterMiddlewares(router *gin.Engine, cfg *config.Config) {
	// Use CORS middleware
	router.Use(CORSMiddleware(cfg))

	// Use logger and recovery middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

}
