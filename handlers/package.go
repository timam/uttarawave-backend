package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/timam/uttaracloud-finance-backend/storage"
	"net/http"
)

func PackagesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"packages": storage.InternetPackages,
	})
}
