package tracing

import (
	"github.com/timam/uttarawave-finance-backend/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.uber.org/zap"
	"time"
)

func InitializeTracing() error {
	var exporter *jaeger.Exporter
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		exporter, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
		if err == nil {
			break
		}
		logger.Warn("Failed to create Jaeger exporter, retrying...",
			zap.Error(err),
			zap.Int("attempt", i+1),
			zap.String("endpoint", "http://localhost:14268/api/traces"))
		time.Sleep(time.Second * 2)
	}
	if err != nil {
		logger.Error("Failed to create Jaeger exporter after multiple attempts",
			zap.Error(err),
			zap.String("endpoint", "http://localhost:14268/api/traces"))
		return nil // Return nil instead of error to continue application startup
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("uttarawave-finance-backend"),
		)),
	)

	otel.SetTracerProvider(tp)
	logger.Info("Tracing initialized successfully")
	return nil
}
