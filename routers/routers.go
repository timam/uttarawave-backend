package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/timam/uttarawave-backend/handlers"
	"github.com/timam/uttarawave-backend/middlewares"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"github.com/timam/uttarawave-backend/pkg/metrics"
	"github.com/timam/uttarawave-backend/repositories"
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

	packageRepo := repositories.NewGormPackageRepository()
	packageHandler := handlers.NewPackageHandler(packageRepo)
	packageRoutes := apiV1.Group("/packages")
	{
		internetRoutes := packageRoutes.Group("/internet")
		{
			internetRoutes.POST("", packageHandler.CreateInternetPackage())
			internetRoutes.PUT("/:id", packageHandler.UpdateInternetPackage())
			internetRoutes.DELETE("/:id", packageHandler.DeleteInternetPackage())
			internetRoutes.GET("/:id", packageHandler.GetInternetPackage())
			internetRoutes.GET("", packageHandler.GetAllInternetPackages())
		}

		cableTVRoutes := packageRoutes.Group("/cabletv")
		{
			cableTVRoutes.POST("", packageHandler.CreateCableTVPackage())
			cableTVRoutes.PUT("/:id", packageHandler.UpdateCableTVPackage())
			cableTVRoutes.DELETE("/:id", packageHandler.DeleteCableTVPackage())
			cableTVRoutes.GET("/:id", packageHandler.GetCableTVPackage())
			cableTVRoutes.GET("", packageHandler.GetAllCableTVPackages())
		}
	}

	buildingRepo := repositories.NewGormBuildingRepository()
	buildingRoutes := apiV1.Group("/buildings")
	{
		buildingHandler := handlers.NewBuildingHandler()
		buildingRoutes.POST("", buildingHandler.AddBuilding())
		buildingRoutes.GET("", buildingHandler.GetAllBuildings())
		buildingRoutes.GET("/:id", buildingHandler.GetBuilding())
		buildingRoutes.PATCH("/:id", buildingHandler.UpdateBuilding())
		buildingRoutes.DELETE("/:id", buildingHandler.DeleteBuilding())
	}

	deviceRepo := repositories.NewGormDeviceRepository()
	deviceHandler := handlers.NewDeviceHandler(deviceRepo, buildingRepo)
	deviceRoutes := apiV1.Group("/devices")
	{
		deviceRoutes.POST("", deviceHandler.CreateDevice())
		deviceRoutes.GET("/:id", deviceHandler.GetDevice())
		deviceRoutes.PUT("/:id", deviceHandler.UpdateDevice())
		deviceRoutes.DELETE("/:id", deviceHandler.DeleteDevice())
		deviceRoutes.GET("", deviceHandler.GetAllDevices())
		deviceRoutes.POST("/:id/assign-to-subscription", deviceHandler.AssignDeviceToSubscription())
		deviceRoutes.POST("/:id/assign-to-building", deviceHandler.AssignDeviceToBuilding())
		deviceRoutes.POST("/:id/unassign", deviceHandler.UnassignDevice())
		deviceRoutes.GET("/by-subscription", deviceHandler.GetDeviceBySubscriptionID())
	}

	subscriptionRepo := repositories.NewGormSubscriptionRepository()
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionRepo, packageRepo, deviceRepo)
	subscriptionRoutes := apiV1.Group("/subscriptions")
	{
		subscriptionRoutes.POST("", subscriptionHandler.CreateSubscription())
		subscriptionRoutes.GET("/:id", subscriptionHandler.GetSubscription())
		subscriptionRoutes.PUT("/:id", subscriptionHandler.UpdateSubscription())
		subscriptionRoutes.DELETE("/:id", subscriptionHandler.DeleteSubscription())
		subscriptionRoutes.GET("", subscriptionHandler.GetAllSubscriptions())
	}

	customerRepo := repositories.NewGormCustomerRepository()
	customerHandler := handlers.NewCustomerHandler(customerRepo, buildingRepo, subscriptionRepo, deviceRepo)
	customerRoutes := apiV1.Group("/customers")
	{
		customerRoutes.POST("", customerHandler.CreateCustomer())
		customerRoutes.GET("", customerHandler.GetCustomer())
		customerRoutes.PUT("/:id", customerHandler.UpdateCustomer())
		customerRoutes.DELETE("", customerHandler.DeleteCustomer())
		customerRoutes.GET("/full-details", customerHandler.GetAllCustomersFullDetails())
	}

	incomeRepo := repositories.NewGormIncomeRepository()
	expenseRepo := repositories.NewGormExpenseRepository()
	incomeHandler := handlers.NewIncomeHandler(incomeRepo, subscriptionRepo)
	expenseHandler := handlers.NewExpenseHandler(expenseRepo)

	incomeRoutes := apiV1.Group("/incomes")
	{
		incomeRoutes.POST("", incomeHandler.CreateIncome())
		// Add other income routes here
	}

	expenseRoutes := apiV1.Group("/expenses")
	{
		expenseRoutes.POST("", expenseHandler.CreateExpense())
		// Add other expense routes here
	}

	logger.Info("Router initialized successfully")
	return router
}
