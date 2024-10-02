package main

import (
	"github.com/timam/uttaracloud-finance-backend/cmd/server"
	"github.com/timam/uttaracloud-finance-backend/internals/configs"
	"github.com/timam/uttaracloud-finance-backend/internals/packages"
	"github.com/timam/uttaracloud-finance-backend/pkg/db"
	"github.com/timam/uttaracloud-finance-backend/pkg/logger"
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

	err = packages.InitializePackages()
	if err != nil {
		logger.Error("Packages initialization failed", zap.Error(err))
	}
	logger.Info("Package initialized successfully")

	err = db.InitializeDynamoDB()
	if err != nil {
		logger.Fatal("DynamoDB initialization failed", zap.Error(err))
	}
	logger.Info("DynamoDB initialized successfully")
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
