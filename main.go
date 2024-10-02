package main

import (
	"github.com/timam/uttaracloud-finance-backend/cmd/server"
	"github.com/timam/uttaracloud-finance-backend/internals/configs"
	"github.com/timam/uttaracloud-finance-backend/internals/packages"
	"github.com/timam/uttaracloud-finance-backend/pkg/logger"
	"go.uber.org/zap"
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

}

func main() {
	if err := configs.InitializeConfig(); err != nil {
		logger.Fatal("Error initializing configs", zap.Error(err))
	}
	server.StartServer()
	select {}
}
