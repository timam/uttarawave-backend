package middlewares

import (
	"bytes"
	"go.opentelemetry.io/otel/trace"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/timam/uttarawave-finance-backend/pkg/logger"
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

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Try to get the trace ID from the context
		var traceID string
		if span := trace.SpanFromContext(c.Request.Context()); span != nil {
			traceID = span.SpanContext().TraceID().String()
		}

		// Log request
		logFields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("clientIP", c.ClientIP()),
		}
		if traceID != "" {
			logFields = append(logFields, zap.String("traceID", traceID))
		}
		logger.Info("Request received", logFields...)

		c.Next()

		// Log response
		duration := time.Since(start)
		status := c.Writer.Status()

		logFields = []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", status),
			zap.Duration("duration", duration),
		}
		if traceID != "" {
			logFields = append(logFields, zap.String("traceID", traceID))
		}

		if viper.GetBool("server.debug") {
			logFields = append(logFields,
				zap.String("clientIP", c.ClientIP()),
				zap.String("response", c.Writer.(*responseWriter).body.String()),
			)
		}

		logger.Info("Request processed", logFields...)
	}
}
