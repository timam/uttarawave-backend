package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	"syscall"
)

// InitializeLogger initializes the custom logger
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
				sig := make(chan os.Signal, 1)
				signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
				<-sig
				syncLogger()
			}()
		}

		// Set Gin's default writer to use zap logger
		gin.DefaultWriter = newZapWriter(loggerInstance, zapcore.InfoLevel)
		gin.DefaultErrorWriter = newZapWriter(loggerInstance, zapcore.ErrorLevel)
	})
	return err
}
