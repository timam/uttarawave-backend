package middlewares

import (
	"bytes"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/timam/uttaracloud-finance-backend/pkg/logger"
	"go.uber.org/zap"
)

// responseWriter is a custom wrapper for gin.ResponseWriter
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the response body
func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Create a custom response writer to capture response body
		writer := &responseWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = writer

		// Process request
		c.Next()

		// Calculate resolution time
		duration := time.Since(start)

		// Log details of the request and response
		logger.Info("Request processed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.String("clientIP", c.ClientIP()),
			zap.String("response", writer.body.String()), // Log the response body
		)
	}
}
