package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/timam/uttaracloud-finance-backend/internals/packages"
	"net/http"
)

func CurrentInternetPackagesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"packages": packages.CurrentInternetPackages,
	})
}

func CurrentCableTVPackagesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"packages": packages.CurrentCableTVPackages,
	})
}
