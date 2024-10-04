package tracing

import (
	"context"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.uber.org/zap"
	"time"
)

func InitializeTracing() error {
	ctx := context.Background()

	exporter, err := otlptrace.New(
		ctx,
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint("localhost:4318"),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		logger.Error("Failed to create OTLP exporter", zap.Error(err))
		return err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("uttarawave-backend"),
		),
	)
	if err != nil {
		logger.Error("Failed to create resource", zap.Error(err))
		return err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	logger.Info("Tracing initialized successfully")

	// Ensure graceful shutdown
	go func() {
		<-ctx.Done()
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctxShutdown); err != nil {
			logger.Error("Error shutting down tracer provider", zap.Error(err))
		}
	}()

	return nil
}
