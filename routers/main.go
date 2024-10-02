package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/timam/uttaracloud-finance-backend/handlers"
	"github.com/timam/uttaracloud-finance-backend/middlewares"
	"github.com/timam/uttaracloud-finance-backend/pkg/logger"
)

func InitRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	// Always apply the CustomLogger middleware
	router.Use(middlewares.CustomLogger())

	if viper.GetBool("server.debug") {
		logger.Info("Debug mode enabled, request and response logging activated")
	} else {
		logger.Info("Debug mode disabled, request and response logging deactivated")
	}

	logger.Info("Initializing router")

	apiV1 := router.Group("/api/v1")

	packageRoutes := apiV1.Group("/packages")
	{
		packageRoutes.GET("/", handlers.PackagesHandler)
	}

	logger.Info("Router initialized successfully")
	return router
}
