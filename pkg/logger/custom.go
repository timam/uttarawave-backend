package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
)

type zapWriter struct {
	logger *zap.Logger
	level  zapcore.Level
}

type zapErrorWriter struct{}

func (zw *zapWriter) Write(p []byte) (n int, err error) {
	message := string(p)
	switch zw.level {
	case zapcore.DebugLevel:
		zw.logger.Debug(message)
	case zapcore.InfoLevel:
		zw.logger.Info(message)
	case zapcore.WarnLevel:
		zw.logger.Warn(message)
	case zapcore.ErrorLevel:
		zw.logger.Error(message)
	case zapcore.DPanicLevel:
		zw.logger.DPanic(message)
	case zapcore.PanicLevel:
		zw.logger.Panic(message)
	case zapcore.FatalLevel:
		zw.logger.Fatal(message)
	}
	return len(p), nil
}

func newZapWriter(logger *zap.Logger, level zapcore.Level) io.Writer {
	return &zapWriter{
		logger: logger,
		level:  level,
	}
}

func (w zapErrorWriter) Write(p []byte) (n int, err error) {
	Error(string(p))
	return len(p), nil
}

func NewZapErrorWriter() io.Writer {
	return &zapErrorWriter{}
}
