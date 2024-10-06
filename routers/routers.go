package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/timam/uttarawave-backend/handlers"
	"github.com/timam/uttarawave-backend/middlewares"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"github.com/timam/uttarawave-backend/pkg/metrics"
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

	// Always apply the middlewares
	router.Use(middlewares.TracingLoggerMiddleware())
	router.Use(middlewares.MetricsMiddleware())

	logger.Info("Initializing router")
	router.GET("/metrics", metrics.MetricsHandler())

	apiV1 := router.Group("/api/v1")
	packageRoutes := apiV1.Group("/packages")
	{
		packageRoutes.GET("/internet", handlers.CurrentInternetPackagesHandler)
		packageRoutes.GET("/cabletv", handlers.CurrentCableTVPackagesHandler)
	}

	buildingRoutes := apiV1.Group("/buildings")
	{
		buildingHandler := handlers.NewBuildingHandler()
		buildingRoutes.POST("", buildingHandler.AddBuilding())
		buildingRoutes.PATCH("/:id", buildingHandler.UpdateBuilding())
		buildingRoutes.DELETE("/:id", buildingHandler.DeleteBuilding())
	}

	customerHandler := handlers.NewCustomerHandler()
	customerRoutes := apiV1.Group("/customers")
	{
		customerRoutes.POST("/", customerHandler.CreateCustomer())
		customerRoutes.GET("/", customerHandler.GetCustomer())
		customerRoutes.PUT("/", customerHandler.UpdateCustomer())
		customerRoutes.DELETE("/", customerHandler.DeleteCustomer())
	}

	logger.Info("Router initialized successfully")
	return router
}
