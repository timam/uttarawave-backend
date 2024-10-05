package logger

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"time"
)

const (
	Silent logger.LogLevel = iota + 1
)

type GormLogger struct {
	ZapLogger *zap.Logger
	LogLevel  logger.LogLevel
}

func NewGormLogger() *GormLogger {
	return &GormLogger{
		ZapLogger: GetLogger(),
		LogLevel:  logger.Info,
	}
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.ZapLogger.Sugar().Infof(msg, data...)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.ZapLogger.Sugar().Warnf(msg, data...)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.ZapLogger.Sugar().Errorf(msg, data...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()

	span := trace.SpanFromContext(ctx)
	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()

	fields := []zap.Field{
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
		zap.String("sql", sql),
		zap.String("traceID", traceID),
		zap.String("spanID", spanID),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		l.ZapLogger.Error("GORM Query", fields...)
	} else {
		l.ZapLogger.Debug("GORM Query", fields...)
	}
}
