package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/timam/uttarawave-backend/internals/packages"
	"github.com/timam/uttarawave-backend/models"
	"net/http"
)

type InternetPackagesResponse struct {
	Version  string                   `json:"version"`
	Packages []models.InternetPackage `json:"packages"`
}

type CableTVPackagesResponse struct {
	Version  string                  `json:"version"`
	Packages []models.CableTVPackage `json:"packages"`
}

func CurrentInternetPackagesHandler(c *gin.Context) {
	latestFile := packages.LatestInternetPackagesFile
	internetPackages := packages.CurrentInternetPackages

	resp := InternetPackagesResponse{
		Version:  latestFile,
		Packages: internetPackages,
	}

	c.JSON(http.StatusOK, resp)
}

func CurrentCableTVPackagesHandler(c *gin.Context) {
	latestFile := packages.LatestCableTVPackagesFile
	cableTVPackages := packages.CurrentCableTVPackages

	resp := CableTVPackagesResponse{
		Version:  latestFile,
		Packages: cableTVPackages,
	}

	c.JSON(http.StatusOK, resp)
}
