package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/timam/uttaracloud-finance-backend/handlers"
	"github.com/timam/uttaracloud-finance-backend/middlewares"
	"github.com/timam/uttaracloud-finance-backend/pkg/logger"
)

func InitRouter() *gin.Engine {

	if viper.GetBool("server.debug") {
		gin.SetMode(gin.DebugMode)
		logger.Info("Debug mode enabled, setting Gin to DebugMode")
	} else {
		gin.SetMode(gin.ReleaseMode)
		logger.Info("Debug mode disabled, setting Gin to ReleaseMode")
	}

	router := gin.New()
	router.Use(gin.Recovery())

	// Always apply the CustomLogger middleware
	router.Use(middlewares.CustomLogger())

	logger.Info("Initializing router")

	apiV1 := router.Group("/api/v1")

	packageRoutes := apiV1.Group("/packages")
	{
		packageRoutes.GET("/", handlers.PackagesHandler)
	}

	logger.Info("Router initialized successfully")
	return router
}
