package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"syscall"
)

var loggerInstance *zap.Logger

func InitializeLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	loggerInstance, err = config.Build()
	if err != nil {
		return nil, err
	}

	// Ensure logger synchronization on exit
	if os.Getenv("SKIP_LOGGER_SYNC") != "true" {
		defer func() {
			if err := loggerInstance.Sync(); err != nil {
				if pathErr, ok := err.(*os.PathError); ok && pathErr.Err == syscall.ENOTTY {
					loggerInstance.Warn("Failed to sync logger; inappropriate ioctl for device")
				} else {
					_, _ = os.Stderr.WriteString("Failed to sync logger: " + err.Error() + "\n")
				}
			}
		}()
	}

	return loggerInstance, nil
}

func GetLogger() *zap.Logger {
	return loggerInstance
}
