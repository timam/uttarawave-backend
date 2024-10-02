package middlewares

import (
	"bytes"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/timam/uttaracloud-finance-backend/pkg/logger"
	"go.uber.org/zap"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !viper.GetBool("server.debug") {
			c.Next()
			return
		}

		start := time.Now()

		// Create a custom response writer to capture response body
		w := &responseWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		// Process request
		c.Next()

		// Calculate resolution time
		duration := time.Since(start)

		// Ensure we capture the status
		status := c.Writer.Status()

		// Log details of the request and response
		logger.Info("Request processed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", status),
			zap.Duration("duration", duration),
			zap.String("clientIP", c.ClientIP()),
			zap.String("response", w.body.String()),
		)
	}
}
