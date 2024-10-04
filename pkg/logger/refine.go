package logger

import (
	"go.uber.org/zap"
	"os"
	"sync"
	"syscall"
)

var (
	loggerInstance *zap.Logger
	once           sync.Once
)

func getLogger() *zap.Logger {
	if loggerInstance == nil {
		panic("Logger not initialized. Call InitializeLogger() first.")
	}
	return loggerInstance
}

func SyncLogger() {
	if loggerInstance == nil {
		return
	}
	err := loggerInstance.Sync()
	if err != nil {
		if pathErr, ok := err.(*os.PathError); ok && pathErr.Err == syscall.ENOTTY {
			loggerInstance.Warn("Failed to sync logger; inappropriate ioctl for device")
		} else {
			_, _ = os.Stderr.WriteString("Failed to sync logger: " + err.Error() + "\n")
		}
	}
}

// Info Convenience functions
func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}
