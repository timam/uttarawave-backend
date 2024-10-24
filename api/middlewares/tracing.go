package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
	"time"
)

func TracingLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		ctx := c.Request.Context()

		// Extract tracing info from the incoming request
		propagator := otel.GetTextMapPropagator()
		ctx = propagator.Extract(ctx, propagation.HeaderCarrier(c.Request.Header))

		tracer := otel.Tracer(viper.GetString("server.name"))
		ctx, span := tracer.Start(ctx, c.FullPath())
		defer span.End()

		// Add trace ID to the response headers
		traceID := span.SpanContext().TraceID().String()
		c.Header("X-Trace-ID", traceID)

		// Create a new logger with trace information for this request
		requestLogger := logger.GetLogger().With(
			zap.String("traceID", traceID),
			zap.String("spanID", span.SpanContext().SpanID().String()),
		)

		// Use a defer function to reset the logger after the request is completed
		defer logger.SetLogger(logger.GetLogger())

		// Set the request-scoped logger
		logger.SetLogger(requestLogger)

		// Log the incoming request
		logger.Info("Request received",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("clientIP", c.ClientIP()),
		)

		// Use the new context for the rest of the request
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		// Log the request completion
		duration := time.Since(start)
		logger.Info("Request completed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
		)
	}
}
