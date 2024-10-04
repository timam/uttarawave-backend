package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/timam/uttarawave-finance-backend/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
)

func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Extract tracing info from the incoming request
		propagator := otel.GetTextMapPropagator()
		ctx = propagator.Extract(ctx, propagation.HeaderCarrier(c.Request.Header))

		tracer := otel.Tracer("uttarawave-finance-backend") //TODO: get value from env variable or config
		ctx, span := tracer.Start(ctx, c.FullPath())
		defer span.End()

		// Add trace ID to the response headers
		c.Header("X-Trace-ID", span.SpanContext().TraceID().String())

		// Create a request-scoped logger with trace information
		requestLogger := logger.GetLogger().With(
			zap.String("traceID", span.SpanContext().TraceID().String()),
			zap.String("spanID", span.SpanContext().SpanID().String()),
		)

		// Store the logger and span in the context
		c.Set("logger", requestLogger)
		c.Set("span", span)

		// Use the new context for the rest of the request
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
