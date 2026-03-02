package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"urlshortener/internal/logging"
)

// GinLogger is a simple middleware that logs incoming HTTP requests
// with method, path, status, and latency using the application's logger.
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// determine path for logging
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		method := c.Request.Method
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()

		// Log with structured context
		if logging.AppLogger != nil {
			logging.AppLogger.Info("HTTP request completed",
				"method", method,
				"path", path,
				"status", status,
				"latency_ms", latency.Milliseconds(),
				"client_ip", c.ClientIP(),
				"user_agent", c.Request.UserAgent(),
			)
		}
	}
}
