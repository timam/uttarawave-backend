package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func PackagesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Packages Endpoint",
	})

}
