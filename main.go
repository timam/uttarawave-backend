package main

import (
	"fmt"
	"github.com/timam/uttarawave-backend/cmd"
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

func serve() error {
	serverInstance, err := cmd.InitializeServer()
	if err != nil {
		return fmt.Errorf("server initialization failed: %w", err)
	}
	logger.Info("Server initialized successfully")

	errChan := make(chan error, 1)
	go func() {
		errChan <- serverInstance.RunServer()
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("server error: %w", err)
		}
	case sig := <-sigChan:
		logger.Warn("Received signal", zap.String("signal", sig.String()))
	}

	if err := serverInstance.ShutdownServer(); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	} else {
		logger.Info("Server exited gracefully")
	}

	return nil
}

func migrate() error {
	logger.Info("TODO: Implement database migration!!")
	return nil
}

func main() {
	cmd.Serve = serve
	cmd.Migrate = migrate

	if err := cmd.Execute(); err != nil {
		logger.Error("Failed to execute", zap.Error(err))
		os.Exit(1)
	}
}
