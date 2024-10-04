package tracing

import (
	"github.com/timam/uttarawave-finance-backend/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"
)

func InitializeTracing() error {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces"))) //TODO: get value from env variable or config
	if err != nil {
		logger.Error("Failed to create Jaeger exporter", zap.Error(err))
		return err
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
