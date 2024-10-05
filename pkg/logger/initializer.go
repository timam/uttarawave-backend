package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func InitializeLogger() error {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.OutputPaths = []string{"stdout"}

	if viper.GetBool("server.debug") {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		config.Development = true
	} else {
		config.DisableStacktrace = true
	}

	logger, err := config.Build(zap.AddCallerSkip(1))
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
