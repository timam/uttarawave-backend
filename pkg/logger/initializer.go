package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	"syscall"
)

var globalLogger *zap.Logger

func InitializeLogger() error {
	var err error
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	globalLogger, err = config.Build()
	if err != nil {
		return err
	}

	// Ensure logger synchronization on exit
	go handleLoggerSync()

	// Set Gin's default writer to use zap logger
	gin.DefaultWriter = newZapWriter(globalLogger, zapcore.InfoLevel)
	gin.DefaultErrorWriter = newZapWriter(globalLogger, zapcore.ErrorLevel)

	return nil
}

func GetLogger() *zap.Logger {
	return globalLogger
}

func handleLoggerSync() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	_ = globalLogger.Sync()
}
