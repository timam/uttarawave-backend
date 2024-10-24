package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	handlers2 "github.com/timam/uttarawave-backend/api/handlers"
	middlewares2 "github.com/timam/uttarawave-backend/api/middlewares"
	repositories2 "github.com/timam/uttarawave-backend/internals/repositories"
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
	router.Use(middlewares2.TracingLoggerMiddleware())
	router.Use(middlewares2.MetricsMiddleware())

	logger.Info("Initializing router")
	router.GET("/metrics", metrics.MetricsHandler())

	apiV1 := router.Group("/api/v1")

	packageRepo := repositories2.NewGormPackageRepository()
	packageHandler := handlers2.NewPackageHandler(packageRepo)
	packageRoutes := apiV1.Group("/packages")
	{
		packageRoutes.POST("", packageHandler.CreatePackage())
		packageRoutes.GET("", packageHandler.GetAllPackages())
		packageRoutes.GET("/:id", packageHandler.GetPackageByID())
		packageRoutes.DELETE("/:id", packageHandler.DeletePackage())
	}

	buildingRepo := repositories2.NewGormBuildingRepository()
	buildingRoutes := apiV1.Group("/buildings")
	{
		buildingHandler := handlers2.NewBuildingHandler()
		buildingRoutes.POST("", buildingHandler.AddBuilding())
		buildingRoutes.GET("", buildingHandler.GetAllBuildings())
		buildingRoutes.GET("/:id", buildingHandler.GetBuilding())
		buildingRoutes.PATCH("/:id", buildingHandler.UpdateBuilding())
		buildingRoutes.DELETE("/:id", buildingHandler.DeleteBuilding())
	}

	deviceRepo := repositories2.NewGormDeviceRepository()
	deviceHandler := handlers2.NewDeviceHandler(deviceRepo)
	deviceRoutes := apiV1.Group("/devices")
	{
		deviceRoutes.POST("", deviceHandler.CreateDevice())
		deviceRoutes.GET("", deviceHandler.GetAllDevices())
		deviceRoutes.GET("/:id", deviceHandler.GetDevice())
		deviceRoutes.PUT("/:id", deviceHandler.UpdateDevice())
		deviceRoutes.POST("/:id/assign", deviceHandler.AssignDevice())
		deviceRoutes.POST("/:id/unassign", deviceHandler.UnassignDevice())
		deviceRoutes.DELETE("/:id", deviceHandler.DeleteDevice())
		deviceRoutes.GET("/by-assignment", deviceHandler.GetDeviceByAssignment())
	}

	invoiceRepo := repositories2.NewGormInvoiceRepository()

	subscriptionRepo := repositories2.NewGormSubscriptionRepository()
	subscriptionHandler := handlers2.NewSubscriptionHandler(subscriptionRepo, packageRepo, deviceRepo, invoiceRepo)
	subscriptionRoutes := apiV1.Group("/subscriptions")
	{
		subscriptionRoutes.POST("", subscriptionHandler.CreateSubscription())
		subscriptionRoutes.GET("/:id", subscriptionHandler.GetSubscription())
		subscriptionRoutes.PUT("/:id", subscriptionHandler.UpdateSubscription())
		subscriptionRoutes.DELETE("/:id", subscriptionHandler.DeleteSubscription())
		subscriptionRoutes.GET("", subscriptionHandler.GetAllSubscriptions())
	}

	customerRepo := repositories2.NewGormCustomerRepository()
	customerHandler := handlers2.NewCustomerHandler(customerRepo, buildingRepo)
	customerRoutes := apiV1.Group("/customers")
	{
		customerRoutes.POST("", customerHandler.CreateCustomer())
		customerRoutes.GET("", customerHandler.GetCustomer())
		customerRoutes.PUT("/:id", customerHandler.UpdateCustomer())
		customerRoutes.DELETE("", customerHandler.DeleteCustomer())
		customerRoutes.GET("/all", customerHandler.GetAllCustomers())
	}

	paymentRepo := repositories2.NewGormPaymentRepository()
	paymentHandler := handlers2.NewPaymentHandler(paymentRepo, subscriptionRepo, invoiceRepo)
	invoiceHandler := handlers2.NewInvoiceHandler(invoiceRepo, subscriptionRepo)

	paymentRoutes := apiV1.Group("/payments")
	{
		paymentRoutes.POST("", paymentHandler.CreatePayment())
		// Add other payment routes here
	}

	invoiceRoutes := apiV1.Group("/invoices")
	{
		invoiceRoutes.POST("", invoiceHandler.CreateInvoice())
		//invoiceRoutes.GET("/:id", invoiceHandler.GetInvoice())
		//invoiceRoutes.PUT("/:id", invoiceHandler.UpdateInvoice())
		//invoiceRoutes.GET("", invoiceHandler.GetAllInvoices())
	}

	expenseRepo := repositories2.NewGormExpenseRepository()
	expenseHandler := handlers2.NewExpenseHandler(expenseRepo)

	expenseRoutes := apiV1.Group("/expenses")
	{
		expenseRoutes.POST("", expenseHandler.CreateExpense())
		// Add other expense routes here
	}

	logger.Info("Router initialized successfully")
	return router
}
