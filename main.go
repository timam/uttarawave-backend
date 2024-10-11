package main

import (
	"github.com/timam/uttarawave-backend/cmd/server"
	"github.com/timam/uttarawave-backend/internals/configs"
	"github.com/timam/uttarawave-backend/pkg/db"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"github.com/timam/uttarawave-backend/pkg/metrics"
	"github.com/timam/uttarawave-backend/pkg/tracing"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	err := logger.InitializeLogger()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	logger.Info("Logger initialized successfully")

	err = configs.InitializeConfig()
	if err != nil {
		logger.Fatal("Config initialization failed", zap.Error(err))
	}
	logger.Info("Config initialized successfully")

	if err := db.InitializePostgreSQL(); err != nil {
		logger.Fatal("Failed to initialize PostgreSQL", zap.Error(err))
	}
	logger.Info("PostgreSQL initialized successfully")

	err = metrics.InitializeMetrics()
	if err != nil {
		logger.Error("Metrics initialization failed", zap.Error(err))
	}
	logger.Info("Metrics initialized successfully")

	err = tracing.InitializeTracing()
	if err != nil {
		logger.Error("Tracing initialization failed", zap.Error(err))
	}
	logger.Info("Tracing initialized successfully")
}

func main() {
	// Start server in a goroutine
	go server.StartServer()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for interrupt signal
	sig := <-sigChan
	logger.Warn("Received signal", zap.String("signal", sig.String()))

	// Shutdown server
	if err := server.ShutdownServer(); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	} else {
		logger.Info("Server exited gracefully")
	}

	// Sync logger
	logger.SyncLogger()

	// Exit
	os.Exit(0)
}
