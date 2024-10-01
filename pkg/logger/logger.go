package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
	"syscall"
)

var (
	loggerInstance *zap.Logger
	once           sync.Once
)

func InitializeLogger() error {
	var err error
	once.Do(func() {
		config := zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		loggerInstance, err = config.Build()
		if err != nil {
			return
		}

		// Ensure logger synchronization on exit
		if os.Getenv("SKIP_LOGGER_SYNC") != "true" {
			go func() {
				<-make(chan os.Signal, 1)
				if err := loggerInstance.Sync(); err != nil {
					if pathErr, ok := err.(*os.PathError); ok && pathErr.Err == syscall.ENOTTY {
						loggerInstance.Warn("Failed to sync logger; inappropriate ioctl for device")
					} else {
						_, _ = os.Stderr.WriteString("Failed to sync logger: " + err.Error() + "\n")
					}
				}
			}()
		}
	})
	return err
}

func getLogger() *zap.Logger {
	if loggerInstance == nil {
		panic("Logger not initialized. Call InitializeLogger() first.")
	}
	return loggerInstance
}

// Convenience functions

func Info(msg string, fields ...zap.Field) {
	getLogger().Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	getLogger().Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	getLogger().Warn(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	getLogger().Debug(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	getLogger().Fatal(msg, fields...)
}
