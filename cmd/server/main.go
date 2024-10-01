package server

import (
	"github.com/timam/uttaracloud-finance-backend/pkg/logger"
	"github.com/timam/uttaracloud-finance-backend/routers"
	"go.uber.org/zap"
)

func StartServer() {

	router := routers.InitRouter()
	err := router.Run(":8080")
	if err != nil {
		logger.Error("Failed to start server: %v", zap.Error(err))
	}

	logger.Info("Server has been started")

}
