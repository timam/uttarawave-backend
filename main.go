package main

import (
	"github.com/timam/uttaracloud-finance-backend/cmd/server"
	"github.com/timam/uttaracloud-finance-backend/internals/packages"
	"github.com/timam/uttaracloud-finance-backend/pkg/logger"
	"go.uber.org/zap"
)

var log *zap.Logger

func init() {
	log, err := logger.InitializeLogger()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	err = packages.InitializePackages()
	if err != nil {
		log.Error("Initialization failed: %v", zap.Error(err))
	}
}

func main() {
	server.StartServer()
}
