package logger

import (
	"go.uber.org/zap"
	"os"
	"sync"
	"syscall"
)

var (
	globalLogger *zap.Logger
	loggerMutex  sync.RWMutex
)

func GetLogger() *zap.Logger {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()
	return globalLogger
}

func SetLogger(l *zap.Logger) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	globalLogger = l
}

func SyncLogger() {
	logger := GetLogger()
	if logger == nil {
		return
	}
	err := logger.Sync()
	if err != nil {
		if pathErr, ok := err.(*os.PathError); ok && pathErr.Err == syscall.ENOTTY {
			logger.Warn("Failed to sync logger; inappropriate ioctl for device")
		} else {
			_, _ = os.Stderr.WriteString("Failed to sync logger: " + err.Error() + "\n")
		}
	}
}

// Convenience functions
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}

func With(fields ...zap.Field) *zap.Logger {
	return GetLogger().With(fields...)
}
