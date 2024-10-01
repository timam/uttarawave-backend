package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/timam/uttaracloud-finance-backend/handlers"
	"github.com/timam/uttaracloud-finance-backend/middlewares"
	"github.com/timam/uttaracloud-finance-backend/pkg/logger"
)

func InitRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
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
