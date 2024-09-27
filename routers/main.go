package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/timam/uttaracloud-finance-backend/handlers"
)

func InitRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	apiV1 := router.Group("/api/v1")

	packageRoutes := apiV1.Group("/packages")
	{
		packageRoutes.GET("/", handlers.PackagesHandler)
	}

	return router
}
