package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/timam/uttarawave-backend/api/response"
	"github.com/timam/uttarawave-backend/internals/models"
	"github.com/timam/uttarawave-backend/internals/repositories"
	"net/http"
	"strconv"
)

type PackageHandler struct {
	repo repositories.PackageRepository
}

func NewPackageHandler(repo repositories.PackageRepository) *PackageHandler {
	return &PackageHandler{repo: repo}
}

func (h *PackageHandler) CreatePackage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pkg models.Package
		if err := c.ShouldBindJSON(&pkg); err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
			return
		}

		pkg.ID = uuid.New().String()
		pkg.IsActive = true // Set the package as active by default

		// Validate mandatory fields for all package types
		if pkg.Type == "" || pkg.Name == "" || pkg.Price == 0 {
			response.Error(c, http.StatusBadRequest, "Missing required fields", "Type, name, and price are required")
			return
		}

		// Validate package type and required fields
		switch pkg.Type {
		case models.CableTVPackage:
			if pkg.ChannelCount == nil || *pkg.ChannelCount == 0 || pkg.TVCount == nil || *pkg.TVCount == 0 {
				response.Error(c, http.StatusBadRequest, "Missing required fields for Cable TV package", "Channel count and TV count are required")
				return
			}
		case models.InternetPackage:
			if pkg.Bandwidth == nil || *pkg.Bandwidth == 0 || pkg.BandwidthType == nil || *pkg.BandwidthType == "" {
				response.Error(c, http.StatusBadRequest, "Missing required fields for Internet package", "Bandwidth and bandwidth type are required")
				return
			}
		default:
			response.Error(c, http.StatusBadRequest, "Invalid package type", "Specified package type is not supported")
			return
		}

		err := h.repo.CreatePackage(c.Request.Context(), &pkg)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to create package", err.Error())
			return
		}

		var responseData interface{}
		switch pkg.Type {
		case models.CableTVPackage:
			responseData = response.NewTVPackageResponse(&pkg)
		case models.InternetPackage:
			responseData = response.NewInternetPackageResponse(&pkg)
		}

		response.Success(c, http.StatusCreated, "Package created successfully", responseData)
	}
}

func (h *PackageHandler) GetAllPackages() gin.HandlerFunc {
	return func(c *gin.Context) {
		packageType := c.Query("type")
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

		packages, total, err := h.repo.GetAllPackages(c.Request.Context(), packageType, page, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.NewPackageResponse(http.StatusInternalServerError, "Failed to fetch packages", nil))
			return
		}

		var responsePackages []interface{}
		for _, pkg := range packages {
			switch pkg.Type {
			case models.CableTVPackage:
				responsePackages = append(responsePackages, response.NewTVPackageResponse(&pkg))
			case models.InternetPackage:
				responsePackages = append(responsePackages, response.NewInternetPackageResponse(&pkg))
			}
		}

		listResponse := response.NewPackageListResponse(responsePackages, total, page, pageSize)
		data := gin.H{"packages": listResponse}
		c.JSON(http.StatusOK, response.NewPackageResponse(http.StatusOK, "Packages retrieved successfully", data))
	}
}

func (h *PackageHandler) GetPackageByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		pkg, err := h.repo.GetPackageByID(c.Request.Context(), id)
		if err != nil {
			response.Error(c, http.StatusNotFound, "Package not found", err.Error())
			return
		}

		var responseData interface{}
		switch pkg.Type {
		case models.CableTVPackage:
			responseData = response.NewTVPackageResponse(pkg)
		case models.InternetPackage:
			responseData = response.NewInternetPackageResponse(pkg)
		default:
			response.Error(c, http.StatusInternalServerError, "Invalid package type", "Unknown package type")
			return
		}

		response.Success(c, http.StatusOK, "Package retrieved successfully", responseData)
	}
}

func (h *PackageHandler) DeletePackage() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		_, err := h.repo.GetPackageByID(c.Request.Context(), id)
		if err != nil {
			response.Error(c, http.StatusNotFound, "Package not found", err.Error())
			return
		}

		err = h.repo.DeletePackage(c.Request.Context(), id)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to delete package", err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Package deleted successfully", nil)
	}
}
