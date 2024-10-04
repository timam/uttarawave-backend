package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func InitializeLogger() error {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		return err
	}

	SetLogger(logger)

	// Set the output for the standard logger
	log.SetOutput(NewZapErrorWriter())

	// Ensure logger synchronization on exit
	go handleLoggerSync()

	// Set Gin's default writer to use zap logger
	gin.DefaultWriter = newZapWriter(GetLogger(), zapcore.InfoLevel)
	gin.DefaultErrorWriter = newZapWriter(GetLogger(), zapcore.ErrorLevel)

	return nil
}
func handleLoggerSync() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	_ = globalLogger.Sync()
}
